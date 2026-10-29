# Standalone presentation components

The standalone renderer expands bundled Lucide shortcodes through the existing
icon plugin without network access or filesystem writes. Authored HTML remains
opt-in. Unknown icons and code examples remain unchanged.

`markata render --appearance` emits versioned JSON containing picker markup,
ordered styles/scripts, and a map of reusable assets. Asset bytes use JSON base64
encoding. Generation uses the same embedded templates and palette, aesthetic,
and font catalogs as a site build. It performs no site discovery or writes.

Python `appearance()` loads this bundle once at application startup. It exposes
`picker_html`, `head_html`, `body_end_html`, and `assets` (path to bytes).
`head_html` and `body_end_html` link to a configurable asset base. Host applications
serve assets by exact lookup; missing assets return 404. Generated font URLs use
the same base. Appearance defaults and overrides are validated against the catalogs.

`RenderResult.document(appearance=..., palette=..., mode=..., aesthetic=...,
fontpack=..., text_size=...)` uses site article wrappers and selected settings.
Without appearance it still provides the site article hierarchy and defaults.
Optional `appearance_bridge=True` supports parent-frame appearance updates using
validated attributes and catalog values. It also accepts `markdown-render`
messages from the parent to replace `.post-content`, run the supplied generated
body scripts, update render-specific head styles, and dispatch the view-transition
completion event used by site components. It does not replay `DOMContentLoaded`.
Hosts can keep the initial preview document and its shared stylesheets loaded
between renders. No bridge is installed by default.

The FastAPI editor owns CodeMirror, keymaps, editing, debounce, persistence, and
metrics. It uses these APIs without modifying generated HTML or duplicating icons,
picker templates, fonts, palettes, or site layout rules.

The editor requests the full preview document on its initial render. Later render
requests can omit that shell and return only the article fragment, conditional
styles, scripts, and metrics. Reusing the sandboxed frame avoids retransmitting
shared appearance styles with every edit.

Document color mode and selected appearance are applied in the initial HTML.
An early inline background and native color scheme prevent a white first paint
while linked theme styles load. The editor sends current appearance with each
render request and initializes its first preview after the picker bootstrap.

The host editor explicitly owns its full-width shell. Site main-column styles
must not constrain the workbench; preview documents retain the site's layout.
Responsive Markdown containers use the preview viewport width.

`Appearance.controls_head_html` supplies picker styles without site main-column or
article component rules. Host applications use this entry point; `document()` uses
the complete appearance stylesheet set inside its own document.

Integration benchmarks reuse the deterministic site benchmark post and rendering
fixtures. Reports separate process-backed rendering, server rendering, HTTP wall
time, and response size. They record sample counts, fixture hashes, environment,
and source revision. Initial document requests and fragment updates use separate
measurements. Timing results describe the measured host, not a performance limit.
