package tokens

// TokenRepository holds all the methods needed to save, delete, load devices.
type TokenRepository interface {
	Load(user string) ([]*Token, error)
	Delete(hash string) error
	Save(token *Token) error
	List(users []string) ([]*Token, error)
	Push(users []string, message string, cert string, key string) error
}
