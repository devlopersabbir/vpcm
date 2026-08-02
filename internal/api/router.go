package api

import (
	"fmt"
	"time"

	"github.com/devlopersabbir/vpcm/internal/inventory"
	"github.com/devlopersabbir/vpcm/internal/notes"
	"github.com/gin-gonic/gin"
)

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
	// Public endpoint (Unprotected)
	s.router.POST("/verify-cloud-access", s.handleVerifyCloudGuard)

	// Protected endpoints (Requires authentication)
	apiGroup := s.router.Group("/")
	apiGroup.Use(AuthMiddleware())
	{
		apiGroup.GET("/", s.handleRoot)

		apiGroup.GET("/servers", s.handleListServers)
		apiGroup.POST("/servers", s.handleCreateServer)
		apiGroup.GET("/servers/:id", s.handleGetServer)
		apiGroup.DELETE("/servers/:id", s.handleDeleteServer)
		apiGroup.POST("/servers/:id/scan", s.handleScanServer)
		apiGroup.POST("/servers/:id/favorite", s.handleFavoriteServer)
		apiGroup.GET("/servers/:id/history", s.handleGetConnectionHistory)

		apiGroup.GET("/notes", s.handleListNotes)
		apiGroup.GET("/events", s.handleListEvents)

		apiGroup.GET("/terminal/preferences", s.handleGetTerminalPreferences)
		apiGroup.POST("/terminal/preferences", s.handleSaveTerminalPreferences)
	}
}

func (s *Server) Start(host string, port int) error {
	addr := fmt.Sprintf("%s:%d", host, port)
	return s.router.Run(addr)
}
