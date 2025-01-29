"""
The `markata.plugins.generator` plugin adds a meta generator tag to each generated HTML page,
indicating that the page was generated by Markata and including the version number.

# Installation

This plugin is built-in and enabled by default through the 'default' plugin. 
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.generator",
]
```

# Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.generator",
]
```

# Configuration

The plugin uses the global `prettify_html` configuration:

```toml
[markata]
prettify_html = true  # Set to false to disable HTML prettification
```

# Functionality

## Meta Tag Generation

The plugin adds a meta tag to the head of each HTML document:
```html
<meta name="generator" content="markata X.Y.Z">
```
where X.Y.Z is the current Markata version.

## Dependencies

This plugin depends on:
- The `render_markdown` plugin to provide HTML content
"""

from bs4 import BeautifulSoup

from markata import Markata, __version__
from markata.hookspec import hook_impl


@hook_impl(trylast=True)
def render(markata: Markata) -> None:
    should_prettify = markata.config.prettify_html
    for article in markata.iter_articles("add ssg tag"):
        soup = BeautifulSoup(article.html, features="lxml")
        tag = soup.new_tag("meta")
        tag.attrs["content"] = f"markata {__version__}"
        tag.attrs["name"] = "generator"
        soup.head.append(tag)

        if should_prettify:
            article.html = soup.prettify()
        else:
            article.html = str(soup)
