package auth

import (
	"github.com/fallra1n/product-keeper/internal/domain/models"
	"github.com/jmoiron/sqlx"
)

type Authrepo interface {
	CreateUser(tx *sqlx.Tx, user models.User) error
	FindPassword(tx *sqlx.Tx, name string) (string, error)
}
