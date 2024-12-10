---
content: "Default glob plugin\n\n\n!! class <h2 id='GlobConfig' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>GlobConfig <em class='small'>class</em></h2>\n\n???+
  source \"GlobConfig <em class='small'>source</em>\"\n\n```python\n\n        class
  GlobConfig(pydantic.BaseModel):\n            glob_patterns: Union[List[str], str]
  = [\"**/*.md\"]\n            use_gitignore: bool = True\n\n            @pydantic.validator(\"glob_patterns\")\n
  \           def convert_to_list(cls, v):\n                if not isinstance(v, list):\n
  \                   return v.split(\",\")\n                return v\n```\n\n\n!!
  class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Config
  <em class='small'>class</em></h2>\n\n???+ source \"Config <em class='small'>source</em>\"\n\n```python\n\n
  \       class Config(pydantic.BaseModel):\n            glob: GlobConfig = GlobConfig()\n```\n\n\n!!
  function <h2 id='config_model' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>config_model <em class='small'>function</em></h2>\n\n???+ source \"config_model
  <em class='small'>source</em>\"\n\n```python\n\n        def config_model(markata:
  \"Markata\") -> None:\n            markata.config_models.append(Config)\n```\n\n\n!!
  function <h2 id='glob' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>glob
  <em class='small'>function</em></h2>\n\n???+ source \"glob <em class='small'>source</em>\"\n\n```python\n\n
  \       def glob(markata: \"Markata\") -> None:\n            markata.files = list(\n
  \               flatten(\n                    [\n                        Path().glob(str(pattern))\n
  \                       for pattern in markata.config.glob.glob_patterns\n                    ],\n
  \               ),\n            )\n            markata.content_directories = list({f.parent
  for f in markata.files})\n\n            try:\n                ignore = True\n            except
  KeyError:\n                ignore = True\n\n            if ignore and (Path(\".gitignore\").exists()
  or Path(\".markataignore\").exists()):\n                import pathspec\n\n                lines
  = []\n\n                if Path(\".gitignore\").exists():\n                    lines.extend(Path(\".gitignore\").read_text().splitlines())\n\n
  \               if Path(\".markataignore\").exists():\n                    lines.extend(Path(\".markataignore\").read_text().splitlines())\n\n
  \               key = markata.make_hash(\"glob\", \"spec\", lines)\n                spec
  = markata.precache.get(key)\n                if spec is None:\n                    spec
  = pathspec.PathSpec.from_lines(\"gitwildmatch\", lines)\n                    with
  markata.cache as cache:\n                        cache.set(key, spec)\n\n                @background.task\n
  \               def check_spec(file: str) -> bool:\n                    key = markata.make_hash(\"glob\",
  \"check_spec\", file)\n                    check = markata.precache.get(key)\n                    if
  check is not None:\n                        return check\n\n                    check
  = spec.match_file(str(file))\n                    with markata.cache as cache:\n
  \                       cache.set(key, check)\n                    return check\n\n
  \               file_checks = [(file, check_spec(str(file))) for file in markata.files]\n
  \               [check.result() for _, check in file_checks]\n                markata.files
  = [file for file, check in file_checks if check]\n```\n\n\n!! method <h2 id='convert_to_list'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>convert_to_list <em
  class='small'>method</em></h2>\n\n???+ source \"convert_to_list <em class='small'>source</em>\"\n\n```python\n\n
  \       def convert_to_list(cls, v):\n                if not isinstance(v, list):\n
  \                   return v.split(\",\")\n                return v\n```\n\n\n!!
  function <h2 id='check_spec' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>check_spec <em class='small'>function</em></h2>\n\n???+ source \"check_spec
  <em class='small'>source</em>\"\n\n```python\n\n        def check_spec(file: str)
  -> bool:\n                    key = markata.make_hash(\"glob\", \"check_spec\",
  file)\n                    check = markata.precache.get(key)\n                    if
  check is not None:\n                        return check\n\n                    check
  = spec.match_file(str(file))\n                    with markata.cache as cache:\n
  \                       cache.set(key, check)\n                    return check\n```\n\n"
date: 0001-01-01
description: 'Default glob plugin ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
  source  ! ???+ source  ! ???+ source '
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Glob.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"Default glob plugin ! ???+ source  ! ???+ source
    \ ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Glob.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Default glob plugin ! ???+ source  !
    ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source \" />\n
    <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\"
    href=\"/post.css\" />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script
    src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    </head>\n    <body>\n<div class='container flex flex-row min-h-screen'>\n
    \   <div>\n    </div>\n    <div class='flex-grow px-8 mx-auto min-h-screen'>\n<header
    class='flex justify-center items-center p-8'>\n\n    <nav class='flex justify-center
    items-center my-8'>\n        <a\n            href='/'>markata</a>\n        <a\n
    \           href='https://github.com/WaylonWalker/markata'>GitHub</a>\n        <a\n
    \           href='https://markata.dev/docs/'>docs</a>\n        <a\n            href='https://markata.dev/plugins/'>plugins</a>\n
    \   </nav>\n\n    <div>\n        <label id=\"theme-switch\" class=\"theme-switch\"
    for=\"checkbox-theme\" title=\"light/dark mode toggle\">\n            <input type=\"checkbox\"
    id=\"checkbox-theme\" />\n            <div class=\"slider round\"></div>\n        </label>\n
    \   </div>\n</header><article class='w-full'>\n<section class=\"title\">\n    <h1
    id=\"title\">\n        Glob.Py\n    </h1>\n</section>    <section class=\"body\">\n
    \       <p>Default glob plugin</p>\n<p>!! class <h2 id='GlobConfig' class='admonition-title'
    style='margin:0;padding:.5rem 1rem;'>GlobConfig <em class='small'>class</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">GlobConfig
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
    <span class=\"nc\">GlobConfig</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">glob_patterns</span><span class=\"p\">:</span> <span
    class=\"n\">Union</span><span class=\"p\">[</span><span class=\"n\">List</span><span
    class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">],</span> <span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"p\">[</span><span class=\"s2\">&quot;**/*.md&quot;</span><span class=\"p\">]</span>\n
    \           <span class=\"n\">use_gitignore</span><span class=\"p\">:</span> <span
    class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"kc\">True</span>\n\n
    \           <span class=\"nd\">@pydantic</span><span class=\"o\">.</span><span
    class=\"n\">validator</span><span class=\"p\">(</span><span class=\"s2\">&quot;glob_patterns&quot;</span><span
    class=\"p\">)</span>\n            <span class=\"k\">def</span> <span class=\"nf\">convert_to_list</span><span
    class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
    class=\"n\">v</span><span class=\"p\">):</span>\n                <span class=\"k\">if</span>
    <span class=\"ow\">not</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"nb\">list</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">v</span><span class=\"o\">.</span><span class=\"n\">split</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;,&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    \           <span class=\"n\">glob</span><span class=\"p\">:</span> <span class=\"n\">GlobConfig</span>
    <span class=\"o\">=</span> <span class=\"n\">GlobConfig</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    function <h2 id='glob' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>glob <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">glob
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
    <span class=\"nf\">glob</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">files</span>
    <span class=\"o\">=</span> <span class=\"nb\">list</span><span class=\"p\">(</span>\n
    \               <span class=\"n\">flatten</span><span class=\"p\">(</span>\n                    <span
    class=\"p\">[</span>\n                        <span class=\"n\">Path</span><span
    class=\"p\">()</span><span class=\"o\">.</span><span class=\"n\">glob</span><span
    class=\"p\">(</span><span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">pattern</span><span class=\"p\">))</span>\n                        <span
    class=\"k\">for</span> <span class=\"n\">pattern</span> <span class=\"ow\">in</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">glob</span><span class=\"o\">.</span><span
    class=\"n\">glob_patterns</span>\n                    <span class=\"p\">],</span>\n
    \               <span class=\"p\">),</span>\n            <span class=\"p\">)</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">content_directories</span>
    <span class=\"o\">=</span> <span class=\"nb\">list</span><span class=\"p\">({</span><span
    class=\"n\">f</span><span class=\"o\">.</span><span class=\"n\">parent</span>
    <span class=\"k\">for</span> <span class=\"n\">f</span> <span class=\"ow\">in</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">files</span><span
    class=\"p\">})</span>\n\n            <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">ignore</span> <span class=\"o\">=</span> <span
    class=\"kc\">True</span>\n            <span class=\"k\">except</span> <span class=\"ne\">KeyError</span><span
    class=\"p\">:</span>\n                <span class=\"n\">ignore</span> <span class=\"o\">=</span>
    <span class=\"kc\">True</span>\n\n            <span class=\"k\">if</span> <span
    class=\"n\">ignore</span> <span class=\"ow\">and</span> <span class=\"p\">(</span><span
    class=\"n\">Path</span><span class=\"p\">(</span><span class=\"s2\">&quot;.gitignore&quot;</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">exists</span><span
    class=\"p\">()</span> <span class=\"ow\">or</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;.markataignore&quot;</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">exists</span><span
    class=\"p\">()):</span>\n                <span class=\"kn\">import</span> <span
    class=\"nn\">pathspec</span>\n\n                <span class=\"n\">lines</span>
    <span class=\"o\">=</span> <span class=\"p\">[]</span>\n\n                <span
    class=\"k\">if</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;.gitignore&quot;</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">exists</span><span class=\"p\">():</span>\n                    <span
    class=\"n\">lines</span><span class=\"o\">.</span><span class=\"n\">extend</span><span
    class=\"p\">(</span><span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;.gitignore&quot;</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">read_text</span><span class=\"p\">()</span><span class=\"o\">.</span><span
    class=\"n\">splitlines</span><span class=\"p\">())</span>\n\n                <span
    class=\"k\">if</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;.markataignore&quot;</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">():</span>\n
    \                   <span class=\"n\">lines</span><span class=\"o\">.</span><span
    class=\"n\">extend</span><span class=\"p\">(</span><span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;.markataignore&quot;</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">read_text</span><span
    class=\"p\">()</span><span class=\"o\">.</span><span class=\"n\">splitlines</span><span
    class=\"p\">())</span>\n\n                <span class=\"n\">key</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">make_hash</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;glob&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;spec&quot;</span><span class=\"p\">,</span> <span class=\"n\">lines</span><span
    class=\"p\">)</span>\n                <span class=\"n\">spec</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">precache</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">)</span>\n                <span class=\"k\">if</span>
    <span class=\"n\">spec</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">spec</span> <span
    class=\"o\">=</span> <span class=\"n\">pathspec</span><span class=\"o\">.</span><span
    class=\"n\">PathSpec</span><span class=\"o\">.</span><span class=\"n\">from_lines</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;gitwildmatch&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">lines</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">with</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">cache</span> <span class=\"k\">as</span> <span class=\"n\">cache</span><span
    class=\"p\">:</span>\n                        <span class=\"n\">cache</span><span
    class=\"o\">.</span><span class=\"n\">set</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">,</span> <span class=\"n\">spec</span><span
    class=\"p\">)</span>\n\n                <span class=\"nd\">@background</span><span
    class=\"o\">.</span><span class=\"n\">task</span>\n                <span class=\"k\">def</span>
    <span class=\"nf\">check_spec</span><span class=\"p\">(</span><span class=\"n\">file</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"nb\">bool</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">key</span> <span class=\"o\">=</span> <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">make_hash</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;glob&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;check_spec&quot;</span><span class=\"p\">,</span> <span
    class=\"n\">file</span><span class=\"p\">)</span>\n                    <span class=\"n\">check</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">precache</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"n\">key</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">if</span> <span class=\"n\">check</span> <span class=\"ow\">is</span>
    <span class=\"ow\">not</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                       <span class=\"k\">return</span> <span class=\"n\">check</span>\n\n
    \                   <span class=\"n\">check</span> <span class=\"o\">=</span>
    <span class=\"n\">spec</span><span class=\"o\">.</span><span class=\"n\">match_file</span><span
    class=\"p\">(</span><span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">file</span><span class=\"p\">))</span>\n                    <span
    class=\"k\">with</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">cache</span> <span class=\"k\">as</span> <span class=\"n\">cache</span><span
    class=\"p\">:</span>\n                        <span class=\"n\">cache</span><span
    class=\"o\">.</span><span class=\"n\">set</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">,</span> <span class=\"n\">check</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">check</span>\n\n                <span class=\"n\">file_checks</span>
    <span class=\"o\">=</span> <span class=\"p\">[(</span><span class=\"n\">file</span><span
    class=\"p\">,</span> <span class=\"n\">check_spec</span><span class=\"p\">(</span><span
    class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">file</span><span
    class=\"p\">)))</span> <span class=\"k\">for</span> <span class=\"n\">file</span>
    <span class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">files</span><span class=\"p\">]</span>\n                <span class=\"p\">[</span><span
    class=\"n\">check</span><span class=\"o\">.</span><span class=\"n\">result</span><span
    class=\"p\">()</span> <span class=\"k\">for</span> <span class=\"n\">_</span><span
    class=\"p\">,</span> <span class=\"n\">check</span> <span class=\"ow\">in</span>
    <span class=\"n\">file_checks</span><span class=\"p\">]</span>\n                <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">files</span>
    <span class=\"o\">=</span> <span class=\"p\">[</span><span class=\"n\">file</span>
    <span class=\"k\">for</span> <span class=\"n\">file</span><span class=\"p\">,</span>
    <span class=\"n\">check</span> <span class=\"ow\">in</span> <span class=\"n\">file_checks</span>
    <span class=\"k\">if</span> <span class=\"n\">check</span><span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='convert_to_list' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>convert_to_list <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">convert_to_list
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
    <span class=\"nf\">convert_to_list</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">):</span>\n                <span
    class=\"k\">if</span> <span class=\"ow\">not</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"nb\">list</span><span class=\"p\">):</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">v</span><span class=\"o\">.</span><span
    class=\"n\">split</span><span class=\"p\">(</span><span class=\"s2\">&quot;,&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='check_spec' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>check_spec <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">check_spec
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
    <span class=\"nf\">check_spec</span><span class=\"p\">(</span><span class=\"n\">file</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"nb\">bool</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">key</span> <span class=\"o\">=</span> <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">make_hash</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;glob&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;check_spec&quot;</span><span class=\"p\">,</span> <span
    class=\"n\">file</span><span class=\"p\">)</span>\n                    <span class=\"n\">check</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">precache</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"n\">key</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">if</span> <span class=\"n\">check</span> <span class=\"ow\">is</span>
    <span class=\"ow\">not</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                       <span class=\"k\">return</span> <span class=\"n\">check</span>\n\n
    \                   <span class=\"n\">check</span> <span class=\"o\">=</span>
    <span class=\"n\">spec</span><span class=\"o\">.</span><span class=\"n\">match_file</span><span
    class=\"p\">(</span><span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">file</span><span class=\"p\">))</span>\n                    <span
    class=\"k\">with</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">cache</span> <span class=\"k\">as</span> <span class=\"n\">cache</span><span
    class=\"p\">:</span>\n                        <span class=\"n\">cache</span><span
    class=\"o\">.</span><span class=\"n\">set</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">,</span> <span class=\"n\">check</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">check</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>
    \   </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Glob.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"Default glob plugin ! ???+ source  ! ???+ source
    \ ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Glob.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Default glob plugin ! ???+ source  !
    ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source \" />\n
    <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\"
    href=\"/post.css\" />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script
    src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    </head>\n    <body>\n<article style=\"text-align: center;\">\n    <style>\n
    \       section {\n            font-size: 200%;\n        }\n\n\n        .edit
    {\n            display: none;\n        }\n    </style>\n<section class=\"title\">\n
    \   <h1 id=\"title\">\n        Glob.Py\n    </h1>\n</section></article>\n     </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Glob.Py\n    </h1>\n</section>    <section class=\"body\">\n        <p>Default
    glob plugin</p>\n<p>!! class <h2 id='GlobConfig' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>GlobConfig <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">GlobConfig
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
    <span class=\"nc\">GlobConfig</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">glob_patterns</span><span class=\"p\">:</span> <span
    class=\"n\">Union</span><span class=\"p\">[</span><span class=\"n\">List</span><span
    class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">],</span> <span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"p\">[</span><span class=\"s2\">&quot;**/*.md&quot;</span><span class=\"p\">]</span>\n
    \           <span class=\"n\">use_gitignore</span><span class=\"p\">:</span> <span
    class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"kc\">True</span>\n\n
    \           <span class=\"nd\">@pydantic</span><span class=\"o\">.</span><span
    class=\"n\">validator</span><span class=\"p\">(</span><span class=\"s2\">&quot;glob_patterns&quot;</span><span
    class=\"p\">)</span>\n            <span class=\"k\">def</span> <span class=\"nf\">convert_to_list</span><span
    class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
    class=\"n\">v</span><span class=\"p\">):</span>\n                <span class=\"k\">if</span>
    <span class=\"ow\">not</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"nb\">list</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">v</span><span class=\"o\">.</span><span class=\"n\">split</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;,&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    \           <span class=\"n\">glob</span><span class=\"p\">:</span> <span class=\"n\">GlobConfig</span>
    <span class=\"o\">=</span> <span class=\"n\">GlobConfig</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    function <h2 id='glob' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>glob <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">glob
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
    <span class=\"nf\">glob</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">files</span>
    <span class=\"o\">=</span> <span class=\"nb\">list</span><span class=\"p\">(</span>\n
    \               <span class=\"n\">flatten</span><span class=\"p\">(</span>\n                    <span
    class=\"p\">[</span>\n                        <span class=\"n\">Path</span><span
    class=\"p\">()</span><span class=\"o\">.</span><span class=\"n\">glob</span><span
    class=\"p\">(</span><span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">pattern</span><span class=\"p\">))</span>\n                        <span
    class=\"k\">for</span> <span class=\"n\">pattern</span> <span class=\"ow\">in</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">glob</span><span class=\"o\">.</span><span
    class=\"n\">glob_patterns</span>\n                    <span class=\"p\">],</span>\n
    \               <span class=\"p\">),</span>\n            <span class=\"p\">)</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">content_directories</span>
    <span class=\"o\">=</span> <span class=\"nb\">list</span><span class=\"p\">({</span><span
    class=\"n\">f</span><span class=\"o\">.</span><span class=\"n\">parent</span>
    <span class=\"k\">for</span> <span class=\"n\">f</span> <span class=\"ow\">in</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">files</span><span
    class=\"p\">})</span>\n\n            <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">ignore</span> <span class=\"o\">=</span> <span
    class=\"kc\">True</span>\n            <span class=\"k\">except</span> <span class=\"ne\">KeyError</span><span
    class=\"p\">:</span>\n                <span class=\"n\">ignore</span> <span class=\"o\">=</span>
    <span class=\"kc\">True</span>\n\n            <span class=\"k\">if</span> <span
    class=\"n\">ignore</span> <span class=\"ow\">and</span> <span class=\"p\">(</span><span
    class=\"n\">Path</span><span class=\"p\">(</span><span class=\"s2\">&quot;.gitignore&quot;</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">exists</span><span
    class=\"p\">()</span> <span class=\"ow\">or</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;.markataignore&quot;</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">exists</span><span
    class=\"p\">()):</span>\n                <span class=\"kn\">import</span> <span
    class=\"nn\">pathspec</span>\n\n                <span class=\"n\">lines</span>
    <span class=\"o\">=</span> <span class=\"p\">[]</span>\n\n                <span
    class=\"k\">if</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;.gitignore&quot;</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">exists</span><span class=\"p\">():</span>\n                    <span
    class=\"n\">lines</span><span class=\"o\">.</span><span class=\"n\">extend</span><span
    class=\"p\">(</span><span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;.gitignore&quot;</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">read_text</span><span class=\"p\">()</span><span class=\"o\">.</span><span
    class=\"n\">splitlines</span><span class=\"p\">())</span>\n\n                <span
    class=\"k\">if</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;.markataignore&quot;</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">():</span>\n
    \                   <span class=\"n\">lines</span><span class=\"o\">.</span><span
    class=\"n\">extend</span><span class=\"p\">(</span><span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;.markataignore&quot;</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">read_text</span><span
    class=\"p\">()</span><span class=\"o\">.</span><span class=\"n\">splitlines</span><span
    class=\"p\">())</span>\n\n                <span class=\"n\">key</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">make_hash</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;glob&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;spec&quot;</span><span class=\"p\">,</span> <span class=\"n\">lines</span><span
    class=\"p\">)</span>\n                <span class=\"n\">spec</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">precache</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">)</span>\n                <span class=\"k\">if</span>
    <span class=\"n\">spec</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">spec</span> <span
    class=\"o\">=</span> <span class=\"n\">pathspec</span><span class=\"o\">.</span><span
    class=\"n\">PathSpec</span><span class=\"o\">.</span><span class=\"n\">from_lines</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;gitwildmatch&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">lines</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">with</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">cache</span> <span class=\"k\">as</span> <span class=\"n\">cache</span><span
    class=\"p\">:</span>\n                        <span class=\"n\">cache</span><span
    class=\"o\">.</span><span class=\"n\">set</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">,</span> <span class=\"n\">spec</span><span
    class=\"p\">)</span>\n\n                <span class=\"nd\">@background</span><span
    class=\"o\">.</span><span class=\"n\">task</span>\n                <span class=\"k\">def</span>
    <span class=\"nf\">check_spec</span><span class=\"p\">(</span><span class=\"n\">file</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"nb\">bool</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">key</span> <span class=\"o\">=</span> <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">make_hash</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;glob&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;check_spec&quot;</span><span class=\"p\">,</span> <span
    class=\"n\">file</span><span class=\"p\">)</span>\n                    <span class=\"n\">check</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">precache</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"n\">key</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">if</span> <span class=\"n\">check</span> <span class=\"ow\">is</span>
    <span class=\"ow\">not</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                       <span class=\"k\">return</span> <span class=\"n\">check</span>\n\n
    \                   <span class=\"n\">check</span> <span class=\"o\">=</span>
    <span class=\"n\">spec</span><span class=\"o\">.</span><span class=\"n\">match_file</span><span
    class=\"p\">(</span><span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">file</span><span class=\"p\">))</span>\n                    <span
    class=\"k\">with</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">cache</span> <span class=\"k\">as</span> <span class=\"n\">cache</span><span
    class=\"p\">:</span>\n                        <span class=\"n\">cache</span><span
    class=\"o\">.</span><span class=\"n\">set</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">,</span> <span class=\"n\">check</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">check</span>\n\n                <span class=\"n\">file_checks</span>
    <span class=\"o\">=</span> <span class=\"p\">[(</span><span class=\"n\">file</span><span
    class=\"p\">,</span> <span class=\"n\">check_spec</span><span class=\"p\">(</span><span
    class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">file</span><span
    class=\"p\">)))</span> <span class=\"k\">for</span> <span class=\"n\">file</span>
    <span class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">files</span><span class=\"p\">]</span>\n                <span class=\"p\">[</span><span
    class=\"n\">check</span><span class=\"o\">.</span><span class=\"n\">result</span><span
    class=\"p\">()</span> <span class=\"k\">for</span> <span class=\"n\">_</span><span
    class=\"p\">,</span> <span class=\"n\">check</span> <span class=\"ow\">in</span>
    <span class=\"n\">file_checks</span><span class=\"p\">]</span>\n                <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">files</span>
    <span class=\"o\">=</span> <span class=\"p\">[</span><span class=\"n\">file</span>
    <span class=\"k\">for</span> <span class=\"n\">file</span><span class=\"p\">,</span>
    <span class=\"n\">check</span> <span class=\"ow\">in</span> <span class=\"n\">file_checks</span>
    <span class=\"k\">if</span> <span class=\"n\">check</span><span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='convert_to_list' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>convert_to_list <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">convert_to_list
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
    <span class=\"nf\">convert_to_list</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">):</span>\n                <span
    class=\"k\">if</span> <span class=\"ow\">not</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"nb\">list</span><span class=\"p\">):</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">v</span><span class=\"o\">.</span><span
    class=\"n\">split</span><span class=\"p\">(</span><span class=\"s2\">&quot;,&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='check_spec' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>check_spec <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">check_spec
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
    <span class=\"nf\">check_spec</span><span class=\"p\">(</span><span class=\"n\">file</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"nb\">bool</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">key</span> <span class=\"o\">=</span> <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">make_hash</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;glob&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;check_spec&quot;</span><span class=\"p\">,</span> <span
    class=\"n\">file</span><span class=\"p\">)</span>\n                    <span class=\"n\">check</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">precache</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"n\">key</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">if</span> <span class=\"n\">check</span> <span class=\"ow\">is</span>
    <span class=\"ow\">not</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                       <span class=\"k\">return</span> <span class=\"n\">check</span>\n\n
    \                   <span class=\"n\">check</span> <span class=\"o\">=</span>
    <span class=\"n\">spec</span><span class=\"o\">.</span><span class=\"n\">match_file</span><span
    class=\"p\">(</span><span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">file</span><span class=\"p\">))</span>\n                    <span
    class=\"k\">with</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">cache</span> <span class=\"k\">as</span> <span class=\"n\">cache</span><span
    class=\"p\">:</span>\n                        <span class=\"n\">cache</span><span
    class=\"o\">.</span><span class=\"n\">set</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">,</span> <span class=\"n\">check</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">check</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>"
  raw.md: ''
published: true
slug: markata/plugins/glob
title: Glob.Py


---

Default glob plugin


!! class <h2 id='GlobConfig' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>GlobConfig <em class='small'>class</em></h2>

???+ source "GlobConfig <em class='small'>source</em>"

```python

        class GlobConfig(pydantic.BaseModel):
            glob_patterns: Union[List[str], str] = ["**/*.md"]
            use_gitignore: bool = True

            @pydantic.validator("glob_patterns")
            def convert_to_list(cls, v):
                if not isinstance(v, list):
                    return v.split(",")
                return v
```


!! class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Config <em class='small'>class</em></h2>

???+ source "Config <em class='small'>source</em>"

```python

        class Config(pydantic.BaseModel):
            glob: GlobConfig = GlobConfig()
```


!! function <h2 id='config_model' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>config_model <em class='small'>function</em></h2>

???+ source "config_model <em class='small'>source</em>"

```python

        def config_model(markata: "Markata") -> None:
            markata.config_models.append(Config)
```


!! function <h2 id='glob' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>glob <em class='small'>function</em></h2>

???+ source "glob <em class='small'>source</em>"

```python

        def glob(markata: "Markata") -> None:
            markata.files = list(
                flatten(
                    [
                        Path().glob(str(pattern))
                        for pattern in markata.config.glob.glob_patterns
                    ],
                ),
            )
            markata.content_directories = list({f.parent for f in markata.files})

            try:
                ignore = True
            except KeyError:
                ignore = True

            if ignore and (Path(".gitignore").exists() or Path(".markataignore").exists()):
                import pathspec

                lines = []

                if Path(".gitignore").exists():
                    lines.extend(Path(".gitignore").read_text().splitlines())

                if Path(".markataignore").exists():
                    lines.extend(Path(".markataignore").read_text().splitlines())

                key = markata.make_hash("glob", "spec", lines)
                spec = markata.precache.get(key)
                if spec is None:
                    spec = pathspec.PathSpec.from_lines("gitwildmatch", lines)
                    with markata.cache as cache:
                        cache.set(key, spec)

                @background.task
                def check_spec(file: str) -> bool:
                    key = markata.make_hash("glob", "check_spec", file)
                    check = markata.precache.get(key)
                    if check is not None:
                        return check

                    check = spec.match_file(str(file))
                    with markata.cache as cache:
                        cache.set(key, check)
                    return check

                file_checks = [(file, check_spec(str(file))) for file in markata.files]
                [check.result() for _, check in file_checks]
                markata.files = [file for file, check in file_checks if check]
```


!! method <h2 id='convert_to_list' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>convert_to_list <em class='small'>method</em></h2>

???+ source "convert_to_list <em class='small'>source</em>"

```python

        def convert_to_list(cls, v):
                if not isinstance(v, list):
                    return v.split(",")
                return v
```


!! function <h2 id='check_spec' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>check_spec <em class='small'>function</em></h2>

???+ source "check_spec <em class='small'>source</em>"

```python

        def check_spec(file: str) -> bool:
                    key = markata.make_hash("glob", "check_spec", file)
                    check = markata.precache.get(key)
                    if check is not None:
                        return check

                    check = spec.match_file(str(file))
                    with markata.cache as cache:
                        cache.set(key, check)
                    return check
```

