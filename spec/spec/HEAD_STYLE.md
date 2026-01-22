# Head and Style Injection Specification

The head and style injection system allows global and per-post customization of HTML `<head>` elements and CSS styling. This enables adding meta tags, scripts, stylesheets, and custom CSS without modifying templates.

## Integration with Themes

This system works alongside the [THEMES.md](./THEMES.md) theming system:

```
┌─────────────────────────────────────────────────────────────────────┐
│                    STYLE RESOLUTION ORDER                            │
├─────────────────────────────────────────────────────────────────────┤
│  1. Theme CSS (base styles from theme)                              │
│     └─ themes/default/static/css/main.css                           │
│                                                                      │
│  2. Theme CSS Variables (from theme.toml options)                   │
│     └─ [name.theme.variables] overrides                             │
│                                                                      │
│  3. Global head/style config ([name.head], [name.style])            │
│     └─ Site-wide meta tags, scripts, color overrides                │
│                                                                      │
│  4. Per-post config_overrides (frontmatter)                         │
│     └─ Post-specific customizations                                 │
│                                                                      │
│  LATER SOURCES OVERRIDE EARLIER ONES                                 │
└─────────────────────────────────────────────────────────────────────┘
```

**When to use each:**
- **Theme** (`[name.theme]`): Visual identity, fonts, layout, component styles
- **Theme Variables** (`[name.theme.variables]`): Quick color/spacing tweaks
- **Head/Style** (`[name.head]`, `[name.style]`): Meta tags, scripts, analytics, color overrides
- **Post Overrides** (`config_overrides`): Per-post customizations

## Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                    HEAD/STYLE INJECTION FLOW                         │
├─────────────────────────────────────────────────────────────────────┤
│  CONFIGURATION SOURCES (in order of precedence)                      │
│                                                                      │
│  1. Post frontmatter config_overrides                               │
│     └─ Highest priority, applies to single post                     │
│                                                                      │
│  2. Global config [name.head] and [name.style]                      │
│     └─ Applies to all posts                                          │
│                                                                      │
│  3. Theme defaults                                                   │
│     └─ Fallback values from theme                                   │
│                                                                      │
│  MERGE STRATEGY                                                      │
│     - Post overrides extend (not replace) global config              │
│     - Lists are concatenated (meta tags, links, scripts)            │
│     - Style values override individually                             │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Head Configuration

### Meta Tags

Add meta tags to every page:

```toml
[[name.head.meta]]
name = "author"
content = "Jane Doe"

[[name.head.meta]]
name = "robots"
content = "index, follow"

# Open Graph tags use 'property' instead of 'name'
[[name.head.meta]]
property = "og:type"
content = "article"

[[name.head.meta]]
property = "og:site_name"
content = "My Blog"
```

**Generated HTML:**
```html
<meta name="author" content="Jane Doe" />
<meta name="robots" content="index, follow" />
<meta property="og:type" content="article" />
<meta property="og:site_name" content="My Blog" />
```

### Link Tags

Add link elements (stylesheets, canonical URLs, favicons):

```toml
[[name.head.link]]
rel = "stylesheet"
href = "/css/custom.css"

[[name.head.link]]
rel = "icon"
href = "/favicon.ico"

[[name.head.link]]
rel = "preconnect"
href = "https://fonts.googleapis.com"
```

**Generated HTML:**
```html
<link rel="stylesheet" href="/css/custom.css" />
<link rel="icon" href="/favicon.ico" />
<link rel="preconnect" href="https://fonts.googleapis.com" />
```

### Script Tags

Add scripts to the head:

```toml
[[name.head.script]]
src = "https://cdn.tailwindcss.com"

[[name.head.script]]
src = "/js/analytics.js"
```

**Generated HTML:**
```html
<script src="https://cdn.tailwindcss.com"></script>
<script src="/js/analytics.js"></script>
```

### Raw HTML Text

For complex head content, use raw text:

```toml
[name.head]
text = '''
<style>
  :root {
    --custom-color: #ff6600;
  }
</style>
<script>
  console.log('Site loaded');
</script>
'''
```

Or as a list of text blocks:

```toml
[[name.head.text]]
value = '''
<style>
  .custom-class { color: red; }
</style>
'''

[[name.head.text]]
value = '''
<script type="application/ld+json">
  {"@context": "https://schema.org", "@type": "WebSite"}
</script>
'''
```

---

