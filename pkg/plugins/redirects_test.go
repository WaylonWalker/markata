package plugins

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// TestRedirectsPlugin_Name tests the plugin name.
func TestRedirectsPlugin_Name(t *testing.T) {
	p := NewRedirectsPlugin()
	if got := p.Name(); got != "redirects" {
		t.Errorf("Name() = %q, want %q", got, "redirects")
	}
}

// TestRedirectsPlugin_ParseRedirects tests redirect file parsing.
func TestRedirectsPlugin_ParseRedirects(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []Redirect
	}{
		{
			name: "basic redirects",
			content: `/old-path    /new-path
/legacy-url  /current-url`,
			want: []Redirect{
				{Original: "/old-path", New: "/new-path"},
				{Original: "/legacy-url", New: "/current-url"},
			},
		},
		{
			name: "with comments",
			content: `# This is a comment
/old    /new
# Another comment
/foo    /bar`,
			want: []Redirect{
				{Original: "/old", New: "/new"},
				{Original: "/foo", New: "/bar"},
			},
		},
		{
			name: "skip wildcards",
			content: `/old/*    /new/*
/valid    /path
/another/*  /dest/*`,
			want: []Redirect{
				{Original: "/valid", New: "/path"},
			},
		},
		{
			name: "skip malformed lines",
			content: `/only-source
/valid    /destination
incomplete
/another  /target`,
			want: []Redirect{
				{Original: "/valid", New: "/destination"},
				{Original: "/another", New: "/target"},
			},
		},
		{
			name: "skip non-absolute paths",
			content: `old-path    /new-path
/valid    /destination
/source   relative-dest`,
			want: []Redirect{
				{Original: "/valid", New: "/destination"},
			},
		},
		{
			name: "external URL redirects",
			content: `/go/github    https://github.com/myorg
/go/docs      http://docs.example.com
/internal     /local-path`,
			want: []Redirect{
				{Original: "/go/github", New: "https://github.com/myorg"},
				{Original: "/go/docs", New: "http://docs.example.com"},
				{Original: "/internal", New: "/local-path"},
			},
		},
		{
			name:    "empty content",
			content: "",
			want:    nil,
		},
		{
			name: "only comments and whitespace",
			content: `# Comment 1

# Comment 2
   `,
			want: nil,
		},
		{
			name: "extra whitespace handling",
			content: `   /old    /new
	/tabbed	/path	`,
			want: []Redirect{
				{Original: "/old", New: "/new"},
				{Original: "/tabbed", New: "/path"},
			},
		},
		{
			name: "multiple spaces between paths",
			content: `/source      /destination
/a          /b`,
			want: []Redirect{
				{Original: "/source", New: "/destination"},
				{Original: "/a", New: "/b"},
			},
		},
		{
			name: "extra fields ignored",
			content: `/old    /new    301
/foo    /bar    302    extra`,
			want: []Redirect{
				{Original: "/old", New: "/new"},
				{Original: "/foo", New: "/bar"},
			},
		},
	}

	p := NewRedirectsPlugin()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.parseRedirects(tt.content)

			if len(got) != len(tt.want) {
				t.Errorf("parseRedirects() returned %d redirects, want %d", len(got), len(tt.want))
				t.Errorf("got: %+v", got)
				t.Errorf("want: %+v", tt.want)
				return
			}

			for i, r := range got {
				if r.Original != tt.want[i].Original {
					t.Errorf("redirect[%d].Original = %q, want %q", i, r.Original, tt.want[i].Original)
				}
				if r.New != tt.want[i].New {
					t.Errorf("redirect[%d].New = %q, want %q", i, r.New, tt.want[i].New)
				}
			}
		})
	}
}

