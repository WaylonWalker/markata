---
title: "Live Markdown editor with FastAPI"
description: "Run a single-file editor with a live Markata preview"
date: 2026-10-28
published: true
tags:
  - documentation
  - python
  - markdown
---

Run the single-file example:

```bash
uv run examples/fastapi_editor.py
```

Open <http://127.0.0.1:8765>. Edit Markdown on the left. The preview updates after
250 milliseconds without edits. Narrow screens stack the editor and preview.
The browser stores your draft locally. Use **Download .md** to save a file.

The example requires a Markata wheel containing the Python API from PR #1835.
Until that API is released, install a wheel built from that PR:

```bash
uv run --with /absolute/path/to/markata-VERSION-PLATFORM.whl examples/fastapi_editor.py
```

Replace the wheel path with your local wheel. With Python and dependencies already
installed, `python examples/fastapi_editor.py` also works. The server listens on
localhost, port 8765. Press Ctrl+C to stop it.

Each preview uses an application-owned `RendererPool` through FastAPI. Two warm
renderers start with the application's lifespan and retain the workspace's
Markdown configuration. Inline assets support
code highlighting, Mermaid, and Chart.js without separate asset routes. A fresh
iframe document executes initialization scripts after each successful render.
The iframe permits scripts and keeps the preview separate from the editor origin.
Authored HTML and inline styles are enabled for this local authoring tool.
The preview remains isolated from the editor origin.

The example limits input to 200,000 UTF-8 bytes and rendering to 15 seconds.
Errors preserve the last successful preview. A crashed or timed-out renderer is
replaced automatically; no new subprocess is started for ordinary requests.
This example is a local development tool. It does not include accounts or shared storage.

Use the appearance button beside **Download .md** to open Markata's full picker.
Search color themes, switch light and dark modes, or select Style and Font tabs.
Selections apply to the preview and persist across edits. Markata supplies the
palette catalog, styles, scripts, and fonts. No site build runs during editing.

The footer shows Markata time, request time, input bytes, document bytes, and frame
load time. Markata time includes queue wait, rendering, asset collection, and JSON parsing,
but normally excludes process startup. Request time also includes HTTP transfer and browser JSON parsing.
Frame load time measures iframe loading, not completed asynchronous diagrams.
The 250-millisecond typing delay is separate from these measurements.

The preview uses the default site's article wrappers, reading size, and spacing.
It previews article content. It does not add site navigation, feeds, or post metadata.

Lucide shortcodes such as `:lucide-globe:` and `:lucide/heart:` display as inline
icons in prose and headings. Markata bundles all 2,118 icons from Markata's pinned
Lucide pack, including its license. Code examples and unknown shortcodes stay unchanged.

The Markdown pane uses bundled CodeMirror with syntax highlighting, line numbers,
line wrapping, and bracket matching. Select **Keymap** above the editor to use
Default, Vim, Emacs, or Sublime bindings. Your choice persists in this browser.
With Vim selected, press `i` to insert text and `Escape` to return to normal mode.
Editor colors follow the appearance picker. CodeMirror requires no CDN requests.


## Run the demo integration tests

Use a virtual environment containing the built Markata wheel, FastAPI, Uvicorn,
HTTPX, and Playwright. Run the endpoint checks against that installed wheel:

```bash
python -m unittest discover -s python/integration -v
```

Install Chromium and run the browser checks:

```bash
python -m playwright install chromium
python python/test_editor_browser.py -v
```

If Chromium is already installed, set `MARKATA_TEST_CHROMIUM` to its executable path.
The tests reserve an available localhost port and stop their server afterward.
They cover live edits, saved drafts, keymaps, initial appearance, picker changes,
Lucide icons, responsive rails, and full-width desktop and mobile layouts.
The wheel workflow runs both suites against the installed Linux wheel.

The first render returns the full preview document. Later requests set
`include_document` to `false`, so the server returns only the rendered fragment,
conditional styles, scripts, and metrics. The browser keeps its document and
appearance styles loaded, then replaces the article content in the existing
sandboxed frame. This avoids sending the large site CSS payload or reloading
theme stylesheets after every edit.

See [integration benchmarks](/docs/guides/renderer-integration-benchmarks/) for measured fixture timings and reproduction steps.

## Shared editor behavior

Use **Open .md** to import a Markdown file. An oversized file leaves the current draft intact.
Use **Download .md** to keep a copy outside browser storage.
The footer reports whether the draft was saved. Storage failures leave editing and downloads available.
See [shared editor features and runtime differences](/docs/guides/renderer-integration-benchmarks/#shared-editor-features-and-runtime-differences) for the comparison.

## Insert Markdown commands

Select **Commands** or press **Ctrl+P** (**Cmd+P** on macOS).
Search for special syntax, then press **Enter** to insert it at your selection.
See the [command palette guide](editor-command-palette.md) for the catalog and keyboard controls.

## Language assistance

The Python and Go editors connect to Markata's LSP through a same-origin WebSocket.
The server selects the workspace with `MARKATA_LSP_WORKSPACE`; browser messages cannot change it.
The endpoint accepts `localhost` and loopback IP hosts only. It is intended for local examples.

Both editors complete wiki links and Markdown document paths. They also complete
command-catalog containers, their known attributes, and Mermaid diagram syntax.
Press **Ctrl+Space** to request suggestions at the cursor. File-tree navigation
is a separate feature.

To run the browser checks against both servers, build the Go binaries and set
`MARKATA_LSP_BINARY` and `MARKATA_WORKBENCH_BINARY` as shown below:

```bash
go build -o /tmp/markata-editor-lsp ./cmd/markata-go
go build -o /tmp/markata-editor-workbench ./cmd/markata-workbench
MARKATA_LSP_BINARY=/tmp/markata-editor-lsp \
MARKATA_WORKBENCH_BINARY=/tmp/markata-editor-workbench \
python python/test_browser_lsp.py -v
```

## Workspace files

Set `MARKATA_LSP_WORKSPACE` to a cloned site before starting the server.
The Files tree uses Markata glob and load plugins through the bundled Go binary.
Select a file to edit a browser draft; download it to keep a copy. See
[workspace navigation](./renderer-site-tree/).
