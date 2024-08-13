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
		jwt      *mockshared.MockJwt
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
				jwt:      mockshared.NewMockJwt(ctrl),
				authRepo: mockauth.NewMockAuthRepo(ctrl),
			}
			if row.prepare != nil {
				row.prepare(&f)
			}

			service := auth.NewAuthService(
				s.log,
				f.authRepo,
				f.crypto,
				f.jwt,
			)

			err := service.CreateUser(f.tx, row.args)
			s.Equal(row.err, err)
		})
	}
}

func (s *RunAuthSuite) TestLoginUser() {
	type fields struct {
		tx       *sqlx.Tx
		crypto   *mockshared.MockCrypto
		jwt      *mockshared.MockJwt
		authRepo *mockauth.MockAuthRepo
	}

	var (
		mockUser           = auth.User{Name: "test name", Password: "test pass"}
		mockHashedPassword = "test hashed pass"
		mockToken          = "test jwt token"
	)

	testList := []struct {
		name     string
		prepare  func(f *fields)
		args     auth.User
		expected string
		err      error
	}{
		{
			name: "successful launch",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.authRepo.EXPECT().FindPassword(f.tx, mockUser.Name).Return(mockHashedPassword, nil),
					f.crypto.EXPECT().CompareHashAndPassword(mockHashedPassword, mockUser.Password).Return(nil),
					f.jwt.EXPECT().GenerateToken(mockUser.Name).Return(mockToken, nil),
				)
			},
			args:     mockUser,
			expected: mockToken,
			err:      nil,
		},
		{
			name: "user not found",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.authRepo.EXPECT().FindPassword(f.tx, mockUser.Name).Return("", auth.ErrUserNotFound),
				)
			},
			args:     mockUser,
			expected: "",
			err:      auth.ErrUserNotFound,
		},
		{
			name: "incorrect password",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.authRepo.EXPECT().FindPassword(f.tx, mockUser.Name).Return(mockHashedPassword, nil),
					f.crypto.EXPECT().CompareHashAndPassword(mockHashedPassword, mockUser.Password).Return(shared.ErrNoData),
				)
			},
			args:     mockUser,
			expected: "",
			err:      auth.ErrIncorrectPassword,
		},
		{
			name: "failed to generate token",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.authRepo.EXPECT().FindPassword(f.tx, mockUser.Name).Return(mockHashedPassword, nil),
					f.crypto.EXPECT().CompareHashAndPassword(mockHashedPassword, mockUser.Password).Return(nil),
					f.jwt.EXPECT().GenerateToken(mockUser.Name).Return("", shared.ErrNoData),
				)
			},
			args:     mockUser,
			expected: "",
			err:      shared.ErrInternal,
		},
	}

	for _, row := range testList {
		s.Run(row.name, func() {
			ctrl := gomock.NewController(s.T())
			defer ctrl.Finish()

			f := fields{
				tx:       &sqlx.Tx{},
				crypto:   mockshared.NewMockCrypto(ctrl),
				jwt:      mockshared.NewMockJwt(ctrl),
				authRepo: mockauth.NewMockAuthRepo(ctrl),
			}
			if row.prepare != nil {
				row.prepare(&f)
			}

			service := auth.NewAuthService(
				s.log,
				f.authRepo,
				f.crypto,
				f.jwt,
			)

			data, err := service.LoginUser(f.tx, row.args)
			s.Equal(row.err, err)
			s.Equal(row.expected, data)
		})
	}
}
