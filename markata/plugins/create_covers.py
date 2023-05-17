from pathlib import Path
from typing import TYPE_CHECKING

from PIL import Image, ImageDraw, ImageFont

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
