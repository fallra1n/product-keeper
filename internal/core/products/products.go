package products

import (
	"errors"
	"log/slog"

	"github.com/jmoiron/sqlx"

	"github.com/fallra1n/product-keeper/internal/core/shared"
)

// ProductsService ...
type ProductsService struct {
	log  *slog.Logger
	date shared.DateTool

	productsRepo       ProductsRepo
	productsStatistics ProductsStatistics
}

// NewProductsService ...
func NewProductsService(
	log *slog.Logger,
	date shared.DateTool,

	productsRepo ProductsRepo,
	productsStatistics ProductsStatistics,
) *ProductsService {
	return &ProductsService{
		log:  log,
		date: date,

		productsRepo:       productsRepo,
		productsStatistics: productsStatistics,
	}
}

// CreateProduct ...
func (s *ProductsService) CreateProduct(tx *sqlx.Tx, product Product) (uint64, error) {
	product.CreatedAt = s.date.Now()

	id, err := s.productsRepo.CreateProduct(tx, product)
	if err != nil {
		s.log.Error("failed to create product", "error", err)
		return 0, shared.ErrInternal
	}

	s.log.Info("product has been created", "id", id)
	return id, nil
}

// FindProduct ...
func (s *ProductsService) FindProduct(tx *sqlx.Tx, id uint64, username string) (Product, error) {
	product, err := s.productsRepo.FindProduct(tx, id)
	if err != nil {
		s.log.Error("failed to find product by id", "error", err, "id", id)
		if errors.Is(err, shared.ErrNoData) {
			return Product{}, ErrProductNotFound
		}

		return Product{}, err
	}

	if product.OwnerName != username {
		s.log.Error(ErrPermissionDenied.Error(), "username", username, "id", id, "ownername", product.OwnerName)
		return Product{}, ErrPermissionDenied
	}

	if err := s.productsStatistics.Send(product); err != nil {
		s.log.Error("failed to send product view to statistics", "error", err, "id", id)
		return Product{}, shared.ErrInternal
	}

	return product, nil
}

// UpdateProduct ...
func (s *ProductsService) UpdateProduct(tx *sqlx.Tx, newProduct Product) (Product, error) {
	product, err := s.productsRepo.FindProduct(tx, newProduct.ID)
	if err != nil {
		s.log.Error("failed to find product by id", "error", err, "id", newProduct.ID)
		if errors.Is(err, shared.ErrNoData) {
			return Product{}, ErrProductNotFound
		}

		return Product{}, err
	}

	if product.OwnerName != newProduct.OwnerName {
		s.log.Error(ErrPermissionDenied.Error(), "username", newProduct.OwnerName, "id", newProduct.ID, "ownername", product.OwnerName)
		return Product{}, ErrPermissionDenied
	}

	data, err := s.productsRepo.UpdateProduct(tx, newProduct)
	if err != nil {
		s.log.Error("failed to update product", "error", err, "id", newProduct.ID)
		return Product{}, shared.ErrInternal
	}

	return data, nil
}

// DeleteProduct ...
func (s *ProductsService) DeleteProduct(tx *sqlx.Tx, id uint64, username string) error {
	product, err := s.productsRepo.FindProduct(tx, id)
	if err != nil {
		if errors.Is(err, shared.ErrNoData) {
			return ErrProductNotFound
		}

		return err
	}

	if product.OwnerName != username {
		return ErrPermissionDenied
	}

	return s.productsRepo.DeleteProduct(tx, id)
}

// FindProductList ...
func (s *ProductsService) FindProductList(tx *sqlx.Tx, username string, productName string, sortBy SortType) ([]Product, error) {
	return s.productsRepo.FindProductList(tx, username, productName, sortBy)
}
