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

// User user info for auth
type User struct {
	Name     string `db:"name"`
	Password string `db:"password"`
}

func NewUser(name string, password string) User {
	return User{
		Name:     name,
		Password: password,
	}
}
