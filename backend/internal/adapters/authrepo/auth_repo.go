package authrepo

import (
	"github.com/fallra1n/product-keeper/internal/adapters/authrepo/postgres"
)

// NewPostgresAuth ...
func NewPostgresAuth() *postgres.AuthRepository {
	return postgres.NewAuth()
}
