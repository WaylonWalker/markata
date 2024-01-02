---
content: "### Configuration\n\nExample configuration.  Covers supports multiple covers
  to be configured.  Here\nis an example from my blog where we have a template sized
  for dev.to and one\nsized for open graph.  Each image takes it's own configuration.\n\n```
  toml\n[[markata.covers]]\nname='-dev'\ntemplate = \"static/cover-template.png\"\nfont
  = \"./static/JosefinSans-Regular.ttf\"\ntext_font = \"./static/JosefinSans-Regular.ttf\"\nfont_color
  = \"rgb(185,155,165)\"\ntext_font_color = \"rgb(255,255,255)\"\ntext_key = 'description'\npadding
  = [0, 40, 100, 300]\ntext_padding = [0,0]\n\n[[markata.covers]]\nname=''\ntemplate
  = \"static/og-template.png\"\nfont = \"./static/JosefinSans-Regular.ttf\"\nfont_color
  = \"rgb(255,255,255)\"\ntext_font = \"./static/JosefinSans-Regular.ttf\"\ntext_font_color
  = \"rgb(200,200,200)\"\ntext_key = 'description'\npadding = [10, 10, 100, 300]\ntext_padding
  = [0,0]\n```\n\n\n!! function <h2 id='_load_font' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>_load_font <em class='small'>function</em></h2>\n\n???+ source \"_load_font
  <em class='small'>source</em>\"\n\n```python\n\n        def _load_font(path: Path,
  size: int) -> ImageFont.FreeTypeFont:\n            return ImageFont.truetype(path,
  size=size)\n```\n\n\n!! function <h2 id='get_font' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>get_font <em class='small'>function</em></h2>\n\n???+ source \"get_font <em
  class='small'>source</em>\"\n\n```python\n\n        def get_font(\n            path:
  Path,\n            draw: ImageDraw.Draw,\n            title: str,\n            size:
  int = 250,\n            max_size: tuple = (800, 220),\n        ) -> ImageFont.FreeTypeFont:\n
  \           title = title or \"\"\n            font = _load_font(path, size)\n            current_size
  = draw.textsize(title, font=font)\n\n            if current_size[0] > max_size[0]
  or current_size[1] > max_size[1]:\n                return get_font(path, draw, title,
  size - 10, max_size=max_size)\n            return font\n```\n\n\n!! class <h2 id='PaddingError'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>PaddingError <em class='small'>class</em></h2>\n\n???+
  source \"PaddingError <em class='small'>source</em>\"\n\n```python\n\n        class
  PaddingError(BaseException):\n            def __init__(\n                self,\n
  \               msg: str = \"\",\n            ) -> None:\n                super().__init__(\n
  \                   \"Padding must be an iterable of length 1, 2, 3, or 4.\\n\"
  + msg,\n                )\n```\n\n\n!! function <h2 id='draw_text' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>draw_text <em class='small'>function</em></h2>\n\n???+
  source \"draw_text <em class='small'>source</em>\"\n\n```python\n\n        def draw_text(\n
  \           image: Image,\n            font_path: Optional[Path],\n            text:
  str,\n            color: Union[str, None],\n            padding: Tuple[int, ...],\n
  \           markata: \"Markata\",\n        ) -> None:\n            text = text or
  \"\"\n            draw = ImageDraw.Draw(image)\n            padding = resolve_padding(padding,
  markata)\n            width = image.size[0]\n            height = image.size[1]\n
  \           bounding_box = [padding[0], padding[1], width - padding[0], height -
  padding[1]]\n            bounding_box = [padding[0], padding[1], width - padding[2],
  height - padding[3]]\n            max_size = (bounding_box[2] - bounding_box[0],
  bounding_box[3] - bounding_box[1])\n            x1, y1, x2, y2 = bounding_box\n
  \           font = get_font(font_path, draw, text, max_size=max_size) if font_path
  else None\n            w, h = draw.textsize(text, font=font)\n            x = (x2
  - x1 - w) / 2 + x1\n            y = (y2 - y1 - h) / 2 + y1\n            draw.text((x,
  y), text, fill=color, font=font, align=\"center\")\n```\n\n\n!! function <h2 id='resolve_padding'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>resolve_padding <em
  class='small'>function</em></h2>\n    Convert padding to a len 4 tuple\n???+ source
  \"resolve_padding <em class='small'>source</em>\"\n\n```python\n\n        def resolve_padding(padding:
  Tuple[int, ...], markata: \"Markata\") -> Tuple[int, ...]:\n            \"\"\"Convert
  padding to a len 4 tuple\"\"\"\n            if len(padding) == 4:\n                return
  padding\n            if len(padding) == 3:\n                return (*padding, padding[1])\n
  \           if len(padding) == 2:\n                return padding * 2\n            if
  len(padding) == 1:\n                return padding * 4\n            raise PaddingError(f\"recieved
  padding: {padding}\")\n```\n\n\n!! function <h2 id='make_cover' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>make_cover <em class='small'>function</em></h2>\n\n???+
  source \"make_cover <em class='small'>source</em>\"\n\n```python\n\n        def
  make_cover(\n            title: str,\n            color: str,\n            output_path:
  Path,\n            template_path: Path,\n            font_path: Optional[Path],\n
  \           padding: Tuple[int, ...],\n            text_font: Path,\n            text:
  str = None,\n            text_font_color: str = None,\n            text_padding:
  Tuple[int, ...] = None,\n            resizes: List[int] = None,\n            markata:
  \"Markata\" = None,\n        ) -> None:\n            if output_path.exists():\n
  \               return\n            image = Image.open(template_path) if template_path
  else Image.new(\"RGB\", (800, 450))\n\n            draw_text(\n                image=image,\n
  \               font_path=font_path,\n                title=title,\n                color=color,\n
  \               padding=padding,\n                markata=markata,\n            )\n
  \           if text is not None:\n                if text_padding is None:\n                    text_padding
  = (\n                        image.size[1] - image.size[1] / 5,\n                        image.size[0]
  / 5,\n                        image.size[1] - image.size[1] / 10,\n                    )\n
  \               draw_text(image, text_font, text, text_font_color, text_padding)\n\n
  \           image.save(output_path)\n            ratio = image.size[1] / image.size[0]\n\n
  \           covers = []\n            if resizes:\n                for width in resizes:\n
  \                   re_img = image.resize((width, int(width * ratio)), Image.ANTIALIAS)\n
  \                   filename = (\n                        f\"{output_path.stem}_{width}x{int(width*ratio)}{output_path.suffix}\"\n
  \                   )\n                    covers.append(filename)\n\n                    filepath
  = Path(output_path.parent / filename)\n                    re_img.save(filepath)\n```\n\n\n!!
  function <h2 id='save' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>save
  <em class='small'>function</em></h2>\n\n???+ source \"save <em class='small'>source</em>\"\n\n```python\n\n
  \       def save(markata: \"Markata\") -> None:\n            futures = []\n\n            if
  \"covers\" not in markata.config.keys():\n                return\n\n            for
  article in markata.iter_articles(\"making covers\"):\n                for cover
  in markata.config[\"covers\"]:\n                    try:\n                        padding
  = cover[\"padding\"]\n                    except KeyError:\n                        padding
  = (\n                            200,\n                            100,\n                        )\n
  \                   try:\n                        text_padding = cover[\"text_padding\"]\n
  \                   except KeyError:\n                        text_padding = (\n
  \                           200,\n                            100,\n                        )\n
  \                   if \"text_key\" in cover:\n                        try:\n                            text
  = article.metadata[cover[\"text_key\"]]\n                        except AttributeError:\n
  \                           text = article[cover[\"text_key\"]]\n                        try:\n
  \                           text = text.replace(\"\\n\", \"\")\n                            from
  more_itertools import chunked\n\n                            text = \"\\n\".join([\"\".join(c)
  for c in chunked(text, 60)])\n                        except AttributeError:\n                            #
  text is likely None\n                            pass\n\n                        text_font
  = cover[\"text_font\"]\n                        text_font_color = cover[\"text_font_color\"]\n
  \                   else:\n                        text = None\n                        text_font
  = None\n                        text_font_color = None\n                    try:\n
  \                       title = article.metadata[\"title\"]\n                    except
  AttributeError:\n                        title = article[\"title\"]\n                    futures.append(\n
  \                       make_cover(\n                            title=title,\n
  \                           color=cover[\"font_color\"],\n                            output_path=Path(markata.config.output_dir)\n
  \                           / (article[\"slug\"] + cover[\"name\"] + \".png\"),\n
  \                           template_path=cover.get(\"template\", None),\n                            font_path=cover.get(\"font\",
  None),\n                            padding=padding,\n                            text_font=text_font,\n
  \                           text=text,\n                            text_font_color=text_font_color,\n
  \                           text_padding=text_padding,\n                            resizes=cover.get(\"resizes\"),\n
  \                           markata=markata,\n                        ),\n                    )\n\n
  \           progress = Progress(\n                BarColumn(bar_width=None),\n                transient=True,\n
  \               console=markata.console,\n            )\n            task_id = progress.add_task(\"loading
  markdown\")\n            progress.update(task_id, total=len(futures))\n            with
  progress:\n                while not all(f.done() for f in futures):\n                    time.sleep(0.1)\n
  \                   progress.update(task_id, total=len([f for f in futures if f.done()]))\n
  \           [f.result() for f in futures]\n```\n\n\n!! method <h2 id='__init__'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__init__ <em class='small'>method</em></h2>\n\n???+
  source \"__init__ <em class='small'>source</em>\"\n\n```python\n\n        def __init__(\n
  \               self,\n                msg: str = \"\",\n            ) -> None:\n
  \               super().__init__(\n                    \"Padding must be an iterable
  of length 1, 2, 3, or 4.\\n\" + msg,\n                )\n```\n"
date: 0001-01-01
description: Example configuration.  Covers supports multiple covers to be configured.  Here
  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ! ???+ source  ! ?
html: "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n<title>Covers.Py</title>\n<meta
  charset=\"utf-8\"/>\n<meta content=\"width=device-width, initial-scale=1\" name=\"viewport\"/>\n<meta
  content=\"Example configuration.  Covers supports multiple covers to be configured.
  \ Here ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ! ???+ source
  \ ! ?\" name=\"description\"/>\n<link href=\"/static/favicon.ico\" rel=\"icon\"
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
  name=\"og:type\" property=\"og:type\"/><meta content=\"Example configuration.  Covers
  supports multiple covers to be configured.  Here ! ???+ source  ! ???+ source  !
  ???+ source  ! ???+ source  ! ! ???+ source  ! ?\" name=\"description\" property=\"description\"/><meta
  content=\"Example configuration.  Covers supports multiple covers to be configured.
  \ Here ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ! ???+ source
  \ ! ?\" name=\"og:description\" property=\"og:description\"/><meta content=\"Example
  configuration.  Covers supports multiple covers to be configured.  Here ! ???+ source
  \ ! ???+ source  ! ???+ source  ! ???+ source  ! ! ???+ source  ! ?\" name=\"twitter:description\"
  property=\"twitter:description\"/><meta content=\"Covers.Py | Markata\" name=\"og:title\"
  property=\"og:title\"/><meta content=\"Covers.Py | Markata\" name=\"twitter:title\"
  property=\"twitter:title\"/><meta content=\"https://markata.dev//markata/plugins/covers-og.png\"
  name=\"og:image\" property=\"og:image\"/><meta content=\"https://markata.dev//markata/plugins/covers-og.png\"
  name=\"twitter:image\" property=\"twitter:image\"/><meta content=\"1600\" name=\"og:image:width\"
  property=\"og:image:width\"/><meta content=\"900\" name=\"og:image:width\" property=\"og:image:width\"/><meta
  content=\"summary_large_image\" name=\"twitter:card\" property=\"twitter:card\"/><meta
  content=\"Markata\" name=\"og:site_name\" property=\"og:site_name\"/><meta content=\"@_waylonwalker\"
  name=\"twitter:creator\" property=\"twitter:creator\"/><meta content=\"Covers.Py\"
  name=\"title\" property=\"title\"/><meta content=\"markata 0.8.0.dev6\" name=\"generator\"
  property=\"generator\"/><link href=\"https://markata.dev//markata/plugins/covers/\"
  rel=\"canonical\"/><meta content=\"https://markata.dev//markata/plugins/covers/\"
  name=\"og:url\" property=\"og:url\"/></head>\n<body><nav>\n<a href=\"/\">markata</a>\n<a
  href=\"https://github.com/WaylonWalker/markata\">GitHub</a>\n</nav>\n<div>\n<label
  class=\"theme-switch\" for=\"checkbox-theme\" id=\"theme-switch\" title=\"light/dark
  mode toggle\">\n<input id=\"checkbox-theme\" type=\"checkbox\"/>\n<div class=\"slider
  round\"></div>\n</label>\n</div>\n<section class=\"title\">\n<h1 id=\"title\">\n
  \           Covers.Py \n            \n        </h1>\n</section>\n<main><h3>Configuration</h3>\n<p>Example
  configuration.  Covers supports multiple covers to be configured.  Here\nis an example
  from my blog where we have a template sized for <a href=\"http://dev.to\">dev.to</a>
  and one\nsized for open graph.  Each image takes it's own configuration.</p>\n<pre
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
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.covers]]</span>\n<span
  class=\"n\">name</span><span class=\"o\">=</span><span class=\"s1\">'-dev'</span>\n<span
  class=\"n\">template</span><span class=\"w\"> </span><span class=\"o\">=</span><span
  class=\"w\"> </span><span class=\"s2\">\"static/cover-template.png\"</span>\n<span
  class=\"n\">font</span><span class=\"w\"> </span><span class=\"o\">=</span><span
  class=\"w\"> </span><span class=\"s2\">\"./static/JosefinSans-Regular.ttf\"</span>\n<span
  class=\"n\">text_font</span><span class=\"w\"> </span><span class=\"o\">=</span><span
  class=\"w\"> </span><span class=\"s2\">\"./static/JosefinSans-Regular.ttf\"</span>\n<span
  class=\"n\">font_color</span><span class=\"w\"> </span><span class=\"o\">=</span><span
  class=\"w\"> </span><span class=\"s2\">\"rgb(185,155,165)\"</span>\n<span class=\"n\">text_font_color</span><span
  class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s2\">\"rgb(255,255,255)\"</span>\n<span
  class=\"n\">text_key</span><span class=\"w\"> </span><span class=\"o\">=</span><span
  class=\"w\"> </span><span class=\"s1\">'description'</span>\n<span class=\"n\">padding</span><span
  class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"p\">[</span><span
  class=\"mi\">0</span><span class=\"p\">,</span><span class=\"w\"> </span><span class=\"mi\">40</span><span
  class=\"p\">,</span><span class=\"w\"> </span><span class=\"mi\">100</span><span
  class=\"p\">,</span><span class=\"w\"> </span><span class=\"mi\">300</span><span
  class=\"p\">]</span>\n<span class=\"n\">text_padding</span><span class=\"w\"> </span><span
  class=\"o\">=</span><span class=\"w\"> </span><span class=\"p\">[</span><span class=\"mi\">0</span><span
  class=\"p\">,</span><span class=\"mi\">0</span><span class=\"p\">]</span>\n\n<span
  class=\"k\">[[markata.covers]]</span>\n<span class=\"n\">name</span><span class=\"o\">=</span><span
  class=\"s1\">''</span>\n<span class=\"n\">template</span><span class=\"w\"> </span><span
  class=\"o\">=</span><span class=\"w\"> </span><span class=\"s2\">\"static/og-template.png\"</span>\n<span
  class=\"n\">font</span><span class=\"w\"> </span><span class=\"o\">=</span><span
  class=\"w\"> </span><span class=\"s2\">\"./static/JosefinSans-Regular.ttf\"</span>\n<span
  class=\"n\">font_color</span><span class=\"w\"> </span><span class=\"o\">=</span><span
  class=\"w\"> </span><span class=\"s2\">\"rgb(255,255,255)\"</span>\n<span class=\"n\">text_font</span><span
  class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s2\">\"./static/JosefinSans-Regular.ttf\"</span>\n<span
  class=\"n\">text_font_color</span><span class=\"w\"> </span><span class=\"o\">=</span><span
  class=\"w\"> </span><span class=\"s2\">\"rgb(200,200,200)\"</span>\n<span class=\"n\">text_key</span><span
  class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s1\">'description'</span>\n<span
  class=\"n\">padding</span><span class=\"w\"> </span><span class=\"o\">=</span><span
  class=\"w\"> </span><span class=\"p\">[</span><span class=\"mi\">10</span><span
  class=\"p\">,</span><span class=\"w\"> </span><span class=\"mi\">10</span><span
  class=\"p\">,</span><span class=\"w\"> </span><span class=\"mi\">100</span><span
  class=\"p\">,</span><span class=\"w\"> </span><span class=\"mi\">300</span><span
  class=\"p\">]</span>\n<span class=\"n\">text_padding</span><span class=\"w\"> </span><span
  class=\"o\">=</span><span class=\"w\"> </span><span class=\"p\">[</span><span class=\"mi\">0</span><span
  class=\"p\">,</span><span class=\"mi\">0</span><span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"_load_font\" style=\"margin:0;padding:.5rem
  1rem;\">_load_font <em class=\"small\">function</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"_load_font
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
  <span class=\"nf\">_load_font</span><span class=\"p\">(</span><span class=\"n\">path</span><span
  class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">,</span> <span
  class=\"n\">size</span><span class=\"p\">:</span> <span class=\"nb\">int</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">ImageFont</span><span
  class=\"o\">.</span><span class=\"n\">FreeTypeFont</span><span class=\"p\">:</span>\n
  \           <span class=\"k\">return</span> <span class=\"n\">ImageFont</span><span
  class=\"o\">.</span><span class=\"n\">truetype</span><span class=\"p\">(</span><span
  class=\"n\">path</span><span class=\"p\">,</span> <span class=\"n\">size</span><span
  class=\"o\">=</span><span class=\"n\">size</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"get_font\" style=\"margin:0;padding:.5rem
  1rem;\">get_font <em class=\"small\">function</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"get_font
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
  <span class=\"nf\">get_font</span><span class=\"p\">(</span>\n            <span
  class=\"n\">path</span><span class=\"p\">:</span> <span class=\"n\">Path</span><span
  class=\"p\">,</span>\n            <span class=\"n\">draw</span><span class=\"p\">:</span>
  <span class=\"n\">ImageDraw</span><span class=\"o\">.</span><span class=\"n\">Draw</span><span
  class=\"p\">,</span>\n            <span class=\"n\">title</span><span class=\"p\">:</span>
  <span class=\"nb\">str</span><span class=\"p\">,</span>\n            <span class=\"n\">size</span><span
  class=\"p\">:</span> <span class=\"nb\">int</span> <span class=\"o\">=</span> <span
  class=\"mi\">250</span><span class=\"p\">,</span>\n            <span class=\"n\">max_size</span><span
  class=\"p\">:</span> <span class=\"nb\">tuple</span> <span class=\"o\">=</span>
  <span class=\"p\">(</span><span class=\"mi\">800</span><span class=\"p\">,</span>
  <span class=\"mi\">220</span><span class=\"p\">),</span>\n        <span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"n\">ImageFont</span><span class=\"o\">.</span><span
  class=\"n\">FreeTypeFont</span><span class=\"p\">:</span>\n            <span class=\"n\">title</span>
  <span class=\"o\">=</span> <span class=\"n\">title</span> <span class=\"ow\">or</span>
  <span class=\"s2\">\"\"</span>\n            <span class=\"n\">font</span> <span
  class=\"o\">=</span> <span class=\"n\">_load_font</span><span class=\"p\">(</span><span
  class=\"n\">path</span><span class=\"p\">,</span> <span class=\"n\">size</span><span
  class=\"p\">)</span>\n            <span class=\"n\">current_size</span> <span class=\"o\">=</span>
  <span class=\"n\">draw</span><span class=\"o\">.</span><span class=\"n\">textsize</span><span
  class=\"p\">(</span><span class=\"n\">title</span><span class=\"p\">,</span> <span
  class=\"n\">font</span><span class=\"o\">=</span><span class=\"n\">font</span><span
  class=\"p\">)</span>\n\n            <span class=\"k\">if</span> <span class=\"n\">current_size</span><span
  class=\"p\">[</span><span class=\"mi\">0</span><span class=\"p\">]</span> <span
  class=\"o\">&gt;</span> <span class=\"n\">max_size</span><span class=\"p\">[</span><span
  class=\"mi\">0</span><span class=\"p\">]</span> <span class=\"ow\">or</span> <span
  class=\"n\">current_size</span><span class=\"p\">[</span><span class=\"mi\">1</span><span
  class=\"p\">]</span> <span class=\"o\">&gt;</span> <span class=\"n\">max_size</span><span
  class=\"p\">[</span><span class=\"mi\">1</span><span class=\"p\">]:</span>\n                <span
  class=\"k\">return</span> <span class=\"n\">get_font</span><span class=\"p\">(</span><span
  class=\"n\">path</span><span class=\"p\">,</span> <span class=\"n\">draw</span><span
  class=\"p\">,</span> <span class=\"n\">title</span><span class=\"p\">,</span> <span
  class=\"n\">size</span> <span class=\"o\">-</span> <span class=\"mi\">10</span><span
  class=\"p\">,</span> <span class=\"n\">max_size</span><span class=\"o\">=</span><span
  class=\"n\">max_size</span><span class=\"p\">)</span>\n            <span class=\"k\">return</span>
  <span class=\"n\">font</span>\n</pre></div>\n\n</pre>\n<p>!! class </p><h2 class=\"admonition-title\"
  id=\"PaddingError\" style=\"margin:0;padding:.5rem 1rem;\">PaddingError <em class=\"small\">class</em></h2>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"PaddingError
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
  <span class=\"nc\">PaddingError</span><span class=\"p\">(</span><span class=\"ne\">BaseException</span><span
  class=\"p\">):</span>\n            <span class=\"k\">def</span> <span class=\"fm\">__init__</span><span
  class=\"p\">(</span>\n                <span class=\"bp\">self</span><span class=\"p\">,</span>\n
  \               <span class=\"n\">msg</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
  <span class=\"o\">=</span> <span class=\"s2\">\"\"</span><span class=\"p\">,</span>\n
  \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n                <span class=\"nb\">super</span><span class=\"p\">()</span><span
  class=\"o\">.</span><span class=\"fm\">__init__</span><span class=\"p\">(</span>\n
  \                   <span class=\"s2\">\"Padding must be an iterable of length 1,
  2, 3, or 4.</span><span class=\"se\">\\n</span><span class=\"s2\">\"</span> <span
  class=\"o\">+</span> <span class=\"n\">msg</span><span class=\"p\">,</span>\n                <span
  class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!! function </p><h2 class=\"admonition-title\"
  id=\"draw_text\" style=\"margin:0;padding:.5rem 1rem;\">draw_text <em class=\"small\">function</em></h2>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"draw_text
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
  <span class=\"nf\">draw_text</span><span class=\"p\">(</span>\n            <span
  class=\"n\">image</span><span class=\"p\">:</span> <span class=\"n\">Image</span><span
  class=\"p\">,</span>\n            <span class=\"n\">font_path</span><span class=\"p\">:</span>
  <span class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"n\">Path</span><span
  class=\"p\">],</span>\n            <span class=\"n\">text</span><span class=\"p\">:</span>
  <span class=\"nb\">str</span><span class=\"p\">,</span>\n            <span class=\"n\">color</span><span
  class=\"p\">:</span> <span class=\"n\">Union</span><span class=\"p\">[</span><span
  class=\"nb\">str</span><span class=\"p\">,</span> <span class=\"kc\">None</span><span
  class=\"p\">],</span>\n            <span class=\"n\">padding</span><span class=\"p\">:</span>
  <span class=\"n\">Tuple</span><span class=\"p\">[</span><span class=\"nb\">int</span><span
  class=\"p\">,</span> <span class=\"o\">...</span><span class=\"p\">],</span>\n            <span
  class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span
  class=\"p\">,</span>\n        <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"n\">text</span>
  <span class=\"o\">=</span> <span class=\"n\">text</span> <span class=\"ow\">or</span>
  <span class=\"s2\">\"\"</span>\n            <span class=\"n\">draw</span> <span
  class=\"o\">=</span> <span class=\"n\">ImageDraw</span><span class=\"o\">.</span><span
  class=\"n\">Draw</span><span class=\"p\">(</span><span class=\"n\">image</span><span
  class=\"p\">)</span>\n            <span class=\"n\">padding</span> <span class=\"o\">=</span>
  <span class=\"n\">resolve_padding</span><span class=\"p\">(</span><span class=\"n\">padding</span><span
  class=\"p\">,</span> <span class=\"n\">markata</span><span class=\"p\">)</span>\n
  \           <span class=\"n\">width</span> <span class=\"o\">=</span> <span class=\"n\">image</span><span
  class=\"o\">.</span><span class=\"n\">size</span><span class=\"p\">[</span><span
  class=\"mi\">0</span><span class=\"p\">]</span>\n            <span class=\"n\">height</span>
  <span class=\"o\">=</span> <span class=\"n\">image</span><span class=\"o\">.</span><span
  class=\"n\">size</span><span class=\"p\">[</span><span class=\"mi\">1</span><span
  class=\"p\">]</span>\n            <span class=\"n\">bounding_box</span> <span class=\"o\">=</span>
  <span class=\"p\">[</span><span class=\"n\">padding</span><span class=\"p\">[</span><span
  class=\"mi\">0</span><span class=\"p\">],</span> <span class=\"n\">padding</span><span
  class=\"p\">[</span><span class=\"mi\">1</span><span class=\"p\">],</span> <span
  class=\"n\">width</span> <span class=\"o\">-</span> <span class=\"n\">padding</span><span
  class=\"p\">[</span><span class=\"mi\">0</span><span class=\"p\">],</span> <span
  class=\"n\">height</span> <span class=\"o\">-</span> <span class=\"n\">padding</span><span
  class=\"p\">[</span><span class=\"mi\">1</span><span class=\"p\">]]</span>\n            <span
  class=\"n\">bounding_box</span> <span class=\"o\">=</span> <span class=\"p\">[</span><span
  class=\"n\">padding</span><span class=\"p\">[</span><span class=\"mi\">0</span><span
  class=\"p\">],</span> <span class=\"n\">padding</span><span class=\"p\">[</span><span
  class=\"mi\">1</span><span class=\"p\">],</span> <span class=\"n\">width</span>
  <span class=\"o\">-</span> <span class=\"n\">padding</span><span class=\"p\">[</span><span
  class=\"mi\">2</span><span class=\"p\">],</span> <span class=\"n\">height</span>
  <span class=\"o\">-</span> <span class=\"n\">padding</span><span class=\"p\">[</span><span
  class=\"mi\">3</span><span class=\"p\">]]</span>\n            <span class=\"n\">max_size</span>
  <span class=\"o\">=</span> <span class=\"p\">(</span><span class=\"n\">bounding_box</span><span
  class=\"p\">[</span><span class=\"mi\">2</span><span class=\"p\">]</span> <span
  class=\"o\">-</span> <span class=\"n\">bounding_box</span><span class=\"p\">[</span><span
  class=\"mi\">0</span><span class=\"p\">],</span> <span class=\"n\">bounding_box</span><span
  class=\"p\">[</span><span class=\"mi\">3</span><span class=\"p\">]</span> <span
  class=\"o\">-</span> <span class=\"n\">bounding_box</span><span class=\"p\">[</span><span
  class=\"mi\">1</span><span class=\"p\">])</span>\n            <span class=\"n\">x1</span><span
  class=\"p\">,</span> <span class=\"n\">y1</span><span class=\"p\">,</span> <span
  class=\"n\">x2</span><span class=\"p\">,</span> <span class=\"n\">y2</span> <span
  class=\"o\">=</span> <span class=\"n\">bounding_box</span>\n            <span class=\"n\">font</span>
  <span class=\"o\">=</span> <span class=\"n\">get_font</span><span class=\"p\">(</span><span
  class=\"n\">font_path</span><span class=\"p\">,</span> <span class=\"n\">draw</span><span
  class=\"p\">,</span> <span class=\"n\">text</span><span class=\"p\">,</span> <span
  class=\"n\">max_size</span><span class=\"o\">=</span><span class=\"n\">max_size</span><span
  class=\"p\">)</span> <span class=\"k\">if</span> <span class=\"n\">font_path</span>
  <span class=\"k\">else</span> <span class=\"kc\">None</span>\n            <span
  class=\"n\">w</span><span class=\"p\">,</span> <span class=\"n\">h</span> <span
  class=\"o\">=</span> <span class=\"n\">draw</span><span class=\"o\">.</span><span
  class=\"n\">textsize</span><span class=\"p\">(</span><span class=\"n\">text</span><span
  class=\"p\">,</span> <span class=\"n\">font</span><span class=\"o\">=</span><span
  class=\"n\">font</span><span class=\"p\">)</span>\n            <span class=\"n\">x</span>
  <span class=\"o\">=</span> <span class=\"p\">(</span><span class=\"n\">x2</span>
  <span class=\"o\">-</span> <span class=\"n\">x1</span> <span class=\"o\">-</span>
  <span class=\"n\">w</span><span class=\"p\">)</span> <span class=\"o\">/</span>
  <span class=\"mi\">2</span> <span class=\"o\">+</span> <span class=\"n\">x1</span>\n
  \           <span class=\"n\">y</span> <span class=\"o\">=</span> <span class=\"p\">(</span><span
  class=\"n\">y2</span> <span class=\"o\">-</span> <span class=\"n\">y1</span> <span
  class=\"o\">-</span> <span class=\"n\">h</span><span class=\"p\">)</span> <span
  class=\"o\">/</span> <span class=\"mi\">2</span> <span class=\"o\">+</span> <span
  class=\"n\">y1</span>\n            <span class=\"n\">draw</span><span class=\"o\">.</span><span
  class=\"n\">text</span><span class=\"p\">((</span><span class=\"n\">x</span><span
  class=\"p\">,</span> <span class=\"n\">y</span><span class=\"p\">),</span> <span
  class=\"n\">text</span><span class=\"p\">,</span> <span class=\"n\">fill</span><span
  class=\"o\">=</span><span class=\"n\">color</span><span class=\"p\">,</span> <span
  class=\"n\">font</span><span class=\"o\">=</span><span class=\"n\">font</span><span
  class=\"p\">,</span> <span class=\"n\">align</span><span class=\"o\">=</span><span
  class=\"s2\">\"center\"</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"resolve_padding\" style=\"margin:0;padding:.5rem
  1rem;\">resolve_padding <em class=\"small\">function</em></h2>\nConvert padding
  to a len 4 tuple\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
  class=\"admonition-title\">\"resolve_padding <em class=\"small\">source</em>\"</p>\n</div>\n<pre
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
  <span class=\"nf\">resolve_padding</span><span class=\"p\">(</span><span class=\"n\">padding</span><span
  class=\"p\">:</span> <span class=\"n\">Tuple</span><span class=\"p\">[</span><span
  class=\"nb\">int</span><span class=\"p\">,</span> <span class=\"o\">...</span><span
  class=\"p\">],</span> <span class=\"n\">markata</span><span class=\"p\">:</span>
  <span class=\"s2\">\"Markata\"</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"n\">Tuple</span><span class=\"p\">[</span><span class=\"nb\">int</span><span
  class=\"p\">,</span> <span class=\"o\">...</span><span class=\"p\">]:</span>\n<span
  class=\"w\">            </span><span class=\"sd\">\"\"\"Convert padding to a len
  4 tuple\"\"\"</span>\n            <span class=\"k\">if</span> <span class=\"nb\">len</span><span
  class=\"p\">(</span><span class=\"n\">padding</span><span class=\"p\">)</span> <span
  class=\"o\">==</span> <span class=\"mi\">4</span><span class=\"p\">:</span>\n                <span
  class=\"k\">return</span> <span class=\"n\">padding</span>\n            <span class=\"k\">if</span>
  <span class=\"nb\">len</span><span class=\"p\">(</span><span class=\"n\">padding</span><span
  class=\"p\">)</span> <span class=\"o\">==</span> <span class=\"mi\">3</span><span
  class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"p\">(</span><span
  class=\"o\">*</span><span class=\"n\">padding</span><span class=\"p\">,</span> <span
  class=\"n\">padding</span><span class=\"p\">[</span><span class=\"mi\">1</span><span
  class=\"p\">])</span>\n            <span class=\"k\">if</span> <span class=\"nb\">len</span><span
  class=\"p\">(</span><span class=\"n\">padding</span><span class=\"p\">)</span> <span
  class=\"o\">==</span> <span class=\"mi\">2</span><span class=\"p\">:</span>\n                <span
  class=\"k\">return</span> <span class=\"n\">padding</span> <span class=\"o\">*</span>
  <span class=\"mi\">2</span>\n            <span class=\"k\">if</span> <span class=\"nb\">len</span><span
  class=\"p\">(</span><span class=\"n\">padding</span><span class=\"p\">)</span> <span
  class=\"o\">==</span> <span class=\"mi\">1</span><span class=\"p\">:</span>\n                <span
  class=\"k\">return</span> <span class=\"n\">padding</span> <span class=\"o\">*</span>
  <span class=\"mi\">4</span>\n            <span class=\"k\">raise</span> <span class=\"n\">PaddingError</span><span
  class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">\"recieved padding:
  </span><span class=\"si\">{</span><span class=\"n\">padding</span><span class=\"si\">}</span><span
  class=\"s2\">\"</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"make_cover\" style=\"margin:0;padding:.5rem
  1rem;\">make_cover <em class=\"small\">function</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"make_cover
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
  <span class=\"nf\">make_cover</span><span class=\"p\">(</span>\n            <span
  class=\"n\">title</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
  class=\"p\">,</span>\n            <span class=\"n\">color</span><span class=\"p\">:</span>
  <span class=\"nb\">str</span><span class=\"p\">,</span>\n            <span class=\"n\">output_path</span><span
  class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">,</span>\n            <span
  class=\"n\">template_path</span><span class=\"p\">:</span> <span class=\"n\">Path</span><span
  class=\"p\">,</span>\n            <span class=\"n\">font_path</span><span class=\"p\">:</span>
  <span class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"n\">Path</span><span
  class=\"p\">],</span>\n            <span class=\"n\">padding</span><span class=\"p\">:</span>
  <span class=\"n\">Tuple</span><span class=\"p\">[</span><span class=\"nb\">int</span><span
  class=\"p\">,</span> <span class=\"o\">...</span><span class=\"p\">],</span>\n            <span
  class=\"n\">text_font</span><span class=\"p\">:</span> <span class=\"n\">Path</span><span
  class=\"p\">,</span>\n            <span class=\"n\">text</span><span class=\"p\">:</span>
  <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"kc\">None</span><span
  class=\"p\">,</span>\n            <span class=\"n\">text_font_color</span><span
  class=\"p\">:</span> <span class=\"nb\">str</span> <span class=\"o\">=</span> <span
  class=\"kc\">None</span><span class=\"p\">,</span>\n            <span class=\"n\">text_padding</span><span
  class=\"p\">:</span> <span class=\"n\">Tuple</span><span class=\"p\">[</span><span
  class=\"nb\">int</span><span class=\"p\">,</span> <span class=\"o\">...</span><span
  class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span><span
  class=\"p\">,</span>\n            <span class=\"n\">resizes</span><span class=\"p\">:</span>
  <span class=\"n\">List</span><span class=\"p\">[</span><span class=\"nb\">int</span><span
  class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span><span
  class=\"p\">,</span>\n            <span class=\"n\">markata</span><span class=\"p\">:</span>
  <span class=\"s2\">\"Markata\"</span> <span class=\"o\">=</span> <span class=\"kc\">None</span><span
  class=\"p\">,</span>\n        <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"k\">if</span>
  <span class=\"n\">output_path</span><span class=\"o\">.</span><span class=\"n\">exists</span><span
  class=\"p\">():</span>\n                <span class=\"k\">return</span>\n            <span
  class=\"n\">image</span> <span class=\"o\">=</span> <span class=\"n\">Image</span><span
  class=\"o\">.</span><span class=\"n\">open</span><span class=\"p\">(</span><span
  class=\"n\">template_path</span><span class=\"p\">)</span> <span class=\"k\">if</span>
  <span class=\"n\">template_path</span> <span class=\"k\">else</span> <span class=\"n\">Image</span><span
  class=\"o\">.</span><span class=\"n\">new</span><span class=\"p\">(</span><span
  class=\"s2\">\"RGB\"</span><span class=\"p\">,</span> <span class=\"p\">(</span><span
  class=\"mi\">800</span><span class=\"p\">,</span> <span class=\"mi\">450</span><span
  class=\"p\">))</span>\n\n            <span class=\"n\">draw_text</span><span class=\"p\">(</span>\n
  \               <span class=\"n\">image</span><span class=\"o\">=</span><span class=\"n\">image</span><span
  class=\"p\">,</span>\n                <span class=\"n\">font_path</span><span class=\"o\">=</span><span
  class=\"n\">font_path</span><span class=\"p\">,</span>\n                <span class=\"n\">title</span><span
  class=\"o\">=</span><span class=\"n\">title</span><span class=\"p\">,</span>\n                <span
  class=\"n\">color</span><span class=\"o\">=</span><span class=\"n\">color</span><span
  class=\"p\">,</span>\n                <span class=\"n\">padding</span><span class=\"o\">=</span><span
  class=\"n\">padding</span><span class=\"p\">,</span>\n                <span class=\"n\">markata</span><span
  class=\"o\">=</span><span class=\"n\">markata</span><span class=\"p\">,</span>\n
  \           <span class=\"p\">)</span>\n            <span class=\"k\">if</span>
  <span class=\"n\">text</span> <span class=\"ow\">is</span> <span class=\"ow\">not</span>
  <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
  class=\"k\">if</span> <span class=\"n\">text_padding</span> <span class=\"ow\">is</span>
  <span class=\"kc\">None</span><span class=\"p\">:</span>\n                    <span
  class=\"n\">text_padding</span> <span class=\"o\">=</span> <span class=\"p\">(</span>\n
  \                       <span class=\"n\">image</span><span class=\"o\">.</span><span
  class=\"n\">size</span><span class=\"p\">[</span><span class=\"mi\">1</span><span
  class=\"p\">]</span> <span class=\"o\">-</span> <span class=\"n\">image</span><span
  class=\"o\">.</span><span class=\"n\">size</span><span class=\"p\">[</span><span
  class=\"mi\">1</span><span class=\"p\">]</span> <span class=\"o\">/</span> <span
  class=\"mi\">5</span><span class=\"p\">,</span>\n                        <span class=\"n\">image</span><span
  class=\"o\">.</span><span class=\"n\">size</span><span class=\"p\">[</span><span
  class=\"mi\">0</span><span class=\"p\">]</span> <span class=\"o\">/</span> <span
  class=\"mi\">5</span><span class=\"p\">,</span>\n                        <span class=\"n\">image</span><span
  class=\"o\">.</span><span class=\"n\">size</span><span class=\"p\">[</span><span
  class=\"mi\">1</span><span class=\"p\">]</span> <span class=\"o\">-</span> <span
  class=\"n\">image</span><span class=\"o\">.</span><span class=\"n\">size</span><span
  class=\"p\">[</span><span class=\"mi\">1</span><span class=\"p\">]</span> <span
  class=\"o\">/</span> <span class=\"mi\">10</span><span class=\"p\">,</span>\n                    <span
  class=\"p\">)</span>\n                <span class=\"n\">draw_text</span><span class=\"p\">(</span><span
  class=\"n\">image</span><span class=\"p\">,</span> <span class=\"n\">text_font</span><span
  class=\"p\">,</span> <span class=\"n\">text</span><span class=\"p\">,</span> <span
  class=\"n\">text_font_color</span><span class=\"p\">,</span> <span class=\"n\">text_padding</span><span
  class=\"p\">)</span>\n\n            <span class=\"n\">image</span><span class=\"o\">.</span><span
  class=\"n\">save</span><span class=\"p\">(</span><span class=\"n\">output_path</span><span
  class=\"p\">)</span>\n            <span class=\"n\">ratio</span> <span class=\"o\">=</span>
  <span class=\"n\">image</span><span class=\"o\">.</span><span class=\"n\">size</span><span
  class=\"p\">[</span><span class=\"mi\">1</span><span class=\"p\">]</span> <span
  class=\"o\">/</span> <span class=\"n\">image</span><span class=\"o\">.</span><span
  class=\"n\">size</span><span class=\"p\">[</span><span class=\"mi\">0</span><span
  class=\"p\">]</span>\n\n            <span class=\"n\">covers</span> <span class=\"o\">=</span>
  <span class=\"p\">[]</span>\n            <span class=\"k\">if</span> <span class=\"n\">resizes</span><span
  class=\"p\">:</span>\n                <span class=\"k\">for</span> <span class=\"n\">width</span>
  <span class=\"ow\">in</span> <span class=\"n\">resizes</span><span class=\"p\">:</span>\n
  \                   <span class=\"n\">re_img</span> <span class=\"o\">=</span> <span
  class=\"n\">image</span><span class=\"o\">.</span><span class=\"n\">resize</span><span
  class=\"p\">((</span><span class=\"n\">width</span><span class=\"p\">,</span> <span
  class=\"nb\">int</span><span class=\"p\">(</span><span class=\"n\">width</span>
  <span class=\"o\">*</span> <span class=\"n\">ratio</span><span class=\"p\">)),</span>
  <span class=\"n\">Image</span><span class=\"o\">.</span><span class=\"n\">ANTIALIAS</span><span
  class=\"p\">)</span>\n                    <span class=\"n\">filename</span> <span
  class=\"o\">=</span> <span class=\"p\">(</span>\n                        <span class=\"sa\">f</span><span
  class=\"s2\">\"</span><span class=\"si\">{</span><span class=\"n\">output_path</span><span
  class=\"o\">.</span><span class=\"n\">stem</span><span class=\"si\">}</span><span
  class=\"s2\">_</span><span class=\"si\">{</span><span class=\"n\">width</span><span
  class=\"si\">}</span><span class=\"s2\">x</span><span class=\"si\">{</span><span
  class=\"nb\">int</span><span class=\"p\">(</span><span class=\"n\">width</span><span
  class=\"o\">*</span><span class=\"n\">ratio</span><span class=\"p\">)</span><span
  class=\"si\">}{</span><span class=\"n\">output_path</span><span class=\"o\">.</span><span
  class=\"n\">suffix</span><span class=\"si\">}</span><span class=\"s2\">\"</span>\n
  \                   <span class=\"p\">)</span>\n                    <span class=\"n\">covers</span><span
  class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
  class=\"n\">filename</span><span class=\"p\">)</span>\n\n                    <span
  class=\"n\">filepath</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
  class=\"p\">(</span><span class=\"n\">output_path</span><span class=\"o\">.</span><span
  class=\"n\">parent</span> <span class=\"o\">/</span> <span class=\"n\">filename</span><span
  class=\"p\">)</span>\n                    <span class=\"n\">re_img</span><span class=\"o\">.</span><span
  class=\"n\">save</span><span class=\"p\">(</span><span class=\"n\">filepath</span><span
  class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!! function </p><h2 class=\"admonition-title\"
  id=\"save\" style=\"margin:0;padding:.5rem 1rem;\">save <em class=\"small\">function</em></h2>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"save
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
  <span class=\"nf\">save</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \           <span class=\"n\">futures</span> <span class=\"o\">=</span> <span class=\"p\">[]</span>\n\n
  \           <span class=\"k\">if</span> <span class=\"s2\">\"covers\"</span> <span
  class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
  class=\"n\">keys</span><span class=\"p\">():</span>\n                <span class=\"k\">return</span>\n\n
  \           <span class=\"k\">for</span> <span class=\"n\">article</span> <span
  class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">iter_articles</span><span class=\"p\">(</span><span class=\"s2\">\"making
  covers\"</span><span class=\"p\">):</span>\n                <span class=\"k\">for</span>
  <span class=\"n\">cover</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">[</span><span
  class=\"s2\">\"covers\"</span><span class=\"p\">]:</span>\n                    <span
  class=\"k\">try</span><span class=\"p\">:</span>\n                        <span
  class=\"n\">padding</span> <span class=\"o\">=</span> <span class=\"n\">cover</span><span
  class=\"p\">[</span><span class=\"s2\">\"padding\"</span><span class=\"p\">]</span>\n
  \                   <span class=\"k\">except</span> <span class=\"ne\">KeyError</span><span
  class=\"p\">:</span>\n                        <span class=\"n\">padding</span> <span
  class=\"o\">=</span> <span class=\"p\">(</span>\n                            <span
  class=\"mi\">200</span><span class=\"p\">,</span>\n                            <span
  class=\"mi\">100</span><span class=\"p\">,</span>\n                        <span
  class=\"p\">)</span>\n                    <span class=\"k\">try</span><span class=\"p\">:</span>\n
  \                       <span class=\"n\">text_padding</span> <span class=\"o\">=</span>
  <span class=\"n\">cover</span><span class=\"p\">[</span><span class=\"s2\">\"text_padding\"</span><span
  class=\"p\">]</span>\n                    <span class=\"k\">except</span> <span
  class=\"ne\">KeyError</span><span class=\"p\">:</span>\n                        <span
  class=\"n\">text_padding</span> <span class=\"o\">=</span> <span class=\"p\">(</span>\n
  \                           <span class=\"mi\">200</span><span class=\"p\">,</span>\n
  \                           <span class=\"mi\">100</span><span class=\"p\">,</span>\n
  \                       <span class=\"p\">)</span>\n                    <span class=\"k\">if</span>
  <span class=\"s2\">\"text_key\"</span> <span class=\"ow\">in</span> <span class=\"n\">cover</span><span
  class=\"p\">:</span>\n                        <span class=\"k\">try</span><span
  class=\"p\">:</span>\n                            <span class=\"n\">text</span>
  <span class=\"o\">=</span> <span class=\"n\">article</span><span class=\"o\">.</span><span
  class=\"n\">metadata</span><span class=\"p\">[</span><span class=\"n\">cover</span><span
  class=\"p\">[</span><span class=\"s2\">\"text_key\"</span><span class=\"p\">]]</span>\n
  \                       <span class=\"k\">except</span> <span class=\"ne\">AttributeError</span><span
  class=\"p\">:</span>\n                            <span class=\"n\">text</span>
  <span class=\"o\">=</span> <span class=\"n\">article</span><span class=\"p\">[</span><span
  class=\"n\">cover</span><span class=\"p\">[</span><span class=\"s2\">\"text_key\"</span><span
  class=\"p\">]]</span>\n                        <span class=\"k\">try</span><span
  class=\"p\">:</span>\n                            <span class=\"n\">text</span>
  <span class=\"o\">=</span> <span class=\"n\">text</span><span class=\"o\">.</span><span
  class=\"n\">replace</span><span class=\"p\">(</span><span class=\"s2\">\"</span><span
  class=\"se\">\\n</span><span class=\"s2\">\"</span><span class=\"p\">,</span> <span
  class=\"s2\">\"\"</span><span class=\"p\">)</span>\n                            <span
  class=\"kn\">from</span> <span class=\"nn\">more_itertools</span> <span class=\"kn\">import</span>
  <span class=\"n\">chunked</span>\n\n                            <span class=\"n\">text</span>
  <span class=\"o\">=</span> <span class=\"s2\">\"</span><span class=\"se\">\\n</span><span
  class=\"s2\">\"</span><span class=\"o\">.</span><span class=\"n\">join</span><span
  class=\"p\">([</span><span class=\"s2\">\"\"</span><span class=\"o\">.</span><span
  class=\"n\">join</span><span class=\"p\">(</span><span class=\"n\">c</span><span
  class=\"p\">)</span> <span class=\"k\">for</span> <span class=\"n\">c</span> <span
  class=\"ow\">in</span> <span class=\"n\">chunked</span><span class=\"p\">(</span><span
  class=\"n\">text</span><span class=\"p\">,</span> <span class=\"mi\">60</span><span
  class=\"p\">)])</span>\n                        <span class=\"k\">except</span>
  <span class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n                            <span
  class=\"c1\"># text is likely None</span>\n                            <span class=\"k\">pass</span>\n\n
  \                       <span class=\"n\">text_font</span> <span class=\"o\">=</span>
  <span class=\"n\">cover</span><span class=\"p\">[</span><span class=\"s2\">\"text_font\"</span><span
  class=\"p\">]</span>\n                        <span class=\"n\">text_font_color</span>
  <span class=\"o\">=</span> <span class=\"n\">cover</span><span class=\"p\">[</span><span
  class=\"s2\">\"text_font_color\"</span><span class=\"p\">]</span>\n                    <span
  class=\"k\">else</span><span class=\"p\">:</span>\n                        <span
  class=\"n\">text</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n
  \                       <span class=\"n\">text_font</span> <span class=\"o\">=</span>
  <span class=\"kc\">None</span>\n                        <span class=\"n\">text_font_color</span>
  <span class=\"o\">=</span> <span class=\"kc\">None</span>\n                    <span
  class=\"k\">try</span><span class=\"p\">:</span>\n                        <span
  class=\"n\">title</span> <span class=\"o\">=</span> <span class=\"n\">article</span><span
  class=\"o\">.</span><span class=\"n\">metadata</span><span class=\"p\">[</span><span
  class=\"s2\">\"title\"</span><span class=\"p\">]</span>\n                    <span
  class=\"k\">except</span> <span class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n
  \                       <span class=\"n\">title</span> <span class=\"o\">=</span>
  <span class=\"n\">article</span><span class=\"p\">[</span><span class=\"s2\">\"title\"</span><span
  class=\"p\">]</span>\n                    <span class=\"n\">futures</span><span
  class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span>\n
  \                       <span class=\"n\">make_cover</span><span class=\"p\">(</span>\n
  \                           <span class=\"n\">title</span><span class=\"o\">=</span><span
  class=\"n\">title</span><span class=\"p\">,</span>\n                            <span
  class=\"n\">color</span><span class=\"o\">=</span><span class=\"n\">cover</span><span
  class=\"p\">[</span><span class=\"s2\">\"font_color\"</span><span class=\"p\">],</span>\n
  \                           <span class=\"n\">output_path</span><span class=\"o\">=</span><span
  class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
  class=\"n\">output_dir</span><span class=\"p\">)</span>\n                            <span
  class=\"o\">/</span> <span class=\"p\">(</span><span class=\"n\">article</span><span
  class=\"p\">[</span><span class=\"s2\">\"slug\"</span><span class=\"p\">]</span>
  <span class=\"o\">+</span> <span class=\"n\">cover</span><span class=\"p\">[</span><span
  class=\"s2\">\"name\"</span><span class=\"p\">]</span> <span class=\"o\">+</span>
  <span class=\"s2\">\".png\"</span><span class=\"p\">),</span>\n                            <span
  class=\"n\">template_path</span><span class=\"o\">=</span><span class=\"n\">cover</span><span
  class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
  class=\"s2\">\"template\"</span><span class=\"p\">,</span> <span class=\"kc\">None</span><span
  class=\"p\">),</span>\n                            <span class=\"n\">font_path</span><span
  class=\"o\">=</span><span class=\"n\">cover</span><span class=\"o\">.</span><span
  class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">\"font\"</span><span
  class=\"p\">,</span> <span class=\"kc\">None</span><span class=\"p\">),</span>\n
  \                           <span class=\"n\">padding</span><span class=\"o\">=</span><span
  class=\"n\">padding</span><span class=\"p\">,</span>\n                            <span
  class=\"n\">text_font</span><span class=\"o\">=</span><span class=\"n\">text_font</span><span
  class=\"p\">,</span>\n                            <span class=\"n\">text</span><span
  class=\"o\">=</span><span class=\"n\">text</span><span class=\"p\">,</span>\n                            <span
  class=\"n\">text_font_color</span><span class=\"o\">=</span><span class=\"n\">text_font_color</span><span
  class=\"p\">,</span>\n                            <span class=\"n\">text_padding</span><span
  class=\"o\">=</span><span class=\"n\">text_padding</span><span class=\"p\">,</span>\n
  \                           <span class=\"n\">resizes</span><span class=\"o\">=</span><span
  class=\"n\">cover</span><span class=\"o\">.</span><span class=\"n\">get</span><span
  class=\"p\">(</span><span class=\"s2\">\"resizes\"</span><span class=\"p\">),</span>\n
  \                           <span class=\"n\">markata</span><span class=\"o\">=</span><span
  class=\"n\">markata</span><span class=\"p\">,</span>\n                        <span
  class=\"p\">),</span>\n                    <span class=\"p\">)</span>\n\n            <span
  class=\"n\">progress</span> <span class=\"o\">=</span> <span class=\"n\">Progress</span><span
  class=\"p\">(</span>\n                <span class=\"n\">BarColumn</span><span class=\"p\">(</span><span
  class=\"n\">bar_width</span><span class=\"o\">=</span><span class=\"kc\">None</span><span
  class=\"p\">),</span>\n                <span class=\"n\">transient</span><span class=\"o\">=</span><span
  class=\"kc\">True</span><span class=\"p\">,</span>\n                <span class=\"n\">console</span><span
  class=\"o\">=</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">console</span><span class=\"p\">,</span>\n            <span class=\"p\">)</span>\n
  \           <span class=\"n\">task_id</span> <span class=\"o\">=</span> <span class=\"n\">progress</span><span
  class=\"o\">.</span><span class=\"n\">add_task</span><span class=\"p\">(</span><span
  class=\"s2\">\"loading markdown\"</span><span class=\"p\">)</span>\n            <span
  class=\"n\">progress</span><span class=\"o\">.</span><span class=\"n\">update</span><span
  class=\"p\">(</span><span class=\"n\">task_id</span><span class=\"p\">,</span> <span
  class=\"n\">total</span><span class=\"o\">=</span><span class=\"nb\">len</span><span
  class=\"p\">(</span><span class=\"n\">futures</span><span class=\"p\">))</span>\n
  \           <span class=\"k\">with</span> <span class=\"n\">progress</span><span
  class=\"p\">:</span>\n                <span class=\"k\">while</span> <span class=\"ow\">not</span>
  <span class=\"nb\">all</span><span class=\"p\">(</span><span class=\"n\">f</span><span
  class=\"o\">.</span><span class=\"n\">done</span><span class=\"p\">()</span> <span
  class=\"k\">for</span> <span class=\"n\">f</span> <span class=\"ow\">in</span> <span
  class=\"n\">futures</span><span class=\"p\">):</span>\n                    <span
  class=\"n\">time</span><span class=\"o\">.</span><span class=\"n\">sleep</span><span
  class=\"p\">(</span><span class=\"mf\">0.1</span><span class=\"p\">)</span>\n                    <span
  class=\"n\">progress</span><span class=\"o\">.</span><span class=\"n\">update</span><span
  class=\"p\">(</span><span class=\"n\">task_id</span><span class=\"p\">,</span> <span
  class=\"n\">total</span><span class=\"o\">=</span><span class=\"nb\">len</span><span
  class=\"p\">([</span><span class=\"n\">f</span> <span class=\"k\">for</span> <span
  class=\"n\">f</span> <span class=\"ow\">in</span> <span class=\"n\">futures</span>
  <span class=\"k\">if</span> <span class=\"n\">f</span><span class=\"o\">.</span><span
  class=\"n\">done</span><span class=\"p\">()]))</span>\n            <span class=\"p\">[</span><span
  class=\"n\">f</span><span class=\"o\">.</span><span class=\"n\">result</span><span
  class=\"p\">()</span> <span class=\"k\">for</span> <span class=\"n\">f</span> <span
  class=\"ow\">in</span> <span class=\"n\">futures</span><span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"__init__\" style=\"margin:0;padding:.5rem
  1rem;\"><strong>init</strong> <em class=\"small\">method</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"<strong>init</strong>
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
  <span class=\"fm\">__init__</span><span class=\"p\">(</span>\n                <span
  class=\"bp\">self</span><span class=\"p\">,</span>\n                <span class=\"n\">msg</span><span
  class=\"p\">:</span> <span class=\"nb\">str</span> <span class=\"o\">=</span> <span
  class=\"s2\">\"\"</span><span class=\"p\">,</span>\n            <span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \               <span class=\"nb\">super</span><span class=\"p\">()</span><span
  class=\"o\">.</span><span class=\"fm\">__init__</span><span class=\"p\">(</span>\n
  \                   <span class=\"s2\">\"Padding must be an iterable of length 1,
  2, 3, or 4.</span><span class=\"se\">\\n</span><span class=\"s2\">\"</span> <span
  class=\"o\">+</span> <span class=\"n\">msg</span><span class=\"p\">,</span>\n                <span
  class=\"p\">)</span>\n</pre></div>\n\n</pre>\n</main>\n<footer>\xA9 2024</footer>\n</body></html>"
published: true
slug: markata/plugins/covers
title: Covers.Py


---

### Configuration

Example configuration.  Covers supports multiple covers to be configured.  Here
is an example from my blog where we have a template sized for dev.to and one
sized for open graph.  Each image takes it's own configuration.

``` toml
[[markata.covers]]
name='-dev'
template = "static/cover-template.png"
font = "./static/JosefinSans-Regular.ttf"
text_font = "./static/JosefinSans-Regular.ttf"
font_color = "rgb(185,155,165)"
text_font_color = "rgb(255,255,255)"
text_key = 'description'
padding = [0, 40, 100, 300]
text_padding = [0,0]

[[markata.covers]]
name=''
template = "static/og-template.png"
font = "./static/JosefinSans-Regular.ttf"
font_color = "rgb(255,255,255)"
text_font = "./static/JosefinSans-Regular.ttf"
text_font_color = "rgb(200,200,200)"
text_key = 'description'
padding = [10, 10, 100, 300]
text_padding = [0,0]
```


!! function <h2 id='_load_font' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_load_font <em class='small'>function</em></h2>

???+ source "_load_font <em class='small'>source</em>"

```python

        def _load_font(path: Path, size: int) -> ImageFont.FreeTypeFont:
            return ImageFont.truetype(path, size=size)
```


!! function <h2 id='get_font' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get_font <em class='small'>function</em></h2>

???+ source "get_font <em class='small'>source</em>"

```python

        def get_font(
            path: Path,
            draw: ImageDraw.Draw,
            title: str,
            size: int = 250,
            max_size: tuple = (800, 220),
        ) -> ImageFont.FreeTypeFont:
            title = title or ""
            font = _load_font(path, size)
            current_size = draw.textsize(title, font=font)

            if current_size[0] > max_size[0] or current_size[1] > max_size[1]:
                return get_font(path, draw, title, size - 10, max_size=max_size)
            return font
```


!! class <h2 id='PaddingError' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>PaddingError <em class='small'>class</em></h2>

???+ source "PaddingError <em class='small'>source</em>"

```python

        class PaddingError(BaseException):
            def __init__(
                self,
                msg: str = "",
            ) -> None:
                super().__init__(
                    "Padding must be an iterable of length 1, 2, 3, or 4.\n" + msg,
                )
```


!! function <h2 id='draw_text' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>draw_text <em class='small'>function</em></h2>

???+ source "draw_text <em class='small'>source</em>"

```python

        def draw_text(
            image: Image,
            font_path: Optional[Path],
            text: str,
            color: Union[str, None],
            padding: Tuple[int, ...],
            markata: "Markata",
        ) -> None:
            text = text or ""
            draw = ImageDraw.Draw(image)
            padding = resolve_padding(padding, markata)
            width = image.size[0]
            height = image.size[1]
            bounding_box = [padding[0], padding[1], width - padding[0], height - padding[1]]
            bounding_box = [padding[0], padding[1], width - padding[2], height - padding[3]]
            max_size = (bounding_box[2] - bounding_box[0], bounding_box[3] - bounding_box[1])
            x1, y1, x2, y2 = bounding_box
            font = get_font(font_path, draw, text, max_size=max_size) if font_path else None
            w, h = draw.textsize(text, font=font)
            x = (x2 - x1 - w) / 2 + x1
            y = (y2 - y1 - h) / 2 + y1
            draw.text((x, y), text, fill=color, font=font, align="center")
```


!! function <h2 id='resolve_padding' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>resolve_padding <em class='small'>function</em></h2>
    Convert padding to a len 4 tuple
???+ source "resolve_padding <em class='small'>source</em>"

```python

        def resolve_padding(padding: Tuple[int, ...], markata: "Markata") -> Tuple[int, ...]:
            """Convert padding to a len 4 tuple"""
            if len(padding) == 4:
                return padding
            if len(padding) == 3:
                return (*padding, padding[1])
            if len(padding) == 2:
                return padding * 2
            if len(padding) == 1:
                return padding * 4
            raise PaddingError(f"recieved padding: {padding}")
```


!! function <h2 id='make_cover' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>make_cover <em class='small'>function</em></h2>

???+ source "make_cover <em class='small'>source</em>"

```python

        def make_cover(
            title: str,
            color: str,
            output_path: Path,
            template_path: Path,
            font_path: Optional[Path],
            padding: Tuple[int, ...],
            text_font: Path,
            text: str = None,
            text_font_color: str = None,
            text_padding: Tuple[int, ...] = None,
            resizes: List[int] = None,
            markata: "Markata" = None,
        ) -> None:
            if output_path.exists():
                return
            image = Image.open(template_path) if template_path else Image.new("RGB", (800, 450))

            draw_text(
                image=image,
                font_path=font_path,
                title=title,
                color=color,
                padding=padding,
                markata=markata,
            )
            if text is not None:
                if text_padding is None:
                    text_padding = (
                        image.size[1] - image.size[1] / 5,
                        image.size[0] / 5,
                        image.size[1] - image.size[1] / 10,
                    )
                draw_text(image, text_font, text, text_font_color, text_padding)

            image.save(output_path)
            ratio = image.size[1] / image.size[0]

            covers = []
            if resizes:
                for width in resizes:
                    re_img = image.resize((width, int(width * ratio)), Image.ANTIALIAS)
                    filename = (
                        f"{output_path.stem}_{width}x{int(width*ratio)}{output_path.suffix}"
                    )
                    covers.append(filename)

                    filepath = Path(output_path.parent / filename)
                    re_img.save(filepath)
```


!! function <h2 id='save' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>save <em class='small'>function</em></h2>

???+ source "save <em class='small'>source</em>"

```python

        def save(markata: "Markata") -> None:
            futures = []

            if "covers" not in markata.config.keys():
                return

            for article in markata.iter_articles("making covers"):
                for cover in markata.config["covers"]:
                    try:
                        padding = cover["padding"]
                    except KeyError:
                        padding = (
                            200,
                            100,
                        )
                    try:
                        text_padding = cover["text_padding"]
                    except KeyError:
                        text_padding = (
                            200,
                            100,
                        )
                    if "text_key" in cover:
                        try:
                            text = article.metadata[cover["text_key"]]
                        except AttributeError:
                            text = article[cover["text_key"]]
                        try:
                            text = text.replace("\n", "")
                            from more_itertools import chunked

                            text = "\n".join(["".join(c) for c in chunked(text, 60)])
                        except AttributeError:
                            # text is likely None
                            pass

                        text_font = cover["text_font"]
                        text_font_color = cover["text_font_color"]
                    else:
                        text = None
                        text_font = None
                        text_font_color = None
                    try:
                        title = article.metadata["title"]
                    except AttributeError:
                        title = article["title"]
                    futures.append(
                        make_cover(
                            title=title,
                            color=cover["font_color"],
                            output_path=Path(markata.config.output_dir)
                            / (article["slug"] + cover["name"] + ".png"),
                            template_path=cover.get("template", None),
                            font_path=cover.get("font", None),
                            padding=padding,
                            text_font=text_font,
                            text=text,
                            text_font_color=text_font_color,
                            text_padding=text_padding,
                            resizes=cover.get("resizes"),
                            markata=markata,
                        ),
                    )

            progress = Progress(
                BarColumn(bar_width=None),
                transient=True,
                console=markata.console,
            )
            task_id = progress.add_task("loading markdown")
            progress.update(task_id, total=len(futures))
            with progress:
                while not all(f.done() for f in futures):
                    time.sleep(0.1)
                    progress.update(task_id, total=len([f for f in futures if f.done()]))
            [f.result() for f in futures]
```


!! method <h2 id='__init__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__init__ <em class='small'>method</em></h2>

???+ source "__init__ <em class='small'>source</em>"

```python

        def __init__(
                self,
                msg: str = "",
            ) -> None:
                super().__init__(
                    "Padding must be an iterable of length 1, 2, 3, or 4.\n" + msg,
                )
```
