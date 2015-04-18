package graph

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/nathanborror/gommon/settings"
)

func unpack(s []string, vars ...*string) {
	for i, str := range s {
		*vars[i] = str
	}
}

/*
  Query lets you use a GraphQL like syntax to query your database. For example
  http://localhost?q=book(1){notes,editions} would return a JSON response with
  a book object and a list of notes and editions that match the book's key.

  Usage:

      query := "book(1){notes,editions}"

    	types := map[string]interface{}{
    		"book":      &Book{},
    		"edition":   &Edition{},
    		"note":      &Note{},
    	}

    	lists := map[string]interface{}{
    		"books":      &[]*Book{},
    		"editions":   &[]*Edition{},
    		"notes":      &[]*Note{},
    	}

      graph.Query(query, types, lists)

*/

func Query(query string, types map[string]interface{}, lists map[string]interface{}) map[string]interface{} {

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
	obj := types[node]
	statement := fmt.Sprintf("SELECT * FROM %s WHERE key = '%s'", node, key)
	err = db.Get(obj, statement)
	context[node] = obj

	// Retrieve children
	if rawChildren != "" {
		children := strings.Split(rawChildren, ",")
		for _, child := range children {
			objList := lists[child]

			re := regexp.MustCompile(`s$`)
			table := re.ReplaceAllString(child, "")

			statement := fmt.Sprintf("SELECT * FROM %s WHERE %s = '%s'", table, node, key)
			err = db.Select(objList, statement)
			context[child] = objList
		}
	}

	return context
}
