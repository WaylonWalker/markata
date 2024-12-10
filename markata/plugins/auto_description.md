---
content: "A Markata plugin to create automatic descriptions for markdown documents.
  \ It\ndoes this by grabbing the first `{len}` number of characters from the document\nthat
  are in a paragraph.\n\n## Configuration\n\nOpen up your `markata.toml` file and
  add new entries for your\nauto_descriptions.  You can have multiple desriptions,
  each one will be named\nafter the key you give it in your config.\n\n``` toml\n[markata]\n\n#
  make sure its in your list of hooks\nhooks=[\n   \"markata.plugins.auto_description\",\n
  \  ]\n\n[markata.auto_description.description]\nlen=160\n[markata.auto_description.long_description]\nlen=250\n[markata.auto_description.super_description]\nlen=500\n```\n\n!!!
  note\n   Make sure that you have the auto_description plugin in your configured
  hooks.\n\nIn the above we will end up with three different descritpions,\n(`description`,
  `long_description`, and `super_description`) each will be the\nfirst number of characters
  from the document as specified in the config.\n\n### Defaults\n\nBy default markata
  will set `description` to 160 and `long_description` to 250,\nif they are not set
  in your config.\n\n### Using the Description\n\nDownstream hooks can now use the
  description for things such as seo, or feeds.\nHere is a simple example that lists
  all of the descriptions in all posts.  This\nis a handy thing you can do right from
  a repl.\n\n``` python\nfrom markata import Markata\nm = Markata()\n[p[\"description\"]
  for p in m.articles]\n```\n\n\n!! function <h2 id='get_description' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>get_description <em class='small'>function</em></h2>\n
  \   Get the full-length description for a single post using the commonmark\n    parser.
  \ Only paragraph nodes will count as text towards the description.\n???+ source
  \"get_description <em class='small'>source</em>\"\n\n```python\n\n        def get_description(article:
  \"Post\") -> str:\n            \"\"\"\n            Get the full-length description
  for a single post using the commonmark\n            parser.  Only paragraph nodes
  will count as text towards the description.\n            \"\"\"\n            ast
  = _parser.parse(article.content)\n\n            # find all paragraph nodes\n            paragraph_nodes
  = [\n                n[0]\n                for n in ast.walker()\n                if
  n[0].t == \"paragraph\" and n[0].first_child.literal is not None\n            ]\n
  \           # for reasons unknown to me commonmark duplicates nodes, dedupe based
  on sourcepos\n            sourcepos = [p.sourcepos for p in paragraph_nodes]\n            #
  find first occurence of node based on source position\n            unique_mask =
  [sourcepos.index(s) == i for i, s in enumerate(sourcepos)]\n            # deduplicate
  paragraph_nodes based on unique source position\n            unique_paragraph_nodes
  = list(compress(paragraph_nodes, unique_mask))\n            paragraphs = \" \".join([p.first_child.literal
  for p in unique_paragraph_nodes])\n            paragraphs = html.escape(paragraphs)\n
  \           return paragraphs\n```\n\n\n!! function <h2 id='set_description' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>set_description <em class='small'>function</em></h2>\n
  \   For a given `article`, find the description, put it in the cache, and set\n
  \   the configured descriptions for the article.\n???+ source \"set_description
  <em class='small'>source</em>\"\n\n```python\n\n        def set_description(\n            markata:
  \"Markata\",\n            article: \"Post\",\n            cache: \"FanoutCache\",\n
  \           config: Dict,\n            max_description: int = 500,\n            plugin_text:
  None = \"\",\n        ) -> None:\n            \"\"\"\n            For a given `article`,
  find the description, put it in the cache, and set\n            the configured descriptions
  for the article.\n            \"\"\"\n            key = markata.make_hash(\n                \"auto_description\",\n
  \               article.content,\n                plugin_text,\n                config,\n
  \           )\n\n            description_from_cache = markata.precache.get(key)\n
  \           if description_from_cache is None:\n                description = get_description(article)[:max_description]\n
  \               markata.cache.add(key, description, expire=markata.config.default_cache_expire)\n
  \           else:\n                description = description_from_cache\n\n            for
  description_key in config:\n                if description_key not in [\"cache_expire\",
  \"config_key\"]:\n                    # overwrites missing (None) and empty ('')\n
  \                   if not article.metadata.get(description_key):\n                        article.metadata[description_key]
  = description[\n                            : config[description_key][\"len\"]\n
  \                       ]\n```\n\n\n!! function <h2 id='pre_render' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>pre_render <em class='small'>function</em></h2>\n
  \   The Markata hook that will set descriptions for all posts in the pre-render
  phase.\n???+ source \"pre_render <em class='small'>source</em>\"\n\n```python\n\n
  \       def pre_render(markata: \"Markata\") -> None:\n            \"\"\"\n            The
  Markata hook that will set descriptions for all posts in the pre-render phase.\n
  \           \"\"\"\n            config = markata.get_plugin_config(__file__)\n\n
  \           if \"description\" not in config.keys():\n                config[\"description\"]
  = {}\n                config[\"description\"][\"len\"] = 160\n\n            if \"long_description\"
  not in config.keys():\n                config[\"long_description\"] = {}\n                config[\"long_description\"][\"len\"]
  = 250\n\n            def try_config_get(key: str) -> Any:\n                try:\n
  \                   return config.get(key).get(\"len\") or None\n                except
  AttributeError:\n                    return None\n\n            max_description
  = max(\n                [\n                    value\n                    for description_key
  in config\n                    if (value := try_config_get(description_key))\n                ],\n
  \           )\n\n            with markata.cache as cache:\n                for article
  in markata.iter_articles(\"setting auto description\"):\n                    set_description(\n
  \                       markata=markata,\n                        article=article,\n
  \                       cache=cache,\n                        config=config,\n                        max_description=max_description,\n
  \                       plugin_text=Path(__file__).read_text(),\n                    )\n```\n\n\n!!
  function <h2 id='try_config_get' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>try_config_get <em class='small'>function</em></h2>\n\n???+ source \"try_config_get
  <em class='small'>source</em>\"\n\n```python\n\n        def try_config_get(key:
  str) -> Any:\n                try:\n                    return config.get(key).get(\"len\")
  or None\n                except AttributeError:\n                    return None\n```\n\n"
date: 0001-01-01
description: 'A Markata plugin to create automatic descriptions for markdown documents.  It
  Open up your  ! In the above we will end up with three different descritpions, By '
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Auto_Description.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"A Markata plugin to create automatic
    descriptions for markdown documents.  It Open up your  ! In the above we will
    end up with three different descritpions, By \" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Auto_Description.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"A Markata plugin to create automatic
    descriptions for markdown documents.  It Open up your  ! In the above we will
    end up with three different descritpions, By \" />\n <link href=\"/favicon.ico\"
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Auto_Description.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <p>A Markata plugin
    to create automatic descriptions for markdown documents.  It\ndoes this by grabbing
    the first <code>{len}</code> number of characters from the document\nthat are
    in a paragraph.</p>\n<h2 id=\"configuration\">Configuration <a class=\"header-anchor\"
    href=\"#configuration\"><svg class=\"heading-permalink\" aria-hidden=\"true\"
    fill=\"currentColor\" focusable=\"false\" height=\"1em\" viewBox=\"0 0 24 24\"
    width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M9.199 13.599a5.99
    5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992 5.992 0
    0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242 6.003
    6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
    3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
    3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
    13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
    2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0
    0 0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
    19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997
    0 0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Open up your <code>markata.toml</code>
    file and add new entries for your\nauto_descriptions.  You can have multiple desriptions,
    each one will be named\nafter the key you give it in your config.</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n\n<span
    class=\"c1\"># make sure its in your list of hooks</span>\n<span class=\"n\">hooks</span><span
    class=\"o\">=</span><span class=\"p\">[</span>\n<span class=\"w\">   </span><span
    class=\"s2\">&quot;markata.plugins.auto_description&quot;</span><span class=\"p\">,</span>\n<span
    class=\"w\">   </span><span class=\"p\">]</span>\n\n<span class=\"k\">[markata.auto_description.description]</span>\n<span
    class=\"n\">len</span><span class=\"o\">=</span><span class=\"mi\">160</span>\n<span
    class=\"k\">[markata.auto_description.long_description]</span>\n<span class=\"n\">len</span><span
    class=\"o\">=</span><span class=\"mi\">250</span>\n<span class=\"k\">[markata.auto_description.super_description]</span>\n<span
    class=\"n\">len</span><span class=\"o\">=</span><span class=\"mi\">500</span>\n</pre></div>\n\n</pre>\n\n<div
    class=\"admonition note\">\n<p class=\"admonition-title\">Note</p>\n</div>\n<p>Make
    sure that you have the auto_description plugin in your configured hooks.</p>\n<p>In
    the above we will end up with three different descritpions,\n(<code>description</code>,
    <code>long_description</code>, and <code>super_description</code>) each will be
    the\nfirst number of characters from the document as specified in the config.</p>\n<h3>Defaults</h3>\n<p>By
    default markata will set <code>description</code> to 160 and <code>long_description</code>
    to 250,\nif they are not set in your config.</p>\n<h3>Using the Description</h3>\n<p>Downstream
    hooks can now use the description for things such as seo, or feeds.\nHere is a
    simple example that lists all of the descriptions in all posts.  This\nis a handy
    thing you can do right from a repl.</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    <span class=\"nn\">markata</span> <span class=\"kn\">import</span> <span class=\"n\">Markata</span>\n<span
    class=\"n\">m</span> <span class=\"o\">=</span> <span class=\"n\">Markata</span><span
    class=\"p\">()</span>\n<span class=\"p\">[</span><span class=\"n\">p</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;description&quot;</span><span class=\"p\">]</span>
    <span class=\"k\">for</span> <span class=\"n\">p</span> <span class=\"ow\">in</span>
    <span class=\"n\">m</span><span class=\"o\">.</span><span class=\"n\">articles</span><span
    class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='get_description'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get_description
    <em class='small'>function</em></h2>\nGet the full-length description for a single
    post using the commonmark\nparser.  Only paragraph nodes will count as text towards
    the description.</p>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">get_description <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">get_description</span><span class=\"p\">(</span><span class=\"n\">article</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            Get the full-length description for a single post using
    the commonmark</span>\n<span class=\"sd\">            parser.  Only paragraph
    nodes will count as text towards the description.</span>\n<span class=\"sd\">
    \           &quot;&quot;&quot;</span>\n            <span class=\"n\">ast</span>
    <span class=\"o\">=</span> <span class=\"n\">_parser</span><span class=\"o\">.</span><span
    class=\"n\">parse</span><span class=\"p\">(</span><span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">content</span><span class=\"p\">)</span>\n\n
    \           <span class=\"c1\"># find all paragraph nodes</span>\n            <span
    class=\"n\">paragraph_nodes</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
    \               <span class=\"n\">n</span><span class=\"p\">[</span><span class=\"mi\">0</span><span
    class=\"p\">]</span>\n                <span class=\"k\">for</span> <span class=\"n\">n</span>
    <span class=\"ow\">in</span> <span class=\"n\">ast</span><span class=\"o\">.</span><span
    class=\"n\">walker</span><span class=\"p\">()</span>\n                <span class=\"k\">if</span>
    <span class=\"n\">n</span><span class=\"p\">[</span><span class=\"mi\">0</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">t</span> <span
    class=\"o\">==</span> <span class=\"s2\">&quot;paragraph&quot;</span> <span class=\"ow\">and</span>
    <span class=\"n\">n</span><span class=\"p\">[</span><span class=\"mi\">0</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">first_child</span><span
    class=\"o\">.</span><span class=\"n\">literal</span> <span class=\"ow\">is</span>
    <span class=\"ow\">not</span> <span class=\"kc\">None</span>\n            <span
    class=\"p\">]</span>\n            <span class=\"c1\"># for reasons unknown to
    me commonmark duplicates nodes, dedupe based on sourcepos</span>\n            <span
    class=\"n\">sourcepos</span> <span class=\"o\">=</span> <span class=\"p\">[</span><span
    class=\"n\">p</span><span class=\"o\">.</span><span class=\"n\">sourcepos</span>
    <span class=\"k\">for</span> <span class=\"n\">p</span> <span class=\"ow\">in</span>
    <span class=\"n\">paragraph_nodes</span><span class=\"p\">]</span>\n            <span
    class=\"c1\"># find first occurence of node based on source position</span>\n
    \           <span class=\"n\">unique_mask</span> <span class=\"o\">=</span> <span
    class=\"p\">[</span><span class=\"n\">sourcepos</span><span class=\"o\">.</span><span
    class=\"n\">index</span><span class=\"p\">(</span><span class=\"n\">s</span><span
    class=\"p\">)</span> <span class=\"o\">==</span> <span class=\"n\">i</span> <span
    class=\"k\">for</span> <span class=\"n\">i</span><span class=\"p\">,</span> <span
    class=\"n\">s</span> <span class=\"ow\">in</span> <span class=\"nb\">enumerate</span><span
    class=\"p\">(</span><span class=\"n\">sourcepos</span><span class=\"p\">)]</span>\n
    \           <span class=\"c1\"># deduplicate paragraph_nodes based on unique source
    position</span>\n            <span class=\"n\">unique_paragraph_nodes</span> <span
    class=\"o\">=</span> <span class=\"nb\">list</span><span class=\"p\">(</span><span
    class=\"n\">compress</span><span class=\"p\">(</span><span class=\"n\">paragraph_nodes</span><span
    class=\"p\">,</span> <span class=\"n\">unique_mask</span><span class=\"p\">))</span>\n
    \           <span class=\"n\">paragraphs</span> <span class=\"o\">=</span> <span
    class=\"s2\">&quot; &quot;</span><span class=\"o\">.</span><span class=\"n\">join</span><span
    class=\"p\">([</span><span class=\"n\">p</span><span class=\"o\">.</span><span
    class=\"n\">first_child</span><span class=\"o\">.</span><span class=\"n\">literal</span>
    <span class=\"k\">for</span> <span class=\"n\">p</span> <span class=\"ow\">in</span>
    <span class=\"n\">unique_paragraph_nodes</span><span class=\"p\">])</span>\n            <span
    class=\"n\">paragraphs</span> <span class=\"o\">=</span> <span class=\"n\">html</span><span
    class=\"o\">.</span><span class=\"n\">escape</span><span class=\"p\">(</span><span
    class=\"n\">paragraphs</span><span class=\"p\">)</span>\n            <span class=\"k\">return</span>
    <span class=\"n\">paragraphs</span>\n</pre></div>\n\n</pre>\n\n<p>!! function
    <h2 id='set_description' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>set_description <em class='small'>function</em></h2>\nFor a given <code>article</code>,
    find the description, put it in the cache, and set\nthe configured descriptions
    for the article.</p>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">set_description <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">set_description</span><span class=\"p\">(</span>\n            <span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span>\n            <span class=\"n\">article</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">,</span>\n            <span
    class=\"n\">cache</span><span class=\"p\">:</span> <span class=\"s2\">&quot;FanoutCache&quot;</span><span
    class=\"p\">,</span>\n            <span class=\"n\">config</span><span class=\"p\">:</span>
    <span class=\"n\">Dict</span><span class=\"p\">,</span>\n            <span class=\"n\">max_description</span><span
    class=\"p\">:</span> <span class=\"nb\">int</span> <span class=\"o\">=</span>
    <span class=\"mi\">500</span><span class=\"p\">,</span>\n            <span class=\"n\">plugin_text</span><span
    class=\"p\">:</span> <span class=\"kc\">None</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;&quot;</span><span class=\"p\">,</span>\n        <span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n<span class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            For a given `article`, find the description, put it in
    the cache, and set</span>\n<span class=\"sd\">            the configured descriptions
    for the article.</span>\n<span class=\"sd\">            &quot;&quot;&quot;</span>\n
    \           <span class=\"n\">key</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">make_hash</span><span class=\"p\">(</span>\n
    \               <span class=\"s2\">&quot;auto_description&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">content</span><span class=\"p\">,</span>\n                <span class=\"n\">plugin_text</span><span
    class=\"p\">,</span>\n                <span class=\"n\">config</span><span class=\"p\">,</span>\n
    \           <span class=\"p\">)</span>\n\n            <span class=\"n\">description_from_cache</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">precache</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"n\">key</span><span class=\"p\">)</span>\n            <span
    class=\"k\">if</span> <span class=\"n\">description_from_cache</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
    class=\"n\">description</span> <span class=\"o\">=</span> <span class=\"n\">get_description</span><span
    class=\"p\">(</span><span class=\"n\">article</span><span class=\"p\">)[:</span><span
    class=\"n\">max_description</span><span class=\"p\">]</span>\n                <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">cache</span><span
    class=\"o\">.</span><span class=\"n\">add</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">,</span> <span class=\"n\">description</span><span
    class=\"p\">,</span> <span class=\"n\">expire</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">default_cache_expire</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">else</span><span class=\"p\">:</span>\n                <span
    class=\"n\">description</span> <span class=\"o\">=</span> <span class=\"n\">description_from_cache</span>\n\n
    \           <span class=\"k\">for</span> <span class=\"n\">description_key</span>
    <span class=\"ow\">in</span> <span class=\"n\">config</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">description_key</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"p\">[</span><span
    class=\"s2\">&quot;cache_expire&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;config_key&quot;</span><span
    class=\"p\">]:</span>\n                    <span class=\"c1\"># overwrites missing
    (None) and empty (&#39;&#39;)</span>\n                    <span class=\"k\">if</span>
    <span class=\"ow\">not</span> <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">metadata</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"n\">description_key</span><span class=\"p\">):</span>\n
    \                       <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">metadata</span><span class=\"p\">[</span><span class=\"n\">description_key</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"n\">description</span><span
    class=\"p\">[</span>\n                            <span class=\"p\">:</span> <span
    class=\"n\">config</span><span class=\"p\">[</span><span class=\"n\">description_key</span><span
    class=\"p\">][</span><span class=\"s2\">&quot;len&quot;</span><span class=\"p\">]</span>\n
    \                       <span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='pre_render' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>pre_render <em class='small'>function</em></h2>\nThe Markata hook that
    will set descriptions for all posts in the pre-render phase.</p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">pre_render
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
    <span class=\"nf\">pre_render</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            The Markata hook that will set descriptions for all posts
    in the pre-render phase.</span>\n<span class=\"sd\">            &quot;&quot;&quot;</span>\n
    \           <span class=\"n\">config</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">get_plugin_config</span><span class=\"p\">(</span><span
    class=\"vm\">__file__</span><span class=\"p\">)</span>\n\n            <span class=\"k\">if</span>
    <span class=\"s2\">&quot;description&quot;</span> <span class=\"ow\">not</span>
    <span class=\"ow\">in</span> <span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">keys</span><span class=\"p\">():</span>\n                <span class=\"n\">config</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;description&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"p\">{}</span>\n                <span
    class=\"n\">config</span><span class=\"p\">[</span><span class=\"s2\">&quot;description&quot;</span><span
    class=\"p\">][</span><span class=\"s2\">&quot;len&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"mi\">160</span>\n\n            <span
    class=\"k\">if</span> <span class=\"s2\">&quot;long_description&quot;</span> <span
    class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">keys</span><span class=\"p\">():</span>\n
    \               <span class=\"n\">config</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;long_description&quot;</span><span class=\"p\">]</span> <span
    class=\"o\">=</span> <span class=\"p\">{}</span>\n                <span class=\"n\">config</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;long_description&quot;</span><span
    class=\"p\">][</span><span class=\"s2\">&quot;len&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"mi\">250</span>\n\n            <span
    class=\"k\">def</span> <span class=\"nf\">try_config_get</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span
    class=\"p\">:</span>\n                <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;len&quot;</span><span
    class=\"p\">)</span> <span class=\"ow\">or</span> <span class=\"kc\">None</span>\n
    \               <span class=\"k\">except</span> <span class=\"ne\">AttributeError</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
    class=\"kc\">None</span>\n\n            <span class=\"n\">max_description</span>
    <span class=\"o\">=</span> <span class=\"nb\">max</span><span class=\"p\">(</span>\n
    \               <span class=\"p\">[</span>\n                    <span class=\"n\">value</span>\n
    \                   <span class=\"k\">for</span> <span class=\"n\">description_key</span>
    <span class=\"ow\">in</span> <span class=\"n\">config</span>\n                    <span
    class=\"k\">if</span> <span class=\"p\">(</span><span class=\"n\">value</span>
    <span class=\"o\">:=</span> <span class=\"n\">try_config_get</span><span class=\"p\">(</span><span
    class=\"n\">description_key</span><span class=\"p\">))</span>\n                <span
    class=\"p\">],</span>\n            <span class=\"p\">)</span>\n\n            <span
    class=\"k\">with</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">cache</span> <span class=\"k\">as</span> <span class=\"n\">cache</span><span
    class=\"p\">:</span>\n                <span class=\"k\">for</span> <span class=\"n\">article</span>
    <span class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">iter_articles</span><span class=\"p\">(</span><span class=\"s2\">&quot;setting
    auto description&quot;</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">set_description</span><span class=\"p\">(</span>\n                        <span
    class=\"n\">markata</span><span class=\"o\">=</span><span class=\"n\">markata</span><span
    class=\"p\">,</span>\n                        <span class=\"n\">article</span><span
    class=\"o\">=</span><span class=\"n\">article</span><span class=\"p\">,</span>\n
    \                       <span class=\"n\">cache</span><span class=\"o\">=</span><span
    class=\"n\">cache</span><span class=\"p\">,</span>\n                        <span
    class=\"n\">config</span><span class=\"o\">=</span><span class=\"n\">config</span><span
    class=\"p\">,</span>\n                        <span class=\"n\">max_description</span><span
    class=\"o\">=</span><span class=\"n\">max_description</span><span class=\"p\">,</span>\n
    \                       <span class=\"n\">plugin_text</span><span class=\"o\">=</span><span
    class=\"n\">Path</span><span class=\"p\">(</span><span class=\"vm\">__file__</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">read_text</span><span
    class=\"p\">(),</span>\n                    <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='try_config_get' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>try_config_get <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">try_config_get
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
    <span class=\"nf\">try_config_get</span><span class=\"p\">(</span><span class=\"n\">key</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">try</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"n\">key</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;len&quot;</span><span class=\"p\">)</span>
    <span class=\"ow\">or</span> <span class=\"kc\">None</span>\n                <span
    class=\"k\">except</span> <span class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"kc\">None</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>    </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Auto_Description.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"A Markata plugin to create automatic
    descriptions for markdown documents.  It Open up your  ! In the above we will
    end up with three different descritpions, By \" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Auto_Description.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"A Markata plugin to create automatic
    descriptions for markdown documents.  It Open up your  ! In the above we will
    end up with three different descritpions, By \" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<article style=\"text-align: center;\">\n    <style>\n        section
    {\n            font-size: 200%;\n        }\n\n\n        .edit {\n            display:
    none;\n        }\n    </style>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Auto_Description.Py\n    </h1>\n</section></article>\n     </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Auto_Description.Py\n    </h1>\n</section>    <section class=\"body\">\n
    \       <p>A Markata plugin to create automatic descriptions for markdown documents.
    \ It\ndoes this by grabbing the first <code>{len}</code> number of characters
    from the document\nthat are in a paragraph.</p>\n<h2 id=\"configuration\">Configuration
    <a class=\"header-anchor\" href=\"#configuration\"><svg class=\"heading-permalink\"
    aria-hidden=\"true\" fill=\"currentColor\" focusable=\"false\" height=\"1em\"
    viewBox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
    d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
    5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
    6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
    3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
    3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
    13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
    2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0
    0 0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
    19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997
    0 0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Open up your <code>markata.toml</code>
    file and add new entries for your\nauto_descriptions.  You can have multiple desriptions,
    each one will be named\nafter the key you give it in your config.</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n\n<span
    class=\"c1\"># make sure its in your list of hooks</span>\n<span class=\"n\">hooks</span><span
    class=\"o\">=</span><span class=\"p\">[</span>\n<span class=\"w\">   </span><span
    class=\"s2\">&quot;markata.plugins.auto_description&quot;</span><span class=\"p\">,</span>\n<span
    class=\"w\">   </span><span class=\"p\">]</span>\n\n<span class=\"k\">[markata.auto_description.description]</span>\n<span
    class=\"n\">len</span><span class=\"o\">=</span><span class=\"mi\">160</span>\n<span
    class=\"k\">[markata.auto_description.long_description]</span>\n<span class=\"n\">len</span><span
    class=\"o\">=</span><span class=\"mi\">250</span>\n<span class=\"k\">[markata.auto_description.super_description]</span>\n<span
    class=\"n\">len</span><span class=\"o\">=</span><span class=\"mi\">500</span>\n</pre></div>\n\n</pre>\n\n<div
    class=\"admonition note\">\n<p class=\"admonition-title\">Note</p>\n</div>\n<p>Make
    sure that you have the auto_description plugin in your configured hooks.</p>\n<p>In
    the above we will end up with three different descritpions,\n(<code>description</code>,
    <code>long_description</code>, and <code>super_description</code>) each will be
    the\nfirst number of characters from the document as specified in the config.</p>\n<h3>Defaults</h3>\n<p>By
    default markata will set <code>description</code> to 160 and <code>long_description</code>
    to 250,\nif they are not set in your config.</p>\n<h3>Using the Description</h3>\n<p>Downstream
    hooks can now use the description for things such as seo, or feeds.\nHere is a
    simple example that lists all of the descriptions in all posts.  This\nis a handy
    thing you can do right from a repl.</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    <span class=\"nn\">markata</span> <span class=\"kn\">import</span> <span class=\"n\">Markata</span>\n<span
    class=\"n\">m</span> <span class=\"o\">=</span> <span class=\"n\">Markata</span><span
    class=\"p\">()</span>\n<span class=\"p\">[</span><span class=\"n\">p</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;description&quot;</span><span class=\"p\">]</span>
    <span class=\"k\">for</span> <span class=\"n\">p</span> <span class=\"ow\">in</span>
    <span class=\"n\">m</span><span class=\"o\">.</span><span class=\"n\">articles</span><span
    class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='get_description'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get_description
    <em class='small'>function</em></h2>\nGet the full-length description for a single
    post using the commonmark\nparser.  Only paragraph nodes will count as text towards
    the description.</p>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">get_description <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">get_description</span><span class=\"p\">(</span><span class=\"n\">article</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            Get the full-length description for a single post using
    the commonmark</span>\n<span class=\"sd\">            parser.  Only paragraph
    nodes will count as text towards the description.</span>\n<span class=\"sd\">
    \           &quot;&quot;&quot;</span>\n            <span class=\"n\">ast</span>
    <span class=\"o\">=</span> <span class=\"n\">_parser</span><span class=\"o\">.</span><span
    class=\"n\">parse</span><span class=\"p\">(</span><span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">content</span><span class=\"p\">)</span>\n\n
    \           <span class=\"c1\"># find all paragraph nodes</span>\n            <span
    class=\"n\">paragraph_nodes</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
    \               <span class=\"n\">n</span><span class=\"p\">[</span><span class=\"mi\">0</span><span
    class=\"p\">]</span>\n                <span class=\"k\">for</span> <span class=\"n\">n</span>
    <span class=\"ow\">in</span> <span class=\"n\">ast</span><span class=\"o\">.</span><span
    class=\"n\">walker</span><span class=\"p\">()</span>\n                <span class=\"k\">if</span>
    <span class=\"n\">n</span><span class=\"p\">[</span><span class=\"mi\">0</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">t</span> <span
    class=\"o\">==</span> <span class=\"s2\">&quot;paragraph&quot;</span> <span class=\"ow\">and</span>
    <span class=\"n\">n</span><span class=\"p\">[</span><span class=\"mi\">0</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">first_child</span><span
    class=\"o\">.</span><span class=\"n\">literal</span> <span class=\"ow\">is</span>
    <span class=\"ow\">not</span> <span class=\"kc\">None</span>\n            <span
    class=\"p\">]</span>\n            <span class=\"c1\"># for reasons unknown to
    me commonmark duplicates nodes, dedupe based on sourcepos</span>\n            <span
    class=\"n\">sourcepos</span> <span class=\"o\">=</span> <span class=\"p\">[</span><span
    class=\"n\">p</span><span class=\"o\">.</span><span class=\"n\">sourcepos</span>
    <span class=\"k\">for</span> <span class=\"n\">p</span> <span class=\"ow\">in</span>
    <span class=\"n\">paragraph_nodes</span><span class=\"p\">]</span>\n            <span
    class=\"c1\"># find first occurence of node based on source position</span>\n
    \           <span class=\"n\">unique_mask</span> <span class=\"o\">=</span> <span
    class=\"p\">[</span><span class=\"n\">sourcepos</span><span class=\"o\">.</span><span
    class=\"n\">index</span><span class=\"p\">(</span><span class=\"n\">s</span><span
    class=\"p\">)</span> <span class=\"o\">==</span> <span class=\"n\">i</span> <span
    class=\"k\">for</span> <span class=\"n\">i</span><span class=\"p\">,</span> <span
    class=\"n\">s</span> <span class=\"ow\">in</span> <span class=\"nb\">enumerate</span><span
    class=\"p\">(</span><span class=\"n\">sourcepos</span><span class=\"p\">)]</span>\n
    \           <span class=\"c1\"># deduplicate paragraph_nodes based on unique source
    position</span>\n            <span class=\"n\">unique_paragraph_nodes</span> <span
    class=\"o\">=</span> <span class=\"nb\">list</span><span class=\"p\">(</span><span
    class=\"n\">compress</span><span class=\"p\">(</span><span class=\"n\">paragraph_nodes</span><span
    class=\"p\">,</span> <span class=\"n\">unique_mask</span><span class=\"p\">))</span>\n
    \           <span class=\"n\">paragraphs</span> <span class=\"o\">=</span> <span
    class=\"s2\">&quot; &quot;</span><span class=\"o\">.</span><span class=\"n\">join</span><span
    class=\"p\">([</span><span class=\"n\">p</span><span class=\"o\">.</span><span
    class=\"n\">first_child</span><span class=\"o\">.</span><span class=\"n\">literal</span>
    <span class=\"k\">for</span> <span class=\"n\">p</span> <span class=\"ow\">in</span>
    <span class=\"n\">unique_paragraph_nodes</span><span class=\"p\">])</span>\n            <span
    class=\"n\">paragraphs</span> <span class=\"o\">=</span> <span class=\"n\">html</span><span
    class=\"o\">.</span><span class=\"n\">escape</span><span class=\"p\">(</span><span
    class=\"n\">paragraphs</span><span class=\"p\">)</span>\n            <span class=\"k\">return</span>
    <span class=\"n\">paragraphs</span>\n</pre></div>\n\n</pre>\n\n<p>!! function
    <h2 id='set_description' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>set_description <em class='small'>function</em></h2>\nFor a given <code>article</code>,
    find the description, put it in the cache, and set\nthe configured descriptions
    for the article.</p>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">set_description <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">set_description</span><span class=\"p\">(</span>\n            <span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span>\n            <span class=\"n\">article</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">,</span>\n            <span
    class=\"n\">cache</span><span class=\"p\">:</span> <span class=\"s2\">&quot;FanoutCache&quot;</span><span
    class=\"p\">,</span>\n            <span class=\"n\">config</span><span class=\"p\">:</span>
    <span class=\"n\">Dict</span><span class=\"p\">,</span>\n            <span class=\"n\">max_description</span><span
    class=\"p\">:</span> <span class=\"nb\">int</span> <span class=\"o\">=</span>
    <span class=\"mi\">500</span><span class=\"p\">,</span>\n            <span class=\"n\">plugin_text</span><span
    class=\"p\">:</span> <span class=\"kc\">None</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;&quot;</span><span class=\"p\">,</span>\n        <span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n<span class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            For a given `article`, find the description, put it in
    the cache, and set</span>\n<span class=\"sd\">            the configured descriptions
    for the article.</span>\n<span class=\"sd\">            &quot;&quot;&quot;</span>\n
    \           <span class=\"n\">key</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">make_hash</span><span class=\"p\">(</span>\n
    \               <span class=\"s2\">&quot;auto_description&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">content</span><span class=\"p\">,</span>\n                <span class=\"n\">plugin_text</span><span
    class=\"p\">,</span>\n                <span class=\"n\">config</span><span class=\"p\">,</span>\n
    \           <span class=\"p\">)</span>\n\n            <span class=\"n\">description_from_cache</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">precache</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"n\">key</span><span class=\"p\">)</span>\n            <span
    class=\"k\">if</span> <span class=\"n\">description_from_cache</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
    class=\"n\">description</span> <span class=\"o\">=</span> <span class=\"n\">get_description</span><span
    class=\"p\">(</span><span class=\"n\">article</span><span class=\"p\">)[:</span><span
    class=\"n\">max_description</span><span class=\"p\">]</span>\n                <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">cache</span><span
    class=\"o\">.</span><span class=\"n\">add</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">,</span> <span class=\"n\">description</span><span
    class=\"p\">,</span> <span class=\"n\">expire</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">default_cache_expire</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">else</span><span class=\"p\">:</span>\n                <span
    class=\"n\">description</span> <span class=\"o\">=</span> <span class=\"n\">description_from_cache</span>\n\n
    \           <span class=\"k\">for</span> <span class=\"n\">description_key</span>
    <span class=\"ow\">in</span> <span class=\"n\">config</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">description_key</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"p\">[</span><span
    class=\"s2\">&quot;cache_expire&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;config_key&quot;</span><span
    class=\"p\">]:</span>\n                    <span class=\"c1\"># overwrites missing
    (None) and empty (&#39;&#39;)</span>\n                    <span class=\"k\">if</span>
    <span class=\"ow\">not</span> <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">metadata</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"n\">description_key</span><span class=\"p\">):</span>\n
    \                       <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">metadata</span><span class=\"p\">[</span><span class=\"n\">description_key</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"n\">description</span><span
    class=\"p\">[</span>\n                            <span class=\"p\">:</span> <span
    class=\"n\">config</span><span class=\"p\">[</span><span class=\"n\">description_key</span><span
    class=\"p\">][</span><span class=\"s2\">&quot;len&quot;</span><span class=\"p\">]</span>\n
    \                       <span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='pre_render' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>pre_render <em class='small'>function</em></h2>\nThe Markata hook that
    will set descriptions for all posts in the pre-render phase.</p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">pre_render
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
    <span class=\"nf\">pre_render</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            The Markata hook that will set descriptions for all posts
    in the pre-render phase.</span>\n<span class=\"sd\">            &quot;&quot;&quot;</span>\n
    \           <span class=\"n\">config</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">get_plugin_config</span><span class=\"p\">(</span><span
    class=\"vm\">__file__</span><span class=\"p\">)</span>\n\n            <span class=\"k\">if</span>
    <span class=\"s2\">&quot;description&quot;</span> <span class=\"ow\">not</span>
    <span class=\"ow\">in</span> <span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">keys</span><span class=\"p\">():</span>\n                <span class=\"n\">config</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;description&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"p\">{}</span>\n                <span
    class=\"n\">config</span><span class=\"p\">[</span><span class=\"s2\">&quot;description&quot;</span><span
    class=\"p\">][</span><span class=\"s2\">&quot;len&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"mi\">160</span>\n\n            <span
    class=\"k\">if</span> <span class=\"s2\">&quot;long_description&quot;</span> <span
    class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">keys</span><span class=\"p\">():</span>\n
    \               <span class=\"n\">config</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;long_description&quot;</span><span class=\"p\">]</span> <span
    class=\"o\">=</span> <span class=\"p\">{}</span>\n                <span class=\"n\">config</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;long_description&quot;</span><span
    class=\"p\">][</span><span class=\"s2\">&quot;len&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"mi\">250</span>\n\n            <span
    class=\"k\">def</span> <span class=\"nf\">try_config_get</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span
    class=\"p\">:</span>\n                <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;len&quot;</span><span
    class=\"p\">)</span> <span class=\"ow\">or</span> <span class=\"kc\">None</span>\n
    \               <span class=\"k\">except</span> <span class=\"ne\">AttributeError</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
    class=\"kc\">None</span>\n\n            <span class=\"n\">max_description</span>
    <span class=\"o\">=</span> <span class=\"nb\">max</span><span class=\"p\">(</span>\n
    \               <span class=\"p\">[</span>\n                    <span class=\"n\">value</span>\n
    \                   <span class=\"k\">for</span> <span class=\"n\">description_key</span>
    <span class=\"ow\">in</span> <span class=\"n\">config</span>\n                    <span
    class=\"k\">if</span> <span class=\"p\">(</span><span class=\"n\">value</span>
    <span class=\"o\">:=</span> <span class=\"n\">try_config_get</span><span class=\"p\">(</span><span
    class=\"n\">description_key</span><span class=\"p\">))</span>\n                <span
    class=\"p\">],</span>\n            <span class=\"p\">)</span>\n\n            <span
    class=\"k\">with</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">cache</span> <span class=\"k\">as</span> <span class=\"n\">cache</span><span
    class=\"p\">:</span>\n                <span class=\"k\">for</span> <span class=\"n\">article</span>
    <span class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">iter_articles</span><span class=\"p\">(</span><span class=\"s2\">&quot;setting
    auto description&quot;</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">set_description</span><span class=\"p\">(</span>\n                        <span
    class=\"n\">markata</span><span class=\"o\">=</span><span class=\"n\">markata</span><span
    class=\"p\">,</span>\n                        <span class=\"n\">article</span><span
    class=\"o\">=</span><span class=\"n\">article</span><span class=\"p\">,</span>\n
    \                       <span class=\"n\">cache</span><span class=\"o\">=</span><span
    class=\"n\">cache</span><span class=\"p\">,</span>\n                        <span
    class=\"n\">config</span><span class=\"o\">=</span><span class=\"n\">config</span><span
    class=\"p\">,</span>\n                        <span class=\"n\">max_description</span><span
    class=\"o\">=</span><span class=\"n\">max_description</span><span class=\"p\">,</span>\n
    \                       <span class=\"n\">plugin_text</span><span class=\"o\">=</span><span
    class=\"n\">Path</span><span class=\"p\">(</span><span class=\"vm\">__file__</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">read_text</span><span
    class=\"p\">(),</span>\n                    <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='try_config_get' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>try_config_get <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">try_config_get
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
    <span class=\"nf\">try_config_get</span><span class=\"p\">(</span><span class=\"n\">key</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">try</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"n\">key</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;len&quot;</span><span class=\"p\">)</span>
    <span class=\"ow\">or</span> <span class=\"kc\">None</span>\n                <span
    class=\"k\">except</span> <span class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"kc\">None</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>"
  raw.md: ''
published: true
slug: markata/plugins/auto-description
title: Auto_Description.Py


---

A Markata plugin to create automatic descriptions for markdown documents.  It
does this by grabbing the first `{len}` number of characters from the document
that are in a paragraph.

## Configuration

Open up your `markata.toml` file and add new entries for your
auto_descriptions.  You can have multiple desriptions, each one will be named
after the key you give it in your config.

``` toml
[markata]

# make sure its in your list of hooks
hooks=[
   "markata.plugins.auto_description",
   ]

[markata.auto_description.description]
len=160
[markata.auto_description.long_description]
len=250
[markata.auto_description.super_description]
len=500
```

!!! note
   Make sure that you have the auto_description plugin in your configured hooks.

In the above we will end up with three different descritpions,
(`description`, `long_description`, and `super_description`) each will be the
first number of characters from the document as specified in the config.

### Defaults

By default markata will set `description` to 160 and `long_description` to 250,
if they are not set in your config.

### Using the Description

Downstream hooks can now use the description for things such as seo, or feeds.
Here is a simple example that lists all of the descriptions in all posts.  This
is a handy thing you can do right from a repl.

``` python
from markata import Markata
m = Markata()
[p["description"] for p in m.articles]
```


!! function <h2 id='get_description' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get_description <em class='small'>function</em></h2>
    Get the full-length description for a single post using the commonmark
    parser.  Only paragraph nodes will count as text towards the description.
???+ source "get_description <em class='small'>source</em>"

```python

        def get_description(article: "Post") -> str:
            """
            Get the full-length description for a single post using the commonmark
            parser.  Only paragraph nodes will count as text towards the description.
            """
            ast = _parser.parse(article.content)

            # find all paragraph nodes
            paragraph_nodes = [
                n[0]
                for n in ast.walker()
                if n[0].t == "paragraph" and n[0].first_child.literal is not None
            ]
            # for reasons unknown to me commonmark duplicates nodes, dedupe based on sourcepos
            sourcepos = [p.sourcepos for p in paragraph_nodes]
            # find first occurence of node based on source position
            unique_mask = [sourcepos.index(s) == i for i, s in enumerate(sourcepos)]
            # deduplicate paragraph_nodes based on unique source position
            unique_paragraph_nodes = list(compress(paragraph_nodes, unique_mask))
            paragraphs = " ".join([p.first_child.literal for p in unique_paragraph_nodes])
            paragraphs = html.escape(paragraphs)
            return paragraphs
```


!! function <h2 id='set_description' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>set_description <em class='small'>function</em></h2>
    For a given `article`, find the description, put it in the cache, and set
    the configured descriptions for the article.
???+ source "set_description <em class='small'>source</em>"

```python

        def set_description(
            markata: "Markata",
            article: "Post",
            cache: "FanoutCache",
            config: Dict,
            max_description: int = 500,
            plugin_text: None = "",
        ) -> None:
            """
            For a given `article`, find the description, put it in the cache, and set
            the configured descriptions for the article.
            """
            key = markata.make_hash(
                "auto_description",
                article.content,
                plugin_text,
                config,
            )

            description_from_cache = markata.precache.get(key)
            if description_from_cache is None:
                description = get_description(article)[:max_description]
                markata.cache.add(key, description, expire=markata.config.default_cache_expire)
            else:
                description = description_from_cache

            for description_key in config:
                if description_key not in ["cache_expire", "config_key"]:
                    # overwrites missing (None) and empty ('')
                    if not article.metadata.get(description_key):
                        article.metadata[description_key] = description[
                            : config[description_key]["len"]
                        ]
```


!! function <h2 id='pre_render' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>pre_render <em class='small'>function</em></h2>
    The Markata hook that will set descriptions for all posts in the pre-render phase.
???+ source "pre_render <em class='small'>source</em>"

```python

        def pre_render(markata: "Markata") -> None:
            """
            The Markata hook that will set descriptions for all posts in the pre-render phase.
            """
            config = markata.get_plugin_config(__file__)

            if "description" not in config.keys():
                config["description"] = {}
                config["description"]["len"] = 160

            if "long_description" not in config.keys():
                config["long_description"] = {}
                config["long_description"]["len"] = 250

            def try_config_get(key: str) -> Any:
                try:
                    return config.get(key).get("len") or None
                except AttributeError:
                    return None

            max_description = max(
                [
                    value
                    for description_key in config
                    if (value := try_config_get(description_key))
                ],
            )

            with markata.cache as cache:
                for article in markata.iter_articles("setting auto description"):
                    set_description(
                        markata=markata,
                        article=article,
                        cache=cache,
                        config=config,
                        max_description=max_description,
                        plugin_text=Path(__file__).read_text(),
                    )
```


!! function <h2 id='try_config_get' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>try_config_get <em class='small'>function</em></h2>

???+ source "try_config_get <em class='small'>source</em>"

```python

        def try_config_get(key: str) -> Any:
                try:
                    return config.get(key).get("len") or None
                except AttributeError:
                    return None
```

