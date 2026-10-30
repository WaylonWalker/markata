# Shared renderer profiles (incremental implementation)

Tracking: [#1854](https://github.com/WaylonWalker/markata/issues/1854).

## Goal and scope

Share **one render profile vocabulary and one Go renderer** across:
- the native `plugins.RenderStandalone` API,
- the one-shot JSON CLI and Python sync/async wrappers,
- Go HTTP and browser/WASM adapters.

`markdown.render_profile`, `markdown.render_asset_mode`, and `markdown.render_asset_base` are **renderer integration defaults**, not a site-builder
plugin lifecycle setting. The existing `markata serve` settings sidebar discovers
this typed setting as part of the Markdown section and pins it under Common.
Changes can be Previewed and Baked into the actual site config (TOML, YAML, JSON).
**Changing the setting does not itself rewrite normal site build asset behavior**;
that is a separate, explicitly scoped task.

## Profiles and capabilities

| Profile | Standalone behavior | Limitation |
| --- | --- | --- |
| `web` (default) | Existing HTML + ordered asset manifest; linked assets by default | May reference remote authored resources |
| `offline-interactive` | Force managed CSS/JS inline; reject `unsafe`, external manifest URLs, remote authored media and stylesheet URLs that the renderer can detect | This is **not yet a certified zero-network browser build**. Vendor script network behavior, custom fonts, plugin assets, and broader URL surfaces still need browser testing |
| `offline-static` | **Reject with clear unsupported error** | Mermaid and Chart.js currently require JS in standalone rendering. Never silently fall back |

The policy is enforced centrally in Go, and the Python wrapper forwards the
same profile string. The JSON response includes optional `render_profile` for
non-default consumers; older readers may ignore the additive key.

## Use

```sh
printf '# Hello\n' | markata render --json --profile offline-interactive
```

```python
from markata import render
result = render("# Hello", profile="offline-interactive")
```

To use the site config saved from the Serve sidebar, explicitly opt in:

```toml
[markata-go.markdown]
render_profile = "offline-interactive"
render_asset_mode = "inline"
render_asset_base = "/markata-assets"
```

```sh
printf '# Hello\n' | markata --config markata-go.toml render --json --site-render-profile
```

```python
result = render("# Hello", config_path="markata-go.toml")
```

Explicit `--profile` overrides `--site-render-profile`. Ordinary standalone
calls continue to render without loading site config or writing site files.

## Asset ownership and next parity work

The JSON result keeps `schema_version=1`, ordered `assets`, `head_html`,
`body_end_html`, and `features`. The profile only changes how those managed
assets are delivered, not the Markdown AST or feature detection. Generated
initializer scripts are per-document and must not be deduplicated like static
libraries.

All three integrations share **the same profile and asset options**:
- Go workbench HTTP: `profile`, `inline_assets`; its asset base remains fixed to the local `/appearance` server. Reuse
  `ResolveStandaloneProfile`, rather than a parallel asset-policy switch.
- WASM: accepts web/offline-interactive, inline vs linked assets, and reports
  native-only unsupported modes explicitly. Service-worker cache and sandboxed
  preview rules remain platform-specific.
- Python: covered by this PR; the CLI and Go implement the actual policy.

#1854 remains open until strict browser-network assertions, Mermaid SVG prerender,
static chart fallback and full cross-runtime parity are proved.
