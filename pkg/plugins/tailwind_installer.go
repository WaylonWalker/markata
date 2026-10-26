// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
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

var tailwindInstallLog = logging.Component("tailwind_installer").Phase("cleanup")

const (
	tailwindBinaryName        = "tailwindcss"
	tailwindBinaryNameWindows = "tailwindcss.exe"
	tailwindReleaseBaseURL    = "https://github.com/tailwindlabs/tailwindcss/releases"
	tailwindHTTPTimeout       = 120 * time.Second
	tailwindMaxBinarySize     = 200 * 1024 * 1024
	tailwindOSWindows         = "windows"
)

// TailwindInstaller handles automatic downloading and caching of Tailwind CLI binaries.
type TailwindInstaller struct {
	CacheDir string
	Version  string
	Verbose  bool
	client   *http.Client
}

// TailwindInstallError indicates an error during Tailwind installation.
type TailwindInstallError struct {
	Operation string
	Message   string
	Err       error
}

func (e *TailwindInstallError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("tailwind install error during %s: %s: %v", e.Operation, e.Message, e.Err)
	}
	return fmt.Sprintf("tailwind install error during %s: %s", e.Operation, e.Message)
}

func (e *TailwindInstallError) Unwrap() error {
	return e.Err
}

// NewTailwindInstallError creates a new TailwindInstallError.
func NewTailwindInstallError(operation, message string, err error) *TailwindInstallError {
	return &TailwindInstallError{
		Operation: operation,
		Message:   message,
		Err:       err,
	}
}

// TailwindInstallerConfig configures the Tailwind installer.
type TailwindInstallerConfig struct {
	AutoInstall *bool  `json:"auto_install,omitempty" yaml:"auto_install,omitempty" toml:"auto_install,omitempty"`
	Version     string `json:"version,omitempty" yaml:"version,omitempty" toml:"version,omitempty"`
	CacheDir    string `json:"cache_dir,omitempty" yaml:"cache_dir,omitempty" toml:"cache_dir,omitempty"`
	Verbose     *bool  `json:"verbose,omitempty" yaml:"verbose,omitempty" toml:"verbose,omitempty"`
}

// IsAutoInstallEnabled returns whether auto-install is enabled.
// Defaults to true if not explicitly set.
func (c *TailwindInstallerConfig) IsAutoInstallEnabled() bool {
	if c.AutoInstall == nil {
		return true
	}
	return *c.AutoInstall
}

// NewTailwindInstallerConfig creates a new TailwindInstallerConfig with default values.
func NewTailwindInstallerConfig() TailwindInstallerConfig {
	autoInstall := true
	return TailwindInstallerConfig{
		AutoInstall: &autoInstall,
		Version:     Latest,
		CacheDir:    "",
	}
}

// tailwindPlatformMapping maps Go's GOOS/GOARCH to Tailwind CLI asset naming.
var tailwindPlatformMapping = map[string]map[string]string{
	"darwin": {
		"amd64": "macos-x64",
		"arm64": "macos-arm64",
	},
	"linux": {
		"amd64": "linux-x64",
		"arm64": "linux-arm64",
	},
	tailwindOSWindows: {
		"amd64": "windows-x64",
	},
}

// NewTailwindInstaller creates a new TailwindInstaller with default settings.
func NewTailwindInstaller() *TailwindInstaller {
	client := &http.Client{
		Timeout: tailwindHTTPTimeout,
	}
	buildstats.InstrumentHTTPClient(client)

	return &TailwindInstaller{
		CacheDir: "",
		Version:  Latest,
		client:   client,
	}
}

// NewTailwindInstallerWithConfig creates a new TailwindInstaller from config.
func NewTailwindInstallerWithConfig(config TailwindInstallerConfig) *TailwindInstaller {
	installer := NewTailwindInstaller()
	if config.CacheDir != "" {
		installer.CacheDir = expandPath(config.CacheDir)
	}
	if config.Version != "" {
		installer.Version = config.Version
	}
	if config.Verbose != nil {
		installer.Verbose = *config.Verbose
	}
	return installer
}

