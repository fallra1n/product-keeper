package auth

import (
	"errors"
)

var (
	ErrFailedHashingPassword = errors.New("failed to hash password")
	ErrIncorrectPassword     = errors.New("incorrect password")
	ErrUserNotFound          = errors.New("user not found")
	ErrUserAlreadyExist      = errors.New("user already exists")
)

type User struct {
	Name     string `db:"name"`
	Password string `db:"password"`
}
