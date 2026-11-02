---
title: "Image Optimization Plugin Design"
description: "Add an optional imageoptimization plugin that generates AVIF/WebP variants for local images at build time and rewrites HTML to use <picture> with a stable..."
date: 2026-02-15
published: true
tags:
  - documentation
---


## Goal

Add an optional `image_optimization` plugin that generates AVIF/WebP variants for local images at build time and rewrites HTML to use `<picture>` with a stable `<img>` fallback. External images are unchanged.

## Constraints

- Use external CLI encoders (`avifenc`, `cwebp`) when available.
- Do not fail the build if an encoder is missing; warn and skip the format.
- Do not add responsive sizes in this issue.
- Write optimized variants next to the original output image.

## User Configuration

```toml
[markata-go]
hooks = ["default", "image_optimization"]

[markata-go.image_optimization]
enabled = true
formats = ["avif", "webp"]
quality = 80
avif_quality = 80
webp_quality = 80
cache_dir = ".markata/image-cache"
avifenc_path = ""
cwebp_path = ""
```

## Behavior

1. Render stage scans `post.ArticleHTML` for `<img>` tags.
2. For local `src` values, rewrite each to a `<picture>` wrapper:
   - `<source type="image/avif" srcset="...">` when AVIF is enabled
   - `<source type="image/webp" srcset="...">` when WebP is enabled
   - Keep the original `<img>` as the fallback
3. External URLs, protocol-relative URLs, and data URIs are skipped.

## Outputs

- `image.jpg` -> `image.avif` and `image.webp` in the same output folder.
- Original image remains unchanged and serves as the fallback.

## Cache Strategy

Cache metadata in `.markata/image-cache/` keyed by:

- Source path
- File size + mod time
- Encoder + quality settings

If the key matches, skip re-encoding.

## Error Handling

- Missing encoder: warn once per build; skip that format.
- Encode failure: warn for that file; keep original `<img>`.
- Missing source file: warn and skip.

## Tests

- HTML rewrite: local images get `<picture>`, external images stay unchanged.
- Cache hit skips encode; cache miss re-encodes.
- Missing encoder does not fail build.
- Failed encode does not corrupt output HTML.

## Docs and Spec

- Add `image_optimization` to `spec/spec/OPTIONAL_PLUGINS.md`.
- Document config and usage in `docs/guides/configuration.md` and `docs/reference/plugins.md`.
- Add a new user guide in `docs/guides/image-optimization.md`.

## CI/Builder Image

- Add `avifenc` and `cwebp` to the builder image in PR #781.
