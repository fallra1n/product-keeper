package handlers

import (
	"github.com/gin-gonic/gin"

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

type authHandler struct {
	authService service.Auth
}

func NewAuthHandler(authService service.Auth) AuthHandler {
	return &authHandler{authService}
}

func (ah *authHandler) UserRegister(c *gin.Context) {
	// TODO call authservice.createuser
}

func (ah *authHandler) UserLogin(c *gin.Context) {
	// TODO call authservice.generatetoken...
}
