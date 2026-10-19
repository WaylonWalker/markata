# Font packs

Markata uses `fontpack` to select typography. The default is `system`, a
zero-download pack that emits semantic `--font-*` variables but no font
binary or `@font-face` rule. Built-in catalog data, manifests, licenses, and
font files are embedded in the Markata executable; they never require a
source checkout or a catalog in the site's working directory.

Bundled packs use prebuilt WOFF2 files and stable catalog tiers. Builds copy
only referenced tiers into `output/assets/fonts` and write one shared
`output/css/fonts.css`; page content never creates a new subset. The optional
FontTools workflow is intentionally separate from ordinary builds.

Before the stylesheets, HTML heads preload only the emitted WOFF2 tier used
for the active pack's body and first-level heading (display, or heading when
display is absent). Duplicate sources are preloaded once; system stacks need
no preload. A page's `fontpack` override replaces the site default preload
list. When the picker is enabled, the synchronous theme bootstrap may preload
the same files for a stored `theme-fontpack` choice instead, without preloading
every pack on every visit. Preloads use `as="font"`, `type="font/woff2"` and
`crossorigin` so the CSS fetch can reuse them. The shared stylesheet has a
content-hashed URL and retains `font-display: swap`.
Nonempty editorial feeds MUST also preload the active pack's code face, because
their first visible week label uses the monospace role. The preload MUST use the
same emitted tier as CSS, deduplicate sources shared with body/headings, respect
a stored picker choice, and add no font hint for system-stack code roles.
JavaScript-enabled editorial calendars MUST reserve the view switcher's space
before its deferred controller initializes, preventing the earlier first paint
from exposing a pagination shift. Pending controls remain invisible and
non-interactive; no-JavaScript layout and the final toolbar remain unchanged.
Empty editorial feeds, classic/simple feeds, and article pages retain their
existing preload selection. Code-face preload metadata MUST be cached and
validated with the ordinary preload metadata; changing it invalidates old HTML.

Built-in packs cache the hash and per-pack preload URLs alongside the font
output, keyed by the rendered content and catalog; a missing or invalid
cache entry triggers ordinary resolution instead of serving stale hints.

## Resolution cost and cache identity

A build MUST compute a canonical set of distinct visible runes once and reuse
it for the coverage signature and all requested packs. The extraction MUST
ignore script/style contents, decode HTML entities, ignore replacement runes,
normalize HTML tabs, line breaks, form feeds, and carriage returns to U+0020 as
browsers do when collapsing whitespace, and preserve concatenated-article
parsing semantics. Non-breaking spaces and supported Unicode text remain
distinct coverage runes.
It SHOULD accept a reader so article HTML need not be copied into a site-sized
intermediate string. Reader errors MUST be surfaced.

Tier selection MUST preserve existing source-aware fallback and missing-tier
errors. Once a source requires `full`, further coverage checks for that source
MUST stop because full supersedes smaller tiers. Subset profiles SHOULD be
compiled once per resolution, and source manifests and validated asset metadata
SHOULD be reused within that resolution, not across mutable custom-catalog builds.
Multi-pack resolution MUST not repeat visible-text extraction or discard
redundantly generated per-pack CSS.

Built-in cache identity MUST include the default resolved pack separately from
the sorted, deduplicated set of effective resolved packs, picker enablement,
visible coverage, the catalog, and a digest of embedded catalog/lock/manifest
metadata. Pack declaration order and equivalent aliases MUST not cause misses;
default/picker changes and bundled-manifest revisions MUST invalidate old output.
Cache schema changes MUST invalidate older cache records.

The intact-output fast path MUST continue to require the stylesheet, managed
font files, and valid preload metadata. Missing output MUST be regenerated with
the same accelerated resolution; a cache hit MUST NOT bypass custom-catalog
checksum or role-capability validation. No change to typography, picker choices,
preload selection, font binaries, or CSS is implied by these optimizations.

## Page annotation and warm builds

The fontpack writer MUST ensure that a complete HTML document has exactly one
`data-fontpack` attribute on its opening HTML element, using the resolved
per-page pack. It MUST preserve unrelated attributes, quoted values, and
whitespace. Attribute and tag names are case-insensitive; fontpack values are
case-sensitive. Attribute values MUST be HTML-escaped when written.

If the opening HTML element already selects the resolved pack, annotation MUST
reuse the existing document unchanged. Finding and checking that element MUST
NOT lowercase, copy, or otherwise process the entire page body. A page without
an HTML element remains unchanged by root-element annotation. Malformed
unfinished tags MUST NOT be partially rewritten.

The font stylesheet remains shared. Annotation MUST NOT add an unversioned
stylesheet when the page already references the hashed stylesheet. These
shortcuts do not bypass catalog, visible-glyph coverage, preload, or missing
output-file validation.

Custom catalogs may be selected with `fontpacks_file`. Relative paths in a
custom catalog resolve relative to the catalog file. Markata records generated
font filenames in `output/assets/fonts/.markata-fonts.json` and removes only
stale files listed by that manifest on later builds.

Config validation MUST accept a `fontpack` or `theme.fontpack` value that
resolves to a pack or alias in the built-in font catalog, even when the name
is not listed in the rendering contract's `fontpacks` enum. When
`fontpacks_file` is set, any name is accepted and the fontpack plugin reports
unknown packs at build time. Other names are rejected with an
`unsupported value` error.

`fonts verify` checks manifest and lockfile provenance, license metadata, full
64-character SHA-256 hashes, and WOFF2 assets. Every bundled manifest must
record a non-empty `source.files` map, and its keys and hashes must exactly
match the lockfile. Short hashes are display-only. With no pack argument it
verifies every unique source used by every built-in bundled pack;
`fonts verify <pack>` limits the check to that pack and uses the configured
catalog when one is present. The built-in
`painted-sign` pack uses Finger Paint for display lettering, Rock Salt for
handwritten accents, Source Sans 3 for prose, and DM Mono for code.

The manifest generator reads every source from the catalog and its locked
metadata; it is not permitted to maintain a separate family list. Generated
tiers with different Unicode ranges must not have identical content. A family
that has no additional glyphs for a tier omits that tier rather than shipping a
duplicate asset.

Tier selection is source-aware. `latin-ext` is requested only when the selected
family manifest provides it; otherwise the resolver selects that family's
`full` tier. Unsupported scripts likewise select `full`. A missing required
tier still fails verification when the family has no `full` fallback.

Role `weight`, `style`, `size`, and `optical_size` are functional. They are
emitted as role variables and applied by the generated stylesheet only when
the corresponding role property is configured. If a role is absent, its
entire rule falls back to `body`; if it exists, omitted optional properties
remain omitted rather than inheriting from `body`. This prevents a specialized
role from accidentally receiving body typography such as optical sizing. In a
multi-pack stylesheet, both role variables and role rules are scoped by
`data-fontpack`, so an optional property configured by one pack cannot affect
another pack. The `display` role is applied to `h1`; the `heading` role is
applied to `h2` through `h6`, with `display` falling back to `heading` when it
is absent. Variable-font manifests record the actual `fvar` axis bounds
discovered from the source font; the generator never assumes a `[300, 900]`
weight range. A non-zero `optical_size` requires an `opsz` axis whose recorded
bounds contain the requested value. `fonts verify` rejects role weights,
styles, and optical sizes that no bundled face can satisfy.
