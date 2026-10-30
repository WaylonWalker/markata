---
title: "Renderer profiles and cross-application assets"
description: "Choose rendering and asset defaults for Python, CLI and Go workbench and WASM adapters"
date: 2026-10-29
published: true
tags:
  - markdown
  - rendering
  - configuration
---

# Renderer profiles and cross-application assets

Markata's standalone renderer is reusable from scripts and applications without
building a site. Profiles are **opt-in defaults**, not a replacement for full
site config or the site-build lifecycle.

| Profile | Native Go / CLI / Python | Limitations |
| --- | --- | --- |
| `web` (default) | Existing feature-sensitive asset manifest, linked CSS/JS unless inlined explicitly | Authored URLs may load remote resources |
| `offline-interactive` | Inline managed CSS/JS; reject unsafe HTML and detected remote media/styles | Not yet a fully verified offline browser build |
| `offline-static` | Explicit error | Static Mermaid SVG and Chart.js export not implemented |

The result includes article HTML, feature signals, an **ordered** asset manifest,
and generated tags for the head and body. Shared libraries may be cached by ID;
per-document initializer scripts must not be deduplicated across independent
documents.

## In the `markata serve` settings sidebar

Open **Settings** (Alt+,), then look under **Common** or **Markdown → Render profile**.
This uses the existing Preview and Bake flow, so the chosen value can be written
to the config file that owns it.

```toml
[markata-go.markdown]
render_profile = "offline-interactive"
render_asset_mode = "inline" # linked or inline; profile may force inline
render_asset_base = "/markata-assets" # used for linked assets
```

The **asset mode** and **asset base** are independent of the regular site\'s `assets.mode` and `assets.output_dir`; they configure standalone renderer integration output only. A configured profile and these asset settings are inherited by `--site-render-profile` (or Python `config_path=`). Explicit CLI overrides still take precedence.\n\nThe site build itself still uses its configured Mermaid, link-avatar, theme and
asset plugins. This new setting is for **standalone renderer consumers that opt
in to loading site config**; it does not yet guarantee that a complete static
site is offline. See the full config and plugin settings for site behavior.

## Apply a profile

```sh
printf '# Example\n' | markata render --json --profile offline-interactive
```

Or read the profile saved in site config:

```sh
printf '# Example\n' | markata --config markata-go.toml render --json --site-render-profile
```

```python
from markata import render

result = render("# Example", profile="offline-interactive")
site_result = render("# Example", config_path="markata-go.toml")
```

A supplied CLI `--profile` wins over `--site-render-profile`.
`render()` and `render_async()` use the same binary and profile policy.
No site config is loaded by default.

## Adapter parity checkpoints

The Go workbench and browser/WASM workbench currently share the same underlying
Markdown renderer, but **their asset transports differ**:

| Adapter | Default assets | Current authoring model |
| --- | --- | --- |
| Native Go CLI | Bundled manifest, linked or inline | Direct Go call |
| Python/FastAPI | Manifest and optional inlined tags | One Go subprocess per render |
| Go workbench | Local linked assets + persistent preview, with `profile`/`inline_assets` request options and ordered asset metadata | In-process render |
| Browser WASM | Linked pre-cached assets by default, optional fully inline assets and shared profiles | Persistent Go WASM worker |

**Follow-up:** complete cross-adapter golden fixtures and capability discovery.
Go and WASM accept `profile` and `inline_assets` without a second Markdown parser.
WASM reports unsupported native-only static rendering modes rather than silently
falling back to client-side Mermaid.

Strict browser-network guarantees, asset font audits and static diagram/chart
prerendering remain tracked in
[#1854](https://github.com/WaylonWalker/markata/issues/1854).
