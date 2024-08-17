package products_test

import (
	"log/slog"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/fallra1n/product-keeper/internal/core/products"
	"github.com/fallra1n/product-keeper/internal/core/shared"
	mockproducts "github.com/fallra1n/product-keeper/internal/mocks/products"
	mockshared "github.com/fallra1n/product-keeper/internal/mocks/shared"
	"github.com/fallra1n/product-keeper/pkg/logging"
)

type RunProductsSuite struct {
	suite.Suite
	log *slog.Logger
}

func TestRunProductsSuite(t *testing.T) {
	suite.Run(t, new(RunProductsSuite))
}

func (s *RunProductsSuite) SetupTest() {
	s.log = logging.SetupLogger("local")
}

func (s *RunProductsSuite) TestCreateProduct() {
	type fields struct {
		tx   *sqlx.Tx
		date *mockshared.MockDateTool

		productsRepo       *mockproducts.MockProductsRepo
		productsStatistics *mockproducts.MockProductsStatistics
	}

	var (
		now = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

		mockProductID uint64 = 567
	)

	testList := []struct {
		name         string
		prepare      func(f *fields)
		args         products.Product
		expectedData uint64
		err          error
	}{
		{
			name: "successful launch",
			prepare: func(f *fields) {
				mockProduct := products.Product{
					Name:      "test product",
					Price:     123,
					CreatedAt: now,
				}

				gomock.InOrder(
					f.date.EXPECT().Now().Return(now),
					f.productsRepo.EXPECT().CreateProduct(f.tx, mockProduct).Return(mockProductID, nil),
				)
			},
			args: products.Product{
				Name:  "test product",
				Price: 123,
			},
			expectedData: mockProductID,
			err:          nil,
		},
		{
			name: "unsuccessful launch",
			prepare: func(f *fields) {
				mockProduct := products.Product{
					Name:      "test product",
					Price:     123,
					CreatedAt: now,
				}

				gomock.InOrder(
					f.date.EXPECT().Now().Return(now),
					f.productsRepo.EXPECT().CreateProduct(f.tx, mockProduct).Return(uint64(0), products.ErrPermissionDenied),
				)
			},
			args: products.Product{
				Name:  "test product",
				Price: 123,
			},
			expectedData: uint64(0),
			err:          shared.ErrInternal,
		},
	}

	for _, row := range testList {
		s.Run(row.name, func() {
			ctrl := gomock.NewController(s.T())
			defer ctrl.Finish()

			f := fields{
				tx:   &sqlx.Tx{},
				date: mockshared.NewMockDateTool(ctrl),

				productsRepo:       mockproducts.NewMockProductsRepo(ctrl),
				productsStatistics: mockproducts.NewMockProductsStatistics(ctrl),
			}
			if row.prepare != nil {
				row.prepare(&f)
			}

			service := products.NewProductsService(
				s.log,
				f.date,

				f.productsRepo,
				f.productsStatistics,
			)

			data, err := service.CreateProduct(f.tx, row.args)
			s.Equal(row.err, err)
			s.Equal(row.expectedData, data)
		})
	}
}
