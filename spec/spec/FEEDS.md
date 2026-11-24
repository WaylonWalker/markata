# Feed System Specification

Feeds are the core differentiator of this static site generator. A feed is a **filtered, sorted, paginated collection of posts** that can output to **multiple formats** simultaneously.

## Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                           FEED DEFINITION                            │
│  filter: "published == True and 'python' in tags"                   │
│  sort: "date"                                                        │
│  reverse: true                                                       │
│  items_per_page: 10                                                  │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                         OUTPUT FORMATS                               │
├──────────┬──────────┬──────────┬──────────┬──────────┬─────────────┤
│   HTML   │  Simple  │   RSS    │   Atom   │   JSON   │ Markdown │    Text     │ Sitemap  │
│  index   │  list    │   feed   │   feed   │   API    │   list   │   list      │   XML    │
└──────────┴──────────┴──────────┴──────────┴──────────┴─────────────┴──────────┘
```

## Why Feeds Matter

1. **One definition, many outputs** - Define a collection once, get HTML pages, RSS, Atom, JSON API, sitemap, etc.
2. **Composable** - Feeds can reference other feeds or share filters
3. **Flexible** - Every "index page" is a feed: home, archives, tags, categories, search indexes
4. **Familiar** - Mirrors how content platforms work (RSS readers, APIs, syndication)

Feed collections MUST exclude posts that did not produce renderable page output. Posts with empty rendered HTML, skipped posts, and drafts MUST NOT appear in HTML, RSS, Atom, JSON, or derived feed collections.

Private posts are excluded from every feed unless the feed sets
`include_private = true` (the compatibility alias `private = true` has the same
meaning). This rule applies before filtering, sorting, limits, offsets,
pagination, and every output format. A public feed MUST never acquire a private
post through a preset post list or an auto-generated feed.

An opted-in private feed may list private post identity and authored metadata,
but it MUST NOT publish raw Markdown, unencrypted rendered content, encryption
key names, content-derived descriptions, private media references, or other
secret values. If private article content is included in RSS, Atom, or JSON
Feed, it MUST be the encrypted wrapper produced by the encryption plugin. HTML,
Markdown, and text feed listings use the same safe metadata projection. Private
feed entries remove the wrapper's key-selection attribute, so they decrypt only
after manual password entry and do not support cross-entry password storage or
unlocking. Private feed entries are not included in public sitemaps. The
`/feeds/` directory hides private feed definitions unless
`feeds_page.show_private_feeds` explicitly lists them.

---

## Reader Preview Hierarchy

The `/reader/` page should prefer stable, curated preview treatments over noisy screenshot thumbnails.

For each reader entry, preview selection follows this order:

1. article image extracted from the entry itself
2. source image from feed avatar or feed image metadata
3. screenshot fallback from `markata-go.blogroll.fallback_image_service`
4. branded source tile with favicon, source name, and domain

For YouTube Atom feeds, reader entries should resolve the per-video thumbnail from the feed entry metadata or stable `yt:video:<id>` entry ID before falling back to source-level imagery.

Reader template contexts may expose `preview_url`, `preview_kind`, `source_image_url`, and existing source metadata so themes can style each preview type differently.

Reader entries MAY expose `image_srcset` for a supported image resizing service.
Ghost CDN raster previews on `storage.ghost.io` with an existing
`/content/images/size/w<width>/` path SHOULD include smaller width candidates,
without exceeding the original requested width. The original `image_url` MUST
remain the `src` fallback. Unknown services and GIF/SVG images MUST
retain their original URLs without invented resizing parameters. Query strings
and escaped path segments MUST be preserved. Templates and built-in fallback
rendering MUST use the same candidates, retain lazy loading and async decoding,
and use a conservative `100vw` size hint that does not introduce intrinsic-size
containment or alter existing preview geometry.
Image selection MUST work without JavaScript and retain existing error tiles.

## Blogroll Add Metadata

When `markata-go blogroll add <feed-url>` fetches metadata for a new external feed, the command MUST use site-level HTML or OpenGraph metadata for the editable title, description, and tags shown to the user.

Feed metadata may still be fetched for canonical feed identity, channel resolution, or related enrichment, but it MUST NOT be the default source for the user-facing title, description, or tags in `blogroll add`.

When the discovered root config composes additional config files via `markata-go.include`, `markata-go blogroll add` SHOULD append the new feed to the single included file that already owns the blogroll section when that owner is unambiguous.

`markata-go blogroll add` MAY also accept YouTube-specific inputs and normalize them to the canonical channel feed URL before validation and metadata fetches. Supported normalization targets SHOULD include:

- `yt:<handle>` shortcuts such as `yt:devtoolsfm`
- channel handle URLs such as `https://www.youtube.com/@devtoolsfm`
- individual video URLs such as `https://www.youtube.com/watch?v=...`

For video URLs, the command SHOULD resolve the creator channel first and then subscribe to that creator's feed instead of treating the single video page as a feed.

## Reader Types

Blogroll feeds MAY declare a reader `type` to support multiple reader streams from the same source list.

Supported feed types are:

- `written`
- `video`
- `podcast`

`markata-go blogroll add` SHOULD let the user set the type interactively and SHOULD infer an initial value when the source is obvious, for example YouTube feeds defaulting to `video`.

Reader output requirements:

- `/reader/` MUST remain the combined stream across all feed types.
- Additional typed reader pages MUST be generated at `/reader/<type>/` for every configured type present in the blogroll.
- Typed reader pagination MUST mirror the main reader pagination structure, for example `/reader/video/page/2/`.
- Reader template contexts and fallback rendering MUST expose links for switching between the combined reader and the typed readers.
- Existing feeds without an explicit `type` MUST continue to work and SHOULD default to `written`.

---

## Feed Configuration

### Basic Feed

```toml
[[markata-go.feeds]]
slug = "blog"
title = "Blog"
filter = "published == True"
sort = "date"
reverse = true
```

### Full Configuration

```toml
[[markata-go.feeds]]
# Identity
slug = "blog"                      # URL path: /blog/
title = "All Posts"                # Display title
description = "Latest blog posts"  # Meta description

# Content Selection
filter = "published == True and date <= today"
sort = "date"
reverse = true                     # Newest first
include_private = false            # Include private posts (default: false)

# Pagination
items_per_page = 10                # 0 = no pagination (all on one page)
orphan_threshold = 3               # If last page has ≤3 items, merge with previous

# Limit/Offset
limit = 0                          # Hard cap on total items (0 = unlimited)
offset = 0                         # Skip the first N items (0 = none)

# Output Formats (all optional, defaults shown)
[markata-go.feeds.formats]
html = true                        # /blog/index.html, /blog/page/2/index.html
simple_html = true                 # /blog/simple/index.html (compact list view)
rss = true                         # /blog/rss.xml
atom = true                        # /blog/atom.xml
json = true                        # /blog/feed.json
markdown = false                   # /blog.md
text = false                       # /blog.txt
sitemap = false                    # /blog/sitemap.xml

# Templates (can override per-format)
[markata-go.feeds.templates]
html = "feed.html"                 # Template for HTML pages
simple_html = "simple-feed.html"   # Template for simple list pages
card = "partials/card.html"        # Template for post cards in list
rss = "rss.xml"                    # Template for RSS (usually default)
atom = "atom.xml"                  # Template for Atom
json = "feed.json"                 # Template for JSON
sitemap = "sitemap.xml"            # Template for sitemap
```

