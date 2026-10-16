package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionCommand_ShortUsesStdoutOnly(t *testing.T) {
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	versionCmd.SetOut(stdout)
	versionCmd.SetErr(stderr)
	currentCmd = versionCmd
	defer func() { currentCmd = nil }()

	if err := versionCmd.Flags().Set("short", "true"); err != nil {
		t.Fatalf("set short flag: %v", err)
	}
	defer func() {
		if err := versionCmd.Flags().Set("short", "false"); err != nil {
			t.Fatalf("reset short flag: %v", err)
		}
	}()

	versionCmd.Run(versionCmd, nil)

	if got := strings.TrimSpace(stdout.String()); got != GetVersion() {
		t.Fatalf("stdout = %q, want %q", got, GetVersion())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr output, got %q", stderr.String())
	}
}

func TestCanonicalVersionBranding(t *testing.T) {
	previousName := activeCLIName
	activeCLIName = "markata"
	defer func() { activeCLIName = previousName }()

	if got := GetVersionInfo(); !strings.HasPrefix(got, "markata "+Version+"\n") {
		t.Fatalf("canonical version output = %q", got)
	}
	if got := GetShortVersionInfo(); !strings.HasPrefix(got, "markata/"+Version+" ") {
		t.Fatalf("canonical short version = %q", got)
	}
}
