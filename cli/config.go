package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/devlopersabbir/vpcm/internal/config"
	"github.com/spf13/cobra"
)

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
		rows = append(rows, []string{"API Mode", cfg.API.Mode})
		rows = append(rows, []string{"API Host", cfg.API.Host})
		rows = append(rows, []string{"API Port", strconv.Itoa(cfg.API.Port)})
		rows = append(rows, []string{"Global SaaS URL", cfg.API.GlobalURL})
		if cfg.API.Token != "" {
			rows = append(rows, []string{"Global SaaS Token", "********"})
		} else {
			rows = append(rows, []string{"Global SaaS Token", "None"})
		}
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
				if val == "Enabled" || val == "local" || val == "cloud" {
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
				defaultURI = "mongodb://127.0.0.1:27017"
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

		fmt.Print("Enable REST API server? [Y/n] (default: Y): ")
		var enableAPI string
		_, _ = fmt.Scanln(&enableAPI)
		enableAPI = strings.ToLower(strings.TrimSpace(enableAPI))
		if enableAPI == "" || enableAPI == "y" || enableAPI == "yes" {
			cfg.API.Enabled = true
		} else {
			cfg.API.Enabled = false
		}

		if cfg.API.Enabled {
			fmt.Printf("Choose API server mode [local/cloud] (default: local): ")
			var mode string
			_, _ = fmt.Scanln(&mode)
			mode = strings.ToLower(strings.TrimSpace(mode))
			if mode == "" {
				mode = "local"
			}
			if mode != "local" && mode != "cloud" {
				return fmt.Errorf("invalid API mode choice: must be 'local' or 'cloud'")
			}
			cfg.API.Mode = mode

			if mode == "cloud" {
				fmt.Print("Enter Global SaaS Authentication Token (optional): ")
				var token string
				_, _ = fmt.Scanln(&token)
				cfg.API.Token = strings.TrimSpace(token)
			}
		}

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		fmt.Println("\n✨ Configuration saved successfully to ~/.config/vpsm/config.yaml!")
		return nil
	},
}

func configExists() bool {
	configHome := filepath.Join(os.Getenv("HOME"), ".config", "vpsm")
	configPath := filepath.Join(configHome, "config.yaml")
	_, err := os.Stat(configPath)
	return err == nil
}

var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Interactively edit your existing universe settings",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.Load()

		fmt.Println("📝 Editing VPSM Configuration...")

		fmt.Printf("Choose database driver [sqlite/mongodb] (current: %s): ", cfg.Database.Driver)
		var driver string
		_, _ = fmt.Scanln(&driver)
		driver = strings.ToLower(strings.TrimSpace(driver))
		if driver != "" {
			if driver != "sqlite" && driver != "mongodb" {
				return fmt.Errorf("invalid driver choice: must be 'sqlite' or 'mongodb'")
			}
			cfg.Database.Driver = driver
		}

		if cfg.Database.Driver == "sqlite" {
			fmt.Printf("Enter SQLite database path (current: %s): ", cfg.Database.Path)
			var path string
			_, _ = fmt.Scanln(&path)
			path = strings.TrimSpace(path)
			if path != "" {
				cfg.Database.Path = path
			}
		} else {
			fmt.Printf("Enter MongoDB Connection URI (current: %s): ", cfg.Database.URI)
			var uri string
			_, _ = fmt.Scanln(&uri)
			uri = strings.TrimSpace(uri)
			if uri != "" {
				cfg.Database.URI = uri
			}

			fmt.Printf("Enter MongoDB Database Name (current: %s): ", cfg.Database.Name)
			var name string
			_, _ = fmt.Scanln(&name)
			name = strings.TrimSpace(name)
			if name != "" {
				cfg.Database.Name = name
			}
		}

		apiStatus := "n"
		if cfg.API.Enabled {
			apiStatus = "y"
		}
		fmt.Printf("Enable REST API server? [y/N] (current: %s): ", apiStatus)
		var enableAPI string
		_, _ = fmt.Scanln(&enableAPI)
		enableAPI = strings.ToLower(strings.TrimSpace(enableAPI))
		if enableAPI != "" {
			if enableAPI == "y" || enableAPI == "yes" {
				cfg.API.Enabled = true
			} else {
				cfg.API.Enabled = false
			}
		}

		if cfg.API.Enabled {
			fmt.Printf("Choose API server mode [local/cloud] (current: %s): ", cfg.API.Mode)
			var mode string
			_, _ = fmt.Scanln(&mode)
			mode = strings.ToLower(strings.TrimSpace(mode))
			if mode != "" {
				if mode != "local" && mode != "cloud" {
					return fmt.Errorf("invalid API mode choice: must be 'local' or 'cloud'")
				}
				cfg.API.Mode = mode
			}

			if cfg.API.Mode == "cloud" {
				tokenPlaceholder := "None"
				if cfg.API.Token != "" {
					tokenPlaceholder = "********"
				}
				fmt.Printf("Enter Global SaaS Authentication Token (current: %s): ", tokenPlaceholder)
				var token string
				_, _ = fmt.Scanln(&token)
				token = strings.TrimSpace(token)
				if token != "" {
					cfg.API.Token = token
				}
			}
		}

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		fmt.Println("\n✨ Configuration updated successfully!")
		return nil
	},
}

var configReloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "Format, validate, and reload the active settings",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔄 Reloading VPSM Configuration...")

		// Load and validate configuration
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to parse configuration file: %w", err)
		}

		// Validation Checks
		if cfg.Database.Driver != "sqlite" && cfg.Database.Driver != "mongodb" {
			return fmt.Errorf("validation failed: database driver must be 'sqlite' or 'mongodb'")
		}
		if cfg.API.Port < 1 || cfg.API.Port > 65535 {
			return fmt.Errorf("validation failed: API port %d out of bounds (1-65535)", cfg.API.Port)
		}
		if cfg.API.Enabled && cfg.API.Mode != "local" && cfg.API.Mode != "cloud" {
			return fmt.Errorf("validation failed: API mode must be 'local' or 'cloud'")
		}

		// Check directory write permissions for SQLite
		if cfg.Database.Driver == "sqlite" {
			dir := filepath.Dir(cfg.Database.Path)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("validation failed: sqlite database folder is not writeable: %w", err)
			}
		}

		// Format / Save configuration (Viper cleans up layout/formatting)
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to format and save configuration: %w", err)
		}
		fmt.Println("[✓] Config formatted and validated successfully.")

		// Restart local API daemon if it is running in background
		if isAPIRunning() {
			fmt.Println("API server daemon detected in background. Restarting daemon...")
			// Simulate calling apiRestartCmd
			if err := apiRestartCmd.RunE(cmd, nil); err != nil {
				return fmt.Errorf("failed to restart API daemon: %w", err)
			}
			fmt.Println("[✓] API server daemon reloaded successfully.")
		} else {
			fmt.Println("No active background API daemon found. New config will take effect on next start.")
		}

		return nil
	},
}
