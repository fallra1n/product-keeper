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

		mockProductID = uint64(567)
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

func (s *RunProductsSuite) TestFindProduct() {
	type fields struct {
		tx   *sqlx.Tx
		date *mockshared.MockDateTool

		productsRepo       *mockproducts.MockProductsRepo
		productsStatistics *mockproducts.MockProductsStatistics
	}

	type args struct {
		id       uint64
		username string
	}

	var (
		mockProductID = uint64(123)
		mockUsername  = "test username"
	)

	testList := []struct {
		name         string
		prepare      func(f *fields)
		args         args
		expectedData products.Product
		err          error
	}{
		{
			name: "successful launch",
			prepare: func(f *fields) {
				mockProduct := products.Product{
					ID:        mockProductID,
					OwnerName: mockUsername,
				}

				gomock.InOrder(
					f.productsRepo.EXPECT().FindProduct(f.tx, mockProductID).Return(mockProduct, nil),
					f.productsStatistics.EXPECT().Send(mockProduct).Return(nil),
				)
			},
			args: args{
				id:       mockProductID,
				username: mockUsername,
			},
			expectedData: products.Product{
				ID:        mockProductID,
				OwnerName: mockUsername,
			},
			err: nil,
		},
		{
			name: "product not found",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.productsRepo.EXPECT().FindProduct(f.tx, mockProductID).Return(products.Product{}, shared.ErrNoData),
				)
			},
			args: args{
				id:       mockProductID,
				username: mockUsername,
			},
			expectedData: products.Product{},
			err:          products.ErrProductNotFound,
		},
		{
			name: "permission denied",
			prepare: func(f *fields) {
				mockProduct := products.Product{
					ID:        mockProductID,
					OwnerName: "other username",
				}

				gomock.InOrder(
					f.productsRepo.EXPECT().FindProduct(f.tx, mockProductID).Return(mockProduct, nil),
				)
			},
			args: args{
				id:       mockProductID,
				username: mockUsername,
			},
			expectedData: products.Product{},
			err:          products.ErrPermissionDenied,
		},
		{
			name: "failed to send statistics",
			prepare: func(f *fields) {
				mockProduct := products.Product{
					ID:        mockProductID,
					OwnerName: mockUsername,
				}

				gomock.InOrder(
					f.productsRepo.EXPECT().FindProduct(f.tx, mockProductID).Return(mockProduct, nil),
					f.productsStatistics.EXPECT().Send(mockProduct).Return(shared.ErrNoData),
				)
			},
			args: args{
				id:       mockProductID,
				username: mockUsername,
			},
			expectedData: products.Product{},
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

			data, err := service.FindProduct(f.tx, row.args.id, row.args.username)
			s.Equal(row.err, err)
			s.Equal(row.expectedData, data)
		})
	}
}

func (s *RunProductsSuite) TestUpdateProduct() {
	type fields struct {
		tx   *sqlx.Tx
		date *mockshared.MockDateTool

		productsRepo       *mockproducts.MockProductsRepo
		productsStatistics *mockproducts.MockProductsStatistics
	}

	var (
		mockProductID = uint64(123)
		mockUsername  = "test new username"
	)

	testList := []struct {
		name         string
		prepare      func(f *fields)
		args         products.Product
		expectedData products.Product
		err          error
	}{
		{
			name: "successful launch",
			prepare: func(f *fields) {
				mockProduct := products.Product{
					ID:        mockProductID,
					OwnerName: mockUsername,
					Name:      "test product",
					Price:     123,
				}

				mockNewProduct := products.Product{
					ID:        mockProductID,
					OwnerName: mockUsername,
					Name:      "new test product",
					Price:     1234,
				}

				gomock.InOrder(
					f.productsRepo.EXPECT().FindProduct(f.tx, mockProductID).Return(mockProduct, nil),
					f.productsRepo.EXPECT().UpdateProduct(f.tx, mockNewProduct).Return(mockNewProduct, nil),
				)
			},
			args: products.Product{
				ID:        mockProductID,
				OwnerName: mockUsername,
				Name:      "new test product",
				Price:     1234,
			},
			expectedData: products.Product{
				ID:        mockProductID,
				OwnerName: mockUsername,
				Name:      "new test product",
				Price:     1234,
			},
			err: nil,
		},
		{
			name: "product not found",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.productsRepo.EXPECT().FindProduct(f.tx, mockProductID).Return(products.Product{}, shared.ErrNoData),
				)
			},
			args: products.Product{
				ID:        mockProductID,
				OwnerName: mockUsername,
			},
			expectedData: products.Product{},
			err:          products.ErrProductNotFound,
		},
		{
			name: "permission denied",
			prepare: func(f *fields) {
				mockProduct := products.Product{
					ID:        mockProductID,
					OwnerName: "other username",
				}

				gomock.InOrder(
					f.productsRepo.EXPECT().FindProduct(f.tx, mockProductID).Return(mockProduct, nil),
				)
			},
			args: products.Product{
				ID:        mockProductID,
				OwnerName: mockUsername,
			},
			expectedData: products.Product{},
			err:          products.ErrPermissionDenied,
		},
		{
			name: "internal error(update product)",
			prepare: func(f *fields) {
				mockProduct := products.Product{
					ID:        mockProductID,
					OwnerName: mockUsername,
					Name:      "test product",
					Price:     123,
				}

				mockNewProduct := products.Product{
					ID:        mockProductID,
					OwnerName: mockUsername,
					Name:      "new test product",
					Price:     1234,
				}

				gomock.InOrder(
					f.productsRepo.EXPECT().FindProduct(f.tx, mockProductID).Return(mockProduct, nil),
					f.productsRepo.EXPECT().UpdateProduct(f.tx, mockNewProduct).Return(products.Product{}, shared.ErrNoData),
				)
			},
			args: products.Product{
				ID:        mockProductID,
				OwnerName: mockUsername,
				Name:      "new test product",
				Price:     1234,
			},
			expectedData: products.Product{},
			err:          shared.ErrInternal,
		},
		{
			name: "internal error(find product)",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.productsRepo.EXPECT().FindProduct(f.tx, mockProductID).Return(products.Product{}, products.ErrPermissionDenied),
				)
			},
			args: products.Product{
				ID:        mockProductID,
				OwnerName: mockUsername,
			},
			expectedData: products.Product{},
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

			data, err := service.UpdateProduct(f.tx, row.args)
			s.Equal(row.err, err)
			s.Equal(row.expectedData, data)
		})
	}
}

