package productshttphandler

import "time"

type DefaultResponse struct {
	Message string `json:"message"`
}

type ProductRequest struct {
	Name     string `json:"name" binding:"required"`
	Price    uint64 `json:"price" binding:"required"`
	Quantity uint64 `json:"quantity" binding:"required"`
}

type ProductResponse struct {
	ID        uint64    `json:"id" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	Price     uint64    `json:"price" binding:"required"`
	Quantity  uint64    `json:"quantity" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
}
