package postgresdb

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgresDB(postgresURL string, timeout time.Duration) *sqlx.DB {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	timeoutChan := time.After(timeout)

	for {
		select {
		case <-ticker.C:
			db, err := sqlx.Connect("postgres", postgresURL)
			if err != nil {
				continue
			}

			if err := db.Ping(); err != nil {
				continue
			}

			log.Println("successful connection to the postgres")
			return db

		case <-timeoutChan:
			log.Fatalln("cannot connect to the postgres, time's up")
		}
	}

	// TODO setup connection parameters
}
