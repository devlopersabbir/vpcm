package main

import (
	"fmt"
	"os"

	"github.com/devlopersabbir/vpcm/internal/config"
	"github.com/devlopersabbir/vpcm/internal/database"
	"github.com/devlopersabbir/vpcm/internal/inventory"
	"github.com/devlopersabbir/vpcm/internal/logger"
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

		if len(servers) == 0 {
			fmt.Println("No servers found. Add one with 'vpsm server add'.")
			return nil
		}

		fmt.Printf("%-5s | %-15s | %-20s | %-5s | %-10s\n", "ID", "Name", "Host", "Port", "Provider")
		fmt.Println("-----------------------------------------------------------------------")
		for _, s := range servers {
			fmt.Printf("%-5d | %-15s | %-20s | %-5d | %-10s\n", s.ID, s.Name, s.Host, s.Port, s.Provider)
		}
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

func init() {
	serverCmd.AddCommand(serverListCmd, serverAddCmd)
	rootCmd.AddCommand(versionCmd, configCmd, doctorCmd, serverCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
