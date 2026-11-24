---
title: "Python Source Docs"
description: "Generate API reference pages from Python source files with the optional python_docs plugin"
date: 2026-04-13
published: true
tags:
  - documentation
  - plugins
  - python
---


Use the optional `python_docs` plugin to turn Python modules into source-backed documentation pages.

The plugin only runs when it is explicitly listed in `hooks` and enabled in its config section.

## Example

```toml
[markata-go]
hooks = ["default", "python_docs"]

[markata-go.python_docs]
enabled = true
patterns = ["src/**/*.py", "pkg/**/*.py"]
exclude = ["**/tests/**", "**/.venv/**"]
slug_prefix = "api"
template = "docs"
include_source = true
include_module_code = false
published = false
```

This creates one generated post per discovered Python module.

Examples:

- `pkg/util.py` -> `/api/pkg/util/`
- `pkg/client/http.py` -> `/api/pkg/client/http/`
- `pkg/__init__.py` -> `/api/pkg/`

## What Gets Rendered

For each module, the plugin generates markdown that includes:

- the module name and source path
- the module docstring rendered as markdown
- an imports section
- an API index for classes and functions
- per-symbol sections with signatures and docstrings
- collapsible source snippets without repeated leading docstrings

## Cross-Linking

When multiple modules are documented together, the plugin links internal references where possible.

Supported cases include:

- import lists like `from pkg.util import greet`
- backticked references like `` `pkg.util` ``
- Sphinx-style refs like `:func:`pkg.util.greet``

## Common Options

```toml
[markata-go]
hooks = ["default", "python_docs"]

[markata-go.python_docs]
enabled = true
patterns = ["src/**/*.py"]
directories = ["pkg", "scripts"]
exclude = ["**/tests/**"]
use_gitignore = true
slug_prefix = "reference"
template = "docs"
published = false
include_private = false
include_source = true
include_module_code = false
tags = ["python", "docs"]
interpreter = "python3"
```

## Notes

- The plugin requires a Python interpreter on `PATH`.
- The plugin must be listed in `[markata-go].hooks` as `python_docs`.
- Markdown content loading is unchanged; Python docs are added alongside normal posts.
- Generated docs are unpublished by default so they stay out of feeds and sitemap unless you opt in.

## Related

- [Built-in Plugins Reference](../reference/plugins.md)
- [Configuration Reference](./configuration.md)
