package tokens

import "time"

// Token defines a token for an auth.User
type Token struct {
	Token    string    `json:"token"`
	Platform string    `json:"platform"`
	User     string    `json:"user"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}
