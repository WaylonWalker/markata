# Build-time page weight diagnostics

This feature measures **uncompressed generated HTML source bytes**, not network
transfer, Brotli/gzip savings, render-blocking time, Core Web Vitals, or
Lighthouse scores. It is a deterministic, fast signal for agents and authors.

- An ordinary `markata-go build` scans emitted HTML after a successful build.
  `--dry-run` does not scan or write a report.
- The report is written to `<output_dir>/.markata/page-weight.json`, with
  schema version, measurement name, budgets, summary and an array of pages.
  Pages are sorted by descending HTML bytes, then stable relative path.
- Each page includes raw HTML bytes, inline JS/CSS bytes, embedded JSON bytes,
  separately identified copy/sidebar JSON, graph attribute bytes, declared
  stylesheet/script counts, image/iframe counts, and actionable findings.
- Stable codes: `HTML_OVER_BUDGET`, `EMBEDDED_JSON_OVER_BUDGET`,
  `INLINE_JS_OVER_BUDGET`, `SIDEBAR_JSON_OVER_BUDGET`.
  Finding fields include bytes, budget bytes, and advice.
- Default budgets are 160 KiB for post HTML, 320 KiB for other HTML,
  32 KiB embedded JSON, 24 KiB inline JS, and 16 KiB sidebar JSON.
- Regular builds do **not** fail when a page exceeds its budget. Console output
  caps example warnings but the JSON report includes every detected page.
  `--page-weight=false` disables scanning; `--page-weight-html-kb` overrides
  the post-HTML warning threshold.
- Missing/unreadable output or report paths surface actionable build errors when
  diagnostics are enabled.
- The analysis does not request external resources and must not treat HTML
  markup counts as actual network requests.
- All page data is generated locally; no remote telemetry or uploads occur.
