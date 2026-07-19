package cli

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
	"github.com/spf13/cobra"
)

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
		repo, err := inventory.NewSQLiteRepository(db)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to init inventory repo: %w", err)
		}
		return repo, inventory.NewService(repo), nil
	}

	db, err := database.InitMongo(cfg.Database.URI, cfg.Database.Name)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize MongoDB: %w", err)
	}
	repo := inventory.NewMongoRepository(db)
	return repo, inventory.NewService(repo), nil
}

var interactiveList bool

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

		if interactiveList {
			selected, err := runTUI(cmd.Context(), servers)
			if err != nil {
				return err
			}
			if selected == nil {
				fmt.Println("No server selected.")
				return nil
			}
			return runSSHConnection(cmd, []string{selected.Name})
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

var serverRenameCmd = &cobra.Command{
	Use:   "rename [id | name] [new_name]",
	Short: "Rename a server in inventory without changing credentials",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]
		newName := strings.TrimSpace(args[1])
		if newName == "" {
			return fmt.Errorf("new name cannot be empty")
		}

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

		if err := svc.RenameServer(cmd.Context(), targetID, newName); err != nil {
			return err
		}

		fmt.Printf("Successfully renamed server '%s' (ID: %d) to '%s'.\n", input, targetID, newName)
		return nil
	},
}
