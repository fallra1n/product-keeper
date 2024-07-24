package authrepo

import (
	"github.com/fallra1n/product-keeper/internal/adapters/authrepo/postgres"
)

func NewPostgresAuth() *postgres.AuthRepository {
	return postgres.NewAuth()
}
