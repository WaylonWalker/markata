package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	cdruntime "github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"

	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/palettes"
	"github.com/WaylonWalker/markata/pkg/runtimeenv"
	gmermaid "go.abhg.dev/goldmark/mermaid"
	"go.abhg.dev/goldmark/mermaid/mermaidcdp"
)

// mermaidRenderer is an interface for rendering Mermaid diagrams in different modes
type mermaidRenderer interface {
	// render processes a mermaid diagram code and returns the rendered output
	// For client mode: returns the diagram code wrapped in mermaid markers
	// For cli/chromium modes: returns the SVG output
	render(diagramCode string) (string, error)
	// close cleans up any resources (e.g., browser connections)
	close() error
}

// mermaidPaletteColors holds resolved palette colors for Mermaid theming.
// These are resolved at build time from the site's palette and passed to
// mermaid.initialize() as themeVariables for consistent diagram styling.
type mermaidPaletteColors struct {
	Background string // bg-primary -> --color-background
	TextColor  string // text-primary -> --color-text
	Accent     string // accent -> --color-border fallback
	Surface    string // bg-surface -> --color-surface
}

// resolvePaletteColors loads and resolves palette colors for mermaid theming.
// Returns nil if no palette is configured or colors cannot be resolved.
func resolvePaletteColors(extra map[string]interface{}) *mermaidPaletteColors {
	if extra == nil {
		return nil
	}

	// Extract palette name from config.Extra["theme"]
	// The theme config may be a map[string]interface{} (raw TOML) or
	// models.ThemeConfig (if already parsed by another plugin).
	var paletteName string
	switch theme := extra["theme"].(type) {
	case map[string]interface{}:
		if p, ok := theme["palette"].(string); ok {
			paletteName = p
		}
	case models.ThemeConfig:
		paletteName = theme.Palette
	case *models.ThemeConfig:
		if theme != nil {
			paletteName = theme.Palette
		}
	}

	if paletteName == "" {
		return nil
	}

	loader := palettes.NewLoader()
	palette, err := loader.Load(paletteName)
	if err != nil {
		log.Printf("[mermaid] could not load palette %q: %v", paletteName, err)
		return nil
	}

	colors := &mermaidPaletteColors{
		Background: palette.Resolve("bg-primary"),
		TextColor:  palette.Resolve("text-primary"),
		Accent:     palette.Resolve("accent"),
		Surface:    palette.Resolve("bg-surface"),
	}

	// Only return if we got at least the accent color
	if colors.Accent == "" {
		log.Printf("[mermaid] palette %q has no accent color, using default theme", paletteName)
		return nil
	}

	log.Printf("[mermaid] resolved palette %q colors for theming (accent=%s, bg=%s, text=%s)",
		paletteName, colors.Accent, colors.Background, colors.TextColor)
	return colors
}

// Default fallback colors for mermaid theme variables when palette values are missing.
const (
	defaultMermaidBg      = "#ffffff"
	defaultMermaidText    = "#1f2937"
	defaultMermaidSurface = "#f9fafb"
)

// themeVariablesJSON returns the mermaid themeVariables object as a JSON string.
func (c *mermaidPaletteColors) themeVariablesJSON() string {
	bg := c.Background
	if bg == "" {
		bg = defaultMermaidBg
	}
	text := c.TextColor
	if text == "" {
		text = defaultMermaidText
	}
	surface := c.Surface
	if surface == "" {
		surface = defaultMermaidSurface
	}

	vars := map[string]interface{}{
		"background":          bg,
		"primaryColor":        surface,
		"primaryTextColor":    text,
		"primaryBorderColor":  c.Accent,
		"lineColor":           c.Accent,
		"textColor":           text,
		"nodeBkg":             surface,
		"nodeBorder":          c.Accent,
		"nodeTextColor":       text,
		"fontSize":            "16px",
		"clusterBkg":          bg,
		"clusterBorder":       c.Accent,
		"clusterTextColor":    text,
		"titleColor":          text,
		"edgeLabelBackground": surface,
	}

	b, err := json.Marshal(vars)
	if err != nil {
		log.Printf("[mermaid] failed to marshal theme variables: %v", err)
		return "{}"
	}
	return string(b)
}

