package workbench

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/editorworkspace"
)

func TestWorkspaceHTTP(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "markata-go.toml"), []byte("[markata-go.glob]\npatterns = ['*.md']\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "post.md"), []byte("---\ntitle: Workspace post\n---\n# Source"), 0o600); err != nil {
		t.Fatal(err)
	}
	workspace, err := editorworkspace.New(root)
	if err != nil {
		t.Fatal(err)
	}
	handler := &Handler{workspace: workspace}
	for _, test := range []struct {
		name, target, peer string
		status             int
		contains           string
	}{
		{"listing", "http://127.0.0.1/workspace", "127.0.0.1:1234", 200, "post.md"},
		{"metadata", "http://127.0.0.1/workspace?metadata=true", "127.0.0.1:1234", 200, "Workspace post"},
		{"source", "http://127.0.0.1/workspace/file?path=post.md", "127.0.0.1:1234", 200, "# Source"},
		{"traversal", "http://127.0.0.1/workspace/file?path=../post.md", "127.0.0.1:1234", 404, ""},
		{"host", "http://example.com/workspace", "127.0.0.1:1234", 403, ""},
		{"peer", "http://127.0.0.1/workspace", "192.0.2.1:1234", 403, ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.target, http.NoBody)
			request.RemoteAddr = test.peer
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != test.status || !strings.Contains(response.Body.String(), test.contains) {
				t.Fatalf("status %d body %s", response.Code, response.Body.String())
			}
		})
	}
}

func TestWorkspaceFileMutationRequiresSameOrigin(t *testing.T) {
	tests := []struct {
		name, method, path, origin, fetchSite string
		status                                int
	}{
		{"PUT cross-origin", http.MethodPut, "post.md", "http://attacker.test", "cross-site", http.StatusForbidden},
		{"POST cross-origin", http.MethodPost, "new-post.md", "http://attacker.test", "cross-site", http.StatusForbidden},
		{"PUT missing origin", http.MethodPut, "post.md", "", "", http.StatusForbidden},
		{"POST null origin", http.MethodPost, "new-post.md", "null", "cross-site", http.StatusForbidden},
		{"PUT wrong scheme", http.MethodPut, "post.md", "https://127.0.0.1", "same-origin", http.StatusForbidden},
		{"PUT cross-site metadata", http.MethodPut, "post.md", "http://127.0.0.1", "cross-site", http.StatusForbidden},
		{"PUT same-origin", http.MethodPut, "post.md", "http://127.0.0.1", "same-origin", http.StatusOK},
		{"POST same-origin", http.MethodPost, "new-post.md", "http://127.0.0.1", "same-origin", http.StatusCreated},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			if err := os.WriteFile(filepath.Join(root, "markata-go.toml"), []byte("[markata-go.glob]\npatterns = ['*.md']\n"), 0o600); err != nil {
				t.Fatal(err)
			}
			original := "# Original"
			if err := os.WriteFile(filepath.Join(root, "post.md"), []byte(original), 0o600); err != nil {
				t.Fatal(err)
			}
			workspace, err := editorworkspace.New(root)
			if err != nil {
				t.Fatal(err)
			}
			handler := &Handler{workspace: workspace}
			request := httptest.NewRequestWithContext(context.Background(), test.method, "http://127.0.0.1/workspace/file?path="+test.path, strings.NewReader("# Replaced"))
			request.RemoteAddr = "127.0.0.1:1234"
			if test.origin != "" {
				request.Header.Set("Origin", test.origin)
			}
			if test.fetchSite != "" {
				request.Header.Set("Sec-Fetch-Site", test.fetchSite)
			}
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != test.status {
				t.Fatalf("status %d body %s; want %d", response.Code, response.Body.String(), test.status)
			}

			writtenPath := filepath.Join(root, filepath.FromSlash(test.path))
			content, err := os.ReadFile(writtenPath)
			if test.method == http.MethodPost && test.status == http.StatusForbidden {
				if !os.IsNotExist(err) {
					t.Fatalf("cross-origin POST created a file: content=%q err=%v", content, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			want := original
			if test.status != http.StatusForbidden {
				want = "# Replaced"
			}
			if string(content) != want {
				t.Fatalf("file content = %q, want %q", content, want)
			}
		})
	}
}

func TestWorkspaceSiteAppearance(t *testing.T) {
	root := t.TempDir()
	config := "[markata-go.theme]\npalette = 'everforest'\npalette_light = 'everforest-light'\npalette_dark = 'everforest-dark'\nfallback_mode = 'light'\naesthetic = 'minimal'\ntext_size = 'small'\ncustom_css = 'editor.css'\n[markata-go.markdown.highlight]\nline_numbers = true\n"
	if err := os.WriteFile(filepath.Join(root, "markata-go.toml"), []byte(config), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "editor.css"), []byte(".post-content { letter-spacing: 1px; }"), 0o600); err != nil {
		t.Fatal(err)
	}
	handler, err := NewWithWorkspace(context.Background(), root)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = handler.Close() })
	for _, attr := range []string{`data-palette="everforest-light"`, `data-theme="light"`, `data-text-size="small"`} {
		if !strings.Contains(string(handler.index), attr) {
			t.Fatalf("editor missing %s", attr)
		}
	}
	request := httptest.NewRequest(http.MethodPost, "http://127.0.0.1/render", strings.NewReader("{\"markdown\":\"```go\\npackage main\\n```\",\"include_document\":true}"))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	var rendered renderResponse
	if err := json.Unmarshal(response.Body.Bytes(), &rendered); err != nil {
		t.Fatal(err)
	}
	if response.Code != http.StatusOK || !strings.Contains(rendered.HTML, `class="ln"`) {
		t.Fatalf("configured highlighting missing: %s", response.Body.String())
	}
	if rendered.Document == nil || !strings.Contains(*rendered.Document, `data-palette="everforest-light"`) {
		t.Fatal("preview missing site appearance")
	}
	if handler.siteConfig == nil || !handler.siteConfig.MarkdownConfig.Highlight.LineNumbers {
		t.Fatal("site renderer config not loaded")
	}
	if !strings.Contains(string(handler.assets["css/site-custom.css"]), "letter-spacing") {
		t.Fatal("site CSS missing")
	}
}

