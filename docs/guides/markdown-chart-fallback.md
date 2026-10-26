# Accessible Chart.js fallbacks

`chartjs` fenced blocks render a standard HTML **data table at build time**.
They retain a JavaScript-free representation even if the Chart.js CDN is
unavailable, the reader disables JavaScript, or interactive rendering fails.

The theme adds a Chart.js canvas alongside an expanded **Chart data** disclosure.
Once Chart.js successfully initializes, the disclosure becomes collapsed but
remains available for exact values or screen-reader users. If the library fails
to load, the table remains open.

## Example

````chartjs
{
  "type": "bar",
  "data": {
    "labels": ["Full", "Incremental"],
    "datasets": [{"label": "Build seconds", "data": [46, 21]}]
  }
}
````

The renderer emits semantic table headers and cells without a second chart
library or a network request. For unusual configurations such as scatter
charts with object-valued `x`/`y` entries, the original JSON remains
available in a **Chart configuration** disclosure instead of fabricating
tabular labels. Huge tables (more than 200 labels or 12 datasets) also fall
back to the original configuration to bound generated HTML.

Chart data is treated as untrusted Markdown input: HTML labels are escaped
and chart JSON is safely encoded before inclusion in an inline script.

## Behavior

- No-JavaScript: canvas is hidden; table or configuration is visible.
- JavaScript working: canvas becomes visible; data disclosure is collapsed,
  but can always be reopened.
- Chart.js loading or rendering error: canvas remains hidden; fallback remains
  visible.
- Multiple charts per post: each chart has its own fallback.
- Pages without `chartjs` blocks are unchanged.

This is an accessibility and progressive-enhancement improvement to the
existing Chart.js plugin. It is **not** a new chart data format, SVG exporter,
or automatic visual data verification.
