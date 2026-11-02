---
title: "Embeds"
description: "Embed internal posts and external URLs as rich preview cards"
date: 2024-01-24
published: true
tags:
  - documentation
  - plugins
  - content
---


The embeds plugin lets you embed rich preview cards for both internal posts and external URLs in your markdown content.

## Quick Start

### Internal Post Embeds

Embed another post from your site using the `![[slug]]` syntax:

```markdown
Check out this related article:

![[getting-started]]
```

This creates a preview card showing the embedded post's title, description, and date.

### External URL Embeds

Embed external URLs using `![embed](url)`, `[!embed](url)`, or Obsidian-style `![[https://url]]`:

```markdown
Here's a great resource:

![embed](https://example.com/article)

![embed](https://example.com/article|card)
![embed](https://example.com/article|no_title center)

[!embed](https://example.com/article)
[!embed](https://example.com/article|hover)

![[https://example.com/article]]
![[https://example.com/article|Custom Title]]
![[https://example.com/article|Custom Title|center full_width]]
```

This resolves metadata via oEmbed (when available), then falls back to Open Graph, and displays a card with the title, description, and image.

If the URL is a Hacker News discussion link, markata-go resolves the outbound article URL first so the card can use the real article metadata while still showing a "Found on HN" link to the original discussion.

Supported options (space-separated):

- `no_title`, `no_description`, `no_meta`
- `image_only`, `performance`, `hover`, `card`, `rich`
- `center`, `full_width`, `link`

## Internal Embeds

### Basic Syntax

```markdown
![[slug]]
```

The slug is the URL path of the post you want to embed. For example, if a post lives at `/guides/configuration/`, use:

```markdown
![[guides/configuration]]
```

### Custom Title

Override the display title with a pipe:

```markdown
![[guides/configuration|Configuration Guide]]
```

### What Gets Displayed

Internal embed cards show:
- **Title** - From the post or your custom title
- **Description** - First 200 characters
- **Date** - If the post has a date

### Not Found Behavior

If the embedded post doesn't exist, markata-go leaves the original syntax and adds an HTML comment:

```html
<!-- embed not found: nonexistent-post -->
![[nonexistent-post]]
```

This helps you spot broken embeds without breaking your build.

## External Embeds

### Syntax

```markdown
![embed](https://example.com/article)
![embed](https://example.com/article|card)
![embed](https://example.com/article|no_description)
[!embed](https://example.com/article|hover)
![[https://example.com/article]]
![[https://example.com/article|Custom Title]]
![[https://example.com/article|Custom Title|center full_width]]
```

> **Note:** The alt text must be exactly `embed`. Regular images are not affected. The Obsidian-style form is only recognized for full URLs (`http`/`https`).

### Metadata Resolution

External embeds fetch and display (in order based on strategy):
- **oEmbed** - Title, provider name, thumbnail image (if available)
- **Open Graph** - `og:title`, `og:description`, `og:image`, `og:site_name`

### Embed Modes

Embed modes control how external content is rendered:

- `card` (default): image + title + description + domain
- `rich`: full oEmbed HTML (iframes, scripts). YouTube defaults to Lite YouTube embeds.
- `performance`: image only, no text
- `hover`: image preview, loads embed HTML on hover
- `image_only`: image-only rendering (useful for photo providers)

### GitHub Gist Rendering

Gist embeds render the first file as a code block styled by your configured
code theme. The filename appears above the code block with a link to the
original gist.

### Caching

Metadata is cached for 7 days by default (configurable) to avoid repeated HTTP requests. Cache files are stored in `.cache/embeds/` by default.

### Fallback Behavior

If metadata can't be fetched:
- Title shows "External Link" (configurable)
- No image is displayed
- Domain name is shown

## Configuration

Add to your `markata-go.toml`:

