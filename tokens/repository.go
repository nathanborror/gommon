package tokens

// TokenRepository holds all the methods needed to save, delete, load devices.
type TokenRepository interface {
	Get(user string) ([]*Token, error)
	Update(token *Token) error
	Insert(token *Token) error
	Delete(token string) error
	List(users []string) ([]*Token, error)
	Push(users []string, message string, cert string, key string) error
}
