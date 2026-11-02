---
title: "Static SVG visualization foundation"
description: "Developing backend-rendered SVG adapters with deterministic assets, dimensions, caching, and explicit limits."
date: 2026-10-31
published: true
tags: [documentation, rendering, performance, development]
---


Markata is developing a shared backend rendering pipeline so diagrams and charts
can arrive as **finished SVG assets**, with no visualization engine running in the
reader's browser.

> **Availability:** This is a **developer-only foundation** in
> `pkg/visualizations`. There is **no supported `d2` or `vega-lite` Markdown
> fence**, site configuration section, or new `markata render` capability yet.
> The existing `mermaid` and `chartjs` syntax and output are unchanged.

## What the foundation provides

A backend renderer produces SVG bytes from a request. The common package then
checks bounded SVG output, derives intrinsic dimensions, assigns a SHA-256-based
filename, publishes the asset and emits normal HTML for a `figure`:

```text
source + palette + font + renderer/version + options
       |
       v
backend Renderer.RenderSVG(ctx, request)
       |
       v
Pipeline (in-process cache and concurrent coalescing)
       |
       v
NewSVGAsset -> Publish(outputDir, asset)
       |
       v
<figure><img src="/assets/visualizations/<sha256>.svg"
             width="960" height="540" alt="..."></figure>
```

The filename depends on the exact SVG bytes; renderer cache identity also
includes the renderer version, theme content hash, and effective options.
Repeated assets can be reused by the browser or a CDN.

## Integrate a backend in Go

An adapter implements this interface (shown here for reference, not as a
working D2 integration):

```go
type Renderer interface {
    ID() string
    Version() string
    RenderSVG(context.Context, visualizations.Request) ([]byte, error)
}
```

Use `visualizations.Pipeline.Compile(ctx, backend, request)` for a rendered
`Asset`. `visualizations.Publish(outputDir, asset)` writes a final static
asset; `visualizations.FigureHTML(asset, alt, caption, aboveFold)` returns
escaped markup with an explicit image width and height.

The request must include all layout-affecting inputs, especially theme and
font contents. Renderers must obey context cancellation and return actionable
errors for unsupported syntax, oversized SVG output, or unavailable layout
dependencies.

```bash
go test ./pkg/visualizations/...
go test -race ./pkg/visualizations/...
go vet ./pkg/visualizations/...
```

## Why the constraints matter

- **No layout shift from unknown image dimensions:** the figure includes
  intrinsic `width` and `height`. Responsive site CSS still has to preserve
  the intended image aspect ratio.
- **No visualization JavaScript on the published page:** the output is a
  normal external SVG image, not a browser chart/canvas library.
- **Smaller repeat views:** a fingerprinted path enables stable cache
  identities, but Brotli/gzip and appropriate response headers remain
  hosting responsibilities.
- **Bounded processing:** one SVG has a 4 MiB uncompressed size limit, and
  the in-process cache holds up to 128 assets or 32 MiB.
- **Clear trust boundary:** the XML checks reject many active/remote
  constructs but are not a complete SVG sanitizer, and do **not** guarantee
  a fully offline page. Author-supplied images and linked resources still
  need a separate output audit.

The current restrictions also reject some legitimate SVG features such
as hyperlinks, embedded images, animation, and `foreignObject`. Test
representative outputs before using a new engine. Large embedded fonts and
complex SVG paths can still increase network and browser-paint cost.

## What comes next

D2 is the intended first backend adapter, followed by Vega-Lite through
`vl-convert`, once each meets compatibility and performance checks.
Site-builder integration will need base-path-aware asset links and a
version-compatible addition to the existing render asset manifest.

Mermaid and Chart.js remain the supported authoring tools until those
follow-ups land. See the [Markdown guide](/docs/guides/markdown/),
[accessible Chart.js figures](/docs/guides/markdown-chart-fallback/),
and [renderer interoperability contract](/docs/guides/python-renderer/).

For exact technical invariants, see
[`spec/spec/STATIC_VISUALIZATIONS.md`](https://github.com/WaylonWalker/markata/blob/main/spec/spec/STATIC_VISUALIZATIONS.md)
(after this PR merges).
