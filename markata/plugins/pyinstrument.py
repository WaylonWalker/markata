"""
Markata plugin to create a pyinstrument profile if pyinstrument is installed.

The profile will be saved to <output_dir>/_profile/index.html
"""
from pathlib import Path

from markata.hookspec import hook_impl

try:
    from pyinstrument import Profiler
except ModuleNotFoundError:
    "ignore if pyinstrument does not exist"
    ...


@hook_impl(tryfirst=True)
def glob(markata: "Markata") -> None:
    "start the profiler as soon as possible"
    try:
        markata.profiler = Profiler()
        markata.profiler.start()
    except NameError:
        "ignore if Profiler does not exist"
        ...


@hook_impl(trylast=True)
def save(markata: "Markata") -> None:
    "stop the profiler and save as late as possible"
    try:
        output_file = Path(markata.output_dir) / "_profile" / "index.html"
        output_file.parent.mkdir(parents=True, exist_ok=True)
        markata.profiler.stop()
        html = markata.profiler.output_html()
        output_file.write_text(html)
        markata.console.print(markata.profiler.output_text())

    except AttributeError:
        "ignore if markata does not have a profiler attribute"
        ...
