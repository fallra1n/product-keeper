package postgres

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/fallra1n/product-keeper/internal/core/auth"
	"github.com/fallra1n/product-keeper/internal/core/shared"
)

// AuthRepository ...
type AuthRepository struct{}

// NewAuth constructor for AuthRepository
func NewAuth() *AuthRepository {
	return &AuthRepository{}
}

// CreateUser ...
func (r *AuthRepository) CreateUser(tx *sqlx.Tx, user auth.User) error {
	sqlQuery := `
		INSERT INTO auth$users (name, password)
		VALUES ($1, $2);
	`

	if _, err := tx.Exec(sqlQuery, user.Name, user.Password); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "02000":
				return shared.ErrNoData
			case "23505":
				return auth.ErrUserAlreadyExist
			default:
			}
		}

		return err
	}

	return nil
}

// FindPassword ...
func (r *AuthRepository) FindPassword(tx *sqlx.Tx, name string) (string, error) {
	sqlQuery := `
		SELECT password
		FROM auth$users
		WHERE name = $1;
	`

	var user auth.User
	err := tx.Get(&user, sqlQuery, name)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return "", auth.ErrUserNotFound
	case err == nil:
		return user.Password, nil
	default:
		return "", err
	}
}