// TestRedirectsPlugin_Write tests the full write workflow.
func TestRedirectsPlugin_Write(t *testing.T) {
	// Create temporary directories
	tmpDir := t.TempDir()
	staticDir := filepath.Join(tmpDir, "static")
	outputDir := filepath.Join(tmpDir, "output")
	_ = os.MkdirAll(staticDir, 0o755) //nolint:errcheck // test setup
	_ = os.MkdirAll(outputDir, 0o755) //nolint:errcheck // test setup

	// Create redirects file
	redirectsContent := `# Blog migration
/old-post    /new-post
/legacy      /current
/blog/2023/article    /posts/article
`
	redirectsFile := filepath.Join(staticDir, "_redirects")
	//nolint:gosec // test file
	if err := os.WriteFile(redirectsFile, []byte(redirectsContent), 0o644); err != nil {
		t.Fatalf("failed to create redirects file: %v", err)
	}

	// Setup manager
	m := lifecycle.NewManager()
	cfg := m.Config()
	cfg.OutputDir = outputDir

	// Create and configure plugin
	p := NewRedirectsPlugin()
	p.SetConfig(RedirectsConfig{
		RedirectsFile: redirectsFile,
	})

	// Run write
	if err := p.Write(m); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Verify output files
	expectedFiles := []struct {
		path     string
		contains []string
	}{
		{
			path: filepath.Join(outputDir, "old-post", "index.html"),
			contains: []string{
				`url='/new-post'`,
				`href="/new-post"`,
				`/old-post`,
				"Page Moved",
			},
		},
		{
			path: filepath.Join(outputDir, "legacy", "index.html"),
			contains: []string{
				`url='/current'`,
				`href="/current"`,
			},
		},
		{
			path: filepath.Join(outputDir, "blog", "2023", "article", "index.html"),
			contains: []string{
				`url='/posts/article'`,
				`href="/posts/article"`,
			},
		},
	}

	for _, ef := range expectedFiles {
		content, err := os.ReadFile(ef.path)
		if err != nil {
			t.Errorf("expected file %s not found: %v", ef.path, err)
			continue
		}

		contentStr := string(content)
		for _, substr := range ef.contains {
			if !strings.Contains(contentStr, substr) {
				t.Errorf("file %s missing expected content %q", ef.path, substr)
			}
		}

		// Verify it's valid HTML
		if !strings.Contains(contentStr, "<!DOCTYPE html>") {
			t.Errorf("file %s missing DOCTYPE", ef.path)
		}
	}
}

// TestRedirectsPlugin_Write_MissingFile tests handling of missing redirects file.
func TestRedirectsPlugin_Write_MissingFile(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "output")
	_ = os.MkdirAll(outputDir, 0o755) //nolint:errcheck // test setup

	m := lifecycle.NewManager()
	cfg := m.Config()
	cfg.OutputDir = outputDir

	p := NewRedirectsPlugin()
	p.SetConfig(RedirectsConfig{
		RedirectsFile: filepath.Join(tmpDir, "nonexistent", "_redirects"),
	})

	// Should not error on missing file
	if err := p.Write(m); err != nil {
		t.Errorf("Write() should not error on missing file, got: %v", err)
	}
}

// TestRedirectsPlugin_Write_EmptyFile tests handling of empty redirects file.
func TestRedirectsPlugin_Write_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "output")
	_ = os.MkdirAll(outputDir, 0o755) //nolint:errcheck // test setup

	// Create empty redirects file
	redirectsFile := filepath.Join(tmpDir, "_redirects")
	//nolint:gosec // test file
	if err := os.WriteFile(redirectsFile, []byte(""), 0o644); err != nil {
		t.Fatalf("failed to create redirects file: %v", err)
	}

	m := lifecycle.NewManager()
	cfg := m.Config()
	cfg.OutputDir = outputDir

	p := NewRedirectsPlugin()
	p.SetConfig(RedirectsConfig{
		RedirectsFile: redirectsFile,
	})

	if err := p.Write(m); err != nil {
		t.Errorf("Write() error on empty file = %v", err)
	}

	// Output dir should have no redirect files
	entries, _ := os.ReadDir(outputDir) //nolint:errcheck // test code
	if len(entries) > 0 {
		t.Errorf("expected no files in output dir for empty redirects, got %d entries", len(entries))
	}
}

