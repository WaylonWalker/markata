// Package editorhub provides site-scoped, read-only hosted-editor discovery.
// It does not expose the trusted-local workbench or builder operator API.
package editorhub

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

var siteIDPattern = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]{0,62}[a-z0-9])?$`)

// Grant represents explicit, per-subject site permissions. These are not
// inferred from authentication alone or from any headers sent by a browser.
type Grant struct {
	Read    bool `json:"read"`
	Write   bool `json:"write"`
	Publish bool `json:"publish"`
}

// Site is a trusted server-side configuration record. Grants map stable
// identity-provider subject IDs to capabilities; no wildcard grants are allowed.
type Site struct {
	ID     string
	Name   string
	Grants map[string]Grant
}

// SiteSummary is the complete, public shape of a site-picker entry. It
// intentionally contains no storage path, Git remote, or builder endpoint.
type SiteSummary struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Permissions Grant  `json:"permissions"`
}

// Registry holds a validated, immutable copy of the configured sites.
type Registry struct {
	sites []Site
}

// NewRegistry validates all ACLs at startup. Publishing requires write access;
// writing requires read access. Invalid or ambiguous configuration fails closed.
func NewRegistry(sites []Site) (*Registry, error) {
	registry := &Registry{sites: make([]Site, 0, len(sites))}
	seen := make(map[string]struct{}, len(sites))
	for _, site := range sites {
		if !siteIDPattern.MatchString(site.ID) || len(site.ID) > 64 {
			return nil, fmt.Errorf("invalid site id %q", site.ID)
		}
		if _, exists := seen[site.ID]; exists {
			return nil, fmt.Errorf("duplicate site id %q", site.ID)
		}
		seen[site.ID] = struct{}{}
		if name := strings.TrimSpace(site.Name); name == "" || len(name) > 128 || strings.ContainsAny(name, "\r\n\x00") {
			return nil, fmt.Errorf("invalid display name for site %q", site.ID)
		}
		copySite := Site{ID: site.ID, Name: site.Name, Grants: make(map[string]Grant, len(site.Grants))}
		for subject, grant := range site.Grants {
			if subject == "" || len(subject) > 256 || strings.TrimSpace(subject) != subject || strings.ContainsAny(subject, "\r\n\x00") || subject == "*" {
				return nil, fmt.Errorf("invalid subject for site %q", site.ID)
			}
			if (grant.Write && !grant.Read) || (grant.Publish && !grant.Write) {
				return nil, fmt.Errorf("inconsistent grants for site %q", site.ID)
			}
			copySite.Grants[subject] = grant
		}
		registry.sites = append(registry.sites, copySite)
	}
	sort.Slice(registry.sites, func(i, j int) bool { return registry.sites[i].ID < registry.sites[j].ID })
	return registry, nil
}

// Permission identifies an action on a site. Document-level visibility must
// still be checked separately before serving private drafts or media.
type Permission string

const (
	PermissionRead    Permission = "read"
	PermissionWrite   Permission = "write"
	PermissionPublish Permission = "publish"
)

// Allows evaluates a site-scoped capability. Unknown permissions and subjects
// fail closed; this never grants access to private documents by itself.
func (r *Registry) Allows(subject, siteID string, permission Permission) bool {
	if r == nil || subject == "" || siteID == "" {
		return false
	}
	i := sort.Search(len(r.sites), func(i int) bool { return r.sites[i].ID >= siteID })
	if i == len(r.sites) || r.sites[i].ID != siteID {
		return false
	}
	grant, ok := r.sites[i].Grants[subject]
	if !ok {
		return false
	}
	switch permission {
	case PermissionRead:
		return grant.Read
	case PermissionWrite:
		return grant.Write
	case PermissionPublish:
		return grant.Publish
	default:
		return false
	}
}

// ListFor returns only explicitly readable sites. It does not reveal that
// hidden sites exist or how many private sites are configured.
func (r *Registry) ListFor(subject string) []SiteSummary {
	visible := make([]SiteSummary, 0)
	if subject == "" || r == nil {
		return visible
	}
	for _, site := range r.sites {
		if grant, ok := site.Grants[subject]; ok && grant.Read {
			visible = append(visible, SiteSummary{ID: site.ID, Name: site.Name, Permissions: grant})
		}
	}
	return visible
}

// Authenticator must return a stable subject verified by a trusted auth layer.
// Never implement it by reading unvalidated X-Hlab-* or forwarded headers.
type Authenticator func(*http.Request) (string, error)

// Handler exposes only GET /api/v1/sites. The hosting application must supply
// an authenticator that validates ingress identity before invoking this API.
type Handler struct {
	registry     *Registry
	authenticate Authenticator
}

func NewHandler(registry *Registry, authenticate Authenticator) (*Handler, error) {
	if registry == nil || authenticate == nil {
		return nil, errors.New("site registry and trusted authenticator are required")
	}
	return &Handler{registry: registry, authenticate: authenticate}, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "private, no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Vary", "Cookie, Authorization")
	if r.URL.Path != "/api/v1/sites" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Sec-Fetch-Site") == "cross-site" {
		http.Error(w, "cross-site access denied", http.StatusForbidden)
		return
	}
	subject, err := h.authenticate(r)
	if err != nil || subject == "" || len(subject) > 256 || strings.TrimSpace(subject) != subject || strings.ContainsAny(subject, "\r\n\x00") {
		http.Error(w, "authentication required", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(struct {
		Sites []SiteSummary `json:"sites"`
	}{Sites: h.registry.ListFor(subject)}); err != nil {
		return // Response may already have been partially written.
	}
}
