---
title: "Series Guide"
description: "Create ordered, multi-part post collections for tutorials, courses, and sequential reading"
date: 2024-01-15
published: true
slug: /docs/guides/series/
tags:
  - documentation
  - series
  - feeds
---

## Series

Series are **ordered collections of posts meant to be read sequentially**. They are perfect for multi-part tutorials, article series, and course chapters.

A series automatically generates a feed with an index page, RSS/Atom feeds, sidebar navigation, prev/next links, and keyboard shortcuts -- all from a single frontmatter field.

## Quick Start

Add `series` to your post's frontmatter:

```yaml
---
title: "Building a REST API - Part 1: Setup"
series: "building-a-rest-api"
---
```

That's it. All posts with the same `series` value are automatically grouped, ordered by date, and given navigation.

## How It Works

```
  Frontmatter                 Auto-Generated Feed              Output
  +--------------+            +------------------+           +------------------+
  | series:      |            | slug: series/X   |           | /series/X/       |
  |   "rest-api" |---scan--->| type: series     |--render-->|   index.html     |
  | series_order:|            | sort: order/date |           |   rss.xml        |
  |   1          |            | sidebar: true    |           |   atom.xml       |
  +--------------+            +------------------+           |   feed.json      |
                                                              +------------------+
```

The series plugin scans all posts during the build's **Collect** stage, groups them by series name, sorts them, and injects feed configurations. The feeds plugin then generates all the output formats.

## Frontmatter Fields

### `series` (required)

The series identifier. Becomes the URL slug (prefixed with `series/`).

```yaml
series: "building-a-rest-api"
```

The value is slugified, so `"Building a REST API"` and `"building-a-rest-api"` are treated as the same series.

A post can belong to **at most one series**.

### `series_order` (optional)

Explicit position in the series (1-indexed).

```yaml
series: "building-a-rest-api"
series_order: 1
```

## Ordering

Posts within a series are ordered according to these rules:

1. If **any** post has `series_order`, sort by `series_order` ascending
2. Posts without `series_order` in a mixed series are placed after ordered posts, sorted by date
3. If **no** post has `series_order`, sort by date ascending (oldest first)
4. Ties are broken by file path for deterministic ordering

Series default to ascending order (oldest/lowest first) because they are meant to be read from beginning to end, unlike blog feeds which show newest first.

### Example: Date-based ordering

When no `series_order` is specified, posts are ordered by date:

```yaml
# part-1.md (2024-01-01) -> Position 1
---
title: "Setup"
series: "go-tutorial"
date: 2024-01-01
---

# part-2.md (2024-01-15) -> Position 2
---
title: "Basics"
series: "go-tutorial"
date: 2024-01-15
---

# part-3.md (2024-02-01) -> Position 3
---
title: "Advanced"
series: "go-tutorial"
date: 2024-02-01
---
```

### Example: Explicit ordering

Use `series_order` when you want control over the sequence regardless of dates:

```yaml
# intro.md -> Position 1 (despite being written last)
---
title: "Introduction"
series: "go-tutorial"
series_order: 1
date: 2024-03-01
---

# basics.md -> Position 2
---
title: "The Basics"
series: "go-tutorial"
series_order: 2
date: 2024-01-01
---
```

## Navigation

### Feed Sidebar

When a post belongs to a series, a sidebar automatically appears showing all posts as a table of contents. The current post is highlighted.

```
+---------------------------+
|  Building a REST API      |
|  --------------------     |
|  1. Project Setup      <  |  <- current post
|  2. Database Models       |
|  3. Handlers              |
|  4. Auth & Middleware     |
|  5. Testing & Deploy      |
|                           |
|  Part 1 of 5              |
+---------------------------+
```

The series sidebar **takes precedence** over any other feed sidebar. If a post has `series` frontmatter, the series sidebar is shown instead of feed-based sidebars.

### Prev/Next Links

Series posts get prev/next navigation that follows series order. The first post has no "previous" link, and the last post has no "next" link.

