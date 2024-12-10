---
content: "Default load plugin.\n\n\n!! class <h2 id='ValidationError' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>ValidationError <em class='small'>class</em></h2>\n\n???+
  source \"ValidationError <em class='small'>source</em>\"\n\n```python\n\n        class
  ValidationError(ValueError): ...\n```\n\n\n!! function <h2 id='load' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>load <em class='small'>function</em></h2>\n\n???+
  source \"load <em class='small'>source</em>\"\n\n```python\n\n        def load(markata:
  \"MarkataMarkdown\") -> None:\n            Progress(\n                BarColumn(bar_width=None),\n
  \               transient=True,\n                console=markata.console,\n            )\n
  \           markata.console.log(f\"found {len(markata.files)} posts\")\n            post_futures
  = [get_post(article, markata) for article in markata.files]\n            posts =
  [post.result() for post in post_futures if post is not None]\n\n            markata.posts_obj
  = markata.Posts.parse_obj(\n                {\"posts\": posts},\n            )\n
  \           markata.posts = markata.posts_obj.posts\n            markata.articles
  = markata.posts\n```\n\n\n!! function <h2 id='get_post' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>get_post <em class='small'>function</em></h2>\n\n???+
  source \"get_post <em class='small'>source</em>\"\n\n```python\n\n        def get_post(path:
  Path, markata: \"Markata\") -> Optional[Callable]:\n            if markata.Post:\n
  \               post = pydantic_get_post(path=path, markata=markata)\n                return
  post\n            else:\n                return legacy_get_post(path=path, markata=markata)\n```\n\n\n!!
  function <h2 id='get_models' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>get_models <em class='small'>function</em></h2>\n\n???+ source \"get_models
  <em class='small'>source</em>\"\n\n```python\n\n        def get_models(markata:
  \"Markata\", error: pydantic.ValidationError) -> List:\n            fields = []\n
  \           for err in error.errors():\n                fields.extend(err[\"loc\"])\n\n
  \           models = {field: f\"{field} used by \" for field in fields}\n\n            for
  field, model in set(\n                itertools.product(\n                    fields,\n
  \                   markata.post_models,\n                ),\n            ):\n                if
  field in model.__fields__:\n                    models[field] += f\"'{model.__module__}.{model.__name__}'\"\n\n
  \           return models\n```\n\n\n!! function <h2 id='pydantic_get_post' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>pydantic_get_post <em class='small'>function</em></h2>\n\n???+
  source \"pydantic_get_post <em class='small'>source</em>\"\n\n```python\n\n        def
  pydantic_get_post(path: Path, markata: \"Markata\") -> Optional[Callable]:\n            try:\n
  \               post = markata.Post.parse_file(markata=markata, path=path)\n                markata.Post.validate(post)\n\n
  \           except pydantic.ValidationError as e:\n                models = get_models(markata=markata,
  error=e)\n                models = list(models.values())\n                models
  = \"\\n\".join(models)\n                raise ValidationError(f\"{e}\\n\\n{models}\\nfailed
  to load {path}\") from e\n\n            return post\n```\n\n\n!! function <h2 id='legacy_get_post'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>legacy_get_post <em
  class='small'>function</em></h2>\n\n???+ source \"legacy_get_post <em class='small'>source</em>\"\n\n```python\n\n
  \       def legacy_get_post(path: Path, markata: \"Markata\") -> Optional[Callable]:\n
  \           default = {\n                \"cover\": \"\",\n                \"title\":
  \"\",\n                \"tags\": [],\n                \"published\": \"False\",\n
  \               \"templateKey\": \"\",\n                \"path\": str(path),\n                \"description\":
  \"\",\n                \"content\": \"\",\n            }\n            try:\n                post:
  \"Post\" = frontmatter.load(path)\n                post.metadata = {**default, **post.metadata}\n
  \               post[\"content\"] = post.content\n            except ParserError:\n
  \               return None\n                post = default\n            except
  ValueError:\n                return None\n                post = default\n            post.metadata[\"path\"]
  = str(path)\n            post[\"edit_link\"] = (\n                markata.config.repo_url
  + \"edit/\" + markata.config.repo_branch + \"/\" + post.path\n            )\n            return
  post\n```\n\n\n!! class <h2 id='MarkataMarkdown' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>MarkataMarkdown <em class='small'>class</em></h2>\n\n???+ source \"MarkataMarkdown
  <em class='small'>source</em>\"\n\n```python\n\n        class MarkataMarkdown(Markata):\n
  \               articles: List = []\n```\n\n"
date: 0001-01-01
description: 'Default load plugin. ! ???+ source  ! ???+ source  ! ???+ source  !
  ???+ source  ! ???+ source  ! ???+ source  ! ???+ source '
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Load.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"Default load plugin. ! ???+ source  ! ???+ source
    \ ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source \"
    />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\"
    href=\"/post.css\" />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script
    src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    <head>\n<title>Load.Py</title>\n<meta charset=\"UTF-8\" />\n<meta name=\"viewport\"
    content=\"width=device-width, initial-scale=1\" />\n<meta name=\"description\"
    content=\"Default load plugin. ! ???+ source  ! ???+ source  ! ???+ source  !
    ???+ source  ! ???+ source  ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\"
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Load.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <p>Default load plugin.</p>\n<p>!!
    class <h2 id='ValidationError' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>ValidationError <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">ValidationError
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
    <span class=\"nc\">ValidationError</span><span class=\"p\">(</span><span class=\"ne\">ValueError</span><span
    class=\"p\">):</span> <span class=\"o\">...</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='load' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>load <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">load
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
    <span class=\"nf\">load</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;MarkataMarkdown&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n            <span class=\"n\">Progress</span><span class=\"p\">(</span>\n
    \               <span class=\"n\">BarColumn</span><span class=\"p\">(</span><span
    class=\"n\">bar_width</span><span class=\"o\">=</span><span class=\"kc\">None</span><span
    class=\"p\">),</span>\n                <span class=\"n\">transient</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">console</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"p\">,</span>\n            <span class=\"p\">)</span>\n            <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">log</span><span class=\"p\">(</span><span
    class=\"sa\">f</span><span class=\"s2\">&quot;found </span><span class=\"si\">{</span><span
    class=\"nb\">len</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">files</span><span class=\"p\">)</span><span
    class=\"si\">}</span><span class=\"s2\"> posts&quot;</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">post_futures</span> <span class=\"o\">=</span> <span
    class=\"p\">[</span><span class=\"n\">get_post</span><span class=\"p\">(</span><span
    class=\"n\">article</span><span class=\"p\">,</span> <span class=\"n\">markata</span><span
    class=\"p\">)</span> <span class=\"k\">for</span> <span class=\"n\">article</span>
    <span class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">files</span><span class=\"p\">]</span>\n            <span class=\"n\">posts</span>
    <span class=\"o\">=</span> <span class=\"p\">[</span><span class=\"n\">post</span><span
    class=\"o\">.</span><span class=\"n\">result</span><span class=\"p\">()</span>
    <span class=\"k\">for</span> <span class=\"n\">post</span> <span class=\"ow\">in</span>
    <span class=\"n\">post_futures</span> <span class=\"k\">if</span> <span class=\"n\">post</span>
    <span class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span><span
    class=\"p\">]</span>\n\n            <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">posts_obj</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">Posts</span><span class=\"o\">.</span><span
    class=\"n\">parse_obj</span><span class=\"p\">(</span>\n                <span
    class=\"p\">{</span><span class=\"s2\">&quot;posts&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">posts</span><span class=\"p\">},</span>\n            <span class=\"p\">)</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">posts</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">posts_obj</span><span class=\"o\">.</span><span class=\"n\">posts</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">articles</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">posts</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='get_post'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get_post <em class='small'>function</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">get_post
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
    <span class=\"nf\">get_post</span><span class=\"p\">(</span><span class=\"n\">path</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">,</span> <span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Optional</span><span
    class=\"p\">[</span><span class=\"n\">Callable</span><span class=\"p\">]:</span>\n
    \           <span class=\"k\">if</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">Post</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">post</span> <span class=\"o\">=</span> <span
    class=\"n\">pydantic_get_post</span><span class=\"p\">(</span><span class=\"n\">path</span><span
    class=\"o\">=</span><span class=\"n\">path</span><span class=\"p\">,</span> <span
    class=\"n\">markata</span><span class=\"o\">=</span><span class=\"n\">markata</span><span
    class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"n\">post</span>\n
    \           <span class=\"k\">else</span><span class=\"p\">:</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">legacy_get_post</span><span class=\"p\">(</span><span
    class=\"n\">path</span><span class=\"o\">=</span><span class=\"n\">path</span><span
    class=\"p\">,</span> <span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='get_models' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>get_models <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">get_models
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
    <span class=\"nf\">get_models</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">error</span><span class=\"p\">:</span> <span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">ValidationError</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">List</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">fields</span> <span class=\"o\">=</span> <span class=\"p\">[]</span>\n
    \           <span class=\"k\">for</span> <span class=\"n\">err</span> <span class=\"ow\">in</span>
    <span class=\"n\">error</span><span class=\"o\">.</span><span class=\"n\">errors</span><span
    class=\"p\">():</span>\n                <span class=\"n\">fields</span><span class=\"o\">.</span><span
    class=\"n\">extend</span><span class=\"p\">(</span><span class=\"n\">err</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;loc&quot;</span><span class=\"p\">])</span>\n\n
    \           <span class=\"n\">models</span> <span class=\"o\">=</span> <span class=\"p\">{</span><span
    class=\"n\">field</span><span class=\"p\">:</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;</span><span class=\"si\">{</span><span class=\"n\">field</span><span
    class=\"si\">}</span><span class=\"s2\"> used by &quot;</span> <span class=\"k\">for</span>
    <span class=\"n\">field</span> <span class=\"ow\">in</span> <span class=\"n\">fields</span><span
    class=\"p\">}</span>\n\n            <span class=\"k\">for</span> <span class=\"n\">field</span><span
    class=\"p\">,</span> <span class=\"n\">model</span> <span class=\"ow\">in</span>
    <span class=\"nb\">set</span><span class=\"p\">(</span>\n                <span
    class=\"n\">itertools</span><span class=\"o\">.</span><span class=\"n\">product</span><span
    class=\"p\">(</span>\n                    <span class=\"n\">fields</span><span
    class=\"p\">,</span>\n                    <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">post_models</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">),</span>\n            <span class=\"p\">):</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">field</span> <span
    class=\"ow\">in</span> <span class=\"n\">model</span><span class=\"o\">.</span><span
    class=\"n\">__fields__</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">models</span><span class=\"p\">[</span><span class=\"n\">field</span><span
    class=\"p\">]</span> <span class=\"o\">+=</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;&#39;</span><span class=\"si\">{</span><span class=\"n\">model</span><span
    class=\"o\">.</span><span class=\"vm\">__module__</span><span class=\"si\">}</span><span
    class=\"s2\">.</span><span class=\"si\">{</span><span class=\"n\">model</span><span
    class=\"o\">.</span><span class=\"vm\">__name__</span><span class=\"si\">}</span><span
    class=\"s2\">&#39;&quot;</span>\n\n            <span class=\"k\">return</span>
    <span class=\"n\">models</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2
    id='pydantic_get_post' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>pydantic_get_post <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">pydantic_get_post
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
    <span class=\"nf\">pydantic_get_post</span><span class=\"p\">(</span><span class=\"n\">path</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">,</span> <span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Optional</span><span
    class=\"p\">[</span><span class=\"n\">Callable</span><span class=\"p\">]:</span>\n
    \           <span class=\"k\">try</span><span class=\"p\">:</span>\n                <span
    class=\"n\">post</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">Post</span><span class=\"o\">.</span><span
    class=\"n\">parse_file</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">path</span><span class=\"o\">=</span><span class=\"n\">path</span><span
    class=\"p\">)</span>\n                <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">Post</span><span class=\"o\">.</span><span class=\"n\">validate</span><span
    class=\"p\">(</span><span class=\"n\">post</span><span class=\"p\">)</span>\n\n
    \           <span class=\"k\">except</span> <span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">ValidationError</span> <span class=\"k\">as</span>
    <span class=\"n\">e</span><span class=\"p\">:</span>\n                <span class=\"n\">models</span>
    <span class=\"o\">=</span> <span class=\"n\">get_models</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"o\">=</span><span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"n\">error</span><span class=\"o\">=</span><span
    class=\"n\">e</span><span class=\"p\">)</span>\n                <span class=\"n\">models</span>
    <span class=\"o\">=</span> <span class=\"nb\">list</span><span class=\"p\">(</span><span
    class=\"n\">models</span><span class=\"o\">.</span><span class=\"n\">values</span><span
    class=\"p\">())</span>\n                <span class=\"n\">models</span> <span
    class=\"o\">=</span> <span class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span
    class=\"s2\">&quot;</span><span class=\"o\">.</span><span class=\"n\">join</span><span
    class=\"p\">(</span><span class=\"n\">models</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">raise</span> <span class=\"n\">ValidationError</span><span
    class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">e</span><span class=\"si\">}</span><span
    class=\"se\">\\n\\n</span><span class=\"si\">{</span><span class=\"n\">models</span><span
    class=\"si\">}</span><span class=\"se\">\\n</span><span class=\"s2\">failed to
    load </span><span class=\"si\">{</span><span class=\"n\">path</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span><span class=\"p\">)</span> <span class=\"kn\">from</span>
    <span class=\"nn\">e</span>\n\n            <span class=\"k\">return</span> <span
    class=\"n\">post</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='legacy_get_post'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>legacy_get_post
    <em class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">legacy_get_post <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">legacy_get_post</span><span class=\"p\">(</span><span class=\"n\">path</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">,</span> <span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Optional</span><span
    class=\"p\">[</span><span class=\"n\">Callable</span><span class=\"p\">]:</span>\n
    \           <span class=\"n\">default</span> <span class=\"o\">=</span> <span
    class=\"p\">{</span>\n                <span class=\"s2\">&quot;cover&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"s2\">&quot;title&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;&quot;</span><span class=\"p\">,</span>\n                <span
    class=\"s2\">&quot;tags&quot;</span><span class=\"p\">:</span> <span class=\"p\">[],</span>\n
    \               <span class=\"s2\">&quot;published&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;False&quot;</span><span class=\"p\">,</span>\n                <span
    class=\"s2\">&quot;templateKey&quot;</span><span class=\"p\">:</span> <span class=\"s2\">&quot;&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;path&quot;</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">path</span><span class=\"p\">),</span>\n                <span class=\"s2\">&quot;description&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"s2\">&quot;content&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;&quot;</span><span class=\"p\">,</span>\n            <span
    class=\"p\">}</span>\n            <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">post</span><span class=\"p\">:</span> <span
    class=\"s2\">&quot;Post&quot;</span> <span class=\"o\">=</span> <span class=\"n\">frontmatter</span><span
    class=\"o\">.</span><span class=\"n\">load</span><span class=\"p\">(</span><span
    class=\"n\">path</span><span class=\"p\">)</span>\n                <span class=\"n\">post</span><span
    class=\"o\">.</span><span class=\"n\">metadata</span> <span class=\"o\">=</span>
    <span class=\"p\">{</span><span class=\"o\">**</span><span class=\"n\">default</span><span
    class=\"p\">,</span> <span class=\"o\">**</span><span class=\"n\">post</span><span
    class=\"o\">.</span><span class=\"n\">metadata</span><span class=\"p\">}</span>\n
    \               <span class=\"n\">post</span><span class=\"p\">[</span><span class=\"s2\">&quot;content&quot;</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"n\">post</span><span
    class=\"o\">.</span><span class=\"n\">content</span>\n            <span class=\"k\">except</span>
    <span class=\"n\">ParserError</span><span class=\"p\">:</span>\n                <span
    class=\"k\">return</span> <span class=\"kc\">None</span>\n                <span
    class=\"n\">post</span> <span class=\"o\">=</span> <span class=\"n\">default</span>\n
    \           <span class=\"k\">except</span> <span class=\"ne\">ValueError</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"kc\">None</span>\n
    \               <span class=\"n\">post</span> <span class=\"o\">=</span> <span
    class=\"n\">default</span>\n            <span class=\"n\">post</span><span class=\"o\">.</span><span
    class=\"n\">metadata</span><span class=\"p\">[</span><span class=\"s2\">&quot;path&quot;</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"nb\">str</span><span
    class=\"p\">(</span><span class=\"n\">path</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">post</span><span class=\"p\">[</span><span class=\"s2\">&quot;edit_link&quot;</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"p\">(</span>\n                <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">repo_url</span> <span class=\"o\">+</span>
    <span class=\"s2\">&quot;edit/&quot;</span> <span class=\"o\">+</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">repo_branch</span> <span class=\"o\">+</span> <span class=\"s2\">&quot;/&quot;</span>
    <span class=\"o\">+</span> <span class=\"n\">post</span><span class=\"o\">.</span><span
    class=\"n\">path</span>\n            <span class=\"p\">)</span>\n            <span
    class=\"k\">return</span> <span class=\"n\">post</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='MarkataMarkdown' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>MarkataMarkdown <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">MarkataMarkdown
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
    <span class=\"nc\">MarkataMarkdown</span><span class=\"p\">(</span><span class=\"n\">Markata</span><span
    class=\"p\">):</span>\n                <span class=\"n\">articles</span><span
    class=\"p\">:</span> <span class=\"n\">List</span> <span class=\"o\">=</span>
    <span class=\"p\">[]</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>
    \   </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Load.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"Default load plugin. ! ???+ source  ! ???+ source
    \ ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source \"
    />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\"
    href=\"/post.css\" />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script
    src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    <head>\n<title>Load.Py</title>\n<meta charset=\"UTF-8\" />\n<meta name=\"viewport\"
    content=\"width=device-width, initial-scale=1\" />\n<meta name=\"description\"
    content=\"Default load plugin. ! ???+ source  ! ???+ source  ! ???+ source  !
    ???+ source  ! ???+ source  ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<article style=\"text-align: center;\">\n    <style>\n        section
    {\n            font-size: 200%;\n        }\n\n\n        .edit {\n            display:
    none;\n        }\n    </style>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Load.Py\n    </h1>\n</section></article>\n     </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Load.Py\n    </h1>\n</section>    <section class=\"body\">\n        <p>Default
    load plugin.</p>\n<p>!! class <h2 id='ValidationError' class='admonition-title'
    style='margin:0;padding:.5rem 1rem;'>ValidationError <em class='small'>class</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">ValidationError
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
    <span class=\"nc\">ValidationError</span><span class=\"p\">(</span><span class=\"ne\">ValueError</span><span
    class=\"p\">):</span> <span class=\"o\">...</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='load' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>load <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">load
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
    <span class=\"nf\">load</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;MarkataMarkdown&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n            <span class=\"n\">Progress</span><span class=\"p\">(</span>\n
    \               <span class=\"n\">BarColumn</span><span class=\"p\">(</span><span
    class=\"n\">bar_width</span><span class=\"o\">=</span><span class=\"kc\">None</span><span
    class=\"p\">),</span>\n                <span class=\"n\">transient</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">console</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"p\">,</span>\n            <span class=\"p\">)</span>\n            <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">log</span><span class=\"p\">(</span><span
    class=\"sa\">f</span><span class=\"s2\">&quot;found </span><span class=\"si\">{</span><span
    class=\"nb\">len</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">files</span><span class=\"p\">)</span><span
    class=\"si\">}</span><span class=\"s2\"> posts&quot;</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">post_futures</span> <span class=\"o\">=</span> <span
    class=\"p\">[</span><span class=\"n\">get_post</span><span class=\"p\">(</span><span
    class=\"n\">article</span><span class=\"p\">,</span> <span class=\"n\">markata</span><span
    class=\"p\">)</span> <span class=\"k\">for</span> <span class=\"n\">article</span>
    <span class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">files</span><span class=\"p\">]</span>\n            <span class=\"n\">posts</span>
    <span class=\"o\">=</span> <span class=\"p\">[</span><span class=\"n\">post</span><span
    class=\"o\">.</span><span class=\"n\">result</span><span class=\"p\">()</span>
    <span class=\"k\">for</span> <span class=\"n\">post</span> <span class=\"ow\">in</span>
    <span class=\"n\">post_futures</span> <span class=\"k\">if</span> <span class=\"n\">post</span>
    <span class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span><span
    class=\"p\">]</span>\n\n            <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">posts_obj</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">Posts</span><span class=\"o\">.</span><span
    class=\"n\">parse_obj</span><span class=\"p\">(</span>\n                <span
    class=\"p\">{</span><span class=\"s2\">&quot;posts&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">posts</span><span class=\"p\">},</span>\n            <span class=\"p\">)</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">posts</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">posts_obj</span><span class=\"o\">.</span><span class=\"n\">posts</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">articles</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">posts</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='get_post'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get_post <em class='small'>function</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">get_post
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
    <span class=\"nf\">get_post</span><span class=\"p\">(</span><span class=\"n\">path</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">,</span> <span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Optional</span><span
    class=\"p\">[</span><span class=\"n\">Callable</span><span class=\"p\">]:</span>\n
    \           <span class=\"k\">if</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">Post</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">post</span> <span class=\"o\">=</span> <span
    class=\"n\">pydantic_get_post</span><span class=\"p\">(</span><span class=\"n\">path</span><span
    class=\"o\">=</span><span class=\"n\">path</span><span class=\"p\">,</span> <span
    class=\"n\">markata</span><span class=\"o\">=</span><span class=\"n\">markata</span><span
    class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"n\">post</span>\n
    \           <span class=\"k\">else</span><span class=\"p\">:</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">legacy_get_post</span><span class=\"p\">(</span><span
    class=\"n\">path</span><span class=\"o\">=</span><span class=\"n\">path</span><span
    class=\"p\">,</span> <span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='get_models' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>get_models <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">get_models
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
    <span class=\"nf\">get_models</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">error</span><span class=\"p\">:</span> <span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">ValidationError</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">List</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">fields</span> <span class=\"o\">=</span> <span class=\"p\">[]</span>\n
    \           <span class=\"k\">for</span> <span class=\"n\">err</span> <span class=\"ow\">in</span>
    <span class=\"n\">error</span><span class=\"o\">.</span><span class=\"n\">errors</span><span
    class=\"p\">():</span>\n                <span class=\"n\">fields</span><span class=\"o\">.</span><span
    class=\"n\">extend</span><span class=\"p\">(</span><span class=\"n\">err</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;loc&quot;</span><span class=\"p\">])</span>\n\n
    \           <span class=\"n\">models</span> <span class=\"o\">=</span> <span class=\"p\">{</span><span
    class=\"n\">field</span><span class=\"p\">:</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;</span><span class=\"si\">{</span><span class=\"n\">field</span><span
    class=\"si\">}</span><span class=\"s2\"> used by &quot;</span> <span class=\"k\">for</span>
    <span class=\"n\">field</span> <span class=\"ow\">in</span> <span class=\"n\">fields</span><span
    class=\"p\">}</span>\n\n            <span class=\"k\">for</span> <span class=\"n\">field</span><span
    class=\"p\">,</span> <span class=\"n\">model</span> <span class=\"ow\">in</span>
    <span class=\"nb\">set</span><span class=\"p\">(</span>\n                <span
    class=\"n\">itertools</span><span class=\"o\">.</span><span class=\"n\">product</span><span
    class=\"p\">(</span>\n                    <span class=\"n\">fields</span><span
    class=\"p\">,</span>\n                    <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">post_models</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">),</span>\n            <span class=\"p\">):</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">field</span> <span
    class=\"ow\">in</span> <span class=\"n\">model</span><span class=\"o\">.</span><span
    class=\"n\">__fields__</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">models</span><span class=\"p\">[</span><span class=\"n\">field</span><span
    class=\"p\">]</span> <span class=\"o\">+=</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;&#39;</span><span class=\"si\">{</span><span class=\"n\">model</span><span
    class=\"o\">.</span><span class=\"vm\">__module__</span><span class=\"si\">}</span><span
    class=\"s2\">.</span><span class=\"si\">{</span><span class=\"n\">model</span><span
    class=\"o\">.</span><span class=\"vm\">__name__</span><span class=\"si\">}</span><span
    class=\"s2\">&#39;&quot;</span>\n\n            <span class=\"k\">return</span>
    <span class=\"n\">models</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2
    id='pydantic_get_post' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>pydantic_get_post <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">pydantic_get_post
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
    <span class=\"nf\">pydantic_get_post</span><span class=\"p\">(</span><span class=\"n\">path</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">,</span> <span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Optional</span><span
    class=\"p\">[</span><span class=\"n\">Callable</span><span class=\"p\">]:</span>\n
    \           <span class=\"k\">try</span><span class=\"p\">:</span>\n                <span
    class=\"n\">post</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">Post</span><span class=\"o\">.</span><span
    class=\"n\">parse_file</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">path</span><span class=\"o\">=</span><span class=\"n\">path</span><span
    class=\"p\">)</span>\n                <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">Post</span><span class=\"o\">.</span><span class=\"n\">validate</span><span
    class=\"p\">(</span><span class=\"n\">post</span><span class=\"p\">)</span>\n\n
    \           <span class=\"k\">except</span> <span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">ValidationError</span> <span class=\"k\">as</span>
    <span class=\"n\">e</span><span class=\"p\">:</span>\n                <span class=\"n\">models</span>
    <span class=\"o\">=</span> <span class=\"n\">get_models</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"o\">=</span><span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"n\">error</span><span class=\"o\">=</span><span
    class=\"n\">e</span><span class=\"p\">)</span>\n                <span class=\"n\">models</span>
    <span class=\"o\">=</span> <span class=\"nb\">list</span><span class=\"p\">(</span><span
    class=\"n\">models</span><span class=\"o\">.</span><span class=\"n\">values</span><span
    class=\"p\">())</span>\n                <span class=\"n\">models</span> <span
    class=\"o\">=</span> <span class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span
    class=\"s2\">&quot;</span><span class=\"o\">.</span><span class=\"n\">join</span><span
    class=\"p\">(</span><span class=\"n\">models</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">raise</span> <span class=\"n\">ValidationError</span><span
    class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">e</span><span class=\"si\">}</span><span
    class=\"se\">\\n\\n</span><span class=\"si\">{</span><span class=\"n\">models</span><span
    class=\"si\">}</span><span class=\"se\">\\n</span><span class=\"s2\">failed to
    load </span><span class=\"si\">{</span><span class=\"n\">path</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span><span class=\"p\">)</span> <span class=\"kn\">from</span>
    <span class=\"nn\">e</span>\n\n            <span class=\"k\">return</span> <span
    class=\"n\">post</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='legacy_get_post'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>legacy_get_post
    <em class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">legacy_get_post <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">legacy_get_post</span><span class=\"p\">(</span><span class=\"n\">path</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">,</span> <span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Optional</span><span
    class=\"p\">[</span><span class=\"n\">Callable</span><span class=\"p\">]:</span>\n
    \           <span class=\"n\">default</span> <span class=\"o\">=</span> <span
    class=\"p\">{</span>\n                <span class=\"s2\">&quot;cover&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"s2\">&quot;title&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;&quot;</span><span class=\"p\">,</span>\n                <span
    class=\"s2\">&quot;tags&quot;</span><span class=\"p\">:</span> <span class=\"p\">[],</span>\n
    \               <span class=\"s2\">&quot;published&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;False&quot;</span><span class=\"p\">,</span>\n                <span
    class=\"s2\">&quot;templateKey&quot;</span><span class=\"p\">:</span> <span class=\"s2\">&quot;&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;path&quot;</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">path</span><span class=\"p\">),</span>\n                <span class=\"s2\">&quot;description&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"s2\">&quot;content&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;&quot;</span><span class=\"p\">,</span>\n            <span
    class=\"p\">}</span>\n            <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">post</span><span class=\"p\">:</span> <span
    class=\"s2\">&quot;Post&quot;</span> <span class=\"o\">=</span> <span class=\"n\">frontmatter</span><span
    class=\"o\">.</span><span class=\"n\">load</span><span class=\"p\">(</span><span
    class=\"n\">path</span><span class=\"p\">)</span>\n                <span class=\"n\">post</span><span
    class=\"o\">.</span><span class=\"n\">metadata</span> <span class=\"o\">=</span>
    <span class=\"p\">{</span><span class=\"o\">**</span><span class=\"n\">default</span><span
    class=\"p\">,</span> <span class=\"o\">**</span><span class=\"n\">post</span><span
    class=\"o\">.</span><span class=\"n\">metadata</span><span class=\"p\">}</span>\n
    \               <span class=\"n\">post</span><span class=\"p\">[</span><span class=\"s2\">&quot;content&quot;</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"n\">post</span><span
    class=\"o\">.</span><span class=\"n\">content</span>\n            <span class=\"k\">except</span>
    <span class=\"n\">ParserError</span><span class=\"p\">:</span>\n                <span
    class=\"k\">return</span> <span class=\"kc\">None</span>\n                <span
    class=\"n\">post</span> <span class=\"o\">=</span> <span class=\"n\">default</span>\n
    \           <span class=\"k\">except</span> <span class=\"ne\">ValueError</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"kc\">None</span>\n
    \               <span class=\"n\">post</span> <span class=\"o\">=</span> <span
    class=\"n\">default</span>\n            <span class=\"n\">post</span><span class=\"o\">.</span><span
    class=\"n\">metadata</span><span class=\"p\">[</span><span class=\"s2\">&quot;path&quot;</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"nb\">str</span><span
    class=\"p\">(</span><span class=\"n\">path</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">post</span><span class=\"p\">[</span><span class=\"s2\">&quot;edit_link&quot;</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"p\">(</span>\n                <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">repo_url</span> <span class=\"o\">+</span>
    <span class=\"s2\">&quot;edit/&quot;</span> <span class=\"o\">+</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">repo_branch</span> <span class=\"o\">+</span> <span class=\"s2\">&quot;/&quot;</span>
    <span class=\"o\">+</span> <span class=\"n\">post</span><span class=\"o\">.</span><span
    class=\"n\">path</span>\n            <span class=\"p\">)</span>\n            <span
    class=\"k\">return</span> <span class=\"n\">post</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='MarkataMarkdown' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>MarkataMarkdown <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">MarkataMarkdown
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
    <span class=\"nc\">MarkataMarkdown</span><span class=\"p\">(</span><span class=\"n\">Markata</span><span
    class=\"p\">):</span>\n                <span class=\"n\">articles</span><span
    class=\"p\">:</span> <span class=\"n\">List</span> <span class=\"o\">=</span>
    <span class=\"p\">[]</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>"
  raw.md: ''
published: true
slug: markata/plugins/load
title: Load.Py


---

Default load plugin.


!! class <h2 id='ValidationError' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>ValidationError <em class='small'>class</em></h2>

???+ source "ValidationError <em class='small'>source</em>"

```python

        class ValidationError(ValueError): ...
```


!! function <h2 id='load' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>load <em class='small'>function</em></h2>

???+ source "load <em class='small'>source</em>"

```python

        def load(markata: "MarkataMarkdown") -> None:
            Progress(
                BarColumn(bar_width=None),
                transient=True,
                console=markata.console,
            )
            markata.console.log(f"found {len(markata.files)} posts")
            post_futures = [get_post(article, markata) for article in markata.files]
            posts = [post.result() for post in post_futures if post is not None]

            markata.posts_obj = markata.Posts.parse_obj(
                {"posts": posts},
            )
            markata.posts = markata.posts_obj.posts
            markata.articles = markata.posts
```


!! function <h2 id='get_post' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get_post <em class='small'>function</em></h2>

???+ source "get_post <em class='small'>source</em>"

```python

        def get_post(path: Path, markata: "Markata") -> Optional[Callable]:
            if markata.Post:
                post = pydantic_get_post(path=path, markata=markata)
                return post
            else:
                return legacy_get_post(path=path, markata=markata)
```


!! function <h2 id='get_models' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get_models <em class='small'>function</em></h2>

???+ source "get_models <em class='small'>source</em>"

```python

        def get_models(markata: "Markata", error: pydantic.ValidationError) -> List:
            fields = []
            for err in error.errors():
                fields.extend(err["loc"])

            models = {field: f"{field} used by " for field in fields}

            for field, model in set(
                itertools.product(
                    fields,
                    markata.post_models,
                ),
            ):
                if field in model.__fields__:
                    models[field] += f"'{model.__module__}.{model.__name__}'"

            return models
```


!! function <h2 id='pydantic_get_post' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>pydantic_get_post <em class='small'>function</em></h2>

???+ source "pydantic_get_post <em class='small'>source</em>"

```python

        def pydantic_get_post(path: Path, markata: "Markata") -> Optional[Callable]:
            try:
                post = markata.Post.parse_file(markata=markata, path=path)
                markata.Post.validate(post)

            except pydantic.ValidationError as e:
                models = get_models(markata=markata, error=e)
                models = list(models.values())
                models = "\n".join(models)
                raise ValidationError(f"{e}\n\n{models}\nfailed to load {path}") from e

            return post
```


!! function <h2 id='legacy_get_post' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>legacy_get_post <em class='small'>function</em></h2>

???+ source "legacy_get_post <em class='small'>source</em>"

```python

        def legacy_get_post(path: Path, markata: "Markata") -> Optional[Callable]:
            default = {
                "cover": "",
                "title": "",
                "tags": [],
                "published": "False",
                "templateKey": "",
                "path": str(path),
                "description": "",
                "content": "",
            }
            try:
                post: "Post" = frontmatter.load(path)
                post.metadata = {**default, **post.metadata}
                post["content"] = post.content
            except ParserError:
                return None
                post = default
            except ValueError:
                return None
                post = default
            post.metadata["path"] = str(path)
            post["edit_link"] = (
                markata.config.repo_url + "edit/" + markata.config.repo_branch + "/" + post.path
            )
            return post
```


!! class <h2 id='MarkataMarkdown' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>MarkataMarkdown <em class='small'>class</em></h2>

???+ source "MarkataMarkdown <em class='small'>source</em>"

```python

        class MarkataMarkdown(Markata):
                articles: List = []
```

