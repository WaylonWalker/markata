package builderadmin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/diagnostics"
	"github.com/WaylonWalker/markata/pkg/servecontrol"
)

func TestProjectControlState_PreservesBuilderJobs(t *testing.T) {
	now := time.Date(2026, 9, 26, 12, 0, 0, 0, time.UTC)
	state := State{
		Queue: []QueuedOperation{
			{ID: "build-q", Kind: "build", Label: "Build", TriggerType: "file-watch", EnqueuedAt: now, Changed: []string{"posts/a.md"}},
			{ID: "refresh-r", Kind: "refresh", Label: "Refresh reader", TriggerType: "schedule", EnqueuedAt: now.Add(time.Minute)},
		},
		Running: &RunningOperation{ID: "refresh-r", Kind: "refresh", Label: "Refresh reader", TriggerType: "schedule", EnqueuedAt: now.Add(time.Minute), StartedAt: now.Add(2 * time.Minute), Phase: "refresh"},
		Builds: []BuildRecord{
			{ID: "build-done", Kind: "build", Status: "failed", TriggerType: "manual-ui", EnqueuedAt: now.Add(-time.Hour), StartedAt: now.Add(-59 * time.Minute), FinishedAt: now.Add(-58 * time.Minute), PrepareMS: 100, BuildMS: 200, Error: "template failed", ChangedPaths: []string{"posts/b.md"}},
		},
		Refresh: []RefreshRecord{
			{ID: "refresh-done", TaskName: "reader", Status: "success", TriggerType: "schedule", EnqueuedAt: now.Add(-2 * time.Hour), StartedAt: now.Add(-119 * time.Minute), FinishedAt: now.Add(-118 * time.Minute), RunMS: 500},
		},
	}
	snapshot := projectControlState(state)
	if len(snapshot.Jobs) != 4 {
		t.Fatalf("jobs = %d, want 4 (running queue item must not duplicate)", len(snapshot.Jobs))
	}
	jobs := make(map[string]servecontrol.Job, len(snapshot.Jobs))
	for _, job := range snapshot.Jobs {
		jobs[job.ID] = job
	}
	if jobs["build-q"].State != servecontrol.StateQueued || !jobs["build-q"].QueuedAt.Equal(now) || len(jobs["build-q"].Pages) != 1 {
		t.Errorf("queued build = %+v", jobs["build-q"])
	}
	if jobs["refresh-r"].State != servecontrol.StateRunning || !jobs["refresh-r"].QueuedAt.Equal(now.Add(time.Minute)) || !jobs["refresh-r"].StartedAt.Equal(now.Add(2*time.Minute)) || len(jobs["refresh-r"].Steps) != 1 {
		t.Errorf("running refresh = %+v", jobs["refresh-r"])
	}
	if jobs["build-done"].State != servecontrol.StateFailed || !jobs["build-done"].EndedAt.Equal(now.Add(-58*time.Minute)) || len(jobs["build-done"].Steps) != 2 || len(jobs["build-done"].Diagnostics) != 1 {
		t.Errorf("completed build = %+v", jobs["build-done"])
	}
	if jobs["refresh-done"].State != servecontrol.StateSuccess || len(jobs["refresh-done"].Steps) != 1 {
		t.Errorf("completed refresh = %+v", jobs["refresh-done"])
	}
	if len(snapshot.Diagnostics) != 1 || snapshot.Diagnostics[0].JobID != "build-done" {
		t.Errorf("completed error diagnostic = %+v", snapshot.Diagnostics)
	}
}

