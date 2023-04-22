"""add generator meta tag"""
from bs4 import BeautifulSoup

from markata import Markata, __version__
from markata.hookspec import hook_impl


@hook_impl(trylast=True)
def render(markata: Markata) -> None:
    should_prettify = markata.config.prettify_html
    for article in markata.iter_articles("add ssg tag"):
        soup = BeautifulSoup(article.html, features="lxml")
        tag = soup.new_tag("meta")
        tag.attrs["content"] = f"markata {__version__}"
        tag.attrs["name"] = "generator"
        soup.head.append(tag)

        if should_prettify:
            article.html = soup.prettify()
        else:
            article.html = str(soup)
