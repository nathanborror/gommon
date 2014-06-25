package auth

// UserRepository holds all the methods needed to save load and list User objects.
type UserRepository interface {
	Load(hash string) (*User, error)
	LoadWithPassword(email string, password string) (*User, error)
	LoadWithEmail(email string) (*User, error)
	Save(user *User) error
	List(limit int) ([]*User, error)
}
