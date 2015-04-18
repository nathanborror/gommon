// Package graph package is a very simplistic implementation of Facebook's GraphQL.
// This is a very naive implementation and doesn't by any means support all the
// great features GraphQL offers.
//
// Usage
//
//    query := "book(1){notes,editions}"
//
//   	objects := ObjectMap{
//   		"book":    Book{},
//   		"edition": Edition{},
//   		"note":    Note{},
//   	}
//
//    lists := ObjectListMap{
//   		"books":    []*Book{},
//   		"editions": []*Edition{},
//   		"notes":    []*Note{},
//   	}
//
//    graph.Query(query, objects, lists)
package graph

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/nathanborror/gommon/settings"
)

// ObjectMap is a map of all available objects
type ObjectMap map[string]interface{}

// ObjectListMap is a map of all available object lists
type ObjectListMap map[string]interface{}

func unpack(s []string, vars ...*string) {
	for i, str := range s {
		*vars[i] = str
	}
}

// Query lets you use a GraphQL like syntax to query your database. For example
// http://localhost?q=book(1){notes,editions} would return a JSON response with
// a book object and a list of notes and editions that match the book's key.
func Query(query string, objects ObjectMap, lists ObjectListMap) map[string]interface{} {

	// Database
	settings := settings.NewSettings("settings.json")
	db, err := sqlx.Connect(settings.DriverName(), settings.DataSource())
	if err != nil {
		log.Println(err)
	}

	// Parse query into {raw, node, key, children[]}
	var raw, node, key, rawChildren string
	re := regexp.MustCompile(`^(\w+)\(('?\w+'?)?\){([\w,]+)?}$`)
	matches := re.FindAllStringSubmatch(query, -1)
	unpack(matches[0], &raw, &node, &key, &rawChildren)

	context := map[string]interface{}{}

	// Retrieve primary node
	obj := objects[node]
	statement := fmt.Sprintf("SELECT * FROM %s WHERE key = '%s'", node, key)
	err = db.Get(&obj, statement)
	context[node] = obj

	// Retrieve children
	if rawChildren != "" {
		children := strings.Split(rawChildren, ",")
		for _, child := range children {
			list := lists[child]

			re := regexp.MustCompile(`s$`)
			table := re.ReplaceAllString(child, "")

			statement := fmt.Sprintf("SELECT * FROM %s WHERE %s = '%s'", table, node, key)
			err = db.Select(&list, statement)
			context[child] = list
		}
	}

	return context
}
