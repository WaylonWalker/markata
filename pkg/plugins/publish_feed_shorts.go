package plugins

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/WaylonWalker/markata-go/pkg/lifecycle"
	"github.com/WaylonWalker/markata-go/pkg/models"
	"github.com/WaylonWalker/markata-go/pkg/templates"
)

const shortsChunkSize = 128

const (
	shortsPlaceholderWidth = 72
	shortsPreviewWidth     = 240
	shortsMainImageWidth   = 1280
)

// validateShortsRoutes makes an opt-in peer route a safe, unique output path.
// It is called before any feed publisher workers start writing.
func validateShortsRoutes(feeds []models.FeedConfig, posts []*models.Post) error {
	owners := make(map[string]string)
	for i := range feeds {
		feed := &feeds[i]
		if feed.Formats.HTML {
			route := "/" + strings.Trim(feed.Slug, "/") + "/"
			if feed.Slug == "" {
				route = "/"
			}
			owners[route] = "feed " + feed.Slug
		}
	}
	for _, post := range posts {
		if post == nil || post.Skip || post.Draft || !post.Published {
			continue
		}
		owners["/"+strings.Trim(post.Slug, "/")+"/"] = "post " + post.Slug
	}
	for i := range feeds {
		feed := &feeds[i]
		if !feed.HasView(models.FeedViewShorts) {
			continue
		}
		route := feed.ShortsURL()
		if err := validateShortsPath(route); err != nil {
			return fmt.Errorf("shorts for feed %q: %w", feed.Slug, err)
		}
		if owner, taken := owners[route]; taken {
			return fmt.Errorf("shorts path %q for feed %q conflicts with %s", route, feed.Slug, owner)
		}
		for used, owner := range owners {
			if strings.HasPrefix(used, route+"data/") || strings.HasPrefix(route, used+"data/") {
				return fmt.Errorf("shorts path %q for feed %q overlaps %s (%s)", route, feed.Slug, used, owner)
			}
		}
		owners[route] = "shorts view of " + feed.Slug
	}
	return nil
}

func validateShortsPath(route string) error {
	if len(route) < 3 || route[0] != '/' || !strings.HasSuffix(route, "/") || strings.Contains(route, "//") {
		return fmt.Errorf("invalid shorts_path %q: expected an absolute path ending in /", route)
	}
	for _, part := range strings.Split(strings.Trim(route, "/"), "/") {
		if part == "." || part == ".." || part == "" {
			return fmt.Errorf("invalid shorts_path %q", route)
		}
		for _, ch := range part {
			if !(ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' ||
				ch >= '0' && ch <= '9' || ch == '-' || ch == '_') {
				return fmt.Errorf("invalid shorts_path %q: only letters, digits, _ and - are supported", route)
			}
		}
	}
	return nil
}

// cleanupDisabledShorts only deletes output previously identified as a Shorts
// artifact, never an arbitrary directory that happens to share a path.
func cleanupDisabledShorts(outputDir, route, feedSlug string) error {
	if err := validateShortsPath(route); err != nil {
		return nil // A legacy disabled setting must not delete arbitrary files.
	}
	dir := filepath.Join(outputDir, strings.Trim(route, "/"))
	marker := filepath.Join(dir, "data", "index.json")
	raw, err := os.ReadFile(marker)
	if err != nil {
		return nil
	}
	var identity struct {
		FeedSlug string `json:"feed_slug"`
	}
	if json.Unmarshal(raw, &identity) != nil || identity.FeedSlug != feedSlug {
		return nil // Another feed owns this Shorts route.
	}
	if err := os.Remove(filepath.Join(dir, "index.html")); err != nil && !os.IsNotExist(err) {
		return err
	}
	return os.RemoveAll(filepath.Join(dir, "data"))
}

