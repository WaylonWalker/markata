---
content: "## Run it\n\n``` bash\npython -m markata.cli.summary\n```\n\n## Configuration\n\nThere
  are two main things currently supported by summary, It can count the\nnumber of
  posts based on a filter (`filter_count'), and it can automatically\nlist all the
  values of an attribute and the number of posts that have that\nattribute (`grid_attr`).\n\n###
  grid_attr\n\n`grid_attr` will map over all posts, find all values for each attribute\nconfigured,
  then report the number of posts for each value.\n\n``` toml\n[markata.summary]\ngrid_attr
  = ['tags', 'series']\n```\n\nExample output that will be shown in the summary, counting
  all posts with each\ntag value.\n\n```\nTAGS\n247  python\n90   linux\n68   cli\n49
  \  kedro\n46   bash\n```\n\n### filter_counts\n\n`filter_count` will pass a given
  filter into `markata.map` and return the number\nof posts.\n\n```\n[[markata.summary.filter_count]]\nname='drafts'\nfilter=\"published
  == 'False'\"\ncolor='red'\n\n[[markata.summary.filter_count]]\nname='articles'\ncolor='dark_orange'\n\n[[markata.summary.filter_count]]\nname='py_modules'\nfilter='\"plugin\"
  not in slug and \"docs\" not in str(path)'\ncolor=\"yellow1\"\n\n[markata.summary.filter_count.published]\nfilter=\"published
  == 'True'\"\ncolor='green1'\n\n[markata.summary.filter_count.plugins]\nfilter='\"plugin\"
  in slug and \"docs\" not in str(path)'\ncolor=\"blue\"\n\n[markata.summary.filter_count.docs]\nfilter=\"'docs'
  in str(path)\"\ncolor='purple'\n```\n\nExample output might look like this, showing
  the number of posts that contained\nwithin each filter specified.\n\n```\n8 drafts\n66
  articles\n20 py_modules\n58 published\n38 plugins\n8 docs\n```\n\n\n!! class <h2
  id='FilterCount' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>FilterCount
  <em class='small'>class</em></h2>\n\n???+ source \"FilterCount <em class='small'>source</em>\"\n\n```python\n\n
  \       class FilterCount(pydantic.BaseModel):\n            name: str\n            filter:
  str = \"True\"\n            color: str = \"white\"\n```\n\n\n!! class <h2 id='SummaryConfig'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>SummaryConfig <em
  class='small'>class</em></h2>\n\n???+ source \"SummaryConfig <em class='small'>source</em>\"\n\n```python\n\n
  \       class SummaryConfig(pydantic.BaseModel):\n            grid_attr: List[str]
  = [\"tags\", \"series\"]\n            filter_count: List[FilterCount] = FilterCount(\n
  \               name=\"drafts\", filter=\"published == 'False'\", color=\"red\"\n
  \           )\n```\n\n\n!! class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>Config <em class='small'>class</em></h2>\n\n???+ source \"Config <em class='small'>source</em>\"\n\n```python\n\n
  \       class Config(pydantic.BaseModel):\n            summary: SummaryConfig =
  SummaryConfig()\n```\n\n\n!! function <h2 id='config_model' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>config_model <em class='small'>function</em></h2>\n\n???+
  source \"config_model <em class='small'>source</em>\"\n\n```python\n\n        def
  config_model(markata: \"Markata\") -> None:\n            markata.config_models.append(Config)\n```\n\n\n!!
  class <h2 id='Summary' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Summary
  <em class='small'>class</em></h2>\n\n???+ source \"Summary <em class='small'>source</em>\"\n\n```python\n\n
  \       class Summary:\n            def __init__(self, m: \"Markata\", simple: bool
  = False) -> None:\n                self.m = m\n                self.simple = simple\n\n
  \           def get_grid(self) -> None:\n                \"create a rich grid to
  display the summary\"\n                self.grid = Table.grid(expand=True)\n\n                for
  filter_count in self.m.config.summary.filter_count:\n                    self.filter_count(filter_count)\n\n
  \               for attr in self.m.config.summary.grid_attr:\n                    self.grid_attr(attr)\n\n
  \               return self.grid\n\n            def filter_count(\n                self,\n
  \               fc: FilterCount,\n            ) -> None:\n                \"add
  a row in the grid for the number of items in a filter config\"\n                self.grid.add_row(\n
  \                   f\"[{fc.color}]{len(self.m.map(filter=fc.filter))}[/] {fc.name}\"\n
  \               )\n\n            def grid_attr(self, attr: str) -> None:\n                \"add
  attribute the the object grid\"\n                posts = list(\n                    flatten(\n
  \                       [\n                            tags if isinstance(tags,
  list) else [tags]\n                            for a in self.m.posts\n                            if
  (tags := a.get(attr, None)) is not None\n                        ],\n                    ),\n
  \               )\n                if len(posts) > 0:\n                    self.grid.add_row()\n
  \                   self.grid.add_row(f\"[bold gold1]{attr.upper()}[/]\")\n                    for
  post, count in Counter(posts).most_common():\n                        self.grid.add_row(f'{count}
  {\" \"*(3-len(str(count)))} {post}')\n\n            def __rich__(self) -> Union[Panel,
  Table]:\n                grid = self.get_grid()\n\n                if self.simple:\n
  \                   return grid\n                else:\n                    return
  Panel(\n                        grid,\n                        title=\"[gold1]summary[/]\",\n
  \                       border_style=\"magenta\",\n                        expand=False,\n
  \                   )\n```\n\n\n!! function <h2 id='configure' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>configure <em class='small'>function</em></h2>\n\n???+
  source \"configure <em class='small'>source</em>\"\n\n```python\n\n        def configure(markata:
  \"Markata\") -> None:\n            def get_summary(self):\n                try:\n
  \                   return self._summary\n                except AttributeError:\n
  \                   self._summary: Summary = Summary(self)\n                    return
  self._summary\n\n            from markata import Markata\n\n            Markata.summary
  = property(get_summary)\n```\n\n\n!! function <h2 id='cli' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>cli <em class='small'>function</em></h2>\n
  \   Markata hook to implement base cli commands.\n???+ source \"cli <em class='small'>source</em>\"\n\n```python\n\n
  \       def cli(app: typer.Typer, markata: \"Markata\") -> None:\n            \"\"\"\n
  \           Markata hook to implement base cli commands.\n            \"\"\"\n            summary_app
  = typer.Typer()\n            app.add_typer(summary_app, name=\"summary\")\n\n            @summary_app.callback(invoke_without_command=True)\n
  \           def summary():\n                \"show the application summary\"\n                from
  rich import print\n\n                markata.console.quiet = True\n\n                print(Summary(markata,
  simple=True))\n```\n\n\n!! method <h2 id='__init__' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>__init__ <em class='small'>method</em></h2>\n\n???+ source \"__init__ <em
  class='small'>source</em>\"\n\n```python\n\n        def __init__(self, m: \"Markata\",
  simple: bool = False) -> None:\n                self.m = m\n                self.simple
  = simple\n```\n\n\n!! method <h2 id='get_grid' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>get_grid <em class='small'>method</em></h2>\n    create a rich grid to display
  the summary\n???+ source \"get_grid <em class='small'>source</em>\"\n\n```python\n\n
  \       def get_grid(self) -> None:\n                \"create a rich grid to display
  the summary\"\n                self.grid = Table.grid(expand=True)\n\n                for
  filter_count in self.m.config.summary.filter_count:\n                    self.filter_count(filter_count)\n\n
  \               for attr in self.m.config.summary.grid_attr:\n                    self.grid_attr(attr)\n\n
  \               return self.grid\n```\n\n\n!! method <h2 id='filter_count' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>filter_count <em class='small'>method</em></h2>\n
  \   add a row in the grid for the number of items in a filter config\n???+ source
  \"filter_count <em class='small'>source</em>\"\n\n```python\n\n        def filter_count(\n
  \               self,\n                fc: FilterCount,\n            ) -> None:\n
  \               \"add a row in the grid for the number of items in a filter config\"\n
  \               self.grid.add_row(\n                    f\"[{fc.color}]{len(self.m.map(filter=fc.filter))}[/]
  {fc.name}\"\n                )\n```\n\n\n!! method <h2 id='grid_attr' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>grid_attr <em class='small'>method</em></h2>\n
  \   add attribute the the object grid\n???+ source \"grid_attr <em class='small'>source</em>\"\n\n```python\n\n
  \       def grid_attr(self, attr: str) -> None:\n                \"add attribute
  the the object grid\"\n                posts = list(\n                    flatten(\n
  \                       [\n                            tags if isinstance(tags,
  list) else [tags]\n                            for a in self.m.posts\n                            if
  (tags := a.get(attr, None)) is not None\n                        ],\n                    ),\n
  \               )\n                if len(posts) > 0:\n                    self.grid.add_row()\n
  \                   self.grid.add_row(f\"[bold gold1]{attr.upper()}[/]\")\n                    for
  post, count in Counter(posts).most_common():\n                        self.grid.add_row(f'{count}
  {\" \"*(3-len(str(count)))} {post}')\n```\n\n\n!! method <h2 id='__rich__' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>__rich__ <em class='small'>method</em></h2>\n\n???+
  source \"__rich__ <em class='small'>source</em>\"\n\n```python\n\n        def __rich__(self)
  -> Union[Panel, Table]:\n                grid = self.get_grid()\n\n                if
  self.simple:\n                    return grid\n                else:\n                    return
  Panel(\n                        grid,\n                        title=\"[gold1]summary[/]\",\n
  \                       border_style=\"magenta\",\n                        expand=False,\n
  \                   )\n```\n\n\n!! function <h2 id='get_summary' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>get_summary <em class='small'>function</em></h2>\n\n???+
  source \"get_summary <em class='small'>source</em>\"\n\n```python\n\n        def
  get_summary(self):\n                try:\n                    return self._summary\n
  \               except AttributeError:\n                    self._summary: Summary
  = Summary(self)\n                    return self._summary\n```\n\n\n!! function
  <h2 id='summary' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>summary
  <em class='small'>function</em></h2>\n    show the application summary\n???+ source
  \"summary <em class='small'>source</em>\"\n\n```python\n\n        def summary():\n
  \               \"show the application summary\"\n                from rich import
  print\n\n                markata.console.quiet = True\n\n                print(Summary(markata,
  simple=True))\n```\n"
date: 0001-01-01
description: There are two main things currently supported by summary, It can count
  the grid_attr Example output that will be shown in the summary, counting all posts
  with e
html: "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n<title>Summary.Py</title>\n<meta
  charset=\"utf-8\"/>\n<meta content=\"width=device-width, initial-scale=1\" name=\"viewport\"/>\n<meta
  content=\"There are two main things currently supported by summary, It can count
  the grid_attr Example output that will be shown in the summary, counting all posts
  with e\" name=\"description\"/>\n<link href=\"/static/favicon.ico\" rel=\"icon\"
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
  name=\"og:type\" property=\"og:type\"/><meta content=\"There are two main things
  currently supported by summary, It can count the grid_attr Example output that will
  be shown in the summary, counting all posts with e\" name=\"description\" property=\"description\"/><meta
  content=\"There are two main things currently supported by summary, It can count
  the grid_attr Example output that will be shown in the summary, counting all posts
  with e\" name=\"og:description\" property=\"og:description\"/><meta content=\"There
  are two main things currently supported by summary, It can count the grid_attr Example
  output that will be shown in the summary, counting all posts with e\" name=\"twitter:description\"
  property=\"twitter:description\"/><meta content=\"Summary.Py | Markata\" name=\"og:title\"
  property=\"og:title\"/><meta content=\"Summary.Py | Markata\" name=\"twitter:title\"
  property=\"twitter:title\"/><meta content=\"https://markata.dev//markata/cli/summary-og.png\"
  name=\"og:image\" property=\"og:image\"/><meta content=\"https://markata.dev//markata/cli/summary-og.png\"
  name=\"twitter:image\" property=\"twitter:image\"/><meta content=\"1600\" name=\"og:image:width\"
  property=\"og:image:width\"/><meta content=\"900\" name=\"og:image:width\" property=\"og:image:width\"/><meta
  content=\"summary_large_image\" name=\"twitter:card\" property=\"twitter:card\"/><meta
  content=\"Markata\" name=\"og:site_name\" property=\"og:site_name\"/><meta content=\"@_waylonwalker\"
  name=\"twitter:creator\" property=\"twitter:creator\"/><meta content=\"Summary.Py\"
  name=\"title\" property=\"title\"/><meta content=\"markata 0.8.0.dev6\" name=\"generator\"
  property=\"generator\"/><link href=\"https://markata.dev//markata/cli/summary/\"
  rel=\"canonical\"/><meta content=\"https://markata.dev//markata/cli/summary/\" name=\"og:url\"
  property=\"og:url\"/></head>\n<body><nav>\n<a href=\"/\">markata</a>\n<a href=\"https://github.com/WaylonWalker/markata\">GitHub</a>\n</nav>\n<div>\n<label
  class=\"theme-switch\" for=\"checkbox-theme\" id=\"theme-switch\" title=\"light/dark
  mode toggle\">\n<input id=\"checkbox-theme\" type=\"checkbox\"/>\n<div class=\"slider
  round\"></div>\n</label>\n</div>\n<section class=\"title\">\n<h1 id=\"title\">\n
  \           Summary.Py \n            \n        </h1>\n</section>\n<main><h2 id=\"run-it\">Run
  it <a class=\"header-anchor\" href=\"#run-it\"><svg aria-hidden=\"true\" class=\"heading-permalink\"
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
  3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632 1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<pre
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
  \       \n<div class=\"highlight\"><pre><span></span>python<span class=\"w\"> </span>-m<span
  class=\"w\"> </span>markata.cli.summary\n</pre></div>\n\n</pre>\n<h2 id=\"configuration\">Configuration
  <a class=\"header-anchor\" href=\"#configuration\"><svg aria-hidden=\"true\" class=\"heading-permalink\"
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
  3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632 1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>There
  are two main things currently supported by summary, It can count the\nnumber of
  posts based on a filter (<code>filter_count'), and it can automatically list all
  the values of an attribute and the number of posts that have that attribute (</code>grid_attr`).</p>\n<h3>grid_attr</h3>\n<p><code>grid_attr</code>
  will map over all posts, find all values for each attribute\nconfigured, then report
  the number of posts for each value.</p>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
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
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata.summary]</span>\n<span
  class=\"n\">grid_attr</span><span class=\"w\"> </span><span class=\"o\">=</span><span
  class=\"w\"> </span><span class=\"p\">[</span><span class=\"s1\">'tags'</span><span
  class=\"p\">,</span><span class=\"w\"> </span><span class=\"s1\">'series'</span><span
  class=\"p\">]</span>\n</pre></div>\n\n</pre>\n<p>Example output that will be shown
  in the summary, counting all posts with each\ntag value.</p>\n<pre class=\"wrapper\">\n\n<div
  class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
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
  \       \n<div class=\"highlight\"><pre><span></span>TAGS\n247  python\n90   linux\n68
  \  cli\n49   kedro\n46   bash\n</pre></div>\n\n</pre>\n<h3>filter_counts</h3>\n<p><code>filter_count</code>
  will pass a given filter into <code>markata.map</code> and return the number\nof
  posts.</p>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button class=\"copy\"
  onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
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
  \       \n<div class=\"highlight\"><pre><span></span>[[markata.summary.filter_count]]\nname='drafts'\nfilter=\"published
  == 'False'\"\ncolor='red'\n\n[[markata.summary.filter_count]]\nname='articles'\ncolor='dark_orange'\n\n[[markata.summary.filter_count]]\nname='py_modules'\nfilter='\"plugin\"
  not in slug and \"docs\" not in str(path)'\ncolor=\"yellow1\"\n\n[markata.summary.filter_count.published]\nfilter=\"published
  == 'True'\"\ncolor='green1'\n\n[markata.summary.filter_count.plugins]\nfilter='\"plugin\"
  in slug and \"docs\" not in str(path)'\ncolor=\"blue\"\n\n[markata.summary.filter_count.docs]\nfilter=\"'docs'
  in str(path)\"\ncolor='purple'\n</pre></div>\n\n</pre>\n<p>Example output might
  look like this, showing the number of posts that contained\nwithin each filter specified.</p>\n<pre
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
  \       \n<div class=\"highlight\"><pre><span></span>8 drafts\n66 articles\n20 py_modules\n58
  published\n38 plugins\n8 docs\n</pre></div>\n\n</pre>\n<p>!! class </p><h2 class=\"admonition-title\"
  id=\"FilterCount\" style=\"margin:0;padding:.5rem 1rem;\">FilterCount <em class=\"small\">class</em></h2>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"FilterCount
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
  <span class=\"nc\">FilterCount</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
  class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
  \           <span class=\"n\">name</span><span class=\"p\">:</span> <span class=\"nb\">str</span>\n
  \           <span class=\"nb\">filter</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
  <span class=\"o\">=</span> <span class=\"s2\">\"True\"</span>\n            <span
  class=\"n\">color</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
  <span class=\"o\">=</span> <span class=\"s2\">\"white\"</span>\n</pre></div>\n\n</pre>\n<p>!!
  class </p><h2 class=\"admonition-title\" id=\"SummaryConfig\" style=\"margin:0;padding:.5rem
  1rem;\">SummaryConfig <em class=\"small\">class</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"SummaryConfig
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
  <span class=\"nc\">SummaryConfig</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
  class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
  \           <span class=\"n\">grid_attr</span><span class=\"p\">:</span> <span class=\"n\">List</span><span
  class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">]</span> <span
  class=\"o\">=</span> <span class=\"p\">[</span><span class=\"s2\">\"tags\"</span><span
  class=\"p\">,</span> <span class=\"s2\">\"series\"</span><span class=\"p\">]</span>\n
  \           <span class=\"n\">filter_count</span><span class=\"p\">:</span> <span
  class=\"n\">List</span><span class=\"p\">[</span><span class=\"n\">FilterCount</span><span
  class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"n\">FilterCount</span><span
  class=\"p\">(</span>\n                <span class=\"n\">name</span><span class=\"o\">=</span><span
  class=\"s2\">\"drafts\"</span><span class=\"p\">,</span> <span class=\"nb\">filter</span><span
  class=\"o\">=</span><span class=\"s2\">\"published == 'False'\"</span><span class=\"p\">,</span>
  <span class=\"n\">color</span><span class=\"o\">=</span><span class=\"s2\">\"red\"</span>\n
  \           <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!! class </p><h2
  class=\"admonition-title\" id=\"Config\" style=\"margin:0;padding:.5rem 1rem;\">Config
  <em class=\"small\">class</em></h2>\n<div class=\"admonition source is-collapsible
  collapsible-open\">\n<p class=\"admonition-title\">\"Config <em class=\"small\">source</em>\"</p>\n</div>\n<pre
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
  <span class=\"nc\">Config</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
  class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
  \           <span class=\"n\">summary</span><span class=\"p\">:</span> <span class=\"n\">SummaryConfig</span>
  <span class=\"o\">=</span> <span class=\"n\">SummaryConfig</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n<p>!!
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
  class </p><h2 class=\"admonition-title\" id=\"Summary\" style=\"margin:0;padding:.5rem
  1rem;\">Summary <em class=\"small\">class</em></h2>\n<div class=\"admonition source
  is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"Summary <em
  class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
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
  <span class=\"nc\">Summary</span><span class=\"p\">:</span>\n            <span class=\"k\">def</span>
  <span class=\"fm\">__init__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">,</span> <span class=\"n\">m</span><span class=\"p\">:</span> <span
  class=\"s2\">\"Markata\"</span><span class=\"p\">,</span> <span class=\"n\">simple</span><span
  class=\"p\">:</span> <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span
  class=\"kc\">False</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">m</span> <span
  class=\"o\">=</span> <span class=\"n\">m</span>\n                <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">simple</span> <span class=\"o\">=</span> <span
  class=\"n\">simple</span>\n\n            <span class=\"k\">def</span> <span class=\"nf\">get_grid</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">)</span> <span
  class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \               <span class=\"s2\">\"create a rich grid to display the summary\"</span>\n
  \               <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">grid</span>
  <span class=\"o\">=</span> <span class=\"n\">Table</span><span class=\"o\">.</span><span
  class=\"n\">grid</span><span class=\"p\">(</span><span class=\"n\">expand</span><span
  class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n\n
  \               <span class=\"k\">for</span> <span class=\"n\">filter_count</span>
  <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">m</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">summary</span><span class=\"o\">.</span><span
  class=\"n\">filter_count</span><span class=\"p\">:</span>\n                    <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">filter_count</span><span
  class=\"p\">(</span><span class=\"n\">filter_count</span><span class=\"p\">)</span>\n\n
  \               <span class=\"k\">for</span> <span class=\"n\">attr</span> <span
  class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">m</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">summary</span><span class=\"o\">.</span><span
  class=\"n\">grid_attr</span><span class=\"p\">:</span>\n                    <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">grid_attr</span><span
  class=\"p\">(</span><span class=\"n\">attr</span><span class=\"p\">)</span>\n\n
  \               <span class=\"k\">return</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">grid</span>\n\n            <span class=\"k\">def</span>
  <span class=\"nf\">filter_count</span><span class=\"p\">(</span>\n                <span
  class=\"bp\">self</span><span class=\"p\">,</span>\n                <span class=\"n\">fc</span><span
  class=\"p\">:</span> <span class=\"n\">FilterCount</span><span class=\"p\">,</span>\n
  \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n                <span class=\"s2\">\"add a row in the grid
  for the number of items in a filter config\"</span>\n                <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">grid</span><span class=\"o\">.</span><span
  class=\"n\">add_row</span><span class=\"p\">(</span>\n                    <span
  class=\"sa\">f</span><span class=\"s2\">\"[</span><span class=\"si\">{</span><span
  class=\"n\">fc</span><span class=\"o\">.</span><span class=\"n\">color</span><span
  class=\"si\">}</span><span class=\"s2\">]</span><span class=\"si\">{</span><span
  class=\"nb\">len</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">m</span><span class=\"o\">.</span><span class=\"n\">map</span><span
  class=\"p\">(</span><span class=\"nb\">filter</span><span class=\"o\">=</span><span
  class=\"n\">fc</span><span class=\"o\">.</span><span class=\"n\">filter</span><span
  class=\"p\">))</span><span class=\"si\">}</span><span class=\"s2\">[/] </span><span
  class=\"si\">{</span><span class=\"n\">fc</span><span class=\"o\">.</span><span
  class=\"n\">name</span><span class=\"si\">}</span><span class=\"s2\">\"</span>\n
  \               <span class=\"p\">)</span>\n\n            <span class=\"k\">def</span>
  <span class=\"nf\">grid_attr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">,</span> <span class=\"n\">attr</span><span class=\"p\">:</span> <span
  class=\"nb\">str</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
  class=\"s2\">\"add attribute the the object grid\"</span>\n                <span
  class=\"n\">posts</span> <span class=\"o\">=</span> <span class=\"nb\">list</span><span
  class=\"p\">(</span>\n                    <span class=\"n\">flatten</span><span
  class=\"p\">(</span>\n                        <span class=\"p\">[</span>\n                            <span
  class=\"n\">tags</span> <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
  class=\"p\">(</span><span class=\"n\">tags</span><span class=\"p\">,</span> <span
  class=\"nb\">list</span><span class=\"p\">)</span> <span class=\"k\">else</span>
  <span class=\"p\">[</span><span class=\"n\">tags</span><span class=\"p\">]</span>\n
  \                           <span class=\"k\">for</span> <span class=\"n\">a</span>
  <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">m</span><span class=\"o\">.</span><span class=\"n\">posts</span>\n                            <span
  class=\"k\">if</span> <span class=\"p\">(</span><span class=\"n\">tags</span> <span
  class=\"o\">:=</span> <span class=\"n\">a</span><span class=\"o\">.</span><span
  class=\"n\">get</span><span class=\"p\">(</span><span class=\"n\">attr</span><span
  class=\"p\">,</span> <span class=\"kc\">None</span><span class=\"p\">))</span> <span
  class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span>\n
  \                       <span class=\"p\">],</span>\n                    <span class=\"p\">),</span>\n
  \               <span class=\"p\">)</span>\n                <span class=\"k\">if</span>
  <span class=\"nb\">len</span><span class=\"p\">(</span><span class=\"n\">posts</span><span
  class=\"p\">)</span> <span class=\"o\">&gt;</span> <span class=\"mi\">0</span><span
  class=\"p\">:</span>\n                    <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">grid</span><span class=\"o\">.</span><span class=\"n\">add_row</span><span
  class=\"p\">()</span>\n                    <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">grid</span><span class=\"o\">.</span><span class=\"n\">add_row</span><span
  class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">\"[bold gold1]</span><span
  class=\"si\">{</span><span class=\"n\">attr</span><span class=\"o\">.</span><span
  class=\"n\">upper</span><span class=\"p\">()</span><span class=\"si\">}</span><span
  class=\"s2\">[/]\"</span><span class=\"p\">)</span>\n                    <span class=\"k\">for</span>
  <span class=\"n\">post</span><span class=\"p\">,</span> <span class=\"n\">count</span>
  <span class=\"ow\">in</span> <span class=\"n\">Counter</span><span class=\"p\">(</span><span
  class=\"n\">posts</span><span class=\"p\">)</span><span class=\"o\">.</span><span
  class=\"n\">most_common</span><span class=\"p\">():</span>\n                        <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">grid</span><span
  class=\"o\">.</span><span class=\"n\">add_row</span><span class=\"p\">(</span><span
  class=\"sa\">f</span><span class=\"s1\">'</span><span class=\"si\">{</span><span
  class=\"n\">count</span><span class=\"si\">}</span><span class=\"s1\"> </span><span
  class=\"si\">{</span><span class=\"s2\">\" \"</span><span class=\"o\">*</span><span
  class=\"p\">(</span><span class=\"mi\">3</span><span class=\"o\">-</span><span class=\"nb\">len</span><span
  class=\"p\">(</span><span class=\"nb\">str</span><span class=\"p\">(</span><span
  class=\"n\">count</span><span class=\"p\">)))</span><span class=\"si\">}</span><span
  class=\"s1\"> </span><span class=\"si\">{</span><span class=\"n\">post</span><span
  class=\"si\">}</span><span class=\"s1\">'</span><span class=\"p\">)</span>\n\n            <span
  class=\"k\">def</span> <span class=\"nf\">__rich__</span><span class=\"p\">(</span><span
  class=\"bp\">self</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"n\">Union</span><span class=\"p\">[</span><span class=\"n\">Panel</span><span
  class=\"p\">,</span> <span class=\"n\">Table</span><span class=\"p\">]:</span>\n
  \               <span class=\"n\">grid</span> <span class=\"o\">=</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">get_grid</span><span class=\"p\">()</span>\n\n
  \               <span class=\"k\">if</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">simple</span><span class=\"p\">:</span>\n
  \                   <span class=\"k\">return</span> <span class=\"n\">grid</span>\n
  \               <span class=\"k\">else</span><span class=\"p\">:</span>\n                    <span
  class=\"k\">return</span> <span class=\"n\">Panel</span><span class=\"p\">(</span>\n
  \                       <span class=\"n\">grid</span><span class=\"p\">,</span>\n
  \                       <span class=\"n\">title</span><span class=\"o\">=</span><span
  class=\"s2\">\"[gold1]summary[/]\"</span><span class=\"p\">,</span>\n                        <span
  class=\"n\">border_style</span><span class=\"o\">=</span><span class=\"s2\">\"magenta\"</span><span
  class=\"p\">,</span>\n                        <span class=\"n\">expand</span><span
  class=\"o\">=</span><span class=\"kc\">False</span><span class=\"p\">,</span>\n
  \                   <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!! function
  </p><h2 class=\"admonition-title\" id=\"configure\" style=\"margin:0;padding:.5rem
  1rem;\">configure <em class=\"small\">function</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"configure
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
  \           <span class=\"k\">def</span> <span class=\"nf\">get_summary</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">):</span>\n
  \               <span class=\"k\">try</span><span class=\"p\">:</span>\n                    <span
  class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">_summary</span>\n                <span class=\"k\">except</span> <span
  class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n                    <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_summary</span><span
  class=\"p\">:</span> <span class=\"n\">Summary</span> <span class=\"o\">=</span>
  <span class=\"n\">Summary</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">)</span>\n                    <span class=\"k\">return</span> <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_summary</span>\n\n
  \           <span class=\"kn\">from</span> <span class=\"nn\">markata</span> <span
  class=\"kn\">import</span> <span class=\"n\">Markata</span>\n\n            <span
  class=\"n\">Markata</span><span class=\"o\">.</span><span class=\"n\">summary</span>
  <span class=\"o\">=</span> <span class=\"nb\">property</span><span class=\"p\">(</span><span
  class=\"n\">get_summary</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"cli\" style=\"margin:0;padding:.5rem
  1rem;\">cli <em class=\"small\">function</em></h2>\nMarkata hook to implement base
  cli commands.\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
  class=\"admonition-title\">\"cli <em class=\"small\">source</em>\"</p>\n</div>\n<pre
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
  <span class=\"nf\">cli</span><span class=\"p\">(</span><span class=\"n\">app</span><span
  class=\"p\">:</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
  class=\"n\">Typer</span><span class=\"p\">,</span> <span class=\"n\">markata</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
  class=\"w\">            </span><span class=\"sd\">\"\"\"</span>\n<span class=\"sd\">
  \           Markata hook to implement base cli commands.</span>\n<span class=\"sd\">
  \           \"\"\"</span>\n            <span class=\"n\">summary_app</span> <span
  class=\"o\">=</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
  class=\"n\">Typer</span><span class=\"p\">()</span>\n            <span class=\"n\">app</span><span
  class=\"o\">.</span><span class=\"n\">add_typer</span><span class=\"p\">(</span><span
  class=\"n\">summary_app</span><span class=\"p\">,</span> <span class=\"n\">name</span><span
  class=\"o\">=</span><span class=\"s2\">\"summary\"</span><span class=\"p\">)</span>\n\n
  \           <span class=\"nd\">@summary_app</span><span class=\"o\">.</span><span
  class=\"n\">callback</span><span class=\"p\">(</span><span class=\"n\">invoke_without_command</span><span
  class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n            <span
  class=\"k\">def</span> <span class=\"nf\">summary</span><span class=\"p\">():</span>\n
  \               <span class=\"s2\">\"show the application summary\"</span>\n                <span
  class=\"kn\">from</span> <span class=\"nn\">rich</span> <span class=\"kn\">import</span>
  <span class=\"nb\">print</span>\n\n                <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
  class=\"n\">quiet</span> <span class=\"o\">=</span> <span class=\"kc\">True</span>\n\n
  \               <span class=\"nb\">print</span><span class=\"p\">(</span><span class=\"n\">Summary</span><span
  class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">,</span> <span
  class=\"n\">simple</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
  class=\"p\">))</span>\n</pre></div>\n\n</pre>\n<p>!! method </p><h2 class=\"admonition-title\"
  id=\"__init__\" style=\"margin:0;padding:.5rem 1rem;\"><strong>init</strong> <em
  class=\"small\">method</em></h2>\n<div class=\"admonition source is-collapsible
  collapsible-open\">\n<p class=\"admonition-title\">\"<strong>init</strong> <em class=\"small\">source</em>\"</p>\n</div>\n<pre
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
  <span class=\"fm\">__init__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">,</span> <span class=\"n\">m</span><span class=\"p\">:</span> <span
  class=\"s2\">\"Markata\"</span><span class=\"p\">,</span> <span class=\"n\">simple</span><span
  class=\"p\">:</span> <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span
  class=\"kc\">False</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">m</span> <span
  class=\"o\">=</span> <span class=\"n\">m</span>\n                <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">simple</span> <span class=\"o\">=</span> <span
  class=\"n\">simple</span>\n</pre></div>\n\n</pre>\n<p>!! method </p><h2 class=\"admonition-title\"
  id=\"get_grid\" style=\"margin:0;padding:.5rem 1rem;\">get_grid <em class=\"small\">method</em></h2>\ncreate
  a rich grid to display the summary\n<div class=\"admonition source is-collapsible
  collapsible-open\">\n<p class=\"admonition-title\">\"get_grid <em class=\"small\">source</em>\"</p>\n</div>\n<pre
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
  <span class=\"nf\">get_grid</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n                <span class=\"s2\">\"create a rich grid to
  display the summary\"</span>\n                <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">grid</span> <span class=\"o\">=</span> <span
  class=\"n\">Table</span><span class=\"o\">.</span><span class=\"n\">grid</span><span
  class=\"p\">(</span><span class=\"n\">expand</span><span class=\"o\">=</span><span
  class=\"kc\">True</span><span class=\"p\">)</span>\n\n                <span class=\"k\">for</span>
  <span class=\"n\">filter_count</span> <span class=\"ow\">in</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">m</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">summary</span><span class=\"o\">.</span><span
  class=\"n\">filter_count</span><span class=\"p\">:</span>\n                    <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">filter_count</span><span
  class=\"p\">(</span><span class=\"n\">filter_count</span><span class=\"p\">)</span>\n\n
  \               <span class=\"k\">for</span> <span class=\"n\">attr</span> <span
  class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">m</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">summary</span><span class=\"o\">.</span><span
  class=\"n\">grid_attr</span><span class=\"p\">:</span>\n                    <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">grid_attr</span><span
  class=\"p\">(</span><span class=\"n\">attr</span><span class=\"p\">)</span>\n\n
  \               <span class=\"k\">return</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">grid</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"filter_count\" style=\"margin:0;padding:.5rem
  1rem;\">filter_count <em class=\"small\">method</em></h2>\nadd a row in the grid
  for the number of items in a filter config\n<div class=\"admonition source is-collapsible
  collapsible-open\">\n<p class=\"admonition-title\">\"filter_count <em class=\"small\">source</em>\"</p>\n</div>\n<pre
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
  <span class=\"nf\">filter_count</span><span class=\"p\">(</span>\n                <span
  class=\"bp\">self</span><span class=\"p\">,</span>\n                <span class=\"n\">fc</span><span
  class=\"p\">:</span> <span class=\"n\">FilterCount</span><span class=\"p\">,</span>\n
  \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n                <span class=\"s2\">\"add a row in the grid
  for the number of items in a filter config\"</span>\n                <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">grid</span><span class=\"o\">.</span><span
  class=\"n\">add_row</span><span class=\"p\">(</span>\n                    <span
  class=\"sa\">f</span><span class=\"s2\">\"[</span><span class=\"si\">{</span><span
  class=\"n\">fc</span><span class=\"o\">.</span><span class=\"n\">color</span><span
  class=\"si\">}</span><span class=\"s2\">]</span><span class=\"si\">{</span><span
  class=\"nb\">len</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">m</span><span class=\"o\">.</span><span class=\"n\">map</span><span
  class=\"p\">(</span><span class=\"nb\">filter</span><span class=\"o\">=</span><span
  class=\"n\">fc</span><span class=\"o\">.</span><span class=\"n\">filter</span><span
  class=\"p\">))</span><span class=\"si\">}</span><span class=\"s2\">[/] </span><span
  class=\"si\">{</span><span class=\"n\">fc</span><span class=\"o\">.</span><span
  class=\"n\">name</span><span class=\"si\">}</span><span class=\"s2\">\"</span>\n
  \               <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!! method
  </p><h2 class=\"admonition-title\" id=\"grid_attr\" style=\"margin:0;padding:.5rem
  1rem;\">grid_attr <em class=\"small\">method</em></h2>\nadd attribute the the object
  grid\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"grid_attr
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
  <span class=\"nf\">grid_attr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">,</span> <span class=\"n\">attr</span><span class=\"p\">:</span> <span
  class=\"nb\">str</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
  class=\"s2\">\"add attribute the the object grid\"</span>\n                <span
  class=\"n\">posts</span> <span class=\"o\">=</span> <span class=\"nb\">list</span><span
  class=\"p\">(</span>\n                    <span class=\"n\">flatten</span><span
  class=\"p\">(</span>\n                        <span class=\"p\">[</span>\n                            <span
  class=\"n\">tags</span> <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
  class=\"p\">(</span><span class=\"n\">tags</span><span class=\"p\">,</span> <span
  class=\"nb\">list</span><span class=\"p\">)</span> <span class=\"k\">else</span>
  <span class=\"p\">[</span><span class=\"n\">tags</span><span class=\"p\">]</span>\n
  \                           <span class=\"k\">for</span> <span class=\"n\">a</span>
  <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">m</span><span class=\"o\">.</span><span class=\"n\">posts</span>\n                            <span
  class=\"k\">if</span> <span class=\"p\">(</span><span class=\"n\">tags</span> <span
  class=\"o\">:=</span> <span class=\"n\">a</span><span class=\"o\">.</span><span
  class=\"n\">get</span><span class=\"p\">(</span><span class=\"n\">attr</span><span
  class=\"p\">,</span> <span class=\"kc\">None</span><span class=\"p\">))</span> <span
  class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span>\n
  \                       <span class=\"p\">],</span>\n                    <span class=\"p\">),</span>\n
  \               <span class=\"p\">)</span>\n                <span class=\"k\">if</span>
  <span class=\"nb\">len</span><span class=\"p\">(</span><span class=\"n\">posts</span><span
  class=\"p\">)</span> <span class=\"o\">&gt;</span> <span class=\"mi\">0</span><span
  class=\"p\">:</span>\n                    <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">grid</span><span class=\"o\">.</span><span class=\"n\">add_row</span><span
  class=\"p\">()</span>\n                    <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">grid</span><span class=\"o\">.</span><span class=\"n\">add_row</span><span
  class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">\"[bold gold1]</span><span
  class=\"si\">{</span><span class=\"n\">attr</span><span class=\"o\">.</span><span
  class=\"n\">upper</span><span class=\"p\">()</span><span class=\"si\">}</span><span
  class=\"s2\">[/]\"</span><span class=\"p\">)</span>\n                    <span class=\"k\">for</span>
  <span class=\"n\">post</span><span class=\"p\">,</span> <span class=\"n\">count</span>
  <span class=\"ow\">in</span> <span class=\"n\">Counter</span><span class=\"p\">(</span><span
  class=\"n\">posts</span><span class=\"p\">)</span><span class=\"o\">.</span><span
  class=\"n\">most_common</span><span class=\"p\">():</span>\n                        <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">grid</span><span
  class=\"o\">.</span><span class=\"n\">add_row</span><span class=\"p\">(</span><span
  class=\"sa\">f</span><span class=\"s1\">'</span><span class=\"si\">{</span><span
  class=\"n\">count</span><span class=\"si\">}</span><span class=\"s1\"> </span><span
  class=\"si\">{</span><span class=\"s2\">\" \"</span><span class=\"o\">*</span><span
  class=\"p\">(</span><span class=\"mi\">3</span><span class=\"o\">-</span><span class=\"nb\">len</span><span
  class=\"p\">(</span><span class=\"nb\">str</span><span class=\"p\">(</span><span
  class=\"n\">count</span><span class=\"p\">)))</span><span class=\"si\">}</span><span
  class=\"s1\"> </span><span class=\"si\">{</span><span class=\"n\">post</span><span
  class=\"si\">}</span><span class=\"s1\">'</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"__rich__\" style=\"margin:0;padding:.5rem
  1rem;\"><strong>rich</strong> <em class=\"small\">method</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"<strong>rich</strong>
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
  <span class=\"nf\">__rich__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Union</span><span
  class=\"p\">[</span><span class=\"n\">Panel</span><span class=\"p\">,</span> <span
  class=\"n\">Table</span><span class=\"p\">]:</span>\n                <span class=\"n\">grid</span>
  <span class=\"o\">=</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">get_grid</span><span class=\"p\">()</span>\n\n                <span
  class=\"k\">if</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">simple</span><span class=\"p\">:</span>\n                    <span class=\"k\">return</span>
  <span class=\"n\">grid</span>\n                <span class=\"k\">else</span><span
  class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
  class=\"n\">Panel</span><span class=\"p\">(</span>\n                        <span
  class=\"n\">grid</span><span class=\"p\">,</span>\n                        <span
  class=\"n\">title</span><span class=\"o\">=</span><span class=\"s2\">\"[gold1]summary[/]\"</span><span
  class=\"p\">,</span>\n                        <span class=\"n\">border_style</span><span
  class=\"o\">=</span><span class=\"s2\">\"magenta\"</span><span class=\"p\">,</span>\n
  \                       <span class=\"n\">expand</span><span class=\"o\">=</span><span
  class=\"kc\">False</span><span class=\"p\">,</span>\n                    <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"get_summary\" style=\"margin:0;padding:.5rem
  1rem;\">get_summary <em class=\"small\">function</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"get_summary
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
  <span class=\"nf\">get_summary</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">):</span>\n                <span class=\"k\">try</span><span class=\"p\">:</span>\n
  \                   <span class=\"k\">return</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">_summary</span>\n                <span class=\"k\">except</span>
  <span class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n                    <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_summary</span><span
  class=\"p\">:</span> <span class=\"n\">Summary</span> <span class=\"o\">=</span>
  <span class=\"n\">Summary</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">)</span>\n                    <span class=\"k\">return</span> <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_summary</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"summary\" style=\"margin:0;padding:.5rem
  1rem;\">summary <em class=\"small\">function</em></h2>\nshow the application summary\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"summary
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
  <span class=\"nf\">summary</span><span class=\"p\">():</span>\n                <span
  class=\"s2\">\"show the application summary\"</span>\n                <span class=\"kn\">from</span>
  <span class=\"nn\">rich</span> <span class=\"kn\">import</span> <span class=\"nb\">print</span>\n\n
  \               <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">quiet</span>
  <span class=\"o\">=</span> <span class=\"kc\">True</span>\n\n                <span
  class=\"nb\">print</span><span class=\"p\">(</span><span class=\"n\">Summary</span><span
  class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">,</span> <span
  class=\"n\">simple</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
  class=\"p\">))</span>\n</pre></div>\n\n</pre>\n</main>\n<footer>\xA9 2024</footer>\n</body></html>"
