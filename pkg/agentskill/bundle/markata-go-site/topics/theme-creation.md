# Theme Creation

Use this topic for palettes, colors, typography, texture, motifs, and CSS.

## Process

Follow these steps for theme work:

1. Discover the active rendering contract and its version.
2. Choose a canonical palette and, when needed, a palette family.
3. Make the smallest semantic configuration change.
4. Use CSS variables or `custom_css` as escape hatches.
5. Add palettes and presentation dials to the contract before adding projections.
6. Regenerate projections and run conformance checks.

The contract source is `spec/rendering-contract/contract-v1.json`. Read
`spec/rendering-contract/README.md` and `MIGRATIONS.md` before changing shared
IDs, defaults, bounds, or aliases. Do not create a second registry in a site
or consumer.

## Choose the Smallest Layer

- Choose a palette for a site-wide color system.
- Use `[markata-go.theme.variables]` for a few generated CSS variable overrides.
- Use `[markata-go.theme]` `custom_css` for selectors or components.
- Override templates only when the markup or layout must change.

Inspect the active config, `palettes/`, `templates/`, and site CSS before editing.
Preserve the site's existing layout and typography unless the task changes them.

`[markata-go.theme.variables]` entries are emitted in the `overrides` cascade
layer and win over text-size presets and palette tokens, so `--content-width`
or `--text-base` set there take effect without `!important`. The resolved
`palette`, `aesthetic`, `fontpack` and texture are stamped on `<html>` for
every page (posts and feeds alike); if a feed page looks untyped, check the
`data-fontpack` attribute in the built HTML before touching CSS.

Theme CSS uses `@layer reset, tokens, base, components, utilities, overrides`.
Any unlayered rule outranks all layered rules, so custom CSS that must lose to
palette tokens should be wrapped in `@layer components { ... }`, and a rule
that should win everything belongs in `@layer overrides`.

## Page Stylesheet Performance

For a page with many blocking stylesheets, consider
`[markata-go.css_bundle] enabled = true` with `pages = ["index.html"]`.
Validate with a full build: fast mode skips page bundles. Measure before and
after through the same server and compression policy, and check theme controls
and view-transition navigation. Keep asset URLs root-relative when bundling;
files with relative URLs remain separate to preserve their meaning.

## Canonical Theme Configuration

Keep theme-level keys in `[markata-go.theme]`. Put each nested group under its
own table. In TOML, a key remains in the current table until another table
header appears. Therefore, keep `custom_css` in the theme table, not in
`[markata-go.theme.motif]` or another nested table.

```toml
[markata-go.theme]
contract_version = 1
palette = "ayu-dark"
aesthetic = "minimal"
fontpack = "system"
custom_css = "custom.css"
text_size = "large"
show_text_size_control = true

[markata-go.theme.texture]
kind = "screenprint"
color_mix = 0.35
scale = 1.0
scope = "all"

[markata-go.theme.heading_texture]
kind = "inherit"
color_mix = 0.45
scale = 1.0

[markata-go.theme.motif]
kind = "block-w"
glyph = "W"
size = "78px"
gap = "10px"
row_offset = 0.24
wobble = 0.18
scatter = 0.0
layer = "sandwich"
color = "ink"
color_mix = 0.01
url = "https://waylonwalker.com/w.svg"

[markata-go.theme.variables]
"--content-width" = "68ch"
"--color-primary" = "#8b5cf6"
```

The contract uses normalized numeric dials. `color_mix`, `row_offset`,
`wobble`, and `scatter` use `0..1`. Texture scales use `0.25..3`. Use enum
values from the contract for palettes, aesthetics, textures, heading textures,
motifs, motif colors, motif layers, and scopes.

The default theme's reading presets are `small`, `medium`, `large`, and
`x-large`. `text_size = "large"` selects the default 18px site base, 22px
article text and a 64ch measure (`--content-width` is in `ch`, so it tracks
the article font size). Article text also scales with `--reading-scale` on
wide viewports (1.08x at ≥1800px, 1.16x at ≥2200px, 1.25x at ≥3000px), so do
not add breakpoint overrides for large monitors. Set `show_text_size_control = false` to hide the visitor selector
while keeping the configured default. Visitor selections are saved per site in
browser-local storage. Do not add `--text-base` or `--content-width`
overrides just to make articles "readable"; the presets already target a
65-70 character measure. Override them only for a deliberate design choice.

Borders in the default theme are soft by design: `--color-border` is the text
ink mixed 16% into the background, `--color-border-soft` is used for section
dividers, and `--color-border-strong` is the full-contrast ink for focus
rings. Headings use a short accent bar (`--heading-rule`) rather than
full-width rules. Radii follow `--radius-sm`/`--radius`/`--radius-lg`/
`--radius-xl`. Prefer adjusting these tokens over re-styling components.

