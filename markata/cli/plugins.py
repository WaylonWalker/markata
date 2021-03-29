from rich.panel import Panel
from rich.table import Table

from .cli import RichM


class Plugins(RichM):
    def __rich__(self) -> Panel:
        grid = Table.grid(expand=True)
        grid.add_row(f"[bright_blue]{len(self.m._pm.get_plugins())}[/] plugins")
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
        return Panel(grid, title="plugins", border_style="gold1")
