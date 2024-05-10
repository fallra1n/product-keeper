package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	AuthHeader = "Authorization"
)

func (ah *authHandler) UserIdentity(c *gin.Context) {
	header := c.GetHeader(AuthHeader)
	if header == "" {
		ah.logger.Error("UserIdentity: auth header is empty")
		c.JSON(http.StatusUnauthorized, DefaultResponse{"empty auth header"})
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		ah.logger.Error("UserIdentity: invalid auth header")
		c.JSON(http.StatusUnauthorized, DefaultResponse{"invalid auth header"})
		return
	}

	username, err := ah.services.ParseToken(headerParts[1])
	if err != nil {
		ah.logger.Error("UserIdentity: " + err.Error())
		c.JSON(http.StatusUnauthorized, DefaultResponse{"invalid auth token"})
		return
	}

	c.Set("username", username)
}
