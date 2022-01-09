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

    try:
        ignore = markata.config["glob"]["use_gitignore"] or True
    except KeyError:
        ignore = True

    if ignore and (Path(".gitignore").exists() or Path(".markataignore").exists()):
        import pathspec

        lines = []

        if Path(".gitignore").exists():
            lines.extend(Path(".gitignore").read_text().splitlines())

        if Path(".markataignore").exists():
            lines.extend(Path(".markataignore").read_text().splitlines())

        spec = pathspec.PathSpec.from_lines("gitwildmatch", lines)

        markata.files = [
            file for file in markata.files if not spec.match_file(str(file))
        ]
