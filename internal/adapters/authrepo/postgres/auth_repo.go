package postgres

import (
	"database/sql"

	"github.com/fallra1n/product-keeper/internal/core/auth"
	"github.com/fallra1n/product-keeper/internal/core/shared"
	"github.com/fallra1n/product-keeper/internal/domain/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type AuthRepository struct{}

func NewAuth() *AuthRepository {
	return &AuthRepository{}
}

func CreateTable(tx *sqlx.Tx) error {
	sqlQuery := `
		CREATE TABLE IF NOT EXISTS auth$users
		(
			name VARCHAR(255) NOT NULL UNIQUE,
		    password VARCHAR(255) NOT NULL
		);
	`

	_, err := tx.Exec(sqlQuery)
	return err
}

func (r *AuthRepository) CreateUser(tx *sqlx.Tx, user models.User) error {
	sqlQuery := `
		INSERT INTO auth$users (name, password)
		VALUES ($1, $2);
	`

	if _, err := tx.Exec(sqlQuery, user.Name, user.Password); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return auth.ErrUserAlreadyExist
			}
		}

		return err
	}

	return nil
}

func (r *AuthRepository) FindPassword(tx *sqlx.Tx, name string) (string, error) {
	sqlQuery := `
		SELECT password
		FROM auth$users
		WHERE name = $1;
	`

	var user models.User
	err := tx.Get(&user, sqlQuery, name)

	switch err {
	case sql.ErrNoRows:
		return "", shared.ErrNoData
	case nil:
		return user.Password, nil
	default:
		return "", err
	}
}
