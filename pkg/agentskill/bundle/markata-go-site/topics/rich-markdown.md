# Rich Markdown Authoring

Use this topic when the user asks for modern article UI such as icon cards,
roadmap feature rows, responsive comparisons, diagrams, status metrics,
sortable tables, code copy, checklists, chart fallbacks, or source disclosures.

For broad requests to **visualize**, create a **visual post**, draw a
**diagram**, or build a **dashboard**, read
`topics/visual-authoring.md` to compose a responsive article layout.
When the content involves **bar charts, lines, trends, data comparisons,
distributions or scatter plots**, also read `topics/chart-authoring.md`
and author actual `chartjs` fenced JSON. Numerical charts complement
`md-stat`, cards, and prose rather than being replaced by them.
The skill is vendor-neutral: do not copy proprietary chat-chart syntax
into a Markata Markdown page.

**A visually designed article is neither Mermaid-plus-table nor
a wall of cards, but a balanced sequence of useful visual roles.** A workflow
may need a compact overview, a `road steps` sequence, a distinct
`md-columns` / `md-card` decision group, and occasional
`md-feature` rows. Keep these when they answer *different questions*;
remove repeated descriptions, not useful orientation. Place longer
instructions in quiet prose beneath meaningful H2 headings instead of
squeezing them into narrow cards. Keep prominent H2 headings outside
boxes so section hierarchy remains strongest; nested container headings
also retain their separate size and weight hierarchy. Use
`examples/editorial-visual-composition.md` as a generic reference for
purposeful visual variety, not a component quota. When authoring a single
Markdown file outside an existing Markata site, you can still write these
containers; only their styled rendering depends on Markata.

**During visual revisions**, compare the actual before/after pages:
reducing density by removing overview and side-route cards can make a page
less useful even if it has more whitespace. Prefer shorter card copy and
more separation between groups. If the site template already renders
frontmatter.title, avoid duplicating it as the first body H1.

## Prefer existing primitives before adding dependencies

1. Inspect the site's existing Markdown and theme; do not invent an icon
   library, client-side framework, or custom plugin for these patterns.
2. Use `:lucide-rocket:` (or `:lucide/rocket:`) for ordinary inline icons.
   The built-in `icons` plugin resolves project-local SVGs first, then
   optionally vendors the pinned Lucide pack on first use. It emits inline
   SVG without browser-time image requests or JavaScript.
3. For card layouts, use `md-columns`, `md-grid`, `md-card`, `md-box`.
   For editorial stories, use `content_width: wide` frontmatter and `md-section`
   containers. A leading H2 plus its next paragraph automatically become the
   sticky section head; `## 2026 · Title` gets a display kicker without changing
   its full-text anchor. Use `md-grid rail` for linked overview cards (four
   desktop columns by default, or `data-columns="2|3|4"`, stacked below 1000px)
   and `md-grid pending` for dashed unresolved-item cards. Rail ticks remain
   visible above cards; an immediately following `md-section` gets responsive
   spacing automatically. No local CSS or JS
   is necessary; see the Markdown Containers guide for complete examples.
4. For a centered icon in a card, use `md-icon`. For a left-hand icon tile
   with body text (as in visual reports and feature lists), use `md-feature` with
   `md-feature-icon` and `md-feature-body`.
5. For data, prefer semantic Markdown tables, `md-status`, `md-stat`,
   `md-meter`, `md-progress` and `md-delta`.
6. For measured trends and category comparisons, use `chartjs` blocks rather
   than Mermaid or a static diagram. Read `topics/chart-authoring.md`
   for valid bar/line/stacked/doughnut examples, scales, and fallbacks.
7. For interactive behavior, use existing opt-ins: Chart.js fences with
   build-time table fallback, `md-checklist` with optional local persistence,
   `md-source` native details, or code blocks with optional Copy controls.
8. Validate render at mobile widths and with JavaScript disabled. Always
   show visible heading/label text when an icon communicates meaning.
9. If the site enables CSS purge, ensure its preserve patterns include
   `md-code-*` for the JS-injected Copy button. Built-in defaults now include
   it; custom `preserve` lists replace the defaults. Verify the *published*
   code.css still has `.md-code-copy` and `.md-code-enhanced` positioning.

## Icons alongside text

````markdown
::: md-feature

:::: md-feature-icon
:lucide-rocket:

::::

:::: md-feature-body
### Verify production

Check the deployed image and the live site.

::::

:::

::: md-feature

:::: md-feature-icon
:lucide-gauge:

::::

:::: md-feature-body
### Protect performance

Validate slow 3G and page-weight budgets.

::::

:::
````

The icon is decorative and aria-hidden; do not make the icon the sole carrier
of meaning. For truly interactive controls, use a readable button label too.

## Side-by-side cards with centered Lucide icons

````markdown
::: md-columns

:::: md-card {data-align="center"}

::::: md-icon {data-size="lg"}
:lucide-server:

:::::

### Falcon2
K3s node

::::

