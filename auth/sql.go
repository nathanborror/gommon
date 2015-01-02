package auth

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // needed
)

type sqlRepository struct {
	db *sqlx.DB
}

// AuthSQLRepository returns a new sqlRepository or panics if it cannot
func AuthSQLRepository(filename string) UserRepository {
	db, err := sqlx.Connect("sqlite3", filename)
	if err != nil {
		panic("Error connecting to db: " + err.Error())
	}

	repo := &sqlRepository{
		db: db,
	}

	schema := `CREATE TABLE IF NOT EXISTS user (
		key text not null primary key,
		name text,
		email text,
		password text,
		created datetime,
		modified datetime,
		lastactive datetime
	);`

	_, err = repo.db.Exec(schema)
	return repo
}

func (r *sqlRepository) Get(key string) (*User, error) {
	user := User{}
	err := r.db.Get(&user, `SELECT * FROM user WHERE key = ?`, key)
	return &user, err
}

func (r *sqlRepository) GetWithPassword(email string, password string) (*User, error) {
	user := User{}
	err := r.db.Get(&user, `SELECT * FROM user WHERE email = ? AND password = ?`, email, password)
	return &user, err
}

func (r *sqlRepository) GetWithEmail(email string) (*User, error) {
	user := User{}
	err := r.db.Get(&user, `SELECT * FROM user WHERE email = ?`, email)
	return &user, err
}

func (r *sqlRepository) Update(user *User) error {
	user.Modified = time.Now()
	statement := `UPDATE user SET name = :name, email = :email, password = :password, modified = :modified, lastactive = :lastactive WHERE key = :key`
	_, err := r.db.NamedExec(statement, &user)
	return err
}

func (r *sqlRepository) Insert(user *User) error {
	user.Created = time.Now()
	user.Modified = time.Now()

	statement := `INSERT INTO user
			(key, name, email, password, created, modified, lastactive)
		VALUES
			(:key, :name, :email, :password, :created, :modified, :lastactive)`

	_, err := r.db.NamedExec(statement, &user)
	return err
}

func (r *sqlRepository) Delete(key string) error {
	_, err := r.db.Exec(`DELETE FROM user WHERE key = ?`, key)
	return err
}

func (r *sqlRepository) List(limit int) ([]*User, error) {
	users := []*User{}
	err := r.db.Select(&users, `SELECT * FROM user ORDER BY modified DESC LIMIT ?`, limit)
	return users, err
}

func (r *sqlRepository) UpdateLastActive(key string) error {
	_, err := r.db.Exec(`UPDATE user SET lastactive = ? WHERE key = ?`, time.Now(), key)
	return err
}
