package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
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

		expectedPass := cfg.API.GuardPassword
		expectedToken := cfg.API.Token
		envPass := os.Getenv("CLOUD_GUARD_PASSWORD")

		// If no password/token is set yet on server, allow access (unseeded initial mode)
		if expectedPass == "" && expectedToken == "" && envPass == "" {
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

		if (expectedPass != "" && token == expectedPass) ||
			(expectedToken != "" && token == expectedToken) ||
			(envPass != "" && token == envPass) {
			c.Next()
			return
		}

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: Cloud guard authentication required. Verify access via /verify-cloud-access first.",
		})
		c.Abort()
	}
}
