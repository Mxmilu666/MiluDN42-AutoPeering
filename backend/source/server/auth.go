package server

import (
	"net/http"

	"github.com/Mxmilu666/MiluDN42-AutoPeering/backend/source"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware 校验请求头中的 Token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			return
		}
		token := header
		const bearerPrefix = "Bearer "
		if len(header) > len(bearerPrefix) && header[:len(bearerPrefix)] == bearerPrefix {
			token = header[len(bearerPrefix):]
		}
		if source.AppConfig == nil || token != source.AppConfig.Token {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		c.Next()
	}
}
