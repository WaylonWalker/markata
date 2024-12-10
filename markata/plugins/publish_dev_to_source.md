---
content: "None\n\n\n!! function <h2 id='should_join' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>should_join <em class='small'>function</em></h2>\n\n???+ source \"should_join
  <em class='small'>source</em>\"\n\n```python\n\n        def should_join(line):\n
  \           if line == \"\":\n                return False\n            if line
  is None:\n                return False\n            if line[0].isalpha():\n                return
  True\n            if line[0].startswith(\"[\"):\n                return True\n            if
  line[0].startswith(\"!\"):\n                return True\n            return False\n```\n\n\n!!
  function <h2 id='join_lines' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>join_lines <em class='small'>function</em></h2>\n\n???+ source \"join_lines
  <em class='small'>source</em>\"\n\n```python\n\n        def join_lines(article):\n
  \           lines = article.split(\"\\n\")\n            line_number = 0\n            while
  line_number + 1 < len(lines):\n                line = lines[line_number]\n                nextline
  = lines[line_number + 1]\n                if should_join(line) and should_join(nextline):\n
  \                   lines[line_number] = f\"{line} {nextline}\"\n                    lines.pop(line_number
  + 1)\n                else:\n                    line_number += 1\n\n            return
  \"\\n\".join(lines)\n```\n\n\n!! class <h2 id='PublishDevToSourcePost' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>PublishDevToSourcePost <em class='small'>class</em></h2>\n\n???+
  source \"PublishDevToSourcePost <em class='small'>source</em>\"\n\n```python\n\n
  \       class PublishDevToSourcePost(BaseModel):\n            markata: Any = Field(None,
  exclude=True)\n            canonical_url: Optional[str] = None\n```\n\n\n!! function
  <h2 id='post_model' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>post_model
  <em class='small'>function</em></h2>\n\n???+ source \"post_model <em class='small'>source</em>\"\n\n```python\n\n
  \       def post_model(markata: \"Markata\") -> None:\n            markata.post_models.append(PublishDevToSourcePost)\n```\n\n\n!!
  function <h2 id='post_render' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>post_render <em class='small'>function</em></h2>\n\n???+ source \"post_render
  <em class='small'>source</em>\"\n\n```python\n\n        def post_render(markata:
  \"Markata\") -> None:\n            for post in markata.iter_articles(description=\"saving
  source documents\"):\n                article = frontmatter.Post(\n                    post.content,\n
  \                   **{k: v for k, v in post.metadata.items() if k in DEV_TO_FRONTMATTER},\n
  \               )\n\n                article.content = join_lines(article.content)\n
  \               article.content = join_lines(article.content)\n\n                if
  \"canonical_url\" not in article:\n                    article[\"canonical_url\"]
  = f\"{markata.config.url}/{post.slug}/\"\n\n                if \"published\" not
  in article:\n                    article[\"published\"] = True\n\n                if
  \"cover_image\" not in article:\n                    article[\"cover_image\"] =
  f\"{markata.config.images_url}/{post.slug}.png\"\n                post.dev_to =
  article\n```\n\n\n!! function <h2 id='save' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>save <em class='small'>function</em></h2>\n\n???+ source \"save <em class='small'>source</em>\"\n\n```python\n\n
  \       def save(markata: \"Markata\") -> None:\n            output_dir = Path(str(markata.config.output_dir))\n
  \           with markata.console.status(\"Saving source documents...\") as status:\n
  \               for post in markata.iter_articles(description=\"saving source documents\"):\n
  \                   status.update(f\"Saving {post['slug']}...\")\n                    with
  open(output_dir / Path(post[\"slug\"]) / \"dev.md\", \"w+\") as f:\n                        f.write(frontmatter.dumps(post.dev_to))\n```\n\n"
date: 0001-01-01
description: 'None ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
  source  ! ???+ source '
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Publish_Dev_To_Source.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None ! ???+ source  ! ???+ source  !
    ???+ source  ! ???+ source  ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Publish_Dev_To_Source.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None ! ???+ source  ! ???+ source  !
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Publish_Dev_To_Source.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <p>None</p>\n<p>!!
    function <h2 id='should_join' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>should_join <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">should_join
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
    <span class=\"nf\">should_join</span><span class=\"p\">(</span><span class=\"n\">line</span><span
    class=\"p\">):</span>\n            <span class=\"k\">if</span> <span class=\"n\">line</span>
    <span class=\"o\">==</span> <span class=\"s2\">&quot;&quot;</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">return</span> <span class=\"kc\">False</span>\n
    \           <span class=\"k\">if</span> <span class=\"n\">line</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
    class=\"k\">return</span> <span class=\"kc\">False</span>\n            <span class=\"k\">if</span>
    <span class=\"n\">line</span><span class=\"p\">[</span><span class=\"mi\">0</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">isalpha</span><span
    class=\"p\">():</span>\n                <span class=\"k\">return</span> <span
    class=\"kc\">True</span>\n            <span class=\"k\">if</span> <span class=\"n\">line</span><span
    class=\"p\">[</span><span class=\"mi\">0</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">startswith</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;[&quot;</span><span class=\"p\">):</span>\n                <span
    class=\"k\">return</span> <span class=\"kc\">True</span>\n            <span class=\"k\">if</span>
    <span class=\"n\">line</span><span class=\"p\">[</span><span class=\"mi\">0</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">startswith</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;!&quot;</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">return</span> <span class=\"kc\">True</span>\n
    \           <span class=\"k\">return</span> <span class=\"kc\">False</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='join_lines' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>join_lines <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">join_lines
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
    <span class=\"nf\">join_lines</span><span class=\"p\">(</span><span class=\"n\">article</span><span
    class=\"p\">):</span>\n            <span class=\"n\">lines</span> <span class=\"o\">=</span>
    <span class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">split</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span
    class=\"s2\">&quot;</span><span class=\"p\">)</span>\n            <span class=\"n\">line_number</span>
    <span class=\"o\">=</span> <span class=\"mi\">0</span>\n            <span class=\"k\">while</span>
    <span class=\"n\">line_number</span> <span class=\"o\">+</span> <span class=\"mi\">1</span>
    <span class=\"o\">&lt;</span> <span class=\"nb\">len</span><span class=\"p\">(</span><span
    class=\"n\">lines</span><span class=\"p\">):</span>\n                <span class=\"n\">line</span>
    <span class=\"o\">=</span> <span class=\"n\">lines</span><span class=\"p\">[</span><span
    class=\"n\">line_number</span><span class=\"p\">]</span>\n                <span
    class=\"n\">nextline</span> <span class=\"o\">=</span> <span class=\"n\">lines</span><span
    class=\"p\">[</span><span class=\"n\">line_number</span> <span class=\"o\">+</span>
    <span class=\"mi\">1</span><span class=\"p\">]</span>\n                <span class=\"k\">if</span>
    <span class=\"n\">should_join</span><span class=\"p\">(</span><span class=\"n\">line</span><span
    class=\"p\">)</span> <span class=\"ow\">and</span> <span class=\"n\">should_join</span><span
    class=\"p\">(</span><span class=\"n\">nextline</span><span class=\"p\">):</span>\n
    \                   <span class=\"n\">lines</span><span class=\"p\">[</span><span
    class=\"n\">line_number</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">line</span><span class=\"si\">}</span><span class=\"s2\"> </span><span
    class=\"si\">{</span><span class=\"n\">nextline</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span>\n                    <span class=\"n\">lines</span><span
    class=\"o\">.</span><span class=\"n\">pop</span><span class=\"p\">(</span><span
    class=\"n\">line_number</span> <span class=\"o\">+</span> <span class=\"mi\">1</span><span
    class=\"p\">)</span>\n                <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">line_number</span> <span class=\"o\">+=</span>
    <span class=\"mi\">1</span>\n\n            <span class=\"k\">return</span> <span
    class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span><span
    class=\"o\">.</span><span class=\"n\">join</span><span class=\"p\">(</span><span
    class=\"n\">lines</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='PublishDevToSourcePost' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>PublishDevToSourcePost <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">PublishDevToSourcePost
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
    <span class=\"nc\">PublishDevToSourcePost</span><span class=\"p\">(</span><span
    class=\"n\">BaseModel</span><span class=\"p\">):</span>\n            <span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"n\">Any</span> <span class=\"o\">=</span> <span
    class=\"n\">Field</span><span class=\"p\">(</span><span class=\"kc\">None</span><span
    class=\"p\">,</span> <span class=\"n\">exclude</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n            <span class=\"n\">canonical_url</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"kc\">None</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='post_model'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>post_model <em class='small'>function</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">post_model
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
    <span class=\"nf\">post_model</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">post_models</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">PublishDevToSourcePost</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='post_render' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>post_render <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">post_render
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
    <span class=\"nf\">post_render</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">for</span> <span class=\"n\">post</span> <span class=\"ow\">in</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">iter_articles</span><span
    class=\"p\">(</span><span class=\"n\">description</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;saving source documents&quot;</span><span class=\"p\">):</span>\n
    \               <span class=\"n\">article</span> <span class=\"o\">=</span> <span
    class=\"n\">frontmatter</span><span class=\"o\">.</span><span class=\"n\">Post</span><span
    class=\"p\">(</span>\n                    <span class=\"n\">post</span><span class=\"o\">.</span><span
    class=\"n\">content</span><span class=\"p\">,</span>\n                    <span
    class=\"o\">**</span><span class=\"p\">{</span><span class=\"n\">k</span><span
    class=\"p\">:</span> <span class=\"n\">v</span> <span class=\"k\">for</span> <span
    class=\"n\">k</span><span class=\"p\">,</span> <span class=\"n\">v</span> <span
    class=\"ow\">in</span> <span class=\"n\">post</span><span class=\"o\">.</span><span
    class=\"n\">metadata</span><span class=\"o\">.</span><span class=\"n\">items</span><span
    class=\"p\">()</span> <span class=\"k\">if</span> <span class=\"n\">k</span> <span
    class=\"ow\">in</span> <span class=\"n\">DEV_TO_FRONTMATTER</span><span class=\"p\">},</span>\n
    \               <span class=\"p\">)</span>\n\n                <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">content</span> <span class=\"o\">=</span>
    <span class=\"n\">join_lines</span><span class=\"p\">(</span><span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">content</span><span class=\"p\">)</span>\n
    \               <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">content</span> <span class=\"o\">=</span> <span class=\"n\">join_lines</span><span
    class=\"p\">(</span><span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">content</span><span class=\"p\">)</span>\n\n                <span
    class=\"k\">if</span> <span class=\"s2\">&quot;canonical_url&quot;</span> <span
    class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">article</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">article</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;canonical_url&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">url</span><span
    class=\"si\">}</span><span class=\"s2\">/</span><span class=\"si\">{</span><span
    class=\"n\">post</span><span class=\"o\">.</span><span class=\"n\">slug</span><span
    class=\"si\">}</span><span class=\"s2\">/&quot;</span>\n\n                <span
    class=\"k\">if</span> <span class=\"s2\">&quot;published&quot;</span> <span class=\"ow\">not</span>
    <span class=\"ow\">in</span> <span class=\"n\">article</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">article</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;published&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"kc\">True</span>\n\n                <span class=\"k\">if</span>
    <span class=\"s2\">&quot;cover_image&quot;</span> <span class=\"ow\">not</span>
    <span class=\"ow\">in</span> <span class=\"n\">article</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">article</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;cover_image&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">images_url</span><span class=\"si\">}</span><span
    class=\"s2\">/</span><span class=\"si\">{</span><span class=\"n\">post</span><span
    class=\"o\">.</span><span class=\"n\">slug</span><span class=\"si\">}</span><span
    class=\"s2\">.png&quot;</span>\n                <span class=\"n\">post</span><span
    class=\"o\">.</span><span class=\"n\">dev_to</span> <span class=\"o\">=</span>
    <span class=\"n\">article</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2
    id='save' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>save <em
    class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">save <em class='small'>source</em></p>\n</div>\n<pre
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
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">output_dir</span> <span class=\"o\">=</span> <span
    class=\"n\">Path</span><span class=\"p\">(</span><span class=\"nb\">str</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">output_dir</span><span
    class=\"p\">))</span>\n            <span class=\"k\">with</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">status</span><span class=\"p\">(</span><span class=\"s2\">&quot;Saving
    source documents...&quot;</span><span class=\"p\">)</span> <span class=\"k\">as</span>
    <span class=\"n\">status</span><span class=\"p\">:</span>\n                <span
    class=\"k\">for</span> <span class=\"n\">post</span> <span class=\"ow\">in</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">iter_articles</span><span
    class=\"p\">(</span><span class=\"n\">description</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;saving source documents&quot;</span><span class=\"p\">):</span>\n
    \                   <span class=\"n\">status</span><span class=\"o\">.</span><span
    class=\"n\">update</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;Saving </span><span class=\"si\">{</span><span class=\"n\">post</span><span
    class=\"p\">[</span><span class=\"s1\">&#39;slug&#39;</span><span class=\"p\">]</span><span
    class=\"si\">}</span><span class=\"s2\">...&quot;</span><span class=\"p\">)</span>\n
    \                   <span class=\"k\">with</span> <span class=\"nb\">open</span><span
    class=\"p\">(</span><span class=\"n\">output_dir</span> <span class=\"o\">/</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">post</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">])</span>
    <span class=\"o\">/</span> <span class=\"s2\">&quot;dev.md&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;w+&quot;</span><span class=\"p\">)</span> <span class=\"k\">as</span>
    <span class=\"n\">f</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">f</span><span class=\"o\">.</span><span class=\"n\">write</span><span
    class=\"p\">(</span><span class=\"n\">frontmatter</span><span class=\"o\">.</span><span
    class=\"n\">dumps</span><span class=\"p\">(</span><span class=\"n\">post</span><span
    class=\"o\">.</span><span class=\"n\">dev_to</span><span class=\"p\">))</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>    </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Publish_Dev_To_Source.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None ! ???+ source  ! ???+ source  !
    ???+ source  ! ???+ source  ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Publish_Dev_To_Source.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None ! ???+ source  ! ???+ source  !
    ???+ source  ! ???+ source  ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<article style=\"text-align: center;\">\n    <style>\n        section
    {\n            font-size: 200%;\n        }\n\n\n        .edit {\n            display:
    none;\n        }\n    </style>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Publish_Dev_To_Source.Py\n    </h1>\n</section></article>\n     </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Publish_Dev_To_Source.Py\n    </h1>\n</section>    <section class=\"body\">\n
    \       <p>None</p>\n<p>!! function <h2 id='should_join' class='admonition-title'
    style='margin:0;padding:.5rem 1rem;'>should_join <em class='small'>function</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">should_join
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
    <span class=\"nf\">should_join</span><span class=\"p\">(</span><span class=\"n\">line</span><span
    class=\"p\">):</span>\n            <span class=\"k\">if</span> <span class=\"n\">line</span>
    <span class=\"o\">==</span> <span class=\"s2\">&quot;&quot;</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">return</span> <span class=\"kc\">False</span>\n
    \           <span class=\"k\">if</span> <span class=\"n\">line</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
    class=\"k\">return</span> <span class=\"kc\">False</span>\n            <span class=\"k\">if</span>
    <span class=\"n\">line</span><span class=\"p\">[</span><span class=\"mi\">0</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">isalpha</span><span
    class=\"p\">():</span>\n                <span class=\"k\">return</span> <span
    class=\"kc\">True</span>\n            <span class=\"k\">if</span> <span class=\"n\">line</span><span
    class=\"p\">[</span><span class=\"mi\">0</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">startswith</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;[&quot;</span><span class=\"p\">):</span>\n                <span
    class=\"k\">return</span> <span class=\"kc\">True</span>\n            <span class=\"k\">if</span>
    <span class=\"n\">line</span><span class=\"p\">[</span><span class=\"mi\">0</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">startswith</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;!&quot;</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">return</span> <span class=\"kc\">True</span>\n
    \           <span class=\"k\">return</span> <span class=\"kc\">False</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='join_lines' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>join_lines <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">join_lines
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
    <span class=\"nf\">join_lines</span><span class=\"p\">(</span><span class=\"n\">article</span><span
    class=\"p\">):</span>\n            <span class=\"n\">lines</span> <span class=\"o\">=</span>
    <span class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">split</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span
    class=\"s2\">&quot;</span><span class=\"p\">)</span>\n            <span class=\"n\">line_number</span>
    <span class=\"o\">=</span> <span class=\"mi\">0</span>\n            <span class=\"k\">while</span>
    <span class=\"n\">line_number</span> <span class=\"o\">+</span> <span class=\"mi\">1</span>
    <span class=\"o\">&lt;</span> <span class=\"nb\">len</span><span class=\"p\">(</span><span
    class=\"n\">lines</span><span class=\"p\">):</span>\n                <span class=\"n\">line</span>
    <span class=\"o\">=</span> <span class=\"n\">lines</span><span class=\"p\">[</span><span
    class=\"n\">line_number</span><span class=\"p\">]</span>\n                <span
    class=\"n\">nextline</span> <span class=\"o\">=</span> <span class=\"n\">lines</span><span
    class=\"p\">[</span><span class=\"n\">line_number</span> <span class=\"o\">+</span>
    <span class=\"mi\">1</span><span class=\"p\">]</span>\n                <span class=\"k\">if</span>
    <span class=\"n\">should_join</span><span class=\"p\">(</span><span class=\"n\">line</span><span
    class=\"p\">)</span> <span class=\"ow\">and</span> <span class=\"n\">should_join</span><span
    class=\"p\">(</span><span class=\"n\">nextline</span><span class=\"p\">):</span>\n
    \                   <span class=\"n\">lines</span><span class=\"p\">[</span><span
    class=\"n\">line_number</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">line</span><span class=\"si\">}</span><span class=\"s2\"> </span><span
    class=\"si\">{</span><span class=\"n\">nextline</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span>\n                    <span class=\"n\">lines</span><span
    class=\"o\">.</span><span class=\"n\">pop</span><span class=\"p\">(</span><span
    class=\"n\">line_number</span> <span class=\"o\">+</span> <span class=\"mi\">1</span><span
    class=\"p\">)</span>\n                <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">line_number</span> <span class=\"o\">+=</span>
    <span class=\"mi\">1</span>\n\n            <span class=\"k\">return</span> <span
    class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span><span
    class=\"o\">.</span><span class=\"n\">join</span><span class=\"p\">(</span><span
    class=\"n\">lines</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='PublishDevToSourcePost' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>PublishDevToSourcePost <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">PublishDevToSourcePost
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
    <span class=\"nc\">PublishDevToSourcePost</span><span class=\"p\">(</span><span
    class=\"n\">BaseModel</span><span class=\"p\">):</span>\n            <span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"n\">Any</span> <span class=\"o\">=</span> <span
    class=\"n\">Field</span><span class=\"p\">(</span><span class=\"kc\">None</span><span
    class=\"p\">,</span> <span class=\"n\">exclude</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n            <span class=\"n\">canonical_url</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"kc\">None</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='post_model'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>post_model <em class='small'>function</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">post_model
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
    <span class=\"nf\">post_model</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">post_models</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">PublishDevToSourcePost</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='post_render' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>post_render <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">post_render
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
    <span class=\"nf\">post_render</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">for</span> <span class=\"n\">post</span> <span class=\"ow\">in</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">iter_articles</span><span
    class=\"p\">(</span><span class=\"n\">description</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;saving source documents&quot;</span><span class=\"p\">):</span>\n
    \               <span class=\"n\">article</span> <span class=\"o\">=</span> <span
    class=\"n\">frontmatter</span><span class=\"o\">.</span><span class=\"n\">Post</span><span
    class=\"p\">(</span>\n                    <span class=\"n\">post</span><span class=\"o\">.</span><span
    class=\"n\">content</span><span class=\"p\">,</span>\n                    <span
    class=\"o\">**</span><span class=\"p\">{</span><span class=\"n\">k</span><span
    class=\"p\">:</span> <span class=\"n\">v</span> <span class=\"k\">for</span> <span
    class=\"n\">k</span><span class=\"p\">,</span> <span class=\"n\">v</span> <span
    class=\"ow\">in</span> <span class=\"n\">post</span><span class=\"o\">.</span><span
    class=\"n\">metadata</span><span class=\"o\">.</span><span class=\"n\">items</span><span
    class=\"p\">()</span> <span class=\"k\">if</span> <span class=\"n\">k</span> <span
    class=\"ow\">in</span> <span class=\"n\">DEV_TO_FRONTMATTER</span><span class=\"p\">},</span>\n
    \               <span class=\"p\">)</span>\n\n                <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">content</span> <span class=\"o\">=</span>
    <span class=\"n\">join_lines</span><span class=\"p\">(</span><span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">content</span><span class=\"p\">)</span>\n
    \               <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">content</span> <span class=\"o\">=</span> <span class=\"n\">join_lines</span><span
    class=\"p\">(</span><span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">content</span><span class=\"p\">)</span>\n\n                <span
    class=\"k\">if</span> <span class=\"s2\">&quot;canonical_url&quot;</span> <span
    class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">article</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">article</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;canonical_url&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">url</span><span
    class=\"si\">}</span><span class=\"s2\">/</span><span class=\"si\">{</span><span
    class=\"n\">post</span><span class=\"o\">.</span><span class=\"n\">slug</span><span
    class=\"si\">}</span><span class=\"s2\">/&quot;</span>\n\n                <span
    class=\"k\">if</span> <span class=\"s2\">&quot;published&quot;</span> <span class=\"ow\">not</span>
    <span class=\"ow\">in</span> <span class=\"n\">article</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">article</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;published&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"kc\">True</span>\n\n                <span class=\"k\">if</span>
    <span class=\"s2\">&quot;cover_image&quot;</span> <span class=\"ow\">not</span>
    <span class=\"ow\">in</span> <span class=\"n\">article</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">article</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;cover_image&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">images_url</span><span class=\"si\">}</span><span
    class=\"s2\">/</span><span class=\"si\">{</span><span class=\"n\">post</span><span
    class=\"o\">.</span><span class=\"n\">slug</span><span class=\"si\">}</span><span
    class=\"s2\">.png&quot;</span>\n                <span class=\"n\">post</span><span
    class=\"o\">.</span><span class=\"n\">dev_to</span> <span class=\"o\">=</span>
    <span class=\"n\">article</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2
    id='save' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>save <em
    class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">save <em class='small'>source</em></p>\n</div>\n<pre
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
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">output_dir</span> <span class=\"o\">=</span> <span
    class=\"n\">Path</span><span class=\"p\">(</span><span class=\"nb\">str</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">output_dir</span><span
    class=\"p\">))</span>\n            <span class=\"k\">with</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">status</span><span class=\"p\">(</span><span class=\"s2\">&quot;Saving
    source documents...&quot;</span><span class=\"p\">)</span> <span class=\"k\">as</span>
    <span class=\"n\">status</span><span class=\"p\">:</span>\n                <span
    class=\"k\">for</span> <span class=\"n\">post</span> <span class=\"ow\">in</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">iter_articles</span><span
    class=\"p\">(</span><span class=\"n\">description</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;saving source documents&quot;</span><span class=\"p\">):</span>\n
    \                   <span class=\"n\">status</span><span class=\"o\">.</span><span
    class=\"n\">update</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;Saving </span><span class=\"si\">{</span><span class=\"n\">post</span><span
    class=\"p\">[</span><span class=\"s1\">&#39;slug&#39;</span><span class=\"p\">]</span><span
    class=\"si\">}</span><span class=\"s2\">...&quot;</span><span class=\"p\">)</span>\n
    \                   <span class=\"k\">with</span> <span class=\"nb\">open</span><span
    class=\"p\">(</span><span class=\"n\">output_dir</span> <span class=\"o\">/</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">post</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">])</span>
    <span class=\"o\">/</span> <span class=\"s2\">&quot;dev.md&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;w+&quot;</span><span class=\"p\">)</span> <span class=\"k\">as</span>
    <span class=\"n\">f</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">f</span><span class=\"o\">.</span><span class=\"n\">write</span><span
    class=\"p\">(</span><span class=\"n\">frontmatter</span><span class=\"o\">.</span><span
    class=\"n\">dumps</span><span class=\"p\">(</span><span class=\"n\">post</span><span
    class=\"o\">.</span><span class=\"n\">dev_to</span><span class=\"p\">))</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>"
  raw.md: ''
published: true
slug: markata/plugins/publish-dev-to-source
title: Publish_Dev_To_Source.Py


---

None


!! function <h2 id='should_join' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>should_join <em class='small'>function</em></h2>

???+ source "should_join <em class='small'>source</em>"

```python

        def should_join(line):
            if line == "":
                return False
            if line is None:
                return False
            if line[0].isalpha():
                return True
            if line[0].startswith("["):
                return True
            if line[0].startswith("!"):
                return True
            return False
```


!! function <h2 id='join_lines' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>join_lines <em class='small'>function</em></h2>

???+ source "join_lines <em class='small'>source</em>"

```python

        def join_lines(article):
            lines = article.split("\n")
            line_number = 0
            while line_number + 1 < len(lines):
                line = lines[line_number]
                nextline = lines[line_number + 1]
                if should_join(line) and should_join(nextline):
                    lines[line_number] = f"{line} {nextline}"
                    lines.pop(line_number + 1)
                else:
                    line_number += 1

            return "\n".join(lines)
```


!! class <h2 id='PublishDevToSourcePost' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>PublishDevToSourcePost <em class='small'>class</em></h2>

???+ source "PublishDevToSourcePost <em class='small'>source</em>"

```python

        class PublishDevToSourcePost(BaseModel):
            markata: Any = Field(None, exclude=True)
            canonical_url: Optional[str] = None
```


!! function <h2 id='post_model' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>post_model <em class='small'>function</em></h2>

???+ source "post_model <em class='small'>source</em>"

```python

        def post_model(markata: "Markata") -> None:
            markata.post_models.append(PublishDevToSourcePost)
```


!! function <h2 id='post_render' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>post_render <em class='small'>function</em></h2>

???+ source "post_render <em class='small'>source</em>"

```python

        def post_render(markata: "Markata") -> None:
            for post in markata.iter_articles(description="saving source documents"):
                article = frontmatter.Post(
                    post.content,
                    **{k: v for k, v in post.metadata.items() if k in DEV_TO_FRONTMATTER},
                )

                article.content = join_lines(article.content)
                article.content = join_lines(article.content)

                if "canonical_url" not in article:
                    article["canonical_url"] = f"{markata.config.url}/{post.slug}/"

                if "published" not in article:
                    article["published"] = True

                if "cover_image" not in article:
                    article["cover_image"] = f"{markata.config.images_url}/{post.slug}.png"
                post.dev_to = article
```


!! function <h2 id='save' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>save <em class='small'>function</em></h2>

???+ source "save <em class='small'>source</em>"

```python

        def save(markata: "Markata") -> None:
            output_dir = Path(str(markata.config.output_dir))
            with markata.console.status("Saving source documents...") as status:
                for post in markata.iter_articles(description="saving source documents"):
                    status.update(f"Saving {post['slug']}...")
                    with open(output_dir / Path(post["slug"]) / "dev.md", "w+") as f:
                        f.write(frontmatter.dumps(post.dev_to))
```

