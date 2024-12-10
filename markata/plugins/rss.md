---
content: "Default glob plugin\n\n\n!! function <h2 id='render' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>render <em class='small'>function</em></h2>\n\n???+
  source \"render <em class='small'>source</em>\"\n\n```python\n\n        def render(markata:
  \"MarkataRss\") -> None:\n            fg = FeedGenerator()\n            url = markata.config.url
  or \"\"\n            title = markata.config.title\n            name = markata.config.author_name\n
  \           email = markata.config.author_email\n            icon = str(markata.config.icon)\n
  \           lang = markata.config.lang\n            rss_description = markata.config.rss_description
  or \"rss feed\"\n\n            fg.id(str(url) + \"/rss.xml\")\n            fg.title(title)\n
  \           fg.author(\n                {\n                    \"name\": name,\n
  \                   \"email\": email,\n                },\n            )\n            fg.link(href=str(url),
  rel=\"alternate\")\n            fg.logo(icon)\n            fg.subtitle(rss_description)\n
  \           fg.link(href=str(url) + \"/rss.xml\", rel=\"self\")\n            fg.language(lang)\n\n
  \           try:\n                all_posts = sorted(markata.articles, key=lambda
  x: x[\"date\"], reverse=True)\n                posts = [post for post in all_posts
  if post[\"published\"] == \"True\"]\n            except BaseException:\n                posts
  = markata.articles\n\n            for article in posts:\n                fe = fg.add_entry()\n
  \               fe.id(str(url + \"/\" + article.slug))\n                fe.title(article.title)\n
  \               fe.published(\n                    datetime.datetime.combine(\n
  \                       article.date or datetime.datetime.min.date(),\n                        datetime.datetime.min.time(),\n
  \                       pytz.UTC,\n                    )\n                )\n                fe.description(article.description)\n
  \               fe.summary(article.long_description)\n                fe.link(href=str(url)
  + \"/\" + article.slug)\n                fe.content(article.article_html.translate(dict.fromkeys(range(32))))\n\n
  \           markata.fg = fg\n            markata.rss = fg.rss_str(pretty=True)\n```\n\n\n!!
  function <h2 id='save' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>save
  <em class='small'>function</em></h2>\n\n???+ source \"save <em class='small'>source</em>\"\n\n```python\n\n
  \       def save(markata: \"MarkataRss\") -> None:\n            output_dir = Path(markata.config[\"output_dir\"])\n
  \           markata.fg.rss_file(str(output_dir / \"rss.xml\"))\n```\n\n\n!! class
  <h2 id='MarkataRss' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>MarkataRss
  <em class='small'>class</em></h2>\n\n???+ source \"MarkataRss <em class='small'>source</em>\"\n\n```python\n\n
  \       class MarkataRss(Markata):\n                fg: \"FeedGenerator\"\n                rss:
  str\n```\n\n"
date: 0001-01-01
description: 'Default glob plugin ! ???+ source  ! ???+ source  ! ???+ source '
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Rss.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"Default glob plugin ! ???+ source  ! ???+ source
    \ ! ???+ source \" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Rss.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Default glob plugin ! ???+ source  !
    ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Rss.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <p>Default glob plugin</p>\n<p>!!
    function <h2 id='render' class='admonition-title' style='margin:0;padding:.5rem
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
    class=\"p\">:</span> <span class=\"s2\">&quot;MarkataRss&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">fg</span> <span class=\"o\">=</span> <span class=\"n\">FeedGenerator</span><span
    class=\"p\">()</span>\n            <span class=\"n\">url</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">url</span> <span class=\"ow\">or</span>
    <span class=\"s2\">&quot;&quot;</span>\n            <span class=\"n\">title</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">title</span>\n
    \           <span class=\"n\">name</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">author_name</span>\n            <span class=\"n\">email</span> <span
    class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">author_email</span>\n
    \           <span class=\"n\">icon</span> <span class=\"o\">=</span> <span class=\"nb\">str</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">icon</span><span
    class=\"p\">)</span>\n            <span class=\"n\">lang</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">lang</span>\n            <span class=\"n\">rss_description</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">rss_description</span>
    <span class=\"ow\">or</span> <span class=\"s2\">&quot;rss feed&quot;</span>\n\n
    \           <span class=\"n\">fg</span><span class=\"o\">.</span><span class=\"n\">id</span><span
    class=\"p\">(</span><span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">url</span><span class=\"p\">)</span> <span class=\"o\">+</span> <span
    class=\"s2\">&quot;/rss.xml&quot;</span><span class=\"p\">)</span>\n            <span
    class=\"n\">fg</span><span class=\"o\">.</span><span class=\"n\">title</span><span
    class=\"p\">(</span><span class=\"n\">title</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">fg</span><span class=\"o\">.</span><span class=\"n\">author</span><span
    class=\"p\">(</span>\n                <span class=\"p\">{</span>\n                    <span
    class=\"s2\">&quot;name&quot;</span><span class=\"p\">:</span> <span class=\"n\">name</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;email&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">email</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">},</span>\n            <span class=\"p\">)</span>\n
    \           <span class=\"n\">fg</span><span class=\"o\">.</span><span class=\"n\">link</span><span
    class=\"p\">(</span><span class=\"n\">href</span><span class=\"o\">=</span><span
    class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">url</span><span
    class=\"p\">),</span> <span class=\"n\">rel</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;alternate&quot;</span><span class=\"p\">)</span>\n            <span
    class=\"n\">fg</span><span class=\"o\">.</span><span class=\"n\">logo</span><span
    class=\"p\">(</span><span class=\"n\">icon</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">fg</span><span class=\"o\">.</span><span class=\"n\">subtitle</span><span
    class=\"p\">(</span><span class=\"n\">rss_description</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">fg</span><span class=\"o\">.</span><span class=\"n\">link</span><span
    class=\"p\">(</span><span class=\"n\">href</span><span class=\"o\">=</span><span
    class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">url</span><span
    class=\"p\">)</span> <span class=\"o\">+</span> <span class=\"s2\">&quot;/rss.xml&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">rel</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;self&quot;</span><span class=\"p\">)</span>\n            <span
    class=\"n\">fg</span><span class=\"o\">.</span><span class=\"n\">language</span><span
    class=\"p\">(</span><span class=\"n\">lang</span><span class=\"p\">)</span>\n\n
    \           <span class=\"k\">try</span><span class=\"p\">:</span>\n                <span
    class=\"n\">all_posts</span> <span class=\"o\">=</span> <span class=\"nb\">sorted</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">articles</span><span class=\"p\">,</span> <span class=\"n\">key</span><span
    class=\"o\">=</span><span class=\"k\">lambda</span> <span class=\"n\">x</span><span
    class=\"p\">:</span> <span class=\"n\">x</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;date&quot;</span><span class=\"p\">],</span> <span class=\"n\">reverse</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \               <span class=\"n\">posts</span> <span class=\"o\">=</span> <span
    class=\"p\">[</span><span class=\"n\">post</span> <span class=\"k\">for</span>
    <span class=\"n\">post</span> <span class=\"ow\">in</span> <span class=\"n\">all_posts</span>
    <span class=\"k\">if</span> <span class=\"n\">post</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;published&quot;</span><span class=\"p\">]</span> <span class=\"o\">==</span>
    <span class=\"s2\">&quot;True&quot;</span><span class=\"p\">]</span>\n            <span
    class=\"k\">except</span> <span class=\"ne\">BaseException</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">posts</span> <span class=\"o\">=</span> <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">articles</span>\n\n
    \           <span class=\"k\">for</span> <span class=\"n\">article</span> <span
    class=\"ow\">in</span> <span class=\"n\">posts</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">fe</span> <span class=\"o\">=</span> <span class=\"n\">fg</span><span
    class=\"o\">.</span><span class=\"n\">add_entry</span><span class=\"p\">()</span>\n
    \               <span class=\"n\">fe</span><span class=\"o\">.</span><span class=\"n\">id</span><span
    class=\"p\">(</span><span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">url</span> <span class=\"o\">+</span> <span class=\"s2\">&quot;/&quot;</span>
    <span class=\"o\">+</span> <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">slug</span><span class=\"p\">))</span>\n                <span class=\"n\">fe</span><span
    class=\"o\">.</span><span class=\"n\">title</span><span class=\"p\">(</span><span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">title</span><span
    class=\"p\">)</span>\n                <span class=\"n\">fe</span><span class=\"o\">.</span><span
    class=\"n\">published</span><span class=\"p\">(</span>\n                    <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">combine</span><span class=\"p\">(</span>\n
    \                       <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">date</span> <span class=\"ow\">or</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">min</span><span class=\"o\">.</span><span class=\"n\">date</span><span
    class=\"p\">(),</span>\n                        <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">min</span><span class=\"o\">.</span><span class=\"n\">time</span><span
    class=\"p\">(),</span>\n                        <span class=\"n\">pytz</span><span
    class=\"o\">.</span><span class=\"n\">UTC</span><span class=\"p\">,</span>\n                    <span
    class=\"p\">)</span>\n                <span class=\"p\">)</span>\n                <span
    class=\"n\">fe</span><span class=\"o\">.</span><span class=\"n\">description</span><span
    class=\"p\">(</span><span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">description</span><span class=\"p\">)</span>\n                <span
    class=\"n\">fe</span><span class=\"o\">.</span><span class=\"n\">summary</span><span
    class=\"p\">(</span><span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">long_description</span><span class=\"p\">)</span>\n                <span
    class=\"n\">fe</span><span class=\"o\">.</span><span class=\"n\">link</span><span
    class=\"p\">(</span><span class=\"n\">href</span><span class=\"o\">=</span><span
    class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">url</span><span
    class=\"p\">)</span> <span class=\"o\">+</span> <span class=\"s2\">&quot;/&quot;</span>
    <span class=\"o\">+</span> <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">slug</span><span class=\"p\">)</span>\n                <span class=\"n\">fe</span><span
    class=\"o\">.</span><span class=\"n\">content</span><span class=\"p\">(</span><span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">article_html</span><span
    class=\"o\">.</span><span class=\"n\">translate</span><span class=\"p\">(</span><span
    class=\"nb\">dict</span><span class=\"o\">.</span><span class=\"n\">fromkeys</span><span
    class=\"p\">(</span><span class=\"nb\">range</span><span class=\"p\">(</span><span
    class=\"mi\">32</span><span class=\"p\">))))</span>\n\n            <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">fg</span> <span class=\"o\">=</span> <span
    class=\"n\">fg</span>\n            <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">rss</span> <span class=\"o\">=</span> <span class=\"n\">fg</span><span
    class=\"o\">.</span><span class=\"n\">rss_str</span><span class=\"p\">(</span><span
    class=\"n\">pretty</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='save' class='admonition-title'
    style='margin:0;padding:.5rem 1rem;'>save <em class='small'>function</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">save
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
    class=\"p\">:</span> <span class=\"s2\">&quot;MarkataRss&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">output_dir</span> <span class=\"o\">=</span> <span
    class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;output_dir&quot;</span><span class=\"p\">])</span>\n            <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">fg</span><span
    class=\"o\">.</span><span class=\"n\">rss_file</span><span class=\"p\">(</span><span
    class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">output_dir</span>
    <span class=\"o\">/</span> <span class=\"s2\">&quot;rss.xml&quot;</span><span
    class=\"p\">))</span>\n</pre></div>\n\n</pre>\n\n<p>!! class <h2 id='MarkataRss'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>MarkataRss <em class='small'>class</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">MarkataRss
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
    <span class=\"nc\">MarkataRss</span><span class=\"p\">(</span><span class=\"n\">Markata</span><span
    class=\"p\">):</span>\n                <span class=\"n\">fg</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;FeedGenerator&quot;</span>\n                <span class=\"n\">rss</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>    </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Rss.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"Default glob plugin ! ???+ source  ! ???+ source
    \ ! ???+ source \" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Rss.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Default glob plugin ! ???+ source  !
    ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    </head>\n    <body>\n<article style=\"text-align:
    center;\">\n    <style>\n        section {\n            font-size: 200%;\n        }\n\n\n
    \       .edit {\n            display: none;\n        }\n    </style>\n<section
    class=\"title\">\n    <h1 id=\"title\">\n        Rss.Py\n    </h1>\n</section></article>\n
    \    </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Rss.Py\n    </h1>\n</section>    <section class=\"body\">\n        <p>Default
    glob plugin</p>\n<p>!! function <h2 id='render' class='admonition-title' style='margin:0;padding:.5rem
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
    class=\"p\">:</span> <span class=\"s2\">&quot;MarkataRss&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">fg</span> <span class=\"o\">=</span> <span class=\"n\">FeedGenerator</span><span
    class=\"p\">()</span>\n            <span class=\"n\">url</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">url</span> <span class=\"ow\">or</span>
    <span class=\"s2\">&quot;&quot;</span>\n            <span class=\"n\">title</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">title</span>\n
    \           <span class=\"n\">name</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">author_name</span>\n            <span class=\"n\">email</span> <span
    class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">author_email</span>\n
    \           <span class=\"n\">icon</span> <span class=\"o\">=</span> <span class=\"nb\">str</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">icon</span><span
    class=\"p\">)</span>\n            <span class=\"n\">lang</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">lang</span>\n            <span class=\"n\">rss_description</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">rss_description</span>
    <span class=\"ow\">or</span> <span class=\"s2\">&quot;rss feed&quot;</span>\n\n
    \           <span class=\"n\">fg</span><span class=\"o\">.</span><span class=\"n\">id</span><span
    class=\"p\">(</span><span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">url</span><span class=\"p\">)</span> <span class=\"o\">+</span> <span
    class=\"s2\">&quot;/rss.xml&quot;</span><span class=\"p\">)</span>\n            <span
    class=\"n\">fg</span><span class=\"o\">.</span><span class=\"n\">title</span><span
    class=\"p\">(</span><span class=\"n\">title</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">fg</span><span class=\"o\">.</span><span class=\"n\">author</span><span
    class=\"p\">(</span>\n                <span class=\"p\">{</span>\n                    <span
    class=\"s2\">&quot;name&quot;</span><span class=\"p\">:</span> <span class=\"n\">name</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;email&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">email</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">},</span>\n            <span class=\"p\">)</span>\n
    \           <span class=\"n\">fg</span><span class=\"o\">.</span><span class=\"n\">link</span><span
    class=\"p\">(</span><span class=\"n\">href</span><span class=\"o\">=</span><span
    class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">url</span><span
    class=\"p\">),</span> <span class=\"n\">rel</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;alternate&quot;</span><span class=\"p\">)</span>\n            <span
    class=\"n\">fg</span><span class=\"o\">.</span><span class=\"n\">logo</span><span
    class=\"p\">(</span><span class=\"n\">icon</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">fg</span><span class=\"o\">.</span><span class=\"n\">subtitle</span><span
    class=\"p\">(</span><span class=\"n\">rss_description</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">fg</span><span class=\"o\">.</span><span class=\"n\">link</span><span
    class=\"p\">(</span><span class=\"n\">href</span><span class=\"o\">=</span><span
    class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">url</span><span
    class=\"p\">)</span> <span class=\"o\">+</span> <span class=\"s2\">&quot;/rss.xml&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">rel</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;self&quot;</span><span class=\"p\">)</span>\n            <span
    class=\"n\">fg</span><span class=\"o\">.</span><span class=\"n\">language</span><span
    class=\"p\">(</span><span class=\"n\">lang</span><span class=\"p\">)</span>\n\n
    \           <span class=\"k\">try</span><span class=\"p\">:</span>\n                <span
    class=\"n\">all_posts</span> <span class=\"o\">=</span> <span class=\"nb\">sorted</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">articles</span><span class=\"p\">,</span> <span class=\"n\">key</span><span
    class=\"o\">=</span><span class=\"k\">lambda</span> <span class=\"n\">x</span><span
    class=\"p\">:</span> <span class=\"n\">x</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;date&quot;</span><span class=\"p\">],</span> <span class=\"n\">reverse</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \               <span class=\"n\">posts</span> <span class=\"o\">=</span> <span
    class=\"p\">[</span><span class=\"n\">post</span> <span class=\"k\">for</span>
    <span class=\"n\">post</span> <span class=\"ow\">in</span> <span class=\"n\">all_posts</span>
    <span class=\"k\">if</span> <span class=\"n\">post</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;published&quot;</span><span class=\"p\">]</span> <span class=\"o\">==</span>
    <span class=\"s2\">&quot;True&quot;</span><span class=\"p\">]</span>\n            <span
    class=\"k\">except</span> <span class=\"ne\">BaseException</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">posts</span> <span class=\"o\">=</span> <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">articles</span>\n\n
    \           <span class=\"k\">for</span> <span class=\"n\">article</span> <span
    class=\"ow\">in</span> <span class=\"n\">posts</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">fe</span> <span class=\"o\">=</span> <span class=\"n\">fg</span><span
    class=\"o\">.</span><span class=\"n\">add_entry</span><span class=\"p\">()</span>\n
    \               <span class=\"n\">fe</span><span class=\"o\">.</span><span class=\"n\">id</span><span
    class=\"p\">(</span><span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">url</span> <span class=\"o\">+</span> <span class=\"s2\">&quot;/&quot;</span>
    <span class=\"o\">+</span> <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">slug</span><span class=\"p\">))</span>\n                <span class=\"n\">fe</span><span
    class=\"o\">.</span><span class=\"n\">title</span><span class=\"p\">(</span><span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">title</span><span
    class=\"p\">)</span>\n                <span class=\"n\">fe</span><span class=\"o\">.</span><span
    class=\"n\">published</span><span class=\"p\">(</span>\n                    <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">combine</span><span class=\"p\">(</span>\n
    \                       <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">date</span> <span class=\"ow\">or</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">min</span><span class=\"o\">.</span><span class=\"n\">date</span><span
    class=\"p\">(),</span>\n                        <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">min</span><span class=\"o\">.</span><span class=\"n\">time</span><span
    class=\"p\">(),</span>\n                        <span class=\"n\">pytz</span><span
    class=\"o\">.</span><span class=\"n\">UTC</span><span class=\"p\">,</span>\n                    <span
    class=\"p\">)</span>\n                <span class=\"p\">)</span>\n                <span
    class=\"n\">fe</span><span class=\"o\">.</span><span class=\"n\">description</span><span
    class=\"p\">(</span><span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">description</span><span class=\"p\">)</span>\n                <span
    class=\"n\">fe</span><span class=\"o\">.</span><span class=\"n\">summary</span><span
    class=\"p\">(</span><span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">long_description</span><span class=\"p\">)</span>\n                <span
    class=\"n\">fe</span><span class=\"o\">.</span><span class=\"n\">link</span><span
    class=\"p\">(</span><span class=\"n\">href</span><span class=\"o\">=</span><span
    class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">url</span><span
    class=\"p\">)</span> <span class=\"o\">+</span> <span class=\"s2\">&quot;/&quot;</span>
    <span class=\"o\">+</span> <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">slug</span><span class=\"p\">)</span>\n                <span class=\"n\">fe</span><span
    class=\"o\">.</span><span class=\"n\">content</span><span class=\"p\">(</span><span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">article_html</span><span
    class=\"o\">.</span><span class=\"n\">translate</span><span class=\"p\">(</span><span
    class=\"nb\">dict</span><span class=\"o\">.</span><span class=\"n\">fromkeys</span><span
    class=\"p\">(</span><span class=\"nb\">range</span><span class=\"p\">(</span><span
    class=\"mi\">32</span><span class=\"p\">))))</span>\n\n            <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">fg</span> <span class=\"o\">=</span> <span
    class=\"n\">fg</span>\n            <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">rss</span> <span class=\"o\">=</span> <span class=\"n\">fg</span><span
    class=\"o\">.</span><span class=\"n\">rss_str</span><span class=\"p\">(</span><span
    class=\"n\">pretty</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='save' class='admonition-title'
    style='margin:0;padding:.5rem 1rem;'>save <em class='small'>function</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">save
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
    class=\"p\">:</span> <span class=\"s2\">&quot;MarkataRss&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">output_dir</span> <span class=\"o\">=</span> <span
    class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;output_dir&quot;</span><span class=\"p\">])</span>\n            <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">fg</span><span
    class=\"o\">.</span><span class=\"n\">rss_file</span><span class=\"p\">(</span><span
    class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">output_dir</span>
    <span class=\"o\">/</span> <span class=\"s2\">&quot;rss.xml&quot;</span><span
    class=\"p\">))</span>\n</pre></div>\n\n</pre>\n\n<p>!! class <h2 id='MarkataRss'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>MarkataRss <em class='small'>class</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">MarkataRss
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
    <span class=\"nc\">MarkataRss</span><span class=\"p\">(</span><span class=\"n\">Markata</span><span
    class=\"p\">):</span>\n                <span class=\"n\">fg</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;FeedGenerator&quot;</span>\n                <span class=\"n\">rss</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>"
  raw.md: ''
published: true
slug: markata/plugins/rss
title: Rss.Py


---

Default glob plugin


!! function <h2 id='render' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>render <em class='small'>function</em></h2>

???+ source "render <em class='small'>source</em>"

```python

        def render(markata: "MarkataRss") -> None:
            fg = FeedGenerator()
            url = markata.config.url or ""
            title = markata.config.title
            name = markata.config.author_name
            email = markata.config.author_email
            icon = str(markata.config.icon)
            lang = markata.config.lang
            rss_description = markata.config.rss_description or "rss feed"

            fg.id(str(url) + "/rss.xml")
            fg.title(title)
            fg.author(
                {
                    "name": name,
                    "email": email,
                },
            )
            fg.link(href=str(url), rel="alternate")
            fg.logo(icon)
            fg.subtitle(rss_description)
            fg.link(href=str(url) + "/rss.xml", rel="self")
            fg.language(lang)

            try:
                all_posts = sorted(markata.articles, key=lambda x: x["date"], reverse=True)
                posts = [post for post in all_posts if post["published"] == "True"]
            except BaseException:
                posts = markata.articles

            for article in posts:
                fe = fg.add_entry()
                fe.id(str(url + "/" + article.slug))
                fe.title(article.title)
                fe.published(
                    datetime.datetime.combine(
                        article.date or datetime.datetime.min.date(),
                        datetime.datetime.min.time(),
                        pytz.UTC,
                    )
                )
                fe.description(article.description)
                fe.summary(article.long_description)
                fe.link(href=str(url) + "/" + article.slug)
                fe.content(article.article_html.translate(dict.fromkeys(range(32))))

            markata.fg = fg
            markata.rss = fg.rss_str(pretty=True)
```


!! function <h2 id='save' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>save <em class='small'>function</em></h2>

???+ source "save <em class='small'>source</em>"

```python

        def save(markata: "MarkataRss") -> None:
            output_dir = Path(markata.config["output_dir"])
            markata.fg.rss_file(str(output_dir / "rss.xml"))
```


!! class <h2 id='MarkataRss' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>MarkataRss <em class='small'>class</em></h2>

???+ source "MarkataRss <em class='small'>source</em>"

```python

        class MarkataRss(Markata):
                fg: "FeedGenerator"
                rss: str
```

