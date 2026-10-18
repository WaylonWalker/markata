package builderadmin

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/internal/testbinary"
	"github.com/WaylonWalker/markata/pkg/diagnostics"
)

func TestBuilderAdminWebhook_RejectsInvalidSignatureAndIgnoresOtherBranches(t *testing.T) {
	t.Parallel()
	svc, err := New(Config{
		SiteDir: t.TempDir(),
		Webhook: WebhookConfig{Enabled: true, Branch: "main", Secret: "test-secret"},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = svc.leaderLock.Close() })
	svc.leader = true
	mux := http.NewServeMux()
	svc.registerRoutes(mux)
	body := []byte(`{"ref":"refs/heads/dev","after":"abcdef"}`)

	invalid := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(string(body)))
	invalid.Header.Set("X-GitHub-Event", "push")
	invalid.Header.Set("X-Hub-Signature-256", "sha256=invalid")
	invalidRecorder := httptest.NewRecorder()
	mux.ServeHTTP(invalidRecorder, invalid)
	if invalidRecorder.Code != http.StatusUnauthorized {
		t.Fatalf("invalid signature status=%d, want %d", invalidRecorder.Code, http.StatusUnauthorized)
	}

	branchMismatch := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(string(body)))
	branchMismatch.Header.Set("X-GitHub-Event", "push")
	branchMismatch.Header.Set("X-Hub-Signature-256", webhookSignature(body, "test-secret"))
	branchMismatchRecorder := httptest.NewRecorder()
	mux.ServeHTTP(branchMismatchRecorder, branchMismatch)
	if branchMismatchRecorder.Code != http.StatusNoContent {
		t.Fatalf("branch mismatch status=%d, want %d", branchMismatchRecorder.Code, http.StatusNoContent)
	}
	if len(svc.snapshotState().Queue) != 0 {
		t.Fatal("branch mismatch queued a build")
	}
}

func TestNew_RejectsEnabledWebhookWithoutSecret(t *testing.T) {
	t.Parallel()
	_, err := New(Config{SiteDir: t.TempDir(), Webhook: WebhookConfig{Enabled: true}})
	if err == nil {
		t.Fatal("New() accepted an enabled webhook without a secret")
	}
}

func TestBuilderAdminWebhook_RejectsQueueWhenFull(t *testing.T) {
	t.Parallel()
	svc, err := New(Config{SiteDir: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = svc.leaderLock.Close() })
	for range cap(svc.queueCh) {
		svc.queueCh <- queueRequest{}
	}
	if err := svc.enqueueWebhookBuild("GitHub push"); !errors.Is(err, ErrBuildQueueFull) {
		t.Fatalf("enqueueWebhookBuild() error = %v, want ErrBuildQueueFull", err)
	}
	if len(svc.snapshotState().Queue) != 0 {
		t.Fatal("full queue left a persisted webhook operation")
	}
}

func TestBuilderAdminPullSource_ReportsChangedCommit(t *testing.T) {
	t.Parallel()
	remote := filepath.Join(t.TempDir(), "remote.git")
	runGit(t, "init", "--bare", remote)
	source := filepath.Join(t.TempDir(), "source")
	runGit(t, "clone", remote, source)
	runGitAt(t, source, "config", "user.email", "builder@example.com")
	runGitAt(t, source, "config", "user.name", "Builder")
	runGitAt(t, source, "checkout", "-b", "main")
	if err := os.WriteFile(filepath.Join(source, "post.md"), []byte("first"), 0o600); err != nil {
		t.Fatal(err)
	}
	runGitAt(t, source, "add", "post.md")
	runGitAt(t, source, "commit", "-m", "initial")
	runGitAt(t, source, "push", "-u", "origin", "main")

	updater := filepath.Join(t.TempDir(), "updater")
	runGit(t, "clone", "--branch", "main", remote, updater)
	runGitAt(t, updater, "config", "user.email", "updater@example.com")
	runGitAt(t, updater, "config", "user.name", "Updater")
	if err := os.WriteFile(filepath.Join(updater, "post.md"), []byte("second"), 0o600); err != nil {
		t.Fatal(err)
	}
	runGitAt(t, updater, "add", "post.md")
	runGitAt(t, updater, "commit", "-m", "update")
	runGitAt(t, updater, "push")

	svc, err := New(Config{SourceDir: source, SiteDir: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = svc.leaderLock.Close() })
	changed, err := svc.pullSource(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !changed {
		t.Fatal("pullSource() = unchanged, want changed")
	}
	changed, err = svc.pullSource(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if changed {
		t.Fatal("pullSource() = changed after synchronization")
	}
}

func runGit(t *testing.T, args ...string) {
	t.Helper()
	if output, err := exec.Command("git", args...).CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, output)
	}
}

