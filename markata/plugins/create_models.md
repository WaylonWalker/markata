---
content: "None\n\n\n!! class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>Config <em class='small'>class</em></h2>\n\n???+ source \"Config <em class='small'>source</em>\"\n\n```python\n\n
  \       class Config:\n            env_prefix = \"markata_\"\n            extras
  = \"allow\"\n```\n\n\n!! class <h2 id='PostConfig' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>PostConfig <em class='small'>class</em></h2>\n\n???+ source \"PostConfig
  <em class='small'>source</em>\"\n\n```python\n\n        class PostConfig:\n            title
  = \"Markata.Post\"\n            arbitrary_types_allowed = True\n            copy_on_model_validation
  = False\n```\n\n\n!! function <h2 id='create_models' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>create_models <em class='small'>function</em></h2>\n\n???+ source \"create_models
  <em class='small'>source</em>\"\n\n```python\n\n        def create_models(markata:
  \"Markata\") -> None:\n            post_models = tuple(unique_everseen(markata.post_models))\n
  \           markata.Post = create_model(\n                \"Post\",\n                __base__=post_models,\n
  \           )\n            markata.Posts = create_model(\n                \"Posts\",\n
  \               posts=(List[markata.Post], ...),\n            )\n            markata.Post.markata
  = markata\n            markata.Config = create_model(\n                \"Config\",\n
  \               __base__=tuple(unique_everseen(markata.config_models)),\n            )\n```\n\n"
date: 0001-01-01
description: 'None ! ???+ source  ! ???+ source  ! ???+ source '
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Create_Models.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None ! ???+ source  ! ???+ source  !
    ???+ source \" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Create_Models.Py</title>\n<meta
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Create_Models.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <p>None</p>\n<p>!!
    class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Config
    <em class='small'>class</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">Config <em class='small'>source</em></p>\n</div>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
    <span class=\"nc\">Config</span><span class=\"p\">:</span>\n            <span
    class=\"n\">env_prefix</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;markata_&quot;</span>\n
    \           <span class=\"n\">extras</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;allow&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='PostConfig' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>PostConfig <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">PostConfig
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
    <span class=\"nc\">PostConfig</span><span class=\"p\">:</span>\n            <span
    class=\"n\">title</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;Markata.Post&quot;</span>\n
    \           <span class=\"n\">arbitrary_types_allowed</span> <span class=\"o\">=</span>
    <span class=\"kc\">True</span>\n            <span class=\"n\">copy_on_model_validation</span>
    <span class=\"o\">=</span> <span class=\"kc\">False</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='create_models' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>create_models <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">create_models
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
    <span class=\"nf\">create_models</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">post_models</span> <span class=\"o\">=</span> <span
    class=\"nb\">tuple</span><span class=\"p\">(</span><span class=\"n\">unique_everseen</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">post_models</span><span class=\"p\">))</span>\n            <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">Post</span> <span class=\"o\">=</span> <span
    class=\"n\">create_model</span><span class=\"p\">(</span>\n                <span
    class=\"s2\">&quot;Post&quot;</span><span class=\"p\">,</span>\n                <span
    class=\"n\">__base__</span><span class=\"o\">=</span><span class=\"n\">post_models</span><span
    class=\"p\">,</span>\n            <span class=\"p\">)</span>\n            <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">Posts</span>
    <span class=\"o\">=</span> <span class=\"n\">create_model</span><span class=\"p\">(</span>\n
    \               <span class=\"s2\">&quot;Posts&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">posts</span><span class=\"o\">=</span><span
    class=\"p\">(</span><span class=\"n\">List</span><span class=\"p\">[</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">Post</span><span
    class=\"p\">],</span> <span class=\"o\">...</span><span class=\"p\">),</span>\n
    \           <span class=\"p\">)</span>\n            <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">Post</span><span class=\"o\">.</span><span
    class=\"n\">markata</span> <span class=\"o\">=</span> <span class=\"n\">markata</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">Config</span>
    <span class=\"o\">=</span> <span class=\"n\">create_model</span><span class=\"p\">(</span>\n
    \               <span class=\"s2\">&quot;Config&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">__base__</span><span class=\"o\">=</span><span
    class=\"nb\">tuple</span><span class=\"p\">(</span><span class=\"n\">unique_everseen</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config_models</span><span class=\"p\">)),</span>\n            <span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>    </div>\n
    \   <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Create_Models.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None ! ???+ source  ! ???+ source  !
    ???+ source \" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Create_Models.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None ! ???+ source  ! ???+ source  !
    ???+ source \" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    </head>\n    <body>\n<article style=\"text-align:
    center;\">\n    <style>\n        section {\n            font-size: 200%;\n        }\n\n\n
    \       .edit {\n            display: none;\n        }\n    </style>\n<section
    class=\"title\">\n    <h1 id=\"title\">\n        Create_Models.Py\n    </h1>\n</section></article>\n
    \    </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Create_Models.Py\n    </h1>\n</section>    <section class=\"body\">\n
    \       <p>None</p>\n<p>!! class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>Config <em class='small'>class</em></h2></p>\n<div class=\"admonition source
    is-collapsible collapsible-open\">\n<p class=\"admonition-title\">Config <em class='small'>source</em></p>\n</div>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
    <span class=\"nc\">Config</span><span class=\"p\">:</span>\n            <span
    class=\"n\">env_prefix</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;markata_&quot;</span>\n
    \           <span class=\"n\">extras</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;allow&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='PostConfig' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>PostConfig <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">PostConfig
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
    <span class=\"nc\">PostConfig</span><span class=\"p\">:</span>\n            <span
    class=\"n\">title</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;Markata.Post&quot;</span>\n
    \           <span class=\"n\">arbitrary_types_allowed</span> <span class=\"o\">=</span>
    <span class=\"kc\">True</span>\n            <span class=\"n\">copy_on_model_validation</span>
    <span class=\"o\">=</span> <span class=\"kc\">False</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='create_models' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>create_models <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">create_models
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
    <span class=\"nf\">create_models</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">post_models</span> <span class=\"o\">=</span> <span
    class=\"nb\">tuple</span><span class=\"p\">(</span><span class=\"n\">unique_everseen</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">post_models</span><span class=\"p\">))</span>\n            <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">Post</span> <span class=\"o\">=</span> <span
    class=\"n\">create_model</span><span class=\"p\">(</span>\n                <span
    class=\"s2\">&quot;Post&quot;</span><span class=\"p\">,</span>\n                <span
    class=\"n\">__base__</span><span class=\"o\">=</span><span class=\"n\">post_models</span><span
    class=\"p\">,</span>\n            <span class=\"p\">)</span>\n            <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">Posts</span>
    <span class=\"o\">=</span> <span class=\"n\">create_model</span><span class=\"p\">(</span>\n
    \               <span class=\"s2\">&quot;Posts&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">posts</span><span class=\"o\">=</span><span
    class=\"p\">(</span><span class=\"n\">List</span><span class=\"p\">[</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">Post</span><span
    class=\"p\">],</span> <span class=\"o\">...</span><span class=\"p\">),</span>\n
    \           <span class=\"p\">)</span>\n            <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">Post</span><span class=\"o\">.</span><span
    class=\"n\">markata</span> <span class=\"o\">=</span> <span class=\"n\">markata</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">Config</span>
    <span class=\"o\">=</span> <span class=\"n\">create_model</span><span class=\"p\">(</span>\n
    \               <span class=\"s2\">&quot;Config&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">__base__</span><span class=\"o\">=</span><span
    class=\"nb\">tuple</span><span class=\"p\">(</span><span class=\"n\">unique_everseen</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config_models</span><span class=\"p\">)),</span>\n            <span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>"
  raw.md: ''
published: true
slug: markata/plugins/create-models
title: Create_Models.Py


---

None


!! class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Config <em class='small'>class</em></h2>

???+ source "Config <em class='small'>source</em>"

```python

        class Config:
            env_prefix = "markata_"
            extras = "allow"
```


!! class <h2 id='PostConfig' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>PostConfig <em class='small'>class</em></h2>

???+ source "PostConfig <em class='small'>source</em>"

```python

        class PostConfig:
            title = "Markata.Post"
            arbitrary_types_allowed = True
            copy_on_model_validation = False
```


!! function <h2 id='create_models' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>create_models <em class='small'>function</em></h2>

???+ source "create_models <em class='small'>source</em>"

```python

        def create_models(markata: "Markata") -> None:
            post_models = tuple(unique_everseen(markata.post_models))
            markata.Post = create_model(
                "Post",
                __base__=post_models,
            )
            markata.Posts = create_model(
                "Posts",
                posts=(List[markata.Post], ...),
            )
            markata.Post.markata = markata
            markata.Config = create_model(
                "Config",
                __base__=tuple(unique_everseen(markata.config_models)),
            )
```

