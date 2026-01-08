package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Jwt struct for working with jwt
type Jwt struct{}

// NewJwt constructor for Jwt
func NewJwt() *Jwt {
	return &Jwt{}
}

// GenerateToken generate token
func (j Jwt) GenerateToken(username string) (string, error) {
	return GenerateToken(username)
}

// ParseToken parse token
func (j Jwt) ParseToken(tokenString string) (string, error) {
	return ParseToken(tokenString)
}

// GenerateToken generate token
func GenerateToken(username string) (string, error) {
	claims := &tokenClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", ErrFailedGenerateToken
	}

	return tokenString, nil
}

// ParseToken parse token
func ParseToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return "", ErrFailedParseToken
	}

	if !token.Valid {
		return "", ErrInvalidToken
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", ErrInvalidTokenClaimsType
	}

	return claims.Username, err
}
