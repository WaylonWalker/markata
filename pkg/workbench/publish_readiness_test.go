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

func publishReadinessRequest(t *testing.T, handler *Handler) *httptest.ResponseRecorder {
	t.Helper()
	r := httptest.NewRequestWithContext(context.Background(), http.MethodGet,
		"http://localhost:8766/workspace/publish/readiness?path=posts%2Fone.md", nil)
	r.RemoteAddr = "127.0.0.1:55555"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, r)
	return rec
}

func TestPublishReadinessRequiresRemoteBranchParity(t *testing.T) {
	root, handler := initPublishRepository(t)
	// A missing origin must never authorize publishing.
	if rec := publishReadinessRequest(t, handler); rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("missing origin: HTTP %d %s", rec.Code, rec.Body.String())
	}
	remote := filepath.Join(t.TempDir(), "origin.git")
	gitForReviewTest(t, root, "init", "--bare", remote)
	gitForReviewTest(t, root, "remote", "add", "origin", remote)
	branch, err := gitOutputForTest(root, "symbolic-ref", "-q", "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	branch = strings.TrimSpace(branch)
	if branch == "" {
		t.Fatal("missing test branch")
	}
	gitForReviewTest(t, root, "push", "origin", branch+":"+branch)
	rec := publishReadinessRequest(t, handler)
	if rec.Code != http.StatusOK {
		t.Fatalf("clean branch rejected: HTTP %d %s", rec.Code, rec.Body.String())
	}
	var ready publishReadiness
	if err := json.Unmarshal(rec.Body.Bytes(), &ready); err != nil {
		t.Fatal(err)
	}
	if !ready.Ready || ready.Branch != branch {
		t.Fatalf("wrong readiness response: %+v", ready)
	}
	// A hidden, unrelated local commit would also be pushed. Fail closed.
	if err := os.WriteFile(filepath.Join(root, "unrelated.md"), []byte("# local-only commit\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	gitForReviewTest(t, root, "add", "--", "unrelated.md")
	gitForReviewTest(t, root, "commit", "-qm", "unrelated local commit")
	rec = publishReadinessRequest(t, handler)
	if rec.Code != http.StatusConflict {
		t.Fatalf("unrelated commit authorized: HTTP %d %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "differs from origin") {
		t.Fatalf("expected actionable mismatch: %s", rec.Body.String())
	}
}
