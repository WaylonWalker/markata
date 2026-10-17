# Default Plugins Specification

This document specifies all built-in plugins that ship with the static site generator.

## Plugin Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                        DEFAULT PLUGIN SET                            │
├─────────────────────────────────────────────────────────────────────┤
│  CONFIGURATION PHASE                                                 │
│    └─ config_defaults    Set default configuration values           │
│                                                                      │
│  CONTENT DISCOVERY                                                   │
│    └─ glob               Find content files                          │
│                                                                      │
│  CONTENT LOADING                                                     │
│    └─ load               Parse frontmatter + content                 │
│                                                                      │
│  PRE-RENDER PROCESSING                                               │
│    ├─ auto_description   Generate descriptions from content         │
│    └─ jinja_md           Process template expressions in content     │
│    └─ tailwind           Run Tailwind CLI, set CSS/JS includes        │
│                                                                      │
│  RENDER                                                              │
│    ├─ render_markdown    Convert markdown to HTML                    │
│    └─ wikilinks          Resolve [[internal links]]                  │
│                                                                      │
│  POST-RENDER                                                         │
│    ├─ heading_anchors    Add anchor links to headings               │
│    ├─ md_video           Convert video images to video tags         │
│    ├─ toc                Generate table of contents                  │
│    ├─ link_collector     Track inlinks/outlinks between posts       │
│    └─ feeds              Generate feed collections                   │
│                                                                      │
│  COLLECT                                                             │
│    ├─ series             Auto-generate series feeds (PriorityEarly)  │
│    ├─ overwrite_check    Detect conflicting output paths             │
│    └─ prevnext           Calculate prev/next from feeds/series       │
│                                                                      │
│  OUTPUT                                                              │
│    ├─ publish_feeds      Write HTML/RSS/Atom/JSON/MD/TXT/Sitemap    │
│    ├─ publish_html       Write individual post HTML files           │
│    ├─ random_post        Write /random/ redirect page               │
│    ├─ well_known         Generate .well-known endpoints             │
│    ├─ llms_txt           Generate public /llms.txt discovery index   │
│    ├─ copy_assets        Copy static files                           │
│    ├─ redirects          Generate HTML redirect pages                │
│    ├─ css_minify         Minify CSS files (PriorityLast)             │
│    └─ js_minify          Minify JS files (PriorityLast)              │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Configuration Phase

### `config_defaults`

**Stage:** `load_config`

**Purpose:** Set sensible defaults for all configuration.

**Defaults:**
```toml
[name]
output_dir = "output"
content_dir = "."
assets_dir = "static"
templates_dir = "templates"

[markata-go.glob]
glob_patterns = ["**/*.md"]
use_gitignore = true
exclude_patterns = ["node_modules/**", ".git/**", "output/**"]

[markata-go.feeds]
default_items_per_page = 10
default_orphan_threshold = 3

[markata-go.feeds.default_formats]
html = true
rss = true
atom = false
json = false
markdown = false
text = false
sitemap = false
```

---

## Content Discovery

### `glob`

**Stage:** `glob`

**Purpose:** Discover content files matching configured patterns.

**Configuration:**
```toml
[markata-go.glob]
glob_patterns = ["posts/**/*.md", "pages/*.md"]
use_gitignore = true
exclude_patterns = ["**/draft-*", "**/wip-*"]
```

**Behavior:**
1. For each pattern in `glob_patterns`, find matching files
2. If `use_gitignore`, exclude files matching `.gitignore` rules
3. Apply `exclude_patterns` to filter results
4. Deduplicate (same file from multiple patterns)
5. Store in `core.files`

**Hook behavior:**

```
for pattern in config.glob.glob_patterns:
    for path in glob(pattern):
        if should_include(path, config.glob):
            core.files.append(path)
```

---

## Content Loading

### `load`

**Stage:** `load`

**Purpose:** Parse files into Post objects.

**Behavior:**
1. Read file content (UTF-8)
2. Extract YAML frontmatter between `---` delimiters
3. Parse frontmatter into dict
4. Create Post object with merged frontmatter + computed fields
5. Store in `core.posts`

**Computed fields:**
- `path` - Source file path
- `content` - Raw content (after frontmatter extraction)
- `slug` - From frontmatter, title, or filename
- `href` - `/{slug}/`

**Error handling:**
| Error | Behavior |
|-------|----------|
| File read error | Skip file, log warning |
| Invalid YAML | Skip file, log error with path + line |
| Invalid encoding | Try fallback encodings, then skip |

**Hook behavior:**

```
for path in core.files:
    try:
        raw = read_text(path, encoding='utf-8')
        frontmatter, content = parse_frontmatter(raw)
        post = Post(
            path=path,
            content=content,
            **frontmatter
        )
        core.posts.append(post)
    except Exception as e:
        log_warning("Failed to load {path}: {e}")
```
```

---

## Pre-Render Processing

### `auto_description`

**Stage:** `pre_render`

**Purpose:** Generate descriptions from content for posts that don't have one.

**Configuration:**
```toml
[markata-go.auto_description]
enabled = true
max_length = 160               # Characters
strip_html = true              # Remove HTML tags
strip_markdown = true          # Remove markdown syntax
fallback = ""                  # If content too short
```

**Behavior:**
1. For posts where `description` is None or empty
2. Take first paragraph of content
3. Strip HTML/markdown if configured
4. Truncate to `max_length` at word boundary
5. Add ellipsis if truncated
6. Set `post.description`

**Hook behavior:**

```
config = core.config.auto_description
if not config.enabled:
    return

for post in core.filter("description == None or description == ''"):
    post.description = generate_description(post.content, config)
```

---

### `jinja_md`

**Stage:** `pre_render`

**Purpose:** Process Jinja template expressions within markdown content.

**Configuration:**
```toml
[markata-go.jinja_md]
enabled = true
default_enabled = false        # Require explicit jinja: true in frontmatter
```

**Activation:**
- If `default_enabled = true`: All posts processed unless `jinja: false`
- If `default_enabled = false`: Only posts with `jinja: true` processed

**Template context:**
| Variable | Type | Description |
|----------|------|-------------|
| `post` | Post | Current post |
| `core` | Core | Core instance (for filtering) |
| `config` | Config | Site configuration |
| `today` | date | Current date |
| `now` | datetime | Current datetime |

**Example content:**
```markdown
---
title: Posts Index
jinja: true
---
# All Posts

