package storage

import (
	"errors"

	"github.com/fallra1n/product-service/internal/domain/models"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserAlreadyExist = errors.New("user exists")
	ErrProductNotFound  = errors.New("product not found")
)

type Users interface {
	CreateUser(user models.User) error
	GetPasswordByName(name string) (string, error)
}

type Products interface {
	CreateProduct(product models.Product) (uint64, error)
	GetProductByID(id uint64) (models.Product, error)
	UpdateProductByID(id uint64, product models.Product) (models.Product, error)
	DeleteProductByID(id uint64) error
}

type Storage interface {
	CreateTables() error
	Users
	Products
}
