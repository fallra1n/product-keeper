package services

import (
	"github.com/fallra1n/product-keeper/internal/domain/models"
	"golang.org/x/mod/sumdb/storage"
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
	GetProducts(username string, productName string, sortBy models.SortType) ([]models.Product, error)
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
