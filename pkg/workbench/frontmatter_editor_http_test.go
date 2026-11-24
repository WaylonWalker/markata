package workbench

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func frontmatterRequestHTTP(t *testing.T, h *Handler, endpoint, origin string, payload any) *httptest.ResponseRecorder {
	t.Helper()
	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "http://localhost"+endpoint, bytes.NewReader(raw))
	r.RemoteAddr = "127.0.0.1:55551"
	r.Header.Set("Content-Type", "application/json")
	if origin != "" {
		r.Header.Set("Origin", origin)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func TestFrontmatterEditorHTTPInspectAndPatch(t *testing.T) {
	h := newHandler(t)
	source := "---\ntitle: 'Hello'\ncustom:\n  nested: [1, 2]\ndraft: true\n---\n# Body\n"
	inspect := frontmatterRequestHTTP(t, h, "/workspace/frontmatter/inspect", "http://localhost", map[string]any{"source": source})
	if inspect.Code != http.StatusOK || !strings.Contains(inspect.Body.String(), "\"name\":\"custom\"") {
		t.Fatalf("inspect: %d %s", inspect.Code, inspect.Body.String())
	}
	edit := frontmatterRequestHTTP(t, h, "/workspace/frontmatter/edit", "http://localhost", map[string]any{
		"source": source, "field": "published", "value": true,
	})
	if edit.Code != http.StatusOK {
		t.Fatalf("edit: %d %s", edit.Code, edit.Body.String())
	}
	var result struct {
		Source string `json:"source"`
	}
	if err := json.Unmarshal(edit.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result.Source, "published: true") || !strings.Contains(result.Source, "custom:\n  nested: [1, 2]") || !strings.HasSuffix(result.Source, "# Body\n") {
		t.Fatalf("bad patch: %q", result.Source)
	}
}

func TestFrontmatterEditorHTTPRejectsNonLocalWrites(t *testing.T) {
	h := newHandler(t)
	payload := map[string]any{"source": "---\ntitle: A\n---\n", "field": "title", "value": "B"}
	for _, origin := range []string{"", "https://elsewhere.test", "http://localhost.evil.test"} {
		w := frontmatterRequestHTTP(t, h, "/workspace/frontmatter/edit", origin, payload)
		if w.Code != http.StatusForbidden {
			t.Errorf("origin %q unexpectedly accepted: %d", origin, w.Code)
		}
	}
}
