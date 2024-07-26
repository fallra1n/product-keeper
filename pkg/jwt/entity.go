package jwt

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
	ErrFailedGenerateToken    = errors.New("failed to generate token")
	ErrInvalidTokenClaimsType = errors.New("invalid token claims type")
	ErrFailedParseToken       = errors.New("failed to parse token")
	ErrInvalidToken           = errors.New("invalid token")
)

var (
	jwtSecret = os.Getenv("JWT_SECRET")
)

type tokenClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
