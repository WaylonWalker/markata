# Rich Markdown Authoring

Use this topic when the user asks for modern article UI such as icon cards,
roadmap feature rows, responsive comparisons, diagrams, status metrics,
sortable tables, code copy, checklists, chart fallbacks, or source disclosures.

## Prefer existing primitives before adding dependencies

1. Inspect the site's existing Markdown and theme; do not invent an icon
   library, client-side framework, or custom plugin for these patterns.
2. Use `:lucide-rocket:` (or `:lucide/rocket:`) for ordinary inline icons.
   The built-in `icons` plugin resolves project-local SVGs first, then
   optionally vendors the pinned Lucide pack on first use. It emits inline
   SVG without browser-time image requests or JavaScript.
3. For card layouts, use `md-columns`, `md-grid`, `md-card`, `md-box`.
4. For a centered icon in a card, use `md-icon`. For a left-hand icon tile
   with body text (as in ChatGPT feature lists), use `md-feature` with
   `md-feature-icon` and `md-feature-body`.
5. For data, prefer semantic Markdown tables, `md-status`, `md-stat`,
   `md-meter`, `md-progress` and `md-delta`.
6. For interactive behavior, use existing opt-ins: Chart.js fences with
   build-time table fallback, `md-checklist` with optional local persistence,
   `md-source` native details, or code blocks with optional Copy controls.
7. Validate render at mobile widths and with JavaScript disabled. Always
   show visible heading/label text when an icon communicates meaning.

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
