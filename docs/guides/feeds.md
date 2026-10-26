---
title: "Feeds Guide"
description: "Deep dive into markata-go's powerful feed system for creating filtered, sorted, paginated collections"
date: 2024-01-15
published: true
slug: /docs/guides/feeds/
tags:
  - documentation
  - feeds
---

# Feeds

For a photo-first shots archive, set `html = "feed-photo-grid.html"` under the feed's `[markata-go.feeds.templates]` table. See the [photo galleries guide](/docs/guides/galleries/).

Editorial feed cards keep thumbnails lazy loaded, including featured weekly
leads. Use eager loading and high fetch priority only when a page's actual
above-the-fold image is likely to be its Largest Contentful Paint element.

Feeds are the core differentiator of markata-go. A feed is a **filtered, sorted, paginated collection of posts** that can output to **multiple formats** simultaneously from a single definition.

Feeds only include posts that produced renderable page output. If a page is skipped, remains a draft, or has no rendered HTML because its content is empty, markata-go omits it from feed pages and syndication outputs instead of publishing a broken entry.

> **Prerequisites:** Before diving into feeds, you should understand:
> - [Frontmatter Guide](/docs/guides/frontmatter/) - How post metadata works (feeds filter based on frontmatter)
> - [Configuration Guide](/docs/guides/configuration/) - Basic config file structure

## Why Feeds Matter

Feeds solve a common problem in static site generation: you need the same collection of posts in different formats for different consumers.

```
                         FEED DEFINITION
  filter: "published == True and 'python' in tags"
  sort: "date"
  reverse: true
  items_per_page: 10
                              |
                              v
                        OUTPUT FORMATS
  +-------+--------+-------+-------+-------+----------+------+---------+
  | HTML  | Simple |  RSS  | Atom  | JSON  | Markdown | Text | Sitemap |
  | index |  list  | feed  | feed  |  API  |   list   | list |   XML   |
  +-------+--------+-------+-------+-------+----------+------+---------+
```

**Benefits:**

1. **One definition, many outputs** - Define a collection once, get HTML pages, a simple list view, RSS, Atom, JSON API, sitemap, etc.
2. **Composable** - Feeds can share filters or inherit from defaults
3. **Flexible** - Every "index page" is a feed: home, archives, tags, categories, search indexes
4. **Familiar** - Mirrors how content platforms work (RSS readers, APIs, syndication)

## Embedding Feeds in Content

The `jinja_md` plugin and page templates expose two helpers that surface feed data without littering Markdown with raw HTML.

### `render_feed` snippet helper

```jinja
{{ render_feed("blog", 5, "card") }}
```

- Renders the latest posts for a feed slug using the default partial `partials/feed_preview.html`. The `card` layout now includes `partials/cards/card-router.html`, so each post renders with the card for its template (article, note, photo, etc.) and photo/shot cards display their media inside a `<figure>` plus `<figcaption>` caption. `variant` accepts `card` or `list` and defaults to `card`.
- `limit` truncates the rendered posts even if the feed definition includes more items. The helper gracefully falls back to a simple `<section>` if the template or engine is unavailable.
- You can pass options as positional args (`limit`, `variant`) in markdown/Jinja content, or with an options map in full templates. The template override key remains `template` (for example, `{"template": "partials/custom-feed.html"}` in template files).
- The rendered HTML is emitted as a single block with no blank lines, so Markdown never re-parses card text as lists, marks or code. Card excerpts use the `summary` filter for the same reason.
- An unknown feed slug renders nothing and logs a `jinja_md` warning naming the slug — check the build output if a section comes up empty.

#### Last 5 blog/shot/ping/thought example

```jinja
{{ render_feed("blog", 5, "card") }}
{{ render_feed("shot", 5, "list") }}
{{ render_feed("ping", 5, "card") }}
{{ render_feed("thought", 5, "card") }}
```

Each `render_feed` call delimits the latest posts for that feed slug. The default partial already styles cards and lists, letting Markdown stay declarative.

### `feed_posts` data helper

When you need total control over markup (e.g., a custom grid or table), use `feed_posts` to iterate over post maps.

```jinja
{% set latest_blog = feed_posts("blog", {"limit": 3}) %}
<ul class="latest-blog">
  {% for post in latest_blog %}
  <li>
    <a href="{{ post.href }}">{{ post.title }}</a>
    {% if post.date %}
    <time datetime="{{ post.date | atom_date }}">{{ post.date | date:"Jan 2, 2006" }}</time>
    {% endif %}
  </li>
  {% endfor %}
</ul>
```

`feed_posts` returns the same map shape that templates receive (`href`, `title`, `description`, `date`, `tags`, etc.), making it ideal for layouts that differ from the default card/list styling.

## Basic Feed Configuration

A feed is defined in your `markata-go.toml` using the `[[markata-go.feeds]]` array:

```toml
[[markata-go.feeds]]
slug = "blog"
title = "Blog"
filter = "published == True"
sort = "date"
reverse = true
```

This creates a feed at `/blog/` with:
- All published posts
- Sorted by date (newest first)
- HTML, simple HTML, RSS, Atom, JSON, and sitemap output (default formats)

### Minimal Example

```toml
[[markata-go.feeds]]
slug = "posts"
title = "All Posts"
filter = "published == True"
```

If you use the default config, markata-go already includes an archive feed at `/archive/`.

## Split Feeds Across Config Files

Large sites often keep each feed in its own file.

