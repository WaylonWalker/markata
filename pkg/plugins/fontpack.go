package plugins

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/fontpacks"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/templates"
)

const (
	fontpackCacheFile        = ".markata-fontpack-cache"
	fontpackPreloadCacheFile = ".markata-fontpack-preloads.json"
	fontpackCacheVersion     = "12"
)

const fontpackRoleHeading = "heading"

// FontpackPlugin installs one site-wide typography stylesheet. It never calls
// a subsetter: bundled tiers are immutable catalog artifacts.
type FontpackPlugin struct {
	name     string
	source   *fontpacks.CatalogSource
	prepared *fontpackBuild
}

type fontpackBuild struct {
	resolved   *fontpacks.Resolved
	cacheKey   string
	cacheReady bool
	preloads   fontpackPreloadCache
}

type fontpackPreloadCache struct {
	Hash          string              `json:"hash"`
	URLs          map[string][]string `json:"urls"`
	EditorialURLs map[string][]string `json:"editorial_urls"`
}

func NewFontpackPlugin() *FontpackPlugin { return &FontpackPlugin{} }
func (p *FontpackPlugin) Name() string   { return "fontpack" }
func (p *FontpackPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageRender {
		return lifecycle.PriorityLate - 1 // After Markdown, before page templates.
	}
	if stage == lifecycle.StageWrite {
		return lifecycle.PriorityFirst
	}
	return lifecycle.PriorityDefault
}

func (p *FontpackPlugin) Configure(m *lifecycle.Manager) error {
	if m.Config().Extra == nil {
		m.Config().Extra = make(map[string]any)
	}
	name := configuredFontpackName(m.Config().Extra)
	p.name = name
	path := ""
	if v, ok := m.Config().Extra["fontpacks_file"].(string); ok && v != "" {
		path = v
	}
	var err error
	if path != "" {
		p.source, err = fontpacks.LoadSource(path)
	} else {
		p.source, err = fontpacks.BuiltinSource()
	}
	if err != nil {
		return fmt.Errorf("load font catalog for pack %q: %w", name, err)
	}
	m.Config().Extra["fontpack_css"] = true
	p.prepared = nil
	return nil
}

func configuredFontpackName(extra map[string]any) string {
	const brushPosterFontpack = "brush-poster"
	canonicalize := func(name string) string {
		if name == brushPosterFontpack {
			return renderingFontpackBrush
		}
		return name
	}
	if configured, ok := extra["models_config"].(*models.Config); ok && configured.Theme.Fontpack != "" {
		return canonicalize(configured.Theme.Fontpack)
	}
	if value, ok := extra["fontpack"].(string); ok && value != "" {
		return canonicalize(value)
	}
	return "system"
}

// Render resolves the actual emitted tiers after Markdown and before the
// templates, so preload URLs and the stylesheet hash agree with Write.
func (p *FontpackPlugin) Render(m *lifecycle.Manager) error {
	_, err := p.prepare(m)
	return err
}

