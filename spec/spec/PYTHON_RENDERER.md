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
