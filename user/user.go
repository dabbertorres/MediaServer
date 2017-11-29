package user

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type Config struct {
	Name string
	Session string
}

const (
	newUserQuery = "insert into users (name, password) values ($1, $2)"
	loginUserQuery = "select password from users where name = $1"
)

func New(db *sql.DB, username, password string) (Config, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return Config{}, err
	}
	
	_, err = db.Exec(newUserQuery, username, hashed)
	if err != nil {
		// TODO check if error is due to 'username' already in database
		return Config{}, err
	}

	return Config{
		Name: username,
	}, nil
}

func Login(db *sql.DB, username, password string) (Config, error) {
	row := db.QueryRow(loginUserQuery, username)
	
	hashed := ""
	if err := row.Scan(&hashed); err != nil {
		// TODO error filtering/handling
		return Config{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)); err != nil {
		return Config{}, err
	}

	// TODO make a session
	
	return Config{
		Name: username,
	}, nil
}

func (c Config) Logout() error {
	// TODO delete user's session from database

	// TODO shutdown user's station

	return nil
}