Built-in HTML feed templates MAY be selected per feed with `templates.html`.

Example:

```toml
[[markata-go.feeds]]
slug = "shots"
title = "Shots"
filter = "templateKey == \"shots\""

[markata-go.feeds.templates]
html = "feed-photo-grid.html"
```

The built-in `feed-photo-grid.html` template renders a dense photo-first grid with square cards by default and hover/focus overlays for title and description. It MUST preserve `post.card_classes` frontmatter hooks such as `col-span-2` and `row-span-2` so selected cards can expand across the grid.

### Default Card Media Sizing

The built-in default card may render media from the `image`, `cover`,
`cover_image`, or `og_image` frontmatter fields. Because these fields can point
to images with different aspect ratios, the default card MUST render media in a
cinematic 12:5 frame:

- request a width-only media variant (for example, `with_size:"1200"`) so the
  source image keeps its natural proportions
- use matching `width` and `height` attributes on the image
- crop the image with `object-fit: cover`

The default theme MUST keep this media within the card width. The fixed-ratio
crop is intentional for compact archive and daily note cards. Other card types
may use different fixed-ratio treatments when their design calls for them.

### Split Feed Definitions Across Files

Feed definitions may be split across composed config files. Repeated feed slugs merge into one resolved feed.

```toml
# markata-go.toml
[markata-go]
include = ["config/feeds/blog.toml", "config/feeds/blog.outputs.toml"]
```

```toml
# config/feeds/blog.toml
[[markata-go.feeds]]
slug = "blog"
title = "Blog"
filter = "published == True"
sort = "date"
reverse = true
```

```toml
# config/feeds/blog.outputs.toml
[[markata-go.feeds]]
slug = "blog"
items_per_page = 15

[markata-go.feeds.formats]
rss = true
atom = true
json = false
```

The resolved config contains one `blog` feed with values merged by `slug`.

---

## Feed Model

### Feed Object

| Field | Type | Description |
|-------|------|-------------|
| `slug` | string | URL-safe identifier |
| `title` | string | Display title |
| `description` | string? | Feed description |
| `href` | string | Base URL path (`/blog/`) |
| `filter` | string | Filter expression |
| `sort` | string? | Sort field |
| `reverse` | bool | Sort direction |
| `include_private` | bool | Include private posts in feed (default: false) |
| `limit` | int | Hard cap on total items (0 = unlimited) |
| `offset` | int | Skip the first N items (0 = none) |
| `posts` | List[Post] | All matching posts (pre-pagination) |
| `page_posts` | List[Post] | Posts for current page |
| `pagination` | Pagination | Pagination info |
| `formats` | FeedFormats | Enabled output formats (html, simple_html, rss, atom, json, markdown, text, sitemap) |

### Feed Merge Rule

When composed config files contribute multiple `[[markata-go.feeds]]` entries with the same `slug`:

- the first matching slug creates the feed
- later entries with the same slug deep-merge into that feed
- nested tables such as `formats` and `templates` also merge
- later explicit values win

For each parsed feed format, an omitted key inherits the corresponding
`feed_defaults.formats` value. An explicit `false` disables that format even
when the default is `true`. This rule applies to TOML, YAML, and JSON. Direct
Go `FeedConfig` values without parsed presence information keep their existing
all-or-nothing default behavior.

### Feed Timestamps

Feed timestamps are deterministic and based on content, not build time:

- **Atom `<updated>`** and **RSS `<lastBuildDate>`** use the most recent post date in the feed.
- If no posts have dates, these fields are omitted or left empty (no build-time fallback).

This ensures no-change incremental builds produce identical feed outputs.

Feed rebuild invalidation is based on feed membership/order, pagination structure,
and semantic post fields used by syndication outputs (title, slug, href, dates,
description, tags, image, and raw content). Downstream decorative HTML transforms
such as avatar injection do not force feed rebuilds.

### WebSub Discovery

When WebSub is enabled, feed outputs include discovery links so hubs and subscribers can find the feed endpoint:

- RSS includes `<atom:link rel="self">` plus one `<atom:link rel="hub">` per configured hub.
- Atom includes `<link rel="self">` plus one `<link rel="hub">` per configured hub.

HTML pages include `<link rel="hub">` tags in the `<head>` for each hub.

### Pagination Object

| Field | Type | Description |
|-------|------|-------------|
| `current_page` | int | Current page number (1-indexed) |
| `total_pages` | int | Total number of pages |
| `total_items` | int | Total posts in feed |
| `items_per_page` | int | Posts per page |
| `has_prev` | bool | Has previous page |
| `has_next` | bool | Has next page |
| `prev_url` | string? | URL to previous page |
| `next_url` | string? | URL to next page |
| `page_urls` | List[string] | URLs for all pages |
| `pagination_type` | string | Pagination strategy (htmx, manual, js) |

### Limit + Offset

Feeds can hard-cap the result set before pagination. After filtering and sorting,
the feed applies `offset` (skip N posts) and `limit` (keep at most N posts). The
remaining posts are then paginated.

Defaults: `limit = 0` and `offset = 0` (disabled).

---

## Pagination Types

Feeds support multiple pagination strategies for different use cases and user experiences.

### Configuration

```toml
[[markata-go.feeds]]
slug = "blog"
title = "Blog"
filter = "published == True"
enabled = true                     # Enable pagination
items_per_page = 10                # Posts per page
pagination_type = "htmx"           # Pagination strategy
```

### Pagination Type Options

| Type | Description | Use Case |
|------|-------------|----------|
| `htmx` | AJAX-based loading via HTMX | Modern, seamless UX without full page reloads |
| `manual` | Traditional page links | SEO-friendly, works without JavaScript |
| `js` | Custom JavaScript pagination | Full control over pagination behavior |

---

## HTMX Pagination

HTMX pagination provides a seamless user experience by loading new pages without full page reloads. Content is fetched via AJAX and swapped into the page dynamically.

### How It Works

```
┌─────────────────────────────────────────────────────────────────────┐
│                      HTMX PAGINATION FLOW                            │
├─────────────────────────────────────────────────────────────────────┤
│  1. User clicks "Next" or page number                               │
│  2. HTMX intercepts click, sends GET request                        │
│  3. Server returns partial HTML (just the posts list)               │
│  4. HTMX swaps new content into existing container                  │
│  5. URL updates via pushState (browser history works)               │
└─────────────────────────────────────────────────────────────────────┘
```

### HTMX Configuration

