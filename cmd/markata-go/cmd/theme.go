package cmd

import (
	"encoding/json"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/palettes"
	"github.com/WaylonWalker/markata/pkg/renderingcontract"
	"github.com/spf13/cobra"
)

// themeCmd represents the theme command group.
var themeCmd = &cobra.Command{
	Use:   "theme",
	Short: "Theme testing and validation commands",
	Long: `Commands for testing and validating themes and palettes.

Subcommands:
  render-all  - Render all theme/palette combinations
  gallery     - Generate a preview gallery of all themes
  check-all   - Run accessibility checks on all themes`,
}

// themeRenderAllCmd renders all theme/palette combinations.
var themeRenderAllCmd = &cobra.Command{
	Use:   "render-all",
	Short: "Render all theme/palette combinations",
	Long: `Render sample content with each available theme/palette combination.

This command generates test sites for each palette, useful for visual
inspection and automated testing.

Example usage:
  markata-go theme render-all
  markata-go theme render-all --output /tmp/theme-gallery/
  markata-go theme render-all --sample-content ./samples/`,
	RunE: runThemeRenderAllCommand,
}

// themeGalleryCmd generates a theme gallery.
var themeGalleryCmd = &cobra.Command{
	Use:   "gallery",
	Short: "Generate theme preview gallery",
	Long: `Generate an HTML gallery showing all themes side-by-side.

The gallery includes:
  - Color swatches for each palette
  - Accessibility scores (WCAG compliance)
  - Theme metadata and variant information
  - Color blindness simulation warnings

Example usage:
  markata-go theme gallery
  markata-go theme gallery --output gallery.html
  markata-go theme gallery --open`,
	RunE: runThemeGalleryCommand,
}

// themeCheckAllCmd runs accessibility checks on all themes.
var themeCheckAllCmd = &cobra.Command{
	Use:   "check-all",
	Short: "Run accessibility checks on all themes",
	Long: `Run comprehensive accessibility checks on all available themes.

Checks include:
  - WCAG AA contrast ratio compliance (16 required combinations)
  - WCAG AAA compliance (optional, with --strict)
  - Color blindness simulation warnings
  - Missing semantic color warnings

Example usage:
  markata-go theme check-all
  markata-go theme check-all --strict         # Include AAA checks
  markata-go theme check-all --json           # Output as JSON
  markata-go theme check-all --colorblindness # Include color blindness warnings`,
	RunE: runThemeCheckAllCommand,
}

// themeCalendarCmd manages seasonal theme calendar.
var themeCalendarCmd = &cobra.Command{
	Use:   "calendar",
	Short: "Manage seasonal theme calendar",
	Long: `Commands for working with the seasonal theme calendar.

The theme calendar allows you to automatically switch themes based on
date ranges. For example, apply a Christmas theme from Dec 15-26.

Subcommands:
  list    - List all configured calendar rules
  preview - Preview which theme applies on a specific date`,
}

// themeCalendarListCmd lists calendar rules.
var themeCalendarListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all calendar rules",
	Long: `List all theme calendar rules from your configuration.

Shows:
  - Rule name
  - Date range (start - end)
  - Palette to apply
  - Current status (active/inactive)

Example usage:
  markata-go theme calendar list`,
	RunE: runThemeCalendarListCommand,
}

// themeCalendarPreviewCmd previews theme for a date.
var themeCalendarPreviewCmd = &cobra.Command{
	Use:   "preview [MM-DD]",
	Short: "Preview theme for a specific date",
	Long: `Preview which theme/palette would be applied on a specific date.

If no date is provided, uses today's date.

Example usage:
  markata-go theme calendar preview          # Check today
  markata-go theme calendar preview 12-25    # Check Christmas
  markata-go theme calendar preview 01-01    # Check New Year`,
	RunE: runThemeCalendarPreviewCommand,
}

var (
	// themeOutputDir is the output directory for rendered themes.
	themeOutputDir string

	// themeSampleContent is the path to sample content for rendering.
	themeSampleContent string

	// themeGalleryOutput is the output file for the gallery.
	themeGalleryOutput string

	// themeGalleryOpen opens the gallery in browser.
	themeGalleryOpen bool

	// themeCheckStrict includes AAA level checks.
	themeCheckStrict bool

	// themeCheckJSON outputs results as JSON.
	themeCheckJSON bool

	// themeCheckColorblindness includes color blindness warnings.
	themeCheckColorblindness bool
)

func init() {
	rootCmd.AddCommand(themeCmd)

	// render-all subcommand
	themeCmd.AddCommand(themeRenderAllCmd)
	themeRenderAllCmd.Flags().StringVarP(&themeOutputDir, "output", "o", "theme-gallery", "Output directory for rendered themes")
	themeRenderAllCmd.Flags().StringVar(&themeSampleContent, "sample-content", "", "Path to sample content (optional)")

	// gallery subcommand
	themeCmd.AddCommand(themeGalleryCmd)
	themeGalleryCmd.Flags().StringVarP(&themeGalleryOutput, "output", "o", "theme-gallery.html", "Output file for the gallery")
	themeGalleryCmd.Flags().BoolVar(&themeGalleryOpen, "open", false, "Open gallery in browser after generation")

	// check-all subcommand
	themeCmd.AddCommand(themeCheckAllCmd)
	themeCheckAllCmd.Flags().BoolVar(&themeCheckStrict, "strict", false, "Include AAA level checks")
	themeCheckAllCmd.Flags().BoolVar(&themeCheckJSON, "json", false, "Output results as JSON")
	themeCheckAllCmd.Flags().BoolVar(&themeCheckColorblindness, "colorblindness", false, "Include color blindness simulation warnings")

	// calendar subcommand
	themeCmd.AddCommand(themeCalendarCmd)
	themeCalendarCmd.AddCommand(themeCalendarListCmd)
	themeCalendarCmd.AddCommand(themeCalendarPreviewCmd)
}

// ThemeRenderResult holds the result of rendering a single theme.
type ThemeRenderResult struct {
	Palette   string `json:"palette"`
	Variant   string `json:"variant"`
	OutputDir string `json:"output_dir"`
	Success   bool   `json:"success"`
	Error     string `json:"error,omitempty"`
}

