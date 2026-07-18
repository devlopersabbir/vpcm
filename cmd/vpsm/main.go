package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/devlopersabbir/vpcm/internal/config"
	"github.com/devlopersabbir/vpcm/internal/database"
	"github.com/devlopersabbir/vpcm/internal/inventory"
	"github.com/devlopersabbir/vpcm/internal/logger"
	"github.com/devlopersabbir/vpcm/internal/ssh"
	"github.com/spf13/cobra"
)

var identityFile string

var rootCmd = &cobra.Command{
	Use:   "vpsm",
	Short: "VPSM - VPS Manager CLI",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		logger.Init(cfg.Logging.Level, cfg.Logging.Format, os.Stdout)
		return nil
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version info",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v0.1.0")
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage and show configuration status",
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show active configuration in a developer-friendly table",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.Load()

		re := lipgloss.NewRenderer(os.Stdout)
		cyan := lipgloss.Color("#00FFFF")
		purple := lipgloss.Color("#7D56F4")
		gray := lipgloss.Color("#888888")
		white := lipgloss.Color("#FFFFFF")
		green := lipgloss.Color("#00FF00")
		red := lipgloss.Color("#FF0000")

		var rows [][]string
		rows = append(rows, []string{"Database Driver", cfg.Database.Driver})
		if cfg.Database.Driver == "sqlite" {
			rows = append(rows, []string{"SQLite DB Path", cfg.Database.Path})
		} else {
			rows = append(rows, []string{"MongoDB URI", cfg.Database.URI})
			rows = append(rows, []string{"MongoDB Name", cfg.Database.Name})
		}

		apiStatus := "Disabled"
		if cfg.API.Enabled {
			apiStatus = "Enabled"
		}
		rows = append(rows, []string{"API Server Status", apiStatus})
		rows = append(rows, []string{"API Host", cfg.API.Host})
		rows = append(rows, []string{"API Port", strconv.Itoa(cfg.API.Port)})
		rows = append(rows, []string{"SSH Timeout", cfg.SSH.Timeout.String()})
		rows = append(rows, []string{"Logging Level", cfg.Logging.Level})
		rows = append(rows, []string{"Logging Format", cfg.Logging.Format})

		tbl := table.New().
			Border(lipgloss.RoundedBorder()).
			BorderStyle(re.NewStyle().Foreground(gray)).
			Headers("Configuration Key", "Active Value").
			Rows(rows...)

		tbl.StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == -1:
				return re.NewStyle().Bold(true).Foreground(cyan).Padding(0, 1)
			default:
				if col == 0 {
					return re.NewStyle().Bold(true).Foreground(purple).Padding(0, 1)
				}
				val := rows[row][1]
				if val == "Enabled" {
					return re.NewStyle().Foreground(green).Padding(0, 1)
				}
				if val == "Disabled" {
					return re.NewStyle().Foreground(red).Padding(0, 1)
				}
				return re.NewStyle().Foreground(white).Padding(0, 1)
			}
		})

		fmt.Println("\n 🛠️  VPSM Current Configuration:")
		fmt.Println(tbl)
		return nil
	},
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Interactively initialize/configure your universe settings",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.Load()

		fmt.Println("🚀 Initializing VPSM Configuration...")

		fmt.Print("Choose database driver [sqlite/mongodb] (default: sqlite): ")
		var driver string
		_, _ = fmt.Scanln(&driver)
		driver = strings.ToLower(strings.TrimSpace(driver))
		if driver == "" {
			driver = "sqlite"
		}
		if driver != "sqlite" && driver != "mongodb" {
			return fmt.Errorf("invalid driver choice: must be 'sqlite' or 'mongodb'")
		}
		cfg.Database.Driver = driver

		if driver == "sqlite" {
			defaultPath := cfg.Database.Path
			if defaultPath == "" {
				defaultPath = "~/.local/share/vpsm/vpsm.db"
			}
			fmt.Printf("Enter SQLite database path (default: %s): ", defaultPath)
			var path string
			_, _ = fmt.Scanln(&path)
			path = strings.TrimSpace(path)
			if path != "" {
				cfg.Database.Path = path
			} else {
				cfg.Database.Path = defaultPath
			}
		} else {
			defaultURI := cfg.Database.URI
			if defaultURI == "" {
				defaultURI = "mongodb://localhost:27017"
			}
			fmt.Printf("Enter MongoDB Connection URI (default: %s): ", defaultURI)
			var uri string
			_, _ = fmt.Scanln(&uri)
			uri = strings.TrimSpace(uri)
			if uri != "" {
				cfg.Database.URI = uri
			} else {
				cfg.Database.URI = defaultURI
			}

			defaultName := cfg.Database.Name
			if defaultName == "" {
				defaultName = "vpsm"
			}
			fmt.Printf("Enter MongoDB Database Name (default: %s): ", defaultName)
			var name string
			_, _ = fmt.Scanln(&name)
			name = strings.TrimSpace(name)
			if name != "" {
				cfg.Database.Name = name
			} else {
				cfg.Database.Name = defaultName
			}
		}

		fmt.Print("Enable local REST API server? [y/N] (default: n): ")
		var enableAPI string
		_, _ = fmt.Scanln(&enableAPI)
		enableAPI = strings.ToLower(strings.TrimSpace(enableAPI))
		if enableAPI == "y" || enableAPI == "yes" {
			cfg.API.Enabled = true
		} else {
			cfg.API.Enabled = false
		}

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		fmt.Println("\n✨ Configuration saved successfully to ~/.config/vpsm/config.yaml!")
		return nil
	},
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Verify systems setup",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[✓] CLI framework initialized")
		fmt.Println("[✓] Logger initialized")
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("[✗] Config error: %v\n", err)
			return
		}
		fmt.Println("[✓] Config validated")

		var dbErr error
		if cfg.Database.Driver == "sqlite" {
			_, dbErr = database.InitSQLite(cfg.Database.Path)
		} else {
			_, dbErr = database.InitMongo(cfg.Database.URI, cfg.Database.Name)
		}
		if dbErr != nil {
			fmt.Printf("[✗] Database connection failed: %v\n", dbErr)
			return
		}
		fmt.Println("[✓] Database connection succeeded")
	},
}

