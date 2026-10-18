package plugins

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestCSSMinifyPlugin_Name(t *testing.T) {
	p := NewCSSMinifyPlugin()
	if got := p.Name(); got != "css_minify" {
		t.Errorf("Name() = %q, want %q", got, "css_minify")
	}
}

func TestCSSMinifyPlugin_Configure(t *testing.T) {
	tests := []struct {
		name            string
		extra           map[string]interface{}
		wantEnabled     bool
		wantExcludeLen  int
		wantPreserveLen int
	}{
		{
			name:            "default config",
			extra:           nil,
			wantEnabled:     true,
			wantExcludeLen:  0,
			wantPreserveLen: 0,
		},
		{
			name: "disabled",
			extra: map[string]interface{}{
				"css_minify": map[string]interface{}{
					"enabled": false,
				},
			},
			wantEnabled:     false,
			wantExcludeLen:  0,
			wantPreserveLen: 0,
		},
		{
			name: "with exclude patterns",
			extra: map[string]interface{}{
				"css_minify": map[string]interface{}{
					"enabled": true,
					"exclude": []interface{}{"variables.css", "vendor/*.css"},
				},
			},
			wantEnabled:    true,
			wantExcludeLen: 2,
		},
		{
			name: "with preserve comments",
			extra: map[string]interface{}{
				"css_minify": map[string]interface{}{
					"enabled":           true,
					"preserve_comments": []interface{}{"/*! Copyright */", "License"},
				},
			},
			wantEnabled:     true,
			wantPreserveLen: 2,
		},
		{
			name: "with typed config",
			extra: map[string]interface{}{
				"css_minify": models.CSSMinifyConfig{
					Enabled:          true,
					Exclude:          []string{"test.css"},
					PreserveComments: []string{"Copyright"},
				},
			},
			wantEnabled:     true,
			wantExcludeLen:  1,
			wantPreserveLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewCSSMinifyPlugin()
			m := lifecycle.NewManager()
			m.Config().Extra = tt.extra

			err := p.Configure(m)
			if err != nil {
				t.Fatalf("Configure error: %v", err)
			}

			if p.config.Enabled != tt.wantEnabled {
				t.Errorf("Enabled = %v, want %v", p.config.Enabled, tt.wantEnabled)
			}

			if len(p.config.Exclude) != tt.wantExcludeLen {
				t.Errorf("len(Exclude) = %d, want %d", len(p.config.Exclude), tt.wantExcludeLen)
			}

			if tt.wantPreserveLen > 0 && len(p.config.PreserveComments) != tt.wantPreserveLen {
				t.Errorf("len(PreserveComments) = %d, want %d", len(p.config.PreserveComments), tt.wantPreserveLen)
			}
		})
	}
}

func TestCSSMinifyPlugin_Write(t *testing.T) {
	// Create a temp directory for output
	tmpDir := t.TempDir()
	cssDir := filepath.Join(tmpDir, "css")
	if err := os.MkdirAll(cssDir, 0o755); err != nil {
		t.Fatalf("failed to create css dir: %v", err)
	}

	// Create a test CSS file with whitespace and comments
	testCSS := `/* This is a comment */
body {
    margin: 0;
    padding: 0;
    font-family: sans-serif;
}

/* Another comment */
.container {
    max-width: 1200px;
    margin: 0 auto;
}
`
	cssPath := filepath.Join(cssDir, "test.css")
	//nolint:gosec // G306: test file permissions
	if err := os.WriteFile(cssPath, []byte(testCSS), 0o644); err != nil {
		t.Fatalf("failed to write test CSS: %v", err)
	}

	p := NewCSSMinifyPlugin()
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		ContentDir: t.TempDir(),
		OutputDir:  tmpDir,
		Extra: map[string]interface{}{
			"css_minify": map[string]interface{}{
				"enabled": true,
			},
		},
	})

	// Configure
	if err := p.Configure(m); err != nil {
		t.Fatalf("Configure error: %v", err)
	}

	// Write (minify)
	if err := p.Write(m); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	// Read the minified file
	content, err := os.ReadFile(cssPath)
	if err != nil {
		t.Fatalf("failed to read minified CSS: %v", err)
	}

	minifiedCSS := string(content)

	// Verify minification occurred
	if strings.Contains(minifiedCSS, "    ") {
		t.Error("minified CSS should not contain multiple spaces")
	}

	if strings.Contains(minifiedCSS, "\n\n") {
		t.Error("minified CSS should not contain multiple newlines")
	}

	// Verify essential CSS is preserved
	if !strings.Contains(minifiedCSS, "body") {
		t.Error("minified CSS should contain 'body' selector")
	}

	if !strings.Contains(minifiedCSS, "margin:0") || !strings.Contains(minifiedCSS, "margin: 0") {
		// Check for either format since minifier may use different formats
		if !strings.Contains(minifiedCSS, "margin") {
			t.Error("minified CSS should contain 'margin' property")
		}
	}

	// Verify size reduction
	if len(minifiedCSS) >= len(testCSS) {
		t.Errorf("minified CSS (%d bytes) should be smaller than original (%d bytes)",
			len(minifiedCSS), len(testCSS))
	}
}

