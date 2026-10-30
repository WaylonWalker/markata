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
