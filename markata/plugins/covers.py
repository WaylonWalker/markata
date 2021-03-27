import time
from functools import lru_cache
from pathlib import Path
from typing import TYPE_CHECKING, Tuple, Union

from PIL import Image, ImageDraw, ImageFont

from markata import background
from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from markata import Markata


@lru_cache(maxsize=64)
def _load_font(path: Path, size: int) -> ImageFont.FreeTypeFont:
    return ImageFont.truetype(path, size=size)


def get_font(
    path: Path,
    draw: ImageDraw.Draw,
    title: str,
    size: int = 250,
    max_size: tuple = (800, 220),
) -> ImageFont.ImageFont.FreeTypeFont:
    font = _load_font(path, size)
    # font = ImageFont.truetype(path, size=size)
    current_size = draw.textsize(title, font=font)
    if current_size[0] > max_size[0] or current_size[1] > max_size[1]:
        return get_font(path, draw, title, size - 10)
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
    image: Image,
    font_path: Path,
    text: str,
    color: Union[str, None],
    padding: Tuple[int, ...],
) -> None:
    draw = ImageDraw.Draw(image)
    font = get_font(font_path, draw, text)
    padding = resolve_padding(padding)
    width = image.size[0]
    height = image.size[1]
    bounding_box = [padding[0], padding[1], width - padding[0], height - padding[1]]
    x1, y1, x2, y2 = bounding_box
    w, h = draw.textsize(text, font=font)
    x = (x2 - x1 - w) / 2 + x1
    y = (y2 - y1 - h) / 2 + y1
    draw.text((x, y), text, fill=color, font=font, align="center")


def resolve_padding(padding: Tuple[int, ...]) -> Tuple[int, ...]:
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
    font_path: Path,
    padding: Tuple[int, ...],
    text_font: Path,
    text: str = None,
    text_font_color: str = None,
    text_padding: Tuple[int, ...] = None,
) -> None:
    # image = Image.open(template_path)
    image = Image.open(template_path)
    draw_text(image, font_path, title, color, padding)
    if text is not None:
        if text_padding is None:
            text_padding = (
                image.size[1] / 2,
                image.size[0] / 5,
                image.size[1] - image.size[1] / 5,
            )
        draw_text(image, text_font, text, text_font_color, text_padding)

    image.save(output_path)
    # image.save(output_path.with_suffix(".webp"), quality=80, optimize=True)
    ratio = image.size[1] / image.size[0]

    covers = []
    for width in [
        # 32,
        250,
        # 500,
    ]:

        re_img = image.resize((width, int(width * ratio)), Image.ANTIALIAS)
        filename = f"{output_path.stem}_{width}x{int(width*ratio)}{output_path.suffix}"
        covers.append(filename)

        filepath = Path(output_path.parent / filename)
        re_img.save(filepath)
        # filename = f"{output_path.stem}_{width}x{int(width*ratio)}.webp"
        # covers.append(filename)
        # filepath = Path(output_path.parent / filename)
        # re_img.save(filepath)


@hook_impl
def save(markata: "Markata") -> None:
    futures = []
    for article in markata.articles:
        for cover in markata.config["covers"]:
            try:
                padding = cover["padding"]
            except KeyError:
                padding = (
                    200,
                    100,
                )
            if "text_key" in cover:
                text = article.metadata[cover["text_key"]]
                text_font = cover["text_font"]
                text_font_color = cover["text_font_color"]
            else:
                text = None
                text_font = None
                text_font_color = None
            futures.append(
                make_cover(
                    article.metadata["title"],
                    cover["font_color"],
                    Path(markata.config["output_dir"])
                    / (article["slug"] + cover["name"] + ".png"),
                    cover["template"],
                    cover["font"],
                    padding,
                    text_font,
                    text,
                    text_font_color,
                )
            )

    with markata.iter_articles("creating cover images") as pbar:
        while not all([f.done() for f in futures]):
            time.sleep(0.1)
            for _ in range(len([f for f in futures if f.done()]) - pbar.n):
                pbar.update()
    # necessary  to not fail silently
    [f.result() for f in futures]
