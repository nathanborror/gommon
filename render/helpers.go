package render

import (
	"encoding/json"
	"net/http"
)

// Check is a shorthand function for error checking
func Check(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ReturnJSON returns json
func ReturnJSON(context map[string]interface{}) string {
	obj, _ := json.MarshalIndent(context, "", "  ")
	return string(obj[:])
}
