from datetime import datetime

import time
from checksumdir import dirhash
from rich import box
from rich.align import Align
from rich.console import Console, RenderGroup
from rich.layout import Layout
from rich.live import Live
from rich.panel import Panel
from rich.progress import BarColumn, Progress, SpinnerColumn, TextColumn
from rich.syntax import Syntax
from rich.table import Table
from rich.text import Text

from markata import Markata


class MarkataCli(Markata):
    count: int = 0
    _dirhash: str = ""


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
        Layout(name="server"),
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
    def __init__(self, markata: MarkataCli):
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
        return Panel(plugin_table)


def find_port(port=8000):
    """Find a port not in ues starting at given port"""
    import socket, errno

    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        if s.connect_ex(("localhost", port)) == 0:
            return find_port(port=port + 1)
        else:
            return port


from rich import box


class Server:
    def __init__(
        self, auto_restart: bool = True, directory: str = ".", port: int = 8000
    ):

        self.auto_restart = auto_restart
        self.directory = directory
        self.port = find_port(port=port)
        self.start_server()

    def start_server(self):
        import subprocess

        self.proc = subprocess.Popen(
            [
                "python",
                "-m",
                "http.server",
                str(self.port),
                "--directory",
                self.directory,
            ],
            stderr=subprocess.PIPE,
            stdout=subprocess.PIPE,
        )
        self.start_time = time.time()

    @property
    def uptime(self):
        return round(time.time() - self.start_time)

    def __rich__(self) -> Panel:
        if not self.proc.poll():
            return Panel(
                f"[green]serving on port: [gold1]{self.port} [green]using pid: [gold1]{self.proc.pid} [green]uptime: [gold1]{self.uptime}[/]",
                border_style="blue",
                title="server",
            )

        else:
            if self.auto_restart:
                self.start_server()

            return Panel(f"[red]server died", title="server", border_style="red")


def run() -> None:
    layout = make_layout()
    layout["header"].update(Header())

    m = MarkataCli()
    server = Server(directory=m.output_dir)
    layout["server"].update(server)

    layout["plugins"].update(Plugins(m))
    m.count = 0
    m._dirhash = ""
    layout["footer"].update(Footer(m))

    with Live(layout, refresh_per_second=1, screen=True, console=m.console) as live:
        m._dirhash = m.content_dir_hash

        while m._dirhash == m.content_dir_hash:
            m.count += 1
    print("files changed")
    run()


if __name__ == "__main__":
    run()
