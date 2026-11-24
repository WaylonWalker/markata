package assets

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/runtimeenv"
)

func TestNewDownloader(t *testing.T) {
	d := NewDownloader("/tmp/test-cache", true)
	if d.cacheDir != "/tmp/test-cache" {
		t.Errorf("expected cacheDir /tmp/test-cache, got %s", d.cacheDir)
	}
	if !d.verifyIntegrity {
		t.Error("expected verifyIntegrity to be true")
	}
	if d.httpClient == nil {
		t.Error("expected httpClient to be non-nil")
	}
}

func TestRegistry_IncludesWebAwesomeArchive(t *testing.T) {
	asset := GetAsset("webawesome")
	if asset == nil {
		t.Fatal("expected webawesome asset in default registry")
	}
	if asset.Type != "archive" {
		t.Fatalf("webawesome type = %q, want archive", asset.Type)
	}
	if asset.ExtractPath != "package/dist-cdn" {
		t.Fatalf("webawesome ExtractPath = %q, want package/dist-cdn", asset.ExtractPath)
	}
	if asset.OutputPath != "webawesome" {
		t.Fatalf("webawesome OutputPath = %q, want webawesome", asset.OutputPath)
	}
}

func TestDownloader_Download(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("test content")); err != nil {
			t.Logf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	// Create temp cache dir
	cacheDir := t.TempDir()
	d := NewDownloader(cacheDir, false)

	asset := Asset{
		Name:      "test-asset",
		URL:       server.URL + "/test.js",
		LocalPath: "test/test.js",
		Type:      "js",
	}

	ctx := context.Background()
	result, err := d.Download(ctx, asset)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Cached {
		t.Error("expected Cached to be false on first download")
	}
	if result.Size != 12 { // "test content" = 12 bytes
		t.Errorf("expected size 12, got %d", result.Size)
	}
	if result.Error != nil {
		t.Errorf("unexpected result error: %v", result.Error)
	}

	// Verify file exists
	cachedPath := filepath.Join(cacheDir, asset.LocalPath)
	if _, err := os.Stat(cachedPath); os.IsNotExist(err) {
		t.Error("expected cached file to exist")
	}

	// Second download should be cached
	result2, err := d.Download(ctx, asset)
	if err != nil {
		t.Fatalf("unexpected error on second download: %v", err)
	}
	if !result2.Cached {
		t.Error("expected Cached to be true on second download")
	}
}

func TestDownloader_DownloadArchiveAsset(t *testing.T) {
	archiveData := buildTestTarGz(t, map[string]string{
		"package/dist-cdn/webawesome.loader.js":  "loader",
		"package/dist-cdn/styles/webawesome.css": "body {}",
		"package/README.md":                      "ignored",
	})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(archiveData); err != nil {
			t.Logf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	cacheDir := t.TempDir()
	d := NewDownloader(cacheDir, false)
	asset := Asset{
		Name:        "webawesome",
		URL:         server.URL + "/webawesome.tgz",
		LocalPath:   "webawesome",
		Type:        "archive",
		ExtractPath: "package/dist-cdn",
	}

	result, err := d.Download(context.Background(), asset)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Cached {
		t.Fatal("expected first archive download to be uncached")
	}

	loaderPath := filepath.Join(cacheDir, "webawesome", "webawesome.loader.js")
	cssPath := filepath.Join(cacheDir, "webawesome", "styles", "webawesome.css")
	for _, path := range []string{loaderPath, cssPath} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected extracted file %s: %v", path, err)
		}
	}
	if _, err := os.Stat(filepath.Join(cacheDir, "webawesome", "README.md")); !os.IsNotExist(err) {
		t.Fatalf("unexpected non-extracted file copied from outside ExtractPath")
	}

	result2, err := d.Download(context.Background(), asset)
	if err != nil {
		t.Fatalf("unexpected second download error: %v", err)
	}
	if !result2.Cached {
		t.Fatal("expected second archive download to be cached")
	}
}

func TestDownloader_Download_HTTP404(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	cacheDir := t.TempDir()
	d := NewDownloader(cacheDir, false)

	asset := Asset{
		Name:      "missing-asset",
		URL:       server.URL + "/missing.js",
		LocalPath: "missing/missing.js",
		Type:      "js",
	}

	ctx := context.Background()
	result, err := d.Download(ctx, asset)
	if err == nil {
		t.Error("expected error for 404 response")
	}
	if result.Error == nil {
		t.Error("expected result.Error to be set")
	}
}

