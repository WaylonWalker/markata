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
