package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/devlopersabbir/vpcm/internal/config"
	"github.com/devlopersabbir/vpcm/internal/inventory"
	"github.com/devlopersabbir/vpcm/internal/version"
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

// GetTerminalPreference fetches terminal preference from REST API
func (a *App) GetTerminalPreference() (*inventory.TerminalPreference, error) {
	url := fmt.Sprintf("%s/terminal/preferences", a.getAPIURL())
	ctx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
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
		return nil, fmt.Errorf("VPSM Desktop %s\nVPS Manager — Remote Server Inventory & SSH Terminal Panel\n\nDeveloped by @devlopersabbir", version.Version)
	}

	var pref inventory.TerminalPreference
	if err := json.NewDecoder(resp.Body).Decode(&pref); err != nil {
		return nil, err
	}
	return &pref, nil
}

// SaveTerminalPreference updates terminal preference via REST API
func (a *App) SaveTerminalPreference(pref inventory.TerminalPreference) error {
	url := fmt.Sprintf("%s/terminal/preferences", a.getAPIURL())
	ctx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
	defer cancel()

	bodyBytes, err := json.Marshal(pref)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}
	return nil
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

// OpenNativeOSTerminal opens the server in the host's native OS terminal application (Windows Terminal/PowerShell/CMD on Windows, Terminal.app on Mac, gnome-terminal/xterm on Linux)
func (a *App) OpenNativeOSTerminal(serverID uint, params SSHConnectionParams) error {
	host := params.Host
	port := params.Port
	if port == 0 {
		port = 22
	}
	user := params.Username
	if user == "" {
		user = "root"
	}

	sshCmd := fmt.Sprintf("vpcm ssh %s@%s -p %d", user, host, port)

	switch runtime.GOOS {
	case "windows":
		// Try Windows Terminal (wt.exe), fallback to powershell, fallback to cmd
		if _, err := exec.LookPath("wt.exe"); err == nil {
			cmd := exec.Command("wt.exe", "powershell", "-NoExit", "-Command", sshCmd)
			return cmd.Start()
		}
		if _, err := exec.LookPath("powershell.exe"); err == nil {
			cmd := exec.Command("cmd.exe", "/c", "start", "powershell", "-NoExit", "-Command", sshCmd)
			return cmd.Start()
		}
		cmd := exec.Command("cmd.exe", "/c", "start", "cmd.exe", "/k", sshCmd)
		return cmd.Start()

	case "darwin":
		script := fmt.Sprintf(`tell application "Terminal" to do script "%s"`, sshCmd)
		cmd := exec.Command("osascript", "-e", script)
		if err := cmd.Run(); err != nil {
			return err
		}
		return exec.Command("osascript", "-e", `tell application "Terminal" to activate`).Run()

	default: // Linux
		terminals := []string{"x-terminal-emulator", "gnome-terminal", "konsole", "xfce4-terminal", "alacritty", "kitty", "xterm"}
		for _, term := range terminals {
			if _, err := exec.LookPath(term); err == nil {
				var cmd *exec.Cmd
				if term == "gnome-terminal" || term == "xfce4-terminal" {
					cmd = exec.Command(term, "--", "bash", "-c", sshCmd+"; exec bash")
				} else {
					cmd = exec.Command(term, "-e", "bash", "-c", sshCmd+"; exec bash")
				}
				return cmd.Start()
			}
		}
		return fmt.Errorf("no supported terminal emulator found")
	}
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

// VerifyCloudPassword checks password with the specified cloud target URL
func (a *App) VerifyCloudPassword(targetURL string, password string) (bool, error) {
	if targetURL == "" {
		targetURL = a.getAPIURL()
	}
	url := fmt.Sprintf("%s/verify-cloud-access", targetURL)
	payload := map[string]string{"password": password}
	data, err := json.Marshal(payload)
	if err != nil {
		return false, err
	}

	ctx, cancel := context.WithTimeout(a.ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("Unable to connect to Cloud API server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	var res struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err == nil && res.Message != "" {
		return false, fmt.Errorf("%s", res.Message)
	}

	return false, fmt.Errorf("Access denied (status: %d)", resp.StatusCode)
}
