package render

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"text/template"

	"github.com/gorilla/sessions"
	"github.com/russross/blackfriday"
)

var store = sessions.NewCookieStore([]byte("something-very-very-secret"))

var funcMap = template.FuncMap{
	"markdown":        markDowner,
	"initials":        initials,
	"isAuthenticated": isAuthenticated,
}

// Render returns a rendered template or JSON depending on the origin
// of the request
func Render(w http.ResponseWriter, r *http.Request, tmpl string, context interface{}) {
	// TOOD: There is a better way to detect XHR requests,
	// this is not that way.
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		RenderJSON(w, context)
		return
	}
	RenderTemplate(w, tmpl, context)
}

// RenderTemplate renders a given template along with any data passed
func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	templates := template.New("").Funcs(funcMap)
	_, err := templates.ParseGlob("templates/*")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// RenderJSON returns marshalled JSON
func RenderJSON(w http.ResponseWriter, data interface{}) {
	obj, _ := json.MarshalIndent(data, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(obj)
}

func markDowner(args ...interface{}) string {
	s := blackfriday.MarkdownCommon([]byte(fmt.Sprintf("%s", args...)))
	return string(s)
}

func initials(args ...interface{}) string {
	s := fmt.Sprintf("%s", args...)
	re := regexp.MustCompile("[^A-Z]")
	return re.ReplaceAllString(s, "")
}

func isAuthenticated(r *http.Request) bool {
	session, _ := store.Get(r, "authenticated-user")
	hash := session.Values["hash"]
	if hash != nil {
		return true
	}
	return false
}