func extraMediaString(extra map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		if v, ok := extra[key].(string); ok && strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

// safeShortsMediaURL rejects dangerous URI schemes when media was entered in frontmatter.
func safeShortsMediaURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	if parsed.Scheme != "" && parsed.Scheme != "https" && parsed.Scheme != "http" {
		return ""
	}
	if parsed.Scheme == "" && strings.HasPrefix(raw, "//") {
		return ""
	}
	return raw
}

// shortsSizedImageURL applies Dropper sizing only to absolute trusted media
// URLs. Relative site assets and untrusted external images keep their original
// URL instead of receiving query parameters that the host may ignore.
func shortsCanResizeImageURL(raw string) bool {
	parsed, err := url.Parse(raw)
	return err == nil && parsed.Host != "" && (parsed.Scheme == "http" || parsed.Scheme == "https") && templates.IsTrustedMediaURL(raw)
}

func shortsSizedImageURL(raw string, width int) string {
	if !shortsCanResizeImageURL(raw) {
		return raw
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	query := parsed.Query()
	// Dropper accepts both short and long aliases. Remove stale aliases before
	// writing w so an earlier width value cannot conflict with the new size.
	query.Del("width")
	query.Del("height")
	parsed.RawQuery = query.Encode()
	return templates.WithSize(parsed.String(), width, 0)
}

func shortsPostItem(post *models.Post) map[string]interface{} {
	extra := post.Extra
	if extra == nil {
		extra = map[string]interface{}{}
	}
	image := safeShortsMediaURL(extraMediaString(extra, embedOptionImage, "cover", "cover_image", "og_image"))
	video := safeShortsMediaURL(extraMediaString(extra, templateTypeVideo))
	media := image
	if video != "" {
		media = video
	}
	kind := embedOptionImage
	if templates.IsVideoURL(media) {
		kind = templateTypeVideo
	}
	poster := ""
	if kind == templateTypeVideo {
		poster = safeShortsMediaURL(templates.PosterURLFromMap(extra, media))
		if poster == "" && image != "" && !templates.IsVideoURL(image) {
			poster = image
		}
	}
	thumb := ""
	placeholder := ""
	src := media
	if kind == embedOptionImage && shortsCanResizeImageURL(media) {
		thumb = shortsSizedImageURL(media, shortsPreviewWidth)
		placeholder = shortsSizedImageURL(media, shortsPlaceholderWidth)
		src = shortsSizedImageURL(media, shortsMainImageWidth)
	}
	if kind == templateTypeVideo {
		src = media
		if poster != "" {
			thumb = shortsSizedImageURL(poster, shortsPreviewWidth)
			placeholder = shortsSizedImageURL(poster, shortsPlaceholderWidth)
		}
	}
	posterSized := ""
	if poster != "" {
		posterSized = shortsSizedImageURL(poster, 720)
	}
	href := post.Href
	if href == "" {
		href = "/" + strings.Trim(post.Slug, "/") + "/"
	}
	description := ""
	if post.Description != nil {
		description = *post.Description
	}
	title := post.PlainTitle()
	if title == "" {
		title = post.Slug
	}
	return map[string]interface{}{
		"id":          post.Slug,
		"href":        href,
		"title":       title,
		"description": description,
		"alt":         extraMediaString(extra, "image_alt", "alt", "caption"),
		"kind":        kind,
		"src":         src,
		"thumb":       thumb,
		"placeholder": placeholder,
		"poster":      posterSized,
		"mime":        templates.VideoMIMEType(media),
	}
}

// shortsPostCount matches the intentionally public-only Shorts manifest projection.
func shortsPostCount(posts []*models.Post) int {
	count := 0
	for _, post := range posts {
		if post != nil && !post.Private && post.Published && !post.Draft && !post.Skip {
			count++
		}
	}
	return count
}

func (p *PublishFeedsPlugin) publishShortsPages(feed *models.FeedConfig, cfg *lifecycle.Config, outputDir string) error {
	route := feed.ShortsURL()
	if err := validateShortsPath(route); err != nil {
		return err
	}
	dir := filepath.Join(outputDir, strings.Trim(route, "/"))
	dataDir := filepath.Join(dir, "data")
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return err
	}
	if err := p.writeShortsManifest(feed, dataDir); err != nil {
		return err
	}
	return p.writeShortsHTML(feed, cfg, dir)
}

