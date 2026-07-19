package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/devlopersabbir/vpcm/internal/config"
	"github.com/devlopersabbir/vpcm/internal/inventory"
	"github.com/devlopersabbir/vpcm/internal/ssh"
	"github.com/spf13/cobra"
)

var sshCmd = &cobra.Command{
	Use:   "ssh [id | name | username@host]",
	Short: "Connect to a server via SSH (saves credentials to database)",
	Args:  cobra.ExactArgs(1),
	RunE:  runSSHConnection,
}

func promptServerName(defaultName string) string {
	fmt.Printf("Enter a custom name for this server (default: %s): ", defaultName)
	var name string
	_, _ = fmt.Scanln(&name)
	name = strings.TrimSpace(name)
	if name == "" {
		return defaultName
	}
	return name
}

func runSSHConnection(cmd *cobra.Command, args []string) error {
	input := args[0]

	repo, svc, err := initRepoAndService(cmd.Context())
	if err != nil {
		return err
	}

	// Get all servers to check by ID/Name/Host
	servers, err := repo.List(cmd.Context())
	if err != nil {
		return err
	}

	var target *inventory.Server

	// Try ID match first
	var idVal uint64
	if id, err := strconv.ParseUint(input, 10, 32); err == nil {
		idVal = id
		for _, s := range servers {
			if s.ID == uint(idVal) {
				target = &s
				break
			}
		}
	}

	// Try Name match
	if target == nil {
		for _, s := range servers {
			if s.Name == input {
				target = &s
				break
			}
		}
	}

	// Try Host match or Parse username@host
	username := "root"
	host := input
	if target == nil {
		if idx := strings.Index(input, "@"); idx != -1 {
			username = input[:idx]
			host = input[idx+1:]
		}
		for _, s := range servers {
			if s.Host == host {
				target = &s
				break
			}
		}
	}

	// If still not found, create new record placeholder
	if target == nil {
		target = &inventory.Server{
			UUID:     "uuid-" + host,
			Name:     username + "@" + host,
			Host:     host,
			Port:     22,
			Username: username,
		}
	}

	if identityFile != "" {
		keyBytes, err := os.ReadFile(identityFile)
		if err != nil {
			return fmt.Errorf("failed to read identity file %s: %w", identityFile, err)
		}
		target.AuthType = "key"
		target.AuthSecret = string(keyBytes)
	}

	cfg, _ := config.Load()
	sshSvc := ssh.NewService(cfg.SSH.Timeout)
	var client ssh.Client

	// Try saved credentials first (or loaded key)
	if target.AuthSecret != "" {
		client, err = sshSvc.Connect(cmd.Context(), target.Host, target.Port, target.Username, target.AuthType, target.AuthSecret)
		if err == nil {
			// Save new server or update existing credentials
			if target.ID == 0 {
				target.Name = promptServerName(target.Name)
				target.Provider = inventory.DetectProvider(cmd.Context(), client, target.Host)
				target.OSFamily, target.OSVersion = inventory.DetectOS(cmd.Context(), client)
				target.CPUModel, target.CPUCores, target.RAMTotal, target.DiskTotal = inventory.DetectSpecs(cmd.Context(), client)
				if err := svc.AddServer(cmd.Context(), target); err != nil {
					return fmt.Errorf("failed to save server to database: %w", err)
				}
			} else {
				if target.Provider == "" || target.Provider == "Generic VPS" {
					target.Provider = inventory.DetectProvider(cmd.Context(), client, target.Host)
				}
				target.OSFamily, target.OSVersion = inventory.DetectOS(cmd.Context(), client)
				target.CPUModel, target.CPUCores, target.RAMTotal, target.DiskTotal = inventory.DetectSpecs(cmd.Context(), client)
				if err := svc.UpdateServer(cmd.Context(), target); err != nil {
					return fmt.Errorf("failed to update server: %w", err)
				}
			}
		}
	}

	// Prompt for password if no credentials or connection failed
	if client == nil || err != nil {
		if identityFile != "" {
			// Log failed attempt
			failLog := &inventory.ConnectionLog{
				ServerID:     target.ID,
				ServerName:   target.Name,
				Username:     target.Username,
				Host:         target.Host,
				LoggedInAt:   time.Now(),
				Status:       "failed",
				ErrorMessage: err.Error(),
			}
			_ = svc.LogConnectionEnd(cmd.Context(), failLog, err)
			return fmt.Errorf("SSH key connection failed: %w", err)
		}
		fmt.Printf("Enter SSH password for %s@%s: ", target.Username, target.Host)
		var password string
		_, err = fmt.Scanln(&password)
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}

		client, err = sshSvc.Connect(cmd.Context(), target.Host, target.Port, target.Username, "password", password)
		if err != nil {
			// Log failed attempt
			failLog := &inventory.ConnectionLog{
				ServerID:     target.ID,
				ServerName:   target.Name,
				Username:     target.Username,
				Host:         target.Host,
				LoggedInAt:   time.Now(),
				Status:       "failed",
				ErrorMessage: err.Error(),
			}
			_ = svc.LogConnectionEnd(cmd.Context(), failLog, err)
			return fmt.Errorf("SSH connection failed: %w", err)
		}

		// Save credentials on successful login
		target.AuthSecret = password
		target.AuthType = "password"

		if target.ID == 0 {
			target.Name = promptServerName(target.Name)
			target.Provider = inventory.DetectProvider(cmd.Context(), client, target.Host)
			target.OSFamily, target.OSVersion = inventory.DetectOS(cmd.Context(), client)
			target.CPUModel, target.CPUCores, target.RAMTotal, target.DiskTotal = inventory.DetectSpecs(cmd.Context(), client)
			if err := svc.AddServer(cmd.Context(), target); err != nil {
				return fmt.Errorf("failed to save server to database: %w", err)
			}
		} else {
			if target.Provider == "" || target.Provider == "Generic VPS" {
				target.Provider = inventory.DetectProvider(cmd.Context(), client, target.Host)
			}
			target.OSFamily, target.OSVersion = inventory.DetectOS(cmd.Context(), client)
			target.CPUModel, target.CPUCores, target.RAMTotal, target.DiskTotal = inventory.DetectSpecs(cmd.Context(), client)
			if err := svc.UpdateServer(cmd.Context(), target); err != nil {
				return fmt.Errorf("failed to update server credentials: %w", err)
			}
		}
	}

	// Log successful connection start
	logRecord, logErr := svc.LogConnectionStart(cmd.Context(), target)

	defer func() {
		if logErr == nil && logRecord != nil {
			_ = svc.LogConnectionEnd(cmd.Context(), logRecord, nil)
		}
	}()

	defer client.Close()
	return client.Shell(cmd.Context(), os.Stdin, os.Stdout, os.Stderr)
}
