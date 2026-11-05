package editorhub

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/editorsync"
)

type verifiedSubjectKey struct{}

const privateHTTPPath = "/api/v1/sites/family-blog/mutations"
const privateHTTPOrigin = "https://editor.example.test"
const privateMutationID = "8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafa"
const otherMutationID = "8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafb"

func newHTTPMutationFixture(t *testing.T) (http.Handler, *editorsync.Engine) {
	t.Helper()
	registry, err := NewRegistry([]Site{{
		ID: "family-blog", Name: "Family Blog", Grants: map[string]Grant{
			"alice": {Read: true, Write: true},
			"bob":   {Read: true},
		},
	}})
	if err != nil {
		t.Fatal(err)
	}
	db, err := editorsync.OpenSQLite(context.Background(), filepath.Join(t.TempDir(), "private-http.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("close SQLite: %v", err)
		}
	})
	engine, err := editorsync.NewEngine(db, func(_ context.Context, subject, site string) error {
		if !registry.Allows(subject, site, PermissionWrite) {
			return errors.New("not granted")
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	auth := func(r *http.Request) (string, error) {
		// The test's authenticated identity only exists in server context,
		// not in any client-provided X-Hlab-* or JSON request property.
		subject, ok := r.Context().Value(verifiedSubjectKey{}).(string)
		if !ok || subject == "" {
			return "", errors.New("no verified session")
		}
		return subject, nil
	}
	handler, err := NewPrivateMutationHandler(registry, auth, engine, privateHTTPOrigin)
	if err != nil {
		t.Fatal(err)
	}
	return handler, engine
}

func mutationRequest(subject, body string) *http.Request {
	request := httptest.NewRequest(http.MethodPost, privateHTTPOrigin+privateHTTPPath, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Origin", privateHTTPOrigin)
	request.Header.Set("X-Markata-Sync-Intent", "write")
	if subject != "" {
		request = request.WithContext(context.WithValue(request.Context(), verifiedSubjectKey{}, subject))
	}
	return request
}

func acceptedPrivateBody() string {
	return `{"document_id":"post-001","mutation_id":"` + privateMutationID +
		`","base_revision":0,"source":"# Hello from offline","deleted":false}`
}

func TestPrivateHTTPOnlyAcknowledgesCommittedAuthorizedMutations(t *testing.T) {
	handler, _ := newHTTPMutationFixture(t)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, mutationRequest("alice", acceptedPrivateBody()))
	if response.Code != http.StatusOK {
		t.Fatalf("expected durable receipt, got %d: %s", response.Code, response.Body.String())
	}
	var accepted map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &accepted); err != nil {
		t.Fatal(err)
	}
	if accepted["mutation_id"] != privateMutationID || accepted["revision"] != float64(1) ||
		accepted["digest"] == "" {
		t.Fatalf("wrong acceptance receipt: %+v", accepted)
	}
	if _, leaked := accepted["sequence"]; leaked {
		t.Fatal("internal per-site sequence leaked to browser")
	}
	if got := response.Header().Get("Cache-Control"); got != "private, no-store" {
		t.Fatalf("private response incorrectly cached: %q", got)
	}
	replay := httptest.NewRecorder()
	handler.ServeHTTP(replay, mutationRequest("alice", acceptedPrivateBody()))
	if replay.Code != http.StatusOK || replay.Body.String() != response.Body.String() {
		t.Fatalf("lost-response retry must be idempotent: %+v", replay)
	}
	changedPayload := strings.Replace(acceptedPrivateBody(), "# Hello from offline", "# Changed", 1)
	conflictingID := httptest.NewRecorder()
	handler.ServeHTTP(conflictingID, mutationRequest("alice", changedPayload))
	if conflictingID.Code != http.StatusConflict {
		t.Fatalf("UUID reuse accepted: %d", conflictingID.Code)
	}
	// A distinct stale UUID is persisted as a conflict before a 409 reply.
	stale := strings.Replace(acceptedPrivateBody(), privateMutationID, otherMutationID, 1)
	conflict := httptest.NewRecorder()
	handler.ServeHTTP(conflict, mutationRequest("alice", stale))
	if conflict.Code != http.StatusConflict {
		t.Fatalf("stale mutation accepted: %d: %s", conflict.Code, conflict.Body.String())
	}
	var body struct {
		Code            string `json:"code"`
		CurrentRevision int64  `json:"current_revision"`
		Preserved       bool   `json:"preserved"`
	}
	if err := json.Unmarshal(conflict.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Code != "revision_conflict" || body.CurrentRevision != 1 || !body.Preserved {
		t.Fatalf("conflict not retained before response: %+v", body)
	}
}

func TestPrivateHTTPRejectsForgedHeadersAndUnauthorizedWriters(t *testing.T) {
	handler, _ := newHTTPMutationFixture(t)
	for _, tc := range []struct {
		name  string
		req   *http.Request
		wants int
	}{
		{"no session with forged HLab headers", func() *http.Request {
			r := mutationRequest("", acceptedPrivateBody())
			r.Header.Set("X-Hlab-User-Id", "alice")
			r.Header.Set("X-Hlab-Auth-Method", "session")
			return r
		}(), http.StatusUnauthorized},
		{"read-only same-site user", mutationRequest("bob", acceptedPrivateBody()), http.StatusForbidden},
		{"unknown grant", mutationRequest("stranger", acceptedPrivateBody()), http.StatusForbidden},
		{"cross-site request", func() *http.Request {
			r := mutationRequest("alice", acceptedPrivateBody())
			r.Header.Set("Origin", "https://evil.example")
			return r
		}(), http.StatusForbidden},
		{"missing Origin", func() *http.Request {
			r := mutationRequest("alice", acceptedPrivateBody())
			r.Header.Del("Origin")
			return r
		}(), http.StatusForbidden},
		{"duplicated Origin", func() *http.Request {
			r := mutationRequest("alice", acceptedPrivateBody())
			r.Header.Add("Origin", "https://evil.example")
			return r
		}(), http.StatusForbidden},
		{"wrong sync intent", func() *http.Request {
			r := mutationRequest("alice", acceptedPrivateBody())
			r.Header.Set("X-Markata-Sync-Intent", "publish")
			return r
		}(), http.StatusForbidden},
		{"cross-site fetch metadata", func() *http.Request {
			r := mutationRequest("alice", acceptedPrivateBody())
			r.Header.Set("Sec-Fetch-Site", "cross-site")
			return r
		}(), http.StatusForbidden},
		{"spoofed body subject", mutationRequest("alice", strings.TrimSuffix(acceptedPrivateBody(), "}")+`,"subject_id":"bob"}`), http.StatusBadRequest},
		{"wrong content type", func() *http.Request {
			r := mutationRequest("alice", acceptedPrivateBody())
			r.Header.Set("Content-Type", "text/plain")
			return r
		}(), http.StatusBadRequest},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, tc.req)
			if w.Code != tc.wants {
				t.Fatalf("got status %d, want %d", w.Code, tc.wants)
			}
		})
	}
}

