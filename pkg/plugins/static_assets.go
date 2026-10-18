package plugins

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/templates"
	"github.com/WaylonWalker/markata/pkg/themes"
)

// StaticAssetsPlugin copies static assets from themes and project directories to output.
// It handles:
// 1. Theme static files (themes/[theme]/static/*)
// 2. Project static files (static/*)
// Project files take precedence over theme files (local override).
type StaticAssetsPlugin struct{}

// NewStaticAssetsPlugin creates a new StaticAssetsPlugin.
func NewStaticAssetsPlugin() *StaticAssetsPlugin {
	return &StaticAssetsPlugin{}
}

// Name returns the unique name of the plugin.
func (p *StaticAssetsPlugin) Name() string {
	return "static_assets"
}

// Configure computes content hashes for JS/CSS assets before templates are rendered.
// This enables cache busting via the theme_asset_hashed filter.
func (p *StaticAssetsPlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()

	// Get theme name
	themeName := ThemeDefault
	if extra := config.Extra; extra != nil {
		if theme, ok := extra["theme"].(map[string]interface{}); ok {
			if name, ok := theme["name"].(string); ok && name != "" {
				themeName = name
			}
		}
		if name, ok := extra["theme"].(string); ok && name != "" {
			themeName = name
		}
	}

	// Compute hashes for assets in priority order (last wins)
	assetHashes := make(map[string]string)

	// 1. Hash embedded assets (base layer)
	if themeName == ThemeDefault {
		if err := p.hashEmbeddedAssets(assetHashes); err != nil {
			return fmt.Errorf("hashing embedded assets: %w", err)
		}
	}

	// 2. Hash filesystem theme assets (overrides embedded)
	themeStaticDir := p.findThemeStaticDir(themeName)
	if themeStaticDir != "" {
		if err := p.hashDirectoryAssets(themeStaticDir, "", assetHashes); err != nil {
			return fmt.Errorf("hashing theme assets: %w", err)
		}
	}

	// 3. Hash project assets (highest priority, overrides theme)
	projectStaticDir := p.projectStaticDir(config)
	if _, err := os.Stat(projectStaticDir); err == nil {
		if err := p.hashDirectoryAssets(projectStaticDir, "", assetHashes); err != nil {
			return fmt.Errorf("hashing project assets: %w", err)
		}
	}

	// Store hashes in Manager for Write stage use
	for path, hash := range assetHashes {
		m.SetAssetHash(path, hash)
	}

	// Set hashes in templates package for theme_asset_hashed filter
	templates.SetAssetHashes(assetHashes)

	return nil
}

// Validate fingerprints the final asset registry after all Configure hooks.
// Generated styles such as palette-core.css register after static files.
func (p *StaticAssetsPlugin) Validate(m *lifecycle.Manager) error {
	if cache := GetBuildCache(m); cache != nil {
		if cache.SetAssetsHash(buildcache.HashAssetMap(m.AssetHashes())) {
			log.Printf("[static_assets] JS/CSS assets changed, full rebuild required")
		}
	}
	return nil
}

// Write copies static assets to the output directory.
// Files are copied in layers with increasing priority:
// 1. Embedded theme static files (lowest priority, base layer)
// 2. Filesystem theme static files (can override embedded)
// 3. Project static files (highest priority, can override all)
// Hashed copies are created in Cleanup stage after all transformations.
func (p *StaticAssetsPlugin) Write(m *lifecycle.Manager) error {
	config := m.Config()
	outputDir := config.OutputDir
	outputFS, err := openOutputRoot(outputDir)
	if err != nil {
		return fmt.Errorf("open static asset output directory: %w", err)
	}
	defer outputFS.Close()

	// Get theme name from config (default to "default")
	themeName := ThemeDefault
	if extra := config.Extra; extra != nil {
		if theme, ok := extra["theme"].(map[string]interface{}); ok {
			if name, ok := theme["name"].(string); ok && name != "" {
				themeName = name
			}
		}
		// Also check for simple theme string
		if name, ok := extra["theme"].(string); ok && name != "" {
			themeName = name
		}
	}

	// Layer 1: Copy embedded static files for default theme (base layer)
	// This ensures all default assets are present even if filesystem theme is incomplete
	if themeName == ThemeDefault {
		if err := p.copyEmbeddedStaticWithRoot(outputDir, outputFS); err != nil {
			return fmt.Errorf("copying embedded static files: %w", err)
		}
	}

	// Layer 2: Copy filesystem theme static files (overrides embedded)
	themeStaticDir := p.findThemeStaticDir(themeName)
	if themeStaticDir != "" {
		if err := p.copyDirWithRoot(themeStaticDir, outputDir, outputFS); err != nil {
			return fmt.Errorf("copying theme static files: %w", err)
		}
	}

	// Layer 3: Copy project static files (highest priority, overrides theme files)
	projectStaticDir := p.projectStaticDir(config)
	if _, err := os.Stat(projectStaticDir); err == nil {
		if err := p.copyDirWithRoot(projectStaticDir, outputDir, outputFS); err != nil {
			return fmt.Errorf("copying project static files: %w", err)
		}
	}

	return nil
}

