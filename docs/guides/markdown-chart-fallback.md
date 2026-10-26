---
title: "Accessible Chart.js figures and fallbacks"
description: "Add titles, captions, and JavaScript-free data to charts"
date: 2026-10-26
published: true
---

# Accessible Chart.js figures and fallbacks

`chartjs` fenced blocks render a standard HTML **data table at build time**.
They retain a JavaScript-free representation even if the Chart.js CDN is
unavailable, the reader disables JavaScript, or interactive rendering fails.

The theme adds a Chart.js canvas alongside an expanded **Chart data** disclosure.
For charts using the options below, successful Chart.js initialization collapses
the disclosure, which remains available for exact values or screen-reader users.
If the library fails to load, the table keeps its configured state.

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

## Titles, captions, and compact figures

Add plain-text `title` and/or `caption` keys at the top level of the chart JSON:

````markdown
```chartjs
{
  "title": "Build time",
  "caption": "Seconds per build; lower is better.",
  "data_display": "collapsed",
  "type": "bar",
  "data": {
    "labels": ["Full", "Incremental"],
    "datasets": [{"label": "Seconds", "data": [46, 21]}]
  },
  "options": {"indexAxis": "y"}
}
```
````

The built-in theme renders a responsive surface card with a left-aligned,
non-italic figcaption. Its title is **not a heading**, so it does not enter the
TOC. The canvas is labelled by the title and references the figcaption.
Short horizontal bar figures (up to four labels) use a compact canvas, capped at
130px; other figures retain normal Chart.js sizing. You do not need page CSS or
JavaScript. Neither title nor caption interprets HTML or Markdown.

These authoring keys are stripped before Chart.js receives the JSON. We use JSON
metadata rather than fence attributes because the Markdown renderer retains
the language but not arbitrary fence-info options (the legacy CSV HTML matcher
is not a reliable authoring path for this). Chart.js's own chart title remains
available separately under `options.plugins.title`.

## Initial data disclosure state

```toml
[markata-go.chartjs]
data_display = "open" # default; or "collapsed"
```

Set top-level `"data_display": "open"` or `"data_display": "collapsed"` in an
individual chart to override the site default, with or without a figure.
Collapsed details still have a visible, keyboard-operable summary without JS.
Successful rendering of an opted-in chart collapses the disclosure regardless
of its initial state; a loading/rendering error leaves that state intact.
Invalid values produce a configuration error (site setting) or a visible chart
error (per chart).
With no new options, existing chart markup and the default-open fallback remain
unchanged.

Legacy charts without these options also retain their original initialization
script: it attempts to collapse the wrapper div rather than the disclosure,
so their data stays open. Add a title, caption, or per-chart `data_display` (or
set the site's `data_display` to `"collapsed"`) to opt into corrected behavior.

## Behavior

- No-JavaScript: canvas is hidden; table or configuration is reachable in details
  (expanded by default).
- JavaScript working: canvas becomes visible; opted-in charts collapse the data
  disclosure, which can always be reopened.
- Chart.js loading or rendering error: canvas remains hidden; fallback remains
  reachable in its configured initial state.
- Multiple charts per post: each chart has its own fallback.
- Pages without `chartjs` blocks are unchanged.

This is an accessibility and progressive-enhancement improvement to the
existing Chart.js plugin. It is **not** a new chart data format, SVG exporter,
or automatic visual data verification.
