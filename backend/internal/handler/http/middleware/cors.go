package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORS middleware для разрешения кросс-доменных запросов
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Разрешаем запросы с любых доменов (в продакшене лучше указать конкретные домены)
		c.Header("Access-Control-Allow-Origin", "*")
		
		// Разрешаем определенные методы
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		
		// Разрешаем определенные заголовки
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		
		// Разрешаем отправку credentials (cookies, authorization headers)
		c.Header("Access-Control-Allow-Credentials", "true")
		
		// Обрабатываем preflight запросы
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
}
