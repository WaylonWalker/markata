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

- **Persistent JSON Lines / stdio:** implemented with one request per line, one correlated response per line, bounded input, and per-request error recovery. Stdout is reserved for protocol data.
- **HTTP:** use identical render semantics with explicit resource limits, authentication/CORS strategy and safe defaults. #1851 is a *local trusted editor*, not a public rendering service.
- **LSP:** preserve existing editor features and translate rendering diagnostics to editor ranges where available. Generic LSP does not define an HTML-preview window.
- **MCP:** expose small render/validate/capabilities tools reusing the engine; reference or deduplicate large assets rather than returning megabytes repeatedly.
- **WASM:** report its runtime-specific capabilities rather than claiming native prerendering features it does not support.

Config presets, strict offline verification, favicon controls, Mermaid SVG prerendering and static charts belong in [#1854](https://github.com/WaylonWalker/markata/issues/1854); this contract must be able to represent them without implying they already exist.

## Compatibility and verification gates

- Run the capability fixture and standalone render tests when changing capabilities or the result schema.
- Add golden fixtures for successful Markdown, optional assets, malformed inputs and unsupported options when request envelopes are implemented.
- Compare normalized render output across native Go, Python, HTTP and WASM where supported. Document deliberate differences; don't normalize away asset order or security-sensitive output.
- Never increment `schema_version` merely for adding a capability option. Only change it for a genuine incompatible result structure.
