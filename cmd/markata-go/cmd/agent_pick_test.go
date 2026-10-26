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
		picker, _ = updateAgentPicker(t, picker, tea.KeyMsg{Type: tea.KeyUp})
	}
	last := supportedAgentTargets[len(supportedAgentTargets)-1].Name
	if !strings.Contains(picker.View(), "> "+last) {
		t.Fatalf("selected agent is outside the visible list: %s", picker.View())
	}
	for range len(supportedAgentTargets) + 2 {
		picker, _ = updateAgentPicker(t, picker, tea.KeyMsg{Type: tea.KeyDown})
	}

	if picker.cursor != 0 {
		t.Fatalf("cursor moved past first agent: %d", picker.cursor)
	}
}

func TestAgentPicker_BottomUpLayout(t *testing.T) {
	picker := newAgentPicker()
	view := picker.View()
	first := strings.Index(view, "> adal")
	second := strings.Index(view, "  amp")
	prompt := strings.LastIndex(view, "\n> ")
	if second < 0 || first <= second || prompt <= first {
		t.Fatalf("expected results above the bottom prompt, first result nearest it:\n%s", view)
	}
	if !strings.HasSuffix(view, "\n> ") {
		t.Fatalf("expected the query prompt on the final line:\n%s", view)
	}
}

func TestAgentPicker_Colors(t *testing.T) {
	oldForceColor, oldNoColor, oldLogFormat := forceColor, noColor, logFormat
	t.Cleanup(func() {
		forceColor, noColor, logFormat = oldForceColor, oldNoColor, oldLogFormat
	})
	forceColor, noColor, logFormat = true, false, ""
	picker, _ := updateAgentPicker(t, newAgentPicker(), tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("ghcp")})
	view := picker.View()
	if !strings.Contains(view, "\033[") ||
		!strings.Contains(view, colorize("g", currentLogTheme.Warning, true)) {
		t.Fatalf("expected colored selection and fuzzy matches:\n%q", view)
	}
	noColor = true
	picker = newAgentPicker()
	if strings.Contains(picker.View(), "\033[") {
		t.Fatal("--no-color must disable picker colors")
	}
	forceColor, noColor = false, false
	t.Setenv("NO_COLOR", "1")
	if strings.Contains(newAgentPicker().View(), "\033[") {
		t.Fatal("NO_COLOR must disable picker colors")
	}
}
