// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"fmt"
	"html"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

const iconPluginName = "icons"

var (
	iconShortcodeRegex = regexp.MustCompile(`:([A-Za-z0-9][A-Za-z0-9._/-]*):`)
	iconWidthRegex     = regexp.MustCompile(`(?i)\swidth\s*=\s*(?:"[^"]*"|'[^']*'|[^\s>]+)`)
	iconHeightRegex    = regexp.MustCompile(`(?i)\sheight\s*=\s*(?:"[^"]*"|'[^']*'|[^\s>]+)`)
	iconEventAttrRegex = regexp.MustCompile(`(?i)\son[a-z]+\s*=`)
)

type iconAsset struct {
	name string
	svg  string
}

// IconsPlugin expands Zensical-style icon shortcodes from local SVG packs.
//
// Icons are loaded from .icons and static/.icons by default. Given an asset at
// static/.icons/lucide/smile.svg, both :lucide-smile: (Zensical's canonical
// Markdown spelling) and :lucide/smile: resolve to the same inline SVG. When
// Lucide is referenced without a local asset, the pinned default pack is
// vendored on demand and then indexed through the same local path.
type IconsPlugin struct {
	enabled         bool
	autoVendor      bool
	paths           []string
	packs           map[string]struct{}
	icons           map[string]iconAsset
	iconMu          sync.Mutex
	vendorAttempted bool
}

// NewIconsPlugin creates an icon shortcode plugin with local asset defaults and
// zero-config Lucide vendoring enabled.
func NewIconsPlugin() *IconsPlugin {
	return &IconsPlugin{
		enabled:    true,
		autoVendor: true,
		paths:      []string{".icons", "static/.icons"},
		packs:      make(map[string]struct{}),
		icons:      make(map[string]iconAsset),
	}
}

// Name returns the unique plugin name.
func (p *IconsPlugin) Name() string { return iconPluginName }

// Priority runs icon expansion after other Markdown source transforms (notably
// Jinja), but before the Render stage converts Markdown to HTML.
func (p *IconsPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageTransform {
		return lifecycle.PriorityLast
	}
	return lifecycle.PriorityDefault
}

// Configure loads icon settings and indexes local SVG packs.
func (p *IconsPlugin) Configure(m *lifecycle.Manager) error {
	p.enabled = true
	p.autoVendor = true
	p.paths = []string{".icons", "static/.icons"}
	p.packs = make(map[string]struct{})
	p.icons = make(map[string]iconAsset)
	p.vendorAttempted = false

	if m != nil && m.Config() != nil && m.Config().Extra != nil {
		p.applyConfig(m.Config().Extra[iconPluginName])
	}
	if !p.enabled {
		return nil
	}
	return p.loadIcons()
}

func (p *IconsPlugin) applyConfig(raw interface{}) {
	cfg, ok := raw.(map[string]interface{})
	if !ok {
		return
	}
	if enabled, ok := cfg["enabled"].(bool); ok {
		p.enabled = enabled
	}
	if autoVendor, ok := cfg["auto_vendor"].(bool); ok {
		p.autoVendor = autoVendor
	}
	if path, ok := cfg["path"].(string); ok && strings.TrimSpace(path) != "" {
		p.paths = []string{strings.TrimSpace(path)}
	}
	if paths := iconStringSlice(cfg["paths"]); len(paths) > 0 {
		p.paths = paths
	}
	if packs := iconStringSlice(cfg["packs"]); len(packs) > 0 {
		for _, pack := range packs {
			pack = strings.ToLower(strings.TrimSpace(pack))
			if pack != "" {
				p.packs[pack] = struct{}{}
			}
		}
	}
}

func iconStringSlice(raw interface{}) []string {
	switch values := raw.(type) {
	case []string:
		out := make([]string, 0, len(values))
		for _, value := range values {
			if value = strings.TrimSpace(value); value != "" {
				out = append(out, value)
			}
		}
		return out
	case []interface{}:
		out := make([]string, 0, len(values))
		for _, value := range values {
			if text, ok := value.(string); ok {
				if text = strings.TrimSpace(text); text != "" {
					out = append(out, text)
				}
			}
		}
		return out
	default:
		return nil
	}
}

func (p *IconsPlugin) loadIcons() error {
	p.icons = make(map[string]iconAsset)
	blockedAliases := make(map[string]struct{})

	for _, root := range p.paths {
		info, err := os.Stat(root)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return fmt.Errorf("icons: stat %s: %w", root, err)
		}
		if !info.IsDir() {
			continue
		}

		err = filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if entry.IsDir() || !strings.EqualFold(filepath.Ext(path), ".svg") {
				return nil
			}

			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			name := strings.TrimSuffix(filepath.ToSlash(rel), filepath.Ext(rel))
			name = strings.Trim(name, "/")
			if name == "" || !p.packAllowed(name) {
				return nil
			}

			contents, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			svg, ok := prepareIconSVG(string(contents))
			if !ok {
				return nil
			}
			asset := iconAsset{name: name, svg: svg}

			// Earlier paths win, so a project-level .icons override can shadow
			// an identically named asset in static/.icons.
			if _, exists := p.icons[name]; !exists {
				p.icons[name] = asset
			}

			alias := strings.ReplaceAll(name, "/", "-")
			if alias == name {
				return nil
			}
			if _, blocked := blockedAliases[alias]; blocked {
				return nil
			}
			if existing, exists := p.icons[alias]; exists && existing.name != name {
				delete(p.icons, alias)
				blockedAliases[alias] = struct{}{}
				return nil
			}
			p.icons[alias] = asset
			return nil
		})
		if err != nil {
			return fmt.Errorf("icons: index %s: %w", root, err)
		}
	}
	return nil
}

