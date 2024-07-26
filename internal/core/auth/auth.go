package auth

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db  *sqlx.DB
	log *slog.Logger

	authRepo Authrepo
}

func NewAuthService(authRepo Authrepo, db *sqlx.DB, log *slog.Logger) *AuthService {
	return &AuthService{
		db:  db,
		log: log,

		authRepo: authRepo,
	}
}

func (s *AuthService) CreateUser(user User) error {
	tx, err := s.db.Beginx()
	if err != nil {
		s.log.Error(fmt.Sprintf("cannot start transaction: %s", err))
		return err
	}
	defer tx.Rollback()

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

	if err := tx.Commit(); err != nil {
		s.log.Error(fmt.Sprintf("cannot commit transaction: %s", err))
		return err
	}

	return nil
}

func (s *AuthService) LoginUser(user User) (string, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		s.log.Error(fmt.Sprintf("cannot start transaction: %s", err))
		return "", err
	}
	defer tx.Rollback()

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

	token, err := s.generateToken(user.Name)
	if err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		s.log.Error(fmt.Sprintf("cannot commit transaction: %s", err))
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

func (s *AuthService) generateToken(username string) (string, error) {
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

func (s *AuthService) ParseToken(tokenString string) (string, error) {
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
