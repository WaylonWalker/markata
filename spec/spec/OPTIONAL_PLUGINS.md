# Optional Plugins Specification

This document specifies optional plugins that extend the static site generator with additional features. Unlike default plugins, these must be explicitly enabled in the configuration.

## Plugin Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                       OPTIONAL PLUGIN SET                           │
├─────────────────────────────────────────────────────────────────────┤
│  CONTENT ENHANCEMENT                                                │
│    ├─ glossary          Auto-link terms to definition posts        │
│    ├─ mermaid           Render Mermaid diagrams                    │
│    ├─ chartjs           Render Chart.js charts                     │
│    ├─ csv_fence         Convert CSV code blocks to tables          │
│    └─ md_video          Convert image syntax to video tags         │
│                                                                      │
│  LINK ENHANCEMENT                                                    │
│    ├─ one_line_link     Rich previews for URLs on own line         │
│    └─ wikilink_hover    Hover previews for wikilinks               │
│                                                                      │
│  OUTPUT GENERATION                                                   │
│    └─ qrcode            Generate QR codes for posts                │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Enabling Optional Plugins

```toml
[your-ssg]
hooks = [
    "default",
    "glossary",
    "mermaid",
    "chartjs",
]
```

Or with full module path:

```toml
[your-ssg]
hooks = [
    "default",
    "your_ssg.plugins.glossary",
    "your_ssg.plugins.mermaid",
]
```

---

## Content Enhancement Plugins

### `glossary`

**Stage:** `post_render`

**Purpose:** Automatically link glossary terms in post content to their definition pages.

**Dependencies:** None (uses standard library)

**Configuration:**

```toml
[your-ssg.glossary]
enabled = true
link_class = "glossary-term"       # CSS class for glossary links
case_sensitive = false             # Match terms case-insensitively
tooltip = true                     # Add tooltip with description
max_links_per_term = 1             # Link only first occurrence (0 = all)
exclude_tags = ["glossary"]        # Don't link in glossary posts themselves
export_json = true                 # Export glossary.json to output
```

**Post Frontmatter (Definition Posts):**

```yaml
---
title: API
templateKey: glossary
description: Application Programming Interface - a set of protocols...
aliases:
  - APIs
  - Application Programming Interface
---
```

**Behavior:**

1. Scan posts with `templateKey: glossary` or in configured glossary path
2. Build term → definition lookup including aliases
3. For each non-glossary post, find term occurrences in `article_html`
4. Replace with linked version: `<a href="/glossary/api/" class="glossary-term" title="...">API</a>`
5. Optionally export `glossary.json` with all terms

**Output Files:**

| File | Description |
|------|-------------|
| `{output}/glossary.json` | JSON export of all terms |

**Example glossary.json:**

```json
{
  "terms": [
    {
      "term": "API",
      "slug": "api",
      "description": "Application Programming Interface",
      "aliases": ["APIs"],
      "href": "/glossary/api/"
    }
  ]
}
```

**Hook Signature:**

```python
@hook_impl
def post_render(core):
    config = core.config.glossary
    if not config.enabled:
        return
    
    # Build glossary from definition posts
    glossary = build_glossary(core)
    
    # Process each post
    for post in core.filter("not skip and templateKey != 'glossary'"):
        post.article_html = link_glossary_terms(
            post.article_html,
            glossary,
            config
        )
```

---

### `mermaid`

**Stage:** `post_render`

**Purpose:** Convert Mermaid code blocks into rendered diagrams.

**Dependencies:** None (client-side rendering via CDN)

**Configuration:**

```toml
[your-ssg.mermaid]
enabled = true
cdn_url = "https://cdn.jsdelivr.net/npm/mermaid@10/dist/mermaid.esm.min.mjs"
theme = "default"                  # default, dark, forest, neutral
init_config = {}                   # Additional mermaid.initialize() options
```

**Syntax:**

