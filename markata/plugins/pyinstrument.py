"""
Markata plugin to create a pyinstrument profile if pyinstrument is installed.

The profile will be saved to <output_dir>/_profile/index.html
"""
from pathlib import Path
import typer
from typing import TYPE_CHECKING

from markata.hookspec import hook_impl, register_attr

if TYPE_CHECKING:
    from markata import Markata

try:
    from pyinstrument import Profiler
except ModuleNotFoundError:
    "ignore if pyinstrument does not exist"
    ...


@hook_impl(tryfirst=True)
def cli(app, markata: "Markata") -> None:
    def should_profile_callback(value: bool):
        if value:
            markata.should_profile = True

    @app.callback()
    def should_profile(
        version: bool = typer.Option(
            None, "--profile", callback=should_profile_callback, is_eager=True
        ),
    ):
        # Do other global stuff, handle other global options here
        return

    @app.command()
    def profile():
        ...


@hook_impl(tryfirst=True)
@register_attr("should_profile")
def configure(markata: "Markata") -> None:
    "set the should_profile variable"

    try:
        markata.should_profile = markata.config["pyinstrument"]["should_profile"]
    except KeyError:
        markata.should_profile = False


@hook_impl(tryfirst=True)
@register_attr("profiler")
def glob(markata: "Markata") -> None:
    "start the profiler as soon as possible"
    if markata.should_profile:
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
