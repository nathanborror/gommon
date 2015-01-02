package tokens

import (
	"log"
	"strings"
	"time"

	"github.com/anachronistic/apns"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // Needed
)

type sqlRepository struct {
	db *sqlx.DB
}

// TokenSQLRepository returns a new sqlRepository or panics if it cannot
func TokenSQLRepository(filename string) TokenRepository {
	db, err := sqlx.Connect("sqlite3", filename)
	if err != nil {
		panic("Error connecting to db: " + err.Error())
	}

	repo := &sqlRepository{
		db: db,
	}

	schema := `CREATE TABLE IF NOT EXISTS token (
		token text not null primary key,
		platform text NULL,
		user text NULL,
		created datetime,
		modified datetime
	);`

	_, err = repo.db.Exec(schema)
	return repo
}

func (r *sqlRepository) Get(user string) ([]*Token, error) {
	tokens := []*Token{}
	err := r.db.Select(&tokens, `SELECT * FROM token WHERE user = ?`, user)
	return tokens, err
}

func (r *sqlRepository) Update(token *Token) error {
	token.Modified = time.Now()
	statement := `UPDATE token SET token = :token, platform = :platform, modified = :modified WHERE user = :user`
	_, err := r.db.NamedExec(statement, &token)
	return err
}

func (r *sqlRepository) Insert(token *Token) error {
	token.Created = time.Now()
	token.Modified = time.Now()

	statement := `INSERT INTO token
			(token, platform, user, created, modified)
		VALUES
			(:token, :platform, :user, :created, :modified)`
	_, err := r.db.NamedExec(statement, &token)
	return err
}

func (r *sqlRepository) Delete(token string) error {
	_, err := r.db.Exec("DELETE FROM token WHERE token = ?", token)
	return err
}

func (r *sqlRepository) List(users []string) ([]*Token, error) {
	tokens := []*Token{}
	err := r.db.Select(&tokens, "SELECT * FROM token WHERE user IN (?)", strings.Join(users, ", "))
	return tokens, err
}

func (r *sqlRepository) Push(users []string, message string, cert string, key string) error {
	tokens := []*Token{}
	err := r.db.Select(&tokens, "SELECT * FROM token WHERE user IN (?)", strings.Join(users, ", "))
	if err != nil {
		return err
	}

	payload := apns.NewPayload()
	payload.Alert = message
	payload.Badge = 1 // TODO: Make this more accurate
	payload.Sound = "bingbong.aiff"

	client := apns.NewClient("gateway.push.apple.com:2195", cert, key)

	for _, token := range tokens {
		pn := apns.NewPushNotification()
		pn.DeviceToken = token.Token
		pn.AddPayload(payload)
		resp := client.Send(pn)

		alert, _ := pn.PayloadString()
		if resp.Error != nil {
			log.Println("APNS Error: ", resp.Error)
		} else {
			log.Println("APNS Alert: ", alert)
			log.Println("APNS Success: ", resp.Success)
		}
	}

	return nil
}
