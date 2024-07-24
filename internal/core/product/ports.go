package product

import "github.com/fallra1n/product-keeper/internal/domain/models"

type Products interface {
	CreateProduct(product models.Product) (uint64, error)
	GetProductByID(id uint64) (models.Product, error)
	UpdateProductByID(newProduct models.Product) (models.Product, error)
	DeleteProductByID(id uint64) error
	GetProducts(username string, productName string, sortBy models.SortType) ([]models.Product, error)
}
