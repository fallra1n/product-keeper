package products

import (
	"errors"
	"time"
)

var (
	// ErrProductNotFound product not found
	ErrProductNotFound = errors.New("product not found")

	// ErrPermissionDenied user does not have access to this product
	ErrPermissionDenied = errors.New("user does not have access to this product")
)

type SortType string

const (
	LastCreate SortType = "last_create"
	Name       SortType = "name"
	Empty      SortType = ""
)

type Product struct {
	ID        uint64    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Price     uint64    `json:"price" db:"price"`
	Quantity  uint64    `json:"quantity" db:"quantity"`
	OwnerName string    `json:"owner_name" db:"owner_name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
