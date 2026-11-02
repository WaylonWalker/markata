---
title: "Syndication Feeds"
description: "Guide to generating Sitemap, RSS, Atom, and JSON Feed outputs for your site"
date: 2024-01-15
published: true
tags:
  - documentation
  - feeds
  - rss
  - atom
---


This guide covers generating syndication feeds for your markata-go site: **Sitemap**, **RSS**, **Atom**, and **JSON Feed**. Each format serves different consumers - search engines, feed readers, and JavaScript applications.

## Quick Reference

| Format | File | MIME Type | Use Case |
|--------|------|-----------|----------|
| Sitemap | `sitemap.xml` | `application/xml` | Search engine indexing |
| RSS 2.0 | `rss.xml` | `application/rss+xml` | Traditional feed readers |
| Atom | `atom.xml` | `application/atom+xml` | Modern feed readers |
| JSON Feed | `feed.json` | `application/feed+json` | JavaScript apps, APIs |

## Prerequisites

Before generating syndication feeds, ensure your `markata-go.toml` has the required site metadata:

```toml
[markata-go]
title = "My Blog"                      # Required for all feeds
description = "A blog about things"    # Used in RSS/Atom/JSON
url = "https://example.com"            # Required for absolute URLs
author = "Jane Doe"                    # Used in Atom/JSON feeds
```

---

## 1. Generate a Sitemap

The sitemap plugin automatically generates a `sitemap.xml` file containing all published posts and feed index pages.

### Enable the Sitemap

The sitemap is generated automatically when the `sitemap` plugin runs. No additional configuration is required for basic usage.

**Output:** `/sitemap.xml` in your output directory

### Example Sitemap Output

```xml
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
    <url>
        <loc>https://example.com/</loc>
        <lastmod>2024-01-15</lastmod>
        <changefreq>daily</changefreq>
        <priority>1.0</priority>
    </url>
    <url>
        <loc>https://example.com/my-first-post/</loc>
        <lastmod>2024-01-15</lastmod>
        <changefreq>weekly</changefreq>
        <priority>0.8</priority>
    </url>
    <url>
        <loc>https://example.com/blog/</loc>
        <lastmod>2024-01-15</lastmod>
        <changefreq>weekly</changefreq>
        <priority>0.6</priority>
    </url>
</urlset>
```

### Customize Sitemap Values via Frontmatter

Override `changefreq` and `priority` per-post in your markdown frontmatter:

```yaml
---
title: "Important Announcement"
date: 2024-01-15
published: true
sitemap:
  changefreq: daily      # Options: always, hourly, daily, weekly, monthly, yearly, never
  priority: 0.9          # Range: 0.0 to 1.0
---
```

### Default Sitemap Values

| Page Type | changefreq | priority |
|-----------|------------|----------|
| Home page | `daily` | `1.0` |
| Posts | `weekly` | `0.8` |
| Feed index pages | `weekly` | `0.6` |

### Exclude Posts from Sitemap

Use `skip: true` in frontmatter to exclude a post from the sitemap:

```yaml
---
title: "Draft Post"
skip: true
---
```

---

## 2. Generate RSS Feeds

RSS 2.0 feeds are the most widely supported format for feed readers.

### Enable RSS for a Feed

```toml
[[markata-go.feeds]]
slug = "blog"
title = "My Blog"
filter = "published == True"
sort = "date"
reverse = true

[markata-go.feeds.formats]
html = true
rss = true                  # Enable RSS output
```

**Output:** `/blog/rss.xml`

### Example RSS Output

```xml
<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
    <title>My Blog</title>
    <link>https://example.com</link>
    <description>A blog about things</description>
    <language>en-us</language>
    <lastBuildDate>Mon, 15 Jan 2024 12:00:00 +0000</lastBuildDate>
    <atom:link href="https://example.com/blog/rss.xml" rel="self" type="application/rss+xml"/>

    <item>
      <title>My First Post</title>
      <link>https://example.com/my-first-post/</link>
      <description>This is the post description...</description>
      <pubDate>Mon, 15 Jan 2024 12:00:00 +0000</pubDate>
      <guid isPermaLink="true">https://example.com/my-first-post/</guid>
    </item>
  </channel>
</rss>
```

