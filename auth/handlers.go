package auth

import (
	"net/http"
	"strings"

	"github.com/nathanborror/gommon/crypto"
	"github.com/nathanborror/gommon/render"
)

var repo = SqlRepository()

// LoginHandler logs a user in
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		email := strings.TrimSpace(r.FormValue("email"))
		password := r.FormValue("password")
		u, err := Authenticate(email, password, w, r)
		if err != nil {
			u = &User{Email: email}
			render.RenderTemplate(w, "auth_register", map[string]interface{}{
				"request": r,
				"user":    u,
			})
			return
		}
		render.Redirect(w, r, "/")
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
		name := strings.TrimSpace(r.FormValue("name"))
		email := strings.TrimSpace(r.FormValue("email"))
		password := r.FormValue("password")

		// If email or password are blank then redirect to register page
		// TODO: provide a sensible error to people so they understand what
		// they did wrong.
		if email == "" || password == "" {
			render.Redirect(w, r, "/register")
			return
		}

		// Check to see if person already exists by attempting to log them in.
		passwordHash := crypto.PasswordHash(password)
		user, err := Authenticate(email, password, w, r)

		// If they do exist, redirect them home else create a new user and
		// log them into the site.
		if user != nil {
			render.Redirect(w, r, "/")
			return
		}

		key := crypto.UniqueHash(name)
		user = &User{Key: key, Name: name, Email: email, Password: passwordHash}
		err = repo.Insert(user)
		render.Check(err, w)

		// Auth user and redirect them
		user, _ = Authenticate(email, password, w, r)
		render.Redirect(w, r, "/")
		return
	}

	render.Render(w, r, "auth_register", nil)
}

// LoginRequired allows you to check for a logged in user on any handler
func LoginRequired(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key, err := GetAuthenticatedUserKey(r)

		_, err = repo.Get(key)
		if err != nil {
			Deauthenticate(w, r)
		}

		if IsAuthenticated(r) {
			fn(w, r)
			return
		}

		if render.IsJSONRequest(r) {
			render.RenderJSON(w, map[string]interface{}{
				"error": "User not authenticated",
			})
		} else {
			render.Redirect(w, r, "/register")
		}
	}
}
