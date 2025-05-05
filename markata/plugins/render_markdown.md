---
content: "---\n\nThe `markata.plugins.render_markdown` plugin converts markdown content
  to HTML.\nThis plugin is essential for rendering markdown files loaded by the `load`
  plugin.\n\n## Installation\n\nThis plugin is built-in and enabled by default through
  the 'default' plugin.\nIf you want to be explicit, you can add it to your list of
  plugins:\n\n```toml\nhooks = [\n    \"markata.plugins.render_markdown\",\n]\n```\n\n##
  Uninstallation\n\nSince this plugin is included in the default plugin set, to disable
  it you must explicitly\nadd it to the disabled_hooks list if you are using the 'default'
  plugin:\n\n```toml\ndisabled_hooks = [\n    \"markata.plugins.render_markdown\",\n]\n```\n\nNote:
  Disabling this plugin will prevent markdown files from being rendered to HTML.\n\n##
  Configuration\n\n## Markdown Backend Selection\n\nChoose from 3 supported markdown
  backends by setting `markdown_backend` in your `markata.toml`:\n\n```toml\n## choose
  your markdown backend\n# markdown_backend='markdown'      # Python-Markdown\n# markdown_backend='markdown2'
  \    # markdown2\nmarkdown_backend='markdown-it-py'  # markdown-it-py (default)\n```\n\n##
  Backend-Specific Configuration\n\n### markdown-it-py\n\nConfigure markdown-it-py
  behavior in your `markata.toml`:\n\n```toml\n[markata.markdown_it_py]\n# Set the
  flavor - options: 'zero', 'commonmark', 'gfm-like'\nconfig = 'gfm-like'\n\n# Enable
  specific plugins\nenable = [\n    'table',\n    'strikethrough',\n    'footnote',\n]\n\n#
  Disable specific plugins\ndisable = [\n    'linkify',\n]\n\n# Configure plugins\n[markata.markdown_it_py.plugins.footnote]\n#
  Plugin-specific settings here\n```\n\nRead more about markdown-it-py settings in
  their [documentation](https://markdown-it-py.readthedocs.io/en/latest/).\n\n## Cache
  Configuration\n\nControl markdown rendering cache:\n\n```toml\n[markata.render_markdown]\ncache_expire
  = 3600  # Cache expiration in seconds\n```\n\n## Functionality\n\n## Registered
  Attributes\n\nThe plugin registers the following attributes on Post objects:\n-
  `html`: The rendered HTML content from the markdown source\n\n## Dependencies\n\nThis
  plugin depends on:\n- One of: python-markdown, markdown2, or markdown-it-py (based
  on configuration)\n- The `load` plugin to provide markdown content for rendering\n\n---\n\n!!!
  function\n    <h2 id=\"configure\" class=\"admonition-title\" style=\"margin: 0;
  padding: .5rem 1rem;\">configure <em class=\"small\">function</em></h2>\n\n    Sets
  up a markdown instance as md\n\n???+ source \"configure <em class='small'>source</em>\"\n
  \   ```python\n    def configure(markata: \"Markata\") -> None:\n        \"Sets
  up a markdown instance as md\"\n        # if \"markdown_extensions\" not in markata.config:\n
  \       #     markdown_extensions = [\"\"]\n        # if isinstance(markata.config[\"markdown_extensions\"],
  str):\n        #     markdown_extensions = [markata.config[\"markdown_extensions\"]]\n
  \       # if isinstance(markata.config[\"markdown_extensions\"], list):\n        #
  \    markdown_extensions = markata.config[\"markdown_extensions\"]\n        # else:\n
  \       #     raise TypeError(\"markdown_extensions should be List[str]\")\n\n        #
  markata.markdown_extensions = [*DEFAULT_MD_EXTENSIONS, *markdown_extensions]\n\n
  \       if (\n            markata.config.get(\"markdown_backend\", \"\")\n            .lower()\n
  \           .replace(\" \", \"-\")\n            .replace(\"_\", \"-\")\n            ==
  \"markdown-it-py\"\n        ):\n            from markdown_it import MarkdownIt\n\n
  \           config_update = markata.config.get(\"markdown_it_py\", {}).get(\n                \"options_update\",\n
  \               {\n                    \"linkify\": True,\n                    \"html\":
  True,\n                    \"typographer\": True,\n                    \"highlight\":
  highlight_code,\n                },\n            )\n            if isinstance(config_update.get(\"highlight\"),
  str):\n                module = config_update[\"highlight\"].split(\":\")[0]\n                func
  = config_update[\"highlight\"].split(\":\")[1]\n                config_update[\"highlight\"]
  = getattr(\n                    importlib.import_module(module),\n                    func,\n
  \               )\n\n            markata.md = MarkdownIt(\n                markata.config.get(\"markdown_it_py\",
  {}).get(\"config\", \"gfm-like\"),\n                config_update,\n            )\n
  \           for plugin in markata.config.get(\"markdown_it_py\", {}).get(\"enable\",
  []):\n                markata.md.enable(plugin)\n            for plugin in markata.config.get(\"markdown_it_py\",
  {}).get(\"disable\", []):\n                markata.md.disable(plugin)\n\n            plugins
  = copy.deepcopy(\n                markata.config.get(\"markdown_it_py\", {}).get(\"plugins\",
  []),\n            )\n            for plugin in plugins:\n                if isinstance(plugin[\"plugin\"],
  str):\n                    plugin[\"plugin_str\"] = plugin[\"plugin\"]\n                    plugin_module
  = plugin[\"plugin\"].split(\":\")[0]\n                    plugin_func = plugin[\"plugin\"].split(\":\")[1]\n
  \                   plugin[\"plugin\"] = getattr(\n                        importlib.import_module(plugin_module),\n
  \                       plugin_func,\n                    )\n                plugin[\"config\"]
  = plugin.get(\"config\", {})\n                for k, _v in plugin[\"config\"].items():\n
  \                   if k == \"markata\":\n                        plugin[\"config\"][k]
  = markata\n\n                markata.md = markata.md.use(plugin[\"plugin\"], **plugin[\"config\"])\n\n
  \           markata.md.convert = markata.md.render\n            markata.md.toc =
  \"\"\n        elif (\n            markata.config.get(\"markdown_backend\", \"\")\n
  \           .lower()\n            .replace(\" \", \"-\")\n            .replace(\"_\",
  \"-\")\n            == \"markdown2\"\n        ):\n            import markdown2\n\n
  \           markata.md = markdown2.Markdown(\n                extras=markata.config.render_markdown.extensions\n
  \           )\n            markata.md.toc = \"\"\n        else:\n            import
  markdown\n\n            markata.md = markdown.Markdown(\n                extensions=markata.config.render_markdown.extensions\n
  \           )\n    ```\n!!! function\n    <h2 id=\"render\" class=\"admonition-title\"
  style=\"margin: 0; padding: .5rem 1rem;\">render <em class=\"small\">function</em></h2>\n\n
  \   Render markdown content in parallel.\n\n???+ source \"render <em class='small'>source</em>\"\n
  \   ```python\n    def render(markata: \"Markata\") -> None:\n        \"\"\"Render
  markdown content in parallel.\"\"\"\n        config = markata.config.render_markdown\n
  \       articles = list(markata.filter(\"not skip\"))\n\n        with markata.cache
  as cache:\n            with concurrent.futures.ThreadPoolExecutor() as executor:\n
  \               render_func = partial(render_article_parallel, markata, config,
  cache)\n                args_list = [(article,) for article in articles]\n\n                for
  article, html in executor.map(render_func, args_list):\n                    article.html
  = html\n                    article.article_html = copy.deepcopy(html)\n\n        markata.rendered_posts
  = markata.posts\n    ```"
date: 2025-05-05
description: "The   plugin converts markdown content to HTML. This plugin is essential
  for rendering markdown files loaded by the   plugin. Installation This plugin is\u2026"
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>render_markdown.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"The   plugin converts markdown content
    to HTML. This plugin is essential for rendering markdown files loaded by the   plugin.
    Installation This plugin is\u2026\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n
    \       <meta property=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    <head>\n        <script src=\"https://unpkg.com/@tailwindcss/browser@4\"></script>\n<title>render_markdown.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"The   plugin converts markdown content
    to HTML. This plugin is essential for rendering markdown files loaded by the   plugin.
    Installation This plugin is\u2026\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        render_markdown.py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <hr />\n<p>The <code>markata.plugins.render_markdown</code>
    plugin converts markdown content to HTML.\nThis plugin is essential for rendering
    markdown files loaded by the <code>load</code> plugin.</p>\n<h2 id=\"installation\">Installation
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>This plugin is built-in
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
    class=\"p\">[</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;markata.plugins.render_markdown&quot;</span><span
    class=\"p\">,</span>\n<span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<h2
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Since this plugin is
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
    class=\"p\">[</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;markata.plugins.render_markdown&quot;</span><span
    class=\"p\">,</span>\n<span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>Note:
    Disabling this plugin will prevent markdown files from being rendered to HTML.</p>\n<h2
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<h2 id=\"markdown-backend-selection\">Markdown
    Backend Selection <a class=\"header-anchor\" href=\"#markdown-backend-selection\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Choose from 3 supported
    markdown backends by setting <code>markdown_backend</code> in your <code>markata.toml</code>:</p>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"c1\">## choose
    your markdown backend</span>\n<span class=\"c1\"># markdown_backend=&#39;markdown&#39;
    \     # Python-Markdown</span>\n<span class=\"c1\"># markdown_backend=&#39;markdown2&#39;
    \    # markdown2</span>\n<span class=\"n\">markdown_backend</span><span class=\"o\">=</span><span
    class=\"s1\">&#39;markdown-it-py&#39;</span><span class=\"w\">  </span><span class=\"c1\">#
    markdown-it-py (default)</span>\n</pre></div>\n\n</pre>\n\n<h2 id=\"backend-specific-configuration\">Backend-Specific
    Configuration <a class=\"header-anchor\" href=\"#backend-specific-configuration\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<h3>markdown-it-py</h3>\n<p>Configure
    markdown-it-py behavior in your <code>markata.toml</code>:</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata.markdown_it_py]</span>\n<span
    class=\"c1\"># Set the flavor - options: &#39;zero&#39;, &#39;commonmark&#39;,
    &#39;gfm-like&#39;</span>\n<span class=\"n\">config</span><span class=\"w\"> </span><span
    class=\"o\">=</span><span class=\"w\"> </span><span class=\"s1\">&#39;gfm-like&#39;</span>\n\n<span
    class=\"c1\"># Enable specific plugins</span>\n<span class=\"n\">enable</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"p\">[</span>\n<span class=\"w\">    </span><span class=\"s1\">&#39;table&#39;</span><span
    class=\"p\">,</span>\n<span class=\"w\">    </span><span class=\"s1\">&#39;strikethrough&#39;</span><span
    class=\"p\">,</span>\n<span class=\"w\">    </span><span class=\"s1\">&#39;footnote&#39;</span><span
    class=\"p\">,</span>\n<span class=\"p\">]</span>\n\n<span class=\"c1\"># Disable
    specific plugins</span>\n<span class=\"n\">disable</span><span class=\"w\"> </span><span
    class=\"o\">=</span><span class=\"w\"> </span><span class=\"p\">[</span>\n<span
    class=\"w\">    </span><span class=\"s1\">&#39;linkify&#39;</span><span class=\"p\">,</span>\n<span
    class=\"p\">]</span>\n\n<span class=\"c1\"># Configure plugins</span>\n<span class=\"k\">[markata.markdown_it_py.plugins.footnote]</span>\n<span
    class=\"c1\"># Plugin-specific settings here</span>\n</pre></div>\n\n</pre>\n\n<p>Read
    more about markdown-it-py settings in their <a href=\"https://markdown-it-py.readthedocs.io/en/latest/\">documentation</a>.</p>\n<h2
    id=\"cache-configuration\">Cache Configuration <a class=\"header-anchor\" href=\"#cache-configuration\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Control markdown rendering
    cache:</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button class='copy'
    title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata.render_markdown]</span>\n<span
    class=\"n\">cache_expire</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"mi\">3600</span><span class=\"w\">  </span><span
    class=\"c1\"># Cache expiration in seconds</span>\n</pre></div>\n\n</pre>\n\n<h2
    id=\"functionality\">Functionality <a class=\"header-anchor\" href=\"#functionality\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<h2 id=\"registered-attributes\">Registered
    Attributes <a class=\"header-anchor\" href=\"#registered-attributes\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>The plugin registers
    the following attributes on Post objects:</p>\n<ul>\n<li><code>html</code>: The
    rendered HTML content from the markdown source</li>\n</ul>\n<h2 id=\"dependencies\">Dependencies
    <a class=\"header-anchor\" href=\"#dependencies\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>This plugin depends
    on:</p>\n<ul>\n<li>One of: python-markdown, markdown2, or markdown-it-py (based
    on configuration)</li>\n<li>The <code>load</code> plugin to provide markdown content
    for rendering</li>\n</ul>\n<hr />\n<div class=\"admonition function\">\n<p class=\"admonition-title\">Function</p>\n<h2
    id=\"configure\" class=\"admonition-title\" style=\"margin: 0; padding: .5rem
    1rem;\">configure <em class=\"small\">function</em></h2>\n<p>Sets up a markdown
    instance as md</p>\n</div>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">configure <em class='small'>source</em></p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">def</span><span
    class=\"w\"> </span><span class=\"nf\">configure</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n    <span class=\"s2\">&quot;Sets up a markdown instance
    as md&quot;</span>\n    <span class=\"c1\"># if &quot;markdown_extensions&quot;
    not in markata.config:</span>\n    <span class=\"c1\">#     markdown_extensions
    = [&quot;&quot;]</span>\n    <span class=\"c1\"># if isinstance(markata.config[&quot;markdown_extensions&quot;],
    str):</span>\n    <span class=\"c1\">#     markdown_extensions = [markata.config[&quot;markdown_extensions&quot;]]</span>\n
    \   <span class=\"c1\"># if isinstance(markata.config[&quot;markdown_extensions&quot;],
    list):</span>\n    <span class=\"c1\">#     markdown_extensions = markata.config[&quot;markdown_extensions&quot;]</span>\n
    \   <span class=\"c1\"># else:</span>\n    <span class=\"c1\">#     raise TypeError(&quot;markdown_extensions
    should be List[str]&quot;)</span>\n\n    <span class=\"c1\"># markata.markdown_extensions
    = [*DEFAULT_MD_EXTENSIONS, *markdown_extensions]</span>\n\n    <span class=\"k\">if</span>
    <span class=\"p\">(</span>\n        <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;markdown_backend&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;&quot;</span><span class=\"p\">)</span>\n
    \       <span class=\"o\">.</span><span class=\"n\">lower</span><span class=\"p\">()</span>\n
    \       <span class=\"o\">.</span><span class=\"n\">replace</span><span class=\"p\">(</span><span
    class=\"s2\">&quot; &quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;-&quot;</span><span
    class=\"p\">)</span>\n        <span class=\"o\">.</span><span class=\"n\">replace</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;_&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;-&quot;</span><span class=\"p\">)</span>\n        <span
    class=\"o\">==</span> <span class=\"s2\">&quot;markdown-it-py&quot;</span>\n    <span
    class=\"p\">):</span>\n        <span class=\"kn\">from</span><span class=\"w\">
    </span><span class=\"nn\">markdown_it</span><span class=\"w\"> </span><span class=\"kn\">import</span>
    <span class=\"n\">MarkdownIt</span>\n\n        <span class=\"n\">config_update</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;markdown_it_py&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">{})</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span>\n            <span class=\"s2\">&quot;options_update&quot;</span><span
    class=\"p\">,</span>\n            <span class=\"p\">{</span>\n                <span
    class=\"s2\">&quot;linkify&quot;</span><span class=\"p\">:</span> <span class=\"kc\">True</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;html&quot;</span><span
    class=\"p\">:</span> <span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \               <span class=\"s2\">&quot;typographer&quot;</span><span class=\"p\">:</span>
    <span class=\"kc\">True</span><span class=\"p\">,</span>\n                <span
    class=\"s2\">&quot;highlight&quot;</span><span class=\"p\">:</span> <span class=\"n\">highlight_code</span><span
    class=\"p\">,</span>\n            <span class=\"p\">},</span>\n        <span class=\"p\">)</span>\n
    \       <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">config_update</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;highlight&quot;</span><span
    class=\"p\">),</span> <span class=\"nb\">str</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">module</span> <span class=\"o\">=</span> <span class=\"n\">config_update</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;highlight&quot;</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">split</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;:&quot;</span><span class=\"p\">)[</span><span class=\"mi\">0</span><span
    class=\"p\">]</span>\n            <span class=\"n\">func</span> <span class=\"o\">=</span>
    <span class=\"n\">config_update</span><span class=\"p\">[</span><span class=\"s2\">&quot;highlight&quot;</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">split</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;:&quot;</span><span class=\"p\">)[</span><span
    class=\"mi\">1</span><span class=\"p\">]</span>\n            <span class=\"n\">config_update</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;highlight&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"nb\">getattr</span><span class=\"p\">(</span>\n
    \               <span class=\"n\">importlib</span><span class=\"o\">.</span><span
    class=\"n\">import_module</span><span class=\"p\">(</span><span class=\"n\">module</span><span
    class=\"p\">),</span>\n                <span class=\"n\">func</span><span class=\"p\">,</span>\n
    \           <span class=\"p\">)</span>\n\n        <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">md</span> <span class=\"o\">=</span> <span
    class=\"n\">MarkdownIt</span><span class=\"p\">(</span>\n            <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;markdown_it_py&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">{})</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;config&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;gfm-like&quot;</span><span class=\"p\">),</span>\n
    \           <span class=\"n\">config_update</span><span class=\"p\">,</span>\n
    \       <span class=\"p\">)</span>\n        <span class=\"k\">for</span> <span
    class=\"n\">plugin</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;markdown_it_py&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">{})</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;enable&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">[]):</span>\n            <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">md</span><span class=\"o\">.</span><span
    class=\"n\">enable</span><span class=\"p\">(</span><span class=\"n\">plugin</span><span
    class=\"p\">)</span>\n        <span class=\"k\">for</span> <span class=\"n\">plugin</span>
    <span class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;markdown_it_py&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">{})</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;disable&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">[]):</span>\n            <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">md</span><span class=\"o\">.</span><span
    class=\"n\">disable</span><span class=\"p\">(</span><span class=\"n\">plugin</span><span
    class=\"p\">)</span>\n\n        <span class=\"n\">plugins</span> <span class=\"o\">=</span>
    <span class=\"n\">copy</span><span class=\"o\">.</span><span class=\"n\">deepcopy</span><span
    class=\"p\">(</span>\n            <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;markdown_it_py&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">{})</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;plugins&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">[]),</span>\n        <span class=\"p\">)</span>\n
    \       <span class=\"k\">for</span> <span class=\"n\">plugin</span> <span class=\"ow\">in</span>
    <span class=\"n\">plugins</span><span class=\"p\">:</span>\n            <span
    class=\"k\">if</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">plugin</span><span class=\"p\">[</span><span class=\"s2\">&quot;plugin&quot;</span><span
    class=\"p\">],</span> <span class=\"nb\">str</span><span class=\"p\">):</span>\n
    \               <span class=\"n\">plugin</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;plugin_str&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"n\">plugin</span><span class=\"p\">[</span><span class=\"s2\">&quot;plugin&quot;</span><span
    class=\"p\">]</span>\n                <span class=\"n\">plugin_module</span> <span
    class=\"o\">=</span> <span class=\"n\">plugin</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;plugin&quot;</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">split</span><span class=\"p\">(</span><span class=\"s2\">&quot;:&quot;</span><span
    class=\"p\">)[</span><span class=\"mi\">0</span><span class=\"p\">]</span>\n                <span
    class=\"n\">plugin_func</span> <span class=\"o\">=</span> <span class=\"n\">plugin</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;plugin&quot;</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">split</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;:&quot;</span><span class=\"p\">)[</span><span class=\"mi\">1</span><span
    class=\"p\">]</span>\n                <span class=\"n\">plugin</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;plugin&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"nb\">getattr</span><span class=\"p\">(</span>\n                    <span
    class=\"n\">importlib</span><span class=\"o\">.</span><span class=\"n\">import_module</span><span
    class=\"p\">(</span><span class=\"n\">plugin_module</span><span class=\"p\">),</span>\n
    \                   <span class=\"n\">plugin_func</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">)</span>\n            <span class=\"n\">plugin</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;config&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"n\">plugin</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;config&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">{})</span>\n            <span class=\"k\">for</span>
    <span class=\"n\">k</span><span class=\"p\">,</span> <span class=\"n\">_v</span>
    <span class=\"ow\">in</span> <span class=\"n\">plugin</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;config&quot;</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">items</span><span class=\"p\">():</span>\n                <span class=\"k\">if</span>
    <span class=\"n\">k</span> <span class=\"o\">==</span> <span class=\"s2\">&quot;markata&quot;</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">plugin</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;config&quot;</span><span class=\"p\">][</span><span
    class=\"n\">k</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"n\">markata</span>\n\n            <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">md</span> <span class=\"o\">=</span> <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">md</span><span
    class=\"o\">.</span><span class=\"n\">use</span><span class=\"p\">(</span><span
    class=\"n\">plugin</span><span class=\"p\">[</span><span class=\"s2\">&quot;plugin&quot;</span><span
    class=\"p\">],</span> <span class=\"o\">**</span><span class=\"n\">plugin</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;config&quot;</span><span class=\"p\">])</span>\n\n
    \       <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">md</span><span
    class=\"o\">.</span><span class=\"n\">convert</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">md</span><span
    class=\"o\">.</span><span class=\"n\">render</span>\n        <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">md</span><span class=\"o\">.</span><span
    class=\"n\">toc</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;&quot;</span>\n
    \   <span class=\"k\">elif</span> <span class=\"p\">(</span>\n        <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;markdown_backend&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;&quot;</span><span class=\"p\">)</span>\n
    \       <span class=\"o\">.</span><span class=\"n\">lower</span><span class=\"p\">()</span>\n
    \       <span class=\"o\">.</span><span class=\"n\">replace</span><span class=\"p\">(</span><span
    class=\"s2\">&quot; &quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;-&quot;</span><span
    class=\"p\">)</span>\n        <span class=\"o\">.</span><span class=\"n\">replace</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;_&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;-&quot;</span><span class=\"p\">)</span>\n        <span
    class=\"o\">==</span> <span class=\"s2\">&quot;markdown2&quot;</span>\n    <span
    class=\"p\">):</span>\n        <span class=\"kn\">import</span><span class=\"w\">
    </span><span class=\"nn\">markdown2</span>\n\n        <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">md</span> <span class=\"o\">=</span> <span
    class=\"n\">markdown2</span><span class=\"o\">.</span><span class=\"n\">Markdown</span><span
    class=\"p\">(</span>\n            <span class=\"n\">extras</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">render_markdown</span><span class=\"o\">.</span><span
    class=\"n\">extensions</span>\n        <span class=\"p\">)</span>\n        <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">md</span><span
    class=\"o\">.</span><span class=\"n\">toc</span> <span class=\"o\">=</span> <span
    class=\"s2\">&quot;&quot;</span>\n    <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \       <span class=\"kn\">import</span><span class=\"w\"> </span><span class=\"nn\">markdown</span>\n\n
    \       <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">md</span>
    <span class=\"o\">=</span> <span class=\"n\">markdown</span><span class=\"o\">.</span><span
    class=\"n\">Markdown</span><span class=\"p\">(</span>\n            <span class=\"n\">extensions</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">render_markdown</span><span
    class=\"o\">.</span><span class=\"n\">extensions</span>\n        <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition function\">\n<p class=\"admonition-title\">Function</p>\n<h2
    id=\"render\" class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">render
    <em class=\"small\">function</em></h2>\n<p>Render markdown content in parallel.</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">render
    <em class='small'>source</em></p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">def</span><span
    class=\"w\"> </span><span class=\"nf\">render</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Render
    markdown content in parallel.&quot;&quot;&quot;</span>\n    <span class=\"n\">config</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">render_markdown</span>\n
    \   <span class=\"n\">articles</span> <span class=\"o\">=</span> <span class=\"nb\">list</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">filter</span><span class=\"p\">(</span><span class=\"s2\">&quot;not
    skip&quot;</span><span class=\"p\">))</span>\n\n    <span class=\"k\">with</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">cache</span>
    <span class=\"k\">as</span> <span class=\"n\">cache</span><span class=\"p\">:</span>\n
    \       <span class=\"k\">with</span> <span class=\"n\">concurrent</span><span
    class=\"o\">.</span><span class=\"n\">futures</span><span class=\"o\">.</span><span
    class=\"n\">ThreadPoolExecutor</span><span class=\"p\">()</span> <span class=\"k\">as</span>
    <span class=\"n\">executor</span><span class=\"p\">:</span>\n            <span
    class=\"n\">render_func</span> <span class=\"o\">=</span> <span class=\"n\">partial</span><span
    class=\"p\">(</span><span class=\"n\">render_article_parallel</span><span class=\"p\">,</span>
    <span class=\"n\">markata</span><span class=\"p\">,</span> <span class=\"n\">config</span><span
    class=\"p\">,</span> <span class=\"n\">cache</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">args_list</span> <span class=\"o\">=</span> <span
    class=\"p\">[(</span><span class=\"n\">article</span><span class=\"p\">,)</span>
    <span class=\"k\">for</span> <span class=\"n\">article</span> <span class=\"ow\">in</span>
    <span class=\"n\">articles</span><span class=\"p\">]</span>\n\n            <span
    class=\"k\">for</span> <span class=\"n\">article</span><span class=\"p\">,</span>
    <span class=\"n\">html</span> <span class=\"ow\">in</span> <span class=\"n\">executor</span><span
    class=\"o\">.</span><span class=\"n\">map</span><span class=\"p\">(</span><span
    class=\"n\">render_func</span><span class=\"p\">,</span> <span class=\"n\">args_list</span><span
    class=\"p\">):</span>\n                <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">html</span> <span class=\"o\">=</span> <span class=\"n\">html</span>\n
    \               <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">article_html</span> <span class=\"o\">=</span> <span class=\"n\">copy</span><span
    class=\"o\">.</span><span class=\"n\">deepcopy</span><span class=\"p\">(</span><span
    class=\"n\">html</span><span class=\"p\">)</span>\n\n    <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">rendered_posts</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">posts</span>\n</pre></div>\n\n</pre>\n\n</div>\n\n
    \   </section>\n</article>    </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>render_markdown.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"The   plugin converts markdown content
    to HTML. This plugin is essential for rendering markdown files loaded by the   plugin.
    Installation This plugin is\u2026\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n
    \       <meta property=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    <head>\n        <script src=\"https://unpkg.com/@tailwindcss/browser@4\"></script>\n<title>render_markdown.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"The   plugin converts markdown content
    to HTML. This plugin is essential for rendering markdown files loaded by the   plugin.
    Installation This plugin is\u2026\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
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
    class=\"title\">\n    <h1 id=\"title\">\n        render_markdown.py\n    </h1>\n</section></article>\n
    \    </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       render_markdown.py\n    </h1>\n</section>    <section class=\"body\">\n
    \       <hr />\n<p>The <code>markata.plugins.render_markdown</code> plugin converts
    markdown content to HTML.\nThis plugin is essential for rendering markdown files
    loaded by the <code>load</code> plugin.</p>\n<h2 id=\"installation\">Installation
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>This plugin is built-in
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
    class=\"p\">[</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;markata.plugins.render_markdown&quot;</span><span
    class=\"p\">,</span>\n<span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<h2
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Since this plugin is
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
    class=\"p\">[</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;markata.plugins.render_markdown&quot;</span><span
    class=\"p\">,</span>\n<span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>Note:
    Disabling this plugin will prevent markdown files from being rendered to HTML.</p>\n<h2
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<h2 id=\"markdown-backend-selection\">Markdown
    Backend Selection <a class=\"header-anchor\" href=\"#markdown-backend-selection\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Choose from 3 supported
    markdown backends by setting <code>markdown_backend</code> in your <code>markata.toml</code>:</p>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"c1\">## choose
    your markdown backend</span>\n<span class=\"c1\"># markdown_backend=&#39;markdown&#39;
    \     # Python-Markdown</span>\n<span class=\"c1\"># markdown_backend=&#39;markdown2&#39;
    \    # markdown2</span>\n<span class=\"n\">markdown_backend</span><span class=\"o\">=</span><span
    class=\"s1\">&#39;markdown-it-py&#39;</span><span class=\"w\">  </span><span class=\"c1\">#
    markdown-it-py (default)</span>\n</pre></div>\n\n</pre>\n\n<h2 id=\"backend-specific-configuration\">Backend-Specific
    Configuration <a class=\"header-anchor\" href=\"#backend-specific-configuration\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<h3>markdown-it-py</h3>\n<p>Configure
    markdown-it-py behavior in your <code>markata.toml</code>:</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata.markdown_it_py]</span>\n<span
    class=\"c1\"># Set the flavor - options: &#39;zero&#39;, &#39;commonmark&#39;,
    &#39;gfm-like&#39;</span>\n<span class=\"n\">config</span><span class=\"w\"> </span><span
    class=\"o\">=</span><span class=\"w\"> </span><span class=\"s1\">&#39;gfm-like&#39;</span>\n\n<span
    class=\"c1\"># Enable specific plugins</span>\n<span class=\"n\">enable</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"p\">[</span>\n<span class=\"w\">    </span><span class=\"s1\">&#39;table&#39;</span><span
    class=\"p\">,</span>\n<span class=\"w\">    </span><span class=\"s1\">&#39;strikethrough&#39;</span><span
    class=\"p\">,</span>\n<span class=\"w\">    </span><span class=\"s1\">&#39;footnote&#39;</span><span
    class=\"p\">,</span>\n<span class=\"p\">]</span>\n\n<span class=\"c1\"># Disable
    specific plugins</span>\n<span class=\"n\">disable</span><span class=\"w\"> </span><span
    class=\"o\">=</span><span class=\"w\"> </span><span class=\"p\">[</span>\n<span
    class=\"w\">    </span><span class=\"s1\">&#39;linkify&#39;</span><span class=\"p\">,</span>\n<span
    class=\"p\">]</span>\n\n<span class=\"c1\"># Configure plugins</span>\n<span class=\"k\">[markata.markdown_it_py.plugins.footnote]</span>\n<span
    class=\"c1\"># Plugin-specific settings here</span>\n</pre></div>\n\n</pre>\n\n<p>Read
    more about markdown-it-py settings in their <a href=\"https://markdown-it-py.readthedocs.io/en/latest/\">documentation</a>.</p>\n<h2
    id=\"cache-configuration\">Cache Configuration <a class=\"header-anchor\" href=\"#cache-configuration\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Control markdown rendering
    cache:</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button class='copy'
    title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata.render_markdown]</span>\n<span
    class=\"n\">cache_expire</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"mi\">3600</span><span class=\"w\">  </span><span
    class=\"c1\"># Cache expiration in seconds</span>\n</pre></div>\n\n</pre>\n\n<h2
    id=\"functionality\">Functionality <a class=\"header-anchor\" href=\"#functionality\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<h2 id=\"registered-attributes\">Registered
    Attributes <a class=\"header-anchor\" href=\"#registered-attributes\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>The plugin registers
    the following attributes on Post objects:</p>\n<ul>\n<li><code>html</code>: The
    rendered HTML content from the markdown source</li>\n</ul>\n<h2 id=\"dependencies\">Dependencies
    <a class=\"header-anchor\" href=\"#dependencies\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>This plugin depends
    on:</p>\n<ul>\n<li>One of: python-markdown, markdown2, or markdown-it-py (based
    on configuration)</li>\n<li>The <code>load</code> plugin to provide markdown content
    for rendering</li>\n</ul>\n<hr />\n<div class=\"admonition function\">\n<p class=\"admonition-title\">Function</p>\n<h2
    id=\"configure\" class=\"admonition-title\" style=\"margin: 0; padding: .5rem
    1rem;\">configure <em class=\"small\">function</em></h2>\n<p>Sets up a markdown
    instance as md</p>\n</div>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">configure <em class='small'>source</em></p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">def</span><span
    class=\"w\"> </span><span class=\"nf\">configure</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n    <span class=\"s2\">&quot;Sets up a markdown instance
    as md&quot;</span>\n    <span class=\"c1\"># if &quot;markdown_extensions&quot;
    not in markata.config:</span>\n    <span class=\"c1\">#     markdown_extensions
    = [&quot;&quot;]</span>\n    <span class=\"c1\"># if isinstance(markata.config[&quot;markdown_extensions&quot;],
    str):</span>\n    <span class=\"c1\">#     markdown_extensions = [markata.config[&quot;markdown_extensions&quot;]]</span>\n
    \   <span class=\"c1\"># if isinstance(markata.config[&quot;markdown_extensions&quot;],
    list):</span>\n    <span class=\"c1\">#     markdown_extensions = markata.config[&quot;markdown_extensions&quot;]</span>\n
    \   <span class=\"c1\"># else:</span>\n    <span class=\"c1\">#     raise TypeError(&quot;markdown_extensions
    should be List[str]&quot;)</span>\n\n    <span class=\"c1\"># markata.markdown_extensions
    = [*DEFAULT_MD_EXTENSIONS, *markdown_extensions]</span>\n\n    <span class=\"k\">if</span>
    <span class=\"p\">(</span>\n        <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;markdown_backend&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;&quot;</span><span class=\"p\">)</span>\n
    \       <span class=\"o\">.</span><span class=\"n\">lower</span><span class=\"p\">()</span>\n
    \       <span class=\"o\">.</span><span class=\"n\">replace</span><span class=\"p\">(</span><span
    class=\"s2\">&quot; &quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;-&quot;</span><span
    class=\"p\">)</span>\n        <span class=\"o\">.</span><span class=\"n\">replace</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;_&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;-&quot;</span><span class=\"p\">)</span>\n        <span
    class=\"o\">==</span> <span class=\"s2\">&quot;markdown-it-py&quot;</span>\n    <span
    class=\"p\">):</span>\n        <span class=\"kn\">from</span><span class=\"w\">
    </span><span class=\"nn\">markdown_it</span><span class=\"w\"> </span><span class=\"kn\">import</span>
    <span class=\"n\">MarkdownIt</span>\n\n        <span class=\"n\">config_update</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;markdown_it_py&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">{})</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span>\n            <span class=\"s2\">&quot;options_update&quot;</span><span
    class=\"p\">,</span>\n            <span class=\"p\">{</span>\n                <span
    class=\"s2\">&quot;linkify&quot;</span><span class=\"p\">:</span> <span class=\"kc\">True</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;html&quot;</span><span
    class=\"p\">:</span> <span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \               <span class=\"s2\">&quot;typographer&quot;</span><span class=\"p\">:</span>
    <span class=\"kc\">True</span><span class=\"p\">,</span>\n                <span
    class=\"s2\">&quot;highlight&quot;</span><span class=\"p\">:</span> <span class=\"n\">highlight_code</span><span
    class=\"p\">,</span>\n            <span class=\"p\">},</span>\n        <span class=\"p\">)</span>\n
    \       <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">config_update</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;highlight&quot;</span><span
    class=\"p\">),</span> <span class=\"nb\">str</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">module</span> <span class=\"o\">=</span> <span class=\"n\">config_update</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;highlight&quot;</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">split</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;:&quot;</span><span class=\"p\">)[</span><span class=\"mi\">0</span><span
    class=\"p\">]</span>\n            <span class=\"n\">func</span> <span class=\"o\">=</span>
    <span class=\"n\">config_update</span><span class=\"p\">[</span><span class=\"s2\">&quot;highlight&quot;</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">split</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;:&quot;</span><span class=\"p\">)[</span><span
    class=\"mi\">1</span><span class=\"p\">]</span>\n            <span class=\"n\">config_update</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;highlight&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"nb\">getattr</span><span class=\"p\">(</span>\n
    \               <span class=\"n\">importlib</span><span class=\"o\">.</span><span
    class=\"n\">import_module</span><span class=\"p\">(</span><span class=\"n\">module</span><span
    class=\"p\">),</span>\n                <span class=\"n\">func</span><span class=\"p\">,</span>\n
    \           <span class=\"p\">)</span>\n\n        <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">md</span> <span class=\"o\">=</span> <span
    class=\"n\">MarkdownIt</span><span class=\"p\">(</span>\n            <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;markdown_it_py&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">{})</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;config&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;gfm-like&quot;</span><span class=\"p\">),</span>\n
    \           <span class=\"n\">config_update</span><span class=\"p\">,</span>\n
    \       <span class=\"p\">)</span>\n        <span class=\"k\">for</span> <span
    class=\"n\">plugin</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;markdown_it_py&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">{})</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;enable&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">[]):</span>\n            <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">md</span><span class=\"o\">.</span><span
    class=\"n\">enable</span><span class=\"p\">(</span><span class=\"n\">plugin</span><span
    class=\"p\">)</span>\n        <span class=\"k\">for</span> <span class=\"n\">plugin</span>
    <span class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;markdown_it_py&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">{})</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;disable&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">[]):</span>\n            <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">md</span><span class=\"o\">.</span><span
    class=\"n\">disable</span><span class=\"p\">(</span><span class=\"n\">plugin</span><span
    class=\"p\">)</span>\n\n        <span class=\"n\">plugins</span> <span class=\"o\">=</span>
    <span class=\"n\">copy</span><span class=\"o\">.</span><span class=\"n\">deepcopy</span><span
    class=\"p\">(</span>\n            <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;markdown_it_py&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">{})</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;plugins&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">[]),</span>\n        <span class=\"p\">)</span>\n
    \       <span class=\"k\">for</span> <span class=\"n\">plugin</span> <span class=\"ow\">in</span>
    <span class=\"n\">plugins</span><span class=\"p\">:</span>\n            <span
    class=\"k\">if</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">plugin</span><span class=\"p\">[</span><span class=\"s2\">&quot;plugin&quot;</span><span
    class=\"p\">],</span> <span class=\"nb\">str</span><span class=\"p\">):</span>\n
    \               <span class=\"n\">plugin</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;plugin_str&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"n\">plugin</span><span class=\"p\">[</span><span class=\"s2\">&quot;plugin&quot;</span><span
    class=\"p\">]</span>\n                <span class=\"n\">plugin_module</span> <span
    class=\"o\">=</span> <span class=\"n\">plugin</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;plugin&quot;</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">split</span><span class=\"p\">(</span><span class=\"s2\">&quot;:&quot;</span><span
    class=\"p\">)[</span><span class=\"mi\">0</span><span class=\"p\">]</span>\n                <span
    class=\"n\">plugin_func</span> <span class=\"o\">=</span> <span class=\"n\">plugin</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;plugin&quot;</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">split</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;:&quot;</span><span class=\"p\">)[</span><span class=\"mi\">1</span><span
    class=\"p\">]</span>\n                <span class=\"n\">plugin</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;plugin&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"nb\">getattr</span><span class=\"p\">(</span>\n                    <span
    class=\"n\">importlib</span><span class=\"o\">.</span><span class=\"n\">import_module</span><span
    class=\"p\">(</span><span class=\"n\">plugin_module</span><span class=\"p\">),</span>\n
    \                   <span class=\"n\">plugin_func</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">)</span>\n            <span class=\"n\">plugin</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;config&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"n\">plugin</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;config&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">{})</span>\n            <span class=\"k\">for</span>
    <span class=\"n\">k</span><span class=\"p\">,</span> <span class=\"n\">_v</span>
    <span class=\"ow\">in</span> <span class=\"n\">plugin</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;config&quot;</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">items</span><span class=\"p\">():</span>\n                <span class=\"k\">if</span>
    <span class=\"n\">k</span> <span class=\"o\">==</span> <span class=\"s2\">&quot;markata&quot;</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">plugin</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;config&quot;</span><span class=\"p\">][</span><span
    class=\"n\">k</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"n\">markata</span>\n\n            <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">md</span> <span class=\"o\">=</span> <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">md</span><span
    class=\"o\">.</span><span class=\"n\">use</span><span class=\"p\">(</span><span
    class=\"n\">plugin</span><span class=\"p\">[</span><span class=\"s2\">&quot;plugin&quot;</span><span
    class=\"p\">],</span> <span class=\"o\">**</span><span class=\"n\">plugin</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;config&quot;</span><span class=\"p\">])</span>\n\n
    \       <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">md</span><span
    class=\"o\">.</span><span class=\"n\">convert</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">md</span><span
    class=\"o\">.</span><span class=\"n\">render</span>\n        <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">md</span><span class=\"o\">.</span><span
    class=\"n\">toc</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;&quot;</span>\n
    \   <span class=\"k\">elif</span> <span class=\"p\">(</span>\n        <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;markdown_backend&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;&quot;</span><span class=\"p\">)</span>\n
    \       <span class=\"o\">.</span><span class=\"n\">lower</span><span class=\"p\">()</span>\n
    \       <span class=\"o\">.</span><span class=\"n\">replace</span><span class=\"p\">(</span><span
    class=\"s2\">&quot; &quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;-&quot;</span><span
    class=\"p\">)</span>\n        <span class=\"o\">.</span><span class=\"n\">replace</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;_&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;-&quot;</span><span class=\"p\">)</span>\n        <span
    class=\"o\">==</span> <span class=\"s2\">&quot;markdown2&quot;</span>\n    <span
    class=\"p\">):</span>\n        <span class=\"kn\">import</span><span class=\"w\">
    </span><span class=\"nn\">markdown2</span>\n\n        <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">md</span> <span class=\"o\">=</span> <span
    class=\"n\">markdown2</span><span class=\"o\">.</span><span class=\"n\">Markdown</span><span
    class=\"p\">(</span>\n            <span class=\"n\">extras</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">render_markdown</span><span class=\"o\">.</span><span
    class=\"n\">extensions</span>\n        <span class=\"p\">)</span>\n        <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">md</span><span
    class=\"o\">.</span><span class=\"n\">toc</span> <span class=\"o\">=</span> <span
    class=\"s2\">&quot;&quot;</span>\n    <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \       <span class=\"kn\">import</span><span class=\"w\"> </span><span class=\"nn\">markdown</span>\n\n
    \       <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">md</span>
    <span class=\"o\">=</span> <span class=\"n\">markdown</span><span class=\"o\">.</span><span
    class=\"n\">Markdown</span><span class=\"p\">(</span>\n            <span class=\"n\">extensions</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">render_markdown</span><span
    class=\"o\">.</span><span class=\"n\">extensions</span>\n        <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition function\">\n<p class=\"admonition-title\">Function</p>\n<h2
    id=\"render\" class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">render
    <em class=\"small\">function</em></h2>\n<p>Render markdown content in parallel.</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">render
    <em class='small'>source</em></p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">def</span><span
    class=\"w\"> </span><span class=\"nf\">render</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Render
    markdown content in parallel.&quot;&quot;&quot;</span>\n    <span class=\"n\">config</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">render_markdown</span>\n
    \   <span class=\"n\">articles</span> <span class=\"o\">=</span> <span class=\"nb\">list</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">filter</span><span class=\"p\">(</span><span class=\"s2\">&quot;not
    skip&quot;</span><span class=\"p\">))</span>\n\n    <span class=\"k\">with</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">cache</span>
    <span class=\"k\">as</span> <span class=\"n\">cache</span><span class=\"p\">:</span>\n
    \       <span class=\"k\">with</span> <span class=\"n\">concurrent</span><span
    class=\"o\">.</span><span class=\"n\">futures</span><span class=\"o\">.</span><span
    class=\"n\">ThreadPoolExecutor</span><span class=\"p\">()</span> <span class=\"k\">as</span>
    <span class=\"n\">executor</span><span class=\"p\">:</span>\n            <span
    class=\"n\">render_func</span> <span class=\"o\">=</span> <span class=\"n\">partial</span><span
    class=\"p\">(</span><span class=\"n\">render_article_parallel</span><span class=\"p\">,</span>
    <span class=\"n\">markata</span><span class=\"p\">,</span> <span class=\"n\">config</span><span
    class=\"p\">,</span> <span class=\"n\">cache</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">args_list</span> <span class=\"o\">=</span> <span
    class=\"p\">[(</span><span class=\"n\">article</span><span class=\"p\">,)</span>
    <span class=\"k\">for</span> <span class=\"n\">article</span> <span class=\"ow\">in</span>
    <span class=\"n\">articles</span><span class=\"p\">]</span>\n\n            <span
    class=\"k\">for</span> <span class=\"n\">article</span><span class=\"p\">,</span>
    <span class=\"n\">html</span> <span class=\"ow\">in</span> <span class=\"n\">executor</span><span
    class=\"o\">.</span><span class=\"n\">map</span><span class=\"p\">(</span><span
    class=\"n\">render_func</span><span class=\"p\">,</span> <span class=\"n\">args_list</span><span
    class=\"p\">):</span>\n                <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">html</span> <span class=\"o\">=</span> <span class=\"n\">html</span>\n
    \               <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">article_html</span> <span class=\"o\">=</span> <span class=\"n\">copy</span><span
    class=\"o\">.</span><span class=\"n\">deepcopy</span><span class=\"p\">(</span><span
    class=\"n\">html</span><span class=\"p\">)</span>\n\n    <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">rendered_posts</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">posts</span>\n</pre></div>\n\n</pre>\n\n</div>\n\n
    \   </section>\n</article>"
  raw.md: ''