// TestRedirectsPlugin_Write_CustomTemplate tests custom template support.
func TestRedirectsPlugin_Write_CustomTemplate(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "output")
	_ = os.MkdirAll(outputDir, 0o755) //nolint:errcheck // test setup

	// Create redirects file
	redirectsFile := filepath.Join(tmpDir, "_redirects")
	//nolint:gosec // test file
	if err := os.WriteFile(redirectsFile, []byte("/old /new"), 0o644); err != nil {
		t.Fatalf("failed to create redirects file: %v", err)
	}

	// Create custom template
	customTemplate := `<!DOCTYPE html>
<html>
<head><meta http-equiv="Refresh" content="0; url='{{ .New }}'" /></head>
<body>CUSTOM: {{ .Original }} -> {{ .New }}</body>
</html>`
	templateFile := filepath.Join(tmpDir, "redirect.html")
	//nolint:gosec // test file
	if err := os.WriteFile(templateFile, []byte(customTemplate), 0o644); err != nil {
		t.Fatalf("failed to create template file: %v", err)
	}

	m := lifecycle.NewManager()
	cfg := m.Config()
	cfg.OutputDir = outputDir

	p := NewRedirectsPlugin()
	p.SetConfig(RedirectsConfig{
		RedirectsFile:    redirectsFile,
		RedirectTemplate: templateFile,
	})

	if err := p.Write(m); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Verify custom template was used
	content, err := os.ReadFile(filepath.Join(outputDir, "old", "index.html"))
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	if !strings.Contains(string(content), "CUSTOM: /old -> /new") {
		t.Errorf("custom template not used, got: %s", string(content))
	}
}

// TestRedirectsPlugin_Configure tests configuration loading.
func TestRedirectsPlugin_Configure(t *testing.T) {
	tests := []struct {
		name     string
		extra    map[string]interface{}
		wantFile string
		wantTmpl string
	}{
		{
			name:     "default config",
			extra:    nil,
			wantFile: "static/_redirects",
			wantTmpl: "",
		},
		{
			name: "top-level redirects config",
			extra: map[string]interface{}{
				"redirects": map[string]interface{}{
					"redirects_file":    "custom/_redirects",
					"redirect_template": "custom/redirect.html",
				},
			},
			wantFile: "custom/_redirects",
			wantTmpl: "custom/redirect.html",
		},
		{
			name: "nested markata-go config",
			extra: map[string]interface{}{
				"markata-go": map[string]interface{}{
					"redirects": map[string]interface{}{
						"redirects_file": "nested/_redirects",
					},
				},
			},
			wantFile: "nested/_redirects",
			wantTmpl: "",
		},
		{
			name: "markata-go takes precedence",
			extra: map[string]interface{}{
				"redirects": map[string]interface{}{
					"redirects_file": "top-level/_redirects",
				},
				"markata-go": map[string]interface{}{
					"redirects": map[string]interface{}{
						"redirects_file": "nested/_redirects",
					},
				},
			},
			wantFile: "nested/_redirects",
			wantTmpl: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := lifecycle.NewManager()
			cfg := m.Config()
			cfg.Extra = tt.extra

			p := NewRedirectsPlugin()
			if err := p.Configure(m); err != nil {
				t.Fatalf("Configure() error = %v", err)
			}

			got := p.GetConfig()
			if got.RedirectsFile != tt.wantFile {
				t.Errorf("RedirectsFile = %q, want %q", got.RedirectsFile, tt.wantFile)
			}
			if got.RedirectTemplate != tt.wantTmpl {
				t.Errorf("RedirectTemplate = %q, want %q", got.RedirectTemplate, tt.wantTmpl)
			}
		})
	}
}

// TestRedirectsPlugin_Priority tests the plugin priority.
func TestRedirectsPlugin_Priority(t *testing.T) {
	p := NewRedirectsPlugin()

	if got := p.Priority(lifecycle.StageWrite); got != lifecycle.PriorityLate {
		t.Errorf("Priority(StageWrite) = %d, want %d (PriorityLate)", got, lifecycle.PriorityLate)
	}

	if got := p.Priority(lifecycle.StageRender); got != lifecycle.PriorityDefault {
		t.Errorf("Priority(StageRender) = %d, want %d (PriorityDefault)", got, lifecycle.PriorityDefault)
	}
}

