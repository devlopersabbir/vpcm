package cli

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/devlopersabbir/vpcm/internal/inventory"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var importFormat string
var importInputFile string
var importOnConflict string
var importDryRun bool

var serverImportCmd = &cobra.Command{
	Use:   "import [file]",
	Short: "Import servers from a previously exported file",
	Args:  cobra.MaximumNArgs(1),
	Long: `Import servers into your inventory database from a file produced by 'vpsm server export'.
Supported formats are JSON, YAML, CSV and standard SSH config files (~/.ssh/config format).

The format is auto-detected from the file extension and contents, so --format / -f is only
needed when reading from stdin or when the extension is misleading.

The input file is read from the path given as an argument or via --in / -i.
If neither is supplied, the data is read from standard input (stdin), so exports can be piped in.

Servers already present in the database are matched by UUID first and then by name.
Use --on-conflict to decide what happens to those matches, and --dry-run to preview
the outcome without writing anything.`,
	Example: `  # Import servers from a JSON backup
  vpsm server import backups/servers.json
  vpsm server import -f json -i backups/servers.json

  # Preview what a YAML backup would change, without touching the database
  vpsm server import backups/servers.yaml --dry-run

  # Update servers that already exist instead of skipping them
  vpsm server import backups/servers.json --on-conflict overwrite

  # Keep both copies by importing conflicting servers under a suffixed name
  vpsm server import backups/servers.csv --on-conflict rename

  # Pipe an export straight from another machine
  ssh admin@backup-box 'vpsm server export -f json' | vpsm server import -f json

  # Adopt hosts from an existing SSH config
  vpsm server import ~/.ssh/config -f ssh`,
	RunE: func(cmd *cobra.Command, args []string) error {
		conflictMode := strings.ToLower(strings.TrimSpace(importOnConflict))
		switch conflictMode {
		case "skip", "overwrite", "rename", "fail":
		default:
			return fmt.Errorf("unsupported conflict strategy '%s'; choose from skip, overwrite, rename, fail", importOnConflict)
		}

		path := strings.TrimSpace(importInputFile)
		if len(args) > 0 {
			if path != "" {
				return fmt.Errorf("input file given twice; pass it either as an argument or via --in / -i, not both")
			}
			path = strings.TrimSpace(args[0])
		}

		data, err := readImportSource(path)
		if err != nil {
			return err
		}
		if len(bytes.TrimSpace(data)) == 0 {
			return fmt.Errorf("no data to import: %s is empty", describeImportSource(path))
		}

		format, err := resolveImportFormat(importFormat, path, data)
		if err != nil {
			return err
		}

		incoming, err := parseImportData(format, data)
		if err != nil {
			return err
		}
		if len(incoming) == 0 {
			return fmt.Errorf("no servers found in %s (parsed as %s)", describeImportSource(path), format)
		}

		repo, _, err := initRepoAndService(cmd.Context())
		if err != nil {
			return err
		}

		existing, err := repo.List(cmd.Context())
		if err != nil {
			return err
		}

		byUUID := make(map[string]inventory.Server, len(existing))
		byName := make(map[string]inventory.Server, len(existing))
		takenNames := make(map[string]bool, len(existing))
		takenUUIDs := make(map[string]bool, len(existing))
		for _, s := range existing {
			byUUID[s.UUID] = s
			byName[strings.ToLower(s.Name)] = s
			takenNames[strings.ToLower(s.Name)] = true
			takenUUIDs[s.UUID] = true
		}

		var results []importResult
		var created, updated, skipped, invalid int

		// A write failure part-way through leaves earlier records applied, so the
		// summary is always printed before the error is surfaced.
		var importErr error

	records:
		for _, candidate := range incoming {
			server, err := normalizeImportedServer(candidate)
			if err != nil {
				results = append(results, importResult{Name: candidate.Name, Host: candidate.Host, Action: "invalid", Detail: err.Error()})
				invalid++
				continue
			}

			match, matchedOn := findImportConflict(server, byUUID, byName)

			if match == nil {
				if server.UUID == "" || takenUUIDs[server.UUID] {
					server.UUID = uuid.NewString()
				}
				if !importDryRun {
					if err := repo.Create(cmd.Context(), &server); err != nil {
						importErr = fmt.Errorf("failed to import server '%s': %w", server.Name, err)
						break records
					}
				}
				takenNames[strings.ToLower(server.Name)] = true
				takenUUIDs[server.UUID] = true
				byName[strings.ToLower(server.Name)] = server
				byUUID[server.UUID] = server
				results = append(results, importResult{Name: server.Name, Host: server.Host, Action: "created"})
				created++
				continue
			}

			switch conflictMode {
			case "fail":
				importErr = fmt.Errorf("server '%s' already exists (matched by %s); rerun with --on-conflict skip, overwrite or rename", server.Name, matchedOn)
				break records

			case "skip":
				results = append(results, importResult{Name: server.Name, Host: server.Host, Action: "skipped", Detail: "already exists, matched by " + matchedOn})
				skipped++

			case "overwrite":
				merged := mergeImportedServer(*match, server)
				if !importDryRun {
					if err := repo.Update(cmd.Context(), &merged); err != nil {
						importErr = fmt.Errorf("failed to update server '%s': %w", merged.Name, err)
						break records
					}
				}
				// A rename frees the old key, which must not keep matching later records.
				if !strings.EqualFold(match.Name, merged.Name) {
					delete(byName, strings.ToLower(match.Name))
				}
				byUUID[merged.UUID] = merged
				byName[strings.ToLower(merged.Name)] = merged
				takenNames[strings.ToLower(merged.Name)] = true
				results = append(results, importResult{Name: merged.Name, Host: merged.Host, Action: "updated", Detail: "matched by " + matchedOn})
				updated++

			case "rename":
				server.Name = uniqueImportName(server.Name, takenNames)
				server.UUID = uuid.NewString()
				if !importDryRun {
					if err := repo.Create(cmd.Context(), &server); err != nil {
						importErr = fmt.Errorf("failed to import server '%s': %w", server.Name, err)
						break records
					}
				}
				takenNames[strings.ToLower(server.Name)] = true
				takenUUIDs[server.UUID] = true
				byName[strings.ToLower(server.Name)] = server
				byUUID[server.UUID] = server
				results = append(results, importResult{Name: server.Name, Host: server.Host, Action: "created", Detail: "renamed, conflicted by " + matchedOn})
				created++
			}
		}

		printImportResults(results, format, path, importSummary{
			Created:   created,
			Updated:   updated,
			Skipped:   skipped,
			Invalid:   invalid,
			DryRun:    importDryRun,
			Aborted:   importErr != nil,
			Remaining: len(incoming) - (created + updated + skipped + invalid),
		})
		return importErr
	},
}