// Server commands
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage servers in inventory",
}

func initRepoAndService(ctx context.Context) (inventory.ServerRepository, inventory.ServerService, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, nil, err
	}

	if cfg.Database.Driver == "sqlite" {
		db, err := database.InitSQLite(cfg.Database.Path)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to initialize SQLite: %w", err)
		}
		repo := inventory.NewSQLiteRepository(db)
		return repo, inventory.NewService(repo), nil
	}

	db, err := database.InitMongo(cfg.Database.URI, cfg.Database.Name)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize MongoDB: %w", err)
	}
	repo := inventory.NewMongoRepository(db)
	return repo, inventory.NewService(repo), nil
}

var serverListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all monitored servers",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, _, err := initRepoAndService(cmd.Context())
		if err != nil {
			return err
		}
		servers, err := repo.List(cmd.Context())
		if err != nil {
			return err
		}

		re := lipgloss.NewRenderer(os.Stdout)
		purple := lipgloss.Color("#7D56F4")
		gray := lipgloss.Color("#3C3C3C")

		var rows [][]string
		for _, s := range servers {
			rows = append(rows, []string{
				strconv.Itoa(int(s.ID)),
				s.Name,
				s.Username,
				s.Host,
				strconv.Itoa(s.Port),
				s.AuthType,
				s.Provider,
			})
		}

		tbl := table.New().
			Border(lipgloss.NormalBorder()).
			BorderStyle(re.NewStyle().Foreground(gray)).
			Headers("ID", "Name", "Username", "Host", "Port", "Auth Type", "Provider").
			Rows(rows...)

		tbl.StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == -1:
				return re.NewStyle().Bold(true).Foreground(purple).Padding(0, 1)
			default:
				return re.NewStyle().Padding(0, 1)
			}
		})

		fmt.Println(tbl)
		return nil
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all monitored servers (alias for server list)",
	RunE:  serverListCmd.RunE,
}

var serverAddCmd = &cobra.Command{
	Use:   "add [name] [host]",
	Short: "Add a new server to inventory",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		host := args[1]

		_, svc, err := initRepoAndService(cmd.Context())
		if err != nil {
			return err
		}

		server := &inventory.Server{
			UUID:     "dummy-uuid-for-v0.0.1",
			Name:     name,
			Host:     host,
			Port:     22,
			Username: "root",
			Provider: inventory.DetectProvider(cmd.Context(), nil, host),
		}

		if err := svc.AddServer(cmd.Context(), server); err != nil {
			return err
		}

		fmt.Printf("Successfully added server %s (%s)\n", name, host)
		return nil
	},
}

