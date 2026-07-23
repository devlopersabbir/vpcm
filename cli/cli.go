package cli

import (
	"fmt"
	"os"

	"github.com/devlopersabbir/vpcm/internal/config"
	"github.com/devlopersabbir/vpcm/internal/logger"
	"github.com/spf13/cobra"
)

const Version = "v1.1.0"

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
	serverListCmd.Flags().BoolVarP(&listFavorites, "favorites", "f", false, "Filter to only show favorite servers")
	serverListCmd.Flags().BoolVarP(&listRecents, "recents", "r", false, "List recently connected servers first")

	listCmd.Flags().BoolVarP(&interactiveList, "interactive", "i", false, "Open interactive TUI server explorer")
	listCmd.Flags().BoolVarP(&listFavorites, "favorites", "f", false, "Filter to only show favorite servers")
	listCmd.Flags().BoolVarP(&listRecents, "recents", "r", false, "List recently connected servers first")

	// audit flags
	auditCmd.Flags().StringVarP(&auditFlagName, "name", "n", "", "Server name to audit")
	auditCmd.Flags().StringVarP(&auditFlagHost, "host", "H", "", "Server host IP/address to audit")
	auditCmd.Flags().UintVar(&auditFlagID, "id", 0, "Server database ID to audit")
	_ = auditCmd.RegisterFlagCompletionFunc("name", serverNameCompletions)
	_ = auditCmd.RegisterFlagCompletionFunc("host", serverHostCompletions)
	_ = auditCmd.RegisterFlagCompletionFunc("id", serverIDCompletions)

	serverCmd.AddCommand(serverListCmd, serverAddCmd, serverRemoveCmd, serverFlushCmd, serverRenameCmd, serverExportCmd, serverFavoriteCmd)
	sshCmd.Flags().StringVarP(&identityFile, "identity", "i", "", "identity file (private key)")
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
	rootCmd.AddCommand(versionCmd, configCmd, doctorCmd, serverCmd, sshCmd, listCmd, apiCmd, auditCmd, completionCmd, shellCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		PrintError(err)
		os.Exit(1)
	}
}
