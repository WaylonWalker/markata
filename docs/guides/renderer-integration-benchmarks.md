---
title: "Renderer integration benchmarks"
description: "Measured rendering, process, HTTP, and response costs for existing Markata fixtures"
date: 2026-10-29
published: true
tags:
  - documentation
  - benchmarks
  - fastapi
---

These measurements use existing repository fixtures with the [FastAPI editor](/docs/guides/fastapi-editor/).
The benchmark includes the exact post that `BenchmarkBuild_Incremental` edits: `benchmarks/site/posts/blog/2024/01/post-001.md`.
The other cases come from `rendering-fixtures/`; they are visual regression fixtures, rather than existing performance benchmarks.

## Fixture coverage

| File | Coverage |
| --- | --- |
| `post-001.md` | Deterministic site benchmark post, including code |
| `home.md` | Small page with minimal assets |
| `markdown-chart-fallback.md` | Chart.js output and accessible data fallbacks |
| `md-card-headings.md` | Cards and expressive headings |
| `roads.md` | Roads, layout styles, and initialization script |
| `story-layout.md` | Longer structured page |
| `test-headings.md` | Heading syntax and typography |

Inputs retain their complete file contents, including frontmatter.
These calls measure the standalone renderer; they do not run the site build pipeline or its frontmatter loading stage.

## Measured timings

All values are milliseconds. The Go column shows the median of three benchmark runs.
The other columns show medians across 15 calls, except the final column.

| Fixture | Input bytes | Go renderer | Python render call | FastAPI render timer, edit | HTTP edit | HTTP edit p95 |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| `post-001.md` | 1,518 | 3.05 | 137.51 | 135.02 | 136.25 | 148.45 |
| `home.md` | 426 | 0.40 | 128.50 | 128.24 | 129.15 | 136.19 |
| `markdown-chart-fallback.md` | 1,550 | 33.04 | 476.93 | 260.01 | 264.50 | 470.84 |
| `md-card-headings.md` | 1,104 | 0.89 | 151.58 | 152.03 | 153.45 | 167.69 |
| `roads.md` | 1,598 | 1.08 | 139.34 | 174.50 | 177.22 | 224.00 |
| `story-layout.md` | 3,900 | 1.46 | 156.73 | 153.57 | 154.97 | 174.24 |
| `test-headings.md` | 1,245 | 0.82 | 168.01 | 151.57 | 152.70 | 170.24 |

For the existing site benchmark post, Go rendering takes about 3 ms, while a Python render call takes about 138 ms.
The Python API starts a Go process for each call and decodes its JSON output.
The gap includes process initialization, serialization, pipe transfer, and Python decoding; these measurements do not isolate each part.

Chart.js has higher rendering and transfer costs because its response includes vendored JavaScript and generated initialization code.
The current fragment response still sends those feature assets with each chart edit.
Keeping shared CSS loaded does not remove all repeated feature-asset bytes.

## Initial document versus later edits

Response sizes below are uncompressed JSON body bytes, rounded to KiB (1,024 bytes).
They exclude HTTP headers and separately requested appearance assets.
Initial requests include a full document; edit requests set `include_document` to `false`.

| Fixture | Initial HTTP median, ms | Edit HTTP median, ms | Initial response, KiB | Edit response, KiB | Fewer response bytes |
| --- | ---: | ---: | ---: | ---: | ---: |
| `post-001.md` | 135.39 | 136.25 | 46.44 | 20.01 | 56.9% |
| `home.md` | 130.10 | 129.15 | 8.13 | 0.86 | 89.4% |
| `markdown-chart-fallback.md` | 432.89 | 264.50 | 460.25 | 226.92 | 50.7% |
| `md-card-headings.md` | 136.74 | 153.45 | 45.55 | 19.57 | 57.0% |
| `roads.md` | 167.05 | 177.22 | 93.22 | 43.40 | 53.4% |
| `story-layout.md` | 152.85 | 154.97 | 25.48 | 9.54 | 62.6% |
| `test-headings.md` | 154.95 | 152.70 | 13.25 | 3.42 | 74.2% |

Omitting the document shell reduces response bytes in every fixture.
It does not establish a latency improvement: each mode runs separately, and host load varies between modes.
Some fragment requests therefore take longer despite their smaller responses.

