"""Default glob plugin"""
from pathlib import Path

from more_itertools import flatten

from markata.hookspec import hook_impl

from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from markata import Markata


@hook_impl
def glob(markata: "Markata") -> None:
    patterns = list(markata.config["glob_patterns"])
    markata.files = list(flatten([Path().glob(str(pattern)) for pattern in patterns]))
    markata.content_directories = list(set([f.parent for f in markata.files]))
