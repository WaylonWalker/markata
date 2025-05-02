"""
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
"""

from pathlib import Path
from typing import TYPE_CHECKING

from PIL import Image
from PIL import ImageDraw
from PIL import ImageFont

from markata import background
from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from markata import Markata


def get_font(
    path: Path,
    draw: ImageDraw.Draw,
    title: str,
    size: int = 250,
) -> ImageFont.FreeTypeFont:
    font = ImageFont.truetype(path, size=size)
    if draw.textsize(title, font=font)[0] > 800:
        return get_font(path, draw, title, size - 10)
    return font


@background.task
def make_cover(
    title: str,
    color: str,
    output_path: Path,
    template_path: Path,
    font_path: Path,
) -> None:
    image = Image.open(template_path)

    draw = ImageDraw.Draw(image)

    font = get_font(font_path, draw, title)

    color = "rgb(255,255,255)"
    padding = (200, 100)
    bounding_box = [padding[0], padding[1], 1000 - padding[0], 420 - padding[1]]
    x1, y1, x2, y2 = bounding_box
    w, h = draw.textsize(title, font=font)
    x = (x2 - x1 - w) / 2 + x1
    y = (y2 - y1 - h) / 2 + y1
    draw.text((x, y), title, fill=color, font=font, align="center")
    image.save(output_path)


@hook_impl
def save(markata: "Markata") -> None:
    for article in markata.articles:
        output_path = Path(markata.output_dir) / (
            Path(article.metadata["path"]).stem + ".png"
        )

        make_cover(
            article.metadata["title"],
            markata.config["cover_font_color"],
            output_path,
            markata.config["cover_template"],
            markata.config["cover_font"],
        )
