---
content: "Renders your markdown as a jinja template during pre_render.\n\nThe markata
  instance is passed into the template, giving you access to things\nsuch as all of
  your articles, config, and this post as post.\n\n# Examples\n\nfirst we can grab
  a few things out of the frontmatter of this post.\n\n``` markdown\n# {{ post.title
  }}\n{{ post.description }}\n```\n\n# one-liner list of links\n\nThis one-liner will
  render a list of markdown links into your markdown at build\ntime.  It's quite handy
  to pop into posts.\n\n``` markdown\n{{ '\\n'.join(markata.map('f\"* [{title}]({slug})\"',
  sort='slug')) }}\n```\n\n# jinja for to markdown list of links\n\nSometimes quoting
  things like your filters are hard to do in a one line without\nrunning out of quote
  variants.  Jinja for loops can make this much easier.\n\n``` markdown\n{% for post
  in markata.map('post', filter='\"git\" in tags') %}\n* [{{ post.title }}]({{ post.slug
  }})\n{% endfor %}\n```\n\n# jinja for to html list of links\n\nSince markdown is
  a superset of html, you can just render out html into your\npost and it is still
  valid.\n\n``` markdown\n<ul>\n{% for post in markata.map('post', filter='\"git\"
  in tags') %}\n    <li><a href=\"{{ post.slug }}\">{{ post.title }}</a></li>\n{%
  endfor %}\n</ul>\n```\n\n# Ignoring files\n\nIt is possible to ignore files by adding
  an ignore to your `markata.jinja_md`\nconfig in your `markata.toml` file.  This
  ignore follows the `gitwildmatch`\nrules, so think of it the same as writing a gitignore.\n\n```toml\n[markata.jinja_md]\nignore=[\n'jinja_md.md',\n]\n```\n\n!!note\n
  \ Docs such as this jinja_md.py file will get converted to jinja_md.md during\n
  \ build time, so use `.md` extensions instead of `.py`.\n\n# Ignoring a single file\n\nYou
  can also ignore a single file right from the articles frontmatter, by\nadding `jinja:
  false`.\n\n```markdown\n---\njinja: false\n\n---\n```\n\n# Escaping\n\nSometimes
  you want the ability to have jinja templates in a post, but also the\nability to
  keep a raw jinja template.  There are a couple of techniques that\nare covered mroe
  in the jinja docs for\n[escaping](https://jinja.palletsprojects.com/en/3.1.x/templates/#escaping)\n\n```markdown\n{%
  raw %}\n{{ '\\n'.join(markata.map('f\"* [{title}]({slug})\"', sort='slug')) }}\n{%
  endraw %}\n\n{{ '{{' }} '\\n'.join(markata.map('f\"* [{title}]({slug})\"', sort='slug'))
  {{ '}}' }}\n```\n\n# Creating a jinja extension\n\nHere is a bit of a boilerplate
  example of a jinja extension.\n\n``` python\nfrom jinja2 import nodes\nfrom jinja2.ext
  import Extension\n\nclass ExampleExtension(Extension):\n    tags = {\"example\"}\n\n
  \   def __init__(self, environment):\n        super().__init__(environment)\n\n
  \   def parse(self, parser):\n        line_number = next(parser.stream).lineno\n
  \       arg = [parser.parse_expression()]\n        return nodes.CallBlock(self.call_method(\"run\",
  arg), [], [], \"\").set_lineno(\n            line_number\n        )\n\n    def run(self,
  arg, caller):\n        return f'hello {arg}'\n```\n\nSo that markata picks up your
  extension, you will need to register an\nentrypoint named `markata.jinja_md`.  Once
  installed markata will automatically\nload this extension to its list of  jinja
  extensions.\n\n``` toml\n[project.entry-points.\"markata.jinja_md\"]\nmarkta_gh
  = \"example_extension:ExampleExtension\"\n```\n\nOnce you have your extension created
  and ready to use you can use it in your\nmarkdown.\n\n``` markdown\n{% example world
  %}\n```\n\n\n!! function <h2 id='register_jinja_extensions' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>register_jinja_extensions <em class='small'>function</em></h2>\n
  \   Gets jinja extensions from entrypoints and loads them in.\n\n    Returns: List
  of jinja Extensions\n???+ source \"register_jinja_extensions <em class='small'>source</em>\"\n\n```python\n\n
  \       def register_jinja_extensions(config: dict) -> List[Extension]:\n            \"\"\"\n
  \           Gets jinja extensions from entrypoints and loads them in.\n\n            Returns:
  List of jinja Extensions\n            \"\"\"\n\n            return [\n                ep.load()
  for ep in pkg_resources.iter_entry_points(group=\"markata.jinja_md\")\n            ]\n```\n\n\n!!
  class <h2 id='IncludeRawExtension' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>IncludeRawExtension <em class='small'>class</em></h2>\n\n???+ source \"IncludeRawExtension
  <em class='small'>source</em>\"\n\n```python\n\n        class IncludeRawExtension(Extension):\n
  \           tags = {\"include_raw\"}\n\n            def parse(self, parser):\n                line_number
  = next(parser.stream).lineno\n                file = [parser.parse_expression()]\n
  \               return nodes.CallBlock(\n                    self.call_method(\"_read_file\",
  file),\n                    [],\n                    [],\n                    \"\",\n
  \               ).set_lineno(line_number)\n\n            def _read_file(self, file,
  caller):\n                return Path(file).read_text()\n```\n\n\n!! class <h2 id='_SilentUndefined'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_SilentUndefined <em
  class='small'>class</em></h2>\n    silence undefined variable errors in jinja templates.\n\n
  \   # Example\n    ```python\n    template = '{{ variable }}'\n    article.content
  = Template( template, undefined=_SilentUndefined).render()\n    ```\n???+ source
  \"_SilentUndefined <em class='small'>source</em>\"\n\n```python\n\n        class
  _SilentUndefined(Undefined):\n            \"\"\"\n            silence undefined
  variable errors in jinja templates.\n\n            # Example\n            ```python\n
  \           template = '{{ variable }}'\n            article.content = Template(
  template, undefined=_SilentUndefined).render()\n            ```\n            \"\"\"\n\n
  \           def _fail_with_undefined_error(self, *args, **kwargs):\n                return
  \"\"\n```\n\n\n!! class <h2 id='PostTemplateSyntaxError' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>PostTemplateSyntaxError <em class='small'>class</em></h2>\n
  \   Custom error message for post template syntax errors.\n???+ source \"PostTemplateSyntaxError
  <em class='small'>source</em>\"\n\n```python\n\n        class PostTemplateSyntaxError(TemplateSyntaxError):\n
  \           \"\"\"\n            Custom error message for post template syntax errors.\n
  \           \"\"\"\n```\n\n\n!! class <h2 id='JinjaMd' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>JinjaMd <em class='small'>class</em></h2>\n\n???+
  source \"JinjaMd <em class='small'>source</em>\"\n\n```python\n\n        class JinjaMd(pydantic.BaseModel):\n
  \           jinja: bool = True\n```\n\n\n!! function <h2 id='post_model' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>post_model <em class='small'>function</em></h2>\n\n???+
  source \"post_model <em class='small'>source</em>\"\n\n```python\n\n        def
  post_model(markata: \"Markata\") -> None:\n            markata.post_models.append(JinjaMd)\n```\n\n\n!!
  function <h2 id='pre_render' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>pre_render <em class='small'>function</em></h2>\n    jinja_md hook for markata
  to render your markdown post as a jinja template.\n\n    The post itself is exposed
  as `post`, and the markata instance is exposed\n    as `markata`.\n???+ source \"pre_render
  <em class='small'>source</em>\"\n\n```python\n\n        def pre_render(markata:
  \"Markata\") -> None:\n            \"\"\"\n            jinja_md hook for markata
  to render your markdown post as a jinja template.\n\n            The post itself
  is exposed as `post`, and the markata instance is exposed\n            as `markata`.\n
  \           \"\"\"\n\n            config = markata.config.jinja_md\n            ignore_spec
  = pathspec.PathSpec.from_lines(\"gitwildmatch\", config.ignore)\n            # for
  article in markata.iter_articles(description=\"jinja_md\"):\n            jinja_env
  = jinja2.Environment(\n                extensions=[IncludeRawExtension, *register_jinja_extensions(config)],\n
  \           )\n\n            for article in markata.articles:\n                if
  article.get(\"jinja\", True) and not ignore_spec.match_file(article[\"path\"]):\n
  \                   try:\n                        key = markata.make_hash(article.content)\n
  \                       content_from_cache = markata.precache.get(key)\n                        if
  content_from_cache is None:\n                            article.content = jinja_env.from_string(article.content).render(\n
  \                               __version__=__version__,\n                                **article,\n
  \                               post=article,\n                            )\n                            with
  markata.cache:\n                                markata.cache.set(key, article.content)\n
  \                       else:\n                            article.content = content_from_cache\n
  \                       # prevent double rendering\n                        article.jinja
  = False\n                    except TemplateSyntaxError as e:\n                        errorline
  = article.content.split(\"\\n\")[e.lineno - 1]\n                        msg = f\"\"\"\n
  \                       Error while processing post {article['path']}\n\n                        {errorline}\n
  \                       \"\"\"\n\n                        raise PostTemplateSyntaxError(msg,
  lineno=e.lineno)\n                    except UndefinedError as e:\n                        raise
  UndefinedError(f'{e} in {article[\"path\"]}')\n```\n\n\n!! class <h2 id='JinjaMdConfig'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>JinjaMdConfig <em
  class='small'>class</em></h2>\n\n???+ source \"JinjaMdConfig <em class='small'>source</em>\"\n\n```python\n\n
  \       class JinjaMdConfig(pydantic.BaseModel):\n            ignore: List[str]
  = []\n```\n\n\n!! class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>Config <em class='small'>class</em></h2>\n\n???+ source \"Config <em class='small'>source</em>\"\n\n```python\n\n
  \       class Config(pydantic.BaseModel):\n            jinja_md: JinjaMdConfig =
  JinjaMdConfig()\n```\n\n\n!! function <h2 id='config_model' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>config_model <em class='small'>function</em></h2>\n\n???+
  source \"config_model <em class='small'>source</em>\"\n\n```python\n\n        def
  config_model(markata: \"Markata\") -> None:\n            markata.config_models.append(Config)\n```\n\n\n!!
  method <h2 id='parse' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>parse
  <em class='small'>method</em></h2>\n\n???+ source \"parse <em class='small'>source</em>\"\n\n```python\n\n
  \       def parse(self, parser):\n                line_number = next(parser.stream).lineno\n
  \               file = [parser.parse_expression()]\n                return nodes.CallBlock(\n
  \                   self.call_method(\"_read_file\", file),\n                    [],\n
  \                   [],\n                    \"\",\n                ).set_lineno(line_number)\n```\n\n\n!!
  method <h2 id='_read_file' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>_read_file <em class='small'>method</em></h2>\n\n???+ source \"_read_file
  <em class='small'>source</em>\"\n\n```python\n\n        def _read_file(self, file,
  caller):\n                return Path(file).read_text()\n```\n\n\n!! method <h2
  id='_fail_with_undefined_error' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>_fail_with_undefined_error <em class='small'>method</em></h2>\n\n???+ source
  \"_fail_with_undefined_error <em class='small'>source</em>\"\n\n```python\n\n        def
  _fail_with_undefined_error(self, *args, **kwargs):\n                return \"\"\n```\n\n"
date: 0001-01-01
description: Renders your markdown as a jinja template during pre The markata instance
  is passed into the template, giving you access to things first we can grab a few
  thing
html: "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n<title>Jinja_Md.Py</title>\n<meta
  charset=\"utf-8\"/>\n<meta content=\"width=device-width, initial-scale=1\" name=\"viewport\"/>\n<meta
  content=\"Renders your markdown as a jinja template during pre The markata instance
  is passed into the template, giving you access to things first we can grab a few
  thing\" name=\"description\"/>\n<link href=\"/static/favicon.ico\" rel=\"icon\"
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
  name=\"og:type\" property=\"og:type\"/><meta content=\"Renders your markdown as
  a jinja template during pre The markata instance is passed into the template, giving
  you access to things first we can grab a few thing\" name=\"description\" property=\"description\"/><meta
  content=\"Renders your markdown as a jinja template during pre The markata instance
  is passed into the template, giving you access to things first we can grab a few
  thing\" name=\"og:description\" property=\"og:description\"/><meta content=\"Renders
  your markdown as a jinja template during pre The markata instance is passed into
  the template, giving you access to things first we can grab a few thing\" name=\"twitter:description\"
  property=\"twitter:description\"/><meta content=\"Jinja_Md.Py | Markata\" name=\"og:title\"
  property=\"og:title\"/><meta content=\"Jinja_Md.Py | Markata\" name=\"twitter:title\"
  property=\"twitter:title\"/><meta content=\"https://markata.dev//markata/plugins/jinja-md-og.png\"
  name=\"og:image\" property=\"og:image\"/><meta content=\"https://markata.dev//markata/plugins/jinja-md-og.png\"
  name=\"twitter:image\" property=\"twitter:image\"/><meta content=\"1600\" name=\"og:image:width\"
  property=\"og:image:width\"/><meta content=\"900\" name=\"og:image:width\" property=\"og:image:width\"/><meta
  content=\"summary_large_image\" name=\"twitter:card\" property=\"twitter:card\"/><meta
  content=\"Markata\" name=\"og:site_name\" property=\"og:site_name\"/><meta content=\"@_waylonwalker\"
  name=\"twitter:creator\" property=\"twitter:creator\"/><meta content=\"Jinja_Md.Py\"
  name=\"title\" property=\"title\"/><meta content=\"markata 0.8.0.dev6\" name=\"generator\"
  property=\"generator\"/><link href=\"https://markata.dev//markata/plugins/jinja-md/\"
  rel=\"canonical\"/><meta content=\"https://markata.dev//markata/plugins/jinja-md/\"
  name=\"og:url\" property=\"og:url\"/></head>\n<body><nav>\n<a href=\"/\">markata</a>\n<a
  href=\"https://github.com/WaylonWalker/markata\">GitHub</a>\n</nav>\n<div>\n<label
  class=\"theme-switch\" for=\"checkbox-theme\" id=\"theme-switch\" title=\"light/dark
  mode toggle\">\n<input id=\"checkbox-theme\" type=\"checkbox\"/>\n<div class=\"slider
  round\"></div>\n</label>\n</div>\n<section class=\"title\">\n<h1 id=\"title\">\n
  \           Jinja_Md.Py \n            \n        </h1>\n</section>\n<main><p>Renders
  your markdown as a jinja template during pre_render.</p>\n<p>The markata instance
  is passed into the template, giving you access to things\nsuch as all of your articles,
  config, and this post as post.</p>\n<h1 id=\"examples\">Examples <a class=\"header-anchor\"
  href=\"#examples\"><svg aria-hidden=\"true\" class=\"heading-permalink\" fill=\"currentColor\"
  focusable=\"false\" height=\"1em\" viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
  d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>first we can grab a few
  things out of the frontmatter of this post.</p>\n<pre class=\"wrapper\">\n\n<div
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
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"gh\"># {{ post.title
  }}</span>\n{{ post.description }}\n</pre></div>\n\n</pre>\n<h1 id=\"one-liner-list-of-links\">one-liner
  list of links <a class=\"header-anchor\" href=\"#one-liner-list-of-links\"><svg
  aria-hidden=\"true\" class=\"heading-permalink\" fill=\"currentColor\" focusable=\"false\"
  height=\"1em\" viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
  d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>This one-liner will render
  a list of markdown links into your markdown at build\ntime.  It's quite handy to
  pop into posts.</p>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
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
  \       \n<div class=\"highlight\"><pre><span></span>{{ '\\n'.join(markata.map('f\"*
  [<span class=\"nt\">{title}</span>](<span class=\"na\">{slug}</span>)\"', sort='slug'))
  }}\n</pre></div>\n\n</pre>\n<h1 id=\"jinja-for-to-markdown-list-of-links\">jinja
  for to markdown list of links <a class=\"header-anchor\" href=\"#jinja-for-to-markdown-list-of-links\"><svg
  aria-hidden=\"true\" class=\"heading-permalink\" fill=\"currentColor\" focusable=\"false\"
  height=\"1em\" viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
  d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>Sometimes quoting things
  like your filters are hard to do in a one line without\nrunning out of quote variants.
  \ Jinja for loops can make this much easier.</p>\n<pre class=\"wrapper\">\n\n<div
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
  \       \n<div class=\"highlight\"><pre><span></span>{% for post in markata.map('post',
  filter='\"git\" in tags') %}\n<span class=\"k\">*</span><span class=\"w\"> </span>[<span
  class=\"nt\">{{ post.title }}</span>](<span class=\"na\">{{ post.slug }}</span>)\n{%
  endfor %}\n</pre></div>\n\n</pre>\n<h1 id=\"jinja-for-to-html-list-of-links\">jinja
  for to html list of links <a class=\"header-anchor\" href=\"#jinja-for-to-html-list-of-links\"><svg
  aria-hidden=\"true\" class=\"heading-permalink\" fill=\"currentColor\" focusable=\"false\"
  height=\"1em\" viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
  d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>Since markdown is a superset
  of html, you can just render out html into your\npost and it is still valid.</p>\n<pre
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
  \       \n<div class=\"highlight\"><pre><span></span>&lt;ul&gt;\n{% for post in
  markata.map('post', filter='\"git\" in tags') %}\n    &lt;li&gt;&lt;a href=\"{{
  post.slug }}\"&gt;{{ post.title }}&lt;/a&gt;&lt;/li&gt;\n{% endfor %}\n&lt;/ul&gt;\n</pre></div>\n\n</pre>\n<h1
  id=\"ignoring-files\">Ignoring files <a class=\"header-anchor\" href=\"#ignoring-files\"><svg
  aria-hidden=\"true\" class=\"heading-permalink\" fill=\"currentColor\" focusable=\"false\"
  height=\"1em\" viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
  d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>It is possible to ignore
  files by adding an ignore to your <code>markata.jinja_md</code>\nconfig in your
  <code>markata.toml</code> file.  This ignore follows the <code>gitwildmatch</code>\nrules,
  so think of it the same as writing a gitignore.</p>\n<pre class=\"wrapper\">\n\n<div
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
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata.jinja_md]</span>\n<span
  class=\"n\">ignore</span><span class=\"o\">=</span><span class=\"p\">[</span>\n<span
  class=\"s1\">'jinja_md.md'</span><span class=\"p\">,</span>\n<span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n<p>!!note\nDocs
  such as this jinja_md.py file will get converted to jinja_md.md during\nbuild time,
  so use <code>.md</code> extensions instead of <code>.py</code>.</p>\n<h1 id=\"ignoring-a-single-file\">Ignoring
  a single file <a class=\"header-anchor\" href=\"#ignoring-a-single-file\"><svg aria-hidden=\"true\"
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
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>You can also ignore a
  single file right from the articles frontmatter, by\nadding <code>jinja: false</code>.</p>\n<pre
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
  \       \n<div class=\"highlight\"><pre><span></span>---\njinja: false\n\n---\n</pre></div>\n\n</pre>\n<h1
  id=\"escaping\">Escaping <a class=\"header-anchor\" href=\"#escaping\"><svg aria-hidden=\"true\"
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
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>Sometimes you want the
  ability to have jinja templates in a post, but also the\nability to keep a raw jinja
  template.  There are a couple of techniques that\nare covered mroe in the jinja
  docs for\n<a href=\"https://jinja.palletsprojects.com/en/3.1.x/templates/#escaping\">escaping</a></p>\n<pre
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
  \       \n<div class=\"highlight\"><pre><span></span>{% raw %}\n{{ '\\n'.join(markata.map('f\"*
  [<span class=\"nt\">{title}</span>](<span class=\"na\">{slug}</span>)\"', sort='slug'))
  }}\n{% endraw %}\n\n{{ '{{' }} '\\n'.join(markata.map('f\"* [<span class=\"nt\">{title}</span>](<span
  class=\"na\">{slug}</span>)\"', sort='slug')) {{ '}}' }}\n</pre></div>\n\n</pre>\n<h1
  id=\"creating-a-jinja-extension\">Creating a jinja extension <a class=\"header-anchor\"
  href=\"#creating-a-jinja-extension\"><svg aria-hidden=\"true\" class=\"heading-permalink\"
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
  3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632 1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>Here
  is a bit of a boilerplate example of a jinja extension.</p>\n<pre class=\"wrapper\">\n\n<div
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
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"kn\">from</span>
  <span class=\"nn\">jinja2</span> <span class=\"kn\">import</span> <span class=\"n\">nodes</span>\n<span
  class=\"kn\">from</span> <span class=\"nn\">jinja2.ext</span> <span class=\"kn\">import</span>
  <span class=\"n\">Extension</span>\n\n<span class=\"k\">class</span> <span class=\"nc\">ExampleExtension</span><span
  class=\"p\">(</span><span class=\"n\">Extension</span><span class=\"p\">):</span>\n
  \   <span class=\"n\">tags</span> <span class=\"o\">=</span> <span class=\"p\">{</span><span
  class=\"s2\">\"example\"</span><span class=\"p\">}</span>\n\n    <span class=\"k\">def</span>
  <span class=\"fm\">__init__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">,</span> <span class=\"n\">environment</span><span class=\"p\">):</span>\n
  \       <span class=\"nb\">super</span><span class=\"p\">()</span><span class=\"o\">.</span><span
  class=\"fm\">__init__</span><span class=\"p\">(</span><span class=\"n\">environment</span><span
  class=\"p\">)</span>\n\n    <span class=\"k\">def</span> <span class=\"nf\">parse</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
  class=\"n\">parser</span><span class=\"p\">):</span>\n        <span class=\"n\">line_number</span>
  <span class=\"o\">=</span> <span class=\"nb\">next</span><span class=\"p\">(</span><span
  class=\"n\">parser</span><span class=\"o\">.</span><span class=\"n\">stream</span><span
  class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">lineno</span>\n
  \       <span class=\"n\">arg</span> <span class=\"o\">=</span> <span class=\"p\">[</span><span
  class=\"n\">parser</span><span class=\"o\">.</span><span class=\"n\">parse_expression</span><span
  class=\"p\">()]</span>\n        <span class=\"k\">return</span> <span class=\"n\">nodes</span><span
  class=\"o\">.</span><span class=\"n\">CallBlock</span><span class=\"p\">(</span><span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">call_method</span><span
  class=\"p\">(</span><span class=\"s2\">\"run\"</span><span class=\"p\">,</span>
  <span class=\"n\">arg</span><span class=\"p\">),</span> <span class=\"p\">[],</span>
  <span class=\"p\">[],</span> <span class=\"s2\">\"\"</span><span class=\"p\">)</span><span
  class=\"o\">.</span><span class=\"n\">set_lineno</span><span class=\"p\">(</span>\n
  \           <span class=\"n\">line_number</span>\n        <span class=\"p\">)</span>\n\n
  \   <span class=\"k\">def</span> <span class=\"nf\">run</span><span class=\"p\">(</span><span
  class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"n\">arg</span><span
  class=\"p\">,</span> <span class=\"n\">caller</span><span class=\"p\">):</span>\n
  \       <span class=\"k\">return</span> <span class=\"sa\">f</span><span class=\"s1\">'hello
  </span><span class=\"si\">{</span><span class=\"n\">arg</span><span class=\"si\">}</span><span
  class=\"s1\">'</span>\n</pre></div>\n\n</pre>\n<p>So that markata picks up your
  extension, you will need to register an\nentrypoint named <code>markata.jinja_md</code>.
  \ Once installed markata will automatically\nload this extension to its list of
  \ jinja extensions.</p>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
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
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[project.entry-points.</span><span
  class=\"s2\">\"markata.jinja_md\"</span><span class=\"k\">]</span>\n<span class=\"n\">markta_gh</span><span
  class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s2\">\"example_extension:ExampleExtension\"</span>\n</pre></div>\n\n</pre>\n<p>Once
  you have your extension created and ready to use you can use it in your\nmarkdown.</p>\n<pre
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
  \       \n<div class=\"highlight\"><pre><span></span>{% example world %}\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"register_jinja_extensions\" style=\"margin:0;padding:.5rem
  1rem;\">register_jinja_extensions <em class=\"small\">function</em></h2>\nGets jinja
  extensions from entrypoints and loads them in.\n<pre><code>Returns: List of jinja
  Extensions\n</code></pre>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
  class=\"admonition-title\">\"register_jinja_extensions <em class=\"small\">source</em>\"</p>\n</div>\n<pre
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
  <span class=\"nf\">register_jinja_extensions</span><span class=\"p\">(</span><span
  class=\"n\">config</span><span class=\"p\">:</span> <span class=\"nb\">dict</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">List</span><span
  class=\"p\">[</span><span class=\"n\">Extension</span><span class=\"p\">]:</span>\n<span
  class=\"w\">            </span><span class=\"sd\">\"\"\"</span>\n<span class=\"sd\">
  \           Gets jinja extensions from entrypoints and loads them in.</span>\n\n<span
  class=\"sd\">            Returns: List of jinja Extensions</span>\n<span class=\"sd\">
  \           \"\"\"</span>\n\n            <span class=\"k\">return</span> <span class=\"p\">[</span>\n
  \               <span class=\"n\">ep</span><span class=\"o\">.</span><span class=\"n\">load</span><span
  class=\"p\">()</span> <span class=\"k\">for</span> <span class=\"n\">ep</span> <span
  class=\"ow\">in</span> <span class=\"n\">pkg_resources</span><span class=\"o\">.</span><span
  class=\"n\">iter_entry_points</span><span class=\"p\">(</span><span class=\"n\">group</span><span
  class=\"o\">=</span><span class=\"s2\">\"markata.jinja_md\"</span><span class=\"p\">)</span>\n
  \           <span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n<p>!! class </p><h2
  class=\"admonition-title\" id=\"IncludeRawExtension\" style=\"margin:0;padding:.5rem
  1rem;\">IncludeRawExtension <em class=\"small\">class</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"IncludeRawExtension
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
  <span class=\"nc\">IncludeRawExtension</span><span class=\"p\">(</span><span class=\"n\">Extension</span><span
  class=\"p\">):</span>\n            <span class=\"n\">tags</span> <span class=\"o\">=</span>
  <span class=\"p\">{</span><span class=\"s2\">\"include_raw\"</span><span class=\"p\">}</span>\n\n
  \           <span class=\"k\">def</span> <span class=\"nf\">parse</span><span class=\"p\">(</span><span
  class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"n\">parser</span><span
  class=\"p\">):</span>\n                <span class=\"n\">line_number</span> <span
  class=\"o\">=</span> <span class=\"nb\">next</span><span class=\"p\">(</span><span
  class=\"n\">parser</span><span class=\"o\">.</span><span class=\"n\">stream</span><span
  class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">lineno</span>\n
  \               <span class=\"n\">file</span> <span class=\"o\">=</span> <span class=\"p\">[</span><span
  class=\"n\">parser</span><span class=\"o\">.</span><span class=\"n\">parse_expression</span><span
  class=\"p\">()]</span>\n                <span class=\"k\">return</span> <span class=\"n\">nodes</span><span
  class=\"o\">.</span><span class=\"n\">CallBlock</span><span class=\"p\">(</span>\n
  \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">call_method</span><span class=\"p\">(</span><span class=\"s2\">\"_read_file\"</span><span
  class=\"p\">,</span> <span class=\"n\">file</span><span class=\"p\">),</span>\n
  \                   <span class=\"p\">[],</span>\n                    <span class=\"p\">[],</span>\n
  \                   <span class=\"s2\">\"\"</span><span class=\"p\">,</span>\n                <span
  class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">set_lineno</span><span
  class=\"p\">(</span><span class=\"n\">line_number</span><span class=\"p\">)</span>\n\n
  \           <span class=\"k\">def</span> <span class=\"nf\">_read_file</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
  class=\"n\">file</span><span class=\"p\">,</span> <span class=\"n\">caller</span><span
  class=\"p\">):</span>\n                <span class=\"k\">return</span> <span class=\"n\">Path</span><span
  class=\"p\">(</span><span class=\"n\">file</span><span class=\"p\">)</span><span
  class=\"o\">.</span><span class=\"n\">read_text</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n<p>!!
  class </p><h2 class=\"admonition-title\" id=\"_SilentUndefined\" style=\"margin:0;padding:.5rem
  1rem;\">_SilentUndefined <em class=\"small\">class</em></h2>\nsilence undefined
  variable errors in jinja templates.\n<pre><code># Example\n```python\ntemplate =
  '{{ variable }}'\narticle.content = Template( template, undefined=_SilentUndefined).render()\n```\n</code></pre>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"_SilentUndefined
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
  <span class=\"nc\">_SilentUndefined</span><span class=\"p\">(</span><span class=\"n\">Undefined</span><span
  class=\"p\">):</span>\n<span class=\"w\">            </span><span class=\"sd\">\"\"\"</span>\n<span
  class=\"sd\">            silence undefined variable errors in jinja templates.</span>\n\n<span
  class=\"sd\">            # Example</span>\n<span class=\"sd\">            ```python</span>\n<span
  class=\"sd\">            template = '{{ variable }}'</span>\n<span class=\"sd\">
  \           article.content = Template( template, undefined=_SilentUndefined).render()</span>\n<span
  class=\"sd\">            ```</span>\n<span class=\"sd\">            \"\"\"</span>\n\n
  \           <span class=\"k\">def</span> <span class=\"nf\">_fail_with_undefined_error</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
  class=\"o\">*</span><span class=\"n\">args</span><span class=\"p\">,</span> <span
  class=\"o\">**</span><span class=\"n\">kwargs</span><span class=\"p\">):</span>\n
  \               <span class=\"k\">return</span> <span class=\"s2\">\"\"</span>\n</pre></div>\n\n</pre>\n<p>!!
  class </p><h2 class=\"admonition-title\" id=\"PostTemplateSyntaxError\" style=\"margin:0;padding:.5rem
  1rem;\">PostTemplateSyntaxError <em class=\"small\">class</em></h2>\nCustom error
  message for post template syntax errors.\n<div class=\"admonition source is-collapsible
  collapsible-open\">\n<p class=\"admonition-title\">\"PostTemplateSyntaxError <em
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
  <span class=\"nc\">PostTemplateSyntaxError</span><span class=\"p\">(</span><span
  class=\"n\">TemplateSyntaxError</span><span class=\"p\">):</span>\n<span class=\"w\">
  \           </span><span class=\"sd\">\"\"\"</span>\n<span class=\"sd\">            Custom
  error message for post template syntax errors.</span>\n<span class=\"sd\">            \"\"\"</span>\n</pre></div>\n\n</pre>\n<p>!!
  class </p><h2 class=\"admonition-title\" id=\"JinjaMd\" style=\"margin:0;padding:.5rem
  1rem;\">JinjaMd <em class=\"small\">class</em></h2>\n<div class=\"admonition source
  is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"JinjaMd <em
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
  <span class=\"nc\">JinjaMd</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
  class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
  \           <span class=\"n\">jinja</span><span class=\"p\">:</span> <span class=\"nb\">bool</span>
  <span class=\"o\">=</span> <span class=\"kc\">True</span>\n</pre></div>\n\n</pre>\n<p>!!
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
  class=\"n\">JinjaMd</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"pre_render\" style=\"margin:0;padding:.5rem
  1rem;\">pre_render <em class=\"small\">function</em></h2>\njinja_md hook for markata
  to render your markdown post as a jinja template.\n<pre><code>The post itself is
  exposed as `post`, and the markata instance is exposed\nas `markata`.\n</code></pre>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"pre_render
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
  <span class=\"nf\">pre_render</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
  class=\"w\">            </span><span class=\"sd\">\"\"\"</span>\n<span class=\"sd\">
  \           jinja_md hook for markata to render your markdown post as a jinja template.</span>\n\n<span
  class=\"sd\">            The post itself is exposed as `post`, and the markata instance
  is exposed</span>\n<span class=\"sd\">            as `markata`.</span>\n<span class=\"sd\">
  \           \"\"\"</span>\n\n            <span class=\"n\">config</span> <span class=\"o\">=</span>
  <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">jinja_md</span>\n            <span class=\"n\">ignore_spec</span>
  <span class=\"o\">=</span> <span class=\"n\">pathspec</span><span class=\"o\">.</span><span
  class=\"n\">PathSpec</span><span class=\"o\">.</span><span class=\"n\">from_lines</span><span
  class=\"p\">(</span><span class=\"s2\">\"gitwildmatch\"</span><span class=\"p\">,</span>
  <span class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">ignore</span><span
  class=\"p\">)</span>\n            <span class=\"c1\"># for article in markata.iter_articles(description=\"jinja_md\"):</span>\n
  \           <span class=\"n\">jinja_env</span> <span class=\"o\">=</span> <span
  class=\"n\">jinja2</span><span class=\"o\">.</span><span class=\"n\">Environment</span><span
  class=\"p\">(</span>\n                <span class=\"n\">extensions</span><span class=\"o\">=</span><span
  class=\"p\">[</span><span class=\"n\">IncludeRawExtension</span><span class=\"p\">,</span>
  <span class=\"o\">*</span><span class=\"n\">register_jinja_extensions</span><span
  class=\"p\">(</span><span class=\"n\">config</span><span class=\"p\">)],</span>\n
  \           <span class=\"p\">)</span>\n\n            <span class=\"k\">for</span>
  <span class=\"n\">article</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">articles</span><span class=\"p\">:</span>\n
  \               <span class=\"k\">if</span> <span class=\"n\">article</span><span
  class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
  class=\"s2\">\"jinja\"</span><span class=\"p\">,</span> <span class=\"kc\">True</span><span
  class=\"p\">)</span> <span class=\"ow\">and</span> <span class=\"ow\">not</span>
  <span class=\"n\">ignore_spec</span><span class=\"o\">.</span><span class=\"n\">match_file</span><span
  class=\"p\">(</span><span class=\"n\">article</span><span class=\"p\">[</span><span
  class=\"s2\">\"path\"</span><span class=\"p\">]):</span>\n                    <span
  class=\"k\">try</span><span class=\"p\">:</span>\n                        <span
  class=\"n\">key</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">make_hash</span><span class=\"p\">(</span><span
  class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">content</span><span
  class=\"p\">)</span>\n                        <span class=\"n\">content_from_cache</span>
  <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">precache</span><span class=\"o\">.</span><span class=\"n\">get</span><span
  class=\"p\">(</span><span class=\"n\">key</span><span class=\"p\">)</span>\n                        <span
  class=\"k\">if</span> <span class=\"n\">content_from_cache</span> <span class=\"ow\">is</span>
  <span class=\"kc\">None</span><span class=\"p\">:</span>\n                            <span
  class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">content</span>
  <span class=\"o\">=</span> <span class=\"n\">jinja_env</span><span class=\"o\">.</span><span
  class=\"n\">from_string</span><span class=\"p\">(</span><span class=\"n\">article</span><span
  class=\"o\">.</span><span class=\"n\">content</span><span class=\"p\">)</span><span
  class=\"o\">.</span><span class=\"n\">render</span><span class=\"p\">(</span>\n
  \                               <span class=\"n\">__version__</span><span class=\"o\">=</span><span
  class=\"n\">__version__</span><span class=\"p\">,</span>\n                                <span
  class=\"o\">**</span><span class=\"n\">article</span><span class=\"p\">,</span>\n
  \                               <span class=\"n\">post</span><span class=\"o\">=</span><span
  class=\"n\">article</span><span class=\"p\">,</span>\n                            <span
  class=\"p\">)</span>\n                            <span class=\"k\">with</span>
  <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">cache</span><span
  class=\"p\">:</span>\n                                <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">cache</span><span class=\"o\">.</span><span
  class=\"n\">set</span><span class=\"p\">(</span><span class=\"n\">key</span><span
  class=\"p\">,</span> <span class=\"n\">article</span><span class=\"o\">.</span><span
  class=\"n\">content</span><span class=\"p\">)</span>\n                        <span
  class=\"k\">else</span><span class=\"p\">:</span>\n                            <span
  class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">content</span>
  <span class=\"o\">=</span> <span class=\"n\">content_from_cache</span>\n                        <span
  class=\"c1\"># prevent double rendering</span>\n                        <span class=\"n\">article</span><span
  class=\"o\">.</span><span class=\"n\">jinja</span> <span class=\"o\">=</span> <span
  class=\"kc\">False</span>\n                    <span class=\"k\">except</span> <span
  class=\"n\">TemplateSyntaxError</span> <span class=\"k\">as</span> <span class=\"n\">e</span><span
  class=\"p\">:</span>\n                        <span class=\"n\">errorline</span>
  <span class=\"o\">=</span> <span class=\"n\">article</span><span class=\"o\">.</span><span
  class=\"n\">content</span><span class=\"o\">.</span><span class=\"n\">split</span><span
  class=\"p\">(</span><span class=\"s2\">\"</span><span class=\"se\">\\n</span><span
  class=\"s2\">\"</span><span class=\"p\">)[</span><span class=\"n\">e</span><span
  class=\"o\">.</span><span class=\"n\">lineno</span> <span class=\"o\">-</span> <span
  class=\"mi\">1</span><span class=\"p\">]</span>\n                        <span class=\"n\">msg</span>
  <span class=\"o\">=</span> <span class=\"sa\">f</span><span class=\"s2\">\"\"\"</span>\n<span
  class=\"s2\">                        Error while processing post </span><span class=\"si\">{</span><span
  class=\"n\">article</span><span class=\"p\">[</span><span class=\"s1\">'path'</span><span
  class=\"p\">]</span><span class=\"si\">}</span>\n\n<span class=\"s2\">                        </span><span
  class=\"si\">{</span><span class=\"n\">errorline</span><span class=\"si\">}</span>\n<span
  class=\"s2\">                        \"\"\"</span>\n\n                        <span
  class=\"k\">raise</span> <span class=\"n\">PostTemplateSyntaxError</span><span class=\"p\">(</span><span
  class=\"n\">msg</span><span class=\"p\">,</span> <span class=\"n\">lineno</span><span
  class=\"o\">=</span><span class=\"n\">e</span><span class=\"o\">.</span><span class=\"n\">lineno</span><span
  class=\"p\">)</span>\n                    <span class=\"k\">except</span> <span
  class=\"n\">UndefinedError</span> <span class=\"k\">as</span> <span class=\"n\">e</span><span
  class=\"p\">:</span>\n                        <span class=\"k\">raise</span> <span
  class=\"n\">UndefinedError</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
  class=\"s1\">'</span><span class=\"si\">{</span><span class=\"n\">e</span><span
  class=\"si\">}</span><span class=\"s1\"> in </span><span class=\"si\">{</span><span
  class=\"n\">article</span><span class=\"p\">[</span><span class=\"s2\">\"path\"</span><span
  class=\"p\">]</span><span class=\"si\">}</span><span class=\"s1\">'</span><span
  class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!! class </p><h2 class=\"admonition-title\"
  id=\"JinjaMdConfig\" style=\"margin:0;padding:.5rem 1rem;\">JinjaMdConfig <em class=\"small\">class</em></h2>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"JinjaMdConfig
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
  <span class=\"nc\">JinjaMdConfig</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
  class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
  \           <span class=\"n\">ignore</span><span class=\"p\">:</span> <span class=\"n\">List</span><span
  class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">]</span> <span
  class=\"o\">=</span> <span class=\"p\">[]</span>\n</pre></div>\n\n</pre>\n<p>!!
  class </p><h2 class=\"admonition-title\" id=\"Config\" style=\"margin:0;padding:.5rem
  1rem;\">Config <em class=\"small\">class</em></h2>\n<div class=\"admonition source
  is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"Config <em class=\"small\">source</em>\"</p>\n</div>\n<pre
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
  \           <span class=\"n\">jinja_md</span><span class=\"p\">:</span> <span class=\"n\">JinjaMdConfig</span>
  <span class=\"o\">=</span> <span class=\"n\">JinjaMdConfig</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n<p>!!
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
  method </p><h2 class=\"admonition-title\" id=\"parse\" style=\"margin:0;padding:.5rem
  1rem;\">parse <em class=\"small\">method</em></h2>\n<div class=\"admonition source
  is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"parse <em class=\"small\">source</em>\"</p>\n</div>\n<pre
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
  <span class=\"nf\">parse</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">,</span> <span class=\"n\">parser</span><span class=\"p\">):</span>\n
  \               <span class=\"n\">line_number</span> <span class=\"o\">=</span>
  <span class=\"nb\">next</span><span class=\"p\">(</span><span class=\"n\">parser</span><span
  class=\"o\">.</span><span class=\"n\">stream</span><span class=\"p\">)</span><span
  class=\"o\">.</span><span class=\"n\">lineno</span>\n                <span class=\"n\">file</span>
  <span class=\"o\">=</span> <span class=\"p\">[</span><span class=\"n\">parser</span><span
  class=\"o\">.</span><span class=\"n\">parse_expression</span><span class=\"p\">()]</span>\n
  \               <span class=\"k\">return</span> <span class=\"n\">nodes</span><span
  class=\"o\">.</span><span class=\"n\">CallBlock</span><span class=\"p\">(</span>\n
  \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">call_method</span><span class=\"p\">(</span><span class=\"s2\">\"_read_file\"</span><span
  class=\"p\">,</span> <span class=\"n\">file</span><span class=\"p\">),</span>\n
  \                   <span class=\"p\">[],</span>\n                    <span class=\"p\">[],</span>\n
  \                   <span class=\"s2\">\"\"</span><span class=\"p\">,</span>\n                <span
  class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">set_lineno</span><span
  class=\"p\">(</span><span class=\"n\">line_number</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"_read_file\" style=\"margin:0;padding:.5rem
  1rem;\">_read_file <em class=\"small\">method</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"_read_file
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
  <span class=\"nf\">_read_file</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">,</span> <span class=\"n\">file</span><span class=\"p\">,</span> <span
  class=\"n\">caller</span><span class=\"p\">):</span>\n                <span class=\"k\">return</span>
  <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">file</span><span
  class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">read_text</span><span
  class=\"p\">()</span>\n</pre></div>\n\n</pre>\n<p>!! method </p><h2 class=\"admonition-title\"
  id=\"_fail_with_undefined_error\" style=\"margin:0;padding:.5rem 1rem;\">_fail_with_undefined_error
  <em class=\"small\">method</em></h2>\n<div class=\"admonition source is-collapsible
  collapsible-open\">\n<p class=\"admonition-title\">\"_fail_with_undefined_error
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
  <span class=\"nf\">_fail_with_undefined_error</span><span class=\"p\">(</span><span
  class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
  class=\"n\">args</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
  class=\"n\">kwargs</span><span class=\"p\">):</span>\n                <span class=\"k\">return</span>
  <span class=\"s2\">\"\"</span>\n</pre></div>\n\n</pre>\n</main>\n<footer>\xA9 2024</footer>\n</body></html>"
published: true
slug: markata/plugins/jinja-md
title: Jinja_Md.Py


---

Renders your markdown as a jinja template during pre_render.

The markata instance is passed into the template, giving you access to things
such as all of your articles, config, and this post as post.

# Examples

first we can grab a few things out of the frontmatter of this post.

``` markdown
# {{ post.title }}
{{ post.description }}
```

# one-liner list of links

This one-liner will render a list of markdown links into your markdown at build
time.  It's quite handy to pop into posts.

``` markdown
{{ '\n'.join(markata.map('f"* [{title}]({slug})"', sort='slug')) }}
```

# jinja for to markdown list of links

Sometimes quoting things like your filters are hard to do in a one line without
running out of quote variants.  Jinja for loops can make this much easier.

``` markdown
{% for post in markata.map('post', filter='"git" in tags') %}
* [{{ post.title }}]({{ post.slug }})
{% endfor %}
```

# jinja for to html list of links

Since markdown is a superset of html, you can just render out html into your
post and it is still valid.

``` markdown
<ul>
{% for post in markata.map('post', filter='"git" in tags') %}
    <li><a href="{{ post.slug }}">{{ post.title }}</a></li>
{% endfor %}
</ul>
```

# Ignoring files

It is possible to ignore files by adding an ignore to your `markata.jinja_md`
config in your `markata.toml` file.  This ignore follows the `gitwildmatch`
rules, so think of it the same as writing a gitignore.

```toml
[markata.jinja_md]
ignore=[
'jinja_md.md',
]
```

!!note
  Docs such as this jinja_md.py file will get converted to jinja_md.md during
  build time, so use `.md` extensions instead of `.py`.

# Ignoring a single file

You can also ignore a single file right from the articles frontmatter, by
adding `jinja: false`.

```markdown
---
jinja: false

---
```

# Escaping

Sometimes you want the ability to have jinja templates in a post, but also the
ability to keep a raw jinja template.  There are a couple of techniques that
are covered mroe in the jinja docs for
[escaping](https://jinja.palletsprojects.com/en/3.1.x/templates/#escaping)

```markdown
{% raw %}
{{ '\n'.join(markata.map('f"* [{title}]({slug})"', sort='slug')) }}
{% endraw %}

{{ '{{' }} '\n'.join(markata.map('f"* [{title}]({slug})"', sort='slug')) {{ '}}' }}
```

# Creating a jinja extension

Here is a bit of a boilerplate example of a jinja extension.

``` python
from jinja2 import nodes
from jinja2.ext import Extension

class ExampleExtension(Extension):
    tags = {"example"}

    def __init__(self, environment):
        super().__init__(environment)

    def parse(self, parser):
        line_number = next(parser.stream).lineno
        arg = [parser.parse_expression()]
        return nodes.CallBlock(self.call_method("run", arg), [], [], "").set_lineno(
            line_number
        )

    def run(self, arg, caller):
        return f'hello {arg}'
```

So that markata picks up your extension, you will need to register an
entrypoint named `markata.jinja_md`.  Once installed markata will automatically
load this extension to its list of  jinja extensions.

``` toml
[project.entry-points."markata.jinja_md"]
markta_gh = "example_extension:ExampleExtension"
```

Once you have your extension created and ready to use you can use it in your
markdown.

``` markdown
{% example world %}
```


!! function <h2 id='register_jinja_extensions' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>register_jinja_extensions <em class='small'>function</em></h2>
    Gets jinja extensions from entrypoints and loads them in.

    Returns: List of jinja Extensions
???+ source "register_jinja_extensions <em class='small'>source</em>"

```python

        def register_jinja_extensions(config: dict) -> List[Extension]:
            """
            Gets jinja extensions from entrypoints and loads them in.

            Returns: List of jinja Extensions
            """

            return [
                ep.load() for ep in pkg_resources.iter_entry_points(group="markata.jinja_md")
            ]
```


!! class <h2 id='IncludeRawExtension' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>IncludeRawExtension <em class='small'>class</em></h2>

???+ source "IncludeRawExtension <em class='small'>source</em>"

```python

        class IncludeRawExtension(Extension):
            tags = {"include_raw"}

            def parse(self, parser):
                line_number = next(parser.stream).lineno
                file = [parser.parse_expression()]
                return nodes.CallBlock(
                    self.call_method("_read_file", file),
                    [],
                    [],
                    "",
                ).set_lineno(line_number)

            def _read_file(self, file, caller):
                return Path(file).read_text()
```


!! class <h2 id='_SilentUndefined' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_SilentUndefined <em class='small'>class</em></h2>
    silence undefined variable errors in jinja templates.

    # Example
    ```python
    template = '{{ variable }}'
    article.content = Template( template, undefined=_SilentUndefined).render()
    ```
???+ source "_SilentUndefined <em class='small'>source</em>"

```python

        class _SilentUndefined(Undefined):
            """
            silence undefined variable errors in jinja templates.

            # Example
            ```python
            template = '{{ variable }}'
            article.content = Template( template, undefined=_SilentUndefined).render()
            ```
            """

            def _fail_with_undefined_error(self, *args, **kwargs):
                return ""
```


!! class <h2 id='PostTemplateSyntaxError' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>PostTemplateSyntaxError <em class='small'>class</em></h2>
    Custom error message for post template syntax errors.
???+ source "PostTemplateSyntaxError <em class='small'>source</em>"

```python

        class PostTemplateSyntaxError(TemplateSyntaxError):
            """
            Custom error message for post template syntax errors.
            """
```


!! class <h2 id='JinjaMd' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>JinjaMd <em class='small'>class</em></h2>

???+ source "JinjaMd <em class='small'>source</em>"

```python

        class JinjaMd(pydantic.BaseModel):
            jinja: bool = True
```


!! function <h2 id='post_model' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>post_model <em class='small'>function</em></h2>

???+ source "post_model <em class='small'>source</em>"

```python

        def post_model(markata: "Markata") -> None:
            markata.post_models.append(JinjaMd)
```


!! function <h2 id='pre_render' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>pre_render <em class='small'>function</em></h2>
    jinja_md hook for markata to render your markdown post as a jinja template.

    The post itself is exposed as `post`, and the markata instance is exposed
    as `markata`.
???+ source "pre_render <em class='small'>source</em>"

```python

        def pre_render(markata: "Markata") -> None:
            """
            jinja_md hook for markata to render your markdown post as a jinja template.

            The post itself is exposed as `post`, and the markata instance is exposed
            as `markata`.
            """

            config = markata.config.jinja_md
            ignore_spec = pathspec.PathSpec.from_lines("gitwildmatch", config.ignore)
            # for article in markata.iter_articles(description="jinja_md"):
            jinja_env = jinja2.Environment(
                extensions=[IncludeRawExtension, *register_jinja_extensions(config)],
            )

            for article in markata.articles:
                if article.get("jinja", True) and not ignore_spec.match_file(article["path"]):
                    try:
                        key = markata.make_hash(article.content)
                        content_from_cache = markata.precache.get(key)
                        if content_from_cache is None:
                            article.content = jinja_env.from_string(article.content).render(
                                __version__=__version__,
                                **article,
                                post=article,
                            )
                            with markata.cache:
                                markata.cache.set(key, article.content)
                        else:
                            article.content = content_from_cache
                        # prevent double rendering
                        article.jinja = False
                    except TemplateSyntaxError as e:
                        errorline = article.content.split("\n")[e.lineno - 1]
                        msg = f"""
                        Error while processing post {article['path']}

                        {errorline}
                        """

                        raise PostTemplateSyntaxError(msg, lineno=e.lineno)
                    except UndefinedError as e:
                        raise UndefinedError(f'{e} in {article["path"]}')
```


!! class <h2 id='JinjaMdConfig' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>JinjaMdConfig <em class='small'>class</em></h2>

???+ source "JinjaMdConfig <em class='small'>source</em>"

```python

        class JinjaMdConfig(pydantic.BaseModel):
            ignore: List[str] = []
```


!! class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Config <em class='small'>class</em></h2>

???+ source "Config <em class='small'>source</em>"

```python

        class Config(pydantic.BaseModel):
            jinja_md: JinjaMdConfig = JinjaMdConfig()
```


!! function <h2 id='config_model' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>config_model <em class='small'>function</em></h2>

???+ source "config_model <em class='small'>source</em>"

```python

        def config_model(markata: "Markata") -> None:
            markata.config_models.append(Config)
```


!! method <h2 id='parse' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>parse <em class='small'>method</em></h2>

???+ source "parse <em class='small'>source</em>"

```python

        def parse(self, parser):
                line_number = next(parser.stream).lineno
                file = [parser.parse_expression()]
                return nodes.CallBlock(
                    self.call_method("_read_file", file),
                    [],
                    [],
                    "",
                ).set_lineno(line_number)
```


!! method <h2 id='_read_file' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_read_file <em class='small'>method</em></h2>

???+ source "_read_file <em class='small'>source</em>"

```python

        def _read_file(self, file, caller):
                return Path(file).read_text()
```


!! method <h2 id='_fail_with_undefined_error' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_fail_with_undefined_error <em class='small'>method</em></h2>

???+ source "_fail_with_undefined_error <em class='small'>source</em>"

```python

        def _fail_with_undefined_error(self, *args, **kwargs):
                return ""
```