{% for p in core.filter("published == True")[:10] %}
- [{{ p.title }}]({{ p.href }})
{% endfor %}
```

**Hook behavior:**

```
config = core.config.jinja_md
if not config.enabled:
    return

filter_expr = "jinja == True" if not config.default_enabled else "jinja != False"

for post in core.filter(filter_expr):
    template = template_engine.from_string(post.content)
    post.content = template.render(
        post=post,
        core=core,
        config=core.config,
        today=today(),
        now=now()
    )
```

---

## Collect Phase

### `series`

**Stage:** `collect` (with early priority, PriorityEarly = -100)

**Purpose:** Scan posts for `series` frontmatter and auto-generate series feed configs with prev/next navigation for guided sequential reading.

**Configuration:**
```toml
[markata-go.series]
slug_prefix = "series"     # URL prefix for series feeds (default: "series")
auto_sidebar = true        # Auto-enable feed sidebar on series posts (default: true)

[markata-go.series.defaults]
items_per_page = 0         # No pagination by default
sidebar = true             # Show sidebar on series posts

[markata-go.series.defaults.formats]
html = true
rss = true
atom = false
json = false

# Per-series overrides (keyed by slugified series name)
[markata-go.series.overrides."building-a-cli-in-go"]
title = "Building a CLI in Go"
description = "A step-by-step guide"
```

**Frontmatter:**
```yaml
series: "Building a CLI in Go"   # Series name (becomes slug)
series_order: 1                   # Optional explicit ordering
```

**Behavior:**
1. Scan all posts for `series` frontmatter
2. Group posts by series name (slugified)
3. Sort posts within each series by `series_order` (if set) or by date ascending
4. Set guide navigation (Prev/Next pointers) on published posts
5. Set `PrevNextContext` with position info (part X of Y)
6. Build `FeedConfig` entries with `FeedTypeSeries` type
7. Inject configs into `config.Extra["feeds"]` for the `feeds` plugin to process

**Constraints:**
- Posts can only belong to one series
- A series must have at least one post to generate a feed
- Runs before `feeds` and `overwrite_check` in the collect stage

---

### `overwrite_check`

**Stage:** `collect` (with early priority)

**Purpose:** Detect when multiple posts or feeds would write to the same output path, preventing accidental content overwrites.

**Configuration:**
```toml
[markata-go.overwrite_check]
enabled = true
warn_only = false              # If true, warn instead of fail
```

**Behavior:**
1. Build a map of all output paths and their content sources
2. For each post (not skipped/draft), calculate its output path: `{output_dir}/{slug}/index.html`
3. For each feed, calculate all output paths based on enabled formats (HTML, RSS, Atom, JSON)
4. Detect conflicts where multiple sources would write to the same path
5. If conflicts found:
   - `warn_only = false`: Return error with all conflicts listed
   - `warn_only = true`: Log warning but continue

**Output paths checked:**

| Content Type | Path Pattern |
|--------------|--------------|
| Post | `{output_dir}/{slug}/index.html` |
| Feed HTML | `{output_dir}/{feed_slug}/index.html` |
| Feed RSS | `{output_dir}/{feed_slug}/rss.xml` |
| Feed Atom | `{output_dir}/{feed_slug}/atom.xml` |
| Feed JSON | `{output_dir}/{feed_slug}/index.json` |
| Homepage | `{output_dir}/index.html` |

**PathConflict structure:**

| Field | Type | Description |
|-------|------|-------------|
| `output_path` | string | The conflicting output path |
| `sources` | list of strings | Content sources (e.g., "post:path/to/file.md", "feed:blog") |

**Error message format:**
```
detected 2 output path conflict(s):
  - output/blog/index.html: post:posts/blog.md, feed:blog
  - output/about/index.html: post:about.md, post:pages/about.md
```

**Common conflict scenarios:**

| Scenario | Example | Resolution |
|----------|---------|------------|
| Post slug matches feed slug | Post with `slug: blog` and feed with `slug: blog` | Rename post slug or feed slug |
| Duplicate post slugs | Two posts with same `slug` frontmatter | Use unique slugs |
| Empty slug conflicts | Homepage feed and post with `slug: ""` | Only one can be the homepage |

**Hook behavior:**

```
for post in posts where not (skip or draft):
    output_path = join(output_dir, post.slug, "index.html")
    add to path_sources[output_path]: "post:{post.path}"

for feed_config in feed_configs:
    for path in get_feed_output_paths(output_dir, feed_config):
        add to path_sources[path]: "feed:{feed_config.slug}"

conflicts = []
for path, sources in path_sources:
    if length(sources) > 1:
        add to conflicts: PathConflict(output_path=path, sources=sources)

if length(conflicts) > 0 and not warn_only:
    return error "detected N output path conflict(s)..."
```

**Interface requirements:**

The plugin MUST implement:
- `Plugin` - Basic plugin interface with `Name()` method
- `CollectPlugin` - To run during the collect stage
- `PriorityPlugin` - To ensure early execution before other collect plugins

---

### `prevnext`

**Stage:** `collect` (after feeds are created)

**Purpose:** Calculate previous/next post links for navigation based on feeds.

**Configuration:**
```toml
[markata-go.prevnext]
enabled = true
strategy = "first_feed"        # Resolution strategy (see below)
default_feed = "blog"          # Feed to use when strategy = "explicit_feed"
```

**Strategy Options:**

| Strategy | Description |
|----------|-------------|
| `first_feed` | Use first feed the post appears in (default) |
| `explicit_feed` | Always use `default_feed` for all posts |
| `series` | Use post's `series` frontmatter to find matching feed, fall back to `first_feed` |
| `frontmatter` | Use post's `prevnext_feed` frontmatter to find matching feed, fall back to `first_feed` |

**Post frontmatter options:**
```yaml
---
title: My Post
series: python-tutorial       # Feed slug to use for navigation (with strategy="series")
prevnext_feed: blog           # Explicit feed for this post's navigation (with strategy="frontmatter")
---
```

**Post fields added:**
| Field | Type | Description |
|-------|------|-------------|
| `Prev` | optional Post | Previous post in sequence (null if first) |
| `Next` | optional Post | Next post in sequence (null if last) |
| `PrevNextFeed` | string | Feed slug used for navigation |
| `PrevNextContext` | optional PrevNextContext | Full navigation context |

**PrevNextContext structure:**

| Field | Type | Description |
|-------|------|-------------|
| `feed_slug` | string | Feed slug |
| `feed_title` | string | Feed title |
| `position` | integer | Position in sequence (1-indexed) |
| `total` | integer | Total posts in sequence |
| `prev` | optional Post | Previous post (null if first) |
| `next` | optional Post | Next post (null if last) |

**Behavior:**

1. **Build post-to-feeds mapping** - For each feed, map post slugs to feeds containing them
2. **For each post, determine navigation context based on strategy:**
   - `first_feed`: Find first feed containing this post
   - `explicit_feed`: Always use `default_feed`
   - `series`: Check post's `Extra["series"]`, look up matching feed, fall back to first_feed
   - `frontmatter`: Check post's `Extra["prevnext_feed"]`, look up matching feed, fall back to first_feed
3. **Find post's position in the feed's posts list**
4. **Set Prev/Next based on position**

**Hook behavior:**

```
config = get_prevnext_config(manager.config)
if not config.enabled:
    return

