"""manifest plugin"""
import json
from pathlib import Path

from bs4 import BeautifulSoup

from markata.hookspec import hook_impl
from markata import Markata, __version__


def _create_seo(markata: Markata, soup, article):
    if article.metadata["description"] == "":
        article.metadata["description"] = " ".join(
            [p.text for p in soup.find(id="post-body").find_all("p")]
        ).strip()[:120]

    seo = [
        *markata.config["seo"],
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


def _add_seo_tags(seo, article, soup):
    for meta in seo:
        soup.head.append(_create_seo_tag(meta, soup))
    return


def _create_seo_tag(meta: dict, soup):
    tag = soup.new_tag("meta")
    for k in meta:
        tag.attrs[k] = meta[k]
    return tag


@hook_impl
def render(markata: Markata):
    for article in markata.iter_articles("add seo tags from seo.py"):
        key = markata.make_hash(
            "seo",
            "render",
            article["content_hash"],
            markata.config["site_name"],
            markata.config["url"],
            article.metadata["slug"],
            markata.config["twitter_card"],
            article.metadata["title"],
            markata.config["site_name"],
            str(markata.config["seo"]),
        )
        html_from_cache = markata.cache.get(key)

        if html_from_cache is None:
            soup = BeautifulSoup(article.html, features="lxml")
            seo = _create_seo(markata, soup, article)
            _add_seo_tags(seo, article, soup)
            html = soup.prettify()
            markata.cache.add(key, html, expire=15 * 24 * 60)
        else:
            html = html_from_cache
        article.html = html
