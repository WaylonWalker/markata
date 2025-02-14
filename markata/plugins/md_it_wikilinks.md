---
content: "---\n\nThe `markata.plugins.md_it_wikilinks` plugin adds support for wiki-style
  links using\ndouble brackets (`[[link]]`). It automatically resolves links to other
  posts in your\nsite using file names or slugs.\n\n# Installation\n\nThis plugin
  is built-in and enabled by default through the 'default' plugin.\nIf you want to
  be explicit, you can add it to your list of plugins:\n\n```toml\nhooks = [\n    \"markata.plugins.md_it_wikilinks\",\n]\n```\n\n#
  Uninstallation\n\nSince this plugin is included in the default plugin set, to disable
  it you must explicitly\nadd it to the disabled_hooks list if you are using the 'default'
  plugin:\n\n```toml\ndisabled_hooks = [\n    \"markata.plugins.md_it_wikilinks\",\n]\n```\n\n#
  Configuration\n\nThis plugin requires no explicit configuration. It automatically
  processes wikilinks\nin your markdown content.\n\n# Functionality\n\n## Basic Wikilinks\n\nSimple
  file-based linking:\n```markdown\n[[nav]]              -> links to docs/nav.md as
  /nav\n[[blog/post]]        -> links to blog/post.md as /blog/post\n[[about|About
  Me]]   -> links to about.md with \"About Me\" as text\n```\n\n## Smart Slug Resolution\n\nThe
  plugin:\n1. Looks up the target file in your content\n2. Finds its generated slug\n3.
  Creates a link to the final URL\n\nExample:\n```markdown\n# File: posts/2024-01-my-post.md\nslug:
  /blog/my-post\n\n# In another file:\n[[2024-01-my-post]]  -> links to /blog/my-post\n```\n\n##
  Link Formats\n\nSupports multiple link styles:\n- Basic: `[[filename]]`\n- With
  text: `[[filename|Link Text]]`\n- With path: `[[folder/file]]`\n- With extension:
  `[[file.md]]` (extension stripped in output)\n\n## HTML Output\n\nGenerated HTML
  structure:\n```html\n<a class=\"wikilink\" href=\"/target-slug\">Link Text</a>\n```\n\n##
  Error Handling\n\nFor broken links:\n- Maintains the wikilink syntax\n- Adds a 'broken-link'
  class\n- Optionally logs warnings\n\n## Dependencies\n\nThis plugin depends on:\n-
  markdown-it-py for markdown parsing\n- The `render_markdown` plugin for final HTML
  output\n\n---\n\n\n\n\n\n!!! function\n    <h2 id=\"wikilinks_plugin\" class=\"admonition-title\"
  style=\"margin: 0; padding: .5rem 1rem;\">wikilinks_plugin <em class=\"small\">function</em></h2>\n\n
  \   A plugin to create wikilinks tokens.\n    These, token should be handled by
  the renderer.\n\n    ???+ example\n\n        ```md title=markdown\n        [[nav]]\n
  \       ```\n\n        ```html title=html\n        <a class=\"wikilink\" href=\"/nav\">load</a>\n
  \       ```\n\n???+ source \"wikilinks_plugin <em class='small'>source</em>\"\n
  \   ```python\n    def wikilinks_plugin(\n        md: MarkdownIt,\n        start_delimiter:
  str = \"[\",\n        end_delimiter: str = \"]\",\n        markata=None,\n    ):\n
  \       \"\"\"A plugin to create wikilinks tokens.\n        These, token should
  be handled by the renderer.\n\n        ???+ example\n\n            ```md title=markdown\n
  \           [[nav]]\n            ```\n\n            ```html title=html\n            <a
  class=\"wikilink\" href=\"/nav\">load</a>\n            ```\n        \"\"\"\n\n        start_char
  = ord(start_delimiter)\n        end_char = ord(end_delimiter)\n\n        def _wikilinks_inline(state:
  StateInline, silent: bool):\n            try:\n                if (\n                    state.srcCharCode[state.pos]
  != start_char\n                    or state.srcCharCode[state.pos + 1] != start_char\n
  \               ):\n                    return False\n            except IndexError:\n
  \               return False\n\n            pos = state.pos + 2\n            found_closing
  = False\n            while True:\n                try:\n                    end
  = state.srcCharCode.index(end_char, pos)\n                except ValueError:\n                    return
  False\n                try:\n                    if state.srcCharCode[end + 1] ==
  end_char:\n                        found_closing = True\n                        break\n
  \               except IndexError:\n                    return False\n                pos
  = end + 2\n\n            if not found_closing:\n                return False\n\n
  \           text = state.src[state.pos + 2 : end].strip()\n            state.pos
  = end + 2\n\n            if silent:\n                return True\n\n            token
  = state.push(\"link_open\", \"a\", 1)\n            token.block = False\n            token.attrSet(\"class\",
  \"wikilink\")\n            if \"#\" in text:\n                link, id = text.split(\"#\")\n
  \               link = link.strip(\"/\")\n            else:\n                link,
  id = text, None\n\n            # possible_pages = markata.filter(\n            #
  \    f'str(path).split(\"/\")[-1].split(\".\")[0].replace(\"_\", \"-\") == \"{link.replace(\"_\",
  \"-\")}\"',\n            # )\n            possible_pages = markata.possible_wikilink.get(link,
  [])\n            if len(possible_pages) == 1:\n                link = possible_pages[0]\n
  \           elif len(possible_pages) > 1:\n                if md.options[\"article\"]
  is None:\n                    debug_value = \"UNKNOWN\"\n                else:\n
  \                   debug_value = md.options[\"article\"].get(\n                        \"path\",\n
  \                       md.options[\"article\"].get(\n                            \"title\",
  md.options[\"article\"].get(\"slug\", \"\")\n                        ),\n                    )\n
  \               logger.warning(\n                    f\"wikilink [[{text}]] has
  duplicate matches ({possible_pages}) in file '{debug_value}', defaulting to the
  first match ({possible_pages[0]})\",\n                )\n                link =
  possible_pages[0]\n            else:\n                if md.options[\"article\"]
  is None:\n                    debug_value = \"UNKNOWN\"\n                else:\n
  \                   debug_value = md.options[\"article\"].get(\n                        \"path\",\n
  \                       md.options[\"article\"].get(\n                            \"title\",
  md.options[\"article\"].get(\"slug\", \"\")\n                        ),\n                    )\n
  \               logger.warning(\n                    f\"wikilink [[{text}]] no matches
  in file '{debug_value}', defaulting to '/{text}'\",\n                )\n                link
  = text\n\n            if id and not link.endswith(f\"#{id}\"):\n                link
  = f\"{link}#{id}\"\n\n            token.attrSet(\"href\", f\"/{link}\")\n            content_token
  = state.push(\"text\", \"\", 0)\n            content_token.content = text\n\n            token
  = state.push(\"link_close\", \"a\", -1)\n            token.content = text\n\n            return
  True\n\n        md.inline.ruler.before(\"escape\", \"wikilinks_inline\", _wikilinks_inline)\n
  \   ```"
date: 2025-02-14
description: 'The `markata.plugins.md_it_wikilinks` plugin adds support for wiki-style
  links using

  double brackets (`[[link]]`). It automatically resolves links to other post'
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>md_it_wikilinks.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"The `markata.plugins.md_it_wikilinks`
    plugin adds support for wiki-style links using\ndouble brackets (`[[link]]`).
    It automatically resolves links to other post\" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n
    \       <meta property=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    <head>\n        <script src=\"https://unpkg.com/@tailwindcss/browser@4\"></script>\n<title>md_it_wikilinks.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"The `markata.plugins.md_it_wikilinks`
    plugin adds support for wiki-style links using\ndouble brackets (`[[link]]`).
    It automatically resolves links to other post\" />\n <link href=\"/favicon.ico\"
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        md_it_wikilinks.py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <hr />\n<p>The <code>markata.plugins.md_it_wikilinks</code>
    plugin adds support for wiki-style links using\ndouble brackets (<code>[[link]]</code>).
    It automatically resolves links to other posts in your\nsite using file names
    or slugs.</p>\n<h1 id=\"installation\">Installation <a class=\"header-anchor\"
    href=\"#installation\"><svg class=\"heading-permalink\" aria-hidden=\"true\" fill=\"currentColor\"
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
    class=\"p\">[</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;markata.plugins.md_it_wikilinks&quot;</span><span
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
    class=\"p\">[</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;markata.plugins.md_it_wikilinks&quot;</span><span
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>This plugin requires
    no explicit configuration. It automatically processes wikilinks\nin your markdown
    content.</p>\n<h1 id=\"functionality\">Functionality <a class=\"header-anchor\"
    href=\"#functionality\"><svg class=\"heading-permalink\" aria-hidden=\"true\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<h2 id=\"basic-wikilinks\">Basic
    Wikilinks <a class=\"header-anchor\" href=\"#basic-wikilinks\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Simple file-based linking:</p>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span>[[nav]]              -&gt;
    links to docs/nav.md as /nav\n[[blog/post]]        -&gt; links to blog/post.md
    as /blog/post\n[[about|About Me]]   -&gt; links to about.md with &quot;About Me&quot;
    as text\n</pre></div>\n\n</pre>\n\n<h2 id=\"smart-slug-resolution\">Smart Slug
    Resolution <a class=\"header-anchor\" href=\"#smart-slug-resolution\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>The plugin:</p>\n<ol>\n<li>Looks
    up the target file in your content</li>\n<li>Finds its generated slug</li>\n<li>Creates
    a link to the final URL</li>\n</ol>\n<p>Example:</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"gh\"># File:
    posts/2024-01-my-post.md</span>\nslug: /blog/my-post\n\n<span class=\"gh\"># In
    another file:</span>\n[[2024-01-my-post]]  -&gt; links to /blog/my-post\n</pre></div>\n\n</pre>\n\n<h2
    id=\"link-formats\">Link Formats <a class=\"header-anchor\" href=\"#link-formats\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Supports multiple link
    styles:</p>\n<ul>\n<li>Basic: <code>[[filename]]</code></li>\n<li>With text: <code>[[filename|Link
    Text]]</code></li>\n<li>With path: <code>[[folder/file]]</code></li>\n<li>With
    extension: <code>[[file.md]]</code> (extension stripped in output)</li>\n</ul>\n<h2
    id=\"html-output\">HTML Output <a class=\"header-anchor\" href=\"#html-output\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Generated HTML structure:</p>\n<pre
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
    class=\"s\">&quot;/target-slug&quot;</span><span class=\"p\">&gt;</span>Link Text<span
    class=\"p\">&lt;/</span><span class=\"nt\">a</span><span class=\"p\">&gt;</span>\n</pre></div>\n\n</pre>\n\n<h2
    id=\"error-handling\">Error Handling <a class=\"header-anchor\" href=\"#error-handling\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>For broken links:</p>\n<ul>\n<li>Maintains
    the wikilink syntax</li>\n<li>Adds a 'broken-link' class</li>\n<li>Optionally
    logs warnings</li>\n</ul>\n<h2 id=\"dependencies\">Dependencies <a class=\"header-anchor\"
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
    on:</p>\n<ul>\n<li>markdown-it-py for markdown parsing</li>\n<li>The <code>render_markdown</code>
    plugin for final HTML output</li>\n</ul>\n<hr />\n<div class=\"admonition function\">\n<p
    class=\"admonition-title\">Function</p>\n<h2 id=\"wikilinks_plugin\" class=\"admonition-title\"
    style=\"margin: 0; padding: .5rem 1rem;\">wikilinks_plugin <em class=\"small\">function</em></h2>\n<p>A
    plugin to create wikilinks tokens.\nThese, token should be handled by the renderer.</p>\n<div
    class=\"admonition example is-collapsible collapsible-open\">\n<p class=\"admonition-title\">Example</p>\n<pre
    class='wrapper'>\n\n<div class='filepath'>\n<p>markdown</p>\n\n<div class='right'>\n\n<button
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n</div>\n\n<div
    class=\"highlight\"><pre><span></span>[[nav]]\n</pre></div>\n\n</pre>\n\n<pre
    class='wrapper'>\n\n<div class='filepath'>\n<p>html</p>\n\n<div class='right'>\n\n<button
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n</div>\n\n<div
    class=\"highlight\"><pre><span></span><span class=\"p\">&lt;</span><span class=\"nt\">a</span>
    <span class=\"na\">class</span><span class=\"o\">=</span><span class=\"s\">&quot;wikilink&quot;</span>
    <span class=\"na\">href</span><span class=\"o\">=</span><span class=\"s\">&quot;/nav&quot;</span><span
    class=\"p\">&gt;</span>load<span class=\"p\">&lt;/</span><span class=\"nt\">a</span><span
    class=\"p\">&gt;</span>\n</pre></div>\n\n</pre>\n\n</div>\n</div>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">wikilinks_plugin
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
    class=\"w\"> </span><span class=\"nf\">wikilinks_plugin</span><span class=\"p\">(</span>\n
    \   <span class=\"n\">md</span><span class=\"p\">:</span> <span class=\"n\">MarkdownIt</span><span
    class=\"p\">,</span>\n    <span class=\"n\">start_delimiter</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;[&quot;</span><span
    class=\"p\">,</span>\n    <span class=\"n\">end_delimiter</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;]&quot;</span><span
    class=\"p\">,</span>\n    <span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"kc\">None</span><span class=\"p\">,</span>\n<span class=\"p\">):</span>\n<span
    class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;A plugin to create
    wikilinks tokens.</span>\n<span class=\"sd\">    These, token should be handled
    by the renderer.</span>\n\n<span class=\"sd\">    ???+ example</span>\n\n<span
    class=\"sd\">        ```md title=markdown</span>\n<span class=\"sd\">        [[nav]]</span>\n<span
    class=\"sd\">        ```</span>\n\n<span class=\"sd\">        ```html title=html</span>\n<span
    class=\"sd\">        &lt;a class=&quot;wikilink&quot; href=&quot;/nav&quot;&gt;load&lt;/a&gt;</span>\n<span
    class=\"sd\">        ```</span>\n<span class=\"sd\">    &quot;&quot;&quot;</span>\n\n
    \   <span class=\"n\">start_char</span> <span class=\"o\">=</span> <span class=\"nb\">ord</span><span
    class=\"p\">(</span><span class=\"n\">start_delimiter</span><span class=\"p\">)</span>\n
    \   <span class=\"n\">end_char</span> <span class=\"o\">=</span> <span class=\"nb\">ord</span><span
    class=\"p\">(</span><span class=\"n\">end_delimiter</span><span class=\"p\">)</span>\n\n
    \   <span class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">_wikilinks_inline</span><span
    class=\"p\">(</span><span class=\"n\">state</span><span class=\"p\">:</span> <span
    class=\"n\">StateInline</span><span class=\"p\">,</span> <span class=\"n\">silent</span><span
    class=\"p\">:</span> <span class=\"nb\">bool</span><span class=\"p\">):</span>\n
    \       <span class=\"k\">try</span><span class=\"p\">:</span>\n            <span
    class=\"k\">if</span> <span class=\"p\">(</span>\n                <span class=\"n\">state</span><span
    class=\"o\">.</span><span class=\"n\">srcCharCode</span><span class=\"p\">[</span><span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">pos</span><span
    class=\"p\">]</span> <span class=\"o\">!=</span> <span class=\"n\">start_char</span>\n
    \               <span class=\"ow\">or</span> <span class=\"n\">state</span><span
    class=\"o\">.</span><span class=\"n\">srcCharCode</span><span class=\"p\">[</span><span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">pos</span>
    <span class=\"o\">+</span> <span class=\"mi\">1</span><span class=\"p\">]</span>
    <span class=\"o\">!=</span> <span class=\"n\">start_char</span>\n            <span
    class=\"p\">):</span>\n                <span class=\"k\">return</span> <span class=\"kc\">False</span>\n
    \       <span class=\"k\">except</span> <span class=\"ne\">IndexError</span><span
    class=\"p\">:</span>\n            <span class=\"k\">return</span> <span class=\"kc\">False</span>\n\n
    \       <span class=\"n\">pos</span> <span class=\"o\">=</span> <span class=\"n\">state</span><span
    class=\"o\">.</span><span class=\"n\">pos</span> <span class=\"o\">+</span> <span
    class=\"mi\">2</span>\n        <span class=\"n\">found_closing</span> <span class=\"o\">=</span>
    <span class=\"kc\">False</span>\n        <span class=\"k\">while</span> <span
    class=\"kc\">True</span><span class=\"p\">:</span>\n            <span class=\"k\">try</span><span
    class=\"p\">:</span>\n                <span class=\"n\">end</span> <span class=\"o\">=</span>
    <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">srcCharCode</span><span
    class=\"o\">.</span><span class=\"n\">index</span><span class=\"p\">(</span><span
    class=\"n\">end_char</span><span class=\"p\">,</span> <span class=\"n\">pos</span><span
    class=\"p\">)</span>\n            <span class=\"k\">except</span> <span class=\"ne\">ValueError</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"kc\">False</span>\n
    \           <span class=\"k\">try</span><span class=\"p\">:</span>\n                <span
    class=\"k\">if</span> <span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">srcCharCode</span><span class=\"p\">[</span><span class=\"n\">end</span>
    <span class=\"o\">+</span> <span class=\"mi\">1</span><span class=\"p\">]</span>
    <span class=\"o\">==</span> <span class=\"n\">end_char</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">found_closing</span> <span class=\"o\">=</span>
    <span class=\"kc\">True</span>\n                    <span class=\"k\">break</span>\n
    \           <span class=\"k\">except</span> <span class=\"ne\">IndexError</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"kc\">False</span>\n
    \           <span class=\"n\">pos</span> <span class=\"o\">=</span> <span class=\"n\">end</span>
    <span class=\"o\">+</span> <span class=\"mi\">2</span>\n\n        <span class=\"k\">if</span>
    <span class=\"ow\">not</span> <span class=\"n\">found_closing</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">return</span> <span class=\"kc\">False</span>\n\n
    \       <span class=\"n\">text</span> <span class=\"o\">=</span> <span class=\"n\">state</span><span
    class=\"o\">.</span><span class=\"n\">src</span><span class=\"p\">[</span><span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">pos</span>
    <span class=\"o\">+</span> <span class=\"mi\">2</span> <span class=\"p\">:</span>
    <span class=\"n\">end</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">strip</span><span class=\"p\">()</span>\n        <span class=\"n\">state</span><span
    class=\"o\">.</span><span class=\"n\">pos</span> <span class=\"o\">=</span> <span
    class=\"n\">end</span> <span class=\"o\">+</span> <span class=\"mi\">2</span>\n\n
    \       <span class=\"k\">if</span> <span class=\"n\">silent</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">return</span> <span class=\"kc\">True</span>\n\n
    \       <span class=\"n\">token</span> <span class=\"o\">=</span> <span class=\"n\">state</span><span
    class=\"o\">.</span><span class=\"n\">push</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;link_open&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;a&quot;</span><span
    class=\"p\">,</span> <span class=\"mi\">1</span><span class=\"p\">)</span>\n        <span
    class=\"n\">token</span><span class=\"o\">.</span><span class=\"n\">block</span>
    <span class=\"o\">=</span> <span class=\"kc\">False</span>\n        <span class=\"n\">token</span><span
    class=\"o\">.</span><span class=\"n\">attrSet</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;class&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;wikilink&quot;</span><span
    class=\"p\">)</span>\n        <span class=\"k\">if</span> <span class=\"s2\">&quot;#&quot;</span>
    <span class=\"ow\">in</span> <span class=\"n\">text</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">link</span><span class=\"p\">,</span> <span class=\"nb\">id</span>
    <span class=\"o\">=</span> <span class=\"n\">text</span><span class=\"o\">.</span><span
    class=\"n\">split</span><span class=\"p\">(</span><span class=\"s2\">&quot;#&quot;</span><span
    class=\"p\">)</span>\n            <span class=\"n\">link</span> <span class=\"o\">=</span>
    <span class=\"n\">link</span><span class=\"o\">.</span><span class=\"n\">strip</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;/&quot;</span><span class=\"p\">)</span>\n
    \       <span class=\"k\">else</span><span class=\"p\">:</span>\n            <span
    class=\"n\">link</span><span class=\"p\">,</span> <span class=\"nb\">id</span>
    <span class=\"o\">=</span> <span class=\"n\">text</span><span class=\"p\">,</span>
    <span class=\"kc\">None</span>\n\n        <span class=\"c1\"># possible_pages
    = markata.filter(</span>\n        <span class=\"c1\">#     f&#39;str(path).split(&quot;/&quot;)[-1].split(&quot;.&quot;)[0].replace(&quot;_&quot;,
    &quot;-&quot;) == &quot;{link.replace(&quot;_&quot;, &quot;-&quot;)}&quot;&#39;,</span>\n
    \       <span class=\"c1\"># )</span>\n        <span class=\"n\">possible_pages</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">possible_wikilink</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"n\">link</span><span class=\"p\">,</span> <span
    class=\"p\">[])</span>\n        <span class=\"k\">if</span> <span class=\"nb\">len</span><span
    class=\"p\">(</span><span class=\"n\">possible_pages</span><span class=\"p\">)</span>
    <span class=\"o\">==</span> <span class=\"mi\">1</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">link</span> <span class=\"o\">=</span> <span class=\"n\">possible_pages</span><span
    class=\"p\">[</span><span class=\"mi\">0</span><span class=\"p\">]</span>\n        <span
    class=\"k\">elif</span> <span class=\"nb\">len</span><span class=\"p\">(</span><span
    class=\"n\">possible_pages</span><span class=\"p\">)</span> <span class=\"o\">&gt;</span>
    <span class=\"mi\">1</span><span class=\"p\">:</span>\n            <span class=\"k\">if</span>
    <span class=\"n\">md</span><span class=\"o\">.</span><span class=\"n\">options</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;article&quot;</span><span class=\"p\">]</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">debug_value</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;UNKNOWN&quot;</span>\n            <span class=\"k\">else</span><span
    class=\"p\">:</span>\n                <span class=\"n\">debug_value</span> <span
    class=\"o\">=</span> <span class=\"n\">md</span><span class=\"o\">.</span><span
    class=\"n\">options</span><span class=\"p\">[</span><span class=\"s2\">&quot;article&quot;</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span>\n                    <span class=\"s2\">&quot;path&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"n\">md</span><span class=\"o\">.</span><span
    class=\"n\">options</span><span class=\"p\">[</span><span class=\"s2\">&quot;article&quot;</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span>\n                        <span class=\"s2\">&quot;title&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">md</span><span class=\"o\">.</span><span
    class=\"n\">options</span><span class=\"p\">[</span><span class=\"s2\">&quot;article&quot;</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"p\">),</span>\n                <span class=\"p\">)</span>\n            <span
    class=\"n\">logger</span><span class=\"o\">.</span><span class=\"n\">warning</span><span
    class=\"p\">(</span>\n                <span class=\"sa\">f</span><span class=\"s2\">&quot;wikilink
    [[</span><span class=\"si\">{</span><span class=\"n\">text</span><span class=\"si\">}</span><span
    class=\"s2\">]] has duplicate matches (</span><span class=\"si\">{</span><span
    class=\"n\">possible_pages</span><span class=\"si\">}</span><span class=\"s2\">)
    in file &#39;</span><span class=\"si\">{</span><span class=\"n\">debug_value</span><span
    class=\"si\">}</span><span class=\"s2\">&#39;, defaulting to the first match (</span><span
    class=\"si\">{</span><span class=\"n\">possible_pages</span><span class=\"p\">[</span><span
    class=\"mi\">0</span><span class=\"p\">]</span><span class=\"si\">}</span><span
    class=\"s2\">)&quot;</span><span class=\"p\">,</span>\n            <span class=\"p\">)</span>\n
    \           <span class=\"n\">link</span> <span class=\"o\">=</span> <span class=\"n\">possible_pages</span><span
    class=\"p\">[</span><span class=\"mi\">0</span><span class=\"p\">]</span>\n        <span
    class=\"k\">else</span><span class=\"p\">:</span>\n            <span class=\"k\">if</span>
    <span class=\"n\">md</span><span class=\"o\">.</span><span class=\"n\">options</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;article&quot;</span><span class=\"p\">]</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">debug_value</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;UNKNOWN&quot;</span>\n            <span class=\"k\">else</span><span
    class=\"p\">:</span>\n                <span class=\"n\">debug_value</span> <span
    class=\"o\">=</span> <span class=\"n\">md</span><span class=\"o\">.</span><span
    class=\"n\">options</span><span class=\"p\">[</span><span class=\"s2\">&quot;article&quot;</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span>\n                    <span class=\"s2\">&quot;path&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"n\">md</span><span class=\"o\">.</span><span
    class=\"n\">options</span><span class=\"p\">[</span><span class=\"s2\">&quot;article&quot;</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span>\n                        <span class=\"s2\">&quot;title&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">md</span><span class=\"o\">.</span><span
    class=\"n\">options</span><span class=\"p\">[</span><span class=\"s2\">&quot;article&quot;</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"p\">),</span>\n                <span class=\"p\">)</span>\n            <span
    class=\"n\">logger</span><span class=\"o\">.</span><span class=\"n\">warning</span><span
    class=\"p\">(</span>\n                <span class=\"sa\">f</span><span class=\"s2\">&quot;wikilink
    [[</span><span class=\"si\">{</span><span class=\"n\">text</span><span class=\"si\">}</span><span
    class=\"s2\">]] no matches in file &#39;</span><span class=\"si\">{</span><span
    class=\"n\">debug_value</span><span class=\"si\">}</span><span class=\"s2\">&#39;,
    defaulting to &#39;/</span><span class=\"si\">{</span><span class=\"n\">text</span><span
    class=\"si\">}</span><span class=\"s2\">&#39;&quot;</span><span class=\"p\">,</span>\n
    \           <span class=\"p\">)</span>\n            <span class=\"n\">link</span>
    <span class=\"o\">=</span> <span class=\"n\">text</span>\n\n        <span class=\"k\">if</span>
    <span class=\"nb\">id</span> <span class=\"ow\">and</span> <span class=\"ow\">not</span>
    <span class=\"n\">link</span><span class=\"o\">.</span><span class=\"n\">endswith</span><span
    class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">&quot;#</span><span
    class=\"si\">{</span><span class=\"nb\">id</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span><span class=\"p\">):</span>\n            <span class=\"n\">link</span>
    <span class=\"o\">=</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">link</span><span class=\"si\">}</span><span
    class=\"s2\">#</span><span class=\"si\">{</span><span class=\"nb\">id</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span>\n\n        <span class=\"n\">token</span><span
    class=\"o\">.</span><span class=\"n\">attrSet</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;href&quot;</span><span class=\"p\">,</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;/</span><span class=\"si\">{</span><span class=\"n\">link</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">)</span>\n
    \       <span class=\"n\">content_token</span> <span class=\"o\">=</span> <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">push</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;text&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;&quot;</span><span class=\"p\">,</span> <span class=\"mi\">0</span><span
    class=\"p\">)</span>\n        <span class=\"n\">content_token</span><span class=\"o\">.</span><span
    class=\"n\">content</span> <span class=\"o\">=</span> <span class=\"n\">text</span>\n\n
    \       <span class=\"n\">token</span> <span class=\"o\">=</span> <span class=\"n\">state</span><span
    class=\"o\">.</span><span class=\"n\">push</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;link_close&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;a&quot;</span><span
    class=\"p\">,</span> <span class=\"o\">-</span><span class=\"mi\">1</span><span
    class=\"p\">)</span>\n        <span class=\"n\">token</span><span class=\"o\">.</span><span
    class=\"n\">content</span> <span class=\"o\">=</span> <span class=\"n\">text</span>\n\n
    \       <span class=\"k\">return</span> <span class=\"kc\">True</span>\n\n    <span
    class=\"n\">md</span><span class=\"o\">.</span><span class=\"n\">inline</span><span
    class=\"o\">.</span><span class=\"n\">ruler</span><span class=\"o\">.</span><span
    class=\"n\">before</span><span class=\"p\">(</span><span class=\"s2\">&quot;escape&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;wikilinks_inline&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">_wikilinks_inline</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n\n
    \   </section>\n</article>    </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>md_it_wikilinks.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"The `markata.plugins.md_it_wikilinks`
    plugin adds support for wiki-style links using\ndouble brackets (`[[link]]`).
    It automatically resolves links to other post\" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n
    \       <meta property=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    <head>\n        <script src=\"https://unpkg.com/@tailwindcss/browser@4\"></script>\n<title>md_it_wikilinks.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"The `markata.plugins.md_it_wikilinks`
    plugin adds support for wiki-style links using\ndouble brackets (`[[link]]`).
    It automatically resolves links to other post\" />\n <link href=\"/favicon.ico\"
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
    class=\"title\">\n    <h1 id=\"title\">\n        md_it_wikilinks.py\n    </h1>\n</section></article>\n
    \    </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       md_it_wikilinks.py\n    </h1>\n</section>    <section class=\"body\">\n
    \       <hr />\n<p>The <code>markata.plugins.md_it_wikilinks</code> plugin adds
    support for wiki-style links using\ndouble brackets (<code>[[link]]</code>). It
    automatically resolves links to other posts in your\nsite using file names or
    slugs.</p>\n<h1 id=\"installation\">Installation <a class=\"header-anchor\" href=\"#installation\"><svg
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
    class=\"p\">[</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;markata.plugins.md_it_wikilinks&quot;</span><span
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
    class=\"p\">[</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;markata.plugins.md_it_wikilinks&quot;</span><span
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>This plugin requires
    no explicit configuration. It automatically processes wikilinks\nin your markdown
    content.</p>\n<h1 id=\"functionality\">Functionality <a class=\"header-anchor\"
    href=\"#functionality\"><svg class=\"heading-permalink\" aria-hidden=\"true\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<h2 id=\"basic-wikilinks\">Basic
    Wikilinks <a class=\"header-anchor\" href=\"#basic-wikilinks\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Simple file-based linking:</p>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span>[[nav]]              -&gt;
    links to docs/nav.md as /nav\n[[blog/post]]        -&gt; links to blog/post.md
    as /blog/post\n[[about|About Me]]   -&gt; links to about.md with &quot;About Me&quot;
    as text\n</pre></div>\n\n</pre>\n\n<h2 id=\"smart-slug-resolution\">Smart Slug
    Resolution <a class=\"header-anchor\" href=\"#smart-slug-resolution\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>The plugin:</p>\n<ol>\n<li>Looks
    up the target file in your content</li>\n<li>Finds its generated slug</li>\n<li>Creates
    a link to the final URL</li>\n</ol>\n<p>Example:</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"gh\"># File:
    posts/2024-01-my-post.md</span>\nslug: /blog/my-post\n\n<span class=\"gh\"># In
    another file:</span>\n[[2024-01-my-post]]  -&gt; links to /blog/my-post\n</pre></div>\n\n</pre>\n\n<h2
    id=\"link-formats\">Link Formats <a class=\"header-anchor\" href=\"#link-formats\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Supports multiple link
    styles:</p>\n<ul>\n<li>Basic: <code>[[filename]]</code></li>\n<li>With text: <code>[[filename|Link
    Text]]</code></li>\n<li>With path: <code>[[folder/file]]</code></li>\n<li>With
    extension: <code>[[file.md]]</code> (extension stripped in output)</li>\n</ul>\n<h2
    id=\"html-output\">HTML Output <a class=\"header-anchor\" href=\"#html-output\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Generated HTML structure:</p>\n<pre
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
    class=\"s\">&quot;/target-slug&quot;</span><span class=\"p\">&gt;</span>Link Text<span
    class=\"p\">&lt;/</span><span class=\"nt\">a</span><span class=\"p\">&gt;</span>\n</pre></div>\n\n</pre>\n\n<h2
    id=\"error-handling\">Error Handling <a class=\"header-anchor\" href=\"#error-handling\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>For broken links:</p>\n<ul>\n<li>Maintains
    the wikilink syntax</li>\n<li>Adds a 'broken-link' class</li>\n<li>Optionally
    logs warnings</li>\n</ul>\n<h2 id=\"dependencies\">Dependencies <a class=\"header-anchor\"
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
    on:</p>\n<ul>\n<li>markdown-it-py for markdown parsing</li>\n<li>The <code>render_markdown</code>
    plugin for final HTML output</li>\n</ul>\n<hr />\n<div class=\"admonition function\">\n<p
    class=\"admonition-title\">Function</p>\n<h2 id=\"wikilinks_plugin\" class=\"admonition-title\"
    style=\"margin: 0; padding: .5rem 1rem;\">wikilinks_plugin <em class=\"small\">function</em></h2>\n<p>A
    plugin to create wikilinks tokens.\nThese, token should be handled by the renderer.</p>\n<div
    class=\"admonition example is-collapsible collapsible-open\">\n<p class=\"admonition-title\">Example</p>\n<pre
    class='wrapper'>\n\n<div class='filepath'>\n<p>markdown</p>\n\n<div class='right'>\n\n<button
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n</div>\n\n<div
    class=\"highlight\"><pre><span></span>[[nav]]\n</pre></div>\n\n</pre>\n\n<pre
    class='wrapper'>\n\n<div class='filepath'>\n<p>html</p>\n\n<div class='right'>\n\n<button
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n</div>\n\n<div
    class=\"highlight\"><pre><span></span><span class=\"p\">&lt;</span><span class=\"nt\">a</span>
    <span class=\"na\">class</span><span class=\"o\">=</span><span class=\"s\">&quot;wikilink&quot;</span>
    <span class=\"na\">href</span><span class=\"o\">=</span><span class=\"s\">&quot;/nav&quot;</span><span
    class=\"p\">&gt;</span>load<span class=\"p\">&lt;/</span><span class=\"nt\">a</span><span
    class=\"p\">&gt;</span>\n</pre></div>\n\n</pre>\n\n</div>\n</div>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">wikilinks_plugin
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
    class=\"w\"> </span><span class=\"nf\">wikilinks_plugin</span><span class=\"p\">(</span>\n
    \   <span class=\"n\">md</span><span class=\"p\">:</span> <span class=\"n\">MarkdownIt</span><span
    class=\"p\">,</span>\n    <span class=\"n\">start_delimiter</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;[&quot;</span><span
    class=\"p\">,</span>\n    <span class=\"n\">end_delimiter</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;]&quot;</span><span
    class=\"p\">,</span>\n    <span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"kc\">None</span><span class=\"p\">,</span>\n<span class=\"p\">):</span>\n<span
    class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;A plugin to create
    wikilinks tokens.</span>\n<span class=\"sd\">    These, token should be handled
    by the renderer.</span>\n\n<span class=\"sd\">    ???+ example</span>\n\n<span
    class=\"sd\">        ```md title=markdown</span>\n<span class=\"sd\">        [[nav]]</span>\n<span
    class=\"sd\">        ```</span>\n\n<span class=\"sd\">        ```html title=html</span>\n<span
    class=\"sd\">        &lt;a class=&quot;wikilink&quot; href=&quot;/nav&quot;&gt;load&lt;/a&gt;</span>\n<span
    class=\"sd\">        ```</span>\n<span class=\"sd\">    &quot;&quot;&quot;</span>\n\n
    \   <span class=\"n\">start_char</span> <span class=\"o\">=</span> <span class=\"nb\">ord</span><span
    class=\"p\">(</span><span class=\"n\">start_delimiter</span><span class=\"p\">)</span>\n
    \   <span class=\"n\">end_char</span> <span class=\"o\">=</span> <span class=\"nb\">ord</span><span
    class=\"p\">(</span><span class=\"n\">end_delimiter</span><span class=\"p\">)</span>\n\n
    \   <span class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">_wikilinks_inline</span><span
    class=\"p\">(</span><span class=\"n\">state</span><span class=\"p\">:</span> <span
    class=\"n\">StateInline</span><span class=\"p\">,</span> <span class=\"n\">silent</span><span
    class=\"p\">:</span> <span class=\"nb\">bool</span><span class=\"p\">):</span>\n
    \       <span class=\"k\">try</span><span class=\"p\">:</span>\n            <span
    class=\"k\">if</span> <span class=\"p\">(</span>\n                <span class=\"n\">state</span><span
    class=\"o\">.</span><span class=\"n\">srcCharCode</span><span class=\"p\">[</span><span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">pos</span><span
    class=\"p\">]</span> <span class=\"o\">!=</span> <span class=\"n\">start_char</span>\n
    \               <span class=\"ow\">or</span> <span class=\"n\">state</span><span
    class=\"o\">.</span><span class=\"n\">srcCharCode</span><span class=\"p\">[</span><span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">pos</span>
    <span class=\"o\">+</span> <span class=\"mi\">1</span><span class=\"p\">]</span>
    <span class=\"o\">!=</span> <span class=\"n\">start_char</span>\n            <span
    class=\"p\">):</span>\n                <span class=\"k\">return</span> <span class=\"kc\">False</span>\n
    \       <span class=\"k\">except</span> <span class=\"ne\">IndexError</span><span
    class=\"p\">:</span>\n            <span class=\"k\">return</span> <span class=\"kc\">False</span>\n\n
    \       <span class=\"n\">pos</span> <span class=\"o\">=</span> <span class=\"n\">state</span><span
    class=\"o\">.</span><span class=\"n\">pos</span> <span class=\"o\">+</span> <span
    class=\"mi\">2</span>\n        <span class=\"n\">found_closing</span> <span class=\"o\">=</span>
    <span class=\"kc\">False</span>\n        <span class=\"k\">while</span> <span
    class=\"kc\">True</span><span class=\"p\">:</span>\n            <span class=\"k\">try</span><span
    class=\"p\">:</span>\n                <span class=\"n\">end</span> <span class=\"o\">=</span>
    <span class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">srcCharCode</span><span
    class=\"o\">.</span><span class=\"n\">index</span><span class=\"p\">(</span><span
    class=\"n\">end_char</span><span class=\"p\">,</span> <span class=\"n\">pos</span><span
    class=\"p\">)</span>\n            <span class=\"k\">except</span> <span class=\"ne\">ValueError</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"kc\">False</span>\n
    \           <span class=\"k\">try</span><span class=\"p\">:</span>\n                <span
    class=\"k\">if</span> <span class=\"n\">state</span><span class=\"o\">.</span><span
    class=\"n\">srcCharCode</span><span class=\"p\">[</span><span class=\"n\">end</span>
    <span class=\"o\">+</span> <span class=\"mi\">1</span><span class=\"p\">]</span>
    <span class=\"o\">==</span> <span class=\"n\">end_char</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">found_closing</span> <span class=\"o\">=</span>
    <span class=\"kc\">True</span>\n                    <span class=\"k\">break</span>\n
    \           <span class=\"k\">except</span> <span class=\"ne\">IndexError</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"kc\">False</span>\n
    \           <span class=\"n\">pos</span> <span class=\"o\">=</span> <span class=\"n\">end</span>
    <span class=\"o\">+</span> <span class=\"mi\">2</span>\n\n        <span class=\"k\">if</span>
    <span class=\"ow\">not</span> <span class=\"n\">found_closing</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">return</span> <span class=\"kc\">False</span>\n\n
    \       <span class=\"n\">text</span> <span class=\"o\">=</span> <span class=\"n\">state</span><span
    class=\"o\">.</span><span class=\"n\">src</span><span class=\"p\">[</span><span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">pos</span>
    <span class=\"o\">+</span> <span class=\"mi\">2</span> <span class=\"p\">:</span>
    <span class=\"n\">end</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">strip</span><span class=\"p\">()</span>\n        <span class=\"n\">state</span><span
    class=\"o\">.</span><span class=\"n\">pos</span> <span class=\"o\">=</span> <span
    class=\"n\">end</span> <span class=\"o\">+</span> <span class=\"mi\">2</span>\n\n
    \       <span class=\"k\">if</span> <span class=\"n\">silent</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">return</span> <span class=\"kc\">True</span>\n\n
    \       <span class=\"n\">token</span> <span class=\"o\">=</span> <span class=\"n\">state</span><span
    class=\"o\">.</span><span class=\"n\">push</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;link_open&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;a&quot;</span><span
    class=\"p\">,</span> <span class=\"mi\">1</span><span class=\"p\">)</span>\n        <span
    class=\"n\">token</span><span class=\"o\">.</span><span class=\"n\">block</span>
    <span class=\"o\">=</span> <span class=\"kc\">False</span>\n        <span class=\"n\">token</span><span
    class=\"o\">.</span><span class=\"n\">attrSet</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;class&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;wikilink&quot;</span><span
    class=\"p\">)</span>\n        <span class=\"k\">if</span> <span class=\"s2\">&quot;#&quot;</span>
    <span class=\"ow\">in</span> <span class=\"n\">text</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">link</span><span class=\"p\">,</span> <span class=\"nb\">id</span>
    <span class=\"o\">=</span> <span class=\"n\">text</span><span class=\"o\">.</span><span
    class=\"n\">split</span><span class=\"p\">(</span><span class=\"s2\">&quot;#&quot;</span><span
    class=\"p\">)</span>\n            <span class=\"n\">link</span> <span class=\"o\">=</span>
    <span class=\"n\">link</span><span class=\"o\">.</span><span class=\"n\">strip</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;/&quot;</span><span class=\"p\">)</span>\n
    \       <span class=\"k\">else</span><span class=\"p\">:</span>\n            <span
    class=\"n\">link</span><span class=\"p\">,</span> <span class=\"nb\">id</span>
    <span class=\"o\">=</span> <span class=\"n\">text</span><span class=\"p\">,</span>
    <span class=\"kc\">None</span>\n\n        <span class=\"c1\"># possible_pages
    = markata.filter(</span>\n        <span class=\"c1\">#     f&#39;str(path).split(&quot;/&quot;)[-1].split(&quot;.&quot;)[0].replace(&quot;_&quot;,
    &quot;-&quot;) == &quot;{link.replace(&quot;_&quot;, &quot;-&quot;)}&quot;&#39;,</span>\n
    \       <span class=\"c1\"># )</span>\n        <span class=\"n\">possible_pages</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">possible_wikilink</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"n\">link</span><span class=\"p\">,</span> <span
    class=\"p\">[])</span>\n        <span class=\"k\">if</span> <span class=\"nb\">len</span><span
    class=\"p\">(</span><span class=\"n\">possible_pages</span><span class=\"p\">)</span>
    <span class=\"o\">==</span> <span class=\"mi\">1</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">link</span> <span class=\"o\">=</span> <span class=\"n\">possible_pages</span><span
    class=\"p\">[</span><span class=\"mi\">0</span><span class=\"p\">]</span>\n        <span
    class=\"k\">elif</span> <span class=\"nb\">len</span><span class=\"p\">(</span><span
    class=\"n\">possible_pages</span><span class=\"p\">)</span> <span class=\"o\">&gt;</span>
    <span class=\"mi\">1</span><span class=\"p\">:</span>\n            <span class=\"k\">if</span>
    <span class=\"n\">md</span><span class=\"o\">.</span><span class=\"n\">options</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;article&quot;</span><span class=\"p\">]</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">debug_value</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;UNKNOWN&quot;</span>\n            <span class=\"k\">else</span><span
    class=\"p\">:</span>\n                <span class=\"n\">debug_value</span> <span
    class=\"o\">=</span> <span class=\"n\">md</span><span class=\"o\">.</span><span
    class=\"n\">options</span><span class=\"p\">[</span><span class=\"s2\">&quot;article&quot;</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span>\n                    <span class=\"s2\">&quot;path&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"n\">md</span><span class=\"o\">.</span><span
    class=\"n\">options</span><span class=\"p\">[</span><span class=\"s2\">&quot;article&quot;</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span>\n                        <span class=\"s2\">&quot;title&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">md</span><span class=\"o\">.</span><span
    class=\"n\">options</span><span class=\"p\">[</span><span class=\"s2\">&quot;article&quot;</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"p\">),</span>\n                <span class=\"p\">)</span>\n            <span
    class=\"n\">logger</span><span class=\"o\">.</span><span class=\"n\">warning</span><span
    class=\"p\">(</span>\n                <span class=\"sa\">f</span><span class=\"s2\">&quot;wikilink
    [[</span><span class=\"si\">{</span><span class=\"n\">text</span><span class=\"si\">}</span><span
    class=\"s2\">]] has duplicate matches (</span><span class=\"si\">{</span><span
    class=\"n\">possible_pages</span><span class=\"si\">}</span><span class=\"s2\">)
    in file &#39;</span><span class=\"si\">{</span><span class=\"n\">debug_value</span><span
    class=\"si\">}</span><span class=\"s2\">&#39;, defaulting to the first match (</span><span
    class=\"si\">{</span><span class=\"n\">possible_pages</span><span class=\"p\">[</span><span
    class=\"mi\">0</span><span class=\"p\">]</span><span class=\"si\">}</span><span
    class=\"s2\">)&quot;</span><span class=\"p\">,</span>\n            <span class=\"p\">)</span>\n
    \           <span class=\"n\">link</span> <span class=\"o\">=</span> <span class=\"n\">possible_pages</span><span
    class=\"p\">[</span><span class=\"mi\">0</span><span class=\"p\">]</span>\n        <span
    class=\"k\">else</span><span class=\"p\">:</span>\n            <span class=\"k\">if</span>
    <span class=\"n\">md</span><span class=\"o\">.</span><span class=\"n\">options</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;article&quot;</span><span class=\"p\">]</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">debug_value</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;UNKNOWN&quot;</span>\n            <span class=\"k\">else</span><span
    class=\"p\">:</span>\n                <span class=\"n\">debug_value</span> <span
    class=\"o\">=</span> <span class=\"n\">md</span><span class=\"o\">.</span><span
    class=\"n\">options</span><span class=\"p\">[</span><span class=\"s2\">&quot;article&quot;</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span>\n                    <span class=\"s2\">&quot;path&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"n\">md</span><span class=\"o\">.</span><span
    class=\"n\">options</span><span class=\"p\">[</span><span class=\"s2\">&quot;article&quot;</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span>\n                        <span class=\"s2\">&quot;title&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">md</span><span class=\"o\">.</span><span
    class=\"n\">options</span><span class=\"p\">[</span><span class=\"s2\">&quot;article&quot;</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"p\">),</span>\n                <span class=\"p\">)</span>\n            <span
    class=\"n\">logger</span><span class=\"o\">.</span><span class=\"n\">warning</span><span
    class=\"p\">(</span>\n                <span class=\"sa\">f</span><span class=\"s2\">&quot;wikilink
    [[</span><span class=\"si\">{</span><span class=\"n\">text</span><span class=\"si\">}</span><span
    class=\"s2\">]] no matches in file &#39;</span><span class=\"si\">{</span><span
    class=\"n\">debug_value</span><span class=\"si\">}</span><span class=\"s2\">&#39;,
    defaulting to &#39;/</span><span class=\"si\">{</span><span class=\"n\">text</span><span
    class=\"si\">}</span><span class=\"s2\">&#39;&quot;</span><span class=\"p\">,</span>\n
    \           <span class=\"p\">)</span>\n            <span class=\"n\">link</span>
    <span class=\"o\">=</span> <span class=\"n\">text</span>\n\n        <span class=\"k\">if</span>
    <span class=\"nb\">id</span> <span class=\"ow\">and</span> <span class=\"ow\">not</span>
    <span class=\"n\">link</span><span class=\"o\">.</span><span class=\"n\">endswith</span><span
    class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">&quot;#</span><span
    class=\"si\">{</span><span class=\"nb\">id</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span><span class=\"p\">):</span>\n            <span class=\"n\">link</span>
    <span class=\"o\">=</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">link</span><span class=\"si\">}</span><span
    class=\"s2\">#</span><span class=\"si\">{</span><span class=\"nb\">id</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span>\n\n        <span class=\"n\">token</span><span
    class=\"o\">.</span><span class=\"n\">attrSet</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;href&quot;</span><span class=\"p\">,</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;/</span><span class=\"si\">{</span><span class=\"n\">link</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">)</span>\n
    \       <span class=\"n\">content_token</span> <span class=\"o\">=</span> <span
    class=\"n\">state</span><span class=\"o\">.</span><span class=\"n\">push</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;text&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;&quot;</span><span class=\"p\">,</span> <span class=\"mi\">0</span><span
    class=\"p\">)</span>\n        <span class=\"n\">content_token</span><span class=\"o\">.</span><span
    class=\"n\">content</span> <span class=\"o\">=</span> <span class=\"n\">text</span>\n\n
    \       <span class=\"n\">token</span> <span class=\"o\">=</span> <span class=\"n\">state</span><span
    class=\"o\">.</span><span class=\"n\">push</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;link_close&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;a&quot;</span><span
    class=\"p\">,</span> <span class=\"o\">-</span><span class=\"mi\">1</span><span
    class=\"p\">)</span>\n        <span class=\"n\">token</span><span class=\"o\">.</span><span
    class=\"n\">content</span> <span class=\"o\">=</span> <span class=\"n\">text</span>\n\n
    \       <span class=\"k\">return</span> <span class=\"kc\">True</span>\n\n    <span
    class=\"n\">md</span><span class=\"o\">.</span><span class=\"n\">inline</span><span
    class=\"o\">.</span><span class=\"n\">ruler</span><span class=\"o\">.</span><span
    class=\"n\">before</span><span class=\"p\">(</span><span class=\"s2\">&quot;escape&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;wikilinks_inline&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">_wikilinks_inline</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n\n
    \   </section>\n</article>"
  raw.md: ''
published: false
slug: markata/plugins/md-it-wikilinks
title: md_it_wikilinks.py


---

---

The `markata.plugins.md_it_wikilinks` plugin adds support for wiki-style links using
double brackets (`[[link]]`). It automatically resolves links to other posts in your
site using file names or slugs.

# Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.md_it_wikilinks",
]
```

# Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.md_it_wikilinks",
]
```

# Configuration

This plugin requires no explicit configuration. It automatically processes wikilinks
in your markdown content.

# Functionality

## Basic Wikilinks

Simple file-based linking:
```markdown
[[nav]]              -> links to docs/nav.md as /nav
[[blog/post]]        -> links to blog/post.md as /blog/post
[[about|About Me]]   -> links to about.md with "About Me" as text
```

## Smart Slug Resolution

The plugin:
1. Looks up the target file in your content
2. Finds its generated slug
3. Creates a link to the final URL

Example:
```markdown
# File: posts/2024-01-my-post.md
slug: /blog/my-post

# In another file:
[[2024-01-my-post]]  -> links to /blog/my-post
```

## Link Formats

Supports multiple link styles:
- Basic: `[[filename]]`
- With text: `[[filename|Link Text]]`
- With path: `[[folder/file]]`
- With extension: `[[file.md]]` (extension stripped in output)

## HTML Output

Generated HTML structure:
```html
<a class="wikilink" href="/target-slug">Link Text</a>
```

## Error Handling

For broken links:
- Maintains the wikilink syntax
- Adds a 'broken-link' class
- Optionally logs warnings

## Dependencies

This plugin depends on:
- markdown-it-py for markdown parsing
- The `render_markdown` plugin for final HTML output

---





!!! function
    <h2 id="wikilinks_plugin" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">wikilinks_plugin <em class="small">function</em></h2>

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

            # possible_pages = markata.filter(
            #     f'str(path).split("/")[-1].split(".")[0].replace("_", "-") == "{link.replace("_", "-")}"',
            # )
            possible_pages = markata.possible_wikilink.get(link, [])
            if len(possible_pages) == 1:
                link = possible_pages[0]
            elif len(possible_pages) > 1:
                if md.options["article"] is None:
                    debug_value = "UNKNOWN"
                else:
                    debug_value = md.options["article"].get(
                        "path",
                        md.options["article"].get(
                            "title", md.options["article"].get("slug", "")
                        ),
                    )
                logger.warning(
                    f"wikilink [[{text}]] has duplicate matches ({possible_pages}) in file '{debug_value}', defaulting to the first match ({possible_pages[0]})",
                )
                link = possible_pages[0]
            else:
                if md.options["article"] is None:
                    debug_value = "UNKNOWN"
                else:
                    debug_value = md.options["article"].get(
                        "path",
                        md.options["article"].get(
                            "title", md.options["article"].get("slug", "")
                        ),
                    )
                logger.warning(
                    f"wikilink [[{text}]] no matches in file '{debug_value}', defaulting to '/{text}'",
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