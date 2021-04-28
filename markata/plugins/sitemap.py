from pathlib import Path

import anyconfig

from markata import Markata
from markata.hookspec import hook_impl


@hook_impl
def render(markata: Markata) -> None:
    sitemap = {
        "urlset": [
            {
                "url": {
                    "loc": markata.url + "/" + article["slug"] + "/",
                    "changefreq": "daily",
                    "priority": "0.7",
                }
            }
            for article in markata.articles
            if article["status"] == "published"
        ]
    }

    sitemap = (
        anyconfig.dumps(sitemap, "xml")
        .decode("utf-8")
        .replace(
            "<urlset>",
            '<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:news="http://www.google.com/schemas/sitemap-news/0.9" xmlns:xhtml="http://www.w3.org/1999/xhtml" xmlns:mobile="http://www.google.com/schemas/sitemap-mobile/1.0" xmlns:image="http://www.google.com/schemas/sitemap-image/1.1" xmlns:video="http://www.google.com/schemas/sitemap-video/1.1">',
        )
        .replace("</url>", "</url>\n")
    )
    setattr(markata, "sitemap", sitemap)


@hook_impl
def save(markata: Markata) -> None:
    with open(markata.output_dir / "sitemap.xml", "w") as f:
        f.write(markata.sitemap)  # type: ignore