// ─── Source & format handling ─────────────────────────────────────────────────

func readImportSource(path string) ([]byte, error) {
	if path == "" {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, fmt.Errorf("failed to read import data from stdin: %w", err)
		}
		return data, nil
	}

	expanded, err := expandImportPath(path)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(expanded)
	if err != nil {
		return nil, fmt.Errorf("failed to read import file %s: %w", path, err)
	}
	return data, nil
}

func expandImportPath(path string) (string, error) {
	if path != "~" && !strings.HasPrefix(path, "~/") {
		return path, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to resolve home directory for path %s: %w", path, err)
	}
	if path == "~" {
		return home, nil
	}
	return filepath.Join(home, path[2:]), nil
}

func describeImportSource(path string) string {
	if path == "" {
		return "stdin"
	}
	return path
}

// resolveImportFormat honours an explicit --format, then falls back to the file
// extension and finally to sniffing the payload itself.
func resolveImportFormat(requested, path string, data []byte) (string, error) {
	switch format := strings.ToLower(strings.TrimSpace(requested)); format {
	case "", "auto":
	case "json":
		return "json", nil
	case "yaml", "yml":
		return "yaml", nil
	case "csv":
		return "csv", nil
	case "ssh":
		return "ssh", nil
	default:
		return "", fmt.Errorf("unsupported format '%s'; choose from ssh, json, csv, yaml (or auto)", requested)
	}

	switch strings.ToLower(filepath.Ext(path)) {
	case ".json":
		return "json", nil
	case ".yaml", ".yml":
		return "yaml", nil
	case ".csv":
		return "csv", nil
	case ".conf", ".cfg", ".sshconfig":
		return "ssh", nil
	}

	if base := strings.ToLower(filepath.Base(path)); base == "config" || strings.HasPrefix(base, "ssh_config") {
		return "ssh", nil
	}

	return sniffImportFormat(data)
}

