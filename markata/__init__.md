---
content: "---\n\nMarkata is a tool for handling directories of markdown.\n\n---\n\n!!!
  method\n    <h2 id=\"teardown\" class=\"admonition-title\" style=\"margin: 0; padding:
  .5rem 1rem;\">teardown <em class=\"small\">method</em></h2>\n\n    Cleanup and print
  statistics when Markata is done.\n\n???+ source \"teardown <em class='small'>source</em>\"\n
  \   ```python\n    def teardown(self: \"Markata\"):\n            \"\"\"Cleanup and
  print statistics when Markata is done.\"\"\"\n            # Print map cache statistics
  if they exist\n            if hasattr(self, \"_map_cache_stats\"):\n                stats
  = self._map_cache_stats\n                total = stats[\"total\"]\n                if
  total > 0:\n                    hit_rate = (stats[\"hits\"] / total) * 100\n                    self.console.print(\"\\n[yellow]Map
  Cache Statistics:[/yellow]\")\n                    self.console.print(f\"Total calls:
  {total}\")\n                    self.console.print(f\"Cache hits: {stats['hits']}\")\n
  \                   self.console.print(f\"Cache misses: {stats['misses']}\")\n                    self.console.print(f\"Hit
  rate: {hit_rate:.1f}%\")\n                    self.console.print(\n                        f\"Cache
  size: {len(getattr(self, '_filtered_cache', {}))}\"\n                    )\n            if
  self.stages_ran:\n                self._pm.hook.teardown(markata=self)\n            return
  self\n    ```\n!!! method\n    <h2 id=\"_compile_sort_key\" class=\"admonition-title\"
  style=\"margin: 0; padding: .5rem 1rem;\">_compile_sort_key <em class=\"small\">method</em></h2>\n\n
  \   Compile a sort key function for better performance\n\n???+ source \"_compile_sort_key
  <em class='small'>source</em>\"\n    ```python\n    def _compile_sort_key(self,
  sort: str):\n            \"\"\"Compile a sort key function for better performance\"\"\"\n
  \           if \"datetime\" in sort.lower():\n                return lambda a: a.get(sort,
  datetime.datetime(1970, 1, 1))\n            if \"date\" in sort.lower():\n                return
  lambda a: a.get(sort, datetime.date(1970, 1, 1))\n\n            # Create a compiled
  function for complex sort expressions\n            try:\n                code =
  compile(sort, \"<string>\", \"eval\")\n\n                def sort_key(a):\n                    try:\n
  \                       value = eval(code, a.to_dict(), {})\n                        if
  isinstance(value, (int, float)):\n                            return value\n                        if
  hasattr(value, \"timestamp\"):\n                            return value.timestamp()\n
  \                       if isinstance(value, datetime.date):\n                            return
  datetime.datetime.combine(\n                                value,\n                                datetime.datetime.min.time(),\n
  \                           ).timestamp()\n                        return sum(ord(c)
  for c in str(value))\n                    except Exception:\n                        return
  -1\n\n                return sort_key\n            except Exception:\n                return
  lambda _: -1\n    ```\n!!! method\n    <h2 id=\"_get_sort_key\" class=\"admonition-title\"
  style=\"margin: 0; padding: .5rem 1rem;\">_get_sort_key <em class=\"small\">method</em></h2>\n\n
  \   Cache compiled sort key functions\n\n???+ source \"_get_sort_key <em class='small'>source</em>\"\n
  \   ```python\n    def _get_sort_key(self, sort: str):\n            \"\"\"Cache
  compiled sort key functions\"\"\"\n            return self._compile_sort_key(sort)\n
  \   ```\n!!! method\n    <h2 id=\"_get_eval_globals\" class=\"admonition-title\"
  style=\"margin: 0; padding: .5rem 1rem;\">_get_eval_globals <em class=\"small\">method</em></h2>\n\n
  \   Get common globals used in eval operations\n\n???+ source \"_get_eval_globals
  <em class='small'>source</em>\"\n    ```python\n    def _get_eval_globals(self):\n
  \           \"\"\"Get common globals used in eval operations\"\"\"\n            if
  not hasattr(self, \"_eval_globals\"):\n                self._eval_globals = {\"timedelta\":
  timedelta}\n            return self._eval_globals\n    ```\n!!! method\n    <h2
  id=\"_eval_with_article\" class=\"admonition-title\" style=\"margin: 0; padding:
  .5rem 1rem;\">_eval_with_article <em class=\"small\">method</em></h2>\n\n    Evaluate
  code with article context, reusing dict where possible\n\n???+ source \"_eval_with_article
  <em class='small'>source</em>\"\n    ```python\n    def _eval_with_article(self,
  code, article, extra_globals=None):\n            \"\"\"Evaluate code with article
  context, reusing dict where possible\"\"\"\n            if not hasattr(article,
  \"_eval_dict\"):\n                article._eval_dict = article.to_dict()\n                article._eval_dict.update({\"post\":
  article, \"m\": self})\n\n            globals_dict = self._get_eval_globals()\n
  \           if extra_globals:\n                globals_dict.update(extra_globals)\n\n
  \           try:\n                return eval(code, article._eval_dict, globals_dict)\n
  \           except Exception:\n                return None\n    ```"
date: 2025-05-05
description: "Markata is a tool for handling directories of markdown. !!! method teardown
  \ method ???+ source \"teardown  source \" !!! method _compile_sort_key  method
  ???+\u2026"
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>__init__.py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"Markata is a tool for handling directories of markdown.
    !!! method teardown  method ???+ source \"teardown  source \" !!! method _compile_sort_key
    \ method ???+\u2026\" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n        <meta property=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    <head>\n        <script src=\"https://unpkg.com/@tailwindcss/browser@4\"></script>\n<title>__init__.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Markata is a tool for handling directories
    of markdown. !!! method teardown  method ???+ source \"teardown  source \" !!!
    method _compile_sort_key  method ???+\u2026\" />\n <link href=\"/favicon.ico\"
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        __init__.py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <hr />\n<p>Markata
    is a tool for handling directories of markdown.</p>\n<hr />\n<div class=\"admonition
    method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"teardown\" class=\"admonition-title\"
    style=\"margin: 0; padding: .5rem 1rem;\">teardown <em class=\"small\">method</em></h2>\n<p>Cleanup
    and print statistics when Markata is done.</p>\n</div>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">teardown
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
    class=\"w\"> </span><span class=\"nf\">teardown</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">):</span>\n<span class=\"w\">        </span><span class=\"sd\">&quot;&quot;&quot;Cleanup
    and print statistics when Markata is done.&quot;&quot;&quot;</span>\n        <span
    class=\"c1\"># Print map cache statistics if they exist</span>\n        <span
    class=\"k\">if</span> <span class=\"nb\">hasattr</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"s2\">&quot;_map_cache_stats&quot;</span><span
    class=\"p\">):</span>\n            <span class=\"n\">stats</span> <span class=\"o\">=</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_map_cache_stats</span>\n
    \           <span class=\"n\">total</span> <span class=\"o\">=</span> <span class=\"n\">stats</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;total&quot;</span><span class=\"p\">]</span>\n
    \           <span class=\"k\">if</span> <span class=\"n\">total</span> <span class=\"o\">&gt;</span>
    <span class=\"mi\">0</span><span class=\"p\">:</span>\n                <span class=\"n\">hit_rate</span>
    <span class=\"o\">=</span> <span class=\"p\">(</span><span class=\"n\">stats</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;hits&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">/</span> <span class=\"n\">total</span><span class=\"p\">)</span>
    <span class=\"o\">*</span> <span class=\"mi\">100</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">print</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span class=\"s2\">[yellow]Map
    Cache Statistics:[/yellow]&quot;</span><span class=\"p\">)</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">print</span><span class=\"p\">(</span><span
    class=\"sa\">f</span><span class=\"s2\">&quot;Total calls: </span><span class=\"si\">{</span><span
    class=\"n\">total</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">print</span><span
    class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">&quot;Cache
    hits: </span><span class=\"si\">{</span><span class=\"n\">stats</span><span class=\"p\">[</span><span
    class=\"s1\">&#39;hits&#39;</span><span class=\"p\">]</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span><span class=\"p\">)</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">print</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;Cache misses: </span><span class=\"si\">{</span><span class=\"n\">stats</span><span
    class=\"p\">[</span><span class=\"s1\">&#39;misses&#39;</span><span class=\"p\">]</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">print</span><span
    class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">&quot;Hit rate:
    </span><span class=\"si\">{</span><span class=\"n\">hit_rate</span><span class=\"si\">:</span><span
    class=\"s2\">.1f</span><span class=\"si\">}</span><span class=\"s2\">%&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">print</span><span
    class=\"p\">(</span>\n                    <span class=\"sa\">f</span><span class=\"s2\">&quot;Cache
    size: </span><span class=\"si\">{</span><span class=\"nb\">len</span><span class=\"p\">(</span><span
    class=\"nb\">getattr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span><span class=\"w\"> </span><span class=\"s1\">&#39;_filtered_cache&#39;</span><span
    class=\"p\">,</span><span class=\"w\"> </span><span class=\"p\">{}))</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span>\n                <span class=\"p\">)</span>\n
    \       <span class=\"k\">if</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">stages_ran</span><span class=\"p\">:</span>\n            <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_pm</span><span class=\"o\">.</span><span
    class=\"n\">hook</span><span class=\"o\">.</span><span class=\"n\">teardown</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"bp\">self</span><span class=\"p\">)</span>\n        <span class=\"k\">return</span>
    <span class=\"bp\">self</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div class=\"admonition
    method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"_compile_sort_key\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">_compile_sort_key
    <em class=\"small\">method</em></h2>\n<p>Compile a sort key function for better
    performance</p>\n</div>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">_compile_sort_key <em class='small'>source</em></p>\n<pre
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
    class=\"w\"> </span><span class=\"nf\">_compile_sort_key</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"n\">sort</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">):</span>\n<span
    class=\"w\">        </span><span class=\"sd\">&quot;&quot;&quot;Compile a sort
    key function for better performance&quot;&quot;&quot;</span>\n        <span class=\"k\">if</span>
    <span class=\"s2\">&quot;datetime&quot;</span> <span class=\"ow\">in</span> <span
    class=\"n\">sort</span><span class=\"o\">.</span><span class=\"n\">lower</span><span
    class=\"p\">():</span>\n            <span class=\"k\">return</span> <span class=\"k\">lambda</span>
    <span class=\"n\">a</span><span class=\"p\">:</span> <span class=\"n\">a</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"n\">sort</span><span class=\"p\">,</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"p\">(</span><span
    class=\"mi\">1970</span><span class=\"p\">,</span> <span class=\"mi\">1</span><span
    class=\"p\">,</span> <span class=\"mi\">1</span><span class=\"p\">))</span>\n
    \       <span class=\"k\">if</span> <span class=\"s2\">&quot;date&quot;</span>
    <span class=\"ow\">in</span> <span class=\"n\">sort</span><span class=\"o\">.</span><span
    class=\"n\">lower</span><span class=\"p\">():</span>\n            <span class=\"k\">return</span>
    <span class=\"k\">lambda</span> <span class=\"n\">a</span><span class=\"p\">:</span>
    <span class=\"n\">a</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"n\">sort</span><span class=\"p\">,</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">date</span><span
    class=\"p\">(</span><span class=\"mi\">1970</span><span class=\"p\">,</span> <span
    class=\"mi\">1</span><span class=\"p\">,</span> <span class=\"mi\">1</span><span
    class=\"p\">))</span>\n\n        <span class=\"c1\"># Create a compiled function
    for complex sort expressions</span>\n        <span class=\"k\">try</span><span
    class=\"p\">:</span>\n            <span class=\"n\">code</span> <span class=\"o\">=</span>
    <span class=\"nb\">compile</span><span class=\"p\">(</span><span class=\"n\">sort</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;&lt;string&gt;&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;eval&quot;</span><span class=\"p\">)</span>\n\n
    \           <span class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">sort_key</span><span
    class=\"p\">(</span><span class=\"n\">a</span><span class=\"p\">):</span>\n                <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                    <span class=\"n\">value</span>
    <span class=\"o\">=</span> <span class=\"nb\">eval</span><span class=\"p\">(</span><span
    class=\"n\">code</span><span class=\"p\">,</span> <span class=\"n\">a</span><span
    class=\"o\">.</span><span class=\"n\">to_dict</span><span class=\"p\">(),</span>
    <span class=\"p\">{})</span>\n                    <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">value</span><span
    class=\"p\">,</span> <span class=\"p\">(</span><span class=\"nb\">int</span><span
    class=\"p\">,</span> <span class=\"nb\">float</span><span class=\"p\">)):</span>\n
    \                       <span class=\"k\">return</span> <span class=\"n\">value</span>\n
    \                   <span class=\"k\">if</span> <span class=\"nb\">hasattr</span><span
    class=\"p\">(</span><span class=\"n\">value</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;timestamp&quot;</span><span class=\"p\">):</span>\n                        <span
    class=\"k\">return</span> <span class=\"n\">value</span><span class=\"o\">.</span><span
    class=\"n\">timestamp</span><span class=\"p\">()</span>\n                    <span
    class=\"k\">if</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">value</span><span class=\"p\">,</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">date</span><span class=\"p\">):</span>\n
    \                       <span class=\"k\">return</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">combine</span><span class=\"p\">(</span>\n                            <span
    class=\"n\">value</span><span class=\"p\">,</span>\n                            <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">min</span><span class=\"o\">.</span><span
    class=\"n\">time</span><span class=\"p\">(),</span>\n                        <span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">timestamp</span><span
    class=\"p\">()</span>\n                    <span class=\"k\">return</span> <span
    class=\"nb\">sum</span><span class=\"p\">(</span><span class=\"nb\">ord</span><span
    class=\"p\">(</span><span class=\"n\">c</span><span class=\"p\">)</span> <span
    class=\"k\">for</span> <span class=\"n\">c</span> <span class=\"ow\">in</span>
    <span class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">value</span><span
    class=\"p\">))</span>\n                <span class=\"k\">except</span> <span class=\"ne\">Exception</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
    class=\"o\">-</span><span class=\"mi\">1</span>\n\n            <span class=\"k\">return</span>
    <span class=\"n\">sort_key</span>\n        <span class=\"k\">except</span> <span
    class=\"ne\">Exception</span><span class=\"p\">:</span>\n            <span class=\"k\">return</span>
    <span class=\"k\">lambda</span> <span class=\"n\">_</span><span class=\"p\">:</span>
    <span class=\"o\">-</span><span class=\"mi\">1</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"_get_sort_key\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">_get_sort_key
    <em class=\"small\">method</em></h2>\n<p>Cache compiled sort key functions</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_get_sort_key
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
    class=\"w\"> </span><span class=\"nf\">_get_sort_key</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"n\">sort</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">):</span>\n<span
    class=\"w\">        </span><span class=\"sd\">&quot;&quot;&quot;Cache compiled
    sort key functions&quot;&quot;&quot;</span>\n        <span class=\"k\">return</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_compile_sort_key</span><span
    class=\"p\">(</span><span class=\"n\">sort</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"_get_eval_globals\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">_get_eval_globals
    <em class=\"small\">method</em></h2>\n<p>Get common globals used in eval operations</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_get_eval_globals
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
    class=\"w\"> </span><span class=\"nf\">_get_eval_globals</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">):</span>\n<span class=\"w\">        </span><span
    class=\"sd\">&quot;&quot;&quot;Get common globals used in eval operations&quot;&quot;&quot;</span>\n
    \       <span class=\"k\">if</span> <span class=\"ow\">not</span> <span class=\"nb\">hasattr</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;_eval_globals&quot;</span><span class=\"p\">):</span>\n            <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_eval_globals</span>
    <span class=\"o\">=</span> <span class=\"p\">{</span><span class=\"s2\">&quot;timedelta&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">timedelta</span><span class=\"p\">}</span>\n
    \       <span class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_eval_globals</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div class=\"admonition
    method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"_eval_with_article\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">_eval_with_article
    <em class=\"small\">method</em></h2>\n<p>Evaluate code with article context, reusing
    dict where possible</p>\n</div>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">_eval_with_article <em class='small'>source</em></p>\n<pre
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
    class=\"w\"> </span><span class=\"nf\">_eval_with_article</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"n\">code</span><span
    class=\"p\">,</span> <span class=\"n\">article</span><span class=\"p\">,</span>
    <span class=\"n\">extra_globals</span><span class=\"o\">=</span><span class=\"kc\">None</span><span
    class=\"p\">):</span>\n<span class=\"w\">        </span><span class=\"sd\">&quot;&quot;&quot;Evaluate
    code with article context, reusing dict where possible&quot;&quot;&quot;</span>\n
    \       <span class=\"k\">if</span> <span class=\"ow\">not</span> <span class=\"nb\">hasattr</span><span
    class=\"p\">(</span><span class=\"n\">article</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;_eval_dict&quot;</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">_eval_dict</span>
    <span class=\"o\">=</span> <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">to_dict</span><span class=\"p\">()</span>\n            <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">_eval_dict</span><span class=\"o\">.</span><span
    class=\"n\">update</span><span class=\"p\">({</span><span class=\"s2\">&quot;post&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">article</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;m&quot;</span><span class=\"p\">:</span> <span class=\"bp\">self</span><span
    class=\"p\">})</span>\n\n        <span class=\"n\">globals_dict</span> <span class=\"o\">=</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_get_eval_globals</span><span
    class=\"p\">()</span>\n        <span class=\"k\">if</span> <span class=\"n\">extra_globals</span><span
    class=\"p\">:</span>\n            <span class=\"n\">globals_dict</span><span class=\"o\">.</span><span
    class=\"n\">update</span><span class=\"p\">(</span><span class=\"n\">extra_globals</span><span
    class=\"p\">)</span>\n\n        <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">return</span> <span class=\"nb\">eval</span><span
    class=\"p\">(</span><span class=\"n\">code</span><span class=\"p\">,</span> <span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">_eval_dict</span><span
    class=\"p\">,</span> <span class=\"n\">globals_dict</span><span class=\"p\">)</span>\n
    \       <span class=\"k\">except</span> <span class=\"ne\">Exception</span><span
    class=\"p\">:</span>\n            <span class=\"k\">return</span> <span class=\"kc\">None</span>\n</pre></div>\n\n</pre>\n\n</div>\n\n
    \   </section>\n</article>    </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>__init__.py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"Markata is a tool for handling directories of markdown.
    !!! method teardown  method ???+ source \"teardown  source \" !!! method _compile_sort_key
    \ method ???+\u2026\" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n        <meta property=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    <head>\n        <script src=\"https://unpkg.com/@tailwindcss/browser@4\"></script>\n<title>__init__.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Markata is a tool for handling directories
    of markdown. !!! method teardown  method ???+ source \"teardown  source \" !!!
    method _compile_sort_key  method ???+\u2026\" />\n <link href=\"/favicon.ico\"
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
    class=\"title\">\n    <h1 id=\"title\">\n        __init__.py\n    </h1>\n</section></article>\n
    \    </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       __init__.py\n    </h1>\n</section>    <section class=\"body\">\n        <hr
    />\n<p>Markata is a tool for handling directories of markdown.</p>\n<hr />\n<div
    class=\"admonition method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"teardown\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">teardown
    <em class=\"small\">method</em></h2>\n<p>Cleanup and print statistics when Markata
    is done.</p>\n</div>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">teardown <em class='small'>source</em></p>\n<pre class='wrapper'>\n\n<div
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
    class=\"w\"> </span><span class=\"nf\">teardown</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">):</span>\n<span class=\"w\">        </span><span class=\"sd\">&quot;&quot;&quot;Cleanup
    and print statistics when Markata is done.&quot;&quot;&quot;</span>\n        <span
    class=\"c1\"># Print map cache statistics if they exist</span>\n        <span
    class=\"k\">if</span> <span class=\"nb\">hasattr</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"s2\">&quot;_map_cache_stats&quot;</span><span
    class=\"p\">):</span>\n            <span class=\"n\">stats</span> <span class=\"o\">=</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_map_cache_stats</span>\n
    \           <span class=\"n\">total</span> <span class=\"o\">=</span> <span class=\"n\">stats</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;total&quot;</span><span class=\"p\">]</span>\n
    \           <span class=\"k\">if</span> <span class=\"n\">total</span> <span class=\"o\">&gt;</span>
    <span class=\"mi\">0</span><span class=\"p\">:</span>\n                <span class=\"n\">hit_rate</span>
    <span class=\"o\">=</span> <span class=\"p\">(</span><span class=\"n\">stats</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;hits&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">/</span> <span class=\"n\">total</span><span class=\"p\">)</span>
    <span class=\"o\">*</span> <span class=\"mi\">100</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">print</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span class=\"s2\">[yellow]Map
    Cache Statistics:[/yellow]&quot;</span><span class=\"p\">)</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">print</span><span class=\"p\">(</span><span
    class=\"sa\">f</span><span class=\"s2\">&quot;Total calls: </span><span class=\"si\">{</span><span
    class=\"n\">total</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">print</span><span
    class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">&quot;Cache
    hits: </span><span class=\"si\">{</span><span class=\"n\">stats</span><span class=\"p\">[</span><span
    class=\"s1\">&#39;hits&#39;</span><span class=\"p\">]</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span><span class=\"p\">)</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">print</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;Cache misses: </span><span class=\"si\">{</span><span class=\"n\">stats</span><span
    class=\"p\">[</span><span class=\"s1\">&#39;misses&#39;</span><span class=\"p\">]</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">print</span><span
    class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">&quot;Hit rate:
    </span><span class=\"si\">{</span><span class=\"n\">hit_rate</span><span class=\"si\">:</span><span
    class=\"s2\">.1f</span><span class=\"si\">}</span><span class=\"s2\">%&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">print</span><span
    class=\"p\">(</span>\n                    <span class=\"sa\">f</span><span class=\"s2\">&quot;Cache
    size: </span><span class=\"si\">{</span><span class=\"nb\">len</span><span class=\"p\">(</span><span
    class=\"nb\">getattr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span><span class=\"w\"> </span><span class=\"s1\">&#39;_filtered_cache&#39;</span><span
    class=\"p\">,</span><span class=\"w\"> </span><span class=\"p\">{}))</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span>\n                <span class=\"p\">)</span>\n
    \       <span class=\"k\">if</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">stages_ran</span><span class=\"p\">:</span>\n            <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_pm</span><span class=\"o\">.</span><span
    class=\"n\">hook</span><span class=\"o\">.</span><span class=\"n\">teardown</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"bp\">self</span><span class=\"p\">)</span>\n        <span class=\"k\">return</span>
    <span class=\"bp\">self</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div class=\"admonition
    method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"_compile_sort_key\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">_compile_sort_key
    <em class=\"small\">method</em></h2>\n<p>Compile a sort key function for better
    performance</p>\n</div>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">_compile_sort_key <em class='small'>source</em></p>\n<pre
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
    class=\"w\"> </span><span class=\"nf\">_compile_sort_key</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"n\">sort</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">):</span>\n<span
    class=\"w\">        </span><span class=\"sd\">&quot;&quot;&quot;Compile a sort
    key function for better performance&quot;&quot;&quot;</span>\n        <span class=\"k\">if</span>
    <span class=\"s2\">&quot;datetime&quot;</span> <span class=\"ow\">in</span> <span
    class=\"n\">sort</span><span class=\"o\">.</span><span class=\"n\">lower</span><span
    class=\"p\">():</span>\n            <span class=\"k\">return</span> <span class=\"k\">lambda</span>
    <span class=\"n\">a</span><span class=\"p\">:</span> <span class=\"n\">a</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"n\">sort</span><span class=\"p\">,</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"p\">(</span><span
    class=\"mi\">1970</span><span class=\"p\">,</span> <span class=\"mi\">1</span><span
    class=\"p\">,</span> <span class=\"mi\">1</span><span class=\"p\">))</span>\n
    \       <span class=\"k\">if</span> <span class=\"s2\">&quot;date&quot;</span>
    <span class=\"ow\">in</span> <span class=\"n\">sort</span><span class=\"o\">.</span><span
    class=\"n\">lower</span><span class=\"p\">():</span>\n            <span class=\"k\">return</span>
    <span class=\"k\">lambda</span> <span class=\"n\">a</span><span class=\"p\">:</span>
    <span class=\"n\">a</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"n\">sort</span><span class=\"p\">,</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">date</span><span
    class=\"p\">(</span><span class=\"mi\">1970</span><span class=\"p\">,</span> <span
    class=\"mi\">1</span><span class=\"p\">,</span> <span class=\"mi\">1</span><span
    class=\"p\">))</span>\n\n        <span class=\"c1\"># Create a compiled function
    for complex sort expressions</span>\n        <span class=\"k\">try</span><span
    class=\"p\">:</span>\n            <span class=\"n\">code</span> <span class=\"o\">=</span>
    <span class=\"nb\">compile</span><span class=\"p\">(</span><span class=\"n\">sort</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;&lt;string&gt;&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;eval&quot;</span><span class=\"p\">)</span>\n\n
    \           <span class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">sort_key</span><span
    class=\"p\">(</span><span class=\"n\">a</span><span class=\"p\">):</span>\n                <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                    <span class=\"n\">value</span>
    <span class=\"o\">=</span> <span class=\"nb\">eval</span><span class=\"p\">(</span><span
    class=\"n\">code</span><span class=\"p\">,</span> <span class=\"n\">a</span><span
    class=\"o\">.</span><span class=\"n\">to_dict</span><span class=\"p\">(),</span>
    <span class=\"p\">{})</span>\n                    <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">value</span><span
    class=\"p\">,</span> <span class=\"p\">(</span><span class=\"nb\">int</span><span
    class=\"p\">,</span> <span class=\"nb\">float</span><span class=\"p\">)):</span>\n
    \                       <span class=\"k\">return</span> <span class=\"n\">value</span>\n
    \                   <span class=\"k\">if</span> <span class=\"nb\">hasattr</span><span
    class=\"p\">(</span><span class=\"n\">value</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;timestamp&quot;</span><span class=\"p\">):</span>\n                        <span
    class=\"k\">return</span> <span class=\"n\">value</span><span class=\"o\">.</span><span
    class=\"n\">timestamp</span><span class=\"p\">()</span>\n                    <span
    class=\"k\">if</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">value</span><span class=\"p\">,</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">date</span><span class=\"p\">):</span>\n
    \                       <span class=\"k\">return</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">combine</span><span class=\"p\">(</span>\n                            <span
    class=\"n\">value</span><span class=\"p\">,</span>\n                            <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">min</span><span class=\"o\">.</span><span
    class=\"n\">time</span><span class=\"p\">(),</span>\n                        <span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">timestamp</span><span
    class=\"p\">()</span>\n                    <span class=\"k\">return</span> <span
    class=\"nb\">sum</span><span class=\"p\">(</span><span class=\"nb\">ord</span><span
    class=\"p\">(</span><span class=\"n\">c</span><span class=\"p\">)</span> <span
    class=\"k\">for</span> <span class=\"n\">c</span> <span class=\"ow\">in</span>
    <span class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">value</span><span
    class=\"p\">))</span>\n                <span class=\"k\">except</span> <span class=\"ne\">Exception</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
    class=\"o\">-</span><span class=\"mi\">1</span>\n\n            <span class=\"k\">return</span>
    <span class=\"n\">sort_key</span>\n        <span class=\"k\">except</span> <span
    class=\"ne\">Exception</span><span class=\"p\">:</span>\n            <span class=\"k\">return</span>
    <span class=\"k\">lambda</span> <span class=\"n\">_</span><span class=\"p\">:</span>
    <span class=\"o\">-</span><span class=\"mi\">1</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"_get_sort_key\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">_get_sort_key
    <em class=\"small\">method</em></h2>\n<p>Cache compiled sort key functions</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_get_sort_key
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
    class=\"w\"> </span><span class=\"nf\">_get_sort_key</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"n\">sort</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">):</span>\n<span
    class=\"w\">        </span><span class=\"sd\">&quot;&quot;&quot;Cache compiled
    sort key functions&quot;&quot;&quot;</span>\n        <span class=\"k\">return</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_compile_sort_key</span><span
    class=\"p\">(</span><span class=\"n\">sort</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"_get_eval_globals\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">_get_eval_globals
    <em class=\"small\">method</em></h2>\n<p>Get common globals used in eval operations</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_get_eval_globals
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
    class=\"w\"> </span><span class=\"nf\">_get_eval_globals</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">):</span>\n<span class=\"w\">        </span><span
    class=\"sd\">&quot;&quot;&quot;Get common globals used in eval operations&quot;&quot;&quot;</span>\n
    \       <span class=\"k\">if</span> <span class=\"ow\">not</span> <span class=\"nb\">hasattr</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;_eval_globals&quot;</span><span class=\"p\">):</span>\n            <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_eval_globals</span>
    <span class=\"o\">=</span> <span class=\"p\">{</span><span class=\"s2\">&quot;timedelta&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">timedelta</span><span class=\"p\">}</span>\n
    \       <span class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_eval_globals</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div class=\"admonition
    method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"_eval_with_article\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">_eval_with_article
    <em class=\"small\">method</em></h2>\n<p>Evaluate code with article context, reusing
    dict where possible</p>\n</div>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">_eval_with_article <em class='small'>source</em></p>\n<pre
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
    class=\"w\"> </span><span class=\"nf\">_eval_with_article</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"n\">code</span><span
    class=\"p\">,</span> <span class=\"n\">article</span><span class=\"p\">,</span>
    <span class=\"n\">extra_globals</span><span class=\"o\">=</span><span class=\"kc\">None</span><span
    class=\"p\">):</span>\n<span class=\"w\">        </span><span class=\"sd\">&quot;&quot;&quot;Evaluate
    code with article context, reusing dict where possible&quot;&quot;&quot;</span>\n
    \       <span class=\"k\">if</span> <span class=\"ow\">not</span> <span class=\"nb\">hasattr</span><span
    class=\"p\">(</span><span class=\"n\">article</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;_eval_dict&quot;</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">_eval_dict</span>
    <span class=\"o\">=</span> <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">to_dict</span><span class=\"p\">()</span>\n            <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">_eval_dict</span><span class=\"o\">.</span><span
    class=\"n\">update</span><span class=\"p\">({</span><span class=\"s2\">&quot;post&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">article</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;m&quot;</span><span class=\"p\">:</span> <span class=\"bp\">self</span><span
    class=\"p\">})</span>\n\n        <span class=\"n\">globals_dict</span> <span class=\"o\">=</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_get_eval_globals</span><span
    class=\"p\">()</span>\n        <span class=\"k\">if</span> <span class=\"n\">extra_globals</span><span
    class=\"p\">:</span>\n            <span class=\"n\">globals_dict</span><span class=\"o\">.</span><span
    class=\"n\">update</span><span class=\"p\">(</span><span class=\"n\">extra_globals</span><span
    class=\"p\">)</span>\n\n        <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">return</span> <span class=\"nb\">eval</span><span
    class=\"p\">(</span><span class=\"n\">code</span><span class=\"p\">,</span> <span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">_eval_dict</span><span
    class=\"p\">,</span> <span class=\"n\">globals_dict</span><span class=\"p\">)</span>\n
    \       <span class=\"k\">except</span> <span class=\"ne\">Exception</span><span
    class=\"p\">:</span>\n            <span class=\"k\">return</span> <span class=\"kc\">None</span>\n</pre></div>\n\n</pre>\n\n</div>\n\n
    \   </section>\n</article>"
  raw.md: ''
published: false
slug: markata/init
title: __init__.py


---

---

Markata is a tool for handling directories of markdown.

---

!!! method
    <h2 id="teardown" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">teardown <em class="small">method</em></h2>

    Cleanup and print statistics when Markata is done.

???+ source "teardown <em class='small'>source</em>"
    ```python
    def teardown(self: "Markata"):
            """Cleanup and print statistics when Markata is done."""
            # Print map cache statistics if they exist
            if hasattr(self, "_map_cache_stats"):
                stats = self._map_cache_stats
                total = stats["total"]
                if total > 0:
                    hit_rate = (stats["hits"] / total) * 100
                    self.console.print("\n[yellow]Map Cache Statistics:[/yellow]")
                    self.console.print(f"Total calls: {total}")
                    self.console.print(f"Cache hits: {stats['hits']}")
                    self.console.print(f"Cache misses: {stats['misses']}")
                    self.console.print(f"Hit rate: {hit_rate:.1f}%")
                    self.console.print(
                        f"Cache size: {len(getattr(self, '_filtered_cache', {}))}"
                    )
            if self.stages_ran:
                self._pm.hook.teardown(markata=self)
            return self
    ```
!!! method
    <h2 id="_compile_sort_key" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">_compile_sort_key <em class="small">method</em></h2>

    Compile a sort key function for better performance

???+ source "_compile_sort_key <em class='small'>source</em>"
    ```python
    def _compile_sort_key(self, sort: str):
            """Compile a sort key function for better performance"""
            if "datetime" in sort.lower():
                return lambda a: a.get(sort, datetime.datetime(1970, 1, 1))
            if "date" in sort.lower():
                return lambda a: a.get(sort, datetime.date(1970, 1, 1))

            # Create a compiled function for complex sort expressions
            try:
                code = compile(sort, "<string>", "eval")

                def sort_key(a):
                    try:
                        value = eval(code, a.to_dict(), {})
                        if isinstance(value, (int, float)):
                            return value
                        if hasattr(value, "timestamp"):
                            return value.timestamp()
                        if isinstance(value, datetime.date):
                            return datetime.datetime.combine(
                                value,
                                datetime.datetime.min.time(),
                            ).timestamp()
                        return sum(ord(c) for c in str(value))
                    except Exception:
                        return -1

                return sort_key
            except Exception:
                return lambda _: -1
    ```
!!! method
    <h2 id="_get_sort_key" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">_get_sort_key <em class="small">method</em></h2>

    Cache compiled sort key functions

???+ source "_get_sort_key <em class='small'>source</em>"
    ```python
    def _get_sort_key(self, sort: str):
            """Cache compiled sort key functions"""
            return self._compile_sort_key(sort)
    ```
!!! method
    <h2 id="_get_eval_globals" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">_get_eval_globals <em class="small">method</em></h2>

    Get common globals used in eval operations

???+ source "_get_eval_globals <em class='small'>source</em>"
    ```python
    def _get_eval_globals(self):
            """Get common globals used in eval operations"""
            if not hasattr(self, "_eval_globals"):
                self._eval_globals = {"timedelta": timedelta}
            return self._eval_globals
    ```
!!! method
    <h2 id="_eval_with_article" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">_eval_with_article <em class="small">method</em></h2>

    Evaluate code with article context, reusing dict where possible

???+ source "_eval_with_article <em class='small'>source</em>"
    ```python
    def _eval_with_article(self, code, article, extra_globals=None):
            """Evaluate code with article context, reusing dict where possible"""
            if not hasattr(article, "_eval_dict"):
                article._eval_dict = article.to_dict()
                article._eval_dict.update({"post": article, "m": self})

            globals_dict = self._get_eval_globals()
            if extra_globals:
                globals_dict.update(extra_globals)

            try:
                return eval(code, article._eval_dict, globals_dict)
            except Exception:
                return None
    ```