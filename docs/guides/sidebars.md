---
title: "Sidebar Navigation"
description: "Complete guide to configuring sidebar navigation with path-based, feed-linked, and multi-feed sidebars"
date: 2024-01-23
published: true
slug: /docs/guides/sidebars/
tags:
  - documentation
  - layout
  - navigation
  - sidebar
---


markata-go provides a powerful sidebar system that adapts to your site structure. You can have different sidebars for different sections, auto-generate navigation from feeds, or combine multiple feeds into a unified sidebar.

Large sites reuse sidebar link projections within each build, avoiding repeated
URL processing for the same post and feed. This requires no configuration:
titles, links, feed selection, and the current-page highlight remain page-correct,
and the temporary cache is reset for the next template-render stage.

## Quick Start

The simplest way to add a sidebar is to define navigation items directly:

```toml
[markata-go.sidebar]
enabled = true
title = "Navigation"

[[markata-go.sidebar.nav]]
title = "Home"
href = "/"

[[markata-go.sidebar.nav]]
title = "Documentation"
href = "/docs/"

[[markata-go.sidebar.nav]]
title = "Blog"
href = "/blog/"
```

---

## Sidebar Modes

markata-go supports three main sidebar modes:

| Mode | Use Case | Configuration |
|------|----------|---------------|
| **Manual** | Hand-crafted navigation | `sidebar.nav` items |
| **Path-based** | Different sidebars for different URL paths | `sidebar.paths` mapping |
| **Multi-feed** | Combined navigation from multiple feeds | `sidebar.multi_feed = true` |

---

## Basic Configuration

### Sidebar Settings

```toml
[markata-go.sidebar]
enabled = true           # Show the sidebar (default: true)
position = "left"        # Position: "left" or "right" (default: "left")
width = "280px"          # Sidebar width (default: "280px")
collapsible = true       # Allow collapsing on mobile (default: true)
default_open = true      # Start expanded (default: true)
title = "Navigation"     # Optional sidebar header
```

### Manual Navigation Items

Define navigation items with nested structure:

```toml
[[markata-go.sidebar.nav]]
title = "Getting Started"
href = "/getting-started/"

[[markata-go.sidebar.nav]]
title = "Guides"
children = [
    { title = "Configuration", href = "/guides/configuration/" },
    { title = "Themes", href = "/guides/themes/" },
    { title = "Templates", href = "/guides/templates/" },
]

[[markata-go.sidebar.nav]]
title = "API Reference"
href = "/reference/"
children = [
    { title = "CLI", href = "/reference/cli/" },
    { title = "Plugins", href = "/reference/plugins/" },
]
```

### Navigation Item Properties

| Property | Type | Description |
|----------|------|-------------|
| `title` | string | Display text for the link |
| `href` | string | URL path (optional for parent items) |
| `children` | array | Nested navigation items |
| `external` | bool | Opens in new tab if true |
| `icon` | string | Optional icon name |

---

## Path-Based Sidebars

Show different sidebars based on the current page's URL path. This is ideal for sites with distinct sections like docs, blog, and tutorials.

### Configuration

```toml
[markata-go.sidebar]
enabled = true

# Sidebar for /docs/ pages
[markata-go.sidebar.paths."/docs/"]
title = "Documentation"
items = [
    { title = "Introduction", href = "/docs/" },
    { title = "Installation", href = "/docs/install/" },
    { title = "Configuration", href = "/docs/config/" },
]

# Sidebar for /blog/ pages
[markata-go.sidebar.paths."/blog/"]
title = "Blog"
items = [
    { title = "All Posts", href = "/blog/" },
    { title = "Archive", href = "/blog/archive/" },
    { title = "Tags", href = "/blog/tags/" },
]

# Sidebar for /tutorials/ pages
[markata-go.sidebar.paths."/tutorials/"]
title = "Tutorials"
feed = "tutorials"  # Auto-generate from feed
```

### Path Matching

Paths are matched using **longest prefix wins**:

| Page URL | Matching Path | Sidebar |
|----------|---------------|---------|
| `/docs/install/` | `/docs/` | Documentation |
| `/docs/v2/getting-started/` | `/docs/v2/` | Docs v2 (if defined) |
| `/blog/my-post/` | `/blog/` | Blog |
| `/about/` | (none) | Default sidebar |

