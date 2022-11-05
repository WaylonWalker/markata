"""
Markata plugin to create a pyinstrument profile if pyinstrument is installed.

The profile will be saved to <output_dir>/_profile/index.html
"""
from pathlib import Path

from markata import Markata
from markata.hookspec import hook_impl, register_attr

try:
    from pyinstrument import Profiler
except ModuleNotFoundError:
    "ignore if pyinstrument does not exist"
    ...


class MarkataInstrument(Markata):
    should_profile = False
    profiler = None


@hook_impl(tryfirst=True)
@register_attr("should_profile", "profiler")
def configure(markata: MarkataInstrument) -> None:
    "set the should_profile variable"
    markata.profiler = None

    if "should_profile" not in markata.__dict__.keys():
        markata.should_profile = markata.config.get("pyinstrument", {}).get(
            "should_profile", True
        )

    if markata.should_profile and "profiler" not in markata.__dict__.keys():
        try:
            markata.profiler = Profiler(async_mode="disabled")
            markata.profiler.start()
        except NameError:
            "ignore if Profiler does not exist"
            ...


@hook_impl(trylast=True)
def save(markata: MarkataInstrument) -> None:
    "stop the profiler and save as late as possible"
    if markata.should_profile:
        try:
            if "profiler" in markata.__dict__.keys():
                output_file = (
                    Path(markata.config["output_dir"]) / "_profile" / "index.html"
                )
                output_file.parent.mkdir(parents=True, exist_ok=True)
                markata.profiler.stop()
                html = markata.profiler.output_html()
                output_file.write_text(html)
                markata.console.print(markata.profiler.output_text())

        except AttributeError:
            "ignore if markata does not have a profiler attribute"
            ...


@hook_impl
def teardown(markata: MarkataInstrument) -> None:
    "stop the profiler on exit"

    if markata.should_profile:
        if hasattr(markata, "profiler"):
            if markata.profiler is not None:
                if markata.profiler.is_running:
                    markata.profiler.stop()
