package cli

import (
	"context"
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
	Use:               "ssh [id | name | username@host]",
	Short:             "Connect to a server via SSH (saves credentials to database)",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: serverNameCompletions,
	RunE:              runSSHConnection,
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

// collectAndPersistMetadata runs all remote detectors and upserts results into
// the three child tables (server_network, server_hardware, server_os).
// It is called after every successful SSH connect.
func collectAndPersistMetadata(ctx context.Context, svc inventory.ServerService, client ssh.Client, target *inventory.Server) {
	// Provider (only re-detect if unknown)
	if target.Provider == "" || target.Provider == "Generic VPS" {
		target.Provider = inventory.DetectProvider(ctx, client, target.Host)
		_ = svc.UpdateServer(ctx, target)
	}

	// OS info
	osFamily, osVersion := inventory.DetectOS(ctx, client)

	// Hardware
	cpuModel, cpuCores, ramTotal, diskTotal := inventory.DetectSpecs(ctx, client)

	// Full server info (network, firmware, OS extras)
	si := inventory.DetectServerInfo(ctx, client)

	// Upsert child tables
	_ = svc.UpsertServerNetwork(ctx, &inventory.ServerNetwork{
		ServerID:         target.ID,
		Hostname:         si.Hostname,
		PublicIP:         si.PublicIP,
		PrivateIP:        si.PrivateIP,
		MACAddress:       si.MACAddress,
		Region:           si.Region,
		AvailabilityZone: si.AvailabilityZone,
	})

	_ = svc.UpsertServerHardware(ctx, &inventory.ServerHardware{
		ServerID:       target.ID,
		CPUModel:       cpuModel,
		CPUCores:       cpuCores,
		RAMTotal:       ramTotal,
		SwapTotal:      si.SwapTotal,
		DiskTotal:      diskTotal,
		Virtualization: si.Virtualization,
		InstanceType:   si.InstanceType,
		SerialNumber:   si.SerialNumber,
		BIOSVersion:    si.BIOSVersion,
		Uptime:         si.Uptime,
	})

	_ = svc.UpsertServerOS(ctx, &inventory.ServerOS{
		ServerID:       target.ID,
		OSFamily:       osFamily,
		OSVersion:      osVersion,
		KernelVersion:  si.KernelVersion,
		Architecture:   si.Architecture,
		InitSystem:     si.InitSystem,
		Timezone:       si.Timezone,
		Locale:         si.Locale,
		PackageManager: si.PackageManager,
	})
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
		}	}

	if identityFile != "" {
		if _, statErr := os.Stat(identityFile); statErr != nil {
			if os.IsNotExist(statErr) {
				return fmt.Errorf("SSH key file not found: %s", identityFile)
			}
			return fmt.Errorf("failed to access SSH key file %s: %w", identityFile, statErr)
		}
		keyBytes, err := os.ReadFile(identityFile)
		if err == nil {
			target.AuthType = "key"
			target.AuthSecret = string(keyBytes)
		} else {
			target.AuthType = "keyfile"
			target.AuthSecret = identityFile
		}
	}

	// Always update database record if existing server found
	if target != nil && target.ID != 0 {
		now := time.Now()
		target.LastSeen = &now
		_ = svc.UpdateServer(cmd.Context(), target)
	}

	cfg, _ := config.Load()
	sshSvc := ssh.NewService(cfg.SSH.Timeout)
	var client ssh.Client

	// ── Try saved credentials (or loaded key) ─────────────────────────────────
	if target.AuthSecret != "" {
		client, err = sshSvc.Connect(cmd.Context(), target.Host, target.Port, target.Username, target.AuthType, target.AuthSecret)
		if err == nil {
			if target.ID == 0 {
				target.Name = promptServerName(target.Name)
				target.Provider = inventory.DetectProvider(cmd.Context(), client, target.Host)
				if err := svc.AddServer(cmd.Context(), target); err != nil {
					return fmt.Errorf("failed to save server to database: %w", err)
				}
			} else {
				if err := svc.UpdateServer(cmd.Context(), target); err != nil {
					return fmt.Errorf("failed to update server: %w", err)
				}
			}
			collectAndPersistMetadata(cmd.Context(), svc, client, target)
		}
	}

	// ── Prompt for password if no credentials or connection failed ─────────────
	if client == nil || err != nil {
		if identityFile != "" {
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

		target.AuthSecret = password
		target.AuthType = "password"

		if target.ID == 0 {
			target.Name = promptServerName(target.Name)
			target.Provider = inventory.DetectProvider(cmd.Context(), client, target.Host)
			if err := svc.AddServer(cmd.Context(), target); err != nil {
				return fmt.Errorf("failed to save server to database: %w", err)
			}
		} else {
			if err := svc.UpdateServer(cmd.Context(), target); err != nil {
				return fmt.Errorf("failed to update server credentials: %w", err)
			}
		}
		collectAndPersistMetadata(cmd.Context(), svc, client, target)
	}

	// ── Log connection ─────────────────────────────────────────────────────────
	logRecord, logErr := svc.LogConnectionStart(cmd.Context(), target)
	defer func() {
		if logErr == nil && logRecord != nil {
			_ = svc.LogConnectionEnd(cmd.Context(), logRecord, nil)
		}
	}()

	defer client.Close()
	return client.Shell(cmd.Context(), os.Stdin, os.Stdout, os.Stderr)
}
