package auth

import (
	"crypto/md5"
	"fmt"
	"hash/fnv"
	"io"
	"net/http"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("something-very-very-secret"))

// GeneratePasswordHash returns a hashed password
func GeneratePasswordHash(password string) (hash string) {
	hasher := md5.New()
	io.WriteString(hasher, password)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

// GenerateUserHash returns a hash that represents a unique user ID
func GenerateUserHash(s string) (hash string) {
	hasher := fnv.New32a()
	io.WriteString(hasher, s)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

// Authenticate authenticates and returns a user object
func Authenticate(email string, password string, w http.ResponseWriter, r *http.Request) (*User, error) {
	hash := GeneratePasswordHash(password)
	u, err := repo.LoadWithPassword(email, hash)
	if err != nil {
		return nil, err
	}

	// Update session
	// TODO: Should just save the entire User object here
	if w != nil && r != nil {
		session, _ := store.Get(r, "authenticated-user")
		session.Values["hash"] = u.Hash
		session.Save(r, w)
	}

	return u, nil
}

// Deauthenticate clears authentication credentials from the client's session
func Deauthenticate(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "authenticated-user")
	session.Values["hash"] = nil
	session.Save(r, w)
}

// IsAuthenticated checks whether someone has been authenticated
func IsAuthenticated(r *http.Request) bool {
	_, err := GetAuthenticatedUser(r)
	if err != nil {
		return false
	}
	return true
}

// GetAuthenticatedUser returns a User object for an authenticated user
func GetAuthenticatedUser(r *http.Request) (*User, error) {
	session, _ := store.Get(r, "authenticated-user")
	hash := session.Values["hash"]
	u, err := repo.Load(hash.(string))
	if err != nil {
		return nil, err
	}
	return u, nil
}
