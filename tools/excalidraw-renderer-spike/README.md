# Excalidraw renderer comparison spike

This directory contains a reproducible fixture corpus and temporary runners
for comparing browserless SVG exporters with Markata's current Chromium
renderer. The spike is research tooling. None of its optional Node or Rust
dependencies are part of Markata's normal build.

## Fixture corpus

Generate the scenes with:

```sh
node tools/excalidraw-renderer-spike/generate-fixtures.mjs
```

The generated fixtures cover rough shapes and connectors, bound text and five
font families, rotation, frames and grouping, freehand strokes, inline and
external images, large scenes with unusual coordinates, malformed JSON,
newer scene versions, and unsupported element types.

## Run each candidate

Install the exact Node dependencies before running these commands:

```sh
cd tools/excalidraw-renderer-spike/node
npm ci
cd ../../..
node tools/excalidraw-renderer-spike/node/benchmark.cjs

cd tools/excalidraw-renderer-spike/node-official
npm ci
cd ../../..
node tools/excalidraw-renderer-spike/node-official/benchmark.mjs
```

The first Node runner tests `@aldinokemal2104/excalidraw-to-svg`; the second
calls the official `@excalidraw/utils` export API directly with a small JSDOM
environment adapter. Both run six exports per fixture, report process peak
RSS and determinism, and write SVGs to `/tmp/markata-excalidraw-spike/` by
default. Override the output and fixture with `EXCALIDRAW_SPIKE_OUTPUT` and
`EXCALIDRAW_SPIKE_FIXTURE`. Run one fixture in a new process to measure its
first call independently:

```sh
EXCALIDRAW_SPIKE_FIXTURE=shapes-rough.excalidraw \
  node tools/excalidraw-renderer-spike/node-official/benchmark.mjs
```

Run the pinned Rust renderer with:

```sh
cargo run --release --manifest-path tools/excalidraw-renderer-spike/rust/Cargo.toml -- \
  pkg/plugins/testdata/excalidraw/renderer-corpus \
  /tmp/markata-excalidraw-spike/rust
```

Run the Chromium reference manually. Its build tag keeps Chromium out of
ordinary Go tests:

```sh
go test -tags excalidraw_spike ./pkg/plugins \
  -run '^TestExcalidrawRendererSpike$' -count=1 -v
```

Set `EXCALIDRAW_SPIKE_FIXTURE` to isolate a scene and
`EXCALIDRAW_SPIKE_OUTPUT` to choose the directory for reference SVGs. The
reference test uses the browser already installed on the machine; it does not
download one.

## Measurement notes

- The JavaScript runners use six exports per fixture. `first_ms` is the first
  export call in that process; `warm_median_ms` is the median of the following
  five. Use a fresh process for each fixture when comparing first calls.
- The official Node adapter supplies an approximate text-width function
  because this spike does not install a native canvas package. Its text layout
  measurements are therefore not a production adapter.
- Peak RSS is process RSS for Node and Rust. The Chromium test samples the Go
  test process and its descendants because most memory is in Chromium.
- Outputs are also passed to `visualizations.NewSVGAsset`. This preserves the
  shared SVG trust boundary; `<image>` output is expected to fail until a
  separate bounded inline-raster policy is implemented.
- Renderer version, operating system, measurements, and known limitations are
  recorded in
  [`docs/reports/excalidraw-renderer-comparison.md`](../../docs/reports/excalidraw-renderer-comparison.md).
