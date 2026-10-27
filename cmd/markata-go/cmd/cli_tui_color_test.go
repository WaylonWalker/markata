package cmd

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func TestConfigureTUIColorForUsesSharedPolicy(t *testing.T) {
	tests := []struct {
		name       string
		noColor    bool
		force      bool
		format     string
		isTTY      bool
		noColorEnv string
		forceEnv   string
		term       string
		want       termenv.Profile
	}{
		{
			name:       "no color flag disables TUI styles",
			noColor:    true,
			isTTY:      true,
			noColorEnv: "",
			term:       "xterm-256color",
			want:       termenv.Ascii,
		},
		{
			name:       "color flag overrides environment",
			force:      true,
			isTTY:      false,
			noColorEnv: "1",
			term:       "dumb",
			want:       termenv.TrueColor,
		},
		{
			name:       "plain log format disables TUI styles",
			force:      true,
			format:     "plain",
			isTTY:      true,
			noColorEnv: "",
			term:       "xterm-256color",
			want:       termenv.Ascii,
		},
		{
			name:       "redirected output disables automatic TUI styles",
			isTTY:      false,
			noColorEnv: "",
			term:       "xterm-256color",
			want:       termenv.Ascii,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldNoColor, oldForceColor, oldLogFormat := noColor, forceColor, logFormat
			oldProfile := lipgloss.DefaultRenderer().ColorProfile()
			t.Cleanup(func() {
				noColor, forceColor, logFormat = oldNoColor, oldForceColor, oldLogFormat
				lipgloss.SetColorProfile(oldProfile)
			})

			t.Setenv("NO_COLOR", tt.noColorEnv)
			t.Setenv("FORCE_COLOR", tt.forceEnv)
			t.Setenv("TERM", tt.term)
			noColor, forceColor, logFormat = tt.noColor, tt.force, tt.format

			configureTUIColorFor(tt.isTTY)
			if got := lipgloss.DefaultRenderer().ColorProfile(); got != tt.want {
				t.Fatalf("color profile = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigureTUIColorForForceColorEnv(t *testing.T) {
	oldNoColor, oldForceColor, oldLogFormat := noColor, forceColor, logFormat
	oldProfile := lipgloss.DefaultRenderer().ColorProfile()
	t.Cleanup(func() {
		noColor, forceColor, logFormat = oldNoColor, oldForceColor, oldLogFormat
		lipgloss.SetColorProfile(oldProfile)
	})
	t.Setenv("NO_COLOR", "1")
	t.Setenv("FORCE_COLOR", "1")
	t.Setenv("TERM", "dumb")
	noColor, forceColor, logFormat = false, false, "auto"

	configureTUIColorFor(false)
	if got := lipgloss.DefaultRenderer().ColorProfile(); got != termenv.TrueColor {
		t.Fatalf("FORCE_COLOR profile = %v, want %v", got, termenv.TrueColor)
	}
}
