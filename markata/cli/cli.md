---
content: "None\n\n\n!! function <h2 id='make_layout' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>make_layout <em class='small'>function</em></h2>\n    Define the layout.\n???+
  source \"make_layout <em class='small'>source</em>\"\n\n```python\n\n        def
  make_layout() -> Layout:\n            \"\"\"Define the layout.\"\"\"\n            layout
  = Layout(name=\"root\")\n\n            layout.split(\n                Layout(name=\"header\",
  size=3),\n                Layout(name=\"main\"),\n            )\n            layout[\"main\"].split_row(\n
  \               Layout(name=\"side\", ratio=50),\n                Layout(name=\"mid\",
  ratio=30),\n                Layout(name=\"describe\", ratio=20),\n            )\n
  \           layout[\"mid\"].split(\n                Layout(name=\"server\"),\n                Layout(name=\"runner\"),\n
  \           )\n            layout[\"side\"].split(\n                Layout(name=\"plugins\"),\n
  \           )\n            return layout\n```\n\n\n!! function <h2 id='run_until_keyboard_interrupt'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>run_until_keyboard_interrupt
  <em class='small'>function</em></h2>\n\n???+ source \"run_until_keyboard_interrupt
  <em class='small'>source</em>\"\n\n```python\n\n        def run_until_keyboard_interrupt()
  -> None:\n            try:\n                while True:\n                    time.sleep(0.2)\n
  \           except KeyboardInterrupt:\n                pass\n```\n\n\n!! function
  <h2 id='version_callback' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>version_callback <em class='small'>function</em></h2>\n\n???+ source \"version_callback
  <em class='small'>source</em>\"\n\n```python\n\n        def version_callback(value:
  bool) -> None:\n            if value:\n                from markata import __version__\n\n
  \               typer.echo(f\"Markata CLI Version: {__version__}\")\n                raise
  typer.Exit\n```\n\n\n!! function <h2 id='json_callback' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>json_callback <em class='small'>function</em></h2>\n\n???+
  source \"json_callback <em class='small'>source</em>\"\n\n```python\n\n        def
  json_callback(value: bool) -> None:\n            if value:\n                from
  markata import Markata\n\n                typer.echo(Markata().to_json())\n                raise
  typer.Exit\n```\n\n\n!! function <h2 id='main' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>main <em class='small'>function</em></h2>\n\n???+ source \"main <em class='small'>source</em>\"\n\n```python\n\n
  \       def main(\n            version: bool = typer.Option(\n                None,\n
  \               \"--version\",\n                callback=version_callback,\n                is_eager=True,\n
  \           ),\n            to_json: bool = typer.Option(\n                None,\n
  \               \"--to-json\",\n                callback=json_callback,\n                is_eager=True,\n
  \           ),\n        ) -> None:\n            # Do other global stuff, handle
  other global options here\n            return\n```\n\n\n!! function <h2 id='cli'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>cli <em class='small'>function</em></h2>\n\n???+
  source \"cli <em class='small'>source</em>\"\n\n```python\n\n        def cli() ->
  None:\n            from markata import Markata\n\n            m = Markata()\n            m._pm.hook.cli(markata=m,
  app=app)\n            app()\n```\n\n"
date: 0001-01-01
description: 'None ! ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  !
  ???+ source '
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Cli.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"None ! ! ???+ source  ! ???+ source  ! ???+ source
    \ ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Cli.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None ! ! ???+ source  ! ???+ source  !
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Cli.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <p>None</p>\n<p>!!
    function <h2 id='make_layout' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>make_layout <em class='small'>function</em></h2>\nDefine the layout.</p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">make_layout
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
    <span class=\"nf\">make_layout</span><span class=\"p\">()</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Layout</span><span class=\"p\">:</span>\n<span class=\"w\">
    \           </span><span class=\"sd\">&quot;&quot;&quot;Define the layout.&quot;&quot;&quot;</span>\n
    \           <span class=\"n\">layout</span> <span class=\"o\">=</span> <span class=\"n\">Layout</span><span
    class=\"p\">(</span><span class=\"n\">name</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;root&quot;</span><span class=\"p\">)</span>\n\n            <span
    class=\"n\">layout</span><span class=\"o\">.</span><span class=\"n\">split</span><span
    class=\"p\">(</span>\n                <span class=\"n\">Layout</span><span class=\"p\">(</span><span
    class=\"n\">name</span><span class=\"o\">=</span><span class=\"s2\">&quot;header&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">size</span><span class=\"o\">=</span><span
    class=\"mi\">3</span><span class=\"p\">),</span>\n                <span class=\"n\">Layout</span><span
    class=\"p\">(</span><span class=\"n\">name</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;main&quot;</span><span class=\"p\">),</span>\n            <span
    class=\"p\">)</span>\n            <span class=\"n\">layout</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;main&quot;</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">split_row</span><span class=\"p\">(</span>\n                <span
    class=\"n\">Layout</span><span class=\"p\">(</span><span class=\"n\">name</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;side&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">ratio</span><span class=\"o\">=</span><span class=\"mi\">50</span><span
    class=\"p\">),</span>\n                <span class=\"n\">Layout</span><span class=\"p\">(</span><span
    class=\"n\">name</span><span class=\"o\">=</span><span class=\"s2\">&quot;mid&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">ratio</span><span class=\"o\">=</span><span
    class=\"mi\">30</span><span class=\"p\">),</span>\n                <span class=\"n\">Layout</span><span
    class=\"p\">(</span><span class=\"n\">name</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;describe&quot;</span><span class=\"p\">,</span> <span class=\"n\">ratio</span><span
    class=\"o\">=</span><span class=\"mi\">20</span><span class=\"p\">),</span>\n
    \           <span class=\"p\">)</span>\n            <span class=\"n\">layout</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;mid&quot;</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">split</span><span class=\"p\">(</span>\n
    \               <span class=\"n\">Layout</span><span class=\"p\">(</span><span
    class=\"n\">name</span><span class=\"o\">=</span><span class=\"s2\">&quot;server&quot;</span><span
    class=\"p\">),</span>\n                <span class=\"n\">Layout</span><span class=\"p\">(</span><span
    class=\"n\">name</span><span class=\"o\">=</span><span class=\"s2\">&quot;runner&quot;</span><span
    class=\"p\">),</span>\n            <span class=\"p\">)</span>\n            <span
    class=\"n\">layout</span><span class=\"p\">[</span><span class=\"s2\">&quot;side&quot;</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">split</span><span
    class=\"p\">(</span>\n                <span class=\"n\">Layout</span><span class=\"p\">(</span><span
    class=\"n\">name</span><span class=\"o\">=</span><span class=\"s2\">&quot;plugins&quot;</span><span
    class=\"p\">),</span>\n            <span class=\"p\">)</span>\n            <span
    class=\"k\">return</span> <span class=\"n\">layout</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='run_until_keyboard_interrupt' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>run_until_keyboard_interrupt <em class='small'>function</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">run_until_keyboard_interrupt
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
    <span class=\"nf\">run_until_keyboard_interrupt</span><span class=\"p\">()</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">try</span><span class=\"p\">:</span>\n                <span
    class=\"k\">while</span> <span class=\"kc\">True</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">time</span><span class=\"o\">.</span><span
    class=\"n\">sleep</span><span class=\"p\">(</span><span class=\"mf\">0.2</span><span
    class=\"p\">)</span>\n            <span class=\"k\">except</span> <span class=\"ne\">KeyboardInterrupt</span><span
    class=\"p\">:</span>\n                <span class=\"k\">pass</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='version_callback' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>version_callback <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">version_callback
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
    <span class=\"nf\">version_callback</span><span class=\"p\">(</span><span class=\"n\">value</span><span
    class=\"p\">:</span> <span class=\"nb\">bool</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">if</span> <span class=\"n\">value</span><span class=\"p\">:</span>\n
    \               <span class=\"kn\">from</span> <span class=\"nn\">markata</span>
    <span class=\"kn\">import</span> <span class=\"n\">__version__</span>\n\n                <span
    class=\"n\">typer</span><span class=\"o\">.</span><span class=\"n\">echo</span><span
    class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">&quot;Markata
    CLI Version: </span><span class=\"si\">{</span><span class=\"n\">__version__</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">raise</span> <span class=\"n\">typer</span><span
    class=\"o\">.</span><span class=\"n\">Exit</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='json_callback' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>json_callback <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">json_callback
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
    <span class=\"nf\">json_callback</span><span class=\"p\">(</span><span class=\"n\">value</span><span
    class=\"p\">:</span> <span class=\"nb\">bool</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">if</span> <span class=\"n\">value</span><span class=\"p\">:</span>\n
    \               <span class=\"kn\">from</span> <span class=\"nn\">markata</span>
    <span class=\"kn\">import</span> <span class=\"n\">Markata</span>\n\n                <span
    class=\"n\">typer</span><span class=\"o\">.</span><span class=\"n\">echo</span><span
    class=\"p\">(</span><span class=\"n\">Markata</span><span class=\"p\">()</span><span
    class=\"o\">.</span><span class=\"n\">to_json</span><span class=\"p\">())</span>\n
    \               <span class=\"k\">raise</span> <span class=\"n\">typer</span><span
    class=\"o\">.</span><span class=\"n\">Exit</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='main' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>main <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">main
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
    <span class=\"nf\">main</span><span class=\"p\">(</span>\n            <span class=\"n\">version</span><span
    class=\"p\">:</span> <span class=\"nb\">bool</span> <span class=\"o\">=</span>
    <span class=\"n\">typer</span><span class=\"o\">.</span><span class=\"n\">Option</span><span
    class=\"p\">(</span>\n                <span class=\"kc\">None</span><span class=\"p\">,</span>\n
    \               <span class=\"s2\">&quot;--version&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">callback</span><span class=\"o\">=</span><span
    class=\"n\">version_callback</span><span class=\"p\">,</span>\n                <span
    class=\"n\">is_eager</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span>\n            <span class=\"p\">),</span>\n            <span
    class=\"n\">to_json</span><span class=\"p\">:</span> <span class=\"nb\">bool</span>
    <span class=\"o\">=</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
    class=\"n\">Option</span><span class=\"p\">(</span>\n                <span class=\"kc\">None</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;--to-json&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">callback</span><span class=\"o\">=</span><span
    class=\"n\">json_callback</span><span class=\"p\">,</span>\n                <span
    class=\"n\">is_eager</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span>\n            <span class=\"p\">),</span>\n        <span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"c1\"># Do other global stuff, handle other global options
    here</span>\n            <span class=\"k\">return</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='cli' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>cli
    <em class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">cli <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">cli</span><span class=\"p\">()</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"kn\">from</span>
    <span class=\"nn\">markata</span> <span class=\"kn\">import</span> <span class=\"n\">Markata</span>\n\n
    \           <span class=\"n\">m</span> <span class=\"o\">=</span> <span class=\"n\">Markata</span><span
    class=\"p\">()</span>\n            <span class=\"n\">m</span><span class=\"o\">.</span><span
    class=\"n\">_pm</span><span class=\"o\">.</span><span class=\"n\">hook</span><span
    class=\"o\">.</span><span class=\"n\">cli</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"o\">=</span><span class=\"n\">m</span><span
    class=\"p\">,</span> <span class=\"n\">app</span><span class=\"o\">=</span><span
    class=\"n\">app</span><span class=\"p\">)</span>\n            <span class=\"n\">app</span><span
    class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>
    \   </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Cli.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"None ! ! ???+ source  ! ???+ source  ! ???+ source
    \ ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Cli.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None ! ! ???+ source  ! ???+ source  !
    ???+ source  ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<article style=\"text-align: center;\">\n    <style>\n        section
    {\n            font-size: 200%;\n        }\n\n\n        .edit {\n            display:
    none;\n        }\n    </style>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Cli.Py\n    </h1>\n</section></article>\n     </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Cli.Py\n    </h1>\n</section>    <section class=\"body\">\n        <p>None</p>\n<p>!!
    function <h2 id='make_layout' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>make_layout <em class='small'>function</em></h2>\nDefine the layout.</p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">make_layout
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
    <span class=\"nf\">make_layout</span><span class=\"p\">()</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Layout</span><span class=\"p\">:</span>\n<span class=\"w\">
    \           </span><span class=\"sd\">&quot;&quot;&quot;Define the layout.&quot;&quot;&quot;</span>\n
    \           <span class=\"n\">layout</span> <span class=\"o\">=</span> <span class=\"n\">Layout</span><span
    class=\"p\">(</span><span class=\"n\">name</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;root&quot;</span><span class=\"p\">)</span>\n\n            <span
    class=\"n\">layout</span><span class=\"o\">.</span><span class=\"n\">split</span><span
    class=\"p\">(</span>\n                <span class=\"n\">Layout</span><span class=\"p\">(</span><span
    class=\"n\">name</span><span class=\"o\">=</span><span class=\"s2\">&quot;header&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">size</span><span class=\"o\">=</span><span
    class=\"mi\">3</span><span class=\"p\">),</span>\n                <span class=\"n\">Layout</span><span
    class=\"p\">(</span><span class=\"n\">name</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;main&quot;</span><span class=\"p\">),</span>\n            <span
    class=\"p\">)</span>\n            <span class=\"n\">layout</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;main&quot;</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">split_row</span><span class=\"p\">(</span>\n                <span
    class=\"n\">Layout</span><span class=\"p\">(</span><span class=\"n\">name</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;side&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">ratio</span><span class=\"o\">=</span><span class=\"mi\">50</span><span
    class=\"p\">),</span>\n                <span class=\"n\">Layout</span><span class=\"p\">(</span><span
    class=\"n\">name</span><span class=\"o\">=</span><span class=\"s2\">&quot;mid&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">ratio</span><span class=\"o\">=</span><span
    class=\"mi\">30</span><span class=\"p\">),</span>\n                <span class=\"n\">Layout</span><span
    class=\"p\">(</span><span class=\"n\">name</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;describe&quot;</span><span class=\"p\">,</span> <span class=\"n\">ratio</span><span
    class=\"o\">=</span><span class=\"mi\">20</span><span class=\"p\">),</span>\n
    \           <span class=\"p\">)</span>\n            <span class=\"n\">layout</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;mid&quot;</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">split</span><span class=\"p\">(</span>\n
    \               <span class=\"n\">Layout</span><span class=\"p\">(</span><span
    class=\"n\">name</span><span class=\"o\">=</span><span class=\"s2\">&quot;server&quot;</span><span
    class=\"p\">),</span>\n                <span class=\"n\">Layout</span><span class=\"p\">(</span><span
    class=\"n\">name</span><span class=\"o\">=</span><span class=\"s2\">&quot;runner&quot;</span><span
    class=\"p\">),</span>\n            <span class=\"p\">)</span>\n            <span
    class=\"n\">layout</span><span class=\"p\">[</span><span class=\"s2\">&quot;side&quot;</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">split</span><span
    class=\"p\">(</span>\n                <span class=\"n\">Layout</span><span class=\"p\">(</span><span
    class=\"n\">name</span><span class=\"o\">=</span><span class=\"s2\">&quot;plugins&quot;</span><span
    class=\"p\">),</span>\n            <span class=\"p\">)</span>\n            <span
    class=\"k\">return</span> <span class=\"n\">layout</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='run_until_keyboard_interrupt' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>run_until_keyboard_interrupt <em class='small'>function</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">run_until_keyboard_interrupt
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
    <span class=\"nf\">run_until_keyboard_interrupt</span><span class=\"p\">()</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">try</span><span class=\"p\">:</span>\n                <span
    class=\"k\">while</span> <span class=\"kc\">True</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">time</span><span class=\"o\">.</span><span
    class=\"n\">sleep</span><span class=\"p\">(</span><span class=\"mf\">0.2</span><span
    class=\"p\">)</span>\n            <span class=\"k\">except</span> <span class=\"ne\">KeyboardInterrupt</span><span
    class=\"p\">:</span>\n                <span class=\"k\">pass</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='version_callback' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>version_callback <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">version_callback
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
    <span class=\"nf\">version_callback</span><span class=\"p\">(</span><span class=\"n\">value</span><span
    class=\"p\">:</span> <span class=\"nb\">bool</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">if</span> <span class=\"n\">value</span><span class=\"p\">:</span>\n
    \               <span class=\"kn\">from</span> <span class=\"nn\">markata</span>
    <span class=\"kn\">import</span> <span class=\"n\">__version__</span>\n\n                <span
    class=\"n\">typer</span><span class=\"o\">.</span><span class=\"n\">echo</span><span
    class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">&quot;Markata
    CLI Version: </span><span class=\"si\">{</span><span class=\"n\">__version__</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">raise</span> <span class=\"n\">typer</span><span
    class=\"o\">.</span><span class=\"n\">Exit</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='json_callback' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>json_callback <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">json_callback
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
    <span class=\"nf\">json_callback</span><span class=\"p\">(</span><span class=\"n\">value</span><span
    class=\"p\">:</span> <span class=\"nb\">bool</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">if</span> <span class=\"n\">value</span><span class=\"p\">:</span>\n
    \               <span class=\"kn\">from</span> <span class=\"nn\">markata</span>
    <span class=\"kn\">import</span> <span class=\"n\">Markata</span>\n\n                <span
    class=\"n\">typer</span><span class=\"o\">.</span><span class=\"n\">echo</span><span
    class=\"p\">(</span><span class=\"n\">Markata</span><span class=\"p\">()</span><span
    class=\"o\">.</span><span class=\"n\">to_json</span><span class=\"p\">())</span>\n
    \               <span class=\"k\">raise</span> <span class=\"n\">typer</span><span
    class=\"o\">.</span><span class=\"n\">Exit</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='main' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>main <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">main
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
    <span class=\"nf\">main</span><span class=\"p\">(</span>\n            <span class=\"n\">version</span><span
    class=\"p\">:</span> <span class=\"nb\">bool</span> <span class=\"o\">=</span>
    <span class=\"n\">typer</span><span class=\"o\">.</span><span class=\"n\">Option</span><span
    class=\"p\">(</span>\n                <span class=\"kc\">None</span><span class=\"p\">,</span>\n
    \               <span class=\"s2\">&quot;--version&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">callback</span><span class=\"o\">=</span><span
    class=\"n\">version_callback</span><span class=\"p\">,</span>\n                <span
    class=\"n\">is_eager</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span>\n            <span class=\"p\">),</span>\n            <span
    class=\"n\">to_json</span><span class=\"p\">:</span> <span class=\"nb\">bool</span>
    <span class=\"o\">=</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
    class=\"n\">Option</span><span class=\"p\">(</span>\n                <span class=\"kc\">None</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;--to-json&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">callback</span><span class=\"o\">=</span><span
    class=\"n\">json_callback</span><span class=\"p\">,</span>\n                <span
    class=\"n\">is_eager</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span>\n            <span class=\"p\">),</span>\n        <span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"c1\"># Do other global stuff, handle other global options
    here</span>\n            <span class=\"k\">return</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='cli' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>cli
    <em class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">cli <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">cli</span><span class=\"p\">()</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"kn\">from</span>
    <span class=\"nn\">markata</span> <span class=\"kn\">import</span> <span class=\"n\">Markata</span>\n\n
    \           <span class=\"n\">m</span> <span class=\"o\">=</span> <span class=\"n\">Markata</span><span
    class=\"p\">()</span>\n            <span class=\"n\">m</span><span class=\"o\">.</span><span
    class=\"n\">_pm</span><span class=\"o\">.</span><span class=\"n\">hook</span><span
    class=\"o\">.</span><span class=\"n\">cli</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"o\">=</span><span class=\"n\">m</span><span
    class=\"p\">,</span> <span class=\"n\">app</span><span class=\"o\">=</span><span
    class=\"n\">app</span><span class=\"p\">)</span>\n            <span class=\"n\">app</span><span
    class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>"
  raw.md: ''
published: true
slug: markata/cli/cli
title: Cli.Py


---

None


!! function <h2 id='make_layout' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>make_layout <em class='small'>function</em></h2>
    Define the layout.
???+ source "make_layout <em class='small'>source</em>"

```python

        def make_layout() -> Layout:
            """Define the layout."""
            layout = Layout(name="root")

            layout.split(
                Layout(name="header", size=3),
                Layout(name="main"),
            )
            layout["main"].split_row(
                Layout(name="side", ratio=50),
                Layout(name="mid", ratio=30),
                Layout(name="describe", ratio=20),
            )
            layout["mid"].split(
                Layout(name="server"),
                Layout(name="runner"),
            )
            layout["side"].split(
                Layout(name="plugins"),
            )
            return layout
```


!! function <h2 id='run_until_keyboard_interrupt' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>run_until_keyboard_interrupt <em class='small'>function</em></h2>

???+ source "run_until_keyboard_interrupt <em class='small'>source</em>"

```python

        def run_until_keyboard_interrupt() -> None:
            try:
                while True:
                    time.sleep(0.2)
            except KeyboardInterrupt:
                pass
```


!! function <h2 id='version_callback' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>version_callback <em class='small'>function</em></h2>

???+ source "version_callback <em class='small'>source</em>"

```python

        def version_callback(value: bool) -> None:
            if value:
                from markata import __version__

                typer.echo(f"Markata CLI Version: {__version__}")
                raise typer.Exit
```


!! function <h2 id='json_callback' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>json_callback <em class='small'>function</em></h2>

???+ source "json_callback <em class='small'>source</em>"

```python

        def json_callback(value: bool) -> None:
            if value:
                from markata import Markata

                typer.echo(Markata().to_json())
                raise typer.Exit
```


!! function <h2 id='main' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>main <em class='small'>function</em></h2>

???+ source "main <em class='small'>source</em>"

```python

        def main(
            version: bool = typer.Option(
                None,
                "--version",
                callback=version_callback,
                is_eager=True,
            ),
            to_json: bool = typer.Option(
                None,
                "--to-json",
                callback=json_callback,
                is_eager=True,
            ),
        ) -> None:
            # Do other global stuff, handle other global options here
            return
```


!! function <h2 id='cli' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>cli <em class='small'>function</em></h2>

???+ source "cli <em class='small'>source</em>"

```python

        def cli() -> None:
            from markata import Markata

            m = Markata()
            m._pm.hook.cli(markata=m, app=app)
            app()
```

