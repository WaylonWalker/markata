import os
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

    _dirhash = ""
    time = time.time()
    std = ""

    def __init__(self, markata: "Markata") -> None:
        self.m = markata

    def run(self) -> None:
        self.cmd = ["markata", "build"]

        with open("markata.log", "w", 1) as f:
            self.proc = subprocess.Popen(
                self.cmd,
                cwd=os.getcwd(),
                stderr=subprocess.PIPE,
                stdout=f,
            )

    @property
    def is_running(self) -> bool:
        return self.proc.poll() is None

    @property
    def status(self) -> str:
        if self._status == "running":
            if not self.is_running:
                self.status = "waiting"
                self.time = time.time()

        elif self.is_running:
            self.status = "running"
        return self._status

    @status.setter
    def status(self, value) -> None:
        if value not in ["running", "waiting"]:
            raise ValueError(f"{value} is not a valid state")
        self._status = value

    @property
    def color(self) -> str:
        if self.status == "running":
            return "green"
        return "white"

    @property
    def phase(self) -> str:
        return self.m.phase_file.read_text()

    def __rich__(self) -> Panel:

        if self._dirhash != self.m.content_dir_hash:
            self.run()
            self._dirhash = self.m.content_dir_hash

        if self.status == "running":
            s = f"{self.status} {self.proc.pid} {self.phase} {round(time.time() - self.time)}"

        else:
            s = f"{self.status} {self.phase} {round(time.time() - self.time)}\n {self.m.content_dir_hash}"
        return Panel(Text(s + self.std), border_style=self.color, title="runner")


if __name__ == "__main__":
    from rich.live import Live

    from markata import Markata

    from .cli import run_until_keyboard_interrupt

    with Live(Runner(Markata()), refresh_per_second=30, screen=True):
        run_until_keyboard_interrupt()
