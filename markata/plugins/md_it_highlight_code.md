---
content: "---\n\nThe `markata.plugins.md_it_highlight_code` plugin adds syntax highlighting
  to code blocks\nin your markdown content. It uses Pygments for highlighting and
  adds a copy button for\neasy code sharing.\n\n# Installation\n\nThis plugin is built-in
  and enabled by default through the 'default' plugin.\nIf you want to be explicit,
  you can add it to your list of plugins:\n\n```toml\nhooks = [\n    \"markata.plugins.md_it_highlight_code\",\n]\n```\n\n#
  Uninstallation\n\nSince this plugin is included in the default plugin set, to disable
  it you must explicitly\nadd it to the disabled_hooks list if you are using the 'default'
  plugin:\n\n```toml\ndisabled_hooks = [\n    \"markata.plugins.md_it_highlight_code\",\n]\n```\n\n#
  Configuration\n\nConfigure syntax highlighting in your `markata.toml`:\n\n```toml\n[markata.highlight_code]\n#
  Optional: Custom CSS class for code blocks\ncode_class = \"highlight\"\n\n# Optional:
  Custom CSS class for copy button\ncopy_button_class = \"copy-button\"\n\n# Optional:
  Custom copy button text\ncopy_button_text = \"Copy\"\ncopy_button_copied_text =
  \"Copied!\"\n```\n\n# Functionality\n\n## Code Highlighting\n\nThe plugin:\n1. Detects
  language from code block info\n2. Applies Pygments syntax highlighting\n3. Adds
  line numbers (optional)\n4. Wraps code in proper HTML structure\n5. Adds a copy-to-clipboard
  button\n\n## Example Usage\n\nIn markdown:\n````markdown\n```python\ndef hello_world():\n
  \   print(\"Hello, World!\")\n```\n\n```javascript\nconsole.log('Hello, World!');\n```\n````\n\nGenerates
  HTML:\n```html\n<div class=\"highlight\">\n  <button class=\"copy-button\">\n    <svg><!--
  Copy icon SVG --></svg>\n  </button>\n  <pre><code class=\"language-python\">\n
  \   <span class=\"def\">def</span> <span class=\"name\">hello_world</span>():\n
  \       <span class=\"keyword\">print</span>(<span class=\"string\">\"Hello, World!\"</span>)\n
  \ </code></pre>\n</div>\n```\n\n## Supported Languages\n\nThe plugin supports all
  languages that Pygments recognizes, including:\n- Python\n- JavaScript\n- HTML/CSS\n-
  Bash/Shell\n- And many more\n\n## Copy Button\n\nFeatures:\n- Hover-to-show button\n-
  Click-to-copy functionality\n- Success feedback\n- Accessible keyboard support\n\n##
  Dependencies\n\nThis plugin depends on:\n- Pygments for syntax highlighting\n- markdown-it-py
  for markdown parsing\n\n---\n\n\n\n!!! function\n    <h2 id=\"highlight_code\" class=\"admonition-title\"
  style=\"margin: 0; padding: .5rem 1rem;\">highlight_code <em class=\"small\">function</em></h2>\n\n
  \   Code highlighter for markdown-it-py.\n\n???+ source \"highlight_code <em class='small'>source</em>\"\n
  \   ```python\n    def highlight_code(code, name, attrs, markata=None):\n        \"\"\"Code
  highlighter for markdown-it-py.\"\"\"\n        # from pygments import highlight\n
  \       # from pygments.formatters import HtmlFormatter\n        # from pygments.lexers
  import ClassNotFound, get_lexer_by_name\n\n        try:\n            lexer = get_lexer_by_name(name
  or \"text\")\n        except ClassNotFound:\n            lexer = get_lexer_by_name(\"text\")\n\n
  \       import re\n\n        pattern = r'(\\w+)\\s*=\\s*(\".*?\"|\\S+)'\n        matches
  = re.findall(pattern, attrs)\n        attrs = dict(matches)\n\n        if attrs.get(\"hl_lines\"):\n
  \           formatter = HtmlFormatter(hl_lines=attrs.get(\"hl_lines\"))\n        else:\n
  \           formatter = HtmlFormatter()\n\n        copy_button = f\"\"\"<button
  class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\">{COPY_ICON}</button>\"\"\"\n\n
  \       from markdown_it import MarkdownIt\n\n        md = MarkdownIt(\n            \"commonmark\",\n
  \           {\n                \"html\": True,\n                \"typographer\":
  True,\n            },\n        )\n\n        if attrs.get(\"help\"):\n            help
  = f\"\"\"\n            <a href={attrs.get(\"help\").strip(\"<\").strip(\">\").strip('\"').strip(\"'\")}
  title='help link' class='help'>{HELP_ICON}</a>\n            \"\"\"\n        else:\n
  \           help = \"\"\n        if attrs.get(\"title\"):\n            file = f\"\"\"\n
  \   <div class='filepath'>\n    {md.render(attrs.get(\"title\").strip('\"').strip(\"'\"))}\n
  \   <div class='right'>\n    {help}\n    {copy_button}\n    </div>\n    </div>\n
  \   \"\"\"\n        else:\n            file = f\"\"\"\n    <div class='copy-wrapper'>\n
  \   {help}\n    {copy_button}\n    </div>\n            \"\"\"\n        return f\"\"\"<pre
  class='wrapper'>\n    {file}\n    {highlight(code, lexer, formatter)}\n    </pre>\n
  \   \"\"\"\n    ```"
date: 2025-02-14
description: 'The `markata.plugins.md_it_highlight_code` plugin adds syntax highlighting
  to code blocks

  in your markdown content. It uses Pygments for highlighting and adds a'
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>md_it_highlight_code.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"The `markata.plugins.md_it_highlight_code`
    plugin adds syntax highlighting to code blocks\nin your markdown content. It uses
    Pygments for highlighting and adds a\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n
    \       <meta property=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    <head>\n        <script src=\"https://unpkg.com/@tailwindcss/browser@4\"></script>\n<title>md_it_highlight_code.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"The `markata.plugins.md_it_highlight_code`
    plugin adds syntax highlighting to code blocks\nin your markdown content. It uses
    Pygments for highlighting and adds a\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n
    \       <meta property=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n        <script>\n            document.addEventListener(\"DOMContentLoaded\",
    () => {\n                const collapsibleElements = document.querySelectorAll('.is-collapsible');\n
    \               collapsibleElements.forEach(el => {\n                    const
    summary = el.querySelector('.admonition-title');\n                    if (summary)
    {\n                        summary.style.cursor = 'pointer';\n                        summary.addEventListener('click',
    () => {\n                            el.classList.toggle('collapsible-open');\n
    \                       });\n                    }\n                });\n            });\n
    \       </script>\n\n        <style>\n\n            .admonition.source {\n                padding-bottom:
    0;\n            }\n            .admonition.source pre.wrapper {\n                margin:
    0;\n                padding: 0;\n            }\n            .is-collapsible {\n
    \               overflow: hidden;\n                transition: max-height 0.3s
    ease;\n            }\n            .is-collapsible:not(.collapsible-open) {\n                max-height:
    0;\n                padding-bottom: 2.5rem;\n            }\n            .admonition-title
    {\n                font-weight: bold;\n                margin-bottom: 8px;\n            }\n
    \       </style>\n    </head>\n    <body>\n<div class='container flex flex-row
    min-h-screen'>\n    <div>\n    </div>\n    <div class='flex-grow px-8 mx-auto
    min-h-screen'>\n<header class='flex justify-center items-center p-8'>\n\n    <nav
    class='flex justify-center items-center my-8'>\n        <a\n            href='/'>markata</a>\n
    \       <a\n            href='https://github.com/WaylonWalker/markata'>GitHub</a>\n
    \       <a\n            href='https://markata.dev/docs/'>docs</a>\n        <a\n
    \           href='https://markata.dev/plugins/'>plugins</a>\n    </nav>\n\n    <div>\n
    \       <label id=\"theme-switch\" class=\"theme-switch\" for=\"checkbox-theme\"
    title=\"light/dark mode toggle\">\n            <input type=\"checkbox\" id=\"checkbox-theme\"
    />\n            <div class=\"slider round\"></div>\n        </label>\n    </div>\n</header><div
    id='didyoumean'>\n    <div class=\"mb-0\">\n        <!-- <label for=\"search\"
    class=\"block text-sm font-medium mb-2\">Search for a page</label> -->\n        <input
    type=\"text\" id=\"search\"\n               class=\"w-full p-2 border rounded-md
    bg-gray-50 dark:bg-gray-800 focus:ring-2 focus:ring-pink-500\"\n               placeholder=\"'/'
    Search for a page\">\n    </div>\n\n    <!-- <div id=\"didyoumean_results\" class=\"grid
    gap-4 grid-cols-1 md:grid-cols-2 lg:grid-cols-3\"> -->\n    <ul id=\"didyoumean_results\"
    class='grid gap-4'>\n        <!-- Results will be populated here -->\n    </ul>\n</div>\n<script
    type='module'>\n// All available pages from Markata\n    // const pages =  markata.map(\"{'slug':slug,'title':title,'description':description,'tags':tags}\",
    filter=config.didyoumean_filter, sort='True')|tojson;\n    // fetch pages from
    config.output_dir / didyoumean.json\n\n    const pages = await fetch('/didyoumean.json').then(response
    => response.json());\n    const populate_search_input = false\n    const search_hotkey
    = \"/\"\n\n// Get current path from URL, removing leading/trailing slashes\n    if
    (populate_search_input) {\n        const currentPath = window.location.pathname.replace(/^\\/|\\/$/g,
    '');\n        document.getElementById('search').value = currentPath;\n    }\n\n//
    Search across all fields in an object\n    function searchObject(needle, obj)
    {\n        needle = needle.toLowerCase();\n        let score = 0;\n\n    // Helper
    to search a single field\n        const searchField = (value) => {\n            if
    (!value) return 0;\n            value = String(value).toLowerCase();\n\n            //
    Exact matches\n            if (value === needle) return 15;\n\n            //
    Word boundary matches (complete words)\n            if (value.match(new RegExp(`\\\\b${needle}\\\\b`)))
    return 10;\n\n            // Contains full search term\n            if (value.includes(needle))
    return 8;\n\n            // Most parts match (for multi-word searches)\n            const
    needleParts = needle.split(/\\W+/).filter(p => p.length > 2);\n            const
    valueParts = value.split(/\\W+/).filter(p => p.length > 2);\n\n            if
    (needleParts.length === 0) return 0;\n\n            let matchCount = 0;\n            for
    (const part of needleParts) {\n                for (const valuePart of valueParts)
    {\n                    if (valuePart.includes(part) || part.includes(valuePart))
    {\n                        matchCount++;\n                        break;\n                    }\n
    \               }\n            }\n\n            // Only count if most parts match\n
    \           const matchRatio = matchCount / needleParts.length;\n            if
    (matchRatio >= 0.75) {\n                return matchRatio * 6;\n            }\n\n
    \           return 0;\n        };\n\n    // Search each field with different weights\n
    \       const slugScore = searchField(obj.slug) * 3;  // Slug is most important\n
    \       const titleScore = searchField(obj.title) * 2;  // Title is next\n        const
    descScore = searchField(obj.description) * 1;  // Description\n        const tagScore
    = (obj.tags || []).reduce((sum, tag) => sum + searchField(tag), 0);  // Tags\n\n
    \       score = slugScore + titleScore + descScore + tagScore;\n\n    // Path
    segment matches for slug (only if we have some other match)\n        if (score
    > 0 && obj.slug) {\n            const inputParts = needle.split('/').filter(p
    => p.length > 0);\n            const slugParts = obj.slug.toLowerCase().split('/');\n\n
    \           // Bonus for matching path structure\n            for (let i = 0;
    i < inputParts.length && i < slugParts.length; i++) {\n                if (slugParts[i].includes(inputParts[i]))
    {\n                    score += 5;  // Matching segments in order is valuable\n
    \               }\n            }\n        }\n\n        return score;\n    }\n\n//
    Find similar pages\n    function findSimilar(input) {\n        if (!input || input.length
    < 2) return [];\n        const normalizedInput = input.toLowerCase().trim();\n\n
    \   // Score each page\n        const scored = pages.map(page => ({\n            ...page,\n
    \           score: searchObject(normalizedInput, page)\n        }));\n\n    //
    Sort by score (higher is better) and take top matches\n        return scored\n
    \           .sort((a, b) => b.score - a.score)\n            .slice(0, 12)  //
    Show more results in the grid\n            .filter(item => item.score > 15); //
    Only show strong matches\n    }\n\n// Update results in the DOM\n    function
    updateResults(results) {\n        const resultsDiv = document.getElementById('didyoumean_results');\n\n
    \       if (results.length === 0) {\n            resultsDiv.innerHTML = '<p class=\"text-gray-500
    col-span-full text-center py-8\">No similar pages found.</p>';\n            return;\n
    \       }\n\n        const html = results.map(page => `\n        <li class=\"p-4
    bg-gray-50 dark:bg-gray-800 rounded-lg hover:shadow-lg transition-shadow first:mt-4\">\n
    \           <a href=\"/${page.slug}\" class=\"block\">\n                <h3 class=\"text-lg
    font-semibold text-pink-500 hover:text-pink-600 dark:text-pink-400 dark:hover:text-pink-300
    mb-2\">\n                    ${page.title || page.slug}\n                </h3>\n
    \               ${page.description ? `\n            <p class=\"text-sm text-gray-600
    dark:text-gray-300 mb-3 line-clamp-2\">\n            ${page.description}\n            </p>\n
    \           ` : ''}\n                <div class=\"flex flex-wrap gap-2 text-xs
    text-gray-500 dark:text-gray-400\">\n                </div>\n                ${page.tags
    && page.tags.length > 0 ? `\n            <div class=\"mt-3 flex flex-wrap gap-2\">\n
    \           ${page.tags.map(tag => `\n                            <span class=\"px-2
    py-1 bg-gray-100 dark:bg-gray-700 rounded text-xs\">\n                                ${tag}\n
    \                           </span>\n                        `).join('')}\n            </div>\n
    \           ` : ''}\n            </a>\n        </li>\n    `).join('');\n\n        resultsDiv.innerHTML
    = html;\n    }\n\n// Set up hotkey for search if configured\n    if (search_hotkey)
    {\n        document.addEventListener('keydown', (e) => {\n            // Don't
    trigger if user is typing in an input or textarea\n            if (e.target.tagName
    === 'INPUT' || e.target.tagName === 'TEXTAREA') {\n                return;\n            }\n\n
    \           // Check if the pressed key matches the hotkey\n            if (e.key
    === search_hotkey) {\n                e.preventDefault();  // Prevent the '/'
    from being typed\n                const searchInput = document.getElementById('search');\n
    \               searchInput.focus();\n                searchInput.select();  //
    Select any existing text\n            }\n        });\n    }\n\n// Set up search
    input handler with debounce\n    let debounceTimeout;\n    const searchInput =
    document.getElementById('search');\n    searchInput.addEventListener('input',
    (e) => {\n        clearTimeout(debounceTimeout);\n        debounceTimeout = setTimeout(()
    => {\n            const results = findSimilar(e.target.value);\n            updateResults(results);\n
    \       }, 100);\n    });\n\n// Initial search with current path\n    if (populate_search_input)
    {\n        updateResults(findSimilar(currentPath));\n    }\n</script><article
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        md_it_highlight_code.py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <hr />\n<p>The <code>markata.plugins.md_it_highlight_code</code>
    plugin adds syntax highlighting to code blocks\nin your markdown content. It uses
    Pygments for highlighting and adds a copy button for\neasy code sharing.</p>\n<h1
    id=\"installation\">Installation <a class=\"header-anchor\" href=\"#installation\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>This plugin is built-in
    and enabled by default through the 'default' plugin.\nIf you want to be explicit,
    you can add it to your list of plugins:</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"n\">hooks</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"p\">[</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;markata.plugins.md_it_highlight_code&quot;</span><span
    class=\"p\">,</span>\n<span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<h1
    id=\"uninstallation\">Uninstallation <a class=\"header-anchor\" href=\"#uninstallation\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>Since this plugin is
    included in the default plugin set, to disable it you must explicitly\nadd it
    to the disabled_hooks list if you are using the 'default' plugin:</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"n\">disabled_hooks</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"p\">[</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;markata.plugins.md_it_highlight_code&quot;</span><span
    class=\"p\">,</span>\n<span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<h1
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>Configure syntax highlighting
    in your <code>markata.toml</code>:</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata.highlight_code]</span>\n<span
    class=\"c1\"># Optional: Custom CSS class for code blocks</span>\n<span class=\"n\">code_class</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;highlight&quot;</span>\n\n<span class=\"c1\"># Optional: Custom
    CSS class for copy button</span>\n<span class=\"n\">copy_button_class</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;copy-button&quot;</span>\n\n<span class=\"c1\"># Optional:
    Custom copy button text</span>\n<span class=\"n\">copy_button_text</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;Copy&quot;</span>\n<span class=\"n\">copy_button_copied_text</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;Copied!&quot;</span>\n</pre></div>\n\n</pre>\n\n<h1 id=\"functionality\">Functionality
    <a class=\"header-anchor\" href=\"#functionality\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<h2 id=\"code-highlighting\">Code
    Highlighting <a class=\"header-anchor\" href=\"#code-highlighting\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>The plugin:</p>\n<ol>\n<li>Detects
    language from code block info</li>\n<li>Applies Pygments syntax highlighting</li>\n<li>Adds
    line numbers (optional)</li>\n<li>Wraps code in proper HTML structure</li>\n<li>Adds
    a copy-to-clipboard button</li>\n</ol>\n<h2 id=\"example-usage\">Example Usage
    <a class=\"header-anchor\" href=\"#example-usage\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>In markdown:</p>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"sb\">```python</span>\n<span
    class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">hello_world</span><span
    class=\"p\">():</span>\n    <span class=\"nb\">print</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;Hello, World!&quot;</span><span class=\"p\">)</span>\n<span
    class=\"sb\">```</span>\n\n<span class=\"sb\">```javascript</span>\n<span class=\"nx\">console</span><span
    class=\"p\">.</span><span class=\"nx\">log</span><span class=\"p\">(</span><span
    class=\"s1\">&#39;Hello, World!&#39;</span><span class=\"p\">);</span>\n<span
    class=\"sb\">```</span>\n</pre></div>\n\n</pre>\n\n<p>Generates HTML:</p>\n<pre
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
    class=\"nt\">div</span> <span class=\"na\">class</span><span class=\"o\">=</span><span
    class=\"s\">&quot;highlight&quot;</span><span class=\"p\">&gt;</span>\n  <span
    class=\"p\">&lt;</span><span class=\"nt\">button</span> <span class=\"na\">class</span><span
    class=\"o\">=</span><span class=\"s\">&quot;copy-button&quot;</span><span class=\"p\">&gt;</span>\n
    \   <span class=\"p\">&lt;</span><span class=\"nt\">svg</span><span class=\"p\">&gt;</span><span
    class=\"cm\">&lt;!-- Copy icon SVG --&gt;</span><span class=\"p\">&lt;/</span><span
    class=\"nt\">svg</span><span class=\"p\">&gt;</span>\n  <span class=\"p\">&lt;/</span><span
    class=\"nt\">button</span><span class=\"p\">&gt;</span>\n  <span class=\"p\">&lt;</span><span
    class=\"nt\">pre</span><span class=\"p\">&gt;&lt;</span><span class=\"nt\">code</span>
    <span class=\"na\">class</span><span class=\"o\">=</span><span class=\"s\">&quot;language-python&quot;</span><span
    class=\"p\">&gt;</span>\n    <span class=\"p\">&lt;</span><span class=\"nt\">span</span>
    <span class=\"na\">class</span><span class=\"o\">=</span><span class=\"s\">&quot;def&quot;</span><span
    class=\"p\">&gt;</span>def<span class=\"p\">&lt;/</span><span class=\"nt\">span</span><span
    class=\"p\">&gt;</span> <span class=\"p\">&lt;</span><span class=\"nt\">span</span>
    <span class=\"na\">class</span><span class=\"o\">=</span><span class=\"s\">&quot;name&quot;</span><span
    class=\"p\">&gt;</span>hello_world<span class=\"p\">&lt;/</span><span class=\"nt\">span</span><span
    class=\"p\">&gt;</span>():\n        <span class=\"p\">&lt;</span><span class=\"nt\">span</span>
    <span class=\"na\">class</span><span class=\"o\">=</span><span class=\"s\">&quot;keyword&quot;</span><span
    class=\"p\">&gt;</span>print<span class=\"p\">&lt;/</span><span class=\"nt\">span</span><span
    class=\"p\">&gt;</span>(<span class=\"p\">&lt;</span><span class=\"nt\">span</span>
    <span class=\"na\">class</span><span class=\"o\">=</span><span class=\"s\">&quot;string&quot;</span><span
    class=\"p\">&gt;</span>&quot;Hello, World!&quot;<span class=\"p\">&lt;/</span><span
    class=\"nt\">span</span><span class=\"p\">&gt;</span>)\n  <span class=\"p\">&lt;/</span><span
    class=\"nt\">code</span><span class=\"p\">&gt;&lt;/</span><span class=\"nt\">pre</span><span
    class=\"p\">&gt;</span>\n<span class=\"p\">&lt;/</span><span class=\"nt\">div</span><span
    class=\"p\">&gt;</span>\n</pre></div>\n\n</pre>\n\n<h2 id=\"supported-languages\">Supported
    Languages <a class=\"header-anchor\" href=\"#supported-languages\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>The plugin supports
    all languages that Pygments recognizes, including:</p>\n<ul>\n<li>Python</li>\n<li>JavaScript</li>\n<li>HTML/CSS</li>\n<li>Bash/Shell</li>\n<li>And
    many more</li>\n</ul>\n<h2 id=\"copy-button\">Copy Button <a class=\"header-anchor\"
    href=\"#copy-button\"><svg class=\"heading-permalink\" aria-hidden=\"true\" fill=\"currentColor\"
    focusable=\"false\" height=\"1em\" viewBox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Features:</p>\n<ul>\n<li>Hover-to-show
    button</li>\n<li>Click-to-copy functionality</li>\n<li>Success feedback</li>\n<li>Accessible
    keyboard support</li>\n</ul>\n<h2 id=\"dependencies\">Dependencies <a class=\"header-anchor\"
    href=\"#dependencies\"><svg class=\"heading-permalink\" aria-hidden=\"true\" fill=\"currentColor\"
    focusable=\"false\" height=\"1em\" viewBox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>This plugin depends
    on:</p>\n<ul>\n<li>Pygments for syntax highlighting</li>\n<li>markdown-it-py for
    markdown parsing</li>\n</ul>\n<hr />\n<div class=\"admonition function\">\n<p
    class=\"admonition-title\">Function</p>\n<h2 id=\"highlight_code\" class=\"admonition-title\"
    style=\"margin: 0; padding: .5rem 1rem;\">highlight_code <em class=\"small\">function</em></h2>\n<p>Code
    highlighter for markdown-it-py.</p>\n</div>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">highlight_code <em class='small'>source</em></p>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">def</span><span
    class=\"w\"> </span><span class=\"nf\">highlight_code</span><span class=\"p\">(</span><span
    class=\"n\">code</span><span class=\"p\">,</span> <span class=\"n\">name</span><span
    class=\"p\">,</span> <span class=\"n\">attrs</span><span class=\"p\">,</span>
    <span class=\"n\">markata</span><span class=\"o\">=</span><span class=\"kc\">None</span><span
    class=\"p\">):</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Code
    highlighter for markdown-it-py.&quot;&quot;&quot;</span>\n    <span class=\"c1\">#
    from pygments import highlight</span>\n    <span class=\"c1\"># from pygments.formatters
    import HtmlFormatter</span>\n    <span class=\"c1\"># from pygments.lexers import
    ClassNotFound, get_lexer_by_name</span>\n\n    <span class=\"k\">try</span><span
    class=\"p\">:</span>\n        <span class=\"n\">lexer</span> <span class=\"o\">=</span>
    <span class=\"n\">get_lexer_by_name</span><span class=\"p\">(</span><span class=\"n\">name</span>
    <span class=\"ow\">or</span> <span class=\"s2\">&quot;text&quot;</span><span class=\"p\">)</span>\n
    \   <span class=\"k\">except</span> <span class=\"n\">ClassNotFound</span><span
    class=\"p\">:</span>\n        <span class=\"n\">lexer</span> <span class=\"o\">=</span>
    <span class=\"n\">get_lexer_by_name</span><span class=\"p\">(</span><span class=\"s2\">&quot;text&quot;</span><span
    class=\"p\">)</span>\n\n    <span class=\"kn\">import</span><span class=\"w\">
    </span><span class=\"nn\">re</span>\n\n    <span class=\"n\">pattern</span> <span
    class=\"o\">=</span> <span class=\"sa\">r</span><span class=\"s1\">&#39;(\\w+)\\s*=\\s*(&quot;.*?&quot;|\\S+)&#39;</span>\n
    \   <span class=\"n\">matches</span> <span class=\"o\">=</span> <span class=\"n\">re</span><span
    class=\"o\">.</span><span class=\"n\">findall</span><span class=\"p\">(</span><span
    class=\"n\">pattern</span><span class=\"p\">,</span> <span class=\"n\">attrs</span><span
    class=\"p\">)</span>\n    <span class=\"n\">attrs</span> <span class=\"o\">=</span>
    <span class=\"nb\">dict</span><span class=\"p\">(</span><span class=\"n\">matches</span><span
    class=\"p\">)</span>\n\n    <span class=\"k\">if</span> <span class=\"n\">attrs</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;hl_lines&quot;</span><span class=\"p\">):</span>\n        <span
    class=\"n\">formatter</span> <span class=\"o\">=</span> <span class=\"n\">HtmlFormatter</span><span
    class=\"p\">(</span><span class=\"n\">hl_lines</span><span class=\"o\">=</span><span
    class=\"n\">attrs</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;hl_lines&quot;</span><span class=\"p\">))</span>\n
    \   <span class=\"k\">else</span><span class=\"p\">:</span>\n        <span class=\"n\">formatter</span>
    <span class=\"o\">=</span> <span class=\"n\">HtmlFormatter</span><span class=\"p\">()</span>\n\n
    \   <span class=\"n\">copy_button</span> <span class=\"o\">=</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;&quot;&quot;&lt;button class=&#39;copy&#39; title=&#39;copy
    code to clipboard&#39; onclick=&quot;navigator.clipboard.writeText(this.parentElement.parentElement.querySelector(&#39;pre&#39;).textContent)&quot;&gt;</span><span
    class=\"si\">{</span><span class=\"n\">COPY_ICON</span><span class=\"si\">}</span><span
    class=\"s2\">&lt;/button&gt;&quot;&quot;&quot;</span>\n\n    <span class=\"kn\">from</span><span
    class=\"w\"> </span><span class=\"nn\">markdown_it</span><span class=\"w\"> </span><span
    class=\"kn\">import</span> <span class=\"n\">MarkdownIt</span>\n\n    <span class=\"n\">md</span>
    <span class=\"o\">=</span> <span class=\"n\">MarkdownIt</span><span class=\"p\">(</span>\n
    \       <span class=\"s2\">&quot;commonmark&quot;</span><span class=\"p\">,</span>\n
    \       <span class=\"p\">{</span>\n            <span class=\"s2\">&quot;html&quot;</span><span
    class=\"p\">:</span> <span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \           <span class=\"s2\">&quot;typographer&quot;</span><span class=\"p\">:</span>
    <span class=\"kc\">True</span><span class=\"p\">,</span>\n        <span class=\"p\">},</span>\n
    \   <span class=\"p\">)</span>\n\n    <span class=\"k\">if</span> <span class=\"n\">attrs</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;help&quot;</span><span class=\"p\">):</span>\n        <span
    class=\"n\">help</span> <span class=\"o\">=</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;&quot;&quot;</span>\n<span class=\"s2\">        &lt;a href=</span><span
    class=\"si\">{</span><span class=\"n\">attrs</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;help&quot;</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">strip</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;&lt;&quot;</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">strip</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;&gt;&quot;</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">strip</span><span class=\"p\">(</span><span class=\"s1\">&#39;&quot;&#39;</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">strip</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;&#39;&quot;</span><span class=\"p\">)</span><span
    class=\"si\">}</span><span class=\"s2\"> title=&#39;help link&#39; class=&#39;help&#39;&gt;</span><span
    class=\"si\">{</span><span class=\"n\">HELP_ICON</span><span class=\"si\">}</span><span
    class=\"s2\">&lt;/a&gt;</span>\n<span class=\"s2\">        &quot;&quot;&quot;</span>\n
    \   <span class=\"k\">else</span><span class=\"p\">:</span>\n        <span class=\"n\">help</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;&quot;</span>\n    <span class=\"k\">if</span>
    <span class=\"n\">attrs</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;title&quot;</span><span class=\"p\">):</span>\n
    \       <span class=\"n\">file</span> <span class=\"o\">=</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;&quot;&quot;</span>\n<span class=\"s2\">&lt;div class=&#39;filepath&#39;&gt;</span>\n<span
    class=\"si\">{</span><span class=\"n\">md</span><span class=\"o\">.</span><span
    class=\"n\">render</span><span class=\"p\">(</span><span class=\"n\">attrs</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;title&quot;</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">strip</span><span class=\"p\">(</span><span class=\"s1\">&#39;&quot;&#39;</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">strip</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;&#39;&quot;</span><span class=\"p\">))</span><span
    class=\"si\">}</span>\n<span class=\"s2\">&lt;div class=&#39;right&#39;&gt;</span>\n<span
    class=\"si\">{</span><span class=\"n\">help</span><span class=\"si\">}</span>\n<span
    class=\"si\">{</span><span class=\"n\">copy_button</span><span class=\"si\">}</span>\n<span
    class=\"s2\">&lt;/div&gt;</span>\n<span class=\"s2\">&lt;/div&gt;</span>\n<span
    class=\"s2\">&quot;&quot;&quot;</span>\n    <span class=\"k\">else</span><span
    class=\"p\">:</span>\n        <span class=\"n\">file</span> <span class=\"o\">=</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;&quot;&quot;</span>\n<span
    class=\"s2\">&lt;div class=&#39;copy-wrapper&#39;&gt;</span>\n<span class=\"si\">{</span><span
    class=\"n\">help</span><span class=\"si\">}</span>\n<span class=\"si\">{</span><span
    class=\"n\">copy_button</span><span class=\"si\">}</span>\n<span class=\"s2\">&lt;/div&gt;</span>\n<span
    class=\"s2\">        &quot;&quot;&quot;</span>\n    <span class=\"k\">return</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;&quot;&quot;&lt;pre class=&#39;wrapper&#39;&gt;</span>\n<span
    class=\"si\">{</span><span class=\"n\">file</span><span class=\"si\">}</span>\n<span
    class=\"si\">{</span><span class=\"n\">highlight</span><span class=\"p\">(</span><span
    class=\"n\">code</span><span class=\"p\">,</span><span class=\"w\"> </span><span
    class=\"n\">lexer</span><span class=\"p\">,</span><span class=\"w\"> </span><span
    class=\"n\">formatter</span><span class=\"p\">)</span><span class=\"si\">}</span>\n<span
    class=\"s2\">&lt;/pre&gt;</span>\n<span class=\"s2\">&quot;&quot;&quot;</span>\n</pre></div>\n\n</pre>\n\n</div>\n\n
    \   </section>\n</article>    </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>md_it_highlight_code.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"The `markata.plugins.md_it_highlight_code`
    plugin adds syntax highlighting to code blocks\nin your markdown content. It uses
    Pygments for highlighting and adds a\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n
    \       <meta property=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    <head>\n        <script src=\"https://unpkg.com/@tailwindcss/browser@4\"></script>\n<title>md_it_highlight_code.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"The `markata.plugins.md_it_highlight_code`
    plugin adds syntax highlighting to code blocks\nin your markdown content. It uses
    Pygments for highlighting and adds a\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n
    \       <meta property=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n        <script>\n            document.addEventListener(\"DOMContentLoaded\",
    () => {\n                const collapsibleElements = document.querySelectorAll('.is-collapsible');\n
    \               collapsibleElements.forEach(el => {\n                    const
    summary = el.querySelector('.admonition-title');\n                    if (summary)
    {\n                        summary.style.cursor = 'pointer';\n                        summary.addEventListener('click',
    () => {\n                            el.classList.toggle('collapsible-open');\n
    \                       });\n                    }\n                });\n            });\n
    \       </script>\n\n        <style>\n\n            .admonition.source {\n                padding-bottom:
    0;\n            }\n            .admonition.source pre.wrapper {\n                margin:
    0;\n                padding: 0;\n            }\n            .is-collapsible {\n
    \               overflow: hidden;\n                transition: max-height 0.3s
    ease;\n            }\n            .is-collapsible:not(.collapsible-open) {\n                max-height:
    0;\n                padding-bottom: 2.5rem;\n            }\n            .admonition-title
    {\n                font-weight: bold;\n                margin-bottom: 8px;\n            }\n
    \       </style>\n    </head>\n    <body>\n<article style=\"text-align: center;\">\n
    \   <style>\n        section {\n            font-size: 200%;\n        }\n\n\n
    \       .edit {\n            display: none;\n        }\n    </style>\n<section
    class=\"title\">\n    <h1 id=\"title\">\n        md_it_highlight_code.py\n    </h1>\n</section></article>\n
    \    </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       md_it_highlight_code.py\n    </h1>\n</section>    <section class=\"body\">\n
    \       <hr />\n<p>The <code>markata.plugins.md_it_highlight_code</code> plugin
    adds syntax highlighting to code blocks\nin your markdown content. It uses Pygments
    for highlighting and adds a copy button for\neasy code sharing.</p>\n<h1 id=\"installation\">Installation
    <a class=\"header-anchor\" href=\"#installation\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>This plugin is built-in
    and enabled by default through the 'default' plugin.\nIf you want to be explicit,
    you can add it to your list of plugins:</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"n\">hooks</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"p\">[</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;markata.plugins.md_it_highlight_code&quot;</span><span
    class=\"p\">,</span>\n<span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<h1
    id=\"uninstallation\">Uninstallation <a class=\"header-anchor\" href=\"#uninstallation\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>Since this plugin is
    included in the default plugin set, to disable it you must explicitly\nadd it
    to the disabled_hooks list if you are using the 'default' plugin:</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"n\">disabled_hooks</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"p\">[</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;markata.plugins.md_it_highlight_code&quot;</span><span
    class=\"p\">,</span>\n<span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<h1
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>Configure syntax highlighting
    in your <code>markata.toml</code>:</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata.highlight_code]</span>\n<span
    class=\"c1\"># Optional: Custom CSS class for code blocks</span>\n<span class=\"n\">code_class</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;highlight&quot;</span>\n\n<span class=\"c1\"># Optional: Custom
    CSS class for copy button</span>\n<span class=\"n\">copy_button_class</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;copy-button&quot;</span>\n\n<span class=\"c1\"># Optional:
    Custom copy button text</span>\n<span class=\"n\">copy_button_text</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;Copy&quot;</span>\n<span class=\"n\">copy_button_copied_text</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;Copied!&quot;</span>\n</pre></div>\n\n</pre>\n\n<h1 id=\"functionality\">Functionality
    <a class=\"header-anchor\" href=\"#functionality\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<h2 id=\"code-highlighting\">Code
    Highlighting <a class=\"header-anchor\" href=\"#code-highlighting\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>The plugin:</p>\n<ol>\n<li>Detects
    language from code block info</li>\n<li>Applies Pygments syntax highlighting</li>\n<li>Adds
    line numbers (optional)</li>\n<li>Wraps code in proper HTML structure</li>\n<li>Adds
    a copy-to-clipboard button</li>\n</ol>\n<h2 id=\"example-usage\">Example Usage
    <a class=\"header-anchor\" href=\"#example-usage\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>In markdown:</p>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"sb\">```python</span>\n<span
    class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">hello_world</span><span
    class=\"p\">():</span>\n    <span class=\"nb\">print</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;Hello, World!&quot;</span><span class=\"p\">)</span>\n<span
    class=\"sb\">```</span>\n\n<span class=\"sb\">```javascript</span>\n<span class=\"nx\">console</span><span
    class=\"p\">.</span><span class=\"nx\">log</span><span class=\"p\">(</span><span
    class=\"s1\">&#39;Hello, World!&#39;</span><span class=\"p\">);</span>\n<span
    class=\"sb\">```</span>\n</pre></div>\n\n</pre>\n\n<p>Generates HTML:</p>\n<pre
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
    class=\"nt\">div</span> <span class=\"na\">class</span><span class=\"o\">=</span><span
    class=\"s\">&quot;highlight&quot;</span><span class=\"p\">&gt;</span>\n  <span
    class=\"p\">&lt;</span><span class=\"nt\">button</span> <span class=\"na\">class</span><span
    class=\"o\">=</span><span class=\"s\">&quot;copy-button&quot;</span><span class=\"p\">&gt;</span>\n
    \   <span class=\"p\">&lt;</span><span class=\"nt\">svg</span><span class=\"p\">&gt;</span><span
    class=\"cm\">&lt;!-- Copy icon SVG --&gt;</span><span class=\"p\">&lt;/</span><span
    class=\"nt\">svg</span><span class=\"p\">&gt;</span>\n  <span class=\"p\">&lt;/</span><span
    class=\"nt\">button</span><span class=\"p\">&gt;</span>\n  <span class=\"p\">&lt;</span><span
    class=\"nt\">pre</span><span class=\"p\">&gt;&lt;</span><span class=\"nt\">code</span>
    <span class=\"na\">class</span><span class=\"o\">=</span><span class=\"s\">&quot;language-python&quot;</span><span
    class=\"p\">&gt;</span>\n    <span class=\"p\">&lt;</span><span class=\"nt\">span</span>
    <span class=\"na\">class</span><span class=\"o\">=</span><span class=\"s\">&quot;def&quot;</span><span
    class=\"p\">&gt;</span>def<span class=\"p\">&lt;/</span><span class=\"nt\">span</span><span
    class=\"p\">&gt;</span> <span class=\"p\">&lt;</span><span class=\"nt\">span</span>
    <span class=\"na\">class</span><span class=\"o\">=</span><span class=\"s\">&quot;name&quot;</span><span
    class=\"p\">&gt;</span>hello_world<span class=\"p\">&lt;/</span><span class=\"nt\">span</span><span
    class=\"p\">&gt;</span>():\n        <span class=\"p\">&lt;</span><span class=\"nt\">span</span>
    <span class=\"na\">class</span><span class=\"o\">=</span><span class=\"s\">&quot;keyword&quot;</span><span
    class=\"p\">&gt;</span>print<span class=\"p\">&lt;/</span><span class=\"nt\">span</span><span
    class=\"p\">&gt;</span>(<span class=\"p\">&lt;</span><span class=\"nt\">span</span>
    <span class=\"na\">class</span><span class=\"o\">=</span><span class=\"s\">&quot;string&quot;</span><span
    class=\"p\">&gt;</span>&quot;Hello, World!&quot;<span class=\"p\">&lt;/</span><span
    class=\"nt\">span</span><span class=\"p\">&gt;</span>)\n  <span class=\"p\">&lt;/</span><span
    class=\"nt\">code</span><span class=\"p\">&gt;&lt;/</span><span class=\"nt\">pre</span><span
    class=\"p\">&gt;</span>\n<span class=\"p\">&lt;/</span><span class=\"nt\">div</span><span
    class=\"p\">&gt;</span>\n</pre></div>\n\n</pre>\n\n<h2 id=\"supported-languages\">Supported
    Languages <a class=\"header-anchor\" href=\"#supported-languages\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>The plugin supports
    all languages that Pygments recognizes, including:</p>\n<ul>\n<li>Python</li>\n<li>JavaScript</li>\n<li>HTML/CSS</li>\n<li>Bash/Shell</li>\n<li>And
    many more</li>\n</ul>\n<h2 id=\"copy-button\">Copy Button <a class=\"header-anchor\"
    href=\"#copy-button\"><svg class=\"heading-permalink\" aria-hidden=\"true\" fill=\"currentColor\"
    focusable=\"false\" height=\"1em\" viewBox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Features:</p>\n<ul>\n<li>Hover-to-show
    button</li>\n<li>Click-to-copy functionality</li>\n<li>Success feedback</li>\n<li>Accessible
    keyboard support</li>\n</ul>\n<h2 id=\"dependencies\">Dependencies <a class=\"header-anchor\"
    href=\"#dependencies\"><svg class=\"heading-permalink\" aria-hidden=\"true\" fill=\"currentColor\"
    focusable=\"false\" height=\"1em\" viewBox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>This plugin depends
    on:</p>\n<ul>\n<li>Pygments for syntax highlighting</li>\n<li>markdown-it-py for
    markdown parsing</li>\n</ul>\n<hr />\n<div class=\"admonition function\">\n<p
    class=\"admonition-title\">Function</p>\n<h2 id=\"highlight_code\" class=\"admonition-title\"
    style=\"margin: 0; padding: .5rem 1rem;\">highlight_code <em class=\"small\">function</em></h2>\n<p>Code
    highlighter for markdown-it-py.</p>\n</div>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">highlight_code <em class='small'>source</em></p>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">def</span><span
    class=\"w\"> </span><span class=\"nf\">highlight_code</span><span class=\"p\">(</span><span
    class=\"n\">code</span><span class=\"p\">,</span> <span class=\"n\">name</span><span
    class=\"p\">,</span> <span class=\"n\">attrs</span><span class=\"p\">,</span>
    <span class=\"n\">markata</span><span class=\"o\">=</span><span class=\"kc\">None</span><span
    class=\"p\">):</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Code
    highlighter for markdown-it-py.&quot;&quot;&quot;</span>\n    <span class=\"c1\">#
    from pygments import highlight</span>\n    <span class=\"c1\"># from pygments.formatters
    import HtmlFormatter</span>\n    <span class=\"c1\"># from pygments.lexers import
    ClassNotFound, get_lexer_by_name</span>\n\n    <span class=\"k\">try</span><span
    class=\"p\">:</span>\n        <span class=\"n\">lexer</span> <span class=\"o\">=</span>
    <span class=\"n\">get_lexer_by_name</span><span class=\"p\">(</span><span class=\"n\">name</span>
    <span class=\"ow\">or</span> <span class=\"s2\">&quot;text&quot;</span><span class=\"p\">)</span>\n
    \   <span class=\"k\">except</span> <span class=\"n\">ClassNotFound</span><span
    class=\"p\">:</span>\n        <span class=\"n\">lexer</span> <span class=\"o\">=</span>
    <span class=\"n\">get_lexer_by_name</span><span class=\"p\">(</span><span class=\"s2\">&quot;text&quot;</span><span
    class=\"p\">)</span>\n\n    <span class=\"kn\">import</span><span class=\"w\">
    </span><span class=\"nn\">re</span>\n\n    <span class=\"n\">pattern</span> <span
    class=\"o\">=</span> <span class=\"sa\">r</span><span class=\"s1\">&#39;(\\w+)\\s*=\\s*(&quot;.*?&quot;|\\S+)&#39;</span>\n
    \   <span class=\"n\">matches</span> <span class=\"o\">=</span> <span class=\"n\">re</span><span
    class=\"o\">.</span><span class=\"n\">findall</span><span class=\"p\">(</span><span
    class=\"n\">pattern</span><span class=\"p\">,</span> <span class=\"n\">attrs</span><span
    class=\"p\">)</span>\n    <span class=\"n\">attrs</span> <span class=\"o\">=</span>
    <span class=\"nb\">dict</span><span class=\"p\">(</span><span class=\"n\">matches</span><span
    class=\"p\">)</span>\n\n    <span class=\"k\">if</span> <span class=\"n\">attrs</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;hl_lines&quot;</span><span class=\"p\">):</span>\n        <span
    class=\"n\">formatter</span> <span class=\"o\">=</span> <span class=\"n\">HtmlFormatter</span><span
    class=\"p\">(</span><span class=\"n\">hl_lines</span><span class=\"o\">=</span><span
    class=\"n\">attrs</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;hl_lines&quot;</span><span class=\"p\">))</span>\n
    \   <span class=\"k\">else</span><span class=\"p\">:</span>\n        <span class=\"n\">formatter</span>
    <span class=\"o\">=</span> <span class=\"n\">HtmlFormatter</span><span class=\"p\">()</span>\n\n
    \   <span class=\"n\">copy_button</span> <span class=\"o\">=</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;&quot;&quot;&lt;button class=&#39;copy&#39; title=&#39;copy
    code to clipboard&#39; onclick=&quot;navigator.clipboard.writeText(this.parentElement.parentElement.querySelector(&#39;pre&#39;).textContent)&quot;&gt;</span><span
    class=\"si\">{</span><span class=\"n\">COPY_ICON</span><span class=\"si\">}</span><span
    class=\"s2\">&lt;/button&gt;&quot;&quot;&quot;</span>\n\n    <span class=\"kn\">from</span><span
    class=\"w\"> </span><span class=\"nn\">markdown_it</span><span class=\"w\"> </span><span
    class=\"kn\">import</span> <span class=\"n\">MarkdownIt</span>\n\n    <span class=\"n\">md</span>
    <span class=\"o\">=</span> <span class=\"n\">MarkdownIt</span><span class=\"p\">(</span>\n
    \       <span class=\"s2\">&quot;commonmark&quot;</span><span class=\"p\">,</span>\n
    \       <span class=\"p\">{</span>\n            <span class=\"s2\">&quot;html&quot;</span><span
    class=\"p\">:</span> <span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \           <span class=\"s2\">&quot;typographer&quot;</span><span class=\"p\">:</span>
    <span class=\"kc\">True</span><span class=\"p\">,</span>\n        <span class=\"p\">},</span>\n
    \   <span class=\"p\">)</span>\n\n    <span class=\"k\">if</span> <span class=\"n\">attrs</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;help&quot;</span><span class=\"p\">):</span>\n        <span
    class=\"n\">help</span> <span class=\"o\">=</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;&quot;&quot;</span>\n<span class=\"s2\">        &lt;a href=</span><span
    class=\"si\">{</span><span class=\"n\">attrs</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;help&quot;</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">strip</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;&lt;&quot;</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">strip</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;&gt;&quot;</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">strip</span><span class=\"p\">(</span><span class=\"s1\">&#39;&quot;&#39;</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">strip</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;&#39;&quot;</span><span class=\"p\">)</span><span
    class=\"si\">}</span><span class=\"s2\"> title=&#39;help link&#39; class=&#39;help&#39;&gt;</span><span
    class=\"si\">{</span><span class=\"n\">HELP_ICON</span><span class=\"si\">}</span><span
    class=\"s2\">&lt;/a&gt;</span>\n<span class=\"s2\">        &quot;&quot;&quot;</span>\n
    \   <span class=\"k\">else</span><span class=\"p\">:</span>\n        <span class=\"n\">help</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;&quot;</span>\n    <span class=\"k\">if</span>
    <span class=\"n\">attrs</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;title&quot;</span><span class=\"p\">):</span>\n
    \       <span class=\"n\">file</span> <span class=\"o\">=</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;&quot;&quot;</span>\n<span class=\"s2\">&lt;div class=&#39;filepath&#39;&gt;</span>\n<span
    class=\"si\">{</span><span class=\"n\">md</span><span class=\"o\">.</span><span
    class=\"n\">render</span><span class=\"p\">(</span><span class=\"n\">attrs</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;title&quot;</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">strip</span><span class=\"p\">(</span><span class=\"s1\">&#39;&quot;&#39;</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">strip</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;&#39;&quot;</span><span class=\"p\">))</span><span
    class=\"si\">}</span>\n<span class=\"s2\">&lt;div class=&#39;right&#39;&gt;</span>\n<span
    class=\"si\">{</span><span class=\"n\">help</span><span class=\"si\">}</span>\n<span
    class=\"si\">{</span><span class=\"n\">copy_button</span><span class=\"si\">}</span>\n<span
    class=\"s2\">&lt;/div&gt;</span>\n<span class=\"s2\">&lt;/div&gt;</span>\n<span
    class=\"s2\">&quot;&quot;&quot;</span>\n    <span class=\"k\">else</span><span
    class=\"p\">:</span>\n        <span class=\"n\">file</span> <span class=\"o\">=</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;&quot;&quot;</span>\n<span
    class=\"s2\">&lt;div class=&#39;copy-wrapper&#39;&gt;</span>\n<span class=\"si\">{</span><span
    class=\"n\">help</span><span class=\"si\">}</span>\n<span class=\"si\">{</span><span
    class=\"n\">copy_button</span><span class=\"si\">}</span>\n<span class=\"s2\">&lt;/div&gt;</span>\n<span
    class=\"s2\">        &quot;&quot;&quot;</span>\n    <span class=\"k\">return</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;&quot;&quot;&lt;pre class=&#39;wrapper&#39;&gt;</span>\n<span
    class=\"si\">{</span><span class=\"n\">file</span><span class=\"si\">}</span>\n<span
    class=\"si\">{</span><span class=\"n\">highlight</span><span class=\"p\">(</span><span
    class=\"n\">code</span><span class=\"p\">,</span><span class=\"w\"> </span><span
    class=\"n\">lexer</span><span class=\"p\">,</span><span class=\"w\"> </span><span
    class=\"n\">formatter</span><span class=\"p\">)</span><span class=\"si\">}</span>\n<span
    class=\"s2\">&lt;/pre&gt;</span>\n<span class=\"s2\">&quot;&quot;&quot;</span>\n</pre></div>\n\n</pre>\n\n</div>\n\n
    \   </section>\n</article>"
  raw.md: ''