## Style Configuration

### Color Scheme

```toml
[name.style]
# Dark mode colors (default)
color_bg = "#1f2022"
color_bg_code = "#1f2022"
color_text = "#eefbfe"
color_link = "#fb30c4"
color_accent = "#e1bd00c9"

# Light mode colors
color_bg_light = "#eefbfe"
color_bg_code_light = "#eefbfe"
color_text_light = "#1f2022"
color_link_light = "#fb30c4"
color_accent_light = "#ffeb00"

# Layout
body_width = "800px"
overlay_brightness = ".85"
overlay_brightness_light = ".95"
```

### Style Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `color_bg` | Background color (dark) | `#1f2022` |
| `color_bg_code` | Code block background (dark) | `#1f2022` |
| `color_text` | Text color (dark) | `#eefbfe` |
| `color_link` | Link color (dark) | `#fb30c4` |
| `color_accent` | Accent color (dark) | `#e1bd00c9` |
| `color_bg_light` | Background color (light) | `#eefbfe` |
| `color_bg_code_light` | Code block background (light) | `#eefbfe` |
| `color_text_light` | Text color (light) | `#1f2022` |
| `color_link_light` | Link color (light) | `#fb30c4` |
| `color_accent_light` | Accent color (light) | `#ffeb00` |
| `body_width` | Maximum content width | `800px` |
| `overlay_brightness` | Image overlay brightness (dark) | `.85` |
| `overlay_brightness_light` | Image overlay brightness (light) | `.95` |

---

## Per-Post Overrides

The `config_overrides` frontmatter key allows posts to override **any** configuration value, not just head and style. This is the primary mechanism for per-post customization.

### Generic Config Override Syntax

```yaml
---
title: Special Post
config_overrides:
  # Override ANY config key using dot notation or nested structure
  output_dir: "special-output"           # Override top-level config
  
  markdown:
    extensions:
      tables: false                      # Override nested config
  
  feeds:
    defaults:
      items_per_page: 5                  # Override feed defaults
  
  head:
    meta:
      - name: robots
        content: noindex
  
  style:
    color_bg: "#000000"
  
  theme:
    options:
      show_toc: false                    # Override theme option
  
  # Custom plugin configuration
  my_plugin:
    enabled: false
---
```

### Override Resolution

```
┌─────────────────────────────────────────────────────────────────────┐
│                    CONFIG OVERRIDE RESOLUTION                        │
├─────────────────────────────────────────────────────────────────────┤
│  For each post, configuration is resolved as:                        │
│                                                                      │
│  1. Start with global config (merged from all sources)              │
│  2. Deep merge post.config_overrides                                 │
│  3. Result is used when rendering this specific post                │
│                                                                      │
│  MERGE RULES:                                                        │
│  - Scalars: Post value replaces global value                        │
│  - Objects: Deep merge (post keys override, others preserved)       │
│  - Arrays: Post array replaces global array (no concatenation)      │
│  - Special: head.meta/link/script arrays are APPENDED               │
└─────────────────────────────────────────────────────────────────────┘
```

### Common Override Patterns

**Disable a feature for one post:**
```yaml
---
title: No TOC Post
config_overrides:
  toc:
    enabled: false
---
```

**Use different template:**
```yaml
---
title: Landing Page
template: landing.html
config_overrides:
  style:
    body_width: "100%"
---
```

**Override markdown processing:**
```yaml
---
title: Code-Heavy Post
config_overrides:
  markdown:
    highlight:
      theme: "monokai"
      line_numbers: true
---
```

**Skip from feeds:**
```yaml
---
title: Hidden Post
config_overrides:
  feeds:
    exclude: true                        # Don't include in any feed
---
```

### Head/Style Overrides (Detailed)

For backwards compatibility and clarity, head and style have special merge behavior:

```yaml
---
title: Special Post
config_overrides:
  head:
    meta:
      - name: robots
        content: noindex
      - property: og:type
        content: video
    text: |
      <script src="/js/special-feature.js"></script>
  style:
    color_bg: "#000000"
    color_text: "#ffffff"
---
```

### Override Behavior

| Element | Behavior |
|---------|----------|
| `meta` | Post meta tags are **added** to global meta tags |
| `link` | Post links are **added** to global links |
| `script` | Post scripts are **added** to global scripts |
| `text` | Post text is **appended** to global text |
| `style.*` | Post style values **override** global values |