func TestDownloader_IsCached(t *testing.T) {
	cacheDir := t.TempDir()
	d := NewDownloader(cacheDir, false)

	asset := Asset{
		Name:      "test-asset",
		LocalPath: "test/file.js",
	}

	// Not cached initially
	if d.IsCached(asset) {
		t.Error("expected IsCached to return false initially")
	}

	// Create the cached file
	cachedPath := filepath.Join(cacheDir, asset.LocalPath)
	if err := os.MkdirAll(filepath.Dir(cachedPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cachedPath, []byte("content"), 0o600); err != nil {
		t.Fatal(err)
	}

	// Now should be cached
	if !d.IsCached(asset) {
		t.Error("expected IsCached to return true after file created")
	}
}

func TestDownloader_GetCachedPath(t *testing.T) {
	cacheDir := t.TempDir()
	d := NewDownloader(cacheDir, false)

	asset := Asset{
		Name:      "test-asset",
		LocalPath: "test/file.js",
	}

	// Should return empty string when not cached
	path := d.GetCachedPath(asset)
	if path != "" {
		t.Errorf("expected empty string, got %s", path)
	}

	// Create cached file
	cachedPath := filepath.Join(cacheDir, asset.LocalPath)
	if err := os.MkdirAll(filepath.Dir(cachedPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cachedPath, []byte("content"), 0o600); err != nil {
		t.Fatal(err)
	}

	// Should return path when cached
	path = d.GetCachedPath(asset)
	if path != cachedPath {
		t.Errorf("expected %s, got %s", cachedPath, path)
	}
}

func TestDownloader_CopyToOutput(t *testing.T) {
	cacheDir := t.TempDir()
	outputDir := t.TempDir()
	d := NewDownloader(cacheDir, false)

	asset := Asset{
		Name:      "test-asset",
		LocalPath: "test/file.js",
	}

	// Create cached file
	cachedPath := filepath.Join(cacheDir, asset.LocalPath)
	if err := os.MkdirAll(filepath.Dir(cachedPath), 0o755); err != nil {
		t.Fatal(err)
	}
	content := []byte("test content for copy")
	if err := os.WriteFile(cachedPath, content, 0o600); err != nil {
		t.Fatal(err)
	}

	// Copy to output
	if err := d.CopyToOutput(asset, outputDir); err != nil {
		t.Fatalf("CopyToOutput failed: %v", err)
	}

	// Verify output file
	outputPath := filepath.Join(outputDir, asset.LocalPath)
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	if !bytes.Equal(data, content) {
		t.Errorf("content mismatch: expected %s, got %s", content, data)
	}
}

func TestDownloader_CopyArchiveToOutput(t *testing.T) {
	cacheDir := t.TempDir()
	outputDir := t.TempDir()
	d := NewDownloader(cacheDir, false)
	asset := Asset{
		Name:        "webawesome",
		LocalPath:   "webawesome",
		Type:        "archive",
		ExtractPath: "package/dist-cdn",
	}

	loaderPath := filepath.Join(cacheDir, "webawesome", "webawesome.loader.js")
	if err := os.MkdirAll(filepath.Dir(loaderPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(loaderPath, []byte("loader"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cacheDir, "webawesome.complete"), []byte("3.5.0"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := d.CopyToOutput(asset, outputDir); err != nil {
		t.Fatalf("CopyToOutput failed: %v", err)
	}

	outputPath := filepath.Join(outputDir, "webawesome", "webawesome.loader.js")
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read copied archive file: %v", err)
	}
	if string(data) != "loader" {
		t.Fatalf("unexpected output content: %s", data)
	}
}

func TestDownloader_CopyToOutput_NotCached(t *testing.T) {
	cacheDir := t.TempDir()
	outputDir := t.TempDir()
	d := NewDownloader(cacheDir, false)

	asset := Asset{
		Name:      "uncached-asset",
		LocalPath: "uncached/file.js",
	}

	err := d.CopyToOutput(asset, outputDir)
	if err == nil {
		t.Error("expected error when copying uncached asset")
	}
}

func TestDownloader_Download_UsesBundledFallbackCacheWhenOffline(t *testing.T) {
	cacheDir := t.TempDir()
	bundledDir := t.TempDir()
	t.Setenv(runtimeenv.EnvOffline, "true")
	t.Setenv(runtimeenv.EnvBundledAssetsCacheDir, bundledDir)

	asset := Asset{
		Name:      "bundled-asset",
		URL:       "https://example.invalid/bundled.js",
		LocalPath: "bundled/file.js",
		Type:      "js",
	}
	bundledPath := filepath.Join(bundledDir, asset.LocalPath)
	if err := os.MkdirAll(filepath.Dir(bundledPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(bundledPath, []byte("bundled content"), 0o600); err != nil {
		t.Fatal(err)
	}

	d := NewDownloader(cacheDir, false)
	result, err := d.Download(context.Background(), asset)
	if err != nil {
		t.Fatalf("unexpected error loading bundled cached asset: %v", err)
	}
	if !result.Cached {
		t.Fatal("expected bundled asset to be treated as cached")
	}
	if got := d.GetCachedPath(asset); got != bundledPath {
		t.Fatalf("GetCachedPath() = %q, want %q", got, bundledPath)
	}
}

func TestDownloader_Download_OfflineMissingFails(t *testing.T) {
	cacheDir := t.TempDir()
	t.Setenv(runtimeenv.EnvOffline, "true")
	t.Setenv(runtimeenv.EnvBundledAssetsCacheDir, t.TempDir())

	d := NewDownloader(cacheDir, false)
	asset := Asset{
		Name:      "missing-asset",
		URL:       "https://example.invalid/missing.js",
		LocalPath: "missing/file.js",
		Type:      "js",
	}

	_, err := d.Download(context.Background(), asset)
	if err == nil {
		t.Fatal("expected offline download to fail when asset is not cached")
	}
}

func TestDownloader_Clean(t *testing.T) {
	cacheDir := t.TempDir()
	d := NewDownloader(cacheDir, false)

	// Create some files
	testFile := filepath.Join(cacheDir, "test", "file.js")
	if err := os.MkdirAll(filepath.Dir(testFile), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(testFile, []byte("content"), 0o600); err != nil {
		t.Fatal(err)
	}

	// Clean
	if err := d.Clean(); err != nil {
		t.Fatalf("Clean failed: %v", err)
	}

	// Verify cache dir is removed
	if _, err := os.Stat(cacheDir); !os.IsNotExist(err) {
		t.Error("expected cache dir to be removed")
	}
}

func TestDownloader_Status(t *testing.T) {
	cacheDir := t.TempDir()
	d := NewDownloader(cacheDir, false)

	statuses := d.Status()
	if len(statuses) == 0 {
		t.Error("expected non-empty status list")
	}

	// All should be uncached initially
	for _, s := range statuses {
		if s.Cached {
			t.Errorf("expected asset %s to be uncached", s.Asset.Name)
		}
	}
}

func TestDownloader_DownloadAll(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("test content")); err != nil {
			t.Logf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	cacheDir := t.TempDir()
	d := NewDownloader(cacheDir, false)

	// Override the registry temporarily (we can't do this in current impl)
	// Instead, just test that DownloadAll runs without panic
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// This will try to download real assets which may fail in test environment
	// Just verify it doesn't panic
	_ = d.DownloadAll(ctx, 2)
}

func TestVerifyIntegrity(t *testing.T) {
	tests := []struct {
		name      string
		data      []byte
		integrity string
		wantErr   bool
	}{
		{
			name:      "valid sha256",
			data:      []byte("hello"),
			integrity: "sha256-LPJNul+wow4m6DsqxbninhsWHlwfp0JecwQzYpOLmCQ=",
			wantErr:   false,
		},
		{
			name:      "invalid sha256",
			data:      []byte("hello"),
			integrity: "sha256-invalidhash",
			wantErr:   true,
		},
		{
			name:      "invalid format",
			data:      []byte("hello"),
			integrity: "nohyphen",
			wantErr:   true,
		},
		{
			name:      "unsupported algorithm",
			data:      []byte("hello"),
			integrity: "md5-rL0Y20zC+Fzt72VPzMSk2A==",
			wantErr:   true,
		},
		{
			name:      "valid sha384",
			data:      []byte("hello"),
			integrity: "sha384-WeF0h3dEjGnea4ANejO7+5/xtGPkQ1TDVTvNucZm+pASWjx5+QOXvfX2oT3oKGhP",
			wantErr:   false,
		},
		{
			name:      "valid sha512",
			data:      []byte("hello"),
			integrity: "sha512-m3HSJL1i83hdltRq0+o9czGb+8KJDKra4t/3JRlnPKcjI8PZm6XBHXx6zG4UuMXaDEZjR1wuXDre9G9zvN7AQw==",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := verifyIntegrity(tt.data, tt.integrity)
			if (err != nil) != tt.wantErr {
				t.Errorf("verifyIntegrity() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func buildTestTarGz(t *testing.T, files map[string]string) []byte {
	t.Helper()

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)

	for name, content := range files {
		header := &tar.Header{
			Name: name,
			Mode: 0o644,
			Size: int64(len(content)),
		}
		if err := tw.WriteHeader(header); err != nil {
			t.Fatalf("write tar header: %v", err)
		}
		if _, err := tw.Write([]byte(content)); err != nil {
			t.Fatalf("write tar content: %v", err)
		}
	}

	if err := tw.Close(); err != nil {
		t.Fatalf("close tar writer: %v", err)
	}
	if err := gz.Close(); err != nil {
		t.Fatalf("close gzip writer: %v", err)
	}

	return buf.Bytes()
}
