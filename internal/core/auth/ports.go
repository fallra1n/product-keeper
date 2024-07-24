package auth

import "github.com/fallra1n/product-keeper/internal/domain/models"

type Users interface {
	CreateUser(user models.User) error
	GetPasswordByName(name string) (string, error)
}
