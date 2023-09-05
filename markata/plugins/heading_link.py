"""
Creates links next to all heading tags to make it easier for users to share a
specific heading.
"""
from pathlib import Path
import re
from typing import TYPE_CHECKING

from bs4 import BeautifulSoup

from markata import Markata
from markata.hookspec import hook_impl

if TYPE_CHECKING:
    import bs4
    from frontmatter import Post


@hook_impl(trylast=True)
def post_render(markata: Markata) -> None:
    """
    This plugin creates a link svg next to all headings.
    """

    with markata.cache as cache:
        for article in markata.iter_articles("link headers"):
            key = markata.make_hash(
                "heading_link",
                "post_render",
                Path(__file__).read_text(),
                article.content,
                article.html,
            )

            html_from_cache = markata.precache.get(key)

            if html_from_cache is None:
                html = link_headings(article)
                cache.add(
                    key,
                    html,
                    expire=markata.config.default_cache_expire,
                )
            else:
                html = html_from_cache
            article.html = html


def link_headings(article: "Post") -> str:
    """
    Use BeautifulSoup to find all headings and run link_heading on them.
    """
    soup = BeautifulSoup(article.html, "html.parser")
    for heading in soup.find_all(re.compile("^h[1-6]$")):
        if (
            not heading.find("a", {"class": "heading-permalink"})
            and heading.get("id", "") != "title"
        ):
            link_heading(soup, heading)
    return str(soup)


def link_heading(soup: "bs4.BeautifulSoup", heading: "bs4.element.Tag") -> None:
    """
    Mutate soup to include an svg link at the heading passed in.
    """
    id = heading.get("id")

    link = soup.new_tag(
        "a",
        alt="id",
        title=f"link to #{id}",
        href=f"#{id}",
        **{"class": "heading-permalink"},
    )
    span = soup.new_tag("span", **{"class": "visually-hidden"})
    svg = soup.new_tag(
        "svg",
        fill="currentColor",
        focusable="false",
        width="1em",
        height="1em",
        xmlns="http://www.w3.org/2000/svg",
        viewBox="0 0 24 24",
        **{
            "aria-hidden": "true",
        },
    )

    path = soup.new_tag(
        "path",
        d="M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992 5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242 6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985 3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005 3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201 13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995 2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0 0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836 19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0 0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632 1.563 1 1 0 0 0 1.602-1.198z",
    )
    svg.append(path)
    link.append(span)
    link.append(svg)
    heading.append(link)
