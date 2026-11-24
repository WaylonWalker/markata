// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildstats"
	"github.com/WaylonWalker/markata/pkg/logging"
)

var pagefindInstallLog = logging.Component("pagefind_installer").Phase("cleanup")

// Constants for Pagefind installer.
const (
	// pagefindBinaryName is the name of the Pagefind binary.
	pagefindBinaryName = "pagefind"

	// pagefindBinaryNameWindows is the name of the Pagefind binary on Windows.
	pagefindBinaryNameWindows = "pagefind.exe"

	// osWindows is the GOOS value for Windows.
	osWindows = "windows"

	// defaultReleaseBaseURL is the base URL for Pagefind releases on GitHub.
	// Note: Pagefind moved from CloudCannon/pagefind to Pagefind/pagefind
	defaultReleaseBaseURL = "https://github.com/Pagefind/pagefind/releases"

	// defaultHTTPTimeout is the default timeout for HTTP requests.
	defaultHTTPTimeout = 120 * time.Second

	// maxBinarySize is the maximum allowed binary size (50MB) to prevent decompression bombs.
	maxBinarySize = 50 * 1024 * 1024
)

// PagefindInstaller handles automatic downloading and caching of Pagefind binaries.
type PagefindInstaller struct {
	// CacheDir is the directory where binaries are cached.
	CacheDir string

	// Version is the Pagefind version to install (e.g., "v1.4.0" or "latest").
	Version string

	// Verbose enables verbose output during installation.
	Verbose bool

	// client is the HTTP client used for downloads.
	client *http.Client
}

// PagefindInstallError indicates an error during Pagefind installation.
type PagefindInstallError struct {
	Operation string
	Message   string
	Err       error
}

func (e *PagefindInstallError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("pagefind install error during %s: %s: %v", e.Operation, e.Message, e.Err)
	}
	return fmt.Sprintf("pagefind install error during %s: %s", e.Operation, e.Message)
}

func (e *PagefindInstallError) Unwrap() error {
	return e.Err
}

// NewPagefindInstallError creates a new PagefindInstallError.
func NewPagefindInstallError(operation, message string, err error) *PagefindInstallError {
	return &PagefindInstallError{
		Operation: operation,
		Message:   message,
		Err:       err,
	}
}

// PagefindInstallerConfig configures the Pagefind installer.
type PagefindInstallerConfig struct {
	// AutoInstall enables automatic Pagefind installation (default: true).
	AutoInstall *bool `json:"auto_install,omitempty" yaml:"auto_install,omitempty" toml:"auto_install,omitempty"`

	// Version is the Pagefind version to install (default: "latest").
	Version string `json:"version,omitempty" yaml:"version,omitempty" toml:"version,omitempty"`

	// CacheDir is the directory for caching binaries (default: ~/.markata-go/bin/).
	CacheDir string `json:"cache_dir,omitempty" yaml:"cache_dir,omitempty" toml:"cache_dir,omitempty"`
}

// IsAutoInstallEnabled returns whether auto-install is enabled.
// Defaults to true if not explicitly set.
func (c *PagefindInstallerConfig) IsAutoInstallEnabled() bool {
	if c.AutoInstall == nil {
		return true
	}
	return *c.AutoInstall
}

// NewPagefindInstallerConfig creates a new PagefindInstallerConfig with default values.
func NewPagefindInstallerConfig() PagefindInstallerConfig {
	autoInstall := true
	return PagefindInstallerConfig{
		AutoInstall: &autoInstall,
		Version:     Latest,
		CacheDir:    "",
	}
}

// platformMapping maps Go's GOOS/GOARCH to Pagefind's asset naming convention.
var platformMapping = map[string]map[string]string{
	"darwin": {
		"amd64": "x86_64-apple-darwin",
		"arm64": "aarch64-apple-darwin",
	},
	"linux": {
		"amd64": "x86_64-unknown-linux-musl",
		"arm64": "aarch64-unknown-linux-musl",
	},
	osWindows: {
		"amd64": "x86_64-pc-windows-msvc",
	},
	"freebsd": {
		"amd64": "x86_64-unknown-freebsd",
	},
}

