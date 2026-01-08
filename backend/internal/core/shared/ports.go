package shared

import (
	"time"
)

// Crypto interface for working with cryptography
type Crypto interface {
	HashPassword(password string) (string, error)
	CompareHashAndPassword(hashedPassword, password string) error
}

// Jwt interface for working with jwt tokens
type Jwt interface {
	GenerateToken(username string) (string, error)
	ParseToken(tokenString string) (string, error)
}

// DateTool interface for wor working with time
type DateTool interface {
	Now() time.Time
}
