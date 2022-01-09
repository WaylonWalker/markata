"""Default glob plugin"""
from pathlib import Path
from typing import TYPE_CHECKING, cast

import pytz
from feedgen.feed import FeedGenerator
from more_itertools import flatten

from markata import Markata
from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from markata import Markata

    class MarkataRss(Markata):
        fg: "FeedGenerator"
        rss: str


@hook_impl(trylast=True)
def render(markata: "MarkataRss") -> None:
    fg = FeedGenerator()
    fg.id(markata.config["url"] + "/rss.xml")
    fg.title(markata.config["title"])
    fg.author(
        {"name": markata.config["author_name"], "email": markata.config["author_email"]}
    )
    fg.link(href=markata.config["url"], rel="alternate")
    fg.logo(markata.config["icon"])
    fg.subtitle(markata.config["rss_description"])
    fg.link(href=markata.config["url"] + "/rss.xml", rel="self")
    fg.language(markata.config["lang"])

    try:
        all_posts = reversed(sorted(markata.articles, key=lambda x: x["date"]))
        posts = [post for post in all_posts if post["status"] == "published"]
    except BaseException:
        posts = markata.articles

    for article in posts:
        fe = fg.add_entry()
        fe.id(markata.config["url"] + "/" + article["slug"])
        fe.title(article.metadata["title"])
        fe.published(article.metadata["datetime"])
        fe.description(article.metadata["description"])
        fe.summary(article.metadata["long_description"])
        fe.link(href=markata.config["url"] + "/" + article["slug"])
        fe.content(article.article_html.translate(dict.fromkeys(range(32))))

    markata.fg = fg
    markata.rss = fg.rss_str(pretty=True)


@hook_impl
def save(markata: "MarkataRss") -> None:
    output_dir = Path(markata.config["output_dir"])
    markata.fg.rss_file(str(output_dir / "rss.xml"))