// NewPagefindInstaller creates a new PagefindInstaller with default settings.
func NewPagefindInstaller() *PagefindInstaller {
	client := &http.Client{
		Timeout: defaultHTTPTimeout,
	}
	buildstats.InstrumentHTTPClient(client)

	return &PagefindInstaller{
		CacheDir: "",
		Version:  Latest,
		client:   client,
	}
}

// NewPagefindInstallerWithConfig creates a new PagefindInstaller from config.
func NewPagefindInstallerWithConfig(config PagefindInstallerConfig) *PagefindInstaller {
	installer := NewPagefindInstaller()

	if config.CacheDir != "" {
		installer.CacheDir = expandPath(config.CacheDir)
	}
	if config.Version != "" {
		installer.Version = config.Version
	}

	return installer
}

// GetPlatformAssetName returns the Pagefind asset name for the current platform.
func GetPlatformAssetName() (string, error) {
	return getPlatformAssetNameForOS(runtime.GOOS, runtime.GOARCH)
}

// getPlatformAssetNameForOS returns the asset name for a given OS/arch combination.
func getPlatformAssetNameForOS(goos, goarch string) (string, error) {
	archMap, ok := platformMapping[goos]
	if !ok {
		return "", NewPagefindInstallError(
			"platform_detection",
			fmt.Sprintf("unsupported operating system: %s", goos),
			nil,
		)
	}

	assetName, ok := archMap[goarch]
	if !ok {
		return "", NewPagefindInstallError(
			"platform_detection",
			fmt.Sprintf("unsupported architecture %s on %s", goarch, goos),
			nil,
		)
	}

	return assetName, nil
}

// GetCacheDir returns the cache directory, creating it if necessary.
func (i *PagefindInstaller) GetCacheDir() (string, error) {
	if i.CacheDir != "" {
		cacheDir := expandPath(i.CacheDir)
		if err := os.MkdirAll(cacheDir, 0o755); err != nil {
			return "", NewPagefindInstallError("cache_setup", "failed to create cache directory", err)
		}
		return cacheDir, nil
	}

	// Use XDG_CACHE_HOME or fallback to ~/.cache
	cacheBase := os.Getenv("XDG_CACHE_HOME")
	if cacheBase == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", NewPagefindInstallError("cache_setup", "failed to get user home directory", err)
		}
		cacheBase = filepath.Join(homeDir, ".cache")
	}

	cacheDir := filepath.Join(cacheBase, "markata-go", "bin")
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return "", NewPagefindInstallError("cache_setup", "failed to create cache directory", err)
	}

	return cacheDir, nil
}

// getBinaryName returns the appropriate binary name for the current OS.
func getBinaryName() string {
	if runtime.GOOS == osWindows {
		return pagefindBinaryNameWindows
	}
	return pagefindBinaryName
}

// GetCachedBinaryPath returns the path to a cached Pagefind binary for the given version.
func (i *PagefindInstaller) GetCachedBinaryPath(version string) (string, error) {
	cacheDir, err := i.GetCacheDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(cacheDir, version, getBinaryName()), nil
}

// IsCached checks if a specific version is already cached and valid.
func (i *PagefindInstaller) IsCached(version string) (bool, error) {
	binaryPath, err := i.GetCachedBinaryPath(version)
	if err != nil {
		return false, err
	}

	info, err := os.Stat(binaryPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, NewPagefindInstallError("cache_check", "failed to check cached binary", err)
	}

	// Check that it's executable (on Unix) and not empty
	if info.Size() == 0 {
		return false, nil
	}

	return true, nil
}

// GetLatestVersion fetches the latest Pagefind version from GitHub releases.
// It follows redirect chains (e.g., repo renames) until it finds the final URL with the version tag.
func (i *PagefindInstaller) GetLatestVersion() (string, error) {
	v, err := getLatestGitHubReleaseVersion(defaultReleaseBaseURL)
	if err != nil {
		return "", NewPagefindInstallError("version_check", "failed to check latest version", err)
	}
	if v == "" {
		return "", NewPagefindInstallError("version_check", "could not find version in final URL or too many redirects", nil)
	}
	return v, nil
}

// ResolveVersion resolves "latest" to an actual version number.
func (i *PagefindInstaller) ResolveVersion() (string, error) {
	if i.Version != Latest && i.Version != "" {
		return i.Version, nil
	}

	return i.GetLatestVersion()
}

