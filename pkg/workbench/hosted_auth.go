package workbench

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// HostedAuthConfig enables the workbench behind an authenticated reverse proxy.
// Configure it before serving requests. Local workbenches do not need this.
type HostedAuthConfig struct {
	PublicOrigin      string
	TrustedProxyCIDRs []string
	RequiredScope     string
}

type hostedAuth struct {
	origin      string
	host        string
	required    string
	proxyRanges []*net.IPNet
}

var (
	hostedScopePattern  = regexp.MustCompile("^[a-z0-9][a-z0-9._:-]*$")
	hostedUserIDPattern = regexp.MustCompile("^[A-Za-z0-9_-]{1,128}$")
)

// ValidateHostedAuthConfig checks the complete fail-closed hosted-mode config.
func ValidateHostedAuthConfig(config HostedAuthConfig) error {
	_, err := parseHostedAuthConfig(config)
	return err
}

func parseHostedAuthConfig(config HostedAuthConfig) (*hostedAuth, error) {
	host, canonicalOrigin, err := parseHostedPublicOrigin(config.PublicOrigin)
	if err != nil {
		return nil, err
	}
	if !hostedScopePattern.MatchString(config.RequiredScope) ||
		!strings.Contains(config.RequiredScope, ":") ||
		(!strings.HasSuffix(config.RequiredScope, ":admin") && !strings.HasSuffix(config.RequiredScope, ":write")) {
		return nil, fmt.Errorf("hosted editor required scope must be an app-specific scope")
	}
	proxyRanges, err := parseHostedProxyCIDRs(config.TrustedProxyCIDRs)
	if err != nil {
		return nil, err
	}

	return &hostedAuth{origin: canonicalOrigin, host: host, required: config.RequiredScope, proxyRanges: proxyRanges}, nil
}

func parseHostedPublicOrigin(raw string) (host, canonicalOrigin string, err error) {
	origin, err := url.Parse(raw)
	if err != nil || origin == nil || origin.Scheme != secureScheme || origin.User != nil || origin.Opaque != "" || origin.Host == "" || origin.Path != "" || origin.RawPath != "" || origin.RawQuery != "" || origin.Fragment != "" {
		return "", "", fmt.Errorf("hosted editor public origin must be a canonical HTTPS origin")
	}
	host = origin.Hostname()
	if host == "" || origin.Port() != "" || net.ParseIP(host) != nil || strings.Contains(host, "..") || !strings.Contains(host, ".") {
		return "", "", fmt.Errorf("hosted editor public origin must use a DNS hostname without a port")
	}
	canonicalOrigin = secureScheme + "://" + strings.ToLower(origin.Host)
	if raw != canonicalOrigin {
		return "", "", fmt.Errorf("hosted editor public origin must be canonical: %s", canonicalOrigin)
	}
	return strings.ToLower(host), canonicalOrigin, nil
}

func parseHostedProxyCIDRs(cidrs []string) ([]*net.IPNet, error) {
	if len(cidrs) == 0 {
		return nil, fmt.Errorf("hosted editor requires at least one trusted proxy CIDR")
	}
	ranges := make([]*net.IPNet, 0, len(cidrs))
	for _, cidr := range cidrs {
		_, network, parseErr := net.ParseCIDR(strings.TrimSpace(cidr))
		if parseErr != nil {
			return nil, fmt.Errorf("invalid trusted proxy CIDR %q: %w", cidr, parseErr)
		}
		ones, _ := network.Mask.Size()
		if ones == 0 {
			return nil, fmt.Errorf("trusted proxy CIDR %q must not be universal", cidr)
		}
		ranges = append(ranges, network)
	}
	return ranges, nil
}

// EnableHostedAuth turns on the hosted request boundary. Call it before serving.
func (h *Handler) EnableHostedAuth(config HostedAuthConfig) error {
	if h == nil {
		return fmt.Errorf("hosted editor handler is nil")
	}
	if h.hostedAuth != nil {
		return fmt.Errorf("hosted editor authentication is already configured")
	}
	parsed, err := parseHostedAuthConfig(config)
	if err != nil {
		return err
	}
	h.hostedAuth = parsed
	return nil
}

func (h *Handler) authorizeHostedRequest(r *http.Request) bool {
	if h.hostedAuth == nil || !hostedPeerAllowed(r.RemoteAddr, h.hostedAuth.proxyRanges) {
		return false
	}
	if !strings.EqualFold(r.Host, h.hostedAuth.host) {
		return false
	}
	userIDs := r.Header.Values("X-Hlab-User-Id")
	if len(userIDs) != 1 || !hostedUserIDPattern.MatchString(userIDs[0]) {
		return false
	}
	values := r.Header.Values("X-Hlab-Scopes")
	if len(values) == 0 {
		return false
	}
	found := false
	for _, value := range values {
		for _, item := range strings.Split(value, ",") {
			scope := strings.TrimSpace(item)
			if !hostedScopePattern.MatchString(scope) {
				return false
			}
			if scope == h.hostedAuth.required {
				found = true
			}
		}
	}
	return found
}

type privateNoStoreResponseWriter struct {
	http.ResponseWriter
}

func (w *privateNoStoreResponseWriter) WriteHeader(status int) {
	w.Header().Set("Cache-Control", "private, no-store")
	w.ResponseWriter.WriteHeader(status)
}

func (w *privateNoStoreResponseWriter) Write(data []byte) (int, error) {
	w.Header().Set("Cache-Control", "private, no-store")
	return w.ResponseWriter.Write(data)
}

func (w *privateNoStoreResponseWriter) Flush() {
	w.Header().Set("Cache-Control", "private, no-store")
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (w *privateNoStoreResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	w.Header().Set("Cache-Control", "private, no-store")
	hijacker, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("response writer does not support connection hijacking")
	}
	return hijacker.Hijack()
}

func (w *privateNoStoreResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func hostedPeerAllowed(remoteAddr string, ranges []*net.IPNet) bool {
	peer, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return false
	}
	ip := net.ParseIP(peer)
	if ip == nil {
		return false
	}
	for _, network := range ranges {
		if network.Contains(ip) || (ip.To4() != nil && network.Contains(ip.To4())) {
			return true
		}
	}
	return false
}
