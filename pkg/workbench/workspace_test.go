package workbench

import (
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
		{"listing", "http://127.0.0.1/workspace", "127.0.0.1:1234", 200, "Workspace post"},
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
