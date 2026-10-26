package cmd

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/fsnotify/fsnotify"
)

func TestAddDirRecursive(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir")
	nestedDir := filepath.Join(subDir, "nested")

	if err := os.MkdirAll(nestedDir, 0o755); err != nil {
		t.Fatalf("Failed to create test directories: %v", err)
	}

	// Create a watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatalf("Failed to create watcher: %v", err)
	}
	defer watcher.Close()

	// Clear the global output path to avoid filtering
	origOutputPath := serveOutputPath
	serveOutputPath = ""
	defer func() { serveOutputPath = origOutputPath }()

	// Add the directory recursively
	if err := addDirRecursive(watcher, tmpDir); err != nil {
		t.Fatalf("addDirRecursive failed: %v", err)
	}

	// Verify all directories are being watched
	// We can't directly check the watch list, but we can verify
	// events are received for files in nested directories
	testFile := filepath.Join(nestedDir, "test.md")

	// Write a file and check for events
	done := make(chan bool, 1)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Create != 0 && filepath.Base(event.Name) == "test.md" {
					done <- true
					return
				}
			case <-time.After(2 * time.Second):
				done <- false
				return
			}
		}
	}()

	// Create the test file
	if err := os.WriteFile(testFile, []byte("# Test"), 0o600); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Wait for the event
	select {
	case received := <-done:
		if !received {
			t.Error("Did not receive create event for nested file")
		}
	case <-time.After(3 * time.Second):
		t.Error("Timeout waiting for file event")
	}
}

func TestAddDirRecursive_SkipsHiddenDirs(t *testing.T) {
	// Skip on Windows as hidden dirs are handled via file attributes, not dot prefix
	if filepath.Separator == '\\' {
		t.Skip("Skipping on Windows: hidden directories use file attributes, not dot prefix")
	}

	tmpDir := t.TempDir()
	hiddenDir := filepath.Join(tmpDir, ".hidden")
	fileInHidden := filepath.Join(hiddenDir, "file.md")

	if err := os.MkdirAll(hiddenDir, 0o755); err != nil {
		t.Fatalf("Failed to create hidden directory: %v", err)
	}
	if err := os.WriteFile(fileInHidden, []byte("test"), 0o600); err != nil {
		t.Fatalf("Failed to create file in hidden dir: %v", err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatalf("Failed to create watcher: %v", err)
	}
	defer watcher.Close()

	// Clear the global output path
	origOutputPath := serveOutputPath
	serveOutputPath = ""
	defer func() { serveOutputPath = origOutputPath }()

	if err := addDirRecursive(watcher, tmpDir); err != nil {
		t.Fatalf("addDirRecursive failed: %v", err)
	}

	// Verify hidden directory is NOT watched by creating a file
	// and checking that no event is received
	newFile := filepath.Join(hiddenDir, "new.md")

	eventReceived := make(chan bool, 1)
	timeout := time.After(500 * time.Millisecond)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if filepath.Base(event.Name) == "new.md" {
					eventReceived <- true
					return
				}
			case <-timeout:
				eventReceived <- false
				return
			}
		}
	}()

	if err := os.WriteFile(newFile, []byte("new"), 0o600); err != nil {
		t.Fatalf("Failed to write new file: %v", err)
	}

	select {
	case received := <-eventReceived:
		if received {
			t.Error("Should not receive events for files in hidden directories")
		}
	case <-time.After(1 * time.Second):
		// No event received, which is expected
	}
}

func TestAddDirRecursive_SkipsOutputDir(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "output")
	contentDir := filepath.Join(tmpDir, "content")

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}
	if err := os.MkdirAll(contentDir, 0o755); err != nil {
		t.Fatalf("Failed to create content directory: %v", err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatalf("Failed to create watcher: %v", err)
	}
	defer watcher.Close()

	// Set the output path to skip
	origOutputPath := serveOutputPath
	serveOutputPath = outputDir
	defer func() { serveOutputPath = origOutputPath }()

	if err := addDirRecursive(watcher, tmpDir); err != nil {
		t.Fatalf("addDirRecursive failed: %v", err)
	}

	// Create a file in output dir - should NOT trigger event
	outputFile := filepath.Join(outputDir, "index.html")

	eventReceived := make(chan bool, 1)
	go func() {
		select {
		case event := <-watcher.Events:
			if filepath.Base(event.Name) == "index.html" {
				eventReceived <- true
			}
		case <-time.After(500 * time.Millisecond):
			eventReceived <- false
		}
	}()

	if err := os.WriteFile(outputFile, []byte("<html></html>"), 0o600); err != nil {
		t.Fatalf("Failed to write output file: %v", err)
	}

	received := <-eventReceived
	if received {
		t.Error("Should not receive events for files in output directory")
	}
}

