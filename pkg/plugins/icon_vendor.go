// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/WaylonWalker/markata/pkg/assets"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/runtimeenv"
)

const (
	iconVendorPluginName    = "icon_vendor"
	defaultIconVendorCache  = ".markata/cache/icon-packs"
	defaultIconVendorTarget = "static/.icons"
	iconVendorSourceNPM     = "npm"
	iconVendorSourceURL     = "url"
)

var (
	iconVendorNameRegex    = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)
	iconVendorPackageRegex = regexp.MustCompile(`^(?:@[A-Za-z0-9._-]+/)?[A-Za-z0-9][A-Za-z0-9._-]*$`)
	iconVendorVersionRegex = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9.+_-]*$`)
)

type iconVendorConfig struct {
	enabled  bool
	cacheDir string
	target   string
	packs    []iconVendorPack
}

type iconVendorPack struct {
	name        string
	version     string
	source      string
	packageName string
	url         string
	archivePath string
	iconsPath   string
	licensePath string
}

// IconVendorPlugin downloads pinned icon-pack archives and materializes their
// SVGs under a local icon root before IconsPlugin indexes them.
type IconVendorPlugin struct {
	config iconVendorConfig
}

// NewIconVendorPlugin creates a disabled-by-default icon vendor plugin.
func NewIconVendorPlugin() *IconVendorPlugin {
	return &IconVendorPlugin{config: defaultIconVendorConfig()}
}

// Name returns the unique plugin name.
func (p *IconVendorPlugin) Name() string { return iconVendorPluginName }

// Priority downloads icon packs early enough for IconsPlugin.Configure to see
// the resulting static/.icons tree in the same build.
func (p *IconVendorPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageConfigure {
		return lifecycle.PriorityEarly
	}
	return lifecycle.PriorityDefault
}

// Configure parses [icons.vendor], restores/downloads requested archives, and
// materializes validated SVGs plus license text under the configured target.
func (p *IconVendorPlugin) Configure(m *lifecycle.Manager) error {
	p.config = defaultIconVendorConfig()
	if m == nil || m.Config() == nil || m.Config().Extra == nil {
		return nil
	}
	iconsRaw, ok := m.Config().Extra[iconPluginName].(map[string]interface{})
	if !ok {
		return nil
	}
	vendorRaw, ok := iconsRaw["vendor"].(map[string]interface{})
	if !ok {
		return nil
	}
	if err := p.config.apply(vendorRaw); err != nil {
		return err
	}
	if !p.config.enabled || len(p.config.packs) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(p.config.packs))
	for i := range p.config.packs {
		pack := &p.config.packs[i]
		if _, exists := seen[pack.name]; exists {
			return fmt.Errorf("icon vendor: duplicate pack name %q", pack.name)
		}
		seen[pack.name] = struct{}{}
		if err := pack.validate(); err != nil {
			return err
		}
		if err := p.vendorPack(context.Background(), *pack); err != nil {
			return err
		}
	}
	return nil
}

func defaultIconVendorConfig() iconVendorConfig {
	return iconVendorConfig{
		enabled:  false,
		cacheDir: defaultIconVendorCache,
		target:   defaultIconVendorTarget,
	}
}

func (c *iconVendorConfig) apply(raw map[string]interface{}) error {
	if enabled, ok := raw["enabled"].(bool); ok {
		c.enabled = enabled
	}
	if cacheDir, ok := raw["cache_dir"].(string); ok && strings.TrimSpace(cacheDir) != "" {
		c.cacheDir = strings.TrimSpace(cacheDir)
	}
	if target, ok := raw["target"].(string); ok && strings.TrimSpace(target) != "" {
		c.target = strings.TrimSpace(target)
	}
	packMaps := iconVendorPackMaps(raw["packs"])
	c.packs = make([]iconVendorPack, 0, len(packMaps))
	for _, packMap := range packMaps {
		pack := iconVendorPack{
			name:        iconVendorString(packMap, "name"),
			version:     iconVendorString(packMap, "version"),
			source:      strings.ToLower(iconVendorString(packMap, "source")),
			packageName: iconVendorString(packMap, "package"),
			url:         iconVendorString(packMap, "url"),
			archivePath: iconVendorString(packMap, "archive_path"),
			iconsPath:   iconVendorString(packMap, "icons_path"),
			licensePath: iconVendorString(packMap, "license_path"),
		}
		if pack.source == "" {
			pack.source = iconVendorSourceNPM
		}
		if pack.archivePath == "" {
			pack.archivePath = "package"
		}
		if pack.iconsPath == "" {
			pack.iconsPath = "icons"
		}
		if pack.licensePath == "" {
			pack.licensePath = "LICENSE"
		}
		c.packs = append(c.packs, pack)
	}
	return nil
}

func iconVendorPackMaps(raw interface{}) []map[string]interface{} {
	switch values := raw.(type) {
	case []map[string]interface{}:
		return values
	case []interface{}:
		out := make([]map[string]interface{}, 0, len(values))
		for _, value := range values {
			if pack, ok := value.(map[string]interface{}); ok {
				out = append(out, pack)
			}
		}
		return out
	default:
		return nil
	}
}

func iconVendorString(values map[string]interface{}, key string) string {
	value, ok := values[key].(string)
	if !ok {
		return ""
	}
	return strings.TrimSpace(value)
}

func (p iconVendorPack) validate() error {
	if !iconVendorNameRegex.MatchString(p.name) {
		return fmt.Errorf("icon vendor: invalid pack name %q", p.name)
	}
	if !iconVendorVersionRegex.MatchString(p.version) {
		return fmt.Errorf("icon vendor %s: version must be pinned and path-safe", p.name)
	}
	if p.source != iconVendorSourceNPM && p.source != iconVendorSourceURL {
		return fmt.Errorf("icon vendor %s: unsupported source %q (want npm or url)", p.name, p.source)
	}
	if p.source == iconVendorSourceNPM && !iconVendorPackageRegex.MatchString(p.packageName) {
		return fmt.Errorf("icon vendor %s: invalid npm package %q", p.name, p.packageName)
	}
	if p.source == iconVendorSourceURL && !strings.HasPrefix(p.url, "https://") && !strings.HasPrefix(p.url, "http://") {
		return fmt.Errorf("icon vendor %s: url source requires an http(s) url", p.name)
	}
	for field, value := range map[string]string{
		"archive_path": p.archivePath,
		"icons_path":   p.iconsPath,
		"license_path": p.licensePath,
	} {
		if !safeIconVendorRelativePath(value) {
			return fmt.Errorf("icon vendor %s: %s must be a safe relative path", p.name, field)
		}
	}
	return nil
}

func safeIconVendorRelativePath(value string) bool {
	if value == "" || filepath.IsAbs(value) {
		return false
	}
	clean := filepath.Clean(filepath.FromSlash(value))
	return clean != "." && clean != ".." && !strings.HasPrefix(clean, ".."+string(filepath.Separator))
}

func (p iconVendorPack) archiveURL() string {
	if p.source == iconVendorSourceURL {
		return p.url
	}
	archiveName := p.packageName
	if slash := strings.LastIndexByte(archiveName, '/'); slash >= 0 {
		archiveName = archiveName[slash+1:]
	}
	return "https://registry.npmjs.org/" + p.packageName + "/-/" + archiveName + "-" + p.version + ".tgz"
}

func (p iconVendorPack) fingerprint() string {
	return strings.Join([]string{
		p.name,
		p.version,
		p.source,
		p.packageName,
		p.archiveURL(),
		p.archivePath,
		p.iconsPath,
		p.licensePath,
	}, "\n") + "\n"
}

func (p *IconVendorPlugin) vendorPack(ctx context.Context, pack iconVendorPack) error {
	targetDir := filepath.Join(p.config.target, pack.name)
	markerPath := filepath.Join(p.config.cacheDir, ".published", pack.name+".txt")
	fingerprint := pack.fingerprint()

	if iconVendorPublished(markerPath, targetDir, fingerprint) {
		return nil
	}
	if runtimeenv.OfflineEnabled() && iconVendorHasSVGs(targetDir) {
		log.Printf("[icon_vendor] offline: using existing vendored pack %s", pack.name)
		return nil
	}

	asset := assets.Asset{
		Name:        "icon-pack-" + pack.name,
		URL:         pack.archiveURL(),
		LocalPath:   filepath.Join("archives", pack.name, pack.version),
		Version:     pack.version,
		Type:        "archive",
		ExtractPath: filepath.ToSlash(pack.archivePath),
	}
	downloader := assets.NewDownloader(p.config.cacheDir, false)
	result, err := downloader.Download(ctx, asset)
	if err != nil {
		return fmt.Errorf("icon vendor %s: %w", pack.name, err)
	}
	if result.Cached {
		log.Printf("[icon_vendor] using cached %s@%s", pack.name, pack.version)
	} else {
		log.Printf("[icon_vendor] downloaded %s@%s", pack.name, pack.version)
	}

	cachedRoot := downloader.GetCachedPath(asset)
	if cachedRoot == "" {
		return fmt.Errorf("icon vendor %s: archive cache path is unavailable", pack.name)
	}
	iconsRoot := filepath.Join(cachedRoot, filepath.FromSlash(pack.iconsPath))
	licenseFile := filepath.Join(cachedRoot, filepath.FromSlash(pack.licensePath))
	if err := p.materializePack(pack, iconsRoot, licenseFile, targetDir); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(markerPath), 0o755); err != nil {
		return fmt.Errorf("icon vendor %s: create publish marker dir: %w", pack.name, err)
	}
	if err := os.WriteFile(markerPath, []byte(fingerprint), 0o644); err != nil { //nolint:gosec // local cache metadata
		return fmt.Errorf("icon vendor %s: write publish marker: %w", pack.name, err)
	}
	return nil
}

func iconVendorPublished(markerPath, targetDir, fingerprint string) bool {
	if !iconVendorHasSVGs(targetDir) {
		return false
	}
	data, err := os.ReadFile(markerPath)
	return err == nil && string(data) == fingerprint
}

func iconVendorHasSVGs(root string) bool {
	found := false
	err := filepath.WalkDir(root, func(_ string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !entry.IsDir() && strings.EqualFold(filepath.Ext(entry.Name()), ".svg") {
			found = true
			return fs.SkipAll
		}
		return nil
	})
	return err == nil && found
}

func (p *IconVendorPlugin) materializePack(pack iconVendorPack, iconsRoot, licenseFile, targetDir string) error {
	info, err := os.Stat(iconsRoot)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("icon vendor %s: icons path %s not found in archive", pack.name, pack.iconsPath)
	}
	license, err := os.ReadFile(licenseFile)
	if err != nil {
		return fmt.Errorf("icon vendor %s: license %s not found in archive: %w", pack.name, pack.licensePath, err)
	}

	tmpDir := targetDir + ".tmp"
	if err := os.RemoveAll(tmpDir); err != nil {
		return fmt.Errorf("icon vendor %s: reset temp dir: %w", pack.name, err)
	}
	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		return fmt.Errorf("icon vendor %s: create temp dir: %w", pack.name, err)
	}

	count, err := copyVendoredSVGs(pack, iconsRoot, tmpDir)
	if err != nil {
		cleanupIconVendorTemp(tmpDir)
		return err
	}
	if count == 0 {
		cleanupIconVendorTemp(tmpDir)
		return fmt.Errorf("icon vendor %s: no SVG files found under %s", pack.name, pack.iconsPath)
	}

	if err := os.RemoveAll(targetDir); err != nil {
		cleanupIconVendorTemp(tmpDir)
		return fmt.Errorf("icon vendor %s: replace target: %w", pack.name, err)
	}
	if err := os.Rename(tmpDir, targetDir); err != nil {
		cleanupIconVendorTemp(tmpDir)
		return fmt.Errorf("icon vendor %s: publish target: %w", pack.name, err)
	}

	if err := p.writeVendorLicense(pack, license); err != nil {
		return err
	}
	log.Printf("[icon_vendor] vendored %d SVGs for %s into %s", count, pack.name, targetDir)
	return nil
}

func copyVendoredSVGs(pack iconVendorPack, iconsRoot, tmpDir string) (int, error) {
	count := 0
	err := filepath.WalkDir(iconsRoot, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return nil
		}
		if entry.IsDir() || !strings.EqualFold(filepath.Ext(path), ".svg") {
			return nil
		}
		rel, err := filepath.Rel(iconsRoot, path)
		if err != nil {
			return err
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if _, safe := prepareIconSVG(string(raw)); !safe {
			return fmt.Errorf("icon vendor %s: refusing unsafe SVG %s", pack.name, filepath.ToSlash(rel))
		}
		dest := filepath.Join(tmpDir, rel)
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(dest, raw, 0o644); err != nil { //nolint:gosec // vendored static assets are intentionally readable
			return err
		}
		count++
		return nil
	})
	return count, err
}

func cleanupIconVendorTemp(path string) {
	if err := os.RemoveAll(path); err != nil {
		log.Printf("[icon_vendor] cleanup temp %s: %v", path, err)
	}
}

func (p *IconVendorPlugin) writeVendorLicense(pack iconVendorPack, license []byte) error {
	licenseDir := filepath.Join(p.config.target, "licenses")
	if err := os.MkdirAll(licenseDir, 0o755); err != nil {
		return fmt.Errorf("icon vendor %s: create license dir: %w", pack.name, err)
	}
	licenseTarget := filepath.Join(licenseDir, pack.name+".txt")
	if err := os.WriteFile(licenseTarget, license, 0o644); err != nil { //nolint:gosec // vendored license text is intentionally readable
		return fmt.Errorf("icon vendor %s: write license: %w", pack.name, err)
	}
	return nil
}

var (
	_ lifecycle.Plugin          = (*IconVendorPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*IconVendorPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*IconVendorPlugin)(nil)
)
