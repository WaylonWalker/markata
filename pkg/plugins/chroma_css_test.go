package plugins

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
)

func TestChromaCSSPlugin_Name(t *testing.T) {
	p := NewChromaCSSPlugin()
	if got := p.Name(); got != "chroma_css" {
		t.Errorf("Name() = %q, want %q", got, "chroma_css")
	}
}

func TestChromaCSSPlugin_Configure(t *testing.T) {
	tests := []struct {
		name         string
		extra        map[string]interface{}
		wantTheme    string
		wantExplicit bool
	}{
		{
			name:         "default theme",
			extra:        map[string]interface{}{},
			wantTheme:    "github-dark",
			wantExplicit: false,
		},
		{
			name: "explicit theme",
			extra: map[string]interface{}{
				"markdown": map[string]interface{}{
					"highlight": map[string]interface{}{
						"theme": "monokai",
					},
				},
			},
			wantTheme:    "monokai",
			wantExplicit: true,
		},
		{
			name: "palette uses palette-native CSS",
			extra: map[string]interface{}{
				"theme": map[string]interface{}{
					"palette": "catppuccin-mocha",
				},
			},
			wantTheme:    "github-dark",
			wantExplicit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewChromaCSSPlugin()
			m := lifecycle.NewManager()
			m.Config().Extra = tt.extra

			err := p.Configure(m)
			if err != nil {
				t.Fatalf("Configure error: %v", err)
			}

			if p.chromaTheme != tt.wantTheme {
				t.Errorf("chromaTheme = %q, want %q", p.chromaTheme, tt.wantTheme)
			}
			if p.explicit != tt.wantExplicit {
				t.Errorf("explicit = %v, want %v", p.explicit, tt.wantExplicit)
			}
		})
	}
}

func TestChromaCSSPlugin_Write(t *testing.T) {
	// Create a temp directory for output
	tmpDir := t.TempDir()

	p := NewChromaCSSPlugin()
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		OutputDir: tmpDir,
		Extra:     map[string]interface{}{},
	})

	// Configure first
	if err := p.Configure(m); err != nil {
		t.Fatalf("Configure error: %v", err)
	}

	// Write the CSS
	if err := p.Write(m); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	// Check the file was created
	cssPath := filepath.Join(tmpDir, "css", "chroma.css")
	content, err := os.ReadFile(cssPath)
	if err != nil {
		t.Fatalf("failed to read chroma.css: %v", err)
	}

	// Verify content
	css := string(content)
	if !strings.Contains(css, "/* Syntax highlighting") {
		t.Error("expected header comment in CSS")
	}
	if !strings.Contains(css, ".chroma") {
		t.Error("expected .chroma class in CSS")
	}
	if !strings.Contains(css, "--color-code-keyword") {
		t.Error("expected palette-native code color variables in CSS")
	}
	// Check for typical chroma classes
	if !strings.Contains(css, ".kd") { // keyword declaration
		t.Error("expected .kd class in CSS")
	}
}

func TestChromaCSSPlugin_Write_DifferentThemes(t *testing.T) {
	tests := []struct {
		name      string
		theme     string
		wantClass string // A class that should exist in the CSS
	}{
		{
			name:      "github-dark",
			theme:     "github-dark",
			wantClass: ".chroma",
		},
		{
			name:      "dracula",
			theme:     "dracula",
			wantClass: ".chroma",
		},
		{
			name:      "monokai",
			theme:     "monokai",
			wantClass: ".chroma",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			p := NewChromaCSSPlugin()
			m := lifecycle.NewManager()
			m.SetConfig(&lifecycle.Config{
				OutputDir: tmpDir,
				Extra: map[string]interface{}{
					"markdown": map[string]interface{}{
						"highlight": map[string]interface{}{
							"theme": tt.theme,
						},
					},
				},
			})

			if err := p.Configure(m); err != nil {
				t.Fatalf("Configure error: %v", err)
			}

			if err := p.Write(m); err != nil {
				t.Fatalf("Write error: %v", err)
			}

			cssPath := filepath.Join(tmpDir, "css", "chroma.css")
			content, err := os.ReadFile(cssPath)
			if err != nil {
				t.Fatalf("failed to read chroma.css: %v", err)
			}

			if !strings.Contains(string(content), tt.wantClass) {
				t.Errorf("expected %q class in CSS for theme %s", tt.wantClass, tt.theme)
			}

			// Verify theme name is in header
			if !strings.Contains(string(content), tt.theme) {
				t.Errorf("expected theme name %q in CSS header", tt.theme)
			}
		})
	}
}

func TestChromaCSSPlugin_Priority(t *testing.T) {
	p := NewChromaCSSPlugin()

	// Should have default priority for write stage
	if got := p.Priority(lifecycle.StageWrite); got != lifecycle.PriorityDefault {
		t.Errorf("Priority(StageWrite) = %d, want %d", got, lifecycle.PriorityDefault)
	}
}
