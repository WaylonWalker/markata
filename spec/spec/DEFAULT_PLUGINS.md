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
│                                                                      │
│  RENDER                                                              │
│    ├─ render_markdown    Convert markdown to HTML                    │
│    └─ wikilinks          Resolve [[internal links]]                  │
│                                                                      │
│  POST-RENDER                                                         │
│    ├─ heading_anchors    Add anchor links to headings               │
│    ├─ toc                Generate table of contents                  │
│    ├─ link_collector     Track inlinks/outlinks between posts       │
│    ├─ feeds              Generate feed collections                   │
│    └─ prevnext           Calculate prev/next from feeds/series       │
│                                                                      │
│  OUTPUT                                                              │
│    ├─ publish_feeds      Write HTML/RSS/Atom/JSON/MD/TXT/Sitemap    │
│    ├─ publish_html       Write individual post HTML files           │
│    └─ copy_assets        Copy static files                           │
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

[name.glob]
glob_patterns = ["**/*.md"]
use_gitignore = true
exclude_patterns = ["node_modules/**", ".git/**", "output/**"]

[name.feeds]
default_items_per_page = 10
default_orphan_threshold = 3

[name.feeds.default_formats]
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
[name.glob]
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

**Hook signature:**
```python
@hook_impl
def glob(core):
    for pattern in core.config.glob.glob_patterns:
        for path in Path().glob(pattern):
            if should_include(path, core.config.glob):
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

**Hook signature:**
```python
@hook_impl
def load(core):
    for path in core.files:
        try:
            raw = path.read_text(encoding='utf-8')
            frontmatter, content = parse_frontmatter(raw)
            post = core.Post(
                path=path,
                content=content,
                **frontmatter
            )
            core.posts.append(post)
        except Exception as e:
            logger.warning(f"Failed to load {path}: {e}")
```

---

## Pre-Render Processing

### `auto_description`

**Stage:** `pre_render`

**Purpose:** Generate descriptions from content for posts that don't have one.

**Configuration:**
```toml
[name.auto_description]
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

**Hook signature:**
```python
@hook_impl
def pre_render(core):
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
[name.jinja_md]
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

**Hook signature:**
```python
@hook_impl
def pre_render(core):
    config = core.config.jinja_md
    if not config.enabled:
        return
    
    filter_expr = "jinja == True" if not config.default_enabled else "jinja != False"
    
    for post in core.filter(filter_expr):
        template = core.jinja_env.from_string(post.content)
        post.content = template.render(
            post=post,
            core=core,
            config=core.config,
            today=date.today(),
            now=datetime.now()
        )
```

---

### `prevnext`

**Stage:** `post_render` (after feeds are created)

**Purpose:** Calculate previous/next post links for navigation based on feeds or series.

**Configuration:**
```toml
[name.prevnext]
enabled = true
strategy = "first_feed"        # Resolution strategy (see below)
default_feed = "blog"          # Feed to use when strategy = "explicit_feed"
```

**Strategy Options:**

| Strategy | Description |
|----------|-------------|
| `first_feed` | Use first feed the post appears in (default) |
| `explicit_feed` | Always use `default_feed` for all posts |
| `series` | Use post's `series` frontmatter, fall back to `first_feed` |
| `frontmatter` | Use post's `prevnext_feed` frontmatter, fall back to `first_feed` |

**Post frontmatter options:**
```yaml
---
title: My Post
series: python-tutorial       # Series name (acts like feed slug)
prevnext_feed: blog           # Explicit feed for this post's navigation
---
```

**Series Definition:**

Series are lightweight feed aliases for navigation purposes. Define them inline or reference existing feeds:

```toml
# Inline series definition
[[name.series]]
slug = "python-tutorial"
title = "Python Tutorial Series"
filter = "series == 'python-tutorial'"
sort = "date"
reverse = false               # Chronological for tutorials

# Series as feed alias
[[name.series]]
slug = "getting-started"
feed = "tutorials"            # Use existing feed's posts and sort order
```

**Post fields added:**
| Field | Type | Description |
|-------|------|-------------|
| `prev` | Post? | Previous post in sequence |
| `next` | Post? | Next post in sequence |
| `prevnext_feed` | string? | Feed/series slug used for navigation |
| `prevnext_context` | PrevNextContext | Full navigation context |

**PrevNextContext structure:**
```python
class PrevNextContext:
    feed_slug: str              # Feed or series slug
    feed_title: str             # Feed or series title
    position: int               # Position in sequence (1-indexed)
    total: int                  # Total posts in sequence
    prev: Post | None           # Previous post
    next: Post | None           # Next post