```toml
# markata-go.toml
[markata-go]
include = ["config/feeds/*.toml"]
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

You can also split one feed across multiple files:

```toml
# markata-go.toml
[markata-go]
include = ["config/feeds/blog.toml", "config/feeds/blog.outputs.toml"]
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

Feeds merge by `slug`, so both fragments resolve to one `blog` feed.

### Full Configuration Example

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
orphan_threshold = 3               # If last page has <=3 items, merge with previous

# Limit/Offset
limit = 0                          # Hard cap on total items (0 = unlimited)
offset = 0                         # Skip the first N items (0 = none)

# Output Formats
[markata-go.feeds.formats]
html = true                        # /blog/index.html, /blog/page/2/index.html
simple_html = true                 # /blog/simple/index.html
rss = true                         # /blog/rss.xml
atom = true                        # /blog/atom.xml
json = true                        # /blog/feed.json
sitemap = true                     # /blog/sitemap.xml
markdown = false                   # /blog.md
text = false                       # /blog.txt

# Custom Templates
[markata-go.feeds.templates]
html = "feed.html"                 # Template for HTML pages
card = "partials/card.html"        # Template for post cards in list
rss = "rss.xml"                    # Template for RSS
atom = "atom.xml"                  # Template for Atom
json = "feed.json"                 # Template for JSON
```

Private posts are excluded unless a feed explicitly sets `include_private =
true`. The compatibility setting `private = true` has the same effect. An
opted-in feed may show safe post metadata and encrypted article HTML, but it
never publishes raw Markdown, unencrypted content, encryption key names,
content-derived descriptions, or private media references. Private feed
definitions are hidden from `/feeds/` unless `feeds_page.show_private_feeds`
lists them.

## Filtering Posts

The `filter` field uses a Python-like expression syntax to select which posts appear in a feed.

For template-based filters, `template` and `templateKey` are interchangeable. Equality checks also respect built-in template aliases, so `templateKey == "shots"` matches posts using `shot`, `shots`, or `photo`.

### Filter Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `==` | Equal | `published == True` |
| `!=` | Not equal | `draft != True` |
| `>`, `>=` | Greater than | `date >= '2024-01-01'` |
| `<`, `<=` | Less than | `date < '2025-01-01'` |
| `in` | Contains (tags) | `'python' in tags` |
| `contains` | String contains | `tags contains 'python'` |
| `and` | Logical AND | `published == True and featured == True` |
| `or` | Logical OR | `'python' in tags or 'go' in tags` |
| `not` | Logical NOT | `not draft` |

### Filter Examples

```toml
# Published posts only
filter = "published == True"

# Exclude drafts
filter = "draft == False"

# Posts with a specific tag
filter = "'tutorial' in tags"

# Multiple tags (OR)
filter = "'python' in tags or 'go' in tags"

# Combined conditions
filter = "published == True and featured == True"

# Date range
filter = "date >= '2024-01-01' and date < '2025-01-01'"

# By slug
filter = "slug != 'about'"

# Template aliases
filter = "templateKey == 'shots' and 'family' in tags"

# Complex filter
filter = "published == True and ('python' in tags or 'go' in tags) and featured == True"
```

## Sorting Posts

Control the order of posts with `sort` and `reverse`:

```toml
[[markata-go.feeds]]
slug = "blog"
sort = "date"        # Field to sort by
reverse = true       # true = descending (newest first)
```

### Sort Fields

You can sort by any post field:

| Field | Description |
|-------|-------------|
| `date` | Publication date (most common) |
| `title` | Alphabetical by title |
| `slug` | Alphabetical by slug |
| `reading_time` | By reading time |
| Any frontmatter | Sort by custom fields |

### Sort Examples

```toml
# Newest first (most common)
sort = "date"
reverse = true

# Oldest first (chronological)
sort = "date"
reverse = false

# Alphabetical by title
sort = "title"
reverse = false

# Custom field (priority)
sort = "priority"
reverse = true
```

## Pagination

Feeds support three pagination strategies for different use cases.

The feed sidebar lists posts on the current page. Use the feed pagination
controls to reach older and newer pages; the sidebar does not repeat every
post in a long feed.

### Pagination Configuration

```toml
[[markata-go.feeds]]
slug = "blog"
items_per_page = 10           # Posts per page (0 = all on one page)
orphan_threshold = 3          # Merge last page if <= N items
pagination_type = "manual"    # Pagination strategy
```

## Limit + Offset

Use `limit` and `offset` to hard-cap a feed before pagination. This is useful for
"latest N" lists or skipping a fixed number of entries.

```toml
[[markata-go.feeds]]
slug = "blog"
filter = "published == True"
sort = "date"
reverse = true
limit = 3      # Only include 3 posts total
offset = 1     # Skip the newest post
```

## Feed Timestamps

Feed timestamps are deterministic and based on content, not build time:

- Atom `<updated>` and RSS `<lastBuildDate>` use the most recent post date in the feed.
- If no posts have dates, these fields are omitted.

This avoids no-change incremental builds rewriting feeds.

### Pagination Types

| Type | Description | Use Case |
|------|-------------|----------|
| `manual` | Traditional page links | SEO-friendly, works without JavaScript |
| `htmx` | AJAX-based loading via HTMX | Modern, seamless UX without full page reloads |
| `js` | Custom JavaScript pagination | Full control over pagination behavior |

### Manual Pagination (Default)

Standard pagination with full page reloads:

```toml
[[markata-go.feeds]]
slug = "blog"
items_per_page = 10
pagination_type = "manual"
```

Generates:
```
/blog/index.html           # Page 1
/blog/page/2/index.html    # Page 2
/blog/page/3/index.html    # Page 3
```

### HTMX Pagination

HTMX pagination provides seamless content loading without full page reloads:

```toml
[markata-go]
htmx_version = "2.0.8"             # HTMX version to use

