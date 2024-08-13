package auth

import (
	"errors"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	"github.com/fallra1n/product-keeper/internal/core/shared"
	"github.com/fallra1n/product-keeper/pkg/jwt"
)

type AuthService struct {
	log    *slog.Logger
	crypto shared.Crypto

	authRepo AuthRepo
}

func NewAuthService(log *slog.Logger, authRepo AuthRepo, crypto shared.Crypto) *AuthService {
	return &AuthService{
		log:    log,
		crypto: crypto,

		authRepo: authRepo,
	}
}

func (s *AuthService) CreateUser(tx *sqlx.Tx, user User) error {
	hash, err := s.crypto.HashPassword(user.Password)
	if err != nil {
		s.log.Error("failed to hash password", "error", err.Error())
		return shared.ErrInternal
	}

	hashedUser := NewUser(user.Name, hash)

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

		return "", shared.ErrInternal
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
