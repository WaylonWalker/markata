---
title: "Native Mermaid SVG renderer"
description: "Opt-in browser-free Mermaid diagrams with Go, Markata palettes and warm Python rendering"
published: true
tags:
  - mermaid
  - performance
  - diagrams
---

# Native Mermaid SVG renderer (experimental)

Markata can render Mermaid diagrams entirely inside its Go process, without
starting Chromium, launching mmdc or shipping Mermaid.js to the browser.

Native SVG is **opt-in**. Existing `client`, `cli` and `chromium` modes remain
available. Client mode remains the default until syntax, typography and visual
verification are complete.

## Enable in a site

```toml
[markata-go.mermaid]
enabled = true
mode = "go"
use_css_variables = true
lightbox = true

[markata-go.mermaid.go]
font_family = "sans-serif"
font_size = 16
node_spacing = 50
rank_spacing = 50
padding = 16
curved_edges = false
transparent_background = false
```

The `go` backend reads semantic palette colors from the same site palette as
the other pre-rendering backends, when `use_css_variables = true`. The SVG
contains colors resolved at **build time**. Switching dark/light themes in
the browser does **not** recolor an already generated SVG; build separately
per palette or use client-mode Mermaid for live recoloring.

Font family and size are passed to the native renderer. The generated SVG may
request any CSS font family, but this library currently estimates glyph widths
using its built-in metrics; it does not read Markata's bundled WOFF2 font faces
for exact layout. Test decorative font packs for clipped or overlapping labels
before relying on them. `sans-serif` is the conservative baseline.

## Python/FastAPI pool

```python
async with RendererPool(workers=2) as pool:
    response = await pool.render(
        "```mermaid\nflowchart LR\n  A[Load] --> B[Serve]\n```",
        mermaid_mode="go",
    )
    assert "<svg" in response.html
```

The one-shot `render(..., mermaid_mode="go")` and
`render_async(..., mermaid_mode="go")` APIs work too.

Native mode returns SVG immediately and does not require the Mermaid.js asset.
The existing presentation CSS still applies. No extra subprocess is created
per diagram; the warm pool's Go process renders each request.

## Compatibility, performance and CI

Native SVG supports many common Mermaid families, but not every upstream
Mermaid.js diagram type. Unsupported syntax returns an explicit error rather
than silently turning into a wrong diagram or falling back to browser work.

- Native renderer version: pinned `go-mermaid` v0.1.6.
- Basic SVG validity and syntax are checked in Go tests and a wider fixture
  comparison suite.
- Dedicated fast tests check concurrent palette isolation, unsafe input,
  native SVG protocol output, and no bundled Mermaid browser JS.
- The render-pool benchmark measures the *same* Markdown + Mermaid source
  with client and native modes through identical Python APIs.
- The concurrency-1 latency bar chart includes both real full-request medians
  and separate, clearly labeled source-to-SVG microbenchmarks.

The full cross-browser/font-pack visual comparison is not complete, and a
native SVG is not a security sanitizer for authored diagrams.

Refer to [Native Mermaid research comparison](https://github.com/WaylonWalker/markata/pull/1891)
and [Native renderer integration](https://github.com/WaylonWalker/markata/pull/1893).
