package products

import (
	"errors"
	"time"
)

var (
	ErrProductNotFound  = errors.New("product not found")
	ErrPermissionDenied = errors.New("user does not have access to this product")
)

type SortType string

const (
	LastCreate SortType = "last_create"
	Name       SortType = "name"
	Empty      SortType = ""
)

type Product struct {
	ID        uint64    `db:"id"`
	Name      string    `db:"name"`
	Price     uint64    `db:"price"`
	Quantity  uint64    `db:"quantity"`
	OwnerName string    `db:"owner_name"`
	CreatedAt time.Time `db:"created_at"`
}
