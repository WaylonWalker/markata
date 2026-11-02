---
title: "Rendering presentation contract"
description: "Configure portable palettes, typography, texture, and motif settings with the shared rendering contract."
date: 2026-08-09
published: true
tags:
  - documentation
  - themes
  - rendering
---


Markata-Go uses rendering contract version 1 for portable presentation settings.
The authoring source is `spec/rendering-contract/contract-v1.json`.

## Canonical configuration

Put theme-level settings under `[markata-go.theme]`. Put texture, heading
texture, motif, and variable settings under their nested tables.

```toml
[markata-go.theme]
contract_version = 1
palette = "ayu-dark"
aesthetic = "minimal"
fontpack = "brush-poster"
custom_css = "custom.css"

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
```

The motif is off by default (`kind = "off"`); the example above enables the
`block-w` motif explicitly.

In TOML, a key belongs to the current table until another table header
appears. Keep `custom_css` in `[markata-go.theme]`, not in `motif`.

## Contract values

The contract defines the shared IDs, defaults, enums, and bounds. Use its
palette IDs, aesthetic IDs, fontpack IDs, texture kinds, heading texture kinds,
motif kinds, colors, layers, and scopes.

The canonical dials use numbers:

- `color_mix`, `row_offset`, `wobble`, and `scatter`: `0..1`
- `texture.scale` and `heading_texture.scale`: `0.25..3`

`heading_texture.kind = "inherit"` uses the surface texture kind while keeping
the heading texture's own mix and scale.

Fontpacks assign the named roles `body`, `heading`, and `mono`. Select a
fontpack instead of inventing role names in site configuration.

## Migration from flat settings

Contract version 1 accepts legacy flat settings as migration inputs. Canonical
nested values always win when both forms are present. The loader reports a
warning when the values conflict. Canonical serialization emits only nested
settings.

The legacy `texture_scope = "headings"` value means headings-only behavior.
Do not convert it to `theme.texture.scope = "headings"`; `headings` is not a
canonical scope. The migration keeps the texture kind in
`theme.heading_texture.kind` and sets the surface texture scope to `quiet`.

Other documented aliases include `texture_strength`,
`heading_texture_strength`, `motif_color_distance`, `texture_scale`,
`heading_texture_scale`, `motif_row_offset`, `motif_wobble`,
`motif_scatter`, `motif_layer`, and `motif_color`. See
`spec/rendering-contract/MIGRATIONS.md` for the complete mapping.

## Add or change contract data

Change `contract-v1.json` first when you add a shared palette family, palette
variant, fontpack, aesthetic, texture, motif value, default, bound, or alias.
Do not edit generated projections by hand. Regenerate them from the repository
root, then run the generated-artifact check:

```bash
go run ./scripts/rendering-contract
go run ./scripts/rendering-contract --check
```

The generator also targets supported consumer repositories. Run it only when
those repositories are available, and do not add a local registry instead.

### New site

Start with the canonical example above. Select one palette ID, aesthetic, and
fontpack. `theme.fontpack` accepts contract fontpack IDs and any pack or alias
in the font catalog (`markata-go fonts packs`, or your `fontpacks_file`). Validate the site with the contract check and the
focused configuration tests before adding local variables.

### New palette

Add the family and both `light` and `dark` variants to `contract-v1.json`.
Include explicit family, variant, name, and semantic role metadata. Regenerate
all projections and run the 120-variant final-CSS sweep.

### New dial

Add the field, default, bounds, and migration alias to the contract first.
Then update normalization, serialization, render plans, CSS projection,
completion, and browser fixtures in Go, Plaindown, and Theme Lab. The drift
check must pass before the dial is accepted.

## Validation

Check a palette after changing its roles or colors:

```bash
markata-go palette check ayu-dark
markata-go palette check ayu-dark --strict
```

Run `go run ./scripts/rendering-contract --check` after contract changes.
Run `go test ./pkg/config ./pkg/renderingcontract` to exercise configuration
migration and contract loading.

## WCAG AA color gate

The rendering/theme gate checks the final semantic projection used by generated
CSS. Normal text, metadata, headings, links, code tokens, and button text use
4.5:1. Controls, focus indicators, borders that identify controls, status
colors, and other non-text UI use 3:1. A role is not treated as large text
because its name says muted, secondary, or comment.

Decorative texture, motif, and shadow layers have no independent contrast target,
but their mix and opacity must not lower readable content below these limits.
AAA is reported separately and is optional.

The canonical palette sweep validates every catalog variant and records the
final foreground/background pair and ratio. Plaindown Editor themes use a
separate sweep because Editor mode owns its chrome; Site mode uses the canonical
Markata palette projection. Both sweeps are required before changing a palette
or editor theme.
