---
title: "Advanced Usage"
description: "Power user guide covering complex configurations, optimization techniques, and advanced patterns"
date: 2024-01-15
published: true
tags:
  - documentation
  - advanced
---


This guide is for power users who want to leverage markata-go's full capabilities. It covers complex configurations, optimization techniques, and advanced patterns that go beyond the basics.

## Table of Contents

- [Complex Feed Configurations](#complex-feed-configurations)
- [Dynamic Content with Jinja-in-Markdown](#dynamic-content-with-jinja-in-markdown)
- [Advanced Filtering](#advanced-filtering)
- [Custom Plugin Development](#custom-plugin-development)
- [Performance Optimization](#performance-optimization)
- [Multi-Environment Builds](#multi-environment-builds)
- [Advanced Templates](#advanced-templates)
- [Wikilinks and Internal Linking](#wikilinks-and-internal-linking)
- [Table of Contents Generation](#table-of-contents-generation)
- [Admonitions and Callouts](#admonitions-and-callouts)

---

## Complex Feed Configurations

The feed system is markata-go's most powerful feature. Beyond basic usage, you can create sophisticated content architectures with multiple feeds, nested hierarchies, and custom API endpoints.

**Related guides:** [[feeds-guide|Feeds Guide]], [[syndication-feeds|Syndication Feeds]]

### Multiple Feeds with Different Filters

Create specialized feeds that slice your content in different ways:

```toml
[markata-go]
title = "Tech Blog"
url = "https://example.com"

# =============================================================================
# FEED DEFAULTS
# =============================================================================
[markata-go.feeds.defaults]
items_per_page = 10

[markata-go.feeds.defaults.formats]
html = true
rss = true

# =============================================================================
# CONTENT FEEDS
# =============================================================================

# Home page - curated recent posts
[[markata-go.feeds]]
slug = ""
title = "Home"
filter = "published == True and date <= today"
sort = "date"
reverse = true
items_per_page = 5
formats = { rss = false }

# Main blog archive
[[markata-go.feeds]]
slug = "blog"
title = "All Posts"
filter = "published == True"
sort = "date"
reverse = true

# Tutorials only
[[markata-go.feeds]]
slug = "tutorials"
title = "Tutorials"
description = "Step-by-step programming guides"
filter = "published == True and 'tutorial' in tags"
sort = "date"
reverse = true

# Featured/curated content
[[markata-go.feeds]]
slug = "featured"
title = "Featured Articles"
description = "Hand-picked must-read articles"
filter = "published == True and featured == True"
sort = "date"
reverse = true
items_per_page = 6

# Quick tips (short-form content)
[[markata-go.feeds]]
slug = "tips"
title = "Quick Tips"
filter = "published == True and 'tip' in tags"
sort = "date"
reverse = true
items_per_page = 20

# Long-form essays
[[markata-go.feeds]]
slug = "essays"
title = "Essays"
filter = "published == True and category == 'essay'"
sort = "date"
reverse = true
items_per_page = 5
```

### Nested Feeds (Categories within Tags)

Create hierarchical content organization with nested URL structures:

```toml
# Top-level language categories
[[markata-go.feeds]]
slug = "go"
title = "Go Articles"
filter = "published == True and 'go' in tags"
sort = "date"
reverse = true

[[markata-go.feeds]]
slug = "python"
title = "Python Articles"
filter = "published == True and 'python' in tags"
sort = "date"
reverse = true

# Nested: Go tutorials
[[markata-go.feeds]]
slug = "go/tutorials"
title = "Go Tutorials"
filter = "published == True and 'go' in tags and 'tutorial' in tags"
sort = "date"
reverse = true

# Nested: Go tips
[[markata-go.feeds]]
slug = "go/tips"
title = "Go Tips"
filter = "published == True and 'go' in tags and 'tip' in tags"
sort = "date"
reverse = true

# Nested: Python tutorials
[[markata-go.feeds]]
slug = "python/tutorials"
title = "Python Tutorials"
filter = "published == True and 'python' in tags and 'tutorial' in tags"
sort = "date"
reverse = true

# Nested: Python advanced
[[markata-go.feeds]]
slug = "python/advanced"
title = "Advanced Python"
filter = "published == True and 'python' in tags and difficulty == 'advanced'"
sort = "date"
reverse = true
```

**Generated structure:**
```
/go/
  index.html
  rss.xml
  tutorials/
    index.html
    rss.xml
  tips/
    index.html
    rss.xml
/python/
  index.html
  tutorials/
    index.html
  advanced/
    index.html
```

### Custom JSON API Endpoints

Create JSON-only feeds to serve as lightweight APIs:

```toml
# Full posts API (all fields)
[[markata-go.feeds]]
slug = "api/posts"
title = "Posts API"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 0                  # No pagination - all posts

[markata-go.feeds.formats]
html = false
rss = false
json = true

# Lightweight posts list (for autocomplete, etc.)
[[markata-go.feeds]]
slug = "api/posts/list"
title = "Posts List"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 0

[markata-go.feeds.formats]
json = true
html = false
rss = false

[markata-go.feeds.templates]
json = "api-list.json"              # Custom minimal template

# Posts by year
[[markata-go.feeds]]
slug = "api/posts/2024"
title = "2024 Posts"
filter = "published == True and date >= '2024-01-01' and date < '2025-01-01'"
sort = "date"
reverse = true
items_per_page = 0

[markata-go.feeds.formats]
json = true
html = false
```

Custom API template (`templates/api-list.json`):

```json
{
  "count": {{ feed.posts|length }},
  "posts": [
    {% for post in feed.posts %}
    {
      "title": {{ post.title|tojson }},
      "slug": {{ post.slug|tojson }},
      "href": {{ post.href|tojson }},
      "date": {{ post.date|atom_date|tojson }},
      "tags": {{ post.tags|tojson }}
    }{% if not forloop.Last %},{% endif %}
    {% endfor %}
  ]
}
```

### Search Indexes

Generate a JSON index optimized for client-side search (e.g., with Lunr.js, Fuse.js, or Pagefind):

```toml
[[markata-go.feeds]]
slug = "search-index"
title = "Search Index"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 0

[markata-go.feeds.formats]
html = false
rss = false
json = true

[markata-go.feeds.templates]
json = "search-index.json"
```

Search index template (`templates/search-index.json`):

```json
{
  "index": [
    {% for post in feed.posts %}
    {
      "id": {{ forloop.Counter0 }},
      "title": {{ post.title|tojson }},
      "href": {{ post.href|tojson }},
      "content": {{ post.content|striptags|truncate:1000|tojson }},
      "description": {{ post.description|default_if_none:""|tojson }},
      "tags": {{ post.tags|tojson }},
      "date": {{ post.date|date_format:"2006-01-02"|tojson }}
    }{% if not forloop.Last %},{% endif %}
    {% endfor %}
  ]
}
```

JavaScript integration with Fuse.js:

```javascript
// Load search index
fetch('/search-index/feed.json')
  .then(r => r.json())
  .then(data => {
    const fuse = new Fuse(data.index, {
      keys: ['title', 'content', 'tags', 'description'],
      threshold: 0.3,
      includeScore: true
    });

    // Search function
    window.search = (query) => {
      const results = fuse.search(query);
      return results.map(r => ({
        title: r.item.title,
        href: r.item.href,
        score: r.score
      }));
    };
  });
```

---

## Dynamic Content with Jinja-in-Markdown

The `jinja_md` plugin enables powerful dynamic content generation directly within your Markdown files.

**Related guide:** [[dynamic-content|Dynamic Content]]

### Embedding Post Lists in Content

Create dynamic index pages that automatically update:

```markdown
---
title: "Documentation"
jinja: true
---

# Documentation

Welcome to the docs. Here's everything organized by topic.

## Getting Started

{% for p in core.filter("published == True and 'getting-started' in tags") %}
- [{{ p.Title }}]({{ p.Href }}) - {{ p.Description|default_if_none:""|truncate:80 }}
{% endfor %}

## Core Concepts

{% for p in core.filter("published == True and 'core-concepts' in tags") %}
- [{{ p.Title }}]({{ p.Href }})
{% endfor %}

## API Reference

{% for p in core.filter("published == True and 'api' in tags") %}
- [{{ p.Title }}]({{ p.Href }})
{% endfor %}
```

### Series Navigation

Build automatic series navigation that stays in sync:

```markdown
---
title: "Building a CLI in Go: Part 2"
series: "go-cli"
series_order: 2
jinja: true
---

<nav class="series-nav">
<strong>Series: Building a CLI in Go</strong>

{% set series_posts = [] %}
{% for p in core.filter("published == True") %}
{% if p.Extra.series == "go-cli" %}
{% set _ = series_posts.append(p) %}
{% endif %}
{% endfor %}

<ol>
{% for p in series_posts|sort(attribute='Extra.series_order') %}
{% if p.Extra.series_order == post.Extra.series_order %}
<li><strong>{{ p.Title }}</strong> (current)</li>
{% else %}
<li><a href="{{ p.Href }}">{{ p.Title }}</a></li>
{% endif %}
{% endfor %}
</ol>

{% set current_order = post.Extra.series_order %}
<div class="series-pagination">
{% for p in series_posts %}
{% if p.Extra.series_order == current_order - 1 %}
<a href="{{ p.Href }}" class="prev">Previous: {{ p.Title }}</a>
{% endif %}
{% if p.Extra.series_order == current_order + 1 %}
<a href="{{ p.Href }}" class="next">Next: {{ p.Title }}</a>
{% endif %}
{% endfor %}
</div>
</nav>

# Building a CLI in Go: Part 2

Now let's add subcommands to our CLI...
```

### Dynamic Table of Contents

Generate a table of contents from your content structure:

```markdown
---
title: "Complete Go Guide"
jinja: true
---

# Complete Go Guide

## Table of Contents

{% set beginner = core.filter("published == True and 'go' in tags and difficulty == 'beginner'") %}
{% set intermediate = core.filter("published == True and 'go' in tags and difficulty == 'intermediate'") %}
{% set advanced = core.filter("published == True and 'go' in tags and difficulty == 'advanced'") %}

### Beginner ({{ beginner|length }} articles)

{% for p in beginner %}
1. [{{ p.Title }}]({{ p.Href }})
{% endfor %}

### Intermediate ({{ intermediate|length }} articles)

{% for p in intermediate %}
1. [{{ p.Title }}]({{ p.Href }})
{% endfor %}

### Advanced ({{ advanced|length }} articles)

{% for p in advanced %}
1. [{{ p.Title }}]({{ p.Href }})
{% endfor %}
```

### Related Posts Widget

Show contextually related posts:

```markdown
---
title: "Understanding Go Interfaces"
tags: [go, interfaces, tutorial]
jinja: true
---

# Understanding Go Interfaces

[Your article content...]

---

## Related Articles

{% set related = [] %}
{% for p in core.filter("published == True") %}
{% if p.Slug != post.Slug %}
{% for tag in post.Tags %}
{% if tag in p.Tags and p not in related %}
{% set _ = related.append(p) %}
{% endif %}
{% endfor %}
{% endif %}
{% endfor %}

{% for p in related[:5] %}
- [{{ p.Title }}]({{ p.Href }})
{% endfor %}
```

---

## Advanced Filtering

The filter expression system supports complex queries for precise content selection.

**Related guides:** [[frontmatter-guide|Frontmatter]], [[feeds-guide|Feeds Guide]]

### Complex Filter Expressions

Combine multiple conditions with logical operators:

```toml
# Published, non-draft, with specific tag, from current year
filter = "published == True and draft == False and 'tutorial' in tags and date >= '2024-01-01'"

# Multiple tags (OR logic)
filter = "published == True and ('go' in tags or 'python' in tags or 'rust' in tags)"

# Exclude specific content
filter = "published == True and slug != 'about' and slug != 'contact'"

# Complex nested logic
filter = "(published == True and featured == True) or (published == True and 'highlight' in tags)"
```

### Date-Based Filtering

Use dynamic date comparisons:

```toml
# Only past/present posts (no scheduled future posts)
filter = "published == True and date <= today"

# Posts from the last 30 days
filter = "published == True and date >= today - 30"

# Posts from a specific year
filter = "published == True and date >= '2024-01-01' and date < '2025-01-01'"

# Posts from a specific month
filter = "published == True and date >= '2024-06-01' and date < '2024-07-01'"

# Evergreen content (no date or old date is fine)
filter = "published == True and (date == None or date <= today)"
```

### Custom Field Filtering

Filter on any frontmatter field:

```toml
# By author
filter = "published == True and author == 'Jane Doe'"

# By category
filter = "published == True and category == 'Backend'"

# By difficulty level
filter = "published == True and difficulty == 'beginner'"

# By custom boolean
filter = "published == True and sponsored == False"

# By series membership
filter = "published == True and series == 'Building a Blog'"

# By numeric field
filter = "published == True and word_count >= 1000"

# By reading time
filter = "published == True and reading_time_minutes <= 5"
```

### String Method Filters

Use string methods for pattern matching:

```toml
# Slugs with a prefix
filter = "published == True and slug.startswith('tutorials/')"

# Slugs with a suffix
filter = "published == True and slug.endswith('-guide')"

# Title contains keyword
filter = "published == True and title.lower().contains('docker')"

# Category starts with
filter = "published == True and category.startswith('Web')"
```

### Filter Cheat Sheet

| Pattern | Example |
|---------|---------|
| Equality | `field == 'value'` |
| Inequality | `field != 'value'` |
| Boolean check | `published == True` |
| Tag membership | `'tag' in tags` |
| Date comparison | `date >= '2024-01-01'` |
| Today's date | `date <= today` |
| Null check | `field == None` |
| String prefix | `slug.startswith('prefix/')` |
| String contains | `title.lower().contains('word')` |
| AND | `cond1 and cond2` |
| OR | `cond1 or cond2` |
| NOT | `not condition` |
| Grouping | `(cond1 or cond2) and cond3` |

---

## Custom Plugin Development

When built-in functionality isn't enough, write custom plugins to extend markata-go.

**Related guide:** [[plugin-development|Plugin Development]]

### When to Write a Plugin

Consider a custom plugin when you need to:

- Add computed fields to posts (reading time, word count, etc.)
- Process content in a custom way (shortcodes, custom syntax)
- Integrate with external services (APIs, databases)
- Generate additional output files
- Implement custom validation rules

### Plugin Architecture Overview

Plugins hook into markata-go's 9-stage lifecycle:

```
configure → validate → glob → load → transform → render → collect → write → cleanup
```

Each stage has a specific purpose:

| Stage | Purpose | Common Uses |
|-------|---------|-------------|
| Configure | Initialize plugin state | Read config, set up clients |
| Validate | Validate configuration | Check required fields |
| Glob | Discover files | Custom file patterns |
| Load | Parse files to posts | Custom parsers |
| Transform | Pre-render processing | Computed fields, shortcodes |
| Render | Convert to HTML | Markdown extensions |
| Collect | Build aggregations | Custom feeds, indexes |
| Write | Write output | Custom file generation |
| Cleanup | Release resources | Close connections |

### Example: Custom Shortcode Plugin

Here's a complete plugin that adds shortcode support (e.g., `{{< youtube id="..." >}}`):

```go
package plugins

import (
    "fmt"
    "regexp"
    "strings"

    "github.com/example/markata-go/pkg/lifecycle"
    "github.com/example/markata-go/pkg/models"
)

// ShortcodePlugin processes shortcodes in markdown content.
type ShortcodePlugin struct {
    shortcodes map[string]ShortcodeFunc
}

// ShortcodeFunc is a function that processes a shortcode.
type ShortcodeFunc func(params map[string]string) string

// NewShortcodePlugin creates a new ShortcodePlugin with built-in shortcodes.
func NewShortcodePlugin() *ShortcodePlugin {
    p := &ShortcodePlugin{
        shortcodes: make(map[string]ShortcodeFunc),
    }

    // Register built-in shortcodes
    p.Register("youtube", youtubeShortcode)
    p.Register("gist", gistShortcode)
    p.Register("tweet", tweetShortcode)
    p.Register("figure", figureShortcode)

    return p
}

// Name returns the plugin name.
func (p *ShortcodePlugin) Name() string {
    return "shortcodes"
}

// Register adds a custom shortcode.
func (p *ShortcodePlugin) Register(name string, fn ShortcodeFunc) {
    p.shortcodes[name] = fn
}

// Transform processes shortcodes in all posts.
func (p *ShortcodePlugin) Transform(m *lifecycle.Manager) error {
    return m.ProcessPostsConcurrently(func(post *models.Post) error {
        if post.Skip || post.Content == "" {
            return nil
        }
        post.Content = p.processShortcodes(post.Content)
        return nil
    })
}

// Shortcode pattern: {{< name param="value" >}}
var shortcodePattern = regexp.MustCompile(`\{\{<\s*(\w+)\s*([^>]*)\s*>\}\}`)

// Param pattern: key="value" or key='value'
var paramPattern = regexp.MustCompile(`(\w+)=["']([^"']*)["']`)

func (p *ShortcodePlugin) processShortcodes(content string) string {
    return shortcodePattern.ReplaceAllStringFunc(content, func(match string) string {
        submatches := shortcodePattern.FindStringSubmatch(match)
        if len(submatches) < 2 {
            return match
        }

        name := submatches[1]
        paramsStr := ""
        if len(submatches) >= 3 {
            paramsStr = submatches[2]
        }

        // Parse parameters
        params := make(map[string]string)
        paramMatches := paramPattern.FindAllStringSubmatch(paramsStr, -1)
        for _, pm := range paramMatches {
            if len(pm) >= 3 {
                params[pm[1]] = pm[2]
            }
        }

        // Execute shortcode
        if fn, ok := p.shortcodes[name]; ok {
            return fn(params)
        }

        return match // Unknown shortcode, leave as-is
    })
}

// Built-in shortcodes

func youtubeShortcode(params map[string]string) string {
    id := params["id"]
    if id == "" {
        return "<!-- youtube: missing id -->"
    }
    return fmt.Sprintf(`<div class="video-container">
<iframe src="https://www.youtube.com/embed/%s"
        frameborder="0"
        allowfullscreen
        loading="lazy"></iframe>
</div>`, id)
}

func gistShortcode(params map[string]string) string {
    user := params["user"]
    id := params["id"]
    if user == "" || id == "" {
        return "<!-- gist: missing user or id -->"
    }
    return fmt.Sprintf(`<script src="https://gist.github.com/%s/%s.js"></script>`, user, id)
}

func tweetShortcode(params map[string]string) string {
    id := params["id"]
    if id == "" {
        return "<!-- tweet: missing id -->"
    }
    return fmt.Sprintf(`<blockquote class="twitter-tweet" data-dnt="true">
<a href="https://twitter.com/x/status/%s"></a>
</blockquote>
<script async src="https://platform.twitter.com/widgets.js"></script>`, id)
}

func figureShortcode(params map[string]string) string {
    src := params["src"]
    alt := params["alt"]
    caption := params["caption"]

    if src == "" {
        return "<!-- figure: missing src -->"
    }

    html := fmt.Sprintf(`<figure>
<img src="%s" alt="%s" loading="lazy">`, src, alt)

    if caption != "" {
        html += fmt.Sprintf(`
<figcaption>%s</figcaption>`, caption)
    }

    html += `
</figure>`

    return html
}

// Interface verification
var (
    _ lifecycle.Plugin          = (*ShortcodePlugin)(nil)
    _ lifecycle.TransformPlugin = (*ShortcodePlugin)(nil)
)
```

**Usage in markdown:**

```markdown
---
title: "Video Tutorial"
---

# Getting Started with Go

Watch the video tutorial:

{{< youtube id="YS4e4q9oBaU" >}}

Check out this code snippet:

{{< gist user="example" id="abc123" >}}

{{< figure src="/images/diagram.png" alt="Architecture" caption="System architecture diagram" >}}
```

---

## Performance Optimization

Optimize markata-go for large sites with thousands of posts.

### Concurrent Processing

markata-go processes posts concurrently by default. Control concurrency via configuration:

```toml
[markata-go]
concurrency = 0    # 0 = auto-detect (uses all CPU cores)
# concurrency = 4  # Limit to 4 workers
# concurrency = 1  # Sequential processing (for debugging)
```

For custom plugins, use the built-in concurrent processor:

```go
func (p *MyPlugin) Transform(m *lifecycle.Manager) error {
    return m.ProcessPostsConcurrently(func(post *models.Post) error {
        // This runs in parallel across posts
        return p.processPost(post)
    })
}
```

### Incremental Builds

For faster development, markata-go can skip unchanged files:

```bash
# Standard build (processes all files)
markata-go build

# Serve with watch mode (enabled by default)
markata-go serve

# Serve without file watching
markata-go serve --no-watch
```

By default, the dev server rebuilds deterministically on every change while still using build cache to skip expensive per-post work. For the fastest loop, use:

```bash
markata-go serve --fast
```

Fast mode rebuilds only the changed posts plus any dependents discovered from the dependency graph, and disables blogroll/mentions to keep the loop tight. Feed collection also skips when no changed posts touch a feed. If a change alters dependencies (links, embeds), the server performs a second pass to update the affected posts.

### Caching Strategies

For expensive operations, implement caching in your plugins:

```go
func (p *MyPlugin) Transform(m *lifecycle.Manager) error {
    cache := m.Cache()

    for _, post := range m.Posts() {
        // Generate cache key from content hash
        key := fmt.Sprintf("processed:%s", post.ContentHash())

        // Check cache
        if cached, ok := cache.Get(key); ok {
            post.Set("processed_data", cached)
            continue
        }

        // Expensive operation
        result := p.expensiveOperation(post)

        // Store in cache
        cache.Set(key, result)
        post.Set("processed_data", result)
    }

    return nil
}
```

### Large Site Considerations

For sites with 1000+ posts:

1. **Limit feed items:**
   ```toml
   [markata-go.feeds.syndication]
   max_items = 50    # Don't include all posts in RSS/Atom
   ```

2. **Use pagination:**
   ```toml
   [[markata-go.feeds]]
   items_per_page = 20    # Reasonable page size
   orphan_threshold = 5
   ```

3. **Minimize Jinja-in-Markdown:**
   - Pre-compute expensive queries
   - Cache dynamic content
   - Use feeds instead of inline queries

4. **Optimize glob patterns:**
   ```toml
   [markata-go.glob]
   patterns = ["posts/**/*.md"]    # Specific patterns
   use_gitignore = true            # Skip ignored files
   ```

5. **Profile builds:**
   ```bash
   markata-go build -v --profile
   ```

### Build Time Benchmarks

| Posts | Cold Build | Incremental | With Watch |
|-------|------------|-------------|------------|
| 100 | ~1s | ~0.2s | ~50ms |
| 1,000 | ~5s | ~0.5s | ~100ms |
| 10,000 | ~30s | ~2s | ~200ms |

*Times vary based on content complexity, plugins, and hardware.*

---

## Multi-Environment Builds

Configure markata-go for different deployment environments.

**Related guides:** [[configuration-guide|Configuration]], [[deployment-guide|Deployment]]

### Staging vs Production

Use environment variables to customize builds:

```bash
# Production build
MARKATA_GO_URL=https://example.com markata-go build --clean

# Staging build
MARKATA_GO_URL=https://staging.example.com markata-go build --clean

# Preview build (e.g., for PR previews)
MARKATA_GO_URL=https://preview-123.example.com markata-go build
```

### Environment Variables Reference

All configuration can be overridden via environment variables:

```bash
# Core settings
export MARKATA_GO_URL=https://example.com
export MARKATA_GO_TITLE="My Site"
export MARKATA_GO_OUTPUT_DIR=dist
export MARKATA_GO_CONCURRENCY=4

# Feature flags
export MARKATA_GO_DRAFT_MODE=true           # Include drafts
export MARKATA_GO_FUTURE_POSTS=true         # Include scheduled posts

# Feed settings
export MARKATA_GO_FEED_DEFAULTS_ITEMS_PER_PAGE=20
```

### CI/CD Integration

#### GitHub Actions

```yaml
name: Deploy

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'
          cache: true

      - name: Install markata-go
        run: go install github.com/example/markata-go/cmd/markata-go@latest

      - name: Build (Production)
        if: github.ref == 'refs/heads/main'
        run: markata-go build --clean
        env:
          MARKATA_GO_URL: https://example.com

      - name: Build (Preview)
        if: github.event_name == 'pull_request'
        run: markata-go build --clean
        env:
          MARKATA_GO_URL: https://preview-${{ github.event.number }}.example.com

      - name: Deploy
        if: github.ref == 'refs/heads/main'
        uses: actions/deploy-pages@v4
```

#### Netlify

```toml
# netlify.toml
[build]
  command = "go install github.com/example/markata-go/cmd/markata-go@latest && markata-go build --clean"
  publish = "public"

[context.production]
  environment = { MARKATA_GO_URL = "https://example.com" }

[context.deploy-preview]
  environment = { MARKATA_GO_URL = "" }  # Uses deploy preview URL
```

#### Docker

```dockerfile
# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /build

ARG MARKATA_GO_URL=https://example.com
ENV MARKATA_GO_URL=$MARKATA_GO_URL

RUN go install github.com/example/markata-go/cmd/markata-go@latest
COPY . .
RUN markata-go build --clean

# Production stage
FROM nginx:alpine
COPY --from=builder /build/public /usr/share/nginx/html
```

Build for different environments:

```bash
# Production
docker build --build-arg MARKATA_GO_URL=https://example.com -t site:prod .

# Staging
docker build --build-arg MARKATA_GO_URL=https://staging.example.com -t site:staging .
```

---

## Advanced Templates

Master template inheritance, conditionals, and reusable components.

**Related guide:** [[templates-guide|Templates Guide]]

### Template Inheritance Patterns

Create a flexible base template with multiple extension points:

```html
{# templates/base.html #}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">

    {# SEO block - can be fully overridden #}
    {% block seo %}
    <title>{% block title %}{{ config.Title }}{% endblock %}</title>
    <meta name="description" content="{% block description %}{{ config.Description }}{% endblock %}">
    {% endblock %}

    {# Open Graph - extendable #}
    {% block opengraph %}
    <meta property="og:title" content="{% block og_title %}{{ config.Title }}{% endblock %}">
    <meta property="og:type" content="{% block og_type %}website{% endblock %}">
    <meta property="og:url" content="{{ config.URL }}{% block og_url %}/{% endblock %}">
    {% endblock %}

    {# Styles - extendable #}
    {% block styles %}
    <link rel="stylesheet" href="/css/main.css">
    {% endblock %}

    {# Extra head content #}
    {% block head %}{% endblock %}
</head>
<body class="{% block body_class %}{% endblock %}">
    {% block skip_link %}
    <a href="#main" class="skip-link">Skip to content</a>
    {% endblock %}

    {% block header %}
    {% include "partials/header.html" %}
    {% endblock %}

    {% block main %}
    <main id="main">
        {% block content %}{% endblock %}
    </main>
    {% endblock %}

    {% block footer %}
    {% include "partials/footer.html" %}
    {% endblock %}

    {% block scripts %}
    <script src="/js/main.js"></script>
    {% endblock %}
</body>
</html>
```

### Conditional Layouts

Use different layouts based on post metadata:

```html
{# templates/post.html #}
{% extends "base.html" %}

{% block body_class %}
post
{% if post.Extra.featured %}featured{% endif %}
{% if post.Extra.wide %}wide-layout{% endif %}
{% endblock %}

{% block content %}
<article class="{% if post.Extra.layout %}{{ post.Extra.layout }}{% else %}standard{% endif %}">

    {# Hero section for featured posts #}
    {% if post.Extra.hero_image %}
    <div class="hero" style="background-image: url('{{ post.Extra.hero_image }}')">
        <h1>{{ post.Title }}</h1>
    </div>
    {% else %}
    <header>
        <h1>{{ post.Title }}</h1>
    </header>
    {% endif %}

    {# Sidebar for posts with TOC #}
    {% if post.Extra.toc %}
    <div class="with-sidebar">
        <aside class="toc">
            {{ post.Extra.toc_html|safe }}
        </aside>
        <div class="content">
            {{ body|safe }}
        </div>
    </div>
    {% else %}
    <div class="content">
        {{ body|safe }}
    </div>
    {% endif %}

    {# Series navigation #}
    {% if post.Extra.series %}
    {% include "partials/series-nav.html" %}
    {% endif %}

</article>
{% endblock %}
```

### Reusable Components

Create composable partial templates:

```html
{# templates/partials/card.html #}
{#
  Expects: post object
  Optional: show_image, show_tags, show_excerpt
#}
<article class="card {% if post.Extra.featured %}card--featured{% endif %}">
    {% if show_image|default_if_none:true and post.Extra.cover_image %}
    <a href="{{ post.Href }}" class="card__image">
        <img src="{{ post.Extra.cover_image }}"
             alt="{{ post.Title }}"
             loading="lazy">
    </a>
    {% endif %}

    <div class="card__body">
        <h3 class="card__title">
            <a href="{{ post.Href }}">{{ post.Title }}</a>
        </h3>

        {% if show_excerpt|default_if_none:true and post.Description %}
        <p class="card__excerpt">{{ post.Description|truncate:120 }}</p>
        {% endif %}

        <footer class="card__meta">
            {% if post.Date %}
            <time datetime="{{ post.Date|atom_date }}">
                {{ post.Date|date_format:"Jan 2, 2006" }}
            </time>
            {% endif %}

            {% if post.Extra.reading_time %}
            <span class="reading-time">{{ post.Extra.reading_time }}</span>
            {% endif %}
        </footer>

        {% if show_tags|default_if_none:false and post.Tags %}
        <ul class="card__tags">
            {% for tag in post.Tags|slice:":3" %}
            <li><a href="/tags/{{ tag|slugify }}/">{{ tag }}</a></li>
            {% endfor %}
        </ul>
        {% endif %}
    </div>
</article>
```

Use the component with different configurations:

```html
{# Full card with image and tags #}
{% with show_image=true, show_tags=true %}
{% include "partials/card.html" %}
{% endwith %}

{# Minimal card without image #}
{% with show_image=false %}
{% include "partials/card.html" %}
{% endwith %}

{# Text-only card #}
{% with show_image=false, show_excerpt=false %}
{% include "partials/card.html" %}
{% endwith %}
```

---

## Wikilinks and Internal Linking

Use `[[wikilink]]` syntax for easy cross-referencing between posts.

### Basic Wikilink Syntax

Link to other posts using their slug:

```markdown
---
title: "Go Interfaces"
---

# Understanding Go Interfaces

Interfaces in Go are implicit. See [[go-structs]] for how structs work.

For more on type systems, check out [[go-type-system|the type system guide]].

Related: [[go-generics]], [[go-methods]]
```

**Syntax variants:**

| Syntax | Output |
|--------|--------|
| `[[slug]]` | Link with auto-title |
| `[[slug|Custom Text]]` | Link with custom text |
| `[[slug#section]]` | Link to specific section |
| `[[slug#section|Text]]` | Section link with text |

### Cross-Referencing Posts

Create bidirectional links for a wiki-like experience:

```markdown
---
title: "Go Structs"
slug: "go-structs"
---

# Go Structs

Structs are the primary way to create custom types in Go.

See also:
- [[go-interfaces]] - How interfaces work with structs
- [[go-methods]] - Adding methods to structs
- [[go-embedding]] - Struct embedding patterns
```

### Broken Link Detection

markata-go validates wikilinks during build and warns about broken links:

```
$ markata-go build
Warning: Broken wikilink in posts/go-structs.md: [[go-embedding]] (post not found)
Warning: Broken wikilink in posts/intro.md: [[getting-started#setup]] (section not found)
```

Enable strict mode to fail builds on broken links:

```toml
[markata-go]
strict_wikilinks = true    # Fail build on broken links
```

### Backlinks

The link collector provides backlinks (posts that link to the current post):

```html
{# In post.html template #}
{% if post.inlinks %}
<aside class="backlinks">
    <h3>Linked from</h3>
    <ul>
    {% for link in post.inlinks %}
        <li><a href="{{ link.source_post.href }}">{{ link.source_post.title }}</a></li>
    {% endfor %}
    </ul>
</aside>
{% endif %}
```

---

## Table of Contents Generation

Automatically generate navigation for long-form content.

### Enable TOC Generation

Add `toc: true` to your frontmatter:

```yaml
---
title: "Complete Go Guide"
toc: true
toc_depth: 3       # Include h1, h2, h3 (default: 2)
toc_min_items: 3   # Only show TOC if 3+ headings (default: 2)
---

# Introduction

## Getting Started

### Installation

### Configuration

## Core Concepts

### Variables

### Functions

## Advanced Topics
```

### TOC Output

The TOC is available as `post.Extra.toc_html`:

```html
<nav class="toc">
    <h2>Contents</h2>
    <ol>
        <li><a href="#introduction">Introduction</a></li>
        <li>
            <a href="#getting-started">Getting Started</a>
            <ol>
                <li><a href="#installation">Installation</a></li>
                <li><a href="#configuration">Configuration</a></li>
            </ol>
        </li>
        <li>
            <a href="#core-concepts">Core Concepts</a>
            <ol>
                <li><a href="#variables">Variables</a></li>
                <li><a href="#functions">Functions</a></li>
            </ol>
        </li>
        <li><a href="#advanced-topics">Advanced Topics</a></li>
    </ol>
</nav>
```

### Template Integration

Include the TOC in your post template:

```html
{% block content %}
<article class="post">
    <header>
        <h1>{{ post.Title }}</h1>
    </header>

    {% if post.Extra.toc_html %}
    <aside class="toc-sidebar">
        {{ post.Extra.toc_html|safe }}
    </aside>
    {% endif %}

    <div class="post-content">
        {{ body|safe }}
    </div>
</article>
{% endblock %}
```

### Styling the TOC

```css
.toc {
    position: sticky;
    top: 2rem;
    max-height: calc(100vh - 4rem);
    overflow-y: auto;
    padding: 1rem;
    background: var(--bg-secondary);
    border-radius: 8px;
}

.toc ol {
    list-style: none;
    padding-left: 0;
}

.toc ol ol {
    padding-left: 1rem;
    margin-top: 0.5rem;
}

.toc a {
    display: block;
    padding: 0.25rem 0;
    color: var(--text-secondary);
    text-decoration: none;
}

.toc a:hover,
.toc a.active {
    color: var(--text-primary);
}

/* Highlight current section with JavaScript */
.toc a.active {
    font-weight: 600;
    color: var(--accent);
}
```

---

## Admonitions and Callouts

Create visually distinct callout blocks for notes, warnings, tips, and more.

### Admonition Syntax

Use fenced blocks with type indicators:

```markdown
:::note
This is a note. Use it for additional information.
:::

:::tip
Pro tip! This helps users work more efficiently.
:::

:::warning
Be careful! This action has consequences.
:::

:::danger
Critical warning! This could cause data loss.
:::

:::info
Informational callout for context.
:::
```

### Admonition with Custom Title

```markdown
:::note Custom Title
This note has a custom title instead of "Note".
:::

:::warning Before You Continue
Make sure you've completed the prerequisites.
:::
```

### Rendered Output

Admonitions render as styled HTML:

```html
<div class="admonition admonition-note">
    <div class="admonition-title">Note</div>
    <div class="admonition-content">
        <p>This is a note. Use it for additional information.</p>
    </div>
</div>
```

### Available Types

| Type | Icon | Use Case |
|------|------|----------|
| `note` | Info icon | Additional context |
| `tip` | Lightbulb | Helpful suggestions |
| `info` | Info circle | Background information |
| `warning` | Warning triangle | Potential issues |
| `danger` | X circle | Critical warnings |
| `success` | Checkmark | Confirmations |
| `example` | Code icon | Code examples |
| `quote` | Quote icon | Quotations |
| `abstract` | Document | Summaries |

### Custom Admonition Types

Define custom types in your configuration:

```toml
[markata-go.admonitions]
# Define a custom "exercise" type
[markata-go.admonitions.exercise]
icon = "pencil"
color = "#6366f1"    # Indigo
```

Use in markdown:

```markdown
:::exercise Practice Problem
Write a function that reverses a string.
:::
```

### Styling Admonitions

```css
.admonition {
    margin: 1.5rem 0;
    padding: 1rem;
    border-left: 4px solid;
    border-radius: 4px;
    background: var(--bg-secondary);
}

.admonition-title {
    font-weight: 600;
    margin-bottom: 0.5rem;
    display: flex;
    align-items: center;
    gap: 0.5rem;
}

.admonition-note { border-color: #3b82f6; }
.admonition-tip { border-color: #22c55e; }
.admonition-warning { border-color: #f59e0b; }
.admonition-danger { border-color: #ef4444; }
.admonition-info { border-color: #06b6d4; }

.admonition-note .admonition-title { color: #3b82f6; }
.admonition-tip .admonition-title { color: #22c55e; }
.admonition-warning .admonition-title { color: #f59e0b; }
.admonition-danger .admonition-title { color: #ef4444; }
.admonition-info .admonition-title { color: #06b6d4; }
```

### Collapsible Admonitions

Make admonitions collapsible with the `collapse` modifier:

```markdown
:::tip collapse Expand for Pro Tips
These tips are hidden by default to keep the content focused.

- Tip 1: Use keyboard shortcuts
- Tip 2: Enable caching
- Tip 3: Run in parallel
:::
```

Renders as:

```html
<details class="admonition admonition-tip">
    <summary class="admonition-title">Expand for Pro Tips</summary>
    <div class="admonition-content">
        <p>These tips are hidden by default...</p>
    </div>
</details>
```

---

## See Also

- [[feeds-guide|Feeds Guide]] - Complete feed system documentation
- [[templates-guide|Templates Guide]] - Template system reference
- [[configuration-guide|Configuration Guide]] - Full configuration options
- [[plugin-development|Plugin Development]] - Create custom plugins
- [[dynamic-content|Dynamic Content]] - Jinja-in-Markdown guide
- [[deployment-guide|Deployment]] - Deployment and CI/CD
- [[frontmatter-guide|Frontmatter]] - Frontmatter field reference
- [[syndication-feeds|Syndication Feeds]] - RSS, Atom, JSON feeds
