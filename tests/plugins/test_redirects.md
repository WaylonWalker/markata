---
content: "Tests the redirects plugin\n\n\n!! function <h2 id='set_directory' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>set_directory <em class='small'>function</em></h2>\n
  \   context manager to set the directory\n???+ source \"set_directory <em class='small'>source</em>\"\n\n```python\n\n
  \       def set_directory(path):\n            \"\"\"\n            context manager
  to set the directory\n            \"\"\"\n            cwd = Path.cwd()\n            try:\n
  \               os.chdir(path)\n                yield\n            finally:\n                os.chdir(cwd)\n```\n\n\n!!
  function <h2 id='test_redirect_exists' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>test_redirect_exists <em class='small'>function</em></h2>\n    ensure that
  the default workflow works\n???+ source \"test_redirect_exists <em class='small'>source</em>\"\n\n```python\n\n
  \       def test_redirect_exists(tmp_files: Path, old: str, new: str) -> None:\n
  \           \"ensure that the default workflow works\"\n            with set_directory(tmp_files):\n
  \               m = Markata()\n                redirects.save(m)\n                redirect_file
  = Path(\"markout\") / old / \"index.html\"\n                assert redirect_file.exists()\n
  \               assert (\n                    f'<meta http-equiv=\"Refresh\" content=\"0;
  url=\\'{new}\\'\" />'\n                    in redirect_file.read_text()\n                )\n```\n\n\n!!
  function <h2 id='test_redirect_ignore_splat' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>test_redirect_ignore_splat <em class='small'>function</em></h2>\n    splats
  cannot be supported statically, test that they are ignored\n???+ source \"test_redirect_ignore_splat
  <em class='small'>source</em>\"\n\n```python\n\n        def test_redirect_ignore_splat(tmp_files:
  Path, old: str) -> None:\n            \"splats cannot be supported statically, test
  that they are ignored\"\n            with set_directory(tmp_files):\n                m
  = Markata()\n                redirects.save(m)\n                redirect_file =
  Path(\"markout\") / old / \"index.html\"\n                assert not redirect_file.exists()\n```\n\n\n!!
  function <h2 id='test_redirect_ignore_more_params' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>test_redirect_ignore_more_params <em class='small'>function</em></h2>\n    status
  codes cannot be supported statically as they are issued by the server\n???+ source
  \"test_redirect_ignore_more_params <em class='small'>source</em>\"\n\n```python\n\n
  \       def test_redirect_ignore_more_params(tmp_files: Path, old: str) -> None:\n
  \           \"status codes cannot be supported statically as they are issued by
  the server\"\n            with set_directory(tmp_files):\n                m = Markata()\n
  \               redirects.save(m)\n                redirect_file = Path(\"markout\")
  / old / \"index.html\"\n                assert not redirect_file.exists()\n```\n\n\n!!
  function <h2 id='test_redirect_configure_redirect_file' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>test_redirect_configure_redirect_file <em class='small'>function</em></h2>\n
  \   ensure that the redirects file can be configured\n???+ source \"test_redirect_configure_redirect_file
  <em class='small'>source</em>\"\n\n```python\n\n        def test_redirect_configure_redirect_file(\n
  \           tmp_files: Path, redirect_file: str, old: str, new: str\n        ) ->
  None:\n            \"ensure that the redirects file can be configured\"\n            with
  set_directory(tmp_files):\n                m = Markata()\n                m.config[\"redirects\"]
  = redirect_file\n                redirects.save(m)\n                redirect_html
  = Path(\"markout\") / old / \"index.html\"\n                assert redirect_html.exists()\n
  \               assert (\n                    f'<meta http-equiv=\"Refresh\" content=\"0;
  url=\\'{new}\\'\" />'\n                    in redirect_html.read_text()\n                )\n```\n\n\n!!
  function <h2 id='test_redirect_custom_template' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>test_redirect_custom_template <em class='small'>function</em></h2>\n    ensure
  the template can be configured\n???+ source \"test_redirect_custom_template <em
  class='small'>source</em>\"\n\n```python\n\n        def test_redirect_custom_template(\n
  \           tmp_files: Path, redirect_template: str, old: str, new: str\n        )
  -> None:\n            \"ensure the template can be configured\"\n            with
  set_directory(tmp_files):\n                m = Markata()\n                m.config[\"redirect_template\"]
  = redirect_template\n                redirects.save(m)\n                redirect_file
  = Path(\"markout\") / old / \"index.html\"\n                assert redirect_file.exists()\n
  \               assert f\"{old} is now {new}\" in redirect_file.read_text()\n```\n\n\n!!
  function <h2 id='test_redirect_empty' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>test_redirect_empty <em class='small'>function</em></h2>\n    ensure empty
  redirects files work\n???+ source \"test_redirect_empty <em class='small'>source</em>\"\n\n```python\n\n
  \       def test_redirect_empty(tmp_files: Path, old: str, new: str) -> None:\n
  \           \"ensure empty redirects files work\"\n            with set_directory(tmp_files):\n
  \               m = Markata()\n                redirects.save(m)\n```\n\n\n!! function
  <h2 id='test_redirect_file_missing' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>test_redirect_file_missing <em class='small'>function</em></h2>\n    ensure
  missing redirects file works\n???+ source \"test_redirect_file_missing <em class='small'>source</em>\"\n\n```python\n\n
  \       def test_redirect_file_missing(tmpdir: Path) -> None:\n            \"ensure
  missing redirects file works\"\n            with set_directory(tmpdir):\n                m
  = Markata()\n                redirects.save(m)\n```\n\n"
date: 0001-01-01
description: Tests the redirects plugin ! ! ! ! ! ! ! !
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Test_Redirects.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Tests the redirects plugin ! ! ! ! !
    ! ! !\" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\" />\n\n<link
    rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\" href=\"/app.css\"
    />\n<script src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    <head>\n<title>Test_Redirects.Py</title>\n<meta charset=\"UTF-8\" />\n<meta
    name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta name=\"description\"
    content=\"Tests the redirects plugin ! ! ! ! ! ! ! !\" />\n <link href=\"/favicon.ico\"
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Test_Redirects.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <p>Tests the redirects
    plugin</p>\n<p>!! function <h2 id='set_directory' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>set_directory <em class='small'>function</em></h2>\ncontext manager to
    set the directory</p>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">set_directory <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">set_directory</span><span class=\"p\">(</span><span class=\"n\">path</span><span
    class=\"p\">):</span>\n<span class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            context manager to set the directory</span>\n<span class=\"sd\">
    \           &quot;&quot;&quot;</span>\n            <span class=\"n\">cwd</span>
    <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"o\">.</span><span
    class=\"n\">cwd</span><span class=\"p\">()</span>\n            <span class=\"k\">try</span><span
    class=\"p\">:</span>\n                <span class=\"n\">os</span><span class=\"o\">.</span><span
    class=\"n\">chdir</span><span class=\"p\">(</span><span class=\"n\">path</span><span
    class=\"p\">)</span>\n                <span class=\"k\">yield</span>\n            <span
    class=\"k\">finally</span><span class=\"p\">:</span>\n                <span class=\"n\">os</span><span
    class=\"o\">.</span><span class=\"n\">chdir</span><span class=\"p\">(</span><span
    class=\"n\">cwd</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='test_redirect_exists' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>test_redirect_exists <em class='small'>function</em></h2>\nensure that
    the default workflow works</p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">test_redirect_exists <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">test_redirect_exists</span><span class=\"p\">(</span><span
    class=\"n\">tmp_files</span><span class=\"p\">:</span> <span class=\"n\">Path</span><span
    class=\"p\">,</span> <span class=\"n\">old</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span><span class=\"p\">,</span> <span class=\"n\">new</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"s2\">&quot;ensure that the default workflow works&quot;</span>\n
    \           <span class=\"k\">with</span> <span class=\"n\">set_directory</span><span
    class=\"p\">(</span><span class=\"n\">tmp_files</span><span class=\"p\">):</span>\n
    \               <span class=\"n\">m</span> <span class=\"o\">=</span> <span class=\"n\">Markata</span><span
    class=\"p\">()</span>\n                <span class=\"n\">redirects</span><span
    class=\"o\">.</span><span class=\"n\">save</span><span class=\"p\">(</span><span
    class=\"n\">m</span><span class=\"p\">)</span>\n                <span class=\"n\">redirect_file</span>
    <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;markout&quot;</span><span class=\"p\">)</span> <span class=\"o\">/</span>
    <span class=\"n\">old</span> <span class=\"o\">/</span> <span class=\"s2\">&quot;index.html&quot;</span>\n
    \               <span class=\"k\">assert</span> <span class=\"n\">redirect_file</span><span
    class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">()</span>\n
    \               <span class=\"k\">assert</span> <span class=\"p\">(</span>\n                    <span
    class=\"sa\">f</span><span class=\"s1\">&#39;&lt;meta http-equiv=&quot;Refresh&quot;
    content=&quot;0; url=</span><span class=\"se\">\\&#39;</span><span class=\"si\">{</span><span
    class=\"n\">new</span><span class=\"si\">}</span><span class=\"se\">\\&#39;</span><span
    class=\"s1\">&quot; /&gt;&#39;</span>\n                    <span class=\"ow\">in</span>
    <span class=\"n\">redirect_file</span><span class=\"o\">.</span><span class=\"n\">read_text</span><span
    class=\"p\">()</span>\n                <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='test_redirect_ignore_splat' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>test_redirect_ignore_splat <em class='small'>function</em></h2>\nsplats
    cannot be supported statically, test that they are ignored</p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">test_redirect_ignore_splat
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
    <span class=\"nf\">test_redirect_ignore_splat</span><span class=\"p\">(</span><span
    class=\"n\">tmp_files</span><span class=\"p\">:</span> <span class=\"n\">Path</span><span
    class=\"p\">,</span> <span class=\"n\">old</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"s2\">&quot;splats
    cannot be supported statically, test that they are ignored&quot;</span>\n            <span
    class=\"k\">with</span> <span class=\"n\">set_directory</span><span class=\"p\">(</span><span
    class=\"n\">tmp_files</span><span class=\"p\">):</span>\n                <span
    class=\"n\">m</span> <span class=\"o\">=</span> <span class=\"n\">Markata</span><span
    class=\"p\">()</span>\n                <span class=\"n\">redirects</span><span
    class=\"o\">.</span><span class=\"n\">save</span><span class=\"p\">(</span><span
    class=\"n\">m</span><span class=\"p\">)</span>\n                <span class=\"n\">redirect_file</span>
    <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;markout&quot;</span><span class=\"p\">)</span> <span class=\"o\">/</span>
    <span class=\"n\">old</span> <span class=\"o\">/</span> <span class=\"s2\">&quot;index.html&quot;</span>\n
    \               <span class=\"k\">assert</span> <span class=\"ow\">not</span>
    <span class=\"n\">redirect_file</span><span class=\"o\">.</span><span class=\"n\">exists</span><span
    class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='test_redirect_ignore_more_params'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>test_redirect_ignore_more_params
    <em class='small'>function</em></h2>\nstatus codes cannot be supported statically
    as they are issued by the server</p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">test_redirect_ignore_more_params
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
    <span class=\"nf\">test_redirect_ignore_more_params</span><span class=\"p\">(</span><span
    class=\"n\">tmp_files</span><span class=\"p\">:</span> <span class=\"n\">Path</span><span
    class=\"p\">,</span> <span class=\"n\">old</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"s2\">&quot;status
    codes cannot be supported statically as they are issued by the server&quot;</span>\n
    \           <span class=\"k\">with</span> <span class=\"n\">set_directory</span><span
    class=\"p\">(</span><span class=\"n\">tmp_files</span><span class=\"p\">):</span>\n
    \               <span class=\"n\">m</span> <span class=\"o\">=</span> <span class=\"n\">Markata</span><span
    class=\"p\">()</span>\n                <span class=\"n\">redirects</span><span
    class=\"o\">.</span><span class=\"n\">save</span><span class=\"p\">(</span><span
    class=\"n\">m</span><span class=\"p\">)</span>\n                <span class=\"n\">redirect_file</span>
    <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;markout&quot;</span><span class=\"p\">)</span> <span class=\"o\">/</span>
    <span class=\"n\">old</span> <span class=\"o\">/</span> <span class=\"s2\">&quot;index.html&quot;</span>\n
    \               <span class=\"k\">assert</span> <span class=\"ow\">not</span>
    <span class=\"n\">redirect_file</span><span class=\"o\">.</span><span class=\"n\">exists</span><span
    class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='test_redirect_configure_redirect_file'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>test_redirect_configure_redirect_file
    <em class='small'>function</em></h2>\nensure that the redirects file can be configured</p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">test_redirect_configure_redirect_file
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
    <span class=\"nf\">test_redirect_configure_redirect_file</span><span class=\"p\">(</span>\n
    \           <span class=\"n\">tmp_files</span><span class=\"p\">:</span> <span
    class=\"n\">Path</span><span class=\"p\">,</span> <span class=\"n\">redirect_file</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">,</span> <span
    class=\"n\">old</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"n\">new</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span>\n        <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"s2\">&quot;ensure
    that the redirects file can be configured&quot;</span>\n            <span class=\"k\">with</span>
    <span class=\"n\">set_directory</span><span class=\"p\">(</span><span class=\"n\">tmp_files</span><span
    class=\"p\">):</span>\n                <span class=\"n\">m</span> <span class=\"o\">=</span>
    <span class=\"n\">Markata</span><span class=\"p\">()</span>\n                <span
    class=\"n\">m</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;redirects&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"n\">redirect_file</span>\n                <span
    class=\"n\">redirects</span><span class=\"o\">.</span><span class=\"n\">save</span><span
    class=\"p\">(</span><span class=\"n\">m</span><span class=\"p\">)</span>\n                <span
    class=\"n\">redirect_html</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;markout&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">/</span> <span class=\"n\">old</span> <span class=\"o\">/</span>
    <span class=\"s2\">&quot;index.html&quot;</span>\n                <span class=\"k\">assert</span>
    <span class=\"n\">redirect_html</span><span class=\"o\">.</span><span class=\"n\">exists</span><span
    class=\"p\">()</span>\n                <span class=\"k\">assert</span> <span class=\"p\">(</span>\n
    \                   <span class=\"sa\">f</span><span class=\"s1\">&#39;&lt;meta
    http-equiv=&quot;Refresh&quot; content=&quot;0; url=</span><span class=\"se\">\\&#39;</span><span
    class=\"si\">{</span><span class=\"n\">new</span><span class=\"si\">}</span><span
    class=\"se\">\\&#39;</span><span class=\"s1\">&quot; /&gt;&#39;</span>\n                    <span
    class=\"ow\">in</span> <span class=\"n\">redirect_html</span><span class=\"o\">.</span><span
    class=\"n\">read_text</span><span class=\"p\">()</span>\n                <span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='test_redirect_custom_template'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>test_redirect_custom_template
    <em class='small'>function</em></h2>\nensure the template can be configured</p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">test_redirect_custom_template
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
    <span class=\"nf\">test_redirect_custom_template</span><span class=\"p\">(</span>\n
    \           <span class=\"n\">tmp_files</span><span class=\"p\">:</span> <span
    class=\"n\">Path</span><span class=\"p\">,</span> <span class=\"n\">redirect_template</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">,</span> <span
    class=\"n\">old</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"n\">new</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span>\n        <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"s2\">&quot;ensure
    the template can be configured&quot;</span>\n            <span class=\"k\">with</span>
    <span class=\"n\">set_directory</span><span class=\"p\">(</span><span class=\"n\">tmp_files</span><span
    class=\"p\">):</span>\n                <span class=\"n\">m</span> <span class=\"o\">=</span>
    <span class=\"n\">Markata</span><span class=\"p\">()</span>\n                <span
    class=\"n\">m</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;redirect_template&quot;</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"n\">redirect_template</span>\n
    \               <span class=\"n\">redirects</span><span class=\"o\">.</span><span
    class=\"n\">save</span><span class=\"p\">(</span><span class=\"n\">m</span><span
    class=\"p\">)</span>\n                <span class=\"n\">redirect_file</span> <span
    class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;markout&quot;</span><span class=\"p\">)</span> <span class=\"o\">/</span>
    <span class=\"n\">old</span> <span class=\"o\">/</span> <span class=\"s2\">&quot;index.html&quot;</span>\n
    \               <span class=\"k\">assert</span> <span class=\"n\">redirect_file</span><span
    class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">()</span>\n
    \               <span class=\"k\">assert</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;</span><span class=\"si\">{</span><span class=\"n\">old</span><span
    class=\"si\">}</span><span class=\"s2\"> is now </span><span class=\"si\">{</span><span
    class=\"n\">new</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span>
    <span class=\"ow\">in</span> <span class=\"n\">redirect_file</span><span class=\"o\">.</span><span
    class=\"n\">read_text</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='test_redirect_empty' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>test_redirect_empty <em class='small'>function</em></h2>\nensure empty
    redirects files work</p>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">test_redirect_empty <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">test_redirect_empty</span><span class=\"p\">(</span><span class=\"n\">tmp_files</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">,</span> <span
    class=\"n\">old</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"n\">new</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"s2\">&quot;ensure
    empty redirects files work&quot;</span>\n            <span class=\"k\">with</span>
    <span class=\"n\">set_directory</span><span class=\"p\">(</span><span class=\"n\">tmp_files</span><span
    class=\"p\">):</span>\n                <span class=\"n\">m</span> <span class=\"o\">=</span>
    <span class=\"n\">Markata</span><span class=\"p\">()</span>\n                <span
    class=\"n\">redirects</span><span class=\"o\">.</span><span class=\"n\">save</span><span
    class=\"p\">(</span><span class=\"n\">m</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='test_redirect_file_missing' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>test_redirect_file_missing <em class='small'>function</em></h2>\nensure
    missing redirects file works</p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">test_redirect_file_missing
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
    <span class=\"nf\">test_redirect_file_missing</span><span class=\"p\">(</span><span
    class=\"n\">tmpdir</span><span class=\"p\">:</span> <span class=\"n\">Path</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n            <span class=\"s2\">&quot;ensure missing redirects
    file works&quot;</span>\n            <span class=\"k\">with</span> <span class=\"n\">set_directory</span><span
    class=\"p\">(</span><span class=\"n\">tmpdir</span><span class=\"p\">):</span>\n
    \               <span class=\"n\">m</span> <span class=\"o\">=</span> <span class=\"n\">Markata</span><span
    class=\"p\">()</span>\n                <span class=\"n\">redirects</span><span
    class=\"o\">.</span><span class=\"n\">save</span><span class=\"p\">(</span><span
    class=\"n\">m</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>
    \   </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Test_Redirects.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Tests the redirects plugin ! ! ! ! !
    ! ! !\" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\" />\n\n<link
    rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\" href=\"/app.css\"
    />\n<script src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    <head>\n<title>Test_Redirects.Py</title>\n<meta charset=\"UTF-8\" />\n<meta
    name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta name=\"description\"
    content=\"Tests the redirects plugin ! ! ! ! ! ! ! !\" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<article style=\"text-align: center;\">\n    <style>\n        section
    {\n            font-size: 200%;\n        }\n\n\n        .edit {\n            display:
    none;\n        }\n    </style>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Test_Redirects.Py\n    </h1>\n</section></article>\n     </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Test_Redirects.Py\n    </h1>\n</section>    <section class=\"body\">\n
    \       <p>Tests the redirects plugin</p>\n<p>!! function <h2 id='set_directory'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>set_directory <em
    class='small'>function</em></h2>\ncontext manager to set the directory</p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">set_directory
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
    <span class=\"nf\">set_directory</span><span class=\"p\">(</span><span class=\"n\">path</span><span
    class=\"p\">):</span>\n<span class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            context manager to set the directory</span>\n<span class=\"sd\">
    \           &quot;&quot;&quot;</span>\n            <span class=\"n\">cwd</span>
    <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"o\">.</span><span
    class=\"n\">cwd</span><span class=\"p\">()</span>\n            <span class=\"k\">try</span><span
    class=\"p\">:</span>\n                <span class=\"n\">os</span><span class=\"o\">.</span><span
    class=\"n\">chdir</span><span class=\"p\">(</span><span class=\"n\">path</span><span
    class=\"p\">)</span>\n                <span class=\"k\">yield</span>\n            <span
    class=\"k\">finally</span><span class=\"p\">:</span>\n                <span class=\"n\">os</span><span
    class=\"o\">.</span><span class=\"n\">chdir</span><span class=\"p\">(</span><span
    class=\"n\">cwd</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='test_redirect_exists' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>test_redirect_exists <em class='small'>function</em></h2>\nensure that
    the default workflow works</p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">test_redirect_exists <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">test_redirect_exists</span><span class=\"p\">(</span><span
    class=\"n\">tmp_files</span><span class=\"p\">:</span> <span class=\"n\">Path</span><span
    class=\"p\">,</span> <span class=\"n\">old</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span><span class=\"p\">,</span> <span class=\"n\">new</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"s2\">&quot;ensure that the default workflow works&quot;</span>\n
    \           <span class=\"k\">with</span> <span class=\"n\">set_directory</span><span
    class=\"p\">(</span><span class=\"n\">tmp_files</span><span class=\"p\">):</span>\n
    \               <span class=\"n\">m</span> <span class=\"o\">=</span> <span class=\"n\">Markata</span><span
    class=\"p\">()</span>\n                <span class=\"n\">redirects</span><span
    class=\"o\">.</span><span class=\"n\">save</span><span class=\"p\">(</span><span
    class=\"n\">m</span><span class=\"p\">)</span>\n                <span class=\"n\">redirect_file</span>
    <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;markout&quot;</span><span class=\"p\">)</span> <span class=\"o\">/</span>
    <span class=\"n\">old</span> <span class=\"o\">/</span> <span class=\"s2\">&quot;index.html&quot;</span>\n
    \               <span class=\"k\">assert</span> <span class=\"n\">redirect_file</span><span
    class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">()</span>\n
    \               <span class=\"k\">assert</span> <span class=\"p\">(</span>\n                    <span
    class=\"sa\">f</span><span class=\"s1\">&#39;&lt;meta http-equiv=&quot;Refresh&quot;
    content=&quot;0; url=</span><span class=\"se\">\\&#39;</span><span class=\"si\">{</span><span
    class=\"n\">new</span><span class=\"si\">}</span><span class=\"se\">\\&#39;</span><span
    class=\"s1\">&quot; /&gt;&#39;</span>\n                    <span class=\"ow\">in</span>
    <span class=\"n\">redirect_file</span><span class=\"o\">.</span><span class=\"n\">read_text</span><span
    class=\"p\">()</span>\n                <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='test_redirect_ignore_splat' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>test_redirect_ignore_splat <em class='small'>function</em></h2>\nsplats
    cannot be supported statically, test that they are ignored</p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">test_redirect_ignore_splat
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
    <span class=\"nf\">test_redirect_ignore_splat</span><span class=\"p\">(</span><span
    class=\"n\">tmp_files</span><span class=\"p\">:</span> <span class=\"n\">Path</span><span
    class=\"p\">,</span> <span class=\"n\">old</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"s2\">&quot;splats
    cannot be supported statically, test that they are ignored&quot;</span>\n            <span
    class=\"k\">with</span> <span class=\"n\">set_directory</span><span class=\"p\">(</span><span
    class=\"n\">tmp_files</span><span class=\"p\">):</span>\n                <span
    class=\"n\">m</span> <span class=\"o\">=</span> <span class=\"n\">Markata</span><span
    class=\"p\">()</span>\n                <span class=\"n\">redirects</span><span
    class=\"o\">.</span><span class=\"n\">save</span><span class=\"p\">(</span><span
    class=\"n\">m</span><span class=\"p\">)</span>\n                <span class=\"n\">redirect_file</span>
    <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;markout&quot;</span><span class=\"p\">)</span> <span class=\"o\">/</span>
    <span class=\"n\">old</span> <span class=\"o\">/</span> <span class=\"s2\">&quot;index.html&quot;</span>\n
    \               <span class=\"k\">assert</span> <span class=\"ow\">not</span>
    <span class=\"n\">redirect_file</span><span class=\"o\">.</span><span class=\"n\">exists</span><span
    class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='test_redirect_ignore_more_params'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>test_redirect_ignore_more_params
    <em class='small'>function</em></h2>\nstatus codes cannot be supported statically
    as they are issued by the server</p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">test_redirect_ignore_more_params
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
    <span class=\"nf\">test_redirect_ignore_more_params</span><span class=\"p\">(</span><span
    class=\"n\">tmp_files</span><span class=\"p\">:</span> <span class=\"n\">Path</span><span
    class=\"p\">,</span> <span class=\"n\">old</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"s2\">&quot;status
    codes cannot be supported statically as they are issued by the server&quot;</span>\n
    \           <span class=\"k\">with</span> <span class=\"n\">set_directory</span><span
    class=\"p\">(</span><span class=\"n\">tmp_files</span><span class=\"p\">):</span>\n
    \               <span class=\"n\">m</span> <span class=\"o\">=</span> <span class=\"n\">Markata</span><span
    class=\"p\">()</span>\n                <span class=\"n\">redirects</span><span
    class=\"o\">.</span><span class=\"n\">save</span><span class=\"p\">(</span><span
    class=\"n\">m</span><span class=\"p\">)</span>\n                <span class=\"n\">redirect_file</span>
    <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;markout&quot;</span><span class=\"p\">)</span> <span class=\"o\">/</span>
    <span class=\"n\">old</span> <span class=\"o\">/</span> <span class=\"s2\">&quot;index.html&quot;</span>\n
    \               <span class=\"k\">assert</span> <span class=\"ow\">not</span>
    <span class=\"n\">redirect_file</span><span class=\"o\">.</span><span class=\"n\">exists</span><span
    class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='test_redirect_configure_redirect_file'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>test_redirect_configure_redirect_file
    <em class='small'>function</em></h2>\nensure that the redirects file can be configured</p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">test_redirect_configure_redirect_file
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
    <span class=\"nf\">test_redirect_configure_redirect_file</span><span class=\"p\">(</span>\n
    \           <span class=\"n\">tmp_files</span><span class=\"p\">:</span> <span
    class=\"n\">Path</span><span class=\"p\">,</span> <span class=\"n\">redirect_file</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">,</span> <span
    class=\"n\">old</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"n\">new</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span>\n        <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"s2\">&quot;ensure
    that the redirects file can be configured&quot;</span>\n            <span class=\"k\">with</span>
    <span class=\"n\">set_directory</span><span class=\"p\">(</span><span class=\"n\">tmp_files</span><span
    class=\"p\">):</span>\n                <span class=\"n\">m</span> <span class=\"o\">=</span>
    <span class=\"n\">Markata</span><span class=\"p\">()</span>\n                <span
    class=\"n\">m</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;redirects&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"n\">redirect_file</span>\n                <span
    class=\"n\">redirects</span><span class=\"o\">.</span><span class=\"n\">save</span><span
    class=\"p\">(</span><span class=\"n\">m</span><span class=\"p\">)</span>\n                <span
    class=\"n\">redirect_html</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;markout&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">/</span> <span class=\"n\">old</span> <span class=\"o\">/</span>
    <span class=\"s2\">&quot;index.html&quot;</span>\n                <span class=\"k\">assert</span>
    <span class=\"n\">redirect_html</span><span class=\"o\">.</span><span class=\"n\">exists</span><span
    class=\"p\">()</span>\n                <span class=\"k\">assert</span> <span class=\"p\">(</span>\n
    \                   <span class=\"sa\">f</span><span class=\"s1\">&#39;&lt;meta
    http-equiv=&quot;Refresh&quot; content=&quot;0; url=</span><span class=\"se\">\\&#39;</span><span
    class=\"si\">{</span><span class=\"n\">new</span><span class=\"si\">}</span><span
    class=\"se\">\\&#39;</span><span class=\"s1\">&quot; /&gt;&#39;</span>\n                    <span
    class=\"ow\">in</span> <span class=\"n\">redirect_html</span><span class=\"o\">.</span><span
    class=\"n\">read_text</span><span class=\"p\">()</span>\n                <span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='test_redirect_custom_template'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>test_redirect_custom_template
    <em class='small'>function</em></h2>\nensure the template can be configured</p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">test_redirect_custom_template
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
    <span class=\"nf\">test_redirect_custom_template</span><span class=\"p\">(</span>\n
    \           <span class=\"n\">tmp_files</span><span class=\"p\">:</span> <span
    class=\"n\">Path</span><span class=\"p\">,</span> <span class=\"n\">redirect_template</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">,</span> <span
    class=\"n\">old</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"n\">new</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span>\n        <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"s2\">&quot;ensure
    the template can be configured&quot;</span>\n            <span class=\"k\">with</span>
    <span class=\"n\">set_directory</span><span class=\"p\">(</span><span class=\"n\">tmp_files</span><span
    class=\"p\">):</span>\n                <span class=\"n\">m</span> <span class=\"o\">=</span>
    <span class=\"n\">Markata</span><span class=\"p\">()</span>\n                <span
    class=\"n\">m</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;redirect_template&quot;</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"n\">redirect_template</span>\n
    \               <span class=\"n\">redirects</span><span class=\"o\">.</span><span
    class=\"n\">save</span><span class=\"p\">(</span><span class=\"n\">m</span><span
    class=\"p\">)</span>\n                <span class=\"n\">redirect_file</span> <span
    class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;markout&quot;</span><span class=\"p\">)</span> <span class=\"o\">/</span>
    <span class=\"n\">old</span> <span class=\"o\">/</span> <span class=\"s2\">&quot;index.html&quot;</span>\n
    \               <span class=\"k\">assert</span> <span class=\"n\">redirect_file</span><span
    class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">()</span>\n
    \               <span class=\"k\">assert</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;</span><span class=\"si\">{</span><span class=\"n\">old</span><span
    class=\"si\">}</span><span class=\"s2\"> is now </span><span class=\"si\">{</span><span
    class=\"n\">new</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span>
    <span class=\"ow\">in</span> <span class=\"n\">redirect_file</span><span class=\"o\">.</span><span
    class=\"n\">read_text</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='test_redirect_empty' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>test_redirect_empty <em class='small'>function</em></h2>\nensure empty
    redirects files work</p>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">test_redirect_empty <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">test_redirect_empty</span><span class=\"p\">(</span><span class=\"n\">tmp_files</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">,</span> <span
    class=\"n\">old</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"n\">new</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"s2\">&quot;ensure
    empty redirects files work&quot;</span>\n            <span class=\"k\">with</span>
    <span class=\"n\">set_directory</span><span class=\"p\">(</span><span class=\"n\">tmp_files</span><span
    class=\"p\">):</span>\n                <span class=\"n\">m</span> <span class=\"o\">=</span>
    <span class=\"n\">Markata</span><span class=\"p\">()</span>\n                <span
    class=\"n\">redirects</span><span class=\"o\">.</span><span class=\"n\">save</span><span
    class=\"p\">(</span><span class=\"n\">m</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='test_redirect_file_missing' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>test_redirect_file_missing <em class='small'>function</em></h2>\nensure
    missing redirects file works</p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">test_redirect_file_missing
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
    <span class=\"nf\">test_redirect_file_missing</span><span class=\"p\">(</span><span
    class=\"n\">tmpdir</span><span class=\"p\">:</span> <span class=\"n\">Path</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n            <span class=\"s2\">&quot;ensure missing redirects
    file works&quot;</span>\n            <span class=\"k\">with</span> <span class=\"n\">set_directory</span><span
    class=\"p\">(</span><span class=\"n\">tmpdir</span><span class=\"p\">):</span>\n
    \               <span class=\"n\">m</span> <span class=\"o\">=</span> <span class=\"n\">Markata</span><span
    class=\"p\">()</span>\n                <span class=\"n\">redirects</span><span
    class=\"o\">.</span><span class=\"n\">save</span><span class=\"p\">(</span><span
    class=\"n\">m</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>"
  raw.md: ''
published: true
slug: tests/plugins/test-redirects
title: Test_Redirects.Py


---

Tests the redirects plugin


!! function <h2 id='set_directory' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>set_directory <em class='small'>function</em></h2>
    context manager to set the directory
???+ source "set_directory <em class='small'>source</em>"

```python

        def set_directory(path):
            """
            context manager to set the directory
            """
            cwd = Path.cwd()
            try:
                os.chdir(path)
                yield
            finally:
                os.chdir(cwd)
```


!! function <h2 id='test_redirect_exists' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>test_redirect_exists <em class='small'>function</em></h2>
    ensure that the default workflow works
???+ source "test_redirect_exists <em class='small'>source</em>"

```python

        def test_redirect_exists(tmp_files: Path, old: str, new: str) -> None:
            "ensure that the default workflow works"
            with set_directory(tmp_files):
                m = Markata()
                redirects.save(m)
                redirect_file = Path("markout") / old / "index.html"
                assert redirect_file.exists()
                assert (
                    f'<meta http-equiv="Refresh" content="0; url=\'{new}\'" />'
                    in redirect_file.read_text()
                )
```


!! function <h2 id='test_redirect_ignore_splat' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>test_redirect_ignore_splat <em class='small'>function</em></h2>
    splats cannot be supported statically, test that they are ignored
???+ source "test_redirect_ignore_splat <em class='small'>source</em>"

```python

        def test_redirect_ignore_splat(tmp_files: Path, old: str) -> None:
            "splats cannot be supported statically, test that they are ignored"
            with set_directory(tmp_files):
                m = Markata()
                redirects.save(m)
                redirect_file = Path("markout") / old / "index.html"
                assert not redirect_file.exists()
```


!! function <h2 id='test_redirect_ignore_more_params' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>test_redirect_ignore_more_params <em class='small'>function</em></h2>
    status codes cannot be supported statically as they are issued by the server
???+ source "test_redirect_ignore_more_params <em class='small'>source</em>"

```python

        def test_redirect_ignore_more_params(tmp_files: Path, old: str) -> None:
            "status codes cannot be supported statically as they are issued by the server"
            with set_directory(tmp_files):
                m = Markata()
                redirects.save(m)
                redirect_file = Path("markout") / old / "index.html"
                assert not redirect_file.exists()
```


!! function <h2 id='test_redirect_configure_redirect_file' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>test_redirect_configure_redirect_file <em class='small'>function</em></h2>
    ensure that the redirects file can be configured
???+ source "test_redirect_configure_redirect_file <em class='small'>source</em>"

```python

        def test_redirect_configure_redirect_file(
            tmp_files: Path, redirect_file: str, old: str, new: str
        ) -> None:
            "ensure that the redirects file can be configured"
            with set_directory(tmp_files):
                m = Markata()
                m.config["redirects"] = redirect_file
                redirects.save(m)
                redirect_html = Path("markout") / old / "index.html"
                assert redirect_html.exists()
                assert (
                    f'<meta http-equiv="Refresh" content="0; url=\'{new}\'" />'
                    in redirect_html.read_text()
                )
```


!! function <h2 id='test_redirect_custom_template' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>test_redirect_custom_template <em class='small'>function</em></h2>
    ensure the template can be configured
???+ source "test_redirect_custom_template <em class='small'>source</em>"

```python

        def test_redirect_custom_template(
            tmp_files: Path, redirect_template: str, old: str, new: str
        ) -> None:
            "ensure the template can be configured"
            with set_directory(tmp_files):
                m = Markata()
                m.config["redirect_template"] = redirect_template
                redirects.save(m)
                redirect_file = Path("markout") / old / "index.html"
                assert redirect_file.exists()
                assert f"{old} is now {new}" in redirect_file.read_text()
```


!! function <h2 id='test_redirect_empty' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>test_redirect_empty <em class='small'>function</em></h2>
    ensure empty redirects files work
???+ source "test_redirect_empty <em class='small'>source</em>"

```python

        def test_redirect_empty(tmp_files: Path, old: str, new: str) -> None:
            "ensure empty redirects files work"
            with set_directory(tmp_files):
                m = Markata()
                redirects.save(m)
```


!! function <h2 id='test_redirect_file_missing' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>test_redirect_file_missing <em class='small'>function</em></h2>
    ensure missing redirects file works
???+ source "test_redirect_file_missing <em class='small'>source</em>"

```python

        def test_redirect_file_missing(tmpdir: Path) -> None:
            "ensure missing redirects file works"
            with set_directory(tmpdir):
                m = Markata()
                redirects.save(m)
```