func TestBuilderAdminStateAPI_ContainsSharedControlRuntime(t *testing.T) {
	dir := t.TempDir()
	svc, err := New(Config{SourceDir: dir, SiteDir: filepath.Join(dir, "site"), HistoryDir: filepath.Join(dir, "history")})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = svc.leaderLock.Close() })
	now := time.Now().UTC()
	svc.stateMu.Lock()
	svc.state.Queue = append(svc.state.Queue, QueuedOperation{ID: "build-123", Kind: "build", Label: "Build", TriggerType: "manual-ui", EnqueuedAt: now})
	svc.saveStateLocked()
	svc.stateMu.Unlock()
	runtimeSnapshot := svc.controlRuntime.Snapshot()
	if len(runtimeSnapshot.Jobs) != 1 || runtimeSnapshot.Jobs[0].ID != "build-123" {
		t.Fatalf("runtime snapshot = %+v", runtimeSnapshot.Jobs)
	}
	response := httptest.NewRecorder()
	svc.handleState(response, httptest.NewRequest(http.MethodGet, "/api/state", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("state status = %d", response.Code)
	}
	var payload struct {
		State         State                 `json:"state"`
		ControlCenter servecontrol.Snapshot `json:"control_center"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if len(payload.State.Queue) != 1 || len(payload.ControlCenter.Jobs) != 1 || payload.ControlCenter.Jobs[0].ID != payload.State.Queue[0].ID {
		t.Fatalf("API projections diverged: legacy=%+v shared=%+v", payload.State.Queue, payload.ControlCenter.Jobs)
	}
}

func TestBuilderAdminStateAPI_UsesCachedControlProjection(t *testing.T) {
	dir := t.TempDir()
	siteDir := filepath.Join(dir, "site")
	releaseDir := filepath.Join(siteDir, "releases", "release-1")
	artifactPath := filepath.Join(releaseDir, diagnostics.DefaultArtifactPath)
	if err := os.MkdirAll(filepath.Dir(artifactPath), 0o755); err != nil {
		t.Fatal(err)
	}
	writeArtifact := func(page string) {
		t.Helper()
		data, err := diagnostics.MarshalArtifact(diagnostics.ContentLedgerSnapshot{
			Entries: []diagnostics.ContentDisposition{{Path: page, Emitted: true}},
		}, diagnostics.ArtifactBuildInfo{BuiltAt: time.Now().UTC()})
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(artifactPath, data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	writeArtifact("posts/first.md")
	if err := os.Symlink(filepath.Join("releases", "release-1"), filepath.Join(siteDir, "current")); err != nil {
		t.Fatal(err)
	}
	svc, err := New(Config{SourceDir: dir, SiteDir: siteDir, HistoryDir: filepath.Join(dir, "history")})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = svc.leaderLock.Close() })

	// Change the artifact without changing Builder Admin state. Polling must
	// return the cached projection instead of rereading and reparsing the file.
	writeArtifact("posts/second.md")
	for range 2 {
		response := httptest.NewRecorder()
		svc.handleState(response, httptest.NewRequest(http.MethodGet, "/api/state", nil))
		var payload struct {
			ControlCenter servecontrol.Snapshot `json:"control_center"`
		}
		if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
			t.Fatal(err)
		}
		if len(payload.ControlCenter.Pages) != 1 || payload.ControlCenter.Pages[0].Path != "posts/first.md" {
			t.Fatalf("poll refreshed files unexpectedly: pages = %+v", payload.ControlCenter.Pages)
		}
	}

	// A state transition refreshes the projection and incorporates the new
	// artifact contents for subsequent polls.
	svc.stateMu.Lock()
	svc.state.Queue = append(svc.state.Queue, QueuedOperation{ID: "build-next", Kind: "build", Label: "Build", EnqueuedAt: time.Now().UTC()})
	svc.saveStateLocked()
	svc.stateMu.Unlock()
	snapshot := svc.controlRuntime.Snapshot()
	if len(snapshot.Pages) != 1 || snapshot.Pages[0].Path != "posts/second.md" {
		t.Fatalf("state transition did not refresh file projection: pages = %+v", snapshot.Pages)
	}
}

func TestBuilderAdminControlRuntime_CurrentReleaseDiagnostics(t *testing.T) {
	dir := t.TempDir()
	siteDir := filepath.Join(dir, "site")
	releaseDir := filepath.Join(siteDir, "releases", "release-1")
	if err := os.MkdirAll(filepath.Join(releaseDir, ".markata"), 0o755); err != nil {
		t.Fatal(err)
	}
	artifact, err := diagnostics.MarshalArtifact(diagnostics.ContentLedgerSnapshot{
		Entries: []diagnostics.ContentDisposition{{
			Path: "posts/foo.md", Emitted: true,
			Diagnostics: []diagnostics.Issue{{
				File: "posts/foo.md", Code: "duplicate-key", Severity: diagnostics.SeverityWarning,
				Range: diagnostics.Range{StartLine: 3, StartCol: 2},
			}},
		}},
	}, diagnostics.ArtifactBuildInfo{BuiltAt: time.Now().UTC()})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(releaseDir, diagnostics.DefaultArtifactPath), artifact, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join("releases", "release-1"), filepath.Join(siteDir, "current")); err != nil {
		t.Fatal(err)
	}
	svc, err := New(Config{SourceDir: dir, SiteDir: siteDir, HistoryDir: filepath.Join(dir, "history")})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = svc.leaderLock.Close() })
	svc.stateMu.Lock()
	svc.state.Builds = []BuildRecord{{ID: "build-1", Kind: "build", Status: "success", ReleaseID: "release-1", FinishedAt: time.Now().UTC()}}
	svc.saveStateLocked()
	svc.stateMu.Unlock()
	svc.refreshControlStateWithFiles(svc.viewState())
	snapshot := svc.controlRuntime.Snapshot()
	if len(snapshot.Pages) != 1 || len(snapshot.Pages[0].Diagnostics) != 1 || snapshot.Pages[0].Diagnostics[0].Line != 4 || snapshot.Pages[0].Diagnostics[0].JobID != "build-1" {
		t.Fatalf("projected page diagnostic = %+v", snapshot.Pages)
	}
	if len(snapshot.Jobs) != 1 || len(snapshot.Jobs[0].Diagnostics) != 1 {
		t.Fatalf("projected job diagnostic = %+v", snapshot.Jobs)
	}
}

func TestBuilderAdminControlRuntime_PersistsLogDiagnosticsAndShowsRecentLogs(t *testing.T) {
	dir := t.TempDir()
	svc, err := New(Config{SourceDir: dir, SiteDir: filepath.Join(dir, "site"), HistoryDir: filepath.Join(dir, "history")})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = svc.leaderLock.Close() })
	logPath := "build-1.log"
	if err := os.WriteFile(filepath.Join(svc.logDir, logPath), []byte("warning: missing image alt in posts/foo.md\nordinary log\nerror: template failed\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	svc.finishBuild(BuildRecord{ID: "build-1", Kind: "build", Status: "failed", LogPath: logPath, Error: "command failed with exit code 1", EnqueuedAt: time.Now().UTC()})
	if len(svc.state.Builds) != 1 || len(svc.state.Builds[0].LogDiagnostics) != 2 {
		t.Fatalf("persisted log diagnostics = %+v", svc.state.Builds)
	}
	if err := os.WriteFile(filepath.Join(svc.logDir, logPath), []byte("ordinary log written later\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	svc.refreshControlStateWithFiles(svc.viewState())
	snapshot := svc.controlRuntime.Snapshot()
	if len(snapshot.Jobs) != 1 || len(snapshot.Jobs[0].Logs) != 1 || len(snapshot.Jobs[0].Diagnostics) != 3 {
		t.Fatalf("job log or diagnostics lost after later output: %+v", snapshot.Jobs)
	}
	if snapshot.Jobs[0].Diagnostics[1].Code != "MARKATA-W900" || snapshot.Jobs[0].Diagnostics[2].Code != "MARKATA-E900" {
		t.Fatalf("log diagnostic codes = %+v", snapshot.Jobs[0].Diagnostics)
	}
}
