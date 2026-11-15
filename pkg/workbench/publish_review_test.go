package workbench

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/editorworkspace"
	"github.com/WaylonWalker/markata/pkg/models"
)

type publishWorkspaceFixture struct{ root string }

func (f publishWorkspaceFixture) Config() (*models.Config, error) { return nil, nil }
func (f publishWorkspaceFixture) Paths() ([]editorworkspace.File, error) {
	return nil, nil
}
func (f publishWorkspaceFixture) Files() ([]editorworkspace.File, error) {
	return nil, nil
}
func (f publishWorkspaceFixture) Create(string, string) error { return nil }
func (f publishWorkspaceFixture) CreateWithImage(string, string, string, string, []byte) (string, error) {
	return "", nil
}
func (f publishWorkspaceFixture) Write(string, string) error { return nil }
func (f publishWorkspaceFixture) Read(path string) (string, error) {
	body, err := os.ReadFile(filepath.Join(f.root, filepath.FromSlash(path)))
	return string(body), err
}
func (publishWorkspaceFixture) ReadImage(string) ([]byte, error) { return nil, os.ErrNotExist }

func gitForReviewTest(t *testing.T, dir string, args ...string) {
	t.Helper()
	// #nosec G204 -- test helper invokes fixed Git subcommands without a shell.
	command := exec.CommandContext(context.Background(), "git", append([]string{"-C", dir}, args...)...)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %s: %v", args, output, err)
	}
}

func TestPublishReviewUsesExactCommittedSourceNotWorkingIndex(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	dir := t.TempDir()
	gitForReviewTest(t, dir, "init", "-q")
	gitForReviewTest(t, dir, "config", "user.name", "Review Test")
	gitForReviewTest(t, dir, "config", "user.email", "review@example.com")
	if err := os.MkdirAll(filepath.Join(dir, "posts"), 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "posts", "example.md")
	before := "---\npublished: false\n---\n\n# Original\n"
	after := "---\npublished: true\n---\n\n# Changed\n"
	if err := os.WriteFile(path, []byte(before), 0o600); err != nil {
		t.Fatal(err)
	}
	gitForReviewTest(t, dir, "add", "posts/example.md")
	gitForReviewTest(t, dir, "commit", "-qm", "baseline")
	if err := os.WriteFile(path, []byte(after), 0o600); err != nil {
		t.Fatal(err)
	}
	handler := &Handler{root: dir, workspace: publishWorkspaceFixture{root: dir}}
	request := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/workspace/publish/review?path=posts%2Fexample.md", http.NoBody)
	recorder := httptest.NewRecorder()
	handler.servePublishReview(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status %d: %s", recorder.Code, recorder.Body.String())
	}
	var got publishReview
	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if !got.BaselineExists || !got.Changed || got.Baseline != before || got.Current != after {
		t.Fatalf("wrong review snapshot: %+v", got)
	}
	if got.CanPublish || got.BaselineKind != "last-committed" {
		t.Fatalf("must not advertise a deployment: %+v", got)
	}
	if got.Head == "" || strings.Contains(got.Head, "\n") {
		t.Fatalf("bad HEAD: %q", got.Head)
	}
}

func TestPublishReviewRejectsUnsafePaths(t *testing.T) {
	dir := t.TempDir()
	handler := &Handler{root: dir, workspace: publishWorkspaceFixture{root: dir}}
	for _, path := range []string{"../escape.md", "/tmp/elsewhere.md", ""} {
		recorder := httptest.NewRecorder()
		handler.servePublishReview(recorder, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/workspace/publish/review?path="+path, http.NoBody))
		if recorder.Code != http.StatusBadRequest {
			t.Errorf("%q returned %d, expected 400", path, recorder.Code)
		}
	}
}

func TestPublishReviewUntrackedDocument(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	dir := t.TempDir()
	gitForReviewTest(t, dir, "init", "-q")
	gitForReviewTest(t, dir, "config", "user.name", "Review Test")
	gitForReviewTest(t, dir, "config", "user.email", "review@example.com")
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("base\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	gitForReviewTest(t, dir, "add", "README.md")
	gitForReviewTest(t, dir, "commit", "-qm", "baseline")
	if err := os.WriteFile(filepath.Join(dir, "new.md"), []byte("# New\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	handler := &Handler{root: dir, workspace: publishWorkspaceFixture{root: dir}}
	recorder := httptest.NewRecorder()
	handler.servePublishReview(recorder, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/workspace/publish/review?path=new.md", http.NoBody))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status %d: %s", recorder.Code, recorder.Body.String())
	}
	var got publishReview
	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.BaselineExists || !got.Changed || got.Baseline != "" || got.Current != "# New\n" {
		t.Fatalf("incorrect new-document review: %+v", got)
	}
}

func TestPublishReviewChangedNonMarkdownFile(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	dir := t.TempDir()
	gitForReviewTest(t, dir, "init", "-q")
	gitForReviewTest(t, dir, "config", "user.name", "Review Test")
	gitForReviewTest(t, dir, "config", "user.email", "review@example.com")
	if err := os.MkdirAll(filepath.Join(dir, "config"), 0o755); err != nil {
		t.Fatal(err)
	}
	before := "[feeds]\nname = \"old\"\n"
	after := "[feeds]\nname = \"new\"\n"
	feeds := filepath.Join(dir, "config", "feeds.toml")
	unchanged := filepath.Join(dir, "config", "authors.toml")
	binary := filepath.Join(dir, "config", "blob.bin")
	for path, body := range map[string]string{feeds: before, unchanged: "x = 1\n", binary: "a\x00b"} {
		if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	gitForReviewTest(t, dir, "add", ".")
	gitForReviewTest(t, dir, "commit", "-qm", "baseline")
	if err := os.WriteFile(feeds, []byte(after), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(binary, []byte("c\x00d"), 0o600); err != nil {
		t.Fatal(err)
	}
	handler := &Handler{root: dir, workspace: publishWorkspaceFixture{root: dir}}
	serve := func(path string) *httptest.ResponseRecorder {
		recorder := httptest.NewRecorder()
		handler.servePublishReview(recorder, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/workspace/publish/review?path="+path, http.NoBody))
		return recorder
	}

	recorder := serve("config%2Ffeeds.toml")
	if recorder.Code != http.StatusOK {
		t.Fatalf("changed toml status %d: %s", recorder.Code, recorder.Body.String())
	}
	var got publishReview
	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Baseline != before || got.Current != after || !got.Changed {
		t.Fatalf("wrong toml review snapshot: %+v", got)
	}
	if recorder := serve("config%2Fauthors.toml"); recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("unchanged file returned %d, expected 422", recorder.Code)
	}
	if recorder := serve("config%2Fblob.bin"); recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("binary file returned %d, expected 422", recorder.Code)
	}
}
