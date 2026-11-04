package editorhub

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func trustedConfig() ForwardAuthConfig {
	return ForwardAuthConfig{
		TrustedProxyCIDRs:  []string{"10.42.7.19/32"},
		SubjectHeader:      "X-Hlab-User-Id",
		AuthMethodHeader:   "X-Hlab-Auth-Method",
		AllowedAuthMethods: []string{"session"},
	}
}

func TestForwardAuthEndToEndSiteDiscovery(t *testing.T) {
	registry, err := NewRegistry([]Site{
		{ID: "family", Name: "Family", Grants: map[string]Grant{
			"alice": {Read: true, Write: true},
			"bob":   {Read: true},
		}},
		{ID: "private", Name: "Personal", Grants: map[string]Grant{
			"alice": {Read: true, Write: true, Publish: true},
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	auth, err := NewForwardAuthAuthenticator(trustedConfig())
	if err != nil {
		t.Fatal(err)
	}
	handler, err := NewHandler(registry, auth)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		subject   string
		wantCount int
	}{
		{"alice", 2},
		{"bob", 1},
		{"unknown", 0},
	} {
		req := httptest.NewRequestWithContext(context.Background(), "GET", "/api/v1/sites", http.NoBody)
		req.RemoteAddr = "10.42.7.19:53535"
		req.Header.Set("X-Hlab-User-Id", tc.subject)
		req.Header.Set("X-Hlab-Auth-Method", "session")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("%s: status %d", tc.subject, w.Code)
		}
		var payload struct {
			Sites []SiteSummary `json:"sites"`
		}
		if err := json.NewDecoder(w.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		if len(payload.Sites) != tc.wantCount {
			t.Errorf("%s: got %d sites, want %d", tc.subject, len(payload.Sites), tc.wantCount)
		}
	}
}

func TestForwardAuthRejectsUntrustedRequests(t *testing.T) {
	auth, err := NewForwardAuthAuthenticator(trustedConfig())
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name, remote, subject, method string
		dupSubject, dupMethod         bool
	}{
		{"untrusted remote", "10.42.7.20:1234", "alice", "session", false, false},
		{"loopback", "127.0.0.1:1234", "alice", "session", false, false},
		{"no port", "10.42.7.19", "alice", "session", false, false},
		{"malformed remote", "not-an-ip:1234", "alice", "session", false, false},
		{"missing subject", "10.42.7.19:1234", "", "session", false, false},
		{"missing method", "10.42.7.19:1234", "alice", "", false, false},
		{"wrong auth method", "10.42.7.19:1234", "alice", "bearer", false, false},
		{"duplicated subject", "10.42.7.19:1234", "alice", "session", true, false},
		{"duplicated method", "10.42.7.19:1234", "alice", "session", false, true},
		{"comma subject", "10.42.7.19:1234", "alice,bob", "session", false, false},
		{"space subject", "10.42.7.19:1234", " alice", "session", false, false},
		{"control subject", "10.42.7.19:1234", "a	b", "session", false, false},
		{"comma method", "10.42.7.19:1234", "alice", "session,bearer", false, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequestWithContext(context.Background(), "GET", "/api/v1/sites", http.NoBody)
			req.RemoteAddr = tc.remote
			if tc.subject != "" {
				req.Header.Add("X-Hlab-User-Id", tc.subject)
			}
			if tc.method != "" {
				req.Header.Add("X-Hlab-Auth-Method", tc.method)
			}
			// A forged client X-Forwarded-For is never a substitute for
			// the direct network peer used by the trust decision.
			req.Header.Set("X-Forwarded-For", "10.42.7.19")
			if tc.dupSubject {
				req.Header.Add("X-Hlab-User-Id", "mallory")
			}
			if tc.dupMethod {
				req.Header.Add("X-Hlab-Auth-Method", "bearer")
			}
			if subject, err := auth(req); err == nil {
				t.Errorf("accepted subject %q", subject)
			}
		})
	}
}

func TestForwardAuthConfigFailsClosed(t *testing.T) {
	cases := []struct {
		name   string
		change func(*ForwardAuthConfig)
	}{
		{"no proxies", func(c *ForwardAuthConfig) { c.TrustedProxyCIDRs = nil }},
		{"world CIDR", func(c *ForwardAuthConfig) { c.TrustedProxyCIDRs = []string{"0.0.0.0/0"} }},
		{"loopback CIDR", func(c *ForwardAuthConfig) { c.TrustedProxyCIDRs = []string{"127.0.0.1/32"} }},
		{"link-local CIDR", func(c *ForwardAuthConfig) { c.TrustedProxyCIDRs = []string{"169.254.1.2/32"} }},
		{"overlapping CIDR", func(c *ForwardAuthConfig) { c.TrustedProxyCIDRs = []string{"0.0.0.0/1"} }},
		{"invalid CIDR", func(c *ForwardAuthConfig) { c.TrustedProxyCIDRs = []string{"not-a-cidr"} }},
		{"broad private network", func(c *ForwardAuthConfig) { c.TrustedProxyCIDRs = []string{"10.0.0.0/8"} }},
		{"broad public network", func(c *ForwardAuthConfig) { c.TrustedProxyCIDRs = []string{"128.0.0.0/1"} }},
		{"broad IPv6 network", func(c *ForwardAuthConfig) { c.TrustedProxyCIDRs = []string{"2001:db8::/32"} }},
		{"no subject header", func(c *ForwardAuthConfig) { c.SubjectHeader = "" }},
		{"invalid subject header", func(c *ForwardAuthConfig) { c.SubjectHeader = "X-Identity
X-Other:" }},
		{"same header", func(c *ForwardAuthConfig) { c.SubjectHeader = c.AuthMethodHeader }},
		{"no method header", func(c *ForwardAuthConfig) { c.AuthMethodHeader = "" }},
		{"no allowed methods", func(c *ForwardAuthConfig) { c.AllowedAuthMethods = nil }},
		{"wildcard method", func(c *ForwardAuthConfig) { c.AllowedAuthMethods = []string{"*"} }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := trustedConfig()
			tc.change(&c)
			if _, err := NewForwardAuthAuthenticator(c); err == nil {
				t.Error("accepted unsafe configuration")
			}
		})
	}
}

func TestForwardAuthAcceptedMethodIsCaseInsensitive(t *testing.T) {
	auth, err := NewForwardAuthAuthenticator(trustedConfig())
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequestWithContext(context.Background(), "GET", "/api/v1/sites", http.NoBody)
	req.RemoteAddr = "10.42.7.19:1234"
	req.Header.Set("X-Hlab-User-Id", "alice")
	req.Header.Set("X-Hlab-Auth-Method", "SESSION")
	if subject, err := auth(req); err != nil || subject != "alice" {
		t.Fatalf("expected verified alice, got %q, %v", subject, err)
	}
}
