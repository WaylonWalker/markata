package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
)

func TestTriagePickerCancellationAndErrors(t *testing.T) {
	options := []triageOption{{target: "codex", label: "Codex"}, {target: triageStdout, label: "Print"}}
	if target, err := triagePickerResult(options, -1, fuzzyfinder.ErrAbort); target != "" || err != nil {
		t.Fatalf("canceled picker launched %q or failed: %v", target, err)
	}
	if target, err := triagePickerResult(options, 1, nil); target != triageStdout || err != nil {
		t.Fatalf("selected wrong destination: %q, %v", target, err)
	}
	if _, err := triagePickerResult(options, -1, nil); err == nil {
		t.Fatal("expected out-of-range picker result error")
	}
	if _, err := triagePickerResult(options, 0, io.ErrUnexpectedEOF); !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("expected wrapped UI error, got %v", err)
	}
}

func TestTriageHarnessArgumentContract(t *testing.T) {
	prompt := "A real triage prompt with flags like --foo and quotes"
	cases := []struct {
		name, flag string
	}{
		{name: "opencode", flag: "--prompt"},
		{name: "claude"},
		{name: "codex"},
		{name: "gemini", flag: "-i"},
		{name: "cursor"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			agent := triageHarness{Name: tc.name, PromptFlag: tc.flag}
			args := triageHarnessArgs(agent, prompt)
			if len(args) == 0 || args[len(args)-1] != prompt {
				t.Fatalf("prompt must be one unchanged argument: %q", args)
			}
			if tc.flag == "" && len(args) != 1 {
				t.Fatalf("unexpected extra flags: %q", args)
			}
			if tc.flag != "" && (len(args) != 2 || args[0] != tc.flag) {
				t.Fatalf("unexpected initial-prompt syntax: %q", args)
			}
		})
	}
}

func TestTriageProjectEvidenceMissingSiteAndBrokenConfig(t *testing.T) {
	before := siteDir
	defer func() { siteDir = before }()
	missing := filepath.Join(t.TempDir(), "missing-site")
	siteDir = missing
	label, base := triageProjectEvidence()
	if base != "" || !strings.Contains(label, "missing or inaccessible") {
		t.Fatalf("missing site should be reported, not used: %q, %q", label, base)
	}
	if got := triageConfigEvidence(base); !strings.Contains(got, "skipped") {
		t.Fatalf("config discovery should be skipped with no site: %q", got)
	}
}

func TestTriageConfigDiscoversSelectedSite(t *testing.T) {
	before := siteDir
	defer func() { siteDir = before }()
	siteDir = t.TempDir()
	if err := os.WriteFile(filepath.Join(siteDir, "markata-go.toml"), []byte("unparseable TOML = !"), 0o600); err != nil {
		t.Fatal(err)
	}
	label, base := triageProjectEvidence()
	if base != siteDir || !strings.Contains(label, "selected site directory") {
		t.Fatalf("selected site not reported: %q, %q", label, base)
	}
	if got := triageConfigEvidence(base); !strings.Contains(got, "markata-go.toml") {
		t.Fatalf("config name missing from report: %q", got)
	}
	if got := triageConfigEvidence(base); strings.Contains(got, "unparseable TOML") {
		t.Fatalf("config contents leaked: %q", got)
	}
}

func TestTriageChildCLIWithInvalidConfiguration(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing")
	// #nosec G204 -- Executes this Go test binary with a fixed test name, not user input.
	proc := exec.Command(os.Args[0], "-test.run=^TestTriageCLISubprocess$")
	proc.Env = append(os.Environ(),
		"MARKATA_TRIAGE_CLI_TEST=1",
		"MARKATA_TRIAGE_MISSING_SITE="+missing,
	)
	var stderr bytes.Buffer
	proc.Stderr = &stderr
	out, err := proc.Output()
	if err != nil {
		t.Fatalf("triage with missing site/config failed: %v\nstderr: %s", err, stderr.String())
	}
	if !strings.HasPrefix(string(out), "# Markata troubleshooting handoff\n") {
		t.Fatalf("stdout is not a prompt: %q", out)
	}
	if !strings.Contains(string(out), "requested site directory missing or inaccessible") {
		t.Fatalf("missing site was not reported: %q", out)
	}
	if !strings.Contains(string(out), "config") {
		t.Fatal("missing configuration diagnostics")
	}
	if stderr.Len() != 0 || strings.Contains(string(out), "\x1b[") {
		t.Fatalf("unexpected status or color: stderr=%q", stderr.String())
	}
}

