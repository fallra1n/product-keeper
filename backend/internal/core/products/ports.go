package products

import (
	"github.com/jmoiron/sqlx"
)

// ProductsRepo ...
type ProductsRepo interface {
	CreateProduct(tx *sqlx.Tx, product Product) (uint64, error)
	FindProduct(tx *sqlx.Tx, id uint64) (Product, error)
	UpdateProduct(tx *sqlx.Tx, newProduct Product) (Product, error)
	DeleteProduct(tx *sqlx.Tx, id uint64) error
	FindProductList(tx *sqlx.Tx, username string, productName string, sortBy SortType) ([]Product, error)
}

// ProductsStatistics ...
type ProductsStatistics interface {
	Send(p Product) error
}
