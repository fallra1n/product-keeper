package services

import (
	"crypto/sha1"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/fallra1n/product-service/internal/domain/models"
	"github.com/fallra1n/product-service/internal/storage/postgres"
)

type Auth interface {
	CreateUser(user models.User) error
}

type authService struct {
	storage *postgres.Storage
	logger  *slog.Logger
}

func NewAuthService(storage *postgres.Storage, logger *slog.Logger) Auth {
	return &authService{storage, logger}
}

func (as *authService) CreateUser(user models.User) error {
	return as.storage.CreateUser(models.User{
		Name:     user.Name,
		Password: as.hashPassword(user.Password),
	})
}

func (as *authService) hashPassword(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	salt := as.fetchSalt()
	as.logger.Warn("salt is empty")

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (as *authService) fetchSalt() string {
	return os.Getenv("PASSWORD_SALT")
}

func (as *authService) generateToken(username, password string) (string, error) {
	// TODO verify that the user exists in db

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
