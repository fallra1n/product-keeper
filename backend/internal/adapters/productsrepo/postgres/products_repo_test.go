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
	"github.com/fallra1n/product-keeper/internal/core/shared"
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

func createUser(tx *sqlx.Tx, user auth.User) error {
	sqlQuery := `
		INSERT INTO auth$users (name, password)
		VALUES ($1, $2);
	`

	_, err := tx.Exec(sqlQuery, user.Name, user.Password)
	return err
}

func createProduct(tx *sqlx.Tx, product products.Product) (uint64, error) {
	sqlQuery := `
		INSERT INTO products (name, price, quantity, owner_name, created_at) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`

	row := tx.QueryRow(sqlQuery, product.Name, product.Price, product.Quantity, product.OwnerName, product.CreatedAt)

	var id uint64
	err := row.Scan(&id)

	return id, err
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

func (s *Suite) TestFindProduct() {
	now := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	mockUser := auth.NewUser("test name", "test password")
	mockProduct := products.NewProduct(0, "test product", 42, 42, "test name", now)

	s.Run("preparing data", func() {
		tx, err := s.db.Beginx()
		s.NoError(err)
		defer tx.Rollback()

		// creating user and product
		err = createUser(tx, mockUser)
		s.NoError(err)

		mockProduct.ID, err = createProduct(tx, mockProduct)
		s.NoError(err)
		s.NotEqual(0, mockProduct.ID)

		s.Run("checking data", func() {
			// call with non-existent id
			_, err := s.repo.FindProduct(tx, 0)
			s.ErrorIs(err, shared.ErrNoData)

			data, err := s.repo.FindProduct(tx, mockProduct.ID)
			s.NoError(err)

			data.CreatedAt = data.CreatedAt.In(time.UTC)
			s.Equal(mockProduct, data)
		})
	})
}

func (s *Suite) TestUpdateProduct() {
	now := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	mockUser := auth.NewUser("test name", "test password")
	mockProduct := products.NewProduct(0, "test product", 42, 42, "test name", now)

	s.Run("preparing data", func() {
		tx, err := s.db.Beginx()
		s.NoError(err)
		defer tx.Rollback()

		// creating product and user
		err = createUser(tx, mockUser)
		s.NoError(err)

		mockProduct.ID, err = createProduct(tx, mockProduct)
		s.NoError(err)
		s.NotEqual(0, mockProduct.ID)

		// update with non-existent id
		_, err = s.repo.UpdateProduct(tx, products.Product{ID: 0})
		s.ErrorIs(err, shared.ErrNoData)

		mockUpdatedProduct := products.NewProduct(mockProduct.ID, "test updated product", 43, 43, "test name", now)

		data, err := s.repo.UpdateProduct(tx, mockUpdatedProduct)
		s.NoError(err)

		data.CreatedAt = data.CreatedAt.In(time.UTC)
		s.Equal(mockUpdatedProduct, data)

		s.Run("checking data", func() {
			sqlQuery := `
				SELECT * 
				FROM products 
				WHERE id = $1;
			`

			var data products.Product
			err := tx.Get(&data, sqlQuery, mockUpdatedProduct.ID)
			s.NoError(err)

			data.CreatedAt = data.CreatedAt.In(time.UTC)
			s.Equal(mockUpdatedProduct, data)
		})
	})
}

func (s *Suite) TestDeleteProduct() {
	now := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	mockUser := auth.NewUser("test name", "test password")
	mockProduct := products.NewProduct(0, "test product", 42, 42, "test name", now)

	s.Run("preparing data", func() {
		tx, err := s.db.Beginx()
		s.NoError(err)
		defer tx.Rollback()

		// creating product and user
		err = createUser(tx, mockUser)
		s.NoError(err)

		mockProduct.ID, err = createProduct(tx, mockProduct)
		s.NoError(err)
		s.NotEqual(0, mockProduct.ID)

		sqlQuery := `
			SELECT * 
			FROM products 
			WHERE id = $1;
		`

		var data products.Product
		err = tx.Get(&data, sqlQuery, mockProduct.ID)
		s.NoError(err)
		s.NotEmpty(data)

		s.Run("checking data", func() {
			// delete product
			err = s.repo.DeleteProduct(tx, mockProduct.ID)
			s.NoError(err)

			_, err = tx.Exec(sqlQuery, mockProduct.ID)
			s.NoError(err)
		})
	})
}

func (s *Suite) TestFindProductList() {
	now := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	mockUser := auth.NewUser("test name", "test password")
	mockProduct1 := products.NewProduct(0, "test product1", 42, 42, "test name", now)
	mockProduct2 := products.NewProduct(0, "test product2", 43, 43, "test name", now.Add(5*time.Hour))

	s.Run("preparing data", func() {
		tx, err := s.db.Beginx()
		s.NoError(err)
		defer tx.Rollback()

		// creating product and user
		err = createUser(tx, mockUser)
		s.NoError(err)

		mockProduct1.ID, err = createProduct(tx, mockProduct1)
		s.NoError(err)
		s.NotEqual(0, mockProduct1.ID)

		mockProduct2.ID, err = createProduct(tx, mockProduct2)
		s.NoError(err)
		s.NotEqual(0, mockProduct2.ID)

		s.Run("checking data", func() {
			// get all products, sort by time
			data, err := s.repo.FindProductList(tx, "test name", "", products.LastCreate)
			s.NoError(err)

			data[0].CreatedAt = data[0].CreatedAt.In(time.UTC)
			data[1].CreatedAt = data[1].CreatedAt.In(time.UTC)
			s.Equal([]products.Product{mockProduct2, mockProduct1}, data)

			// get all products, sort by name
			data, err = s.repo.FindProductList(tx, "test name", "", products.Name)
			s.NoError(err)

			data[0].CreatedAt = data[0].CreatedAt.In(time.UTC)
			data[1].CreatedAt = data[1].CreatedAt.In(time.UTC)
			s.Equal([]products.Product{mockProduct1, mockProduct2}, data)

			// get named product
			data, err = s.repo.FindProductList(tx, "test name", "test product1", products.Name)
			s.NoError(err)

			data[0].CreatedAt = data[0].CreatedAt.In(time.UTC)
			s.Equal([]products.Product{mockProduct1}, data)

			data, err = s.repo.FindProductList(tx, "test name", "test product2", products.Name)
			s.NoError(err)

			data[0].CreatedAt = data[0].CreatedAt.In(time.UTC)
			s.Equal([]products.Product{mockProduct2}, data)
		})
	})
}
