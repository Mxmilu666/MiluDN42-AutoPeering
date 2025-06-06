package server

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/Mxmilu666/MiluDN42-AutoPeering/center/source"
	logger "github.com/Mxmilu666/MiluDN42-AutoPeering/center/source/logger"
	"github.com/gin-gonic/gin"
)

func Setupserver() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(filterLogs())
	r = initRouter(r)

	// start http server
	address := fmt.Sprintf("%s:%d", source.AppConfig.Listen.IP, source.AppConfig.Listen.Port)
	logger.Info("Starting server on", "address", address)
	if err := r.Run(address); err != nil {
		logger.Fatal("Could not start server", "error", err)
	}
}

func filterLogs() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		fullURL, err := url.QueryUnescape(c.Request.URL.String())
		if err != nil {
			logger.Error("Error decoding URL", "error", err)
			return
		}

		fullURL = strings.ReplaceAll(fullURL, "%", "%%")

		c.Next()

		latency := time.Since(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		userAgent := c.Request.UserAgent()

		logger.Info("HTTP Request",
			"status", statusCode,
			"latency", latency,
			"ip", clientIP,
			"method", method,
			"agent", userAgent,
			"url", fullURL,
		)
	}
}
