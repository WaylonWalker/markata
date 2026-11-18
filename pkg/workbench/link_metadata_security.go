package workbench

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"strings"
	"time"
)

// A browser workbench must not turn URL metadata lookup into a way of
// requesting private services reachable from the Markata host. Reject unsafe
// destinations before and after DNS resolution, and pin the validated address
// for the actual connection (including every redirect).
var forbiddenMetadataNetworks = []netip.Prefix{
	netip.MustParsePrefix("100.64.0.0/10"), // carrier-grade NAT
	netip.MustParsePrefix("192.0.0.0/24"),  // protocol assignments
	netip.MustParsePrefix("192.0.2.0/24"),  // documentation
	netip.MustParsePrefix("198.18.0.0/15"), // benchmarking
	netip.MustParsePrefix("198.51.100.0/24"),
	netip.MustParsePrefix("203.0.113.0/24"),
	netip.MustParsePrefix("240.0.0.0/4"), // reserved/broadcast
	netip.MustParsePrefix("2001:db8::/32"),
}

func publicMetadataIP(ip netip.Addr) bool {
	ip = ip.Unmap()
	if !ip.IsValid() || !ip.IsGlobalUnicast() || ip.IsPrivate() {
		return false
	}
	for _, network := range forbiddenMetadataNetworks {
		if network.Contains(ip) {
			return false
		}
	}
	return true
}

func publicMetadataURL(destination *url.URL) bool {
	if destination == nil || destination.User != nil || destination.Hostname() == "" ||
		(destination.Scheme != "http" && destination.Scheme != "https") {
		return false
	}
	host := strings.ToLower(strings.TrimSuffix(destination.Hostname(), "."))
	if ip, err := netip.ParseAddr(host); err == nil {
		return publicMetadataIP(ip)
	}
	if !strings.Contains(host, ".") {
		return false
	}
	for _, suffix := range []string{".localhost", ".local", ".internal", ".test", ".invalid"} {
		if strings.HasSuffix(host, suffix) {
			return false
		}
	}
	return true
}

type publicMetadataTransport struct {
	base   *http.Transport
	parent context.Context
}

type publicMetadataBody struct {
	io.ReadCloser
	stop   func() bool
	cancel context.CancelFunc
}

func (b *publicMetadataBody) Close() error {
	err := b.ReadCloser.Close()
	b.stop()
	b.cancel()
	return err
}

func (t *publicMetadataTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if !publicMetadataURL(request.URL) {
		return nil, fmt.Errorf("metadata URL must be a public HTTP(S) destination")
	}
	ctx, cancel := context.WithCancel(request.Context())
	stop := context.AfterFunc(t.parent, cancel)
	response, err := t.base.RoundTrip(request.WithContext(ctx))
	if err != nil {
		stop()
		cancel()
		return nil, err
	}
	response.Body = &publicMetadataBody{ReadCloser: response.Body, stop: stop, cancel: cancel}
	return response, nil
}

func newPublicMetadataHTTPClient(parent context.Context, seconds int) *http.Client {
	if seconds <= 0 {
		seconds = 10
	}
	dialer := &net.Dialer{Timeout: 5 * time.Second, KeepAlive: 30 * time.Second}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = nil // never let environment proxies bypass the address policy
	transport.DisableKeepAlives = true
	transport.DialContext = func(ctx context.Context, _ string, address string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(address)
		if err != nil {
			return nil, err
		}
		var addresses []netip.Addr
		if ip, err := netip.ParseAddr(host); err == nil {
			addresses = []netip.Addr{ip}
		} else {
			addresses, err = net.DefaultResolver.LookupNetIP(ctx, "ip", host)
			if err != nil {
				return nil, err
			}
		}
		if len(addresses) == 0 {
			return nil, fmt.Errorf("metadata host has no DNS addresses")
		}
		for _, ip := range addresses {
			if !publicMetadataIP(ip) {
				return nil, fmt.Errorf("metadata host resolves to a non-public address")
			}
		}
		var dialErr error
		for _, ip := range addresses {
			connection, err := dialer.DialContext(ctx, "tcp", net.JoinHostPort(ip.String(), port))
			if err == nil {
				return connection, nil
			}
			dialErr = err
		}
		return nil, dialErr
	}
	return &http.Client{
		Timeout: time.Duration(seconds) * time.Second,
		Transport: &publicMetadataTransport{base: transport, parent: parent},
		CheckRedirect: func(request *http.Request, via []*http.Request) error {
			if len(via) >= 5 {
				return fmt.Errorf("metadata lookup exceeded five redirects")
			}
			if !publicMetadataURL(request.URL) {
				return fmt.Errorf("metadata redirect is not a public HTTP(S) URL")
			}
			return nil
		},
	}
}
