package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/devlopersabbir/vpcm/internal/inventory"
	"github.com/devlopersabbir/vpcm/internal/notes"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router           *gin.Engine
	inventoryService inventory.ServerService
	noteService      notes.NoteService
}

func NewServer(invSvc inventory.ServerService, noteSvc notes.NoteService) *Server {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	s := &Server{
		router:           r,
		inventoryService: invSvc,
		noteService:      noteSvc,
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router.GET("/servers", s.handleListServers)
	s.router.POST("/servers/:id/scan", s.handleScanServer)

	s.router.GET("/notes", s.handleListNotes)
	s.router.GET("/events", s.handleListEvents)
}

func (s *Server) Start(host string, port int) error {
	addr := fmt.Sprintf("%s:%d", host, port)
	return s.router.Run(addr)
}

func (s *Server) handleListServers(c *gin.Context) {
	servers, err := s.inventoryService.ListServers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, servers)
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
