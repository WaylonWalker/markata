"""manifest plugin"""
import json
from pathlib import Path
from typing import TYPE_CHECKING, List

from bs4 import BeautifulSoup

from markata import Markata, __version__
from markata.hookspec import hook_impl

if TYPE_CHECKING:
    import frontmatter
    from bs4.element import Tag


def _create_seo(
    markata: Markata, soup: BeautifulSoup, article: "frontmatter.Post"
) -> List:
    if article.metadata["description"] == "" or None:
        article.metadata["description"] = " ".join(
            [p.text for p in soup.find(id="post-body").find_all("p")]
        ).strip()[:120]

    seo = [
        *markata.seo,
        {
            "name": "description",
            "property": "description",
            "content": article.metadata["description"],
        },
        {
            "name": "og:description",
            "property": "og:description",
            "content": article.metadata["description"],
        },
        {
            "name": "twitter:description",
            "property": "twitter:description",
            "content": article.metadata["description"],
        },
        {
            "name": "og:title",
            "property": "og:title",
            "content": f'{article.metadata["title"]} | {markata.config["site_name"]}'[
                :60
            ],
        },
        {
            "name": "twitter:title",
            "property": "twitter:title",
            "content": f'{article.metadata["title"]} | {markata.config["site_name"]}'[
                :60
            ],
        },
        {
            "name": "og:image",
            "property": "og:image",
            "content": f'{markata.config["url"]}/{article.metadata["slug"]}-og.png',
        },
        {
            "name": "twitter:image",
            "property": "twitter:image",
            "content": f'{markata.config["url"]}/{article.metadata["slug"]}-og.png',
        },
        {
            "name": "og:image:width",
            "property": "og:image:width",
            "content": "1600",
        },
        {
            "name": "og:image:width",
            "property": "og:image:width",
            "content": "900",
        },
        {
            "name": "twitter:card",
            "property": "twitter:card",
            "content": markata.config["twitter_card"],
        },
        {
            "name": "og:site_name",
            "property": "og:site_name",
            "content": markata.config["site_name"],
        },
        {
            "name": "twitter:creator",
            "property": "twitter:creator",
            "content": markata.config["twitter_creator"],
        },
        {
            "name": "title",
            "property": "title",
            "content": article.metadata["title"],
        },
        {
            "name": "generator",
            "property": "generator",
            "content": f"markata {__version__}",
        },
    ]
    return seo


def _add_seo_tags(seo: List, article: "frontmatter.Post", soup: BeautifulSoup) -> None:
    for meta in seo:
        soup.head.append(_create_seo_tag(meta, soup))


def _create_seo_tag(meta: dict, soup: BeautifulSoup) -> "Tag":
    tag = soup.new_tag("meta")
    for k in meta:
        tag.attrs[k] = meta[k]
    return tag


@hook_impl
def render(markata: Markata) -> None:
    for article in markata.iter_articles("add seo tags from seo.py"):
        key = markata.make_hash(
            "seo",
            "render",
            article["content_hash"],
            markata.site_name,
            markata.url,
            article.metadata["slug"],
            markata.twitter_card,
            article.metadata["title"],
            markata.site_name,
            str(markata.seo),
        )
        html_from_cache = markata.cache.get(key)

        if html_from_cache is None:
            soup = BeautifulSoup(article.html, features="lxml")
            seo = _create_seo(markata, soup, article)
            _add_seo_tags(seo, article, soup)
            canonical_link = soup.new_tag("link")
            canonical_link.attrs["rel"] = "canonical"
            canonical_link.attrs["href"] = f'{markata.url}/{article.metadata["slug"]}/'
            soup.head.append(canonical_link)

            html = soup.prettify()
            markata.cache.add(key, html, expire=15 * 24 * 60)
        else:
            html = html_from_cache
        article.html = html
