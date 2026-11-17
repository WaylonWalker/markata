# Static visualization foundation specification

**Status:** implemented foundation with an Excalidraw static-SVG adapter.
**Scope:** shared SVG validation/publication and the Excalidraw Markdown adapter.
**Related:** [Renderer contract](RENDERER_CONTRACT.md), [Mermaid prerendering](MERMAID_PRERENDER.md), [renderer presets and offline policy #1854](https://github.com/WaylonWalker/markata/issues/1854).

## Goals and non-goals

Markata should be able to publish backend-rendered SVG as ordinary immutable
static assets, with intrinsic dimensions so browsers reserve layout space.
The authoring engine must not be downloaded to a visitor's browser.

The shared foundation currently powers the `excalidraw` Markdown fence. It
does not render D2 or Vega-Lite, change Mermaid/Chart.js behavior, or add an
asset kind to the standalone renderer's schema-v1 `assets` list. A build cannot
claim `offline-static` solely because this package exists.

## Go interface and render identity

A backend implements `Renderer` with `ID()`, `Version()`, and
`RenderSVG(context.Context, Request) ([]byte, error)`. Versions must
include changes in the renderer/layout engine that alter SVG bytes.

`Request` specifies `Source`, `Palette`, `Font`, `ThemeHash`, and
`Options`. `ThemeHash` identifies the effective theme/font **contents**,
not only their names. The caller must provide all inputs affecting the SVG
and must not mutate the options map during a call.

The per-process cache key is SHA-256 over a domain separator and
length-prefixed ID, version, source, palette, font, theme hash, and options
sorted by key. This avoids delimiter and map-order collisions.

- Successful results may be reused; errors are never cached.
- Identical concurrent requests are coalesced. Waiting callers may cancel
  via their own context without cancelling an already running leader.
- The leader's context controls the actual engine render. Its failure is
  returned to waiters and is retryable on a later call.
- Successful entries are kept in an insertion-order cache bounded to
  **128 assets or 32 MiB**; these are implementation limits, not promises
  of persistent site-build caching. Eviction means a future re-render.
- SVG bytes are copied when inserted or returned, so callers do not mutate
  cached assets.

Persisting render-key mappings between builds is **not** implemented. Final
SVG files can still be deduplicated by their content hash.

## Output validation and trust boundary

`NewSVGAsset` accepts only a nonempty well-formed SVG document up to
**4 MiB** of uncompressed bytes. It rejects external entities/DOCTYPE,
non-XML processing instructions, SVG scripting/active elements, links,
embedded images, animation elements, event-handler attributes, external
`href`, and selected network-capable CSS forms. Only fragment `href`
values are supported.

It accepts an SVG root with no namespace or the normal SVG namespace,
rejecting other root namespaces. An explicit `viewBox` must always be
a finite four-number box with positive size, even if explicit dimensions
are present. Dimensions must be positive finite CSS pixel values (maximum
100,000 px). If dimensions are missing, they derive from the `viewBox`;
if one is present, the viewBox aspect ratio determines the other. Relative
units such as percentages are not supported by this foundation.

These restrictions are conservative and may reject legitimate renderer SVG
containing `foreignObject`, `<a>`, `<image>`, animations, or remote
fonts. Adapters must report incompatibilities clearly; **no silent client-JS
fallback** when static output is requested.

**Security boundary:** this is a defensive validation of **trusted backend
renderer output**, not a general-purpose SVG or CSS sanitizer for untrusted
author input. It does not constitute a complete runtime-network audit or a
guarantee of SVG browser behavior. Future adapters must test real output and
add policy checks before claiming `offline-static`.

## Asset publication and HTML

Final SVG identity is SHA-256 of the **exact SVG bytes**, independent of the
render-key hash. The site-relative path is
`assets/visualizations/<sha256>.svg`. Different sources yielding identical
bytes share the same final asset.

`Publish(outputDir, asset)` validates the source bytes and expected identity,
writes to a temporary file and renames within the destination directory, and
reuses existing identical bytes. A conflicting existing asset is an error.
The output directory must be trusted. Asset garbage collection and
cross-process publication coordination are outside this slice.

`FigureHTML(asset, alt, caption, aboveFold)` and
`FigureHTMLWithBasePath(asset, alt, caption, aboveFold, basePath)` require
nonblank descriptive alternative text, validate asset identity/dimensions,
and produce escaped `figure`/`img` markup. The latter projects the URL below a
site URL path prefix. `width` and `height` reserve layout space;
`loading="lazy"` is the default, with eager loading for critical above-fold
images. It emits no JavaScript. Site CSS should make images responsive.

The default generated URL is rooted at `/`, suitable for a root-hosted site.
The site-builder adapter MUST pass the path component of the configured site
URL to `FigureHTMLWithBasePath` when the site is hosted below a path prefix.

## Integration contract for follow-up PRs

1. Add D2 and then Vega-Lite adapters without duplicating Markdown parsing.
2. Add explicit backend capabilities, failures, and resource limits to
   CLI/Python/HTTP as they are implemented; preserve schema v1.
3. Measure actual engine SVGs (including embedded fonts), compressed bytes,
   warm/cold throughput, memory, CLS, and offline browser requests.
4. Update the bundled authoring skill and chart/diagram guides only when the
   corresponding Markdown syntax works.

## Conformance tests

The package tests cover render-key determinism and invalidation, concurrent
coalescing, error retry, cache size, copied bytes, SVG size/dimensions, active
resource restrictions, idempotent publication, and escaping.

Full real-engine fixtures, cross-interface result parity, browser layout
measurements, and external resource audits remain **outstanding requirements**
before the Excalidraw adapter can claim complete conformance or offline output.
