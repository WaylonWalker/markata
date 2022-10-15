import atexit
import subprocess
import time
from typing import TYPE_CHECKING

from rich.panel import Panel
from rich.text import Text

if TYPE_CHECKING:
    from markata import Markata


class Runner:
    """Display Footer"""

    _status = "waiting"
    status = "starting"
    last_error = ""
    title = "runner"
    border = "green"

    _dirhash = ""
    time = time.time()

    def __init__(self, markata: "Markata") -> None:
        print("registering kill")
        atexit.register(self.kill)
        self.m = markata
        self.build()

    def build(self) -> subprocess.Popen:
        self.status = "running"
        self.time = time.time()
        self.proc = subprocess.Popen(
            ["markata", "build"],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
        )
        print(f"starting pid {self.proc.pid}")

    def kill(self):
        self.proc.stdout.close()
        self.proc.stderr.close()
        self.proc.kill()
        self.proc.wait()

    def run(self) -> None:
        """ """
        if self.proc.poll() is not None:
            self.proc.kill()
            self.build()
            print(f"starting pid {self.proc.pid}")
        print("already running")

    @property
    def status_message(self) -> str:
        """"""
        num_lines = self.m.console.height - 4
        last_error = "\n".join(self.last_error.split("\n")[-num_lines:])
        if self.status == "running":
            self.title = "runner running"
            self.border = "gold1"
        elif last_error == "":
            self.title = "runner succeded"
            self.border = "green"
        else:
            self.title = "runner failed"
            self.border = "red"
        self.title = f"{self.title} [blue]({round(time.time() - self.time)}s)[/]"

        return f"runner is {self.status} {round(time.time() - self.time)}\npid: {self.proc.pid}\nhash: {self.m.content_dir_hash}\n{last_error}"

    def __rich__(self) -> Panel:

        if self.proc:
            if self.proc.poll() is None:
                return Panel(
                    Text(self.status_message),
                    border_style=self.border,
                    title=self.title,
                    expand=True,
                )

        if self.status == "running":
            self.status = "waiting"
            self.time = time.time()
            if self.proc:
                self.last_error = self.proc.stderr.read().decode()

        if self._dirhash != self.m.content_dir_hash:
            self.run()
            self._dirhash = self.m.content_dir_hash

        return Panel(
            Text(self.status_message),
            border_style=self.border,
            title=self.title,
            expand=True,
        )


if __name__ == "__main__":
    from rich.live import Live

    from markata import Markata

    from .cli import run_until_keyboard_interrupt

    with Live(Runner(Markata()), refresh_per_second=30, screen=True):
        run_until_keyboard_interrupt()
