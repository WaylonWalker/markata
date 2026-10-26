# Template System Specification

This document specifies the template system for rendering content.

## Navigation previews

The default header and footer navigation components resolve local URLs against generated feeds, published posts, and the enabled Reader, Blogroll, and Random routes. A matching link may show a compact preview on hover or keyboard focus. Feed previews show their description when present, the count of public posts, total words and estimated reading time when available, and a monthly publication rhythm for dated posts. Post previews show their description, tags, word count, and reading time when available. Reader and Blogroll previews show the number of active followed sources. Random previews show the number of posts eligible for its configured draw. Missing values are omitted. External and unmatched links remain ordinary links.

Preview content is generated at build time and escaped by the template engine. Private or unpublished post metadata must never appear in a preview. The nav link remains directly navigable, and preview content is supplementary. On narrow viewports, previews must not obstruct the scrollable navigation strip. Keyboard focus and reduced-motion preferences must be supported. Pointer previews stay open while the pointer moves from the link into the panel, close after leaving the link and panel, and stack above page content and other navigation previews.

Local Markdown links and wikilinks in a post also resolve against public generated feeds and posts. Matching links show the same feed or post statistics as navigation previews, with the destination title rather than the link label. Absolute links to the configured site or a loopback development host resolve by path; other external links keep their existing hover behavior. A post's tag links show a feed preview when a public tag feed exists. Unmatched links and tags remain ordinary links. These article hover cards support keyboard focus and do not expose private or unpublished metadata.

## Overview

The `render_feed` helper's `summary` variant MUST render linked titles, bounded
plain-text summaries, and dates without embedding article media or code blocks.
The default `card` variant MUST retain its existing routing behavior.

The `media_srcset` filter MUST emit width candidates only for trusted remote
resize services. Relative and untrusted URLs MUST return an empty srcset. It
MUST preserve unrelated query parameters and remove fixed-height parameters.
Default photo and article cards SHOULD use these candidates with native lazy
loading and `sizes="auto, ..."` so small cards do not request full-size media.

Templates wrap rendered markdown content in HTML layouts. The system supports:

- Template inheritance (base templates)
- Includes (partials)
- Variables and expressions
- Control flow (if/for)
- Filters and functions
- Custom template per post
- Format-specific templates for `.txt`, `.ansi`, `.md`, and OG outputs
- Presentation templates such as `slides.html` for reveal.js decks

Each template execution MUST receive its own top-level configuration map. The
cached configuration projection is shared across renders and MUST remain
immutable after creation; render-specific values such as the active theme must
be applied to a per-execution copy. Concurrent feed rendering MUST NOT mutate
maps read by another template execution.

---

## Cached HTML Restoration

Warm builds MAY restore rendered article HTML and complete page HTML from the
build cache in either the DAG or legacy pipeline. Restoration MUST be
byte-identical to the cached content, including empty content, NUL bytes and
invalid UTF-8; it MUST NOT normalize, decode, truncate or otherwise transform
the HTML. Byte-preserving disk restoration MUST NOT itself change cache keys,
versions, paths or freshness checks (including the article content hash).
Rendering behavior revisions MAY invalidate derived caches through the build
identity so obsolete article, page and feed output is recomputed consistently.

The HTML disk reader MUST avoid an additional body-sized temporary byte buffer
followed by a whole-body byte-to-string copy. Returned strings MUST own immutable
storage and remain unchanged after subsequent reads or reuse of transfer
buffers. Reusable transfer buffers MUST have bounded individual capacity and
MUST NOT back returned strings.

File metadata is only an allocation hint: preallocation MUST be bounded and
safe for the platform's integer size, including suspiciously large, unknown or
nonregular file sizes. This MUST NOT impose a content-size limit. Reads MUST
continue to EOF even if the file grows or shrinks relative to its reported size.
Open, stat, read and close failures MUST be reported by the internal reader
without returning partial HTML. Missing or unreadable disk cache entries MUST
remain cache misses so the caller can re-render, never publish a partial page.
Existing in-memory hits MUST remain usable without accessing the disk.
Concurrent full-page reads sharing a cache path MUST return the same published
backing string rather than retaining one body allocation per worker.

Template cache diagnostics MUST observe the existing first-failing gates,
without changing selection, freshness, rendering, restoration, or best-effort
storage behavior. The bounded `template_cache` aggregate and reconciliation
rules are specified in [Content Diagnostics](CONTENT_DIAGNOSTICS.md).
An empty full-page getter result remains a render fallback, including an empty
cached page; skipped posts are not classified, while empty article bodies are.

Full-page restoration MUST remain eager: every usable cache hit MUST populate
public `Post.HTML` before any cache-miss template executes, including templates
that inspect peer posts through Core. Restoration MAY use a bounded worker pool,
limited by manager concurrency and a restoration-specific cap of four without
changing manager concurrency. Concurrent workers MUST join before HTML assignment,
counter aggregation, canonical serve filtering, rendering, or publication.
Unavailable full pages MUST append after classification misses in their original
cacheable input order; only nonempty successful results may replace `Post.HTML`.
Existing eligibility, encryption invalidation, and missing-cache fallback rules
MUST remain unchanged. Complete pages remain available to hooks, publication,
Fontpack, and Tailwind; this optimization does not remove their live-memory floor.

Template phase logs MUST measure classification, full-page restoration, and
fresh rendering separately, including invocations requiring no fresh rendering.
Restore time MUST NOT be attributed to the rendering phase. Aggregate diagnostic
counts and their JSON schema remain unchanged.

### Canonical semantic hash handoff

For every freshly parsed post, Load MUST capture the previous feed, tag, and
garden hashes before updating cache metadata. The build-local handoff MUST
contain only source paths and these three hashes, be concurrency-safe, capture
each path only once, and bind to the actual build-cache instance. An explicitly
captured empty triplet represents a new post, not an absent handoff.

