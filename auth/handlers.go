package auth

import (
	"net/http"

	"github.com/nathanborror/gommon/render"
)

var repo = AuthSQLRepository("db.sqlite3")

// LoginHandler logs a user in
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")
		u, err := Authenticate(email, password, w, r)
		if err != nil {
			u = &User{Email: email}
			render.RenderTemplate(w, "auth_register", map[string]interface{}{
				"user": u,
			})
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	render.RenderTemplate(w, "auth_login", nil)
}

// LogoutHandler signs a user out
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	Deauthenticate(w, r)
	render.RenderTemplate(w, "auth_logout", nil)
}

// RegisterHandler registers a new user
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" && password == "" {
			http.Redirect(w, r, "/register", http.StatusFound)
			return
		}

		hash := GenerateUserHash(email)
		passwordHash := GeneratePasswordHash(password)
		u, err := Authenticate(email, password, w, r)

		// If user already exists, sign them in
		if err == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		u = &User{Email: email, Password: passwordHash, Hash: hash, Name: name}
		err = repo.Save(u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Auth user and redirect them
		u, _ = Authenticate(email, password, w, r)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	render.RenderTemplate(w, "auth_register", nil)
}

// AuthHandler allows you to check for a logged in user on any handler
func AuthHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if IsAuthenticated(r) == true {
			fn(w, r)
			return
		}
		http.Redirect(w, r, "/register", http.StatusFound)
	}
}
