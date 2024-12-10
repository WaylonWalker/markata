---
content: "Define hook specs.\n\n\n!! class <h2 id='MarkataSpecs' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>MarkataSpecs <em class='small'>class</em></h2>\n
  \   Namespace that defines all specifications for Load hooks.\n\n    configure ->
  glob -> load -> render -> save\n???+ source \"MarkataSpecs <em class='small'>source</em>\"\n\n```python\n\n
  \       class MarkataSpecs:\n            \"\"\"\n            Namespace that defines
  all specifications for Load hooks.\n\n            configure -> glob -> load -> render
  -> save\n            \"\"\"\n```\n\n\n!! function <h2 id='generic_lifecycle_method'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>generic_lifecycle_method
  <em class='small'>function</em></h2>\n\n???+ source \"generic_lifecycle_method <em
  class='small'>source</em>\"\n\n```python\n\n        def generic_lifecycle_method(\n
  \           markata: \"Markata\",\n        ) -> Any: ...\n```\n\n\n!! function <h2
  id='cli_lifecycle_method' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>cli_lifecycle_method <em class='small'>function</em></h2>\n    A Markata
  lifecycle methos that includes a typer app used for cli's\n???+ source \"cli_lifecycle_method
  <em class='small'>source</em>\"\n\n```python\n\n        def cli_lifecycle_method(markata:
  \"Markata\", app: \"typer.Typer\") -> Any:\n            \"A Markata lifecycle methos
  that includes a typer app used for cli's\"\n```\n\n\n!! function <h2 id='register_attr'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>register_attr <em
  class='small'>function</em></h2>\n\n???+ source \"register_attr <em class='small'>source</em>\"\n\n```python\n\n
  \       def register_attr(*attrs: Any) -> Callable:\n            def decorator_register(\n
  \               func: Callable,\n            ) -> Callable:\n                for
  attr in attrs:\n                    if attr not in registered_attrs:\n                        registered_attrs[attr]
  = []\n                    registered_attrs[attr].append(\n                        {\n
  \                           \"func\": func,\n                            \"funcname\":
  func.__code__.co_name,\n                            \"lifecycle\": getattr(LifeCycle,
  func.__code__.co_name),\n                        },\n                    )\n\n                @functools.wraps(func)\n
  \               def wrapper_register(markata: \"Markata\", *args: Any, **kwargs:
  Any) -> Any:\n                    return func(markata, *args, **kwargs)\n\n                return
  wrapper_register\n\n            return decorator_register\n```\n\n\n!! function
  <h2 id='decorator_register' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>decorator_register <em class='small'>function</em></h2>\n\n???+ source \"decorator_register
  <em class='small'>source</em>\"\n\n```python\n\n        def decorator_register(\n
  \               func: Callable,\n            ) -> Callable:\n                for
  attr in attrs:\n                    if attr not in registered_attrs:\n                        registered_attrs[attr]
  = []\n                    registered_attrs[attr].append(\n                        {\n
  \                           \"func\": func,\n                            \"funcname\":
  func.__code__.co_name,\n                            \"lifecycle\": getattr(LifeCycle,
  func.__code__.co_name),\n                        },\n                    )\n\n                @functools.wraps(func)\n
  \               def wrapper_register(markata: \"Markata\", *args: Any, **kwargs:
  Any) -> Any:\n                    return func(markata, *args, **kwargs)\n\n                return
  wrapper_register\n```\n\n\n!! function <h2 id='wrapper_register' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>wrapper_register <em class='small'>function</em></h2>\n\n???+
  source \"wrapper_register <em class='small'>source</em>\"\n\n```python\n\n        def
  wrapper_register(markata: \"Markata\", *args: Any, **kwargs: Any) -> Any:\n                    return
  func(markata, *args, **kwargs)\n```\n\n"
date: 0001-01-01
description: 'Define hook specs. ! ???+ source  ! ???+ source  ! ! ???+ source  !
  ???+ source  ! ???+ source '
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Hookspec.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"Define hook specs. ! ???+ source  ! ???+ source
    \ ! ! ???+ source  ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Hookspec.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Define hook specs. ! ???+ source  ! ???+
    source  ! ! ???+ source  ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\"
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Hookspec.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <p>Define hook specs.</p>\n<p>!!
    class <h2 id='MarkataSpecs' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>MarkataSpecs <em class='small'>class</em></h2>\nNamespace that defines
    all specifications for Load hooks.</p>\n<pre><code>configure -&gt; glob -&gt;
    load -&gt; render -&gt; save\n</code></pre>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">MarkataSpecs <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nc\">MarkataSpecs</span><span class=\"p\">:</span>\n<span class=\"w\">
    \           </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span class=\"sd\">
    \           Namespace that defines all specifications for Load hooks.</span>\n\n<span
    class=\"sd\">            configure -&gt; glob -&gt; load -&gt; render -&gt; save</span>\n<span
    class=\"sd\">            &quot;&quot;&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='generic_lifecycle_method' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>generic_lifecycle_method <em class='small'>function</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">generic_lifecycle_method
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
    <span class=\"nf\">generic_lifecycle_method</span><span class=\"p\">(</span>\n
    \           <span class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span>\n        <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Any</span><span class=\"p\">:</span> <span class=\"o\">...</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='cli_lifecycle_method' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>cli_lifecycle_method <em class='small'>function</em></h2>\nA Markata lifecycle
    methos that includes a typer app used for cli's</p>\n<div class=\"admonition source
    is-collapsible collapsible-open\">\n<p class=\"admonition-title\">cli_lifecycle_method
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
    <span class=\"nf\">cli_lifecycle_method</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">app</span><span class=\"p\">:</span> <span
    class=\"s2\">&quot;typer.Typer&quot;</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Any</span><span class=\"p\">:</span>\n            <span class=\"s2\">&quot;A
    Markata lifecycle methos that includes a typer app used for cli&#39;s&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='register_attr' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>register_attr <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">register_attr
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
    <span class=\"nf\">register_attr</span><span class=\"p\">(</span><span class=\"o\">*</span><span
    class=\"n\">attrs</span><span class=\"p\">:</span> <span class=\"n\">Any</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Callable</span><span
    class=\"p\">:</span>\n            <span class=\"k\">def</span> <span class=\"nf\">decorator_register</span><span
    class=\"p\">(</span>\n                <span class=\"n\">func</span><span class=\"p\">:</span>
    <span class=\"n\">Callable</span><span class=\"p\">,</span>\n            <span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Callable</span><span
    class=\"p\">:</span>\n                <span class=\"k\">for</span> <span class=\"n\">attr</span>
    <span class=\"ow\">in</span> <span class=\"n\">attrs</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">if</span> <span class=\"n\">attr</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">registered_attrs</span><span
    class=\"p\">:</span>\n                        <span class=\"n\">registered_attrs</span><span
    class=\"p\">[</span><span class=\"n\">attr</span><span class=\"p\">]</span> <span
    class=\"o\">=</span> <span class=\"p\">[]</span>\n                    <span class=\"n\">registered_attrs</span><span
    class=\"p\">[</span><span class=\"n\">attr</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span>\n
    \                       <span class=\"p\">{</span>\n                            <span
    class=\"s2\">&quot;func&quot;</span><span class=\"p\">:</span> <span class=\"n\">func</span><span
    class=\"p\">,</span>\n                            <span class=\"s2\">&quot;funcname&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">func</span><span class=\"o\">.</span><span
    class=\"vm\">__code__</span><span class=\"o\">.</span><span class=\"n\">co_name</span><span
    class=\"p\">,</span>\n                            <span class=\"s2\">&quot;lifecycle&quot;</span><span
    class=\"p\">:</span> <span class=\"nb\">getattr</span><span class=\"p\">(</span><span
    class=\"n\">LifeCycle</span><span class=\"p\">,</span> <span class=\"n\">func</span><span
    class=\"o\">.</span><span class=\"vm\">__code__</span><span class=\"o\">.</span><span
    class=\"n\">co_name</span><span class=\"p\">),</span>\n                        <span
    class=\"p\">},</span>\n                    <span class=\"p\">)</span>\n\n                <span
    class=\"nd\">@functools</span><span class=\"o\">.</span><span class=\"n\">wraps</span><span
    class=\"p\">(</span><span class=\"n\">func</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">def</span> <span class=\"nf\">wrapper_register</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"n\">args</span><span class=\"p\">:</span> <span
    class=\"n\">Any</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"n\">kwargs</span><span class=\"p\">:</span> <span class=\"n\">Any</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">func</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"o\">*</span><span class=\"n\">args</span><span
    class=\"p\">,</span> <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">return</span> <span
    class=\"n\">wrapper_register</span>\n\n            <span class=\"k\">return</span>
    <span class=\"n\">decorator_register</span>\n</pre></div>\n\n</pre>\n\n<p>!! function
    <h2 id='decorator_register' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>decorator_register <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">decorator_register
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
    <span class=\"nf\">decorator_register</span><span class=\"p\">(</span>\n                <span
    class=\"n\">func</span><span class=\"p\">:</span> <span class=\"n\">Callable</span><span
    class=\"p\">,</span>\n            <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Callable</span><span class=\"p\">:</span>\n                <span
    class=\"k\">for</span> <span class=\"n\">attr</span> <span class=\"ow\">in</span>
    <span class=\"n\">attrs</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">if</span> <span class=\"n\">attr</span> <span class=\"ow\">not</span>
    <span class=\"ow\">in</span> <span class=\"n\">registered_attrs</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">registered_attrs</span><span class=\"p\">[</span><span
    class=\"n\">attr</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"p\">[]</span>\n                    <span class=\"n\">registered_attrs</span><span
    class=\"p\">[</span><span class=\"n\">attr</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span>\n
    \                       <span class=\"p\">{</span>\n                            <span
    class=\"s2\">&quot;func&quot;</span><span class=\"p\">:</span> <span class=\"n\">func</span><span
    class=\"p\">,</span>\n                            <span class=\"s2\">&quot;funcname&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">func</span><span class=\"o\">.</span><span
    class=\"vm\">__code__</span><span class=\"o\">.</span><span class=\"n\">co_name</span><span
    class=\"p\">,</span>\n                            <span class=\"s2\">&quot;lifecycle&quot;</span><span
    class=\"p\">:</span> <span class=\"nb\">getattr</span><span class=\"p\">(</span><span
    class=\"n\">LifeCycle</span><span class=\"p\">,</span> <span class=\"n\">func</span><span
    class=\"o\">.</span><span class=\"vm\">__code__</span><span class=\"o\">.</span><span
    class=\"n\">co_name</span><span class=\"p\">),</span>\n                        <span
    class=\"p\">},</span>\n                    <span class=\"p\">)</span>\n\n                <span
    class=\"nd\">@functools</span><span class=\"o\">.</span><span class=\"n\">wraps</span><span
    class=\"p\">(</span><span class=\"n\">func</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">def</span> <span class=\"nf\">wrapper_register</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"n\">args</span><span class=\"p\">:</span> <span
    class=\"n\">Any</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"n\">kwargs</span><span class=\"p\">:</span> <span class=\"n\">Any</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">func</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"o\">*</span><span class=\"n\">args</span><span
    class=\"p\">,</span> <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">return</span> <span
    class=\"n\">wrapper_register</span>\n</pre></div>\n\n</pre>\n\n<p>!! function
    <h2 id='wrapper_register' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>wrapper_register <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">wrapper_register
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
    <span class=\"nf\">wrapper_register</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span>
    <span class=\"o\">*</span><span class=\"n\">args</span><span class=\"p\">:</span>
    <span class=\"n\">Any</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"n\">kwargs</span><span class=\"p\">:</span> <span class=\"n\">Any</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">func</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"o\">*</span><span class=\"n\">args</span><span
    class=\"p\">,</span> <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>    </div>\n
    \   <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Hookspec.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"Define hook specs. ! ???+ source  ! ???+ source
    \ ! ! ???+ source  ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Hookspec.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Define hook specs. ! ???+ source  ! ???+
    source  ! ! ???+ source  ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<article style=\"text-align: center;\">\n    <style>\n        section
    {\n            font-size: 200%;\n        }\n\n\n        .edit {\n            display:
    none;\n        }\n    </style>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Hookspec.Py\n    </h1>\n</section></article>\n     </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Hookspec.Py\n    </h1>\n</section>    <section class=\"body\">\n        <p>Define
    hook specs.</p>\n<p>!! class <h2 id='MarkataSpecs' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>MarkataSpecs <em class='small'>class</em></h2>\nNamespace that defines
    all specifications for Load hooks.</p>\n<pre><code>configure -&gt; glob -&gt;
    load -&gt; render -&gt; save\n</code></pre>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">MarkataSpecs <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nc\">MarkataSpecs</span><span class=\"p\">:</span>\n<span class=\"w\">
    \           </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span class=\"sd\">
    \           Namespace that defines all specifications for Load hooks.</span>\n\n<span
    class=\"sd\">            configure -&gt; glob -&gt; load -&gt; render -&gt; save</span>\n<span
    class=\"sd\">            &quot;&quot;&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='generic_lifecycle_method' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>generic_lifecycle_method <em class='small'>function</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">generic_lifecycle_method
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
    <span class=\"nf\">generic_lifecycle_method</span><span class=\"p\">(</span>\n
    \           <span class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span>\n        <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Any</span><span class=\"p\">:</span> <span class=\"o\">...</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='cli_lifecycle_method' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>cli_lifecycle_method <em class='small'>function</em></h2>\nA Markata lifecycle
    methos that includes a typer app used for cli's</p>\n<div class=\"admonition source
    is-collapsible collapsible-open\">\n<p class=\"admonition-title\">cli_lifecycle_method
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
    <span class=\"nf\">cli_lifecycle_method</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">app</span><span class=\"p\">:</span> <span
    class=\"s2\">&quot;typer.Typer&quot;</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Any</span><span class=\"p\">:</span>\n            <span class=\"s2\">&quot;A
    Markata lifecycle methos that includes a typer app used for cli&#39;s&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='register_attr' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>register_attr <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">register_attr
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
    <span class=\"nf\">register_attr</span><span class=\"p\">(</span><span class=\"o\">*</span><span
    class=\"n\">attrs</span><span class=\"p\">:</span> <span class=\"n\">Any</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Callable</span><span
    class=\"p\">:</span>\n            <span class=\"k\">def</span> <span class=\"nf\">decorator_register</span><span
    class=\"p\">(</span>\n                <span class=\"n\">func</span><span class=\"p\">:</span>
    <span class=\"n\">Callable</span><span class=\"p\">,</span>\n            <span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Callable</span><span
    class=\"p\">:</span>\n                <span class=\"k\">for</span> <span class=\"n\">attr</span>
    <span class=\"ow\">in</span> <span class=\"n\">attrs</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">if</span> <span class=\"n\">attr</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">registered_attrs</span><span
    class=\"p\">:</span>\n                        <span class=\"n\">registered_attrs</span><span
    class=\"p\">[</span><span class=\"n\">attr</span><span class=\"p\">]</span> <span
    class=\"o\">=</span> <span class=\"p\">[]</span>\n                    <span class=\"n\">registered_attrs</span><span
    class=\"p\">[</span><span class=\"n\">attr</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span>\n
    \                       <span class=\"p\">{</span>\n                            <span
    class=\"s2\">&quot;func&quot;</span><span class=\"p\">:</span> <span class=\"n\">func</span><span
    class=\"p\">,</span>\n                            <span class=\"s2\">&quot;funcname&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">func</span><span class=\"o\">.</span><span
    class=\"vm\">__code__</span><span class=\"o\">.</span><span class=\"n\">co_name</span><span
    class=\"p\">,</span>\n                            <span class=\"s2\">&quot;lifecycle&quot;</span><span
    class=\"p\">:</span> <span class=\"nb\">getattr</span><span class=\"p\">(</span><span
    class=\"n\">LifeCycle</span><span class=\"p\">,</span> <span class=\"n\">func</span><span
    class=\"o\">.</span><span class=\"vm\">__code__</span><span class=\"o\">.</span><span
    class=\"n\">co_name</span><span class=\"p\">),</span>\n                        <span
    class=\"p\">},</span>\n                    <span class=\"p\">)</span>\n\n                <span
    class=\"nd\">@functools</span><span class=\"o\">.</span><span class=\"n\">wraps</span><span
    class=\"p\">(</span><span class=\"n\">func</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">def</span> <span class=\"nf\">wrapper_register</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"n\">args</span><span class=\"p\">:</span> <span
    class=\"n\">Any</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"n\">kwargs</span><span class=\"p\">:</span> <span class=\"n\">Any</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">func</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"o\">*</span><span class=\"n\">args</span><span
    class=\"p\">,</span> <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">return</span> <span
    class=\"n\">wrapper_register</span>\n\n            <span class=\"k\">return</span>
    <span class=\"n\">decorator_register</span>\n</pre></div>\n\n</pre>\n\n<p>!! function
    <h2 id='decorator_register' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>decorator_register <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">decorator_register
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
    <span class=\"nf\">decorator_register</span><span class=\"p\">(</span>\n                <span
    class=\"n\">func</span><span class=\"p\">:</span> <span class=\"n\">Callable</span><span
    class=\"p\">,</span>\n            <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Callable</span><span class=\"p\">:</span>\n                <span
    class=\"k\">for</span> <span class=\"n\">attr</span> <span class=\"ow\">in</span>
    <span class=\"n\">attrs</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">if</span> <span class=\"n\">attr</span> <span class=\"ow\">not</span>
    <span class=\"ow\">in</span> <span class=\"n\">registered_attrs</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">registered_attrs</span><span class=\"p\">[</span><span
    class=\"n\">attr</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"p\">[]</span>\n                    <span class=\"n\">registered_attrs</span><span
    class=\"p\">[</span><span class=\"n\">attr</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span>\n
    \                       <span class=\"p\">{</span>\n                            <span
    class=\"s2\">&quot;func&quot;</span><span class=\"p\">:</span> <span class=\"n\">func</span><span
    class=\"p\">,</span>\n                            <span class=\"s2\">&quot;funcname&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">func</span><span class=\"o\">.</span><span
    class=\"vm\">__code__</span><span class=\"o\">.</span><span class=\"n\">co_name</span><span
    class=\"p\">,</span>\n                            <span class=\"s2\">&quot;lifecycle&quot;</span><span
    class=\"p\">:</span> <span class=\"nb\">getattr</span><span class=\"p\">(</span><span
    class=\"n\">LifeCycle</span><span class=\"p\">,</span> <span class=\"n\">func</span><span
    class=\"o\">.</span><span class=\"vm\">__code__</span><span class=\"o\">.</span><span
    class=\"n\">co_name</span><span class=\"p\">),</span>\n                        <span
    class=\"p\">},</span>\n                    <span class=\"p\">)</span>\n\n                <span
    class=\"nd\">@functools</span><span class=\"o\">.</span><span class=\"n\">wraps</span><span
    class=\"p\">(</span><span class=\"n\">func</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">def</span> <span class=\"nf\">wrapper_register</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"n\">args</span><span class=\"p\">:</span> <span
    class=\"n\">Any</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"n\">kwargs</span><span class=\"p\">:</span> <span class=\"n\">Any</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">func</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"o\">*</span><span class=\"n\">args</span><span
    class=\"p\">,</span> <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">return</span> <span
    class=\"n\">wrapper_register</span>\n</pre></div>\n\n</pre>\n\n<p>!! function
    <h2 id='wrapper_register' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>wrapper_register <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">wrapper_register
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
    <span class=\"nf\">wrapper_register</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span>
    <span class=\"o\">*</span><span class=\"n\">args</span><span class=\"p\">:</span>
    <span class=\"n\">Any</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"n\">kwargs</span><span class=\"p\">:</span> <span class=\"n\">Any</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">func</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"o\">*</span><span class=\"n\">args</span><span
    class=\"p\">,</span> <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>"
  raw.md: ''
published: true
slug: markata/hookspec
title: Hookspec.Py


---

Define hook specs.


!! class <h2 id='MarkataSpecs' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>MarkataSpecs <em class='small'>class</em></h2>
    Namespace that defines all specifications for Load hooks.

    configure -> glob -> load -> render -> save
???+ source "MarkataSpecs <em class='small'>source</em>"

```python

        class MarkataSpecs:
            """
            Namespace that defines all specifications for Load hooks.

            configure -> glob -> load -> render -> save
            """
```


!! function <h2 id='generic_lifecycle_method' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>generic_lifecycle_method <em class='small'>function</em></h2>

???+ source "generic_lifecycle_method <em class='small'>source</em>"

```python

        def generic_lifecycle_method(
            markata: "Markata",
        ) -> Any: ...
```


!! function <h2 id='cli_lifecycle_method' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>cli_lifecycle_method <em class='small'>function</em></h2>
    A Markata lifecycle methos that includes a typer app used for cli's
???+ source "cli_lifecycle_method <em class='small'>source</em>"

```python

        def cli_lifecycle_method(markata: "Markata", app: "typer.Typer") -> Any:
            "A Markata lifecycle methos that includes a typer app used for cli's"
```


!! function <h2 id='register_attr' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>register_attr <em class='small'>function</em></h2>

???+ source "register_attr <em class='small'>source</em>"

```python

        def register_attr(*attrs: Any) -> Callable:
            def decorator_register(
                func: Callable,
            ) -> Callable:
                for attr in attrs:
                    if attr not in registered_attrs:
                        registered_attrs[attr] = []
                    registered_attrs[attr].append(
                        {
                            "func": func,
                            "funcname": func.__code__.co_name,
                            "lifecycle": getattr(LifeCycle, func.__code__.co_name),
                        },
                    )

                @functools.wraps(func)
                def wrapper_register(markata: "Markata", *args: Any, **kwargs: Any) -> Any:
                    return func(markata, *args, **kwargs)

                return wrapper_register

            return decorator_register
```


!! function <h2 id='decorator_register' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>decorator_register <em class='small'>function</em></h2>

???+ source "decorator_register <em class='small'>source</em>"

```python

        def decorator_register(
                func: Callable,
            ) -> Callable:
                for attr in attrs:
                    if attr not in registered_attrs:
                        registered_attrs[attr] = []
                    registered_attrs[attr].append(
                        {
                            "func": func,
                            "funcname": func.__code__.co_name,
                            "lifecycle": getattr(LifeCycle, func.__code__.co_name),
                        },
                    )

                @functools.wraps(func)
                def wrapper_register(markata: "Markata", *args: Any, **kwargs: Any) -> Any:
                    return func(markata, *args, **kwargs)

                return wrapper_register
```


!! function <h2 id='wrapper_register' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>wrapper_register <em class='small'>function</em></h2>

???+ source "wrapper_register <em class='small'>source</em>"

```python

        def wrapper_register(markata: "Markata", *args: Any, **kwargs: Any) -> Any:
                    return func(markata, *args, **kwargs)
```

