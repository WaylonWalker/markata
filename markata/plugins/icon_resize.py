"""Icon Resize Plugin"""
from pathlib import Path
from typing import TYPE_CHECKING

from PIL import Image

from markata.hookspec import hook_impl

if TYPE_CHECKING:

    from typing import Dict, List

    from markata import Markata

    class MarkataIcons(Markata):
        icons: List[Dict[str, str]]


@hook_impl
def render(markata: "MarkataIcons") -> None:
    base_out_file = Path(markata.output_dir) / markata.icon
    img = Image.open(Path(markata.assets_dir) / markata.icon)
    ratio = img.size[1] / img.size[0]
    markata.icons = []
    for width in [48, 72, 96, 144, 192, 256, 384, 512]:
        height = int(float(img.size[1]) * float(width / float(img.size[0])))
        filename = Path(f"{base_out_file.stem}_{width}x{height}{base_out_file.suffix}")
        markata.icons.append(
            {
                "src": str(filename),
                "sizes": f"{width}x{width}",
                "type": f"image/{img.format}".lower(),
                "purpose": "any maskable",
            }
        )


@hook_impl
def save(markata: "MarkataIcons") -> None:
    base_out_file = Path(markata.output_dir) / markata.icon
    for width in [48, 72, 96, 144, 192, 256, 384, 512]:
        img = Image.open(Path(markata.assets_dir) / markata.icon)
        height = int(float(img.size[1]) * float(width / float(img.size[0])))
        img = img.resize((width, height), Image.ANTIALIAS)
        filename = Path(f"{base_out_file.stem}_{width}x{height}{base_out_file.suffix}")
        out_file = Path(markata.output_dir) / filename
        img.save(out_file)
