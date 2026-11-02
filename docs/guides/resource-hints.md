---
title: "Resource Hints"
description: "Automatically generate resource hints (dns-prefetch, preconnect, preload, prefetch) to improve page load performance"
date: 2025-01-25
published: true
tags:
  - documentation
  - performance
  - resource-hints
  - optimization
---


The Resource Hints plugin automatically generates and injects resource hints (`dns-prefetch`, `preconnect`, `preload`, `prefetch`) into your HTML pages to improve performance.

## Overview

Resource hints help browsers prepare for external resources before they're needed, improving page load performance by:

- **DNS Prefetch**: Performing DNS lookups in advance
- **Preconnect**: Establishing full connections (DNS + TCP + TLS) early
- **Preload**: Fetching critical resources before the browser discovers them
- **Prefetch**: Loading resources for future navigation

## Page-Specific Hints

As of the latest version, each page gets resource hints **only for external domains detected on that specific page**, rather than a site-wide list.

### Benefits

- ✅ **Reduced HTML bloat** - Typically 5-20 hints per page instead of 400+
- ✅ **Better performance** - Browsers honor relevant hints (ignore excess)
- ✅ **More accurate** - Hints match actual page needs
- ✅ **Cleaner code** - No irrelevant dns-prefetch pollution

### Before vs After

**Before** (site-wide approach):
- Scanned all pages, collected all domains
- Generated one comprehensive list
- Injected same 400+ hints everywhere

**After** (page-specific approach):
- Each page scanned individually
- Hints generated per page
- Only relevant domains included

## Configuration

Add resource hints configuration to your `markata.toml`:

```toml
[resource_hints]
# Enable/disable resource hints generation (default: true)
enabled = true

# Enable auto-detection of external domains in HTML/CSS (default: true)
auto_detect = true

# Exclude specific domains from auto-detection
# Useful for pages with many external links (blogrolls, reader pages, etc.)
exclude_domains = [
    "example.com",
    "reader.waylonwalker.com",  # Exclude links from your reader
]

# Manually configure hints for specific domains
[[resource_hints.domains]]
domain = "cdn.jsdelivr.net"
hint_types = ["dns-prefetch"]

[[resource_hints.domains]]
domain = "fonts.googleapis.com"
hint_types = ["preconnect"]
crossorigin = "anonymous"
```

## Hint Types

### DNS Prefetch

Performs DNS lookup in advance with minimal overhead. Good for third-party resources.

```toml
[[resource_hints.domains]]
domain = "cdn.jsdelivr.net"
hint_types = ["dns-prefetch"]
```

**When to use**:
- Third-party CDNs
- Analytics domains
- Social media widgets
- Any domain not critical but likely to be used

### Preconnect

Establishes full connection (DNS + TCP + TLS handshake). Higher overhead but faster when resource is actually needed.

```toml
[[resource_hints.domains]]
domain = "fonts.googleapis.com"
hint_types = ["preconnect"]
crossorigin = "anonymous"
```

**When to use**:
- Critical resources (fonts, critical CSS/JS)
- Domains you know will be used immediately
- A media origin only when the media starts immediately; YouTube and Vimeo
  discovered from page links use `dns-prefetch` so a link alone does not open a
  TCP/TLS connection
- Limit to 3-5 per page for best results

### Preload

Fetches specific resources early. Requires the `as` attribute.

```toml
[[resource_hints.domains]]
domain = "fonts.gstatic.com"
hint_types = ["preload"]
as = "font"
crossorigin = "anonymous"
```

**When to use**:
- Critical fonts
- Above-the-fold images
- Essential scripts/styles

### Prefetch

Low-priority fetch for future navigation.

```toml
[[resource_hints.domains]]
domain = "next-page-cdn.com"
hint_types = ["prefetch"]
```

**When to use**:
- Resources for next likely page
- Assets for common navigation paths

## Excluding Domains

If you have pages with many external links (blogroll, reader page), exclude those domains from auto-detection:

```toml
[resource_hints]
exclude_domains = [
    "reader.waylonwalker.com",  # Don't generate hints for blogroll links
    "twitter.com",              # Social media often not critical
    "github.com",               # Code hosting links
]
```

### Why Exclude?

Pages like blogrolls or RSS readers may have hundreds of external links. These domains:
- Are rarely visited from the current page
- Would generate excessive hints (400+)
- Waste browser resources on unused DNS lookups

**Example**: Your archive page shows posts that link to many blogs. Those blog domains shouldn't generate hints because users don't typically visit them from the archive.

## Performance Impact

### Build Time

- **Impact**: Minimal, typically +100-200ms for 100-500 page sites
- **Scaling**: Linear with number of pages
- **Trade-off**: Slightly longer builds for significantly better UX

### User Experience

- **HTML Size**: ~20KB less per page (without irrelevant hints)
- **DNS Lookups**: Only relevant domains, better cache utilization
- **Page Load**: Faster due to early DNS resolution of actually-used domains
- **Browser Efficiency**: Browsers honor reasonable hint counts (10-50)

