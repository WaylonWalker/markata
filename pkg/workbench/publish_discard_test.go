package workbench

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/WaylonWalker/markata/pkg/editorworkspace"
)

type discardWorkspaceFixture struct{ publishWorkspaceFixture }

func (f discardWorkspaceFixture) Write(path, source string) error {
	return os.WriteFile(filepath.Join(f.root, filepath.FromSlash(path)), []byte(source), 0o600)
}

func TestPublishDiscardTrackedFileAndStaleChecks(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git unavailable")
	}
	dir := t.TempDir()
	gitForReviewTest(t, dir, "init", "-q")
	gitForReviewTest(t, dir, "config", "user.name", "Review Test")
	gitForReviewTest(t, dir, "config", "user.email", "review@example.com")
	if err := os.MkdirAll(filepath.Join(dir, "posts"), 0o700); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(dir, "posts", "example.md")
	baseline := "---\npublished: true\n---\n# Before\n"
	changed := "---\npublished: true\n---\n# After\n"
	if err := os.WriteFile(target, []byte(baseline), 0o600); err != nil {
		t.Fatal(err)
	}
	gitForReviewTest(t, dir, "add", "posts/example.md")
	gitForReviewTest(t, dir, "commit", "-qm", "initial")
	headCommand := exec.CommandContext(context.Background(), "git", "-C", dir, "rev-parse", "HEAD")
	headOutput, err := headCommand.Output()
	if err != nil {
		t.Fatal(err)
	}
	head := string(bytes.TrimSpace(headOutput))
	h := &Handler{root: dir, workspace: discardWorkspaceFixture{publishWorkspaceFixture{root: dir}}}
	send := func(version, expectedHead string, origin bool) *httptest.ResponseRecorder {
		t.Helper()
		data, err := json.Marshal(discardRequest{Path: "posts/example.md", Version: version, Head: expectedHead})
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "http://127.0.0.1:8766/workspace/publish/discard", bytes.NewReader(data))
		if origin {
			req.Header.Set("Origin", "http://127.0.0.1:8766")
		}
		req.Header.Set("Content-Type", "application/json")
		out := httptest.NewRecorder()
		h.servePublishDiscard(out, req)
		return out
	}
	if err := os.WriteFile(target, []byte(changed), 0o600); err != nil {
		t.Fatal(err)
	}
	version := editorworkspace.SourceVersion(changed)
	if got := send(version, head, false); got.Code != http.StatusForbidden {
		t.Fatalf("missing Origin got %d", got.Code)
	}
	if got := send(editorworkspace.SourceVersion(baseline), head, true); got.Code != http.StatusConflict {
		t.Fatalf("stale draft got %d: %s", got.Code, got.Body.String())
	}
	if got := send(version, "0000000000000000000000000000000000000000", true); got.Code != http.StatusConflict {
		t.Fatalf("stale HEAD got %d: %s", got.Code, got.Body.String())
	}
	gitForReviewTest(t, dir, "add", "posts/example.md")
	if got := send(version, head, true); got.Code != http.StatusConflict {
		t.Fatalf("staged path got %d: %s", got.Code, got.Body.String())
	}
	gitForReviewTest(t, dir, "reset", "-q", "HEAD", "--", "posts/example.md")
	if got := send(version, head, true); got.Code != http.StatusOK {
		t.Fatalf("discard got %d: %s", got.Code, got.Body.String())
	}
	actual, err := os.ReadFile(target)
	if err != nil || string(actual) != baseline {
		t.Fatalf("restore failed: %v %q", err, actual)
	}
	if got := send(version, head, true); got.Code != http.StatusConflict {
		t.Fatalf("retry should not overwrite restored draft: %d", got.Code)
	}
}

func TestPublishDiscardUntrackedFileIsPreserved(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git unavailable")
	}
	dir := t.TempDir()
	gitForReviewTest(t, dir, "init", "-q")
	gitForReviewTest(t, dir, "config", "user.name", "Review Test")
	gitForReviewTest(t, dir, "config", "user.email", "review@example.com")
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("base"), 0o600); err != nil {
		t.Fatal(err)
	}
	gitForReviewTest(t, dir, "add", "README.md")
	gitForReviewTest(t, dir, "commit", "-qm", "initial")
	headCommand := exec.CommandContext(context.Background(), "git", "-C", dir, "rev-parse", "HEAD")
	out, err := headCommand.Output()
	if err != nil {
		t.Fatal(err)
	}
	head := string(bytes.TrimSpace(out))
	content := "# New draft\n"
	path := filepath.Join(dir, "new.md")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	reqData, err := json.Marshal(discardRequest{Path: "new.md", Head: head, Version: editorworkspace.SourceVersion(content)})
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "http://127.0.0.1:8766/workspace/publish/discard", bytes.NewReader(reqData))
	req.Header.Set("Origin", "http://127.0.0.1:8766")
	h := &Handler{root: dir, workspace: discardWorkspaceFixture{publishWorkspaceFixture{root: dir}}}
	response := httptest.NewRecorder()
	h.servePublishDiscard(response, req)
	if response.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected safe rejection: %d %s", response.Code, response.Body.String())
	}
	actual, err := os.ReadFile(path)
	if err != nil || string(actual) != content {
		t.Fatalf("untracked draft was lost: %v %s", err, actual)
	}
}
