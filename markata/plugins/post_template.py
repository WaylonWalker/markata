from typing import TYPE_CHECKING

from jinja2 import Template

from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from markata import Markata


@hook_impl
def render(markata: "Markata") -> None:
    template_file = markata.post_template
    with open(template_file) as f:
        template = Template(f.read())
    for article in markata.iter_articles("apply template"):

        article.html = template.render(
            body=article.html,
            toc=markata.md.toc,  # type: ignore
            **article,
        )
