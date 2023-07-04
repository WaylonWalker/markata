import shutil
from typing import TYPE_CHECKING

from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from markata import Markata


@hook_impl
def save(markata: "Markata") -> None:
    with markata.console.status("copying assets", spinner="aesthetic", speed=0.2):
        if markata.config.assets_dir.exists():
            shutil.copytree(
                markata.config.assets_dir,
                markata.config.output_dir,
                dirs_exist_ok=True,
            )
