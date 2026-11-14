package builderadmin

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBuilderAdminLLMsTXT_AuthenticatedRoutes(t *testing.T) {
	t.Parallel()

	svc, err := New(Config{
		SiteDir:           t.TempDir(),
		TrustedProxyCIDRs: []string{"10.42.0.0/24"},
		PublicOrigin:      "https://builder.example.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = svc.leaderLock.Close() })
	mux := http.NewServeMux()
	svc.registerRoutes(mux)

	cases := []struct {
		name       string
		path       string
		method     string
		remoteAddr string
		userID     string
		authMethod string
		wantStatus int
	}{
		{"browser SSO session", "/llms.txt", http.MethodGet, "10.42.0.12:443", "operator-123", "", http.StatusOK},
		{"browser unauthenticated", "/llms.txt", http.MethodGet, "10.42.0.12:443", "", "", http.StatusUnauthorized},
		{"browser untrusted", "/llms.txt", http.MethodGet, "192.0.2.8:443", "operator-123", "", http.StatusUnauthorized},
		{"agent authorized", "/api/v1/llms.txt", http.MethodGet, "10.42.0.12:443", "operator-123", "bearer", http.StatusOK},
		{"agent without bearer proof", "/api/v1/llms.txt", http.MethodGet, "10.42.0.12:443", "operator-123", "", http.StatusUnauthorized},
		{"agent browser-only identity", "/api/v1/llms.txt", http.MethodGet, "10.42.0.12:443", "operator-123", "session", http.StatusUnauthorized},
		{"agent unauthenticated", "/api/v1/llms.txt", http.MethodGet, "10.42.0.12:443", "", "", http.StatusUnauthorized},
		{"agent untrusted", "/api/v1/llms.txt", http.MethodGet, "192.0.2.8:443", "operator-123", "bearer", http.StatusUnauthorized},
		{"browser HEAD", "/llms.txt", http.MethodHead, "10.42.0.12:443", "operator-123", "", http.StatusOK},
		{"agent HEAD", "/api/v1/llms.txt", http.MethodHead, "10.42.0.12:443", "operator-123", "bearer", http.StatusOK},
		{"browser write disallowed", "/llms.txt", http.MethodPost, "10.42.0.12:443", "operator-123", "", http.StatusMethodNotAllowed},
		{"agent write disallowed", "/api/v1/llms.txt", http.MethodPost, "10.42.0.12:443", "operator-123", "bearer", http.StatusMethodNotAllowed},
	}
	bodies := make(map[string]string)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			req.RemoteAddr = tc.remoteAddr
			if tc.userID != "" {
				req.Header.Set(DefaultAuthHeaders().UserID, tc.userID)
			}
			if tc.authMethod != "" {
				req.Header.Set(DefaultAuthHeaders().AuthMethod, tc.authMethod)
			}
			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, req)
			if recorder.Code != tc.wantStatus {
				t.Fatalf("%s %s: status = %d, want %d: %s", tc.method, tc.path, recorder.Code, tc.wantStatus, recorder.Body.String())
			}
			if tc.wantStatus == http.StatusMethodNotAllowed {
				if recorder.Header().Get("Allow") != "GET, HEAD" {
					t.Fatalf("wrong allowed methods: %q", recorder.Header().Get("Allow"))
				}
			}
			if tc.wantStatus != http.StatusOK {
				if strings.Contains(recorder.Body.String(), "# markata-go Builder Admin") {
					t.Fatal("non-success response exposed protected guide")
				}
				return
			}
			if got := recorder.Header().Get("Content-Type"); got != "text/plain; charset=utf-8" {
				t.Fatalf("Content-Type = %q", got)
			}
			if got := recorder.Header().Get("Cache-Control"); got != "private, no-store" {
				t.Fatalf("Cache-Control = %q", got)
			}
			if got := recorder.Header().Get("X-Content-Type-Options"); got != "nosniff" {
				t.Fatalf("X-Content-Type-Options = %q", got)
			}
			if len(recorder.Result().Cookies()) != 0 {
				t.Fatal("guide response should not mint session or CSRF cookies")
			}
			if tc.method == http.MethodHead {
				if recorder.Body.Len() != 0 {
					t.Fatal("HEAD request returned a response body")
				}
				return
			}
			body := recorder.Body.String()
			for _, mustContain := range []string{
				"# markata-go Builder Admin",
				"## How to sign in",
				"Authorization: Bearer",
				"/api/v1/state",
				"/api/v1/builds/{build-id}",
				"POST /api/v1/builds",
				"POST /api/v1/releases/{release-id}/rollback",
				"explicit operator approval",
				"ForwardAuth",
			} {
				if !strings.Contains(body, mustContain) {
					t.Errorf("guide missing %q", mustContain)
				}
			}
			bodies[tc.path] = body
		})
	}
	if bodies["/llms.txt"] == "" || bodies["/llms.txt"] != bodies["/api/v1/llms.txt"] {
		t.Fatal("browser and bearer routes must serve the same embedded guide")
	}
}
