import atexit
import subprocess
import time
from typing import TYPE_CHECKING

from rich.panel import Panel
from rich.text import Text

from markata.hookspec import hook_impl, register_attr

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
        self._dirhash = self.m.content_dir_hash
        self._run()
        atexit.register(self.kill)

    def kill(self) -> None:
        self.proc.stdout.close()
        self.proc.stderr.close()
        self.proc.kill()
        self.proc.wait()

    def run(self) -> None:
        "Runs the build only if one is not already running."
        if self.proc.poll() is not None:
            self._run()

    def _run(self) -> None:
        "Runs the build and sets the proc"
        self.status = "running"
        self.time = time.time()
        self.proc = subprocess.Popen(
            ["markata", "build"],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
        )

    @property
    def status_message(self) -> str:
        "returns the status message to display"
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

        return (
            f"runner is {self.status}"
            f"{round(time.time() - self.time)}\n"
            f"pid: {self.proc.pid}\n"
            f"hash: {self.m.content_dir_hash}\n"
            f"{last_error}"
        )

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


@hook_impl
@register_attr("runner")
def configure(markata: "Markata") -> None:
    def get_runner(self):
        try:
            return self._runner
        except AttributeError:
            self._runner: Runner = Runner(self)
            return self._runner

    from markata import Markata

    Markata.runner = property(get_runner)


if __name__ == "__main__":
    from rich.live import Live

    from .cli import run_until_keyboard_interrupt

    with Live(Runner(Markata()), refresh_per_second=30, screen=True):
        run_until_keyboard_interrupt()
