package handlers

import "github.com/gin-gonic/gin"

type AuthHandler interface {
	UserRegister(c *gin.Context)
	UserLogin(c *gin.Context)
}

type authHandler struct{}

func NewAuthHandler() AuthHandler {
	return &authHandler{}
}

func (h *authHandler) UserRegister(c *gin.Context) {}
func (h *authHandler) UserLogin(c *gin.Context)    {}