````markdown
```mermaid
graph TD
    A[Start] --> B{Decision}
    B -->|Yes| C[Action 1]
    B -->|No| D[Action 2]
```
````

**Behavior:**

1. Find all `<pre><code class="language-mermaid">` blocks
2. Replace with `<pre class="mermaid">{diagram code}</pre>`
3. Inject Mermaid.js script (once per page with mermaid content)

**Output:**

```html
<pre class="mermaid">
graph TD
    A[Start] --> B{Decision}
    B -->|Yes| C[Action 1]
    B -->|No| D[Action 2]
</pre>

<script type="module">
  import mermaid from 'https://cdn.jsdelivr.net/npm/mermaid@10/dist/mermaid.esm.min.mjs';
  mermaid.initialize({ startOnLoad: true, theme: 'default' });
</script>
```

**Supported Diagram Types:**

| Type | Description |
|------|-------------|
| `graph` / `flowchart` | Flow diagrams |
| `sequenceDiagram` | Sequence diagrams |
| `classDiagram` | Class diagrams |
| `stateDiagram` | State diagrams |
| `erDiagram` | Entity relationship |
| `gantt` | Gantt charts |
| `pie` | Pie charts |
| `gitGraph` | Git graphs |
| `mindmap` | Mind maps |
| `timeline` | Timelines |

**Hook Signature:**

```python
@hook_impl
def post_render(core):
    config = core.config.mermaid
    if not config.enabled:
        return
    
    for post in core.filter("not skip"):
        if has_mermaid_blocks(post.article_html):
            post.article_html = process_mermaid_blocks(post.article_html, config)
            post.article_html = inject_mermaid_script(post.article_html, config)
```

---

### `chartjs`

**Stage:** `post_render`

**Purpose:** Convert Chart.js JSON blocks into rendered charts.

**Dependencies:** None (client-side rendering via CDN)

**Configuration:**

```toml
[your-ssg.chartjs]
enabled = true
cdn_url = "https://cdn.jsdelivr.net/npm/chart.js"
default_options = {}               # Default Chart.js options
```

**Syntax:**

````markdown
```chartjs
{
  "type": "bar",
  "data": {
    "labels": ["Red", "Blue", "Yellow"],
    "datasets": [{
      "label": "Votes",
      "data": [12, 19, 3],
      "backgroundColor": ["#f87171", "#60a5fa", "#fbbf24"]
    }]
  }
}
```
````

**Behavior:**

1. Find all `<pre><code class="language-chartjs">` blocks
2. Parse JSON content
3. Replace with canvas element and initialization script
4. Inject Chart.js script (once per page)

**Output:**

```html
<div class="chartjs-container">
  <canvas id="chart-1"></canvas>
</div>

<script>
  new Chart(document.getElementById('chart-1'), {
    type: 'bar',
    data: { ... }
  });
</script>

<script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
```

**Supported Chart Types:**

| Type | Description |
|------|-------------|
| `bar` | Bar chart |
| `line` | Line chart |
| `pie` | Pie chart |
| `doughnut` | Doughnut chart |
| `radar` | Radar chart |
| `polarArea` | Polar area chart |
| `bubble` | Bubble chart |
| `scatter` | Scatter plot |

**Hook Signature:**

```python
@hook_impl
def post_render(core):
    config = core.config.chartjs
    if not config.enabled:
        return
    
    for post in core.filter("not skip"):
        if has_chartjs_blocks(post.article_html):
            post.article_html = process_chartjs_blocks(post.article_html, config)
```

---

### `csv_fence`

**Stage:** `post_render`

**Purpose:** Convert CSV code blocks into HTML tables.

**Dependencies:** None (uses standard library csv module)

**Configuration:**

```toml
[your-ssg.csv_fence]
enabled = true
table_class = "csv-table"          # CSS class for table
has_header = true                  # First row is header
delimiter = ","                    # CSV delimiter
```

**Syntax:**

