package auth_test

import (
	"log/slog"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/fallra1n/product-keeper/internal/core/auth"
	"github.com/fallra1n/product-keeper/internal/core/shared"
	mockauth "github.com/fallra1n/product-keeper/internal/mocks/auth"
	mockshared "github.com/fallra1n/product-keeper/internal/mocks/shared"
	"github.com/fallra1n/product-keeper/pkg/logging"
)

type RunAuthSuite struct {
	suite.Suite
	log *slog.Logger
}

func TestRunAuthSuite(t *testing.T) {
	suite.Run(t, new(RunAuthSuite))
}

func (s *RunAuthSuite) SetupTest() {
	s.log = logging.SetupLogger("local")
}

func (s *RunAuthSuite) TestCreateUser() {
	type fields struct {
		tx       *sqlx.Tx
		crypto   *mockshared.MockCrypto
		authRepo *mockauth.MockAuthRepo
	}

	var (
		mockUser       = auth.User{Name: "test name", Password: "test pass"}
		mockHashedUser = auth.User{Name: "test name", Password: "test hashed pass"}
	)

	testList := []struct {
		name    string
		prepare func(f *fields)
		args    auth.User
		err     error
	}{
		{
			name: "successful launch",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.crypto.EXPECT().HashPassword(mockUser.Password).Return(mockHashedUser.Password, nil),
					f.authRepo.EXPECT().CreateUser(f.tx, mockHashedUser).Return(nil),
				)
			},
			args: mockUser,
			err:  nil,
		},
		{
			name: "hashing failed",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.crypto.EXPECT().HashPassword(mockUser.Password).Return("", auth.ErrIncorrectPassword),
				)
			},
			args: mockUser,
			err:  shared.ErrInternal,
		},
		{
			name: "user already exist",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.crypto.EXPECT().HashPassword(mockUser.Password).Return(mockHashedUser.Password, nil),
					f.authRepo.EXPECT().CreateUser(f.tx, mockHashedUser).Return(auth.ErrUserAlreadyExist),
				)
			},
			args: mockUser,
			err:  auth.ErrUserAlreadyExist,
		},
	}

	for _, row := range testList {
		s.Run(row.name, func() {
			ctrl := gomock.NewController(s.T())
			defer ctrl.Finish()

			f := fields{
				tx:       &sqlx.Tx{},
				crypto:   mockshared.NewMockCrypto(ctrl),
				authRepo: mockauth.NewMockAuthRepo(ctrl),
			}
			if row.prepare != nil {
				row.prepare(&f)
			}

			service := auth.NewAuthService(
				s.log,
				f.authRepo,
				f.crypto,
			)

			err := service.CreateUser(f.tx, row.args)
			s.Equal(row.err, err)
		})
	}
}
