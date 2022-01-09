from typing import Optional

import typer
from rich.console import Console

from markata import Markata
from markata.hookspec import hook_impl


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
        use_pager: bool = typer.Option(True, "--pager", "--no-pager"),
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

        if markata.console.is_terminal and use_pager:
            with markata.console.pager():
                for a in filtered:
                    markata.console.print(a, style="purple")
        else:
            for a in filtered:
                markata.console.print(a)

    @app.command()
    def console():
        if markata.console.is_terminal:
            markata.console.print("console options:", markata.console.options)
        else:
            print("here")

    @app.command()
    def build(
        rich: bool = False,
        quiet: bool = False,
        to_dict: bool = False,
        watch: bool = False,
        verbose: bool = typer.Option(
            False,
            "--verbose",
            "-v",
        ),
    ) -> None:
        import time

        from rich import pretty, traceback

        if not rich:
            pretty.install()
            traceback.install()

        m = Markata()

        if quiet:
            m.console.quiet = True

        if verbose:
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