A "View All" link points to the series index page.

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `[` | Navigate to previous post in series |
| `]` | Navigate to next post in series |

These use the existing keyboard navigation infrastructure.

### Position Indicator

Posts show their position in the series. In templates:

```jinja2
Part {{ post.PrevNextContext.Position }} of {{ post.PrevNextContext.Total }}
```

## Output

For a series named `building-a-rest-api`, the following files are generated:

```
output/
  series/
    building-a-rest-api/
      index.html       # Series listing page
      simple/
        index.html     # Simple list view
      rss.xml          # RSS feed
      atom.xml         # Atom feed
      feed.json        # JSON feed
      sitemap.xml      # Sitemap
```

The index page lists all posts in order with their titles and descriptions.

## Configuration

Series work with zero configuration -- just add `series` to your frontmatter. For advanced use cases, you can customize behavior in `markata-go.toml`.

### Global Series Defaults

```toml
[markata-go.series]
slug_prefix = "series"        # URL prefix (default: "series")
auto_sidebar = true            # Auto-enable sidebar (default: true)

[markata-go.series.defaults]
items_per_page = 0             # No pagination (default: 0 = show all)
sidebar = true                 # Show sidebar on series posts (default: true)

[markata-go.series.defaults.formats]
html = true
simple_html = true
rss = true
atom = true
json = true
sitemap = true
```

### Per-Series Overrides

Override settings for a specific series:

```toml
[markata-go.series.overrides."building-a-rest-api"]
title = "Building a REST API with Go"
description = "A complete guide from zero to production"

[markata-go.series.overrides."building-a-rest-api".formats]
markdown = true                # Also generate markdown for this series
```

### Series vs Feed Defaults

Series have different defaults than regular feeds:

| Setting | Series Default | Feed Default |
|---------|---------------|--------------|
| `items_per_page` | 0 (no pagination) | 10 |
| `sidebar` | true | false |
| `type` | `"series"` | `"blog"` |
| `reverse` | false (ascending) | true (descending) |

## Template Variables

### Post Templates

When a post belongs to a series:

| Variable | Type | Description |
|----------|------|-------------|
| `post.Extra.series` | string | Series identifier |
| `post.Extra.series_order` | int | Explicit position (if set) |
| `post.Extra.series_slug` | string | Full series feed slug |
| `post.Extra.series_total` | int | Total posts in series |
| `post.Prev` | Post | Previous post in series |
| `post.Next` | Post | Next post in series |
| `post.PrevNextFeed` | string | Series feed slug |
| `post.PrevNextContext.FeedSlug` | string | Series feed slug |
| `post.PrevNextContext.FeedTitle` | string | Series title |
| `post.PrevNextContext.Position` | int | 1-indexed position |
| `post.PrevNextContext.Total` | int | Total posts in series |

### Example: Series Banner

```jinja2
{% if post.Extra.series %}
<div class="series-banner">
    <h3>{{ post.PrevNextContext.FeedTitle }}</h3>
    <p>Part {{ post.PrevNextContext.Position }} of {{ post.PrevNextContext.Total }}</p>
    <a href="/{{ post.PrevNextContext.FeedSlug }}/">View all parts</a>
</div>
{% endif %}
```

## Edge Cases

### Unpublished Posts

Posts with `published: false` are excluded from the series feed output but still count for ordering purposes. Prev/next links skip unpublished posts.

Private posts are treated as unpublished for series navigation and output.

### Single-Post Series

A series with one post is valid. No prev/next navigation is generated, but the sidebar and index page are still created.

### Empty Series

If all posts in a series are unpublished or in draft, no series feed is generated.

## See Also

- [Feeds Guide](/docs/guides/feeds/) - The underlying feed system
- [Frontmatter Guide](/docs/guides/frontmatter/) - All frontmatter fields
- [Sidebar Guide](/docs/guides/sidebars/) - Sidebar configuration
- [Keyboard Navigation](/docs/guides/keyboard-navigation/) - Navigation shortcuts
