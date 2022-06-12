"""
Setup Logging hook sets up the RichHandler for pretty console logs, and file logs to the configured markata's configured `log_dir`, or
`output_dir/_logs` if `log_dir` is not configured.  The log file will be named after the `<levelname>.log`

## The log files

There will be 6 log files created based on log level and file type.

```
markout/_logs
├── debug
│   └── index.html
├── debug.log
├── info
│   └── index.html
├── info.log
├── warning
│   └── index.html
└── warning.log
```

## Configuration

Ensure that setup_logging is in your hooks.  You can check if `setup_logging`
is in your hooks by running `markata list --hooks` from your terminal and
checking the output, or creating an instance of `Markata()` and checking the
`Markata().hooks` attribute.  If its missing or you wan to be more explicit,
you can add `setup_logging` to your `markata.toml` `[markata.hooks]`.

``` toml
[markata]

# make sure its in your list of hooks
hooks=[
   "markata.plugins.setup_logging",
   ]
```

## Log Template
``` toml
[markata]

# make sure its in your list of hooks
hooks=[
   "markata.plugins.setup_logging",
   ]

# point log template to the path of your logging template
log_template='templates/log_template.html'
```

You can see the latest default `log_template` on
[GitHub](https://github.com/WaylonWalker/markata/blob/main/markata/plugins/default_log_template.html)

## Disable Logging

If you do not want logging, you can explicityly disable it by adding it to your
`[markata.disabled_hooks]` array in your `[markata.toml]`

``` toml
[markata]

# make sure its in your list of hooks
disabled_hooks=[
   "markata.plugins.setup_logging",
   ]
```

"""
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


def has_rich_handler() -> bool:
    """
    Returns a boolean whether or not there is a RichHandler attached to the
    root logger.
    """

    logger = logging.getLogger()
    return bool([h for h in logger.handlers if isinstance(h, RichHandler)])


def has_file_handler(log_file: Path) -> bool:
    logger = logging.getLogger()
    existing_logger_files = [
        handler.baseFilename
        for handler in logger.handlers
        if isinstance(handler, logging.FileHandler)
    ]
    return str(log_file.absolute()) in existing_logger_files


def setup_log(markata: "Markata", level: int = logging.INFO) -> Path:
    path = setup_html_log(markata, level)
    setup_text_log(markata, level)
    return path


def setup_text_log(markata: "Markata", level: int = logging.INFO) -> Path:
    """
    sets up a plain text log in markata's configured `log_dir`, or
    `output_dir/_logs` if `log_dir` is not configured.  The log file will be
    named after the `<levelname>.log`
    """
    log_file = Path(
        str(
            markata.config.get(
                "log_dir",
                Path(str(markata.config.get("output_dir", "markout")))
                / "_logs"
                / (logging.getLevelName(level).lower() + ".log"),
            )
        )
    )
    if has_file_handler(log_file):
        return log_file

    if not log_file.parent.exists():
        log_file.parent.mkdir(parents=True)
    fh = logging.FileHandler(log_file)
    fh.setLevel(level)
    fh_formatter = logging.Formatter(
        "%(asctime)s %(name)-12s %(levelname)-8s %(message)s",
    )
    fh.setFormatter(fh_formatter)
    logging.getLogger("").addHandler(fh)

    return log_file


def setup_html_log(markata: "Markata", level: int = logging.INFO) -> Path:
    """
    sets up an html log in markata's configured `log_dir`, or
    `output_dir/_logs` if `log_dir` is not configured.  The log file will be
    named after the `<levelname>/index.html`.  The goal of this is to give
    """

    log_file = Path(
        str(
            markata.config.get(
                "log_dir",
                Path(str(markata.config.get("output_dir", "markout")))
                / "_logs"
                / logging.getLevelName(level).lower()
                / "index.html",
            )
        )
    )

    if has_file_handler(log_file):
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
    setup_log(markata, logging.DEBUG)
    setup_log(markata, logging.INFO)
    setup_log(markata, logging.WARNING)

    if not has_rich_handler():
        console = RichHandler(rich_tracebacks=True)
        console.setLevel(logging.INFO)
        formatter = logging.Formatter("%(message)s")
        console.setFormatter(formatter)
        logging.getLogger("").addHandler(console)