````markdown
```csv
Name,Age,City
Alice,30,New York
Bob,25,Los Angeles
Charlie,35,Chicago
```
````

**Behavior:**

1. Find all `<pre><code class="language-csv">` blocks
2. Parse CSV content
3. Generate HTML table with `<thead>` (if header) and `<tbody>`

**Output:**

```html
<table class="csv-table">
  <thead>
    <tr>
      <th>Name</th>
      <th>Age</th>
      <th>City</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>Alice</td>
      <td>30</td>
      <td>New York</td>
    </tr>
    <tr>
      <td>Bob</td>
      <td>25</td>
      <td>Los Angeles</td>
    </tr>
    <tr>
      <td>Charlie</td>
      <td>35</td>
      <td>Chicago</td>
    </tr>
  </tbody>
</table>
```

**Advanced Usage:**

````markdown
```csv delimiter=";" has_header="false"
1;2;3
4;5;6
```
````

**Hook Signature:**

```python
@hook_impl
def post_render(core):
    config = core.config.csv_fence
    if not config.enabled:
        return
    
    for post in core.filter("not skip"):
        post.article_html = process_csv_blocks(post.article_html, config)
```

---

### `md_video`

**Stage:** `post_render`

**Purpose:** Convert markdown image syntax for video files into HTML video elements.

**Dependencies:** None

**Configuration:**

```toml
[your-ssg.md_video]
enabled = true
video_extensions = [".mp4", ".webm", ".ogg", ".mov"]
video_class = "md-video"
controls = true
autoplay = false
loop = false
muted = false
preload = "metadata"               # "none", "metadata", "auto"
```

**Syntax:**

```markdown
![Video description](video.mp4)

![Animated demo](demo.webm)
```

**Behavior:**

1. Find all `<img>` tags where `src` ends with video extension
2. Replace with `<video>` element with configured attributes
3. Preserve alt text as fallback content

**Output:**

```html
<video class="md-video" controls preload="metadata">
  <source src="video.mp4" type="video/mp4">
  Video description
</video>
```

**GIF-like Behavior:**

For silent, looping videos (like animated GIFs):

```toml
[your-ssg.md_video]
autoplay = true
loop = true
muted = true
controls = false
```

**Hook Signature:**

```python
@hook_impl
def post_render(core):
    config = core.config.md_video
    if not config.enabled:
        return
    
    for post in core.filter("not skip"):
        post.article_html = convert_video_images(post.article_html, config)
```

---

## Link Enhancement Plugins

### `one_line_link`

**Stage:** `post_render`

**Purpose:** Expand URLs that appear alone on a line into rich preview cards.

**Dependencies:** `lxml` (optional, for parsing fetched pages)

**Configuration:**

```toml
[your-ssg.one_line_link]
enabled = true
card_class = "link-card"
fetch_metadata = true              # Fetch title/description from URL
cache_metadata = true              # Cache fetched metadata
fallback_title = "Link"            # Title when fetch fails
timeout = 5                        # Fetch timeout in seconds

# URL patterns to exclude from expansion
exclude_patterns = [
    "^https://twitter.com",
    "^https://x.com",
]

# Custom templates per domain
[your-ssg.one_line_link.templates]
"github.com" = "github-card.html"
"youtube.com" = "youtube-card.html"
```

**Syntax:**

```markdown
Check out this article:

https://example.com/awesome-article

And continue reading...
```

**Behavior:**

1. Find URLs that are alone on their own line (paragraph containing only a URL)
2. Fetch page metadata (title, description, image) if enabled
3. Replace with rich card HTML

**Output:**

```html
<p>Check out this article:</p>

<a href="https://example.com/awesome-article" class="link-card">
  <div class="link-card-image" style="background-image: url('...')"></div>
  <div class="link-card-content">
    <div class="link-card-title">Awesome Article Title</div>
    <div class="link-card-description">A brief description...</div>
    <div class="link-card-url">example.com</div>
  </div>
</a>

<p>And continue reading...</p>
```

