package builderadmin

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	csrfCookieName        = "__Host-markata-builder-admin-csrf"
	forwardedLeaderHeader = "X-Markata-Builder-Admin-Forwarded"
	apiAuthMethodBearer   = "bearer"
)

var forbiddenTrustedProxyPrefixes = []netip.Prefix{
	netip.MustParsePrefix("127.0.0.0/8"),
	netip.MustParsePrefix("::1/128"),
	netip.MustParsePrefix("169.254.0.0/16"),
	netip.MustParsePrefix("fe80::/10"),
}

func newCSRFToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("read CSRF randomness: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func csrfCookie(token string) *http.Cookie {
	return &http.Cookie{Name: csrfCookieName, Value: token, Path: "/", Secure: true, HttpOnly: true, SameSite: http.SameSiteStrictMode}
}

// validateCSRF runs only on the active leader, after a standby has forwarded a
// mutation. This keeps cookie validation local to the process that changes state.
func (s *Service) validateCSRF(w http.ResponseWriter, r *http.Request) bool {
	if s.cfg.PublicOrigin == "" || r.Header.Get("Origin") != s.cfg.PublicOrigin {
		http.Error(w, "invalid CSRF origin", http.StatusForbidden)
		return false
	}
	if site := r.Header.Get("Sec-Fetch-Site"); site != "" && site != "same-origin" {
		http.Error(w, "invalid CSRF fetch site", http.StatusForbidden)
		return false
	}
	cookie, err := r.Cookie(csrfCookieName)
	if err != nil || cookie.Value == "" {
		http.Error(w, "missing CSRF token", http.StatusForbidden)
		return false
	}
	token := r.FormValue("csrf_token")
	if headerToken := r.Header.Get("X-CSRF-Token"); headerToken != "" {
		token = headerToken
	}
	if token == "" || subtle.ConstantTimeCompare([]byte(cookie.Value), []byte(token)) != 1 {
		http.Error(w, "invalid CSRF token", http.StatusForbidden)
		return false
	}
	return true
}

// OperatorProfile contains display-only assertions supplied by hlab-auth. UserID
// is the durable identity key; no authorization decision is based on the other fields.
type OperatorProfile struct {
	UserID      string
	Username    string
	DisplayName string
	Email       string
	Groups      string
	Roles       string
	Scopes      string
}

// AuthHeaders maps identity assertions and the machine-authentication marker
// supplied by a trusted ForwardAuth proxy. Empty optional header names disable
// that assertion. The default preserves the stable hlab-auth header contract.
type AuthHeaders struct {
	UserID      string
	Username    string
	DisplayName string
	Email       string
	Groups      string
	Roles       string
	Scopes      string
	AuthMethod  string
}

// DefaultAuthHeaders returns the header contract used by hlab-auth.
func DefaultAuthHeaders() AuthHeaders {
	return AuthHeaders{
		UserID:      "X-Hlab-User-Id",
		Username:    "X-Hlab-Username",
		DisplayName: "X-Hlab-Display-Name",
		Email:       "X-Hlab-Email",
		Groups:      "X-Hlab-Groups",
		Roles:       "X-Hlab-Roles",
		Scopes:      "X-Hlab-Scopes",
		AuthMethod:  "X-Hlab-Auth-Method",
	}
}

func normalizeAuthHeaders(headers AuthHeaders) (AuthHeaders, error) {
	if headers == (AuthHeaders{}) {
		headers = DefaultAuthHeaders()
	} else if headers.UserID == "" {
		return AuthHeaders{}, fmt.Errorf("auth user ID header is required")
	}
	for _, field := range []struct {
		name  string
		value string
	}{
		{"user ID", headers.UserID},
		{"username", headers.Username},
		{"display name", headers.DisplayName},
		{"email", headers.Email},
		{"groups", headers.Groups},
		{"roles", headers.Roles},
		{"scopes", headers.Scopes},
		{"auth method", headers.AuthMethod},
	} {
		if field.value != "" && !validHTTPHeaderName(field.value) {
			return AuthHeaders{}, fmt.Errorf("invalid auth %s header %q", field.name, field.value)
		}
	}
	if headers.AuthMethod != "" {
		for _, identityHeader := range []string{
			headers.UserID,
			headers.Username,
			headers.DisplayName,
			headers.Email,
			headers.Groups,
			headers.Roles,
			headers.Scopes,
		} {
			if identityHeader != "" && strings.EqualFold(headers.AuthMethod, identityHeader) {
				return AuthHeaders{}, fmt.Errorf("auth method header must differ from identity header %q", identityHeader)
			}
		}
	}
	return headers, nil
}

func validHTTPHeaderName(name string) bool {
	if name == "" {
		return false
	}
	for _, r := range name {
		if r > unicode.MaxASCII || (!unicode.IsLetter(r) && !unicode.IsDigit(r) && !strings.ContainsRune("!#$%&'*+-.^_`|~", r)) {
			return false
		}
	}
	return true
}