### RSS Configuration Options

Control RSS behavior via feed defaults:

```toml
[markata-go.feed_defaults.syndication]
max_items = 20              # Maximum items in RSS feed (default: 20)
include_content = true      # Include full content or just description
```

### RSS with Full Content

To include full post content instead of just the description:

```toml
[markata-go.feed_defaults.syndication]
include_content = true
```

When `include_content = true`, the `<description>` element contains the rendered HTML content (truncated to 500 characters). When `false`, it uses the post's `description` frontmatter field.

---

## 3. Generate Atom Feeds

Atom feeds (RFC 4287) provide richer metadata than RSS and are well-supported by modern feed readers.

### Enable Atom for a Feed

```toml
[[markata-go.feeds]]
slug = "blog"
title = "My Blog"
filter = "published == True"
sort = "date"
reverse = true

[markata-go.feeds.formats]
html = true
rss = true
atom = true                 # Enable Atom output
```

**Output:** `/blog/atom.xml`

### Example Atom Output

```xml
<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title>My Blog</title>
  <id>https://example.com/blog/atom.xml</id>
  <updated>2024-01-15T12:00:00Z</updated>
  <link href="https://example.com" rel="alternate" type="text/html"/>
  <link href="https://example.com/blog/atom.xml" rel="self" type="application/atom+xml"/>
  <author>
    <name>Jane Doe</name>
  </author>

  <entry>
    <title>My First Post</title>
    <id>https://example.com/my-first-post/</id>
    <updated>2024-01-15T12:00:00Z</updated>
    <published>2024-01-15T12:00:00Z</published>
    <link href="https://example.com/my-first-post/" rel="alternate" type="text/html"/>
    <summary type="text">This is the post description...</summary>
    <content type="html">&lt;p&gt;Full post content...&lt;/p&gt;</content>
  </entry>
</feed>
```

### Atom vs RSS: Key Differences

| Feature | RSS 2.0 | Atom |
|---------|---------|------|
| Author info | Limited | Full (name, email, URI) |
| Content types | Text only | Text, HTML, XHTML |
| Update tracking | `lastBuildDate` | `updated` per-entry |
| Self-reference | Via extension | Built-in `rel="self"` |
| ID format | URL (guid) | IRI (more flexible) |

### When to Use Atom

- You need author information (name, email)
- You want both summary and full content
- You need per-entry update timestamps
- Your audience uses modern feed readers

### When to Use RSS

- Maximum compatibility with older readers
- Simpler format requirements
- Podcast feeds (RSS is the standard)

**Recommendation:** Enable both RSS and Atom - let your users choose.

---

## 4. Generate JSON Feeds

JSON Feed (version 1.1) is a modern feed format designed for easy consumption by JavaScript applications.

### Enable JSON Feed

```toml
[[markata-go.feeds]]
slug = "blog"
title = "My Blog"
filter = "published == True"
sort = "date"
reverse = true

[markata-go.feeds.formats]
html = true
rss = true
atom = true
json = true                 # Enable JSON Feed output
```

**Output:** `/blog/feed.json`

### Example JSON Feed Output

```json
{
  "version": "https://jsonfeed.org/version/1.1",
  "title": "My Blog",
  "home_page_url": "https://example.com",
  "feed_url": "https://example.com/blog/feed.json",
  "description": "A blog about things",
  "language": "en",
  "authors": [
    {
      "name": "Jane Doe"
    }
  ],
  "items": [
    {
      "id": "https://example.com/my-first-post/",
      "url": "https://example.com/my-first-post/",
      "title": "My First Post",
      "content_html": "<p>Full post content...</p>",
      "content_text": "Full post content...",
      "summary": "This is the post description...",
      "date_published": "2024-01-15T12:00:00Z",
      "date_modified": "2024-01-15T12:00:00Z",
      "tags": ["go", "tutorial"]
    }
  ]
}
```

### JSON Feed Use Cases

1. **JavaScript consumption** - Parse with `fetch()` and `JSON.parse()`
2. **API endpoints** - Serve as a simple read-only API
3. **Static search indexes** - Power client-side search
4. **Cross-origin requests** - JSON has better CORS support than XML

