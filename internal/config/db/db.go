package db

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var DB *sql.DB

type Database struct {
	*sql.DB
}

func (d Database) New(adr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", adr)
	if err != nil {
		return nil, err
	}
	DB = db

	return db, nil

}
