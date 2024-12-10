---
content: "The prevnext plugin, creates previous and next links inside each post.\n\n##
  Example config\n\nIn this example we have two maps of posts to look through.  prevnext
  will look\nthrough each of these lists of posts for the current post, then return
  the post\nbefore and after this post as the prevnext posts.\n\n``` toml\n\n[markata]\n#
  default colors will be taken from markata's color_text and color_accent\ncolor_text=white\ncolor_text_light=black\ncolor_accent=white\ncolor_accent_light=black\n\n[markata.prevnext]\n#
  strategy can be 'first' or 'all'\n# 'first' will cycle through the first map the
  post is found in.\n# 'all' will cycle through all of the maps\nstrategy='first'\n\n#
  if you want different colors than your main color_text and color_accent, then\n#
  you can override it here\n# colors can be any valid css color format\n\nprevnext_color_text=white\nprevnext_color_text_light=black\nprevnext_color_angle=white\nprevnext_color_angle_light=black\n\n\n#
  you can have multiple maps, the order they appear will determine their preference\n[markata.feeds.python]\nfilter='\"python\"
  in tags'\nsort='slug'\n\n[markata.feeds.others]\nfilter='\"python\" not in tags'\nsort='slug'\n```\n\nThe
  configuration below will setup two maps, one where \"python\" is in the list\nof
  tags, and another where it is not.  This will link all python posts together\nwith
  a prevnext cycle, and all non-python posts in a separate prevnext cycle.\n\n## strategy\n\nThere
  are currently two supported strategies.\n\n* first\n* all\n\n### first\n\n`first`
  will cycle through only the posts contained within the first map that\ncontains
  the post.\n\n### all\n\n`all` will cycle through all of the posts aggregated from
  any prevnext map.\n\n\n!! class <h2 id='UnsupportedPrevNextStrategy' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>UnsupportedPrevNextStrategy <em class='small'>class</em></h2>\n
  \   A custom error class to raise when an unsupporte prevnext strategy is\n    defined.\n???+
  source \"UnsupportedPrevNextStrategy <em class='small'>source</em>\"\n\n```python\n\n
  \       class UnsupportedPrevNextStrategy(NotImplementedError):\n            \"\"\"\n
  \           A custom error class to raise when an unsupporte prevnext strategy is\n
  \           defined.\n            \"\"\"\n```\n\n\n!! class <h2 id='PrevNext' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>PrevNext <em class='small'>class</em></h2>\n\n???+
  source \"PrevNext <em class='small'>source</em>\"\n\n```python\n\n        class
  PrevNext:\n            prev: str\n            next: str\n```\n\n\n!! function <h2
  id='prevnext' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>prevnext
  <em class='small'>function</em></h2>\n\n???+ source \"prevnext <em class='small'>source</em>\"\n\n```python\n\n
  \       def prevnext(\n            markata: \"Markata\",\n            post: \"Post\",\n
  \           conf: List[Dict[str, str]],\n            strategy: str = \"first\",\n
  \       ) -> Optional[PrevNext]:\n            posts = []\n            for map_conf
  in conf.values():\n                _posts = markata.map(\"post\", **map_conf)\n
  \               # if the strategy is first, cycle back to the beginning after each
  map\n                if strategy == \"first\" and _posts:\n                    _posts.append(_posts[0])\n
  \               posts.extend(_posts)\n            # if the strategy is 'all', cycle
  back to the beginning after all of the maps.\n            if strategy == \"all\":\n
  \               posts.append(posts[0])\n\n            try:\n                post_idx
  = posts.index(post)\n                return PrevNext(prev=posts[post_idx - 1], next=posts[post_idx
  + 1])\n            except ValueError:\n                # post is not in posts\n
  \               return None\n```\n\n\n!! function <h2 id='pre_render' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>pre_render <em class='small'>function</em></h2>\n\n???+
  source \"pre_render <em class='small'>source</em>\"\n\n```python\n\n        def
  pre_render(markata: \"Markata\") -> None:\n            config = markata.config.get(\"prevnext\",
  {})\n            feed_config = markata.config.get(\"feeds\", {})\n            strategy
  = config.get(\"strategy\", \"first\")\n            if strategy not in SUPPORTED_STRATEGIES:\n
  \               msg = f\"\"\"\n                \"{strategy}\" is not a supported
  prevnext strategy\n\n                configure prevnext in your markata.toml to
  use one of {SUPPORTED_STRATEGIES}\n                \"\"\"\n                raise
  UnsupportedPrevNextStrategy(msg)\n            template = config.get(\"template\",
  None)\n            if template is None:\n                template = Template(TEMPLATE)\n
  \           else:\n                template = Template(Path(template).read_text())\n\n
  \           _full_config = copy.deepcopy(markata.config)\n            for article
  in set(markata.articles):\n                article[\"prevnext\"] = prevnext(\n                    markata,\n
  \                   article,\n                    feed_config,\n                    strategy=strategy,\n
  \               )\n                if \"prevnext\" not in article.content and article[\"prevnext\"]:\n
  \                   article.content += template.render(\n                        config=always_merger.merge(\n
  \                           _full_config,\n                            copy.deepcopy(\n
  \                               article.get(\n                                    \"config_overrides\",\n
  \                                   {},\n                                ),\n                            ),\n
  \                       ),\n                        **article,\n                    )\n```\n\n"
date: 0001-01-01
description: The prevnext plugin, creates previous and next links inside each post.
  In this example we have two maps of posts to look through.  prevnext will look The
  config
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Prevnext.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"The prevnext plugin, creates previous and next
    links inside each post. In this example we have two maps of posts to look through.
    \ prevnext will look The config\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Prevnext.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"The prevnext plugin, creates previous
    and next links inside each post. In this example we have two maps of posts to
    look through.  prevnext will look The config\" />\n <link href=\"/favicon.ico\"
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Prevnext.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <p>The prevnext plugin,
    creates previous and next links inside each post.</p>\n<h2 id=\"example-config\">Example
    config <a class=\"header-anchor\" href=\"#example-config\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>In this example we have
    two maps of posts to look through.  prevnext will look\nthrough each of these
    lists of posts for the current post, then return the post\nbefore and after this
    post as the prevnext posts.</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n<span
    class=\"c1\"># default colors will be taken from markata&#39;s color_text and
    color_accent</span>\n<span class=\"n\">color_text</span><span class=\"o\">=</span><span
    class=\"err\">white</span>\n<span class=\"n\">color_text_light</span><span class=\"o\">=</span><span
    class=\"err\">black</span>\n<span class=\"n\">color_accent</span><span class=\"o\">=</span><span
    class=\"err\">white</span>\n<span class=\"n\">color_accent_light</span><span class=\"o\">=</span><span
    class=\"err\">black</span>\n\n<span class=\"k\">[markata.prevnext]</span>\n<span
    class=\"c1\"># strategy can be &#39;first&#39; or &#39;all&#39;</span>\n<span
    class=\"c1\"># &#39;first&#39; will cycle through the first map the post is found
    in.</span>\n<span class=\"c1\"># &#39;all&#39; will cycle through all of the maps</span>\n<span
    class=\"n\">strategy</span><span class=\"o\">=</span><span class=\"s1\">&#39;first&#39;</span>\n\n<span
    class=\"c1\"># if you want different colors than your main color_text and color_accent,
    then</span>\n<span class=\"c1\"># you can override it here</span>\n<span class=\"c1\">#
    colors can be any valid css color format</span>\n\n<span class=\"n\">prevnext_color_text</span><span
    class=\"o\">=</span><span class=\"err\">white</span>\n<span class=\"n\">prevnext_color_text_light</span><span
    class=\"o\">=</span><span class=\"err\">black</span>\n<span class=\"n\">prevnext_color_angle</span><span
    class=\"o\">=</span><span class=\"err\">white</span>\n<span class=\"n\">prevnext_color_angle_light</span><span
    class=\"o\">=</span><span class=\"err\">black</span>\n\n\n<span class=\"c1\">#
    you can have multiple maps, the order they appear will determine their preference</span>\n<span
    class=\"k\">[markata.feeds.python]</span>\n<span class=\"n\">filter</span><span
    class=\"o\">=</span><span class=\"s1\">&#39;&quot;python&quot; in tags&#39;</span>\n<span
    class=\"n\">sort</span><span class=\"o\">=</span><span class=\"s1\">&#39;slug&#39;</span>\n\n<span
    class=\"k\">[markata.feeds.others]</span>\n<span class=\"n\">filter</span><span
    class=\"o\">=</span><span class=\"s1\">&#39;&quot;python&quot; not in tags&#39;</span>\n<span
    class=\"n\">sort</span><span class=\"o\">=</span><span class=\"s1\">&#39;slug&#39;</span>\n</pre></div>\n\n</pre>\n\n<p>The
    configuration below will setup two maps, one where &quot;python&quot; is in the
    list\nof tags, and another where it is not.  This will link all python posts together\nwith
    a prevnext cycle, and all non-python posts in a separate prevnext cycle.</p>\n<h2
    id=\"strategy\">strategy <a class=\"header-anchor\" href=\"#strategy\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>There are currently
    two supported strategies.</p>\n<ul>\n<li>first</li>\n<li>all</li>\n</ul>\n<h3>first</h3>\n<p><code>first</code>
    will cycle through only the posts contained within the first map that\ncontains
    the post.</p>\n<h3>all</h3>\n<p><code>all</code> will cycle through all of the
    posts aggregated from any prevnext map.</p>\n<p>!! class <h2 id='UnsupportedPrevNextStrategy'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>UnsupportedPrevNextStrategy
    <em class='small'>class</em></h2>\nA custom error class to raise when an unsupporte
    prevnext strategy is\ndefined.</p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">UnsupportedPrevNextStrategy
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
    <span class=\"nc\">UnsupportedPrevNextStrategy</span><span class=\"p\">(</span><span
    class=\"ne\">NotImplementedError</span><span class=\"p\">):</span>\n<span class=\"w\">
    \           </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span class=\"sd\">
    \           A custom error class to raise when an unsupporte prevnext strategy
    is</span>\n<span class=\"sd\">            defined.</span>\n<span class=\"sd\">
    \           &quot;&quot;&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!! class <h2
    id='PrevNext' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>PrevNext
    <em class='small'>class</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">PrevNext <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nc\">PrevNext</span><span class=\"p\">:</span>\n            <span
    class=\"n\">prev</span><span class=\"p\">:</span> <span class=\"nb\">str</span>\n
    \           <span class=\"nb\">next</span><span class=\"p\">:</span> <span class=\"nb\">str</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='prevnext' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>prevnext <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">prevnext
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
    <span class=\"nf\">prevnext</span><span class=\"p\">(</span>\n            <span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span>\n            <span class=\"n\">post</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">,</span>\n            <span
    class=\"n\">conf</span><span class=\"p\">:</span> <span class=\"n\">List</span><span
    class=\"p\">[</span><span class=\"n\">Dict</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
    class=\"p\">]],</span>\n            <span class=\"n\">strategy</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;first&quot;</span><span
    class=\"p\">,</span>\n        <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"n\">PrevNext</span><span
    class=\"p\">]:</span>\n            <span class=\"n\">posts</span> <span class=\"o\">=</span>
    <span class=\"p\">[]</span>\n            <span class=\"k\">for</span> <span class=\"n\">map_conf</span>
    <span class=\"ow\">in</span> <span class=\"n\">conf</span><span class=\"o\">.</span><span
    class=\"n\">values</span><span class=\"p\">():</span>\n                <span class=\"n\">_posts</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">map</span><span class=\"p\">(</span><span class=\"s2\">&quot;post&quot;</span><span
    class=\"p\">,</span> <span class=\"o\">**</span><span class=\"n\">map_conf</span><span
    class=\"p\">)</span>\n                <span class=\"c1\"># if the strategy is
    first, cycle back to the beginning after each map</span>\n                <span
    class=\"k\">if</span> <span class=\"n\">strategy</span> <span class=\"o\">==</span>
    <span class=\"s2\">&quot;first&quot;</span> <span class=\"ow\">and</span> <span
    class=\"n\">_posts</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">_posts</span><span class=\"o\">.</span><span class=\"n\">append</span><span
    class=\"p\">(</span><span class=\"n\">_posts</span><span class=\"p\">[</span><span
    class=\"mi\">0</span><span class=\"p\">])</span>\n                <span class=\"n\">posts</span><span
    class=\"o\">.</span><span class=\"n\">extend</span><span class=\"p\">(</span><span
    class=\"n\">_posts</span><span class=\"p\">)</span>\n            <span class=\"c1\">#
    if the strategy is &#39;all&#39;, cycle back to the beginning after all of the
    maps.</span>\n            <span class=\"k\">if</span> <span class=\"n\">strategy</span>
    <span class=\"o\">==</span> <span class=\"s2\">&quot;all&quot;</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">posts</span><span class=\"o\">.</span><span
    class=\"n\">append</span><span class=\"p\">(</span><span class=\"n\">posts</span><span
    class=\"p\">[</span><span class=\"mi\">0</span><span class=\"p\">])</span>\n\n
    \           <span class=\"k\">try</span><span class=\"p\">:</span>\n                <span
    class=\"n\">post_idx</span> <span class=\"o\">=</span> <span class=\"n\">posts</span><span
    class=\"o\">.</span><span class=\"n\">index</span><span class=\"p\">(</span><span
    class=\"n\">post</span><span class=\"p\">)</span>\n                <span class=\"k\">return</span>
    <span class=\"n\">PrevNext</span><span class=\"p\">(</span><span class=\"n\">prev</span><span
    class=\"o\">=</span><span class=\"n\">posts</span><span class=\"p\">[</span><span
    class=\"n\">post_idx</span> <span class=\"o\">-</span> <span class=\"mi\">1</span><span
    class=\"p\">],</span> <span class=\"nb\">next</span><span class=\"o\">=</span><span
    class=\"n\">posts</span><span class=\"p\">[</span><span class=\"n\">post_idx</span>
    <span class=\"o\">+</span> <span class=\"mi\">1</span><span class=\"p\">])</span>\n
    \           <span class=\"k\">except</span> <span class=\"ne\">ValueError</span><span
    class=\"p\">:</span>\n                <span class=\"c1\"># post is not in posts</span>\n
    \               <span class=\"k\">return</span> <span class=\"kc\">None</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='pre_render' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>pre_render <em class='small'>function</em></h2></p>\n<div class=\"admonition
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
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">config</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;prevnext&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">{})</span>\n            <span class=\"n\">feed_config</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;feeds&quot;</span><span class=\"p\">,</span>
    <span class=\"p\">{})</span>\n            <span class=\"n\">strategy</span> <span
    class=\"o\">=</span> <span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;strategy&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;first&quot;</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">if</span> <span class=\"n\">strategy</span> <span
    class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">SUPPORTED_STRATEGIES</span><span
    class=\"p\">:</span>\n                <span class=\"n\">msg</span> <span class=\"o\">=</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;&quot;&quot;</span>\n<span
    class=\"s2\">                &quot;</span><span class=\"si\">{</span><span class=\"n\">strategy</span><span
    class=\"si\">}</span><span class=\"s2\">&quot; is not a supported prevnext strategy</span>\n\n<span
    class=\"s2\">                configure prevnext in your markata.toml to use one
    of </span><span class=\"si\">{</span><span class=\"n\">SUPPORTED_STRATEGIES</span><span
    class=\"si\">}</span>\n<span class=\"s2\">                &quot;&quot;&quot;</span>\n
    \               <span class=\"k\">raise</span> <span class=\"n\">UnsupportedPrevNextStrategy</span><span
    class=\"p\">(</span><span class=\"n\">msg</span><span class=\"p\">)</span>\n            <span
    class=\"n\">template</span> <span class=\"o\">=</span> <span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;template&quot;</span><span class=\"p\">,</span> <span class=\"kc\">None</span><span
    class=\"p\">)</span>\n            <span class=\"k\">if</span> <span class=\"n\">template</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">template</span> <span class=\"o\">=</span> <span
    class=\"n\">Template</span><span class=\"p\">(</span><span class=\"n\">TEMPLATE</span><span
    class=\"p\">)</span>\n            <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">template</span> <span class=\"o\">=</span> <span
    class=\"n\">Template</span><span class=\"p\">(</span><span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">template</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">read_text</span><span class=\"p\">())</span>\n\n
    \           <span class=\"n\">_full_config</span> <span class=\"o\">=</span> <span
    class=\"n\">copy</span><span class=\"o\">.</span><span class=\"n\">deepcopy</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"p\">)</span>\n            <span class=\"k\">for</span>
    <span class=\"n\">article</span> <span class=\"ow\">in</span> <span class=\"nb\">set</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">articles</span><span class=\"p\">):</span>\n                <span
    class=\"n\">article</span><span class=\"p\">[</span><span class=\"s2\">&quot;prevnext&quot;</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"n\">prevnext</span><span
    class=\"p\">(</span>\n                    <span class=\"n\">markata</span><span
    class=\"p\">,</span>\n                    <span class=\"n\">article</span><span
    class=\"p\">,</span>\n                    <span class=\"n\">feed_config</span><span
    class=\"p\">,</span>\n                    <span class=\"n\">strategy</span><span
    class=\"o\">=</span><span class=\"n\">strategy</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">)</span>\n                <span class=\"k\">if</span>
    <span class=\"s2\">&quot;prevnext&quot;</span> <span class=\"ow\">not</span> <span
    class=\"ow\">in</span> <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">content</span> <span class=\"ow\">and</span> <span class=\"n\">article</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;prevnext&quot;</span><span class=\"p\">]:</span>\n
    \                   <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">content</span> <span class=\"o\">+=</span> <span class=\"n\">template</span><span
    class=\"o\">.</span><span class=\"n\">render</span><span class=\"p\">(</span>\n
    \                       <span class=\"n\">config</span><span class=\"o\">=</span><span
    class=\"n\">always_merger</span><span class=\"o\">.</span><span class=\"n\">merge</span><span
    class=\"p\">(</span>\n                            <span class=\"n\">_full_config</span><span
    class=\"p\">,</span>\n                            <span class=\"n\">copy</span><span
    class=\"o\">.</span><span class=\"n\">deepcopy</span><span class=\"p\">(</span>\n
    \                               <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span>\n                                    <span
    class=\"s2\">&quot;config_overrides&quot;</span><span class=\"p\">,</span>\n                                    <span
    class=\"p\">{},</span>\n                                <span class=\"p\">),</span>\n
    \                           <span class=\"p\">),</span>\n                        <span
    class=\"p\">),</span>\n                        <span class=\"o\">**</span><span
    class=\"n\">article</span><span class=\"p\">,</span>\n                    <span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>    </div>\n
    \   <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Prevnext.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"The prevnext plugin, creates previous and next
    links inside each post. In this example we have two maps of posts to look through.
    \ prevnext will look The config\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Prevnext.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"The prevnext plugin, creates previous
    and next links inside each post. In this example we have two maps of posts to
    look through.  prevnext will look The config\" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<article style=\"text-align: center;\">\n    <style>\n        section
    {\n            font-size: 200%;\n        }\n\n\n        .edit {\n            display:
    none;\n        }\n    </style>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Prevnext.Py\n    </h1>\n</section></article>\n     </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Prevnext.Py\n    </h1>\n</section>    <section class=\"body\">\n        <p>The
    prevnext plugin, creates previous and next links inside each post.</p>\n<h2 id=\"example-config\">Example
    config <a class=\"header-anchor\" href=\"#example-config\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>In this example we have
    two maps of posts to look through.  prevnext will look\nthrough each of these
    lists of posts for the current post, then return the post\nbefore and after this
    post as the prevnext posts.</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n<span
    class=\"c1\"># default colors will be taken from markata&#39;s color_text and
    color_accent</span>\n<span class=\"n\">color_text</span><span class=\"o\">=</span><span
    class=\"err\">white</span>\n<span class=\"n\">color_text_light</span><span class=\"o\">=</span><span
    class=\"err\">black</span>\n<span class=\"n\">color_accent</span><span class=\"o\">=</span><span
    class=\"err\">white</span>\n<span class=\"n\">color_accent_light</span><span class=\"o\">=</span><span
    class=\"err\">black</span>\n\n<span class=\"k\">[markata.prevnext]</span>\n<span
    class=\"c1\"># strategy can be &#39;first&#39; or &#39;all&#39;</span>\n<span
    class=\"c1\"># &#39;first&#39; will cycle through the first map the post is found
    in.</span>\n<span class=\"c1\"># &#39;all&#39; will cycle through all of the maps</span>\n<span
    class=\"n\">strategy</span><span class=\"o\">=</span><span class=\"s1\">&#39;first&#39;</span>\n\n<span
    class=\"c1\"># if you want different colors than your main color_text and color_accent,
    then</span>\n<span class=\"c1\"># you can override it here</span>\n<span class=\"c1\">#
    colors can be any valid css color format</span>\n\n<span class=\"n\">prevnext_color_text</span><span
    class=\"o\">=</span><span class=\"err\">white</span>\n<span class=\"n\">prevnext_color_text_light</span><span
    class=\"o\">=</span><span class=\"err\">black</span>\n<span class=\"n\">prevnext_color_angle</span><span
    class=\"o\">=</span><span class=\"err\">white</span>\n<span class=\"n\">prevnext_color_angle_light</span><span
    class=\"o\">=</span><span class=\"err\">black</span>\n\n\n<span class=\"c1\">#
    you can have multiple maps, the order they appear will determine their preference</span>\n<span
    class=\"k\">[markata.feeds.python]</span>\n<span class=\"n\">filter</span><span
    class=\"o\">=</span><span class=\"s1\">&#39;&quot;python&quot; in tags&#39;</span>\n<span
    class=\"n\">sort</span><span class=\"o\">=</span><span class=\"s1\">&#39;slug&#39;</span>\n\n<span
    class=\"k\">[markata.feeds.others]</span>\n<span class=\"n\">filter</span><span
    class=\"o\">=</span><span class=\"s1\">&#39;&quot;python&quot; not in tags&#39;</span>\n<span
    class=\"n\">sort</span><span class=\"o\">=</span><span class=\"s1\">&#39;slug&#39;</span>\n</pre></div>\n\n</pre>\n\n<p>The
    configuration below will setup two maps, one where &quot;python&quot; is in the
    list\nof tags, and another where it is not.  This will link all python posts together\nwith
    a prevnext cycle, and all non-python posts in a separate prevnext cycle.</p>\n<h2
    id=\"strategy\">strategy <a class=\"header-anchor\" href=\"#strategy\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>There are currently
    two supported strategies.</p>\n<ul>\n<li>first</li>\n<li>all</li>\n</ul>\n<h3>first</h3>\n<p><code>first</code>
    will cycle through only the posts contained within the first map that\ncontains
    the post.</p>\n<h3>all</h3>\n<p><code>all</code> will cycle through all of the
    posts aggregated from any prevnext map.</p>\n<p>!! class <h2 id='UnsupportedPrevNextStrategy'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>UnsupportedPrevNextStrategy
    <em class='small'>class</em></h2>\nA custom error class to raise when an unsupporte
    prevnext strategy is\ndefined.</p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">UnsupportedPrevNextStrategy
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
    <span class=\"nc\">UnsupportedPrevNextStrategy</span><span class=\"p\">(</span><span
    class=\"ne\">NotImplementedError</span><span class=\"p\">):</span>\n<span class=\"w\">
    \           </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span class=\"sd\">
    \           A custom error class to raise when an unsupporte prevnext strategy
    is</span>\n<span class=\"sd\">            defined.</span>\n<span class=\"sd\">
    \           &quot;&quot;&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!! class <h2
    id='PrevNext' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>PrevNext
    <em class='small'>class</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">PrevNext <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nc\">PrevNext</span><span class=\"p\">:</span>\n            <span
    class=\"n\">prev</span><span class=\"p\">:</span> <span class=\"nb\">str</span>\n
    \           <span class=\"nb\">next</span><span class=\"p\">:</span> <span class=\"nb\">str</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='prevnext' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>prevnext <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">prevnext
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
    <span class=\"nf\">prevnext</span><span class=\"p\">(</span>\n            <span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span>\n            <span class=\"n\">post</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">,</span>\n            <span
    class=\"n\">conf</span><span class=\"p\">:</span> <span class=\"n\">List</span><span
    class=\"p\">[</span><span class=\"n\">Dict</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
    class=\"p\">]],</span>\n            <span class=\"n\">strategy</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;first&quot;</span><span
    class=\"p\">,</span>\n        <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"n\">PrevNext</span><span
    class=\"p\">]:</span>\n            <span class=\"n\">posts</span> <span class=\"o\">=</span>
    <span class=\"p\">[]</span>\n            <span class=\"k\">for</span> <span class=\"n\">map_conf</span>
    <span class=\"ow\">in</span> <span class=\"n\">conf</span><span class=\"o\">.</span><span
    class=\"n\">values</span><span class=\"p\">():</span>\n                <span class=\"n\">_posts</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">map</span><span class=\"p\">(</span><span class=\"s2\">&quot;post&quot;</span><span
    class=\"p\">,</span> <span class=\"o\">**</span><span class=\"n\">map_conf</span><span
    class=\"p\">)</span>\n                <span class=\"c1\"># if the strategy is
    first, cycle back to the beginning after each map</span>\n                <span
    class=\"k\">if</span> <span class=\"n\">strategy</span> <span class=\"o\">==</span>
    <span class=\"s2\">&quot;first&quot;</span> <span class=\"ow\">and</span> <span
    class=\"n\">_posts</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">_posts</span><span class=\"o\">.</span><span class=\"n\">append</span><span
    class=\"p\">(</span><span class=\"n\">_posts</span><span class=\"p\">[</span><span
    class=\"mi\">0</span><span class=\"p\">])</span>\n                <span class=\"n\">posts</span><span
    class=\"o\">.</span><span class=\"n\">extend</span><span class=\"p\">(</span><span
    class=\"n\">_posts</span><span class=\"p\">)</span>\n            <span class=\"c1\">#
    if the strategy is &#39;all&#39;, cycle back to the beginning after all of the
    maps.</span>\n            <span class=\"k\">if</span> <span class=\"n\">strategy</span>
    <span class=\"o\">==</span> <span class=\"s2\">&quot;all&quot;</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">posts</span><span class=\"o\">.</span><span
    class=\"n\">append</span><span class=\"p\">(</span><span class=\"n\">posts</span><span
    class=\"p\">[</span><span class=\"mi\">0</span><span class=\"p\">])</span>\n\n
    \           <span class=\"k\">try</span><span class=\"p\">:</span>\n                <span
    class=\"n\">post_idx</span> <span class=\"o\">=</span> <span class=\"n\">posts</span><span
    class=\"o\">.</span><span class=\"n\">index</span><span class=\"p\">(</span><span
    class=\"n\">post</span><span class=\"p\">)</span>\n                <span class=\"k\">return</span>
    <span class=\"n\">PrevNext</span><span class=\"p\">(</span><span class=\"n\">prev</span><span
    class=\"o\">=</span><span class=\"n\">posts</span><span class=\"p\">[</span><span
    class=\"n\">post_idx</span> <span class=\"o\">-</span> <span class=\"mi\">1</span><span
    class=\"p\">],</span> <span class=\"nb\">next</span><span class=\"o\">=</span><span
    class=\"n\">posts</span><span class=\"p\">[</span><span class=\"n\">post_idx</span>
    <span class=\"o\">+</span> <span class=\"mi\">1</span><span class=\"p\">])</span>\n
    \           <span class=\"k\">except</span> <span class=\"ne\">ValueError</span><span
    class=\"p\">:</span>\n                <span class=\"c1\"># post is not in posts</span>\n
    \               <span class=\"k\">return</span> <span class=\"kc\">None</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='pre_render' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>pre_render <em class='small'>function</em></h2></p>\n<div class=\"admonition
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
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">config</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;prevnext&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">{})</span>\n            <span class=\"n\">feed_config</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;feeds&quot;</span><span class=\"p\">,</span>
    <span class=\"p\">{})</span>\n            <span class=\"n\">strategy</span> <span
    class=\"o\">=</span> <span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;strategy&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;first&quot;</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">if</span> <span class=\"n\">strategy</span> <span
    class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">SUPPORTED_STRATEGIES</span><span
    class=\"p\">:</span>\n                <span class=\"n\">msg</span> <span class=\"o\">=</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;&quot;&quot;</span>\n<span
    class=\"s2\">                &quot;</span><span class=\"si\">{</span><span class=\"n\">strategy</span><span
    class=\"si\">}</span><span class=\"s2\">&quot; is not a supported prevnext strategy</span>\n\n<span
    class=\"s2\">                configure prevnext in your markata.toml to use one
    of </span><span class=\"si\">{</span><span class=\"n\">SUPPORTED_STRATEGIES</span><span
    class=\"si\">}</span>\n<span class=\"s2\">                &quot;&quot;&quot;</span>\n
    \               <span class=\"k\">raise</span> <span class=\"n\">UnsupportedPrevNextStrategy</span><span
    class=\"p\">(</span><span class=\"n\">msg</span><span class=\"p\">)</span>\n            <span
    class=\"n\">template</span> <span class=\"o\">=</span> <span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;template&quot;</span><span class=\"p\">,</span> <span class=\"kc\">None</span><span
    class=\"p\">)</span>\n            <span class=\"k\">if</span> <span class=\"n\">template</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">template</span> <span class=\"o\">=</span> <span
    class=\"n\">Template</span><span class=\"p\">(</span><span class=\"n\">TEMPLATE</span><span
    class=\"p\">)</span>\n            <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">template</span> <span class=\"o\">=</span> <span
    class=\"n\">Template</span><span class=\"p\">(</span><span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">template</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">read_text</span><span class=\"p\">())</span>\n\n
    \           <span class=\"n\">_full_config</span> <span class=\"o\">=</span> <span
    class=\"n\">copy</span><span class=\"o\">.</span><span class=\"n\">deepcopy</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"p\">)</span>\n            <span class=\"k\">for</span>
    <span class=\"n\">article</span> <span class=\"ow\">in</span> <span class=\"nb\">set</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">articles</span><span class=\"p\">):</span>\n                <span
    class=\"n\">article</span><span class=\"p\">[</span><span class=\"s2\">&quot;prevnext&quot;</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"n\">prevnext</span><span
    class=\"p\">(</span>\n                    <span class=\"n\">markata</span><span
    class=\"p\">,</span>\n                    <span class=\"n\">article</span><span
    class=\"p\">,</span>\n                    <span class=\"n\">feed_config</span><span
    class=\"p\">,</span>\n                    <span class=\"n\">strategy</span><span
    class=\"o\">=</span><span class=\"n\">strategy</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">)</span>\n                <span class=\"k\">if</span>
    <span class=\"s2\">&quot;prevnext&quot;</span> <span class=\"ow\">not</span> <span
    class=\"ow\">in</span> <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">content</span> <span class=\"ow\">and</span> <span class=\"n\">article</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;prevnext&quot;</span><span class=\"p\">]:</span>\n
    \                   <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">content</span> <span class=\"o\">+=</span> <span class=\"n\">template</span><span
    class=\"o\">.</span><span class=\"n\">render</span><span class=\"p\">(</span>\n
    \                       <span class=\"n\">config</span><span class=\"o\">=</span><span
    class=\"n\">always_merger</span><span class=\"o\">.</span><span class=\"n\">merge</span><span
    class=\"p\">(</span>\n                            <span class=\"n\">_full_config</span><span
    class=\"p\">,</span>\n                            <span class=\"n\">copy</span><span
    class=\"o\">.</span><span class=\"n\">deepcopy</span><span class=\"p\">(</span>\n
    \                               <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span>\n                                    <span
    class=\"s2\">&quot;config_overrides&quot;</span><span class=\"p\">,</span>\n                                    <span
    class=\"p\">{},</span>\n                                <span class=\"p\">),</span>\n
    \                           <span class=\"p\">),</span>\n                        <span
    class=\"p\">),</span>\n                        <span class=\"o\">**</span><span
    class=\"n\">article</span><span class=\"p\">,</span>\n                    <span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>"
  raw.md: ''
published: true
slug: markata/plugins/prevnext
title: Prevnext.Py


---

The prevnext plugin, creates previous and next links inside each post.

## Example config

In this example we have two maps of posts to look through.  prevnext will look
through each of these lists of posts for the current post, then return the post
before and after this post as the prevnext posts.

``` toml

[markata]
# default colors will be taken from markata's color_text and color_accent
color_text=white
color_text_light=black
color_accent=white
color_accent_light=black

[markata.prevnext]
# strategy can be 'first' or 'all'
# 'first' will cycle through the first map the post is found in.
# 'all' will cycle through all of the maps
strategy='first'

# if you want different colors than your main color_text and color_accent, then
# you can override it here
# colors can be any valid css color format

prevnext_color_text=white
prevnext_color_text_light=black
prevnext_color_angle=white
prevnext_color_angle_light=black


# you can have multiple maps, the order they appear will determine their preference
[markata.feeds.python]
filter='"python" in tags'
sort='slug'

[markata.feeds.others]
filter='"python" not in tags'
sort='slug'
```

The configuration below will setup two maps, one where "python" is in the list
of tags, and another where it is not.  This will link all python posts together
with a prevnext cycle, and all non-python posts in a separate prevnext cycle.

## strategy

There are currently two supported strategies.

* first
* all

### first

`first` will cycle through only the posts contained within the first map that
contains the post.

### all

`all` will cycle through all of the posts aggregated from any prevnext map.


!! class <h2 id='UnsupportedPrevNextStrategy' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>UnsupportedPrevNextStrategy <em class='small'>class</em></h2>
    A custom error class to raise when an unsupporte prevnext strategy is
    defined.
???+ source "UnsupportedPrevNextStrategy <em class='small'>source</em>"

```python

        class UnsupportedPrevNextStrategy(NotImplementedError):
            """
            A custom error class to raise when an unsupporte prevnext strategy is
            defined.
            """
```


!! class <h2 id='PrevNext' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>PrevNext <em class='small'>class</em></h2>

???+ source "PrevNext <em class='small'>source</em>"

```python

        class PrevNext:
            prev: str
            next: str
```


!! function <h2 id='prevnext' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>prevnext <em class='small'>function</em></h2>

???+ source "prevnext <em class='small'>source</em>"

```python

        def prevnext(
            markata: "Markata",
            post: "Post",
            conf: List[Dict[str, str]],
            strategy: str = "first",
        ) -> Optional[PrevNext]:
            posts = []
            for map_conf in conf.values():
                _posts = markata.map("post", **map_conf)
                # if the strategy is first, cycle back to the beginning after each map
                if strategy == "first" and _posts:
                    _posts.append(_posts[0])
                posts.extend(_posts)
            # if the strategy is 'all', cycle back to the beginning after all of the maps.
            if strategy == "all":
                posts.append(posts[0])

            try:
                post_idx = posts.index(post)
                return PrevNext(prev=posts[post_idx - 1], next=posts[post_idx + 1])
            except ValueError:
                # post is not in posts
                return None
```


!! function <h2 id='pre_render' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>pre_render <em class='small'>function</em></h2>

???+ source "pre_render <em class='small'>source</em>"

```python

        def pre_render(markata: "Markata") -> None:
            config = markata.config.get("prevnext", {})
            feed_config = markata.config.get("feeds", {})
            strategy = config.get("strategy", "first")
            if strategy not in SUPPORTED_STRATEGIES:
                msg = f"""
                "{strategy}" is not a supported prevnext strategy

                configure prevnext in your markata.toml to use one of {SUPPORTED_STRATEGIES}
                """
                raise UnsupportedPrevNextStrategy(msg)
            template = config.get("template", None)
            if template is None:
                template = Template(TEMPLATE)
            else:
                template = Template(Path(template).read_text())

            _full_config = copy.deepcopy(markata.config)
            for article in set(markata.articles):
                article["prevnext"] = prevnext(
                    markata,
                    article,
                    feed_config,
                    strategy=strategy,
                )
                if "prevnext" not in article.content and article["prevnext"]:
                    article.content += template.render(
                        config=always_merger.merge(
                            _full_config,
                            copy.deepcopy(
                                article.get(
                                    "config_overrides",
                                    {},
                                ),
                            ),
                        ),
                        **article,
                    )
```