feeds = manager.feeds()
post_to_feeds = build_post_to_feeds_map(feeds)

for post in manager.posts():
    feed = resolve_feed(post, config, feeds, post_to_feeds)
    if feed is null:
        continue

    # Find position and set prev/next
    for i, feed_post in enumerate(feed.posts):
        if feed_post.slug == post.slug:
            if i > 0:
                post.prev = feed.posts[i-1]
            if i < length(feed.posts) - 1:
                post.next = feed.posts[i+1]
            post.prev_next_feed = feed.name
            post.prev_next_context = PrevNextContext(
                feed_slug=feed.name,
                feed_title=feed.title,
                position=i + 1,
                total=length(feed.posts),
                prev=post.prev,
                next=post.next
            )
            break
```

**Template usage:**
```jinja2
{# Basic prev/next navigation #}
<nav class="post-navigation">
  {% if post.Prev %}
  <a href="{{ post.Prev.Href }}" class="nav-prev">
    <span class="nav-label">Previous</span>
    <span class="nav-title">{{ post.Prev.Title }}</span>
  </a>
  {% endif %}

  {% if post.Next %}
  <a href="{{ post.Next.Href }}" class="nav-next">
    <span class="nav-label">Next</span>
    <span class="nav-title">{{ post.Next.Title }}</span>
  </a>
  {% endif %}
</nav>

{# Navigation with position indicator #}
{% if post.PrevNextContext %}
<nav class="series-navigation">
  <div class="series-info">
    <span class="series-title">{{ post.PrevNextContext.FeedTitle }}</span>
    <span class="series-position">
      Part {{ post.PrevNextContext.Position }} of {{ post.PrevNextContext.Total }}
    </span>
  </div>

  <div class="series-links">
    {% if post.PrevNextContext.Prev %}
    <a href="{{ post.PrevNextContext.Prev.Href }}">← {{ post.PrevNextContext.Prev.Title }}</a>
    {% endif %}
    {% if post.PrevNextContext.Next %}
    <a href="{{ post.PrevNextContext.Next.Href }}">{{ post.PrevNextContext.Next.Title }} →</a>
    {% endif %}
  </div>
</nav>
{% endif %}
```

**Example configurations:**

```toml
# Strategy 1: Use first feed (default)
[markata-go.prevnext]
enabled = true
strategy = "first_feed"

# Strategy 2: Always use specific feed
[markata-go.prevnext]
enabled = true
strategy = "explicit_feed"
default_feed = "all-posts"

# Strategy 3: Series-based navigation (uses series frontmatter to find feed)
[markata-go.prevnext]
enabled = true
strategy = "series"
```

**Frontmatter examples:**

```yaml
---
# Post using series frontmatter (strategy="series")
title: "Variables in Python"
series: python-basics    # Must match a feed slug
---

---
# Post with explicit navigation feed (strategy="frontmatter")
title: "Announcement Post"
prevnext_feed: announcements
---
```

---

## Render

### `render_markdown`

**Stage:** `render`

**Purpose:** Convert markdown content to HTML.

**Configuration:**
```toml
[markata-go.markdown]
backend = "auto"               # "markdown-it", "commonmark", etc.

[markata-go.markdown.extensions]
tables = true
admonitions = true
footnotes = true
strikethrough = true
task_lists = true
heading_ids = true

[markata-go.markdown.highlight]
enabled = true
theme = "github-dark"
line_numbers = false
guess_language = true
```

**Supported syntax:**
- CommonMark base
- GFM tables
- Admonitions (`!!! note "Title"`)
- Fenced code blocks with language
- Syntax highlighting
- Footnotes
- Strikethrough
- Task lists

**Output:**
Sets `post.article_html` to rendered HTML (content only, no template wrapper).

**Hook behavior:**

```
for post in core.filter("not skip"):
    post.article_html = markdown_parser.render(post.content)
```

---

### `wikilinks`

**Stage:** `render` (after `render_markdown`)

**Purpose:** Resolve `[[internal links]]` to actual post URLs.

**Configuration:**
```toml
[markata-go.wikilinks]
enabled = true
warn_broken = true             # Warn about broken links
broken_class = "broken-link"   # CSS class for broken links
```

**Syntax:**
```markdown
Link to post: [[other-post-slug]]
With custom text: [[other-post-slug|Click here]]
Link to section: [[other-post-slug#section-id]]
Section with text: [[other-post-slug#section-id|See this section]]
```

**Resolution:**
1. Split `link_target` on the first `#` into `slug` and optional `fragment`
2. Find post where `slug == slug` (or alias match, case-insensitive)
3. If found: `<a href="{post.href}#{fragment}">{text or post.title}</a>` (fragment omitted if empty)
4. If not found: Leave as `[[link]]` or wrap with broken-link class

**Hook behavior:**

This plugin runs after `render_markdown` due to late priority.

```
config = core.config.wikilinks
if not config.enabled:
    return

# Build slug -> post lookup
slug_map = {p.slug: p for p in core.posts}

for post in core.filter("not skip"):
    post.article_html = resolve_wikilinks(
        post.article_html,
        slug_map,
        config
    )
```

---

## Post-Render

### `heading_anchors`

**Stage:** `render` (with late priority, after `render_markdown`)

**Purpose:** Add anchor links to headings for direct linking.

**Configuration:**
```toml
[markata-go.heading_anchors]
enabled = true
min_level = 2                  # Start at h2
max_level = 4                  # End at h4
position = "end"               # "start" or "end"
symbol = "#"                   # Link text
class = "heading-anchor"       # CSS class
```

**Behavior:**
1. Find all headings in `article_html`
2. Skip headings that already contain a goldmark anchor (`class="anchor"`) to prevent duplicates
3. Generate ID from heading text (slugified)
4. Handle duplicate IDs by appending numbers
5. Insert anchor link at configured position

**Note:** The goldmark markdown renderer has its own anchor extension (`[markdown.extensions] anchor`), which is disabled by default. The `heading_anchors` plugin is the recommended way to add heading anchors because it offers more configuration options. If the goldmark anchor extension is explicitly re-enabled, `heading_anchors` will detect and skip headings that already have goldmark anchors to avoid duplicates.

**Example output:**
```html
<h2 id="my-section">
  My Section
  <a href="#my-section" class="heading-anchor">#</a>
</h2>
```

**Hook behavior:**

```
config = core.config.heading_anchors
if not config.enabled:
    return

for post in core.filter("not skip"):
    post.article_html = add_heading_anchors(post.article_html, config)
```

---

### `md_video`

**Stage:** `render` (with late priority, after `render_markdown`)

**Purpose:** Convert markdown image syntax for video files into native video elements that load when the reader requests playback.

**Configuration:**
```toml
[markata-go.md_video]
enabled = true
video_extensions = [".mp4", ".webm", ".ogg", ".ogv", ".mov", ".m4v", ".avi", ".mkv"]
video_class = "md-video"
controls = true
autoplay = false                   # Opt in to automatic playback
loop = true                        # GIF-like behavior
muted = true                       # Required for autoplay
playsinline = true                 # Inline on mobile
preload = "none"
```

**Configuration Fields:**

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `true` | Enable/disable the plugin |
| `video_extensions` | []string | `[".mp4", ".webm", ...]` | Extensions to treat as video |
| `video_class` | string | `"md-video"` | CSS class for video elements |
| `controls` | bool | `true` | Show video controls |
| `autoplay` | bool | `false` | Auto-start playback; off by default to avoid surprise media downloads |
| `loop` | bool | `true` | Loop continuously |
| `muted` | bool | `true` | Mute audio |
| `playsinline` | bool | `true` | Play inline on iOS |
| `preload` | string | `"none"` | Preload hint |

**Behavior:**
1. Find all `<img>` tags in `article_html` after markdown rendering
2. Check if `src` ends with a video extension (handles query params)
3. Replace with `<video>` element with configured attributes
4. Preserve alt text as fallback content
5. Auto-detect MIME type from extension
6. Resolve a poster image for trusted video URLs across every supported video format

**Markdown syntax:**
```markdown
![kickflip down the 3 stair](https://example.com/video.mp4)

![Demo](demo.webm?width=500)
```

**Example output:**
```html
<video loop muted playsinline controls preload="none" class="md-video">
  <source src="https://example.com/video.mp4" type="video/mp4">
  kickflip down the 3 stair
</video>
```

**MIME type detection:**

| Extension | MIME Type |
|-----------|-----------|
| `.mp4` | `video/mp4` |
| `.webm` | `video/webm` |
| `.ogg`, `.ogv` | `video/ogg` |
| `.mov` | `video/quicktime` |
| `.m4v` | `video/x-m4v` |

**Interface requirements:**

The plugin MUST implement:
- `Plugin` - Basic plugin interface with `Name()` method
- `ConfigurePlugin` - To read configuration
- `RenderPlugin` - To process posts during the render stage
- `PriorityPlugin` - To ensure late execution after markdown rendering

---

### `toc`

**Stage:** `post_render`

**Purpose:** Generate table of contents from headings.

**Configuration:**
```toml
[markata-go.toc]
enabled = true
min_level = 2
max_level = 4
min_headings = 3               # Only generate if >= 3 headings
```

**Post fields added:**
| Field | Type | Description |
|-------|------|-------------|
| `toc` | string | HTML table of contents |
| `toc_items` | list of TocItem | Structured TOC data |

**TocItem structure:**

| Field | Type | Description |
|-------|------|-------------|
| `level` | integer | Heading level (2-6) |
| `id` | string | Heading ID |
| `text` | string | Heading text |
| `children` | list of TocItem | Nested headings |

**Example output:**
```html
<nav class="toc">
  <ul>
    <li><a href="#introduction">Introduction</a></li>
    <li>
      <a href="#main">Main</a>
      <ul>
        <li><a href="#subsection">Subsection</a></li>
      </ul>
    </li>
  </ul>
</nav>
```

**Template usage:**
```jinja2
{% if post.toc %}
<aside class="sidebar">
  {{ post.toc | safe }}
</aside>
{% endif %}
```

---

### `link_collector`

**Stage:** `render` (after `render_markdown`)

**Purpose:** Collect all hyperlinks from posts and track inlinks (pages linking TO this post) and outlinks (pages this post links TO). Enables backlink navigation and link graph visualization.

**Configuration:**
```toml
[markata-go.link_collector]
enabled = true
include_feeds = false          # Exclude feed pages from inlinks by default
include_index = false          # Exclude index page from inlinks by default
```

**Post fields added:**
| Field | Type | Description |
|-------|------|-------------|
| `hrefs` | list of strings | Raw href values from all links in post |
| `inlinks` | list of Link | Links pointing TO this post from other posts |
| `outlinks` | list of Link | Links FROM this post to other pages |

**Link structure:**

| Field | Type | Description |
|-------|------|-------------|
| `source_url` | string | Absolute URL of the source post |
| `source_post` | Post | Reference to source post object |
| `target_post` | optional Post | Reference to target post (null if external) |
| `raw_target` | string | Original href value |
| `target_url` | string | Resolved absolute URL |
| `target_domain` | optional string | Domain extracted from target_url |
| `is_internal` | boolean | True if link points to same site |
| `is_self` | boolean | True if link points to same post |
| `source_text` | string | Cleaned link text from source |
| `target_text` | string | Cleaned link text from target |

**Behavior:**
1. Parse `article_html` for all `<a href="...">` elements
2. For each href:
   - Resolve relative URLs against post's base URL
   - Determine if internal (same domain) or external
   - Look up target post by slug if internal
   - Create Link object with all metadata
3. Store all links in `core.links`
4. For each post, populate:
   - `inlinks`: Links where `target_post == this_post` (deduplicated by source_url)
   - `outlinks`: Links where `source_post == this_post` (deduplicated by target_url)
5. Exclude self-links from both inlinks and outlinks

**Caching:**
Per-post href extraction is cached using a hash of `article_html`. If the hash
matches on a subsequent build, the plugin reuses cached hrefs instead of
re-parsing HTML. Link objects are still rebuilt each run to keep target
resolution and inlinks/outlinks consistent.

**Hook behavior:**

```
links = []
site_domain = parse_url(config.url).host

# Build slug -> post lookup
post_by_slug = {p.slug: p for p in posts}

for post in posts:
    base_url = url_join(config.url, post.slug)
    soup = parse_html(post.article_html)

    # Optionally limit to post-body section
    post_body = soup.find_by_id("post-body")
    if post_body:
        soup = post_body

    post.hrefs = [a.href for a in soup.find_all("a", has_href=true)]

    for href in post.hrefs:
        target_url = url_join(base_url, href)
        domain = parse_url(target_url).host
        is_internal = (domain == site_domain)

        target_post = null
        if is_internal:
            target_slug = parse_url(target_url).path.strip("/")
            target_post = post_by_slug.get(target_slug)

        links.append(Link(
            source_url=base_url,
            source_post=post,
            target_post=target_post,
            raw_target=href,
            target_url=target_url,
            target_domain=domain,
            is_internal=is_internal,
            is_self=(target_post and post.slug == target_post.slug)
        ))

core.links = links

# Assign inlinks/outlinks to each post
for post in posts:
    post.inlinks = [
        link for link in links
        if link.target_post == post and not link.is_self
    ]
    post.outlinks = [
        link for link in links
        if link.source_post == post and not link.is_self
    ]
```

**Template usage - Basic inlinks/outlinks:**
```jinja2
{% if post.inlinks %}
<section class="inlinks">
  <h2>Pages that link here</h2>
  <ul>
    {% for link in post.inlinks %}
    {% if link.source_post.slug not in core.feeds.keys() %}
    <li><a href="{{ link.source_url }}">{{ link.source_post.title }}</a></li>
    {% endif %}
    {% endfor %}
  </ul>
</section>
{% endif %}

{% if post.outlinks %}
<section class="outlinks">
  <h2>Links from this page</h2>
  <ul>
    {% for link in post.outlinks %}
    {% if link.target_post %}
    <li><a href="{{ link.target_url }}">{{ link.target_post.title }}</a></li>
    {% else %}
    <li><a href="{{ link.target_url }}">{{ link.target_url }}</a></li>
    {% endif %}
    {% endfor %}
  </ul>
</section>
{% endif %}
```

**Template usage - Link graph visualization (Mermaid):**
```jinja2
{% if post.inlinks and post.outlinks %}
<section class="link-graph">
  <h2>Link Graph</h2>
  <pre class="mermaid">
graph LR

{% for link in post.inlinks %}
{% if link.source_post.slug != 'index' and link.source_post.slug not in core.feeds.keys() %}
    {{ link.source_post.slug }}:::inlink --> {{ post.slug }}:::this
{% endif %}
{% endfor %}

{% for link in post.outlinks %}
{% if link.target_post %}
    {{ post.slug }}:::this --> {{ link.target_post.slug }}:::outlink
{% else %}
    {{ post.slug }}:::this --> {{ link.target_text }}:::outlink
{% endif %}
{% endfor %}

    classDef this stroke:#ffcc00
    classDef outlink stroke:#50C878
    classDef inlink stroke:#75E6DA
  </pre>
</section>
{% endif %}
```

**CSS for link sections:**
```css
.inlinks, .outlinks, .link-graph {
  margin-top: 3rem;
  padding: 1rem;
  border-top: 1px solid var(--border-color, #e0e0e0);
}

.inlinks h2, .outlinks h2, .link-graph h2 {
  font-size: 1.25rem;
  margin-bottom: 0.5rem;
}

.inlinks ul, .outlinks ul {
  list-style: none;
  padding: 0;
}

.inlinks li, .outlinks li {
  margin: 0.25rem 0;
}

/* Mermaid graph styling */
.link-graph .mermaid {
  background: var(--bg-secondary, #f5f5f5);
  padding: 1rem;
  border-radius: 4px;
}
```

**Utility functions:**

| Function | Description |
|----------|-------------|
| `count_links(links)` | Count target_url frequency across all links |
| `count_domains(links)` | Count target_domain frequency across all links |
| `get_external_links(links)` | Filter to only external links |
| `get_broken_links(links)` | Find internal links that don't resolve to a post |

---

## Feeds & Output

### `feeds`

**Stage:** `post_render`

**Purpose:** Generate feed collections from post queries.

See [FEEDS.md](./FEEDS.md) for complete specification.

**Core behavior:**
1. Read feed definitions from config
2. Generate auto-feeds (tags, categories, dates) if enabled
3. For each feed:
   - Run filter expression
   - Sort results
   - Paginate if configured
   - Store in `core.feeds`

**Hook behavior:**

```
core.feeds = []

# Process explicit feeds
for feed_config in core.config.feeds:
    feed = create_feed(feed_config, core)
    core.feeds.append(feed)

# Process auto-feeds
if core.config.feeds.auto_tags.enabled:
    for tag in get_all_tags(core.posts):
        feed = create_tag_feed(tag, core)
        core.feeds.append(feed)
```

---

### `publish_feeds`

**Stage:** `save`

**Purpose:** Write feed output files in all configured formats.

**Behavior:**
For each feed, for each enabled format:
1. Load format template
2. Render with feed context
3. Write to output path

**Performance notes:**
- Feed filtering results may be reused across feeds that share identical filter expressions and privacy settings within a single build.

**Output paths:**
| Format | Path Pattern |
|--------|--------------|
| HTML | `/{slug}/index.html`, `/{slug}/page/{n}/index.html` |
| RSS | `/{slug}/rss.xml` |
| Atom | `/{slug}/atom.xml` |
| JSON | `/{slug}/feed.json` |
| Markdown | `/{slug}.md` |
| Text | `/{slug}.txt` |
| Sitemap | `/{slug}/sitemap.xml` |

**Hook behavior:**

```
for feed in core.feeds:
    for format_name, enabled in feed.formats.items():
        if not enabled:
            continue

        template = get_feed_template(format_name, feed, core)

        if format_name == 'html':
            # Paginated output
            for page_num, page_posts in enumerate(feed.pages, start=1):
                output = template.render(
                    feed=feed,
                    page_posts=page_posts,
                    pagination=get_pagination(feed, page_num)
                )
                write_page(feed, page_num, output, core)
        else:
            # Single file output
            output = template.render(feed=feed)
            write_feed_file(feed, format_name, output, core)
```

---

### `publish_html`

**Stage:** `save`

**Purpose:** Write individual post HTML files.

**Behavior:**
1. For each post (not skipped)
2. Load post template
3. Render with post context
4. Write to `{output_dir}/{slug}/index.html`

**Template context:**
| Variable | Type | Description |
|----------|------|-------------|
| `post` | Post | The post being rendered |
| `body` | string | `post.article_html` |
| `config` | Config | Site configuration |
| `core` | Core | Core instance |

**Hook behavior:**

```
for post in core.filter("not skip"):
    template_name = post.template or "post.html"
    template = template_engine.get_template(template_name)

    html = template.render(
        post=post,
        body=post.article_html,
        config=core.config,
        core=core
    )

    output_path = join(core.config.output_dir, post.slug, "index.html")
    create_parent_dirs(output_path)
    write_text(output_path, html)
```

---

### `well_known`

**Stage:** `write`

**Purpose:** Generate `.well-known` endpoints from site configuration data.

**Behavior:**
1. Resolve `well_known` configuration
2. Determine enabled entries (defaults + optional)
3. Render templates (if available) or fall back to built-in defaults
4. Write files to the output directory

**Outputs:**
- `/.well-known/host-meta`
- `/.well-known/host-meta.json`
- `/.well-known/webfinger`
- `/.well-known/nodeinfo`
- `/nodeinfo/2.0`
- `/.well-known/time`
- `/.well-known/links`
- `/.well-known/internal-links`
- `/external-links/` (`index.html`)
- `/internal-links/` (`index.html`)
- `/.well-known/sshfp` (optional)
- `/.well-known/keybase.txt` (optional)

**Template context:**
| Variable | Type | Description |
|----------|------|-------------|
| `config` | Config | Site configuration |
| `well_known` | map | Derived values (site URL, host, build time, author, etc.) |
| `well_known_links` | array | Outbound link buckets (domain, count, list of source/target pairs) |
| `well_known_internal_links` | array | Internal link buckets (target URL, count, list of source targets) |

**Hook behavior:**

```
if config.well_known.enabled:
    entries = resolve_entries(config.well_known)
    for entry in entries:
        output = render_template(entry.template, config, well_known_data)
        write_file(output_dir, entry.path, output)
```

---

### `copy_assets`

**Stage:** `save`

**Purpose:** Copy static assets to output directory.

**Configuration:**
```toml
[markata-go.assets]
dir = "static"
output_subdir = ""             # "" = root, "assets" = /assets/
exclude = ["*.psd", "*.ai"]
fingerprint = false

[markata-go.assets.fingerprint]
enabled = false
algorithm = "sha256"
length = 8
exclude = ["robots.txt", "favicon.ico"]
```

**Behavior:**
1. Find all files in `assets_dir`
2. Apply exclusion patterns
3. If fingerprinting enabled:
   - Compute content hash
   - Rename: `style.css` → `style.a1b2c3d4.css`
   - Update manifest
4. Copy to output directory
5. Preserve directory structure

**Hook behavior:**

```
config = core.config.assets
assets_dir = Path(config.dir)

if not assets_dir.exists():
    return

for src in assets_dir.recursive_glob("*"):
    if src.is_file() and not is_excluded(src, config.exclude):
        rel_path = src.relative_to(assets_dir)

        if should_fingerprint(src, config):
            rel_path = fingerprint_path(src, rel_path, config)

        dst = join(core.config.output_dir, config.output_subdir, rel_path)
        create_parent_dirs(dst)
        copy_file(src, dst)
```

---

### `redirects`

**Stage:** `write` (with late priority, after content is written)

**Purpose:** Generate HTML redirect pages from a `_redirects` file. Creates static HTML pages at old URLs that redirect browsers to new URLs.

**Configuration:**
```toml
[markata-go.redirects]
redirects_file = "static/_redirects"    # Path to redirects file (default)
redirect_template = ""                   # Optional custom template path
```

**Redirects file format:**

The `_redirects` file follows Netlify's format - one redirect per line with source and destination paths:

```
# Comments start with #
/old-path /new-path
/blog/old-post /blog/new-post
/legacy/page /modern/page
```

**Syntax rules:**
- Lines starting with `#` are comments
- Empty lines are ignored
- Each line has two space-separated paths: `<original> <new>`
- Both paths must start with `/`
- Wildcards (`*`) are not supported (skipped)

**Behavior:**
1. Read the `_redirects` file (skip silently if not found)
2. Parse each line into redirect rules
3. For each redirect rule:
   - Create directory at `{output_dir}/{original_path}/`
   - Generate `index.html` with meta refresh redirect
4. Cache results based on file content hash

**Output:**
For each redirect rule, creates `{output_dir}/{original_path}/index.html`:

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta http-equiv="Refresh" content="0; url='/new-path'" />
  <meta charset="UTF-8">
  <link rel="canonical" href="/new-path" />
  <meta name="description" content="/old-path has been moved to /new-path." />
  <title>/old-path has been moved to /new-path</title>
  <!-- Styled fallback content for users with JS disabled -->
</head>
<body>
  <h1>Page Moved</h1>
  <p><code>/old-path</code> has moved to <a href="/new-path">/new-path</a></p>
</body>
</html>
```

**Custom templates:**

If `redirect_template` is set, the plugin loads a custom Go template with these variables:

| Variable | Type | Description |
|----------|------|-------------|
| `.Original` | string | The source path |
| `.New` | string | The destination path |
| `.Config` | Config | Site configuration |

**Example custom template:**
```html
<!DOCTYPE html>
<html>
<head>
  <meta http-equiv="Refresh" content="0; url='{{ .New }}'" />
  <link rel="canonical" href="{{ .New }}" />
</head>
<body>
  <p>Redirecting to <a href="{{ .New }}">{{ .New }}</a>...</p>
</body>
</html>
```

**Hook behavior:**

```
content = read_file(config.redirects_file)
if file_not_found:
    return  # Skip silently

# Parse and generate redirect pages
for redirect in parse_redirects(content):
    output_path = join(output_dir, redirect.original, "index.html")
    # Render template and write file
```

**Example configuration:**

```toml
# Basic usage (default settings)
[markata-go.redirects]
# Uses static/_redirects by default

# Custom redirects file location
[markata-go.redirects]
redirects_file = "_redirects"

# With custom template
[markata-go.redirects]
redirects_file = "config/_redirects"
redirect_template = "templates/redirect.html"
```

**Example `_redirects` file:**

```
# Blog post renames
/blog/old-title /blog/new-title
/posts/draft-post /posts/published-post

# Section reorganization
/tutorials/beginner /guides/getting-started
/tutorials/advanced /guides/advanced-topics

# Legacy URLs
/about-me /about
/contact-us /contact
```

---

### `random_post`

**Stage:** `write`

**Purpose:** Generate a static `/random/` endpoint that redirects client-side to a random eligible post.

**Specification:** See `spec/spec/RANDOM_POST.md`.

---

### `css_minify`

**Stage:** `write` (with `PriorityLast`)

**Purpose:** Minify all CSS files in the output directory to reduce file sizes and improve Lighthouse performance scores.

**Configuration:**

```toml
[markata-go.css_minify]
enabled = true                    # Enable CSS minification (default: true)
exclude = ["variables.css"]       # Files to skip (exact names or glob patterns)
preserve_comments = ["Copyright"] # Strings that mark comments to preserve
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | boolean | `true` | Enable/disable CSS minification |
| `exclude` | list of strings | `[]` | File patterns to skip (exact or glob) |
| `preserve_comments` | list of strings | `[]` | Substrings that mark comments to keep |

**Behavior:**

1. Skip if `enabled` is `false` or fast mode is active
2. Walk the output directory recursively to find all `.css` files
3. Skip `_pagefind`, symlinks, and exclusion patterns (exact match and glob)
4. Read content once and apply the persistent reuse rules below
5. On a transform, extract preserved comments and minify using `tdewolff/minify/v2/css`
6. Prepend preserved comments to the engine result without changing comment semantics
7. Persist the relation before atomic target replacement; log separate work counts and transformed size reduction

**Hook behavior:**

```
if not config.enabled:
    return

css_files = find_files(output_dir, "*.css")
for file in css_files:
    if is_excluded(file, config.exclude):
        skip
    apply_cached_or_fresh_css_minification(file, recipe, cache)

log_stats(total_original, total_minified)
```

**Interface requirements:**
- `Plugin` (Name)
- `ConfigurePlugin` (read config from Extra)
- `WritePlugin` (minify files)
- `PriorityPlugin` (return `PriorityLast` for Write stage)

---

### `js_minify`

**Stage:** `write` (with `PriorityLast`)

**Purpose:** Minify all JavaScript files in the output directory to reduce file sizes and improve Lighthouse performance scores.

**Configuration:**

```toml
[markata-go.js_minify]
enabled = true                    # Enable JS minification (default: true)
exclude = ["pagefind-ui.js"]      # Files to skip (exact names or glob patterns)
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | boolean | `true` | Enable/disable JS minification |
| `exclude` | list of strings | `[]` | File patterns to skip (exact or glob) |

**Behavior:**

1. Skip if `enabled` is `false` or fast mode is active
2. Walk the output directory recursively to find all `.js` files
3. Skip `.min.js` files (already minified)
4. Skip `_pagefind`, symlinks, and exclusion patterns (exact match and glob)
5. Read content once and apply the persistent reuse rules below
6. On a transform, minify using `tdewolff/minify/v2/js`
7. Persist the relation before atomic target replacement
8. Log separate work counts and transformed size reduction

**Hook behavior:**

```
if not config.enabled:
    return

js_files = find_files(output_dir, "*.js", exclude="*.min.js")
for file in js_files:
    if is_excluded(file, config.exclude):
        skip
    apply_cached_or_fresh_js_minification(file, recipe, cache)

log_stats(total_original, total_minified)
```

**Interface requirements:**
- `Plugin` (Name)
- `ConfigurePlugin` (read config from Extra)
- `WritePlugin` (minify files)
- `PriorityPlugin` (return `PriorityLast` for Write stage)

**Shared infrastructure:**

Both `css_minify` and `js_minify` use shared helper functions in `minify_helpers.go`:
- `runMinification` - Processes eligible assets with bounded concurrency and persistent result reuse
- `isExcludedByPatterns(filename, excludeMap)` - Checks exact and glob exclusion patterns

**Persistent asset minification (version 1):**

- Cache independently of the build-cache enabled setting, under
  `ContentDir/.markata/asset-minify/v1/<scopehash>/`. A nonempty `cache_dir`
  override is used as supplied, matching build-cache configuration. The scope
  includes the absolute normalized site root, configured output root, and kind.
  Cache storage MUST remain outside published output, including symlink-resolved
  relationships. Validate the derived storage directories and every ancestor
  that creation or permission changes could affect BEFORE performing those
  operations. Private cache parents must not contain published output either.
  Unsafe storage is warned about and disabled without changing output permissions.
  Missing ancestors MUST undergo overlap validation before treating `ENOTDIR`
  as a non-directory fallback: on Windows, `ENOTDIR` also represents a missing
  path. Preserve the fallback for genuine non-directory ancestors and report
  other inspection failures explicitly.
- Private directories use mode `0700`; records and content-addressed blobs use
  `0600`. Per-asset records live in `records/<relativepathhash>.json` and blobs
  in `blobs/<sha256>`. Only discovered output assets authorize target writes.
  Validate schema, scope, relative identity, and all SHA-256 values before use.
- The recipe includes kind, a maintained wrapper revision, minify and parse
  dependency versions, and effective CSS comment-preservation options.
  Dependency upgrades MUST update the maintained identities; transform changes
  MUST bump the wrapper revision. Selection gates (fast, enabled, exclusions)
  are not transform inputs.
- Hash current eligible bytes once in the worker. Reuse ONLY when the digest
  matches recorded `source_hash` (the exact transform input, not original author
  source) AND the recipe matches. Restore the verified result atomically.
  Matching a prior output is not sufficient: the CSS wrapper can be
  non-idempotent for important comments matching preservation patterns.
  Changed recipes ALWAYS transform current stage-input bytes, and diagnostics
  must describe that action without claiming the digest proves source origin.
  Output hashes and historical observed-input hashes MUST NOT select an older
  source snapshot. Stripped comments require regeneration from authoritative
  source, not recovery from this cache.
- Every successful fresh transform atomically replaces the target with mode
  `0644`, even when bytes are identical, preserving uncached publication behavior
  and existing hardlinked releases. The existing empty-JavaScript exception
  remains a no-write operation. Verified cache hits publish by the same atomic
  replacement, even if result bytes already match. Never chmod an inode
  potentially shared with an existing release.
- Source hashes are optional: absent means unavailable, not empty. A missing or
  corrupt result blob needed for exact-input reuse triggers a warned transform
  of current bytes and repair. Historical source snapshots may remain stored,
  but are never read or newly written. New records describe the exact
  current input and its result, not a claim that this input is original source.
- Persist the verified result blob and the per-asset transform relation
  atomically BEFORE replacing the generated target atomically. Failed target
  writes remain retryable using the current-byte check. Never truncate hardlinked
  targets or blobs. On persistence failure, warn and durably invalidate any old
  record before uncached target mutation. If persistence and invalidation both
  fail, leave the target unchanged and report failure; many-to-one outputs must
  never allow obsolete source provenance to be resurrected.
- Historical records may contain an observed-input digest (`input_hash`).
  Ignore old values; never produce new ones or use them to authorize reuse.
  Failed target replacement is retryable only through the exact current-input
  and same-recipe relation; a changed recipe transforms the current bytes.
- Legacy bare hash maps are not promoted or trusted as source/recipe provenance.
  Do not chase or delete legacy sidecars. A source-regenerating rebuild establishes
  clean provenance for assets not normally recopied.
- Retain warning-only transform failure policy. Report transformed (engine
  calls), restored (exact-input result reuse), excluded, and failed counts
  separately; no unchanged/no-write counter is promised. Raw recopies can avoid
  engine work after priming, but retained non-idempotent output may transform
  every time. StaticAssets Cleanup must
  continue creating hashed aliases from the final minified canonical bytes.

---

## Reading-Time Field Ownership

The `reading_time` and `stats` plugins both participate in the Transform stage
and share one pure word-count, reading-time, and formatted-text calculator.
They retain separate configuration and lifecycle contracts. `stats` runs at
`PriorityEarly`; `reading_time` uses the default priority. Early content
transforms such as `jinja_md` run before `reading_time`, so public reading-time
fields count expanded content.

### `reading_time`

For each non-skipped post with non-empty content, `reading_time` writes these
public fields to `Post.Extra`:

- `word_count`
- `reading_time`
- `reading_time_text`

The calculator excludes fenced code, inline code, image definitions, URLs, HTML tags, and Markdown formatting characters.
It counts Unicode letters and digits, uses a default of 200 words per minute, rounds up, and returns zero for zero words.
The plugin reads its words-per-minute value from `config.Extra["words_per_minute"]`.

### `stats`

During Transform, `stats` calculates and stores an independent `PostStats`
value with the shared calculator and its Stats configuration. It also writes
the reading-time fields so the plugin works when used alone. In a normal
default build, `reading_time` runs later and overwrites only the public
reading-time fields; it does not change `PostStats`.

`include_code_in_count` applies to Stats' `PostStats` and to Stats feed/site
aggregates in both standalone and normal default builds. It does not override
the public fields written by `reading_time`.
`track_code_blocks` independently controls Stats' code metrics. When it is
false, `PostStats.CodeBlocks` and `PostStats.CodeLines`, and their feed/site
aggregate totals, are zero.
The plugin reads its Stats-specific settings from `config.Extra["stats"]` and
keeps the top-level words-per-minute fallback.
When both words-per-minute settings exist, the top-level value takes precedence for backward compatibility.

Stats feed and site aggregates are derived from each post's `PostStats`, not
from the public fields that `reading_time` may write later.

This contract keeps one reading-time algorithm while preserving standalone use by callers such as list-cache transforms.

---

## Reading-Time Transform Order

The default plugin registry lists `reading_time` before `stats`, but lifecycle
priority sorting determines the effective Transform order. For these plugins,
`stats` runs first at `PriorityEarly`, `jinja_md` runs next at
`PriorityEarly`, and `reading_time` runs later at the default priority. This
keeps Stats' independent `PostStats` based on pre-Jinja content while public
ReadingTime fields use expanded content.

---

## Plugin Load Order

When `hooks = ["default"]`, plugins load in this order:

```python
DEFAULT_PLUGINS = [
    "config_defaults",      # Set defaults first
    "glob",                 # Find files
    "load",                 # Parse files
    "auto_description",     # Generate descriptions
    "jinja_md",             # Process jinja in markdown
    "tailwind",             # Build Tailwind CSS and inject includes
    "render_markdown",      # Markdown → HTML
    "wikilinks",            # Resolve internal links
    "heading_anchors",      # Add heading anchors
    "toc",                  # Generate TOC
    "link_collector",       # Track inlinks/outlinks
    "series",               # Auto-generate series feeds (PriorityEarly)
    "feeds",                # Create feed collections
    "overwrite_check",      # Detect conflicting output paths
    "prevnext",             # Calculate prev/next from feeds/series
    "publish_feeds",        # Write feed files
    "publish_html",         # Write post files
    "copy_assets",          # Copy static files
    "redirects",            # Generate redirect pages
    "css_minify",           # Minify CSS files (PriorityLast)
    "js_minify",            # Minify JS files (PriorityLast)
]
```

---

## Disabling Plugins

```toml
[name]
hooks = ["default"]
disabled_hooks = [
    "toc",              # Don't generate TOC
    "wikilinks",        # Don't process wikilinks
    "link_collector",   # Don't track inlinks/outlinks
]
```

Or load specific plugins only:

```toml
[name]
hooks = [
    "glob",
    "load",
    "render_markdown",
    "publish_html",
]
```

---

## See Also

- [FEEDS.md](./FEEDS.md) - Feed system specification
- [CONFIG.md](./CONFIG.md) - Plugin configuration
- [SPEC.md](./SPEC.md) - Core specification
- [LIFECYCLE.md](./LIFECYCLE.md) - Lifecycle stages


## Automatically Generated `llms.txt`

The default full-site build runs `llms_txt` in the Write stage at last
priority, after static assets, published posts, and sitemap output. The
plugin creates an H1, a site-description blockquote, and Markdown link-list
sections for the homepage, documentation, and recent published content.
It favors `.md` URLs when that format is enabled and falls back to HTML
links. Only public, published, non-skipped and indexable posts are eligible.

Site config `[markata-go.llms_txt]` accepts `enabled` (default true) and
`max_pages` (default 30; cap 100). An explicitly authored `llms.md`, a
project static `llms.txt`, or a hand-managed output file takes precedence.
Only auto-generated output is removed when disabled. This plugin is not
part of single-page builds.
