package authhttphandler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"github.com/fallra1n/product-keeper/internal/core/auth"
)

// AuthHandler ...
type AuthHandler struct {
	log *slog.Logger
	db  *sqlx.DB

	authService *auth.AuthService
}

// NewAuthHandler constructor for AuthHandler
func NewAuthHandler(log *slog.Logger, db *sqlx.DB, authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		log: log,
		db:  db,

		authService: authService,
	}
}

// UserRegister ...
func (h *AuthHandler) UserRegister(c *gin.Context) {
	var req AuthRequest

	if err := c.BindJSON(&req); err != nil {
		h.log.Error("UserRegister: " + err.Error())
		c.JSON(http.StatusBadRequest, DefaultResponse{"failed to decode request"})
		return
	}

	tx, err := h.db.Beginx()
	if err != nil {
		h.log.Error(fmt.Sprintf("cannot start transaction: %s", err))
		c.JSON(http.StatusInternalServerError, DefaultResponse{"internal error"})
		return
	}
	defer tx.Rollback()

	err = h.authService.CreateUser(tx, auth.NewUser(
		req.Name, req.Password,
	))

	if err != nil {
		if errors.Is(err, auth.ErrUserAlreadyExist) {
			h.log.Error("UserRegister: " + err.Error())
			c.JSON(http.StatusBadRequest, DefaultResponse{"username already exists"})
			return
		}

		h.log.Error("UserRegister: " + err.Error())
		c.JSON(http.StatusInternalServerError, DefaultResponse{"internal error"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, DefaultResponse{"internal error"})
		return
	}

	h.log.Info("UserRegister: a user has been successfully registered")
	c.JSON(http.StatusOK, DefaultResponse{"a user has been successfully registered"})
}

// UserLogin ...
func (h *AuthHandler) UserLogin(c *gin.Context) {
	var req AuthRequest

	if err := c.BindJSON(&req); err != nil {
		h.log.Error("UserLogin: " + err.Error())
		c.JSON(http.StatusBadRequest, DefaultResponse{"failed to decode request"})
		return
	}

	tx, err := h.db.Beginx()
	if err != nil {
		h.log.Error(fmt.Sprintf("cannot start transaction: %s", err))
		c.JSON(http.StatusInternalServerError, DefaultResponse{"internal error"})
		return
	}
	defer tx.Rollback()

	token, err := h.authService.LoginUser(tx, auth.NewUser(
		req.Name,
		req.Password,
	))

	if err != nil {
		if errors.Is(err, auth.ErrIncorrectPassword) {
			h.log.Error("UserLogin: " + err.Error())
			c.JSON(http.StatusBadRequest, DefaultResponse{"incorrect password"})
			return
		}

		if errors.Is(err, auth.ErrUserNotFound) {
			h.log.Error("UserLogin: " + err.Error())
			c.JSON(http.StatusNotFound, DefaultResponse{"user not found"})
			return
		}

		h.log.Error("UserLogin: " + err.Error())
		c.JSON(http.StatusInternalServerError, DefaultResponse{"internal error"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, DefaultResponse{"internal error"})
		return
	}

	h.log.Info("UserLogin: a user has been successfully authorized")
	c.JSON(http.StatusOK, LoginResponse{token})
}