// clientRenderer renders diagrams client-side via JavaScript
type clientRenderer struct {
	config models.MermaidConfig
}

func (r *clientRenderer) render(diagramCode string) (string, error) {
	// Client-side rendering: just return the diagram code, the script will render it
	return diagramCode, nil
}

func (r *clientRenderer) close() error {
	return nil
}

// cliRenderer renders diagrams using the mmdc CLI tool
type cliRenderer struct {
	config models.MermaidConfig
}

func (r *cliRenderer) render(diagramCode string) (string, error) {
	mmdc := r.config.CLIConfig.MMDCPath
	if mmdc == "" {
		var err error
		mmdc, err = exec.LookPath("mmdc")
		if err != nil {
			return "", fmt.Errorf("mmdc binary not found: %w", err)
		}
	}

	// Create temporary directory for diagram files
	tmpDir, err := os.MkdirTemp("", "markata-mermaid-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Write diagram code to temp file
	inputFile := filepath.Join(tmpDir, "diagram.mmd")
	outputFile := filepath.Join(tmpDir, "diagram.svg")

	if err := os.WriteFile(inputFile, []byte(diagramCode), 0o600); err != nil {
		return "", fmt.Errorf("failed to write diagram file: %w", err)
	}

	// Run mmdc to render the diagram
	cmd := exec.Command(mmdc, "-i", inputFile, "-o", outputFile)
	if r.config.CLIConfig.ExtraArgs != "" {
		cmd.Args = append(cmd.Args, strings.Fields(r.config.CLIConfig.ExtraArgs)...)
	}

	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("mmdc rendering failed: %w (output: %s)", err, string(output))
	}

	// Read the generated SVG
	svgBytes, err := os.ReadFile(outputFile)
	if err != nil {
		return "", fmt.Errorf("failed to read generated SVG: %w", err)
	}

	return string(svgBytes), nil
}

func (r *cliRenderer) close() error {
	return nil
}

// chromiumRenderer renders diagrams using Chrome DevTools Protocol.
// When paletteColors is set, it uses a custom chromedp flow that passes
// themeVariables to mermaid.initialize() for palette-aware theming.
// Otherwise, it falls back to mermaidcdp.Compiler for default theming.
type chromiumRenderer struct {
	config        models.MermaidConfig
	paletteColors *mermaidPaletteColors
	once          sync.Once
	initErr       error
	semaphore     chan struct{} // limits concurrent renders

	// Custom chromedp compiler (when palette colors are set)
	browserCtx    context.Context //nolint:containedctx // chromedp requires storing context
	browserCancel context.CancelFunc

	// Fallback mermaidcdp compiler (when no palette colors)
	compiler *mermaidcdp.Compiler
}

// mermaidJSVersion is the MermaidJS version to download for chromium rendering.
const mermaidJSVersion = "10.9.5"

// getMermaidJSCacheDir returns the cache directory for MermaidJS source files,
// following XDG conventions (~/.cache/markata-go/mermaid/).
func getMermaidJSCacheDir() (string, error) {
	cacheBase := os.Getenv("XDG_CACHE_HOME")
	if cacheBase == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %w", err)
		}
		cacheBase = filepath.Join(homeDir, ".cache")
	}
	cacheDir := filepath.Join(cacheBase, "markata-go", "mermaid")
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create mermaid cache directory: %w", err)
	}
	return cacheDir, nil
}

