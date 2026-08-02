package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/devlopersabbir/vpcm/internal/config"
	"github.com/devlopersabbir/vpcm/internal/database"
)

// GetConfig returns the active settings
func (a *App) GetConfig() (*config.Config, error) {
	return config.Load()
}

// SaveConfig updates the active settings
func (a *App) SaveConfig(cfg config.Config) error {
	return config.Save(&cfg)
}

// TestDatabaseConnection checks connection to SQLite or MongoDB based on provided settings
func (a *App) TestDatabaseConnection(driver string, path string, uri string, dbName string) (map[string]interface{}, error) {
	if driver == "mongodb" {
		if uri == "" {
			return map[string]interface{}{"success": false, "message": "MongoDB connection URI is empty"}, nil
		}
		if dbName == "" {
			dbName = "vpsm"
		}
		db, err := database.InitMongo(uri, dbName)
		if err != nil {
			return map[string]interface{}{"success": false, "message": fmt.Sprintf("MongoDB connection failed: %v", err)}, nil
		}
		_ = db
		return map[string]interface{}{"success": true, "message": "MongoDB connection successful!"}, nil
	}

	// SQLite check
	if path == "" {
		path = filepath.Join(os.Getenv("HOME"), ".local", "share", "vpsm", "vpsm.db")
	}
	db, err := database.InitSQLite(path)
	if err != nil {
		return map[string]interface{}{"success": false, "message": fmt.Sprintf("SQLite connection failed: %v", err)}, nil
	}
	_ = db
	return map[string]interface{}{"success": true, "message": "SQLite database path is valid and accessible!"}, nil
}

// TestAPIConnection checks HTTP connection to the API server endpoint
func (a *App) TestAPIConnection(apiURL string) (map[string]interface{}, error) {
	if apiURL == "" {
		return map[string]interface{}{"success": false, "message": "API endpoint URL is empty"}, nil
	}

	ctx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return map[string]interface{}{"success": false, "message": fmt.Sprintf("Invalid URL format: %v", err)}, nil
	}
	a.addAuthHeader(req)

	resp, err := httpClient.Do(req)
	if err != nil {
		return map[string]interface{}{"success": false, "message": fmt.Sprintf("Unable to connect to API server: %v", err)}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return map[string]interface{}{"success": true, "message": fmt.Sprintf("API server reached successfully (HTTP %d)", resp.StatusCode)}, nil
	}

	return map[string]interface{}{"success": false, "message": fmt.Sprintf("API server returned status HTTP %d", resp.StatusCode)}, nil
}
