package builderadmin

import (
	"context"
	"crypto/sha256"
	"io"
	"io/fs"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/internal/fileidentity"
	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestVerifiedNoOpLeavesCurrentReleaseAndWorkspaceUntouched(t *testing.T) {
	if !fileidentity.Available {
		t.Skip("verified no-op requires Linux file identity")
	}
	svc, _ := newNoOpFixture(t)
	snapshot, err := svc.noOpInputSnapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	checkSnapshot, err := svc.noOpInputSnapshot(context.Background())
	if err != nil || checkSnapshot.Fingerprint != snapshot.Fingerprint {
		t.Fatalf("fixture snapshot changed before no-op: before=%+v after=%+v error=%v", snapshot, checkSnapshot, err)
	}
	if err := os.MkdirAll(filepath.Join(svc.cfg.SiteDir, "releases", "release-a"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(svc.cfg.SiteDir, "releases", "release-a", "index.html"), []byte("live"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join("releases", "release-a"), filepath.Join(svc.cfg.SiteDir, "current")); err != nil {
		t.Fatal(err)
	}
	if err := svc.saveNoOpBaseline(context.Background(), snapshot, "release-a", time.Now()); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(svc.cfg.WorkDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(svc.cfg.WorkDir, "index.html"), []byte("workspace"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := markReusableWorkspace(svc.cfg.WorkDir, "release-a"); err != nil {
		t.Fatal(err)
	}
	if got := svc.currentReleaseID(); got != "release-a" {
		t.Fatalf("current release id=%q", got)
	}
	if ok, err := reusableWorkspaceMatches(svc.cfg.WorkDir, "release-a"); err != nil || !ok {
		t.Fatalf("workspace proof=%t, error=%v", ok, err)
	}
	setupSnapshot, err := svc.noOpInputSnapshot(context.Background())
	if err != nil || setupSnapshot.Fingerprint != snapshot.Fingerprint {
		t.Fatalf("inputs changed during runtime setup: before=%+v after=%+v error=%v", snapshot, setupSnapshot, err)
	}
	loaded, err := svc.loadNoOpBaseline()
	if err != nil || loaded.ReleaseID != "release-a" || !time.Now().Before(loaded.ExpiresAt) {
		t.Fatalf("baseline=%+v, error=%v", loaded, err)
	}

	started := time.Now().UTC()
	record := &BuildRecord{ID: "build-noop", Kind: "build", Status: "running", StartedAt: started, EnqueuedAt: started, TriggerType: "manual-ui"}
	req := queueRequest{QueuedOperation: QueuedOperation{ID: record.ID, Kind: "build", TriggerType: "manual-ui", EnqueuedAt: started}}
	if !svc.canUseNoOp(req) {
		t.Fatal("manual request failed no-op preconditions")
	}
	if ok, _ := svc.verifiedNoOp(context.Background(), req, record, io.Discard); !ok {
		t.Fatal("verified unchanged request did not complete as a no-op")
	}
	if record.Status != "success" || record.Result != "unchanged" || record.BecameLive || record.ReleaseID != "release-a" {
		t.Fatalf("no-op record = %+v", record)
	}
	if record.NoOpReleaseEntries != 1 {
		t.Fatalf("release entries = %d, want 1", record.NoOpReleaseEntries)
	}
	if got, err := os.Readlink(filepath.Join(svc.cfg.SiteDir, "current")); err != nil || got != filepath.Join("releases", "release-a") {
		t.Fatalf("current symlink = %q, error=%v", got, err)
	}
	if got, err := os.ReadFile(workspaceReuseMarkerPath(svc.cfg.WorkDir)); err != nil || strings.TrimSpace(string(got)) != "release-a" {
		t.Fatalf("workspace marker = %q, error=%v", got, err)
	}
	if got, err := os.ReadFile(filepath.Join(svc.cfg.WorkDir, "index.html")); err != nil || string(got) != "workspace" {
		t.Fatalf("workspace output = %q, error=%v", got, err)
	}
	if got, err := os.ReadFile(svc.statePath); err != nil || !strings.Contains(string(got), `"result": "unchanged"`) {
		t.Fatalf("persisted build state = %q, error=%v", got, err)
	}

	indexPath := filepath.Join(svc.cfg.SiteDir, "releases", "release-a", "index.html")
	indexInfo, err := os.Stat(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(indexPath, []byte("edit"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(indexPath, indexInfo.ModTime(), indexInfo.ModTime()); err != nil {
		t.Fatal(err)
	}
	changed := &BuildRecord{ID: "build-changed-release", Kind: "build", Status: "running", StartedAt: time.Now().UTC(), TriggerType: "manual-ui"}
	changedReq := queueRequest{QueuedOperation: QueuedOperation{ID: changed.ID, Kind: "build", TriggerType: "manual-ui", EnqueuedAt: changed.StartedAt}}
	if ok, _ := svc.verifiedNoOp(context.Background(), changedReq, changed, io.Discard); ok {
		t.Fatal("live release root change was accepted as unchanged")
	}
}

func TestNoOpRejectsCurrentReleaseOutsideReleaseRoot(t *testing.T) {
	if !fileidentity.Available {
		t.Skip("verified no-op requires Linux file identity")
	}
	svc, _ := newNoOpFixture(t)
	release := filepath.Join(svc.cfg.SiteDir, "releases", "release-a")
	if err := os.MkdirAll(release, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(release, "index.html"), []byte("live"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join("..", "..", "outside"), filepath.Join(svc.cfg.SiteDir, "current")); err != nil {
		t.Fatal(err)
	}
	if _, _, err := svc.noOpReleaseFingerprint(context.Background(), "outside"); err == nil {
		t.Fatal("current symlink outside the release root passed validation")
	}
}

func TestNoOpBuilderFingerprintDetectsExecutableRewrite(t *testing.T) {
	svc := &Service{executable: filepath.Join(t.TempDir(), "builder")}
	if err := os.WriteFile(svc.executable, []byte("first"), 0o700); err != nil {
		t.Fatal(err)
	}
	before, err := svc.noOpBuilderFingerprint()
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(svc.executable)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(10 * time.Millisecond)
	if err := os.WriteFile(svc.executable, []byte("other"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(svc.executable, info.ModTime(), info.ModTime()); err != nil {
		t.Fatal(err)
	}
	after, err := svc.noOpBuilderFingerprint()
	if err != nil {
		t.Fatal(err)
	}
	if before == after {
		t.Fatal("rewritten builder executable did not change its fingerprint")
	}
}

func TestBuildCacheSaveFailureDisablesNoOpBaseline(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "build.log")
	if err := os.WriteFile(logPath, []byte("build completed\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	failed, err := buildCacheSaveFailed(logPath)
	if err != nil || failed {
		t.Fatalf("clean build log result=%t, error=%v", failed, err)
	}
	if err := os.WriteFile(logPath, []byte("[build_cache] failed to save build cache: disk full\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	failed, err = buildCacheSaveFailed(logPath)
	if err != nil || !failed {
		t.Fatalf("failed cache save result=%t, error=%v", failed, err)
	}
}

func TestNoOpSnapshotDetectsSameSizeSameMtimeCacheRewrite(t *testing.T) {
	if !fileidentity.Available {
		t.Skip("verified no-op requires Linux file identity")
	}
	svc, root := newNoOpFixture(t)
	cacheFile := filepath.Join(root, "cache", "mentions", "entry.json")
	if err := os.MkdirAll(filepath.Dir(cacheFile), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cacheFile, []byte("aaaa"), 0o600); err != nil {
		t.Fatal(err)
	}
	before, err := svc.noOpInputSnapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(cacheFile)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(10 * time.Millisecond)
	if err := os.WriteFile(cacheFile, []byte("bbbb"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(cacheFile, info.ModTime(), info.ModTime()); err != nil {
		t.Fatal(err)
	}
	after, err := svc.noOpInputSnapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if before.Fingerprint == after.Fingerprint {
		t.Fatal("same-size, same-mtime cache rewrite did not invalidate the snapshot")
	}
}

func TestNoOpSnapshotTracksCustomWebmentionsCache(t *testing.T) {
	if !fileidentity.Available {
		t.Skip("verified no-op requires Linux file identity")
	}
	svc, root := newNoOpFixture(t)
	cacheFile := filepath.Join(root, "custom", "webmentions", "received_mentions.json")
	if err := os.MkdirAll(filepath.Dir(cacheFile), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cacheFile, []byte("first"), 0o600); err != nil {
		t.Fatal(err)
	}
	before, err := svc.noOpInputSnapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(10 * time.Millisecond)
	if err := os.WriteFile(cacheFile, []byte("other"), 0o600); err != nil {
		t.Fatal(err)
	}
	after, err := svc.noOpInputSnapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if before.Fingerprint == after.Fingerprint {
		t.Fatal("custom webmentions cache change did not invalidate the no-op snapshot")
	}
}

func TestNoOpSnapshotTracksDotEnvFileIdentity(t *testing.T) {
	if !fileidentity.Available {
		t.Skip("verified no-op requires Linux file identity")
	}
	svc, root := newNoOpFixture(t)
	dotEnv := filepath.Join(root, ".env")
	if err := os.WriteFile(dotEnv, []byte("# alpha\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	before, err := svc.noOpInputSnapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(dotEnv)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(10 * time.Millisecond)
	if err := os.WriteFile(dotEnv, []byte("# bravo\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(dotEnv, info.ModTime(), info.ModTime()); err != nil {
		t.Fatal(err)
	}
	after, err := svc.noOpInputSnapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if before.Fingerprint == after.Fingerprint {
		t.Fatal("same-size, same-mtime .env rewrite did not invalidate the no-op snapshot")
	}
}

func TestNoOpSnapshotTracksPluginConfigOutsideSource(t *testing.T) {
	if !fileidentity.Available {
		t.Skip("verified no-op requires Linux file identity")
	}
	svc, _ := newNoOpFixture(t)
	configPath := filepath.Join(t.TempDir(), "external-config.json")
	svc.cfg.ConfigPath = configPath
	writeNoOpFile(t, configPath, `{"markata-go":{"glob":{"use_gitignore":true},"tailwind":{"build":false},"search":{"backend":"bleve"},"well_known":{"auto_generate":["links"]},"custom_plugin":{"enabled":true}}}`)
	before, err := svc.noOpInputSnapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(10 * time.Millisecond)
	writeNoOpFile(t, configPath, `{"markata-go":{"glob":{"use_gitignore":true},"tailwind":{"build":false},"search":{"backend":"bleve"},"well_known":{"auto_generate":["links"]},"custom_plugin":{"enabled":false}}}`)
	after, err := svc.noOpInputSnapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if before.Fingerprint == after.Fingerprint {
		t.Fatal("plugin config outside source did not invalidate the no-op snapshot")
	}
}

func TestNoOpFinalCheckDetectsConfigRewriteOutsideSource(t *testing.T) {
	if !fileidentity.Available {
		t.Skip("verified no-op requires Linux file identity")
	}
	svc, _ := newNoOpFixture(t)
	configPath := filepath.Join(t.TempDir(), "external-config.json")
	svc.cfg.ConfigPath = configPath
	writeNoOpFile(t, configPath, `{"markata-go":{"glob":{"use_gitignore":true},"tailwind":{"build":false},"search":{"backend":"bleve"},"well_known":{"auto_generate":["links"]},"custom_plugin":{"enabled":true}}}`)
	snapshot, err := svc.captureNoOpSnapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(10 * time.Millisecond)
	writeNoOpFile(t, configPath, `{"markata-go":{"glob":{"use_gitignore":true},"tailwind":{"build":false},"search":{"backend":"bleve"},"well_known":{"auto_generate":["links"]},"custom_plugin":{"enabled":false}}}`)
	if err := os.Chtimes(configPath, info.ModTime(), info.ModTime()); err != nil {
		t.Fatal(err)
	}
	confirmation, err := svc.confirmNoOpSnapshot(context.Background(), snapshot)
	if err != nil {
		t.Fatal(err)
	}
	if confirmation.InputFingerprint == snapshot.InputFingerprint {
		t.Fatal("final stability check did not detect external config rewrite")
	}
}

func TestNoOpSnapshotTracksCustomBuildCacheRoot(t *testing.T) {
	if !fileidentity.Available {
		t.Skip("verified no-op requires Linux file identity")
	}
	svc, _ := newNoOpFixture(t)
	root := svc.cfg.SourceDir
	cacheDir := filepath.Join(t.TempDir(), "builder-cache")
	writeNoOpFile(t, filepath.Join(cacheDir, "build-cache.json"), `{"version":1}`)
	configPath := filepath.Join(root, "markata-go.json")
	writeNoOpFile(t, configPath, `{"markata-go":{"glob":{"use_gitignore":true},"tailwind":{"build":false},"search":{"backend":"bleve"},"well_known":{"auto_generate":["links"]},"cache_dir":"`+cacheDir+`"}}`)
	svc.cfg.ConfigPath = configPath
	before, err := svc.noOpInputSnapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(filepath.Join(cacheDir, "build-cache.json"))
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(10 * time.Millisecond)
	if err := os.WriteFile(filepath.Join(cacheDir, "build-cache.json"), []byte(`{"version":2}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(filepath.Join(cacheDir, "build-cache.json"), info.ModTime(), info.ModTime()); err != nil {
		t.Fatal(err)
	}
	after, err := svc.noOpInputSnapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if before.Fingerprint == after.Fingerprint {
		t.Fatal("custom build cache root did not invalidate the no-op snapshot")
	}
}

func TestNoOpRejectsDisabledBuildCacheAndMissingExternalEmbedTransform(t *testing.T) {
	if !fileidentity.Available {
		t.Skip("verified no-op requires Linux file identity")
	}
	svc, root := newNoOpFixture(t)
	siteConfig, _, err := svc.loadNoOpSiteConfig()
	if err != nil {
		t.Fatal(err)
	}
	if siteConfig.Extra == nil {
		siteConfig.Extra = make(map[string]any)
	}
	siteConfig.Extra["build_cache"] = map[string]any{"enabled": false}
	if err := validateNoOpConfig(siteConfig); err == nil {
		t.Fatal("disabled build cache was accepted")
	}
	delete(siteConfig.Extra, "build_cache")
	siteConfig.DisabledHooks = []string{"build_cache"}
	if err := validateNoOpConfig(siteConfig); err == nil {
		t.Fatal("disabled build_cache hook was accepted")
	}
	siteConfig.DisabledHooks = nil
	writeNoOpFile(t, filepath.Join(root, "posts", "remote-embed.md"), "![embed](https://example.com/post)\n")
	if err := checkNoOpEmbedCache(context.Background(), root, "", siteConfig); err == nil {
		t.Fatal("external embed without a persisted transform cache was accepted")
	}
}

func TestNoOpRequiresTailwindBuildDisabled(t *testing.T) {
	config := newNoOpValidationConfig()
	config.Search.Backend = "bleve"
	if err := validateNoOpConfig(config); err == nil {
		t.Fatal("omitted Tailwind build option unexpectedly passed no-op validation")
	}
	config.Extra = map[string]any{"tailwind": map[string]any{"build": true}}
	if err := validateNoOpConfig(config); err == nil {
		t.Fatal("Tailwind build unexpectedly passed no-op validation")
	}
	config.Extra["tailwind"] = map[string]any{"build": false}
	if err := validateNoOpConfig(config); err != nil {
		t.Fatalf("explicitly disabled Tailwind build rejected: %v", err)
	}
}

func TestNoOpRejectsUnfingerprintedDefaultPluginInputs(t *testing.T) {
	tests := []struct {
		name  string
		extra map[string]any
	}{
		{name: "icon vendor", extra: map[string]any{"tailwind": map[string]any{"build": false}, "icons": map[string]any{"vendor": map[string]any{"enabled": true}}}},
		{name: "Python docs", extra: map[string]any{"tailwind": map[string]any{"build": false}, "python_docs": map[string]any{"enabled": true}}},
		{name: "external link hover", extra: map[string]any{"tailwind": map[string]any{"build": false}, "external_link_hover": map[string]any{"enabled": true}}},
		{name: "Mermaid CLI", extra: map[string]any{"tailwind": map[string]any{"build": false}, "mermaid": map[string]any{"mode": "cli"}}},
		{name: "Excalidraw Chromium", extra: map[string]any{"tailwind": map[string]any{"build": false}, "excalidraw": map[string]any{"mode": "chromium"}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := newNoOpValidationConfig()
			config.Search.Backend = "bleve"
			config.Extra = test.extra
			if err := validateNoOpConfig(config); err == nil {
				t.Fatal("unfingerprinted default-plugin input unexpectedly passed validation")
			}
		})
	}
}

func TestNoOpRequiresStableSearchBackend(t *testing.T) {
	config := newNoOpValidationConfig()
	config.Extra = map[string]any{"tailwind": map[string]any{"build": false}}
	if err := validateNoOpConfig(config); err == nil {
		t.Fatal("default Pagefind backend unexpectedly passed no-op validation")
	}
	config.Search.Backend = "bleve"
	if err := validateNoOpConfig(config); err != nil {
		t.Fatalf("Bleve backend rejected: %v", err)
	}
}

func TestNoOpRejectsExternalMermaidBuilderOverride(t *testing.T) {
	config := newNoOpValidationConfig()
	config.Search.Backend = "bleve"
	config.Extra = map[string]any{"tailwind": map[string]any{"build": false}}
	if err := validateNoOpConfig(config, "cli"); err == nil {
		t.Fatal("external Mermaid builder override unexpectedly passed no-op validation")
	}
	if err := validateNoOpConfig(config, "go"); err != nil {
		t.Fatalf("native Go Mermaid override rejected: %v", err)
	}
}

func TestNoOpRequiresBuildTimeLinkAvatarCacheCoverage(t *testing.T) {
	if !fileidentity.Available {
		t.Skip("verified no-op requires Linux file identity")
	}
	svc, root := newNoOpFixture(t)
	siteConfig, _, err := svc.loadNoOpSiteConfig()
	if err != nil {
		t.Fatal(err)
	}
	siteConfig.Extra["link_avatars"] = map[string]any{"enabled": true, "mode": "local"}
	writeNoOpFile(t, filepath.Join(root, "posts", "avatar.md"), "A link to [example](https://example.com/).\n")
	if err := checkNoOpLinkAvatarCache(context.Background(), root, "", siteConfig); err == nil {
		t.Fatal("missing link-avatar transform cache was accepted")
	}
	cache, err := buildcache.Load(noOpBuildCacheDir(root, "", siteConfig))
	if err != nil {
		t.Fatal(err)
	}
	cache.CacheLinkAvatarsHTML("posts/avatar.md", "input-hash", "<p>cached</p>")
	if err := cache.Save(); err != nil {
		t.Fatal(err)
	}
	if err := checkNoOpLinkAvatarCache(context.Background(), root, "", siteConfig); err != nil {
		t.Fatalf("complete link-avatar transform cache rejected: %v", err)
	}
}

func TestNoOpRequiresCompleteSelfHostedAssetCache(t *testing.T) {
	root := t.TempDir()
	siteConfig := models.NewConfig()
	siteConfig.Assets = models.NewAssetsConfig()
	if err := checkNoOpCDNAssetCache(root, siteConfig); err == nil {
		t.Fatal("missing self-hosted asset cache was accepted")
	}
	siteConfig.Assets.Mode = "cdn"
	if err := checkNoOpCDNAssetCache(root, siteConfig); err != nil {
		t.Fatalf("external CDN mode unexpectedly required a local cache: %v", err)
	}
}

func TestNoOpRejectsMentionCandidatesWithoutDomainCoverage(t *testing.T) {
	if !fileidentity.Available {
		t.Skip("verified no-op requires Linux file identity")
	}
	svc, root := newNoOpFixture(t)
	writeNoOpFile(t, filepath.Join(root, "posts", "mention.md"), "A reference to @someone.\n")
	siteConfig, _, err := svc.loadNoOpSiteConfig()
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	if _, err := checkNoOpMentionFreshness(context.Background(), root, siteConfig, now); err == nil {
		t.Fatal("mention references were accepted without proving complete domain cache coverage")
	}
}

func TestNoOpRejectsGitOnlyIgnoredContentPaths(t *testing.T) {
	if !fileidentity.Available {
		t.Skip("verified no-op requires Linux file identity")
	}
	for _, ignoredPath := range []string{"hidden.md", "hidden/"} {
		t.Run(strings.TrimSuffix(ignoredPath, "/"), func(t *testing.T) {
			svc, root := newNoOpFixture(t)
			writeNoOpFile(t, filepath.Join(root, ".git", "info", "exclude"), ignoredPath+"\n")
			if strings.HasSuffix(ignoredPath, "/") {
				writeNoOpFile(t, filepath.Join(root, ignoredPath, "post.md"), "hidden input")
			} else {
				writeNoOpFile(t, filepath.Join(root, ignoredPath), "hidden input")
			}
			if _, err := svc.noOpInputSnapshot(context.Background()); err == nil {
				t.Fatal("Git-only ignored path was accepted as a verified no-op input")
			}
		})
	}
}

func TestNoOpSnapshotHashesExecutableSidecarThemes(t *testing.T) {
	if !fileidentity.Available {
		t.Skip("verified no-op requires Linux file identity")
	}
	svc, root := newNoOpFixture(t)
	writeNoOpFile(t, filepath.Join(root, "markata-go.json"), `{"markata-go":{"templates_dir":"templates","assets_dir":"static","glob":{"use_gitignore":true,"patterns":["**/*.md"]},"tailwind":{"build":false},"search":{"backend":"bleve"},"well_known":{"auto_generate":["links"]},"theme":{"name":"custom"},"webmentions":{"cache_dir":"custom/webmentions"}}}`)
	binDir := filepath.Join(t.TempDir(), "bin")
	svc.executable = filepath.Join(binDir, "markata-go")
	if err := os.MkdirAll(filepath.Join(binDir, "themes", "custom", "templates"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(binDir, "themes", "custom", "static"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeNoOpFile(t, svc.executable, "builder executable")
	templatePath := filepath.Join(binDir, "themes", "custom", "templates", "base.html")
	staticPath := filepath.Join(binDir, "themes", "custom", "static", "theme.css")
	writeNoOpFile(t, templatePath, "template-v1")
	writeNoOpFile(t, staticPath, "static-v1")
	before, err := svc.noOpInputSnapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(templatePath)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(10 * time.Millisecond)
	writeNoOpFile(t, templatePath, "template-v2")
	if err := os.Chtimes(templatePath, info.ModTime(), info.ModTime()); err != nil {
		t.Fatal(err)
	}
	afterTemplate, err := svc.noOpInputSnapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if before.Fingerprint == afterTemplate.Fingerprint {
		t.Fatal("executable-relative theme template change did not change the no-op fingerprint")
	}
	info, err = os.Stat(staticPath)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(10 * time.Millisecond)
	writeNoOpFile(t, staticPath, "static-v2")
	if err := os.Chtimes(staticPath, info.ModTime(), info.ModTime()); err != nil {
		t.Fatal(err)
	}
	afterStatic, err := svc.noOpInputSnapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if afterTemplate.Fingerprint == afterStatic.Fingerprint {
		t.Fatal("executable-relative theme static change did not change the no-op fingerprint")
	}
}

func TestNoOpFinalCheckDetectsNewExecutableSidecarTheme(t *testing.T) {
	if !fileidentity.Available {
		t.Skip("verified no-op requires Linux file identity")
	}
	svc, _ := newNoOpFixture(t)
	binDir := filepath.Join(t.TempDir(), "bin")
	svc.executable = filepath.Join(binDir, "markata-go")
	writeNoOpFile(t, svc.executable, "builder executable")
	snapshot, err := svc.captureNoOpSnapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	newThemeTemplate := filepath.Join(binDir, "themes", "default", "templates", "base.html")
	writeNoOpFile(t, newThemeTemplate, "new sidecar theme template")
	if _, err := svc.confirmNoOpSnapshot(context.Background(), snapshot); err == nil {
		t.Fatal("final stability check accepted a new executable-sidecar theme root")
	}
}

func TestWriteNoOpIdentityRejectsUnavailableIdentity(t *testing.T) {
	if err := writeNoOpIdentity(sha256.New(), "input-file", "cache.json", noOpFileInfoWithoutIdentity{}); err == nil {
		t.Fatal("no-op identity accepted file metadata without an inode")
	}
}

type noOpFileInfoWithoutIdentity struct{}

func (noOpFileInfoWithoutIdentity) Name() string       { return "cache.json" }
func (noOpFileInfoWithoutIdentity) Size() int64        { return 1 }
func (noOpFileInfoWithoutIdentity) Mode() fs.FileMode  { return 0o644 }
func (noOpFileInfoWithoutIdentity) ModTime() time.Time { return time.Time{} }
func (noOpFileInfoWithoutIdentity) IsDir() bool        { return false }
func (noOpFileInfoWithoutIdentity) Sys() any           { return nil }

func TestNoOpConfigRejectsAbsoluteContentPatterns(t *testing.T) {
	config := models.NewConfig()
	config.GlobConfig.UseGitignore = true
	config.GlobConfig.Patterns = []string{filepath.Join(t.TempDir(), "**", "*.md")}
	if err := validateNoOpConfig(config); err == nil {
		t.Fatal("absolute content glob was accepted for a no-op snapshot")
	}
}

func TestNoOpWellKnownTimeFreshnessGuard(t *testing.T) {
	config := models.NewConfig()
	if !noOpWellKnownTimeEnabled(config) {
		t.Fatal("default well-known time endpoint was not detected")
	}
	config.WellKnown.AutoGenerate = []string{"links"}
	if noOpWellKnownTimeEnabled(config) {
		t.Fatal("well-known time endpoint remained enabled after removal")
	}
	config.WellKnown.AutoGenerate = []string{"time"}
	disabled := false
	config.WellKnown.Enabled = &disabled
	if noOpWellKnownTimeEnabled(config) {
		t.Fatal("disabled well-known plugin was treated as a changing input")
	}
}

func TestNoOpHashesIgnoredDirectoriesWithNonMarkdownInputs(t *testing.T) {
	if !fileidentity.Available {
		t.Skip("verified no-op requires Linux file identity")
	}
	svc, root := newNoOpFixture(t)
	writeNoOpFile(t, filepath.Join(root, ".gitignore"), "cache/\n.cache/\n.markata/\ncustom/\n.env\nignored-media/\n")
	imagePath := filepath.Join(root, "ignored-media", "cover.png")
	writeNoOpFile(t, imagePath, "image-data")
	before, err := svc.noOpInputSnapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(imagePath)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(10 * time.Millisecond)
	writeNoOpFile(t, imagePath, "image-edit")
	if err := os.Chtimes(imagePath, info.ModTime(), info.ModTime()); err != nil {
		t.Fatal(err)
	}
	after, err := svc.noOpInputSnapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if before.Fingerprint == after.Fingerprint {
		t.Fatal("change under a Git-ignored asset directory did not change the no-op fingerprint")
	}
}

func TestNoOpSnapshotRejectsUnsafeInputs(t *testing.T) {
	if !fileidentity.Available {
		t.Skip("verified no-op requires Linux file identity")
	}
	tests := []struct {
		name   string
		before func(t *testing.T, root string)
	}{
		{"custom hook", func(t *testing.T, root string) {
			writeNoOpFile(t, filepath.Join(root, "markata-go.json"), `{"markata-go":{"hooks":["default","python_docs"]}}`)
		}},
		{"unsupported now", func(t *testing.T, root string) {
			writeNoOpFile(t, filepath.Join(root, "templates", "footer.html"), `{{ now }}`)
		}},
		{"dynamic Jinja now", func(t *testing.T, root string) {
			writeNoOpFile(t, filepath.Join(root, "posts", "dynamic.md"), "---\njinja: true\n---\n{{ now | date:\"2006-01-02 15:04\" }}\n")
		}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, root := newNoOpFixture(t)
			tc.before(t, root)
			if _, err := svc.validateNoOpSemantics(context.Background()); err == nil {
				t.Fatal("unsafe input was accepted for no-op")
			}
		})
	}
}

func TestNoOpContentPathMatchesConfiguredPatterns(t *testing.T) {
	config := models.NewConfig()
	config.GlobConfig.Patterns = []string{"pages/**/*.md", "notes/*.markdown"}
	for path, want := range map[string]bool{
		"pages/post.md":         true,
		"pages/nested/post.md":  true,
		"notes/legacy.markdown": true,
		"templates/post.md":     false,
		"pages/post.html":       false,
	} {
		if got := noOpContentPathMatches(config, path); got != want {
			t.Errorf("noOpContentPathMatches(%q)=%t, want %t", path, got, want)
		}
	}
}

func TestNoOpTemplateTimeScanIgnoresMarkdownAuthoringTemplates(t *testing.T) {
	root := t.TempDir()
	writeNoOpFile(t, filepath.Join(root, "templates", "post-templates", "daily.md"), "date: <% tp.date.now(\"YYYY-MM-DD HH:mm:ss\") %>\n")
	writeNoOpFile(t, filepath.Join(root, "templates", "resume.html"), "<!-- show Cat yellow for now -->\n")
	if err := checkTemplateTimeUse(root, "templates"); err != nil {
		t.Fatalf("non-rendered time references disabled no-op eligibility: %v", err)
	}
}

func TestCanUseNoOpRequiresManualUnchangedUnqueuedRequest(t *testing.T) {
	svc := &Service{}
	base := queueRequest{QueuedOperation: QueuedOperation{Kind: "build", TriggerType: "manual-ui"}}
	if !svc.canUseNoOp(base) {
		t.Fatal("manual no-change request was rejected")
	}
	scheduled := queueRequest{QueuedOperation: QueuedOperation{Kind: "build", TriggerType: "scheduled-refresh", Detail: "Build enqueued by refresh task reader-update"}}
	if !svc.canUseNoOp(scheduled) {
		t.Fatal("scheduled refresh no-change request was rejected")
	}
	unknownRefresh := queueRequest{QueuedOperation: QueuedOperation{Kind: "build", TriggerType: "scheduled-refresh", Detail: "Build enqueued by refresh task unknown"}}
	if svc.canUseNoOp(unknownRefresh) {
		t.Fatal("unknown refresh task was accepted for no-op")
	}
	for _, req := range []queueRequest{
		{QueuedOperation: QueuedOperation{Kind: "build", TriggerType: "file-watch"}},
		{QueuedOperation: QueuedOperation{Kind: "build", TriggerType: "schedule"}},
		{QueuedOperation: QueuedOperation{Kind: "build", TriggerType: "manual-ui", Force: true}},
		{QueuedOperation: QueuedOperation{Kind: "build", TriggerType: "manual-ui", Changed: []string{"posts/a.md"}}},
		{QueuedOperation: QueuedOperation{Kind: "build", TriggerType: "manual-ui", SyncSource: true}},
	} {
		if svc.canUseNoOp(req) {
			t.Fatalf("request %+v was accepted for no-op", req)
		}
	}
	svc.state.Queue = []QueuedOperation{{Kind: "refresh", TriggerType: "schedule"}}
	if svc.canUseNoOp(base) {
		t.Fatal("no-op was accepted while a refresh is queued")
	}
}

func TestDecodeBuildAPIRequest(t *testing.T) {
	for _, tc := range []struct {
		name    string
		body    string
		want    bool
		wantErr bool
	}{
		{"legacy empty body", "", false, false},
		{"empty object", `{}`, false, false},
		{"force", `{"force":true}`, true, false},
		{"unknown field", `{"override":true}`, false, true},
		{"non-object", `true`, false, true},
		{"null", `null`, false, true},
		{"multiple objects", `{} {}`, false, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequestWithContext(context.Background(), "POST", "/api/v1/builds", strings.NewReader(tc.body))
			got, err := decodeBuildAPIRequest(httptest.NewRecorder(), req)
			if (err != nil) != tc.wantErr {
				t.Fatalf("decode error=%v, wantErr=%t", err, tc.wantErr)
			}
			if err == nil && got.Force != tc.want {
				t.Fatalf("force=%t, want %t", got.Force, tc.want)
			}
		})
	}
}

func BenchmarkVerifiedNoOpProductionShape(b *testing.B) {
	if !fileidentity.Available {
		b.Skip("verified no-op requires Linux file identity")
	}
	b.StopTimer()
	svc, source := newNoOpFixture(b)
	if err := os.MkdirAll(filepath.Join(source, "posts"), 0o755); err != nil {
		b.Fatal(err)
	}

	for i := 0; i < 3341; i++ {
		path := filepath.Join(source, "posts", "post-"+strconv.Itoa(i)+".md")
		if err := os.WriteFile(path, []byte("---\ntitle: Post\n---\n\nProduction-shaped post content.\n"), 0o644); err != nil {
			b.Fatal(err)
		}
	}
	cacheRoot := filepath.Join(source, ".markata", "image-cache")
	if err := os.MkdirAll(cacheRoot, 0o755); err != nil {
		b.Fatal(err)
	}
	for i := 0; i < 3689; i++ {
		path := filepath.Join(cacheRoot, "cache-"+strconv.Itoa(i)+".json")
		if err := os.WriteFile(path, []byte(`{"cached":true}`), 0o644); err != nil {
			b.Fatal(err)
		}
	}
	for _, args := range [][]string{{"add", "-A"}, {"commit", "-qm", "production-shaped source"}} {
		command := exec.CommandContext(context.Background(), "git", append([]string{"-C", source}, args...)...)
		if output, err := command.CombinedOutput(); err != nil {
			b.Fatalf("git %v: %v: %s", args, err, output)
		}
	}
	releaseID := "release-production-shaped"
	releaseRoot := filepath.Join(svc.cfg.SiteDir, "releases", releaseID)
	for i := 0; i < 37655; i++ {
		dir := filepath.Join(releaseRoot, "pages", "shard-"+strconv.Itoa(i/600))
		if err := os.MkdirAll(dir, 0o755); err != nil {
			b.Fatal(err)
		}
		path := filepath.Join(dir, "page-"+strconv.Itoa(i)+".html")
		if err := os.WriteFile(path, []byte("<!doctype html><title>Page</title>"), 0o644); err != nil {
			b.Fatal(err)
		}
	}
	if err := os.MkdirAll(filepath.Dir(filepath.Join(svc.cfg.SiteDir, "current")), 0o755); err != nil {
		b.Fatal(err)
	}
	if err := os.Symlink(filepath.Join("releases", releaseID), filepath.Join(svc.cfg.SiteDir, "current")); err != nil {
		b.Fatal(err)
	}
	snapshot, err := svc.noOpInputSnapshot(context.Background())
	if err != nil {
		b.Fatal(err)
	}
	if err := svc.saveNoOpBaseline(context.Background(), snapshot, releaseID, time.Now()); err != nil {
		b.Fatal(err)
	}
	if err := os.MkdirAll(svc.cfg.WorkDir, 0o755); err != nil {
		b.Fatal(err)
	}
	if err := markReusableWorkspace(svc.cfg.WorkDir, releaseID); err != nil {
		b.Fatal(err)
	}
	b.Run("source-and-input-snapshot", func(b *testing.B) {
		b.ReportMetric(3341, "source-posts")
		b.ReportMetric(3689, "input-cache-entries")
		for i := 0; i < b.N; i++ {
			if _, err := svc.noOpInputSnapshot(context.Background()); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("live-release-fingerprint", func(b *testing.B) {
		b.ReportMetric(37655, "release-files")
		for i := 0; i < b.N; i++ {
			if _, _, err := svc.noOpReleaseFingerprint(context.Background(), releaseID); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("complete-verified-no-op", func(b *testing.B) {
		b.ReportMetric(3341, "source-posts")
		b.ReportMetric(3689, "input-cache-entries")
		b.ReportMetric(37655, "release-files")
		for i := 0; i < b.N; i++ {
			started := time.Now().UTC()
			record := &BuildRecord{ID: "bench-noop-" + strconv.Itoa(i), Kind: "build", Status: "running", StartedAt: started, EnqueuedAt: started, TriggerType: "manual-ui"}
			req := queueRequest{QueuedOperation: QueuedOperation{ID: record.ID, Kind: "build", TriggerType: "manual-ui", EnqueuedAt: started}}
			ok, _ := svc.verifiedNoOp(context.Background(), req, record, io.Discard)
			if !ok {
				b.Fatal("production-shaped unchanged request did not complete as a no-op")
			}
		}
	})
}

func newNoOpFixture(t testing.TB) (*Service, string) {
	t.Helper()
	root := t.TempDir()
	writeNoOpFile(t, filepath.Join(root, ".gitignore"), "cache/\n.cache/\n.markata/\ncustom/\n.env\n")
	writeNoOpFile(t, filepath.Join(root, "markata-go.json"), `{"markata-go":{"templates_dir":"templates","assets_dir":"static","glob":{"use_gitignore":true},"tailwind":{"build":false},"search":{"backend":"bleve"},"well_known":{"auto_generate":["links"]},"webmentions":{"cache_dir":"custom/webmentions"}}}`)
	writeNoOpFile(t, filepath.Join(root, "templates", "footer.html"), `{{ now | date:"2006" }}`)
	writeNoOpFile(t, filepath.Join(root, "static", "site.css"), "body{}\n")
	for _, args := range [][]string{
		{"init", "-q"},
		{"config", "user.email", "test@example.com"},
		{"config", "user.name", "Test"},
		{"add", "-A"},
		{"commit", "-qm", "fixture"},
	} {
		command := exec.CommandContext(context.Background(), "git", append([]string{"-C", root}, args...)...)
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", args, err, output)
		}
	}
	paths := t.TempDir()
	history := filepath.Join(paths, "history")
	if err := os.MkdirAll(history, 0o755); err != nil {
		t.Fatal(err)
	}
	svc := &Service{
		cfg: Config{
			SourceDir:            root,
			SiteDir:              filepath.Join(paths, "site"),
			WorkDir:              filepath.Join(paths, "work"),
			ConfigPath:           filepath.Join(root, "markata-go.json"),
			BuilderCommit:        "0123456789abcdef0123456789abcdef01234567",
			SuccessfulBuildsKeep: 60,
			FailedBuildsKeep:     20,
		},
		statePath: filepath.Join(history, defaultStateName),
		logDir:    filepath.Join(history, defaultLogDirName),
	}
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	svc.executable = executable
	if err := os.MkdirAll(svc.logDir, 0o755); err != nil {
		t.Fatal(err)
	}
	return svc, root
}

func newNoOpValidationConfig() *models.Config {
	config := models.NewConfig()
	config.WellKnown.AutoGenerate = []string{"links"}
	return config
}

func writeNoOpFile(t testing.TB, path, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}
