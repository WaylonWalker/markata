---
date: 2025-12-09
description: "The plugin is responsible for loading and parsing markdown files with
  frontmatter into Post objects. It provides parallel loading capabilities and handles\u2026"
published: false
slug: markata/plugins/load
title: load.py


---

---

The `markata.plugins.load` plugin is responsible for loading and parsing markdown files
with frontmatter into Post objects. It provides parallel loading capabilities and
handles both modern Pydantic-based and legacy frontmatter validation.

## Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.load",
]
```

## Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.load",
]
```

Note: Disabling this plugin will prevent Markata from loading any markdown files.
This will effectively disable most of Markata's functionality.

## Configuration

Configure loading behavior in your `markata.toml`:

```toml
[markata]
# Directories containing markdown content
content_directories = [
    "content",
    "posts"
]

# Optional: Set to true to use legacy frontmatter validation
legacy_frontmatter = false

# Optional: Number of worker processes for parallel loading
load_workers = 4
```

## Functionality

## File Loading

The plugin:
1. Discovers markdown files in content directories
2. Loads file content and frontmatter
3. Validates frontmatter against Post model
4. Creates Post objects for further processing

## Parallel Processing

Loading is parallelized using:
- Process pool for file reading
- Configurable number of workers
- Chunked file processing

## Validation Modes

Supports two validation approaches:
1. Modern Pydantic-based validation (default)
   - Strict type checking
   - Automatic type coercion
   - Detailed validation errors

2. Legacy frontmatter validation
   - Looser type checking
   - Compatible with older content
   - Less strict validation

## Error Handling

The plugin provides:
- Detailed validation error messages
- Per-file error reporting
- Graceful fallback to legacy mode
- Optional strict validation

## Registered Attributes

The plugin adds:
- `articles`: List of loaded Post objects
- `content_directories`: List of content source directories

## Dependencies

This plugin depends on:
- python-frontmatter for YAML parsing
- pydantic for validation
- multiprocessing for parallel loading

---

!!! function
    <h2 id="load_file_content" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">load_file_content <em class="small">function</em></h2>

    Load file content without validation.

???+ source "load_file_content <em class='small'>source</em>"
    ```python
    def load_file_content(path: Path) -> tuple[Path, dict]:
        """Load file content without validation."""
        try:
            with open(path, "r") as f:
                raw_content = f.read()
            try:
                content = frontmatter.loads(raw_content)
            except Exception:
                content = None
            content["raw"] = raw_content
            return path, content
        except Exception:
            return path, None
    ```