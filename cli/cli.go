package cli

import (
	"fmt"
	"os"

	"github.com/devlopersabbir/vpcm/internal/config"
	"github.com/devlopersabbir/vpcm/internal/logger"
	"github.com/devlopersabbir/vpcm/internal/version"
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
		fmt.Println(version.Version)
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

	serverImportCmd.Flags().StringVarP(&importFormat, "format", "f", "auto", "Input format (ssh, json, csv, yaml, auto)")
	serverImportCmd.Flags().StringVarP(&importInputFile, "in", "i", "", "Input file path (default: stdin)")
	serverImportCmd.Flags().StringVar(&importOnConflict, "on-conflict", "skip", "Action for servers that already exist (skip, overwrite, rename, fail)")
	serverImportCmd.Flags().BoolVar(&importDryRun, "dry-run", false, "Preview the import without writing to the database")
	_ = serverImportCmd.RegisterFlagCompletionFunc("format", func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		return []string{"auto", "json", "yaml", "csv", "ssh"}, cobra.ShellCompDirectiveNoFileComp
	})
	_ = serverImportCmd.RegisterFlagCompletionFunc("on-conflict", func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		return []string{"skip", "overwrite", "rename", "fail"}, cobra.ShellCompDirectiveNoFileComp
	})

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

	serverCmd.AddCommand(serverListCmd, serverAddCmd, serverRemoveCmd, serverFlushCmd, serverRenameCmd, serverExportCmd, serverImportCmd, serverFavoriteCmd)
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
