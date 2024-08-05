package postgres_test

import (
	"database/sql"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"

	"github.com/fallra1n/product-keeper/config"
	"github.com/fallra1n/product-keeper/internal/adapters/authrepo/postgres"
	"github.com/fallra1n/product-keeper/internal/core/auth"
	"github.com/fallra1n/product-keeper/pkg/access"
	"github.com/fallra1n/product-keeper/pkg/postgresdb"
)

type Suite struct {
	suite.Suite
	repo *postgres.AuthRepository
	db   *sqlx.DB
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) SetupTest() {
	cfg := config.MustLoad()
	s.db = postgresdb.NewPostgresDB(access.PostgresTestConnect(cfg), cfg.Postgres.Timeout)
	s.repo = postgres.NewAuth()
}

func (s *Suite) TestCreateUser() {
	mockUser := auth.NewUser("test name", "test password")

	s.Run("Подготовка  данных", func() {
		tx, err := s.db.Beginx()
		s.NoError(err)
		defer tx.Rollback()

		// check that similar user doesn't exist
		sqlQuery := `
			SELECT password
			FROM auth$users
			WHERE name = $1;
		`

		var data string
		err = tx.QueryRow(sqlQuery, mockUser.Name).Scan(&data)
		s.ErrorIs(err, sql.ErrNoRows)

		// create user
		err = s.repo.CreateUser(tx, mockUser)
		s.NoError(err)

		s.Run("Загрузка  данных", func() {
			var data string
			err = tx.QueryRow(sqlQuery, mockUser.Name).Scan(&data)
			s.NoError(err)
			s.Equal(mockUser.Password, data)

			// create a user that already exists
			err = s.repo.CreateUser(tx, mockUser)
			s.ErrorIs(err, auth.ErrUserAlreadyExist)
		})
	})
}
