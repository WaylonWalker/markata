from datetime import datetime
from markata import Markata
from checksumdir import dirhash

from rich import box
from rich.align import Align
from rich.console import Console, RenderGroup
from rich.layout import Layout
from rich.panel import Panel
from rich.progress import Progress, SpinnerColumn, BarColumn, TextColumn
from rich.syntax import Syntax
from rich.table import Table
from rich.text import Text


from rich.live import Live


def make_layout() -> Layout:
    """Define the layout."""
    layout = Layout(name="root")

    layout.split(
        Layout(name="header", size=3),
        Layout(name="main", ratio=1),
        Layout(name="footer", size=7),
    )
    layout["main"].split(
        Layout(name="side"),
        Layout(name="body", ratio=2, minimum_size=60),
        direction="horizontal",
    )
    layout["side"].split(
        Layout(name="plugins"),
    )
    return layout


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
        return Panel(grid, style="yellow on black")


class RichM:
    def __init__(self, markata: Markata):
        self.m = markata


class Footer(RichM):
    """Display Footer"""

    def __rich__(self) -> Panel:

        s = f"watching: {self.m.output_dir} - {self.m._dirhash} - {self.m.count}"
        return Panel(Text(s))


class Plugins(RichM):
    def __rich__(self) -> Panel:
        plugin_table = Table.grid(expand=True)
        plugin_table.add_row(f"[bright_blue]{len(self.m._pm.get_plugins())}[/] plugins")
        for plugin in self.m._pm.get_plugins():
            plugin_table.add_row(
                "".join(
                    [
                        "[bright_black]",
                        ".".join(plugin.__name__.split(".")[:-1]),
                        ".[/]",
                        plugin.__name__.split(".")[-1],
                    ]
                )
            )
        return plugin_table


# import logging
# from rich.logging import RichHandler
# from markata.tail_logger import TailLogger

# import io

# tail = TailLogger(10)
# log_capture_string = io.StringIO()
# ch = logging.StreamHandler(log_capture_string)
# ch.setLevel(logging.DEBUG)


# FORMAT = "%(message)s"
# logging.basicConfig(
#     level="INFO",
#     format=FORMAT,
#     datefmt="[%X]",
#     handlers=[RichHandler()],
# )

# ch.setFormatter(FORMAT)
# # tail.setFormatter(FORMAT)

# log = logging.getLogger(__name__)
# log.info("Hello, World!")
# log.addHandler(ch)
# log.addHandler(tail)


def run():
    layout = make_layout()
    layout["header"].update(Header())

    m = Markata()

    layout["plugins"].update(Plugins(m))
    m.count = 0
    m._dirhash = ""
    layout["footer"].update(Footer(m))

    with Live(layout, refresh_per_second=1, screen=True, console=m.console) as live:
        m._dirhash = m.content_dir_hash

        while m._dirhash == m.content_dir_hash:
            # layout["footer"].update(f"watching: {m.output_dir} - {_dirhash} - {count}")
            m.count += 1
    print("files changed")
    run()


if __name__ == "__main__":
    # console = Console()
    # console.log("here")
    run()
