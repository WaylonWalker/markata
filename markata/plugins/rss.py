"""
[DEPRECATED] The `markata.plugins.rss` plugin is deprecated and will be removed in a
future version. Please use `markata.plugins.feeds` instead, which provides more
comprehensive feed generation capabilities.

# Installation

This plugin is deprecated. Use `markata.plugins.feeds` instead:

```toml
hooks = [
    "markata.plugins.feeds",  # Use this instead
    # "markata.plugins.rss",  # Deprecated
]
```

# Migration Guide

To migrate to the new feeds plugin:

1. Remove rss plugin from hooks:
```toml
# Remove or comment out
# "markata.plugins.rss"
```

2. Add feeds plugin:
```toml
hooks = [
    "markata.plugins.feeds"
]
```

3. Update configuration:
```toml
[markata.feeds]
# RSS feed configuration
rss = { output = "rss.xml" }

# Optional: Add other feed formats
atom = { output = "atom.xml" }
json = { output = "feed.json" }
```

See the feeds plugin documentation for more configuration options.

# Legacy Configuration

If you must continue using this plugin temporarily, configure in `markata.toml`:

```toml
[markata]
url = "https://example.com"
title = "Site Title"
author_name = "Author Name"
author_email = "author@example.com"
icon = "favicon.ico"
lang = "en"
```

# Dependencies

This plugin depends on:
- feedgen for RSS generation
- pytz for timezone handling

WARNING: This plugin is deprecated and will be removed in a future version.
Please migrate to `markata.plugins.feeds` as soon as possible.
"""

import datetime
import logging
from pathlib import Path
from typing import TYPE_CHECKING

import pytz
from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from markata import Markata
    from feedgen.feed import FeedGenerator

    class MarkataRss(Markata):
        fg: "FeedGenerator"
        rss: str


logger = logging.getLogger(__name__)


def _show_deprecation_warning():
    logger.warning(
        "DEPRECATION WARNING: The 'rss' plugin is deprecated and will be removed in a "
        "future version. Please migrate to the 'feeds' plugin which provides more "
        "comprehensive feed generation capabilities. See the documentation for migration "
        "instructions: https://markata.dev/markata/plugins/feeds/"
    )


@hook_impl(trylast=True)
def render(markata: "MarkataRss") -> None:
    _show_deprecation_warning()
    from feedgen.feed import FeedGenerator

    fg = FeedGenerator()
    url = markata.config.url or ""
    title = markata.config.title
    name = markata.config.author_name
    email = markata.config.author_email
    icon = str(markata.config.icon)
    lang = markata.config.lang
    rss_description = markata.config.rss_description or "rss feed"

    fg.id(str(url) + "/rss.xml")
    fg.title(title)
    fg.author(
        {
            "name": name,
            "email": email,
        },
    )
    fg.link(href=str(url), rel="alternate")
    fg.logo(icon)
    fg.subtitle(rss_description)
    fg.link(href=str(url) + "/rss.xml", rel="self")
    fg.language(lang)

    try:
        all_posts = sorted(markata.articles, key=lambda x: x["date"], reverse=True)
        posts = [post for post in all_posts if post["published"] == "True"]
    except BaseException:
        posts = markata.articles

    for article in posts:
        fe = fg.add_entry()
        fe.id(str(url + "/" + article.slug))
        fe.title(article.title)
        fe.published(
            datetime.datetime.combine(
                article.date or datetime.datetime.min.date(),
                datetime.datetime.min.time(),
                pytz.UTC,
            )
        )
        fe.description(article.description)
        fe.summary(article.long_description)
        fe.link(href=str(url) + "/" + article.slug)

        # Get the article HTML content
        article_html = article.article_html
        if isinstance(article_html, dict):
            # Try to get the 'index' HTML content first
            article_html = article_html.get('index')
            if article_html is None:
                # Fall back to the first available HTML content
                article_html = next(iter(article_html.values()))

        # Clean up the HTML content by removing control characters
        fe.content(article_html.translate(dict.fromkeys(range(32))))

    markata.fg = fg
    markata.rss = fg.rss_str(pretty=True)


@hook_impl
def save(markata: "MarkataRss") -> None:
    output_dir = Path(markata.config["output_dir"])
    markata.fg.rss_file(str(output_dir / "rss.xml"))
