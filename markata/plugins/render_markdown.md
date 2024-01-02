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
  \           article_html: Optional[str] = None\n            html: Optional[str]
  = None\n```\n\n\n!! function <h2 id='config_model' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>config_model <em class='small'>function</em></h2>\n\n???+ source \"config_model
  <em class='small'>source</em>\"\n\n```python\n\n        def config_model(markata:
  \"Markata\") -> None:\n            markata.config_models.append(Config)\n```\n\n\n!!
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
  markata.cache as cache:\n                for article in markata.iter_articles(\"rendering
  markdown\"):\n                    key = markata.make_hash(\n                        \"render_markdown\",\n
  \                       \"render\",\n                        article.content,\n
  \                   )\n                    html_from_cache = markata.precache.get(key)\n
  \                   if html_from_cache is None:\n                        html =
  markata.md.convert(article.content)\n                        cache.add(key, html,
  expire=config.cache_expire)\n                    else:\n                        html
  = html_from_cache\n                    article.html = html\n                    article.article_html
  = copy.deepcopy(html)\n\n                    article.html = html\n                    article.article_html
  = article.article_html\n```\n\n\n!! method <h2 id='convert_to_list' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>convert_to_list <em class='small'>method</em></h2>\n\n???+
  source \"convert_to_list <em class='small'>source</em>\"\n\n```python\n\n        def
  convert_to_list(cls, v):\n                if not isinstance(v, list):\n                    return
  [v]\n                return v\n```\n"
date: 0001-01-01
description: Renders markdown content as html.  This may be markdown files loaded
  in by way There are 3 supported markdown backends that you can configure markata
  to use mar
html: "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n<title>Render_Markdown.Py</title>\n<meta
  charset=\"utf-8\"/>\n<meta content=\"width=device-width, initial-scale=1\" name=\"viewport\"/>\n<meta
  content=\"Renders markdown content as html.  This may be markdown files loaded in
  by way There are 3 supported markdown backends that you can configure markata to
  use mar\" name=\"description\"/>\n<link href=\"/static/favicon.ico\" rel=\"icon\"
  type=\"image/png\"/>\n<script>\n        function setTheme(theme) {\n            document.documentElement.setAttribute(\"data-theme\",
  theme);\n        }\n\n        function detectColorSchemeOnLoad() {\n            //local
  storage is used to override OS theme settings\n            if (localStorage.getItem(\"theme\"))
  {\n                if (localStorage.getItem(\"theme\") == \"dark\") {\n                    setTheme(\"dark\");\n
  \               } else if (localStorage.getItem(\"theme\") == \"light\") {\n                    setTheme(\"light\");\n
  \               }\n            } else if (!window.matchMedia) {\n                //matchMedia
  method not supported\n                setTheme(\"light\");\n                return
  false;\n            } else if (window.matchMedia(\"(prefers-color-scheme: dark)\").matches)
  {\n                //OS theme setting detected as dark\n                setTheme(\"dark\");\n
  \           } else {\n                setTheme(\"light\");\n            }\n        }\n
  \       detectColorSchemeOnLoad();\n        document.addEventListener(\n            \"DOMContentLoaded\",\n
  \           function () {\n                //identify the toggle switch HTML element\n
  \               const toggleSwitch = document.querySelector(\n                    '#theme-switch
  input[type=\"checkbox\"]',\n                );\n\n                //function that
  changes the theme, and sets a localStorage variable to track the theme between page
  loads\n                function switchTheme(e) {\n                    if (e.target.checked)
  {\n                        localStorage.setItem(\"theme\", \"dark\");\n                        document.documentElement.setAttribute(\"data-theme\",
  \"dark\");\n                        toggleSwitch.checked = true;\n                    }
  else {\n                        localStorage.setItem(\"theme\", \"light\");\n                        document.documentElement.setAttribute(\"data-theme\",
  \"light\");\n                        toggleSwitch.checked = false;\n                    }\n
  \               }\n\n                //listener for changing themes\n                toggleSwitch.addEventListener(\"change\",
  switchTheme, false);\n\n                //pre-check the dark-theme checkbox if dark-theme
  is set\n                if (document.documentElement.getAttribute(\"data-theme\")
  == \"dark\") {\n                    toggleSwitch.checked = true;\n                }\n
  \           },\n            false,\n        );\n    </script>\n<style>\n      :root
  {\n        --color-bg: #1f2022;\n        --color-bg-2: ;\n        --color-bg-code:
  #1f2022;\n        --color-text: #eefbfe;\n        --color-link: #fb30c4; \n        --color-accent:
  #e1bd00c9;\n        --overlay-brightness: .85;\n        --body-width: 800px;\n      }\n
  \     [data-theme=\"dark\"] {\n        --color-bg: #1f2022;\n        --color-bg-2:
  ;\n        --color-bg-code: #1f2022;\n        --color-text: #eefbfe;\n        --color-link:
  #fb30c4; \n        --color-accent: #e1bd00c9;\n        --overlay-brightness: .85;\n
  \       --body-width: 800px;\n      }\n      [data-theme=\"light\"] {\n        --color-bg:
  #eefbfe;\n        --color-bg-2: ;\n        --color-bg-code: #eefbfe;\n        --color-text:
  #1f2022;\n        --color-link: #fb30c4; \n        --color-accent: #ffeb00;\n        --overlay-brightness:
  .95;\n      }\n\n        html {\n            font-family: \"Space Mono\", monospace;\n
  \           background: var(--color-bg);\n            color: var(--color-text);\n
  \       }\n\n        a {\n            color: var(--color-link);\n        }\n\n        main
  a {\n            max-width: 100%;\n        }\n\n        .heading-permalink {\n            font-size:
  .7em;\n        }\n\n        body {\n            max-width: var(--body-width);\n
  \           margin: 5rem auto;\n            padding: 0 .5rem;\n            font-size:
  1rem;\n            line-height: 1.56;\n        }\n\n        blockquote {\n            background:
  var(--color-bg);\n            filter: brightness(var(--overlay-brightness));\n            border-left:
  4px solid var(--color-accent);\n            border-radius: 4px;\n            box-shadow:\n
  \               -0.8rem 0rem 1rem -1rem #f1fa8c,\n                0.2rem 0rem 1rem
  rgb(0, 0, 0, .4);\n            padding-left: 1rem;\n            margin: 1rem;\n
  \       }\n\n        li.post {\n            list-style-type: None;\n            padding:
  .2rem 0;\n        }\n\n        pre.wrapper {\n            padding: 0;\n            box-shadow:
  0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n            display: flex;\n            flex-direction:
  column;\n            position: relative;\n            margin: 2rem;\n        }\n\n
  \       pre {\n            margin: 0;\n            padding: 1rem;\n            min-width:
  -webkit-fill-available;\n            max-width: fit-content;\n            overflow-x:
  auto;\n        }\n\n        pre .filepath {\n            margin: 0;\n            padding-left:
  1rem;\n            border-radius: 4px 4px 0 0;\n            background: black;\n
  \           display: flex;\n            justify-content: space-between;\n            align-items:
  center;\n        }\n\n        pre .filepath p {\n            margin: 0\n        }\n\n
  \       pre .filepath .right {\n            display: flex;\n            gap: .2rem;\n
  \           align-items: center;\n        }\n\n        pre::-webkit-scrollbar {\n
  \           height: 4px;\n            background-color: transparent;\n        }\n\n
  \       pre::-webkit-scrollbar-thumb {\n            background-color: #d3d3d32e;\n
  \           border-radius: 2px;\n        }\n\n        pre::-webkit-scrollbar-track
  {\n            background-color: transparent;\n        }\n\n        .copy-wrapper
  {\n            background: none;\n            position: absolute;\n            width:
  100%;\n            z-index: 100;\n            display: flex;\n            justify-content:
  flex-end;\n        }\n\n        button.copy {\n            z-index: 100;\n            background:
  none;\n            fill: #ffffff45;\n            border: none;\n            width:
  32px;\n            align-self: flex-end;\n            top: 0;\n            right:
  0;\n            margin: 0.5rem 0.2rem;\n\n        }\n\n        button.copy:hover
  {\n            fill: white\n        }\n\n        a.help {\n            fill: #ffffff45;\n
  \       }\n\n        a.help:hover {\n            fill: white;\n        }\n\n        a.help
  svg {\n            height: 24px;\n            width: 24px;\n        }\n\n        .highlight
  {\n            background: var(--color-bg-code);\n            color: var(--color-text);\n
  \           filter: brightness(var(--overlay-brightness));\n            border-radius:
  0 0 4px 4px;\n        }\n\n        .highlight .c {\n            color: #8b8b8b\n
  \       }\n\n        /* Comment */\n        .highlight .err {\n            color:
  #960050;\n            background-color: #1e0010\n        }\n\n        /* Error */\n
  \       .highlight .k {\n            color: #c678dd\n        }\n\n        /* Keyword
  */\n        .highlight .l {\n            color: #ae81ff\n        }\n\n        /*
  Literal */\n        .highlight .n {\n            color: #abb2bf\n        }\n\n        /*
  Name */\n        .highlight .o {\n            color: #c678dd\n        }\n\n        /*
  Operator */\n        .highlight .p {\n            color: #abb2bf\n        }\n\n
  \       /* Punctuation */\n        .highlight .ch {\n            color: #8b8b8b\n
  \       }\n\n        /* Comment.Hashbang */\n        .highlight .cm {\n            color:
  #8b8b8b\n        }\n\n        /* Comment.Multiline */\n        .highlight .cp {\n
  \           color: #8b8b8b\n        }\n\n        /* Comment.Preproc */\n        .highlight
  .cpf {\n            color: #8b8b8b\n        }\n\n        /* Comment.PreprocFile
  */\n        .highlight .c1 {\n            color: #8b8b8b\n        }\n\n        /*
  Comment.Single */\n        .highlight .cs {\n            color: #8b8b8b\n        }\n\n
  \       /* Comment.Special */\n        .highlight .gd {\n            color: #c678dd\n
  \       }\n\n        /* Generic.Deleted */\n        .highlight .ge {\n            font-style:
  italic\n        }\n\n        /* Generic.Emph */\n        .highlight .gi {\n            color:
  #a6e22e\n        }\n\n        /* Generic.Inserted */\n        .highlight .gs {\n
  \           font-weight: bold\n        }\n\n        /* Generic.Strong */\n        .highlight
  .gu {\n            color: #8b8b8b\n        }\n\n        /* Generic.Subheading */\n
  \       .highlight .kc {\n            color: #c678dd\n        }\n\n        /* Keyword.Constant
  */\n        .highlight .kd {\n            color: #c678dd\n        }\n\n        /*
  Keyword.Declaration */\n        .highlight .kn {\n            color: #c678dd\n        }\n\n
  \       /* Keyword.Namespace */\n        .highlight .kp {\n            color: #c678dd\n
  \       }\n\n        /* Keyword.Pseudo */\n        .highlight .kr {\n            color:
  #c678dd\n        }\n\n        /* Keyword.Reserved */\n        .highlight .kt {\n
  \           color: #c678dd\n        }\n\n        /* Keyword.Type */\n        .highlight
  .ld {\n            color: #e6db74\n        }\n\n        /* Literal.Date */\n        .highlight
  .m {\n            color: #ae81ff\n        }\n\n        /* Literal.Number */\n        .highlight
  .s {\n            color: #e6db74\n        }\n\n        /* Literal.String */\n        .highlight
  .na {\n            color: #a6e22e\n        }\n\n        /* Name.Attribute */\n        .highlight
  .nb {\n            color: #98c379\n        }\n\n        /* Name.Builtin */\n        .highlight
  .nc {\n            color: #abb2bf\n        }\n\n        /* Name.Class */\n        .highlight
  .no {\n            color: #c678dd\n        }\n\n        /* Name.Constant */\n        .highlight
  .nd {\n            color: #abb2bf\n        }\n\n        /* Name.Decorator */\n        .highlight
  .ni {\n            color: #abb2bf\n        }\n\n        /* Name.Entity */\n        .highlight
  .ne {\n            color: #a6e22e\n        }\n\n        /* Name.Exception */\n        .highlight
  .nf {\n            color: #61afef\n        }\n\n        /* Name.Function */\n        .highlight
  .nl {\n            color: #abb2bf\n        }\n\n        /* Name.Label */\n        .highlight
  .nn {\n            color: #abb2bf\n        }\n\n        /* Name.Namespace */\n        .highlight
  .nx {\n            color: #a6e22e\n        }\n\n        /* Name.Other */\n        .highlight
  .py {\n            color: #abb2bf\n        }\n\n        /* Name.Property */\n        .highlight
  .nt {\n            color: #c678dd\n        }\n\n        /* Name.Tag */\n        .highlight
  .nv {\n            color: #abb2bf\n        }\n\n        /* Name.Variable */\n        .highlight
  .ow {\n            color: #c678dd\n        }\n\n        /* Operator.Word */\n        .highlight
  .w {\n            color: #abb2bf\n        }\n\n        /* Text.Whitespace */\n        .highlight
  .mb {\n            color: #ae81ff\n        }\n\n        /* Literal.Number.Bin */\n
  \       .highlight .mf {\n            color: #ae81ff\n        }\n\n        /* Literal.Number.Float
  */\n        .highlight .mh {\n            color: #ae81ff\n        }\n\n        /*
  Literal.Number.Hex */\n        .highlight .mi {\n            color: #ae81ff\n        }\n\n
  \       /* Literal.Number.Integer */\n        .highlight .mo {\n            color:
  #ae81ff\n        }\n\n        /* Literal.Number.Oct */\n        .highlight .sa {\n
  \           color: #e6db74\n        }\n\n        /* Literal.String.Affix */\n        .highlight
  .sb {\n            color: #e6db74\n        }\n\n        /* Literal.String.Backtick
  */\n        .highlight .sc {\n            color: #e6db74\n        }\n\n        /*
  Literal.String.Char */\n        .highlight .dl {\n            color: #e6db74\n        }\n\n
  \       /* Literal.String.Delimiter */\n        .highlight .sd {\n            color:
  #98c379\n        }\n\n        /* Literal.String.Doc */\n        .highlight .s2 {\n
  \           color: #98c379\n        }\n\n        /* Literal.String.Double */\n        .highlight
  .se {\n            color: #ae81ff\n        }\n\n        /* Literal.String.Escape
  */\n        .highlight .sh {\n            color: #e6db74\n        }\n\n        /*
  Literal.String.Heredoc */\n        .highlight .si {\n            color: #e6db74\n
  \       }\n\n        /* Literal.String.Interpol */\n        .highlight .sx {\n            color:
  #e6db74\n        }\n\n        /* Literal.String.Other */\n        .highlight .sr
  {\n            color: #e6db74\n        }\n\n        /* Literal.String.Regex */\n
  \       .highlight .s1 {\n            color: #e6db74\n        }\n\n        /* Literal.String.Single
  */\n        .highlight .ss {\n            color: #e6db74\n        }\n\n        /*
  Literal.String.Symbol */\n        .highlight .bp {\n            color: #abb2bf\n
  \       }\n\n        /* Name.Builtin.Pseudo */\n        .highlight .fm {\n            color:
  #61afef\n        }\n\n        /* Name.Function.Magic */\n        .highlight .vc
  {\n            color: #abb2bf\n        }\n\n        /* Name.Variable.Class */\n
  \       .highlight .vg {\n            color: #abb2bf\n        }\n\n        /* Name.Variable.Global
  */\n        .highlight .vi {\n            color: #abb2bf\n        }\n\n        /*
  Name.Variable.Instance */\n        .highlight .vm {\n            color: #abb2bf\n
  \       }\n\n        /* Name.Variable.Magic */\n        .highlight .il {\n            color:
  #ae81ff\n        }\n\n        /* Literal.Number.Integer.Long */\n\n        /* Tab
  style starts here */\n        .tabbed-set {\n            position: relative;\n            display:
  flex;\n            flex-wrap: wrap;\n            margin: 1em 0;\n            border-radius:
  0.1rem;\n        }\n\n        .tabbed-set>input {\n            display: none;\n
  \       }\n\n        .tabbed-set label {\n            width: auto;\n            padding:
  0.9375em 1.25em 0.78125em;\n            font-weight: 700;\n            font-size:
  0.84em;\n            white-space: nowrap;\n            border-bottom: 0.15rem solid
  transparent;\n            border-top-left-radius: 0.1rem;\n            border-top-right-radius:
  0.1rem;\n            cursor: pointer;\n            transition: background-color
  250ms, color 250ms;\n        }\n\n        .tabbed-set .tabbed-content {\n            width:
  100%;\n            display: none;\n            box-shadow: 0 -.05rem #ddd;\n        }\n\n
  \       .tabbed-set input {\n            position: absolute;\n            opacity:
  0;\n        }\n\n        /* fonts */\n        h1 {\n            font-weight: 700;\n
  \       }\n\n        h1#title a {\n            font-size: 16px;\n        }\n\n        h1,\n
  \       h2,\n        h3,\n        h4,\n        h5,\n        h6 {\n            margin-top:
  3rem;\n        }\n\n        h1 {\n            font-size: 2.5em;\n            margin-top:
  5rem;\n        }\n\n        h2 {\n            font-size: 1.63rem;\n            margin-top:
  5rem;\n        }\n\n\n\n        p {\n            font-size: 21px;\n            font-style:
  normal;\n            font-variant: normal;\n            font-weight: 400;\n            line-height:
  1.5;\n        }\n\n        @media only screen and (max-width: 700px) {\n            p
  {\n                font-size: 18px;\n            }\n        }\n\n        @media
  only screen and (max-width: 600px) {\n            p {\n                font-size:
  16px;\n            }\n        }\n\n        @media only screen and (max-width: 500px)
  {\n            p {\n                font-size: 14px;\n            }\n        }\n\n
  \       @media only screen and (max-width: 400px) {\n            p {\n                font-size:
  12px;\n            }\n        }\n\n\n        pre {\n            font-style: normal;\n
  \           font-variant: normal;\n            font-weight: 400;\n            line-height:
  18.5714px;\n            */\n        }\n\n        a {\n            font-weight: 600;\n
  \           text-decoration-color: var(--color-accent);\n            color: var(--color-link);\n
  \           padding: .3rem .5rem;\n            display: inline-block;\n        }\n\n
  \       .admonition,\n        details {\n            box-shadow: 0.2rem 0rem 1rem
  rgb(0, 0, 0, .4);\n            margin: 5rem 0;\n            border: 1px solid transparent;\n
  \           border-radius: 4px;\n            text-align: left;\n            padding:
  0;\n            border: 0;\n\n        }\n\n        .admonition {\n            padding-bottom:
  1rem;\n        }\n\n        details[open] {\n            padding-bottom: .5rem;\n
  \       }\n\n        .admonition p {\n            padding: .2rem .6rem;\n        }\n\n
  \       .admonition-title,\n        .details-title,\n        summary {\n            background:
  var(--color-bg-2);\n            padding: 0;\n            margin: 0;\n            position:
  sticky;\n            top: 0;\n            z-index: 10;\n        }\n\n        summary:hover
  {\n            cursor: pointer;\n        }\n\n        summary.admonition-title,\n
  \       summary.details-title {\n            padding: .5rem;\n            padding-left:
  1rem;\n        }\n\n        .note {\n            border-left: 4px solid #f1fa8c;\n
  \           box-shadow:\n                -0.8rem 0rem 1rem -1rem #f1fa8c,\n                0.2rem
  0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n        .note>.admonition-title {\n            border-bottom:
  1px solid #3c3d2d;\n        }\n\n        .abstract {\n            border-left: 4px
  solid #8be9fd;\n            box-shadow:\n                -0.8rem 0rem 1rem -1rem
  #8be9fd,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n        .abstract>.admonition-title
  {\n            border-bottom: 1px solid #2c3a3f;\n        }\n\n        .info {\n
  \           border-left: 4px solid;\n            box-shadow:\n                -0.8rem
  0rem 1rem -1rem #8bb0fd,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n
  \       .info>.admonition-title {\n            border-bottom: 1px solid #2c313f;\n
  \       }\n\n        .tip {\n            border-left: 4px solid #008080;\n            box-shadow:\n
  \               -0.8rem 0rem 1rem -1rem #008080,\n                0.2rem 0rem 1rem
  rgb(0, 0, 0, .4);\n        }\n\n        .tip>.admonition-title {\n            border-bottom:
  1px solid #1b2a2b;\n        }\n\n        .success {\n            border-left: 4px
  solid #50fa7b;\n            box-shadow:\n                -0.8rem 0rem 1rem -1rem
  #50fa7b,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n        .success>.admonition-title
  {\n            border-bottom: 1px solid #263e2b;\n        }\n\n        .question
  {\n            border-left: 4px solid #a7fcbd;\n            box-shadow:\n                -0.8rem
  0rem 1rem -1rem #a7fcbd,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n
  \       .question>.admonition-title {\n            border-bottom: 1px solid #303e35;\n
  \       }\n\n        .warning {\n            border-left: 4px solid #ffb86c;\n            box-shadow:\n
  \               -0.8rem 0rem 1rem -1rem #ffb86c,\n                0.2rem 0rem 1rem
  rgb(0, 0, 0, .4);\n        }\n\n        .warning>.admonition-title {\n            border-bottom:
  1px solid #3f3328;\n        }\n\n        .failure {\n            border-left: 4px
  solid #b23b3b;\n            box-shadow:\n                -0.8rem 0rem 1rem -1rem
  #b23b3b,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n        .failure>.admonition-title
  {\n            border-bottom: 1px solid #34201f;\n        }\n\n        .danger {\n
  \           border-left: 4px solid #ff5555;\n            box-shadow:\n                -0.8rem
  0rem 1rem -1rem #ff5555,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n
  \       .danger>.admonition-title {\n            border-bottom: 1px solid #402523;\n
  \       }\n\n        .bug {\n            border-left: 4px solid #b2548a;\n            box-shadow:\n
  \               -0.8rem 0rem 1rem -1rem #b2548a,\n                0.2rem 0rem 1rem
  rgb(0, 0, 0, .4);\n        }\n\n        .bug>.admonition-title {\n            border-bottom:
  1px solid #32232c;\n        }\n\n        .example {\n            border-left: 4px
  solid #bd93f9;\n            box-shadow:\n                -0.8rem 0rem 1rem -1rem
  #bd93f9,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n        .example>.admonition-title
  {\n            border-bottom: 1px solid #332d3e;\n        }\n\n        .source {\n
  \           border-left: 4px solid #bd93f9;\n            box-shadow:\n                -0.8rem
  0rem 1rem -1rem #bd93f9,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n
  \       .source>.admonition-title {\n            border-bottom: 1px solid #332d3e;\n
  \       }\n\n        .quote {\n            border-left: 4px solid #999;\n            box-shadow:\n
  \               -0.8rem 0rem 1rem -1rem #999,\n                0.2rem 0rem 1rem
  rgb(0, 0, 0, .4);\n        }\n\n        .quote>.admonition-title {\n            border-bottom:
  1px solid #2d2e2f;\n        }\n\n        table {\n            margin: 1rem 0;\n
  \           border-collapse: collapse;\n            border-spacing: 0;\n            display:
  block;\n            max-width: -moz-fit-content;\n            max-width: fit-content;\n
  \           overflow-x: auto;\n            white-space: nowrap;\n        }\n\n        table
  thead th {\n            border: solid 1px var(--color-text);\n            padding:
  10px;\n            text-align: left;\n        }\n\n        table tbody td {\n            border:
  solid 1px var(--color-text);\n            padding: 10px;\n        }\n\n        .theme-switch
  {\n            z-index: 10;\n            display: inline-block;\n            height:
  34px;\n            position: relative;\n            width: 60px;\n\n            display:
  flex;\n            justify-content: flex-end;\n            margin-right: 1rem;\n
  \           margin-left: auto;\n            position: fixed;\n            right:
  1rem;\n            top: 1rem;\n        }\n\n        .theme-switch input {\n            display:
  none;\n\n        }\n\n        .slider {\n            background-color: #ccc;\n            bottom:
  0;\n            cursor: pointer;\n            left: 0;\n            position: absolute;\n
  \           right: 0;\n            top: 0;\n            transition: .4s;\n        }\n\n
  \       .slider:before {\n            background-color: #fff;\n            bottom:
  4px;\n            content: \"\";\n            height: 26px;\n            left: 4px;\n
  \           position: absolute;\n            transition: .4s;\n            width:
  26px;\n        }\n\n        input:checked+.slider {\n            background-color:
  #343434;\n        }\n\n        input:checked+.slider:before {\n            background-color:
  #848484;\n        }\n\n        input:checked+.slider:before {\n            transform:
  translateX(26px);\n        }\n\n        .slider.round {\n            border-radius:
  34px;\n        }\n\n        .slider.round:before {\n            border-radius: 50%;\n
  \       }\n\n        main p img {\n            width: 100%;\n            width:
  -moz-available;\n            width: -webkit-fill-available;\n            width:
  fill-available;\n        }\n\n        details>* {\n            margin: 1rem;\n        }\n\n
  \       .admonition>* {\n            margin: 1rem;\n        }\n\n        p.admonition-title,\n
  \       summary {\n            margin: 0;\n            padding-left: 1.2rem;\n        }\n\n
  \       .small {\n            font-size: .9rem;\n            color: #888;\n        }\n\n
  \       admonition+admonition {\n            margin-top: 20rem;\n        }\n\n        ::-webkit-scrollbar
  {\n            height: 12px;\n            background-color: transparent;\n        }\n\n
  \       ::-webkit-scrollbar-thumb {\n            background-color: #d3d3d32e;\n
  \           border-radius: 6px;\n        }\n\n        ::-webkit-scrollbar-track
  {\n            background-color: transparent;\n        }\n    </style>\n<script>\n
  \       if (\"serviceWorker\" in navigator) {\n            navigator.serviceWorker.register(\"/service-worker.js\");\n
  \           navigator.serviceWorker.addEventListener(\"controllerchange\", () =>
  {\n                console.log(\"new worker\");\n                window.location.reload();\n
  \           });\n        }\n    </script>\n<meta content=\"waylon@waylonwalker.com\"
  name=\"og:author_email\"/>\n<meta content=\"waylon@waylonwalker.com\" name=\"og:author_email\"/>\n<meta
  content=\"Waylon Walker\" name=\"og:author\" property=\"og:author\"/><meta content=\"waylon@waylonwalaker.com\"
  name=\"og:author_email\" property=\"og:author_email\"/><meta content=\"website\"
  name=\"og:type\" property=\"og:type\"/><meta content=\"Renders markdown content
  as html.  This may be markdown files loaded in by way There are 3 supported markdown
  backends that you can configure markata to use mar\" name=\"description\" property=\"description\"/><meta
  content=\"Renders markdown content as html.  This may be markdown files loaded in
  by way There are 3 supported markdown backends that you can configure markata to
  use mar\" name=\"og:description\" property=\"og:description\"/><meta content=\"Renders
  markdown content as html.  This may be markdown files loaded in by way There are
  3 supported markdown backends that you can configure markata to use mar\" name=\"twitter:description\"
  property=\"twitter:description\"/><meta content=\"Render_Markdown.Py | Markata\"
  name=\"og:title\" property=\"og:title\"/><meta content=\"Render_Markdown.Py | Markata\"
  name=\"twitter:title\" property=\"twitter:title\"/><meta content=\"https://markata.dev//markata/plugins/render-markdown-og.png\"
  name=\"og:image\" property=\"og:image\"/><meta content=\"https://markata.dev//markata/plugins/render-markdown-og.png\"
  name=\"twitter:image\" property=\"twitter:image\"/><meta content=\"1600\" name=\"og:image:width\"
  property=\"og:image:width\"/><meta content=\"900\" name=\"og:image:width\" property=\"og:image:width\"/><meta
  content=\"summary_large_image\" name=\"twitter:card\" property=\"twitter:card\"/><meta
  content=\"Markata\" name=\"og:site_name\" property=\"og:site_name\"/><meta content=\"@_waylonwalker\"
  name=\"twitter:creator\" property=\"twitter:creator\"/><meta content=\"Render_Markdown.Py\"
  name=\"title\" property=\"title\"/><meta content=\"markata 0.8.0.dev6\" name=\"generator\"
  property=\"generator\"/><link href=\"https://markata.dev//markata/plugins/render-markdown/\"
  rel=\"canonical\"/><meta content=\"https://markata.dev//markata/plugins/render-markdown/\"
  name=\"og:url\" property=\"og:url\"/></head>\n<body><nav>\n<a href=\"/\">markata</a>\n<a
  href=\"https://github.com/WaylonWalker/markata\">GitHub</a>\n</nav>\n<div>\n<label
  class=\"theme-switch\" for=\"checkbox-theme\" id=\"theme-switch\" title=\"light/dark
  mode toggle\">\n<input id=\"checkbox-theme\" type=\"checkbox\"/>\n<div class=\"slider
  round\"></div>\n</label>\n</div>\n<section class=\"title\">\n<h1 id=\"title\">\n
  \           Render_Markdown.Py \n            \n        </h1>\n</section>\n<main><p>Renders
  markdown content as html.  This may be markdown files loaded in by way\nof the <a
  class=\"wikilink\" href=\"/markata/plugins/load\">load</a> plugin.</p>\n<h2 id=\"markdown-backend\">Markdown
  backend <a class=\"header-anchor\" href=\"#markdown-backend\"><svg aria-hidden=\"true\"
  class=\"heading-permalink\" fill=\"currentColor\" focusable=\"false\" height=\"1em\"
  viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M9.199
  13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>There are 3 supported
  markdown backends that you can configure markata to use\nby setting the <code>markdown_backend</code>
  in your <code>markata.toml</code>.</p>\n<pre class=\"wrapper\">\n\n<div class=\"filepath\">\n<p>markata.toml</p>\n\n<div
  class=\"right\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n</div>\n\n<div
  class=\"highlight\"><pre><span></span><span class=\"c1\">## choose your markdown
  backend</span>\n<span class=\"c1\"># markdown_backend='markdown'</span>\n<span class=\"c1\">#
  markdown_backend='markdown2'</span>\n<span class=\"n\">markdown_backend</span><span
  class=\"o\">=</span><span class=\"s1\">'markdown-it-py'</span>\n</pre></div>\n\n</pre>\n<h2
  id=\"markdown-it-py-configuration\">markdown-it-py configuration <a class=\"header-anchor\"
  href=\"#markdown-it-py-configuration\"><svg aria-hidden=\"true\" class=\"heading-permalink\"
  fill=\"currentColor\" focusable=\"false\" height=\"1em\" viewbox=\"0 0 24 24\" width=\"1em\"
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
  3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632 1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p><code>markdown-it-py</code>
  has quite a bit of configuration that you can do, you can read\nmore about the settings
  in their\n<a href=\"https://markdown-it-py.readthedocs.io/en/latest/\">docs</a>.</p>\n<pre
  class=\"wrapper\">\n\n<div class=\"filepath\">\n<p>markata.toml</p>\n\n<div class=\"right\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n</div>\n\n<div
  class=\"highlight\"><pre><span></span><span class=\"c1\"># markdown_it flavor</span>\n<span
  class=\"k\">[markata.markdown_it_py]</span>\n<span class=\"n\">config</span><span
  class=\"o\">=</span><span class=\"s1\">'gfm-like'</span>\n</pre></div>\n\n</pre>\n<p>You
  can enable and disable built in plugins using the <code>enable</code> and <code>disable</code>\nlists.</p>\n<pre
  class=\"wrapper\">\n\n<div class=\"filepath\">\n<p>markata.toml</p>\n\n<div class=\"right\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n</div>\n\n<div
  class=\"highlight\"><pre><span></span><span class=\"k\">[markata.markdown_it_py]</span>\n<span
  class=\"c1\"># markdown_it built-in plugins</span>\n<span class=\"n\">enable</span><span
  class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"p\">[</span><span
  class=\"w\"> </span><span class=\"s2\">\"table\"</span><span class=\"w\"> </span><span
  class=\"p\">]</span>\n<span class=\"n\">disable</span><span class=\"w\"> </span><span
  class=\"o\">=</span><span class=\"w\"> </span><span class=\"p\">[</span><span class=\"w\">
  </span><span class=\"s2\">\"image\"</span><span class=\"w\"> </span><span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n<p>You
  can configure the <code>options_update</code> as follows, this is for the built-in\nplugins
  and core configuration.</p>\n<pre class=\"wrapper\">\n\n<div class=\"filepath\">\n<p>markata.toml</p>\n\n<div
  class=\"right\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n</div>\n\n<div
  class=\"highlight\"><pre><span></span><span class=\"k\">[markata.markdown_it_py.options_update]</span>\n<span
  class=\"n\">linkify</span><span class=\"w\"> </span><span class=\"o\">=</span><span
  class=\"w\"> </span><span class=\"kc\">true</span>\n<span class=\"n\">html</span><span
  class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"kc\">true</span>\n<span
  class=\"n\">typographer</span><span class=\"w\"> </span><span class=\"o\">=</span><span
  class=\"w\"> </span><span class=\"kc\">true</span>\n<span class=\"n\">highlight</span><span
  class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s1\">'markata.plugins.md_it_highlight_code:highlight_code'</span>\n</pre></div>\n\n</pre>\n<p>Lastly
  you can add external plugins as follows.  Many of the great plugins that\ncome from
  <a href=\"https://github.com/executablebooks\">executable_books</a> actually comes\nfrom
  a separate library\n<a href=\"https://mdit-py-plugins.readthedocs.io/\">mdit-py-plugins</a>,
  so they will be\nconfigured here.</p>\n<pre class=\"wrapper\">\n\n<div class=\"filepath\">\n<p>markata.toml</p>\n\n<div
  class=\"right\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n</div>\n\n<div
  class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.markdown_it_py.plugins]]</span>\n<span
  class=\"n\">plugin</span><span class=\"w\"> </span><span class=\"o\">=</span><span
  class=\"w\"> </span><span class=\"s2\">\"mdit_py_plugins.admon:admon_plugin\"</span>\n\n<span
  class=\"k\">[[markata.markdown_it_py.plugins]]</span>\n<span class=\"n\">plugin</span><span
  class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s2\">\"mdit_py_plugins.admon:admon_plugin\"</span>\n\n<span
  class=\"k\">[[markata.markdown_it_py.plugins]]</span>\n<span class=\"n\">plugin</span><span
  class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s2\">\"mdit_py_plugins.attrs:attrs_plugin\"</span>\n<span
  class=\"n\">config</span><span class=\"w\"> </span><span class=\"o\">=</span><span
  class=\"w\"> </span><span class=\"p\">{</span><span class=\"n\">spans</span><span
  class=\"w\"> </span><span class=\"p\">=</span><span class=\"w\"> </span><span class=\"kc\">true</span><span
  class=\"p\">}</span>\n\n<span class=\"k\">[[markata.markdown_it_py.plugins]]</span>\n<span
  class=\"n\">plugin</span><span class=\"w\"> </span><span class=\"o\">=</span><span
  class=\"w\"> </span><span class=\"s2\">\"mdit_py_plugins.attrs:attrs_block_plugin\"</span>\n\n<span
  class=\"k\">[[markata.markdown_it_py.plugins]]</span>\n<span class=\"n\">plugin</span><span
  class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s2\">\"markata.plugins.mdit_details:details_plugin\"</span>\n\n<span
  class=\"k\">[[markata.markdown_it_py.plugins]]</span>\n<span class=\"n\">plugin</span><span
  class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s2\">\"mdit_py_plugins.anchors:anchors_plugin\"</span>\n\n<span
  class=\"k\">[markata.markdown_it_py.plugins.config]</span>\n<span class=\"n\">permalink</span><span
  class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"kc\">true</span>\n<span
  class=\"n\">permalinkSymbol</span><span class=\"w\"> </span><span class=\"o\">=</span><span
  class=\"w\"> </span><span class=\"s1\">'&lt;svg class=\"heading-permalink\" aria-hidden=\"true\"
  fill=\"currentColor\" focusable=\"false\" height=\"1em\" viewBox=\"0 0 24 24\" width=\"1em\"
  xmlns=\"http://www.w3.org/2000/svg\"&gt;&lt;path d=\"M9.199 13.599a5.99 5.99 0 0
  0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992 5.992 0 0 0 1.695-4.285
  5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242 6.003 6.003 0 0 0-1.905
  1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985 3.985 0 0 1 2.761-1.123
  3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005 3.006a3.982 3.982
  0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201 13.6zm5.602-3.198a5.99
  5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995 2.994a5.992 5.992 0
  0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0 0 6.431 1.242 6.003
  6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836 19.81a3.985
  3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0 0 1-.111-5.644l3.005-3.006a3.982
  3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632 1.563 1 1 0 0 0 1.602-1.198z\"&gt;&lt;/path&gt;&lt;/svg&gt;'</span>\n\n<span
  class=\"k\">[[markata.markdown_it_py.plugins]]</span>\n<span class=\"n\">plugin</span><span
  class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s2\">\"markata.plugins.md_it_wikilinks:wikilinks_plugin\"</span>\n<span
  class=\"n\">config</span><span class=\"w\"> </span><span class=\"o\">=</span><span
  class=\"w\"> </span><span class=\"p\">{</span><span class=\"n\">markata</span><span
  class=\"w\"> </span><span class=\"p\">=</span><span class=\"w\"> </span><span class=\"s2\">\"markata\"</span><span
  class=\"p\">}</span>\n</pre></div>\n\n</pre>\n<p>!! class </p><h2 class=\"admonition-title\"
  id=\"Backend\" style=\"margin:0;padding:.5rem 1rem;\">Backend <em class=\"small\">class</em></h2>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"Backend
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
  <span class=\"nc\">Backend</span><span class=\"p\">(</span><span class=\"nb\">str</span><span
  class=\"p\">,</span> <span class=\"n\">Enum</span><span class=\"p\">):</span>\n
  \           <span class=\"n\">markdown</span> <span class=\"o\">=</span> <span class=\"s2\">\"markdown\"</span>\n
  \           <span class=\"n\">markdown2</span> <span class=\"o\">=</span> <span
  class=\"s2\">\"markdown2\"</span>\n            <span class=\"n\">markdown_it_py</span>
  <span class=\"o\">=</span> <span class=\"s2\">\"markdown-it-py\"</span>\n</pre></div>\n\n</pre>\n<p>!!
  class </p><h2 class=\"admonition-title\" id=\"MdItExtension\" style=\"margin:0;padding:.5rem
  1rem;\">MdItExtension <em class=\"small\">class</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"MdItExtension
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
  <span class=\"nc\">MdItExtension</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
  class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
  \           <span class=\"n\">plugin</span><span class=\"p\">:</span> <span class=\"nb\">str</span>\n
  \           <span class=\"n\">config</span><span class=\"p\">:</span> <span class=\"n\">Dict</span>
  <span class=\"o\">=</span> <span class=\"kc\">None</span>\n</pre></div>\n\n</pre>\n<p>!!
  class </p><h2 class=\"admonition-title\" id=\"RenderMarkdownConfig\" style=\"margin:0;padding:.5rem
  1rem;\">RenderMarkdownConfig <em class=\"small\">class</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"RenderMarkdownConfig
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
  <span class=\"nc\">RenderMarkdownConfig</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
  class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
  \           <span class=\"n\">backend</span><span class=\"p\">:</span> <span class=\"n\">Backend</span>
  <span class=\"o\">=</span> <span class=\"n\">Backend</span><span class=\"p\">(</span><span
  class=\"s2\">\"markdown-it-py\"</span><span class=\"p\">)</span>\n            <span
  class=\"n\">extensions</span><span class=\"p\">:</span> <span class=\"n\">List</span><span
  class=\"p\">[</span><span class=\"n\">MdItExtension</span><span class=\"p\">]</span>
  <span class=\"o\">=</span> <span class=\"p\">[]</span>\n            <span class=\"n\">cache_expire</span><span
  class=\"p\">:</span> <span class=\"nb\">int</span> <span class=\"o\">=</span> <span
  class=\"mi\">3600</span>\n\n            <span class=\"nd\">@pydantic</span><span
  class=\"o\">.</span><span class=\"n\">validator</span><span class=\"p\">(</span><span
  class=\"s2\">\"extensions\"</span><span class=\"p\">)</span>\n            <span
  class=\"k\">def</span> <span class=\"nf\">convert_to_list</span><span class=\"p\">(</span><span
  class=\"bp\">cls</span><span class=\"p\">,</span> <span class=\"n\">v</span><span
  class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"ow\">not</span>
  <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">v</span><span
  class=\"p\">,</span> <span class=\"nb\">list</span><span class=\"p\">):</span>\n
  \                   <span class=\"k\">return</span> <span class=\"p\">[</span><span
  class=\"n\">v</span><span class=\"p\">]</span>\n                <span class=\"k\">return</span>
  <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n<p>!! class </p><h2 class=\"admonition-title\"
  id=\"Config\" style=\"margin:0;padding:.5rem 1rem;\">Config <em class=\"small\">class</em></h2>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"Config
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
  <span class=\"nc\">Config</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
  class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
  \           <span class=\"n\">render_markdown</span><span class=\"p\">:</span> <span
  class=\"n\">RenderMarkdownConfig</span> <span class=\"o\">=</span> <span class=\"n\">RenderMarkdownConfig</span><span
  class=\"p\">()</span>\n</pre></div>\n\n</pre>\n<p>!! class </p><h2 class=\"admonition-title\"
  id=\"RenderMarkdownPost\" style=\"margin:0;padding:.5rem 1rem;\">RenderMarkdownPost
  <em class=\"small\">class</em></h2>\n<div class=\"admonition source is-collapsible
  collapsible-open\">\n<p class=\"admonition-title\">\"RenderMarkdownPost <em class=\"small\">source</em>\"</p>\n</div>\n<pre
  class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
  <span class=\"nc\">RenderMarkdownPost</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
  class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
  \           <span class=\"n\">article_html</span><span class=\"p\">:</span> <span
  class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
  class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n
  \           <span class=\"n\">html</span><span class=\"p\">:</span> <span class=\"n\">Optional</span><span
  class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">]</span> <span
  class=\"o\">=</span> <span class=\"kc\">None</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"config_model\" style=\"margin:0;padding:.5rem
  1rem;\">config_model <em class=\"small\">function</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"config_model
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">config_model</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config_models</span><span
  class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
  class=\"n\">Config</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"post_model\" style=\"margin:0;padding:.5rem
  1rem;\">post_model <em class=\"small\">function</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"post_model
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">post_model</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">post_models</span><span
  class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
  class=\"n\">RenderMarkdownPost</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"configure\" style=\"margin:0;padding:.5rem
  1rem;\">configure <em class=\"small\">function</em></h2>\nSets up a markdown instance
  as md\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"configure
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">configure</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \           <span class=\"s2\">\"Sets up a markdown instance as md\"</span>\n            <span
  class=\"c1\"># if \"markdown_extensions\" not in markata.config:</span>\n            <span
  class=\"c1\">#     markdown_extensions = [\"\"]</span>\n            <span class=\"c1\">#
  if isinstance(markata.config[\"markdown_extensions\"], str):</span>\n            <span
  class=\"c1\">#     markdown_extensions = [markata.config[\"markdown_extensions\"]]</span>\n
  \           <span class=\"c1\"># if isinstance(markata.config[\"markdown_extensions\"],
  list):</span>\n            <span class=\"c1\">#     markdown_extensions = markata.config[\"markdown_extensions\"]</span>\n
  \           <span class=\"c1\"># else:</span>\n            <span class=\"c1\">#
  \    raise TypeError(\"markdown_extensions should be List[str]\")</span>\n\n            <span
  class=\"c1\"># markata.markdown_extensions = [*DEFAULT_MD_EXTENSIONS, *markdown_extensions]</span>\n\n
  \           <span class=\"k\">if</span> <span class=\"p\">(</span>\n                <span
  class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
  class=\"s2\">\"markdown_backend\"</span><span class=\"p\">,</span> <span class=\"s2\">\"\"</span><span
  class=\"p\">)</span>\n                <span class=\"o\">.</span><span class=\"n\">lower</span><span
  class=\"p\">()</span>\n                <span class=\"o\">.</span><span class=\"n\">replace</span><span
  class=\"p\">(</span><span class=\"s2\">\" \"</span><span class=\"p\">,</span> <span
  class=\"s2\">\"-\"</span><span class=\"p\">)</span>\n                <span class=\"o\">.</span><span
  class=\"n\">replace</span><span class=\"p\">(</span><span class=\"s2\">\"_\"</span><span
  class=\"p\">,</span> <span class=\"s2\">\"-\"</span><span class=\"p\">)</span>\n
  \               <span class=\"o\">==</span> <span class=\"s2\">\"markdown-it-py\"</span>\n
  \           <span class=\"p\">):</span>\n                <span class=\"kn\">from</span>
  <span class=\"nn\">markdown_it</span> <span class=\"kn\">import</span> <span class=\"n\">MarkdownIt</span>\n\n
  \               <span class=\"n\">config_update</span> <span class=\"o\">=</span>
  <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
  class=\"s2\">\"markdown_it_py\"</span><span class=\"p\">,</span> <span class=\"p\">{})</span><span
  class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span>\n                    <span
  class=\"s2\">\"options_update\"</span><span class=\"p\">,</span>\n                    <span
  class=\"p\">{</span>\n                        <span class=\"s2\">\"linkify\"</span><span
  class=\"p\">:</span> <span class=\"kc\">True</span><span class=\"p\">,</span>\n
  \                       <span class=\"s2\">\"html\"</span><span class=\"p\">:</span>
  <span class=\"kc\">True</span><span class=\"p\">,</span>\n                        <span
  class=\"s2\">\"typographer\"</span><span class=\"p\">:</span> <span class=\"kc\">True</span><span
  class=\"p\">,</span>\n                        <span class=\"s2\">\"highlight\"</span><span
  class=\"p\">:</span> <span class=\"n\">highlight_code</span><span class=\"p\">,</span>\n
  \                   <span class=\"p\">},</span>\n                <span class=\"p\">)</span>\n
  \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
  class=\"p\">(</span><span class=\"n\">config_update</span><span class=\"o\">.</span><span
  class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">\"highlight\"</span><span
  class=\"p\">),</span> <span class=\"nb\">str</span><span class=\"p\">):</span>\n
  \                   <span class=\"n\">module</span> <span class=\"o\">=</span> <span
  class=\"n\">config_update</span><span class=\"p\">[</span><span class=\"s2\">\"highlight\"</span><span
  class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">split</span><span
  class=\"p\">(</span><span class=\"s2\">\":\"</span><span class=\"p\">)[</span><span
  class=\"mi\">0</span><span class=\"p\">]</span>\n                    <span class=\"n\">func</span>
  <span class=\"o\">=</span> <span class=\"n\">config_update</span><span class=\"p\">[</span><span
  class=\"s2\">\"highlight\"</span><span class=\"p\">]</span><span class=\"o\">.</span><span
  class=\"n\">split</span><span class=\"p\">(</span><span class=\"s2\">\":\"</span><span
  class=\"p\">)[</span><span class=\"mi\">1</span><span class=\"p\">]</span>\n                    <span
  class=\"n\">config_update</span><span class=\"p\">[</span><span class=\"s2\">\"highlight\"</span><span
  class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"nb\">getattr</span><span
  class=\"p\">(</span>\n                        <span class=\"n\">importlib</span><span
  class=\"o\">.</span><span class=\"n\">import_module</span><span class=\"p\">(</span><span
  class=\"n\">module</span><span class=\"p\">),</span>\n                        <span
  class=\"n\">func</span><span class=\"p\">,</span>\n                    <span class=\"p\">)</span>\n\n
  \               <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">md</span> <span class=\"o\">=</span> <span class=\"n\">MarkdownIt</span><span
  class=\"p\">(</span>\n                    <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
  class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">\"markdown_it_py\"</span><span
  class=\"p\">,</span> <span class=\"p\">{})</span><span class=\"o\">.</span><span
  class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">\"config\"</span><span
  class=\"p\">,</span> <span class=\"s2\">\"gfm-like\"</span><span class=\"p\">),</span>\n
  \                   <span class=\"n\">config_update</span><span class=\"p\">,</span>\n
  \               <span class=\"p\">)</span>\n                <span class=\"k\">for</span>
  <span class=\"n\">plugin</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
  class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">\"markdown_it_py\"</span><span
  class=\"p\">,</span> <span class=\"p\">{})</span><span class=\"o\">.</span><span
  class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">\"enable\"</span><span
  class=\"p\">,</span> <span class=\"p\">[]):</span>\n                    <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">md</span><span class=\"o\">.</span><span class=\"n\">enable</span><span
  class=\"p\">(</span><span class=\"n\">plugin</span><span class=\"p\">)</span>\n
  \               <span class=\"k\">for</span> <span class=\"n\">plugin</span> <span
  class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">get</span><span
  class=\"p\">(</span><span class=\"s2\">\"markdown_it_py\"</span><span class=\"p\">,</span>
  <span class=\"p\">{})</span><span class=\"o\">.</span><span class=\"n\">get</span><span
  class=\"p\">(</span><span class=\"s2\">\"disable\"</span><span class=\"p\">,</span>
  <span class=\"p\">[]):</span>\n                    <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">md</span><span class=\"o\">.</span><span class=\"n\">disable</span><span
  class=\"p\">(</span><span class=\"n\">plugin</span><span class=\"p\">)</span>\n\n
  \               <span class=\"n\">plugins</span> <span class=\"o\">=</span> <span
  class=\"n\">copy</span><span class=\"o\">.</span><span class=\"n\">deepcopy</span><span
  class=\"p\">(</span>\n                    <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
  class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">\"markdown_it_py\"</span><span
  class=\"p\">,</span> <span class=\"p\">{})</span><span class=\"o\">.</span><span
  class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">\"plugins\"</span><span
  class=\"p\">,</span> <span class=\"p\">[]),</span>\n                <span class=\"p\">)</span>\n
  \               <span class=\"k\">for</span> <span class=\"n\">plugin</span> <span
  class=\"ow\">in</span> <span class=\"n\">plugins</span><span class=\"p\">:</span>\n
  \                   <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
  class=\"p\">(</span><span class=\"n\">plugin</span><span class=\"p\">[</span><span
  class=\"s2\">\"plugin\"</span><span class=\"p\">],</span> <span class=\"nb\">str</span><span
  class=\"p\">):</span>\n                        <span class=\"n\">plugin</span><span
  class=\"p\">[</span><span class=\"s2\">\"plugin_str\"</span><span class=\"p\">]</span>
  <span class=\"o\">=</span> <span class=\"n\">plugin</span><span class=\"p\">[</span><span
  class=\"s2\">\"plugin\"</span><span class=\"p\">]</span>\n                        <span
  class=\"n\">plugin_module</span> <span class=\"o\">=</span> <span class=\"n\">plugin</span><span
  class=\"p\">[</span><span class=\"s2\">\"plugin\"</span><span class=\"p\">]</span><span
  class=\"o\">.</span><span class=\"n\">split</span><span class=\"p\">(</span><span
  class=\"s2\">\":\"</span><span class=\"p\">)[</span><span class=\"mi\">0</span><span
  class=\"p\">]</span>\n                        <span class=\"n\">plugin_func</span>
  <span class=\"o\">=</span> <span class=\"n\">plugin</span><span class=\"p\">[</span><span
  class=\"s2\">\"plugin\"</span><span class=\"p\">]</span><span class=\"o\">.</span><span
  class=\"n\">split</span><span class=\"p\">(</span><span class=\"s2\">\":\"</span><span
  class=\"p\">)[</span><span class=\"mi\">1</span><span class=\"p\">]</span>\n                        <span
  class=\"n\">plugin</span><span class=\"p\">[</span><span class=\"s2\">\"plugin\"</span><span
  class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"nb\">getattr</span><span
  class=\"p\">(</span>\n                            <span class=\"n\">importlib</span><span
  class=\"o\">.</span><span class=\"n\">import_module</span><span class=\"p\">(</span><span
  class=\"n\">plugin_module</span><span class=\"p\">),</span>\n                            <span
  class=\"n\">plugin_func</span><span class=\"p\">,</span>\n                        <span
  class=\"p\">)</span>\n                    <span class=\"n\">plugin</span><span class=\"p\">[</span><span
  class=\"s2\">\"config\"</span><span class=\"p\">]</span> <span class=\"o\">=</span>
  <span class=\"n\">plugin</span><span class=\"o\">.</span><span class=\"n\">get</span><span
  class=\"p\">(</span><span class=\"s2\">\"config\"</span><span class=\"p\">,</span>
  <span class=\"p\">{})</span>\n                    <span class=\"k\">for</span> <span
  class=\"n\">k</span><span class=\"p\">,</span> <span class=\"n\">_v</span> <span
  class=\"ow\">in</span> <span class=\"n\">plugin</span><span class=\"p\">[</span><span
  class=\"s2\">\"config\"</span><span class=\"p\">]</span><span class=\"o\">.</span><span
  class=\"n\">items</span><span class=\"p\">():</span>\n                        <span
  class=\"k\">if</span> <span class=\"n\">k</span> <span class=\"o\">==</span> <span
  class=\"s2\">\"markata\"</span><span class=\"p\">:</span>\n                            <span
  class=\"n\">plugin</span><span class=\"p\">[</span><span class=\"s2\">\"config\"</span><span
  class=\"p\">][</span><span class=\"n\">k</span><span class=\"p\">]</span> <span
  class=\"o\">=</span> <span class=\"n\">markata</span>\n\n                    <span
  class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">md</span>
  <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">md</span><span class=\"o\">.</span><span class=\"n\">use</span><span
  class=\"p\">(</span><span class=\"n\">plugin</span><span class=\"p\">[</span><span
  class=\"s2\">\"plugin\"</span><span class=\"p\">],</span> <span class=\"o\">**</span><span
  class=\"n\">plugin</span><span class=\"p\">[</span><span class=\"s2\">\"config\"</span><span
  class=\"p\">])</span>\n\n                <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">md</span><span class=\"o\">.</span><span class=\"n\">convert</span>
  <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">md</span><span class=\"o\">.</span><span class=\"n\">render</span>\n
  \               <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">md</span><span class=\"o\">.</span><span class=\"n\">toc</span> <span
  class=\"o\">=</span> <span class=\"s2\">\"\"</span>\n            <span class=\"k\">elif</span>
  <span class=\"p\">(</span>\n                <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
  class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">\"markdown_backend\"</span><span
  class=\"p\">,</span> <span class=\"s2\">\"\"</span><span class=\"p\">)</span>\n
  \               <span class=\"o\">.</span><span class=\"n\">lower</span><span class=\"p\">()</span>\n
  \               <span class=\"o\">.</span><span class=\"n\">replace</span><span
  class=\"p\">(</span><span class=\"s2\">\" \"</span><span class=\"p\">,</span> <span
  class=\"s2\">\"-\"</span><span class=\"p\">)</span>\n                <span class=\"o\">.</span><span
  class=\"n\">replace</span><span class=\"p\">(</span><span class=\"s2\">\"_\"</span><span
  class=\"p\">,</span> <span class=\"s2\">\"-\"</span><span class=\"p\">)</span>\n
  \               <span class=\"o\">==</span> <span class=\"s2\">\"markdown2\"</span>\n
  \           <span class=\"p\">):</span>\n                <span class=\"kn\">import</span>
  <span class=\"nn\">markdown2</span>\n\n                <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">md</span> <span class=\"o\">=</span> <span
  class=\"n\">markdown2</span><span class=\"o\">.</span><span class=\"n\">Markdown</span><span
  class=\"p\">(</span>\n                    <span class=\"n\">extras</span><span class=\"o\">=</span><span
  class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">render_markdown</span><span class=\"o\">.</span><span
  class=\"n\">extensions</span>\n                <span class=\"p\">)</span>\n                <span
  class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">md</span><span
  class=\"o\">.</span><span class=\"n\">toc</span> <span class=\"o\">=</span> <span
  class=\"s2\">\"\"</span>\n            <span class=\"k\">else</span><span class=\"p\">:</span>\n
  \               <span class=\"kn\">import</span> <span class=\"nn\">markdown</span>\n\n
  \               <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">md</span> <span class=\"o\">=</span> <span class=\"n\">markdown</span><span
  class=\"o\">.</span><span class=\"n\">Markdown</span><span class=\"p\">(</span>\n
  \                   <span class=\"n\">extensions</span><span class=\"o\">=</span><span
  class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">render_markdown</span><span class=\"o\">.</span><span
  class=\"n\">extensions</span>\n                <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"render\" style=\"margin:0;padding:.5rem
  1rem;\">render <em class=\"small\">function</em></h2>\n<div class=\"admonition source
  is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"render <em class=\"small\">source</em>\"</p>\n</div>\n<pre
  class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">render</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \           <span class=\"n\">config</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
  class=\"n\">render_markdown</span>\n            <span class=\"k\">with</span> <span
  class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">cache</span>
  <span class=\"k\">as</span> <span class=\"n\">cache</span><span class=\"p\">:</span>\n
  \               <span class=\"k\">for</span> <span class=\"n\">article</span> <span
  class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">iter_articles</span><span class=\"p\">(</span><span class=\"s2\">\"rendering
  markdown\"</span><span class=\"p\">):</span>\n                    <span class=\"n\">key</span>
  <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">make_hash</span><span class=\"p\">(</span>\n                        <span
  class=\"s2\">\"render_markdown\"</span><span class=\"p\">,</span>\n                        <span
  class=\"s2\">\"render\"</span><span class=\"p\">,</span>\n                        <span
  class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">content</span><span
  class=\"p\">,</span>\n                    <span class=\"p\">)</span>\n                    <span
  class=\"n\">html_from_cache</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">precache</span><span class=\"o\">.</span><span
  class=\"n\">get</span><span class=\"p\">(</span><span class=\"n\">key</span><span
  class=\"p\">)</span>\n                    <span class=\"k\">if</span> <span class=\"n\">html_from_cache</span>
  <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \                       <span class=\"n\">html</span> <span class=\"o\">=</span>
  <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">md</span><span
  class=\"o\">.</span><span class=\"n\">convert</span><span class=\"p\">(</span><span
  class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">content</span><span
  class=\"p\">)</span>\n                        <span class=\"n\">cache</span><span
  class=\"o\">.</span><span class=\"n\">add</span><span class=\"p\">(</span><span
  class=\"n\">key</span><span class=\"p\">,</span> <span class=\"n\">html</span><span
  class=\"p\">,</span> <span class=\"n\">expire</span><span class=\"o\">=</span><span
  class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">cache_expire</span><span
  class=\"p\">)</span>\n                    <span class=\"k\">else</span><span class=\"p\">:</span>\n
  \                       <span class=\"n\">html</span> <span class=\"o\">=</span>
  <span class=\"n\">html_from_cache</span>\n                    <span class=\"n\">article</span><span
  class=\"o\">.</span><span class=\"n\">html</span> <span class=\"o\">=</span> <span
  class=\"n\">html</span>\n                    <span class=\"n\">article</span><span
  class=\"o\">.</span><span class=\"n\">article_html</span> <span class=\"o\">=</span>
  <span class=\"n\">copy</span><span class=\"o\">.</span><span class=\"n\">deepcopy</span><span
  class=\"p\">(</span><span class=\"n\">html</span><span class=\"p\">)</span>\n\n
  \                   <span class=\"n\">article</span><span class=\"o\">.</span><span
  class=\"n\">html</span> <span class=\"o\">=</span> <span class=\"n\">html</span>\n
  \                   <span class=\"n\">article</span><span class=\"o\">.</span><span
  class=\"n\">article_html</span> <span class=\"o\">=</span> <span class=\"n\">article</span><span
  class=\"o\">.</span><span class=\"n\">article_html</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"convert_to_list\" style=\"margin:0;padding:.5rem
  1rem;\">convert_to_list <em class=\"small\">method</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"convert_to_list
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">convert_to_list</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
  class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">):</span>\n                <span
  class=\"k\">if</span> <span class=\"ow\">not</span> <span class=\"nb\">isinstance</span><span
  class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span class=\"nb\">list</span><span
  class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
  class=\"p\">[</span><span class=\"n\">v</span><span class=\"p\">]</span>\n                <span
  class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n</main>\n<footer>\xA9
  2024</footer>\n</body></html>"
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
            html: Optional[str] = None
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
```


!! method <h2 id='convert_to_list' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>convert_to_list <em class='small'>method</em></h2>

???+ source "convert_to_list <em class='small'>source</em>"

```python

        def convert_to_list(cls, v):
                if not isinstance(v, list):
                    return [v]
                return v
```