func (p *StaticAssetsPlugin) projectStaticDir(config *lifecycle.Config) string {
	if config != nil && config.Extra != nil {
		if configured, ok := config.Extra["assets_dir"].(string); ok && configured != "" {
			if filepath.IsAbs(configured) || config.ContentDir == "" {
				return filepath.Clean(configured)
			}
			return filepath.Clean(filepath.Join(config.ContentDir, configured))
		}
	}
	if config != nil && config.ContentDir != "" {
		return filepath.Join(config.ContentDir, StaticDir)
	}
	return StaticDir
}

// Cleanup creates hashed copies of JS/CSS files after all Write plugins have run.
// This ensures the hashes match the final transformed content (after palette_css, minifiers, etc.)
func (p *StaticAssetsPlugin) Cleanup(m *lifecycle.Manager) error {
	config := m.Config()
	outputDir := config.OutputDir
	outputFS, err := openOutputRoot(outputDir)
	if err != nil {
		return fmt.Errorf("open static asset output directory: %w", err)
	}
	defer outputFS.Close()

	if err := p.createHashedCopiesWithRoot(m, outputDir, outputFS); err != nil {
		return fmt.Errorf("creating hashed asset copies: %w", err)
	}
	// A cache-busted feed stylesheet must actually exist with matching bytes.
	// Otherwise a successful build can publish HTML whose CSS URL returns 404
	// (or whose old hash points to stale content on the CDN).
	if err := verifyVersionedFeedCSS(outputFS, m.GetAssetHash("css/feeds.css")); err != nil {
		return fmt.Errorf("feed CSS asset integrity: %w", err)
	}

	return nil
}

// verifyVersionedFeedCSS checks the link target emitted by theme_asset_hashed.
// It is intentionally limited to the feed stylesheet; other generated assets
// may have separate, plugin-specific hashing and transformation rules.
func verifyVersionedFeedCSS(root *os.Root, registeredHash string) error {
	const sourcePath = "css/feeds.css"
	original, err := root.ReadFile(sourcePath)
	if os.IsNotExist(err) && registeredHash == "" {
		// Custom themes may not ship the default feed stylesheet.
		return nil
	}
	if err != nil {
		return fmt.Errorf("reading %s: %w", sourcePath, err)
	}
	if registeredHash == "" {
		return fmt.Errorf("%s exists but has no registered content hash", sourcePath)
	}
	actualHash := fmt.Sprintf("%x", sha256.Sum256(original))[:8]
	if registeredHash != actualHash {
		return fmt.Errorf("%s registered hash %s does not match published bytes %s", sourcePath, registeredHash, actualHash)
	}
	hashedPath := "css/feeds." + registeredHash + ".css"
	published, err := root.ReadFile(hashedPath)
	if err != nil {
		return fmt.Errorf("missing versioned stylesheet %s: %w", hashedPath, err)
	}
	if !bytes.Equal(original, published) {
		return fmt.Errorf("versioned stylesheet %s differs from %s", hashedPath, sourcePath)
	}
	if strings.Contains(string(original), ".feed--editorial") &&
		!strings.Contains(string(published), ".editorial-shot-grid") {
		return fmt.Errorf("versioned stylesheet %s lacks the editorial shot grid", hashedPath)
	}
	return nil
}