```

**Behavior:**

1. **Collect all series definitions** - Parse `[[name.series]]` configs
2. **For each post, determine navigation context:**
   - If `strategy == "frontmatter"` and post has `prevnext_feed`: use that feed
   - If `strategy == "series"` and post has `series`: look up series, use its posts
   - If `strategy == "explicit_feed"`: use `default_feed`
   - If `strategy == "first_feed"`: find first feed containing this post
3. **Get ordered posts from feed/series**
4. **Find post position and set prev/next**

**Resolution order for series/feeds:**
1. Check `core.series` for matching slug
2. Check `core.feeds` for matching slug
3. Log warning if not found, skip prev/next for this post

**Hook signature:**
```python
@hook_impl
def post_render(core):
    config = core.config.prevnext
    if not config.enabled:
        return
    
    # Build series lookup
    series_map = {}
    for series_config in core.config.get("series", []):
        if "feed" in series_config:
            # Series is alias for existing feed
            feed = core.feeds.get(series_config["feed"])
            if feed:
                series_map[series_config["slug"]] = {
                    "title": series_config.get("title", feed.title),
                    "posts": feed.posts
                }
        else:
            # Inline series definition
            posts = core.filter(series_config.get("filter", "True"))
            sort_field = series_config.get("sort", "date")
            reverse = series_config.get("reverse", True)
            posts = sorted(posts, key=lambda p: getattr(p, sort_field, ""), reverse=reverse)
            series_map[series_config["slug"]] = {
                "title": series_config.get("title", series_config["slug"]),
                "posts": posts
            }
    
    # Build feed lookup
    feed_map = {f.slug: f for f in core.feeds}
    
    for post in core.posts:
        # Determine which feed/series to use
        context_slug = None
        
        if config.strategy == "frontmatter" and hasattr(post, "prevnext_feed"):
            context_slug = post.prevnext_feed
        elif config.strategy == "series" and hasattr(post, "series"):
            context_slug = post.series
        elif config.strategy == "explicit_feed":
            context_slug = config.default_feed
        elif config.strategy == "first_feed":
            # Find first feed containing this post
            for feed in core.feeds:
                if post in feed.posts:
                    context_slug = feed.slug
                    break
        
        if not context_slug:
            post.prev = None
            post.next = None
            continue
        
        # Get posts from series or feed
        if context_slug in series_map:
            context = series_map[context_slug]
            posts = context["posts"]
            title = context["title"]
        elif context_slug in feed_map:
            feed = feed_map[context_slug]
            posts = feed.posts
            title = feed.title
        else:
            logger.warning(f"Feed/series '{context_slug}' not found for post '{post.slug}'")
            post.prev = None
            post.next = None
            continue
        
        # Find position and set prev/next
        try:
            idx = posts.index(post)
            post.prev = posts[idx - 1] if idx > 0 else None
            post.next = posts[idx + 1] if idx < len(posts) - 1 else None
            post.prevnext_feed = context_slug
            post.prevnext_context = PrevNextContext(
                feed_slug=context_slug,
                feed_title=title,
                position=idx + 1,
                total=len(posts),
                prev=post.prev,
                next=post.next
            )
        except ValueError:
            post.prev = None
            post.next = None
