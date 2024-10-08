package products

import (
	"errors"
	"time"
)

var (
	// ErrProductNotFound product not found
	ErrProductNotFound = errors.New("product not found")

	// ErrProductListNotFound product list not found
	ErrProductListNotFound = errors.New("product list not found")

	// ErrPermissionDenied user does not have access to this product
	ErrPermissionDenied = errors.New("user does not have access to this product")
)

// SortType FindProductList param
type SortType string

const (
	// LastCreate sort by created_at
	LastCreate SortType = "last_create"

	// Name sort by name
	Name SortType = "name"

	// Empty without sorting
	Empty SortType = ""
)

// Product info about product
type Product struct {
	ID        uint64    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Price     uint64    `json:"price" db:"price"`
	Quantity  uint64    `json:"quantity" db:"quantity"`
	OwnerName string    `json:"owner_name" db:"owner_name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// NewProduct constructor for Product
func NewProduct(
	id uint64,
	name string,
	price,
	quantity uint64,
	ownerName string,
	createdAt time.Time,
) Product {
	return Product{
		ID:        id,
		Name:      name,
		Price:     price,
		Quantity:  quantity,
		OwnerName: ownerName,
		CreatedAt: createdAt,
	}
}
