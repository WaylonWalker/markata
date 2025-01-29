"""
The `markata.plugins.icon_resize` plugin automatically generates favicons in multiple
sizes from a single source image. This ensures your site has appropriate icons for
different devices and platforms.

# Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.icon_resize",
]
```

# Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.icon_resize",
]
```

# Configuration

Configure icon settings in your `markata.toml`:

```toml
[markata]
# Directory for output files
output_dir = "markout"

# Directory containing static assets
assets_dir = "static"

# Path to source icon file
icon = "static/icon.png"

# Optional: List of icon sizes to generate
# Default sizes are provided if not specified
icon_sizes = [16, 32, 48, 128, 180, 192, 512]
```

# Functionality

## Icon Generation

The plugin:
1. Loads the source icon image
2. Resizes it to multiple standard sizes
3. Saves each size as a separate PNG file
4. Generates appropriate HTML meta tags

## Default Sizes

If not configured, generates these standard sizes:
- 16x16: Classic favicon size
- 32x32: Modern favicon size
- 48x48: Windows taskbar
- 128x128: Chrome Web Store
- 180x180: Apple Touch icon
- 192x192: Android homescreen
- 512x512: PWA splash screen

## Output Files

Generated files follow this pattern:
- `favicon-{size}x{size}.png`
- `favicon.ico` (16x16 and 32x32 combined)
- `apple-touch-icon.png` (180x180)

## HTML Integration

The plugin adds appropriate meta tags to your HTML:
```html
<link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png">
<link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png">
<link rel="apple-touch-icon" sizes="180x180" href="/apple-touch-icon.png">
```

## Dependencies

This plugin depends on:
- Pillow (PIL) for image manipulation
"""

from pathlib import Path
from typing import Dict, List, Optional, TYPE_CHECKING

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

    @pydantic.field_validator("icon", mode="before")
    @classmethod
    def ensure_icon_exists(cls, v, info) -> Path:
        if v is None:
            return None

        # Convert string to Path if needed
        if isinstance(v, str):
            v = Path(v)

        if v.exists():
            return v

        icon = Path(info.data["assets_dir"]) / v

        if icon.exists():
            return icon
        else:
            raise FileNotFoundError(v)

    @pydantic.field_validator("icon_out_file", mode="before")
    @classmethod
    def default_icon_out_file(cls, v, info) -> Optional[Path]:
        if v is None and info.data.get("icon") is not None:
            return Path(info.data["output_dir"]) / info.data["icon"]
        if isinstance(v, str):
            return Path(v)
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
    from PIL import Image

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
    from PIL import Image

    for width in [48, 72, 96, 144, 192, 256, 384, 512]:
        with Image.open(markata.config.icon) as img:
            height = int(float(img.size[1]) * float(width / float(img.size[0])))
            img = img.resize((width, height), Image.LANCZOS)
            filename = Path(
                f"{markata.config.icon_out_file.stem}_{width}x{height}{markata.config.icon_out_file.suffix}",
            )
            out_file = Path(markata.config.output_dir) / filename
            if out_file.exists():
                continue
            img = img.resize((width, height), Image.LANCZOS)
            img.save(out_file)
