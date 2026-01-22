# Template System Specification

This document specifies the template system for rendering content.

## Overview

Templates wrap rendered markdown content in HTML layouts. The system supports:

- Template inheritance (base templates)
- Includes (partials)
- Variables and expressions
- Control flow (if/for)
- Filters and functions
- Custom template per post

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

```python
def reading_time(content, wpm=200):
    words = len(content.split())
    minutes = max(1, words // wpm)
    return f"{minutes} min read"

jinja_env.filters["reading_time"] = reading_time
```

Usage:
```jinja2
{{ post.content | reading_time }}
{{ post.content | reading_time(250) }}
```

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

---

## Error Handling

### Undefined Variables

| Mode | Behavior |
|------|----------|
| **Silent** (recommended) | Undefined → empty string |
| **Strict** | Undefined → error |
| **Debug** | Undefined → `{{ undefined }}` |

Configure silent mode:
```python
from jinja2 import Environment, Undefined

class SilentUndefined(Undefined):
    def _fail_with_undefined_error(self, *args, **kwargs):
        return ""

env = Environment(undefined=SilentUndefined)
```

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

Cache compiled templates:

```python
from jinja2 import Environment, FileSystemBytecodeCache

env = Environment(
    bytecode_cache=FileSystemBytecodeCache(
        directory=".tool.cache/templates"
    )
)
```

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

### Python
- **Jinja2** - Standard, fast, full-featured

### JavaScript
- **Nunjucks** - Jinja-like, by Mozilla
- **EJS** - Simple, embedded JS
- **Handlebars** - Logic-less templates

### Go
- **html/template** - Standard library
- **pongo2** - Jinja-like

### Rust
- **Tera** - Jinja-like
- **Handlebars** - Port of JS library

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

**Implementation Requirement:** If using pongo2, implement a `strftime` filter that translates standard format codes:

```go
// Recommended: Add strftime filter for consistency
{{ post.date | strftime:"%B %d, %Y" }}  // January 15, 2024
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

## See Also

- [SPEC.md](./SPEC.md) - Full specification
- [THEMES.md](./THEMES.md) - Theming and template overrides
- [CONFIG.md](./CONFIG.md) - Template configuration
- [CONTENT.md](./CONTENT.md) - Markdown processing
- [PLUGINS.md](./PLUGINS.md) - Plugin development
