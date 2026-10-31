package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/WaylonWalker/markata/pkg/config"
	"github.com/WaylonWalker/markata/pkg/logging"
	"github.com/WaylonWalker/markata/pkg/lsp"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/palettes"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/spf13/cobra"
)

// lspCmd represents the lsp command.
var lspCmd = &cobra.Command{
	Use:   "lsp",
	Short: "Start the Language Server Protocol server",
	Long: `Start the markata-go LSP server for IDE integration.

The LSP server provides IDE features for markdown files with wikilink support:
  - Autocomplete for [[wikilinks]] - type [[ to get suggestions
  - Diagnostics for broken wikilinks - warnings for links to missing posts
  - Hover information - see post title and description on hover
  - Go to definition - Ctrl+click to navigate to linked posts

The server communicates over stdin/stdout using the Language Server Protocol.

Set up an editor with "markata-go lsp setup --editor <editor>".
Check your project and LSP installation with "markata-go lsp doctor".`,
	RunE: runLSPCommand,
}

var lspBrowserWorkspace string
var lspSetupEditor string
var lspSkipEditorVerification bool

var lspDoctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check local LSP prerequisites",
	Long: `Check whether this markata-go binary can serve LSP requests and whether
the current project has a usable markata-go configuration and verifies installed
editors that support headless checks. This can load user startup and plugin code.
Pass --no-verify-editor to detect editors without loading their configuration.`,
	RunE: runLSPDoctor,
}

var lspSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Print editor LSP setup guidance",
	Long: `Print copy-paste setup snippets for installed editors, or for a selected
editor. The command never edits editor configuration because editor clients and
configuration layouts vary.

Examples:
  markata-go lsp setup
  markata-go lsp setup --editor neovim
  markata-go lsp setup --editor helix`,
	RunE: runLSPSetup,
}

func init() {
	rootCmd.AddCommand(lspCmd)
	lspCmd.Flags().StringVar(&lspBrowserWorkspace, "browser-workspace", "", "confine a browser bridge to this workspace and its in-memory draft")
	lspCmd.AddCommand(lspDoctorCmd, lspSetupCmd)
	lspSetupCmd.Flags().StringVarP(&lspSetupEditor, "editor", "e", "", "editor: generic, neovim, helix, emacs, zed, or vscode (default: all installed editors)")
	lspDoctorCmd.Flags().BoolVar(&lspSkipEditorVerification, "no-verify-editor", false, "detect installed editors without loading their configuration")
}

func runLSPCommand(_ *cobra.Command, _ []string) error {
	// Setup logging to stderr (stdout is used for LSP communication)
	var logger *log.Logger
	if verbose {
		logger = log.New(os.Stderr, "[markata-lsp] ", log.LstdFlags|log.Lshortfile)
	} else {
		logger = log.New(os.Stderr, "[markata-lsp] ", log.LstdFlags)
	}

	// Create context with cancellation for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		logger.Println("Received shutdown signal")
		cancel()
	}()

	// Create and run the LSP server
	srv := lsp.New(logger)
	if lspBrowserWorkspace != "" {
		var err error
		srv, err = lsp.NewBrowserSession(logger, lspBrowserWorkspace)
		if err != nil {
			return err
		}
	}
	if err := srv.Run(ctx, os.Stdin, os.Stdout); err != nil {
		return fmt.Errorf("LSP server error: %w", err)
	}

	return nil
}