The feed and document sidebars are drawers at every width, opened from an edge
handle (or `b` / `Shift+B`); they never open on hover. On wide screens
(>= 1201px) an open drawer docks and pushes the article aside, and the open
state is remembered per side. At 1200px and below, drawers are overlays with a
backdrop, opened from bottom-corner pills, one at a time. Width comes from the
`--feed-sidebar-width` / `--doc-sidebar-width` custom properties on
`.page-wrapper`. If a site's custom CSS positions `.feed-sidebar` or
`.doc-sidebar`, remove it and rely on the theme.

Every site ships a live theme picker by default. The header shows one
palette-swatch button and the light/dark toggle, and visitors can preview and
choose any palette. The configured `palette` and `aesthetic` remain the
defaults for first-time visitors. Do not add extra theme buttons to the nav.
To lock the site to its configured palette, set:

```toml
[markata-go.theme.switcher]
enabled = false
```

Use `include` or `exclude` under the same table to limit the palettes offered.
The picker has Colors, Style (aesthetic), and Font (every bundled fontpack,
plus text size) tabs of live preview cards, with ‹ › buttons to step through
the active tab. Colors includes a **Seasonal** option: northern hemisphere
seasons, switching to holiday palettes a few days before and on world
holidays. To make that the default for visitors, add `seasonal = true` to
`[markata-go.theme]` and keep `palette` as the fallback. Under `markata-go serve` only, its **Bake** button writes the current choices
into `[markata-go.theme]` (color mode and text size only when the visitor
picked them). It edits whichever config file already holds the
theme table (including `include`d files such as `config/theme.toml`), keeps
comments, and triggers a rebuild. When a user says "use the look I picked",
have them click Bake in `serve`, then review the diff of the named file. The
default fontpack is the zero-download `system` pack (native system-ui
body/headings and native monospace code). `brush` and `brush-poster` remain
explicit opt-ins for branded poster headings. Sites are dark by default
(`fallback_mode = "dark"`).
Builds with no config file (such as `markata-go build post.md`) use the
contract defaults (`ayu-dark`, `minimal`). The motif is off by default; set
`[markata-go.theme.motif] kind = "block-w"` only when the site wants it.
Enabled motifs recolor automatically when a visitor picks another palette, so
do not hard-code motif colors per palette.

`heading_texture.kind = "inherit"` uses the surface texture kind while keeping
the heading texture's own `color_mix` and `scale`.

## CSS Escape Hatches

Put direct CSS variable overrides in `[markata-go.theme.variables]` when the
change is a token change. Put selector-level rules in the file named by
`[markata-go.theme]` `custom_css`:

```toml
[markata-go.theme]
custom_css = "custom.css"
```

Place that file at `static/custom.css`. Use templates only when CSS cannot
change the required structure.

## Palettes and Semantic Roles

Use `markata-go palette list` to discover available palette IDs. Use
`markata-go palette info <name>` to inspect a palette and
`markata-go palette check <name>` to check its contrast.

Use semantic palette roles such as `accent`, `background`, `ink`, and `surface`
instead of styling individual selectors. For a site-local palette, read
`../reference/palette-reference.md` for the `palettes/<name>.toml` shape.

When adding a shared palette, add its family, ID, variant, and roles to
`contract-v1.json` first. Add the light and dark members in the same change
when the family supports both variants; single-variant palettes get a derived
`<id>-light`/`<id>-dark` counterpart automatically (sync its contract roles
with `go run ./scripts/rendering-contract --derive-palettes`). Regenerate all
projections before using the new ID.

Every palette has a working light/dark toggle. Do not add `palette_light` or
`palette_dark` just to "enable" light mode; only set them to choose a
different hand-tuned partner than the derived one.

```bash
go run ./scripts/rendering-contract
go run ./scripts/rendering-contract --check
```

The generator writes projections for the repository and its supported consumer
repositories. Run it from the markata-go repository root.

## Fonts

Select a `fontpack` from the contract or from `markata-go fonts packs` (built-in
catalog packs and aliases such as `typewriter` or `reader` are valid). A fontpack assigns named roles instead
of requiring each selector to name a font family. The contract roles are
`body`, `heading`, and `mono`; use the roles supplied by the selected pack. The
`mono` role controls code and diagrams. Do not invent a fontpack ID or a role
that the catalog does not provide.

## Migration Rules

Canonical nested values take precedence over legacy flat values. The loader can
accept legacy values as migration inputs, but canonical serialization emits
only nested values. Migration warnings report explicit conflicts once, not differences from built-in defaults.
If both forms conflict, read the migration warning and
replace the legacy value.

The legacy `texture_scope = "headings"` setting needs special care. It means
headings-only behavior; it does not mean that `scope = "headings"` is a valid
canonical value. The migration moves the texture kind to
`theme.heading_texture.kind`, sets the surface texture to `none`, and uses the
quiet surface scope. Preserve this meaning when converting old configuration.

## Conformance Checks

Run the smallest relevant checks after editing:

```bash
markata-go palette check my-brand
markata-go palette check my-brand --strict
go run ./scripts/rendering-contract --check
go test ./pkg/config ./pkg/renderingcontract
```

Use `markata-go theme check-all` to check all discovered palettes. Add
`--colorblindness` for color-vision warnings. The command returns an error when
any palette fails its selected contrast checks.
