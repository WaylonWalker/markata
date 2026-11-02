---
title: "WebMentions & Social Bridging"
description: "Connect your static site to the social web with WebMentions and cross-platform bridging"
date: 2026-01-24
published: true
tags:
  - webmentions
  - indieweb
  - social
  - bridging
---


markata-go supports the IndieWeb WebMention protocol for connecting your static site to the decentralized social web. This includes:

- **Outgoing WebMentions**: Automatically notify sites you link to
- **Social Media Bridging**: Receive likes, reposts, and replies from Twitter, Bluesky, Mastodon, and more

## Quick Start

Enable WebMentions in your `markata-go.toml`:

```toml
[markata-go.webmentions]
enabled = true

# For outgoing mentions
outgoing = true

# For social media bridging
[markata-go.webmentions.bridges]
enabled = true
bluesky = true
twitter = true
mastodon = true
```

## Outgoing WebMentions

When you link to external sites in your posts, markata-go can automatically send WebMentions to those sites (if they support the protocol).

### How It Works

1. During build, markata-go scans your posts for external links
2. For each link, it discovers the target site's WebMention endpoint
3. It sends a notification: "Hey, this URL linked to you!"
4. The target site can then display your link/comment

### Configuration

```toml
[markata-go.webmentions]
enabled = true
outgoing = true

# HTTP request timeout
timeout = "30s"

# Max concurrent requests
concurrent_requests = 5

# Cache directory to avoid re-sending
cache_dir = ".cache/webmentions"

# Custom User-Agent string
user_agent = "MySite/1.0 (WebMention; +https://mysite.com)"
```

### Caching

Sent WebMentions are cached to avoid re-sending on every build. The cache is stored in the configured `cache_dir` as JSON.

## Social Media Bridging

