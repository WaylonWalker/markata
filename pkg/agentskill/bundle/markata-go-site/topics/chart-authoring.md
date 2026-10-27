# Chart Authoring and Data Visualization

Use this topic for **chart**, **charts**, **graph**, **plot**, **data visualization**,
**visualize this data**, **bar chart**, **line chart**, **area chart**,
**stacked chart**, **trend**, **distribution**, **scatter plot**,
**time series**, **dashboard**, **analytics report**, or a visual post
with meaningful numerical comparisons.

Markata supports Chart.js 4 via fenced `chartjs` JSON blocks and a separate
`contribution-graph` renderer for calendar heatmaps. This is **published
Markdown**, not a proprietary chart widget schema. **Do not** output an
HTML-like chat-chart component or assume a screenshot substitutes for the
chart's source data.

A chart is a first-class visual option alongside `md-stat`, `md-card`,
`road`, and prose. A broad "visualize" request does not automatically mean
"draw Mermaid" or "write a table." When numerical data contains an actual
comparison, trend or distribution, choose a chart if it will communicate
that relationship more clearly.

## Choose the chart from the question

| Reader's question | Prefer | Notes |
| --- | --- | --- |
| Which categories are largest or smallest? | Vertical or horizontal `bar` | Use horizontal bars for long names; start a quantitative bar axis at zero |
| How has a value changed across ordered periods? | `line` | Keep chronological labels ordered; indicate units and actual reporting periods |
| How much of the whole comes from each category? | `doughnut` or `pie` | Use few mutually exclusive categories and a meaningful known whole; bars are easier for precise ranking |
| How much did each part contribute over time? | Stacked `bar` | Use a comparable denominator; do not mix counts and percentages |
| Is overall volume changing while fluctuating? | Filled `line` (area-like) | `fill: true` with a `line` chart; only when baseline/area is meaningful |
| Does one numeric variable relate to another? | `scatter` | Use `{x,y}` values, label both axes; note the fallback is source JSON, not a tabular scatter view |
| How did publishing vary across calendar days? | `contribution-graph` | Heatmap is more useful than hundreds of thin daily bars |
| One bounded task-completion value? | `md-progress` or `md-meter` | Do not force Chart.js onto a single known scalar |
| Sequence, architecture, dependency or decision? | `road` or Mermaid | Chart.js is not a graph-topology/flowchart replacement |
| A few exact values with no visual comparison? | Semantic table | Do not make a chart just to decorate the page |

Chart.js also supports `radar`, `polarArea`, and `bubble`. Use those
when their encoding is genuinely the clearest answer. Prefer simpler
bars or lines if readers need accurate comparisons.

## Valid authoring syntax

Write **complete, parseable JSON** inside a fenced `chartjs` block:
`type`, `data.labels`, and `data.datasets`. Do not use the chat UI's
`Chart`/series format, JavaScript functions, comments, trailing commas,
or Markdown tables as the chart source. Chart.js's `options` is an
optional JSON object.

For a built-in figure card, add top-level plain-text `title` and `caption`
strings. The title does not enter the TOC; no page-local JS/CSS is needed.
Optional top-level `"data_display": "collapsed"`
starts the accessible data disclosure collapsed, overriding the site
`[markata-go.chartjs] data_display` setting (`"open"` by default).
All three metadata keys are stripped before Chart.js initialization.
Short horizontal bar figures (up to four labels, `options.indexAxis: "y"`)
automatically use a compact canvas; other charts retain normal sizing.

These examples use **illustrative values** and are not user measurements.

### A clear categorical bar chart

~~~~markdown
## Work by category

```chartjs
{
  "type": "bar",
  "data": {
    "labels": ["Docs", "Tests", "Features", "Fixes"],
    "datasets": [
      {"label": "Completed items", "data": [8, 12, 6, 10]}
    ]
  },
  "options": {
    "responsive": true,
    "plugins": {
      "title": {"display": false},
      "legend": {"display": false}
    },
    "scales": {
      "y": {"beginAtZero": true, "title": {"display": true, "text": "Items"}}
    }
  }
}
```

Tests account for the largest count in this *illustrative* dataset.
~~~~

Use a short prose conclusion tied to the numbers, not a restatement of
every bar. For long labels set `options.indexAxis` to `"y"`, and put
`beginAtZero` on the corresponding quantitative scale.

### A trend line

~~~~markdown
## Trend over time

```chartjs
{
  "type": "line",
  "data": {
    "labels": ["Jan", "Feb", "Mar", "Apr", "May"],
    "datasets": [
      {
        "label": "Published articles",
        "data": [4, 6, 3, 8, 7],
        "tension": 0.2,
        "fill": false
      }
    ]
  },
  "options": {
    "responsive": true,
    "plugins": {"legend": {"position": "bottom"}},
    "scales": {
      "y": {"beginAtZero": true, "title": {"display": true, "text": "Articles"}}
    }
  }
}
```

This illustrates a trend with a temporary dip; it is not evidence of
any real publishing history.
~~~~

