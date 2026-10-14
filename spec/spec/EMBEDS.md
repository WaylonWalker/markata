# Embeds Plugin Specification

The embeds plugin enables rich embedding of both internal posts and external URLs within markdown content.

## Overview

The embeds plugin processes two types of embed syntax:

1. **Internal embeds** (`![[slug]]`) - Embed another post from the same site as a preview card
2. **External embeds** (`![embed](url)`, `[!embed](url)`, or `![[https://url]]`) - Embed external URLs with rich metadata (oEmbed + Open Graph)

## Automatic link-post destinations

Posts with `template: link`, `bookmark`, `stars`, or `til` and a valid
HTTP(S) `link` frontmatter field (or `url` alias) automatically render an
external preview before the authored body. This applies even when the body
is empty. When the body already contains an embed for that destination,
or a standalone URL handled by the render-stage URL plugins, no additional
preview is generated. The destination is included in the transform-cache
key so editing frontmatter invalidates stale previews.

The post page and feed link card must retain the site's own permalink,
and expose a separate, visible, accessible link to the original destination
with `u-bookmark-of` attribution. Do not fetch private destinations.
A missing/invalid URL or unavailable metadata must not break the build.

## Syntax

### Internal Embeds

```markdown
![[slug]]           # Basic internal embed
![[slug|Title]]     # Custom display title
```

### External Embeds

```markdown
![embed](https://example.com/article)         # External embed with OG metadata
![embed](https://example.com/article|card)    # External embed with mode/class options
[!embed](https://example.com/article)         # Bracket-style external embed
[!embed](https://example.com/article|hover)   # Bracket-style with options
![[https://example.com/article]]              # Obsidian-style external embed
![[https://example.com/article|Title]]        # Obsidian-style with custom title
![[https://example.com/article|Title|center]] # Obsidian-style with classes
```

Note: The alt text must be exactly `embed` to trigger external embedding. Regular images with other alt text are not affected.

Supported options (space-separated):
- `no_title`, `no_description`, `no_meta`
- `image_only`, `performance`, `hover`, `card`, `rich`
- `center`, `full_width`, `link`

## Configuration

```toml
[embeds]
enabled = true                                    # Enable/disable the plugin
internal_card_class = "embed-card"               # CSS class for internal cards
external_card_class = "embed-card embed-card-external"  # CSS class for external cards
fetch_external = true                            # Fetch OG metadata from external URLs
oembed_enabled = true                            # Enable oEmbed resolution
oembed_auto_discover = false                     # Discover oEmbed endpoints from HTML link tags
default_mode = "card"                            # Default embed mode
oembed_providers_url = "https://oembed.com/providers.json" # Providers list URL (empty disables)
resolution_strategy = "oembed_first"             # oembed_first | og_first | oembed_only
cache_dir = ".cache/embeds"                      # Cache directory for external metadata
cache_ttl = 604800                               # Cache TTL in seconds
timeout = 10                                     # HTTP timeout in seconds
fallback_title = "External Link"                 # Title when OG title is unavailable
show_image = true                                # Show OG images in external embeds

[embeds.providers]
giphy = { enabled = true, mode = "image_only" }
youtube = { enabled = true, mode = "rich" } # default
vimeo = { enabled = true }
tiktok = { enabled = true }
flickr = { enabled = true }
spotify = { enabled = true }
soundcloud = { enabled = true }
codepen = { enabled = true }
codesandbox = { enabled = true }
jsfiddle = { enabled = true }
observable = { enabled = true }
github = { enabled = true }
slideshare = { enabled = true }
prezi = { enabled = true }
speakerdeck = { enabled = true }
issuu = { enabled = true }
datawrapper = { enabled = true }
flourish = { enabled = true }
infogram = { enabled = true }
reddit = { enabled = true }
dailymotion = { enabled = true }
wistia = { enabled = true }
giphy = { enabled = true }
```

