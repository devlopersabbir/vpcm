package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/devlopersabbir/vpcm/internal/config"
	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		logMsg := fmt.Sprintf("%s %s %s?%s | Status: %d | Latency: %v | IP: %s", method, clientIP, path, query, status, latency, clientIP)
		if status >= 500 {
			slog.Error(logMsg)
		} else if status >= 400 {
			slog.Warn(logMsg)
		} else {
			slog.Info(logMsg)
		}
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg, err := config.Load()
		if err != nil || cfg == nil {
			c.Next()
			return
		}

		expectedToken := cfg.API.Token
		// If no token is configured on server, allow direct access
		if expectedToken == "" {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		cloudAuthHeader := c.GetHeader("X-Cloud-Auth")
		token := ""

		if cloudAuthHeader != "" {
			token = cloudAuthHeader
		} else if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else if authHeader != "" {
			token = authHeader
		} else {
			token = c.Query("token")
		}

		if token == expectedToken {
			c.Next()
			return
		}

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: API token authentication required.",
		})
		c.Abort()
	}
}