### Example: Video Post

```yaml
---
title: My Video Tutorial
config_overrides:
  head:
    meta:
      - property: og:type
        content: video.other
      - property: og:video
        content: https://example.com/video.mp4
    script:
      - src: /js/video-player.js
  style:
    body_width: "1200px"  # Wider for video content
---
```

### Example: Landing Page

```yaml
---
title: Welcome
template: landing.html
config_overrides:
  head:
    link:
      - rel: stylesheet
        href: /css/landing.css
    text: |
      <style>
        body { margin: 0; padding: 0; }
      </style>
  style:
    body_width: "100%"
---
```

---

## Template Integration

### Rendering Head Content

Templates access head configuration through `config.head`:

```jinja2
<head>
  <meta charset="UTF-8">
  <title>{{ post.title }} - {{ config.title }}</title>
  
  {# Render raw text first #}
  {{ config.head.text | safe }}
  
  {# Render meta tags #}
  {% for meta in config.head.meta %}
    {% if meta.name %}
    <meta name="{{ meta.name }}" content="{{ meta.content }}" />
    {% elif meta.property %}
    <meta property="{{ meta.property }}" content="{{ meta.content }}" />
    {% endif %}
  {% endfor %}
  
  {# Render link tags #}
  {% for link in config.head.link %}
    <link rel="{{ link.rel }}" href="{{ link.href }}" />
  {% endfor %}
  
  {# Render script tags #}
  {% for script in config.head.script %}
    <script src="{{ script.src }}"></script>
  {% endfor %}
</head>
```

### Using Style Variables

Access style configuration in templates:

```jinja2
<style>
  :root {
    --color-bg: {{ config.style.color_bg }};
    --color-text: {{ config.style.color_text }};
    --color-link: {{ config.style.color_link }};
    --color-accent: {{ config.style.color_accent }};
    --body-width: {{ config.style.body_width }};
  }
  
  @media (prefers-color-scheme: light) {
    :root {
      --color-bg: {{ config.style.color_bg_light }};
      --color-text: {{ config.style.color_text_light }};
      --color-link: {{ config.style.color_link_light }};
      --color-accent: {{ config.style.color_accent_light }};
    }
  }
  
  body {
    background: var(--color-bg);
    color: var(--color-text);
    max-width: var(--body-width);
    margin: 0 auto;
  }
  
  a {
    color: var(--color-link);
    text-decoration-color: var(--color-accent);
  }
</style>
```

### Dynamic Head Content

Use Jinja expressions in head configuration:

```toml
[[name.head.link]]
rel = "canonical"
href = "{{ config.url }}/{{ post.slug }}/"

[[name.head.meta]]
property = "og:url"
content = "{{ config.url }}/{{ post.slug }}/"

[[name.head.meta]]
property = "og:title"
content = "{{ post.title }}"

[[name.head.meta]]
property = "og:description"
content = "{{ post.description }}"
```

---

## Configuration Model

### HeadConfig

```python
class Meta(pydantic.BaseModel):
    name: Optional[str] = None
    property: Optional[str] = None
    content: str
    
    @validator('name')
    def check_og(cls, v):
        if v and v.startswith('og:'):
            raise ValueError("Use 'property' for og: tags, not 'name'")
        return v
    
    @root_validator
    def check_name_or_property(cls, values):
        if not values.get('name') and not values.get('property'):
            raise ValueError("Either 'name' or 'property' must be set")
        return values


class Link(pydantic.BaseModel):
    rel: str
    href: str


class Script(pydantic.BaseModel):
    src: str


class Text(pydantic.BaseModel):
    value: str


class HeadConfig(pydantic.BaseModel):
    meta: List[Meta] = []
    link: List[Link] = []
    script: List[Script] = []
    text: Union[List[Text], str] = ""
    
    @validator('text', pre=True)
    def normalize_text(cls, v):
        if isinstance(v, list):
            return "\n".join(item['value'] for item in v)
        return v
    
    @property
    def html(self) -> str:
        """Render all head content as HTML string."""
        parts = [self.text]
        for meta in self.meta:
            if meta.name:
                parts.append(f'<meta name="{meta.name}" content="{meta.content}" />')
            elif meta.property:
                parts.append(f'<meta property="{meta.property}" content="{meta.content}" />')
        for link in self.link:
            parts.append(f'<link rel="{link.rel}" href="{link.href}" />')
        for script in self.script:
            parts.append(f'<script src="{script.src}"></script>')
        return "\n".join(parts)
```

