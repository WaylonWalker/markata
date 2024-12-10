---
content: "None\n\n\n!! class <h2 id='TuiKey' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>TuiKey <em class='small'>class</em></h2>\n\n???+ source \"TuiKey <em class='small'>source</em>\"\n\n```python\n\n
  \       class TuiKey(pydantic.BaseModel):\n            name: str\n            key:
  str\n```\n\n\n!! class <h2 id='TuiConfig' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>TuiConfig <em class='small'>class</em></h2>\n\n???+ source \"TuiConfig <em
  class='small'>source</em>\"\n\n```python\n\n        class TuiConfig(pydantic.BaseModel):\n
  \           new_cmd: List[str] = [\"markata\", \"new\", \"post\"]\n            keymap:
  List[TuiKey] = [TuiKey(name=\"new\", key=\"n\")]\n```\n\n\n!! class <h2 id='Config'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Config <em class='small'>class</em></h2>\n\n???+
  source \"Config <em class='small'>source</em>\"\n\n```python\n\n        class Config(pydantic.BaseModel):\n
  \           tui: TuiConfig = TuiConfig()\n```\n\n\n!! function <h2 id='config_model'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>config_model <em class='small'>function</em></h2>\n\n???+
  source \"config_model <em class='small'>source</em>\"\n\n```python\n\n        def
  config_model(markata: \"Markata\") -> None:\n            markata.config_models.append(Config)\n```\n\n\n!!
  function <h2 id='cli' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>cli
  <em class='small'>function</em></h2>\n\n???+ source \"cli <em class='small'>source</em>\"\n\n```python\n\n
  \       def cli(app, markata):\n            @app.command()\n            def tui():\n
  \               MarkataApp.run(log=\"textual.log\")\n```\n\n\n!! class <h2 id='MarkataWidget'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>MarkataWidget <em
  class='small'>class</em></h2>\n\n???+ source \"MarkataWidget <em class='small'>source</em>\"\n\n```python\n\n
  \       class MarkataWidget(Widget):\n                def __init__(self, markata:
  Markata, widget: str = \"server\") -> None:\n                    super().__init__(widget)\n
  \                   self.m = markata\n                    self.widget = widget\n
  \                   self.renderable = getattr(self.m, self.widget)\n\n                    def
  render(self):\n                        return self.renderable\n\n                    async
  def update(self, renderable: RenderableType) -> None:\n                        self.renderable
  = renderable\n                        self.refresh()\n```\n\n\n!! class <h2 id='MarkataApp'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>MarkataApp <em class='small'>class</em></h2>\n\n???+
  source \"MarkataApp <em class='small'>source</em>\"\n\n```python\n\n        class
  MarkataApp(App):\n                async def on_load(self, event):\n                    self.m
  = Markata()\n                    self.m.console.quiet = True\n                    await
  self.bind(\"q\", \"quit\", \"quit\")\n                    await self.bind(\"r\",
  \"refresh\", \"refresh\")\n                    self.new_cmd = self.m.config.get(\"tui\",
  {}).get(\"new_cmd\", \"\")\n                    if self.new_cmd != \"\":\n                        await
  self.bind(\"n\", \"new\", \"new\")\n\n                async def on_mount(self) ->
  None:\n                    self.server = MarkataWidget(self.m, \"server\")\n                    self.runner
  = MarkataWidget(self.m, \"runner\")\n                    self.plugins = MarkataWidget(self.m,
  \"plugins\")\n                    self.summary = MarkataWidget(self.m, \"summary\")\n
  \                   await self.view.dock(Footer(), edge=\"bottom\")\n                    await
  self.view.dock(self.plugins, edge=\"left\", size=30, name=\"plugins\")\n                    await
  self.view.dock(self.summary, edge=\"right\", size=30, name=\"summary\")\n                    await
  self.view.dock(self.server, self.runner, edge=\"top\")\n                    self.set_interval(1,
  self.action_refresh)\n\n                async def action_refresh(self) -> None:\n
  \                   self.refresh()\n                    self.runner.refresh()\n
  \                   self.server.refresh()\n                    self.plugins.refresh()\n
  \                   self.summary.refresh()\n\n                async def action_new(self)
  -> None:\n                    subprocess.Popen(\n                        self.new_cmd,
  stdout=subprocess.PIPE, stderr=subprocess.PIPE\n                    )\n```\n\n\n!!
  function <h2 id='tui' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>tui
  <em class='small'>function</em></h2>\n\n???+ source \"tui <em class='small'>source</em>\"\n\n```python\n\n
  \       def tui():\n                MarkataApp.run(log=\"textual.log\")\n```\n\n\n!!
  method <h2 id='__init__' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>__init__ <em class='small'>method</em></h2>\n\n???+ source \"__init__ <em
  class='small'>source</em>\"\n\n```python\n\n        def __init__(self, markata:
  Markata, widget: str = \"server\") -> None:\n                    super().__init__(widget)\n
  \                   self.m = markata\n                    self.widget = widget\n
  \                   self.renderable = getattr(self.m, self.widget)\n\n                    def
  render(self):\n                        return self.renderable\n\n                    async
  def update(self, renderable: RenderableType) -> None:\n                        self.renderable
  = renderable\n                        self.refresh()\n```\n\n\n!! function <h2 id='render'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>render <em class='small'>function</em></h2>\n\n???+
  source \"render <em class='small'>source</em>\"\n\n```python\n\n        def render(self):\n
  \                       return self.renderable\n```\n\n"
date: 0001-01-01
description: 'None ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
  source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source '
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Tui.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"None ! ???+ source  ! ???+ source  ! ???+ source
    \ ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  !
    ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Tui.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None ! ???+ source  ! ???+ source  !
    ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
    source  ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Tui.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <p>None</p>\n<p>!!
    class <h2 id='TuiKey' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>TuiKey
    <em class='small'>class</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">TuiKey <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nc\">TuiKey</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">name</span><span class=\"p\">:</span> <span class=\"nb\">str</span>\n
    \           <span class=\"n\">key</span><span class=\"p\">:</span> <span class=\"nb\">str</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='TuiConfig' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>TuiConfig <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">TuiConfig
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
    <span class=\"nc\">TuiConfig</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">new_cmd</span><span class=\"p\">:</span> <span class=\"n\">List</span><span
    class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">]</span> <span
    class=\"o\">=</span> <span class=\"p\">[</span><span class=\"s2\">&quot;markata&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;new&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;post&quot;</span><span class=\"p\">]</span>\n            <span
    class=\"n\">keymap</span><span class=\"p\">:</span> <span class=\"n\">List</span><span
    class=\"p\">[</span><span class=\"n\">TuiKey</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"p\">[</span><span class=\"n\">TuiKey</span><span
    class=\"p\">(</span><span class=\"n\">name</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;new&quot;</span><span class=\"p\">,</span> <span class=\"n\">key</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;n&quot;</span><span class=\"p\">)]</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    <span class=\"nc\">Config</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">tui</span><span class=\"p\">:</span> <span class=\"n\">TuiConfig</span>
    <span class=\"o\">=</span> <span class=\"n\">TuiConfig</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='config_model' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>config_model <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">config_model
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
    <span class=\"nf\">config_model</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config_models</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">Config</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    <span class=\"nf\">cli</span><span class=\"p\">(</span><span class=\"n\">app</span><span
    class=\"p\">,</span> <span class=\"n\">markata</span><span class=\"p\">):</span>\n
    \           <span class=\"nd\">@app</span><span class=\"o\">.</span><span class=\"n\">command</span><span
    class=\"p\">()</span>\n            <span class=\"k\">def</span> <span class=\"nf\">tui</span><span
    class=\"p\">():</span>\n                <span class=\"n\">MarkataApp</span><span
    class=\"o\">.</span><span class=\"n\">run</span><span class=\"p\">(</span><span
    class=\"n\">log</span><span class=\"o\">=</span><span class=\"s2\">&quot;textual.log&quot;</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! class <h2 id='MarkataWidget'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>MarkataWidget <em
    class='small'>class</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">MarkataWidget <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nc\">MarkataWidget</span><span class=\"p\">(</span><span class=\"n\">Widget</span><span
    class=\"p\">):</span>\n                <span class=\"k\">def</span> <span class=\"fm\">__init__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"n\">Markata</span><span
    class=\"p\">,</span> <span class=\"n\">widget</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;server&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                    <span class=\"nb\">super</span><span
    class=\"p\">()</span><span class=\"o\">.</span><span class=\"fm\">__init__</span><span
    class=\"p\">(</span><span class=\"n\">widget</span><span class=\"p\">)</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span> <span class=\"o\">=</span> <span class=\"n\">markata</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">widget</span> <span class=\"o\">=</span> <span class=\"n\">widget</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">renderable</span> <span class=\"o\">=</span> <span class=\"nb\">getattr</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"p\">,</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">widget</span><span class=\"p\">)</span>\n\n
    \                   <span class=\"k\">def</span> <span class=\"nf\">render</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">):</span>\n
    \                       <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">renderable</span>\n\n                    <span
    class=\"k\">async</span> <span class=\"k\">def</span> <span class=\"nf\">update</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"n\">renderable</span><span class=\"p\">:</span> <span class=\"n\">RenderableType</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                        <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">renderable</span> <span class=\"o\">=</span>
    <span class=\"n\">renderable</span>\n                        <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">refresh</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    class=\"p\">):</span>\n                <span class=\"k\">async</span> <span class=\"k\">def</span>
    <span class=\"nf\">on_load</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">event</span><span class=\"p\">):</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span> <span class=\"o\">=</span> <span class=\"n\">Markata</span><span
    class=\"p\">()</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">m</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">quiet</span>
    <span class=\"o\">=</span> <span class=\"kc\">True</span>\n                    <span
    class=\"k\">await</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">bind</span><span class=\"p\">(</span><span class=\"s2\">&quot;q&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;quit&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;quit&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">await</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">bind</span><span class=\"p\">(</span><span class=\"s2\">&quot;r&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;refresh&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;refresh&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">new_cmd</span>
    <span class=\"o\">=</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;tui&quot;</span><span class=\"p\">,</span> <span class=\"p\">{})</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;new_cmd&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;&quot;</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">if</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">new_cmd</span> <span class=\"o\">!=</span>
    <span class=\"s2\">&quot;&quot;</span><span class=\"p\">:</span>\n                        <span
    class=\"k\">await</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">bind</span><span class=\"p\">(</span><span class=\"s2\">&quot;n&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;new&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;new&quot;</span><span class=\"p\">)</span>\n\n                <span
    class=\"k\">async</span> <span class=\"k\">def</span> <span class=\"nf\">on_mount</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">server</span> <span class=\"o\">=</span> <span class=\"n\">MarkataWidget</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"p\">,</span> <span class=\"s2\">&quot;server&quot;</span><span
    class=\"p\">)</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">runner</span> <span class=\"o\">=</span>
    <span class=\"n\">MarkataWidget</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">m</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;runner&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">plugins</span>
    <span class=\"o\">=</span> <span class=\"n\">MarkataWidget</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">m</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;plugins&quot;</span><span class=\"p\">)</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">summary</span> <span class=\"o\">=</span> <span class=\"n\">MarkataWidget</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"p\">,</span> <span class=\"s2\">&quot;summary&quot;</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">await</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">view</span><span
    class=\"o\">.</span><span class=\"n\">dock</span><span class=\"p\">(</span><span
    class=\"n\">Footer</span><span class=\"p\">(),</span> <span class=\"n\">edge</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;bottom&quot;</span><span class=\"p\">)</span>\n
    \                   <span class=\"k\">await</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">view</span><span class=\"o\">.</span><span
    class=\"n\">dock</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">plugins</span><span class=\"p\">,</span>
    <span class=\"n\">edge</span><span class=\"o\">=</span><span class=\"s2\">&quot;left&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">size</span><span class=\"o\">=</span><span
    class=\"mi\">30</span><span class=\"p\">,</span> <span class=\"n\">name</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;plugins&quot;</span><span class=\"p\">)</span>\n
    \                   <span class=\"k\">await</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">view</span><span class=\"o\">.</span><span
    class=\"n\">dock</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">summary</span><span class=\"p\">,</span>
    <span class=\"n\">edge</span><span class=\"o\">=</span><span class=\"s2\">&quot;right&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">size</span><span class=\"o\">=</span><span
    class=\"mi\">30</span><span class=\"p\">,</span> <span class=\"n\">name</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;summary&quot;</span><span class=\"p\">)</span>\n
    \                   <span class=\"k\">await</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">view</span><span class=\"o\">.</span><span
    class=\"n\">dock</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">server</span><span class=\"p\">,</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">runner</span><span
    class=\"p\">,</span> <span class=\"n\">edge</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;top&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">set_interval</span><span
    class=\"p\">(</span><span class=\"mi\">1</span><span class=\"p\">,</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">action_refresh</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">async</span> <span class=\"k\">def</span>
    <span class=\"nf\">action_refresh</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">refresh</span><span class=\"p\">()</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">runner</span><span class=\"o\">.</span><span class=\"n\">refresh</span><span
    class=\"p\">()</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">server</span><span class=\"o\">.</span><span
    class=\"n\">refresh</span><span class=\"p\">()</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">plugins</span><span
    class=\"o\">.</span><span class=\"n\">refresh</span><span class=\"p\">()</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">summary</span><span class=\"o\">.</span><span class=\"n\">refresh</span><span
    class=\"p\">()</span>\n\n                <span class=\"k\">async</span> <span
    class=\"k\">def</span> <span class=\"nf\">action_new</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">subprocess</span><span class=\"o\">.</span><span class=\"n\">Popen</span><span
    class=\"p\">(</span>\n                        <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">new_cmd</span><span class=\"p\">,</span>
    <span class=\"n\">stdout</span><span class=\"o\">=</span><span class=\"n\">subprocess</span><span
    class=\"o\">.</span><span class=\"n\">PIPE</span><span class=\"p\">,</span> <span
    class=\"n\">stderr</span><span class=\"o\">=</span><span class=\"n\">subprocess</span><span
    class=\"o\">.</span><span class=\"n\">PIPE</span>\n                    <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='tui' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>tui
    <em class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">tui <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">tui</span><span class=\"p\">():</span>\n                <span
    class=\"n\">MarkataApp</span><span class=\"o\">.</span><span class=\"n\">run</span><span
    class=\"p\">(</span><span class=\"n\">log</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;textual.log&quot;</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                    <span
    class=\"nb\">super</span><span class=\"p\">()</span><span class=\"o\">.</span><span
    class=\"fm\">__init__</span><span class=\"p\">(</span><span class=\"n\">widget</span><span
    class=\"p\">)</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">m</span> <span class=\"o\">=</span> <span
    class=\"n\">markata</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">widget</span> <span class=\"o\">=</span>
    <span class=\"n\">widget</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">renderable</span> <span class=\"o\">=</span>
    <span class=\"nb\">getattr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">m</span><span class=\"p\">,</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">widget</span><span
    class=\"p\">)</span>\n\n                    <span class=\"k\">def</span> <span
    class=\"nf\">render</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">):</span>\n                        <span class=\"k\">return</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">renderable</span>\n\n
    \                   <span class=\"k\">async</span> <span class=\"k\">def</span>
    <span class=\"nf\">update</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">renderable</span><span class=\"p\">:</span>
    <span class=\"n\">RenderableType</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                        <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">renderable</span>
    <span class=\"o\">=</span> <span class=\"n\">renderable</span>\n                        <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">refresh</span><span
    class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='render'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>render <em class='small'>function</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">render
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
    class=\"p\">):</span>\n                        <span class=\"k\">return</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">renderable</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>    </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Tui.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"None ! ???+ source  ! ???+ source  ! ???+ source
    \ ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  !
    ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Tui.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None ! ???+ source  ! ???+ source  !
    ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
    source  ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<article style=\"text-align: center;\">\n    <style>\n        section
    {\n            font-size: 200%;\n        }\n\n\n        .edit {\n            display:
    none;\n        }\n    </style>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Tui.Py\n    </h1>\n</section></article>\n     </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Tui.Py\n    </h1>\n</section>    <section class=\"body\">\n        <p>None</p>\n<p>!!
    class <h2 id='TuiKey' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>TuiKey
    <em class='small'>class</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">TuiKey <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nc\">TuiKey</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">name</span><span class=\"p\">:</span> <span class=\"nb\">str</span>\n
    \           <span class=\"n\">key</span><span class=\"p\">:</span> <span class=\"nb\">str</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='TuiConfig' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>TuiConfig <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">TuiConfig
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
    <span class=\"nc\">TuiConfig</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">new_cmd</span><span class=\"p\">:</span> <span class=\"n\">List</span><span
    class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">]</span> <span
    class=\"o\">=</span> <span class=\"p\">[</span><span class=\"s2\">&quot;markata&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;new&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;post&quot;</span><span class=\"p\">]</span>\n            <span
    class=\"n\">keymap</span><span class=\"p\">:</span> <span class=\"n\">List</span><span
    class=\"p\">[</span><span class=\"n\">TuiKey</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"p\">[</span><span class=\"n\">TuiKey</span><span
    class=\"p\">(</span><span class=\"n\">name</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;new&quot;</span><span class=\"p\">,</span> <span class=\"n\">key</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;n&quot;</span><span class=\"p\">)]</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    <span class=\"nc\">Config</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">tui</span><span class=\"p\">:</span> <span class=\"n\">TuiConfig</span>
    <span class=\"o\">=</span> <span class=\"n\">TuiConfig</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='config_model' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>config_model <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">config_model
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
    <span class=\"nf\">config_model</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config_models</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">Config</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    <span class=\"nf\">cli</span><span class=\"p\">(</span><span class=\"n\">app</span><span
    class=\"p\">,</span> <span class=\"n\">markata</span><span class=\"p\">):</span>\n
    \           <span class=\"nd\">@app</span><span class=\"o\">.</span><span class=\"n\">command</span><span
    class=\"p\">()</span>\n            <span class=\"k\">def</span> <span class=\"nf\">tui</span><span
    class=\"p\">():</span>\n                <span class=\"n\">MarkataApp</span><span
    class=\"o\">.</span><span class=\"n\">run</span><span class=\"p\">(</span><span
    class=\"n\">log</span><span class=\"o\">=</span><span class=\"s2\">&quot;textual.log&quot;</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! class <h2 id='MarkataWidget'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>MarkataWidget <em
    class='small'>class</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">MarkataWidget <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nc\">MarkataWidget</span><span class=\"p\">(</span><span class=\"n\">Widget</span><span
    class=\"p\">):</span>\n                <span class=\"k\">def</span> <span class=\"fm\">__init__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"n\">Markata</span><span
    class=\"p\">,</span> <span class=\"n\">widget</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;server&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                    <span class=\"nb\">super</span><span
    class=\"p\">()</span><span class=\"o\">.</span><span class=\"fm\">__init__</span><span
    class=\"p\">(</span><span class=\"n\">widget</span><span class=\"p\">)</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span> <span class=\"o\">=</span> <span class=\"n\">markata</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">widget</span> <span class=\"o\">=</span> <span class=\"n\">widget</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">renderable</span> <span class=\"o\">=</span> <span class=\"nb\">getattr</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"p\">,</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">widget</span><span class=\"p\">)</span>\n\n
    \                   <span class=\"k\">def</span> <span class=\"nf\">render</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">):</span>\n
    \                       <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">renderable</span>\n\n                    <span
    class=\"k\">async</span> <span class=\"k\">def</span> <span class=\"nf\">update</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"n\">renderable</span><span class=\"p\">:</span> <span class=\"n\">RenderableType</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                        <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">renderable</span> <span class=\"o\">=</span>
    <span class=\"n\">renderable</span>\n                        <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">refresh</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    class=\"p\">):</span>\n                <span class=\"k\">async</span> <span class=\"k\">def</span>
    <span class=\"nf\">on_load</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">event</span><span class=\"p\">):</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span> <span class=\"o\">=</span> <span class=\"n\">Markata</span><span
    class=\"p\">()</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">m</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">quiet</span>
    <span class=\"o\">=</span> <span class=\"kc\">True</span>\n                    <span
    class=\"k\">await</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">bind</span><span class=\"p\">(</span><span class=\"s2\">&quot;q&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;quit&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;quit&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">await</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">bind</span><span class=\"p\">(</span><span class=\"s2\">&quot;r&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;refresh&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;refresh&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">new_cmd</span>
    <span class=\"o\">=</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;tui&quot;</span><span class=\"p\">,</span> <span class=\"p\">{})</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;new_cmd&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;&quot;</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">if</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">new_cmd</span> <span class=\"o\">!=</span>
    <span class=\"s2\">&quot;&quot;</span><span class=\"p\">:</span>\n                        <span
    class=\"k\">await</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">bind</span><span class=\"p\">(</span><span class=\"s2\">&quot;n&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;new&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;new&quot;</span><span class=\"p\">)</span>\n\n                <span
    class=\"k\">async</span> <span class=\"k\">def</span> <span class=\"nf\">on_mount</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">server</span> <span class=\"o\">=</span> <span class=\"n\">MarkataWidget</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"p\">,</span> <span class=\"s2\">&quot;server&quot;</span><span
    class=\"p\">)</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">runner</span> <span class=\"o\">=</span>
    <span class=\"n\">MarkataWidget</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">m</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;runner&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">plugins</span>
    <span class=\"o\">=</span> <span class=\"n\">MarkataWidget</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">m</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;plugins&quot;</span><span class=\"p\">)</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">summary</span> <span class=\"o\">=</span> <span class=\"n\">MarkataWidget</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"p\">,</span> <span class=\"s2\">&quot;summary&quot;</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">await</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">view</span><span
    class=\"o\">.</span><span class=\"n\">dock</span><span class=\"p\">(</span><span
    class=\"n\">Footer</span><span class=\"p\">(),</span> <span class=\"n\">edge</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;bottom&quot;</span><span class=\"p\">)</span>\n
    \                   <span class=\"k\">await</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">view</span><span class=\"o\">.</span><span
    class=\"n\">dock</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">plugins</span><span class=\"p\">,</span>
    <span class=\"n\">edge</span><span class=\"o\">=</span><span class=\"s2\">&quot;left&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">size</span><span class=\"o\">=</span><span
    class=\"mi\">30</span><span class=\"p\">,</span> <span class=\"n\">name</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;plugins&quot;</span><span class=\"p\">)</span>\n
    \                   <span class=\"k\">await</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">view</span><span class=\"o\">.</span><span
    class=\"n\">dock</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">summary</span><span class=\"p\">,</span>
    <span class=\"n\">edge</span><span class=\"o\">=</span><span class=\"s2\">&quot;right&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">size</span><span class=\"o\">=</span><span
    class=\"mi\">30</span><span class=\"p\">,</span> <span class=\"n\">name</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;summary&quot;</span><span class=\"p\">)</span>\n
    \                   <span class=\"k\">await</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">view</span><span class=\"o\">.</span><span
    class=\"n\">dock</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">server</span><span class=\"p\">,</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">runner</span><span
    class=\"p\">,</span> <span class=\"n\">edge</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;top&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">set_interval</span><span
    class=\"p\">(</span><span class=\"mi\">1</span><span class=\"p\">,</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">action_refresh</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">async</span> <span class=\"k\">def</span>
    <span class=\"nf\">action_refresh</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">refresh</span><span class=\"p\">()</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">runner</span><span class=\"o\">.</span><span class=\"n\">refresh</span><span
    class=\"p\">()</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">server</span><span class=\"o\">.</span><span
    class=\"n\">refresh</span><span class=\"p\">()</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">plugins</span><span
    class=\"o\">.</span><span class=\"n\">refresh</span><span class=\"p\">()</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">summary</span><span class=\"o\">.</span><span class=\"n\">refresh</span><span
    class=\"p\">()</span>\n\n                <span class=\"k\">async</span> <span
    class=\"k\">def</span> <span class=\"nf\">action_new</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">subprocess</span><span class=\"o\">.</span><span class=\"n\">Popen</span><span
    class=\"p\">(</span>\n                        <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">new_cmd</span><span class=\"p\">,</span>
    <span class=\"n\">stdout</span><span class=\"o\">=</span><span class=\"n\">subprocess</span><span
    class=\"o\">.</span><span class=\"n\">PIPE</span><span class=\"p\">,</span> <span
    class=\"n\">stderr</span><span class=\"o\">=</span><span class=\"n\">subprocess</span><span
    class=\"o\">.</span><span class=\"n\">PIPE</span>\n                    <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='tui' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>tui
    <em class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">tui <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">tui</span><span class=\"p\">():</span>\n                <span
    class=\"n\">MarkataApp</span><span class=\"o\">.</span><span class=\"n\">run</span><span
    class=\"p\">(</span><span class=\"n\">log</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;textual.log&quot;</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                    <span
    class=\"nb\">super</span><span class=\"p\">()</span><span class=\"o\">.</span><span
    class=\"fm\">__init__</span><span class=\"p\">(</span><span class=\"n\">widget</span><span
    class=\"p\">)</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">m</span> <span class=\"o\">=</span> <span
    class=\"n\">markata</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">widget</span> <span class=\"o\">=</span>
    <span class=\"n\">widget</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">renderable</span> <span class=\"o\">=</span>
    <span class=\"nb\">getattr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">m</span><span class=\"p\">,</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">widget</span><span
    class=\"p\">)</span>\n\n                    <span class=\"k\">def</span> <span
    class=\"nf\">render</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">):</span>\n                        <span class=\"k\">return</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">renderable</span>\n\n
    \                   <span class=\"k\">async</span> <span class=\"k\">def</span>
    <span class=\"nf\">update</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">renderable</span><span class=\"p\">:</span>
    <span class=\"n\">RenderableType</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                        <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">renderable</span>
    <span class=\"o\">=</span> <span class=\"n\">renderable</span>\n                        <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">refresh</span><span
    class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='render'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>render <em class='small'>function</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">render
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
    class=\"p\">):</span>\n                        <span class=\"k\">return</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">renderable</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>"
  raw.md: ''
published: true
slug: markata/plugins/tui
title: Tui.Py


---

None


!! class <h2 id='TuiKey' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>TuiKey <em class='small'>class</em></h2>

???+ source "TuiKey <em class='small'>source</em>"

```python

        class TuiKey(pydantic.BaseModel):
            name: str
            key: str
```


!! class <h2 id='TuiConfig' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>TuiConfig <em class='small'>class</em></h2>

???+ source "TuiConfig <em class='small'>source</em>"

```python

        class TuiConfig(pydantic.BaseModel):
            new_cmd: List[str] = ["markata", "new", "post"]
            keymap: List[TuiKey] = [TuiKey(name="new", key="n")]
```


!! class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Config <em class='small'>class</em></h2>

???+ source "Config <em class='small'>source</em>"

```python

        class Config(pydantic.BaseModel):
            tui: TuiConfig = TuiConfig()
```


!! function <h2 id='config_model' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>config_model <em class='small'>function</em></h2>

???+ source "config_model <em class='small'>source</em>"

```python

        def config_model(markata: "Markata") -> None:
            markata.config_models.append(Config)
```


!! function <h2 id='cli' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>cli <em class='small'>function</em></h2>

???+ source "cli <em class='small'>source</em>"

```python

        def cli(app, markata):
            @app.command()
            def tui():
                MarkataApp.run(log="textual.log")
```


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
                async def on_load(self, event):
                    self.m = Markata()
                    self.m.console.quiet = True
                    await self.bind("q", "quit", "quit")
                    await self.bind("r", "refresh", "refresh")
                    self.new_cmd = self.m.config.get("tui", {}).get("new_cmd", "")
                    if self.new_cmd != "":
                        await self.bind("n", "new", "new")

                async def on_mount(self) -> None:
                    self.server = MarkataWidget(self.m, "server")
                    self.runner = MarkataWidget(self.m, "runner")
                    self.plugins = MarkataWidget(self.m, "plugins")
                    self.summary = MarkataWidget(self.m, "summary")
                    await self.view.dock(Footer(), edge="bottom")
                    await self.view.dock(self.plugins, edge="left", size=30, name="plugins")
                    await self.view.dock(self.summary, edge="right", size=30, name="summary")
                    await self.view.dock(self.server, self.runner, edge="top")
                    self.set_interval(1, self.action_refresh)

                async def action_refresh(self) -> None:
                    self.refresh()
                    self.runner.refresh()
                    self.server.refresh()
                    self.plugins.refresh()
                    self.summary.refresh()

                async def action_new(self) -> None:
                    subprocess.Popen(
                        self.new_cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE
                    )
```


!! function <h2 id='tui' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>tui <em class='small'>function</em></h2>

???+ source "tui <em class='small'>source</em>"

```python

        def tui():
                MarkataApp.run(log="textual.log")
```


!! method <h2 id='__init__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__init__ <em class='small'>method</em></h2>

???+ source "__init__ <em class='small'>source</em>"

```python

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


!! function <h2 id='render' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>render <em class='small'>function</em></h2>

???+ source "render <em class='small'>source</em>"

```python

        def render(self):
                        return self.renderable
```

