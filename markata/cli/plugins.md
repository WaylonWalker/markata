---
content: "None\n\n\n!! class <h2 id='Plugins' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>Plugins <em class='small'>class</em></h2>\n\n???+ source \"Plugins <em class='small'>source</em>\"\n\n```python\n\n
  \       class Plugins:\n            def __init__(self, markata: \"Markata\") ->
  None:\n                self.m = markata\n\n            def __rich__(self) -> Panel:\n
  \               grid = Table.grid(expand=True)\n                num_plugins = f\"[bright_blue]({len(self.m._pm.get_plugins())})[/]\"\n\n
  \               for plugin in self.m._pm.get_plugins():\n                    grid.add_row(\n
  \                       \"\".join(\n                            [\n                                \"[bright_black]\",\n
  \                               \".\".join(plugin.__name__.split(\".\")[:-1]),\n
  \                               \".[/]\",\n                                plugin.__name__.split(\".\")[-1],\n
  \                           ],\n                        ),\n                    )\n
  \               return Panel(\n                    grid,\n                    title=f\"plugins
  {num_plugins}\",\n                    border_style=\"gold1\",\n                    expand=False,\n
  \               )\n```\n\n\n!! function <h2 id='configure' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>configure <em class='small'>function</em></h2>\n\n???+
  source \"configure <em class='small'>source</em>\"\n\n```python\n\n        def configure(markata:
  \"Markata\") -> None:\n            def get_plugins(self):\n                try:\n
  \                   return self._plugins\n                except AttributeError:\n
  \                   self._plugins: Plugins = Plugins(self)\n                    return
  self._plugins\n\n            from markata import Markata\n\n            Markata.plugins
  = property(get_plugins)\n```\n\n\n!! method <h2 id='__init__' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>__init__ <em class='small'>method</em></h2>\n\n???+
  source \"__init__ <em class='small'>source</em>\"\n\n```python\n\n        def __init__(self,
  markata: \"Markata\") -> None:\n                self.m = markata\n```\n\n\n!! method
  <h2 id='__rich__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__rich__
  <em class='small'>method</em></h2>\n\n???+ source \"__rich__ <em class='small'>source</em>\"\n\n```python\n\n
  \       def __rich__(self) -> Panel:\n                grid = Table.grid(expand=True)\n
  \               num_plugins = f\"[bright_blue]({len(self.m._pm.get_plugins())})[/]\"\n\n
  \               for plugin in self.m._pm.get_plugins():\n                    grid.add_row(\n
  \                       \"\".join(\n                            [\n                                \"[bright_black]\",\n
  \                               \".\".join(plugin.__name__.split(\".\")[:-1]),\n
  \                               \".[/]\",\n                                plugin.__name__.split(\".\")[-1],\n
  \                           ],\n                        ),\n                    )\n
  \               return Panel(\n                    grid,\n                    title=f\"plugins
  {num_plugins}\",\n                    border_style=\"gold1\",\n                    expand=False,\n
  \               )\n```\n\n\n!! function <h2 id='get_plugins' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>get_plugins <em class='small'>function</em></h2>\n\n???+
  source \"get_plugins <em class='small'>source</em>\"\n\n```python\n\n        def
  get_plugins(self):\n                try:\n                    return self._plugins\n
  \               except AttributeError:\n                    self._plugins: Plugins
  = Plugins(self)\n                    return self._plugins\n```\n\n"
date: 0001-01-01
description: 'None ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
  source '
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Plugins.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"None ! ???+ source  ! ???+ source  ! ???+ source
    \ ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Plugins.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None ! ???+ source  ! ???+ source  !
    ???+ source  ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\"
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Plugins.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <p>None</p>\n<p>!!
    class <h2 id='Plugins' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>Plugins <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">Plugins
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
    <span class=\"nc\">Plugins</span><span class=\"p\">:</span>\n            <span
    class=\"k\">def</span> <span class=\"fm\">__init__</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span> <span class=\"o\">=</span> <span class=\"n\">markata</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"nf\">__rich__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Panel</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">grid</span> <span class=\"o\">=</span> <span
    class=\"n\">Table</span><span class=\"o\">.</span><span class=\"n\">grid</span><span
    class=\"p\">(</span><span class=\"n\">expand</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n                <span class=\"n\">num_plugins</span>
    <span class=\"o\">=</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;[bright_blue](</span><span
    class=\"si\">{</span><span class=\"nb\">len</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">m</span><span
    class=\"o\">.</span><span class=\"n\">_pm</span><span class=\"o\">.</span><span
    class=\"n\">get_plugins</span><span class=\"p\">())</span><span class=\"si\">}</span><span
    class=\"s2\">)[/]&quot;</span>\n\n                <span class=\"k\">for</span>
    <span class=\"n\">plugin</span> <span class=\"ow\">in</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">m</span><span class=\"o\">.</span><span
    class=\"n\">_pm</span><span class=\"o\">.</span><span class=\"n\">get_plugins</span><span
    class=\"p\">():</span>\n                    <span class=\"n\">grid</span><span
    class=\"o\">.</span><span class=\"n\">add_row</span><span class=\"p\">(</span>\n
    \                       <span class=\"s2\">&quot;&quot;</span><span class=\"o\">.</span><span
    class=\"n\">join</span><span class=\"p\">(</span>\n                            <span
    class=\"p\">[</span>\n                                <span class=\"s2\">&quot;[bright_black]&quot;</span><span
    class=\"p\">,</span>\n                                <span class=\"s2\">&quot;.&quot;</span><span
    class=\"o\">.</span><span class=\"n\">join</span><span class=\"p\">(</span><span
    class=\"n\">plugin</span><span class=\"o\">.</span><span class=\"vm\">__name__</span><span
    class=\"o\">.</span><span class=\"n\">split</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;.&quot;</span><span class=\"p\">)[:</span><span class=\"o\">-</span><span
    class=\"mi\">1</span><span class=\"p\">]),</span>\n                                <span
    class=\"s2\">&quot;.[/]&quot;</span><span class=\"p\">,</span>\n                                <span
    class=\"n\">plugin</span><span class=\"o\">.</span><span class=\"vm\">__name__</span><span
    class=\"o\">.</span><span class=\"n\">split</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;.&quot;</span><span class=\"p\">)[</span><span class=\"o\">-</span><span
    class=\"mi\">1</span><span class=\"p\">],</span>\n                            <span
    class=\"p\">],</span>\n                        <span class=\"p\">),</span>\n                    <span
    class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"n\">Panel</span><span
    class=\"p\">(</span>\n                    <span class=\"n\">grid</span><span class=\"p\">,</span>\n
    \                   <span class=\"n\">title</span><span class=\"o\">=</span><span
    class=\"sa\">f</span><span class=\"s2\">&quot;plugins </span><span class=\"si\">{</span><span
    class=\"n\">num_plugins</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"n\">border_style</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;gold1&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"n\">expand</span><span class=\"o\">=</span><span
    class=\"kc\">False</span><span class=\"p\">,</span>\n                <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='configure' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>configure <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">configure
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
    <span class=\"nf\">configure</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">def</span> <span class=\"nf\">get_plugins</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">try</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_plugins</span>\n                <span class=\"k\">except</span> <span
    class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_plugins</span><span
    class=\"p\">:</span> <span class=\"n\">Plugins</span> <span class=\"o\">=</span>
    <span class=\"n\">Plugins</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">return</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_plugins</span>\n\n
    \           <span class=\"kn\">from</span> <span class=\"nn\">markata</span> <span
    class=\"kn\">import</span> <span class=\"n\">Markata</span>\n\n            <span
    class=\"n\">Markata</span><span class=\"o\">.</span><span class=\"n\">plugins</span>
    <span class=\"o\">=</span> <span class=\"nb\">property</span><span class=\"p\">(</span><span
    class=\"n\">get_plugins</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='__init__' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'><strong>init</strong> <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\"><strong>init</strong>
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
    <span class=\"fm\">__init__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">markata</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span> <span class=\"o\">=</span> <span class=\"n\">markata</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='__rich__' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'><strong>rich</strong> <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\"><strong>rich</strong>
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
    <span class=\"nf\">__rich__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Panel</span><span
    class=\"p\">:</span>\n                <span class=\"n\">grid</span> <span class=\"o\">=</span>
    <span class=\"n\">Table</span><span class=\"o\">.</span><span class=\"n\">grid</span><span
    class=\"p\">(</span><span class=\"n\">expand</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n                <span class=\"n\">num_plugins</span>
    <span class=\"o\">=</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;[bright_blue](</span><span
    class=\"si\">{</span><span class=\"nb\">len</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">m</span><span
    class=\"o\">.</span><span class=\"n\">_pm</span><span class=\"o\">.</span><span
    class=\"n\">get_plugins</span><span class=\"p\">())</span><span class=\"si\">}</span><span
    class=\"s2\">)[/]&quot;</span>\n\n                <span class=\"k\">for</span>
    <span class=\"n\">plugin</span> <span class=\"ow\">in</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">m</span><span class=\"o\">.</span><span
    class=\"n\">_pm</span><span class=\"o\">.</span><span class=\"n\">get_plugins</span><span
    class=\"p\">():</span>\n                    <span class=\"n\">grid</span><span
    class=\"o\">.</span><span class=\"n\">add_row</span><span class=\"p\">(</span>\n
    \                       <span class=\"s2\">&quot;&quot;</span><span class=\"o\">.</span><span
    class=\"n\">join</span><span class=\"p\">(</span>\n                            <span
    class=\"p\">[</span>\n                                <span class=\"s2\">&quot;[bright_black]&quot;</span><span
    class=\"p\">,</span>\n                                <span class=\"s2\">&quot;.&quot;</span><span
    class=\"o\">.</span><span class=\"n\">join</span><span class=\"p\">(</span><span
    class=\"n\">plugin</span><span class=\"o\">.</span><span class=\"vm\">__name__</span><span
    class=\"o\">.</span><span class=\"n\">split</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;.&quot;</span><span class=\"p\">)[:</span><span class=\"o\">-</span><span
    class=\"mi\">1</span><span class=\"p\">]),</span>\n                                <span
    class=\"s2\">&quot;.[/]&quot;</span><span class=\"p\">,</span>\n                                <span
    class=\"n\">plugin</span><span class=\"o\">.</span><span class=\"vm\">__name__</span><span
    class=\"o\">.</span><span class=\"n\">split</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;.&quot;</span><span class=\"p\">)[</span><span class=\"o\">-</span><span
    class=\"mi\">1</span><span class=\"p\">],</span>\n                            <span
    class=\"p\">],</span>\n                        <span class=\"p\">),</span>\n                    <span
    class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"n\">Panel</span><span
    class=\"p\">(</span>\n                    <span class=\"n\">grid</span><span class=\"p\">,</span>\n
    \                   <span class=\"n\">title</span><span class=\"o\">=</span><span
    class=\"sa\">f</span><span class=\"s2\">&quot;plugins </span><span class=\"si\">{</span><span
    class=\"n\">num_plugins</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"n\">border_style</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;gold1&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"n\">expand</span><span class=\"o\">=</span><span
    class=\"kc\">False</span><span class=\"p\">,</span>\n                <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='get_plugins' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>get_plugins <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">get_plugins
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
    <span class=\"nf\">get_plugins</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">):</span>\n                <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_plugins</span>\n                <span class=\"k\">except</span>
    <span class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_plugins</span><span
    class=\"p\">:</span> <span class=\"n\">Plugins</span> <span class=\"o\">=</span>
    <span class=\"n\">Plugins</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">return</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_plugins</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>    </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Plugins.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"None ! ???+ source  ! ???+ source  ! ???+ source
    \ ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Plugins.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None ! ???+ source  ! ???+ source  !
    ???+ source  ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<article style=\"text-align: center;\">\n    <style>\n        section
    {\n            font-size: 200%;\n        }\n\n\n        .edit {\n            display:
    none;\n        }\n    </style>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Plugins.Py\n    </h1>\n</section></article>\n     </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Plugins.Py\n    </h1>\n</section>    <section class=\"body\">\n        <p>None</p>\n<p>!!
    class <h2 id='Plugins' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>Plugins <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">Plugins
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
    <span class=\"nc\">Plugins</span><span class=\"p\">:</span>\n            <span
    class=\"k\">def</span> <span class=\"fm\">__init__</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span> <span class=\"o\">=</span> <span class=\"n\">markata</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"nf\">__rich__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Panel</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">grid</span> <span class=\"o\">=</span> <span
    class=\"n\">Table</span><span class=\"o\">.</span><span class=\"n\">grid</span><span
    class=\"p\">(</span><span class=\"n\">expand</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n                <span class=\"n\">num_plugins</span>
    <span class=\"o\">=</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;[bright_blue](</span><span
    class=\"si\">{</span><span class=\"nb\">len</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">m</span><span
    class=\"o\">.</span><span class=\"n\">_pm</span><span class=\"o\">.</span><span
    class=\"n\">get_plugins</span><span class=\"p\">())</span><span class=\"si\">}</span><span
    class=\"s2\">)[/]&quot;</span>\n\n                <span class=\"k\">for</span>
    <span class=\"n\">plugin</span> <span class=\"ow\">in</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">m</span><span class=\"o\">.</span><span
    class=\"n\">_pm</span><span class=\"o\">.</span><span class=\"n\">get_plugins</span><span
    class=\"p\">():</span>\n                    <span class=\"n\">grid</span><span
    class=\"o\">.</span><span class=\"n\">add_row</span><span class=\"p\">(</span>\n
    \                       <span class=\"s2\">&quot;&quot;</span><span class=\"o\">.</span><span
    class=\"n\">join</span><span class=\"p\">(</span>\n                            <span
    class=\"p\">[</span>\n                                <span class=\"s2\">&quot;[bright_black]&quot;</span><span
    class=\"p\">,</span>\n                                <span class=\"s2\">&quot;.&quot;</span><span
    class=\"o\">.</span><span class=\"n\">join</span><span class=\"p\">(</span><span
    class=\"n\">plugin</span><span class=\"o\">.</span><span class=\"vm\">__name__</span><span
    class=\"o\">.</span><span class=\"n\">split</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;.&quot;</span><span class=\"p\">)[:</span><span class=\"o\">-</span><span
    class=\"mi\">1</span><span class=\"p\">]),</span>\n                                <span
    class=\"s2\">&quot;.[/]&quot;</span><span class=\"p\">,</span>\n                                <span
    class=\"n\">plugin</span><span class=\"o\">.</span><span class=\"vm\">__name__</span><span
    class=\"o\">.</span><span class=\"n\">split</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;.&quot;</span><span class=\"p\">)[</span><span class=\"o\">-</span><span
    class=\"mi\">1</span><span class=\"p\">],</span>\n                            <span
    class=\"p\">],</span>\n                        <span class=\"p\">),</span>\n                    <span
    class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"n\">Panel</span><span
    class=\"p\">(</span>\n                    <span class=\"n\">grid</span><span class=\"p\">,</span>\n
    \                   <span class=\"n\">title</span><span class=\"o\">=</span><span
    class=\"sa\">f</span><span class=\"s2\">&quot;plugins </span><span class=\"si\">{</span><span
    class=\"n\">num_plugins</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"n\">border_style</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;gold1&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"n\">expand</span><span class=\"o\">=</span><span
    class=\"kc\">False</span><span class=\"p\">,</span>\n                <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='configure' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>configure <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">configure
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
    <span class=\"nf\">configure</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">def</span> <span class=\"nf\">get_plugins</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">try</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_plugins</span>\n                <span class=\"k\">except</span> <span
    class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_plugins</span><span
    class=\"p\">:</span> <span class=\"n\">Plugins</span> <span class=\"o\">=</span>
    <span class=\"n\">Plugins</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">return</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_plugins</span>\n\n
    \           <span class=\"kn\">from</span> <span class=\"nn\">markata</span> <span
    class=\"kn\">import</span> <span class=\"n\">Markata</span>\n\n            <span
    class=\"n\">Markata</span><span class=\"o\">.</span><span class=\"n\">plugins</span>
    <span class=\"o\">=</span> <span class=\"nb\">property</span><span class=\"p\">(</span><span
    class=\"n\">get_plugins</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='__init__' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'><strong>init</strong> <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\"><strong>init</strong>
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
    <span class=\"fm\">__init__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">markata</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span> <span class=\"o\">=</span> <span class=\"n\">markata</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='__rich__' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'><strong>rich</strong> <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\"><strong>rich</strong>
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
    <span class=\"nf\">__rich__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Panel</span><span
    class=\"p\">:</span>\n                <span class=\"n\">grid</span> <span class=\"o\">=</span>
    <span class=\"n\">Table</span><span class=\"o\">.</span><span class=\"n\">grid</span><span
    class=\"p\">(</span><span class=\"n\">expand</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n                <span class=\"n\">num_plugins</span>
    <span class=\"o\">=</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;[bright_blue](</span><span
    class=\"si\">{</span><span class=\"nb\">len</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">m</span><span
    class=\"o\">.</span><span class=\"n\">_pm</span><span class=\"o\">.</span><span
    class=\"n\">get_plugins</span><span class=\"p\">())</span><span class=\"si\">}</span><span
    class=\"s2\">)[/]&quot;</span>\n\n                <span class=\"k\">for</span>
    <span class=\"n\">plugin</span> <span class=\"ow\">in</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">m</span><span class=\"o\">.</span><span
    class=\"n\">_pm</span><span class=\"o\">.</span><span class=\"n\">get_plugins</span><span
    class=\"p\">():</span>\n                    <span class=\"n\">grid</span><span
    class=\"o\">.</span><span class=\"n\">add_row</span><span class=\"p\">(</span>\n
    \                       <span class=\"s2\">&quot;&quot;</span><span class=\"o\">.</span><span
    class=\"n\">join</span><span class=\"p\">(</span>\n                            <span
    class=\"p\">[</span>\n                                <span class=\"s2\">&quot;[bright_black]&quot;</span><span
    class=\"p\">,</span>\n                                <span class=\"s2\">&quot;.&quot;</span><span
    class=\"o\">.</span><span class=\"n\">join</span><span class=\"p\">(</span><span
    class=\"n\">plugin</span><span class=\"o\">.</span><span class=\"vm\">__name__</span><span
    class=\"o\">.</span><span class=\"n\">split</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;.&quot;</span><span class=\"p\">)[:</span><span class=\"o\">-</span><span
    class=\"mi\">1</span><span class=\"p\">]),</span>\n                                <span
    class=\"s2\">&quot;.[/]&quot;</span><span class=\"p\">,</span>\n                                <span
    class=\"n\">plugin</span><span class=\"o\">.</span><span class=\"vm\">__name__</span><span
    class=\"o\">.</span><span class=\"n\">split</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;.&quot;</span><span class=\"p\">)[</span><span class=\"o\">-</span><span
    class=\"mi\">1</span><span class=\"p\">],</span>\n                            <span
    class=\"p\">],</span>\n                        <span class=\"p\">),</span>\n                    <span
    class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"n\">Panel</span><span
    class=\"p\">(</span>\n                    <span class=\"n\">grid</span><span class=\"p\">,</span>\n
    \                   <span class=\"n\">title</span><span class=\"o\">=</span><span
    class=\"sa\">f</span><span class=\"s2\">&quot;plugins </span><span class=\"si\">{</span><span
    class=\"n\">num_plugins</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"n\">border_style</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;gold1&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"n\">expand</span><span class=\"o\">=</span><span
    class=\"kc\">False</span><span class=\"p\">,</span>\n                <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='get_plugins' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>get_plugins <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">get_plugins
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
    <span class=\"nf\">get_plugins</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">):</span>\n                <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_plugins</span>\n                <span class=\"k\">except</span>
    <span class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_plugins</span><span
    class=\"p\">:</span> <span class=\"n\">Plugins</span> <span class=\"o\">=</span>
    <span class=\"n\">Plugins</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">return</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_plugins</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>"
  raw.md: ''
published: true
slug: markata/cli/plugins
title: Plugins.Py


---

None


!! class <h2 id='Plugins' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Plugins <em class='small'>class</em></h2>

???+ source "Plugins <em class='small'>source</em>"

```python

        class Plugins:
            def __init__(self, markata: "Markata") -> None:
                self.m = markata

            def __rich__(self) -> Panel:
                grid = Table.grid(expand=True)
                num_plugins = f"[bright_blue]({len(self.m._pm.get_plugins())})[/]"

                for plugin in self.m._pm.get_plugins():
                    grid.add_row(
                        "".join(
                            [
                                "[bright_black]",
                                ".".join(plugin.__name__.split(".")[:-1]),
                                ".[/]",
                                plugin.__name__.split(".")[-1],
                            ],
                        ),
                    )
                return Panel(
                    grid,
                    title=f"plugins {num_plugins}",
                    border_style="gold1",
                    expand=False,
                )