// TestRedirectsPlugin_Write_RecreatesMissingFallback checks that a missing
// generated page is restored while its source rule still exists.
func TestRedirectsPlugin_Write_RecreatesMissingFallback(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "output")
	_ = os.MkdirAll(outputDir, 0o755) //nolint:errcheck // test setup

	// Create redirects file
	redirectsFile := filepath.Join(tmpDir, "_redirects")
	//nolint:gosec // test file
	if err := os.WriteFile(redirectsFile, []byte("/old /new"), 0o644); err != nil {
		t.Fatalf("failed to create redirects file: %v", err)
	}

	m := lifecycle.NewManager()
	cfg := m.Config()
	cfg.OutputDir = outputDir

	p := NewRedirectsPlugin()
	p.SetConfig(RedirectsConfig{
		RedirectsFile: redirectsFile,
	})

	// First write
	if err := p.Write(m); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Verify file was created
	outputPath := filepath.Join(outputDir, "old", "index.html")
	if _, err := os.Stat(outputPath); err != nil {
		t.Fatalf("expected file %s not found after first write", outputPath)
	}

	// Delete the output file
	if err := os.Remove(outputPath); err != nil {
		t.Fatalf("remove generated page: %v", err)
	}

	// The next write restores the generated page from the source rule.
	if err := p.Write(m); err != nil {
		t.Fatalf("Write() error on second call = %v", err)
	}

	if _, err := os.Stat(outputPath); err != nil {
		t.Errorf("file %s was not restored: %v", outputPath, err)
	}
}

// TestRedirectsPlugin_Write_NestedPaths tests handling of deeply nested paths.
func TestRedirectsPlugin_Write_NestedPaths(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "output")
	_ = os.MkdirAll(outputDir, 0o755) //nolint:errcheck // test setup

	// Create redirects file with nested paths
	redirectsContent := `/a/b/c/d/e    /target
/deep/nested/path/here    /simple`
	redirectsFile := filepath.Join(tmpDir, "_redirects")
	//nolint:gosec // test file
	if err := os.WriteFile(redirectsFile, []byte(redirectsContent), 0o644); err != nil {
		t.Fatalf("failed to create redirects file: %v", err)
	}

	m := lifecycle.NewManager()
	cfg := m.Config()
	cfg.OutputDir = outputDir

	p := NewRedirectsPlugin()
	p.SetConfig(RedirectsConfig{
		RedirectsFile: redirectsFile,
	})

	if err := p.Write(m); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Verify nested directories were created
	expectedPaths := []string{
		filepath.Join(outputDir, "a", "b", "c", "d", "e", "index.html"),
		filepath.Join(outputDir, "deep", "nested", "path", "here", "index.html"),
	}

	for _, path := range expectedPaths {
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected file %s not found", path)
		}
	}
}

// TestRedirect_Struct tests the Redirect struct.
func TestRedirect_Struct(t *testing.T) {
	r := Redirect{
		Original: "/old-path",
		New:      "/new-path",
	}

	if r.Original != "/old-path" {
		t.Errorf("Original = %q, want %q", r.Original, "/old-path")
	}
	if r.New != "/new-path" {
		t.Errorf("New = %q, want %q", r.New, "/new-path")
	}
}

// TestHashContent tests the hash function.
func TestHashContent(t *testing.T) {
	// Same content should produce same hash
	h1 := hashContent([]byte("test content"))
	h2 := hashContent([]byte("test content"))
	if h1 != h2 {
		t.Errorf("same content produced different hashes: %d != %d", h1, h2)
	}

	// Different content should produce different hash
	h3 := hashContent([]byte("different content"))
	if h1 == h3 {
		t.Errorf("different content produced same hash: %d", h1)
	}

	// Empty content should work
	h4 := hashContent([]byte(""))
	if h4 == 0 {
		t.Errorf("empty content produced zero hash")
	}
}

