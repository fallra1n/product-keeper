package models

type Product struct {
	ID        uint64 `json:"id" db:"id"`
	Name      string `json:"name" db:"name" binding:"required"`
	Price     uint64 `json:"price" db:"price" binding:"required"`
	Quantity  uint64 `json:"quantity" db:"quantity" binding:"required"`
	OwnerName string `json:"owner_name" db:"owner_name" binding:"required"`
}
