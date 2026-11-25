---
title: "Draw in the Markata editor with Excalidraw"
description: "Create and edit Excalidraw scenes inside Markdown in the Markata workbench."
date: 2026-11-16
published: true
tags:
  - documentation
  - editor
  - diagrams
---

The Go and FastAPI workbenches can open Excalidraw inside the Markdown editor.
The drawing is saved as JSON in an `excalidraw` code fence in the current page.

## Create a drawing

1. Open a site in the Go workbench:

   ```sh
   markata editor --workspace ./my-site
   ```

2. Place the cursor where the drawing should appear.
3. Select **Drawing** above the source, or open **Command Palette** and choose
   **Insert Excalidraw drawing**.
4. Draw on the canvas. Select **Apply drawing** to add it to the Markdown.

Cancel or press Escape to close the canvas without changing the Markdown.

## Edit a drawing

Complete Excalidraw fences start with the JSON folded to keep the source editor
compact. Select **Show JSON** to inspect or edit the scene source, then select
**Hide JSON** to fold it again. You can reach both controls and **Edit drawing**
with Tab, then activate them with Enter or Space.

Select **Edit drawing** on the scene controls to open that drawing, even when
the cursor is elsewhere. You can also place the cursor inside the fence and
select **Edit drawing** above the source, or choose **Edit Excalidraw drawing**
from the command palette. Select **Apply drawing** to save the scene. The
editor keeps the opening and closing Markdown fence and updates the JSON
between them. Read-only documents cannot open a drawing for editing.

````markdown
```excalidraw
{"type":"excalidraw","version":2,"source":"https://excalidraw.com","elements":[],"appState":{"viewBackgroundColor":"#ffffff","gridSize":null},"files":{}}
```
````

The scene stays in the Markdown file. Markata does not send it to Excalidraw's
service. The editor loads its drawing code and bundled fonts from local Markata
assets. The browser can use its normal font fallback for other characters.
Font files are available offline through the local asset path used by the
Excalidraw editor.

The workbench can store documents up to 200,000 bytes. If applying a drawing
would exceed this limit, Markata keeps the source unchanged and shows an error.

## Preview and site output

The Go workbench can attach a compressed SVG export to the same scene block
when you apply it. The export is optional: Markata still applies and saves the
editable JSON if SVG export fails or does not fit inside the document limit.
Published pages use a content-addressed SVG file. They do not load the
Excalidraw editor or its JavaScript.

The SVG follows the site's configured fallback mode when Markata builds the
page. Dark mode is the default when no site configuration is available. The
editor canvas UI follows the active Markata mode. New drawings use that mode's
canvas background; existing drawings keep their saved background while open
for editing.

The default `auto` mode reuses an editor-exported SVG only when its scene,
renderer assets, fonts, theme, and export settings still match. It never
starts Chromium. If an artifact is missing or stale, the build reports an
error. Apply the scene again in the Go workbench to refresh its artifact, or
explicitly select Chromium compatibility mode in `markata-go.toml`:

```toml
[markata-go.excalidraw]
enabled = true
mode = "chromium"

[markata-go.excalidraw.chromium]
browser_path = "/usr/bin/chromium"
timeout = 30
no_sandbox = false
```

`mode = "cached"` also requires a matching editor artifact and never starts a
renderer. The current build does not yet advertise Node or Rust modes; see the
[renderer comparison report](/docs/reports/excalidraw-renderer-comparison/)
for measured options and known compatibility limits.

Set `no_sandbox = true` only when the build runs in a container that requires
it. If the renderer cannot find a browser, or a scene contains an unsupported
feature, the build reports an error instead of publishing a code block as if it
were a rendered drawing. The first static renderer supports shapes, lines,
arrows, freehand strokes, and text. Embedded images, links, and web embeds are
not supported yet.

Set `enabled = false` to keep Excalidraw fences as highlighted JSON code. The
standalone renderer returns an inline SVG data URI and does not write an asset
file.
