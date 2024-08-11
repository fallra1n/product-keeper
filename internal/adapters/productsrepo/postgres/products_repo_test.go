package postgres_test

import (
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"

	"github.com/fallra1n/product-keeper/config"
	"github.com/fallra1n/product-keeper/internal/adapters/productsrepo/postgres"
	"github.com/fallra1n/product-keeper/internal/core/auth"
	"github.com/fallra1n/product-keeper/internal/core/products"
	"github.com/fallra1n/product-keeper/pkg/access"
	"github.com/fallra1n/product-keeper/pkg/postgresdb"
)

type Suite struct {
	suite.Suite
	repo *postgres.ProductsRepository
	db   *sqlx.DB
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) SetupTest() {
	cfg := config.MustLoad()
	s.db = postgresdb.NewPostgresDB(access.PostgresTestConnect(cfg), cfg.Postgres.Timeout)
	s.repo = postgres.NewProducts()
}

func (s *Suite) TestCreateProduct() {
	now := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	mockUser := auth.NewUser("test name", "test password")
	mockProduct := products.NewProduct(0, "test product", 42, 42, "test name", now)

	s.Run("preparing data", func() {
		tx, err := s.db.Beginx()
		s.NoError(err)
		defer tx.Rollback()

		sqlQuery := `
			INSERT INTO auth$users (name, password)
			VALUES ($1, $2);
		`

		_, err = tx.Exec(sqlQuery, mockUser.Name, mockUser.Password)
		s.NoError(err)

		// creating product
		mockProduct.ID, err = s.repo.CreateProduct(tx, mockProduct)
		s.NoError(err)
		s.NotEqual(0, mockProduct.ID)

		s.Run("checking data", func() {
			sqlQuery := `
				SELECT * 
				FROM products 
				WHERE id = $1;
			`

			var data products.Product
			err := tx.Get(&data, sqlQuery, mockProduct.ID)
			s.NoError(err)

			data.CreatedAt = data.CreatedAt.In(time.UTC)
			s.Equal(mockProduct, data)
		})
	})
}
