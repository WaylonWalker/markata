package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/servecontrol"
	"github.com/WaylonWalker/markata/pkg/servefix"
	"github.com/fsnotify/fsnotify"
)

func TestBatchFixQueuesOneBuildWithFilesystemWatcher(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "post.md")
	original := "# Heading\n[link](//example.com)\n"
	if err := os.WriteFile(path, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatal(err)
	}
	if err := watcher.Add(root); err != nil {
		_ = watcher.Close()
		t.Fatal(err)
	}
	defer watcher.Close()
	clearExpectedServeFixWrites()
	defer clearExpectedServeFixWrites()
	isRebuilding.Store(false)
	rebuildPending.Store(false)
	rebuildCh := make(chan struct{}, 8)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	events := make(chan bool, 16)
	absPath, err := filepath.Abs(path)
	if err != nil {
		t.Fatal(err)
	}
	go watchFilesWithObserver(ctx, watcher, rebuildCh, func(event fsnotify.Event, suppressed bool) {
		eventPath, absErr := filepath.Abs(event.Name)
		if absErr == nil && filepath.Clean(eventPath) == filepath.Clean(absPath) {
			events <- suppressed
		}
	})
	var rebuildCount atomic.Int32
	built := make(chan struct{}, 8)
	go handleRebuildsWith(ctx, rebuildCh, 300*time.Millisecond, func() {
		rebuildCount.Add(1)
		built <- struct{}{}
	})

	runtime := servecontrol.NewRuntime()
	actionCount := 0
	runtime.SetActionHandler(func(request servecontrol.ActionRequest) error {
		if request.Kind != serveBuildAction {
			t.Fatalf("unexpected action: %+v", request)
		}
		actionCount++
		rebuildCh <- struct{}{}
		return nil
	})
	handler := servecontrol.NewWebHandlerWithFixHooks(runtime, root, func(previews []servefix.FilePreview) {
		registerExpectedServeFixWrites(root, previews)
	})
	plan, err := servefix.PlanFile(root, "post.md")
	if err != nil {
		t.Fatal(err)
	}
	selection := servefix.Selection{All: true, SafeOnly: true}
	previewBody, err := json.Marshal(map[string]any{
		"files": []servefix.FileSelection{{Path: "post.md", Selection: selection}},
	})
	if err != nil {
		t.Fatal(err)
	}
	previewRequest := httptest.NewRequest(http.MethodPost, "http://localhost/_markata/api/fixes/batch/preview", bytes.NewReader(previewBody))
	previewRequest.RemoteAddr = "127.0.0.1:45678"
	previewRequest.Header.Set("Content-Type", "application/json")
	previewRequest.Header.Set("Origin", "http://localhost")
	previewResponse := httptest.NewRecorder()
	handler.ServeHTTP(previewResponse, previewRequest)
	if previewResponse.Code != http.StatusOK {
		t.Fatalf("preview status=%d body=%s", previewResponse.Code, previewResponse.Body.String())
	}

	body, err := json.Marshal(map[string]any{
		"files": []servefix.FileSelection{{Path: "post.md", Digest: plan.Digest, Selection: selection}},
	})
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, "http://localhost/_markata/api/fixes/batch/apply", bytes.NewReader(body))
	request.RemoteAddr = "127.0.0.1:45678"
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Origin", "http://localhost")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("apply status=%d body=%s", response.Code, response.Body.String())
	}
	if actionCount != 1 {
		t.Fatalf("build action count=%d", actionCount)
	}
	select {
	case <-events:
	case <-time.After(3 * time.Second):
		t.Fatal("watcher did not observe the batch replacement")
	}
	waitForWatcherQuiet(t, events, 500*time.Millisecond)
	waitForBuildCount(t, built, &rebuildCount, 1)
	if got := rebuildCount.Load(); got != 1 {
		t.Fatalf("batch produced %d rebuilds, want exactly one", got)
	}
	updated, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(updated, []byte(original)) {
		t.Fatal("batch fix did not replace source")
	}

	// Exact undo is the important regression: watcher suppression remembers
	// only the replacement bytes Markata actually wrote, so a fast B -> A edit
	// is a real external change and must queue another rebuild.
	if err := os.WriteFile(path, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}
	waitForBuildCount(t, built, &rebuildCount, 2)

	if err := os.WriteFile(path, []byte(original+"external edit\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	waitForBuildCount(t, built, &rebuildCount, 3)
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	waitForBuildCount(t, built, &rebuildCount, 4)
}

func waitForBuildCount(t *testing.T, built <-chan struct{}, rebuildCount *atomic.Int32, want int32) {
	t.Helper()
	deadline := time.After(3 * time.Second)
	for rebuildCount.Load() < want {
		select {
		case <-built:
		case <-deadline:
			t.Fatalf("only %d rebuilds occurred, want %d", rebuildCount.Load(), want)
		}
	}
}

func waitForWatcherQuiet(t *testing.T, events <-chan bool, quietPeriod time.Duration) {
	t.Helper()
	timer := time.NewTimer(quietPeriod)
	defer timer.Stop()
	for {
		select {
		case <-events:
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			timer.Reset(quietPeriod)
		case <-timer.C:
			return
		}
	}
}