```


!! function <h2 id='configure' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>configure <em class='small'>function</em></h2>

???+ source "configure <em class='small'>source</em>"

```python

        def configure(markata: "Markata") -> None:
            def get_plugins(self):
                try:
                    return self._plugins
                except AttributeError:
                    self._plugins: Plugins = Plugins(self)
                    return self._plugins

            from markata import Markata

            Markata.plugins = property(get_plugins)
```


!! method <h2 id='__init__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__init__ <em class='small'>method</em></h2>

???+ source "__init__ <em class='small'>source</em>"

```python

        def __init__(self, markata: "Markata") -> None:
                self.m = markata
```


!! method <h2 id='__rich__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__rich__ <em class='small'>method</em></h2>

???+ source "__rich__ <em class='small'>source</em>"

```python

        def __rich__(self) -> Panel:
                grid = Table.grid(expand=True)
                num_plugins = f"[bright_blue]({len(self.m._pm.get_plugins())})[/]"

                for plugin in self.m._pm.get_plugins():
                    grid.add_row(
                        "".join(
                            [
                                "[bright_black]",
                                ".".join(plugin.__name__.split(".")[:-1]),
                                ".[/]",
                                plugin.__name__.split(".")[-1],
                            ],
                        ),
                    )
                return Panel(
                    grid,
                    title=f"plugins {num_plugins}",
                    border_style="gold1",
                    expand=False,
                )
```


!! function <h2 id='get_plugins' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get_plugins <em class='small'>function</em></h2>

???+ source "get_plugins <em class='small'>source</em>"

```python

        def get_plugins(self):
                try:
                    return self._plugins
                except AttributeError:
                    self._plugins: Plugins = Plugins(self)
                    return self._plugins
```

