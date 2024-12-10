---
content: "None\n\n\n!! class <h2 id='MarkataWidget' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>MarkataWidget <em class='small'>class</em></h2>\n\n???+ source \"MarkataWidget
  <em class='small'>source</em>\"\n\n```python\n\n        class MarkataWidget(Widget):\n
  \           def __init__(self, markata: Markata, widget: str = \"server\") -> None:\n
  \               super().__init__(widget)\n                self.m = markata\n                self.widget
  = widget\n                self.renderable = getattr(self.m, self.widget)\n\n            def
  render(self):\n                return self.renderable\n\n            async def update(self,
  renderable: RenderableType) -> None:\n                self.renderable = renderable\n
  \               self.refresh()\n```\n\n\n!! class <h2 id='MarkataApp' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>MarkataApp <em class='small'>class</em></h2>\n\n???+
  source \"MarkataApp <em class='small'>source</em>\"\n\n```python\n\n        class
  MarkataApp(App):\n            async def on_mount(self) -> None:\n                self.m
  = Markata()\n                self.server = MarkataWidget(self.m, \"server\")\n                self.runner
  = MarkataWidget(self.m, \"runner\")\n                self.plugins = MarkataWidget(self.m,
  \"plugins\")\n                self.summary = MarkataWidget(self.m, \"summary\")\n
  \               await self.view.dock(self.plugins, edge=\"left\", size=30, name=\"plugins\")\n
  \               await self.view.dock(self.summary, edge=\"right\", size=30, name=\"summary\")\n
  \               await self.view.dock(self.server, self.runner, edge=\"top\")\n                self.set_interval(1,
  self.action_refresh)\n\n            async def on_load(self, event):\n                await
  self.bind(\"q\", \"my_quit\")\n                await self.bind(\"r\", \"refresh\")\n\n
  \           async def action_my_quit(self) -> None:\n                await self.action_quit()\n\n
  \           async def action_refresh(self) -> None:\n                self.refresh()\n
  \               self.runner.refresh()\n                self.server.refresh()\n                self.plugins.refresh()\n
  \               self.summary.refresh()\n```\n\n\n!! method <h2 id='__init__' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>__init__ <em class='small'>method</em></h2>\n\n???+
  source \"__init__ <em class='small'>source</em>\"\n\n```python\n\n        def __init__(self,
  markata: Markata, widget: str = \"server\") -> None:\n                super().__init__(widget)\n
  \               self.m = markata\n                self.widget = widget\n                self.renderable
  = getattr(self.m, self.widget)\n```\n\n\n!! method <h2 id='render' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>render <em class='small'>method</em></h2>\n\n???+
  source \"render <em class='small'>source</em>\"\n\n```python\n\n        def render(self):\n
  \               return self.renderable\n```\n\n"
date: 0001-01-01
description: 'None ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source '
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>__Main__.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"None ! ???+ source  ! ???+ source  ! ???+ source
    \ ! ???+ source \" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>__Main__.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None ! ???+ source  ! ???+ source  !
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        __Main__.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <p>None</p>\n<p>!!
    class <h2 id='MarkataWidget' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>MarkataWidget <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">MarkataWidget
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
    <span class=\"nc\">MarkataWidget</span><span class=\"p\">(</span><span class=\"n\">Widget</span><span
    class=\"p\">):</span>\n            <span class=\"k\">def</span> <span class=\"fm\">__init__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"n\">Markata</span><span
    class=\"p\">,</span> <span class=\"n\">widget</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;server&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"nb\">super</span><span class=\"p\">()</span><span
    class=\"o\">.</span><span class=\"fm\">__init__</span><span class=\"p\">(</span><span
    class=\"n\">widget</span><span class=\"p\">)</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">m</span> <span class=\"o\">=</span> <span
    class=\"n\">markata</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">widget</span> <span class=\"o\">=</span>
    <span class=\"n\">widget</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">renderable</span> <span class=\"o\">=</span>
    <span class=\"nb\">getattr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">m</span><span class=\"p\">,</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">widget</span><span
    class=\"p\">)</span>\n\n            <span class=\"k\">def</span> <span class=\"nf\">render</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">renderable</span>\n\n            <span class=\"k\">async</span>
    <span class=\"k\">def</span> <span class=\"nf\">update</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"n\">renderable</span><span
    class=\"p\">:</span> <span class=\"n\">RenderableType</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">renderable</span> <span class=\"o\">=</span> <span class=\"n\">renderable</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">refresh</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='MarkataApp' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>MarkataApp <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">MarkataApp
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
    <span class=\"nc\">MarkataApp</span><span class=\"p\">(</span><span class=\"n\">App</span><span
    class=\"p\">):</span>\n            <span class=\"k\">async</span> <span class=\"k\">def</span>
    <span class=\"nf\">on_mount</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span> <span class=\"o\">=</span> <span class=\"n\">Markata</span><span
    class=\"p\">()</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">server</span> <span class=\"o\">=</span> <span class=\"n\">MarkataWidget</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"p\">,</span> <span class=\"s2\">&quot;server&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">runner</span> <span class=\"o\">=</span> <span class=\"n\">MarkataWidget</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"p\">,</span> <span class=\"s2\">&quot;runner&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">plugins</span> <span class=\"o\">=</span> <span class=\"n\">MarkataWidget</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"p\">,</span> <span class=\"s2\">&quot;plugins&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">summary</span> <span class=\"o\">=</span> <span class=\"n\">MarkataWidget</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"p\">,</span> <span class=\"s2\">&quot;summary&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"k\">await</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">view</span><span class=\"o\">.</span><span
    class=\"n\">dock</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">plugins</span><span class=\"p\">,</span>
    <span class=\"n\">edge</span><span class=\"o\">=</span><span class=\"s2\">&quot;left&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">size</span><span class=\"o\">=</span><span
    class=\"mi\">30</span><span class=\"p\">,</span> <span class=\"n\">name</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;plugins&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">await</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">view</span><span class=\"o\">.</span><span
    class=\"n\">dock</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">summary</span><span class=\"p\">,</span>
    <span class=\"n\">edge</span><span class=\"o\">=</span><span class=\"s2\">&quot;right&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">size</span><span class=\"o\">=</span><span
    class=\"mi\">30</span><span class=\"p\">,</span> <span class=\"n\">name</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;summary&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">await</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">view</span><span class=\"o\">.</span><span
    class=\"n\">dock</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">server</span><span class=\"p\">,</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">runner</span><span
    class=\"p\">,</span> <span class=\"n\">edge</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;top&quot;</span><span class=\"p\">)</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">set_interval</span><span
    class=\"p\">(</span><span class=\"mi\">1</span><span class=\"p\">,</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">action_refresh</span><span
    class=\"p\">)</span>\n\n            <span class=\"k\">async</span> <span class=\"k\">def</span>
    <span class=\"nf\">on_load</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">event</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">await</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">bind</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;q&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;my_quit&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"k\">await</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">bind</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;r&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;refresh&quot;</span><span
    class=\"p\">)</span>\n\n            <span class=\"k\">async</span> <span class=\"k\">def</span>
    <span class=\"nf\">action_my_quit</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"k\">await</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">action_quit</span><span class=\"p\">()</span>\n\n
    \           <span class=\"k\">async</span> <span class=\"k\">def</span> <span
    class=\"nf\">action_refresh</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">refresh</span><span class=\"p\">()</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">runner</span><span class=\"o\">.</span><span
    class=\"n\">refresh</span><span class=\"p\">()</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">server</span><span class=\"o\">.</span><span
    class=\"n\">refresh</span><span class=\"p\">()</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">plugins</span><span class=\"o\">.</span><span
    class=\"n\">refresh</span><span class=\"p\">()</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">summary</span><span class=\"o\">.</span><span
    class=\"n\">refresh</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    <span class=\"n\">Markata</span><span class=\"p\">,</span> <span class=\"n\">widget</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;server&quot;</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
    class=\"nb\">super</span><span class=\"p\">()</span><span class=\"o\">.</span><span
    class=\"fm\">__init__</span><span class=\"p\">(</span><span class=\"n\">widget</span><span
    class=\"p\">)</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span> <span class=\"o\">=</span> <span class=\"n\">markata</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">widget</span> <span class=\"o\">=</span> <span class=\"n\">widget</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">renderable</span> <span class=\"o\">=</span> <span class=\"nb\">getattr</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"p\">,</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">widget</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='render' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>render <em class='small'>method</em></h2></p>\n<div class=\"admonition
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
    <span class=\"nf\">render</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">):</span>\n                <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">renderable</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>    </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>__Main__.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"None ! ???+ source  ! ???+ source  ! ???+ source
    \ ! ???+ source \" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>__Main__.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None ! ???+ source  ! ???+ source  !
    ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    </head>\n    <body>\n<article style=\"text-align:
    center;\">\n    <style>\n        section {\n            font-size: 200%;\n        }\n\n\n
    \       .edit {\n            display: none;\n        }\n    </style>\n<section
    class=\"title\">\n    <h1 id=\"title\">\n        __Main__.Py\n    </h1>\n</section></article>\n
    \    </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       __Main__.Py\n    </h1>\n</section>    <section class=\"body\">\n        <p>None</p>\n<p>!!
    class <h2 id='MarkataWidget' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>MarkataWidget <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">MarkataWidget
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
    <span class=\"nc\">MarkataWidget</span><span class=\"p\">(</span><span class=\"n\">Widget</span><span
    class=\"p\">):</span>\n            <span class=\"k\">def</span> <span class=\"fm\">__init__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"n\">Markata</span><span
    class=\"p\">,</span> <span class=\"n\">widget</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;server&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"nb\">super</span><span class=\"p\">()</span><span
    class=\"o\">.</span><span class=\"fm\">__init__</span><span class=\"p\">(</span><span
    class=\"n\">widget</span><span class=\"p\">)</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">m</span> <span class=\"o\">=</span> <span
    class=\"n\">markata</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">widget</span> <span class=\"o\">=</span>
    <span class=\"n\">widget</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">renderable</span> <span class=\"o\">=</span>
    <span class=\"nb\">getattr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">m</span><span class=\"p\">,</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">widget</span><span
    class=\"p\">)</span>\n\n            <span class=\"k\">def</span> <span class=\"nf\">render</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">renderable</span>\n\n            <span class=\"k\">async</span>
    <span class=\"k\">def</span> <span class=\"nf\">update</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"n\">renderable</span><span
    class=\"p\">:</span> <span class=\"n\">RenderableType</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">renderable</span> <span class=\"o\">=</span> <span class=\"n\">renderable</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">refresh</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='MarkataApp' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>MarkataApp <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">MarkataApp
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
    <span class=\"nc\">MarkataApp</span><span class=\"p\">(</span><span class=\"n\">App</span><span
    class=\"p\">):</span>\n            <span class=\"k\">async</span> <span class=\"k\">def</span>
    <span class=\"nf\">on_mount</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span> <span class=\"o\">=</span> <span class=\"n\">Markata</span><span
    class=\"p\">()</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">server</span> <span class=\"o\">=</span> <span class=\"n\">MarkataWidget</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"p\">,</span> <span class=\"s2\">&quot;server&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">runner</span> <span class=\"o\">=</span> <span class=\"n\">MarkataWidget</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"p\">,</span> <span class=\"s2\">&quot;runner&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">plugins</span> <span class=\"o\">=</span> <span class=\"n\">MarkataWidget</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"p\">,</span> <span class=\"s2\">&quot;plugins&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">summary</span> <span class=\"o\">=</span> <span class=\"n\">MarkataWidget</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"p\">,</span> <span class=\"s2\">&quot;summary&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"k\">await</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">view</span><span class=\"o\">.</span><span
    class=\"n\">dock</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">plugins</span><span class=\"p\">,</span>
    <span class=\"n\">edge</span><span class=\"o\">=</span><span class=\"s2\">&quot;left&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">size</span><span class=\"o\">=</span><span
    class=\"mi\">30</span><span class=\"p\">,</span> <span class=\"n\">name</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;plugins&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">await</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">view</span><span class=\"o\">.</span><span
    class=\"n\">dock</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">summary</span><span class=\"p\">,</span>
    <span class=\"n\">edge</span><span class=\"o\">=</span><span class=\"s2\">&quot;right&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">size</span><span class=\"o\">=</span><span
    class=\"mi\">30</span><span class=\"p\">,</span> <span class=\"n\">name</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;summary&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">await</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">view</span><span class=\"o\">.</span><span
    class=\"n\">dock</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">server</span><span class=\"p\">,</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">runner</span><span
    class=\"p\">,</span> <span class=\"n\">edge</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;top&quot;</span><span class=\"p\">)</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">set_interval</span><span
    class=\"p\">(</span><span class=\"mi\">1</span><span class=\"p\">,</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">action_refresh</span><span
    class=\"p\">)</span>\n\n            <span class=\"k\">async</span> <span class=\"k\">def</span>
    <span class=\"nf\">on_load</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">event</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">await</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">bind</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;q&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;my_quit&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"k\">await</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">bind</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;r&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;refresh&quot;</span><span
    class=\"p\">)</span>\n\n            <span class=\"k\">async</span> <span class=\"k\">def</span>
    <span class=\"nf\">action_my_quit</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"k\">await</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">action_quit</span><span class=\"p\">()</span>\n\n
    \           <span class=\"k\">async</span> <span class=\"k\">def</span> <span
    class=\"nf\">action_refresh</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">refresh</span><span class=\"p\">()</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">runner</span><span class=\"o\">.</span><span
    class=\"n\">refresh</span><span class=\"p\">()</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">server</span><span class=\"o\">.</span><span
    class=\"n\">refresh</span><span class=\"p\">()</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">plugins</span><span class=\"o\">.</span><span
    class=\"n\">refresh</span><span class=\"p\">()</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">summary</span><span class=\"o\">.</span><span
    class=\"n\">refresh</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    <span class=\"n\">Markata</span><span class=\"p\">,</span> <span class=\"n\">widget</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;server&quot;</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
    class=\"nb\">super</span><span class=\"p\">()</span><span class=\"o\">.</span><span
    class=\"fm\">__init__</span><span class=\"p\">(</span><span class=\"n\">widget</span><span
    class=\"p\">)</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span> <span class=\"o\">=</span> <span class=\"n\">markata</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">widget</span> <span class=\"o\">=</span> <span class=\"n\">widget</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">renderable</span> <span class=\"o\">=</span> <span class=\"nb\">getattr</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"p\">,</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">widget</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='render' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>render <em class='small'>method</em></h2></p>\n<div class=\"admonition
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
    <span class=\"nf\">render</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">):</span>\n                <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">renderable</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>"
  raw.md: ''
published: true
slug: markata/cli/main
title: __Main__.Py


---

None


!! class <h2 id='MarkataWidget' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>MarkataWidget <em class='small'>class</em></h2>

???+ source "MarkataWidget <em class='small'>source</em>"

```python

        class MarkataWidget(Widget):
            def __init__(self, markata: Markata, widget: str = "server") -> None:
                super().__init__(widget)
                self.m = markata
                self.widget = widget
                self.renderable = getattr(self.m, self.widget)

            def render(self):
                return self.renderable

            async def update(self, renderable: RenderableType) -> None:
                self.renderable = renderable
                self.refresh()
```


!! class <h2 id='MarkataApp' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>MarkataApp <em class='small'>class</em></h2>

???+ source "MarkataApp <em class='small'>source</em>"

```python

        class MarkataApp(App):
            async def on_mount(self) -> None:
                self.m = Markata()
                self.server = MarkataWidget(self.m, "server")
                self.runner = MarkataWidget(self.m, "runner")
                self.plugins = MarkataWidget(self.m, "plugins")
                self.summary = MarkataWidget(self.m, "summary")
                await self.view.dock(self.plugins, edge="left", size=30, name="plugins")
                await self.view.dock(self.summary, edge="right", size=30, name="summary")
                await self.view.dock(self.server, self.runner, edge="top")
                self.set_interval(1, self.action_refresh)

            async def on_load(self, event):
                await self.bind("q", "my_quit")
                await self.bind("r", "refresh")

            async def action_my_quit(self) -> None:
                await self.action_quit()

            async def action_refresh(self) -> None:
                self.refresh()
                self.runner.refresh()
                self.server.refresh()
                self.plugins.refresh()
                self.summary.refresh()
```


!! method <h2 id='__init__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__init__ <em class='small'>method</em></h2>

???+ source "__init__ <em class='small'>source</em>"

```python

        def __init__(self, markata: Markata, widget: str = "server") -> None:
                super().__init__(widget)
                self.m = markata
                self.widget = widget
                self.renderable = getattr(self.m, self.widget)
```


!! method <h2 id='render' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>render <em class='small'>method</em></h2>

???+ source "render <em class='small'>source</em>"

```python

        def render(self):
                return self.renderable
```

