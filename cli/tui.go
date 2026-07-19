package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/devlopersabbir/vpcm/internal/inventory"
)

type tuiModel struct {
	servers  []inventory.Server
	filtered []inventory.Server
	cursor   int
	search   textinput.Model
	selected *inventory.Server
	quitting bool
}

func initialModel(servers []inventory.Server) tuiModel {
	ti := textinput.New()
	ti.Placeholder = "Type to search servers (e.g. aws, prod, db)..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40

	return tuiModel{
		servers:  servers,
		filtered: servers,
		cursor:   0,
		search:   ti,
	}
}

func (m tuiModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			return m, tea.Quit

		case tea.KeyEnter:
			if len(m.filtered) > 0 && m.cursor >= 0 && m.cursor < len(m.filtered) {
				m.selected = &m.filtered[m.cursor]
			}
			m.quitting = true
			return m, tea.Quit

		case tea.KeyUp, tea.KeyCtrlP:
			if m.cursor > 0 {
				m.cursor--
			} else if len(m.filtered) > 0 {
				m.cursor = len(m.filtered) - 1 // wrap
			}

		case tea.KeyDown, tea.KeyCtrlN:
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
			} else {
				m.cursor = 0 // wrap
			}
		}

	default:
		// Do nothing
	}

	// Update text input
	m.search, cmd = m.search.Update(msg)

	// Filter list based on search term
	term := strings.ToLower(m.search.Value())
	if term == "" {
		m.filtered = m.servers
	} else {
		var matched []inventory.Server
		for _, s := range m.servers {
			// Search hostname, name, provider, username
			if strings.Contains(strings.ToLower(s.Name), term) ||
				strings.Contains(strings.ToLower(s.Host), term) ||
				strings.Contains(strings.ToLower(s.Provider), term) ||
				strings.Contains(strings.ToLower(s.Username), term) {
				matched = append(matched, s)
			}
		}
		m.filtered = matched
	}

	// Adjust cursor if filtered list is smaller
	if m.cursor >= len(m.filtered) {
		m.cursor = len(m.filtered) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}

	return m, cmd
}

func (m tuiModel) View() string {
	if m.quitting {
		return ""
	}

	var sb strings.Builder

	// Header style
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		Padding(0, 1)

	sb.WriteString("\n" + headerStyle.Render("🔍 VPSM Interactive Server Explorer") + "\n")
	sb.WriteString(strings.Repeat("─", 60) + "\n\n")

	// Search bar
	sb.WriteString("  Search: " + m.search.View() + "\n\n")
	sb.WriteString(strings.Repeat("─", 60) + "\n\n")

	// Server list
	if len(m.filtered) == 0 {
		sb.WriteString("  No matching servers found.\n")
	} else {
		for i, s := range m.filtered {
			cursorStr := "  "
			rowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))

			if i == m.cursor {
				cursorStr = "> "
				rowStyle = lipgloss.NewStyle().
					Bold(true).
					Foreground(lipgloss.Color("#00FFFF"))
			}

			// Format row
			provider := s.Provider
			if provider == "" {
				provider = "Generic"
			}
			rowText := fmt.Sprintf("[%d] %s (%s@%s) - %s", s.ID, s.Name, s.Username, s.Host, provider)
			sb.WriteString(cursorStr + rowStyle.Render(rowText) + "\n")
		}
	}

	sb.WriteString("\n" + strings.Repeat("─", 60) + "\n")
	footerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Padding(0, 1)
	sb.WriteString(footerStyle.Render("Use Up/Down to select | Enter to Connect | Esc/Ctrl+C to Quit") + "\n")

	return sb.String()
}

func runTUI(ctx context.Context, servers []inventory.Server) (*inventory.Server, error) {
	p := tea.NewProgram(initialModel(servers), tea.WithContext(ctx))
	m, err := p.Run()
	if err != nil {
		return nil, err
	}

	model := m.(tuiModel)
	return model.selected, nil
}