[[markata-go.feeds]]
slug = "blog"
items_per_page = 10
pagination_type = "htmx"
template = "feed.html"
partial_template = "feed_partial.html"
```

**How it works:**
1. User clicks "Next" or page number
2. HTMX intercepts click, sends GET request
3. Server returns partial HTML (just the posts list)
4. HTMX swaps new content into existing container
5. URL updates via pushState (browser history works)

**Output structure:**
```
/blog/
  index.html                 # Full page (page 1)
  partial/
    index.html               # Partial content (page 1)
  page/
    2/
      index.html             # Full page (page 2)
      partial/
        index.html           # Partial content (page 2)
```

**HTMX template example:**

```html
<!-- feed.html (full page) -->
<div id="feed-content" hx-swap-oob="true">
  {% include "feed_partial.html" %}
</div>

<nav class="pagination">
  {% if pagination.has_prev %}
  <a href="{{ feed.slug }}/page/{{ pagination.prev_page }}/"
     hx-get="{{ feed.slug }}/page/{{ pagination.prev_page }}/partial/"
     hx-target="#feed-content"
     hx-push-url="true">
    Previous
  </a>
  {% endif %}

  {% if pagination.has_next %}
  <a href="{{ feed.slug }}/page/{{ pagination.next_page }}/"
     hx-get="{{ feed.slug }}/page/{{ pagination.next_page }}/partial/"
     hx-target="#feed-content"
     hx-push-url="true">
    Next
  </a>
  {% endif %}
</nav>
```

### JavaScript Pagination

For custom JavaScript implementations:

```toml
[[markata-go.feeds]]
slug = "blog"
items_per_page = 10
pagination_type = "js"
```

Generates a configuration file at `/static/js/pagination-config.js`:

```javascript
window.paginationData = {
  "enabled": true,
  "type": "js",
  "page": 1,
  "totalPages": 5,
  "totalPosts": 47,
  "itemsShown": 10,
  "feedName": "blog",
  "hasNext": true
};
```

### Orphan Threshold

The `orphan_threshold` prevents tiny last pages:

```toml
[[markata-go.feeds]]
slug = "blog"
items_per_page = 10
orphan_threshold = 3    # If last page has <=3 items, merge with previous
```

With 23 posts and `items_per_page = 10`:
- Without threshold: Page 1 (10), Page 2 (10), Page 3 (3)
- With threshold: Page 1 (10), Page 2 (13)

## Output Formats

Each feed can generate multiple output formats simultaneously.

### HTML

Paginated HTML index pages.

```toml
[markata-go.feeds.formats]
html = true
```

**Output:**
```
/blog/index.html           # Page 1
/blog/page/2/index.html    # Page 2
```

### Simple HTML

A compact, dense list view of posts designed to fit many entries on screen. Each post renders as a single line with title (as a link), date, and reading time. The simple feed is generated alongside the standard HTML feed automatically.

```toml
[markata-go.feeds.formats]
simple_html = true
```

**Output:**
```
/blog/simple/index.html           # Page 1
/blog/simple/page/2/index.html    # Page 2
```

The simple feed page includes a navigation switcher to toggle between the full (rich card) view and the simple list view. Each entry uses [microformats2](http://microformats.org/wiki/h-entry) markup (`h-entry`, `p-name`, `dt-published`).

**When to use the simple format:**

- **Dense archives** -- Show many posts at a glance without scrolling through large cards
- **Quick scanning** -- Readers who know what they're looking for can scan titles and dates rapidly
- **Low-bandwidth** -- Minimal HTML and no images, loads fast on slow connections
- **Print-friendly** -- Styled for clean printing via `@media print`

**Disabling the simple format** for a specific feed:

```toml
[[markata-go.feeds]]
slug = "api/posts"
formats = { simple_html = false }
```

### RSS 2.0

RSS feed for feed readers.

```toml
[markata-go.feeds.formats]
rss = true
```

**Output:** `/blog/rss.xml`

```xml
<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
    <title>Blog</title>
    <link>https://example.com/blog/</link>
    <description>Latest blog posts</description>
    <lastBuildDate>Mon, 15 Jan 2024 12:00:00 +0000</lastBuildDate>
    <atom:link href="https://example.com/blog/rss.xml" rel="self" type="application/rss+xml"/>

    <item>
      <title>My Post</title>
      <link>https://example.com/my-post/</link>
      <guid isPermaLink="true">https://example.com/my-post/</guid>
      <pubDate>Mon, 15 Jan 2024 12:00:00 +0000</pubDate>
      <description>Post description...</description>
    </item>
  </channel>
</rss>
```

### Atom

Atom feed (RFC 4287).

```toml
[markata-go.feeds.formats]
atom = true
```

**Output:** `/blog/atom.xml`

```xml
<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title>Blog</title>
  <link href="https://example.com/blog/" rel="alternate"/>
  <link href="https://example.com/blog/atom.xml" rel="self"/>
  <id>https://example.com/blog/</id>
  <updated>2024-01-15T12:00:00Z</updated>

  <entry>
    <title>My Post</title>
    <link href="https://example.com/my-post/" rel="alternate"/>
    <id>https://example.com/my-post/</id>
    <published>2024-01-15T12:00:00Z</published>
    <updated>2024-01-15T12:00:00Z</updated>
    <summary>Post description...</summary>
  </entry>
