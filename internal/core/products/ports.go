package products

import (
	"github.com/fallra1n/product-keeper/internal/domain/models"
	"github.com/jmoiron/sqlx"
)

type ProductsRepo interface {
	CreateProduct(tx *sqlx.Tx, product models.Product) (uint64, error)
	FindProduct(tx *sqlx.Tx, id uint64) (models.Product, error)
	UpdateProduct(tx *sqlx.Tx, newProduct models.Product) (models.Product, error)
	DeleteProduct(tx *sqlx.Tx, id uint64) error
	FindProductList(tx *sqlx.Tx, username string, productName string, sortBy models.SortType) ([]models.Product, error)
}
