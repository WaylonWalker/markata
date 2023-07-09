"""Icon Resize Plugin

Resized favicon to a set of common sizes.

## markata.plugins.icon_resize configuration

```toml title=markata.toml
[markata]
output_dir = "markout"
assets_dir = "static"
icon = "static/icon.png"
```

"""
from pathlib import Path
from typing import Dict, List, Optional, TYPE_CHECKING

from PIL import Image

from markata.hookspec import register_attr

if TYPE_CHECKING:
    from markata import Markata

import pydantic

from markata.hookspec import hook_impl


class Config(pydantic.BaseModel):
    output_dir: pydantic.DirectoryPath = "markout"
    assets_dir: Path = pydantic.Field(
        Path("static"),
        description="The directory to store static assets",
    )
    icon: Optional[Path] = None
    icon_out_file: Optional[Path] = None
    icons: Optional[List[Dict[str, str]]] = []

    @pydantic.validator("icon")
    def ensure_icon_exists(cls, v, *, values: Dict) -> Path:
        if v is None:
            return
        if v.exists():
            return v

        icon = Path(values["assets_dir"]) / v

        if icon.exists():
            return icon
        else:
            raise FileNotFoundError(v)

    @pydantic.validator("icon_out_file", pre=True, always=True)
    def default_icon_out_file(cls, v, *, values: Dict) -> Path:
        if v is None and values["icon"] is not None:
            return Path(values["output_dir"]) / values["icon"]
        return v


@hook_impl()
@register_attr("config_models")
def config_model(markata: "Markata") -> None:
    markata.config_models.append(Config)


@hook_impl
@register_attr("icons")
def render(markata: "Markata") -> None:
    if markata.config.icon is None:
        return

    with Image.open(markata.config.icon) as img:
        for width in [48, 72, 96, 144, 192, 256, 384, 512]:
            height = int(float(img.size[1]) * float(width / float(img.size[0])))
            filename = Path(
                f"{markata.config.icon_out_file.stem}_{width}x{height}{markata.config.icon_out_file.suffix}",
            )
            markata.config.icons.append(
                {
                    "src": str(filename),
                    "sizes": f"{width}x{width}",
                    "type": f"image/{img.format}".lower(),
                    "purpose": "any maskable",
                },
            )


@hook_impl
def save(markata: "Markata") -> None:
    if markata.config.icon is None:
        return
    for width in [48, 72, 96, 144, 192, 256, 384, 512]:
        with Image.open(markata.config.icon) as img:
            height = int(float(img.size[1]) * float(width / float(img.size[0])))
            filename = Path(
                f"{markata.config.icon_out_file.stem}_{width}x{height}{markata.config.icon_out_file.suffix}",
            )
            out_file = Path(markata.config.output_dir) / filename
            if out_file.exists():
                continue
            img = img.resize((width, height), Image.LANCZOS)
            img.save(out_file)
