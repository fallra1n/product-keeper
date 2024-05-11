package http

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/fallra1n/product-keeper/internal/http-server/handlers"
)

func SetupRouter(auth handlers.AuthHandler, productHandlers handlers.ProductHandler, logger *slog.Logger) *gin.Engine {
	router := gin.Default()

	// TODO using custom logger

	router.POST("/user/register", auth.UserRegister)
	router.POST("/user/login", auth.UserLogin)
	router.GET("/products", productHandlers.GetProducts)

	product := router.Group("/product", auth.UserIdentity)
	{
		product.POST("/add", productHandlers.CreateProduct)
		product.GET("/:id", productHandlers.GetProductByID)
		product.PUT("/:id", productHandlers.UpdateProductByID)
		product.DELETE("/:id", productHandlers.DeleteProductByID)
	}

	return router
}
