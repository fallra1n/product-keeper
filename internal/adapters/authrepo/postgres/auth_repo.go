package postgres

import (
	"database/sql"
	"errors"

	"github.com/fallra1n/product-keeper/internal/core/auth"
	"github.com/fallra1n/product-keeper/internal/core/shared"
	"github.com/fallra1n/product-keeper/internal/domain/models"

	"github.com/lib/pq"
)

type AuthRepository struct{}

func NewAuth() *AuthRepository {
	return &AuthRepository{}
}

func (r *AuthRepository) CreateUser(user models.User) error {
	query := `
		INSERT INTO users (name, password) 
		VALUES ($1, $2);`

	if _, err := s.db.Exec(query, user.Name, user.Password); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return auth.ErrUserAlreadyExist
			}
		}

		return err
	}

	return nil
}

func (r *AuthRepository) GetPasswordByName(name string) (string, error) {
	query := "SELECT password FROM users WHERE name = $1;"

	var user models.User
	if err := s.db.Get(&user, query, name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", shared.ErrNoData
		}
		return "", err
	}

	return user.Password, nil
}