// findThemeStaticDir searches for theme static directory in various locations.
func (p *StaticAssetsPlugin) findThemeStaticDir(themeName string) string {
	// 1. Check current working directory
	cwdPath := filepath.Join("themes", themeName, "static")
	if _, err := os.Stat(cwdPath); err == nil {
		return cwdPath
	}

	// 2. Check relative to executable
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)

		// themes next to executable
		exeThemePath := filepath.Join(exeDir, "themes", themeName, "static")
		if _, err := os.Stat(exeThemePath); err == nil {
			return exeThemePath
		}

		// Check parent/share/markata-go/themes (standard install location)
		parentDir := filepath.Dir(exeDir)
		sharePath := filepath.Join(parentDir, "share", "markata-go", "themes", themeName, "static")
		if _, err := os.Stat(sharePath); err == nil {
			return sharePath
		}
	}

	return ""
}

func (p *StaticAssetsPlugin) copyEmbeddedStaticWithRoot(outputDir string, outputFS *os.Root) error {
	staticFS := themes.DefaultStatic()
	if staticFS == nil {
		return nil // No embedded static files
	}

	return fs.WalkDir(staticFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip root directory
		if path == "." {
			return nil
		}

		dstPath := filepath.Join(outputDir, path)
		relative, err := outputRelativePath(outputDir, dstPath)
		if err != nil {
			return err
		}

		if d.IsDir() {
			return outputFS.MkdirAll(relative, 0o755)
		}

		// Read embedded file
		content, err := fs.ReadFile(staticFS, path)
		if err != nil {
			return fmt.Errorf("reading embedded file %s: %w", path, err)
		}

		// Ensure parent directory exists
		if directory := filepath.Dir(relative); directory != "." {
			if err := outputFS.MkdirAll(directory, 0o755); err != nil {
				return fmt.Errorf("creating parent directory: %w", err)
			}
		}

		// Write file
		if err := outputFS.WriteFile(relative, content, 0o644); err != nil {
			return fmt.Errorf("writing file %s: %w", dstPath, err)
		}

		return nil
	})
}

// copyDir recursively copies a directory to the destination.
func (p *StaticAssetsPlugin) copyDir(src, dst string) error {
	outputFS, err := openOutputRoot(dst)
	if err != nil {
		return err
	}
	defer outputFS.Close()
	return p.copyDirWithRoot(src, dst, outputFS)
}

func (p *StaticAssetsPlugin) copyDirWithRoot(src, dst string, outputFS *os.Root) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			log.Printf("[static_assets] Skipping symbolic link %s", path)
			return nil
		}

		// Calculate relative path from source
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return fmt.Errorf("calculating relative path: %w", err)
		}

		// The source root maps to the already-open output root.
		if relPath == "." {
			return nil
		}

		// Calculate destination path
		dstPath := filepath.Join(dst, relPath)
		if err := validateOutputPath(dst, dstPath); err != nil {
			return err
		}
		if info.IsDir() {
			// Create directory
			destination, err := outputRelativePath(dst, dstPath)
			if err != nil {
				return err
			}
			return outputFS.MkdirAll(destination, 0o755)
		}

		// Copy file
		return p.copyFileWithRoot(path, "", dst, dstPath, outputFS)
	})
}

func (p *StaticAssetsPlugin) copyFileWithRoot(src, sourceRelative, outputRoot, dst string, outputFS *os.Root) error {
	relative, err := outputRelativePath(outputRoot, dst)
	if err != nil {
		return err
	}
	directory := filepath.Dir(relative)
	if directory != "." {
		if err := outputFS.MkdirAll(directory, 0o755); err != nil {
			return fmt.Errorf("creating parent directory: %w", err)
		}
	}

	// Open source file. Hashed copies read their source through the same
	// descriptor-backed root as the destination so a path replacement cannot
	// redirect the read between validation and copying.
	var srcFile *os.File
	if sourceRelative != "" {
		srcFile, err = outputFS.Open(sourceRelative)
	} else {
		srcFile, err = os.Open(src)
	}
	if err != nil {
		return fmt.Errorf("opening source file: %w", err)
	}
	defer srcFile.Close()

	// Create destination file below the descriptor-backed output root.
	dstFile, err := outputFS.OpenFile(relative, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("creating destination file: %w", err)
	}
	defer dstFile.Close()

	// Copy content
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("copying file content: %w", err)
	}

	// Preserve file permissions
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return nil // Non-critical, continue without preserving permissions
	}
	return dstFile.Chmod(srcInfo.Mode())
}