func loadBundledMermaidJSSource(version string) (string, error) {
	bundledDir := runtimeenv.BundledMermaidDir()
	if bundledDir == "" {
		return "", os.ErrNotExist
	}
	cacheFile := filepath.Join(bundledDir, fmt.Sprintf("mermaid-v%s.min.js", version))
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return "", err
	}
	if len(data) == 0 {
		return "", fmt.Errorf("bundled MermaidJS source is empty: %s", cacheFile)
	}
	log.Printf("[mermaid] chromium: loaded MermaidJS from bundled cache (%d bytes)", len(data))
	return string(data), nil
}

// loadOrDownloadJSSource loads the MermaidJS source from the local cache,
// downloading it if not already cached. The cached file is stored at
// ~/.cache/markata-go/mermaid/mermaid-v{version}.min.js
func loadOrDownloadJSSource(ctx context.Context, version string) (string, error) {
	cacheDir, err := getMermaidJSCacheDir()
	if err != nil {
		if bundledSource, bundledErr := loadBundledMermaidJSSource(version); bundledErr == nil {
			return bundledSource, nil
		}
		if runtimeenv.OfflineEnabled() {
			return "", fmt.Errorf("failed to access MermaidJS cache in offline mode: %w", err)
		}
		// Cache dir unavailable; fall back to direct download
		log.Printf("[mermaid] chromium: cache unavailable (%v), downloading directly", err)
		return mermaidcdp.DownloadJSSource(ctx, version)
	}

	cacheFile := filepath.Join(cacheDir, fmt.Sprintf("mermaid-v%s.min.js", version))

	// Try loading from cache first
	if data, err := os.ReadFile(cacheFile); err == nil && len(data) > 0 {
		log.Printf("[mermaid] chromium: loaded MermaidJS from cache (%d bytes)", len(data))
		return string(data), nil
	}
	if bundledSource, err := loadBundledMermaidJSSource(version); err == nil {
		return bundledSource, nil
	}
	if runtimeenv.OfflineEnabled() {
		return "", fmt.Errorf("mermaid JS source mermaid-v%s.min.js is not available in cache or bundled assets while offline", version)
	}

	// Download and cache
	log.Println("[mermaid] chromium: downloading MermaidJS source...")
	jsSource, err := mermaidcdp.DownloadJSSource(ctx, version)
	if err != nil {
		return "", err
	}
	log.Printf("[mermaid] chromium: downloaded MermaidJS (%d bytes)", len(jsSource))

	// Write to cache (best-effort, don't fail the build if caching fails)
	if writeErr := os.WriteFile(cacheFile, []byte(jsSource), 0o600); writeErr != nil {
		log.Printf("[mermaid] chromium: warning: failed to cache MermaidJS source: %v", writeErr)
	} else {
		log.Printf("[mermaid] chromium: cached MermaidJS source at %s", cacheFile)
	}

	return jsSource, nil
}

// _renderSVGJS is the JavaScript helper injected into the headless browser
// to render mermaid diagrams. Equivalent to mermaidcdp's extras.js.
const _renderSVGJS = `async function renderSVG(src) {
	const { svg } = await mermaid.render('mermaid', src);
	return svg;
}`

