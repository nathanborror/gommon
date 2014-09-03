package auth

import (
	"time"

	"github.com/nathanborror/gommon/render"
)

// User defines a person in the system
type User struct {
	Hash     string    `json:"hash"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

func init() {
	_ = render.RegisterTemplateFunction("isAuthenticated", IsAuthenticated)
	_ = render.RegisterTemplateFunction("authenticatedUser", GetAuthenticatedUser)
	_ = render.RegisterTemplateFunction("user", GetUser)
}
