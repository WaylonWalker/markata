from typing import TYPE_CHECKING

from bs4 import BeautifulSoup

from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from markata import Markata


@hook_impl
def render(markata: "Markata") -> None:
    config = markata.get_plugin_config(__file__)
    with markata.cache as cache:
        for article in markata.iter_articles("setting long description"):

            key = markata.make_hash(
                "long_description",
                "render",
                article["content_hash"],
            )

            description_from_cache = cache.get(key)

            if description_from_cache is None:

                if "long_description" in article.metadata:
                    description = article.metadata["long_description"]

                else:
                    soup = BeautifulSoup(article.html, features="lxml")
                    description = " ".join(
                        [p.text for p in soup.find(id="post-body").find_all("p")]
                    ).strip()[:250]

                markata.cache.add(key, description, expire=config["cache_expire"])

            else:
                description = description_from_cache

            article.metadata["long_description"] = description
