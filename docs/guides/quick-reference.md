---
title: "Quick Reference"
description: "Common commands, configuration snippets, and troubleshooting tips for markata-go"
date: 2026-01-24
published: true
slug: /docs/guides/quick-reference/
tags:
  - documentation
  - guides
  - reference
---

# Quick Reference

A quick reference for common markata-go commands, configuration snippets, and troubleshooting tips.

---

## Prerequisites

Before using this guide, you should have:
- markata-go installed ([Installation Guide](/docs/getting-started/#installation))
- A basic understanding of TOML configuration format

---

## CLI Commands

### Build and Serve

```bash
# Build site
markata-go build

# Build with clean output directory
markata-go build --clean

# Build to custom output directory
markata-go build -o dist

# Start development server with live reload
markata-go serve

# Serve on custom port
markata-go serve -p 3000

# Serve with verbose output
markata-go serve -v
```

### Content Creation

```bash
# Create a new post
markata-go new "My Post Title"

# Create with specific template
markata-go new "Tutorial" -t tutorial

# Create in specific directory
markata-go new "About" --dir pages

# Create as published (not draft)
markata-go new "Announcement" --draft=false

# List available templates
markata-go new --list
```

### Configuration

```bash
# Initialize config file
markata-go config init

# Show resolved configuration
markata-go config show

# Show as JSON
markata-go config show --json

# Get specific value
markata-go config get output_dir
markata-go config get glob.patterns

# Validate configuration
markata-go config validate
```

### Themes and Palettes

```bash
# List available palettes
markata-go palette list

# Show palette details
markata-go palette info catppuccin-mocha

# Check accessibility
markata-go palette check catppuccin-mocha

# Export palette as CSS
markata-go palette export catppuccin-mocha --format css

# Clone a palette for customization
markata-go palette clone catppuccin-mocha --name "my-theme"
```

---

## Rich Markdown icon shortcuts

Use the existing Lucide shortcode in any Markdown post:

````markdown
:lucide-rocket: Deployment

::: md-feature

:::: md-feature-icon
:lucide-rocket:

::::

:::: md-feature-body
### Verify production
Confirm the deployed image, with visible text for accessibility.

::::

:::
````

Combine `md-icon` with `md-card` for centered icon cards. See
[Icon Shortcodes](/docs/guides/icons/) and
[Markdown Containers](/docs/guides/markdown-containers/). Agents should read
the bundled skill's `topics/rich-markdown.md` before inventing custom icon
code.

---

## Configuration Snippets

### Minimal Site Config

```toml
[markata-go]
title = "My Site"
url = "https://example.com"
output_dir = "public"

[markata-go.glob]
patterns = ["posts/**/*.md"]
```

### Blog with RSS

```toml
[markata-go]
title = "My Blog"
url = "https://myblog.com"
output_dir = "public"

[markata-go.glob]
patterns = ["posts/**/*.md", "pages/*.md"]

[[markata-go.feeds]]
slug = "blog"
title = "Blog"
filter = "published == True"
sort = "date"
reverse = true

[markata-go.feeds.formats]
html = true
rss = true
atom = true
```

### Documentation Site

```toml
[markata-go]
title = "Project Docs"
url = "https://docs.example.com"
output_dir = "site"

[markata-go.glob]
patterns = ["docs/**/*.md"]

[markata-go.markdown]
extensions = ["tables", "tasklist"]

[markata-go.components.doc_sidebar]
enabled = true
position = "right"
min_depth = 2
max_depth = 4

[search]
enabled = true
position = "navbar"
```

### Custom Theme

```toml
[markata-go.theme]
palette = "catppuccin-mocha"
custom_css = "custom.css"

[markata-go.theme.variables]
"--color-primary" = "#8b5cf6"
"--content-width" = "800px"
```

### Navigation Config

```toml
[markata-go.components.nav]
enabled = true
position = "header"

[[markata-go.components.nav.items]]
label = "Home"
url = "/"

[[markata-go.components.nav.items]]
label = "Blog"
url = "/blog/"

[[markata-go.components.nav.items]]
label = "GitHub"
url = "https://github.com/user"
external = true
```

---

## Frontmatter Snippets

### Basic Post

```yaml
---
title: "My Post"
date: 2026-01-24
published: true
tags:
  - topic
---
```

### Full Post with Custom Fields

```yaml
---
title: "Complete Guide to X"
slug: "guide-to-x"
date: 2026-01-24
published: true
tags:
  - tutorial
  - guide
description: "A comprehensive guide covering everything about X."
template: "tutorial.html"

# Custom fields (available via post.Extra)
author: "Jane Doe"
reading_time: "10 min"
featured: true
series: "X Fundamentals"
series_order: 1
---
```

### Draft Post

```yaml
---
title: "Work in Progress"
date: 2026-01-24
published: false
draft: true
---
```

### Landing Page

```yaml
---
title: "Welcome"
slug: ""
published: true
template: "landing.html"
---
```

---

## Filter Expressions

Common filter expressions for feed configuration:

```toml
# Published posts only
filter = "published == True"

# Exclude drafts
filter = "published == True and draft == False"

# Posts with specific tag
filter = "'tutorial' in tags"

# Multiple tags (OR)
filter = "'python' in tags or 'go' in tags"

# Date range
filter = "date >= '2024-01-01' and date < '2025-01-01'"

# No future posts
filter = "published == True and date <= today"

# Featured posts
filter = "published == True and featured == True"

# By author (custom field)
filter = "author == 'Jane Doe'"

# Slug prefix (for sections)
filter = "slug.startswith('tutorials/')"
```

---

## Template Snippets

### Date Formatting

```django
{{ post.Date|date_format:"January 2, 2006" }}
{{ post.Date|date_format:"2006-01-02" }}
{{ post.Date|rss_date }}
{{ post.Date|atom_date }}
```

### Conditional Content

```django
{% if post.Tags %}
<ul class="tags">
  {% for tag in post.Tags %}
  <li><a href="/tags/{{ tag|slugify }}/">{{ tag }}</a></li>
  {% endfor %}
</ul>
{% endif %}
```

### Custom Field Access

```django
{% if post.Extra.featured %}
<span class="badge">Featured</span>
{% endif %}

{% if post.Extra.author %}
<p>By {{ post.Extra.author }}</p>
{% endif %}
```

### Pagination

```django
{% if page.HasPrev or page.HasNext %}
<nav class="pagination">
  {% if page.HasPrev %}
  <a href="{{ page.PrevURL }}">Previous</a>
  {% endif %}
  <span>Page {{ page.Number }}</span>
  {% if page.HasNext %}
  <a href="{{ page.NextURL }}">Next</a>
  {% endif %}
</nav>
{% endif %}
```

---

## Environment Variables

Override any config with `MARKATA_GO_` prefix:

```bash
# Core settings
export MARKATA_GO_OUTPUT_DIR=dist
export MARKATA_GO_URL=https://example.com
export MARKATA_GO_TITLE="My Site"

# Build with overrides
MARKATA_GO_URL=https://staging.example.com markata-go build
```

---

## Quick Troubleshooting

### Site Not Updating

```bash
# Clean build
markata-go build --clean

# Clear browser cache or use incognito mode
```

### Posts Not Appearing

1. Check frontmatter: `published: true` must be set
2. Check filter: Ensure feed filter matches your posts
3. Check date: Future dates won't appear with `date <= today` filter

```bash
# Validate config
markata-go config validate

# Check resolved config
markata-go config show
```

### RSS Feed Empty

1. Verify RSS is enabled in feed formats:
   ```toml
   [markata-go.feeds.formats]
   rss = true
   ```
2. Check that posts match the feed filter
3. Ensure `url` is set in config (required for absolute URLs)

### Template Not Found

1. Check template path in frontmatter matches file name
2. Verify file exists in `templates/` directory
3. Include `.html` extension: `template: "custom.html"`

### Styles Not Loading

1. Check `static/` directory contains CSS files
2. Verify `assets_dir` setting if using custom location
3. Check browser console for 404 errors

### Build Fails

```bash
# Verbose output for debugging
markata-go build -v

# Validate configuration first
markata-go config validate
```

---

## Common File Locations

| Purpose | Default Location |
|---------|-----------------|
| Configuration | `./markata-go.toml` |
| Content | `./posts/**/*.md` |
| Templates | `./templates/` |
| Static assets | `./static/` |
| Build output | `./public/` |
| Custom CSS | `./static/custom.css` |

---

## Next Steps

- **[Getting Started](/docs/getting-started/)** - Full installation tutorial
- **[Configuration Guide](/docs/guides/configuration/)** - Complete config reference
- **[Troubleshooting](/docs/troubleshooting/)** - Detailed problem solutions
- **[Guides Hub](/docs/guides/)** - Learning paths and all guides
