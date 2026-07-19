package cli

import (
	"fmt"
	"os"

	"github.com/devlopersabbir/vpcm/internal/config"
	"github.com/devlopersabbir/vpcm/internal/logger"
	"github.com/spf13/cobra"
)

const Version = "v0.1.7"

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
		fmt.Println(Version)
	},
}

func init() {
	if configExists() {
		configCmd.AddCommand(configShowCmd, configEditCmd, configReloadCmd)
	} else {
		configCmd.AddCommand(configShowCmd, configInitCmd, configReloadCmd)
	}
	serverExportCmd.Flags().StringVarP(&exportFormat, "format", "f", "json", "Output format (ssh, json, csv, yaml)")
	serverExportCmd.Flags().StringVarP(&exportOutputFile, "out", "o", "", "Output file path (default: stdout)")

	serverListCmd.Flags().BoolVarP(&interactiveList, "interactive", "i", false, "Open interactive TUI server explorer")
	listCmd.Flags().BoolVarP(&interactiveList, "interactive", "i", false, "Open interactive TUI server explorer")

	serverCmd.AddCommand(serverListCmd, serverAddCmd, serverRemoveCmd, serverFlushCmd, serverRenameCmd, serverExportCmd)
	sshCmd.Flags().StringVarP(&identityFile, "identity", "i", "", "identity file (private key)")
	rootCmd.AddCommand(versionCmd, configCmd, doctorCmd, serverCmd, sshCmd, listCmd, apiCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