```toml
[name]
htmx_version = "2.0.8"             # HTMX version to use
skip_htmx_integrity_check = false  # Skip SHA-256 verification (not recommended)

[[markata-go.feeds]]
slug = "blog"
enabled = true
pagination_type = "htmx"
items_per_page = 10
template = "feed.html"
partial_template = "feed_partial.html"
```

### HTMX Library Management

The feeds plugin automatically handles HTMX library delivery:

1. **Download**: HTMX is downloaded from unpkg CDN during build
2. **Verification**: SHA-256 integrity hash is verified against known-good versions
3. **Local hosting**: Library is served from `/static/js/htmx.min.js`
4. **Auto-injection**: Script tag is added to `<head>` via head configuration

**Supported HTMX Versions:**

| Version | SHA-256 Hash (first 16 chars) |
|---------|-------------------------------|
| 2.0.8 | `22283ef68cb75459...` |
| 2.0.7 | `60231ae6ba9db382...` |
| 2.0.6 | `b6768eed4f3af85b...` |
| 1.9.12 | `449317ade7881e94...` |
| 1.9.10 | `b3bdcf5c741897a5...` |

### Templates for HTMX Pagination

**Full page template (`feed.html`):**

```jinja2
<!DOCTYPE html>
<html>
<head>
  <title>{{ feed.config.title }}</title>
  {# HTMX is auto-injected via head configuration #}
</head>
<body>
  <h1>{{ feed.config.title }}</h1>

  <div id="feed-content" hx-swap-oob="true">
    {% include "feed_partial.html" %}
  </div>

  {% include "pagination_controls.html" %}
</body>
</html>
```

**Partial template (`feed_partial.html`):**

```jinja2
<ul class="post-list">
{% for post in feed.posts %}
  <li>
    <a href="{{ post.slug }}/">{{ post.title }}</a>
    <time>{{ post.date }}</time>
  </li>
{% endfor %}
</ul>
```

**Pagination controls:**

```jinja2
<nav class="pagination" aria-label="Pagination">
  {% if pagination.has_prev %}
  <a href="{{ feed.config.slug }}/page/{{ pagination.prev_page }}/"
     hx-get="{{ feed.config.slug }}/page/{{ pagination.prev_page }}/partial/"
     hx-target="#feed-content"
     hx-push-url="true"
     class="pagination-prev">
    Previous
  </a>
  {% endif %}

  <span class="pagination-info">
    Page {{ pagination.current_page }} of {{ pagination.total_pages }}
  </span>

  {% if pagination.has_next %}
  <a href="{{ feed.config.slug }}/page/{{ pagination.next_page }}/"
     hx-get="{{ feed.config.slug }}/page/{{ pagination.next_page }}/partial/"
     hx-target="#feed-content"
     hx-push-url="true"
     class="pagination-next">
    Next
  </a>
  {% endif %}
</nav>
```

### HTMX Attributes Reference

| Attribute | Description |
|-----------|-------------|
| `hx-get` | URL to fetch partial content from |
| `hx-target` | CSS selector for content swap target |
| `hx-push-url` | Update browser URL (enables back/forward) |
| `hx-swap` | How to swap content (default: innerHTML) |
| `hx-swap-oob` | Out-of-band swap for multiple targets |

### Output Structure

For a feed with `pagination_type = "htmx"`:

```
output/
└── blog/
    ├── index.html              # Full page (page 1)
    ├── partial/
    │   └── index.html          # Partial content (page 1)
    ├── page/
    │   ├── 2/
    │   │   ├── index.html      # Full page (page 2)
    │   │   └── partial/
    │   │       └── index.html  # Partial content (page 2)
    │   └── 3/
    │       ├── index.html
    │       └── partial/
    │           └── index.html
    ├── rss.xml
    └── sitemap.xml
```

### Security Considerations

1. **Integrity verification**: HTMX downloads are verified against known SHA-256 hashes
2. **Local hosting**: HTMX is served locally, not from CDN at runtime
3. **Path sanitization**: Feed slugs are sanitized to prevent path traversal
4. **CSP compatible**: Works with Content-Security-Policy headers

### Fallback Behavior

If JavaScript is disabled:
- Standard links work normally (full page navigation)
- `hx-get` is ignored, `href` is followed
- SEO crawlers see all pages linked normally

---

## Reader Page (`/reader/`)

The reader page aggregates the latest posts from all followed feeds into a curated, scan-friendly layout.

### Layout Contract

- Entries are grouped by calendar day in reverse chronological order.
- Each day group renders a large date rail and a two-column content grid on wide screens.
- Each entry shows a source favicon up front, the source name, publish date, title, description, and optional image.
- When a feed item points at a Hacker News or Lobsters discussion URL, the collector resolves the outbound article URL for the main card when one exists, preserves the original discussion permalink, and keeps the feed source label visible. Lobsters RSS items commonly provide the article in `<link>` and the discussion in `<guid>` or `<comments>`; these must be preserved without an unnecessary API lookup. Reader discussion labels reflect the actual aggregator.
- The layout collapses to a single column on small screens.
- On paginated reader grids, cards after the first four in each day group SHOULD use `content-visibility: auto` with a reasonable `contain-intrinsic-size` fallback, so offscreen card contents can skip initial layout/paint without client-side rendering. The first four cards remain normally rendered to preserve first paint; browsers without support render all cards normally.
- Reader preview and source images SHOULD use asynchronous decoding while keeping native lazy loading and explicit dimensions.

### Template Context

- `entries` - flat list of entry maps for compatibility
- `day_groups` - grouped reader sections with:
  - `date_iso`
  - `weekday`
  - `title`
  - `year`
  - `count`
  - `entries`
- each entry map includes `feed_url`, `feed_title`, `source_icon_url`, `published_datetime`, `published_label`, `url`, `title`, `description`, and `image_url`
- when available, `url` points at the resolved article URL rather than the discussion page for Hacker News items
- `original_url` preserves the original HN discussion link when a feed item was normalized

### Cache Refresh Command

The CLI MUST expose a reader-cache refresh command:

```bash
markata-go reader update
```

Optional flag:

```bash
markata-go reader update --concurrency 12
```

Behavior requirements:

- the command MUST refresh the external feed cache used by the reader/blogroll system without running a site build
- the command MUST force remote feed fetches even when cached entries are still within `cache_duration`
- the command MUST continue to use stale cached feed data as a fallback when a remote refresh fails
- when `--concurrency` is provided, the command MUST use that value instead of the configured `concurrent_requests` for the current run only
- the command MUST write primary results to `stdout`
- the next `markata-go build` MUST reuse the refreshed cache data through the existing blogroll cache path

When `markata-go.blogroll.refresh_on_build = false`, normal builds MUST reuse cached
reader/blogroll data without making remote refresh requests. In that mode,
`markata-go reader update` is the explicit refresh path and builds MAY use stale cached
blogroll data when it is available.

---

## Manual Pagination

Traditional pagination with full page reloads. Best for SEO and accessibility.

### Configuration

```toml
[[markata-go.feeds]]
slug = "blog"
enabled = true
pagination_type = "manual"
items_per_page = 10
```

### Template Example