### Path Sidebar Options

```toml
[markata-go.sidebar.paths."/docs/"]
title = "Documentation"       # Section title
position = "left"            # Override default position
collapsible = true           # Override collapsible setting

# Option 1: Manual items
items = [...]

# Option 2: Link to a feed
feed = "docs"

# Option 3: Auto-generate from directory
auto_generate = { directory = "docs", order_by = "nav_order" }
```

---

## Feed-Linked Sidebars

Link a sidebar to a feed to automatically generate navigation from the feed's posts.

### Enable Sidebar on a Feed

```toml
[[markata-go.feeds]]
slug = "docs"
filter = "path.startswith('docs/')"
title = "Documentation"
sidebar = true              # Enable sidebar generation
sidebar_title = "Docs"      # Override title in sidebar
sidebar_order = 1           # Position in multi-feed sidebars
sidebar_group_by = "category"  # Group by frontmatter field
```

### Link Path to Feed

```toml
[markata-go.sidebar.paths."/docs/"]
feed = "docs"  # Uses the "docs" feed for navigation
```

### Grouping by Frontmatter

When `sidebar_group_by` is set, posts are grouped by that frontmatter field:

```toml
[[markata-go.feeds]]
slug = "tutorials"
filter = "'tutorial' in tags"
sidebar = true
sidebar_group_by = "difficulty"
```

With posts having frontmatter like:

```yaml
---
title: "Your First Site"
difficulty: "Beginner"
---
```

The sidebar will show:

```
Tutorials
├── Beginner
│   ├── Your First Site
│   └── Basic Configuration
├── Intermediate
│   ├── Custom Templates
│   └── Plugin Development
└── Advanced
    └── Performance Optimization
```

---

## Auto-Generated Sidebars

Automatically generate sidebar navigation from your directory structure.

### From Directory

```toml
[markata-go.sidebar]
auto_generate = { directory = "docs" }
```

### Full Auto-Generate Options

```toml
[markata-go.sidebar.auto_generate]
directory = "docs"           # Source directory
order_by = "nav_order"       # Sort: "title", "date", "nav_order", "filename"
reverse = false              # Reverse sort order
max_depth = 3                # Limit directory depth (0 = unlimited)
exclude = ["drafts/*", "_*"] # Glob patterns to exclude
```

### Ordering with nav_order

Add `nav_order` to your frontmatter to control sidebar order:

```yaml
---
title: "Getting Started"
nav_order: 1
---
```

```yaml
---
title: "Configuration"
nav_order: 2
---
```

Lower numbers appear first. Posts without `nav_order` default to 999.

---

## Multi-Feed Sidebars

Combine multiple feeds into a single sidebar with collapsible sections.

### Basic Multi-Feed

```toml
[markata-go.sidebar]
multi_feed = true
feeds = ["docs", "guides", "reference"]
```

This creates a sidebar like:

```
├── Documentation
│   ├── Getting Started
│   └── Installation
├── Guides
│   ├── Configuration
│   └── Themes
└── Reference
    ├── CLI
    └── Plugins
```

### Detailed Section Configuration

For more control, use `feed_sections`:

```toml
[markata-go.sidebar]
multi_feed = true

[[markata-go.sidebar.feed_sections]]
feed = "docs"
title = "Documentation"    # Override feed title
collapsed = false          # Start expanded
max_items = 10             # Limit items shown

[[markata-go.sidebar.feed_sections]]
feed = "guides"
title = "Guides"
collapsed = true           # Start collapsed

[[markata-go.sidebar.feed_sections]]
feed = "reference"
title = "API Reference"
collapsed = true
max_items = 5
```

### Feed Sidebar Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `sidebar` | bool | `false` | Enable this feed for sidebar |
| `sidebar_title` | string | feed title | Display title in sidebar |
| `sidebar_order` | int | 0 | Sort order in multi-feed mode |
| `sidebar_group_by` | string | - | Group posts by frontmatter field |

---

## Template Integration

Sidebar data is available in templates via the context:

