package workbench

import (
	"bytes"
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

func publishSyncRequestForTest(t *testing.T, handler *Handler, head, source string) (*httptest.ResponseRecorder, publishSyncResult) {
	t.Helper()
	body, err := json.Marshal(publishSyncRequest{Path: "posts/one.md", Head: head, Version: editorworkspace.SourceVersion(source)})
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "http://localhost:8766/workspace/publish/sync", bytes.NewReader(body))
	r.RemoteAddr = "127.0.0.1:55555"
	r.Host = "localhost:8766"
	r.Header.Set("Origin", "http://localhost:8766")
	r.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, r)
	var result publishSyncResult
	if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
		t.Fatalf("decode sync result (%d): %s: %v", rec.Code, rec.Body.String(), err)
	}
	return rec, result
}

func setupPublishSyncRemote(t *testing.T, root string) (remoteWork, branch string) {
	t.Helper()
	bare := filepath.Join(t.TempDir(), "origin.git")
	remoteWork = filepath.Join(t.TempDir(), "remote-work")
	gitForReviewTest(t, root, "init", "--bare", bare)
	gitForReviewTest(t, root, "remote", "add", "origin", bare)
	branchOutput, err := gitOutputForTest(root, "symbolic-ref", "-q", "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	branch = strings.TrimSpace(branchOutput)
	shortBranch := strings.TrimPrefix(branch, "refs/heads/")
	gitForReviewTest(t, root, "push", "--set-upstream", "origin", shortBranch)
	gitForReviewTest(t, filepath.Dir(remoteWork), "clone", "--branch", shortBranch, bare, remoteWork)
	gitForReviewTest(t, remoteWork, "config", "user.name", "Remote Review")
	gitForReviewTest(t, remoteWork, "config", "user.email", "remote@example.com")
	return remoteWork, branch
}

func TestPublishSyncFastForwardsAndPreservesUnrelatedDraft(t *testing.T) {
	root, handler := initPublishRepository(t)
	remoteWork, branch := setupPublishSyncRemote(t, root)
	initial, err := gitOutputForTest(root, "rev-parse", "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	initial = strings.TrimSpace(initial)
	remoteText := "# Remote unrelated change\n"
	if err := os.WriteFile(filepath.Join(remoteWork, "unrelated.md"), []byte(remoteText), 0o600); err != nil {
		t.Fatal(err)
	}
	gitForReviewTest(t, remoteWork, "add", "unrelated.md")
	gitForReviewTest(t, remoteWork, "commit", "-qm", "remote unrelated change")
	gitForReviewTest(t, remoteWork, "push", "origin", strings.TrimPrefix(branch, "refs/heads/"))
	localText := "---\npublished: true\n---\n\n# Local post draft\n"
	if err := os.WriteFile(filepath.Join(root, "posts", "one.md"), []byte(localText), 0o600); err != nil {
		t.Fatal(err)
	}

	rec, result := publishSyncRequestForTest(t, handler, initial, localText)
	if rec.Code != http.StatusOK || result.Status != "synced" {
		t.Fatalf("sync status=%d result=%+v body=%s", rec.Code, result, rec.Body.String())
	}
	if result.Head != result.RemoteHead {
		t.Fatalf("synced head %s differs from origin %s", result.Head, result.RemoteHead)
	}
	if got, err := os.ReadFile(filepath.Join(root, "posts", "one.md")); err != nil || string(got) != localText {
		t.Fatalf("local draft changed: %q %v", got, err)
	}
	if got, err := gitOutputForTest(root, "show", "HEAD:unrelated.md"); err != nil || got != remoteText {
		t.Fatalf("origin change missing after fast-forward: %q %v", got, err)
	}
}

func TestPublishSyncShowsBothVersionsWhenOriginChangedPost(t *testing.T) {
	root, handler := initPublishRepository(t)
	remoteWork, branch := setupPublishSyncRemote(t, root)
	initial, err := gitOutputForTest(root, "rev-parse", "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	initial = strings.TrimSpace(initial)
	remoteText := "---\npublished: true\n---\n\n# Remote version\n"
	if err := os.WriteFile(filepath.Join(remoteWork, "posts", "one.md"), []byte(remoteText), 0o600); err != nil {
		t.Fatal(err)
	}
	gitForReviewTest(t, remoteWork, "add", "posts/one.md")
	gitForReviewTest(t, remoteWork, "commit", "-qm", "remote post edit")
	gitForReviewTest(t, remoteWork, "push", "origin", strings.TrimPrefix(branch, "refs/heads/"))
	localText := "---\npublished: true\n---\n\n# Local version\n"
	if err := os.WriteFile(filepath.Join(root, "posts", "one.md"), []byte(localText), 0o600); err != nil {
		t.Fatal(err)
	}

	rec, result := publishSyncRequestForTest(t, handler, initial, localText)
	if rec.Code != http.StatusConflict || result.Status != "review-required" {
		t.Fatalf("sync status=%d result=%+v body=%s", rec.Code, result, rec.Body.String())
	}
	if !result.RemoteExists || result.RemoteSource != remoteText {
		t.Fatalf("missing origin copy: %+v", result)
	}
	if got, err := os.ReadFile(filepath.Join(root, "posts", "one.md")); err != nil || string(got) != localText {
		t.Fatalf("local copy changed: %q %v", got, err)
	}
	if got, err := gitOutputForTest(root, "rev-parse", "HEAD"); err != nil || strings.TrimSpace(got) != initial {
		t.Fatalf("local HEAD changed during conflict: %q %v", got, err)
	}
}

func TestPublishSyncRebasesOnePageCommitAcrossUnrelatedOriginChange(t *testing.T) {
	root, handler := initPublishRepository(t)
	remoteWork, branch := setupPublishSyncRemote(t, root)
	localText := "---\npublished: true\n---\n\n# Reviewed local post\n"
	if err := os.WriteFile(filepath.Join(root, "posts", "one.md"), []byte(localText), 0o600); err != nil {
		t.Fatal(err)
	}
	gitForReviewTest(t, root, "add", "posts/one.md")
	gitForReviewTest(t, root, "commit", "-qm", "publish reviewed post")
	localHead, err := gitOutputForTest(root, "rev-parse", "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	localHead = strings.TrimSpace(localHead)
	remoteText := "# Remote unrelated change\n"
	if err := os.WriteFile(filepath.Join(remoteWork, "unrelated.md"), []byte(remoteText), 0o600); err != nil {
		t.Fatal(err)
	}
	gitForReviewTest(t, remoteWork, "add", "unrelated.md")
	gitForReviewTest(t, remoteWork, "commit", "-qm", "remote unrelated change")
	gitForReviewTest(t, remoteWork, "push", "origin", strings.TrimPrefix(branch, "refs/heads/"))

	rec, result := publishSyncRequestForTest(t, handler, localHead, localText)
	if rec.Code != http.StatusOK || result.Status != "synced" {
		t.Fatalf("sync status=%d result=%+v body=%s", rec.Code, result, rec.Body.String())
	}
	if result.Head == localHead || result.Head == result.RemoteHead {
		t.Fatalf("sync did not replay the reviewed commit: %+v", result)
	}
	if got, err := gitOutputForTest(root, "rev-parse", "HEAD^"); err != nil || strings.TrimSpace(got) != result.RemoteHead {
		t.Fatalf("rebased parent does not match origin: %q %v", got, err)
	}
	if got, err := gitOutputForTest(root, "show", "HEAD:posts/one.md"); err != nil || got != localText {
		t.Fatalf("reviewed post was lost: %q %v", got, err)
	}
	if got, err := gitOutputForTest(root, "show", "HEAD:unrelated.md"); err != nil || got != remoteText {
		t.Fatalf("origin change was lost: %q %v", got, err)
	}
	if got, err := os.ReadFile(filepath.Join(root, "posts", "one.md")); err != nil || string(got) != localText {
		t.Fatalf("saved source changed: %q %v", got, err)
	}
}

func TestPublishSyncRejectsOverlappingDivergedCommit(t *testing.T) {
	root, handler := initPublishRepository(t)
	remoteWork, branch := setupPublishSyncRemote(t, root)
	localText := "---\npublished: true\n---\n\n# Local version\n"
	if err := os.WriteFile(filepath.Join(root, "posts", "one.md"), []byte(localText), 0o600); err != nil {
		t.Fatal(err)
	}
	gitForReviewTest(t, root, "add", "posts/one.md")
	gitForReviewTest(t, root, "commit", "-qm", "local post edit")
	localHead, err := gitOutputForTest(root, "rev-parse", "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	localHead = strings.TrimSpace(localHead)
	remoteText := "---\npublished: true\n---\n\n# Remote version\n"
	if err := os.WriteFile(filepath.Join(remoteWork, "posts", "one.md"), []byte(remoteText), 0o600); err != nil {
		t.Fatal(err)
	}
	gitForReviewTest(t, remoteWork, "add", "posts/one.md")
	gitForReviewTest(t, remoteWork, "commit", "-qm", "remote post edit")
	gitForReviewTest(t, remoteWork, "push", "origin", strings.TrimPrefix(branch, "refs/heads/"))

	rec, result := publishSyncRequestForTest(t, handler, localHead, localText)
	if rec.Code != http.StatusConflict || result.Status != "review-required" || result.RemoteSource != remoteText {
		t.Fatalf("sync status=%d result=%+v body=%s", rec.Code, result, rec.Body.String())
	}
	if got, err := gitOutputForTest(root, "rev-parse", "HEAD"); err != nil || strings.TrimSpace(got) != localHead {
		t.Fatalf("local commit changed during conflict: %q %v", got, err)
	}
}

func TestPublishSyncRejectsStagedChanges(t *testing.T) {
	root, handler := initPublishRepository(t)
	setupPublishSyncRemote(t, root)
	head, err := gitOutputForTest(root, "rev-parse", "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	head = strings.TrimSpace(head)
	if err := os.WriteFile(filepath.Join(root, "unrelated.md"), []byte("# Staged\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	gitForReviewTest(t, root, "add", "unrelated.md")
	source, err := os.ReadFile(filepath.Join(root, "posts", "one.md"))
	if err != nil {
		t.Fatal(err)
	}
	rec, result := publishSyncRequestForTest(t, handler, head, string(source))
	if rec.Code != http.StatusConflict || result.Status != "review-required" || !strings.Contains(result.Message, "Staged changes") {
		t.Fatalf("staged changes were not rejected: HTTP %d %s", rec.Code, rec.Body.String())
	}
}

func TestPublishSyncRequiresSameOrigin(t *testing.T) {
	_, handler := initPublishRepository(t)
	r := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/workspace/publish/sync", strings.NewReader(`{"path":"posts/one.md"}`))
	r.RemoteAddr = "127.0.0.1:55555"
	r.Host = "localhost:8766"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, r)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("missing Origin was accepted: HTTP %d %s", rec.Code, rec.Body.String())
	}
}