### Using JSON Feed in JavaScript

```javascript
fetch('/blog/feed.json')
  .then(response => response.json())
  .then(feed => {
    feed.items.forEach(item => {
      console.log(item.title, item.url);
    });
  });
```

---

## 5. Combining Multiple Formats

The real power of markata-go's feed system is generating all formats from a single definition.

### Complete Multi-Format Feed

```toml
[markata-go]
title = "My Tech Blog"
description = "Articles about Go, Python, and web development"
url = "https://techblog.example.com"
author = "Jane Doe"
output_dir = "public"

# Global feed defaults
[markata-go.feed_defaults]
items_per_page = 10

[markata-go.feed_defaults.syndication]
max_items = 20
include_content = true

# Main blog feed with ALL formats
[[markata-go.feeds]]
slug = "blog"
title = "Tech Blog"
description = "Latest articles about programming"
filter = "published == True"
sort = "date"
reverse = true

[markata-go.feeds.formats]
html = true                 # /blog/index.html, /blog/page/2/index.html
rss = true                  # /blog/rss.xml
atom = true                 # /blog/atom.xml
json = true                 # /blog/feed.json
```

### Output Structure

```
public/
  sitemap.xml                    # Global sitemap (auto-generated)
  blog/
    index.html                   # HTML page 1
    page/
      2/
        index.html               # HTML page 2
    rss.xml                      # RSS 2.0 feed
    atom.xml                     # Atom feed
    feed.json                    # JSON Feed
```

### Multiple Feeds with Different Formats

```toml
# Home page - HTML only
[[markata-go.feeds]]
slug = ""
title = "Home"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 5

[markata-go.feeds.formats]
html = true
rss = false
atom = false
json = false

# Blog - All syndication formats
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
json = true

# API endpoint - JSON only
[[markata-go.feeds]]
slug = "api/posts"
title = "Posts API"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 0              # No pagination - all posts

[markata-go.feeds.formats]
html = false
rss = false
atom = false
json = true
```

---

## 6. Adding Feed Discovery Links

Help feed readers and browsers automatically discover your feeds by adding `<link>` tags to your HTML templates.

### Template with Feed Discovery

Add these links to your base template's `<head>` section:

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ title }}</title>

    <!-- RSS Feed Discovery -->
    <link rel="alternate"
          type="application/rss+xml"
          title="{{ site.title }} (RSS)"
          href="{{ site.url }}/blog/rss.xml">

    <!-- Atom Feed Discovery -->
    <link rel="alternate"
          type="application/atom+xml"
          title="{{ site.title }} (Atom)"
          href="{{ site.url }}/blog/atom.xml">

    <!-- JSON Feed Discovery -->
    <link rel="alternate"
          type="application/feed+json"
          title="{{ site.title }} (JSON)"
          href="{{ site.url }}/blog/feed.json">
</head>
<body>
    <!-- ... -->
</body>
</html>
```

### Feed Links Per Section

If you have multiple feeds (blog, tutorials, etc.), add discovery links for each:

```html
<head>
    <!-- Blog feeds -->
    <link rel="alternate" type="application/rss+xml"
          title="Blog (RSS)" href="/blog/rss.xml">
    <link rel="alternate" type="application/atom+xml"
          title="Blog (Atom)" href="/blog/atom.xml">

    <!-- Tutorial feeds -->
    <link rel="alternate" type="application/rss+xml"
          title="Tutorials (RSS)" href="/tutorials/rss.xml">
    <link rel="alternate" type="application/atom+xml"
          title="Tutorials (Atom)" href="/tutorials/atom.xml">
