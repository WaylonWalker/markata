---
date: 2025-12-09
description: None
published: false
slug: markata/cli/runner
title: runner.py


---

---

None

---

!!! class
    <h2 id="Runner" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">Runner <em class="small">class</em></h2>

    Display Footer

???+ source "Runner <em class='small'>source</em>"
    ```python
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
    ```
!!! method
    <h2 id="run" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">run <em class="small">method</em></h2>

    Runs the build only if one is not already running.

???+ source "run <em class='small'>source</em>"
    ```python
    def run(self) -> None:
            "Runs the build only if one is not already running."
            if self.proc.poll() is not None:
                self._run()
    ```
!!! method
    <h2 id="_run" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">_run <em class="small">method</em></h2>

    Runs the build and sets the proc

???+ source "_run <em class='small'>source</em>"
    ```python
    def _run(self) -> None:
            "Runs the build and sets the proc"
            self.status = "running"
            self.time = time.time()
            self.proc = subprocess.Popen(
                ["markata", "build"],
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
            )
    ```
!!! method
    <h2 id="status_message" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">status_message <em class="small">method</em></h2>

    returns the status message to display

???+ source "status_message <em class='small'>source</em>"
    ```python
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
    ```