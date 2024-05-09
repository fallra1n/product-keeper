package models

type User struct {
	Name     string `json:"username" db:"name" binding:"required"`
	Password string `json:"password" db:"password" binding:"required"`
}