// shouldHashAsset returns true if the file should be content-hashed for cache busting.
// Only JS and CSS files are hashed to avoid breaking image references.
func (p *StaticAssetsPlugin) shouldHashAsset(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".js" || ext == ".css" || ext == ".mjs"
}

// hashEmbeddedAssets computes SHA-256 hashes for embedded JS/CSS assets.
func (p *StaticAssetsPlugin) hashEmbeddedAssets(hashes map[string]string) error {
	staticFS := themes.DefaultStatic()
	if staticFS == nil {
		return nil
	}

	return fs.WalkDir(staticFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || path == "." {
			return err
		}

		if !p.shouldHashAsset(path) {
			return nil
		}

		// Read file content
		content, err := fs.ReadFile(staticFS, path)
		if err != nil {
			return fmt.Errorf("reading embedded asset %s: %w", path, err)
		}

		// Compute hash (first 8 chars of SHA-256)
		hash := fmt.Sprintf("%x", sha256.Sum256(content))[:8]
		hashes[path] = hash

		return nil
	})
}

// hashDirectoryAssets computes SHA-256 hashes for JS/CSS assets in a directory.
// prefix is the relative path prefix to prepend to file paths in the hash map.
func (p *StaticAssetsPlugin) hashDirectoryAssets(dir, prefix string, hashes map[string]string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}

		if !p.shouldHashAsset(path) {
			return nil
		}

		// Calculate relative path from dir
		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return fmt.Errorf("calculating relative path: %w", err)
		}

		// Add prefix if provided
		if prefix != "" {
			relPath = filepath.Join(prefix, relPath)
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading asset %s: %w", path, err)
		}

		// Compute hash (first 8 chars of SHA-256)
		hash := fmt.Sprintf("%x", sha256.Sum256(content))[:8]
		hashes[relPath] = hash

		return nil
	})
}

// createHashedCopies creates hashed copies of JS/CSS files in the output directory.
// For each file with a hash, creates a copy like main.js -> main.abc12345.js
func (p *StaticAssetsPlugin) createHashedCopies(m *lifecycle.Manager, outputDir string) error {
	outputFS, err := openOutputRoot(outputDir)
	if err != nil {
		return err
	}
	defer outputFS.Close()
	return p.createHashedCopiesWithRoot(m, outputDir, outputFS)
}

func (p *StaticAssetsPlugin) createHashedCopiesWithRoot(m *lifecycle.Manager, outputDir string, outputFS *os.Root) error {
	assetHashes := m.AssetHashes()

	for path, hash := range assetHashes {
		// Original file location in output
		origPath := filepath.Join(outputDir, path)
		origRelative, err := outputRelativePath(outputDir, origPath)
		if err != nil {
			return err
		}

		// Check if file exists (might not if overridden)
		info, err := outputFS.Lstat(origRelative)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("hashed asset source is a symlink: %s", origPath)
		}

		// Compute hashed filename: main.js -> main.abc12345.js
		ext := filepath.Ext(path)
		base := strings.TrimSuffix(path, ext)
		hashedPath := base + "." + hash + ext
		hashedFullPath := filepath.Join(outputDir, hashedPath)
		if err := validateOutputPath(outputDir, hashedFullPath); err != nil {
			return err
		}

		// Create hashed copy
		if err := p.copyFileWithRoot(origPath, origRelative, outputDir, hashedFullPath, outputFS); err != nil {
			return fmt.Errorf("creating hashed copy %s: %w", hashedPath, err)
		}
	}

	return nil
}

// Priority runs static assets before plugins that consume copied media.
func (p *StaticAssetsPlugin) Priority(stage lifecycle.Stage) int {
	// Configure registers asset hashes; Write copies source images before
	// image optimization reads them.
	if stage == lifecycle.StageConfigure || stage == lifecycle.StageWrite {
		return lifecycle.PriorityEarly
	}
	return lifecycle.PriorityDefault
}

// Ensure StaticAssetsPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*StaticAssetsPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*StaticAssetsPlugin)(nil)
	_ lifecycle.ValidatePlugin  = (*StaticAssetsPlugin)(nil)
	_ lifecycle.WritePlugin     = (*StaticAssetsPlugin)(nil)
	_ lifecycle.CleanupPlugin   = (*StaticAssetsPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*StaticAssetsPlugin)(nil)
)
