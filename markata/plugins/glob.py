"""Default glob plugin"""
from pathlib import Path

from more_itertools import flatten

from markata.hookspec import hook_impl


@hook_impl
def glob(markata):
    markata.files = flatten(
        [Path().glob(pattern) for pattern in markata.config["glob_patterns"]]
    )
