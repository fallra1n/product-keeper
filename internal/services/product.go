package services

import (
	"errors"
	"time"

	"github.com/fallra1n/product-keeper/internal/domain/models"
	"github.com/fallra1n/product-keeper/internal/storage"
)

var (
	ErrProductNotFound  = errors.New("product not found")
	ErrPermissionDenied = errors.New("user does not have access to this product")
)

type productService struct {
	storage storage.Storage
}

func NewProductService(storage storage.Storage) Product {
	return &productService{storage}
}

func (s *productService) CreateProduct(product models.Product) (uint64, error) {
	product.CreatedAt = time.Now()
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
	product, err := s.storage.GetProductByID(newProduct.ID)
	if err != nil {
		if errors.Is(err, storage.ErrProductNotFound) {
			return models.Product{}, ErrProductNotFound
		}

		return models.Product{}, err
	}

	if product.OwnerName != newProduct.OwnerName {
		return models.Product{}, ErrPermissionDenied
	}

	return s.storage.UpdateProductByID(newProduct)
}

func (s *productService) DeleteProductByID(id uint64, username string) error {
	product, err := s.storage.GetProductByID(id)
	if err != nil {
		if errors.Is(err, storage.ErrProductNotFound) {
			return ErrProductNotFound
		}

		return err
	}

	if product.OwnerName != username {
		return ErrPermissionDenied
	}

	return s.storage.DeleteProductByID(id)
}

func (s *productService) GetProducts(username string, productName string, sortBy models.SortType) ([]models.Product, error) {
	return s.storage.GetProducts(username, productName, sortBy)
}
