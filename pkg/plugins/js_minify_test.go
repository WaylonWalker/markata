package plugins

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestJSMinifyPlugin_Name(t *testing.T) {
	p := NewJSMinifyPlugin()
	if got := p.Name(); got != "js_minify" {
		t.Errorf("Name() = %q, want %q", got, "js_minify")
	}
}

func TestJSMinifyPlugin_Configure(t *testing.T) {
	tests := []struct {
		name           string
		extra          map[string]interface{}
		wantEnabled    bool
		wantExcludeLen int
	}{
		{
			name:           "default config",
			extra:          nil,
			wantEnabled:    true,
			wantExcludeLen: 0,
		},
		{
			name: "disabled",
			extra: map[string]interface{}{
				"js_minify": map[string]interface{}{
					"enabled": false,
				},
			},
			wantEnabled:    false,
			wantExcludeLen: 0,
		},
		{
			name: "with exclude patterns",
			extra: map[string]interface{}{
				"js_minify": map[string]interface{}{
					"enabled": true,
					"exclude": []interface{}{"pagefind-ui.js", "vendor-*.js"},
				},
			},
			wantEnabled:    true,
			wantExcludeLen: 2,
		},
		{
			name: "with typed config",
			extra: map[string]interface{}{
				"js_minify": models.JSMinifyConfig{
					Enabled: true,
					Exclude: []string{"test.js"},
				},
			},
			wantEnabled:    true,
			wantExcludeLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewJSMinifyPlugin()
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
		})
	}
}

