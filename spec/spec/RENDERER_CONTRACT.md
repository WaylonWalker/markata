# Renderer interoperability contract

**Status:** incremental contract definition, tracked in [#1855](https://github.com/WaylonWalker/markata/issues/1855).
**Baseline:** the standalone renderer introduced by [#1835](https://github.com/WaylonWalker/markata/pull/1835).
**Related:** interoperability roadmap [#1856](https://github.com/WaylonWalker/markata/issues/1856) and configuration/offline policies [#1854](https://github.com/WaylonWalker/markata/issues/1854).

## Principles

1. There is **one Go Markdown renderer**. CLI, Python, HTTP, WASM, LSP, and MCP are adapters, not alternate Markdown implementations.
2. Preserve existing schema-v1 JSON and Python clients. Additive fields may appear only when older consumers can ignore them; incompatible changes require a new version.
3. Preserve ordering of generated assets and distinguish reusable libraries from per-document initializer scripts.
4. Never mistake "bundled libraries" for a complete offline/remote-dependency guarantee.
5. Report unavailable capabilities explicitly. Do not silently fall back from static Mermaid rendering to client-side JS if a caller requires static output.
6. Do not load the site build lifecycle implicitly as a side effect of rendering a fragment.

## Implemented: standalone result v1

`markata render --json` reads a single Markdown document from stdin and emits one JSON object. `schema_version` equals **1**. Existing fields remain:

| Field | Meaning |
| --- | --- |
| `schema_version` | Compatibility version for the render result |
| `html` | Article fragment, excluding renderer-generated scripts |
| `features` | Feature-name to boolean signal |
| `assets` | Ordered styles/scripts with `id`, `kind`, `placement`, optional `content`, `url`, and `module` |
| `external_urls` | External browser dependencies explicitly detected by the renderer; **not** an audit of authored remote images/links |
| `head_html` | Generated ordered head tags for styles |
| `body_end_html` | Generated ordered body-end tags for scripts |

`render` without `--json` emits the HTML fragment. `--inline-assets`, `--asset-base`, and `--unsafe` remain the existing options. The base standalone renderer does not load site configuration; the CLI can apply a selected workspace with `--workspace`. Unsafe HTML remains opt-in and the output is **not** a sanitizer for untrusted content.

## Implemented: capability discovery (first slice)

```shell
markata render --capabilities
```

This prints a JSON description without reading document stdin or rendering Markdown. The reference fixture is `pkg/plugins/testdata/renderer_contract/capabilities.json`.
`--capabilities` describes the base standalone renderer and cannot be combined with `--appearance` or `--workspace`.

- `schema_version`: version of the **render result** supported by this engine.
- `renderer`: identifies the standalone rendering implementation, not a transport.
- `formats`, `options`, `features`, `diagram_modes`, `asset_modes`: currently supported capabilities; values are conservative and deterministic.
- `loads_site_config`: false until configuration is explicitly supported in standalone rendering.
- `external_url_audit`: false until authored and generated resource references are fully audited.

The current standalone engine supports **client-mode Mermaid and Chart.js**. The *site builder* has additional Mermaid render modes; those are not advertised here because this entry point does not expose them.

Clients should check capabilities rather than assuming all Markata site features are available. Unknown future fields should be ignored. New capability names may be added in compatible releases.

## Implemented next slice: transport-neutral request envelope

The JSON Lines stdio transport uses a request envelope resembling:

```json
{
  "id": "doc-42",
  "op": "render",
  "source": "# Heading",
  "options": {"inline_assets": false}
}
```

Responses include the echoed `id`, a typed success or error, and the existing render result. Implemented operations are `render`, `capabilities`, and `shutdown`. Error responses have stable machine-readable codes and readable messages. A rendering failure does not kill the persistent process. See [RENDERER_STDIO.md](RENDERER_STDIO.md) for exact framing, supported fields, size limits, and failure behavior. These operations are available through `markata render --stdio`; `markata render --json` remains a one-shot command.

## Future transports and security boundaries

The Go workbench may add workspace-aware authoring transforms around the
standalone renderer. These transforms are editor conveniences and MUST remain
separate from `markata render`'s site-independent contract. In particular,
resolved wikilinks in a workbench preview use the selected site's LSP index;
they do not imply that standalone rendering now loads a complete site.

The local Go workbench MUST save edits to an existing selected workspace file
automatically. Writes MUST stay inside the selected workspace, replace files
atomically, and reject oversized content. Unsaved scratch documents remain in
browser storage until the author uses an explicit export action. The browser
may retain a recovery copy, but it MUST report whether the site file was saved.
The workbench MUST provide a resizable editor/preview split and editor-only,
preview-only, and split layouts. Internal wikilinks in the preview MUST open
the matching workspace document in the editor.
The Go workbench MUST offer built-in Markdown starters and local Markdown
starters from the workspace `templates/` directory. Creating content MUST
require a new Markdown path inside an existing workspace directory and MUST
not overwrite an existing file.
The New Post command MUST open a fuzzy template picker before the details form.
After a template is selected, the details form MUST start with the template's
title and a default destination folder. The slug field MUST remain empty when
using the suggested slug; its placeholder and the displayed destination path
MUST update as the title changes. If the author enters a slug, that value MUST
take precedence over the title-derived suggestion.
Feed navigation MUST use feed names as tree roots and paths relative to the
workspace root beneath them; absolute filesystem paths MUST NOT appear as feed
labels.
In split layout, source and preview scrolling MUST synchronize in both directions
by default. The active pane controls the position; editing or moving the source
cursor makes the editor the active pane. Users MUST be able to disable scrolling
sync in editor settings, and the setting MUST persist locally. Source line
anchors MUST refer to the original editor document, including any stripped
frontmatter, and MUST be strictly increasing before interpolation. Editor-to-
preview anchors remain in source order, so preview coordinates may decrease for
side-by-side content. Preview-to-editor anchors use visual order; equal pixel
positions select the source line nearest the last synchronized line, with the
greatest source line breaking ties. Without a prior source line, they select the
greatest source line at that position. Preview anchors
MUST use rendered article document coordinates and refresh after content,
viewport, or observed article geometry changes. Scroll synchronization MUST use
the same measured source position in either direction. The editor MUST keep its
native caret at the actual text cursor; a scroll-position marker MUST NOT be
rendered as an editor cursor or caret. Content-only preview updates MUST preserve
the iframe document and unchanged rendered block nodes, and MUST NOT dispatch
global navigation or theme-change events. New or changed feature blocks MUST be
initialized without rerendering unchanged diagrams or components. Preview
position updates MUST be deduplicated when the measured target has not changed.
Only explicit pane input (wheel, touch, scrollbar or pointer drag, and scroll
keys) may make a pane the scroll-sync driver. Scroll changes caused by rendering,
image or font layout, browser scroll anchoring, or synchronization itself MUST
not be treated as user scrolling. Text edits MUST NOT move the preview before
the corresponding content patch is ready. After the patch, a focused editor's
actual caret line MUST remain the synchronization target even when editing did
not change its logical line; the preview marker MUST be resolved again against
the updated source map.
The workbench MAY add start-line and end-line attributes to container openers in
the render input. This instrumentation MUST preserve authored editor text and
line count, ignore fenced code, and use the attributes as additional nested
preview anchors. When several nested containers contain the caret line, the
preview marker MUST select the narrowest source span, then the deepest element.
Preview-only and editor-only layouts do not synchronize hidden panes.

The Go workbench MUST bind Ctrl/Cmd+P to a fuzzy workspace file finder that
searches Markdown paths, titles, slugs, and tags and opens the selected file in
the editor. The Markdown command palette MUST be labeled "Command Palette" and
MUST open from Ctrl/Cmd+Shift+P or Ctrl/Cmd+K. Its visible Commands button MUST
remain available.
Workspace previews MUST resolve internal post embeds written as `![[slug]]` or
`![embed](path/to/post.md)` into linked embed cards using indexed post metadata.
Embed syntax inside fenced code blocks MUST remain unchanged.

- **Persistent JSON Lines / stdio:** implemented with one request per line, one correlated response per line, bounded input, and per-request error recovery. Stdout is reserved for protocol data.
- **HTTP:** use identical render semantics with explicit resource limits, authentication/CORS strategy and safe defaults. #1851 is a *local trusted editor*, not a public rendering service.
- **LSP:** preserve existing editor features and translate rendering diagnostics to editor ranges where available. Generic LSP does not define an HTML-preview window.
- **MCP:** expose small render/validate/capabilities tools reusing the engine; reference or deduplicate large assets rather than returning megabytes repeatedly.
- **WASM:** report its runtime-specific capabilities rather than claiming native prerendering features it does not support.

A separate internal [static SVG foundation](STATIC_VISUALIZATIONS.md) defines backend SVG identities, validation, and publication, but is not yet integrated with this standalone result or capability discovery.

Config presets, strict offline verification, favicon controls, Mermaid SVG prerendering and static charts belong in [#1854](https://github.com/WaylonWalker/markata/issues/1854); this contract must be able to represent them without implying they already exist.

## Compatibility and verification gates

- Run the capability fixture and standalone render tests when changing capabilities or the result schema.
- Add golden fixtures for successful Markdown, optional assets, malformed inputs and unsupported options when request envelopes are implemented.
- Compare normalized render output across native Go, Python, HTTP and WASM where supported. Document deliberate differences; don't normalize away asset order or security-sensitive output.
- Never increment `schema_version` merely for adding a capability option. Only change it for a genuine incompatible result structure.
