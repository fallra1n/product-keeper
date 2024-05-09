package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/fallra1n/product-service/internal/config"
	"github.com/fallra1n/product-service/internal/domain/models"
	"github.com/fallra1n/product-service/internal/storage"
)

type postgres struct {
	db *sqlx.DB
}

func New(cfg *config.Config) (storage.Storage, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DBName)

	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &postgres{db}, nil
}

func (s *postgres) CreateTables() error {
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

func (s *postgres) CreateUser(user models.User) error {
	query := `
		INSERT INTO users (name, password) 
		VALUES ($1, $2);`

	if _, err := s.db.Exec(query, user.Name, user.Password); err != nil {
		return err
	}

	return nil
}

func (s *postgres) GetPasswordByName(name string) (string, error) {
	query := "SELECT password FROM users WHERE name = $1"

	var user models.User
	if err := s.db.Get(&user, query, name); err != nil {
		return "", err
	}

	return user.Password, nil
}
