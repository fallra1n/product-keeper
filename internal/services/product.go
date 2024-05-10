package services

import "github.com/fallra1n/product-service/internal/storage"

type Product interface {
}

type productService struct {
	storage storage.Storage
}

func NewProductService(storage storage.Storage) Product {
	return &productService{storage}
}
