package editorhub

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type testSubjectKey struct{}

func fixture(t *testing.T) (*Registry, *Handler) {
	t.Helper()
	registry, err := NewRegistry([]Site{
		{ID: "z-personal", Name: "Personal", Grants: map[string]Grant{"writer-a": {Read: true, Write: true, Publish: true}}},
		{ID: "a-shared", Name: "Family", Grants: map[string]Grant{"writer-a": {Read: true, Write: true}, "writer-b": {Read: true}}},
		{ID: "secret", Name: "Hidden", Grants: map[string]Grant{"operator": {Read: true, Write: true, Publish: true}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	handler, err := NewHandler(registry, func(r *http.Request) (string, error) {
		subject, ok := r.Context().Value(testSubjectKey{}).(string)
		if !ok {
			return "", errors.New("identity not authenticated")
		}
		return subject, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return registry, handler
}

func serve(t *testing.T, h *Handler, method, path, subject string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, nil)
	if subject != "" {
		req = req.WithContext(context.WithValue(req.Context(), testSubjectKey{}, subject))
	}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	return rr
}

func TestDiscoveryIsolatedBySubject(t *testing.T) {
	_, h := fixture(t)
	for _, tc := range []struct {
		subject string
		want    []string
	}{
		{"writer-a", []string{"a-shared", "z-personal"}},
		{"writer-b", []string{"a-shared"}},
		{"operator", []string{"secret"}},
		{"unknown", []string{}},
	} {
		got := serve(t, h, "GET", "/api/v1/sites", tc.subject)
		if got.Code != 200 {
			t.Fatalf("%q: status %d", tc.subject, got.Code)
		}
		var data struct {
			Sites []SiteSummary `json:"sites"`
		}
		if err := json.Unmarshal(got.Body.Bytes(), &data); err != nil {
			t.Fatal(err)
		}
		ids := make([]string, 0, len(data.Sites))
		for _, site := range data.Sites {
			ids = append(ids, site.ID)
		}
		if !reflect.DeepEqual(ids, tc.want) {
			t.Errorf("%s: got %v, want %v", tc.subject, ids, tc.want)
		}
		if got.Header().Get("Cache-Control") != "private, no-store" {
			t.Error("missing no-store")
		}
		if got.Header().Get("Access-Control-Allow-Origin") != "" {
			t.Error("unexpected CORS")
		}
		if strings.Contains(got.Body.String(), "git@") || strings.Contains(got.Body.String(), "operator") {
			t.Error("sensitive details leaked")
		}
	}
}

func TestFailClosedAndNoWorkbenchSurface(t *testing.T) {
	_, h := fixture(t)
	cases := []struct {
		method, path string
		want         int
	}{
		{"GET", "/api/v1/sites", 401},
		{"GET", "/workspace/files", 404},
		{"GET", "/api/v1/state", 404},
		{"POST", "/api/v1/sites", 405},
		{"GET", "/api/v1/sites/secret", 404},
	}
	for _, tc := range cases {
		rr := serve(t, h, tc.method, tc.path, "")
		if rr.Code != tc.want {
			t.Errorf("%s %s: got %d want %d", tc.method, tc.path, rr.Code, tc.want)
		}
	}
	req := httptest.NewRequest("GET", "/api/v1/sites", nil)
	req.Header.Set("X-Hlab-User-Id", "writer-a")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != 401 {
		t.Error("untrusted identity header bypassed auth")
	}
	req = httptest.NewRequest("GET", "/api/v1/sites", nil)
	req = req.WithContext(context.WithValue(req.Context(), testSubjectKey{}, "writer-a"))
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != 403 {
		t.Error("cross-site request not rejected")
	}
}

func TestRegistryValidationAndCopy(t *testing.T) {
	for _, tc := range []Site{
		{ID: "bad/path", Name: "Bad"},
		{ID: "ok", Name: " "},
		{ID: "ok", Name: "Ok", Grants: map[string]Grant{"*": {Read: true}}},
		{ID: "ok", Name: "Ok", Grants: map[string]Grant{"x": {Write: true}}},
		{ID: "ok", Name: "Ok", Grants: map[string]Grant{"x": {Read: true, Publish: true}}},
	} {
		if _, err := NewRegistry([]Site{tc}); err == nil {
			t.Errorf("accepted invalid site %+v", tc)
		}
	}
	if _, err := NewRegistry([]Site{{ID: "same", Name: "One"}, {ID: "same", Name: "Two"}}); err == nil {
		t.Error("accepted duplicate")
	}
	if _, err := NewHandler(nil, func(*http.Request) (string, error) { return "a", nil }); err == nil {
		t.Error("accepted nil registry")
	}
	input := []Site{{ID: "site", Name: "Site", Grants: map[string]Grant{"a": {Read: true}}}}
	registry, err := NewRegistry(input)
	if err != nil {
		t.Fatal(err)
	}
	input[0].Grants["b"] = Grant{Read: true}
	if len(registry.ListFor("b")) != 0 {
		t.Error("registry references mutable config")
	}
	view := registry.ListFor("a")
	view[0].Name = "tampered"
	if registry.ListFor("a")[0].Name != "Site" {
		t.Error("caller mutated registry")
	}
}
