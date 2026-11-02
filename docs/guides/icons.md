---
title: "Icon Shortcodes"
description: "Use Lucide and local SVG icon packs directly in Markdown"
published: true
slug: /docs/guides/icons/
tags:
  - documentation
  - markdown
  - icons
---


markata-go renders Zensical-style icon shortcodes directly in Markdown.

Lucide is ready by default. You can write:

```markdown
A friendly icon :lucide-smile: inline with text.
```

or the slash alias:

```markdown
A friendly icon :lucide/smile: inline with text.
```

No icon directory or configuration is required. The first unresolved Lucide shortcode vendors a pinned `lucide-static` pack through markata-go's existing icon-vendor cache, indexes the resulting local SVGs, and renders the requested icon. Builds that never use a Lucide shortcode do not download the pack.

The default Lucide version is pinned by markata-go rather than following an upstream `latest` tag. Once vendored, later builds reuse the local materialization/cache and do not need to fetch it again.

## Icons in cards, diagrams, and feature lists

These are the same thin line icons shown in modern ChatGPT-style feature
cards. **No new icon plugin or browser JavaScript is needed.** The `icons`
plugin already emits inline SVG, and the default article theme supplies
opt-in containers to size and position it.

### Centered icons in side-by-side cards

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

`md-icon` also supports `data-size="sm"` (1rem) and the default 1.25rem
size. The `data-size` values affect only presentation; the source SVG
remains a 1em icon that inherits the current text color.

### Icon-first feature rows

````markdown
::: md-feature

:::: md-feature-icon
:lucide-rocket:

::::

:::: md-feature-body
### 1. Verify production

Confirm the new image and the live page revision. Test without JavaScript.

::::

:::

::: md-feature

:::: md-feature-icon
:lucide-gauge:

::::

:::: md-feature-body
### 2. Protect performance

Test slow 3G and keep page-weight budgets.

::::

:::
````

The icon tile and text reflow within the article width, including on narrow
mobile screens. The icon itself is *decorative* and is hidden from assistive
technology; the visible heading conveys the meaning. Avoid icon-only buttons
without a readable text label.

### Which syntax should an agent choose?

- In ordinary prose, use `:lucide-rocket:` (or `:lucide/rocket:`).
- In a card, wrap the shortcode with `::: md-icon` for centered/sized display.
- In a list like a feature roadmap, use a `md-feature` row with nested
  `md-feature-icon` and `md-feature-body` containers.
- For arbitrary connected graphs, prefer Mermaid; use `md-feature` for
  icon-plus-text rows rather than inventing new HTML/CSS.
- In a custom post template, reuse existing icons and styles where appropriate;
  do not hardcode a second SVG library, icon font, or CDN.

The diagram and layout syntax is also documented under
[Markdown Containers](/docs/guides/markdown-containers/). See the bundled
agent skill's `topics/rich-markdown.md` for authoring examples that agents
can use outside the Markata source repository.

## Local icons and overrides

Project-local SVG packs still work and take precedence over automatically vendored Lucide icons. Put SVG files under `.icons/` or `static/.icons/`.

For this asset:

```text
static/.icons/lucide/smile.svg
```

both `:lucide-smile:` and `:lucide/smile:` use your local SVG.

Unknown names are left unchanged, so normal emoji shortcodes such as `:smile:` can still be handled by the Markdown emoji extension.

## Configuration

Icon rendering and the default Lucide fallback are enabled by default:

```toml
[icons]
enabled = true
auto_vendor = true
paths = [".icons", "static/.icons"]
```

To disable the automatic Lucide fallback while continuing to use local icons:

```toml
[icons]
auto_vendor = false
```

To use only selected packs:

```toml
[icons]
packs = ["lucide", "simple-icons"]
```

If `packs` is set and does not include `lucide`, the automatic Lucide fallback is not used.

To use one custom icon root:

```toml
[icons]
path = "assets/icons"
```

To disable icon shortcodes entirely:

```toml
[icons]
enabled = false
```

For explicit pinning, additional icon packs, custom sources, or strict build-time vendoring, see [Vendoring Icon Packs](/docs/guides/icon-vendoring/).

## Offline builds

The automatic Lucide fallback reuses markata-go's icon vendor cache. If the pack has already been cached or materialized, it works offline.

With `MARKATA_GO_OFFLINE=1`, markata-go never downloads a missing pack. If a Lucide shortcode cannot be resolved from local or cached assets, the shortcode is left unchanged instead of failing the whole site build.

## Naming

The path below an icon root becomes the logical icon name. Slashes are replaced with hyphens for the canonical Markdown alias:

| SVG file | Canonical shortcode | Slash alias |
|---|---|---|
| `lucide/smile.svg` | `:lucide-smile:` | `:lucide/smile:` |
| `simple-icons/github.svg` | `:simple-icons-github:` | `:simple-icons/github:` |
| `custom/arrows/left.svg` | `:custom-arrows-left:` | `:custom/arrows/left:` |

If two paths flatten to the same hyphenated alias, markata-go leaves that ambiguous alias unresolved. The explicit slash spellings continue to work.

## Code examples stay literal

Icon shortcodes are not expanded or vendored inside inline code or fenced code blocks:

````markdown
This renders: :lucide-smile:

This stays literal: `:lucide-smile:`

```markdown
:lucide-smile:
```
````

## SVG behavior

Icons are inlined into the generated HTML, so rendered pages do not need browser-time requests for SVG files. Fixed source `width` and `height` values are normalized to `1em` so icons follow surrounding text size.

Because inline SVG becomes part of the page, markata-go rejects SVGs containing scripts, `foreignObject`, JavaScript URLs, or inline event handlers. Local and vendored SVGs go through the same safety check.

Shortcodes are treated as decorative (`aria-hidden`). When an icon communicates meaning, include readable text next to it rather than relying on the glyph alone.