</feed>
```

### Styled RSS and Atom Pages

RSS and Atom feeds link to `/rss.xsl` and `/atom.xsl`, so a browser shows a readable page instead of raw XML. Both pages use your site's theme: the same palette, light or dark mode, fontpack, text size, and aesthetic. They also follow the visitor's own theme picker choices. When you change the theme, the feed pages change with it.

The browser view identifies itself as a feed, explains how to subscribe, and presents entries in a numbered reading list. Feed readers still receive the original RSS or Atom XML; the stylesheet only changes how a browser displays it. The same stylesheet supports either format, so custom views can keep their visual treatment consistent.

To customize a feed page, copy `rss.xsl` or `atom.xsl` into your `templates/` directory. Keep this block inside `<head>` to keep the site theme; markata-go replaces it with the theme's styles and scripts when it builds:

```xml
<!-- markata:theme-head -->
<link rel="stylesheet" href="/css/variables.css" />
<link rel="stylesheet" href="/css/main.css" />
<!-- /markata:theme-head -->
```

A custom stylesheet without this block is copied as-is. To change what is injected, override `partials/feed-xsl-head.html` in your templates.

### JSON Feed

JSON Feed (version 1.1) for modern feed readers and APIs.

```toml
[markata-go.feeds.formats]
json = true
```

**Output:** `/blog/feed.json`

```json
{
  "version": "https://jsonfeed.org/version/1.1",
  "title": "Blog",
  "home_page_url": "https://example.com",
  "feed_url": "https://example.com/blog/feed.json",
  "description": "Latest blog posts",
  "items": [
    {
      "id": "https://example.com/my-post/",
      "url": "https://example.com/my-post/",
      "title": "My Post",
      "content_html": "<p>Post content...</p>",
      "summary": "Post description...",
      "date_published": "2024-01-15T12:00:00Z",
      "tags": ["go", "tutorial"]
    }
  ]
}
```

### Markdown

Markdown list of posts (useful for READMEs, documentation).

```toml
[markata-go.feeds.formats]
markdown = true
```

**Output:** `/blog.md`

```markdown
# Blog

Latest blog posts

- [My Post](/my-post/) - 2024-01-15
- [Another Post](/another-post/) - 2024-01-10
```

### Plain Text

Simple text list (useful for APIs, scripts, minimal readers).

```toml
[markata-go.feeds.formats]
text = true
```

**Output:** `/blog.txt`

```
Blog
====

Latest blog posts

2024-01-15 - My Post
  /my-post/

2024-01-10 - Another Post
  /another-post/
```

Text output always produces clean plain text -- HTML entities are decoded to their literal characters and no HTML tags appear in the output. Links from HTML content are converted to footnote-style references. See the [Templates Guide](/docs/guides/templates/#text-templates) for details on using the `plaintext` filter.

### Sitemap

XML sitemap for search engines.

```toml
[markata-go.feeds.formats]
sitemap = true
```

**Output:** `/blog/sitemap.xml`

```xml
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>https://example.com/my-post/</loc>
    <lastmod>2024-01-15</lastmod>
    <changefreq>weekly</changefreq>
    <priority>0.5</priority>
  </url>
</urlset>
```

## Auto-Generated Tag Feeds

markata-go can automatically create feeds for each unique tag in your posts.

Automatic and configured feeds refresh selection diagnostics on every Collect,
including exclusions and feeds whose output publishing is skipped. Selection
scratch is reused only within that invocation, not cached between builds.
Privacy choices, preset series order, duplicate source occurrences, and
offset/limit exclusion reasons remain unchanged. For timing, inspect the bounded
`auto_feeds` collection debug record described in the [performance guide](performance.md);
it separates generation, filtering/sorting, diagnostic selection recording, and
pagination/preparation without logging post paths or feed names.

### Enable Auto Tags

```toml
[markata-go.feeds.auto_tags]
enabled = true
slug_prefix = "tags"               # /tags/python/, /tags/go/

[markata-go.feeds.auto_tags.formats]
html = true
rss = true
```

### Generated Output

If your posts have tags `python`, `go`, and `tutorial`:

```
/tags/
  python/
    index.html
    rss.xml
  go/
    index.html
    rss.xml
  tutorial/
    index.html
    rss.xml
```

### Auto Categories

Similar to tags, but for the `category` frontmatter field:

```toml
[markata-go.feeds.auto_categories]
enabled = true
slug_prefix = "categories"

[markata-go.feeds.auto_categories.formats]
html = true
rss = true
```

### Auto Archives

Automatic date-based archive feeds:

```toml
[markata-go.feeds.auto_archives]
enabled = true
slug_prefix = "archive"
yearly_feeds = true               # /archive/2024/
monthly_feeds = true              # /archive/2024/01/

[markata-go.feeds.auto_archives.formats]
html = true
```

## Feed Defaults and Inheritance

Configure defaults that apply to all feeds, then override as needed.

### Global Defaults

```toml
# Default values for ALL feeds
[markata-go.feeds.defaults]
items_per_page = 10
orphan_threshold = 3

[markata-go.feeds.defaults.formats]
html = true
rss = true
atom = false
json = false
markdown = false
text = false

[markata-go.feeds.defaults.templates]
html = "feed.html"
card = "partials/card.html"
rss = "rss.xml"
atom = "atom.xml"
json = "feed.json"

