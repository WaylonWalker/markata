import datetime
import logging
import sys
from pathlib import Path
from typing import TYPE_CHECKING

from jinja2 import Template, Undefined
from rich.logging import RichHandler

from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from markata import Markata


class SilentUndefined(Undefined):
    def _fail_with_undefined_error(self, *args, **kwargs):
        return ""


def setup_log(markata: "Markata", level: int = logging.INFO) -> Path:

    log_file = Path(
        str(
            markata.config.get(
                "log_dir",
                Path(str(markata.config.get("output_dir", "markout")))
                / "_logs"
                / logging.getLevelName(level)
                / "index.html",
            )
        )
    )

    logger = logging.getLogger()
    files = [
        l.baseFilename for l in logger.handlers if isinstance(l, logging.FileHandler)
    ]
    if str(log_file.absolute()) in files:
        return log_file

    if not log_file.parent.exists():
        log_file.parent.mkdir(parents=True)
    if not log_file.exists():
        template_file = Path(
            str(
                markata.config.get(
                    "log_template", Path(__file__).parent / "default_log_template.html"
                )
            )
        )
        template = Template(template_file.read_text(), undefined=SilentUndefined)
        log_header = template.render(
            title=str(markata.config.get("title", "markata build")) + " logs",
            config=markata.config,
        )
        log_file.write_text(log_header)
    with open(log_file, "a") as f:
        command = Path(sys.argv[0]).name + " " + " ".join(sys.argv[1:])
        f.write(
            f"""
            <div style="width: 100%; height: 20px; margin-top: 5rem; border-bottom: 1px solid goldenrod; text-align: center">
    <span style="padding: 0 10px;">
    {datetime.datetime.now()} running "{command}"
    </span>
    </div>
    """
        )
    fh = logging.FileHandler(log_file)
    fh.setLevel(level)
    fh_formatter = logging.Formatter(
        "<li><p><span class='time'>%(asctime)s</span> <span class='name %(name)s'>%(name)-12s</span> <span class='levelname %(levelname)s'>%(levelname)-8s</span></p><p class='message'>%(message)s</p></li>"
    )
    fh.setFormatter(fh_formatter)
    logging.getLogger("").addHandler(fh)

    return log_file


@hook_impl(tryfirst=True)
def configure(markata: "Markata") -> None:

    logging.basicConfig(
        level=logging.DEBUG,
        format="%(asctime)s %(name)-12s %(levelname)-8s %(message)s",
        datefmt="%m-%d %H:%M",
        filename="debug.log",
        filemode="w",
    )

    setup_log(markata, logging.DEBUG)
    setup_log(markata, logging.INFO)
    setup_log(markata, logging.WARNING)

    console = RichHandler(rich_tracebacks=True)
    console.setLevel(logging.INFO)
    formatter = logging.Formatter("%(message)s")
    console.setFormatter(formatter)
    logging.getLogger("").addHandler(console)