// runThemeRenderAllCommand renders all theme/palette combinations.
func runThemeRenderAllCommand(_ *cobra.Command, _ []string) error {
	loader := palettes.NewLoader()
	infos, err := loader.Discover()
	if err != nil {
		return fmt.Errorf("failed to discover palettes: %w", err)
	}

	if len(infos) == 0 {
		fmt.Println("No palettes found.")
		return nil
	}

	// Create output directory
	if err := os.MkdirAll(themeOutputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	fmt.Printf("Rendering %d palettes to %s/\n", len(infos), themeOutputDir)
	fmt.Println(strings.Repeat("-", 50))

	results := make([]ThemeRenderResult, 0, len(infos))
	for _, info := range infos {
		result := renderTheme(loader, info, themeOutputDir)
		results = append(results, result)

		status := "OK"
		if !result.Success {
			status = "FAILED: " + result.Error
		}
		fmt.Printf("  %-30s [%s] %s\n", info.Name, info.Variant, status)
	}

	// Summary
	successCount := 0
	for _, r := range results {
		if r.Success {
			successCount++
		}
	}

	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("Rendered %d/%d palettes successfully\n", successCount, len(results))
	fmt.Printf("Output: %s/\n", themeOutputDir)

	return nil
}

// renderTheme renders a single theme to the output directory.
func renderTheme(loader *palettes.Loader, info palettes.PaletteInfo, outputDir string) ThemeRenderResult {
	result := ThemeRenderResult{
		Palette: info.Name,
		Variant: string(info.Variant),
	}

	p, err := loader.Load(info.Name)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	// Create palette-specific output directory
	paletteDirName := normalizeFileName(info.Name)
	paletteDir := filepath.Join(outputDir, paletteDirName)
	result.OutputDir = paletteDir

	if err := os.MkdirAll(paletteDir, 0o755); err != nil {
		result.Error = fmt.Sprintf("failed to create directory: %v", err)
		return result
	}

	// Generate CSS
	css := p.GenerateCSS()
	cssPath := filepath.Join(paletteDir, "variables.css")
	if err := os.WriteFile(cssPath, []byte(css), 0o644); err != nil { //nolint:gosec // preview files should be readable
		result.Error = fmt.Sprintf("failed to write CSS: %v", err)
		return result
	}

	// Generate sample HTML page
	htmlContent := generateThemeSampleHTML(p)
	htmlPath := filepath.Join(paletteDir, "index.html")
	if err := os.WriteFile(htmlPath, []byte(htmlContent), 0o644); err != nil { //nolint:gosec // preview files should be readable
		result.Error = fmt.Sprintf("failed to write HTML: %v", err)
		return result
	}

	result.Success = true
	return result
}

// generateThemeSampleHTML generates a sample HTML page for a theme.
func generateThemeSampleHTML(p *palettes.Palette) string {
	var sb strings.Builder

	sb.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Theme Preview: `)
	sb.WriteString(html.EscapeString(p.Name))
	sb.WriteString(`</title>
  <link rel="stylesheet" href="variables.css">
  <style>
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      font-family: system-ui, -apple-system, sans-serif;
      background: var(--bg-primary);
      color: var(--text-primary);
      line-height: 1.6;
      padding: 2rem;
    }
    .container { max-width: 800px; margin: 0 auto; }
    header { margin-bottom: 2rem; border-bottom: 1px solid var(--border); padding-bottom: 1rem; }
    h1 { color: var(--text-primary); margin-bottom: 0.5rem; }
    .meta { color: var(--text-muted); font-size: 0.875rem; }
    .variant-badge {
      display: inline-block;
      padding: 0.25rem 0.5rem;
      border-radius: 4px;
      font-size: 0.75rem;
      font-weight: 600;
      text-transform: uppercase;
    }
    .variant-dark { background: #1e1e1e; color: #fff; }
    .variant-light { background: #f0f0f0; color: #000; }
    section { margin-bottom: 2rem; }
    h2 { color: var(--text-primary); margin-bottom: 1rem; font-size: 1.25rem; }
    p { margin-bottom: 1rem; }
    a { color: var(--link); }
    a:hover { color: var(--link-hover); }
    .surface { background: var(--bg-surface); padding: 1rem; border-radius: 8px; margin-bottom: 1rem; }
    .elevated { background: var(--bg-elevated); padding: 1rem; border-radius: 8px; }
    .status-colors { display: flex; gap: 1rem; flex-wrap: wrap; margin-bottom: 1rem; }
    .status-pill {
      padding: 0.5rem 1rem;
      border-radius: 9999px;
      font-size: 0.875rem;
      font-weight: 500;
    }
    .status-success { background: var(--success); color: var(--bg-primary); }
    .status-warning { background: var(--warning); color: var(--bg-primary); }
    .status-error { background: var(--error); color: var(--bg-primary); }
    .status-info { background: var(--info); color: var(--bg-primary); }
    .code-block {
      background: var(--code-bg);
      color: var(--code-text);
      padding: 1rem;
      border-radius: 8px;
      font-family: 'JetBrains Mono', 'Fira Code', monospace;
      font-size: 0.875rem;
      overflow-x: auto;
    }
    .code-comment { color: var(--code-comment); }
    .code-keyword { color: var(--code-keyword); }
    .code-string { color: var(--code-string); }
    .code-function { color: var(--code-function); }
    .button-row { display: flex; gap: 0.5rem; flex-wrap: wrap; }
    button {
      padding: 0.5rem 1rem;
      border-radius: 6px;
      border: none;
      cursor: pointer;
      font-size: 0.875rem;
      font-weight: 500;
    }
    .btn-primary { background: var(--button-primary-bg); color: var(--button-primary-text); }
    .btn-secondary { background: var(--button-secondary-bg); color: var(--button-secondary-text); }
    .accent-text { color: var(--accent); font-weight: 600; }
  </style>
</head>
<body>
  <div class="container">
    <header>
      <h1>`)
	sb.WriteString(html.EscapeString(p.Name))
	sb.WriteString(`</h1>
      <p class="meta">
        <span class="variant-badge variant-`)
	sb.WriteString(string(p.Variant))
	sb.WriteString(`">`)
	sb.WriteString(string(p.Variant))
	sb.WriteString(`</span>`)
	if p.Author != "" {
		sb.WriteString(` &middot; by `)
		sb.WriteString(html.EscapeString(p.Author))
	}
	if p.Description != "" {
		sb.WriteString(` &middot; `)
		sb.WriteString(html.EscapeString(p.Description))
	}
	sb.WriteString(`
      </p>
    </header>

    <section>
      <h2>Typography</h2>
      <p>This is <strong>primary text</strong> on the primary background. Links look like <a href="#">this example link</a>.</p>
      <p style="color: var(--text-secondary);">This is secondary text, used for less important content.</p>
      <p style="color: var(--text-muted);">This is muted text, used for timestamps and metadata.</p>
      <p>Here is some <span class="accent-text">accented text</span> for emphasis.</p>
    </section>

    <section>
      <h2>Surfaces</h2>
      <div class="surface">
        <p><strong>Surface:</strong> Secondary content area (cards, sidebars)</p>
        <div class="elevated">
          <p><strong>Elevated:</strong> Dropdowns, modals, tooltips</p>
        </div>
      </div>
    </section>

    <section>
      <h2>Status Colors</h2>
      <div class="status-colors">
        <span class="status-pill status-success">Success</span>
        <span class="status-pill status-warning">Warning</span>
        <span class="status-pill status-error">Error</span>
        <span class="status-pill status-info">Info</span>
      </div>
    </section>

    <section>
      <h2>Code Block</h2>
      <pre class="code-block"><span class="code-comment">// Example Go code</span>
<span class="code-keyword">func</span> <span class="code-function">main</span>() {
    message := <span class="code-string">"Hello, World!"</span>
    fmt.Println(message)
}</pre>
    </section>

    <section>
      <h2>Buttons</h2>
      <div class="button-row">
        <button class="btn-primary">Primary Button</button>
        <button class="btn-secondary">Secondary Button</button>
      </div>
    </section>
  </div>
</body>
</html>
`)

	return sb.String()
}

// ThemeGalleryEntry holds data for a single theme in the gallery.
type ThemeGalleryEntry struct {
	Name           string                  `json:"name"`
	Variant        string                  `json:"variant"`
	Author         string                  `json:"author,omitempty"`
	Description    string                  `json:"description,omitempty"`
	Source         string                  `json:"source"`
	ContrastScore  ThemeContrastScore      `json:"contrast_score"`
	ColorSwatches  []ThemeColorSwatch      `json:"color_swatches"`
	ColorBlindness []ColorBlindnessWarning `json:"colorblindness_warnings,omitempty"`
}

// ThemeContrastScore holds accessibility scoring for a theme.
type ThemeContrastScore struct {
	Passed      int  `json:"passed"`
	Failed      int  `json:"failed"`
	Skipped     int  `json:"skipped"`
	Total       int  `json:"total"`
	PassPercent int  `json:"pass_percent"`
	AllPassed   bool `json:"all_passed"`
}

// ThemeColorSwatch holds a single color for display.
type ThemeColorSwatch struct {
	Name string `json:"name"`
	Hex  string `json:"hex"`
	Type string `json:"type"` // "raw", "semantic", "component"
}

// ColorBlindnessWarning holds a warning about color blindness issues.
type ColorBlindnessWarning struct {
	Type        string `json:"type"` // "protanopia", "deuteranopia", "tritanopia"
	Description string `json:"description"`
	Colors      string `json:"colors"` // Which colors are affected
}

// runThemeGalleryCommand generates the theme gallery.
func runThemeGalleryCommand(_ *cobra.Command, _ []string) error {
	loader := palettes.NewLoader()
	infos, err := loader.Discover()
	if err != nil {
		return fmt.Errorf("failed to discover palettes: %w", err)
	}

	if len(infos) == 0 {
		fmt.Println("No palettes found.")
		return nil
	}

	// Sort palettes: by variant (dark first), then by name
	sort.Slice(infos, func(i, j int) bool {
		if infos[i].Variant != infos[j].Variant {
			return infos[i].Variant == palettes.VariantDark
		}
		return infos[i].Name < infos[j].Name
	})

	fmt.Printf("Generating gallery for %d palettes...\n", len(infos))

	entries := make([]ThemeGalleryEntry, 0, len(infos))
	for _, info := range infos {
		p, err := loader.Load(info.Name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to load %s: %v\n", info.Name, err)
			continue
		}

		entry := buildGalleryEntry(p, themeCheckColorblindness)
		entries = append(entries, entry)
	}

	// Generate HTML gallery
	galleryHTML := generateGalleryHTML(entries)

	if err := os.WriteFile(themeGalleryOutput, []byte(galleryHTML), 0o644); err != nil { //nolint:gosec // gallery files should be readable
		return fmt.Errorf("failed to write gallery: %w", err)
	}

	fmt.Printf("Gallery generated: %s\n", themeGalleryOutput)

	if themeGalleryOpen {
		openBrowser(themeGalleryOutput)
	}

	return nil
}

// buildGalleryEntry builds a gallery entry for a palette.
func buildGalleryEntry(p *palettes.Palette, includeColorblindness bool) ThemeGalleryEntry {
	entry := ThemeGalleryEntry{
		Name:        p.Name,
		Variant:     string(p.Variant),
		Author:      p.Author,
		Description: p.Description,
		Source:      p.Source,
	}

	// Run contrast checks
	results := p.CheckContrast()
	summary := palettes.SummarizeContrast(p.Name, results)
	entry.ContrastScore = ThemeContrastScore{
		Passed:    summary.Passed,
		Failed:    summary.Failed,
		Skipped:   summary.Skipped,
		Total:     summary.Total,
		AllPassed: summary.AllPassed,
	}
	if summary.Total > 0 {
		entry.ContrastScore.PassPercent = (summary.Passed * 100) / summary.Total
	}

	// Build color swatches
	// Key semantic colors first
	semanticColors := []string{"text-primary", "text-secondary", "bg-primary", "bg-surface", "link", "accent", "success", "warning", "error", "info"}
	for _, name := range semanticColors {
		hex := p.Resolve(name)
		if hex != "" {
			entry.ColorSwatches = append(entry.ColorSwatches, ThemeColorSwatch{
				Name: name,
				Hex:  hex,
				Type: "semantic",
			})
		}
	}

	// Color blindness warnings
	if includeColorblindness {
		entry.ColorBlindness = analyzeColorBlindnessRisks(p)
	}

	return entry
}

// analyzeColorBlindnessRisks analyzes a palette for color blindness issues.
func analyzeColorBlindnessRisks(p *palettes.Palette) []ColorBlindnessWarning {
	var warnings []ColorBlindnessWarning

	// Get key status colors
	success := p.Resolve("success")
	warning := p.Resolve("warning")
	errColor := p.Resolve("error")

	// Check for red-green confusion (protanopia/deuteranopia)
	if success != "" && errColor != "" {
		successC, err1 := palettes.ParseHexColor(success)
		errorC, err2 := palettes.ParseHexColor(errColor)
		if err1 == nil && err2 == nil {
			// Check if both colors are primarily in the red-green spectrum
			// and could be confused by someone with red-green color blindness
			if isRedGreenConfusable(successC, errorC) {
				warnings = append(warnings, ColorBlindnessWarning{
					Type:        "protanopia/deuteranopia",
					Description: "Success and error colors may be difficult to distinguish for users with red-green color blindness",
					Colors:      fmt.Sprintf("success (%s) vs error (%s)", success, errColor),
				})
			}
		}
	}

	// Check for yellow-blue confusion (tritanopia)
	if warning != "" {
		warnC, err := palettes.ParseHexColor(warning)
		if err == nil && isLowBlueYellowContrast(warnC) {
			warnings = append(warnings, ColorBlindnessWarning{
				Type:        "tritanopia",
				Description: "Warning color may be difficult to perceive for users with blue-yellow color blindness",
				Colors:      fmt.Sprintf("warning (%s)", warning),
			})
		}
	}

	return warnings
}

// isRedGreenConfusable checks if two colors might be confused by red-green color blind users.
func isRedGreenConfusable(c1, c2 palettes.Color) bool {
	// Simplified check: if both colors have similar blue values but differ mainly in red/green,
	// they might be confusable. This is a heuristic, not a full simulation.
	blueDiff := abs(int(c1.B) - int(c2.B))
	redGreenDiff := abs(int(c1.R)-int(c2.R)) + abs(int(c1.G)-int(c2.G))

	// If colors differ mainly in red/green channel and blue is similar
	return blueDiff < 50 && redGreenDiff < 100
}

// isLowBlueYellowContrast checks if a color might have issues for tritanopia.
func isLowBlueYellowContrast(c palettes.Color) bool {
	// Yellow colors (high R, high G, low B) can be problematic for tritanopia
	// This is a simplified heuristic
	return c.R > 200 && c.G > 150 && c.B < 100
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// generateGalleryHTML generates the HTML gallery page.
func generateGalleryHTML(entries []ThemeGalleryEntry) string {
	var sb strings.Builder

	sb.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Theme Gallery - markata-go</title>
  <style>
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      font-family: system-ui, -apple-system, sans-serif;
      background: #0f0f0f;
      color: #e0e0e0;
      line-height: 1.6;
      padding: 2rem;
    }
    .container { max-width: 1400px; margin: 0 auto; }
    header { margin-bottom: 2rem; text-align: center; }
    h1 { font-size: 2rem; margin-bottom: 0.5rem; }
    .subtitle { color: #888; }
    .stats { display: flex; justify-content: center; gap: 2rem; margin-top: 1rem; }
    .stat { text-align: center; }
    .stat-value { font-size: 2rem; font-weight: 700; color: #7c3aed; }
    .stat-label { font-size: 0.75rem; color: #888; text-transform: uppercase; }
    .filters { display: flex; justify-content: center; gap: 1rem; margin-bottom: 2rem; }
    .filter-btn {
      padding: 0.5rem 1rem;
      border: 1px solid #333;
      background: transparent;
      color: #e0e0e0;
      border-radius: 6px;
      cursor: pointer;
      font-size: 0.875rem;
    }
    .filter-btn:hover, .filter-btn.active { background: #7c3aed; border-color: #7c3aed; }
    .gallery { display: grid; grid-template-columns: repeat(auto-fill, minmax(320px, 1fr)); gap: 1.5rem; }
    .theme-card {
      background: #1a1a1a;
      border-radius: 12px;
      overflow: hidden;
      border: 1px solid #333;
      transition: transform 0.2s, box-shadow 0.2s;
    }
    .theme-card:hover { transform: translateY(-2px); box-shadow: 0 8px 24px rgba(0,0,0,0.3); }
    .theme-preview {
      height: 120px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 1.25rem;
      font-weight: 600;
      position: relative;
    }
    .theme-info { padding: 1rem; }
    .theme-name { font-size: 1.125rem; font-weight: 600; margin-bottom: 0.25rem; }
    .theme-meta { font-size: 0.75rem; color: #888; margin-bottom: 0.75rem; }
    .variant-badge {
      display: inline-block;
      padding: 0.125rem 0.5rem;
      border-radius: 4px;
      font-size: 0.625rem;
      font-weight: 600;
      text-transform: uppercase;
    }
    .variant-dark { background: #333; color: #fff; }
    .variant-light { background: #ddd; color: #000; }
    .score-bar { height: 4px; background: #333; border-radius: 2px; margin-bottom: 0.5rem; overflow: hidden; }
    .score-fill { height: 100%; border-radius: 2px; transition: width 0.3s; }
    .score-100 { background: #10b981; }
    .score-high { background: #84cc16; }
    .score-medium { background: #f59e0b; }
    .score-low { background: #ef4444; }
    .score-text { font-size: 0.75rem; display: flex; justify-content: space-between; }
    .score-label { color: #888; }
    .score-value { font-weight: 600; }
    .pass-all { color: #10b981; }
    .pass-some { color: #f59e0b; }
    .pass-fail { color: #ef4444; }
    .swatches { display: flex; flex-wrap: wrap; gap: 4px; margin-top: 0.75rem; }
    .swatch { width: 24px; height: 24px; border-radius: 4px; border: 1px solid rgba(255,255,255,0.1); cursor: pointer; }
    .swatch:hover { transform: scale(1.2); }
    .warnings { margin-top: 0.5rem; }
    .warning { font-size: 0.75rem; color: #f59e0b; display: flex; align-items: center; gap: 0.25rem; }
    .warning-icon { font-size: 0.875rem; }
    footer { margin-top: 3rem; text-align: center; color: #666; font-size: 0.875rem; }
    .timestamp { font-family: monospace; }
  </style>
</head>
<body>
  <div class="container">
    <header>
      <h1>Theme Gallery</h1>
      <p class="subtitle">All available themes and palettes for markata-go</p>
      <div class="stats">
        <div class="stat">
          <div class="stat-value">`)
	sb.WriteString(fmt.Sprintf("%d", len(entries)))
	sb.WriteString(`</div>
          <div class="stat-label">Total Themes</div>
        </div>
        <div class="stat">
          <div class="stat-value">`)
	darkCount := 0
	lightCount := 0
	passCount := 0
	for i := range entries {
		if entries[i].Variant == "dark" {
			darkCount++
		} else {
			lightCount++
		}
		if entries[i].ContrastScore.AllPassed {
			passCount++
		}
	}
	sb.WriteString(fmt.Sprintf("%d", darkCount))
	sb.WriteString(`</div>
          <div class="stat-label">Dark Themes</div>
        </div>
        <div class="stat">
          <div class="stat-value">`)
	sb.WriteString(fmt.Sprintf("%d", lightCount))
	sb.WriteString(`</div>
          <div class="stat-label">Light Themes</div>
        </div>
        <div class="stat">
          <div class="stat-value">`)
	sb.WriteString(fmt.Sprintf("%d", passCount))
	sb.WriteString(`</div>
          <div class="stat-label">WCAG AA Pass</div>
        </div>
      </div>
    </header>

    <div class="filters">
      <button class="filter-btn active" onclick="filterThemes('all')">All</button>
      <button class="filter-btn" onclick="filterThemes('dark')">Dark</button>
      <button class="filter-btn" onclick="filterThemes('light')">Light</button>
      <button class="filter-btn" onclick="filterThemes('passing')">WCAG Passing</button>
    </div>

    <div class="gallery">
`)

	for i := range entries {
		sb.WriteString(generateThemeCard(entries[i]))
	}

	sb.WriteString(`    </div>

    <footer>
      <p>Generated by <strong>markata-go theme gallery</strong></p>
      <p class="timestamp">`)
	sb.WriteString(time.Now().Format(time.RFC3339))
	sb.WriteString(`</p>
    </footer>
  </div>

  <script>
    function filterThemes(filter) {
      document.querySelectorAll('.filter-btn').forEach(btn => btn.classList.remove('active'));
      event.target.classList.add('active');

      document.querySelectorAll('.theme-card').forEach(card => {
        const variant = card.dataset.variant;
        const passing = card.dataset.passing === 'true';

        let show = true;
        if (filter === 'dark') show = variant === 'dark';
        else if (filter === 'light') show = variant === 'light';
        else if (filter === 'passing') show = passing;

        card.style.display = show ? 'block' : 'none';
      });
    }
  </script>
</body>
</html>
`)

	return sb.String()
}

// generateThemeCard generates HTML for a single theme card.
func generateThemeCard(entry ThemeGalleryEntry) string {
	var sb strings.Builder

	// Get background and text colors for preview
	bgColor := defaultBgColor
	textColor := defaultTextColor
	for _, swatch := range entry.ColorSwatches {
		if swatch.Name == "bg-primary" {
			bgColor = swatch.Hex
		}
		if swatch.Name == "text-primary" {
			textColor = swatch.Hex
		}
	}

	scoreClass := "score-100"
	textClass := "pass-all"
	if !entry.ContrastScore.AllPassed {
		switch {
		case entry.ContrastScore.PassPercent >= 80:
			scoreClass = "score-high"
			textClass = "pass-some"
		case entry.ContrastScore.PassPercent >= 50:
			scoreClass = "score-medium"
			textClass = "pass-some"
		default:
			scoreClass = "score-low"
			textClass = "pass-fail"
		}
	}

	sb.WriteString(fmt.Sprintf(`      <div class="theme-card" data-variant="%s" data-passing="%t">
        <div class="theme-preview" style="background: %s; color: %s;">
          %s
        </div>
        <div class="theme-info">
          <div class="theme-name">%s</div>
          <div class="theme-meta">
            <span class="variant-badge variant-%s">%s</span>`,
		entry.Variant,
		entry.ContrastScore.AllPassed,
		bgColor,
		textColor,
		html.EscapeString(entry.Name),
		html.EscapeString(entry.Name),
		entry.Variant,
		entry.Variant,
	))

	if entry.Author != "" {
		sb.WriteString(fmt.Sprintf(` &middot; %s`, html.EscapeString(entry.Author)))
	}
	sb.WriteString(fmt.Sprintf(` &middot; %s`, entry.Source))
	sb.WriteString(`
          </div>
          <div class="score-bar">
            <div class="score-fill `)
	sb.WriteString(scoreClass)
	sb.WriteString(`" style="width: `)
	sb.WriteString(fmt.Sprintf("%d", entry.ContrastScore.PassPercent))
	sb.WriteString(`%;"></div>
          </div>
          <div class="score-text">
            <span class="score-label">WCAG AA Contrast</span>
            <span class="score-value `)
	sb.WriteString(textClass)
	sb.WriteString(`">`)
	sb.WriteString(fmt.Sprintf("%d/%d", entry.ContrastScore.Passed, entry.ContrastScore.Total))
	sb.WriteString(`</span>
          </div>
          <div class="swatches">
`)

	for _, swatch := range entry.ColorSwatches {
		sb.WriteString(fmt.Sprintf(`            <div class="swatch" style="background: %s;" title="%s: %s"></div>
`, swatch.Hex, swatch.Name, swatch.Hex))
	}

	sb.WriteString(`          </div>`)

	if len(entry.ColorBlindness) > 0 {
		sb.WriteString(`
          <div class="warnings">
`)
		for _, warn := range entry.ColorBlindness {
			sb.WriteString(fmt.Sprintf(`            <div class="warning"><span class="warning-icon">⚠</span> %s</div>
`, html.EscapeString(warn.Type)))
		}
		sb.WriteString(`          </div>`)
	}

	sb.WriteString(`
        </div>
      </div>
`)

	return sb.String()
}

// ThemeCheckResult holds the result of checking a single theme.
type ThemeCheckResult struct {
	Palette         string                   `json:"palette"`
	Variant         string                   `json:"variant"`
	ContrastSummary palettes.ContrastSummary `json:"contrast_summary"`
	ColorBlindness  []ColorBlindnessWarning  `json:"colorblindness_warnings,omitempty"`
	AllPassed       bool                     `json:"all_passed"`
}

// ThemeCheckAllResult holds the result of checking all themes.
type ThemeCheckAllResult struct {
	Timestamp     string             `json:"timestamp"`
	TotalPalettes int                `json:"total_palettes"`
	AllPassing    int                `json:"all_passing"`
	SomeFailing   int                `json:"some_failing"`
	Results       []ThemeCheckResult `json:"results"`
}

// runThemeCheckAllCommand runs accessibility checks on all themes.
func runThemeCheckAllCommand(_ *cobra.Command, _ []string) error {
	contract, err := renderingcontract.Load()
	if err != nil {
		return fmt.Errorf("failed to load rendering contract: %w", err)
	}
	if len(contract.Palettes) == 0 {
		fmt.Println("No palettes found.")
		return nil
	}

	allResult := ThemeCheckAllResult{
		Timestamp:     time.Now().Format(time.RFC3339),
		TotalPalettes: len(contract.Palettes),
	}

	if !themeCheckJSON {
		fmt.Printf("Checking final-render accessibility for %d canonical palettes...\n", len(contract.Palettes))
		fmt.Println(strings.Repeat("=", 60))
	}

	for _, p := range contract.Palettes {
		results := renderingcontract.FinalRenderContrast(p)
		summary := palettes.SummarizeContrast(p.ID, results)

		checkResult := ThemeCheckResult{
			Palette:         p.ID,
			Variant:         p.Variant,
			ContrastSummary: summary,
			AllPassed:       summary.AllPassed,
		}

		allResult.Results = append(allResult.Results, checkResult)

		if summary.AllPassed {
			allResult.AllPassing++
		} else {
			allResult.SomeFailing++
		}

		if !themeCheckJSON {
			printThemeCheckResult(checkResult, themeCheckStrict)
		}
	}

	if themeCheckJSON {
		data, err := json.MarshalIndent(allResult, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
	} else {
		fmt.Println(strings.Repeat("=", 60))
		fmt.Printf("\nSummary: %d/%d palettes pass all WCAG %s checks\n",
			allResult.AllPassing,
			allResult.TotalPalettes,
			getWCAGLevel(themeCheckStrict))

		if allResult.SomeFailing > 0 {
			fmt.Printf("\n%d palettes have failing checks:\n", allResult.SomeFailing)
			for i := range allResult.Results {
				if !allResult.Results[i].AllPassed {
					fmt.Printf("  - %s (%d/%d passed)\n",
						allResult.Results[i].Palette,
						allResult.Results[i].ContrastSummary.Passed,
						allResult.Results[i].ContrastSummary.Total)
				}
			}
		}
	}

	// Return error if any checks failed (useful for CI)
	if allResult.SomeFailing > 0 {
		return fmt.Errorf("%d palettes have failing accessibility checks", allResult.SomeFailing)
	}

	return nil
}

// printThemeCheckResult prints a single theme check result.
func printThemeCheckResult(result ThemeCheckResult, strict bool) {
	status := "\u2713" // checkmark
	statusColor := ""
	if !result.AllPassed {
		status = "\u2717" // X
		statusColor = " (FAIL)"
	}

	fmt.Printf("\n%s %s [%s]%s\n",
		status,
		result.Palette,
		result.Variant,
		statusColor)

	summary := result.ContrastSummary
	fmt.Printf("   Contrast: %d passed, %d failed",
		summary.Passed, summary.Failed)
	if summary.Skipped > 0 {
		fmt.Printf(", %d skipped", summary.Skipped)
	}
	fmt.Printf(" (WCAG %s)\n", getWCAGLevel(strict))

	// Show failed checks
	if len(summary.FailedChecks) > 0 && len(summary.FailedChecks) <= 5 {
		for i := range summary.FailedChecks {
			fmt.Printf("   \u2717 %s on %s: %.1f:1 (need %.1f:1)\n",
				summary.FailedChecks[i].Foreground, summary.FailedChecks[i].Background,
				summary.FailedChecks[i].Ratio, summary.FailedChecks[i].Required)
		}
	} else if len(summary.FailedChecks) > 5 {
		for i := 0; i < 3; i++ {
			fmt.Printf("   \u2717 %s on %s: %.1f:1 (need %.1f:1)\n",
				summary.FailedChecks[i].Foreground, summary.FailedChecks[i].Background,
				summary.FailedChecks[i].Ratio, summary.FailedChecks[i].Required)
		}
		fmt.Printf("   ... and %d more failures\n", len(summary.FailedChecks)-3)
	}

	// Color blindness warnings
	if len(result.ColorBlindness) > 0 {
		fmt.Printf("   Color blindness warnings:\n")
		for _, warn := range result.ColorBlindness {
			fmt.Printf("   \u26A0 %s: %s\n", warn.Type, warn.Description)
		}
	}
}

// getWCAGLevel returns the WCAG level string.
func getWCAGLevel(strict bool) string {
	if strict {
		return "AA+AAA"
	}
	return "AA"
}

// runThemeCalendarListCommand lists all calendar rules.
func runThemeCalendarListCommand(_ *cobra.Command, _ []string) error {
	// Load config
	cfg, err := loadConfigForCalendar()
	if err != nil {
		return err
	}

	rules := getCalendarRules(cfg)
	if len(rules) == 0 {
		fmt.Println("No theme calendar rules configured.")
		fmt.Println("\nTo add rules, add to your markata-go.toml:")
		fmt.Print(`
[markata-go.theme_calendar]
enabled = true

[[markata-go.theme_calendar.rules]]
name = "Christmas Season"
start_date = "12-15"
end_date = "12-26"
palette = "christmas"
`)
		return nil
	}

	// Check if calendar is enabled
	enabled := isCalendarEnabled(cfg)
	if !enabled {
		fmt.Println("Theme calendar is DISABLED")
		fmt.Println("Set 'enabled = true' in [markata-go.theme_calendar] to enable")
		fmt.Println()
	}

	// Get current date for status
	now := time.Now()
	currentMonth := int(now.Month())
	currentDay := now.Day()

	fmt.Printf("Theme Calendar Rules (%d configured)\n", len(rules))
	fmt.Println(strings.Repeat("=", 60))

	for _, rule := range rules {
		printCalendarRule(rule, currentMonth, currentDay, enabled)
	}

	fmt.Println()
	return nil
}

// printCalendarRule prints a single calendar rule.
func printCalendarRule(rule map[string]interface{}, currentMonth, currentDay int, enabled bool) {
	name := getStringFromRule(rule, "name")
	startDate := getStringFromRule(rule, "start_date")
	endDate := getStringFromRule(rule, "end_date")

	// Check if active
	active := isDateInRangeForCLI(currentMonth, currentDay, startDate, endDate)
	status := ""
	if active && enabled {
		status = " [ACTIVE]"
	}

	fmt.Printf("\n%s%s\n", name, status)
	fmt.Printf("  Date Range: %s to %s\n", startDate, endDate)

	if palette := getStringFromRule(rule, "palette"); palette != "" {
		fmt.Printf("  Palette: %s\n", palette)
	}
	if paletteLight := getStringFromRule(rule, "palette_light"); paletteLight != "" {
		fmt.Printf("  Light Palette: %s\n", paletteLight)
	}
	if paletteDark := getStringFromRule(rule, "palette_dark"); paletteDark != "" {
		fmt.Printf("  Dark Palette: %s\n", paletteDark)
	}

	// Show other overrides if present
	if customCSS, ok := rule["custom_css"].(string); ok && customCSS != "" {
		fmt.Printf("  Custom CSS: (defined)\n")
	}
	if vars, ok := rule["variables"].(map[string]interface{}); ok && len(vars) > 0 {
		fmt.Printf("  CSS Variables: %d defined\n", len(vars))
	}
	if _, ok := rule["background"].(map[string]interface{}); ok {
		fmt.Printf("  Background: (custom)\n")
	}
	if _, ok := rule["font"].(map[string]interface{}); ok {
		fmt.Printf("  Font: (custom)\n")
	}
}

// getStringFromRule safely extracts a string value from a rule map.
func getStringFromRule(rule map[string]interface{}, key string) string {
	if v, ok := rule[key].(string); ok {
		return v
	}
	return ""
}

// runThemeCalendarPreviewCommand previews the theme for a specific date.
func runThemeCalendarPreviewCommand(_ *cobra.Command, args []string) error {
	// Load config
	cfg, err := loadConfigForCalendar()
	if err != nil {
		return err
	}

	rules := getCalendarRules(cfg)
	if len(rules) == 0 {
		fmt.Println("No theme calendar rules configured.")
		return nil
	}

	// Parse target date
	targetMonth, targetDay, err := parseTargetDate(args)
	if err != nil {
		return err
	}

	fmt.Printf("Checking theme for date: %02d-%02d\n", targetMonth, targetDay)
	fmt.Println(strings.Repeat("-", 40))

	// Check if calendar is enabled
	if !isCalendarEnabled(cfg) {
		fmt.Println("\nNote: Theme calendar is DISABLED in config")
	}

	// Find matching rule
	matchingRule := findMatchingRule(rules, targetMonth, targetDay)
	if matchingRule == nil {
		fmt.Println("\nNo matching rule found for this date.")
		fmt.Println("The base theme configuration will be used.")
		return nil
	}

	printMatchingRuleDetails(matchingRule)
	return nil
}

// parseTargetDate parses the target date from args or returns current date.
func parseTargetDate(args []string) (month, day int, err error) {
	if len(args) > 0 {
		parsed, parseErr := parseMMDDForCLI(args[0])
		if parseErr != nil {
			return 0, 0, fmt.Errorf("invalid date format %q: %w (expected MM-DD)", args[0], parseErr)
		}
		return parsed.month, parsed.day, nil
	}
	now := time.Now()
	return int(now.Month()), now.Day(), nil
}

// findMatchingRule finds the first rule that matches the given date.
func findMatchingRule(rules []map[string]interface{}, month, day int) map[string]interface{} {
	for _, rule := range rules {
		startDate := getStringFromRule(rule, "start_date")
		endDate := getStringFromRule(rule, "end_date")
		if isDateInRangeForCLI(month, day, startDate, endDate) {
			return rule
		}
	}
	return nil
}

// printMatchingRuleDetails prints the details of a matching rule.
func printMatchingRuleDetails(rule map[string]interface{}) {
	name := getStringFromRule(rule, "name")
	fmt.Printf("\nMatching Rule: %s\n", name)

	startDate := getStringFromRule(rule, "start_date")
	endDate := getStringFromRule(rule, "end_date")
	fmt.Printf("Date Range: %s to %s\n", startDate, endDate)

	// Show what will be applied
	fmt.Println("\nTheme Overrides:")
	printRulePalettes(rule)
	printRuleVariables(rule)
	printRuleBackgroundAndFont(rule)
}

// printRulePalettes prints palette information from a rule.
func printRulePalettes(rule map[string]interface{}) {
	if palette := getStringFromRule(rule, "palette"); palette != "" {
		fmt.Printf("  Palette: %s\n", palette)
	}
	if paletteLight := getStringFromRule(rule, "palette_light"); paletteLight != "" {
		fmt.Printf("  Light Palette: %s\n", paletteLight)
	}
	if paletteDark := getStringFromRule(rule, "palette_dark"); paletteDark != "" {
		fmt.Printf("  Dark Palette: %s\n", paletteDark)
	}
	if customCSS := getStringFromRule(rule, "custom_css"); customCSS != "" {
		fmt.Printf("  Custom CSS: %s\n", customCSS)
	}
}

// printRuleVariables prints CSS variables from a rule.
func printRuleVariables(rule map[string]interface{}) {
	if vars, ok := rule["variables"].(map[string]interface{}); ok && len(vars) > 0 {
		fmt.Printf("  CSS Variables:\n")
		for k, v := range vars {
			fmt.Printf("    %s: %v\n", k, v)
		}
	}
}

// printRuleBackgroundAndFont prints background and font info from a rule.
func printRuleBackgroundAndFont(rule map[string]interface{}) {
	if bg, ok := rule["background"].(map[string]interface{}); ok {
		if enabled, ok := bg["enabled"].(bool); ok && enabled {
			fmt.Printf("  Background: enabled\n")
		}
	}
	if font, ok := rule["font"].(map[string]interface{}); ok {
		if family, ok := font["family"].(string); ok && family != "" {
			fmt.Printf("  Font Family: %s\n", family)
		}
	}
}

// loadConfigForCalendar loads the config file for calendar commands.
func loadConfigForCalendar() (map[string]interface{}, error) {
	configPaths := []string{
		"markata-go.toml",
		"markata.toml",
	}

	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			data, err := os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("failed to read config: %w", err)
			}

			// Simple TOML parsing for the specific section we need
			// For full TOML support, we'd use a proper parser
			cfg := parseSimpleTOML(string(data))
			return cfg, nil
		}
	}

	return nil, fmt.Errorf("no config file found (checked: %s)", strings.Join(configPaths, ", "))
}

// parseSimpleTOML is a simple TOML parser for theme calendar config.
// This extracts the theme_calendar section without full TOML parsing.
func parseSimpleTOML(content string) map[string]interface{} {
	cfg := make(map[string]interface{})

	// Look for markata-go section
	markata := make(map[string]interface{})

	// Extract theme_calendar section
	calendarConfig := extractThemeCalendarSection(content)
	if calendarConfig != nil {
		markata["theme_calendar"] = calendarConfig
	}

	cfg["markata-go"] = markata
	return cfg
}

// tomlParserState tracks the current parsing state.
type tomlParserState struct {
	result            map[string]interface{}
	rules             []map[string]interface{}
	currentRule       map[string]interface{}
	inCalendarSection bool
	inRuleSection     bool
}

// newTOMLParserState creates a new parser state.
func newTOMLParserState() *tomlParserState {
	return &tomlParserState{
		result:      make(map[string]interface{}),
		currentRule: make(map[string]interface{}),
	}
}

// extractThemeCalendarSection extracts the theme_calendar config from TOML content.
func extractThemeCalendarSection(content string) map[string]interface{} {
	state := newTOMLParserState()
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check for section headers
		if strings.HasPrefix(line, "[") {
			state.handleSectionHeader(line)
			continue
		}

		// Parse key-value pairs
		state.handleKeyValue(line)
	}

	return state.finalize()
}

// handleSectionHeader processes a TOML section header.
func (s *tomlParserState) handleSectionHeader(line string) {
	switch {
	case strings.Contains(line, "theme_calendar.rules"):
		s.inRuleSection = true
		s.inCalendarSection = false
		if len(s.currentRule) > 0 {
			s.rules = append(s.rules, s.currentRule)
			s.currentRule = make(map[string]interface{})
		}
	case strings.Contains(line, "theme_calendar") && !strings.Contains(line, "."):
		s.inCalendarSection = true
		s.inRuleSection = false
	default:
		// Different section
		if s.inRuleSection && len(s.currentRule) > 0 {
			s.rules = append(s.rules, s.currentRule)
			s.currentRule = make(map[string]interface{})
		}
		s.inCalendarSection = false
		s.inRuleSection = false
	}
}

// handleKeyValue processes a key-value pair line.
func (s *tomlParserState) handleKeyValue(line string) {
	if !s.inCalendarSection && !s.inRuleSection {
		return
	}

	idx := strings.Index(line, "=")
	if idx <= 0 {
		return
	}

	key := strings.TrimSpace(line[:idx])
	value := strings.TrimSpace(line[idx+1:])
	value = strings.Trim(value, "\"'")

	if s.inCalendarSection {
		s.setCalendarValue(key, value)
	} else if s.inRuleSection {
		s.setRuleValue(key, value)
	}
}

// setCalendarValue sets a value in the calendar section.
func (s *tomlParserState) setCalendarValue(key, value string) {
	switch key {
	case "enabled":
		s.result["enabled"] = value == boolStrTrue
	case "default_palette":
		s.result["default_palette"] = value
	}
}

// setRuleValue sets a value in a rule section.
func (s *tomlParserState) setRuleValue(key, value string) {
	if key == "enabled" {
		s.currentRule[key] = value == boolStrTrue
	} else {
		s.currentRule[key] = value
	}
}

// finalize completes parsing and returns the result.
func (s *tomlParserState) finalize() map[string]interface{} {
	// Don't forget the last rule
	if len(s.currentRule) > 0 {
		s.rules = append(s.rules, s.currentRule)
	}

	if len(s.rules) > 0 {
		s.result["rules"] = s.rules
	}

	if len(s.result) == 0 {
		return nil
	}
	return s.result
}

// getCalendarRules extracts calendar rules from config.
func getCalendarRules(cfg map[string]interface{}) []map[string]interface{} {
	markata, ok := cfg["markata-go"].(map[string]interface{})
	if !ok {
		return nil
	}

	calendar, ok := markata["theme_calendar"].(map[string]interface{})
	if !ok {
		return nil
	}

	rules, ok := calendar["rules"].([]map[string]interface{})
	if !ok {
		// Try interface slice
		if iRules, ok := calendar["rules"].([]interface{}); ok {
			result := make([]map[string]interface{}, 0, len(iRules))
			for _, r := range iRules {
				if rMap, ok := r.(map[string]interface{}); ok {
					result = append(result, rMap)
				}
			}
			return result
		}
		return nil
	}

	return rules
}

// isCalendarEnabled checks if the calendar is enabled.
func isCalendarEnabled(cfg map[string]interface{}) bool {
	markata, ok := cfg["markata-go"].(map[string]interface{})
	if !ok {
		return false
	}

	calendar, ok := markata["theme_calendar"].(map[string]interface{})
	if !ok {
		return false
	}

	enabled, ok := calendar["enabled"].(bool)
	return ok && enabled
}

// parseMMDDForCLI parses a MM-DD date string for CLI commands.
type parsedDate struct {
	month int
	day   int
}

func parseMMDDForCLI(s string) (parsedDate, error) {
	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		return parsedDate{}, fmt.Errorf("expected MM-DD format")
	}

	month, err := strconv.Atoi(parts[0])
	if err != nil || month < 1 || month > 12 {
		return parsedDate{}, fmt.Errorf("invalid month")
	}

	day, err := strconv.Atoi(parts[1])
	if err != nil || day < 1 || day > 31 {
		return parsedDate{}, fmt.Errorf("invalid day")
	}

	return parsedDate{month: month, day: day}, nil
}

// isDateInRangeForCLI checks if a date is in range (for CLI commands).
func isDateInRangeForCLI(month, day int, startDate, endDate string) bool {
	start, err1 := parseMMDDForCLI(startDate)
	end, err2 := parseMMDDForCLI(endDate)
	if err1 != nil || err2 != nil {
		return false
	}

	current := month*100 + day
	startVal := start.month*100 + start.day
	endVal := end.month*100 + end.day

	if startVal <= endVal {
		return current >= startVal && current <= endVal
	}
	// Year boundary crossing
	return current >= startVal || current <= endVal
}
