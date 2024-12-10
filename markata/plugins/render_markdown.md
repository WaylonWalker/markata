---
content: "Renders markdown content as html.  This may be markdown files loaded in
  by way\nof the [[load]] plugin.\n\n\n## Markdown backend\n\nThere are 3 supported
  markdown backends that you can configure markata to use\nby setting the `markdown_backend`
  in your `markata.toml`.\n\n``` toml title=markata.toml\n## choose your markdown
  backend\n# markdown_backend='markdown'\n# markdown_backend='markdown2'\nmarkdown_backend='markdown-it-py'\n```\n\n##
  markdown-it-py configuration\n\n`markdown-it-py` has quite a bit of configuration
  that you can do, you can read\nmore about the settings in their\n[docs](https://markdown-it-py.readthedocs.io/en/latest/).\n\n```
  toml title=markata.toml\n# markdown_it flavor\n[markata.markdown_it_py]\nconfig='gfm-like'\n```\n\nYou
  can enable and disable built in plugins using the `enable` and `disable`\nlists.\n\n```
  toml title=markata.toml\n[markata.markdown_it_py]\n# markdown_it built-in plugins\nenable
  = [ \"table\" ]\ndisable = [ \"image\" ]\n```\n\nYou can configure the `options_update`
  as follows, this is for the built-in\nplugins and core configuration.\n\n``` toml
  title=markata.toml\n[markata.markdown_it_py.options_update]\nlinkify = true\nhtml
  = true\ntypographer = true\nhighlight = 'markata.plugins.md_it_highlight_code:highlight_code'\n```\n\nLastly
  you can add external plugins as follows.  Many of the great plugins that\ncome from
  [executable_books](https://github.com/executablebooks) actually comes\nfrom a separate
  library\n[mdit-py-plugins](https://mdit-py-plugins.readthedocs.io/), so they will
  be\nconfigured here.\n\n``` toml title=markata.toml\n[[markata.markdown_it_py.plugins]]\nplugin
  = \"mdit_py_plugins.admon:admon_plugin\"\n\n[[markata.markdown_it_py.plugins]]\nplugin
  = \"mdit_py_plugins.admon:admon_plugin\"\n\n[[markata.markdown_it_py.plugins]]\nplugin
  = \"mdit_py_plugins.attrs:attrs_plugin\"\nconfig = {spans = true}\n\n[[markata.markdown_it_py.plugins]]\nplugin
  = \"mdit_py_plugins.attrs:attrs_block_plugin\"\n\n[[markata.markdown_it_py.plugins]]\nplugin
  = \"markata.plugins.mdit_details:details_plugin\"\n\n[[markata.markdown_it_py.plugins]]\nplugin
  = \"mdit_py_plugins.anchors:anchors_plugin\"\n\n[markata.markdown_it_py.plugins.config]\npermalink
  = true\npermalinkSymbol = '<svg class=\"heading-permalink\" aria-hidden=\"true\"
  fill=\"currentColor\" focusable=\"false\" height=\"1em\" viewBox=\"0 0 24 24\" width=\"1em\"
  xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949
  2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992 5.992 0 0 0 1.695-4.285 5.976
  5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242 6.003 6.003 0 0 0-1.905 1.24l-1.731
  1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985 3.985 0 0 1 2.761-1.123 3.975
  3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005 3.006a3.982 3.982 0
  0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201 13.6zm5.602-3.198a5.99
  5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995 2.994a5.992 5.992 0
  0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0 0 6.431 1.242 6.003
  6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836 19.81a3.985
  3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0 0 1-.111-5.644l3.005-3.006a3.982
  3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632 1.563 1 1 0 0 0 1.602-1.198z\"></path></svg>'\n\n[[markata.markdown_it_py.plugins]]\nplugin
  = \"markata.plugins.md_it_wikilinks:wikilinks_plugin\"\nconfig = {markata = \"markata\"}\n```\n\n\n!!
  class <h2 id='Backend' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Backend
  <em class='small'>class</em></h2>\n\n???+ source \"Backend <em class='small'>source</em>\"\n\n```python\n\n
  \       class Backend(str, Enum):\n            markdown = \"markdown\"\n            markdown2
  = \"markdown2\"\n            markdown_it_py = \"markdown-it-py\"\n```\n\n\n!! class
  <h2 id='MdItExtension' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>MdItExtension
  <em class='small'>class</em></h2>\n\n???+ source \"MdItExtension <em class='small'>source</em>\"\n\n```python\n\n
  \       class MdItExtension(pydantic.BaseModel):\n            plugin: str\n            config:
  Dict = None\n```\n\n\n!! class <h2 id='RenderMarkdownConfig' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>RenderMarkdownConfig <em class='small'>class</em></h2>\n\n???+
  source \"RenderMarkdownConfig <em class='small'>source</em>\"\n\n```python\n\n        class
  RenderMarkdownConfig(pydantic.BaseModel):\n            backend: Backend = Backend(\"markdown-it-py\")\n
  \           extensions: List[MdItExtension] = []\n            cache_expire: int
  = 3600\n\n            @pydantic.validator(\"extensions\")\n            def convert_to_list(cls,
  v):\n                if not isinstance(v, list):\n                    return [v]\n
  \               return v\n```\n\n\n!! class <h2 id='Config' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>Config <em class='small'>class</em></h2>\n\n???+
  source \"Config <em class='small'>source</em>\"\n\n```python\n\n        class Config(pydantic.BaseModel):\n
  \           render_markdown: RenderMarkdownConfig = RenderMarkdownConfig()\n```\n\n\n!!
  class <h2 id='RenderMarkdownPost' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>RenderMarkdownPost <em class='small'>class</em></h2>\n\n???+ source \"RenderMarkdownPost
  <em class='small'>source</em>\"\n\n```python\n\n        class RenderMarkdownPost(pydantic.BaseModel):\n
  \           article_html: Optional[str] = None\n            html: Optional[str |
  Dict[str, str]] = None\n```\n\n\n!! function <h2 id='config_model' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>config_model <em class='small'>function</em></h2>\n\n???+
  source \"config_model <em class='small'>source</em>\"\n\n```python\n\n        def
  config_model(markata: \"Markata\") -> None:\n            markata.config_models.append(Config)\n```\n\n\n!!
  function <h2 id='post_model' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>post_model <em class='small'>function</em></h2>\n\n???+ source \"post_model
  <em class='small'>source</em>\"\n\n```python\n\n        def post_model(markata:
  \"Markata\") -> None:\n            markata.post_models.append(RenderMarkdownPost)\n```\n\n\n!!
  function <h2 id='configure' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>configure <em class='small'>function</em></h2>\n    Sets up a markdown instance
  as md\n???+ source \"configure <em class='small'>source</em>\"\n\n```python\n\n
  \       def configure(markata: \"Markata\") -> None:\n            \"Sets up a markdown
  instance as md\"\n            # if \"markdown_extensions\" not in markata.config:\n
  \           #     markdown_extensions = [\"\"]\n            # if isinstance(markata.config[\"markdown_extensions\"],
  str):\n            #     markdown_extensions = [markata.config[\"markdown_extensions\"]]\n
  \           # if isinstance(markata.config[\"markdown_extensions\"], list):\n            #
  \    markdown_extensions = markata.config[\"markdown_extensions\"]\n            #
  else:\n            #     raise TypeError(\"markdown_extensions should be List[str]\")\n\n
  \           # markata.markdown_extensions = [*DEFAULT_MD_EXTENSIONS, *markdown_extensions]\n\n
  \           if (\n                markata.config.get(\"markdown_backend\", \"\")\n
  \               .lower()\n                .replace(\" \", \"-\")\n                .replace(\"_\",
  \"-\")\n                == \"markdown-it-py\"\n            ):\n                from
  markdown_it import MarkdownIt\n\n                config_update = markata.config.get(\"markdown_it_py\",
  {}).get(\n                    \"options_update\",\n                    {\n                        \"linkify\":
  True,\n                        \"html\": True,\n                        \"typographer\":
  True,\n                        \"highlight\": highlight_code,\n                    },\n
  \               )\n                if isinstance(config_update.get(\"highlight\"),
  str):\n                    module = config_update[\"highlight\"].split(\":\")[0]\n
  \                   func = config_update[\"highlight\"].split(\":\")[1]\n                    config_update[\"highlight\"]
  = getattr(\n                        importlib.import_module(module),\n                        func,\n
  \                   )\n\n                markata.md = MarkdownIt(\n                    markata.config.get(\"markdown_it_py\",
  {}).get(\"config\", \"gfm-like\"),\n                    config_update,\n                )\n
  \               for plugin in markata.config.get(\"markdown_it_py\", {}).get(\"enable\",
  []):\n                    markata.md.enable(plugin)\n                for plugin
  in markata.config.get(\"markdown_it_py\", {}).get(\"disable\", []):\n                    markata.md.disable(plugin)\n\n
  \               plugins = copy.deepcopy(\n                    markata.config.get(\"markdown_it_py\",
  {}).get(\"plugins\", []),\n                )\n                for plugin in plugins:\n
  \                   if isinstance(plugin[\"plugin\"], str):\n                        plugin[\"plugin_str\"]
  = plugin[\"plugin\"]\n                        plugin_module = plugin[\"plugin\"].split(\":\")[0]\n
  \                       plugin_func = plugin[\"plugin\"].split(\":\")[1]\n                        plugin[\"plugin\"]
  = getattr(\n                            importlib.import_module(plugin_module),\n
  \                           plugin_func,\n                        )\n                    plugin[\"config\"]
  = plugin.get(\"config\", {})\n                    for k, _v in plugin[\"config\"].items():\n
  \                       if k == \"markata\":\n                            plugin[\"config\"][k]
  = markata\n\n                    markata.md = markata.md.use(plugin[\"plugin\"],
  **plugin[\"config\"])\n\n                markata.md.convert = markata.md.render\n
  \               markata.md.toc = \"\"\n            elif (\n                markata.config.get(\"markdown_backend\",
  \"\")\n                .lower()\n                .replace(\" \", \"-\")\n                .replace(\"_\",
  \"-\")\n                == \"markdown2\"\n            ):\n                import
  markdown2\n\n                markata.md = markdown2.Markdown(\n                    extras=markata.config.render_markdown.extensions\n
  \               )\n                markata.md.toc = \"\"\n            else:\n                import
  markdown\n\n                markata.md = markdown.Markdown(\n                    extensions=markata.config.render_markdown.extensions\n
  \               )\n```\n\n\n!! function <h2 id='render' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>render <em class='small'>function</em></h2>\n\n???+
  source \"render <em class='small'>source</em>\"\n\n```python\n\n        def render(markata:
  \"Markata\") -> None:\n            config = markata.config.render_markdown\n            with
  markata.cache as cache:\n                for article in markata.articles:\n                    article.html
  = render_article(markata, config, cache, article)\n                    article.article_html
  = copy.deepcopy(article.html)\n```\n\n\n!! function <h2 id='render_article' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>render_article <em class='small'>function</em></h2>\n\n???+
  source \"render_article <em class='small'>source</em>\"\n\n```python\n\n        def
  render_article(markata: \"Markata\", config, cache, article):\n            key =
  markata.make_hash(\n                \"render_markdown\",\n                \"render\",\n
  \               article.content,\n            )\n            html_from_cache = markata.precache.get(key)\n
  \           if html_from_cache is None:\n                html = markata.md.convert(article.content)\n
  \               cache.add(key, html, expire=config.cache_expire)\n            else:\n
  \               html = html_from_cache\n            return html\n\n            article.html
  = html\n            article.article_html = copy.deepcopy(html)\n\n            article.html
  = html\n            article.article_html = article.article_html\n```\n\n\n!! method
  <h2 id='convert_to_list' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>convert_to_list <em class='small'>method</em></h2>\n\n???+ source \"convert_to_list
  <em class='small'>source</em>\"\n\n```python\n\n        def convert_to_list(cls,
  v):\n                if not isinstance(v, list):\n                    return [v]\n
  \               return v\n```\n\n"
date: 0001-01-01
description: Renders markdown content as html.  This may be markdown files loaded
  in by way There are 3 supported markdown backends that you can configure markata
  to use mar
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Render_Markdown.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Renders markdown content as html.  This
    may be markdown files loaded in by way There are 3 supported markdown backends
    that you can configure markata to use mar\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Render_Markdown.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Renders markdown content as html.  This
    may be markdown files loaded in by way There are 3 supported markdown backends
    that you can configure markata to use mar\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<div class='container flex flex-row min-h-screen'>\n    <div>\n    </div>\n
    \   <div class='flex-grow px-8 mx-auto min-h-screen'>\n<header class='flex justify-center
    items-center p-8'>\n\n    <nav class='flex justify-center items-center my-8'>\n
    \       <a\n            href='/'>markata</a>\n        <a\n            href='https://github.com/WaylonWalker/markata'>GitHub</a>\n
    \       <a\n            href='https://markata.dev/docs/'>docs</a>\n        <a\n
    \           href='https://markata.dev/plugins/'>plugins</a>\n    </nav>\n\n    <div>\n
    \       <label id=\"theme-switch\" class=\"theme-switch\" for=\"checkbox-theme\"
    title=\"light/dark mode toggle\">\n            <input type=\"checkbox\" id=\"checkbox-theme\"
    />\n            <div class=\"slider round\"></div>\n        </label>\n    </div>\n</header><article
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Render_Markdown.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <p>Renders markdown
    content as html.  This may be markdown files loaded in by way\nof the <a class=\"wikilink\"
    href=\"/markata/plugins/load\">load</a> plugin.</p>\n<h2 id=\"markdown-backend\">Markdown
    backend <a class=\"header-anchor\" href=\"#markdown-backend\"><svg class=\"heading-permalink\"
    aria-hidden=\"true\" fill=\"currentColor\" focusable=\"false\" height=\"1em\"
    viewBox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
    d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
    5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
    6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
    3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
    3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
    13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
    2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0
    0 0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
    19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997
    0 0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>There are 3 supported
    markdown backends that you can configure markata to use\nby setting the <code>markdown_backend</code>
    in your <code>markata.toml</code>.</p>\n<pre class='wrapper'>\n\n<div class='filepath'>\n<p>markata.toml</p>\n\n<div
    class='right'>\n\n<button class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n</div>\n\n<div
    class=\"highlight\"><pre><span></span><span class=\"c1\">## choose your markdown
    backend</span>\n<span class=\"c1\"># markdown_backend=&#39;markdown&#39;</span>\n<span
    class=\"c1\"># markdown_backend=&#39;markdown2&#39;</span>\n<span class=\"n\">markdown_backend</span><span
    class=\"o\">=</span><span class=\"s1\">&#39;markdown-it-py&#39;</span>\n</pre></div>\n\n</pre>\n\n<h2
    id=\"markdown-it-py-configuration\">markdown-it-py configuration <a class=\"header-anchor\"
    href=\"#markdown-it-py-configuration\"><svg class=\"heading-permalink\" aria-hidden=\"true\"
    fill=\"currentColor\" focusable=\"false\" height=\"1em\" viewBox=\"0 0 24 24\"
    width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M9.199 13.599a5.99
    5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992 5.992 0
    0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242 6.003
    6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
    3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
    3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
    13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
    2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0
    0 0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
    19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997
    0 0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p><code>markdown-it-py</code>
    has quite a bit of configuration that you can do, you can read\nmore about the
    settings in their\n<a href=\"https://markdown-it-py.readthedocs.io/en/latest/\">docs</a>.</p>\n<pre
    class='wrapper'>\n\n<div class='filepath'>\n<p>markata.toml</p>\n\n<div class='right'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n</div>\n\n<div
    class=\"highlight\"><pre><span></span><span class=\"c1\"># markdown_it flavor</span>\n<span
    class=\"k\">[markata.markdown_it_py]</span>\n<span class=\"n\">config</span><span
    class=\"o\">=</span><span class=\"s1\">&#39;gfm-like&#39;</span>\n</pre></div>\n\n</pre>\n\n<p>You
    can enable and disable built in plugins using the <code>enable</code> and <code>disable</code>\nlists.</p>\n<pre
    class='wrapper'>\n\n<div class='filepath'>\n<p>markata.toml</p>\n\n<div class='right'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n</div>\n\n<div
    class=\"highlight\"><pre><span></span><span class=\"k\">[markata.markdown_it_py]</span>\n<span
    class=\"c1\"># markdown_it built-in plugins</span>\n<span class=\"n\">enable</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"p\">[</span><span class=\"w\"> </span><span class=\"s2\">&quot;table&quot;</span><span
    class=\"w\"> </span><span class=\"p\">]</span>\n<span class=\"n\">disable</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"p\">[</span><span class=\"w\"> </span><span class=\"s2\">&quot;image&quot;</span><span
    class=\"w\"> </span><span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>You
    can configure the <code>options_update</code> as follows, this is for the built-in\nplugins
    and core configuration.</p>\n<pre class='wrapper'>\n\n<div class='filepath'>\n<p>markata.toml</p>\n\n<div
    class='right'>\n\n<button class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n</div>\n\n<div
    class=\"highlight\"><pre><span></span><span class=\"k\">[markata.markdown_it_py.options_update]</span>\n<span
    class=\"n\">linkify</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"kc\">true</span>\n<span class=\"n\">html</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"kc\">true</span>\n<span class=\"n\">typographer</span><span class=\"w\">
    </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"kc\">true</span>\n<span
    class=\"n\">highlight</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s1\">&#39;markata.plugins.md_it_highlight_code:highlight_code&#39;</span>\n</pre></div>\n\n</pre>\n\n<p>Lastly
    you can add external plugins as follows.  Many of the great plugins that\ncome
    from <a href=\"https://github.com/executablebooks\">executable_books</a> actually
    comes\nfrom a separate library\n<a href=\"https://mdit-py-plugins.readthedocs.io/\">mdit-py-plugins</a>,
    so they will be\nconfigured here.</p>\n<pre class='wrapper'>\n\n<div class='filepath'>\n<p>markata.toml</p>\n\n<div
    class='right'>\n\n<button class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n</div>\n\n<div
    class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.markdown_it_py.plugins]]</span>\n<span
    class=\"n\">plugin</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;mdit_py_plugins.admon:admon_plugin&quot;</span>\n\n<span
    class=\"k\">[[markata.markdown_it_py.plugins]]</span>\n<span class=\"n\">plugin</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;mdit_py_plugins.admon:admon_plugin&quot;</span>\n\n<span class=\"k\">[[markata.markdown_it_py.plugins]]</span>\n<span
    class=\"n\">plugin</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;mdit_py_plugins.attrs:attrs_plugin&quot;</span>\n<span
    class=\"n\">config</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"p\">{</span><span class=\"n\">spans</span><span
    class=\"w\"> </span><span class=\"p\">=</span><span class=\"w\"> </span><span
    class=\"kc\">true</span><span class=\"p\">}</span>\n\n<span class=\"k\">[[markata.markdown_it_py.plugins]]</span>\n<span
    class=\"n\">plugin</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;mdit_py_plugins.attrs:attrs_block_plugin&quot;</span>\n\n<span
    class=\"k\">[[markata.markdown_it_py.plugins]]</span>\n<span class=\"n\">plugin</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;markata.plugins.mdit_details:details_plugin&quot;</span>\n\n<span
    class=\"k\">[[markata.markdown_it_py.plugins]]</span>\n<span class=\"n\">plugin</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;mdit_py_plugins.anchors:anchors_plugin&quot;</span>\n\n<span
    class=\"k\">[markata.markdown_it_py.plugins.config]</span>\n<span class=\"n\">permalink</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"kc\">true</span>\n<span class=\"n\">permalinkSymbol</span><span class=\"w\">
    </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s1\">&#39;&lt;svg
    class=&quot;heading-permalink&quot; aria-hidden=&quot;true&quot; fill=&quot;currentColor&quot;
    focusable=&quot;false&quot; height=&quot;1em&quot; viewBox=&quot;0 0 24 24&quot;
    width=&quot;1em&quot; xmlns=&quot;http://www.w3.org/2000/svg&quot;&gt;&lt;path
    d=&quot;M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
    5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
    6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
    3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
    3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
    13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
    2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0
    0 0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
    19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997
    0 0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
    1.563 1 1 0 0 0 1.602-1.198z&quot;&gt;&lt;/path&gt;&lt;/svg&gt;&#39;</span>\n\n<span
    class=\"k\">[[markata.markdown_it_py.plugins]]</span>\n<span class=\"n\">plugin</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;markata.plugins.md_it_wikilinks:wikilinks_plugin&quot;</span>\n<span
    class=\"n\">config</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"p\">{</span><span class=\"n\">markata</span><span
    class=\"w\"> </span><span class=\"p\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;markata&quot;</span><span class=\"p\">}</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='Backend' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>Backend <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">Backend
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
    <span class=\"nc\">Backend</span><span class=\"p\">(</span><span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"n\">Enum</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">markdown</span> <span class=\"o\">=</span> <span
    class=\"s2\">&quot;markdown&quot;</span>\n            <span class=\"n\">markdown2</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;markdown2&quot;</span>\n            <span
    class=\"n\">markdown_it_py</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;markdown-it-py&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='MdItExtension' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>MdItExtension <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">MdItExtension
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
    <span class=\"nc\">MdItExtension</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">plugin</span><span class=\"p\">:</span> <span class=\"nb\">str</span>\n
    \           <span class=\"n\">config</span><span class=\"p\">:</span> <span class=\"n\">Dict</span>
    <span class=\"o\">=</span> <span class=\"kc\">None</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='RenderMarkdownConfig' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>RenderMarkdownConfig <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">RenderMarkdownConfig
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
    <span class=\"nc\">RenderMarkdownConfig</span><span class=\"p\">(</span><span
    class=\"n\">pydantic</span><span class=\"o\">.</span><span class=\"n\">BaseModel</span><span
    class=\"p\">):</span>\n            <span class=\"n\">backend</span><span class=\"p\">:</span>
    <span class=\"n\">Backend</span> <span class=\"o\">=</span> <span class=\"n\">Backend</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;markdown-it-py&quot;</span><span
    class=\"p\">)</span>\n            <span class=\"n\">extensions</span><span class=\"p\">:</span>
    <span class=\"n\">List</span><span class=\"p\">[</span><span class=\"n\">MdItExtension</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"p\">[]</span>\n
    \           <span class=\"n\">cache_expire</span><span class=\"p\">:</span> <span
    class=\"nb\">int</span> <span class=\"o\">=</span> <span class=\"mi\">3600</span>\n\n
    \           <span class=\"nd\">@pydantic</span><span class=\"o\">.</span><span
    class=\"n\">validator</span><span class=\"p\">(</span><span class=\"s2\">&quot;extensions&quot;</span><span
    class=\"p\">)</span>\n            <span class=\"k\">def</span> <span class=\"nf\">convert_to_list</span><span
    class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
    class=\"n\">v</span><span class=\"p\">):</span>\n                <span class=\"k\">if</span>
    <span class=\"ow\">not</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"nb\">list</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"p\">[</span><span class=\"n\">v</span><span class=\"p\">]</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Config
    <em class='small'>class</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">Config <em class='small'>source</em></p>\n</div>\n<pre
    class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button class='copy' title='copy
    code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
    <span class=\"nc\">Config</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">render_markdown</span><span class=\"p\">:</span>
    <span class=\"n\">RenderMarkdownConfig</span> <span class=\"o\">=</span> <span
    class=\"n\">RenderMarkdownConfig</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='RenderMarkdownPost' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>RenderMarkdownPost <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">RenderMarkdownPost
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
    <span class=\"nc\">RenderMarkdownPost</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">article_html</span><span class=\"p\">:</span> <span
    class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n
    \           <span class=\"n\">html</span><span class=\"p\">:</span> <span class=\"n\">Optional</span><span
    class=\"p\">[</span><span class=\"nb\">str</span> <span class=\"o\">|</span> <span
    class=\"n\">Dict</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"nb\">str</span><span class=\"p\">]]</span>
    <span class=\"o\">=</span> <span class=\"kc\">None</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='config_model' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>config_model <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">config_model
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">config_model</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config_models</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">Config</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='post_model' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>post_model <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">post_model
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">post_model</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">post_models</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">RenderMarkdownPost</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='configure' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>configure <em class='small'>function</em></h2>\nSets up a markdown instance
    as md</p>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">configure <em class='small'>source</em></p>\n</div>\n<pre
    class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button class='copy' title='copy
    code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">configure</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"s2\">&quot;Sets up a markdown instance as md&quot;</span>\n
    \           <span class=\"c1\"># if &quot;markdown_extensions&quot; not in markata.config:</span>\n
    \           <span class=\"c1\">#     markdown_extensions = [&quot;&quot;]</span>\n
    \           <span class=\"c1\"># if isinstance(markata.config[&quot;markdown_extensions&quot;],
    str):</span>\n            <span class=\"c1\">#     markdown_extensions = [markata.config[&quot;markdown_extensions&quot;]]</span>\n
    \           <span class=\"c1\"># if isinstance(markata.config[&quot;markdown_extensions&quot;],
    list):</span>\n            <span class=\"c1\">#     markdown_extensions = markata.config[&quot;markdown_extensions&quot;]</span>\n
    \           <span class=\"c1\"># else:</span>\n            <span class=\"c1\">#
    \    raise TypeError(&quot;markdown_extensions should be List[str]&quot;)</span>\n\n
    \           <span class=\"c1\"># markata.markdown_extensions = [*DEFAULT_MD_EXTENSIONS,
    *markdown_extensions]</span>\n\n            <span class=\"k\">if</span> <span
    class=\"p\">(</span>\n                <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;markdown_backend&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"o\">.</span><span class=\"n\">lower</span><span
    class=\"p\">()</span>\n                <span class=\"o\">.</span><span class=\"n\">replace</span><span
    class=\"p\">(</span><span class=\"s2\">&quot; &quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;-&quot;</span><span class=\"p\">)</span>\n                <span
    class=\"o\">.</span><span class=\"n\">replace</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;_&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;-&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"o\">==</span> <span class=\"s2\">&quot;markdown-it-py&quot;</span>\n
    \           <span class=\"p\">):</span>\n                <span class=\"kn\">from</span>
    <span class=\"nn\">markdown_it</span> <span class=\"kn\">import</span> <span class=\"n\">MarkdownIt</span>\n\n
    \               <span class=\"n\">config_update</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;markdown_it_py&quot;</span><span class=\"p\">,</span> <span
    class=\"p\">{})</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span>\n                    <span class=\"s2\">&quot;options_update&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"p\">{</span>\n                        <span
    class=\"s2\">&quot;linkify&quot;</span><span class=\"p\">:</span> <span class=\"kc\">True</span><span
    class=\"p\">,</span>\n                        <span class=\"s2\">&quot;html&quot;</span><span
    class=\"p\">:</span> <span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \                       <span class=\"s2\">&quot;typographer&quot;</span><span
    class=\"p\">:</span> <span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \                       <span class=\"s2\">&quot;highlight&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">highlight_code</span><span class=\"p\">,</span>\n                    <span
    class=\"p\">},</span>\n                <span class=\"p\">)</span>\n                <span
    class=\"k\">if</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">config_update</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;highlight&quot;</span><span class=\"p\">),</span>
    <span class=\"nb\">str</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">module</span> <span class=\"o\">=</span> <span class=\"n\">config_update</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;highlight&quot;</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">split</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;:&quot;</span><span class=\"p\">)[</span><span class=\"mi\">0</span><span
    class=\"p\">]</span>\n                    <span class=\"n\">func</span> <span
    class=\"o\">=</span> <span class=\"n\">config_update</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;highlight&quot;</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">split</span><span class=\"p\">(</span><span class=\"s2\">&quot;:&quot;</span><span
    class=\"p\">)[</span><span class=\"mi\">1</span><span class=\"p\">]</span>\n                    <span
    class=\"n\">config_update</span><span class=\"p\">[</span><span class=\"s2\">&quot;highlight&quot;</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"nb\">getattr</span><span
    class=\"p\">(</span>\n                        <span class=\"n\">importlib</span><span
    class=\"o\">.</span><span class=\"n\">import_module</span><span class=\"p\">(</span><span
    class=\"n\">module</span><span class=\"p\">),</span>\n                        <span
    class=\"n\">func</span><span class=\"p\">,</span>\n                    <span class=\"p\">)</span>\n\n
    \               <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">md</span> <span class=\"o\">=</span> <span class=\"n\">MarkdownIt</span><span
    class=\"p\">(</span>\n                    <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;markdown_it_py&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">{})</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;config&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;gfm-like&quot;</span><span class=\"p\">),</span>\n
    \                   <span class=\"n\">config_update</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">)</span>\n                <span class=\"k\">for</span>
    <span class=\"n\">plugin</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;markdown_it_py&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">{})</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;enable&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">[]):</span>\n                    <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">md</span><span
    class=\"o\">.</span><span class=\"n\">enable</span><span class=\"p\">(</span><span
    class=\"n\">plugin</span><span class=\"p\">)</span>\n                <span class=\"k\">for</span>
    <span class=\"n\">plugin</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;markdown_it_py&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">{})</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;disable&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">[]):</span>\n                    <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">md</span><span
    class=\"o\">.</span><span class=\"n\">disable</span><span class=\"p\">(</span><span
    class=\"n\">plugin</span><span class=\"p\">)</span>\n\n                <span class=\"n\">plugins</span>
    <span class=\"o\">=</span> <span class=\"n\">copy</span><span class=\"o\">.</span><span
    class=\"n\">deepcopy</span><span class=\"p\">(</span>\n                    <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;markdown_it_py&quot;</span><span class=\"p\">,</span> <span
    class=\"p\">{})</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;plugins&quot;</span><span class=\"p\">,</span>
    <span class=\"p\">[]),</span>\n                <span class=\"p\">)</span>\n                <span
    class=\"k\">for</span> <span class=\"n\">plugin</span> <span class=\"ow\">in</span>
    <span class=\"n\">plugins</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">if</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">plugin</span><span class=\"p\">[</span><span class=\"s2\">&quot;plugin&quot;</span><span
    class=\"p\">],</span> <span class=\"nb\">str</span><span class=\"p\">):</span>\n
    \                       <span class=\"n\">plugin</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;plugin_str&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"n\">plugin</span><span class=\"p\">[</span><span class=\"s2\">&quot;plugin&quot;</span><span
    class=\"p\">]</span>\n                        <span class=\"n\">plugin_module</span>
    <span class=\"o\">=</span> <span class=\"n\">plugin</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;plugin&quot;</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">split</span><span class=\"p\">(</span><span class=\"s2\">&quot;:&quot;</span><span
    class=\"p\">)[</span><span class=\"mi\">0</span><span class=\"p\">]</span>\n                        <span
    class=\"n\">plugin_func</span> <span class=\"o\">=</span> <span class=\"n\">plugin</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;plugin&quot;</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">split</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;:&quot;</span><span class=\"p\">)[</span><span class=\"mi\">1</span><span
    class=\"p\">]</span>\n                        <span class=\"n\">plugin</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;plugin&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"nb\">getattr</span><span class=\"p\">(</span>\n
    \                           <span class=\"n\">importlib</span><span class=\"o\">.</span><span
    class=\"n\">import_module</span><span class=\"p\">(</span><span class=\"n\">plugin_module</span><span
    class=\"p\">),</span>\n                            <span class=\"n\">plugin_func</span><span
    class=\"p\">,</span>\n                        <span class=\"p\">)</span>\n                    <span
    class=\"n\">plugin</span><span class=\"p\">[</span><span class=\"s2\">&quot;config&quot;</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"n\">plugin</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;config&quot;</span><span class=\"p\">,</span> <span class=\"p\">{})</span>\n
    \                   <span class=\"k\">for</span> <span class=\"n\">k</span><span
    class=\"p\">,</span> <span class=\"n\">_v</span> <span class=\"ow\">in</span>
    <span class=\"n\">plugin</span><span class=\"p\">[</span><span class=\"s2\">&quot;config&quot;</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">items</span><span
    class=\"p\">():</span>\n                        <span class=\"k\">if</span> <span
    class=\"n\">k</span> <span class=\"o\">==</span> <span class=\"s2\">&quot;markata&quot;</span><span
    class=\"p\">:</span>\n                            <span class=\"n\">plugin</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;config&quot;</span><span class=\"p\">][</span><span
    class=\"n\">k</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"n\">markata</span>\n\n                    <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">md</span> <span class=\"o\">=</span> <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">md</span><span
    class=\"o\">.</span><span class=\"n\">use</span><span class=\"p\">(</span><span
    class=\"n\">plugin</span><span class=\"p\">[</span><span class=\"s2\">&quot;plugin&quot;</span><span
    class=\"p\">],</span> <span class=\"o\">**</span><span class=\"n\">plugin</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;config&quot;</span><span class=\"p\">])</span>\n\n
    \               <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">md</span><span class=\"o\">.</span><span class=\"n\">convert</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">md</span><span class=\"o\">.</span><span class=\"n\">render</span>\n
    \               <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">md</span><span class=\"o\">.</span><span class=\"n\">toc</span> <span
    class=\"o\">=</span> <span class=\"s2\">&quot;&quot;</span>\n            <span
    class=\"k\">elif</span> <span class=\"p\">(</span>\n                <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;markdown_backend&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"o\">.</span><span class=\"n\">lower</span><span
    class=\"p\">()</span>\n                <span class=\"o\">.</span><span class=\"n\">replace</span><span
    class=\"p\">(</span><span class=\"s2\">&quot; &quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;-&quot;</span><span class=\"p\">)</span>\n                <span
    class=\"o\">.</span><span class=\"n\">replace</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;_&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;-&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"o\">==</span> <span class=\"s2\">&quot;markdown2&quot;</span>\n
    \           <span class=\"p\">):</span>\n                <span class=\"kn\">import</span>
    <span class=\"nn\">markdown2</span>\n\n                <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">md</span> <span class=\"o\">=</span> <span
    class=\"n\">markdown2</span><span class=\"o\">.</span><span class=\"n\">Markdown</span><span
    class=\"p\">(</span>\n                    <span class=\"n\">extras</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">render_markdown</span><span
    class=\"o\">.</span><span class=\"n\">extensions</span>\n                <span
    class=\"p\">)</span>\n                <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">md</span><span class=\"o\">.</span><span class=\"n\">toc</span> <span
    class=\"o\">=</span> <span class=\"s2\">&quot;&quot;</span>\n            <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                <span class=\"kn\">import</span>
    <span class=\"nn\">markdown</span>\n\n                <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">md</span> <span class=\"o\">=</span> <span
    class=\"n\">markdown</span><span class=\"o\">.</span><span class=\"n\">Markdown</span><span
    class=\"p\">(</span>\n                    <span class=\"n\">extensions</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">render_markdown</span><span
    class=\"o\">.</span><span class=\"n\">extensions</span>\n                <span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='render'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>render <em class='small'>function</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">render
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">render</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">config</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">render_markdown</span>\n            <span class=\"k\">with</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">cache</span>
    <span class=\"k\">as</span> <span class=\"n\">cache</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">for</span> <span class=\"n\">article</span>
    <span class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">articles</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">html</span>
    <span class=\"o\">=</span> <span class=\"n\">render_article</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">,</span> <span class=\"n\">config</span><span
    class=\"p\">,</span> <span class=\"n\">cache</span><span class=\"p\">,</span>
    <span class=\"n\">article</span><span class=\"p\">)</span>\n                    <span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">article_html</span>
    <span class=\"o\">=</span> <span class=\"n\">copy</span><span class=\"o\">.</span><span
    class=\"n\">deepcopy</span><span class=\"p\">(</span><span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">html</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='render_article' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>render_article <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">render_article
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">render_article</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">config</span><span class=\"p\">,</span> <span class=\"n\">cache</span><span
    class=\"p\">,</span> <span class=\"n\">article</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">key</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">make_hash</span><span class=\"p\">(</span>\n
    \               <span class=\"s2\">&quot;render_markdown&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"s2\">&quot;render&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">content</span><span class=\"p\">,</span>\n            <span class=\"p\">)</span>\n
    \           <span class=\"n\">html_from_cache</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">precache</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">)</span>\n            <span class=\"k\">if</span>
    <span class=\"n\">html_from_cache</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"n\">html</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">md</span><span
    class=\"o\">.</span><span class=\"n\">convert</span><span class=\"p\">(</span><span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">content</span><span
    class=\"p\">)</span>\n                <span class=\"n\">cache</span><span class=\"o\">.</span><span
    class=\"n\">add</span><span class=\"p\">(</span><span class=\"n\">key</span><span
    class=\"p\">,</span> <span class=\"n\">html</span><span class=\"p\">,</span> <span
    class=\"n\">expire</span><span class=\"o\">=</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">cache_expire</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">else</span><span class=\"p\">:</span>\n                <span
    class=\"n\">html</span> <span class=\"o\">=</span> <span class=\"n\">html_from_cache</span>\n
    \           <span class=\"k\">return</span> <span class=\"n\">html</span>\n\n
    \           <span class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">html</span>
    <span class=\"o\">=</span> <span class=\"n\">html</span>\n            <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">article_html</span> <span class=\"o\">=</span>
    <span class=\"n\">copy</span><span class=\"o\">.</span><span class=\"n\">deepcopy</span><span
    class=\"p\">(</span><span class=\"n\">html</span><span class=\"p\">)</span>\n\n
    \           <span class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">html</span>
    <span class=\"o\">=</span> <span class=\"n\">html</span>\n            <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">article_html</span> <span class=\"o\">=</span>
    <span class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">article_html</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='convert_to_list' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>convert_to_list <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">convert_to_list
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">convert_to_list</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">):</span>\n                <span
    class=\"k\">if</span> <span class=\"ow\">not</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"nb\">list</span><span class=\"p\">):</span>\n                    <span
    class=\"k\">return</span> <span class=\"p\">[</span><span class=\"n\">v</span><span
    class=\"p\">]</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>    </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Render_Markdown.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Renders markdown content as html.  This
    may be markdown files loaded in by way There are 3 supported markdown backends
    that you can configure markata to use mar\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Render_Markdown.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Renders markdown content as html.  This
    may be markdown files loaded in by way There are 3 supported markdown backends
    that you can configure markata to use mar\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<article style=\"text-align: center;\">\n    <style>\n        section
    {\n            font-size: 200%;\n        }\n\n\n        .edit {\n            display:
    none;\n        }\n    </style>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Render_Markdown.Py\n    </h1>\n</section></article>\n     </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Render_Markdown.Py\n    </h1>\n</section>    <section class=\"body\">\n
    \       <p>Renders markdown content as html.  This may be markdown files loaded
    in by way\nof the <a class=\"wikilink\" href=\"/markata/plugins/load\">load</a>
    plugin.</p>\n<h2 id=\"markdown-backend\">Markdown backend <a class=\"header-anchor\"
    href=\"#markdown-backend\"><svg class=\"heading-permalink\" aria-hidden=\"true\"
    fill=\"currentColor\" focusable=\"false\" height=\"1em\" viewBox=\"0 0 24 24\"
    width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M9.199 13.599a5.99
    5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992 5.992 0
    0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242 6.003
    6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
    3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
    3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
    13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
    2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0
    0 0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
    19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997
    0 0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>There are 3 supported
    markdown backends that you can configure markata to use\nby setting the <code>markdown_backend</code>
    in your <code>markata.toml</code>.</p>\n<pre class='wrapper'>\n\n<div class='filepath'>\n<p>markata.toml</p>\n\n<div
    class='right'>\n\n<button class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n</div>\n\n<div
    class=\"highlight\"><pre><span></span><span class=\"c1\">## choose your markdown
    backend</span>\n<span class=\"c1\"># markdown_backend=&#39;markdown&#39;</span>\n<span
    class=\"c1\"># markdown_backend=&#39;markdown2&#39;</span>\n<span class=\"n\">markdown_backend</span><span
    class=\"o\">=</span><span class=\"s1\">&#39;markdown-it-py&#39;</span>\n</pre></div>\n\n</pre>\n\n<h2
    id=\"markdown-it-py-configuration\">markdown-it-py configuration <a class=\"header-anchor\"
    href=\"#markdown-it-py-configuration\"><svg class=\"heading-permalink\" aria-hidden=\"true\"
    fill=\"currentColor\" focusable=\"false\" height=\"1em\" viewBox=\"0 0 24 24\"
    width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M9.199 13.599a5.99
    5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992 5.992 0
    0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242 6.003
    6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
    3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
    3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
    13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
    2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0
    0 0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
    19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997
    0 0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p><code>markdown-it-py</code>
    has quite a bit of configuration that you can do, you can read\nmore about the
    settings in their\n<a href=\"https://markdown-it-py.readthedocs.io/en/latest/\">docs</a>.</p>\n<pre
    class='wrapper'>\n\n<div class='filepath'>\n<p>markata.toml</p>\n\n<div class='right'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n</div>\n\n<div
    class=\"highlight\"><pre><span></span><span class=\"c1\"># markdown_it flavor</span>\n<span
    class=\"k\">[markata.markdown_it_py]</span>\n<span class=\"n\">config</span><span
    class=\"o\">=</span><span class=\"s1\">&#39;gfm-like&#39;</span>\n</pre></div>\n\n</pre>\n\n<p>You
    can enable and disable built in plugins using the <code>enable</code> and <code>disable</code>\nlists.</p>\n<pre
    class='wrapper'>\n\n<div class='filepath'>\n<p>markata.toml</p>\n\n<div class='right'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n</div>\n\n<div
    class=\"highlight\"><pre><span></span><span class=\"k\">[markata.markdown_it_py]</span>\n<span
    class=\"c1\"># markdown_it built-in plugins</span>\n<span class=\"n\">enable</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"p\">[</span><span class=\"w\"> </span><span class=\"s2\">&quot;table&quot;</span><span
    class=\"w\"> </span><span class=\"p\">]</span>\n<span class=\"n\">disable</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"p\">[</span><span class=\"w\"> </span><span class=\"s2\">&quot;image&quot;</span><span
    class=\"w\"> </span><span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>You
    can configure the <code>options_update</code> as follows, this is for the built-in\nplugins
    and core configuration.</p>\n<pre class='wrapper'>\n\n<div class='filepath'>\n<p>markata.toml</p>\n\n<div
    class='right'>\n\n<button class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n</div>\n\n<div
    class=\"highlight\"><pre><span></span><span class=\"k\">[markata.markdown_it_py.options_update]</span>\n<span
    class=\"n\">linkify</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"kc\">true</span>\n<span class=\"n\">html</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"kc\">true</span>\n<span class=\"n\">typographer</span><span class=\"w\">
    </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"kc\">true</span>\n<span
    class=\"n\">highlight</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s1\">&#39;markata.plugins.md_it_highlight_code:highlight_code&#39;</span>\n</pre></div>\n\n</pre>\n\n<p>Lastly
    you can add external plugins as follows.  Many of the great plugins that\ncome
    from <a href=\"https://github.com/executablebooks\">executable_books</a> actually
    comes\nfrom a separate library\n<a href=\"https://mdit-py-plugins.readthedocs.io/\">mdit-py-plugins</a>,
    so they will be\nconfigured here.</p>\n<pre class='wrapper'>\n\n<div class='filepath'>\n<p>markata.toml</p>\n\n<div
    class='right'>\n\n<button class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n</div>\n\n<div
    class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.markdown_it_py.plugins]]</span>\n<span
    class=\"n\">plugin</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;mdit_py_plugins.admon:admon_plugin&quot;</span>\n\n<span
    class=\"k\">[[markata.markdown_it_py.plugins]]</span>\n<span class=\"n\">plugin</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;mdit_py_plugins.admon:admon_plugin&quot;</span>\n\n<span class=\"k\">[[markata.markdown_it_py.plugins]]</span>\n<span
    class=\"n\">plugin</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;mdit_py_plugins.attrs:attrs_plugin&quot;</span>\n<span
    class=\"n\">config</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"p\">{</span><span class=\"n\">spans</span><span
    class=\"w\"> </span><span class=\"p\">=</span><span class=\"w\"> </span><span
    class=\"kc\">true</span><span class=\"p\">}</span>\n\n<span class=\"k\">[[markata.markdown_it_py.plugins]]</span>\n<span
    class=\"n\">plugin</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;mdit_py_plugins.attrs:attrs_block_plugin&quot;</span>\n\n<span
    class=\"k\">[[markata.markdown_it_py.plugins]]</span>\n<span class=\"n\">plugin</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;markata.plugins.mdit_details:details_plugin&quot;</span>\n\n<span
    class=\"k\">[[markata.markdown_it_py.plugins]]</span>\n<span class=\"n\">plugin</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;mdit_py_plugins.anchors:anchors_plugin&quot;</span>\n\n<span
    class=\"k\">[markata.markdown_it_py.plugins.config]</span>\n<span class=\"n\">permalink</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"kc\">true</span>\n<span class=\"n\">permalinkSymbol</span><span class=\"w\">
    </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s1\">&#39;&lt;svg
    class=&quot;heading-permalink&quot; aria-hidden=&quot;true&quot; fill=&quot;currentColor&quot;
    focusable=&quot;false&quot; height=&quot;1em&quot; viewBox=&quot;0 0 24 24&quot;
    width=&quot;1em&quot; xmlns=&quot;http://www.w3.org/2000/svg&quot;&gt;&lt;path
    d=&quot;M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
    5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
    6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
    3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
    3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
    13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
    2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0
    0 0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
    19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997
    0 0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
    1.563 1 1 0 0 0 1.602-1.198z&quot;&gt;&lt;/path&gt;&lt;/svg&gt;&#39;</span>\n\n<span
    class=\"k\">[[markata.markdown_it_py.plugins]]</span>\n<span class=\"n\">plugin</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;markata.plugins.md_it_wikilinks:wikilinks_plugin&quot;</span>\n<span
    class=\"n\">config</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"p\">{</span><span class=\"n\">markata</span><span
    class=\"w\"> </span><span class=\"p\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;markata&quot;</span><span class=\"p\">}</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='Backend' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>Backend <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">Backend
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
    <span class=\"nc\">Backend</span><span class=\"p\">(</span><span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"n\">Enum</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">markdown</span> <span class=\"o\">=</span> <span
    class=\"s2\">&quot;markdown&quot;</span>\n            <span class=\"n\">markdown2</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;markdown2&quot;</span>\n            <span
    class=\"n\">markdown_it_py</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;markdown-it-py&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='MdItExtension' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>MdItExtension <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">MdItExtension
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
    <span class=\"nc\">MdItExtension</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">plugin</span><span class=\"p\">:</span> <span class=\"nb\">str</span>\n
    \           <span class=\"n\">config</span><span class=\"p\">:</span> <span class=\"n\">Dict</span>
    <span class=\"o\">=</span> <span class=\"kc\">None</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='RenderMarkdownConfig' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>RenderMarkdownConfig <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">RenderMarkdownConfig
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
    <span class=\"nc\">RenderMarkdownConfig</span><span class=\"p\">(</span><span
    class=\"n\">pydantic</span><span class=\"o\">.</span><span class=\"n\">BaseModel</span><span
    class=\"p\">):</span>\n            <span class=\"n\">backend</span><span class=\"p\">:</span>
    <span class=\"n\">Backend</span> <span class=\"o\">=</span> <span class=\"n\">Backend</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;markdown-it-py&quot;</span><span
    class=\"p\">)</span>\n            <span class=\"n\">extensions</span><span class=\"p\">:</span>
    <span class=\"n\">List</span><span class=\"p\">[</span><span class=\"n\">MdItExtension</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"p\">[]</span>\n
    \           <span class=\"n\">cache_expire</span><span class=\"p\">:</span> <span
    class=\"nb\">int</span> <span class=\"o\">=</span> <span class=\"mi\">3600</span>\n\n
    \           <span class=\"nd\">@pydantic</span><span class=\"o\">.</span><span
    class=\"n\">validator</span><span class=\"p\">(</span><span class=\"s2\">&quot;extensions&quot;</span><span
    class=\"p\">)</span>\n            <span class=\"k\">def</span> <span class=\"nf\">convert_to_list</span><span
    class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
    class=\"n\">v</span><span class=\"p\">):</span>\n                <span class=\"k\">if</span>
    <span class=\"ow\">not</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"nb\">list</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"p\">[</span><span class=\"n\">v</span><span class=\"p\">]</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Config
    <em class='small'>class</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">Config <em class='small'>source</em></p>\n</div>\n<pre
    class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button class='copy' title='copy
    code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
    <span class=\"nc\">Config</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">render_markdown</span><span class=\"p\">:</span>
    <span class=\"n\">RenderMarkdownConfig</span> <span class=\"o\">=</span> <span
    class=\"n\">RenderMarkdownConfig</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='RenderMarkdownPost' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>RenderMarkdownPost <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">RenderMarkdownPost
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
    <span class=\"nc\">RenderMarkdownPost</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">article_html</span><span class=\"p\">:</span> <span
    class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n
    \           <span class=\"n\">html</span><span class=\"p\">:</span> <span class=\"n\">Optional</span><span
    class=\"p\">[</span><span class=\"nb\">str</span> <span class=\"o\">|</span> <span
    class=\"n\">Dict</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"nb\">str</span><span class=\"p\">]]</span>
    <span class=\"o\">=</span> <span class=\"kc\">None</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='config_model' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>config_model <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">config_model
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">config_model</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config_models</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">Config</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='post_model' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>post_model <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">post_model
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">post_model</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">post_models</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">RenderMarkdownPost</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='configure' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>configure <em class='small'>function</em></h2>\nSets up a markdown instance
    as md</p>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">configure <em class='small'>source</em></p>\n</div>\n<pre
    class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button class='copy' title='copy
    code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">configure</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"s2\">&quot;Sets up a markdown instance as md&quot;</span>\n
    \           <span class=\"c1\"># if &quot;markdown_extensions&quot; not in markata.config:</span>\n
    \           <span class=\"c1\">#     markdown_extensions = [&quot;&quot;]</span>\n
    \           <span class=\"c1\"># if isinstance(markata.config[&quot;markdown_extensions&quot;],
    str):</span>\n            <span class=\"c1\">#     markdown_extensions = [markata.config[&quot;markdown_extensions&quot;]]</span>\n
    \           <span class=\"c1\"># if isinstance(markata.config[&quot;markdown_extensions&quot;],
    list):</span>\n            <span class=\"c1\">#     markdown_extensions = markata.config[&quot;markdown_extensions&quot;]</span>\n
    \           <span class=\"c1\"># else:</span>\n            <span class=\"c1\">#
    \    raise TypeError(&quot;markdown_extensions should be List[str]&quot;)</span>\n\n
    \           <span class=\"c1\"># markata.markdown_extensions = [*DEFAULT_MD_EXTENSIONS,
    *markdown_extensions]</span>\n\n            <span class=\"k\">if</span> <span
    class=\"p\">(</span>\n                <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;markdown_backend&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"o\">.</span><span class=\"n\">lower</span><span
    class=\"p\">()</span>\n                <span class=\"o\">.</span><span class=\"n\">replace</span><span
    class=\"p\">(</span><span class=\"s2\">&quot; &quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;-&quot;</span><span class=\"p\">)</span>\n                <span
    class=\"o\">.</span><span class=\"n\">replace</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;_&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;-&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"o\">==</span> <span class=\"s2\">&quot;markdown-it-py&quot;</span>\n
    \           <span class=\"p\">):</span>\n                <span class=\"kn\">from</span>
    <span class=\"nn\">markdown_it</span> <span class=\"kn\">import</span> <span class=\"n\">MarkdownIt</span>\n\n
    \               <span class=\"n\">config_update</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;markdown_it_py&quot;</span><span class=\"p\">,</span> <span
    class=\"p\">{})</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span>\n                    <span class=\"s2\">&quot;options_update&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"p\">{</span>\n                        <span
    class=\"s2\">&quot;linkify&quot;</span><span class=\"p\">:</span> <span class=\"kc\">True</span><span
    class=\"p\">,</span>\n                        <span class=\"s2\">&quot;html&quot;</span><span
    class=\"p\">:</span> <span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \                       <span class=\"s2\">&quot;typographer&quot;</span><span
    class=\"p\">:</span> <span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \                       <span class=\"s2\">&quot;highlight&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">highlight_code</span><span class=\"p\">,</span>\n                    <span
    class=\"p\">},</span>\n                <span class=\"p\">)</span>\n                <span
    class=\"k\">if</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">config_update</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;highlight&quot;</span><span class=\"p\">),</span>
    <span class=\"nb\">str</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">module</span> <span class=\"o\">=</span> <span class=\"n\">config_update</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;highlight&quot;</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">split</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;:&quot;</span><span class=\"p\">)[</span><span class=\"mi\">0</span><span
    class=\"p\">]</span>\n                    <span class=\"n\">func</span> <span
    class=\"o\">=</span> <span class=\"n\">config_update</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;highlight&quot;</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">split</span><span class=\"p\">(</span><span class=\"s2\">&quot;:&quot;</span><span
    class=\"p\">)[</span><span class=\"mi\">1</span><span class=\"p\">]</span>\n                    <span
    class=\"n\">config_update</span><span class=\"p\">[</span><span class=\"s2\">&quot;highlight&quot;</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"nb\">getattr</span><span
    class=\"p\">(</span>\n                        <span class=\"n\">importlib</span><span
    class=\"o\">.</span><span class=\"n\">import_module</span><span class=\"p\">(</span><span
    class=\"n\">module</span><span class=\"p\">),</span>\n                        <span
    class=\"n\">func</span><span class=\"p\">,</span>\n                    <span class=\"p\">)</span>\n\n
    \               <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">md</span> <span class=\"o\">=</span> <span class=\"n\">MarkdownIt</span><span
    class=\"p\">(</span>\n                    <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;markdown_it_py&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">{})</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;config&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;gfm-like&quot;</span><span class=\"p\">),</span>\n
    \                   <span class=\"n\">config_update</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">)</span>\n                <span class=\"k\">for</span>
    <span class=\"n\">plugin</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;markdown_it_py&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">{})</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;enable&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">[]):</span>\n                    <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">md</span><span
    class=\"o\">.</span><span class=\"n\">enable</span><span class=\"p\">(</span><span
    class=\"n\">plugin</span><span class=\"p\">)</span>\n                <span class=\"k\">for</span>
    <span class=\"n\">plugin</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;markdown_it_py&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">{})</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;disable&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">[]):</span>\n                    <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">md</span><span
    class=\"o\">.</span><span class=\"n\">disable</span><span class=\"p\">(</span><span
    class=\"n\">plugin</span><span class=\"p\">)</span>\n\n                <span class=\"n\">plugins</span>
    <span class=\"o\">=</span> <span class=\"n\">copy</span><span class=\"o\">.</span><span
    class=\"n\">deepcopy</span><span class=\"p\">(</span>\n                    <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;markdown_it_py&quot;</span><span class=\"p\">,</span> <span
    class=\"p\">{})</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;plugins&quot;</span><span class=\"p\">,</span>
    <span class=\"p\">[]),</span>\n                <span class=\"p\">)</span>\n                <span
    class=\"k\">for</span> <span class=\"n\">plugin</span> <span class=\"ow\">in</span>
    <span class=\"n\">plugins</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">if</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">plugin</span><span class=\"p\">[</span><span class=\"s2\">&quot;plugin&quot;</span><span
    class=\"p\">],</span> <span class=\"nb\">str</span><span class=\"p\">):</span>\n
    \                       <span class=\"n\">plugin</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;plugin_str&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"n\">plugin</span><span class=\"p\">[</span><span class=\"s2\">&quot;plugin&quot;</span><span
    class=\"p\">]</span>\n                        <span class=\"n\">plugin_module</span>
    <span class=\"o\">=</span> <span class=\"n\">plugin</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;plugin&quot;</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">split</span><span class=\"p\">(</span><span class=\"s2\">&quot;:&quot;</span><span
    class=\"p\">)[</span><span class=\"mi\">0</span><span class=\"p\">]</span>\n                        <span
    class=\"n\">plugin_func</span> <span class=\"o\">=</span> <span class=\"n\">plugin</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;plugin&quot;</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">split</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;:&quot;</span><span class=\"p\">)[</span><span class=\"mi\">1</span><span
    class=\"p\">]</span>\n                        <span class=\"n\">plugin</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;plugin&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"nb\">getattr</span><span class=\"p\">(</span>\n
    \                           <span class=\"n\">importlib</span><span class=\"o\">.</span><span
    class=\"n\">import_module</span><span class=\"p\">(</span><span class=\"n\">plugin_module</span><span
    class=\"p\">),</span>\n                            <span class=\"n\">plugin_func</span><span
    class=\"p\">,</span>\n                        <span class=\"p\">)</span>\n                    <span
    class=\"n\">plugin</span><span class=\"p\">[</span><span class=\"s2\">&quot;config&quot;</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"n\">plugin</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;config&quot;</span><span class=\"p\">,</span> <span class=\"p\">{})</span>\n
    \                   <span class=\"k\">for</span> <span class=\"n\">k</span><span
    class=\"p\">,</span> <span class=\"n\">_v</span> <span class=\"ow\">in</span>
    <span class=\"n\">plugin</span><span class=\"p\">[</span><span class=\"s2\">&quot;config&quot;</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">items</span><span
    class=\"p\">():</span>\n                        <span class=\"k\">if</span> <span
    class=\"n\">k</span> <span class=\"o\">==</span> <span class=\"s2\">&quot;markata&quot;</span><span
    class=\"p\">:</span>\n                            <span class=\"n\">plugin</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;config&quot;</span><span class=\"p\">][</span><span
    class=\"n\">k</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"n\">markata</span>\n\n                    <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">md</span> <span class=\"o\">=</span> <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">md</span><span
    class=\"o\">.</span><span class=\"n\">use</span><span class=\"p\">(</span><span
    class=\"n\">plugin</span><span class=\"p\">[</span><span class=\"s2\">&quot;plugin&quot;</span><span
    class=\"p\">],</span> <span class=\"o\">**</span><span class=\"n\">plugin</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;config&quot;</span><span class=\"p\">])</span>\n\n
    \               <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">md</span><span class=\"o\">.</span><span class=\"n\">convert</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">md</span><span class=\"o\">.</span><span class=\"n\">render</span>\n
    \               <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">md</span><span class=\"o\">.</span><span class=\"n\">toc</span> <span
    class=\"o\">=</span> <span class=\"s2\">&quot;&quot;</span>\n            <span
    class=\"k\">elif</span> <span class=\"p\">(</span>\n                <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;markdown_backend&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"o\">.</span><span class=\"n\">lower</span><span
    class=\"p\">()</span>\n                <span class=\"o\">.</span><span class=\"n\">replace</span><span
    class=\"p\">(</span><span class=\"s2\">&quot; &quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;-&quot;</span><span class=\"p\">)</span>\n                <span
    class=\"o\">.</span><span class=\"n\">replace</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;_&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;-&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"o\">==</span> <span class=\"s2\">&quot;markdown2&quot;</span>\n
    \           <span class=\"p\">):</span>\n                <span class=\"kn\">import</span>
    <span class=\"nn\">markdown2</span>\n\n                <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">md</span> <span class=\"o\">=</span> <span
    class=\"n\">markdown2</span><span class=\"o\">.</span><span class=\"n\">Markdown</span><span
    class=\"p\">(</span>\n                    <span class=\"n\">extras</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">render_markdown</span><span
    class=\"o\">.</span><span class=\"n\">extensions</span>\n                <span
    class=\"p\">)</span>\n                <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">md</span><span class=\"o\">.</span><span class=\"n\">toc</span> <span
    class=\"o\">=</span> <span class=\"s2\">&quot;&quot;</span>\n            <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                <span class=\"kn\">import</span>
    <span class=\"nn\">markdown</span>\n\n                <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">md</span> <span class=\"o\">=</span> <span
    class=\"n\">markdown</span><span class=\"o\">.</span><span class=\"n\">Markdown</span><span
    class=\"p\">(</span>\n                    <span class=\"n\">extensions</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">render_markdown</span><span
    class=\"o\">.</span><span class=\"n\">extensions</span>\n                <span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='render'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>render <em class='small'>function</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">render
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">render</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">config</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">render_markdown</span>\n            <span class=\"k\">with</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">cache</span>
    <span class=\"k\">as</span> <span class=\"n\">cache</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">for</span> <span class=\"n\">article</span>
    <span class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">articles</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">html</span>
    <span class=\"o\">=</span> <span class=\"n\">render_article</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">,</span> <span class=\"n\">config</span><span
    class=\"p\">,</span> <span class=\"n\">cache</span><span class=\"p\">,</span>
    <span class=\"n\">article</span><span class=\"p\">)</span>\n                    <span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">article_html</span>
    <span class=\"o\">=</span> <span class=\"n\">copy</span><span class=\"o\">.</span><span
    class=\"n\">deepcopy</span><span class=\"p\">(</span><span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">html</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='render_article' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>render_article <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">render_article
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">render_article</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">config</span><span class=\"p\">,</span> <span class=\"n\">cache</span><span
    class=\"p\">,</span> <span class=\"n\">article</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">key</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">make_hash</span><span class=\"p\">(</span>\n
    \               <span class=\"s2\">&quot;render_markdown&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"s2\">&quot;render&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">content</span><span class=\"p\">,</span>\n            <span class=\"p\">)</span>\n
    \           <span class=\"n\">html_from_cache</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">precache</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">)</span>\n            <span class=\"k\">if</span>
    <span class=\"n\">html_from_cache</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"n\">html</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">md</span><span
    class=\"o\">.</span><span class=\"n\">convert</span><span class=\"p\">(</span><span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">content</span><span
    class=\"p\">)</span>\n                <span class=\"n\">cache</span><span class=\"o\">.</span><span
    class=\"n\">add</span><span class=\"p\">(</span><span class=\"n\">key</span><span
    class=\"p\">,</span> <span class=\"n\">html</span><span class=\"p\">,</span> <span
    class=\"n\">expire</span><span class=\"o\">=</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">cache_expire</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">else</span><span class=\"p\">:</span>\n                <span
    class=\"n\">html</span> <span class=\"o\">=</span> <span class=\"n\">html_from_cache</span>\n
    \           <span class=\"k\">return</span> <span class=\"n\">html</span>\n\n
    \           <span class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">html</span>
    <span class=\"o\">=</span> <span class=\"n\">html</span>\n            <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">article_html</span> <span class=\"o\">=</span>
    <span class=\"n\">copy</span><span class=\"o\">.</span><span class=\"n\">deepcopy</span><span
    class=\"p\">(</span><span class=\"n\">html</span><span class=\"p\">)</span>\n\n
    \           <span class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">html</span>
    <span class=\"o\">=</span> <span class=\"n\">html</span>\n            <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">article_html</span> <span class=\"o\">=</span>
    <span class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">article_html</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='convert_to_list' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>convert_to_list <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">convert_to_list
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">convert_to_list</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">):</span>\n                <span
    class=\"k\">if</span> <span class=\"ow\">not</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"nb\">list</span><span class=\"p\">):</span>\n                    <span
    class=\"k\">return</span> <span class=\"p\">[</span><span class=\"n\">v</span><span
    class=\"p\">]</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>"
  raw.md: ''
published: true
slug: markata/plugins/render-markdown
title: Render_Markdown.Py


---

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


!! class <h2 id='Backend' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Backend <em class='small'>class</em></h2>

???+ source "Backend <em class='small'>source</em>"

```python

        class Backend(str, Enum):
            markdown = "markdown"
            markdown2 = "markdown2"
            markdown_it_py = "markdown-it-py"
```


!! class <h2 id='MdItExtension' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>MdItExtension <em class='small'>class</em></h2>

???+ source "MdItExtension <em class='small'>source</em>"

```python

        class MdItExtension(pydantic.BaseModel):
            plugin: str
            config: Dict = None
```


!! class <h2 id='RenderMarkdownConfig' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>RenderMarkdownConfig <em class='small'>class</em></h2>

???+ source "RenderMarkdownConfig <em class='small'>source</em>"

```python

        class RenderMarkdownConfig(pydantic.BaseModel):
            backend: Backend = Backend("markdown-it-py")
            extensions: List[MdItExtension] = []
            cache_expire: int = 3600

            @pydantic.validator("extensions")
            def convert_to_list(cls, v):
                if not isinstance(v, list):
                    return [v]
                return v
```


!! class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Config <em class='small'>class</em></h2>

???+ source "Config <em class='small'>source</em>"

```python

        class Config(pydantic.BaseModel):
            render_markdown: RenderMarkdownConfig = RenderMarkdownConfig()
```


!! class <h2 id='RenderMarkdownPost' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>RenderMarkdownPost <em class='small'>class</em></h2>

???+ source "RenderMarkdownPost <em class='small'>source</em>"

```python

        class RenderMarkdownPost(pydantic.BaseModel):
            article_html: Optional[str] = None
            html: Optional[str | Dict[str, str]] = None
```


!! function <h2 id='config_model' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>config_model <em class='small'>function</em></h2>

???+ source "config_model <em class='small'>source</em>"

```python

        def config_model(markata: "Markata") -> None:
            markata.config_models.append(Config)
```


!! function <h2 id='post_model' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>post_model <em class='small'>function</em></h2>

???+ source "post_model <em class='small'>source</em>"

```python

        def post_model(markata: "Markata") -> None:
            markata.post_models.append(RenderMarkdownPost)
```


!! function <h2 id='configure' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>configure <em class='small'>function</em></h2>
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


!! function <h2 id='render' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>render <em class='small'>function</em></h2>

???+ source "render <em class='small'>source</em>"

```python

        def render(markata: "Markata") -> None:
            config = markata.config.render_markdown
            with markata.cache as cache:
                for article in markata.articles:
                    article.html = render_article(markata, config, cache, article)
                    article.article_html = copy.deepcopy(article.html)
```


!! function <h2 id='render_article' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>render_article <em class='small'>function</em></h2>

???+ source "render_article <em class='small'>source</em>"

```python

        def render_article(markata: "Markata", config, cache, article):
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
            return html

            article.html = html
            article.article_html = copy.deepcopy(html)

            article.html = html
            article.article_html = article.article_html
```


!! method <h2 id='convert_to_list' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>convert_to_list <em class='small'>method</em></h2>

???+ source "convert_to_list <em class='small'>source</em>"

```python

        def convert_to_list(cls, v):
                if not isinstance(v, list):
                    return [v]
                return v
```

