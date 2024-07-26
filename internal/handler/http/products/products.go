package productshttphandler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/fallra1n/product-keeper/internal/core/products"
	"github.com/fallra1n/product-keeper/internal/handler/http/middleware"
)

type ProductHandler struct {
	log *slog.Logger

	productsService *products.ProductsService
}

func NewProductHandler(productsService *products.ProductsService, log *slog.Logger) *ProductHandler {
	return &ProductHandler{
		log:             log,
		productsService: productsService,
	}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	username, ok := c.Get(middleware.UserContext)
	if !ok {
		return
	}

	var req ProductRequest
	if err := c.BindJSON(&req); err != nil {
		h.log.Error("CreateProduct: " + err.Error())
		c.JSON(http.StatusBadRequest, DefaultResponse{"failed to decode request"})
		return
	}

	id, err := h.productsService.CreateProduct(products.Product{
		Name:      req.Name,
		Price:     req.Price,
		Quantity:  req.Quantity,
		OwnerName: username.(string),
	})
	if err != nil {
		h.log.Error("CreateProduct: " + err.Error())
		c.JSON(http.StatusInternalServerError, DefaultResponse{"internal server error"})
		return
	}

	h.log.Info("CreateProduct: product has been successfully created")
	c.JSON(http.StatusCreated, map[string]any{
		"product_id": id,
	})
}

func (h *ProductHandler) FindProduct(c *gin.Context) {
	username, ok := c.Get(middleware.UserContext)
	if !ok {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.log.Error("GetProductByID: " + err.Error())
		c.JSON(http.StatusBadRequest, DefaultResponse{"invalid id param"})
		return
	}

	product, err := h.productsService.FindProduct(id, username.(string))
	if err != nil {
		if errors.Is(err, products.ErrProductNotFound) {
			h.log.Error("GetProductByID: " + err.Error())
			c.JSON(http.StatusNotFound, DefaultResponse{"product with such id does not exist"})
			return
		}

		if errors.Is(err, products.ErrPermissionDenied) {
			h.log.Error("GetProductByID: " + err.Error())
			c.JSON(http.StatusForbidden, DefaultResponse{"permission denied"})
			return
		}

		h.log.Error("GetProductByID: " + err.Error())
		c.JSON(http.StatusInternalServerError, DefaultResponse{"internal server error"})
		return
	}

	h.log.Info("GetProductByID: product data has been successfully received")
	c.JSON(http.StatusOK, ProductResponse{
		ID:       product.ID,
		Name:     product.Name,
		Price:    product.Price,
		Quantity: product.Quantity,
	})
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	username, ok := c.Get(middleware.UserContext)
	if !ok {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.log.Error("UpdateProductByID: " + err.Error())
		c.JSON(http.StatusBadRequest, DefaultResponse{"invalid id param"})
		return
	}

	var req ProductRequest
	if err := c.BindJSON(&req); err != nil {
		h.log.Error("UpdateProductByID: " + err.Error())
		c.JSON(http.StatusBadRequest, DefaultResponse{"failed to decode request"})
		return
	}

	updated, err := h.productsService.UpdateProduct(products.Product{
		ID:        id,
		Name:      req.Name,
		Price:     req.Price,
		Quantity:  req.Quantity,
		OwnerName: username.(string),
	})
	if err != nil {
		if errors.Is(err, products.ErrProductNotFound) {
			h.log.Error("UpdateProductByID: " + err.Error())
			c.JSON(http.StatusNotFound, DefaultResponse{"product with such id does not exist"})
			return
		}

		if errors.Is(err, products.ErrPermissionDenied) {
			h.log.Error("UpdateProductByID: " + err.Error())
			c.JSON(http.StatusForbidden, DefaultResponse{"permission denied"})
			return
		}

		h.log.Error("UpdateProductByID: " + err.Error())
		c.JSON(http.StatusInternalServerError, DefaultResponse{"internal server error"})
		return
	}

	h.log.Info("UpdateProductByID: product data has been successfully updated")
	c.JSON(http.StatusOK, ProductResponse{
		ID:       updated.ID,
		Name:     updated.Name,
		Price:    updated.Price,
		Quantity: updated.Quantity,
	})
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	username, ok := c.Get(middleware.UserContext)
	if !ok {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.log.Error("DeleteProductByID: " + err.Error())
		c.JSON(http.StatusBadRequest, DefaultResponse{"invalid id param"})
		return
	}

	if err := h.productsService.DeleteProduct(id, username.(string)); err != nil {
		if errors.Is(err, products.ErrProductNotFound) {
			h.log.Error("DeleteProductByID: " + err.Error())
			c.JSON(http.StatusNotFound, DefaultResponse{"product with such id does not exist"})
			return
		}

		if errors.Is(err, products.ErrPermissionDenied) {
			h.log.Error("DeleteProductByID: " + err.Error())
			c.JSON(http.StatusForbidden, DefaultResponse{"permission denied"})
			return
		}

		h.log.Error("DeleteProductByID: " + err.Error())
		c.JSON(http.StatusInternalServerError, DefaultResponse{"internal server error"})
		return
	}

	h.log.Info("DeleteProductByID: product has been successfully deleted")
	c.JSON(http.StatusOK, DefaultResponse{"product has been successfully deleted"})
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	username, ok := c.Get(middleware.UserContext)
	if !ok {
		return
	}

	productName := c.Query("name")
	sortByString := c.Query("sort_by")

	var sortBy products.SortType
	switch sortByString {
	case "last_create":
		sortBy = products.LastCreate
	case "name":
		sortBy = products.Name
	case "":
		sortBy = products.Empty
	default:
		h.log.Error("GetProducts: bad sort_by param")
		c.JSON(http.StatusBadRequest, DefaultResponse{"invalid sort_by param"})
		return
	}

	products, err := h.productsService.FindProductList(username.(string), productName, sortBy)
	if err != nil {
		h.log.Error("GetProducts: " + err.Error())
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

	h.log.Info("GetProducts: products has been successfully received")
	c.JSON(http.StatusOK, productsResponse)
}
