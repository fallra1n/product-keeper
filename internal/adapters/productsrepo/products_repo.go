package productsrepo

import (
	"github.com/fallra1n/product-keeper/internal/adapters/productsrepo/postgres"
)

func NewPostgresProducts() *postgres.ProductsRepository {
	return postgres.NewProducts()
}