**Detection Rules:**

| Pattern | Expanded? |
|---------|-----------|
| `https://example.com` (alone) | Yes |
| `Check out https://example.com` | No (inline) |
| `[Link](https://example.com)` | No (already link) |
| `<https://example.com>` | No (autolink) |

**Hook Signature:**

```python
@hook_impl
def post_render(core):
    config = core.config.one_line_link
    if not config.enabled:
        return
    
    for post in core.filter("not skip"):
        post.article_html = expand_one_line_links(
            post.article_html,
            config,
            core.cache
        )
```

---

### `wikilink_hover`

**Stage:** `post_render` (after wikilinks)

**Purpose:** Add hover previews to wikilinks showing target post content.

**Dependencies:** `lxml` (for HTML parsing)

**Configuration:**

```toml
[your-ssg.wikilink_hover]
enabled = true
preview_length = 200               # Characters to show in preview
include_image = true               # Include featured image if available
screenshot_service = ""            # URL of screenshot service (optional)
screenshot_width = 400
screenshot_height = 300
```

**Behavior:**

1. Find all wikilink anchors (`<a>` tags created by wikilinks plugin)
2. Add `data-preview` attribute with preview content
3. Optionally add screenshot URL

**Output:**

```html
<a href="/other-post/"
   class="wikilink"
   data-preview="Preview text from the target post..."
   data-preview-image="/other-post/featured.jpg">
  Other Post
</a>
```

**JavaScript Integration:**

The plugin adds data attributes; JavaScript handles the hover display:

```javascript
document.querySelectorAll('[data-preview]').forEach(link => {
  link.addEventListener('mouseenter', showPreview);
  link.addEventListener('mouseleave', hidePreview);
});
```

**Screenshot Service:**

If `screenshot_service` is configured:

```toml
screenshot_service = "https://screenshot.example.com/capture?url="
```

Output:

```html
<a href="/other-post/"
   data-preview-screenshot="https://screenshot.example.com/capture?url=https://mysite.com/other-post/">
  Other Post
</a>
```

**Hook Signature:**

```python
@hook_impl(trylast=True)  # Run after wikilinks
def post_render(core):
    config = core.config.wikilink_hover
    if not config.enabled:
        return
    
    # Build post lookup
    post_lookup = {p.slug: p for p in core.posts}
    
    for post in core.filter("not skip"):
        post.article_html = add_wikilink_previews(
            post.article_html,
            post_lookup,
            config,
            core.config.url
        )
```

---

## Output Generation Plugins

### `qrcode`

**Stage:** `save`

**Purpose:** Generate QR code images for each post's URL.

**Dependencies:** `pyqrcode` (or `qrcode`)

**Installation:**

```bash
pip install pyqrcode
# or
pip install qrcode[pil]
```

**Configuration:**

```toml
[your-ssg.qrcode]
enabled = true
format = "svg"                     # "svg" or "png"
size = 200                         # Size in pixels
output_dir = "qrcodes"             # Subdirectory in output
filename_template = "{slug}.{format}"
error_correction = "M"             # L, M, Q, H
foreground = "#000000"
background = "#ffffff"
include_logo = false               # Embed site logo in center
logo_path = "static/logo.png"
```

**Behavior:**

1. For each post, generate QR code for its absolute URL
2. Save to output directory
3. Add `qrcode_url` field to post model

**Output Files:**

```
output/
  qrcodes/
    hello-world.svg
    another-post.svg
```

**Post Model Extension:**

```python
class QRCodePostFields(BaseModel):
    qrcode_url: Optional[str] = None
```

**Template Usage:**

```jinja2
{% if post.qrcode_url %}
<img src="{{ post.qrcode_url }}" alt="QR Code" class="qr-code" />
{% endif %}
```

**Hook Signature:**

