package products

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/fallra1n/product-keeper/internal/core/shared"
)

type ProductsService struct {
	db  *sqlx.DB
	log *slog.Logger

	productsRepo ProductsRepo
}

func NewProductsService(db *sqlx.DB, log *slog.Logger, productsRepo ProductsRepo) *ProductsService {
	return &ProductsService{
		db: db,
		log: log,
		productsRepo: productsRepo,
	}
}

func (s *ProductsService) CreateProduct(product Product) (uint64, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		s.log.Error(fmt.Sprintf("cannot start transaction: %s", err))
		return 0, err
	}
	defer tx.Rollback()

	product.CreatedAt = time.Now()
	id, err := s.productsRepo.CreateProduct(tx, product)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		s.log.Error(fmt.Sprintf("cannot commit transaction: %s", err))
		return 0, err
	}

	return id, nil
}

func (s *ProductsService) FindProduct(id uint64, username string) (Product, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		s.log.Error(fmt.Sprintf("cannot start transaction: %s", err))
		return Product{}, err
	}
	defer tx.Rollback()

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

	if err := tx.Commit(); err != nil {
		s.log.Error(fmt.Sprintf("cannot commit transaction: %s", err))
		return Product{}, err
	}

	return product, nil
}

func (s *ProductsService) UpdateProduct(newProduct Product) (Product, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		s.log.Error(fmt.Sprintf("cannot start transaction: %s", err))
		return Product{}, err
	}
	defer tx.Rollback()

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

	product, err = s.productsRepo.UpdateProduct(tx, newProduct)
	if err != nil {
		return Product{}, err
	}

	if err := tx.Commit(); err != nil {
		s.log.Error(fmt.Sprintf("cannot commit transaction: %s", err))
		return Product{}, err
	}

	return product, nil
}

func (s *ProductsService) DeleteProduct(id uint64, username string) error {
	tx, err := s.db.Beginx()
	if err != nil {
		s.log.Error(fmt.Sprintf("cannot start transaction: %s", err))
		return err
	}
	defer tx.Rollback()

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

	if err := s.productsRepo.DeleteProduct(tx, id); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		s.log.Error(fmt.Sprintf("cannot commit transaction: %s", err))
		return err
	}

	return nil
}

func (s *ProductsService) FindProductList(username string, productName string, sortBy SortType) ([]Product, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		s.log.Error(fmt.Sprintf("cannot start transaction: %s", err))
		return nil, err
	}
	defer tx.Rollback()

	productList, err := s.productsRepo.FindProductList(tx, username, productName, sortBy)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		s.log.Error(fmt.Sprintf("cannot commit transaction: %s", err))
		return nil, err
	}

	return productList, nil
}
