import atexit
import time
from pathlib import Path
from typing import Union

from rich.panel import Panel


def find_port(port: int = 8000) -> int:
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
            directory = Path(str(m.config["output_dir"]))

        self.auto_restart = auto_restart
        self.directory = directory
        self.port = find_port(port=port)
        self.start_server()
        atexit.register(self.kill)

    def start_server(self) -> None:
        import subprocess

        self.cmd = [
            "python",
            "-m",
            "http.server",
            str(self.port),
            "--directory",
            self.directory,
        ]

        self.proc = subprocess.Popen(
            self.cmd,
            stderr=subprocess.PIPE,
            stdout=subprocess.PIPE,
        )
        self.start_time = time.time()

    @property
    def uptime(self) -> int:
        return round(time.time() - self.start_time)

    def kill(self) -> None:
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

            return Panel(f"[red]server died", title="server", border_style="red")


if __name__ == "__main__":

    from rich.live import Live

    from markata import Markata

    from .cli import run_until_keyboard_interrupt

    m = Markata()
    with Live(Server(), refresh_per_second=1, screen=True):
        run_until_keyboard_interrupt()