func runLSPDoctor(cmd *cobra.Command, _ []string) error {
	currentCmd = cmd
	doctorTheme := logging.DefaultTheme()

	cfg, _, configPaths, err := loadManagerConfig(cfgFile)
	if err != nil {
		failure := fmt.Errorf("configuration is invalid: %w", err)
		lspDoctorHeader(doctorTheme, "markata-go LSP diagnostics")
		lspDoctorLine(doctorTheme, "PASS", "LSP command: this binary can run 'markata-go lsp'.")
		lspDoctorLine(doctorTheme, "FAIL", failure.Error())
		return failure
	}
	validationErrors, _ := config.SplitErrorsAndWarnings(config.ValidateConfig(cfg))
	if len(validationErrors) > 0 {
		failure := fmt.Errorf("configuration is invalid: %w", validationErrors[0])
		lspDoctorHeader(doctorTheme, "markata-go LSP diagnostics")
		lspDoctorLine(doctorTheme, "PASS", "LSP command: this binary can run 'markata-go lsp'.")
		lspDoctorLine(doctorTheme, "FAIL", failure.Error())
		return failure
	}

	doctorTheme = resolveLSPDoctorTheme(cfg, configPaths)
	lspDoctorHeader(doctorTheme, "markata-go LSP diagnostics")
	lspDoctorLine(doctorTheme, "PASS", "LSP command: this binary can run 'markata-go lsp'.")
	if len(configPaths) == 0 {
		lspDoctorLine(doctorTheme, "INFO", "Configuration: none found; the LSP will use the workspace and available Markdown files.")
		lspDoctorLine(doctorTheme, "INFO", "Workspace root: supplied by your editor.")
	} else {
		lspDoctorLine(doctorTheme, "PASS", fmt.Sprintf("Configuration: %s is valid.", strings.Join(configPaths, ", ")))
		for _, configPath := range configPaths {
			if strings.HasSuffix(configPath, ".yml") || strings.HasSuffix(configPath, ".json") {
				lspDoctorLine(doctorTheme, "WARN", fmt.Sprintf("Configuration: %s is valid for builds, but LSP mention indexing reads only TOML and .yaml files.", configPath))
			}
		}
	}

	if runLSPEditorDiagnostics(!lspSkipEditorVerification, outWriter(), doctorTheme) && lspSkipEditorVerification {
		lspDoctorVerificationHint(doctorTheme)
	}
	lspDoctorLine(doctorTheme, "PASS", "LSP prerequisites are ready.")
	return nil
}

func resolveLSPDoctorTheme(cfg *models.Config, configPaths []string) logging.Theme {
	return resolveCLITheme(cfg, configPaths)
}

func lspDoctorHeader(theme logging.Theme, title string) {
	outln(colorizeOutput(title, theme.Component))
}

func lspDoctorLine(theme logging.Theme, status, message string) {
	outlnf("%s  %s", colorizeOutput(status, lspDoctorStatusColor(theme, status)), message)
}

func lspDoctorVerificationHint(theme logging.Theme) {
	errln(colorize("Editor verification was skipped. Run the default doctor command to verify detected editors:", theme.Warning, colorEnabledFor(errorOutputIsTerminal())))
	errln("  markata-go lsp doctor")
	errln("  This may load editor startup and plugin code.")
}

type lspEditor struct {
	name             string
	executable       string
	supportsHeadless bool
	verify           func(context.Context) lspEditorVerification
}

type lspEditorVerification struct {
	status  lspEditorStatus
	message string
}

type lspEditorStatus string

const (
	lspEditorConfigured   lspEditorStatus = "configured"
	lspEditorUnconfigured lspEditorStatus = "unconfigured"
	lspEditorInconclusive lspEditorStatus = "inconclusive"
)

var lspEditors = []lspEditor{
	{name: "Neovim", executable: "nvim", supportsHeadless: true, verify: verifyNeovimLSP},
	{name: "Helix", executable: "hx", supportsHeadless: true, verify: verifyHelixLSP},
	{name: "Emacs", executable: "emacs", supportsHeadless: true, verify: verifyEmacsLSP},
	{name: "VS Code", executable: "code"},
	{name: "Cursor", executable: "cursor"},
	{name: "Zed", executable: "zed"},
}

var lspLookPath = exec.LookPath
var lspRunCommand = func(ctx context.Context, name string, args ...string) ([]byte, error) {
	return exec.CommandContext(ctx, name, args...).CombinedOutput()
}

