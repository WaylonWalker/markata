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
│   HTML   │   RSS    │   Atom   │   JSON   │ Markdown │    Text     │ Sitemap  │
│  index   │   feed   │   feed   │   API    │   list   │   list      │   XML    │
└──────────┴──────────┴──────────┴──────────┴──────────┴─────────────┴──────────┘
```

## Why Feeds Matter

1. **One definition, many outputs** - Define a collection once, get HTML pages, RSS, Atom, JSON API, sitemap, etc.
2. **Composable** - Feeds can reference other feeds or share filters
3. **Flexible** - Every "index page" is a feed: home, archives, tags, categories, search indexes
4. **Familiar** - Mirrors how content platforms work (RSS readers, APIs, syndication)

---

## Feed Configuration

### Basic Feed

```toml
[[name.feeds]]
slug = "blog"
title = "Blog"
filter = "published == True"
sort = "date"
reverse = true
```

### Full Configuration

```toml
[[name.feeds]]
# Identity
slug = "blog"                      # URL path: /blog/
title = "All Posts"                # Display title
description = "Latest blog posts"  # Meta description

# Content Selection
filter = "published == True and date <= today"
sort = "date"
reverse = true                     # Newest first

# Pagination
items_per_page = 10                # 0 = no pagination (all on one page)
orphan_threshold = 3               # If last page has ≤3 items, merge with previous

# Output Formats (all optional, defaults shown)
[name.feeds.formats]
html = true                        # /blog/index.html, /blog/page/2/index.html
rss = true                         # /blog/rss.xml
atom = true                        # /blog/atom.xml
json = true                        # /blog/feed.json
markdown = false                   # /blog/index.md
text = false                       # /blog/index.txt
sitemap = false                    # /blog/sitemap.xml

# Templates (can override per-format)
[name.feeds.templates]
html = "feed.html"                 # Template for HTML pages
card = "partials/card.html"        # Template for post cards in list
rss = "rss.xml"                    # Template for RSS (usually default)
atom = "atom.xml"                  # Template for Atom
json = "feed.json"                 # Template for JSON
sitemap = "sitemap.xml"            # Template for sitemap
```

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
| `posts` | List[Post] | All matching posts (pre-pagination) |
| `page_posts` | List[Post] | Posts for current page |
| `pagination` | Pagination | Pagination info |
| `formats` | FeedFormats | Enabled output formats (html, rss, atom, json, markdown, text, sitemap) |

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

---

## Pagination Types

Feeds support multiple pagination strategies for different use cases and user experiences.

### Configuration

```toml
[[name.feeds]]
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

[[name.feeds]]
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

## Manual Pagination

Traditional pagination with full page reloads. Best for SEO and accessibility.

### Configuration

```toml
[[name.feeds]]
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
[[name.feeds]]
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
    <lastBuildDate>{{ now | rss_date }}</lastBuildDate>
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
  <updated>{{ now | atom_date }}</updated>
  
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

**Output path:** `/blog/index.md`

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

**Output path:** `/blog/index.txt`

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
[name.feeds.sitemap]
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
[[name.feeds]]
slug = "sitemap"
title = "Sitemap"
filter = "published == True"      # All published posts
sort = "date"
reverse = true
formats = { sitemap = true }      # Only sitemap output
```

This generates `/sitemap.xml` at the root when `slug = ""`:

```toml
[[name.feeds]]
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
[[name.feeds]]
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
[[name.feeds]]
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
[name.feeds.auto_tags]
enabled = true
slug_prefix = "tags"               # /tags/python/, /tags/rust/
title_template = "Posts tagged '{{ tag }}'"
filter_template = "published == True and '{{ tag }}' in tags"
formats = { html = true, rss = true }
```

This auto-generates feeds for each unique tag found in posts.

### Category Pages

```toml
[name.feeds.auto_categories]
enabled = true
slug_prefix = "category"
source_field = "category"          # Frontmatter field to use
title_template = "{{ category | title }}"
filter_template = "published == True and category == '{{ category }}'"
```

### Year/Month Archives

```toml
[name.feeds.auto_date]
enabled = true
slug_template = "{{ date.year }}/{{ date.month }}"
title_template = "Posts from {{ date | date('%B %Y') }}"
filter_template = "published == True and date.year == {{ year }} and date.month == {{ month }}"
formats = { html = true }
```

### API Endpoint

JSON-only feed for JavaScript consumption:

