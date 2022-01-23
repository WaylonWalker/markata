import pdb
import sys
import traceback
from typing import TYPE_CHECKING, Callable, Optional

import typer
from rich.console import Console

from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from markata import Markata


def make_pretty() -> None:
    import click
    import pluggy
    import typer
    from rich import pretty as _pretty
    from rich import traceback

    _pretty.install()
    traceback.install(
        show_locals=True,
        suppress=[
            pluggy,
            click,
            typer,
        ],
    )


@hook_impl()
def cli(app: typer.Typer, markata: "Markata") -> None:
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
    ) -> None:
        """
        list posts
        """

        markata.console.quiet = True

        tail = -tail if tail else tail
        filtered = markata.map(map, filter, sort)
        if not include_empty:
            filtered = [a for a in filtered if a != ""]
        filtered = filtered[tail:head]
        if reverse:
            filtered = reversed(filtered)

        markata.console.quiet = False
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
        pretty: bool = True,
        quiet: bool = typer.Option(
            False,
            "--quiet",
            "-q",
        ),
        # to_dict: bool = False,
        watch: bool = False,
        verbose: bool = typer.Option(
            False,
            "--verbose",
            "-v",
        ),
        should_pdb: bool = typer.Option(
            False,
            "--pdb",
        ),
        profile: bool = True,
    ) -> None:
        import time

        if pretty:
            make_pretty()

        if quiet:
            markata.console.quiet = True

        if verbose:
            markata.console.print("console options:", markata.console.options)

        if watch:

            hash = markata.content_dir_hash
            markata.run()
            console = Console()
            with console.status("waiting for change", spinner="aesthetic", speed=0.2):
                while True:
                    if markata.content_dir_hash != hash:
                        hash = markata.content_dir_hash
                        markata.run()
                    time.sleep(0.1)

        if profile:
            markata.should_profile_cli = True
            markata.should_profile = True
            markata.configure()

        if should_pdb:
            pdb_run(markata.run)

        else:
            markata.run()


def pdb_run(func: Callable) -> None:
    try:
        func()
    except Exception:
        extype, value, tb = sys.exc_info()
        traceback.print_exc()
        pdb.post_mortem(tb)
