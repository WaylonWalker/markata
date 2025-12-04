"""
The `markata.plugins.render_markdown` plugin converts markdown content to HTML.
This plugin is essential for rendering markdown files loaded by the `load` plugin.

## Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.render_markdown",
]
```

## Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.render_markdown",
]
```

Note: Disabling this plugin will prevent markdown files from being rendered to HTML.

## Configuration

## Markdown Backend Selection

Choose from 3 supported markdown backends by setting `markdown_backend` in your `markata.toml`:

```toml
## choose your markdown backend
# markdown_backend='markdown'      # Python-Markdown
# markdown_backend='markdown2'     # markdown2
markdown_backend='markdown-it-py'  # markdown-it-py (default)
```

## Backend-Specific Configuration

### markdown-it-py

Configure markdown-it-py behavior in your `markata.toml`:

```toml
[markata.markdown_it_py]
# Set the flavor - options: 'zero', 'commonmark', 'gfm-like'
config = 'gfm-like'

# Enable specific plugins
enable = [
    'table',
    'strikethrough',
    'footnote',
]

# Disable specific plugins
disable = [
    'linkify',
]

# Configure plugins
[markata.markdown_it_py.plugins.footnote]
# Plugin-specific settings here
```

Read more about markdown-it-py settings in their [documentation](https://markdown-it-py.readthedocs.io/en/latest/).

## Cache Configuration

Control markdown rendering cache:

```toml
[markata.render_markdown]
cache_expire = 3600  # Cache expiration in seconds
```

## Functionality

## Registered Attributes

The plugin registers the following attributes on Post objects:
- `html`: The rendered HTML content from the markdown source

## Dependencies

This plugin depends on:
- One of: python-markdown, markdown2, or markdown-it-py (based on configuration)
- The `load` plugin to provide markdown content for rendering
"""

import concurrent.futures
import copy
import importlib
from enum import Enum
from functools import partial
from typing import TYPE_CHECKING
from typing import Dict
from typing import List
from typing import Optional

import pydantic

from markata import __version__
from markata.hookspec import hook_impl
from markata.hookspec import register_attr
from markata.plugins.md_it_highlight_code import highlight_code

if TYPE_CHECKING:
    from markata import Markata


class Backend(str, Enum):
    markdown = "markdown"
    markdown2 = "markdown2"
    markdown_it_py = "markdown-it-py"


class MdItExtension(pydantic.BaseModel):
    plugin: str
    config: Dict = None


class RenderMarkdownConfig(pydantic.BaseModel):
    backend: Backend = Backend("markdown-it-py")
    extensions: List[MdItExtension] = []
    cache_expire: int = 3600

    @pydantic.field_validator("extensions", mode="before")
    @classmethod
    def validate_extensions(cls, v) -> List[str]:
        if not isinstance(v, list):
            return [v]
        return v


class Config(pydantic.BaseModel):
    render_markdown: RenderMarkdownConfig = RenderMarkdownConfig()


class RenderMarkdownPost(pydantic.BaseModel):
    html: Optional[Dict[str, str] | str] = None
    model_config = pydantic.ConfigDict(
        validate_assignment=False,
        arbitrary_types_allowed=True,
        extra="allow",
        str_strip_whitespace=True,
        validate_default=True,
        coerce_numbers_to_str=True,
        populate_by_name=True,
    )


@hook_impl()
@register_attr("config_models")
def config_model(markata: "Markata") -> None:
    markata.config_models.append(Config)


@hook_impl()
@register_attr("post_models")
def post_model(markata: "Markata") -> None:
    markata.post_models.append(RenderMarkdownPost)


