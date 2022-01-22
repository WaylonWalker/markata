"""manifest plugin"""
from typing import TYPE_CHECKING, Dict, List

from bs4 import BeautifulSoup

from markata import Markata, __version__
from markata.hookspec import hook_impl

if TYPE_CHECKING:
    import frontmatter
    from bs4.element import Tag


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
                [p.text for p in soup.find(id="post-body").find_all("p")]
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


@hook_impl
def render(markata: Markata) -> None:

    url = markata.get_config("url") or ""
    images_url = markata.get_config("images_url") or url or ""
    site_name = markata.get_config("site_name") or ""
    author_name = markata.get_config("author_name") or ""
    author_email = markata.get_config("author_email") or ""
    twitter_creator = markata.get_config("twitter_creator") or ""
    twitter_card = markata.get_config("twitter_card") or "summary_large_image"
    config_seo = markata.get_config("seo", warn=False) or dict()

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

            html_from_cache = cache.get(key)

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

                html = soup.prettify()
                cache.add(key, html, expire=markata.config["default_cache_expire"])

            else:
                html = html_from_cache
            article.html = html