## Internal Embeds

### Behavior

1. The plugin looks up the target post by slug (case-insensitive)
2. If found, it generates an embed card with:
   - Linked title (from display text or post title)
   - Description (truncated to 200 characters)
   - Date (if available)
3. If not found, it adds a warning comment and preserves the original syntax

### Self-Reference Protection

A post cannot embed itself. Attempting to do so adds a warning comment:

```html
<!-- cannot embed self -->
![[self-slug]]
```

### Generated HTML

```html
<div class="embed-card">
  <a href="/target-post/" class="embed-card-link">
    <div class="embed-card-content">
      <div class="embed-card-title">Target Post Title</div>
      <div class="embed-card-description">Post description...</div>
      <div class="embed-card-meta">Jan 15, 2024</div>
    </div>
  </a>
</div>
```

## External Embeds

### Behavior

1. Validates the URL (must be http or https)
2. If the URL is a Hacker News discussion link, resolves the outbound article URL through the HN item API and uses that for metadata lookup and the main card target when available.
3. Resolves metadata using the configured strategy:
    - **oembed_first** (default): try oEmbed providers, fall back to OG
    - **og_first**: try OG, fall back to oEmbed providers
    - **oembed_only**: only use oEmbed (no OG fallback)
4. Caches metadata (configurable TTL)
5. Selects a mode (explicit option, provider override, oEmbed type default, or global default)
6. Generates an embed card with:
    - OG image (if available and enabled)
    - OG title (or fallback)
    - OG description (truncated)
    - Site name and domain

When HN normalization succeeds, the card should still show a secondary "Found on HN" link to the original discussion URL instead of hiding the source.

### Metadata Resolution

The plugin extracts:
- **oEmbed**: title, provider name, thumbnail URL (if available)
- **Open Graph**: `og:title`, `og:description`, `og:image`, `og:site_name`

For GitHub Gist, the oEmbed resolver fetches the gist JSON and renders the
first file as a code block using the site's code theme. The HTML includes
the gist filename and links back to the gist. Code content is escaped and
rendered as a plain code block (not a script embed).

### Modes

Embed modes control how external content is rendered:

- `card` (default): image + title + description + domain
- `rich`: full oEmbed HTML (iframes, scripts). For YouTube providers, rich uses Lite YouTube embeds by default.
- `performance`: image only, no text
- `hover`: image preview, loads embed HTML on hover
- `image_only`: image-only rendering (useful for photo providers)

### oEmbed Providers (Phase 2)

Supported providers include Phase 1 plus additional Phase 2 providers:

- YouTube (`https://www.youtube.com/oembed`)
- Vimeo (`https://vimeo.com/api/oembed.json`)
- TikTok (`https://www.tiktok.com/oembed`)
- Flickr (`https://www.flickr.com/services/oembed/`)
- Spotify (`https://open.spotify.com/oembed`)
- SoundCloud (`https://soundcloud.com/oembed`)
- CodePen (`https://codepen.io/api/oembed`)
- CodeSandbox (`https://codesandbox.io/oembed`)
- JSFiddle (`https://jsfiddle.net/services/oembed/`)
- Observable (`https://api.observablehq.com/oembed`)
- GitHub Gist (`https://github.com/services/oembed`)
- SlideShare (`https://www.slideshare.net/api/oembed/2`)
- Prezi (`https://prezi.com/services/oembed/`)
- Speaker Deck (`https://speakerdeck.com/oembed.json`)
- Issuu (`https://issuu.com/oembed`)
- Datawrapper (`https://api.datawrapper.de/v3/oembed/`)
- Flourish (`https://app.flourish.studio/api/v1/oembed`)
- Infogram (`https://infogram.com/oembed`)
- Reddit (`https://www.reddit.com/oembed`)
- Dailymotion (`https://www.dailymotion.com/services/oembed`)
- Wistia (`https://fast.wistia.com/oembed.json`)
- GIPHY (`https://giphy.com/services/oembed`)

