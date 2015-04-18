package render

import (
	"encoding/json"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/sessions"
)

// Context represents a renderer context
type Context map[string]interface{}

// MarshalPreparable can supply an alternative value in preparation for marshalling
type MarshalPreparable interface {
	MarshalPrepare() interface{}
}

var store = sessions.NewCookieStore([]byte("something-very-very-secret"))

var funcMap = template.FuncMap{
	"date": date,
}

// IsJSONRequest returns true if the request is XMLHttpRequest or ?json
func IsJSONRequest(r *http.Request) bool {
	json := r.FormValue("json")

	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" || json == "true" {
		return true
	}
	return false
}

// RegisterTemplateFunction registers functions to be used within templates
func RegisterTemplateFunction(name string, function interface{}) (alreadyRegistered bool) {
	_, alreadyRegistered = funcMap[name]
	funcMap[name] = function
	return alreadyRegistered
}

// Render returns a rendered template or JSON depending on the origin
// of the request
func Render(w http.ResponseWriter, r *http.Request, tmpl string, context map[string]interface{}) {
	if IsJSONRequest(r) {
		RenderJSON(w, context)
		return
	}
	RenderTemplate(w, tmpl, context)
}

// RenderTemplate renders a given template along with any data passed
func RenderTemplate(w http.ResponseWriter, tmpl string, context map[string]interface{}) {
	templates := template.New("").Funcs(funcMap)
	_, err := templates.ParseGlob("templates/*")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = templates.ExecuteTemplate(w, tmpl+".html", context)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// RenderJSON returns marshalled JSON
func RenderJSON(w http.ResponseWriter, context map[string]interface{}) {
	for key, data := range context {
		if preparable, ok := data.(MarshalPreparable); ok {
			context[key] = preparable.MarshalPrepare()
		}
	}

	obj, _ := json.MarshalIndent(context, "", "  ")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(obj)
}

// Redirect will redirect to given url. If it's an XHR request it will
// return JSON
func Redirect(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, http.StatusFound)
	return
}

func date(args ...interface{}) string {
	if len(args) == 2 {
		value := args[1].(time.Time)
		layout := args[0].(string)
		return value.Format(layout)
	}

	return args[0].(time.Time).Format("January 2, 2006 at 3:04PM")
}
