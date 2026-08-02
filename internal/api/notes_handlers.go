package api

import (
	"net/http"
	"strconv"

	"github.com/devlopersabbir/vpcm/internal/inventory"
	"github.com/gin-gonic/gin"
)

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
	c.JSON(http.StatusOK, []any{})
}

func (s *Server) handleGetTerminalPreferences(c *gin.Context) {
	pref, err := s.inventoryService.GetTerminalPreference(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pref)
}

func (s *Server) handleSaveTerminalPreferences(c *gin.Context) {
	var pref inventory.TerminalPreference
	if err := c.ShouldBindJSON(&pref); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.inventoryService.SaveTerminalPreference(c.Request.Context(), &pref); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pref)
}
