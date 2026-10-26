package cmd

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func updateAgentPicker(t *testing.T, picker agentPicker, msg tea.Msg) (agentPicker, tea.Cmd) {
	t.Helper()
	model, command := picker.Update(msg)
	updated, ok := model.(agentPicker)
	if !ok {
		t.Fatalf("unexpected picker model %T", model)
	}
	return updated, command
}

func TestAgentPicker_FuzzySelection(t *testing.T) {
	for _, query := range []string{"copilot", "GHCP"} {
		t.Run(query, func(t *testing.T) {
			picker, _ := updateAgentPicker(t, newAgentPicker(), tea.KeyMsg{Type: tea.KeyDown})
			picker, _ = updateAgentPicker(t, picker, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(query)})
			if !strings.Contains(picker.View(), "github-copilot  (~/.copilot/skills)") {
				t.Fatalf("missing matching agent and destination: %s", picker.View())
			}
			picker, command := updateAgentPicker(t, picker, tea.KeyMsg{Type: tea.KeyEnter})
			if picker.chosen != "github-copilot" || command == nil {
				t.Fatalf("chosen = %q, quit command = %v", picker.chosen, command)
			}
		})
	}
}

func TestAgentPicker_NoMatchesAndBackspace(t *testing.T) {
	picker, _ := updateAgentPicker(t, newAgentPicker(), tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("copilot界")})
	if !strings.Contains(picker.View(), "No matching agents") {
		t.Fatal("expected empty-result guidance")
	}
	picker, command := updateAgentPicker(t, picker, tea.KeyMsg{Type: tea.KeyEnter})
	if command != nil || picker.chosen != "" {
		t.Fatal("Enter must not select or quit with no matches")
	}
	picker, _ = updateAgentPicker(t, picker, tea.KeyMsg{Type: tea.KeyBackspace})
	picker, _ = updateAgentPicker(t, picker, tea.KeyMsg{Type: tea.KeyEnter})
	if picker.chosen != "github-copilot" {
		t.Fatalf("backspace did not restore matches: %q", picker.chosen)
	}
}

func TestAgentPicker_Cancel(t *testing.T) {
	for _, key := range []tea.KeyType{tea.KeyEscape, tea.KeyCtrlC} {
		picker, command := updateAgentPicker(t, newAgentPicker(), tea.KeyMsg{Type: key})
		if picker.chosen != "" || command == nil {
			t.Fatalf("cancel selected %q or did not quit", picker.chosen)
		}
	}
}

func TestAgentPicker_NavigationAndScrolling(t *testing.T) {
	picker, _ := updateAgentPicker(t, newAgentPicker(), tea.WindowSizeMsg{Height: 8})
	for range len(supportedAgentTargets) + 2 {
		picker, _ = updateAgentPicker(t, picker, tea.KeyMsg{Type: tea.KeyDown})
	}
	last := supportedAgentTargets[len(supportedAgentTargets)-1].Name
	if !strings.Contains(picker.View(), "> "+last) {
		t.Fatalf("selected agent is outside the visible list: %s", picker.View())
	}
	for range len(supportedAgentTargets) + 2 {
		picker, _ = updateAgentPicker(t, picker, tea.KeyMsg{Type: tea.KeyUp})
	}
	if picker.cursor != 0 {
		t.Fatalf("cursor moved past first agent: %d", picker.cursor)
	}
}
