---
date: 2025-12-09
description: "The plugin saves processed markdown files to the output directory, preserving
  frontmatter and content modifications. This enables source file access alongside\u2026"
published: false
slug: markata/plugins/publish-source
title: publish_source.py


---

---

The `markata.plugins.publish_source` plugin saves processed markdown files to the output
directory, preserving frontmatter and content modifications. This enables source file
access alongside rendered HTML.

## Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.publish_source",
]
```

## Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.publish_source",
]
```

## Configuration

No explicit configuration is required. The plugin automatically saves source files
alongside HTML output.

## Functionality

## File Output

The plugin:
1. Preserves markdown source
2. Maintains frontmatter
3. Creates output directories
4. Uses post slugs for paths

## Path Resolution

Source files are saved:
- At post's slug location
- With .md extension
- In output directory
- Alongside HTML files

## Frontmatter Handling

Handles frontmatter:
- Strips unserializable values
- Preserves YAML compatibility
- Maintains metadata
- Logs stripped fields

## Error Handling

Includes:
- YAML validation
- Path verification
- Directory creation
- Error logging

## Dependencies

This plugin depends on:
- python-frontmatter for YAML handling
- pyyaml for serialization

---

!!! function
    <h2 id="save" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">save <em class="small">function</em></h2>

    Saves the final modified post to the output site as markdown.

    !!! note
        Any keys that are not yaml serializable will be stripped.

???+ source "save <em class='small'>source</em>"
    ```python
    def save(markata: "Markata") -> None:
        """
        Saves the final modified post to the output site as markdown.

        !!! note
            Any keys that are not yaml serializable will be stripped.

        """
        output_dir = Path(str(markata.config["output_dir"]))
        for article in markata.filter(
            "not skip"
        ):  # iter_articles(description="saving source documents"):
            path = Path(
                output_dir / Path(article["slug"]).parent / Path(article["path"]).name,
            )
            current_content = path.read_text() if path.exists() else ""
            if current_content != article.dumps():
                path.parent.mkdir(parents=True, exist_ok=True)
                path.write_text(article.dumps())
    ```