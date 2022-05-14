from pathlib import Path
from typing import TYPE_CHECKING

from jinja2 import Template

from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from markata import Markata


@hook_impl
def render(markata: "Markata") -> None:
    if "post_template" in markata.config:
        template_file = markata.config["post_template"]
    else:
        template_file = Path(__file__).parent / "default_post_template.html"
    with open(template_file) as f:
        template = Template(f.read())
    for article in markata.iter_articles("apply template"):

        article.html = template.render(
            body=article.html,
            toc=markata.md.toc,  # type: ignore
            config=markata.config,
            **article,
        )
