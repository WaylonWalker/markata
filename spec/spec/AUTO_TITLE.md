# Auto Title Specification

The auto title plugin automatically generates human-readable titles for posts that don't have a title specified in their frontmatter. It derives titles from filenames using simple text transformations.

## Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                       AUTO TITLE WORKFLOW                            │
├─────────────────────────────────────────────────────────────────────┤
│  INPUT                                                               │
│    posts/my-first-post.md     (no title in frontmatter)             │
│    posts/python_tips.md       (no title in frontmatter)             │
│    posts/hello-world.md       (title: "Hello!" in frontmatter)      │
│                                                                      │
│  PROCESS                                                             │
│    1. Filter posts where title == "" or title is None               │
│    2. Extract filename stem (without extension)                      │
│    3. Replace hyphens and underscores with spaces                   │
│    4. Apply title case                                               │
│                                                                      │
│  OUTPUT                                                              │
│    my-first-post.md  → title: "My First Post"                       │
│    python_tips.md    → title: "Python Tips"                         │
│    hello-world.md    → title: "Hello!" (unchanged - has title)      │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Behavior

### Title Generation Rules

| Step | Transformation | Example |
|------|----------------|---------|
| 1 | Extract filename stem | `my-first-post.md` → `my-first-post` |
| 2 | Replace hyphens with spaces | `my-first-post` → `my first post` |
| 3 | Replace underscores with spaces | `python_tips` → `python tips` |
| 4 | Apply title case | `my first post` → `My First Post` |

### Examples

| Filename | Generated Title |
|----------|----------------|
| `my-first-post.md` | My First Post |
| `python_tips_and_tricks.md` | Python Tips And Tricks |
| `hello-world.md` | Hello World |
| `2024-01-15-new-feature.md` | 2024 01 15 New Feature |
| `README.md` | Readme |
| `api-v2-reference.md` | Api V2 Reference |
| `getting_started.md` | Getting Started |

### When Title Is NOT Generated

The plugin skips posts that already have a title:

```yaml
---
title: "My Custom Title"
---
```

Or explicitly set to a non-empty value:

```yaml
---
title: Hello World!
---
```

### When Title IS Generated

The plugin processes posts where:

- `title` is not present in frontmatter
- `title` is empty string (`title: ""`)
- `title` is null (`title: null` or `title: ~`)

---

## Configuration

This plugin requires **no configuration**. It is enabled by default and processes all posts without titles automatically.

### Enabling (Default)

The plugin is included in the default plugin set:

```toml
[name]
hooks = ["default"]
```

### Disabling

To disable auto title generation:

```toml
[name]
hooks = ["default"]
disabled_hooks = ["auto_title"]
```

---

## Hook Specification

### Stage

`pre_render`

### Hook Signature

```python
@hook_impl
def pre_render(core):
    """
    Generate titles for posts that don't have one.
    
    Processes posts where title is empty or None, deriving
    a human-readable title from the filename.
    """
    for post in core.filter('title == "" or title == None'):
        filename_stem = Path(post.path).stem
        title = filename_stem.replace("-", " ").replace("_", " ").title()
        post.title = title
```

### Filter Expression

The plugin uses the filter expression:

```python
'title == "" or title == None'
```

This matches posts where:
- `title` equals empty string
- `title` equals None/null

---

## Post Model Impact

### Field Modified

| Field | Type | Description |
|-------|------|-------------|
| `title` | string | Set to generated title if empty/None |

### Before/After

**Before (in frontmatter):**
```yaml
---
slug: my-first-post
date: 2024-01-15
---
```

**After (in memory):**
```python
post.title = "My First Post"  # Derived from filename
post.slug = "my-first-post"
post.date = datetime(2024, 1, 15)
```

---

## Interaction with Other Plugins

### Slug Generation

If using a slug plugin that derives slugs from titles, auto_title should run **before** slug generation:

```
pre_render order:
  1. auto_title     → Sets title from filename
  2. slug_from_title → Could derive slug from title (if needed)
```

However, typically slugs are derived from filenames directly, making this ordering less critical.

### SEO/Meta Plugins

Auto-generated titles are available for:
- `<title>` tag generation
- Open Graph `og:title` meta tags
- RSS/Atom feed titles
- Search index entries

### Prev/Next Navigation

The `prevnext` plugin uses titles for navigation links, so auto_title should run before it.

---

## Edge Cases

### Date Prefixed Filenames

Files with date prefixes retain the date in the title:

| Filename | Title |
|----------|-------|
| `2024-01-15-new-post.md` | 2024 01 15 New Post |

To avoid this, either:
1. Set explicit titles in frontmatter
2. Use a separate plugin to strip date prefixes
3. Organize files by date in directories instead

### Special Characters

Characters other than hyphens and underscores are preserved:

| Filename | Title |
|----------|-------|
| `c++_tutorial.md` | C++ Tutorial |
| `node.js-guide.md` | Node.Js Guide |
| `q&a.md` | Q&A |

### Acronyms and Capitalization

Title case may not handle acronyms ideally:

| Filename | Generated | Preferred |
|----------|-----------|-----------|
| `api-reference.md` | Api Reference | API Reference |
| `html-basics.md` | Html Basics | HTML Basics |
| `aws-setup.md` | Aws Setup | AWS Setup |

For proper acronym handling, set titles explicitly in frontmatter.

---

## Implementation Notes

### Performance

- Runs once per build during `pre_render`
- Only processes posts matching filter (no title)
- No external dependencies
- No caching needed (fast string operations)

### Unicode Support

The plugin handles unicode filenames:

| Filename | Title |
|----------|-------|
| `café-guide.md` | Café Guide |
| `日本語-post.md` | 日本語 Post |

### Path Handling

Only the filename stem is used, not the full path:

| Path | Title |
|------|-------|
| `posts/2024/my-post.md` | My Post |
| `blog/tutorials/getting-started.md` | Getting Started |

---

## Comparison with Explicit Titles

| Approach | Pros | Cons |
|----------|------|------|
| **Auto Title** | No frontmatter needed, fast authoring | Less control, acronym issues |
| **Explicit Title** | Full control, proper formatting | More typing, can forget |

**Recommendation:** Use auto title for quick drafts and simple posts. Set explicit titles for published content where formatting matters.

---

## See Also

- [DATA_MODEL.md](./DATA_MODEL.md) - Post model specification
- [DEFAULT_PLUGINS.md](./DEFAULT_PLUGINS.md) - Default plugin set
- [LIFECYCLE.md](./LIFECYCLE.md) - Build lifecycle (pre_render stage)
