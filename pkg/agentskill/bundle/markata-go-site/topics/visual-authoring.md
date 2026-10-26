# Visual Authoring, Diagrams, and Dashboards

Use this topic when the user says **visual**, **visualize**, **visualization**, **visual post**, **visual article**, **visual report**, **diagram**, **dashboard**, **infographic**, **scorecard**, **status board**, **progress tracker**, **roadmap**, **timeline**, **flowchart**, **architecture map**, **comparison**, or **at a glance**. Also consider it for "show progress", "make this scannable", "draw the flow", or "show what remains." These are *intent keywords*, not names of a specific chat interface or visual product.

The goal is expressive, responsive **published Markdown** using Markata's existing
building blocks. Do not assume the user wants a screenshot, generated image,
app, or proprietary UI syntax. Select the output medium from the request. Use
plain prose or a normal table when graphics would add no value.

## Decide whether to compose a page or draw one diagram

**A visual document is an information design, not a synonym for Mermaid.**
For requests such as "make this visual", "visualize this workflow", "visually
represent", "visual post", "dashboard", or "make it look like a rich report",
the default deliverable is a **composed document**: hierarchy, a useful opening
summary, responsive sections, theme-aware cards or feature rows, and
appropriate visual relationships. A Mermaid diagram followed by a table
and prose is **not sufficient** when the user asked for visual *design*.
Only choose a diagram-only or table-only answer when the user explicitly
asks for that narrow format or when richer layout would obscure the facts.

For a multi-stage process with decisions, start with **a `road steps` journey
plus `md-columns` / `md-card` choices**; enrich important steps with
`md-feature` and Lucide icons. Add Mermaid *only* when the precise network
of edges or complex loops genuinely needs a graph. If the user asks for a
single standalone Markdown file in a repo that is **not** a Markata site,
still author the rich Markata Markdown in **that one file**, without creating
an unnecessary site, config, CSS, or JS. Explain that Markata is required to
render the special containers, and don't claim a browser preview was tested.

Before writing, choose the primary visual grammar (for example a journey with
choice cards, a compact scorecard, or a comparison gallery). Lead with it.
Keep the number of components proportional to the information; don't fabricate
metrics or add cards merely to hit a quota.

## Choose the visualization that matches the information

| User intent | First choice in a Markata post |
| --- | --- |
| Status, release readiness, project health | `md-status`, `md-stat`, `md-progress`, a short evidence table |
| Numeric measurement with natural range | Native `md-meter` with explicit value and unit |
| Two or more options, before/after | `md-columns` with `md-card`, and semantic tables |
| Ranked tradeoffs or comparison | Markdown table with visible labels; sortable if supported |
| A time sequence, journey, change history | `::: road timeline`, `road steps`, or `road editorial` |
| Workflow with steps and keep/archive decisions | `road steps` plus `md-columns`/`md-card` branches and optional `md-feature` steps; Mermaid only if edges are essential |
| Simple infrastructure topology | `md-layout`, `md-box`, `md-columns`, `md-connector` |
| Arbitrary dependency/network graph | Mermaid; check the deployed Mermaid version first |
| Data trends or category shares | Fenced `chartjs` JSON with build-time accessible table fallback |
| Publishing/activity calendar | Fenced `contribution-graph` block |
| Short feature explanation with icons | `md-feature`, `md-feature-icon`, `md-feature-body`, pinned `:lucide-...:` icons |
| Teaching a composition | Rendered example plus `md-source` native details |

Read `topics/rich-markdown.md`, `topics/roads.md`, and
`topics/analytics-storytelling.md` as needed. For exact syntax check
`docs/guides/markdown-containers.md`, `markdown-metrics.md`, and
`markdown-chart-fallback.md`. Inspect site overrides before adding new CSS.

**Copyable exemplar:** `examples/visual-workflow.md` is a one-file
repo-steward visual workflow with an opening layout, icon cards, a road,
decision branches, and feature rows. Use it as a composition reference,
not as generic boilerplate. It deliberately has no Mermaid, invented
counts, or custom component framework.

## Compose an information-dense visual post

**Visual quality bar for a multi-section visual request:**

- An opening visual anchor tells readers what the document is about: e.g.
  `md-layout`, a small `md-columns` summary, or a well-labeled feature row.
- The main story uses a fitting layout, such as `road steps`, `md-feature`
  rows, a grid of comparison cards, or meaningful measurements.
- Secondary branches and cautions are visually grouped, but semantics and
  accessible text survive with JavaScript disabled.
- A diagram or data table may supplement the composition; neither is a
  substitute for it by default.
- **Design check:** remove the Mermaid and any ordinary table. Would the
  remaining document still communicate a clear visual structure? If not,
  redesign the layout before declaring the visual request complete.

1. **Lead with the answer.** Explain what is done, what is risky, and what should happen next before decorating.
2. **Give a small overview.** Use two to four metric or status cards only when the labels and values are meaningful.
3. **Show relationships.** Prefer a progress list, before/after comparison, timeline, table, or diagram that answers a concrete question.
4. **Surface unresolved work.** Name the owner, blocker, verification state, next action, and evidence location. Never color "pending" green.
5. **End with action.** Finish a progress report with a prioritized next-step list and a handoff; include links to PRs, builds, or screenshots.
6. **Make evidence explicit.** Separate *implemented*, *locally tested*, *CI passed*, *merged*, *deployed*, and *verified live*. Merging does not prove deployment.

