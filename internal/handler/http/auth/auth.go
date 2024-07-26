package authhttphandler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fallra1n/product-keeper/internal/core/auth"
)

type AuthHandler struct {
	log *slog.Logger

	authService *auth.AuthService
}

func NewAuthHandler(authService *auth.AuthService, log *slog.Logger) *AuthHandler {
	return &AuthHandler{
		log: log,

		authService: authService,
	}
}

func (h *AuthHandler) UserRegister(c *gin.Context) {
	var req AuthRequest

	if err := c.BindJSON(&req); err != nil {
		h.log.Error("UserRegister: " + err.Error())
		c.JSON(http.StatusBadRequest, DefaultResponse{"failed to decode request"})
		return
	}

	err := h.authService.CreateUser(auth.User{
		Name:     req.Name,
		Password: req.Password,
	})

	if err != nil {
		if errors.Is(err, auth.ErrFailedHashingPassword) {
			h.log.Error("UserRegister: " + err.Error())
			c.JSON(http.StatusInternalServerError, DefaultResponse{"cannot hash password"})
			return
		}

		if errors.Is(err, auth.ErrUserAlreadyExist) {
			h.log.Error("UserRegister: " + err.Error())
			c.JSON(http.StatusBadRequest, DefaultResponse{"username already exists"})
			return
		}

		h.log.Error("UserRegister: " + err.Error())
		c.JSON(http.StatusInternalServerError, DefaultResponse{"internal server error"})
		return
	}

	h.log.Info("UserRegister: a user has been successfully registered")
	c.JSON(http.StatusOK, DefaultResponse{"a user has been successfully registered"})
}

func (h *AuthHandler) UserLogin(c *gin.Context) {
	var req AuthRequest

	if err := c.BindJSON(&req); err != nil {
		h.log.Error("UserLogin: " + err.Error())
		c.JSON(http.StatusBadRequest, DefaultResponse{"failed to decode request"})
		return
	}

	token, err := h.authService.LoginUser(auth.User{
		Name:     req.Name,
		Password: req.Password,
	})

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
		c.JSON(http.StatusInternalServerError, DefaultResponse{"cannot hash password"})
		return
	}

	h.log.Info("UserLogin: a user has been successfully authorized")
	c.JSON(http.StatusOK, LoginResponse{token})
}
