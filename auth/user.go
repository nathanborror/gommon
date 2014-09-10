package auth

import (
	"time"

	"github.com/nathanborror/gommon/render"
)

// User defines a person in the system
type User struct {
	Hash       string    `json:"hash"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Created    time.Time `json:"created"`
	Modified   time.Time `json:"modified"`
	LastActive time.Time `json:"last_active"`
}

// IsActive returns the duration of time since last active
func (u User) IsActive() bool {
	now := time.Now().UTC()
	dur := now.Sub(u.LastActive.UTC())
	return dur.Seconds() < 300.0
}

func init() {
	_ = render.RegisterTemplateFunction("isAuthenticated", IsAuthenticated)
	_ = render.RegisterTemplateFunction("authenticatedUser", GetAuthenticatedUser)
	_ = render.RegisterTemplateFunction("authenticatedUserHash", GetAuthenticatedUserHash)
	_ = render.RegisterTemplateFunction("user", GetUser)
}
