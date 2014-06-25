package auth

import "time"

// User defines a person in the system
type User struct {
	Hash     string    `json:"hash"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}
