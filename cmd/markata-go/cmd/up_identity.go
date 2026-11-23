package cmd

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	upIdentityPath     = "/_markata/up-identity"
	upIdentityProtocol = "markata-local-workspace-v1"
	upProbeTimeout     = 150 * time.Millisecond
)

type upServiceIdentity struct {
	Protocol  string `json:"protocol"`
	Workspace string `json:"workspace"`
	Service   string `json:"service"`
}

func upWorkspaceIdentity(workspace string) (string, error) {
	absolute, err := filepath.Abs(workspace)
	if err != nil {
		return "", err
	}
	// Resolve existing symlink spellings so a site opened from two aliases
	// still has a single local identity.
	if resolved, resolveErr := filepath.EvalSymlinks(absolute); resolveErr == nil {
		absolute = resolved
	}
	digest := sha256.Sum256([]byte(filepath.Clean(absolute)))
	return fmt.Sprintf("%x", digest), nil
}

// withUpIdentity exposes only a read-only local identity. It never exposes
// the workspace path, PID, source contents, credentials or admin operations.
// Its response is a *hint* for safe local reuse, not remote authentication.
func withUpIdentity(next http.Handler, service, workspace string) http.Handler {
	identity, err := upWorkspaceIdentity(workspace)
	if err != nil {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != upIdentityPath {
			next.ServeHTTP(w, r)
			return
		}
		if r.Method != http.MethodGet || !upIsLoopbackRequest(r) {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "private, no-store")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		if err := json.NewEncoder(w).Encode(upServiceIdentity{
			Protocol: upIdentityProtocol, Workspace: identity, Service: service,
		}); err != nil {
			return
		}
	})
}

func upIsLoopbackRequest(r *http.Request) bool {
	peer, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return false
	}
	ip := net.ParseIP(peer)
	if ip == nil || !ip.IsLoopback() {
		return false
	}
	host := r.Host
	if parsed, _, splitErr := net.SplitHostPort(host); splitErr == nil {
		host = parsed
	}
	hostIP := net.ParseIP(host)
	return strings.EqualFold(host, "localhost") || hostIP != nil && hostIP.IsLoopback()
}

// probeUpService only contacts a numeric IPv4 loopback address. Redirects
// aren't followed, and slow or non-Markata listeners are ignored.
func probeUpService(ctx context.Context, port int) (upServiceIdentity, bool) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"http://"+net.JoinHostPort(upLoopbackAddress, strconv.Itoa(port))+upIdentityPath, http.NoBody)
	if err != nil {
		return upServiceIdentity{}, false
	}
	transport := &http.Transport{
		Proxy:             nil,
		DialContext:       (&net.Dialer{Timeout: upProbeTimeout}).DialContext,
		DisableKeepAlives: true,
	}
	defer transport.CloseIdleConnections()
	client := &http.Client{
		Timeout:   upProbeTimeout,
		Transport: transport,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	response, err := client.Do(request)
	if err != nil {
		return upServiceIdentity{}, false
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK ||
		!strings.HasPrefix(response.Header.Get("Content-Type"), "application/json") {
		return upServiceIdentity{}, false
	}
	var result upServiceIdentity
	err = json.NewDecoder(io.LimitReader(response.Body, 1024)).Decode(&result)
	return result, err == nil && result.Protocol == upIdentityProtocol
}

// findExistingUpServices recognizes editor/preview servers in the bounded
// default port range. It does not touch or stop processes it finds. Explicit
// ports restrict discovery to that exact port.
func findExistingUpServices(ctx context.Context, services []upService, workspace string, policy upPortPolicy) (newServices, reused []upService, err error) {
	identity, err := upWorkspaceIdentity(workspace)
	if err != nil {
		return nil, nil, err
	}
	newServices = make([]upService, 0, len(services))
	for _, service := range services {
		base, portErr := upServicePort(service)
		if portErr != nil {
			return nil, nil, portErr
		}
		last := base
		if !policy.fixed(service) {
			last = min(base+upPortFallbackRange, 65535)
		}
		found := false
		for port := base; port <= last; port++ {
			if err := ctx.Err(); err != nil {
				return nil, nil, err
			}
			instance, ok := probeUpService(ctx, port)
			if !ok || instance.Workspace != identity || instance.Service != service.name {
				continue
			}
			reusable, updateErr := withUpServicePort(service, port)
			if updateErr != nil {
				return nil, nil, updateErr
			}
			reused = append(reused, reusable)
			found = true
			break
		}
		if !found {
			newServices = append(newServices, service)
		}
	}
	return newServices, reused, nil
}
