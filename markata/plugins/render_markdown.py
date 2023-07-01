"""
Renders markdown content as html.  This may be markdown files loaded in by way
of the [[load]] plugin.


## Markdown backend

There are 3 supported markdown backends that you can configure markata to use
by setting the `markdown_backend` in your `markata.toml`.

``` toml title=markata.toml
## choose your markdown backend
# markdown_backend='markdown'
# markdown_backend='markdown2'
markdown_backend='markdown-it-py'
```

## markdown-it-py configuration

`markdown-it-py` has quite a bit of configuration that you can do, you can read
more about the settings in their
[docs](https://markdown-it-py.readthedocs.io/en/latest/).

``` toml title=markata.toml
# markdown_it flavor
[markata.markdown_it_py]
config='gfm-like'
```

You can enable and disable built in plugins using the `enable` and `disable`
lists.

``` toml title=markata.toml
[markata.markdown_it_py]
# markdown_it built-in plugins
enable = [ "table" ]
disable = [ "image" ]
```

You can configure the `options_update` as follows, this is for the built-in
plugins and core configuration.

``` toml title=markata.toml
[markata.markdown_it_py.options_update]
linkify = true
html = true
typographer = true
highlight = 'markata.plugins.md_it_highlight_code:highlight_code'
```

Lastly you can add external plugins as follows.  Many of the great plugins that
come from [executable_books](https://github.com/executablebooks) actually comes
from a separate library
[mdit-py-plugins](https://mdit-py-plugins.readthedocs.io/), so they will be
configured here.

``` toml title=markata.toml
[[markata.markdown_it_py.plugins]]
plugin = "mdit_py_plugins.admon:admon_plugin"

[[markata.markdown_it_py.plugins]]
plugin = "mdit_py_plugins.admon:admon_plugin"

[[markata.markdown_it_py.plugins]]
plugin = "mdit_py_plugins.attrs:attrs_plugin"
config = {spans = true}

[[markata.markdown_it_py.plugins]]
plugin = "mdit_py_plugins.attrs:attrs_block_plugin"

[[markata.markdown_it_py.plugins]]
plugin = "markata.plugins.mdit_details:details_plugin"

[[markata.markdown_it_py.plugins]]
plugin = "mdit_py_plugins.anchors:anchors_plugin"

[markata.markdown_it_py.plugins.config]
permalink = true
permalinkSymbol = '<svg class="heading-permalink" aria-hidden="true" fill="currentColor" focusable="false" height="1em" viewBox="0 0 24 24" width="1em" xmlns="http://www.w3.org/2000/svg"><path d="M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992 5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242 6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985 3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005 3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201 13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995 2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0 0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836 19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0 0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632 1.563 1 1 0 0 0 1.602-1.198z"></path></svg>'

[[markata.markdown_it_py.plugins]]
plugin = "markata.plugins.md_it_wikilinks:wikilinks_plugin"
config = {markata = "markata"}
```

"""
import copy
from enum import Enum
import importlib
from typing import Dict, List, Optional, TYPE_CHECKING

import pydantic

from markata.hookspec import hook_impl, register_attr
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

    @pydantic.validator("extensions")
    def convert_to_list(cls, v):
        if not isinstance(v, list):
            return [v]
        return v


class Config(pydantic.BaseModel):
    render_markdown: RenderMarkdownConfig = RenderMarkdownConfig()


class RenderMarkdownPost(pydantic.BaseModel):
    article_html: Optional[str] = None
    html: Optional[str] = None


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
@register_attr("articles", "posts")
def render(markata: "Markata") -> None:
    config = markata.config.render_markdown
    with markata.cache as cache:
        for article in markata.iter_articles("rendering markdown"):
            key = markata.make_hash(
                "render_markdown",
                "render",
                article.content,
            )
            html_from_cache = markata.precache.get(key)
            if html_from_cache is None:
                html = markata.md.convert(article.content)
                cache.add(key, html, expire=config.cache_expire)
            else:
                html = html_from_cache
            article.html = html
            article.article_html = copy.deepcopy(html)

            article.html = html
            article.article_html = article.article_html
