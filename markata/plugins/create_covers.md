---
content: "None\n\n\n!! function <h2 id='get_font' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>get_font <em class='small'>function</em></h2>\n\n???+ source \"get_font <em
  class='small'>source</em>\"\n\n```python\n\n        def get_font(\n            path:
  Path,\n            draw: ImageDraw.Draw,\n            title: str,\n            size:
  int = 250,\n        ) -> ImageFont.FreeTypeFont:\n            font = ImageFont.truetype(path,
  size=size)\n            if draw.textsize(title, font=font)[0] > 800:\n                return
  get_font(path, draw, title, size - 10)\n            return font\n```\n\n\n!! function
  <h2 id='make_cover' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>make_cover
  <em class='small'>function</em></h2>\n\n???+ source \"make_cover <em class='small'>source</em>\"\n\n```python\n\n
  \       def make_cover(\n            title: str,\n            color: str,\n            output_path:
  Path,\n            template_path: Path,\n            font_path: Path,\n        )
  -> None:\n            image = Image.open(template_path)\n\n            draw = ImageDraw.Draw(image)\n\n
  \           font = get_font(font_path, draw, title)\n\n            color = \"rgb(255,255,255)\"\n
  \           padding = (200, 100)\n            bounding_box = [padding[0], padding[1],
  1000 - padding[0], 420 - padding[1]]\n            x1, y1, x2, y2 = bounding_box\n
  \           w, h = draw.textsize(title, font=font)\n            x = (x2 - x1 - w)
  / 2 + x1\n            y = (y2 - y1 - h) / 2 + y1\n            draw.text((x, y),
  title, fill=color, font=font, align=\"center\")\n            image.save(output_path)\n```\n\n\n!!
  function <h2 id='save' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>save
  <em class='small'>function</em></h2>\n\n???+ source \"save <em class='small'>source</em>\"\n\n```python\n\n
  \       def save(markata: \"Markata\") -> None:\n            for article in markata.articles:\n
  \               output_path = Path(markata.output_dir) / (\n                    Path(article.metadata[\"path\"]).stem
  + \".png\"\n                )\n\n                make_cover(\n                    article.metadata[\"title\"],\n
  \                   markata.config[\"cover_font_color\"],\n                    output_path,\n
  \                   markata.config[\"cover_template\"],\n                    markata.config[\"cover_font\"],\n
  \               )\n```\n\n"
date: 0001-01-01
description: 'None ! ???+ source  ! ???+ source  ! ???+ source '
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Create_Covers.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None ! ???+ source  ! ???+ source  !
    ???+ source \" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Create_Covers.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None ! ???+ source  ! ???+ source  !
    ???+ source \" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    </head>\n    <body>\n<div class='container
    flex flex-row min-h-screen'>\n    <div>\n    </div>\n    <div class='flex-grow
    px-8 mx-auto min-h-screen'>\n<header class='flex justify-center items-center p-8'>\n\n
    \   <nav class='flex justify-center items-center my-8'>\n        <a\n            href='/'>markata</a>\n
    \       <a\n            href='https://github.com/WaylonWalker/markata'>GitHub</a>\n
    \       <a\n            href='https://markata.dev/docs/'>docs</a>\n        <a\n
    \           href='https://markata.dev/plugins/'>plugins</a>\n    </nav>\n\n    <div>\n
    \       <label id=\"theme-switch\" class=\"theme-switch\" for=\"checkbox-theme\"
    title=\"light/dark mode toggle\">\n            <input type=\"checkbox\" id=\"checkbox-theme\"
    />\n            <div class=\"slider round\"></div>\n        </label>\n    </div>\n</header><article
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Create_Covers.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <p>None</p>\n<p>!!
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
    <span class=\"mi\">250</span><span class=\"p\">,</span>\n        <span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">ImageFont</span><span class=\"o\">.</span><span
    class=\"n\">FreeTypeFont</span><span class=\"p\">:</span>\n            <span class=\"n\">font</span>
    <span class=\"o\">=</span> <span class=\"n\">ImageFont</span><span class=\"o\">.</span><span
    class=\"n\">truetype</span><span class=\"p\">(</span><span class=\"n\">path</span><span
    class=\"p\">,</span> <span class=\"n\">size</span><span class=\"o\">=</span><span
    class=\"n\">size</span><span class=\"p\">)</span>\n            <span class=\"k\">if</span>
    <span class=\"n\">draw</span><span class=\"o\">.</span><span class=\"n\">textsize</span><span
    class=\"p\">(</span><span class=\"n\">title</span><span class=\"p\">,</span> <span
    class=\"n\">font</span><span class=\"o\">=</span><span class=\"n\">font</span><span
    class=\"p\">)[</span><span class=\"mi\">0</span><span class=\"p\">]</span> <span
    class=\"o\">&gt;</span> <span class=\"mi\">800</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">get_font</span><span
    class=\"p\">(</span><span class=\"n\">path</span><span class=\"p\">,</span> <span
    class=\"n\">draw</span><span class=\"p\">,</span> <span class=\"n\">title</span><span
    class=\"p\">,</span> <span class=\"n\">size</span> <span class=\"o\">-</span>
    <span class=\"mi\">10</span><span class=\"p\">)</span>\n            <span class=\"k\">return</span>
    <span class=\"n\">font</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='make_cover'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>make_cover <em class='small'>function</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">make_cover
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
    class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">,</span>\n
    \       <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n            <span class=\"n\">image</span> <span class=\"o\">=</span>
    <span class=\"n\">Image</span><span class=\"o\">.</span><span class=\"n\">open</span><span
    class=\"p\">(</span><span class=\"n\">template_path</span><span class=\"p\">)</span>\n\n
    \           <span class=\"n\">draw</span> <span class=\"o\">=</span> <span class=\"n\">ImageDraw</span><span
    class=\"o\">.</span><span class=\"n\">Draw</span><span class=\"p\">(</span><span
    class=\"n\">image</span><span class=\"p\">)</span>\n\n            <span class=\"n\">font</span>
    <span class=\"o\">=</span> <span class=\"n\">get_font</span><span class=\"p\">(</span><span
    class=\"n\">font_path</span><span class=\"p\">,</span> <span class=\"n\">draw</span><span
    class=\"p\">,</span> <span class=\"n\">title</span><span class=\"p\">)</span>\n\n
    \           <span class=\"n\">color</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;rgb(255,255,255)&quot;</span>\n
    \           <span class=\"n\">padding</span> <span class=\"o\">=</span> <span
    class=\"p\">(</span><span class=\"mi\">200</span><span class=\"p\">,</span> <span
    class=\"mi\">100</span><span class=\"p\">)</span>\n            <span class=\"n\">bounding_box</span>
    <span class=\"o\">=</span> <span class=\"p\">[</span><span class=\"n\">padding</span><span
    class=\"p\">[</span><span class=\"mi\">0</span><span class=\"p\">],</span> <span
    class=\"n\">padding</span><span class=\"p\">[</span><span class=\"mi\">1</span><span
    class=\"p\">],</span> <span class=\"mi\">1000</span> <span class=\"o\">-</span>
    <span class=\"n\">padding</span><span class=\"p\">[</span><span class=\"mi\">0</span><span
    class=\"p\">],</span> <span class=\"mi\">420</span> <span class=\"o\">-</span>
    <span class=\"n\">padding</span><span class=\"p\">[</span><span class=\"mi\">1</span><span
    class=\"p\">]]</span>\n            <span class=\"n\">x1</span><span class=\"p\">,</span>
    <span class=\"n\">y1</span><span class=\"p\">,</span> <span class=\"n\">x2</span><span
    class=\"p\">,</span> <span class=\"n\">y2</span> <span class=\"o\">=</span> <span
    class=\"n\">bounding_box</span>\n            <span class=\"n\">w</span><span class=\"p\">,</span>
    <span class=\"n\">h</span> <span class=\"o\">=</span> <span class=\"n\">draw</span><span
    class=\"o\">.</span><span class=\"n\">textsize</span><span class=\"p\">(</span><span
    class=\"n\">title</span><span class=\"p\">,</span> <span class=\"n\">font</span><span
    class=\"o\">=</span><span class=\"n\">font</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">x</span> <span class=\"o\">=</span> <span class=\"p\">(</span><span
    class=\"n\">x2</span> <span class=\"o\">-</span> <span class=\"n\">x1</span> <span
    class=\"o\">-</span> <span class=\"n\">w</span><span class=\"p\">)</span> <span
    class=\"o\">/</span> <span class=\"mi\">2</span> <span class=\"o\">+</span> <span
    class=\"n\">x1</span>\n            <span class=\"n\">y</span> <span class=\"o\">=</span>
    <span class=\"p\">(</span><span class=\"n\">y2</span> <span class=\"o\">-</span>
    <span class=\"n\">y1</span> <span class=\"o\">-</span> <span class=\"n\">h</span><span
    class=\"p\">)</span> <span class=\"o\">/</span> <span class=\"mi\">2</span> <span
    class=\"o\">+</span> <span class=\"n\">y1</span>\n            <span class=\"n\">draw</span><span
    class=\"o\">.</span><span class=\"n\">text</span><span class=\"p\">((</span><span
    class=\"n\">x</span><span class=\"p\">,</span> <span class=\"n\">y</span><span
    class=\"p\">),</span> <span class=\"n\">title</span><span class=\"p\">,</span>
    <span class=\"n\">fill</span><span class=\"o\">=</span><span class=\"n\">color</span><span
    class=\"p\">,</span> <span class=\"n\">font</span><span class=\"o\">=</span><span
    class=\"n\">font</span><span class=\"p\">,</span> <span class=\"n\">align</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;center&quot;</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">image</span><span class=\"o\">.</span><span class=\"n\">save</span><span
    class=\"p\">(</span><span class=\"n\">output_path</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    \           <span class=\"k\">for</span> <span class=\"n\">article</span> <span
    class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">articles</span><span class=\"p\">:</span>\n                <span class=\"n\">output_path</span>
    <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">output_dir</span><span
    class=\"p\">)</span> <span class=\"o\">/</span> <span class=\"p\">(</span>\n                    <span
    class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">metadata</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;path&quot;</span><span class=\"p\">])</span><span class=\"o\">.</span><span
    class=\"n\">stem</span> <span class=\"o\">+</span> <span class=\"s2\">&quot;.png&quot;</span>\n
    \               <span class=\"p\">)</span>\n\n                <span class=\"n\">make_cover</span><span
    class=\"p\">(</span>\n                    <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">metadata</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;title&quot;</span><span class=\"p\">],</span>\n                    <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;cover_font_color&quot;</span><span
    class=\"p\">],</span>\n                    <span class=\"n\">output_path</span><span
    class=\"p\">,</span>\n                    <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;cover_template&quot;</span><span class=\"p\">],</span>\n                    <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;cover_font&quot;</span><span class=\"p\">],</span>\n
    \               <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>
    \   </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Create_Covers.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None ! ???+ source  ! ???+ source  !
    ???+ source \" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Create_Covers.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None ! ???+ source  ! ???+ source  !
    ???+ source \" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    </head>\n    <body>\n<article style=\"text-align:
    center;\">\n    <style>\n        section {\n            font-size: 200%;\n        }\n\n\n
    \       .edit {\n            display: none;\n        }\n    </style>\n<section
    class=\"title\">\n    <h1 id=\"title\">\n        Create_Covers.Py\n    </h1>\n</section></article>\n
    \    </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Create_Covers.Py\n    </h1>\n</section>    <section class=\"body\">\n
    \       <p>None</p>\n<p>!! function <h2 id='get_font' class='admonition-title'
    style='margin:0;padding:.5rem 1rem;'>get_font <em class='small'>function</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">get_font
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
    <span class=\"mi\">250</span><span class=\"p\">,</span>\n        <span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">ImageFont</span><span class=\"o\">.</span><span
    class=\"n\">FreeTypeFont</span><span class=\"p\">:</span>\n            <span class=\"n\">font</span>
    <span class=\"o\">=</span> <span class=\"n\">ImageFont</span><span class=\"o\">.</span><span
    class=\"n\">truetype</span><span class=\"p\">(</span><span class=\"n\">path</span><span
    class=\"p\">,</span> <span class=\"n\">size</span><span class=\"o\">=</span><span
    class=\"n\">size</span><span class=\"p\">)</span>\n            <span class=\"k\">if</span>
    <span class=\"n\">draw</span><span class=\"o\">.</span><span class=\"n\">textsize</span><span
    class=\"p\">(</span><span class=\"n\">title</span><span class=\"p\">,</span> <span
    class=\"n\">font</span><span class=\"o\">=</span><span class=\"n\">font</span><span
    class=\"p\">)[</span><span class=\"mi\">0</span><span class=\"p\">]</span> <span
    class=\"o\">&gt;</span> <span class=\"mi\">800</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">get_font</span><span
    class=\"p\">(</span><span class=\"n\">path</span><span class=\"p\">,</span> <span
    class=\"n\">draw</span><span class=\"p\">,</span> <span class=\"n\">title</span><span
    class=\"p\">,</span> <span class=\"n\">size</span> <span class=\"o\">-</span>
    <span class=\"mi\">10</span><span class=\"p\">)</span>\n            <span class=\"k\">return</span>
    <span class=\"n\">font</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='make_cover'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>make_cover <em class='small'>function</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">make_cover
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
    class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">,</span>\n
    \       <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n            <span class=\"n\">image</span> <span class=\"o\">=</span>
    <span class=\"n\">Image</span><span class=\"o\">.</span><span class=\"n\">open</span><span
    class=\"p\">(</span><span class=\"n\">template_path</span><span class=\"p\">)</span>\n\n
    \           <span class=\"n\">draw</span> <span class=\"o\">=</span> <span class=\"n\">ImageDraw</span><span
    class=\"o\">.</span><span class=\"n\">Draw</span><span class=\"p\">(</span><span
    class=\"n\">image</span><span class=\"p\">)</span>\n\n            <span class=\"n\">font</span>
    <span class=\"o\">=</span> <span class=\"n\">get_font</span><span class=\"p\">(</span><span
    class=\"n\">font_path</span><span class=\"p\">,</span> <span class=\"n\">draw</span><span
    class=\"p\">,</span> <span class=\"n\">title</span><span class=\"p\">)</span>\n\n
    \           <span class=\"n\">color</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;rgb(255,255,255)&quot;</span>\n
    \           <span class=\"n\">padding</span> <span class=\"o\">=</span> <span
    class=\"p\">(</span><span class=\"mi\">200</span><span class=\"p\">,</span> <span
    class=\"mi\">100</span><span class=\"p\">)</span>\n            <span class=\"n\">bounding_box</span>
    <span class=\"o\">=</span> <span class=\"p\">[</span><span class=\"n\">padding</span><span
    class=\"p\">[</span><span class=\"mi\">0</span><span class=\"p\">],</span> <span
    class=\"n\">padding</span><span class=\"p\">[</span><span class=\"mi\">1</span><span
    class=\"p\">],</span> <span class=\"mi\">1000</span> <span class=\"o\">-</span>
    <span class=\"n\">padding</span><span class=\"p\">[</span><span class=\"mi\">0</span><span
    class=\"p\">],</span> <span class=\"mi\">420</span> <span class=\"o\">-</span>
    <span class=\"n\">padding</span><span class=\"p\">[</span><span class=\"mi\">1</span><span
    class=\"p\">]]</span>\n            <span class=\"n\">x1</span><span class=\"p\">,</span>
    <span class=\"n\">y1</span><span class=\"p\">,</span> <span class=\"n\">x2</span><span
    class=\"p\">,</span> <span class=\"n\">y2</span> <span class=\"o\">=</span> <span
    class=\"n\">bounding_box</span>\n            <span class=\"n\">w</span><span class=\"p\">,</span>
    <span class=\"n\">h</span> <span class=\"o\">=</span> <span class=\"n\">draw</span><span
    class=\"o\">.</span><span class=\"n\">textsize</span><span class=\"p\">(</span><span
    class=\"n\">title</span><span class=\"p\">,</span> <span class=\"n\">font</span><span
    class=\"o\">=</span><span class=\"n\">font</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">x</span> <span class=\"o\">=</span> <span class=\"p\">(</span><span
    class=\"n\">x2</span> <span class=\"o\">-</span> <span class=\"n\">x1</span> <span
    class=\"o\">-</span> <span class=\"n\">w</span><span class=\"p\">)</span> <span
    class=\"o\">/</span> <span class=\"mi\">2</span> <span class=\"o\">+</span> <span
    class=\"n\">x1</span>\n            <span class=\"n\">y</span> <span class=\"o\">=</span>
    <span class=\"p\">(</span><span class=\"n\">y2</span> <span class=\"o\">-</span>
    <span class=\"n\">y1</span> <span class=\"o\">-</span> <span class=\"n\">h</span><span
    class=\"p\">)</span> <span class=\"o\">/</span> <span class=\"mi\">2</span> <span
    class=\"o\">+</span> <span class=\"n\">y1</span>\n            <span class=\"n\">draw</span><span
    class=\"o\">.</span><span class=\"n\">text</span><span class=\"p\">((</span><span
    class=\"n\">x</span><span class=\"p\">,</span> <span class=\"n\">y</span><span
    class=\"p\">),</span> <span class=\"n\">title</span><span class=\"p\">,</span>
    <span class=\"n\">fill</span><span class=\"o\">=</span><span class=\"n\">color</span><span
    class=\"p\">,</span> <span class=\"n\">font</span><span class=\"o\">=</span><span
    class=\"n\">font</span><span class=\"p\">,</span> <span class=\"n\">align</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;center&quot;</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">image</span><span class=\"o\">.</span><span class=\"n\">save</span><span
    class=\"p\">(</span><span class=\"n\">output_path</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    \           <span class=\"k\">for</span> <span class=\"n\">article</span> <span
    class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">articles</span><span class=\"p\">:</span>\n                <span class=\"n\">output_path</span>
    <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">output_dir</span><span
    class=\"p\">)</span> <span class=\"o\">/</span> <span class=\"p\">(</span>\n                    <span
    class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">metadata</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;path&quot;</span><span class=\"p\">])</span><span class=\"o\">.</span><span
    class=\"n\">stem</span> <span class=\"o\">+</span> <span class=\"s2\">&quot;.png&quot;</span>\n
    \               <span class=\"p\">)</span>\n\n                <span class=\"n\">make_cover</span><span
    class=\"p\">(</span>\n                    <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">metadata</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;title&quot;</span><span class=\"p\">],</span>\n                    <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;cover_font_color&quot;</span><span
    class=\"p\">],</span>\n                    <span class=\"n\">output_path</span><span
    class=\"p\">,</span>\n                    <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;cover_template&quot;</span><span class=\"p\">],</span>\n                    <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;cover_font&quot;</span><span class=\"p\">],</span>\n
    \               <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>"
  raw.md: ''
published: true
slug: markata/plugins/create-covers
title: Create_Covers.Py


---

None


!! function <h2 id='get_font' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get_font <em class='small'>function</em></h2>

???+ source "get_font <em class='small'>source</em>"

```python

        def get_font(
            path: Path,
            draw: ImageDraw.Draw,
            title: str,
            size: int = 250,
        ) -> ImageFont.FreeTypeFont:
            font = ImageFont.truetype(path, size=size)
            if draw.textsize(title, font=font)[0] > 800:
                return get_font(path, draw, title, size - 10)
            return font
```


!! function <h2 id='make_cover' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>make_cover <em class='small'>function</em></h2>

???+ source "make_cover <em class='small'>source</em>"

```python

        def make_cover(
            title: str,
            color: str,
            output_path: Path,
            template_path: Path,
            font_path: Path,
        ) -> None:
            image = Image.open(template_path)

            draw = ImageDraw.Draw(image)

            font = get_font(font_path, draw, title)

            color = "rgb(255,255,255)"
            padding = (200, 100)
            bounding_box = [padding[0], padding[1], 1000 - padding[0], 420 - padding[1]]
            x1, y1, x2, y2 = bounding_box
            w, h = draw.textsize(title, font=font)
            x = (x2 - x1 - w) / 2 + x1
            y = (y2 - y1 - h) / 2 + y1
            draw.text((x, y), title, fill=color, font=font, align="center")
            image.save(output_path)
```


!! function <h2 id='save' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>save <em class='small'>function</em></h2>

???+ source "save <em class='small'>source</em>"

```python

        def save(markata: "Markata") -> None:
            for article in markata.articles:
                output_path = Path(markata.output_dir) / (
                    Path(article.metadata["path"]).stem + ".png"
                )

                make_cover(
                    article.metadata["title"],
                    markata.config["cover_font_color"],
                    output_path,
                    markata.config["cover_template"],
                    markata.config["cover_font"],
                )
```

