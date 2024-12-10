---
content: "Default datetime plugin\n\n\n!! function <h2 id='load' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>load <em class='small'>function</em></h2>\n\n???+
  source \"load <em class='small'>source</em>\"\n\n```python\n\n        def load(markata:
  \"Markata\") -> None:\n            for article in markata.iter_articles(\"datetime\"):\n
  \               try:\n                    date = article.metadata[\"date\"]\n                except
  KeyError:\n                    date = None\n                if isinstance(date,
  str):\n                    date = dateutil.parser.parse(date)\n                if
  isinstance(date, datetime.date):\n                    date = datetime.datetime(\n
  \                       year=date.year,\n                        month=date.month,\n
  \                       day=date.day,\n                        tzinfo=pytz.utc,\n
  \                   )\n\n                article[\"today\"] = datetime.date.today()\n
  \               article[\"now\"] = datetime.datetime.now()\n                article[\"datetime\"]
  = date\n```\n\n"
date: 0001-01-01
description: 'Default datetime plugin ! ???+ source '
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Datetime.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"Default datetime plugin ! ???+ source \" />\n <link
    href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\"
    href=\"/post.css\" />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script
    src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    <head>\n<title>Datetime.Py</title>\n<meta charset=\"UTF-8\" />\n<meta
    name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta name=\"description\"
    content=\"Default datetime plugin ! ???+ source \" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Datetime.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <p>Default datetime
    plugin</p>\n<p>!! function <h2 id='load' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>load <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">load
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
    <span class=\"nf\">load</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">for</span> <span class=\"n\">article</span> <span
    class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">iter_articles</span><span class=\"p\">(</span><span class=\"s2\">&quot;datetime&quot;</span><span
    class=\"p\">):</span>\n                <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">date</span> <span class=\"o\">=</span> <span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">metadata</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;date&quot;</span><span class=\"p\">]</span>\n
    \               <span class=\"k\">except</span> <span class=\"ne\">KeyError</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">date</span> <span
    class=\"o\">=</span> <span class=\"kc\">None</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">date</span><span
    class=\"p\">,</span> <span class=\"nb\">str</span><span class=\"p\">):</span>\n
    \                   <span class=\"n\">date</span> <span class=\"o\">=</span> <span
    class=\"n\">dateutil</span><span class=\"o\">.</span><span class=\"n\">parser</span><span
    class=\"o\">.</span><span class=\"n\">parse</span><span class=\"p\">(</span><span
    class=\"n\">date</span><span class=\"p\">)</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">date</span><span
    class=\"p\">,</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">date</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">date</span> <span class=\"o\">=</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"p\">(</span>\n
    \                       <span class=\"n\">year</span><span class=\"o\">=</span><span
    class=\"n\">date</span><span class=\"o\">.</span><span class=\"n\">year</span><span
    class=\"p\">,</span>\n                        <span class=\"n\">month</span><span
    class=\"o\">=</span><span class=\"n\">date</span><span class=\"o\">.</span><span
    class=\"n\">month</span><span class=\"p\">,</span>\n                        <span
    class=\"n\">day</span><span class=\"o\">=</span><span class=\"n\">date</span><span
    class=\"o\">.</span><span class=\"n\">day</span><span class=\"p\">,</span>\n                        <span
    class=\"n\">tzinfo</span><span class=\"o\">=</span><span class=\"n\">pytz</span><span
    class=\"o\">.</span><span class=\"n\">utc</span><span class=\"p\">,</span>\n                    <span
    class=\"p\">)</span>\n\n                <span class=\"n\">article</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;today&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">date</span><span class=\"o\">.</span><span class=\"n\">today</span><span
    class=\"p\">()</span>\n                <span class=\"n\">article</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;now&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">now</span><span class=\"p\">()</span>\n
    \               <span class=\"n\">article</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;datetime&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"n\">date</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>
    \   </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Datetime.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"Default datetime plugin ! ???+ source \" />\n <link
    href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\"
    href=\"/post.css\" />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script
    src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    <head>\n<title>Datetime.Py</title>\n<meta charset=\"UTF-8\" />\n<meta
    name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta name=\"description\"
    content=\"Default datetime plugin ! ???+ source \" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<article style=\"text-align: center;\">\n    <style>\n        section
    {\n            font-size: 200%;\n        }\n\n\n        .edit {\n            display:
    none;\n        }\n    </style>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Datetime.Py\n    </h1>\n</section></article>\n     </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Datetime.Py\n    </h1>\n</section>    <section class=\"body\">\n        <p>Default
    datetime plugin</p>\n<p>!! function <h2 id='load' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>load <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">load
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
    <span class=\"nf\">load</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">for</span> <span class=\"n\">article</span> <span
    class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">iter_articles</span><span class=\"p\">(</span><span class=\"s2\">&quot;datetime&quot;</span><span
    class=\"p\">):</span>\n                <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">date</span> <span class=\"o\">=</span> <span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">metadata</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;date&quot;</span><span class=\"p\">]</span>\n
    \               <span class=\"k\">except</span> <span class=\"ne\">KeyError</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">date</span> <span
    class=\"o\">=</span> <span class=\"kc\">None</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">date</span><span
    class=\"p\">,</span> <span class=\"nb\">str</span><span class=\"p\">):</span>\n
    \                   <span class=\"n\">date</span> <span class=\"o\">=</span> <span
    class=\"n\">dateutil</span><span class=\"o\">.</span><span class=\"n\">parser</span><span
    class=\"o\">.</span><span class=\"n\">parse</span><span class=\"p\">(</span><span
    class=\"n\">date</span><span class=\"p\">)</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">date</span><span
    class=\"p\">,</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">date</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">date</span> <span class=\"o\">=</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"p\">(</span>\n
    \                       <span class=\"n\">year</span><span class=\"o\">=</span><span
    class=\"n\">date</span><span class=\"o\">.</span><span class=\"n\">year</span><span
    class=\"p\">,</span>\n                        <span class=\"n\">month</span><span
    class=\"o\">=</span><span class=\"n\">date</span><span class=\"o\">.</span><span
    class=\"n\">month</span><span class=\"p\">,</span>\n                        <span
    class=\"n\">day</span><span class=\"o\">=</span><span class=\"n\">date</span><span
    class=\"o\">.</span><span class=\"n\">day</span><span class=\"p\">,</span>\n                        <span
    class=\"n\">tzinfo</span><span class=\"o\">=</span><span class=\"n\">pytz</span><span
    class=\"o\">.</span><span class=\"n\">utc</span><span class=\"p\">,</span>\n                    <span
    class=\"p\">)</span>\n\n                <span class=\"n\">article</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;today&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">date</span><span class=\"o\">.</span><span class=\"n\">today</span><span
    class=\"p\">()</span>\n                <span class=\"n\">article</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;now&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">now</span><span class=\"p\">()</span>\n
    \               <span class=\"n\">article</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;datetime&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"n\">date</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>"
  raw.md: ''
published: true
slug: markata/plugins/datetime
title: Datetime.Py


---

Default datetime plugin


!! function <h2 id='load' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>load <em class='small'>function</em></h2>

???+ source "load <em class='small'>source</em>"

```python

        def load(markata: "Markata") -> None:
            for article in markata.iter_articles("datetime"):
                try:
                    date = article.metadata["date"]
                except KeyError:
                    date = None
                if isinstance(date, str):
                    date = dateutil.parser.parse(date)
                if isinstance(date, datetime.date):
                    date = datetime.datetime(
                        year=date.year,
                        month=date.month,
                        day=date.day,
                        tzinfo=pytz.utc,
                    )

                article["today"] = datetime.date.today()
                article["now"] = datetime.datetime.now()
                article["datetime"] = date
```

