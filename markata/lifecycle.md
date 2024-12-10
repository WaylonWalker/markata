---
content: "The LifeCycle is a core component for the internal workings of Markata.
  \ It\nsets fourth the hooks available, the methods to run them on the Markata\ninstance,
  and the order they run in.\n\n### Usage\n\n``` python\nfrom markata import Lifecycle\n\nstep
  = Lifecycle.glob\n```\n\n\n!! class <h2 id='LifeCycle' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>LifeCycle <em class='small'>class</em></h2>\n
  \   LifeCycle currently supports the following steps.\n\n\n    * configure - load
  and fix configuration\n    * glob - find files\n    * load - load files\n    * validate_posts\n
  \   * pre_render - clean up files/metadata before render\n    * render - render
  content\n    * post_render - clean up rendered content\n    * save - store results
  to disk\n    * teardown - runs on exit\n???+ source \"LifeCycle <em class='small'>source</em>\"\n\n```python\n\n
  \       class LifeCycle(Enum):\n            \"\"\"\n            LifeCycle currently
  supports the following steps.\n\n\n            * configure - load and fix configuration\n
  \           * glob - find files\n            * load - load files\n            *
  validate_posts\n            * pre_render - clean up files/metadata before render\n
  \           * render - render content\n            * post_render - clean up rendered
  content\n            * save - store results to disk\n            * teardown - runs
  on exit\n\n            \"\"\"\n\n            config_model = auto()\n            post_model
  = auto()\n            create_models = auto()\n            load_config = auto()\n
  \           configure = auto()\n            validate_config = auto()\n            glob
  = auto()\n            load = auto()\n            pre_render = auto()\n            render
  = auto()\n            post_render = auto()\n            save = auto()\n            teardown
  = auto()\n\n            def __lt__(self, other: object) -> bool:\n                \"\"\"\n
  \               Determine whether other is less than this instance.\n                \"\"\"\n
  \               if isinstance(other, LifeCycle):\n                    return self.value
  < other.value\n                if isinstance(other, int):\n                    return
  self.value < other\n                return NotImplemented\n\n            def __eq__(self,
  other: object) -> bool:\n                \"\"\"\n                Determine whether
  other is equal to this instance.\n                \"\"\"\n                if isinstance(other,
  LifeCycle):\n                    return self.value == other.value\n                if
  isinstance(other, int):\n                    return self.value == other\n                return
  NotImplemented\n```\n\n\n!! method <h2 id='__lt__' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>__lt__ <em class='small'>method</em></h2>\n    Determine whether other is
  less than this instance.\n???+ source \"__lt__ <em class='small'>source</em>\"\n\n```python\n\n
  \       def __lt__(self, other: object) -> bool:\n                \"\"\"\n                Determine
  whether other is less than this instance.\n                \"\"\"\n                if
  isinstance(other, LifeCycle):\n                    return self.value < other.value\n
  \               if isinstance(other, int):\n                    return self.value
  < other\n                return NotImplemented\n```\n\n\n!! method <h2 id='__eq__'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__eq__ <em class='small'>method</em></h2>\n
  \   Determine whether other is equal to this instance.\n???+ source \"__eq__ <em
  class='small'>source</em>\"\n\n```python\n\n        def __eq__(self, other: object)
  -> bool:\n                \"\"\"\n                Determine whether other is equal
  to this instance.\n                \"\"\"\n                if isinstance(other,
  LifeCycle):\n                    return self.value == other.value\n                if
  isinstance(other, int):\n                    return self.value == other\n                return
  NotImplemented\n```\n\n"
date: 0001-01-01
description: The LifeCycle is a core component for the internal workings of Markata.  It
  ! ???+ source  ! !
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Lifecycle.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"The LifeCycle is a core component for
    the internal workings of Markata.  It ! ???+ source  ! !\" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Lifecycle.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"The LifeCycle is a core component for
    the internal workings of Markata.  It ! ???+ source  ! !\" />\n <link href=\"/favicon.ico\"
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Lifecycle.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <p>The LifeCycle is
    a core component for the internal workings of Markata.  It\nsets fourth the hooks
    available, the methods to run them on the Markata\ninstance, and the order they
    run in.</p>\n<h3>Usage</h3>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"kn\">from</span>
    <span class=\"nn\">markata</span> <span class=\"kn\">import</span> <span class=\"n\">Lifecycle</span>\n\n<span
    class=\"n\">step</span> <span class=\"o\">=</span> <span class=\"n\">Lifecycle</span><span
    class=\"o\">.</span><span class=\"n\">glob</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='LifeCycle' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>LifeCycle <em class='small'>class</em></h2>\nLifeCycle currently supports
    the following steps.</p>\n<pre><code>* configure - load and fix configuration\n*
    glob - find files\n* load - load files\n* validate_posts\n* pre_render - clean
    up files/metadata before render\n* render - render content\n* post_render - clean
    up rendered content\n* save - store results to disk\n* teardown - runs on exit\n</code></pre>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">LifeCycle
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
    <span class=\"nc\">LifeCycle</span><span class=\"p\">(</span><span class=\"n\">Enum</span><span
    class=\"p\">):</span>\n<span class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            LifeCycle currently supports the following steps.</span>\n\n\n<span
    class=\"sd\">            * configure - load and fix configuration</span>\n<span
    class=\"sd\">            * glob - find files</span>\n<span class=\"sd\">            *
    load - load files</span>\n<span class=\"sd\">            * validate_posts</span>\n<span
    class=\"sd\">            * pre_render - clean up files/metadata before render</span>\n<span
    class=\"sd\">            * render - render content</span>\n<span class=\"sd\">
    \           * post_render - clean up rendered content</span>\n<span class=\"sd\">
    \           * save - store results to disk</span>\n<span class=\"sd\">            *
    teardown - runs on exit</span>\n\n<span class=\"sd\">            &quot;&quot;&quot;</span>\n\n
    \           <span class=\"n\">config_model</span> <span class=\"o\">=</span> <span
    class=\"n\">auto</span><span class=\"p\">()</span>\n            <span class=\"n\">post_model</span>
    <span class=\"o\">=</span> <span class=\"n\">auto</span><span class=\"p\">()</span>\n
    \           <span class=\"n\">create_models</span> <span class=\"o\">=</span>
    <span class=\"n\">auto</span><span class=\"p\">()</span>\n            <span class=\"n\">load_config</span>
    <span class=\"o\">=</span> <span class=\"n\">auto</span><span class=\"p\">()</span>\n
    \           <span class=\"n\">configure</span> <span class=\"o\">=</span> <span
    class=\"n\">auto</span><span class=\"p\">()</span>\n            <span class=\"n\">validate_config</span>
    <span class=\"o\">=</span> <span class=\"n\">auto</span><span class=\"p\">()</span>\n
    \           <span class=\"n\">glob</span> <span class=\"o\">=</span> <span class=\"n\">auto</span><span
    class=\"p\">()</span>\n            <span class=\"n\">load</span> <span class=\"o\">=</span>
    <span class=\"n\">auto</span><span class=\"p\">()</span>\n            <span class=\"n\">pre_render</span>
    <span class=\"o\">=</span> <span class=\"n\">auto</span><span class=\"p\">()</span>\n
    \           <span class=\"n\">render</span> <span class=\"o\">=</span> <span class=\"n\">auto</span><span
    class=\"p\">()</span>\n            <span class=\"n\">post_render</span> <span
    class=\"o\">=</span> <span class=\"n\">auto</span><span class=\"p\">()</span>\n
    \           <span class=\"n\">save</span> <span class=\"o\">=</span> <span class=\"n\">auto</span><span
    class=\"p\">()</span>\n            <span class=\"n\">teardown</span> <span class=\"o\">=</span>
    <span class=\"n\">auto</span><span class=\"p\">()</span>\n\n            <span
    class=\"k\">def</span> <span class=\"fm\">__lt__</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"n\">other</span><span
    class=\"p\">:</span> <span class=\"nb\">object</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">bool</span><span class=\"p\">:</span>\n<span
    class=\"w\">                </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">                Determine whether other is less than this instance.</span>\n<span
    class=\"sd\">                &quot;&quot;&quot;</span>\n                <span
    class=\"k\">if</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">other</span><span class=\"p\">,</span> <span class=\"n\">LifeCycle</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">value</span>
    <span class=\"o\">&lt;</span> <span class=\"n\">other</span><span class=\"o\">.</span><span
    class=\"n\">value</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">other</span><span class=\"p\">,</span> <span
    class=\"nb\">int</span><span class=\"p\">):</span>\n                    <span
    class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">value</span> <span class=\"o\">&lt;</span> <span class=\"n\">other</span>\n
    \               <span class=\"k\">return</span> <span class=\"bp\">NotImplemented</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"fm\">__eq__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"n\">other</span><span class=\"p\">:</span> <span class=\"nb\">object</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">bool</span><span
    class=\"p\">:</span>\n<span class=\"w\">                </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">                Determine whether other is equal to this instance.</span>\n<span
    class=\"sd\">                &quot;&quot;&quot;</span>\n                <span
    class=\"k\">if</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">other</span><span class=\"p\">,</span> <span class=\"n\">LifeCycle</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">value</span>
    <span class=\"o\">==</span> <span class=\"n\">other</span><span class=\"o\">.</span><span
    class=\"n\">value</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">other</span><span class=\"p\">,</span> <span
    class=\"nb\">int</span><span class=\"p\">):</span>\n                    <span
    class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">value</span> <span class=\"o\">==</span> <span class=\"n\">other</span>\n
    \               <span class=\"k\">return</span> <span class=\"bp\">NotImplemented</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='__lt__' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'><strong>lt</strong> <em class='small'>method</em></h2>\nDetermine whether
    other is less than this instance.</p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\"><strong>lt</strong> <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"fm\">__lt__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">other</span><span class=\"p\">:</span>
    <span class=\"nb\">object</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"nb\">bool</span><span class=\"p\">:</span>\n<span class=\"w\">                </span><span
    class=\"sd\">&quot;&quot;&quot;</span>\n<span class=\"sd\">                Determine
    whether other is less than this instance.</span>\n<span class=\"sd\">                &quot;&quot;&quot;</span>\n
    \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">other</span><span class=\"p\">,</span> <span
    class=\"n\">LifeCycle</span><span class=\"p\">):</span>\n                    <span
    class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">value</span> <span class=\"o\">&lt;</span> <span class=\"n\">other</span><span
    class=\"o\">.</span><span class=\"n\">value</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">other</span><span
    class=\"p\">,</span> <span class=\"nb\">int</span><span class=\"p\">):</span>\n
    \                   <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">value</span> <span class=\"o\">&lt;</span>
    <span class=\"n\">other</span>\n                <span class=\"k\">return</span>
    <span class=\"bp\">NotImplemented</span>\n</pre></div>\n\n</pre>\n\n<p>!! method
    <h2 id='__eq__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'><strong>eq</strong>
    <em class='small'>method</em></h2>\nDetermine whether other is equal to this instance.</p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\"><strong>eq</strong>
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
    <span class=\"fm\">__eq__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">other</span><span class=\"p\">:</span>
    <span class=\"nb\">object</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"nb\">bool</span><span class=\"p\">:</span>\n<span class=\"w\">                </span><span
    class=\"sd\">&quot;&quot;&quot;</span>\n<span class=\"sd\">                Determine
    whether other is equal to this instance.</span>\n<span class=\"sd\">                &quot;&quot;&quot;</span>\n
    \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">other</span><span class=\"p\">,</span> <span
    class=\"n\">LifeCycle</span><span class=\"p\">):</span>\n                    <span
    class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">value</span> <span class=\"o\">==</span> <span class=\"n\">other</span><span
    class=\"o\">.</span><span class=\"n\">value</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">other</span><span
    class=\"p\">,</span> <span class=\"nb\">int</span><span class=\"p\">):</span>\n
    \                   <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">value</span> <span class=\"o\">==</span>
    <span class=\"n\">other</span>\n                <span class=\"k\">return</span>
    <span class=\"bp\">NotImplemented</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>
    \   </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Lifecycle.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"The LifeCycle is a core component for the internal
    workings of Markata.  It ! ???+ source  ! !\" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Lifecycle.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"The LifeCycle is a core component for
    the internal workings of Markata.  It ! ???+ source  ! !\" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<article style=\"text-align: center;\">\n    <style>\n        section
    {\n            font-size: 200%;\n        }\n\n\n        .edit {\n            display:
    none;\n        }\n    </style>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Lifecycle.Py\n    </h1>\n</section></article>\n     </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Lifecycle.Py\n    </h1>\n</section>    <section class=\"body\">\n        <p>The
    LifeCycle is a core component for the internal workings of Markata.  It\nsets
    fourth the hooks available, the methods to run them on the Markata\ninstance,
    and the order they run in.</p>\n<h3>Usage</h3>\n<pre class='wrapper'>\n\n<div
    class='copy-wrapper'>\n\n<button class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"kn\">from</span>
    <span class=\"nn\">markata</span> <span class=\"kn\">import</span> <span class=\"n\">Lifecycle</span>\n\n<span
    class=\"n\">step</span> <span class=\"o\">=</span> <span class=\"n\">Lifecycle</span><span
    class=\"o\">.</span><span class=\"n\">glob</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='LifeCycle' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>LifeCycle <em class='small'>class</em></h2>\nLifeCycle currently supports
    the following steps.</p>\n<pre><code>* configure - load and fix configuration\n*
    glob - find files\n* load - load files\n* validate_posts\n* pre_render - clean
    up files/metadata before render\n* render - render content\n* post_render - clean
    up rendered content\n* save - store results to disk\n* teardown - runs on exit\n</code></pre>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">LifeCycle
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
    <span class=\"nc\">LifeCycle</span><span class=\"p\">(</span><span class=\"n\">Enum</span><span
    class=\"p\">):</span>\n<span class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            LifeCycle currently supports the following steps.</span>\n\n\n<span
    class=\"sd\">            * configure - load and fix configuration</span>\n<span
    class=\"sd\">            * glob - find files</span>\n<span class=\"sd\">            *
    load - load files</span>\n<span class=\"sd\">            * validate_posts</span>\n<span
    class=\"sd\">            * pre_render - clean up files/metadata before render</span>\n<span
    class=\"sd\">            * render - render content</span>\n<span class=\"sd\">
    \           * post_render - clean up rendered content</span>\n<span class=\"sd\">
    \           * save - store results to disk</span>\n<span class=\"sd\">            *
    teardown - runs on exit</span>\n\n<span class=\"sd\">            &quot;&quot;&quot;</span>\n\n
    \           <span class=\"n\">config_model</span> <span class=\"o\">=</span> <span
    class=\"n\">auto</span><span class=\"p\">()</span>\n            <span class=\"n\">post_model</span>
    <span class=\"o\">=</span> <span class=\"n\">auto</span><span class=\"p\">()</span>\n
    \           <span class=\"n\">create_models</span> <span class=\"o\">=</span>
    <span class=\"n\">auto</span><span class=\"p\">()</span>\n            <span class=\"n\">load_config</span>
    <span class=\"o\">=</span> <span class=\"n\">auto</span><span class=\"p\">()</span>\n
    \           <span class=\"n\">configure</span> <span class=\"o\">=</span> <span
    class=\"n\">auto</span><span class=\"p\">()</span>\n            <span class=\"n\">validate_config</span>
    <span class=\"o\">=</span> <span class=\"n\">auto</span><span class=\"p\">()</span>\n
    \           <span class=\"n\">glob</span> <span class=\"o\">=</span> <span class=\"n\">auto</span><span
    class=\"p\">()</span>\n            <span class=\"n\">load</span> <span class=\"o\">=</span>
    <span class=\"n\">auto</span><span class=\"p\">()</span>\n            <span class=\"n\">pre_render</span>
    <span class=\"o\">=</span> <span class=\"n\">auto</span><span class=\"p\">()</span>\n
    \           <span class=\"n\">render</span> <span class=\"o\">=</span> <span class=\"n\">auto</span><span
    class=\"p\">()</span>\n            <span class=\"n\">post_render</span> <span
    class=\"o\">=</span> <span class=\"n\">auto</span><span class=\"p\">()</span>\n
    \           <span class=\"n\">save</span> <span class=\"o\">=</span> <span class=\"n\">auto</span><span
    class=\"p\">()</span>\n            <span class=\"n\">teardown</span> <span class=\"o\">=</span>
    <span class=\"n\">auto</span><span class=\"p\">()</span>\n\n            <span
    class=\"k\">def</span> <span class=\"fm\">__lt__</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"n\">other</span><span
    class=\"p\">:</span> <span class=\"nb\">object</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">bool</span><span class=\"p\">:</span>\n<span
    class=\"w\">                </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">                Determine whether other is less than this instance.</span>\n<span
    class=\"sd\">                &quot;&quot;&quot;</span>\n                <span
    class=\"k\">if</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">other</span><span class=\"p\">,</span> <span class=\"n\">LifeCycle</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">value</span>
    <span class=\"o\">&lt;</span> <span class=\"n\">other</span><span class=\"o\">.</span><span
    class=\"n\">value</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">other</span><span class=\"p\">,</span> <span
    class=\"nb\">int</span><span class=\"p\">):</span>\n                    <span
    class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">value</span> <span class=\"o\">&lt;</span> <span class=\"n\">other</span>\n
    \               <span class=\"k\">return</span> <span class=\"bp\">NotImplemented</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"fm\">__eq__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"n\">other</span><span class=\"p\">:</span> <span class=\"nb\">object</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">bool</span><span
    class=\"p\">:</span>\n<span class=\"w\">                </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">                Determine whether other is equal to this instance.</span>\n<span
    class=\"sd\">                &quot;&quot;&quot;</span>\n                <span
    class=\"k\">if</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">other</span><span class=\"p\">,</span> <span class=\"n\">LifeCycle</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">value</span>
    <span class=\"o\">==</span> <span class=\"n\">other</span><span class=\"o\">.</span><span
    class=\"n\">value</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">other</span><span class=\"p\">,</span> <span
    class=\"nb\">int</span><span class=\"p\">):</span>\n                    <span
    class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">value</span> <span class=\"o\">==</span> <span class=\"n\">other</span>\n
    \               <span class=\"k\">return</span> <span class=\"bp\">NotImplemented</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='__lt__' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'><strong>lt</strong> <em class='small'>method</em></h2>\nDetermine whether
    other is less than this instance.</p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\"><strong>lt</strong> <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"fm\">__lt__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">other</span><span class=\"p\">:</span>
    <span class=\"nb\">object</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"nb\">bool</span><span class=\"p\">:</span>\n<span class=\"w\">                </span><span
    class=\"sd\">&quot;&quot;&quot;</span>\n<span class=\"sd\">                Determine
    whether other is less than this instance.</span>\n<span class=\"sd\">                &quot;&quot;&quot;</span>\n
    \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">other</span><span class=\"p\">,</span> <span
    class=\"n\">LifeCycle</span><span class=\"p\">):</span>\n                    <span
    class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">value</span> <span class=\"o\">&lt;</span> <span class=\"n\">other</span><span
    class=\"o\">.</span><span class=\"n\">value</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">other</span><span
    class=\"p\">,</span> <span class=\"nb\">int</span><span class=\"p\">):</span>\n
    \                   <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">value</span> <span class=\"o\">&lt;</span>
    <span class=\"n\">other</span>\n                <span class=\"k\">return</span>
    <span class=\"bp\">NotImplemented</span>\n</pre></div>\n\n</pre>\n\n<p>!! method
    <h2 id='__eq__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'><strong>eq</strong>
    <em class='small'>method</em></h2>\nDetermine whether other is equal to this instance.</p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\"><strong>eq</strong>
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
    <span class=\"fm\">__eq__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">other</span><span class=\"p\">:</span>
    <span class=\"nb\">object</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"nb\">bool</span><span class=\"p\">:</span>\n<span class=\"w\">                </span><span
    class=\"sd\">&quot;&quot;&quot;</span>\n<span class=\"sd\">                Determine
    whether other is equal to this instance.</span>\n<span class=\"sd\">                &quot;&quot;&quot;</span>\n
    \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">other</span><span class=\"p\">,</span> <span
    class=\"n\">LifeCycle</span><span class=\"p\">):</span>\n                    <span
    class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">value</span> <span class=\"o\">==</span> <span class=\"n\">other</span><span
    class=\"o\">.</span><span class=\"n\">value</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">other</span><span
    class=\"p\">,</span> <span class=\"nb\">int</span><span class=\"p\">):</span>\n
    \                   <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">value</span> <span class=\"o\">==</span>
    <span class=\"n\">other</span>\n                <span class=\"k\">return</span>
    <span class=\"bp\">NotImplemented</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>"
  raw.md: ''
published: true
slug: markata/lifecycle
title: Lifecycle.Py


---

The LifeCycle is a core component for the internal workings of Markata.  It
sets fourth the hooks available, the methods to run them on the Markata
instance, and the order they run in.

### Usage

``` python
from markata import Lifecycle

