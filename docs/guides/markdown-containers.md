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

## Native tabs (no Web Awesome required)

Native tabs use the same fenced-container syntax as `md-card`. Every
panel is rendered to HTML on the server. Without JavaScript, all panels remain
visible with headings; when JavaScript loads, a small keyboard-accessible tab
strip replaces the headings. Printing expands all panels.

````markdown
::: md-tabs {data-label="Installation methods" style="boxed"}

:::: md-tab {label="curl"}
```bash
curl -fsSL https://example.com/install.sh | bash
```
::::

:::: md-tab {label="mise"}
```bash
mise use -g github:WaylonWalker/markata
```
::::

:::
````

The outer `md-tabs` uses three colons, with direct child `md-tab`
containers using four. Tabs may contain ordinary Markdown, code fences, cards,
and nested blocks. Use `style="boxed"` (default), `style="underline"`,
or `style="pills"`; `data-style` is also accepted. Use `data-label`
to give the tab group an accessible name. Tab labels are plain text.
On narrow screens the tab strip scrolls horizontally. Native tabs have no
Web Awesome dependency and the enhancement script only loads on pages using
them. Existing `wa-tabs` continue to work separately.

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

For prose-led histories, use `::: road editorial story` with a heading per
entry. A leading `:lucide-name:` becomes a rail badge; an emphasis-only
paragraph directly after the heading becomes its eyebrow label. Add `reveal`
(`::: road editorial story reveal`) for optional, reduced-motion-safe scroll
entrances. `story` alone opts into larger type and diamond-bullet highlights;
plain editorial roads keep their existing typography.

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
| `md-section` | Section with automatic heading/intro column, sticky on wide screens |

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

### Markdown story pages

Set `content_width: wide` in a post's frontmatter to opt into a wider frame,
larger balanced title, and a larger opening paragraph. Other pages keep their
existing appearance. This option applies to the default `post.html` template
(custom templates must opt in themselves).

Start each `md-section` with an H2. The heading and the paragraph immediately
after it become the section head; everything else becomes the body. At 1100px
and wider the head sticks beside the body while scrolling. Smaller screens
stack them; print uses static headings. Without a leading H2, all content stays
in the body.

````markdown
---
title: "From foundations to agent workflows"
content_width: wide
---

A short opening paragraph introduces the story.

::: md-section
## 2026 · Agent workflows and access

Agents became first-class platform users.

### Identity and access

Longer narrative, lists, tables, and cards go here.
:::
````

A short token (up to 16 characters, no spaces) followed by the exact separator
` · ` becomes a large primary-color kicker. The full heading's ID is unchanged:
the example still links to `#2026-agent-workflows-and-access`. No JavaScript
or separate `story-head` container is needed.

For a timeline overview, add the `rail` modifier to an `md-grid`:

````markdown
::: md-grid rail {data-columns="2"}
:::: md-card
::::: md-icon
:lucide-bot:
:::::

_2026_

### [Agent workflows](#2026-agent-workflows-and-access)

Agents inspect runs and read logs.

::::

:::: md-card
_2025_

### [Platform visibility](#2025-platform-visibility)

Telemetry reaches more teams.

::::

:::
````

At 1000px and wider, a rail joins card tops with ring ticks that fill on
hover or keyboard focus. The default is four columns; `data-columns="2"`, `"3"`
or `"4"` overrides it. Below 1000px cards stack along a vertical rail with ring ticks. The ticks
remain visible and fill on hover or keyboard focus. A card's heading
link stretches over the card; keyboard focus outlines the card. Use one heading
link per card and avoid other interactive controls in these overview cards.
An emphasis-only paragraph supplies the small uppercase monospace kicker.
Rail cards allow visible overflow so the ring ticks remain visible above their tops.
An `md-section` immediately after the rail gets responsive spacing from 4.5rem to 7.5rem.
These defaults need no page-local CSS.

For unresolved items, use `::: md-grid pending` with nested `md-card`
containers. Cards get dashed borders, transparent backgrounds and muted icons.
Both variants use ordinary Markdown, work without JavaScript, and retain theme
colors. See `rendering-fixtures/story-layout.md` for a full example.

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
