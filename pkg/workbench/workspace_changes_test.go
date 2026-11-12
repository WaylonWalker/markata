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
)

func TestParseWorkspaceChangesIncludesRenamesAndWhitespace(t *testing.T) {
	raw := " M posts/one.md\x00M  staged.md\x00?? new post.md\x00R  renamed.md\x00original.md\x00 D deleted.md\x00"
	entries, err := parseWorkspaceChanges(raw, ".")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 5 {
		t.Fatalf("expected five records, got %d: %+v", len(entries), entries)
	}
	byPath := map[string]workspaceChange{}
	for _, entry := range entries {
		byPath[entry.Path] = entry
		if !entry.Workspace {
			t.Errorf("root workspace missing path %q", entry.Path)
		}
	}
	if got := byPath["posts/one.md"]; got.Status != " M" || !got.Unstaged || got.Staged || !got.Markdown {
		t.Errorf("unstaged modified post: %+v", got)
	}
	if got := byPath["staged.md"]; !got.Staged || got.Unstaged {
		t.Errorf("staged file: %+v", got)
	}
	if got := byPath["new post.md"]; got.Kind != "untracked" || !got.Untracked {
		t.Errorf("untracked file with spaces: %+v", got)
	}
	if got := byPath["renamed.md"]; got.Kind != "renamed" || got.PreviousPath != "original.md" || !got.Staged {
		t.Errorf("rename: %+v", got)
	}
	if got := byPath["deleted.md"]; got.Kind != "deleted" || !got.Unstaged {
		t.Errorf("deletion: %+v", got)
	}
}

func TestParseWorkspaceChangesFlagsOutsideSitePaths(t *testing.T) {
	entries, err := parseWorkspaceChanges(" M blog/posts/inside.md\x00 M elsewhere/outside.md\x00", "..")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 || !entries[0].Workspace || entries[1].Workspace {
		t.Fatalf("wrong workspace scopes: %+v", entries)
	}
	for _, invalid := range []string{" M broken", "M\x00", "R  target.md\x00", "?? \x00"} {
		if _, err := parseWorkspaceChanges(invalid, "."); err == nil {
			t.Errorf("accepted malformed Git output %q", invalid)
		}
	}
}

func TestWorkspaceChangesReadOnlyListsAllGitStatusKinds(t *testing.T) {
	root, handler := initPublishRepository(t)
	post := filepath.Join(root, "posts", "one.md")
	if err := os.WriteFile(post, []byte("---\npublished: true\n---\n# Work in progress\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "posts", "new note.md"), []byte("# New\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "unrelated.md"), []byte("# staged elsewhere\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	gitForReviewTest(t, root, "add", "--", "unrelated.md")
	before, err := gitOutput(context.Background(), root, "status", "--porcelain=v1", "-z", "--untracked-files=all")
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "http://localhost:8766/workspace/changes", nil)
	request.RemoteAddr = "127.0.0.1:55555"
	reply := httptest.NewRecorder()
	handler.ServeHTTP(reply, request)
	if reply.Code != http.StatusOK {
		t.Fatalf("status %d: %s", reply.Code, reply.Body.String())
	}
	if reply.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("change inventory must not be cached")
	}
	var changes workspaceChanges
	if err := json.Unmarshal(reply.Body.Bytes(), &changes); err != nil {
		t.Fatal(err)
	}
	if changes.Total != 3 || changes.Staged != 1 || changes.Unstaged != 1 || changes.Untracked != 1 {
		t.Fatalf("wrong workspace change counts: %+v", changes)
	}
	seen := map[string]workspaceChange{}
	for _, item := range changes.Entries {
		seen[item.Path] = item
	}
	if !seen["posts/new note.md"].Untracked || !seen["posts/one.md"].Unstaged || !seen["unrelated.md"].Staged {
		t.Errorf("missing dirty files: %+v", changes.Entries)
	}
	after, err := gitOutput(context.Background(), root, "status", "--porcelain=v1", "-z", "--untracked-files=all")
	if err != nil || before != after {
		t.Fatalf("listing mutated Git status: before %q, after %q, err %v", before, after, err)
	}
}

func TestWorkspaceChangesRejectsNonGitAndRemoteAccess(t *testing.T) {
	root := t.TempDir()
	h := &Handler{root: root, workspace: publishWorkspaceFixture{root: root}}
	request := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "http://localhost:8766/workspace/changes", nil)
	request.RemoteAddr = "127.0.0.1:55555"
	reply := httptest.NewRecorder()
	h.ServeHTTP(reply, request)
	if reply.Code != http.StatusUnprocessableEntity {
		t.Fatalf("non-Git workspace returned %d: %s", reply.Code, reply.Body.String())
	}
	request = httptest.NewRequestWithContext(context.Background(), http.MethodGet, "http://localhost:8766/workspace/changes", nil)
	request.RemoteAddr = "203.0.113.9:55555"
	reply = httptest.NewRecorder()
	h.ServeHTTP(reply, request)
	if reply.Code != http.StatusForbidden {
		t.Fatalf("remote peer gained access: %d", reply.Code)
	}
	if strings.Contains(reply.Body.String(), root) {
		t.Fatal("remote response leaked workspace root")
	}
}
