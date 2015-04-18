package auth

import (
	"time"

	"github.com/nathanborror/gommon/render"
)

// User defines a person in the system
type User struct {
	Key        string    `json:"key"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Created    time.Time `json:"created"`
	Modified   time.Time `json:"modified"`
	LastActive time.Time `json:"lastactive"`
}

// IsActive returns the duration of time since last active
func (u User) IsActive() bool {
	now := time.Now().UTC()
	dur := now.Sub(u.LastActive.UTC())
	return dur.Seconds() < 300.0
}

// GetURL returns an absolute url
func (u User) GetURL() string {
	return "/users/" + u.Key
}

// MarshalPrepare output
func (u User) MarshalPrepare() interface{} {
	return struct {
		User
		URL string `json:"url"`
	}{u, u.GetURL()}
}

func init() {
	_ = render.RegisterTemplateFunction("isAuthenticated", IsAuthenticated)
	_ = render.RegisterTemplateFunction("authenticatedUser", GetAuthenticatedUser)
	_ = render.RegisterTemplateFunction("authenticatedUserKey", GetAuthenticatedUserKey)
	_ = render.RegisterTemplateFunction("user", GetUser)
}
