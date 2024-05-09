package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/fallra1n/product-service/internal/config"
	"github.com/fallra1n/product-service/internal/domain/models"
)

type Storage struct {
	db *sqlx.DB
}

func New(cfg *config.Config) (*Storage, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DBName)

	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Storage{db}, nil
}

func (s *Storage) CreateTables() error {
	createUser := `
		CREATE TABLE IF NOT EXISTS users 
		(
			name Varchar(255) NOT NULL UNIQUE,
		    password varchar(255) NOT NULL
		);`

	// TODO create products table

	if _, err := s.db.Exec(createUser); err != nil {
		return err
	}

	return nil
}

func (s *Storage) CreateUser(user models.User) error {
	query := `
		INSERT INTO users (name, password) 
		VALUES ($1, $2);`

	if _, err := s.db.Exec(query, user.Name, user.Password); err != nil {
		return err
	}

	return nil
}
