"""manifest plugin"""
from markata.hookspec import hook_impl
from pathlib import Path
import json
from bs4 import BeautifulSoup


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
def render(markata):
    for article in markata.iter_articles("add seo tags"):
        soup = BeautifulSoup(article.html, features="lxml")
        if article.metadata["description"] == "":
            article.metadata["description"] = " ".join(
                [p.text for p in soup.find(id="post-body").find_all("p")]
            ).strip()[:120]

        if "long_description" not in article.metadata:
            article.metadata["long_description"] = ""

        if article.metadata["long_description"] == "":
            article.metadata["long_description"] = " ".join(
                [p.text for p in soup.find(id="post-body").find_all("p")]
            ).strip()[:250]

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
        ]
        _add_seo_tags(seo, article, soup)

        article.soup = soup
        article.html = soup.prettify()
