package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/devlopersabbir/vpcm/internal/inventory"
	"github.com/devlopersabbir/vpcm/internal/notes"
	"github.com/gin-gonic/gin"
)

// startedAt records when the API server process started (used for uptime).
var startedAt = time.Now()

type Server struct {
	router           *gin.Engine
	inventoryService inventory.ServerService
	noteService      notes.NoteService
}

func NewServer(invSvc inventory.ServerService, noteSvc notes.NoteService) *Server {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(LoggerMiddleware(), gin.Recovery())

	s := &Server{
		router:           r,
		inventoryService: invSvc,
		noteService:      noteSvc,
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router.GET("/", s.handleRoot)

	s.router.GET("/servers", s.handleListServers)
	s.router.GET("/servers/:id", s.handleGetServer)
	s.router.POST("/servers/:id/scan", s.handleScanServer)

	s.router.GET("/notes", s.handleListNotes)
	s.router.GET("/events", s.handleListEvents)
}

func (s *Server) Start(host string, port int) error {
	addr := fmt.Sprintf("%s:%d", host, port)
	return s.router.Run(addr)
}

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
		{"GET", "/servers/:id", "Get a single server by ID with full inventory"},
		{"POST", "/servers/:id/scan", "Trigger a live inventory scan for a server"},
		{"GET", "/notes?server_id=:id", "List notes attached to a server"},
		{"GET", "/events", "Stream or list recent system events"},
	}

	c.JSON(http.StatusOK, gin.H{
		"name":        "VPSM API",
		"description": "VPS Manager — Remote server inventory & SSH session tracking",
		"version":     "v0.1.5",
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

	if err := s.inventoryService.ScanInventory(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "scan initiated"})
}

func (s *Server) handleListNotes(c *gin.Context) {
	serverIDStr := c.Query("server_id")
	if serverIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "server_id parameter required"})
		return
	}

	serverID, err := strconv.ParseUint(serverIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid server_id"})
		return
	}

	notesList, err := s.noteService.GetServerNotes(c.Request.Context(), uint(serverID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notesList)
}

func (s *Server) handleListEvents(c *gin.Context) {
	// Return skeleton event list
	c.JSON(http.StatusOK, []any{})
}

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

		// Log request details
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
