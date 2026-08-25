package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// ─── All known vpsm command tokens ────────────────────────────────────────────

var knownCommands = []string{
	"server list",
	"server add",
	"server remove",
	"server rename",
	"server favorite",
	"server flush",
	"server export",
	"server import",
	"ssh",
	"audit --name",
	"audit --host",
	"audit --id",
	"list",
	"list --favorites",
	"list --recents",
	"list --interactive",
	"config show",
	"config edit",
	"config reload",
	"config init",
	"api start",
	"api stop",
	"api restart",
	"api status",
	"api logs",
	"api logs --follow",
	"completion bash",
	"completion zsh",
	"completion fish",
	"completion powershell",
	"doctor",
	"version",
	"shell",
}

// ─── ANSI helpers ─────────────────────────────────────────────────────────────

const (
	ansiReset    = "\033[0m"
	ansiGray     = "\033[38;5;240m" // dark gray ghost text
	ansiCyan     = "\033[36m"
	ansiBold     = "\033[1m"
	ansiClearEOL = "\033[K" // erase to end of line
)

// ─── Ghost-text suggestion engine ─────────────────────────────────────────────

// suggest returns the best single inline suggestion for the current input.
// It checks dynamic server names first, then falls back to static commands.
func suggest(input string, serverNames, serverHosts []string) string {
	if input == "" {
		return ""
	}
	lower := strings.ToLower(input)

	// Dynamic: server names (for audit --name <partial> and ssh <partial>)
	for _, n := range serverNames {
		if strings.HasPrefix(strings.ToLower(n), lower) && len(n) > len(input) {
			return n[len(input):]
		}
	}
	// Dynamic: server hosts (for audit --host <partial>)
	for _, h := range serverHosts {
		if strings.HasPrefix(strings.ToLower(h), lower) && len(h) > len(input) {
			return h[len(input):]
		}
	}
	// Static commands
	for _, cmd := range knownCommands {
		if strings.HasPrefix(strings.ToLower(cmd), lower) && len(cmd) > len(input) {
			return cmd[len(input):]
		}
	}
	return ""
}

// ─── Raw terminal readline ────────────────────────────────────────────────────

// readlineWithGhost reads a line from stdin in raw mode, drawing inline
// ghost text (fish-shell style) after the cursor as the user types.
func readlineWithGhost(prompt string, serverNames, serverHosts []string) (string, error) {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		// Fallback to plain input if raw mode fails
		fmt.Print(prompt)
		var line string
		_, _ = fmt.Scanln(&line)
		return line, nil
	}
	defer term.Restore(fd, oldState)

	fmt.Print(prompt)

	var buf []rune // the user's actual typed input

	redraw := func() {
		// Move cursor to just after prompt (erase the line from prompt onwards)
		// \r moves to start of line, then re-print prompt + buf + ghost
		ghost := suggest(string(buf), serverNames, serverHosts)

		// Erase from start of typed section
		fmt.Print("\r" + prompt + ansiClearEOL)
		// Print user input
		fmt.Print(string(buf))
		// Print ghost text in gray
		if ghost != "" {
			fmt.Print(ansiGray + ghost + ansiReset)
			// Move cursor back over the ghost text
			fmt.Printf("\033[%dD", len([]rune(ghost)))
		}
	}

	for {
		// Read one byte at a time
		b := make([]byte, 4)
		n, err := os.Stdin.Read(b)
		if err != nil || n == 0 {
			fmt.Println()
			return string(buf), nil
		}

		switch {
		case b[0] == '\r' || b[0] == '\n': // Enter
			// Accept — erase ghost and move to new line
			ghost := suggest(string(buf), serverNames, serverHosts)
			if ghost != "" {
				fmt.Print(ansiGray + ghost + ansiReset)
			}
			fmt.Println()
			return string(buf), nil

		case b[0] == 3: // Ctrl+C
			fmt.Println()
			return "", fmt.Errorf("interrupted")

		case b[0] == 4: // Ctrl+D / EOF
			fmt.Println()
			return "", fmt.Errorf("EOF")

		case b[0] == 127 || b[0] == 8: // Backspace / Delete
			if len(buf) > 0 {
				buf = buf[:len(buf)-1]
			}
			redraw()

		case b[0] == '\033' && n >= 3 && b[1] == '[': // Escape sequences (arrow keys)
			switch b[2] {
			case 'C': // Right arrow → accept ghost suggestion
				ghost := suggest(string(buf), serverNames, serverHosts)
				if ghost != "" {
					buf = append(buf, []rune(ghost)...)
					redraw()
				}
			}

		case b[0] == '\t': // Tab → accept ghost suggestion
			ghost := suggest(string(buf), serverNames, serverHosts)
			if ghost != "" {
				buf = append(buf, []rune(ghost)...)
				redraw()
			}

		case b[0] >= 32 && b[0] < 127: // Printable ASCII
			buf = append(buf, rune(b[0]))
			redraw()

		case b[0] > 127: // Multi-byte UTF-8
			r := []rune(string(b[:n]))
			buf = append(buf, r...)
			redraw()
		}
	}
}

// ─── Shell command ────────────────────────────────────────────────────────────

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Interactive REPL with inline ghost-text auto-completion",
	Long: `Start an interactive vpsm shell with fish-style inline ghost text completion.

As you type, the predicted next token appears grayed out to the right.
  → Press Right Arrow or Tab to accept the suggestion
  → Press Enter to run the command
  → Press Ctrl+C or type 'exit' to quit

Server names and host IPs are loaded live from your local inventory database.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Pre-load server names and hosts from DB for dynamic completions
		var serverNames, serverHosts []string
		repo, _, err := initRepoAndService(context.Background())
		if err == nil {
			servers, err := repo.List(context.Background())
			if err == nil {
				for _, s := range servers {
					serverNames = append(serverNames, s.Name)
					serverHosts = append(serverHosts, s.Host)
				}
			}
		}

		// Print welcome banner
		fmt.Printf("\n%s%s  vpsm interactive shell%s\n", ansiBold, ansiCyan, ansiReset)
		fmt.Printf("%sType a command. Right arrow or Tab accepts ghost suggestion. 'exit' to quit.%s\n\n",
			ansiGray, ansiReset)

		prompt := ansiBold + ansiCyan + "vpsm" + ansiReset + " › "

		for {
			input, err := readlineWithGhost(prompt, serverNames, serverHosts)
			if err != nil {
				// Ctrl+C or EOF — exit cleanly
				fmt.Printf("\n%sGoodbye.%s\n", ansiGray, ansiReset)
				return nil
			}

			input = strings.TrimSpace(input)
			if input == "" {
				continue
			}
			if input == "exit" || input == "quit" || input == "q" {
				fmt.Printf("%sGoodbye.%s\n", ansiGray, ansiReset)
				return nil
			}

			// Execute the command through cobra root
			rootCmd.SetArgs(strings.Fields(input))
			if execErr := rootCmd.Execute(); execErr != nil {
				PrintError(execErr)
			}

			// Reload server list after mutating commands
			if strings.HasPrefix(input, "server add") || strings.HasPrefix(input, "server remove") {
				repo, _, err = initRepoAndService(context.Background())
				if err == nil {
					servers, err := repo.List(context.Background())
					if err == nil {
						serverNames = nil
						serverHosts = nil
						for _, s := range servers {
							serverNames = append(serverNames, s.Name)
							serverHosts = append(serverHosts, s.Host)
						}
					}
				}
			}

			fmt.Println()
		}
	},
}
