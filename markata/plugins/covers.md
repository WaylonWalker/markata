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
  of length 1, 2, 3, or 4.\\n\" + msg,\n                )\n```\n\n"
date: 0001-01-01
description: Example configuration.  Covers supports multiple covers to be configured.  Here
  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ! ???+ source  ! ?
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Covers.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"Example configuration.  Covers supports multiple
    covers to be configured.  Here ! ???+ source  ! ???+ source  ! ???+ source  !
    ???+ source  ! ! ???+ source  ! ?\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Covers.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Example configuration.  Covers supports
    multiple covers to be configured.  Here ! ???+ source  ! ???+ source  ! ???+ source
    \ ! ???+ source  ! ! ???+ source  ! ?\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Covers.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <h3>Configuration</h3>\n<p>Example
    configuration.  Covers supports multiple covers to be configured.  Here\nis an
    example from my blog where we have a template sized for <a href=\"http://dev.to\">dev.to</a>
    and one\nsized for open graph.  Each image takes it's own configuration.</p>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.covers]]</span>\n<span
    class=\"n\">name</span><span class=\"o\">=</span><span class=\"s1\">&#39;-dev&#39;</span>\n<span
    class=\"n\">template</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;static/cover-template.png&quot;</span>\n<span
    class=\"n\">font</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;./static/JosefinSans-Regular.ttf&quot;</span>\n<span
    class=\"n\">text_font</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;./static/JosefinSans-Regular.ttf&quot;</span>\n<span
    class=\"n\">font_color</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;rgb(185,155,165)&quot;</span>\n<span
    class=\"n\">text_font_color</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;rgb(255,255,255)&quot;</span>\n<span
    class=\"n\">text_key</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s1\">&#39;description&#39;</span>\n<span class=\"n\">padding</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"p\">[</span><span class=\"mi\">0</span><span class=\"p\">,</span><span
    class=\"w\"> </span><span class=\"mi\">40</span><span class=\"p\">,</span><span
    class=\"w\"> </span><span class=\"mi\">100</span><span class=\"p\">,</span><span
    class=\"w\"> </span><span class=\"mi\">300</span><span class=\"p\">]</span>\n<span
    class=\"n\">text_padding</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"p\">[</span><span class=\"mi\">0</span><span
    class=\"p\">,</span><span class=\"mi\">0</span><span class=\"p\">]</span>\n\n<span
    class=\"k\">[[markata.covers]]</span>\n<span class=\"n\">name</span><span class=\"o\">=</span><span
    class=\"s1\">&#39;&#39;</span>\n<span class=\"n\">template</span><span class=\"w\">
    </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s2\">&quot;static/og-template.png&quot;</span>\n<span
    class=\"n\">font</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;./static/JosefinSans-Regular.ttf&quot;</span>\n<span
    class=\"n\">font_color</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;rgb(255,255,255)&quot;</span>\n<span
    class=\"n\">text_font</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;./static/JosefinSans-Regular.ttf&quot;</span>\n<span
    class=\"n\">text_font_color</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;rgb(200,200,200)&quot;</span>\n<span
    class=\"n\">text_key</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s1\">&#39;description&#39;</span>\n<span class=\"n\">padding</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"p\">[</span><span class=\"mi\">10</span><span class=\"p\">,</span><span
    class=\"w\"> </span><span class=\"mi\">10</span><span class=\"p\">,</span><span
    class=\"w\"> </span><span class=\"mi\">100</span><span class=\"p\">,</span><span
    class=\"w\"> </span><span class=\"mi\">300</span><span class=\"p\">]</span>\n<span
    class=\"n\">text_padding</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"p\">[</span><span class=\"mi\">0</span><span
    class=\"p\">,</span><span class=\"mi\">0</span><span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='_load_font' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>_load_font <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_load_font
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
    <span class=\"nf\">_load_font</span><span class=\"p\">(</span><span class=\"n\">path</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">,</span> <span
    class=\"n\">size</span><span class=\"p\">:</span> <span class=\"nb\">int</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">ImageFont</span><span
    class=\"o\">.</span><span class=\"n\">FreeTypeFont</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">return</span> <span class=\"n\">ImageFont</span><span
    class=\"o\">.</span><span class=\"n\">truetype</span><span class=\"p\">(</span><span
    class=\"n\">path</span><span class=\"p\">,</span> <span class=\"n\">size</span><span
    class=\"o\">=</span><span class=\"n\">size</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='get_font' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>get_font <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">get_font
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
    <span class=\"nf\">get_font</span><span class=\"p\">(</span>\n            <span
    class=\"n\">path</span><span class=\"p\">:</span> <span class=\"n\">Path</span><span
    class=\"p\">,</span>\n            <span class=\"n\">draw</span><span class=\"p\">:</span>
    <span class=\"n\">ImageDraw</span><span class=\"o\">.</span><span class=\"n\">Draw</span><span
    class=\"p\">,</span>\n            <span class=\"n\">title</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span><span class=\"p\">,</span>\n            <span class=\"n\">size</span><span
    class=\"p\">:</span> <span class=\"nb\">int</span> <span class=\"o\">=</span>
    <span class=\"mi\">250</span><span class=\"p\">,</span>\n            <span class=\"n\">max_size</span><span
    class=\"p\">:</span> <span class=\"nb\">tuple</span> <span class=\"o\">=</span>
    <span class=\"p\">(</span><span class=\"mi\">800</span><span class=\"p\">,</span>
    <span class=\"mi\">220</span><span class=\"p\">),</span>\n        <span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">ImageFont</span><span class=\"o\">.</span><span
    class=\"n\">FreeTypeFont</span><span class=\"p\">:</span>\n            <span class=\"n\">title</span>
    <span class=\"o\">=</span> <span class=\"n\">title</span> <span class=\"ow\">or</span>
    <span class=\"s2\">&quot;&quot;</span>\n            <span class=\"n\">font</span>
    <span class=\"o\">=</span> <span class=\"n\">_load_font</span><span class=\"p\">(</span><span
    class=\"n\">path</span><span class=\"p\">,</span> <span class=\"n\">size</span><span
    class=\"p\">)</span>\n            <span class=\"n\">current_size</span> <span
    class=\"o\">=</span> <span class=\"n\">draw</span><span class=\"o\">.</span><span
    class=\"n\">textsize</span><span class=\"p\">(</span><span class=\"n\">title</span><span
    class=\"p\">,</span> <span class=\"n\">font</span><span class=\"o\">=</span><span
    class=\"n\">font</span><span class=\"p\">)</span>\n\n            <span class=\"k\">if</span>
    <span class=\"n\">current_size</span><span class=\"p\">[</span><span class=\"mi\">0</span><span
    class=\"p\">]</span> <span class=\"o\">&gt;</span> <span class=\"n\">max_size</span><span
    class=\"p\">[</span><span class=\"mi\">0</span><span class=\"p\">]</span> <span
    class=\"ow\">or</span> <span class=\"n\">current_size</span><span class=\"p\">[</span><span
    class=\"mi\">1</span><span class=\"p\">]</span> <span class=\"o\">&gt;</span>
    <span class=\"n\">max_size</span><span class=\"p\">[</span><span class=\"mi\">1</span><span
    class=\"p\">]:</span>\n                <span class=\"k\">return</span> <span class=\"n\">get_font</span><span
    class=\"p\">(</span><span class=\"n\">path</span><span class=\"p\">,</span> <span
    class=\"n\">draw</span><span class=\"p\">,</span> <span class=\"n\">title</span><span
    class=\"p\">,</span> <span class=\"n\">size</span> <span class=\"o\">-</span>
    <span class=\"mi\">10</span><span class=\"p\">,</span> <span class=\"n\">max_size</span><span
    class=\"o\">=</span><span class=\"n\">max_size</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">return</span> <span class=\"n\">font</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='PaddingError' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>PaddingError <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">PaddingError
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
    <span class=\"nc\">PaddingError</span><span class=\"p\">(</span><span class=\"ne\">BaseException</span><span
    class=\"p\">):</span>\n            <span class=\"k\">def</span> <span class=\"fm\">__init__</span><span
    class=\"p\">(</span>\n                <span class=\"bp\">self</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">msg</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;&quot;</span><span class=\"p\">,</span>\n
    \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"nb\">super</span><span class=\"p\">()</span><span
    class=\"o\">.</span><span class=\"fm\">__init__</span><span class=\"p\">(</span>\n
    \                   <span class=\"s2\">&quot;Padding must be an iterable of length
    1, 2, 3, or 4.</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span>
    <span class=\"o\">+</span> <span class=\"n\">msg</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! function
    <h2 id='draw_text' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>draw_text
    <em class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">draw_text <em class='small'>source</em></p>\n</div>\n<pre
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
    class=\"p\">,</span> <span class=\"o\">...</span><span class=\"p\">],</span>\n
    \           <span class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span>\n        <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"n\">text</span>
    <span class=\"o\">=</span> <span class=\"n\">text</span> <span class=\"ow\">or</span>
    <span class=\"s2\">&quot;&quot;</span>\n            <span class=\"n\">draw</span>
    <span class=\"o\">=</span> <span class=\"n\">ImageDraw</span><span class=\"o\">.</span><span
    class=\"n\">Draw</span><span class=\"p\">(</span><span class=\"n\">image</span><span
    class=\"p\">)</span>\n            <span class=\"n\">padding</span> <span class=\"o\">=</span>
    <span class=\"n\">resolve_padding</span><span class=\"p\">(</span><span class=\"n\">padding</span><span
    class=\"p\">,</span> <span class=\"n\">markata</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">width</span> <span class=\"o\">=</span> <span class=\"n\">image</span><span
    class=\"o\">.</span><span class=\"n\">size</span><span class=\"p\">[</span><span
    class=\"mi\">0</span><span class=\"p\">]</span>\n            <span class=\"n\">height</span>
    <span class=\"o\">=</span> <span class=\"n\">image</span><span class=\"o\">.</span><span
    class=\"n\">size</span><span class=\"p\">[</span><span class=\"mi\">1</span><span
    class=\"p\">]</span>\n            <span class=\"n\">bounding_box</span> <span
    class=\"o\">=</span> <span class=\"p\">[</span><span class=\"n\">padding</span><span
    class=\"p\">[</span><span class=\"mi\">0</span><span class=\"p\">],</span> <span
    class=\"n\">padding</span><span class=\"p\">[</span><span class=\"mi\">1</span><span
    class=\"p\">],</span> <span class=\"n\">width</span> <span class=\"o\">-</span>
    <span class=\"n\">padding</span><span class=\"p\">[</span><span class=\"mi\">0</span><span
    class=\"p\">],</span> <span class=\"n\">height</span> <span class=\"o\">-</span>
    <span class=\"n\">padding</span><span class=\"p\">[</span><span class=\"mi\">1</span><span
    class=\"p\">]]</span>\n            <span class=\"n\">bounding_box</span> <span
    class=\"o\">=</span> <span class=\"p\">[</span><span class=\"n\">padding</span><span
    class=\"p\">[</span><span class=\"mi\">0</span><span class=\"p\">],</span> <span
    class=\"n\">padding</span><span class=\"p\">[</span><span class=\"mi\">1</span><span
    class=\"p\">],</span> <span class=\"n\">width</span> <span class=\"o\">-</span>
    <span class=\"n\">padding</span><span class=\"p\">[</span><span class=\"mi\">2</span><span
    class=\"p\">],</span> <span class=\"n\">height</span> <span class=\"o\">-</span>
    <span class=\"n\">padding</span><span class=\"p\">[</span><span class=\"mi\">3</span><span
    class=\"p\">]]</span>\n            <span class=\"n\">max_size</span> <span class=\"o\">=</span>
    <span class=\"p\">(</span><span class=\"n\">bounding_box</span><span class=\"p\">[</span><span
    class=\"mi\">2</span><span class=\"p\">]</span> <span class=\"o\">-</span> <span
    class=\"n\">bounding_box</span><span class=\"p\">[</span><span class=\"mi\">0</span><span
    class=\"p\">],</span> <span class=\"n\">bounding_box</span><span class=\"p\">[</span><span
    class=\"mi\">3</span><span class=\"p\">]</span> <span class=\"o\">-</span> <span
    class=\"n\">bounding_box</span><span class=\"p\">[</span><span class=\"mi\">1</span><span
    class=\"p\">])</span>\n            <span class=\"n\">x1</span><span class=\"p\">,</span>
    <span class=\"n\">y1</span><span class=\"p\">,</span> <span class=\"n\">x2</span><span
    class=\"p\">,</span> <span class=\"n\">y2</span> <span class=\"o\">=</span> <span
    class=\"n\">bounding_box</span>\n            <span class=\"n\">font</span> <span
    class=\"o\">=</span> <span class=\"n\">get_font</span><span class=\"p\">(</span><span
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
    class=\"s2\">&quot;center&quot;</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='resolve_padding' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>resolve_padding <em class='small'>function</em></h2>\nConvert padding to
    a len 4 tuple</p>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">resolve_padding <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">resolve_padding</span><span class=\"p\">(</span><span class=\"n\">padding</span><span
    class=\"p\">:</span> <span class=\"n\">Tuple</span><span class=\"p\">[</span><span
    class=\"nb\">int</span><span class=\"p\">,</span> <span class=\"o\">...</span><span
    class=\"p\">],</span> <span class=\"n\">markata</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Tuple</span><span class=\"p\">[</span><span
    class=\"nb\">int</span><span class=\"p\">,</span> <span class=\"o\">...</span><span
    class=\"p\">]:</span>\n<span class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;Convert
    padding to a len 4 tuple&quot;&quot;&quot;</span>\n            <span class=\"k\">if</span>
    <span class=\"nb\">len</span><span class=\"p\">(</span><span class=\"n\">padding</span><span
    class=\"p\">)</span> <span class=\"o\">==</span> <span class=\"mi\">4</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"n\">padding</span>\n
    \           <span class=\"k\">if</span> <span class=\"nb\">len</span><span class=\"p\">(</span><span
    class=\"n\">padding</span><span class=\"p\">)</span> <span class=\"o\">==</span>
    <span class=\"mi\">3</span><span class=\"p\">:</span>\n                <span class=\"k\">return</span>
    <span class=\"p\">(</span><span class=\"o\">*</span><span class=\"n\">padding</span><span
    class=\"p\">,</span> <span class=\"n\">padding</span><span class=\"p\">[</span><span
    class=\"mi\">1</span><span class=\"p\">])</span>\n            <span class=\"k\">if</span>
    <span class=\"nb\">len</span><span class=\"p\">(</span><span class=\"n\">padding</span><span
    class=\"p\">)</span> <span class=\"o\">==</span> <span class=\"mi\">2</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"n\">padding</span>
    <span class=\"o\">*</span> <span class=\"mi\">2</span>\n            <span class=\"k\">if</span>
    <span class=\"nb\">len</span><span class=\"p\">(</span><span class=\"n\">padding</span><span
    class=\"p\">)</span> <span class=\"o\">==</span> <span class=\"mi\">1</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"n\">padding</span>
    <span class=\"o\">*</span> <span class=\"mi\">4</span>\n            <span class=\"k\">raise</span>
    <span class=\"n\">PaddingError</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;recieved padding: </span><span class=\"si\">{</span><span class=\"n\">padding</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='make_cover' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>make_cover <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">make_cover
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
    <span class=\"nf\">make_cover</span><span class=\"p\">(</span>\n            <span
    class=\"n\">title</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">,</span>\n            <span class=\"n\">color</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span><span class=\"p\">,</span>\n            <span class=\"n\">output_path</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">,</span>\n
    \           <span class=\"n\">template_path</span><span class=\"p\">:</span> <span
    class=\"n\">Path</span><span class=\"p\">,</span>\n            <span class=\"n\">font_path</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"n\">Path</span><span class=\"p\">],</span>\n            <span class=\"n\">padding</span><span
    class=\"p\">:</span> <span class=\"n\">Tuple</span><span class=\"p\">[</span><span
    class=\"nb\">int</span><span class=\"p\">,</span> <span class=\"o\">...</span><span
    class=\"p\">],</span>\n            <span class=\"n\">text_font</span><span class=\"p\">:</span>
    <span class=\"n\">Path</span><span class=\"p\">,</span>\n            <span class=\"n\">text</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span> <span class=\"o\">=</span>
    <span class=\"kc\">None</span><span class=\"p\">,</span>\n            <span class=\"n\">text_font_color</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span> <span class=\"o\">=</span>
    <span class=\"kc\">None</span><span class=\"p\">,</span>\n            <span class=\"n\">text_padding</span><span
    class=\"p\">:</span> <span class=\"n\">Tuple</span><span class=\"p\">[</span><span
    class=\"nb\">int</span><span class=\"p\">,</span> <span class=\"o\">...</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span><span
    class=\"p\">,</span>\n            <span class=\"n\">resizes</span><span class=\"p\">:</span>
    <span class=\"n\">List</span><span class=\"p\">[</span><span class=\"nb\">int</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span><span
    class=\"p\">,</span>\n            <span class=\"n\">markata</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;Markata&quot;</span> <span class=\"o\">=</span> <span
    class=\"kc\">None</span><span class=\"p\">,</span>\n        <span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">if</span> <span class=\"n\">output_path</span><span
    class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">():</span>\n
    \               <span class=\"k\">return</span>\n            <span class=\"n\">image</span>
    <span class=\"o\">=</span> <span class=\"n\">Image</span><span class=\"o\">.</span><span
    class=\"n\">open</span><span class=\"p\">(</span><span class=\"n\">template_path</span><span
    class=\"p\">)</span> <span class=\"k\">if</span> <span class=\"n\">template_path</span>
    <span class=\"k\">else</span> <span class=\"n\">Image</span><span class=\"o\">.</span><span
    class=\"n\">new</span><span class=\"p\">(</span><span class=\"s2\">&quot;RGB&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">(</span><span class=\"mi\">800</span><span
    class=\"p\">,</span> <span class=\"mi\">450</span><span class=\"p\">))</span>\n\n
    \           <span class=\"n\">draw_text</span><span class=\"p\">(</span>\n                <span
    class=\"n\">image</span><span class=\"o\">=</span><span class=\"n\">image</span><span
    class=\"p\">,</span>\n                <span class=\"n\">font_path</span><span
    class=\"o\">=</span><span class=\"n\">font_path</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">title</span><span class=\"o\">=</span><span
    class=\"n\">title</span><span class=\"p\">,</span>\n                <span class=\"n\">color</span><span
    class=\"o\">=</span><span class=\"n\">color</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">padding</span><span class=\"o\">=</span><span
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
    class=\"mi\">5</span><span class=\"p\">,</span>\n                        <span
    class=\"n\">image</span><span class=\"o\">.</span><span class=\"n\">size</span><span
    class=\"p\">[</span><span class=\"mi\">0</span><span class=\"p\">]</span> <span
    class=\"o\">/</span> <span class=\"mi\">5</span><span class=\"p\">,</span>\n                        <span
    class=\"n\">image</span><span class=\"o\">.</span><span class=\"n\">size</span><span
    class=\"p\">[</span><span class=\"mi\">1</span><span class=\"p\">]</span> <span
    class=\"o\">-</span> <span class=\"n\">image</span><span class=\"o\">.</span><span
    class=\"n\">size</span><span class=\"p\">[</span><span class=\"mi\">1</span><span
    class=\"p\">]</span> <span class=\"o\">/</span> <span class=\"mi\">10</span><span
    class=\"p\">,</span>\n                    <span class=\"p\">)</span>\n                <span
    class=\"n\">draw_text</span><span class=\"p\">(</span><span class=\"n\">image</span><span
    class=\"p\">,</span> <span class=\"n\">text_font</span><span class=\"p\">,</span>
    <span class=\"n\">text</span><span class=\"p\">,</span> <span class=\"n\">text_font_color</span><span
    class=\"p\">,</span> <span class=\"n\">text_padding</span><span class=\"p\">)</span>\n\n
    \           <span class=\"n\">image</span><span class=\"o\">.</span><span class=\"n\">save</span><span
    class=\"p\">(</span><span class=\"n\">output_path</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">ratio</span> <span class=\"o\">=</span> <span class=\"n\">image</span><span
    class=\"o\">.</span><span class=\"n\">size</span><span class=\"p\">[</span><span
    class=\"mi\">1</span><span class=\"p\">]</span> <span class=\"o\">/</span> <span
    class=\"n\">image</span><span class=\"o\">.</span><span class=\"n\">size</span><span
    class=\"p\">[</span><span class=\"mi\">0</span><span class=\"p\">]</span>\n\n
    \           <span class=\"n\">covers</span> <span class=\"o\">=</span> <span class=\"p\">[]</span>\n
    \           <span class=\"k\">if</span> <span class=\"n\">resizes</span><span
    class=\"p\">:</span>\n                <span class=\"k\">for</span> <span class=\"n\">width</span>
    <span class=\"ow\">in</span> <span class=\"n\">resizes</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">re_img</span> <span class=\"o\">=</span>
    <span class=\"n\">image</span><span class=\"o\">.</span><span class=\"n\">resize</span><span
    class=\"p\">((</span><span class=\"n\">width</span><span class=\"p\">,</span>
    <span class=\"nb\">int</span><span class=\"p\">(</span><span class=\"n\">width</span>
    <span class=\"o\">*</span> <span class=\"n\">ratio</span><span class=\"p\">)),</span>
    <span class=\"n\">Image</span><span class=\"o\">.</span><span class=\"n\">ANTIALIAS</span><span
    class=\"p\">)</span>\n                    <span class=\"n\">filename</span> <span
    class=\"o\">=</span> <span class=\"p\">(</span>\n                        <span
    class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">output_path</span><span class=\"o\">.</span><span class=\"n\">stem</span><span
    class=\"si\">}</span><span class=\"s2\">_</span><span class=\"si\">{</span><span
    class=\"n\">width</span><span class=\"si\">}</span><span class=\"s2\">x</span><span
    class=\"si\">{</span><span class=\"nb\">int</span><span class=\"p\">(</span><span
    class=\"n\">width</span><span class=\"o\">*</span><span class=\"n\">ratio</span><span
    class=\"p\">)</span><span class=\"si\">}{</span><span class=\"n\">output_path</span><span
    class=\"o\">.</span><span class=\"n\">suffix</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span>\n                    <span class=\"p\">)</span>\n                    <span
    class=\"n\">covers</span><span class=\"o\">.</span><span class=\"n\">append</span><span
    class=\"p\">(</span><span class=\"n\">filename</span><span class=\"p\">)</span>\n\n
    \                   <span class=\"n\">filepath</span> <span class=\"o\">=</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">output_path</span><span
    class=\"o\">.</span><span class=\"n\">parent</span> <span class=\"o\">/</span>
    <span class=\"n\">filename</span><span class=\"p\">)</span>\n                    <span
    class=\"n\">re_img</span><span class=\"o\">.</span><span class=\"n\">save</span><span
    class=\"p\">(</span><span class=\"n\">filepath</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='save' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>save <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">save
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
    <span class=\"nf\">save</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">futures</span> <span class=\"o\">=</span> <span
    class=\"p\">[]</span>\n\n            <span class=\"k\">if</span> <span class=\"s2\">&quot;covers&quot;</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">keys</span><span class=\"p\">():</span>\n                <span class=\"k\">return</span>\n\n
    \           <span class=\"k\">for</span> <span class=\"n\">article</span> <span
    class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">iter_articles</span><span class=\"p\">(</span><span class=\"s2\">&quot;making
    covers&quot;</span><span class=\"p\">):</span>\n                <span class=\"k\">for</span>
    <span class=\"n\">cover</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;covers&quot;</span><span class=\"p\">]:</span>\n                    <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">padding</span> <span class=\"o\">=</span> <span class=\"n\">cover</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;padding&quot;</span><span class=\"p\">]</span>\n
    \                   <span class=\"k\">except</span> <span class=\"ne\">KeyError</span><span
    class=\"p\">:</span>\n                        <span class=\"n\">padding</span>
    <span class=\"o\">=</span> <span class=\"p\">(</span>\n                            <span
    class=\"mi\">200</span><span class=\"p\">,</span>\n                            <span
    class=\"mi\">100</span><span class=\"p\">,</span>\n                        <span
    class=\"p\">)</span>\n                    <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">text_padding</span> <span class=\"o\">=</span>
    <span class=\"n\">cover</span><span class=\"p\">[</span><span class=\"s2\">&quot;text_padding&quot;</span><span
    class=\"p\">]</span>\n                    <span class=\"k\">except</span> <span
    class=\"ne\">KeyError</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">text_padding</span> <span class=\"o\">=</span> <span class=\"p\">(</span>\n
    \                           <span class=\"mi\">200</span><span class=\"p\">,</span>\n
    \                           <span class=\"mi\">100</span><span class=\"p\">,</span>\n
    \                       <span class=\"p\">)</span>\n                    <span
    class=\"k\">if</span> <span class=\"s2\">&quot;text_key&quot;</span> <span class=\"ow\">in</span>
    <span class=\"n\">cover</span><span class=\"p\">:</span>\n                        <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                            <span
    class=\"n\">text</span> <span class=\"o\">=</span> <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">metadata</span><span class=\"p\">[</span><span
    class=\"n\">cover</span><span class=\"p\">[</span><span class=\"s2\">&quot;text_key&quot;</span><span
    class=\"p\">]]</span>\n                        <span class=\"k\">except</span>
    <span class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n                            <span
    class=\"n\">text</span> <span class=\"o\">=</span> <span class=\"n\">article</span><span
    class=\"p\">[</span><span class=\"n\">cover</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;text_key&quot;</span><span class=\"p\">]]</span>\n                        <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                            <span
    class=\"n\">text</span> <span class=\"o\">=</span> <span class=\"n\">text</span><span
    class=\"o\">.</span><span class=\"n\">replace</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;&quot;</span><span class=\"p\">)</span>\n
    \                           <span class=\"kn\">from</span> <span class=\"nn\">more_itertools</span>
    <span class=\"kn\">import</span> <span class=\"n\">chunked</span>\n\n                            <span
    class=\"n\">text</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;</span><span
    class=\"se\">\\n</span><span class=\"s2\">&quot;</span><span class=\"o\">.</span><span
    class=\"n\">join</span><span class=\"p\">([</span><span class=\"s2\">&quot;&quot;</span><span
    class=\"o\">.</span><span class=\"n\">join</span><span class=\"p\">(</span><span
    class=\"n\">c</span><span class=\"p\">)</span> <span class=\"k\">for</span> <span
    class=\"n\">c</span> <span class=\"ow\">in</span> <span class=\"n\">chunked</span><span
    class=\"p\">(</span><span class=\"n\">text</span><span class=\"p\">,</span> <span
    class=\"mi\">60</span><span class=\"p\">)])</span>\n                        <span
    class=\"k\">except</span> <span class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n
    \                           <span class=\"c1\"># text is likely None</span>\n
    \                           <span class=\"k\">pass</span>\n\n                        <span
    class=\"n\">text_font</span> <span class=\"o\">=</span> <span class=\"n\">cover</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;text_font&quot;</span><span class=\"p\">]</span>\n
    \                       <span class=\"n\">text_font_color</span> <span class=\"o\">=</span>
    <span class=\"n\">cover</span><span class=\"p\">[</span><span class=\"s2\">&quot;text_font_color&quot;</span><span
    class=\"p\">]</span>\n                    <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">text</span> <span class=\"o\">=</span>
    <span class=\"kc\">None</span>\n                        <span class=\"n\">text_font</span>
    <span class=\"o\">=</span> <span class=\"kc\">None</span>\n                        <span
    class=\"n\">text_font_color</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n
    \                   <span class=\"k\">try</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">title</span> <span class=\"o\">=</span> <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">metadata</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;title&quot;</span><span class=\"p\">]</span>\n                    <span
    class=\"k\">except</span> <span class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">title</span> <span class=\"o\">=</span>
    <span class=\"n\">article</span><span class=\"p\">[</span><span class=\"s2\">&quot;title&quot;</span><span
    class=\"p\">]</span>\n                    <span class=\"n\">futures</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span>\n
    \                       <span class=\"n\">make_cover</span><span class=\"p\">(</span>\n
    \                           <span class=\"n\">title</span><span class=\"o\">=</span><span
    class=\"n\">title</span><span class=\"p\">,</span>\n                            <span
    class=\"n\">color</span><span class=\"o\">=</span><span class=\"n\">cover</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;font_color&quot;</span><span class=\"p\">],</span>\n
    \                           <span class=\"n\">output_path</span><span class=\"o\">=</span><span
    class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">output_dir</span><span class=\"p\">)</span>\n                            <span
    class=\"o\">/</span> <span class=\"p\">(</span><span class=\"n\">article</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">+</span> <span class=\"n\">cover</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;name&quot;</span><span class=\"p\">]</span> <span class=\"o\">+</span>
    <span class=\"s2\">&quot;.png&quot;</span><span class=\"p\">),</span>\n                            <span
    class=\"n\">template_path</span><span class=\"o\">=</span><span class=\"n\">cover</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;template&quot;</span><span class=\"p\">,</span> <span class=\"kc\">None</span><span
    class=\"p\">),</span>\n                            <span class=\"n\">font_path</span><span
    class=\"o\">=</span><span class=\"n\">cover</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;font&quot;</span><span
    class=\"p\">,</span> <span class=\"kc\">None</span><span class=\"p\">),</span>\n
    \                           <span class=\"n\">padding</span><span class=\"o\">=</span><span
    class=\"n\">padding</span><span class=\"p\">,</span>\n                            <span
    class=\"n\">text_font</span><span class=\"o\">=</span><span class=\"n\">text_font</span><span
    class=\"p\">,</span>\n                            <span class=\"n\">text</span><span
    class=\"o\">=</span><span class=\"n\">text</span><span class=\"p\">,</span>\n
    \                           <span class=\"n\">text_font_color</span><span class=\"o\">=</span><span
    class=\"n\">text_font_color</span><span class=\"p\">,</span>\n                            <span
    class=\"n\">text_padding</span><span class=\"o\">=</span><span class=\"n\">text_padding</span><span
    class=\"p\">,</span>\n                            <span class=\"n\">resizes</span><span
    class=\"o\">=</span><span class=\"n\">cover</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;resizes&quot;</span><span
    class=\"p\">),</span>\n                            <span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"p\">,</span>\n
    \                       <span class=\"p\">),</span>\n                    <span
    class=\"p\">)</span>\n\n            <span class=\"n\">progress</span> <span class=\"o\">=</span>
    <span class=\"n\">Progress</span><span class=\"p\">(</span>\n                <span
    class=\"n\">BarColumn</span><span class=\"p\">(</span><span class=\"n\">bar_width</span><span
    class=\"o\">=</span><span class=\"kc\">None</span><span class=\"p\">),</span>\n
    \               <span class=\"n\">transient</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n                <span class=\"n\">console</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"p\">,</span>\n            <span class=\"p\">)</span>\n
    \           <span class=\"n\">task_id</span> <span class=\"o\">=</span> <span
    class=\"n\">progress</span><span class=\"o\">.</span><span class=\"n\">add_task</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;loading markdown&quot;</span><span
    class=\"p\">)</span>\n            <span class=\"n\">progress</span><span class=\"o\">.</span><span
    class=\"n\">update</span><span class=\"p\">(</span><span class=\"n\">task_id</span><span
    class=\"p\">,</span> <span class=\"n\">total</span><span class=\"o\">=</span><span
    class=\"nb\">len</span><span class=\"p\">(</span><span class=\"n\">futures</span><span
    class=\"p\">))</span>\n            <span class=\"k\">with</span> <span class=\"n\">progress</span><span
    class=\"p\">:</span>\n                <span class=\"k\">while</span> <span class=\"ow\">not</span>
    <span class=\"nb\">all</span><span class=\"p\">(</span><span class=\"n\">f</span><span
    class=\"o\">.</span><span class=\"n\">done</span><span class=\"p\">()</span> <span
    class=\"k\">for</span> <span class=\"n\">f</span> <span class=\"ow\">in</span>
    <span class=\"n\">futures</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">time</span><span class=\"o\">.</span><span class=\"n\">sleep</span><span
    class=\"p\">(</span><span class=\"mf\">0.1</span><span class=\"p\">)</span>\n
    \                   <span class=\"n\">progress</span><span class=\"o\">.</span><span
    class=\"n\">update</span><span class=\"p\">(</span><span class=\"n\">task_id</span><span
    class=\"p\">,</span> <span class=\"n\">total</span><span class=\"o\">=</span><span
    class=\"nb\">len</span><span class=\"p\">([</span><span class=\"n\">f</span> <span
    class=\"k\">for</span> <span class=\"n\">f</span> <span class=\"ow\">in</span>
    <span class=\"n\">futures</span> <span class=\"k\">if</span> <span class=\"n\">f</span><span
    class=\"o\">.</span><span class=\"n\">done</span><span class=\"p\">()]))</span>\n
    \           <span class=\"p\">[</span><span class=\"n\">f</span><span class=\"o\">.</span><span
    class=\"n\">result</span><span class=\"p\">()</span> <span class=\"k\">for</span>
    <span class=\"n\">f</span> <span class=\"ow\">in</span> <span class=\"n\">futures</span><span
    class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='__init__'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'><strong>init</strong>
    <em class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\"><strong>init</strong> <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"fm\">__init__</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">self</span><span class=\"p\">,</span>\n                <span class=\"n\">msg</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;&quot;</span><span class=\"p\">,</span>\n            <span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"nb\">super</span><span class=\"p\">()</span><span
    class=\"o\">.</span><span class=\"fm\">__init__</span><span class=\"p\">(</span>\n
    \                   <span class=\"s2\">&quot;Padding must be an iterable of length
    1, 2, 3, or 4.</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span>
    <span class=\"o\">+</span> <span class=\"n\">msg</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>
    \   </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Covers.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"Example configuration.  Covers supports multiple
    covers to be configured.  Here ! ???+ source  ! ???+ source  ! ???+ source  !
    ???+ source  ! ! ???+ source  ! ?\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Covers.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Example configuration.  Covers supports
    multiple covers to be configured.  Here ! ???+ source  ! ???+ source  ! ???+ source
    \ ! ???+ source  ! ! ???+ source  ! ?\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<article style=\"text-align: center;\">\n    <style>\n        section
    {\n            font-size: 200%;\n        }\n\n\n        .edit {\n            display:
    none;\n        }\n    </style>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Covers.Py\n    </h1>\n</section></article>\n     </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Covers.Py\n    </h1>\n</section>    <section class=\"body\">\n        <h3>Configuration</h3>\n<p>Example
    configuration.  Covers supports multiple covers to be configured.  Here\nis an
    example from my blog where we have a template sized for <a href=\"http://dev.to\">dev.to</a>
    and one\nsized for open graph.  Each image takes it's own configuration.</p>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.covers]]</span>\n<span
    class=\"n\">name</span><span class=\"o\">=</span><span class=\"s1\">&#39;-dev&#39;</span>\n<span
    class=\"n\">template</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;static/cover-template.png&quot;</span>\n<span
    class=\"n\">font</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;./static/JosefinSans-Regular.ttf&quot;</span>\n<span
    class=\"n\">text_font</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;./static/JosefinSans-Regular.ttf&quot;</span>\n<span
    class=\"n\">font_color</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;rgb(185,155,165)&quot;</span>\n<span
    class=\"n\">text_font_color</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;rgb(255,255,255)&quot;</span>\n<span
    class=\"n\">text_key</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s1\">&#39;description&#39;</span>\n<span class=\"n\">padding</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"p\">[</span><span class=\"mi\">0</span><span class=\"p\">,</span><span
    class=\"w\"> </span><span class=\"mi\">40</span><span class=\"p\">,</span><span
    class=\"w\"> </span><span class=\"mi\">100</span><span class=\"p\">,</span><span
    class=\"w\"> </span><span class=\"mi\">300</span><span class=\"p\">]</span>\n<span
    class=\"n\">text_padding</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"p\">[</span><span class=\"mi\">0</span><span
    class=\"p\">,</span><span class=\"mi\">0</span><span class=\"p\">]</span>\n\n<span
    class=\"k\">[[markata.covers]]</span>\n<span class=\"n\">name</span><span class=\"o\">=</span><span
    class=\"s1\">&#39;&#39;</span>\n<span class=\"n\">template</span><span class=\"w\">
    </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s2\">&quot;static/og-template.png&quot;</span>\n<span
    class=\"n\">font</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;./static/JosefinSans-Regular.ttf&quot;</span>\n<span
    class=\"n\">font_color</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;rgb(255,255,255)&quot;</span>\n<span
    class=\"n\">text_font</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;./static/JosefinSans-Regular.ttf&quot;</span>\n<span
    class=\"n\">text_font_color</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;rgb(200,200,200)&quot;</span>\n<span
    class=\"n\">text_key</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s1\">&#39;description&#39;</span>\n<span class=\"n\">padding</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"p\">[</span><span class=\"mi\">10</span><span class=\"p\">,</span><span
    class=\"w\"> </span><span class=\"mi\">10</span><span class=\"p\">,</span><span
    class=\"w\"> </span><span class=\"mi\">100</span><span class=\"p\">,</span><span
    class=\"w\"> </span><span class=\"mi\">300</span><span class=\"p\">]</span>\n<span
    class=\"n\">text_padding</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"p\">[</span><span class=\"mi\">0</span><span
    class=\"p\">,</span><span class=\"mi\">0</span><span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='_load_font' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>_load_font <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_load_font
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
    <span class=\"nf\">_load_font</span><span class=\"p\">(</span><span class=\"n\">path</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">,</span> <span
    class=\"n\">size</span><span class=\"p\">:</span> <span class=\"nb\">int</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">ImageFont</span><span
    class=\"o\">.</span><span class=\"n\">FreeTypeFont</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">return</span> <span class=\"n\">ImageFont</span><span
    class=\"o\">.</span><span class=\"n\">truetype</span><span class=\"p\">(</span><span
    class=\"n\">path</span><span class=\"p\">,</span> <span class=\"n\">size</span><span
    class=\"o\">=</span><span class=\"n\">size</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='get_font' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>get_font <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">get_font
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
    <span class=\"nf\">get_font</span><span class=\"p\">(</span>\n            <span
    class=\"n\">path</span><span class=\"p\">:</span> <span class=\"n\">Path</span><span
    class=\"p\">,</span>\n            <span class=\"n\">draw</span><span class=\"p\">:</span>
    <span class=\"n\">ImageDraw</span><span class=\"o\">.</span><span class=\"n\">Draw</span><span
    class=\"p\">,</span>\n            <span class=\"n\">title</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span><span class=\"p\">,</span>\n            <span class=\"n\">size</span><span
    class=\"p\">:</span> <span class=\"nb\">int</span> <span class=\"o\">=</span>
    <span class=\"mi\">250</span><span class=\"p\">,</span>\n            <span class=\"n\">max_size</span><span
    class=\"p\">:</span> <span class=\"nb\">tuple</span> <span class=\"o\">=</span>
    <span class=\"p\">(</span><span class=\"mi\">800</span><span class=\"p\">,</span>
    <span class=\"mi\">220</span><span class=\"p\">),</span>\n        <span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">ImageFont</span><span class=\"o\">.</span><span
    class=\"n\">FreeTypeFont</span><span class=\"p\">:</span>\n            <span class=\"n\">title</span>
    <span class=\"o\">=</span> <span class=\"n\">title</span> <span class=\"ow\">or</span>
    <span class=\"s2\">&quot;&quot;</span>\n            <span class=\"n\">font</span>
    <span class=\"o\">=</span> <span class=\"n\">_load_font</span><span class=\"p\">(</span><span
    class=\"n\">path</span><span class=\"p\">,</span> <span class=\"n\">size</span><span
    class=\"p\">)</span>\n            <span class=\"n\">current_size</span> <span
    class=\"o\">=</span> <span class=\"n\">draw</span><span class=\"o\">.</span><span
    class=\"n\">textsize</span><span class=\"p\">(</span><span class=\"n\">title</span><span
    class=\"p\">,</span> <span class=\"n\">font</span><span class=\"o\">=</span><span
    class=\"n\">font</span><span class=\"p\">)</span>\n\n            <span class=\"k\">if</span>
    <span class=\"n\">current_size</span><span class=\"p\">[</span><span class=\"mi\">0</span><span
    class=\"p\">]</span> <span class=\"o\">&gt;</span> <span class=\"n\">max_size</span><span
    class=\"p\">[</span><span class=\"mi\">0</span><span class=\"p\">]</span> <span
    class=\"ow\">or</span> <span class=\"n\">current_size</span><span class=\"p\">[</span><span
    class=\"mi\">1</span><span class=\"p\">]</span> <span class=\"o\">&gt;</span>
    <span class=\"n\">max_size</span><span class=\"p\">[</span><span class=\"mi\">1</span><span
    class=\"p\">]:</span>\n                <span class=\"k\">return</span> <span class=\"n\">get_font</span><span
    class=\"p\">(</span><span class=\"n\">path</span><span class=\"p\">,</span> <span
    class=\"n\">draw</span><span class=\"p\">,</span> <span class=\"n\">title</span><span
    class=\"p\">,</span> <span class=\"n\">size</span> <span class=\"o\">-</span>
    <span class=\"mi\">10</span><span class=\"p\">,</span> <span class=\"n\">max_size</span><span
    class=\"o\">=</span><span class=\"n\">max_size</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">return</span> <span class=\"n\">font</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='PaddingError' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>PaddingError <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">PaddingError
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
    <span class=\"nc\">PaddingError</span><span class=\"p\">(</span><span class=\"ne\">BaseException</span><span
    class=\"p\">):</span>\n            <span class=\"k\">def</span> <span class=\"fm\">__init__</span><span
    class=\"p\">(</span>\n                <span class=\"bp\">self</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">msg</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;&quot;</span><span class=\"p\">,</span>\n
    \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"nb\">super</span><span class=\"p\">()</span><span
    class=\"o\">.</span><span class=\"fm\">__init__</span><span class=\"p\">(</span>\n
    \                   <span class=\"s2\">&quot;Padding must be an iterable of length
    1, 2, 3, or 4.</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span>
    <span class=\"o\">+</span> <span class=\"n\">msg</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! function
    <h2 id='draw_text' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>draw_text
    <em class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">draw_text <em class='small'>source</em></p>\n</div>\n<pre
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
    class=\"p\">,</span> <span class=\"o\">...</span><span class=\"p\">],</span>\n
    \           <span class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span>\n        <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"n\">text</span>
    <span class=\"o\">=</span> <span class=\"n\">text</span> <span class=\"ow\">or</span>
    <span class=\"s2\">&quot;&quot;</span>\n            <span class=\"n\">draw</span>
    <span class=\"o\">=</span> <span class=\"n\">ImageDraw</span><span class=\"o\">.</span><span
    class=\"n\">Draw</span><span class=\"p\">(</span><span class=\"n\">image</span><span
    class=\"p\">)</span>\n            <span class=\"n\">padding</span> <span class=\"o\">=</span>
    <span class=\"n\">resolve_padding</span><span class=\"p\">(</span><span class=\"n\">padding</span><span
    class=\"p\">,</span> <span class=\"n\">markata</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">width</span> <span class=\"o\">=</span> <span class=\"n\">image</span><span
    class=\"o\">.</span><span class=\"n\">size</span><span class=\"p\">[</span><span
    class=\"mi\">0</span><span class=\"p\">]</span>\n            <span class=\"n\">height</span>
    <span class=\"o\">=</span> <span class=\"n\">image</span><span class=\"o\">.</span><span
    class=\"n\">size</span><span class=\"p\">[</span><span class=\"mi\">1</span><span
    class=\"p\">]</span>\n            <span class=\"n\">bounding_box</span> <span
    class=\"o\">=</span> <span class=\"p\">[</span><span class=\"n\">padding</span><span
    class=\"p\">[</span><span class=\"mi\">0</span><span class=\"p\">],</span> <span
    class=\"n\">padding</span><span class=\"p\">[</span><span class=\"mi\">1</span><span
    class=\"p\">],</span> <span class=\"n\">width</span> <span class=\"o\">-</span>
    <span class=\"n\">padding</span><span class=\"p\">[</span><span class=\"mi\">0</span><span
    class=\"p\">],</span> <span class=\"n\">height</span> <span class=\"o\">-</span>
    <span class=\"n\">padding</span><span class=\"p\">[</span><span class=\"mi\">1</span><span
    class=\"p\">]]</span>\n            <span class=\"n\">bounding_box</span> <span
    class=\"o\">=</span> <span class=\"p\">[</span><span class=\"n\">padding</span><span
    class=\"p\">[</span><span class=\"mi\">0</span><span class=\"p\">],</span> <span
    class=\"n\">padding</span><span class=\"p\">[</span><span class=\"mi\">1</span><span
    class=\"p\">],</span> <span class=\"n\">width</span> <span class=\"o\">-</span>
    <span class=\"n\">padding</span><span class=\"p\">[</span><span class=\"mi\">2</span><span
    class=\"p\">],</span> <span class=\"n\">height</span> <span class=\"o\">-</span>
    <span class=\"n\">padding</span><span class=\"p\">[</span><span class=\"mi\">3</span><span
    class=\"p\">]]</span>\n            <span class=\"n\">max_size</span> <span class=\"o\">=</span>
    <span class=\"p\">(</span><span class=\"n\">bounding_box</span><span class=\"p\">[</span><span
    class=\"mi\">2</span><span class=\"p\">]</span> <span class=\"o\">-</span> <span
    class=\"n\">bounding_box</span><span class=\"p\">[</span><span class=\"mi\">0</span><span
    class=\"p\">],</span> <span class=\"n\">bounding_box</span><span class=\"p\">[</span><span
    class=\"mi\">3</span><span class=\"p\">]</span> <span class=\"o\">-</span> <span
    class=\"n\">bounding_box</span><span class=\"p\">[</span><span class=\"mi\">1</span><span
    class=\"p\">])</span>\n            <span class=\"n\">x1</span><span class=\"p\">,</span>
    <span class=\"n\">y1</span><span class=\"p\">,</span> <span class=\"n\">x2</span><span
    class=\"p\">,</span> <span class=\"n\">y2</span> <span class=\"o\">=</span> <span
    class=\"n\">bounding_box</span>\n            <span class=\"n\">font</span> <span
    class=\"o\">=</span> <span class=\"n\">get_font</span><span class=\"p\">(</span><span
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
    class=\"s2\">&quot;center&quot;</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='resolve_padding' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>resolve_padding <em class='small'>function</em></h2>\nConvert padding to
    a len 4 tuple</p>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">resolve_padding <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">resolve_padding</span><span class=\"p\">(</span><span class=\"n\">padding</span><span
    class=\"p\">:</span> <span class=\"n\">Tuple</span><span class=\"p\">[</span><span
    class=\"nb\">int</span><span class=\"p\">,</span> <span class=\"o\">...</span><span
    class=\"p\">],</span> <span class=\"n\">markata</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Tuple</span><span class=\"p\">[</span><span
    class=\"nb\">int</span><span class=\"p\">,</span> <span class=\"o\">...</span><span
    class=\"p\">]:</span>\n<span class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;Convert
    padding to a len 4 tuple&quot;&quot;&quot;</span>\n            <span class=\"k\">if</span>
    <span class=\"nb\">len</span><span class=\"p\">(</span><span class=\"n\">padding</span><span
    class=\"p\">)</span> <span class=\"o\">==</span> <span class=\"mi\">4</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"n\">padding</span>\n
    \           <span class=\"k\">if</span> <span class=\"nb\">len</span><span class=\"p\">(</span><span
    class=\"n\">padding</span><span class=\"p\">)</span> <span class=\"o\">==</span>
    <span class=\"mi\">3</span><span class=\"p\">:</span>\n                <span class=\"k\">return</span>
    <span class=\"p\">(</span><span class=\"o\">*</span><span class=\"n\">padding</span><span
    class=\"p\">,</span> <span class=\"n\">padding</span><span class=\"p\">[</span><span
    class=\"mi\">1</span><span class=\"p\">])</span>\n            <span class=\"k\">if</span>
    <span class=\"nb\">len</span><span class=\"p\">(</span><span class=\"n\">padding</span><span
    class=\"p\">)</span> <span class=\"o\">==</span> <span class=\"mi\">2</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"n\">padding</span>
    <span class=\"o\">*</span> <span class=\"mi\">2</span>\n            <span class=\"k\">if</span>
    <span class=\"nb\">len</span><span class=\"p\">(</span><span class=\"n\">padding</span><span
    class=\"p\">)</span> <span class=\"o\">==</span> <span class=\"mi\">1</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"n\">padding</span>
    <span class=\"o\">*</span> <span class=\"mi\">4</span>\n            <span class=\"k\">raise</span>
    <span class=\"n\">PaddingError</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;recieved padding: </span><span class=\"si\">{</span><span class=\"n\">padding</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='make_cover' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>make_cover <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">make_cover
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
    <span class=\"nf\">make_cover</span><span class=\"p\">(</span>\n            <span
    class=\"n\">title</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">,</span>\n            <span class=\"n\">color</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span><span class=\"p\">,</span>\n            <span class=\"n\">output_path</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">,</span>\n
    \           <span class=\"n\">template_path</span><span class=\"p\">:</span> <span
    class=\"n\">Path</span><span class=\"p\">,</span>\n            <span class=\"n\">font_path</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"n\">Path</span><span class=\"p\">],</span>\n            <span class=\"n\">padding</span><span
    class=\"p\">:</span> <span class=\"n\">Tuple</span><span class=\"p\">[</span><span
    class=\"nb\">int</span><span class=\"p\">,</span> <span class=\"o\">...</span><span
    class=\"p\">],</span>\n            <span class=\"n\">text_font</span><span class=\"p\">:</span>
    <span class=\"n\">Path</span><span class=\"p\">,</span>\n            <span class=\"n\">text</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span> <span class=\"o\">=</span>
    <span class=\"kc\">None</span><span class=\"p\">,</span>\n            <span class=\"n\">text_font_color</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span> <span class=\"o\">=</span>
    <span class=\"kc\">None</span><span class=\"p\">,</span>\n            <span class=\"n\">text_padding</span><span
    class=\"p\">:</span> <span class=\"n\">Tuple</span><span class=\"p\">[</span><span
    class=\"nb\">int</span><span class=\"p\">,</span> <span class=\"o\">...</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span><span
    class=\"p\">,</span>\n            <span class=\"n\">resizes</span><span class=\"p\">:</span>
    <span class=\"n\">List</span><span class=\"p\">[</span><span class=\"nb\">int</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span><span
    class=\"p\">,</span>\n            <span class=\"n\">markata</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;Markata&quot;</span> <span class=\"o\">=</span> <span
    class=\"kc\">None</span><span class=\"p\">,</span>\n        <span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">if</span> <span class=\"n\">output_path</span><span
    class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">():</span>\n
    \               <span class=\"k\">return</span>\n            <span class=\"n\">image</span>
    <span class=\"o\">=</span> <span class=\"n\">Image</span><span class=\"o\">.</span><span
    class=\"n\">open</span><span class=\"p\">(</span><span class=\"n\">template_path</span><span
    class=\"p\">)</span> <span class=\"k\">if</span> <span class=\"n\">template_path</span>
    <span class=\"k\">else</span> <span class=\"n\">Image</span><span class=\"o\">.</span><span
    class=\"n\">new</span><span class=\"p\">(</span><span class=\"s2\">&quot;RGB&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">(</span><span class=\"mi\">800</span><span
    class=\"p\">,</span> <span class=\"mi\">450</span><span class=\"p\">))</span>\n\n
    \           <span class=\"n\">draw_text</span><span class=\"p\">(</span>\n                <span
    class=\"n\">image</span><span class=\"o\">=</span><span class=\"n\">image</span><span
    class=\"p\">,</span>\n                <span class=\"n\">font_path</span><span
    class=\"o\">=</span><span class=\"n\">font_path</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">title</span><span class=\"o\">=</span><span
    class=\"n\">title</span><span class=\"p\">,</span>\n                <span class=\"n\">color</span><span
    class=\"o\">=</span><span class=\"n\">color</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">padding</span><span class=\"o\">=</span><span
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
    class=\"mi\">5</span><span class=\"p\">,</span>\n                        <span
    class=\"n\">image</span><span class=\"o\">.</span><span class=\"n\">size</span><span
    class=\"p\">[</span><span class=\"mi\">0</span><span class=\"p\">]</span> <span
    class=\"o\">/</span> <span class=\"mi\">5</span><span class=\"p\">,</span>\n                        <span
    class=\"n\">image</span><span class=\"o\">.</span><span class=\"n\">size</span><span
    class=\"p\">[</span><span class=\"mi\">1</span><span class=\"p\">]</span> <span
    class=\"o\">-</span> <span class=\"n\">image</span><span class=\"o\">.</span><span
    class=\"n\">size</span><span class=\"p\">[</span><span class=\"mi\">1</span><span
    class=\"p\">]</span> <span class=\"o\">/</span> <span class=\"mi\">10</span><span
    class=\"p\">,</span>\n                    <span class=\"p\">)</span>\n                <span
    class=\"n\">draw_text</span><span class=\"p\">(</span><span class=\"n\">image</span><span
    class=\"p\">,</span> <span class=\"n\">text_font</span><span class=\"p\">,</span>
    <span class=\"n\">text</span><span class=\"p\">,</span> <span class=\"n\">text_font_color</span><span
    class=\"p\">,</span> <span class=\"n\">text_padding</span><span class=\"p\">)</span>\n\n
    \           <span class=\"n\">image</span><span class=\"o\">.</span><span class=\"n\">save</span><span
    class=\"p\">(</span><span class=\"n\">output_path</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">ratio</span> <span class=\"o\">=</span> <span class=\"n\">image</span><span
    class=\"o\">.</span><span class=\"n\">size</span><span class=\"p\">[</span><span
    class=\"mi\">1</span><span class=\"p\">]</span> <span class=\"o\">/</span> <span
    class=\"n\">image</span><span class=\"o\">.</span><span class=\"n\">size</span><span
    class=\"p\">[</span><span class=\"mi\">0</span><span class=\"p\">]</span>\n\n
    \           <span class=\"n\">covers</span> <span class=\"o\">=</span> <span class=\"p\">[]</span>\n
    \           <span class=\"k\">if</span> <span class=\"n\">resizes</span><span
    class=\"p\">:</span>\n                <span class=\"k\">for</span> <span class=\"n\">width</span>
    <span class=\"ow\">in</span> <span class=\"n\">resizes</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">re_img</span> <span class=\"o\">=</span>
    <span class=\"n\">image</span><span class=\"o\">.</span><span class=\"n\">resize</span><span
    class=\"p\">((</span><span class=\"n\">width</span><span class=\"p\">,</span>
    <span class=\"nb\">int</span><span class=\"p\">(</span><span class=\"n\">width</span>
    <span class=\"o\">*</span> <span class=\"n\">ratio</span><span class=\"p\">)),</span>
    <span class=\"n\">Image</span><span class=\"o\">.</span><span class=\"n\">ANTIALIAS</span><span
    class=\"p\">)</span>\n                    <span class=\"n\">filename</span> <span
    class=\"o\">=</span> <span class=\"p\">(</span>\n                        <span
    class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">output_path</span><span class=\"o\">.</span><span class=\"n\">stem</span><span
    class=\"si\">}</span><span class=\"s2\">_</span><span class=\"si\">{</span><span
    class=\"n\">width</span><span class=\"si\">}</span><span class=\"s2\">x</span><span
    class=\"si\">{</span><span class=\"nb\">int</span><span class=\"p\">(</span><span
    class=\"n\">width</span><span class=\"o\">*</span><span class=\"n\">ratio</span><span
    class=\"p\">)</span><span class=\"si\">}{</span><span class=\"n\">output_path</span><span
    class=\"o\">.</span><span class=\"n\">suffix</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span>\n                    <span class=\"p\">)</span>\n                    <span
    class=\"n\">covers</span><span class=\"o\">.</span><span class=\"n\">append</span><span
    class=\"p\">(</span><span class=\"n\">filename</span><span class=\"p\">)</span>\n\n
    \                   <span class=\"n\">filepath</span> <span class=\"o\">=</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">output_path</span><span
    class=\"o\">.</span><span class=\"n\">parent</span> <span class=\"o\">/</span>
    <span class=\"n\">filename</span><span class=\"p\">)</span>\n                    <span
    class=\"n\">re_img</span><span class=\"o\">.</span><span class=\"n\">save</span><span
    class=\"p\">(</span><span class=\"n\">filepath</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='save' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>save <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">save
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
    <span class=\"nf\">save</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">futures</span> <span class=\"o\">=</span> <span
    class=\"p\">[]</span>\n\n            <span class=\"k\">if</span> <span class=\"s2\">&quot;covers&quot;</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">keys</span><span class=\"p\">():</span>\n                <span class=\"k\">return</span>\n\n
    \           <span class=\"k\">for</span> <span class=\"n\">article</span> <span
    class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">iter_articles</span><span class=\"p\">(</span><span class=\"s2\">&quot;making
    covers&quot;</span><span class=\"p\">):</span>\n                <span class=\"k\">for</span>
    <span class=\"n\">cover</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;covers&quot;</span><span class=\"p\">]:</span>\n                    <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">padding</span> <span class=\"o\">=</span> <span class=\"n\">cover</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;padding&quot;</span><span class=\"p\">]</span>\n
    \                   <span class=\"k\">except</span> <span class=\"ne\">KeyError</span><span
    class=\"p\">:</span>\n                        <span class=\"n\">padding</span>
    <span class=\"o\">=</span> <span class=\"p\">(</span>\n                            <span
    class=\"mi\">200</span><span class=\"p\">,</span>\n                            <span
    class=\"mi\">100</span><span class=\"p\">,</span>\n                        <span
    class=\"p\">)</span>\n                    <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">text_padding</span> <span class=\"o\">=</span>
    <span class=\"n\">cover</span><span class=\"p\">[</span><span class=\"s2\">&quot;text_padding&quot;</span><span
    class=\"p\">]</span>\n                    <span class=\"k\">except</span> <span
    class=\"ne\">KeyError</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">text_padding</span> <span class=\"o\">=</span> <span class=\"p\">(</span>\n
    \                           <span class=\"mi\">200</span><span class=\"p\">,</span>\n
    \                           <span class=\"mi\">100</span><span class=\"p\">,</span>\n
    \                       <span class=\"p\">)</span>\n                    <span
    class=\"k\">if</span> <span class=\"s2\">&quot;text_key&quot;</span> <span class=\"ow\">in</span>
    <span class=\"n\">cover</span><span class=\"p\">:</span>\n                        <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                            <span
    class=\"n\">text</span> <span class=\"o\">=</span> <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">metadata</span><span class=\"p\">[</span><span
    class=\"n\">cover</span><span class=\"p\">[</span><span class=\"s2\">&quot;text_key&quot;</span><span
    class=\"p\">]]</span>\n                        <span class=\"k\">except</span>
    <span class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n                            <span
    class=\"n\">text</span> <span class=\"o\">=</span> <span class=\"n\">article</span><span
    class=\"p\">[</span><span class=\"n\">cover</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;text_key&quot;</span><span class=\"p\">]]</span>\n                        <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                            <span
    class=\"n\">text</span> <span class=\"o\">=</span> <span class=\"n\">text</span><span
    class=\"o\">.</span><span class=\"n\">replace</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;&quot;</span><span class=\"p\">)</span>\n
    \                           <span class=\"kn\">from</span> <span class=\"nn\">more_itertools</span>
    <span class=\"kn\">import</span> <span class=\"n\">chunked</span>\n\n                            <span
    class=\"n\">text</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;</span><span
    class=\"se\">\\n</span><span class=\"s2\">&quot;</span><span class=\"o\">.</span><span
    class=\"n\">join</span><span class=\"p\">([</span><span class=\"s2\">&quot;&quot;</span><span
    class=\"o\">.</span><span class=\"n\">join</span><span class=\"p\">(</span><span
    class=\"n\">c</span><span class=\"p\">)</span> <span class=\"k\">for</span> <span
    class=\"n\">c</span> <span class=\"ow\">in</span> <span class=\"n\">chunked</span><span
    class=\"p\">(</span><span class=\"n\">text</span><span class=\"p\">,</span> <span
    class=\"mi\">60</span><span class=\"p\">)])</span>\n                        <span
    class=\"k\">except</span> <span class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n
    \                           <span class=\"c1\"># text is likely None</span>\n
    \                           <span class=\"k\">pass</span>\n\n                        <span
    class=\"n\">text_font</span> <span class=\"o\">=</span> <span class=\"n\">cover</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;text_font&quot;</span><span class=\"p\">]</span>\n
    \                       <span class=\"n\">text_font_color</span> <span class=\"o\">=</span>
    <span class=\"n\">cover</span><span class=\"p\">[</span><span class=\"s2\">&quot;text_font_color&quot;</span><span
    class=\"p\">]</span>\n                    <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">text</span> <span class=\"o\">=</span>
    <span class=\"kc\">None</span>\n                        <span class=\"n\">text_font</span>
    <span class=\"o\">=</span> <span class=\"kc\">None</span>\n                        <span
    class=\"n\">text_font_color</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n
    \                   <span class=\"k\">try</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">title</span> <span class=\"o\">=</span> <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">metadata</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;title&quot;</span><span class=\"p\">]</span>\n                    <span
    class=\"k\">except</span> <span class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">title</span> <span class=\"o\">=</span>
    <span class=\"n\">article</span><span class=\"p\">[</span><span class=\"s2\">&quot;title&quot;</span><span
    class=\"p\">]</span>\n                    <span class=\"n\">futures</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span>\n
    \                       <span class=\"n\">make_cover</span><span class=\"p\">(</span>\n
    \                           <span class=\"n\">title</span><span class=\"o\">=</span><span
    class=\"n\">title</span><span class=\"p\">,</span>\n                            <span
    class=\"n\">color</span><span class=\"o\">=</span><span class=\"n\">cover</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;font_color&quot;</span><span class=\"p\">],</span>\n
    \                           <span class=\"n\">output_path</span><span class=\"o\">=</span><span
    class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">output_dir</span><span class=\"p\">)</span>\n                            <span
    class=\"o\">/</span> <span class=\"p\">(</span><span class=\"n\">article</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">+</span> <span class=\"n\">cover</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;name&quot;</span><span class=\"p\">]</span> <span class=\"o\">+</span>
    <span class=\"s2\">&quot;.png&quot;</span><span class=\"p\">),</span>\n                            <span
    class=\"n\">template_path</span><span class=\"o\">=</span><span class=\"n\">cover</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;template&quot;</span><span class=\"p\">,</span> <span class=\"kc\">None</span><span
    class=\"p\">),</span>\n                            <span class=\"n\">font_path</span><span
    class=\"o\">=</span><span class=\"n\">cover</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;font&quot;</span><span
    class=\"p\">,</span> <span class=\"kc\">None</span><span class=\"p\">),</span>\n
    \                           <span class=\"n\">padding</span><span class=\"o\">=</span><span
    class=\"n\">padding</span><span class=\"p\">,</span>\n                            <span
    class=\"n\">text_font</span><span class=\"o\">=</span><span class=\"n\">text_font</span><span
    class=\"p\">,</span>\n                            <span class=\"n\">text</span><span
    class=\"o\">=</span><span class=\"n\">text</span><span class=\"p\">,</span>\n
    \                           <span class=\"n\">text_font_color</span><span class=\"o\">=</span><span
    class=\"n\">text_font_color</span><span class=\"p\">,</span>\n                            <span
    class=\"n\">text_padding</span><span class=\"o\">=</span><span class=\"n\">text_padding</span><span
    class=\"p\">,</span>\n                            <span class=\"n\">resizes</span><span
    class=\"o\">=</span><span class=\"n\">cover</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;resizes&quot;</span><span
    class=\"p\">),</span>\n                            <span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"p\">,</span>\n
    \                       <span class=\"p\">),</span>\n                    <span
    class=\"p\">)</span>\n\n            <span class=\"n\">progress</span> <span class=\"o\">=</span>
    <span class=\"n\">Progress</span><span class=\"p\">(</span>\n                <span
    class=\"n\">BarColumn</span><span class=\"p\">(</span><span class=\"n\">bar_width</span><span
    class=\"o\">=</span><span class=\"kc\">None</span><span class=\"p\">),</span>\n
    \               <span class=\"n\">transient</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n                <span class=\"n\">console</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"p\">,</span>\n            <span class=\"p\">)</span>\n
    \           <span class=\"n\">task_id</span> <span class=\"o\">=</span> <span
    class=\"n\">progress</span><span class=\"o\">.</span><span class=\"n\">add_task</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;loading markdown&quot;</span><span
    class=\"p\">)</span>\n            <span class=\"n\">progress</span><span class=\"o\">.</span><span
    class=\"n\">update</span><span class=\"p\">(</span><span class=\"n\">task_id</span><span
    class=\"p\">,</span> <span class=\"n\">total</span><span class=\"o\">=</span><span
    class=\"nb\">len</span><span class=\"p\">(</span><span class=\"n\">futures</span><span
    class=\"p\">))</span>\n            <span class=\"k\">with</span> <span class=\"n\">progress</span><span
    class=\"p\">:</span>\n                <span class=\"k\">while</span> <span class=\"ow\">not</span>
    <span class=\"nb\">all</span><span class=\"p\">(</span><span class=\"n\">f</span><span
    class=\"o\">.</span><span class=\"n\">done</span><span class=\"p\">()</span> <span
    class=\"k\">for</span> <span class=\"n\">f</span> <span class=\"ow\">in</span>
    <span class=\"n\">futures</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">time</span><span class=\"o\">.</span><span class=\"n\">sleep</span><span
    class=\"p\">(</span><span class=\"mf\">0.1</span><span class=\"p\">)</span>\n
    \                   <span class=\"n\">progress</span><span class=\"o\">.</span><span
    class=\"n\">update</span><span class=\"p\">(</span><span class=\"n\">task_id</span><span
    class=\"p\">,</span> <span class=\"n\">total</span><span class=\"o\">=</span><span
    class=\"nb\">len</span><span class=\"p\">([</span><span class=\"n\">f</span> <span
    class=\"k\">for</span> <span class=\"n\">f</span> <span class=\"ow\">in</span>
    <span class=\"n\">futures</span> <span class=\"k\">if</span> <span class=\"n\">f</span><span
    class=\"o\">.</span><span class=\"n\">done</span><span class=\"p\">()]))</span>\n
    \           <span class=\"p\">[</span><span class=\"n\">f</span><span class=\"o\">.</span><span
    class=\"n\">result</span><span class=\"p\">()</span> <span class=\"k\">for</span>
    <span class=\"n\">f</span> <span class=\"ow\">in</span> <span class=\"n\">futures</span><span
    class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='__init__'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'><strong>init</strong>
    <em class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\"><strong>init</strong> <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"fm\">__init__</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">self</span><span class=\"p\">,</span>\n                <span class=\"n\">msg</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;&quot;</span><span class=\"p\">,</span>\n            <span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"nb\">super</span><span class=\"p\">()</span><span
    class=\"o\">.</span><span class=\"fm\">__init__</span><span class=\"p\">(</span>\n
    \                   <span class=\"s2\">&quot;Padding must be an iterable of length
    1, 2, 3, or 4.</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span>
    <span class=\"o\">+</span> <span class=\"n\">msg</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>"
  raw.md: ''
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

