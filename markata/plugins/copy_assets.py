import shutil
from pathlib import Path

from rich.console import Console

from markata.hookspec import hook_impl

from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from markata import Markata


@hook_impl
def save(markata: "Markata") -> None:
    try:
        output_dir = Path(str(markata.config["output_dir"]))
        assets_dir = Path(str(markata.config["assets_dir"]))
    except KeyError:
        return

    with Console().status("copying assets", spinner="aesthetic", speed=0.2):
        shutil.copytree(assets_dir, output_dir, dirs_exist_ok=True)