// ensureCompilerCustom sets up a custom chromedp flow with themeVariables support.
func (r *chromiumRenderer) ensureCompilerCustom(jsSource string) {
	noSandbox := r.config.ChromiumConfig != nil && r.config.ChromiumConfig.NoSandbox

	ctx := context.Background()
	if noSandbox {
		execOpts := make([]chromedp.ExecAllocatorOption, 0, len(chromedp.DefaultExecAllocatorOptions)+1)
		execOpts = append(execOpts, chromedp.DefaultExecAllocatorOptions[:]...)
		execOpts = append(execOpts, chromedp.NoSandbox)

		var allocCancel context.CancelFunc
		ctx, allocCancel = chromedp.NewExecAllocator(ctx, execOpts...)
		// Store cancel so we clean up on close
		defer func() {
			if r.initErr != nil {
				allocCancel()
			}
		}()
		// Chain the alloc cancel into browserCancel
		origCancel := allocCancel
		defer func() {
			if r.initErr == nil {
				// Wrap both cancels
				innerCancel := r.browserCancel
				r.browserCancel = func() {
					innerCancel()
					origCancel()
				}
			}
		}()
	}

	var browserCancel context.CancelFunc
	ctx, browserCancel = chromedp.NewContext(ctx)
	r.browserCtx = ctx
	r.browserCancel = browserCancel

	// Load MermaidJS source
	var ready *cdruntime.RemoteObject
	if err := chromedp.Run(ctx, chromedp.Evaluate(jsSource, &ready)); err != nil {
		r.initErr = fmt.Errorf("failed to load MermaidJS in browser: %w", err)
		browserCancel()
		return
	}

	// Inject renderSVG helper
	ready = nil
	if err := chromedp.Run(ctx, chromedp.Evaluate(_renderSVGJS, &ready)); err != nil {
		r.initErr = fmt.Errorf("failed to inject renderSVG helper: %w", err)
		browserCancel()
		return
	}

	// Build mermaid.initialize() call with themeVariables
	initJS := fmt.Sprintf(
		`mermaid.initialize({startOnLoad: false, theme: 'base', themeVariables: %s, flowchart: {nodeSpacing: 60, rankSpacing: 90, padding: 12}})`,
		r.paletteColors.themeVariablesJSON(),
	)

	ready = nil
	if err := chromedp.Run(ctx, chromedp.Evaluate(initJS, &ready)); err != nil {
		r.initErr = fmt.Errorf("failed to initialize mermaid with palette theme: %w", err)
		browserCancel()
		return
	}

	log.Println("[mermaid] chromium: initialized with palette-aware theming")
}

// ensureCompilerFallback sets up the standard mermaidcdp.Compiler (no palette).
func (r *chromiumRenderer) ensureCompilerFallback(jsSource string) {
	cfg := &mermaidcdp.Config{
		JSSource:  jsSource,
		Theme:     r.config.Theme,
		NoSandbox: r.config.ChromiumConfig != nil && r.config.ChromiumConfig.NoSandbox,
	}

	var err error
	r.compiler, err = mermaidcdp.New(cfg)
	if err != nil {
		r.initErr = fmt.Errorf("failed to start chromium compiler: %w", err)
		log.Printf("[mermaid] chromium: browser launch failed: %v", err)
		return
	}
}

func (r *chromiumRenderer) ensureCompiler() error {
	r.once.Do(func() {
		dlCtx, dlCancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer dlCancel()

		jsSource, err := loadOrDownloadJSSource(dlCtx, mermaidJSVersion)
		if err != nil {
			r.initErr = fmt.Errorf("failed to obtain MermaidJS source: %w", err)
			return
		}

		log.Println("[mermaid] chromium: launching headless browser...")

		if r.paletteColors != nil {
			r.ensureCompilerCustom(jsSource)
		} else {
			r.ensureCompilerFallback(jsSource)
		}

		if r.initErr != nil {
			return
		}

		log.Println("[mermaid] chromium: browser ready")

		// Initialize the semaphore for concurrent renders
		maxConcurrent := 4
		if r.config.ChromiumConfig != nil && r.config.ChromiumConfig.MaxConcurrent > 0 {
			maxConcurrent = r.config.ChromiumConfig.MaxConcurrent
		}
		r.semaphore = make(chan struct{}, maxConcurrent)
	})
	return r.initErr
}