```

**Template usage:**
```jinja2
{# Basic prev/next navigation #}
<nav class="post-navigation">
  {% if post.prev %}
  <a href="{{ post.prev.href }}" class="nav-prev">
    <span class="nav-label">Previous</span>
    <span class="nav-title">{{ post.prev.title }}</span>
  </a>
  {% endif %}
  
  {% if post.next %}
  <a href="{{ post.next.href }}" class="nav-next">
    <span class="nav-label">Next</span>
    <span class="nav-title">{{ post.next.title }}</span>
  </a>
  {% endif %}
</nav>

{# Series navigation with position indicator #}
{% if post.prevnext_context %}
<nav class="series-navigation">
  <div class="series-info">
    <span class="series-title">{{ post.prevnext_context.feed_title }}</span>
    <span class="series-position">
      Part {{ post.prevnext_context.position }} of {{ post.prevnext_context.total }}
    </span>
  </div>
  
  <div class="series-links">
    {% if post.prev %}
    <a href="{{ post.prev.href }}">← {{ post.prev.title }}</a>
    {% endif %}
    {% if post.next %}
    <a href="{{ post.next.href }}">{{ post.next.title }} →</a>
    {% endif %}
  </div>
</nav>
{% endif %}
```

**Example configurations:**

```toml
# Strategy 1: Use first feed (default)
[name.prevnext]
enabled = true
strategy = "first_feed"

# Strategy 2: Always use specific feed
[name.prevnext]
enabled = true
strategy = "explicit_feed"
default_feed = "all-posts"

# Strategy 3: Series-based navigation
[name.prevnext]
enabled = true
strategy = "series"

# Define series for tutorials
[[name.series]]
slug = "python-basics"
title = "Python Basics Tutorial"
filter = "series == 'python-basics'"
sort = "date"
reverse = false

[[name.series]]
slug = "web-dev-guide"
title = "Web Development Guide"
filter = "'web-dev-guide' in tags"
sort = "order"              # Custom sort field
reverse = false

# Series as alias for existing feed
[[name.series]]
slug = "recent"
feed = "blog"               # Inherit posts from 'blog' feed
title = "Recent Posts"
```

**Frontmatter examples:**

```yaml
---
# Post in a tutorial series
title: "Variables in Python"
series: python-basics
---

---
# Post with explicit navigation feed
title: "Announcement Post"
prevnext_feed: announcements
---

---
# Post with both (series takes precedence with strategy="series")
title: "Advanced Topic"
series: python-advanced
prevnext_feed: all-posts
---
```

---

## Render

### `render_markdown`

**Stage:** `render`

**Purpose:** Convert markdown content to HTML.

**Configuration:**
```toml
[name.markdown]
backend = "auto"               # "markdown-it", "commonmark", etc.

[name.markdown.extensions]
tables = true
admonitions = true
footnotes = true
strikethrough = true
task_lists = true
heading_ids = true

[name.markdown.highlight]
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

**Hook signature:**
```python
@hook_impl
def render(core):
    for post in core.filter("not skip"):
        post.article_html = core.markdown_parser.render(post.content)
```

---

### `wikilinks`

**Stage:** `render` (after `render_markdown`)

**Purpose:** Resolve `[[internal links]]` to actual post URLs.

**Configuration:**
```toml
[name.wikilinks]
enabled = true
warn_broken = true             # Warn about broken links
broken_class = "broken-link"   # CSS class for broken links
```

**Syntax:**
```markdown
Link to post: [[other-post-slug]]
With custom text: [[other-post-slug|Click here]]
```

**Resolution:**
1. Find post where `slug == link_target`
2. If found: `<a href="{post.href}">{text or post.title}</a>`
3. If not found: Leave as `[[link]]` or wrap with broken-link class

**Hook signature:**
```python
@hook_impl(trylast=True)  # Run after render_markdown
def render(core):
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

**Stage:** `post_render`

**Purpose:** Add anchor links to headings for direct linking.

**Configuration:**
```toml
[name.heading_anchors]
enabled = true
min_level = 2                  # Start at h2
max_level = 4                  # End at h4
position = "end"               # "start" or "end"
symbol = "#"                   # Link text
class = "heading-anchor"       # CSS class
```

**Behavior:**
1. Find all headings in `article_html`
2. Generate ID from heading text (slugified)
3. Handle duplicate IDs by appending numbers
4. Insert anchor link at configured position

**Example output:**
```html
<h2 id="my-section">
  My Section
  <a href="#my-section" class="heading-anchor">#</a>
</h2>
```

**Hook signature:**
```python
@hook_impl
def post_render(core):
    config = core.config.heading_anchors
    if not config.enabled:
        return
    
    for post in core.filter("not skip"):
        post.article_html = add_heading_anchors(post.article_html, config)
```

---

### `toc`

**Stage:** `post_render`

**Purpose:** Generate table of contents from headings.

**Configuration:**
```toml
[name.toc]
enabled = true
min_level = 2
max_level = 4
min_headings = 3               # Only generate if >= 3 headings
```

**Post fields added:**
| Field | Type | Description |
|-------|------|-------------|
| `toc` | string | HTML table of contents |
| `toc_items` | List[TocItem] | Structured TOC data |

**TocItem structure:**
```python
class TocItem:
    level: int      # Heading level (2-6)
    id: str         # Heading ID
    text: str       # Heading text
    children: List[TocItem]
```

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
[name.link_collector]
enabled = true
include_feeds = false          # Exclude feed pages from inlinks by default
include_index = false          # Exclude index page from inlinks by default
```

**Post fields added:**
| Field | Type | Description |
|-------|------|-------------|
| `hrefs` | List[str] | Raw href values from all links in post |
| `inlinks` | List[Link] | Links pointing TO this post from other posts |
| `outlinks` | List[Link] | Links FROM this post to other pages |

**Link structure:**
```python
class Link:
    source_url: str            # Absolute URL of the source post
    source_post: Post          # Reference to source post object
    target_post: Post | None   # Reference to target post (None if external)
    raw_target: str            # Original href value
    target_url: str            # Resolved absolute URL
    target_domain: str | None  # Domain extracted from target_url
    is_internal: bool          # True if link points to same site
    is_self: bool              # True if link points to same post
    source_text: str           # Cleaned link text from source
    target_text: str           # Cleaned link text from target
```

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
Results are cached per-post based on hash of:
- Plugin file
- Post slug
- Post title  
- Post content

**Hook signature:**
```python
@hook_impl
@register_attr("links")
def render(core):
    links = []
    site_domain = urlparse(str(core.config.url)).netloc
    
    # Build slug -> post lookup
    post_by_slug = {p.slug: p for p in core.posts}
    
    for post in core.posts:
        base_url = urljoin(str(core.config.url), post.slug)
        soup = BeautifulSoup(post.article_html, "html.parser")
        
        # Optionally limit to post-body section
        post_body = soup.find(id="post-body")
        if post_body:
            soup = post_body
        
        post.hrefs = [a["href"] for a in soup.find_all("a", href=True)]
        
        for href in post.hrefs:
            target_url = urljoin(base_url, href)
            domain = urlparse(target_url).netloc
            is_internal = domain == site_domain
            
            target_post = None
            if is_internal:
                target_slug = urlparse(target_url).path.strip("/")
                target_post = post_by_slug.get(target_slug)
            
            links.append(Link(
                source_url=base_url,
                source_post=post,
                target_post=target_post,
                raw_target=href,
                target_url=target_url,
                target_domain=domain,
                is_internal=is_internal,
                is_self=target_post and post.slug == target_post.slug
            ))
    
    core.links = links
    
    # Assign inlinks/outlinks to each post
    for post in core.posts:
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
```python
from collections import Counter

def count_links(links: List[Link]) -> Counter:
    """Count target_url frequency across all links."""
    return Counter(link.target_url for link in links)

def count_domains(links: List[Link]) -> Counter:
    """Count target_domain frequency across all links."""
    return Counter(link.target_domain for link in links if link.target_domain)

def get_external_links(links: List[Link]) -> List[Link]:
    """Filter to only external links."""
    return [link for link in links if not link.is_internal]

def get_broken_links(links: List[Link]) -> List[Link]:
    """Find internal links that don't resolve to a post."""
    return [link for link in links if link.is_internal and link.target_post is None]
```

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

**Hook signature:**
```python
@hook_impl
def post_render(core):
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

**Output paths:**
| Format | Path Pattern |
|--------|--------------|
| HTML | `/{slug}/index.html`, `/{slug}/page/{n}/index.html` |
| RSS | `/{slug}/rss.xml` |
| Atom | `/{slug}/atom.xml` |
| JSON | `/{slug}/feed.json` |
| Markdown | `/{slug}/index.md` |
| Text | `/{slug}/index.txt` |
| Sitemap | `/{slug}/sitemap.xml` |

**Hook signature:**
```python
@hook_impl
def save(core):
    for feed in core.feeds:
        for format_name, enabled in feed.formats.items():
            if not enabled:
                continue
            
            template = get_feed_template(format_name, feed, core)
            
            if format_name == 'html':
                # Paginated output
                for page_num, page_posts in enumerate(feed.pages, 1):
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

**Hook signature:**
```python
@hook_impl
def save(core):
    for post in core.filter("not skip"):
        template_name = post.template or "post.html"
        template = core.jinja_env.get_template(template_name)
        
        html = template.render(
            post=post,
            body=post.article_html,
            config=core.config,
            core=core
        )
        
        output_path = core.config.output_dir / post.slug / "index.html"
        output_path.parent.mkdir(parents=True, exist_ok=True)
        output_path.write_text(html)
```

---

### `copy_assets`

**Stage:** `save`

**Purpose:** Copy static assets to output directory.

**Configuration:**
```toml
[name.assets]
dir = "static"
output_subdir = ""             # "" = root, "assets" = /assets/
exclude = ["*.psd", "*.ai"]
fingerprint = false

[name.assets.fingerprint]
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

**Hook signature:**
```python
@hook_impl
def save(core):
    config = core.config.assets
    assets_dir = Path(config.dir)
    
    if not assets_dir.exists():
        return
    
    for src in assets_dir.rglob("*"):
        if src.is_file() and not is_excluded(src, config.exclude):
            rel_path = src.relative_to(assets_dir)
            
            if should_fingerprint(src, config):
                rel_path = fingerprint_path(src, rel_path, config)
            
            dst = core.config.output_dir / config.output_subdir / rel_path
            dst.parent.mkdir(parents=True, exist_ok=True)
            shutil.copy2(src, dst)
```

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
    "render_markdown",      # Markdown → HTML
    "wikilinks",            # Resolve internal links
    "heading_anchors",      # Add heading anchors
    "toc",                  # Generate TOC
    "link_collector",       # Track inlinks/outlinks
    "feeds",                # Create feed collections
    "prevnext",             # Calculate prev/next from feeds/series
    "publish_feeds",        # Write feed files
    "publish_html",         # Write post files
    "copy_assets",          # Copy static files
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