```jinja2
<nav class="pagination">
  {% if pagination.has_prev %}
  <a href="{{ feed.config.slug }}/page/{{ pagination.prev_page }}/">Previous</a>
  {% endif %}

  {% for page_num in range(1, pagination.total_pages + 1) %}
  <a href="{{ feed.config.slug }}/page/{{ page_num }}/"
     {% if page_num == pagination.current_page %}class="active"{% endif %}>
    {{ page_num }}
  </a>
  {% endfor %}

  {% if pagination.has_next %}
  <a href="{{ feed.config.slug }}/page/{{ pagination.next_page }}/">Next</a>
  {% endif %}
</nav>
```

---

## JavaScript Pagination

For custom pagination behavior, a JavaScript configuration file is generated.

### Configuration

```toml
[[markata-go.feeds]]
slug = "blog"
enabled = true
pagination_type = "js"
items_per_page = 10
```

### Generated Config File

`/static/js/pagination-config.js`:

```javascript
window.paginationData = {
  "enabled": true,
  "type": "js",
  "page": 1,
  "totalPages": 5,
  "totalPosts": 47,
  "itemsShown": 10,
  "feedName": "blog",
  "hasNext": true,
  "config": {
    "pagination_type": "js",
    "posts_per_page": 10,
    "template": "feed.html"
  }
};
```

### Custom Implementation

```javascript
// pagination.js (custom implementation)
document.addEventListener('DOMContentLoaded', () => {
  const data = window.paginationData;

  if (!data || !data.enabled) return;

  // Implement infinite scroll, load more button, etc.
  if (data.hasNext) {
    loadMoreButton.addEventListener('click', () => {
      fetch(`/${data.feedName}/page/${data.page + 1}/partial/`)
        .then(r => r.text())
        .then(html => {
          container.insertAdjacentHTML('beforeend', html);
          data.page++;
          data.hasNext = data.page < data.totalPages;
        });
    });
  }
});
```

---

## Output Formats

### HTML

Paginated HTML index pages.

**Output paths:**
```
/blog/index.html           # Page 1
/blog/page/2/index.html    # Page 2
/blog/page/3/index.html    # Page 3
```

**Template context:**
```jinja2
{{ feed.title }}
{{ feed.description }}
{{ feed.page_posts }}      {# Posts for this page #}
{{ feed.pagination }}
```

**Default template:** `feed.html`

---

### Simple HTML

A compact, dense list view of posts designed for scanning many posts quickly.
Each post renders as a single line with title (anchor link), date, and reading time.

**Output paths:**
```
/blog/simple/index.html           # Page 1
/blog/simple/page/2/index.html    # Page 2
/blog/simple/page/3/index.html    # Page 3
```

**Design goals:**
- Dense list optimized for scanning many posts at a glance
- Each post fits on one line on desktop: `Title (link)  ·  Date  ·  Reading Time`
- Wraps gracefully on mobile (stacks vertically within each entry)
- Uses semantic HTML (`<ol>` list) with microformats2 (`h-feed`, `h-entry`)
- Links back to the rich HTML feed view for full card layout

**Default template:** `simple-feed.html`

**Template context:** Same as HTML feed (`feed`, `page`, `config`).

**Configuration:**

```toml
# Enable/disable simple HTML format (default: true)
[markata-go.feeds.formats]
simple_html = true                 # /blog/simple/index.html
```

**Custom template:**

```toml
[markata-go.feeds.templates]
simple_html = "simple-feed.html"   # Override the simple feed page template
```

---

### RSS 2.0

RSS feed for feed readers.

**Output path:** `/blog/rss.xml`

**Required elements:**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
    <title>{{ feed.title }}</title>
    <link>{{ config.url }}{{ feed.href }}</link>
    <description>{{ feed.description }}</description>
    <language>{{ config.lang | default('en') }}</language>
    <lastBuildDate>{{ feed.updated | rss_date }}</lastBuildDate>
    <atom:link href="{{ config.url }}{{ feed.href }}rss.xml" rel="self" type="application/rss+xml"/>

    {% for post in feed.posts[:20] %}
    <item>
      <title>{{ post.title | xml_escape }}</title>
      <link>{{ config.url }}{{ post.href }}</link>
      <guid isPermaLink="true">{{ config.url }}{{ post.href }}</guid>
      <pubDate>{{ post.date | rss_date }}</pubDate>
      <description>{{ post.description | xml_escape }}</description>
    </item>
    {% endfor %}
  </channel>
</rss>
```

**RSS date format:** `Mon, 02 Jan 2006 15:04:05 -0700`

---

### Atom

Atom feed (RFC 4287).

**Output path:** `/blog/atom.xml`

**Required elements:**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title>{{ feed.title }}</title>
  <link href="{{ config.url }}{{ feed.href }}" rel="alternate"/>
  <link href="{{ config.url }}{{ feed.href }}atom.xml" rel="self"/>
  <id>{{ config.url }}{{ feed.href }}</id>
  <updated>{{ feed.updated | atom_date }}</updated>

  {% for post in feed.posts[:20] %}
  <entry>
    <title>{{ post.title | xml_escape }}</title>
    <link href="{{ config.url }}{{ post.href }}" rel="alternate"/>
    <id>{{ config.url }}{{ post.href }}</id>
    <published>{{ post.date | atom_date }}</published>
    <updated>{{ post.date | atom_date }}</updated>
    <summary>{{ post.description | xml_escape }}</summary>
    <author>
      <name>{{ post.author | default(config.author) }}</name>
    </author>
  </entry>
  {% endfor %}
</feed>
```

**Atom date format:** `2006-01-02T15:04:05Z` (ISO 8601)

---

### JSON Feed

JSON Feed (version 1.1).

**Output path:** `/blog/feed.json`

**Structure:**
```json
{
  "version": "https://jsonfeed.org/version/1.1",
  "title": "{{ feed.title }}",
  "home_page_url": "{{ config.url }}",
  "feed_url": "{{ config.url }}{{ feed.href }}feed.json",
  "description": "{{ feed.description }}",
  "items": [
    {
      "id": "{{ config.url }}{{ post.href }}",
      "url": "{{ config.url }}{{ post.href }}",
      "title": "{{ post.title }}",
      "content_html": "{{ post.article_html }}",
      "summary": "{{ post.description }}",
      "date_published": "{{ post.date | iso8601 }}",
      "tags": {{ post.tags | tojson }}
    }
  ]
}
```

---

### Markdown

Markdown list of posts (useful for READMEs, documentation).

**Output path:** `/blog.md`

**Default format:**
```markdown
# {{ feed.title }}

{{ feed.description }}

{% for post in feed.posts %}
- [{{ post.title }}]({{ post.href }}) - {{ post.date }}
{% endfor %}
```

---

### Plain Text

Simple text list (useful for APIs, scripts, minimal readers).

**Output path:** `/blog.txt`

**Default format:**
```
{{ feed.title }}
{{ '=' * feed.title|length }}

{% for post in feed.posts %}
{{ post.title }}
{{ post.href }}
{{ post.date }}

{% endfor %}
```