For an area-style visual use a `line` dataset with `"fill": true`.
When datasets have different units, do not silently plot them as though
they share a scale. Use separate charts or explicitly configured labeled
axes. Avoid falsely implying smooth or continuous measurement from
sparse categorical periods.

### Part-to-whole composition

~~~~markdown
## Share by category

```chartjs
{
  "type": "doughnut",
  "data": {
    "labels": ["Guides", "Notes", "Tutorials"],
    "datasets": [
      {"label": "Articles", "data": [12, 9, 6]}
    ]
  },
  "options": {
    "responsive": true,
    "plugins": {"legend": {"position": "bottom"}}
  }
}
```
~~~~

Prefer bars for many groups, tiny slices, long legend names, or comparisons
where small differences matter. Avoid arbitrary 3D, unlabeled slices, or
pie charts for values that do not form a whole.

### A stacked comparison

~~~~markdown
```chartjs
{
  "type": "bar",
  "data": {
    "labels": ["Q1", "Q2", "Q3"],
    "datasets": [
      {"label": "New", "data": [6, 8, 7]},
      {"label": "Revised", "data": [4, 3, 5]}
    ]
  },
  "options": {
    "responsive": true,
    "plugins": {"legend": {"position": "bottom"}},
    "scales": {
      "x": {"stacked": true},
      "y": {"stacked": true, "beginAtZero": true}
    }
  }
}
```
~~~~

## Compose charts into a designed article

A strong data story typically has:

1. **One-sentence lead:** the specific decision or question the data answers.
2. **A small, grounded overview:** one or two `md-stat` values if useful,
   not a grid of redundant totals or invented changes.
3. **Primary chart:** the appropriate bar, line or other chart, given
   enough width and a clear title or adjacent H2.
4. **Interpretation:** a short paragraph explaining the supported pattern,
   with the measurement period, units and any uncertainty.
5. **Optional secondary visual:** only if it answers a different question
   (e.g. trend vs composition), with breathing room between sections.
6. **Exact values:** Chart.js's build-time data disclosure, or a table when
   the numbers themselves are the primary answer.

Do **not** put a detailed Chart.js graph in a tiny two-column card that
crushes axis labels. Let the main chart occupy the article's readable
content width. Use cards beside or above it for a few high-value metrics.
A chart is not a fancy replacement for every table or a diagram.

For a full, portable, single-file example see
`examples/chart-composition.md`. It uses generic invented values,
visually distinct questions, and no project-specific jargon.

## Accessible, static-first, and theme-aware

The built-in Chart.js renderer creates a **semantic data table at build
time** for common charts where `data.labels` aligns with dataset arrays.
The `Chart data` details remain available with JavaScript disabled, when
the chart library fails, and for users who want exact values. The data
disclosure only collapses after successful canvas initialization.

- Always supply a visible chart heading/caption that says *what* and
  *when*, with units. Axes, legend labels, and non-color encoding matter.
- Do not encode the entire meaning only by hue. Avoid indistinguishable
  series, excessive legend entries, or ambiguous periods.
- For standard table fallback, keep each dataset's `data` array length
  equal to `labels`. The renderer uses a configuration disclosure
  instead if data is irregular, absent, or exceeds 200 labels/12 datasets.
  Scatter/bubble arrays with `{x,y}` or `{x,y,r}` entries should be
  treated as advanced: expect the config disclosure, not tabular data.
- Charts use a canvas; unlike pure HTML cards, chart strokes, colors and
  tick text are not automatically guaranteed to follow every palette.
  Existing Chart.js defaults or authored colors may be hard to see in
  dark mode. Inspect the actual chart in light, dark, and custom palettes
  before claiming theme compatibility. Do not assume a `var(--token)`
  string is resolved into a canvas color by Chart.js.
- Reuse an established site chart color treatment when available. If new
  runtime theme adaptation is needed, treat that as a **separate renderer
  change with browser tests**—do not invent an untested chart theme schema
  in content.
- Keep charts responsive and avoid axes overflowing at 320/360/412px.
  Inspect at 768 and 1440px too. Verify no-JS table visibility, keyboard
  access to disclosures, reduced-motion expectations and print.
- Distinguish build-time values from current/live telemetry. Never claim
  chart data was fetched or verified unless it was; label samples,
  missing values, period boundaries, and estimates clearly.
- A `chartjs` block needs the plugin enabled and its vendored Chart.js
  asset available; if it appears as code or blank canvas, inspect the
  active site config, output HTML and browser console before rewriting.

Markata's published `chartjs` support is **not** the same as a native
interactive conversation chart component. It can render Chart.js
visualizations and accessible fallbacks, but it doesn't automatically
select charts, query a database, theme a canvas across all palettes,
or attach complex interactive filtering. Author the correct chart first;
request explicit renderer work for missing runtime behaviors.

Simple horizontal bar charts with one dataset and up to five labels use compact, rounded bars by default. Set dataset bar geometry or `options.datasets.bar` / `options.elements.bar` to choose your own; set `options.aspectRatio` or `options.maintainAspectRatio` to control sizing. Multi-dataset and larger charts keep Chart.js defaults.
