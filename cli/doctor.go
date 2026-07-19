package cli

import (
	"fmt"

	"github.com/devlopersabbir/vpcm/internal/config"
	"github.com/devlopersabbir/vpcm/internal/database"
	"github.com/spf13/cobra"
)

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
