package api

import (
	"context"
	"log/slog"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/devlopersabbir/vpcm/internal/inventory"
	"github.com/gin-gonic/gin"
)

type endpointDoc struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	Description string `json:"description"`
}

func (s *Server) handleRoot(c *gin.Context) {
	uptime := time.Since(startedAt).Round(time.Second).String()

	endpoints := []endpointDoc{
		{"GET", "/", "API info & available endpoints (this page)"},
		{"GET", "/servers", "List all servers with full inventory (network, hardware, OS)"},
		{"POST", "/servers", "Add a new server to the inventory"},
		{"GET", "/servers/:id", "Get a single server by ID with full inventory"},
		{"DELETE", "/servers/:id", "Delete a server from inventory"},
		{"POST", "/servers/:id/scan", "Trigger a live inventory scan for a server"},
		{"POST", "/servers/:id/favorite", "Toggle favorite status for a server"},
		{"GET", "/servers/:id/history", "Get SSH connection history for a server"},
		{"GET", "/notes?server_id=:id", "List notes attached to a server"},
		{"GET", "/events", "Stream or list recent system events"},
	}

	c.JSON(http.StatusOK, gin.H{
		"name":        "VPSM API",
		"description": "VPS Manager — Remote server inventory & SSH session tracking",
		"version":     "v1.2.0",
		"status":      "ok",
		"uptime":      uptime,
		"started_at":  startedAt.UTC().Format(time.RFC3339),
		"runtime":     runtime.Version(),
		"os_arch":     runtime.GOOS + "/" + runtime.GOARCH,
		"docs":        "https://github.com/devlopersabbir/vpcm/blob/main/docs/api.md",
		"endpoints":   endpoints,
	})
}

func (s *Server) handleListServers(c *gin.Context) {
	views, err := s.inventoryService.ListServers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, views)
}

func (s *Server) handleGetServer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	view, err := s.inventoryService.GetServer(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, view)
}

func (s *Server) handleScanServer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	go func() {
		ctx := context.Background()
		if err := s.inventoryService.ScanInventory(ctx, uint(id)); err != nil {
			slog.Error("Background specs scan failed", "server_id", id, "error", err)
		}
	}()
	c.JSON(http.StatusAccepted, gin.H{"status": "scan initiated"})
}

func (s *Server) handleCreateServer(c *gin.Context) {
	var srv inventory.Server
	if err := c.ShouldBindJSON(&srv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.inventoryService.AddServer(c.Request.Context(), &srv); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, srv)
}

func (s *Server) handleDeleteServer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	if err := s.inventoryService.RemoveServer(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "server deleted"})
}

func (s *Server) handleFavoriteServer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	fav, err := s.inventoryService.ToggleFavorite(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"is_favorite": fav})
}

func (s *Server) handleGetConnectionHistory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	history, err := s.inventoryService.GetConnectionHistory(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, history)
}
