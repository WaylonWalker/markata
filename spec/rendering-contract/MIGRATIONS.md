# Rendering contract migrations

Contract version 1 uses nested semantic theme values and normalized numeric
dials. Canonical nested values always win over legacy aliases. A conflict is
actionable: consumers warn with both keys and the canonical replacement.
Serialization emits only the nested form.

| Legacy field | Canonical field | Conversion |
| --- | --- | --- |
| `texture_strength` | `theme.texture.color_mix` | percentage or number divided/clamped to `0..1` |
| `heading_texture_strength` | `theme.heading_texture.color_mix` | percentage or number divided/clamped to `0..1` |
| `motif_color_distance` | `theme.motif.color_mix` | percentage or number divided/clamped to `0..1` |
| `texture_scale` | `theme.texture.scale` | percentage divided by 100 |
| `heading_texture_scale` | `theme.heading_texture.scale` | percentage divided by 100 |
| `motif_row_offset` | `theme.motif.row_offset` | percentage divided by 100 |
| `motif_wobble` | `theme.motif.wobble` | percentage divided by 100 |
| `motif_scatter` | `theme.motif.scatter` | percentage divided by 100 |
| `texture_scope = "headings"` | `theme.heading_texture.kind` | move the old texture kind to heading wear, set `theme.texture.kind = "none"`, and warn |
| `motif_layer` | `theme.motif.layer` | unchanged |
| `motif_color` | `theme.motif.color` | unchanged |

The legacy `texture_scope = "headings"` value is a mode migration, not a
direct scope rename. It means that the texture applies to headings only. The
migration copies the texture kind to `theme.heading_texture.kind` when needed
and sets `theme.texture.scope = "quiet"`. Do not emit `scope = "headings"`.

`heading_texture.kind = "inherit"` means use the surface texture kind while
retaining the heading texture's own `color_mix` and `scale`. The legacy
`texture_scope = "headings"` form is different: it disables surface texture
and moves that texture kind into heading wear.

The portable source is `contract-v1.json`. Browser and Go copies are generated
artifacts and must not gain local registry entries. Add palette families,
variants, and presentation dials to the portable source first, then regenerate
the projections and run the generated-artifact check.

## Default motif is off

The contract default for `theme.motif.kind` is `"off"`. Sites that relied on
the implicit `block-w` motif must set it explicitly:

```toml
[markata-go.theme.motif]
kind = "block-w"
```

## Default fontpack is now system

When `theme.fontpack` and the legacy top-level `fontpack` are both omitted,
the contract resolves to the zero-download `system` stack, matching the
core Markata configuration default. This is a default-value change, not a
migration of explicitly configured typography. Sites choosing `brush` or
`brush-poster` keep their selected fontpack; per-page overrides and visitor
picker choices also remain available.
