package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestSetupCreatesConfigAndIsIdempotent(t *testing.T) {
	workspace := filepath.Join(t.TempDir(), "new-site")
	var output bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&output)
	if err := runSetupCommand(cmd, []string{workspace}); err != nil {
		t.Fatal(err)
	}
	config := filepath.Join(workspace, defaultConfigFilename)
	original, err := os.ReadFile(config)
	if err != nil {
		t.Fatal(err)
	}
	if len(original) == 0 {
		t.Fatal("created an empty config")
	}
	if err := runSetupCommand(cmd, []string{workspace}); err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(config)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(original, after) {
		t.Fatal("setup modified an existing configuration")
	}
	if !strings.Contains(output.String(), "markata editor") {
		t.Fatal("setup did not provide next steps")
	}
}

func TestSetupPreservesExistingAlternativeConfig(t *testing.T) {
	workspace := t.TempDir()
	path := filepath.Join(workspace, "markata-go.yaml")
	if err := os.WriteFile(path, []byte("title: Existing\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	cmd := &cobra.Command{}
	if err := runSetupCommand(cmd, []string{workspace}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(workspace, defaultConfigFilename)); !os.IsNotExist(err) {
		t.Fatalf("unexpected new config: %v", err)
	}
	contents, err := os.ReadFile(path)
	if err != nil || string(contents) != "title: Existing\n" {
		t.Fatalf("existing config changed: %v, %q", err, contents)
	}
}

func TestSetupRejectsFileAsWorkspace(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "file")
	if err := os.WriteFile(file, []byte("a"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := runSetupCommand(&cobra.Command{}, []string{file}); err == nil {
		t.Fatal("expected error for file workspace")
	}
}

func TestSetupRejectsUndiscoverableConfig(t *testing.T) {
	workspace := t.TempDir()
	if err := os.WriteFile(filepath.Join(workspace, setupAlternateTOML), []byte("title = \"Existing\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	err := runSetupCommand(&cobra.Command{}, []string{workspace})
	if err == nil || !strings.Contains(err.Error(), "does not automatically discover") {
		t.Fatalf("expected guidance for unrecognized config name, got %v", err)
	}
	if _, statErr := os.Stat(filepath.Join(workspace, defaultConfigFilename)); !os.IsNotExist(statErr) {
		t.Fatalf("setup created a duplicate config: %v", statErr)
	}
}

func TestSetupAtomicCreationDoesNotReplaceExistingConfig(t *testing.T) {
	workspace := t.TempDir()
	target := filepath.Join(workspace, defaultConfigFilename)
	const original = "# existing user configuration\n"
	if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := writeSetupConfig(target); err == nil {
		t.Fatal("setup must not overwrite an existing file")
	}
	got, err := os.ReadFile(target)
	if err != nil || string(got) != original {
		t.Fatalf("existing configuration changed: %v, %q", err, got)
	}
	matches, err := filepath.Glob(filepath.Join(workspace, ".markata-setup-*.tmp"))
	if err != nil || len(matches) != 0 {
		t.Fatalf("staged files left after failure: %v, %v", matches, err)
	}
}