func TestCSSMinifyPlugin_Write_Disabled(t *testing.T) {
	tmpDir := t.TempDir()
	cssDir := filepath.Join(tmpDir, "css")
	if err := os.MkdirAll(cssDir, 0o755); err != nil {
		t.Fatalf("failed to create css dir: %v", err)
	}

	// Create a test CSS file
	testCSS := `body { margin: 0; }`
	cssPath := filepath.Join(cssDir, "test.css")
	//nolint:gosec // G306: test file permissions
	if err := os.WriteFile(cssPath, []byte(testCSS), 0o644); err != nil {
		t.Fatalf("failed to write test CSS: %v", err)
	}

	p := NewCSSMinifyPlugin()
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		ContentDir: t.TempDir(),
		OutputDir:  tmpDir,
		Extra: map[string]interface{}{
			"css_minify": map[string]interface{}{
				"enabled": false,
			},
		},
	})

	if err := p.Configure(m); err != nil {
		t.Fatalf("Configure error: %v", err)
	}

	if err := p.Write(m); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	// Read the file - should be unchanged
	content, err := os.ReadFile(cssPath)
	if err != nil {
		t.Fatalf("failed to read CSS: %v", err)
	}

	if string(content) != testCSS {
		t.Error("CSS should be unchanged when plugin is disabled")
	}
}

func TestCSSMinifyPlugin_Write_Exclude(t *testing.T) {
	tmpDir := t.TempDir()
	cssDir := filepath.Join(tmpDir, "css")
	if err := os.MkdirAll(cssDir, 0o755); err != nil {
		t.Fatalf("failed to create css dir: %v", err)
	}

	// Create test CSS files
	testCSS := `/* Comment */
body {
    margin: 0;
}
`
	regularPath := filepath.Join(cssDir, "regular.css")
	excludedPath := filepath.Join(cssDir, "variables.css")

	//nolint:gosec // G306: test file permissions
	if err := os.WriteFile(regularPath, []byte(testCSS), 0o644); err != nil {
		t.Fatalf("failed to write regular CSS: %v", err)
	}
	//nolint:gosec // G306: test file permissions
	if err := os.WriteFile(excludedPath, []byte(testCSS), 0o644); err != nil {
		t.Fatalf("failed to write excluded CSS: %v", err)
	}

	p := NewCSSMinifyPlugin()
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		ContentDir: t.TempDir(),
		OutputDir:  tmpDir,
		Extra: map[string]interface{}{
			"css_minify": map[string]interface{}{
				"enabled": true,
				"exclude": []interface{}{"variables.css"},
			},
		},
	})

	if err := p.Configure(m); err != nil {
		t.Fatalf("Configure error: %v", err)
	}

	if err := p.Write(m); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	// Check regular file was minified
	regularContent, err := os.ReadFile(regularPath)
	if err != nil {
		t.Fatalf("failed to read regular CSS: %v", err)
	}
	if len(regularContent) >= len(testCSS) {
		t.Error("regular.css should be minified")
	}

	// Check excluded file was not minified
	excludedContent, err := os.ReadFile(excludedPath)
	if err != nil {
		t.Fatalf("failed to read excluded CSS: %v", err)
	}
	if string(excludedContent) != testCSS {
		t.Error("variables.css should not be minified (excluded)")
	}
}

