package tokens

import (
	"net/http"

	"github.com/nathanborror/gommon/auth"
	"github.com/nathanborror/gommon/render"
)

var repo = TokenSQLRepository("db.sqlite3")
var userRepo = auth.AuthSQLRepository("db.sqlite3")

// SaveHandler saves a item
func SaveHandler(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetAuthenticatedUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	token := r.FormValue("token")
	platform := r.FormValue("platform")

	d := &Token{Token: token, Platform: platform, User: user.Key}
	err = repo.Insert(d)
	render.Check(err, w)

	http.Redirect(w, r, "/", http.StatusFound)
}