```python
@hook_impl
def save(core):
    config = core.config.qrcode
    if not config.enabled:
        return
    
    qr_dir = Path(core.config.output_dir) / config.output_dir
    qr_dir.mkdir(parents=True, exist_ok=True)
    
    for post in core.filter("not skip"):
        url = f"{core.config.url}{post.href}"
        qr_path = generate_qr(url, post.slug, qr_dir, config)
        post.qrcode_url = f"/{config.output_dir}/{qr_path.name}"
```

---

## Plugin Dependencies

| Plugin | Required | Optional |
|--------|----------|----------|
| `glossary` | - | - |
| `mermaid` | - | - |
| `chartjs` | - | - |
| `csv_fence` | - | - |
| `md_video` | - | - |
| `one_line_link` | - | `lxml` (metadata fetching) |
| `wikilink_hover` | `lxml` | - |
| `qrcode` | `pyqrcode` or `qrcode` | `pillow` (for PNG) |

Install optional dependencies:

```bash
pip install your-ssg[qrcode]      # QR code support
pip install your-ssg[hover]       # Wikilink hover support
pip install your-ssg[all-plugins] # All optional plugin deps
```

---

## Plugin Interaction

### Execution Order

When using multiple optional plugins, order matters:

```toml
hooks = [
    "default",           # Includes wikilinks
    "glossary",          # After render (needs article_html)
    "mermaid",           # After render
    "chartjs",           # After render
    "csv_fence",         # After render
    "md_video",          # After render
    "one_line_link",     # After render
    "wikilink_hover",    # After wikilinks (needs wikilink anchors)
    "qrcode",            # At save stage
]
```

### Conflicts

| Plugins | Issue | Resolution |
|---------|-------|------------|
| `glossary` + `wikilinks` | Terms in wikilinks get double-linked | Glossary excludes content inside `<a>` tags |
| `mermaid` + `csv_fence` | None | Process independently |
| `one_line_link` + any | URLs in code blocks | Exclude `<pre>` and `<code>` content |

---

## Disabling Optional Plugins

```toml
[your-ssg]
hooks = ["default", "mermaid", "chartjs"]
disabled_hooks = ["mermaid"]  # Disable mermaid specifically
```

Or per-post via frontmatter:

```yaml
---
title: No Diagrams Here
mermaid: false
chartjs: false
---
```

---

## Creating Custom Optional Plugins

Follow the pattern established by these plugins:

1. **Configuration Model**: Define with sensible defaults
2. **Post Model Extension**: Add fields if needed
3. **Enabled Check**: Always check `config.enabled` first
4. **Filter Usage**: Use `core.filter("not skip")` 
5. **Caching**: Cache expensive operations
6. **Error Handling**: Gracefully handle missing dependencies

Example skeleton:

```python
"""
My Custom Plugin

Configuration:
    [your-ssg.my_plugin]
    enabled = true
    option = "value"
"""
from typing import TYPE_CHECKING
import pydantic

from your_ssg.hookspec import hook_impl

if TYPE_CHECKING:
    from your_ssg import Core


class MyPluginConfig(pydantic.BaseModel):
    enabled: bool = True
    option: str = "default"


class Config(pydantic.BaseModel):
    my_plugin: MyPluginConfig = MyPluginConfig()


@hook_impl
def config_model(core: "Core") -> None:
    core.config_models.append(Config)


@hook_impl
def post_render(core: "Core") -> None:
    config = core.config.my_plugin
    if not config.enabled:
        return
    
    for post in core.filter("not skip"):
        # Process post
        pass
```

---

## See Also

- [PLUGINS.md](./PLUGINS.md) - Plugin development guide
- [DEFAULT_PLUGINS.md](./DEFAULT_PLUGINS.md) - Default plugins
- [CONFIG.md](./CONFIG.md) - Plugin configuration
- [THEMES.md](./THEMES.md) - Styling for plugin output
