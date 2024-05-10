package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fallra1n/product-service/internal/domain/models"
	"github.com/fallra1n/product-service/internal/services"
)

type ProductHandler interface {
	CreateProduct(c *gin.Context)
	GetProducts(c *gin.Context)
	GetProductByID(c *gin.Context)
	ChangeProductByID(c *gin.Context)
	DeleteProductByID(c *gin.Context)
}

type ProductRequest struct {
	Name     string `json:"name" binding:"required"`
	Price    uint64 `json:"price" binding:"required"`
	Quantity uint64 `json:"quantity" binding:"required"`
}

type ProductResponse struct {
	ID       uint64 `json:"id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Price    uint64 `json:"price" binding:"required"`
	Quantity uint64 `json:"quantity" binding:"required"`
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
	var req ProductRequest

	if err := c.BindJSON(&req); err != nil {
		p.logger.Error("CreateProduct: " + err.Error())
		c.JSON(http.StatusBadRequest, DefaultResponse{"failed to decode request"})
		return
	}

	userName, ok := c.Get(UserContext)
	if !ok {
		return
	}

	id, err := p.services.CreateProduct(models.Product{
		Name:      req.Name,
		Price:     req.Price,
		Quantity:  req.Quantity,
		OwnerName: userName.(string),
	})
	if err != nil {
		p.logger.Error("CreateUser: " + err.Error())
		c.JSON(http.StatusInternalServerError, DefaultResponse{"internal server error"})
		return
	}

	p.logger.Info("CreateProduct: product has been successfully created")
	c.JSON(http.StatusCreated, map[string]any{
		"product_id": id,
	})
}

func (p *productHandler) GetProducts(c *gin.Context) {}

func (p *productHandler) GetProductByID(c *gin.Context) {}

func (p *productHandler) ChangeProductByID(c *gin.Context) {}

func (p *productHandler) DeleteProductByID(c *gin.Context) {}
