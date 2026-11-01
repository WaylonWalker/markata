package workbench

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type deployEvidence struct {
	Revision struct {
		LiveSourceCommit string `json:"live_source_commit"`
	} `json:"revision"`
	Releases []struct {
		ID           string `json:"id"`
		SourceCommit string `json:"source_commit"`
		Current      bool   `json:"current"`
	} `json:"releases"`
	CurrentID string `json:"current_release_id"`
}

type deploymentStatus struct {
	Configured       bool   `json:"configured"`
	Status           string `json:"status"`
	Verified         bool   `json:"verified"`
	ExpectedCommit   string `json:"expected_commit"`
	LiveSourceCommit string `json:"live_source_commit,omitempty"`
	ReleaseID        string `json:"release_id,omitempty"`
	Message          string `json:"message"`
}

const (
	secureScheme   = "https"
	insecureScheme = "http"
)

func evaluateDeployment(commit string, evidence deployEvidence) deploymentStatus {
	status := deploymentStatus{
		Configured: true, ExpectedCommit: commit, Status: "unknown",
		Message: "Builder reports no reliable live source revision yet.",
	}
	status.LiveSourceCommit = strings.ToLower(strings.TrimSpace(evidence.Revision.LiveSourceCommit))
	for _, release := range evidence.Releases {
		if !release.Current || (evidence.CurrentID != "" && release.ID != evidence.CurrentID) {
			continue
		}
		status.ReleaseID = release.ID
		if status.LiveSourceCommit == "" {
			status.LiveSourceCommit = strings.ToLower(strings.TrimSpace(release.SourceCommit))
		}
		break
	}
	if !gitHashPattern.MatchString(status.LiveSourceCommit) {
		status.LiveSourceCommit = ""
		return status
	}
	if status.LiveSourceCommit == commit {
		status.Status = "live-verified"
		status.Verified = true
		status.Message = "The current Builder Admin release reports the requested source commit."
	} else {
		status.Status = "not-live"
		status.Message = "The current release has a different source commit. Your build may still be queued, running or unsuccessful."
	}
	return status
}

func publishBuilderURL(raw string) (string, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return "", fmt.Errorf("invalid builder endpoint")
	}
	host := u.Hostname()
	if u.Scheme != secureScheme {
		ip := net.ParseIP(host)
		loopback := strings.EqualFold(host, "localhost") || (ip != nil && ip.IsLoopback())
		if u.Scheme != insecureScheme || !loopback {
			return "", fmt.Errorf("builder endpoint must use HTTPS")
		}
	}
	if u.Path != "" && u.Path != "/" {
		return "", fmt.Errorf("builder endpoint must be an origin with no path")
	}
	return strings.TrimRight(u.String(), "/") + "/api/v1/state", nil
}

func (h *Handler) servePublishDeployment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		fail(w, http.StatusMethodNotAllowed, "GET required")
		return
	}
	commit := strings.ToLower(r.URL.Query().Get("commit"))
	if !gitHashPattern.MatchString(commit) {
		fail(w, http.StatusBadRequest, "expected commit SHA is required")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	base := strings.TrimSpace(os.Getenv("MARKATA_PUBLISH_BUILDER_URL"))
	token := strings.TrimSpace(os.Getenv("MARKATA_PUBLISH_BUILDER_TOKEN"))
	if base == "" || token == "" {
		if err := json.NewEncoder(w).Encode(deploymentStatus{
			Status: "unconfigured", ExpectedCommit: commit,
			Message: "Deployment verification is not configured for this editor.",
		}); err != nil {
			return
		}
		return
	}
	endpoint, err := publishBuilderURL(base)
	if err != nil {
		fail(w, 422, "invalid builder deployment configuration")
		return
	}
	// #nosec G704 -- endpoint is an explicit operator setting validated as an HTTPS origin (loopback HTTP for local use).
	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		fail(w, 422, "cannot prepare builder status request")
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	client := &http.Client{Timeout: 7 * time.Second, CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	}}
	// #nosec G704 -- redirects are disabled and the destination is the operator-configured builder origin.
	resp, err := client.Do(req)
	if err != nil {
		fail(w, http.StatusBadGateway, "builder status is unreachable")
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fail(w, http.StatusBadGateway, "builder status is unavailable or not authorized")
		return
	}
	var evidence deployEvidence
	if err := json.NewDecoder(http.MaxBytesReader(w, resp.Body, 1<<20)).Decode(&evidence); err != nil {
		fail(w, http.StatusBadGateway, "invalid builder state response")
		return
	}
	if err := json.NewEncoder(w).Encode(evaluateDeployment(commit, evidence)); err != nil {
		return
	}
}
