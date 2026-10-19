---
title: "Blogroll & Reader Guide"
description: "Create a blogroll of blogs you follow and a reader page showing their latest posts"
date: 2024-01-15
published: true
slug: /docs/guides/blogroll/
tags:
  - documentation
  - blogroll
  - rss
  - feeds
---

# Blogroll & Reader

The blogroll plugin lets you curate a list of blogs you follow and automatically fetch their RSS/Atom feeds to create:

- **`/blogroll/`** - A directory of blogs you follow, organized by category
- **`/reader/`** - A river-of-news style page showing latest posts from all feeds

This is perfect for sharing your reading list, building community connections, or creating a personal feed reader built into your site.

## Quick Start

Add this to your `markata-go.toml`:

```toml
[markata-go.blogroll]
enabled = true

[[markata-go.blogroll.feeds]]
url = "https://simonwillison.net/atom/everything/"
title = "Simon Willison"
category = "Technology"
```

Run `markata-go build` and you'll have:
- `/blogroll/` - Lists Simon Willison's blog
- `/reader/` - Shows his latest posts

## Manage subscriptions in your browser

Start the local development server with `markata-go serve`, then open
[`/_markata/subscriptions`](http://localhost:8000/_markata/subscriptions).
The Serve Control Center and local preview pages also link to this editor.
Paste a public website URL and select **Find feed** to discover its advertised
RSS/Atom/JSON Feed; review its title, category, type and tags, then subscribe.
You can search, edit or remove existing subscriptions without editing TOML.

The editor writes only the **effective owning TOML config file**, including
when blogroll subscriptions live in an included config. It preserves other
sections and comments. When a later included/merge config replaces the feed
list, the editor never silently writes to a shadowed earlier file. Global
configs, merge overrides and non-TOML owners are read-only; edit those source
files directly. The editor is only served to local loopback sessions, not
included in production site output.

## Configuration

### Basic Settings

```toml
[markata-go.blogroll]
enabled = true                    # Enable the blogroll plugin
cache_dir = "cache/blogroll"      # Where to cache fetched feeds
cache_duration = "24h"            # How long to cache (default: 24 hours)
timeout = 30                      # HTTP request timeout in seconds
concurrent_requests = 5           # Max parallel feed fetches
max_entries_per_feed = 50         # Global default entries per feed
```

### Pagination Settings

The reader page supports pagination to handle large numbers of entries:

```toml
[markata-go.blogroll]
enabled = true
items_per_page = 50               # Entries per page (default: 50)
orphan_threshold = 3              # Min entries for separate page (default: 3)
pagination_type = "manual"        # "manual", "htmx", or "js"
```

**Pagination Types:**

| Type | Description |
|------|-------------|
| `manual` | Traditional page links with full page reloads |
| `htmx` | Seamless AJAX-based navigation using HTMX |
| `js` | Client-side JavaScript pagination |

The paginated reader generates:
```
/reader/              # Page 1
/reader/page/2/       # Page 2
/reader/page/3/       # Page 3
/reader/page/2/partial/   # HTMX partial for page 2
```

### Adding Feeds

Add feeds using the `[[markata-go.blogroll.feeds]]` array:

```toml
[[markata-go.blogroll.feeds]]
url = "https://example.com/feed.xml"    # Required: RSS/Atom feed URL
title = "Example Blog"                   # Optional: display name (auto-fetched if not set)
description = "A great blog about stuff" # Optional: short description
category = "Technology"                  # Optional: groups feeds on blogroll page
type = "written"                         # Optional: written, video, or podcast
tags = ["python", "web"]                 # Optional: additional labels
site_url = "https://example.com"         # Optional: main website URL
image_url = "https://example.com/logo.png" # Optional: logo/icon
active = true                            # Optional: set false to disable without removing
max_entries = 50                         # Optional: override global max_entries_per_feed
```

`markata-go blogroll add <feed-url>` now uses site metadata for the suggested title, description, and tags shown in the prompt. This works better for sources like YouTube, where the feed itself often lacks useful channel-level description and tag metadata.

If your site composes multiple config files with `markata-go.include`, `markata-go blogroll add` tries to append the new feed to the file that already contains the blogroll section instead of rewriting the root config file.

You can also pass common YouTube channel inputs and let markata-go resolve them to the underlying channel feed automatically:

```bash
markata-go blogroll add yt:devtoolsfm
markata-go blogroll add https://www.youtube.com/@devtoolsfm
markata-go blogroll add https://www.youtube.com/watch?v=aEkpFeJoKvk
```

For a single YouTube video URL, markata-go resolves the owning channel and subscribes to that channel's feed.

### Per-Feed Entry Limits

Override the global `max_entries_per_feed` for individual feeds:

```toml
[markata-go.blogroll]
enabled = true
max_entries_per_feed = 50         # Global default

[[markata-go.blogroll.feeds]]
url = "https://prolific-blogger.com/feed.xml"
title = "Prolific Blogger"
max_entries = 100                 # Override: this site posts frequently

[[markata-go.blogroll.feeds]]
url = "https://micro.blog/user.xml"
title = "Micro Blog"
max_entries = 200                 # Override: many small posts

[[markata-go.blogroll.feeds]]
url = "https://infrequent-poster.com/feed.xml"
title = "Infrequent Poster"
max_entries = 10                  # Override: rarely posts, save cache space

[[markata-go.blogroll.feeds]]
url = "https://normal-blog.com/feed.xml"
title = "Normal Blog"
# Uses global default: 50 entries
```

### Feed Configuration Reference

| Field | Required | Default | Description |
|-------|----------|---------|-------------|
| `url` | Yes | - | RSS or Atom feed URL |
| `title` | No | Auto-fetched | Display name for the feed |
| `description` | No | Auto-fetched | Short description |
| `category` | No | "Uncategorized" | Groups feeds together |
| `tags` | No | `[]` | Additional labels for filtering |
| `site_url` | No | Auto-fetched | Main website URL |
| `image_url` | No | Auto-fetched | Logo or icon URL |
| `active` | No | `true` | Set to `false` to disable |
| `max_entries` | No | Global default | Override max entries for this feed |

### Custom Templates

Override the default templates:

```toml
[markata-go.blogroll.templates]
blogroll = "blogroll.html"    # Template for /blogroll/ page
reader = "reader.html"        # Template for /reader/ page
```

## Example: Building a Reading List

Here's a complete example with multiple feeds organized by category:

```toml
[markata-go.blogroll]
enabled = true
cache_duration = "2h"
max_entries_per_feed = 25

# =============================================================================
# TECHNOLOGY
# =============================================================================

[[markata-go.blogroll.feeds]]
url = "https://simonwillison.net/atom/everything/"
title = "Simon Willison"
description = "Creator of Datasette, Django co-creator, AI/LLM enthusiast"
category = "Technology"
tags = ["python", "ai", "llm", "sqlite"]
site_url = "https://simonwillison.net"

[[markata-go.blogroll.feeds]]
url = "https://jvns.ca/atom.xml"
title = "Julia Evans"
description = "Making hard things easy to understand"
category = "Technology"
tags = ["linux", "networking", "zines"]
site_url = "https://jvns.ca"

[[markata-go.blogroll.feeds]]
url = "https://danluu.com/atom.xml"
title = "Dan Luu"
description = "Deep dives into computer systems"
category = "Technology"
tags = ["systems", "performance"]
site_url = "https://danluu.com"

[[markata-go.blogroll.feeds]]
url = "https://blog.codinghorror.com/rss/"
title = "Coding Horror"
description = "Jeff Atwood on programming and human factors"
category = "Technology"
tags = ["programming", "software"]
site_url = "https://blog.codinghorror.com"

# =============================================================================
# DESIGN
# =============================================================================

[[markata-go.blogroll.feeds]]
url = "https://alistapart.com/main/feed/"
title = "A List Apart"
description = "For people who make websites"
category = "Design"
tags = ["web", "ux", "accessibility"]
site_url = "https://alistapart.com"

[[markata-go.blogroll.feeds]]
url = "https://css-tricks.com/feed/"
title = "CSS-Tricks"
description = "Tips, tricks, and techniques on using CSS"
category = "Design"
tags = ["css", "frontend"]
site_url = "https://css-tricks.com"

# =============================================================================
# PERSONAL
# =============================================================================

[[markata-go.blogroll.feeds]]
url = "https://austinkleon.com/feed/"
title = "Austin Kleon"
description = "Writer and artist"
category = "Personal"
tags = ["creativity", "writing"]
site_url = "https://austinkleon.com"
```

## Adding Simon Willison

Simon Willison is a prolific blogger known for creating Datasette, co-creating Django, and writing extensively about AI/LLMs. His feed URL is:

```
https://simonwillison.net/atom/everything/
```

Add him to your blogroll:

```toml
[[markata-go.blogroll.feeds]]
url = "https://simonwillison.net/atom/everything/"
title = "Simon Willison"
description = "Creator of Datasette, Django co-creator, AI/LLM enthusiast"
category = "Technology"
tags = ["python", "ai", "llm", "sqlite", "datasette"]
site_url = "https://simonwillison.net"
```

Simon also has topic-specific feeds if you want to subscribe to specific content:
- `https://simonwillison.net/atom/entries/` - Blog entries only
- `https://simonwillison.net/atom/links/` - Links/bookmarks only

## Finding Feed URLs

Most blogs have RSS/Atom feeds. Here's how to find them:

1. **Look for feed icons** - Usually in the header, footer, or sidebar
2. **Check common paths:**
   - `/feed/`
   - `/rss/`
   - `/atom.xml`
   - `/feed.xml`
   - `/rss.xml`
   - `/index.xml`
3. **View page source** - Search for `application/rss+xml` or `application/atom+xml`
4. **Use browser extensions** - Feed discovery extensions can help

### Common Feed URL Patterns

| Platform | Feed URL Pattern |
|----------|-----------------|
| WordPress | `/feed/` or `/feed/rss/` |
| Ghost | `/rss/` |
| Substack | `/feed` |
| Medium | `/feed` |
| Jekyll | `/feed.xml` |
| Hugo | `/index.xml` |
| dev.to | `/feed` |
| GitHub Releases | `/releases.atom` |

## Generated Pages

### Blogroll Page (`/blogroll/`)

The blogroll page lists all feeds grouped by category:

```
/blogroll/
  index.html
```

**Default layout:**
- Header with title and feed count
- Feeds grouped by category
- Each feed shows: title, description, post count
- Links to the original site

### Reader Page (`/reader/`)

The reader page shows the latest posts from all feeds in reverse chronological order, grouped by day, with pagination:

```
/reader/
  index.html              # Page 1
  partial/
    index.html            # Page 1 partial (HTMX)
  page/
    2/
      index.html          # Page 2
      partial/
        index.html        # Page 2 partial (HTMX)
    3/
      index.html
      partial/
        index.html
```

**Default layout:**
- Editorial header with page stats
- Type switch links between the full reader and any typed readers
- Entries grouped by day with a large date rail
- Two-column entry grid on desktop, single column on mobile
- Each entry shows: source favicon, source name, date, title, description, optional image
- Links to the original article
- Hacker News and Lobsters discussion URLs are normalized to the outbound article URL when one is available. Lobsters RSS feeds that already link directly to articles retain their discussion permalinks from the GUID or comments field. The original discussions remain available as "Found on HN" or "Found on Lobsters", alongside the feed source label
- Pagination navigation (when more than one page)
- Offscreen reader cards after the first four in each day group may skip layout and paint in supporting browsers via CSS `content-visibility`; the full HTML, links, and no-JavaScript navigation remain available. Preview images and source icons load lazily and decode asynchronously.

Ghost-hosted PNG, JPEG, and WebP previews with an existing `storage.ghost.io`
width URL include responsive candidates. The browser selects a width for the
card and screen density, including when JavaScript is disabled. The original
image remains the fallback; other providers and GIF/SVG previews keep their
existing URLs. Custom reader templates can use `entry.image_srcset` with
`sizes="100vw"` on lazy images.

### Community aggregators: Lobsters, Lemmy and Reddit

The reader preserves the original discussion permalink alongside a linked external article. For Lemmy community RSS feeds, `<link>` normally points to the article, and a `<guid>` or `<comments>` permalink can point to the Lemmy discussion. For Reddit's Atom feeds, posts remain linked to Reddit unless an explicit `[link]` anchor in the entry content provides an external article URL; text-only posts remain on Reddit. Unknown or blocked feeds retain the normal fallback behavior, and arbitrary links in a post body are not promoted to articles.

Example subscriptions:

```toml
[[markata-go.blogroll.feeds]]
url = "https://lemmy.world/feeds/c/selfhosted.xml"
title = "Lemmy: Selfhosted"

[[markata-go.blogroll.feeds]]
url = "https://www.reddit.com/r/selfhosted/.rss"
title = "Reddit: Selfhosted"
```

**Reddit deprecation:** Reddit [announced](https://www.reddit.com/r/modnews/comments/1wubgvt/) that its native RSS feeds will stop working on **November 13, 2026**. The Reddit example is temporary and should not be relied on for new long-term subscriptions; no unauthenticated replacement is promised. Lemmy provides supported native RSS feeds. Feed refreshes may also be rate-limited by individual communities.

If you assign feed types, markata-go also generates typed reader pages:

```text
/reader/           # all feeds
/reader/written/   # written feeds only
/reader/video/     # video feeds only
/reader/podcast/   # podcast feeds only
```

**Deterministic ordering:**
1. Most recent publish/update date first
2. Feed URL (ascending)
3. Entry ID (ascending)
4. Title (ascending)

## Custom Templates

Create custom templates for full control over the appearance.

### Blogroll Template

Create `templates/blogroll.html`:

```html
{% extends "base.html" %}

{% block content %}
<main class="blogroll">
  <h1>{{ title }}</h1>
  <p class="subtitle">{{ feed_count }} blogs I follow</p>

  {% for category in categories %}
  <section class="category" id="{{ category.Slug }}">
    <h2>{{ category.Name }}</h2>
    <div class="feed-grid">
      {% for feed in category.Feeds %}
      <article class="feed-card">
        {% if feed.ImageURL %}
        <img src="{{ feed.ImageURL }}" alt="{{ feed.Title }}" class="feed-icon">
        {% endif %}
        <h3>
          {% if feed.SiteURL %}
          <a href="{{ feed.SiteURL }}" target="_blank" rel="noopener">{{ feed.Title }}</a>
          {% else %}
          {{ feed.Title }}
          {% endif %}
        </h3>
        {% if feed.Description %}
        <p class="description">{{ feed.Description }}</p>
        {% endif %}
        <div class="meta">
          <span class="post-count">{{ feed.EntryCount }} posts</span>
          <a href="{{ feed.FeedURL }}" class="feed-link" title="RSS Feed">
            <svg><!-- RSS icon --></svg>
          </a>
        </div>
      </article>
      {% endfor %}
    </div>
  </section>
  {% endfor %}
</main>
{% endblock %}
```

### Reader Template

Create `templates/reader.html`:

```html
{% extends "base.html" %}

{% block content %}
<div class="reader-page">
  <h1>{{ title }}</h1>
  <p class="reader-subtitle">Latest posts from blogs I follow</p>

  <div class="reader-stream posts-list">
    {% for day in day_groups %}
    <section class="reader-day">
      <aside class="reader-day-rail">
        <time{% if day.date_iso %} datetime="{{ day.date_iso }}"{% endif %}>
          <span class="reader-day-weekday">{{ day.weekday }}</span>
          <span class="reader-day-title">{{ day.title }}</span>
          {% if day.year %}<span class="reader-day-year">{{ day.year }}</span>{% endif %}
        </time>
        <span class="reader-day-count">{{ day.count }} stories</span>
      </aside>

      <div class="reader-day-entries">
        {% for entry in day.entries %}
        <article class="reader-entry{% if entry.image_url %} has-image{% endif %}">
          <div class="reader-entry-meta-row">
            <a href="{{ entry.feed_url }}" target="_blank" rel="noopener noreferrer" class="reader-entry-source-link">
              {% if entry.source_icon_url %}
              <img src="{{ entry.source_icon_url }}" alt="" class="reader-entry-source-icon" width="16" height="16" loading="lazy">
              {% else %}
              <span class="reader-entry-source-icon reader-entry-source-icon--fallback" aria-hidden="true"></span>
              {% endif %}
              <span class="reader-entry-source">{{ entry.feed_title }}</span>
            </a>
            {% if entry.published_datetime %}
            <time datetime="{{ entry.published_datetime }}">{{ entry.published_label }}</time>
            {% endif %}
          </div>

          {% if entry.original_url %}
          <div class="reader-entry-found-on">Found on HN: <a href="{{ entry.original_url }}" target="_blank" rel="noopener noreferrer">discussion</a></div>
          {% endif %}

          <h2 class="reader-entry-title">
            <a href="{{ entry.url }}" target="_blank" rel="noopener noreferrer">{{ entry.title }}<span class="visually-hidden">(opens in new tab)</span></a>
          </h2>

          {% if entry.description %}
          <p class="reader-entry-description">{{ entry.description|striptags|truncate:200 }}</p>
          {% endif %}

          {% if entry.image_url %}
          <a href="{{ entry.url }}" target="_blank" rel="noopener noreferrer" class="reader-entry-image-link">
            <img src="{{ entry.image_url }}" alt="" class="reader-entry-image" width="960" height="540" loading="lazy">
          </a>
          {% endif %}
        </article>
        {% endfor %}
      </div>
    </section>
    {% endfor %}
  </div>
</div>
{% endblock %}
```

## Template Variables

### Blogroll Template Variables

| Variable | Type | Description |
|----------|------|-------------|
| `title` | string | Page title ("Blogroll") |
| `description` | string | Page description |
| `feeds` | []ExternalFeed | All feeds |
| `categories` | []BlogrollCategory | Feeds grouped by category |
| `feed_count` | int | Total number of feeds |

### Reader Template Variables

| Variable | Type | Description |
|----------|------|-------------|
| `title` | string | Page title ("Reader") |
| `description` | string | Page description |
| `entries` | []ExternalEntry | Entries for current page (newest first) |
| `day_groups` | []ReaderDayGroup | Entries grouped by day for the curated layout |
| `entry_count` | int | Total number of entries across all pages |
| `page` | ReaderPage | Pagination information |
| `pagination_type` | string | Pagination type ("manual", "htmx", "js") |
| `original_url` | string | Original discussion link when a feed entry was normalized |

### ReaderPage Fields

| Field | Type | Description |
|-------|------|-------------|
| `number` | int | Current page number (1-indexed) |
| `has_prev` | bool | True if previous page exists |
| `has_next` | bool | True if next page exists |
| `prev_url` | string | URL of previous page |
| `next_url` | string | URL of next page |
| `total_pages` | int | Total number of pages |
| `total_items` | int | Total number of entries |
| `items_per_page` | int | Entries per page |
| `page_urls` | []string | URLs for all pages |
| `pagination_type` | string | Pagination type |

### ReaderDayGroup Fields

| Field | Type | Description |
|-------|------|-------------|
| `date_iso` | string | ISO-8601 day timestamp for the rail |
| `weekday` | string | Short weekday label |
| `title` | string | Short month/day label |
| `year` | string | Year label |
| `count` | int | Entries in the day group |
| `entries` | []map[string]interface{} | Entry maps for the group |

### ExternalFeed Fields

| Field | Type | Description |
|-------|------|-------------|
| `Title` | string | Feed title |
| `Description` | string | Feed description |
| `SiteURL` | string | Main website URL |
| `FeedURL` | string | RSS/Atom feed URL |
| `ImageURL` | string | Feed logo/icon |
| `Category` | string | Feed category |
| `Tags` | []string | Feed tags |
| `EntryCount` | int | Number of entries |
| `Entries` | []ExternalEntry | Feed entries |
| `LastFetched` | *time.Time | When feed was last fetched |
| `LastUpdated` | *time.Time | Feed's last update date |
| `Error` | string | Error message if fetch failed |

### ExternalEntry Fields

| Field | Type | Description |
|-------|------|-------------|
| `Title` | string | Entry title |
| `URL` | string | Link to full article |
| `Description` | string | Summary or excerpt |
| `Content` | string | Full content (HTML) |
| `Author` | string | Entry author |
| `Published` | *time.Time | Publication date |
| `Updated` | *time.Time | Last update date |
| `Categories` | []string | Entry categories/tags |
| `ImageURL` | string | Featured image |
| `ReadingTime` | int | Estimated reading time (minutes) |
| `FeedURL` | string | Source feed URL |
| `FeedTitle` | string | Source feed title |

Reader entry maps also include `source_icon_url`, `published_datetime`, `published_label`, and `date_key` for the redesigned layout.

### BlogrollCategory Fields

| Field | Type | Description |
|-------|------|-------------|
| `Name` | string | Category name |
| `Slug` | string | URL-safe identifier |
| `Feeds` | []ExternalFeed | Feeds in this category |

## Caching

The blogroll plugin caches fetched feeds to avoid hitting external servers on every build.

### Cache Configuration

```toml
[markata-go.blogroll]
cache_dir = "cache/blogroll"    # Cache directory
cache_duration = "24h"          # How long to cache feeds
refresh_on_build = false         # Reuse cache during builds; refresh with `reader update`
```

### Cache Behavior

1. `markata-go reader update` is the explicit cache refresh command for blogroll and reader feeds
2. When `refresh_on_build = false`, normal builds reuse cached feed data and do not make remote refresh requests
3. When `refresh_on_build = true`, builds use cached feeds until `cache_duration` expires, then refresh as needed
4. Expiration is based on the cached feed's recorded fetch time, not file mtimes
5. If a refresh fails, the stale cached feed is reused instead of failing open
6. Explicit config overrides like `title`, `description`, `site_url`, `image_url`, `handle`, aliases, category, and tags are reapplied even when cached feed data is reused
7. Delete `cache/blogroll/` to force a fresh fetch

### Cache Duration Examples

```toml
cache_duration = "30m"    # 30 minutes
cache_duration = "1h"     # 1 hour
cache_duration = "6h"     # 6 hours
cache_duration = "24h"    # 1 day (default)
cache_duration = "168h"   # 1 week
```

**Recommendations:**
- Development: `"30m"` to `"1h"` - Good balance without hammering feeds
- Production: `"24h"` to `"168h"` - Prioritize build speed and stable caches
- High-traffic: `"24h"` or more - Reduce external requests

## Error Handling

When a feed fails to fetch, the plugin:

1. Records the error in `feed.Error`
2. Continues processing other feeds
3. Includes the feed in the blogroll (with error indicator)
4. Uses cached data if available

### Handling Errors in Templates

```html
{% for feed in feeds %}
<article class="feed-card {% if feed.Error %}feed-error{% endif %}">
  <h3>{{ feed.Title }}</h3>
  {% if feed.Error %}
  <p class="error">Unable to fetch: {{ feed.Error }}</p>
  {% else %}
  <p>{{ feed.EntryCount }} posts</p>
  {% endif %}
</article>
{% endfor %}
```

## Screenshot Service Integration

Many RSS feed entries don't include featured images. The `fallback_image_service` option lets you automatically generate preview images using screenshot services.

### Configuration

Set a URL template with `{url}` as a placeholder for the entry URL:

```toml
[blogroll]
enabled = true
fallback_image_service = "https://shots.so/s/{url}"
```

When an entry doesn't have an `image_url`, markata-go will:
1. URL-encode the entry's URL
2. Replace `{url}` in the template with the encoded URL
3. Use the resulting URL as the fallback image

### Popular Screenshot Services

#### shots.so (Free)

[shots.so](https://shots.so/) provides free website screenshots with no API key required:

```toml
fallback_image_service = "https://shots.so/s/{url}"
```

Features:
- No registration or API key needed
- Automatic caching of screenshots
- Customizable viewport and device type

For custom parameters:

```toml
fallback_image_service = "https://shots.so/s/{url}?width=1200&height=630"
```

#### API Flash

[API Flash](https://apiflash.com/) offers high-quality screenshots with customization:

```toml
fallback_image_service = "https://api.apiflash.com/v1/urltoimage?access_key=YOUR_ACCESS_KEY&url={url}&width=1200&height=630"
```

Features:
- Requires API key (free tier available)
- Extensive customization options
- Reliable caching and CDN delivery

Get your access key from [apiflash.com](https://apiflash.com/).

#### ScreenshotOne

[ScreenshotOne](https://screenshotone.com/) provides API-based screenshot generation:

```toml
fallback_image_service = "https://api.screenshotone.com/take?access_key=YOUR_ACCESS_KEY&url={url}&viewport_width=1200&viewport_height=630&image_quality=80"
```

Features:
- API key required (free tier available)
- Advanced rendering options
- GDPR compliant

Get your access key from [screenshotone.com](https://screenshotone.com/).

### URL Encoding

The `{url}` placeholder is automatically URL-encoded, so you don't need to manually encode it. For example:

```
Entry URL: https://example.com/blog/post?id=123&ref=home
Encoded:   https%3A%2F%2Fexample.com%2Fblog%2Fpost%3Fid%3D123%26ref%3Dhome
Result:    https://shots.so/s/https%3A%2F%2Fexample.com%2Fblog%2Fpost%3Fid%3D123%26ref%3Dhome
```

### Rate Limits and Caching

**Important considerations:**

1. **Rate Limits** - Most screenshot services have rate limits on their free tiers. If you have many feeds with many entries, you may hit these limits during builds.

2. **Build Time** - Screenshot services generate images on-demand, which can slow down your first build. Subsequent builds are faster as services cache screenshots.

3. **Service Caching** - Screenshots are typically cached by the service, so the same URL will return quickly on subsequent requests.

4. **Local Caching** - Consider the blogroll plugin's `cache_duration` setting to reduce build frequency and screenshot service requests:

```toml
[blogroll]
cache_duration = "6h"    # Reduces rebuild frequency
fallback_image_service = "https://shots.so/s/{url}"
```

5. **Cost Considerations** - For high-traffic sites or large blogrolls, consider:
   - Using services with generous free tiers (like shots.so)
   - Limiting `max_entries_per_feed` to reduce screenshot requests
   - Increasing `cache_duration` to minimize rebuilds
   - Upgrading to paid tiers for higher rate limits

### Example Configuration

```toml
[blogroll]
enabled = true
cache_duration = "6h"              # Cache for 6 hours to reduce rebuilds
max_entries_per_feed = 20          # Limit entries to reduce screenshot requests
fallback_image_service = "https://shots.so/s/{url}?width=1200&height=630"

[[blogroll.feeds]]
url = "https://simonwillison.net/atom/everything/"
title = "Simon Willison"
category = "Technology"
```

With this configuration, the Reader uses a preview hierarchy for entries without article images:

1. article image from the feed entry
2. source image from the feed avatar/logo
3. screenshot fallback from `fallback_image_service`
4. branded source tile when no image is available

This keeps the Reader more consistent and reduces the visual noise from arbitrary webpage screenshots.

YouTube feeds also use the per-video thumbnail automatically, so video channels show the actual video art in `/reader/` instead of a generic source tile or screenshot.

## Performance Tips

### Optimize Build Times

1. **Increase cache duration** for production builds
2. **Limit `max_entries_per_feed`** if you only need recent posts
3. **Reduce `concurrent_requests`** if you're hitting rate limits
4. **Disable inactive feeds** with `active = false` instead of removing them

### Example Production Config

```toml
[markata-go.blogroll]
enabled = true
cache_dir = "cache/blogroll"
cache_duration = "6h"           # Cache for 6 hours
timeout = 15                    # Shorter timeout
concurrent_requests = 3         # Be nice to servers
max_entries_per_feed = 20       # Only recent posts
```

## Configuration Reference

### Full Configuration

```toml
[markata-go.blogroll]
# Enable/disable the entire feature
enabled = true

# Cache settings
cache_dir = "cache/blogroll"
cache_duration = "1h"

# HTTP settings
timeout = 30
concurrent_requests = 5
max_entries_per_feed = 50

# Pagination settings
items_per_page = 50
orphan_threshold = 3
pagination_type = "manual"    # "manual", "htmx", or "js"

# Screenshot service for last-resort Reader previews
fallback_image_service = ""   # Used only after article and source images are exhausted

# Custom templates
[blogroll.templates]
blogroll = "blogroll.html"
reader = "reader.html"

# Feeds
[[markata-go.blogroll.feeds]]
url = "https://example.com/feed.xml"
title = "Example"
description = "Description"
category = "Category"
tags = ["tag1", "tag2"]
site_url = "https://example.com"
image_url = "https://example.com/logo.png"
active = true
max_entries = 50              # Override global max_entries_per_feed
```

### Configuration Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Enable blogroll plugin |
| `cache_dir` | string | `"cache/blogroll"` | Cache directory |
| `cache_duration` | string | `"24h"` | Cache TTL (Go duration) |
| `refresh_on_build` | bool | `false` | Allow builds to refresh remote blogroll feeds |
| `timeout` | int | `30` | HTTP timeout in seconds |
| `concurrent_requests` | int | `5` | Max parallel fetches |
| `max_entries_per_feed` | int | `50` | Global max entries per feed |
| `items_per_page` | int | `50` | Entries per reader page |
| `orphan_threshold` | int | `3` | Min entries for separate page |
| `pagination_type` | string | `"manual"` | Pagination style |
| `fallback_image_service` | string | `""` | Last-resort screenshot template for Reader previews |
| `feeds` | []Feed | `[]` | List of feeds |
| `templates.blogroll` | string | `"blogroll.html"` | Blogroll template |
| `templates.reader` | string | `"reader.html"` | Reader template |

---

## Fast Builds During Development

When working on large sites with many blogroll feeds, the best default is usually to keep blogroll pages enabled while separating refresh from normal builds.

### Keep Blogroll Pages, Skip Refresh During Builds

```toml
[markata-go.blogroll]
enabled = true
refresh_on_build = false
```

Then refresh on demand:

```bash
markata-go reader update
markata-go build
```

For one-off runs, you can also set:

```bash
MARKATA_GO_BLOGROLL_REFRESH_ON_BUILD=false markata-go build
```

If you want the smallest possible loop, you can still disable the plugin entirely:

### Disable Blogroll Via Environment Variable

```bash
# Build without fetching blogroll feeds
MARKATA_GO_BLOGROLL_ENABLED=false markata-go build

# Or combine with glob patterns to build just one post
MARKATA_GO_GLOB_PATTERNS="posts/my-draft.md" MARKATA_GO_BLOGROLL_ENABLED=false markata-go build
```

### Using a Local .env File

Create a `.env.local` file (and add it to `.gitignore`):

```bash
# .env.local
MARKATA_GO_BLOGROLL_ENABLED=false
```

Then source it before building:

```bash
source .env.local && markata-go build
```

### When to Use Fast Builds

**Use fast builds when:**
- Writing or editing posts (blogroll content doesn't affect your writing)
- Testing template changes
- Running local development server
- Debugging build issues

**Use full builds when:**
- Deploying to production
- Testing blogroll template changes
- Verifying complete site output

### Build Time Impact

Blogroll fetch times depend on:
- Number of feeds configured
- Feed response times
- `concurrent_requests` setting (default: 5 parallel)
- Cache hit rate

For a blogroll with 20 feeds, expect:
- **First build:** 5-15 seconds (fetching all feeds)
- **Cached build:** 0.5-2 seconds (using cached feeds)
- **Disabled:** 0 seconds (blogroll plugin skipped entirely)

---

## Next Steps

- [Feeds Guide](/docs/guides/feeds/) - Create feeds from your own content
- [Templates Guide](/docs/guides/templates/) - Customize blogroll appearance
- [Syndication Guide](/docs/guides/syndication-feeds/) - Share your own content via RSS

---

## See Also

- [Configuration Guide](/docs/guides/configuration/) - Full configuration reference
- [Themes Guide](/docs/guides/themes/) - Style your blogroll with themes