func TestCSSMinifyPlugin_IsExcluded_PagefindAssets(t *testing.T) {
	p := NewCSSMinifyPlugin()
	if !p.isExcluded("/tmp/output/_pagefind/pagefind-ui.css") {
		t.Fatal("expected _pagefind/pagefind-ui.css to be excluded")
	}
	if !p.isExcluded("/tmp/output/css/fonts.css") {
		t.Fatal("expected generated font CSS to be excluded")
	}
	if p.isExcluded("/tmp/output/other/fonts.css") || p.isExcluded("/tmp/output/fonts.css") {
		t.Fatal("did not expect unrelated files named fonts.css to be excluded")
	}
	if p.isExcluded("/tmp/output/css/main.css") {
		t.Fatal("did not expect normal CSS asset to be excluded")
	}
}

func TestCSSMinifyPlugin_Write_PreserveComments(t *testing.T) {
	tmpDir := t.TempDir()
	cssDir := filepath.Join(tmpDir, "css")
	if err := os.MkdirAll(cssDir, 0o755); err != nil {
		t.Fatalf("failed to create css dir: %v", err)
	}

	// Create test CSS with copyright comment
	testCSS := `/*! Copyright 2024 Test Inc. - Do not remove */
/* Regular comment to remove */
body {
    margin: 0;
    padding: 0;
}
`
	cssPath := filepath.Join(cssDir, "test.css")
	//nolint:gosec // G306: test file permissions
	if err := os.WriteFile(cssPath, []byte(testCSS), 0o644); err != nil {
		t.Fatalf("failed to write test CSS: %v", err)
	}

	p := NewCSSMinifyPlugin()
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		ContentDir: t.TempDir(),
		OutputDir:  tmpDir,
		Extra: map[string]interface{}{
			"css_minify": map[string]interface{}{
				"enabled":           true,
				"preserve_comments": []interface{}{"Copyright"},
			},
		},
	})

	if err := p.Configure(m); err != nil {
		t.Fatalf("Configure error: %v", err)
	}

	if err := p.Write(m); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	// Read the minified file
	content, err := os.ReadFile(cssPath)
	if err != nil {
		t.Fatalf("failed to read minified CSS: %v", err)
	}

	minifiedCSS := string(content)

	// Verify copyright comment is preserved
	if !strings.Contains(minifiedCSS, "Copyright 2024 Test Inc.") {
		t.Error("copyright comment should be preserved")
	}
}

func TestCSSMinifyPlugin_Priority(t *testing.T) {
	p := NewCSSMinifyPlugin()

	// Should have PriorityLast for write stage (runs after all other CSS plugins)
	if got := p.Priority(lifecycle.StageWrite); got != lifecycle.PriorityLast {
		t.Errorf("Priority(StageWrite) = %d, want %d", got, lifecycle.PriorityLast)
	}

	// Other stages should have default priority
	if got := p.Priority(lifecycle.StageRender); got != lifecycle.PriorityDefault {
		t.Errorf("Priority(StageRender) = %d, want %d", got, lifecycle.PriorityDefault)
	}
}

