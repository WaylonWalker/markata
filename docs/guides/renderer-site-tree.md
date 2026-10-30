---
title: "Editor workspace file tree"
description: "Browse a cloned Markata workspace in the Go and Python editors"
date: 2026-10-30
published: true
tags: [documentation, editor]
---

Start the Go editor for a cloned site:

```bash
markata editor --workspace ~/git/my-site --port 8002
```

For the Python example, set `MARKATA_LSP_WORKSPACE` to the same site directory
before starting `examples/fastapi_editor.py`.

The Files tree lists content discovered by Markata's glob plugin. The editor
reads `markata-go.toml`, YAML, or JSON configuration from the selected workspace.
It uses the configured glob patterns and gitignore behavior. Without a config,
it discovers Markdown under `pages/` and `posts/` with Markata's defaults. The Tags view groups
loaded documents by their frontmatter tags.

Select a file to load its original source, including frontmatter. Edit it and
view the Go-rendered preview. Drafts stay in browser storage, scoped to the
workspace and path. Switching files keeps each draft. **Download .md** exports
the active draft. This iteration does not write changes to the cloned files.
The scratch `draft.md` buffer remains available for imported documents.

Only discovered files inside the selected workspace can be read. Files reached
through symlinks outside the workspace are excluded. Source reads are limited
to 200,000 bytes; a failed read leaves the current editor buffer intact.

## API usage

Go applications can use `editorworkspace.New(root)`, `Files()`, and `Read(path)`
from `github.com/WaylonWalker/markata/pkg/editorworkspace`. Discovery uses the
existing glob and load plugins without running a site build.

Python applications can use the same implementation through the bundled binary:

```python
from markata import workspace_files

listing = workspace_files("/path/to/site")
source = workspace_files("/path/to/site", path="posts/hello.md")["source"]
```

The CLI exposes the JSON interface as `markata workspace --workspace /path/to/site`
and `markata workspace --workspace /path/to/site --file posts/hello.md`.

The Go and FastAPI servers provide `GET /workspace` and
`GET /workspace/file?path=posts/hello.md`. The browser loads source on demand.
See the [Go editor guide](./go-workbench/) and [Python editor guide](./fastapi-editor/).

## Site defaults

When the selected workspace has a config, the editor applies its palette,
light/dark fallback, aesthetic, fontpack, text size, palette variables, custom CSS,
and Markdown highlighting settings. Local `palettes/` and `aesthetics/` catalogs
resolve from that workspace, as does an explicit `fontpacks_file`.
Saved browser appearance preferences take priority over site defaults.
The preview remains an article preview; it does not run page templates or a site build.

Python hosts can pass `workspace="/path/to/site"` to `appearance()`, `render()`,
and `render_async()`. The JSON command equivalent is
`markata render --workspace /path/to/site --json` (Markdown on stdin), or
`markata render --workspace /path/to/site --appearance` for controls and assets.