### Example: visually readable release status

The numbers below are illustrative author-supplied data, **not live Markata measurements**.

~~~~markdown
## Release progress

::: md-layout

:::: md-columns

::::: md-stat {data-tone="success"}
Checks passing

**8 / 10**

2 checks are pending

:::::

::::: md-stat {data-tone="caution"}
Production

**Not verified**

The deployment revision is unknown

:::::

::::

:::: md-progress {data-value="8" data-max="10" data-label="Eight of ten release checks passed"}

8 of 10 release checks passed.

::::

:::: md-status {data-tone="caution"}

Awaiting production smoke test.

::::

:::

| Gate | State | Evidence |
| --- | --- | --- |
| Local render | Passed | Screenshot |
| CI | Passed | Build link |
| Deploy revision | Unknown | Verify rollout |
| Live page | Pending | Browser audit |
~~~~

Choose `md-progress` for a bounded task process; choose `md-meter` for a real measurement. If no denominator or reliable number exists, use descriptive statuses instead of an invented percent. All gauges require numeric `data-value` and `data-max` and a descriptive `data-label`. Keep values and labels readable in ordinary Markdown.

### Example: compose a visual workflow with a branch

A step-by-step journey plus two decision cards communicates most product
workflows better than one large flowchart. This is an excerpt; the complete
one-file example is `examples/visual-workflow.md`.

~~~~markdown
## Choose a route

::: road steps

- [Discover](#discover) 01
- [Decide](#decide) 02
- [Prepare](#prepare) 03
- [Publish](#publish) 04

:::

::: md-columns {data-gap="lg"}

:::: md-card

::::: md-icon {data-size="lg"}
:lucide-git-branch:
:::::

### Keep

Prepare the repository, review the cleanup plan, then commit and publish.

::::

:::: md-card

::::: md-icon {data-size="lg"}
:lucide-archive:
:::::

### Archive

Verify the archive destination and copied data before removing local files.

::::

:::
~~~~

The full document should also include a strong opening layout, step-specific
`md-feature` rows, and clear warnings about destructive operations where
appropriate. Keep the output responsive and readable without JS.

**Mermaid is an optional detail view**, not the default hero for a
multi-section visual post. Choose a fenced Mermaid flowchart when exact
connections, loops, dependencies, or exception routes cannot be conveyed
clearly with cards and a road. Check the deployed Mermaid runtime before
using version-specific syntax: `architecture-beta` is not supported by
Mermaid 10.9.5. Test the rendered browser output, not just the source.

### Example: numerical chart with text fallback

~~~~markdown
```chartjs
{
  "type": "bar",
  "data": {
    "labels": ["Build", "Review", "Deploy"],
    "datasets": [
      {"label": "Completed checks", "data": [5, 3, 0]}
    ]
  }
}
```
~~~~

Chart.js fenced blocks publish a semantic data-table fallback at build time. When JavaScript is unavailable or Chart.js fails, the reader still has the source values. Do not turn real-time dashboards into static claims without disclosing freshness.

## Review the generated document before delivering

For an explicitly **visual** assignment, inspect the actual output, not just
whether the Markdown is syntactically legal. Did it use suitable built-in
Markata containers? Can a reader skim the primary flow without decoding
Mermaid? Are secondary choices grouped into meaningful cards rather than a
dense table? Does the first screen establish hierarchy and purpose? Do
unknown numbers remain unknown? Did the user request just one file? Does the
file work as an authorable Markata document even outside a Markata site?

If the only "visual" elements are a Mermaid block and a Markdown table,
**revise the artifact** before declaring completion. For a one-file authoring
task outside a Markata site, make a transparent note about how to preview it,
but do not replace rich container syntax with generic Markdown just because
the current repository has no Markata config.

## Design and accessibility

- Use a clear title, restrained card count, readable typography, and generous but economical spacing. Prefer one strong visual over decorative charts everywhere.
- Prefer theme tokens and built-in `md-*` layouts; do not hardcode light-only colors, install an icon font, or invent a heavy JavaScript UI.
- Nested `:::` containers use progressively longer opening and matching closing fences, with bare colon-only closing lines and blank lines before closers.
- Use semantic headings, HTML tables, labeled gauges, text equivalents for diagrams, visible status words, and focusable controls. Color and iconography alone cannot communicate state.
- Mobile layouts must stack or scroll intentionally at 320, 360, 412, 768, and 1440 pixels; test dark/light/custom palettes.
- Check with JavaScript disabled, reduced motion, keyboard access, print where relevant, and code-copy controls that do not obscure code. Record screenshots after final build.
- Any screenshot or visual regression baseline must identify viewport, theme, revision, environment, and date.

## Limits and verification

Markata renders rich, static-first articles; it does **not** automatically verify project progress, read CI, or produce proprietary native widgets from Markdown. Gather evidence before presenting values. Distinguish a rendered diagram from a working interactive interface, and a screenshot from executable test coverage.

Run a fast local build while editing, then the normal production-equivalent build before release. Test both HTML/DOM semantics and the final browser appearance. For source-sensitive visuals, keep an example fixture and a regression check. Confirm the deployed revision, purge/cache behavior, and live page before marking release complete.

To install updates from a newly built Markata binary, run `markata agent update` and `markata agent doctor`; editing the bundled skill source alone does not update local agent installations.
