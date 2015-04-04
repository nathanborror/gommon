package auth

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // needed
	"github.com/nathanborror/gommon/settings"
)

type sqlRepository struct {
	db *sqlx.DB
}

// SqlRepository returns a new sqlRepository or panics if it cannot
func SqlRepository() UserRepository {
	settings := settings.NewSettings("settings.json")
	db, err := sqlx.Connect(settings.DriverName(), settings.DataSource())
	if err != nil {
		panic("Error connecting to db: " + err.Error())
	}

	repo := &sqlRepository{
		db: db,
	}

	datetime := "datetime"
	if settings.DriverName() == "postgres" {
		datetime = "timestamp with time zone"
	}

	schema := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS authuser (
		key text not null primary key,
		name text,
		email text,
		password text,
		created %s,
		modified %s,
		lastactive %s
	);`, datetime, datetime, datetime)

	_, err = repo.db.Exec(schema)
	return repo
}

func (r *sqlRepository) Get(key string) (*User, error) {
	user := User{}
	statement := fmt.Sprintf("SELECT * FROM authuser WHERE key = '%s'", key)
	err := r.db.Get(&user, statement)
	return &user, err
}

func (r *sqlRepository) GetWithPassword(email string, password string) (*User, error) {
	user := User{}
	statement := fmt.Sprintf("SELECT * FROM authuser WHERE email = '%s' AND password = '%s'", email, password)
	err := r.db.Get(&user, statement)
	return &user, err
}

func (r *sqlRepository) GetWithEmail(email string) (*User, error) {
	user := User{}
	statement := fmt.Sprintf("SELECT * FROM authuser WHERE email = '%s'", email)
	err := r.db.Get(&user, statement)
	return &user, err
}

func (r *sqlRepository) Update(user *User) error {
	user.Modified = time.Now()
	statement := `UPDATE authuser
		SET
			name = :name, email = :email, password = :password, modified = :modified, lastactive = :lastactive
		WHERE
			key = :key`
	_, err := r.db.NamedExec(statement, &user)
	return err
}

func (r *sqlRepository) Insert(user *User) error {
	user.Created = time.Now()
	user.Modified = time.Now()

	statement := `INSERT INTO authuser
			(key, name, email, password, created, modified, lastactive)
		VALUES
			(:key, :name, :email, :password, :created, :modified, :lastactive)`

	_, err := r.db.NamedExec(statement, &user)
	return err
}

func (r *sqlRepository) Delete(key string) error {
	statement := fmt.Sprintf("DELETE FROM authuser WHERE key = '%s'", key)
	_, err := r.db.Exec(statement)
	return err
}

func (r *sqlRepository) List(limit int) ([]*User, error) {
	users := []*User{}
	statement := fmt.Sprintf("SELECT * FROM authuser ORDER BY modified DESC LIMIT %d", limit)
	err := r.db.Select(&users, statement)
	return users, err
}

func (r *sqlRepository) UpdateLastActive(key string) error {
	statement := fmt.Sprintf("UPDATE authuser SET lastactive = '%s' WHERE key = '%s'", time.Now().String(), key)
	_, err := r.db.Exec(statement)
	return err
}
