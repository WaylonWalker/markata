package cmd

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/logging"
	"github.com/spf13/cobra"
)

func TestRunLSPSetup_PrintsRequestedEditorSnippet(t *testing.T) {
	oldEditor := lspSetupEditor
	defer func() { lspSetupEditor = oldEditor }()

	command := &cobra.Command{Use: "setup"}
	stdout := bytes.NewBuffer(nil)
	command.SetOut(stdout)

	lspSetupEditor = "helix"
	if err := runLSPSetup(command, nil); err != nil {
		t.Fatalf("runLSPSetup() error = %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "[language-server.markata-go]") {
		t.Errorf("setup output = %q, want Helix language server configuration", output)
	}
	if !strings.Contains(output, `args = ["lsp"]`) {
		t.Errorf("setup output = %q, want LSP command arguments", output)
	}
}

func TestRunLSPSetup_PrintsSeparatedInstalledEditorSnippets(t *testing.T) {
	oldEditor := lspSetupEditor
	oldEditors := lspEditors
	oldLookPath := lspLookPath
	defer func() {
		lspSetupEditor = oldEditor
		lspEditors = oldEditors
		lspLookPath = oldLookPath
	}()

	lspSetupEditor = ""
	lspEditors = []lspEditor{
		{name: "Neovim", executable: "nvim"},
		{name: "Helix", executable: "hx"},
	}
	lspLookPath = func(string) (string, error) { return "/usr/bin/editor", nil }
	command := &cobra.Command{Use: "setup"}
	stdout := bytes.NewBuffer(nil)
	command.SetOut(stdout)

	if err := runLSPSetup(command, nil); err != nil {
		t.Fatalf("runLSPSetup() error = %v", err)
	}
	output := stdout.String()
	for _, want := range []string{"Neovim", "Helix", `vim.lsp.config("markata"`, "[language-server.markata-go]"} {
		if !strings.Contains(output, want) {
			t.Errorf("setup output = %q, want %q", output, want)
		}
	}
}

func TestRenderLSPSetupSnippets_SyntaxHighlightsInteractiveOutput(t *testing.T) {
	oldForceColor := forceColor
	defer func() { forceColor = oldForceColor }()
	forceColor = true

	output := renderLSPSetupSnippets([]string{"neovim", "helix"})
	if !strings.Contains(output, "\x1b[") {
		t.Errorf("setup output = %q, want ANSI syntax highlighting", output)
	}
	if !strings.Contains(output, "Neovim") || !strings.Contains(output, "Helix") {
		t.Errorf("setup output = %q, want editor headings", output)
	}
}

func TestRunLSPSetup_RejectsUnsupportedEditor(t *testing.T) {
	oldEditor := lspSetupEditor
	defer func() { lspSetupEditor = oldEditor }()

	lspSetupEditor = "unknown"
	err := runLSPSetup(&cobra.Command{Use: "setup"}, nil)
	if err == nil {
		t.Fatal("runLSPSetup() error = nil, want unsupported editor error")
	}
	if !strings.Contains(err.Error(), "supported editors: emacs, generic, helix, neovim, vscode, zed") {
		t.Errorf("error = %q, want supported editor list", err)
	}
	if got := ExitCodeForError(err); got != exitCodeUsage {
		t.Errorf("ExitCodeForError() = %d, want %d", got, exitCodeUsage)
	}
}

func TestRunLSPDoctor_ReportsValidConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "markata-go.toml")
	if err := os.WriteFile(configPath, []byte("[markata-go]\ntitle = \"Test\"\n"), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	oldConfigFile := cfgFile
	oldCmd := currentCmd
	defer func() {
		cfgFile = oldConfigFile
		currentCmd = oldCmd
	}()

	command := &cobra.Command{Use: "doctor"}
	stdout := bytes.NewBuffer(nil)
	command.SetOut(stdout)
	cfgFile = configPath

	if err := runLSPDoctor(command, nil); err != nil {
		t.Fatalf("runLSPDoctor() error = %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "Configuration: "+configPath+" is valid.") {
		t.Errorf("doctor output = %q, want valid configuration", output)
	}
	if !strings.Contains(output, "LSP prerequisites are ready.") {
		t.Errorf("doctor output = %q, want success summary", output)
	}
}

func TestRunLSPDoctor_AllowsMissingConfig(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("change working directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(oldWd); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	}()

	oldConfigFile := cfgFile
	oldCmd := currentCmd
	defer func() {
		cfgFile = oldConfigFile
		currentCmd = oldCmd
	}()

	command := &cobra.Command{Use: "doctor"}
	stdout := bytes.NewBuffer(nil)
	command.SetOut(stdout)
	cfgFile = ""

	if err := runLSPDoctor(command, nil); err != nil {
		t.Fatalf("runLSPDoctor() error = %v", err)
	}
	if !strings.Contains(stdout.String(), "Configuration: none found") {
		t.Errorf("doctor output = %q, want missing configuration information", stdout.String())
	}
}

func TestRunLSPDoctor_RejectsInvalidConfig(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "markata-go.toml")
	if err := os.WriteFile(configPath, []byte("[markata-go\ntitle = \"broken\"\n"), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	oldConfigFile := cfgFile
	oldCmd := currentCmd
	defer func() {
		cfgFile = oldConfigFile
		currentCmd = oldCmd
	}()
	cfgFile = configPath

	command := &cobra.Command{Use: "doctor"}
	stdout := bytes.NewBuffer(nil)
	command.SetOut(stdout)

	err := runLSPDoctor(command, nil)
	if err == nil {
		t.Fatal("runLSPDoctor() error = nil, want invalid configuration error")
	}
	if !strings.Contains(err.Error(), "configuration is invalid") {
		t.Errorf("error = %q, want invalid configuration error", err)
	}
	if !strings.Contains(stdout.String(), "FAIL  configuration is invalid") {
		t.Errorf("doctor output = %q, want failure diagnostic", stdout.String())
	}
}

func TestInstalledLSPEditors_FiltersUnavailableEditors(t *testing.T) {
	editors := []lspEditor{
		{name: "Available", executable: "available"},
		{name: "Missing", executable: "missing"},
	}
	lookup := func(name string) (string, error) {
		if name == "available" {
			return "/usr/bin/available", nil
		}
		return "", os.ErrNotExist
	}

	installed := installedLSPEditors(editors, lookup)
	if len(installed) != 1 || installed[0].name != "Available" {
		t.Errorf("installedLSPEditors() = %#v, want only Available", installed)
	}
}

func TestRunLSPEditorDiagnostics_ReportsAvailableVerification(t *testing.T) {
	oldEditors := lspEditors
	oldLookPath := lspLookPath
	defer func() {
		lspEditors = oldEditors
		lspLookPath = oldLookPath
	}()

	lspLookPath = func(string) (string, error) { return "/usr/bin/test-editor", nil }

	output := bytes.NewBuffer(nil)
	lspEditors = []lspEditor{{name: "Test editor", executable: "test-editor", supportsHeadless: true}}
	runLSPEditorDiagnostics(false, output, logging.DefaultTheme())
	if !strings.Contains(output.String(), "without --no-verify-editor") {
		t.Errorf("default diagnostics = %q, want verification hint", output.String())
	}

	for _, tt := range []struct {
		status lspEditorStatus
		want   string
	}{
		{lspEditorConfigured, "PASS  Editor: Test editor configured"},
		{lspEditorUnconfigured, "WARN  Editor: Test editor unconfigured"},
		{lspEditorInconclusive, "WARN  Editor: Test editor inconclusive"},
	} {
		t.Run(string(tt.status), func(t *testing.T) {
			lspEditors[0].verify = func(context.Context) lspEditorVerification {
				return lspEditorVerification{status: tt.status, message: "test result"}
			}
			output.Reset()
			runLSPEditorDiagnostics(true, output, logging.DefaultTheme())
			if !strings.Contains(output.String(), tt.want) {
				t.Errorf("verification diagnostics = %q, want %q", output.String(), tt.want)
			}
		})
	}
}

func TestRunLSPDoctor_WarnsForUnsupportedLSPConfigFormat(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "markata-go.yml")
	if err := os.WriteFile(configPath, []byte("title: Test\n"), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	oldConfigFile := cfgFile
	oldCmd := currentCmd
	defer func() {
		cfgFile = oldConfigFile
		currentCmd = oldCmd
	}()
	cfgFile = configPath
	command := &cobra.Command{Use: "doctor"}
	stdout := bytes.NewBuffer(nil)
	command.SetOut(stdout)

	if err := runLSPDoctor(command, nil); err != nil {
		t.Fatalf("runLSPDoctor() error = %v", err)
	}
	if !strings.Contains(stdout.String(), "LSP mention indexing reads only TOML and .yaml files") {
		t.Errorf("doctor output = %q, want unsupported LSP config format warning", stdout.String())
	}
}

func TestVerifyNeovimLSP_ClassifiesCommandResults(t *testing.T) {
	oldRunCommand := lspRunCommand
	defer func() { lspRunCommand = oldRunCommand }()

	tests := []struct {
		name       string
		output     string
		err        error
		wantStatus lspEditorStatus
	}{
		{"configured", "MARKATA_CONFIGURED\n", nil, lspEditorConfigured},
		{"not configured", "MARKATA_NOT_CONFIGURED\n", errors.New("exit status 2"), lspEditorUnconfigured},
		{"invalid", "MARKATA_CONFIG_INVALID\n", errors.New("exit status 2"), lspEditorUnconfigured},
		{"command failure", "startup failure", errors.New("exit status 1"), lspEditorInconclusive},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lspRunCommand = func(_ context.Context, name string, args ...string) ([]byte, error) {
				if name != "nvim" || len(args) != 2 || args[0] != "--headless" {
					t.Fatalf("command = %q %q, want nvim --headless <script>", name, args)
				}
				return []byte(tt.output), tt.err
			}
			result := verifyNeovimLSP(context.Background())
			if result.status != tt.wantStatus {
				t.Errorf("status = %q, want %q", result.status, tt.wantStatus)
			}
		})
	}
}

