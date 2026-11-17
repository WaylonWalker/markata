---
title: "Excalidraw SVG renderer comparison"
description: "Measured browserless Excalidraw SVG renderers and the Chromium reference for Markata issue 2014."
date: 2026-11-17
published: true
tags:
  - documentation
  - excalidraw
  - rendering
  - performance
---

# Excalidraw SVG renderer comparison

This report compares the official Excalidraw exporter, two browserless Node
adapters, a pure Rust renderer, and Markata's current Chromium renderer. It is
the evidence for [issue #2014](https://github.com/WaylonWalker/markata/issues/2014).

At the time of this comparison, issue #2014 is open and no related PR is
listed. The local repository already has Excalidraw authoring, static SVG
rendering, and folded source fences. The static renderer still starts
Chromium automatically when a published scene is present.

## Result

No headless backend is ready to become Markata's universal default. The Rust
renderer had the lowest runtime and produced deterministic, self-contained
SVG, but visual inspection found a bound-text rotation mismatch. It is also a
new 0.1 release qualified only for Linux x86-64 GNU. The official Node exporter
tracks Excalidraw's implementation more directly and embeds fonts, but the
current Node setup used 267–312 MiB RSS and a text-metrics shim. The other Node
adapter used 352–495 MiB RSS, including a dependency on an Excalidraw test
package version.

The next useful experiment is to export SVG from the already bundled official
editor when a user applies a scene, then persist it with a fingerprint of the
scene, exporter assets, fonts, theme, and export options. This could let a
normal build reuse a matching artifact without installing Chromium, Node, or
Rust. Stale or missing output must trigger an explicitly configured
browserless renderer or an actionable error. The artifact format and its
interaction with the workbench's 200,000-byte document limit still need a
design spike.

Keep the Rust renderer as the first optional backend candidate after the
rotated bound-text mismatch and platform packaging are resolved. Keep the
official Node exporter as the compatibility reference. Neither should be
selected automatically yet.

## Test setup

Measurements were captured on 2026-11-17 on Linux x86-64 GNU:

| Component | Version |
|---|---:|
| Go | 1.27.0 |
| Node.js | 26.13.0 |
| Chromium | 155.0.8059.84 |
| Rust | 1.92.0 |
| Official Excalidraw SVG bundle in Markata | 0.18.1 |
| `@excalidraw/utils` | 0.1.8 |
| `@aldinokemal2104/excalidraw-to-svg` | 1.1.1 |
| `excalidraw-render` | 0.1.0 |

The corpus in
[`pkg/plugins/testdata/excalidraw/renderer-corpus/`](../../pkg/plugins/testdata/excalidraw/renderer-corpus/)
contains rectangles, ellipses, diamonds, lines, arrows, rough fills, rotation,
bound text, five font-family values, multiline text, freehand strokes, frames,
grouping, inline and external images, a 360-element large scene with unusual
coordinates, malformed JSON, a version 3 scene, and an unknown element type.

Each reported isolated fixture used six renders in a fresh process. “First” is
the first render call; “warm” is the median of the following five. The Node
process wall time includes module startup and all six calls. The Go test time
includes browser launch, all six renders, and shutdown. RSS is per-process for
Node and Rust; Chromium RSS is a sampled sum of the test process and its
children. These are local measurements, not CI thresholds.

The direct official Node adapter uses JSDOM and a simple approximate text
metrics provider. It embeds Excalidraw fonts from the installed package, but
the metrics provider is not faithful enough for production layout. This limits
the value of its text timings and visual comparison. The official browser
renderer is used as the visual reference; its current SVG output also contains
remote font CSS, so the current Markata output fails the shared offline SVG
validator for text scenes.

## Runtime and output measurements

Each cell reports `first / warm median` in milliseconds. Process RSS is shown
for the isolated fixture process.

| Renderer | Shapes: time / SVG / RSS | Text: time / SVG / RSS | Large: time / SVG / RSS |
|---|---|---|---|
| Chromium reference | 630 / 5.3 ms; 9,620 B; 1,159,684 KiB | 638 / 5.6 ms; 2,914 B; 1,164,472 KiB | 827 / 167.4 ms; 1,785,355 B; 1,239,892 KiB |
| Official `@excalidraw/utils` | 116 / 21.9 ms; 11,833 B; 273,016 KiB | 205 / 106.1 ms; 20,727 B; 314,424 KiB | 176 / 120.6 ms; 1,788,211 B; 319,232 KiB |
| `@aldinokemal2104/excalidraw-to-svg` | 556 / 8.3 ms; 13,404 B; 360,396 KiB | 573 / 5.7 ms; 25,288 B; 374,704 KiB | 787 / 220.6 ms; 1,785,355 B; 506,308 KiB |
| Rust `excalidraw-render` | 9 / 7.9 ms; 17,463 B; 14,696 KiB | 54 / 52.7 ms; 32,932 B; 24,328 KiB | 46 / 45.2 ms; 1,785,443 B; 17,756 KiB |

The browserless Rust runner used `default-features = false` with the `cjk`
feature. All four renderers were deterministic across six calls for these
fixtures. For scenes without images, all browserless outputs passed Markata's
`NewSVGAsset` validator. The Chromium shapes and text SVGs failed because of
external font CSS; its large scene passed. The main output-size differences
come from font embedding and renderer-specific SVG structure, so byte size is
not a visual quality score.

The cold process wall time for six calls was 0.84–1.32 seconds for the official
Node adapter, 0.93–2.26 seconds for the `aldino` adapter, and 0.05–0.32 seconds
for the already compiled Rust executable. The first per-scene render times
above expose backend startup: the `aldino` first call includes worker startup;
Chromium's first call includes launching the browser. A Go test invocation and
Node module loading add their own process overhead.

## Visual and compatibility findings

The SVGs were rasterized in Chromium and reviewed side by side for rough shapes
and text. Rust matched the reference's rough strokes, fills, connectors, font
appearance, and dimensions closely. One concrete mismatch remains: the label
inside the rotated rounded rectangle rotates with the rectangle in Rust, while
the official exporter keeps that bound label horizontal. The Rust compatibility
report emitted no diagnostic for this case, so the issue needs a regression
test before adoption. The direct official Node output kept the label horizontal;
the `aldino` output also kept it horizontal but used visibly different text
metrics and geometry. This was visual inspection, not a pixel-diff score.

| Case | Official Node | `aldino` Node | Rust | Markata gate |
|---|---|---|---|---|
| Malformed JSON | Parse error | Parse error | Parse error | Rejected before rendering |
| Scene version 3 | Returned a small SVG without an error | Returned a small SVG without an error | Helpful unsupported-version parse error | Rejected before rendering |
| Unknown element type | Returned a small SVG without a diagnostic | Returned a small SVG without a diagnostic | Rendered a fallback and reported two diagnostics | Rejected before rendering |
| Inline PNG | SVG contains `<image>` | SVG contains `<image>` | SVG contains `<image>` | Current shared validator rejects it |
| External image URL | Kept an external reference | Kept an external reference | Reported `image.forbidden-resource`, output still contains `<image>` | Input validator rejects it; shared validator also rejects `<image>` |

The renderer must not be trusted to validate scene schema or output safety.
Markata should reject unknown versions and unsupported elements before render,
check compatibility diagnostics, and run every output through the shared SVG
validator. Supporting images requires a separate policy that accepts only
bounded inline PNG/JPEG data, rejects external URLs and nested SVG, and does
not weaken validation for other renderers.

## Candidate review

### Official `@excalidraw/utils`

The official MIT package is version 0.1.8. Its current `exportScene` API has an
environment hook for non-browser hosts, including the SVG `document`, text
metrics, font loading, canvas, image, and path creation hooks. This makes a
browserless Node adapter possible without reimplementing Excalidraw semantics.
The local JSDOM harness produced deterministic, self-contained SVG for regular
scenes, but took more memory than Rust and the test text metrics were only an
approximation. It accepts newer or unsupported scenes without reporting the
loss. Keep Markata's own schema and compatibility checks in front of it.

The 0.1.8 package is MIT licensed. The Excalidraw changelog documents the
non-browser environment hooks and the current export API at
[`@excalidraw/excalidraw` changelog](https://github.com/excalidraw/excalidraw/blob/master/packages/excalidraw/CHANGELOG.md).

### `@aldinokemal2104/excalidraw-to-svg`

Version 1.1.1 is MIT licensed. It offers a warm worker, font embedding, timeouts,
and cancellation, and its SVG passed Markata validation for non-image scenes.
Its first render includes worker setup; the large-scene warm median was 220.6
ms. Its npm dependency tree pins `@excalidraw/utils` to
`0.1.3-test32`, which is a maintenance and reproducibility risk. The GitHub
repository had 20 commits and one star when checked; the package's own README
describes its [worker and font-subsetting design](https://github.com/aldinokemal/excalidraw-to-svg).

`npm audit` reported zero advisories for the two locked Node comparison
dependencies on this date. This checks only those npm lockfiles; it is not a
general security review of either project.

### `excalidraw-render` for Rust

Version 0.1.0 is pure Rust, bundles fonts, requires no network or browser, and
returns a compatibility report. It had the lowest measured memory and large
scene latency. The crate is `MIT OR Apache-2.0`, with bundled/transitive MIT,
BSD-3-Clause, and CC-BY-4.0 components. Its project docs qualify release builds
only for Linux x86-64 GNU; browser WASM is experimental. The crate is new and
0.x output/diagnostics may change. It also renders synchronously with no
cancellation; its resource limits bound work but do not impose a wall-clock
deadline. Pin the crate and lockfile, preserve notices for bundled fonts, and
run it in a bounded worker process before exposing it to untrusted bulk input.

The Rust renderer's [compatibility and safety documentation](https://docs.rs/excalidraw-render/latest/excalidraw_render/)
describes its limits. Do not adopt it as a default until the bound-text issue,
cross-platform packaging, and current-scene compatibility are tested.

### `@moona3k/excalidraw-export`

The 0.2.1 MIT package is a small Rough.js-based implementation with a bundled
Virgil font. Its documentation lists shapes, text, arrows, freehand, images,
frames, rotation, and fills, but it is not the official Excalidraw exporter.
One simple scene rendered in 3.46 ms initially and 0.84 ms warm, but generated
an 81,625-byte SVG because it embedded the full font. It has only 14 repository
commits and four stars; its latest repository push was February 2026. It was
not run against the complete corpus. Its low runtime cost does not offset the
fidelity and release-confidence gap for Markata's default backend. See its
[supported-element and architecture notes](https://github.com/moona3k/excalidraw-export).

### `JonRC/excalirender`

This MIT project has a Bun/node-canvas/Rough.js pipeline and publishes a
self-contained Linux x64 binary; its documentation offers Docker usage for
Linux, macOS, and Windows. The repository had 84 commits, 35 stars, and a last
push in February 2026. It was not run against this corpus, and it is a CLI
binary rather than a Go library. Its documentation does not establish a
complete font compatibility matrix. It remains worth a later adapter test if
cross-platform packaging is practical. See the
[project's platform and rendering description](https://github.com/JonRC/excalirender).

### Official Excalidraw CLI

The official MIT CLI is version 1.0.0 and requires Node.js 22 or newer. It uses
an already installed Chrome, Edge, Chromium, or Firefox. It must remain an
explicit compatibility path, not an automatic build dependency. Its
[documentation](https://github.com/excalidraw/cli) says it does not download a
browser.

## QuickJS and Rough.js

The pure-Go `go-quickjs` runtime can execute JavaScript without cgo or a browser,
and Rough.js exposes path generation without needing a DOM. Embedding that
geometry code would still require Markata to implement Excalidraw element
semantics, text metrics and fonts, frames, bindings, images, migrations, and
scene-version compatibility. It would move the hard maintenance work into Go
without providing the official renderer's fidelity. This does not justify a
native Go rewrite. See [`go-quickjs`](https://pkg.go.dev/github.com/go-quickjs/go-quickjs)
and the [Rough.js generator API](https://github.com/rough-stuff/rough/wiki/RoughGenerator).

## Shared SVG and image policy

`visualizations.NewSVGAsset` accepts at most 4 MiB and rejects all `<image>`
elements. This is the correct global safety default. Every tested renderer
emits `<image>` for the inline PNG scene; both Node renderers also preserve an
external image URL when given one. No Excalidraw adapter may bypass this
validator or relax it for other visualization types. If Excalidraw images are
supported, validate a strict image MIME allowlist and decoded byte/pixel limits,
reject `image/svg+xml` and all network references, then use an Excalidraw-only
validation path with focused tests.

## Reproduction

See [`tools/excalidraw-renderer-spike/README.md`](../../tools/excalidraw-renderer-spike/README.md)
for fixture generation, pinned dependency installs, and exact benchmark
commands. The ordinary Go test suite does not run Chromium; the reference test
is behind the `excalidraw_spike` build tag.
