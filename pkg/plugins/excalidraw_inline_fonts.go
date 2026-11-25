package plugins

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/fs"
	"regexp"
	"strings"

	"github.com/WaylonWalker/markata/pkg/visualizations"
)

const excalidrawFontURLPrefix = "https://esm.sh/@excalidraw/excalidraw@" + visualizations.ExcalidrawPackageVersion + "/dist/prod/fonts/"

var (
	excalidrawStyleElementPattern = regexp.MustCompile(`(?is)(<style(?:\s[^>]*)?>)(.*?)(</style\s*>)`)
	excalidrawCSSURLPattern       = regexp.MustCompile(`(?is)url\(\s*([^)]*?)\s*\)`)
	excalidrawFontPathPattern     = regexp.MustCompile(`^[A-Za-z0-9_.-]+(/[A-Za-z0-9_.-]+)*\.woff2$`)
)

// inlineExcalidrawSVGFonts embeds only fonts listed in Markata's bundled
// Excalidraw font manifest. Other CSS URLs remain subject to SVG validation.
func inlineExcalidrawSVGFonts(svg []byte) ([]byte, error) {
	if len(svg) == 0 || len(svg) > visualizations.MaxSVGBytes {
		return nil, fmt.Errorf("excalidraw SVG must contain 1..%d bytes before font embedding", visualizations.MaxSVGBytes)
	}
	fontPaths, assets, err := excalidrawBundledFontPaths()
	if err != nil {
		return nil, err
	}

	styles := excalidrawStyleElementPattern.FindAllSubmatchIndex(svg, -1)
	if len(styles) == 0 {
		return append([]byte(nil), svg...), nil
	}
	var output bytes.Buffer
	last := 0
	for _, match := range styles {
		output.Write(svg[last:match[0]])
		output.Write(svg[match[2]:match[3]])
		css, err := inlineExcalidrawCSSFonts(svg[match[4]:match[5]], fontPaths, assets)
		if err != nil {
			return nil, err
		}
		output.Write(css)
		output.Write(svg[match[6]:match[7]])
		last = match[1]
		if output.Len()+len(svg)-last > visualizations.MaxSVGBytes {
			return nil, fmt.Errorf("excalidraw SVG exceeds the %d-byte limit after font embedding", visualizations.MaxSVGBytes)
		}
	}
	output.Write(svg[last:])
	if output.Len() > visualizations.MaxSVGBytes {
		return nil, fmt.Errorf("excalidraw SVG exceeds the %d-byte limit after font embedding", visualizations.MaxSVGBytes)
	}
	return output.Bytes(), nil
}

func excalidrawBundledFontPaths() (map[string]struct{}, fs.FS, error) {
	assets, err := visualizations.ExcalidrawAssetFS()
	if err != nil {
		return nil, nil, err
	}
	manifestData, err := fs.ReadFile(assets, "excalidraw-fonts/manifest.json")
	if err != nil {
		return nil, nil, fmt.Errorf("read bundled Excalidraw font manifest: %w", err)
	}
	var manifest struct {
		RendererVersion string   `json:"renderer_version"`
		Fonts           []string `json:"fonts"`
	}
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		return nil, nil, fmt.Errorf("parse bundled Excalidraw font manifest: %w", err)
	}
	if manifest.RendererVersion != visualizations.ExcalidrawSVGExportVersion || len(manifest.Fonts) == 0 {
		return nil, nil, fmt.Errorf("bundled Excalidraw font manifest has an unsupported version or no fonts")
	}
	paths := make(map[string]struct{}, len(manifest.Fonts))
	for _, font := range manifest.Fonts {
		if !fs.ValidPath(font) || !excalidrawFontPathPattern.MatchString(font) {
			return nil, nil, fmt.Errorf("invalid path in bundled Excalidraw font manifest: %q", font)
		}
		if _, exists := paths[font]; exists {
			return nil, nil, fmt.Errorf("duplicate path in bundled Excalidraw font manifest: %q", font)
		}
		paths[font] = struct{}{}
	}
	return paths, assets, nil
}

func inlineExcalidrawCSSFonts(css []byte, fontPaths map[string]struct{}, assets fs.FS) ([]byte, error) {
	matches := excalidrawCSSURLPattern.FindAllSubmatchIndex(css, -1)
	if len(matches) == 0 {
		return append([]byte(nil), css...), nil
	}
	var output bytes.Buffer
	last := 0
	for _, match := range matches {
		output.Write(css[last:match[0]])
		urlValue := strings.Trim(strings.TrimSpace(string(css[match[2]:match[3]])), `"'`)
		if !strings.HasPrefix(urlValue, excalidrawFontURLPrefix) {
			output.Write(css[match[0]:match[1]])
			last = match[1]
			continue
		}
		font := strings.TrimPrefix(urlValue, excalidrawFontURLPrefix)
		if !excalidrawFontPathPattern.MatchString(font) {
			return nil, fmt.Errorf("unsupported Excalidraw font URL %q", urlValue)
		}
		if _, ok := fontPaths[font]; !ok {
			return nil, fmt.Errorf("excalidraw font %q is not listed in the bundled manifest", font)
		}
		fontData, err := fs.ReadFile(assets, "excalidraw-fonts/"+font)
		if err != nil {
			return nil, fmt.Errorf("read bundled Excalidraw font %q: %w", font, err)
		}
		dataURL := "url(data:font/woff2;base64," + base64.StdEncoding.EncodeToString(fontData) + ")"
		output.WriteString(dataURL)
		last = match[1]
		if output.Len()+len(css)-last > visualizations.MaxSVGBytes {
			return nil, fmt.Errorf("excalidraw SVG exceeds the %d-byte limit after font embedding", visualizations.MaxSVGBytes)
		}
	}
	output.Write(css[last:])
	if output.Len() > visualizations.MaxSVGBytes {
		return nil, fmt.Errorf("excalidraw SVG exceeds the %d-byte limit after font embedding", visualizations.MaxSVGBytes)
	}
	return output.Bytes(), nil
}