@hook_impl(tryfirst=True)
@register_attr("md", "markdown_extensions")
def configure(markata: "Markata") -> None:
    "Sets up a markdown instance as md"
    # if "markdown_extensions" not in markata.config:
    #     markdown_extensions = [""]
    # if isinstance(markata.config["markdown_extensions"], str):
    #     markdown_extensions = [markata.config["markdown_extensions"]]
    # if isinstance(markata.config["markdown_extensions"], list):
    #     markdown_extensions = markata.config["markdown_extensions"]
    # else:
    #     raise TypeError("markdown_extensions should be List[str]")

    # markata.markdown_extensions = [*DEFAULT_MD_EXTENSIONS, *markdown_extensions]

    if (
        markata.config.get("markdown_backend", "")
        .lower()
        .replace(" ", "-")
        .replace("_", "-")
        == "markdown-it-py"
    ):
        from markdown_it import MarkdownIt

        config_update = markata.config.get("markdown_it_py", {}).get(
            "options_update",
            {
                "linkify": True,
                "html": True,
                "typographer": True,
                "highlight": highlight_code,
            },
        )
        if isinstance(config_update.get("highlight"), str):
            module = config_update["highlight"].split(":")[0]
            func = config_update["highlight"].split(":")[1]
            config_update["highlight"] = getattr(
                importlib.import_module(module),
                func,
            )

        markata.md = MarkdownIt(
            markata.config.get("markdown_it_py", {}).get("config", "gfm-like"),
            config_update,
        )
        for plugin in markata.config.get("markdown_it_py", {}).get("enable", []):
            markata.md.enable(plugin)
        for plugin in markata.config.get("markdown_it_py", {}).get("disable", []):
            markata.md.disable(plugin)

        plugins = copy.deepcopy(
            markata.config.get("markdown_it_py", {}).get("plugins", []),
        )
        for plugin in plugins:
            if isinstance(plugin["plugin"], str):
                plugin["plugin_str"] = plugin["plugin"]
                plugin_module = plugin["plugin"].split(":")[0]
                plugin_func = plugin["plugin"].split(":")[1]
                plugin["plugin"] = getattr(
                    importlib.import_module(plugin_module),
                    plugin_func,
                )
            plugin["config"] = plugin.get("config", {})
            for k, _v in plugin["config"].items():
                if k == "markata":
                    plugin["config"][k] = markata

            markata.md = markata.md.use(plugin["plugin"], **plugin["config"])

        markata.md.convert = markata.md.render
        markata.md.toc = ""
    elif (
        markata.config.get("markdown_backend", "")
        .lower()
        .replace(" ", "-")
        .replace("_", "-")
        == "markdown2"
    ):
        import markdown2

        markata.md = markdown2.Markdown(
            extras=markata.config.render_markdown.extensions
        )
        markata.md.toc = ""
    else:
        import markdown

        markata.md = markdown.Markdown(
            extensions=markata.config.render_markdown.extensions
        )


@hook_impl(tryfirst=True)
@register_attr("rendered_posts")
def render(markata: "Markata") -> None:
    """Render markdown content in parallel."""
    config = markata.config.render_markdown
    articles = list(markata.filter("not skip"))

    with markata.cache as cache:
        with concurrent.futures.ThreadPoolExecutor() as executor:
            render_func = partial(render_article_parallel, markata, config, cache)
            args_list = [(article,) for article in articles]

            for article, html in executor.map(render_func, args_list):
                article.html = html
                article.article_html = copy.deepcopy(html)

    markata.rendered_posts = markata.posts


def render_article_parallel(markata, config, cache, article):
    # Handle article being passed as a tuple
    if isinstance(article, tuple):
        article = article[0]  # Extract the Post object from the tuple

    # Get content, defaulting to empty string
    content = getattr(article, "content", "")
    if not content:
        article.article_html = ""
        article.html = ""
        return article, ""

    key = markata.make_hash(
        "render_markdown",
        "render",
        content,
        __version__,
        markata.config.render_markdown.backend.value,
        str(markata.config.render_markdown.extensions),
    )
    html_from_cache = markata.precache.get(key)

    if html_from_cache is not None:
        article.article_html = html_from_cache
        article.html = html_from_cache
        return article, html_from_cache

    # Update markata instance with current article for plugin usage
    markata.md.options["article"] = article
    html = markata.md.convert(content)
    # Clear the article reference
    markata.md.options["article"] = None

    cache.set(key, html, expire=markata.config.markdown_cache_expire)
    article.article_html = html
    article.html = html
    return article, html