## Timing boundaries and limits

- **Go renderer:** `RenderStandalone` setup, Markdown processing, asset generation, and inline tag generation. Fixture reads occur before timing starts.
- **Python render call:** synchronous `render(..., inline_assets=True, unsafe=True)`, including process startup, output transfer, and JSON decoding.
- **FastAPI render timer:** asynchronous `render_async(...)`, including process startup and JSON decoding. Document construction and HTTP response encoding occur afterward.
- **HTTP time:** request object construction, a new localhost connection, FastAPI processing, response transfer, and body reading. Request JSON encoding and response JSON decoding are excluded.

Each Python and HTTP mode has two excluded warm-up calls and 15 measured calls per file.
Every render still starts a fresh Go process; these are warm filesystem measurements, rather than a persistent renderer process.
The nearest-rank p95 equals the maximum with only 15 samples. Use more samples for stable tail-latency decisions.

The measurements ran sequentially on a shared workstation with normal background activity.
They are observations from this host, rather than service limits or a controlled before-and-after experiment.
They exclude the editor's 350 ms debounce, browser layout, chart drawing, Mermaid execution, and user-perceived paint time.
No site build cache, encryption key, or site configuration participates in these standalone calls.

## Environment and raw results

- Date: 2026-10-29 UTC.
- Source revision: `f226485ead9f91081e2eb3c848b3290bf1df040a`.
- CPU: AMD Ryzen 5 3600, 6 cores and 12 logical CPUs.
- OS: `Linux-7.2.8-9-omarchy-x86_64-with-glibc2.44`.
- Go: `go1.26.3 linux/amd64`.
- Python: `3.13.13`.
- FastAPI: `0.162.0`; Uvicorn: `0.54.0`.