```html
{% if sidebar_items %}
<nav class="sidebar">
    {% if sidebar_title %}
    <h2>{{ sidebar_title }}</h2>
    {% endif %}

    <ul>
    {% for item in sidebar_items %}
        <li>
            {% if item.href %}
            <a href="{{ item.href }}"
               {% if item.href == post.href %}class="active"{% endif %}>
                {{ item.title }}
            </a>
            {% else %}
            <span>{{ item.title }}</span>
            {% endif %}

            {% if item.children %}
            <ul>
                {% for child in item.children %}
                <li>
                    <a href="{{ child.href }}"
                       {% if child.href == post.href %}class="active"{% endif %}>
                        {{ child.title }}
                    </a>
                </li>
                {% endfor %}
            </ul>
            {% endif %}
        </li>
    {% endfor %}
    </ul>
</nav>
{% endif %}
```

### Available Template Variables

| Variable | Type | Description |
|----------|------|-------------|
| `sidebar_items` | array | Navigation items for current page |
| `sidebar_title` | string | Sidebar title/header |
| `config.sidebar` | object | Full sidebar configuration |

---

## Per-Page Overrides

Override sidebar settings in page frontmatter:

```yaml
---
title: "My Page"
sidebar: false          # Hide sidebar on this page
---
```

```yaml
---
title: "Special Page"
layout: landing         # Landing layout has no sidebar
---
```

---

## Examples

### Documentation Site

```toml
[markata-go.sidebar]
enabled = true

[markata-go.sidebar.paths."/docs/"]
title = "Documentation"
auto_generate = { directory = "docs", order_by = "nav_order" }

[markata-go.sidebar.paths."/api/"]
title = "API Reference"
auto_generate = { directory = "api", order_by = "title" }
```

### Blog with Categories

```toml
[[markata-go.feeds]]
slug = "posts"
filter = "path.startswith('blog/')"
sidebar = true
sidebar_title = "Categories"
sidebar_group_by = "category"

[markata-go.sidebar.paths."/blog/"]
feed = "posts"
```

### Multi-Section Site

```toml
[markata-go.sidebar]
multi_feed = true

[[markata-go.feeds]]
slug = "docs"
filter = "path.startswith('docs/')"
sidebar = true
sidebar_order = 1

[[markata-go.feeds]]
slug = "tutorials"
filter = "'tutorial' in tags"
sidebar = true
sidebar_order = 2

[[markata-go.feeds]]
slug = "blog"
filter = "path.startswith('blog/')"
sidebar = true
sidebar_order = 3
```

### Versioned Documentation

```toml
[markata-go.sidebar.paths."/v1/"]
title = "v1 Documentation"
auto_generate = { directory = "v1" }

[markata-go.sidebar.paths."/v2/"]
title = "v2 Documentation"
auto_generate = { directory = "v2" }
```

---

## Best Practices

1. **Use path-based sidebars** for sites with distinct sections (docs, blog, guides)

2. **Use feed-linked sidebars** when you want navigation to automatically update as you add content

3. **Use multi-feed sidebars** for comprehensive navigation across your entire site

4. **Add nav_order** to frontmatter for predictable ordering instead of relying on alphabetical sort

5. **Keep sidebars focused** - don't try to show everything; use the header nav for top-level navigation

6. **Test mobile behavior** - ensure sidebars collapse properly on small screens

---

## Troubleshooting

### Sidebar not showing

1. Check that `sidebar.enabled = true`
2. Verify your layout supports sidebars (`docs` layout, not `landing` or `bare`)
3. Check frontmatter doesn't have `sidebar: false`

### Wrong sidebar showing

1. Check path matching - longer prefixes take priority
2. Verify feed slugs match in both feed definition and sidebar config

### Items not appearing

1. Check feed filter is matching expected posts
2. Verify posts have `published: true`
3. Check exclude patterns aren't filtering desired content

### Incorrect order

1. Add `nav_order` frontmatter to control explicit ordering
2. Check `order_by` setting matches your intent
3. Verify `reverse` is set correctly

---

## Related

- [[configuration|Configuration Guide]] - Full configuration reference
- [[feeds|Feeds Guide]] - Creating and configuring feeds
- [[templates|Templates Guide]] - Customizing sidebar templates
- [[themes|Themes Guide]] - Styling your sidebar
