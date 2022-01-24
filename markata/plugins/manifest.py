"""manifest plugin"""
import json
from pathlib import Path
from typing import TYPE_CHECKING

from bs4 import BeautifulSoup

from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from markata.plugins.icon_resize import MarkataIcons


@hook_impl
def render(markata: "MarkataIcons") -> None:
    if "icons" in markata.__dict__.keys():
        icons = markata.icons
    else:
        icons = []
    manifest = {
        "name": markata.get_config("site_name") or "",
        "short_name": markata.get_config("short_name") or "",
        "start_url": markata.get_config("start_url") or "",
        "display": markata.get_config("display") or "",
        "background_color": markata.get_config("background_color") or "",
        "theme_color": markata.get_config("theme_color") or "",
        "description": markata.get_config("description") or "",
        "icons": icons,
    }
    filepath = Path(markata.config["output_dir"]) / "manifest.json"
    filepath.parent.mkdir(parents=True, exist_ok=True)
    filepath.touch(exist_ok=True)
    with open(filepath, "w+") as f:
        json.dump(manifest, f, ensure_ascii=True, indent=4)
    config = markata.get_plugin_config(__file__)
    with markata.cache as cache:
        for article in markata.iter_articles("add manifest link"):
            key = markata.make_hash(
                "seo",
                "manifest",
                article.content,
                article.html,
            )
            html_from_cache = cache.get(key)

            if html_from_cache is None:
                soup = BeautifulSoup(article.html, features="lxml")
                link = soup.new_tag("link")
                link.attrs["rel"] = "manifest"
                link.attrs["href"] = "/manifest.json"
                soup.head.append(link)

                html = soup.prettify()
                cache.add(key, html, expire=config["cache_expire"])
            else:
                html = html_from_cache
            article.html = html