func (p *PublishFeedsPlugin) writeShortsManifest(feed *models.FeedConfig, dataDir string) error {
	entries := make([]map[string]interface{}, 0, len(feed.Posts))
	ids := make([]string, 0, len(feed.Posts))
	for _, post := range feed.Posts {
		// Explicitly never export private post metadata or decrypted fields.
		if post == nil || post.Private || !post.Published || post.Skip || post.Draft {
			continue
		}
		entries = append(entries, shortsPostItem(post))
		ids = append(ids, post.Slug)
	}

	for start := 0; start < len(entries); start += shortsChunkSize {
		end := start + shortsChunkSize
		if end > len(entries) {
			end = len(entries)
		}
		chunk, err := json.Marshal(entries[start:end])
		if err != nil {
			return err
		}
		file := filepath.Join(dataDir, fmt.Sprintf("%04d.json", start/shortsChunkSize))
		if err := p.safeWriteFile(file, append(chunk, '\n')); err != nil {
			return err
		}
	}
	// Remove obsolete chunks after a collection shrinks.
	names, err := os.ReadDir(dataDir)
	if err != nil {
		return err
	}
	chunkCount := (len(entries) + shortsChunkSize - 1) / shortsChunkSize
	for _, entry := range names {
		if entry.IsDir() || entry.Name() == "index.json" || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		var number int
		if _, err := fmt.Sscanf(entry.Name(), "%04d.json", &number); err == nil && number >= chunkCount {
			if err := os.Remove(filepath.Join(dataDir, entry.Name())); err != nil {
				return err
			}
		}
	}

	index, err := json.Marshal(map[string]interface{}{
		"version":    1,
		"feed_slug":  feed.Slug,
		"total":      len(entries),
		"chunk_size": shortsChunkSize,
		"ids":        ids,
	})
	if err != nil {
		return err
	}
	return p.safeWriteFile(filepath.Join(dataDir, "index.json"), append(index, '\n'))
}

func (p *PublishFeedsPlugin) writeShortsHTML(feed *models.FeedConfig, cfg *lifecycle.Config, dir string) error {
	route := feed.ShortsURL()
	modelsConfig := ToModelsConfig(cfg)
	templateDir := PluginNameTemplates
	if v, ok := cfg.Extra["templates_dir"].(string); ok && v != "" {
		templateDir = v
	}
	themeName := ThemeDefault
	if theme, ok := cfg.Extra["theme"].(models.ThemeConfig); ok && theme.Name != "" {
		themeName = theme.Name
	} else if theme, ok := cfg.Extra["theme"].(map[string]interface{}); ok {
		if name, ok := theme["name"].(string); ok && name != "" {
			themeName = name
		}
	} else if name, ok := cfg.Extra["theme"].(string); ok && name != "" {
		themeName = name
	}
	engine, err := p.getOrCreateEngine(templateDir, themeName)
	if err != nil {
		return err
	}
	if !engine.TemplateExists("feed-shorts.html") {
		return fmt.Errorf("missing feed-shorts.html template for Shorts view")
	}
	page := &models.FeedPage{Number: 1, TotalPages: 1}
	// The manifest owns the complete collection. Avoid eagerly materializing
	// thousands of post maps while rendering this lightweight HTML shell.
	shellFeed := *feed
	shellFeed.Posts = nil
	shellFeed.Pages = nil
	ctx := templates.NewFeedContext(&shellFeed, page, modelsConfig)
	ctx.Set("shorts_manifest_url", route+"data/index.json")
	ctx.Set("feed_robots", feed.Robots)
	rendered, err := engine.Render("feed-shorts.html", ctx)
	if err != nil {
		return err
	}
	return p.safeWriteFile(filepath.Join(dir, "index.html"), []byte(rendered))
}