func (p *IconsPlugin) packAllowed(name string) bool {
	if len(p.packs) == 0 {
		return true
	}
	pack := name
	if slash := strings.IndexByte(name, '/'); slash >= 0 {
		pack = name[:slash]
	}
	_, ok := p.packs[strings.ToLower(pack)]
	return ok
}

func prepareIconSVG(raw string) (string, bool) {
	raw = strings.TrimSpace(strings.TrimPrefix(raw, "\ufeff"))
	lower := strings.ToLower(raw)
	start := strings.Index(lower, "<svg")
	end := strings.LastIndex(lower, "</svg>")
	if start < 0 || end < start {
		return "", false
	}
	end += len("</svg>")
	svg := raw[start:end]
	// Shortcodes are expanded into Markdown source before Goldmark parses it.
	// Keep the SVG on one line so multiline vendored files remain inline HTML
	// instead of being split into paragraphs and escaped as text.
	svg = strings.NewReplacer("\r\n", " ", "\n", " ", "\r", " ").Replace(svg)
	lower = strings.ToLower(svg)

	// Local SVGs are trusted authoring assets, but reject the most dangerous
	// active-content primitives now so future vendored packs have a safe base.
	if strings.Contains(lower, "<script") ||
		strings.Contains(lower, "<foreignobject") ||
		strings.Contains(lower, "javascript:") ||
		iconEventAttrRegex.MatchString(svg) {
		return "", false
	}

	openEnd := strings.IndexByte(svg, '>')
	if openEnd < 0 {
		return "", false
	}
	opening := svg[:openEnd]
	opening = iconWidthRegex.ReplaceAllString(opening, "")
	opening = iconHeightRegex.ReplaceAllString(opening, "")
	opening += ` width="1em" height="1em" aria-hidden="true" focusable="false"`
	return opening + svg[openEnd:], true
}

// Transform expands icon shortcodes in Markdown source.
func (p *IconsPlugin) Transform(m *lifecycle.Manager) error {
	if !p.enabled {
		return nil
	}
	posts := m.FilterPosts(func(post *models.Post) bool {
		return !post.Skip && post.Content != ""
	})
	if lifecycle.IsServeIncremental(m) {
		if affected := lifecycle.GetServeAffectedPaths(m); len(affected) > 0 {
			filtered := posts[:0]
			for _, post := range posts {
				if affected[post.Path] {
					filtered = append(filtered, post)
				}
			}
			posts = filtered
		}
	}
	return m.ProcessPostsSliceConcurrently(posts, func(post *models.Post) error {
		post.Content = p.processContent(post.Content)
		return nil
	})
}

func (p *IconsPlugin) processContent(content string) string {
	if content == "" {
		return content
	}

	lines := strings.SplitAfter(content, "\n")
	var out strings.Builder
	out.Grow(len(content))
	var fence byte
	var fenceLen int

	for _, line := range lines {
		marker, markerLen := markdownFenceMarker(line)
		if fence != 0 {
			out.WriteString(line)
			if marker == fence && markerLen >= fenceLen {
				fence = 0
				fenceLen = 0
			}
			continue
		}
		if marker != 0 {
			fence = marker
			fenceLen = markerLen
			out.WriteString(line)
			continue
		}
		out.WriteString(p.processInlineText(line))
	}
	return out.String()
}

func markdownFenceMarker(line string) (marker byte, markerLen int) {
	line = strings.TrimSuffix(line, "\n")
	line = strings.TrimSuffix(line, "\r")
	indent := 0
	for indent < len(line) && indent < 4 && line[indent] == ' ' {
		indent++
	}
	if indent > 3 || indent >= len(line) {
		return 0, 0
	}
	marker = line[indent]
	if marker != '`' && marker != '~' {
		return 0, 0
	}
	end := indent
	for end < len(line) && line[end] == marker {
		end++
	}
	if end-indent < 3 {
		return 0, 0
	}
	return marker, end - indent
}

func (p *IconsPlugin) processInlineText(text string) string {
	var out strings.Builder
	out.Grow(len(text))

	for pos := 0; pos < len(text); {
		next := strings.IndexByte(text[pos:], '`')
		if next < 0 {
			out.WriteString(p.replaceShortcodes(text[pos:]))
			break
		}
		next += pos
		out.WriteString(p.replaceShortcodes(text[pos:next]))

		runEnd := next
		for runEnd < len(text) && text[runEnd] == '`' {
			runEnd++
		}
		delimiter := text[next:runEnd]
		closing := strings.Index(text[runEnd:], delimiter)
		if closing < 0 {
			// Unmatched backticks are literal Markdown, so continue shortcode
			// expansion rather than suppressing the rest of the line.
			out.WriteString(p.replaceShortcodes(text[next:]))
			break
		}
		closingEnd := runEnd + closing + len(delimiter)
		out.WriteString(text[next:closingEnd])
		pos = closingEnd
	}
	return out.String()
}

func (p *IconsPlugin) replaceShortcodes(text string) string {
	return iconShortcodeRegex.ReplaceAllStringFunc(text, func(shortcode string) string {
		match := iconShortcodeRegex.FindStringSubmatch(shortcode)
		if len(match) != 2 {
			return shortcode
		}
		asset, ok := p.lookupIcon(match[1])
		if !ok {
			return shortcode
		}
		name := html.EscapeString(asset.name)
		return `<span class="icon" data-icon="` + name + `" aria-hidden="true" style="display:inline-flex;line-height:1;vertical-align:-0.125em">` + asset.svg + `</span>`
	})
}
