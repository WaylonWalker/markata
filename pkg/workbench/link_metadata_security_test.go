package workbench

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"net/url"
	"sync/atomic"
	"testing"
)

func TestPublicMetadataIP(t *testing.T) {
	for _, test := range []struct {
		address string
		allowed bool
	}{
		{"1.1.1.1", true},
		{"2606:4700:4700::1111", true},
		{"127.0.0.1", false},
		{"::1", false},
		{"::ffff:127.0.0.1", false},
		{"10.1.2.3", false},
		{"172.16.1.1", false},
		{"192.168.1.1", false},
		{"169.254.169.254", false},
		{"100.64.0.1", false},
		{"198.18.0.1", false},
		{"0.0.0.0", false},
		{"224.0.0.1", false},
		{"fc00::1", false},
		{"fe80::1", false},
	} {
		t.Run(test.address, func(t *testing.T) {
			if got := publicMetadataIP(netip.MustParseAddr(test.address)); got != test.allowed {
				t.Fatalf("publicMetadataIP(%q) = %v; want %v", test.address, got, test.allowed)
			}
		})
	}
}

func TestPublicMetadataURL(t *testing.T) {
	for _, test := range []struct {
		raw     string
		allowed bool
	}{
		{"https://example.org/article", true},
		{"https://1.1.1.1/", true},
		{"https://[2606:4700:4700::1111]/", true},
		{"http://localhost/", false},
		{"http://localhost./", false},
		{"http://server.internal/", false},
		{"http://printer.local/", false},
		{"http://127.0.0.1:8080/", false},
		{"http://192.168.1.1/", false},
		{"http://100.64.0.1/", false},
		{"http://169.254.169.254/latest/meta-data/", false},
		{"http://[::1]/", false},
		{"http://user:pass@example.org/", false},
		{"file://example.org/etc/passwd", false},
	} {
		t.Run(test.raw, func(t *testing.T) {
			parsed, err := url.Parse(test.raw)
			if err != nil {
				t.Fatal(err)
			}
			if got := publicMetadataURL(parsed); got != test.allowed {
				t.Fatalf("publicMetadataURL(%q) = %v; want %v", test.raw, got, test.allowed)
			}
			if got := validLinkMetadataURL(test.raw); got != test.allowed {
				t.Fatalf("validLinkMetadataURL(%q) = %v; want %v", test.raw, got, test.allowed)
			}
		})
	}
}

func TestPublicMetadataClientDoesNotContactLocalServer(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newPublicMetadataHTTPClient(context.Background(), 2)
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, http.NoBody)
	if err != nil {
		t.Fatal(err)
	}
	response, err := client.Do(request)
	if err == nil {
		response.Body.Close()
		t.Fatal("metadata client unexpectedly fetched a loopback address")
	}
	if requests.Load() != 0 {
		t.Fatal("metadata request reached a loopback server")
	}
}

func TestPublicMetadataClientRejectsUnsafeRedirect(t *testing.T) {
	client := newPublicMetadataHTTPClient(context.Background(), 2)
	start, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.org/article", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}
	for _, target := range []string{"http://127.0.0.1:8000/admin", "http://169.254.169.254/latest/meta-data/", "file:///etc/passwd"} {
		next, err := http.NewRequestWithContext(context.Background(), http.MethodGet, target, http.NoBody)
		if err != nil {
			t.Fatal(err)
		}
		if err := client.CheckRedirect(next, []*http.Request{start}); err == nil {
			t.Errorf("permitted redirect to %s", target)
		}
	}
}
