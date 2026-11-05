package editorhub

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPrivateSessionBootstrapIsFreshAndSiteScoped(t *testing.T) {
	handler := privateSyncFixture(t)
	alice := composedCall(handler, http.MethodGet, "alice", "/api/v1/session", "")
	if alice.Code != http.StatusOK {
		t.Fatalf("Alice session unavailable: %d %s", alice.Code, alice.Body.String())
	}
	var first struct {
		SubjectID string        `json:"subject_id"`
		Sites     []SiteSummary `json:"sites"`
	}
	if err := json.Unmarshal(alice.Body.Bytes(), &first); err != nil {
		t.Fatal(err)
	}
	if first.SubjectID != "alice" || len(first.Sites) != 2 {
		t.Fatalf("unexpected private session bootstrap: %+v", first)
	}
	for _, site := range first.Sites {
		if !site.Permissions.Read || !site.Permissions.Write || site.Permissions.Publish {
			t.Fatalf("incorrect Alice capability leak: %+v", site)
		}
	}
	if got := alice.Header().Get("Cache-Control"); got != "private, no-store" {
		t.Fatalf("session response was cached: %q", got)
	}
	if alice.Header().Get("Cross-Origin-Resource-Policy") != "same-origin" {
		t.Fatal("session response lacks same-origin resource policy")
	}
	if strings.Contains(alice.Body.String(), "X-Hlab") || strings.Contains(alice.Body.String(), "token") {
		t.Fatal("proxy assertion or secret field leaked to browser")
	}

	// Each call reauthenticates and cannot reuse the previous site's identity.
	bob := composedCall(handler, http.MethodGet, "bob", "/api/v1/session", "")
	var second struct {
		SubjectID string        `json:"subject_id"`
		Sites     []SiteSummary `json:"sites"`
	}
	if err := json.Unmarshal(bob.Body.Bytes(), &second); err != nil {
		t.Fatal(err)
	}
	if bob.Code != http.StatusOK || second.SubjectID != "bob" ||
		len(second.Sites) != 1 || second.Sites[0].ID != "family-blog" {
		t.Fatalf("account switch leaked previous subject/site: %d %+v", bob.Code, second)
	}
	carol := composedCall(handler, http.MethodGet, "carol", "/api/v1/session", "")
	var readOnly struct {
		Sites []SiteSummary `json:"sites"`
	}
	if err := json.Unmarshal(carol.Body.Bytes(), &readOnly); err != nil {
		t.Fatal(err)
	}
	if carol.Code != http.StatusOK || len(readOnly.Sites) != 1 ||
		!readOnly.Sites[0].Permissions.Read || readOnly.Sites[0].Permissions.Write {
		t.Fatalf("read-only user received write grant: %d %+v", carol.Code, readOnly)
	}
	stranger := composedCall(handler, http.MethodGet, "stranger", "/api/v1/session", "")
	var hidden struct {
		Sites []SiteSummary `json:"sites"`
	}
	if err := json.Unmarshal(stranger.Body.Bytes(), &hidden); err != nil {
		t.Fatal(err)
	}
	if stranger.Code != http.StatusOK || len(hidden.Sites) != 0 {
		t.Fatalf("unknown subject received site access: %d %+v", stranger.Code, hidden)
	}
}

func TestPrivateSessionRejectsSpoofingAndCrossSiteRequests(t *testing.T) {
	handler := privateSyncFixture(t)
	for _, tc := range []struct {
		name    string
		request *http.Request
		status  int
	}{
		{"missing session", composedRequest(http.MethodGet, "", "/api/v1/session", ""), http.StatusUnauthorized},
		{"invalid subject", composedRequest(http.MethodGet, "alice\n", "/api/v1/session", ""), http.StatusUnauthorized},
		{"wrong origin", func() *http.Request {
			r := composedRequest(http.MethodGet, "alice", "/api/v1/session", "")
			r.Header.Set("Origin", "https://evil.example")
			return r
		}(), http.StatusForbidden},
		{"duplicate origin", func() *http.Request {
			r := composedRequest(http.MethodGet, "alice", "/api/v1/session", "")
			r.Header.Add("Origin", privateHTTPOrigin)
			r.Header.Add("Origin", privateHTTPOrigin)
			return r
		}(), http.StatusForbidden},
		{"cross-site fetch metadata", func() *http.Request {
			r := composedRequest(http.MethodGet, "alice", "/api/v1/session", "")
			r.Header.Set("Sec-Fetch-Site", "cross-site")
			return r
		}(), http.StatusForbidden},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, tc.request)
			if w.Code != tc.status {
				t.Fatalf("got %d want %d: %s", w.Code, tc.status, w.Body.String())
			}
		})
	}
	forged := composedRequest(http.MethodGet, "", "/api/v1/session", "")
	forged.Header.Set("X-Hlab-User-Id", "alice")
	forged.Header.Set("X-Hlab-Auth-Method", "session")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, forged)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("untrusted header authorized session: %d", w.Code)
	}
}

func TestPrivateSessionConstructorRejectsUnsafeInputs(t *testing.T) {
	registry, err := NewRegistry([]Site{{ID: "test-site", Name: "Test Site"}})
	if err != nil {
		t.Fatal(err)
	}
	auth := func(*http.Request) (string, error) { return "alice", errors.New("not verified") }
	for _, origin := range []string{"", "http://editor.example.test", "https://editor.example.test/path", "https://user@editor.example.test"} {
		if handler, err := NewPrivateSessionHandler(registry, auth, origin); err == nil || handler != nil {
			t.Fatalf("invalid origin accepted: %q", origin)
		}
	}
	if handler, err := NewPrivateSessionHandler(registry, nil, privateHTTPOrigin); err == nil || handler != nil {
		t.Fatal("missing authenticator accepted")
	}
	if handler, err := NewPrivateSessionHandler(nil, auth, privateHTTPOrigin); err == nil || handler != nil {
		t.Fatal("missing registry accepted")
	}
}

func TestPrivateSessionRechecksRevocationOnEveryRequest(t *testing.T) {
	registry, err := NewRegistry([]Site{{
		ID: "test-site", Name: "Test Site",
		Grants: map[string]Grant{"alice": {Read: true, Write: true}},
	}})
	if err != nil {
		t.Fatal(err)
	}
	allowed := true
	auth := func(*http.Request) (string, error) {
		if !allowed {
			return "", errors.New("session revoked")
		}
		return "alice", nil
	}
	handler, err := NewPrivateSessionHandler(registry, auth, privateHTTPOrigin)
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodGet, privateHTTPOrigin+"/api/v1/session", http.NoBody)
	first := httptest.NewRecorder()
	handler.ServeHTTP(first, request)
	if first.Code != http.StatusOK {
		t.Fatalf("valid session rejected: %d", first.Code)
	}
	allowed = false
	second := httptest.NewRecorder()
	handler.ServeHTTP(second, request)
	if second.Code != http.StatusUnauthorized || strings.Contains(second.Body.String(), "alice") {
		t.Fatalf("revoked session accepted or identity leaked: %d %s", second.Code, second.Body.String())
	}
	if second.Header().Get("Cache-Control") != "private, no-store" {
		t.Fatal("revoked session response must not be cached")
	}
}
