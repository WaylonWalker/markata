"""
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
"""

import time
from functools import lru_cache
from pathlib import Path
from typing import TYPE_CHECKING
from typing import List
from typing import Optional
from typing import Tuple
from typing import Union

from rich.progress import BarColumn
from rich.progress import Progress

from markata import background
from markata.hookspec import hook_impl

# Lazy imports for PIL
Image = None
ImageDraw = None
ImageFont = None


def _lazy_import_pil():
    """Lazy import PIL modules when needed."""
    global Image, ImageDraw, ImageFont
    if Image is None:
        from PIL import Image
        from PIL import ImageDraw
        from PIL import ImageFont


if TYPE_CHECKING:
    from markata import Markata


@lru_cache(maxsize=64)
def _load_font(path: Path, size: int) -> "ImageFont.FreeTypeFont":
    _lazy_import_pil()
    return ImageFont.truetype(path, size=size)


def get_font(
    path: Path,
    draw: "ImageDraw.Draw",
    title: str,
    size: int = 250,
    max_size: tuple = (800, 220),
) -> "ImageFont.FreeTypeFont":
    title = title or ""
    font = _load_font(path, size)
    current_size = draw.textsize(title, font=font)

    if current_size[0] > max_size[0] or current_size[1] > max_size[1]:
        return get_font(path, draw, title, size - 10, max_size=max_size)
    return font


class PaddingError(BaseException):
    def __init__(
        self,
        msg: str = "",
    ) -> None:
        super().__init__(
            "Padding must be an iterable of length 1, 2, 3, or 4.\n" + msg,
        )


def draw_text(
    image: "Image",
    font_path: Optional[Path],
    text: str,
    color: Union[str, None],
    padding: Tuple[int, ...],
    markata: "Markata",
) -> None:
    _lazy_import_pil()
    text = text or ""
    draw = ImageDraw.Draw(image)
    padding = resolve_padding(padding, markata)
    width = image.size[0]
    height = image.size[1]
    bounding_box = [padding[0], padding[1], width - padding[0], height - padding[1]]
    bounding_box = [padding[0], padding[1], width - padding[2], height - padding[3]]
    max_size = (bounding_box[2] - bounding_box[0], bounding_box[3] - bounding_box[1])
    x1, y1, x2, y2 = bounding_box
    font = get_font(font_path, draw, text, max_size=max_size) if font_path else None
    w, h = draw.textsize(text, font=font)
    x = (x2 - x1 - w) / 2 + x1
    y = (y2 - y1 - h) / 2 + y1
    draw.text((x, y), text, fill=color, font=font, align="center")


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


@background.task
def make_cover(
    title: str,
    color: str,
    output_path: Path,
    template_path: Path,
    font_path: Optional[Path],
    padding: Tuple[int, ...],
    text_font: Path,
    text: str = None,
    text_font_color: str = None,
    text_padding: Tuple[int, ...] = None,
    resizes: List[int] = None,
    markata: "Markata" = None,
) -> None:
    _lazy_import_pil()
    if output_path.exists():
        return
    image = Image.open(template_path) if template_path else Image.new("RGB", (800, 450))

    draw_text(
        image=image,
        font_path=font_path,
        title=title,
        color=color,
        padding=padding,
        markata=markata,
    )
    if text is not None:
        if text_padding is None:
            text_padding = (
                image.size[1] - image.size[1] / 5,
                image.size[0] / 5,
                image.size[1] - image.size[1] / 10,
            )
        draw_text(image, text_font, text, text_font_color, text_padding)

    image.save(output_path)
    ratio = image.size[1] / image.size[0]

    covers = []
    if resizes:
        for width in resizes:
            re_img = image.resize((width, int(width * ratio)), Image.ANTIALIAS)
            filename = (
                f"{output_path.stem}_{width}x{int(width * ratio)}{output_path.suffix}"
            )
            covers.append(filename)

            filepath = Path(output_path.parent / filename)
            re_img.save(filepath)


@hook_impl
def save(markata: "Markata") -> None:
    futures = []

    if "covers" not in markata.config.keys():
        return

    for article in markata.iter_articles("making covers"):
        for cover in markata.config["covers"]:
            try:
                padding = cover["padding"]
            except KeyError:
                padding = (
                    200,
                    100,
                )
            try:
                text_padding = cover["text_padding"]
            except KeyError:
                text_padding = (
                    200,
                    100,
                )
            if "text_key" in cover:
                try:
                    text = article.metadata[cover["text_key"]]
                except AttributeError:
                    text = article[cover["text_key"]]
                try:
                    text = text.replace("\n", "")
                    from more_itertools import chunked

                    text = "\n".join(["".join(c) for c in chunked(text, 60)])
                except AttributeError:
                    # text is likely None
                    pass

                text_font = cover["text_font"]
                text_font_color = cover["text_font_color"]
            else:
                text = None
                text_font = None
                text_font_color = None
            try:
                title = article.metadata["title"]
            except AttributeError:
                title = article["title"]
            futures.append(
                make_cover(
                    title=title,
                    color=cover["font_color"],
                    output_path=Path(markata.config.output_dir)
                    / (article["slug"] + cover["name"] + ".png"),
                    template_path=cover.get("template", None),
                    font_path=cover.get("font", None),
                    padding=padding,
                    text_font=text_font,
                    text=text,
                    text_font_color=text_font_color,
                    text_padding=text_padding,
                    resizes=cover.get("resizes"),
                    markata=markata,
                ),
            )

    progress = Progress(
        BarColumn(bar_width=None),
        transient=True,
        console=markata.console,
    )
    task_id = progress.add_task("loading markdown")
    progress.update(task_id, total=len(futures))
    with progress:
        while not all(f.done() for f in futures):
            time.sleep(0.1)
            progress.update(task_id, total=len([f for f in futures if f.done()]))
    [f.result() for f in futures]
