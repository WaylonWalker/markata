package workbench

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/editorworkspace"
)

func publishPostRequest(t *testing.T, handler *Handler, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	data, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(data))
	request.Header.Set("Origin", "http://example.com")
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	return recorder
}

func initPublishRepository(t *testing.T) (string, *Handler) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	root := t.TempDir()
	gitForReviewTest(t, root, "init", "-q")
	gitForReviewTest(t, root, "config", "user.name", "Markata Review")
	gitForReviewTest(t, root, "config", "user.email", "markata-review@example.com")
	if err := os.MkdirAll(filepath.Join(root, "posts"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "posts", "one.md"), []byte("---\npublished: true\n---\n\n# Original\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "unrelated.md"), []byte("# Initial\n"), 0644); err != nil {
		t.Fatal(err)
	}
	gitForReviewTest(t, root, "add", ".")
	gitForReviewTest(t, root, "commit", "-qm", "baseline")
	return root, &Handler{root: root, workspace: publishWorkspaceFixture{root: root}}
}

func TestPublishCommitIsolatesSelectedFileAndPreservesRealIndex(t *testing.T) {
	root, handler := initPublishRepository(t)
	first, err := gitOutputForTest(root, "rev-parse", "HEAD")
	if err != nil { t.Fatal(err) }
	approved := "---\npublished: true\n---\n\n# Revised\n"
	if err := os.WriteFile(filepath.Join(root, "posts", "one.md"), []byte(approved), 0644); err != nil { t.Fatal(err) }
	if err := os.WriteFile(filepath.Join(root, "unrelated.md"), []byte("# Staged unrelated\n"), 0644); err != nil { t.Fatal(err) }
	gitForReviewTest(t, root, "add", "unrelated.md")
	recorder := publishPostRequest(t, handler, "/workspace/publish/commit", commitRequest{
		Path: "posts/one.md", Head: first, Version: editorworkspace.SourceVersion(approved),
	})
	if recorder.Code != 200 {
		t.Fatalf("status %d: %s", recorder.Code, recorder.Body.String())
	}
	var result commitResult
	if err := json.Unmarshal(recorder.Body.Bytes(), &result); err != nil { t.Fatal(err) }
	if result.Status != "committed-local" { t.Fatalf("unexpected status: %+v",result) }
	selected, err := gitOutputForTest(root, "show", "HEAD:posts/one.md")
	if err != nil || selected != approved { t.Fatalf("approved blob changed: %q %v", selected, err) }
	unrelated, err := gitOutputForTest(root, "show", "HEAD:unrelated.md")
	if err != nil || unrelated != "# Initial\n" {
		t.Fatalf("unrelated file included: %q %v", unrelated, err)
	}
	staged, err := gitOutputForTest(root, "diff", "--cached", "--name-only")
	if err != nil || staged != "unrelated.md" { t.Fatalf("original staged index not preserved: %q %v", staged, err) }
	status, err := gitOutputForTest(root, "status", "--porcelain", "--", "posts/one.md")
	if err != nil || status != "" { t.Fatalf("selected file left dirty: %q %v", status, err) }
}

func TestPublishCommitRejectsStaleAndPrivateContent(t *testing.T) {
	root, handler := initPublishRepository(t)
	head, err := gitOutputForTest(root, "rev-parse", "HEAD")
	if err != nil { t.Fatal(err) }
	content := "---\npublished: false\ndraft: true\n---\n\n# Still draft\n"
	if err := os.WriteFile(filepath.Join(root, "posts", "one.md"), []byte(content), 0644); err != nil { t.Fatal(err) }
	request := commitRequest{Path: "posts/one.md", Head: head, Version: editorworkspace.SourceVersion(content)}
	recorder := publishPostRequest(t, handler, "/workspace/publish/commit", request)
	if recorder.Code != 422 { t.Fatalf("expected 422 on private post: %d: %s", recorder.Code, recorder.Body.String()) }
	request.Version = strings.Repeat("a",64)
	recorder = publishPostRequest(t, handler, "/workspace/publish/commit", request)
	if recorder.Code != http.StatusConflict { t.Fatalf("expected stale source conflict: %d", recorder.Code) }
	request.Version = editorworkspace.SourceVersion(content)
	request.Head = strings.Repeat("b",40)
	// Fail fast on draft metadata before attempting the stale Git ref check.
	recorder = publishPostRequest(t, handler, "/workspace/publish/commit", request)
	if recorder.Code != 422 { t.Fatalf("unexpected status: %d", recorder.Code) }
}

func TestPublishRejectsStagedSelectedPath(t *testing.T) {
	root, handler := initPublishRepository(t)
	head, err := gitOutputForTest(root, "rev-parse", "HEAD")
	if err != nil { t.Fatal(err) }
	content := "---\npublished: true\n---\n\n# Updated\n"
	if err := os.WriteFile(filepath.Join(root, "posts", "one.md"), []byte(content), 0644); err != nil { t.Fatal(err) }
	gitForReviewTest(t, root, "add", "posts/one.md")
	recorder := publishPostRequest(t, handler, "/workspace/publish/commit", commitRequest{
		Path: "posts/one.md", Head: head, Version: editorworkspace.SourceVersion(content),
	})
	if recorder.Code != http.StatusConflict {
		t.Fatalf("expected staged selected path rejection: %d: %s", recorder.Code, recorder.Body.String())
	}
}

func TestPublishCommitMustHaveSameOrigin(t *testing.T) {
	_, handler := initPublishRepository(t)
	request := httptest.NewRequest(http.MethodPost, "/workspace/publish/commit", strings.NewReader("{}"))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("missing Origin must be rejected: %d", recorder.Code)
	}
}

func gitOutputForTest(root string, args ...string) (string, error) {
	out, err := exec.Command("git", append([]string{"-C", root}, args...)...).Output()
	return string(out), err
}
