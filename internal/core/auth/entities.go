package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenTTL = 5 * time.Minute
)

var (
	ErrFailedHashingPassword  = errors.New("failed to hash password")
	ErrFailedGenerateToken    = errors.New("failed to generate token")
	ErrIncorrectPassword      = errors.New("incorrect password")
	ErrUserNotFound           = errors.New("user not found")
	ErrInvalidTokenClaimsType = errors.New("invalid token claims type")
	ErrUserAlreadyExist       = errors.New("user already exists")
	ErrFailedParseToken       = errors.New("failed to parse token")
	ErrInvalidToken           = errors.New("invalid token")
)

var (
	jwtSecret = os.Getenv("JWT_SECRET")
)

type User struct {
	Name     string `json:"username" db:"name" binding:"required"`
	Password string `json:"password" db:"password" binding:"required"`
}

type tokenClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
