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

For a multi-stage process with decisions, start with **one primary story**
(e.g. a `road steps` journey) and reserve `md-columns` / `md-card`
for an important choice. Add `md-feature` only where a step benefits from a
visually emphasized explanation. **Do not reproduce every road step a
second time as an icon row.** Add Mermaid only when the precise network
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

**Generic exemplar:** `examples/editorial-visual-composition.md` shows
a single Markdown document with a strong unboxed title, a short lead,
one road, one feature row, one meaningful card pair, plain editorial
paragraphs, and one status notice. Treat it as a *pacing reference*, not
a mandatory component recipe. It deliberately avoids invented
measurements, project-specific jargon, and custom frontend frameworks.

## Compose an editorially paced visual post

**Visual quality bar for a multi-section visual request:**

- Begin with a clear, *unboxed* title and short lead. Establish what matters
  before presenting a card, meter, or process diagram.
- Use one dominant visual grammar for the main story. Let other content be
  ordinary paragraphs, a compact list, or whitespace-separated sections.
- Reserve a visually prominent box for a significant comparison, choice,
  warning, or quantified fact. Do not box ordinary filler.
- Give each section one idea, with meaningful H2 headings and short
  introductions. Use H3 for subordinated card or step labels, not as the
  only headline size across the page.
- Avoid duplicate representation: a road of steps *and* one identical
  feature row for every step adds density without adding information.
- A diagram or table may supplement the composition; neither is a
  substitute for a designed article when the user asked for that.
- **Density check:** if cards and feature rows run continuously from top
  to bottom, remove or merge some. A few strong moments beat a grid of
  equally weighted components.
- **Hierarchy check:** close the page and recall the three most important
  things. If no distinction stands out, simplify the layout and wording.

1. **Lead with the answer.** Explain what is done, what is risky, and what should happen next before decorating.
2. **Give a calm overview.** A sentence or two often beats a metric-card row. Use two to four cards only when there are meaningful distinct values or categories.
3. **Show relationships.** Prefer a progress list, before/after comparison, timeline, table, or diagram that answers a concrete question.
4. **Surface unresolved work.** Name the owner, blocker, verification state, next action, and evidence location. Never color "pending" green.
5. **End with action.** Finish a progress report with a prioritized next-step list and a handoff; include links to PRs, builds, or screenshots.
6. **Make evidence explicit.** Separate *implemented*, *locally tested*, *CI passed*, *merged*, *deployed*, and *verified live*. Merging does not prove deployment.

## Typography, whitespace, and hierarchy

Think in **reading beats**, not a component checklist.

1. **Opening beat:** A strong H1, a short deck/intro paragraph, and space
   before the first layout. Avoid an H2 nested immediately inside a
   bordered `md-layout` wrapper just because the container exists.
2. **Primary visual beat:** One useful road, chart, choice pair, or scorecard.
   Use `{data-gap="lg"}` on supported `md-layout`, `md-columns`,
   `md-grid`, or `md-stack` when their *children* need breathing room.
3. **Quiet explanatory beat:** Simple H2 + paragraph, perhaps a short
   list, without decoration. Let the eye rest between visual clusters.
4. **Secondary beat:** An optional compact visual that answers a
   *different* question, not the same process repeated.
5. **Closing beat:** Plain-language consequence, uncertainty, or next
   action, with a status component only when it meaningfully signals risk.

**Typographic hierarchy:** Markata's visual containers now preserve
distinct H2–H6 type sizes and weights. Still place real article sections
under unboxed H2 headings, keeping brief H3/H4 labels inside choice
cards. Use the zero-download `system` fontpack by default:
`system-ui` body/headings and a native monospace stack for code.
Expressive brush/poster fonts are opt-in for pages or branded sites,
not the baseline for documentation, reports, and visual workflows.

**Practical type ladder (default theme):**

| Role | Semantic markup | Size / weight intent | Color role |
| --- | --- | --- | --- |
| Page title | H1, outside a container | Largest, bold (about 2.25rem) | Primary text |
| Section | H2, outside cards | Prominent (about 1.875rem), bold | Primary text |
| Choice / component title | H3 in `md-card` or `md-feature` | Moderate (about 1.25rem), strong | Primary text |
| Paragraph | `p` | Comfortable reading size, normal weight | Primary text |
| Labels and supporting context | `small`, captions, short metadata | Smaller (roughly 0.75–0.875rem), medium | Muted text |
| Code or technical identifiers | Inline/fenced `code` | Native monospace, distinct from prose | Theme semantic colors |

Use the tokens already provided by the theme; these are roles, not
instructions to hardcode `px`, custom fonts, or arbitrary hex colors into
content. Use no more than one saturated accent in a small visual grouping.
Weight 700 should be exceptional; normal paragraphs need room to stay light.
A body paragraph, label, status pill and card heading should not all
share the same size or boldness.

**Visual editing pass:**

- Vary long and short forms: single-line lead, medium paragraph,
  compact choice, quieter explanatory paragraph.
- Use short card titles and one concise paragraph before details.
  If a card needs multiple dense paragraphs and several commands,
  move implementation detail into its own prose section.
- Avoid placing a `road steps`, a full row of icon features, and a
  second group of step cards back to back.
- Use whitespace between independent sections rather than a
  `md-divider` after every paragraph. Source blank lines ensure valid
  Markdown but do not guarantee rendered breathing room: preview.
- Prefer visible language and hierarchy to decorative status badges,
  repetitive Lucide icons, or color filling every container.
- View at narrow mobile widths and on desktop. Inspect the first screen
  and a middle section for wrapping, vertical rhythm, and excessive
  uniformity; revise after screenshots where possible.

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

### Example: a spacious visual story

Read `examples/editorial-visual-composition.md` for a complete,
repository-neutral, one-file Markdown example. Notice what is **not**
decorated: the title, intro, H2 section headings, and two explanatory
sections are ordinary text. The main journey uses one road; a single
two-card comparison highlights the meaningful fork, and only one icon
feature row emphasizes a central fact. The final status notice is
reserved for an actual caution.

**Design exercise:** compare that restrained page with one that repeats
every road milestone as an icon feature, every detail as a card, and
every conclusion as a colored box. Prefer the former when the user asks
for breathing room or better typography.

**Mermaid is optional**, not the default hero for a multi-section visual
post. Use it for exact connections or complicated loops that a road
and short prose cannot convey. Check runtime support first:
`architecture-beta` is unavailable in Mermaid 10.9.5. Validate rendered
output, not just fenced source.

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
**revise the artifact** before declaring completion—unless a narrow diagram
or table was explicitly requested. Equally, if every paragraph is inside a
card or icon row, **remove components** to restore hierarchy.
For a one-file authoring task outside a Markata site, explain how it could
be previewed with Markata and avoid claiming an actual visual inspection
that was not performed. If a Markata binary is available and a preview is
feasible, use a **disposable scratch site outside the user's repository**
to build and inspect the authored Markdown, then discard scratch output;
do not leave `output/`, `.markata/`, or new site configuration in a
repository where the user requested exactly one file. If a real browser
preview is unavailable, distinguish source validation from visual acceptance.
Do not replace rich container syntax with generic Markdown merely because
the repository has no Markata config.

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