//nolint:gocyclo // Resolving page overrides, picker packs, cache state, and preloads is one coordinated preparation pass.
func (p *FontpackPlugin) prepare(m *lifecycle.Manager) (*fontpackBuild, error) {
	if p.prepared != nil {
		return p.prepared, nil
	}
	if p.source == nil || p.source.Catalog == nil {
		return nil, fmt.Errorf("fontpack catalog has not been configured")
	}
	defaultName, _, err := p.source.Catalog.ResolvePack(p.name)
	if err != nil {
		return nil, err
	}
	posts := m.Posts()
	readers := make([]io.Reader, 0, len(posts)*2)
	names := []string{defaultName}
	pageNames := make(map[string]string)
	for _, post := range posts {
		readers = append(readers, strings.NewReader(post.ArticleHTML), strings.NewReader("\n"))
		name := p.name
		if value, ok := post.Extra["fontpack"].(string); ok && value != "" {
			if value == "brush-poster" {
				value = renderingFontpackBrush
			}
			name = value
		}
		resolvedName, _, err := p.source.Catalog.ResolvePack(name)
		if err != nil {
			return nil, fmt.Errorf("post %q fontpack %q: %w", post.Path, name, err)
		}
		pageNames[post.Path] = resolvedName
		if post.Extra == nil {
			post.Extra = make(map[string]interface{})
		}
		post.Extra["_resolved_fontpack"] = resolvedName
		if !slices.Contains(names, resolvedName) {
			names = append(names, resolvedName)
		}
	}
	pickerEnabled := themeSwitcherEnabled(m.Config().Extra)
	if pickerEnabled {
		// The theme picker lets visitors choose any pack. @font-face files are
		// fetched only when a pack is actually used, so offering all of them
		// costs CSS bytes, not downloads.
		for _, name := range fontpacks.SortedKeys(p.source.Catalog.FontPacks) {
			if !slices.Contains(names, name) {
				names = append(names, name)
			}
		}
	}
	coverage, err := fontpacks.CollectCoverage(io.MultiReader(readers...))
	if err != nil {
		return nil, fmt.Errorf("collect fontpack visible coverage: %w", err)
	}
	cacheKey, err := fontpackCacheKey(coverage.Signature(), defaultName, names, pickerEnabled, p.source)
	if err != nil {
		return nil, err
	}
	build := &fontpackBuild{cacheKey: cacheKey}
	output := m.Config().OutputDir
	if p.source.Builtin && fontpackOutputCached(output, cacheKey) {
		data, readErr := os.ReadFile(filepath.Join(output, fontpackPreloadCacheFile))
		if readErr == nil && json.Unmarshal(data, &build.preloads) == nil && validFontpackPreloadCache(output, names, p.source.Catalog, build.preloads) {
			build.cacheReady = true
		}
	}
	if !build.cacheReady {
		resolved, err := p.source.Catalog.ResolveManyFSWithCoverage(names, p.source.FS, p.source.Root, coverage, fontpackResolveOptions(p.source))
		if err != nil {
			return nil, err
		}
		if pickerEnabled {
			resolved.CSS += fontpackManifestCSS(p.source.Catalog, defaultName, resolved.Packs)
		}
		build.resolved = resolved
		build.preloads.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(resolved.CSS)))[:8]
		build.preloads.URLs = make(map[string][]string, len(resolved.Packs))
		build.preloads.EditorialURLs = make(map[string][]string, len(resolved.Packs))
		for name, pack := range resolved.Packs {
			build.preloads.URLs[name] = fontpackPreloadURLs(pack, resolved.Assets)
			build.preloads.EditorialURLs[name] = fontpackEditorialPreloadURLs(pack, resolved.Assets)
		}
	}
	if raw, ok := m.Cache().Get("build_cache"); ok {
		if cache, ok := raw.(*buildcache.Cache); ok {
			// Only invalidate rendered HTML when its font CSS/preload URLs have
			// actually changed. Coverage can change without changing the emitted
			// tiers, and should not turn an ordinary content edit into a full page
			// rebuild.
			cache.SetFontpackHash(fontpackHTMLCacheHash(build.preloads))
		}
	}
	m.SetAssetHash("css/fonts.css", build.preloads.Hash)
	templates.SetAssetHashes(map[string]string{"css/fonts.css": build.preloads.Hash})
	m.Config().Extra["fontpack_preload_urls"] = build.preloads.URLs[defaultName]
	m.Config().Extra["fontpack_editorial_preload_urls"] = build.preloads.EditorialURLs[defaultName]
	// JSON is escaped by encoding/json (including HTML metacharacters), and
	// inserted as an object literal only when the picker is enabled.
	if pickerEnabled {
		data, err := json.Marshal(build.preloads.URLs)
		if err != nil {
			return nil, fmt.Errorf("encode font preloads: %w", err)
		}
		m.Config().Extra["fontpack_preloads_json"] = string(data)
		editorialData, err := json.Marshal(build.preloads.EditorialURLs)
		if err != nil {
			return nil, fmt.Errorf("encode editorial font preloads: %w", err)
		}
		m.Config().Extra["fontpack_editorial_preloads_json"] = string(editorialData)
	}
	for _, post := range m.Posts() {
		if name := pageNames[post.Path]; name != "" {
			post.Extra["_fontpack_preload_urls"] = build.preloads.URLs[name]
		}
	}
	p.prepared = build
	return p.prepared, nil
}

