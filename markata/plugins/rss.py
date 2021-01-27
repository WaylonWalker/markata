"""Default glob plugin"""
import pytz
from markata.hookspec import hook_impl
from pathlib import Path
from more_itertools import flatten
from feedgen.feed import FeedGenerator


@hook_impl(trylast=True)
def render(markata):
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

    for article in markata.articles:
        fe = fg.add_entry()
        fe.id(markata.config["url"] + "/" + article["slug"])
        fe.title(article.metadata["title"])
        fe.description(article.metadata["description"])
        fe.link(href=markata.config["url"] + "/" + article["slug"])
        fe.content(article.article_html.translate(dict.fromkeys(range(32))))
        # fe.published(str(article["datetime"]))

    markata.fg = fg
    markata.rss = fg.rss_str(pretty=True)
    # markata.atom = fg.atom_str(pretty=True)


@hook_impl
def save(markata):
    output_dir = Path(markata.config["output_dir"])
    markata.fg.rss_file(str(output_dir / "rss.xml"))
