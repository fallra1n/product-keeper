package httphandler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/fallra1n/product-keeper/internal/domain/models"
	"github.com/fallra1n/product-keeper/internal/services"
	"github.com/gin-gonic/gin"
)

type ProductHandler interface {
	CreateProduct(c *gin.Context)
	GetProductByID(c *gin.Context)
	UpdateProductByID(c *gin.Context)
	DeleteProductByID(c *gin.Context)
	GetProducts(c *gin.Context)
}

type ProductRequest struct {
	Name     string `json:"name" binding:"required"`
	Price    uint64 `json:"price" binding:"required"`
	Quantity uint64 `json:"quantity" binding:"required"`
}

type ProductResponse struct {
	ID        uint64    `json:"id" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	Price     uint64    `json:"price" binding:"required"`
	Quantity  uint64    `json:"quantity" binding:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
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
	userName, ok := c.Get(UserContext)
	if !ok {
		return
	}

	var req ProductRequest
	if err := c.BindJSON(&req); err != nil {
		p.logger.Error("CreateProduct: " + err.Error())
		c.JSON(http.StatusBadRequest, DefaultResponse{"failed to decode request"})
		return
	}

	id, err := p.services.CreateProduct(models.Product{
		Name:      req.Name,
		Price:     req.Price,
		Quantity:  req.Quantity,
		OwnerName: userName.(string),
	})
	if err != nil {
		p.logger.Error("CreateProduct: " + err.Error())
		c.JSON(http.StatusInternalServerError, DefaultResponse{"internal server error"})
		return
	}

	p.logger.Info("CreateProduct: product has been successfully created")
	c.JSON(http.StatusCreated, map[string]any{
		"product_id": id,
	})
}

func (p *productHandler) GetProductByID(c *gin.Context) {
	userName, ok := c.Get(UserContext)
	if !ok {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		p.logger.Error("GetProductByID: " + err.Error())
		c.JSON(http.StatusBadRequest, DefaultResponse{"invalid id param"})
		return
	}

	product, err := p.services.GetProductByID(id, userName.(string))
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			p.logger.Error("GetProductByID: " + err.Error())
			c.JSON(http.StatusNotFound, DefaultResponse{"product with such id does not exist"})
			return
		}

		if errors.Is(err, services.ErrPermissionDenied) {
			p.logger.Error("GetProductByID: " + err.Error())
			c.JSON(http.StatusForbidden, DefaultResponse{"permission denied"})
			return
		}

		p.logger.Error("GetProductByID: " + err.Error())
		c.JSON(http.StatusInternalServerError, DefaultResponse{"internal server error"})
		return
	}

	p.logger.Info("GetProductByID: product data has been successfully received")
	c.JSON(http.StatusOK, ProductResponse{
		ID:       product.ID,
		Name:     product.Name,
		Price:    product.Price,
		Quantity: product.Quantity,
	})
}

func (p *productHandler) UpdateProductByID(c *gin.Context) {
	userName, ok := c.Get(UserContext)
	if !ok {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		p.logger.Error("UpdateProductByID: " + err.Error())
		c.JSON(http.StatusBadRequest, DefaultResponse{"invalid id param"})
		return
	}

	var req ProductRequest
	if err := c.BindJSON(&req); err != nil {
		p.logger.Error("UpdateProductByID: " + err.Error())
		c.JSON(http.StatusBadRequest, DefaultResponse{"failed to decode request"})
		return
	}

	updated, err := p.services.UpdateProductByID(models.Product{
		ID:        id,
		Name:      req.Name,
		Price:     req.Price,
		Quantity:  req.Quantity,
		OwnerName: userName.(string),
	})
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			p.logger.Error("UpdateProductByID: " + err.Error())
			c.JSON(http.StatusNotFound, DefaultResponse{"product with such id does not exist"})
			return
		}

		if errors.Is(err, services.ErrPermissionDenied) {
			p.logger.Error("UpdateProductByID: " + err.Error())
			c.JSON(http.StatusForbidden, DefaultResponse{"permission denied"})
			return
		}

		p.logger.Error("UpdateProductByID: " + err.Error())
		c.JSON(http.StatusInternalServerError, DefaultResponse{"internal server error"})
		return
	}

	p.logger.Info("UpdateProductByID: product data has been successfully updated")
	c.JSON(http.StatusOK, ProductResponse{
		ID:       updated.ID,
		Name:     updated.Name,
		Price:    updated.Price,
		Quantity: updated.Quantity,
	})
}

func (p *productHandler) DeleteProductByID(c *gin.Context) {
	userName, ok := c.Get(UserContext)
	if !ok {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		p.logger.Error("DeleteProductByID: " + err.Error())
		c.JSON(http.StatusBadRequest, DefaultResponse{"invalid id param"})
		return
	}

	if err := p.services.DeleteProductByID(id, userName.(string)); err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			p.logger.Error("DeleteProductByID: " + err.Error())
			c.JSON(http.StatusNotFound, DefaultResponse{"product with such id does not exist"})
			return
		}

		if errors.Is(err, services.ErrPermissionDenied) {
			p.logger.Error("DeleteProductByID: " + err.Error())
			c.JSON(http.StatusForbidden, DefaultResponse{"permission denied"})
			return
		}

		p.logger.Error("DeleteProductByID: " + err.Error())
		c.JSON(http.StatusInternalServerError, DefaultResponse{"internal server error"})
		return
	}

	p.logger.Info("DeleteProductByID: product has been successfully deleted")
	c.JSON(http.StatusOK, DefaultResponse{"product has been successfully deleted"})
}

func (p *productHandler) GetProducts(c *gin.Context) {
	userName, ok := c.Get(UserContext)
	if !ok {
		return
	}

	productName := c.Query("name")
	sortByString := c.Query("sort_by")

	var sortBy models.SortType
	switch sortByString {
	case "last_create":
		sortBy = models.LastCreate
	case "name":
		sortBy = models.Name
	case "":
		sortBy = models.Empty
	default:
		p.logger.Error("GetProducts: bad sort_by param")
		c.JSON(http.StatusBadRequest, DefaultResponse{"invalid sort_by param"})
		return
	}

	products, err := p.services.GetProducts(userName.(string), productName, sortBy)
	if err != nil {
		p.logger.Error("GetProducts: " + err.Error())
		c.JSON(http.StatusInternalServerError, DefaultResponse{"internal server error"})
		return
	}

	var productsResponse []ProductResponse
	for _, product := range products {
		productResponse := ProductResponse{
			ID:        product.ID,
			Name:      product.Name,
			Price:     product.Price,
			Quantity:  product.Quantity,
			CreatedAt: product.CreatedAt,
		}
		productsResponse = append(productsResponse, productResponse)
	}

	p.logger.Info("GetProducts: products has been successfully received")
	c.JSON(http.StatusOK, productsResponse)
}
