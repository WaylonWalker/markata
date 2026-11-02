---
title: "Font packs"
description: "Choose offline system typography or a vendored Markata font pack."
date: 2026-08-06
published: true
tags:
  - documentation
  - typography
---


Markata defaults to the offline `system` pack. It makes no external font
requests and emits no font binaries. To choose another pack, add this to the
site configuration:

```yaml
fontpack: system-reader
```

The same value can be set as `theme.fontpack`. Any pack or alias listed by
`markata-go fonts packs` is valid. When `fontpacks_file` is set, names are
checked against that catalog when the site builds.

Catalog-defined bundled packs use vendored, stable WOFF2 tiers. The built-in
catalog, manifests, licenses, and font assets are embedded in the executable,
so these packs work from an installed or GoReleaser binary in any working
directory. Bundled builds do not run Python or FontTools during `markata build`.

The generated `css/fonts.css` has a content-hashed URL for caching. Pages
preload the emitted WOFF2 files for their body and main heading ahead of CSS;
per-page overrides use their own pack. A stored picker choice is restored
before paint and preloads that choice too. Nonempty editorial feeds also preload
the code face used by their visible week labels, so layout does not have to
discover that font later. Calendar controls reserve their space while loading
to keep the archive toolbar stable. Shared sources are fetched once, and system code
stacks need no download. Other layouts keep code fonts on demand. Additional
script coverage also loads when needed. Font faces use
`font-display: swap`, so slow networks can still briefly show fallback text;
preloading reduces that delay without forcing fallback fonts permanently.

System packs use local system stacks and transfer zero font files. Custom packs
are loaded from an explicit catalog and may reference local assets. Relative
paths in a custom catalog are resolved relative to the catalog file, not the
current working directory.

Inspect the catalog with:

```text
markata-go fonts packs
markata-go fonts show system
markata-go fonts doctor
markata-go fonts verify
markata-go fonts report
```

## Compare packs on one site

### Painted Sign

`painted-sign` combines Finger Paint display lettering with scratchy Rock Salt
headings and accents over a clean Source Sans 3 prose stack. DM Mono remains
available for code, so the expressive hierarchy does not compromise reading.

The `fontpack` frontmatter field overrides the site default for one page:

```yaml
---
title: "Brush comparison"
fontpack: brush-poster
---
```

All selected page packs share one generated stylesheet and deduplicated asset
set. The comparison pages in `pages/post/test-*.md` use this mechanism.

Markata tracks emitted files in `output/assets/fonts/.markata-fonts.json`.
Switching packs removes obsolete Markata-generated files while preserving
unrelated user-managed files in the same directory.

On warm builds, pages that already select the resolved pack keep their existing
HTML rather than copying the page body just to reapply the root attribute.
Changing the site or per-page pack still updates `data-fontpack`; custom
attributes and quoted values on the HTML element are preserved. Font catalog,
glyph coverage, preloads, and missing output files are still checked normally.

Cold resolution streams article HTML through one shared visible-character
analysis. Repeated characters and repeated pack selections do not repeat that
work, and a family's coverage checks stop once its full tier is required.
This does not change font choices or run a subsetter during your build.

Built-in output caching tracks visible coverage, the default pack, effective
pack choices, picker settings, and bundled catalog/manifest revisions.
Reordering equivalent selections or using equivalent aliases does not invalidate
the cache. Missing generated stylesheets or fonts are regenerated through the
same accelerated resolver; custom catalogs still validate their asset checksums.
A new cache schema or bundled metadata revision causes one fresh resolution.

Coverage is selected per font family. A family gets `latin-ext` only when its
manifest provides that tier; otherwise Markata falls back to that family's
`full` tier, so expressive faces without extended subsets still render all
text. Role `weight`, `style`, `size`, and `optical_size` values are emitted as
CSS properties and applied to the corresponding body, heading, code, lead,
quote, and caption elements.
When multiple packs are used, role variables and declarations are scoped to
`[data-fontpack="..."]`; optional properties from one pack cannot affect
another pack. A missing specialized role falls back to `body` as a whole, but
an existing specialized role does not inherit optional properties that it did
not configure. For example, a quote role without `optical_size` will not pick
up the body's optical-size setting.

The maintenance generator reads variable-font `fvar` axes from the source
font, including their real weight and optical-size bounds. To check generated
manifests without rewriting them, run it with `--check` and the pinned
google/fonts checkout.

Regenerate the comparison pages after changing the catalog with:

```bash
uv run --with pyyaml python scripts/generate_fontpack_demos.py
```
