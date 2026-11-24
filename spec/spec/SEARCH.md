# Search Specification

This document specifies the search functionality for markata-go sites.

Pagefind is the default static search implementation. Bleve-backed search is also supported for server-backed and remote-hosted deployments.

## Overview

Search is **enabled by default** for all markata-go sites. The implementation uses [Pagefind](https://pagefind.app/), a static site search library that generates an optimized search index during the build process.

### TUI list filtering

The terminal UI uses `/` as a lightweight filter for the currently visible post
or feed list. It filters as the user types and never searches outside the
current list. Posts match title, path, description, tags, and content; feeds
match name, title, and path. `Esc` cancels the filter and restores the visible
list; `Enter` keeps the matching list selected.

Fast mode (`build --fast`, `serve --fast`) skips Pagefind indexing.

## Architecture

### Build Process

1. **Write Stage** - HTML files are generated normally
2. **Cleanup Stage** - Pagefind plugin runs after all files are written:
   - Executes `pagefind` CLI to index the output directory
   - Generates search index in `_pagefind/` directory
   - Copies Pagefind UI assets
3. **Incremental optimization** - markata-go hashes the searchable corpus and
   skips Pagefind entirely when the indexed content is unchanged.

### Runtime Behavior

1. Pagefind JavaScript loads asynchronously
2. Search UI renders in configured position (navbar by default)
3. User queries are matched against the pre-built index
4. Results display with excerpts and highlighting

For navbar search, the results panel overlays page content and remains above
the article beneath it. The header must establish a stacking context above the
main content; increasing the panel's own `z-index` cannot escape a lower
ancestor stacking context.

### Private Content Rules

Private posts use a reduced search surface.

- Search may expose an explicit frontmatter title.
- Search may expose an explicit frontmatter description.
- Search may expose navigational fields needed to open the post.
- Search must not expose post body content.
- Search must not expose frontmatter media fields such as `image`, `cover_image`, `og_image`, `video`, or derived poster URLs.
- Search must not expose tags or derived content fields such as word count or read time.

If a private post has no explicit frontmatter title or description, search must not synthesize those fields for the result.

## Configuration

### Basic Configuration

```toml
[search]
enabled = true              # Default: true
position = "navbar"         # navbar, sidebar, footer, custom
placeholder = "Search..."   # Search input placeholder text
```

### Advanced Configuration

```toml
[search]
enabled = true
position = "navbar"
placeholder = "Search..."
show_images = true          # Show thumbnails in results (default: true)
excerpt_length = 200        # Characters for result excerpts (default: 200)
ranking_boost_title = 2.0   # Boost title matches (default: 2.0)
ranking_boost_heading = 1.5 # Boost heading matches (default: 1.5)

# Pagefind CLI options
[search.pagefind]
bundle_dir = "_pagefind"    # Output directory for index (default: "_pagefind")
exclude_selectors = [".no-search", "nav", "footer"]  # Elements to exclude
root_selector = "main"      # Element containing searchable content

# Automatic binary installation (new in v0.2.x)
auto_install = true         # Automatically download Pagefind if not in PATH (default: true)
version = "latest"          # Version to install: "latest" or specific like "v1.4.0" (default: "latest")
cache_dir = ""              # Custom cache directory (default: XDG cache ~/.cache/markata-go/bin/)
```

### Feed-Specific Search

For sites with multiple content types, configure feed-specific search instances:

```toml
# Global search (all published content)
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

### Frontmatter Control

Posts can override search behavior:

```yaml
---
title: "My Post"
search_exclude: true        # Exclude from search index
search_instance: "guides"   # Use specific search instance
search_boost: 2.0          # Boost this page in results
---
```

## Data Model

### SearchConfig

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | optional boolean | `true` | Controls whether search is active |
| `position` | string | `"navbar"` | Where search UI appears: "navbar", "sidebar", "footer", "custom" |
| `placeholder` | string | `"Search..."` | Search input placeholder text |
| `show_images` | optional boolean | `true` | Show thumbnails in search results |
| `excerpt_length` | integer | `200` | Character limit for result excerpts |
| `pagefind` | PagefindConfig | (see below) | Pagefind CLI configuration |
| `feeds` | list of SearchFeedConfig | `[]` | Feed-specific search instances |

### PagefindConfig

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `bundle_dir` | string | `"_pagefind"` | Output directory for search index |
| `exclude_selectors` | list of strings | `[]` | CSS selectors for elements to exclude from indexing |
| `root_selector` | string | `""` | CSS selector for the searchable content container |
| `auto_install` | optional boolean | `true` | Enable automatic Pagefind binary installation |
| `version` | string | `"latest"` | Pagefind version to install: "latest" or specific (e.g., "v1.4.0") |
| `cache_dir` | string | `""` | Directory for caching Pagefind binaries (default: XDG cache) |

### SearchFeedConfig

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Search instance identifier |
| `filter` | string | Filter expression for posts in this search |
| `position` | string | Where this search UI appears |
| `placeholder` | string | Search input placeholder text |

### Default Values

When creating a new SearchConfig, these defaults apply:
- `enabled`: true
- `position`: "navbar"
- `placeholder`: "Search..."
- `show_images`: true
- `excerpt_length`: 200
- `pagefind.bundle_dir`: "_pagefind"
- `pagefind.exclude_selectors`: []
- `pagefind.root_selector`: ""
- `feeds`: []

## Template Integration

### Data Attributes

Templates add Pagefind data attributes to enable indexing:

```html
<!-- Main content area -->
<main data-pagefind-body>
    <article>
        <h1 data-pagefind-meta="title">{{ post.title }}</h1>
        <div data-pagefind-meta="excerpt" class="description">
            {{ post.description }}
        </div>

        <!-- Filter by feed membership -->
        {% for feed in post.feeds %}
        <span data-pagefind-filter="feed" style="display:none">{{ feed }}</span>
        {% endfor %}

        <!-- Filter by tags -->
        {% for tag in post.tags %}
        <span data-pagefind-filter="tag" style="display:none">{{ tag }}</span>
        {% endfor %}

        <div class="content">
            {{ post.article_html | safe }}
        </div>
    </article>
</main>
```

### Excluding Content

```html
<!-- Exclude specific elements -->
<nav data-pagefind-ignore>
    <!-- Navigation not indexed -->
</nav>

<aside data-pagefind-ignore="all">
    <!-- Sidebar and all descendants excluded -->
</aside>
```

### Search UI Component

The search component is included via template:

```html
{% if config.search.enabled %}
{% include "components/search.html" %}
{% endif %}
```

### components/search.html

```html
{# Search component template #}
{# Requires: config.search (SearchConfig) #}
{# Pagefind CSS/JS are lazy-loaded on user interaction, not eagerly #}

{% with search = config.search %}
{% if search.enabled %}
<div id="search" class="search search--{{ search.position | default:'navbar' }}">
    <div id="pagefind-search">
        <input type="text" class="pagefind-ui__search-input"
               placeholder="{{ search.placeholder | default:'Search... ( / )' }}" readonly>
    </div>
</div>

<script>
(function() {
    var loaded = false;
    function loadPagefind() {
        if (loaded) return;
        loaded = true;
        var link = document.createElement('link');
        link.rel = 'stylesheet';
        link.href = '/{{ search.pagefind.bundle_dir | default:"_pagefind" }}/pagefind-ui.css';
        document.head.appendChild(link);
        var script = document.createElement('script');
        script.src = '/{{ search.pagefind.bundle_dir | default:"_pagefind" }}/pagefind-ui.js';
        script.onload = function() {
            new PagefindUI({
                element: "#pagefind-search",
                showImages: {{ search.show_images | default:true | lower }},
                excerptLength: {{ search.excerpt_length | default:200 }},
                translations: {
                    placeholder: "{{ search.placeholder | default:'Search...' }}"
                }
            });
        };
        document.head.appendChild(script);
    }
    var searchEl = document.getElementById('search');
    if (searchEl) {
        searchEl.addEventListener('mouseenter', loadPagefind, {once: true});
        searchEl.addEventListener('focusin', loadPagefind, {once: true});
    }
    document.addEventListener('keydown', function(e) {
        var tag = (e.target || {}).tagName;
        if (tag === 'INPUT' || tag === 'TEXTAREA' || tag === 'SELECT') return;
        if (e.key === '/') { e.preventDefault(); loadPagefind(); }
        if ((e.ctrlKey || e.metaKey) && e.key === 'k') { e.preventDefault(); loadPagefind(); }
    });
})();
</script>
{% endif %}
{% endwith %}
```

## Plugin Implementation

### PagefindPlugin

The plugin runs in the **Cleanup stage** (after Write) to ensure all HTML files exist.

**Behavior:**

1. Check if search is enabled in configuration
2. If Pagefind binary is not available:
   - If auto-install is enabled, download and cache the binary
   - Otherwise, log a warning with installation instructions
3. Execute Pagefind CLI with configured options to generate search index

**Plugin interfaces:**

The plugin MUST implement:
- `Plugin` - Basic plugin interface with `Name()` method returning "pagefind"
- `CleanupPlugin` - To run after all files are written
- `PriorityPlugin` - To run with last priority in the cleanup stage

**Pagefind CLI arguments:**

| Argument | Source |
|----------|--------|
| `--site` | config.output_dir |
| `--output-subdir` | search_config.pagefind.bundle_dir |
| `--root-selector` | search_config.pagefind.root_selector (if set) |
| `--exclude-selectors` | search_config.pagefind.exclude_selectors (for each) |

## CSS Styling

### Default Search Styles

```css
/* Search container positioning */
.search {
    position: relative;
}

.search--navbar {
    display: flex;
    align-items: center;
    margin-left: auto;
}

.search--sidebar {
    margin-bottom: var(--spacing-lg);
}

.search--footer {
    margin-top: var(--spacing-lg);
}

/* Pagefind UI overrides for theme integration */
:root {
    --pagefind-ui-scale: 0.9;
    --pagefind-ui-primary: var(--primary);
    --pagefind-ui-text: var(--text);
    --pagefind-ui-background: var(--surface);
    --pagefind-ui-border: var(--border);
    --pagefind-ui-tag: var(--secondary);
    --pagefind-ui-border-width: 1px;
    --pagefind-ui-border-radius: var(--radius);
    --pagefind-ui-font: var(--font-family);
}

/* Compact mode for navbar */
.search--navbar .pagefind-ui__search-input {
    width: 200px;
    transition: width 0.2s ease;
}

.search--navbar .pagefind-ui__search-input:focus {
    width: 300px;
}

/* Results dropdown */
.search--navbar .pagefind-ui__results {
    position: absolute;
    top: 100%;
    right: 0;
    width: 400px;
    max-height: 80vh;
    overflow-y: auto;
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: var(--radius);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}
```

## Accessibility

The search component follows WCAG 2.1 guidelines:

- Proper ARIA labels and roles
- Keyboard navigation support (Tab, Enter, Escape, Arrow keys)
- Focus management for results
- Screen reader announcements for results count
- Sufficient color contrast in all themes

## Performance

### Index Size Optimization

Pagefind automatically optimizes the search index:
- Content is chunked and compressed
- Only relevant chunks are loaded on search
- Typical index size: ~1KB per page indexed

### Loading Strategy

Pagefind resources are **lazy-loaded on user interaction** to avoid blocking page load:

- **CSS**: Injected as a `<link>` tag on first interaction
- **JS**: Injected as a `<script>` tag on first interaction
- **Search index**: Loaded lazily by Pagefind on first query

**Trigger events** (any of these loads Pagefind):
1. `mouseenter` or `focusin` on the search container
2. Pressing `/` (slash) key anywhere on the page
3. Pressing `Ctrl+K` (or `Cmd+K` on macOS)

**Implementation:**

```html
<div id="search" class="search">
    <div id="pagefind-search">
        <input type="text" placeholder="Search... ( / )" readonly>
    </div>
</div>

<script>
(function() {
    var loaded = false;
    function loadPagefind() {
        if (loaded) return;
        loaded = true;

        // Inject CSS
        var link = document.createElement('link');
        link.rel = 'stylesheet';
        link.href = '/_pagefind/pagefind-ui.css';
        document.head.appendChild(link);

        // Inject JS and initialize
        var script = document.createElement('script');
        script.src = '/_pagefind/pagefind-ui.js';
        script.onload = function() {
            new PagefindUI({
                element: "#pagefind-search",
                showImages: true
            });
        };
        document.head.appendChild(script);
    }

    // Load on interaction with search area
    var search = document.getElementById('search');
    if (search) {
        search.addEventListener('mouseenter', loadPagefind, { once: true });
        search.addEventListener('focusin', loadPagefind, { once: true });
    }

    // Load on keyboard shortcuts
    document.addEventListener('keydown', function(e) {
        if (e.key === '/' && !isInputFocused()) loadPagefind();
        if ((e.ctrlKey || e.metaKey) && e.key === 'k') loadPagefind();
    });
})();
</script>
```

**Benefits:**
- Saves ~30-50KB of JS and ~5KB of CSS on initial page load
- Search index chunks load only when a query is made
- Results cached for repeated queries

## Error Handling

### Pagefind Not Found

If `pagefind` CLI is not found in PATH:
1. **Auto-install enabled (default)**: Automatically downloads and caches Pagefind binary
   - Downloads from official GitHub releases
   - Verifies SHA256 checksum before execution
   - Caches in XDG cache directory for subsequent builds
2. **Auto-install disabled**: Warning logged with installation instructions
   - Site functions normally without search
   - Search UI may show placeholder message

### Auto-Install Failures

If automatic installation fails:
1. Network error: Falls back to PATH check, warns user
2. Checksum mismatch: Aborts for security, warns user
3. Unsupported platform: Warns user, suggests manual installation
4. Disk space/permission issues: Warns user with specific error

### Empty Index

If no content is indexed:
1. Warning logged during build
2. Search UI shows "No results" for all queries

### JavaScript Disabled

When JavaScript is disabled:
1. Search UI hidden via CSS
2. Site remains fully navigable
3. Consider adding sitemap link as fallback

## Security

### Binary Verification

All auto-installed Pagefind binaries are verified:
1. Downloaded from official CloudCannon GitHub releases only
2. SHA256 checksum fetched and verified before extraction
3. Binaries are cached with version-specific directories
4. Executable permissions set appropriately per platform

### Supported Platforms

| OS | Architecture | Asset Name |
|----|--------------|------------|
| macOS | x86_64 (Intel) | x86_64-apple-darwin |
| macOS | arm64 (Apple Silicon) | aarch64-apple-darwin |
| Linux | x86_64 | x86_64-unknown-linux-musl |
| Linux | arm64 | aarch64-unknown-linux-musl |
| Windows | x86_64 | x86_64-pc-windows-msvc |
| FreeBSD | x86_64 | x86_64-unknown-freebsd |

## Future Enhancements

1. **Search Analytics** - Track popular queries
2. **Instant Search** - Search-as-you-type with debouncing
3. **Search Suggestions** - Autocomplete based on content
4. **Federated Search** - Search across multiple sites
5. **Voice Search** - Web Speech API integration

## Bleve Search Server Plan

### Goals

The bleve search server exists to support deployment shapes that Pagefind does not cover well:

- remote-hosted search APIs such as `search.example.com`
- shared navbar search across multiple sites or environments
- server-side ranking and filtering for large sites
- production deployments where search should scale independently from the main site

The bleve implementation MUST preserve the same privacy expectations as the rest of the system:

- private posts MAY be searchable by intentionally public metadata such as title, description, tags, and date
- private post body content MUST NOT be indexed
- private media URLs, poster URLs, and other decryptable assets MUST NOT be exposed by search results

### Required Server Modes

The standalone search server MUST support three operating modes.

#### 1. `runtime-index`

This is the local-writer mode.

- The server loads config and content at startup
- The server builds or refreshes a local bleve index as needed
- The server MAY write cache files and index files
- This mode is suitable for local development, simple VPS deployments, and single-instance containers

#### 2. `watch-content`

This is the live-reload server mode for long-running processes.

- The server loads config and content at startup
- The server watches configured content and config roots for changes
- The server ignores generated and cache directories such as `output/`, `.markata/`, `.markata-cache/`, `cache/`, `public/`, and `markout/`
- On change, the server reloads posts, updates its in-memory view, and rebuilds or refreshes its local index
- Rebuilds MUST be debounced and serialized to prevent rebuild storms
- The server SHOULD eagerly warm the rebuilt index before marking the refreshed state ready for queries so the next user request does not pay the rebuild/open cost

#### 3. `read-only-index`

This is the production scale-out mode.

- The server opens a prebuilt index artifact and serves queries without rebuilding
- The server MUST NOT write the bleve index directory or adjacent hash files in this mode
- The server SHOULD fail fast on startup when the index is missing, unreadable, or version-incompatible
- The server SHOULD be safe to run behind a load balancer across multiple pods when all replicas mount the same read-only artifact
- The server SHOULD open and validate the configured index before readiness succeeds so the first search request is not the readiness probe for index health

### Result Hydration Requirements

The current search API hydrates results from loaded posts in memory. To support `read-only-index`, bleve documents MUST eventually store all fields needed to serve API responses directly.

The stored result schema MUST include at least:

- title
- path
- slug
- href
- description
- date
- tags
- word count or read-time source fields
- public media fields for non-private posts only
- private flag when the result represents a private post discoverable by metadata

In `read-only-index` mode the server SHOULD be able to answer queries without loading the full site content tree.

### Index Ownership and Storage Rules

#### Writable indexes

When running in `runtime-index` or `watch-content` mode:

- the search server MUST use a configurable cache dir and index dir
- the search server MUST use a configurable `index_name`
- defaults SHOULD avoid collisions between `serve`, `search-server`, and future workers

#### Shared PVCs

The system MUST NOT assume that multiple writers can safely mutate the same bleve index directory on a shared PVC.

The implementation and docs MUST explicitly distinguish these deployment shapes:

- safe: one writer per writable index path
- safe: many readers mounting the same read-only prebuilt index artifact
- unsafe: builder pods and search pods concurrently mutating the same writable bleve index path
- unsafe: multiple search server replicas sharing the same writable index path without unique index names

### Content Update Model

The standalone search server MUST document and support clear update semantics.

#### `runtime-index`

- Content changes on disk are NOT automatically observed
- A restart or explicit rebuild trigger is required

#### `watch-content`

- Content changes on disk are observed by the watcher
- The server updates its in-memory post set and refreshes the local index automatically

#### `read-only-index`

- Content changes on disk do not affect the running server directly
- A separate builder or indexer is responsible for publishing a new index artifact
- The server reloads that artifact only via restart, explicit reopen, or future hot-swap support

### Required CLI Surface

The CLI MUST evolve to support standalone production deployments.

Required flags or equivalent config for `search-server`:

- `--mode=runtime-index|watch-content|read-only-index`
- `--cache-dir`
- `--index-dir`
- `--index-name`
- `--rebuild-index`
- `--watch-debounce`

Recommended future command:

- `markata-go search build-index`

The index-build command SHOULD:

- build the bleve index without starting the HTTP server
- produce a reusable artifact for containers or object storage
- optionally emit a manifest with schema version and content fingerprint

### Remote Search Endpoint Integration

The search UI MUST support a config-driven remote endpoint.

Example target configuration:

```toml
[search.bleve]
endpoint = "https://search.example.com/api/search"
cors_origins = ["https://example.com"]
```

Behavior requirements:

- when a bleve endpoint is configured, the navbar search UI SHOULD use the bleve client instead of Pagefind
- Pagefind SHOULD remain the fallback when no endpoint is configured
- dev-mode endpoint injection used by `serve` is a convenience, not the long-term production mechanism

### Kubernetes and Container Deployment Targets

The implementation and docs MUST support three recommended deployment patterns.

#### Pattern A: Single search pod with local writable cache

- content and config mounted into the container
- bleve index stored on ephemeral or pod-local writable storage
- restart pod after new content is deployed
- When the source tree contains cache symlinks that target `/data/cache`, the Helm
  search pod MAY mount a separate pod-local `emptyDir` at `/data/cache`.
  `search.localCache.enabled` opts into this mount and defaults to `false`.
  It does not change the source or index volumes, and `read-only-index` pods
  never receive the mount. Its cache contents do not survive pod replacement.

#### Pattern B: Multiple search pods with per-pod writable cache

- each pod mounts the same content/config artifact
- each pod builds and owns its own pod-local writable index
- pods are safe to load balance because they do not share a writable index path
- this is the preferred rolling-update architecture for `watch-content` deployments because old pods may keep serving while new pods pull images, build indexes, and warm their local search state

#### Pattern C: Builder/indexer plus read-only search pods

- a builder job creates the index artifact
- search pods mount that artifact read-only
- search pods run in `read-only-index` mode
- this is the preferred production scale-out architecture

### Readiness And Warmup Requirements

The standalone search server MUST expose health behavior that distinguishes process liveness from query readiness.

- Startup readiness SHOULD fail until the active index has been opened or built successfully for the configured mode.
- In `watch-content` mode, content-change rebuilds SHOULD warm the replacement in-memory/search index before subsequent requests observe the new state.
- The first real search request after startup SHOULD NOT be responsible for paying one-time index-open or index-build latency when the deployment already reported ready.

Recommended Kubernetes deployment ergonomics:

- When source-backed search or builder-admin pods share a PVC with an external
  source bootstrap job, the deployment SHOULD gate startup until the source
  population job writes its completion marker. This readiness gate MUST be
  independent of whether a scheduled build CronJob is enabled. The completion
  marker MUST be a regular file and MUST NOT be a symlink.
- source-archive deployments SHOULD allow host-specific config overrides through environment variables so one content repo can safely drive preview and production hostnames
- ingress configuration SHOULD support optional TLS termination with a configurable secret name and the primary host plus any configured aliases

### Privacy Requirements

Bleve search MUST continue to enforce these rules across all modes:

- draft and skipped posts are excluded entirely
- private posts are searchable only by public metadata
- private body content is never indexed
- private search results never expose media URLs, poster URLs, or decryptable assets
- if a private post description is returned, it MUST be explicit public metadata rather than a body-derived excerpt

### Testing Requirements

The search implementation MUST add regression coverage for:

- private results discoverable by title/tags but not body content
- private results never exposing media fields
- remote endpoint selection overriding Pagefind in the navbar
- `watch-content` ignoring generated/cache trees and rebuilding only from content/config changes
- `read-only-index` mode opening an existing index without writing
- multi-replica operation using unique writable indexes or shared read-only artifacts

### Recommended Implementation Order

1. add config-driven remote bleve endpoint selection in the frontend
2. make `search-server` cache dir, index dir, and index name configurable
3. add an explicit index build command
4. implement `read-only-index` mode
5. implement `watch-content` mode
6. add deployment docs and regression tests

## Search startup geometry

The conditional asset loader MUST initialize search after the initial DOM
load and after SPA navigation without reading navbar geometry or dispatching
a synthetic global resize. Eager CSS reserves the search dimensions. Native
viewport resize events MUST continue to reach consumers normally. Both the
Pagefind lazy-loader callback and the legacy Pagefind initializer retain
their existing priority and invocation behavior. Other providers keep their
existing initialization paths.
