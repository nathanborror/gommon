package crypto

import (
	"crypto/sha1"
	"fmt"
	"io"
	"time"
)

// Hash returns a hash.
func Hash(text string) string {
	hasher := sha1.New()
	io.WriteString(hasher, text)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

// UniqueHash returns a sha1 unique hash based on a given string and the current date and time.
func UniqueHash(text string) string {
	time := time.Now().String()
	return Hash(text + time)
}

// PasswordHash returns a sha1 hashed password.
func PasswordHash(password string) string {
	return Hash(password)
}
