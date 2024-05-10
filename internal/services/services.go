package services

import "github.com/fallra1n/product-service/internal/storage"

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
