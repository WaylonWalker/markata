from datetime import datetime

from rich.panel import Panel
from rich.table import Table


class Header:
    """Display header with clock."""

    def __rich__(self) -> Panel:
        grid = Table.grid(expand=True)
        grid.add_column(justify="center", ratio=1)
        grid.add_column(justify="right")
        grid.add_row(
            "[magenta][b]Markata[/b][/] [bright_black]Live Server[/]",
            datetime.now().ctime(),
        )
        return Panel(grid, style="yellow")
