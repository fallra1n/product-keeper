package service

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Auth interface {
	CreateUser(username, password string) (uint64, error)
	GenerateToken(username, password string) (string, error)
}

type authService struct {
}

func NewAuthService() Auth {
	return &authService{}
}

func (s *authService) CreateUser(username, password string) (uint64, error) {
	// TODO add user to db
	return 0, nil
}

func (s *authService) GenerateToken(username, password string) (string, error) {
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