An existing entry with a recorded input identity but incomplete semantic hashes
does not provide a canonical baseline. This can occur when the cold navigation
reset drops derived metadata before pages are rendered and published again.
Load MUST capture that unavailability once per path, rather than interpreting
missing hashes as a new post or recapturing its later intermediate hashes.
InlineTitles MUST retain the original Update-return comparison for such paths.
Fresh entries without a prior input identity still have a valid empty/partial
triplet. The handoff may encode unavailability without retaining another hash,
post snapshot, or private value.

Every Load invocation MUST reset the handoff, including early-return paths.
InlineTitles MUST persist the current canonical hashes after auto-title
inference, comparing against the captured prior-build hashes for its own
feed/slug invalidation when available. Otherwise it MUST retain the existing
Update-return comparisons. The handoff MUST be consumed and discarded after
InlineTitles completes, including errors; repeated execution and other managers
or cache instances MUST NOT reuse it.

Load-only and nonstandard pipelines retain Load's existing updates and
conservative feed/tag/garden dirty signals. This handoff MUST NOT clear changes
from other producers, or permit decrypted parsed-post/article caching.
Encryption MUST independently invalidate pages when their encrypted wrapper is
regenerated (see [Encryption](ENCRYPTION.md)).

Regenerated encrypted wrappers MUST invalidate the previous full-page cache
reference before fresh page storage. If best-effort full-page storage fails,
the reference MUST remain unavailable across cache save/reload, so the next
unchanged wrapper hit re-renders rather than restores stale ciphertext.

---

## Template Location

Templates are loaded from (in order):

1. **User templates:** `templates/` directory in project root
2. **Cache templates:** `.tool.cache/templates/` (generated)
3. **Package templates:** Built-in default templates

```
my-site/
├── templates/           # User templates (highest priority)
│   ├── base.html
│   ├── post.html
│   └── partials/
│       └── nav.html
├── .tool.cache/
│   └── templates/       # Generated templates
└── posts/
```

The rendered-page cache MUST include a content hash for both project template
files and the embedded default theme templates. A change to either source must
invalidate cached page HTML, including when the builder binary changes but the
site content and project template directory do not.

---

## Base Template

The foundation template that others extend:

```html
<!-- templates/base.html -->
<!DOCTYPE html>
<html lang="{{ config.lang | default('en') }}">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{% block title %}{{ config.title }}{% endblock %}</title>
    <meta name="description" content="{% block description %}{{ config.description }}{% endblock %}">
    {% block head %}{% endblock %}
</head>
<body>
    {% block header %}
    {% include "partials/header.html" %}
    {% endblock %}

    <main>
        {% block content %}{% endblock %}
    </main>

    {% block footer %}
    {% include "partials/footer.html" %}
    {% endblock %}

    {% block scripts %}{% endblock %}
</body>
</html>
```

---

## Post Template

Template for individual posts:

```html
<!-- templates/post.html -->
{% extends "base.html" %}

{% block title %}{{ post.title }} | {{ config.title }}{% endblock %}

{% block description %}{{ post.description | default(config.description) }}{% endblock %}

{% block head %}
<meta property="og:title" content="{{ post.title }}">
<meta property="og:type" content="article">
<meta property="og:url" content="{{ config.url }}{{ post.href }}">
{% if post.date %}
<meta property="article:published_time" content="{{ post.date.isoformat() }}">
{% endif %}
{% endblock %}

{% block content %}
<article>
    <header>
        <h1>{{ post.title }}</h1>
        {% if post.date %}
        <time datetime="{{ post.date.isoformat() }}">{{ post.date.strftime('%B %d, %Y') }}</time>
        {% endif %}
        {% if post.tags %}
        <ul class="tags">
            {% for tag in post.tags %}
            <li><a href="/tags/{{ tag | slugify }}/">{{ tag }}</a></li>
            {% endfor %}
        </ul>
        {% endif %}
    </header>

    <div class="content">
        {{ body | safe }}
    </div>

    {% if post.prev or post.next %}
    <nav class="post-nav">
        {% if post.prev %}
        <a href="{{ post.prev.href }}" class="prev">← {{ post.prev.title }}</a>
        {% endif %}
        {% if post.next %}
        <a href="{{ post.next.href }}" class="next">{{ post.next.title }} →</a>
        {% endif %}
    </nav>
    {% endif %}
</article>
{% endblock %}
```

---

## Slides Template

Markata-go MAY render a post with `template: "slides.html"` for presentation-style decks.

The default `slides.html` behavior is:

