package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	errBadgeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#DC2626")). // Red background
			Padding(0, 1)

	errHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#F87171")) // Soft Red / Rose

	errPathStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FDE047")) // Yellow highlight for key paths

	errHintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#94A3B8")) // Slate / Muted

	errBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#EF4444")).
			Padding(0, 1).
			MarginTop(1).
			MarginBottom(1)
)

// PrintError formats and outputs a highlighted error box to stderr.
func PrintError(err error) {
	if err == nil {
		return
	}

	errStr := err.Error()
	lower := strings.ToLower(errStr)

	badge := errBadgeStyle.Render("ERROR")
	header := errHeaderStyle.Render(errStr)

	lines := []string{
		fmt.Sprintf("%s %s", badge, header),
	}

	// Detect SSH key / location / parsing errors and render helpful troubleshooting guidance
	if strings.Contains(lower, "ssh key") || strings.Contains(lower, "key file") || strings.Contains(lower, "no key found") || strings.Contains(lower, "private key") || strings.Contains(lower, "identity file") {
		lines = append(lines, "")
		lines = append(lines, errHintStyle.Render("💡 Troubleshooting Hint:"))

		if strings.Contains(lower, "not found") || strings.Contains(lower, "no such file") {
			lines = append(lines, errHintStyle.Render("   • The specified SSH key file does not exist at the given path."))
			lines = append(lines, errHintStyle.Render("   • Please verify the file path supplied via -i / --identity or update saved server credentials."))
			lines = append(lines, errHintStyle.Render("   • Ensure the file path is correct and accessible."))
		} else if strings.Contains(lower, "parse") {
			lines = append(lines, errHintStyle.Render("   • Failed to parse the SSH private key."))
			lines = append(lines, errHintStyle.Render("   • Ensure the key is a valid, unencrypted OpenSSH private key file (e.g. ~/.ssh/id_rsa)."))
		} else {
			lines = append(lines, errHintStyle.Render("   • Verify SSH identity file path, file permissions (e.g. chmod 600 key.pem), and server host details."))
		}
	} else if strings.Contains(lower, "accepts") || strings.Contains(lower, "unknown flag") || strings.Contains(lower, "flag needs an argument") {
		lines = append(lines, "")
		lines = append(lines, errHintStyle.Render("💡 Usage Hint:"))
		lines = append(lines, errHintStyle.Render("   • Use --help to view available command options and usage flags."))
	}

	fmt.Fprintln(os.Stderr, errBoxStyle.Render(strings.Join(lines, "\n")))
}
