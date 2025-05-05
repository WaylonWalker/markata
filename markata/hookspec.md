---
content: "---\n\nMarkata's hook specification system for plugin development.\n\n#
  Overview\n\nMarkata uses pluggy to define hooks that plugins can implement. These
  hooks allow plugins\nto modify Markata's behavior at specific points in the build
  process.\n\n# Hook Types\n\n## Configuration Hooks\n\nUsed to set up plugin configuration
  and models:\n\n```python\nfrom markata.hookspec import hook_impl, register_attr\n\n@hook_impl\n@register_attr(\"config\")\ndef
  config_model(markata):\n    \"\"\"Add plugin-specific config.\"\"\"\n    from pydantic
  import BaseModel\n\n    class MyConfig(BaseModel):\n        enabled: bool = True\n
  \       output_file: str = \"output.html\"\n\n    return {\"my_plugin\": MyConfig()}\n\n@hook_impl\n@register_attr(\"my_data\")\ndef
  configure(markata):\n    \"\"\"Initialize plugin using config.\"\"\"\n    if markata.config.my_plugin.enabled:\n
  \       # Set up plugin resources\n        markata.my_data = []\n```\n\n## Content
  Model Hooks\n\nDefine how content is structured:\n\n```python\n@hook_impl\n@register_attr(\"post_models\")\ndef
  post_model(markata):\n    \"\"\"Add fields to post model.\"\"\"\n    from pydantic
  import BaseModel, Field\n\n    class MyPostFields(BaseModel):\n        custom_date:
  str = Field(None, description=\"Custom date field\")\n        tags: list[str] =
  Field(default_factory=list)\n\n    return MyPostFields\n```\n\n## Content Processing
  Hooks\n\nHandle content transformation:\n\n```python\n@hook_impl(trylast=True)  #
  Run after other render hooks\ndef render(markata):\n    \"\"\"Process each article.\"\"\"\n
  \   for article in markata.filter(\"not skip\"):\n        # Add custom processing\n
  \       if article.tags:\n            article.tag_links = [f\"<a href='/tags/{tag}'>{tag}</a>\"\n
  \                              for tag in article.tags]\n```\n\n## Output Generation
  Hooks\n\nControl how content is saved:\n\n```python\n@hook_impl\ndef save(markata):\n
  \   \"\"\"Save processed content.\"\"\"\n    output_dir = Path(markata.config.output_dir)\n\n
  \   # Save custom index\n    if markata.config.my_plugin.enabled:\n        index
  = generate_custom_index(markata.articles)\n        (output_dir / \"custom.html\").write_text(index)\n```\n\n#
  Hook Ordering\n\nControl execution order with decorators:\n\n```python\n# Run first
  in configure stage\n@hook_impl(tryfirst=True)\ndef configure(markata): ...\n\n#
  Run in middle (default)\n@hook_impl\ndef configure(markata): ...\n\n# Run last in
  configure stage\n@hook_impl(trylast=True)\ndef configure(markata): ...\n```\n\n#
  Attribute Registration\n\nRegister data on the Markata instance:\n\n```python\n#
  Single attribute\n@register_attr(\"articles\")\ndef my_hook(markata):\n    markata.articles
  = []\n\n# Multiple attributes\n@register_attr(\"articles\", \"tags\", \"categories\")\ndef
  my_hook(markata):\n    markata.articles = []\n    markata.tags = {}\n    markata.categories
  = {}\n\n# Access in other hooks\n@hook_impl\ndef render(markata):\n    print(markata.articles)
  \ # Access registered data\n```\n\n# Complex Example\n\nHere's a complete plugin
  example combining multiple hooks:\n\n```python\nfrom pathlib import Path\nfrom typing
  import List, Optional\nfrom pydantic import BaseModel, Field\nfrom markata.hookspec
  import hook_impl, register_attr\n\nclass TagConfig(BaseModel):\n    \"\"\"Configuration
  for tag handling.\"\"\"\n    enabled: bool = True\n    min_posts: int = 2\n    output_dir:
  str = \"tags\"\n\nclass TaggedPost(BaseModel):\n    \"\"\"Add tag fields to posts.\"\"\"\n
  \   tags: List[str] = Field(default_factory=list)\n    tag_links: Optional[str]
  = None\n\n@hook_impl\n@register_attr(\"config_models\")\ndef config_model(markata):\n
  \   \"\"\"Add tag configuration.\"\"\"\n    markata.config_models.append(TagConfig)\n\n@hook_impl\n@register_attr(\"post_models\")\ndef
  post_model(markata):\n    \"\"\"Add tag fields to posts.\"\"\"\n    markata.post_models.append(TaggedPost)\n\n\n@hook_impl(trylast=True)\ndef
  render(markata):\n    \"\"\"Add tag links to articles.\"\"\"\n    if not markata.config.tags.enabled:\n
  \       return\n\n    for article in markata.filter(\"not skip\"):\n        article.tag_links
  = \" \".join(\n            f\"<a href='/tags/{tag}'>{tag}</a>\"\n            for
  tag in article.tags\n        )\n\n@hook_impl\ndef save(markata):\n    \"\"\"Generate
  tag pages.\"\"\"\n    if not markata.config.tags.enabled:\n        return\n\n    output_dir
  = Path(markata.config.output_dir)\n    tag_dir = output_dir / markata.config.tags.output_dir\n
  \   tag_dir.mkdir(exist_ok=True)\n\n    for tag, articles in markata.tags.items():\n
  \       if len(articles) >= markata.config.tags.min_posts:\n            content
  = generate_tag_page(tag, articles)\n            (tag_dir / f\"{tag}.html\").write_text(content)\n```\n\nThis
  example shows:\n1. Configuration definition\n2. Model extension\n3. Data processing\n4.
  Content generation\n5. Output handling\n\nSee [[ markata/lifecycle ]] for the exact
  order hooks are executed.\n\n---\n\n!!! class\n    <h2 id=\"MarkataSpecs\" class=\"admonition-title\"
  style=\"margin: 0; padding: .5rem 1rem;\">MarkataSpecs <em class=\"small\">class</em></h2>\n\n
  \   Namespace that defines all specifications for Load hooks.\n\n    configure ->
  glob -> load -> render -> save\n\n???+ source \"MarkataSpecs <em class='small'>source</em>\"\n
  \   ```python\n    class MarkataSpecs:\n        \"\"\"\n        Namespace that defines
  all specifications for Load hooks.\n\n        configure -> glob -> load -> render
  -> save\n        \"\"\"\n    ```\n!!! function\n    <h2 id=\"cli_lifecycle_method\"
  class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">cli_lifecycle_method
  <em class=\"small\">function</em></h2>\n\n    A Markata lifecycle methos that includes
  a typer app used for cli's\n\n???+ source \"cli_lifecycle_method <em class='small'>source</em>\"\n
  \   ```python\n    def cli_lifecycle_method(markata: \"Markata\", app: \"typer.Typer\")
  -> Any:\n        \"A Markata lifecycle methos that includes a typer app used for
  cli's\"\n    ```"