func sniffImportFormat(data []byte) (string, error) {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return "", fmt.Errorf("cannot detect format of empty input; pass --format / -f explicitly")
	}

	if trimmed[0] == '[' || trimmed[0] == '{' {
		return "json", nil
	}

	firstLine := trimmed
	if idx := bytes.IndexByte(trimmed, '\n'); idx >= 0 {
		firstLine = trimmed[:idx]
	}
	lowerFirst := strings.ToLower(strings.TrimSpace(string(firstLine)))

	if strings.HasPrefix(lowerFirst, "host ") || strings.HasPrefix(lowerFirst, "host\t") {
		return "ssh", nil
	}
	if strings.Contains(lowerFirst, ",") && strings.Contains(lowerFirst, "name") && strings.Contains(lowerFirst, "host") {
		return "csv", nil
	}
	if strings.HasPrefix(lowerFirst, "- ") || strings.Contains(lowerFirst, ":") {
		return "yaml", nil
	}

	return "", fmt.Errorf("could not detect the input format; pass --format / -f explicitly (ssh, json, csv, yaml)")
}

// ─── Parsers ──────────────────────────────────────────────────────────────────

func parseImportData(format string, data []byte) ([]inventory.Server, error) {
	switch format {
	case "json":
		return parseJSONServers(data)
	case "yaml":
		return parseYAMLServers(data)
	case "csv":
		return parseCSVServers(data)
	case "ssh":
		return parseSSHConfigServers(data)
	default:
		return nil, fmt.Errorf("unsupported format '%s'; choose from ssh, json, csv, yaml", format)
	}
}

func parseJSONServers(data []byte) ([]inventory.Server, error) {
	var servers []inventory.Server
	if err := json.Unmarshal(data, &servers); err == nil {
		return servers, nil
	}

	var single inventory.Server
	if err := json.Unmarshal(data, &single); err != nil {
		return nil, fmt.Errorf("failed to parse JSON input: %w", err)
	}
	return []inventory.Server{single}, nil
}

func parseYAMLServers(data []byte) ([]inventory.Server, error) {
	var servers []inventory.Server
	if err := yaml.Unmarshal(data, &servers); err == nil {
		return servers, nil
	}

	var single inventory.Server
	if err := yaml.Unmarshal(data, &single); err != nil {
		return nil, fmt.Errorf("failed to parse YAML input: %w", err)
	}
	return []inventory.Server{single}, nil
}

