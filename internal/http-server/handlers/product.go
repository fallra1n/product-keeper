package handlers

import (
	"github.com/fallra1n/product-service/internal/services"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type ProductHandler interface {
	CreateProduct(c *gin.Context)
	GetProducts(c *gin.Context)
	GetProductByID(c *gin.Context)
	ChangeProductByID(c *gin.Context)
	DeleteProductByID(c *gin.Context)
}

type productHandler struct {
	services services.Services
	logger   *slog.Logger
}

func NewProductHandler(services services.Services, logger *slog.Logger) ProductHandler {
	return &productHandler{
		services,
		logger,
	}
}

func (p *productHandler) CreateProduct(c *gin.Context) {

}
func (p *productHandler) GetProducts(c *gin.Context)       {}
func (p *productHandler) GetProductByID(c *gin.Context)    {}
func (p *productHandler) ChangeProductByID(c *gin.Context) {}
func (p *productHandler) DeleteProductByID(c *gin.Context) {}
