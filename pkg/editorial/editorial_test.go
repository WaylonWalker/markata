package editorial

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func testGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	c := exec.Command("git", args...)
	c.Dir = dir
	c.Env = append(os.Environ(), "GIT_AUTHOR_NAME=Test", "GIT_AUTHOR_EMAIL=test@example.com", "GIT_COMMITTER_NAME=Test", "GIT_COMMITTER_EMAIL=test@example.com")
	if b, e := c.CombinedOutput(); e != nil {
		t.Fatalf("git %v: %v\n%s", args, e, b)
	}
	if len(args) > 0 && args[0] == "init" {
		testGit(t, dir, "config", "user.name", "Editorial Test")
		testGit(t, dir, "config", "user.email", "editorial@example.invalid")
	}
}
func mustMkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}
}
func mustWrite(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestSessionIsolatesChangesAndLinksHistory(t *testing.T) {
	d := t.TempDir()
	mustMkdir(t, filepath.Join(d, "pages"))
	mustWrite(t, filepath.Join(d, "pages", "post.md"), []byte("---\ntitle: Post\nauthors:\n  - id: writer\n---\nBody\n"))
	mustWrite(t, filepath.Join(d, "other.txt"), []byte("one\n"))
	testGit(t, d, "init")
	testGit(t, d, "remote", "add", "origin", "https://github.com/example/site.git")
	testGit(t, d, "add", ".")
	testGit(t, d, "commit", "-m", "initial")
	s, e := Start(Options{Root: d}, "agent", "copy-editor", []string{"pages/post.md"})
	if e != nil {
		t.Fatal(e)
	}
	mustWrite(t, filepath.Join(d, "pages", "post.md"), []byte("---\ntitle: Post\nauthors:\n  - id: writer\n---\nEdited body\n"))
	mustWrite(t, filepath.Join(d, "other.txt"), []byte("two\n"))
	done, e := Finish(Options{Root: d}, s.ID, "fixed prose")
	if e != nil {
		t.Fatal(e)
	}
	if len(done.Commits) != 2 {
		t.Fatalf("commits=%v", done.Commits)
	}
	for _, sha := range done.Commits {
		out := exec.Command("git", "show", "--format=", "--name-only", sha)
		out.Dir = d
		b, err := out.Output()
		if err != nil {
			t.Fatal(err)
		}
		if strings.TrimSpace(string(b)) != "pages/post.md" {
			t.Fatalf("commit %s files=%q", sha, b)
		}
	}
	b, err := os.ReadFile(filepath.Join(d, "pages", "post.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "editorial_history:") {
		t.Fatal("history missing")
	}
	b, err = os.ReadFile(filepath.Join(d, "other.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "two\n" {
		t.Fatalf("unrelated edit lost: %q", b)
	}
}

func TestStartRejectsSelectedStagedChanges(t *testing.T) {
	d := t.TempDir()
	mustMkdir(t, filepath.Join(d, "pages"))
	mustWrite(t, filepath.Join(d, "pages", "post.md"), []byte("---\ntitle: Post\n---\nBody\n"))
	testGit(t, d, "init")
	testGit(t, d, "add", ".")
	testGit(t, d, "commit", "-m", "initial")
	mustWrite(t, filepath.Join(d, "pages", "post.md"), []byte("changed"))
	testGit(t, d, "add", "pages/post.md")
	if _, e := Start(Options{Root: d}, "agent", "copy-editor", []string{"pages/post.md"}); e == nil || !strings.Contains(e.Error(), "staged") {
		t.Fatalf("error=%v", e)
	}
}

func TestUnchangedReviewAndMixedAuthorsRemainStable(t *testing.T) {
	got, err := credit("---\nauthors:\n  - writer\n  - id: prior\n    role: copy editor\n---\nBody\n", "agent", "copy-editor", "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "- writer") || !strings.Contains(got, "id: prior") || !strings.Contains(got, "id: agent") {
		t.Fatalf("authors lost during credit:\n%s", got)
	}
	d := t.TempDir()
	mustMkdir(t, filepath.Join(d, "posts"))
	mustWrite(t, filepath.Join(d, "posts", "review.md"), []byte("---\ntitle: Review\n---\nBody\n"))
	testGit(t, d, "init")
	testGit(t, d, "remote", "add", "origin", "https://github.com/example/site.git")
	testGit(t, d, "add", ".")
	testGit(t, d, "commit", "-m", "initial")
	first, err := Start(Options{Root: d, DefaultAuthor: "writer"}, "agent", "copy-editor", []string{"posts/review.md"})
	if err != nil {
		t.Fatal(err)
	}
	done, err := Finish(Options{Root: d, DefaultAuthor: "writer"}, first.ID, "Reviewed prose")
	if err != nil {
		t.Fatal(err)
	}
	if done.Actions["posts/review.md"] != "reviewed" {
		t.Fatalf("action=%q", done.Actions["posts/review.md"])
	}
	b, err := os.ReadFile(filepath.Join(d, "posts", "review.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "- writer") {
		t.Fatalf("default writer missing:\n%s", b)
	}
	second, err := Start(Options{Root: d, DefaultAuthor: "writer"}, "agent", "copy-editor", []string{"posts/review.md"})
	if err != nil {
		t.Fatal(err)
	}
	done2, err := Finish(Options{Root: d, DefaultAuthor: "writer"}, second.ID, "Reviewed again")
	if err != nil {
		t.Fatal(err)
	}
	if done2.Actions["posts/review.md"] != "reviewed" {
		t.Fatalf("second action=%q", done2.Actions["posts/review.md"])
	}
}

func TestCommitPreservesCurrentBodyAndLabelsEdit(t *testing.T) {
	d := t.TempDir()
	mustMkdir(t, filepath.Join(d, "posts"))
	path := filepath.Join(d, "posts", "owned.md")
	mustWrite(t, path, []byte("---\ntitle: Owned\nauthor: writer\n---\nOriginal\n"))
	testGit(t, d, "init")
	testGit(t, d, "remote", "add", "origin", "https://github.com/example/site.git")
	testGit(t, d, "add", ".")
	testGit(t, d, "commit", "-m", "initial")
	mustWrite(t, path, []byte("---\ntitle: Owned\nauthor: writer\n---\nCurrent edited body\n"))
	s, err := Commit(Options{Root: d}, "agent", "copy-editor", "Edited body", []string{"posts/owned.md"})
	if err != nil {
		t.Fatal(err)
	}
	if s.Actions["posts/owned.md"] != "edited" {
		t.Fatalf("action=%q", s.Actions["posts/owned.md"])
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "Current edited body") {
		t.Fatalf("current body replaced:\n%s", b)
	}
}
