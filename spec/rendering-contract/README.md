# Rendering contract

`contract-v1.json` is the single authoring source for presentation IDs,
variants, roles, defaults, bounds, and migration aliases.

The default fontpack is `system`, which uses native browser/platform font
stacks without downloads. Named packs such as `brush` remain explicit opt-ins.

Use the contract version and its migration table before changing theme data.
Add shared palette families and presentation dials here first. Do not add a
second registry in a site or consumer.

Regenerate all consumer projections from the Markata-Go repository root:

```bash
go run ./scripts/rendering-contract
go run ./scripts/rendering-contract --check
```

The generated Go and browser artifacts must be committed with a source change.

## Motif layer invariant

The page background is the base. The motif is an independent decorative pass.
`under` renders the motif below the surface texture; `over` renders it above
the texture; `sandwich` renders one motif pass below and one above the texture.
The surface texture mix controls only the texture pass. A zero texture mix does
not remove a configured sandwich motif. A zero motif mix makes the motif
background-equivalent, while one uses the full selected motif color. Geometry
(`size`, `gap`, `row_offset`, `wobble`, and `scatter`) remains active at every
mix value. `size` is the width of one visible motif mark, not a multi-mark
source tile. `gap` controls the space around that mark. A configured motif
remains fully opaque at every non-off mix value; `color_mix` only interpolates
between the background color and the selected motif color.
