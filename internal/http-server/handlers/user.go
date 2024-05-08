package handlers

import "github.com/gin-gonic/gin"

type UserHandler interface {
	UserRegister(c *gin.Context)
	UserLogin(c *gin.Context)
}

type userHandler struct{}

func NewUserHandler() UserHandler {
	return &userHandler{}
}

func (h *userHandler) UserRegister(c *gin.Context) {}
func (h *userHandler) UserLogin(c *gin.Context)    {}
