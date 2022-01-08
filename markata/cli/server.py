import time
from typing import TYPE_CHECKING, Union

from rich.panel import Panel

if TYPE_CHECKING:
    from pathlib import Path

import atexit


def find_port(port=8000):
    """Find a port not in ues starting at given port"""
    import socket

    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        if s.connect_ex(("localhost", port)) == 0:
            return find_port(port=port + 1)
        else:
            return port


class Server:
    def __init__(
        self,
        auto_restart: bool = True,
        directory: Union[str, "Path"] = None,
        port: int = 8000,
    ):
        if directory is None:
            from markata import Markata

            m = Markata()
            directory = m.config["output_dir"]

        self.auto_restart = auto_restart
        self.directory = directory
        self.port = find_port(port=port)
        self.start_server()
        atexit.register(self.kill)

    def start_server(self):
        import subprocess

        self.cmd = [
            "python",
            "-m",
            "http.server",
            str(self.port),
            "--directory",
            self.directory,
        ]
        # self.cmd = [
        #     "pipx",
        #     "run",
        #     "livereload",
        #     "-p",
        #     str(self.port),
        #     # "-w",
        #     # ".1",
        #     self.directory,
        # ]

        self.proc = subprocess.Popen(
            self.cmd,
            stderr=subprocess.PIPE,
            stdout=subprocess.PIPE,
        )
        self.start_time = time.time()

    @property
    def uptime(self):
        return round(time.time() - self.start_time)

    def kill(self):
        self.auto_restart = False
        self.proc.kill()

    def __rich__(self) -> Panel:
        if not self.proc.poll():
            return Panel(
                f"[green]serving on port: [gold1]{self.port} [green]using pid: [gold1]{self.proc.pid} [green]uptime: [gold1]{self.uptime} [green]link: [gold1] http://localhost:{self.port}[/]",
                border_style="blue",
                title="server",
            )

        else:
            if self.auto_restart:
                self.start_server()

            return Panel(f"[red]server died", title="server", border_style="red")


if __name__ == "__main__":

    from rich.live import Live

    from markata import Markata

    from .cli import run_until_keyboard_interrupt

    m = Markata()
    with Live(Server(), refresh_per_second=1, screen=True):
        run_until_keyboard_interrupt()
