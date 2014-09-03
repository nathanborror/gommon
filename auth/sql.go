package auth

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/modl"
	_ "github.com/mattn/go-sqlite3"
)

type sqlUserRepository struct {
	dbmap *modl.DbMap
}

// AuthSQLRepository returns a new sqlUserRepository or panics if it cannot
func AuthSQLRepository(filename string) UserRepository {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		panic("Error connecting to db: " + err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic("Error connecting to db: " + err.Error())
	}
	r := &sqlUserRepository{
		dbmap: modl.NewDbMap(db, modl.SqliteDialect{}),
	}
	r.dbmap.TraceOn("", log.New(os.Stdout, "db: ", log.Lmicroseconds))
	r.dbmap.AddTable(User{}).SetKeys(false, "hash")
	r.dbmap.CreateTablesIfNotExists()
	return r
}

func (r *sqlUserRepository) Load(hash string) (*User, error) {
	objects := []*User{}
	err := r.dbmap.Select(&objects, "SELECT * FROM user WHERE hash=?", hash)
	if len(objects) != 1 {
		return nil, fmt.Errorf("expected 1 object, got %d", len(objects))
	}
	return objects[0], err
}

func (r *sqlUserRepository) LoadWithPassword(email string, password string) (*User, error) {
	objects := []*User{}
	err := r.dbmap.Select(&objects, "SELECT * FROM user WHERE email=? AND password=?", email, password)
	if len(objects) != 1 {
		return nil, fmt.Errorf("expected 1 object, got %d", len(objects))
	}
	return objects[0], err
}

func (r *sqlUserRepository) LoadWithEmail(email string) (*User, error) {
	objects := []*User{}
	err := r.dbmap.Select(&objects, "SELECT * FROM user WHERE email=?", email)
	if len(objects) != 1 {
		return nil, fmt.Errorf("expected 1 object, got %d", len(objects))
	}
	return objects[0], err
}

func (r *sqlUserRepository) Save(user *User) error {
	n, err := r.dbmap.Update(user)
	if err != nil {
		panic(err)
	}
	if n == 0 {
		err = r.dbmap.Insert(user)
	}
	return err
}

func (r *sqlUserRepository) Delete(hash string) error {
	_, err := r.dbmap.Exec("DELETE FROM user WHERE hash=?", hash)
	return err
}

func (r *sqlUserRepository) List(limit int) ([]*User, error) {
	objects := []*User{}
	err := r.dbmap.Select(&objects, "SELECT * FROM user ORDER BY modified DESC LIMIT ?", limit)
	return objects, err
}

func (r *sqlUserRepository) Ping(user *User) {
	user.LastActive = time.Now()
	r.dbmap.Update(user)
}
