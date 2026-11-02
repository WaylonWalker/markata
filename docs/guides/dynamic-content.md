---
title: "Dynamic Content"
description: "Using Jinja2 template syntax in Markdown files for dynamic content generation"
date: 2024-01-15
published: true
tags:
  - documentation
  - jinja
  - templates
---

## Dynamic Content (Jinja in Markdown)

The `jinja_md` plugin lets you use Jinja2 template syntax directly inside your Markdown files. This enables dynamic content generation: listing recent posts, showing related articles, creating series navigation, and more.

## Table of Contents

- [Overview](#overview)
- [Enabling Jinja in Markdown](#enabling-jinja-in-markdown)
- [Accessing Posts in Markdown](#accessing-posts-in-markdown)
- [Including Recent Posts from a Feed](#including-recent-posts-from-a-feed)
- [Common Use Cases](#common-use-cases)
- [Available Variables](#available-variables)
- [Available Filters](#available-filters)
- [Practical Examples](#practical-examples)
- [Tips and Troubleshooting](#tips-and-troubleshooting)

---

## Overview

The `jinja_md` plugin processes Jinja2 template syntax in your Markdown content **before** it's converted to HTML. This means you can:

- Loop through posts and create dynamic lists
- Filter posts by tags, dates, or any frontmatter field
- Access site configuration and post metadata
- Create reusable content patterns within your posts

The plugin uses the same template engine as HTML templates, so all filters and syntax work identically.

---

## Enabling Jinja in Markdown

The `jinja_md` plugin is **included by default**, but Jinja processing is **opt-in per post**. You must explicitly enable it in the frontmatter:

```yaml
---
title: "My Homepage"
jinja: true
---

## Recent Posts

{% for post in core.filter("published == True")[:5] %}
- [{{ post.Title }}]({{ post.Href }})
{% endfor %}
```

### Why Opt-In?

Jinja processing is opt-in because:

1. **Performance** - Not every post needs template processing
2. **Safety** - Prevents accidental template interpretation in code examples
3. **Clarity** - Makes it obvious which files contain dynamic content

### Valid Values

The `jinja` frontmatter field accepts:

```yaml
jinja: true      # Enable (boolean)
jinja: "true"    # Enable (string)
jinja: "yes"     # Enable (string)
jinja: "1"       # Enable (string)
jinja: false     # Disable (default)
```

---

## Accessing Posts in Markdown

When Jinja is enabled, you have access to three key objects for working with posts:

### `core.Posts()`

Returns **all posts** as a list:

```jinja2
{% for post in core.Posts() %}
- {{ post.Title }}
{% endfor %}
```

### `core.filter(expression)`

Returns posts matching a **filter expression**:

```jinja2
{% for post in core.filter("published == True") %}
- {{ post.Title }}
{% endfor %}
```

The filter syntax matches what you use in feed configuration:

| Expression | Description |
|------------|-------------|
| `"published == True"` | Published posts |
| `"draft == False"` | Non-draft posts |
| `"'python' in tags"` | Posts tagged with "python" |
| `"date >= '2024-01-01'"` | Posts from 2024 onwards |
| `"published == True and 'tutorial' in tags"` | Published tutorials |

### `core.feeds`

Access to configured feed collections (coming soon):

```jinja2
{% for post in core.feeds.blog.posts[:5] %}
- {{ post.Title }}
{% endfor %}
```

---

## Including Recent Posts from a Feed

This is the most common use case: showing a list of recent posts on a homepage or landing page.

### Basic Example: Last 3 Posts as Links

```markdown
---
title: "Welcome to My Blog"
jinja: true
---

# Welcome

Check out my latest articles:

## Recent Posts

{% for post in core.filter("published == True")[:3] %}
- [{{ post.Title }}]({{ post.Href }})
{% endfor %}
```

**Output:**

```markdown
# Welcome

Check out my latest articles:

## Recent Posts

- [Getting Started with Go](/getting-started-with-go/)
- [Why I Switched to Neovim](/why-i-switched-to-neovim/)
- [Building a Static Site Generator](/building-a-static-site-generator/)
```

### With Dates

```markdown
## Recent Posts

{% for post in core.filter("published == True")[:5] %}
- [{{ post.Title }}]({{ post.Href }}) - {{ post.Date|date_format:"Jan 2, 2006" }}
{% endfor %}
```

### With Descriptions

```markdown
## Recent Posts

{% for post in core.filter("published == True")[:3] %}
### [{{ post.Title }}]({{ post.Href }})

{{ post.Description|default_if_none:"" }}

{% endfor %}
```

### Excluding the Current Post

When showing related posts, exclude the current post:

```markdown
## Other Posts You Might Like

{% for post in core.filter("published == True") %}
{% if post.Slug != slug %}
- [{{ post.Title }}]({{ post.Href }})
{% endif %}
{% endfor %}
```

---

## Common Use Cases

### Related Posts by Tag

Show posts that share tags with the current post:

```markdown
---
title: "Python Tips and Tricks"
tags: ["python", "programming"]
jinja: true
---

# Python Tips and Tricks

[Your article content here]

---

## Related Posts

{% for p in core.filter("published == True") %}
{% if p.Slug != post.Slug %}
{% for tag in post.Tags %}
{% if tag in p.Tags %}
- [{{ p.Title }}]({{ p.Href }})
{% endif %}
{% endfor %}
{% endif %}
{% endfor %}
```

### Posts in a Series

Create navigation for a multi-part series:

```markdown
---
title: "Go Tutorial Part 2: Variables"
series: "go-tutorial"
series_order: 2
jinja: true
---

## Series Navigation

{% for p in core.filter("published == True") %}
{% if p.Extra.series == "go-tutorial" %}
{% if p.Extra.series_order == post.Extra.series_order %}
**{{ p.Extra.series_order }}. {{ p.Title }}** (You are here)
{% else %}
{{ p.Extra.series_order }}. [{{ p.Title }}]({{ p.Href }})
{% endif %}
{% endif %}
{% endfor %}
```

### Featured Posts Section

Show posts marked as featured:

```markdown
---
title: "Home"
jinja: true
---

## Featured Articles

{% for post in core.filter("published == True and featured == True")[:3] %}
### [{{ post.Title }}]({{ post.Href }})

{{ post.Description|default_if_none:"" }}

{% endfor %}
```

### Table of Contents from Other Posts

Generate a table of contents from a category:

```markdown
---
title: "Documentation Index"
jinja: true
---

# Documentation

## Getting Started

{% for post in core.filter("published == True and 'getting-started' in tags") %}
- [{{ post.Title }}]({{ post.Href }})
{% endfor %}

## Advanced Topics

{% for post in core.filter("published == True and 'advanced' in tags") %}
- [{{ post.Title }}]({{ post.Href }})
{% endfor %}

## API Reference

{% for post in core.filter("published == True and 'api' in tags") %}
- [{{ post.Title }}]({{ post.Href }})
{% endfor %}
```

### Dynamic Navigation

Create a navigation section that updates automatically:

```markdown
---
title: "Python Guide"
category: "python"
jinja: true
---

<nav class="sidebar">
{% for p in core.filter("published == True and category == 'python'") %}
{% if p.Slug == post.Slug %}
<strong>{{ p.Title }}</strong>
{% else %}
<a href="{{ p.Href }}">{{ p.Title }}</a>
{% endif %}
{% endfor %}
</nav>

# Python Guide

[Content here]
```

---

## Available Variables

When Jinja is enabled in a Markdown file, these variables are available:

### `post` - The Current Post

| Field | Type | Description |
|-------|------|-------------|
| `post.Title` | string | Post title |
| `post.Slug` | string | URL slug |
| `post.Href` | string | Relative URL (e.g., `/my-post/`) |
| `post.Date` | time | Publication date |
| `post.Published` | bool | Whether the post is published |
| `post.Draft` | bool | Whether it's a draft |
| `post.Tags` | []string | List of tags |
| `post.Description` | string | Post description |
| `post.Content` | string | Raw Markdown content |
| `post.Extra` | map | Custom frontmatter fields |

Access custom frontmatter via `post.Extra`:

```jinja2
Series: {{ post.Extra.series }}
Order: {{ post.Extra.series_order }}
Custom Field: {{ post.Extra.my_custom_field }}
```

### `config` - Site Configuration

| Field | Type | Description |
|-------|------|-------------|
| `config.Title` | string | Site title |
| `config.Description` | string | Site description |
| `config.URL` | string | Site base URL |
| `config.Author` | string | Site author |

### `core` - The Lifecycle Manager

| Method | Returns | Description |
|--------|---------|-------------|
| `core.Posts()` | []Post | All posts |
| `core.filter(expr)` | []Post | Posts matching expression |

### `posts` - All Posts (Alias)

A convenience alias for `core.Posts()`:

```jinja2
{% for p in posts %}
- {{ p.Title }}
{% endfor %}
```

### Shorthand Aliases

These aliases are also available for convenience:

| Alias | Equivalent |
|-------|------------|
| `site_title` | `config.Title` |
| `site_url` | `config.URL` |
| `site_description` | `config.Description` |
| `site_author` | `config.Author` |

---

## Available Filters

All template filters work in Jinja Markdown. Here are the most useful ones:

### Date Formatting

```jinja2
{{ post.Date|date_format:"January 2, 2006" }}
{{ post.Date|date_format:"2006-01-02" }}
{{ post.Date|date_format:"Jan 2" }}
```

Go date format reference:

| Format | Example |
|--------|---------|
| `2006-01-02` | 2024-03-15 |
| `January 2, 2006` | March 15, 2024 |
| `Jan 2, 2006` | Mar 15, 2024 |
| `02 Jan 2006` | 15 Mar 2024 |
| `Monday, January 2` | Friday, March 15 |

### String Manipulation

```jinja2
{{ post.Title|slugify }}           {# my-post-title #}
{{ post.Description|truncate:100 }} {# First 100 chars... #}
{{ post.Description|truncatewords:20 }} {# First 20 words... #}
{{ "Hello World"|upper }}          {# HELLO WORLD #}
{{ "Hello World"|lower }}          {# hello world #}
```

### Default Values

```jinja2
{{ post.Description|default_if_none:"No description" }}
{{ post.Extra.custom|default_if_none:"Default value" }}
```

### Collections

```jinja2
{{ post.Tags|length }}             {# Number of tags #}
{{ post.Tags|join:", " }}          {# tag1, tag2, tag3 #}
{{ post.Tags|first }}              {# First tag #}
{{ post.Tags|last }}               {# Last tag #}
{{ post.Tags|sort|join:", " }}     {# Alphabetically sorted #}
```

### HTML/Text

```jinja2
{{ html_content|striptags }}       {# Remove HTML tags #}
{{ text|linebreaksbr }}            {# Convert \n to <br> #}
```

### URLs

```jinja2
{{ post.Href|absolute_url:config.URL }}  {# https://example.com/my-post/ #}
{{ "search query"|urlencode }}           {# search%20query #}
```

---

## Practical Examples

### Blog Homepage

A complete homepage showing featured posts, recent posts, and categories:

```markdown
---
title: "Home"
template: "landing.html"
jinja: true
---

# Welcome to My Blog

I write about programming, technology, and life.

## Featured Posts

{% for p in core.filter("published == True and featured == True")[:3] %}
<article class="featured">
  <h3><a href="{{ p.Href }}">{{ p.Title }}</a></h3>
  <p>{{ p.Description|default_if_none:"" }}</p>
  <time>{{ p.Date|date_format:"January 2, 2006" }}</time>
</article>
{% endfor %}

## Recent Posts

{% for p in core.filter("published == True")[:10] %}
- [{{ p.Title }}]({{ p.Href }}) <small>{{ p.Date|date_format:"Jan 2" }}</small>
{% endfor %}

[View all posts](/archive/)
```

### Series Navigation Component

Reusable series navigation that appears at the top and bottom of series posts:

```markdown
---
title: "Building a CLI in Go: Part 3"
series: "go-cli"
series_order: 3
jinja: true
---

<nav class="series-nav">
  <strong>Series: Building a CLI in Go</strong>
  <ol>
  {% for p in core.filter("published == True") %}
  {% if p.Extra.series == post.Extra.series %}
    {% if p.Slug == post.Slug %}
    <li class="current">{{ p.Title }}</li>
    {% else %}
    <li><a href="{{ p.Href }}">{{ p.Title }}</a></li>
    {% endif %}
  {% endif %}
  {% endfor %}
  </ol>
</nav>

# Building a CLI in Go: Part 3

[Your content here]

---

<nav class="series-nav">
  {% set series_posts = [] %}
  {% for p in core.filter("published == True") %}
  {% if p.Extra.series == post.Extra.series %}
  {% set _ = series_posts.append(p) %}
  {% endif %}
  {% endfor %}

  {% if post.Extra.series_order > 1 %}
  [Previous: Part {{ post.Extra.series_order - 1 }}](#)
  {% endif %}

  {% if post.Extra.series_order < series_posts|length %}
  [Next: Part {{ post.Extra.series_order + 1 }}](#)
  {% endif %}
</nav>
```

### Tag Cloud

Generate a tag cloud showing all tags:

```markdown
---
title: "Tags"
jinja: true
---

# All Tags

<div class="tag-cloud">
{% set all_tags = [] %}
{% for p in core.filter("published == True") %}
  {% for tag in p.Tags %}
    {% if tag not in all_tags %}
      {% set _ = all_tags.append(tag) %}
    {% endif %}
  {% endfor %}
{% endfor %}

{% for tag in all_tags|sort %}
<a href="/tags/{{ tag|slugify }}/" class="tag">{{ tag }}</a>
{% endfor %}
</div>
```

### Archive by Year

Create an archive organized by year:

```markdown
---
title: "Archive"
jinja: true
---

# Archive

{% set current_year = "" %}
{% for p in core.filter("published == True") %}
{% set year = p.Date|date_format:"2006" %}
{% if year != current_year %}
{% set current_year = year %}

## {{ year }}

{% endif %}
- [{{ p.Title }}]({{ p.Href }}) <small>{{ p.Date|date_format:"Jan 2" }}</small>
{% endfor %}
```

### Posts by Category

Group posts by a custom category field:

```markdown
---
title: "All Articles"
jinja: true
---

# All Articles

## Tutorials

{% for p in core.filter("published == True and category == 'tutorial'") %}
- [{{ p.Title }}]({{ p.Href }})
{% endfor %}

## Essays

{% for p in core.filter("published == True and category == 'essay'") %}
- [{{ p.Title }}]({{ p.Href }})
{% endfor %}

## Reviews

{% for p in core.filter("published == True and category == 'review'") %}
- [{{ p.Title }}]({{ p.Href }})
{% endfor %}
```

### Dynamic Changelog

Show recent posts as a changelog:

```markdown
---
title: "Changelog"
jinja: true
---

# Changelog

{% for p in core.filter("published == True and 'changelog' in tags") %}
## {{ p.Date|date_format:"January 2, 2006" }} - {{ p.Title }}

{{ p.Description|default_if_none:"" }}

[Read more]({{ p.Href }})

---

{% endfor %}
```

---

## Tips and Troubleshooting

### Escaping Jinja Syntax

If you need to show literal Jinja syntax (like in documentation), use raw blocks:

```jinja2
{% raw %}
This won't be processed: {{ variable }}
Neither will this: {% for item in list %}
{% endraw %}
```

### Debugging

Print variables to see their values:

```jinja2
<!-- Debug: {{ post.Extra }} -->
<!-- Tags: {{ post.Tags }} -->
<!-- Post count: {{ core.Posts()|length }} -->
```

### Performance Considerations

- Use specific filters rather than filtering all posts in templates
- Limit results with slicing: `[:10]` instead of processing everything
- Consider using feeds for complex aggregations

### Common Mistakes

**Forgetting `jinja: true`:**

```yaml
---
title: "My Post"
jinja: true  # Don't forget this!
---
```

**Wrong field names:**

```jinja2
{{ post.title }}  # Wrong: lowercase
{{ post.Title }}  # Correct: PascalCase
```

**Missing safe filter for HTML:**

```jinja2
{{ post.ArticleHTML }}        # Escaped HTML (wrong)
{{ post.ArticleHTML|safe }}   # Raw HTML (correct)
```

**Comparing to string instead of boolean:**

```jinja2
{% if post.Published == "true" %}  # Wrong
{% if post.Published %}             # Correct
{% if post.Published == True %}     # Also correct
```

---

## See Also

- [[templates-guide|Templates Guide]] - Full template syntax and filters
- [[feeds-guide|Feeds Guide]] - Configure feed collections
- [[configuration-guide|Configuration Guide]] - Site-wide settings
