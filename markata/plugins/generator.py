"""add generator meta tag"""
from markata.hookspec import hook_impl
from bs4 import BeautifulSoup
from markata import Markata, __version__


@hook_impl(trylast=True)
def render(markata: Markata) -> None:
    for article in markata.iter_articles("add seo tags"):
        soup = BeautifulSoup(article.html, features="lxml")
        tag = soup.new_tag("meta")
        tag.attrs["content"] = f"markata {__version__}"
        tag.attrs["name"] = "generator"
        soup.head.append(tag)
        article.soup = soup
        article.html = soup.prettify()
