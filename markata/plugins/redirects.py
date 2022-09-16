from dataclasses import asdict, dataclass
from pathlib import Path
from typing import TYPE_CHECKING

from jinja2 import Template

from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from markata import Markata

DEFAULT_REDIRECT_TEMPLATE = Path(__file__).parent / "default_redirect_template.html"


@dataclass
class Redirect:
    original: str
    new: str


@hook_impl
def save(markata: "Markata") -> None:
    assets_dir: str = str(markata.config.get("assets_dir", "static"))
    redirects_file = Path(
        str(markata.config.get("redirects", Path(assets_dir) / "_redirects"))
    )

    redirects = [
        Redirect(*s)
        for r in redirects_file.read_text().split("\n")
        if "*" not in r and len(s := r.split()) == 2 and not r.strip().startswith("#")
    ]

    if "redirect_template" in markata.config:
        template_file = Path(markata.config.get("redirect_template"))
    else:
        template_file = DEFAULT_REDIRECT_TEMPLATE
    template = Template(template_file.read_text())

    output_dir = Path(markata.config["output_dir"])  # type: ignore
    output_dir.mkdir(parents=True, exist_ok=True)

    for redirect in redirects:
        file = output_dir / redirect.original.strip("/") / "index.html"
        file.parent.mkdir(parents=True, exist_ok=True)
        file.write_text(template.render(**asdict(redirect), config=markata.config))