The [JSON report](https://github.com/WaylonWalker/markata/blob/main/benchmarks/renderer-integration-2026-10-29.json)
contains every measured sample, fixture SHA-256 hashes, and the installed binary and wrapper hashes.
The [Go benchmark output](https://github.com/WaylonWalker/markata/blob/main/benchmarks/renderer-go-2026-10-29.txt)
contains iteration counts, allocation counts, and all three runs.

## Reproduce

Use a virtual environment with the current Markata wheel, FastAPI, and Uvicorn.
Run these commands from the repository root, sequentially:

```bash
python python/benchmark_editor.py
go test ./pkg/plugins -run '^$' -bench '^BenchmarkRenderStandaloneFixtures$' -benchtime=200ms -count=3
```

The Python script starts a temporary localhost server and stops it when the run ends.
It writes `benchmarks/renderer-integration-YYYY-MM-DD.json` using the UTC date.
It does not replace the live editor server.

Before comparing changes, keep the fixture bytes, dependency versions, machine, and renderer options the same.
Compare medians and raw samples as well as response sizes.
Useful next measurements include process initialization, JSON encoding, asset transfer, browser paint, and memory retained after repeated diagram edits.


## Go server workbench

The [Go workbench](/docs/guides/go-workbench/) uses the same seven fixture files and the same HTTP timing boundaries.
All fixture hashes match the Python report above. The server renders directly in process and keeps one previous result in memory.

Each mode uses two excluded warm-up calls and 15 measured calls.
Before every cache-miss call, a separate unmeasured request renders a small invalidator document.
This keeps the measured fixture bytes unchanged while preventing an accidental cache hit.
Every measured response checks its `cached` field. Cached edits use the same source repeatedly without invalidation.
The invalidator request is excluded from the HTTP time and response size.

### HTTP results

Values are milliseconds. Each median and p95 comes from 15 calls.
The Python column comes from the earlier report; the Go columns come from a separate run on the same workstation.

| Fixture | Python HTTP edit median | Go initial miss median | Go edit miss median | Go edit miss p95 | Go cached edit median |
| --- | ---: | ---: | ---: | ---: | ---: |
| `post-001.md` | 136.245 | 2.247 | 2.143 | 3.599 | 0.365 |
| `home.md` | 129.152 | 0.645 | 0.466 | 0.867 | 0.324 |
| `markdown-chart-fallback.md` | 264.495 | 36.032 | 28.571 | 53.724 | 0.295 |
| `md-card-headings.md` | 153.448 | 0.799 | 1.504 | 1.790 | 0.366 |
| `roads.md` | 177.217 | 0.957 | 0.973 | 1.547 | 0.325 |
| `story-layout.md` | 154.968 | 1.368 | 1.326 | 1.827 | 0.300 |
| `test-headings.md` | 152.701 | 0.838 | 0.805 | 1.721 | 0.293 |

For the site benchmark post, an uncached Go edit takes 2.143 ms over loopback HTTP.
A repeated cached edit takes 0.365 ms.
The earlier Python HTTP measurement for that exact file is 136.25 ms.

The Go server avoids a renderer subprocess per request. It also links feature assets instead of embedding their contents in each response.
Both differences contribute to the observed gap. These runs do not isolate language or process-startup costs.
The native Go renderer benchmark uses inline assets, so it can take longer than the Go HTTP path with linked assets.

### Response sizes

Values are uncompressed JSON body KiB. Separately requested stylesheets, libraries, and fonts are excluded.
Initial responses contain a complete document and fragment fields. Edit responses omit the document.

| Fixture | Python edit response | Go initial response | Go edit response | Go initial-to-edit reduction |
| --- | ---: | ---: | ---: | ---: |
| `post-001.md` | 20.01 | 24.62 | 8.82 | 64.2% |
| `home.md` | 0.86 | 8.99 | 1.02 | 88.6% |
| `markdown-chart-fallback.md` | 226.92 | 24.15 | 8.54 | 64.6% |
| `md-card-headings.md` | 19.57 | 12.99 | 2.99 | 77.0% |
| `roads.md` | 43.40 | 19.58 | 6.30 | 67.8% |
| `story-layout.md` | 9.54 | 30.59 | 11.82 | 61.4% |
| `test-headings.md` | 3.42 | 16.24 | 4.65 | 71.4% |

The Go editor retains loaded libraries and shared styles in the preview frame.
The linked-asset policy makes its chart response smaller, but the browser must fetch each dependency when it is first needed.
Response-size reductions do not include those initial asset transfers and do not establish a browser paint-time improvement.

### Go measurements and reproduction

The Go run uses runtime source revision `059bb011f3314c43ca9e0c6e51d7ed735cf68790`.
The merged parent adds benchmark tooling and documentation without changing that runtime.
The host is the same AMD Ryzen 5 3600 workstation, with 12 logical CPUs and Go 1.26.3.
The Go HTTP client uses Python 3.14.8; the earlier Python-reference client uses 3.13.13.
The comparison includes that client-version difference. The JSON report records its timestamp, OS, binary hash, and harness hash.

The Go server timer includes rendering, asset registration, and optional document construction.
It excludes HTTP JSON encoding and network transfer. The Python server timer above excludes document construction.
Use HTTP measurements to compare the complete request paths.

These are sequential measurements on a shared workstation, rather than a controlled before-and-after experiment.
They exclude server startup, appearance generation, cache invalidation calls, separately fetched assets, editor debounce, browser layout, and drawing.
The Go editor debounce is 250 ms. The p95 equals the maximum with 15 samples; larger runs are needed for stable tail estimates.
No site configuration, site build cache, encryption keys, or site lifecycle stages participate.

The [Go HTTP report](https://github.com/WaylonWalker/markata/blob/feat/go-render-workbench/benchmarks/go-workbench-2026-10-29.json)
contains every sample and fixture/binary/harness hashes.
The [native Go output](https://github.com/WaylonWalker/markata/blob/feat/go-render-workbench/benchmarks/go-workbench-native-2026-10-29.txt)
contains three fresh runs of the parent's renderer benchmark on the same fixtures.

Run these commands sequentially from the repository root:

```bash
go build -o /tmp/markata-workbench-benchmark ./cmd/markata-workbench
python python/benchmark_go_editor.py --binary /tmp/markata-workbench-benchmark
go test ./pkg/plugins -run '^$' -bench '^BenchmarkRenderStandaloneFixtures$' -benchtime=200ms -count=3
```

The standard-library Python client starts a temporary Go server on an OS-assigned loopback port and stops it after the run.
It does not restart the live editor or change its cache.
Use `--samples`, `--warmups`, and `--output` to control the run. The default output uses the UTC date.