date: 2025-05-05
description: "Markata's hook specification system for plugin development. Overview
  Markata uses pluggy to define hooks that plugins can implement. These hooks allow
  plugins\u2026"
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>hookspec.py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"Markata's hook specification system for plugin
    development. Overview Markata uses pluggy to define hooks that plugins can implement.
    These hooks allow plugins\u2026\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n
    \       <meta property=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    <head>\n        <script src=\"https://unpkg.com/@tailwindcss/browser@4\"></script>\n<title>hookspec.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Markata's hook specification system for
    plugin development. Overview Markata uses pluggy to define hooks that plugins
    can implement. These hooks allow plugins\u2026\" />\n <link href=\"/favicon.ico\"
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        hookspec.py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <hr />\n<p>Markata's
    hook specification system for plugin development.</p>\n<h1 id=\"overview\">Overview
    <a class=\"header-anchor\" href=\"#overview\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>Markata uses pluggy
    to define hooks that plugins can implement. These hooks allow plugins\nto modify
    Markata's behavior at specific points in the build process.</p>\n<h1 id=\"hook-types\">Hook
    Types <a class=\"header-anchor\" href=\"#hook-types\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<h2 id=\"configuration-hooks\">Configuration
    Hooks <a class=\"header-anchor\" href=\"#configuration-hooks\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Used to set up plugin
    configuration and models:</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"kn\">from</span><span
    class=\"w\"> </span><span class=\"nn\">markata.hookspec</span><span class=\"w\">
    </span><span class=\"kn\">import</span> <span class=\"n\">hook_impl</span><span
    class=\"p\">,</span> <span class=\"n\">register_attr</span>\n\n<span class=\"nd\">@hook_impl</span>\n<span
    class=\"nd\">@register_attr</span><span class=\"p\">(</span><span class=\"s2\">&quot;config&quot;</span><span
    class=\"p\">)</span>\n<span class=\"k\">def</span><span class=\"w\"> </span><span
    class=\"nf\">config_model</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">):</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Add
    plugin-specific config.&quot;&quot;&quot;</span>\n    <span class=\"kn\">from</span><span
    class=\"w\"> </span><span class=\"nn\">pydantic</span><span class=\"w\"> </span><span
    class=\"kn\">import</span> <span class=\"n\">BaseModel</span>\n\n    <span class=\"k\">class</span><span
    class=\"w\"> </span><span class=\"nc\">MyConfig</span><span class=\"p\">(</span><span
    class=\"n\">BaseModel</span><span class=\"p\">):</span>\n        <span class=\"n\">enabled</span><span
    class=\"p\">:</span> <span class=\"nb\">bool</span> <span class=\"o\">=</span>
    <span class=\"kc\">True</span>\n        <span class=\"n\">output_file</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;output.html&quot;</span>\n\n    <span class=\"k\">return</span>
    <span class=\"p\">{</span><span class=\"s2\">&quot;my_plugin&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">MyConfig</span><span class=\"p\">()}</span>\n\n<span
    class=\"nd\">@hook_impl</span>\n<span class=\"nd\">@register_attr</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;my_data&quot;</span><span class=\"p\">)</span>\n<span
    class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">configure</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">):</span>\n<span
    class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Initialize plugin
    using config.&quot;&quot;&quot;</span>\n    <span class=\"k\">if</span> <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">my_plugin</span><span class=\"o\">.</span><span
    class=\"n\">enabled</span><span class=\"p\">:</span>\n        <span class=\"c1\">#
    Set up plugin resources</span>\n        <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">my_data</span> <span class=\"o\">=</span>
    <span class=\"p\">[]</span>\n</pre></div>\n\n</pre>\n\n<h2 id=\"content-model-hooks\">Content
    Model Hooks <a class=\"header-anchor\" href=\"#content-model-hooks\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Define how content is
    structured:</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"nd\">@hook_impl</span>\n<span
    class=\"nd\">@register_attr</span><span class=\"p\">(</span><span class=\"s2\">&quot;post_models&quot;</span><span
    class=\"p\">)</span>\n<span class=\"k\">def</span><span class=\"w\"> </span><span
    class=\"nf\">post_model</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">):</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Add
    fields to post model.&quot;&quot;&quot;</span>\n    <span class=\"kn\">from</span><span
    class=\"w\"> </span><span class=\"nn\">pydantic</span><span class=\"w\"> </span><span
    class=\"kn\">import</span> <span class=\"n\">BaseModel</span><span class=\"p\">,</span>
    <span class=\"n\">Field</span>\n\n    <span class=\"k\">class</span><span class=\"w\">
    </span><span class=\"nc\">MyPostFields</span><span class=\"p\">(</span><span class=\"n\">BaseModel</span><span
    class=\"p\">):</span>\n        <span class=\"n\">custom_date</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"n\">Field</span><span
    class=\"p\">(</span><span class=\"kc\">None</span><span class=\"p\">,</span> <span
    class=\"n\">description</span><span class=\"o\">=</span><span class=\"s2\">&quot;Custom
    date field&quot;</span><span class=\"p\">)</span>\n        <span class=\"n\">tags</span><span
    class=\"p\">:</span> <span class=\"nb\">list</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"n\">Field</span><span class=\"p\">(</span><span class=\"n\">default_factory</span><span
    class=\"o\">=</span><span class=\"nb\">list</span><span class=\"p\">)</span>\n\n
    \   <span class=\"k\">return</span> <span class=\"n\">MyPostFields</span>\n</pre></div>\n\n</pre>\n\n<h2
    id=\"content-processing-hooks\">Content Processing Hooks <a class=\"header-anchor\"
    href=\"#content-processing-hooks\"><svg class=\"heading-permalink\" aria-hidden=\"true\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Handle content transformation:</p>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"nd\">@hook_impl</span><span
    class=\"p\">(</span><span class=\"n\">trylast</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>  <span class=\"c1\"># Run after
    other render hooks</span>\n<span class=\"k\">def</span><span class=\"w\"> </span><span
    class=\"nf\">render</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">):</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Process
    each article.&quot;&quot;&quot;</span>\n    <span class=\"k\">for</span> <span
    class=\"n\">article</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">filter</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;not skip&quot;</span><span class=\"p\">):</span>\n        <span
    class=\"c1\"># Add custom processing</span>\n        <span class=\"k\">if</span>
    <span class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">tags</span><span
    class=\"p\">:</span>\n            <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">tag_links</span> <span class=\"o\">=</span> <span class=\"p\">[</span><span
    class=\"sa\">f</span><span class=\"s2\">&quot;&lt;a href=&#39;/tags/</span><span
    class=\"si\">{</span><span class=\"n\">tag</span><span class=\"si\">}</span><span
    class=\"s2\">&#39;&gt;</span><span class=\"si\">{</span><span class=\"n\">tag</span><span
    class=\"si\">}</span><span class=\"s2\">&lt;/a&gt;&quot;</span>\n                               <span
    class=\"k\">for</span> <span class=\"n\">tag</span> <span class=\"ow\">in</span>
    <span class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">tags</span><span
    class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<h2 id=\"output-generation-hooks\">Output
    Generation Hooks <a class=\"header-anchor\" href=\"#output-generation-hooks\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Control how content
    is saved:</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"nd\">@hook_impl</span>\n<span
    class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">save</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">):</span>\n<span
    class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Save processed content.&quot;&quot;&quot;</span>\n
    \   <span class=\"n\">output_dir</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">output_dir</span><span
    class=\"p\">)</span>\n\n    <span class=\"c1\"># Save custom index</span>\n    <span
    class=\"k\">if</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">my_plugin</span><span
    class=\"o\">.</span><span class=\"n\">enabled</span><span class=\"p\">:</span>\n
    \       <span class=\"n\">index</span> <span class=\"o\">=</span> <span class=\"n\">generate_custom_index</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">articles</span><span class=\"p\">)</span>\n        <span class=\"p\">(</span><span
    class=\"n\">output_dir</span> <span class=\"o\">/</span> <span class=\"s2\">&quot;custom.html&quot;</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">write_text</span><span
    class=\"p\">(</span><span class=\"n\">index</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<h1
    id=\"hook-ordering\">Hook Ordering <a class=\"header-anchor\" href=\"#hook-ordering\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>Control execution order
    with decorators:</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"c1\"># Run
    first in configure stage</span>\n<span class=\"nd\">@hook_impl</span><span class=\"p\">(</span><span
    class=\"n\">tryfirst</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n<span class=\"k\">def</span><span class=\"w\"> </span><span
    class=\"nf\">configure</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">):</span> <span class=\"o\">...</span>\n\n<span class=\"c1\"># Run
    in middle (default)</span>\n<span class=\"nd\">@hook_impl</span>\n<span class=\"k\">def</span><span
    class=\"w\"> </span><span class=\"nf\">configure</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">):</span> <span class=\"o\">...</span>\n\n<span
    class=\"c1\"># Run last in configure stage</span>\n<span class=\"nd\">@hook_impl</span><span
    class=\"p\">(</span><span class=\"n\">trylast</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n<span class=\"k\">def</span><span
    class=\"w\"> </span><span class=\"nf\">configure</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">):</span> <span class=\"o\">...</span>\n</pre></div>\n\n</pre>\n\n<h1
    id=\"attribute-registration\">Attribute Registration <a class=\"header-anchor\"
    href=\"#attribute-registration\"><svg class=\"heading-permalink\" aria-hidden=\"true\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>Register data on the
    Markata instance:</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"c1\"># Single
    attribute</span>\n<span class=\"nd\">@register_attr</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;articles&quot;</span><span class=\"p\">)</span>\n<span class=\"k\">def</span><span
    class=\"w\"> </span><span class=\"nf\">my_hook</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">):</span>\n    <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">articles</span> <span class=\"o\">=</span>
    <span class=\"p\">[]</span>\n\n<span class=\"c1\"># Multiple attributes</span>\n<span
    class=\"nd\">@register_attr</span><span class=\"p\">(</span><span class=\"s2\">&quot;articles&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;tags&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;categories&quot;</span><span class=\"p\">)</span>\n<span
    class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">my_hook</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">):</span>\n
    \   <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">articles</span>
    <span class=\"o\">=</span> <span class=\"p\">[]</span>\n    <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">tags</span> <span class=\"o\">=</span> <span
    class=\"p\">{}</span>\n    <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">categories</span> <span class=\"o\">=</span> <span class=\"p\">{}</span>\n\n<span
    class=\"c1\"># Access in other hooks</span>\n<span class=\"nd\">@hook_impl</span>\n<span
    class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">render</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">):</span>\n
    \   <span class=\"nb\">print</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">articles</span><span class=\"p\">)</span>
    \ <span class=\"c1\"># Access registered data</span>\n</pre></div>\n\n</pre>\n\n<h1
    id=\"complex-example\">Complex Example <a class=\"header-anchor\" href=\"#complex-example\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>Here's a complete plugin
    example combining multiple hooks:</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"kn\">from</span><span
    class=\"w\"> </span><span class=\"nn\">pathlib</span><span class=\"w\"> </span><span
    class=\"kn\">import</span> <span class=\"n\">Path</span>\n<span class=\"kn\">from</span><span
    class=\"w\"> </span><span class=\"nn\">typing</span><span class=\"w\"> </span><span
    class=\"kn\">import</span> <span class=\"n\">List</span><span class=\"p\">,</span>
    <span class=\"n\">Optional</span>\n<span class=\"kn\">from</span><span class=\"w\">
    </span><span class=\"nn\">pydantic</span><span class=\"w\"> </span><span class=\"kn\">import</span>
    <span class=\"n\">BaseModel</span><span class=\"p\">,</span> <span class=\"n\">Field</span>\n<span
    class=\"kn\">from</span><span class=\"w\"> </span><span class=\"nn\">markata.hookspec</span><span
    class=\"w\"> </span><span class=\"kn\">import</span> <span class=\"n\">hook_impl</span><span
    class=\"p\">,</span> <span class=\"n\">register_attr</span>\n\n<span class=\"k\">class</span><span
    class=\"w\"> </span><span class=\"nc\">TagConfig</span><span class=\"p\">(</span><span
    class=\"n\">BaseModel</span><span class=\"p\">):</span>\n<span class=\"w\">    </span><span
    class=\"sd\">&quot;&quot;&quot;Configuration for tag handling.&quot;&quot;&quot;</span>\n
    \   <span class=\"n\">enabled</span><span class=\"p\">:</span> <span class=\"nb\">bool</span>
    <span class=\"o\">=</span> <span class=\"kc\">True</span>\n    <span class=\"n\">min_posts</span><span
    class=\"p\">:</span> <span class=\"nb\">int</span> <span class=\"o\">=</span>
    <span class=\"mi\">2</span>\n    <span class=\"n\">output_dir</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;tags&quot;</span>\n\n<span
    class=\"k\">class</span><span class=\"w\"> </span><span class=\"nc\">TaggedPost</span><span
    class=\"p\">(</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n<span
    class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Add tag fields to
    posts.&quot;&quot;&quot;</span>\n    <span class=\"n\">tags</span><span class=\"p\">:</span>
    <span class=\"n\">List</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"n\">Field</span><span
    class=\"p\">(</span><span class=\"n\">default_factory</span><span class=\"o\">=</span><span
    class=\"nb\">list</span><span class=\"p\">)</span>\n    <span class=\"n\">tag_links</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"kc\">None</span>\n\n<span class=\"nd\">@hook_impl</span>\n<span class=\"nd\">@register_attr</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;config_models&quot;</span><span class=\"p\">)</span>\n<span
    class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">config_model</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">):</span>\n<span
    class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Add tag configuration.&quot;&quot;&quot;</span>\n
    \   <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config_models</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">TagConfig</span><span class=\"p\">)</span>\n\n<span class=\"nd\">@hook_impl</span>\n<span
    class=\"nd\">@register_attr</span><span class=\"p\">(</span><span class=\"s2\">&quot;post_models&quot;</span><span
    class=\"p\">)</span>\n<span class=\"k\">def</span><span class=\"w\"> </span><span
    class=\"nf\">post_model</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">):</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Add
    tag fields to posts.&quot;&quot;&quot;</span>\n    <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">post_models</span><span class=\"o\">.</span><span
    class=\"n\">append</span><span class=\"p\">(</span><span class=\"n\">TaggedPost</span><span
    class=\"p\">)</span>\n\n\n<span class=\"nd\">@hook_impl</span><span class=\"p\">(</span><span
    class=\"n\">trylast</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n<span class=\"k\">def</span><span class=\"w\"> </span><span
    class=\"nf\">render</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">):</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Add
    tag links to articles.&quot;&quot;&quot;</span>\n    <span class=\"k\">if</span>
    <span class=\"ow\">not</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">tags</span><span
    class=\"o\">.</span><span class=\"n\">enabled</span><span class=\"p\">:</span>\n
    \       <span class=\"k\">return</span>\n\n    <span class=\"k\">for</span> <span
    class=\"n\">article</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">filter</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;not skip&quot;</span><span class=\"p\">):</span>\n        <span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">tag_links</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot; &quot;</span><span class=\"o\">.</span><span
    class=\"n\">join</span><span class=\"p\">(</span>\n            <span class=\"sa\">f</span><span
    class=\"s2\">&quot;&lt;a href=&#39;/tags/</span><span class=\"si\">{</span><span
    class=\"n\">tag</span><span class=\"si\">}</span><span class=\"s2\">&#39;&gt;</span><span
    class=\"si\">{</span><span class=\"n\">tag</span><span class=\"si\">}</span><span
    class=\"s2\">&lt;/a&gt;&quot;</span>\n            <span class=\"k\">for</span>
    <span class=\"n\">tag</span> <span class=\"ow\">in</span> <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">tags</span>\n        <span class=\"p\">)</span>\n\n<span
    class=\"nd\">@hook_impl</span>\n<span class=\"k\">def</span><span class=\"w\">
    </span><span class=\"nf\">save</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">):</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Generate
    tag pages.&quot;&quot;&quot;</span>\n    <span class=\"k\">if</span> <span class=\"ow\">not</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">tags</span><span class=\"o\">.</span><span
    class=\"n\">enabled</span><span class=\"p\">:</span>\n        <span class=\"k\">return</span>\n\n
    \   <span class=\"n\">output_dir</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">output_dir</span><span
    class=\"p\">)</span>\n    <span class=\"n\">tag_dir</span> <span class=\"o\">=</span>
    <span class=\"n\">output_dir</span> <span class=\"o\">/</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">tags</span><span class=\"o\">.</span><span class=\"n\">output_dir</span>\n
    \   <span class=\"n\">tag_dir</span><span class=\"o\">.</span><span class=\"n\">mkdir</span><span
    class=\"p\">(</span><span class=\"n\">exist_ok</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n\n    <span class=\"k\">for</span>
    <span class=\"n\">tag</span><span class=\"p\">,</span> <span class=\"n\">articles</span>
    <span class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">tags</span><span class=\"o\">.</span><span class=\"n\">items</span><span
    class=\"p\">():</span>\n        <span class=\"k\">if</span> <span class=\"nb\">len</span><span
    class=\"p\">(</span><span class=\"n\">articles</span><span class=\"p\">)</span>
    <span class=\"o\">&gt;=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">tags</span><span
    class=\"o\">.</span><span class=\"n\">min_posts</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">content</span> <span class=\"o\">=</span> <span
    class=\"n\">generate_tag_page</span><span class=\"p\">(</span><span class=\"n\">tag</span><span
    class=\"p\">,</span> <span class=\"n\">articles</span><span class=\"p\">)</span>\n
    \           <span class=\"p\">(</span><span class=\"n\">tag_dir</span> <span class=\"o\">/</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">tag</span><span class=\"si\">}</span><span class=\"s2\">.html&quot;</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">write_text</span><span
    class=\"p\">(</span><span class=\"n\">content</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>This
    example shows:</p>\n<ol>\n<li>Configuration definition</li>\n<li>Model extension</li>\n<li>Data
    processing</li>\n<li>Content generation</li>\n<li>Output handling</li>\n</ol>\n<p>See
    <a class=\"wikilink\" href=\"/markata/lifecycle\">markata/lifecycle</a> for the
    exact order hooks are executed.</p>\n<hr />\n<div class=\"admonition class\">\n<p
    class=\"admonition-title\">Class</p>\n<h2 id=\"MarkataSpecs\" class=\"admonition-title\"
    style=\"margin: 0; padding: .5rem 1rem;\">MarkataSpecs <em class=\"small\">class</em></h2>\n<p>Namespace
    that defines all specifications for Load hooks.</p>\n<p>configure -&gt; glob -&gt;
    load -&gt; render -&gt; save</p>\n</div>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">MarkataSpecs <em class='small'>source</em></p>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">class</span><span
    class=\"w\"> </span><span class=\"nc\">MarkataSpecs</span><span class=\"p\">:</span>\n<span
    class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span class=\"sd\">
    \   Namespace that defines all specifications for Load hooks.</span>\n\n<span
    class=\"sd\">    configure -&gt; glob -&gt; load -&gt; render -&gt; save</span>\n<span
    class=\"sd\">    &quot;&quot;&quot;</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition function\">\n<p class=\"admonition-title\">Function</p>\n<h2
    id=\"cli_lifecycle_method\" class=\"admonition-title\" style=\"margin: 0; padding:
    .5rem 1rem;\">cli_lifecycle_method <em class=\"small\">function</em></h2>\n<p>A
    Markata lifecycle methos that includes a typer app used for cli's</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">cli_lifecycle_method
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
    class=\"w\"> </span><span class=\"nf\">cli_lifecycle_method</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">app</span><span class=\"p\">:</span> <span
    class=\"s2\">&quot;typer.Typer&quot;</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Any</span><span class=\"p\">:</span>\n    <span class=\"s2\">&quot;A
    Markata lifecycle methos that includes a typer app used for cli&#39;s&quot;</span>\n</pre></div>\n\n</pre>\n\n</div>\n\n
    \   </section>\n</article>    </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>hookspec.py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"Markata's hook specification system for plugin
    development. Overview Markata uses pluggy to define hooks that plugins can implement.
    These hooks allow plugins\u2026\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n
    \       <meta property=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    <head>\n        <script src=\"https://unpkg.com/@tailwindcss/browser@4\"></script>\n<title>hookspec.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Markata's hook specification system for
    plugin development. Overview Markata uses pluggy to define hooks that plugins
    can implement. These hooks allow plugins\u2026\" />\n <link href=\"/favicon.ico\"
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
    class=\"title\">\n    <h1 id=\"title\">\n        hookspec.py\n    </h1>\n</section></article>\n
    \    </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       hookspec.py\n    </h1>\n</section>    <section class=\"body\">\n        <hr
    />\n<p>Markata's hook specification system for plugin development.</p>\n<h1 id=\"overview\">Overview
    <a class=\"header-anchor\" href=\"#overview\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>Markata uses pluggy
    to define hooks that plugins can implement. These hooks allow plugins\nto modify
    Markata's behavior at specific points in the build process.</p>\n<h1 id=\"hook-types\">Hook
    Types <a class=\"header-anchor\" href=\"#hook-types\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<h2 id=\"configuration-hooks\">Configuration
    Hooks <a class=\"header-anchor\" href=\"#configuration-hooks\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Used to set up plugin
    configuration and models:</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"kn\">from</span><span
    class=\"w\"> </span><span class=\"nn\">markata.hookspec</span><span class=\"w\">
    </span><span class=\"kn\">import</span> <span class=\"n\">hook_impl</span><span
    class=\"p\">,</span> <span class=\"n\">register_attr</span>\n\n<span class=\"nd\">@hook_impl</span>\n<span
    class=\"nd\">@register_attr</span><span class=\"p\">(</span><span class=\"s2\">&quot;config&quot;</span><span
    class=\"p\">)</span>\n<span class=\"k\">def</span><span class=\"w\"> </span><span
    class=\"nf\">config_model</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">):</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Add
    plugin-specific config.&quot;&quot;&quot;</span>\n    <span class=\"kn\">from</span><span
    class=\"w\"> </span><span class=\"nn\">pydantic</span><span class=\"w\"> </span><span
    class=\"kn\">import</span> <span class=\"n\">BaseModel</span>\n\n    <span class=\"k\">class</span><span
    class=\"w\"> </span><span class=\"nc\">MyConfig</span><span class=\"p\">(</span><span
    class=\"n\">BaseModel</span><span class=\"p\">):</span>\n        <span class=\"n\">enabled</span><span
    class=\"p\">:</span> <span class=\"nb\">bool</span> <span class=\"o\">=</span>
    <span class=\"kc\">True</span>\n        <span class=\"n\">output_file</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;output.html&quot;</span>\n\n    <span class=\"k\">return</span>
    <span class=\"p\">{</span><span class=\"s2\">&quot;my_plugin&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">MyConfig</span><span class=\"p\">()}</span>\n\n<span
    class=\"nd\">@hook_impl</span>\n<span class=\"nd\">@register_attr</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;my_data&quot;</span><span class=\"p\">)</span>\n<span
    class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">configure</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">):</span>\n<span
    class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Initialize plugin
    using config.&quot;&quot;&quot;</span>\n    <span class=\"k\">if</span> <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">my_plugin</span><span class=\"o\">.</span><span
    class=\"n\">enabled</span><span class=\"p\">:</span>\n        <span class=\"c1\">#
    Set up plugin resources</span>\n        <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">my_data</span> <span class=\"o\">=</span>
    <span class=\"p\">[]</span>\n</pre></div>\n\n</pre>\n\n<h2 id=\"content-model-hooks\">Content
    Model Hooks <a class=\"header-anchor\" href=\"#content-model-hooks\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Define how content is
    structured:</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"nd\">@hook_impl</span>\n<span
    class=\"nd\">@register_attr</span><span class=\"p\">(</span><span class=\"s2\">&quot;post_models&quot;</span><span
    class=\"p\">)</span>\n<span class=\"k\">def</span><span class=\"w\"> </span><span
    class=\"nf\">post_model</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">):</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Add
    fields to post model.&quot;&quot;&quot;</span>\n    <span class=\"kn\">from</span><span
    class=\"w\"> </span><span class=\"nn\">pydantic</span><span class=\"w\"> </span><span
    class=\"kn\">import</span> <span class=\"n\">BaseModel</span><span class=\"p\">,</span>
    <span class=\"n\">Field</span>\n\n    <span class=\"k\">class</span><span class=\"w\">
    </span><span class=\"nc\">MyPostFields</span><span class=\"p\">(</span><span class=\"n\">BaseModel</span><span
    class=\"p\">):</span>\n        <span class=\"n\">custom_date</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"n\">Field</span><span
    class=\"p\">(</span><span class=\"kc\">None</span><span class=\"p\">,</span> <span
    class=\"n\">description</span><span class=\"o\">=</span><span class=\"s2\">&quot;Custom
    date field&quot;</span><span class=\"p\">)</span>\n        <span class=\"n\">tags</span><span
    class=\"p\">:</span> <span class=\"nb\">list</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"n\">Field</span><span class=\"p\">(</span><span class=\"n\">default_factory</span><span
    class=\"o\">=</span><span class=\"nb\">list</span><span class=\"p\">)</span>\n\n
    \   <span class=\"k\">return</span> <span class=\"n\">MyPostFields</span>\n</pre></div>\n\n</pre>\n\n<h2
    id=\"content-processing-hooks\">Content Processing Hooks <a class=\"header-anchor\"
    href=\"#content-processing-hooks\"><svg class=\"heading-permalink\" aria-hidden=\"true\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Handle content transformation:</p>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"nd\">@hook_impl</span><span
    class=\"p\">(</span><span class=\"n\">trylast</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>  <span class=\"c1\"># Run after
    other render hooks</span>\n<span class=\"k\">def</span><span class=\"w\"> </span><span
    class=\"nf\">render</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">):</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Process
    each article.&quot;&quot;&quot;</span>\n    <span class=\"k\">for</span> <span
    class=\"n\">article</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">filter</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;not skip&quot;</span><span class=\"p\">):</span>\n        <span
    class=\"c1\"># Add custom processing</span>\n        <span class=\"k\">if</span>
    <span class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">tags</span><span
    class=\"p\">:</span>\n            <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">tag_links</span> <span class=\"o\">=</span> <span class=\"p\">[</span><span
    class=\"sa\">f</span><span class=\"s2\">&quot;&lt;a href=&#39;/tags/</span><span
    class=\"si\">{</span><span class=\"n\">tag</span><span class=\"si\">}</span><span
    class=\"s2\">&#39;&gt;</span><span class=\"si\">{</span><span class=\"n\">tag</span><span
    class=\"si\">}</span><span class=\"s2\">&lt;/a&gt;&quot;</span>\n                               <span
    class=\"k\">for</span> <span class=\"n\">tag</span> <span class=\"ow\">in</span>
    <span class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">tags</span><span
    class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<h2 id=\"output-generation-hooks\">Output
    Generation Hooks <a class=\"header-anchor\" href=\"#output-generation-hooks\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Control how content
    is saved:</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"nd\">@hook_impl</span>\n<span
    class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">save</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">):</span>\n<span
    class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Save processed content.&quot;&quot;&quot;</span>\n
    \   <span class=\"n\">output_dir</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">output_dir</span><span
    class=\"p\">)</span>\n\n    <span class=\"c1\"># Save custom index</span>\n    <span
    class=\"k\">if</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">my_plugin</span><span
    class=\"o\">.</span><span class=\"n\">enabled</span><span class=\"p\">:</span>\n
    \       <span class=\"n\">index</span> <span class=\"o\">=</span> <span class=\"n\">generate_custom_index</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">articles</span><span class=\"p\">)</span>\n        <span class=\"p\">(</span><span
    class=\"n\">output_dir</span> <span class=\"o\">/</span> <span class=\"s2\">&quot;custom.html&quot;</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">write_text</span><span
    class=\"p\">(</span><span class=\"n\">index</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<h1
    id=\"hook-ordering\">Hook Ordering <a class=\"header-anchor\" href=\"#hook-ordering\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>Control execution order
    with decorators:</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"c1\"># Run
    first in configure stage</span>\n<span class=\"nd\">@hook_impl</span><span class=\"p\">(</span><span
    class=\"n\">tryfirst</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n<span class=\"k\">def</span><span class=\"w\"> </span><span
    class=\"nf\">configure</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">):</span> <span class=\"o\">...</span>\n\n<span class=\"c1\"># Run
    in middle (default)</span>\n<span class=\"nd\">@hook_impl</span>\n<span class=\"k\">def</span><span
    class=\"w\"> </span><span class=\"nf\">configure</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">):</span> <span class=\"o\">...</span>\n\n<span
    class=\"c1\"># Run last in configure stage</span>\n<span class=\"nd\">@hook_impl</span><span
    class=\"p\">(</span><span class=\"n\">trylast</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n<span class=\"k\">def</span><span
    class=\"w\"> </span><span class=\"nf\">configure</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">):</span> <span class=\"o\">...</span>\n</pre></div>\n\n</pre>\n\n<h1
    id=\"attribute-registration\">Attribute Registration <a class=\"header-anchor\"
    href=\"#attribute-registration\"><svg class=\"heading-permalink\" aria-hidden=\"true\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>Register data on the
    Markata instance:</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"c1\"># Single
    attribute</span>\n<span class=\"nd\">@register_attr</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;articles&quot;</span><span class=\"p\">)</span>\n<span class=\"k\">def</span><span
    class=\"w\"> </span><span class=\"nf\">my_hook</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">):</span>\n    <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">articles</span> <span class=\"o\">=</span>
    <span class=\"p\">[]</span>\n\n<span class=\"c1\"># Multiple attributes</span>\n<span
    class=\"nd\">@register_attr</span><span class=\"p\">(</span><span class=\"s2\">&quot;articles&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;tags&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;categories&quot;</span><span class=\"p\">)</span>\n<span
    class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">my_hook</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">):</span>\n
    \   <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">articles</span>
    <span class=\"o\">=</span> <span class=\"p\">[]</span>\n    <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">tags</span> <span class=\"o\">=</span> <span
    class=\"p\">{}</span>\n    <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">categories</span> <span class=\"o\">=</span> <span class=\"p\">{}</span>\n\n<span
    class=\"c1\"># Access in other hooks</span>\n<span class=\"nd\">@hook_impl</span>\n<span
    class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">render</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">):</span>\n
    \   <span class=\"nb\">print</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">articles</span><span class=\"p\">)</span>
    \ <span class=\"c1\"># Access registered data</span>\n</pre></div>\n\n</pre>\n\n<h1
    id=\"complex-example\">Complex Example <a class=\"header-anchor\" href=\"#complex-example\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>Here's a complete plugin
    example combining multiple hooks:</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"kn\">from</span><span
    class=\"w\"> </span><span class=\"nn\">pathlib</span><span class=\"w\"> </span><span
    class=\"kn\">import</span> <span class=\"n\">Path</span>\n<span class=\"kn\">from</span><span
    class=\"w\"> </span><span class=\"nn\">typing</span><span class=\"w\"> </span><span
    class=\"kn\">import</span> <span class=\"n\">List</span><span class=\"p\">,</span>
    <span class=\"n\">Optional</span>\n<span class=\"kn\">from</span><span class=\"w\">
    </span><span class=\"nn\">pydantic</span><span class=\"w\"> </span><span class=\"kn\">import</span>
    <span class=\"n\">BaseModel</span><span class=\"p\">,</span> <span class=\"n\">Field</span>\n<span
    class=\"kn\">from</span><span class=\"w\"> </span><span class=\"nn\">markata.hookspec</span><span
    class=\"w\"> </span><span class=\"kn\">import</span> <span class=\"n\">hook_impl</span><span
    class=\"p\">,</span> <span class=\"n\">register_attr</span>\n\n<span class=\"k\">class</span><span
    class=\"w\"> </span><span class=\"nc\">TagConfig</span><span class=\"p\">(</span><span
    class=\"n\">BaseModel</span><span class=\"p\">):</span>\n<span class=\"w\">    </span><span
    class=\"sd\">&quot;&quot;&quot;Configuration for tag handling.&quot;&quot;&quot;</span>\n
    \   <span class=\"n\">enabled</span><span class=\"p\">:</span> <span class=\"nb\">bool</span>
    <span class=\"o\">=</span> <span class=\"kc\">True</span>\n    <span class=\"n\">min_posts</span><span
    class=\"p\">:</span> <span class=\"nb\">int</span> <span class=\"o\">=</span>
    <span class=\"mi\">2</span>\n    <span class=\"n\">output_dir</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;tags&quot;</span>\n\n<span
    class=\"k\">class</span><span class=\"w\"> </span><span class=\"nc\">TaggedPost</span><span
    class=\"p\">(</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n<span
    class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Add tag fields to
    posts.&quot;&quot;&quot;</span>\n    <span class=\"n\">tags</span><span class=\"p\">:</span>
    <span class=\"n\">List</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"n\">Field</span><span
    class=\"p\">(</span><span class=\"n\">default_factory</span><span class=\"o\">=</span><span
    class=\"nb\">list</span><span class=\"p\">)</span>\n    <span class=\"n\">tag_links</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"kc\">None</span>\n\n<span class=\"nd\">@hook_impl</span>\n<span class=\"nd\">@register_attr</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;config_models&quot;</span><span class=\"p\">)</span>\n<span
    class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">config_model</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">):</span>\n<span
    class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Add tag configuration.&quot;&quot;&quot;</span>\n
    \   <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config_models</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">TagConfig</span><span class=\"p\">)</span>\n\n<span class=\"nd\">@hook_impl</span>\n<span
    class=\"nd\">@register_attr</span><span class=\"p\">(</span><span class=\"s2\">&quot;post_models&quot;</span><span
    class=\"p\">)</span>\n<span class=\"k\">def</span><span class=\"w\"> </span><span
    class=\"nf\">post_model</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">):</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Add
    tag fields to posts.&quot;&quot;&quot;</span>\n    <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">post_models</span><span class=\"o\">.</span><span
    class=\"n\">append</span><span class=\"p\">(</span><span class=\"n\">TaggedPost</span><span
    class=\"p\">)</span>\n\n\n<span class=\"nd\">@hook_impl</span><span class=\"p\">(</span><span
    class=\"n\">trylast</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n<span class=\"k\">def</span><span class=\"w\"> </span><span
    class=\"nf\">render</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">):</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Add
    tag links to articles.&quot;&quot;&quot;</span>\n    <span class=\"k\">if</span>
    <span class=\"ow\">not</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">tags</span><span
    class=\"o\">.</span><span class=\"n\">enabled</span><span class=\"p\">:</span>\n
    \       <span class=\"k\">return</span>\n\n    <span class=\"k\">for</span> <span
    class=\"n\">article</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">filter</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;not skip&quot;</span><span class=\"p\">):</span>\n        <span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">tag_links</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot; &quot;</span><span class=\"o\">.</span><span
    class=\"n\">join</span><span class=\"p\">(</span>\n            <span class=\"sa\">f</span><span
    class=\"s2\">&quot;&lt;a href=&#39;/tags/</span><span class=\"si\">{</span><span
    class=\"n\">tag</span><span class=\"si\">}</span><span class=\"s2\">&#39;&gt;</span><span
    class=\"si\">{</span><span class=\"n\">tag</span><span class=\"si\">}</span><span
    class=\"s2\">&lt;/a&gt;&quot;</span>\n            <span class=\"k\">for</span>
    <span class=\"n\">tag</span> <span class=\"ow\">in</span> <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">tags</span>\n        <span class=\"p\">)</span>\n\n<span
    class=\"nd\">@hook_impl</span>\n<span class=\"k\">def</span><span class=\"w\">
    </span><span class=\"nf\">save</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">):</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Generate
    tag pages.&quot;&quot;&quot;</span>\n    <span class=\"k\">if</span> <span class=\"ow\">not</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">tags</span><span class=\"o\">.</span><span
    class=\"n\">enabled</span><span class=\"p\">:</span>\n        <span class=\"k\">return</span>\n\n
    \   <span class=\"n\">output_dir</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">output_dir</span><span
    class=\"p\">)</span>\n    <span class=\"n\">tag_dir</span> <span class=\"o\">=</span>
    <span class=\"n\">output_dir</span> <span class=\"o\">/</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">tags</span><span class=\"o\">.</span><span class=\"n\">output_dir</span>\n
    \   <span class=\"n\">tag_dir</span><span class=\"o\">.</span><span class=\"n\">mkdir</span><span
    class=\"p\">(</span><span class=\"n\">exist_ok</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n\n    <span class=\"k\">for</span>
    <span class=\"n\">tag</span><span class=\"p\">,</span> <span class=\"n\">articles</span>
    <span class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">tags</span><span class=\"o\">.</span><span class=\"n\">items</span><span
    class=\"p\">():</span>\n        <span class=\"k\">if</span> <span class=\"nb\">len</span><span
    class=\"p\">(</span><span class=\"n\">articles</span><span class=\"p\">)</span>
    <span class=\"o\">&gt;=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">tags</span><span
    class=\"o\">.</span><span class=\"n\">min_posts</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">content</span> <span class=\"o\">=</span> <span
    class=\"n\">generate_tag_page</span><span class=\"p\">(</span><span class=\"n\">tag</span><span
    class=\"p\">,</span> <span class=\"n\">articles</span><span class=\"p\">)</span>\n
    \           <span class=\"p\">(</span><span class=\"n\">tag_dir</span> <span class=\"o\">/</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">tag</span><span class=\"si\">}</span><span class=\"s2\">.html&quot;</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">write_text</span><span
    class=\"p\">(</span><span class=\"n\">content</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>This
    example shows:</p>\n<ol>\n<li>Configuration definition</li>\n<li>Model extension</li>\n<li>Data
    processing</li>\n<li>Content generation</li>\n<li>Output handling</li>\n</ol>\n<p>See
    <a class=\"wikilink\" href=\"/markata/lifecycle\">markata/lifecycle</a> for the
    exact order hooks are executed.</p>\n<hr />\n<div class=\"admonition class\">\n<p
    class=\"admonition-title\">Class</p>\n<h2 id=\"MarkataSpecs\" class=\"admonition-title\"
    style=\"margin: 0; padding: .5rem 1rem;\">MarkataSpecs <em class=\"small\">class</em></h2>\n<p>Namespace
    that defines all specifications for Load hooks.</p>\n<p>configure -&gt; glob -&gt;
    load -&gt; render -&gt; save</p>\n</div>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">MarkataSpecs <em class='small'>source</em></p>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">class</span><span
    class=\"w\"> </span><span class=\"nc\">MarkataSpecs</span><span class=\"p\">:</span>\n<span
    class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span class=\"sd\">
    \   Namespace that defines all specifications for Load hooks.</span>\n\n<span
    class=\"sd\">    configure -&gt; glob -&gt; load -&gt; render -&gt; save</span>\n<span
    class=\"sd\">    &quot;&quot;&quot;</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition function\">\n<p class=\"admonition-title\">Function</p>\n<h2
    id=\"cli_lifecycle_method\" class=\"admonition-title\" style=\"margin: 0; padding:
    .5rem 1rem;\">cli_lifecycle_method <em class=\"small\">function</em></h2>\n<p>A
    Markata lifecycle methos that includes a typer app used for cli's</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">cli_lifecycle_method
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
    class=\"w\"> </span><span class=\"nf\">cli_lifecycle_method</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">app</span><span class=\"p\">:</span> <span
    class=\"s2\">&quot;typer.Typer&quot;</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Any</span><span class=\"p\">:</span>\n    <span class=\"s2\">&quot;A
    Markata lifecycle methos that includes a typer app used for cli&#39;s&quot;</span>\n</pre></div>\n\n</pre>\n\n</div>\n\n
    \   </section>\n</article>"
  raw.md: ''
