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

Place the cursor inside its `excalidraw` code fence, then select **Edit drawing**.
You can also choose **Edit Excalidraw drawing** from the command palette.
Select **Apply drawing** to save the scene. The editor keeps the opening and
closing Markdown fence and updates the JSON between them.

````markdown
```excalidraw
{"type":"excalidraw","version":2,"source":"https://excalidraw.com","elements":[],"appState":{"viewBackgroundColor":"#ffffff","gridSize":null},"files":{}}
```
````

The scene stays in the Markdown file. Markata does not send it to Excalidraw's
service. The editor loads its drawing code and bundled fonts from local Markata
assets. The browser can use its normal font fallback for other characters.

The workbench can store documents up to 200,000 bytes. If applying a drawing
would exceed this limit, Markata keeps the source unchanged and shows an error.

## Preview and site output

The editor preview and generated site render the scene as static SVG. Published
pages use a content-addressed SVG file. They do not load the Excalidraw editor
or its JavaScript.

The native Go renderer uses the local Excalidraw export bundle in headless
Chrome or Chromium. A browser is needed only when a site contains an enabled
Excalidraw fence. Install Chrome or Chromium, or set its path in
`markata-go.toml`:

```toml
[markata-go.excalidraw]
enabled = true

[markata-go.excalidraw.chromium]
browser_path = "/usr/bin/chromium"
timeout = 30
no_sandbox = false
```

Set `no_sandbox = true` only when the build runs in a container that requires
it. If the renderer cannot find a browser, or a scene contains an unsupported
feature, the build reports an error instead of publishing a code block as if it
were a rendered drawing. The first static renderer supports shapes, lines,
arrows, freehand strokes, and text. Embedded images, links, and web embeds are
not supported yet.

Set `enabled = false` to keep Excalidraw fences as highlighted JSON code. The
standalone renderer returns an inline SVG data URI and does not write an asset
file.
