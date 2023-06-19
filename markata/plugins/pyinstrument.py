"""
Markata plugin to create a pyinstrument profile if pyinstrument is installed.

The profile will be saved to <output_dir>/_profile/index.html
"""
from pathlib import Path
from typing import Any, Optional

import pydantic

from markata import Markata
from markata.hookspec import hook_impl, register_attr

try:
    from pyinstrument import Profiler

    SHOULD_PROFILE = True
except ModuleNotFoundError:
    SHOULD_PROFILE = False
    "ignore if pyinstrument does not exist"
    ...


class ProfilerConfig(pydantic.BaseModel):
    output_dir: pydantic.DirectoryPath = Path("markout")
    should_profile: bool = SHOULD_PROFILE
    profiler: Optional[
        Any
    ] = None  # No valicator for type pyinstrument.profiler.Profiler
    output_file: Optional[Path] = None

    @pydantic.validator("output_file", pre=True, always=True)
    def validate_output_file(cls, v, *, values):
        if v is None:
            output_file = values["output_dir"] / "_profile" / "index.html"
            output_file.parent.mkdir(parents=True, exist_ok=True)
            return output_file
        return v


class Config(pydantic.BaseModel):
    profiler: ProfilerConfig = ProfilerConfig()


@hook_impl()
@register_attr("config_models")
def config_model(markata: "Markata") -> None:
    markata.config_models.append(Config)


@hook_impl
def configure(markata: "Markata") -> None:
    # profiler must exist in the same thread and cannot be configured through pydantic validation
    if (
        markata.config.profiler.should_profile
        and markata.config.profiler.profiler is None
    ):
        markata.config.profiler.profiler = Profiler()
        markata.config.profiler.profiler.start()


@hook_impl(trylast=True)
def save(markata: Markata) -> None:
    "stop the profiler and save as late as possible"
    if markata.config.profiler.should_profile:
        if markata.config.profiler.profiler is not None:
            if markata.config.profiler.profiler.is_running:
                try:
                    markata.config.profiler.profiler.stop()
                    html = markata.config.profiler.profiler.output_html()
                    markata.config.profiler.output_file.write_text(html)
                    markata.console.print(
                        markata.config.profiler.profiler.output_text()
                    )

                except AttributeError:
                    markata.console.log(
                        "profiler not available, skipping save pyinstrument save",
                    )
                    markata.console.log(
                        "[red]to enable profiler [wheat1][itallic]pip install 'markata\[pyinstrument]'",
                    )


@hook_impl
def teardown(markata: Markata) -> None:
    "stop the profiler on exit"
    # import logging

    # logger = logging.getLogger()
    # logger.handlers.clear()
    if markata.config.profiler.should_profile:
        if markata.config.profiler.profiler is not None:
            if markata.config.profiler.profiler.is_running:
                markata.config.profiler.profiler.stop()
