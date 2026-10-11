# Slow 3G filmstrip benchmark (experimental)

Tracking: [#1553](https://github.com/WaylonWalker/markata-go/issues/1553).

This is a reproducible **browser capture and comparison tool**, not evidence of a
site performance improvement. It addresses the Chromium/Pyppeteer problems
encountered during Phase 3B/3C investigation. Offline frame-capture and
**synthetic gzip HTTP A/B smoke tests** now pass in GitHub Actions. The HTTP
test serves real CSS, JS and image responses and rejects a deliberate 200/HTML
response for CSS. It has **not** completed a real Markata site A/B.

## Requirements

- Python >=3.11 and `uv`, or another environment with Playwright installed.
- A compatible Chromium executable (`--chromium /path/to/chromium`), or
  Playwright's managed browser.
- A stable HTTP(S) test server hosting complete, independently reproducible
  before and after sites, with matching CSS, JS, image, font, and media assets.

The script uses a PEP 723 dependency declaration; no repository-wide Python
environment is required.

## Offline CDP smoke test

From the repository root:

```sh
uv run --script scripts/slow3g_compare.py \
  --before "fixture:$(pwd)/scripts/fixtures/slow3g-readable-main.html" \
  --after "fixture:$(pwd)/scripts/fixtures/slow3g-readable-main.html" \
  --runs 2 --observe-ms 3500 --checkpoints 1000 2000 3000 \
  --output ./slow3g-smoke
```

Open `slow3g-smoke/filmstrip.html`. The fixture mode uses
`page.set_content()`: it validates screencast frames and asynchronous
acknowledgments, **not network throttling**, document arrival, or browser
resource fetches. FCP/LCP may be null for the fixture.

## Automated gzip HTTP smoke

The [Slow 3G filmstrip smoke workflow](../.github/workflows/slow3g-harness.yml)
runs `scripts/slow3g_fixture_server.py` with before/after HTML, actual
CSS/JS/SVG endpoints and `Content-Encoding: gzip`. It verifies complete CDP
filmstrip frames, fetched asset types, response MIME types and compression.

It also serves deliberately broken CSS (HTML with status 200) and confirms
that the benchmark rejects the run. The workflow uploads filmstrip artifacts
for debugging. This is a **test of the benchmark**, not evidence of a faster
production site.

## Serve two actual generated Markata outputs

The synthetic fixture is **not** a substitute for testing the generated site.
The repository now includes `scripts/slow3g_site_server.py`, a local-only
HTTP server that serves the full contents of one generated output directory
with matching gzip-compressed HTML/CSS/JS/JSON, correct MIME types for images
and fonts, and streamed byte-range responses for videos. Missing assets
return 404 instead of accidentally returning index.html.

Build **the same site content at the same revision** twice using the original
Markata executable and the proposed patch executable. Save their *complete*
generated output folders separately; the filenames below are placeholders,
not prescribed Markata output paths:

```sh
# Terminal 1: unchanged build (replace with your actual generated output)
python scripts/slow3g_site_server.py --root /tmp/markata-before --port 8000

# Terminal 2: changed build (replace with your actual generated output)
python scripts/slow3g_site_server.py --root /tmp/markata-after --port 8001

# Confirm compression, MIME and missing-asset behavior before measuring:
curl -sSI -H 'Accept-Encoding: gzip' http://127.0.0.1:8000/archive/
curl -sSI http://127.0.0.1:8000/css/main.css
curl -sSI http://127.0.0.1:8001/css/main.css

# Terminal 3: run alternating paired measurements
uv run --script scripts/slow3g_compare.py \
  --before http://127.0.0.1:8000/archive/ \
  --after http://127.0.0.1:8001/archive/ \
  --runs 5 --observe-ms 55000 --timeout-ms 70000 \
  --checkpoints 1000 2000 5000 8000 12000 18000 25000 35000 45000 52000 \
  --output ./slow3g-archive
```

Run the same experiment on `/`, a representative article, `/shots/`,
and `/reader/` *where both builds provide the same content*. Open
`slow3g-archive/filmstrip.html` and inspect **main text readability**
manually; do not infer first-readable time from FCP alone.

**External resources are not mirrored by this server**: fonts/images hosted
on separate domains, including Dropper and any CDN assets, still use those
services and can vary between runs. Document which URLs differ, record
failures and late third-party resources, and do not claim an isolated CSS
speedup without making asset inventories comparable. This server uses gzip
rather than production Cloudflare/Brotli; results are controlled lab
comparisons, **not production latencies**. Do not use this test server publicly.

The server has independent regression tests for `GET`, `HEAD`, gzip,
MIME, range requests, missing assets, and symlink/path traversal safety.

## Real A/B validation

Build both variants from the same pinned content revision, serve both with a
real compression-aware HTTP server, and confirm that their resource URLs
resolve to the same images/fonts/CSS/JS. Do **not** use the old temporary HTTP/2
server that returned the HTML document for every asset URL.

```sh
uv run --script scripts/slow3g_compare.py \
  --before 'http://127.0.0.1:8000/archive/' \
  --after 'http://127.0.0.1:8001/archive/' \
  --runs 5 --observe-ms 50000 \
  --output ./slow3g-archive
```

Use `--ignore-https-errors` only for self-signed local test certificates. If
local Chromium isn't available, install the Playwright-managed browser. Pin
the browser and Playwright version when running a published comparison.

Defaults: 390×844 viewport, 400 ms latency, 51,200 bytes/s download, 4× CPU.
The script uses alternating A/B and B/A pairs, fresh browser contexts,
disabled browser cache, and one owned browser process.

Results:

- `results.json`: request accounting, network failures, FCP/LCP/CLS,
  checkpoint/actual frame-receipt timestamps, and per-run diagnostics.
- `summary.json` and `summary.csv`: valid-run counts and FCP/LCP median/p90.
- `filmstrip.html`: visual frame comparisons.

**Manually mark first visibly readable main content using the images.** FCP
may reflect only the nav shell. DOM geometry does not prove that pixels were
painted. `encoded_bytes_completed` counts completed transfers, not bytes
received at FCP or readability.

The harness marks a run invalid if a CSS/JS/image/font response has an
inappropriate MIME type or HTTP status, even if an incorrect server responds
with 200. The `document_delivery` metadata records protocol and compression
when visible through CDP.

Treat major asset failures, HTTP error responses, missing frames, layout
shifts, or inconsistent page content as invalid comparisons. Repeat warm-cache
tests separately if needed. The script's `--no-js` option is a separate
visual smoke check; PerformanceObserver metrics are unavailable in that mode.

For the sidebar-order experiment, review `/archive/`, `/shots/`,
`/reader/`, and a text-heavy article; validate skip links, tab order,
drawer behavior, dark/light themes, CLS and no-JS before opening a performance
PR. **Do not merge or deploy a source-order change based on decoded HTML byte
offset alone.**

The standalone tool does not automatically label readable main-text pixels;
review the filmstrips. The generated-site test server does implement gzip
but does not simulate production Cloudflare edge caching.
