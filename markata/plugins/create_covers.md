---
date: 2025-12-09
description: "The plugin is a simplified version of the covers plugin. It generates
  basic cover images with titles using a single template and font configuration.\u2026"
published: false
slug: markata/plugins/create-covers
title: create_covers.py


---

---

The `markata.plugins.create_covers` plugin is a simplified version of the covers plugin.
It generates basic cover images with titles using a single template and font configuration.

## Installation

This plugin is built-in but NOT enabled by default. To use it, add it to your plugins:

```toml
hooks = [
    "markata.plugins.create_covers",
]
```

## Uninstallation

Since this plugin is not in the default plugin set, simply remove it from your hooks list
to disable it:

```toml
hooks = [
    # Remove or comment out the line below
    # "markata.plugins.create_covers",
]
```

## Configuration

Configure basic cover settings in your `markata.toml`:

```toml
[markata.create_covers]
template = "static/cover-template.png"  # Base template image
font = "./static/OpenSans-Regular.ttf"  # Font file for title
font_color = "rgb(255,255,255)"        # Title color (white)
```

## Functionality

## Cover Generation

The plugin:
1. Loads a single template image
2. Draws the post title using the configured font
3. Automatically sizes text to fit within image bounds
4. Saves the cover as PNG in the output directory

## File Naming

Generated files are named:
- `<slug>.png` in the output directory

## Dependencies

This plugin depends on:
- Pillow (PIL) for image manipulation

Note: For more advanced cover generation with multiple templates and text positioning,
use the `markata.plugins.covers` plugin instead.

---