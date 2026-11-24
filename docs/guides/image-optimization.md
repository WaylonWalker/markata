---
title: "Image Optimization"
description: "Generate AVIF and WebP variants for local images"
date: 2026-02-15
published: true
tags:
  - documentation
  - images
  - performance
---


The `image_optimization` plugin generates AVIF and WebP variants for local images and rewrites HTML to use a `<picture>` element with the original image as the fallback.

## Quick Start

Enable the plugin in your hooks:

```toml
[markata-go]
hooks = ["default", "image_optimization"]

[markata-go.image_optimization]
enabled = true
```

## Requirements

- `avifenc` for AVIF output
- `cwebp` for WebP output

If a tool is missing, the build continues and that format is skipped.

## Example Output

Input HTML:

```html
<img src="/images/cat.jpg" alt="Cat">
```

Output HTML:

```html
<picture>
  <source type="image/avif" srcset="/images/cat-480w.avif 480w, /images/cat-960w.avif 960w" sizes="100vw">
  <source type="image/webp" srcset="/images/cat-480w.webp 480w, /images/cat-960w.webp 960w" sizes="100vw">
  <img src="/images/cat.jpg" alt="Cat">
</picture>
```

The generated `.avif` and `.webp` files are written next to the original image in the output folder.

The optimizer reuses the image cache in `.markata/image-cache` and validates source bytes with SHA-256, so a file timestamp change alone does not trigger re-encoding. It hashes each unique source once per build. Within one build, repeated references to the same local output image share one optimization pass that generates the union of their requested widths. Each image occurrence still keeps its own `<picture>` markup.

Each build logs target occurrences and unique target candidates after deduplication. A candidate counts once in `unique_targets` even if its path cannot be resolved or its source file is missing; those candidates are also counted in `missing_sources`. `duplicate_targets` is the number of occurrences removed by deduplication. The summary also reports variants checked, cache hits and misses, and encoder calls and failures.

## Configuration

```toml
[markata-go.image_optimization]
enabled = true
formats = ["avif", "webp"]
quality = 80
avif_quality = 80
webp_quality = 80
widths = [480, 960]
sizes = "100vw"
cache_dir = ".markata/image-cache"
avifenc_path = ""
cwebp_path = ""
```

### Notes

- Only local images are processed (relative paths and site-root paths).
- External URLs, protocol-relative URLs, and data URIs are skipped.
- Images already inside `<picture>` are left unchanged.

For full details, see the [plugin reference](/docs/reference/plugins/#image_optimization).
