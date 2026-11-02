---
title: "Frontmatter Guide"
description: "Complete guide to post metadata, built-in fields, and custom frontmatter in markata-go"
date: 2024-01-15
published: true
slug: /docs/guides/frontmatter/
tags:
  - documentation
  - frontmatter
  - content
---

## Frontmatter

Frontmatter is the **metadata block** at the top of your Markdown files that tells markata-go how to process, display, and organize your content. This guide covers everything you need to know about using frontmatter effectively.

In the workbench editor, press **Ctrl+Space** to request metadata completion.
Known field values are suggested while editing scalar values and YAML list
items, including existing site tags.

## Table of Contents

- [What is Frontmatter?](#what-is-frontmatter)
- [Basic Frontmatter Fields](#basic-frontmatter-fields)
- [Complete Field Reference](#complete-field-reference)
- [Custom Fields (Extra)](#custom-fields-extra)
- [Examples](#examples)
- [Common Patterns](#common-patterns)
- [Frontmatter in Filtering](#frontmatter-in-filtering)

---

## What is Frontmatter?

Frontmatter is **YAML metadata** placed at the very top of a Markdown file, enclosed between two lines of three dashes (`---`). It defines properties like the post's title, date, tags, and any custom data you want to associate with the content.

### Visual Example

```markdown
---
title: "Getting Started with Go"
date: 2024-01-15
published: true
tags:
  - go
  - tutorial
  - beginners
description: "A beginner-friendly introduction to Go programming."
---

# Welcome to Go!

This is where your actual content begins...
```

### How markata-go Parses Frontmatter

When markata-go processes your Markdown files, it:

1. **Detects the opening delimiter** - Looks for `---` at the very start of the file
2. **Extracts the YAML block** - Reads everything until the closing `---`
3. **Parses the YAML** - Converts the metadata into structured data
4. **Separates the content** - Everything after the second `---` becomes the post body

**Important notes:**

- The opening `---` must be on the very first line of the file
- The closing `---` must be on its own line
- Content before the opening `---` or malformed delimiters will cause errors
- Files without frontmatter are valid - they'll use default values

### Valid Frontmatter

```markdown
---
title: My Post
---
Content here
```

### Invalid Frontmatter (Missing Closing Delimiter)

```markdown
---
title: My Post

Content here (ERROR: unclosed frontmatter)
```

### No Frontmatter (Also Valid)

```markdown
# My Post

Content starts immediately - defaults will be used.
```

---

## Basic Frontmatter Fields

These are the core fields that markata-go recognizes and uses directly.

### title (string)

The display title of your post.

```yaml
title: "Understanding Goroutines"
```

- **Type:** string
- **Default:** None (derived from filename if not set)
- **Used for:** Page title, `<title>` tag, feed listings, slug generation

If not provided, the slug is derived from the filename instead.

Titles also accept HTML character references. markata-go decodes references
before rendering, so `It&rsquo;s just the Carpet&hellip;` displays as
`It’s just the Carpet…` instead of showing the entity names. Use
`title_text` in custom templates when you need the plain-text title.

### slug (string)

The URL-safe identifier for your post. Determines the URL path.

```yaml
slug: "understanding-goroutines"
```

- **Type:** string
- **Default:** Auto-generated from `title` (or filename if no title)
- **Used for:** URL path (`/understanding-goroutines/`)

**Auto-generation rules:**
- Converts to lowercase
- Replaces spaces with hyphens
- Removes special characters
- Collapses multiple hyphens

```yaml
# These titles produce these slugs:
title: "Hello World!"        # slug: hello-world
title: "Go's Best Features"  # slug: gos-best-features
title: "Part 1: The Basics"  # slug: part-1-the-basics
```

#### Custom Slugs

You can explicitly set any slug in frontmatter:

```yaml
# Simple custom slug
slug: my-custom-url

# Nested path (creates /docs/api/overview/)
slug: docs/api/overview

# Leading/trailing slashes are automatically stripped
slug: /about/            # becomes: about
slug: /docs/guides/      # becomes: docs/guides
```

#### Homepage from Custom Slug

Any post can become the homepage by setting an empty or slash slug:

```yaml
---
title: "Welcome to My Site"
slug: ""                    # This post becomes the homepage at /
published: true
---

Welcome to my site!
```

Or equivalently:

```yaml
---
title: "Welcome to My Site"
slug: /                     # Also becomes the homepage at /
published: true
---
```

This is useful when you want a custom homepage from content in any directory, rather than using a feed or `index.md` at the root.

**Note:** If multiple posts have the same slug (including empty slug for homepage), the build will fail with a path conflict error to prevent accidental overwrites.

#### Special Case: `index.md` Files

Files named `index.md` automatically generate slugs based on their directory:

| File Path | Auto-Generated Slug | URL |
|-----------|---------------------|-----|
| `./index.md` | `""` (empty) | `/` (homepage) |
| `docs/index.md` | `docs` | `/docs/` |
| `blog/guides/index.md` | `blog/guides` | `/blog/guides/` |

This allows directory-based URL structures without needing to set explicit slugs.

### date (date)

The publication date of the post.

```yaml
date: 2024-01-15
```

- **Type:** date (YYYY-MM-DD format recommended)
- **Default:** None
- **Used for:** Sorting, display, feeds, scheduled publishing

**Supported date formats:**

markata-go supports a wide variety of date formats for maximum flexibility:

```yaml
# ISO 8601 formats (recommended)
date: 2024-01-15                    # Date only
date: 2024-01-15T10:30:00           # With time
date: 2024-01-15T10:30:00Z          # With UTC timezone
date: 2024-01-15T10:30:00+05:00     # With timezone offset
date: 2024-01-15 10:30:00           # With space separator
date: 2024-01-15 10:30              # Without seconds

# Single-digit hours (automatically normalized)
date: 2024-01-15 1:30:00            # 1am
date: 2024-01-15 9:30:00            # 9am

# Slash-separated dates
date: 2024/01/15                    # YYYY/MM/DD
date: 2024/01/15 10:30:00           # With time
date: 01/15/2024                    # MM/DD/YYYY (US format)

# Named month formats
date: January 15, 2024              # Full month name
date: Jan 15, 2024                  # Abbreviated month
date: 15 January 2024               # Day first
date: 15 Jan 2024                   # Day first, abbreviated

# European format
date: 15-01-2024                    # DD-MM-YYYY
```

**Note:** Malformed time components are automatically corrected. For example, `2024-01-15 8:011:00` (typo) will be parsed as `2024-01-15 08:11:00`.

### published (boolean)

Whether the post should be included in public feeds and listings.

```yaml
published: true
```

- **Type:** boolean
- **Default:** `false`
- **Used for:** Filtering posts in feeds, sitemaps, and RSS

**Shadow Pages:** Posts with `published: false` are still rendered to HTML and accessible via direct URL. They become "shadow pages" - accessible but not discoverable through normal site navigation.

| `published` | Behavior |
|-------------|----------|
| `true` | Rendered + included in feeds, sitemaps, RSS |
| `false` | Rendered as "shadow page" + NOT in feeds, sitemaps, RSS |

**Use cases for shadow pages:**
- Draft content accessible to reviewers via direct URL
- Private documentation not linked publicly
- Work-in-progress content shared with specific people
- Admin pages or staging content

**Accepted values:**

```yaml
published: true    # or: yes, on
published: false   # or: no, off
```

### draft (boolean)

Marks the post as a work-in-progress that should NOT be rendered at all.

```yaml
draft: true
```

- **Type:** boolean
- **Default:** `false`
- **Used for:** Truly private content that shouldn't be accessible

**Important:** Posts with `draft: true` are **never rendered** to HTML. Use this for content you don't want accessible at all, even via direct URL.

| `draft` | Behavior |
|---------|----------|
| `false` | Content is rendered (visibility depends on `published`) |
| `true` | Content is NOT rendered at all |

**Comparison with `published`:**

| Scenario | `published` | `draft` | Rendered? | In Feeds? |
|----------|-------------|---------|-----------|-----------|
| Public post | `true` | `false` | Yes | Yes |
| Shadow page | `false` | `false` | Yes | No |
| Private WIP | any | `true` | No | No |

### tags (list)

A list of tags for categorizing the post.

```yaml
tags:
  - go
  - concurrency
  - tutorial
```

Or inline format:

```yaml
tags: [go, concurrency, tutorial]
```

For compatibility with older content, `tags: go` is also accepted as a single tag. Diagnostics still warns about scalar shorthand and recommends a list, which is the canonical form.

- **Type:** list of strings
- **Default:** `[]` (empty list)
- **Used for:** Filtering, tag pages, SEO, organization

### description (string)

A brief summary of the post content.

```yaml
description: "Learn how to use goroutines for concurrent programming in Go."
```

- **Type:** string
- **Default:** Auto-generated from first ~160 characters of content
- **Used for:** Meta description, SEO, social previews

**Note:** The `description` field is no longer used for article card excerpts in feeds. Excerpts are now auto-generated from your content (first 3 paragraphs or 1500 characters). However, `description` may still be used for SEO meta tags and RSS feeds.

**Auto-generation:** If not provided, markata-go extracts the first paragraph or ~160 characters from your content (with HTML tags stripped).

### template (string)

The HTML template file to use for rendering this post.

```yaml
template: "tutorial.html"
```

- **Type:** string
- **Default:** `"post.html"`
- **Used for:** Custom layouts per post

Templates are looked up in:
1. `templates/` directory in your project
2. Theme templates
3. Default theme fallback

Example built-in template choices include `post.html`, `landing.html`, `docs.html`, and `slides.html`.

`slides.html` is intended for presentation decks and pairs well with the shared assets system when you want to self-host Reveal.js.

### content_width (string)

Set `content_width: wide` to opt a default-theme `post.html` page into a wider
story frame, a larger balanced title, and a lead opening paragraph. It pairs
with `md-section` for sticky section heads and `md-grid rail` for overview
cards. Unset or unknown values keep the normal layout. Custom templates must
provide their own support. See [Markdown story pages](/docs/guides/markdown-containers/#markdown-story-pages).

---

## Complete Field Reference

| Field | Type | Default | Required | Description |
|-------|------|---------|----------|-------------|
| `title` | string | None | No | Display title of the post |
| `slug` | string | Auto-generated | No | URL path identifier |
| `date` | date | None | No | Publication date (YYYY-MM-DD) |
| `modified` | date | None | No | Last modified date (YYYY-MM-DD) |
| `published` | bool | `false` | No | Whether to include in public feeds |
| `draft` | bool | `false` | No | Whether this is a work-in-progress |
| `tags` | []string | `[]` | No | List of categorization tags |
| `description` | string | Auto-generated | No | Brief summary for SEO/meta tags |
| `template` | string | `"post.html"` | No | Template file to use for rendering |
| `content_width` | string | Normal theme width | No | `wide` opts a default-theme post into story layout |
| `skip` | bool | `false` | No | Skip this file during processing entirely |
| `authors` | []string | `[]` | No | List of author IDs (multi-author support) |
| `author` | string | None | No | Legacy single-author name or ID |
| `by` | string | None | No | Alias for `author` |
| `writer` | string | None | No | Alias for `author` |

### Date Field Aliases

For compatibility with other static site generators, markata-go supports multiple frontmatter keys:

**Publication date aliases:** `publishdate`, `date`, `pubdate`

**Modified date aliases:** `lastmod`, `modified`, `updated`, `updated_at`, `last_modified`

When multiple aliases are provided, the **first listed wins**. This precedence matches common SSG conventions.

```yaml
---
# All of these set the publication date:
publishdate: 2024-01-15
date: 2024-01-15
pubdate: 2024-01-15

# All of these set the modified date:
lastmod: 2024-02-20
modified: 2024-02-20
updated: 2024-02-20
updated_at: 2024-02-20
last_modified: 2024-02-20
---
```

### Field Details

#### skip

Completely exclude a file from processing:

```yaml
skip: true
```

Use this for files you want to keep in your content directory but never process (notes, drafts not ready for review, etc.).

#### authors

Assign one or more authors to a post by referencing author IDs defined in your site configuration:

```yaml
---
title: "Collaborative Article"
authors:
  - waylon
  - guest
---
```

Each author ID is resolved against the `[markata-go.authors.authors]` config map during the build. Resolved author objects (with name, avatar, role, bio, etc.) are available in templates as `post.author_objects`.

You can use the extended format to specify per-post roles and details for each author:

```yaml
---
title: "Collaborative Article"
authors:
  - id: waylon
    role: author
    details: wrote the introduction and conclusion
  - id: codex
    role: pair programmer
    details: wrote the code examples
  - id: kimmi
    role: outliner
---
```

The `role` overrides the author's config-level role for this post only. The `details` field provides a short description of what the author did on this specific post, displayed as a tooltip when hovering over the author's name in the byline. Both fields are optional and can be used independently.

**Key aliases:** For convenience, the extended format supports aliases so you can use whichever key name feels natural:

| Canonical Key | Aliases |
|---------------|---------|
| `id` | `name`, `handle` |
| `role` | `job`, `position`, `part`, `title` |
| `details` | `detail`, `description` |

For example, these are all equivalent:

```yaml
---
authors:
  - name: waylon
    title: maintainer
    description: built the feature
  # same as:
  - id: waylon
    role: maintainer
    details: built the feature
---
```

The canonical key always takes precedence when both it and an alias are present.

You can mix simple string IDs and extended format in the same array:

```yaml
---
authors:
  - waylon
  - id: codex
    role: editor
    details: reviewed and edited the draft
---
```

If no `authors` or `author` field is set, the default author from config (the one with `default = true`) is assigned automatically.

#### author

Legacy single-author field. Use this for simple sites or backward compatibility:

```yaml
---
title: "Solo Post"
author: "Jane Doe"
---
```

**Aliases:** The `by` and `writer` frontmatter fields are aliases for `author`. This allows more natural frontmatter syntax:

```yaml
---
title: "Solo Post"
by: "Jane Doe"
---
```

```yaml
---
title: "Solo Post"
writer: "Jane Doe"
---
```

All three (`author`, `by`, `writer`) set the same author field. Priority order: `authors` > `author` > `by` > `writer`. Only the first match is used.

**Priority:** If both `authors` and `author` (or its aliases) are set on the same post, `authors` takes precedence.

**Related:** See [Authors Configuration](configuration.md#authors-configuration) for defining author profiles.

---

## Custom Fields (Extra)

Any frontmatter field that isn't a built-in field is automatically stored in the `Extra` map. This allows you to add **any custom metadata** to your posts.

### Adding Custom Fields

```yaml
---
title: "Building a REST API"
date: 2024-01-15
published: true

# Custom fields - stored in Extra
author: "Jane Doe"
category: "Backend"
series: "API Development"
series_order: 1
featured: true
cover_image: "/images/api-cover.jpg"
difficulty: "intermediate"
reading_time: "8 min"
---
```

### Accessing Custom Fields in Templates

Custom fields are available via `post.Extra` in templates:

```html
{% if post.Extra.featured %}
<span class="badge">Featured</span>
{% endif %}

{% if post.Extra.author %}
<p class="author">By {{ post.Extra.author }}</p>
{% endif %}

{% if post.Extra.cover_image %}
<img src="{{ post.Extra.cover_image }}" alt="{{ post.Title }}">
{% endif %}

{% if post.Extra.series %}
<div class="series-info">
    Part {{ post.Extra.series_order }} of {{ post.Extra.series }}
</div>
{% endif %}
```

### Common Custom Field Use Cases

#### Author Information

```yaml
author: "Jane Doe"
author_email: "jane@example.com"
author_url: "https://janedoe.dev"
author_avatar: "/images/avatars/jane.jpg"
```

#### Series/Collections

```yaml
series: "Building a Blog with Go"
series_order: 3
series_total: 5
```

#### Visual Elements

```yaml
cover_image: "/images/posts/my-cover.jpg"
thumbnail: "/images/posts/my-thumb.jpg"
og_image: "/images/social/my-post-og.jpg"
hero_video: "https://youtube.com/watch?v=..."
```

#### Content Metadata

```yaml
difficulty: "beginner"          # beginner, intermediate, advanced
reading_time: "5 min read"
word_count: 1250
revision: 3
```

**Note:** The `modified` date is now a first-class field (see above). Use `lastmod`, `modified`, or `updated` in frontmatter - they all set the same field. If multiple are present, the first listed alias wins.

#### Categorization

```yaml
category: "Tutorials"
subcategory: "Web Development"
topic: "Go Programming"
```

#### Flags and Features

```yaml
featured: true
pinned: true
sponsored: false
comments_enabled: true
toc: true                       # Enable table of contents
math: true                      # Enable math rendering
mermaid: true                   # Enable diagrams
```

#### External Links

```yaml
canonical_url: "https://example.com/original-post"
github_repo: "https://github.com/user/project"
demo_url: "https://demo.example.com"
```

#### Aliases (for Wikilinks)

Define alternative names that can be used in wikilinks to link to this post:

```yaml
---
title: "ECMAScript Language Specification"
slug: "ecmascript"
aliases:
  - js
  - javascript
  - JavaScript
---
```

With these aliases, any of the following wikilinks will resolve to this post:

- `[[ecmascript]]` - matches the slug
- `[[js]]` - matches an alias
- `[[javascript]]` - matches an alias
- `[[JavaScript]]` - case-insensitive alias match

!!! note "Slug Precedence"

    If another post has `slug: "javascript"`, then `[[javascript]]` will link to
    that post (slug match) rather than the post with the alias. Slugs always take
    precedence over aliases.

**Alias field synonyms:** The frontmatter loader also accepts `alias`, `handle`, and `handles` and normalizes them into `aliases`.

Alias resolution is the same on cold, warm, and edited builds. Loading a post
from cache does not change its alternative names or require clearing caches
before a wikilink can resolve.

---

## Examples

### Minimal Frontmatter

The absolute minimum for a publishable post:

```yaml
---
title: "Hello World"
published: true
---

Your content here.
```

### Full Frontmatter Example

A comprehensive example using all built-in fields:

```yaml
---
title: "Complete Guide to Error Handling in Go"
slug: "go-error-handling-guide"
date: 2024-01-15
published: true
draft: false
tags:
  - go
  - error-handling
  - best-practices
  - tutorial
description: "Master error handling in Go with this comprehensive guide covering best practices, custom errors, and common patterns."
template: "tutorial.html"
---

Content begins here...
```

### Blog Post Example

A typical blog post with custom fields:

```yaml
---
title: "Why I Switched from Python to Go"
date: 2024-01-15
published: true
tags:
  - go
  - python
  - opinion
  - programming
description: "My journey from Python to Go and the lessons learned along the way."

# Custom fields
author: "Alex Chen"
category: "Opinion"
featured: true
cover_image: "/images/python-to-go.jpg"
reading_time: "6 min read"
---

# Why I Switched from Python to Go

After 5 years of Python development, I made the switch to Go...
```

### Documentation Page Example

For technical documentation:

```yaml
---
title: "Configuration Reference"
slug: "docs/configuration"
published: true
tags:
  - documentation
  - reference
description: "Complete reference for markata-go configuration options."
template: "docs.html"

# Custom fields
section: "Reference"
order: 10
toc: true
prev_page: "/docs/getting-started/"
next_page: "/docs/plugins/"
---

# Configuration Reference

This page documents all available configuration options...
```

### Landing Page Example

For standalone pages with custom layouts:

```yaml
---
title: "Welcome to My Site"
slug: ""
published: true
template: "landing.html"

# Custom fields
hero_title: "Build Faster with Go"
hero_subtitle: "A static site generator that respects your time"
cta_text: "Get Started"
cta_url: "/docs/getting-started/"
features:
  - title: "Fast Builds"
    description: "Compile thousands of pages in seconds"
    icon: "lightning"
  - title: "Plugin System"
    description: "Extend functionality with Go plugins"
    icon: "puzzle"
  - title: "Feed Generation"
    description: "RSS, Atom, and JSON feeds built-in"
    icon: "rss"
---

Additional content for the landing page...
```

### Tutorial with Series

For posts that are part of a series:

```yaml
---
title: "Building a CLI Tool - Part 2: Adding Commands"
date: 2024-01-20
published: true
tags:
  - go
  - cli
  - tutorial
description: "Learn how to add subcommands to your Go CLI application using cobra."
template: "tutorial.html"

# Series information
series: "Building a CLI Tool in Go"
series_slug: "go-cli-tutorial"
series_order: 2
series_total: 5

# Tutorial metadata
difficulty: "intermediate"
prerequisites:
  - "Basic Go knowledge"
  - "Part 1 of this series"
code_repo: "https://github.com/example/go-cli-tutorial"
---

# Adding Commands

In the previous part, we set up our project structure...
```

---

## Media Fields

The `image` and `video` frontmatter fields can be used **interchangeably** in photo and video card templates. The system auto-detects whether a URL points to a video or image based on the file extension. The image library also indexes both fields; `image` and `video` are cover fallbacks when no explicit `cover` or `cover_image` is set.

### Recognized Video Extensions

`.avi`, `.mkv`, `.mp4`, `.webm`, `.mov`, `.m4v`, `.ogv`, `.ogg` (case-insensitive)

Any other extension (or no extension) is treated as an image.

### How It Works

Card templates use two filters to resolve media:

1. **`media_url`** -- picks the first non-empty value between `image` and `video` fields
2. **`is_video`** -- checks the file extension to decide whether to render `<video>` or `<img>`

This means all of the following work equally well for a photo card:

```yaml
# Option A: image field with a video file
---
image: "/media/demo.mp4"
---

# Option B: video field with a video file
---
video: "/media/demo.mp4"
---

# Option C: image field with an image file
---
image: "/photos/sunset.jpg"
---

# Option D: both fields (image takes priority)
---
image: "/photos/sunset.jpg"
video: "/media/demo.mp4"
---
```

### Using Media Fields in Custom Templates

You can use the `is_video` and `media_url` filters in your own templates:

```html
{% with post.image|media_url:post.video as media_src %}
{% if media_src %}
  {% if media_src|is_video %}
  <video src="{{ media_src }}" autoplay muted loop playsinline></video>
  {% else %}
  <img src="{{ media_src }}" alt="{{ post.title }}">
  {% endif %}
{% endif %}
{% endwith %}
```

---

## Common Patterns

### Draft Workflow

Use `draft` and `published` together for a clear workflow:

```yaml
# Work in progress - not visible anywhere
---
title: "My Draft Post"
draft: true
published: false
---

# Ready for review - still not public
---
title: "My Draft Post"
draft: true
published: false
---

# Published - visible to everyone
---
title: "My Draft Post"
draft: false
published: true
---
```

**Feed filter for published only:**

```toml
[[markata-go.feeds]]
slug = "blog"
filter = "published == True"
```

**Feed filter excluding drafts:**

```toml
[[markata-go.feeds]]
slug = "blog"
filter = "published == True and draft == False"
```

### Scheduled Publishing

Combine `date` with feed filters for scheduled publishing:

```yaml
---
title: "New Year Announcement"
date: 2025-01-01
published: true
---
```

**Feed filter for past/current dates only:**

```toml
[[markata-go.feeds]]
slug = "blog"
filter = "published == True and date <= today"
```

Posts with future dates won't appear until that date arrives. Rebuild your site daily (via CI/CD) to "publish" scheduled posts.

### Custom Templates Per Post

Use different layouts for different types of content:

```yaml
# Regular blog post
---
title: "My Post"
template: "post.html"
---

# Tutorial with sidebar
---
title: "Go Tutorial"
template: "tutorial.html"
---

# Full-width landing page
---
title: "About Me"
template: "landing.html"
---

# Documentation with navigation
---
title: "API Reference"
template: "docs.html"
---
```

### Series of Posts

Organize related posts into a series:

**Post 1:**
```yaml
---
title: "Web Scraping with Go - Part 1: Basics"
date: 2024-01-10
series: "Web Scraping with Go"
series_order: 1
tags: [go, web-scraping, tutorial]
---
```

**Post 2:**
```yaml
---
title: "Web Scraping with Go - Part 2: Handling JavaScript"
date: 2024-01-17
series: "Web Scraping with Go"
series_order: 2
tags: [go, web-scraping, tutorial]
---
```

**Filter for the series:**
```toml
[[markata-go.feeds]]
slug = "series/web-scraping-go"
title = "Web Scraping with Go"
filter = "Extra.series == 'Web Scraping with Go'"
sort = "Extra.series_order"
reverse = false
```

### Canonical URLs for Cross-Posted Content

When you publish the same content elsewhere:

```yaml
---
title: "My Post"
canonical_url: "https://dev.to/username/my-post"
---
```

Use in templates:

```html
{% if post.Extra.canonical_url %}
<link rel="canonical" href="{{ post.Extra.canonical_url }}">
{% else %}
<link rel="canonical" href="{{ config.URL }}{{ post.Href }}">
{% endif %}
```

---

## Frontmatter in Filtering

Frontmatter fields power the feed filtering system. You can filter posts based on any frontmatter value.

### Filter Syntax

Filters use a Python-like expression syntax:

```toml
filter = "published == True"
filter = "'tutorial' in tags"
filter = "date >= '2024-01-01'"
```

### Filtering by Built-in Fields

```toml
# Published posts only
filter = "published == True"

# Exclude drafts
filter = "draft == False"

# Posts with a specific tag
filter = "'go' in tags"

# Posts from 2024
filter = "date >= '2024-01-01' and date < '2025-01-01'"

# Posts with a specific slug prefix
filter = "slug.startswith('tutorials/')"

# Posts using a specific template
filter = "template == 'tutorial.html'"
```

### Filtering by Custom Fields (Extra)

Access custom fields directly by name:

```toml
# Featured posts
filter = "featured == True"

# Posts by specific author
filter = "author == 'Jane Doe'"

# Posts in a specific category
filter = "category == 'Tutorials'"

# Posts in a series
filter = "series == 'Web Scraping with Go'"

# Intermediate difficulty tutorials
filter = "difficulty == 'intermediate'"
```

### Combined Filters

Use `and`, `or`, and parentheses for complex filters:

```toml
# Published tutorials
filter = "published == True and 'tutorial' in tags"

# Featured posts from 2024
filter = "published == True and featured == True and date >= '2024-01-01'"

# Go or Python tutorials
filter = "published == True and ('go' in tags or 'python' in tags)"

# Published, non-draft, with specific category
filter = "published == True and draft == False and category == 'Backend'"
```

### Filtering with Dates

Special date values are available:

```toml
# Posts up to today (no future posts)
filter = "date <= today"

# Posts from the last 30 days
filter = "date >= today - 30"

# Recent posts (within current timestamp)
filter = "date <= now"
```

### Filter Operators Reference

| Operator | Description | Example |
|----------|-------------|---------|
| `==` | Equal to | `published == True` |
| `!=` | Not equal to | `draft != True` |
| `>` | Greater than | `date > '2024-01-01'` |
| `>=` | Greater than or equal | `date >= '2024-01-01'` |
| `<` | Less than | `date < '2025-01-01'` |
| `<=` | Less than or equal | `date <= today` |
| `in` | Value in collection | `'go' in tags` |
| `and` | Logical AND | `published == True and featured == True` |
| `or` | Logical OR | `'go' in tags or 'python' in tags` |
| `not` | Logical NOT | `not draft` |

### String Methods in Filters

```toml
# Slugs starting with a prefix
filter = "slug.startswith('tutorials/')"

# Slugs ending with a suffix
filter = "slug.endswith('-guide')"

# Slugs containing a substring
filter = "slug.contains('api')"

# Case-insensitive comparison
filter = "title.lower().contains('go')"
```

### Complete Filter Examples

**Home page - recent published posts:**
```toml
[[markata-go.feeds]]
slug = ""
title = "Recent Posts"
filter = "published == True and date <= today"
sort = "date"
reverse = true
items_per_page = 5
```

**Tutorials section:**
```toml
[[markata-go.feeds]]
slug = "tutorials"
title = "Tutorials"
filter = "published == True and 'tutorial' in tags"
sort = "date"
reverse = true
```

**Featured posts:**
```toml
[[markata-go.feeds]]
slug = "featured"
title = "Featured"
filter = "published == True and featured == True"
sort = "date"
reverse = true
items_per_page = 6
```

**Author archive:**
```toml
[[markata-go.feeds]]
slug = "authors/jane-doe"
title = "Posts by Jane Doe"
filter = "published == True and author == 'Jane Doe'"
sort = "date"
reverse = true
```

**Beginner-friendly content:**
```toml
[[markata-go.feeds]]
slug = "beginners"
title = "Beginner Guides"
filter = "published == True and difficulty == 'beginner'"
sort = "date"
reverse = true
```

---

---

## Next Steps

Now that you understand frontmatter, here are recommended next steps based on what you want to do:

**Organize your content with feeds:**
- [Feeds Guide](/docs/guides/feeds/) - Learn how to create filtered, sorted collections using your frontmatter fields

**Customize how your content looks:**
- [Templates Guide](/docs/guides/templates/) - Use frontmatter values in your templates
- [Themes Guide](/docs/guides/themes/) - Style your site with color palettes and custom CSS

**Learn about Markdown features:**
- [Markdown Guide](/docs/guides/markdown/) - Advanced Markdown syntax and extensions

---

## See Also

- [Feeds Guide](/docs/guides/feeds/) - Complete guide to feed filtering and configuration
- [Templates Guide](/docs/guides/templates/) - Using frontmatter in templates
- [Configuration Guide](/docs/guides/configuration/) - Site-wide defaults and settings
- [Quick Reference](/docs/guides/quick-reference/) - Frontmatter snippets and cheat sheets