published: true
slug: markata/cli/summary
title: Summary.Py


---

## Run it

``` bash
python -m markata.cli.summary
```

## Configuration

There are two main things currently supported by summary, It can count the
number of posts based on a filter (`filter_count'), and it can automatically
list all the values of an attribute and the number of posts that have that
attribute (`grid_attr`).

### grid_attr

`grid_attr` will map over all posts, find all values for each attribute
configured, then report the number of posts for each value.

``` toml
[markata.summary]
grid_attr = ['tags', 'series']
```

Example output that will be shown in the summary, counting all posts with each
tag value.

```
TAGS
247  python
90   linux
68   cli
49   kedro
46   bash
```

### filter_counts

`filter_count` will pass a given filter into `markata.map` and return the number
of posts.

```
[[markata.summary.filter_count]]
name='drafts'
filter="published == 'False'"
color='red'

[[markata.summary.filter_count]]
name='articles'
color='dark_orange'

[[markata.summary.filter_count]]
name='py_modules'
filter='"plugin" not in slug and "docs" not in str(path)'
color="yellow1"

[markata.summary.filter_count.published]
filter="published == 'True'"
color='green1'

[markata.summary.filter_count.plugins]
filter='"plugin" in slug and "docs" not in str(path)'
color="blue"

[markata.summary.filter_count.docs]
filter="'docs' in str(path)"
color='purple'
```

Example output might look like this, showing the number of posts that contained
within each filter specified.

```
8 drafts
66 articles
20 py_modules
58 published
38 plugins
8 docs
```


!! class <h2 id='FilterCount' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>FilterCount <em class='small'>class</em></h2>

???+ source "FilterCount <em class='small'>source</em>"

```python

        class FilterCount(pydantic.BaseModel):
            name: str
            filter: str = "True"
            color: str = "white"
```


!! class <h2 id='SummaryConfig' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>SummaryConfig <em class='small'>class</em></h2>

???+ source "SummaryConfig <em class='small'>source</em>"

```python

        class SummaryConfig(pydantic.BaseModel):
            grid_attr: List[str] = ["tags", "series"]
            filter_count: List[FilterCount] = FilterCount(
                name="drafts", filter="published == 'False'", color="red"
            )
```


!! class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Config <em class='small'>class</em></h2>

???+ source "Config <em class='small'>source</em>"

```python

        class Config(pydantic.BaseModel):
            summary: SummaryConfig = SummaryConfig()
```


!! function <h2 id='config_model' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>config_model <em class='small'>function</em></h2>

???+ source "config_model <em class='small'>source</em>"

```python

        def config_model(markata: "Markata") -> None:
            markata.config_models.append(Config)
```


!! class <h2 id='Summary' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Summary <em class='small'>class</em></h2>

???+ source "Summary <em class='small'>source</em>"

```python

        class Summary:
            def __init__(self, m: "Markata", simple: bool = False) -> None:
                self.m = m
                self.simple = simple

            def get_grid(self) -> None:
                "create a rich grid to display the summary"
                self.grid = Table.grid(expand=True)

                for filter_count in self.m.config.summary.filter_count:
                    self.filter_count(filter_count)

                for attr in self.m.config.summary.grid_attr:
                    self.grid_attr(attr)

                return self.grid

            def filter_count(
                self,
                fc: FilterCount,
            ) -> None:
                "add a row in the grid for the number of items in a filter config"
                self.grid.add_row(
                    f"[{fc.color}]{len(self.m.map(filter=fc.filter))}[/] {fc.name}"
                )

            def grid_attr(self, attr: str) -> None:
                "add attribute the the object grid"
                posts = list(
                    flatten(
                        [
                            tags if isinstance(tags, list) else [tags]
                            for a in self.m.posts
                            if (tags := a.get(attr, None)) is not None
                        ],
                    ),
                )
                if len(posts) > 0:
                    self.grid.add_row()
                    self.grid.add_row(f"[bold gold1]{attr.upper()}[/]")
                    for post, count in Counter(posts).most_common():
                        self.grid.add_row(f'{count} {" "*(3-len(str(count)))} {post}')

            def __rich__(self) -> Union[Panel, Table]:
                grid = self.get_grid()

                if self.simple:
                    return grid
                else:
                    return Panel(
                        grid,
                        title="[gold1]summary[/]",
                        border_style="magenta",
                        expand=False,
                    )
```


!! function <h2 id='configure' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>configure <em class='small'>function</em></h2>

???+ source "configure <em class='small'>source</em>"

```python

        def configure(markata: "Markata") -> None:
            def get_summary(self):
                try:
                    return self._summary
                except AttributeError:
                    self._summary: Summary = Summary(self)
                    return self._summary

            from markata import Markata

            Markata.summary = property(get_summary)
```


!! function <h2 id='cli' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>cli <em class='small'>function</em></h2>
    Markata hook to implement base cli commands.
???+ source "cli <em class='small'>source</em>"

```python

        def cli(app: typer.Typer, markata: "Markata") -> None:
            """
            Markata hook to implement base cli commands.
            """
            summary_app = typer.Typer()
            app.add_typer(summary_app, name="summary")

            @summary_app.callback(invoke_without_command=True)
            def summary():
                "show the application summary"
                from rich import print

                markata.console.quiet = True

                print(Summary(markata, simple=True))
```


!! method <h2 id='__init__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__init__ <em class='small'>method</em></h2>

???+ source "__init__ <em class='small'>source</em>"

```python

        def __init__(self, m: "Markata", simple: bool = False) -> None:
                self.m = m
                self.simple = simple
```


!! method <h2 id='get_grid' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get_grid <em class='small'>method</em></h2>
    create a rich grid to display the summary
???+ source "get_grid <em class='small'>source</em>"

```python

        def get_grid(self) -> None:
                "create a rich grid to display the summary"
                self.grid = Table.grid(expand=True)

                for filter_count in self.m.config.summary.filter_count:
                    self.filter_count(filter_count)

                for attr in self.m.config.summary.grid_attr:
                    self.grid_attr(attr)

                return self.grid
```


!! method <h2 id='filter_count' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>filter_count <em class='small'>method</em></h2>
    add a row in the grid for the number of items in a filter config
???+ source "filter_count <em class='small'>source</em>"

```python

        def filter_count(
                self,
                fc: FilterCount,
            ) -> None:
                "add a row in the grid for the number of items in a filter config"
                self.grid.add_row(
                    f"[{fc.color}]{len(self.m.map(filter=fc.filter))}[/] {fc.name}"
                )
```


!! method <h2 id='grid_attr' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>grid_attr <em class='small'>method</em></h2>
    add attribute the the object grid
???+ source "grid_attr <em class='small'>source</em>"

```python

        def grid_attr(self, attr: str) -> None:
                "add attribute the the object grid"
                posts = list(
                    flatten(
                        [
                            tags if isinstance(tags, list) else [tags]
                            for a in self.m.posts
                            if (tags := a.get(attr, None)) is not None
                        ],
                    ),
                )
                if len(posts) > 0:
                    self.grid.add_row()
                    self.grid.add_row(f"[bold gold1]{attr.upper()}[/]")
                    for post, count in Counter(posts).most_common():
                        self.grid.add_row(f'{count} {" "*(3-len(str(count)))} {post}')
```


!! method <h2 id='__rich__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__rich__ <em class='small'>method</em></h2>

???+ source "__rich__ <em class='small'>source</em>"

```python

        def __rich__(self) -> Union[Panel, Table]:
                grid = self.get_grid()

                if self.simple:
                    return grid
                else:
                    return Panel(
                        grid,
                        title="[gold1]summary[/]",
                        border_style="magenta",
                        expand=False,
                    )
```


!! function <h2 id='get_summary' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get_summary <em class='small'>function</em></h2>

???+ source "get_summary <em class='small'>source</em>"

```python

        def get_summary(self):
                try:
                    return self._summary
                except AttributeError:
                    self._summary: Summary = Summary(self)
                    return self._summary
```


!! function <h2 id='summary' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>summary <em class='small'>function</em></h2>
    show the application summary
???+ source "summary <em class='small'>source</em>"

```python

        def summary():
                "show the application summary"
                from rich import print

                markata.console.quiet = True

                print(Summary(markata, simple=True))
```