published: false
slug: markata/plugins/md-it-highlight-code
title: md_it_highlight_code.py


---

---

The `markata.plugins.md_it_highlight_code` plugin adds syntax highlighting to code blocks
in your markdown content. It uses Pygments for highlighting and adds a copy button for
easy code sharing.

# Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.md_it_highlight_code",
]
```

# Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.md_it_highlight_code",
]
```

# Configuration

Configure syntax highlighting in your `markata.toml`:

```toml
[markata.highlight_code]
# Optional: Custom CSS class for code blocks
code_class = "highlight"

# Optional: Custom CSS class for copy button
copy_button_class = "copy-button"

# Optional: Custom copy button text
copy_button_text = "Copy"
copy_button_copied_text = "Copied!"
```

# Functionality

## Code Highlighting

The plugin:
1. Detects language from code block info
2. Applies Pygments syntax highlighting
3. Adds line numbers (optional)
4. Wraps code in proper HTML structure
5. Adds a copy-to-clipboard button

## Example Usage

In markdown:
````markdown
```python
def hello_world():
    print("Hello, World!")
```

```javascript
console.log('Hello, World!');
```
````

Generates HTML:
```html
<div class="highlight">
  <button class="copy-button">
    <svg><!-- Copy icon SVG --></svg>
  </button>
  <pre><code class="language-python">
    <span class="def">def</span> <span class="name">hello_world</span>():
        <span class="keyword">print</span>(<span class="string">"Hello, World!"</span>)
  </code></pre>
</div>
```

