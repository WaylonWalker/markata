---
content: "Writes the final modified markdown and frontmatter to the output directory.\nReplacing
  the trailing slash if its there and adding .md will bring up the raw\nsource.\n\n##
  Configuration\n\nThe only configuration for the publish_source plugin is to make
  sure its in\nyour list of hooks.\n\n\n``` toml\n[markata]\n\n# make sure its in
  your list of hooks\nhooks=[\n   \"markata.plugins.publish_source\",\n   ]\n```\n\n!!!
  note\n    publish_source is included by default, but if you have not included the\n
  \   default set of hooks you will need to explicitly add it.\n\n\n!! function <h2
  id='_save' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_save <em
  class='small'>function</em></h2>\n    saves the article to the output directory
  at its specified slug.\n???+ source \"_save <em class='small'>source</em>\"\n\n```python\n\n
  \       def _save(output_dir: Path, article: frontmatter.Post) -> None:\n            \"\"\"\n
  \           saves the article to the output directory at its specified slug.\n            \"\"\"\n
  \           path = Path(\n                output_dir / Path(article[\"slug\"]).parent
  / Path(article[\"path\"]).name,\n            )\n            path.parent.mkdir(parents=True,
  exist_ok=True)\n            path.write_text(article.dumps())\n```\n\n\n!! function
  <h2 id='_strip_unserializable_values' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>_strip_unserializable_values <em class='small'>function</em></h2>\n    Returns
  an article with only yaml serializable frontmatter.\n???+ source \"_strip_unserializable_values
  <em class='small'>source</em>\"\n\n```python\n\n        def _strip_unserializable_values(\n
  \           markata: \"Markata\",\n            article: frontmatter.Post,\n        )
  -> frontmatter.Post:\n            \"\"\"\n            Returns an article with only
  yaml serializable frontmatter.\n            \"\"\"\n            _article = frontmatter.Post(\n
  \               article.content,\n                **{k: v for k, v in article.metadata.items()
  if k != \"content\"},\n            )\n            kwargs = {\n                \"Dumper\":
  yaml.cyaml.CSafeDumper,\n                \"default_flow_style\": False,\n                \"allow_unicode\":
  True,\n            }\n            for key, value in article.metadata.items():\n
  \               try:\n                    yaml.dump({key: value}, **kwargs)\n                except
  RepresenterError:\n                    del _article[key]\n            if markata.Post:\n
  \               _article = markata.Post(**_article.metadata, path=str(article.path))\n
  \           return _article\n```\n\n\n!! function <h2 id='save' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>save <em class='small'>function</em></h2>\n
  \   Saves the final modified post to the output site as markdown.\n\n    !!! note\n
  \       Any keys that are not yaml serializable will be stripped.\n???+ source \"save
  <em class='small'>source</em>\"\n\n```python\n\n        def save(markata: \"Markata\")
  -> None:\n            \"\"\"\n            Saves the final modified post to the output
  site as markdown.\n\n            !!! note\n                Any keys that are not
  yaml serializable will be stripped.\n\n            \"\"\"\n            output_dir
  = Path(str(markata.config[\"output_dir\"]))\n            for (\n                article\n
  \           ) in markata.articles:  # iter_articles(description=\"saving source
  documents\"):\n                try:\n                    _save(output_dir, article)\n
  \               except RepresenterError:\n                    _article = _strip_unserializable_values(markata,
  article)\n\n                    _save(output_dir, _article)\n```\n\n"
date: 0001-01-01
description: 'Writes the final modified markdown and frontmatter to the output directory.
  The only configuration for the publish ! ! ! ! ???+ source '
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Publish_Source.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Writes the final modified markdown and
    frontmatter to the output directory. The only configuration for the publish !
    ! ! ! ???+ source \" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Publish_Source.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Writes the final modified markdown and
    frontmatter to the output directory. The only configuration for the publish !
    ! ! ! ???+ source \" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Publish_Source.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <p>Writes the final
    modified markdown and frontmatter to the output directory.\nReplacing the trailing
    slash if its there and adding .md will bring up the raw\nsource.</p>\n<h2 id=\"configuration\">Configuration
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>The only configuration
    for the publish_source plugin is to make sure its in\nyour list of hooks.</p>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n\n<span
    class=\"c1\"># make sure its in your list of hooks</span>\n<span class=\"n\">hooks</span><span
    class=\"o\">=</span><span class=\"p\">[</span>\n<span class=\"w\">   </span><span
    class=\"s2\">&quot;markata.plugins.publish_source&quot;</span><span class=\"p\">,</span>\n<span
    class=\"w\">   </span><span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<div
    class=\"admonition note\">\n<p class=\"admonition-title\">Note</p>\n<p>publish_source
    is included by default, but if you have not included the\ndefault set of hooks
    you will need to explicitly add it.</p>\n</div>\n<p>!! function <h2 id='_save'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_save <em class='small'>function</em></h2>\nsaves
    the article to the output directory at its specified slug.</p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_save
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
    <span class=\"nf\">_save</span><span class=\"p\">(</span><span class=\"n\">output_dir</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">,</span> <span
    class=\"n\">article</span><span class=\"p\">:</span> <span class=\"n\">frontmatter</span><span
    class=\"o\">.</span><span class=\"n\">Post</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            saves the article to the output directory at its specified
    slug.</span>\n<span class=\"sd\">            &quot;&quot;&quot;</span>\n            <span
    class=\"n\">path</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span>\n                <span class=\"n\">output_dir</span> <span
    class=\"o\">/</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">article</span><span class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">])</span><span class=\"o\">.</span><span class=\"n\">parent</span>
    <span class=\"o\">/</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">article</span><span class=\"p\">[</span><span class=\"s2\">&quot;path&quot;</span><span
    class=\"p\">])</span><span class=\"o\">.</span><span class=\"n\">name</span><span
    class=\"p\">,</span>\n            <span class=\"p\">)</span>\n            <span
    class=\"n\">path</span><span class=\"o\">.</span><span class=\"n\">parent</span><span
    class=\"o\">.</span><span class=\"n\">mkdir</span><span class=\"p\">(</span><span
    class=\"n\">parents</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span> <span class=\"n\">exist_ok</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n            <span class=\"n\">path</span><span
    class=\"o\">.</span><span class=\"n\">write_text</span><span class=\"p\">(</span><span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">dumps</span><span
    class=\"p\">())</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='_strip_unserializable_values'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_strip_unserializable_values
    <em class='small'>function</em></h2>\nReturns an article with only yaml serializable
    frontmatter.</p>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">_strip_unserializable_values <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">_strip_unserializable_values</span><span class=\"p\">(</span>\n
    \           <span class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span>\n            <span class=\"n\">article</span><span class=\"p\">:</span>
    <span class=\"n\">frontmatter</span><span class=\"o\">.</span><span class=\"n\">Post</span><span
    class=\"p\">,</span>\n        <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">frontmatter</span><span class=\"o\">.</span><span class=\"n\">Post</span><span
    class=\"p\">:</span>\n<span class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            Returns an article with only yaml serializable frontmatter.</span>\n<span
    class=\"sd\">            &quot;&quot;&quot;</span>\n            <span class=\"n\">_article</span>
    <span class=\"o\">=</span> <span class=\"n\">frontmatter</span><span class=\"o\">.</span><span
    class=\"n\">Post</span><span class=\"p\">(</span>\n                <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">content</span><span class=\"p\">,</span>\n
    \               <span class=\"o\">**</span><span class=\"p\">{</span><span class=\"n\">k</span><span
    class=\"p\">:</span> <span class=\"n\">v</span> <span class=\"k\">for</span> <span
    class=\"n\">k</span><span class=\"p\">,</span> <span class=\"n\">v</span> <span
    class=\"ow\">in</span> <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">metadata</span><span class=\"o\">.</span><span class=\"n\">items</span><span
    class=\"p\">()</span> <span class=\"k\">if</span> <span class=\"n\">k</span> <span
    class=\"o\">!=</span> <span class=\"s2\">&quot;content&quot;</span><span class=\"p\">},</span>\n
    \           <span class=\"p\">)</span>\n            <span class=\"n\">kwargs</span>
    <span class=\"o\">=</span> <span class=\"p\">{</span>\n                <span class=\"s2\">&quot;Dumper&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">yaml</span><span class=\"o\">.</span><span
    class=\"n\">cyaml</span><span class=\"o\">.</span><span class=\"n\">CSafeDumper</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;default_flow_style&quot;</span><span
    class=\"p\">:</span> <span class=\"kc\">False</span><span class=\"p\">,</span>\n
    \               <span class=\"s2\">&quot;allow_unicode&quot;</span><span class=\"p\">:</span>
    <span class=\"kc\">True</span><span class=\"p\">,</span>\n            <span class=\"p\">}</span>\n
    \           <span class=\"k\">for</span> <span class=\"n\">key</span><span class=\"p\">,</span>
    <span class=\"n\">value</span> <span class=\"ow\">in</span> <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">metadata</span><span class=\"o\">.</span><span
    class=\"n\">items</span><span class=\"p\">():</span>\n                <span class=\"k\">try</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">yaml</span><span class=\"o\">.</span><span
    class=\"n\">dump</span><span class=\"p\">({</span><span class=\"n\">key</span><span
    class=\"p\">:</span> <span class=\"n\">value</span><span class=\"p\">},</span>
    <span class=\"o\">**</span><span class=\"n\">kwargs</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">except</span> <span class=\"n\">RepresenterError</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">del</span> <span class=\"n\">_article</span><span
    class=\"p\">[</span><span class=\"n\">key</span><span class=\"p\">]</span>\n            <span
    class=\"k\">if</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">Post</span><span class=\"p\">:</span>\n                <span class=\"n\">_article</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">Post</span><span class=\"p\">(</span><span class=\"o\">**</span><span
    class=\"n\">_article</span><span class=\"o\">.</span><span class=\"n\">metadata</span><span
    class=\"p\">,</span> <span class=\"n\">path</span><span class=\"o\">=</span><span
    class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">path</span><span class=\"p\">))</span>\n
    \           <span class=\"k\">return</span> <span class=\"n\">_article</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='save' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>save <em class='small'>function</em></h2>\nSaves the final modified post
    to the output site as markdown.</p>\n<pre><code>!!! note\n    Any keys that are
    not yaml serializable will be stripped.\n</code></pre>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">save
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
    <span class=\"nf\">save</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            Saves the final modified post to the output site as markdown.</span>\n\n<span
    class=\"sd\">            !!! note</span>\n<span class=\"sd\">                Any
    keys that are not yaml serializable will be stripped.</span>\n\n<span class=\"sd\">
    \           &quot;&quot;&quot;</span>\n            <span class=\"n\">output_dir</span>
    <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;output_dir&quot;</span><span class=\"p\">]))</span>\n            <span
    class=\"k\">for</span> <span class=\"p\">(</span>\n                <span class=\"n\">article</span>\n
    \           <span class=\"p\">)</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">articles</span><span class=\"p\">:</span>
    \ <span class=\"c1\"># iter_articles(description=&quot;saving source documents&quot;):</span>\n
    \               <span class=\"k\">try</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">_save</span><span class=\"p\">(</span><span class=\"n\">output_dir</span><span
    class=\"p\">,</span> <span class=\"n\">article</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">except</span> <span class=\"n\">RepresenterError</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">_article</span> <span
    class=\"o\">=</span> <span class=\"n\">_strip_unserializable_values</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">article</span><span class=\"p\">)</span>\n\n                    <span
    class=\"n\">_save</span><span class=\"p\">(</span><span class=\"n\">output_dir</span><span
    class=\"p\">,</span> <span class=\"n\">_article</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>    </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Publish_Source.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Writes the final modified markdown and
    frontmatter to the output directory. The only configuration for the publish !
    ! ! ! ???+ source \" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Publish_Source.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Writes the final modified markdown and
    frontmatter to the output directory. The only configuration for the publish !
    ! ! ! ???+ source \" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    </head>\n    <body>\n<article style=\"text-align:
    center;\">\n    <style>\n        section {\n            font-size: 200%;\n        }\n\n\n
    \       .edit {\n            display: none;\n        }\n    </style>\n<section
    class=\"title\">\n    <h1 id=\"title\">\n        Publish_Source.Py\n    </h1>\n</section></article>\n
    \    </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Publish_Source.Py\n    </h1>\n</section>    <section class=\"body\">\n
    \       <p>Writes the final modified markdown and frontmatter to the output directory.\nReplacing
    the trailing slash if its there and adding .md will bring up the raw\nsource.</p>\n<h2
    id=\"configuration\">Configuration <a class=\"header-anchor\" href=\"#configuration\"><svg
    class=\"heading-permalink\" aria-hidden=\"true\" fill=\"currentColor\" focusable=\"false\"
    height=\"1em\" viewBox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>The only configuration
    for the publish_source plugin is to make sure its in\nyour list of hooks.</p>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n\n<span
    class=\"c1\"># make sure its in your list of hooks</span>\n<span class=\"n\">hooks</span><span
    class=\"o\">=</span><span class=\"p\">[</span>\n<span class=\"w\">   </span><span
    class=\"s2\">&quot;markata.plugins.publish_source&quot;</span><span class=\"p\">,</span>\n<span
    class=\"w\">   </span><span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<div
    class=\"admonition note\">\n<p class=\"admonition-title\">Note</p>\n<p>publish_source
    is included by default, but if you have not included the\ndefault set of hooks
    you will need to explicitly add it.</p>\n</div>\n<p>!! function <h2 id='_save'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_save <em class='small'>function</em></h2>\nsaves
    the article to the output directory at its specified slug.</p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_save
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
    <span class=\"nf\">_save</span><span class=\"p\">(</span><span class=\"n\">output_dir</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">,</span> <span
    class=\"n\">article</span><span class=\"p\">:</span> <span class=\"n\">frontmatter</span><span
    class=\"o\">.</span><span class=\"n\">Post</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            saves the article to the output directory at its specified
    slug.</span>\n<span class=\"sd\">            &quot;&quot;&quot;</span>\n            <span
    class=\"n\">path</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span>\n                <span class=\"n\">output_dir</span> <span
    class=\"o\">/</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">article</span><span class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">])</span><span class=\"o\">.</span><span class=\"n\">parent</span>
    <span class=\"o\">/</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">article</span><span class=\"p\">[</span><span class=\"s2\">&quot;path&quot;</span><span
    class=\"p\">])</span><span class=\"o\">.</span><span class=\"n\">name</span><span
    class=\"p\">,</span>\n            <span class=\"p\">)</span>\n            <span
    class=\"n\">path</span><span class=\"o\">.</span><span class=\"n\">parent</span><span
    class=\"o\">.</span><span class=\"n\">mkdir</span><span class=\"p\">(</span><span
    class=\"n\">parents</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span> <span class=\"n\">exist_ok</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n            <span class=\"n\">path</span><span
    class=\"o\">.</span><span class=\"n\">write_text</span><span class=\"p\">(</span><span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">dumps</span><span
    class=\"p\">())</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='_strip_unserializable_values'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_strip_unserializable_values
    <em class='small'>function</em></h2>\nReturns an article with only yaml serializable
    frontmatter.</p>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">_strip_unserializable_values <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">_strip_unserializable_values</span><span class=\"p\">(</span>\n
    \           <span class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span>\n            <span class=\"n\">article</span><span class=\"p\">:</span>
    <span class=\"n\">frontmatter</span><span class=\"o\">.</span><span class=\"n\">Post</span><span
    class=\"p\">,</span>\n        <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">frontmatter</span><span class=\"o\">.</span><span class=\"n\">Post</span><span
    class=\"p\">:</span>\n<span class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            Returns an article with only yaml serializable frontmatter.</span>\n<span
    class=\"sd\">            &quot;&quot;&quot;</span>\n            <span class=\"n\">_article</span>
    <span class=\"o\">=</span> <span class=\"n\">frontmatter</span><span class=\"o\">.</span><span
    class=\"n\">Post</span><span class=\"p\">(</span>\n                <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">content</span><span class=\"p\">,</span>\n
    \               <span class=\"o\">**</span><span class=\"p\">{</span><span class=\"n\">k</span><span
    class=\"p\">:</span> <span class=\"n\">v</span> <span class=\"k\">for</span> <span
    class=\"n\">k</span><span class=\"p\">,</span> <span class=\"n\">v</span> <span
    class=\"ow\">in</span> <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">metadata</span><span class=\"o\">.</span><span class=\"n\">items</span><span
    class=\"p\">()</span> <span class=\"k\">if</span> <span class=\"n\">k</span> <span
    class=\"o\">!=</span> <span class=\"s2\">&quot;content&quot;</span><span class=\"p\">},</span>\n
    \           <span class=\"p\">)</span>\n            <span class=\"n\">kwargs</span>
    <span class=\"o\">=</span> <span class=\"p\">{</span>\n                <span class=\"s2\">&quot;Dumper&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">yaml</span><span class=\"o\">.</span><span
    class=\"n\">cyaml</span><span class=\"o\">.</span><span class=\"n\">CSafeDumper</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;default_flow_style&quot;</span><span
    class=\"p\">:</span> <span class=\"kc\">False</span><span class=\"p\">,</span>\n
    \               <span class=\"s2\">&quot;allow_unicode&quot;</span><span class=\"p\">:</span>
    <span class=\"kc\">True</span><span class=\"p\">,</span>\n            <span class=\"p\">}</span>\n
    \           <span class=\"k\">for</span> <span class=\"n\">key</span><span class=\"p\">,</span>
    <span class=\"n\">value</span> <span class=\"ow\">in</span> <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">metadata</span><span class=\"o\">.</span><span
    class=\"n\">items</span><span class=\"p\">():</span>\n                <span class=\"k\">try</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">yaml</span><span class=\"o\">.</span><span
    class=\"n\">dump</span><span class=\"p\">({</span><span class=\"n\">key</span><span
    class=\"p\">:</span> <span class=\"n\">value</span><span class=\"p\">},</span>
    <span class=\"o\">**</span><span class=\"n\">kwargs</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">except</span> <span class=\"n\">RepresenterError</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">del</span> <span class=\"n\">_article</span><span
    class=\"p\">[</span><span class=\"n\">key</span><span class=\"p\">]</span>\n            <span
    class=\"k\">if</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">Post</span><span class=\"p\">:</span>\n                <span class=\"n\">_article</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">Post</span><span class=\"p\">(</span><span class=\"o\">**</span><span
    class=\"n\">_article</span><span class=\"o\">.</span><span class=\"n\">metadata</span><span
    class=\"p\">,</span> <span class=\"n\">path</span><span class=\"o\">=</span><span
    class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">path</span><span class=\"p\">))</span>\n
    \           <span class=\"k\">return</span> <span class=\"n\">_article</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='save' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>save <em class='small'>function</em></h2>\nSaves the final modified post
    to the output site as markdown.</p>\n<pre><code>!!! note\n    Any keys that are
    not yaml serializable will be stripped.\n</code></pre>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">save
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
    <span class=\"nf\">save</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            Saves the final modified post to the output site as markdown.</span>\n\n<span
    class=\"sd\">            !!! note</span>\n<span class=\"sd\">                Any
    keys that are not yaml serializable will be stripped.</span>\n\n<span class=\"sd\">
    \           &quot;&quot;&quot;</span>\n            <span class=\"n\">output_dir</span>
    <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;output_dir&quot;</span><span class=\"p\">]))</span>\n            <span
    class=\"k\">for</span> <span class=\"p\">(</span>\n                <span class=\"n\">article</span>\n
    \           <span class=\"p\">)</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">articles</span><span class=\"p\">:</span>
    \ <span class=\"c1\"># iter_articles(description=&quot;saving source documents&quot;):</span>\n
    \               <span class=\"k\">try</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">_save</span><span class=\"p\">(</span><span class=\"n\">output_dir</span><span
    class=\"p\">,</span> <span class=\"n\">article</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">except</span> <span class=\"n\">RepresenterError</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">_article</span> <span
    class=\"o\">=</span> <span class=\"n\">_strip_unserializable_values</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">article</span><span class=\"p\">)</span>\n\n                    <span
    class=\"n\">_save</span><span class=\"p\">(</span><span class=\"n\">output_dir</span><span
    class=\"p\">,</span> <span class=\"n\">_article</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>"
  raw.md: ''
published: true
slug: markata/plugins/publish-source
title: Publish_Source.Py


---

Writes the final modified markdown and frontmatter to the output directory.
Replacing the trailing slash if its there and adding .md will bring up the raw
source.

## Configuration

The only configuration for the publish_source plugin is to make sure its in
your list of hooks.


``` toml
[markata]

# make sure its in your list of hooks
hooks=[
   "markata.plugins.publish_source",
   ]
```

!!! note
    publish_source is included by default, but if you have not included the
    default set of hooks you will need to explicitly add it.


!! function <h2 id='_save' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_save <em class='small'>function</em></h2>
    saves the article to the output directory at its specified slug.
???+ source "_save <em class='small'>source</em>"

```python

        def _save(output_dir: Path, article: frontmatter.Post) -> None:
            """
            saves the article to the output directory at its specified slug.
            """
            path = Path(
                output_dir / Path(article["slug"]).parent / Path(article["path"]).name,
            )
            path.parent.mkdir(parents=True, exist_ok=True)
            path.write_text(article.dumps())
```


!! function <h2 id='_strip_unserializable_values' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_strip_unserializable_values <em class='small'>function</em></h2>
    Returns an article with only yaml serializable frontmatter.
???+ source "_strip_unserializable_values <em class='small'>source</em>"

```python

        def _strip_unserializable_values(
            markata: "Markata",
            article: frontmatter.Post,
        ) -> frontmatter.Post:
            """
            Returns an article with only yaml serializable frontmatter.
            """
            _article = frontmatter.Post(
                article.content,
                **{k: v for k, v in article.metadata.items() if k != "content"},
            )
            kwargs = {
                "Dumper": yaml.cyaml.CSafeDumper,
                "default_flow_style": False,
                "allow_unicode": True,
            }
            for key, value in article.metadata.items():
                try:
                    yaml.dump({key: value}, **kwargs)
                except RepresenterError:
                    del _article[key]
            if markata.Post:
                _article = markata.Post(**_article.metadata, path=str(article.path))
            return _article
```


!! function <h2 id='save' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>save <em class='small'>function</em></h2>
    Saves the final modified post to the output site as markdown.

    !!! note
        Any keys that are not yaml serializable will be stripped.
???+ source "save <em class='small'>source</em>"

```python

        def save(markata: "Markata") -> None:
            """
            Saves the final modified post to the output site as markdown.

            !!! note
                Any keys that are not yaml serializable will be stripped.

            """
            output_dir = Path(str(markata.config["output_dir"]))
            for (
                article
            ) in markata.articles:  # iter_articles(description="saving source documents"):
                try:
                    _save(output_dir, article)
                except RepresenterError:
                    _article = _strip_unserializable_values(markata, article)

                    _save(output_dir, _article)
```

