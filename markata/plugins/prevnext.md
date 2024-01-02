---
content: "The prevnext plugin, creates previous and next links inside each post.\n\n##
  Example config\n\nIn this example we have two maps of posts to look through.  prevnext
  will look\nthrough each of these lists of posts for the current post, then return
  the post\nbefore and after this post as the prevnext posts.\n\n``` toml\n\n[markata]\n#
  default colors will be taken from markata's color_text and color_accent\ncolor_text=white\ncolor_text_light=black\ncolor_accent=white\ncolor_accent_light=black\n\n[markata.prevnext]\n#
  strategy can be 'first' or 'all'\n# 'first' will cycle through the first map the
  post is found in.\n# 'all' will cycle through all of the maps\nstrategy='first'\n\n#
  if you want different colors than your main color_text and color_accent, then\n#
  you can override it here\n# colors can be any valid css color format\n\nprevnext_color_text=white\nprevnext_color_text_light=black\nprevnext_color_angle=white\nprevnext_color_angle_light=black\n\n\n#
  you can have multiple maps, the order they appear will determine their preference\n[markata.feeds.python]\nfilter='\"python\"
  in tags'\nsort='slug'\n\n[markata.feeds.others]\nfilter='\"python\" not in tags'\nsort='slug'\n```\n\nThe
  configuration below will setup two maps, one where \"python\" is in the list\nof
  tags, and another where it is not.  This will link all python posts together\nwith
  a prevnext cycle, and all non-python posts in a separate prevnext cycle.\n\n## strategy\n\nThere
  are currently two supported strategies.\n\n* first\n* all\n\n### first\n\n`first`
  will cycle through only the posts contained within the first map that\ncontains
  the post.\n\n### all\n\n`all` will cycle through all of the posts aggregated from
  any prevnext map.\n\n\n!! class <h2 id='UnsupportedPrevNextStrategy' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>UnsupportedPrevNextStrategy <em class='small'>class</em></h2>\n
  \   A custom error class to raise when an unsupporte prevnext strategy is\n    defined.\n???+
  source \"UnsupportedPrevNextStrategy <em class='small'>source</em>\"\n\n```python\n\n
  \       class UnsupportedPrevNextStrategy(NotImplementedError):\n            \"\"\"\n
  \           A custom error class to raise when an unsupporte prevnext strategy is\n
  \           defined.\n            \"\"\"\n```\n\n\n!! class <h2 id='PrevNext' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>PrevNext <em class='small'>class</em></h2>\n\n???+
  source \"PrevNext <em class='small'>source</em>\"\n\n```python\n\n        class
  PrevNext:\n            prev: str\n            next: str\n```\n\n\n!! function <h2
  id='prevnext' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>prevnext
  <em class='small'>function</em></h2>\n\n???+ source \"prevnext <em class='small'>source</em>\"\n\n```python\n\n
  \       def prevnext(\n            markata: \"Markata\",\n            post: \"Post\",\n
  \           conf: List[Dict[str, str]],\n            strategy: str = \"first\",\n
  \       ) -> Optional[PrevNext]:\n            posts = []\n            for map_conf
  in conf.values():\n                _posts = markata.map(\"post\", **map_conf)\n
  \               # if the strategy is first, cycle back to the beginning after each
  map\n                if strategy == \"first\" and _posts:\n                    _posts.append(_posts[0])\n
  \               posts.extend(_posts)\n            # if the strategy is 'all', cycle
  back to the beginning after all of the maps.\n            if strategy == \"all\":\n
  \               posts.append(posts[0])\n\n            try:\n                post_idx
  = posts.index(post)\n                return PrevNext(prev=posts[post_idx - 1], next=posts[post_idx
  + 1])\n            except ValueError:\n                # post is not in posts\n
  \               return None\n```\n\n\n!! function <h2 id='pre_render' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>pre_render <em class='small'>function</em></h2>\n\n???+
  source \"pre_render <em class='small'>source</em>\"\n\n```python\n\n        def
  pre_render(markata: \"Markata\") -> None:\n            config = markata.config.get(\"prevnext\",
  {})\n            feed_config = markata.config.get(\"feeds\", {})\n            strategy
  = config.get(\"strategy\", \"first\")\n            if strategy not in SUPPORTED_STRATEGIES:\n
  \               msg = f\"\"\"\n                \"{strategy}\" is not a supported
  prevnext strategy\n\n                configure prevnext in your markata.toml to
  use one of {SUPPORTED_STRATEGIES}\n                \"\"\"\n                raise
  UnsupportedPrevNextStrategy(msg)\n            template = config.get(\"template\",
  None)\n            if template is None:\n                template = Template(TEMPLATE)\n
  \           else:\n                template = Template(Path(template).read_text())\n\n
  \           _full_config = copy.deepcopy(markata.config)\n            for article
  in set(markata.articles):\n                article[\"prevnext\"] = prevnext(\n                    markata,\n
  \                   article,\n                    feed_config,\n                    strategy=strategy,\n
  \               )\n                if \"prevnext\" not in article.content and article[\"prevnext\"]:\n
  \                   article.content += template.render(\n                        config=always_merger.merge(\n
  \                           _full_config,\n                            copy.deepcopy(\n
  \                               article.get(\n                                    \"config_overrides\",\n
  \                                   {},\n                                ),\n                            ),\n
  \                       ),\n                        **article,\n                    )\n```\n"
date: 0001-01-01
description: The prevnext plugin, creates previous and next links inside each post.
  In this example we have two maps of posts to look through.  prevnext will look The
  config
html: "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n<title>Prevnext.Py</title>\n<meta
  charset=\"utf-8\"/>\n<meta content=\"width=device-width, initial-scale=1\" name=\"viewport\"/>\n<meta
  content=\"The prevnext plugin, creates previous and next links inside each post.
  In this example we have two maps of posts to look through.  prevnext will look The
  config\" name=\"description\"/>\n<link href=\"/static/favicon.ico\" rel=\"icon\"
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
  name=\"og:type\" property=\"og:type\"/><meta content=\"The prevnext plugin, creates
  previous and next links inside each post. In this example we have two maps of posts
  to look through.  prevnext will look The config\" name=\"description\" property=\"description\"/><meta
  content=\"The prevnext plugin, creates previous and next links inside each post.
  In this example we have two maps of posts to look through.  prevnext will look The
  config\" name=\"og:description\" property=\"og:description\"/><meta content=\"The
  prevnext plugin, creates previous and next links inside each post. In this example
  we have two maps of posts to look through.  prevnext will look The config\" name=\"twitter:description\"
  property=\"twitter:description\"/><meta content=\"Prevnext.Py | Markata\" name=\"og:title\"
  property=\"og:title\"/><meta content=\"Prevnext.Py | Markata\" name=\"twitter:title\"
  property=\"twitter:title\"/><meta content=\"https://markata.dev//markata/plugins/prevnext-og.png\"
  name=\"og:image\" property=\"og:image\"/><meta content=\"https://markata.dev//markata/plugins/prevnext-og.png\"
  name=\"twitter:image\" property=\"twitter:image\"/><meta content=\"1600\" name=\"og:image:width\"
  property=\"og:image:width\"/><meta content=\"900\" name=\"og:image:width\" property=\"og:image:width\"/><meta
  content=\"summary_large_image\" name=\"twitter:card\" property=\"twitter:card\"/><meta
  content=\"Markata\" name=\"og:site_name\" property=\"og:site_name\"/><meta content=\"@_waylonwalker\"
  name=\"twitter:creator\" property=\"twitter:creator\"/><meta content=\"Prevnext.Py\"
  name=\"title\" property=\"title\"/><meta content=\"markata 0.8.0.dev6\" name=\"generator\"
  property=\"generator\"/><link href=\"https://markata.dev//markata/plugins/prevnext/\"
  rel=\"canonical\"/><meta content=\"https://markata.dev//markata/plugins/prevnext/\"
  name=\"og:url\" property=\"og:url\"/></head>\n<body><nav>\n<a href=\"/\">markata</a>\n<a
  href=\"https://github.com/WaylonWalker/markata\">GitHub</a>\n</nav>\n<div>\n<label
  class=\"theme-switch\" for=\"checkbox-theme\" id=\"theme-switch\" title=\"light/dark
  mode toggle\">\n<input id=\"checkbox-theme\" type=\"checkbox\"/>\n<div class=\"slider
  round\"></div>\n</label>\n</div>\n<section class=\"title\">\n<h1 id=\"title\">\n
  \           Prevnext.Py \n            \n        </h1>\n</section>\n<main><p>The
  prevnext plugin, creates previous and next links inside each post.</p>\n<h2 id=\"example-config\">Example
  config <a class=\"header-anchor\" href=\"#example-config\"><svg aria-hidden=\"true\"
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
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>In this example we have
  two maps of posts to look through.  prevnext will look\nthrough each of these lists
  of posts for the current post, then return the post\nbefore and after this post
  as the prevnext posts.</p>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
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
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n<span
  class=\"c1\"># default colors will be taken from markata's color_text and color_accent</span>\n<span
  class=\"n\">color_text</span><span class=\"o\">=</span><span class=\"err\">white</span>\n<span
  class=\"n\">color_text_light</span><span class=\"o\">=</span><span class=\"err\">black</span>\n<span
  class=\"n\">color_accent</span><span class=\"o\">=</span><span class=\"err\">white</span>\n<span
  class=\"n\">color_accent_light</span><span class=\"o\">=</span><span class=\"err\">black</span>\n\n<span
  class=\"k\">[markata.prevnext]</span>\n<span class=\"c1\"># strategy can be 'first'
  or 'all'</span>\n<span class=\"c1\"># 'first' will cycle through the first map the
  post is found in.</span>\n<span class=\"c1\"># 'all' will cycle through all of the
  maps</span>\n<span class=\"n\">strategy</span><span class=\"o\">=</span><span class=\"s1\">'first'</span>\n\n<span
  class=\"c1\"># if you want different colors than your main color_text and color_accent,
  then</span>\n<span class=\"c1\"># you can override it here</span>\n<span class=\"c1\">#
  colors can be any valid css color format</span>\n\n<span class=\"n\">prevnext_color_text</span><span
  class=\"o\">=</span><span class=\"err\">white</span>\n<span class=\"n\">prevnext_color_text_light</span><span
  class=\"o\">=</span><span class=\"err\">black</span>\n<span class=\"n\">prevnext_color_angle</span><span
  class=\"o\">=</span><span class=\"err\">white</span>\n<span class=\"n\">prevnext_color_angle_light</span><span
  class=\"o\">=</span><span class=\"err\">black</span>\n\n\n<span class=\"c1\"># you
  can have multiple maps, the order they appear will determine their preference</span>\n<span
  class=\"k\">[markata.feeds.python]</span>\n<span class=\"n\">filter</span><span
  class=\"o\">=</span><span class=\"s1\">'\"python\" in tags'</span>\n<span class=\"n\">sort</span><span
  class=\"o\">=</span><span class=\"s1\">'slug'</span>\n\n<span class=\"k\">[markata.feeds.others]</span>\n<span
  class=\"n\">filter</span><span class=\"o\">=</span><span class=\"s1\">'\"python\"
  not in tags'</span>\n<span class=\"n\">sort</span><span class=\"o\">=</span><span
  class=\"s1\">'slug'</span>\n</pre></div>\n\n</pre>\n<p>The configuration below will
  setup two maps, one where \"python\" is in the list\nof tags, and another where
  it is not.  This will link all python posts together\nwith a prevnext cycle, and
  all non-python posts in a separate prevnext cycle.</p>\n<h2 id=\"strategy\">strategy
  <a class=\"header-anchor\" href=\"#strategy\"><svg aria-hidden=\"true\" class=\"heading-permalink\"
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
  are currently two supported strategies.</p>\n<ul>\n<li>first</li>\n<li>all</li>\n</ul>\n<h3>first</h3>\n<p><code>first</code>
  will cycle through only the posts contained within the first map that\ncontains
  the post.</p>\n<h3>all</h3>\n<p><code>all</code> will cycle through all of the posts
  aggregated from any prevnext map.</p>\n<p>!! class </p><h2 class=\"admonition-title\"
  id=\"UnsupportedPrevNextStrategy\" style=\"margin:0;padding:.5rem 1rem;\">UnsupportedPrevNextStrategy
  <em class=\"small\">class</em></h2>\nA custom error class to raise when an unsupporte
  prevnext strategy is\ndefined.\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
  class=\"admonition-title\">\"UnsupportedPrevNextStrategy <em class=\"small\">source</em>\"</p>\n</div>\n<pre
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
  <span class=\"nc\">UnsupportedPrevNextStrategy</span><span class=\"p\">(</span><span
  class=\"ne\">NotImplementedError</span><span class=\"p\">):</span>\n<span class=\"w\">
  \           </span><span class=\"sd\">\"\"\"</span>\n<span class=\"sd\">            A
  custom error class to raise when an unsupporte prevnext strategy is</span>\n<span
  class=\"sd\">            defined.</span>\n<span class=\"sd\">            \"\"\"</span>\n</pre></div>\n\n</pre>\n<p>!!
  class </p><h2 class=\"admonition-title\" id=\"PrevNext\" style=\"margin:0;padding:.5rem
  1rem;\">PrevNext <em class=\"small\">class</em></h2>\n<div class=\"admonition source
  is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"PrevNext <em
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
  <span class=\"nc\">PrevNext</span><span class=\"p\">:</span>\n            <span
  class=\"n\">prev</span><span class=\"p\">:</span> <span class=\"nb\">str</span>\n
  \           <span class=\"nb\">next</span><span class=\"p\">:</span> <span class=\"nb\">str</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"prevnext\" style=\"margin:0;padding:.5rem
  1rem;\">prevnext <em class=\"small\">function</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"prevnext
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
  <span class=\"nf\">prevnext</span><span class=\"p\">(</span>\n            <span
  class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span
  class=\"p\">,</span>\n            <span class=\"n\">post</span><span class=\"p\">:</span>
  <span class=\"s2\">\"Post\"</span><span class=\"p\">,</span>\n            <span
  class=\"n\">conf</span><span class=\"p\">:</span> <span class=\"n\">List</span><span
  class=\"p\">[</span><span class=\"n\">Dict</span><span class=\"p\">[</span><span
  class=\"nb\">str</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
  class=\"p\">]],</span>\n            <span class=\"n\">strategy</span><span class=\"p\">:</span>
  <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">\"first\"</span><span
  class=\"p\">,</span>\n        <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"n\">PrevNext</span><span
  class=\"p\">]:</span>\n            <span class=\"n\">posts</span> <span class=\"o\">=</span>
  <span class=\"p\">[]</span>\n            <span class=\"k\">for</span> <span class=\"n\">map_conf</span>
  <span class=\"ow\">in</span> <span class=\"n\">conf</span><span class=\"o\">.</span><span
  class=\"n\">values</span><span class=\"p\">():</span>\n                <span class=\"n\">_posts</span>
  <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">map</span><span class=\"p\">(</span><span class=\"s2\">\"post\"</span><span
  class=\"p\">,</span> <span class=\"o\">**</span><span class=\"n\">map_conf</span><span
  class=\"p\">)</span>\n                <span class=\"c1\"># if the strategy is first,
  cycle back to the beginning after each map</span>\n                <span class=\"k\">if</span>
  <span class=\"n\">strategy</span> <span class=\"o\">==</span> <span class=\"s2\">\"first\"</span>
  <span class=\"ow\">and</span> <span class=\"n\">_posts</span><span class=\"p\">:</span>\n
  \                   <span class=\"n\">_posts</span><span class=\"o\">.</span><span
  class=\"n\">append</span><span class=\"p\">(</span><span class=\"n\">_posts</span><span
  class=\"p\">[</span><span class=\"mi\">0</span><span class=\"p\">])</span>\n                <span
  class=\"n\">posts</span><span class=\"o\">.</span><span class=\"n\">extend</span><span
  class=\"p\">(</span><span class=\"n\">_posts</span><span class=\"p\">)</span>\n
  \           <span class=\"c1\"># if the strategy is 'all', cycle back to the beginning
  after all of the maps.</span>\n            <span class=\"k\">if</span> <span class=\"n\">strategy</span>
  <span class=\"o\">==</span> <span class=\"s2\">\"all\"</span><span class=\"p\">:</span>\n
  \               <span class=\"n\">posts</span><span class=\"o\">.</span><span class=\"n\">append</span><span
  class=\"p\">(</span><span class=\"n\">posts</span><span class=\"p\">[</span><span
  class=\"mi\">0</span><span class=\"p\">])</span>\n\n            <span class=\"k\">try</span><span
  class=\"p\">:</span>\n                <span class=\"n\">post_idx</span> <span class=\"o\">=</span>
  <span class=\"n\">posts</span><span class=\"o\">.</span><span class=\"n\">index</span><span
  class=\"p\">(</span><span class=\"n\">post</span><span class=\"p\">)</span>\n                <span
  class=\"k\">return</span> <span class=\"n\">PrevNext</span><span class=\"p\">(</span><span
  class=\"n\">prev</span><span class=\"o\">=</span><span class=\"n\">posts</span><span
  class=\"p\">[</span><span class=\"n\">post_idx</span> <span class=\"o\">-</span>
  <span class=\"mi\">1</span><span class=\"p\">],</span> <span class=\"nb\">next</span><span
  class=\"o\">=</span><span class=\"n\">posts</span><span class=\"p\">[</span><span
  class=\"n\">post_idx</span> <span class=\"o\">+</span> <span class=\"mi\">1</span><span
  class=\"p\">])</span>\n            <span class=\"k\">except</span> <span class=\"ne\">ValueError</span><span
  class=\"p\">:</span>\n                <span class=\"c1\"># post is not in posts</span>\n
  \               <span class=\"k\">return</span> <span class=\"kc\">None</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"pre_render\" style=\"margin:0;padding:.5rem
  1rem;\">pre_render <em class=\"small\">function</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"pre_render
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
  <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \           <span class=\"n\">config</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
  class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">\"prevnext\"</span><span
  class=\"p\">,</span> <span class=\"p\">{})</span>\n            <span class=\"n\">feed_config</span>
  <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">get</span><span
  class=\"p\">(</span><span class=\"s2\">\"feeds\"</span><span class=\"p\">,</span>
  <span class=\"p\">{})</span>\n            <span class=\"n\">strategy</span> <span
  class=\"o\">=</span> <span class=\"n\">config</span><span class=\"o\">.</span><span
  class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">\"strategy\"</span><span
  class=\"p\">,</span> <span class=\"s2\">\"first\"</span><span class=\"p\">)</span>\n
  \           <span class=\"k\">if</span> <span class=\"n\">strategy</span> <span
  class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">SUPPORTED_STRATEGIES</span><span
  class=\"p\">:</span>\n                <span class=\"n\">msg</span> <span class=\"o\">=</span>
  <span class=\"sa\">f</span><span class=\"s2\">\"\"\"</span>\n<span class=\"s2\">
  \               \"</span><span class=\"si\">{</span><span class=\"n\">strategy</span><span
  class=\"si\">}</span><span class=\"s2\">\" is not a supported prevnext strategy</span>\n\n<span
  class=\"s2\">                configure prevnext in your markata.toml to use one
  of </span><span class=\"si\">{</span><span class=\"n\">SUPPORTED_STRATEGIES</span><span
  class=\"si\">}</span>\n<span class=\"s2\">                \"\"\"</span>\n                <span
  class=\"k\">raise</span> <span class=\"n\">UnsupportedPrevNextStrategy</span><span
  class=\"p\">(</span><span class=\"n\">msg</span><span class=\"p\">)</span>\n            <span
  class=\"n\">template</span> <span class=\"o\">=</span> <span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
  class=\"s2\">\"template\"</span><span class=\"p\">,</span> <span class=\"kc\">None</span><span
  class=\"p\">)</span>\n            <span class=\"k\">if</span> <span class=\"n\">template</span>
  <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \               <span class=\"n\">template</span> <span class=\"o\">=</span> <span
  class=\"n\">Template</span><span class=\"p\">(</span><span class=\"n\">TEMPLATE</span><span
  class=\"p\">)</span>\n            <span class=\"k\">else</span><span class=\"p\">:</span>\n
  \               <span class=\"n\">template</span> <span class=\"o\">=</span> <span
  class=\"n\">Template</span><span class=\"p\">(</span><span class=\"n\">Path</span><span
  class=\"p\">(</span><span class=\"n\">template</span><span class=\"p\">)</span><span
  class=\"o\">.</span><span class=\"n\">read_text</span><span class=\"p\">())</span>\n\n
  \           <span class=\"n\">_full_config</span> <span class=\"o\">=</span> <span
  class=\"n\">copy</span><span class=\"o\">.</span><span class=\"n\">deepcopy</span><span
  class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">config</span><span class=\"p\">)</span>\n            <span class=\"k\">for</span>
  <span class=\"n\">article</span> <span class=\"ow\">in</span> <span class=\"nb\">set</span><span
  class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">articles</span><span class=\"p\">):</span>\n                <span class=\"n\">article</span><span
  class=\"p\">[</span><span class=\"s2\">\"prevnext\"</span><span class=\"p\">]</span>
  <span class=\"o\">=</span> <span class=\"n\">prevnext</span><span class=\"p\">(</span>\n
  \                   <span class=\"n\">markata</span><span class=\"p\">,</span>\n
  \                   <span class=\"n\">article</span><span class=\"p\">,</span>\n
  \                   <span class=\"n\">feed_config</span><span class=\"p\">,</span>\n
  \                   <span class=\"n\">strategy</span><span class=\"o\">=</span><span
  class=\"n\">strategy</span><span class=\"p\">,</span>\n                <span class=\"p\">)</span>\n
  \               <span class=\"k\">if</span> <span class=\"s2\">\"prevnext\"</span>
  <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">article</span><span
  class=\"o\">.</span><span class=\"n\">content</span> <span class=\"ow\">and</span>
  <span class=\"n\">article</span><span class=\"p\">[</span><span class=\"s2\">\"prevnext\"</span><span
  class=\"p\">]:</span>\n                    <span class=\"n\">article</span><span
  class=\"o\">.</span><span class=\"n\">content</span> <span class=\"o\">+=</span>
  <span class=\"n\">template</span><span class=\"o\">.</span><span class=\"n\">render</span><span
  class=\"p\">(</span>\n                        <span class=\"n\">config</span><span
  class=\"o\">=</span><span class=\"n\">always_merger</span><span class=\"o\">.</span><span
  class=\"n\">merge</span><span class=\"p\">(</span>\n                            <span
  class=\"n\">_full_config</span><span class=\"p\">,</span>\n                            <span
  class=\"n\">copy</span><span class=\"o\">.</span><span class=\"n\">deepcopy</span><span
  class=\"p\">(</span>\n                                <span class=\"n\">article</span><span
  class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span>\n                                    <span
  class=\"s2\">\"config_overrides\"</span><span class=\"p\">,</span>\n                                    <span
  class=\"p\">{},</span>\n                                <span class=\"p\">),</span>\n
  \                           <span class=\"p\">),</span>\n                        <span
  class=\"p\">),</span>\n                        <span class=\"o\">**</span><span
  class=\"n\">article</span><span class=\"p\">,</span>\n                    <span
  class=\"p\">)</span>\n</pre></div>\n\n</pre>\n</main>\n<footer>\xA9 2024</footer>\n</body></html>"
published: true
slug: markata/plugins/prevnext
title: Prevnext.Py


---

The prevnext plugin, creates previous and next links inside each post.

## Example config

In this example we have two maps of posts to look through.  prevnext will look
through each of these lists of posts for the current post, then return the post
before and after this post as the prevnext posts.

``` toml

[markata]
# default colors will be taken from markata's color_text and color_accent
color_text=white
color_text_light=black
color_accent=white
color_accent_light=black

[markata.prevnext]
# strategy can be 'first' or 'all'
# 'first' will cycle through the first map the post is found in.
# 'all' will cycle through all of the maps
strategy='first'

# if you want different colors than your main color_text and color_accent, then
# you can override it here
# colors can be any valid css color format

prevnext_color_text=white
prevnext_color_text_light=black
prevnext_color_angle=white
prevnext_color_angle_light=black


# you can have multiple maps, the order they appear will determine their preference
[markata.feeds.python]
filter='"python" in tags'
sort='slug'

[markata.feeds.others]
filter='"python" not in tags'
sort='slug'
```

The configuration below will setup two maps, one where "python" is in the list
of tags, and another where it is not.  This will link all python posts together
with a prevnext cycle, and all non-python posts in a separate prevnext cycle.

## strategy

There are currently two supported strategies.

* first
* all

### first

`first` will cycle through only the posts contained within the first map that
contains the post.

### all

`all` will cycle through all of the posts aggregated from any prevnext map.


!! class <h2 id='UnsupportedPrevNextStrategy' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>UnsupportedPrevNextStrategy <em class='small'>class</em></h2>
    A custom error class to raise when an unsupporte prevnext strategy is
    defined.
???+ source "UnsupportedPrevNextStrategy <em class='small'>source</em>"

```python

        class UnsupportedPrevNextStrategy(NotImplementedError):
            """
            A custom error class to raise when an unsupporte prevnext strategy is
            defined.
            """
```


!! class <h2 id='PrevNext' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>PrevNext <em class='small'>class</em></h2>

???+ source "PrevNext <em class='small'>source</em>"

```python

        class PrevNext:
            prev: str
            next: str
```


!! function <h2 id='prevnext' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>prevnext <em class='small'>function</em></h2>

???+ source "prevnext <em class='small'>source</em>"

```python

        def prevnext(
            markata: "Markata",
            post: "Post",
            conf: List[Dict[str, str]],
            strategy: str = "first",
        ) -> Optional[PrevNext]:
            posts = []
            for map_conf in conf.values():
                _posts = markata.map("post", **map_conf)
                # if the strategy is first, cycle back to the beginning after each map
                if strategy == "first" and _posts:
                    _posts.append(_posts[0])
                posts.extend(_posts)
            # if the strategy is 'all', cycle back to the beginning after all of the maps.
            if strategy == "all":
                posts.append(posts[0])

            try:
                post_idx = posts.index(post)
                return PrevNext(prev=posts[post_idx - 1], next=posts[post_idx + 1])
            except ValueError:
                # post is not in posts
                return None
```


!! function <h2 id='pre_render' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>pre_render <em class='small'>function</em></h2>

???+ source "pre_render <em class='small'>source</em>"

```python

        def pre_render(markata: "Markata") -> None:
            config = markata.config.get("prevnext", {})
            feed_config = markata.config.get("feeds", {})
            strategy = config.get("strategy", "first")
            if strategy not in SUPPORTED_STRATEGIES:
                msg = f"""
                "{strategy}" is not a supported prevnext strategy

                configure prevnext in your markata.toml to use one of {SUPPORTED_STRATEGIES}
                """
                raise UnsupportedPrevNextStrategy(msg)
            template = config.get("template", None)
            if template is None:
                template = Template(TEMPLATE)
            else:
                template = Template(Path(template).read_text())

            _full_config = copy.deepcopy(markata.config)
            for article in set(markata.articles):
                article["prevnext"] = prevnext(
                    markata,
                    article,
                    feed_config,
                    strategy=strategy,
                )
                if "prevnext" not in article.content and article["prevnext"]:
                    article.content += template.render(
                        config=always_merger.merge(
                            _full_config,
                            copy.deepcopy(
                                article.get(
                                    "config_overrides",
                                    {},
                                ),
                            ),
                        ),
                        **article,
                    )
```
