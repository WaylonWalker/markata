import shutil
from pathlib import Path
from typing import TYPE_CHECKING

from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from markata import Markata


@hook_impl
def save(markata: "Markata") -> None:
    output_dir = Path(str(markata.config.get("output_dir", "markout")))
    assets_dir = Path(str(markata.config.get("assets_dir", "static")))

    with markata.console.status("copying assets", spinner="aesthetic", speed=0.2):
        if assets_dir.exists():
            shutil.copytree(assets_dir, output_dir, dirs_exist_ok=True)