func parseCSVServers(data []byte) ([]inventory.Server, error) {
	reader := csv.NewReader(bytes.NewReader(data))
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV input: %w", err)
	}
	if len(records) == 0 {
		return nil, fmt.Errorf("CSV input contains no rows")
	}

	columns := make(map[string]int, len(records[0]))
	for i, header := range records[0] {
		key := strings.ToLower(strings.TrimSpace(header))
		key = strings.NewReplacer(" ", "", "_", "", "-", "").Replace(key)
		columns[key] = i
	}
	if _, ok := columns["name"]; !ok {
		return nil, fmt.Errorf("CSV input is missing a 'Name' column header")
	}
	if _, ok := columns["host"]; !ok {
		return nil, fmt.Errorf("CSV input is missing a 'Host' column header")
	}

	field := func(row []string, names ...string) string {
		for _, name := range names {
			if idx, ok := columns[name]; ok && idx < len(row) {
				return strings.TrimSpace(row[idx])
			}
		}
		return ""
	}

	var servers []inventory.Server
	for _, row := range records[1:] {
		if len(row) == 0 || strings.TrimSpace(strings.Join(row, "")) == "" {
			continue
		}

		server := inventory.Server{
			UUID:       field(row, "uuid"),
			Name:       field(row, "name"),
			Host:       field(row, "host"),
			Username:   field(row, "username", "user"),
			AuthType:   field(row, "authtype"),
			AuthSecret: field(row, "authsecret"),
			Provider:   field(row, "provider"),
		}

		if raw := field(row, "port"); raw != "" {
			port, err := strconv.Atoi(raw)
			if err != nil {
				return nil, fmt.Errorf("invalid port '%s' for server '%s' in CSV input", raw, server.Name)
			}
			server.Port = port
		}

		switch strings.ToLower(field(row, "isfavorite", "favorite")) {
		case "1", "true", "yes", "y":
			server.IsFavorite = true
		}

		for _, tag := range strings.FieldsFunc(field(row, "tags"), func(r rune) bool { return r == ';' || r == '|' }) {
			if tag = strings.TrimSpace(tag); tag != "" {
				server.Tags = append(server.Tags, inventory.Tag{Name: tag})
			}
		}

		servers = append(servers, server)
	}

	return servers, nil
}