## Best Practices

### 1. Use DNS Prefetch for Most Domains

Low overhead, good for any external resource:

```toml
[[resource_hints.domains]]
domain = "cdn.jsdelivr.net"
hint_types = ["dns-prefetch"]
```

### 2. Reserve Preconnect for Critical Origins

Limit to 3-5 per page:

```toml
[[resource_hints.domains]]
domain = "fonts.googleapis.com"
hint_types = ["preconnect"]

[[resource_hints.domains]]
domain = "cdn.critical-resource.com"
hint_types = ["preconnect"]
```

### 3. Exclude Blogroll/Reader Links

```toml
[resource_hints]
exclude_domains = ["reader.waylonwalker.com"]
```

### 4. Monitor Hint Counts

Each page should have < 20 hints for optimal results. Check your HTML:

```bash
# Count hints on a page
grep -c 'rel="dns-prefetch"' output/archive/index.html
```

### 5. Combine with Other Optimizations

- Minimize external dependencies
- Self-host critical resources when possible
- Use CDNs with good performance

## Example Output

### Typical Blog Post

```html
<head>
  <meta charset="UTF-8">

  <!-- Auto-generated resource hints -->
  <link rel="dns-prefetch" href="https://cdn.jsdelivr.net">
  <link rel="dns-prefetch" href="https://www.youtube.com">
  <link rel="dns-prefetch" href="https://github.com">
  <link rel="dns-prefetch" href="https://i.ytimg.com">
  <!-- End resource hints -->

  <title>Blog Post</title>
  ...
</head>
```

### Archive Page (After Fix)

```html
<head>
  <meta charset="UTF-8">

  <!-- Auto-generated resource hints -->
  <link rel="dns-prefetch" href="https://cdn.jsdelivr.net">
  <link rel="dns-prefetch" href="https://cdn.tailwindcss.com">
  <!-- End resource hints -->

  <title>Archive</title>
  ...
</head>
```

## Troubleshooting

### Too Many Hints on a Page?

**Symptoms**: Page has 50+ hints

**Solutions**:
1. Check if page has many external links (blogroll, reader)
2. Add those domains to `exclude_domains`
3. Consider disabling auto-detection for specific content types

```toml
[resource_hints]
exclude_domains = [
    "reader.waylonwalker.com",
    "blogroll.example.com",
]
```

### Hints Not Showing Up?

**Check**:
1. `enabled = true` in config
2. `auto_detect = true` in config
3. Page has external URLs in `href` or `src` attributes
4. Page has `<head>` tag
5. Build logs for errors

### Build Time Increased?

**Expected**: +100-200ms for typical sites

**If excessive** (> 1 second):
1. Check number of pages (scales linearly)
2. Consider disabling auto-detection:

```toml
[resource_hints]
auto_detect = false
```

3. Manually configure critical domains only

### Duplicate Hints?

**Symptoms**: Multiple hints for same domain

**Cause**: Both manual config and auto-detection enabled

**Solution**: Choose one approach:

```toml
# Option 1: Auto-detect only
[resource_hints]
auto_detect = true

# Option 2: Manual only
[resource_hints]
auto_detect = false

[[resource_hints.domains]]
domain = "cdn.example.com"
hint_types = ["dns-prefetch"]
```

## Advanced Usage

### Per-Domain Hint Types

Some domains benefit from multiple hint types:

```toml
[[resource_hints.domains]]
domain = "fonts.googleapis.com"
hint_types = ["dns-prefetch", "preconnect"]
```

This generates both:
```html
<link rel="dns-prefetch" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.googleapis.com">
```

### Cross-Origin Hints

For resources requiring CORS:

```toml
[[resource_hints.domains]]
domain = "fonts.gstatic.com"
hint_types = ["preconnect"]
crossorigin = "anonymous"
```

Generates:
```html
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
```

### Preload with Type

For specific resource types:

```toml
[[resource_hints.domains]]
domain = "cdn.example.com"
hint_types = ["preload"]
as = "font"
crossorigin = "anonymous"
```

Generates:
```html
<link rel="preload" href="https://cdn.example.com" as="font" crossorigin>
```

## Migration from Site-Wide Hints

If you're upgrading from an older version:

1. **No action required** - Works automatically
2. **Optional**: Add `exclude_domains` for fine-tuning
3. **Rebuild** to see reduced hint counts
4. **Verify** pages load correctly

Your site will automatically benefit from:
- Smaller HTML files
- Faster page loads
- More relevant resource hints

## Related Documentation

- [[performance|Performance Guide]] - Overall site optimization
- [[configuration|Configuration Reference]] - All config options
- [[plugin-development|Plugin Development]] - Extending markata-go

## See Also

- [Resource Hints Spec](https://www.w3.org/TR/resource-hints/)
- [DNS Prefetch](https://developer.mozilla.org/en-US/docs/Web/Performance/dns-prefetch)
- [Preconnect](https://developer.mozilla.org/en-US/docs/Web/Performance/Speculative_loading#preconnect)
