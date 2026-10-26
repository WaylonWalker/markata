// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma/v2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/styles"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/palettes"
	"github.com/WaylonWalker/markata/pkg/templates"
)

// ChromaCSSPlugin generates CSS for syntax highlighting from Chroma themes.
// It runs during the Write stage and creates css/chroma.css with the
// syntax highlighting styles that correspond to the configured theme.
//
// This plugin works in conjunction with RenderMarkdownPlugin which uses
// CSS classes for syntax highlighting (via WithClasses option).
type ChromaCSSPlugin struct {
	chromaTheme string
	chromaCSS   string // Pre-generated CSS content (computed in Configure)
	chromaHash  string // Content hash for cache busting
	explicit    bool
}

// NewChromaCSSPlugin creates a new ChromaCSSPlugin.
func NewChromaCSSPlugin() *ChromaCSSPlugin {
	return &ChromaCSSPlugin{
		chromaTheme: palettes.DefaultChromaThemeDark,
	}
}

// Name returns the unique name of the plugin.
func (p *ChromaCSSPlugin) Name() string {
	return "chroma_css"
}

// Configure reads the highlight theme configuration and pre-generates the CSS.
// This runs before templates are rendered, allowing us to register the hash
// for cache busting (css/chroma.css -> css/chroma.abc12345.css).
func (p *ChromaCSSPlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()
	extra := config.Extra

	configuredTheme, explicitTheme := p.getExplicitHighlightTheme(extra)
	p.explicit = explicitTheme
	if configuredTheme != "" {
		p.chromaTheme = configuredTheme
	}

	var (
		css string
		err error
	)
	if p.explicit {
		style := styles.Get(p.chromaTheme)
		if style == nil {
			style = styles.Fallback
		}
		css, err = p.generateThemeCSS(style)
	} else {
		css = p.generatePaletteCSS()
	}
	if err != nil {
		return fmt.Errorf("generating chroma CSS: %w", err)
	}
	p.chromaCSS = css

	// Compute hash for cache busting (first 8 chars of SHA-256)
	hash := sha256.Sum256([]byte(css))
	p.chromaHash = fmt.Sprintf("%x", hash[:4]) // 4 bytes = 8 hex chars

	// Register the hash with the template engine for theme_asset_hashed filter
	assetHashes := map[string]string{
		"css/chroma.css": p.chromaHash,
	}
	templates.SetAssetHashes(assetHashes)

	// Also register with the lifecycle manager for tracking
	m.SetAssetHash("css/chroma.css", p.chromaHash)

	return nil
}

// Write generates the Chroma CSS file.
// The CSS content was already generated in Configure() for cache busting.
func (p *ChromaCSSPlugin) Write(m *lifecycle.Manager) error {
	config := m.Config()
	outputDir := config.OutputDir

	// CSS was pre-generated in Configure for hash computation
	css := p.chromaCSS
	if css == "" {
		// Fallback: generate now if Configure didn't run (shouldn't happen)
		if p.explicit {
			style := styles.Get(p.chromaTheme)
			if style == nil {
				style = styles.Fallback
			}
			var err error
			css, err = p.generateThemeCSS(style)
			if err != nil {
				return fmt.Errorf("generating chroma CSS: %w", err)
			}
		} else {
			css = p.generatePaletteCSS()
		}
	}

	// Write to output directory
	cssDir := filepath.Join(outputDir, "css")
	if err := os.MkdirAll(cssDir, 0o755); err != nil {
		return fmt.Errorf("creating css directory: %w", err)
	}

	// Write original filename (css/chroma.css)
	cssPath := filepath.Join(cssDir, "chroma.css")
	//nolint:gosec // G306: chroma.css is a public CSS file, 0644 is appropriate
	if err := os.WriteFile(cssPath, []byte(css), 0o644); err != nil {
		return fmt.Errorf("writing chroma CSS: %w", err)
	}

	// Write hashed version (css/chroma.abc12345.css) for cache busting
	if p.chromaHash != "" {
		hashedFilename := fmt.Sprintf("chroma.%s.css", p.chromaHash)
		hashedPath := filepath.Join(cssDir, hashedFilename)
		//nolint:gosec // G306: chroma CSS is a public file, 0644 is appropriate
		if err := os.WriteFile(hashedPath, []byte(css), 0o644); err != nil {
			return fmt.Errorf("writing hashed chroma CSS: %w", err)
		}
	}

	return nil
}