published: false
slug: markata/hookspec
title: hookspec.py


---

---

Markata's hook specification system for plugin development.

# Overview

Markata uses pluggy to define hooks that plugins can implement. These hooks allow plugins
to modify Markata's behavior at specific points in the build process.

# Hook Types

## Configuration Hooks

Used to set up plugin configuration and models:

```python
from markata.hookspec import hook_impl, register_attr

@hook_impl
@register_attr("config")
def config_model(markata):
    """Add plugin-specific config."""
    from pydantic import BaseModel

    class MyConfig(BaseModel):
        enabled: bool = True
        output_file: str = "output.html"

    return {"my_plugin": MyConfig()}

@hook_impl
@register_attr("my_data")
def configure(markata):
    """Initialize plugin using config."""
    if markata.config.my_plugin.enabled:
        # Set up plugin resources
        markata.my_data = []
```

## Content Model Hooks

Define how content is structured:

```python
@hook_impl
@register_attr("post_models")
def post_model(markata):
    """Add fields to post model."""
    from pydantic import BaseModel, Field

    class MyPostFields(BaseModel):
        custom_date: str = Field(None, description="Custom date field")
        tags: list[str] = Field(default_factory=list)

    return MyPostFields
```

## Content Processing Hooks

Handle content transformation:

```python
@hook_impl(trylast=True)  # Run after other render hooks
def render(markata):
    """Process each article."""
    for article in markata.filter("not skip"):
        # Add custom processing
        if article.tags:
            article.tag_links = [f"<a href='/tags/{tag}'>{tag}</a>"
                               for tag in article.tags]
```

