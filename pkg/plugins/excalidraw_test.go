package plugins

import (
	"html"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/visualizations"
)

const testExcalidrawTextScene = `{"type":"excalidraw","version":2,"elements":[{"id":"text-1","type":"text","text":"Flow   overview","isDeleted":false}],"appState":{},"files":{}}`

func TestValidateExcalidrawScene(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		wantAlt string
		wantErr string
	}{
		{name: "text alt", source: testExcalidrawTextScene, wantAlt: "Flow overview"},
		{name: "empty scene fallback", source: `{"type":"excalidraw","version":2,"elements":[],"appState":{},"files":{}}`, wantAlt: "Excalidraw diagram with 0 elements"},
		{name: "invalid json", source: `{`, wantErr: "invalid scene JSON"},
		{name: "unsupported version", source: `{"type":"excalidraw","version":3,"elements":[]}`, wantErr: "version 2"},
		{name: "embedded image", source: `{"type":"excalidraw","version":2,"elements":[{"type":"image"}]}`, wantErr: "embedded image"},
		{name: "linked element", source: `{"type":"excalidraw","version":2,"elements":[{"type":"rectangle","link":"https://example.com"}]}`, wantErr: "linked elements"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateExcalidrawScene(tt.source)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("validateExcalidrawScene() error = %v, want substring %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("validateExcalidrawScene() error = %v", err)
			}
			if got != tt.wantAlt {
				t.Errorf("validateExcalidrawScene() alt = %q, want %q", got, tt.wantAlt)
			}
		})
	}
}

func TestExcalidrawDarkMode_FollowsSiteFallback(t *testing.T) {
	light := models.NewConfig()
	light.Theme.FallbackMode = "light"
	dark := models.NewConfig()
	dark.Theme.FallbackMode = "dark"
	tests := []struct {
		name   string
		config *models.Config
		want   bool
	}{
		{name: "no site config defaults dark", want: true},
		{name: "dark site", config: dark, want: true},
		{name: "light site", config: light, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := excalidrawDarkMode(tt.config); got != tt.want {
				t.Errorf("excalidrawDarkMode() = %t, want %t", got, tt.want)
			}
		})
	}
}

func TestExcalidrawVersionForTheme_SeparatesSVGCacheModes(t *testing.T) {
	dark := excalidrawVersionForTheme("renderer-v1", true)
	light := excalidrawVersionForTheme("renderer-v1", false)
	if dark == light {
		t.Fatalf("dark and light renderer versions both equal %q", dark)
	}
	if !strings.HasSuffix(dark, ":dark") || !strings.HasSuffix(light, ":light") {
		t.Errorf("theme versions = %q and %q, want distinct dark/light suffixes", dark, light)
	}
}

func TestReplaceExcalidrawBlocks_RendersStaticFigure(t *testing.T) {
	asset, err := visualizations.NewSVGAsset([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="20" height="10"><rect width="20" height="10"/></svg>`))
	if err != nil {
		t.Fatalf("NewSVGAsset() error = %v", err)
	}
	article := `<pre><code class="language-excalidraw">` + html.EscapeString(testExcalidrawTextScene) + `</code></pre>`
	renderedSource := ""
	updated, assets, err := replaceExcalidrawBlocks(article, func(source string) (visualizations.Asset, error) {
		renderedSource = source
		return asset, nil
	}, func(_ visualizations.Asset, alt string) (string, error) {
		return `<figure>` + html.EscapeString(alt) + `</figure>`, nil
	})
	if err != nil {
		t.Fatalf("replaceExcalidrawBlocks() error = %v", err)
	}
	if renderedSource != testExcalidrawTextScene {
		t.Errorf("renderer source = %q, want %q", renderedSource, testExcalidrawTextScene)
	}
	if updated != `<figure>Flow overview</figure>` {
		t.Errorf("updated article = %q, want a static figure", updated)
	}
	if len(assets) != 1 || assets[0].Path != asset.Path {
		t.Errorf("assets = %#v, want one asset at %q", assets, asset.Path)
	}
}

func TestInlineExcalidrawFigure_UsesSVGDataURI(t *testing.T) {
	asset, err := visualizations.NewSVGAsset([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="20" height="10"><rect width="20" height="10"/></svg>`))
	if err != nil {
		t.Fatalf("NewSVGAsset() error = %v", err)
	}
	markup, err := inlineExcalidrawFigure(asset, `A "diagram"`)
	if err != nil {
		t.Fatalf("inlineExcalidrawFigure() error = %v", err)
	}
	if !strings.Contains(markup, `src="data:image/svg+xml;base64,`) {
		t.Errorf("markup does not use an SVG data URI: %s", markup)
	}
	if !strings.Contains(markup, `alt="A &#34;diagram&#34;"`) {
		t.Errorf("markup does not escape alternative text: %s", markup)
	}
}
