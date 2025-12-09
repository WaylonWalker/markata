---
date: 2025-12-09
description: "The plugin converts markdown content to HTML. This plugin is essential
  for rendering markdown files loaded by the plugin. Installation This plugin is built-in\u2026"
published: false
slug: markata/plugins/render-markdown
title: render_markdown.py


---

---

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

---

!!! function
    <h2 id="configure" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">configure <em class="small">function</em></h2>

    Sets up a markdown instance as md

???+ source "configure <em class='small'>source</em>"
    ```python
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
    ```
!!! function
    <h2 id="render" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">render <em class="small">function</em></h2>

    Render markdown content in parallel.

???+ source "render <em class='small'>source</em>"
    ```python
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
    ```