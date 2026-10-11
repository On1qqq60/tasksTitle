package repository

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type Repository interface {
	Add(uuid uuid.UUID, name string, password string) error
	IsCorrect(name string, password string) (bool, error)
}

type DB struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &DB{db: db}
}

func (d *DB) Add(uuid uuid.UUID, name string, password string) error {
	var findName string
	err := d.db.QueryRow(`SELECT name FROM "user" WHERE name = $1`, name).Scan(&findName)
	if err != nil {
		_, err = d.db.Exec(`INSERT INTO "user" (id, name, password) VALUES ($1, $2, $3)`, uuid, name, password)
		return err
	} else {
		return errors.New("Имя уже есть в базе")
	}
}

func (d *DB) IsCorrect(name string, password string) (bool, error) {
	var findPass string
	err := d.db.QueryRow(`SELECT password FROM "user" WHERE name = $1`, name).Scan(&findPass)
	if err != nil {
		return false, err
	}
	if password != findPass {
		err = errors.New("Неверный пароль")
		return false, err
	}

	return true, nil
}

func InitSgemaUsers(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS "user"(
		id 	 		UUID PRIMARY KEY,
		name 		TEXT NOT NULL,
		password 	TEXT NOT NULL
	)`)

	return err
}
