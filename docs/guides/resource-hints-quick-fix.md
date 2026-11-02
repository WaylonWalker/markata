---
title: "Resource Hints Quick Fix"
description: "Quick reference for fixing excessive DNS-prefetch links on archive and blogroll pages"
date: 2025-01-25
published: true
tags:
  - documentation
  - quick-reference
  - resource-hints
  - troubleshooting
---

## Quick Fix: Excessive DNS-Prefetch Links

## TL;DR

**Problem**: Archive page has 400+ dns-prefetch links (most from blogroll/reader)  
**Solution**: Changed to page-specific hints  
**Result**: Archive now has ~5-10 relevant hints instead of 400+  

## What Changed

✅ Each page now gets hints only for external domains **on that specific page**  
✅ No more site-wide hint lists injected everywhere  
✅ Blogroll links no longer pollute other pages  

## Do You Need to Change Anything?

**No!** This works automatically. But you can optionally exclude domains:

### Optional: Exclude Reader Domains

Add to your `markata.toml`:

```toml
[resource_hints]
exclude_domains = [
    "reader.waylonwalker.com",
]
```

This will prevent reader links from generating hints (useful if you don't want any hints for blogroll domains).

## Will This Slow Down Builds?

**No significant impact**:
- Adds ~100-200ms to build time (for 100-500 pages)
- You'll barely notice it
- Much better user experience makes it worth it

## What If I Want the Old Behavior?

You can disable auto-detection and manually configure hints:

```toml
[resource_hints]
auto_detect = false

# Manually add critical domains
[[resource_hints.domains]]
domain = "cdn.jsdelivr.net"
hint_types = ["dns-prefetch"]
```

## Next Steps

1. **Test it**: Build your site and check an archive page
2. **Optionally**: Add `exclude_domains` if you want to filter reader links
3. **Done**: Enjoy cleaner HTML and better performance!

## Verification

Check hint count on a page:

```bash
# Count dns-prefetch hints on archive page
grep -c 'rel="dns-prefetch"' markout/archive/index.html

# View actual hints
grep 'dns-prefetch' markout/archive/index.html
```

**Expected results**:
- Archive: ~5-10 hints (was 400+)
- Blog posts: ~4-8 hints (was 400+)
- Blogroll: Still high, but legitimate (those links ARE on that page)

## See Also

- [[resource-hints|Resource Hints Guide]] - Full documentation
- [[performance|Performance Guide]] - Overall optimization tips
- `examples/resource_hints.toml` - Configuration examples