func TestVerifyHelixLSP_ClassifiesCommandResults(t *testing.T) {
	oldRunCommand := lspRunCommand
	defer func() { lspRunCommand = oldRunCommand }()

	tests := []struct {
		name       string
		output     string
		err        error
		wantStatus lspEditorStatus
	}{
		{"configured", "markdown markata-go", nil, lspEditorConfigured},
		{"not configured", "markdown marksman", nil, lspEditorUnconfigured},
		{"command failure", "", errors.New("exit status 1"), lspEditorInconclusive},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lspRunCommand = func(_ context.Context, name string, args ...string) ([]byte, error) {
				if name != "hx" || len(args) != 2 || args[0] != "--health" || args[1] != "markdown" {
					t.Fatalf("command = %q %q, want hx --health markdown", name, args)
				}
				return []byte(tt.output), tt.err
			}
			if result := verifyHelixLSP(context.Background()); result.status != tt.wantStatus {
				t.Errorf("status = %q, want %q", result.status, tt.wantStatus)
			}
		})
	}
}

func TestVerifyEmacsLSP_ClassifiesCommandResults(t *testing.T) {
	oldRunCommand := lspRunCommand
	defer func() { lspRunCommand = oldRunCommand }()

	tests := []struct {
		name       string
		output     string
		err        error
		wantStatus lspEditorStatus
	}{
		{"configured", "((markdown-mode . (markata-go lsp)))", nil, lspEditorConfigured},
		{"not configured", "NO_EGLOT", nil, lspEditorUnconfigured},
		{"command failure", "", errors.New("exit status 1"), lspEditorInconclusive},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lspRunCommand = func(_ context.Context, name string, args ...string) ([]byte, error) {
				if name != "emacs" || len(args) != 3 || args[0] != "--batch" || args[1] != "--eval" {
					t.Fatalf("command = %q %q, want emacs --batch --eval <expression>", name, args)
				}
				return []byte(tt.output), tt.err
			}
			if result := verifyEmacsLSP(context.Background()); result.status != tt.wantStatus {
				t.Errorf("status = %q, want %q", result.status, tt.wantStatus)
			}
		})
	}
}

