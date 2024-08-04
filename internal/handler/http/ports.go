package httphandler

import "github.com/gin-gonic/gin"

type AuthHandler interface {
	UserRegister(c *gin.Context)
	UserLogin(c *gin.Context)
}

type ProductsHandler interface {
	CreateProduct(c *gin.Context)
	FindProduct(c *gin.Context)
	UpdateProduct(c *gin.Context)
	DeleteProduct(c *gin.Context)
	FindProductList(c *gin.Context)
}