func TestCSSMinifyPlugin_SizeReduction(t *testing.T) {
	tmpDir := t.TempDir()
	cssDir := filepath.Join(tmpDir, "css")
	if err := os.MkdirAll(cssDir, 0o755); err != nil {
		t.Fatalf("failed to create css dir: %v", err)
	}

	// Create a larger CSS file to test meaningful reduction
	largeCSS := `
/* =========================================
   Main Stylesheet
   ========================================= */

/* Reset styles */
html, body, div, span, applet, object, iframe,
h1, h2, h3, h4, h5, h6, p, blockquote, pre,
a, abbr, acronym, address, big, cite, code,
del, dfn, em, img, ins, kbd, q, s, samp,
small, strike, strong, sub, sup, tt, var,
b, u, i, center {
    margin: 0;
    padding: 0;
    border: 0;
    font-size: 100%;
    font: inherit;
    vertical-align: baseline;
}

/* Body styles */
body {
    line-height: 1;
    background-color: #ffffff;
    color: #333333;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
}

/* Container */
.container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px;
}

/* Typography */
h1 {
    font-size: 2.5rem;
    font-weight: 700;
    margin-bottom: 1rem;
}

h2 {
    font-size: 2rem;
    font-weight: 600;
    margin-bottom: 0.875rem;
}

h3 {
    font-size: 1.5rem;
    font-weight: 600;
    margin-bottom: 0.75rem;
}

p {
    margin-bottom: 1rem;
    line-height: 1.6;
}

/* Links */
a {
    color: #0066cc;
    text-decoration: none;
    transition: color 0.2s ease;
}

a:hover {
    color: #004499;
    text-decoration: underline;
}

/* Buttons */
.button {
    display: inline-block;
    padding: 12px 24px;
    background-color: #0066cc;
    color: #ffffff;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    transition: background-color 0.2s ease;
}

.button:hover {
    background-color: #004499;
}

/* Card component */
.card {
    background-color: #ffffff;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    padding: 24px;
    margin-bottom: 24px;
}

.card-title {
    font-size: 1.25rem;
    font-weight: 600;
    margin-bottom: 12px;
}

.card-content {
    color: #666666;
    line-height: 1.5;
}
`
	cssPath := filepath.Join(cssDir, "main.css")
	//nolint:gosec // G306: test file permissions
	if err := os.WriteFile(cssPath, []byte(largeCSS), 0o644); err != nil {
		t.Fatalf("failed to write test CSS: %v", err)
	}

	originalSize := len(largeCSS)

	p := NewCSSMinifyPlugin()
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		ContentDir: t.TempDir(),
		OutputDir:  tmpDir,
		Extra: map[string]interface{}{
			"css_minify": map[string]interface{}{
				"enabled": true,
			},
		},
	})

	if err := p.Configure(m); err != nil {
		t.Fatalf("Configure error: %v", err)
	}

	if err := p.Write(m); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	// Read minified content
	content, err := os.ReadFile(cssPath)
	if err != nil {
		t.Fatalf("failed to read minified CSS: %v", err)
	}

	minifiedSize := len(content)
	reduction := float64(originalSize-minifiedSize) / float64(originalSize) * 100

	// Verify meaningful reduction (should be at least 15%)
	if reduction < 15 {
		t.Errorf("CSS reduction = %.1f%%, want at least 15%%", reduction)
	}

	t.Logf("CSS minification: %d -> %d bytes (%.1f%% reduction)", originalSize, minifiedSize, reduction)
}

func TestCSSMinifyPlugin_GlobExclude(t *testing.T) {
	tmpDir := t.TempDir()
	cssDir := filepath.Join(tmpDir, "css")
	if err := os.MkdirAll(cssDir, 0o755); err != nil {
		t.Fatalf("failed to create css dir: %v", err)
	}

	testCSS := `/* Test */
body { margin: 0; }
`

	// Create test files
	files := []string{"main.css", "vendor-lib.css", "vendor-other.css"}
	for _, f := range files {
		//nolint:gosec // G306: test file permissions
		if err := os.WriteFile(filepath.Join(cssDir, f), []byte(testCSS), 0o644); err != nil {
			t.Fatalf("failed to write %s: %v", f, err)
		}
	}

	p := NewCSSMinifyPlugin()
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		ContentDir: t.TempDir(),
		OutputDir:  tmpDir,
		Extra: map[string]interface{}{
			"css_minify": map[string]interface{}{
				"enabled": true,
				"exclude": []interface{}{"vendor-*.css"},
			},
		},
	})

	if err := p.Configure(m); err != nil {
		t.Fatalf("Configure error: %v", err)
	}

	if err := p.Write(m); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	// main.css should be minified
	mainContent, err := os.ReadFile(filepath.Join(cssDir, "main.css"))
	if err != nil {
		t.Fatalf("failed to read main.css: %v", err)
	}
	if len(mainContent) >= len(testCSS) {
		t.Error("main.css should be minified")
	}

	// vendor files should not be minified
	for _, f := range []string{"vendor-lib.css", "vendor-other.css"} {
		content, err := os.ReadFile(filepath.Join(cssDir, f))
		if err != nil {
			t.Fatalf("failed to read %s: %v", f, err)
		}
		if string(content) != testCSS {
			t.Errorf("%s should not be minified (excluded by glob pattern)", f)
		}
	}
}