// GetTailwindPlatformAssetName returns the Tailwind asset name for the current platform.
func GetTailwindPlatformAssetName() (string, error) {
	return getTailwindPlatformAssetNameForOS(runtime.GOOS, runtime.GOARCH)
}

// getTailwindPlatformAssetNameForOS returns the asset name for a given OS/arch combination.
func getTailwindPlatformAssetNameForOS(goos, goarch string) (string, error) {
	archMap, ok := tailwindPlatformMapping[goos]
	if !ok {
		return "", NewTailwindInstallError(
			"platform_detection",
			fmt.Sprintf("unsupported operating system: %s", goos),
			nil,
		)
	}

	assetName, ok := archMap[goarch]
	if !ok {
		return "", NewTailwindInstallError(
			"platform_detection",
			fmt.Sprintf("unsupported architecture %s on %s", goarch, goos),
			nil,
		)
	}

	return assetName, nil
}

// GetTailwindCacheDir returns the cache directory, creating it if necessary.
func (i *TailwindInstaller) GetTailwindCacheDir() (string, error) {
	if i.CacheDir != "" {
		cacheDir := expandPath(i.CacheDir)
		if err := os.MkdirAll(cacheDir, 0o755); err != nil {
			return "", NewTailwindInstallError("cache_setup", "failed to create cache directory", err)
		}
		return cacheDir, nil
	}

	cacheBase := os.Getenv("XDG_CACHE_HOME")
	if cacheBase == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", NewTailwindInstallError("cache_setup", "failed to get user home directory", err)
		}
		cacheBase = filepath.Join(homeDir, ".cache")
	}

	cacheDir := filepath.Join(cacheBase, "markata-go", "bin")
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return "", NewTailwindInstallError("cache_setup", "failed to create cache directory", err)
	}

	return cacheDir, nil
}

// getTailwindBinaryName returns the appropriate binary name for the current OS.
func getTailwindBinaryName() string {
	if runtime.GOOS == tailwindOSWindows {
		return tailwindBinaryNameWindows
	}
	return tailwindBinaryName
}

// GetCachedTailwindBinaryPath returns the path to a cached Tailwind binary for the given version.
func (i *TailwindInstaller) GetCachedTailwindBinaryPath(version string) (string, error) {
	cacheDir, err := i.GetTailwindCacheDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(cacheDir, "tailwindcss", version, getTailwindBinaryName()), nil
}

// IsTailwindCached checks if a specific version is already cached and valid.
func (i *TailwindInstaller) IsTailwindCached(version string) (bool, error) {
	binaryPath, err := i.GetCachedTailwindBinaryPath(version)
	if err != nil {
		return false, err
	}

	info, err := os.Stat(binaryPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, NewTailwindInstallError("cache_check", "failed to check cached binary", err)
	}

	if info.Size() == 0 {
		return false, nil
	}

	return true, nil
}

// GetLatestTailwindVersion fetches the latest Tailwind version from GitHub releases.
func (i *TailwindInstaller) GetLatestTailwindVersion() (string, error) {
	v, err := getLatestGitHubReleaseVersion(tailwindReleaseBaseURL)
	if err != nil {
		return "", NewTailwindInstallError("version_check", "failed to check latest version", err)
	}
	if v == "" {
		return "", NewTailwindInstallError("version_check", "could not find version in final URL or too many redirects", nil)
	}
	return v, nil
}

// ResolveTailwindVersion resolves "latest" to an actual version number.
func (i *TailwindInstaller) ResolveTailwindVersion() (string, error) {
	if i.Version != Latest && i.Version != "" {
		return i.Version, nil
	}
	return i.GetLatestTailwindVersion()
}

// buildTailwindAssetURL constructs the download URL for a Tailwind release asset.
func buildTailwindAssetURL(version, platformAsset string) string {
	filename := fmt.Sprintf("tailwindcss-%s", platformAsset)
	if strings.HasPrefix(platformAsset, "windows-") {
		filename += ".exe"
	}
	return fmt.Sprintf("%s/download/%s/%s", tailwindReleaseBaseURL, version, filename)
}

// buildTailwindChecksumsURL constructs the URL for the SHA256 sums file.
func buildTailwindChecksumsURL(version string) string {
	return fmt.Sprintf("%s/download/%s/sha256sums.txt", tailwindReleaseBaseURL, version)
}

