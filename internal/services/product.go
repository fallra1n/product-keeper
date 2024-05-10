package services

import (
	"github.com/fallra1n/product-service/internal/domain/models"
	"github.com/fallra1n/product-service/internal/storage"
)

type Product interface {
	CreateProduct(product models.Product) (uint64, error)
	GetProductByID(id uint64) (models.Product, error)
	UpdateProductByID(newProduct models.Product) (models.Product, error)
	DeleteProductByID(id uint64) error
}

type productService struct {
	storage storage.Storage
}

func NewProductService(storage storage.Storage) Product {
	return &productService{storage}
}

func (s *productService) CreateProduct(product models.Product) (uint64, error) {
	return s.storage.CreateProduct(product)
}

func (s *productService) GetProductByID(id uint64) (models.Product, error) {
	return models.Product{}, nil
}

func (s *productService) UpdateProductByID(newProduct models.Product) (models.Product, error) {
	return models.Product{}, nil
}

func (s *productService) DeleteProductByID(id uint64) error {
	return nil
}
