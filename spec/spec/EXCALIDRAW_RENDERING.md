# Excalidraw Static Rendering

**Status:** renderer comparison complete; browserless publishing integration
is not implemented. The current static renderer still uses headless Chromium.
**Tracked by:** [#2014](https://github.com/WaylonWalker/markata/issues/2014).
**Related:** [Excalidraw authoring](EXCALIDRAW.md),
[static visualization foundation](STATIC_VISUALIZATIONS.md).

## Rendering contract

The site builder and standalone Markdown renderer MUST NOT start a browser,
download browser assets, or contact a rendering service unless the author
explicitly selects `mode = "chromium"`.

The default `auto` mode MUST use an Excalidraw SVG artifact only when its source
fingerprint matches the scene JSON, renderer identity, font inputs, theme, and
all output options. If no matching artifact exists, it MUST use an available
browserless renderer. If neither is available, it MUST report which renderer
is missing and how to select an explicit alternative. A renderer failure MUST
NOT fall back to Chromium.

The `chromium` mode MUST be an explicit configuration choice. Chromium mode
MUST remain isolated from `auto` and browserless modes so an unavailable or
failed browserless renderer cannot start a browser.

## Static SVG requirements

Rendered SVG MUST be deterministic for a fixed scene and renderer identity.
The SVG MUST NOT reference remote fonts, scripts, images, stylesheets, or
other resources. Fonts MUST be embedded in the SVG or published as local,
content-addressed assets that use Markata's base-path-aware URL handling.

The shared SVG validator MUST continue to reject active content and external
resources. If Excalidraw output needs embedded raster images, use an explicit
Excalidraw validation policy. That policy MUST allow only bounded, inline raster
payloads with an approved MIME type; it MUST reject nested SVG payloads and
external links. The global validator MUST remain strict for every other
renderer. Every SVG size limit MUST remain explicit and tested.

Scene JSON remains authoritative. A cached SVG MUST be discarded when any
render input changes. Unsupported or malformed scenes MUST produce an
actionable error while preserving the editable Markdown source.

## Renderer comparison

Renderer selection MUST consider the shared fixture corpus in
`pkg/plugins/testdata/excalidraw/renderer-corpus/`. The corpus covers rough
shapes and connectors, text and fonts, frames and groups, freehand paths,
embedded raster images, large scenes and unusual coordinates, malformed JSON,
and newer scene versions.

Candidate results MUST record renderer and dependency versions, generated SVG
size, cold and warm render time, determinism, validator result, network
references, and resource use where the platform exposes a reliable measure.
The official Excalidraw exporter is the visual reference. A candidate that
reports approximations or omissions MUST surface those diagnostics rather
than silently dropping content.

The default backend MUST be selected from measured fidelity, offline behavior,
cross-platform support, maintenance, packaging, and runtime cost. No candidate
is considered production-ready from its language or project description alone.

The current comparison and its measurements are recorded in
[`docs/reports/excalidraw-renderer-comparison.md`](../../docs/reports/excalidraw-renderer-comparison.md).
The comparison does not select a universal renderer backend: the Rust
candidate is promising but has a known bound-text rotation mismatch and only
Linux x86-64 release qualification; Node candidates add substantial memory
and dependency cost. The editor's official browser export should be evaluated
as the source of a fingerprinted pre-rendered artifact before adding a build
runtime dependency.
