package editorhub

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/editorsync"
)

func privateSyncFixture(t *testing.T) http.Handler {
	t.Helper()
	registry, err := NewRegistry([]Site{
		{ID: "family-blog", Name: "Family Blog", Grants: map[string]Grant{
			"alice": {Read: true, Write: true},
			"bob":   {Read: true, Write: true},
			"carol": {Read: true},
		}},
		{ID: "alice-only", Name: "Alice Only", Grants: map[string]Grant{
			"alice": {Read: true, Write: true},
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	db, err := editorsync.OpenSQLite(context.Background(), filepath.Join(t.TempDir(), "composed.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("close private sync SQLite: %v", err)
		}
	})
	writer, err := editorsync.NewEngine(db, func(_ context.Context, subject, site string) error {
		if !registry.Allows(subject, site, PermissionWrite) {
			return editorsync.ErrWriteForbidden
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	reader, err := editorsync.NewPrivateReader(db, func(_ context.Context, subject, site string) error {
		if !registry.Allows(subject, site, PermissionRead) {
			return editorsync.ErrReadForbidden
		}
		return nil
	}, []byte("0123456789abcdef0123456789abcdef"), time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	auth := func(r *http.Request) (string, error) {
		s, ok := r.Context().Value(verifiedSubjectKey{}).(string)
		if !ok || s == "" {
			return "", errors.New("missing verified session")
		}
		return s, nil
	}
	handler, err := NewPrivateSyncHandler(registry, auth, writer, reader, privateHTTPOrigin)
	if err != nil {
		t.Fatal(err)
	}
	return handler
}

func composedRequest(method, subject, route, body string) *http.Request {
	r := httptest.NewRequestWithContext(context.Background(), method, privateHTTPOrigin+route, strings.NewReader(body))
	if method == http.MethodPost {
		r.Header.Set("Origin", privateHTTPOrigin)
	}
	if strings.HasSuffix(route, "/mutations") {
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("X-Markata-Sync-Intent", "write")
	}
	if subject != "" {
		r = r.WithContext(context.WithValue(r.Context(), verifiedSubjectKey{}, subject))
	}
	return r
}

func composedCall(handler http.Handler, method, subject, route, body string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, composedRequest(method, subject, route, body))
	return w
}

func TestComposedPrivateSyncReadWriteReplayWithoutLeakingOtherUsers(t *testing.T) {
	handler := privateSyncFixture(t)
	sites := composedCall(handler, http.MethodGet, "alice", "/api/v1/sites", "")
	if sites.Code != http.StatusOK || !strings.Contains(sites.Body.String(), "alice-only") ||
		!strings.Contains(sites.Body.String(), "family-blog") {
		t.Fatalf("site discovery failed: %d %s", sites.Code, sites.Body.String())
	}
	bobsSites := composedCall(handler, http.MethodGet, "bob", "/api/v1/sites", "")
	if bobsSites.Code != http.StatusOK || strings.Contains(bobsSites.Body.String(), "alice-only") {
		t.Fatalf("site discovery leaked unrelated site: %d %s", bobsSites.Code, bobsSites.Body.String())
	}

	saved := composedCall(handler, http.MethodPost, "alice", privateHTTPPath, acceptedPrivateBody())
	if saved.Code != http.StatusOK {
		t.Fatalf("initial post failed: %d %s", saved.Code, saved.Body.String())
	}
	replay := composedCall(handler, http.MethodPost, "alice", privateHTTPPath, acceptedPrivateBody())
	if replay.Code != http.StatusOK || replay.Body.String() != saved.Body.String() {
		t.Fatal("lost HTTP acknowledgement did not replay identical durable receipt")
	}
	bobBody := `{"document_id":"bob-secret","mutation_id":"` + otherMutationID +
		`","base_revision":0,"source":"# Bob's private entry","deleted":false}`
	bobWrite := composedCall(handler, http.MethodPost, "bob", privateHTTPPath, bobBody)
	if bobWrite.Code != http.StatusOK {
		t.Fatalf("Bob could not write his own private draft: %d %s", bobWrite.Code, bobWrite.Body.String())
	}

	snapshot := composedCall(handler, http.MethodPost, "alice", "/api/v1/sites/family-blog/snapshots?limit=1", "")
	if snapshot.Code != http.StatusOK {
		t.Fatalf("snapshot failed: %d %s", snapshot.Code, snapshot.Body.String())
	}
	var initial privateHTTPSnapshotPage
	if err := json.Unmarshal(snapshot.Body.Bytes(), &initial); err != nil {
		t.Fatal(err)
	}
	if len(initial.Items) != 1 || initial.Items[0].DocumentID != "post-001" ||
		initial.Items[0].Revision != 1 || initial.ChangesCursor == "" ||
		strings.Contains(snapshot.Body.String(), "bob-secret") {
		t.Fatalf("private snapshot leaked or omitted source metadata: %+v", initial)
	}

	updateID := "8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafc"
	updated := `{"document_id":"post-001","mutation_id":"` + updateID +
		`","base_revision":1,"source":"# Alice second draft","deleted":false}`
	update := composedCall(handler, http.MethodPost, "alice", privateHTTPPath, updated)
	if update.Code != http.StatusOK {
		t.Fatalf("update failed: %d %s", update.Code, update.Body.String())
	}
	source := composedCall(handler, http.MethodGet, "alice", "/api/v1/sites/family-blog/documents/post-001?revision=1", "")
	if source.Code != http.StatusOK {
		t.Fatalf("hydrate old version failed: %d %s", source.Code, source.Body.String())
	}
	var hydrated struct {
		SourceBase64 string `json:"source_base64"`
	}
	if err := json.Unmarshal(source.Body.Bytes(), &hydrated); err != nil {
		t.Fatal(err)
	}
	raw, err := base64.StdEncoding.DecodeString(hydrated.SourceBase64)
	if err != nil || string(raw) != "# Hello from offline" {
		t.Fatalf("historical source altered after update: %q %v", raw, err)
	}
	feedRoute := "/api/v1/sites/family-blog/changes?cursor=" + url.QueryEscape(initial.ChangesCursor)
	changes := composedCall(handler, http.MethodGet, "alice", feedRoute, "")
	if changes.Code != http.StatusOK {
		t.Fatalf("change feed failed: %d %s", changes.Code, changes.Body.String())
	}
	var newer privateHTTPChangesPage
	if err := json.Unmarshal(changes.Body.Bytes(), &newer); err != nil {
		t.Fatal(err)
	}
	if len(newer.Items) != 1 || newer.Items[0].DocumentID != "post-001" ||
		newer.Items[0].Revision != 2 || strings.Contains(changes.Body.String(), "bob-secret") {
		t.Fatalf("private change feed included another user's data: %+v", newer)
	}
	bobCursor := composedCall(handler, http.MethodGet, "bob", feedRoute, "")
	if bobCursor.Code != http.StatusBadRequest {
		t.Fatalf("other subject reused opaque cursor: %d %s", bobCursor.Code, bobCursor.Body.String())
	}
	bobSource := composedCall(handler, http.MethodGet, "bob", "/api/v1/sites/family-blog/documents/post-001?revision=1", "")
	if bobSource.Code != http.StatusNotFound {
		t.Fatalf("other writer hydrated Alice's private source: %d", bobSource.Code)
	}

	// With the combined handler, rejected and orphaned source recovery must
	// remain author-only even when another user can write to the same site.
	conflictID := "8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafd"
	stale := `{"document_id":"post-001","mutation_id":"` + conflictID +
		`","base_revision":0,"source":"# Rejected offline update","deleted":false}`
	rejected := composedCall(handler, http.MethodPost, "alice", privateHTTPPath, stale)
	if rejected.Code != http.StatusConflict {
		t.Fatalf("stale write was not preserved: %d %s", rejected.Code, rejected.Body.String())
	}
	conflictURL := "/api/v1/sites/family-blog/conflicts/" + conflictID
	ownerRecovery := composedCall(handler, http.MethodGet, "alice", conflictURL, "")
	if ownerRecovery.Code != http.StatusOK {
		t.Fatalf("original author cannot recover conflict: %d %s", ownerRecovery.Code, ownerRecovery.Body.String())
	}
	var recovered struct {
		SourceBase64 string `json:"source_base64"`
	}
	if err := json.Unmarshal(ownerRecovery.Body.Bytes(), &recovered); err != nil {
		t.Fatal(err)
	}
	copyBytes, err := base64.StdEncoding.DecodeString(recovered.SourceBase64)
	if err != nil || string(copyBytes) != "# Rejected offline update" {
		t.Fatalf("durable conflict copy lost: %q %v", copyBytes, err)
	}
	if w := composedCall(handler, http.MethodGet, "bob", conflictURL, ""); w.Code != http.StatusNotFound {
		t.Fatalf("another site writer read private conflict: %d", w.Code)
	}
	if w := composedCall(handler, http.MethodGet, "carol", conflictURL, ""); w.Code != http.StatusForbidden {
		t.Fatalf("read-only site member recovered private conflict: %d", w.Code)
	}
}

func TestComposedPrivateSyncExcludesOperatorRoutesAndSpoofedHeaders(t *testing.T) {
	handler := privateSyncFixture(t)
	for _, route := range []string{
		"/api/v1/sites/family-blog/builder-admin",
		"/api/v1/sites/family-blog/git",
		"/api/v1/sites/family-blog/lsp",
	} {
		w := composedCall(handler, http.MethodGet, "alice", route, "")
		if w.Code != http.StatusNotFound {
			t.Fatalf("unexpected operator route %q: %d", route, w.Code)
		}
	}
	forged := composedRequest(http.MethodPost, "", privateHTTPPath, acceptedPrivateBody())
	forged.Header.Set("X-Hlab-User-Id", "alice")
	forged.Header.Set("X-Hlab-Auth-Method", "session")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, forged)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("forged browser headers authenticated a mutation: %d", w.Code)
	}
	unauthorized := composedCall(handler, http.MethodPost, "carol", privateHTTPPath, acceptedPrivateBody())
	if unauthorized.Code != http.StatusForbidden {
		t.Fatalf("site reader gained write access: %d", unauthorized.Code)
	}
	if w := composedCall(handler, http.MethodGet, "", "/api/v1/sites", ""); w.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous site discovery was allowed: %d", w.Code)
	}
}

func TestPrivateSyncCompositionRejectsMissingDependencies(t *testing.T) {
	registry, err := NewRegistry([]Site{{ID: "family-blog", Name: "Family Blog"}})
	if err != nil {
		t.Fatal(err)
	}
	auth := func(*http.Request) (string, error) { return "alice", nil }
	for _, origin := range []string{"", "http://editor.example.test", "https://editor.example.test/path"} {
		if handler, err := NewPrivateSyncHandler(registry, auth, nil, nil, origin); err == nil || handler != nil {
			t.Fatalf("accepted invalid composition with origin %q", origin)
		}
	}
}