func parseSSHConfigServers(data []byte) ([]inventory.Server, error) {
	var servers []inventory.Server
	var current *inventory.Server

	flush := func() {
		if current != nil && current.Name != "" {
			servers = append(servers, *current)
		}
		current = nil
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		keyword, value := splitSSHConfigLine(line)
		if value == "" {
			continue
		}

		switch strings.ToLower(keyword) {
		case "host":
			flush()
			// Patterns cannot be resolved to a single host, so they are not importable.
			if strings.ContainsAny(value, "*?!") {
				continue
			}
			// Only the first alias of a multi-alias Host line becomes the server name.
			alias := strings.Fields(value)[0]
			current = &inventory.Server{Name: alias, Host: alias}
		case "hostname":
			if current != nil {
				current.Host = value
			}
		case "user":
			if current != nil {
				current.Username = value
			}
		case "port":
			if current == nil {
				continue
			}
			port, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("invalid Port '%s' for Host '%s' in SSH config input", value, current.Name)
			}
			current.Port = port
		case "identityfile":
			if current != nil {
				current.AuthType = "key"
				current.AuthSecret = value
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to parse SSH config input: %w", err)
	}
	flush()

	return servers, nil
}

func splitSSHConfigLine(line string) (string, string) {
	if idx := strings.IndexAny(line, " \t="); idx >= 0 {
		keyword := line[:idx]
		value := strings.TrimSpace(strings.TrimLeft(strings.TrimSpace(line[idx+1:]), "="))
		return keyword, value
	}
	return line, ""
}

// ─── Record handling ──────────────────────────────────────────────────────────

func normalizeImportedServer(server inventory.Server) (inventory.Server, error) {
	server.ID = 0
	server.UUID = strings.TrimSpace(server.UUID)
	server.Name = strings.TrimSpace(server.Name)
	server.Host = strings.TrimSpace(server.Host)
	server.Username = strings.TrimSpace(server.Username)
	server.AuthType = strings.ToLower(strings.TrimSpace(server.AuthType))
	server.Provider = strings.TrimSpace(server.Provider)

	if server.Host == "" {
		return server, fmt.Errorf("missing host")
	}
	if server.Name == "" {
		server.Name = server.Host
	}
	if server.Port == 0 {
		server.Port = 22
	}
	if server.Port < 1 || server.Port > 65535 {
		return server, fmt.Errorf("port %d is out of range", server.Port)
	}
	if server.Username == "" {
		server.Username = "root"
	}
	switch server.AuthType {
	case "", "key", "keyfile", "password":
	default:
		return server, fmt.Errorf("unknown auth type '%s'", server.AuthType)
	}

	return server, nil
}

func findImportConflict(server inventory.Server, byUUID, byName map[string]inventory.Server) (*inventory.Server, string) {
	if server.UUID != "" {
		if match, ok := byUUID[server.UUID]; ok {
			return &match, "uuid"
		}
	}
	if match, ok := byName[strings.ToLower(server.Name)]; ok {
		return &match, "name"
	}
	return nil, ""
}

// mergeImportedServer layers the imported values onto the stored record, keeping
// the existing identity and never clearing credentials the import does not carry.
func mergeImportedServer(existing, incoming inventory.Server) inventory.Server {
	merged := existing
	merged.Name = incoming.Name
	merged.Host = incoming.Host
	merged.Port = incoming.Port
	merged.Username = incoming.Username
	merged.IsFavorite = incoming.IsFavorite

	if incoming.AuthType != "" {
		merged.AuthType = incoming.AuthType
	}
	if incoming.AuthSecret != "" {
		merged.AuthSecret = incoming.AuthSecret
	}
	if incoming.Provider != "" {
		merged.Provider = incoming.Provider
	}
	if len(incoming.Tags) > 0 {
		merged.Tags = incoming.Tags
	}
	if incoming.LastSeen != nil {
		merged.LastSeen = incoming.LastSeen
	}

	return merged
}

func uniqueImportName(base string, taken map[string]bool) string {
	if !taken[strings.ToLower(base)] {
		return base
	}
	for i := 1; ; i++ {
		candidate := fmt.Sprintf("%s-%d", base, i)
		if !taken[strings.ToLower(candidate)] {
			return candidate
		}
	}
}

// ─── Reporting ────────────────────────────────────────────────────────────────

type importResult struct {
	Name   string
	Host   string
	Action string
	Detail string
}

type importSummary struct {
	Created   int
	Updated   int
	Skipped   int
	Invalid   int
	DryRun    bool
	Aborted   bool
	Remaining int
}

func printImportResults(results []importResult, format, path string, summary importSummary) {
	re := lipgloss.NewRenderer(os.Stdout)
	purple := lipgloss.Color("#7D56F4")
	gray := lipgloss.Color("#3C3C3C")
	green := lipgloss.Color("#4ADE80")
	amber := lipgloss.Color("#FBBF24")
	red := lipgloss.Color("#F87171")

	if summary.DryRun {
		fmt.Println(re.NewStyle().Bold(true).Foreground(amber).Render(" 🔍  DRY RUN — no changes were written"))
	}

	var rows [][]string
	for _, r := range results {
		rows = append(rows, []string{r.Name, r.Host, r.Action, r.Detail})
	}

	tbl := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(re.NewStyle().Foreground(gray)).
		Headers("Name", "Host", "Action", "Details").
		Rows(rows...)

	tbl.StyleFunc(func(row, col int) lipgloss.Style {
		if row < 0 {
			return re.NewStyle().Bold(true).Foreground(purple).Padding(0, 1)
		}
		style := re.NewStyle().Padding(0, 1)
		switch results[row].Action {
		case "created":
			return style.Foreground(green)
		case "updated":
			return style.Foreground(purple)
		case "invalid":
			return style.Foreground(red)
		default:
			return style.Foreground(amber)
		}
	})

	fmt.Println(tbl)

	verb := "Imported"
	switch {
	case summary.DryRun:
		verb = "Would import"
	case summary.Aborted:
		verb = "Applied before stopping"
	}
	line := fmt.Sprintf("%s from %s (%s format): %d created, %d updated, %d skipped",
		verb, describeImportSource(path), format, summary.Created, summary.Updated, summary.Skipped)
	if summary.Invalid > 0 {
		line += fmt.Sprintf(", %s", re.NewStyle().Foreground(red).Render(fmt.Sprintf("%d invalid", summary.Invalid)))
	}
	fmt.Println(line + ".")

	if summary.Aborted && summary.Remaining > 0 {
		fmt.Println(re.NewStyle().Foreground(amber).Render(fmt.Sprintf(
			"%d record(s) were not processed. Re-running the import will pick up where it stopped.", summary.Remaining)))
	}
}
