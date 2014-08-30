package tokens

import (
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/anachronistic/apns"
	"github.com/jmoiron/modl"
	_ "github.com/mattn/go-sqlite3"
)

type sqlTokenRepository struct {
	dbmap *modl.DbMap
}

// TokenSQLRepository returns a new sqlTokenRepository or panics if it cannot
func TokenSQLRepository(filename string) TokenRepository {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		panic("Error connecting to db: " + err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic("Error connecting to db: " + err.Error())
	}
	r := &sqlTokenRepository{
		dbmap: modl.NewDbMap(db, modl.SqliteDialect{}),
	}
	r.dbmap.TraceOn("", log.New(os.Stdout, "db: ", log.Lmicroseconds))
	r.dbmap.AddTable(Token{}).SetKeys(false, "token")
	r.dbmap.CreateTablesIfNotExists()
	return r
}

func (r *sqlTokenRepository) Save(token *Token) error {
	n, err := r.dbmap.Update(token)
	if err != nil {
		panic(err)
	}
	if n == 0 {
		err = r.dbmap.Insert(token)
	}
	return err
}

func (r *sqlTokenRepository) Delete(hash string) error {
	_, err := r.dbmap.Exec("DELETE FROM token WHERE hash=?", hash)
	return err
}

func (r *sqlTokenRepository) Load(user string) ([]*Token, error) {
	tokens := []*Token{}
	err := r.dbmap.Select(&tokens, "SELECT * FROM token WHERE user = ?", user)
	return tokens, err
}

func (r *sqlTokenRepository) List(users []string) ([]*Token, error) {
	tokens := []*Token{}
	err := r.dbmap.Select(&tokens, "SELECT * FROM token WHERE user IN (?)", strings.Join(users, ", "))
	return tokens, err
}

func (r *sqlTokenRepository) Push(users []string, message string, cert string, key string) error {
	tokens := []*Token{}
	err := r.dbmap.Select(&tokens, "SELECT * FROM token WHERE user IN (?)", strings.Join(users, ", "))
	if err != nil {
		return err
	}

	payload := apns.NewPayload()
	payload.Alert = message
	payload.Badge = 1 // TODO: Make this more accurate
	payload.Sound = "bingbong.aiff"

	client := apns.NewClient("gateway.sandbox.push.apple.com:2195", cert, key)

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
