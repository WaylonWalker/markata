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

    def __init__(self, markata: "Markata") -> None:
        self.m = markata

    def run(self) -> None:
        self.status = "running"
        self.m.run()
        self.time = time.time()
        self.status = "waiting"

    def __rich__(self) -> Panel:

        if self._dirhash != self.m.content_dir_hash:
            self.run()
            self._dirhash = self.m.content_dir_hash

        s = f"runner is waiting {round(time.time() - self.time)}"

        return Panel(Text(s), border_style="green", title="runner")


if __name__ == "__main__":
    from rich.live import Live

    from .cli import run_until_keyboard_interrupt

    with Live(Runner(Markata()), refresh_per_second=30, screen=True):
        run_until_keyboard_interrupt()