// buildAssetURL constructs the download URL for a Pagefind release asset.
func buildAssetURL(version, platformAsset string) string {
	// Format: pagefind-v1.4.0-x86_64-apple-darwin.tar.gz
	filename := fmt.Sprintf("pagefind-%s-%s.tar.gz", version, platformAsset)
	return fmt.Sprintf("%s/download/%s/%s", defaultReleaseBaseURL, version, filename)
}

// buildChecksumURL constructs the URL for the SHA256 checksum file.
func buildChecksumURL(version, platformAsset string) string {
	filename := fmt.Sprintf("pagefind-%s-%s.tar.gz.sha256", version, platformAsset)
	return fmt.Sprintf("%s/download/%s/%s", defaultReleaseBaseURL, version, filename)
}

// fetchChecksum downloads and parses the SHA256 checksum for verification.
func (i *PagefindInstaller) fetchChecksum(version, platformAsset string) (string, error) {
	checksumURL := buildChecksumURL(version, platformAsset)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", checksumURL, http.NoBody)
	if err != nil {
		return "", NewPagefindInstallError("checksum_fetch", "failed to create request", err)
	}

	resp, err := i.client.Do(req)
	if err != nil {
		return "", NewPagefindInstallError("checksum_fetch", "failed to download checksum", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", NewPagefindInstallError(
			"checksum_fetch",
			fmt.Sprintf("checksum download failed with status %d", resp.StatusCode),
			nil,
		)
	}

	// Read the checksum file content
	content, err := io.ReadAll(io.LimitReader(resp.Body, 1024))
	if err != nil {
		return "", NewPagefindInstallError("checksum_fetch", "failed to read checksum", err)
	}

	// Parse checksum (format: "hash  filename" or just "hash")
	checksumStr := strings.TrimSpace(string(content))
	parts := strings.Fields(checksumStr)
	if len(parts) == 0 {
		return "", NewPagefindInstallError("checksum_fetch", "empty checksum file", nil)
	}

	checksum := parts[0]
	if len(checksum) != 64 {
		return "", NewPagefindInstallError(
			"checksum_fetch",
			fmt.Sprintf("invalid checksum length: expected 64, got %d", len(checksum)),
			nil,
		)
	}

	return strings.ToLower(checksum), nil
}