// fetchTailwindChecksum downloads and parses the SHA256 checksum for verification.
func (i *TailwindInstaller) fetchTailwindChecksum(version, platformAsset string) (string, error) {
	checksumURL := buildTailwindChecksumsURL(version)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", checksumURL, http.NoBody)
	if err != nil {
		return "", NewTailwindInstallError("checksum_fetch", "failed to create request", err)
	}

	resp, err := i.client.Do(req)
	if err != nil {
		return "", NewTailwindInstallError("checksum_fetch", "failed to download checksum", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", NewTailwindInstallError(
			"checksum_fetch",
			fmt.Sprintf("checksum download failed with status %d", resp.StatusCode),
			nil,
		)
	}

	content, err := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
	if err != nil {
		return "", NewTailwindInstallError("checksum_fetch", "failed to read checksum", err)
	}

	targetName := fmt.Sprintf("tailwindcss-%s", platformAsset)
	if strings.HasPrefix(platformAsset, "windows-") {
		targetName += ".exe"
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) < 2 {
			continue
		}
		filename := strings.TrimPrefix(fields[1], "./")
		if filename == targetName {
			checksum := strings.ToLower(fields[0])
			if len(checksum) != 64 {
				return "", NewTailwindInstallError(
					"checksum_fetch",
					fmt.Sprintf("invalid checksum length: expected 64, got %d", len(checksum)),
					nil,
				)
			}
			return checksum, nil
		}
	}

	return "", NewTailwindInstallError("checksum_fetch", "checksum entry not found", nil)
}

// downloadTailwindAsset downloads the Tailwind binary to a temporary file.
func (i *TailwindInstaller) downloadTailwindAsset(version, platformAsset string) (string, error) {
	assetURL := buildTailwindAssetURL(version, platformAsset)

	if i.Verbose {
		tailwindInstallLog.Printf("Downloading Tailwind %s for %s...", version, platformAsset)
	}

	ctx, cancel := context.WithTimeout(context.Background(), tailwindHTTPTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", assetURL, http.NoBody)
	if err != nil {
		return "", NewTailwindInstallError("download", "failed to create request", err)
	}

	resp, err := i.client.Do(req)
	if err != nil {
		return "", NewTailwindInstallError("download", "failed to download asset", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", NewTailwindInstallError(
			"download",
			fmt.Sprintf("download failed with status %d", resp.StatusCode),
			nil,
		)
	}

	tmpFile, err := os.CreateTemp("", "tailwind-*")
	if err != nil {
		return "", NewTailwindInstallError("download", "failed to create temp file", err)
	}
	defer tmpFile.Close()

	written, err := io.Copy(tmpFile, io.LimitReader(resp.Body, tailwindMaxBinarySize+1))
	if err != nil {
		os.Remove(tmpFile.Name())
		return "", NewTailwindInstallError("download", "failed to write downloaded file", err)
	}

	if written > tailwindMaxBinarySize {
		os.Remove(tmpFile.Name())
		return "", NewTailwindInstallError("download", "binary exceeds maximum allowed size", nil)
	}

	if i.Verbose {
		tailwindInstallLog.Printf("Downloaded %d bytes", written)
	}

	return tmpFile.Name(), nil
}

// verifyTailwindChecksum verifies the SHA256 checksum of a downloaded file.
func verifyTailwindChecksum(filePath, expectedChecksum string, verbose bool) error {
	file, err := os.Open(filePath)
	if err != nil {
		return NewTailwindInstallError("verify", "failed to open file for verification", err)
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return NewTailwindInstallError("verify", "failed to compute checksum", err)
	}

	actualChecksum := hex.EncodeToString(hasher.Sum(nil))
	if actualChecksum != expectedChecksum {
		return NewTailwindInstallError(
			"verify",
			fmt.Sprintf("checksum mismatch: expected %s, got %s", expectedChecksum, actualChecksum),
			nil,
		)
	}

	if verbose {
		tailwindInstallLog.Printf("Checksum verified: %s", actualChecksum[:16]+"...")
	}

	return nil
}

// installTailwindBinary installs the Tailwind binary into the cache dir.
func (i *TailwindInstaller) installTailwindBinary(downloadPath, version string) (string, error) {
	cacheDir, err := i.GetTailwindCacheDir()
	if err != nil {
		return "", err
	}

	versionDir := filepath.Join(cacheDir, "tailwindcss", version)
	if err := os.MkdirAll(versionDir, 0o755); err != nil {
		return "", NewTailwindInstallError("install", "failed to create version directory", err)
	}

	binaryName := getTailwindBinaryName()
	destPath := filepath.Join(versionDir, binaryName)

	input, err := os.Open(downloadPath)
	if err != nil {
		return "", NewTailwindInstallError("install", "failed to open downloaded binary", err)
	}
	defer input.Close()

	output, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return "", NewTailwindInstallError("install", "failed to create binary file", err)
	}

	if _, err := io.Copy(output, io.LimitReader(input, tailwindMaxBinarySize)); err != nil {
		output.Close()
		os.Remove(destPath)
		return "", NewTailwindInstallError("install", "failed to write binary", err)
	}
	if err := output.Close(); err != nil {
		os.Remove(destPath)
		return "", NewTailwindInstallError("install", "failed to close binary", err)
	}

	if i.Verbose {
		tailwindInstallLog.Printf("Installed Tailwind %s to %s", version, destPath)
	}

	return destPath, nil
}

