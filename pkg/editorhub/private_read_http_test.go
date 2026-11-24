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

type readVerifiedSubjectKey struct{}

const readOrigin = "https://editor.example.test"
const readSite = "family-blog"
const readMutationA = "8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafa"
const readMutationB = "8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafb"
const readMutationC = "8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafc"
const readMutationD = "8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafd"

func readHTTPFixture(t *testing.T) (http.Handler, *editorsync.Engine) {
	t.Helper()
	registry, err := NewRegistry([]Site{{
		ID: readSite, Name: "Family Blog", Grants: map[string]Grant{
			"alice": {Read: true, Write: true},
			"bob":   {Read: true, Write: true},
		},
	}})
	if err != nil {
		t.Fatal(err)
	}
	db, err := editorsync.OpenSQLite(context.Background(), filepath.Join(t.TempDir(), "read-http.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("close private read SQLite: %v", err)
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
		s, ok := r.Context().Value(readVerifiedSubjectKey{}).(string)
		if !ok || s == "" {
			return "", errors.New("not signed in")
		}
		return s, nil
	}
	handler, err := NewPrivateReadHandler(registry, auth, reader, readOrigin)
	if err != nil {
		t.Fatal(err)
	}
	return handler, writer
}

func readRequest(method, subject, route string) *http.Request {
	r := httptest.NewRequestWithContext(context.Background(), method, readOrigin+route, http.NoBody)
	if method == http.MethodPost {
		r.Header.Set("Origin", readOrigin)
	}
	if subject != "" {
		r = r.WithContext(context.WithValue(r.Context(), readVerifiedSubjectKey{}, subject))
	}
	return r
}

func saveForReadHTTP(t *testing.T, writer *editorsync.Engine, subject, document, id string, base uint64, source string) {
	t.Helper()
	m := editorsync.Mutation{
		SiteID: readSite, SubjectID: subject, DocumentID: document, MutationID: id,
		BaseRevision: base, Visibility: editorsync.Private, Source: []byte(source),
	}
	if _, err := writer.Apply(context.Background(), subject, m); err != nil {
		t.Fatal(err)
	}
}

func snapshotFromResponse(t *testing.T, w *httptest.ResponseRecorder) privateHTTPSnapshotPage {
	t.Helper()
	if w.Code != http.StatusOK {
		t.Fatalf("snapshot status %d: %s", w.Code, w.Body.String())
	}
	var page privateHTTPSnapshotPage
	if err := json.Unmarshal(w.Body.Bytes(), &page); err != nil {
		t.Fatal(err)
	}
	return page
}

func TestPrivateReadHTTPSnapshotIsStableAndLazy(t *testing.T) {
	handler, writer := readHTTPFixture(t)
	saveForReadHTTP(t, writer, "alice", "post-a", readMutationA, 0, "# A old")
	saveForReadHTTP(t, writer, "alice", "post-b", readMutationB, 0, "# B old")
	saveForReadHTTP(t, writer, "bob", "secret-bob", readMutationC, 0, "# Bob secret")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, readRequest(http.MethodPost, "alice", "/api/v1/sites/family-blog/snapshots?limit=1"))
	first := snapshotFromResponse(t, w)
	if len(first.Items) != 1 || first.Items[0].DocumentID != "post-a" ||
		first.NextCursor == "" || first.ChangesCursor != "" {
		t.Fatalf("first page leaked or lost snapshot: %+v", first)
	}
	if strings.Contains(w.Body.String(), "secret-bob") || strings.Contains(w.Body.String(), "Bob secret") {
		t.Fatal("snapshot exposed other private draft")
	}
	if !strings.Contains(w.Body.String(), "\"items\"") ||
		!strings.Contains(w.Body.String(), "\"document_id\"") ||
		strings.Contains(w.Body.String(), "\"Items\"") {
		t.Fatalf("snapshot JSON contract is not versionable snake_case: %s", w.Body.String())
	}
	if got := w.Header().Get("Cache-Control"); got != "private, no-store" {
		t.Fatalf("snapshot cache policy %q", got)
	}
	saveForReadHTTP(t, writer, "alice", "post-b", readMutationD, 1, "# B updated")
	nextURL := "/api/v1/sites/family-blog/documents?limit=1&cursor=" + url.QueryEscape(first.NextCursor)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, readRequest(http.MethodGet, "alice", nextURL))
	second := snapshotFromResponse(t, w)
	if len(second.Items) != 1 || second.Items[0].DocumentID != "post-b" ||
		second.Items[0].Revision != 1 || second.NextCursor != "" || second.ChangesCursor == "" {
		t.Fatalf("snapshot mixed later revisions: %+v", second)
	}
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, readRequest(http.MethodGet, "alice", "/api/v1/sites/family-blog/documents/post-b?revision=1"))
	if w.Code != http.StatusOK {
		t.Fatalf("source hydration %d: %s", w.Code, w.Body.String())
	}
	var source struct {
		Revision int64  `json:"revision"`
		Base64   string `json:"source_base64"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &source); err != nil {
		t.Fatal(err)
	}
	raw, err := base64.StdEncoding.DecodeString(source.Base64)
	if err != nil || source.Revision != 1 || string(raw) != "# B old" {
		t.Fatalf("snapshot source changed during pagination: %+v %v", source, err)
	}
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, readRequest(http.MethodGet, "alice", "/api/v1/sites/family-blog/changes?limit=10&cursor="+url.QueryEscape(second.ChangesCursor)))
	if w.Code != http.StatusOK {
		t.Fatalf("change feed status %d: %s", w.Code, w.Body.String())
	}
	var changed privateHTTPChangesPage
	if err := json.Unmarshal(w.Body.Bytes(), &changed); err != nil {
		t.Fatal(err)
	}
	if len(changed.Items) != 1 || changed.Items[0].DocumentID != "post-b" || changed.Items[0].Revision != 2 || changed.NextCursor == "" {
		t.Fatalf("incremental feed missed newer change: %+v", changed)
	}
}

func TestPrivateReadHTTPRejectsSpoofingAndCursorMisuse(t *testing.T) {
	handler, writer := readHTTPFixture(t)
	saveForReadHTTP(t, writer, "alice", "post-a", readMutationA, 0, "# Alice")
	saveForReadHTTP(t, writer, "alice", "post-b", readMutationB, 0, "# Alice two")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, readRequest(http.MethodPost, "alice", "/api/v1/sites/family-blog/snapshots?limit=1"))
	first := snapshotFromResponse(t, w)
	if first.NextCursor == "" {
		t.Fatal("expected paginated cursor")
	}
	cases := []struct {
		name    string
		request *http.Request
		status  int
	}{
		{"missing authentication", readRequest(http.MethodPost, "", "/api/v1/sites/family-blog/snapshots"), http.StatusUnauthorized},
		{"other subject cursor", readRequest(http.MethodGet, "bob", "/api/v1/sites/family-blog/documents?cursor="+url.QueryEscape(first.NextCursor)), http.StatusBadRequest},
		{"ungranted site", readRequest(http.MethodGet, "alice", "/api/v1/sites/unknown/changes?cursor="+url.QueryEscape(first.NextCursor)), http.StatusForbidden},
		{"cross-site fetch", func() *http.Request {
			r := readRequest(http.MethodGet, "alice", "/api/v1/sites/family-blog/documents?cursor="+url.QueryEscape(first.NextCursor))
			r.Header.Set("Sec-Fetch-Site", "cross-site")
			return r
		}(), http.StatusForbidden},
		{"missing POST origin", func() *http.Request {
			r := readRequest(http.MethodPost, "alice", "/api/v1/sites/family-blog/snapshots")
			r.Header.Del("Origin")
			return r
		}(), http.StatusForbidden},
		{"forged POST origin", func() *http.Request {
			r := readRequest(http.MethodPost, "alice", "/api/v1/sites/family-blog/snapshots")
			r.Header.Set("Origin", "https://evil.example")
			return r
		}(), http.StatusForbidden},
		{"duplicate POST origin", func() *http.Request {
			r := readRequest(http.MethodPost, "alice", "/api/v1/sites/family-blog/snapshots")
			r.Header.Add("Origin", "https://evil.example")
			return r
		}(), http.StatusForbidden},
		{"invalid cursor", readRequest(http.MethodGet, "alice", "/api/v1/sites/family-blog/changes?cursor=invalid"), http.StatusBadRequest},
		{"snapshot cursor in changes", readRequest(http.MethodGet, "alice", "/api/v1/sites/family-blog/changes?cursor="+url.QueryEscape(first.NextCursor)), http.StatusBadRequest},
		{"no change cursor", readRequest(http.MethodGet, "alice", "/api/v1/sites/family-blog/changes"), http.StatusBadRequest},
		{"duplicate limit", readRequest(http.MethodGet, "alice", "/api/v1/sites/family-blog/documents?cursor="+url.QueryEscape(first.NextCursor)+"&limit=1&limit=2"), http.StatusBadRequest},
		{"large limit", readRequest(http.MethodPost, "alice", "/api/v1/sites/family-blog/snapshots?limit=100000"), http.StatusBadRequest},
		{"wrong owner hydration", readRequest(http.MethodGet, "bob", "/api/v1/sites/family-blog/documents/post-a?revision=1"), http.StatusNotFound},
		{"unknown HTTP path", readRequest(http.MethodGet, "alice", "/api/v1/sites/family-blog/builder-admin"), http.StatusNotFound},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, tc.request)
			if response.Code != tc.status {
				t.Fatalf("got %d, want %d: %s", response.Code, tc.status, response.Body.String())
			}
		})
	}
	forged := readRequest(http.MethodPost, "", "/api/v1/sites/family-blog/snapshots")
	forged.Header.Set("X-Hlab-User-Id", "alice")
	forged.Header.Set("X-Hlab-Auth-Method", "session")
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, forged)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("untrusted identity headers authorized snapshot: %d", w.Code)
	}
}

func TestPrivateReadHTTPConstructorFailsClosed(t *testing.T) {
	registry, err := NewRegistry([]Site{{ID: "family-blog", Name: "Family"}})
	if err != nil {
		t.Fatal(err)
	}
	auth := func(*http.Request) (string, error) { return "alice", nil }
	for _, origin := range []string{"", "http://editor.example.test", "https://user@editor.example.test", "https://editor.example.test/path"} {
		if handler, err := NewPrivateReadHandler(registry, auth, nil, origin); err == nil || handler != nil {
			t.Fatalf("accepted invalid reader/origin %q", origin)
		}
	}
}
