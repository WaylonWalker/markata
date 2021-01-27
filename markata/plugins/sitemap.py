from markata.hookspec import hook_impl
from markata import Markata
import anyconfig
from pathlib import Path


@hook_impl
def render(markata: Markata) -> None:
    sitemap = {
        "urlset": [
            {
                "url": {
                    "loc": markata.config["url"] + "/" + article["slug"],
                    "changefreq": "daily",
                    "priority": "0.7",
                }
            }
            for article in markata.articles
            if article["status"] == "published"
        ]
    }

    markata.sitemap = (
        anyconfig.dumps(sitemap, "xml")
        .decode("utf-8")
        .replace(
            "<urlset>",
            '<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:news="http://www.google.com/schemas/sitemap-news/0.9" xmlns:xhtml="http://www.w3.org/1999/xhtml" xmlns:mobile="http://www.google.com/schemas/sitemap-mobile/1.0" xmlns:image="http://www.google.com/schemas/sitemap-image/1.1" xmlns:video="http://www.google.com/schemas/sitemap-video/1.1">',
        )
        .replace("</url>", "</url>\n")
    )


@hook_impl
def save(markata: Markata) -> None:
    output_dir = Path(markata.config["output_dir"])
    with open(output_dir / "sitemap.xml", "w") as f:
        f.write(markata.sitemap)
