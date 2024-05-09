package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/fallra1n/product-service/internal/domain/models"
	"github.com/fallra1n/product-service/internal/storage"
)

var (
	ErrFailedHashingPassword = errors.New("failed to hash password")
	ErrFailedGenerateToken   = errors.New("failed to generate token")
	ErrIncorrectPassword     = errors.New("incorrect password")
	ErrUserNotFound          = errors.New("url not found")
	ErrUserAlreadyExist      = errors.New("url exists")
)

type Auth interface {
	CreateUser(user models.User) error
	LoginUser(user models.User) (string, error)
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
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Minute).Unix()

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", ErrFailedGenerateToken
	}

	return tokenString, nil
}
