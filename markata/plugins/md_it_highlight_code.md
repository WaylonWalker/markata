---
content: "None\n\n\n!! function <h2 id='highlight_code' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>highlight_code <em class='small'>function</em></h2>\n    Code highlighter
  for markdown-it-py.\n???+ source \"highlight_code <em class='small'>source</em>\"\n\n```python\n\n
  \       def highlight_code(code, name, attrs, markata=None):\n            \"\"\"Code
  highlighter for markdown-it-py.\"\"\"\n\n            try:\n                lexer
  = get_lexer_by_name(name or \"text\")\n            except ClassNotFound:\n                lexer
  = get_lexer_by_name(\"text\")\n\n            import re\n\n            pattern =
  r'(\\w+)\\s*=\\s*(\".*?\"|\\S+)'\n            matches = re.findall(pattern, attrs)\n
  \           attrs = dict(matches)\n\n            if attrs.get(\"hl_lines\"):\n                formatter
  = HtmlFormatter(hl_lines=attrs.get(\"hl_lines\"))\n            else:\n                formatter
  = HtmlFormatter()\n\n            copy_button = f\"\"\"<button class='copy' title='copy
  code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\">{COPY_ICON}</button>\"\"\"\n\n
  \           from markdown_it import MarkdownIt\n\n            md = MarkdownIt(\n
  \               \"commonmark\",\n                {\n                    \"html\":
  True,\n                    \"typographer\": True,\n                },\n            )\n\n
  \           if attrs.get(\"help\"):\n                help = f\"\"\"\n                <a
  href={attrs.get('help').strip('<').strip('>').strip('\"').strip(\"'\")} title='help
  link' class='help'>{HELP_ICON}</a>\n                \"\"\"\n            else:\n
  \               help = \"\"\n            if attrs.get(\"title\"):\n                file
  = f\"\"\"\n        <div class='filepath'>\n        {md.render(attrs.get('title').strip('\"').strip(\"'\"))}\n
  \       <div class='right'>\n        {help}\n        {copy_button}\n        </div>\n
  \       </div>\n        \"\"\"\n            else:\n                file = f\"\"\"\n
  \       <div class='copy-wrapper'>\n        {help}\n        {copy_button}\n        </div>\n
  \               \"\"\"\n            return f\"\"\"<pre class='wrapper'>\n        {file}\n
  \       {highlight(code, lexer, formatter)}\n        </pre>\n        \"\"\"\n```\n\n"
date: 0001-01-01
description: None !
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Md_It_Highlight_Code.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None !\" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Md_It_Highlight_Code.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None !\" />\n <link href=\"/favicon.ico\"
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Md_It_Highlight_Code.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <p>None</p>\n<p>!!
    function <h2 id='highlight_code' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>highlight_code <em class='small'>function</em></h2>\nCode highlighter for
    markdown-it-py.</p>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">highlight_code <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">highlight_code</span><span class=\"p\">(</span><span class=\"n\">code</span><span
    class=\"p\">,</span> <span class=\"n\">name</span><span class=\"p\">,</span> <span
    class=\"n\">attrs</span><span class=\"p\">,</span> <span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"kc\">None</span><span class=\"p\">):</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;Code highlighter
    for markdown-it-py.&quot;&quot;&quot;</span>\n\n            <span class=\"k\">try</span><span
    class=\"p\">:</span>\n                <span class=\"n\">lexer</span> <span class=\"o\">=</span>
    <span class=\"n\">get_lexer_by_name</span><span class=\"p\">(</span><span class=\"n\">name</span>
    <span class=\"ow\">or</span> <span class=\"s2\">&quot;text&quot;</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">except</span> <span class=\"n\">ClassNotFound</span><span
    class=\"p\">:</span>\n                <span class=\"n\">lexer</span> <span class=\"o\">=</span>
    <span class=\"n\">get_lexer_by_name</span><span class=\"p\">(</span><span class=\"s2\">&quot;text&quot;</span><span
    class=\"p\">)</span>\n\n            <span class=\"kn\">import</span> <span class=\"nn\">re</span>\n\n
    \           <span class=\"n\">pattern</span> <span class=\"o\">=</span> <span
    class=\"sa\">r</span><span class=\"s1\">&#39;(\\w+)\\s*=\\s*(&quot;.*?&quot;|\\S+)&#39;</span>\n
    \           <span class=\"n\">matches</span> <span class=\"o\">=</span> <span
    class=\"n\">re</span><span class=\"o\">.</span><span class=\"n\">findall</span><span
    class=\"p\">(</span><span class=\"n\">pattern</span><span class=\"p\">,</span>
    <span class=\"n\">attrs</span><span class=\"p\">)</span>\n            <span class=\"n\">attrs</span>
    <span class=\"o\">=</span> <span class=\"nb\">dict</span><span class=\"p\">(</span><span
    class=\"n\">matches</span><span class=\"p\">)</span>\n\n            <span class=\"k\">if</span>
    <span class=\"n\">attrs</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;hl_lines&quot;</span><span class=\"p\">):</span>\n
    \               <span class=\"n\">formatter</span> <span class=\"o\">=</span>
    <span class=\"n\">HtmlFormatter</span><span class=\"p\">(</span><span class=\"n\">hl_lines</span><span
    class=\"o\">=</span><span class=\"n\">attrs</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;hl_lines&quot;</span><span
    class=\"p\">))</span>\n            <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">formatter</span> <span class=\"o\">=</span>
    <span class=\"n\">HtmlFormatter</span><span class=\"p\">()</span>\n\n            <span
    class=\"n\">copy_button</span> <span class=\"o\">=</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;&quot;&quot;&lt;button class=&#39;copy&#39; title=&#39;copy
    code to clipboard&#39; onclick=&quot;navigator.clipboard.writeText(this.parentElement.parentElement.querySelector(&#39;pre&#39;).textContent)&quot;&gt;</span><span
    class=\"si\">{</span><span class=\"n\">COPY_ICON</span><span class=\"si\">}</span><span
    class=\"s2\">&lt;/button&gt;&quot;&quot;&quot;</span>\n\n            <span class=\"kn\">from</span>
    <span class=\"nn\">markdown_it</span> <span class=\"kn\">import</span> <span class=\"n\">MarkdownIt</span>\n\n
    \           <span class=\"n\">md</span> <span class=\"o\">=</span> <span class=\"n\">MarkdownIt</span><span
    class=\"p\">(</span>\n                <span class=\"s2\">&quot;commonmark&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"p\">{</span>\n                    <span
    class=\"s2\">&quot;html&quot;</span><span class=\"p\">:</span> <span class=\"kc\">True</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;typographer&quot;</span><span
    class=\"p\">:</span> <span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">},</span>\n            <span class=\"p\">)</span>\n\n
    \           <span class=\"k\">if</span> <span class=\"n\">attrs</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;help&quot;</span><span
    class=\"p\">):</span>\n                <span class=\"n\">help</span> <span class=\"o\">=</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;&quot;&quot;</span>\n<span
    class=\"s2\">                &lt;a href=</span><span class=\"si\">{</span><span
    class=\"n\">attrs</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s1\">&#39;help&#39;</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">strip</span><span class=\"p\">(</span><span
    class=\"s1\">&#39;&lt;&#39;</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">strip</span><span class=\"p\">(</span><span class=\"s1\">&#39;&gt;&#39;</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">strip</span><span
    class=\"p\">(</span><span class=\"s1\">&#39;&quot;&#39;</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">strip</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;&#39;&quot;</span><span class=\"p\">)</span><span class=\"si\">}</span><span
    class=\"s2\"> title=&#39;help link&#39; class=&#39;help&#39;&gt;</span><span class=\"si\">{</span><span
    class=\"n\">HELP_ICON</span><span class=\"si\">}</span><span class=\"s2\">&lt;/a&gt;</span>\n<span
    class=\"s2\">                &quot;&quot;&quot;</span>\n            <span class=\"k\">else</span><span
    class=\"p\">:</span>\n                <span class=\"n\">help</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;&quot;</span>\n            <span class=\"k\">if</span>
    <span class=\"n\">attrs</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;title&quot;</span><span class=\"p\">):</span>\n
    \               <span class=\"n\">file</span> <span class=\"o\">=</span> <span
    class=\"sa\">f</span><span class=\"s2\">&quot;&quot;&quot;</span>\n<span class=\"s2\">
    \       &lt;div class=&#39;filepath&#39;&gt;</span>\n<span class=\"s2\">        </span><span
    class=\"si\">{</span><span class=\"n\">md</span><span class=\"o\">.</span><span
    class=\"n\">render</span><span class=\"p\">(</span><span class=\"n\">attrs</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s1\">&#39;title&#39;</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">strip</span><span class=\"p\">(</span><span class=\"s1\">&#39;&quot;&#39;</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">strip</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;&#39;&quot;</span><span class=\"p\">))</span><span
    class=\"si\">}</span>\n<span class=\"s2\">        &lt;div class=&#39;right&#39;&gt;</span>\n<span
    class=\"s2\">        </span><span class=\"si\">{</span><span class=\"n\">help</span><span
    class=\"si\">}</span>\n<span class=\"s2\">        </span><span class=\"si\">{</span><span
    class=\"n\">copy_button</span><span class=\"si\">}</span>\n<span class=\"s2\">
    \       &lt;/div&gt;</span>\n<span class=\"s2\">        &lt;/div&gt;</span>\n<span
    class=\"s2\">        &quot;&quot;&quot;</span>\n            <span class=\"k\">else</span><span
    class=\"p\">:</span>\n                <span class=\"n\">file</span> <span class=\"o\">=</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;&quot;&quot;</span>\n<span
    class=\"s2\">        &lt;div class=&#39;copy-wrapper&#39;&gt;</span>\n<span class=\"s2\">
    \       </span><span class=\"si\">{</span><span class=\"n\">help</span><span class=\"si\">}</span>\n<span
    class=\"s2\">        </span><span class=\"si\">{</span><span class=\"n\">copy_button</span><span
    class=\"si\">}</span>\n<span class=\"s2\">        &lt;/div&gt;</span>\n<span class=\"s2\">
    \               &quot;&quot;&quot;</span>\n            <span class=\"k\">return</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;&quot;&quot;&lt;pre class=&#39;wrapper&#39;&gt;</span>\n<span
    class=\"s2\">        </span><span class=\"si\">{</span><span class=\"n\">file</span><span
    class=\"si\">}</span>\n<span class=\"s2\">        </span><span class=\"si\">{</span><span
    class=\"n\">highlight</span><span class=\"p\">(</span><span class=\"n\">code</span><span
    class=\"p\">,</span><span class=\"w\"> </span><span class=\"n\">lexer</span><span
    class=\"p\">,</span><span class=\"w\"> </span><span class=\"n\">formatter</span><span
    class=\"p\">)</span><span class=\"si\">}</span>\n<span class=\"s2\">        &lt;/pre&gt;</span>\n<span
    class=\"s2\">        &quot;&quot;&quot;</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>
    \   </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Md_It_Highlight_Code.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None !\" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Md_It_Highlight_Code.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None !\" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<article style=\"text-align: center;\">\n    <style>\n        section
    {\n            font-size: 200%;\n        }\n\n\n        .edit {\n            display:
    none;\n        }\n    </style>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Md_It_Highlight_Code.Py\n    </h1>\n</section></article>\n     </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Md_It_Highlight_Code.Py\n    </h1>\n</section>    <section class=\"body\">\n
    \       <p>None</p>\n<p>!! function <h2 id='highlight_code' class='admonition-title'
    style='margin:0;padding:.5rem 1rem;'>highlight_code <em class='small'>function</em></h2>\nCode
    highlighter for markdown-it-py.</p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">highlight_code <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">highlight_code</span><span class=\"p\">(</span><span class=\"n\">code</span><span
    class=\"p\">,</span> <span class=\"n\">name</span><span class=\"p\">,</span> <span
    class=\"n\">attrs</span><span class=\"p\">,</span> <span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"kc\">None</span><span class=\"p\">):</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;Code highlighter
    for markdown-it-py.&quot;&quot;&quot;</span>\n\n            <span class=\"k\">try</span><span
    class=\"p\">:</span>\n                <span class=\"n\">lexer</span> <span class=\"o\">=</span>
    <span class=\"n\">get_lexer_by_name</span><span class=\"p\">(</span><span class=\"n\">name</span>
    <span class=\"ow\">or</span> <span class=\"s2\">&quot;text&quot;</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">except</span> <span class=\"n\">ClassNotFound</span><span
    class=\"p\">:</span>\n                <span class=\"n\">lexer</span> <span class=\"o\">=</span>
    <span class=\"n\">get_lexer_by_name</span><span class=\"p\">(</span><span class=\"s2\">&quot;text&quot;</span><span
    class=\"p\">)</span>\n\n            <span class=\"kn\">import</span> <span class=\"nn\">re</span>\n\n
    \           <span class=\"n\">pattern</span> <span class=\"o\">=</span> <span
    class=\"sa\">r</span><span class=\"s1\">&#39;(\\w+)\\s*=\\s*(&quot;.*?&quot;|\\S+)&#39;</span>\n
    \           <span class=\"n\">matches</span> <span class=\"o\">=</span> <span
    class=\"n\">re</span><span class=\"o\">.</span><span class=\"n\">findall</span><span
    class=\"p\">(</span><span class=\"n\">pattern</span><span class=\"p\">,</span>
    <span class=\"n\">attrs</span><span class=\"p\">)</span>\n            <span class=\"n\">attrs</span>
    <span class=\"o\">=</span> <span class=\"nb\">dict</span><span class=\"p\">(</span><span
    class=\"n\">matches</span><span class=\"p\">)</span>\n\n            <span class=\"k\">if</span>
    <span class=\"n\">attrs</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;hl_lines&quot;</span><span class=\"p\">):</span>\n
    \               <span class=\"n\">formatter</span> <span class=\"o\">=</span>
    <span class=\"n\">HtmlFormatter</span><span class=\"p\">(</span><span class=\"n\">hl_lines</span><span
    class=\"o\">=</span><span class=\"n\">attrs</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;hl_lines&quot;</span><span
    class=\"p\">))</span>\n            <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">formatter</span> <span class=\"o\">=</span>
    <span class=\"n\">HtmlFormatter</span><span class=\"p\">()</span>\n\n            <span
    class=\"n\">copy_button</span> <span class=\"o\">=</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;&quot;&quot;&lt;button class=&#39;copy&#39; title=&#39;copy
    code to clipboard&#39; onclick=&quot;navigator.clipboard.writeText(this.parentElement.parentElement.querySelector(&#39;pre&#39;).textContent)&quot;&gt;</span><span
    class=\"si\">{</span><span class=\"n\">COPY_ICON</span><span class=\"si\">}</span><span
    class=\"s2\">&lt;/button&gt;&quot;&quot;&quot;</span>\n\n            <span class=\"kn\">from</span>
    <span class=\"nn\">markdown_it</span> <span class=\"kn\">import</span> <span class=\"n\">MarkdownIt</span>\n\n
    \           <span class=\"n\">md</span> <span class=\"o\">=</span> <span class=\"n\">MarkdownIt</span><span
    class=\"p\">(</span>\n                <span class=\"s2\">&quot;commonmark&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"p\">{</span>\n                    <span
    class=\"s2\">&quot;html&quot;</span><span class=\"p\">:</span> <span class=\"kc\">True</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;typographer&quot;</span><span
    class=\"p\">:</span> <span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">},</span>\n            <span class=\"p\">)</span>\n\n
    \           <span class=\"k\">if</span> <span class=\"n\">attrs</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;help&quot;</span><span
    class=\"p\">):</span>\n                <span class=\"n\">help</span> <span class=\"o\">=</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;&quot;&quot;</span>\n<span
    class=\"s2\">                &lt;a href=</span><span class=\"si\">{</span><span
    class=\"n\">attrs</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s1\">&#39;help&#39;</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">strip</span><span class=\"p\">(</span><span
    class=\"s1\">&#39;&lt;&#39;</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">strip</span><span class=\"p\">(</span><span class=\"s1\">&#39;&gt;&#39;</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">strip</span><span
    class=\"p\">(</span><span class=\"s1\">&#39;&quot;&#39;</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">strip</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;&#39;&quot;</span><span class=\"p\">)</span><span class=\"si\">}</span><span
    class=\"s2\"> title=&#39;help link&#39; class=&#39;help&#39;&gt;</span><span class=\"si\">{</span><span
    class=\"n\">HELP_ICON</span><span class=\"si\">}</span><span class=\"s2\">&lt;/a&gt;</span>\n<span
    class=\"s2\">                &quot;&quot;&quot;</span>\n            <span class=\"k\">else</span><span
    class=\"p\">:</span>\n                <span class=\"n\">help</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;&quot;</span>\n            <span class=\"k\">if</span>
    <span class=\"n\">attrs</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;title&quot;</span><span class=\"p\">):</span>\n
    \               <span class=\"n\">file</span> <span class=\"o\">=</span> <span
    class=\"sa\">f</span><span class=\"s2\">&quot;&quot;&quot;</span>\n<span class=\"s2\">
    \       &lt;div class=&#39;filepath&#39;&gt;</span>\n<span class=\"s2\">        </span><span
    class=\"si\">{</span><span class=\"n\">md</span><span class=\"o\">.</span><span
    class=\"n\">render</span><span class=\"p\">(</span><span class=\"n\">attrs</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s1\">&#39;title&#39;</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">strip</span><span class=\"p\">(</span><span class=\"s1\">&#39;&quot;&#39;</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">strip</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;&#39;&quot;</span><span class=\"p\">))</span><span
    class=\"si\">}</span>\n<span class=\"s2\">        &lt;div class=&#39;right&#39;&gt;</span>\n<span
    class=\"s2\">        </span><span class=\"si\">{</span><span class=\"n\">help</span><span
    class=\"si\">}</span>\n<span class=\"s2\">        </span><span class=\"si\">{</span><span
    class=\"n\">copy_button</span><span class=\"si\">}</span>\n<span class=\"s2\">
    \       &lt;/div&gt;</span>\n<span class=\"s2\">        &lt;/div&gt;</span>\n<span
    class=\"s2\">        &quot;&quot;&quot;</span>\n            <span class=\"k\">else</span><span
    class=\"p\">:</span>\n                <span class=\"n\">file</span> <span class=\"o\">=</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;&quot;&quot;</span>\n<span
    class=\"s2\">        &lt;div class=&#39;copy-wrapper&#39;&gt;</span>\n<span class=\"s2\">
    \       </span><span class=\"si\">{</span><span class=\"n\">help</span><span class=\"si\">}</span>\n<span
    class=\"s2\">        </span><span class=\"si\">{</span><span class=\"n\">copy_button</span><span
    class=\"si\">}</span>\n<span class=\"s2\">        &lt;/div&gt;</span>\n<span class=\"s2\">
    \               &quot;&quot;&quot;</span>\n            <span class=\"k\">return</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;&quot;&quot;&lt;pre class=&#39;wrapper&#39;&gt;</span>\n<span
    class=\"s2\">        </span><span class=\"si\">{</span><span class=\"n\">file</span><span
    class=\"si\">}</span>\n<span class=\"s2\">        </span><span class=\"si\">{</span><span
    class=\"n\">highlight</span><span class=\"p\">(</span><span class=\"n\">code</span><span
    class=\"p\">,</span><span class=\"w\"> </span><span class=\"n\">lexer</span><span
    class=\"p\">,</span><span class=\"w\"> </span><span class=\"n\">formatter</span><span
    class=\"p\">)</span><span class=\"si\">}</span>\n<span class=\"s2\">        &lt;/pre&gt;</span>\n<span
    class=\"s2\">        &quot;&quot;&quot;</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>"
  raw.md: ''
published: true
slug: markata/plugins/md-it-highlight-code
title: Md_It_Highlight_Code.Py


---

None


!! function <h2 id='highlight_code' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>highlight_code <em class='small'>function</em></h2>
    Code highlighter for markdown-it-py.
???+ source "highlight_code <em class='small'>source</em>"

```python

        def highlight_code(code, name, attrs, markata=None):
            """Code highlighter for markdown-it-py."""

            try:
                lexer = get_lexer_by_name(name or "text")
            except ClassNotFound:
                lexer = get_lexer_by_name("text")

            import re

            pattern = r'(\w+)\s*=\s*(".*?"|\S+)'
            matches = re.findall(pattern, attrs)
            attrs = dict(matches)

            if attrs.get("hl_lines"):
                formatter = HtmlFormatter(hl_lines=attrs.get("hl_lines"))
            else:
                formatter = HtmlFormatter()

            copy_button = f"""<button class='copy' title='copy code to clipboard' onclick="navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)">{COPY_ICON}</button>"""

            from markdown_it import MarkdownIt

            md = MarkdownIt(
                "commonmark",
                {
                    "html": True,
                    "typographer": True,
                },
            )

            if attrs.get("help"):
                help = f"""
                <a href={attrs.get('help').strip('<').strip('>').strip('"').strip("'")} title='help link' class='help'>{HELP_ICON}</a>
                """
            else:
                help = ""
            if attrs.get("title"):
                file = f"""
        <div class='filepath'>
        {md.render(attrs.get('title').strip('"').strip("'"))}
        <div class='right'>
        {help}
        {copy_button}
        </div>
        </div>
        """
            else:
                file = f"""
        <div class='copy-wrapper'>
        {help}
        {copy_button}
        </div>
                """
            return f"""<pre class='wrapper'>
        {file}
        {highlight(code, lexer, formatter)}
        </pre>
        """
```

