import subprocess
import time
from multiprocessing import Process

from rich.panel import Panel
from rich.text import Text

from markata import Markata


class Runner:
    """Display Footer"""

    _status = "waiting"

    m = Markata()
    m.console.quiet = True
    _dirhash = ""
    time = time.time()
    std = ""

    def run(self) -> None:
        self.cmd = ["markata"]
        import os

        self.time = time.time()

        with open("markata.log", "w", 1) as f:
            self.proc = subprocess.Popen(
                self.cmd,
                cwd=os.getcwd(),
                stderr=subprocess.PIPE,
                stdout=f,
            )

    @property
    def is_running(self):
        return self.proc.poll() == None

    @property
    def status(self):
        if self._status == "running":
            if not self.is_running:
                self.status = "waiting"
                self.time = time.time()

        elif self.is_running:
            self.status = "running"
        return self._status

    @status.setter
    def status(self, value):
        if value not in ["running", "waiting"]:
            raise ValueError(f"{value} is not a valid state")
        self._status = value

    @property
    def color(self):
        if self.status == "running":
            return "green"
        return "white"

    @property
    def phase(self):
        return self.m.phase_file.read_text()

    def __rich__(self) -> Panel:

        if self._dirhash != self.m.content_dir_hash:
            self.run()
            self._dirhash = self.m.content_dir_hash

        if self.status == "running":
            s = f"{self.status} {self.proc.pid} {self.phase} {round(time.time() - self.time)}"

        else:
            s = f"{self.status} {self.phase} {round(time.time() - self.time)}"
        return Panel(Text(s + self.std), border_style=self.color)


if __name__ == "__main__":
    from rich.live import Live

    from .cli import run_until_keyboard_interrupt

    with Live(Runner(), refresh_per_second=30, screen=True):
        run_until_keyboard_interrupt()
