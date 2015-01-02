package auth

// UserRepository holds all the methods needed to save load and list User objects.
type UserRepository interface {
	Get(key string) (*User, error)
	GetWithPassword(email string, password string) (*User, error)
	GetWithEmail(email string) (*User, error)
	Update(user *User) error
	Insert(user *User) error
	Delete(key string) error
	List(limit int) ([]*User, error)
	UpdateLastActive(key string) error
}
