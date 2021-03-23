"""manifest plugin"""
import json
from pathlib import Path

from bs4 import BeautifulSoup

from markata.hookspec import hook_impl

from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from markata import Markata


@hook_impl
def render(markata: "Markata") -> None:
    manifest = {
        "name": markata.config["site_name"],
        "short_name": markata.config["short_name"],
        "start_url": markata.config["start_url"],
        "display": markata.config["display"],
        "background_color": markata.config["background_color"],
        "theme_color": markata.config["theme_color"],
        "description": markata.config["description"],
        "icons": markata.icons,
    }
    filepath = Path(markata.output_dir) / "manifest.json"
    filepath.parent.mkdir(parents=True, exist_ok=True)
    filepath.touch(exist_ok=True)
    with open(filepath, "w+") as f:
        json.dump(manifest, f, ensure_ascii=True, indent=4)
    for article in markata.iter_articles("add manifest link"):
        key = markata.make_hash(
            "seo",
            "manifest",
            article["content_hash"],
        )
        html_from_cache = markata.cache.get(key)

        if html_from_cache is None:
            soup = BeautifulSoup(article.html, features="lxml")
            link = soup.new_tag("link")
            link.attrs["rel"] = "manifest"
            link.attrs["href"] = "/manifest.json"
            soup.head.append(link)

            html = soup.prettify()
            markata.cache.add(key, html, expire=15 * 24 * 60)
        else:
            html = html_from_cache
        article.html = html
