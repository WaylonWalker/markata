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

For a multi-stage process, choose a **primary wayfinding element**
(e.g. `road steps`), then preserve other *complementary* visual jobs:
a compact overview that orients the reader, cards that compare decisions,
and selected `md-feature` rows that explain a consequential action.
**Different visual roles can coexist.** Do not reproduce every road
milestone as a nearly identical icon row, but also do not remove an
overview or options gallery merely because a road already exists.
Add Mermaid only when precise graph topology genuinely matters. If the user asks for a
single standalone Markdown file in a repo that is **not** a Markata site,
still author the rich Markata Markdown in **that one file**, without creating
an unnecessary site, config, CSS, or JS. Explain that Markata is required to
render the special containers, and don't claim a browser preview was tested.

Before writing, choose the primary visual grammar (for example a journey with
choice cards, a compact scorecard, or a comparison gallery). Lead with it.
Keep the number of components proportional to the information; don't fabricate
metrics or add cards merely to hit a quota.

## Chart-first visualizations

**Numeric data should get a genuine data visualization when it improves
comprehension.** For explicit **bar chart**, **line chart**, **area chart**,
**stacked chart**, **scatter plot**, **time-series**, **trend**, or
**data dashboard** requests, read `topics/chart-authoring.md` before
writing the page. Choose the encoding from the question and the scale of
the data rather than treating Mermaid, cards, or tables as the only visuals.
Use real `chartjs` fenced JSON for numerical charts, and
`contribution-graph` for publishing calendar heatmaps.

A composed data story can include one or two grounded KPI cards, a
full-width primary chart, and quiet interpretation. The exact-value
HTML table is produced by Markata's Chart.js fallback when ordinary
labels and dataset arrays are aligned. Avoid inventing measurements
or mixing incompatible units on the same quantitative axis.

**Current visual limitation:** Chart.js's canvas and default colors do
not automatically adapt to every site's light/dark/custom palette just
because the surrounding Markdown containers use theme tokens. Check
rendered contrast before claiming a polished chart. A native chart
widget from a chat application is not an authorable Markata Markdown
primitive.

## Choose the visualization that matches the information

| User intent | First choice in a Markata post |
| --- | --- |
| Status, release readiness, project health | `md-status`, `md-stat`, `md-progress`, a short evidence table |
| Numeric measurement with natural range | Native `md-meter` with explicit value and unit |
| Trend behind a headline number | `md-sparkline` (`data-values`, `data-label`; add `area` to shade under the line) inside the `md-stat` |
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

Read `topics/chart-authoring.md` for numerical plots,
`topics/rich-markdown.md` for card layouts,
`topics/roads.md` for flows, and
`topics/analytics-storytelling.md` for site-derived data as needed. For exact syntax check
`docs/guides/markdown-containers.md`, `markdown-metrics.md`, and
`markdown-chart-fallback.md`. Inspect site overrides before adding new CSS.

**Historical storytelling:** For a sourced, multi-period project retrospective, see `examples/evidence-led-history.md`. Begin with a thesis and source scope, use editorial roads for milestones, and keep unresolved work distinct.

**Generic exemplar:** `examples/editorial-visual-composition.md` shows
a single Markdown document with an introductory pair of orientation
cards, a road for the sequence, another pair of cards for decisions,
one icon feature for a key explanation, ordinary prose for details,
and a meaningful status notice. The two card groups answer **different
questions**, rather than duplicating the road. Treat it as a *pacing
reference*, not a mandatory recipe. It avoids invented measurements,
project-specific jargon, and custom frontend frameworks.

## Compose an editorially paced visual post

**Visual quality bar for a multi-section visual request:**

- Begin with a clear, *unboxed* title and short lead. Establish what matters
  before presenting a card, meter, or process diagram.
- Use one main **navigation** or storytelling grammar, but retain
  complementary orientation, comparison, or exception views where they
  make the information easier to understand. Not every visual is a repeat.
- Reserve visually prominent boxes for **distinct purposes**: short
  summary/orientation, an actual choice, a warning, a supporting route, or
  grounded measurement. Do not box filler or turn substantial instructions
  into tiny cards.
- Give each section one idea, with meaningful H2 headings and short
  introductions. Use H3 for subordinated card or step labels, not as the
  only headline size across the page.
- Avoid duplicate representation: a road of steps *and* one identical
  feature row for every step adds density without adding information.
- A diagram or table may supplement the composition; neither is a
  substitute for a designed article when the user asked for that.
- **Density check:** when visual groups run together, shorten the copy,
  separate them with a calm reading beat, improve labels, and group related
  items before deleting anything. A few strong moments beat equal-weight
  decoration, but a wall of numbered prose is not automatically better.