## Supported Languages

The plugin supports all languages that Pygments recognizes, including:
- Python
- JavaScript
- HTML/CSS
- Bash/Shell
- And many more

## Copy Button

Features:
- Hover-to-show button
- Click-to-copy functionality
- Success feedback
- Accessible keyboard support

## Dependencies

This plugin depends on:
- Pygments for syntax highlighting
- markdown-it-py for markdown parsing

---



!!! function
    <h2 id="highlight_code" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">highlight_code <em class="small">function</em></h2>

    Code highlighter for markdown-it-py.

???+ source "highlight_code <em class='small'>source</em>"
    ```python
    def highlight_code(code, name, attrs, markata=None):
        """Code highlighter for markdown-it-py."""
        # from pygments import highlight
        # from pygments.formatters import HtmlFormatter
        # from pygments.lexers import ClassNotFound, get_lexer_by_name

        try:
            lexer = get_lexer_by_name(name or "text")
        except ClassNotFound:
            lexer = get_lexer_by_name("text")

        import re

        pattern = r'(\w+)\s*=\s*(".*?"|\S+)'
        matches = re.findall(pattern, attrs)
        attrs = dict(matches)

        if attrs.get("hl_lines"):
            formatter = HtmlFormatter(hl_lines=attrs.get("hl_lines"))
        else:
            formatter = HtmlFormatter()

        copy_button = f"""<button class='copy' title='copy code to clipboard' onclick="navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)">{COPY_ICON}</button>"""

        from markdown_it import MarkdownIt

        md = MarkdownIt(
            "commonmark",
            {
                "html": True,
                "typographer": True,
            },
        )

        if attrs.get("help"):
            help = f"""
            <a href={attrs.get("help").strip("<").strip(">").strip('"').strip("'")} title='help link' class='help'>{HELP_ICON}</a>
            """
        else:
            help = ""
        if attrs.get("title"):
            file = f"""
    <div class='filepath'>
    {md.render(attrs.get("title").strip('"').strip("'"))}
    <div class='right'>
    {help}
    {copy_button}
    </div>
    </div>
    """
        else:
            file = f"""
    <div class='copy-wrapper'>
    {help}
    {copy_button}
    </div>
            """
        return f"""<pre class='wrapper'>
    {file}
    {highlight(code, lexer, formatter)}
    </pre>
    """
    ```