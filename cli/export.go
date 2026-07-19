package cli

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var exportFormat string
var exportOutputFile string

var serverExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export all servers in various formats",
	Long: `Export your entire server inventory database into machine-readable formats.
Supported formats include JSON (default), YAML, CSV sheets, and standard SSH config files (~/.ssh/config format).

If the --out / -o flag is provided, the formatted output will be saved directly to the specified file path.
Otherwise, it prints directly to standard output (stdout), allowing you to pipe the data to other commands.`,
	Example: `  # Export all servers as JSON (printed to stdout)
  vpsm server export --format json
  vpsm server export -f json

  # Export all servers as YAML and save to a backup file
  vpsm server export -f yaml --out backups/servers.yaml
  vpsm server export -f yaml -o backups/servers.yaml

  # Export to CSV for spreadsheet imports
  vpsm server export -f csv -o ~/Documents/servers.csv

  # Export to SSH config format to append to ~/.ssh/config
  vpsm server export -f ssh >> ~/.ssh/config`,
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, _, err := initRepoAndService(cmd.Context())
		if err != nil {
			return err
		}

		servers, err := repo.List(cmd.Context())
		if err != nil {
			return err
		}

		var outputBytes []byte

		format := strings.ToLower(strings.TrimSpace(exportFormat))
		switch format {
		case "ssh":
			var sb strings.Builder
			for _, s := range servers {
				sb.WriteString(fmt.Sprintf("Host %s\n", s.Name))
				sb.WriteString(fmt.Sprintf("    HostName %s\n", s.Host))
				sb.WriteString(fmt.Sprintf("    User %s\n", s.Username))
				sb.WriteString(fmt.Sprintf("    Port %d\n", s.Port))
				if s.AuthType == "key" {
					sb.WriteString(fmt.Sprintf("    # AuthType: key (key stored in database)\n"))
				}
				sb.WriteString("\n")
			}
			outputBytes = []byte(sb.String())

		case "json":
			bytes, err := json.MarshalIndent(servers, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal to JSON: %w", err)
			}
			outputBytes = bytes

		case "yaml", "yml":
			bytes, err := yaml.Marshal(servers)
			if err != nil {
				return fmt.Errorf("failed to marshal to YAML: %w", err)
			}
			outputBytes = bytes

		case "csv":
			var sb strings.Builder
			writer := csv.NewWriter(&sb)
			err := writer.Write([]string{"ID", "UUID", "Name", "Host", "Port", "Username", "AuthType", "Provider"})
			if err != nil {
				return fmt.Errorf("failed to write CSV header: %w", err)
			}
			for _, s := range servers {
				err = writer.Write([]string{
					strconv.Itoa(int(s.ID)),
					s.UUID,
					s.Name,
					s.Host,
					strconv.Itoa(s.Port),
					s.Username,
					s.AuthType,
					s.Provider,
				})
				if err != nil {
					return fmt.Errorf("failed to write CSV row: %w", err)
				}
			}
			writer.Flush()
			outputBytes = []byte(sb.String())

		default:
			return fmt.Errorf("unsupported format '%s'; choose from ssh, json, csv, yaml", exportFormat)
		}

		if exportOutputFile != "" {
			err = os.WriteFile(exportOutputFile, outputBytes, 0600)
			if err != nil {
				return fmt.Errorf("failed to write output to file %s: %w", exportOutputFile, err)
			}
			fmt.Printf("Successfully exported servers to %s (%s format).\n", exportOutputFile, format)
		} else {
			fmt.Println(string(outputBytes))
		}

		return nil
	},
}