```toml
[embeds]
enabled = true

# CSS classes for styling
internal_card_class = "embed-card"
external_card_class = "embed-card embed-card-external"

# External fetch settings
fetch_external = true        # Set false to skip HTTP requests
oembed_enabled = true        # Enable oEmbed resolution
oembed_auto_discover = false # Use HTML auto-discovery when no provider matches
default_mode = "card"        # Default embed mode
oembed_providers_url = "https://oembed.com/providers.json" # Providers list URL (empty disables)
resolution_strategy = "oembed_first"  # oembed_first | og_first | oembed_only
cache_dir = ".cache/embeds"  # Where to store cached metadata
cache_ttl = 604800           # Cache TTL in seconds (default 7 days)
timeout = 10                 # HTTP timeout in seconds

# Display settings
fallback_title = "External Link"  # Title when metadata is unavailable
show_image = true                 # Show images

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

### Disabling External Fetching

For faster builds or offline development:

```toml
[embeds]
fetch_external = false
```

External embeds will use the fallback title and show no image.

### Disabling oEmbed

```toml
[embeds]
oembed_enabled = false
resolution_strategy = "og_first"
```

This skips oEmbed providers and falls back to Open Graph metadata only.

### Enabling oEmbed Auto-Discovery

```toml
[embeds]
oembed_auto_discover = true
```

When enabled, markata-go looks for oEmbed link tags in HTML pages and uses the discovered endpoint if no provider matches.

### Provider Mode Overrides

You can set a default mode for all embeds and override it per provider:

```toml
[embeds]
default_mode = "card"

[embeds.providers.youtube]
mode = "rich"

[embeds.providers.giphy]
mode = "image_only"
```

### YouTube oEmbed (Lite YouTube)

YouTube oEmbeds in rich mode render as Lite YouTube embeds by default for faster loads.
Lite YouTube assets are loaded only when needed and follow the vendor assets mode
(`assets.mode`) for CDN vs self-hosted paths.

## Styling

Embed cards use these CSS classes:

| Class | Description |
|-------|-------------|
| `.embed-card` | Base container |
| `.embed-card-external` | Added for external embeds |
| `.embed-card-link` | The clickable anchor |
| `.embed-card-image` | Image wrapper (external only) |
| `.embed-card-rich` | Rich HTML container |
| `.embed-card-hover` | Hover-mode container |
| `.embed-card-hover-embed` | Hover-mode embed HTML container |
| `.embed-card-content` | Text content area |
| `.embed-card-title` | Title |
| `.embed-card-description` | Description |
| `.embed-card-meta` | Date or domain |

### Custom Styling

Override the default styles in your custom CSS:

```css
/* Make internal embeds more prominent */
.embed-card:not(.embed-card-external) {
  border-left-width: 6px;
  border-left-color: var(--color-accent);
}

/* Larger images for external embeds */
.embed-card-image {
  width: 250px;
}
```

## Code Blocks

Embed syntax inside code blocks is preserved:

````markdown
Normal embed:
![[my-post]]

Code example (not processed):
```
![[my-post]]
```
````

## Examples

### Blog Post with Related Reading

```markdown
---
title: "Advanced Configuration"
---

# Advanced Configuration

This guide builds on the basics. If you're new, start here:

![[getting-started]]

## Deep Dive

...content...

## Further Reading

Check out these external resources:

![embed](https://gohugo.io/getting-started/)

![embed](https://www.markdownguide.org/)
```

### Documentation with Cross-References

```markdown
# API Reference

## Authentication

For setup instructions, see:

![[guides/authentication|Authentication Guide]]

## Rate Limiting

Related:
![[api/errors]]
![[guides/best-practices]]
```

## Comparison with Wikilinks

| Feature | Embeds (`![[]]`) | Wikilinks (`[[]]`) |
|---------|------------------|-------------------|
| Output | Preview card | Inline link |
| Shows description | Yes | No |
| Shows date | Yes | No |
| Best for | Related content, references | In-text mentions |

Use embeds when you want readers to see a preview. Use wikilinks for inline mentions.

## Troubleshooting

### Embed Not Rendering

1. Check the slug matches the target post's URL path
2. Ensure the target post is published
3. Look for `<!-- embed not found -->` comments in the HTML

### External Embed Shows Fallback

1. Check your internet connection
2. Verify the URL is accessible
3. Some sites block robots - the fallback is expected
4. Increase timeout if the site is slow

### Images Not Showing

1. Ensure `show_image = true` in config
2. The page may not have an `og:image` tag
3. Image URL may be invalid or blocked

### Build Is Slow

External fetching adds latency. Options:

1. Set `fetch_external = false` during development
2. Reduce `timeout` value
3. Let the cache warm up (first build is slowest)
