---
content: "leading docstring\n\n\n!! function <h2 id='add_parents' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>add_parents <em class='small'>function</em></h2>\n\n???+
  source \"add_parents <em class='small'>source</em>\"\n\n```python\n\n        def
  add_parents(tree: ast.AST) -> None:\n            for node in ast.walk(tree):\n                for
  child in ast.iter_child_nodes(node):\n                    child.parent = node\n
  \                   if not hasattr(child, \"parents\"):\n                        child.parents
  = [node]\n                    child.parents.append(node)\n                    if
  isinstance(node, ast.ClassDef) and isinstance(child, ast.FunctionDef):\n                        child.type
  = \"method\"\n                    elif isinstance(child, ast.FunctionDef):\n                        child.type
  = \"function\"\n                    elif isinstance(child, ast.ClassDef):\n                        child.type
  = \"class\"\n```\n\n\n!! function <h2 id='glob' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>glob <em class='small'>function</em></h2>\n    finds k\n\n    ## Parameters\n\n
  \   `markata` the markata object\n???+ source \"glob <em class='small'>source</em>\"\n\n```python\n\n
  \       def glob(markata: \"MarkataDocs\") -> None:\n            \"\"\"\n            finds
  k\n\n            ## Parameters\n\n            `markata` the markata object\n\n            \"\"\"\n\n
  \           import glob\n\n            markata.py_files = [Path(f) for f in glob.glob(\"**/*.py\",
  recursive=True)]\n\n            content_directories = list({f.parent for f in markata.py_files})\n
  \           if \"content_directories\" in markata.__dict__:\n                markata.content_directories.extend(content_directories)\n
  \           else:\n                markata.content_directories = content_directories\n\n
  \           try:\n                ignore = True\n            except KeyError:\n
  \               ignore = True\n\n            if ignore and (Path(\".gitignore\").exists()
  or Path(\".markataignore\").exists()):\n                import pathspec\n\n                lines
  = []\n\n                if Path(\".gitignore\").exists():\n                    lines.extend(Path(\".gitignore\").read_text().splitlines())\n\n
  \               if Path(\".markataignore\").exists():\n                    lines.extend(Path(\".markataignore\").read_text().splitlines())\n\n
  \           spec = pathspec.PathSpec.from_lines(\"gitwildmatch\", lines)\n\n            markata.py_files
  = [\n                file for file in markata.py_files if not spec.match_file(str(file))\n
  \           ]\n```\n\n\n!! function <h2 id='get_template' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>get_template <em class='small'>function</em></h2>\n\n???+
  source \"get_template <em class='small'>source</em>\"\n\n```python\n\n        def
  get_template():\n            jinja_env = jinja2.Environment()\n            template
  = jinja_env.from_string(\n                (Path(__file__).parent / \"default_doc_template.md\").read_text(),\n
  \           )\n            return template\n```\n\n\n!! function <h2 id='make_article'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>make_article <em class='small'>function</em></h2>\n\n???+
  source \"make_article <em class='small'>source</em>\"\n\n```python\n\n        def
  make_article(markata: \"Markata\", file: Path, cache) -> frontmatter.Post:\n            with
  open(file) as f:\n                raw_source = f.read()\n            key = markata.make_hash(\"docs\",
  \"file\", raw_source)\n            slug = f\"{file.parent}/{file.stem}\".lstrip(\"/\").lstrip(\"./\")\n
  \           edit_link = (\n                str(markata.config.get(\"repo_url\",
  \"https://github.com/\"))\n                + \"edit/\"\n                + str(markata.config.get(\"repo_branch\",
  \"main\"))\n                + \"/\"\n                + str(file)\n            )\n
  \           article_from_cache = markata.precache.get(key)\n            if article_from_cache
  is not None:\n                article = article_from_cache\n            else:\n
  \               tree = ast.parse(raw_source)\n                add_parents(tree)\n
  \               nodes = [\n                    n for n in ast.walk(tree) if isinstance(n,
  (ast.FunctionDef, ast.ClassDef))\n                ]\n\n                article =
  get_template().render(\n                    ast=ast,\n                    file=file,\n
  \                   slug=slug,\n                    edit_link=edit_link,\n                    tree=tree,\n
  \                   datetime=datetime,\n                    nodes=nodes,\n                    raw_source=raw_source,\n
  \                   indent=textwrap.indent,\n                )\n                cache.add(\n
  \                   key,\n                    article,\n                    expire=markata.config.default_cache_expire,\n
  \               )\n\n            try:\n                article = markata.Post(\n
  \                   markata=markata,\n                    path=str(file).replace(\".py\",
  \".md\"),\n                    title=file.name,\n                    content=article,\n
  \                   ast=ast,\n                    file=file,\n                    slug=slug,\n
  \                   edit_link=edit_link,\n                    datetime=datetime,\n
  \               )\n\n            except pydantic.ValidationError as e:\n                from
  markata.plugins.load import ValidationError, get_models\n\n                models
  = get_models(markata=markata, error=e)\n                models = list(models.values())\n
  \               models = \"\\n\".join(models)\n                raise ValidationError(f\"{e}\\n\\n{models}\\nfailed
  to load {path}\") from e\n\n            return article\n```\n\n\n!! function <h2
  id='load' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>load <em
  class='small'>function</em></h2>\n    similar to [glob](../glob)\n???+ source \"load
  <em class='small'>source</em>\"\n\n```python\n\n        def load(markata: \"MarkataDocs\")
  -> None:\n            \"\"\"\n            similar to [glob](../glob)\n            \"\"\"\n
  \           if \"articles\" not in markata.__dict__:\n                markata.articles
  = []\n            for py_file in markata.py_files:\n                with markata.cache
  as cache:\n                    markata.articles.append(make_article(markata, py_file,
  cache))\n```\n\n\n!! class <h2 id='MarkataDocs' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>MarkataDocs <em class='small'>class</em></h2>\n\n???+ source \"MarkataDocs
  <em class='small'>source</em>\"\n\n```python\n\n        class MarkataDocs(Markata):\n
  \               py_files: List = []\n                content_directories: List =
  []\n```\n"
date: 0001-01-01
description: 'leading docstring ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
  source  ! ! ???+ source '
html: "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n<title>Docs.Py</title>\n<meta
  charset=\"utf-8\"/>\n<meta content=\"width=device-width, initial-scale=1\" name=\"viewport\"/>\n<meta
  content=\"leading docstring ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
  source  ! ! ???+ source \" name=\"description\"/>\n<link href=\"/static/favicon.ico\"
  rel=\"icon\" type=\"image/png\"/>\n<script>\n        function setTheme(theme) {\n
  \           document.documentElement.setAttribute(\"data-theme\", theme);\n        }\n\n
  \       function detectColorSchemeOnLoad() {\n            //local storage is used
  to override OS theme settings\n            if (localStorage.getItem(\"theme\"))
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
  name=\"og:type\" property=\"og:type\"/><meta content=\"leading docstring ! ???+
  source  ! ???+ source  ! ???+ source  ! ???+ source  ! ! ???+ source \" name=\"description\"
  property=\"description\"/><meta content=\"leading docstring ! ???+ source  ! ???+
  source  ! ???+ source  ! ???+ source  ! ! ???+ source \" name=\"og:description\"
  property=\"og:description\"/><meta content=\"leading docstring ! ???+ source  !
  ???+ source  ! ???+ source  ! ???+ source  ! ! ???+ source \" name=\"twitter:description\"
  property=\"twitter:description\"/><meta content=\"Docs.Py | Markata\" name=\"og:title\"
  property=\"og:title\"/><meta content=\"Docs.Py | Markata\" name=\"twitter:title\"
  property=\"twitter:title\"/><meta content=\"https://markata.dev//markata/plugins/docs-og.png\"
  name=\"og:image\" property=\"og:image\"/><meta content=\"https://markata.dev//markata/plugins/docs-og.png\"
  name=\"twitter:image\" property=\"twitter:image\"/><meta content=\"1600\" name=\"og:image:width\"
  property=\"og:image:width\"/><meta content=\"900\" name=\"og:image:width\" property=\"og:image:width\"/><meta
  content=\"summary_large_image\" name=\"twitter:card\" property=\"twitter:card\"/><meta
  content=\"Markata\" name=\"og:site_name\" property=\"og:site_name\"/><meta content=\"@_waylonwalker\"
  name=\"twitter:creator\" property=\"twitter:creator\"/><meta content=\"Docs.Py\"
  name=\"title\" property=\"title\"/><meta content=\"markata 0.8.0.dev6\" name=\"generator\"
  property=\"generator\"/><link href=\"https://markata.dev//markata/plugins/docs/\"
  rel=\"canonical\"/><meta content=\"https://markata.dev//markata/plugins/docs/\"
  name=\"og:url\" property=\"og:url\"/></head>\n<body><nav>\n<a href=\"/\">markata</a>\n<a
  href=\"https://github.com/WaylonWalker/markata\">GitHub</a>\n</nav>\n<div>\n<label
  class=\"theme-switch\" for=\"checkbox-theme\" id=\"theme-switch\" title=\"light/dark
  mode toggle\">\n<input id=\"checkbox-theme\" type=\"checkbox\"/>\n<div class=\"slider
  round\"></div>\n</label>\n</div>\n<section class=\"title\">\n<h1 id=\"title\">\n
  \           Docs.Py \n            \n        </h1>\n</section>\n<main><p>leading
  docstring</p>\n<p>!! function </p><h2 class=\"admonition-title\" id=\"add_parents\"
  style=\"margin:0;padding:.5rem 1rem;\">add_parents <em class=\"small\">function</em></h2>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"add_parents
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
  <span class=\"nf\">add_parents</span><span class=\"p\">(</span><span class=\"n\">tree</span><span
  class=\"p\">:</span> <span class=\"n\">ast</span><span class=\"o\">.</span><span
  class=\"n\">AST</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"k\">for</span>
  <span class=\"n\">node</span> <span class=\"ow\">in</span> <span class=\"n\">ast</span><span
  class=\"o\">.</span><span class=\"n\">walk</span><span class=\"p\">(</span><span
  class=\"n\">tree</span><span class=\"p\">):</span>\n                <span class=\"k\">for</span>
  <span class=\"n\">child</span> <span class=\"ow\">in</span> <span class=\"n\">ast</span><span
  class=\"o\">.</span><span class=\"n\">iter_child_nodes</span><span class=\"p\">(</span><span
  class=\"n\">node</span><span class=\"p\">):</span>\n                    <span class=\"n\">child</span><span
  class=\"o\">.</span><span class=\"n\">parent</span> <span class=\"o\">=</span> <span
  class=\"n\">node</span>\n                    <span class=\"k\">if</span> <span class=\"ow\">not</span>
  <span class=\"nb\">hasattr</span><span class=\"p\">(</span><span class=\"n\">child</span><span
  class=\"p\">,</span> <span class=\"s2\">\"parents\"</span><span class=\"p\">):</span>\n
  \                       <span class=\"n\">child</span><span class=\"o\">.</span><span
  class=\"n\">parents</span> <span class=\"o\">=</span> <span class=\"p\">[</span><span
  class=\"n\">node</span><span class=\"p\">]</span>\n                    <span class=\"n\">child</span><span
  class=\"o\">.</span><span class=\"n\">parents</span><span class=\"o\">.</span><span
  class=\"n\">append</span><span class=\"p\">(</span><span class=\"n\">node</span><span
  class=\"p\">)</span>\n                    <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
  class=\"p\">(</span><span class=\"n\">node</span><span class=\"p\">,</span> <span
  class=\"n\">ast</span><span class=\"o\">.</span><span class=\"n\">ClassDef</span><span
  class=\"p\">)</span> <span class=\"ow\">and</span> <span class=\"nb\">isinstance</span><span
  class=\"p\">(</span><span class=\"n\">child</span><span class=\"p\">,</span> <span
  class=\"n\">ast</span><span class=\"o\">.</span><span class=\"n\">FunctionDef</span><span
  class=\"p\">):</span>\n                        <span class=\"n\">child</span><span
  class=\"o\">.</span><span class=\"n\">type</span> <span class=\"o\">=</span> <span
  class=\"s2\">\"method\"</span>\n                    <span class=\"k\">elif</span>
  <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">child</span><span
  class=\"p\">,</span> <span class=\"n\">ast</span><span class=\"o\">.</span><span
  class=\"n\">FunctionDef</span><span class=\"p\">):</span>\n                        <span
  class=\"n\">child</span><span class=\"o\">.</span><span class=\"n\">type</span>
  <span class=\"o\">=</span> <span class=\"s2\">\"function\"</span>\n                    <span
  class=\"k\">elif</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
  class=\"n\">child</span><span class=\"p\">,</span> <span class=\"n\">ast</span><span
  class=\"o\">.</span><span class=\"n\">ClassDef</span><span class=\"p\">):</span>\n
  \                       <span class=\"n\">child</span><span class=\"o\">.</span><span
  class=\"n\">type</span> <span class=\"o\">=</span> <span class=\"s2\">\"class\"</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"glob\" style=\"margin:0;padding:.5rem
  1rem;\">glob <em class=\"small\">function</em></h2>\nfinds k\n<pre><code>## Parameters\n\n`markata`
  the markata object\n</code></pre>\n<div class=\"admonition source is-collapsible
  collapsible-open\">\n<p class=\"admonition-title\">\"glob <em class=\"small\">source</em>\"</p>\n</div>\n<pre
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
  <span class=\"nf\">glob</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
  class=\"p\">:</span> <span class=\"s2\">\"MarkataDocs\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
  class=\"w\">            </span><span class=\"sd\">\"\"\"</span>\n<span class=\"sd\">
  \           finds k</span>\n\n<span class=\"sd\">            ## Parameters</span>\n\n<span
  class=\"sd\">            `markata` the markata object</span>\n\n<span class=\"sd\">
  \           \"\"\"</span>\n\n            <span class=\"kn\">import</span> <span
  class=\"nn\">glob</span>\n\n            <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">py_files</span> <span class=\"o\">=</span> <span class=\"p\">[</span><span
  class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">f</span><span
  class=\"p\">)</span> <span class=\"k\">for</span> <span class=\"n\">f</span> <span
  class=\"ow\">in</span> <span class=\"n\">glob</span><span class=\"o\">.</span><span
  class=\"n\">glob</span><span class=\"p\">(</span><span class=\"s2\">\"**/*.py\"</span><span
  class=\"p\">,</span> <span class=\"n\">recursive</span><span class=\"o\">=</span><span
  class=\"kc\">True</span><span class=\"p\">)]</span>\n\n            <span class=\"n\">content_directories</span>
  <span class=\"o\">=</span> <span class=\"nb\">list</span><span class=\"p\">({</span><span
  class=\"n\">f</span><span class=\"o\">.</span><span class=\"n\">parent</span> <span
  class=\"k\">for</span> <span class=\"n\">f</span> <span class=\"ow\">in</span> <span
  class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">py_files</span><span
  class=\"p\">})</span>\n            <span class=\"k\">if</span> <span class=\"s2\">\"content_directories\"</span>
  <span class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"vm\">__dict__</span><span class=\"p\">:</span>\n                <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">content_directories</span><span class=\"o\">.</span><span
  class=\"n\">extend</span><span class=\"p\">(</span><span class=\"n\">content_directories</span><span
  class=\"p\">)</span>\n            <span class=\"k\">else</span><span class=\"p\">:</span>\n
  \               <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">content_directories</span> <span class=\"o\">=</span> <span class=\"n\">content_directories</span>\n\n
  \           <span class=\"k\">try</span><span class=\"p\">:</span>\n                <span
  class=\"n\">ignore</span> <span class=\"o\">=</span> <span class=\"kc\">True</span>\n
  \           <span class=\"k\">except</span> <span class=\"ne\">KeyError</span><span
  class=\"p\">:</span>\n                <span class=\"n\">ignore</span> <span class=\"o\">=</span>
  <span class=\"kc\">True</span>\n\n            <span class=\"k\">if</span> <span
  class=\"n\">ignore</span> <span class=\"ow\">and</span> <span class=\"p\">(</span><span
  class=\"n\">Path</span><span class=\"p\">(</span><span class=\"s2\">\".gitignore\"</span><span
  class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">exists</span><span
  class=\"p\">()</span> <span class=\"ow\">or</span> <span class=\"n\">Path</span><span
  class=\"p\">(</span><span class=\"s2\">\".markataignore\"</span><span class=\"p\">)</span><span
  class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">()):</span>\n
  \               <span class=\"kn\">import</span> <span class=\"nn\">pathspec</span>\n\n
  \               <span class=\"n\">lines</span> <span class=\"o\">=</span> <span
  class=\"p\">[]</span>\n\n                <span class=\"k\">if</span> <span class=\"n\">Path</span><span
  class=\"p\">(</span><span class=\"s2\">\".gitignore\"</span><span class=\"p\">)</span><span
  class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">():</span>\n
  \                   <span class=\"n\">lines</span><span class=\"o\">.</span><span
  class=\"n\">extend</span><span class=\"p\">(</span><span class=\"n\">Path</span><span
  class=\"p\">(</span><span class=\"s2\">\".gitignore\"</span><span class=\"p\">)</span><span
  class=\"o\">.</span><span class=\"n\">read_text</span><span class=\"p\">()</span><span
  class=\"o\">.</span><span class=\"n\">splitlines</span><span class=\"p\">())</span>\n\n
  \               <span class=\"k\">if</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
  class=\"s2\">\".markataignore\"</span><span class=\"p\">)</span><span class=\"o\">.</span><span
  class=\"n\">exists</span><span class=\"p\">():</span>\n                    <span
  class=\"n\">lines</span><span class=\"o\">.</span><span class=\"n\">extend</span><span
  class=\"p\">(</span><span class=\"n\">Path</span><span class=\"p\">(</span><span
  class=\"s2\">\".markataignore\"</span><span class=\"p\">)</span><span class=\"o\">.</span><span
  class=\"n\">read_text</span><span class=\"p\">()</span><span class=\"o\">.</span><span
  class=\"n\">splitlines</span><span class=\"p\">())</span>\n\n            <span class=\"n\">spec</span>
  <span class=\"o\">=</span> <span class=\"n\">pathspec</span><span class=\"o\">.</span><span
  class=\"n\">PathSpec</span><span class=\"o\">.</span><span class=\"n\">from_lines</span><span
  class=\"p\">(</span><span class=\"s2\">\"gitwildmatch\"</span><span class=\"p\">,</span>
  <span class=\"n\">lines</span><span class=\"p\">)</span>\n\n            <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">py_files</span> <span class=\"o\">=</span>
  <span class=\"p\">[</span>\n                <span class=\"n\">file</span> <span
  class=\"k\">for</span> <span class=\"n\">file</span> <span class=\"ow\">in</span>
  <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">py_files</span>
  <span class=\"k\">if</span> <span class=\"ow\">not</span> <span class=\"n\">spec</span><span
  class=\"o\">.</span><span class=\"n\">match_file</span><span class=\"p\">(</span><span
  class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">file</span><span
  class=\"p\">))</span>\n            <span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"get_template\" style=\"margin:0;padding:.5rem
  1rem;\">get_template <em class=\"small\">function</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"get_template
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
  <span class=\"nf\">get_template</span><span class=\"p\">():</span>\n            <span
  class=\"n\">jinja_env</span> <span class=\"o\">=</span> <span class=\"n\">jinja2</span><span
  class=\"o\">.</span><span class=\"n\">Environment</span><span class=\"p\">()</span>\n
  \           <span class=\"n\">template</span> <span class=\"o\">=</span> <span class=\"n\">jinja_env</span><span
  class=\"o\">.</span><span class=\"n\">from_string</span><span class=\"p\">(</span>\n
  \               <span class=\"p\">(</span><span class=\"n\">Path</span><span class=\"p\">(</span><span
  class=\"vm\">__file__</span><span class=\"p\">)</span><span class=\"o\">.</span><span
  class=\"n\">parent</span> <span class=\"o\">/</span> <span class=\"s2\">\"default_doc_template.md\"</span><span
  class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">read_text</span><span
  class=\"p\">(),</span>\n            <span class=\"p\">)</span>\n            <span
  class=\"k\">return</span> <span class=\"n\">template</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"make_article\" style=\"margin:0;padding:.5rem
  1rem;\">make_article <em class=\"small\">function</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"make_article
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
  <span class=\"nf\">make_article</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">,</span>
  <span class=\"n\">file</span><span class=\"p\">:</span> <span class=\"n\">Path</span><span
  class=\"p\">,</span> <span class=\"n\">cache</span><span class=\"p\">)</span> <span
  class=\"o\">-&gt;</span> <span class=\"n\">frontmatter</span><span class=\"o\">.</span><span
  class=\"n\">Post</span><span class=\"p\">:</span>\n            <span class=\"k\">with</span>
  <span class=\"nb\">open</span><span class=\"p\">(</span><span class=\"n\">file</span><span
  class=\"p\">)</span> <span class=\"k\">as</span> <span class=\"n\">f</span><span
  class=\"p\">:</span>\n                <span class=\"n\">raw_source</span> <span
  class=\"o\">=</span> <span class=\"n\">f</span><span class=\"o\">.</span><span class=\"n\">read</span><span
  class=\"p\">()</span>\n            <span class=\"n\">key</span> <span class=\"o\">=</span>
  <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">make_hash</span><span
  class=\"p\">(</span><span class=\"s2\">\"docs\"</span><span class=\"p\">,</span>
  <span class=\"s2\">\"file\"</span><span class=\"p\">,</span> <span class=\"n\">raw_source</span><span
  class=\"p\">)</span>\n            <span class=\"n\">slug</span> <span class=\"o\">=</span>
  <span class=\"sa\">f</span><span class=\"s2\">\"</span><span class=\"si\">{</span><span
  class=\"n\">file</span><span class=\"o\">.</span><span class=\"n\">parent</span><span
  class=\"si\">}</span><span class=\"s2\">/</span><span class=\"si\">{</span><span
  class=\"n\">file</span><span class=\"o\">.</span><span class=\"n\">stem</span><span
  class=\"si\">}</span><span class=\"s2\">\"</span><span class=\"o\">.</span><span
  class=\"n\">lstrip</span><span class=\"p\">(</span><span class=\"s2\">\"/\"</span><span
  class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">lstrip</span><span
  class=\"p\">(</span><span class=\"s2\">\"./\"</span><span class=\"p\">)</span>\n
  \           <span class=\"n\">edit_link</span> <span class=\"o\">=</span> <span
  class=\"p\">(</span>\n                <span class=\"nb\">str</span><span class=\"p\">(</span><span
  class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
  class=\"s2\">\"repo_url\"</span><span class=\"p\">,</span> <span class=\"s2\">\"https://github.com/\"</span><span
  class=\"p\">))</span>\n                <span class=\"o\">+</span> <span class=\"s2\">\"edit/\"</span>\n
  \               <span class=\"o\">+</span> <span class=\"nb\">str</span><span class=\"p\">(</span><span
  class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
  class=\"s2\">\"repo_branch\"</span><span class=\"p\">,</span> <span class=\"s2\">\"main\"</span><span
  class=\"p\">))</span>\n                <span class=\"o\">+</span> <span class=\"s2\">\"/\"</span>\n
  \               <span class=\"o\">+</span> <span class=\"nb\">str</span><span class=\"p\">(</span><span
  class=\"n\">file</span><span class=\"p\">)</span>\n            <span class=\"p\">)</span>\n
  \           <span class=\"n\">article_from_cache</span> <span class=\"o\">=</span>
  <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">precache</span><span
  class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
  class=\"n\">key</span><span class=\"p\">)</span>\n            <span class=\"k\">if</span>
  <span class=\"n\">article_from_cache</span> <span class=\"ow\">is</span> <span class=\"ow\">not</span>
  <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
  class=\"n\">article</span> <span class=\"o\">=</span> <span class=\"n\">article_from_cache</span>\n
  \           <span class=\"k\">else</span><span class=\"p\">:</span>\n                <span
  class=\"n\">tree</span> <span class=\"o\">=</span> <span class=\"n\">ast</span><span
  class=\"o\">.</span><span class=\"n\">parse</span><span class=\"p\">(</span><span
  class=\"n\">raw_source</span><span class=\"p\">)</span>\n                <span class=\"n\">add_parents</span><span
  class=\"p\">(</span><span class=\"n\">tree</span><span class=\"p\">)</span>\n                <span
  class=\"n\">nodes</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
  \                   <span class=\"n\">n</span> <span class=\"k\">for</span> <span
  class=\"n\">n</span> <span class=\"ow\">in</span> <span class=\"n\">ast</span><span
  class=\"o\">.</span><span class=\"n\">walk</span><span class=\"p\">(</span><span
  class=\"n\">tree</span><span class=\"p\">)</span> <span class=\"k\">if</span> <span
  class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">n</span><span
  class=\"p\">,</span> <span class=\"p\">(</span><span class=\"n\">ast</span><span
  class=\"o\">.</span><span class=\"n\">FunctionDef</span><span class=\"p\">,</span>
  <span class=\"n\">ast</span><span class=\"o\">.</span><span class=\"n\">ClassDef</span><span
  class=\"p\">))</span>\n                <span class=\"p\">]</span>\n\n                <span
  class=\"n\">article</span> <span class=\"o\">=</span> <span class=\"n\">get_template</span><span
  class=\"p\">()</span><span class=\"o\">.</span><span class=\"n\">render</span><span
  class=\"p\">(</span>\n                    <span class=\"n\">ast</span><span class=\"o\">=</span><span
  class=\"n\">ast</span><span class=\"p\">,</span>\n                    <span class=\"n\">file</span><span
  class=\"o\">=</span><span class=\"n\">file</span><span class=\"p\">,</span>\n                    <span
  class=\"n\">slug</span><span class=\"o\">=</span><span class=\"n\">slug</span><span
  class=\"p\">,</span>\n                    <span class=\"n\">edit_link</span><span
  class=\"o\">=</span><span class=\"n\">edit_link</span><span class=\"p\">,</span>\n
  \                   <span class=\"n\">tree</span><span class=\"o\">=</span><span
  class=\"n\">tree</span><span class=\"p\">,</span>\n                    <span class=\"n\">datetime</span><span
  class=\"o\">=</span><span class=\"n\">datetime</span><span class=\"p\">,</span>\n
  \                   <span class=\"n\">nodes</span><span class=\"o\">=</span><span
  class=\"n\">nodes</span><span class=\"p\">,</span>\n                    <span class=\"n\">raw_source</span><span
  class=\"o\">=</span><span class=\"n\">raw_source</span><span class=\"p\">,</span>\n
  \                   <span class=\"n\">indent</span><span class=\"o\">=</span><span
  class=\"n\">textwrap</span><span class=\"o\">.</span><span class=\"n\">indent</span><span
  class=\"p\">,</span>\n                <span class=\"p\">)</span>\n                <span
  class=\"n\">cache</span><span class=\"o\">.</span><span class=\"n\">add</span><span
  class=\"p\">(</span>\n                    <span class=\"n\">key</span><span class=\"p\">,</span>\n
  \                   <span class=\"n\">article</span><span class=\"p\">,</span>\n
  \                   <span class=\"n\">expire</span><span class=\"o\">=</span><span
  class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">default_cache_expire</span><span class=\"p\">,</span>\n
  \               <span class=\"p\">)</span>\n\n            <span class=\"k\">try</span><span
  class=\"p\">:</span>\n                <span class=\"n\">article</span> <span class=\"o\">=</span>
  <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">Post</span><span
  class=\"p\">(</span>\n                    <span class=\"n\">markata</span><span
  class=\"o\">=</span><span class=\"n\">markata</span><span class=\"p\">,</span>\n
  \                   <span class=\"n\">path</span><span class=\"o\">=</span><span
  class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">file</span><span
  class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">replace</span><span
  class=\"p\">(</span><span class=\"s2\">\".py\"</span><span class=\"p\">,</span>
  <span class=\"s2\">\".md\"</span><span class=\"p\">),</span>\n                    <span
  class=\"n\">title</span><span class=\"o\">=</span><span class=\"n\">file</span><span
  class=\"o\">.</span><span class=\"n\">name</span><span class=\"p\">,</span>\n                    <span
  class=\"n\">content</span><span class=\"o\">=</span><span class=\"n\">article</span><span
  class=\"p\">,</span>\n                    <span class=\"n\">ast</span><span class=\"o\">=</span><span
  class=\"n\">ast</span><span class=\"p\">,</span>\n                    <span class=\"n\">file</span><span
  class=\"o\">=</span><span class=\"n\">file</span><span class=\"p\">,</span>\n                    <span
  class=\"n\">slug</span><span class=\"o\">=</span><span class=\"n\">slug</span><span
  class=\"p\">,</span>\n                    <span class=\"n\">edit_link</span><span
  class=\"o\">=</span><span class=\"n\">edit_link</span><span class=\"p\">,</span>\n
  \                   <span class=\"n\">datetime</span><span class=\"o\">=</span><span
  class=\"n\">datetime</span><span class=\"p\">,</span>\n                <span class=\"p\">)</span>\n\n
  \           <span class=\"k\">except</span> <span class=\"n\">pydantic</span><span
  class=\"o\">.</span><span class=\"n\">ValidationError</span> <span class=\"k\">as</span>
  <span class=\"n\">e</span><span class=\"p\">:</span>\n                <span class=\"kn\">from</span>
  <span class=\"nn\">markata.plugins.load</span> <span class=\"kn\">import</span>
  <span class=\"n\">ValidationError</span><span class=\"p\">,</span> <span class=\"n\">get_models</span>\n\n
  \               <span class=\"n\">models</span> <span class=\"o\">=</span> <span
  class=\"n\">get_models</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
  class=\"o\">=</span><span class=\"n\">markata</span><span class=\"p\">,</span> <span
  class=\"n\">error</span><span class=\"o\">=</span><span class=\"n\">e</span><span
  class=\"p\">)</span>\n                <span class=\"n\">models</span> <span class=\"o\">=</span>
  <span class=\"nb\">list</span><span class=\"p\">(</span><span class=\"n\">models</span><span
  class=\"o\">.</span><span class=\"n\">values</span><span class=\"p\">())</span>\n
  \               <span class=\"n\">models</span> <span class=\"o\">=</span> <span
  class=\"s2\">\"</span><span class=\"se\">\\n</span><span class=\"s2\">\"</span><span
  class=\"o\">.</span><span class=\"n\">join</span><span class=\"p\">(</span><span
  class=\"n\">models</span><span class=\"p\">)</span>\n                <span class=\"k\">raise</span>
  <span class=\"n\">ValidationError</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
  class=\"s2\">\"</span><span class=\"si\">{</span><span class=\"n\">e</span><span
  class=\"si\">}</span><span class=\"se\">\\n\\n</span><span class=\"si\">{</span><span
  class=\"n\">models</span><span class=\"si\">}</span><span class=\"se\">\\n</span><span
  class=\"s2\">failed to load </span><span class=\"si\">{</span><span class=\"n\">path</span><span
  class=\"si\">}</span><span class=\"s2\">\"</span><span class=\"p\">)</span> <span
  class=\"kn\">from</span> <span class=\"nn\">e</span>\n\n            <span class=\"k\">return</span>
  <span class=\"n\">article</span>\n</pre></div>\n\n</pre>\n<p>!! function </p><h2
  class=\"admonition-title\" id=\"load\" style=\"margin:0;padding:.5rem 1rem;\">load
  <em class=\"small\">function</em></h2>\nsimilar to <a href=\"../glob\">glob</a>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"load
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
  <span class=\"nf\">load</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
  class=\"p\">:</span> <span class=\"s2\">\"MarkataDocs\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
  class=\"w\">            </span><span class=\"sd\">\"\"\"</span>\n<span class=\"sd\">
  \           similar to [glob](../glob)</span>\n<span class=\"sd\">            \"\"\"</span>\n
  \           <span class=\"k\">if</span> <span class=\"s2\">\"articles\"</span> <span
  class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"vm\">__dict__</span><span class=\"p\">:</span>\n
  \               <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">articles</span> <span class=\"o\">=</span> <span class=\"p\">[]</span>\n
  \           <span class=\"k\">for</span> <span class=\"n\">py_file</span> <span
  class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">py_files</span><span class=\"p\">:</span>\n                <span class=\"k\">with</span>
  <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">cache</span>
  <span class=\"k\">as</span> <span class=\"n\">cache</span><span class=\"p\">:</span>\n
  \                   <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">articles</span><span class=\"o\">.</span><span class=\"n\">append</span><span
  class=\"p\">(</span><span class=\"n\">make_article</span><span class=\"p\">(</span><span
  class=\"n\">markata</span><span class=\"p\">,</span> <span class=\"n\">py_file</span><span
  class=\"p\">,</span> <span class=\"n\">cache</span><span class=\"p\">))</span>\n</pre></div>\n\n</pre>\n<p>!!
  class </p><h2 class=\"admonition-title\" id=\"MarkataDocs\" style=\"margin:0;padding:.5rem
  1rem;\">MarkataDocs <em class=\"small\">class</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"MarkataDocs
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
  <span class=\"nc\">MarkataDocs</span><span class=\"p\">(</span><span class=\"n\">Markata</span><span
  class=\"p\">):</span>\n                <span class=\"n\">py_files</span><span class=\"p\">:</span>
  <span class=\"n\">List</span> <span class=\"o\">=</span> <span class=\"p\">[]</span>\n
  \               <span class=\"n\">content_directories</span><span class=\"p\">:</span>
  <span class=\"n\">List</span> <span class=\"o\">=</span> <span class=\"p\">[]</span>\n</pre></div>\n\n</pre>\n</main>\n<footer>\xA9
  2024</footer>\n</body></html>"
published: true
slug: markata/plugins/docs
title: Docs.Py


---

leading docstring


!! function <h2 id='add_parents' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>add_parents <em class='small'>function</em></h2>

???+ source "add_parents <em class='small'>source</em>"

```python

        def add_parents(tree: ast.AST) -> None:
            for node in ast.walk(tree):
                for child in ast.iter_child_nodes(node):
                    child.parent = node
                    if not hasattr(child, "parents"):
                        child.parents = [node]
                    child.parents.append(node)
                    if isinstance(node, ast.ClassDef) and isinstance(child, ast.FunctionDef):
                        child.type = "method"
                    elif isinstance(child, ast.FunctionDef):
                        child.type = "function"
                    elif isinstance(child, ast.ClassDef):
                        child.type = "class"
```


!! function <h2 id='glob' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>glob <em class='small'>function</em></h2>
    finds k

    ## Parameters

    `markata` the markata object
???+ source "glob <em class='small'>source</em>"

```python

        def glob(markata: "MarkataDocs") -> None:
            """
            finds k

            ## Parameters

            `markata` the markata object

            """

            import glob

            markata.py_files = [Path(f) for f in glob.glob("**/*.py", recursive=True)]

            content_directories = list({f.parent for f in markata.py_files})
            if "content_directories" in markata.__dict__:
                markata.content_directories.extend(content_directories)
            else:
                markata.content_directories = content_directories

            try:
                ignore = True
            except KeyError:
                ignore = True

            if ignore and (Path(".gitignore").exists() or Path(".markataignore").exists()):
                import pathspec

                lines = []

                if Path(".gitignore").exists():
                    lines.extend(Path(".gitignore").read_text().splitlines())

                if Path(".markataignore").exists():
                    lines.extend(Path(".markataignore").read_text().splitlines())

            spec = pathspec.PathSpec.from_lines("gitwildmatch", lines)

            markata.py_files = [
                file for file in markata.py_files if not spec.match_file(str(file))
            ]
```


!! function <h2 id='get_template' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get_template <em class='small'>function</em></h2>

???+ source "get_template <em class='small'>source</em>"

```python

        def get_template():
            jinja_env = jinja2.Environment()
            template = jinja_env.from_string(
                (Path(__file__).parent / "default_doc_template.md").read_text(),
            )
            return template
```


!! function <h2 id='make_article' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>make_article <em class='small'>function</em></h2>

???+ source "make_article <em class='small'>source</em>"

```python

        def make_article(markata: "Markata", file: Path, cache) -> frontmatter.Post:
            with open(file) as f:
                raw_source = f.read()
            key = markata.make_hash("docs", "file", raw_source)
            slug = f"{file.parent}/{file.stem}".lstrip("/").lstrip("./")
            edit_link = (
                str(markata.config.get("repo_url", "https://github.com/"))
                + "edit/"
                + str(markata.config.get("repo_branch", "main"))
                + "/"
                + str(file)
            )
            article_from_cache = markata.precache.get(key)
            if article_from_cache is not None:
                article = article_from_cache
            else:
                tree = ast.parse(raw_source)
                add_parents(tree)
                nodes = [
                    n for n in ast.walk(tree) if isinstance(n, (ast.FunctionDef, ast.ClassDef))
                ]

                article = get_template().render(
                    ast=ast,
                    file=file,
                    slug=slug,
                    edit_link=edit_link,
                    tree=tree,
                    datetime=datetime,
                    nodes=nodes,
                    raw_source=raw_source,
                    indent=textwrap.indent,
                )
                cache.add(
                    key,
                    article,
                    expire=markata.config.default_cache_expire,
                )

            try:
                article = markata.Post(
                    markata=markata,
                    path=str(file).replace(".py", ".md"),
                    title=file.name,
                    content=article,
                    ast=ast,
                    file=file,
                    slug=slug,
                    edit_link=edit_link,
                    datetime=datetime,
                )

            except pydantic.ValidationError as e:
                from markata.plugins.load import ValidationError, get_models

                models = get_models(markata=markata, error=e)
                models = list(models.values())
                models = "\n".join(models)
                raise ValidationError(f"{e}\n\n{models}\nfailed to load {path}") from e

            return article
```


!! function <h2 id='load' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>load <em class='small'>function</em></h2>
    similar to [glob](../glob)
???+ source "load <em class='small'>source</em>"

```python

        def load(markata: "MarkataDocs") -> None:
            """
            similar to [glob](../glob)
            """
            if "articles" not in markata.__dict__:
                markata.articles = []
            for py_file in markata.py_files:
                with markata.cache as cache:
                    markata.articles.append(make_article(markata, py_file, cache))
```


!! class <h2 id='MarkataDocs' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>MarkataDocs <em class='small'>class</em></h2>

???+ source "MarkataDocs <em class='small'>source</em>"

```python

        class MarkataDocs(Markata):
                py_files: List = []
                content_directories: List = []
```