:::: md-card {data-align="center"}

::::: md-icon {data-size="lg"}
:lucide-database:

:::::

### Falcon3
NAS + ZFS

::::

:::
````

Other sizes are `data-size="sm"` and the default icon size. These CSS-only
containers retain readable content without JS, stack on narrow widths and
inherit the active theme. Markdown container closers must be bare colon lines
that match their corresponding opening length; use blank lines before closers.

## Distinguish icon shortcodes from component libraries

- `:lucide-rocket:` is an existing built-in Markdown shortcode, **not** a
  Web Awesome component, JS dependency or new Markdown extension.
- If a shortcode stays literal in output, check that `icons` is enabled,
  that the name exists in the pinned/local pack, that offline builds have the
  pack cached, and that the shortcode is not inside inline/fenced code.
- `wa-*` components require Web Awesome configuration and must not be used
  just to show a simple line icon.
- Icons may use local overrides at `.icons/<pack>/<name>.svg` or
  `static/.icons/<pack>/<name>.svg`. They must pass the SVG safety checks.
- When composing metrics, prefer visible values and text labels. Never claim
  a chart, operational value or deployment is verified by the renderer.

## Verify and install this guidance

In a site repo, `markata agent doctor` checks whether its installed skill
matches the current binary bundle. After installing a build that contains
this guidance, use `markata agent update` then `markata agent doctor`.
Updating the source repository alone does not update older installed skill
copies or deployed sites.

Reference docs in the Markata repository:
- `docs/guides/icons.md`
- `docs/guides/markdown-containers.md`
- `docs/guides/markdown-metrics.md`
- `docs/guides/markdown-checklists.md`
- `docs/guides/markdown-chart-fallback.md`
- `docs/guides/markdown-source-disclosure.md`

Rail card grids stay connected when responsive: horizontal above 1000px, vertical
with ring ticks below 1000px. Use the rail modifier instead of custom connector CSS.

### Standalone rendering and Python applications

`markata render --json` reads Markdown from stdin and returns HTML, feature
signals, and ordered assets without building a site. `--inline-assets` includes
bundled CSS/JS (including Mermaid and Chart.js). The official Python wheels expose
`from markata import render, render_async`; use `render_async` in async FastAPI
handlers and `result.document()` for a complete HTML response. For existing
layouts, insert the fragment inside `.post-content` and include both asset tag
fields. See `docs/guides/python-renderer.md` for asset serving and API scope.

For application integrations, `markata.appearance()` provides the site's picker and
asset catalog. Use `result.document(appearance=...)` for site article layout instead
of copying templates or modifying rendered HTML. Standalone rendering includes
bundled Lucide icons. CodeMirror and application editing behavior stay application-side.

## Browser-local WASM rendering

For browser-only Markdown apps, build `scripts/wasm-workbench.py --output DIR`.
Use its persistent `worker.js` rather than a rendering server. Wait for `ready`,
then send `{source, revision}`. Accept only the newest result. Serve the generated
`appearance/` assets and the matching `wasm_exec.js` beside `markata.wasm`.
Raw HTML is disabled; site-dependent plugins still require a full site build.
Wait for the demo's **Ready offline** status before testing offline reload.
See `docs/guides/wasm-workbench.md` for hosting headers and browser checks.
For a server backed Go editor, run `go run ./cmd/markata-workbench` from the
Markata source repository. The default URL is `http://127.0.0.1:8766/`; use
`--listen 127.0.0.1:8767` when that port is occupied. Applications can embed
`workbench.New()` from `pkg/workbench` at an origin root. It renders in process,
serves local appearance assets, and preserves drafts in browser storage.
This trusted local editor enables authored HTML inside a sandboxed preview.
It does not load a site's configuration or run its full build lifecycle.
See `docs/guides/go-workbench.md` for the HTTP contract and embedding example.

The Go workbench follows the Python preview protocol: request the full document
once, then set `include_document: false` and update the existing sandboxed frame
with `html`, `feature_head_html`, and `body_end_html`. Keep shared appearance
styles loaded between edits. Do not replace the iframe document on each render.

For renderer integrations, reuse the benchmark site post and rendering fixtures.
Run `python python/benchmark_editor.py` with an installed wheel and FastAPI dependencies.
See `docs/guides/renderer-integration-benchmarks.md` for timing boundaries and host-specific results.


For Go HTTP benchmarks, build `cmd/markata-workbench` and run
`python python/benchmark_go_editor.py --binary /path/to/markata-workbench`.
Use the same fixture hashes as the Python report. Keep cache misses separate
from repeated cached renders, and distinguish linked assets from inline assets.
See the integration benchmark guide for timing boundaries and stacking order.

For comparable WASM fixture measurements, serve the browser workbench and run
`uv run --with playwright python scripts/benchmark-wasm.py --output /tmp/wasm-benchmark.json`.
The report measures worker rendering, excluding startup, transfer, layout, and paint.
Keep these timings separate from HTTP latency.