// TestRedirectsPlugin_Write_ExternalURLs tests external URL redirects.
func TestRedirectsPlugin_Write_ExternalURLs(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "output")
	_ = os.MkdirAll(outputDir, 0o755) //nolint:errcheck // test setup

	// Create redirects file with external URLs
	redirectsContent := `/go/github    https://github.com/myorg
/go/docs      http://docs.example.com/guide
`
	redirectsFile := filepath.Join(tmpDir, "_redirects")
	//nolint:gosec // test file
	if err := os.WriteFile(redirectsFile, []byte(redirectsContent), 0o644); err != nil {
		t.Fatalf("failed to create redirects file: %v", err)
	}

	m := lifecycle.NewManager()
	cfg := m.Config()
	cfg.OutputDir = outputDir

	p := NewRedirectsPlugin()
	p.SetConfig(RedirectsConfig{
		RedirectsFile: redirectsFile,
	})

	if err := p.Write(m); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Verify external URL redirects were created
	expectedFiles := []struct {
		path     string
		contains []string
	}{
		{
			path: filepath.Join(outputDir, "go", "github", "index.html"),
			contains: []string{
				`url='https://github.com/myorg'`,
				`href="https://github.com/myorg"`,
			},
		},
		{
			path: filepath.Join(outputDir, "go", "docs", "index.html"),
			contains: []string{
				`url='http://docs.example.com/guide'`,
				`href="http://docs.example.com/guide"`,
			},
		},
	}

	for _, ef := range expectedFiles {
		content, err := os.ReadFile(ef.path)
		if err != nil {
			t.Errorf("expected file %s not found: %v", ef.path, err)
			continue
		}

		contentStr := string(content)
		for _, substr := range ef.contains {
			if !strings.Contains(contentStr, substr) {
				t.Errorf("file %s missing expected content %q", ef.path, substr)
			}
		}
	}
}

// TestRedirectsPlugin_Write_PathTraversal tests that path traversal is blocked.
func TestRedirectsPlugin_Write_PathTraversal(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "output")
	_ = os.MkdirAll(outputDir, 0o755) //nolint:errcheck // test setup

	// Create redirects file with path traversal attempt
	redirectsContent := `/../../etc/passwd    /hacked
/valid/path    /destination
`
	redirectsFile := filepath.Join(tmpDir, "_redirects")
	//nolint:gosec // test file
	if err := os.WriteFile(redirectsFile, []byte(redirectsContent), 0o644); err != nil {
		t.Fatalf("failed to create redirects file: %v", err)
	}

	m := lifecycle.NewManager()
	cfg := m.Config()
	cfg.OutputDir = outputDir

	p := NewRedirectsPlugin()
	p.SetConfig(RedirectsConfig{
		RedirectsFile: redirectsFile,
	})

	// Write should not error (it logs warnings for path traversal)
	if err := p.Write(m); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Valid path should be created
	validPath := filepath.Join(outputDir, "valid", "path", "index.html")
	if _, err := os.Stat(validPath); err != nil {
		t.Errorf("expected valid path %s to be created", validPath)
	}

	// Path traversal should NOT create a file outside output dir
	// Check that nothing was created at parent directories
	badPath := filepath.Join(outputDir, "..", "etc", "passwd", "index.html")
	if _, err := os.Stat(badPath); err == nil {
		t.Errorf("path traversal should have been blocked, but %s exists", badPath)
	}
}