func fontpackHTMLCacheHash(preloads fontpackPreloadCache) string {
	urls, err := json.Marshal(preloads)
	if err != nil {
		return preloads.Hash
	}
	sum := sha256.Sum256(append([]byte(preloads.Hash+"\x00"), urls...))
	return hex.EncodeToString(sum[:])
}

func validFontpackPreloadCache(output string, names []string, catalog *fontpacks.Catalog, cached fontpackPreloadCache) bool {
	if len(cached.Hash) != 8 || len(cached.URLs) == 0 {
		return false
	}
	if _, err := hex.DecodeString(cached.Hash); err != nil {
		return false
	}
	for _, name := range names {
		resolved, _, err := catalog.ResolvePack(name)
		if err != nil {
			return false
		}
		for _, lists := range []map[string][]string{cached.URLs, cached.EditorialURLs} {
			urls, ok := lists[resolved]
			if !ok {
				return false
			}
			for _, url := range urls {
				file := strings.TrimPrefix(url, "/assets/fonts/")
				if file == url || file != filepath.Base(file) {
					return false
				}
				if info, err := os.Stat(filepath.Join(output, "assets", "fonts", file)); err != nil || !info.Mode().IsRegular() {
					return false
				}
			}
		}
	}
	return true
}

// Ordinary pages preload body and heading faces. Other roles load on demand.
func fontpackPreloadURLs(pack fontpacks.FontPack, assets []fontpacks.Asset) []string {
	return fontpackPreloadURLsForRoles(pack, assets, "body", "display", fontpackRoleHeading)
}

func fontpackEditorialPreloadURLs(pack fontpacks.FontPack, assets []fontpacks.Asset) []string {
	return fontpackPreloadURLsForRoles(pack, assets, "body", "display", fontpackRoleHeading, "code")
}

// Select the requested emitted tier, falling back to another core tier or full.
func fontpackPreloadURLsForRoles(pack fontpacks.FontPack, assets []fontpacks.Asset, roles ...string) []string {
	urls := []string{}
	seen := map[string]bool{}
	for _, role := range roles {
		source := pack.Roles[role].Source
		tier := pack.Roles[role].Tier
		if source == "" || seen[source] {
			continue
		}
		seen[source] = true
		var best *fontpacks.Asset
		var full *fontpacks.Asset
		for i := range assets {
			a := &assets[i]
			if a.Source != source {
				continue
			}
			if a.Tier == tier {
				best = a
				break
			}
			if a.Tier == "full" {
				full = a
				continue
			}
			if best == nil && (a.Tier == "prose-core" || a.Tier == "display-core" || a.Tier == "code-core") {
				best = a
			}
		}
		if best == nil {
			best = full
		}
		if best != nil {
			urls = append(urls, best.URL)
		}
	}
	return urls
}

func (p *FontpackPlugin) Write(m *lifecycle.Manager) error {
	build, err := p.prepare(m)
	if err != nil {
		return err
	}
	output := m.Config().OutputDir
	if !build.cacheReady {
		if err := build.resolved.CopyFS(p.source.FS, p.source.Root, output); err != nil {
			return err
		}
		if p.source.Builtin {
			data, err := json.Marshal(build.preloads)
			if err != nil {
				return fmt.Errorf("encode font preload cache: %w", err)
			}
			if err := os.WriteFile(filepath.Join(output, fontpackPreloadCacheFile), data, 0o600); err != nil {
				return err
			}
			if err := os.WriteFile(filepath.Join(output, fontpackCacheFile), []byte(build.cacheKey), 0o600); err != nil {
				return err
			}
		} else {
			for _, name := range []string{fontpackCacheFile, fontpackPreloadCacheFile} {
				if err := os.Remove(filepath.Join(output, name)); err != nil && !os.IsNotExist(err) {
					return err
				}
			}
		}
	}
	for _, post := range m.Posts() {
		if name, ok := post.Extra["_resolved_fontpack"].(string); ok && name != "" {
			post.HTML = markPostFontpack(post.HTML, name)
		}
	}
	return nil
}

