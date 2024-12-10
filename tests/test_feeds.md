---
content: "None\n\n\n!! class <h2 id='DummyMarkata' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>DummyMarkata <em class='small'>class</em></h2>\n\n???+ source \"DummyMarkata
  <em class='small'>source</em>\"\n\n```python\n\n        class DummyMarkata:\n            def
  __init__(self):\n                self.config = {\"feeds\": {\"all-posts\": {\"filter\":
  True}}}\n\n            def map(self, *args, **kwargs):\n                return []\n```\n\n\n!!
  function <h2 id='test_feed_map' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>test_feed_map <em class='small'>function</em></h2>\n\n???+ source \"test_feed_map
  <em class='small'>source</em>\"\n\n```python\n\n        def test_feed_map(mocker):\n
  \           mocked_markata = mocker.patch.object(markata, \"Markata\", DummyMarkata())\n
  \           feed = Feed(name=\"archive\", config={}, posts=[], _m=mocked_markata)\n
  \           feed.map()\n```\n\n\n!! function <h2 id='test_feeds' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>test_feeds <em class='small'>function</em></h2>\n\n???+
  source \"test_feeds <em class='small'>source</em>\"\n\n```python\n\n        def
  test_feeds(mocker):\n            mocked_markata = mocker.patch.object(markata, \"Markata\",
  DummyMarkata())\n            Feeds(markata=mocked_markata)\n```\n\n\n!! function
  <h2 id='test_feeds_iter' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>test_feeds_iter <em class='small'>function</em></h2>\n\n???+ source \"test_feeds_iter
  <em class='small'>source</em>\"\n\n```python\n\n        def test_feeds_iter(mocker):\n
  \           mocked_markata = mocker.patch.object(markata, \"Markata\", DummyMarkata())\n
  \           feeds = Feeds(markata=mocked_markata)\n\n            assert [feed for
  feed in feeds] == [\"all-posts\"]\n```\n\n\n!! function <h2 id='test_feeds_keys'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>test_feeds_keys <em
  class='small'>function</em></h2>\n\n???+ source \"test_feeds_keys <em class='small'>source</em>\"\n\n```python\n\n
  \       def test_feeds_keys(mocker):\n            mocked_markata = mocker.patch.object(markata,
  \"Markata\", DummyMarkata())\n            feeds = Feeds(markata=mocked_markata)\n\n
  \           assert [feed for feed in feeds.keys()] == [\"all-posts\"]\n```\n\n\n!!
  function <h2 id='test_feeds_items' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>test_feeds_items <em class='small'>function</em></h2>\n\n???+ source \"test_feeds_items
  <em class='small'>source</em>\"\n\n```python\n\n        def test_feeds_items(mocker):\n
  \           mocked_markata = mocker.patch.object(markata, \"Markata\", DummyMarkata())\n
  \           feeds = Feeds(markata=mocked_markata)\n            iterate_feeds = {name:
  feed for name, feed in feeds.items()}\n\n            assert list(iterate_feeds.keys())
  == [\"all-posts\"]\n```\n\n\n!! function <h2 id='test_feeds___get_item__' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>test_feeds___get_item__ <em class='small'>function</em></h2>\n\n???+
  source \"test_feeds___get_item__ <em class='small'>source</em>\"\n\n```python\n\n
  \       def test_feeds___get_item__(mocker):\n            mocked_markata = mocker.patch.object(markata,
  \"Markata\", DummyMarkata())\n            feeds = Feeds(markata=mocked_markata)\n\n
  \           feeds[\"all-posts\"].name == \"all_posts\"\n```\n\n\n!! function <h2
  id='test_feeds___rich__' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>test_feeds___rich__ <em class='small'>function</em></h2>\n\n???+ source \"test_feeds___rich__
  <em class='small'>source</em>\"\n\n```python\n\n        def test_feeds___rich__(mocker):\n
  \           mocked_markata = mocker.patch.object(markata, \"Markata\", DummyMarkata())\n
  \           feeds = Feeds(markata=mocked_markata)\n\n            assert isinstance(feeds.__rich__(),
  rich.table.Table)\n```\n\n\n!! function <h2 id='test_feeds_values' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>test_feeds_values <em class='small'>function</em></h2>\n\n???+
  source \"test_feeds_values <em class='small'>source</em>\"\n\n```python\n\n        def
  test_feeds_values(mocker):\n            mocked_markata = mocker.patch.object(markata,
  \"Markata\", DummyMarkata())\n            values = Feeds(markata=mocked_markata).values()\n
  \           assert all([isinstance(v, Feed) for v in values])\n```\n\n\n!! method
  <h2 id='__init__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__init__
  <em class='small'>method</em></h2>\n\n???+ source \"__init__ <em class='small'>source</em>\"\n\n```python\n\n
  \       def __init__(self):\n                self.config = {\"feeds\": {\"all-posts\":
  {\"filter\": True}}}\n```\n\n\n!! method <h2 id='map' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>map <em class='small'>method</em></h2>\n\n???+ source \"map <em class='small'>source</em>\"\n\n```python\n\n
  \       def map(self, *args, **kwargs):\n                return []\n```\n\n"
date: 0001-01-01
description: None ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
  source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  !
  ???
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Test_Feeds.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None ! ???+ source  ! ???+ source  !
    ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
    source  ! ???+ source  ! ???+ source  ! ???\" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Test_Feeds.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None ! ???+ source  ! ???+ source  !
    ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
    source  ! ???+ source  ! ???+ source  ! ???\" />\n <link href=\"/favicon.ico\"
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Test_Feeds.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <p>None</p>\n<p>!!
    class <h2 id='DummyMarkata' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>DummyMarkata <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">DummyMarkata
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
    <span class=\"nc\">DummyMarkata</span><span class=\"p\">:</span>\n            <span
    class=\"k\">def</span> <span class=\"fm\">__init__</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">):</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">config</span> <span class=\"o\">=</span>
    <span class=\"p\">{</span><span class=\"s2\">&quot;feeds&quot;</span><span class=\"p\">:</span>
    <span class=\"p\">{</span><span class=\"s2\">&quot;all-posts&quot;</span><span
    class=\"p\">:</span> <span class=\"p\">{</span><span class=\"s2\">&quot;filter&quot;</span><span
    class=\"p\">:</span> <span class=\"kc\">True</span><span class=\"p\">}}}</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"nf\">map</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"n\">args</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"n\">kwargs</span><span class=\"p\">):</span>\n                <span class=\"k\">return</span>
    <span class=\"p\">[]</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='test_feed_map'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>test_feed_map <em
    class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">test_feed_map <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">test_feed_map</span><span class=\"p\">(</span><span class=\"n\">mocker</span><span
    class=\"p\">):</span>\n            <span class=\"n\">mocked_markata</span> <span
    class=\"o\">=</span> <span class=\"n\">mocker</span><span class=\"o\">.</span><span
    class=\"n\">patch</span><span class=\"o\">.</span><span class=\"n\">object</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span> <span
    class=\"n\">DummyMarkata</span><span class=\"p\">())</span>\n            <span
    class=\"n\">feed</span> <span class=\"o\">=</span> <span class=\"n\">Feed</span><span
    class=\"p\">(</span><span class=\"n\">name</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;archive&quot;</span><span class=\"p\">,</span> <span class=\"n\">config</span><span
    class=\"o\">=</span><span class=\"p\">{},</span> <span class=\"n\">posts</span><span
    class=\"o\">=</span><span class=\"p\">[],</span> <span class=\"n\">_m</span><span
    class=\"o\">=</span><span class=\"n\">mocked_markata</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">feed</span><span class=\"o\">.</span><span class=\"n\">map</span><span
    class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='test_feeds'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>test_feeds <em class='small'>function</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">test_feeds
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
    <span class=\"nf\">test_feeds</span><span class=\"p\">(</span><span class=\"n\">mocker</span><span
    class=\"p\">):</span>\n            <span class=\"n\">mocked_markata</span> <span
    class=\"o\">=</span> <span class=\"n\">mocker</span><span class=\"o\">.</span><span
    class=\"n\">patch</span><span class=\"o\">.</span><span class=\"n\">object</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span> <span
    class=\"n\">DummyMarkata</span><span class=\"p\">())</span>\n            <span
    class=\"n\">Feeds</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"n\">mocked_markata</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='test_feeds_iter' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>test_feeds_iter <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">test_feeds_iter
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
    <span class=\"nf\">test_feeds_iter</span><span class=\"p\">(</span><span class=\"n\">mocker</span><span
    class=\"p\">):</span>\n            <span class=\"n\">mocked_markata</span> <span
    class=\"o\">=</span> <span class=\"n\">mocker</span><span class=\"o\">.</span><span
    class=\"n\">patch</span><span class=\"o\">.</span><span class=\"n\">object</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span> <span
    class=\"n\">DummyMarkata</span><span class=\"p\">())</span>\n            <span
    class=\"n\">feeds</span> <span class=\"o\">=</span> <span class=\"n\">Feeds</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"n\">mocked_markata</span><span class=\"p\">)</span>\n\n            <span
    class=\"k\">assert</span> <span class=\"p\">[</span><span class=\"n\">feed</span>
    <span class=\"k\">for</span> <span class=\"n\">feed</span> <span class=\"ow\">in</span>
    <span class=\"n\">feeds</span><span class=\"p\">]</span> <span class=\"o\">==</span>
    <span class=\"p\">[</span><span class=\"s2\">&quot;all-posts&quot;</span><span
    class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='test_feeds_keys'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>test_feeds_keys
    <em class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">test_feeds_keys <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">test_feeds_keys</span><span class=\"p\">(</span><span class=\"n\">mocker</span><span
    class=\"p\">):</span>\n            <span class=\"n\">mocked_markata</span> <span
    class=\"o\">=</span> <span class=\"n\">mocker</span><span class=\"o\">.</span><span
    class=\"n\">patch</span><span class=\"o\">.</span><span class=\"n\">object</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span> <span
    class=\"n\">DummyMarkata</span><span class=\"p\">())</span>\n            <span
    class=\"n\">feeds</span> <span class=\"o\">=</span> <span class=\"n\">Feeds</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"n\">mocked_markata</span><span class=\"p\">)</span>\n\n            <span
    class=\"k\">assert</span> <span class=\"p\">[</span><span class=\"n\">feed</span>
    <span class=\"k\">for</span> <span class=\"n\">feed</span> <span class=\"ow\">in</span>
    <span class=\"n\">feeds</span><span class=\"o\">.</span><span class=\"n\">keys</span><span
    class=\"p\">()]</span> <span class=\"o\">==</span> <span class=\"p\">[</span><span
    class=\"s2\">&quot;all-posts&quot;</span><span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='test_feeds_items' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>test_feeds_items <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">test_feeds_items
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
    <span class=\"nf\">test_feeds_items</span><span class=\"p\">(</span><span class=\"n\">mocker</span><span
    class=\"p\">):</span>\n            <span class=\"n\">mocked_markata</span> <span
    class=\"o\">=</span> <span class=\"n\">mocker</span><span class=\"o\">.</span><span
    class=\"n\">patch</span><span class=\"o\">.</span><span class=\"n\">object</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span> <span
    class=\"n\">DummyMarkata</span><span class=\"p\">())</span>\n            <span
    class=\"n\">feeds</span> <span class=\"o\">=</span> <span class=\"n\">Feeds</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"n\">mocked_markata</span><span class=\"p\">)</span>\n            <span
    class=\"n\">iterate_feeds</span> <span class=\"o\">=</span> <span class=\"p\">{</span><span
    class=\"n\">name</span><span class=\"p\">:</span> <span class=\"n\">feed</span>
    <span class=\"k\">for</span> <span class=\"n\">name</span><span class=\"p\">,</span>
    <span class=\"n\">feed</span> <span class=\"ow\">in</span> <span class=\"n\">feeds</span><span
    class=\"o\">.</span><span class=\"n\">items</span><span class=\"p\">()}</span>\n\n
    \           <span class=\"k\">assert</span> <span class=\"nb\">list</span><span
    class=\"p\">(</span><span class=\"n\">iterate_feeds</span><span class=\"o\">.</span><span
    class=\"n\">keys</span><span class=\"p\">())</span> <span class=\"o\">==</span>
    <span class=\"p\">[</span><span class=\"s2\">&quot;all-posts&quot;</span><span
    class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='test_feeds___get_item__'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>test_feeds___get_item__
    <em class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">test_feeds___get_item__ <em
    class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    <span class=\"nf\">test_feeds___get_item__</span><span class=\"p\">(</span><span
    class=\"n\">mocker</span><span class=\"p\">):</span>\n            <span class=\"n\">mocked_markata</span>
    <span class=\"o\">=</span> <span class=\"n\">mocker</span><span class=\"o\">.</span><span
    class=\"n\">patch</span><span class=\"o\">.</span><span class=\"n\">object</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span> <span
    class=\"n\">DummyMarkata</span><span class=\"p\">())</span>\n            <span
    class=\"n\">feeds</span> <span class=\"o\">=</span> <span class=\"n\">Feeds</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"n\">mocked_markata</span><span class=\"p\">)</span>\n\n            <span
    class=\"n\">feeds</span><span class=\"p\">[</span><span class=\"s2\">&quot;all-posts&quot;</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">name</span> <span
    class=\"o\">==</span> <span class=\"s2\">&quot;all_posts&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='test_feeds___rich__' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>test_feeds___rich__ <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">test_feeds___rich__
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
    <span class=\"nf\">test_feeds___rich__</span><span class=\"p\">(</span><span class=\"n\">mocker</span><span
    class=\"p\">):</span>\n            <span class=\"n\">mocked_markata</span> <span
    class=\"o\">=</span> <span class=\"n\">mocker</span><span class=\"o\">.</span><span
    class=\"n\">patch</span><span class=\"o\">.</span><span class=\"n\">object</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span> <span
    class=\"n\">DummyMarkata</span><span class=\"p\">())</span>\n            <span
    class=\"n\">feeds</span> <span class=\"o\">=</span> <span class=\"n\">Feeds</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"n\">mocked_markata</span><span class=\"p\">)</span>\n\n            <span
    class=\"k\">assert</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">feeds</span><span class=\"o\">.</span><span class=\"n\">__rich__</span><span
    class=\"p\">(),</span> <span class=\"n\">rich</span><span class=\"o\">.</span><span
    class=\"n\">table</span><span class=\"o\">.</span><span class=\"n\">Table</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='test_feeds_values'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>test_feeds_values
    <em class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">test_feeds_values <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">test_feeds_values</span><span class=\"p\">(</span><span class=\"n\">mocker</span><span
    class=\"p\">):</span>\n            <span class=\"n\">mocked_markata</span> <span
    class=\"o\">=</span> <span class=\"n\">mocker</span><span class=\"o\">.</span><span
    class=\"n\">patch</span><span class=\"o\">.</span><span class=\"n\">object</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span> <span
    class=\"n\">DummyMarkata</span><span class=\"p\">())</span>\n            <span
    class=\"n\">values</span> <span class=\"o\">=</span> <span class=\"n\">Feeds</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"n\">mocked_markata</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">values</span><span class=\"p\">()</span>\n            <span class=\"k\">assert</span>
    <span class=\"nb\">all</span><span class=\"p\">([</span><span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"n\">Feed</span><span class=\"p\">)</span> <span class=\"k\">for</span>
    <span class=\"n\">v</span> <span class=\"ow\">in</span> <span class=\"n\">values</span><span
    class=\"p\">])</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='__init__'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'><strong>init</strong>
    <em class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\"><strong>init</strong> <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"fm\">__init__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">):</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">config</span> <span class=\"o\">=</span> <span class=\"p\">{</span><span
    class=\"s2\">&quot;feeds&quot;</span><span class=\"p\">:</span> <span class=\"p\">{</span><span
    class=\"s2\">&quot;all-posts&quot;</span><span class=\"p\">:</span> <span class=\"p\">{</span><span
    class=\"s2\">&quot;filter&quot;</span><span class=\"p\">:</span> <span class=\"kc\">True</span><span
    class=\"p\">}}}</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='map' class='admonition-title'
    style='margin:0;padding:.5rem 1rem;'>map <em class='small'>method</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">map
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
    <span class=\"nf\">map</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"o\">*</span><span class=\"n\">args</span><span
    class=\"p\">,</span> <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">):</span>\n                <span class=\"k\">return</span> <span class=\"p\">[]</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>    </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Test_Feeds.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"None ! ???+ source  ! ???+ source  ! ???+ source
    \ ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  !
    ???+ source  ! ???+ source  ! ???\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Test_Feeds.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None ! ???+ source  ! ???+ source  !
    ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
    source  ! ???+ source  ! ???+ source  ! ???\" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<article style=\"text-align: center;\">\n    <style>\n        section
    {\n            font-size: 200%;\n        }\n\n\n        .edit {\n            display:
    none;\n        }\n    </style>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Test_Feeds.Py\n    </h1>\n</section></article>\n     </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Test_Feeds.Py\n    </h1>\n</section>    <section class=\"body\">\n        <p>None</p>\n<p>!!
    class <h2 id='DummyMarkata' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>DummyMarkata <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">DummyMarkata
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
    <span class=\"nc\">DummyMarkata</span><span class=\"p\">:</span>\n            <span
    class=\"k\">def</span> <span class=\"fm\">__init__</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">):</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">config</span> <span class=\"o\">=</span>
    <span class=\"p\">{</span><span class=\"s2\">&quot;feeds&quot;</span><span class=\"p\">:</span>
    <span class=\"p\">{</span><span class=\"s2\">&quot;all-posts&quot;</span><span
    class=\"p\">:</span> <span class=\"p\">{</span><span class=\"s2\">&quot;filter&quot;</span><span
    class=\"p\">:</span> <span class=\"kc\">True</span><span class=\"p\">}}}</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"nf\">map</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"n\">args</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"n\">kwargs</span><span class=\"p\">):</span>\n                <span class=\"k\">return</span>
    <span class=\"p\">[]</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='test_feed_map'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>test_feed_map <em
    class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">test_feed_map <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">test_feed_map</span><span class=\"p\">(</span><span class=\"n\">mocker</span><span
    class=\"p\">):</span>\n            <span class=\"n\">mocked_markata</span> <span
    class=\"o\">=</span> <span class=\"n\">mocker</span><span class=\"o\">.</span><span
    class=\"n\">patch</span><span class=\"o\">.</span><span class=\"n\">object</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span> <span
    class=\"n\">DummyMarkata</span><span class=\"p\">())</span>\n            <span
    class=\"n\">feed</span> <span class=\"o\">=</span> <span class=\"n\">Feed</span><span
    class=\"p\">(</span><span class=\"n\">name</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;archive&quot;</span><span class=\"p\">,</span> <span class=\"n\">config</span><span
    class=\"o\">=</span><span class=\"p\">{},</span> <span class=\"n\">posts</span><span
    class=\"o\">=</span><span class=\"p\">[],</span> <span class=\"n\">_m</span><span
    class=\"o\">=</span><span class=\"n\">mocked_markata</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">feed</span><span class=\"o\">.</span><span class=\"n\">map</span><span
    class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='test_feeds'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>test_feeds <em class='small'>function</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">test_feeds
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
    <span class=\"nf\">test_feeds</span><span class=\"p\">(</span><span class=\"n\">mocker</span><span
    class=\"p\">):</span>\n            <span class=\"n\">mocked_markata</span> <span
    class=\"o\">=</span> <span class=\"n\">mocker</span><span class=\"o\">.</span><span
    class=\"n\">patch</span><span class=\"o\">.</span><span class=\"n\">object</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span> <span
    class=\"n\">DummyMarkata</span><span class=\"p\">())</span>\n            <span
    class=\"n\">Feeds</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"n\">mocked_markata</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='test_feeds_iter' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>test_feeds_iter <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">test_feeds_iter
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
    <span class=\"nf\">test_feeds_iter</span><span class=\"p\">(</span><span class=\"n\">mocker</span><span
    class=\"p\">):</span>\n            <span class=\"n\">mocked_markata</span> <span
    class=\"o\">=</span> <span class=\"n\">mocker</span><span class=\"o\">.</span><span
    class=\"n\">patch</span><span class=\"o\">.</span><span class=\"n\">object</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span> <span
    class=\"n\">DummyMarkata</span><span class=\"p\">())</span>\n            <span
    class=\"n\">feeds</span> <span class=\"o\">=</span> <span class=\"n\">Feeds</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"n\">mocked_markata</span><span class=\"p\">)</span>\n\n            <span
    class=\"k\">assert</span> <span class=\"p\">[</span><span class=\"n\">feed</span>
    <span class=\"k\">for</span> <span class=\"n\">feed</span> <span class=\"ow\">in</span>
    <span class=\"n\">feeds</span><span class=\"p\">]</span> <span class=\"o\">==</span>
    <span class=\"p\">[</span><span class=\"s2\">&quot;all-posts&quot;</span><span
    class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='test_feeds_keys'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>test_feeds_keys
    <em class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">test_feeds_keys <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">test_feeds_keys</span><span class=\"p\">(</span><span class=\"n\">mocker</span><span
    class=\"p\">):</span>\n            <span class=\"n\">mocked_markata</span> <span
    class=\"o\">=</span> <span class=\"n\">mocker</span><span class=\"o\">.</span><span
    class=\"n\">patch</span><span class=\"o\">.</span><span class=\"n\">object</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span> <span
    class=\"n\">DummyMarkata</span><span class=\"p\">())</span>\n            <span
    class=\"n\">feeds</span> <span class=\"o\">=</span> <span class=\"n\">Feeds</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"n\">mocked_markata</span><span class=\"p\">)</span>\n\n            <span
    class=\"k\">assert</span> <span class=\"p\">[</span><span class=\"n\">feed</span>
    <span class=\"k\">for</span> <span class=\"n\">feed</span> <span class=\"ow\">in</span>
    <span class=\"n\">feeds</span><span class=\"o\">.</span><span class=\"n\">keys</span><span
    class=\"p\">()]</span> <span class=\"o\">==</span> <span class=\"p\">[</span><span
    class=\"s2\">&quot;all-posts&quot;</span><span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='test_feeds_items' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>test_feeds_items <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">test_feeds_items
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
    <span class=\"nf\">test_feeds_items</span><span class=\"p\">(</span><span class=\"n\">mocker</span><span
    class=\"p\">):</span>\n            <span class=\"n\">mocked_markata</span> <span
    class=\"o\">=</span> <span class=\"n\">mocker</span><span class=\"o\">.</span><span
    class=\"n\">patch</span><span class=\"o\">.</span><span class=\"n\">object</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span> <span
    class=\"n\">DummyMarkata</span><span class=\"p\">())</span>\n            <span
    class=\"n\">feeds</span> <span class=\"o\">=</span> <span class=\"n\">Feeds</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"n\">mocked_markata</span><span class=\"p\">)</span>\n            <span
    class=\"n\">iterate_feeds</span> <span class=\"o\">=</span> <span class=\"p\">{</span><span
    class=\"n\">name</span><span class=\"p\">:</span> <span class=\"n\">feed</span>
    <span class=\"k\">for</span> <span class=\"n\">name</span><span class=\"p\">,</span>
    <span class=\"n\">feed</span> <span class=\"ow\">in</span> <span class=\"n\">feeds</span><span
    class=\"o\">.</span><span class=\"n\">items</span><span class=\"p\">()}</span>\n\n
    \           <span class=\"k\">assert</span> <span class=\"nb\">list</span><span
    class=\"p\">(</span><span class=\"n\">iterate_feeds</span><span class=\"o\">.</span><span
    class=\"n\">keys</span><span class=\"p\">())</span> <span class=\"o\">==</span>
    <span class=\"p\">[</span><span class=\"s2\">&quot;all-posts&quot;</span><span
    class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='test_feeds___get_item__'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>test_feeds___get_item__
    <em class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">test_feeds___get_item__ <em
    class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    <span class=\"nf\">test_feeds___get_item__</span><span class=\"p\">(</span><span
    class=\"n\">mocker</span><span class=\"p\">):</span>\n            <span class=\"n\">mocked_markata</span>
    <span class=\"o\">=</span> <span class=\"n\">mocker</span><span class=\"o\">.</span><span
    class=\"n\">patch</span><span class=\"o\">.</span><span class=\"n\">object</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span> <span
    class=\"n\">DummyMarkata</span><span class=\"p\">())</span>\n            <span
    class=\"n\">feeds</span> <span class=\"o\">=</span> <span class=\"n\">Feeds</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"n\">mocked_markata</span><span class=\"p\">)</span>\n\n            <span
    class=\"n\">feeds</span><span class=\"p\">[</span><span class=\"s2\">&quot;all-posts&quot;</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">name</span> <span
    class=\"o\">==</span> <span class=\"s2\">&quot;all_posts&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='test_feeds___rich__' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>test_feeds___rich__ <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">test_feeds___rich__
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
    <span class=\"nf\">test_feeds___rich__</span><span class=\"p\">(</span><span class=\"n\">mocker</span><span
    class=\"p\">):</span>\n            <span class=\"n\">mocked_markata</span> <span
    class=\"o\">=</span> <span class=\"n\">mocker</span><span class=\"o\">.</span><span
    class=\"n\">patch</span><span class=\"o\">.</span><span class=\"n\">object</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span> <span
    class=\"n\">DummyMarkata</span><span class=\"p\">())</span>\n            <span
    class=\"n\">feeds</span> <span class=\"o\">=</span> <span class=\"n\">Feeds</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"n\">mocked_markata</span><span class=\"p\">)</span>\n\n            <span
    class=\"k\">assert</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">feeds</span><span class=\"o\">.</span><span class=\"n\">__rich__</span><span
    class=\"p\">(),</span> <span class=\"n\">rich</span><span class=\"o\">.</span><span
    class=\"n\">table</span><span class=\"o\">.</span><span class=\"n\">Table</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='test_feeds_values'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>test_feeds_values
    <em class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">test_feeds_values <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">test_feeds_values</span><span class=\"p\">(</span><span class=\"n\">mocker</span><span
    class=\"p\">):</span>\n            <span class=\"n\">mocked_markata</span> <span
    class=\"o\">=</span> <span class=\"n\">mocker</span><span class=\"o\">.</span><span
    class=\"n\">patch</span><span class=\"o\">.</span><span class=\"n\">object</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span> <span
    class=\"n\">DummyMarkata</span><span class=\"p\">())</span>\n            <span
    class=\"n\">values</span> <span class=\"o\">=</span> <span class=\"n\">Feeds</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"n\">mocked_markata</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">values</span><span class=\"p\">()</span>\n            <span class=\"k\">assert</span>
    <span class=\"nb\">all</span><span class=\"p\">([</span><span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"n\">Feed</span><span class=\"p\">)</span> <span class=\"k\">for</span>
    <span class=\"n\">v</span> <span class=\"ow\">in</span> <span class=\"n\">values</span><span
    class=\"p\">])</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='__init__'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'><strong>init</strong>
    <em class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\"><strong>init</strong> <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"fm\">__init__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">):</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">config</span> <span class=\"o\">=</span> <span class=\"p\">{</span><span
    class=\"s2\">&quot;feeds&quot;</span><span class=\"p\">:</span> <span class=\"p\">{</span><span
    class=\"s2\">&quot;all-posts&quot;</span><span class=\"p\">:</span> <span class=\"p\">{</span><span
    class=\"s2\">&quot;filter&quot;</span><span class=\"p\">:</span> <span class=\"kc\">True</span><span
    class=\"p\">}}}</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='map' class='admonition-title'
    style='margin:0;padding:.5rem 1rem;'>map <em class='small'>method</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">map
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
    <span class=\"nf\">map</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"o\">*</span><span class=\"n\">args</span><span
    class=\"p\">,</span> <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">):</span>\n                <span class=\"k\">return</span> <span class=\"p\">[]</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>"
  raw.md: ''
published: true
slug: tests/test-feeds
title: Test_Feeds.Py


---

None


!! class <h2 id='DummyMarkata' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>DummyMarkata <em class='small'>class</em></h2>

???+ source "DummyMarkata <em class='small'>source</em>"

```python

        class DummyMarkata:
            def __init__(self):
                self.config = {"feeds": {"all-posts": {"filter": True}}}

            def map(self, *args, **kwargs):
                return []
```


!! function <h2 id='test_feed_map' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>test_feed_map <em class='small'>function</em></h2>

???+ source "test_feed_map <em class='small'>source</em>"

```python

        def test_feed_map(mocker):
            mocked_markata = mocker.patch.object(markata, "Markata", DummyMarkata())
            feed = Feed(name="archive", config={}, posts=[], _m=mocked_markata)
            feed.map()
```


!! function <h2 id='test_feeds' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>test_feeds <em class='small'>function</em></h2>

???+ source "test_feeds <em class='small'>source</em>"

```python

        def test_feeds(mocker):
            mocked_markata = mocker.patch.object(markata, "Markata", DummyMarkata())
            Feeds(markata=mocked_markata)
```


!! function <h2 id='test_feeds_iter' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>test_feeds_iter <em class='small'>function</em></h2>

???+ source "test_feeds_iter <em class='small'>source</em>"

```python

        def test_feeds_iter(mocker):
            mocked_markata = mocker.patch.object(markata, "Markata", DummyMarkata())
            feeds = Feeds(markata=mocked_markata)

            assert [feed for feed in feeds] == ["all-posts"]
```


!! function <h2 id='test_feeds_keys' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>test_feeds_keys <em class='small'>function</em></h2>

???+ source "test_feeds_keys <em class='small'>source</em>"

```python

        def test_feeds_keys(mocker):
            mocked_markata = mocker.patch.object(markata, "Markata", DummyMarkata())
            feeds = Feeds(markata=mocked_markata)

            assert [feed for feed in feeds.keys()] == ["all-posts"]
```


!! function <h2 id='test_feeds_items' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>test_feeds_items <em class='small'>function</em></h2>

???+ source "test_feeds_items <em class='small'>source</em>"

```python

        def test_feeds_items(mocker):
            mocked_markata = mocker.patch.object(markata, "Markata", DummyMarkata())
            feeds = Feeds(markata=mocked_markata)
            iterate_feeds = {name: feed for name, feed in feeds.items()}

            assert list(iterate_feeds.keys()) == ["all-posts"]
```


!! function <h2 id='test_feeds___get_item__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>test_feeds___get_item__ <em class='small'>function</em></h2>

???+ source "test_feeds___get_item__ <em class='small'>source</em>"

```python

        def test_feeds___get_item__(mocker):
            mocked_markata = mocker.patch.object(markata, "Markata", DummyMarkata())
            feeds = Feeds(markata=mocked_markata)

            feeds["all-posts"].name == "all_posts"
```


!! function <h2 id='test_feeds___rich__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>test_feeds___rich__ <em class='small'>function</em></h2>

???+ source "test_feeds___rich__ <em class='small'>source</em>"

```python

        def test_feeds___rich__(mocker):
            mocked_markata = mocker.patch.object(markata, "Markata", DummyMarkata())
            feeds = Feeds(markata=mocked_markata)

            assert isinstance(feeds.__rich__(), rich.table.Table)
```


!! function <h2 id='test_feeds_values' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>test_feeds_values <em class='small'>function</em></h2>

???+ source "test_feeds_values <em class='small'>source</em>"

```python

        def test_feeds_values(mocker):
            mocked_markata = mocker.patch.object(markata, "Markata", DummyMarkata())
            values = Feeds(markata=mocked_markata).values()
            assert all([isinstance(v, Feed) for v in values])
```


!! method <h2 id='__init__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__init__ <em class='small'>method</em></h2>

???+ source "__init__ <em class='small'>source</em>"

```python

        def __init__(self):
                self.config = {"feeds": {"all-posts": {"filter": True}}}
```


!! method <h2 id='map' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>map <em class='small'>method</em></h2>

???+ source "map <em class='small'>source</em>"

```python

        def map(self, *args, **kwargs):
                return []
```