func (r *chromiumRenderer) render(diagramCode string) (string, error) {
	if err := r.ensureCompiler(); err != nil {
		return "", err
	}

	// Acquire semaphore slot
	r.semaphore <- struct{}{}
	defer func() { <-r.semaphore }()

	timeout := 30
	if r.config.ChromiumConfig != nil && r.config.ChromiumConfig.Timeout > 0 {
		timeout = r.config.ChromiumConfig.Timeout
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	// Custom flow: render directly via chromedp
	if r.browserCtx != nil {
		return r.renderCustom(ctx, diagramCode)
	}

	// Fallback: render via mermaidcdp.Compiler
	return r.renderFallback(ctx, diagramCode)
}

// renderCustom renders a diagram using the custom chromedp browser context.
func (r *chromiumRenderer) renderCustom(timeoutCtx context.Context, diagramCode string) (string, error) {
	// Build the renderSVG() call with JSON-encoded source
	var script strings.Builder
	script.WriteString("renderSVG(")
	if err := json.NewEncoder(&script).Encode(diagramCode); err != nil {
		return "", fmt.Errorf("failed to encode diagram source: %w", err)
	}
	script.WriteString(")")

	var result string
	render := chromedp.Evaluate(
		script.String(),
		&result,
		func(p *cdruntime.EvaluateParams) *cdruntime.EvaluateParams {
			return p.WithAwaitPromise(true)
		},
	)

	// Merge browser context (for chromedp state) with timeout context.
	// chromedp needs its context for browser access, but we want the
	// timeout from timeoutCtx.
	mergedCtx, mergeCancel := mergeContextLifetime(r.browserCtx, timeoutCtx)
	defer mergeCancel()

	if err := chromedp.Run(mergedCtx, render); err != nil {
		return "", fmt.Errorf("mermaid rendering failed: %w", err)
	}

	svgOutput := html.UnescapeString(result)
	return strings.TrimSpace(svgOutput), nil
}

// renderFallback renders a diagram via the mermaidcdp.Compiler.
func (r *chromiumRenderer) renderFallback(ctx context.Context, diagramCode string) (string, error) {
	resp, err := r.compiler.Compile(ctx, &gmermaid.CompileRequest{Source: diagramCode})
	if err != nil {
		return "", fmt.Errorf("mermaidcdp rendering failed: %w", err)
	}

	svgOutput := html.UnescapeString(resp.SVG)
	return strings.TrimSpace(svgOutput), nil
}

func (r *chromiumRenderer) close() error {
	if r.browserCancel != nil {
		log.Println("[mermaid] chromium: closing browser (custom)")
		r.browserCancel()
		return nil
	}
	if r.compiler != nil {
		log.Println("[mermaid] chromium: closing browser")
		return r.compiler.Close()
	}
	return nil
}

// mergeContextLifetime creates a child of parentCtx that is also canceled
// when timeCtx is done. This allows chromedp operations to use the browser
// context (parentCtx) while respecting a render timeout (timeCtx).
func mergeContextLifetime(parentCtx, timeCtx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancelCause(parentCtx)
	stop := context.AfterFunc(timeCtx, func() {
		cancel(context.Cause(timeCtx))
	})
	return ctx, func() {
		stop()
		cancel(context.Canceled)
	}
}

// newMermaidRenderer creates the appropriate renderer based on the config mode
func newMermaidRenderer(config models.MermaidConfig, paletteColors *mermaidPaletteColors) (mermaidRenderer, error) {
	switch config.Mode {
	case "client":
		return &clientRenderer{config: config}, nil

	case "cli":
		// Validate CLI dependencies
		info := checkCLIDependency(config.CLIConfig.MMDCPath)
		if !info.IsInstalled {
			err := models.NewMermaidRenderError("", config.Mode, "mmdc binary not found", nil)
			err.Suggestion = info.InstallInstructions + "\n\n" + info.FallbackSuggestion
			return nil, err
		}
		return &cliRenderer{config: config}, nil

	case "chromium":
		// Validate Chromium dependencies
		info := checkChromiumDependency(config.ChromiumConfig.BrowserPath)
		if !info.IsInstalled {
			err := models.NewMermaidRenderError("", config.Mode, "browser not found", nil)
			err.Suggestion = info.InstallInstructions + "\n\n" + info.FallbackSuggestion
			return nil, err
		}
		return &chromiumRenderer{config: config, paletteColors: paletteColors}, nil

	default:
		return nil, fmt.Errorf("invalid mermaid rendering mode: %q", config.Mode)
	}
}
