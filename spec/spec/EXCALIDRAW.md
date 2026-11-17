# Excalidraw authoring in the workbench

**Status:** initial editor integration; static site rendering is a follow-up.
**Tracked by:** [#2014](https://github.com/WaylonWalker/markata/issues/2014).
**Related:** [Editor workspace](EDITOR_WORKSPACE.md), [static visualization foundation](STATIC_VISUALIZATIONS.md).

## Markdown representation

An Excalidraw drawing MUST be stored in a Markdown fenced code block whose
info string starts with `excalidraw`:

````markdown
```excalidraw
{"type":"excalidraw","version":2,"source":"https://excalidraw.com","elements":[],"appState":{"viewBackgroundColor":"#ffffff","gridSize":null},"files":{}}
```
````

The block body MUST be a JSON object accepted by Excalidraw's local scene
restore API. The block is self-contained. Image data MUST remain in the
scene's `files` data when Excalidraw serializes it. The authoring UI MUST NOT
convert the scene into a remote Excalidraw URL, expose collaboration controls,
enable backend exports, or require an Excalidraw account. Local scene and image
exports may remain available.

## Workbench controls

The Go workbench MUST provide a visible Drawing control and command-palette
actions to insert a drawing or edit the drawing at the current cursor.
The controls MUST be disabled when the active document is read-only.

- If the cursor is inside a valid `excalidraw` fenced block, **Drawing** opens
  that scene in the canvas.
- Otherwise, **Drawing** opens a new blank scene. Applying it inserts a new
  fenced block at the cursor.
- **Insert Excalidraw drawing** always starts a new scene, even when the cursor
  is inside an existing scene.
- **Edit Excalidraw drawing** is available only when the cursor is inside an
  Excalidraw block.
- The editor MUST preserve the opening and closing fence lines and replace only
  the JSON body when an existing scene is applied.
- A malformed scene MUST remain byte-for-byte unchanged. The editor MUST show
  a readable error and MUST NOT open an empty replacement scene over it.
- Apply MUST update the CodeMirror document in one undoable operation. Normal
  draft saving, workspace compare-and-swap saving, diagnostics, and preview
  scheduling MUST continue to use the existing editor change path.
- Cancel, Escape, or closing the dialog MUST discard canvas changes and MUST
  leave Markdown source unchanged. For a new scene, no fence is inserted until
  Apply.
- Before Apply, the complete Markdown document MUST remain within the existing
  200,000-byte editor limit. If it would exceed the limit, Apply MUST fail with
  a clear message and preserve the source.

## Runtime assets

The Excalidraw editor JavaScript, styles, and fonts MUST be served from local
Markata assets. The workbench MUST lazy-load the editor JavaScript when the
author opens a scene. Rendering or editing a scene MUST NOT require a request
to Excalidraw's CDN or a remote account. When a bundled font does not contain a
character, the browser MAY use its normal local font fallback.

The same browser authoring controls MUST work in the Go and FastAPI workbenches,
which share the editor assets. The separate offline WASM workbench is outside
this first integration slice.

## Rendering boundary

This editor integration does not add Excalidraw rendering to the site builder
or standalone Markdown renderer. Until that follow-up is implemented, an
`excalidraw` block is ordinary fenced code in generated pages. A rendering
implementation MUST use the saved scene JSON, publish static output by default,
and report unsupported scene features clearly. It MUST NOT silently inject the
Excalidraw editor runtime into published pages.
