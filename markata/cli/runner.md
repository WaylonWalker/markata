---
content: "---\n\nNone\n\n---\n\n\n\n!!! class\n    <h2 id=\"Runner\" class=\"admonition-title\"
  style=\"margin: 0; padding: .5rem 1rem;\">Runner <em class=\"small\">class</em></h2>\n\n
  \   Display Footer\n\n???+ source \"Runner <em class='small'>source</em>\"\n    ```python\n
  \   class Runner:\n        \"\"\"Display Footer\"\"\"\n\n        _status = \"waiting\"\n
  \       status = \"starting\"\n        last_error = \"\"\n        title = \"runner\"\n
  \       border = \"green\"\n\n        _dirhash = \"\"\n        time = time.time()\n\n
  \       def __init__(self, markata: \"Markata\") -> None:\n            self.m =
  markata\n            self._dirhash = self.m.content_dir_hash\n            self._run()\n
  \           atexit.register(self.kill)\n\n        def kill(self) -> None:\n            self.proc.stdout.close()\n
  \           self.proc.stderr.close()\n            self.proc.kill()\n            self.proc.wait()\n\n
  \       def run(self) -> None:\n            \"Runs the build only if one is not
  already running.\"\n            if self.proc.poll() is not None:\n                self._run()\n\n
  \       def _run(self) -> None:\n            \"Runs the build and sets the proc\"\n
  \           self.status = \"running\"\n            self.time = time.time()\n            self.proc
  = subprocess.Popen(\n                [\"markata\", \"build\"],\n                stdout=subprocess.PIPE,\n
  \               stderr=subprocess.PIPE,\n            )\n\n        @property\n        def
  status_message(self) -> str:\n            \"returns the status message to display\"\n
  \           num_lines = self.m.console.height - 4\n            last_error = \"\\n\".join(self.last_error.split(\"\\n\")[-num_lines:])\n
  \           if self.status == \"running\":\n                self.title = \"runner
  running\"\n                self.border = \"gold1\"\n            elif last_error
  == \"\":\n                self.title = \"runner succeded\"\n                self.border
  = \"green\"\n            else:\n                self.title = \"runner failed\"\n
  \               self.border = \"red\"\n            self.title = f\"{self.title}
  [blue]({round(time.time() - self.time)}s)[/]\"\n\n            return (\n                f\"runner
  is {self.status}\"\n                f\"{round(time.time() - self.time)}\\n\"\n                f\"pid:
  {self.proc.pid}\\n\"\n                f\"hash: {self.m.content_dir_hash}\\n\"\n
  \               f\"{last_error}\"\n            )\n\n        def __rich__(self) ->
  Panel:\n            if self.proc:\n                if self.proc.poll() is None:\n
  \                   return Panel(\n                        Text(self.status_message),\n
  \                       border_style=self.border,\n                        title=self.title,\n
  \                       expand=True,\n                    )\n\n            if self.status
  == \"running\":\n                self.status = \"waiting\"\n                self.time
  = time.time()\n                if self.proc:\n                    self.last_error
  = self.proc.stderr.read().decode()\n\n            if self._dirhash != self.m.content_dir_hash:\n
  \               self.run()\n                self._dirhash = self.m.content_dir_hash\n\n
  \           return Panel(\n                Text(self.status_message),\n                border_style=self.border,\n
  \               title=self.title,\n                expand=True,\n            )\n
  \   ```\n\n\n\n\n\n\n\n\n\n!!! method\n    <h2 id=\"run\" class=\"admonition-title\"
  style=\"margin: 0; padding: .5rem 1rem;\">run <em class=\"small\">method</em></h2>\n\n
  \   Runs the build only if one is not already running.\n\n???+ source \"run <em
  class='small'>source</em>\"\n    ```python\n    def run(self) -> None:\n            \"Runs
  the build only if one is not already running.\"\n            if self.proc.poll()
  is not None:\n                self._run()\n    ```\n\n\n\n!!! method\n    <h2 id=\"_run\"
  class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">_run <em class=\"small\">method</em></h2>\n\n
  \   Runs the build and sets the proc\n\n???+ source \"_run <em class='small'>source</em>\"\n
  \   ```python\n    def _run(self) -> None:\n            \"Runs the build and sets
  the proc\"\n            self.status = \"running\"\n            self.time = time.time()\n
  \           self.proc = subprocess.Popen(\n                [\"markata\", \"build\"],\n
  \               stdout=subprocess.PIPE,\n                stderr=subprocess.PIPE,\n
  \           )\n    ```\n\n\n\n!!! method\n    <h2 id=\"status_message\" class=\"admonition-title\"
  style=\"margin: 0; padding: .5rem 1rem;\">status_message <em class=\"small\">method</em></h2>\n\n
  \   returns the status message to display\n\n???+ source \"status_message <em class='small'>source</em>\"\n
  \   ```python\n    def status_message(self) -> str:\n            \"returns the status
  message to display\"\n            num_lines = self.m.console.height - 4\n            last_error
  = \"\\n\".join(self.last_error.split(\"\\n\")[-num_lines:])\n            if self.status
  == \"running\":\n                self.title = \"runner running\"\n                self.border
  = \"gold1\"\n            elif last_error == \"\":\n                self.title =
  \"runner succeded\"\n                self.border = \"green\"\n            else:\n
  \               self.title = \"runner failed\"\n                self.border = \"red\"\n
  \           self.title = f\"{self.title} [blue]({round(time.time() - self.time)}s)[/]\"\n\n
  \           return (\n                f\"runner is {self.status}\"\n                f\"{round(time.time()
  - self.time)}\\n\"\n                f\"pid: {self.proc.pid}\\n\"\n                f\"hash:
  {self.m.content_dir_hash}\\n\"\n                f\"{last_error}\"\n            )\n
  \   ```"
date: 2025-02-14
description: "None !!! class\n    &lt;h2 id=&quot;Runner&quot; class=&quot;admonition-title&quot;
  style=&quot;margin: 0; padding: .5rem 1rem;&quot;&gt;Runner &lt;em class=&quo"
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>runner.py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"None !!! class\n    &lt;h2 id=&quot;Runner&quot;
    class=&quot;admonition-title&quot; style=&quot;margin: 0; padding: .5rem 1rem;&quot;&gt;Runner
    &lt;em class=&quo\" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n        <meta property=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    <head>\n        <script src=\"https://unpkg.com/@tailwindcss/browser@4\"></script>\n<title>runner.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None !!! class\n    &lt;h2 id=&quot;Runner&quot;
    class=&quot;admonition-title&quot; style=&quot;margin: 0; padding: .5rem 1rem;&quot;&gt;Runner
    &lt;em class=&quo\" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n        <meta property=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n        <script>\n            document.addEventListener(\"DOMContentLoaded\",
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        runner.py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <hr />\n<p>None</p>\n<hr
    />\n<div class=\"admonition class\">\n<p class=\"admonition-title\">Class</p>\n<h2
    id=\"Runner\" class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">Runner
    <em class=\"small\">class</em></h2>\n<p>Display Footer</p>\n</div>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">Runner
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">class</span><span
    class=\"w\"> </span><span class=\"nc\">Runner</span><span class=\"p\">:</span>\n<span
    class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Display Footer&quot;&quot;&quot;</span>\n\n
    \   <span class=\"n\">_status</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;waiting&quot;</span>\n
    \   <span class=\"n\">status</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;starting&quot;</span>\n
    \   <span class=\"n\">last_error</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;&quot;</span>\n
    \   <span class=\"n\">title</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;runner&quot;</span>\n
    \   <span class=\"n\">border</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;green&quot;</span>\n\n
    \   <span class=\"n\">_dirhash</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;&quot;</span>\n
    \   <span class=\"n\">time</span> <span class=\"o\">=</span> <span class=\"n\">time</span><span
    class=\"o\">.</span><span class=\"n\">time</span><span class=\"p\">()</span>\n\n
    \   <span class=\"k\">def</span><span class=\"w\"> </span><span class=\"fm\">__init__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n        <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span> <span class=\"o\">=</span> <span class=\"n\">markata</span>\n
    \       <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_dirhash</span>
    <span class=\"o\">=</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"o\">.</span><span class=\"n\">content_dir_hash</span>\n
    \       <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_run</span><span
    class=\"p\">()</span>\n        <span class=\"n\">atexit</span><span class=\"o\">.</span><span
    class=\"n\">register</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">kill</span><span class=\"p\">)</span>\n\n
    \   <span class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">kill</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \       <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">proc</span><span
    class=\"o\">.</span><span class=\"n\">stdout</span><span class=\"o\">.</span><span
    class=\"n\">close</span><span class=\"p\">()</span>\n        <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">proc</span><span class=\"o\">.</span><span
    class=\"n\">stderr</span><span class=\"o\">.</span><span class=\"n\">close</span><span
    class=\"p\">()</span>\n        <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">proc</span><span class=\"o\">.</span><span class=\"n\">kill</span><span
    class=\"p\">()</span>\n        <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">proc</span><span class=\"o\">.</span><span class=\"n\">wait</span><span
    class=\"p\">()</span>\n\n    <span class=\"k\">def</span><span class=\"w\"> </span><span
    class=\"nf\">run</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n        <span class=\"s2\">&quot;Runs the build only if
    one is not already running.&quot;</span>\n        <span class=\"k\">if</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">proc</span><span
    class=\"o\">.</span><span class=\"n\">poll</span><span class=\"p\">()</span> <span
    class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n            <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_run</span><span class=\"p\">()</span>\n\n    <span class=\"k\">def</span><span
    class=\"w\"> </span><span class=\"nf\">_run</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n        <span class=\"s2\">&quot;Runs
    the build and sets the proc&quot;</span>\n        <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">status</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;running&quot;</span>\n        <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">time</span> <span class=\"o\">=</span> <span
    class=\"n\">time</span><span class=\"o\">.</span><span class=\"n\">time</span><span
    class=\"p\">()</span>\n        <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">proc</span> <span class=\"o\">=</span> <span class=\"n\">subprocess</span><span
    class=\"o\">.</span><span class=\"n\">Popen</span><span class=\"p\">(</span>\n
    \           <span class=\"p\">[</span><span class=\"s2\">&quot;markata&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;build&quot;</span><span class=\"p\">],</span>\n
    \           <span class=\"n\">stdout</span><span class=\"o\">=</span><span class=\"n\">subprocess</span><span
    class=\"o\">.</span><span class=\"n\">PIPE</span><span class=\"p\">,</span>\n
    \           <span class=\"n\">stderr</span><span class=\"o\">=</span><span class=\"n\">subprocess</span><span
    class=\"o\">.</span><span class=\"n\">PIPE</span><span class=\"p\">,</span>\n
    \       <span class=\"p\">)</span>\n\n    <span class=\"nd\">@property</span>\n
    \   <span class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">status_message</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n
    \       <span class=\"s2\">&quot;returns the status message to display&quot;</span>\n
    \       <span class=\"n\">num_lines</span> <span class=\"o\">=</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">m</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">height</span>
    <span class=\"o\">-</span> <span class=\"mi\">4</span>\n        <span class=\"n\">last_error</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span
    class=\"s2\">&quot;</span><span class=\"o\">.</span><span class=\"n\">join</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">last_error</span><span class=\"o\">.</span><span class=\"n\">split</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span
    class=\"s2\">&quot;</span><span class=\"p\">)[</span><span class=\"o\">-</span><span
    class=\"n\">num_lines</span><span class=\"p\">:])</span>\n        <span class=\"k\">if</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">status</span>
    <span class=\"o\">==</span> <span class=\"s2\">&quot;running&quot;</span><span
    class=\"p\">:</span>\n            <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">title</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;runner
    running&quot;</span>\n            <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">border</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;gold1&quot;</span>\n
    \       <span class=\"k\">elif</span> <span class=\"n\">last_error</span> <span
    class=\"o\">==</span> <span class=\"s2\">&quot;&quot;</span><span class=\"p\">:</span>\n
    \           <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">title</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;runner succeded&quot;</span>\n
    \           <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">border</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;green&quot;</span>\n        <span
    class=\"k\">else</span><span class=\"p\">:</span>\n            <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">title</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;runner failed&quot;</span>\n            <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">border</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;red&quot;</span>\n        <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">title</span> <span class=\"o\">=</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">title</span><span
    class=\"si\">}</span><span class=\"s2\"> [blue](</span><span class=\"si\">{</span><span
    class=\"nb\">round</span><span class=\"p\">(</span><span class=\"n\">time</span><span
    class=\"o\">.</span><span class=\"n\">time</span><span class=\"p\">()</span><span
    class=\"w\"> </span><span class=\"o\">-</span><span class=\"w\"> </span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">time</span><span
    class=\"p\">)</span><span class=\"si\">}</span><span class=\"s2\">s)[/]&quot;</span>\n\n
    \       <span class=\"k\">return</span> <span class=\"p\">(</span>\n            <span
    class=\"sa\">f</span><span class=\"s2\">&quot;runner is </span><span class=\"si\">{</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">status</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span>\n            <span class=\"sa\">f</span><span
    class=\"s2\">&quot;</span><span class=\"si\">{</span><span class=\"nb\">round</span><span
    class=\"p\">(</span><span class=\"n\">time</span><span class=\"o\">.</span><span
    class=\"n\">time</span><span class=\"p\">()</span><span class=\"w\"> </span><span
    class=\"o\">-</span><span class=\"w\"> </span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">time</span><span class=\"p\">)</span><span
    class=\"si\">}</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span>\n
    \           <span class=\"sa\">f</span><span class=\"s2\">&quot;pid: </span><span
    class=\"si\">{</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">proc</span><span class=\"o\">.</span><span class=\"n\">pid</span><span
    class=\"si\">}</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span>\n
    \           <span class=\"sa\">f</span><span class=\"s2\">&quot;hash: </span><span
    class=\"si\">{</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"o\">.</span><span class=\"n\">content_dir_hash</span><span
    class=\"si\">}</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span>\n
    \           <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">last_error</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span>\n
    \       <span class=\"p\">)</span>\n\n    <span class=\"k\">def</span><span class=\"w\">
    </span><span class=\"nf\">__rich__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Panel</span><span
    class=\"p\">:</span>\n        <span class=\"k\">if</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">proc</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">if</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">proc</span><span class=\"o\">.</span><span class=\"n\">poll</span><span
    class=\"p\">()</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"n\">Panel</span><span
    class=\"p\">(</span>\n                    <span class=\"n\">Text</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">status_message</span><span
    class=\"p\">),</span>\n                    <span class=\"n\">border_style</span><span
    class=\"o\">=</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">border</span><span class=\"p\">,</span>\n                    <span
    class=\"n\">title</span><span class=\"o\">=</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">title</span><span class=\"p\">,</span>\n
    \                   <span class=\"n\">expand</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n                <span class=\"p\">)</span>\n\n
    \       <span class=\"k\">if</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">status</span> <span class=\"o\">==</span> <span class=\"s2\">&quot;running&quot;</span><span
    class=\"p\">:</span>\n            <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">status</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;waiting&quot;</span>\n
    \           <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">time</span>
    <span class=\"o\">=</span> <span class=\"n\">time</span><span class=\"o\">.</span><span
    class=\"n\">time</span><span class=\"p\">()</span>\n            <span class=\"k\">if</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">proc</span><span
    class=\"p\">:</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">last_error</span> <span class=\"o\">=</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">proc</span><span class=\"o\">.</span><span
    class=\"n\">stderr</span><span class=\"o\">.</span><span class=\"n\">read</span><span
    class=\"p\">()</span><span class=\"o\">.</span><span class=\"n\">decode</span><span
    class=\"p\">()</span>\n\n        <span class=\"k\">if</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_dirhash</span> <span class=\"o\">!=</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">m</span><span
    class=\"o\">.</span><span class=\"n\">content_dir_hash</span><span class=\"p\">:</span>\n
    \           <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">run</span><span
    class=\"p\">()</span>\n            <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_dirhash</span> <span class=\"o\">=</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">m</span><span class=\"o\">.</span><span
    class=\"n\">content_dir_hash</span>\n\n        <span class=\"k\">return</span>
    <span class=\"n\">Panel</span><span class=\"p\">(</span>\n            <span class=\"n\">Text</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">status_message</span><span class=\"p\">),</span>\n            <span
    class=\"n\">border_style</span><span class=\"o\">=</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">border</span><span class=\"p\">,</span>\n
    \           <span class=\"n\">title</span><span class=\"o\">=</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">title</span><span class=\"p\">,</span>\n
    \           <span class=\"n\">expand</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span>\n        <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"run\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">run <em class=\"small\">method</em></h2>\n<p>Runs
    the build only if one is not already running.</p>\n</div>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">run <em
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
    class=\"w\"> </span><span class=\"nf\">run</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n        <span class=\"s2\">&quot;Runs
    the build only if one is not already running.&quot;</span>\n        <span class=\"k\">if</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">proc</span><span
    class=\"o\">.</span><span class=\"n\">poll</span><span class=\"p\">()</span> <span
    class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n            <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_run</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"_run\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">_run <em
    class=\"small\">method</em></h2>\n<p>Runs the build and sets the proc</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_run
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
    class=\"w\"> </span><span class=\"nf\">_run</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n        <span class=\"s2\">&quot;Runs
    the build and sets the proc&quot;</span>\n        <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">status</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;running&quot;</span>\n        <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">time</span> <span class=\"o\">=</span> <span
    class=\"n\">time</span><span class=\"o\">.</span><span class=\"n\">time</span><span
    class=\"p\">()</span>\n        <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">proc</span> <span class=\"o\">=</span> <span class=\"n\">subprocess</span><span
    class=\"o\">.</span><span class=\"n\">Popen</span><span class=\"p\">(</span>\n
    \           <span class=\"p\">[</span><span class=\"s2\">&quot;markata&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;build&quot;</span><span class=\"p\">],</span>\n
    \           <span class=\"n\">stdout</span><span class=\"o\">=</span><span class=\"n\">subprocess</span><span
    class=\"o\">.</span><span class=\"n\">PIPE</span><span class=\"p\">,</span>\n
    \           <span class=\"n\">stderr</span><span class=\"o\">=</span><span class=\"n\">subprocess</span><span
    class=\"o\">.</span><span class=\"n\">PIPE</span><span class=\"p\">,</span>\n
    \       <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div class=\"admonition
    method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"status_message\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">status_message
    <em class=\"small\">method</em></h2>\n<p>returns the status message to display</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">status_message
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
    class=\"w\"> </span><span class=\"nf\">status_message</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"nb\">str</span><span class=\"p\">:</span>\n        <span class=\"s2\">&quot;returns
    the status message to display&quot;</span>\n        <span class=\"n\">num_lines</span>
    <span class=\"o\">=</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">height</span> <span class=\"o\">-</span>
    <span class=\"mi\">4</span>\n        <span class=\"n\">last_error</span> <span
    class=\"o\">=</span> <span class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span
    class=\"s2\">&quot;</span><span class=\"o\">.</span><span class=\"n\">join</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">last_error</span><span class=\"o\">.</span><span class=\"n\">split</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span
    class=\"s2\">&quot;</span><span class=\"p\">)[</span><span class=\"o\">-</span><span
    class=\"n\">num_lines</span><span class=\"p\">:])</span>\n        <span class=\"k\">if</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">status</span>
    <span class=\"o\">==</span> <span class=\"s2\">&quot;running&quot;</span><span
    class=\"p\">:</span>\n            <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">title</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;runner
    running&quot;</span>\n            <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">border</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;gold1&quot;</span>\n
    \       <span class=\"k\">elif</span> <span class=\"n\">last_error</span> <span
    class=\"o\">==</span> <span class=\"s2\">&quot;&quot;</span><span class=\"p\">:</span>\n
    \           <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">title</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;runner succeded&quot;</span>\n
    \           <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">border</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;green&quot;</span>\n        <span
    class=\"k\">else</span><span class=\"p\">:</span>\n            <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">title</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;runner failed&quot;</span>\n            <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">border</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;red&quot;</span>\n        <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">title</span> <span class=\"o\">=</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">title</span><span
    class=\"si\">}</span><span class=\"s2\"> [blue](</span><span class=\"si\">{</span><span
    class=\"nb\">round</span><span class=\"p\">(</span><span class=\"n\">time</span><span
    class=\"o\">.</span><span class=\"n\">time</span><span class=\"p\">()</span><span
    class=\"w\"> </span><span class=\"o\">-</span><span class=\"w\"> </span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">time</span><span
    class=\"p\">)</span><span class=\"si\">}</span><span class=\"s2\">s)[/]&quot;</span>\n\n
    \       <span class=\"k\">return</span> <span class=\"p\">(</span>\n            <span
    class=\"sa\">f</span><span class=\"s2\">&quot;runner is </span><span class=\"si\">{</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">status</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span>\n            <span class=\"sa\">f</span><span
    class=\"s2\">&quot;</span><span class=\"si\">{</span><span class=\"nb\">round</span><span
    class=\"p\">(</span><span class=\"n\">time</span><span class=\"o\">.</span><span
    class=\"n\">time</span><span class=\"p\">()</span><span class=\"w\"> </span><span
    class=\"o\">-</span><span class=\"w\"> </span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">time</span><span class=\"p\">)</span><span
    class=\"si\">}</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span>\n
    \           <span class=\"sa\">f</span><span class=\"s2\">&quot;pid: </span><span
    class=\"si\">{</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">proc</span><span class=\"o\">.</span><span class=\"n\">pid</span><span
    class=\"si\">}</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span>\n
    \           <span class=\"sa\">f</span><span class=\"s2\">&quot;hash: </span><span
    class=\"si\">{</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"o\">.</span><span class=\"n\">content_dir_hash</span><span
    class=\"si\">}</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span>\n
    \           <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">last_error</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span>\n
    \       <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n\n    </section>\n</article>
    \   </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>runner.py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"None !!! class\n    &lt;h2 id=&quot;Runner&quot;
    class=&quot;admonition-title&quot; style=&quot;margin: 0; padding: .5rem 1rem;&quot;&gt;Runner
    &lt;em class=&quo\" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n        <meta property=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    <head>\n        <script src=\"https://unpkg.com/@tailwindcss/browser@4\"></script>\n<title>runner.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None !!! class\n    &lt;h2 id=&quot;Runner&quot;
    class=&quot;admonition-title&quot; style=&quot;margin: 0; padding: .5rem 1rem;&quot;&gt;Runner
    &lt;em class=&quo\" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n        <meta property=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n        <script>\n            document.addEventListener(\"DOMContentLoaded\",
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
    class=\"title\">\n    <h1 id=\"title\">\n        runner.py\n    </h1>\n</section></article>\n
    \    </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       runner.py\n    </h1>\n</section>    <section class=\"body\">\n        <hr
    />\n<p>None</p>\n<hr />\n<div class=\"admonition class\">\n<p class=\"admonition-title\">Class</p>\n<h2
    id=\"Runner\" class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">Runner
    <em class=\"small\">class</em></h2>\n<p>Display Footer</p>\n</div>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">Runner
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">class</span><span
    class=\"w\"> </span><span class=\"nc\">Runner</span><span class=\"p\">:</span>\n<span
    class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Display Footer&quot;&quot;&quot;</span>\n\n
    \   <span class=\"n\">_status</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;waiting&quot;</span>\n
    \   <span class=\"n\">status</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;starting&quot;</span>\n
    \   <span class=\"n\">last_error</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;&quot;</span>\n
    \   <span class=\"n\">title</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;runner&quot;</span>\n
    \   <span class=\"n\">border</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;green&quot;</span>\n\n
    \   <span class=\"n\">_dirhash</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;&quot;</span>\n
    \   <span class=\"n\">time</span> <span class=\"o\">=</span> <span class=\"n\">time</span><span
    class=\"o\">.</span><span class=\"n\">time</span><span class=\"p\">()</span>\n\n
    \   <span class=\"k\">def</span><span class=\"w\"> </span><span class=\"fm\">__init__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n        <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span> <span class=\"o\">=</span> <span class=\"n\">markata</span>\n
    \       <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_dirhash</span>
    <span class=\"o\">=</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"o\">.</span><span class=\"n\">content_dir_hash</span>\n
    \       <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_run</span><span
    class=\"p\">()</span>\n        <span class=\"n\">atexit</span><span class=\"o\">.</span><span
    class=\"n\">register</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">kill</span><span class=\"p\">)</span>\n\n
    \   <span class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">kill</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \       <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">proc</span><span
    class=\"o\">.</span><span class=\"n\">stdout</span><span class=\"o\">.</span><span
    class=\"n\">close</span><span class=\"p\">()</span>\n        <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">proc</span><span class=\"o\">.</span><span
    class=\"n\">stderr</span><span class=\"o\">.</span><span class=\"n\">close</span><span
    class=\"p\">()</span>\n        <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">proc</span><span class=\"o\">.</span><span class=\"n\">kill</span><span
    class=\"p\">()</span>\n        <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">proc</span><span class=\"o\">.</span><span class=\"n\">wait</span><span
    class=\"p\">()</span>\n\n    <span class=\"k\">def</span><span class=\"w\"> </span><span
    class=\"nf\">run</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n        <span class=\"s2\">&quot;Runs the build only if
    one is not already running.&quot;</span>\n        <span class=\"k\">if</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">proc</span><span
    class=\"o\">.</span><span class=\"n\">poll</span><span class=\"p\">()</span> <span
    class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n            <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_run</span><span class=\"p\">()</span>\n\n    <span class=\"k\">def</span><span
    class=\"w\"> </span><span class=\"nf\">_run</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n        <span class=\"s2\">&quot;Runs
    the build and sets the proc&quot;</span>\n        <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">status</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;running&quot;</span>\n        <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">time</span> <span class=\"o\">=</span> <span
    class=\"n\">time</span><span class=\"o\">.</span><span class=\"n\">time</span><span
    class=\"p\">()</span>\n        <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">proc</span> <span class=\"o\">=</span> <span class=\"n\">subprocess</span><span
    class=\"o\">.</span><span class=\"n\">Popen</span><span class=\"p\">(</span>\n
    \           <span class=\"p\">[</span><span class=\"s2\">&quot;markata&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;build&quot;</span><span class=\"p\">],</span>\n
    \           <span class=\"n\">stdout</span><span class=\"o\">=</span><span class=\"n\">subprocess</span><span
    class=\"o\">.</span><span class=\"n\">PIPE</span><span class=\"p\">,</span>\n
    \           <span class=\"n\">stderr</span><span class=\"o\">=</span><span class=\"n\">subprocess</span><span
    class=\"o\">.</span><span class=\"n\">PIPE</span><span class=\"p\">,</span>\n
    \       <span class=\"p\">)</span>\n\n    <span class=\"nd\">@property</span>\n
    \   <span class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">status_message</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n
    \       <span class=\"s2\">&quot;returns the status message to display&quot;</span>\n
    \       <span class=\"n\">num_lines</span> <span class=\"o\">=</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">m</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">height</span>
    <span class=\"o\">-</span> <span class=\"mi\">4</span>\n        <span class=\"n\">last_error</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span
    class=\"s2\">&quot;</span><span class=\"o\">.</span><span class=\"n\">join</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">last_error</span><span class=\"o\">.</span><span class=\"n\">split</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span
    class=\"s2\">&quot;</span><span class=\"p\">)[</span><span class=\"o\">-</span><span
    class=\"n\">num_lines</span><span class=\"p\">:])</span>\n        <span class=\"k\">if</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">status</span>
    <span class=\"o\">==</span> <span class=\"s2\">&quot;running&quot;</span><span
    class=\"p\">:</span>\n            <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">title</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;runner
    running&quot;</span>\n            <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">border</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;gold1&quot;</span>\n
    \       <span class=\"k\">elif</span> <span class=\"n\">last_error</span> <span
    class=\"o\">==</span> <span class=\"s2\">&quot;&quot;</span><span class=\"p\">:</span>\n
    \           <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">title</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;runner succeded&quot;</span>\n
    \           <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">border</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;green&quot;</span>\n        <span
    class=\"k\">else</span><span class=\"p\">:</span>\n            <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">title</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;runner failed&quot;</span>\n            <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">border</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;red&quot;</span>\n        <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">title</span> <span class=\"o\">=</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">title</span><span
    class=\"si\">}</span><span class=\"s2\"> [blue](</span><span class=\"si\">{</span><span
    class=\"nb\">round</span><span class=\"p\">(</span><span class=\"n\">time</span><span
    class=\"o\">.</span><span class=\"n\">time</span><span class=\"p\">()</span><span
    class=\"w\"> </span><span class=\"o\">-</span><span class=\"w\"> </span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">time</span><span
    class=\"p\">)</span><span class=\"si\">}</span><span class=\"s2\">s)[/]&quot;</span>\n\n
    \       <span class=\"k\">return</span> <span class=\"p\">(</span>\n            <span
    class=\"sa\">f</span><span class=\"s2\">&quot;runner is </span><span class=\"si\">{</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">status</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span>\n            <span class=\"sa\">f</span><span
    class=\"s2\">&quot;</span><span class=\"si\">{</span><span class=\"nb\">round</span><span
    class=\"p\">(</span><span class=\"n\">time</span><span class=\"o\">.</span><span
    class=\"n\">time</span><span class=\"p\">()</span><span class=\"w\"> </span><span
    class=\"o\">-</span><span class=\"w\"> </span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">time</span><span class=\"p\">)</span><span
    class=\"si\">}</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span>\n
    \           <span class=\"sa\">f</span><span class=\"s2\">&quot;pid: </span><span
    class=\"si\">{</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">proc</span><span class=\"o\">.</span><span class=\"n\">pid</span><span
    class=\"si\">}</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span>\n
    \           <span class=\"sa\">f</span><span class=\"s2\">&quot;hash: </span><span
    class=\"si\">{</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"o\">.</span><span class=\"n\">content_dir_hash</span><span
    class=\"si\">}</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span>\n
    \           <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">last_error</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span>\n
    \       <span class=\"p\">)</span>\n\n    <span class=\"k\">def</span><span class=\"w\">
    </span><span class=\"nf\">__rich__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Panel</span><span
    class=\"p\">:</span>\n        <span class=\"k\">if</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">proc</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">if</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">proc</span><span class=\"o\">.</span><span class=\"n\">poll</span><span
    class=\"p\">()</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"n\">Panel</span><span
    class=\"p\">(</span>\n                    <span class=\"n\">Text</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">status_message</span><span
    class=\"p\">),</span>\n                    <span class=\"n\">border_style</span><span
    class=\"o\">=</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">border</span><span class=\"p\">,</span>\n                    <span
    class=\"n\">title</span><span class=\"o\">=</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">title</span><span class=\"p\">,</span>\n
    \                   <span class=\"n\">expand</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n                <span class=\"p\">)</span>\n\n
    \       <span class=\"k\">if</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">status</span> <span class=\"o\">==</span> <span class=\"s2\">&quot;running&quot;</span><span
    class=\"p\">:</span>\n            <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">status</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;waiting&quot;</span>\n
    \           <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">time</span>
    <span class=\"o\">=</span> <span class=\"n\">time</span><span class=\"o\">.</span><span
    class=\"n\">time</span><span class=\"p\">()</span>\n            <span class=\"k\">if</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">proc</span><span
    class=\"p\">:</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">last_error</span> <span class=\"o\">=</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">proc</span><span class=\"o\">.</span><span
    class=\"n\">stderr</span><span class=\"o\">.</span><span class=\"n\">read</span><span
    class=\"p\">()</span><span class=\"o\">.</span><span class=\"n\">decode</span><span
    class=\"p\">()</span>\n\n        <span class=\"k\">if</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_dirhash</span> <span class=\"o\">!=</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">m</span><span
    class=\"o\">.</span><span class=\"n\">content_dir_hash</span><span class=\"p\">:</span>\n
    \           <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">run</span><span
    class=\"p\">()</span>\n            <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_dirhash</span> <span class=\"o\">=</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">m</span><span class=\"o\">.</span><span
    class=\"n\">content_dir_hash</span>\n\n        <span class=\"k\">return</span>
    <span class=\"n\">Panel</span><span class=\"p\">(</span>\n            <span class=\"n\">Text</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">status_message</span><span class=\"p\">),</span>\n            <span
    class=\"n\">border_style</span><span class=\"o\">=</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">border</span><span class=\"p\">,</span>\n
    \           <span class=\"n\">title</span><span class=\"o\">=</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">title</span><span class=\"p\">,</span>\n
    \           <span class=\"n\">expand</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span>\n        <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"run\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">run <em class=\"small\">method</em></h2>\n<p>Runs
    the build only if one is not already running.</p>\n</div>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">run <em
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
    class=\"w\"> </span><span class=\"nf\">run</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n        <span class=\"s2\">&quot;Runs
    the build only if one is not already running.&quot;</span>\n        <span class=\"k\">if</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">proc</span><span
    class=\"o\">.</span><span class=\"n\">poll</span><span class=\"p\">()</span> <span
    class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n            <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_run</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"_run\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">_run <em
    class=\"small\">method</em></h2>\n<p>Runs the build and sets the proc</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_run
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
    class=\"w\"> </span><span class=\"nf\">_run</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n        <span class=\"s2\">&quot;Runs
    the build and sets the proc&quot;</span>\n        <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">status</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;running&quot;</span>\n        <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">time</span> <span class=\"o\">=</span> <span
    class=\"n\">time</span><span class=\"o\">.</span><span class=\"n\">time</span><span
    class=\"p\">()</span>\n        <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">proc</span> <span class=\"o\">=</span> <span class=\"n\">subprocess</span><span
    class=\"o\">.</span><span class=\"n\">Popen</span><span class=\"p\">(</span>\n
    \           <span class=\"p\">[</span><span class=\"s2\">&quot;markata&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;build&quot;</span><span class=\"p\">],</span>\n
    \           <span class=\"n\">stdout</span><span class=\"o\">=</span><span class=\"n\">subprocess</span><span
    class=\"o\">.</span><span class=\"n\">PIPE</span><span class=\"p\">,</span>\n
    \           <span class=\"n\">stderr</span><span class=\"o\">=</span><span class=\"n\">subprocess</span><span
    class=\"o\">.</span><span class=\"n\">PIPE</span><span class=\"p\">,</span>\n
    \       <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div class=\"admonition
    method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"status_message\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">status_message
    <em class=\"small\">method</em></h2>\n<p>returns the status message to display</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">status_message
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
    class=\"w\"> </span><span class=\"nf\">status_message</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"nb\">str</span><span class=\"p\">:</span>\n        <span class=\"s2\">&quot;returns
    the status message to display&quot;</span>\n        <span class=\"n\">num_lines</span>
    <span class=\"o\">=</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">height</span> <span class=\"o\">-</span>
    <span class=\"mi\">4</span>\n        <span class=\"n\">last_error</span> <span
    class=\"o\">=</span> <span class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span
    class=\"s2\">&quot;</span><span class=\"o\">.</span><span class=\"n\">join</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">last_error</span><span class=\"o\">.</span><span class=\"n\">split</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span
    class=\"s2\">&quot;</span><span class=\"p\">)[</span><span class=\"o\">-</span><span
    class=\"n\">num_lines</span><span class=\"p\">:])</span>\n        <span class=\"k\">if</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">status</span>
    <span class=\"o\">==</span> <span class=\"s2\">&quot;running&quot;</span><span
    class=\"p\">:</span>\n            <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">title</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;runner
    running&quot;</span>\n            <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">border</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;gold1&quot;</span>\n
    \       <span class=\"k\">elif</span> <span class=\"n\">last_error</span> <span
    class=\"o\">==</span> <span class=\"s2\">&quot;&quot;</span><span class=\"p\">:</span>\n
    \           <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">title</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;runner succeded&quot;</span>\n
    \           <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">border</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;green&quot;</span>\n        <span
    class=\"k\">else</span><span class=\"p\">:</span>\n            <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">title</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;runner failed&quot;</span>\n            <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">border</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;red&quot;</span>\n        <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">title</span> <span class=\"o\">=</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">title</span><span
    class=\"si\">}</span><span class=\"s2\"> [blue](</span><span class=\"si\">{</span><span
    class=\"nb\">round</span><span class=\"p\">(</span><span class=\"n\">time</span><span
    class=\"o\">.</span><span class=\"n\">time</span><span class=\"p\">()</span><span
    class=\"w\"> </span><span class=\"o\">-</span><span class=\"w\"> </span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">time</span><span
    class=\"p\">)</span><span class=\"si\">}</span><span class=\"s2\">s)[/]&quot;</span>\n\n
    \       <span class=\"k\">return</span> <span class=\"p\">(</span>\n            <span
    class=\"sa\">f</span><span class=\"s2\">&quot;runner is </span><span class=\"si\">{</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">status</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span>\n            <span class=\"sa\">f</span><span
    class=\"s2\">&quot;</span><span class=\"si\">{</span><span class=\"nb\">round</span><span
    class=\"p\">(</span><span class=\"n\">time</span><span class=\"o\">.</span><span
    class=\"n\">time</span><span class=\"p\">()</span><span class=\"w\"> </span><span
    class=\"o\">-</span><span class=\"w\"> </span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">time</span><span class=\"p\">)</span><span
    class=\"si\">}</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span>\n
    \           <span class=\"sa\">f</span><span class=\"s2\">&quot;pid: </span><span
    class=\"si\">{</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">proc</span><span class=\"o\">.</span><span class=\"n\">pid</span><span
    class=\"si\">}</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span>\n
    \           <span class=\"sa\">f</span><span class=\"s2\">&quot;hash: </span><span
    class=\"si\">{</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"o\">.</span><span class=\"n\">content_dir_hash</span><span
    class=\"si\">}</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span>\n
    \           <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">last_error</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span>\n
    \       <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n\n    </section>\n</article>"
  raw.md: ''
published: false
slug: markata/cli/runner
title: runner.py


---

---

None

---



!!! class
    <h2 id="Runner" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">Runner <em class="small">class</em></h2>

    Display Footer

???+ source "Runner <em class='small'>source</em>"
    ```python
    class Runner:
        """Display Footer"""

        _status = "waiting"
        status = "starting"
        last_error = ""
        title = "runner"
        border = "green"

        _dirhash = ""
        time = time.time()

        def __init__(self, markata: "Markata") -> None:
            self.m = markata
            self._dirhash = self.m.content_dir_hash
            self._run()
            atexit.register(self.kill)

        def kill(self) -> None:
            self.proc.stdout.close()
            self.proc.stderr.close()
            self.proc.kill()
            self.proc.wait()

        def run(self) -> None:
            "Runs the build only if one is not already running."
            if self.proc.poll() is not None:
                self._run()

        def _run(self) -> None:
            "Runs the build and sets the proc"
            self.status = "running"
            self.time = time.time()
            self.proc = subprocess.Popen(
                ["markata", "build"],
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
            )

        @property
        def status_message(self) -> str:
            "returns the status message to display"
            num_lines = self.m.console.height - 4
            last_error = "\n".join(self.last_error.split("\n")[-num_lines:])
            if self.status == "running":
                self.title = "runner running"
                self.border = "gold1"
            elif last_error == "":
                self.title = "runner succeded"
                self.border = "green"
            else:
                self.title = "runner failed"
                self.border = "red"
            self.title = f"{self.title} [blue]({round(time.time() - self.time)}s)[/]"

            return (
                f"runner is {self.status}"
                f"{round(time.time() - self.time)}\n"
                f"pid: {self.proc.pid}\n"
                f"hash: {self.m.content_dir_hash}\n"
                f"{last_error}"
            )

        def __rich__(self) -> Panel:
            if self.proc:
                if self.proc.poll() is None:
                    return Panel(
                        Text(self.status_message),
                        border_style=self.border,
                        title=self.title,
                        expand=True,
                    )

            if self.status == "running":
                self.status = "waiting"
                self.time = time.time()
                if self.proc:
                    self.last_error = self.proc.stderr.read().decode()

            if self._dirhash != self.m.content_dir_hash:
                self.run()
                self._dirhash = self.m.content_dir_hash

            return Panel(
                Text(self.status_message),
                border_style=self.border,
                title=self.title,
                expand=True,
            )
    ```









!!! method
    <h2 id="run" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">run <em class="small">method</em></h2>

    Runs the build only if one is not already running.

???+ source "run <em class='small'>source</em>"
    ```python
    def run(self) -> None:
            "Runs the build only if one is not already running."
            if self.proc.poll() is not None:
                self._run()
    ```



!!! method
    <h2 id="_run" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">_run <em class="small">method</em></h2>

    Runs the build and sets the proc

???+ source "_run <em class='small'>source</em>"
    ```python
    def _run(self) -> None:
            "Runs the build and sets the proc"
            self.status = "running"
            self.time = time.time()
            self.proc = subprocess.Popen(
                ["markata", "build"],
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
            )
    ```



!!! method
    <h2 id="status_message" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">status_message <em class="small">method</em></h2>

    returns the status message to display

???+ source "status_message <em class='small'>source</em>"
    ```python
    def status_message(self) -> str:
            "returns the status message to display"
            num_lines = self.m.console.height - 4
            last_error = "\n".join(self.last_error.split("\n")[-num_lines:])
            if self.status == "running":
                self.title = "runner running"
                self.border = "gold1"
            elif last_error == "":
                self.title = "runner succeded"
                self.border = "green"
            else:
                self.title = "runner failed"
                self.border = "red"
            self.title = f"{self.title} [blue]({round(time.time() - self.time)}s)[/]"

            return (
                f"runner is {self.status}"
                f"{round(time.time() - self.time)}\n"
                f"pid: {self.proc.pid}\n"
                f"hash: {self.m.content_dir_hash}\n"
                f"{last_error}"
            )
    ```