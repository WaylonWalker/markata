---
content: "Markata plugin to create a pyinstrument profile if pyinstrument is installed.\n\nThe
  profile will be saved to <output_dir>/_profile/index.html\n\n\n!! class <h2 id='ProfilerConfig'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>ProfilerConfig <em
  class='small'>class</em></h2>\n\n???+ source \"ProfilerConfig <em class='small'>source</em>\"\n\n```python\n\n
  \       class ProfilerConfig(pydantic.BaseModel):\n            output_dir: pydantic.DirectoryPath
  = Path(\"markout\")\n            should_profile: bool = SHOULD_PROFILE\n            profiler:
  Optional[Any] = (\n                None  # No valicator for type pyinstrument.profiler.Profiler\n
  \           )\n            output_file: Optional[Path] = None\n\n            @pydantic.validator(\"output_file\",
  pre=True, always=True)\n            def validate_output_file(cls, v, *, values):\n
  \               if v is None:\n                    output_file = values[\"output_dir\"]
  / \"_profile\" / \"index.html\"\n                    output_file.parent.mkdir(parents=True,
  exist_ok=True)\n                    return output_file\n                return v\n```\n\n\n!!
  class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Config
  <em class='small'>class</em></h2>\n\n???+ source \"Config <em class='small'>source</em>\"\n\n```python\n\n
  \       class Config(pydantic.BaseModel):\n            should_profile: bool = SHOULD_PROFILE\n
  \           profiler: ProfilerConfig = ProfilerConfig()\n```\n\n\n!! function <h2
  id='config_model' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>config_model
  <em class='small'>function</em></h2>\n\n???+ source \"config_model <em class='small'>source</em>\"\n\n```python\n\n
  \       def config_model(markata: \"Markata\") -> None:\n            markata.config_models.append(Config)\n```\n\n\n!!
  function <h2 id='configure' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>configure <em class='small'>function</em></h2>\n\n???+ source \"configure
  <em class='small'>source</em>\"\n\n```python\n\n        def configure(markata: \"Markata\")
  -> None:\n            # profiler must exist in the same thread and cannot be configured
  through pydantic validation\n            if (\n                markata.config.profiler.should_profile\n
  \               and markata.config.profiler.profiler is None\n            ):\n                markata.config.profiler.profiler
  = Profiler()\n                markata.config.profiler.profiler.start()\n```\n\n\n!!
  function <h2 id='save' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>save
  <em class='small'>function</em></h2>\n    stop the profiler and save as late as
  possible\n???+ source \"save <em class='small'>source</em>\"\n\n```python\n\n        def
  save(markata: Markata) -> None:\n            \"stop the profiler and save as late
  as possible\"\n            if markata.config.profiler.should_profile:\n                if
  markata.config.profiler.profiler is not None:\n                    if markata.config.profiler.profiler.is_running:\n
  \                       try:\n                            markata.config.profiler.profiler.stop()\n
  \                           html = markata.config.profiler.profiler.output_html()\n
  \                           markata.config.profiler.output_file.write_text(html)\n
  \                           markata.console.print(\n                                markata.config.profiler.profiler.output_text()\n
  \                           )\n\n                        except AttributeError:\n
  \                           markata.console.log(\n                                \"profiler
  not available, skipping save pyinstrument save\",\n                            )\n
  \                           markata.console.log(\n                                \"[red]to
  enable profiler [wheat1][itallic]pip install 'markata\\[pyinstrument]'\",\n                            )\n```\n\n\n!!
  function <h2 id='teardown' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>teardown <em class='small'>function</em></h2>\n    stop the profiler on exit\n???+
  source \"teardown <em class='small'>source</em>\"\n\n```python\n\n        def teardown(markata:
  Markata) -> None:\n            \"stop the profiler on exit\"\n            # import
  logging\n\n            # logger = logging.getLogger()\n            # logger.handlers.clear()\n
  \           if markata.config.profiler.should_profile:\n                if markata.config.profiler.profiler
  is not None:\n                    if markata.config.profiler.profiler.is_running:\n
  \                       markata.config.profiler.profiler.stop()\n```\n\n\n!! method
  <h2 id='validate_output_file' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>validate_output_file <em class='small'>method</em></h2>\n\n???+ source \"validate_output_file
  <em class='small'>source</em>\"\n\n```python\n\n        def validate_output_file(cls,
  v, *, values):\n                if v is None:\n                    output_file =
  values[\"output_dir\"] / \"_profile\" / \"index.html\"\n                    output_file.parent.mkdir(parents=True,
  exist_ok=True)\n                    return output_file\n                return v\n```\n\n"
date: 0001-01-01
description: 'Markata plugin to create a pyinstrument profile if pyinstrument is installed.
  The profile will be saved to  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ '
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Pyinstrument.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Markata plugin to create a pyinstrument
    profile if pyinstrument is installed. The profile will be saved to  ! ???+ source
    \ ! ???+ source  ! ???+ source  ! ???+ \" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Pyinstrument.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Markata plugin to create a pyinstrument
    profile if pyinstrument is installed. The profile will be saved to  ! ???+ source
    \ ! ???+ source  ! ???+ source  ! ???+ \" />\n <link href=\"/favicon.ico\" rel=\"icon\"
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Pyinstrument.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <p>Markata plugin to
    create a pyinstrument profile if pyinstrument is installed.</p>\n<p>The profile
    will be saved to &lt;output_dir&gt;/_profile/index.html</p>\n<p>!! class <h2 id='ProfilerConfig'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>ProfilerConfig <em
    class='small'>class</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">ProfilerConfig <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nc\">ProfilerConfig</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">output_dir</span><span class=\"p\">:</span> <span
    class=\"n\">pydantic</span><span class=\"o\">.</span><span class=\"n\">DirectoryPath</span>
    <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;markout&quot;</span><span class=\"p\">)</span>\n            <span
    class=\"n\">should_profile</span><span class=\"p\">:</span> <span class=\"nb\">bool</span>
    <span class=\"o\">=</span> <span class=\"n\">SHOULD_PROFILE</span>\n            <span
    class=\"n\">profiler</span><span class=\"p\">:</span> <span class=\"n\">Optional</span><span
    class=\"p\">[</span><span class=\"n\">Any</span><span class=\"p\">]</span> <span
    class=\"o\">=</span> <span class=\"p\">(</span>\n                <span class=\"kc\">None</span>
    \ <span class=\"c1\"># No valicator for type pyinstrument.profiler.Profiler</span>\n
    \           <span class=\"p\">)</span>\n            <span class=\"n\">output_file</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"n\">Path</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"kc\">None</span>\n\n            <span class=\"nd\">@pydantic</span><span
    class=\"o\">.</span><span class=\"n\">validator</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;output_file&quot;</span><span class=\"p\">,</span> <span class=\"n\">pre</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
    class=\"n\">always</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n            <span class=\"k\">def</span> <span class=\"nf\">validate_output_file</span><span
    class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"p\">,</span> <span class=\"n\">values</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">output_file</span> <span class=\"o\">=</span> <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;output_dir&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">/</span> <span class=\"s2\">&quot;_profile&quot;</span> <span
    class=\"o\">/</span> <span class=\"s2\">&quot;index.html&quot;</span>\n                    <span
    class=\"n\">output_file</span><span class=\"o\">.</span><span class=\"n\">parent</span><span
    class=\"o\">.</span><span class=\"n\">mkdir</span><span class=\"p\">(</span><span
    class=\"n\">parents</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span> <span class=\"n\">exist_ok</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">output_file</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    \           <span class=\"n\">should_profile</span><span class=\"p\">:</span>
    <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"n\">SHOULD_PROFILE</span>\n
    \           <span class=\"n\">profiler</span><span class=\"p\">:</span> <span
    class=\"n\">ProfilerConfig</span> <span class=\"o\">=</span> <span class=\"n\">ProfilerConfig</span><span
    class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='config_model'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>config_model <em
    class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">config_model <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">config_model</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config_models</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">Config</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    \           <span class=\"c1\"># profiler must exist in the same thread and cannot
    be configured through pydantic validation</span>\n            <span class=\"k\">if</span>
    <span class=\"p\">(</span>\n                <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">profiler</span><span class=\"o\">.</span><span class=\"n\">should_profile</span>\n
    \               <span class=\"ow\">and</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">profiler</span><span class=\"o\">.</span><span class=\"n\">profiler</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span>\n            <span
    class=\"p\">):</span>\n                <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">profiler</span><span
    class=\"o\">.</span><span class=\"n\">profiler</span> <span class=\"o\">=</span>
    <span class=\"n\">Profiler</span><span class=\"p\">()</span>\n                <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">profiler</span><span class=\"o\">.</span><span
    class=\"n\">profiler</span><span class=\"o\">.</span><span class=\"n\">start</span><span
    class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='save'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>save <em class='small'>function</em></h2>\nstop
    the profiler and save as late as possible</p>\n<div class=\"admonition source
    is-collapsible collapsible-open\">\n<p class=\"admonition-title\">save <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">save</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"n\">Markata</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"s2\">&quot;stop the profiler and save as late as possible&quot;</span>\n
    \           <span class=\"k\">if</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">profiler</span><span class=\"o\">.</span><span class=\"n\">should_profile</span><span
    class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">profiler</span><span class=\"o\">.</span><span class=\"n\">profiler</span>
    <span class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">if</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">profiler</span><span class=\"o\">.</span><span class=\"n\">profiler</span><span
    class=\"o\">.</span><span class=\"n\">is_running</span><span class=\"p\">:</span>\n
    \                       <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                           <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">profiler</span><span
    class=\"o\">.</span><span class=\"n\">profiler</span><span class=\"o\">.</span><span
    class=\"n\">stop</span><span class=\"p\">()</span>\n                            <span
    class=\"n\">html</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">profiler</span><span class=\"o\">.</span><span class=\"n\">profiler</span><span
    class=\"o\">.</span><span class=\"n\">output_html</span><span class=\"p\">()</span>\n
    \                           <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">profiler</span><span
    class=\"o\">.</span><span class=\"n\">output_file</span><span class=\"o\">.</span><span
    class=\"n\">write_text</span><span class=\"p\">(</span><span class=\"n\">html</span><span
    class=\"p\">)</span>\n                            <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">print</span><span class=\"p\">(</span>\n                                <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">profiler</span><span class=\"o\">.</span><span
    class=\"n\">profiler</span><span class=\"o\">.</span><span class=\"n\">output_text</span><span
    class=\"p\">()</span>\n                            <span class=\"p\">)</span>\n\n
    \                       <span class=\"k\">except</span> <span class=\"ne\">AttributeError</span><span
    class=\"p\">:</span>\n                            <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">log</span><span class=\"p\">(</span>\n                                <span
    class=\"s2\">&quot;profiler not available, skipping save pyinstrument save&quot;</span><span
    class=\"p\">,</span>\n                            <span class=\"p\">)</span>\n
    \                           <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
    class=\"p\">(</span>\n                                <span class=\"s2\">&quot;[red]to
    enable profiler [wheat1][itallic]pip install &#39;markata\\[pyinstrument]&#39;&quot;</span><span
    class=\"p\">,</span>\n                            <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='teardown' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>teardown <em class='small'>function</em></h2>\nstop the profiler on exit</p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">teardown
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
    <span class=\"nf\">teardown</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"n\">Markata</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"s2\">&quot;stop the profiler on exit&quot;</span>\n
    \           <span class=\"c1\"># import logging</span>\n\n            <span class=\"c1\">#
    logger = logging.getLogger()</span>\n            <span class=\"c1\"># logger.handlers.clear()</span>\n
    \           <span class=\"k\">if</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">profiler</span><span class=\"o\">.</span><span class=\"n\">should_profile</span><span
    class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">profiler</span><span class=\"o\">.</span><span class=\"n\">profiler</span>
    <span class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">if</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">profiler</span><span class=\"o\">.</span><span class=\"n\">profiler</span><span
    class=\"o\">.</span><span class=\"n\">is_running</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">profiler</span><span
    class=\"o\">.</span><span class=\"n\">profiler</span><span class=\"o\">.</span><span
    class=\"n\">stop</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='validate_output_file' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>validate_output_file <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">validate_output_file
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
    <span class=\"nf\">validate_output_file</span><span class=\"p\">(</span><span
    class=\"bp\">cls</span><span class=\"p\">,</span> <span class=\"n\">v</span><span
    class=\"p\">,</span> <span class=\"o\">*</span><span class=\"p\">,</span> <span
    class=\"n\">values</span><span class=\"p\">):</span>\n                <span class=\"k\">if</span>
    <span class=\"n\">v</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">output_file</span>
    <span class=\"o\">=</span> <span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;output_dir&quot;</span><span class=\"p\">]</span> <span class=\"o\">/</span>
    <span class=\"s2\">&quot;_profile&quot;</span> <span class=\"o\">/</span> <span
    class=\"s2\">&quot;index.html&quot;</span>\n                    <span class=\"n\">output_file</span><span
    class=\"o\">.</span><span class=\"n\">parent</span><span class=\"o\">.</span><span
    class=\"n\">mkdir</span><span class=\"p\">(</span><span class=\"n\">parents</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
    class=\"n\">exist_ok</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">output_file</span>\n                <span class=\"k\">return</span>
    <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>
    \   </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Pyinstrument.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Markata plugin to create a pyinstrument
    profile if pyinstrument is installed. The profile will be saved to  ! ???+ source
    \ ! ???+ source  ! ???+ source  ! ???+ \" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Pyinstrument.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Markata plugin to create a pyinstrument
    profile if pyinstrument is installed. The profile will be saved to  ! ???+ source
    \ ! ???+ source  ! ???+ source  ! ???+ \" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<article style=\"text-align: center;\">\n    <style>\n        section
    {\n            font-size: 200%;\n        }\n\n\n        .edit {\n            display:
    none;\n        }\n    </style>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Pyinstrument.Py\n    </h1>\n</section></article>\n     </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Pyinstrument.Py\n    </h1>\n</section>    <section class=\"body\">\n        <p>Markata
    plugin to create a pyinstrument profile if pyinstrument is installed.</p>\n<p>The
    profile will be saved to &lt;output_dir&gt;/_profile/index.html</p>\n<p>!! class
    <h2 id='ProfilerConfig' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>ProfilerConfig <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">ProfilerConfig
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
    <span class=\"nc\">ProfilerConfig</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">output_dir</span><span class=\"p\">:</span> <span
    class=\"n\">pydantic</span><span class=\"o\">.</span><span class=\"n\">DirectoryPath</span>
    <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;markout&quot;</span><span class=\"p\">)</span>\n            <span
    class=\"n\">should_profile</span><span class=\"p\">:</span> <span class=\"nb\">bool</span>
    <span class=\"o\">=</span> <span class=\"n\">SHOULD_PROFILE</span>\n            <span
    class=\"n\">profiler</span><span class=\"p\">:</span> <span class=\"n\">Optional</span><span
    class=\"p\">[</span><span class=\"n\">Any</span><span class=\"p\">]</span> <span
    class=\"o\">=</span> <span class=\"p\">(</span>\n                <span class=\"kc\">None</span>
    \ <span class=\"c1\"># No valicator for type pyinstrument.profiler.Profiler</span>\n
    \           <span class=\"p\">)</span>\n            <span class=\"n\">output_file</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"n\">Path</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"kc\">None</span>\n\n            <span class=\"nd\">@pydantic</span><span
    class=\"o\">.</span><span class=\"n\">validator</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;output_file&quot;</span><span class=\"p\">,</span> <span class=\"n\">pre</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
    class=\"n\">always</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n            <span class=\"k\">def</span> <span class=\"nf\">validate_output_file</span><span
    class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"p\">,</span> <span class=\"n\">values</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">output_file</span> <span class=\"o\">=</span> <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;output_dir&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">/</span> <span class=\"s2\">&quot;_profile&quot;</span> <span
    class=\"o\">/</span> <span class=\"s2\">&quot;index.html&quot;</span>\n                    <span
    class=\"n\">output_file</span><span class=\"o\">.</span><span class=\"n\">parent</span><span
    class=\"o\">.</span><span class=\"n\">mkdir</span><span class=\"p\">(</span><span
    class=\"n\">parents</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span> <span class=\"n\">exist_ok</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">output_file</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    \           <span class=\"n\">should_profile</span><span class=\"p\">:</span>
    <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"n\">SHOULD_PROFILE</span>\n
    \           <span class=\"n\">profiler</span><span class=\"p\">:</span> <span
    class=\"n\">ProfilerConfig</span> <span class=\"o\">=</span> <span class=\"n\">ProfilerConfig</span><span
    class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='config_model'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>config_model <em
    class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">config_model <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">config_model</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config_models</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">Config</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    \           <span class=\"c1\"># profiler must exist in the same thread and cannot
    be configured through pydantic validation</span>\n            <span class=\"k\">if</span>
    <span class=\"p\">(</span>\n                <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">profiler</span><span class=\"o\">.</span><span class=\"n\">should_profile</span>\n
    \               <span class=\"ow\">and</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">profiler</span><span class=\"o\">.</span><span class=\"n\">profiler</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span>\n            <span
    class=\"p\">):</span>\n                <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">profiler</span><span
    class=\"o\">.</span><span class=\"n\">profiler</span> <span class=\"o\">=</span>
    <span class=\"n\">Profiler</span><span class=\"p\">()</span>\n                <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">profiler</span><span class=\"o\">.</span><span
    class=\"n\">profiler</span><span class=\"o\">.</span><span class=\"n\">start</span><span
    class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='save'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>save <em class='small'>function</em></h2>\nstop
    the profiler and save as late as possible</p>\n<div class=\"admonition source
    is-collapsible collapsible-open\">\n<p class=\"admonition-title\">save <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">save</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"n\">Markata</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"s2\">&quot;stop the profiler and save as late as possible&quot;</span>\n
    \           <span class=\"k\">if</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">profiler</span><span class=\"o\">.</span><span class=\"n\">should_profile</span><span
    class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">profiler</span><span class=\"o\">.</span><span class=\"n\">profiler</span>
    <span class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">if</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">profiler</span><span class=\"o\">.</span><span class=\"n\">profiler</span><span
    class=\"o\">.</span><span class=\"n\">is_running</span><span class=\"p\">:</span>\n
    \                       <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                           <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">profiler</span><span
    class=\"o\">.</span><span class=\"n\">profiler</span><span class=\"o\">.</span><span
    class=\"n\">stop</span><span class=\"p\">()</span>\n                            <span
    class=\"n\">html</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">profiler</span><span class=\"o\">.</span><span class=\"n\">profiler</span><span
    class=\"o\">.</span><span class=\"n\">output_html</span><span class=\"p\">()</span>\n
    \                           <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">profiler</span><span
    class=\"o\">.</span><span class=\"n\">output_file</span><span class=\"o\">.</span><span
    class=\"n\">write_text</span><span class=\"p\">(</span><span class=\"n\">html</span><span
    class=\"p\">)</span>\n                            <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">print</span><span class=\"p\">(</span>\n                                <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">profiler</span><span class=\"o\">.</span><span
    class=\"n\">profiler</span><span class=\"o\">.</span><span class=\"n\">output_text</span><span
    class=\"p\">()</span>\n                            <span class=\"p\">)</span>\n\n
    \                       <span class=\"k\">except</span> <span class=\"ne\">AttributeError</span><span
    class=\"p\">:</span>\n                            <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">log</span><span class=\"p\">(</span>\n                                <span
    class=\"s2\">&quot;profiler not available, skipping save pyinstrument save&quot;</span><span
    class=\"p\">,</span>\n                            <span class=\"p\">)</span>\n
    \                           <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
    class=\"p\">(</span>\n                                <span class=\"s2\">&quot;[red]to
    enable profiler [wheat1][itallic]pip install &#39;markata\\[pyinstrument]&#39;&quot;</span><span
    class=\"p\">,</span>\n                            <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='teardown' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>teardown <em class='small'>function</em></h2>\nstop the profiler on exit</p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">teardown
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
    <span class=\"nf\">teardown</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"n\">Markata</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"s2\">&quot;stop the profiler on exit&quot;</span>\n
    \           <span class=\"c1\"># import logging</span>\n\n            <span class=\"c1\">#
    logger = logging.getLogger()</span>\n            <span class=\"c1\"># logger.handlers.clear()</span>\n
    \           <span class=\"k\">if</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">profiler</span><span class=\"o\">.</span><span class=\"n\">should_profile</span><span
    class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">profiler</span><span class=\"o\">.</span><span class=\"n\">profiler</span>
    <span class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">if</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">profiler</span><span class=\"o\">.</span><span class=\"n\">profiler</span><span
    class=\"o\">.</span><span class=\"n\">is_running</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">profiler</span><span
    class=\"o\">.</span><span class=\"n\">profiler</span><span class=\"o\">.</span><span
    class=\"n\">stop</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='validate_output_file' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>validate_output_file <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">validate_output_file
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
    <span class=\"nf\">validate_output_file</span><span class=\"p\">(</span><span
    class=\"bp\">cls</span><span class=\"p\">,</span> <span class=\"n\">v</span><span
    class=\"p\">,</span> <span class=\"o\">*</span><span class=\"p\">,</span> <span
    class=\"n\">values</span><span class=\"p\">):</span>\n                <span class=\"k\">if</span>
    <span class=\"n\">v</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">output_file</span>
    <span class=\"o\">=</span> <span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;output_dir&quot;</span><span class=\"p\">]</span> <span class=\"o\">/</span>
    <span class=\"s2\">&quot;_profile&quot;</span> <span class=\"o\">/</span> <span
    class=\"s2\">&quot;index.html&quot;</span>\n                    <span class=\"n\">output_file</span><span
    class=\"o\">.</span><span class=\"n\">parent</span><span class=\"o\">.</span><span
    class=\"n\">mkdir</span><span class=\"p\">(</span><span class=\"n\">parents</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
    class=\"n\">exist_ok</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">output_file</span>\n                <span class=\"k\">return</span>
    <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>"
  raw.md: ''
published: true
slug: markata/plugins/pyinstrument
title: Pyinstrument.Py


---

Markata plugin to create a pyinstrument profile if pyinstrument is installed.

The profile will be saved to <output_dir>/_profile/index.html


!! class <h2 id='ProfilerConfig' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>ProfilerConfig <em class='small'>class</em></h2>

???+ source "ProfilerConfig <em class='small'>source</em>"

```python

        class ProfilerConfig(pydantic.BaseModel):
            output_dir: pydantic.DirectoryPath = Path("markout")
            should_profile: bool = SHOULD_PROFILE
            profiler: Optional[Any] = (
                None  # No valicator for type pyinstrument.profiler.Profiler
            )
            output_file: Optional[Path] = None

            @pydantic.validator("output_file", pre=True, always=True)
            def validate_output_file(cls, v, *, values):
                if v is None:
                    output_file = values["output_dir"] / "_profile" / "index.html"
                    output_file.parent.mkdir(parents=True, exist_ok=True)
                    return output_file
                return v
```


!! class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Config <em class='small'>class</em></h2>

???+ source "Config <em class='small'>source</em>"

```python

        class Config(pydantic.BaseModel):
            should_profile: bool = SHOULD_PROFILE
            profiler: ProfilerConfig = ProfilerConfig()
```


!! function <h2 id='config_model' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>config_model <em class='small'>function</em></h2>

???+ source "config_model <em class='small'>source</em>"

```python

        def config_model(markata: "Markata") -> None:
            markata.config_models.append(Config)
```


!! function <h2 id='configure' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>configure <em class='small'>function</em></h2>

???+ source "configure <em class='small'>source</em>"

```python

        def configure(markata: "Markata") -> None:
            # profiler must exist in the same thread and cannot be configured through pydantic validation
            if (
                markata.config.profiler.should_profile
                and markata.config.profiler.profiler is None
            ):
                markata.config.profiler.profiler = Profiler()
                markata.config.profiler.profiler.start()
```


!! function <h2 id='save' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>save <em class='small'>function</em></h2>
    stop the profiler and save as late as possible
???+ source "save <em class='small'>source</em>"

```python

        def save(markata: Markata) -> None:
            "stop the profiler and save as late as possible"
            if markata.config.profiler.should_profile:
                if markata.config.profiler.profiler is not None:
                    if markata.config.profiler.profiler.is_running:
                        try:
                            markata.config.profiler.profiler.stop()
                            html = markata.config.profiler.profiler.output_html()
                            markata.config.profiler.output_file.write_text(html)
                            markata.console.print(
                                markata.config.profiler.profiler.output_text()
                            )

                        except AttributeError:
                            markata.console.log(
                                "profiler not available, skipping save pyinstrument save",
                            )
                            markata.console.log(
                                "[red]to enable profiler [wheat1][itallic]pip install 'markata\[pyinstrument]'",
                            )
```


!! function <h2 id='teardown' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>teardown <em class='small'>function</em></h2>
    stop the profiler on exit
???+ source "teardown <em class='small'>source</em>"

```python

        def teardown(markata: Markata) -> None:
            "stop the profiler on exit"
            # import logging

            # logger = logging.getLogger()
            # logger.handlers.clear()
            if markata.config.profiler.should_profile:
                if markata.config.profiler.profiler is not None:
                    if markata.config.profiler.profiler.is_running:
                        markata.config.profiler.profiler.stop()
```


!! method <h2 id='validate_output_file' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>validate_output_file <em class='small'>method</em></h2>

???+ source "validate_output_file <em class='small'>source</em>"

```python

        def validate_output_file(cls, v, *, values):
                if v is None:
                    output_file = values["output_dir"] / "_profile" / "index.html"
                    output_file.parent.mkdir(parents=True, exist_ok=True)
                    return output_file
                return v
```

