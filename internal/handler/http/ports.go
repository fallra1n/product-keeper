package httphandler

import "github.com/gin-gonic/gin"

// AuthHandler ...
type AuthHandler interface {
	UserRegister(c *gin.Context)
	UserLogin(c *gin.Context)
}

// ProductsHandler ...
type ProductsHandler interface {
	CreateProduct(c *gin.Context)
	FindProduct(c *gin.Context)
	UpdateProduct(c *gin.Context)
	DeleteProduct(c *gin.Context)
	FindProductList(c *gin.Context)
}
