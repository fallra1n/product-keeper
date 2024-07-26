package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/fallra1n/product-keeper/pkg/jwt"
)

func UserIdentity() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader(AuthHeader)
		if header == "" {
			c.JSON(http.StatusUnauthorized, DefaultResponse{"empty auth header"})
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			c.JSON(http.StatusUnauthorized, DefaultResponse{"invalid auth header"})
			return
		}

		username, err := jwt.ParseToken(headerParts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, DefaultResponse{"invalid auth token"})
			return
		}

		c.Set(UserContext, username)
	}
}
