package handlers

import "github.com/gin-gonic/gin"

type ProductHandler interface {
	CreateProduct(c *gin.Context)
	GetProducts(c *gin.Context)
	GetProductByID(c *gin.Context)
	ChangeProductByID(c *gin.Context)
	DeleteProductByID(c *gin.Context)
}

type productHandler struct{}

func NewProductHandler() ProductHandler {
	return &productHandler{}
}

func (h *productHandler) CreateProduct(c *gin.Context)     {}
func (h *productHandler) GetProducts(c *gin.Context)       {}
func (h *productHandler) GetProductByID(c *gin.Context)    {}
func (h *productHandler) ChangeProductByID(c *gin.Context) {}
func (h *productHandler) DeleteProductByID(c *gin.Context) {}
