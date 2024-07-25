package postgresdb

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func NewPostgresDB(postgresURL string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", postgresURL)
	if err != nil {
		log.Fatalln(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalln("база данных недоступна:", err)
	}

	// TODO setup connection parameters

	return db
}
