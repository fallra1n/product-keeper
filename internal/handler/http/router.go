package http

import (
	"log/slog"

	"github.com/fallra1n/product-keeper/internal/handler/http/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(auth AuthHandler, productHandlers ProductsHandler, logger *slog.Logger) *gin.Engine {
	router := gin.Default()

	// TODO using custom logger

	router.POST("/user/register", auth.UserRegister)
	router.POST("/user/login", auth.UserLogin)

	products := router.Group("/products", middleware.UserIdentity())
	{
		products.GET("", productHandlers.GetProducts)
	}

	product := router.Group("/product", middleware.UserIdentity())
	{
		product.POST("/add", productHandlers.CreateProduct)
		product.GET("/:id", productHandlers.GetProductByID)
		product.PUT("/:id", productHandlers.UpdateProductByID)
		product.DELETE("/:id", productHandlers.DeleteProductByID)
	}

	return router
}
