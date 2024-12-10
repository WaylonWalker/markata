---
content: "Wikilinks depicted with `\\[[wikilinks]]` can be enabled for sites using
  the\n`markdown-it-py` plugin markata will look through all posts matching up the\nfile
  stem to the wiki link and inserting the slug as the href.\n\n???+ note normal wikilink\n\n
  \   Here is a link to a markdown file `docs/nav.md`, and the url becomes\n    [/nav](/nav).\n\n
  \   ```md\n    [[nav]]\n    ```\n\n    When rendered out this wikilink will become
  an anchor link.\n\n    ```html\n    <a class=\"wikilink\" href=\"/nav\">load</a>\n
  \   ```\n\n    > This behaves just like the standard wikilink\n\nA special feature
  that markata brings is slug lookup.  It is able to not only\nblindly link to the
  route specified, but will look up the slug of an article.\n\n\n???+ note markata
  slug lookup\n    Markata has a load plugin that is generated with the [[docs]] plugin.
  \ It's\n    filepath is `markata/plugins/load.py`, so it can be referenced by the
  file\n    stem `load`.\n\n    ```md\n    [[load]]\n    ```\n\n    Markata will look
  up the article by the file stem, grab the first article,\n    and use its slug as
  the href.  This turns it into an anchor link that looks\n    like this.\n\n    ```html\n
  \   <a class=\"wikilink\" href=\"/markata/plugins/load\">load</a>\n    ```\n\n\n!!
  function <h2 id='wikilinks_plugin' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>wikilinks_plugin <em class='small'>function</em></h2>\n    A plugin to create
  wikilinks tokens.\n    These, token should be handled by the renderer.\n\n    ???+
  example\n\n        ```md title=markdown\n        [[nav]]\n        ```\n\n        ```html
  title=html\n        <a class=\"wikilink\" href=\"/nav\">load</a>\n        ```\n???+
  source \"wikilinks_plugin <em class='small'>source</em>\"\n\n```python\n\n        def
  wikilinks_plugin(\n            md: MarkdownIt,\n            start_delimiter: str
  = \"[\",\n            end_delimiter: str = \"]\",\n            markata=None,\n        ):\n
  \           \"\"\"A plugin to create wikilinks tokens.\n            These, token
  should be handled by the renderer.\n\n            ???+ example\n\n                ```md
  title=markdown\n                [[nav]]\n                ```\n\n                ```html
  title=html\n                <a class=\"wikilink\" href=\"/nav\">load</a>\n                ```\n
  \           \"\"\"\n\n            start_char = ord(start_delimiter)\n            end_char
  = ord(end_delimiter)\n\n            def _wikilinks_inline(state: StateInline, silent:
  bool):\n                try:\n                    if (\n                        state.srcCharCode[state.pos]
  != start_char\n                        or state.srcCharCode[state.pos + 1] != start_char\n
  \                   ):\n                        return False\n                except
  IndexError:\n                    return False\n\n                pos = state.pos
  + 2\n                found_closing = False\n                while True:\n                    try:\n
  \                       end = state.srcCharCode.index(end_char, pos)\n                    except
  ValueError:\n                        return False\n                    try:\n                        if
  state.srcCharCode[end + 1] == end_char:\n                            found_closing
  = True\n                            break\n                    except IndexError:\n
  \                       return False\n                    pos = end + 2\n\n                if
  not found_closing:\n                    return False\n\n                text = state.src[state.pos
  + 2 : end].strip()\n                state.pos = end + 2\n\n                if silent:\n
  \                   return True\n\n                token = state.push(\"link_open\",
  \"a\", 1)\n                token.block = False\n                token.attrSet(\"class\",
  \"wikilink\")\n                if \"#\" in text:\n                    link, id =
  text.split(\"#\")\n                    link = link.strip(\"/\")\n                else:\n
  \                   link, id = text, None\n                possible_pages = markata.filter(\n
  \                   f'str(path).split(\"/\")[-1].split(\".\")[0].replace(\"_\",
  \"-\") == \"{link.replace(\"_\", \"-\")}\"',\n                )\n                if
  len(possible_pages) == 1:\n                    link = possible_pages[0].get(\"slug\",
  f\"/{text}\")\n                elif len(possible_pages) > 1:\n                    logger.warning(\n
  \                       f\"wikilink [[{text}]] ({link}, {id}) has duplicate matches,
  defaulting to the first\",\n                    )\n                    link = possible_pages[0].get(\"slug\",
  f\"/{text}\")\n                else:\n                    logger.warning(\n                        f\"wikilink
  [[{text}]] ({link}, {id}) no matches, defaulting to '/{text}'\",\n                    )\n
  \                   link = text\n\n                if id and not link.endswith(f\"#{id}\"):\n
  \                   link = f\"{link}#{id}\"\n\n                token.attrSet(\"href\",
  f\"/{link}\")\n                content_token = state.push(\"text\", \"\", 0)\n                content_token.content
  = text\n\n                token = state.push(\"link_close\", \"a\", -1)\n                token.content
  = text\n\n                return True\n\n            md.inline.ruler.before(\"escape\",
  \"wikilinks_inline\", _wikilinks_inline)\n```\n\n\n!! function <h2 id='_wikilinks_inline'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_wikilinks_inline
  <em class='small'>function</em></h2>\n\n???+ source \"_wikilinks_inline <em class='small'>source</em>\"\n\n```python\n\n
  \       def _wikilinks_inline(state: StateInline, silent: bool):\n                try:\n
  \                   if (\n                        state.srcCharCode[state.pos] !=
  start_char\n                        or state.srcCharCode[state.pos + 1] != start_char\n
  \                   ):\n                        return False\n                except
  IndexError:\n                    return False\n\n                pos = state.pos
  + 2\n                found_closing = False\n                while True:\n                    try:\n
  \                       end = state.srcCharCode.index(end_char, pos)\n                    except
  ValueError:\n                        return False\n                    try:\n                        if
  state.srcCharCode[end + 1] == end_char:\n                            found_closing
  = True\n                            break\n                    except IndexError:\n
  \                       return False\n                    pos = end + 2\n\n                if
  not found_closing:\n                    return False\n\n                text = state.src[state.pos
  + 2 : end].strip()\n                state.pos = end + 2\n\n                if silent:\n
  \                   return True\n\n                token = state.push(\"link_open\",
  \"a\", 1)\n                token.block = False\n                token.attrSet(\"class\",
  \"wikilink\")\n                if \"#\" in text:\n                    link, id =
  text.split(\"#\")\n                    link = link.strip(\"/\")\n                else:\n
  \                   link, id = text, None\n                possible_pages = markata.filter(\n
  \                   f'str(path).split(\"/\")[-1].split(\".\")[0].replace(\"_\",
  \"-\") == \"{link.replace(\"_\", \"-\")}\"',\n                )\n                if
  len(possible_pages) == 1:\n                    link = possible_pages[0].get(\"slug\",
  f\"/{text}\")\n                elif len(possible_pages) > 1:\n                    logger.warning(\n
  \                       f\"wikilink [[{text}]] ({link}, {id}) has duplicate matches,
  defaulting to the first\",\n                    )\n                    link = possible_pages[0].get(\"slug\",
  f\"/{text}\")\n                else:\n                    logger.warning(\n                        f\"wikilink
  [[{text}]] ({link}, {id}) no matches, defaulting to '/{text}'\",\n                    )\n
  \                   link = text\n\n                if id and not link.endswith(f\"#{id}\"):\n
  \                   link = f\"{link}#{id}\"\n\n                token.attrSet(\"href\",
  f\"/{link}\")\n                content_token = state.push(\"text\", \"\", 0)\n                content_token.content
  = text\n\n                token = state.push(\"link_close\", \"a\", -1)\n                token.content
  = text\n\n                return True\n```\n\n"
date: 0001-01-01
description: Wikilinks depicted with  ???+ note normal wikilink A special feature
  that markata brings is slug lookup.  It is able to not only ???+ note markata slug
  lookup !
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Md_It_Wikilinks.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Wikilinks depicted with  ???+ note normal
    wikilink A special feature that markata brings is slug lookup.  It is able to
    not only ???+ note markata slug lookup !\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Md_It_Wikilinks.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Wikilinks depicted with  ???+ note normal
    wikilink A special feature that markata brings is slug lookup.  It is able to
    not only ???+ note markata slug lookup !\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Md_It_Wikilinks.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <p>Wikilinks depicted
    with <code>\\[[wikilinks]]</code> can be enabled for sites using the\n<code>markdown-it-py</code>
    plugin markata will look through all posts matching up the\nfile stem to the wiki
    link and inserting the slug as the href.</p>\n<div class=\"admonition note is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">normal wikilink</p>\n<p>Here
    is a link to a markdown file <code>docs/nav.md</code>, and the url becomes\n<a
    href=\"/nav\">/nav</a>.</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span>[[nav]]\n</pre></div>\n\n</pre>\n\n<p>When
    rendered out this wikilink will become an anchor link.</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"p\">&lt;</span><span
    class=\"nt\">a</span> <span class=\"na\">class</span><span class=\"o\">=</span><span
    class=\"s\">&quot;wikilink&quot;</span> <span class=\"na\">href</span><span class=\"o\">=</span><span
    class=\"s\">&quot;/nav&quot;</span><span class=\"p\">&gt;</span>load<span class=\"p\">&lt;/</span><span
    class=\"nt\">a</span><span class=\"p\">&gt;</span>\n</pre></div>\n\n</pre>\n\n<blockquote>\n<p>This
    behaves just like the standard wikilink</p>\n</blockquote>\n</div>\n<p>A special
    feature that markata brings is slug lookup.  It is able to not only\nblindly link
    to the route specified, but will look up the slug of an article.</p>\n<div class=\"admonition
    note is-collapsible collapsible-open\">\n<p class=\"admonition-title\">markata
    slug lookup</p>\n<p>Markata has a load plugin that is generated with the <a class=\"wikilink\"
    href=\"/markata/plugins/docs\">docs</a> plugin.  It's\nfilepath is <code>markata/plugins/load.py</code>,
    so it can be referenced by the file\nstem <code>load</code>.</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span>[[load]]\n</pre></div>\n\n</pre>\n\n<p>Markata
    will look up the article by the file stem, grab the first article,\nand use its
    slug as the href.  This turns it into an anchor link that looks\nlike this.</p>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"p\">&lt;</span><span
    class=\"nt\">a</span> <span class=\"na\">class</span><span class=\"o\">=</span><span
    class=\"s\">&quot;wikilink&quot;</span> <span class=\"na\">href</span><span class=\"o\">=</span><span
    class=\"s\">&quot;/markata/plugins/load&quot;</span><span class=\"p\">&gt;</span>load<span
    class=\"p\">&lt;/</span><span class=\"nt\">a</span><span class=\"p\">&gt;</span>\n</pre></div>\n\n</pre>\n\n</div>\n<p>!!
    function <h2 id='wikilinks_plugin' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>wikilinks_plugin <em class='small'>function</em></h2>\nA plugin to create
    wikilinks tokens.\nThese, token should be handled by the renderer.</p>\n<pre><code>???+
    example\n\n    ```md title=markdown\n    [[nav]]\n    ```\n\n    ```html title=html\n
    \   &lt;a class=&quot;wikilink&quot; href=&quot;/nav&quot;&gt;load&lt;/a&gt;\n
    \   ```\n</code></pre>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">wikilinks_plugin <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">wikilinks_plugin</span><span class=\"p\">(</span>\n            <span
    class=\"n\">md</span><span class=\"p\">:</span> <span class=\"n\">MarkdownIt</span><span
    class=\"p\">,</span>\n            <span class=\"n\">start_delimiter</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;[&quot;</span><span class=\"p\">,</span>\n            <span
    class=\"n\">end_delimiter</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;]&quot;</span><span class=\"p\">,</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">=</span><span class=\"kc\">None</span><span
    class=\"p\">,</span>\n        <span class=\"p\">):</span>\n<span class=\"w\">
    \           </span><span class=\"sd\">&quot;&quot;&quot;A plugin to create wikilinks
    tokens.</span>\n<span class=\"sd\">            These, token should be handled
    by the renderer.</span>\n\n<span class=\"sd\">            ???+ example</span>\n\n<span
    class=\"sd\">                ```md title=markdown</span>\n<span class=\"sd\">
    \               [[nav]]</span>\n<span class=\"sd\">                ```</span>\n\n<span
    class=\"sd\">                ```html title=html</span>\n<span class=\"sd\">                &lt;a
    class=&quot;wikilink&quot; href=&quot;/nav&quot;&gt;load&lt;/a&gt;</span>\n<span
    class=\"sd\">                ```</span>\n<span class=\"sd\">            &quot;&quot;&quot;</span>\n\n
    \           <span class=\"n\">start_char</span> <span class=\"o\">=</span> <span
    class=\"nb\">ord</span><span class=\"p\">(</span><span class=\"n\">start_delimiter</span><span
    class=\"p\">)</span>\n            <span class=\"n\">end_char</span> <span class=\"o\">=</span>
    <span class=\"nb\">ord</span><span class=\"p\">(</span><span class=\"n\">end_delimiter</span><span
    class=\"p\">)</span>\n\n            <span class=\"k\">def</span> <span class=\"nf\">_wikilinks_inline</span><span
    class=\"p\">(</span><span class=\"n\">state</span><span class=\"p\">:</span> <span
    class=\"n\">StateInline</span><span class=\"p\">,</span> <span class=\"n\">silent</span><span
    class=\"p\">:</span> <span class=\"nb\">bool</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">try</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">if</span> <span class=\"p\">(</span>\n                        <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">srcCharCode</span><span
    class=\"p\">[</span><span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">pos</span><span class=\"p\">]</span> <span class=\"o\">!=</span> <span
    class=\"n\">start_char</span>\n                        <span class=\"ow\">or</span>
    <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">srcCharCode</span><span
    class=\"p\">[</span><span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">pos</span> <span class=\"o\">+</span> <span class=\"mi\">1</span><span
    class=\"p\">]</span> <span class=\"o\">!=</span> <span class=\"n\">start_char</span>\n
    \                   <span class=\"p\">):</span>\n                        <span
    class=\"k\">return</span> <span class=\"kc\">False</span>\n                <span
    class=\"k\">except</span> <span class=\"ne\">IndexError</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"kc\">False</span>\n\n
    \               <span class=\"n\">pos</span> <span class=\"o\">=</span> <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">pos</span>
    <span class=\"o\">+</span> <span class=\"mi\">2</span>\n                <span
    class=\"n\">found_closing</span> <span class=\"o\">=</span> <span class=\"kc\">False</span>\n
    \               <span class=\"k\">while</span> <span class=\"kc\">True</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">end</span> <span class=\"o\">=</span>
    <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">srcCharCode</span><span
    class=\"o\">.</span><span class=\"n\">index</span><span class=\"p\">(</span><span
    class=\"n\">end_char</span><span class=\"p\">,</span> <span class=\"n\">pos</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">except</span> <span
    class=\"ne\">ValueError</span><span class=\"p\">:</span>\n                        <span
    class=\"k\">return</span> <span class=\"kc\">False</span>\n                    <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                        <span
    class=\"k\">if</span> <span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">srcCharCode</span><span class=\"p\">[</span><span class=\"n\">end</span>
    <span class=\"o\">+</span> <span class=\"mi\">1</span><span class=\"p\">]</span>
    <span class=\"o\">==</span> <span class=\"n\">end_char</span><span class=\"p\">:</span>\n
    \                           <span class=\"n\">found_closing</span> <span class=\"o\">=</span>
    <span class=\"kc\">True</span>\n                            <span class=\"k\">break</span>\n
    \                   <span class=\"k\">except</span> <span class=\"ne\">IndexError</span><span
    class=\"p\">:</span>\n                        <span class=\"k\">return</span>
    <span class=\"kc\">False</span>\n                    <span class=\"n\">pos</span>
    <span class=\"o\">=</span> <span class=\"n\">end</span> <span class=\"o\">+</span>
    <span class=\"mi\">2</span>\n\n                <span class=\"k\">if</span> <span
    class=\"ow\">not</span> <span class=\"n\">found_closing</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"kc\">False</span>\n\n
    \               <span class=\"n\">text</span> <span class=\"o\">=</span> <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">src</span><span
    class=\"p\">[</span><span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">pos</span> <span class=\"o\">+</span> <span class=\"mi\">2</span>
    <span class=\"p\">:</span> <span class=\"n\">end</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">strip</span><span class=\"p\">()</span>\n
    \               <span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">pos</span> <span class=\"o\">=</span> <span class=\"n\">end</span>
    <span class=\"o\">+</span> <span class=\"mi\">2</span>\n\n                <span
    class=\"k\">if</span> <span class=\"n\">silent</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"kc\">True</span>\n\n
    \               <span class=\"n\">token</span> <span class=\"o\">=</span> <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">push</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;link_open&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;a&quot;</span><span class=\"p\">,</span> <span class=\"mi\">1</span><span
    class=\"p\">)</span>\n                <span class=\"n\">token</span><span class=\"o\">.</span><span
    class=\"n\">block</span> <span class=\"o\">=</span> <span class=\"kc\">False</span>\n
    \               <span class=\"n\">token</span><span class=\"o\">.</span><span
    class=\"n\">attrSet</span><span class=\"p\">(</span><span class=\"s2\">&quot;class&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;wikilink&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">if</span> <span class=\"s2\">&quot;#&quot;</span>
    <span class=\"ow\">in</span> <span class=\"n\">text</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">link</span><span class=\"p\">,</span> <span
    class=\"nb\">id</span> <span class=\"o\">=</span> <span class=\"n\">text</span><span
    class=\"o\">.</span><span class=\"n\">split</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;#&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"n\">link</span> <span class=\"o\">=</span> <span class=\"n\">link</span><span
    class=\"o\">.</span><span class=\"n\">strip</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;/&quot;</span><span class=\"p\">)</span>\n                <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                    <span class=\"n\">link</span><span
    class=\"p\">,</span> <span class=\"nb\">id</span> <span class=\"o\">=</span> <span
    class=\"n\">text</span><span class=\"p\">,</span> <span class=\"kc\">None</span>\n
    \               <span class=\"n\">possible_pages</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">filter</span><span
    class=\"p\">(</span>\n                    <span class=\"sa\">f</span><span class=\"s1\">&#39;str(path).split(&quot;/&quot;)[-1].split(&quot;.&quot;)[0].replace(&quot;_&quot;,
    &quot;-&quot;) == &quot;</span><span class=\"si\">{</span><span class=\"n\">link</span><span
    class=\"o\">.</span><span class=\"n\">replace</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;_&quot;</span><span class=\"p\">,</span><span class=\"w\">
    </span><span class=\"s2\">&quot;-&quot;</span><span class=\"p\">)</span><span
    class=\"si\">}</span><span class=\"s1\">&quot;&#39;</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">)</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">len</span><span class=\"p\">(</span><span class=\"n\">possible_pages</span><span
    class=\"p\">)</span> <span class=\"o\">==</span> <span class=\"mi\">1</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">link</span> <span
    class=\"o\">=</span> <span class=\"n\">possible_pages</span><span class=\"p\">[</span><span
    class=\"mi\">0</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">,</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;/</span><span
    class=\"si\">{</span><span class=\"n\">text</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span><span class=\"p\">)</span>\n                <span class=\"k\">elif</span>
    <span class=\"nb\">len</span><span class=\"p\">(</span><span class=\"n\">possible_pages</span><span
    class=\"p\">)</span> <span class=\"o\">&gt;</span> <span class=\"mi\">1</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">logger</span><span
    class=\"o\">.</span><span class=\"n\">warning</span><span class=\"p\">(</span>\n
    \                       <span class=\"sa\">f</span><span class=\"s2\">&quot;wikilink
    [[</span><span class=\"si\">{</span><span class=\"n\">text</span><span class=\"si\">}</span><span
    class=\"s2\">]] (</span><span class=\"si\">{</span><span class=\"n\">link</span><span
    class=\"si\">}</span><span class=\"s2\">, </span><span class=\"si\">{</span><span
    class=\"nb\">id</span><span class=\"si\">}</span><span class=\"s2\">) has duplicate
    matches, defaulting to the first&quot;</span><span class=\"p\">,</span>\n                    <span
    class=\"p\">)</span>\n                    <span class=\"n\">link</span> <span
    class=\"o\">=</span> <span class=\"n\">possible_pages</span><span class=\"p\">[</span><span
    class=\"mi\">0</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">,</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;/</span><span
    class=\"si\">{</span><span class=\"n\">text</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span><span class=\"p\">)</span>\n                <span class=\"k\">else</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">logger</span><span
    class=\"o\">.</span><span class=\"n\">warning</span><span class=\"p\">(</span>\n
    \                       <span class=\"sa\">f</span><span class=\"s2\">&quot;wikilink
    [[</span><span class=\"si\">{</span><span class=\"n\">text</span><span class=\"si\">}</span><span
    class=\"s2\">]] (</span><span class=\"si\">{</span><span class=\"n\">link</span><span
    class=\"si\">}</span><span class=\"s2\">, </span><span class=\"si\">{</span><span
    class=\"nb\">id</span><span class=\"si\">}</span><span class=\"s2\">) no matches,
    defaulting to &#39;/</span><span class=\"si\">{</span><span class=\"n\">text</span><span
    class=\"si\">}</span><span class=\"s2\">&#39;&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"p\">)</span>\n                    <span class=\"n\">link</span>
    <span class=\"o\">=</span> <span class=\"n\">text</span>\n\n                <span
    class=\"k\">if</span> <span class=\"nb\">id</span> <span class=\"ow\">and</span>
    <span class=\"ow\">not</span> <span class=\"n\">link</span><span class=\"o\">.</span><span
    class=\"n\">endswith</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;#</span><span class=\"si\">{</span><span class=\"nb\">id</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">):</span>\n
    \                   <span class=\"n\">link</span> <span class=\"o\">=</span> <span
    class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">link</span><span class=\"si\">}</span><span class=\"s2\">#</span><span
    class=\"si\">{</span><span class=\"nb\">id</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span>\n\n                <span class=\"n\">token</span><span
    class=\"o\">.</span><span class=\"n\">attrSet</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;href&quot;</span><span class=\"p\">,</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;/</span><span class=\"si\">{</span><span class=\"n\">link</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"n\">content_token</span> <span class=\"o\">=</span>
    <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">push</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;text&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;&quot;</span><span class=\"p\">,</span> <span class=\"mi\">0</span><span
    class=\"p\">)</span>\n                <span class=\"n\">content_token</span><span
    class=\"o\">.</span><span class=\"n\">content</span> <span class=\"o\">=</span>
    <span class=\"n\">text</span>\n\n                <span class=\"n\">token</span>
    <span class=\"o\">=</span> <span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">push</span><span class=\"p\">(</span><span class=\"s2\">&quot;link_close&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;a&quot;</span><span class=\"p\">,</span>
    <span class=\"o\">-</span><span class=\"mi\">1</span><span class=\"p\">)</span>\n
    \               <span class=\"n\">token</span><span class=\"o\">.</span><span
    class=\"n\">content</span> <span class=\"o\">=</span> <span class=\"n\">text</span>\n\n
    \               <span class=\"k\">return</span> <span class=\"kc\">True</span>\n\n
    \           <span class=\"n\">md</span><span class=\"o\">.</span><span class=\"n\">inline</span><span
    class=\"o\">.</span><span class=\"n\">ruler</span><span class=\"o\">.</span><span
    class=\"n\">before</span><span class=\"p\">(</span><span class=\"s2\">&quot;escape&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;wikilinks_inline&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">_wikilinks_inline</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='_wikilinks_inline' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>_wikilinks_inline <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_wikilinks_inline
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
    <span class=\"nf\">_wikilinks_inline</span><span class=\"p\">(</span><span class=\"n\">state</span><span
    class=\"p\">:</span> <span class=\"n\">StateInline</span><span class=\"p\">,</span>
    <span class=\"n\">silent</span><span class=\"p\">:</span> <span class=\"nb\">bool</span><span
    class=\"p\">):</span>\n                <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">if</span> <span class=\"p\">(</span>\n                        <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">srcCharCode</span><span
    class=\"p\">[</span><span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">pos</span><span class=\"p\">]</span> <span class=\"o\">!=</span> <span
    class=\"n\">start_char</span>\n                        <span class=\"ow\">or</span>
    <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">srcCharCode</span><span
    class=\"p\">[</span><span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">pos</span> <span class=\"o\">+</span> <span class=\"mi\">1</span><span
    class=\"p\">]</span> <span class=\"o\">!=</span> <span class=\"n\">start_char</span>\n
    \                   <span class=\"p\">):</span>\n                        <span
    class=\"k\">return</span> <span class=\"kc\">False</span>\n                <span
    class=\"k\">except</span> <span class=\"ne\">IndexError</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"kc\">False</span>\n\n
    \               <span class=\"n\">pos</span> <span class=\"o\">=</span> <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">pos</span>
    <span class=\"o\">+</span> <span class=\"mi\">2</span>\n                <span
    class=\"n\">found_closing</span> <span class=\"o\">=</span> <span class=\"kc\">False</span>\n
    \               <span class=\"k\">while</span> <span class=\"kc\">True</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">end</span> <span class=\"o\">=</span>
    <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">srcCharCode</span><span
    class=\"o\">.</span><span class=\"n\">index</span><span class=\"p\">(</span><span
    class=\"n\">end_char</span><span class=\"p\">,</span> <span class=\"n\">pos</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">except</span> <span
    class=\"ne\">ValueError</span><span class=\"p\">:</span>\n                        <span
    class=\"k\">return</span> <span class=\"kc\">False</span>\n                    <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                        <span
    class=\"k\">if</span> <span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">srcCharCode</span><span class=\"p\">[</span><span class=\"n\">end</span>
    <span class=\"o\">+</span> <span class=\"mi\">1</span><span class=\"p\">]</span>
    <span class=\"o\">==</span> <span class=\"n\">end_char</span><span class=\"p\">:</span>\n
    \                           <span class=\"n\">found_closing</span> <span class=\"o\">=</span>
    <span class=\"kc\">True</span>\n                            <span class=\"k\">break</span>\n
    \                   <span class=\"k\">except</span> <span class=\"ne\">IndexError</span><span
    class=\"p\">:</span>\n                        <span class=\"k\">return</span>
    <span class=\"kc\">False</span>\n                    <span class=\"n\">pos</span>
    <span class=\"o\">=</span> <span class=\"n\">end</span> <span class=\"o\">+</span>
    <span class=\"mi\">2</span>\n\n                <span class=\"k\">if</span> <span
    class=\"ow\">not</span> <span class=\"n\">found_closing</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"kc\">False</span>\n\n
    \               <span class=\"n\">text</span> <span class=\"o\">=</span> <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">src</span><span
    class=\"p\">[</span><span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">pos</span> <span class=\"o\">+</span> <span class=\"mi\">2</span>
    <span class=\"p\">:</span> <span class=\"n\">end</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">strip</span><span class=\"p\">()</span>\n
    \               <span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">pos</span> <span class=\"o\">=</span> <span class=\"n\">end</span>
    <span class=\"o\">+</span> <span class=\"mi\">2</span>\n\n                <span
    class=\"k\">if</span> <span class=\"n\">silent</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"kc\">True</span>\n\n
    \               <span class=\"n\">token</span> <span class=\"o\">=</span> <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">push</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;link_open&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;a&quot;</span><span class=\"p\">,</span> <span class=\"mi\">1</span><span
    class=\"p\">)</span>\n                <span class=\"n\">token</span><span class=\"o\">.</span><span
    class=\"n\">block</span> <span class=\"o\">=</span> <span class=\"kc\">False</span>\n
    \               <span class=\"n\">token</span><span class=\"o\">.</span><span
    class=\"n\">attrSet</span><span class=\"p\">(</span><span class=\"s2\">&quot;class&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;wikilink&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">if</span> <span class=\"s2\">&quot;#&quot;</span>
    <span class=\"ow\">in</span> <span class=\"n\">text</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">link</span><span class=\"p\">,</span> <span
    class=\"nb\">id</span> <span class=\"o\">=</span> <span class=\"n\">text</span><span
    class=\"o\">.</span><span class=\"n\">split</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;#&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"n\">link</span> <span class=\"o\">=</span> <span class=\"n\">link</span><span
    class=\"o\">.</span><span class=\"n\">strip</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;/&quot;</span><span class=\"p\">)</span>\n                <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                    <span class=\"n\">link</span><span
    class=\"p\">,</span> <span class=\"nb\">id</span> <span class=\"o\">=</span> <span
    class=\"n\">text</span><span class=\"p\">,</span> <span class=\"kc\">None</span>\n
    \               <span class=\"n\">possible_pages</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">filter</span><span
    class=\"p\">(</span>\n                    <span class=\"sa\">f</span><span class=\"s1\">&#39;str(path).split(&quot;/&quot;)[-1].split(&quot;.&quot;)[0].replace(&quot;_&quot;,
    &quot;-&quot;) == &quot;</span><span class=\"si\">{</span><span class=\"n\">link</span><span
    class=\"o\">.</span><span class=\"n\">replace</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;_&quot;</span><span class=\"p\">,</span><span class=\"w\">
    </span><span class=\"s2\">&quot;-&quot;</span><span class=\"p\">)</span><span
    class=\"si\">}</span><span class=\"s1\">&quot;&#39;</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">)</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">len</span><span class=\"p\">(</span><span class=\"n\">possible_pages</span><span
    class=\"p\">)</span> <span class=\"o\">==</span> <span class=\"mi\">1</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">link</span> <span
    class=\"o\">=</span> <span class=\"n\">possible_pages</span><span class=\"p\">[</span><span
    class=\"mi\">0</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">,</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;/</span><span
    class=\"si\">{</span><span class=\"n\">text</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span><span class=\"p\">)</span>\n                <span class=\"k\">elif</span>
    <span class=\"nb\">len</span><span class=\"p\">(</span><span class=\"n\">possible_pages</span><span
    class=\"p\">)</span> <span class=\"o\">&gt;</span> <span class=\"mi\">1</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">logger</span><span
    class=\"o\">.</span><span class=\"n\">warning</span><span class=\"p\">(</span>\n
    \                       <span class=\"sa\">f</span><span class=\"s2\">&quot;wikilink
    [[</span><span class=\"si\">{</span><span class=\"n\">text</span><span class=\"si\">}</span><span
    class=\"s2\">]] (</span><span class=\"si\">{</span><span class=\"n\">link</span><span
    class=\"si\">}</span><span class=\"s2\">, </span><span class=\"si\">{</span><span
    class=\"nb\">id</span><span class=\"si\">}</span><span class=\"s2\">) has duplicate
    matches, defaulting to the first&quot;</span><span class=\"p\">,</span>\n                    <span
    class=\"p\">)</span>\n                    <span class=\"n\">link</span> <span
    class=\"o\">=</span> <span class=\"n\">possible_pages</span><span class=\"p\">[</span><span
    class=\"mi\">0</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">,</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;/</span><span
    class=\"si\">{</span><span class=\"n\">text</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span><span class=\"p\">)</span>\n                <span class=\"k\">else</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">logger</span><span
    class=\"o\">.</span><span class=\"n\">warning</span><span class=\"p\">(</span>\n
    \                       <span class=\"sa\">f</span><span class=\"s2\">&quot;wikilink
    [[</span><span class=\"si\">{</span><span class=\"n\">text</span><span class=\"si\">}</span><span
    class=\"s2\">]] (</span><span class=\"si\">{</span><span class=\"n\">link</span><span
    class=\"si\">}</span><span class=\"s2\">, </span><span class=\"si\">{</span><span
    class=\"nb\">id</span><span class=\"si\">}</span><span class=\"s2\">) no matches,
    defaulting to &#39;/</span><span class=\"si\">{</span><span class=\"n\">text</span><span
    class=\"si\">}</span><span class=\"s2\">&#39;&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"p\">)</span>\n                    <span class=\"n\">link</span>
    <span class=\"o\">=</span> <span class=\"n\">text</span>\n\n                <span
    class=\"k\">if</span> <span class=\"nb\">id</span> <span class=\"ow\">and</span>
    <span class=\"ow\">not</span> <span class=\"n\">link</span><span class=\"o\">.</span><span
    class=\"n\">endswith</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;#</span><span class=\"si\">{</span><span class=\"nb\">id</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">):</span>\n
    \                   <span class=\"n\">link</span> <span class=\"o\">=</span> <span
    class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">link</span><span class=\"si\">}</span><span class=\"s2\">#</span><span
    class=\"si\">{</span><span class=\"nb\">id</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span>\n\n                <span class=\"n\">token</span><span
    class=\"o\">.</span><span class=\"n\">attrSet</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;href&quot;</span><span class=\"p\">,</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;/</span><span class=\"si\">{</span><span class=\"n\">link</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"n\">content_token</span> <span class=\"o\">=</span>
    <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">push</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;text&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;&quot;</span><span class=\"p\">,</span> <span class=\"mi\">0</span><span
    class=\"p\">)</span>\n                <span class=\"n\">content_token</span><span
    class=\"o\">.</span><span class=\"n\">content</span> <span class=\"o\">=</span>
    <span class=\"n\">text</span>\n\n                <span class=\"n\">token</span>
    <span class=\"o\">=</span> <span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">push</span><span class=\"p\">(</span><span class=\"s2\">&quot;link_close&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;a&quot;</span><span class=\"p\">,</span>
    <span class=\"o\">-</span><span class=\"mi\">1</span><span class=\"p\">)</span>\n
    \               <span class=\"n\">token</span><span class=\"o\">.</span><span
    class=\"n\">content</span> <span class=\"o\">=</span> <span class=\"n\">text</span>\n\n
    \               <span class=\"k\">return</span> <span class=\"kc\">True</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>    </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Md_It_Wikilinks.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Wikilinks depicted with  ???+ note normal
    wikilink A special feature that markata brings is slug lookup.  It is able to
    not only ???+ note markata slug lookup !\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Md_It_Wikilinks.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Wikilinks depicted with  ???+ note normal
    wikilink A special feature that markata brings is slug lookup.  It is able to
    not only ???+ note markata slug lookup !\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<article style=\"text-align: center;\">\n    <style>\n        section
    {\n            font-size: 200%;\n        }\n\n\n        .edit {\n            display:
    none;\n        }\n    </style>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Md_It_Wikilinks.Py\n    </h1>\n</section></article>\n     </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Md_It_Wikilinks.Py\n    </h1>\n</section>    <section class=\"body\">\n
    \       <p>Wikilinks depicted with <code>\\[[wikilinks]]</code> can be enabled
    for sites using the\n<code>markdown-it-py</code> plugin markata will look through
    all posts matching up the\nfile stem to the wiki link and inserting the slug as
    the href.</p>\n<div class=\"admonition note is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">normal wikilink</p>\n<p>Here is a link to a markdown
    file <code>docs/nav.md</code>, and the url becomes\n<a href=\"/nav\">/nav</a>.</p>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span>[[nav]]\n</pre></div>\n\n</pre>\n\n<p>When
    rendered out this wikilink will become an anchor link.</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"p\">&lt;</span><span
    class=\"nt\">a</span> <span class=\"na\">class</span><span class=\"o\">=</span><span
    class=\"s\">&quot;wikilink&quot;</span> <span class=\"na\">href</span><span class=\"o\">=</span><span
    class=\"s\">&quot;/nav&quot;</span><span class=\"p\">&gt;</span>load<span class=\"p\">&lt;/</span><span
    class=\"nt\">a</span><span class=\"p\">&gt;</span>\n</pre></div>\n\n</pre>\n\n<blockquote>\n<p>This
    behaves just like the standard wikilink</p>\n</blockquote>\n</div>\n<p>A special
    feature that markata brings is slug lookup.  It is able to not only\nblindly link
    to the route specified, but will look up the slug of an article.</p>\n<div class=\"admonition
    note is-collapsible collapsible-open\">\n<p class=\"admonition-title\">markata
    slug lookup</p>\n<p>Markata has a load plugin that is generated with the <a class=\"wikilink\"
    href=\"/markata/plugins/docs\">docs</a> plugin.  It's\nfilepath is <code>markata/plugins/load.py</code>,
    so it can be referenced by the file\nstem <code>load</code>.</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span>[[load]]\n</pre></div>\n\n</pre>\n\n<p>Markata
    will look up the article by the file stem, grab the first article,\nand use its
    slug as the href.  This turns it into an anchor link that looks\nlike this.</p>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"p\">&lt;</span><span
    class=\"nt\">a</span> <span class=\"na\">class</span><span class=\"o\">=</span><span
    class=\"s\">&quot;wikilink&quot;</span> <span class=\"na\">href</span><span class=\"o\">=</span><span
    class=\"s\">&quot;/markata/plugins/load&quot;</span><span class=\"p\">&gt;</span>load<span
    class=\"p\">&lt;/</span><span class=\"nt\">a</span><span class=\"p\">&gt;</span>\n</pre></div>\n\n</pre>\n\n</div>\n<p>!!
    function <h2 id='wikilinks_plugin' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>wikilinks_plugin <em class='small'>function</em></h2>\nA plugin to create
    wikilinks tokens.\nThese, token should be handled by the renderer.</p>\n<pre><code>???+
    example\n\n    ```md title=markdown\n    [[nav]]\n    ```\n\n    ```html title=html\n
    \   &lt;a class=&quot;wikilink&quot; href=&quot;/nav&quot;&gt;load&lt;/a&gt;\n
    \   ```\n</code></pre>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">wikilinks_plugin <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">wikilinks_plugin</span><span class=\"p\">(</span>\n            <span
    class=\"n\">md</span><span class=\"p\">:</span> <span class=\"n\">MarkdownIt</span><span
    class=\"p\">,</span>\n            <span class=\"n\">start_delimiter</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;[&quot;</span><span class=\"p\">,</span>\n            <span
    class=\"n\">end_delimiter</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;]&quot;</span><span class=\"p\">,</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">=</span><span class=\"kc\">None</span><span
    class=\"p\">,</span>\n        <span class=\"p\">):</span>\n<span class=\"w\">
    \           </span><span class=\"sd\">&quot;&quot;&quot;A plugin to create wikilinks
    tokens.</span>\n<span class=\"sd\">            These, token should be handled
    by the renderer.</span>\n\n<span class=\"sd\">            ???+ example</span>\n\n<span
    class=\"sd\">                ```md title=markdown</span>\n<span class=\"sd\">
    \               [[nav]]</span>\n<span class=\"sd\">                ```</span>\n\n<span
    class=\"sd\">                ```html title=html</span>\n<span class=\"sd\">                &lt;a
    class=&quot;wikilink&quot; href=&quot;/nav&quot;&gt;load&lt;/a&gt;</span>\n<span
    class=\"sd\">                ```</span>\n<span class=\"sd\">            &quot;&quot;&quot;</span>\n\n
    \           <span class=\"n\">start_char</span> <span class=\"o\">=</span> <span
    class=\"nb\">ord</span><span class=\"p\">(</span><span class=\"n\">start_delimiter</span><span
    class=\"p\">)</span>\n            <span class=\"n\">end_char</span> <span class=\"o\">=</span>
    <span class=\"nb\">ord</span><span class=\"p\">(</span><span class=\"n\">end_delimiter</span><span
    class=\"p\">)</span>\n\n            <span class=\"k\">def</span> <span class=\"nf\">_wikilinks_inline</span><span
    class=\"p\">(</span><span class=\"n\">state</span><span class=\"p\">:</span> <span
    class=\"n\">StateInline</span><span class=\"p\">,</span> <span class=\"n\">silent</span><span
    class=\"p\">:</span> <span class=\"nb\">bool</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">try</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">if</span> <span class=\"p\">(</span>\n                        <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">srcCharCode</span><span
    class=\"p\">[</span><span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">pos</span><span class=\"p\">]</span> <span class=\"o\">!=</span> <span
    class=\"n\">start_char</span>\n                        <span class=\"ow\">or</span>
    <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">srcCharCode</span><span
    class=\"p\">[</span><span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">pos</span> <span class=\"o\">+</span> <span class=\"mi\">1</span><span
    class=\"p\">]</span> <span class=\"o\">!=</span> <span class=\"n\">start_char</span>\n
    \                   <span class=\"p\">):</span>\n                        <span
    class=\"k\">return</span> <span class=\"kc\">False</span>\n                <span
    class=\"k\">except</span> <span class=\"ne\">IndexError</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"kc\">False</span>\n\n
    \               <span class=\"n\">pos</span> <span class=\"o\">=</span> <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">pos</span>
    <span class=\"o\">+</span> <span class=\"mi\">2</span>\n                <span
    class=\"n\">found_closing</span> <span class=\"o\">=</span> <span class=\"kc\">False</span>\n
    \               <span class=\"k\">while</span> <span class=\"kc\">True</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">end</span> <span class=\"o\">=</span>
    <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">srcCharCode</span><span
    class=\"o\">.</span><span class=\"n\">index</span><span class=\"p\">(</span><span
    class=\"n\">end_char</span><span class=\"p\">,</span> <span class=\"n\">pos</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">except</span> <span
    class=\"ne\">ValueError</span><span class=\"p\">:</span>\n                        <span
    class=\"k\">return</span> <span class=\"kc\">False</span>\n                    <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                        <span
    class=\"k\">if</span> <span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">srcCharCode</span><span class=\"p\">[</span><span class=\"n\">end</span>
    <span class=\"o\">+</span> <span class=\"mi\">1</span><span class=\"p\">]</span>
    <span class=\"o\">==</span> <span class=\"n\">end_char</span><span class=\"p\">:</span>\n
    \                           <span class=\"n\">found_closing</span> <span class=\"o\">=</span>
    <span class=\"kc\">True</span>\n                            <span class=\"k\">break</span>\n
    \                   <span class=\"k\">except</span> <span class=\"ne\">IndexError</span><span
    class=\"p\">:</span>\n                        <span class=\"k\">return</span>
    <span class=\"kc\">False</span>\n                    <span class=\"n\">pos</span>
    <span class=\"o\">=</span> <span class=\"n\">end</span> <span class=\"o\">+</span>
    <span class=\"mi\">2</span>\n\n                <span class=\"k\">if</span> <span
    class=\"ow\">not</span> <span class=\"n\">found_closing</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"kc\">False</span>\n\n
    \               <span class=\"n\">text</span> <span class=\"o\">=</span> <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">src</span><span
    class=\"p\">[</span><span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">pos</span> <span class=\"o\">+</span> <span class=\"mi\">2</span>
    <span class=\"p\">:</span> <span class=\"n\">end</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">strip</span><span class=\"p\">()</span>\n
    \               <span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">pos</span> <span class=\"o\">=</span> <span class=\"n\">end</span>
    <span class=\"o\">+</span> <span class=\"mi\">2</span>\n\n                <span
    class=\"k\">if</span> <span class=\"n\">silent</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"kc\">True</span>\n\n
    \               <span class=\"n\">token</span> <span class=\"o\">=</span> <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">push</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;link_open&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;a&quot;</span><span class=\"p\">,</span> <span class=\"mi\">1</span><span
    class=\"p\">)</span>\n                <span class=\"n\">token</span><span class=\"o\">.</span><span
    class=\"n\">block</span> <span class=\"o\">=</span> <span class=\"kc\">False</span>\n
    \               <span class=\"n\">token</span><span class=\"o\">.</span><span
    class=\"n\">attrSet</span><span class=\"p\">(</span><span class=\"s2\">&quot;class&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;wikilink&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">if</span> <span class=\"s2\">&quot;#&quot;</span>
    <span class=\"ow\">in</span> <span class=\"n\">text</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">link</span><span class=\"p\">,</span> <span
    class=\"nb\">id</span> <span class=\"o\">=</span> <span class=\"n\">text</span><span
    class=\"o\">.</span><span class=\"n\">split</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;#&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"n\">link</span> <span class=\"o\">=</span> <span class=\"n\">link</span><span
    class=\"o\">.</span><span class=\"n\">strip</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;/&quot;</span><span class=\"p\">)</span>\n                <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                    <span class=\"n\">link</span><span
    class=\"p\">,</span> <span class=\"nb\">id</span> <span class=\"o\">=</span> <span
    class=\"n\">text</span><span class=\"p\">,</span> <span class=\"kc\">None</span>\n
    \               <span class=\"n\">possible_pages</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">filter</span><span
    class=\"p\">(</span>\n                    <span class=\"sa\">f</span><span class=\"s1\">&#39;str(path).split(&quot;/&quot;)[-1].split(&quot;.&quot;)[0].replace(&quot;_&quot;,
    &quot;-&quot;) == &quot;</span><span class=\"si\">{</span><span class=\"n\">link</span><span
    class=\"o\">.</span><span class=\"n\">replace</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;_&quot;</span><span class=\"p\">,</span><span class=\"w\">
    </span><span class=\"s2\">&quot;-&quot;</span><span class=\"p\">)</span><span
    class=\"si\">}</span><span class=\"s1\">&quot;&#39;</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">)</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">len</span><span class=\"p\">(</span><span class=\"n\">possible_pages</span><span
    class=\"p\">)</span> <span class=\"o\">==</span> <span class=\"mi\">1</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">link</span> <span
    class=\"o\">=</span> <span class=\"n\">possible_pages</span><span class=\"p\">[</span><span
    class=\"mi\">0</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">,</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;/</span><span
    class=\"si\">{</span><span class=\"n\">text</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span><span class=\"p\">)</span>\n                <span class=\"k\">elif</span>
    <span class=\"nb\">len</span><span class=\"p\">(</span><span class=\"n\">possible_pages</span><span
    class=\"p\">)</span> <span class=\"o\">&gt;</span> <span class=\"mi\">1</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">logger</span><span
    class=\"o\">.</span><span class=\"n\">warning</span><span class=\"p\">(</span>\n
    \                       <span class=\"sa\">f</span><span class=\"s2\">&quot;wikilink
    [[</span><span class=\"si\">{</span><span class=\"n\">text</span><span class=\"si\">}</span><span
    class=\"s2\">]] (</span><span class=\"si\">{</span><span class=\"n\">link</span><span
    class=\"si\">}</span><span class=\"s2\">, </span><span class=\"si\">{</span><span
    class=\"nb\">id</span><span class=\"si\">}</span><span class=\"s2\">) has duplicate
    matches, defaulting to the first&quot;</span><span class=\"p\">,</span>\n                    <span
    class=\"p\">)</span>\n                    <span class=\"n\">link</span> <span
    class=\"o\">=</span> <span class=\"n\">possible_pages</span><span class=\"p\">[</span><span
    class=\"mi\">0</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">,</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;/</span><span
    class=\"si\">{</span><span class=\"n\">text</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span><span class=\"p\">)</span>\n                <span class=\"k\">else</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">logger</span><span
    class=\"o\">.</span><span class=\"n\">warning</span><span class=\"p\">(</span>\n
    \                       <span class=\"sa\">f</span><span class=\"s2\">&quot;wikilink
    [[</span><span class=\"si\">{</span><span class=\"n\">text</span><span class=\"si\">}</span><span
    class=\"s2\">]] (</span><span class=\"si\">{</span><span class=\"n\">link</span><span
    class=\"si\">}</span><span class=\"s2\">, </span><span class=\"si\">{</span><span
    class=\"nb\">id</span><span class=\"si\">}</span><span class=\"s2\">) no matches,
    defaulting to &#39;/</span><span class=\"si\">{</span><span class=\"n\">text</span><span
    class=\"si\">}</span><span class=\"s2\">&#39;&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"p\">)</span>\n                    <span class=\"n\">link</span>
    <span class=\"o\">=</span> <span class=\"n\">text</span>\n\n                <span
    class=\"k\">if</span> <span class=\"nb\">id</span> <span class=\"ow\">and</span>
    <span class=\"ow\">not</span> <span class=\"n\">link</span><span class=\"o\">.</span><span
    class=\"n\">endswith</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;#</span><span class=\"si\">{</span><span class=\"nb\">id</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">):</span>\n
    \                   <span class=\"n\">link</span> <span class=\"o\">=</span> <span
    class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">link</span><span class=\"si\">}</span><span class=\"s2\">#</span><span
    class=\"si\">{</span><span class=\"nb\">id</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span>\n\n                <span class=\"n\">token</span><span
    class=\"o\">.</span><span class=\"n\">attrSet</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;href&quot;</span><span class=\"p\">,</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;/</span><span class=\"si\">{</span><span class=\"n\">link</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"n\">content_token</span> <span class=\"o\">=</span>
    <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">push</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;text&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;&quot;</span><span class=\"p\">,</span> <span class=\"mi\">0</span><span
    class=\"p\">)</span>\n                <span class=\"n\">content_token</span><span
    class=\"o\">.</span><span class=\"n\">content</span> <span class=\"o\">=</span>
    <span class=\"n\">text</span>\n\n                <span class=\"n\">token</span>
    <span class=\"o\">=</span> <span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">push</span><span class=\"p\">(</span><span class=\"s2\">&quot;link_close&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;a&quot;</span><span class=\"p\">,</span>
    <span class=\"o\">-</span><span class=\"mi\">1</span><span class=\"p\">)</span>\n
    \               <span class=\"n\">token</span><span class=\"o\">.</span><span
    class=\"n\">content</span> <span class=\"o\">=</span> <span class=\"n\">text</span>\n\n
    \               <span class=\"k\">return</span> <span class=\"kc\">True</span>\n\n
    \           <span class=\"n\">md</span><span class=\"o\">.</span><span class=\"n\">inline</span><span
    class=\"o\">.</span><span class=\"n\">ruler</span><span class=\"o\">.</span><span
    class=\"n\">before</span><span class=\"p\">(</span><span class=\"s2\">&quot;escape&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;wikilinks_inline&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">_wikilinks_inline</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='_wikilinks_inline' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>_wikilinks_inline <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_wikilinks_inline
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
    <span class=\"nf\">_wikilinks_inline</span><span class=\"p\">(</span><span class=\"n\">state</span><span
    class=\"p\">:</span> <span class=\"n\">StateInline</span><span class=\"p\">,</span>
    <span class=\"n\">silent</span><span class=\"p\">:</span> <span class=\"nb\">bool</span><span
    class=\"p\">):</span>\n                <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">if</span> <span class=\"p\">(</span>\n                        <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">srcCharCode</span><span
    class=\"p\">[</span><span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">pos</span><span class=\"p\">]</span> <span class=\"o\">!=</span> <span
    class=\"n\">start_char</span>\n                        <span class=\"ow\">or</span>
    <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">srcCharCode</span><span
    class=\"p\">[</span><span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">pos</span> <span class=\"o\">+</span> <span class=\"mi\">1</span><span
    class=\"p\">]</span> <span class=\"o\">!=</span> <span class=\"n\">start_char</span>\n
    \                   <span class=\"p\">):</span>\n                        <span
    class=\"k\">return</span> <span class=\"kc\">False</span>\n                <span
    class=\"k\">except</span> <span class=\"ne\">IndexError</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"kc\">False</span>\n\n
    \               <span class=\"n\">pos</span> <span class=\"o\">=</span> <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">pos</span>
    <span class=\"o\">+</span> <span class=\"mi\">2</span>\n                <span
    class=\"n\">found_closing</span> <span class=\"o\">=</span> <span class=\"kc\">False</span>\n
    \               <span class=\"k\">while</span> <span class=\"kc\">True</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">end</span> <span class=\"o\">=</span>
    <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">srcCharCode</span><span
    class=\"o\">.</span><span class=\"n\">index</span><span class=\"p\">(</span><span
    class=\"n\">end_char</span><span class=\"p\">,</span> <span class=\"n\">pos</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">except</span> <span
    class=\"ne\">ValueError</span><span class=\"p\">:</span>\n                        <span
    class=\"k\">return</span> <span class=\"kc\">False</span>\n                    <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                        <span
    class=\"k\">if</span> <span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">srcCharCode</span><span class=\"p\">[</span><span class=\"n\">end</span>
    <span class=\"o\">+</span> <span class=\"mi\">1</span><span class=\"p\">]</span>
    <span class=\"o\">==</span> <span class=\"n\">end_char</span><span class=\"p\">:</span>\n
    \                           <span class=\"n\">found_closing</span> <span class=\"o\">=</span>
    <span class=\"kc\">True</span>\n                            <span class=\"k\">break</span>\n
    \                   <span class=\"k\">except</span> <span class=\"ne\">IndexError</span><span
    class=\"p\">:</span>\n                        <span class=\"k\">return</span>
    <span class=\"kc\">False</span>\n                    <span class=\"n\">pos</span>
    <span class=\"o\">=</span> <span class=\"n\">end</span> <span class=\"o\">+</span>
    <span class=\"mi\">2</span>\n\n                <span class=\"k\">if</span> <span
    class=\"ow\">not</span> <span class=\"n\">found_closing</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"kc\">False</span>\n\n
    \               <span class=\"n\">text</span> <span class=\"o\">=</span> <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">src</span><span
    class=\"p\">[</span><span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">pos</span> <span class=\"o\">+</span> <span class=\"mi\">2</span>
    <span class=\"p\">:</span> <span class=\"n\">end</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">strip</span><span class=\"p\">()</span>\n
    \               <span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">pos</span> <span class=\"o\">=</span> <span class=\"n\">end</span>
    <span class=\"o\">+</span> <span class=\"mi\">2</span>\n\n                <span
    class=\"k\">if</span> <span class=\"n\">silent</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"kc\">True</span>\n\n
    \               <span class=\"n\">token</span> <span class=\"o\">=</span> <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">push</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;link_open&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;a&quot;</span><span class=\"p\">,</span> <span class=\"mi\">1</span><span
    class=\"p\">)</span>\n                <span class=\"n\">token</span><span class=\"o\">.</span><span
    class=\"n\">block</span> <span class=\"o\">=</span> <span class=\"kc\">False</span>\n
    \               <span class=\"n\">token</span><span class=\"o\">.</span><span
    class=\"n\">attrSet</span><span class=\"p\">(</span><span class=\"s2\">&quot;class&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;wikilink&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">if</span> <span class=\"s2\">&quot;#&quot;</span>
    <span class=\"ow\">in</span> <span class=\"n\">text</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">link</span><span class=\"p\">,</span> <span
    class=\"nb\">id</span> <span class=\"o\">=</span> <span class=\"n\">text</span><span
    class=\"o\">.</span><span class=\"n\">split</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;#&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"n\">link</span> <span class=\"o\">=</span> <span class=\"n\">link</span><span
    class=\"o\">.</span><span class=\"n\">strip</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;/&quot;</span><span class=\"p\">)</span>\n                <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                    <span class=\"n\">link</span><span
    class=\"p\">,</span> <span class=\"nb\">id</span> <span class=\"o\">=</span> <span
    class=\"n\">text</span><span class=\"p\">,</span> <span class=\"kc\">None</span>\n
    \               <span class=\"n\">possible_pages</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">filter</span><span
    class=\"p\">(</span>\n                    <span class=\"sa\">f</span><span class=\"s1\">&#39;str(path).split(&quot;/&quot;)[-1].split(&quot;.&quot;)[0].replace(&quot;_&quot;,
    &quot;-&quot;) == &quot;</span><span class=\"si\">{</span><span class=\"n\">link</span><span
    class=\"o\">.</span><span class=\"n\">replace</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;_&quot;</span><span class=\"p\">,</span><span class=\"w\">
    </span><span class=\"s2\">&quot;-&quot;</span><span class=\"p\">)</span><span
    class=\"si\">}</span><span class=\"s1\">&quot;&#39;</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">)</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">len</span><span class=\"p\">(</span><span class=\"n\">possible_pages</span><span
    class=\"p\">)</span> <span class=\"o\">==</span> <span class=\"mi\">1</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">link</span> <span
    class=\"o\">=</span> <span class=\"n\">possible_pages</span><span class=\"p\">[</span><span
    class=\"mi\">0</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">,</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;/</span><span
    class=\"si\">{</span><span class=\"n\">text</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span><span class=\"p\">)</span>\n                <span class=\"k\">elif</span>
    <span class=\"nb\">len</span><span class=\"p\">(</span><span class=\"n\">possible_pages</span><span
    class=\"p\">)</span> <span class=\"o\">&gt;</span> <span class=\"mi\">1</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">logger</span><span
    class=\"o\">.</span><span class=\"n\">warning</span><span class=\"p\">(</span>\n
    \                       <span class=\"sa\">f</span><span class=\"s2\">&quot;wikilink
    [[</span><span class=\"si\">{</span><span class=\"n\">text</span><span class=\"si\">}</span><span
    class=\"s2\">]] (</span><span class=\"si\">{</span><span class=\"n\">link</span><span
    class=\"si\">}</span><span class=\"s2\">, </span><span class=\"si\">{</span><span
    class=\"nb\">id</span><span class=\"si\">}</span><span class=\"s2\">) has duplicate
    matches, defaulting to the first&quot;</span><span class=\"p\">,</span>\n                    <span
    class=\"p\">)</span>\n                    <span class=\"n\">link</span> <span
    class=\"o\">=</span> <span class=\"n\">possible_pages</span><span class=\"p\">[</span><span
    class=\"mi\">0</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">,</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;/</span><span
    class=\"si\">{</span><span class=\"n\">text</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span><span class=\"p\">)</span>\n                <span class=\"k\">else</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">logger</span><span
    class=\"o\">.</span><span class=\"n\">warning</span><span class=\"p\">(</span>\n
    \                       <span class=\"sa\">f</span><span class=\"s2\">&quot;wikilink
    [[</span><span class=\"si\">{</span><span class=\"n\">text</span><span class=\"si\">}</span><span
    class=\"s2\">]] (</span><span class=\"si\">{</span><span class=\"n\">link</span><span
    class=\"si\">}</span><span class=\"s2\">, </span><span class=\"si\">{</span><span
    class=\"nb\">id</span><span class=\"si\">}</span><span class=\"s2\">) no matches,
    defaulting to &#39;/</span><span class=\"si\">{</span><span class=\"n\">text</span><span
    class=\"si\">}</span><span class=\"s2\">&#39;&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"p\">)</span>\n                    <span class=\"n\">link</span>
    <span class=\"o\">=</span> <span class=\"n\">text</span>\n\n                <span
    class=\"k\">if</span> <span class=\"nb\">id</span> <span class=\"ow\">and</span>
    <span class=\"ow\">not</span> <span class=\"n\">link</span><span class=\"o\">.</span><span
    class=\"n\">endswith</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;#</span><span class=\"si\">{</span><span class=\"nb\">id</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">):</span>\n
    \                   <span class=\"n\">link</span> <span class=\"o\">=</span> <span
    class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">link</span><span class=\"si\">}</span><span class=\"s2\">#</span><span
    class=\"si\">{</span><span class=\"nb\">id</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span>\n\n                <span class=\"n\">token</span><span
    class=\"o\">.</span><span class=\"n\">attrSet</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;href&quot;</span><span class=\"p\">,</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;/</span><span class=\"si\">{</span><span class=\"n\">link</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"n\">content_token</span> <span class=\"o\">=</span>
    <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">push</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;text&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;&quot;</span><span class=\"p\">,</span> <span class=\"mi\">0</span><span
    class=\"p\">)</span>\n                <span class=\"n\">content_token</span><span
    class=\"o\">.</span><span class=\"n\">content</span> <span class=\"o\">=</span>
    <span class=\"n\">text</span>\n\n                <span class=\"n\">token</span>
    <span class=\"o\">=</span> <span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">push</span><span class=\"p\">(</span><span class=\"s2\">&quot;link_close&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;a&quot;</span><span class=\"p\">,</span>
    <span class=\"o\">-</span><span class=\"mi\">1</span><span class=\"p\">)</span>\n
    \               <span class=\"n\">token</span><span class=\"o\">.</span><span
    class=\"n\">content</span> <span class=\"o\">=</span> <span class=\"n\">text</span>\n\n
    \               <span class=\"k\">return</span> <span class=\"kc\">True</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>"
  raw.md: ''
published: true
slug: markata/plugins/md-it-wikilinks
title: Md_It_Wikilinks.Py


---

Wikilinks depicted with `\[[wikilinks]]` can be enabled for sites using the
`markdown-it-py` plugin markata will look through all posts matching up the
file stem to the wiki link and inserting the slug as the href.

???+ note normal wikilink

    Here is a link to a markdown file `docs/nav.md`, and the url becomes
    [/nav](/nav).

    ```md
    [[nav]]
    ```

    When rendered out this wikilink will become an anchor link.

    ```html
    <a class="wikilink" href="/nav">load</a>
    ```

    > This behaves just like the standard wikilink

A special feature that markata brings is slug lookup.  It is able to not only
blindly link to the route specified, but will look up the slug of an article.


???+ note markata slug lookup
    Markata has a load plugin that is generated with the [[docs]] plugin.  It's
    filepath is `markata/plugins/load.py`, so it can be referenced by the file
    stem `load`.

    ```md
    [[load]]
    ```

    Markata will look up the article by the file stem, grab the first article,
    and use its slug as the href.  This turns it into an anchor link that looks
    like this.

    ```html
    <a class="wikilink" href="/markata/plugins/load">load</a>
    ```


!! function <h2 id='wikilinks_plugin' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>wikilinks_plugin <em class='small'>function</em></h2>
    A plugin to create wikilinks tokens.
    These, token should be handled by the renderer.

    ???+ example

        ```md title=markdown
        [[nav]]
        ```

        ```html title=html
        <a class="wikilink" href="/nav">load</a>
        ```
???+ source "wikilinks_plugin <em class='small'>source</em>"

```python

        def wikilinks_plugin(
            md: MarkdownIt,
            start_delimiter: str = "[",
            end_delimiter: str = "]",
            markata=None,
        ):
            """A plugin to create wikilinks tokens.
            These, token should be handled by the renderer.

            ???+ example

                ```md title=markdown
                [[nav]]
                ```

                ```html title=html
                <a class="wikilink" href="/nav">load</a>
                ```
            """

            start_char = ord(start_delimiter)
            end_char = ord(end_delimiter)

            def _wikilinks_inline(state: StateInline, silent: bool):
                try:
                    if (
                        state.srcCharCode[state.pos] != start_char
                        or state.srcCharCode[state.pos + 1] != start_char
                    ):
                        return False
                except IndexError:
                    return False

                pos = state.pos + 2
                found_closing = False
                while True:
                    try:
                        end = state.srcCharCode.index(end_char, pos)
                    except ValueError:
                        return False
                    try:
                        if state.srcCharCode[end + 1] == end_char:
                            found_closing = True
                            break
                    except IndexError:
                        return False
                    pos = end + 2

                if not found_closing:
                    return False

                text = state.src[state.pos + 2 : end].strip()
                state.pos = end + 2

                if silent:
                    return True

                token = state.push("link_open", "a", 1)
                token.block = False
                token.attrSet("class", "wikilink")
                if "#" in text:
                    link, id = text.split("#")
                    link = link.strip("/")
                else:
                    link, id = text, None
                possible_pages = markata.filter(
                    f'str(path).split("/")[-1].split(".")[0].replace("_", "-") == "{link.replace("_", "-")}"',
                )
                if len(possible_pages) == 1:
                    link = possible_pages[0].get("slug", f"/{text}")
                elif len(possible_pages) > 1:
                    logger.warning(
                        f"wikilink [[{text}]] ({link}, {id}) has duplicate matches, defaulting to the first",
                    )
                    link = possible_pages[0].get("slug", f"/{text}")
                else:
                    logger.warning(
                        f"wikilink [[{text}]] ({link}, {id}) no matches, defaulting to '/{text}'",
                    )
                    link = text

                if id and not link.endswith(f"#{id}"):
                    link = f"{link}#{id}"

                token.attrSet("href", f"/{link}")
                content_token = state.push("text", "", 0)
                content_token.content = text

                token = state.push("link_close", "a", -1)
                token.content = text

                return True

            md.inline.ruler.before("escape", "wikilinks_inline", _wikilinks_inline)
```


!! function <h2 id='_wikilinks_inline' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_wikilinks_inline <em class='small'>function</em></h2>

???+ source "_wikilinks_inline <em class='small'>source</em>"

```python

        def _wikilinks_inline(state: StateInline, silent: bool):
                try:
                    if (
                        state.srcCharCode[state.pos] != start_char
                        or state.srcCharCode[state.pos + 1] != start_char
                    ):
                        return False
                except IndexError:
                    return False

                pos = state.pos + 2
                found_closing = False
                while True:
                    try:
                        end = state.srcCharCode.index(end_char, pos)
                    except ValueError:
                        return False
                    try:
                        if state.srcCharCode[end + 1] == end_char:
                            found_closing = True
                            break
                    except IndexError:
                        return False
                    pos = end + 2

                if not found_closing:
                    return False

                text = state.src[state.pos + 2 : end].strip()
                state.pos = end + 2

                if silent:
                    return True

                token = state.push("link_open", "a", 1)
                token.block = False
                token.attrSet("class", "wikilink")
                if "#" in text:
                    link, id = text.split("#")
                    link = link.strip("/")
                else:
                    link, id = text, None
                possible_pages = markata.filter(
                    f'str(path).split("/")[-1].split(".")[0].replace("_", "-") == "{link.replace("_", "-")}"',
                )
                if len(possible_pages) == 1:
                    link = possible_pages[0].get("slug", f"/{text}")
                elif len(possible_pages) > 1:
                    logger.warning(
                        f"wikilink [[{text}]] ({link}, {id}) has duplicate matches, defaulting to the first",
                    )
                    link = possible_pages[0].get("slug", f"/{text}")
                else:
                    logger.warning(
                        f"wikilink [[{text}]] ({link}, {id}) no matches, defaulting to '/{text}'",
                    )
                    link = text

                if id and not link.endswith(f"#{id}"):
                    link = f"{link}#{id}"

                token.attrSet("href", f"/{link}")
                content_token = state.push("text", "", 0)
                content_token.content = text

                token = state.push("link_close", "a", -1)
                token.content = text

                return True
```

