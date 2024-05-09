package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fallra1n/product-service/internal/domain/models"
	"github.com/fallra1n/product-service/internal/services"
)

type AuthHandler interface {
	UserRegister(c *gin.Context)
	UserLogin(c *gin.Context)
}

type Request struct {
	Name     string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type DefaultResponse struct {
	Message string `json:"message"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type authHandler struct {
	authService services.Auth
	logger      *slog.Logger
}

func NewAuthHandler(authService services.Auth, logger *slog.Logger) AuthHandler {
	return &authHandler{authService, logger}
}

func (ah *authHandler) UserRegister(c *gin.Context) {
	var req Request

	if err := c.BindJSON(&req); err != nil {
		ah.logger.Error("UserRegister: " + err.Error())
		c.JSON(http.StatusBadRequest, DefaultResponse{err.Error()})
		return
	}

	err := ah.authService.CreateUser(models.User{
		Name:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		ah.logger.Error("UserRegister: " + err.Error())
		c.JSON(http.StatusInternalServerError, DefaultResponse{err.Error()})
		return
	}

	ah.logger.Info("UserRegister: a user has been successfully registered with username: %s", req.Name)
	c.JSON(http.StatusOK, DefaultResponse{"a user has been successfully registered"})
}

func (ah *authHandler) UserLogin(c *gin.Context) {
	var req Request

	if err := c.BindJSON(&req); err != nil {
		ah.logger.Error("UserLogin: " + err.Error())
		c.JSON(http.StatusBadRequest, DefaultResponse{err.Error()})
		return
	}

	token, err := ah.authService.LoginUser(models.User{
		Name:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		ah.logger.Error("UserLogin: " + err.Error())
		c.JSON(http.StatusInternalServerError, DefaultResponse{err.Error()})
		return
	}

	ah.logger.Info("UserLogin: a user has been successfully authorized: %s", req.Name)
	c.JSON(http.StatusOK, LoginResponse{token})
}