func TestServeSearchPostsSnapshotCopiesMutablePostState(t *testing.T) {
	post := models.NewPost("posts/example.md")
	post.Path = "posts/example.md"
	title := "Before"
	post.Title = &title
	post.Tags = []string{"before"}
	post.Extra["image"] = "https://example.test/before.jpg"
	setServeSearchPosts([]*models.Post{post})
	t.Cleanup(func() { setServeSearchPosts(nil) })

	snapshot := getServeSearchPostsSnapshot()
	if len(snapshot) != 1 || snapshot[0] == post {
		t.Fatalf("snapshot = %#v, want an independent post", snapshot)
	}

	*post.Title = "After"
	post.Tags[0] = "after"
	post.Extra["image"] = "https://example.test/after.jpg"

	if got := snapshot[0].PlainTitle(); got != "Before" {
		t.Fatalf("snapshot title = %q, want Before", got)
	}
	if got := snapshot[0].Tags[0]; got != "before" {
		t.Fatalf("snapshot tag = %q, want before", got)
	}
	if got := snapshot[0].Extra["image"]; got != "https://example.test/before.jpg" {
		t.Fatalf("snapshot image = %v, want before URL", got)
	}
}

func TestIsPathWithinDir(t *testing.T) {
	tests := []struct {
		name     string
		pathname string
		dir      string
		want     bool
	}{
		{
			name:     "file within dir",
			pathname: "/home/user/project/file.txt",
			dir:      "/home/user/project",
			want:     true,
		},
		{
			name:     "nested file within dir",
			pathname: "/home/user/project/sub/file.txt",
			dir:      "/home/user/project",
			want:     true,
		},
		{
			name:     "file outside dir",
			pathname: "/home/user/other/file.txt",
			dir:      "/home/user/project",
			want:     false,
		},
		{
			name:     "path traversal attempt",
			pathname: "/home/user/project/../other/file.txt",
			dir:      "/home/user/project",
			want:     false,
		},
		{
			name:     "same path",
			pathname: "/home/user/project",
			dir:      "/home/user/project",
			want:     true,
		},
		{
			name:     "empty dir",
			pathname: "/home/user/file.txt",
			dir:      "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isPathWithinDir(tt.pathname, tt.dir)
			if got != tt.want {
				t.Errorf("isPathWithinDir(%q, %q) = %v, want %v",
					tt.pathname, tt.dir, got, tt.want)
			}
		})
	}
}

func TestInjectDevScripts_AddsBannerScript(t *testing.T) {
	html := "<html><body><h1>Hi</h1></body></html>"
	status := BuildStatus{Status: buildStatusBuilding}

	updated := injectDevScripts(html, status)

	if !strings.Contains(updated, "markata-build-banner") {
		t.Error("expected build banner script injection")
	}
	if !strings.Contains(updated, "/__livereload") {
		t.Error("expected live reload EventSource")
	}
	if !strings.Contains(updated, `"status":"building"`) {
		t.Error("expected build status payload")
	}
}

func TestServe404Page_FallbackIncludesBanner(t *testing.T) {
	outputDir := t.TempDir()
	recorder := httptest.NewRecorder()

	serve404Page(recorder, outputDir, BuildStatus{Status: buildStatusBuilding})

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected 404 status, got %d", recorder.Code)
	}

	body := recorder.Body.String()
	if !strings.Contains(body, "Page Not Found") {
		t.Error("expected fallback 404 content")
	}
	if !strings.Contains(body, "markata-build-banner") {
		t.Error("expected build banner in fallback 404")
	}
}

func TestBuildStatusPayload_IncludesLicenseWarning(t *testing.T) {
	payload := buildStatusPayload(BuildStatus{
		Status:         buildStatusSuccess,
		LicenseWarning: "set license in config",
	})

	if !strings.Contains(payload, "license_warning") {
		t.Fatalf("expected license_warning field in payload, got %s", payload)
	}
}

func TestBuildDevScript_IncludesLicenseToast(t *testing.T) {
	script := buildDevScript(BuildStatus{Status: buildStatusSuccess})

	if !strings.Contains(script, "markata-license-toast") {
		t.Fatalf("expected license toast id in dev script")
	}
	if !strings.Contains(script, "license_warning") {
		t.Fatalf("expected license_warning handling in dev script")
	}
}

func TestInjectDevScripts_UsesConfiguredSearchEndpoint(t *testing.T) {
	original := serveSearchEndpoint
	serveSearchEndpoint = "https://search.example.com/api/search"
	defer func() { serveSearchEndpoint = original }()

	updated := injectDevScripts("<html><head></head><body></body></html>", BuildStatus{Status: buildStatusBuilding})
	if !strings.Contains(updated, `window.__markataSearchEndpoint = "https://search.example.com/api/search";`) {
		t.Fatalf("expected configured search endpoint injection, got %s", updated)
	}
}

func TestHandleSearchFallback_EncodesRedirectQuery(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/_search", strings.NewReader("q=go+%26+bleve"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	recorder := httptest.NewRecorder()

	handleSearchFallback(recorder, req, t.TempDir())

	if recorder.Code != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d", recorder.Code)
	}
	if location := recorder.Header().Get("Location"); location != "/?q=go+%26+bleve" {
		t.Fatalf("location = %q, want %q", location, "/?q=go+%26+bleve")
	}
}
