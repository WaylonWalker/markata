package cmd

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestDetectTriageHarnesses(t *testing.T) {
	found := detectTriageHarnesses(func(name string) (string, error) {
		if name == "opencode" || name == "codex" {
			return "/bin/" + name, nil
		}
		return "", errors.New("missing")
	})
	if len(found) != 2 || found[0].Name != "opencode" || found[1].Name != "codex" {
		t.Fatalf("unexpected harness detection: %+v", found)
	}
	if found[0].Path != "/bin/opencode" || found[1].Path != "/bin/codex" {
		t.Fatalf("missing discovered paths: %+v", found)
	}
}

func TestTriagePromptAsksBeforeChangingAnything(t *testing.T) {
	facts := triageFacts{
		Version: "0.99.1", Commit: "abc123", BuildDate: "2026-10-19",
		GoVersion: "go1.26", OS: "linux", Arch: "amd64", Distro: "Arch",
		Executable: "~/bin/markata", Install: "user-local", Project: "~/site",
		Config: "markata.toml", Env: "CI",
	}
	prompt := buildTriagePrompt(facts, "feed page is missing")
	for _, required := range []string{
		"0.99.1", "abc123", "~/bin/markata", "~/site",
		"feed page is missing", "Ask", "affected version",
		"before making changes", "https://github.com/WaylonWalker/markata",
		"No raw environment values", "not verified as the affected installation",
	} {
		if !strings.Contains(prompt, required) {
			t.Fatalf("expected %q in triage prompt", required)
		}
	}
}

func TestTriagePromptDoesNotDumpEnvironmentValues(t *testing.T) {
	t.Setenv("MARKATA_GO_CONFIG", "super-secret-config-value")
	t.Setenv("CI", "CI-secret-token")
	evidence := triageEnvironmentEvidence()
	if !strings.Contains(evidence, "MARKATA_GO_CONFIG") || !strings.Contains(evidence, "CI") {
		t.Fatalf("missing allowlisted variable names: %q", evidence)
	}
	if strings.Contains(evidence, "super-secret") || strings.Contains(evidence, "CI-secret") {
		t.Fatalf("triage leaked environment values: %q", evidence)
	}
}

func TestTriageSafePathRedactsHome(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		t.Skip("home directory unavailable")
	}
	path := filepath.Join(home, "site", "markata.toml")
	safe := triageSafePath(path)
	if !strings.HasPrefix(safe, "~") {
		t.Fatalf("expected tilde path, got %q", safe)
	}
	if strings.Contains(safe, home) {
		t.Fatalf("expected home to be redacted, got %q", safe)
	}
}

func TestTriageInstallHintIsExplicitlyTentative(t *testing.T) {
	if got := triageInstallHint("/home/example/.local/share/mise/installs/markata/1/bin/markata"); !strings.Contains(got, "guess") {
		t.Fatalf("expected tentative detection, got %q", got)
	}
	if got := triageInstallHint("/usr/bin/markata"); !strings.Contains(got, "unknown") {
		t.Fatalf("expected unknown install source, got %q", got)
	}
}

func TestTriageClipboardCandidates(t *testing.T) {
	if got := triageClipboardCommands("darwin"); len(got) != 1 || got[0].Binary != "pbcopy" {
		t.Fatalf("unexpected macOS candidates: %+v", got)
	}
	if got := triageClipboardCommands("windows"); len(got) != 1 || got[0].Binary != "clip.exe" {
		t.Fatalf("unexpected Windows candidates: %+v", got)
	}
	if got := triageClipboardCommands("linux"); len(got) < 3 || got[0].Binary != "wl-copy" {
		t.Fatalf("unexpected Linux candidates: %+v", got)
	}
}

func TestTriageStdoutIsPromptOnly(t *testing.T) {
	beforeTo, beforeAgent, beforeIssue, beforeNoInput := triageTo, triageAgent, triageIssue, noInput
	defer func() {
		triageTo, triageAgent, triageIssue, noInput = beforeTo, beforeAgent, beforeIssue, beforeNoInput
	}()

	triageTo, triageAgent, triageIssue, noInput = "auto", "", "", true
	var stdout, stderr bytes.Buffer
	cmd := &cobra.Command{Use: "triage"}
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	if err := runTriageCommand(cmd, nil); err != nil {
		t.Fatalf("run triage: %v", err)
	}
	if !strings.HasPrefix(stdout.String(), "# Markata troubleshooting handoff\n") {
		t.Fatalf("stdout is not a prompt: %q", stdout.String())
	}
	if stderr.Len() > 0 {
		t.Fatalf("unexpected stderr: %q", stderr.String())
	}
	if strings.Contains(stdout.String(), "\x1b[") {
		t.Fatalf("stdout should contain no terminal escape sequences")
	}
}

func TestTriageInitialIssueBounded(t *testing.T) {
	beforeTo, beforeAgent, beforeIssue := triageTo, triageAgent, triageIssue
	defer func() {
		triageTo, triageAgent, triageIssue = beforeTo, beforeAgent, beforeIssue
	}()
	triageTo, triageAgent, triageIssue = "stdout", "", strings.Repeat("x", 4001)
	cmd := &cobra.Command{Use: "triage"}
	err := runTriageCommand(cmd, nil)
	if err == nil || !strings.Contains(err.Error(), "4000") {
		t.Fatalf("expected a bounded issue error, got %v", err)
	}
}

func TestTriageShortRemovesControlCharacters(t *testing.T) {
	if got := triageShort("hello\nworld\tsecret", 40); got != "hello world secret" {
		t.Fatalf("unexpected sanitized field %q", got)
	}
}
