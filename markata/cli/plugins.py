from typing import TYPE_CHECKING

import rich
from rich.panel import Panel
from rich.table import Table

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


if __name__ == "__main__":

    plugins = Plugins(Markata())
    rich.print(plugins)
