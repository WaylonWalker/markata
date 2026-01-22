# Static Redirects Specification

Static redirects enable URL migrations and content reorganization without losing traffic or breaking bookmarks. This plugin generates HTML redirect pages from a simple configuration file, compatible with static hosting platforms.

## Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                        REDIRECTS WORKFLOW                            │
├─────────────────────────────────────────────────────────────────────┤
│  1. READ                                                             │
│     - Load _redirects file from configured location                  │
│     - Parse redirect rules (source → destination)                    │
│                                                                      │
│  2. FILTER                                                           │
│     - Skip comments (lines starting with #)                          │
│     - Skip wildcard patterns (contain *)                             │
│     - Skip malformed entries                                         │
│                                                                      │
│  3. GENERATE                                                         │
│     - For each redirect, create index.html at source path            │
│     - Use meta refresh + canonical link for redirect                 │
│     - Apply template with styling from config                        │
│                                                                      │
│  4. CACHE                                                            │
│     - Hash redirects file content                                    │
│     - Skip regeneration if unchanged                                 │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Configuration

### Basic Configuration

```toml
[name.redirects]
redirects_file = "static/_redirects"
```

### Full Configuration

```toml
[name.redirects]
# Path to redirects file (relative to project root)
redirects_file = "static/_redirects"

# Custom template for redirect pages (optional)
# If not specified, uses built-in default template
redirect_template = "templates/redirect.html"
```

---

## Redirects File Format

The redirects file uses a simple space-separated format, compatible with Cloudflare Pages and Netlify.

### Basic Syntax

```text
# Comments start with #
/old-path    /new-path
/legacy-url  /current-url
```

### Rules

| Rule | Description |
|------|-------------|
| Comments | Lines starting with `#` are ignored |
| Format | `<source> <destination>` separated by whitespace |
| Paths | Must start with `/` |
| Wildcards | Patterns containing `*` are skipped (not supported for static generation) |

### Example Redirects File

```text
# Blog reorganization (2024)
/blog/old-post    /posts/new-post
/articles         /blog

# Legacy URLs
/about-me         /about
/contact-us       /contact

# Renamed sections
/tutorials/python-basics    /learn/python/getting-started
/tutorials/rust-intro       /learn/rust/introduction

# These are skipped (wildcards not supported):
# /old-blog/*    /blog/*
# /api/*         /v2/api/*
```

---

## Redirect Model

### Redirect Object

| Field | Type | Description |
|-------|------|-------------|
| `original` | string | Source path (the old URL) |
| `new` | string | Destination path (the new URL) |

---

## Generated Output

### Output Structure

For each redirect rule, an `index.html` file is created:

```
output/
├── old-path/
│   └── index.html          # Redirects to /new-path
├── legacy-url/
│   └── index.html          # Redirects to /current-url
└── about-me/
    └── index.html          # Redirects to /about
```

### Default Template

The default redirect template provides:

1. **Instant redirect** via `<meta http-equiv="Refresh">`
2. **SEO-friendly** canonical link to new URL
3. **Fallback content** for users/bots that don't follow redirects
4. **Styled page** using site configuration colors

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta http-equiv="Refresh" content="0; url='{{ new }}'" />
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <link rel="canonical" href="{{ new }}" />
  <meta name="description" content="{{ original }} has been moved to {{ new }}." />
  <title>{{ original }} has been moved to {{ new }}</title>
  <style>
    html {
      font-family: system-ui, sans-serif;
      background: {{ config.style.color_bg | default('#1f2022') }};
      color: {{ config.style.color_text | default('#eefbfe') }};
    }
    body {
      margin: 5rem auto;
      max-width: {{ config.style.body_width | default('800px') }};
    }
    a {
      color: {{ config.style.color_link | default('#fb30c4') }};
      text-decoration-color: {{ config.style.color_accent | default('#e1bd00c9') }};
    }
  </style>
</head>
<body>
  <h1>Page Moved</h1>
  <p>
    <code>{{ original }}</code> has moved to 
    <a href="{{ new }}">{{ new }}</a>
  </p>
</body>
</html>
```

### Template Context

| Variable | Type | Description |
|----------|------|-------------|
| `original` | string | Original/source path |
| `new` | string | New/destination path |
| `config` | Config | Site configuration (for styling) |

---

## Behavior

### Processing Rules

1. **Skip comments**: Lines starting with `#` are ignored
2. **Skip wildcards**: Patterns containing `*` are not processed (static sites can't handle dynamic wildcards)
3. **Skip malformed**: Lines with fewer than 2 parts are skipped
4. **Trim paths**: Leading/trailing whitespace is removed from paths

### Caching

The plugin caches based on the content of the redirects file:

```python
key = hash("redirects", raw_redirects_content)
if cache.get(key) == "done":
    return  # Skip regeneration
```

This ensures redirects are only regenerated when the `_redirects` file changes.

### Directory Creation

For each redirect, the plugin:

1. Creates the parent directory structure
2. Writes `index.html` inside the source path directory

Example: `/old/nested/path` → creates `output/old/nested/path/index.html`

---

## Hook Specification

### Stage

`save`

### Hook Signature

```python
@hook_impl
def save(core):
    config = core.config.redirects
    redirects_file = Path(config.redirects_file)
    
    if not redirects_file.exists():
        return
    
    raw_redirects = redirects_file.read_text().split("\n")
    
    # Cache check
    key = core.make_hash("redirects", raw_redirects)
    if core.cache.get(key) == "done":
        return
    
    # Parse redirects
    redirects = []
    for line in raw_redirects:
        line = line.strip()
        if not line or line.startswith("#") or "*" in line:
            continue
        parts = line.split()
        if len(parts) >= 2:
            redirects.append(Redirect(original=parts[0], new=parts[1]))
    
    # Load template
    if config.redirect_template:
        template = load_template(config.redirect_template)
    else:
        template = load_default_redirect_template()
    
    # Generate redirect pages
    for redirect in redirects:
        output_path = core.config.output_dir / redirect.original.strip("/") / "index.html"
        output_path.parent.mkdir(parents=True, exist_ok=True)
        output_path.write_text(template.render(
            original=redirect.original,
            new=redirect.new,
            config=core.config
        ))
    
    core.cache.set(key, "done")
```

---

## Configuration Model

```python
class RedirectsConfig(pydantic.BaseModel):
    redirects_file: Path = Path("static/_redirects")
    redirect_template: Optional[Path] = None

class Config(pydantic.BaseModel):
    redirects: RedirectsConfig = RedirectsConfig()
```

---

## Platform Compatibility

### Supported Platforms

| Platform | Native Support | Static Fallback |
|----------|---------------|-----------------|
| Cloudflare Pages | Yes (`_redirects` file) | Yes (HTML redirects) |
| Netlify | Yes (`_redirects` file) | Yes (HTML redirects) |
| Vercel | Yes (`vercel.json`) | Yes (HTML redirects) |
| GitHub Pages | No | Yes (HTML redirects) |
| S3/Static hosting | No | Yes (HTML redirects) |

### Why HTML Redirects?

While platforms like Cloudflare and Netlify support `_redirects` files natively, generating HTML redirect pages provides:

1. **Universal compatibility**: Works on any static host
2. **SEO preservation**: Canonical links maintain search rankings
3. **User experience**: Provides fallback content if redirect fails
4. **Debugging**: Easy to inspect and verify redirect targets

---

## Limitations

| Limitation | Reason |
|------------|--------|
| No wildcard support | Static HTML can't handle dynamic patterns |
| No status codes | HTML meta refresh is always a 302-equivalent |
| No query parameters | Query strings are not preserved in meta refresh |
| No conditional redirects | Static files can't evaluate conditions |

For advanced redirect needs (wildcards, status codes, conditions), use your hosting platform's native redirect features alongside this plugin.

---

## Examples

### Basic Blog Migration

```text
# _redirects
/blog/2023/post-one    /posts/post-one
/blog/2023/post-two    /posts/post-two
/blog/2024/new-post    /posts/new-post
```

### Section Reorganization

```text
# _redirects
/tutorials    /learn
/guides       /learn
/howto        /learn
/docs/api     /reference/api
/docs/cli     /reference/cli
```

### Shortened URLs

```text
# _redirects
/go/github     https://github.com/myorg/myproject
/go/discord    https://discord.gg/invite-code
/go/docs       /documentation
```

---

## Error Handling

| Error | Behavior |
|-------|----------|
| Missing redirects file | Skip silently (no redirects generated) |
| Malformed line | Skip line, continue processing |
| Write error | Log error, continue with other redirects |
| Template error | Use default template, log warning |

---

## See Also

- [SPEC.md](./SPEC.md) - Core specification
- [CONFIG.md](./CONFIG.md) - Configuration system
- [LIFECYCLE.md](./LIFECYCLE.md) - Build lifecycle (save stage)