Social media bridging allows you to receive interactions from Twitter, Bluesky, Mastodon, and other platforms as WebMentions. This is powered by services like [Bridgy Fed](https://fed.brid.gy/).

### Supported Platforms

| Platform | Interactions | Detection |
|----------|-------------|-----------|
| Bluesky | Likes, reposts, replies | bsky.app URLs, Bridgy |
| Twitter/X | Likes, retweets, replies | twitter.com/x.com URLs |
| Mastodon | Favorites, boosts, replies | Fediverse patterns |
| GitHub | Stars, comments, issues | github.com URLs |
| Flickr | Favorites, comments | flickr.com URLs |

### Setting Up Bridging

1. **Sign up for webmention.io** to receive incoming mentions
2. **Connect Bridgy Fed** to your social accounts at [fed.brid.gy](https://fed.brid.gy/)
3. **Configure markata-go**:

```toml
[markata-go.webmentions]
enabled = true
webmention_io_token = "your_token_here"

[markata-go.webmentions.bridges]
enabled = true
bridgy_fediverse = true

# Enable specific platforms
bluesky = true
twitter = true
mastodon = true
github = true
flickr = false
```

### Platform Detection

markata-go automatically detects the source platform of incoming WebMentions by analyzing:

1. **URL patterns**: `brid.gy/publish/bluesky`, `bsky.app`, etc.
2. **Domain matching**: Known Mastodon instances, GitHub, etc.
3. **Content patterns**: Platform-specific indicators in mention content

Each detected mention is enriched with:
- **Platform name**: "bluesky", "twitter", "mastodon", etc.
- **Handle**: Platform-specific username (e.g., `@alice.bsky.social`)
- **Original URL**: Link back to the original interaction

### Filtering Bridged Mentions

You can filter incoming mentions by platform, interaction type, or content:

```toml
[markata-go.webmentions.bridges.filters]
# Only accept from these platforms
platforms = ["bluesky", "mastodon"]

# Only accept these interaction types
# Options: "like", "repost", "reply", "bookmark", "mention"
interaction_types = ["like", "repost", "reply"]

# Minimum content length for replies
min_content_length = 10

# Block specific domains
blocked_domains = ["spam-site.com", "bad-actor.net"]
```

## Template Integration

Display WebMentions in your templates:

```html
{% if post.webmentions %}
<section class="webmentions">
  <h3>Reactions</h3>

  {% for mention in post.webmentions %}
  <article class="webmention webmention--{{ mention.platform }}">
    <!-- Platform badge -->
    {% if mention.platform != "web" %}
    <span class="platform-badge">
      {% if mention.platform == "bluesky" %}Bluesky
      {% elif mention.platform == "twitter" %}Twitter
      {% elif mention.platform == "mastodon" %}Mastodon
      {% elif mention.platform == "github" %}GitHub
      {% endif %}
    </span>
    {% endif %}

    <!-- Author info -->
    <div class="author">
      {% if mention.author.photo %}
      <img src="{{ mention.author.photo }}" alt="{{ mention.author.name }}">
      {% endif %}
      <a href="{{ mention.author.url }}">
        {{ mention.handle | default(mention.author.name) }}
      </a>
    </div>

    <!-- Interaction type -->
    {% if mention.wm_property == "like-of" %}
      <p>Liked this post</p>
    {% elif mention.wm_property == "repost-of" %}
      <p>Reposted this post</p>
    {% elif mention.wm_property == "in-reply-to" %}
      <blockquote>{{ mention.content.text }}</blockquote>
    {% endif %}

    <!-- Link to original -->
    {% if mention.original_url %}
    <a href="{{ mention.original_url }}">View on {{ mention.platform }}</a>
    {% endif %}
  </article>
  {% endfor %}
</section>
{% endif %}
```

### Platform-Specific Styling

```css
/* Platform colors */
.platform-badge--bluesky { background: #0085ff; }
.platform-badge--twitter { background: #1da1f2; }
.platform-badge--mastodon { background: #6364ff; }
.platform-badge--github { background: #333; }

/* Platform borders */
.webmention--bluesky { border-left: 3px solid #0085ff; }
.webmention--twitter { border-left: 3px solid #1da1f2; }
.webmention--mastodon { border-left: 3px solid #6364ff; }
.webmention--github { border-left: 3px solid #333; }
```

## Fetching WebMentions

WebMentions are fetched from webmention.io and cached locally. The build process reads from the cache, so you need to periodically fetch fresh mentions.

### Fetching with CLI

```bash
# Fetch all webmentions for your site
markata-go webmentions fetch

# The mentions are saved to .cache/webmentions/received_mentions.json
```

### Environment Variable

You can set your webmention.io token via environment variable instead of config:

```bash
export WEBMENTION_IO_TOKEN="your_token_here"
markata-go webmentions fetch
```

### Automated Fetching

Set up a cron job or CI workflow to fetch mentions regularly:

```yaml
# .github/workflows/fetch-webmentions.yml
name: Fetch WebMentions
on:
  schedule:
    - cron: '0 */6 * * *'  # Every 6 hours
  workflow_dispatch:

jobs:
  fetch:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Fetch webmentions
        run: markata-go webmentions fetch
        env:
          WEBMENTION_IO_TOKEN: ${{ secrets.WEBMENTION_IO_TOKEN }}
      - name: Commit cache
        run: |
          git config user.name github-actions
          git config user.email github-actions@github.com
          git add .cache/webmentions/
          git diff --staged --quiet || git commit -m "chore: update webmention cache"
          git push
```

## Displaying Counts on Post Cards

The default theme includes a webmention counts partial that shows like/repost/reply counts on post cards.

### Including the Partial

Add to your card templates (e.g., `article-card.html`):

```html
<article class="article-card">
  <a href="{{ post.Href }}">
    <h3>{{ post.Title }}</h3>
    <p>{{ post.Description }}</p>
  </a>
  <footer>
    <time>{{ post.Date.Format "Jan 2, 2006" }}</time>
    {% include "partials/webmention-counts.html" %}
  </footer>
</article>
```

### The Counts Partial

The `partials/webmention-counts.html` partial renders compact counts:

```html
{% if post.Extra.webmentions %}
{% set likes = post.Extra.webmentions|selectattr("WMProperty", "equalto", "like-of")|list|length %}
{% set reposts = post.Extra.webmentions|selectattr("WMProperty", "equalto", "repost-of")|list|length %}
{% set replies = post.Extra.webmentions|selectattr("WMProperty", "equalto", "in-reply-to")|list|length %}

{% if likes > 0 or reposts > 0 or replies > 0 %}
<div class="webmention-counts">
  {% if likes > 0 %}<span class="wm-count wm-likes" title="Likes">{{ likes }}</span>{% endif %}
  {% if reposts > 0 %}<span class="wm-count wm-reposts" title="Reposts">{{ reposts }}</span>{% endif %}
  {% if replies > 0 %}<span class="wm-count wm-replies" title="Replies">{{ replies }}</span>{% endif %}
</div>
{% endif %}
{% endif %}
```

### Styling Counts

The default theme uses semantic color variables:

```css
.webmention-counts {
  display: flex;
  gap: 0.75rem;
  font-size: 0.875rem;
}

.wm-count {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
}

.wm-count::before {
  font-size: 1rem;
}

.wm-likes::before { content: ""; color: var(--color-error); }
.wm-reposts::before { content: ""; color: var(--color-success); }
.wm-replies::before { content: ""; color: var(--color-primary); }
```

## Engagement Leaderboard

The `webmentions_leaderboard` plugin calculates top posts by engagement and makes the data available for analytics pages.

### Enabling the Leaderboard

The leaderboard plugin runs automatically when webmentions are present. No additional configuration needed.

### Accessing Leaderboard Data

In any Jinja-enabled page, access the leaderboard via `config.Extra.webmention_leaderboard`:

```markdown
---
title: "Site Analytics"
jinja: true
---

# Content Analytics

{% if config.Extra.webmention_leaderboard %}

## Most Liked Posts

| Likes | Post |
|------:|------|
{% for entry in config.Extra.webmention_leaderboard.TopLiked %}| {{ entry.Likes }} | [{{ entry.Title }}]({{ entry.Href }}) |
{% endfor %}

## Most Discussed Posts

| Replies | Post |
|--------:|------|
{% for entry in config.Extra.webmention_leaderboard.TopReplied %}| {{ entry.Replies }} | [{{ entry.Title }}]({{ entry.Href }}) |
{% endfor %}

## Most Shared Posts

| Reposts | Post |
|--------:|------|
{% for entry in config.Extra.webmention_leaderboard.TopReposted %}| {{ entry.Reposts }} | [{{ entry.Title }}]({{ entry.Href }}) |
{% endfor %}

## Top Engaged Posts (Total)

| Total | Likes | Reposts | Replies | Post |
|------:|------:|--------:|--------:|------|
{% for entry in config.Extra.webmention_leaderboard.TopTotal %}| {{ entry.Total }} | {{ entry.Likes }} | {{ entry.Reposts }} | {{ entry.Replies }} | [{{ entry.Title }}]({{ entry.Href }}) |
{% endfor %}

---

## Site Totals

- **Total Likes:** {{ config.Extra.webmention_leaderboard.TotalLikes }}
- **Total Reposts:** {{ config.Extra.webmention_leaderboard.TotalReposts }}
- **Total Replies:** {{ config.Extra.webmention_leaderboard.TotalReplies }}
- **Total Mentions:** {{ config.Extra.webmention_leaderboard.TotalMentions }}

{% else %}
*No webmention data available.*
{% endif %}
```

### Leaderboard Data Structure

| Field | Type | Description |
|-------|------|-------------|
| `TopLiked` | `[]LeaderboardEntry` | Top 20 posts by likes |
| `TopReposted` | `[]LeaderboardEntry` | Top 20 posts by reposts |
| `TopReplied` | `[]LeaderboardEntry` | Top 20 posts by replies |
| `TopTotal` | `[]LeaderboardEntry` | Top 20 posts by total engagement |
| `TotalLikes` | `int` | Site-wide total likes |
| `TotalReposts` | `int` | Site-wide total reposts |
| `TotalReplies` | `int` | Site-wide total replies |
| `TotalMentions` | `int` | Site-wide total mentions |

Each `LeaderboardEntry` contains:

| Field | Type | Description |
|-------|------|-------------|
| `Href` | `string` | Post URL |
| `Title` | `string` | Post title |
| `Likes` | `int` | Number of likes |
| `Reposts` | `int` | Number of reposts |
| `Replies` | `int` | Number of replies |
| `Bookmarks` | `int` | Number of bookmarks |
| `Mentions` | `int` | Number of generic mentions |
| `Total` | `int` | Total engagement |

## Full Configuration Reference

```toml
[markata-go.webmentions]
# Enable the webmentions plugin
enabled = true

# Send outgoing webmentions when linking to external sites
outgoing = true

# HTTP request timeout for sending/receiving
timeout = "30s"

# Directory to cache sent webmentions (avoids re-sending)
cache_dir = ".cache/webmentions"

# Maximum concurrent HTTP requests
concurrent_requests = 5

# User-Agent string for HTTP requests
user_agent = "markata-go/1.0 (WebMention; +https://github.com/WaylonWalker/markata-go)"

# API token for webmention.io (for receiving mentions)
webmention_io_token = ""

[markata-go.webmentions.bridges]
# Enable social media bridging detection
enabled = false

# Use Bridgy Fediverse for bridging
bridgy_fediverse = true

# Platform toggles
bluesky = true
twitter = true
mastodon = true
github = true
flickr = false

[markata-go.webmentions.bridges.filters]
# Limit to specific platforms (empty = all enabled)
platforms = []

# Limit to specific interaction types (empty = all)
# Options: "like", "repost", "reply", "bookmark", "mention"
interaction_types = []

# Minimum content length for replies
min_content_length = 0

# Block mentions from these domains
blocked_domains = []
```

## Troubleshooting

### WebMentions not being sent

1. Ensure `enabled = true` and `outgoing = true`
2. Check that target sites have WebMention endpoints
3. Review the cache directory for sent status
4. Check HTTP timeout settings

### Bridged mentions not appearing

1. Verify your webmention.io token is correct
2. Check that Bridgy is connected to your social accounts
3. Ensure the platform is enabled in `bridges` config
4. Check filters aren't blocking the mentions

### Platform not detected correctly

The platform detection uses URL patterns. If a mention isn't being detected:

1. Check the source URL in the raw WebMention data
2. Ensure the platform is enabled in config
3. File an issue if detection should be improved

## Resources

- [WebMention Spec](https://www.w3.org/TR/webmention/)
- [Bridgy Fed](https://fed.brid.gy/) - Social media bridging
- [webmention.io](https://webmention.io/) - Hosted WebMention receiver
- [IndieWeb Wiki](https://indieweb.org/Webmention)