#### Text Output Encoding

Text output (`.txt` files for both feeds and individual posts) MUST produce clean plain text with no HTML artifacts:

1. **No HTML entities** - Characters like `&`, `<`, `>`, `"` appear as literal characters, never as `&amp;`, `&lt;`, `&gt;`, `&quot;`
2. **No HTML tags** - All HTML tags are stripped from the output
3. **Footnote-style links** - Links from HTML content are converted to numbered footnote references:
   ```
   Read the Go docs [1] for more info.

   References:
   [1]: https://go.dev/doc/
   ```
4. **Block structure preserved** - Block-level HTML elements produce appropriate line breaks in text output

This applies to:
- Individual post `.txt` output (rendered through `default.txt` template using the `plaintext` filter)
- The `buildTextContentFallback()` code path (when no template engine is available)
- Feed listing `.txt` output already handles entity decoding via `html.UnescapeString()`

**Implementation notes:**
- The `plaintext` template filter handles HTML-to-text conversion and marks output as safe to prevent template engine re-escaping
- Plain text metadata fields (title, description) that never contain HTML should use the `safe` filter to prevent auto-escaping
- RSS/Atom/XML output is NOT affected - entities in XML output are correct and expected

---

### Sitemap

XML sitemap for search engines (follows sitemaps.org protocol).

**Output path:** `/blog/sitemap.xml`

**Default format:**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  {% for post in feed.posts %}
  <url>
    <loc>{{ config.url }}{{ post.href }}</loc>
    {% if post.date %}
    <lastmod>{{ post.date | iso8601 }}</lastmod>
    {% endif %}
    <changefreq>{{ post.changefreq | default('weekly') }}</changefreq>
    <priority>{{ post.priority | default('0.5') }}</priority>
  </url>
  {% endfor %}
</urlset>
```

**Sitemap fields:**

| Field | Source | Default |
|-------|--------|---------|
| `loc` | `config.url` + `post.href` | Required |
| `lastmod` | `post.date` or `post.updated` | Optional |
| `changefreq` | `post.changefreq` frontmatter | `weekly` |
| `priority` | `post.priority` frontmatter | `0.5` |

**Configuration options:**
```toml
[markata-go.feeds.sitemap]
include_feeds = true              # Include feed index pages in sitemap
default_changefreq = "weekly"     # daily, weekly, monthly, yearly, never
default_priority = 0.5            # 0.0 to 1.0
```

**Per-post frontmatter:**
```yaml
---
title: Important Page
changefreq: daily
priority: 0.8
---
```

**Global sitemap:**

To generate a site-wide sitemap with all posts (not just feed posts), create a feed with no filter:

```toml
[[markata-go.feeds]]
slug = "sitemap"
title = "Sitemap"
filter = "published == True"      # All published posts
sort = "date"
reverse = true
formats = { sitemap = true }      # Only sitemap output
```

This generates `/sitemap.xml` at the root when `slug = ""`:

```toml
[[markata-go.feeds]]
slug = ""                         # Root sitemap
title = "Sitemap"
filter = "published == True"
formats = { html = false, sitemap = true }
```

---

## Common Feed Patterns

### Home Page

The root index is just a feed with `slug = ""`:

```toml
[[markata-go.feeds]]
slug = ""                          # Outputs to /index.html
title = "Home"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 5
formats = { html = true, rss = false }
```

### Archive (All Posts)

```toml
[[markata-go.feeds]]
slug = "archive"
title = "Archive"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 0                 # No pagination - all on one page
formats = { html = true }
```

### Tag Pages

Generate a feed per tag dynamically:

```toml
[markata-go.feeds.auto_tags]
enabled = true
slug_prefix = "tags"               # /tags/python/, /tags/rust/
title_template = "Posts tagged '{{ tag }}'"
filter_template = "published == True and '{{ tag }}' in tags"
formats = { html = true, rss = true }
```

This auto-generates feeds for each unique tag found in posts.

### Category Pages

```toml
[markata-go.feeds.auto_categories]
enabled = true
slug_prefix = "category"
source_field = "category"          # Frontmatter field to use
title_template = "{{ category | title }}"
filter_template = "published == True and category == '{{ category }}'"
```

### Year/Month Archives

```toml
[markata-go.feeds.auto_date]
enabled = true
slug_template = "{{ date.year }}/{{ date.month }}"
title_template = "Posts from {{ date | date('%B %Y') }}"
filter_template = "published == True and date.year == {{ year }} and date.month == {{ month }}"
formats = { html = true }
```

### API Endpoint

JSON-only feed for JavaScript consumption:

```toml
[[markata-go.feeds]]
slug = "api/posts"
title = "Posts API"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 0
formats = { json = true }          # Only JSON output
```

### Search Index

Generate a JSON index for client-side search:

```toml
[[markata-go.feeds]]
slug = "search"
title = "Search Index"
filter = "published == True"
formats = { json = true }
templates.json = "search-index.json"
```

With custom template:
```json
{
  "posts": [
    {% for post in feed.posts %}
    {
      "title": {{ post.title | tojson }},
      "href": {{ post.href | tojson }},
      "content": {{ post.content | striptags | truncate(500) | tojson }},
      "tags": {{ post.tags | tojson }}
    }{% if not loop.last %},{% endif %}
    {% endfor %}
  ]
}
```

---

## Feed Processing Pipeline

```
┌─────────────────────────────────────────────────────────────────────┐
│                          FEED LIFECYCLE                              │
├─────────────────────────────────────────────────────────────────────┤
│  1. COLLECT                                                          │
│     - Read feed definitions from config                              │
│     - Generate auto-feeds (tags, categories, dates)                  │
│                                                                      │
│  2. FILTER                                                           │
│     - For each feed, run filter expression against all posts        │
│     - `template` and `templateKey` MUST behave as the same field    │
│     - Template equality checks MUST respect built-in aliases        │
│     - Store matching posts in feed.posts                             │
│                                                                      │
│  3. SORT                                                             │
│     - Sort feed.posts by sort field                                  │
│     - Apply reverse if specified                                     │
│                                                                      │
│  4. PAGINATE                                                         │
│     - Split into pages based on items_per_page                       │
│     - Generate pagination metadata                                   │
│     - Handle orphan threshold                                        │
│                                                                      │
│  5. RENDER                                                           │
│     - For each enabled format:                                       │
│       - Load format template                                         │
│       - Render with feed context                                     │
│       - Write to output path                                         │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Template Filters for Feeds

| Filter | Description | Example |
|--------|-------------|---------|
| `rss_date` | Format as RSS date | `{{ date \| rss_date }}` |
| `atom_date` | Format as Atom/ISO date | `{{ date \| atom_date }}` |
| `iso8601` | Format as ISO 8601 | `{{ date \| iso8601 }}` |
| `xml_escape` | Escape for XML | `{{ title \| xml_escape }}` |
| `cdata` | Wrap in CDATA | `{{ html \| cdata }}` |
| `tojson` | Convert to JSON | `{{ tags \| tojson }}` |
| `absolute_url` | Add base URL | `{{ href \| absolute_url }}` |

