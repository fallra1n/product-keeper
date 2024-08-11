package postgres

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/fallra1n/product-keeper/internal/core/products"
	"github.com/fallra1n/product-keeper/internal/core/shared"
)

type ProductsRepository struct{}

func NewProducts() *ProductsRepository {
	return &ProductsRepository{}
}

func (r *ProductsRepository) CreateProduct(tx *sqlx.Tx, product products.Product) (uint64, error) {
	sqlQuery := `
		INSERT INTO products (name, price, quantity, owner_name, created_at) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`

	row := tx.QueryRow(sqlQuery, product.Name, product.Price, product.Quantity, product.OwnerName, product.CreatedAt)

	var id uint64
	err := row.Scan(&id)

	switch err {
	case sql.ErrNoRows:
		return 0, shared.ErrNoData
	case nil:
		return id, nil
	default:
		return 0, err
	}
}

func (r *ProductsRepository) FindProduct(tx *sqlx.Tx, id uint64) (products.Product, error) {
	sqlQuery := `
		SELECT * 
		FROM products 
		WHERE id = $1;
	`

	var data products.Product
	err := tx.Get(&data, sqlQuery, id)

	switch err {
	case sql.ErrNoRows:
		return products.Product{}, shared.ErrNoData
	case nil:
		return data, nil
	default:
		return products.Product{}, err
	}
}

func (r *ProductsRepository) UpdateProduct(tx *sqlx.Tx, newProduct products.Product) (products.Product, error) {
	sqlQuery := `
    UPDATE products
    SET name = $1, price = $2, quantity = $3
    WHERE id = $4
    RETURNING *;
	`

	var data products.Product
	err := tx.Get(&data, sqlQuery, newProduct.Name, newProduct.Price, newProduct.Quantity, newProduct.ID)

	switch err {
	case sql.ErrNoRows:
		return products.Product{}, shared.ErrNoData
	case nil:
		return data, nil
	default:
		return products.Product{}, err
	}
}

func (r *ProductsRepository) DeleteProduct(tx *sqlx.Tx, id uint64) error {
	sqlQuery := `
		DELETE 
		FROM products
		WHERE id = $1;
	`

	_, err := tx.Exec(sqlQuery, id)
	return err
}

func (r *ProductsRepository) FindProductList(tx *sqlx.Tx, username string, productName string, sortBy products.SortType) ([]products.Product, error) {
	sqlQuery := `
		SELECT * 
		FROM products 
		WHERE owner_name = $1
	`

	if productName != "" {
		sqlQuery += fmt.Sprintf(" AND name = '%s'", productName)
	}

	switch sortBy {
	case products.Name:
		sqlQuery += " ORDER BY name"
	case products.LastCreate:
		sqlQuery += " ORDER BY created_at DESC"
	default:
	}

	sqlQuery += ";"

	var data []products.Product
	err := tx.Select(&data, sqlQuery, username)

	switch err {
	case sql.ErrNoRows:
		return nil, shared.ErrNoData
	case nil:
		return data, nil
	default:
		return nil, err
	}
}
