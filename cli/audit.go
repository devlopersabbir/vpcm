package cli

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/devlopersabbir/vpcm/internal/inventory"
	internalssh "github.com/devlopersabbir/vpcm/internal/ssh"
	"github.com/spf13/cobra"
)

// ─── Lipgloss colour palette ──────────────────────────────────────────────────

var (
	auditStyleInfo    = lipgloss.NewStyle().Foreground(lipgloss.Color("#94A3B8"))            // slate-400
	auditStyleSuccess = lipgloss.NewStyle().Foreground(lipgloss.Color("#34D399"))            // emerald-400
	auditStyleSpecs   = lipgloss.NewStyle().Foreground(lipgloss.Color("#A5B4FC"))            // indigo-300
	auditStyleSync    = lipgloss.NewStyle().Foreground(lipgloss.Color("#67E8F9"))            // cyan-300
	auditStyleError   = lipgloss.NewStyle().Foreground(lipgloss.Color("#F87171")).Bold(true) // rose-400
	auditStyleLabel   = lipgloss.NewStyle().Foreground(lipgloss.Color("#22D3EE")).Bold(true) // cyan-400
	auditStyleDim     = lipgloss.NewStyle().Foreground(lipgloss.Color("#475569"))            // slate-600
	auditStyleBox     = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#1E293B")).
				Padding(0, 2).
				MarginTop(1).
				MarginBottom(1)
)

// ─── Flags ────────────────────────────────────────────────────────────────────

var (
	auditFlagName string
	auditFlagHost string
	auditFlagID   uint
)

// ─── Command ──────────────────────────────────────────────────────────────────

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Live SSH inventory audit on a registered server",
	Long: `Run a live SSH inventory audit directly against a server in your inventory.
Each step (OS, hardware, software) is collected over SSH in real-time and
displayed as it completes. Results are saved back to your local inventory.

Examples:
  vpsm audit --name my-vps
  vpsm audit --host 52.54.164.79
  vpsm audit --id 3`,
	RunE: func(cmd *cobra.Command, args []string) error {
		auditFlagName = strings.TrimSpace(auditFlagName)
		auditFlagHost = strings.TrimSpace(auditFlagHost)

		if auditFlagName == "" && auditFlagHost == "" && auditFlagID == 0 {
			return fmt.Errorf("provide at least one of --name, --host, or --id")
		}

		ctx := cmd.Context()

		// ── 1. Resolve server from inventory ─────────────────────────────────
		server, err := resolveAuditTarget(ctx, auditFlagName, auditFlagHost, auditFlagID)
		if err != nil {
			auditPrint(auditStyleError, "error", err.Error())
			return err
		}

		// ── 2. Print header ───────────────────────────────────────────────────
		printAuditHeader(server)

		// ── 3. SSH connect ────────────────────────────────────────────────────
		auditPrint(auditStyleInfo, "info", fmt.Sprintf(
			"Connecting to %s@%s:%d over SSH (%s)...",
			server.Username, server.Host, server.Port, server.AuthType,
		))

		sshSvc := internalssh.NewService(20 * time.Second)
		var client internalssh.Client
		if server.AuthType == "key" {
			client, err = sshSvc.Connect(ctx, server.Host, server.Port, server.Username, "key", server.AuthSecret)
		} else {
			client, err = sshSvc.Connect(ctx, server.Host, server.Port, server.Username, "password", server.AuthSecret)
		}
		if err != nil {
			auditPrint(auditStyleError, "error", "SSH connection failed: "+err.Error())
			return err
		}
		defer client.Close()

		auditPrint(auditStyleSuccess, "success", fmt.Sprintf("Connected to %s@%s", server.Username, server.Host))

		// ── 4. Detect OS ──────────────────────────────────────────────────────
		auditPrint(auditStyleInfo, "info", "Detecting operating system...")
		osFamily, osVersion := inventory.DetectOS(ctx, client)

		osDisplay := strings.Title(osFamily)
		if osVersion != "" {
			osDisplay += " " + osVersion
		}
		auditPrint(auditStyleSuccess, "success", "Discovered "+osDisplay)

		// ── 5. Detect hardware specs ──────────────────────────────────────────
		auditPrint(auditStyleInfo, "info", "Collecting hardware specifications...")
		cpuModel, cpuCores, ramTotal, diskTotal := inventory.DetectSpecs(ctx, client)

		var specParts []string
		if cpuCores > 0 {
			specParts = append(specParts, fmt.Sprintf("%d Cores", cpuCores))
		}
		if ramTotal != "" {
			specParts = append(specParts, ramTotal+" RAM")
		}
		if diskTotal != "" {
			specParts = append(specParts, diskTotal+" Disk")
		}
		auditPrint(auditStyleSpecs, "specs", strings.Join(specParts, " • "))

		// ── 6. Detect server info (network, kernel, uptime…) ─────────────────
		auditPrint(auditStyleInfo, "info", "Fetching network & system info...")
		si := inventory.DetectServerInfo(ctx, client)

		if si.PublicIP != "" {
			auditPrint(auditStyleSpecs, "specs", fmt.Sprintf("Public IP: %s  •  Hostname: %s", si.PublicIP, si.Hostname))
		}
		if si.KernelVersion != "" {
			auditPrint(auditStyleSpecs, "specs", "Kernel: "+si.KernelVersion+"  •  Arch: "+si.Architecture)
		}
		if si.Uptime != "" {
			auditPrint(auditStyleSpecs, "specs", "Uptime: "+si.Uptime)
		}

		// ── 7. Detect installed software ──────────────────────────────────────
		auditPrint(auditStyleInfo, "info", "Scanning installed software & services...")
		swList := inventory.DetectSoftware(ctx, client)

		var swNames []string
		for _, sw := range swList {
			swNames = append(swNames, sw.Name+" "+sw.Version)
		}
		if len(swNames) > 0 {
			auditPrint(auditStyleSpecs, "specs", "Software: "+strings.Join(swNames, " • "))
		} else {
			auditPrint(auditStyleInfo, "info", "No known software detected")
		}

		// ── 8. Detect provider ────────────────────────────────────────────────
		auditPrint(auditStyleInfo, "info", "Identifying cloud provider...")
		provider := inventory.DetectProvider(ctx, client, server.Host)
		auditPrint(auditStyleSuccess, "success", "Provider: "+provider)

		// ── 9. Persist results to local inventory ─────────────────────────────
		auditPrint(auditStyleInfo, "info", "Saving results to local inventory...")

		repo, _, err := initRepoAndService(ctx)
		if err != nil {
			auditPrint(auditStyleError, "error", "Could not open inventory: "+err.Error())
			return err
		}

		if server.Provider == "" || server.Provider == "Generic VPS" {
			server.Provider = provider
			_ = repo.Update(ctx, server)
		}

		_ = repo.ReplaceSoftware(ctx, server.ID, swList)

		_ = repo.UpsertNetwork(ctx, &inventory.ServerNetwork{
			ServerID:         server.ID,
			Hostname:         si.Hostname,
			PublicIP:         si.PublicIP,
			PrivateIP:        si.PrivateIP,
			MACAddress:       si.MACAddress,
			Region:           si.Region,
			AvailabilityZone: si.AvailabilityZone,
		})

		_ = repo.UpsertHardware(ctx, &inventory.ServerHardware{
			ServerID:       server.ID,
			CPUModel:       cpuModel,
			CPUCores:       cpuCores,
			RAMTotal:       ramTotal,
			SwapTotal:      si.SwapTotal,
			DiskTotal:      diskTotal,
			Virtualization: si.Virtualization,
			InstanceType:   si.InstanceType,
			SerialNumber:   si.SerialNumber,
			BIOSVersion:    si.BIOSVersion,
			Uptime:         si.Uptime,
		})

		_ = repo.UpsertOS(ctx, &inventory.ServerOS{
			ServerID:       server.ID,
			OSFamily:       osFamily,
			OSVersion:      osVersion,
			KernelVersion:  si.KernelVersion,
			Architecture:   si.Architecture,
			InitSystem:     si.InitSystem,
			Timezone:       si.Timezone,
			Locale:         si.Locale,
			PackageManager: si.PackageManager,
		})

		auditPrint(auditStyleSync, "sync", "Inventory updated in local database")

		// ── 10. Summary footer ────────────────────────────────────────────────
		printAuditFooter(server, osDisplay, specParts, swNames, si, provider)
		return nil
	},
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func resolveAuditTarget(ctx context.Context, name, host string, id uint) (*inventory.Server, error) {
	repo, _, err := initRepoAndService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to inventory: %w", err)
	}

	// --id: direct lookup by primary key
	if id != 0 {
		s, err := repo.GetByID(ctx, id)
		if err != nil || s == nil {
			return nil, fmt.Errorf("no server found with id %s", strconv.Itoa(int(id)))
		}
		return s, nil
	}

	servers, err := repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list servers: %w", err)
	}
	for _, s := range servers {
		if name != "" && strings.EqualFold(s.Name, name) {
			return &s, nil
		}
		if host != "" && strings.EqualFold(s.Host, host) {
			return &s, nil
		}
	}
	if name != "" {
		return nil, fmt.Errorf("no server found with name %q", name)
	}
	return nil, fmt.Errorf("no server found with host %q", host)
}

