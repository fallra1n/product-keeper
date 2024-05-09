package services

import (
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/fallra1n/product-service/internal/domain/models"
	"github.com/fallra1n/product-service/internal/storage"
)

type Auth interface {
	CreateUser(user models.User) error
	LoginUser(user models.User) (string, error)
}

type authService struct {
	storage storage.Storage
	logger  *slog.Logger
}

func NewAuthService(storage storage.Storage, logger *slog.Logger) Auth {
	return &authService{storage, logger}
}

func (as *authService) CreateUser(user models.User) error {
	hash, err := as.hashPassword(user.Password)
	if err != nil {
		return err
	}

	return as.storage.CreateUser(models.User{
		Name:     user.Name,
		Password: hash,
	})
}

func (as *authService) LoginUser(user models.User) (string, error) {
	hashedPassword, err := as.storage.GetPasswordByName(user.Name)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password)); err != nil {
		return "", err
	}

	return as.generateToken(user.Name)
}

func (as *authService) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		as.logger.Error("Error hashing password:", err)
		return "", err
	}

	return string(hash), nil
}

func (as *authService) generateToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Minute).Unix()

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
