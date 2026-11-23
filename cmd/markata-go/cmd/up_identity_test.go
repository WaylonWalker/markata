package cmd

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func startUpIdentityTestServer(t *testing.T, service, workspace string) int {
	t.Helper()
	server := httptest.NewServer(withUpIdentity(http.NotFoundHandler(), service, workspace))
	t.Cleanup(server.Close)
	parsed, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	port, err := strconv.Atoi(parsed.Port())
	if err != nil {
		t.Fatal(err)
	}
	return port
}

func TestUpIdentityIsLocalAndDoesNotExposeWorkspacePath(t *testing.T) {
	workspace := t.TempDir()
	identity, err := upWorkspaceIdentity(workspace)
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(withUpIdentity(http.NotFoundHandler(), upServiceEditor, workspace))
	defer server.Close()
	response, err := http.Get(server.URL + upIdentityPath) //nolint:gosec // Local test listener
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	var value upServiceIdentity
	if err := json.NewDecoder(response.Body).Decode(&value); err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusOK || value.Workspace != identity || value.Service != upServiceEditor || value.Protocol != upIdentityProtocol {
		t.Fatalf("unexpected identity: %#v status=%d", value, response.StatusCode)
	}
	if strings.Contains(value.Workspace, filepath.Base(workspace)) {
		t.Fatal("identity response must not expose local directory name")
	}
	remote := httptest.NewRequest(http.MethodGet, "http://other.invalid"+upIdentityPath, nil)
	remote.RemoteAddr = "127.0.0.1:48888"
	writer := httptest.NewRecorder()
	withUpIdentity(http.NotFoundHandler(), upServiceEditor, workspace).ServeHTTP(writer, remote)
	if writer.Code != http.StatusNotFound {
		t.Fatalf("non-loopback Host must not access identity: %d", writer.Code)
	}
	remote.Host = "127.0.0.1:8766"
	remote.RemoteAddr = "192.0.2.4:48888"
	writer = httptest.NewRecorder()
	withUpIdentity(http.NotFoundHandler(), upServiceEditor, workspace).ServeHTTP(writer, remote)
	if writer.Code != http.StatusNotFound {
		t.Fatalf("non-loopback peer must not access identity: %d", writer.Code)
	}
}

func TestUpDiscoversExistingEditorWithoutStartingDuplicate(t *testing.T) {
	workspace := t.TempDir()
	port := startUpIdentityTestServer(t, upServiceEditor, workspace)
	services, err := planUpServices(upOptions{editorOnly: true, editorPort: port})
	if err != nil {
		t.Fatal(err)
	}
	launch, reused, err := findExistingUpServices(context.Background(), services, workspace, upPortPolicy{})
	if err != nil {
		t.Fatal(err)
	}
	if len(launch) != 0 || len(reused) != 1 || reused[0].url != services[0].url {
		t.Fatalf("same-workspace editor should be reused, launch=%v reused=%v", launch, reused)
	}
}

func TestUpDiscoversBothEditorAndPreview(t *testing.T) {
	workspace := t.TempDir()
	editorPort := startUpIdentityTestServer(t, upServiceEditor, workspace)
	previewPort := startUpIdentityTestServer(t, upServiceBuilder, workspace)
	services, err := planUpServices(upOptions{editorPort: editorPort, previewPort: previewPort})
	if err != nil {
		t.Fatal(err)
	}
	launch, reused, err := findExistingUpServices(context.Background(), services, workspace, upPortPolicy{})
	if err != nil {
		t.Fatal(err)
	}
	if len(launch) != 0 || len(reused) != 2 {
		t.Fatalf("expected both services reused: launch=%v reused=%v", launch, reused)
	}
}

func TestUpRecognizesOnlyMatchingWorkspaceAndRole(t *testing.T) {
	workspace := t.TempDir()
	other := t.TempDir()
	port := startUpIdentityTestServer(t, upServiceEditor, other)
	services, err := planUpServices(upOptions{editorOnly: true, editorPort: port})
	if err != nil {
		t.Fatal(err)
	}
	launch, reused, err := findExistingUpServices(context.Background(), services, workspace, upPortPolicy{})
	if err != nil {
		t.Fatal(err)
	}
	if len(reused) != 0 || len(launch) != 1 {
		t.Fatalf("must not attach another workspace: launch=%v reused=%v", launch, reused)
	}
	port = startUpIdentityTestServer(t, upServiceBuilder, workspace)
	services, err = planUpServices(upOptions{editorOnly: true, editorPort: port})
	if err != nil {
		t.Fatal(err)
	}
	launch, reused, err = findExistingUpServices(context.Background(), services, workspace, upPortPolicy{})
	if err != nil || len(reused) != 0 || len(launch) != 1 {
		t.Fatalf("must not attach wrong role: launch=%v reused=%v err=%v", launch, reused, err)
	}
}

func TestUpRecognizesWorkspaceOnAdjacentPort(t *testing.T) {
	workspace := t.TempDir()
	port := startUpIdentityTestServer(t, upServiceEditor, workspace)
	if port <= 1 {
		t.Skip("no preceding port to use")
	}
	services, err := planUpServices(upOptions{editorOnly: true, editorPort: port - 1})
	if err != nil {
		t.Fatal(err)
	}
	launch, reused, err := findExistingUpServices(context.Background(), services, workspace, upPortPolicy{})
	if err != nil {
		t.Fatal(err)
	}
	if len(reused) != 1 || len(launch) != 0 {
		t.Fatalf("should recognize existing service in fallback range: launch=%v reused=%v", launch, reused)
	}
	if got, err := upServicePort(reused[0]); err != nil || got != port {
		t.Fatalf("wrong reused editor port: %d, %v", got, err)
	}
}

func TestUpDoesNotTrustUnidentifiedLocalWebServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == upIdentityPath {
			http.Redirect(w, r, "https://example.com/secret", http.StatusFound)
			return
		}
		_, _ = io.WriteString(w, "<html>unrelated app</html>")
	}))
	defer server.Close()
	parsed, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	port, err := strconv.Atoi(parsed.Port())
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := probeUpService(context.Background(), port); ok {
		t.Fatal("redirect or unrelated service must never be considered Markata")
	}
}

func TestUpIdentityOnlyUsesLoopback(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:8766/", nil)
	r.RemoteAddr = net.JoinHostPort("127.0.0.1", "42345")
	if !upIsLoopbackRequest(r) {
		t.Fatal("expected loopback request")
	}
	r.RemoteAddr = "10.0.0.4:42345"
	if upIsLoopbackRequest(r) {
		t.Fatal("remote network must not qualify")
	}
}
