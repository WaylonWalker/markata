---
content: "manifest plugin\n\n\n!! function <h2 id='render' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>render <em class='small'>function</em></h2>\n\n???+
  source \"render <em class='small'>source</em>\"\n\n```python\n\n        def render(markata:
  \"MarkataIcons\") -> None:\n            icons = markata.icons if \"icons\" in markata.__dict__
  else []\n            manifest = {\n                \"name\": markata.config.site_name,\n
  \               \"short_name\": markata.config.short_name,\n                \"start_url\":
  markata.config.start_url,\n                \"display\": markata.config.display,\n
  \               \"background_color\": str(markata.config.background_color),\n                \"theme_color\":
  str(markata.config.theme_color),\n                \"description\": markata.config.description,\n
  \               \"icons\": icons,\n            }\n            filepath = Path(markata.config[\"output_dir\"])
  / \"manifest.json\"\n            filepath.touch(exist_ok=True)\n            with
  open(filepath, \"w+\") as f:\n                json.dump(manifest, f, ensure_ascii=True,
  indent=4)\n```\n\n"
date: 0001-01-01
description: 'manifest plugin ! ???+ source '
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Manifest.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"manifest plugin ! ???+ source \" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Manifest.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"manifest plugin ! ???+ source \" />\n
    <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\"
    href=\"/post.css\" />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script
    src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    </head>\n    <body>\n<div class='container flex flex-row min-h-screen'>\n
    \   <div>\n    </div>\n    <div class='flex-grow px-8 mx-auto min-h-screen'>\n<header
    class='flex justify-center items-center p-8'>\n\n    <nav class='flex justify-center
    items-center my-8'>\n        <a\n            href='/'>markata</a>\n        <a\n
    \           href='https://github.com/WaylonWalker/markata'>GitHub</a>\n        <a\n
    \           href='https://markata.dev/docs/'>docs</a>\n        <a\n            href='https://markata.dev/plugins/'>plugins</a>\n
    \   </nav>\n\n    <div>\n        <label id=\"theme-switch\" class=\"theme-switch\"
    for=\"checkbox-theme\" title=\"light/dark mode toggle\">\n            <input type=\"checkbox\"
    id=\"checkbox-theme\" />\n            <div class=\"slider round\"></div>\n        </label>\n
    \   </div>\n</header><article class='w-full'>\n<section class=\"title\">\n    <h1
    id=\"title\">\n        Manifest.Py\n    </h1>\n</section>    <section class=\"body\">\n
    \       <p>manifest plugin</p>\n<p>!! function <h2 id='render' class='admonition-title'
    style='margin:0;padding:.5rem 1rem;'>render <em class='small'>function</em></h2></p>\n<div
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
    class=\"p\">:</span> <span class=\"s2\">&quot;MarkataIcons&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">icons</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">icons</span> <span class=\"k\">if</span>
    <span class=\"s2\">&quot;icons&quot;</span> <span class=\"ow\">in</span> <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"vm\">__dict__</span>
    <span class=\"k\">else</span> <span class=\"p\">[]</span>\n            <span class=\"n\">manifest</span>
    <span class=\"o\">=</span> <span class=\"p\">{</span>\n                <span class=\"s2\">&quot;name&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">site_name</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;short_name&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">short_name</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;start_url&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">start_url</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;display&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">display</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;background_color&quot;</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">background_color</span><span class=\"p\">),</span>\n
    \               <span class=\"s2\">&quot;theme_color&quot;</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">theme_color</span><span class=\"p\">),</span>\n                <span
    class=\"s2\">&quot;description&quot;</span><span class=\"p\">:</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">description</span><span class=\"p\">,</span>\n                <span
    class=\"s2\">&quot;icons&quot;</span><span class=\"p\">:</span> <span class=\"n\">icons</span><span
    class=\"p\">,</span>\n            <span class=\"p\">}</span>\n            <span
    class=\"n\">filepath</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"p\">[</span><span class=\"s2\">&quot;output_dir&quot;</span><span
    class=\"p\">])</span> <span class=\"o\">/</span> <span class=\"s2\">&quot;manifest.json&quot;</span>\n
    \           <span class=\"n\">filepath</span><span class=\"o\">.</span><span class=\"n\">touch</span><span
    class=\"p\">(</span><span class=\"n\">exist_ok</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n            <span class=\"k\">with</span>
    <span class=\"nb\">open</span><span class=\"p\">(</span><span class=\"n\">filepath</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;w+&quot;</span><span class=\"p\">)</span>
    <span class=\"k\">as</span> <span class=\"n\">f</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">json</span><span class=\"o\">.</span><span class=\"n\">dump</span><span
    class=\"p\">(</span><span class=\"n\">manifest</span><span class=\"p\">,</span>
    <span class=\"n\">f</span><span class=\"p\">,</span> <span class=\"n\">ensure_ascii</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
    class=\"n\">indent</span><span class=\"o\">=</span><span class=\"mi\">4</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>    </div>\n
    \   <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Manifest.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"manifest plugin ! ???+ source \" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Manifest.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"manifest plugin ! ???+ source \" />\n
    <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\"
    href=\"/post.css\" />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script
    src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    </head>\n    <body>\n<article style=\"text-align: center;\">\n    <style>\n
    \       section {\n            font-size: 200%;\n        }\n\n\n        .edit
    {\n            display: none;\n        }\n    </style>\n<section class=\"title\">\n
    \   <h1 id=\"title\">\n        Manifest.Py\n    </h1>\n</section></article>\n
    \    </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Manifest.Py\n    </h1>\n</section>    <section class=\"body\">\n        <p>manifest
    plugin</p>\n<p>!! function <h2 id='render' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>render <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">render
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
    class=\"p\">:</span> <span class=\"s2\">&quot;MarkataIcons&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">icons</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">icons</span> <span class=\"k\">if</span>
    <span class=\"s2\">&quot;icons&quot;</span> <span class=\"ow\">in</span> <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"vm\">__dict__</span>
    <span class=\"k\">else</span> <span class=\"p\">[]</span>\n            <span class=\"n\">manifest</span>
    <span class=\"o\">=</span> <span class=\"p\">{</span>\n                <span class=\"s2\">&quot;name&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">site_name</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;short_name&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">short_name</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;start_url&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">start_url</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;display&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">display</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;background_color&quot;</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">background_color</span><span class=\"p\">),</span>\n
    \               <span class=\"s2\">&quot;theme_color&quot;</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">theme_color</span><span class=\"p\">),</span>\n                <span
    class=\"s2\">&quot;description&quot;</span><span class=\"p\">:</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">description</span><span class=\"p\">,</span>\n                <span
    class=\"s2\">&quot;icons&quot;</span><span class=\"p\">:</span> <span class=\"n\">icons</span><span
    class=\"p\">,</span>\n            <span class=\"p\">}</span>\n            <span
    class=\"n\">filepath</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"p\">[</span><span class=\"s2\">&quot;output_dir&quot;</span><span
    class=\"p\">])</span> <span class=\"o\">/</span> <span class=\"s2\">&quot;manifest.json&quot;</span>\n
    \           <span class=\"n\">filepath</span><span class=\"o\">.</span><span class=\"n\">touch</span><span
    class=\"p\">(</span><span class=\"n\">exist_ok</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n            <span class=\"k\">with</span>
    <span class=\"nb\">open</span><span class=\"p\">(</span><span class=\"n\">filepath</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;w+&quot;</span><span class=\"p\">)</span>
    <span class=\"k\">as</span> <span class=\"n\">f</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">json</span><span class=\"o\">.</span><span class=\"n\">dump</span><span
    class=\"p\">(</span><span class=\"n\">manifest</span><span class=\"p\">,</span>
    <span class=\"n\">f</span><span class=\"p\">,</span> <span class=\"n\">ensure_ascii</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
    class=\"n\">indent</span><span class=\"o\">=</span><span class=\"mi\">4</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>"
  raw.md: ''
published: true
slug: markata/plugins/manifest
title: Manifest.Py


---

manifest plugin


!! function <h2 id='render' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>render <em class='small'>function</em></h2>

???+ source "render <em class='small'>source</em>"

```python

        def render(markata: "MarkataIcons") -> None:
            icons = markata.icons if "icons" in markata.__dict__ else []
            manifest = {
                "name": markata.config.site_name,
                "short_name": markata.config.short_name,
                "start_url": markata.config.start_url,
                "display": markata.config.display,
                "background_color": str(markata.config.background_color),
                "theme_color": str(markata.config.theme_color),
                "description": markata.config.description,
                "icons": icons,
            }
            filepath = Path(markata.config["output_dir"]) / "manifest.json"
            filepath.touch(exist_ok=True)
            with open(filepath, "w+") as f:
                json.dump(manifest, f, ensure_ascii=True, indent=4)
```

