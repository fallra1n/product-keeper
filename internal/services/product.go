package services

import (
	"errors"

	"github.com/fallra1n/product-service/internal/domain/models"
	"github.com/fallra1n/product-service/internal/storage"
)

var (
	ErrProductNotFound  = errors.New("product not found")
	ErrPermissionDenied = errors.New("user does not have access to this product")
)

type Product interface {
	CreateProduct(product models.Product) (uint64, error)
	GetProductByID(id uint64, username string) (models.Product, error)
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

func (s *productService) GetProductByID(id uint64, username string) (models.Product, error) {
	product, err := s.storage.GetProductByID(id)
	if err != nil {
		if errors.Is(err, storage.ErrProductNotFound) {
			return models.Product{}, ErrProductNotFound
		}

		return models.Product{}, err
	}

	if product.OwnerName != username {
		return models.Product{}, ErrPermissionDenied
	}

	return product, nil
}

func (s *productService) UpdateProductByID(newProduct models.Product) (models.Product, error) {
	return models.Product{}, nil
}

func (s *productService) DeleteProductByID(id uint64) error {
	return nil
}