# Syndication settings (RSS/Atom/JSON limits)
[markata-go.feeds.syndication]
max_items = 20                     # Max items in RSS/Atom feeds
include_content = false            # Include full content or just summary
```

### Inheritance Example

```toml
# Global defaults
[markata-go.feeds.defaults]
items_per_page = 10

[markata-go.feeds.defaults.formats]
html = true
rss = true
atom = false
json = false

# Home page - override items_per_page and formats
[[markata-go.feeds]]
slug = ""
title = "Home"
filter = "published == True"
items_per_page = 5                 # Override: fewer items
formats = { rss = false }          # Override: no RSS for home

# Blog - inherits all defaults
[[markata-go.feeds]]
slug = "blog"
title = "Blog"
filter = "published == True"
# Inherits: items_per_page=10, html=true, rss=true

# API endpoint - completely different formats
[[markata-go.feeds]]
slug = "api/posts"
title = "Posts API"
filter = "published == True"
items_per_page = 0                 # Override: no pagination
formats = { html = false, rss = false, json = true }  # Override: JSON only
```

### Format Override Behavior

When overriding formats:

```toml
# Defaults: html=true, rss=true, atom=false, json=false

# Add to defaults (merge)
formats = { atom = true }
# Result: html=true, rss=true, atom=true, json=false

# Disable specific format
formats = { rss = false }
# Result: html=true, rss=false, atom=false, json=false

# Complete replacement
formats = { html = false, rss = false, atom = false, json = true }
# Result: json only
```

## Common Patterns

### Home Page Feed

The home page is a feed with an empty slug:

```toml
[[markata-go.feeds]]
slug = ""                          # Outputs to /index.html
title = "Recent Posts"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 5

[markata-go.feeds.formats]
html = true
rss = false                        # Usually no RSS on home page
```

### Archive (All Posts, No Pagination)

```toml
[[markata-go.feeds]]
slug = "archive"
title = "All Posts"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 0                 # No pagination - all posts on one page

[markata-go.feeds.formats]
html = true
rss = false
```

### Private Archive (Including Private Posts)

If you want a feed that includes private posts (useful for personal archives or admin pages):

```toml
[[markata-go.feeds]]
slug = "private-archive"
title = "Private Archive"
filter = "published == True"
sort = "date"
reverse = true
include_private = true             # Include private posts in this feed

[markata-go.feeds.formats]
html = true
rss = false                        # Usually don't syndicate private posts
```

If RSS, Atom, or JSON Feed is enabled for this feed, private entries contain
only the encrypted article wrapper and safe metadata. The feed wrapper omits
the encryption key name, so entries decrypt individually and do not share
session storage. Public feeds continue to exclude private posts.

### Tag Pages (Manual)

If you don't want auto-generated tag feeds:

```toml
[[markata-go.feeds]]
slug = "tags/python"
title = "Python Posts"
filter = "published == True and 'python' in tags"
sort = "date"
reverse = true
```

### Featured Posts

```toml
[[markata-go.feeds]]
slug = "featured"
title = "Featured"
description = "Hand-picked articles"
filter = "published == True and featured == True"
sort = "date"
reverse = true
items_per_page = 6
```

### JSON API Endpoint

```toml
[[markata-go.feeds]]
slug = "api/posts"
title = "Posts API"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 0                 # All posts

[markata-go.feeds.formats]
html = false
rss = false
json = true                        # JSON only
```

### Search Index

Generate a JSON index for client-side search:

```toml
[[markata-go.feeds]]
slug = "search"
title = "Search Index"
filter = "published == True"

[markata-go.feeds.formats]
json = true

[markata-go.feeds.templates]
json = "search-index.json"
```

Custom template (`templates/search-index.json`):

```json
{
  "posts": [
    {% for post in feed.posts %}
    {
      "title": {{ post.title | tojson }},
      "href": {{ post.href | tojson }},
      "content": {{ post.content | striptags | truncate:500 | tojson }},
      "tags": {{ post.tags | tojson }}
    }{% if not loop.last %},{% endif %}
    {% endfor %}
  ]
}
```

### Global Sitemap

Root sitemap with all posts:

```toml
[[markata-go.feeds]]
slug = ""                          # Root
title = "Sitemap"
filter = "published == True"
sort = "date"
reverse = true

[markata-go.feeds.formats]
html = false                       # No HTML
sitemap = true                     # Sitemap only
```

Outputs `/sitemap.xml` at the root.

### Tutorials Section

```toml
[[markata-go.feeds]]
slug = "tutorials"
title = "Tutorials"
description = "Step-by-step guides"
filter = "published == True and 'tutorial' in tags"
sort = "date"
reverse = true
items_per_page = 10