// downloadAsset downloads the Pagefind release archive to a temporary file.
func (i *PagefindInstaller) downloadAsset(version, platformAsset string) (string, error) {
	assetURL := buildAssetURL(version, platformAsset)

	if i.Verbose {
		pagefindInstallLog.Printf("Downloading Pagefind %s for %s...", version, platformAsset)
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultHTTPTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", assetURL, http.NoBody)
	if err != nil {
		return "", NewPagefindInstallError("download", "failed to create request", err)
	}

	resp, err := i.client.Do(req)
	if err != nil {
		return "", NewPagefindInstallError("download", "failed to download asset", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", NewPagefindInstallError(
			"download",
			fmt.Sprintf("download failed with status %d", resp.StatusCode),
			nil,
		)
	}

	// Create a temporary file for the download
	tmpFile, err := os.CreateTemp("", "pagefind-*.tar.gz")
	if err != nil {
		return "", NewPagefindInstallError("download", "failed to create temp file", err)
	}
	defer tmpFile.Close()

	// Download with progress (if terminal)
	written, err := io.Copy(tmpFile, resp.Body)
	if err != nil {
		os.Remove(tmpFile.Name())
		return "", NewPagefindInstallError("download", "failed to write downloaded file", err)
	}

	if i.Verbose {
		pagefindInstallLog.Printf("Downloaded %d bytes", written)
	}

	return tmpFile.Name(), nil
}

// verifyChecksum verifies the SHA256 checksum of a downloaded file.
func verifyChecksum(filePath, expectedChecksum string, verbose bool) error {
	file, err := os.Open(filePath)
	if err != nil {
		return NewPagefindInstallError("verify", "failed to open file for verification", err)
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return NewPagefindInstallError("verify", "failed to compute checksum", err)
	}

	actualChecksum := hex.EncodeToString(hasher.Sum(nil))

	if actualChecksum != expectedChecksum {
		return NewPagefindInstallError(
			"verify",
			fmt.Sprintf("checksum mismatch: expected %s, got %s", expectedChecksum, actualChecksum),
			nil,
		)
	}

	if verbose {
		pagefindInstallLog.Printf("Checksum verified: %s", actualChecksum[:16]+"...")
	}
	return nil
}

// extractBinary extracts the pagefind binary from the downloaded tar.gz archive.
func (i *PagefindInstaller) extractBinary(archivePath, version string) (string, error) {
	cacheDir, err := i.GetCacheDir()
	if err != nil {
		return "", err
	}

	// Create version-specific directory
	versionDir := filepath.Join(cacheDir, version)
	if err := os.MkdirAll(versionDir, 0o755); err != nil {
		return "", NewPagefindInstallError("extract", "failed to create version directory", err)
	}

	// Open the archive
	file, err := os.Open(archivePath)
	if err != nil {
		return "", NewPagefindInstallError("extract", "failed to open archive", err)
	}
	defer file.Close()

	// Create gzip reader
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return "", NewPagefindInstallError("extract", "failed to create gzip reader", err)
	}
	defer gzipReader.Close()

	// Create tar reader
	tarReader := tar.NewReader(gzipReader)

	binaryName := getBinaryName()

	var extractedPath string

	// Extract files
	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", NewPagefindInstallError("extract", "failed to read tar header", err)
		}

		// Only extract the pagefind binary
		baseName := filepath.Base(header.Name)
		if baseName != binaryName {
			continue
		}

		extractedPath = filepath.Join(versionDir, binaryName)

		// Create the output file
		outFile, err := os.OpenFile(extractedPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
		if err != nil {
			return "", NewPagefindInstallError("extract", "failed to create output file", err)
		}

		// Copy the binary content with size limit to prevent decompression bombs
		written, err := io.Copy(outFile, io.LimitReader(tarReader, maxBinarySize))
		outFile.Close()
		if err != nil {
			os.Remove(extractedPath)
			return "", NewPagefindInstallError("extract", "failed to extract binary", err)
		}

		if written == maxBinarySize {
			os.Remove(extractedPath)
			return "", NewPagefindInstallError("extract", "binary exceeds maximum allowed size", nil)
		}

		if i.Verbose {
			pagefindInstallLog.Printf("Extracted %s to %s", binaryName, versionDir)
		}
		break
	}

	if extractedPath == "" {
		return "", NewPagefindInstallError("extract", fmt.Sprintf("binary '%s' not found in archive", binaryName), nil)
	}

	return extractedPath, nil
}

// Install downloads and caches the Pagefind binary for the current platform.
// Returns the path to the installed binary.
func (i *PagefindInstaller) Install() (string, error) {
	// Resolve version
	version, err := i.ResolveVersion()
	if err != nil {
		return "", err
	}

	// Check if already cached
	cached, err := i.IsCached(version)
	if err != nil {
		return "", err
	}

	if cached {
		binaryPath, err := i.GetCachedBinaryPath(version)
		if err != nil {
			return "", err
		}
		if i.Verbose {
			pagefindInstallLog.Printf("Using cached Pagefind %s", version)
		}
		return binaryPath, nil
	}

	// Get platform asset name
	platformAsset, err := GetPlatformAssetName()
	if err != nil {
		return "", err
	}

	// Fetch expected checksum
	expectedChecksum, err := i.fetchChecksum(version, platformAsset)
	if err != nil {
		return "", err
	}

	// Download the archive
	archivePath, err := i.downloadAsset(version, platformAsset)
	if err != nil {
		return "", err
	}
	defer os.Remove(archivePath)

	// Verify checksum (CRITICAL for security)
	if err := verifyChecksum(archivePath, expectedChecksum, i.Verbose); err != nil {
		return "", err
	}

	// Extract the binary
	binaryPath, err := i.extractBinary(archivePath, version)
	if err != nil {
		return "", err
	}

	if i.Verbose {
		pagefindInstallLog.Printf("Successfully installed Pagefind %s", version)
	}
	return binaryPath, nil
}

// expandPath expands ~ to the user's home directory.
func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(homeDir, path[2:])
		}
	}
	return path
}
