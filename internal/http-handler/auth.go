package httphandler

import (
	"errors"
	"log/slog"
	"net/http"

	services "github.com/fallra1n/product-keeper/internal/core"
	"github.com/fallra1n/product-keeper/internal/domain/models"
	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	UserRegister(c *gin.Context)
	UserLogin(c *gin.Context)
	UserIdentity(c *gin.Context)
}

type AuthRequest struct {
	Name     string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type authHandler struct {
	services services.Services
	logger   *slog.Logger
}

func NewAuthHandler(services services.Services, logger *slog.Logger) AuthHandler {
	return &authHandler{
		services,
		logger,
	}
}

func (ah *authHandler) UserRegister(c *gin.Context) {
	var req AuthRequest

	if err := c.BindJSON(&req); err != nil {
		ah.logger.Error("UserRegister: " + err.Error())
		c.JSON(http.StatusBadRequest, DefaultResponse{"failed to decode request"})
		return
	}

	err := ah.services.CreateUser(models.User{
		Name:     req.Name,
		Password: req.Password,
	})

	if err != nil {
		if errors.Is(err, services.ErrFailedHashingPassword) {
			ah.logger.Error("UserRegister: " + err.Error())
			c.JSON(http.StatusInternalServerError, DefaultResponse{"cannot hash password"})
			return
		}

		if errors.Is(err, services.ErrUserAlreadyExist) {
			ah.logger.Error("UserRegister: " + err.Error())
			c.JSON(http.StatusBadRequest, DefaultResponse{"username already exists"})
			return
		}

		ah.logger.Error("UserRegister: " + err.Error())
		c.JSON(http.StatusInternalServerError, DefaultResponse{"internal server error"})
		return
	}

	ah.logger.Info("UserRegister: a user has been successfully registered")
	c.JSON(http.StatusOK, DefaultResponse{"a user has been successfully registered"})
}

func (ah *authHandler) UserLogin(c *gin.Context) {
	var req AuthRequest

	if err := c.BindJSON(&req); err != nil {
		ah.logger.Error("UserLogin: " + err.Error())
		c.JSON(http.StatusBadRequest, DefaultResponse{"failed to decode request"})
		return
	}

	token, err := ah.services.LoginUser(models.User{
		Name:     req.Name,
		Password: req.Password,
	})

	if err != nil {
		if errors.Is(err, services.ErrIncorrectPassword) {
			ah.logger.Error("UserLogin: " + err.Error())
			c.JSON(http.StatusBadRequest, DefaultResponse{"incorrect password"})
			return
		}

		if errors.Is(err, services.ErrUserNotFound) {
			ah.logger.Error("UserLogin: " + err.Error())
			c.JSON(http.StatusNotFound, DefaultResponse{"user not found"})
			return
		}

		ah.logger.Error("UserLogin: " + err.Error())
		c.JSON(http.StatusInternalServerError, DefaultResponse{"cannot hash password"})
		return
	}

	ah.logger.Info("UserLogin: a user has been successfully authorized")
	c.JSON(http.StatusOK, LoginResponse{token})
}
