from typing import TYPE_CHECKING

import rich
import typer
from rich import print as rich_print
from rich.panel import Panel
from rich.table import Table

from markata.hookspec import hook_impl, register_attr

if TYPE_CHECKING:
    from markata import Markata


class Plugins:
    def __init__(self, markata: "Markata"):
        self.m = markata

    def __rich__(self) -> Panel:
        grid = Table.grid(expand=True)
        num_plugins = f"[bright_blue]({len(self.m._pm.get_plugins())})[/]"

        for plugin in self.m._pm.get_plugins():
            grid.add_row(
                "".join(
                    [
                        "[bright_black]",
                        ".".join(plugin.__name__.split(".")[:-1]),
                        ".[/]",
                        plugin.__name__.split(".")[-1],
                    ]
                )
            )
        return Panel(
            grid, title=f"plugins {num_plugins}", border_style="gold1", expand=False
        )


@hook_impl
@register_attr("plugins")
def configure(markata: "Markata") -> None:
    def get_plugins(self):
        try:
            return self._plugins
        except AttributeError:
            self._plugins: Plugins = Plugins(self)
            return self._plugins

    from markata import Markata

    Markata.plugins = property(get_plugins)


@hook_impl
def cli(app: typer.Typer, markata: "Markata") -> None:
    feeds_app = typer.Typer()
    app.add_typer(feeds_app)

    @feeds_app.callback()
    def feeds():
        "feeds cli"

    @feeds_app.command()
    def show() -> None:
        markata.console.quiet = True
        feeds = markata.feeds
        markata.console.quiet = False
        rich_print(feeds)


if __name__ == "__main__":

    plugins = Plugins(Markata())
    rich.print(plugins)