func TestWorkspaceWithoutConfig_DefaultPaletteIsAvailable(t *testing.T) {
	handler, err := NewWithWorkspace(context.Background(), t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = handler.Close() })
	selected := handler.appearance.Defaults["data-palette"]
	if selected != "default-dark" {
		t.Fatalf("default palette = %q", selected)
	}
	found := false
	for _, palette := range handler.choices["data-palette"] {
		if palette == selected {
			found = true
		}
	}
	if !found {
		t.Fatal("default palette absent from renderer choices")
	}
}

func TestWorkspace_ReusesRendererAndReloadsConfiguration(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "markata-go.toml")
	writeConfig := func(size string) {
		t.Helper()
		if err := os.WriteFile(configPath, []byte("[markata-go.theme]\ntext_size='"+size+"'\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	writeConfig("small")
	handler, err := NewWithWorkspace(context.Background(), root)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = handler.Close() })
	render := func() {
		t.Helper()
		request := httptest.NewRequest(http.MethodPost, "http://127.0.0.1/render", strings.NewReader(`{"markdown":"# Cached renderer"}`))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), "Cached renderer") {
			t.Fatalf("render: %d %s", response.Code, response.Body.String())
		}
	}
	original := handler.renderer
	render()
	render()
	if handler.renderer != original {
		t.Fatal("renderer rebuilt for unchanged configuration")
	}
	revision := handler.store.ConfigRevision()
	writeConfig("large")
	deadline := time.Now().Add(5 * time.Second)
	for handler.store.ConfigRevision() == revision {
		if time.Now().After(deadline) {
			t.Fatal("configuration invalidation did not arrive")
		}
		time.Sleep(10 * time.Millisecond)
	}
	render()
	if handler.renderer == original {
		t.Fatal("renderer retained stale configuration")
	}
	if handler.appearance.Defaults["data-text-size"] != "large" {
		t.Fatal("appearance retained stale defaults")
	}
	refreshed := handler.renderer
	render()
	if handler.renderer != refreshed {
		t.Fatal("refreshed renderer was not reused")
	}
}
