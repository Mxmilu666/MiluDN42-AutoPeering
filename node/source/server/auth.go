package server

import (
	"net/http"

	"github.com/Mxmilu666/MiluDN42-AutoPeering/node/source"
	"github.com/Mxmilu666/MiluDN42-AutoPeering/node/source/server/handles"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware 校验请求头中的 Token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			handles.SendResponse(c, http.StatusUnauthorized, "error", "Authorization token required")
			return
		}
		token := header
		const bearerPrefix = "Bearer "
		if len(header) > len(bearerPrefix) && header[:len(bearerPrefix)] == bearerPrefix {
			token = header[len(bearerPrefix):]
		}
		if source.AppConfig == nil || token != source.AppConfig.Token {
			handles.SendResponse(c, http.StatusUnauthorized, "error", "Invalid token")
			return
		}
		c.Next()
	}
}