func TestJSMinifyPlugin_Write(t *testing.T) {
	tmpDir := t.TempDir()
	jsDir := filepath.Join(tmpDir, "js")
	if err := os.MkdirAll(jsDir, 0o755); err != nil {
		t.Fatalf("failed to create js dir: %v", err)
	}

	// Create a test JS file with whitespace and comments
	testJS := `// This is a comment
function greet(name) {
    // Say hello
    var message = "Hello, " + name + "!";
    console.log(message);
    return message;
}

/* Another comment block */
function add(a, b) {
    return a + b;
}

// Call the function
greet("World");
`
	jsPath := filepath.Join(jsDir, "test.js")
	//nolint:gosec // G306: test file permissions
	if err := os.WriteFile(jsPath, []byte(testJS), 0o644); err != nil {
		t.Fatalf("failed to write test JS: %v", err)
	}

	p := NewJSMinifyPlugin()
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		ContentDir: t.TempDir(),
		OutputDir:  tmpDir,
		Extra: map[string]interface{}{
			"js_minify": map[string]interface{}{
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

	// Read the minified file
	content, err := os.ReadFile(jsPath)
	if err != nil {
		t.Fatalf("failed to read minified JS: %v", err)
	}

	minifiedJS := string(content)

	// Verify minification occurred - comments should be removed
	if strings.Contains(minifiedJS, "// This is a comment") {
		t.Error("minified JS should not contain line comments")
	}

	if strings.Contains(minifiedJS, "/* Another comment block */") {
		t.Error("minified JS should not contain block comments")
	}

	// Verify essential JS is preserved
	if !strings.Contains(minifiedJS, "greet") {
		t.Error("minified JS should contain 'greet' function name")
	}

	if !strings.Contains(minifiedJS, "console.log") {
		t.Error("minified JS should contain 'console.log' call")
	}

	// Verify size reduction
	if len(minifiedJS) >= len(testJS) {
		t.Errorf("minified JS (%d bytes) should be smaller than original (%d bytes)",
			len(minifiedJS), len(testJS))
	}
}

func TestJSMinifyPlugin_Write_Disabled(t *testing.T) {
	tmpDir := t.TempDir()
	jsDir := filepath.Join(tmpDir, "js")
	if err := os.MkdirAll(jsDir, 0o755); err != nil {
		t.Fatalf("failed to create js dir: %v", err)
	}

	testJS := `function hello() { return "world"; }`
	jsPath := filepath.Join(jsDir, "test.js")
	//nolint:gosec // G306: test file permissions
	if err := os.WriteFile(jsPath, []byte(testJS), 0o644); err != nil {
		t.Fatalf("failed to write test JS: %v", err)
	}

	p := NewJSMinifyPlugin()
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		ContentDir: t.TempDir(),
		OutputDir:  tmpDir,
		Extra: map[string]interface{}{
			"js_minify": map[string]interface{}{
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

	// File should be unchanged
	content, err := os.ReadFile(jsPath)
	if err != nil {
		t.Fatalf("failed to read JS: %v", err)
	}

	if string(content) != testJS {
		t.Error("JS should be unchanged when plugin is disabled")
	}
}

func TestJSMinifyPlugin_IsExcluded_PagefindAssets(t *testing.T) {
	p := NewJSMinifyPlugin()
	if !p.isExcluded("/tmp/output/_pagefind/pagefind.js") {
		t.Fatal("expected _pagefind/pagefind.js to be excluded")
	}
	if p.isExcluded("/tmp/output/js/app.js") {
		t.Fatal("did not expect normal JS asset to be excluded")
	}
}

func TestJSMinifyPlugin_Write_Exclude(t *testing.T) {
	tmpDir := t.TempDir()
	jsDir := filepath.Join(tmpDir, "js")
	if err := os.MkdirAll(jsDir, 0o755); err != nil {
		t.Fatalf("failed to create js dir: %v", err)
	}

	testJS := `// Comment to remove
function test() {
    return true;
}
`
	regularPath := filepath.Join(jsDir, "app.js")
	excludedPath := filepath.Join(jsDir, "vendor-lib.js")

	//nolint:gosec // G306: test file permissions
	if err := os.WriteFile(regularPath, []byte(testJS), 0o644); err != nil {
		t.Fatalf("failed to write regular JS: %v", err)
	}
	//nolint:gosec // G306: test file permissions
	if err := os.WriteFile(excludedPath, []byte(testJS), 0o644); err != nil {
		t.Fatalf("failed to write excluded JS: %v", err)
	}

	p := NewJSMinifyPlugin()
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		ContentDir: t.TempDir(),
		OutputDir:  tmpDir,
		Extra: map[string]interface{}{
			"js_minify": map[string]interface{}{
				"enabled": true,
				"exclude": []interface{}{"vendor-*.js"},
			},
		},
	})

	if err := p.Configure(m); err != nil {
		t.Fatalf("Configure error: %v", err)
	}

	if err := p.Write(m); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	// Regular file should be minified
	regularContent, err := os.ReadFile(regularPath)
	if err != nil {
		t.Fatalf("failed to read regular JS: %v", err)
	}
	if len(regularContent) >= len(testJS) {
		t.Error("app.js should be minified")
	}

	// Excluded file should be unchanged
	excludedContent, err := os.ReadFile(excludedPath)
	if err != nil {
		t.Fatalf("failed to read excluded JS: %v", err)
	}
	if string(excludedContent) != testJS {
		t.Error("vendor-lib.js should not be minified (excluded by glob pattern)")
	}
}

func TestJSMinifyPlugin_Write_SkipsMinJS(t *testing.T) {
	tmpDir := t.TempDir()
	jsDir := filepath.Join(tmpDir, "js")
	if err := os.MkdirAll(jsDir, 0o755); err != nil {
		t.Fatalf("failed to create js dir: %v", err)
	}

	// Already-minified file should be completely skipped (not even found)
	minifiedJS := `function a(){return 1}`
	minPath := filepath.Join(jsDir, "lib.min.js")
	//nolint:gosec // G306: test file permissions
	if err := os.WriteFile(minPath, []byte(minifiedJS), 0o644); err != nil {
		t.Fatalf("failed to write min.js: %v", err)
	}

	p := NewJSMinifyPlugin()
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		ContentDir: t.TempDir(),
		OutputDir:  tmpDir,
		Extra: map[string]interface{}{
			"js_minify": map[string]interface{}{
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

	// .min.js file should be unchanged
	content, err := os.ReadFile(minPath)
	if err != nil {
		t.Fatalf("failed to read min.js: %v", err)
	}
	if string(content) != minifiedJS {
		t.Error("lib.min.js should not be modified")
	}
}

func TestJSMinifyPlugin_Priority(t *testing.T) {
	p := NewJSMinifyPlugin()

	if got := p.Priority(lifecycle.StageWrite); got != lifecycle.PriorityLast {
		t.Errorf("Priority(StageWrite) = %d, want %d", got, lifecycle.PriorityLast)
	}

	if got := p.Priority(lifecycle.StageRender); got != lifecycle.PriorityDefault {
		t.Errorf("Priority(StageRender) = %d, want %d", got, lifecycle.PriorityDefault)
	}
}

func TestJSMinifyPlugin_SizeReduction(t *testing.T) {
	tmpDir := t.TempDir()
	jsDir := filepath.Join(tmpDir, "js")
	if err := os.MkdirAll(jsDir, 0o755); err != nil {
		t.Fatalf("failed to create js dir: %v", err)
	}

	// Create a larger JS file to test meaningful reduction
	largeJS := `
// =========================================
// View Transitions Handler
// =========================================

/**
 * Handles view transitions for SPA-like navigation.
 * This manages the transition between pages without full reloads.
 */

(function() {
    "use strict";

    // Configuration
    var CONFIG = {
        transitionDuration: 300,
        easing: "ease-in-out",
        historyScrollRestoration: "manual"
    };

    /**
     * Initialize the view transition handler.
     * Sets up event listeners and configures the browser history API.
     */
    function initialize() {
        // Set up scroll restoration
        if ("scrollRestoration" in history) {
            history.scrollRestoration = CONFIG.historyScrollRestoration;
        }

        // Bind click handler to all internal links
        document.addEventListener("click", function(event) {
            var link = event.target.closest("a");
            if (!link) {
                return;
            }

            // Check if this is an internal link
            if (link.hostname !== window.location.hostname) {
                return;
            }

            // Prevent default navigation
            event.preventDefault();

            // Perform the transition
            navigateTo(link.href);
        });

        // Handle browser back/forward
        window.addEventListener("popstate", function(event) {
            if (event.state && event.state.url) {
                navigateTo(event.state.url, true);
            }
        });
    }

    /**
     * Navigate to a new URL with a view transition.
     * @param {string} url - The URL to navigate to
     * @param {boolean} isPopState - Whether this is from browser navigation
     */
    function navigateTo(url, isPopState) {
        // Fetch the new page
        fetch(url)
            .then(function(response) {
                if (!response.ok) {
                    throw new Error("HTTP error: " + response.status);
                }
                return response.text();
            })
            .then(function(html) {
                // Parse the HTML
                var parser = new DOMParser();
                var doc = parser.parseFromString(html, "text/html");

                // Update the page content
                updatePage(doc, url, isPopState);
            })
            .catch(function(error) {
                console.error("Navigation failed:", error);
                // Fallback to traditional navigation
                window.location.href = url;
            });
    }

    /**
     * Update the page content with the new document.
     * @param {Document} doc - The parsed HTML document
     * @param {string} url - The new URL
     * @param {boolean} isPopState - Whether from browser navigation
     */
    function updatePage(doc, url, isPopState) {
        // Swap the body content
        document.body.innerHTML = doc.body.innerHTML;

        // Update the title
        document.title = doc.title;

        // Update browser history
        if (!isPopState) {
            history.pushState({ url: url }, doc.title, url);
        }

        // Scroll to top
        window.scrollTo(0, 0);

        // Dispatch completion event
        document.dispatchEvent(new CustomEvent("view-transition-complete", {
            detail: { url: url }
        }));
    }

    // Initialize when DOM is ready
    if (document.readyState === "loading") {
        document.addEventListener("DOMContentLoaded", initialize);
    } else {
        initialize();
    }
})();
`
	jsPath := filepath.Join(jsDir, "view-transitions.js")
	//nolint:gosec // G306: test file permissions
	if err := os.WriteFile(jsPath, []byte(largeJS), 0o644); err != nil {
		t.Fatalf("failed to write test JS: %v", err)
	}

	originalSize := len(largeJS)

	p := NewJSMinifyPlugin()
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		ContentDir: t.TempDir(),
		OutputDir:  tmpDir,
		Extra: map[string]interface{}{
			"js_minify": map[string]interface{}{
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

	content, err := os.ReadFile(jsPath)
	if err != nil {
		t.Fatalf("failed to read minified JS: %v", err)
	}

	minifiedSize := len(content)
	reduction := float64(originalSize-minifiedSize) / float64(originalSize) * 100

	// Verify meaningful reduction (should be at least 25% for a file with comments and whitespace)
	if reduction < 25 {
		t.Errorf("JS reduction = %.1f%%, want at least 25%%", reduction)
	}

	t.Logf("JS minification: %d -> %d bytes (%.1f%% reduction)", originalSize, minifiedSize, reduction)
}
