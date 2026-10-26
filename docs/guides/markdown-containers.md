---
title: "Markdown Containers"
description: "Use ::: fenced containers safely, including predictable nested containers"
date: 2026-09-27
published: true
slug: /docs/guides/markdown-containers/
tags:
  - documentation
  - markdown
  - containers
---

# Markdown Containers

markata-go supports fenced Markdown containers for grouping content in a `<div>` with classes, an optional ID, and custom attributes.

## Basic container

Use three or more colons followed by a class name or attributes to open a container. Close it with a bare line containing only the same number of colons.

````markdown
::: card {#intro .featured data-kind="note"}
This Markdown is inside the card.

:::
````

The closing line must be **colon-only**. A line such as `::: card`, `::: {.card}`, or a decorated/indented marker should not be used as a closer.

## Nested containers

Use a longer fence for each nested level and close each level with its matching bare fence. This makes the structure predictable and keeps content after an inner close inside its parent container.

````markdown
::: cards
Outer content before the card.

:::: card
Inner card content.

::::

Outer content after the card.

:::
````

The example above renders the `card` inside `cards`, then continues rendering the final paragraph inside the outer `cards` container.

A useful rule is:

| Purpose | Example |
|---|---|
| Open outer container | `::: cards` |
| Open nested container | `:::: card` |
| Close nested container | `::::` |
| Close outer container | `:::` |

Do not rely on a shorter outer closing fence to implicitly close deeper nested containers. Close the innermost container first with the same number of colons it used to open, then close the parent.

## Classes and attributes

The opening line can include multiple classes and an attribute block:

````markdown
:::: panel wide {#example .accent data-state="ready"}
Content

::::
````

This produces a container with the classes `panel`, `wide`, and `accent`, the ID `example`, and the custom `data-state` attribute.

## Web Awesome example

The same nesting rule is useful for components that transform Markdown containers after rendering:

````markdown
::: wa-tabs

:::: wa-tab {label="macOS"}
```bash
brew install markata-go
```

::::

:::: wa-tab {label="Linux"}
```bash
curl -fsSL https://example.com/install.sh | sh
```

::::


:::
````

The outer `wa-tabs` container uses three colons and each nested `wa-tab` uses four, so the structure is unambiguous.

## Roads

The default theme turns a `road` container into a scroll-drawn path of linked
milestones:

````markdown
::: road timeline

- [Kickoff](/notes/kickoff/) Jan 2026
- [Launch](/notes/launch/) Jun 2026

:::
````

See [Roads](/docs/guides/roads/) for all styles and options.

## Responsive article layouts (static HTML and CSS)

Markata also provides **opt-in, article-scoped** layout classes built from the
same fenced containers. They generate plain HTML at build time: no JavaScript,
external assets, or custom component framework. These classes intentionally use
`md-` prefixes to avoid changing existing `::: card`, `::: columns`,
`!!! vsplit`, and feed cards.

| Container | Result |
|---|---|
| `md-layout` | Bordered grouping for an architecture diagram |
| `md-box` | Simple bounded box |
| `md-card` | Soft-surface card (scrollable if it contains a wide table) |
| `md-columns` | Auto-fitting equal-width columns, stacking when narrow |
| `md-grid` | Responsive grid, with optional fixed column count |
| `md-stack` | Vertical grouping |
| `md-divider` | Subtle rule |
| `md-connector` | Decorative arrow between rows |

Use `{data-gap="sm"}` or `{data-gap="lg"}` on grid/stack/layout containers,
`{data-columns="2"}` (also 3 or 4) for explicit desktop columns, and
`{data-align="center"}` on a card or box. Explicit columns stack at narrow
viewports. Unknown values retain the defaults; arbitrary CSS is not accepted.
Use descriptive text to convey relationships, rather than relying on arrow
direction alone for accessibility.

Heading permalink links stay inline inside `md-card`, including nested headings,
so ordinary titles do not add a horizontal scrollbar. The links remain available
on hover and keyboard focus. Wide code blocks still scroll independently; a wide
comparison table can scroll within its card without widening the page.

### Architecture diagram

````markdown
::: md-layout

:::: md-box {data-align="center"}
### Falcon3 — Kubernetes
Two deployments, shared published site volume.


::::

:::: md-columns {data-gap="sm"}

::::: md-card {data-align="center"}
### Website
Nginx serves current release.


:::::

::::: md-card {data-align="center"}
### Builder Admin
Builds and publishes releases.


:::::



::::

:::: md-connector {data-direction="both"}


::::

:::: md-box {data-align="center"}
### Shared `/data/site` Longhorn PVC
Replicated storage: NVMe + HDD.


::::



:::
````

The connector arrows are only a decorative aid, not a replacement for the
descriptive labels. For connected arbitrary graphs, use Mermaid.

### Icons and feature rows

The built-in `icons` plugin already accepts Lucide shortcodes like
`:lucide-rocket:`. Combine those shortcodes with `md-icon` (centered icon,
optional `{data-size="lg"}`), or with an `md-feature` row containing a
`md-feature-icon` and `md-feature-body`. No new icon font, CDN, or
JavaScript is necessary.

````markdown
::: md-feature

:::: md-feature-icon
:lucide-rocket:

::::

:::: md-feature-body
### Verify production
Smoke-test the new Markdown components at mobile widths.

::::

:::
````

The icon remains decorative, so provide readable heading text. The row stays
within its article column at narrow widths. For side-by-side icon cards or
more icon names, see [Icon Shortcodes](/docs/guides/icons/).

### Side-by-side tables

````markdown
::: md-columns

:::: md-card
### Current storage

| Node | Capacity |
| --- | ---: |
| Falcon2 | 3 TB |
| Falcon3 | 4 TB |


::::

:::: md-card
### Proposed storage

| Node | Capacity |
| --- | ---: |
| Falcon2 | 3 TB |
| Falcon3 | 14 TB |


::::



:::
````

Each table remains normal semantic Markdown HTML and can independently use
table sorting when enabled. Cards keep their own horizontal overflow instead
of forcing the whole article to scroll. For more examples, see
`examples/gallery-site/pages/layout-primitives.md`.

### Status and metric reports

You can also compose native, accessible status badges, metric cards, progress bars,
and before/after comparisons. See the [Markdown metrics guide](/docs/guides/markdown-metrics/)
for the required values, labels, and data-tone options.
