package cmd

import (
	"errors"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/logging"
)

func TestCanonicalAgentErrorsSuggestMarkata(t *testing.T) {
	previous := activeCLIName
	activeCLIName = canonicalCLIName
	t.Cleanup(func() { activeCLIName = previous })

	_, err := normalizeAgentInstallTarget("github-copilot,claude-code")
	if err == nil || !strings.Contains(err.Error(), "markata agent list-agents") || strings.Contains(err.Error(), "markata-go agent") {
		t.Fatalf("unknown-agent diagnostic = %v; want canonical markata suggestion", err)
	}

	_, err = resolveAgentInstallTarget("", "", true)
	if err == nil || !strings.Contains(err.Error(), "markata agent list-agents") || strings.Contains(err.Error(), "markata-go agent") {
		t.Fatalf("missing global agent diagnostic = %v; want canonical markata suggestion", err)
	}
}

func TestLegacyAgentErrorsKeepMarkataGo(t *testing.T) {
	previous := activeCLIName
	activeCLIName = "markata-go"
	t.Cleanup(func() { activeCLIName = previous })

	_, err := normalizeAgentInstallTarget("invalid-agent")
	if err == nil || !strings.Contains(err.Error(), "markata-go agent list-agents") {
		t.Fatalf("legacy unknown-agent diagnostic = %v; want legacy executable name", err)
	}
}

func TestFormatErrorAsRebrandsSuggestionsAfterExecuteAs(t *testing.T) {
	// ExecuteAs restores the shared command state before main formats errors.
	// The canonical entry point therefore passes the executable name explicitly.
	previous := activeCLIName
	activeCLIName = "markata-go"
	t.Cleanup(func() { activeCLIName = previous })

	err := errors.New("Run 'markata-go agent list-agents' or inspect markata-go.toml and markata-go-site")
	got := FormatErrorAs(err, canonicalCLIName)
	if !strings.Contains(got, "markata agent list-agents") || strings.Contains(got, "markata-go agent list-agents") {
		t.Fatalf("canonical formatted diagnostic = %q", got)
	}
	if !strings.Contains(got, "markata-go.toml") || !strings.Contains(got, "markata-go-site") {
		t.Fatalf("canonical formatting must preserve file and skill names: %q", got)
	}

	legacy := FormatErrorAs(err, "markata-go")
	if !strings.Contains(legacy, "markata-go agent list-agents") {
		t.Fatalf("legacy formatted diagnostic = %q", legacy)
	}
}

func TestCanonicalHelpExamplesUseMarkata(t *testing.T) {
	previous := activeCLIName
	t.Cleanup(func() { activeCLIName = previous })

	description := "Examples:\n  markata-go agent doctor\n  markata-go build\nConfiguration: markata-go.toml\nSkill: markata-go-site"
	var canonical strings.Builder
	activeCLIName = canonicalCLIName
	renderHelpDescription(&canonical, description, logging.DefaultTheme())
	if !strings.Contains(canonical.String(), "markata agent doctor") ||
		!strings.Contains(canonical.String(), "markata build") ||
		strings.Contains(canonical.String(), "markata-go agent doctor") {
		t.Fatalf("canonical help = %q", canonical.String())
	}
	if !strings.Contains(canonical.String(), "markata-go.toml") || !strings.Contains(canonical.String(), "markata-go-site") {
		t.Fatalf("canonical help must preserve file and skill names: %q", canonical.String())
	}

	var legacy strings.Builder
	activeCLIName = "markata-go"
	renderHelpDescription(&legacy, description, logging.DefaultTheme())
	if !strings.Contains(legacy.String(), "markata-go agent doctor") || !strings.Contains(legacy.String(), "markata-go build") {
		t.Fatalf("legacy help = %q", legacy.String())
	}
}