[markata-go.feeds.formats]
html = true
rss = true
atom = true
json = true
```

## Template Variables

When rendering feed templates, these variables are available:

### Feed Context (`feed`)

| Variable | Type | Description |
|----------|------|-------------|
| `feed.slug` | string | URL-safe identifier |
| `feed.title` | string | Display title |
| `feed.description` | string | Feed description |
| `feed.href` | string | Base URL path (`/blog/`) |
| `feed.posts` | []Post | All matching posts (pre-pagination) |

### Pagination Context (`pagination`)

| Variable | Type | Description |
|----------|------|-------------|
| `pagination.current_page` | int | Current page number (1-indexed) |
| `pagination.total_pages` | int | Total number of pages |
| `pagination.total_items` | int | Total posts in feed |
| `pagination.items_per_page` | int | Posts per page |
| `pagination.has_prev` | bool | Has previous page |
| `pagination.has_next` | bool | Has next page |
| `pagination.prev_url` | string | URL to previous page |
| `pagination.next_url` | string | URL to next page |
| `pagination.prev_page` | int | Previous page number |
| `pagination.next_page` | int | Next page number |
| `pagination.page_urls` | []string | URLs for all pages |

### Page Context (`page`)

| Variable | Type | Description |
|----------|------|-------------|
| `page.number` | int | Page number |
| `page.posts` | []Post | Posts on this page |
| `page.has_prev` | bool | Has previous page |
| `page.has_next` | bool | Has next page |
| `page.prev_url` | string | Previous page URL |
| `page.next_url` | string | Next page URL |

### Post Fields

Each post in `feed.posts` or `page.posts` has:

| Field | Type | Description |
|-------|------|-------------|
| `post.Title` | string | Post title |
| `post.Slug` | string | URL slug |
| `post.Href` | string | Relative URL (`/my-post/`) |
| `post.Date` | time | Publication date |
| `post.Published` | bool | Whether published |
| `post.Tags` | []string | List of tags |
| `post.Description` | string | Post description |
| `post.Content` | string | Raw markdown |
| `post.ArticleHTML` | string | Rendered HTML |
| `post.Extra` | map | Additional frontmatter |

### Template Example

```html
{% extends "base.html" %}

{% block content %}
<main>
  <h1>{{ feed.title }}</h1>
  {% if feed.description %}
  <p class="description">{{ feed.description }}</p>
  {% endif %}

  <ul class="posts">
  {% for post in page.posts %}
    <li class="post">
      <a href="{{ post.Href }}">
        <h2>{{ post.Title }}</h2>
      </a>
      <time datetime="{{ post.Date|atom_date }}">
        {{ post.Date|date_format:"January 2, 2006" }}
      </time>
      {% if post.Description %}
      <p>{{ post.Description }}</p>
      {% endif %}
    </li>
  {% endfor %}
  </ul>

  <nav class="pagination">
    {% if page.has_prev %}
    <a href="{{ page.prev_url }}" rel="prev">Previous</a>
    {% endif %}

    <span>Page {{ pagination.current_page }} of {{ pagination.total_pages }}</span>

    {% if page.has_next %}
    <a href="{{ page.next_url }}" rel="next">Next</a>
    {% endif %}
  </nav>
</main>
{% endblock %}
```

## Template Filters for Feeds

| Filter | Description | Example |
|--------|-------------|---------|
| `rss_date` | Format as RSS date | `{{ date\|rss_date }}` |
| `atom_date` | Format as Atom/ISO date | `{{ date\|atom_date }}` |
| `date_format` | Custom date format | `{{ date\|date_format:"2006-01-02" }}` |
| `tojson` | Convert to JSON | `{{ tags\|tojson }}` |
| `absolute_url` | Add base URL | `{{ href\|absolute_url:config.URL }}` |
| `striptags` | Remove HTML tags | `{{ html\|striptags }}` |
| `truncate` | Truncate to length | `{{ text\|truncate:100 }}` |

## Built-In Feed Templates

Feeds can opt into a different built-in HTML layout without copying templates into the site.

```toml
[[markata-go.feeds]]
slug = "shots"
title = "Shots"
filter = "templateKey == \"shots\""