// generateThemeCSS creates CSS from an explicit Chroma style override.
func (p *ChromaCSSPlugin) generateThemeCSS(style *chroma.Style) (string, error) {
	formatter := chromahtml.New(chromahtml.WithClasses(true), chromahtml.WithAllClasses(true))

	var sb strings.Builder
	sb.WriteString("/* Syntax highlighting - generated from Chroma theme: ")
	sb.WriteString(p.chromaTheme)
	sb.WriteString(" */\n\n")

	// Write the CSS for the style
	if err := formatter.WriteCSS(&sb, style); err != nil {
		return "", err
	}

	return sb.String(), nil
}

// generatePaletteCSS styles Chroma token classes from the active palette's code colors.
func (p *ChromaCSSPlugin) generatePaletteCSS() string {
	var sb strings.Builder
	sb.WriteString("/* Syntax highlighting - generated from palette code colors */\n\n")
	sb.WriteString(`.chroma { color: var(--color-code-text, var(--color-text, currentColor)); background-color: var(--color-code-bg, var(--color-surface, transparent)); }
.chroma .err { color: var(--color-error, #dc2626); background-color: color-mix(in srgb, var(--color-error, #dc2626) 14%, transparent); }
.chroma .lntd { vertical-align: top; padding: 0; margin: 0; border: 0; }
.chroma .lntable { border-spacing: 0; padding: 0; margin: 0; border: 0; }
.chroma .hl { background-color: color-mix(in srgb, var(--color-code-text, currentColor) 10%, var(--color-code-bg, transparent)); }
.chroma .ln, .chroma .lnt { color: color-mix(in srgb, var(--color-code-text, currentColor) 45%, transparent); }
.chroma .c, .chroma .ch, .chroma .cm, .chroma .c1, .chroma .cs, .chroma .cpf { color: var(--color-code-comment, var(--color-text-muted, #6b7280)); font-style: italic; }
.chroma .k, .chroma .kc, .chroma .kd, .chroma .kn, .chroma .kp, .chroma .kr, .chroma .kt { color: var(--color-code-keyword, var(--color-primary, #7c3aed)); font-weight: 600; }
.chroma .s, .chroma .sa, .chroma .sb, .chroma .sc, .chroma .dl, .chroma .sd, .chroma .s2, .chroma .se, .chroma .sh, .chroma .si, .chroma .sx, .chroma .sr, .chroma .s1, .chroma .ss { color: var(--color-code-string, var(--color-success, #059669)); }
.chroma .m, .chroma .mb, .chroma .mf, .chroma .mh, .chroma .mi, .chroma .il, .chroma .mo, .chroma .bin, .chroma .oct, .chroma .hex { color: var(--color-code-number, var(--color-code-keyword, #7c3aed)); }
.chroma .nf, .chroma .fm { color: var(--color-code-function, var(--color-link, #2563eb)); }
.chroma .nc, .chroma .nn, .chroma .no, .chroma .nt, .chroma .nd, .chroma .ne, .chroma .nl { color: var(--color-code-type, var(--color-code-function, #2563eb)); }
.chroma .o, .chroma .ow { color: var(--color-code-operator, var(--color-code-keyword, #7c3aed)); }
.chroma .na, .chroma .py, .chroma .bp { color: var(--color-code-function, var(--color-link, #2563eb)); }
.chroma .n, .chroma .nb, .chroma .nv, .chroma .vc, .chroma .vg, .chroma .vi, .chroma .vm { color: var(--color-code-text, var(--color-text, currentColor)); }
.chroma .gi { color: var(--color-success, #059669); background-color: color-mix(in srgb, var(--color-success, #059669) 14%, transparent); }
.chroma .gd { color: var(--color-error, #dc2626); background-color: color-mix(in srgb, var(--color-error, #dc2626) 14%, transparent); }
`)
	return sb.String()
}

// getExplicitHighlightTheme returns the configured Chroma theme override, if any.
func (p *ChromaCSSPlugin) getExplicitHighlightTheme(extra map[string]interface{}) (string, bool) {
	if extra == nil {
		return "", false
	}

	if highlight, ok := highlightConfigFromExtra(extra); ok && highlight.Theme != "" {
		return highlight.Theme, true
	}

	return "", false
}

// Priority returns the plugin priority for the write stage.
// Should run after static_assets so it can add to the css directory.
func (p *ChromaCSSPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageWrite {
		return lifecycle.PriorityDefault // After static_assets (PriorityEarly)
	}
	return lifecycle.PriorityDefault
}

// Ensure ChromaCSSPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*ChromaCSSPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*ChromaCSSPlugin)(nil)
	_ lifecycle.WritePlugin     = (*ChromaCSSPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*ChromaCSSPlugin)(nil)
)
