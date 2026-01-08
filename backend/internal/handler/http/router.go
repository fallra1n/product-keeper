package httphandler

import (
	"log/slog"

	"github.com/fallra1n/product-keeper/internal/handler/http/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRouter ...
func SetupRouter(log *slog.Logger, auth AuthHandler, productHandlers ProductsHandler) *gin.Engine {
	router := gin.Default()

	// TODO using custom logger

	router.POST("/user/register", auth.UserRegister)
	router.POST("/user/login", auth.UserLogin)

	products := router.Group("/products", middleware.UserIdentity())
	{
		products.GET("", productHandlers.FindProductList)
	}

	product := router.Group("/product", middleware.UserIdentity())
	{
		product.POST("/add", productHandlers.CreateProduct)
		product.GET("/:id", productHandlers.FindProduct)
		product.PUT("/:id", productHandlers.UpdateProduct)
		product.DELETE("/:id", productHandlers.DeleteProduct)
	}

	return router
}