// TestRedirectsPlugin_Write_TemplateChangeInvalidatesCache tests that template changes invalidate cache.
func TestRedirectsPlugin_Write_TemplateChangeInvalidatesCache(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "output")
	_ = os.MkdirAll(outputDir, 0o755) //nolint:errcheck // test setup

	// Create redirects file
	redirectsFile := filepath.Join(tmpDir, "_redirects")
	//nolint:gosec // test file
	if err := os.WriteFile(redirectsFile, []byte("/old /new"), 0o644); err != nil {
		t.Fatalf("failed to create redirects file: %v", err)
	}

	// Create custom template v1
	templateV1 := `<!DOCTYPE html><html><body>TEMPLATE_V1</body></html>`
	templateFile := filepath.Join(tmpDir, "redirect.html")
	//nolint:gosec // test file
	if err := os.WriteFile(templateFile, []byte(templateV1), 0o644); err != nil {
		t.Fatalf("failed to create template file: %v", err)
	}

	m := lifecycle.NewManager()
	cfg := m.Config()
	cfg.OutputDir = outputDir

	p := NewRedirectsPlugin()
	p.SetConfig(RedirectsConfig{
		RedirectsFile:    redirectsFile,
		RedirectTemplate: templateFile,
	})

	// First write with template v1
	if err := p.Write(m); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	outputPath := filepath.Join(outputDir, "old", "index.html")
	content1, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	if !strings.Contains(string(content1), "TEMPLATE_V1") {
		t.Errorf("first write should use template v1")
	}

	// Update template to v2
	templateV2 := `<!DOCTYPE html><html><body>TEMPLATE_V2</body></html>`
	//nolint:gosec // test file
	if err := os.WriteFile(templateFile, []byte(templateV2), 0o644); err != nil {
		t.Fatalf("failed to update template file: %v", err)
	}

	// Second write - cache should be invalidated due to template change
	if err := p.Write(m); err != nil {
		t.Fatalf("Write() error on second call = %v", err)
	}

	content2, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output file after second write: %v", err)
	}
	if !strings.Contains(string(content2), "TEMPLATE_V2") {
		t.Errorf("second write should use template v2 (cache should be invalidated), got: %s", string(content2))
	}
}

func TestRedirectsPluginPreservesAuthoredHTMLAndNativeRules(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "_redirects")
	if err := os.WriteFile(file, []byte("/article /missing\n/legacy /new\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(dir, "output")
	page := filepath.Join(output, "article", "index.html")
	if err := os.MkdirAll(filepath.Dir(page), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(page, []byte("published article"), 0o600); err != nil {
		t.Fatal(err)
	}
	m := lifecycle.NewManager()
	m.Config().OutputDir = output
	p := NewRedirectsPlugin()
	p.SetConfig(RedirectsConfig{RedirectsFile: file})
	if err := p.Write(m); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(page)
	if err != nil || string(got) != "published article" {
		t.Fatalf("redirect overwrote authored HTML: %q, %v", got, err)
	}
	native, err := os.ReadFile(filepath.Join(output, "redirects.conf"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(native), "/article") || !strings.Contains(string(native), "/legacy") {
		t.Fatalf("nginx rules can shadow authored page: %s", native)
	}
	if _, err := os.Stat(filepath.Join(output, "legacy", "index.html")); err != nil {
		t.Fatal(err)
	}
}

func TestRedirectsPluginProtectsPublishedRoutesBeforeHTMLExists(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "_redirects")
	if err := os.WriteFile(file, []byte("/published /missing\n/old /new\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	m := lifecycle.NewManager()
	m.Config().OutputDir = filepath.Join(dir, "output")
	m.SetPosts([]*models.Post{{Slug: "published", Href: "/published/", Published: true}})
	p := NewRedirectsPlugin()
	p.SetConfig(RedirectsConfig{RedirectsFile: file})
	if err := p.Write(m); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(m.Config().OutputDir, "published", "index.html")); !os.IsNotExist(err) {
		t.Fatalf("redirect claimed published route: %v", err)
	}
	native, err := os.ReadFile(filepath.Join(m.Config().OutputDir, "redirects.conf"))
	if err != nil || strings.Contains(string(native), "/published") {
		t.Fatalf("published route leaked to native config: %q, %v", native, err)
	}
	if _, err := os.Stat(filepath.Join(m.Config().OutputDir, "old", "index.html")); err != nil {
		t.Fatal(err)
	}
}

func TestRedirectsPluginRefreshesItsOwnFallback(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "_redirects")
	m := lifecycle.NewManager()
	m.Config().OutputDir = filepath.Join(dir, "output")
	p := NewRedirectsPlugin()
	p.SetConfig(RedirectsConfig{RedirectsFile: file})
	for _, dest := range []string{"/first", "/second"} {
		if err := os.WriteFile(file, []byte("/old "+dest+"\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := p.Write(m); err != nil {
			t.Fatal(err)
		}
		content, err := os.ReadFile(filepath.Join(m.Config().OutputDir, "old", "index.html"))
		if err != nil || !strings.Contains(string(content), "url='"+dest+"'") {
			t.Fatalf("redirect fallback did not update to %s: %q, %v", dest, content, err)
		}
	}
}