var serverRemoveCmd = &cobra.Command{
	Use:   "remove [id | name]",
	Short: "Remove a server from inventory",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]

		repo, svc, err := initRepoAndService(cmd.Context())
		if err != nil {
			return err
		}

		servers, err := repo.List(cmd.Context())
		if err != nil {
			return err
		}

		var targetID uint
		var found bool

		// Check by ID first
		if id, err := strconv.ParseUint(input, 10, 32); err == nil {
			for _, s := range servers {
				if s.ID == uint(id) {
					targetID = s.ID
					found = true
					break
				}
			}
		}

		// Check by Name
		if !found {
			for _, s := range servers {
				if s.Name == input {
					targetID = s.ID
					found = true
					break
				}
			}
		}

		if !found {
			return fmt.Errorf("server '%s' not found in database", input)
		}

		if err := svc.RemoveServer(cmd.Context(), targetID); err != nil {
			return err
		}

		fmt.Printf("Successfully removed server '%s' (ID: %d) from database.\n", input, targetID)
		return nil
	},
}

var serverFlushCmd = &cobra.Command{
	Use:     "flush",
	Aliases: []string{"flash"},
	Short:   "Fully clean/flush the database of all servers (requires double confirmation)",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Print("Are you sure you want to flush all servers? [y/N]: ")
		var confirmation1 string
		_, _ = fmt.Scanln(&confirmation1)
		confirmation1 = strings.ToLower(strings.TrimSpace(confirmation1))
		if confirmation1 != "y" && confirmation1 != "yes" {
			fmt.Println("Aborted.")
			return nil
		}

		fmt.Print("This action is irreversible. Type 'FLUSH' to confirm: ")
		var confirmation2 string
		_, _ = fmt.Scanln(&confirmation2)
		confirmation2 = strings.TrimSpace(confirmation2)
		if confirmation2 != "FLUSH" {
			fmt.Println("Aborted (confirmation text did not match 'FLUSH').")
			return nil
		}

		_, svc, err := initRepoAndService(cmd.Context())
		if err != nil {
			return err
		}

		if err := svc.FlushServers(cmd.Context()); err != nil {
			return fmt.Errorf("failed to flush database: %w", err)
		}

		fmt.Println("Successfully flushed/cleaned all servers from the database.")
		return nil
	},
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
				if err := svc.AddServer(cmd.Context(), target); err != nil {
					return fmt.Errorf("failed to save server to database: %w", err)
				}
			} else {
				if target.Provider == "" || target.Provider == "Generic VPS" {
					target.Provider = inventory.DetectProvider(cmd.Context(), client, target.Host)
				}
				if identityFile != "" || target.Provider != "" {
					if err := svc.UpdateServer(cmd.Context(), target); err != nil {
						return fmt.Errorf("failed to update server: %w", err)
					}
				}
			}
		}
	}

	// Prompt for password if no credentials or connection failed
	if client == nil || err != nil {
		if identityFile != "" {
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
			return fmt.Errorf("SSH connection failed: %w", err)
		}

		// Save credentials on successful login
		target.AuthSecret = password
		target.AuthType = "password"

		if target.ID == 0 {
			target.Name = promptServerName(target.Name)
			target.Provider = inventory.DetectProvider(cmd.Context(), client, target.Host)
			if err := svc.AddServer(cmd.Context(), target); err != nil {
				return fmt.Errorf("failed to save server to database: %w", err)
			}
		} else {
			if target.Provider == "" || target.Provider == "Generic VPS" {
				target.Provider = inventory.DetectProvider(cmd.Context(), client, target.Host)
			}
			if err := svc.UpdateServer(cmd.Context(), target); err != nil {
				return fmt.Errorf("failed to update server credentials: %w", err)
			}
		}
	}

	defer client.Close()
	return client.Shell(cmd.Context(), os.Stdin, os.Stdout, os.Stderr)
}

var sshCmd = &cobra.Command{
	Use:   "ssh [id | name | username@host]",
	Short: "Connect to a server via SSH (saves credentials to database)",
	Args:  cobra.ExactArgs(1),
	RunE:  runSSHConnection,
}

func init() {
	configCmd.AddCommand(configShowCmd, configInitCmd)
	serverCmd.AddCommand(serverListCmd, serverAddCmd, serverRemoveCmd, serverFlushCmd)
	sshCmd.Flags().StringVarP(&identityFile, "identity", "i", "", "identity file (private key)")
	rootCmd.AddCommand(versionCmd, configCmd, doctorCmd, serverCmd, sshCmd, listCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
