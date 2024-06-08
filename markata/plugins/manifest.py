"""manifest plugin"""

import json
from pathlib import Path
from typing import TYPE_CHECKING


from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from markata.plugins.icon_resize import MarkataIcons


@hook_impl
def render(markata: "MarkataIcons") -> None:
    icons = markata.icons if "icons" in markata.__dict__ else []
    manifest = {
        "name": markata.config.site_name,
        "short_name": markata.config.short_name,
        "start_url": markata.config.start_url,
        "display": markata.config.display,
        "background_color": str(markata.config.background_color),
        "theme_color": str(markata.config.theme_color),
        "description": markata.config.description,
        "icons": icons,
    }
    filepath = Path(markata.config["output_dir"]) / "manifest.json"
    filepath.touch(exist_ok=True)
    with open(filepath, "w+") as f:
        json.dump(manifest, f, ensure_ascii=True, indent=4)