step = Lifecycle.glob
```


!! class <h2 id='LifeCycle' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>LifeCycle <em class='small'>class</em></h2>
    LifeCycle currently supports the following steps.


    * configure - load and fix configuration
    * glob - find files
    * load - load files
    * validate_posts
    * pre_render - clean up files/metadata before render
    * render - render content
    * post_render - clean up rendered content
    * save - store results to disk
    * teardown - runs on exit
???+ source "LifeCycle <em class='small'>source</em>"

```python

        class LifeCycle(Enum):
            """
            LifeCycle currently supports the following steps.


            * configure - load and fix configuration
            * glob - find files
            * load - load files
            * validate_posts
            * pre_render - clean up files/metadata before render
            * render - render content
            * post_render - clean up rendered content
            * save - store results to disk
            * teardown - runs on exit

            """

            config_model = auto()
            post_model = auto()
            create_models = auto()
            load_config = auto()
            configure = auto()
            validate_config = auto()
            glob = auto()
            load = auto()
            pre_render = auto()
            render = auto()
            post_render = auto()
            save = auto()
            teardown = auto()

            def __lt__(self, other: object) -> bool:
                """
                Determine whether other is less than this instance.
                """
                if isinstance(other, LifeCycle):
                    return self.value < other.value
                if isinstance(other, int):
                    return self.value < other
                return NotImplemented

            def __eq__(self, other: object) -> bool:
                """
                Determine whether other is equal to this instance.
                """
                if isinstance(other, LifeCycle):
                    return self.value == other.value
                if isinstance(other, int):
                    return self.value == other
                return NotImplemented
```


!! method <h2 id='__lt__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__lt__ <em class='small'>method</em></h2>
    Determine whether other is less than this instance.
???+ source "__lt__ <em class='small'>source</em>"

```python

        def __lt__(self, other: object) -> bool:
                """
                Determine whether other is less than this instance.
                """
                if isinstance(other, LifeCycle):
                    return self.value < other.value
                if isinstance(other, int):
                    return self.value < other
                return NotImplemented
```


!! method <h2 id='__eq__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__eq__ <em class='small'>method</em></h2>
    Determine whether other is equal to this instance.
???+ source "__eq__ <em class='small'>source</em>"

```python

        def __eq__(self, other: object) -> bool:
                """
                Determine whether other is equal to this instance.
                """
                if isinstance(other, LifeCycle):
                    return self.value == other.value
                if isinstance(other, int):
                    return self.value == other
                return NotImplemented
```