func parseTrustedProxyPrefixes(cidrs []string) ([]netip.Prefix, error) {
	prefixes := make([]netip.Prefix, 0, len(cidrs))
	for _, cidr := range cidrs {
		prefix, err := netip.ParsePrefix(strings.TrimSpace(cidr))
		if err != nil {
			return nil, fmt.Errorf("invalid trusted proxy CIDR %q: %w", cidr, err)
		}
		if prefix.Bits() == 0 {
			return nil, fmt.Errorf("trusted proxy CIDR %q must not match all addresses", cidr)
		}
		prefix = prefix.Masked()
		for _, forbidden := range forbiddenTrustedProxyPrefixes {
			if prefix.Overlaps(forbidden) {
				return nil, fmt.Errorf("trusted proxy CIDR %q must not include loopback or link-local addresses", cidr)
			}
		}
		prefixes = append(prefixes, prefix)
	}
	return prefixes, nil
}

func normalizePublicAuthOrigin(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", nil
	}
	origin, err := url.Parse(raw)
	if err != nil || origin.Scheme != "https" || origin.Host == "" || origin.User != nil || origin.RawQuery != "" || origin.Fragment != "" {
		return "", fmt.Errorf("public auth origin must be an HTTPS origin without credentials, query, or fragment")
	}
	if origin.Path != "" && origin.Path != "/" {
		return "", fmt.Errorf("public auth origin must not include a path")
	}
	return strings.TrimSuffix(origin.String(), "/"), nil
}

func profilePictureURL(publicAuthOrigin, userID string) string {
	if publicAuthOrigin == "" || userID == "" {
		return ""
	}
	return publicAuthOrigin + "/users/" + url.PathEscape(userID) + "/picture"
}

func (s *Service) requireTrustedOperator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The leader-forwarding marker distinguishes peer forwarding, but it
		// never substitutes for source provenance.
		if !s.isTrustedProxyRequest(r) {
			http.Error(w, "unauthenticated", http.StatusUnauthorized)
			return
		}
		if _, err := operatorProfileFromHeaders(r.Header, s.cfg.AuthHeaders); err != nil {
			http.Error(w, "unauthenticated", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// requireBearerOperator protects the machine API with a bearer-authentication
// assertion returned by ForwardAuth. The configured header is stripped at the
// public proxy before authentication and only copied from a successful auth
// response. Local headers and scope display assertions never grant access.
func (s *Service) requireBearerOperator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !s.isTrustedProxyRequest(r) {
			writeAPIError(w, http.StatusUnauthorized, "authentication_required", "A trusted ForwardAuth identity is required.")
			return
		}
		if _, err := operatorProfileFromHeaders(r.Header, s.cfg.AuthHeaders); err != nil {
			writeAPIError(w, http.StatusUnauthorized, "authentication_required", "A trusted ForwardAuth identity is required.")
			return
		}
		method, err := singleTrustedHeader(r.Header, s.cfg.AuthHeaders.AuthMethod, true)
		if err != nil || !strings.EqualFold(method, apiAuthMethodBearer) {
			writeAPIError(w, http.StatusUnauthorized, "authentication_required", "A validated bearer token is required.")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Service) isTrustedProxyRequest(r *http.Request) bool {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return false
	}
	address, err := netip.ParseAddr(host)
	if err != nil {
		return false
	}
	for _, prefix := range s.trustedProxyPrefixes {
		if prefix.Contains(address) {
			return true
		}
	}
	return false
}

func operatorProfileFromHeaders(headers http.Header, authHeaders AuthHeaders) (OperatorProfile, error) {
	userID, err := singleTrustedHeader(headers, authHeaders.UserID, true)
	if err != nil {
		return OperatorProfile{}, err
	}
	profile := OperatorProfile{UserID: userID}
	for _, field := range []struct {
		header string
		target *string
	}{
		{authHeaders.Username, &profile.Username},
		{authHeaders.DisplayName, &profile.DisplayName},
		{authHeaders.Email, &profile.Email},
		{authHeaders.Groups, &profile.Groups},
		{authHeaders.Roles, &profile.Roles},
		{authHeaders.Scopes, &profile.Scopes},
	} {
		if field.header == "" {
			continue
		}
		*field.target, err = singleTrustedHeader(headers, field.header, false)
		if err != nil {
			return OperatorProfile{}, err
		}
	}
	return profile, nil
}

func singleTrustedHeader(headers http.Header, name string, required bool) (string, error) {
	values := headers.Values(name)
	if len(values) == 0 {
		if required {
			return "", fmt.Errorf("missing %s", name)
		}
		return "", nil
	}
	if len(values) != 1 {
		return "", fmt.Errorf("multiple %s values", name)
	}
	value := values[0]
	if !utf8.ValidString(value) || len(value) > 4096 || strings.ContainsAny(value, "\r\n") {
		return "", fmt.Errorf("malformed %s", name)
	}
	value = strings.TrimSpace(value)
	if required && value == "" {
		return "", fmt.Errorf("missing %s", name)
	}
	return value, nil
}
