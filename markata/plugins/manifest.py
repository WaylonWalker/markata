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
    icons = markata.icons if "icons" in markata.__dict__ else []
    manifest = {
        "name": markata.config.site_name,
        "short_name": markata.config.short_name,
        "start_url": markata.config.start_url,
        "display": markata.config.display,
        "background_color": str(markata.config.background_color),
        "theme_color": str(markata.config.theme_color),
        "description": markata.config.description,
        "icons": icons,
    }
    filepath = Path(markata.config["output_dir"]) / "manifest.json"
    filepath.touch(exist_ok=True)
    with open(filepath, "w+") as f:
        json.dump(manifest, f, ensure_ascii=True, indent=4)
    config = markata.get_plugin_config(__file__)
    should_prettify = markata.config.get("prettify_html", False)
    with markata.cache as cache:
        for article in markata.iter_articles("add manifest link"):
            key = markata.make_hash(
                "seo",
                "manifest",
                article.content,
                article.html,
            )
            html_from_cache = markata.precache.get(key)

            if html_from_cache is None:
                soup = BeautifulSoup(article.html, features="lxml")
                link = soup.new_tag("link")
                link.attrs["rel"] = "manifest"
                link.attrs["href"] = "/manifest.json"
                soup.head.append(link)

                html = soup.prettify() if should_prettify else str(soup)
                cache.add(key, html, expire=config["cache_expire"])
            else:
                html = html_from_cache
            article.html = html
