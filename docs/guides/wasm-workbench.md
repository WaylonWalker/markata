---
title: "Browser WASM renderer and local workbench"
description: "Render Markata Markdown in a browser worker, with saved drafts and offline editing."
date: 2026-10-28
published: true
tags:
  - documentation
  - wasm
  - markdown
---

Build and serve the workbench from the repository root:

```bash
python3 scripts/wasm-workbench.py --serve
```

Open <http://127.0.0.1:8766/>. Go and Python 3 are build requirements.
The browser downloads the Go runtime once. No Python package, Node build,
CDN, or rendering API is required. Port 8765 remains available for the
[FastAPI editor](fastapi-editor).

Write Markdown on the left and read its preview on the right. The editor
supports Vim, Emacs, Sublime, and default keymaps. Use the appearance picker
to change the palette, aesthetic, fonts, and text size.

Use **Open .md** to import a document. Use **Download .md** to export it.
The document limit is 200,000 UTF-8 bytes. The footer reports render time
and draft storage status. A storage error does not stop editing or export.

## Drafts and offline use

Drafts and keymap preferences stay in browser storage. Drafts use a separate
key from the FastAPI editor. On first use, the workbench can copy that editor's
draft when both apps share an origin. Different ports have separate storage.

Wait for **Ready offline** before disconnecting. The service worker saves the
app, WASM renderer, editor, theme assets, and fonts. You can then reload the
page and continue editing without a network connection.

Browser storage can be cleared or evicted. Private browsing can restrict it.
Export files to keep durable copies. The demo edits one document at a time;
it does not sync drafts between devices or resolve edits from multiple tabs.

Markdown stays on the device. External images and links in authored content
can still contact their destinations. Raw HTML is disabled. The preview runs
in a sandbox with a separate origin and cannot read the editor's storage.

## Use the worker in another app

Build the static files, then copy the generated directory to your host:

```bash
python3 scripts/wasm-workbench.py --output /tmp/markata-browser
```

The generated `worker.js`, `wasm_exec.js`, and `markata.wasm` form the rendering
entrypoint. Keep these files together. The build copies `wasm_exec.js` from
the Go installation that builds the module.

```javascript
const worker = new Worker('./worker.js');
worker.onmessage = ({data}) => {
  if (data.type === 'ready') {
    worker.postMessage({source: '# Hello from Markata', revision: 1});
  } else if (data.type === 'result') {
    if (data.result.error) console.error(data.result.error);
    else console.log(data.result.html, data.result.assets, data.render_ms);
  } else if (data.type === 'fatal') {
    console.error(data.error);
  }
};
```

The result follows the standalone renderer's schema version 1. It includes
HTML, feature flags, ordered assets, and head/body tags. Bundled assets contain
IDs rather than their file contents. Serve those files from `appearance/`.
Generated initializer scripts retain their contents. The sample parent app loads each dependency once. It sends cached asset contents
to the sandboxed `preview.js`, which disposes old Chart.js instances.

Inside the initialized runtime, `markataRender(source)` returns a JSON string.
Invalid argument types, oversized documents, and render errors return an
`error` field. The worker preserves the request's `revision`. Hosts must discard
stale results. Keep one render in flight and replace pending edits with the
latest source. The sample editor uses a 250 ms debounce.

The renderer uses the same Go Markdown implementation as the Python API.
It supports syntax highlighting, admonitions, containers, bundled Lucide icons,
Mermaid, and Chart.js. Site discovery, frontmatter processing, feeds, and
site-dependent plugins require a full Markata build. Mermaid runs in client mode.

## Static hosting

Serve the output over HTTPS or loopback HTTP. Service workers do not work
from `file://` URLs. Subdirectory hosting is supported; keep relative paths intact.

Serve `.wasm` as `application/wasm`. The loader also supports hosts with an
incorrect MIME type. Enable gzip or Brotli compression. The generated `.gz`
files require `Content-Encoding: gzip` and the original file's content type.
The Python server handles these headers automatically.

The parent app reads assets from its service-worker cache. It embeds preview
fonts and styles so the isolated preview can render offline. The included server
also sets CORS headers for hosts that reuse these assets. Keep `sw.js` available
for revalidation. Each build gives its offline cache a content-derived version.
A new version activates after its complete asset bundle is cached. Reload the
tab to use a newly installed version.

The initial module is larger than a JavaScript-only Markdown parser because
it includes Go, Chroma syntax highlighting, and Markata's extensions. On the
local benchmark build, gzip reduced the module to about 6.5 MiB. Subsequent
renders reuse the runtime and parser. The worker keeps editing responsive.

To serve an existing build:

```bash
python3 scripts/wasm-workbench.py --output /tmp/markata-browser --no-build --serve --port 8766
```

Stop the server with Ctrl+C. To remove an offline copy, unregister its service
worker and clear its site data in browser settings. Export the draft first.

## Browser checks

```bash
python3 scripts/wasm-workbench.py
uv run --with playwright python -m playwright install chromium
uv run --with playwright python tests/wasm/test_workbench.py
```

The checks cover offline reload, import/export, draft recovery, rapid edits,
repeated diagrams and charts, keymaps, mobile width, storage failures, input
limits, worker errors, and local-only network requests. Warm-render measurements
are reported without machine-specific timing assertions.

## Shared editor behavior

Use **Open .md** to import a Markdown file. An oversized file leaves the current draft intact.
Use **Download .md** to keep a copy outside browser storage.
The footer reports whether the draft was saved. Storage failures leave editing and downloads available.
See [shared editor features and runtime differences](/docs/guides/renderer-integration-benchmarks/#shared-editor-features-and-runtime-differences) for the comparison.

## Insert Markdown commands

Select **Commands** or press **Ctrl+Shift+P** (**Cmd+Shift+P** on macOS).
Search for special syntax, then press **Enter** to insert it at your selection.
See the [command palette guide](editor-command-palette.md) for the catalog and keyboard controls.

Open the workbench with `?example` to load a Markdown example and open the command palette.