[markata-go.feeds.templates]
html = "feed-photo-grid.html"
```

`feed-photo-grid.html` is a photo-grid archive designed for image-heavy feeds. It renders:

- a tight multi-column grid
- square cards by default
- title and description overlays on hover and keyboard focus
- `post.card_classes` span hooks such as `col-span-2`, `row-span-2`, and `col-span-4`

That makes it a good fit for shots, screenshots, or gallery-style feeds where individual items should sometimes break out of the default square rhythm.

## Feed Discovery

markata-go provides automatic feed discovery so visitors and feed readers can find the right feed for each page.

### How It Works

Every page automatically includes `<link rel="alternate">` tags in the `<head>` for feed discovery. The feed advertised depends on the page:

1. **Posts with a sidebar feed** - If a post participates in a feed sidebar (e.g., a series or guide), the discovery links point to that feed.
2. **Other pages** - All other pages advertise the site's default subscription feeds (`/rss.xml`, `/atom.xml`).

This ensures readers always discover the most relevant feed for the content they're viewing.

### Template Context

When rendering pages, markata-go injects a `discovery_feed` variable into the template context:

| Field | Type | Description |
|-------|------|-------------|
| `discovery_feed.title` | string | Feed title |
| `discovery_feed.has_rss` | bool | RSS format available |
| `discovery_feed.has_atom` | bool | Atom format available |
| `discovery_feed.has_json` | bool | JSON format available |
| `discovery_feed.rss_url` | string | RSS feed URL |
| `discovery_feed.atom_url` | string | Atom feed URL |
| `discovery_feed.json_url` | string | JSON feed URL |

### Default Template Behavior

The default `base.html` template handles discovery automatically:

```html
<head>
  {% if config.head.alternate_feeds %}
    {# Explicit override - user-configured feeds take precedence #}
    {% for feed in config.head.alternate_feeds %}
    <link rel="alternate" type="{{ feed.type }}" title="{{ feed.title }}" href="{{ feed.href }}">
    {% endfor %}
  {% elif discovery_feed %}
    {# Automatic per-page discovery #}
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
  {% else %}
    {# Fallback #}
    <link rel="alternate" type="application/rss+xml" title="RSS Feed" href="/rss.xml">
  {% endif %}
</head>
```

### Override with `head.alternate_feeds`

If you need full control over discovery links, configure `head.alternate_feeds` explicitly:

```toml
[[markata-go.head.alternate_feeds]]
type = "application/rss+xml"
title = "My Custom Feed"
href = "/custom/rss.xml"
```

When `head.alternate_feeds` is set, it completely overrides the automatic discovery behavior.

## Built-in Subscription Feeds

markata-go automatically creates site-wide subscription feeds so visitors always have a way to subscribe.

### Generated Feeds

By default, these feeds are created automatically:

| Path | Description |
|------|-------------|
| `/rss.xml` | Site RSS feed (all published posts) |
| `/atom.xml` | Site Atom feed (all published posts) |
| `/archive/rss.xml` | Archive RSS feed (same content as root) |
| `/archive/atom.xml` | Archive Atom feed (same content as root) |

### How It Works

The `subscription_feeds` plugin creates two internal feeds:

1. **Root feed** (`slug = ""`) - Generates `/rss.xml` and `/atom.xml` (no HTML to avoid overwriting your home page)
2. **Archive feed** (`slug = "archive"`) - Generates `/archive/rss.xml` and `/archive/atom.xml`

Both feeds contain the same items: all published posts sorted by date (newest first).

### Customizing Subscription Feeds

The built-in feeds use sensible defaults, but you can override them by defining your own feeds with the same slugs:

```toml
# Override the root subscription feed
[[markata-go.feeds]]
slug = ""
title = "My Site"
filter = "published == True and featured == True"  # Only featured posts
sort = "date"
reverse = true

[markata-go.feeds.formats]
html = false      # Still no HTML (don't overwrite home page)
rss = true
atom = true
```

### Why Both Root and Archive?

- **Root feeds** (`/rss.xml`, `/atom.xml`) - The standard location feed readers check first
- **Archive feeds** (`/archive/...`) - Alternative paths for services that expect feeds in subdirectories

Both contain identical content, ensuring maximum compatibility with feed readers and syndication services.

## Manual Feed Discovery Links

If you need to add discovery links manually in a custom template:

```html
<head>
  <!-- RSS Feed -->
  <link rel="alternate" type="application/rss+xml"
        title="{{ feed.title }} (RSS)"
        href="{{ feed.href }}rss.xml">

  <!-- Atom Feed -->
  <link rel="alternate" type="application/atom+xml"
        title="{{ feed.title }} (Atom)"
        href="{{ feed.href }}atom.xml">

  <!-- JSON Feed -->
  <link rel="alternate" type="application/feed+json"
        title="{{ feed.title }} (JSON)"
        href="{{ feed.href }}feed.json">
</head>
```

### WebSub Discovery

If you enable WebSub, markata-go also emits hub discovery links:

```html
<head>
  <link rel="hub" href="https://hub.example.com/">
</head>
```

RSS/Atom feeds include `rel="hub"` and `rel="self"` links when WebSub is enabled.

## Configuration Reference

### Feed Config Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `slug` | string | Required | URL path (`""` for root) |
| `title` | string | Required | Display title |
| `description` | string | `""` | Feed description |
| `filter` | string | `""` | Filter expression |
| `sort` | string | `"date"` | Sort field |
| `reverse` | bool | `true` | Sort direction (true=descending) |
| `include_private` | bool | `false` | Include private posts in feed |
| `items_per_page` | int | `10` | Posts per page (0=no pagination) |
| `orphan_threshold` | int | `3` | Min items for separate page |
| `pagination_type` | string | `"manual"` | `manual`, `htmx`, or `js` |
| `limit` | int | `0` | Hard cap on total items (0 = unlimited) |
| `offset` | int | `0` | Skip the first N items (0 = none) |
| `formats` | object | - | Output formats |
| `templates` | object | - | Custom templates |

### Feed Formats

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `html` | bool | `true` | Generate HTML pages |
| `simple_html` | bool | `true` | Generate simple list HTML pages |
| `rss` | bool | `true` | Generate RSS feed |
| `atom` | bool | `false` | Generate Atom feed |
| `json` | bool | `false` | Generate JSON feed |
| `markdown` | bool | `false` | Generate Markdown file |
| `text` | bool | `false` | Generate text file |
| `sitemap` | bool | `false` | Generate sitemap |

### Feed Templates

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `html` | string | `"feed.html"` | HTML template |
| `simple_html` | string | `"simple-feed.html"` | Simple list HTML template |
| `card` | string | `"card.html"` | Post card template |
| `rss` | string | `"feed.xml"` | RSS template |
| `atom` | string | `"atom.xml"` | Atom template |
| `json` | string | `"feed.json"` | JSON template |

### Feed Defaults

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `items_per_page` | int | `10` | Default posts per page |
| `orphan_threshold` | int | `3` | Default orphan threshold |
| `formats` | object | - | Default formats |
| `templates` | object | - | Default templates |

### Syndication Settings

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `max_items` | int | `20` | Max items in RSS/Atom |
| `include_content` | bool | `false` | Include full content |

### Auto Tags Config

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Enable auto tag feeds |
| `slug_prefix` | string | `"tags"` | URL prefix |
| `formats` | object | - | Output formats |

### Auto Archives Config

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Enable auto archives |
| `slug_prefix` | string | `"archive"` | URL prefix |
| `yearly_feeds` | bool | `true` | Generate yearly feeds |
| `monthly_feeds` | bool | `false` | Generate monthly feeds |
| `formats` | object | - | Output formats |

## Complete Example

```toml
[markata-go]
title = "My Blog"
url = "https://myblog.com"
author = "Jane Doe"

# =============================================================================
# FEED DEFAULTS
# =============================================================================
[markata-go.feeds.defaults]
items_per_page = 10
orphan_threshold = 3

[markata-go.feeds.defaults.formats]
html = true
rss = true
atom = true
json = false

[markata-go.feeds.defaults.templates]
html = "feed.html"
card = "partials/card.html"

[markata-go.feeds.syndication]
max_items = 20
include_content = false

# =============================================================================
# INDIVIDUAL FEEDS
# =============================================================================

# Home page
[[markata-go.feeds]]
slug = ""
title = "Recent Posts"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 5
formats = { rss = false, atom = false }

# Main blog feed
[[markata-go.feeds]]
slug = "blog"
title = "Blog"
description = "All blog posts"
filter = "published == True"
sort = "date"
reverse = true

# Tutorials
[[markata-go.feeds]]
slug = "tutorials"
title = "Tutorials"
filter = "published == True and 'tutorial' in tags"
sort = "date"
reverse = true

# Archive (no pagination)
[[markata-go.feeds]]
slug = "archive"
title = "Archive"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 0
formats = { rss = false, atom = false }

# JSON API
[[markata-go.feeds]]
slug = "api/posts"
title = "Posts API"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 0
formats = { html = false, rss = false, atom = false, json = true }

# =============================================================================
# AUTO-GENERATED FEEDS
# =============================================================================

[markata-go.feeds.auto_tags]
enabled = true
slug_prefix = "tags"

[markata-go.feeds.auto_tags.formats]
html = true
rss = true
```

**Generated output:**

```
public/
  index.html                     # Home (5 posts)
  blog/
    index.html                   # Blog page 1
    page/2/index.html            # Blog page 2
    simple/
      index.html                 # Simple list page 1
      page/2/index.html          # Simple list page 2
    rss.xml
    atom.xml
  tutorials/
    index.html
    simple/
      index.html
    rss.xml
    atom.xml
  archive/
    index.html                   # All posts, no pagination
    simple/
      index.html
  api/
    posts/
      feed.json                  # JSON API
  tags/
    python/
      index.html
      rss.xml
    go/
      index.html
      rss.xml
```

---

## Next Steps

Now that you understand feeds, here are recommended next steps:

**Customize feed appearance:**
- [Templates Guide](/docs/guides/templates/) - Customize how feeds and cards are rendered

**Add syndication and discovery:**
- [Syndication Guide](/docs/guides/syndication/) - Share your feeds on Mastodon, Twitter, and other platforms

**Deploy your site:**
- [Deployment Guide](/docs/guides/deployment/) - Deploy to production with CI/CD

---

## Calendar View

Open Calendar from a built-in feed's view controls to browse publishing activity
one year at a time. The newest year with posts appears first. Use the older/newer
buttons or year selector to jump between populated years; empty years are skipped.
Browser Back and Forward restore your selection.

```text
/blog/?view=calendar&year=2024
```

Months use the available calendar width, with larger day targets and fewer
columns in narrow containers. The selected year shows its post count and number
of publishing days. Day previews still open on hover, focus, or activation.

Calendar grids are built only when Calendar is opened, and only twelve months
are built at once. On regular feeds, compact full-history metadata comes from a
JSON sidecar fetched only when Calendar is opened. The initial feed HTML keeps
the current page's posts without duplicating every calendar entry. With
JavaScript disabled, the ordinary paginated feed list remains available.
Dedicated calendar-first feeds keep their visible dated list as a fallback.
Ordinary feed pagination does not truncate the years available in Calendar.

For a calendar-first feed, set `templates.html = "calendar-feed.html"` on the
feed. The same annual navigation applies, using `?year=2024` on its URL.

Calendar panels and the year selector follow the active theme surface and text colors. On wide calendar pages, twelve months form four columns and three rows; smaller containers keep fewer, readable columns.

## Publishing Rhythm Charts

The built-in `/feeds/` directory can show a compact monthly publishing-rhythm chart
for each feed. Hover over a chart to read a month's post count. Keyboard users can
focus a chart and use the left and right arrow keys to move between months; Home and
End jump to the first and latest month. The chart's line and overall appearance stay
the same while its month details are available on demand.

## See Also

- [Configuration Guide](/docs/guides/configuration/) - Full configuration reference
- [Templates Guide](/docs/guides/templates/) - Template system documentation
- [Frontmatter Guide](/docs/guides/frontmatter/) - Post metadata for filtering
- [Quick Reference](/docs/guides/quick-reference/) - Filter expression cheat sheet

### Editorial loading

Editorial archives render the visible cards immediately and let the browser
defer layout for offscreen cards and later weeks until they approach the
viewport. All posts are still in the HTML, including when JavaScript is
disabled. Browser find and printing keep access to the complete page. Browsers
without content visibility support render the complete page normally.

Keyboard navigation setup does not download the next archive page. Link hover,
focus, and explicit feed cycling can still prepare the next navigation.

Calendar feeds reserve their full list width before optional calendar styles
arrive. Navbar search also reserves its final width before search hydration,
so opening an archive does not move the content as those styles finish loading.

Editorial weeks below the first viewport reserve space according to the number
of items they contain. This keeps unrendered weeks from crowding following
sections. The browser remembers their actual height after revealing them.
