---
title: "Structured Data & SEO"
description: "Configure JSON-LD Schema.org markup, OpenGraph, and Twitter Cards for better SEO and social sharing"
date: 2024-01-15
published: true
tags:
  - documentation
  - seo
  - structured-data
---


markata-go automatically generates comprehensive structured data for your posts to improve search engine visibility and social media sharing.

## Features

- **JSON-LD Schema.org** - Structured data for Google Rich Results
- **OpenGraph** - Meta tags for Facebook, LinkedIn, and general social sharing  
- **Twitter Cards** - Meta tags for Twitter/X sharing previews

## Quick Start

Structured data is enabled by default. For basic usage, just ensure your posts have good frontmatter:

```yaml
---
title: "My Post Title"
description: "A compelling description for SEO"
date: 2024-01-15
tags: ["tag1", "tag2"]
---
```

## Configuration

### Site-Level SEO Configuration

Configure SEO settings in your `markata-go.toml`:

```toml
[markata-go.seo]
# Twitter/X username (without @)
twitter_handle = "yourusername"

# Default image for posts without images
default_image = "/images/og-default.jpg"

# Site logo for Schema.org
logo_url = "/images/logo.png"

[markata-go.seo.structured_data]
# Enable/disable structured data (default: true)
enabled = true

# Publisher information
[markata-go.seo.structured_data.publisher]
type = "Organization"  # or "Person"
name = "Your Site Name"
url = "https://example.com"
logo = "/images/logo.png"

# Default author for posts without explicit author
[markata-go.seo.structured_data.default_author]
type = "Person"
name = "Author Name"
url = "https://example.com/about"
```

### Frontmatter Fields

Override structured data per-post via frontmatter:

```yaml
---
title: "Post Title"
description: "Post description for meta tags"
date: 2024-01-15

# Optional fields
author: "Author Name"          # Override default author
image: "/images/post.jpg"      # OG/Twitter image
social_image: "/images/og.jpg" # Specific OG image override
twitter: "authorhandle"        # Author's Twitter (without @)
modified: "2024-01-16"         # Last modified date
tags: ["tag1", "tag2"]
---
```

## Generated Output

### JSON-LD Schema

For each post, markata-go generates a `BlogPosting` schema:

```json
{
  "@context": "https://schema.org",
  "@type": "BlogPosting",
  "headline": "Post Title",
  "description": "Post description",
  "datePublished": "2024-01-15T00:00:00Z",
  "dateModified": "2024-01-16T00:00:00Z",
  "author": {
    "@type": "Person",
    "name": "Author Name",
    "url": "https://example.com/about"
  },
  "publisher": {
    "@type": "Organization",
    "name": "Site Name",
    "logo": {
      "@type": "ImageObject",
      "url": "https://example.com/logo.png"
    }
  },
  "mainEntityOfPage": {
    "@type": "WebPage",
    "@id": "https://example.com/post-slug/"
  },
  "image": "https://example.com/images/post.jpg",
  "keywords": ["tag1", "tag2"]
}
```

### OpenGraph Tags

Generated tags include:

| Tag | Description |
|-----|-------------|
| `og:title` | Post title |
| `og:description` | Post description |
| `og:type` | "article" for posts, "website" for pages |
| `og:url` | Canonical URL |
| `og:site_name` | Site title |
| `og:image` | Post or default image |
| `og:locale` | Content locale |
| `article:published_time` | Publication date |
| `article:modified_time` | Last modified date |
| `article:author` | Author URL |
| `article:tag` | Post tags |

### Twitter Card Tags

Generated tags include:

| Tag | Description |
|-----|-------------|
| `twitter:card` | "summary_large_image" or "summary" |
| `twitter:site` | Site Twitter handle |
| `twitter:creator` | Author Twitter handle |
| `twitter:title` | Post title |
| `twitter:description` | Post description (truncated to 200 chars) |
| `twitter:image` | Post or default image |

## Image Priority

Images are selected in this order:

1. `social_image` from frontmatter (OG-specific override)
2. `image` from frontmatter
3. `default_image` from SEO config

## Author Priority

Authors are determined in this order:

1. `author` field from frontmatter
2. `default_author` from structured data config
3. `author` from site config

## Disabling Structured Data

To disable structured data generation:

```toml
[markata-go.seo.structured_data]
enabled = false
```

## Testing Your Structured Data

Use these tools to validate your structured data:

- [Google Rich Results Test](https://search.google.com/test/rich-results)
- [Facebook Sharing Debugger](https://developers.facebook.com/tools/debug/)
- [Twitter Card Validator](https://cards-dev.twitter.com/validator)
- [Schema.org Validator](https://validator.schema.org/)

## Best Practices

1. **Always include descriptions** - Good descriptions improve both SEO and social sharing
2. **Use quality images** - Social images should be at least 1200x630 pixels
3. **Keep titles under 60 characters** - Prevents truncation in search results
4. **Set publication dates** - Helps search engines understand content freshness
5. **Use relevant tags** - Tags become keywords in schema markup

## Template Access

Structured data is available in templates via `post.structured_data`:

```html
{# JSON-LD #}
{% if post.structured_data.jsonld %}
<script type="application/ld+json">
{{ post.structured_data.jsonld | safe }}
</script>
{% endif %}

{# OpenGraph #}
{% for meta in post.structured_data.opengraph %}
<meta property="{{ meta.property }}" content="{{ meta.content }}">
{% endfor %}

{# Twitter Cards #}
{% for meta in post.structured_data.twitter %}
<meta name="{{ meta.name }}" content="{{ meta.content }}">
{% endfor %}
```

## See Also

- [Configuration Reference](/docs/guides/configuration/)
- [Post Formats](/docs/guides/post-formats/)
- [Themes](/docs/guides/themes/)
