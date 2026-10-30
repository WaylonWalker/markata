---
title: "Status badges, metrics, gauges, and sparklines"
description: "Author static-first status reports, before/after stats, accessible meters, and trend sparklines with Markdown."
date: 2026-10-20
published: true
slug: /docs/guides/markdown-metrics/
tags:
  - documentation
  - markdown
  - containers
---

# Status, metrics, gauges, and sparklines

Markata's Markdown containers can render compact, theme-aware operational
summaries without using JavaScript. Combine these components with the
[responsive layout primitives](/docs/guides/markdown-containers/) to produce
the same kind of dashboard-like reporting found in modern chat interfaces.

The components are **presentation**, not automatic monitoring. Authors supply
the measurements; Markata does not calculate or verify them.

## Components

| Container | Role |
| --- | --- |
| `md-status` | Compact inline status with an optional semantic tone |
| `md-stat` | First paragraph is label; second is large value; third is optional secondary label |
| `md-delta` | Side-by-side comparison caption and emphasized result |
| `md-meter` | Native HTML `<meter>` for measurements like free disk space |
| `md-progress` | Native HTML `<progress>` for partially completed tasks |
| `md-sparkline` | Static inline SVG trend line from a short series of numbers |

`data-tone="success"`, `data-tone="caution"`, and
`data-tone="danger"` use the site's existing semantic palette colors.
Unknown tones use the ordinary theme accent. Both gauges require finite
numbers `data-value` and `data-max` where max is positive and
0 ≤ value ≤ max. They also **require** a descriptive `data-label`, which
becomes the accessible label on the native HTML element.

When an attribute is missing or invalid, the container still renders its
Markdown body normally and no graphical gauge is added. The text should
always include the exact value so the page is understandable in text-only
readers and feeds.

## Verified result badge

````markdown
::: md-status {data-tone="success"}

NVMe cleanup verified

:::
````

Do not use green color alone to communicate success: include the literal
status wording, as above.

## Before/after metrics and native meters

````markdown
::: md-layout

Falcon3 NVMe — before and after cleanup

:::: md-columns

::::: md-stat {data-tone="caution"}
Before

**83 GiB**

18% free

:::::: md-meter {data-value="18" data-max="100" data-label="Before cleanup: 18 percent free" data-tone="caution"}

::::::

:::::

::::: md-stat {data-tone="success"}
After

**212 GiB**

45% free

:::::: md-meter {data-value="45" data-max="100" data-label="After cleanup: 45 percent free" data-tone="success"}

::::::

:::::

::::

:::: md-delta {data-tone="success"}
Space recovered

**+129 GiB**

::::

:::
````

These remain normal Markdown paragraphs and nested containers. Use longer
colon fences for deeper nesting, and place blank lines before bare closers.

## Trend sparklines

A sparkline shows the shape of a series next to the number it explains.
Give `data-values` as two or more finite numbers, separated by spaces or
commas, and a descriptive `data-label`. The line is drawn at build time as an
inline SVG, so it needs no JavaScript. The Markdown body stays as a visible
caption.

````markdown
::: md-stat {data-tone="success"}
Jobs run

**222,421**

Year 4 ran 3x the jobs of year 3

:::: md-sparkline {data-values="8120 21400 48300 144601" data-label="Jobs per year, years 1 to 4" data-tone="success"}

Years 1–4

::::

:::
````

The line fills the container width and scales between the lowest and highest
values. It does not start at zero, so it shows shape, not size. Put the exact
numbers in the text or a table when the reader needs them. The SVG's
accessible name gives the label, point count, first and last values, low, and
high.

A sparkline accepts up to 500 points. With fewer than two points, a missing
label, or a value that is not a finite number, no line is drawn and the
caption renders as normal.

## A build progress bar

A task-completion indicator is not the same as a measured storage percentage.
Use `md-progress` only when a process is underway.

````markdown
::: md-progress {data-value="7" data-max="10" data-label="Build steps completed"}

7 of 10 build steps finished.

:::
````

## Full report with a normal Markdown table

Use the [example page](/layout-metrics/) for the complete composition,
including the existing semantic GFM table. Future sortable table behavior
will work independently of these static metrics and should never be needed
to see the content.

## Compatibility and accessibility

- Fully rendered at build time; no component JavaScript or network requests.
- The source Markdown body and numeric labels remain readable in plain text.
- Native `meter` and `progress` provide standard accessible semantics.
- Sparklines are `role="img"` SVGs with a generated accessible summary.
- Color tones inherit from the active theme, including alternate palettes.
- Nested grids stack at mobile widths using the layout primitives.