func runLSPEditorDiagnostics(verify bool, outputWriter io.Writer, theme logging.Theme) bool {
	installed := installedLSPEditors(lspEditors, lspLookPath)
	if len(installed) == 0 {
		lspDoctorWriterLine(outputWriter, theme, "INFO", "Editors: no supported editor executable found on PATH.")
		return false
	}

	hasHeadlessEditor := false
	for _, editor := range installed {
		hasHeadlessEditor = hasHeadlessEditor || editor.supportsHeadless
		if !verify {
			if editor.supportsHeadless {
				lspDoctorWriterLine(outputWriter, theme, "INFO", fmt.Sprintf("Editor: %s is installed. Run doctor without --no-verify-editor to validate its LSP configuration.", editor.name))
			} else {
				lspDoctorWriterLine(outputWriter, theme, "INFO", fmt.Sprintf("Editor: %s is installed; runtime verification is not supported.", editor.name))
			}
			continue
		}
		if !editor.supportsHeadless {
			lspDoctorWriterLine(outputWriter, theme, "INFO", fmt.Sprintf("Editor: %s is installed; runtime verification is not supported.", editor.name))
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		result := editor.verify(ctx)
		cancel()
		switch result.status {
		case lspEditorConfigured:
			lspDoctorWriterLine(outputWriter, theme, "PASS", fmt.Sprintf("Editor: %s configured: %s", editor.name, result.message))
		case lspEditorUnconfigured:
			lspDoctorWriterLine(outputWriter, theme, "WARN", fmt.Sprintf("Editor: %s unconfigured: %s", editor.name, result.message))
		default:
			lspDoctorWriterLine(outputWriter, theme, "WARN", fmt.Sprintf("Editor: %s inconclusive: %s", editor.name, result.message))
		}
	}
	return hasHeadlessEditor
}

func lspDoctorWriterLine(writer io.Writer, theme logging.Theme, status, message string) {
	_, _ = fmt.Fprintf(writer, "%s  %s\n", colorizeOutput(status, lspDoctorStatusColor(theme, status)), message)
}

func lspDoctorStatusColor(theme logging.Theme, status string) string {
	switch status {
	case "PASS":
		return theme.Success
	case "INFO":
		return theme.Component
	case "WARN":
		return theme.Warning
	case "FAIL":
		return theme.Error
	default:
		return theme.Component
	}
}

func installedLSPEditors(editors []lspEditor, lookPath func(string) (string, error)) []lspEditor {
	installed := make([]lspEditor, 0, len(editors))
	for _, editor := range editors {
		if _, err := lookPath(editor.executable); err == nil {
			installed = append(installed, editor)
		}
	}
	return installed
}

func verifyNeovimLSP(ctx context.Context) lspEditorVerification {
	const script = `local c=vim.lsp.config["markata"]; if not c then io.stderr:write("MARKATA_NOT_CONFIGURED\n"); vim.cmd("cquit 2") end; local cmd=c.cmd or {}; local filetypes=c.filetypes or {}; local markdown=false; for _, ft in ipairs(filetypes) do if ft == "markdown" then markdown=true end end; if cmd[1] ~= "markata-go" or cmd[2] ~= "lsp" or not markdown then io.stderr:write("MARKATA_CONFIG_INVALID\n"); vim.cmd("cquit 2") end; print("MARKATA_CONFIGURED"); vim.cmd("qa")`
	output, err := lspRunCommand(ctx, "nvim", "--headless", "+lua "+script)
	if err != nil {
		if strings.Contains(string(output), "MARKATA_NOT_CONFIGURED") || strings.Contains(string(output), "MARKATA_CONFIG_INVALID") {
			return lspEditorVerification{status: lspEditorUnconfigured, message: "run 'markata-go lsp setup --editor neovim'."}
		}
		return lspEditorVerification{status: lspEditorInconclusive, message: "could not load a Neovim LSP configuration."}
	}
	if strings.Contains(string(output), "MARKATA_CONFIGURED") {
		return lspEditorVerification{status: lspEditorConfigured, message: "markata-go LSP is registered for Markdown."}
	}
	return lspEditorVerification{status: lspEditorInconclusive, message: "did not return a recognized verification result."}
}

func verifyHelixLSP(ctx context.Context) lspEditorVerification {
	output, err := lspRunCommand(ctx, "hx", "--health", "markdown")
	if err != nil {
		return lspEditorVerification{status: lspEditorInconclusive, message: "could not run 'hx --health markdown'."}
	}
	if strings.Contains(string(output), "markata-go") {
		return lspEditorVerification{status: lspEditorConfigured, message: "health check recognizes markata-go for Markdown."}
	}
	return lspEditorVerification{status: lspEditorUnconfigured, message: "health check does not recognize markata-go for Markdown; run 'markata-go lsp setup --editor helix'."}
}

func verifyEmacsLSP(ctx context.Context) lspEditorVerification {
	const expression = `(progn (require 'eglot nil t) (if (boundp 'eglot-server-programs) (princ (prin1-to-string eglot-server-programs)) (princ "NO_EGLOT")))`
	output, err := lspRunCommand(ctx, "emacs", "--batch", "--eval", expression)
	if err != nil {
		return lspEditorVerification{status: lspEditorInconclusive, message: "could not load Eglot configuration in batch mode."}
	}
	if strings.Contains(string(output), "markata-go") && strings.Contains(string(output), "lsp") {
		return lspEditorVerification{status: lspEditorConfigured, message: "Eglot configuration references markata-go lsp."}
	}
	return lspEditorVerification{status: lspEditorUnconfigured, message: "does not expose a recognized Eglot markata-go LSP configuration."}
}

func runLSPSetup(cmd *cobra.Command, _ []string) error {
	currentCmd = cmd
	if lspSetupEditor != "" {
		_, ok := lspSetupSnippets[lspSetupEditor]
		if !ok {
			return newUsageError(fmt.Errorf("unsupported editor %q; supported editors: %s", lspSetupEditor, strings.Join(supportedLSPEditors(), ", ")))
		}
		outText(renderLSPSetupSnippets([]string{lspSetupEditor}))
		return nil
	}

	setupEditors := installedLSPSetupEditors(installedLSPEditors(lspEditors, lspLookPath))
	if len(setupEditors) == 0 {
		outText(lspSetupSnippets["generic"])
		return nil
	}

	outText(renderLSPSetupSnippets(setupEditors))
	return nil
}

func renderLSPSetupSnippets(editors []string) string {
	theme, chromaStyle := resolveLSPSetupRendering()
	if len(editors) == 1 {
		editor := editors[0]
		return highlightLSPSetupSnippet(lspSetupSnippets[editor], lspSetupLanguages[editor], chromaStyle)
	}

	sections := make([]string, 0, len(editors))
	for _, editor := range editors {
		header := cliSectionRule(lspSetupEditorNames[editor])
		if colorEnabledOnOutput() {
			header = colorizeOutput(header, theme.Component)
		}
		sections = append(sections, header+"\n"+highlightLSPSetupSnippet(lspSetupSnippets[editor], lspSetupLanguages[editor], chromaStyle))
	}
	return strings.Join(sections, "\n\n")
}

func resolveLSPSetupRendering() (theme logging.Theme, chromaStyle string) {
	theme = logging.DefaultTheme()
	cfg, _, configPaths, err := loadManagerConfig(cfgFile)
	if err == nil && cfg != nil {
		extra := make(map[string]any)
		if len(configPaths) > 0 {
			extra["config_path"] = configPaths[0]
		}
		if palette, ok := loadLoggerPalette(cfg.Theme, extra); ok {
			theme = logging.ThemeFromPalette(palette)
			if style := palettes.ChromaTheme(palette.Name); style != "" {
				return theme, style
			}
			return theme, palettes.ChromaThemeForVariant(palette.Variant)
		}
	}
	return theme, palettes.DefaultChromaThemeDark
}

func highlightLSPSetupSnippet(snippet, language, chromaStyle string) string {
	if !colorEnabledOnOutput() {
		return snippet
	}
	lexer := lexers.Get(language)
	if lexer == nil {
		lexer = lexers.Get("plaintext")
	}
	style := styles.Get(chromaStyle)
	if style == nil {
		style = styles.Fallback
	}
	iterator, err := lexer.Tokenise(nil, snippet)
	if err != nil {
		return snippet
	}
	var output strings.Builder
	formatter := formatters.Get("terminal16m")
	if formatter == nil {
		return snippet
	}
	if err := formatter.Format(&output, style, iterator); err != nil {
		return snippet
	}
	return output.String()
}

func installedLSPSetupEditors(editors []lspEditor) []string {
	setupEditorByExecutable := map[string]string{
		"nvim":   "neovim",
		"hx":     "helix",
		"emacs":  "emacs",
		"code":   "vscode",
		"cursor": "vscode",
		"zed":    "zed",
	}
	seen := make(map[string]bool)
	setupEditors := make([]string, 0, len(editors))
	for _, editor := range editors {
		setupEditor, ok := setupEditorByExecutable[editor.executable]
		if !ok || seen[setupEditor] {
			continue
		}
		seen[setupEditor] = true
		setupEditors = append(setupEditors, setupEditor)
	}
	return setupEditors
}

func supportedLSPEditors() []string {
	editors := make([]string, 0, len(lspSetupSnippets))
	for editor := range lspSetupSnippets {
		editors = append(editors, editor)
	}
	sort.Strings(editors)
	return editors
}

var lspSetupSnippets = map[string]string{
	"generic": `# markata-go LSP generic client configuration
#
# Start this command over stdio for Markdown files. Set the workspace root to
# the directory containing markata-go.toml (or your site's content directory).
command: markata-go lsp
filetypes: markdown
root: markata-go configuration directory
`,
	"neovim": `-- Neovim 0.11+ (init.lua)
-- Add this before opening Markdown files. The root marker makes each site a workspace.
vim.lsp.config("markata", {
  cmd = { "markata-go", "lsp" },
  filetypes = { "markdown" },
  root_markers = { "markata-go.toml", "markata-go.yaml", "markata-go.yml", "markata-go.json", ".git" },
})
vim.lsp.enable("markata")
`,
	"helix": `# Helix languages.toml (usually ~/.config/helix/languages.toml)
[language-server.markata-go]
command = "markata-go"
args = ["lsp"]

[[language]]
name = "markdown"
language-servers = ["markata-go"]
`,
	"emacs": `;; Emacs with Eglot (init.el)
;; Install/enable eglot, then register markata-go for Markdown buffers.
(with-eval-after-load 'eglot
  (add-to-list 'eglot-server-programs
               '((markdown-mode gfm-mode) . ("markata-go" "lsp"))))
`,
	"zed": `// Zed settings.json
// Zed configuration evolves quickly. Register the command with your installed
// Markdown LSP client/extension using these stable values:
// command: "markata-go"
// arguments: ["lsp"]
// language: Markdown
// workspace root: directory containing markata-go.toml
//
// See the Zed language-server settings for your installed version, then run:
// markata-go lsp doctor
`,
	"vscode": `// VS Code / Cursor
// VS Code does not let settings.json register an arbitrary stdio LSP server by
// itself. Install an LSP-client extension or a Markata-aware extension, then
// configure its server launch values as:
//
// command: "markata-go"
// args: ["lsp"]
// filetypes/language: ["markdown"]
// workspace root: ${workspaceFolder} (the directory with markata-go.toml)
//
// Do not add a task as a substitute: tasks run once and are not LSP clients.
// Ask an agent to adapt these values to the LSP extension already installed in
// this workspace, then run "markata-go lsp doctor".
`,
}

var lspSetupEditorNames = map[string]string{
	"emacs":   "Emacs",
	"generic": "Generic LSP client",
	"helix":   "Helix",
	"neovim":  "Neovim",
	"vscode":  "VS Code / Cursor",
	"zed":     "Zed",
}

var lspSetupLanguages = map[string]string{
	"emacs":   "emacs-lisp",
	"generic": "yaml",
	"helix":   "toml",
	"neovim":  "lua",
	"vscode":  "json",
	"zed":     "json",
}
