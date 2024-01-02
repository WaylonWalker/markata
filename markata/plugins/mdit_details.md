---
content: "None\n\n\n!! function <h2 id='get_tag' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>get_tag <em class='small'>function</em></h2>\n\n???+ source \"get_tag <em
  class='small'>source</em>\"\n\n```python\n\n        def get_tag(params: str) ->
  Tuple[str, str, bool]:\n            open = False\n            if not params.strip():\n
  \               return \"\", \"\", open\n\n            if params.strip().startswith(\"+\"):\n
  \               params = params.strip(\"+\")\n                open = True\n\n            tag,
  *_title = params.strip().split(\" \")\n            joined = \" \".join(_title).strip('\"').strip(\"'\")\n\n
  \           title = \"\"\n            if not joined:\n                title = tag.title()\n
  \           elif joined != '\"\"':\n                title = joined\n            return
  (tag.lower(), title, open)\n```\n\n\n!! function <h2 id='validate' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>validate <em class='small'>function</em></h2>\n\n???+
  source \"validate <em class='small'>source</em>\"\n\n```python\n\n        def validate(params:
  str) -> bool:\n            tag = params.strip().split(\" \", 1)[-1] or \"\"\n            return
  bool(tag)\n```\n\n\n!! function <h2 id='details' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>details <em class='small'>function</em></h2>\n\n???+ source \"details <em
  class='small'>source</em>\"\n\n```python\n\n        def details(state: StateBlock,
  startLine: int, endLine: int, silent: bool) -> bool:\n            start = state.bMarks[startLine]
  + state.tShift[startLine]\n            maximum = state.eMarks[startLine]\n\n            #
  Check out the first character quickly, which should filter out most of non-containers\n
  \           if ord(state.src[start]) != MARKER_CHAR:\n                return False\n\n
  \           # Check out the rest of the marker string\n            pos = start +
  1\n            while pos <= maximum and MARKER_STR[(pos - start) % MARKER_LEN] ==
  state.src[pos]:\n                pos += 1\n\n            marker_count = math.floor((pos
  - start) / MARKER_LEN)\n            if marker_count < MIN_MARKERS:\n                return
  False\n            marker_pos = pos - ((pos - start) % MARKER_LEN)\n            params
  = state.src[marker_pos:maximum]\n            markup = state.src[start:marker_pos]\n\n
  \           if not validate(params):\n                return False\n\n            #
  Since start is found, we can report success here in validation mode\n            if
  silent:\n                return True\n\n            old_parent = state.parentType\n
  \           old_line_max = state.lineMax\n            old_indent = state.blkIndent\n\n
  \           blk_start = pos\n            while blk_start < maximum and state.src[blk_start]
  == \" \":\n                blk_start += 1\n\n            state.parentType = \"details\"\n
  \           state.blkIndent += blk_start - start\n\n            was_empty = False\n\n
  \           # Search for the end of the block\n            next_line = startLine\n
  \           while True:\n                next_line += 1\n                if next_line
  >= endLine:\n                    # unclosed block should be autoclosed by end of
  document.\n                    # also block seems to be autoclosed by end of parent\n
  \                   break\n                pos = state.bMarks[next_line] + state.tShift[next_line]\n
  \               maximum = state.eMarks[next_line]\n                is_empty = state.sCount[next_line]
  < state.blkIndent\n\n                # two consecutive empty lines autoclose the
  block\n                if is_empty and was_empty:\n                    break\n                was_empty
  = is_empty\n\n                if pos < maximum and state.sCount[next_line] < state.blkIndent:\n
  \                   # non-empty line with negative indent should stop the block:\n
  \                   # - !!!\n                    #  test\n                    break\n\n
  \           # this will prevent lazy continuations from ever going past our end
  marker\n            state.lineMax = next_line\n\n            tag, title, open =
  get_tag(params)\n            attrs = {\"class\": f\"details {tag}\"}\n            if
  open:\n                attrs[\"open\"] = \"\"\n\n            token = state.push(\"details_open\",
  \"details\", 1)\n            token.markup = markup\n            token.block = True\n
  \           token.attrs = attrs\n            token.meta = {\"tag\": tag}\n            token.content
  = title\n            token.info = params\n            token.map = [startLine, next_line]\n\n
  \           if title:\n                title_markup = f\"{markup} {tag}\"\n                token
  = state.push(\"details_title_open\", \"summary\", 1)\n                token.markup
  = title_markup\n                token.attrs = {\"class\": \"admonition-title\"}\n
  \               token.map = [startLine, startLine + 1]\n\n                token
  = state.push(\"inline\", \"\", 0)\n                token.content = title\n                token.map
  = [startLine, startLine + 1]\n                token.children = []\n\n                token
  = state.push(\"details_title_close\", \"summary\", -1)\n                token.markup
  = title_markup\n\n            state.md.block.tokenize(state, startLine + 1, next_line)\n\n
  \           token = state.push(\"details_close\", \"details\", -1)\n            token.markup
  = state.src[start:pos]\n            token.block = True\n\n            state.parentType
  = old_parent\n            state.lineMax = old_line_max\n            state.blkIndent
  = old_indent\n            state.line = next_line\n\n            return True\n```\n\n\n!!
  function <h2 id='details_plugin' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>details_plugin <em class='small'>function</em></h2>\n    Plugin to use\n
  \   `python-markdown style detailss\n    <https://python-markdown.github.io/extensions/details>`_.\n\n
  \   .. code-block:: md\n\n        ??? note\n            *content*\n\n    Note, this
  is ported from\n    `markdown-it-admon\n    <https://github.com/commenthol/markdown-it-admon>`_.\n???+
  source \"details_plugin <em class='small'>source</em>\"\n\n```python\n\n        def
  details_plugin(md: MarkdownIt, render: Optional[Callable] = None) -> None:\n            \"\"\"Plugin
  to use\n            `python-markdown style detailss\n            <https://python-markdown.github.io/extensions/details>`_.\n\n
  \           .. code-block:: md\n\n                ??? note\n                    *content*\n\n
  \           Note, this is ported from\n            `markdown-it-admon\n            <https://github.com/commenthol/markdown-it-admon>`_.\n
  \           \"\"\"\n\n            def renderDefault(self, tokens, idx, _options,
  env):\n                return self.renderToken(tokens, idx, _options, env)\n\n            render
  = render or renderDefault\n\n            md.add_render_rule(\"details_open\", render)\n
  \           md.add_render_rule(\"details_close\", render)\n            md.add_render_rule(\"details_title_open\",
  render)\n            md.add_render_rule(\"details_title_close\", render)\n\n            md.block.ruler.before(\n
  \               \"fence\",\n                \"details\",\n                details,\n
  \               {\"alt\": [\"paragraph\", \"reference\", \"blockquote\", \"list\"]},\n
  \           )\n```\n\n\n!! function <h2 id='renderDefault' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>renderDefault <em class='small'>function</em></h2>\n\n???+
  source \"renderDefault <em class='small'>source</em>\"\n\n```python\n\n        def
  renderDefault(self, tokens, idx, _options, env):\n                return self.renderToken(tokens,
  idx, _options, env)\n```\n"
date: 0001-01-01
description: 'None ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
  source '
html: "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n<title>Mdit_Details.Py</title>\n<meta
  charset=\"utf-8\"/>\n<meta content=\"width=device-width, initial-scale=1\" name=\"viewport\"/>\n<meta
  content=\"None ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
  source \" name=\"description\"/>\n<link href=\"/static/favicon.ico\" rel=\"icon\"
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
  name=\"og:type\" property=\"og:type\"/><meta content=\"None ! ???+ source  ! ???+
  source  ! ???+ source  ! ???+ source  ! ???+ source \" name=\"description\" property=\"description\"/><meta
  content=\"None ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
  source \" name=\"og:description\" property=\"og:description\"/><meta content=\"None
  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source \" name=\"twitter:description\"
  property=\"twitter:description\"/><meta content=\"Mdit_Details.Py | Markata\" name=\"og:title\"
  property=\"og:title\"/><meta content=\"Mdit_Details.Py | Markata\" name=\"twitter:title\"
  property=\"twitter:title\"/><meta content=\"https://markata.dev//markata/plugins/mdit-details-og.png\"
  name=\"og:image\" property=\"og:image\"/><meta content=\"https://markata.dev//markata/plugins/mdit-details-og.png\"
  name=\"twitter:image\" property=\"twitter:image\"/><meta content=\"1600\" name=\"og:image:width\"
  property=\"og:image:width\"/><meta content=\"900\" name=\"og:image:width\" property=\"og:image:width\"/><meta
  content=\"summary_large_image\" name=\"twitter:card\" property=\"twitter:card\"/><meta
  content=\"Markata\" name=\"og:site_name\" property=\"og:site_name\"/><meta content=\"@_waylonwalker\"
  name=\"twitter:creator\" property=\"twitter:creator\"/><meta content=\"Mdit_Details.Py\"
  name=\"title\" property=\"title\"/><meta content=\"markata 0.8.0.dev6\" name=\"generator\"
  property=\"generator\"/><link href=\"https://markata.dev//markata/plugins/mdit-details/\"
  rel=\"canonical\"/><meta content=\"https://markata.dev//markata/plugins/mdit-details/\"
  name=\"og:url\" property=\"og:url\"/></head>\n<body><nav>\n<a href=\"/\">markata</a>\n<a
  href=\"https://github.com/WaylonWalker/markata\">GitHub</a>\n</nav>\n<div>\n<label
  class=\"theme-switch\" for=\"checkbox-theme\" id=\"theme-switch\" title=\"light/dark
  mode toggle\">\n<input id=\"checkbox-theme\" type=\"checkbox\"/>\n<div class=\"slider
  round\"></div>\n</label>\n</div>\n<section class=\"title\">\n<h1 id=\"title\">\n
  \           Mdit_Details.Py \n            \n        </h1>\n</section>\n<main><p>None</p>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"get_tag\" style=\"margin:0;padding:.5rem
  1rem;\">get_tag <em class=\"small\">function</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"get_tag
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
  <span class=\"nf\">get_tag</span><span class=\"p\">(</span><span class=\"n\">params</span><span
  class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">)</span> <span
  class=\"o\">-&gt;</span> <span class=\"n\">Tuple</span><span class=\"p\">[</span><span
  class=\"nb\">str</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
  class=\"p\">,</span> <span class=\"nb\">bool</span><span class=\"p\">]:</span>\n
  \           <span class=\"nb\">open</span> <span class=\"o\">=</span> <span class=\"kc\">False</span>\n
  \           <span class=\"k\">if</span> <span class=\"ow\">not</span> <span class=\"n\">params</span><span
  class=\"o\">.</span><span class=\"n\">strip</span><span class=\"p\">():</span>\n
  \               <span class=\"k\">return</span> <span class=\"s2\">\"\"</span><span
  class=\"p\">,</span> <span class=\"s2\">\"\"</span><span class=\"p\">,</span> <span
  class=\"nb\">open</span>\n\n            <span class=\"k\">if</span> <span class=\"n\">params</span><span
  class=\"o\">.</span><span class=\"n\">strip</span><span class=\"p\">()</span><span
  class=\"o\">.</span><span class=\"n\">startswith</span><span class=\"p\">(</span><span
  class=\"s2\">\"+\"</span><span class=\"p\">):</span>\n                <span class=\"n\">params</span>
  <span class=\"o\">=</span> <span class=\"n\">params</span><span class=\"o\">.</span><span
  class=\"n\">strip</span><span class=\"p\">(</span><span class=\"s2\">\"+\"</span><span
  class=\"p\">)</span>\n                <span class=\"nb\">open</span> <span class=\"o\">=</span>
  <span class=\"kc\">True</span>\n\n            <span class=\"n\">tag</span><span
  class=\"p\">,</span> <span class=\"o\">*</span><span class=\"n\">_title</span> <span
  class=\"o\">=</span> <span class=\"n\">params</span><span class=\"o\">.</span><span
  class=\"n\">strip</span><span class=\"p\">()</span><span class=\"o\">.</span><span
  class=\"n\">split</span><span class=\"p\">(</span><span class=\"s2\">\" \"</span><span
  class=\"p\">)</span>\n            <span class=\"n\">joined</span> <span class=\"o\">=</span>
  <span class=\"s2\">\" \"</span><span class=\"o\">.</span><span class=\"n\">join</span><span
  class=\"p\">(</span><span class=\"n\">_title</span><span class=\"p\">)</span><span
  class=\"o\">.</span><span class=\"n\">strip</span><span class=\"p\">(</span><span
  class=\"s1\">'\"'</span><span class=\"p\">)</span><span class=\"o\">.</span><span
  class=\"n\">strip</span><span class=\"p\">(</span><span class=\"s2\">\"'\"</span><span
  class=\"p\">)</span>\n\n            <span class=\"n\">title</span> <span class=\"o\">=</span>
  <span class=\"s2\">\"\"</span>\n            <span class=\"k\">if</span> <span class=\"ow\">not</span>
  <span class=\"n\">joined</span><span class=\"p\">:</span>\n                <span
  class=\"n\">title</span> <span class=\"o\">=</span> <span class=\"n\">tag</span><span
  class=\"o\">.</span><span class=\"n\">title</span><span class=\"p\">()</span>\n
  \           <span class=\"k\">elif</span> <span class=\"n\">joined</span> <span
  class=\"o\">!=</span> <span class=\"s1\">'\"\"'</span><span class=\"p\">:</span>\n
  \               <span class=\"n\">title</span> <span class=\"o\">=</span> <span
  class=\"n\">joined</span>\n            <span class=\"k\">return</span> <span class=\"p\">(</span><span
  class=\"n\">tag</span><span class=\"o\">.</span><span class=\"n\">lower</span><span
  class=\"p\">(),</span> <span class=\"n\">title</span><span class=\"p\">,</span>
  <span class=\"nb\">open</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"validate\" style=\"margin:0;padding:.5rem
  1rem;\">validate <em class=\"small\">function</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"validate
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
  <span class=\"nf\">validate</span><span class=\"p\">(</span><span class=\"n\">params</span><span
  class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">)</span> <span
  class=\"o\">-&gt;</span> <span class=\"nb\">bool</span><span class=\"p\">:</span>\n
  \           <span class=\"n\">tag</span> <span class=\"o\">=</span> <span class=\"n\">params</span><span
  class=\"o\">.</span><span class=\"n\">strip</span><span class=\"p\">()</span><span
  class=\"o\">.</span><span class=\"n\">split</span><span class=\"p\">(</span><span
  class=\"s2\">\" \"</span><span class=\"p\">,</span> <span class=\"mi\">1</span><span
  class=\"p\">)[</span><span class=\"o\">-</span><span class=\"mi\">1</span><span
  class=\"p\">]</span> <span class=\"ow\">or</span> <span class=\"s2\">\"\"</span>\n
  \           <span class=\"k\">return</span> <span class=\"nb\">bool</span><span
  class=\"p\">(</span><span class=\"n\">tag</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"details\" style=\"margin:0;padding:.5rem
  1rem;\">details <em class=\"small\">function</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"details
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
  <span class=\"nf\">details</span><span class=\"p\">(</span><span class=\"n\">state</span><span
  class=\"p\">:</span> <span class=\"n\">StateBlock</span><span class=\"p\">,</span>
  <span class=\"n\">startLine</span><span class=\"p\">:</span> <span class=\"nb\">int</span><span
  class=\"p\">,</span> <span class=\"n\">endLine</span><span class=\"p\">:</span>
  <span class=\"nb\">int</span><span class=\"p\">,</span> <span class=\"n\">silent</span><span
  class=\"p\">:</span> <span class=\"nb\">bool</span><span class=\"p\">)</span> <span
  class=\"o\">-&gt;</span> <span class=\"nb\">bool</span><span class=\"p\">:</span>\n
  \           <span class=\"n\">start</span> <span class=\"o\">=</span> <span class=\"n\">state</span><span
  class=\"o\">.</span><span class=\"n\">bMarks</span><span class=\"p\">[</span><span
  class=\"n\">startLine</span><span class=\"p\">]</span> <span class=\"o\">+</span>
  <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">tShift</span><span
  class=\"p\">[</span><span class=\"n\">startLine</span><span class=\"p\">]</span>\n
  \           <span class=\"n\">maximum</span> <span class=\"o\">=</span> <span class=\"n\">state</span><span
  class=\"o\">.</span><span class=\"n\">eMarks</span><span class=\"p\">[</span><span
  class=\"n\">startLine</span><span class=\"p\">]</span>\n\n            <span class=\"c1\">#
  Check out the first character quickly, which should filter out most of non-containers</span>\n
  \           <span class=\"k\">if</span> <span class=\"nb\">ord</span><span class=\"p\">(</span><span
  class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">src</span><span
  class=\"p\">[</span><span class=\"n\">start</span><span class=\"p\">])</span> <span
  class=\"o\">!=</span> <span class=\"n\">MARKER_CHAR</span><span class=\"p\">:</span>\n
  \               <span class=\"k\">return</span> <span class=\"kc\">False</span>\n\n
  \           <span class=\"c1\"># Check out the rest of the marker string</span>\n
  \           <span class=\"n\">pos</span> <span class=\"o\">=</span> <span class=\"n\">start</span>
  <span class=\"o\">+</span> <span class=\"mi\">1</span>\n            <span class=\"k\">while</span>
  <span class=\"n\">pos</span> <span class=\"o\">&lt;=</span> <span class=\"n\">maximum</span>
  <span class=\"ow\">and</span> <span class=\"n\">MARKER_STR</span><span class=\"p\">[(</span><span
  class=\"n\">pos</span> <span class=\"o\">-</span> <span class=\"n\">start</span><span
  class=\"p\">)</span> <span class=\"o\">%</span> <span class=\"n\">MARKER_LEN</span><span
  class=\"p\">]</span> <span class=\"o\">==</span> <span class=\"n\">state</span><span
  class=\"o\">.</span><span class=\"n\">src</span><span class=\"p\">[</span><span
  class=\"n\">pos</span><span class=\"p\">]:</span>\n                <span class=\"n\">pos</span>
  <span class=\"o\">+=</span> <span class=\"mi\">1</span>\n\n            <span class=\"n\">marker_count</span>
  <span class=\"o\">=</span> <span class=\"n\">math</span><span class=\"o\">.</span><span
  class=\"n\">floor</span><span class=\"p\">((</span><span class=\"n\">pos</span>
  <span class=\"o\">-</span> <span class=\"n\">start</span><span class=\"p\">)</span>
  <span class=\"o\">/</span> <span class=\"n\">MARKER_LEN</span><span class=\"p\">)</span>\n
  \           <span class=\"k\">if</span> <span class=\"n\">marker_count</span> <span
  class=\"o\">&lt;</span> <span class=\"n\">MIN_MARKERS</span><span class=\"p\">:</span>\n
  \               <span class=\"k\">return</span> <span class=\"kc\">False</span>\n
  \           <span class=\"n\">marker_pos</span> <span class=\"o\">=</span> <span
  class=\"n\">pos</span> <span class=\"o\">-</span> <span class=\"p\">((</span><span
  class=\"n\">pos</span> <span class=\"o\">-</span> <span class=\"n\">start</span><span
  class=\"p\">)</span> <span class=\"o\">%</span> <span class=\"n\">MARKER_LEN</span><span
  class=\"p\">)</span>\n            <span class=\"n\">params</span> <span class=\"o\">=</span>
  <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">src</span><span
  class=\"p\">[</span><span class=\"n\">marker_pos</span><span class=\"p\">:</span><span
  class=\"n\">maximum</span><span class=\"p\">]</span>\n            <span class=\"n\">markup</span>
  <span class=\"o\">=</span> <span class=\"n\">state</span><span class=\"o\">.</span><span
  class=\"n\">src</span><span class=\"p\">[</span><span class=\"n\">start</span><span
  class=\"p\">:</span><span class=\"n\">marker_pos</span><span class=\"p\">]</span>\n\n
  \           <span class=\"k\">if</span> <span class=\"ow\">not</span> <span class=\"n\">validate</span><span
  class=\"p\">(</span><span class=\"n\">params</span><span class=\"p\">):</span>\n
  \               <span class=\"k\">return</span> <span class=\"kc\">False</span>\n\n
  \           <span class=\"c1\"># Since start is found, we can report success here
  in validation mode</span>\n            <span class=\"k\">if</span> <span class=\"n\">silent</span><span
  class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"kc\">True</span>\n\n
  \           <span class=\"n\">old_parent</span> <span class=\"o\">=</span> <span
  class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">parentType</span>\n
  \           <span class=\"n\">old_line_max</span> <span class=\"o\">=</span> <span
  class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">lineMax</span>\n
  \           <span class=\"n\">old_indent</span> <span class=\"o\">=</span> <span
  class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">blkIndent</span>\n\n
  \           <span class=\"n\">blk_start</span> <span class=\"o\">=</span> <span
  class=\"n\">pos</span>\n            <span class=\"k\">while</span> <span class=\"n\">blk_start</span>
  <span class=\"o\">&lt;</span> <span class=\"n\">maximum</span> <span class=\"ow\">and</span>
  <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">src</span><span
  class=\"p\">[</span><span class=\"n\">blk_start</span><span class=\"p\">]</span>
  <span class=\"o\">==</span> <span class=\"s2\">\" \"</span><span class=\"p\">:</span>\n
  \               <span class=\"n\">blk_start</span> <span class=\"o\">+=</span> <span
  class=\"mi\">1</span>\n\n            <span class=\"n\">state</span><span class=\"o\">.</span><span
  class=\"n\">parentType</span> <span class=\"o\">=</span> <span class=\"s2\">\"details\"</span>\n
  \           <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">blkIndent</span>
  <span class=\"o\">+=</span> <span class=\"n\">blk_start</span> <span class=\"o\">-</span>
  <span class=\"n\">start</span>\n\n            <span class=\"n\">was_empty</span>
  <span class=\"o\">=</span> <span class=\"kc\">False</span>\n\n            <span
  class=\"c1\"># Search for the end of the block</span>\n            <span class=\"n\">next_line</span>
  <span class=\"o\">=</span> <span class=\"n\">startLine</span>\n            <span
  class=\"k\">while</span> <span class=\"kc\">True</span><span class=\"p\">:</span>\n
  \               <span class=\"n\">next_line</span> <span class=\"o\">+=</span> <span
  class=\"mi\">1</span>\n                <span class=\"k\">if</span> <span class=\"n\">next_line</span>
  <span class=\"o\">&gt;=</span> <span class=\"n\">endLine</span><span class=\"p\">:</span>\n
  \                   <span class=\"c1\"># unclosed block should be autoclosed by
  end of document.</span>\n                    <span class=\"c1\"># also block seems
  to be autoclosed by end of parent</span>\n                    <span class=\"k\">break</span>\n
  \               <span class=\"n\">pos</span> <span class=\"o\">=</span> <span class=\"n\">state</span><span
  class=\"o\">.</span><span class=\"n\">bMarks</span><span class=\"p\">[</span><span
  class=\"n\">next_line</span><span class=\"p\">]</span> <span class=\"o\">+</span>
  <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">tShift</span><span
  class=\"p\">[</span><span class=\"n\">next_line</span><span class=\"p\">]</span>\n
  \               <span class=\"n\">maximum</span> <span class=\"o\">=</span> <span
  class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">eMarks</span><span
  class=\"p\">[</span><span class=\"n\">next_line</span><span class=\"p\">]</span>\n
  \               <span class=\"n\">is_empty</span> <span class=\"o\">=</span> <span
  class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">sCount</span><span
  class=\"p\">[</span><span class=\"n\">next_line</span><span class=\"p\">]</span>
  <span class=\"o\">&lt;</span> <span class=\"n\">state</span><span class=\"o\">.</span><span
  class=\"n\">blkIndent</span>\n\n                <span class=\"c1\"># two consecutive
  empty lines autoclose the block</span>\n                <span class=\"k\">if</span>
  <span class=\"n\">is_empty</span> <span class=\"ow\">and</span> <span class=\"n\">was_empty</span><span
  class=\"p\">:</span>\n                    <span class=\"k\">break</span>\n                <span
  class=\"n\">was_empty</span> <span class=\"o\">=</span> <span class=\"n\">is_empty</span>\n\n
  \               <span class=\"k\">if</span> <span class=\"n\">pos</span> <span class=\"o\">&lt;</span>
  <span class=\"n\">maximum</span> <span class=\"ow\">and</span> <span class=\"n\">state</span><span
  class=\"o\">.</span><span class=\"n\">sCount</span><span class=\"p\">[</span><span
  class=\"n\">next_line</span><span class=\"p\">]</span> <span class=\"o\">&lt;</span>
  <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">blkIndent</span><span
  class=\"p\">:</span>\n                    <span class=\"c1\"># non-empty line with
  negative indent should stop the block:</span>\n                    <span class=\"c1\">#
  - !!!</span>\n                    <span class=\"c1\">#  test</span>\n                    <span
  class=\"k\">break</span>\n\n            <span class=\"c1\"># this will prevent lazy
  continuations from ever going past our end marker</span>\n            <span class=\"n\">state</span><span
  class=\"o\">.</span><span class=\"n\">lineMax</span> <span class=\"o\">=</span>
  <span class=\"n\">next_line</span>\n\n            <span class=\"n\">tag</span><span
  class=\"p\">,</span> <span class=\"n\">title</span><span class=\"p\">,</span> <span
  class=\"nb\">open</span> <span class=\"o\">=</span> <span class=\"n\">get_tag</span><span
  class=\"p\">(</span><span class=\"n\">params</span><span class=\"p\">)</span>\n
  \           <span class=\"n\">attrs</span> <span class=\"o\">=</span> <span class=\"p\">{</span><span
  class=\"s2\">\"class\"</span><span class=\"p\">:</span> <span class=\"sa\">f</span><span
  class=\"s2\">\"details </span><span class=\"si\">{</span><span class=\"n\">tag</span><span
  class=\"si\">}</span><span class=\"s2\">\"</span><span class=\"p\">}</span>\n            <span
  class=\"k\">if</span> <span class=\"nb\">open</span><span class=\"p\">:</span>\n
  \               <span class=\"n\">attrs</span><span class=\"p\">[</span><span class=\"s2\">\"open\"</span><span
  class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"s2\">\"\"</span>\n\n
  \           <span class=\"n\">token</span> <span class=\"o\">=</span> <span class=\"n\">state</span><span
  class=\"o\">.</span><span class=\"n\">push</span><span class=\"p\">(</span><span
  class=\"s2\">\"details_open\"</span><span class=\"p\">,</span> <span class=\"s2\">\"details\"</span><span
  class=\"p\">,</span> <span class=\"mi\">1</span><span class=\"p\">)</span>\n            <span
  class=\"n\">token</span><span class=\"o\">.</span><span class=\"n\">markup</span>
  <span class=\"o\">=</span> <span class=\"n\">markup</span>\n            <span class=\"n\">token</span><span
  class=\"o\">.</span><span class=\"n\">block</span> <span class=\"o\">=</span> <span
  class=\"kc\">True</span>\n            <span class=\"n\">token</span><span class=\"o\">.</span><span
  class=\"n\">attrs</span> <span class=\"o\">=</span> <span class=\"n\">attrs</span>\n
  \           <span class=\"n\">token</span><span class=\"o\">.</span><span class=\"n\">meta</span>
  <span class=\"o\">=</span> <span class=\"p\">{</span><span class=\"s2\">\"tag\"</span><span
  class=\"p\">:</span> <span class=\"n\">tag</span><span class=\"p\">}</span>\n            <span
  class=\"n\">token</span><span class=\"o\">.</span><span class=\"n\">content</span>
  <span class=\"o\">=</span> <span class=\"n\">title</span>\n            <span class=\"n\">token</span><span
  class=\"o\">.</span><span class=\"n\">info</span> <span class=\"o\">=</span> <span
  class=\"n\">params</span>\n            <span class=\"n\">token</span><span class=\"o\">.</span><span
  class=\"n\">map</span> <span class=\"o\">=</span> <span class=\"p\">[</span><span
  class=\"n\">startLine</span><span class=\"p\">,</span> <span class=\"n\">next_line</span><span
  class=\"p\">]</span>\n\n            <span class=\"k\">if</span> <span class=\"n\">title</span><span
  class=\"p\">:</span>\n                <span class=\"n\">title_markup</span> <span
  class=\"o\">=</span> <span class=\"sa\">f</span><span class=\"s2\">\"</span><span
  class=\"si\">{</span><span class=\"n\">markup</span><span class=\"si\">}</span><span
  class=\"s2\"> </span><span class=\"si\">{</span><span class=\"n\">tag</span><span
  class=\"si\">}</span><span class=\"s2\">\"</span>\n                <span class=\"n\">token</span>
  <span class=\"o\">=</span> <span class=\"n\">state</span><span class=\"o\">.</span><span
  class=\"n\">push</span><span class=\"p\">(</span><span class=\"s2\">\"details_title_open\"</span><span
  class=\"p\">,</span> <span class=\"s2\">\"summary\"</span><span class=\"p\">,</span>
  <span class=\"mi\">1</span><span class=\"p\">)</span>\n                <span class=\"n\">token</span><span
  class=\"o\">.</span><span class=\"n\">markup</span> <span class=\"o\">=</span> <span
  class=\"n\">title_markup</span>\n                <span class=\"n\">token</span><span
  class=\"o\">.</span><span class=\"n\">attrs</span> <span class=\"o\">=</span> <span
  class=\"p\">{</span><span class=\"s2\">\"class\"</span><span class=\"p\">:</span>
  <span class=\"s2\">\"admonition-title\"</span><span class=\"p\">}</span>\n                <span
  class=\"n\">token</span><span class=\"o\">.</span><span class=\"n\">map</span> <span
  class=\"o\">=</span> <span class=\"p\">[</span><span class=\"n\">startLine</span><span
  class=\"p\">,</span> <span class=\"n\">startLine</span> <span class=\"o\">+</span>
  <span class=\"mi\">1</span><span class=\"p\">]</span>\n\n                <span class=\"n\">token</span>
  <span class=\"o\">=</span> <span class=\"n\">state</span><span class=\"o\">.</span><span
  class=\"n\">push</span><span class=\"p\">(</span><span class=\"s2\">\"inline\"</span><span
  class=\"p\">,</span> <span class=\"s2\">\"\"</span><span class=\"p\">,</span> <span
  class=\"mi\">0</span><span class=\"p\">)</span>\n                <span class=\"n\">token</span><span
  class=\"o\">.</span><span class=\"n\">content</span> <span class=\"o\">=</span>
  <span class=\"n\">title</span>\n                <span class=\"n\">token</span><span
  class=\"o\">.</span><span class=\"n\">map</span> <span class=\"o\">=</span> <span
  class=\"p\">[</span><span class=\"n\">startLine</span><span class=\"p\">,</span>
  <span class=\"n\">startLine</span> <span class=\"o\">+</span> <span class=\"mi\">1</span><span
  class=\"p\">]</span>\n                <span class=\"n\">token</span><span class=\"o\">.</span><span
  class=\"n\">children</span> <span class=\"o\">=</span> <span class=\"p\">[]</span>\n\n
  \               <span class=\"n\">token</span> <span class=\"o\">=</span> <span
  class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">push</span><span
  class=\"p\">(</span><span class=\"s2\">\"details_title_close\"</span><span class=\"p\">,</span>
  <span class=\"s2\">\"summary\"</span><span class=\"p\">,</span> <span class=\"o\">-</span><span
  class=\"mi\">1</span><span class=\"p\">)</span>\n                <span class=\"n\">token</span><span
  class=\"o\">.</span><span class=\"n\">markup</span> <span class=\"o\">=</span> <span
  class=\"n\">title_markup</span>\n\n            <span class=\"n\">state</span><span
  class=\"o\">.</span><span class=\"n\">md</span><span class=\"o\">.</span><span class=\"n\">block</span><span
  class=\"o\">.</span><span class=\"n\">tokenize</span><span class=\"p\">(</span><span
  class=\"n\">state</span><span class=\"p\">,</span> <span class=\"n\">startLine</span>
  <span class=\"o\">+</span> <span class=\"mi\">1</span><span class=\"p\">,</span>
  <span class=\"n\">next_line</span><span class=\"p\">)</span>\n\n            <span
  class=\"n\">token</span> <span class=\"o\">=</span> <span class=\"n\">state</span><span
  class=\"o\">.</span><span class=\"n\">push</span><span class=\"p\">(</span><span
  class=\"s2\">\"details_close\"</span><span class=\"p\">,</span> <span class=\"s2\">\"details\"</span><span
  class=\"p\">,</span> <span class=\"o\">-</span><span class=\"mi\">1</span><span
  class=\"p\">)</span>\n            <span class=\"n\">token</span><span class=\"o\">.</span><span
  class=\"n\">markup</span> <span class=\"o\">=</span> <span class=\"n\">state</span><span
  class=\"o\">.</span><span class=\"n\">src</span><span class=\"p\">[</span><span
  class=\"n\">start</span><span class=\"p\">:</span><span class=\"n\">pos</span><span
  class=\"p\">]</span>\n            <span class=\"n\">token</span><span class=\"o\">.</span><span
  class=\"n\">block</span> <span class=\"o\">=</span> <span class=\"kc\">True</span>\n\n
  \           <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">parentType</span>
  <span class=\"o\">=</span> <span class=\"n\">old_parent</span>\n            <span
  class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">lineMax</span>
  <span class=\"o\">=</span> <span class=\"n\">old_line_max</span>\n            <span
  class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">blkIndent</span>
  <span class=\"o\">=</span> <span class=\"n\">old_indent</span>\n            <span
  class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">line</span>
  <span class=\"o\">=</span> <span class=\"n\">next_line</span>\n\n            <span
  class=\"k\">return</span> <span class=\"kc\">True</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"details_plugin\" style=\"margin:0;padding:.5rem
  1rem;\">details_plugin <em class=\"small\">function</em></h2>\nPlugin to use\n<code>python-markdown
  style detailss     &lt;https://python-markdown.github.io/extensions/details&gt;</code>_.\n<pre><code>..
  code-block:: md\n\n    ??? note\n        *content*\n\nNote, this is ported from\n`markdown-it-admon\n&lt;https://github.com/commenthol/markdown-it-admon&gt;`_.\n</code></pre>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"details_plugin
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
  <span class=\"nf\">details_plugin</span><span class=\"p\">(</span><span class=\"n\">md</span><span
  class=\"p\">:</span> <span class=\"n\">MarkdownIt</span><span class=\"p\">,</span>
  <span class=\"n\">render</span><span class=\"p\">:</span> <span class=\"n\">Optional</span><span
  class=\"p\">[</span><span class=\"n\">Callable</span><span class=\"p\">]</span>
  <span class=\"o\">=</span> <span class=\"kc\">None</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
  class=\"w\">            </span><span class=\"sd\">\"\"\"Plugin to use</span>\n<span
  class=\"sd\">            `python-markdown style detailss</span>\n<span class=\"sd\">
  \           &lt;https://python-markdown.github.io/extensions/details&gt;`_.</span>\n\n<span
  class=\"sd\">            .. code-block:: md</span>\n\n<span class=\"sd\">                ???
  note</span>\n<span class=\"sd\">                    *content*</span>\n\n<span class=\"sd\">
  \           Note, this is ported from</span>\n<span class=\"sd\">            `markdown-it-admon</span>\n<span
  class=\"sd\">            &lt;https://github.com/commenthol/markdown-it-admon&gt;`_.</span>\n<span
  class=\"sd\">            \"\"\"</span>\n\n            <span class=\"k\">def</span>
  <span class=\"nf\">renderDefault</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">,</span> <span class=\"n\">tokens</span><span class=\"p\">,</span> <span
  class=\"n\">idx</span><span class=\"p\">,</span> <span class=\"n\">_options</span><span
  class=\"p\">,</span> <span class=\"n\">env</span><span class=\"p\">):</span>\n                <span
  class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">renderToken</span><span class=\"p\">(</span><span class=\"n\">tokens</span><span
  class=\"p\">,</span> <span class=\"n\">idx</span><span class=\"p\">,</span> <span
  class=\"n\">_options</span><span class=\"p\">,</span> <span class=\"n\">env</span><span
  class=\"p\">)</span>\n\n            <span class=\"n\">render</span> <span class=\"o\">=</span>
  <span class=\"n\">render</span> <span class=\"ow\">or</span> <span class=\"n\">renderDefault</span>\n\n
  \           <span class=\"n\">md</span><span class=\"o\">.</span><span class=\"n\">add_render_rule</span><span
  class=\"p\">(</span><span class=\"s2\">\"details_open\"</span><span class=\"p\">,</span>
  <span class=\"n\">render</span><span class=\"p\">)</span>\n            <span class=\"n\">md</span><span
  class=\"o\">.</span><span class=\"n\">add_render_rule</span><span class=\"p\">(</span><span
  class=\"s2\">\"details_close\"</span><span class=\"p\">,</span> <span class=\"n\">render</span><span
  class=\"p\">)</span>\n            <span class=\"n\">md</span><span class=\"o\">.</span><span
  class=\"n\">add_render_rule</span><span class=\"p\">(</span><span class=\"s2\">\"details_title_open\"</span><span
  class=\"p\">,</span> <span class=\"n\">render</span><span class=\"p\">)</span>\n
  \           <span class=\"n\">md</span><span class=\"o\">.</span><span class=\"n\">add_render_rule</span><span
  class=\"p\">(</span><span class=\"s2\">\"details_title_close\"</span><span class=\"p\">,</span>
  <span class=\"n\">render</span><span class=\"p\">)</span>\n\n            <span class=\"n\">md</span><span
  class=\"o\">.</span><span class=\"n\">block</span><span class=\"o\">.</span><span
  class=\"n\">ruler</span><span class=\"o\">.</span><span class=\"n\">before</span><span
  class=\"p\">(</span>\n                <span class=\"s2\">\"fence\"</span><span class=\"p\">,</span>\n
  \               <span class=\"s2\">\"details\"</span><span class=\"p\">,</span>\n
  \               <span class=\"n\">details</span><span class=\"p\">,</span>\n                <span
  class=\"p\">{</span><span class=\"s2\">\"alt\"</span><span class=\"p\">:</span>
  <span class=\"p\">[</span><span class=\"s2\">\"paragraph\"</span><span class=\"p\">,</span>
  <span class=\"s2\">\"reference\"</span><span class=\"p\">,</span> <span class=\"s2\">\"blockquote\"</span><span
  class=\"p\">,</span> <span class=\"s2\">\"list\"</span><span class=\"p\">]},</span>\n
  \           <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!! function </p><h2
  class=\"admonition-title\" id=\"renderDefault\" style=\"margin:0;padding:.5rem 1rem;\">renderDefault
  <em class=\"small\">function</em></h2>\n<div class=\"admonition source is-collapsible
  collapsible-open\">\n<p class=\"admonition-title\">\"renderDefault <em class=\"small\">source</em>\"</p>\n</div>\n<pre
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
  <span class=\"nf\">renderDefault</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">,</span> <span class=\"n\">tokens</span><span class=\"p\">,</span> <span
  class=\"n\">idx</span><span class=\"p\">,</span> <span class=\"n\">_options</span><span
  class=\"p\">,</span> <span class=\"n\">env</span><span class=\"p\">):</span>\n                <span
  class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">renderToken</span><span class=\"p\">(</span><span class=\"n\">tokens</span><span
  class=\"p\">,</span> <span class=\"n\">idx</span><span class=\"p\">,</span> <span
  class=\"n\">_options</span><span class=\"p\">,</span> <span class=\"n\">env</span><span
  class=\"p\">)</span>\n</pre></div>\n\n</pre>\n</main>\n<footer>\xA9 2024</footer>\n</body></html>"
published: true
slug: markata/plugins/mdit-details
title: Mdit_Details.Py


---

None


!! function <h2 id='get_tag' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get_tag <em class='small'>function</em></h2>

???+ source "get_tag <em class='small'>source</em>"

```python

        def get_tag(params: str) -> Tuple[str, str, bool]:
            open = False
            if not params.strip():
                return "", "", open

            if params.strip().startswith("+"):
                params = params.strip("+")
                open = True

            tag, *_title = params.strip().split(" ")
            joined = " ".join(_title).strip('"').strip("'")

            title = ""
            if not joined:
                title = tag.title()
            elif joined != '""':
                title = joined
            return (tag.lower(), title, open)
```


!! function <h2 id='validate' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>validate <em class='small'>function</em></h2>

???+ source "validate <em class='small'>source</em>"

```python

        def validate(params: str) -> bool:
            tag = params.strip().split(" ", 1)[-1] or ""
            return bool(tag)
```


!! function <h2 id='details' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>details <em class='small'>function</em></h2>

???+ source "details <em class='small'>source</em>"

```python

        def details(state: StateBlock, startLine: int, endLine: int, silent: bool) -> bool:
            start = state.bMarks[startLine] + state.tShift[startLine]
            maximum = state.eMarks[startLine]

            # Check out the first character quickly, which should filter out most of non-containers
            if ord(state.src[start]) != MARKER_CHAR:
                return False

            # Check out the rest of the marker string
            pos = start + 1
            while pos <= maximum and MARKER_STR[(pos - start) % MARKER_LEN] == state.src[pos]:
                pos += 1

            marker_count = math.floor((pos - start) / MARKER_LEN)
            if marker_count < MIN_MARKERS:
                return False
            marker_pos = pos - ((pos - start) % MARKER_LEN)
            params = state.src[marker_pos:maximum]
            markup = state.src[start:marker_pos]

            if not validate(params):
                return False

            # Since start is found, we can report success here in validation mode
            if silent:
                return True

            old_parent = state.parentType
            old_line_max = state.lineMax
            old_indent = state.blkIndent

            blk_start = pos
            while blk_start < maximum and state.src[blk_start] == " ":
                blk_start += 1

            state.parentType = "details"
            state.blkIndent += blk_start - start

            was_empty = False

            # Search for the end of the block
            next_line = startLine
            while True:
                next_line += 1
                if next_line >= endLine:
                    # unclosed block should be autoclosed by end of document.
                    # also block seems to be autoclosed by end of parent
                    break
                pos = state.bMarks[next_line] + state.tShift[next_line]
                maximum = state.eMarks[next_line]
                is_empty = state.sCount[next_line] < state.blkIndent

                # two consecutive empty lines autoclose the block
                if is_empty and was_empty:
                    break
                was_empty = is_empty

                if pos < maximum and state.sCount[next_line] < state.blkIndent:
                    # non-empty line with negative indent should stop the block:
                    # - !!!
                    #  test
                    break

            # this will prevent lazy continuations from ever going past our end marker
            state.lineMax = next_line

            tag, title, open = get_tag(params)
            attrs = {"class": f"details {tag}"}
            if open:
                attrs["open"] = ""

            token = state.push("details_open", "details", 1)
            token.markup = markup
            token.block = True
            token.attrs = attrs
            token.meta = {"tag": tag}
            token.content = title
            token.info = params
            token.map = [startLine, next_line]

            if title:
                title_markup = f"{markup} {tag}"
                token = state.push("details_title_open", "summary", 1)
                token.markup = title_markup
                token.attrs = {"class": "admonition-title"}
                token.map = [startLine, startLine + 1]

                token = state.push("inline", "", 0)
                token.content = title
                token.map = [startLine, startLine + 1]
                token.children = []

                token = state.push("details_title_close", "summary", -1)
                token.markup = title_markup

            state.md.block.tokenize(state, startLine + 1, next_line)

            token = state.push("details_close", "details", -1)
            token.markup = state.src[start:pos]
            token.block = True

            state.parentType = old_parent
            state.lineMax = old_line_max
            state.blkIndent = old_indent
            state.line = next_line

            return True
```


!! function <h2 id='details_plugin' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>details_plugin <em class='small'>function</em></h2>
    Plugin to use
    `python-markdown style detailss
    <https://python-markdown.github.io/extensions/details>`_.

    .. code-block:: md

        ??? note
            *content*

    Note, this is ported from
    `markdown-it-admon
    <https://github.com/commenthol/markdown-it-admon>`_.
???+ source "details_plugin <em class='small'>source</em>"

```python

        def details_plugin(md: MarkdownIt, render: Optional[Callable] = None) -> None:
            """Plugin to use
            `python-markdown style detailss
            <https://python-markdown.github.io/extensions/details>`_.

            .. code-block:: md

                ??? note
                    *content*

            Note, this is ported from
            `markdown-it-admon
            <https://github.com/commenthol/markdown-it-admon>`_.
            """

            def renderDefault(self, tokens, idx, _options, env):
                return self.renderToken(tokens, idx, _options, env)

            render = render or renderDefault

            md.add_render_rule("details_open", render)
            md.add_render_rule("details_close", render)
            md.add_render_rule("details_title_open", render)
            md.add_render_rule("details_title_close", render)

            md.block.ruler.before(
                "fence",
                "details",
                details,
                {"alt": ["paragraph", "reference", "blockquote", "list"]},
            )
```


!! function <h2 id='renderDefault' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>renderDefault <em class='small'>function</em></h2>

???+ source "renderDefault <em class='small'>source</em>"

```python

        def renderDefault(self, tokens, idx, _options, env):
                return self.renderToken(tokens, idx, _options, env)
```
