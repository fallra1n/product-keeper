package http

import "github.com/gin-gonic/gin"

type AuthHandler interface {
	UserRegister(c *gin.Context)
	UserLogin(c *gin.Context)
}

type ProductsHandler interface {
	CreateProduct(c *gin.Context)
	GetProductByID(c *gin.Context)
	UpdateProductByID(c *gin.Context)
	DeleteProductByID(c *gin.Context)
	GetProducts(c *gin.Context)
}
