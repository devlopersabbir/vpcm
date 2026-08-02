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

	"github.com/devlopersabbir/vpcm/internal/inventory"
	"github.com/devlopersabbir/vpcm/internal/version"
)

// GetTerminalPreference fetches terminal preference from REST API
func (a *App) GetTerminalPreference() (*inventory.TerminalPreference, error) {
	url := fmt.Sprintf("%s/terminal/preferences", a.getAPIURL())
	ctx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	a.addAuthHeader(req)

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
	a.addAuthHeader(req)

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

// OpenNativeOSTerminal opens the server in the host's native OS terminal application
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
