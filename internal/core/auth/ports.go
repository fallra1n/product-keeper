package auth

import (
	"github.com/jmoiron/sqlx"
)

type Authrepo interface {
	CreateUser(tx *sqlx.Tx, user User) error
	FindPassword(tx *sqlx.Tx, name string) (string, error)
}
