package products

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"

	"github.com/fallra1n/product-keeper/internal/core/shared"
)

type ProductsService struct {
	log  *slog.Logger
	date shared.DateTool

	productsRepo       ProductsRepo
	productsStatistics ProductsStatistics
}

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

func (s *ProductsService) CreateProduct(tx *sqlx.Tx, product Product) (uint64, error) {
	product.CreatedAt = s.date.Now()

	id, err := s.productsRepo.CreateProduct(tx, product)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *ProductsService) FindProduct(tx *sqlx.Tx, id uint64, username string) (Product, error) {
	product, err := s.productsRepo.FindProduct(tx, id)
	if err != nil {
		if errors.Is(err, shared.ErrNoData) {
			return Product{}, ErrProductNotFound
		}

		return Product{}, err
	}

	if product.OwnerName != username {
		return Product{}, ErrPermissionDenied
	}

	if err := s.productsStatistics.Send(product); err != nil {
		s.log.Error(fmt.Sprintf("cannot send product view to statistics: %s", err))
		return Product{}, shared.ErrInternal
	}

	return product, nil
}

func (s *ProductsService) UpdateProduct(tx *sqlx.Tx, newProduct Product) (Product, error) {
	product, err := s.productsRepo.FindProduct(tx, newProduct.ID)
	if err != nil {
		if errors.Is(err, shared.ErrNoData) {
			return Product{}, ErrProductNotFound
		}

		return Product{}, err
	}

	if product.OwnerName != newProduct.OwnerName {
		return Product{}, ErrPermissionDenied
	}

	return s.productsRepo.UpdateProduct(tx, newProduct)
}

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

func (s *ProductsService) FindProductList(tx *sqlx.Tx, username string, productName string, sortBy SortType) ([]Product, error) {
	return s.productsRepo.FindProductList(tx, username, productName, sortBy)
}
