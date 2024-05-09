package storage

import "github.com/fallra1n/product-service/internal/domain/models"

type Users interface {
	CreateTables() error
	CreateUser(user models.User) error
	GetPasswordByName(name string) (string, error)
}

type Products interface{}

type Storage interface {
	Users
	Products
}