## Output Generation Hooks

Control how content is saved:

```python
@hook_impl
def save(markata):
    """Save processed content."""
    output_dir = Path(markata.config.output_dir)

    # Save custom index
    if markata.config.my_plugin.enabled:
        index = generate_custom_index(markata.articles)
        (output_dir / "custom.html").write_text(index)
```

# Hook Ordering

Control execution order with decorators:

```python
# Run first in configure stage
@hook_impl(tryfirst=True)
def configure(markata): ...

# Run in middle (default)
@hook_impl
def configure(markata): ...

# Run last in configure stage
@hook_impl(trylast=True)
def configure(markata): ...
```

# Attribute Registration

Register data on the Markata instance:

```python
# Single attribute
@register_attr("articles")
def my_hook(markata):
    markata.articles = []

# Multiple attributes
@register_attr("articles", "tags", "categories")
def my_hook(markata):
    markata.articles = []
    markata.tags = {}
    markata.categories = {}

# Access in other hooks
@hook_impl
def render(markata):
    print(markata.articles)  # Access registered data
```

# Complex Example

Here's a complete plugin example combining multiple hooks:

```python
from pathlib import Path
from typing import List, Optional
from pydantic import BaseModel, Field
from markata.hookspec import hook_impl, register_attr

class TagConfig(BaseModel):
    """Configuration for tag handling."""
    enabled: bool = True
    min_posts: int = 2
    output_dir: str = "tags"

class TaggedPost(BaseModel):
    """Add tag fields to posts."""
    tags: List[str] = Field(default_factory=list)
    tag_links: Optional[str] = None

@hook_impl
@register_attr("config_models")
def config_model(markata):
    """Add tag configuration."""
    markata.config_models.append(TagConfig)

@hook_impl
@register_attr("post_models")
def post_model(markata):
    """Add tag fields to posts."""
    markata.post_models.append(TaggedPost)


@hook_impl(trylast=True)
def render(markata):
    """Add tag links to articles."""
    if not markata.config.tags.enabled:
        return

    for article in markata.filter("not skip"):
        article.tag_links = " ".join(
            f"<a href='/tags/{tag}'>{tag}</a>"
            for tag in article.tags
        )

@hook_impl
def save(markata):
    """Generate tag pages."""
    if not markata.config.tags.enabled:
        return

    output_dir = Path(markata.config.output_dir)
    tag_dir = output_dir / markata.config.tags.output_dir
    tag_dir.mkdir(exist_ok=True)

    for tag, articles in markata.tags.items():
        if len(articles) >= markata.config.tags.min_posts:
            content = generate_tag_page(tag, articles)
            (tag_dir / f"{tag}.html").write_text(content)
```

This example shows:
1. Configuration definition
2. Model extension
3. Data processing
4. Content generation
5. Output handling

See [[ markata/lifecycle ]] for the exact order hooks are executed.

---

!!! class
    <h2 id="MarkataSpecs" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">MarkataSpecs <em class="small">class</em></h2>

    Namespace that defines all specifications for Load hooks.

    configure -> glob -> load -> render -> save

???+ source "MarkataSpecs <em class='small'>source</em>"
    ```python
    class MarkataSpecs:
        """
        Namespace that defines all specifications for Load hooks.

        configure -> glob -> load -> render -> save
        """
    ```
!!! function
    <h2 id="cli_lifecycle_method" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">cli_lifecycle_method <em class="small">function</em></h2>

    A Markata lifecycle methos that includes a typer app used for cli's

???+ source "cli_lifecycle_method <em class='small'>source</em>"
    ```python
    def cli_lifecycle_method(markata: "Markata", app: "typer.Typer") -> Any:
        "A Markata lifecycle methos that includes a typer app used for cli's"
    ```