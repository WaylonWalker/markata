from markata import Markata
from markata.hookspec import hook_impl

from rich.console import Console
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

    @app.command()
    def build(
        rich: bool = False,
        quiet: bool = False,
        to_dict: bool = False,
        watch: bool = False,
    ) -> None:
        import time

        from rich import pretty, traceback

        if not rich:
            pretty.install()
            traceback.install()

        m = Markata()

        if quiet:
            m.console.quiet = True

        else:
            m.console.print("console options:", m.console.options)

        if to_dict:
            m.console.quiet = True
            data = m.to_dict()
            m.console.quiet = False
            m.console.print(data)
            return

        if watch:

            hash = m.content_dir_hash
            m.run()
            console = Console()
            with console.status("waiting for change", spinner="aesthetic", speed=0.2):
                while True:
                    if m.content_dir_hash != hash:
                        hash = m.content_dir_hash
                        m.run()
                    time.sleep(0.1)

        m.run()
