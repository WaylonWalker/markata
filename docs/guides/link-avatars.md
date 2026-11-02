---
title: "Link Avatars"
description: "Add favicon icons next to external links for better visual identification"
date: 2026-02-09
published: true
slug: /docs/guides/link-avatars/
tags:
  - plugins
  - links
  - customization
---


The `link_avatars` plugin automatically adds small favicon icons next to external links in your content. This helps readers quickly identify link destinations before clicking.

## Quick Start

Enable link avatars with minimal configuration:

```toml
[markata-go.link_avatars]
enabled = true
```

This will add a 16x16 pixel favicon icon before each external link using DuckDuckGo's icon service.

## How It Works

The plugin supports three modes:

- **js** (default): client-side enhancement with runtime favicon fetch
- **local**: build-time caching, HTML points to site-relative icon URLs
- **hosted**: build-time caching, HTML points to a hosted base URL (CDN)

In all modes, the plugin applies the same ignore rules and uses the same CSS for sizing and positioning.

## Configuration

### Basic Options

```toml
[markata-go.link_avatars]
enabled = true

# CSS selector for links to enhance (default: "a[href^='http']")
selector = "a[href^='http']"

# Avatar size in pixels (default: 16)
size = 16

# Position: "before" or "after" the link text (default: "before")
position = "before"

# Mode: "js", "local", or "hosted" (default: "js")
mode = "js"

# Hosted base URL for mode = "hosted"
# hosted_base_url = "https://cdn.example.com/markata/link-avatars"
```

### Service Selection

Choose your favicon service provider:

```toml
[markata-go.link_avatars]
enabled = true

# Options: "duckduckgo", "google", "custom"
service = "duckduckgo"
```

| Service | URL Template | Notes |
|---------|--------------|-------|
| `duckduckgo` | `icons.duckduckgo.com/ip3/{host}.ico` | Default, privacy-focused |
| `google` | `google.com/s2/favicons?domain={host}` | Reliable, supports size param |
| `custom` | User-provided template | For self-hosted services |

### Custom Template

Use a custom favicon service:

```toml
[markata-go.link_avatars]
enabled = true
service = "custom"
template = "https://favicon.splitbee.io/?url={origin}"
```

Template placeholders:
- `{host}` - Domain name (e.g., `github.com`)
- `{origin}` - Full origin URL-encoded (e.g., `https%3A%2F%2Fgithub.com`)

### Ignore Rules

Exclude specific links from getting avatars. Links that wrap images are always skipped.

```toml
[markata-go.link_avatars]
enabled = true

# Skip links to these domains
ignore_domains = ["localhost", "127.0.0.1", "example.com"]

# Skip links to these exact origins
ignore_origins = ["https://internal.example.com"]

# Skip links matching these CSS selectors
ignore_selectors = ["nav a", ".footer a", "a.plain-link"]

# Skip links with these classes
ignore_classes = ["no-avatar", "internal"]

# Skip links inside elements with these IDs
ignore_ids = ["site-nav", "footer"]
```

## Full Configuration Example

```toml
[markata-go.link_avatars]
enabled = true
selector = "article a[href^='http']"  # Only article links
service = "google"
size = 14
position = "after"
ignore_domains = ["localhost", "127.0.0.1"]
ignore_classes = ["no-avatar"]
ignore_selectors = ["nav a", ".sidebar a"]
```

## CSS Customization

The plugin adds these classes to enhanced links:

- `.has-avatar` - Applied to all links with avatars
- `.has-avatar-before` - Avatar appears before text
- `.has-avatar-after` - Avatar appears after text

### Override Default Styles

```css
/* Change avatar appearance */
a.has-avatar::before,
a.has-avatar::after {
  opacity: 0.7;
  border-radius: 2px;
  margin: 0 0.4em;
}

/* Add hover effect */
a.has-avatar:hover::before,
a.has-avatar:hover::after {
  opacity: 1;
  transform: scale(1.1);
}

/* Hide avatars in specific contexts */
.prose a.has-avatar::before,
.prose a.has-avatar::after {
  display: none;
}
```

### Dark Mode Considerations

The CSS handles both light and dark modes automatically. For custom styling:

```css
/* Dark mode avatar styling */
@media (prefers-color-scheme: dark) {
  a.has-avatar::before,
  a.has-avatar::after {
    filter: brightness(1.2);
  }
}
```

## Opt-Out Per Link

Add the `no-avatar` class (or your configured ignore class) to skip specific links:

```markdown
[Regular link](https://github.com) gets an avatar.

[No avatar link](https://github.com){.no-avatar} is excluded.
```

Or use inline HTML:

```html
<a href="https://example.com" class="no-avatar">Plain link</a>
```

## Performance

The plugin is designed for minimal performance impact:

- **js mode**: client-side, lazy favicon loading via Intersection Observer
- **local/hosted modes**: build-time caching, no runtime JS or external requests
- **Deterministic builds** - Generated assets are stable between builds

## Generated Files

When enabled, the plugin creates:

```
output/
├── css/
│   └── link-avatars.css  # Styling
├── js/
│   └── link-avatars.js   # Client-side JavaScript (js mode only)
└── assets/
    └── markata/
        └── link-avatars/     # Cached icons (local/hosted modes)
            └── example.com.ico
```

These are automatically included in your page's `<head>`.

## Troubleshooting

### Favicons Not Showing

1. **Check browser console** for errors
2. **Verify links match selector** - Default is `a[href^='http']`
3. **Check ignore rules** - Link might match an ignore pattern
4. **Test favicon service** - Try the URL directly in browser

### Favicon Service Blocked

Some ad blockers may block favicon services. Consider:
- Using a self-hosted favicon service
- Adding exceptions for trusted services
- Providing fallback styling

### Layout Shift

If avatars cause layout shift:

```css
/* Reserve space for avatar */
a.has-avatar {
  padding-left: 1.5em;
  position: relative;
}

a.has-avatar::before {
  position: absolute;
  left: 0;
}
```

## See Also

- [Configuration Reference](/docs/guides/configuration/)
- [Plugin Development](/docs/guides/plugin-development/)
- [Themes and Styling](/docs/guides/themes/)
### Build-Time Modes

Use build-time caching to avoid runtime JavaScript and third-party fetches:

```toml
[markata-go.link_avatars]
enabled = true
mode = "local"
```

Hosted mode uploads the cached icons to your CDN and points HTML to that base URL:

```toml
[markata-go.link_avatars]
enabled = true
mode = "hosted"
hosted_base_url = "https://cdn.example.com/markata/link-avatars"
```
