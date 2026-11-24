---
title: "Migrating from Python markata"
description: "Guide to migrating your site from Python markata to markata-go"
date: 2024-01-15
published: true
tags:
  - documentation
  - migration
  - getting-started
---


This guide helps you migrate your existing Python markata site to markata-go. The migration tool automates most of the process while highlighting areas that need manual attention.

## Quick Start

Run the migration tool to analyze your current configuration:

```bash
markata-go migrate
```

This will show a detailed report of:
- Configuration changes needed
- Filter expression migrations
- Template compatibility issues
- Warnings for unsupported features

## Migration Commands

### Full Migration

```bash
# Analyze and show migration report
markata-go migrate

# Analyze without writing (dry run)
markata-go migrate --dry-run

# Write migrated config to file
markata-go migrate -o markata-go.toml
```

### Config Migration Only

```bash
# Migrate config from markata.toml
markata-go migrate config

# Migrate from pyproject.toml
markata-go migrate config -i pyproject.toml -o markata-go.toml
```

### Filter Expression Check

```bash
# Check a specific filter expression
markata-go migrate filter "published == 'True'"

# Check all filters in your config
markata-go migrate filter
```

### Template Compatibility Check

```bash
# Check templates directory
markata-go migrate templates

# Check specific directory
markata-go migrate templates ./my-templates
```

## Configuration Changes

### Namespace

The configuration namespace changes from `[markata]` to `[markata-go]`:

```toml
# Before (Python markata)
[markata]
output = "public"

# After (markata-go)
[markata-go]
output_dir = "public"
```

### Key Renames

Several configuration keys have been renamed:

| Python markata | markata-go | Notes |
|----------------|------------|-------|
| `output` | `output_dir` | Output directory |
| `glob_patterns` | `patterns` | Under `[markata-go.glob]` |
| `author_name` | `author` | Site author |
| `site_name` | `title` | Site title |
| `site_description` | `description` | Site description |

### Navigation

Python markata uses a map for navigation:

```toml
# Python markata
[markata.nav]
home = "/"
blog = "/blog"
about = "/about"
```

markata-go uses an array with explicit labels:

```toml
# markata-go
[[markata-go.nav]]
label = "Home"
url = "/"

[[markata-go.nav]]
label = "Blog"
url = "/blog"

[[markata-go.nav]]
label = "About"
url = "/about"
```

## Filter Expression Changes

### Boolean Literals

Python markata often uses quoted boolean strings. markata-go uses unquoted booleans:

```python
# Python markata
filter = "published == 'True'"

# markata-go
filter = "published == True"
```

### `in` Operator

The `in` operator with lists must be converted to `or` expressions:

```python
# Python markata
filter = "templateKey in ['blog-post', 'til']"

# markata-go
filter = "templateKey == 'blog-post' or templateKey == 'til'"
```

### Operator Spacing

Ensure operators have surrounding whitespace:

```python
# Python markata (may work)
filter = "date<=today"

# markata-go (required)
filter = "date <= today"
```

### None Comparisons

Python-style `is None` must use `==`:

```python
# Python markata
filter = "image is None"
filter = "image is not None"

# markata-go
filter = "image == None"
filter = "image != None"
```

## Template Changes

markata-go uses pongo2, which is Jinja2-compatible with some differences.

### Variable Changes

| Python markata | markata-go | Notes |
|----------------|------------|-------|
| `post.markata.config` | `config` | Direct access |
| `post.markata.feeds` | `feeds` | Direct access |
| `post.article_html` | `post.content` | Renamed |

### Unsupported Features

The following Jinja2 features are not supported in pongo2:

| Feature | Alternative |
|---------|-------------|
| `{% macro %}` | Use `{% include %}` with variables |
| `{% do %}` | Use `{% set %}` |
| `{% call %}` | Restructure to use includes |
| `{% import %}` | Use `{% include %}` |
| List comprehensions | Pre-compute in Go |

### Filter Syntax

Python string methods should be converted to pongo2 filters:

```jinja2
{# Python markata #}
{{ title.lower() }}
{{ text.strip() }}

{# markata-go #}
{{ title|lower }}
{{ text|trim }}
```

## Step-by-Step Migration

### 1. Analyze Your Site

```bash
markata-go migrate --dry-run
```

Review the migration report carefully.

### 2. Migrate Configuration

```bash
markata-go migrate -o markata-go.toml
```

### 3. Update Templates

Check templates for compatibility:

```bash
markata-go migrate templates
```

Update any flagged issues manually.

### 4. Test the Build

```bash
markata-go build --dry-run
```

### 5. Full Build

```bash
markata-go build
```

### 6. Verify Output

Compare the output with your Python markata build to ensure everything looks correct.

## Common Issues

### "No markata config found"

The migration tool looks for:
- `markata.toml`
- `pyproject.toml` (with `[tool.markata]` section)
- `markata.yaml`

Specify your config file explicitly:

```bash
markata-go migrate -i my-config.toml
```

### Filter Validation Errors

If a migrated filter fails validation, you may need to simplify it. Complex Python expressions are not supported:

```python
# Not supported
filter = "len(tags) > 0 and published"

# Supported alternative
filter = "has_tags == True and published == True"
```

Add a `has_tags` field in your frontmatter or preprocessing.

### Template Macros

Convert macros to include files:

```jinja2
{# Before: macro in template #}
{% macro render_card(post) %}
  <div class="card">{{ post.title }}</div>
{% endmacro %}

{{ render_card(post) }}

{# After: separate include file #}
{# _card.html #}
<div class="card">{{ post.title }}</div>

{# main template #}
{% include "_card.html" %}
```

## Getting Help

- Check the [troubleshooting guide](/docs/troubleshooting)
- Open an issue on [GitHub](https://github.com/WaylonWalker/markata-go/issues)
- Join the discussion on the repository

## See Also

- [Configuration Guide](/docs/guides/configuration) - Full configuration reference
- [Templates Guide](/docs/guides/templates) - Template system documentation
- [Feeds Guide](/docs/guides/feeds) - Feed configuration
