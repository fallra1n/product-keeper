package services

import (
	"github.com/fallra1n/product-service/internal/domain/models"
	"github.com/fallra1n/product-service/internal/storage"
)

type Auth interface {
	CreateUser(user models.User) error
	LoginUser(user models.User) (string, error)
	ParseToken(token string) (string, error)
}

type Product interface {
	CreateProduct(product models.Product) (uint64, error)
	GetProductByID(id uint64, username string) (models.Product, error)
	UpdateProductByID(newProduct models.Product) (models.Product, error)
	DeleteProductByID(id uint64, username string) error
}

type Services interface {
	Auth
	Product
}

type services struct {
	Auth
	Product
}

func NewServices(storage storage.Storage) Services {
	return &services{
		NewAuthService(storage),
		NewProductService(storage),
	}
}
