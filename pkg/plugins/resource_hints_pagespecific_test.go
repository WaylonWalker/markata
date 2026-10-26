package plugins

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestResourceHintsPageSpecific(t *testing.T) {
	// Create a temporary directory for test output
	tempDir := t.TempDir()

	// Create test HTML files with different external links
	page1 := `<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>Page 1</title>
</head>
<body>
  <a href="https://github.com">GitHub</a>
  <img src="https://cdn.jsdelivr.net/image.png">
</body>
</html>`

	page2 := `<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>Page 2</title>
</head>
<body>
  <a href="https://youtube.com">YouTube</a>
  <a href="https://dev.to">Dev.to</a>
  <script src="https://cdn.tailwindcss.com/script.js"></script>
</body>
</html>`

	page1Path := filepath.Join(tempDir, "page1.html")
	page2Path := filepath.Join(tempDir, "page2.html")

	if err := os.WriteFile(page1Path, []byte(page1), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(page2Path, []byte(page2), 0o600); err != nil {
		t.Fatal(err)
	}

	// Create plugin and configure
	plugin := NewResourceHintsPlugin()

	config := &lifecycle.Config{
		OutputDir: tempDir,
		Extra: map[string]interface{}{
			"resource_hints": models.ResourceHintsConfig{
				Enabled:    boolPtr(true),
				AutoDetect: boolPtr(true),
			},
		},
	}

	manager := &lifecycle.Manager{}
	manager.SetConfig(config)

	if err := plugin.Configure(manager); err != nil {
		t.Fatal(err)
	}

	// Run the plugin
	if err := plugin.Write(manager); err != nil {
		t.Fatal(err)
	}

	// Read modified pages
	page1Content, err := os.ReadFile(page1Path)
	if err != nil {
		t.Fatal(err)
	}
	page2Content, err := os.ReadFile(page2Path)
	if err != nil {
		t.Fatal(err)
	}

	page1Str := string(page1Content)
	page2Str := string(page2Content)

	// Verify page 1 has only its domains
	if !strings.Contains(page1Str, "github.com") {
		t.Error("Page 1 should have github.com hint")
	}
	if !strings.Contains(page1Str, "cdn.jsdelivr.net") {
		t.Error("Page 1 should have cdn.jsdelivr.net hint")
	}
	if strings.Contains(page1Str, "youtube.com") {
		t.Error("Page 1 should NOT have youtube.com hint (that's on page 2)")
	}
	if strings.Contains(page1Str, "dev.to") {
		t.Error("Page 1 should NOT have dev.to hint (that's on page 2)")
	}

	// Verify page 2 has only its domains
	if !strings.Contains(page2Str, "youtube.com") {
		t.Error("Page 2 should have youtube.com hint")
	}
	if !strings.Contains(page2Str, "dev.to") {
		t.Error("Page 2 should have dev.to hint")
	}
	if !strings.Contains(page2Str, "cdn.tailwindcss.com") {
		t.Error("Page 2 should have cdn.tailwindcss.com hint")
	}
	if strings.Contains(page2Str, "github.com") {
		t.Error("Page 2 should NOT have github.com hint (that's on page 1)")
	}
	if strings.Contains(page2Str, "cdn.jsdelivr.net") {
		t.Error("Page 2 should NOT have cdn.jsdelivr.net hint (that's on page 1)")
	}

	// Verify both pages have resource hint comments
	if !strings.Contains(page1Str, "<!-- Auto-generated resource hints -->") {
		t.Error("Page 1 should have resource hints comment")
	}
	if !strings.Contains(page2Str, "<!-- Auto-generated resource hints -->") {
		t.Error("Page 2 should have resource hints comment")
	}

	// Count hints on each page (should be low, not hundreds)
	page1HintCount := strings.Count(page1Str, `rel="dns-prefetch"`)
	page2HintCount := strings.Count(page2Str, `rel="dns-prefetch"`)

	if page1HintCount > 10 {
		t.Errorf("Page 1 has too many hints: %d (expected < 10)", page1HintCount)
	}
	if page2HintCount > 10 {
		t.Errorf("Page 2 has too many hints: %d (expected < 10)", page2HintCount)
	}

	t.Logf("Page 1 hints: %d", page1HintCount)
	t.Logf("Page 2 hints: %d", page2HintCount)
}

func boolPtr(b bool) *bool {
	return &b
}
