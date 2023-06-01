import atexit
import time
from pathlib import Path
from typing import TYPE_CHECKING, Union

import typer
from rich.panel import Panel

from markata.hookspec import hook_impl, register_attr

if TYPE_CHECKING:
    from markata import Markata


def find_port(port: int = 8000) -> int:
    """Find a port not in ues starting at given port"""
    import socket

    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        if s.connect_ex(("localhost", port)) == 0:
            return find_port(port=port + 1)
        return port


class Server:
    def __init__(
        self: "Server",
        *,
        auto_restart: bool = True,
        directory: Union[str, "Path"] = None,
        port: int = 8000,
    ) -> None:
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

    @property
    def title(self) -> str:
        return f"server ({self.uptime})"

    def kill(self) -> None:
        self.auto_restart = False
        self.proc.stdout.close()
        self.proc.stderr.close()
        self.proc.kill()
        self.proc.wait()

    def __rich__(self) -> Panel:
        if not self.proc.poll():
            return Panel(
                (
                    f"[green]serving on port: [gold1]{self.port} "
                    f"[green]using pid: [gold1]{self.proc.pid} "
                    f"[green]uptime: [gold1]{self.uptime} "
                    f"[green]link: [gold1] http://localhost:{self.port}[/]"
                ),
                border_style="blue",
                title=self.title,
                expand=True,
            )

        else:
            return Panel(
                "[red]server died",
                title=self.title,
                border_style="red",
                expand=True,
            )


@hook_impl
@register_attr("server")
def configure(markata: "Markata") -> None:
    def get_server(self):
        try:
            return self._server
        except AttributeError:
            self._server: Server = Server(directory=str(self.config["output_dir"]))
            return self._server

    from markata import Markata

    Markata.server = property(get_server)


def run_server() -> None:
    from rich.live import Live

    from .cli import run_until_keyboard_interrupt

    with Live(Server(), refresh_per_second=1, screen=True):
        run_until_keyboard_interrupt()


@hook_impl()
def cli(app: typer.Typer, markata: "Markata") -> None:
    server_app = typer.Typer()
    app.add_typer(server_app)

    @server_app.callback(invoke_without_command=True)
    def serve():
        """
        Serve the site locally.
        """
        run_server()


if __name__ == "__main__":
    run_server()
