# Standalone Markdown rendering

`markata render --json` reads UTF-8 Markdown from stdin, without loading site
configuration or writing output files. The version 1 JSON result contains:
`html`, `features` (boolean signals), ordered `assets`, `external_urls`, and
`head_html` / `body_end_html`. Each asset has a stable `id`, `kind` (style,
script), `content` or `url`, `placement` (head, body_end), and optional `module`.
The fragment contains article markup only; renderer-generated scripts are assets.
Authored raw HTML is disabled unless `--unsafe` is supplied. This does not make
the rich output an HTML sanitizer; applications must sanitize untrusted content.

Core theme styles are included, with conditional code/Chroma, admonition,
Mermaid, road, sortable-table, and interactive-checklist assets. Mermaid and Chart.js use the existing render plugins.
Only rendering plugins run: site-dependent wikilinks, embeds, image processing,
and other build lifecycle features are outside this API. Any external browser dependencies must be explicitly listed in `external_urls`;
rendering never downloads assets. `--inline-assets` embeds bundled styles/scripts into the returned tag
strings. Otherwise `--asset-base` supplies their serving prefix; content is still
returned so applications can serve or cache assets. Generated script ordering
preserves dependency order. Vendored Mermaid and Chart.js libraries are bundled assets too, allowing fully
inline output without third-party script requests.

Official PyPI wheels expose `markata.render` and `markata.render_async`, returning
a result object with the JSON fields as attributes and a `document()`
helper. Python calls the bundled binary over stdin with argv, never a shell.
One process is launched per call; timeouts and cancellation terminate and reap it.
The async interface does not block the ASGI event loop. Protocol mismatch,
missing executables, and renderer errors raise `RenderError`.

## Editor feature parity

The Python, Go, and WASM demos share the Go workbench's system typography,
editor keymaps, appearance controls, Markdown import/export, and visible draft
save feedback. Drafts remain recoverable after reload; storage failures must
leave editing and export available. Each origin owns its own draft/preferences.

All editors reject documents above 200,000 UTF-8 bytes before rendering or
import replacement. Rejected imports preserve the current draft. Successful
imports use the normal edit/save/render path. Errors clear after a valid render.
HTTP endpoints enforce the same byte limit independently of browser validation.
Edits use a 250 ms debounce and coalesce behind one active render; stale results
must not replace newer previews. Shared CSS must not load full site layout
styles into the editor shell. Preview article styles remain available.

Runtime-specific differences remain: Python subprocess versus native Go versus
browser WASM, inline versus linked/cached assets, browser offline support, and
raw HTML policy. Documentation distinguishes those from shared editor features.