Falls back to `<title>` and `<meta name="description">` if OG tags are missing.

### Caching

External metadata is cached as JSON files:
- Location: `.cache/embeds/` (configurable)
- File name: SHA-256 hash of URL (first 8 bytes) + source suffix (`-oembed` or `-og`)
- Expiration: configurable via `cache_ttl` (default 7 days)

Example cache file:
```json
{
  "title": "Article Title",
  "description": "Article description",
  "image": "https://example.com/og-image.jpg",
  "site_name": "Example Site",
  "type": "article",
  "fetched_at": 1705350000,
  "source": "oembed"
}
```

### Generated HTML

```html
<div class="embed-card embed-card-external">
  <a href="https://example.com/article" class="embed-card-link" target="_blank" rel="noopener noreferrer">
    <div class="embed-card-image">
      <img src="https://example.com/og-image.jpg" alt="" loading="lazy">
    </div>
    <div class="embed-card-content">
      <div class="embed-card-title">Article Title</div>
      <div class="embed-card-description">Article description...</div>
      <div class="embed-card-meta">Example Site &middot; example.com</div>
    </div>
  </a>
</div>
```

### YouTube oEmbed (Lite YouTube)

When a YouTube oEmbed resolves in rich mode, markata-go renders a Lite YouTube embed by default:

```html
<lite-youtube videoid="VIDEO_ID" title="Video Title" playlabel="Play: Video Title"></lite-youtube>
```

The Lite YouTube CSS/JS assets are loaded only when needed and follow the vendor assets mode
(`assets.mode`) for CDN vs self-hosted paths.

## Code Block Protection

Embed syntax inside fenced code blocks is preserved and not processed:

```markdown
Normal embed: ![[my-post]]

` ` `
Code example: ![[my-post]]
` ` `
```

## Lifecycle Stage

The embeds plugin runs in the **Transform** stage with `PriorityEarly` (-100), ensuring it processes content before:
- Wikilinks plugin
- Table of Contents extraction
- Jinja-MD processing

## Error Handling

| Scenario | Behavior |
|----------|----------|
| Internal embed not found | Warning comment + original syntax preserved |
| Self-reference | Warning comment + original syntax preserved |
| External URL invalid | Original syntax preserved |
| Obsidian external URL invalid | Original syntax preserved |
| External fetch fails | Uses fallback title, no image |
| oEmbed provider disabled | Treat as matched, fall back if allowed |
| oEmbed discovery disabled | Treat as not matched |
| External timeout | Uses fallback title, no image |

## CSS Classes

| Class | Purpose |
|-------|---------|
| `.embed-card` | Base container |
| `.embed-card-external` | External embed modifier |
| `.embed-card-link` | Clickable anchor |
| `.embed-card-image` | Image container (external only) |
| `.embed-card-content` | Text content wrapper |
| `.embed-card-title` | Title element |
| `.embed-card-description` | Description element |
| `.embed-card-meta` | Metadata (date, domain) |

## Performance Considerations

1. **Caching** - External metadata is cached to avoid repeated HTTP requests
2. **Timeout** - Configurable HTTP timeout (default 10s) prevents slow builds
3. **Concurrent Processing** - Posts are processed concurrently
4. **Body Limit** - External pages are limited to 1MB to prevent memory issues
5. **Disable Fetching** - Set `fetch_external = false` to skip OG HTTP requests entirely
6. **Disable oEmbed** - Set `oembed_enabled = false` or use `oembed_only` to avoid OG fallback
7. **Auto-Discovery** - Set `oembed_auto_discover = true` to discover endpoints via HTML link tags

## Related Features

- **Wikilinks** (`[[slug]]`) - Simple internal links without preview cards
- **One-line Links** - Automatic link card for standalone URLs
- **Link Collector** - Tracks outlinks for backlink generation