// Install downloads and caches the Tailwind binary for the current platform.
// Returns the path to the installed binary.
func (i *TailwindInstaller) Install() (string, error) {
	version, err := i.ResolveTailwindVersion()
	if err != nil {
		return "", err
	}

	cached, err := i.IsTailwindCached(version)
	if err != nil {
		return "", err
	}
	if cached {
		binaryPath, err := i.GetCachedTailwindBinaryPath(version)
		if err != nil {
			return "", err
		}
		if i.Verbose {
			tailwindInstallLog.Printf("Using cached Tailwind %s", version)
		}
		return binaryPath, nil
	}

	platformAsset, err := GetTailwindPlatformAssetName()
	if err != nil {
		return "", err
	}

	expectedChecksum, err := i.fetchTailwindChecksum(version, platformAsset)
	if err != nil {
		return "", err
	}

	downloadPath, err := i.downloadTailwindAsset(version, platformAsset)
	if err != nil {
		return "", err
	}
	defer os.Remove(downloadPath)

	if err := verifyTailwindChecksum(downloadPath, expectedChecksum, i.Verbose); err != nil {
		return "", err
	}

	installedPath, err := i.installTailwindBinary(downloadPath, version)
	if err != nil {
		return "", err
	}

	return installedPath, nil
}

func tailwindConfigBoolean(value interface{}, fallback bool) bool {
	if value == nil {
		return fallback
	}
	switch v := value.(type) {
	case bool:
		return v
	case string:
		v = strings.TrimSpace(strings.ToLower(v))
		if v == "" {
			return fallback
		}
		if v == "true" || v == "1" || v == "yes" || v == "on" {
			return true
		}
		if v == "false" || v == "0" || v == "no" || v == "off" {
			return false
		}
	}
	return fallback
}

func tailwindConfigString(value interface{}) string {
	if value == nil {
		return ""
	}
	if v, ok := value.(string); ok {
		return strings.TrimSpace(v)
	}
	return ""
}

func tailwindConfigStringSlice(value interface{}) []string {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case []string:
		return v
	case []interface{}:
		out := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok {
				trimmed := strings.TrimSpace(s)
				if trimmed != "" {
					out = append(out, trimmed)
				}
			}
		}
		return out
	case string:
		trimmed := strings.TrimSpace(v)
		if trimmed != "" {
			return []string{trimmed}
		}
	}
	return nil
}

func tailwindNormalizeVersion(version string) (string, error) {
	trimmed := strings.TrimSpace(version)
	if trimmed == "" || trimmed == Latest {
		return Latest, nil
	}
	if strings.HasPrefix(trimmed, "v") {
		return trimmed, nil
	}
	if strings.Contains(trimmed, ".") {
		return "v" + trimmed, nil
	}
	return "", NewTailwindInstallError("version_check", "invalid version format", errors.New("version must be semver"))
}
