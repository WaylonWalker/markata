---
content: "None\n\n\n!! function <h2 id='get_tag' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>get_tag <em class='small'>function</em></h2>\n\n???+ source \"get_tag <em
  class='small'>source</em>\"\n\n```python\n\n        def get_tag(params: str) ->
  Tuple[str, str, bool]:\n            open = False\n            if not params.strip():\n
  \               return \"\", \"\", open\n\n            if params.strip().startswith(\"+\"):\n
  \               params = params.strip(\"+\")\n                open = True\n\n            tag,
  *_title = params.strip().split(\" \")\n            joined = \" \".join(_title).strip('\"').strip(\"'\")\n\n
  \           title = \"\"\n            if not joined:\n                title = tag.title()\n
  \           elif joined != '\"\"':\n                title = joined\n            return
  (tag.lower(), title, open)\n```\n\n\n!! function <h2 id='validate' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>validate <em class='small'>function</em></h2>\n\n???+
  source \"validate <em class='small'>source</em>\"\n\n```python\n\n        def validate(params:
  str) -> bool:\n            tag = params.strip().split(\" \", 1)[-1] or \"\"\n            return
  bool(tag)\n```\n\n\n!! function <h2 id='details' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>details <em class='small'>function</em></h2>\n\n???+ source \"details <em
  class='small'>source</em>\"\n\n```python\n\n        def details(state: StateBlock,
  startLine: int, endLine: int, silent: bool) -> bool:\n            start = state.bMarks[startLine]
  + state.tShift[startLine]\n            maximum = state.eMarks[startLine]\n\n            #
  Check out the first character quickly, which should filter out most of non-containers\n
  \           if ord(state.src[start]) != MARKER_CHAR:\n                return False\n\n
  \           # Check out the rest of the marker string\n            pos = start +
  1\n            while pos <= maximum and MARKER_STR[(pos - start) % MARKER_LEN] ==
  state.src[pos]:\n                pos += 1\n\n            marker_count = math.floor((pos
  - start) / MARKER_LEN)\n            if marker_count < MIN_MARKERS:\n                return
  False\n            marker_pos = pos - ((pos - start) % MARKER_LEN)\n            params
  = state.src[marker_pos:maximum]\n            markup = state.src[start:marker_pos]\n\n
  \           if not validate(params):\n                return False\n\n            #
  Since start is found, we can report success here in validation mode\n            if
  silent:\n                return True\n\n            old_parent = state.parentType\n
  \           old_line_max = state.lineMax\n            old_indent = state.blkIndent\n\n
  \           blk_start = pos\n            while blk_start < maximum and state.src[blk_start]
  == \" \":\n                blk_start += 1\n\n            state.parentType = \"details\"\n
  \           state.blkIndent += blk_start - start\n\n            was_empty = False\n\n
  \           # Search for the end of the block\n            next_line = startLine\n
  \           while True:\n                next_line += 1\n                if next_line
  >= endLine:\n                    # unclosed block should be autoclosed by end of
  document.\n                    # also block seems to be autoclosed by end of parent\n
  \                   break\n                pos = state.bMarks[next_line] + state.tShift[next_line]\n
  \               maximum = state.eMarks[next_line]\n                is_empty = state.sCount[next_line]
  < state.blkIndent\n\n                # two consecutive empty lines autoclose the
  block\n                if is_empty and was_empty:\n                    break\n                was_empty
  = is_empty\n\n                if pos < maximum and state.sCount[next_line] < state.blkIndent:\n
  \                   # non-empty line with negative indent should stop the block:\n
  \                   # - !!!\n                    #  test\n                    break\n\n
  \           # this will prevent lazy continuations from ever going past our end
  marker\n            state.lineMax = next_line\n\n            tag, title, open =
  get_tag(params)\n            attrs = {\"class\": f\"details {tag}\"}\n            if
  open:\n                attrs[\"open\"] = \"\"\n\n            token = state.push(\"details_open\",
  \"details\", 1)\n            token.markup = markup\n            token.block = True\n
  \           token.attrs = attrs\n            token.meta = {\"tag\": tag}\n            token.content
  = title\n            token.info = params\n            token.map = [startLine, next_line]\n\n
  \           if title:\n                title_markup = f\"{markup} {tag}\"\n                token
  = state.push(\"details_title_open\", \"summary\", 1)\n                token.markup
  = title_markup\n                token.attrs = {\"class\": \"admonition-title\"}\n
  \               token.map = [startLine, startLine + 1]\n\n                token
  = state.push(\"inline\", \"\", 0)\n                token.content = title\n                token.map
  = [startLine, startLine + 1]\n                token.children = []\n\n                token
  = state.push(\"details_title_close\", \"summary\", -1)\n                token.markup
  = title_markup\n\n            state.md.block.tokenize(state, startLine + 1, next_line)\n\n
  \           token = state.push(\"details_close\", \"details\", -1)\n            token.markup
  = state.src[start:pos]\n            token.block = True\n\n            state.parentType
  = old_parent\n            state.lineMax = old_line_max\n            state.blkIndent
  = old_indent\n            state.line = next_line\n\n            return True\n```\n\n\n!!
  function <h2 id='details_plugin' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>details_plugin <em class='small'>function</em></h2>\n    Plugin to use\n
  \   `python-markdown style detailss\n    <https://python-markdown.github.io/extensions/details>`_.\n\n
  \   .. code-block:: md\n\n        ??? note\n            *content*\n\n    Note, this
  is ported from\n    `markdown-it-admon\n    <https://github.com/commenthol/markdown-it-admon>`_.\n???+
  source \"details_plugin <em class='small'>source</em>\"\n\n```python\n\n        def
  details_plugin(md: MarkdownIt, render: Optional[Callable] = None) -> None:\n            \"\"\"Plugin
  to use\n            `python-markdown style detailss\n            <https://python-markdown.github.io/extensions/details>`_.\n\n
  \           .. code-block:: md\n\n                ??? note\n                    *content*\n\n
  \           Note, this is ported from\n            `markdown-it-admon\n            <https://github.com/commenthol/markdown-it-admon>`_.\n
  \           \"\"\"\n\n            def renderDefault(self, tokens, idx, _options,
  env):\n                return self.renderToken(tokens, idx, _options, env)\n\n            render
  = render or renderDefault\n\n            md.add_render_rule(\"details_open\", render)\n
  \           md.add_render_rule(\"details_close\", render)\n            md.add_render_rule(\"details_title_open\",
  render)\n            md.add_render_rule(\"details_title_close\", render)\n\n            md.block.ruler.before(\n
  \               \"fence\",\n                \"details\",\n                details,\n
  \               {\"alt\": [\"paragraph\", \"reference\", \"blockquote\", \"list\"]},\n
  \           )\n```\n\n\n!! function <h2 id='renderDefault' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>renderDefault <em class='small'>function</em></h2>\n\n???+
  source \"renderDefault <em class='small'>source</em>\"\n\n```python\n\n        def
  renderDefault(self, tokens, idx, _options, env):\n                return self.renderToken(tokens,
  idx, _options, env)\n```\n\n"
date: 0001-01-01
description: 'None ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
  source '
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Mdit_Details.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None ! ???+ source  ! ???+ source  !
    ???+ source  ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Mdit_Details.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None ! ???+ source  ! ???+ source  !
    ???+ source  ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\"
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Mdit_Details.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <p>None</p>\n<p>!!
    function <h2 id='get_tag' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>get_tag <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">get_tag
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
    <span class=\"nf\">get_tag</span><span class=\"p\">(</span><span class=\"n\">params</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Tuple</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"nb\">bool</span><span class=\"p\">]:</span>\n
    \           <span class=\"nb\">open</span> <span class=\"o\">=</span> <span class=\"kc\">False</span>\n
    \           <span class=\"k\">if</span> <span class=\"ow\">not</span> <span class=\"n\">params</span><span
    class=\"o\">.</span><span class=\"n\">strip</span><span class=\"p\">():</span>\n
    \               <span class=\"k\">return</span> <span class=\"s2\">&quot;&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;&quot;</span><span class=\"p\">,</span>
    <span class=\"nb\">open</span>\n\n            <span class=\"k\">if</span> <span
    class=\"n\">params</span><span class=\"o\">.</span><span class=\"n\">strip</span><span
    class=\"p\">()</span><span class=\"o\">.</span><span class=\"n\">startswith</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;+&quot;</span><span class=\"p\">):</span>\n
    \               <span class=\"n\">params</span> <span class=\"o\">=</span> <span
    class=\"n\">params</span><span class=\"o\">.</span><span class=\"n\">strip</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;+&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"nb\">open</span> <span class=\"o\">=</span> <span
    class=\"kc\">True</span>\n\n            <span class=\"n\">tag</span><span class=\"p\">,</span>
    <span class=\"o\">*</span><span class=\"n\">_title</span> <span class=\"o\">=</span>
    <span class=\"n\">params</span><span class=\"o\">.</span><span class=\"n\">strip</span><span
    class=\"p\">()</span><span class=\"o\">.</span><span class=\"n\">split</span><span
    class=\"p\">(</span><span class=\"s2\">&quot; &quot;</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">joined</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;
    &quot;</span><span class=\"o\">.</span><span class=\"n\">join</span><span class=\"p\">(</span><span
    class=\"n\">_title</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">strip</span><span class=\"p\">(</span><span class=\"s1\">&#39;&quot;&#39;</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">strip</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;&#39;&quot;</span><span class=\"p\">)</span>\n\n
    \           <span class=\"n\">title</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;&quot;</span>\n
    \           <span class=\"k\">if</span> <span class=\"ow\">not</span> <span class=\"n\">joined</span><span
    class=\"p\">:</span>\n                <span class=\"n\">title</span> <span class=\"o\">=</span>
    <span class=\"n\">tag</span><span class=\"o\">.</span><span class=\"n\">title</span><span
    class=\"p\">()</span>\n            <span class=\"k\">elif</span> <span class=\"n\">joined</span>
    <span class=\"o\">!=</span> <span class=\"s1\">&#39;&quot;&quot;&#39;</span><span
    class=\"p\">:</span>\n                <span class=\"n\">title</span> <span class=\"o\">=</span>
    <span class=\"n\">joined</span>\n            <span class=\"k\">return</span> <span
    class=\"p\">(</span><span class=\"n\">tag</span><span class=\"o\">.</span><span
    class=\"n\">lower</span><span class=\"p\">(),</span> <span class=\"n\">title</span><span
    class=\"p\">,</span> <span class=\"nb\">open</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='validate' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>validate <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">validate
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
    <span class=\"nf\">validate</span><span class=\"p\">(</span><span class=\"n\">params</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"nb\">bool</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">tag</span> <span class=\"o\">=</span> <span class=\"n\">params</span><span
    class=\"o\">.</span><span class=\"n\">strip</span><span class=\"p\">()</span><span
    class=\"o\">.</span><span class=\"n\">split</span><span class=\"p\">(</span><span
    class=\"s2\">&quot; &quot;</span><span class=\"p\">,</span> <span class=\"mi\">1</span><span
    class=\"p\">)[</span><span class=\"o\">-</span><span class=\"mi\">1</span><span
    class=\"p\">]</span> <span class=\"ow\">or</span> <span class=\"s2\">&quot;&quot;</span>\n
    \           <span class=\"k\">return</span> <span class=\"nb\">bool</span><span
    class=\"p\">(</span><span class=\"n\">tag</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='details' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>details <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">details
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
    <span class=\"nf\">details</span><span class=\"p\">(</span><span class=\"n\">state</span><span
    class=\"p\">:</span> <span class=\"n\">StateBlock</span><span class=\"p\">,</span>
    <span class=\"n\">startLine</span><span class=\"p\">:</span> <span class=\"nb\">int</span><span
    class=\"p\">,</span> <span class=\"n\">endLine</span><span class=\"p\">:</span>
    <span class=\"nb\">int</span><span class=\"p\">,</span> <span class=\"n\">silent</span><span
    class=\"p\">:</span> <span class=\"nb\">bool</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">bool</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">start</span> <span class=\"o\">=</span> <span class=\"n\">state</span><span
    class=\"o\">.</span><span class=\"n\">bMarks</span><span class=\"p\">[</span><span
    class=\"n\">startLine</span><span class=\"p\">]</span> <span class=\"o\">+</span>
    <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">tShift</span><span
    class=\"p\">[</span><span class=\"n\">startLine</span><span class=\"p\">]</span>\n
    \           <span class=\"n\">maximum</span> <span class=\"o\">=</span> <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">eMarks</span><span
    class=\"p\">[</span><span class=\"n\">startLine</span><span class=\"p\">]</span>\n\n
    \           <span class=\"c1\"># Check out the first character quickly, which
    should filter out most of non-containers</span>\n            <span class=\"k\">if</span>
    <span class=\"nb\">ord</span><span class=\"p\">(</span><span class=\"n\">state</span><span
    class=\"o\">.</span><span class=\"n\">src</span><span class=\"p\">[</span><span
    class=\"n\">start</span><span class=\"p\">])</span> <span class=\"o\">!=</span>
    <span class=\"n\">MARKER_CHAR</span><span class=\"p\">:</span>\n                <span
    class=\"k\">return</span> <span class=\"kc\">False</span>\n\n            <span
    class=\"c1\"># Check out the rest of the marker string</span>\n            <span
    class=\"n\">pos</span> <span class=\"o\">=</span> <span class=\"n\">start</span>
    <span class=\"o\">+</span> <span class=\"mi\">1</span>\n            <span class=\"k\">while</span>
    <span class=\"n\">pos</span> <span class=\"o\">&lt;=</span> <span class=\"n\">maximum</span>
    <span class=\"ow\">and</span> <span class=\"n\">MARKER_STR</span><span class=\"p\">[(</span><span
    class=\"n\">pos</span> <span class=\"o\">-</span> <span class=\"n\">start</span><span
    class=\"p\">)</span> <span class=\"o\">%</span> <span class=\"n\">MARKER_LEN</span><span
    class=\"p\">]</span> <span class=\"o\">==</span> <span class=\"n\">state</span><span
    class=\"o\">.</span><span class=\"n\">src</span><span class=\"p\">[</span><span
    class=\"n\">pos</span><span class=\"p\">]:</span>\n                <span class=\"n\">pos</span>
    <span class=\"o\">+=</span> <span class=\"mi\">1</span>\n\n            <span class=\"n\">marker_count</span>
    <span class=\"o\">=</span> <span class=\"n\">math</span><span class=\"o\">.</span><span
    class=\"n\">floor</span><span class=\"p\">((</span><span class=\"n\">pos</span>
    <span class=\"o\">-</span> <span class=\"n\">start</span><span class=\"p\">)</span>
    <span class=\"o\">/</span> <span class=\"n\">MARKER_LEN</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">if</span> <span class=\"n\">marker_count</span>
    <span class=\"o\">&lt;</span> <span class=\"n\">MIN_MARKERS</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">return</span> <span class=\"kc\">False</span>\n
    \           <span class=\"n\">marker_pos</span> <span class=\"o\">=</span> <span
    class=\"n\">pos</span> <span class=\"o\">-</span> <span class=\"p\">((</span><span
    class=\"n\">pos</span> <span class=\"o\">-</span> <span class=\"n\">start</span><span
    class=\"p\">)</span> <span class=\"o\">%</span> <span class=\"n\">MARKER_LEN</span><span
    class=\"p\">)</span>\n            <span class=\"n\">params</span> <span class=\"o\">=</span>
    <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">src</span><span
    class=\"p\">[</span><span class=\"n\">marker_pos</span><span class=\"p\">:</span><span
    class=\"n\">maximum</span><span class=\"p\">]</span>\n            <span class=\"n\">markup</span>
    <span class=\"o\">=</span> <span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">src</span><span class=\"p\">[</span><span class=\"n\">start</span><span
    class=\"p\">:</span><span class=\"n\">marker_pos</span><span class=\"p\">]</span>\n\n
    \           <span class=\"k\">if</span> <span class=\"ow\">not</span> <span class=\"n\">validate</span><span
    class=\"p\">(</span><span class=\"n\">params</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">return</span> <span class=\"kc\">False</span>\n\n
    \           <span class=\"c1\"># Since start is found, we can report success here
    in validation mode</span>\n            <span class=\"k\">if</span> <span class=\"n\">silent</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"kc\">True</span>\n\n
    \           <span class=\"n\">old_parent</span> <span class=\"o\">=</span> <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">parentType</span>\n
    \           <span class=\"n\">old_line_max</span> <span class=\"o\">=</span> <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">lineMax</span>\n
    \           <span class=\"n\">old_indent</span> <span class=\"o\">=</span> <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">blkIndent</span>\n\n
    \           <span class=\"n\">blk_start</span> <span class=\"o\">=</span> <span
    class=\"n\">pos</span>\n            <span class=\"k\">while</span> <span class=\"n\">blk_start</span>
    <span class=\"o\">&lt;</span> <span class=\"n\">maximum</span> <span class=\"ow\">and</span>
    <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">src</span><span
    class=\"p\">[</span><span class=\"n\">blk_start</span><span class=\"p\">]</span>
    <span class=\"o\">==</span> <span class=\"s2\">&quot; &quot;</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">blk_start</span> <span class=\"o\">+=</span>
    <span class=\"mi\">1</span>\n\n            <span class=\"n\">state</span><span
    class=\"o\">.</span><span class=\"n\">parentType</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;details&quot;</span>\n            <span class=\"n\">state</span><span
    class=\"o\">.</span><span class=\"n\">blkIndent</span> <span class=\"o\">+=</span>
    <span class=\"n\">blk_start</span> <span class=\"o\">-</span> <span class=\"n\">start</span>\n\n
    \           <span class=\"n\">was_empty</span> <span class=\"o\">=</span> <span
    class=\"kc\">False</span>\n\n            <span class=\"c1\"># Search for the end
    of the block</span>\n            <span class=\"n\">next_line</span> <span class=\"o\">=</span>
    <span class=\"n\">startLine</span>\n            <span class=\"k\">while</span>
    <span class=\"kc\">True</span><span class=\"p\">:</span>\n                <span
    class=\"n\">next_line</span> <span class=\"o\">+=</span> <span class=\"mi\">1</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">next_line</span>
    <span class=\"o\">&gt;=</span> <span class=\"n\">endLine</span><span class=\"p\">:</span>\n
    \                   <span class=\"c1\"># unclosed block should be autoclosed by
    end of document.</span>\n                    <span class=\"c1\"># also block seems
    to be autoclosed by end of parent</span>\n                    <span class=\"k\">break</span>\n
    \               <span class=\"n\">pos</span> <span class=\"o\">=</span> <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">bMarks</span><span
    class=\"p\">[</span><span class=\"n\">next_line</span><span class=\"p\">]</span>
    <span class=\"o\">+</span> <span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">tShift</span><span class=\"p\">[</span><span class=\"n\">next_line</span><span
    class=\"p\">]</span>\n                <span class=\"n\">maximum</span> <span class=\"o\">=</span>
    <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">eMarks</span><span
    class=\"p\">[</span><span class=\"n\">next_line</span><span class=\"p\">]</span>\n
    \               <span class=\"n\">is_empty</span> <span class=\"o\">=</span> <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">sCount</span><span
    class=\"p\">[</span><span class=\"n\">next_line</span><span class=\"p\">]</span>
    <span class=\"o\">&lt;</span> <span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">blkIndent</span>\n\n                <span class=\"c1\"># two consecutive
    empty lines autoclose the block</span>\n                <span class=\"k\">if</span>
    <span class=\"n\">is_empty</span> <span class=\"ow\">and</span> <span class=\"n\">was_empty</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">break</span>\n                <span
    class=\"n\">was_empty</span> <span class=\"o\">=</span> <span class=\"n\">is_empty</span>\n\n
    \               <span class=\"k\">if</span> <span class=\"n\">pos</span> <span
    class=\"o\">&lt;</span> <span class=\"n\">maximum</span> <span class=\"ow\">and</span>
    <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">sCount</span><span
    class=\"p\">[</span><span class=\"n\">next_line</span><span class=\"p\">]</span>
    <span class=\"o\">&lt;</span> <span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">blkIndent</span><span class=\"p\">:</span>\n                    <span
    class=\"c1\"># non-empty line with negative indent should stop the block:</span>\n
    \                   <span class=\"c1\"># - !!!</span>\n                    <span
    class=\"c1\">#  test</span>\n                    <span class=\"k\">break</span>\n\n
    \           <span class=\"c1\"># this will prevent lazy continuations from ever
    going past our end marker</span>\n            <span class=\"n\">state</span><span
    class=\"o\">.</span><span class=\"n\">lineMax</span> <span class=\"o\">=</span>
    <span class=\"n\">next_line</span>\n\n            <span class=\"n\">tag</span><span
    class=\"p\">,</span> <span class=\"n\">title</span><span class=\"p\">,</span>
    <span class=\"nb\">open</span> <span class=\"o\">=</span> <span class=\"n\">get_tag</span><span
    class=\"p\">(</span><span class=\"n\">params</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">attrs</span> <span class=\"o\">=</span> <span class=\"p\">{</span><span
    class=\"s2\">&quot;class&quot;</span><span class=\"p\">:</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;details </span><span class=\"si\">{</span><span class=\"n\">tag</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">}</span>\n
    \           <span class=\"k\">if</span> <span class=\"nb\">open</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">attrs</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;open&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;&quot;</span>\n\n            <span class=\"n\">token</span>
    <span class=\"o\">=</span> <span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">push</span><span class=\"p\">(</span><span class=\"s2\">&quot;details_open&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;details&quot;</span><span class=\"p\">,</span>
    <span class=\"mi\">1</span><span class=\"p\">)</span>\n            <span class=\"n\">token</span><span
    class=\"o\">.</span><span class=\"n\">markup</span> <span class=\"o\">=</span>
    <span class=\"n\">markup</span>\n            <span class=\"n\">token</span><span
    class=\"o\">.</span><span class=\"n\">block</span> <span class=\"o\">=</span>
    <span class=\"kc\">True</span>\n            <span class=\"n\">token</span><span
    class=\"o\">.</span><span class=\"n\">attrs</span> <span class=\"o\">=</span>
    <span class=\"n\">attrs</span>\n            <span class=\"n\">token</span><span
    class=\"o\">.</span><span class=\"n\">meta</span> <span class=\"o\">=</span> <span
    class=\"p\">{</span><span class=\"s2\">&quot;tag&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">tag</span><span class=\"p\">}</span>\n            <span class=\"n\">token</span><span
    class=\"o\">.</span><span class=\"n\">content</span> <span class=\"o\">=</span>
    <span class=\"n\">title</span>\n            <span class=\"n\">token</span><span
    class=\"o\">.</span><span class=\"n\">info</span> <span class=\"o\">=</span> <span
    class=\"n\">params</span>\n            <span class=\"n\">token</span><span class=\"o\">.</span><span
    class=\"n\">map</span> <span class=\"o\">=</span> <span class=\"p\">[</span><span
    class=\"n\">startLine</span><span class=\"p\">,</span> <span class=\"n\">next_line</span><span
    class=\"p\">]</span>\n\n            <span class=\"k\">if</span> <span class=\"n\">title</span><span
    class=\"p\">:</span>\n                <span class=\"n\">title_markup</span> <span
    class=\"o\">=</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">markup</span><span class=\"si\">}</span><span
    class=\"s2\"> </span><span class=\"si\">{</span><span class=\"n\">tag</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span>\n                <span class=\"n\">token</span>
    <span class=\"o\">=</span> <span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">push</span><span class=\"p\">(</span><span class=\"s2\">&quot;details_title_open&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;summary&quot;</span><span class=\"p\">,</span>
    <span class=\"mi\">1</span><span class=\"p\">)</span>\n                <span class=\"n\">token</span><span
    class=\"o\">.</span><span class=\"n\">markup</span> <span class=\"o\">=</span>
    <span class=\"n\">title_markup</span>\n                <span class=\"n\">token</span><span
    class=\"o\">.</span><span class=\"n\">attrs</span> <span class=\"o\">=</span>
    <span class=\"p\">{</span><span class=\"s2\">&quot;class&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;admonition-title&quot;</span><span class=\"p\">}</span>\n
    \               <span class=\"n\">token</span><span class=\"o\">.</span><span
    class=\"n\">map</span> <span class=\"o\">=</span> <span class=\"p\">[</span><span
    class=\"n\">startLine</span><span class=\"p\">,</span> <span class=\"n\">startLine</span>
    <span class=\"o\">+</span> <span class=\"mi\">1</span><span class=\"p\">]</span>\n\n
    \               <span class=\"n\">token</span> <span class=\"o\">=</span> <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">push</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;inline&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;&quot;</span><span class=\"p\">,</span> <span class=\"mi\">0</span><span
    class=\"p\">)</span>\n                <span class=\"n\">token</span><span class=\"o\">.</span><span
    class=\"n\">content</span> <span class=\"o\">=</span> <span class=\"n\">title</span>\n
    \               <span class=\"n\">token</span><span class=\"o\">.</span><span
    class=\"n\">map</span> <span class=\"o\">=</span> <span class=\"p\">[</span><span
    class=\"n\">startLine</span><span class=\"p\">,</span> <span class=\"n\">startLine</span>
    <span class=\"o\">+</span> <span class=\"mi\">1</span><span class=\"p\">]</span>\n
    \               <span class=\"n\">token</span><span class=\"o\">.</span><span
    class=\"n\">children</span> <span class=\"o\">=</span> <span class=\"p\">[]</span>\n\n
    \               <span class=\"n\">token</span> <span class=\"o\">=</span> <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">push</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;details_title_close&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;summary&quot;</span><span class=\"p\">,</span>
    <span class=\"o\">-</span><span class=\"mi\">1</span><span class=\"p\">)</span>\n
    \               <span class=\"n\">token</span><span class=\"o\">.</span><span
    class=\"n\">markup</span> <span class=\"o\">=</span> <span class=\"n\">title_markup</span>\n\n
    \           <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">md</span><span
    class=\"o\">.</span><span class=\"n\">block</span><span class=\"o\">.</span><span
    class=\"n\">tokenize</span><span class=\"p\">(</span><span class=\"n\">state</span><span
    class=\"p\">,</span> <span class=\"n\">startLine</span> <span class=\"o\">+</span>
    <span class=\"mi\">1</span><span class=\"p\">,</span> <span class=\"n\">next_line</span><span
    class=\"p\">)</span>\n\n            <span class=\"n\">token</span> <span class=\"o\">=</span>
    <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">push</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;details_close&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;details&quot;</span><span class=\"p\">,</span> <span
    class=\"o\">-</span><span class=\"mi\">1</span><span class=\"p\">)</span>\n            <span
    class=\"n\">token</span><span class=\"o\">.</span><span class=\"n\">markup</span>
    <span class=\"o\">=</span> <span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">src</span><span class=\"p\">[</span><span class=\"n\">start</span><span
    class=\"p\">:</span><span class=\"n\">pos</span><span class=\"p\">]</span>\n            <span
    class=\"n\">token</span><span class=\"o\">.</span><span class=\"n\">block</span>
    <span class=\"o\">=</span> <span class=\"kc\">True</span>\n\n            <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">parentType</span>
    <span class=\"o\">=</span> <span class=\"n\">old_parent</span>\n            <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">lineMax</span>
    <span class=\"o\">=</span> <span class=\"n\">old_line_max</span>\n            <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">blkIndent</span>
    <span class=\"o\">=</span> <span class=\"n\">old_indent</span>\n            <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">line</span>
    <span class=\"o\">=</span> <span class=\"n\">next_line</span>\n\n            <span
    class=\"k\">return</span> <span class=\"kc\">True</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='details_plugin' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>details_plugin <em class='small'>function</em></h2>\nPlugin to use\n<code>python-markdown
    style detailss     &lt;https://python-markdown.github.io/extensions/details&gt;</code>_.</p>\n<pre><code>..
    code-block:: md\n\n    ??? note\n        *content*\n\nNote, this is ported from\n`markdown-it-admon\n&lt;https://github.com/commenthol/markdown-it-admon&gt;`_.\n</code></pre>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">details_plugin
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
    <span class=\"nf\">details_plugin</span><span class=\"p\">(</span><span class=\"n\">md</span><span
    class=\"p\">:</span> <span class=\"n\">MarkdownIt</span><span class=\"p\">,</span>
    <span class=\"n\">render</span><span class=\"p\">:</span> <span class=\"n\">Optional</span><span
    class=\"p\">[</span><span class=\"n\">Callable</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"kc\">None</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;Plugin to
    use</span>\n<span class=\"sd\">            `python-markdown style detailss</span>\n<span
    class=\"sd\">            &lt;https://python-markdown.github.io/extensions/details&gt;`_.</span>\n\n<span
    class=\"sd\">            .. code-block:: md</span>\n\n<span class=\"sd\">                ???
    note</span>\n<span class=\"sd\">                    *content*</span>\n\n<span
    class=\"sd\">            Note, this is ported from</span>\n<span class=\"sd\">
    \           `markdown-it-admon</span>\n<span class=\"sd\">            &lt;https://github.com/commenthol/markdown-it-admon&gt;`_.</span>\n<span
    class=\"sd\">            &quot;&quot;&quot;</span>\n\n            <span class=\"k\">def</span>
    <span class=\"nf\">renderDefault</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">tokens</span><span class=\"p\">,</span>
    <span class=\"n\">idx</span><span class=\"p\">,</span> <span class=\"n\">_options</span><span
    class=\"p\">,</span> <span class=\"n\">env</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">renderToken</span><span class=\"p\">(</span><span
    class=\"n\">tokens</span><span class=\"p\">,</span> <span class=\"n\">idx</span><span
    class=\"p\">,</span> <span class=\"n\">_options</span><span class=\"p\">,</span>
    <span class=\"n\">env</span><span class=\"p\">)</span>\n\n            <span class=\"n\">render</span>
    <span class=\"o\">=</span> <span class=\"n\">render</span> <span class=\"ow\">or</span>
    <span class=\"n\">renderDefault</span>\n\n            <span class=\"n\">md</span><span
    class=\"o\">.</span><span class=\"n\">add_render_rule</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;details_open&quot;</span><span class=\"p\">,</span> <span class=\"n\">render</span><span
    class=\"p\">)</span>\n            <span class=\"n\">md</span><span class=\"o\">.</span><span
    class=\"n\">add_render_rule</span><span class=\"p\">(</span><span class=\"s2\">&quot;details_close&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">render</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">md</span><span class=\"o\">.</span><span class=\"n\">add_render_rule</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;details_title_open&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">render</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">md</span><span class=\"o\">.</span><span class=\"n\">add_render_rule</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;details_title_close&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">render</span><span class=\"p\">)</span>\n\n
    \           <span class=\"n\">md</span><span class=\"o\">.</span><span class=\"n\">block</span><span
    class=\"o\">.</span><span class=\"n\">ruler</span><span class=\"o\">.</span><span
    class=\"n\">before</span><span class=\"p\">(</span>\n                <span class=\"s2\">&quot;fence&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;details&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">details</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">{</span><span class=\"s2\">&quot;alt&quot;</span><span
    class=\"p\">:</span> <span class=\"p\">[</span><span class=\"s2\">&quot;paragraph&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;reference&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;blockquote&quot;</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;list&quot;</span><span class=\"p\">]},</span>\n            <span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='renderDefault'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>renderDefault <em
    class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">renderDefault <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">renderDefault</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">tokens</span><span class=\"p\">,</span>
    <span class=\"n\">idx</span><span class=\"p\">,</span> <span class=\"n\">_options</span><span
    class=\"p\">,</span> <span class=\"n\">env</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">renderToken</span><span class=\"p\">(</span><span
    class=\"n\">tokens</span><span class=\"p\">,</span> <span class=\"n\">idx</span><span
    class=\"p\">,</span> <span class=\"n\">_options</span><span class=\"p\">,</span>
    <span class=\"n\">env</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>    </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Mdit_Details.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None ! ???+ source  ! ???+ source  !
    ???+ source  ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Mdit_Details.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None ! ???+ source  ! ???+ source  !
    ???+ source  ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<article style=\"text-align: center;\">\n    <style>\n        section
    {\n            font-size: 200%;\n        }\n\n\n        .edit {\n            display:
    none;\n        }\n    </style>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Mdit_Details.Py\n    </h1>\n</section></article>\n     </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Mdit_Details.Py\n    </h1>\n</section>    <section class=\"body\">\n        <p>None</p>\n<p>!!
    function <h2 id='get_tag' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>get_tag <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">get_tag
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
    <span class=\"nf\">get_tag</span><span class=\"p\">(</span><span class=\"n\">params</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Tuple</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"nb\">bool</span><span class=\"p\">]:</span>\n
    \           <span class=\"nb\">open</span> <span class=\"o\">=</span> <span class=\"kc\">False</span>\n
    \           <span class=\"k\">if</span> <span class=\"ow\">not</span> <span class=\"n\">params</span><span
    class=\"o\">.</span><span class=\"n\">strip</span><span class=\"p\">():</span>\n
    \               <span class=\"k\">return</span> <span class=\"s2\">&quot;&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;&quot;</span><span class=\"p\">,</span>
    <span class=\"nb\">open</span>\n\n            <span class=\"k\">if</span> <span
    class=\"n\">params</span><span class=\"o\">.</span><span class=\"n\">strip</span><span
    class=\"p\">()</span><span class=\"o\">.</span><span class=\"n\">startswith</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;+&quot;</span><span class=\"p\">):</span>\n
    \               <span class=\"n\">params</span> <span class=\"o\">=</span> <span
    class=\"n\">params</span><span class=\"o\">.</span><span class=\"n\">strip</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;+&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"nb\">open</span> <span class=\"o\">=</span> <span
    class=\"kc\">True</span>\n\n            <span class=\"n\">tag</span><span class=\"p\">,</span>
    <span class=\"o\">*</span><span class=\"n\">_title</span> <span class=\"o\">=</span>
    <span class=\"n\">params</span><span class=\"o\">.</span><span class=\"n\">strip</span><span
    class=\"p\">()</span><span class=\"o\">.</span><span class=\"n\">split</span><span
    class=\"p\">(</span><span class=\"s2\">&quot; &quot;</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">joined</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;
    &quot;</span><span class=\"o\">.</span><span class=\"n\">join</span><span class=\"p\">(</span><span
    class=\"n\">_title</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">strip</span><span class=\"p\">(</span><span class=\"s1\">&#39;&quot;&#39;</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">strip</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;&#39;&quot;</span><span class=\"p\">)</span>\n\n
    \           <span class=\"n\">title</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;&quot;</span>\n
    \           <span class=\"k\">if</span> <span class=\"ow\">not</span> <span class=\"n\">joined</span><span
    class=\"p\">:</span>\n                <span class=\"n\">title</span> <span class=\"o\">=</span>
    <span class=\"n\">tag</span><span class=\"o\">.</span><span class=\"n\">title</span><span
    class=\"p\">()</span>\n            <span class=\"k\">elif</span> <span class=\"n\">joined</span>
    <span class=\"o\">!=</span> <span class=\"s1\">&#39;&quot;&quot;&#39;</span><span
    class=\"p\">:</span>\n                <span class=\"n\">title</span> <span class=\"o\">=</span>
    <span class=\"n\">joined</span>\n            <span class=\"k\">return</span> <span
    class=\"p\">(</span><span class=\"n\">tag</span><span class=\"o\">.</span><span
    class=\"n\">lower</span><span class=\"p\">(),</span> <span class=\"n\">title</span><span
    class=\"p\">,</span> <span class=\"nb\">open</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='validate' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>validate <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">validate
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
    <span class=\"nf\">validate</span><span class=\"p\">(</span><span class=\"n\">params</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"nb\">bool</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">tag</span> <span class=\"o\">=</span> <span class=\"n\">params</span><span
    class=\"o\">.</span><span class=\"n\">strip</span><span class=\"p\">()</span><span
    class=\"o\">.</span><span class=\"n\">split</span><span class=\"p\">(</span><span
    class=\"s2\">&quot; &quot;</span><span class=\"p\">,</span> <span class=\"mi\">1</span><span
    class=\"p\">)[</span><span class=\"o\">-</span><span class=\"mi\">1</span><span
    class=\"p\">]</span> <span class=\"ow\">or</span> <span class=\"s2\">&quot;&quot;</span>\n
    \           <span class=\"k\">return</span> <span class=\"nb\">bool</span><span
    class=\"p\">(</span><span class=\"n\">tag</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='details' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>details <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">details
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
    <span class=\"nf\">details</span><span class=\"p\">(</span><span class=\"n\">state</span><span
    class=\"p\">:</span> <span class=\"n\">StateBlock</span><span class=\"p\">,</span>
    <span class=\"n\">startLine</span><span class=\"p\">:</span> <span class=\"nb\">int</span><span
    class=\"p\">,</span> <span class=\"n\">endLine</span><span class=\"p\">:</span>
    <span class=\"nb\">int</span><span class=\"p\">,</span> <span class=\"n\">silent</span><span
    class=\"p\">:</span> <span class=\"nb\">bool</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">bool</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">start</span> <span class=\"o\">=</span> <span class=\"n\">state</span><span
    class=\"o\">.</span><span class=\"n\">bMarks</span><span class=\"p\">[</span><span
    class=\"n\">startLine</span><span class=\"p\">]</span> <span class=\"o\">+</span>
    <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">tShift</span><span
    class=\"p\">[</span><span class=\"n\">startLine</span><span class=\"p\">]</span>\n
    \           <span class=\"n\">maximum</span> <span class=\"o\">=</span> <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">eMarks</span><span
    class=\"p\">[</span><span class=\"n\">startLine</span><span class=\"p\">]</span>\n\n
    \           <span class=\"c1\"># Check out the first character quickly, which
    should filter out most of non-containers</span>\n            <span class=\"k\">if</span>
    <span class=\"nb\">ord</span><span class=\"p\">(</span><span class=\"n\">state</span><span
    class=\"o\">.</span><span class=\"n\">src</span><span class=\"p\">[</span><span
    class=\"n\">start</span><span class=\"p\">])</span> <span class=\"o\">!=</span>
    <span class=\"n\">MARKER_CHAR</span><span class=\"p\">:</span>\n                <span
    class=\"k\">return</span> <span class=\"kc\">False</span>\n\n            <span
    class=\"c1\"># Check out the rest of the marker string</span>\n            <span
    class=\"n\">pos</span> <span class=\"o\">=</span> <span class=\"n\">start</span>
    <span class=\"o\">+</span> <span class=\"mi\">1</span>\n            <span class=\"k\">while</span>
    <span class=\"n\">pos</span> <span class=\"o\">&lt;=</span> <span class=\"n\">maximum</span>
    <span class=\"ow\">and</span> <span class=\"n\">MARKER_STR</span><span class=\"p\">[(</span><span
    class=\"n\">pos</span> <span class=\"o\">-</span> <span class=\"n\">start</span><span
    class=\"p\">)</span> <span class=\"o\">%</span> <span class=\"n\">MARKER_LEN</span><span
    class=\"p\">]</span> <span class=\"o\">==</span> <span class=\"n\">state</span><span
    class=\"o\">.</span><span class=\"n\">src</span><span class=\"p\">[</span><span
    class=\"n\">pos</span><span class=\"p\">]:</span>\n                <span class=\"n\">pos</span>
    <span class=\"o\">+=</span> <span class=\"mi\">1</span>\n\n            <span class=\"n\">marker_count</span>
    <span class=\"o\">=</span> <span class=\"n\">math</span><span class=\"o\">.</span><span
    class=\"n\">floor</span><span class=\"p\">((</span><span class=\"n\">pos</span>
    <span class=\"o\">-</span> <span class=\"n\">start</span><span class=\"p\">)</span>
    <span class=\"o\">/</span> <span class=\"n\">MARKER_LEN</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">if</span> <span class=\"n\">marker_count</span>
    <span class=\"o\">&lt;</span> <span class=\"n\">MIN_MARKERS</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">return</span> <span class=\"kc\">False</span>\n
    \           <span class=\"n\">marker_pos</span> <span class=\"o\">=</span> <span
    class=\"n\">pos</span> <span class=\"o\">-</span> <span class=\"p\">((</span><span
    class=\"n\">pos</span> <span class=\"o\">-</span> <span class=\"n\">start</span><span
    class=\"p\">)</span> <span class=\"o\">%</span> <span class=\"n\">MARKER_LEN</span><span
    class=\"p\">)</span>\n            <span class=\"n\">params</span> <span class=\"o\">=</span>
    <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">src</span><span
    class=\"p\">[</span><span class=\"n\">marker_pos</span><span class=\"p\">:</span><span
    class=\"n\">maximum</span><span class=\"p\">]</span>\n            <span class=\"n\">markup</span>
    <span class=\"o\">=</span> <span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">src</span><span class=\"p\">[</span><span class=\"n\">start</span><span
    class=\"p\">:</span><span class=\"n\">marker_pos</span><span class=\"p\">]</span>\n\n
    \           <span class=\"k\">if</span> <span class=\"ow\">not</span> <span class=\"n\">validate</span><span
    class=\"p\">(</span><span class=\"n\">params</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">return</span> <span class=\"kc\">False</span>\n\n
    \           <span class=\"c1\"># Since start is found, we can report success here
    in validation mode</span>\n            <span class=\"k\">if</span> <span class=\"n\">silent</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"kc\">True</span>\n\n
    \           <span class=\"n\">old_parent</span> <span class=\"o\">=</span> <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">parentType</span>\n
    \           <span class=\"n\">old_line_max</span> <span class=\"o\">=</span> <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">lineMax</span>\n
    \           <span class=\"n\">old_indent</span> <span class=\"o\">=</span> <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">blkIndent</span>\n\n
    \           <span class=\"n\">blk_start</span> <span class=\"o\">=</span> <span
    class=\"n\">pos</span>\n            <span class=\"k\">while</span> <span class=\"n\">blk_start</span>
    <span class=\"o\">&lt;</span> <span class=\"n\">maximum</span> <span class=\"ow\">and</span>
    <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">src</span><span
    class=\"p\">[</span><span class=\"n\">blk_start</span><span class=\"p\">]</span>
    <span class=\"o\">==</span> <span class=\"s2\">&quot; &quot;</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">blk_start</span> <span class=\"o\">+=</span>
    <span class=\"mi\">1</span>\n\n            <span class=\"n\">state</span><span
    class=\"o\">.</span><span class=\"n\">parentType</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;details&quot;</span>\n            <span class=\"n\">state</span><span
    class=\"o\">.</span><span class=\"n\">blkIndent</span> <span class=\"o\">+=</span>
    <span class=\"n\">blk_start</span> <span class=\"o\">-</span> <span class=\"n\">start</span>\n\n
    \           <span class=\"n\">was_empty</span> <span class=\"o\">=</span> <span
    class=\"kc\">False</span>\n\n            <span class=\"c1\"># Search for the end
    of the block</span>\n            <span class=\"n\">next_line</span> <span class=\"o\">=</span>
    <span class=\"n\">startLine</span>\n            <span class=\"k\">while</span>
    <span class=\"kc\">True</span><span class=\"p\">:</span>\n                <span
    class=\"n\">next_line</span> <span class=\"o\">+=</span> <span class=\"mi\">1</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">next_line</span>
    <span class=\"o\">&gt;=</span> <span class=\"n\">endLine</span><span class=\"p\">:</span>\n
    \                   <span class=\"c1\"># unclosed block should be autoclosed by
    end of document.</span>\n                    <span class=\"c1\"># also block seems
    to be autoclosed by end of parent</span>\n                    <span class=\"k\">break</span>\n
    \               <span class=\"n\">pos</span> <span class=\"o\">=</span> <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">bMarks</span><span
    class=\"p\">[</span><span class=\"n\">next_line</span><span class=\"p\">]</span>
    <span class=\"o\">+</span> <span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">tShift</span><span class=\"p\">[</span><span class=\"n\">next_line</span><span
    class=\"p\">]</span>\n                <span class=\"n\">maximum</span> <span class=\"o\">=</span>
    <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">eMarks</span><span
    class=\"p\">[</span><span class=\"n\">next_line</span><span class=\"p\">]</span>\n
    \               <span class=\"n\">is_empty</span> <span class=\"o\">=</span> <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">sCount</span><span
    class=\"p\">[</span><span class=\"n\">next_line</span><span class=\"p\">]</span>
    <span class=\"o\">&lt;</span> <span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">blkIndent</span>\n\n                <span class=\"c1\"># two consecutive
    empty lines autoclose the block</span>\n                <span class=\"k\">if</span>
    <span class=\"n\">is_empty</span> <span class=\"ow\">and</span> <span class=\"n\">was_empty</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">break</span>\n                <span
    class=\"n\">was_empty</span> <span class=\"o\">=</span> <span class=\"n\">is_empty</span>\n\n
    \               <span class=\"k\">if</span> <span class=\"n\">pos</span> <span
    class=\"o\">&lt;</span> <span class=\"n\">maximum</span> <span class=\"ow\">and</span>
    <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">sCount</span><span
    class=\"p\">[</span><span class=\"n\">next_line</span><span class=\"p\">]</span>
    <span class=\"o\">&lt;</span> <span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">blkIndent</span><span class=\"p\">:</span>\n                    <span
    class=\"c1\"># non-empty line with negative indent should stop the block:</span>\n
    \                   <span class=\"c1\"># - !!!</span>\n                    <span
    class=\"c1\">#  test</span>\n                    <span class=\"k\">break</span>\n\n
    \           <span class=\"c1\"># this will prevent lazy continuations from ever
    going past our end marker</span>\n            <span class=\"n\">state</span><span
    class=\"o\">.</span><span class=\"n\">lineMax</span> <span class=\"o\">=</span>
    <span class=\"n\">next_line</span>\n\n            <span class=\"n\">tag</span><span
    class=\"p\">,</span> <span class=\"n\">title</span><span class=\"p\">,</span>
    <span class=\"nb\">open</span> <span class=\"o\">=</span> <span class=\"n\">get_tag</span><span
    class=\"p\">(</span><span class=\"n\">params</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">attrs</span> <span class=\"o\">=</span> <span class=\"p\">{</span><span
    class=\"s2\">&quot;class&quot;</span><span class=\"p\">:</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;details </span><span class=\"si\">{</span><span class=\"n\">tag</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">}</span>\n
    \           <span class=\"k\">if</span> <span class=\"nb\">open</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">attrs</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;open&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;&quot;</span>\n\n            <span class=\"n\">token</span>
    <span class=\"o\">=</span> <span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">push</span><span class=\"p\">(</span><span class=\"s2\">&quot;details_open&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;details&quot;</span><span class=\"p\">,</span>
    <span class=\"mi\">1</span><span class=\"p\">)</span>\n            <span class=\"n\">token</span><span
    class=\"o\">.</span><span class=\"n\">markup</span> <span class=\"o\">=</span>
    <span class=\"n\">markup</span>\n            <span class=\"n\">token</span><span
    class=\"o\">.</span><span class=\"n\">block</span> <span class=\"o\">=</span>
    <span class=\"kc\">True</span>\n            <span class=\"n\">token</span><span
    class=\"o\">.</span><span class=\"n\">attrs</span> <span class=\"o\">=</span>
    <span class=\"n\">attrs</span>\n            <span class=\"n\">token</span><span
    class=\"o\">.</span><span class=\"n\">meta</span> <span class=\"o\">=</span> <span
    class=\"p\">{</span><span class=\"s2\">&quot;tag&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">tag</span><span class=\"p\">}</span>\n            <span class=\"n\">token</span><span
    class=\"o\">.</span><span class=\"n\">content</span> <span class=\"o\">=</span>
    <span class=\"n\">title</span>\n            <span class=\"n\">token</span><span
    class=\"o\">.</span><span class=\"n\">info</span> <span class=\"o\">=</span> <span
    class=\"n\">params</span>\n            <span class=\"n\">token</span><span class=\"o\">.</span><span
    class=\"n\">map</span> <span class=\"o\">=</span> <span class=\"p\">[</span><span
    class=\"n\">startLine</span><span class=\"p\">,</span> <span class=\"n\">next_line</span><span
    class=\"p\">]</span>\n\n            <span class=\"k\">if</span> <span class=\"n\">title</span><span
    class=\"p\">:</span>\n                <span class=\"n\">title_markup</span> <span
    class=\"o\">=</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">markup</span><span class=\"si\">}</span><span
    class=\"s2\"> </span><span class=\"si\">{</span><span class=\"n\">tag</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span>\n                <span class=\"n\">token</span>
    <span class=\"o\">=</span> <span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">push</span><span class=\"p\">(</span><span class=\"s2\">&quot;details_title_open&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;summary&quot;</span><span class=\"p\">,</span>
    <span class=\"mi\">1</span><span class=\"p\">)</span>\n                <span class=\"n\">token</span><span
    class=\"o\">.</span><span class=\"n\">markup</span> <span class=\"o\">=</span>
    <span class=\"n\">title_markup</span>\n                <span class=\"n\">token</span><span
    class=\"o\">.</span><span class=\"n\">attrs</span> <span class=\"o\">=</span>
    <span class=\"p\">{</span><span class=\"s2\">&quot;class&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;admonition-title&quot;</span><span class=\"p\">}</span>\n
    \               <span class=\"n\">token</span><span class=\"o\">.</span><span
    class=\"n\">map</span> <span class=\"o\">=</span> <span class=\"p\">[</span><span
    class=\"n\">startLine</span><span class=\"p\">,</span> <span class=\"n\">startLine</span>
    <span class=\"o\">+</span> <span class=\"mi\">1</span><span class=\"p\">]</span>\n\n
    \               <span class=\"n\">token</span> <span class=\"o\">=</span> <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">push</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;inline&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;&quot;</span><span class=\"p\">,</span> <span class=\"mi\">0</span><span
    class=\"p\">)</span>\n                <span class=\"n\">token</span><span class=\"o\">.</span><span
    class=\"n\">content</span> <span class=\"o\">=</span> <span class=\"n\">title</span>\n
    \               <span class=\"n\">token</span><span class=\"o\">.</span><span
    class=\"n\">map</span> <span class=\"o\">=</span> <span class=\"p\">[</span><span
    class=\"n\">startLine</span><span class=\"p\">,</span> <span class=\"n\">startLine</span>
    <span class=\"o\">+</span> <span class=\"mi\">1</span><span class=\"p\">]</span>\n
    \               <span class=\"n\">token</span><span class=\"o\">.</span><span
    class=\"n\">children</span> <span class=\"o\">=</span> <span class=\"p\">[]</span>\n\n
    \               <span class=\"n\">token</span> <span class=\"o\">=</span> <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">push</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;details_title_close&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;summary&quot;</span><span class=\"p\">,</span>
    <span class=\"o\">-</span><span class=\"mi\">1</span><span class=\"p\">)</span>\n
    \               <span class=\"n\">token</span><span class=\"o\">.</span><span
    class=\"n\">markup</span> <span class=\"o\">=</span> <span class=\"n\">title_markup</span>\n\n
    \           <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">md</span><span
    class=\"o\">.</span><span class=\"n\">block</span><span class=\"o\">.</span><span
    class=\"n\">tokenize</span><span class=\"p\">(</span><span class=\"n\">state</span><span
    class=\"p\">,</span> <span class=\"n\">startLine</span> <span class=\"o\">+</span>
    <span class=\"mi\">1</span><span class=\"p\">,</span> <span class=\"n\">next_line</span><span
    class=\"p\">)</span>\n\n            <span class=\"n\">token</span> <span class=\"o\">=</span>
    <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">push</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;details_close&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;details&quot;</span><span class=\"p\">,</span> <span
    class=\"o\">-</span><span class=\"mi\">1</span><span class=\"p\">)</span>\n            <span
    class=\"n\">token</span><span class=\"o\">.</span><span class=\"n\">markup</span>
    <span class=\"o\">=</span> <span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">src</span><span class=\"p\">[</span><span class=\"n\">start</span><span
    class=\"p\">:</span><span class=\"n\">pos</span><span class=\"p\">]</span>\n            <span
    class=\"n\">token</span><span class=\"o\">.</span><span class=\"n\">block</span>
    <span class=\"o\">=</span> <span class=\"kc\">True</span>\n\n            <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">parentType</span>
    <span class=\"o\">=</span> <span class=\"n\">old_parent</span>\n            <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">lineMax</span>
    <span class=\"o\">=</span> <span class=\"n\">old_line_max</span>\n            <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">blkIndent</span>
    <span class=\"o\">=</span> <span class=\"n\">old_indent</span>\n            <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">line</span>
    <span class=\"o\">=</span> <span class=\"n\">next_line</span>\n\n            <span
    class=\"k\">return</span> <span class=\"kc\">True</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='details_plugin' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>details_plugin <em class='small'>function</em></h2>\nPlugin to use\n<code>python-markdown
    style detailss     &lt;https://python-markdown.github.io/extensions/details&gt;</code>_.</p>\n<pre><code>..
    code-block:: md\n\n    ??? note\n        *content*\n\nNote, this is ported from\n`markdown-it-admon\n&lt;https://github.com/commenthol/markdown-it-admon&gt;`_.\n</code></pre>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">details_plugin
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
    <span class=\"nf\">details_plugin</span><span class=\"p\">(</span><span class=\"n\">md</span><span
    class=\"p\">:</span> <span class=\"n\">MarkdownIt</span><span class=\"p\">,</span>
    <span class=\"n\">render</span><span class=\"p\">:</span> <span class=\"n\">Optional</span><span
    class=\"p\">[</span><span class=\"n\">Callable</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"kc\">None</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;Plugin to
    use</span>\n<span class=\"sd\">            `python-markdown style detailss</span>\n<span
    class=\"sd\">            &lt;https://python-markdown.github.io/extensions/details&gt;`_.</span>\n\n<span
    class=\"sd\">            .. code-block:: md</span>\n\n<span class=\"sd\">                ???
    note</span>\n<span class=\"sd\">                    *content*</span>\n\n<span
    class=\"sd\">            Note, this is ported from</span>\n<span class=\"sd\">
    \           `markdown-it-admon</span>\n<span class=\"sd\">            &lt;https://github.com/commenthol/markdown-it-admon&gt;`_.</span>\n<span
    class=\"sd\">            &quot;&quot;&quot;</span>\n\n            <span class=\"k\">def</span>
    <span class=\"nf\">renderDefault</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">tokens</span><span class=\"p\">,</span>
    <span class=\"n\">idx</span><span class=\"p\">,</span> <span class=\"n\">_options</span><span
    class=\"p\">,</span> <span class=\"n\">env</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">renderToken</span><span class=\"p\">(</span><span
    class=\"n\">tokens</span><span class=\"p\">,</span> <span class=\"n\">idx</span><span
    class=\"p\">,</span> <span class=\"n\">_options</span><span class=\"p\">,</span>
    <span class=\"n\">env</span><span class=\"p\">)</span>\n\n            <span class=\"n\">render</span>
    <span class=\"o\">=</span> <span class=\"n\">render</span> <span class=\"ow\">or</span>
    <span class=\"n\">renderDefault</span>\n\n            <span class=\"n\">md</span><span
    class=\"o\">.</span><span class=\"n\">add_render_rule</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;details_open&quot;</span><span class=\"p\">,</span> <span class=\"n\">render</span><span
    class=\"p\">)</span>\n            <span class=\"n\">md</span><span class=\"o\">.</span><span
    class=\"n\">add_render_rule</span><span class=\"p\">(</span><span class=\"s2\">&quot;details_close&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">render</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">md</span><span class=\"o\">.</span><span class=\"n\">add_render_rule</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;details_title_open&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">render</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">md</span><span class=\"o\">.</span><span class=\"n\">add_render_rule</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;details_title_close&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">render</span><span class=\"p\">)</span>\n\n
    \           <span class=\"n\">md</span><span class=\"o\">.</span><span class=\"n\">block</span><span
    class=\"o\">.</span><span class=\"n\">ruler</span><span class=\"o\">.</span><span
    class=\"n\">before</span><span class=\"p\">(</span>\n                <span class=\"s2\">&quot;fence&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;details&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">details</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">{</span><span class=\"s2\">&quot;alt&quot;</span><span
    class=\"p\">:</span> <span class=\"p\">[</span><span class=\"s2\">&quot;paragraph&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;reference&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;blockquote&quot;</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;list&quot;</span><span class=\"p\">]},</span>\n            <span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='renderDefault'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>renderDefault <em
    class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">renderDefault <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">renderDefault</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">tokens</span><span class=\"p\">,</span>
    <span class=\"n\">idx</span><span class=\"p\">,</span> <span class=\"n\">_options</span><span
    class=\"p\">,</span> <span class=\"n\">env</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">renderToken</span><span class=\"p\">(</span><span
    class=\"n\">tokens</span><span class=\"p\">,</span> <span class=\"n\">idx</span><span
    class=\"p\">,</span> <span class=\"n\">_options</span><span class=\"p\">,</span>
    <span class=\"n\">env</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>"
  raw.md: ''
published: true
slug: markata/plugins/mdit-details
title: Mdit_Details.Py


---

None


!! function <h2 id='get_tag' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get_tag <em class='small'>function</em></h2>

???+ source "get_tag <em class='small'>source</em>"

```python

        def get_tag(params: str) -> Tuple[str, str, bool]:
            open = False
            if not params.strip():
                return "", "", open

            if params.strip().startswith("+"):
                params = params.strip("+")
                open = True

            tag, *_title = params.strip().split(" ")
            joined = " ".join(_title).strip('"').strip("'")

            title = ""
            if not joined:
                title = tag.title()
            elif joined != '""':
                title = joined
            return (tag.lower(), title, open)
```


!! function <h2 id='validate' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>validate <em class='small'>function</em></h2>

???+ source "validate <em class='small'>source</em>"

```python

        def validate(params: str) -> bool:
            tag = params.strip().split(" ", 1)[-1] or ""
            return bool(tag)
```


!! function <h2 id='details' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>details <em class='small'>function</em></h2>

???+ source "details <em class='small'>source</em>"

```python

        def details(state: StateBlock, startLine: int, endLine: int, silent: bool) -> bool:
            start = state.bMarks[startLine] + state.tShift[startLine]
            maximum = state.eMarks[startLine]

            # Check out the first character quickly, which should filter out most of non-containers
            if ord(state.src[start]) != MARKER_CHAR:
                return False

            # Check out the rest of the marker string
            pos = start + 1
            while pos <= maximum and MARKER_STR[(pos - start) % MARKER_LEN] == state.src[pos]:
                pos += 1

            marker_count = math.floor((pos - start) / MARKER_LEN)
            if marker_count < MIN_MARKERS:
                return False
            marker_pos = pos - ((pos - start) % MARKER_LEN)
            params = state.src[marker_pos:maximum]
            markup = state.src[start:marker_pos]

            if not validate(params):
                return False

            # Since start is found, we can report success here in validation mode
            if silent:
                return True

            old_parent = state.parentType
            old_line_max = state.lineMax
            old_indent = state.blkIndent

            blk_start = pos
            while blk_start < maximum and state.src[blk_start] == " ":
                blk_start += 1

            state.parentType = "details"
            state.blkIndent += blk_start - start

            was_empty = False

            # Search for the end of the block
            next_line = startLine
            while True:
                next_line += 1
                if next_line >= endLine:
                    # unclosed block should be autoclosed by end of document.
                    # also block seems to be autoclosed by end of parent
                    break
                pos = state.bMarks[next_line] + state.tShift[next_line]
                maximum = state.eMarks[next_line]
                is_empty = state.sCount[next_line] < state.blkIndent

                # two consecutive empty lines autoclose the block
                if is_empty and was_empty:
                    break
                was_empty = is_empty

                if pos < maximum and state.sCount[next_line] < state.blkIndent:
                    # non-empty line with negative indent should stop the block:
                    # - !!!
                    #  test
                    break

            # this will prevent lazy continuations from ever going past our end marker
            state.lineMax = next_line

            tag, title, open = get_tag(params)
            attrs = {"class": f"details {tag}"}
            if open:
                attrs["open"] = ""

            token = state.push("details_open", "details", 1)
            token.markup = markup
            token.block = True
            token.attrs = attrs
            token.meta = {"tag": tag}
            token.content = title
            token.info = params
            token.map = [startLine, next_line]

            if title:
                title_markup = f"{markup} {tag}"
                token = state.push("details_title_open", "summary", 1)
                token.markup = title_markup
                token.attrs = {"class": "admonition-title"}
                token.map = [startLine, startLine + 1]

                token = state.push("inline", "", 0)
                token.content = title
                token.map = [startLine, startLine + 1]
                token.children = []

                token = state.push("details_title_close", "summary", -1)
                token.markup = title_markup

            state.md.block.tokenize(state, startLine + 1, next_line)

            token = state.push("details_close", "details", -1)
            token.markup = state.src[start:pos]
            token.block = True

            state.parentType = old_parent
            state.lineMax = old_line_max
            state.blkIndent = old_indent
            state.line = next_line

            return True
```


!! function <h2 id='details_plugin' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>details_plugin <em class='small'>function</em></h2>
    Plugin to use
    `python-markdown style detailss
    <https://python-markdown.github.io/extensions/details>`_.

    .. code-block:: md

        ??? note
            *content*

    Note, this is ported from
    `markdown-it-admon
    <https://github.com/commenthol/markdown-it-admon>`_.
???+ source "details_plugin <em class='small'>source</em>"

```python

        def details_plugin(md: MarkdownIt, render: Optional[Callable] = None) -> None:
            """Plugin to use
            `python-markdown style detailss
            <https://python-markdown.github.io/extensions/details>`_.

            .. code-block:: md

                ??? note
                    *content*

            Note, this is ported from
            `markdown-it-admon
            <https://github.com/commenthol/markdown-it-admon>`_.
            """

            def renderDefault(self, tokens, idx, _options, env):
                return self.renderToken(tokens, idx, _options, env)

            render = render or renderDefault

            md.add_render_rule("details_open", render)
            md.add_render_rule("details_close", render)
            md.add_render_rule("details_title_open", render)
            md.add_render_rule("details_title_close", render)

            md.block.ruler.before(
                "fence",
                "details",
                details,
                {"alt": ["paragraph", "reference", "blockquote", "list"]},
            )
```


!! function <h2 id='renderDefault' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>renderDefault <em class='small'>function</em></h2>

???+ source "renderDefault <em class='small'>source</em>"

```python

        def renderDefault(self, tokens, idx, _options, env):
                return self.renderToken(tokens, idx, _options, env)
```

