"""manifest plugin"""
import logging
from typing import Any, Dict, List, TYPE_CHECKING

from bs4 import BeautifulSoup

from markata import Markata, __version__
from markata.hookspec import hook_impl

if TYPE_CHECKING:
    import frontmatter
    from bs4.element import Tag

logger = logging.getLogger(__file__)


def _create_seo(
    markata: Markata,
    soup: BeautifulSoup,
    article: "frontmatter.Post",
    site_name: str,
    author_name: str,
    author_email: str,
    twitter_card: str,
    twitter_creator: str,
    config_seo: Dict,
    images_url: str,
) -> List:
    if article.metadata["description"] == "" or None:
        try:
            article.metadata["description"] = " ".join(
                [p.text for p in soup.find(id="post-body").find_all("p")],
            ).strip()[:120]
        except AttributeError:
            article.metadata["description"] = ""

    seo = [
        *config_seo,
        {
            "name": "og:author",
            "property": "og:author",
            "content": author_name,
        },
        {
            "name": "og:author_email",
            "property": "og:author_email",
            "content": author_email,
        },
        {
            "name": "og:type",
            "property": "og:type",
            "content": "website",
        },
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
            "content": f'{article.metadata["title"]} | {site_name}'[:60],
        },
        {
            "name": "twitter:title",
            "property": "twitter:title",
            "content": f'{article.metadata["title"]} | {site_name}'[:60],
        },
        {
            "name": "og:image",
            "property": "og:image",
            "content": f'{images_url}/{article.metadata["slug"]}-og.png',
        },
        {
            "name": "twitter:image",
            "property": "twitter:image",
            "content": f'{images_url}/{article.metadata["slug"]}-og.png',
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
            "content": twitter_card,
        },
        {
            "name": "og:site_name",
            "property": "og:site_name",
            "content": site_name,
        },
        {
            "name": "twitter:creator",
            "property": "twitter:creator",
            "content": twitter_creator,
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


def _get_or_warn(config: Dict, key: str, default: str) -> Any:
    if key not in config.keys():
        logger.warning(
            f"{key} is missing from markata.toml config, using default value {default}",
        )
    return config.get(key, default)


@hook_impl
def render(markata: Markata) -> None:
    url = _get_or_warn(markata.config, "url", "")
    images_url = markata.config.get("images_url", url)
    site_name = _get_or_warn(markata.config, "site_name", "")
    author_name = _get_or_warn(markata.config, "author_name", "")
    author_email = _get_or_warn(markata.config, "author_email", "")
    twitter_creator = _get_or_warn(markata.config, "twitter_creator", "")
    twitter_card = _get_or_warn(markata.config, "twitter_card", "summary_large_image")
    config_seo = markata.config.get("seo", {})
    should_prettify = markata.config.get("prettify_html", False)

    with markata.cache as cache:
        for article in markata.iter_articles("add seo tags from seo.py"):
            key = markata.make_hash(
                "seo",
                "render",
                article.html,
                site_name,
                url,
                article.metadata["slug"],
                twitter_card,
                article.metadata["title"],
                str(config_seo),
            )

            html_from_cache = markata.precache.get(key)

            if html_from_cache is None:
                soup = BeautifulSoup(article.html, features="lxml")
                seo = _create_seo(
                    markata=markata,
                    soup=soup,
                    article=article,
                    site_name=site_name,
                    author_name=author_name,
                    author_email=author_email,
                    twitter_card=twitter_card,
                    twitter_creator=twitter_creator,
                    config_seo=config_seo,
                    images_url=images_url,
                )
                _add_seo_tags(seo, article, soup)
                canonical_link = soup.new_tag("link")
                canonical_link.attrs["rel"] = "canonical"
                if article.metadata["slug"] == "index":
                    canonical_link.attrs["href"] = f"{url}/"
                else:
                    canonical_link.attrs["href"] = f'{url}/{article.metadata["slug"]}/'
                soup.head.append(canonical_link)

                meta_url = soup.new_tag("meta")
                meta_url.attrs["name"] = "og:url"
                meta_url.attrs["property"] = "og:url"
                if article.metadata["slug"] == "index":
                    meta_url.attrs["content"] = f"{url}/"
                else:
                    meta_url.attrs["content"] = f'{url}/{article.metadata["slug"]}/'
                soup.head.append(meta_url)

                html = soup.prettify() if should_prettify else str(soup)
                cache.add(key, html, expire=markata.config.default_cache_expire)

            else:
                html = html_from_cache
            article.html = html