// FontpackManifestEntry describes one pack offered by the theme picker.
type FontpackManifestEntry struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Heading     string `json:"heading"`
	Body        string `json:"body"`
	Code        string `json:"code,omitempty"`
	// CSS font-family stacks and heading weight used by picker previews.
	HeadingFont   string  `json:"headingFont,omitempty"`
	HeadingWeight float64 `json:"headingWeight,omitempty"`
	BodyFont      string  `json:"bodyFont,omitempty"`
	CodeFont      string  `json:"codeFont,omitempty"`
}

// fontpackManifestCSS exposes the offered packs to the theme picker through
// custom properties, mirroring the palette and aesthetic manifests.
func fontpackManifestCSS(catalog *fontpacks.Catalog, defaultName string, packs map[string]fontpacks.FontPack) string {
	entries := make([]FontpackManifestEntry, 0, len(packs))
	for _, name := range fontpacks.SortedKeys(packs) {
		pack := packs[name]
		display := pack.Name
		if display == "" {
			display = name
		}
		// Single-quoted families keep the JSON free of backslash escapes, which
		// would not survive the CSS custom property round trip.
		headingFont, headingWeight := catalog.RoleFontFamily(pack, "display", "heading")
		bodyFont, _ := catalog.RoleFontFamily(pack, "body")
		codeFont, _ := catalog.RoleFontFamily(pack, "code", "mono")
		headingFont = strings.ReplaceAll(headingFont, `"`, `'`)
		bodyFont = strings.ReplaceAll(bodyFont, `"`, `'`)
		codeFont = strings.ReplaceAll(codeFont, `"`, `'`)
		entries = append(entries, FontpackManifestEntry{
			Name:          name,
			DisplayName:   display,
			Heading:       fontpackRoleFamily(catalog, pack, "display", "heading"),
			Body:          fontpackRoleFamily(catalog, pack, "body"),
			Code:          fontpackRoleFamily(catalog, pack, "code", "mono"),
			HeadingFont:   headingFont,
			HeadingWeight: headingWeight,
			BodyFont:      bodyFont,
			CodeFont:      codeFont,
		})
	}
	data, err := json.Marshal(entries)
	if err != nil {
		return ""
	}
	escaped := strings.NewReplacer(`\`, `\\`, `'`, `\'`, "\n", "").Replace(string(data))
	return fmt.Sprintf(":root{--fontpack-default:%q;--fontpack-manifest:'%s'}\n", defaultName, escaped)
}

// fontpackRoleFamily returns a human label for the first matching role:
// the font family for bundled sources or a generic name for system stacks.
func fontpackRoleFamily(catalog *fontpacks.Catalog, pack fontpacks.FontPack, roles ...string) string {
	for _, role := range roles {
		r, ok := pack.Roles[role]
		if !ok {
			continue
		}
		if src, ok := catalog.FontSources[r.Source]; ok && src.Family != "" {
			return src.Family
		}
		if r.Stack != "" {
			return "System " + strings.ReplaceAll(r.Stack, "-", " ")
		}
	}
	return ""
}

func fontpackOutputCached(output, key string) bool {
	data, err := os.ReadFile(filepath.Join(output, fontpackCacheFile))
	if err != nil || strings.TrimSpace(string(data)) != key {
		return false
	}
	if info, err := os.Stat(filepath.Join(output, "css", "fonts.css")); err != nil || !info.Mode().IsRegular() {
		return false
	}
	files, err := fontpacks.ManagedFontFiles(output)
	if err != nil {
		return false
	}
	for _, file := range files {
		info, err := os.Stat(filepath.Join(output, "assets", "fonts", file))
		if err != nil || !info.Mode().IsRegular() {
			return false
		}
	}
	return true
}

func fontpackResolveOptions(source *fontpacks.CatalogSource) fontpacks.ResolveOptions {
	// Bundled assets are immutable and already validated when the release is
	// built. Re-hashing every WOFF2 file on every site build dominates warm
	// builds, while custom catalogs must retain runtime checksum validation.
	return fontpacks.ResolveOptions{ValidateChecksums: !source.Builtin}
}

var _ lifecycle.ConfigurePlugin = (*FontpackPlugin)(nil)
var _ lifecycle.RenderPlugin = (*FontpackPlugin)(nil)
var _ lifecycle.WritePlugin = (*FontpackPlugin)(nil)
var _ lifecycle.PriorityPlugin = (*FontpackPlugin)(nil)
