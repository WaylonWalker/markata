# Icon Shortcodes Specification

The icons plugin expands local SVG icon packs into inline HTML during the Transform stage and provides a zero-config Lucide fallback.

## Goals

- Support Zensical-compatible icon shortcodes in Markdown.
- Make `:lucide-…:` shortcodes useful on a fresh install without manual asset setup.
- Keep builds deterministic by pinning automatically vendored Lucide.
- Keep ordinary builds network-free unless a missing Lucide shortcode is actually encountered.
- Prefer project-local SVG assets over automatically vendored assets.
- Avoid changing emoji shortcodes such as `:smile:` when no matching icon asset exists.
- Never expand or vendor icon syntax inside inline code or fenced code blocks.

## Asset Layout

The plugin searches these directories by default, in order:

1. `.icons/`
2. `static/.icons/`

SVG paths below an icon root define pack and icon names. For example:

```text
static/.icons/
└── lucide/
    └── smile.svg
```

produces the logical icon name `lucide/smile`.

When the same logical icon exists in more than one configured path, the first configured path wins. Local icons MUST win before the automatic Lucide fallback is considered.

## Markdown Syntax

For an asset named `lucide/smile`, the canonical Zensical-style Markdown syntax is:

```markdown
:lucide-smile:
```

markata-go also accepts the slash form originally proposed by issue #895:

```markdown
:lucide/smile:
```

Both forms resolve to the same icon. Nested paths are flattened by replacing `/` with `-` for the canonical alias. If two different icon paths would produce the same flattened alias, the ambiguous hyphen alias is disabled; the slash forms remain available.

Unknown shortcodes are preserved verbatim. This lets Goldmark's emoji extension continue handling ordinary emoji shortcodes independently.

## Configuration

```toml
[icons]
enabled = true
auto_vendor = true
paths = [".icons", "static/.icons"]
packs = ["lucide", "simple-icons"]
```

Fields:

| Field | Type | Default | Description |
|---|---|---|---|
| `enabled` | bool | `true` | Enable icon shortcode expansion. |
| `auto_vendor` | bool | `true` | Allow the built-in pinned Lucide fallback when a local Lucide icon is missing. |
| `path` | string | — | Convenience form for a single icon root. |
| `paths` | list[string] | `[".icons", "static/.icons"]` | Ordered local SVG roots. |
| `packs` | list[string] | all packs | Optional allowlist of first path components. |

`paths` takes precedence over `path` when both are set.

If `packs` is configured and does not include `lucide`, the automatic Lucide fallback MUST NOT run. `auto_vendor = false` disables only the zero-config Lucide fallback; explicitly configured `[icons.vendor]` packs remain independent.

## Rendering

Resolved icons are emitted as inline SVG wrapped in a span:

```html
<span class="icon" data-icon="lucide/smile" aria-hidden="true">
  <svg width="1em" height="1em" aria-hidden="true" focusable="false">...</svg>
</span>
```

The wrapper is vertically aligned with surrounding text. Source SVG `width` and `height` attributes are normalized to `1em`; the SVG `viewBox` and drawing attributes are preserved.

Icon shortcodes are decorative by default and are hidden from assistive technology. Authors should include visible text when an icon carries meaning.

The plugin serializes SVG source on one line before inserting it into Markdown.
This keeps multiline SVG files from being split into Markdown paragraphs and
escaped instead of rendered as inline SVG.

## Article Compositions

In the default article theme, the existing inline SVG shortcode may be
composed with static Markdown container classes. The icon resolver MUST remain
independent of the presentation class so icons work in prose too.

- `md-icon` centers and sizes an icon within a card or block, using
  `data-size="sm"` (1rem), default (1.25rem), or `data-size="lg"` (1.75rem).
- `md-feature` lays out an icon tile and Markdown content without JS.
- `md-feature-icon` contains a regular `:lucide-…:` shortcode.
- `md-feature-body` contains meaningful headings, paragraphs, links or lists.

These article-scoped layouts MUST NOT affect existing global `.icon`,
`.card`, `wa-*` or feed-card styles. They MUST maintain a mobile layout that
does not overflow the article width. Unknown icon names retain their literal
shortcode; do not silently replace missing glyphs with an unrelated icon.
SVGs remain decorative and aria-hidden; an icon conveying meaning MUST have
visible adjacent text.

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
````

## Code Protection

The plugin must not transform or vendor shortcodes inside code:

````markdown
Visible :lucide-smile:

`literal :lucide-smile:`

```text
:lucide-smile:
```
````

Only the first shortcode is eligible for expansion and automatic vendoring.

## SVG Safety

Local and vendored SVGs are rejected when they contain active-content primitives that would be unsafe to inline:

- `<script>`
- `<foreignObject>`
- inline event-handler attributes such as `onload=`
- `javascript:` URLs

Rejected assets behave like missing icons and leave the shortcode unchanged.

## Lifecycle

The plugin runs at `PriorityLast` in the Transform stage. This lets source transforms such as Jinja run first while still expanding icons before Markdown rendering.

Local icon roots are indexed during Configure. If Transform encounters an unresolved Lucide shortcode and `auto_vendor` is enabled, the plugin MUST:

1. synchronize the fallback so concurrent post transforms do not download the same pack repeatedly;
2. vendor the pinned default `lucide-static` pack using the existing icon-vendor implementation;
3. materialize its SVGs and license under an icon root scanned by the renderer;
4. re-index the icon roots; and
5. retry the shortcode lookup.

The automatic fallback MUST be attempted at most once per plugin instance/build.

## Default Lucide Vendoring

The built-in fallback uses a markata-go-pinned `lucide-static` version. It MUST NOT follow an upstream `latest` tag. The current pin is `1.48.0`.

The fallback reuses the same archive cache, SVG safety validation, materialization, and license preservation implemented by `IconVendorPlugin` rather than maintaining a separate downloader.

If the default pack cannot be materialized—for example, a first-use build is offline with no cached/materialized pack—the site build MUST continue and the unresolved shortcode MUST remain literal. This best-effort behavior applies only to the automatic fallback; explicitly configured `[icons.vendor]` packs retain their strict error behavior.

Builds with no unresolved Lucide shortcodes MUST NOT trigger the automatic vendor path or network access.
