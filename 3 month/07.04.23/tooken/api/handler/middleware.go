package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		value := c.GetHeader("Password")
		if value != h.cfg.SecretKey {
			c.AbortWithError(http.StatusForbidden, errors.New("password invalid"))
			return
		}

		c.Next()
	}
}
