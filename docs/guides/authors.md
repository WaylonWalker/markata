---
title: "Multi-Author Support"
description: "Configure multiple authors with roles, bios, and avatars for your blog"
date: 2026-02-12
published: true
tags:
  - authors
  - configuration
  - frontmatter
---

Markata-go supports multiple authors with rich metadata including roles, bios, avatars, and social links. This guide covers how to configure authors and reference them in your posts.

## Quick Start

### 1. Configure Authors

Add authors to your `markata-go.toml`:

```toml
[markata-go.authors]

[markata-go.authors.authors.waylon]
name = "Waylon Walker"
role = "author"
url = "https://waylonwalker.com"
avatar = "/images/waylon.jpg"
bio = "Python and Go developer"
active = true
default = true

[markata-go.authors.authors.guest]
name = "Guest Writer"
role = "editor"
bio = "Occasional contributor"
guest = true
active = true
```

### 2. Reference Authors in Posts

```yaml
---
title: My Post
authors:
  - waylon
  - guest
---
```

## Author Configuration

### Author Fields

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | **Required** - Display name |
| `role` | string | Simple role label (author, editor, etc.) |
| `bio` | string | Short biography |
| `avatar` | string | Avatar image URL or path |
| `url` | string | Personal website URL |
| `email` | string | Contact email |
| `social` | map | Social media links |
| `default` | bool | Default author for posts without explicit authors |
| `guest` | bool | Guest author flag |
| `active` | bool | Whether author is currently active |

### Default Author

Set `default = true` on one author. Posts without an explicit `authors` frontmatter field will automatically use this author:

```toml
[markata-go.authors.authors.waylon]
name = "Waylon Walker"
default = true
```

### Author Pages

Generate individual author profile pages:

```toml
[markata-go.authors]
generate_pages = true
url_pattern = "/authors/{id}/"
```

## Frontmatter

### Simple Format

Reference authors by ID:

```yaml
---
title: My Post
authors:
  - waylon
  - guest
---
```

### Extended Format

Specify per-post roles and details:

```yaml
---
title: Collaborative Post
authors:
  - id: waylon
    role: author
    details: wrote the introduction
  - id: codex
    role: pair programmer
    details: wrote the code examples
---
```

### Key Aliases

Use convenient aliases for common keys:

```yaml
---
authors:
  - name: waylon     # same as id: waylon
    title: author    # same as role: author
    description: wrote it  # same as details: ...
---
```

### Legacy Format

The older single-author format still works:

```yaml
---
title: My Post
author: waylon
---
```

## Template Access

In templates, access author data:

```html
{% for author in post.author_objects %}
  <span>{{ author.name }}</span>
  {% if author.role %}({{ author.role }}){% endif %}
{% endfor %}
```

Available author properties:
- `id`, `name`, `bio`, `email`, `avatar`, `url`
- `role` - Simple role label
- `role_display` - Pre-computed display string
- `guest`, `active`, `default`
- `social` - Map of platform -> URL
- `contributions` - CReDiT taxonomy roles
- `details` - Per-post contribution details

## Role System

### Simple Roles

Use the `role` field for blog-friendly roles:

```toml
[markata-go.authors.authors.waylon]
name = "Waylon Walker"
role = "author"
```

Common roles: `author`, `editor`, `reviewer`, `photographer`, `illustrator`, `translator`

### CReDiT Roles

For academic-style contribution tracking, use `contributions`:

```toml
[markata-go.authors.authors.waylon]
name = "Waylon Walker"
contributions = ["Conceptualization", "Software", "Writing"]
```

## Examples

### Multiple Authors with Different Roles

```yaml
---
title: Building a CLI Tool
authors:
  - id: waylon
    role: author
    details: wrote the core implementation
  - id: kimmi
    role: reviewer
    details: reviewed and tested
---
```

### Mixed Author Types

```yaml
---
title: Guest Post
authors:
  - waylon  # uses default role from config
  - id: guest-writer
    role: guest author
    details: wrote this guest post
---
```

## Editorial credits and revision history

The default byline keeps writers prominent and groups editing/reviewing roles
below them. The first editor is displayed by name; additional editors are
stacked behind a disclosure with all contributor names and roles.

Use [Editorial history](../editorial-history/) to show recorded reviews,
GitHub file diffs, and lazily loaded inline source changes. Humans and agents
can use the [editorial authoring commands](../editorial-authoring/) to commit
selected edits and record their attribution automatically.
