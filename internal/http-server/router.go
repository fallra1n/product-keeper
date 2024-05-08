package http

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/fallra1n/product-service/internal/http-server/handlers"
)

func SetupRouter(userHandlers handlers.UserHandler, productHandlers handlers.ProductHandler, logger *slog.Logger) *gin.Engine {
	router := gin.Default()

	// TODO using custom logger

	router.POST("/user/register", userHandlers.UserRegister)
	router.POST("/user/login", userHandlers.UserLogin)
	router.POST("/product/add", productHandlers.CreateProduct)
	router.GET("/products", productHandlers.GetProducts)
	router.GET("/product/:id", productHandlers.GetProductByID)
	router.PUT("/product/:id", productHandlers.ChangeProductByID)
	router.DELETE("/product/:id", productHandlers.DeleteProductByID)

	return router
}
