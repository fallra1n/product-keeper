package auth

import (
	"errors"
)

var (
	// ErrIncorrectPassword incorrect password
	ErrIncorrectPassword = errors.New("incorrect password")

	// ErrUserNotFound user not found
	ErrUserNotFound = errors.New("user not found")

	// ErrUserAlreadyExist user already exists
	ErrUserAlreadyExist = errors.New("user already exists")
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
