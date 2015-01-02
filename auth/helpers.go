package auth

import (
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/nathanborror/gommon/crypto"
)

var store = sessions.NewCookieStore([]byte("something-very-very-secret"))

// Authenticate authenticates and returns a user object
func Authenticate(email string, password string, w http.ResponseWriter, r *http.Request) (*User, error) {
	hash := crypto.PasswordHash(password)
	user, err := repo.GetWithPassword(email, hash)
	if err != nil {
		return nil, err
	}

	// Update session
	// TODO: Should just save the entire User object here
	if w != nil && r != nil {
		session, _ := store.Get(r, "authenticated-user")
		session.Values["key"] = user.Key
		session.Save(r, w)
	}

	return user, nil
}

// Deauthenticate clears authentication credentials from the client's session
func Deauthenticate(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "authenticated-user")
	session.Values["key"] = nil
	session.Save(r, w)
}

// IsAuthenticated checks whether someone has been authenticated
func IsAuthenticated(r *http.Request) bool {
	_, err := GetAuthenticatedUserKey(r)
	if err != nil {
		return false
	}
	return true
}

// GetAuthenticatedUserKey returns a User hash for an authenticated user
func GetAuthenticatedUserKey(r *http.Request) (string, error) {
	session, _ := store.Get(r, "authenticated-user")
	key := session.Values["key"]
	if key == nil {
		return "", errors.New("User is not authenticated")
	}
	return key.(string), nil
}

// GetAuthenticatedUser returns a User object for an authenticated user
func GetAuthenticatedUser(r *http.Request) (*User, error) {
	session, _ := store.Get(r, "authenticated-user")
	key := session.Values["key"]
	if key == nil {
		return nil, errors.New("User is not authenticated")
	}

	user, err := repo.Get(key.(string))
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUser returns a user for a given hash
func GetUser(key string) *User {
	user, _ := repo.Get(key)
	return user
}
