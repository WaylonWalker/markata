package cmd

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func pickInstallAgent(cmd *cobra.Command) (string, error) {
	result, err := tea.NewProgram(newAgentPicker(),
		tea.WithInput(cmd.InOrStdin()),
		tea.WithOutput(cmd.OutOrStdout()),
		tea.WithContext(cmd.Context()),
	).Run()
	if err != nil {
		return "", fmt.Errorf("agent picker: %w", err)
	}
	picker, ok := result.(agentPicker)
	if !ok || picker.chosen == "" {
		return "", fmt.Errorf("agent selection canceled")
	}
	return picker.chosen, nil
}

type agentPicker struct {
	filtered []agentInstallTarget
	query    string
	cursor   int
	height   int
	chosen   string
}

func newAgentPicker() agentPicker {
	return agentPicker{filtered: supportedAgentTargets, height: 20}
}

func (m agentPicker) Init() tea.Cmd { return nil }

func (m agentPicker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEscape:
			return m, tea.Quit
		case tea.KeyEnter:
			if len(m.filtered) > 0 {
				m.chosen = m.filtered[m.cursor].Name
				return m, tea.Quit
			}
		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			}
		case tea.KeyDown:
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
			}
		case tea.KeyBackspace, tea.KeyDelete:
			query := []rune(m.query)
			if len(query) > 0 {
				m.query = string(query[:len(query)-1])
				m.applyFilter()
			}
		case tea.KeyRunes:
			m.query += string(msg.Runes)
			m.applyFilter()
		default:
			return m, nil
		}
	}
	return m, nil
}

func (m *agentPicker) applyFilter() {
	m.filtered = nil
	m.cursor = 0
	for _, target := range supportedAgentTargets {
		if fuzzyMatch(target.Name, strings.ToLower(m.query)) {
			m.filtered = append(m.filtered, target)
		}
	}
}

func (m agentPicker) View() string {
	var view strings.Builder
	view.WriteString("Choose an agent for global installation\n\n> " + m.query + "\n\n")
	height := m.height - 7
	if height < 1 {
		height = 1
	}
	start, end := scrollRange(m.cursor, len(m.filtered), height)
	for i := start; i < end; i++ {
		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}
		fmt.Fprintf(&view, "%s%s  (~/%s)\n", cursor, m.filtered[i].Name, m.filtered[i].GlobalPath)
	}
	if len(m.filtered) == 0 {
		view.WriteString("No matching agents\n")
	}
	view.WriteString("\nType to fuzzy-filter | Up/Down to navigate | Enter to select | Esc/Ctrl+C to cancel\n")
	return view.String()
}
