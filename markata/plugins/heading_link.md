---
content: "Creates links next to all heading tags to make it easier for users to share
  a\nspecific heading.\n\n\n!! function <h2 id='post_render' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>post_render <em class='small'>function</em></h2>\n
  \   This plugin creates a link svg next to all headings.\n???+ source \"post_render
  <em class='small'>source</em>\"\n\n```python\n\n        def post_render(markata:
  Markata) -> None:\n            \"\"\"\n            This plugin creates a link svg
  next to all headings.\n            \"\"\"\n\n            with markata.cache as cache:\n
  \               for article in markata.iter_articles(\"link headers\"):\n                    key
  = markata.make_hash(\n                        \"heading_link\",\n                        \"post_render\",\n
  \                       Path(__file__).read_text(),\n                        article.content,\n
  \                       article.html,\n                    )\n\n                    html_from_cache
  = markata.precache.get(key)\n\n                    if html_from_cache is None:\n
  \                       html = link_headings(article)\n                        cache.add(\n
  \                           key,\n                            html,\n                            expire=markata.config.default_cache_expire,\n
  \                       )\n                    else:\n                        html
  = html_from_cache\n                    article.html = html\n```\n\n\n!! function
  <h2 id='link_headings' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>link_headings
  <em class='small'>function</em></h2>\n    Use BeautifulSoup to find all headings
  and run link_heading on them.\n???+ source \"link_headings <em class='small'>source</em>\"\n\n```python\n\n
  \       def link_headings(article: \"Post\") -> str:\n            \"\"\"\n            Use
  BeautifulSoup to find all headings and run link_heading on them.\n            \"\"\"\n
  \           soup = BeautifulSoup(article.html, \"html.parser\")\n            for
  heading in soup.find_all(re.compile(\"^h[1-6]$\")):\n                if (\n                    not
  heading.find(\"a\", {\"class\": \"heading-permalink\"})\n                    and
  heading.get(\"id\", \"\") != \"title\"\n                ):\n                    link_heading(soup,
  heading)\n            return str(soup)\n```\n\n\n!! function <h2 id='link_heading'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>link_heading <em class='small'>function</em></h2>\n
  \   Mutate soup to include an svg link at the heading passed in.\n???+ source \"link_heading
  <em class='small'>source</em>\"\n\n```python\n\n        def link_heading(soup: \"bs4.BeautifulSoup\",
  heading: \"bs4.element.Tag\") -> None:\n            \"\"\"\n            Mutate soup
  to include an svg link at the heading passed in.\n            \"\"\"\n            id
  = heading.get(\"id\")\n\n            link = soup.new_tag(\n                \"a\",\n
  \               alt=\"id\",\n                title=f\"link to #{id}\",\n                href=f\"#{id}\",\n
  \               **{\"class\": \"heading-permalink\"},\n            )\n            span
  = soup.new_tag(\"span\", **{\"class\": \"visually-hidden\"})\n            svg =
  soup.new_tag(\n                \"svg\",\n                fill=\"currentColor\",\n
  \               focusable=\"false\",\n                width=\"1em\",\n                height=\"1em\",\n
  \               xmlns=\"http://www.w3.org/2000/svg\",\n                viewBox=\"0
  0 24 24\",\n                **{\n                    \"aria-hidden\": \"true\",\n
  \               },\n            )\n\n            path = soup.new_tag(\n                \"path\",\n
  \               d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0
  5.105-1.702l2.995-2.994a5.992 5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211
  5.99 5.99 0 0 0-6.431-1.242 6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0
  1 0 1.41 1.418l1.709-1.699a3.985 3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799
  1.122 3.997 3.997 0 0 1 .111 5.644l-3.005 3.006a3.982 3.982 0 0 1-3.395 1.126 3.987
  3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201 13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345
  5.987 5.987 0 0 0-5.105 1.702l-2.995 2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976
  0 0 0 1.831 4.211 5.99 5.99 0 0 0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999
  0 1 0-1.414-1.414L9.836 19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122
  3.997 3.997 0 0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987
  0 0 1 2.632 1.563 1 1 0 0 0 1.602-1.198z\",\n            )\n            svg.append(path)\n
  \           link.append(span)\n            link.append(svg)\n            heading.append(link)\n```\n\n"
date: 0001-01-01
description: Creates links next to all heading tags to make it easier for users to
  share a ! ! !
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Heading_Link.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Creates links next to all heading tags
    to make it easier for users to share a ! ! !\" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Heading_Link.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Creates links next to all heading tags
    to make it easier for users to share a ! ! !\" />\n <link href=\"/favicon.ico\"
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Heading_Link.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <p>Creates links next
    to all heading tags to make it easier for users to share a\nspecific heading.</p>\n<p>!!
    function <h2 id='post_render' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>post_render <em class='small'>function</em></h2>\nThis plugin creates a
    link svg next to all headings.</p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">post_render <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">post_render</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"n\">Markata</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            This plugin creates a link svg next to all headings.</span>\n<span
    class=\"sd\">            &quot;&quot;&quot;</span>\n\n            <span class=\"k\">with</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">cache</span>
    <span class=\"k\">as</span> <span class=\"n\">cache</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">for</span> <span class=\"n\">article</span>
    <span class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">iter_articles</span><span class=\"p\">(</span><span class=\"s2\">&quot;link
    headers&quot;</span><span class=\"p\">):</span>\n                    <span class=\"n\">key</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">make_hash</span><span class=\"p\">(</span>\n                        <span
    class=\"s2\">&quot;heading_link&quot;</span><span class=\"p\">,</span>\n                        <span
    class=\"s2\">&quot;post_render&quot;</span><span class=\"p\">,</span>\n                        <span
    class=\"n\">Path</span><span class=\"p\">(</span><span class=\"vm\">__file__</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">read_text</span><span
    class=\"p\">(),</span>\n                        <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">content</span><span class=\"p\">,</span>\n
    \                       <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">html</span><span class=\"p\">,</span>\n                    <span class=\"p\">)</span>\n\n
    \                   <span class=\"n\">html_from_cache</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">precache</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">)</span>\n\n                    <span
    class=\"k\">if</span> <span class=\"n\">html_from_cache</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">html</span> <span class=\"o\">=</span> <span class=\"n\">link_headings</span><span
    class=\"p\">(</span><span class=\"n\">article</span><span class=\"p\">)</span>\n
    \                       <span class=\"n\">cache</span><span class=\"o\">.</span><span
    class=\"n\">add</span><span class=\"p\">(</span>\n                            <span
    class=\"n\">key</span><span class=\"p\">,</span>\n                            <span
    class=\"n\">html</span><span class=\"p\">,</span>\n                            <span
    class=\"n\">expire</span><span class=\"o\">=</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">default_cache_expire</span><span class=\"p\">,</span>\n                        <span
    class=\"p\">)</span>\n                    <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">html</span> <span class=\"o\">=</span>
    <span class=\"n\">html_from_cache</span>\n                    <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">html</span> <span class=\"o\">=</span> <span
    class=\"n\">html</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='link_headings'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>link_headings <em
    class='small'>function</em></h2>\nUse BeautifulSoup to find all headings and run
    link_heading on them.</p>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">link_headings <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">link_headings</span><span class=\"p\">(</span><span class=\"n\">article</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            Use BeautifulSoup to find all headings and run link_heading
    on them.</span>\n<span class=\"sd\">            &quot;&quot;&quot;</span>\n            <span
    class=\"n\">soup</span> <span class=\"o\">=</span> <span class=\"n\">BeautifulSoup</span><span
    class=\"p\">(</span><span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">html</span><span class=\"p\">,</span> <span class=\"s2\">&quot;html.parser&quot;</span><span
    class=\"p\">)</span>\n            <span class=\"k\">for</span> <span class=\"n\">heading</span>
    <span class=\"ow\">in</span> <span class=\"n\">soup</span><span class=\"o\">.</span><span
    class=\"n\">find_all</span><span class=\"p\">(</span><span class=\"n\">re</span><span
    class=\"o\">.</span><span class=\"n\">compile</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;^h[1-6]$&quot;</span><span class=\"p\">)):</span>\n                <span
    class=\"k\">if</span> <span class=\"p\">(</span>\n                    <span class=\"ow\">not</span>
    <span class=\"n\">heading</span><span class=\"o\">.</span><span class=\"n\">find</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;a&quot;</span><span class=\"p\">,</span>
    <span class=\"p\">{</span><span class=\"s2\">&quot;class&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;heading-permalink&quot;</span><span class=\"p\">})</span>\n
    \                   <span class=\"ow\">and</span> <span class=\"n\">heading</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;id&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">!=</span> <span class=\"s2\">&quot;title&quot;</span>\n
    \               <span class=\"p\">):</span>\n                    <span class=\"n\">link_heading</span><span
    class=\"p\">(</span><span class=\"n\">soup</span><span class=\"p\">,</span> <span
    class=\"n\">heading</span><span class=\"p\">)</span>\n            <span class=\"k\">return</span>
    <span class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">soup</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='link_heading'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>link_heading <em
    class='small'>function</em></h2>\nMutate soup to include an svg link at the heading
    passed in.</p>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">link_heading <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">link_heading</span><span class=\"p\">(</span><span class=\"n\">soup</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;bs4.BeautifulSoup&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">heading</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;bs4.element.Tag&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            Mutate soup to include an svg link at the heading passed
    in.</span>\n<span class=\"sd\">            &quot;&quot;&quot;</span>\n            <span
    class=\"nb\">id</span> <span class=\"o\">=</span> <span class=\"n\">heading</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;id&quot;</span><span class=\"p\">)</span>\n\n            <span
    class=\"n\">link</span> <span class=\"o\">=</span> <span class=\"n\">soup</span><span
    class=\"o\">.</span><span class=\"n\">new_tag</span><span class=\"p\">(</span>\n
    \               <span class=\"s2\">&quot;a&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">alt</span><span class=\"o\">=</span><span class=\"s2\">&quot;id&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">title</span><span class=\"o\">=</span><span
    class=\"sa\">f</span><span class=\"s2\">&quot;link to #</span><span class=\"si\">{</span><span
    class=\"nb\">id</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">href</span><span class=\"o\">=</span><span
    class=\"sa\">f</span><span class=\"s2\">&quot;#</span><span class=\"si\">{</span><span
    class=\"nb\">id</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"o\">**</span><span class=\"p\">{</span><span
    class=\"s2\">&quot;class&quot;</span><span class=\"p\">:</span> <span class=\"s2\">&quot;heading-permalink&quot;</span><span
    class=\"p\">},</span>\n            <span class=\"p\">)</span>\n            <span
    class=\"n\">span</span> <span class=\"o\">=</span> <span class=\"n\">soup</span><span
    class=\"o\">.</span><span class=\"n\">new_tag</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;span&quot;</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"p\">{</span><span class=\"s2\">&quot;class&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;visually-hidden&quot;</span><span class=\"p\">})</span>\n
    \           <span class=\"n\">svg</span> <span class=\"o\">=</span> <span class=\"n\">soup</span><span
    class=\"o\">.</span><span class=\"n\">new_tag</span><span class=\"p\">(</span>\n
    \               <span class=\"s2\">&quot;svg&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">fill</span><span class=\"o\">=</span><span class=\"s2\">&quot;currentColor&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">focusable</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;false&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">width</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;1em&quot;</span><span class=\"p\">,</span>\n                <span
    class=\"n\">height</span><span class=\"o\">=</span><span class=\"s2\">&quot;1em&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">xmlns</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;http://www.w3.org/2000/svg&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">viewBox</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;0 0 24 24&quot;</span><span class=\"p\">,</span>\n                <span
    class=\"o\">**</span><span class=\"p\">{</span>\n                    <span class=\"s2\">&quot;aria-hidden&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;true&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">},</span>\n            <span class=\"p\">)</span>\n\n
    \           <span class=\"n\">path</span> <span class=\"o\">=</span> <span class=\"n\">soup</span><span
    class=\"o\">.</span><span class=\"n\">new_tag</span><span class=\"p\">(</span>\n
    \               <span class=\"s2\">&quot;path&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">d</span><span class=\"o\">=</span><span class=\"s2\">&quot;M9.199
    13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
    5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
    6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
    3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
    3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
    13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
    2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0
    0 0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
    19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997
    0 0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
    1.563 1 1 0 0 0 1.602-1.198z&quot;</span><span class=\"p\">,</span>\n            <span
    class=\"p\">)</span>\n            <span class=\"n\">svg</span><span class=\"o\">.</span><span
    class=\"n\">append</span><span class=\"p\">(</span><span class=\"n\">path</span><span
    class=\"p\">)</span>\n            <span class=\"n\">link</span><span class=\"o\">.</span><span
    class=\"n\">append</span><span class=\"p\">(</span><span class=\"n\">span</span><span
    class=\"p\">)</span>\n            <span class=\"n\">link</span><span class=\"o\">.</span><span
    class=\"n\">append</span><span class=\"p\">(</span><span class=\"n\">svg</span><span
    class=\"p\">)</span>\n            <span class=\"n\">heading</span><span class=\"o\">.</span><span
    class=\"n\">append</span><span class=\"p\">(</span><span class=\"n\">link</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>    </div>\n
    \   <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Heading_Link.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Creates links next to all heading tags
    to make it easier for users to share a ! ! !\" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Heading_Link.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Creates links next to all heading tags
    to make it easier for users to share a ! ! !\" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<article style=\"text-align: center;\">\n    <style>\n        section
    {\n            font-size: 200%;\n        }\n\n\n        .edit {\n            display:
    none;\n        }\n    </style>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Heading_Link.Py\n    </h1>\n</section></article>\n     </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Heading_Link.Py\n    </h1>\n</section>    <section class=\"body\">\n        <p>Creates
    links next to all heading tags to make it easier for users to share a\nspecific
    heading.</p>\n<p>!! function <h2 id='post_render' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>post_render <em class='small'>function</em></h2>\nThis plugin creates a
    link svg next to all headings.</p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">post_render <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">post_render</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"n\">Markata</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            This plugin creates a link svg next to all headings.</span>\n<span
    class=\"sd\">            &quot;&quot;&quot;</span>\n\n            <span class=\"k\">with</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">cache</span>
    <span class=\"k\">as</span> <span class=\"n\">cache</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">for</span> <span class=\"n\">article</span>
    <span class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">iter_articles</span><span class=\"p\">(</span><span class=\"s2\">&quot;link
    headers&quot;</span><span class=\"p\">):</span>\n                    <span class=\"n\">key</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">make_hash</span><span class=\"p\">(</span>\n                        <span
    class=\"s2\">&quot;heading_link&quot;</span><span class=\"p\">,</span>\n                        <span
    class=\"s2\">&quot;post_render&quot;</span><span class=\"p\">,</span>\n                        <span
    class=\"n\">Path</span><span class=\"p\">(</span><span class=\"vm\">__file__</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">read_text</span><span
    class=\"p\">(),</span>\n                        <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">content</span><span class=\"p\">,</span>\n
    \                       <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">html</span><span class=\"p\">,</span>\n                    <span class=\"p\">)</span>\n\n
    \                   <span class=\"n\">html_from_cache</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">precache</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">)</span>\n\n                    <span
    class=\"k\">if</span> <span class=\"n\">html_from_cache</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">html</span> <span class=\"o\">=</span> <span class=\"n\">link_headings</span><span
    class=\"p\">(</span><span class=\"n\">article</span><span class=\"p\">)</span>\n
    \                       <span class=\"n\">cache</span><span class=\"o\">.</span><span
    class=\"n\">add</span><span class=\"p\">(</span>\n                            <span
    class=\"n\">key</span><span class=\"p\">,</span>\n                            <span
    class=\"n\">html</span><span class=\"p\">,</span>\n                            <span
    class=\"n\">expire</span><span class=\"o\">=</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">default_cache_expire</span><span class=\"p\">,</span>\n                        <span
    class=\"p\">)</span>\n                    <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">html</span> <span class=\"o\">=</span>
    <span class=\"n\">html_from_cache</span>\n                    <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">html</span> <span class=\"o\">=</span> <span
    class=\"n\">html</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='link_headings'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>link_headings <em
    class='small'>function</em></h2>\nUse BeautifulSoup to find all headings and run
    link_heading on them.</p>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">link_headings <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">link_headings</span><span class=\"p\">(</span><span class=\"n\">article</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            Use BeautifulSoup to find all headings and run link_heading
    on them.</span>\n<span class=\"sd\">            &quot;&quot;&quot;</span>\n            <span
    class=\"n\">soup</span> <span class=\"o\">=</span> <span class=\"n\">BeautifulSoup</span><span
    class=\"p\">(</span><span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">html</span><span class=\"p\">,</span> <span class=\"s2\">&quot;html.parser&quot;</span><span
    class=\"p\">)</span>\n            <span class=\"k\">for</span> <span class=\"n\">heading</span>
    <span class=\"ow\">in</span> <span class=\"n\">soup</span><span class=\"o\">.</span><span
    class=\"n\">find_all</span><span class=\"p\">(</span><span class=\"n\">re</span><span
    class=\"o\">.</span><span class=\"n\">compile</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;^h[1-6]$&quot;</span><span class=\"p\">)):</span>\n                <span
    class=\"k\">if</span> <span class=\"p\">(</span>\n                    <span class=\"ow\">not</span>
    <span class=\"n\">heading</span><span class=\"o\">.</span><span class=\"n\">find</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;a&quot;</span><span class=\"p\">,</span>
    <span class=\"p\">{</span><span class=\"s2\">&quot;class&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;heading-permalink&quot;</span><span class=\"p\">})</span>\n
    \                   <span class=\"ow\">and</span> <span class=\"n\">heading</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;id&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">!=</span> <span class=\"s2\">&quot;title&quot;</span>\n
    \               <span class=\"p\">):</span>\n                    <span class=\"n\">link_heading</span><span
    class=\"p\">(</span><span class=\"n\">soup</span><span class=\"p\">,</span> <span
    class=\"n\">heading</span><span class=\"p\">)</span>\n            <span class=\"k\">return</span>
    <span class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">soup</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='link_heading'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>link_heading <em
    class='small'>function</em></h2>\nMutate soup to include an svg link at the heading
    passed in.</p>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">link_heading <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">link_heading</span><span class=\"p\">(</span><span class=\"n\">soup</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;bs4.BeautifulSoup&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">heading</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;bs4.element.Tag&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            Mutate soup to include an svg link at the heading passed
    in.</span>\n<span class=\"sd\">            &quot;&quot;&quot;</span>\n            <span
    class=\"nb\">id</span> <span class=\"o\">=</span> <span class=\"n\">heading</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;id&quot;</span><span class=\"p\">)</span>\n\n            <span
    class=\"n\">link</span> <span class=\"o\">=</span> <span class=\"n\">soup</span><span
    class=\"o\">.</span><span class=\"n\">new_tag</span><span class=\"p\">(</span>\n
    \               <span class=\"s2\">&quot;a&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">alt</span><span class=\"o\">=</span><span class=\"s2\">&quot;id&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">title</span><span class=\"o\">=</span><span
    class=\"sa\">f</span><span class=\"s2\">&quot;link to #</span><span class=\"si\">{</span><span
    class=\"nb\">id</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">href</span><span class=\"o\">=</span><span
    class=\"sa\">f</span><span class=\"s2\">&quot;#</span><span class=\"si\">{</span><span
    class=\"nb\">id</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"o\">**</span><span class=\"p\">{</span><span
    class=\"s2\">&quot;class&quot;</span><span class=\"p\">:</span> <span class=\"s2\">&quot;heading-permalink&quot;</span><span
    class=\"p\">},</span>\n            <span class=\"p\">)</span>\n            <span
    class=\"n\">span</span> <span class=\"o\">=</span> <span class=\"n\">soup</span><span
    class=\"o\">.</span><span class=\"n\">new_tag</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;span&quot;</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"p\">{</span><span class=\"s2\">&quot;class&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;visually-hidden&quot;</span><span class=\"p\">})</span>\n
    \           <span class=\"n\">svg</span> <span class=\"o\">=</span> <span class=\"n\">soup</span><span
    class=\"o\">.</span><span class=\"n\">new_tag</span><span class=\"p\">(</span>\n
    \               <span class=\"s2\">&quot;svg&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">fill</span><span class=\"o\">=</span><span class=\"s2\">&quot;currentColor&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">focusable</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;false&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">width</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;1em&quot;</span><span class=\"p\">,</span>\n                <span
    class=\"n\">height</span><span class=\"o\">=</span><span class=\"s2\">&quot;1em&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">xmlns</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;http://www.w3.org/2000/svg&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">viewBox</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;0 0 24 24&quot;</span><span class=\"p\">,</span>\n                <span
    class=\"o\">**</span><span class=\"p\">{</span>\n                    <span class=\"s2\">&quot;aria-hidden&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;true&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">},</span>\n            <span class=\"p\">)</span>\n\n
    \           <span class=\"n\">path</span> <span class=\"o\">=</span> <span class=\"n\">soup</span><span
    class=\"o\">.</span><span class=\"n\">new_tag</span><span class=\"p\">(</span>\n
    \               <span class=\"s2\">&quot;path&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">d</span><span class=\"o\">=</span><span class=\"s2\">&quot;M9.199
    13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
    5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
    6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
    3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
    3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
    13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
    2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0
    0 0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
    19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997
    0 0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
    1.563 1 1 0 0 0 1.602-1.198z&quot;</span><span class=\"p\">,</span>\n            <span
    class=\"p\">)</span>\n            <span class=\"n\">svg</span><span class=\"o\">.</span><span
    class=\"n\">append</span><span class=\"p\">(</span><span class=\"n\">path</span><span
    class=\"p\">)</span>\n            <span class=\"n\">link</span><span class=\"o\">.</span><span
    class=\"n\">append</span><span class=\"p\">(</span><span class=\"n\">span</span><span
    class=\"p\">)</span>\n            <span class=\"n\">link</span><span class=\"o\">.</span><span
    class=\"n\">append</span><span class=\"p\">(</span><span class=\"n\">svg</span><span
    class=\"p\">)</span>\n            <span class=\"n\">heading</span><span class=\"o\">.</span><span
    class=\"n\">append</span><span class=\"p\">(</span><span class=\"n\">link</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>"
  raw.md: ''
published: true
slug: markata/plugins/heading-link
title: Heading_Link.Py


---

Creates links next to all heading tags to make it easier for users to share a
specific heading.


!! function <h2 id='post_render' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>post_render <em class='small'>function</em></h2>
    This plugin creates a link svg next to all headings.
???+ source "post_render <em class='small'>source</em>"

```python

        def post_render(markata: Markata) -> None:
            """
            This plugin creates a link svg next to all headings.
            """

            with markata.cache as cache:
                for article in markata.iter_articles("link headers"):
                    key = markata.make_hash(
                        "heading_link",
                        "post_render",
                        Path(__file__).read_text(),
                        article.content,
                        article.html,
                    )

                    html_from_cache = markata.precache.get(key)

                    if html_from_cache is None:
                        html = link_headings(article)
                        cache.add(
                            key,
                            html,
                            expire=markata.config.default_cache_expire,
                        )
                    else:
                        html = html_from_cache
                    article.html = html
```


!! function <h2 id='link_headings' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>link_headings <em class='small'>function</em></h2>
    Use BeautifulSoup to find all headings and run link_heading on them.
???+ source "link_headings <em class='small'>source</em>"

```python

        def link_headings(article: "Post") -> str:
            """
            Use BeautifulSoup to find all headings and run link_heading on them.
            """
            soup = BeautifulSoup(article.html, "html.parser")
            for heading in soup.find_all(re.compile("^h[1-6]$")):
                if (
                    not heading.find("a", {"class": "heading-permalink"})
                    and heading.get("id", "") != "title"
                ):
                    link_heading(soup, heading)
            return str(soup)
```


!! function <h2 id='link_heading' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>link_heading <em class='small'>function</em></h2>
    Mutate soup to include an svg link at the heading passed in.
???+ source "link_heading <em class='small'>source</em>"

```python

        def link_heading(soup: "bs4.BeautifulSoup", heading: "bs4.element.Tag") -> None:
            """
            Mutate soup to include an svg link at the heading passed in.
            """
            id = heading.get("id")

            link = soup.new_tag(
                "a",
                alt="id",
                title=f"link to #{id}",
                href=f"#{id}",
                **{"class": "heading-permalink"},
            )
            span = soup.new_tag("span", **{"class": "visually-hidden"})
            svg = soup.new_tag(
                "svg",
                fill="currentColor",
                focusable="false",
                width="1em",
                height="1em",
                xmlns="http://www.w3.org/2000/svg",
                viewBox="0 0 24 24",
                **{
                    "aria-hidden": "true",
                },
            )

            path = soup.new_tag(
                "path",
                d="M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992 5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242 6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985 3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005 3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201 13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995 2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0 0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836 19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0 0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632 1.563 1 1 0 0 0 1.602-1.198z",
            )
            svg.append(path)
            link.append(span)
            link.append(svg)
            heading.append(link)
```

