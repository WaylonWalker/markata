---
date: 2025-12-09
description: "The plugin generates custom cover images for your posts using configurable
  templates. It supports multiple cover image formats with different sizes, fonts,\u2026"
published: false
slug: markata/plugins/covers
title: covers.py


---

---

The `markata.plugins.covers` plugin generates custom cover images for your posts using
configurable templates. It supports multiple cover image formats with different sizes,
fonts, and layouts for various platforms (e.g., blog, social media, dev.to).

## Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.covers",
]
```

## Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.covers",
]
```

## Configuration

Configure multiple cover image templates in your `markata.toml`:

```toml
[[markata.covers]]
name = '-dev'                                    # Suffix for generated files
template = "static/cover-template.png"           # Base template image
font = "./static/JosefinSans-Regular.ttf"        # Title font
text_font = "./static/JosefinSans-Regular.ttf"   # Description font
font_color = "rgb(185,155,165)"                 # Title color
text_font_color = "rgb(255,255,255)"            # Description color
text_key = 'description'                        # Post attribute to use for description
padding = [0, 40, 100, 300]                     # Title padding [top, right, bottom, left]
text_padding = [0, 0]                           # Description padding

[[markata.covers]]
name = ''                                       # No suffix (default cover)
template = "static/og-template.png"             # Different template for social media
font = "./static/JosefinSans-Regular.ttf"
font_color = "rgb(255,255,255)"
text_font = "./static/JosefinSans-Regular.ttf"
text_font_color = "rgb(200,200,200)"
text_key = 'description'
padding = [10, 10, 100, 300]
text_padding = [0, 0]
```

### Configuration Options

Each cover configuration supports:
- `name`: Suffix for generated files (e.g., '-dev' creates 'post-dev.png')
- `template`: Path to template image
- `font`: Path to title font file
- `text_font`: Path to description font file
- `font_color`: RGB color for title
- `text_font_color`: RGB color for description
- `text_key`: Post attribute to use for description text
- `padding`: List of 1-4 integers for title padding
- `text_padding`: List of 1-4 integers for description padding

## Functionality

## Cover Generation

The plugin:
1. Loads the template image for each configuration
2. Draws the post title using specified font and color
3. Optionally draws description text
4. Applies configured padding
5. Saves the generated cover with appropriate suffix

## File Naming

Generated files follow this pattern:
- Main cover: `<slug>.png`
- Named covers: `<slug><name>.png`

Example:
- `my-post.png` (default cover)
- `my-post-dev.png` (dev.to cover)

## Dependencies

This plugin depends on:
- Pillow (PIL) for image manipulation

---

!!! function
    <h2 id="_lazy_import_pil" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">_lazy_import_pil <em class="small">function</em></h2>

    Lazy import PIL modules when needed.

???+ source "_lazy_import_pil <em class='small'>source</em>"
    ```python
    def _lazy_import_pil():
        """Lazy import PIL modules when needed."""
        global Image, ImageDraw, ImageFont
        if Image is None:
            from PIL import Image
            from PIL import ImageDraw
            from PIL import ImageFont
    ```
!!! function
    <h2 id="resolve_padding" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">resolve_padding <em class="small">function</em></h2>

    Convert padding to a len 4 tuple

???+ source "resolve_padding <em class='small'>source</em>"
    ```python
    def resolve_padding(padding: Tuple[int, ...], markata: "Markata") -> Tuple[int, ...]:
        """Convert padding to a len 4 tuple"""
        if len(padding) == 4:
            return padding
        if len(padding) == 3:
            return (*padding, padding[1])
        if len(padding) == 2:
            return padding * 2
        if len(padding) == 1:
            return padding * 4
        raise PaddingError(f"recieved padding: {padding}")
    ```