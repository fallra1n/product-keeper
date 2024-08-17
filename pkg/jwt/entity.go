package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// TokenTTL token validity period
	TokenTTL = 20 * time.Minute
)

var (
	// ErrFailedGenerateToken failed to generate token
	ErrFailedGenerateToken = errors.New("failed to generate token")
	// ErrInvalidTokenClaimsType invalid token claims type
	ErrInvalidTokenClaimsType = errors.New("invalid token claims type")
	// ErrFailedParseToken failed to parse token
	ErrFailedParseToken = errors.New("failed to parse token")
	// ErrInvalidToken invalid token
	ErrInvalidToken = errors.New("invalid token")
)

var (
	jwtSecret = os.Getenv("JWT_SECRET")
)

type tokenClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