---

## Feed Discovery

Help feed readers find feeds with `<link>` tags:

```html
<!-- In base.html <head> -->
{% for feed in core.feeds %}
  {% if feed.formats.rss %}
  <link rel="alternate" type="application/rss+xml"
        title="{{ feed.title }} (RSS)"
        href="{{ feed.href }}rss.xml">
  {% endif %}
  {% if feed.formats.atom %}
  <link rel="alternate" type="application/atom+xml"
        title="{{ feed.title }} (Atom)"
        href="{{ feed.href }}atom.xml">
  {% endif %}
  {% if feed.formats.json %}
  <link rel="alternate" type="application/feed+json"
        title="{{ feed.title }} (JSON)"
        href="{{ feed.href }}feed.json">
  {% endif %}
{% endfor %}
```

### Per-Page Feed Discovery

Pages can advertise different feeds based on their context. For example, a post in a series should advertise the series feed, not the site-wide feed.

**Discovery feed selection rules:**

1. If `config.head.alternate_feeds` is configured, use those feeds (explicit override)
2. If the post has a `sidebar_feed` (belongs to a feed sidebar), use that feed
3. Otherwise, use the site default feed (root subscription feed at `/rss.xml`, `/atom.xml`)

**Template context:**

The `discovery_feed` variable is injected into template context with:

```jinja2
{% if discovery_feed %}
  {% if discovery_feed.has_rss %}
  <link rel="alternate" type="application/rss+xml"
        title="{{ discovery_feed.title }} (RSS)"
        href="{{ discovery_feed.rss_url }}">
  {% endif %}
  {% if discovery_feed.has_atom %}
  <link rel="alternate" type="application/atom+xml"
        title="{{ discovery_feed.title }} (Atom)"
        href="{{ discovery_feed.atom_url }}">
  {% endif %}
  {% if discovery_feed.has_json %}
  <link rel="alternate" type="application/feed+json"
        title="{{ discovery_feed.title }} (JSON)"
        href="{{ discovery_feed.json_url }}">
  {% endif %}
{% endif %}
```

**discovery_feed fields:**

| Field | Type | Description |
|-------|------|-------------|
| `slug` | string | Feed slug (e.g., "tags/python", "") |
| `title` | string | Feed title |
| `rss_url` | string | RSS feed URL |
| `atom_url` | string | Atom feed URL |
| `json_url` | string | JSON feed URL |
| `has_rss` | bool | RSS format enabled |
| `has_atom` | bool | Atom format enabled |
| `has_json` | bool | JSON format enabled |

---

## Built-in Subscription Feeds

markata-go automatically creates subscription feeds at the site root:

### Root Feed (slug="")

