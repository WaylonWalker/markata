---
title: "Card Types"
description: "How posts are displayed differently in feeds based on their template type"
date: 2026-01-29
published: true
tags:
  - templates
  - feeds
  - cards
---


Cards are how posts appear in feed listings. Different types of content benefit from different visual treatments - a photo post should emphasize the image, while a blog post should show the title and excerpt.

## How Cards Are Selected

Cards are selected based on the `template` frontmatter field:

```yaml
---
title: "My Photo"
template: photo
---
```

> **Note:** `templateKey` is supported as an alias for backwards compatibility with Python markata.

The `template` field determines:
1. Which **card template** renders the post in feeds
2. Which **page template** renders the full post (if defined)

## Field Name Comparison

markata-go uses `template` as the primary field, following Zola and Pelican conventions:

| SSG | Field Name |
|-----|------------|
| Hugo, Jekyll, Eleventy, Astro | `layout` |
| Zola, Pelican, **markata-go** | `template` |

Both `template` and `templateKey` work in markata-go frontmatter.

## Card Types

| Card Type | template Values | Best For |
|-----------|-----------------|----------|
| `article` | `blog-post`, `article`, `post`, `essay`, `tutorial` | Long-form content with title, excerpt, tags |
| `note` | `note`, `ping`, `thought`, `status`, `tweet` | Short thoughts, micro-posts |
| `photo` | `photo`, `shot`, `shots`, `image`, `gallery` | Image-focused posts |
| `video` | `video`, `clip`, `cast`, `stream` | Video content with thumbnail |
| `link` | `link`, `bookmark`, `stars` | External links, bookmarks |
| `quote` | `quote`, `quotation` | Quoted text with attribution |
| `guide` | `guide`, `series`, `step`, `chapter` | Multi-part tutorials |
| `inline` | `gratitude`, `inline`, `micro`, `til` | Full content shown in feed |
| `contact` | `contact`, `character`, `person` | Person/character profile cards |
| `default` | (any other value) | Fallback for unmapped types |

## Linked posts and source previews

A `link`, `bookmark`, `stars`, or `til` post with a `link` or `url`
frontmatter destination preserves its local permalink and shows a visible
**Source** link on the post page. Feed cards also link to the external domain.

```yaml
---
title: "Worth watching"
template: link
link: "https://www.youtube.com/watch?v=R_uf5OfMGio"
published: true
---

My thoughts on the video.
```

The destination is automatically embedded before your commentary through
the normal external embed system. YouTube uses the site's Lite YouTube
renderer; other URLs use cached oEmbed/Open Graph previews and fall back to
a clickable source card when metadata is unavailable. No new Markdown is
required. Explicit `![embed](url)` or `![[url]]` in the body is not
duplicated. Leave the `link`/`url` field unset if the post is not sharing
an external destination.

## Card Anatomy

Each card type has different visual elements:

### Article Card
- Large title (linked)
- Content excerpt (first 3 paragraphs or 1500 characters)
- Date, reading time
- Tags

### Note Card
- Compact layout with left accent border
- Title (optional)
- Short text content
- Date

### Photo Card
- Large image (16:9 aspect ratio)
- Caption below
- Date

### Inline Card
- Full `article_html` content rendered in the feed
- No "read more" - everything is visible
- Useful for gratitude journals, micro-posts

### Link Card
- Domain indicator
- Title and description
- External link styling

### Contact Card
- Avatar (from `avatar`, `image`, or `icon` frontmatter fields)
- Name (post title)
- Handle (`@handle`)
- Short bio (description)
- External site link (from `url` frontmatter field, if set)
- Tags

### Default Card Images

Posts without a mapped `template` use the default card. If the frontmatter
contains `image`, `cover`, `cover_image`, or `og_image`, the card displays that
media in a cinematic 12:5 frame. The source uses width-only sizing so the CDN
does not distort the image; CSS crops it with `object-fit: cover` so portrait
images do not make daily note cards too tall.

## Examples

### Blog Post (article card)
```yaml
---
title: "Understanding Go Interfaces"
template: blog-post
tags:
  - go
  - programming
---

Your article content here. The first three paragraphs (or up to 1500 characters) will automatically be shown as the excerpt in feed listings.
```

### Quick Thought (note card)
```yaml
---
title: "Ping 42"
template: ping
---

Just discovered that vim has a built-in terminal. Mind blown.
```

### Photo (photo card)
```yaml
---
title: "Sunset at the Beach"
template: shots
image: /images/sunset.jpg
---
```

### Gratitude Entry (inline card)
```yaml
---
template: gratitude
---

Today I'm grateful for good coffee and working CI pipelines.
```

## Configuring Excerpts

Article cards automatically show an excerpt from your content - the first 3 paragraphs or 1500 characters, whichever is shorter.

### Default Behavior

By default, article cards extract:
- **Up to 3 paragraphs** from your rendered HTML content
- **OR up to 1500 characters** (whichever comes first)
- With intelligent truncation at word boundaries
- An ellipsis (`...`) is added when content is truncated

### How It Works

1. The `excerpt` filter extracts `<p>` tags from your rendered `article_html`
2. It collects paragraphs until reaching 3 paragraphs OR 1500 characters
3. It strips inner HTML tags but preserves text
4. It adds ellipsis (`...`) if content was truncated
5. Output is safe HTML ready to render

### Tips

- **Short posts:** If your post is shorter than the limit, the entire post shows (no ellipsis)
- **No paragraphs:** If content has no `<p>` tags, the filter falls back to plain text truncation
- **Removal of description:** The old `description` frontmatter field is no longer used for article cards - excerpts are auto-generated from content

## Customizing Cards

To override a built-in card, create a file in your site's `templates/partials/cards/` directory:

```
templates/
  partials/
    cards/
      article-card.html   # Overrides built-in article card
      my-custom-card.html # New card type
```

To add a new card type, create a custom `card-router.html`:

```html
{% elif post.template == "my-type" %}
{% include "partials/cards/my-custom-card.html" %}
```

## Card CSS

Card styles are in `themes/default/static/css/cards.css`. Each card has a class:

- `.card-article`
- `.card-note`
- `.card-photo`
- `.card-video`
- `.card-link`
- `.card-quote`
- `.card-guide`
- `.card-inline`
- `.card-contact`
- `.card-default`

## How It Works

1. Feed template (`feed.html`) loops through posts
2. For each post, includes `partials/cards/card-router.html`
3. Card router checks `post.template` and includes the matching card template
4. Card template renders the post with appropriate HTML structure

```html
<!-- feed.html -->
{% for post in page.posts %}
{% include "partials/cards/card-router.html" %}
{% endfor %}
```

```html
<!-- card-router.html -->
{% if post.template == "blog-post" %}
{% include "partials/cards/article-card.html" %}
{% elif post.template == "ping" %}
{% include "partials/cards/note-card.html" %}
...
{% endif %}
```

## Template Access in Custom Templates

In your templates, you can access the template value via either name:

```html
{# Both work - template is preferred #}
{{ post.template }}
{{ post.templateKey }}
```