func (s *RunProductsSuite) TestDeleteProduct() {
	type fields struct {
		tx   *sqlx.Tx
		date *mockshared.MockDateTool

		productsRepo       *mockproducts.MockProductsRepo
		productsStatistics *mockproducts.MockProductsStatistics
	}

	type args struct {
		id       uint64
		username string
	}

	var (
		mockProductID = uint64(123)
		mockUsername  = "test username"
	)

	testList := []struct {
		name    string
		prepare func(f *fields)
		args    args
		err     error
	}{
		{
			name: "successful launch",
			prepare: func(f *fields) {
				mockProduct := products.Product{
					ID:        mockProductID,
					OwnerName: mockUsername,
					Quantity:  123,
				}

				gomock.InOrder(
					f.productsRepo.EXPECT().FindProduct(f.tx, mockProductID).Return(mockProduct, nil),
					f.productsRepo.EXPECT().DeleteProduct(f.tx, mockProductID).Return(nil),
				)
			},
			args: args{
				id:       mockProductID,
				username: mockUsername,
			},
			err: nil,
		},
		{
			name: "product not found",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.productsRepo.EXPECT().FindProduct(f.tx, mockProductID).Return(products.Product{}, shared.ErrNoData),
				)
			},
			args: args{
				id:       mockProductID,
				username: mockUsername,
			},
			err: products.ErrProductNotFound,
		},
		{
			name: "permission denied",
			prepare: func(f *fields) {
				mockProduct := products.Product{
					ID:        mockProductID,
					OwnerName: "other username",
				}

				gomock.InOrder(
					f.productsRepo.EXPECT().FindProduct(f.tx, mockProductID).Return(mockProduct, nil),
				)
			},
			args: args{
				id:       mockProductID,
				username: mockUsername,
			},
			err: products.ErrPermissionDenied,
		},
		{
			name: "internal error(delete product)",
			prepare: func(f *fields) {
				mockProduct := products.Product{
					ID:        mockProductID,
					OwnerName: mockUsername,
					Quantity:  123,
				}

				gomock.InOrder(
					f.productsRepo.EXPECT().FindProduct(f.tx, mockProductID).Return(mockProduct, nil),
					f.productsRepo.EXPECT().DeleteProduct(f.tx, mockProductID).Return(products.ErrProductNotFound),
				)
			},
			args: args{
				id:       mockProductID,
				username: mockUsername,
			},
			err: shared.ErrInternal,
		},
		{
			name: "internal error(find product)",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.productsRepo.EXPECT().FindProduct(f.tx, mockProductID).Return(products.Product{}, products.ErrPermissionDenied),
				)
			},
			args: args{
				id:       mockProductID,
				username: mockUsername,
			},
			err: shared.ErrInternal,
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

			err := service.DeleteProduct(f.tx, row.args.id, row.args.username)
			s.Equal(row.err, err)
		})
	}
}

func (s *RunProductsSuite) TestFindProductList() {
	type fields struct {
		tx   *sqlx.Tx
		date *mockshared.MockDateTool

		productsRepo       *mockproducts.MockProductsRepo
		productsStatistics *mockproducts.MockProductsStatistics
	}

	type args struct {
		username    string
		productName string
		sortBy      products.SortType
	}

	var (
		mockUsername  = "test username"
		mockProductID = uint64(123)
	)

	testList := []struct {
		name         string
		prepare      func(f *fields)
		args         args
		expectedData []products.Product
		err          error
	}{
		{
			name: "successful launch",
			prepare: func(f *fields) {
				mockProduct := products.Product{
					ID:        mockProductID,
					OwnerName: mockUsername,
				}

				gomock.InOrder(
					f.productsRepo.EXPECT().FindProductList(f.tx, mockUsername, "", products.LastCreate).Return([]products.Product{mockProduct}, nil),
				)
			},
			args: args{
				username:    mockUsername,
				productName: "",
				sortBy:      products.LastCreate,
			},
			expectedData: []products.Product{
				{
					ID:        mockProductID,
					OwnerName: mockUsername,
				},
			},
			err: nil,
		},
		{
			name: "product list not found",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.productsRepo.EXPECT().FindProductList(f.tx, mockUsername, "123", products.Empty).Return(nil, shared.ErrNoData),
				)
			},
			args: args{
				username:    mockUsername,
				productName: "123",
				sortBy:      products.Empty,
			},
			expectedData: nil,
			err:          products.ErrProductListNotFound,
		},
		{
			name: "internal error",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.productsRepo.EXPECT().FindProductList(f.tx, mockUsername, "", products.Empty).Return(nil, products.ErrProductListNotFound),
				)
			},
			args: args{
				username:    mockUsername,
				productName: "",
				sortBy:      products.Empty,
			},
			expectedData: nil,
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

			data, err := service.FindProductList(f.tx, row.args.username, row.args.productName, row.args.sortBy)
			s.Equal(row.err, err)
			s.Equal(row.expectedData, data)
		})
	}
}
