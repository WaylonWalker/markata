---
content: "---\n\nThe `markata.plugins.auto_description` plugin automatically generates
  descriptions for your\nposts by extracting text from the first paragraphs of your
  markdown content. It can create\nmultiple descriptions of different lengths for
  different purposes (e.g., SEO, previews).\n\n# Installation\n\nThis plugin is built-in
  and enabled by default through the 'default' plugin.\nIf you want to be explicit,
  you can add it to your list of plugins:\n\n```toml\nhooks = [\n    \"markata.plugins.auto_description\",\n]\n```\n\n#
  Uninstallation\n\nSince this plugin is included in the default plugin set, to disable
  it you must explicitly\nadd it to the disabled_hooks list if you are using the 'default'
  plugin:\n\n```toml\ndisabled_hooks = [\n    \"markata.plugins.auto_description\",\n]\n```\n\n#
  Configuration\n\nConfigure multiple description types in your `markata.toml`:\n\n```toml\n#
  Standard description (160 characters)\n[markata.auto_description.description]\nlen
  = 160\n\n# Longer description (250 characters)\n[markata.auto_description.long_description]\nlen
  = 250\n\n# Full description (500 characters)\n[markata.auto_description.super_description]\nlen
  = 500\n```\n\nEach description configuration:\n- Must be under `markata.auto_description.[name]`\n-
  Requires a `len` parameter specifying maximum character length\n- Creates an attribute
  named after the configuration key\n\n## Frontmatter Override\n\nYou can override
  automatic descriptions by setting them in frontmatter:\n\n```markdown\n---\ntitle:
  My Post\ndescription: My custom description\nlong_description: A longer custom description\n---\n```\n\n#
  Functionality\n\n## Description Generation\n\nThe plugin:\n1. Converts markdown
  to plain text\n2. Finds the first meaningful paragraphs\n3. Truncates to the specified
  length\n4. Ensures clean word breaks\n5. Caches results for performance\n\n## Registered
  Attributes\n\nFor each configured description (e.g., `description`, `long_description`),
  the plugin adds:\n- The description attribute with the generated/specified text\n-
  Truncated to the configured length\n- Preserving complete words\n\n## Dependencies\n\nThis
  plugin depends on:\n- markdown-it-py for markdown parsing\n\n---\n\n\n\n!!! function\n
  \   <h2 id=\"get_description\" class=\"admonition-title\" style=\"margin: 0; padding:
  .5rem 1rem;\">get_description <em class=\"small\">function</em></h2>\n\n    Get
  the full-length description for a single post by converting markdown to plain text.\n
  \   Uses markdown-it-py to parse the markdown and extracts text content from all
  nodes.\n\n???+ source \"get_description <em class='small'>source</em>\"\n    ```python\n
  \   def get_description(article: \"Post\") -> str:\n        \"\"\"\n        Get
  the full-length description for a single post by converting markdown to plain text.\n
  \       Uses markdown-it-py to parse the markdown and extracts text content from
  all nodes.\n        \"\"\"\n        from markdown_it import MarkdownIt\n\n        md
  = MarkdownIt(\"commonmark\")\n        tokens = md.parse(article.content)\n\n        #
  Extract text content from all inline tokens\n        text_chunks = []\n        for
  token in tokens:\n            if token.type == \"inline\" and token.content:\n                text_chunks.append(token.content)\n\n
  \       description = \" \".join(text_chunks)\n        return html.escape(description)\n
  \   ```\n\n\n\n!!! function\n    <h2 id=\"set_description\" class=\"admonition-title\"
  style=\"margin: 0; padding: .5rem 1rem;\">set_description <em class=\"small\">function</em></h2>\n\n
  \   For a given `article`, find the description, put it in the cache, and set\n
  \   the configured descriptions for the article.\n\n???+ source \"set_description
  <em class='small'>source</em>\"\n    ```python\n    def set_description(\n        markata:
  \"Markata\",\n        article: \"Post\",\n        cache: \"FanoutCache\",\n        config:
  Dict,\n        max_description: int = 500,\n        plugin_text: None = \"\",\n
  \   ) -> None:\n        \"\"\"\n        For a given `article`, find the description,
  put it in the cache, and set\n        the configured descriptions for the article.\n
  \       \"\"\"\n        key = markata.make_hash(\n            \"auto_description\",\n
  \           article.content,\n            plugin_text,\n            config,\n        )\n\n
  \       description_from_cache = markata.precache.get(key)\n        if description_from_cache
  is None:\n            description = get_description(article)[:max_description]\n
  \           markata.cache.set(key, description, expire=markata.config.default_cache_expire)\n
  \       else:\n            description = description_from_cache\n\n        for description_key
  in config:\n            if description_key not in [\"cache_expire\", \"config_key\"]:\n
  \               # overwrites missing (None) and empty ('')\n                if not
  article.metadata.get(description_key):\n                    article.metadata[description_key]
  = description[\n                        : config[description_key][\"len\"]\n                    ]\n
  \   ```\n\n\n\n!!! function\n    <h2 id=\"pre_render\" class=\"admonition-title\"
  style=\"margin: 0; padding: .5rem 1rem;\">pre_render <em class=\"small\">function</em></h2>\n\n
  \   The Markata hook that will set descriptions for all posts in the pre-render
  phase.\n\n???+ source \"pre_render <em class='small'>source</em>\"\n    ```python\n
  \   def pre_render(markata: \"Markata\") -> None:\n        \"\"\"\n        The Markata
  hook that will set descriptions for all posts in the pre-render phase.\n        \"\"\"\n
  \       config = markata.get_plugin_config(__file__)\n\n        if \"description\"
  not in config.keys():\n            config[\"description\"] = {}\n            config[\"description\"][\"len\"]
  = 160\n\n        if \"long_description\" not in config.keys():\n            config[\"long_description\"]
  = {}\n            config[\"long_description\"][\"len\"] = 250\n\n        def try_config_get(key:
  str) -> Any:\n            try:\n                return config.get(key).get(\"len\")
  or None\n            except AttributeError:\n                return None\n\n        max_description
  = max(\n            [\n                value\n                for description_key
  in config\n                if (value := try_config_get(description_key))\n            ],\n
  \       )\n\n        with markata.cache as cache:\n            for article in markata.iter_articles(\"setting
  auto description\"):\n                set_description(\n                    markata=markata,\n
  \                   article=article,\n                    cache=cache,\n                    config=config,\n
  \                   max_description=max_description,\n                    plugin_text=Path(__file__).read_text(),\n
  \               )\n    ```"
date: 2025-02-14
description: 'The `markata.plugins.auto_description` plugin automatically generates
  descriptions for your

  posts by extracting text from the first paragraphs of your markdown '
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>auto_description.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"The `markata.plugins.auto_description`
    plugin automatically generates descriptions for your\nposts by extracting text
    from the first paragraphs of your markdown \" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n
    \       <meta property=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    <head>\n        <script src=\"https://unpkg.com/@tailwindcss/browser@4\"></script>\n<title>auto_description.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"The `markata.plugins.auto_description`
    plugin automatically generates descriptions for your\nposts by extracting text
    from the first paragraphs of your markdown \" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        auto_description.py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <hr />\n<p>The <code>markata.plugins.auto_description</code>
    plugin automatically generates descriptions for your\nposts by extracting text
    from the first paragraphs of your markdown content. It can create\nmultiple descriptions
    of different lengths for different purposes (e.g., SEO, previews).</p>\n<h1 id=\"installation\">Installation
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
    class=\"p\">[</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;markata.plugins.auto_description&quot;</span><span
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
    class=\"p\">[</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;markata.plugins.auto_description&quot;</span><span
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>Configure multiple description
    types in your <code>markata.toml</code>:</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"c1\"># Standard
    description (160 characters)</span>\n<span class=\"k\">[markata.auto_description.description]</span>\n<span
    class=\"n\">len</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"mi\">160</span>\n\n<span class=\"c1\"># Longer
    description (250 characters)</span>\n<span class=\"k\">[markata.auto_description.long_description]</span>\n<span
    class=\"n\">len</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"mi\">250</span>\n\n<span class=\"c1\"># Full
    description (500 characters)</span>\n<span class=\"k\">[markata.auto_description.super_description]</span>\n<span
    class=\"n\">len</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"mi\">500</span>\n</pre></div>\n\n</pre>\n\n<p>Each
    description configuration:</p>\n<ul>\n<li>Must be under <code>markata.auto_description.[name]</code></li>\n<li>Requires
    a <code>len</code> parameter specifying maximum character length</li>\n<li>Creates
    an attribute named after the configuration key</li>\n</ul>\n<h2 id=\"frontmatter-override\">Frontmatter
    Override <a class=\"header-anchor\" href=\"#frontmatter-override\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>You can override automatic
    descriptions by setting them in frontmatter:</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span>---\ntitle: My Post\ndescription:
    My custom description\n<span class=\"gu\">long_description: A longer custom description</span>\n<span
    class=\"gu\">---</span>\n</pre></div>\n\n</pre>\n\n<h1 id=\"functionality\">Functionality
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<h2 id=\"description-generation\">Description
    Generation <a class=\"header-anchor\" href=\"#description-generation\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>The plugin:</p>\n<ol>\n<li>Converts
    markdown to plain text</li>\n<li>Finds the first meaningful paragraphs</li>\n<li>Truncates
    to the specified length</li>\n<li>Ensures clean word breaks</li>\n<li>Caches results
    for performance</li>\n</ol>\n<h2 id=\"registered-attributes\">Registered Attributes
    <a class=\"header-anchor\" href=\"#registered-attributes\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>For each configured
    description (e.g., <code>description</code>, <code>long_description</code>), the
    plugin adds:</p>\n<ul>\n<li>The description attribute with the generated/specified
    text</li>\n<li>Truncated to the configured length</li>\n<li>Preserving complete
    words</li>\n</ul>\n<h2 id=\"dependencies\">Dependencies <a class=\"header-anchor\"
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
    on:</p>\n<ul>\n<li>markdown-it-py for markdown parsing</li>\n</ul>\n<hr />\n<div
    class=\"admonition function\">\n<p class=\"admonition-title\">Function</p>\n<h2
    id=\"get_description\" class=\"admonition-title\" style=\"margin: 0; padding:
    .5rem 1rem;\">get_description <em class=\"small\">function</em></h2>\n<p>Get the
    full-length description for a single post by converting markdown to plain text.\nUses
    markdown-it-py to parse the markdown and extracts text content from all nodes.</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">get_description
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
    class=\"w\"> </span><span class=\"nf\">get_description</span><span class=\"p\">(</span><span
    class=\"n\">article</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span
    class=\"p\">:</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">    Get the full-length description for a single post by converting
    markdown to plain text.</span>\n<span class=\"sd\">    Uses markdown-it-py to
    parse the markdown and extracts text content from all nodes.</span>\n<span class=\"sd\">
    \   &quot;&quot;&quot;</span>\n    <span class=\"kn\">from</span><span class=\"w\">
    </span><span class=\"nn\">markdown_it</span><span class=\"w\"> </span><span class=\"kn\">import</span>
    <span class=\"n\">MarkdownIt</span>\n\n    <span class=\"n\">md</span> <span class=\"o\">=</span>
    <span class=\"n\">MarkdownIt</span><span class=\"p\">(</span><span class=\"s2\">&quot;commonmark&quot;</span><span
    class=\"p\">)</span>\n    <span class=\"n\">tokens</span> <span class=\"o\">=</span>
    <span class=\"n\">md</span><span class=\"o\">.</span><span class=\"n\">parse</span><span
    class=\"p\">(</span><span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">content</span><span class=\"p\">)</span>\n\n    <span class=\"c1\">#
    Extract text content from all inline tokens</span>\n    <span class=\"n\">text_chunks</span>
    <span class=\"o\">=</span> <span class=\"p\">[]</span>\n    <span class=\"k\">for</span>
    <span class=\"n\">token</span> <span class=\"ow\">in</span> <span class=\"n\">tokens</span><span
    class=\"p\">:</span>\n        <span class=\"k\">if</span> <span class=\"n\">token</span><span
    class=\"o\">.</span><span class=\"n\">type</span> <span class=\"o\">==</span>
    <span class=\"s2\">&quot;inline&quot;</span> <span class=\"ow\">and</span> <span
    class=\"n\">token</span><span class=\"o\">.</span><span class=\"n\">content</span><span
    class=\"p\">:</span>\n            <span class=\"n\">text_chunks</span><span class=\"o\">.</span><span
    class=\"n\">append</span><span class=\"p\">(</span><span class=\"n\">token</span><span
    class=\"o\">.</span><span class=\"n\">content</span><span class=\"p\">)</span>\n\n
    \   <span class=\"n\">description</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;
    &quot;</span><span class=\"o\">.</span><span class=\"n\">join</span><span class=\"p\">(</span><span
    class=\"n\">text_chunks</span><span class=\"p\">)</span>\n    <span class=\"k\">return</span>
    <span class=\"n\">html</span><span class=\"o\">.</span><span class=\"n\">escape</span><span
    class=\"p\">(</span><span class=\"n\">description</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition function\">\n<p class=\"admonition-title\">Function</p>\n<h2
    id=\"set_description\" class=\"admonition-title\" style=\"margin: 0; padding:
    .5rem 1rem;\">set_description <em class=\"small\">function</em></h2>\n<p>For a
    given <code>article</code>, find the description, put it in the cache, and set\nthe
    configured descriptions for the article.</p>\n</div>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">set_description
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
    class=\"w\"> </span><span class=\"nf\">set_description</span><span class=\"p\">(</span>\n
    \   <span class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span>\n    <span class=\"n\">article</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">,</span>\n    <span
    class=\"n\">cache</span><span class=\"p\">:</span> <span class=\"s2\">&quot;FanoutCache&quot;</span><span
    class=\"p\">,</span>\n    <span class=\"n\">config</span><span class=\"p\">:</span>
    <span class=\"n\">Dict</span><span class=\"p\">,</span>\n    <span class=\"n\">max_description</span><span
    class=\"p\">:</span> <span class=\"nb\">int</span> <span class=\"o\">=</span>
    <span class=\"mi\">500</span><span class=\"p\">,</span>\n    <span class=\"n\">plugin_text</span><span
    class=\"p\">:</span> <span class=\"kc\">None</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;&quot;</span><span class=\"p\">,</span>\n<span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
    class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span class=\"sd\">
    \   For a given `article`, find the description, put it in the cache, and set</span>\n<span
    class=\"sd\">    the configured descriptions for the article.</span>\n<span class=\"sd\">
    \   &quot;&quot;&quot;</span>\n    <span class=\"n\">key</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">make_hash</span><span
    class=\"p\">(</span>\n        <span class=\"s2\">&quot;auto_description&quot;</span><span
    class=\"p\">,</span>\n        <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">content</span><span class=\"p\">,</span>\n        <span class=\"n\">plugin_text</span><span
    class=\"p\">,</span>\n        <span class=\"n\">config</span><span class=\"p\">,</span>\n
    \   <span class=\"p\">)</span>\n\n    <span class=\"n\">description_from_cache</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">precache</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"n\">key</span><span class=\"p\">)</span>\n    <span
    class=\"k\">if</span> <span class=\"n\">description_from_cache</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n        <span class=\"n\">description</span>
    <span class=\"o\">=</span> <span class=\"n\">get_description</span><span class=\"p\">(</span><span
    class=\"n\">article</span><span class=\"p\">)[:</span><span class=\"n\">max_description</span><span
    class=\"p\">]</span>\n        <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">cache</span><span class=\"o\">.</span><span class=\"n\">set</span><span
    class=\"p\">(</span><span class=\"n\">key</span><span class=\"p\">,</span> <span
    class=\"n\">description</span><span class=\"p\">,</span> <span class=\"n\">expire</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">default_cache_expire</span><span
    class=\"p\">)</span>\n    <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \       <span class=\"n\">description</span> <span class=\"o\">=</span> <span
    class=\"n\">description_from_cache</span>\n\n    <span class=\"k\">for</span>
    <span class=\"n\">description_key</span> <span class=\"ow\">in</span> <span class=\"n\">config</span><span
    class=\"p\">:</span>\n        <span class=\"k\">if</span> <span class=\"n\">description_key</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"p\">[</span><span
    class=\"s2\">&quot;cache_expire&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;config_key&quot;</span><span
    class=\"p\">]:</span>\n            <span class=\"c1\"># overwrites missing (None)
    and empty (&#39;&#39;)</span>\n            <span class=\"k\">if</span> <span class=\"ow\">not</span>
    <span class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">metadata</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"n\">description_key</span><span class=\"p\">):</span>\n                <span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">metadata</span><span
    class=\"p\">[</span><span class=\"n\">description_key</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"n\">description</span><span class=\"p\">[</span>\n
    \                   <span class=\"p\">:</span> <span class=\"n\">config</span><span
    class=\"p\">[</span><span class=\"n\">description_key</span><span class=\"p\">][</span><span
    class=\"s2\">&quot;len&quot;</span><span class=\"p\">]</span>\n                <span
    class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div class=\"admonition
    function\">\n<p class=\"admonition-title\">Function</p>\n<h2 id=\"pre_render\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">pre_render
    <em class=\"small\">function</em></h2>\n<p>The Markata hook that will set descriptions
    for all posts in the pre-render phase.</p>\n</div>\n<div class=\"admonition source
    is-collapsible collapsible-open\">\n<p class=\"admonition-title\">pre_render <em
    class='small'>source</em></p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    class=\"w\"> </span><span class=\"nf\">pre_render</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">    The Markata hook that will set descriptions for all posts in
    the pre-render phase.</span>\n<span class=\"sd\">    &quot;&quot;&quot;</span>\n
    \   <span class=\"n\">config</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">get_plugin_config</span><span class=\"p\">(</span><span
    class=\"vm\">__file__</span><span class=\"p\">)</span>\n\n    <span class=\"k\">if</span>
    <span class=\"s2\">&quot;description&quot;</span> <span class=\"ow\">not</span>
    <span class=\"ow\">in</span> <span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">keys</span><span class=\"p\">():</span>\n        <span class=\"n\">config</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;description&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"p\">{}</span>\n        <span class=\"n\">config</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;description&quot;</span><span class=\"p\">][</span><span
    class=\"s2\">&quot;len&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"mi\">160</span>\n\n    <span class=\"k\">if</span> <span class=\"s2\">&quot;long_description&quot;</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">keys</span><span class=\"p\">():</span>\n
    \       <span class=\"n\">config</span><span class=\"p\">[</span><span class=\"s2\">&quot;long_description&quot;</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"p\">{}</span>\n
    \       <span class=\"n\">config</span><span class=\"p\">[</span><span class=\"s2\">&quot;long_description&quot;</span><span
    class=\"p\">][</span><span class=\"s2\">&quot;len&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"mi\">250</span>\n\n    <span class=\"k\">def</span><span
    class=\"w\"> </span><span class=\"nf\">try_config_get</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span
    class=\"p\">:</span>\n        <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">return</span> <span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;len&quot;</span><span
    class=\"p\">)</span> <span class=\"ow\">or</span> <span class=\"kc\">None</span>\n
    \       <span class=\"k\">except</span> <span class=\"ne\">AttributeError</span><span
    class=\"p\">:</span>\n            <span class=\"k\">return</span> <span class=\"kc\">None</span>\n\n
    \   <span class=\"n\">max_description</span> <span class=\"o\">=</span> <span
    class=\"nb\">max</span><span class=\"p\">(</span>\n        <span class=\"p\">[</span>\n
    \           <span class=\"n\">value</span>\n            <span class=\"k\">for</span>
    <span class=\"n\">description_key</span> <span class=\"ow\">in</span> <span class=\"n\">config</span>\n
    \           <span class=\"k\">if</span> <span class=\"p\">(</span><span class=\"n\">value</span>
    <span class=\"o\">:=</span> <span class=\"n\">try_config_get</span><span class=\"p\">(</span><span
    class=\"n\">description_key</span><span class=\"p\">))</span>\n        <span class=\"p\">],</span>\n
    \   <span class=\"p\">)</span>\n\n    <span class=\"k\">with</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">cache</span> <span class=\"k\">as</span>
    <span class=\"n\">cache</span><span class=\"p\">:</span>\n        <span class=\"k\">for</span>
    <span class=\"n\">article</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">iter_articles</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;setting auto description&quot;</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">set_description</span><span class=\"p\">(</span>\n
    \               <span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"p\">,</span>\n                <span class=\"n\">article</span><span
    class=\"o\">=</span><span class=\"n\">article</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">cache</span><span class=\"o\">=</span><span
    class=\"n\">cache</span><span class=\"p\">,</span>\n                <span class=\"n\">config</span><span
    class=\"o\">=</span><span class=\"n\">config</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">max_description</span><span class=\"o\">=</span><span
    class=\"n\">max_description</span><span class=\"p\">,</span>\n                <span
    class=\"n\">plugin_text</span><span class=\"o\">=</span><span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"vm\">__file__</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">read_text</span><span class=\"p\">(),</span>\n
    \           <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n\n    </section>\n</article>
    \   </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>auto_description.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"The `markata.plugins.auto_description`
    plugin automatically generates descriptions for your\nposts by extracting text
    from the first paragraphs of your markdown \" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n
    \       <meta property=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    <head>\n        <script src=\"https://unpkg.com/@tailwindcss/browser@4\"></script>\n<title>auto_description.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"The `markata.plugins.auto_description`
    plugin automatically generates descriptions for your\nposts by extracting text
    from the first paragraphs of your markdown \" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n
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
    class=\"title\">\n    <h1 id=\"title\">\n        auto_description.py\n    </h1>\n</section></article>\n
    \    </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       auto_description.py\n    </h1>\n</section>    <section class=\"body\">\n
    \       <hr />\n<p>The <code>markata.plugins.auto_description</code> plugin automatically
    generates descriptions for your\nposts by extracting text from the first paragraphs
    of your markdown content. It can create\nmultiple descriptions of different lengths
    for different purposes (e.g., SEO, previews).</p>\n<h1 id=\"installation\">Installation
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
    class=\"p\">[</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;markata.plugins.auto_description&quot;</span><span
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
    class=\"p\">[</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;markata.plugins.auto_description&quot;</span><span
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>Configure multiple description
    types in your <code>markata.toml</code>:</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"c1\"># Standard
    description (160 characters)</span>\n<span class=\"k\">[markata.auto_description.description]</span>\n<span
    class=\"n\">len</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"mi\">160</span>\n\n<span class=\"c1\"># Longer
    description (250 characters)</span>\n<span class=\"k\">[markata.auto_description.long_description]</span>\n<span
    class=\"n\">len</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"mi\">250</span>\n\n<span class=\"c1\"># Full
    description (500 characters)</span>\n<span class=\"k\">[markata.auto_description.super_description]</span>\n<span
    class=\"n\">len</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"mi\">500</span>\n</pre></div>\n\n</pre>\n\n<p>Each
    description configuration:</p>\n<ul>\n<li>Must be under <code>markata.auto_description.[name]</code></li>\n<li>Requires
    a <code>len</code> parameter specifying maximum character length</li>\n<li>Creates
    an attribute named after the configuration key</li>\n</ul>\n<h2 id=\"frontmatter-override\">Frontmatter
    Override <a class=\"header-anchor\" href=\"#frontmatter-override\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>You can override automatic
    descriptions by setting them in frontmatter:</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span>---\ntitle: My Post\ndescription:
    My custom description\n<span class=\"gu\">long_description: A longer custom description</span>\n<span
    class=\"gu\">---</span>\n</pre></div>\n\n</pre>\n\n<h1 id=\"functionality\">Functionality
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<h2 id=\"description-generation\">Description
    Generation <a class=\"header-anchor\" href=\"#description-generation\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>The plugin:</p>\n<ol>\n<li>Converts
    markdown to plain text</li>\n<li>Finds the first meaningful paragraphs</li>\n<li>Truncates
    to the specified length</li>\n<li>Ensures clean word breaks</li>\n<li>Caches results
    for performance</li>\n</ol>\n<h2 id=\"registered-attributes\">Registered Attributes
    <a class=\"header-anchor\" href=\"#registered-attributes\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>For each configured
    description (e.g., <code>description</code>, <code>long_description</code>), the
    plugin adds:</p>\n<ul>\n<li>The description attribute with the generated/specified
    text</li>\n<li>Truncated to the configured length</li>\n<li>Preserving complete
    words</li>\n</ul>\n<h2 id=\"dependencies\">Dependencies <a class=\"header-anchor\"
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
    on:</p>\n<ul>\n<li>markdown-it-py for markdown parsing</li>\n</ul>\n<hr />\n<div
    class=\"admonition function\">\n<p class=\"admonition-title\">Function</p>\n<h2
    id=\"get_description\" class=\"admonition-title\" style=\"margin: 0; padding:
    .5rem 1rem;\">get_description <em class=\"small\">function</em></h2>\n<p>Get the
    full-length description for a single post by converting markdown to plain text.\nUses
    markdown-it-py to parse the markdown and extracts text content from all nodes.</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">get_description
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
    class=\"w\"> </span><span class=\"nf\">get_description</span><span class=\"p\">(</span><span
    class=\"n\">article</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span
    class=\"p\">:</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">    Get the full-length description for a single post by converting
    markdown to plain text.</span>\n<span class=\"sd\">    Uses markdown-it-py to
    parse the markdown and extracts text content from all nodes.</span>\n<span class=\"sd\">
    \   &quot;&quot;&quot;</span>\n    <span class=\"kn\">from</span><span class=\"w\">
    </span><span class=\"nn\">markdown_it</span><span class=\"w\"> </span><span class=\"kn\">import</span>
    <span class=\"n\">MarkdownIt</span>\n\n    <span class=\"n\">md</span> <span class=\"o\">=</span>
    <span class=\"n\">MarkdownIt</span><span class=\"p\">(</span><span class=\"s2\">&quot;commonmark&quot;</span><span
    class=\"p\">)</span>\n    <span class=\"n\">tokens</span> <span class=\"o\">=</span>
    <span class=\"n\">md</span><span class=\"o\">.</span><span class=\"n\">parse</span><span
    class=\"p\">(</span><span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">content</span><span class=\"p\">)</span>\n\n    <span class=\"c1\">#
    Extract text content from all inline tokens</span>\n    <span class=\"n\">text_chunks</span>
    <span class=\"o\">=</span> <span class=\"p\">[]</span>\n    <span class=\"k\">for</span>
    <span class=\"n\">token</span> <span class=\"ow\">in</span> <span class=\"n\">tokens</span><span
    class=\"p\">:</span>\n        <span class=\"k\">if</span> <span class=\"n\">token</span><span
    class=\"o\">.</span><span class=\"n\">type</span> <span class=\"o\">==</span>
    <span class=\"s2\">&quot;inline&quot;</span> <span class=\"ow\">and</span> <span
    class=\"n\">token</span><span class=\"o\">.</span><span class=\"n\">content</span><span
    class=\"p\">:</span>\n            <span class=\"n\">text_chunks</span><span class=\"o\">.</span><span
    class=\"n\">append</span><span class=\"p\">(</span><span class=\"n\">token</span><span
    class=\"o\">.</span><span class=\"n\">content</span><span class=\"p\">)</span>\n\n
    \   <span class=\"n\">description</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;
    &quot;</span><span class=\"o\">.</span><span class=\"n\">join</span><span class=\"p\">(</span><span
    class=\"n\">text_chunks</span><span class=\"p\">)</span>\n    <span class=\"k\">return</span>
    <span class=\"n\">html</span><span class=\"o\">.</span><span class=\"n\">escape</span><span
    class=\"p\">(</span><span class=\"n\">description</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition function\">\n<p class=\"admonition-title\">Function</p>\n<h2
    id=\"set_description\" class=\"admonition-title\" style=\"margin: 0; padding:
    .5rem 1rem;\">set_description <em class=\"small\">function</em></h2>\n<p>For a
    given <code>article</code>, find the description, put it in the cache, and set\nthe
    configured descriptions for the article.</p>\n</div>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">set_description
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
    class=\"w\"> </span><span class=\"nf\">set_description</span><span class=\"p\">(</span>\n
    \   <span class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span>\n    <span class=\"n\">article</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">,</span>\n    <span
    class=\"n\">cache</span><span class=\"p\">:</span> <span class=\"s2\">&quot;FanoutCache&quot;</span><span
    class=\"p\">,</span>\n    <span class=\"n\">config</span><span class=\"p\">:</span>
    <span class=\"n\">Dict</span><span class=\"p\">,</span>\n    <span class=\"n\">max_description</span><span
    class=\"p\">:</span> <span class=\"nb\">int</span> <span class=\"o\">=</span>
    <span class=\"mi\">500</span><span class=\"p\">,</span>\n    <span class=\"n\">plugin_text</span><span
    class=\"p\">:</span> <span class=\"kc\">None</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;&quot;</span><span class=\"p\">,</span>\n<span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
    class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span class=\"sd\">
    \   For a given `article`, find the description, put it in the cache, and set</span>\n<span
    class=\"sd\">    the configured descriptions for the article.</span>\n<span class=\"sd\">
    \   &quot;&quot;&quot;</span>\n    <span class=\"n\">key</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">make_hash</span><span
    class=\"p\">(</span>\n        <span class=\"s2\">&quot;auto_description&quot;</span><span
    class=\"p\">,</span>\n        <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">content</span><span class=\"p\">,</span>\n        <span class=\"n\">plugin_text</span><span
    class=\"p\">,</span>\n        <span class=\"n\">config</span><span class=\"p\">,</span>\n
    \   <span class=\"p\">)</span>\n\n    <span class=\"n\">description_from_cache</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">precache</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"n\">key</span><span class=\"p\">)</span>\n    <span
    class=\"k\">if</span> <span class=\"n\">description_from_cache</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n        <span class=\"n\">description</span>
    <span class=\"o\">=</span> <span class=\"n\">get_description</span><span class=\"p\">(</span><span
    class=\"n\">article</span><span class=\"p\">)[:</span><span class=\"n\">max_description</span><span
    class=\"p\">]</span>\n        <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">cache</span><span class=\"o\">.</span><span class=\"n\">set</span><span
    class=\"p\">(</span><span class=\"n\">key</span><span class=\"p\">,</span> <span
    class=\"n\">description</span><span class=\"p\">,</span> <span class=\"n\">expire</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">default_cache_expire</span><span
    class=\"p\">)</span>\n    <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \       <span class=\"n\">description</span> <span class=\"o\">=</span> <span
    class=\"n\">description_from_cache</span>\n\n    <span class=\"k\">for</span>
    <span class=\"n\">description_key</span> <span class=\"ow\">in</span> <span class=\"n\">config</span><span
    class=\"p\">:</span>\n        <span class=\"k\">if</span> <span class=\"n\">description_key</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"p\">[</span><span
    class=\"s2\">&quot;cache_expire&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;config_key&quot;</span><span
    class=\"p\">]:</span>\n            <span class=\"c1\"># overwrites missing (None)
    and empty (&#39;&#39;)</span>\n            <span class=\"k\">if</span> <span class=\"ow\">not</span>
    <span class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">metadata</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"n\">description_key</span><span class=\"p\">):</span>\n                <span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">metadata</span><span
    class=\"p\">[</span><span class=\"n\">description_key</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"n\">description</span><span class=\"p\">[</span>\n
    \                   <span class=\"p\">:</span> <span class=\"n\">config</span><span
    class=\"p\">[</span><span class=\"n\">description_key</span><span class=\"p\">][</span><span
    class=\"s2\">&quot;len&quot;</span><span class=\"p\">]</span>\n                <span
    class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div class=\"admonition
    function\">\n<p class=\"admonition-title\">Function</p>\n<h2 id=\"pre_render\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">pre_render
    <em class=\"small\">function</em></h2>\n<p>The Markata hook that will set descriptions
    for all posts in the pre-render phase.</p>\n</div>\n<div class=\"admonition source
    is-collapsible collapsible-open\">\n<p class=\"admonition-title\">pre_render <em
    class='small'>source</em></p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    class=\"w\"> </span><span class=\"nf\">pre_render</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">    The Markata hook that will set descriptions for all posts in
    the pre-render phase.</span>\n<span class=\"sd\">    &quot;&quot;&quot;</span>\n
    \   <span class=\"n\">config</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">get_plugin_config</span><span class=\"p\">(</span><span
    class=\"vm\">__file__</span><span class=\"p\">)</span>\n\n    <span class=\"k\">if</span>
    <span class=\"s2\">&quot;description&quot;</span> <span class=\"ow\">not</span>
    <span class=\"ow\">in</span> <span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">keys</span><span class=\"p\">():</span>\n        <span class=\"n\">config</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;description&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"p\">{}</span>\n        <span class=\"n\">config</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;description&quot;</span><span class=\"p\">][</span><span
    class=\"s2\">&quot;len&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"mi\">160</span>\n\n    <span class=\"k\">if</span> <span class=\"s2\">&quot;long_description&quot;</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">keys</span><span class=\"p\">():</span>\n
    \       <span class=\"n\">config</span><span class=\"p\">[</span><span class=\"s2\">&quot;long_description&quot;</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"p\">{}</span>\n
    \       <span class=\"n\">config</span><span class=\"p\">[</span><span class=\"s2\">&quot;long_description&quot;</span><span
    class=\"p\">][</span><span class=\"s2\">&quot;len&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"mi\">250</span>\n\n    <span class=\"k\">def</span><span
    class=\"w\"> </span><span class=\"nf\">try_config_get</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span
    class=\"p\">:</span>\n        <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">return</span> <span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;len&quot;</span><span
    class=\"p\">)</span> <span class=\"ow\">or</span> <span class=\"kc\">None</span>\n
    \       <span class=\"k\">except</span> <span class=\"ne\">AttributeError</span><span
    class=\"p\">:</span>\n            <span class=\"k\">return</span> <span class=\"kc\">None</span>\n\n
    \   <span class=\"n\">max_description</span> <span class=\"o\">=</span> <span
    class=\"nb\">max</span><span class=\"p\">(</span>\n        <span class=\"p\">[</span>\n
    \           <span class=\"n\">value</span>\n            <span class=\"k\">for</span>
    <span class=\"n\">description_key</span> <span class=\"ow\">in</span> <span class=\"n\">config</span>\n
    \           <span class=\"k\">if</span> <span class=\"p\">(</span><span class=\"n\">value</span>
    <span class=\"o\">:=</span> <span class=\"n\">try_config_get</span><span class=\"p\">(</span><span
    class=\"n\">description_key</span><span class=\"p\">))</span>\n        <span class=\"p\">],</span>\n
    \   <span class=\"p\">)</span>\n\n    <span class=\"k\">with</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">cache</span> <span class=\"k\">as</span>
    <span class=\"n\">cache</span><span class=\"p\">:</span>\n        <span class=\"k\">for</span>
    <span class=\"n\">article</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">iter_articles</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;setting auto description&quot;</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">set_description</span><span class=\"p\">(</span>\n
    \               <span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"p\">,</span>\n                <span class=\"n\">article</span><span
    class=\"o\">=</span><span class=\"n\">article</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">cache</span><span class=\"o\">=</span><span
    class=\"n\">cache</span><span class=\"p\">,</span>\n                <span class=\"n\">config</span><span
    class=\"o\">=</span><span class=\"n\">config</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">max_description</span><span class=\"o\">=</span><span
    class=\"n\">max_description</span><span class=\"p\">,</span>\n                <span
    class=\"n\">plugin_text</span><span class=\"o\">=</span><span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"vm\">__file__</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">read_text</span><span class=\"p\">(),</span>\n
    \           <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n\n    </section>\n</article>"
  raw.md: ''
published: false
slug: markata/plugins/auto-description
title: auto_description.py


---

---

The `markata.plugins.auto_description` plugin automatically generates descriptions for your
posts by extracting text from the first paragraphs of your markdown content. It can create
multiple descriptions of different lengths for different purposes (e.g., SEO, previews).

# Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.auto_description",
]
```

# Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.auto_description",
]
```

# Configuration

Configure multiple description types in your `markata.toml`:

```toml
# Standard description (160 characters)
[markata.auto_description.description]
len = 160

# Longer description (250 characters)
[markata.auto_description.long_description]
len = 250

# Full description (500 characters)
[markata.auto_description.super_description]
len = 500
```

Each description configuration:
- Must be under `markata.auto_description.[name]`
- Requires a `len` parameter specifying maximum character length
- Creates an attribute named after the configuration key

## Frontmatter Override

You can override automatic descriptions by setting them in frontmatter:

```markdown
---
title: My Post
description: My custom description
long_description: A longer custom description
---
```

# Functionality

## Description Generation

The plugin:
1. Converts markdown to plain text
2. Finds the first meaningful paragraphs
3. Truncates to the specified length
4. Ensures clean word breaks
5. Caches results for performance

## Registered Attributes

For each configured description (e.g., `description`, `long_description`), the plugin adds:
- The description attribute with the generated/specified text
- Truncated to the configured length
- Preserving complete words

## Dependencies

This plugin depends on:
- markdown-it-py for markdown parsing

---



!!! function
    <h2 id="get_description" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">get_description <em class="small">function</em></h2>

    Get the full-length description for a single post by converting markdown to plain text.
    Uses markdown-it-py to parse the markdown and extracts text content from all nodes.

???+ source "get_description <em class='small'>source</em>"
    ```python
    def get_description(article: "Post") -> str:
        """
        Get the full-length description for a single post by converting markdown to plain text.
        Uses markdown-it-py to parse the markdown and extracts text content from all nodes.
        """
        from markdown_it import MarkdownIt

        md = MarkdownIt("commonmark")
        tokens = md.parse(article.content)

        # Extract text content from all inline tokens
        text_chunks = []
        for token in tokens:
            if token.type == "inline" and token.content:
                text_chunks.append(token.content)

        description = " ".join(text_chunks)
        return html.escape(description)
    ```



!!! function
    <h2 id="set_description" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">set_description <em class="small">function</em></h2>

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
            markata.cache.set(key, description, expire=markata.config.default_cache_expire)
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



!!! function
    <h2 id="pre_render" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">pre_render <em class="small">function</em></h2>

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