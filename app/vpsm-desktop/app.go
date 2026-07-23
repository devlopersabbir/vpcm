package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/devlopersabbir/vpcm/internal/config"
	"github.com/devlopersabbir/vpcm/internal/inventory"
)

var httpClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     5 * time.Second,
	},
}

// App struct
type App struct {
	ctx              context.Context
	isTerminalOnly   bool
	terminalServerID uint
	terminalParams   SSHConnectionParams
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

	resp, err := httpClient.Do(req)
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

	resp, err := httpClient.Do(req)
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

	resp, err := httpClient.Do(req)
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

	resp, err := httpClient.Do(req)
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

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return []inventory.ConnectionLog{}, nil
	}

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

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
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

	resp, err := httpClient.Do(req)
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

// OpenStandaloneTerminalWindow spawns an independent native application window containing our custom terminal
func (a *App) OpenStandaloneTerminalWindow(serverID uint, params SSHConnectionParams) error {
	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	args := []string{}
	if serverID > 0 {
		args = append(args, fmt.Sprintf("-terminal-server-id=%d", serverID))
	}
	if params.Host != "" {
		args = append(args, fmt.Sprintf("-terminal-host=%s", params.Host))
		args = append(args, fmt.Sprintf("-terminal-port=%d", params.Port))
		args = append(args, fmt.Sprintf("-terminal-user=%s", params.Username))
		args = append(args, fmt.Sprintf("-terminal-authtype=%s", params.AuthType))
		args = append(args, fmt.Sprintf("-terminal-authsecret=%s", params.AuthSecret))
	}

	cmd := exec.Command(execPath, args...)
	cmd.Env = os.Environ()
	return cmd.Start()
}

// GetTerminalInitialParams returns initial params if app was launched as terminal-only window
func (a *App) GetTerminalInitialParams() (map[string]interface{}, error) {
	if !a.isTerminalOnly {
		return map[string]interface{}{"is_terminal_only": false}, nil
	}

	// If server ID was passed, fetch the full fresh server record directly from the REST API
	if a.terminalServerID > 0 {
		srv, err := a.GetServer(a.terminalServerID)
		if err == nil && srv != nil {
			return map[string]interface{}{
				"is_terminal_only": true,
				"params": map[string]interface{}{
					"host":        srv.Host,
					"port":        srv.Port,
					"username":    srv.Username,
					"auth_type":   srv.AuthType,
					"auth_secret": srv.AuthSecret,
				},
			}, nil
		}
	}

	return map[string]interface{}{
		"is_terminal_only": true,
		"params":           a.terminalParams,
	}, nil
}
