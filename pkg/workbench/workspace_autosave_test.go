package workbench

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/editorworkspace"
)

func TestWorkspaceAutosaveRejectsStaleFilesystemRevisions(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "markata-go.toml"),
		[]byte("[markata-go.glob]\npatterns = ['*.md']\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	post := filepath.Join(root, "post.md")
	initial := "# First version\n"
	if err := os.WriteFile(post, []byte(initial), 0o600); err != nil {
		t.Fatal(err)
	}
	workspace, err := editorworkspace.New(root)
	if err != nil {
		t.Fatal(err)
	}
	handler := &Handler{workspace: workspace}
	const url = "http://127.0.0.1/workspace/file?path=post.md"
	const versionHeader = "X-Markata-Expected-Version"
	put := func(version, source string) *httptest.ResponseRecorder {
		t.Helper()
		req := httptest.NewRequestWithContext(context.Background(),
			http.MethodPut, url, strings.NewReader(source))
		req.RemoteAddr = "127.0.0.1:1234"
		req.Header.Set("Origin", "http://127.0.0.1")
		if version != "" {
			req.Header.Set(versionHeader, version)
		}
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, req)
		return response
	}
	read := func() string {
		t.Helper()
		body, readErr := os.ReadFile(post)
		if readErr != nil {
			t.Fatal(readErr)
		}
		return string(body)
	}
	if got := put("", "# Browser change\n"); got.Code != http.StatusPreconditionRequired {
		t.Fatalf("missing version returned %d: %s", got.Code, got.Body.String())
	}
	if read() != initial {
		t.Fatal("PUT without version changed the file")
	}

	// Simulate a second editor writing on disk after the browser loaded v0.
	external := "# External edit\n"
	if err := os.WriteFile(post, []byte(external), 0o600); err != nil {
		t.Fatal(err)
	}
	if got := put(editorworkspace.SourceVersion(initial), "# Stale browser change\n"); got.Code != http.StatusConflict {
		t.Fatalf("stale autosave returned %d: %s", got.Code, got.Body.String())
	}
	if read() != external {
		t.Fatal("stale browser autosave overwrote external work")
	}

	current := "# Merged browser change\n"
	if got := put(editorworkspace.SourceVersion(external), current); got.Code != http.StatusOK {
		t.Fatalf("fresh revision returned %d: %s", got.Code, got.Body.String())
	}
	if read() != current {
		t.Fatal("fresh autosave did not persist")
	}
	if got := put(editorworkspace.SourceVersion(external), "# Too late\n"); got.Code != http.StatusConflict {
		t.Fatalf("repeated old revision returned %d: %s", got.Code, got.Body.String())
	}
	if read() != current {
		t.Fatal("stale retry overwrote newer work")
	}
}
