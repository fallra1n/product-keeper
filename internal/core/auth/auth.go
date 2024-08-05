package auth

import (
	"errors"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	"github.com/fallra1n/product-keeper/pkg/jwt"
)

type AuthService struct {
	log *slog.Logger

	authRepo AuthRepo
}

func NewAuthService(log *slog.Logger, authRepo AuthRepo) *AuthService {
	return &AuthService{
		log: log,

		authRepo: authRepo,
	}
}

func (s *AuthService) CreateUser(tx *sqlx.Tx, user User) error {
	hash, err := s.hashPassword(user.Password)
	if err != nil {
		return err
	}

	hashedUser := User{
		Name:     user.Name,
		Password: hash,
	}

	if err := s.authRepo.CreateUser(tx, hashedUser); err != nil {
		if errors.Is(err, ErrUserAlreadyExist) {
			return ErrUserAlreadyExist
		}
		return err
	}

	return nil
}

func (s *AuthService) LoginUser(tx *sqlx.Tx, user User) (string, error) {
	hashedPassword, err := s.authRepo.FindPassword(tx, user.Name)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return "", ErrUserNotFound
		}

		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password)); err != nil {
		return "", ErrIncorrectPassword
	}

	token, err := jwt.GenerateToken(user.Name)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", ErrFailedHashingPassword
	}

	return string(hash), nil
}