```toml
[[name.feeds]]
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
[[name.feeds]]
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

---

## Configuration Inheritance

Feed configuration follows a **defaults → override** pattern. Global defaults apply to all feeds, and individual feeds can override any setting.

### Inheritance Rules

```
┌─────────────────────────────────────────────────────────────────────┐
│                    CONFIGURATION RESOLUTION                          │
├─────────────────────────────────────────────────────────────────────┤
│  1. Start with built-in defaults                                     │
│  2. Apply [name.feeds.defaults.*] settings                          │
│  3. Apply [[name.feeds]] individual feed settings                   │
│  4. Individual feed values WIN over defaults                         │
└─────────────────────────────────────────────────────────────────────┘
```

### Configurable at Both Levels

| Setting | Global Key | Feed Override |
|---------|------------|---------------|
| Items per page | `defaults.items_per_page` | `items_per_page` |
| Orphan threshold | `defaults.orphan_threshold` | `orphan_threshold` |
| HTML format | `defaults.formats.html` | `formats.html` |
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
[name.feeds.defaults]
items_per_page = 10
orphan_threshold = 3

[name.feeds.defaults.formats]
html = true
rss = true
atom = false                       # Atom OFF by default
json = false                       # JSON OFF by default
markdown = false
text = false
sitemap = false                    # Sitemap OFF by default

[name.feeds.defaults.templates]
html = "feed.html"
card = "partials/card.html"
rss = "rss.xml"

[name.feeds.syndication]
max_items = 20
include_content = false

# =============================================================================
# INDIVIDUAL FEEDS (inherit from defaults, override as needed)
# =============================================================================

# Home page - uses most defaults, but fewer items and no RSS
[[name.feeds]]
slug = ""
title = "Home"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 5                 # Override: fewer items on home
formats = { rss = false }          # Override: no RSS for home page

# Blog - uses all defaults
[[name.feeds]]
slug = "blog"
title = "Blog"
filter = "published == True"
sort = "date"
reverse = true
# items_per_page: inherits 10 from defaults
# formats: inherits html=true, rss=true from defaults

# API endpoint - completely different format set
[[name.feeds]]
slug = "api/posts"
title = "Posts API"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 0                 # Override: no pagination
formats = { html = false, rss = false, json = true }  # Override: JSON only

# Archive - different pagination
[[name.feeds]]
slug = "archive"
title = "All Posts"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 0                 # Override: all posts on one page
formats = { rss = false }          # Override: no RSS for archive

# Tutorials - wants Atom too
[[name.feeds]]
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
# Defaults: html=true, rss=true, atom=false, json=false
formats = { atom = true }
# Result: html=true, rss=true, atom=true, json=false
```

**2. Disable specific formats:**
```toml
# Defaults: html=true, rss=true
formats = { rss = false }
# Result: html=true, rss=false
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
[name.feeds.defaults.templates]
html = "feed.html"
card = "partials/card.html"

# Individual feed override
[[name.feeds]]
slug = "featured"
templates = { html = "featured-feed.html" }
# Result: html="featured-feed.html", card="partials/card.html" (inherited)
```

---

## Configuration Reference

### Global Feed Defaults

```toml
[name.feeds.defaults]
# Pagination
items_per_page = 10
orphan_threshold = 3

# Formats enabled by default
[name.feeds.defaults.formats]
html = true
rss = true
atom = false
json = false
markdown = false
text = false
sitemap = false

# Templates used by default
[name.feeds.defaults.templates]
html = "feed.html"
card = "partials/card.html"
rss = "rss.xml"
atom = "atom.xml"
json = "feed.json"
markdown = "feed.md"
text = "feed.txt"
sitemap = "sitemap.xml"

# Syndication settings (RSS/Atom/JSON)
[name.feeds.syndication]
max_items = 20                     # Max items in RSS/Atom feeds
include_content = false            # Include full content or just summary
```

### Built-in Defaults

If no `[name.feeds.defaults]` is specified, these built-in values apply:

| Setting | Built-in Default |
|---------|------------------|
| `items_per_page` | 10 |
| `orphan_threshold` | 3 |
| `formats.html` | true |
| `formats.rss` | true |
| `formats.atom` | false |
| `formats.json` | false |
| `formats.markdown` | false |
| `formats.text` | false |
| `formats.sitemap` | false |
| `syndication.max_items` | 20 |
| `syndication.include_content` | false |

### Feed Item Limits

| Format | Default Limit | Configurable |
|--------|---------------|--------------|
| HTML | Paginated | `items_per_page` |
| RSS | 20 | `syndication.max_items` or `max_items` |
| Atom | 20 | `syndication.max_items` or `max_items` |
| JSON | All | `items_per_page` |
| Markdown | All | `items_per_page` |
| Text | All | `items_per_page` |
| Sitemap | All | N/A (always includes all posts) |

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
[name.feeds.defaults]
items_per_page = 10
orphan_threshold = 3

[name.feeds.defaults.formats]
html = true
rss = true
atom = true
json = true

[name.feeds.defaults.templates]
html = "feed.html"
card = "partials/card.html"

[name.feeds.syndication]
max_items = 20
include_content = false

# =============================================================================
# INDIVIDUAL FEEDS
# =============================================================================

# Home page - fewer items, no syndication feeds
[[name.feeds]]
slug = ""
title = "Recent Posts"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 5                 # Override default
formats = { rss = false, atom = false, json = false }  # Override: HTML only

# Main blog feed - uses all defaults
[[name.feeds]]
slug = "blog"
title = "Blog"
description = "All blog posts"
filter = "published == True and templateKey == 'post'"
sort = "date"
reverse = true
# Inherits: items_per_page=10, all formats enabled

# Tutorials section - uses defaults
[[name.feeds]]
slug = "tutorials"
title = "Tutorials"
filter = "published == True and 'tutorial' in tags"
sort = "date"
reverse = true

# Archive - all posts, no pagination, HTML only
[[name.feeds]]
slug = "archive"
title = "Archive"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 0                 # Override: no pagination
formats = { rss = false, atom = false, json = false }  # Override: HTML only

# API endpoint - JSON only, no pagination
[[name.feeds]]
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
[name.feeds.auto_tags]
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
