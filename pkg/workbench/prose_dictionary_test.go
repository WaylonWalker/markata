package workbench

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/prose"
)

func proseDictionaryRequest(origin, body string) *http.Request {
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost,
		"http://localhost/workspace/prose/dictionary/add", strings.NewReader(body))
	req.RemoteAddr = "127.0.0.1:42000"
	req.Header.Set("Origin", origin)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestProseEditorDictionaryAddAtWorkspaceRoot(t *testing.T) {
	root := t.TempDir()
	handler := &Handler{root: root}
	req := proseDictionaryRequest("http://localhost", `{"word":"QuickSync"}`)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	if response.Code != http.StatusOK {
		t.Fatalf("dictionary returned HTTP %d: %s", response.Code, response.Body.String())
	}
	prefs, err := prose.LoadPreferences(filepath.Join(root, "post.md"))
	if err != nil || !prefs.Options().AcceptedWords["quicksync"] {
		t.Fatalf("editor did not share a project dictionary: %+v, %v", prefs, err)
	}
}

func TestProseEditorDictionaryRejectsBadWrites(t *testing.T) {
	for _, tt := range []struct {
		name   string
		origin string
		body   string
	}{
		{"cross-origin", "https://example.org", `{"word":"QuickSync"}`},
		{"unknown field", "http://localhost", `{"word":"QuickSync","path":"../../outside"}`},
		{"whitespace", "http://localhost", `{"word":"two words"}`},
		{"file path", "http://localhost", `{"word":"../outside"}`},
		{"malformed", "http://localhost", `{"word":`},
		{"multiple", "http://localhost", `{"word":"QuickSync"}{"word":"other"}`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			handler := &Handler{root: root}
			req := proseDictionaryRequest(tt.origin, tt.body)
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, req)
			if response.Code < 400 {
				t.Fatalf("expected rejection, got %d: %s", response.Code, response.Body.String())
			}
			if _, err := os.Stat(filepath.Join(root, prose.PreferencesFile)); !os.IsNotExist(err) {
				t.Fatalf("rejected dictionary request changed workspace: %v", err)
			}
		})
	}
}

func TestProseEditorDictionaryRejectsSymlinkedConfig(t *testing.T) {
	root := t.TempDir()
	external := filepath.Join(t.TempDir(), "outside.json")
	if err := os.WriteFile(external, []byte("{}"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(external, filepath.Join(root, prose.PreferencesFile)); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	handler := &Handler{root: root}
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, proseDictionaryRequest("http://localhost", `{"word":"QuickSync"}`))
	if response.Code < 400 {
		t.Fatalf("symlinked config should be rejected: HTTP %d", response.Code)
	}
	content, err := os.ReadFile(external)
	if err != nil || string(content) != "{}" {
		t.Fatalf("external config changed: %q, %v", content, err)
	}
}
