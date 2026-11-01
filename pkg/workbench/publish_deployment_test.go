package workbench

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPublishDeploymentNeedsPositiveLiveEvidence(t *testing.T) {
	expected := strings.Repeat("a", 40)
	var evidence deployEvidence
	status := evaluateDeployment(expected, evidence)
	if status.Verified || status.Status != "unknown" {
		t.Fatalf("missing evidence must not verify publication: %+v", status)
	}
	evidence.Revision.LiveSourceCommit = strings.Repeat("b", 40)
	status = evaluateDeployment(expected, evidence)
	if status.Verified || status.Status != "not-live" {
		t.Fatalf("stale live source must not verify publication: %+v", status)
	}
	evidence.Revision.LiveSourceCommit = expected
	status = evaluateDeployment(expected, evidence)
	if !status.Verified || status.Status != "live-verified" {
		t.Fatalf("matching live source must verify publication: %+v", status)
	}
}

func TestPublishDeploymentUsesCurrentReleaseFallback(t *testing.T) {
	expected := strings.Repeat("a", 40)
	var evidence deployEvidence
	evidence.CurrentID = "release-2"
	evidence.Releases = append(evidence.Releases,
		struct {
			ID           string `json:"id"`
			SourceCommit string `json:"source_commit"`
			Current      bool   `json:"current"`
		}{ID: "release-2", SourceCommit: expected, Current: true})
	status := evaluateDeployment(expected, evidence)
	if !status.Verified || status.ReleaseID != "release-2" {
		t.Fatalf("expected live release fallback: %+v", status)
	}
}

func TestPublishBuilderURLDoesNotAllowInsecureRemoteAuth(t *testing.T) {
	for _, bad := range []string{
		"http://builder.example.com",
		"https://user:pass@builder.example.com",
		"https://builder.example.com/some/path",
		"https://builder.example.com/?token=bad",
	} {
		if _, err := publishBuilderURL(bad); err == nil {
			t.Errorf("must reject %q", bad)
		}
	}
	endpoint, err := publishBuilderURL("https://builder.example.com")
	if err != nil || endpoint != "https://builder.example.com/api/v1/state" {
		t.Fatalf("unexpected builder URL %q: %v", endpoint, err)
	}
	if _, err := publishBuilderURL("http://127.0.0.1:9090"); err != nil {
		t.Fatalf("local development should support loopback: %v", err)
	}
}

func TestPublishDeploymentWithoutConfiguredBuilderIsExplicit(t *testing.T) {
	t.Setenv("MARKATA_PUBLISH_BUILDER_URL", "")
	t.Setenv("MARKATA_PUBLISH_BUILDER_TOKEN", "")
	handler := &Handler{}
	request := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/workspace/publish/deployment?commit="+strings.Repeat("a", 40), http.NoBody)
	recorder := httptest.NewRecorder()
	handler.servePublishDeployment(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status: %d", recorder.Code)
	}
	var result deploymentStatus
	if err := json.Unmarshal(recorder.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Verified || result.Configured || result.Status != "unconfigured" {
		t.Fatalf("must not claim publication without builder configuration: %+v", result)
	}
}
