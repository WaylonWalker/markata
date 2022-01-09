"""Default glob plugin"""
from pathlib import Path
from typing import TYPE_CHECKING

from more_itertools import flatten

from markata.hookspec import hook_impl, register_attr

if TYPE_CHECKING:
    from markata import Markata


@hook_impl
@register_attr("content_directories", "files")
def glob(markata: "Markata") -> None:
    markata.files = list(
        flatten([Path().glob(str(pattern)) for pattern in markata.glob_patterns])
    )
    markata.content_directories = list(set([f.parent for f in markata.files]))
