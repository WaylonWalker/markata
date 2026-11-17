# Excalidraw authoring in the workbench

**Status:** editor integration and static SVG rendering.
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

The site builder and standalone Markdown renderer MUST interpret valid
`excalidraw` fences as saved scene JSON and render static SVG. The native Go
renderer MUST use the locally bundled Excalidraw export implementation in
headless Chromium. It MUST NOT access the network, and published pages MUST
NOT load the Excalidraw editor runtime. The site build MUST publish a
content-addressed SVG asset through `pkg/visualizations`; standalone rendering
MUST return inert, self-contained SVG markup or a data URI without writing
files. The default configuration MUST enable rendering. Setting
`excalidraw.enabled = false` leaves the fence as ordinary highlighted code.

Rendering MUST validate the input as an Excalidraw scene and reject unsupported
or unsafe scene features with a message that identifies the feature. The first
static renderer MUST support Excalidraw shapes, lines, arrows, freedraw strokes,
and text. It MUST report embedded images and linked/embeddable elements as
unsupported instead of dropping them. It MUST enforce the shared 4 MiB SVG
output limit and publish no partially rendered asset on failure. Alternative
text MUST use scene text when present and a descriptive element-count fallback
otherwise.

Native rendering requires an installed Chromium-compatible browser only when
an enabled scene is present. `excalidraw.chromium.browser_path`, `timeout`,
and `no_sandbox` configure that renderer. Missing browser or
invalid scene errors MUST explain how to resolve the problem. The separate
WASM workbench remains outside this rendering implementation.
