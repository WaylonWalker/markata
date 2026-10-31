# Markata native Mermaid research harness

This module is **not** a production renderer and adds no dependencies to Markata's main `go.mod`.
It benchmarks the actual Go APIs of `mmdg` v0.2.0 and `go-mermaid` v0.1.6 in one
persistent process, not their CLI subprocess startup times.

```sh
cd tools/mermaid-compare
go build -mod=mod -o /tmp/mermaid-compare .
cd ../..
python tools/mermaid-compare/run.py --binary /tmp/mermaid-compare --full --iterations 8
python tools/mermaid-compare/gallery.py /tmp/mermaid-comparison
```

Outputs: `comparison.json`, `comparison.md`, `performance.svg`, one SVG per
fixture/engine, `theme-font-{engine}.svg`, and `index.html` for a side-by-side gallery.

The fixture collector extracts every Mermaid fenced block from Markdown under
`docs/`, `examples/`, `rendering-fixtures/` and `benchmarks/` and adds
13 representative diagram families. It deduplicates exact source blocks.

**Interpretation:** The measured numbers are warm **source-to-SVG engine costs**,
not the Markdown rendering API, Python pool, or visitor's browser. Do not
interpret the chart's native timings as the total Markata request improvement.
The previous pool benchmark (230.33ms one-shot, 60.33ms pooled for the Mermaid
example) uses a different timing boundary, and is labeled separately.
The two renderers generate different SVG structures; `index.html` and
`comparison.md` expose those differences, without demanding byte equality.

Both candidates expose theme color settings and `font-family`/size. This test
passes Markata-style colors and an alternative font to each, stores the result,
and checks the resulting XML for changes. **Important:** font-family in SVG is
not proof that the engine measured glyphs correctly. In particular, `mmdg`
searches installed system fonts while the `go-mermaid` API supplies a font
name and size. The visual gallery must be checked for clipping and overflow
with the actual Markata font packs on light and dark palettes before selecting
a production backend. Never bundle proprietary font files.

The comparison runs as an added step of the existing renderer-pool performance
workflow on changes to this harness or Mermaid-related files. GitHub
`workflow_dispatch` performs the complete scan with a longer benchmark.
Regular PR CI uses a smaller repetition count to remain fast; there is **no
browser or mmdc** in this fast suite. The full visual baseline against pinned
Mermaid.js 10.9.5 remains a separate, heavier validation step and has not yet
been run. No native library should be made the default before that comparison
and Markata's real palette/font snapshot checks pass.

The harness deliberately records unsupported syntax and renderer failures.
Its minimum gate is successful SVG for both engines on the six core families;
all other example failures are listed, not hidden or misidentified as passing.
Full Mermaid syntax parity has **not** been established.
