---
title: "Site Search"
description: "Add full-text search to your markata-go site using Pagefind or a bleve search server"
date: 2024-01-15
published: true
slug: /docs/guides/search/
tags:
  - documentation
  - search
  - pagefind
  - bleve
---

# Site Search

markata-go supports two search architectures:

- [Pagefind](https://pagefind.app/) for static, client-side search built into the generated site
- bleve for server-backed search APIs that can run locally or on a separate host such as `search.example.com`

Pagefind is the default search implementation. Bleve is the path for remote-hosted search, Kubernetes deployments, and future server-side search features.

## Choose a Search Mode

### Pagefind

Use Pagefind when you want:

- a fully static site
- no search server to operate
- simple hosting on Netlify, Vercel, GitHub Pages, or object storage

### Bleve Search Server

Use bleve when you want:

- a search API hosted separately from the main site
- a navbar that queries `https://search.example.com/api/search`
- server-side ranking and filtering
- Kubernetes deployments with independently scalable search pods

The current bleve implementation already works well for local development and simple standalone hosting. The production roadmap adds read-only index mode, content watching, and stronger container ergonomics.

## Private Posts

Bleve search treats private posts differently from public posts.

- Private results may include an explicit frontmatter title.
- Private results may include an explicit frontmatter description.
- Private results keep the link fields needed to open the post.
- Private results do not expose body content.
- Private results do not expose image, cover, video, poster, or thumbnail URLs from frontmatter.
- Private results do not expose tags, word count, or read time.

If a private post has no explicit `title` or `description` in frontmatter, search does not synthesize those fields for the result.

## Features

- **Enabled by default** - Search works out of the box with no configuration
- **Automatic installation** - Pagefind binary is automatically downloaded if not found in PATH
- **Fast and lightweight** - Only loads the index chunks needed for each query
- **Offline capable** - Works entirely client-side, no server required
- **Theme integrated** - Automatically matches your site's color palette
- **Navbar overlay** - Results display above page content without being obscured by the current article
- **Secure** - SHA256 checksum verification for all downloaded binaries

## Quick Start

Search is enabled by default for all markata-go sites. Just build your site and the search will work:

```bash
# Build your site - Pagefind is automatically installed if needed
markata-go build
```

That's it! The search box appears in your navigation bar.

## Bleve Search Server

The standalone search API is available through:

```bash
markata-go search-server
```

Example:

```bash
markata-go search-server --host 0.0.0.0 --port 3001
markata-go search-server --mode watch-content
markata-go search-server --mode read-only-index --index-dir /data/search.bleve
```

The server exposes:

- `/api/search`
- `/health`

### Separate Search Host

You can host search separately from the main site, for example:

- main site: `https://waylonwalker.com`
- search API: `https://search.waylonwalker.com/api/search`

To make that work cleanly, the site navbar needs a config-driven bleve endpoint and the search server must allow CORS from the main site origin.

Recommended CORS configuration:

```toml
[search.bleve]
cors_origins = ["https://waylonwalker.com"]
```

The long-term production shape is:

1. host the site separately from the search API
2. configure the navbar to query the remote bleve endpoint
3. keep Pagefind as a fallback when no bleve endpoint is configured

### Server Modes

The standalone bleve server currently supports three modes:

- `runtime-index` builds or refreshes a local index from mounted content
- `watch-content` watches content/config roots and refreshes the local index when source files change
- `read-only-index` serves a prebuilt index artifact without loading site content at runtime

Operational guidance:

- `watch-content` is the best fit when each pod owns its own writable local index directory
- for Kubernetes rolling updates, prefer pod-local writable storage for `watch-content` so old pods can
  keep serving while new pods pull images, build indexes, and warm search state
- `read-only-index` remains the right fit when a separate builder/indexer publishes a shared artifact
- readiness should only succeed after the configured index has been opened or built so the first real
  query does not pay a one-time cold-start penalty

Examples:

```bash
# Build a reusable index artifact
markata-go search build-index --index-dir /data/search.bleve --hash-path /data/search.hash

# Load content and keep a local index fresh
markata-go search-server --mode watch-content

# Serve only from a prebuilt read-only index
markata-go search-server --mode read-only-index --index-dir /data/search.bleve
```

## Automatic Pagefind Installation

As of markata-go v0.2.x, Pagefind is automatically downloaded and cached when you build your site. This eliminates the need for manual installation.

### How It Works

1. markata-go checks if `pagefind` is in your PATH
2. If not found, it automatically downloads the appropriate binary for your platform
3. The binary is verified using SHA256 checksums from official GitHub releases
4. Downloaded binaries are cached in `~/.cache/markata-go/bin/` for future builds

### Supported Platforms

| OS | Architecture | Status |
|----|--------------|--------|
| macOS | x86_64 (Intel) | Supported |
| macOS | arm64 (Apple Silicon) | Supported |
| Linux | x86_64 | Supported |
| Linux | arm64 | Supported |
| Windows | x86_64 | Supported |
| FreeBSD | x86_64 | Supported |

### Configuring Auto-Installation

```toml
[search.pagefind]
# Disable auto-install (requires manual Pagefind installation)
auto_install = false

# Pin to a specific version (default: "latest")
version = "v1.4.0"

# Custom cache directory (default: XDG cache)
cache_dir = "~/.my-cache/pagefind/"
```

### Offline Environments

After the first build, Pagefind is cached locally. Subsequent builds work offline. For CI/CD environments, consider:

1. **Pre-cache in Docker image**: Include Pagefind in your CI image
2. **Use `version` pinning**: Pin to a specific version for reproducible builds
3. **Disable auto-install**: Set `auto_install = false` and install manually

## Manual Installation (Optional)

If you prefer to install Pagefind manually or need to use a custom build:

### npm (Recommended)

```bash
npm install -g pagefind
```

### Homebrew (macOS)

```bash
brew install pagefind
```

### Binary Download

Download from [Pagefind releases](https://github.com/CloudCannon/pagefind/releases).

### Disabling Auto-Install

To require manual installation:

```toml
[search.pagefind]
auto_install = false
```

If Pagefind is not installed and auto-install is disabled, markata-go logs a warning but continues building. The search UI will be hidden.

## Configuration

### Search Backends

Pagefind remains the default backend today. The bleve server uses the same top-level search config plus additional server-side settings.

Planned production-oriented bleve configuration:

```toml
[search]
enabled = true
position = "navbar"
placeholder = "Search..."

[search.bleve]
endpoint = "https://search.example.com/api/search"
fuzzy = false
limit = 20
max_limit = 100
cors_origins = ["https://example.com"]
```

When `endpoint` is configured, the navbar should use the bleve client instead of Pagefind. When it is not configured, the site should fall back to Pagefind.

### Basic Configuration

Configure search in your `markata-go.toml`:

```toml
[search]
enabled = true              # Default: true
position = "navbar"         # Where to show search: navbar, sidebar, footer
placeholder = "Search..."   # Search input placeholder text
```

### Disable Search

To disable search entirely:

```toml
[search]
enabled = false
```

### Advanced Configuration

```toml
[search]
enabled = true
position = "navbar"
placeholder = "Search..."
show_images = true          # Show thumbnails in results (default: true)
excerpt_length = 200        # Characters for result excerpts (default: 200)

# Pagefind CLI options
[search.pagefind]
bundle_dir = "_pagefind"    # Output directory for index (default: "_pagefind")
root_selector = "main"      # Element containing searchable content
exclude_selectors = [       # Elements to exclude from indexing
    ".no-search",
    "nav",
    "footer"
]

# Auto-installation options
auto_install = true         # Automatically download Pagefind (default: true)
version = "latest"          # Version to install (default: "latest")
cache_dir = ""              # Custom cache directory (default: XDG cache)
```

## Search Positions

The `position` option controls where the search UI appears:

| Position | Description |
|----------|-------------|
| `navbar` | Right side of the navigation bar (default) |
| `sidebar` | In the sidebar (if enabled) |
| `footer` | In the footer section |
| `custom` | Hidden; use for custom placement via templates |

### Custom Placement

For full control over search placement, set `position = "custom"` and include the search component manually in your templates:

```html
{% include "components/search.html" %}
```

## What Gets Indexed

### Pagefind

By default, Pagefind indexes:

- **Post content** - The main article body
- **Title** - Displayed in search results
- **Description** - Shown as result excerpt
- **Tags** - Available as filters
- **Feed membership** - Available as filters

### Excluding Content from Search

#### Via Frontmatter

Exclude a post from the search index:

```yaml
---
title: "Draft Post"
search_exclude: true
---
```

#### Via CSS Classes

Content with `data-pagefind-ignore` is excluded:

```html
<div data-pagefind-ignore>
    This content won't be indexed
</div>
```

## Privacy Rules for Bleve

Bleve search follows stricter privacy rules than a generic full-text index.

- draft and skipped posts are excluded entirely
- private posts may be searchable by intentionally public metadata such as title, description, tags, and date
- private post body content is never indexed
- private media URLs are never exposed in bleve search results
- private descriptions shown in search results must be explicit public metadata, not body-derived excerpts

This allows users to discover an encrypted post by title or tags while still requiring decryption to read the protected content.

## Search Results

### Pagefind

Results include:

- **Title** - Clickable link to the post
- **Excerpt** - Highlighted text showing where the match was found
- **Thumbnail** - Post image (if available and `show_images` is enabled)

### Filtering Results

Pagefind supports filtering by metadata. Users can click on tags or feeds in results to filter.

### Bleve

Bleve result payloads currently include:

- title
- href
- description
- date
- tags
- read time
- media for non-private posts only

The bleve navbar UI currently renders a compact card-style dropdown. It does not yet provide the same highlighted excerpt behavior as Pagefind.

## Theming

The search UI automatically inherits your site's theme colors through CSS variables:

```css
:root {
    --pagefind-ui-primary: var(--primary);
    --pagefind-ui-text: var(--text);
    --pagefind-ui-background: var(--surface);
    --pagefind-ui-border: var(--border);
}
```

### Custom Styling

Override Pagefind styles in your custom CSS:

```css
/* Larger search input */
.pagefind-ui__search-input {
    font-size: 1.1rem;
    padding: 0.75rem 1rem;
}

/* Wider results dropdown */
.search--navbar .pagefind-ui__results {
    width: 500px;
}
```

## Keyboard Navigation

The current bleve navbar implementation supports:

- `/` to focus search
- `Escape` to dismiss results
- arrow keys to move through results
- `Enter` to open the active result
- `Tab` and `Shift+Tab` to move between the input and result links

Pagefind behavior depends on the upstream Pagefind UI implementation plus markata-go shortcut integration.

| Key | Action |
|-----|--------|
| `/` | Focus search input (when implemented) |
| `Enter` | Go to first result |
| `↑` / `↓` | Navigate results |
| `Escape` | Close results dropdown |

## Performance

Pagefind is highly optimized:

- **Index size**: ~1KB per page indexed
- **Initial load**: Only loads the search UI (~10KB)
- **On search**: Loads only relevant index chunks
- **Caching**: Results cached for repeated queries

### Large Sites

For sites with thousands of pages, Pagefind handles it efficiently:

- Index is split into chunks
- Only relevant chunks are loaded per query
- Typical search latency: 10-50ms

## Troubleshooting

### Bleve result is missing from the navbar

The current bleve navbar UI shows a limited number of results. A result may exist in the API response but rank below the visible dropdown cutoff for broad queries.

Check the API directly:

```bash
curl "http://localhost:3001/api/search?q=kubernetes&limit=30"
```

If the result is present in the API but missing from the dropdown, increase the navbar result count or improve ranking for exact title and slug matches.

### Private content is leaking through search

Bleve should never expose private body content or private media.

Check these first:

1. Is the text an explicit frontmatter `description`? Those are intentionally public metadata.
2. Is the leaked field a media URL, poster URL, or body-derived excerpt? That is a bug and should be fixed in the search API or indexer.

Private posts should be discoverable by public metadata only.

### Search server did not update after content changed

The current standalone `search-server` snapshots posts at startup.

- changing content on disk does not automatically update the running server
- restart the server after content changes
- future `watch-content` mode is intended to close this gap

For production today, prefer a restart-on-deploy model.

### Shared PVC or multiple pods

Do not rely on multiple processes mutating the same writable bleve index directory.

Safer patterns:

1. one search pod with local writable cache
2. many search pods, each with its own writable local cache
3. one builder job producing a read-only index artifact consumed by many search pods

### Search box not appearing

1. **Check build logs for Pagefind output:**
   Look for `[pagefind]` messages indicating download or execution.

2. **Verify search is enabled:**
   ```toml
   [search]
   enabled = true  # Should be true (default)
   ```

3. **Check network connectivity:**
   If auto-install is enabled and Pagefind isn't cached, internet access is required.

### Auto-install failing

1. **Network issues:**
   ```
   pagefind install error during download: failed to download asset
   ```
   Check your internet connection or firewall settings.

2. **Unsupported platform:**
   ```
   pagefind install error during platform_detection: unsupported operating system
   ```
   Install Pagefind manually for unsupported platforms.

3. **Checksum verification failure:**
   ```
   pagefind install error during verify: checksum mismatch
   ```
   This is a security feature. Try clearing the cache:
   ```bash
   rm -rf ~/.cache/markata-go/bin/
   ```

4. **Fallback to manual installation:**
   ```toml
   [search.pagefind]
   auto_install = false
   ```
   Then install Pagefind manually via npm, Homebrew, or direct download.

### No search results

1. **Check if content has `data-pagefind-body`:**
   The default post template includes this. Custom templates need it too.

2. **Verify posts are published:**
   Draft posts (`published: false`) may not be indexed.

3. **Check exclude selectors:**
   Ensure your content isn't excluded by `exclude_selectors`.

### Search results missing content

1. **Check `root_selector`:**
   If set, only content within that selector is indexed.

2. **Verify template structure:**
   Content must be within `data-pagefind-body` element.

## Feed-Specific Search

For sites with multiple content types, you can configure feed-specific search instances:

```toml
[search]
enabled = true
position = "navbar"

# Feed-specific search instances
[[search.feeds]]
name = "guides"
filter = "published && feed:guides"
position = "sidebar"
placeholder = "Search guides..."

[[search.feeds]]
name = "blog"
filter = "published && feed:blog"
position = "navbar"
placeholder = "Search blog..."
```

## CI/CD Integration

Thanks to automatic Pagefind installation, CI/CD setup is simpler than ever.

### GitHub Actions

```yaml
- name: Build site
  run: markata-go build  # Pagefind auto-installs
```

Or pin to a specific version for reproducible builds:

```yaml
- name: Build site
  run: markata-go build
  env:
    MARKATA_GO_SEARCH_PAGEFIND_VERSION: v1.4.0
```

### Netlify

Add to your `netlify.toml`:

```toml
[build]
  command = "markata-go build"
```

### Vercel

Add to your `vercel.json`:

```json
{
  "buildCommand": "markata-go build"
}
```

### Docker

Pagefind is automatically cached. For faster builds, you can pre-install:

```dockerfile
# Option 1: Let auto-install handle it (simple)
RUN markata-go build

# Option 2: Pre-install for faster CI
RUN npm install -g pagefind
RUN markata-go build
```

### Search Server Containers

The current `search-server` works best in a restart-on-deploy model.

Recommended today:

1. mount content and config into the container
2. give the server a writable local cache directory
3. restart or redeploy search pods when content changes

Recommended long-term production architecture:

1. builder pod or CI job creates the bleve index artifact
2. search pods mount that artifact read-only
3. search pods run in `read-only-index` mode
4. the main site navbar queries the remote API

This avoids shared writable index problems and makes horizontal scaling safer.

## How It Works

### Pagefind

1. **Build stage**: markata-go generates HTML files with `data-pagefind-*` attributes
2. **Cleanup stage**: PagefindPlugin runs `pagefind` CLI to index the output directory
3. **Index generation**: Pagefind creates optimized search index in `_pagefind/`
4. **Runtime**: Pagefind JS loads the index and handles search queries client-side

### Bleve

1. **Startup**: `search-server` loads config and content
2. **Indexing**: bleve index is built or opened from cache
3. **Runtime**: the navbar search UI queries `/api/search`
4. **Results**: the API returns ranked JSON results with privacy filtering applied

Planned future modes:

- `runtime-index` for simple deployments
- `watch-content` for long-running local-writer servers
- `read-only-index` for builder-produced index artifacts and load-balanced search pods

## See Also

- [Configuration Reference](/docs/guides/configuration/) - Full configuration options
- [Themes](/docs/guides/themes/) - Customize search appearance
- [Pagefind Documentation](https://pagefind.app/docs/) - Official Pagefind docs

Navbar search reserves its dimensions in the initial stylesheet. Search
initialization and SPA navigation do not force a layout read or broadcast a
resize event. Real viewport changes still resize the navbar normally.
