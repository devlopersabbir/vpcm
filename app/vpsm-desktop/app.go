package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/devlopersabbir/vpcm/internal/config"
	"github.com/devlopersabbir/vpcm/internal/inventory"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) getAPIURL() string {
	cfg, err := config.Load()
	if err == nil && cfg.API.GlobalURL != "" {
		return cfg.API.GlobalURL
	}
	return "http://187.77.151.75:8080"
}

// GetServers returns the list of all registered servers via REST API
func (a *App) GetServers() ([]inventory.ServerView, error) {
	url := fmt.Sprintf("%s/servers", a.getAPIURL())
	ctx, cancel := context.WithTimeout(a.ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	var list []inventory.ServerView
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, err
	}
	return list, nil
}

// GetServer returns the detailed view of a single server via REST API
func (a *App) GetServer(id uint) (*inventory.ServerView, error) {
	url := fmt.Sprintf("%s/servers/%d", a.getAPIURL(), id)
	ctx, cancel := context.WithTimeout(a.ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	var view inventory.ServerView
	if err := json.NewDecoder(resp.Body).Decode(&view); err != nil {
		return nil, err
	}
	return &view, nil
}

// AddServer registers a new server in the inventory via REST API
func (a *App) AddServer(name string, host string, port int, username string, authType string, authSecret string) error {
	url := fmt.Sprintf("%s/servers", a.getAPIURL())
	srv := &inventory.Server{
		Name:       name,
		Host:       host,
		Port:       port,
		Username:   username,
		AuthType:   authType,
		AuthSecret: authSecret,
	}

	data, err := json.Marshal(srv)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(a.ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status: %d", resp.StatusCode)
	}
	return nil
}

// DeleteServer removes a server from the inventory via REST API
func (a *App) DeleteServer(id uint) error {
	url := fmt.Sprintf("%s/servers/%d", a.getAPIURL(), id)
	ctx, cancel := context.WithTimeout(a.ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status: %d", resp.StatusCode)
	}
	return nil
}

// GetConnectionHistory returns connection history logs for a server via REST API
func (a *App) GetConnectionHistory(id uint) ([]inventory.ConnectionLog, error) {
	url := fmt.Sprintf("%s/servers/%d/history", a.getAPIURL(), id)
	ctx, cancel := context.WithTimeout(a.ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	var logs []inventory.ConnectionLog
	if err := json.NewDecoder(resp.Body).Decode(&logs); err != nil {
		return nil, err
	}
	return logs, nil
}

// ScanServer triggers background data collection (OS, Hardware, Network, etc.) via REST API
func (a *App) ScanServer(id uint) error {
	url := fmt.Sprintf("%s/servers/%d/scan", a.getAPIURL(), id)
	ctx, cancel := context.WithTimeout(a.ctx, 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status: %d", resp.StatusCode)
	}
	return nil
}

// ToggleFavorite switches the favorite state of a server via REST API
func (a *App) ToggleFavorite(id uint) (bool, error) {
	url := fmt.Sprintf("%s/servers/%d/favorite", a.getAPIURL(), id)
	ctx, cancel := context.WithTimeout(a.ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return false, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	var res struct {
		IsFavorite bool `json:"is_favorite"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return false, err
	}
	return res.IsFavorite, nil
}

// GetConfig returns the active settings
func (a *App) GetConfig() (*config.Config, error) {
	return config.Load()
}

// SaveConfig updates the active settings
func (a *App) SaveConfig(cfg config.Config) error {
	return config.Save(&cfg)
}