### StyleConfig

```python
class Style(pydantic.BaseModel):
    # Dark mode
    color_bg: str = "#1f2022"
    color_bg_code: str = "#1f2022"
    color_text: str = "#eefbfe"
    color_link: str = "#fb30c4"
    color_accent: str = "#e1bd00c9"
    overlay_brightness: str = ".85"
    
    # Light mode
    color_bg_light: str = "#eefbfe"
    color_bg_code_light: str = "#eefbfe"
    color_text_light: str = "#1f2022"
    color_link_light: str = "#fb30c4"
    color_accent_light: str = "#ffeb00"
    overlay_brightness_light: str = ".95"
    
    # Layout
    body_width: str = "800px"
```

### PostOverrides

```python
class PostOverrides(pydantic.BaseModel):
    head: HeadConfig = HeadConfig()
    style: Style = Style()


class Post(pydantic.BaseModel):
    config_overrides: PostOverrides = PostOverrides()
    # ... other post fields
```

---

## Hook Specification

### Stages

- `config_model`: Register HeadConfig and Style models
- `post_model`: Register PostOverrides on post model
- `pre_render`: Merge post overrides with global config

### Hook Signatures

```python
@hook_impl(tryfirst=True)
def config_model(core):
    """Register head and style configuration models."""
    core.config_models.append(Config)


@hook_impl(tryfirst=True)
def post_model(core):
    """Register post override model."""
    core.post_models.append(Post)


@hook_impl(tryfirst=True)
def pre_render(core):
    """Process post config overrides, merging text lists."""
    for post in core.posts:
        if 'config_overrides' not in post:
            continue
        
        # Normalize text list to string
        raw_text = post.config_overrides.get('head', {}).get('text', '')
        if isinstance(raw_text, list):
            post.config_overrides['head']['text'] = "\n".join(
                item['value'] for item in raw_text
            )
```

---

## Merged Configuration

When rendering, the template receives a merged configuration:

```python
def get_merged_config(core, post):
    """Merge global config with post overrides."""
    config = core.config.copy()
    
    if post.config_overrides:
        # Merge head
        config.head.meta.extend(post.config_overrides.head.meta)
        config.head.link.extend(post.config_overrides.head.link)
        config.head.script.extend(post.config_overrides.head.script)
        config.head.text += "\n" + post.config_overrides.head.text
        
        # Merge style (override individual values)
        for field in post.config_overrides.style.__fields__:
            value = getattr(post.config_overrides.style, field)
            if value is not None:
                setattr(config.style, field, value)
    
    return config
```

---

## Common Patterns

### Analytics Integration

```toml
[[name.head.script]]
src = "https://www.googletagmanager.com/gtag/js?id=GA_MEASUREMENT_ID"

[name.head]
text = '''
<script>
  window.dataLayer = window.dataLayer || [];
  function gtag(){dataLayer.push(arguments);}
  gtag('js', new Date());
  gtag('config', 'GA_MEASUREMENT_ID');
</script>
'''
```

### Font Loading

```toml
[[name.head.link]]
rel = "preconnect"
href = "https://fonts.googleapis.com"

[[name.head.link]]
rel = "preconnect"
href = "https://fonts.gstatic.com"
crossorigin = true

[[name.head.link]]
rel = "stylesheet"
href = "https://fonts.googleapis.com/css2?family=Inter:wght@400;600;700&display=swap"
```

### Social Media Cards

```toml
[[name.head.meta]]
property = "og:type"
content = "website"

[[name.head.meta]]
property = "og:image"
content = "{{ config.url }}/og-image.png"

[[name.head.meta]]
name = "twitter:card"
content = "summary_large_image"

[[name.head.meta]]
name = "twitter:site"
content = "@myhandle"
```

### Structured Data

```toml
[name.head]
text = '''
<script type="application/ld+json">
{
  "@context": "https://schema.org",
  "@type": "WebSite",
  "name": "{{ config.title }}",
  "url": "{{ config.url }}"
}
</script>
'''
```

---

## See Also

- [TEMPLATES.md](./TEMPLATES.md) - Template system specification
- [CONFIG.md](./CONFIG.md) - Configuration system
- [THEMES.md](./THEMES.md) - Theme customization
- [DATA_MODEL.md](./DATA_MODEL.md) - Post model specification