1. Use [reveal.js](https://revealjs.com/) for keyboard navigation and slide rendering.
2. Convert rendered article HTML into nested `<section>` elements.
3. Start a new horizontal slide when the rendered body encounters:
   - an `<h2>` heading
   - an `<hr>` element, which matches the common `---` Markdown slide separator used by reveal.js and Marp
4. Start a new vertical slide within the current horizontal stack when the rendered body encounters an `<h3>` heading.
5. Preserve any leading content before the first `<h2>` as the first slide in the deck.
6. Load reveal.js from the shared CDN asset registry when `[markata-go.assets].mode` is self-hosted or auto and the assets are available; otherwise fall back to CDN URLs.
7. Size the reveal.js deck responsively so slides remain readable on narrow mobile viewports and continue to relayout on viewport changes such as rotation.

This allows authors to write decks in plain Markdown using common slide-authoring conventions:

```markdown
## Opening

Welcome.

### Details

Vertical slide under Opening.

---

## Next Chapter
```

The template system MUST expose a `slides_reveal` filter that accepts rendered HTML and returns safe reveal.js section markup.

---

## Feed Template

Template for index/listing pages:

```html
<!-- templates/feed.html -->
{% extends "base.html" %}

{% block title %}{{ feed.title }} | {{ config.title }}{% endblock %}

{% block content %}
<section class="feed">
    <h1>{{ feed.title }}</h1>

    <ul class="post-list">
        {% for post in feed.posts %}
        <li>
            {% include "partials/card.html" %}
        </li>
        {% endfor %}
    </ul>

    {% if feed.pagination.total_pages > 1 %}
    <nav class="pagination">
        {% if feed.pagination.has_prev %}
        <a href="{{ feed.pagination.prev_url }}">← Previous</a>
        {% endif %}

        <span>Page {{ feed.pagination.current_page }} of {{ feed.pagination.total_pages }}</span>

        {% if feed.pagination.has_next %}
        <a href="{{ feed.pagination.next_url }}">Next →</a>
        {% endif %}
    </nav>
    {% endif %}
</section>
{% endblock %}
```

---

## Template Context

### Post Templates

| Variable | Type | Description |
|----------|------|-------------|
| `post` | Post | The post being rendered |
| `body` | str | Rendered HTML content |
| `config` | Config | Site configuration |
| `core` | Core | Core instance |

**Link graph fields (from link_collector):**

- `post.hrefs` - list of href strings found in the post
- `post.inlinks` - list of link maps pointing to this post
- `post.outlinks` - list of link maps pointing from this post

Each link map includes:

- `source_url`, `source_text`
- `target_url`, `target_domain`, `target_text`
- `is_internal`, `is_self`, `raw_target`
- `source_post`, `target_post` (maps with `href`, `slug`, `title`)

### OG Templates

OpenGraph card pages use dedicated templates. The template name is derived from the
post template when possible:

- `post.html` -> `post-og.html`
- If no post template exists, fall back to `og-card.html`

The default theme ships both `post-og.html` and `og-card.html` so OG cards match the
site palette and background styling.

When an OG card template renders post media, it MUST use this frontmatter fallback order:

1. `image`
2. `cover_image`
3. `og_image`
4. `video`

If the resolved media is a video file, the OG template should render a video tag instead of an image tag. If none of these fields are present, the OG card renders a text-only layout.

### Terminal Templates

Plain terminal and ANSI terminal outputs use dedicated format templates resolved with the same per-format rules as HTML and OG output:

- `post.html` -> `post.txt`
- `post.html` -> `post.ansi`
- if no format-specific template exists, fall back to `default.txt` or `default.ansi`

Terminal templates receive:

- `body`: the fully rendered terminal page for the target format
- `post`: post metadata and rendered fields
- `config`: resolved configuration including theme and post format settings

Terminal templates SHOULD render `{{ body|safe }}` when they want the default terminal page body. Custom templates MAY mix `body` with additional metadata or custom separators.

`default.txt` MUST emit plain text only. `default.ansi` MAY include ANSI escape sequences and SHOULD preserve semantic emphasis, links, admonitions, blockquotes, tables, and syntax-highlighted code blocks in a terminal-safe form.

### Feed Templates

| Variable | Type | Description |
|----------|------|-------------|
| `feed` | Feed | Feed configuration and posts |
| `feed.posts` | List[Post] | Posts in this feed page |
| `feed.pagination` | Pagination | Pagination info |
| `config` | Config | Site configuration |
| `core` | Core | Core instance |

### Global Variables

Always available:

| Variable | Type | Description |
|----------|------|-------------|
| `today` | date | Current date |
| `now` | datetime | Current datetime |

---

## Feed Embed Helpers

Two helper callables make it easy to surface feed posts inside templates and `jinja_md`-enabled markdown without embedding raw HTML.

### `feed_posts(feed_slug, limit=...)`

- Returns a list of post maps (same shape as `templates.PostsToMaps`) for the requested feed slug.
- Accepts either a numeric positional argument or an options map containing `"limit"` to cap the results; zero or missing limit returns the full feed output.
- Respects the feed configuration (`filter`, `sort`, `offset`, `limit`, `include_private`) but never fails the build when the slug is missing—an empty list is returned instead.

### `render_feed(feed_slug, options={})`

- Renders a safe HTML snippet for the feed's latest posts using the default partial `partials/feed_preview.html`. The `card` layout now embeds `partials/cards/card-router.html`, giving each post type (article, note, photo, video, link, quote, guide, etc.) its own card markup. Photo/shot cards wrap the media in a `<figure>` and render a `<figcaption>` caption derived from the description/title so the visual output resembles a figure with caption. The `list` layout keeps the compact list styles.
- Supported options: `limit` (numeric override), `variant` (`"card"` or `"list"`, default `"card"`), and `template` (custom template path). If `template` is provided, its output must accept `posts`, `feed`, and `variant` in the context.
- The helper looks up the shared `templates.engine` from the manager cache and gracefully falls back to simple markup when the engine or template is unavailable so that markdown renderers can use it even before the full render stage runs.
---

## Template Syntax

### Variables

```jinja2
{{ post.title }}
{{ config.url }}
{{ post.date.year }}
```

### Attribute Access

```jinja2
{{ post.title }}             {# Attribute #}
{{ post["title"] }}          {# Item access #}
{{ post.get("title", "") }}  {# With default #}
```

### Filters

```jinja2
{{ post.title | upper }}
{{ post.title | lower }}
{{ post.title | title }}
{{ post.description | truncate(160) }}
{{ post.date | date('%Y-%m-%d') }}
{{ tags | join(', ') }}
{{ body | safe }}            {# Don't escape HTML #}
{{ value | default('N/A') }} {# Default value #}
```

### Control Flow

```jinja2
{% if post.published %}
    <span class="status">Published</span>
{% elif post.draft %}
    <span class="status">Draft</span>
{% else %}
    <span class="status">Private</span>
{% endif %}

{% for tag in post.tags %}
    <a href="/tags/{{ tag | slugify }}/">{{ tag }}</a>
    {% if not loop.last %}, {% endif %}
{% endfor %}

{% for post in core.filter("published == True") %}
    <li>{{ post.title }}</li>
{% empty %}
    <li>No posts found</li>
{% endfor %}
```

### Loop Variables

Inside `{% for %}` loops:

| Variable | Description |
|----------|-------------|
| `loop.index` | Current iteration (1-indexed) |
| `loop.index0` | Current iteration (0-indexed) |
| `loop.first` | True if first iteration |
| `loop.last` | True if last iteration |
| `loop.length` | Total number of items |
| `loop.revindex` | Iterations remaining (1-indexed) |

### Template Inheritance

```jinja2
{# child.html #}
{% extends "base.html" %}

{% block content %}
    <h1>Child content</h1>
    {{ super() }}  {# Include parent block content #}
{% endblock %}
```

### Includes

```jinja2
{% include "partials/header.html" %}

{% include "partials/card.html" with context %}

{% include "partials/sidebar.html" ignore missing %}
```

### Macros

```jinja2
{% macro post_card(post) %}
<article class="card">
    <h2><a href="{{ post.href }}">{{ post.title }}</a></h2>
    <p>{{ post.description }}</p>
</article>
{% endmacro %}

{# Usage #}
{{ post_card(post) }}

{% for p in posts %}
    {{ post_card(p) }}
{% endfor %}
```

### Comments

```jinja2
{# This is a comment and won't appear in output #}

{#
Multi-line
comment
#}
```

---

## Built-in Filters

### String Filters

| Filter | Example | Output |
|--------|---------|--------|
| `upper` | `{{ "hello" \| upper }}` | `HELLO` |
| `lower` | `{{ "HELLO" \| lower }}` | `hello` |
| `title` | `{{ "hello world" \| title }}` | `Hello World` |
| `capitalize` | `{{ "hello" \| capitalize }}` | `Hello` |
| `trim` | `{{ "  hi  " \| trim }}` | `hi` |
| `truncate(n)` | `{{ text \| truncate(100) }}` | First 100 chars... |
| `striptags` | `{{ html \| striptags }}` | Remove HTML tags |
| `plaintext` | `{{ html \| plaintext }}` | Convert HTML to clean plain text (see below) |
| `summary(n)` | `{{ html \| summary:300 }}` | Single-line escaped card excerpt (see below) |
| `escape` | `{{ html \| escape }}` | HTML escape (default) |
| `safe` | `{{ html \| safe }}` | Don't escape |
| `slugify` | `{{ "Hello World" \| slugify }}` | `hello-world` |

### List Filters

| Filter | Example | Output |
|--------|---------|--------|
| `length` | `{{ list \| length }}` | Count |
| `first` | `{{ list \| first }}` | First item |
| `last` | `{{ list \| last }}` | Last item |
| `join(sep)` | `{{ list \| join(", ") }}` | Joined string |
| `sort` | `{{ list \| sort }}` | Sorted list |
| `reverse` | `{{ list \| reverse }}` | Reversed list |
| `unique` | `{{ list \| unique }}` | Deduplicated |
| `map(attr)` | `{{ posts \| map(attribute='title') }}` | Extract field |
| `selectattr` | `{{ posts \| selectattr('published') }}` | Filter by attr |
| `rejectattr` | `{{ posts \| rejectattr('draft') }}` | Exclude by attr |
| `batch(n)` | `{{ list \| batch(3) }}` | Group into chunks |

### Number Filters

| Filter | Example | Output |
|--------|---------|--------|
| `round` | `{{ 3.7 \| round }}` | `4` |
| `int` | `{{ "42" \| int }}` | `42` |
| `float` | `{{ "3.14" \| float }}` | `3.14` |
| `abs` | `{{ -5 \| abs }}` | `5` |
| `filesizeformat` | `{{ 1024 \| filesizeformat }}` | `1.0 KB` |

### Date Filters

| Filter | Example | Output |
|--------|---------|--------|
| `date(fmt)` | `{{ post.date \| date('%Y-%m-%d') }}` | `2024-01-15` |
| `strftime(fmt)` | `{{ post.date.strftime('%B %d') }}` | `January 15` |

### Utility Filters

| Filter | Example | Output |
|--------|---------|--------|
| `default(val)` | `{{ x \| default('N/A') }}` | Value or default |
| `tojson` | `{{ obj \| tojson }}` | JSON string |
| `pprint` | `{{ obj \| pprint }}` | Pretty print |

---

## Custom Filters

Implementations can add custom filters:

**Example (pseudocode):**

```
function reading_time(content, wpm=200):
    words = length(content.split())
    minutes = max(1, words / wpm)
    return "{minutes} min read"

template_engine.add_filter("reading_time", reading_time)
```

Usage:
```jinja2
{{ post.content | reading_time }}
{{ post.content | reading_time(250) }}
```

### The `pin_preview` Filter

The built-in `pin_preview` filter accepts a post map and returns `image` and
`commentary` strings derived from that post's fields and rendered HTML. It MUST
perform no network requests and MUST NOT mutate post data. Pins templates MUST
escape the returned strings normally. Selection and note extraction semantics
are specified in [Pins feeds](./FEEDS.md#pins-feed-slugpins).

### The `summary` Filter

The `summary` filter produces a short, single-line excerpt from rendered HTML
for use in cards and other listings. Unlike `plaintext`, whose output is
Markdown-like and multi-line, `summary` output is safe to inject into a
Markdown document (for example through `render_feed`) because it never
contains blank lines or raw HTML that goldmark could re-parse.

**Behavior:**

1. Removes `<pre>`, `<code>`, `<script>`, `<style>`, `<figure>`, `<table>` and
   `<svg>` elements together with their content
2. Strips all remaining tags and decodes entities
3. Collapses whitespace to single spaces and repairs `" ."` gaps left by
   removed inline elements
4. Truncates to `n` characters (default 300) on a rune boundary and appends `…`
5. Returns an HTML-escaped value marked safe

### The `plaintext` Filter

The `plaintext` filter converts HTML content to clean plain text. It is designed for `.txt` template output where HTML tags and entities must not appear.

**Behavior:**

1. **Entity decoding** - HTML entities (`&amp;`, `&lt;`, `&gt;`, `&quot;`, `&#39;`, numeric entities) are decoded to their character equivalents
2. **Tag stripping** - All HTML tags are removed
3. **Block structure** - Block-level elements (`<p>`, `<div>`, `<h1>`-`<h6>`, `<li>`, `<br>`, `<hr>`) produce appropriate line breaks
4. **Footnote-style links** - Anchor tags (`<a href="...">`) are converted to footnote-style references following the Lynx/Pandoc convention
5. **Auto-escaping bypass** - Output is marked safe to prevent pongo2 from re-escaping it

**Link footnote format:**

Links in the HTML are replaced with numbered references in the text body, and a references section is appended at the end:

```
Input HTML:
  <p>Read the <a href="https://go.dev/doc/">Go docs</a> for more info.</p>

Output text:
  Read the Go docs [1] for more info.

  References:
  [1]: https://go.dev/doc/
```

**Footnote rules:**
- Duplicate URLs reuse the same reference number
- When link text is identical to the URL, no footnote marker is added (the URL is already visible)
- Reference numbers are sequential integers starting at 1

**Usage in templates:**

```jinja2
{# In .txt templates - use plaintext for HTML content #}
{{ post.content | plaintext }}

{# For plain text strings that don't contain HTML, use safe instead #}
{{ post.title | safe }}
{{ post.description | safe }}
```

**Difference from `striptags`:** The `striptags` filter only removes HTML tags. The `plaintext` filter additionally decodes HTML entities, preserves block structure with line breaks, converts links to footnotes, and marks the output as safe.

---

## Querying in Templates

Access the full query API:

```jinja2
{# Recent posts #}
{% for post in core.filter("published == True and date <= today")[:5] %}
    <li>{{ post.title }}</li>
{% endfor %}

{# Posts by tag #}
{% for post in core.filter("'python' in tags") %}
    <li>{{ post.title }}</li>
{% endfor %}

{# Get specific post #}
{% set about = core.one("slug == 'about'") %}
<a href="{{ about.href }}">{{ about.title }}</a>

{# Map to get field list #}
{% set all_tags = core.map("tags") | flatten | unique | sort %}
```

---

## Custom Templates Per Post

Posts can specify their template:

```yaml
---
title: Special Post
template: special.html
---
```

Or multiple templates for different contexts:

```yaml
---
title: Feature Post
template:
  default: feature.html
  card: feature-card.html
  feed: feature-feed-item.html
---
```

---

## Partials

Small reusable template fragments:

```html
<!-- templates/partials/card.html -->
<article class="card">
    <a href="{{ post.href }}">
        {% if post.cover_image %}
        <img src="{{ post.cover_image }}" alt="{{ post.title }}">
        {% endif %}
        <h2>{{ post.title }}</h2>
    </a>
    {% if post.description %}
    <p>{{ post.description }}</p>
    {% endif %}
    <footer>
        {% if post.date %}
        <time>{{ post.date | date('%b %d, %Y') }}</time>
        {% endif %}
        {% if post.reading_time %}
        <span>{{ post.reading_time }}</span>
        {% endif %}
    </footer>
</article>
```

## Build-Local Sidebar Projections

Sidebar post projections MAY be shared across page renders within one build.
Their identity MUST include the post slug, effective plain title, original href,
and selected feed slug. The active-post flag MUST be applied to a value copy for
each page, never stored as shared page-specific state. A changed identity MUST
produce a fresh projection, and a new template-render stage MUST reset the cache.
Feed selection, windowing, ordering, URL encoding, and JSON output MUST remain
unchanged.

## Footer License Display

When `config.license` contains a string key the default footer appends the license attribution to the copyright line (next to the copyright symbol). Custom footers should honor the same guard to avoid losing attribution:

```jinja
{% if footer.show_copyright %}
<p class="footer-copyright">
  &copy; {{ now | date:"2006" }}{% if config.author or config.title %} {{ config.author | default:config.title }}.{% endif %}
  {% if config.license and config.license.name %}
  Content licensed under
  {% if config.license.url %}
  <a href="{{ config.license.url }}" target="_blank" rel="noopener noreferrer">{{ config.license.name }}</a>.
  {% else %}
  {{ config.license.name }}.
  {% endif %}
  {% endif %}
</p>
{% elif config.license and config.license.name %}
<p class="footer-license">
  Content licensed under
  {% if config.license.url %}
  <a href="{{ config.license.url }}" target="_blank" rel="noopener noreferrer">{{ config.license.name }}</a>
  {% else %}
  {{ config.license.name }}
  {% endif %}
</p>
{% endif %}
```

The `config.license` object is intentionally empty when `license = false` so the paragraph is skipped, and is `nil` when the value is omitted (which triggers the live serve warning until an explicit string is configured).

---

## Error Handling

### Undefined Variables

| Mode | Behavior |
|------|----------|
| **Silent** (recommended) | Undefined → empty string |
| **Strict** | Undefined → error |
| **Debug** | Undefined → `{{ undefined }}` |

Configure silent mode in your template engine's configuration to handle undefined variables gracefully.

### Template Not Found

If a post's template doesn't exist:
1. Fall back to `post.html`
2. If `post.html` missing, error

### Template Syntax Errors

Report with:
- Template file path
- Line number
- Error description

---

## Performance

### Template Caching

Cache compiled templates using your template engine's bytecode caching feature to improve performance.

### Avoid N+1 Queries

Bad:
```jinja2
{% for post in posts %}
    {# This queries for each post #}
    {% set related = core.filter("'" ~ post.tags[0] ~ "' in tags")[:3] %}
{% endfor %}
```

Good:
```jinja2
{# Pre-compute in plugin, pass to template #}
{% for post in posts %}
    {% for related in post.related_posts %}
        ...
    {% endfor %}
{% endfor %}
```

---

## Library Recommendations

The template engine SHOULD support Jinja2-like syntax for cross-platform consistency.

### Common Template Libraries

| Language | Library | Notes |
|----------|---------|-------|
| Python | Jinja2 | Reference implementation |
| JavaScript | Nunjucks | Jinja-compatible, by Mozilla |
| Go | pongo2 | Jinja2-like syntax |
| Rust | Tera | Jinja2 inspired |

---

## Configuration

```toml
[tool-name.templates]
# Template directory
dir = "templates"

# Autoescape HTML
autoescape = true

# Undefined variable behavior
undefined = "silent"  # "silent", "strict", "debug"

# Cache compiled templates
cache = true

# Extensions to treat as templates
extensions = [".html", ".xml", ".txt"]

# Custom filters module
filters = "my_site.filters"
```

---

## Template Engine Differences

When implementing in different languages, template engines have subtle differences. This section documents key variations to ensure consistent behavior.

### Recommended Libraries by Language

| Language | Library | Notes |
|----------|---------|-------|
| Python | Jinja2 | Reference implementation |
| JavaScript | Nunjucks | Jinja-compatible, by Mozilla |
| Go | pongo2 | Jinja2-like syntax |
| Rust | Tera | Jinja2 inspired |

### Date Formatting Differences

Date formatting is a common source of incompatibility. The spec uses a **strftime-like** format as the canonical reference.

#### Format String Mapping

| Format | Description | Jinja2 (Python) | Nunjucks (JS) | pongo2 (Go) | Tera (Rust) |
|--------|-------------|-----------------|---------------|-------------|-------------|
| Full date | 2024-01-15 | `%Y-%m-%d` | `%Y-%m-%d` | `2006-01-02` | `%Y-%m-%d` |
| Month name | January | `%B` | `%B` | `January` | `%B` |
| Month abbr | Jan | `%b` | `%b` | `Jan` | `%b` |
| Day of month | 15 | `%d` | `%d` | `02` | `%d` |
| Year | 2024 | `%Y` | `%Y` | `2006` | `%Y` |
| Hour (24h) | 14 | `%H` | `%H` | `15` | `%H` |
| Minute | 30 | `%M` | `%M` | `04` | `%M` |
| ISO 8601 | 2024-01-15T14:30:00Z | `.isoformat()` | `.toISOString()` | `.Format(time.RFC3339)` | `%+` |

#### Go's Unique Format System

Go uses reference time (`Mon Jan 2 15:04:05 MST 2006`) instead of format codes:

```go
// Go (pongo2)
{{ post.date | date:"January 2, 2006" }}     // January 15, 2024
{{ post.date | date:"2006-01-02" }}          // 2024-01-15
{{ post.date | date:"Jan 02" }}              // Jan 15
```

**Implementation Requirement:** If the template engine uses a different date format system (like Go's reference time format), implement a `strftime` filter that translates standard format codes for cross-platform consistency:

```jinja2
{# Recommended: Add strftime filter for consistency #}
{{ post.date | strftime:"%B %d, %Y" }}  {# January 15, 2024 #}
```

#### Examples Across Engines

**Display "January 15, 2024":**

```jinja2
{# Jinja2 (Python) #}
{{ post.date.strftime('%B %d, %Y') }}
{{ post.date | date('%B %d, %Y') }}

{# Nunjucks (JavaScript) #}
{{ post.date | date('%B %d, %Y') }}

{# pongo2 (Go) - native format #}
{{ post.date | date:"January 02, 2006" }}

{# pongo2 (Go) - with strftime filter #}
{{ post.date | strftime:"%B %d, %Y" }}

{# Tera (Rust) #}
{{ post.date | date(format="%B %d, %Y") }}
```

**Display ISO 8601 for `<time>` element:**

```jinja2
{# Jinja2 #}
<time datetime="{{ post.date.isoformat() }}">

{# Nunjucks #}
<time datetime="{{ post.date.toISOString() }}">

{# pongo2 #}
<time datetime="{{ post.date | date:"2006-01-02T15:04:05Z07:00" }}">

{# Tera #}
<time datetime="{{ post.date | date(format="%+") }}">
```

### Filter Syntax Differences

| Operation | Jinja2 | Nunjucks | pongo2 | Tera |
|-----------|--------|----------|--------|------|
| Filter args | `\| truncate(100)` | `\| truncate(100)` | `\| truncate:100` | `\| truncate(length=100)` |
| Multiple args | `\| replace("a", "b")` | `\| replace("a", "b")` | N/A (use custom) | `\| replace(from="a", to="b")` |
| Chaining | `\| upper \| trim` | `\| upper \| trim` | `\| upper \| trim` | `\| upper \| trim` |

### Boolean/Truthiness Differences

| Value | Jinja2 | Nunjucks | pongo2 | Tera |
|-------|--------|----------|--------|------|
| Empty string `""` | Falsy | Falsy | Falsy | Falsy |
| Empty list `[]` | Falsy | Falsy | Truthy* | Falsy |
| Zero `0` | Falsy | Falsy | Falsy | Falsy |
| `None`/`nil`/`null` | Falsy | Falsy | Falsy | Falsy |

*pongo2 note: Empty slices are truthy by Go convention; use `| length > 0` for explicit checks.

### Safe/Raw HTML Output

| Engine | Syntax |
|--------|--------|
| Jinja2 | `{{ html \| safe }}` |
| Nunjucks | `{{ html \| safe }}` |
| pongo2 | `{{ html \| safe }}` |
| Tera | `{{ html \| safe }}` |

### Undefined Variable Handling

Configure silent undefined behavior for all engines:

```python
# Jinja2
from jinja2 import Environment, ChainableUndefined
env = Environment(undefined=ChainableUndefined)
```

```javascript
// Nunjucks
const env = nunjucks.configure({ throwOnUndefined: false });
```

```go
// pongo2 - undefined variables are empty by default
```

```rust
// Tera
let mut tera = Tera::new("templates/**/*")?;
tera.set_undefined_behavior(UndefinedBehavior::Silent);
```

### Required Custom Filters

Regardless of engine, implementations MUST provide these filters with consistent behavior:

| Filter | Signature | Description |
|--------|-----------|-------------|
| `slugify` | `(str) -> str` | Convert to URL-safe slug |
| `date` | `(datetime, format) -> str` | Format date (strftime syntax) |
| `reading_time` | `(str, wpm=200) -> str` | Estimate reading time |
| `filesizeformat` | `(int) -> str` | Human-readable file size |

### Testing Template Compatibility

Include these test cases to verify cross-engine behavior:

```yaml
template_engine_compat:
  - name: "date formatting consistency"
    input:
      template: '{{ date | date("%Y-%m-%d") }}'
      context:
        date: "2024-01-15T10:30:00Z"
    output: "2024-01-15"

  - name: "empty list truthiness"
    input:
      template: '{% if items | length > 0 %}yes{% else %}no{% endif %}'
      context:
        items: []
    output: "no"

  - name: "undefined variable silent"
    input:
      template: "Hello {{ name | default('World') }}"
      context: {}
    output: "Hello World"

  - name: "slugify filter"
    input:
      template: "{{ title | slugify }}"
      context:
        title: "Hello World!"
    output: "hello-world"
```

---

## Microformats2 Semantic Markup

Default templates include Microformats2 classes for IndieWeb compatibility.

### Post Pages (h-entry)

Single post templates MUST include `h-entry` markup:

Visible HTML dates SHOULD use the shared `human_date` filter so post pages, cards, archive views, and reader metadata render as `Apr 15, 2026` while `datetime` attributes remain machine-readable.

```html
<article class="post h-entry">
  <a class="u-url" href="{{ config.url }}{{ post.href }}" hidden></a>
  <h1 class="p-name">{{ post.title }}</h1>
  <time class="dt-published" datetime="{{ post.date | atom_date }}">{{ post.date | human_date }}</time>
  <div class="post-content e-content">{{ body | safe }}</div>
  {% for tag in post.tags %}
  <a class="p-category" href="/tags/{{ tag | slugify }}/">{{ tag }}</a>
  {% endfor %}
<span class="p-author h-card" hidden>
    <a class="u-url p-name" href="{{ config.url }}">{{ config.author }}</a>
</span>
</article>

### Article Reading Progress Indicator

To reinforce long-form reading flows, the default blog layout renders a sticky progress track before the `<article>` content. This indicator updates smoothly as the visitor scrolls and only appears on article pages (`article.post`). The markup looks like:

```html
<div class="article-progress" role="presentation" aria-hidden="true">
  <span class="article-progress__indicator"></span>
</div>
```

The CSS uses the following variables to control appearance (these default to palette-aware values defined in `variables.css`, but they can be overridden at `[markata-go.theme.variables]`):

| Variable | Purpose |
|----------|---------|
| `--article-progress-height` | Track thickness (default `4px`). |
| `--article-progress-track` | Background for the empty track. |
| `--article-progress-start` / `--article-progress-end` | Gradient for the filled portion. |
| `--article-progress-glow` | Subtle glow applied to the indicator for depth. |

JavaScript (the `initArticleProgressIndicator` initializer) keeps the indicator in sync with the viewport. It:

1. Locates `.article-progress__indicator` and `article.post`; aborts early when either is missing.
2. Computes progress as `clamp((window.innerHeight + scrollY - articleTop) / articleHeight, 0, 1)` so the indicator is full once the viewport reaches the bottom of the article.
3. Uses `requestAnimationFrame` to throttle updates and avoids redundant DOM writes when the ratio is unchanged.
4. Attaches `scroll` and `resize` listeners with `{ passive: true }`.
5. Respects `prefers-reduced-motion` by delegating transition control to CSS (`@media (prefers-reduced-motion: reduce)` disables the transform animation).
```

### Feed Pages (h-feed)

Feed/listing templates MUST include `h-feed` markup:

```html
<div class="feed h-feed">
  <h1 class="p-name">{{ feed.title }}</h1>
  <p class="p-summary">{{ feed.description }}</p>
  <span class="p-author h-card" hidden>
    <a class="u-url p-name" href="{{ config.url }}">{{ config.author }}</a>
  </span>
  {% for post in posts %}
  <article class="card h-entry">...</article>
  {% endfor %}
</div>
```

### Required Microformat Classes

| Class | Element | Description |
|-------|---------|-------------|
| `h-entry` | article | Entry container |
| `p-name` | h1/h2 | Entry title |
| `u-url` | a | Canonical permalink |
| `dt-published` | time | Publication datetime |
| `e-content` | div | Entry content (HTML) |
| `p-summary` | p | Entry summary/excerpt |
| `p-category` | a/span | Tags/categories |
| `p-author h-card` | span | Author information |
| `h-feed` | div | Feed container |
| `u-photo` | img | Photo content |
| `u-video` | video | Video content |

---

## Media Detection Filters

Templates can auto-detect media type (image vs video) using file extension analysis. This allows `image` and `video` frontmatter fields to be used interchangeably.

### `is_video` Filter

Returns `true` if the input string has a video file extension.

**Recognized video extensions:** `.mp4`, `.webm`, `.mov`, `.m4v`, `.ogv`, `.ogg`

Extension matching is case-insensitive.

```jinja2
{% if post.image|is_video %}
<video src="{{ post.image }}" autoplay muted loop playsinline></video>
{% else %}
<img src="{{ post.image }}" alt="{{ post.title }}">
{% endif %}
```

### `media_url` Filter

Resolves a media URL from multiple fields. Returns the first non-empty value from the input (primary) and parameter (fallback).
Trusted media URLs are normalized to `https` before output so templates can reuse the filter for raw image metadata.

```jinja2
{# Use image field first, fall back to video field #}
{% with post.image|media_url:post.video as media_src %}
{% if media_src %}
  {# Render media_src #}
{% endif %}
{% endwith %}
```

### Card Template Behavior

Photo and video card templates use both filters together to support interchangeable `image` and `video` frontmatter fields:

1. `media_url` resolves which field has a value (`image` preferred over `video`)
2. `is_video` determines whether to render a `<video>` or `<img>` element
3. Video MIME type is inferred from the file extension

### Card Media Sizing & Posters

OG cards, feed cards, and embed cards share the same media rules so that headless screenshotters, RSS readers, and embedded previews all see equivalent metadata. When any of these templates renders `image`/`cover_image`/`og_image`/`video` content, the helpers below guarantee consistent query parameters, poster selection, and host restrictions:

1. Use the `with_size(width, height)` helper (e.g., `post.cover_image|with_size:"1200,630"`) so the rendered media URL always includes explicit `w` and `h` query parameters that match the template's rendered width/height. These query params must stay in sync with the actual layout so caching layers and screenshot tools receive accurate dimensions. The helper only decorates relative URLs or URLs hosted on the trusted allowlist, and trusted media URLs are normalized to `https` before rendering to avoid mixed-content warnings.
2. Use `responsive_srcset` for a trusted remote image service that supports width-based resizing. Supply comma-separated candidate widths and a matching `sizes` value. For lazy images whose rendered width varies with layout, `sizes="auto, ...fallback..."` lets supporting browsers use the actual rendered width while older browsers use the conservative fallback. The filter emits a `srcset` only for trusted remote hosts; local images remain the responsibility of `image_optimization`, and unknown hosts keep their existing single-URL behavior.
   ```jinja2
   {% with post.image|responsive_srcset:"160,320,640" as image_srcset %}
   <img src="{{ post.image|with_size:'640' }}" srcset="{{ image_srcset }}"
        sizes="(max-width: 480px) 50vw, 33vw" alt="{{ post.title }}">
   {% endwith %}
   ```
   Use `original_media_url` for a deferred lightbox URL when the reader should open the original image instead of the selected thumbnail. It removes known size parameters while retaining tokens and unrelated query parameters.
3. Detect video media with a query/fragment-safe `is_video` filter so that `video.mp4?token=…` is still treated as a video, and use the companion `video_mime` filter to infer the correct `video/…` MIME type without relying on brittle `endswith` checks.
4. Resolve video posters with the `poster_url()` helper. The helper checks the following frontmatter aliases in order: `poster_image`, `poster`, `video_poster`, `video_thumbnail`, `thumbnail`, `thumb`. If none of those values exist and the resolved video URL is hosted on a trusted domain, `poster_url()` derives a `.webp` poster from the video path so you still get a preview image for allowlisted hosts. Templates can re-run `with_size` against the poster URL to keep query parameters aligned with the rendered dimensions.

Only media URLs that are either relative or hosted on the trusted allowlist are decorated with sizing parameters and posters. The default allowlist contains `dropper.wayl.one`, `dropper.waylonwalker.com`, and `dropper-dev.wayl.one`; the `[markata-go.templates.media]` configuration block lets sites add or replace hosts when they run their own CDN. Non-allowlisted hosts (and non-relative URLs) fall back to height/width-agnostic rendering without auto-generated posters to avoid leaking content from unknown vendors. Templates can read `config.templates.media.trusted_domains` to keep template logic aligned with the configured allowlist.

---

## See Also

- [SPEC.md](./SPEC.md) - Full specification
- [THEMES.md](./THEMES.md) - Theming and template overrides
- [CONFIG.md](./CONFIG.md) - Template configuration
- [CONTENT.md](./CONTENT.md) - Markdown processing
- [PLUGINS.md](./PLUGINS.md) - Plugin development

---

## Accessibility Requirements

Templates MUST follow these accessibility guidelines:

### Image Dimensions

All `<img>` tags MUST include explicit `width` and `height` attributes to prevent
Cumulative Layout Shift (CLS). Large content images SHOULD also include `loading="lazy"`.

### External Link Hints

Links that open in a new tab (`target="_blank"`) MUST include a visually-hidden
screen reader hint such as `<span class="visually-hidden">(opens in new tab)</span>`
so that assistive technology users are warned about the navigation change.

### Reduced Motion

CSS hover/transition effects MUST be disabled or reduced inside a
`@media (prefers-reduced-motion: reduce)` block to respect user motion preferences.
