package workbench

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValidateHostedAuthConfig(t *testing.T) {
	valid := HostedAuthConfig{
		PublicOrigin:      "https://edit.waylonwalker.com",
		TrustedProxyCIDRs: []string{"10.42.0.0/24"},
		RequiredScope:     "markata-editor:admin",
	}
	if err := ValidateHostedAuthConfig(valid); err != nil {
		t.Fatalf("valid config rejected: %v", err)
	}

	tests := []struct {
		name   string
		change func(*HostedAuthConfig)
	}{
		{name: "missing origin", change: func(c *HostedAuthConfig) { c.PublicOrigin = "" }},
		{name: "HTTP origin", change: func(c *HostedAuthConfig) { c.PublicOrigin = "http://edit.waylonwalker.com" }},
		{name: "origin path", change: func(c *HostedAuthConfig) { c.PublicOrigin += "/editor" }},
		{name: "origin port", change: func(c *HostedAuthConfig) { c.PublicOrigin = "https://edit.waylonwalker.com:8443" }},
		{name: "non-canonical origin", change: func(c *HostedAuthConfig) { c.PublicOrigin = "https://EDIT.waylonwalker.com" }},
		{name: "missing proxy", change: func(c *HostedAuthConfig) { c.TrustedProxyCIDRs = nil }},
		{name: "universal proxy", change: func(c *HostedAuthConfig) { c.TrustedProxyCIDRs = []string{"0.0.0.0/0"} }},
		{name: "invalid proxy", change: func(c *HostedAuthConfig) { c.TrustedProxyCIDRs = []string{"traefik"} }},
		{name: "non-app scope", change: func(c *HostedAuthConfig) { c.RequiredScope = "admin" }},
		{name: "viewer scope", change: func(c *HostedAuthConfig) { c.RequiredScope = "markata-editor:read" }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := valid
			config.TrustedProxyCIDRs = append([]string(nil), valid.TrustedProxyCIDRs...)
			test.change(&config)
			if err := ValidateHostedAuthConfig(config); err == nil {
				t.Fatal("expected configuration error")
			}
		})
	}
}

func TestHostedAuthGuardsEveryRouteAndDisablesSharedCaching(t *testing.T) {
	h := newHandler(t)
	if err := h.EnableHostedAuth(HostedAuthConfig{
		PublicOrigin:      "https://edit.waylonwalker.com",
		TrustedProxyCIDRs: []string{"10.42.0.0/24"},
		RequiredScope:     "markata-editor:admin",
	}); err != nil {
		t.Fatal(err)
	}

	for _, path := range []string{"/", "/editor/lsp.js", "/appearance/css/main.css", "/lsp", "/workspace/paths"} {
		t.Run(path, func(t *testing.T) {
			request := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "https://edit.waylonwalker.com"+path, http.NoBody)
			request.RemoteAddr = "192.0.2.5:443"
			request.Header.Set("X-Hlab-User-Id", "usr_123")
			request.Header.Set("X-Hlab-Scopes", "markata-editor:admin")
			response := httptest.NewRecorder()
			h.ServeHTTP(response, request)
			if response.Code != http.StatusUnauthorized {
				t.Fatalf("untrusted peer status = %d, want %d", response.Code, http.StatusUnauthorized)
			}
			if response.Header().Get("Cache-Control") != "private, no-store" {
				t.Fatalf("unauthorized response cache policy = %q", response.Header().Get("Cache-Control"))
			}
		})
	}

	request := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "https://edit.waylonwalker.com/editor/lsp.js", http.NoBody)
	request.RemoteAddr = "10.42.0.12:8443"
	request.Header.Set("X-Hlab-User-Id", "usr_123")
	request.Header.Set("X-Hlab-Scopes", "markata-editor:admin")
	response := httptest.NewRecorder()
	h.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("trusted authenticated asset status = %d: %s", response.Code, response.Body.String())
	}
	if response.Header().Get("Cache-Control") != "private, no-store" {
		t.Fatalf("authenticated response cache policy = %q", response.Header().Get("Cache-Control"))
	}
	if response.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Fatalf("hosted asset unexpectedly enabled CORS: %q", response.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestHostedAuthRejectsMalformedIdentityAndInsufficientScope(t *testing.T) {
	h := newHandler(t)
	if err := h.EnableHostedAuth(HostedAuthConfig{
		PublicOrigin:      "https://edit.waylonwalker.com",
		TrustedProxyCIDRs: []string{"10.42.0.0/24"},
		RequiredScope:     "markata-editor:admin",
	}); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		prepare func(*http.Request)
	}{
		{name: "missing identity", prepare: func(*http.Request) {}},
		{name: "duplicate identity", prepare: func(r *http.Request) {
			r.Header.Add("X-Hlab-User-Id", "usr_123")
			r.Header.Add("X-Hlab-User-Id", "usr_456")
			r.Header.Set("X-Hlab-Scopes", "markata-editor:admin")
		}},
		{name: "malformed identity", prepare: func(r *http.Request) {
			r.Header.Set("X-Hlab-User-Id", "usr_123,usr_456")
			r.Header.Set("X-Hlab-Scopes", "markata-editor:admin")
		}},
		{name: "insufficient scope", prepare: func(r *http.Request) {
			r.Header.Set("X-Hlab-User-Id", "usr_123")
			r.Header.Set("X-Hlab-Scopes", "markata-editor:read,markata-editor:admin-extra")
		}},
		{name: "malformed scope list", prepare: func(r *http.Request) {
			r.Header.Set("X-Hlab-User-Id", "usr_123")
			r.Header.Set("X-Hlab-Scopes", "markata-editor:admin,")
		}},
		{name: "wrong host", prepare: func(r *http.Request) {
			r.Host = "build.waylonwalker.com"
			r.Header.Set("X-Hlab-User-Id", "usr_123")
			r.Header.Set("X-Hlab-Scopes", "markata-editor:admin")
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "https://edit.waylonwalker.com/", http.NoBody)
			request.RemoteAddr = "10.42.0.12:8443"
			test.prepare(request)
			response := httptest.NewRecorder()
			h.ServeHTTP(response, request)
			if response.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d, want %d: %s", response.Code, http.StatusUnauthorized, response.Body.String())
			}
		})
	}
}