func TestLSPDoctorVerificationHint_WritesToStandardError(t *testing.T) {
	oldCmd := currentCmd
	defer func() { currentCmd = oldCmd }()

	command := &cobra.Command{Use: "doctor"}
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	command.SetOut(stdout)
	command.SetErr(stderr)
	currentCmd = command

	lspDoctorVerificationHint(logging.DefaultTheme())
	if stdout.Len() != 0 {
		t.Errorf("stdout = %q, want no verification hint", stdout.String())
	}
	if !strings.Contains(stderr.String(), "markata-go lsp doctor") {
		t.Errorf("stderr = %q, want verification command", stderr.String())
	}
}

func TestLSPDoctorStatusColor_UsesSemanticThemeColors(t *testing.T) {
	theme := logging.Theme{
		Component: "#111111",
		Success:   "#222222",
		Warning:   "#333333",
		Error:     "#444444",
	}

	for _, tt := range []struct {
		status string
		want   string
	}{
		{"PASS", "#222222"},
		{"INFO", "#111111"},
		{"WARN", "#333333"},
		{"FAIL", "#444444"},
	} {
		if got := lspDoctorStatusColor(theme, tt.status); got != tt.want {
			t.Errorf("lspDoctorStatusColor(%q) = %q, want %q", tt.status, got, tt.want)
		}
	}
}