func auditPrint(style lipgloss.Style, tag, msg string) {
	bracket := auditStyleDim.Render("[")
	tagStr := style.Copy().UnsetBold().Bold(tag == "error").Render(tag)
	closeBracket := auditStyleDim.Render("]")
	line := fmt.Sprintf("%s%s%s %s", bracket, tagStr, closeBracket, lipgloss.NewStyle().Foreground(style.GetForeground()).Render(msg))
	fmt.Fprintln(os.Stdout, line)
}

func printAuditHeader(s *inventory.Server) {
	re := lipgloss.NewRenderer(os.Stdout)
	title := re.NewStyle().Bold(true).Foreground(lipgloss.Color("#22D3EE")).
		Render(fmt.Sprintf("⚡  VPSM Live Audit — %s", s.Name))
	sub := re.NewStyle().Foreground(lipgloss.Color("#64748B")).
		Render(fmt.Sprintf("   %s@%s:%d", s.Username, s.Host, s.Port))
	fmt.Fprintln(os.Stdout, auditStyleBox.Render(title+"\n"+sub))
}

func printAuditFooter(
	s *inventory.Server,
	osDisplay string,
	specParts []string,
	swNames []string,
	si inventory.ServerInfo,
	provider string,
) {
	_ = s
	lines := []string{
		auditStyleLabel.Render("  Audit Complete ✓"),
		"",
	}

	add := func(k, v string) {
		if v != "" {
			lines = append(lines, fmt.Sprintf("  %-16s %s",
				auditStyleDim.Render(k+":"),
				auditStyleInfo.Render(v),
			))
		}
	}

	add("OS", osDisplay)
	add("Provider", provider)
	if len(specParts) > 0 {
		add("Hardware", strings.Join(specParts, " • "))
	}
	add("Public IP", si.PublicIP)
	add("Hostname", si.Hostname)
	add("Kernel", si.KernelVersion)
	add("Arch", si.Architecture)
	add("Uptime", si.Uptime)
	if len(swNames) > 0 {
		add("Software", strings.Join(swNames, ", "))
	}

	fmt.Fprintln(os.Stdout, auditStyleBox.Render(strings.Join(lines, "\n")))
}