func TestPrivateHTTPRejectsUnknownRoutesAndOversizeBody(t *testing.T) {
	handler, _ := newHTTPMutationFixture(t)
	route := httptest.NewRecorder()
	req := mutationRequest("alice", acceptedPrivateBody())
	req.URL.Path = "/api/v1/sites/family-blog/builder-admin"
	handler.ServeHTTP(route, req)
	if route.Code != http.StatusNotFound {
		t.Fatalf("private writer exposed non-sync route: %d", route.Code)
	}
	tooLarge := strings.Repeat("a", 2<<20+4096)
	oversized := httptest.NewRecorder()
	handler.ServeHTTP(oversized, mutationRequest("alice", `{"document_id":"post","mutation_id":"`+
		privateMutationID+`","base_revision":0,"source":"`+tooLarge+`"}`))
	if oversized.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("large mutation accepted: %d", oversized.Code)
	}
}

func TestPrivateHTTPConstructorRequiresTrustedInputs(t *testing.T) {
	registry, err := NewRegistry([]Site{{ID: "site-a", Name: "Test"}})
	if err != nil {
		t.Fatal(err)
	}
	auth := func(*http.Request) (string, error) { return "alice", nil }
	_, engine := newHTTPMutationFixture(t)
	for _, origin := range []string{"", "http://editor.example.test", "https://editor.example.test/path", "https://user@editor.example.test"} {
		if handler, err := NewPrivateMutationHandler(registry, auth, engine, origin); err == nil || handler != nil {
			t.Fatalf("accepted invalid editor origin %q", origin)
		}
	}
	if handler, err := NewPrivateMutationHandler(registry, auth, nil, privateHTTPOrigin); err == nil || handler != nil {
		t.Fatal("missing SQLite writer accepted")
	}
	if handler, err := NewPrivateMutationHandler(registry, nil, engine, privateHTTPOrigin); err == nil || handler != nil {
		t.Fatal("missing trusted authenticator accepted")
	}
	if handler, err := NewPrivateMutationHandler(nil, auth, engine, privateHTTPOrigin); err == nil || handler != nil {
		t.Fatal("missing site registry accepted")
	}
}
