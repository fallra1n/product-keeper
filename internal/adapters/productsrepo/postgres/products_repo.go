package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"golang.org/x/mod/sumdb/storage"

	"github.com/fallra1n/product-keeper/config"
	"github.com/fallra1n/product-keeper/internal/core/shared"
	"github.com/fallra1n/product-keeper/internal/domain/models"
)

type postgres struct {
	db *sqlx.DB
}

type ProductsRepository struct{}

func NewProducts() *ProductsRepository {
	return &ProductsRepository{}
}

func New(cfg *config.Config) (storage.Storage, error) {
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DBName)

	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &postgres{db}, nil
}

func (r *ProductsRepository) CreateTables() error {
	createUser := `
		CREATE TABLE IF NOT EXISTS users 
		(
			name VARCHAR(255) NOT NULL UNIQUE,
		    password VARCHAR(255) NOT NULL
		);`

	createProduct := `
		CREATE TABLE IF NOT EXISTS products
		(
		    id SERIAL PRIMARY KEY,
		    name VARCHAR(255) NOT NULL,
		    price INT NOT NULL,
		    quantity INT NOT NULL,
		    owner_name VARCHAR(255) NOT NULL,
		    created_at TIMESTAMP NOT NULL,
		    FOREIGN KEY (owner_name) REFERENCES users(name)
		);`

	if _, err := s.db.Exec(createUser); err != nil {
		return err
	}

	if _, err := s.db.Exec(createProduct); err != nil {
		return err
	}

	return nil
}

func (r *ProductsRepository) CreateProduct(product models.Product) (uint64, error) {
	query := `
		INSERT INTO products (name, price, quantity, owner_name, created_at) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;`

	row := s.db.QueryRow(query, product.Name, product.Price, product.Quantity, product.OwnerName, product.CreatedAt)

	var id uint64
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *ProductsRepository) GetProductByID(id uint64) (models.Product, error) {
	query := "SELECT * FROM products WHERE id = $1;"

	var product models.Product
	if err := s.db.Get(&product, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Product{}, shared.ErrNoData
		}
		return models.Product{}, err
	}

	return product, nil
}

func (r *ProductsRepository) UpdateProductByID(newProduct models.Product) (models.Product, error) {
	query := `
        UPDATE products
        SET name = $1, price = $2, quantity = $3
        WHERE id = $4
        RETURNING *;`

	var updated models.Product
	if err := s.db.Get(&updated, query, newProduct.Name, newProduct.Price, newProduct.Quantity, newProduct.ID); err != nil {
		return models.Product{}, err
	}

	return updated, nil
}

func (r *ProductsRepository) DeleteProductByID(id uint64) error {
	query := "DELETE FROM products WHERE id = $1;"
	if _, err := s.db.Exec(query, id); err != nil {
		return err
	}

	return nil
}

func (r *ProductsRepository) GetProducts(username string, productName string, sortBy models.SortType) ([]models.Product, error) {
	query := "SELECT * FROM products WHERE owner_name = $1"

	if productName != "" {
		query += fmt.Sprintf(" AND name = '%s'", productName)
	}

	switch sortBy {
	case models.Name:
		query += " ORDER BY name"
	case models.LastCreate:
		query += " ORDER BY created_at DESC"
	}

	query += ";"

	var products []models.Product
	if err := s.db.Select(&products, query, username); err != nil {
		return nil, err
	}

	return products, nil
}