// Run the real Cobra root as a subprocess so global CLI flags and pre-run
// validation are tested independently of test-package state.
func TestTriageCLISubprocess(_ *testing.T) {
	if os.Getenv("MARKATA_TRIAGE_CLI_TEST") != "1" {
		return
	}
	rootCmd.SetArgs([]string{
		"--quiet", "--no-input", "--no-color",
		"--site-dir", os.Getenv("MARKATA_TRIAGE_MISSING_SITE"),
		"--config", filepath.Join(os.Getenv("MARKATA_TRIAGE_MISSING_SITE"), "broken.toml"),
		"triage",
	})
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}

type triageClosedPipeWriter struct{}

func (triageClosedPipeWriter) Write([]byte) (int, error) { return 0, syscall.EPIPE }

func TestTriageStdoutHandlesClosedPipe(t *testing.T) {
	beforeTo, beforeAgent, beforeIssue, beforeNoInput := triageTo, triageAgent, triageIssue, noInput
	defer func() {
		triageTo, triageAgent, triageIssue, noInput = beforeTo, beforeAgent, beforeIssue, beforeNoInput
	}()
	triageTo, triageAgent, triageIssue, noInput = triageStdout, "", "", true
	cmd := &cobra.Command{Use: "triage"}
	cmd.SetOut(triageClosedPipeWriter{})
	if err := runTriageCommand(cmd, nil); err != nil {
		t.Fatalf("closed consumer should be handled quietly: %v", err)
	}
}

func TestTriageStdioRespectsNoInputNoColorAndQuiet(t *testing.T) {
	beforeTo, beforeAgent, beforeIssue, beforeNoInput, beforeQuiet, beforeNoColor :=
		triageTo, triageAgent, triageIssue, noInput, quiet, noColor
	defer func() {
		triageTo, triageAgent, triageIssue, noInput, quiet, noColor =
			beforeTo, beforeAgent, beforeIssue, beforeNoInput, beforeQuiet, beforeNoColor
	}()
	t.Setenv("NO_COLOR", "1")
	t.Setenv("TERM", "dumb")
	triageTo, triageAgent, triageIssue, noInput, quiet, noColor = triageAuto, "", "", true, true, true
	var out, errOut bytes.Buffer
	cmd := &cobra.Command{Use: "triage"}
	cmd.SetIn(strings.NewReader("don't launch any agents"))
	cmd.SetOut(&out)
	cmd.SetErr(&errOut)
	if err := runTriageCommand(cmd, nil); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(out.String(), "# Markata troubleshooting handoff\n") || !strings.HasSuffix(out.String(), "\n") {
		t.Fatalf("stdout contract failed: %q", out.String())
	}
	if errOut.Len() != 0 || strings.Contains(out.String(), "\x1b[") {
		t.Fatalf("plain quiet stdout contained UI or logs: stderr=%q", errOut.String())
	}
}

func TestTriageAgentSubprocessStub(_ *testing.T) {
	if os.Getenv("MARKATA_TRIAGE_AGENT_STUB") != "1" {
		return
	}
	// The supplied prompt is one positional argv value, not shell-executed.
	_, _ = fmt.Fprintf(os.Stdout, "stub-received:%s", os.Args[len(os.Args)-1])
	os.Exit(0)
}

func TestTriageAgentLaunchForwardsPromptAndStreams(t *testing.T) {
	t.Setenv("MARKATA_TRIAGE_AGENT_STUB", "1")
	payload := "text with spaces; $(echo NOT-A-COMMAND); --flag"
	agent := triageHarness{
		Name:       "stub",
		Label:      "Fake agent",
		Path:       os.Args[0],
		PromptFlag: "-test.run=^TestTriageAgentSubprocessStub$",
	}
	var stdout, stderr bytes.Buffer
	cmd := &cobra.Command{Use: "triage"}
	cmd.SetContext(context.Background())
	cmd.SetIn(strings.NewReader(""))
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	if err := runTriageHarness(cmd, agent, payload); err != nil {
		t.Fatalf("could not launch test agent: %v, stderr=%q", err, stderr.String())
	}
	if want := "stub-received:" + payload; stdout.String() != want {
		t.Fatalf("agent did not receive exactly one prompt argument: got %q, want %q", stdout.String(), want)
	}
	if stderr.Len() != 0 {
		t.Fatalf("unexpected agent stderr %q", stderr.String())
	}
}

func TestTriageCancelledContextDoesNotLaunchAgent(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	cmd := &cobra.Command{Use: "triage"}
	cmd.SetContext(ctx)
	var out bytes.Buffer
	cmd.SetOut(&out)
	agent := triageHarness{
		Name:       "fake-agent",
		Label:      "Fake agent",
		Path:       os.Args[0],
		PromptFlag: "-test.run=^TestTriageAgentSubprocessStub$",
	}
	if err := runTriageHarness(cmd, agent, "do not launch"); err == nil {
		t.Fatal("expected canceled context to stop the subprocess")
	}
	if out.Len() != 0 {
		t.Fatalf("canceled subprocess produced output: %q", out.String())
	}
}
