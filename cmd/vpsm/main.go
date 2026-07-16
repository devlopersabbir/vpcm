package main

import (
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
		fmt.Println("v0.0.1")
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show configuration status",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.Load()
		fmt.Printf("Database Driver: %s\n", cfg.Database.Driver)
		fmt.Printf("Database Path:   %s\n", cfg.Database.Path)
		fmt.Printf("API Enabled:     %v\n", cfg.API.Enabled)
		fmt.Printf("Logging Level:   %s\n", cfg.Logging.Level)
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

		db, err := database.Init(cfg.Database.Path)
		if err != nil {
			fmt.Printf("[✗] Database connection failed: %v\n", err)
			return
		}
		fmt.Println("[✓] Database connection succeeded")

		// Migrate skeleton schema
		err = db.AutoMigrate(&inventory.Server{}, &inventory.Tag{}, &inventory.Software{})
		if err != nil {
			fmt.Printf("[✗] Migration failed: %v\n", err)
			return
		}
		fmt.Println("[✓] Schema integrity checked")
	},
}

// Server commands
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage servers in inventory",
}

var serverListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all monitored servers",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.Load()
		db, err := database.Init(cfg.Database.Path)
		if err != nil {
			return err
		}
		repo := inventory.NewRepository(db)
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
				s.Provider,
			})
		}

		tbl := table.New().
			Border(lipgloss.NormalBorder()).
			BorderStyle(re.NewStyle().Foreground(gray)).
			Headers("ID", "Name", "Username", "Host", "Port", "Provider").
			Rows(rows...)

		tbl.StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == 0:
				return re.NewStyle().Bold(true).Foreground(purple).Padding(0, 1)
			default:
				return re.NewStyle().Padding(0, 1)
			}
		})

		fmt.Println(tbl)
		return nil
	},
}

var serverAddCmd = &cobra.Command{
	Use:   "add [name] [host]",
	Short: "Add a new server to inventory",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		host := args[1]

		cfg, _ := config.Load()
		db, err := database.Init(cfg.Database.Path)
		if err != nil {
			return err
		}

		repo := inventory.NewRepository(db)
		svc := inventory.NewService(repo)

		server := &inventory.Server{
			UUID:     "dummy-uuid-for-v0.0.1",
			Name:     name,
			Host:     host,
			Port:     22,
			Username: "root",
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

		cfg, _ := config.Load()
		db, err := database.Init(cfg.Database.Path)
		if err != nil {
			return err
		}

		repo := inventory.NewRepository(db)
		svc := inventory.NewService(repo)

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

func runSSHConnection(cmd *cobra.Command, args []string) error {
	input := args[0]

	cfg, _ := config.Load()
	db, err := database.Init(cfg.Database.Path)
	if err != nil {
		return err
	}

	repo := inventory.NewRepository(db)
	svc := inventory.NewService(repo)

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

	sshSvc := ssh.NewService(cfg.SSH.Timeout)
	var client ssh.Client

	// Try saved credentials first
	if target.AuthSecret != "" {
		client, err = sshSvc.Connect(cmd.Context(), target.Host, target.Port, target.Username, "password", target.AuthSecret)
	}

	// Prompt if no credentials or connection failed
	if client == nil || err != nil {
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
			if err := svc.AddServer(cmd.Context(), target); err != nil {
				return fmt.Errorf("failed to save server to database: %w", err)
			}
		} else {
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
	serverCmd.AddCommand(serverListCmd, serverAddCmd, serverRemoveCmd)
	rootCmd.AddCommand(versionCmd, configCmd, doctorCmd, serverCmd, sshCmd)
}

func main() {
	if len(os.Args) > 1 {
		cmdName := os.Args[1]
		// Check if it's not a known subcommand and not a flag
		isKnownCommand := false
		knownSubcommands := []string{"completion", "config", "doctor", "help", "server", "version", "ssh"}
		for _, sc := range knownSubcommands {
			if cmdName == sc {
				isKnownCommand = true
				break
			}
		}

		if !isKnownCommand && !strings.HasPrefix(cmdName, "-") {
			// Rewrite os.Args to insert "ssh" at index 1
			newArgs := make([]string, 0, len(os.Args)+1)
			newArgs = append(newArgs, os.Args[0], "ssh")
			newArgs = append(newArgs, os.Args[1:]...)
			os.Args = newArgs
		}
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
