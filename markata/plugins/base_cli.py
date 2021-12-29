from markata.hookspec import hook_impl

import typer
from typing import Optional


@hook_impl()
def cli(app, markata):
    @app.command()
    def list(
        map: str = "title",
        filter: str = "True",
        sort: str = "True",
        head: Optional[int] = None,
        tail: Optional[int] = None,
        include_empty: bool = False,
        reverse: bool = False,
    ):
        """
        list posts
        """
        tail = -tail if tail else tail
        filtered = markata.map(map, filter, sort)
        if not include_empty:
            filtered = [a for a in filtered if a != ""]
        filtered = filtered[tail:head]
        if reverse:
            filtered = reversed(filtered)

        for a in filtered:
            typer.echo(a)
