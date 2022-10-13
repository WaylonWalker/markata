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
        self.m = markata
        self.run()

    def run(self) -> None:
        self.status = "running"
        self.proc = subprocess.Popen(
            ["markata", "build"],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
        )
        self.time = time.time()

    @property
    def status_message(self) -> str:
        num_lines = self.m.console.height - 4
        last_error = "\n".join(self.last_error.split("\n")[-num_lines:])
        if last_error == "":
            self.title = "runner success"
            self.border = "green"
        else:
            self.title = "runner failed"
            self.border = "red"

        return f"runner is {self.status} {round(time.time() - self.time)}\nhash: {self.m.content_dir_hash}\n{last_error}"

    def __rich__(self) -> Panel:

        if self.proc.poll() is None:
            self.title = "runner running"
            self.border = "gold1"
            return Panel(
                Text(self.status_message),
                border_style=self.border,
                title=self.title,
                expand=True,
            )

        if self.status == "running":
            self.status = "waiting"
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
