package auth

import (
	"errors"
	"log/slog"

	"github.com/jmoiron/sqlx"

	"github.com/fallra1n/product-keeper/internal/core/shared"
)

// AuthService ...
type AuthService struct {
	log    *slog.Logger
	crypto shared.Crypto
	jwt    shared.Jwt

	authRepo AuthRepo
}

// NewAuthService constructor for AuthService
func NewAuthService(
	log *slog.Logger,
	crypto shared.Crypto,
	jwt shared.Jwt,

	authRepo AuthRepo,
) *AuthService {
	return &AuthService{
		log:    log,
		crypto: crypto,
		jwt:    jwt,

		authRepo: authRepo,
	}
}

// CreateUser ...
func (s *AuthService) CreateUser(tx *sqlx.Tx, user User) error {
	hash, err := s.crypto.HashPassword(user.Password)
	if err != nil {
		s.log.Error("failed to hash password", "error", err, "password", user.Password)
		return shared.ErrInternal
	}

	hashedUser := NewUser(user.Name, hash)

	if err := s.authRepo.CreateUser(tx, hashedUser); err != nil {
		s.log.Error("failed to create user", "error", err, "username", user.Name, "hashed_password", hashedUser.Password)

		if errors.Is(err, ErrUserAlreadyExist) {
			return ErrUserAlreadyExist
		}
		return err
	}

	return nil
}

// LoginUser ...
func (s *AuthService) LoginUser(tx *sqlx.Tx, user User) (string, error) {
	hashedPassword, err := s.authRepo.FindPassword(tx, user.Name)
	if err != nil {
		s.log.Error("failed to find password", "error", err, "username", user.Name)

		if errors.Is(err, ErrUserNotFound) {
			return "", ErrUserNotFound
		}

		return "", shared.ErrInternal
	}

	if err := s.crypto.CompareHashAndPassword(hashedPassword, user.Password); err != nil {
		s.log.Error("incorrect password", "password", user.Password, "hashed_password", hashedPassword)
		return "", ErrIncorrectPassword
	}

	token, err := s.jwt.GenerateToken(user.Name)
	if err != nil {
		s.log.Error("failed to generate token", "error", err, "username", user.Name)
		return "", shared.ErrInternal
	}

	return token, nil
}