- **Hierarchy check:** close the page and recall the three most important
  things. If no distinction stands out, simplify the layout and wording.

1. **Lead with the answer.** Explain what is done, what is risky, and what should happen next before decorating.
2. **Give a compelling overview.** Use a short lead plus compact cards or a
   feature pair when readers benefit from seeing two distinct responsibilities,
   choices or categories at a glance. Don't invent metrics or force a card row
   where a sentence communicates the same information.
3. **Show relationships.** Prefer a progress list, before/after comparison, timeline, table, or diagram that answers a concrete question.
4. **Surface unresolved work.** Name the owner, blocker, verification state, next action, and evidence location. Never color "pending" green.
5. **End with action.** Finish a progress report with a prioritized next-step list and a handoff; include links to PRs, builds, or screenshots.
6. **Make evidence explicit.** Separate *implemented*, *locally tested*, *CI passed*, *merged*, *deployed*, and *verified live*. Merging does not prove deployment.

## Typography, whitespace, and hierarchy

Think in **reading beats**, not a component checklist.

1. **Opening beat:** A strong H1, a short deck/intro paragraph, and space
   before the first layout. Avoid an H2 nested immediately inside a
   bordered `md-layout` wrapper just because the container exists.
2. **Visual orientation beat:** A small overview or two-part comparison
   can orient the reader. A road or chart may then show *sequence* or
   *measurements* rather than restating those same cards. Use
   `{data-gap="lg"}` on supported layouts when their children need space.
3. **Quiet explanatory beat:** Simple H2 + paragraph, perhaps a short
   list, without decoration. Let the eye rest between visual clusters.
4. **Secondary beat:** Keep a compact visual if it answers a *different*
   question—such as which route to take after the road shows the sequence.
   Distinguish visual redundancy from legitimate information structure.
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
- Avoid placing a road, identical icon-feature recap, and identical
  step-card recap back to back. A road plus cards for a **separate decision**
  or **useful side routes** is a stronger design, not automatically clutter.
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
repository-neutral, one-file Markdown example. Notice the **different
jobs**: orientation cards tell readers what matters; the road shows
sequence; decision cards compare possible routes; one icon feature
explains a consequential moment; ordinary paragraphs carry detailed
instructions. The page has both visual anchors **and** room to read.
The final status notice is reserved for an actual caution.

**Design exercise:** compare two revisions. One may be too dense because
it repeats each road milestone as a card; another may be too flat because
it removed the overview, meaningful comparison, and useful wayfinding.
The strongest revision preserves the distinct visual jobs and edits
spacing, copy length, and typography instead of collapsing the document
into a numbered outline.

**Mermaid is optional**, not the default hero for a multi-section visual
post. Use it for exact connections or complicated loops that a road
and short prose cannot convey. Check runtime support first:
`architecture-beta` is unavailable in Mermaid 10.9.5. Validate rendered
output, not just fenced source.

### Example: a chart-driven article

Use `examples/chart-composition.md` for a complete generic one-file
data story with valid Chart.js **bar**, **line**, and **doughnut** blocks.
The page opens with two small, grounded overview cards, then gives
each chart full reading width, with prose between charts and a
build-time data disclosure for the exact values.

Read `topics/chart-authoring.md` for the choice between those charts,
filled lines, horizontal and stacked bars, scatter plots, label and
axis design, and accessible fallback behavior. Do not insert a chart
merely to make a report look visual when no measured data exists.

## Compare renders, not just Markdown

When improving an existing visual document, identify three things in the
previous render: its **strongest visual orientation**, its **clearest
decision/comparison**, and the **content that is genuinely repetitive**.
Do not delete the first two while fixing the third. Side-by-side preview
helps identify a regression that source-based checks will miss: a more
spacious document can still be a less helpful visual document.

Check the page at its actual reading width, not only at a full browser
window. Narrow columns can make otherwise concise card copy wrap densely.
Before removing cards, try shorter titles, shorter card descriptions,
a more readable column measure, and spacing between separate groups.

**Title ownership:** Check whether the Markata template already renders the
frontmatter title above the article. If so, avoid repeating the same H1
in the Markdown body. A standalone file without template title chrome
may use its own H1. The rendered page should not show two identical
adjacent article titles.

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
or table was explicitly requested. If every paragraph is inside a card,
move detailed prose out while retaining useful overview and comparison
cards. If a redesign leaves only a road and long numbered paragraphs,
consider restoring compact orientation and choice visuals. Neither extreme
is the desired outcome.
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
