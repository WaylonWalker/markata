package editorhub

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"strings"
	"unicode"
	"unicode/utf8"
)

// ForwardAuthConfig describes assertions that a trusted reverse proxy copies
// from a successful ForwardAuth decision. It is not a token validator itself.
// The hosting ingress MUST remove browser-supplied assertion headers first.
type ForwardAuthConfig struct {
	TrustedProxyCIDRs  []string
	SubjectHeader      string
	AuthMethodHeader   string
	AllowedAuthMethods []string
}

var deniedProxyRanges = []netip.Prefix{
	netip.MustParsePrefix("127.0.0.0/8"),
	netip.MustParsePrefix("::1/128"),
	netip.MustParsePrefix("169.254.0.0/16"),
	netip.MustParsePrefix("fe80::/10"),
}

// NewForwardAuthAuthenticator creates an Authenticator for a *verified* ingress
// path. Both the direct network peer and the proxy-provided assertions must
// match explicit allowlists. Use NetworkPolicy to ensure that only the intended
// reverse proxy can access this endpoint; a broad pod CIDR is not sufficient.
func NewForwardAuthAuthenticator(config ForwardAuthConfig) (Authenticator, error) {
	if len(config.TrustedProxyCIDRs) == 0 {
		return nil, errors.New("at least one trusted proxy CIDR is required")
	}
	if len(config.AllowedAuthMethods) == 0 {
		return nil, errors.New("at least one auth method is required")
	}
	if !validAssertionHeader(config.SubjectHeader) || !validAssertionHeader(config.AuthMethodHeader) {
		return nil, errors.New("invalid or missing ForwardAuth assertion header")
	}
	if strings.EqualFold(config.SubjectHeader, config.AuthMethodHeader) {
		return nil, errors.New("subject and auth method headers must differ")
	}
	prefixes := make([]netip.Prefix, 0, len(config.TrustedProxyCIDRs))
	for _, raw := range config.TrustedProxyCIDRs {
		prefix, err := netip.ParsePrefix(raw)
		if err != nil || prefix.Bits() == 0 || prefix.Addr().Is4In6() {
			return nil, fmt.Errorf("invalid trusted proxy CIDR %q", raw)
		}
		prefix = prefix.Masked()
		for _, forbidden := range deniedProxyRanges {
			if prefix.Overlaps(forbidden) {
				return nil, fmt.Errorf("trusted proxy CIDR %q overlaps loopback or link-local space", raw)
			}
		}
		prefixes = append(prefixes, prefix)
	}
	methods := make(map[string]struct{}, len(config.AllowedAuthMethods))
	for _, method := range config.AllowedAuthMethods {
		if method == "" || len(method) > 64 || !validMethodName(method) {
			return nil, fmt.Errorf("invalid allowed auth method %q", method)
		}
		methods[strings.ToLower(method)] = struct{}{}
	}
	subjectHeader := http.CanonicalHeaderKey(config.SubjectHeader)
	methodHeader := http.CanonicalHeaderKey(config.AuthMethodHeader)
	return func(r *http.Request) (string, error) {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return "", errors.New("request has no trusted network peer")
		}
		peer, err := netip.ParseAddr(host)
		if err != nil || !peer.IsValid() || peer.Is4In6() {
			return "", errors.New("request has invalid network peer")
		}
		trusted := false
		for _, prefix := range prefixes {
			if prefix.Contains(peer) {
				trusted = true
				break
			}
		}
		if !trusted {
			return "", errors.New("request did not originate from a trusted proxy")
		}
		subject, err := singleAssertion(r.Header, subjectHeader, 256)
		if err != nil {
			return "", err
		}
		method, err := singleAssertion(r.Header, methodHeader, 64)
		if err != nil {
			return "", err
		}
		if _, ok := methods[strings.ToLower(method)]; !ok {
			return "", errors.New("auth method not allowed")
		}
		return subject, nil
	}, nil
}

func validAssertionHeader(name string) bool {
	if name == "" {
		return false
	}
	for _, r := range name {
		if r > unicode.MaxASCII || !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-') {
			return false
		}
	}
	return true
}

func validMethodName(method string) bool {
	for _, r := range method {
		if r > unicode.MaxASCII || !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_') {
			return false
		}
	}
	return true
}

func singleAssertion(headers http.Header, name string, maxLength int) (string, error) {
	values := headers.Values(name)
	if len(values) != 1 {
		return "", errors.New("missing or duplicated ForwardAuth assertion")
	}
	value := values[0]
	if value == "" || len(value) > maxLength || !utf8.ValidString(value) ||
		strings.TrimSpace(value) != value || strings.ContainsAny(value, ",\r\n\x00") {
		return "", errors.New("malformed ForwardAuth assertion")
	}
	for _, r := range value {
		if unicode.IsControl(r) {
			return "", errors.New("control character in ForwardAuth assertion")
		}
	}
	return value, nil
}