- **RSS:** `/rss.xml`
- **Atom:** `/atom.xml`
- **HTML:** Not generated (won't overwrite `/index.html`)
- **Filter:** `published == true`

### Archive Feed (slug="archive")

- **RSS:** `/archive/rss.xml`
- **Atom:** `/archive/atom.xml`
- **HTML:** Not generated
- **Filter:** `published == true`

These feeds provide consistent subscription endpoints for readers without requiring manual configuration.

### Disabling Subscription Feeds

To disable the auto-generated subscription feeds:

```toml
[markata-go]
subscription_feeds_disabled = true
```

Or define your own feeds with `slug = ""` or `slug = "archive"` to override the defaults.

### Pins Feed (slug="pins")

When subscription feeds are enabled, the plugin MUST also provide an implicit
HTML-only `/pins/` feed using `pins.html`, the filter
`published == true and (link or url)`,
and descending date order, with 100 posts per page and manual pagination.
The Pins source link MUST prefer non-empty `link` and use non-empty `url` when
`link` is absent.
A configured archive MUST remain unchanged alongside
the implicit root and Pins feeds.

A configured `pins` feed MUST remain authoritative. A published, non-draft,
non-skipped authored post owning the `pins` slug MUST prevent implicit Pins feed
injection. Generic `link` frontmatter MUST remain render-neutral outside Pins.

The default Pins layout MUST use the available viewport width rather than the
article reading width, with responsive masonry columns and an inset header.
Card footers MUST share the card surface. Hover effects MUST stay within cards
and MUST NOT reveal floating commentary or move cards.

The `pin_preview` template filter MUST return `image` and `commentary` strings.
It MUST prefer non-empty `image`, `cover`, `cover_image`, then `og_image` fields.
Otherwise it MUST reuse images in already-rendered embed previews, valid
YouTube embed thumbnails, or ordinary post images. It MUST make no network
requests. Explicit cover fields MUST remain authoritative.

Commentary MUST come from authored paragraphs and list items, excluding embed
metadata, code, media, and standalone source URLs. If no authored commentary
exists, a non-URL description MAY provide the note. Notes MUST render as escaped
text, with a short preview and native inline details expansion on click, tap,
or keyboard activation. Long words and notes MUST remain contained by the card.
Cover art and commentary MUST work without JavaScript.

See the [Pins guide](../../docs/guides/pins.md) for usage.

---

## Calendar Year Pagination

Built-in feed Calendar views and `calendar-feed.html` MUST construct only one
populated year at a time, with at most twelve month grids. The default year MUST
be the newest year with dated posts; undated posts remain in the list fallback.
Inactive years MUST NOT have calendar day DOM. Primary/list feed views MUST defer
calendar construction until Calendar is opened.

Older/newer buttons and a native year selector MUST skip empty years, disable
unavailable directions, and support keyboard input. A valid `year=YYYY` query
parameter MUST restore the selected year. Missing or invalid years MUST fall
back to the newest populated year. Year changes MUST preserve other query
parameters and browser Back/Forward MUST restore the selected year and view.

Month sizing MUST respond to the actual available calendar width rather than
viewport breakpoints that squeeze six months into an article-width column.
Each month SHOULD have at least 17rem of width when space allows, with a single
column on narrow screens and readable day targets. An annual header MUST identify
the selected year and its post count. Existing day hover/focus/click previews
and the regular feed list fallback MUST remain available.

For primary/list feeds and calendar-enabled photo grids, calendar metadata MUST
be published as a compact JSON sidecar and fetched only when the calendar view
is opened. Those HTML views MUST NOT embed a duplicate full-history calendar
list. A dedicated calendar-first feed MAY use its visible dated list as the
calendar input and no-JavaScript fallback. Without JavaScript, the regular
paginated feed list remains available.

Feed sidebar navigation on a paginated feed MUST render only the current page's
posts. Older and newer posts remain reachable through the feed's pagination
controls; the sidebar MUST NOT duplicate links for every post in the feed.

These rules bound initial HTML and DOM size while retaining full post cards,
calendar years, and ordinary feed navigation.

Calendar month panels and year selector options use theme surface/text colors. Dedicated calendar feeds may use the page wrapper width; containers at least 73.25rem wide show four month columns in three rows. Smaller containers retain responsive readable columns.

## Configuration Inheritance

Feed configuration follows a **defaults → override** pattern. Global defaults apply to all feeds, and individual feeds can override any setting.

### Inheritance Rules

```
┌─────────────────────────────────────────────────────────────────────┐
│                    CONFIGURATION RESOLUTION                          │
├─────────────────────────────────────────────────────────────────────┤
│  1. Start with built-in defaults                                     │
│  2. Apply [markata-go.feeds.defaults.*] settings                          │
│  3. Apply [[markata-go.feeds]] individual feed settings                   │
│  4. Individual feed values WIN over defaults                         │
└─────────────────────────────────────────────────────────────────────┘
```

### Configurable at Both Levels

| Setting | Global Key | Feed Override |
|---------|------------|---------------|
| Items per page | `defaults.items_per_page` | `items_per_page` |
| Orphan threshold | `defaults.orphan_threshold` | `orphan_threshold` |
| HTML format | `defaults.formats.html` | `formats.html` |
| Simple HTML format | `defaults.formats.simple_html` | `formats.simple_html` |
| RSS format | `defaults.formats.rss` | `formats.rss` |
| Atom format | `defaults.formats.atom` | `formats.atom` |
| JSON format | `defaults.formats.json` | `formats.json` |
| Markdown format | `defaults.formats.markdown` | `formats.markdown` |
| Text format | `defaults.formats.text` | `formats.text` |
| Sitemap format | `defaults.formats.sitemap` | `formats.sitemap` |
| HTML template | `defaults.templates.html` | `templates.html` |
| RSS template | `defaults.templates.rss` | `templates.rss` |
| Card template | `defaults.templates.card` | `templates.card` |
| RSS max items | `syndication.max_items` | `max_items` |
| Include content | `syndication.include_content` | `include_content` |

### Example: Defaults with Overrides

```toml
# =============================================================================
# GLOBAL FEED DEFAULTS
# =============================================================================
[markata-go.feeds.defaults]
items_per_page = 10
orphan_threshold = 3

[markata-go.feeds.defaults.formats]
html = true
rss = true
atom = false                       # Atom OFF by default
json = false                       # JSON OFF by default
markdown = false
text = false
sitemap = false                    # Sitemap OFF by default

[markata-go.feeds.defaults.templates]
html = "feed.html"
card = "partials/card.html"
rss = "rss.xml"

[markata-go.feeds.syndication]
max_items = 20
include_content = false

# =============================================================================
# INDIVIDUAL FEEDS (inherit from defaults, override as needed)
# =============================================================================

# Home page - uses most defaults, but fewer items and no RSS
[[markata-go.feeds]]
slug = ""
title = "Home"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 5                 # Override: fewer items on home
formats = { rss = false }          # Override: no RSS for home page

# Blog - uses all defaults
[[markata-go.feeds]]
slug = "blog"
title = "Blog"
filter = "published == True"
sort = "date"
reverse = true
# items_per_page: inherits 10 from defaults
# formats: inherits html=true, rss=true from defaults

# API endpoint - completely different format set
[[markata-go.feeds]]
slug = "api/posts"
title = "Posts API"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 0                 # Override: no pagination
formats = { html = false, rss = false, json = true }  # Override: JSON only

# Archive - different pagination
[[markata-go.feeds]]
slug = "archive"
title = "All Posts"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 0                 # Override: all posts on one page
formats = { rss = false }          # Override: no RSS for archive

# Tutorials - wants Atom too
[[markata-go.feeds]]
slug = "tutorials"
title = "Tutorials"
filter = "published == True and 'tutorial' in tags"
sort = "date"
reverse = true
formats = { atom = true }          # Override: enable Atom (adds to defaults)
```

### Format Override Behavior

When overriding formats, you can:

**1. Override specific formats (merge with defaults):**
```toml
# Defaults: html=true, simple_html=true, rss=true, atom=true, json=true, sitemap=true
formats = { atom = true }
# Result: html=true, simple_html=true, rss=true, atom=true, json=true, sitemap=true
```

**2. Disable specific formats:**
```toml
# Defaults: html=true, simple_html=true, rss=true, atom=true, json=true, sitemap=true
formats = { rss = false }
# Result: html=true, simple_html=true, rss=false, atom=true, json=true, sitemap=true
```

**3. Completely replace formats (explicit all):**
```toml
formats = { html = false, rss = false, atom = false, json = true }
# Result: only JSON enabled
```

### Template Override Behavior

Templates follow the same pattern:

```toml
# Global default
[markata-go.feeds.defaults.templates]
html = "feed.html"
card = "partials/card.html"

# Individual feed override
[[markata-go.feeds]]
slug = "featured"
templates = { html = "featured-feed.html" }
# Result: html="featured-feed.html", card="partials/card.html" (inherited)
```

---

## Configuration Reference

### Global Feed Defaults

```toml
[markata-go.feeds.defaults]
# Pagination
items_per_page = 10
orphan_threshold = 3

# Formats enabled by default
[markata-go.feeds.defaults.formats]
html = true
simple_html = true
rss = true
atom = true
json = true
markdown = false
text = false
sitemap = true

# Templates used by default
[markata-go.feeds.defaults.templates]
html = "feed.html"
simple_html = "simple-feed.html"
card = "partials/card.html"
rss = "rss.xml"
atom = "atom.xml"
json = "feed.json"
markdown = "feed.md"
text = "feed.txt"
sitemap = "sitemap.xml"

# Syndication settings (RSS/Atom/JSON)
[markata-go.feeds.syndication]
max_items = 20                     # Max items in RSS/Atom feeds
include_content = true             # Include full content or just summary
```

### Built-in Defaults

If no `[markata-go.feeds.defaults]` is specified, these built-in values apply:

| Setting | Built-in Default |
|---------|------------------|
| `items_per_page` | 10 |
| `orphan_threshold` | 3 |
| `formats.html` | true |
| `formats.simple_html` | true |
| `formats.rss` | true |
| `formats.atom` | true |
| `formats.json` | true |
| `formats.markdown` | false |
| `formats.text` | false |
| `formats.sitemap` | true |
| `syndication.max_items` | 20 |
| `syndication.include_content` | true |

### Feed Item Limits

| Format | Default Limit | Configurable |
|--------|---------------|--------------|
| HTML | Paginated | `items_per_page` + `limit` / `offset` |
| Simple HTML | Paginated | `items_per_page` + `limit` / `offset` |
| RSS | 20 | `syndication.max_items` + `limit` / `offset` |
| Atom | 20 | `syndication.max_items` + `limit` / `offset` |
| JSON | All | `limit` / `offset` |
| Markdown | All | `limit` / `offset` |
| Text | All | `limit` / `offset` |
| Sitemap | All | `limit` / `offset` |

---

## Complete Example

```toml
# Site config
[name]
url = "https://myblog.com"
title = "My Blog"
author = "Jane Doe"

# =============================================================================
# GLOBAL FEED DEFAULTS
# =============================================================================
[markata-go.feeds.defaults]
items_per_page = 10
orphan_threshold = 3

[markata-go.feeds.defaults.formats]
html = true
simple_html = true
rss = true
atom = true
json = true
sitemap = true

[markata-go.feeds.defaults.templates]
html = "feed.html"
card = "partials/card.html"

[markata-go.feeds.syndication]
max_items = 20
include_content = true

# =============================================================================
# INDIVIDUAL FEEDS
# =============================================================================

# Home page - fewer items, no syndication feeds
[[markata-go.feeds]]
slug = ""
title = "Recent Posts"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 5                 # Override default
formats = { rss = false, atom = false, json = false }  # Override: HTML only

# Main blog feed - uses all defaults
[[markata-go.feeds]]
slug = "blog"
title = "Blog"
description = "All blog posts"
filter = "published == True and templateKey == 'post'"
sort = "date"
reverse = true
# Inherits: items_per_page=10, all formats enabled

# Tutorials section - uses defaults
[[markata-go.feeds]]
slug = "tutorials"
title = "Tutorials"
filter = "published == True and 'tutorial' in tags"
sort = "date"
reverse = true

# Archive - all posts, no pagination, HTML only
[[markata-go.feeds]]
slug = "archive"
title = "Archive"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 0                 # Override: no pagination
formats = { rss = false, atom = false, json = false }  # Override: HTML only

# API endpoint - JSON only, no pagination
[[markata-go.feeds]]
slug = "api/posts"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 0
formats = { html = false, rss = false, atom = false, json = true }

# =============================================================================
# AUTO-GENERATED FEEDS
# =============================================================================

# Tag pages - inherit some defaults, override formats
[markata-go.feeds.auto_tags]
enabled = true
slug_prefix = "tags"
formats = { atom = false, json = false }  # Override: HTML + RSS only
filter = "published == True"
sort = "date"
reverse = true
formats = { html = false, rss = false, atom = false, json = true }
```

**Generated outputs:**
```
public/
├── index.html                     # Home (5 posts)
├── blog/
│   ├── index.html                 # Blog page 1
│   ├── page/2/index.html          # Blog page 2
│   ├── rss.xml
│   ├── atom.xml
│   └── feed.json
├── tutorials/
│   ├── index.html
│   ├── rss.xml
│   ├── atom.xml
│   └── feed.json
├── archive/
│   └── index.html                 # All posts, no pagination
├── tags/
│   ├── python/
│   │   ├── index.html
│   │   └── rss.xml
│   ├── rust/
│   │   ├── index.html
│   │   └── rss.xml
│   └── ...
└── api/
    └── posts/
        └── feed.json              # JSON API
```

---

## See Also

- [SPEC.md](./SPEC.md) - Core specification
- [CONFIG.md](./CONFIG.md) - Feed configuration details
- [TEMPLATES.md](./TEMPLATES.md) - Template system
- [DATA_MODEL.md](./DATA_MODEL.md) - Post model

## Collection Observation Freshness

Configured and automatic feeds MUST record every considered source occurrence,
including exclusions, even when publishing is skipped. A synchronous Collect
invocation MAY prepare raw paths and eligibility flags once for its stable source
sequence. It MUST NOT deduplicate occurrences by pointer or raw/normalized path.
Preset series lists MUST be prepared separately. Each invocation MUST refresh
flags and membership; prepared observations MUST NOT survive Collect or be
cached across stages/builds. Source flags remain stable during this synchronous
hook; collection does not invoke cross-feed callbacks that mutate them.
Selection MUST retain raw-path membership, the first matched index, per-feed
privacy, eligibility reason ordering, fallback reasons, and offset-before-limit
precedence. Window reasons MUST NOT mutate shared eligibility storage.

## TUI Feed Inventory

The terminal UI feed view is an inventory, not the public `/feeds/` listing. It
MUST include configured feeds, private-enabled feeds, feeds hidden from sidebar
navigation, and automatically generated feeds. Feed rows MUST identify these
states in the title without relying on color alone:

- `private` for feeds with `include_private = true`
- `hidden` for feeds with `sidebar = false`
- `automated` for feeds generated by auto-feed configuration

The feed configuration key `private = true` is accepted as a compatibility
alias for `include_private = true`.

The title for each state SHOULD use a distinct semantic color so the inventory
can be scanned quickly.

## Editorial Feed Images

Editorial feed cards MUST defer thumbnail loading, including featured weekly
lead cards. The leading content on an editorial archive can be text, so weekly
lead status alone is not enough evidence to prioritize an image request. A
template SHOULD reserve eager loading and high fetch priority for a verified
above-the-fold image that is the page's likely Largest Contentful Paint element.

The current pagination page marker MUST maintain readable text contrast across
theme palettes. When the primary palette color is too bright for white text,
the marker background SHOULD darken that color enough to retain at least 4.5:1
contrast.

## Feed Listings

The generated `/feeds/` listing MUST hide feeds with `include_private = true` by default.
Sites MAY opt specific private-enabled feeds into the listing with:

```toml
[markata-go.feeds_page]
show_private_feeds = ["tags/gratitude"]
```

The `show_private_feeds` entries are exact feed slugs without a leading site URL.
When a private-enabled feed is listed, the listing page may show feed titles,
descriptions, counts, and feed variant links, but the underlying feed outputs remain
responsible for rendering private post content through the encryption flow rather than
exposing plaintext content.

Feed listing pages MAY show a monthly publishing-rhythm sparkline for each feed.
The chart MUST keep its accessible summary and allow users to inspect monthly counts
with pointer movement and keyboard navigation. Implementations SHOULD keep the
initial SVG DOM constant in size per chart rather than creating SVG nodes for every
month. The monthly values MAY be stored as compact chart data on the SVG and used to
update a single guide and marker. Keyboard users MUST be able to focus each chart
and use the arrow keys to inspect values; the live readout MUST announce the selected
month and count. Pointer exit MUST restore the chart's latest-month summary.

### Editorial startup work

Editorial weeks after the first use native `content-visibility: auto` with a
remembered intrinsic block size. All content remains server rendered and
available without JavaScript, to browser find, and to assistive technology.
Printing must render every week without containment. The first week renders
normally so its initial paint does not depend on visibility observation.

Navigation shortcut initialization must not fetch adjacent pages. Adjacent
documents may be prefetched on explicit navigation intent (link hover/focus or
feed cycling), respecting the existing Save Data policy.

## Deferred Feed Style Geometry

Calendar/list feed widths and navbar search widths MUST be established by
eager styles before the first paint. Applying deferred calendar/search styles
or hydrating either search provider MUST preserve the visible header, feed,
and story positions and widths at mobile and desktop sizes, including feeds
with a closed sidebar. Deferred styles MAY supply optional control appearance.

## Offscreen Editorial Week Sizing

Offscreen editorial weeks reserve an intrinsic height scaled to their editorial
item count, with a conservative fallback for custom templates. The estimate
should avoid geometric overlap with following sections while actual heights
are unknown. Native reveal remembers measured heights. Find, focus, print and
no-JavaScript access remain available.
