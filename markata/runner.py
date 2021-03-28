import time
from checksumdir import dirhash
from multiprocessing import Process
from rich.panel import Panel
from rich.live import Live
from rich.text import Text
from rich.layout import Layout
from markata import Markata


def run_until_keyboard_interrupt():
    try:
        while True:
            time.sleep(0.2)
    except KeyboardInterrupt:
        pass


class Runner:
    """Display Footer"""

    _status = "waiting"

    m = Markata()
    m.console.quiet = True
    _dirhash = ""
    time = time.time()

    def _run(self) -> None:
        self.m.run()
        # time.sleep(1)

    def run(self) -> None:
        if not self.is_running:
            self.p = Process(target=self._run)
            self.p.start()

    @property
    def is_running(self):
        try:
            return self.p.is_alive()
        except AttributeError:
            return False

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

    def __rich__(self) -> Panel:
        from random import choice

        if self._dirhash != self.m.content_dir_hash:
            self.run()
            self._dirhash = self.m.content_dir_hash

        if self.status == "running":
            s = f"{self.status} {self.p.pid} {self.m.phase} {round(time.time() - self.time)}"

        else:
            s = f"{self.status} {self.m.phase} {round(time.time() - self.time)}"
        return Panel(Text(s), border_style=self.color)


def main():
    layout = Layout()
    layout.update(Runner())

    with Live(layout, refresh_per_second=4, screen=True):
        # time.sleep(9999)
        run_until_keyboard_interrupt()


if __name__ == "__main__":
    main()