func TestHostedWorkspaceMutationRequiresExactPublicOrigin(t *testing.T) {
	h := newHandler(t)
	if err := h.EnableHostedAuth(HostedAuthConfig{
		PublicOrigin:      "https://edit.waylonwalker.com",
		TrustedProxyCIDRs: []string{"10.42.0.0/24"},
		RequiredScope:     "markata-editor:admin",
	}); err != nil {
		t.Fatal(err)
	}

	for _, origin := range []string{
		"https://edit.waylonwalker.com",
		"https://edit.waylonwalker.com.evil.example",
		"http://edit.waylonwalker.com",
		"https://edit.waylonwalker.com:443",
	} {
		t.Run(origin, func(t *testing.T) {
			request := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "https://edit.waylonwalker.com/workspace/file", strings.NewReader(""))
			request.RemoteAddr = "10.42.0.12:8443"
			request.Header.Set("Origin", origin)
			if got := h.sameOriginWorkspaceMutation(request); got != (origin == "https://edit.waylonwalker.com") {
				t.Fatalf("same-origin result = %t", got)
			}
		})
	}
	request := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "https://edit.waylonwalker.com/workspace/file", strings.NewReader(""))
	request.Header.Set("Origin", "https://edit.waylonwalker.com")
	request.Header.Set("Sec-Fetch-Site", "cross-site")
	if h.sameOriginWorkspaceMutation(request) {
		t.Fatal("cross-site fetch metadata was accepted")
	}
}

func TestHostedLSPOriginRequiresAuthenticatedExactOrigin(t *testing.T) {
	h := newHandler(t)
	if err := h.EnableHostedAuth(HostedAuthConfig{
		PublicOrigin:      "https://edit.waylonwalker.com",
		TrustedProxyCIDRs: []string{"10.42.0.0/24"},
		RequiredScope:     "markata-editor:admin",
	}); err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "https://edit.waylonwalker.com/lsp", http.NoBody)
	request.RemoteAddr = "10.42.0.12:8443"
	request.Header.Set("X-Hlab-User-Id", "usr_123")
	request.Header.Set("X-Hlab-Scopes", "markata-editor:admin")
	request.Header.Set("Origin", "https://edit.waylonwalker.com")
	if !h.validLSPOrigin(request) {
		t.Fatal("authenticated same-origin LSP request was rejected")
	}
	request.Header.Set("Origin", "https://attacker.example")
	if h.validLSPOrigin(request) {
		t.Fatal("cross-origin LSP request was accepted")
	}
}
