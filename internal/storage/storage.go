package storage

import (
	"errors"

	"github.com/fallra1n/product-service/internal/domain/models"
)

var (
	ErrUserNotFound     = errors.New("url not found")
	ErrUserAlreadyExist = errors.New("url exists")
)

type Users interface {
	CreateTables() error
	CreateUser(user models.User) error
	GetPasswordByName(name string) (string, error)
}

type Products interface {
	CreateProduct(product models.Product) (uint64, error)
}

type Storage interface {
	Users
	Products
}