func runGitAt(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git -C %s %s: %v\n%s", dir, strings.Join(args, " "), err, output)
	}
}

func webhookSignature(body []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(body)
	return fmt.Sprintf("sha256=%x", mac.Sum(nil))
}

func TestBuilderAdminAuthentication_RejectsUntrustedOrMissingIdentity(t *testing.T) {
	t.Parallel()
	svc, err := New(Config{SiteDir: t.TempDir(), TrustedProxyCIDRs: []string{"10.42.0.0/24"}})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = svc.leaderLock.Close() })
	mux := http.NewServeMux()
	svc.registerRoutes(mux)

	for _, tt := range []struct {
		name       string
		remoteAddr string
		identity   string
	}{
		{name: "missing identity", remoteAddr: "10.42.0.10:443"},
		{name: "untrusted source", remoteAddr: "192.0.2.10:443", identity: "operator"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = tt.remoteAddr
			if tt.identity != "" {
				req.Header.Set(DefaultAuthHeaders().UserID, tt.identity)
			}
			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, req)
			if recorder.Code != http.StatusUnauthorized {
				t.Fatalf("status=%d, want %d", recorder.Code, http.StatusUnauthorized)
			}
		})
	}
}

func TestBuilderAdminAuthentication_TrustedIdentityAndCSRF(t *testing.T) {
	t.Parallel()
	svc, err := New(Config{
		SiteDir:           t.TempDir(),
		TrustedProxyCIDRs: []string{"10.42.0.0/24"},
		PublicOrigin:      "https://builder.example.com",
		PublicAuthOrigin:  "https://auth.example.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = svc.leaderLock.Close() })
	svc.leader = true
	mux := http.NewServeMux()
	svc.registerRoutes(mux)

	indexRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	indexRequest.RemoteAddr = "10.42.0.10:443"
	indexRequest.Header.Set(DefaultAuthHeaders().UserID, "operator")
	indexRecorder := httptest.NewRecorder()
	mux.ServeHTTP(indexRecorder, indexRequest)
	if indexRecorder.Code != http.StatusOK {
		t.Fatalf("index status=%d, want %d", indexRecorder.Code, http.StatusOK)
	}
	indexBody := indexRecorder.Body.String()
	if !strings.Contains(indexBody, `operator-avatar-fallback[hidden] { display: none; }`) {
		t.Fatal("index is missing the hidden avatar fallback rule")
	}
	if !strings.Contains(indexBody, `src="https://auth.example.com/users/operator/picture"`) {
		t.Fatal("index is missing the operator profile picture URL")
	}
	if !strings.Contains(indexBody, `onerror="this.hidden=true;this.nextElementSibling.hidden=false"`) || !strings.Contains(indexBody, `hidden role="img" aria-label="Profile picture unavailable"`) || !strings.Contains(indexBody, `<svg viewBox="0 0 24 24"`) {
		t.Fatal("index is missing the accessible avatar fallback icon")
	}
	cookies := indexRecorder.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Name != csrfCookieName || !cookies[0].Secure || !cookies[0].HttpOnly {
		t.Fatalf("csrf cookie=%+v", cookies)
	}

	values := url.Values{"csrf_token": {cookies[0].Value}}
	buildRequest := httptest.NewRequest(http.MethodPost, "/api/builds", strings.NewReader(values.Encode()))
	buildRequest.RemoteAddr = "10.42.0.10:443"
	buildRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	buildRequest.Header.Set("Origin", "https://builder.example.com")
	buildRequest.Header.Set(DefaultAuthHeaders().UserID, "operator")
	buildRequest.AddCookie(cookies[0])
	buildRecorder := httptest.NewRecorder()
	mux.ServeHTTP(buildRecorder, buildRequest)
	if buildRecorder.Code != http.StatusSeeOther {
		t.Fatalf("build status=%d, want %d", buildRecorder.Code, http.StatusSeeOther)
	}
}

func TestBuilderAdminAuthentication_AuthentikHeaders(t *testing.T) {
	t.Parallel()
	authHeaders := AuthHeaders{
		UserID:      "X-Authentik-Uid",
		Username:    "X-Authentik-Username",
		DisplayName: "X-Authentik-Name",
		Email:       "X-Authentik-Email",
		Groups:      "X-Authentik-Groups",
	}
	svc, err := New(Config{
		SiteDir:           t.TempDir(),
		TrustedProxyCIDRs: []string{"10.42.0.0/24"},
		AuthHeaders:       authHeaders,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = svc.leaderLock.Close() })
	mux := http.NewServeMux()
	svc.registerRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "10.42.0.10:443"
	req.Header.Set("X-Authentik-Uid", "3c8d0ec5-023e-4a7f-ae09-f6985cebd4dc")
	req.Header.Set("X-Authentik-Username", "operator")
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status=%d, want %d", recorder.Code, http.StatusOK)
	}
}

func TestNew_RejectsInvalidAuthHeaderName(t *testing.T) {
	t.Parallel()
	_, err := New(Config{SiteDir: t.TempDir(), AuthHeaders: AuthHeaders{UserID: "X-Auth\nUser"}})
	if err == nil {
		t.Fatal("New() accepted an invalid auth header name")
	}
}

func TestNew_RejectsUnsafeTrustedProxyCIDRs(t *testing.T) {
	t.Parallel()
	for _, cidr := range []string{
		"0.0.0.0/0", "::/0", "127.0.0.0/8", "127.0.0.1/32", "::1/128",
		"169.254.0.0/16", "fe80::/10", "126.0.0.0/7", "fe00::/8",
	} {
		t.Run(cidr, func(t *testing.T) {
			_, err := New(Config{SiteDir: t.TempDir(), TrustedProxyCIDRs: []string{cidr}})
			if err == nil {
				t.Fatalf("New() with trusted proxy CIDR %q succeeded, want error", cidr)
			}
		})
	}

	service, err := New(Config{SiteDir: t.TempDir(), TrustedProxyCIDRs: []string{"10.42.0.0/16"}})
	if err != nil {
		t.Fatalf("New() with pod CIDR failed: %v", err)
	}
	t.Cleanup(func() { _ = service.leaderLock.Close() })
}

func TestNew_RejectsEmptyConfiguredUserIDHeader(t *testing.T) {
	t.Parallel()
	_, err := New(Config{
		SiteDir:     t.TempDir(),
		AuthHeaders: AuthHeaders{Username: "X-Authenticated-Username"},
	})
	if err == nil || !strings.Contains(err.Error(), "user ID header") {
		t.Fatalf("New() error = %v, want required user ID header error", err)
	}
}

func TestIgnoreWatchPath(t *testing.T) {
	t.Parallel()
	root := "/tmp/site"
	tests := []struct {
		path string
		want bool
	}{
		{path: "/tmp/site/pages/post.md", want: false},
		{path: "/tmp/site/.git/index", want: true},
		{path: "/tmp/site/.markata/cache.json", want: true},
		{path: "/tmp/site/.builder-admin/state.json", want: true},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := ignoreWatchPath(root, tt.path); got != tt.want {
				t.Fatalf("ignoreWatchPath(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestExtractPerfSummaryFromFileMissing(t *testing.T) {
	t.Parallel()
	if got := extractPerfSummaryFromFile("/does/not/exist"); got != nil {
		t.Fatalf("extractPerfSummaryFromFile() = %#v, want nil", got)
	}
}

func TestIndexHTMLIncludesDynamicFavicon(t *testing.T) {
	t.Parallel()
	checks := []string{
		`id="app-favicon"`,
		`function updateFavicon(stateName)`,
		`function faviconState(state)`,
		`updateFavicon('error');`,
	}
	for _, check := range checks {
		if !strings.Contains(indexHTML, check) {
			t.Fatalf("indexHTML missing %q", check)
		}
	}
}

func TestIndexHTMLPrioritizesOperationalBuildData(t *testing.T) {
	t.Parallel()
	checks := []string{
		`<h2 id="live-work-heading">Active work</h2>`,
		`<h2>Jobs</h2>`,
		`<th>Job</th>`,
		`<th>Run</th>`,
		`Running and queued work first`,
		`class="build-details"`,
		`Open raw log`,
		`.responsive-table td::before`,
		`.responsive-table thead { position: absolute;`,
	}
	for _, check := range checks {
		if !strings.Contains(indexHTML, check) {
			t.Fatalf("indexHTML missing %q", check)
		}
	}
	if strings.Contains(indexHTML, `background-size: auto, auto, 11px 11px`) {
		t.Fatal("index retains the decorative grid background")
	}
	for _, unwanted := range []string{`.Operator.Email`, `.Operator.Groups`, `.Operator.Roles`, `.Operator.Scopes`, `.Operator.UserID`} {
		if strings.Contains(indexHTML, unwanted) {
			t.Fatalf("index exposes %q in the default identity view", unwanted)
		}
	}
}

func TestHandleIndex_RendersOperationalSummary(t *testing.T) {
	t.Parallel()
	svc, err := New(Config{SiteDir: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = svc.leaderLock.Close() })
	svc.leader = true
	svc.state = State{
		Queue: []QueuedOperation{{
			ID:         "build-queued",
			Label:      "Build",
			EnqueuedAt: time.Now().Add(-2 * time.Minute),
		}},
		Builds: []BuildRecord{{
			ID:         "build-success",
			Status:     "success",
			FinishedAt: time.Now(),
		}},
	}
	recorder := httptest.NewRecorder()
	svc.handleIndex(recorder, httptest.NewRequest(http.MethodGet, "/", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status=%d, want %d", recorder.Code, http.StatusOK)
	}
	body := recorder.Body.String()
	for _, want := range []string{
		`<style>:root{color-scheme:dark;--markata-background:#09090b`,
		`data-capabilities=`, `sourceMutation&#34;:false`, `enqueueBuild&#34;:true`,
		`const capabilities = JSON.parse(document.body.dataset.capabilities || '{}')`,
		`data-capability="enqueueBuild"`,
		`id="active-work-detail"`,
		`let buildsFingerprint = ''`,
		`let refreshFingerprint = ''`,
		`let releasesFingerprint = ''`,
		`let pollInFlight = false`,
		`renderBuilds(state)`,
		`function buildTimeBaseline(builds)`,
		`slice(0, 60)`,
		`30 * 24 * 60 * 60 * 1000`,
		`data-expandable`,
		`View build record →`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("index body missing %q", want)
		}
	}
}

func TestCompletedJobs_SortsBuildsAndRefreshesByCompletion(t *testing.T) {
	t.Parallel()
	older := time.Date(2026, 7, 25, 20, 0, 0, 0, time.UTC)
	newer := older.Add(time.Minute)
	jobs := completedJobs(State{
		Builds:  []BuildRecord{{ID: "build", FinishedAt: older}},
		Refresh: []RefreshRecord{{ID: "refresh", TaskName: "reader", FinishedAt: newer}},
	})
	if len(jobs) != 2 {
		t.Fatalf("len(jobs)=%d, want 2", len(jobs))
	}
	if jobs[0].ID != "refresh" || jobs[1].ID != "build" {
		t.Fatalf("jobs=%#v, want refresh before build", jobs)
	}
}

func TestBuildDetail_ReturnsBuildAndNotFound(t *testing.T) {
	t.Parallel()
	svc, err := New(Config{SiteDir: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = svc.leaderLock.Close() })
	svc.leader = true
	svc.state.Builds = []BuildRecord{{ID: "build-1", Status: "success", FinishedAt: time.Now()}}
	for _, tt := range []struct {
		path string
		want int
	}{{"/builds/build-1", http.StatusOK}, {"/builds/missing", http.StatusNotFound}} {
		recorder := httptest.NewRecorder()
		svc.handleBuildDetail(recorder, httptest.NewRequest(http.MethodGet, tt.path, nil))
		if recorder.Code != tt.want {
			t.Errorf("%s status=%d, want %d", tt.path, recorder.Code, tt.want)
		}
	}
}

func TestReleaseTimestampFromID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		id   string
		want string
		ok   bool
	}{
		{name: "builder admin style", id: "20260711T203402Z-pod-name", want: "2026-07-11T20:34:02Z", ok: true},
		{name: "legacy numeric", id: "20260711160013", want: "2026-07-11T16:00:13Z", ok: true},
		{name: "invalid", id: "not-a-release", ok: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := releaseTimestampFromID(tt.id)
			if ok != tt.ok {
				t.Fatalf("releaseTimestampFromID(%q) ok=%v want %v", tt.id, ok, tt.ok)
			}
			if !tt.ok {
				return
			}
			if got.Format(time.RFC3339) != tt.want {
				t.Fatalf("releaseTimestampFromID(%q)=%s want %s", tt.id, got.Format(time.RFC3339), tt.want)
			}
		})
	}
}

func TestDiscoverReleasesPrefersBuildFinishedAtAndCurrentFirst(t *testing.T) {
	t.Parallel()
	siteDir := t.TempDir()
	historyDir := filepath.Join(siteDir, ".builder-admin")
	releasesDir := filepath.Join(siteDir, "releases")
	if err := os.MkdirAll(releasesDir, 0o755); err != nil {
		t.Fatal(err)
	}
	older := "20260711T153738Z-old"
	current := "20260711T203402Z-current"
	for _, releaseID := range []string{older, current} {
		if err := os.MkdirAll(filepath.Join(releasesDir, releaseID), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Symlink(filepath.Join("releases", current), filepath.Join(siteDir, "current")); err != nil {
		t.Fatal(err)
	}
	svc, err := New(Config{SiteDir: siteDir, HistoryDir: historyDir})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if svc.leaderLock != nil {
			_ = svc.leaderLock.Close()
		}
	})
	svc.state.Builds = []BuildRecord{
		{ID: "build-new", ReleaseID: current, Status: "success", FinishedAt: time.Date(2026, 7, 11, 20, 34, 4, 0, time.UTC)},
		{ID: "build-old", ReleaseID: older, Status: "success", FinishedAt: time.Date(2026, 7, 11, 15, 37, 50, 0, time.UTC)},
	}
	svc.leader = true
	views := svc.discoverReleases()
	if len(views) != 2 {
		t.Fatalf("discoverReleases() len=%d want 2", len(views))
	}
	if !views[0].Current || views[0].ID != current {
		t.Fatalf("views[0]=%+v want current release first", views[0])
	}
	if views[0].BuildID != "build-new" {
		t.Fatalf("views[0].BuildID=%q want build-new", views[0].BuildID)
	}
	if got := views[0].CreatedAt.Format(time.RFC3339); got != "2026-07-11T20:34:04Z" {
		t.Fatalf("views[0].CreatedAt=%s want 2026-07-11T20:34:04Z", got)
	}
}

func TestBuilderAdminBuild_RemovesDeletedPostFromNewRelease(t *testing.T) {
	sourceDir := t.TempDir()
	contentDir := filepath.Join(sourceDir, "content")
	if err := os.MkdirAll(contentDir, 0o755); err != nil {
		t.Fatal(err)
	}
	config := `[markata-go]
title = "Builder Admin deletion fixture"
description = "A small deletion fixture."
license = false
output_dir = "output"

[markata-go.assets]
mode = "cdn"

# Minification is covered separately; this fixture tests deleted-output pruning.
[markata-go.js_minify]
enabled = false

[markata-go.css_minify]
enabled = false

[markata-go.glob]
patterns = ["content/**/*.md"]
use_gitignore = false

[markata-go.post_formats]
html = true
markdown = true
text = true
ansi = true
og = true
`
	if err := os.WriteFile(filepath.Join(sourceDir, "markata-go.toml"), []byte(config), 0o600); err != nil {
		t.Fatal(err)
	}
	for name, contents := range map[string]string{
		"keep.md": `---
title: Keep
published: true
---

Keep this post.
`,
		"deleted.md": `---
title: Deleted
published: true
---

Delete this post.
`,
	} {
		if err := os.WriteFile(filepath.Join(contentDir, name), []byte(contents), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller() failed")
	}
	moduleRoot := filepath.Clean(filepath.Join(filepath.Dir(sourceFile), "..", ".."))
	binary, err := testbinary.Resolve(moduleRoot, t.TempDir())
	if err != nil {
		t.Fatalf("resolve markata-go test binary: %v", err)
	}

	siteDir := t.TempDir()
	svc, err := New(Config{
		SourceDir:    sourceDir,
		SiteDir:      siteDir,
		ConfigPath:   "markata-go.toml",
		ReleasesKeep: 10,
		BuildTimeout: 5 * time.Minute,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = svc.leaderLock.Close() })
	svc.executable = binary

	// runBuild schedules release pruning in the background, and that goroutine
	// persists state under siteDir/.builder-admin. Wait for it so TempDir
	// cleanup does not race with the final state write.
	waitForPrune := func() {
		svc.pruneScheduleMu.Lock()
		done := svc.pruneDone
		svc.pruneScheduleMu.Unlock()
		if done != nil {
			<-done
		}
	}
	t.Cleanup(waitForPrune)

	run := func(id string) BuildRecord {
		t.Helper()
		svc.runBuild(context.Background(), queueRequest{QueuedOperation: QueuedOperation{
			ID:          id,
			Kind:        "build",
			Label:       "Build",
			TriggerType: "test",
			EnqueuedAt:  time.Now().UTC(),
		}})
		waitForPrune()
		state := svc.snapshotState()
		if len(state.Builds) == 0 {
			t.Fatal("build did not record a result")
		}
		if state.Builds[0].Status != "success" {
			t.Fatalf("build %s status = %q, error = %q", id, state.Builds[0].Status, state.Builds[0].Error)
		}
		return state.Builds[0]
	}

	first := run("build-first")
	readArtifact := func(releasePath string) diagnostics.Artifact {
		t.Helper()
		data, err := os.ReadFile(filepath.Join(releasePath, diagnostics.DefaultArtifactPath))
		if err != nil {
			t.Fatalf("read diagnostics artifact from release: %v", err)
		}
		artifact, err := diagnostics.ParseArtifact(data)
		if err != nil {
			t.Fatalf("parse diagnostics artifact from release: %v", err)
		}
		return artifact
	}
	firstArtifact := readArtifact(first.ReleasePath)
	if firstArtifact.Summary.Discovered != 2 {
		t.Fatalf("first release diagnostics summary = %+v, want two discovered files", firstArtifact.Summary)
	}
	oldDeleted := filepath.Join(first.ReleasePath, "deleted", "index.html")
	oldContents, err := os.ReadFile(oldDeleted)
	if err != nil {
		t.Fatalf("read first release deleted output: %v", err)
	}
	oldKeep := filepath.Join(first.ReleasePath, "keep", "index.html")
	oldKeepContents, err := os.ReadFile(oldKeep)
	if err != nil {
		t.Fatalf("read first release keep output: %v", err)
	}
	if err := os.Remove(filepath.Join(contentDir, "deleted.md")); err != nil {
		t.Fatal(err)
	}
	second := run("build-second")
	secondArtifact := readArtifact(second.ReleasePath)
	if secondArtifact.Summary.Discovered != 1 {
		t.Fatalf("second release diagnostics summary = %+v, want one discovered file", secondArtifact.Summary)
	}
	if first.ReleasePath == second.ReleasePath {
		t.Fatalf("release paths are identical: %q", first.ReleasePath)
	}

	for _, relative := range []string{
		"deleted/index.html",
		"deleted/index.md/index.html",
		"deleted/index.txt/index.html",
		"deleted/index.ansi/index.html",
		"deleted/og/index.html",
		"deleted.md",
		"deleted.txt",
		"deleted.ansi",
	} {
		if _, err := os.Stat(filepath.Join(second.ReleasePath, filepath.FromSlash(relative))); !os.IsNotExist(err) {
			t.Errorf("new release retained deleted output %s: %v", relative, err)
		}
	}
	if _, err := os.Stat(filepath.Join(second.ReleasePath, "keep", "index.html")); err != nil {
		t.Fatalf("new release lost kept output: %v", err)
	}
	currentTarget, err := os.Readlink(filepath.Join(siteDir, "current"))
	if err != nil {
		t.Fatalf("read current release link: %v", err)
	}
	wantCurrentTarget := filepath.Join("releases", filepath.Base(second.ReleasePath))
	if currentTarget != wantCurrentTarget {
		t.Fatalf("current release target = %q, want %q", currentTarget, wantCurrentTarget)
	}
	if got, err := os.ReadFile(oldDeleted); err != nil || string(got) != string(oldContents) {
		t.Fatalf("historical deleted output changed: error = %v", err)
	}
	if got, err := os.ReadFile(oldKeep); err != nil || string(got) != string(oldKeepContents) {
		t.Fatalf("historical kept output changed: error = %v", err)
	}
}