published: false
slug: markata/plugins/render-markdown
title: render_markdown.py


---

---

The `markata.plugins.render_markdown` plugin converts markdown content to HTML.
This plugin is essential for rendering markdown files loaded by the `load` plugin.

## Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.render_markdown",
]
```

## Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.render_markdown",
]
```

Note: Disabling this plugin will prevent markdown files from being rendered to HTML.

## Configuration

## Markdown Backend Selection

Choose from 3 supported markdown backends by setting `markdown_backend` in your `markata.toml`:

```toml
## choose your markdown backend
# markdown_backend='markdown'      # Python-Markdown
# markdown_backend='markdown2'     # markdown2
markdown_backend='markdown-it-py'  # markdown-it-py (default)
```

## Backend-Specific Configuration

### markdown-it-py

Configure markdown-it-py behavior in your `markata.toml`:

```toml
[markata.markdown_it_py]
# Set the flavor - options: 'zero', 'commonmark', 'gfm-like'
config = 'gfm-like'

# Enable specific plugins
enable = [
    'table',
    'strikethrough',
    'footnote',
]

# Disable specific plugins
disable = [
    'linkify',
]

# Configure plugins
[markata.markdown_it_py.plugins.footnote]
# Plugin-specific settings here
```

Read more about markdown-it-py settings in their [documentation](https://markdown-it-py.readthedocs.io/en/latest/).

## Cache Configuration

Control markdown rendering cache:

```toml
[markata.render_markdown]
cache_expire = 3600  # Cache expiration in seconds
```

## Functionality

## Registered Attributes

The plugin registers the following attributes on Post objects:
- `html`: The rendered HTML content from the markdown source

## Dependencies

This plugin depends on:
- One of: python-markdown, markdown2, or markdown-it-py (based on configuration)
- The `load` plugin to provide markdown content for rendering

---

!!! function
    <h2 id="configure" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">configure <em class="small">function</em></h2>

    Sets up a markdown instance as md

???+ source "configure <em class='small'>source</em>"
    ```python
    def configure(markata: "Markata") -> None:
        "Sets up a markdown instance as md"
        # if "markdown_extensions" not in markata.config:
        #     markdown_extensions = [""]
        # if isinstance(markata.config["markdown_extensions"], str):
        #     markdown_extensions = [markata.config["markdown_extensions"]]
        # if isinstance(markata.config["markdown_extensions"], list):
        #     markdown_extensions = markata.config["markdown_extensions"]
        # else:
        #     raise TypeError("markdown_extensions should be List[str]")

        # markata.markdown_extensions = [*DEFAULT_MD_EXTENSIONS, *markdown_extensions]

        if (
            markata.config.get("markdown_backend", "")
            .lower()
            .replace(" ", "-")
            .replace("_", "-")
            == "markdown-it-py"
        ):
            from markdown_it import MarkdownIt

            config_update = markata.config.get("markdown_it_py", {}).get(
                "options_update",
                {
                    "linkify": True,
                    "html": True,
                    "typographer": True,
                    "highlight": highlight_code,
                },
            )
            if isinstance(config_update.get("highlight"), str):
                module = config_update["highlight"].split(":")[0]
                func = config_update["highlight"].split(":")[1]
                config_update["highlight"] = getattr(
                    importlib.import_module(module),
                    func,
                )

            markata.md = MarkdownIt(
                markata.config.get("markdown_it_py", {}).get("config", "gfm-like"),
                config_update,
            )
            for plugin in markata.config.get("markdown_it_py", {}).get("enable", []):
                markata.md.enable(plugin)
            for plugin in markata.config.get("markdown_it_py", {}).get("disable", []):
                markata.md.disable(plugin)

            plugins = copy.deepcopy(
                markata.config.get("markdown_it_py", {}).get("plugins", []),
            )
            for plugin in plugins:
                if isinstance(plugin["plugin"], str):
                    plugin["plugin_str"] = plugin["plugin"]
                    plugin_module = plugin["plugin"].split(":")[0]
                    plugin_func = plugin["plugin"].split(":")[1]
                    plugin["plugin"] = getattr(
                        importlib.import_module(plugin_module),
                        plugin_func,
                    )
                plugin["config"] = plugin.get("config", {})
                for k, _v in plugin["config"].items():
                    if k == "markata":
                        plugin["config"][k] = markata

                markata.md = markata.md.use(plugin["plugin"], **plugin["config"])

            markata.md.convert = markata.md.render
            markata.md.toc = ""
        elif (
            markata.config.get("markdown_backend", "")
            .lower()
            .replace(" ", "-")
            .replace("_", "-")
            == "markdown2"
        ):
            import markdown2

            markata.md = markdown2.Markdown(
                extras=markata.config.render_markdown.extensions
            )
            markata.md.toc = ""
        else:
            import markdown

            markata.md = markdown.Markdown(
                extensions=markata.config.render_markdown.extensions
            )
    ```
!!! function
    <h2 id="render" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">render <em class="small">function</em></h2>

    Render markdown content in parallel.

???+ source "render <em class='small'>source</em>"
    ```python
    def render(markata: "Markata") -> None:
        """Render markdown content in parallel."""
        config = markata.config.render_markdown
        articles = list(markata.filter("not skip"))

        with markata.cache as cache:
            with concurrent.futures.ThreadPoolExecutor() as executor:
                render_func = partial(render_article_parallel, markata, config, cache)
                args_list = [(article,) for article in articles]

                for article, html in executor.map(render_func, args_list):
                    article.html = html
                    article.article_html = copy.deepcopy(html)

        markata.rendered_posts = markata.posts
    ```