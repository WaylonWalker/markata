---
date: 2025-12-09
description: "The plugin handles saving rendered HTML content to files. It determines
  the output path for each article and ensures files are saved in the correct location\u2026"
published: false
slug: markata/plugins/publish-html
title: publish_html.py


---

---

The `markata.plugins.publish_html` plugin handles saving rendered HTML content to files.
It determines the output path for each article and ensures files are saved in the correct
location within the output directory.

## Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.publish_html",
]
```

## Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.publish_html",
]
```

Note: Disabling this plugin will prevent HTML files from being written to disk.

## Configuration

Configure HTML output in `markata.toml`:

```toml
[markata]
# Base output directory
output_dir = "dist"

# Custom output paths
[[markata.output_paths]]
pattern = "blog/*"
output = "posts/{stem}.html"

[[markata.output_paths]]
pattern = "docs/*"
output = "documentation/{stem}/index.html"
```

## Functionality

## Path Resolution

The plugin:
1. Determines output path for each post
2. Creates necessary directories
3. Validates paths are within output_dir
4. Handles custom path mappings

## Output Model

Extends the base Post model with:
- output_html path
- Path validation
- Slug resolution
- Directory creation

## File Operations

Handles:
- Directory creation
- File writing
- Path validation
- Error logging

## Path Customization

Supports:
- Custom output paths
- Path patterns
- Directory structures
- Index files

## Safety Features

Includes:
- Path validation
- Directory verification
- Error handling
- Logging

## Dependencies

This plugin depends on:
- pathlib for path operations
- pydantic for model validation

---

!!! function
    <h2 id="save" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">save <em class="small">function</em></h2>

    Saves all the articles to their set `output_html` location if that location
    is relative to the specified `output_dir`.  If its not relative to the
    `output_dir` it will log an error and move on.

???+ source "save <em class='small'>source</em>"
    ```python
    def save(markata: "Markata") -> None:
        """
        Saves all the articles to their set `output_html` location if that location
        is relative to the specified `output_dir`.  If its not relative to the
        `output_dir` it will log an error and move on.
        """
        from slugify import slugify

        for article in markata.filter("not skip"):
            if article.html is None:
                continue

            if isinstance(article.html, str):
                # Create parent directories before writing
                current_html = (
                    article.output_html.read_text() if article.output_html.exists() else ""
                )
                if current_html != article.html:
                    article.output_html.parent.mkdir(parents=True, exist_ok=True)
                    article.output_html.write_text(article.html)
            elif isinstance(article.html, dict):
                for slug, html in article.html.items():
                    # Handle special case for index
                    if slug == "index":
                        output_path = article.output_html
                    # Handle files with extensions
                    elif "." in slug:
                        output_path = article.output_html.parent / slug
                    # Handle other slugs by creating subdirectories
                    else:
                        slug_path = slugify(slug)
                        output_path = article.output_html.parent / slug_path / "index.html"

                    # Create parent directories and write file
                    output_path.parent.mkdir(parents=True, exist_ok=True)
                    current_html = output_path.read_text() if output_path.exists() else ""
                    if current_html != html:
                        output_path.write_text(html)
    ```