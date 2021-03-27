"""Default glob plugin"""
from pathlib import Path
from typing import TYPE_CHECKING

import pytz
from feedgen.feed import FeedGenerator
from more_itertools import flatten

from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from markata import Markata

    class MarkataRss(Markata):
        fg: "FeedGenerator"
        rss: str


@hook_impl(trylast=True)
def render(markata: "MarkataRss") -> None:
    fg = FeedGenerator()
    fg.id(markata.url + "/rss.xml")
    fg.title(markata.title)
    fg.author({"name": markata.author_name, "email": markata.author_email})
    fg.link(href=markata.url, rel="alternate")
    fg.logo(markata.icon)
    fg.subtitle(markata.rss_description)
    fg.link(href=markata.url + "/rss.xml", rel="self")
    fg.language(markata.lang)

    for article in markata.articles:
        fe = fg.add_entry()
        fe.id(markata.url + "/" + article["slug"])
        fe.title(article.metadata["title"])
        fe.description(article.metadata["description"])
        fe.link(href=markata.url + "/" + article["slug"])
        fe.content(article.article_html.translate(dict.fromkeys(range(32))))

    markata.fg = fg
    markata.rss = fg.rss_str(pretty=True)


from typing import cast

from markata import Markata


@hook_impl
def save(markata: "MarkataRss") -> None:
    output_dir = Path(markata.output_dir)
    markata.fg.rss_file(str(output_dir / "rss.xml"))
