package api

import (
	"net/http"
	"os"

	"github.com/devlopersabbir/vpcm/internal/config"
	"github.com/gin-gonic/gin"
)

type verifyCloudGuardRequest struct {
	Password string `json:"password"`
}

func (s *Server) handleVerifyCloudGuard(c *gin.Context) {
	var req verifyCloudGuardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"valid": false, "message": "Invalid request body"})
		return
	}

	cfg, _ := config.Load()
	expectedPass := ""
	if cfg != nil && cfg.API.GuardPassword != "" {
		expectedPass = cfg.API.GuardPassword
	} else {
		expectedPass = os.Getenv("CLOUD_GUARD_PASSWORD")
	}

	// If no guard password is set yet, seed the password on first input
	if expectedPass == "" {
		if req.Password != "" && cfg != nil {
			cfg.API.GuardPassword = req.Password
			_ = config.Save(cfg)
		}
		c.JSON(http.StatusOK, gin.H{
			"valid":   true,
			"message": "Cloud access authorized",
		})
		return
	}

	if req.Password == expectedPass {
		c.JSON(http.StatusOK, gin.H{
			"valid":   true,
			"message": "Cloud access authorized",
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"valid":   false,
			"message": "Invalid Cloud access password",
		})
	}
}
