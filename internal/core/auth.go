package services

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/fallra1n/product-keeper/internal/domain/models"
)

const (
	TokenTTL = 5 * time.Minute
)

var (
	ErrFailedHashingPassword  = errors.New("failed to hash password")
	ErrFailedGenerateToken    = errors.New("failed to generate token")
	ErrIncorrectPassword      = errors.New("incorrect password")
	ErrUserNotFound           = errors.New("user not found")
	ErrUserAlreadyExist       = errors.New("user exists")
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

type authService struct {
	storage storage.Storage
}

func NewAuthService(storage storage.Storage) Auth {
	return &authService{storage}
}

func (as *authService) CreateUser(user models.User) error {
	hash, err := as.hashPassword(user.Password)
	if err != nil {
		return err
	}

	hashedUser := models.User{
		Name:     user.Name,
		Password: hash,
	}

	if err := as.storage.CreateUser(hashedUser); err != nil {
		if errors.Is(err, storage.ErrUserAlreadyExist) {
			return ErrUserAlreadyExist
		}
		return err
	}

	return nil
}

func (as *authService) LoginUser(user models.User) (string, error) {
	hashedPassword, err := as.storage.GetPasswordByName(user.Name)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return "", ErrUserNotFound
		}

		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password)); err != nil {
		return "", ErrIncorrectPassword
	}

	return as.generateToken(user.Name)
}

func (as *authService) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", ErrFailedHashingPassword
	}

	return string(hash), nil
}

func (as *authService) generateToken(username string) (string, error) {
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

func (as *authService) ParseToken(tokenString string) (string, error) {
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