</head>
```

### Dynamic Feed Discovery in Templates

For feed-specific pages, use template variables:

```html
{% if feed %}
    {% if feed.formats.rss %}
    <link rel="alternate" type="application/rss+xml"
          title="{{ feed.title }} (RSS)"
          href="{{ feed.href }}rss.xml">
    {% endif %}

    {% if feed.formats.atom %}
    <link rel="alternate" type="application/atom+xml"
          title="{{ feed.title }} (Atom)"
          href="{{ feed.href }}atom.xml">
    {% endif %}

    {% if feed.formats.json %}
    <link rel="alternate" type="application/feed+json"
          title="{{ feed.title }} (JSON)"
          href="{{ feed.href }}feed.json">
    {% endif %}
{% endif %}
```

### Visible Feed Links

Also provide visible links for users who want to manually subscribe:

```html
<footer>
    <h3>Subscribe</h3>
    <ul class="feed-links">
        <li><a href="/blog/rss.xml">RSS Feed</a></li>
        <li><a href="/blog/atom.xml">Atom Feed</a></li>
        <li><a href="/blog/feed.json">JSON Feed</a></li>
    </ul>
</footer>
```

---

## Complete Example Configuration

Here's a full `markata-go.toml` with all syndication features configured:

```toml
[markata-go]
# Site metadata (required for syndication feeds)
title = "My Tech Blog"
description = "Articles about Go, Python, and web development"
url = "https://techblog.example.com"
author = "Jane Doe"

# Build settings
output_dir = "public"
templates_dir = "templates"

# Content discovery
[markata-go.glob]
patterns = ["posts/**/*.md", "pages/*.md"]

# Markdown extensions
[markata-go.markdown]
extensions = ["tables", "strikethrough", "tasklist"]

# =============================================================================
# FEED DEFAULTS
# =============================================================================

[markata-go.feed_defaults]
items_per_page = 10
orphan_threshold = 3

[markata-go.feed_defaults.formats]
html = true
rss = true
atom = true
json = false

[markata-go.feed_defaults.syndication]
max_items = 20                  # Items in RSS/Atom/JSON feeds
include_content = true          # Include full content in feeds

# =============================================================================
# FEEDS
# =============================================================================

# Home page - HTML only, no syndication
[[markata-go.feeds]]
slug = ""
title = "Recent Posts"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 5

[markata-go.feeds.formats]
html = true
rss = false
atom = false

# Main blog - Full syndication
[[markata-go.feeds]]
slug = "blog"
title = "Blog"
description = "All blog posts"
filter = "published == True"
sort = "date"
reverse = true

[markata-go.feeds.formats]
html = true
rss = true
atom = true
json = true

# Tutorials - RSS and Atom only
[[markata-go.feeds]]
slug = "tutorials"
title = "Tutorials"
description = "Step-by-step programming guides"
filter = "published == True and 'tutorial' in tags"
sort = "date"
reverse = true

[markata-go.feeds.formats]
html = true
rss = true
atom = true
json = false

# JSON API endpoint
[[markata-go.feeds]]
slug = "api/posts"
title = "Posts API"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 0              # All posts, no pagination

[markata-go.feeds.formats]
html = false
rss = false
atom = false
json = true
```

### Generated Output

```
public/
  index.html                     # Home page
  sitemap.xml                    # Global sitemap

  blog/
    index.html                   # Blog page 1
    page/2/index.html            # Blog page 2
    rss.xml                      # RSS feed
    atom.xml                     # Atom feed
    feed.json                    # JSON feed

  tutorials/
    index.html                   # Tutorials page 1
    rss.xml                      # RSS feed
    atom.xml                     # Atom feed

  api/
    posts/
      feed.json                  # JSON API (all posts)
```

---

## Troubleshooting

### Feed URLs show `https://example.com`

Ensure you've set the `url` field in your config:

```toml
[markata-go]
url = "https://your-actual-domain.com"
```

### Feeds are empty

Check that posts have `published = true` (or your filter matches):

```yaml
---
title: "My Post"
published: true    # Required for most filters
---
```

### Author missing in Atom/JSON feeds

Set the `author` field:

```toml
[markata-go]
author = "Your Name"
```

### Feed not generating

Ensure the format is enabled:

```toml
[markata-go.feeds.formats]
rss = true     # Must be explicitly true
```

---

## See Also

- [[feeds-guide|Feeds Guide]] - Complete feed system documentation
- [[configuration-guide|Configuration Guide]] - Full configuration reference
- [[templates-guide|Templates Guide]] - Template system and filters
