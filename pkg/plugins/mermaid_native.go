//go:build !js

package plugins

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	"github.com/WaylonWalker/markata/pkg/models"
	gomermaid "github.com/zkrebbekx/go-mermaid"
)

// This backend does not start a browser, subprocess, or network request.
// It uses the Go library's SVG renderer and its own, bounded render options.

type nativeMermaidRenderer struct {
	config  models.MermaidConfig
	palette *mermaidPaletteColors
}

func availableMermaidModes() []string {
	return []string{mermaidModeClient, mermaidModeGo}
}

func newNativeMermaidRenderer(config models.MermaidConfig, palette *mermaidPaletteColors) mermaidRenderer {
	return &nativeMermaidRenderer{config: config, palette: palette}
}

func (r *nativeMermaidRenderer) render(source string) (string, error) {
	options := []gomermaid.Option{}
	goCfg := r.config.GoConfig
	if goCfg == nil {
		goCfg = &models.GoMermaidRendererConfig{}
	}
	theme := gomermaid.Theme(r.config.Theme)
	if r.palette != nil && r.config.UseCSSVariables {
		bg := r.palette.Background
		if bg == "" {
			bg = defaultMermaidBg
		}
		fill := r.palette.Surface
		if fill == "" {
			fill = defaultMermaidSurface
		}
		textColor := r.palette.TextColor
		if textColor == "" {
			textColor = defaultMermaidText
		}
		accent := r.palette.Accent
		if accent == "" {
			accent = textColor
		}
		palette := gomermaid.Palette{
			Background: bg, NodeFill: fill, NodeStroke: accent,
			Text: textColor, Edge: textColor,
		}
		// go-mermaid registers palettes globally; never use a common mutable
		// name for different site palettes because concurrent builds would
		// mix colors. Hash the resolved semantic colors into the name.
		sum := sha256.Sum256([]byte(strings.Join(
			[]string{bg, fill, accent, textColor}, "\x00")))
		name := "markata-" + hex.EncodeToString(sum[:])
		options = append(options, gomermaid.WithCustomTheme(name, palette))
	} else {
		options = append(options, gomermaid.WithTheme(theme))
	}

	fontFamily := goCfg.FontFamily
	if fontFamily == "" {
		fontFamily = "sans-serif"
	}
	fontSize := goCfg.FontSize
	if fontSize <= 0 {
		fontSize = 16
	}
	options = append(options, gomermaid.WithFont(fontFamily, fontSize))
	if goCfg.Padding > 0 {
		options = append(options, gomermaid.WithPadding(goCfg.Padding))
	}
	if goCfg.NodeSpacing > 0 && goCfg.RankSpacing > 0 {
		options = append(options, gomermaid.WithSpacing(goCfg.NodeSpacing, goCfg.RankSpacing))
	}
	if goCfg.CurvedEdges {
		options = append(options, gomermaid.WithCurvedEdges(true))
	}
	if goCfg.TransparentBackground {
		options = append(options, gomermaid.WithTransparentBackground())
	}

	svg, err := gomermaid.Render(source, options...)
	if err != nil {
		return "", fmt.Errorf("native Mermaid: %w", err)
	}
	if r.palette != nil && r.config.UseCSSVariables {
		svg = nativeMermaidPaletteCSS(svg, r.palette)
	}
	if err := validateNativeMermaidSVG(svg); err != nil {
		return "", fmt.Errorf("reject unsafe native Mermaid output: %w", err)
	}
	return string(svg), nil
}

func (r *nativeMermaidRenderer) close() error { return nil }


// Only rewrite SVG paint attributes, never author-supplied node labels or URLs.
// The fallback preserves the server-rendered appearance until the CSS palette
// loads; CSS variables let the same SVG respond to live light/dark changes.
var nativeMermaidPaintAttr = regexp.MustCompile(`(?i)\b(?:fill|stroke|color|style|stop-color)="[^"]*"`)

func nativeMermaidPaletteCSS(svg []byte, palette *mermaidPaletteColors) []byte {
	if palette == nil {
		return svg
	}
	type paint struct {
		color string
		variable string
	}
	colors := []paint{
		{palette.Background, "--color-background"},
		{palette.Surface, "--color-surface"},
		{palette.TextColor, "--color-text"},
		{palette.Accent, "--color-primary"},
	}
	seen := make(map[string]bool, len(colors))
	replacements := make([]string, 0, len(colors)*2)
	for _, entry := range colors {
		if entry.color == "" || seen[entry.color] {
			continue
		}
		seen[entry.color] = true
		replacements = append(replacements, entry.color,
			"var("+entry.variable+", "+entry.color+")")
	}
	if len(replacements) == 0 {
		return svg
	}
	replacer := strings.NewReplacer(replacements...)
	return []byte(nativeMermaidPaintAttr.ReplaceAllStringFunc(string(svg), replacer.Replace))
}
