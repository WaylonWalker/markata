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
350 milliseconds without edits. Narrow screens stack the editor and preview.
The browser stores your draft locally. Use **Download .md** to save a file.

The example requires a Markata wheel containing the Python API from PR #1835.
Until that API is released, install a wheel built from that PR:

```bash
uv run --with /absolute/path/to/markata-VERSION-PLATFORM.whl examples/fastapi_editor.py
```

Replace the wheel path with your local wheel. With Python and dependencies already
installed, `python examples/fastapi_editor.py` also works. The server listens on
localhost, port 8765. Press Ctrl+C to stop it.

Each preview calls `markata.render_async()` through FastAPI. Inline assets support
code highlighting, Mermaid, and Chart.js without separate asset routes. A fresh
iframe document executes initialization scripts after each successful render.
The iframe permits scripts and keeps the preview separate from the editor origin.
Authored HTML and inline styles are enabled for this local authoring tool.
The preview remains isolated from the editor origin.

The example limits input to 200,000 characters and rendering to 15 seconds.
Errors preserve the last successful preview. Each call starts a renderer process.
This example is a local development tool. It does not include accounts or shared storage.

Use the appearance button beside **Download .md** to open Markata's full picker.
Search color themes, switch light and dark modes, or select Style and Font tabs.
Selections apply to the preview and persist across edits. Markata supplies the
palette catalog, styles, scripts, and fonts. No site build runs during editing.

The footer shows Markata time, request time, input bytes, document bytes, and frame
load time. Markata time includes process startup, rendering, asset collection, and
JSON parsing. Request time also includes HTTP transfer and browser JSON parsing.
Frame load time measures iframe loading, not completed asynchronous diagrams.
The 350-millisecond typing delay is separate from these measurements.

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
