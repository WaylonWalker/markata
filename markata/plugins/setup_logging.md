---
content: "---\n\nSetup Logging hook sets up the RichHandler for pretty console logs,
  and file\nlogs to the configured markata's configured `log_dir`, or `output_dir/_logs`
  if\n`log_dir` is not configured.  The log file will be named after the\n`<levelname>.log`\n\n#
  The log files\n\nThere will be 6 log files created based on log level and file type.\n\n```\nmarkout/_logs\n\u251C\u2500\u2500
  debug\n\u2502   \u2514\u2500\u2500 index.html\n\u251C\u2500\u2500 debug.log\n\u251C\u2500\u2500
  info\n\u2502   \u2514\u2500\u2500 index.html\n\u251C\u2500\u2500 info.log\n\u251C\u2500\u2500
  warning\n\u2502   \u2514\u2500\u2500 index.html\n\u2514\u2500\u2500 warning.log\n```\n\n##
  Configuration\n\nEnsure that setup_logging is in your hooks.  You can check if `setup_logging`\nis
  in your hooks by running `markata list --hooks` from your terminal and\nchecking
  the output, or creating an instance of `Markata()` and checking the\n`Markata().hooks`
  attribute.  If its missing or you wan to be more explicit,\nyou can add `setup_logging`
  to your `markata.toml` `[markata.hooks]`.\n\n``` toml\n[markata]\n\n# make sure
  its in your list of hooks\nhooks=[\n   \"markata.plugins.setup_logging\",\n   ]\n```\n\n#
  Log Template\n``` toml\n[markata]\n\n# make sure its in your list of hooks\nhooks=[\n
  \  \"markata.plugins.setup_logging\",\n   ]\n\n# point log template to the path
  of your logging template\nlog_template='templates/log_template.html'\n```\n\nYou
  can see the latest default `log_template` on\n[GitHub](https://github.com/WaylonWalker/markata/blob/main/markata/plugins/default_log_template.html)\n\n#
  Disable Logging\n\nIf you do not want logging, you can explicityly disable it by
  adding it to your\n`[markata.disabled_hooks]` array in your `[markata.toml]`\n\n```
  toml\n[markata]\n\n# make sure its in your list of hooks\ndisabled_hooks=[\n   \"markata.plugins.setup_logging\",\n
  \  ]\n```\n\n---\n\n!!! function\n    <h2 id=\"has_rich_handler\" class=\"admonition-title\"
  style=\"margin: 0; padding: .5rem 1rem;\">has_rich_handler <em class=\"small\">function</em></h2>\n\n
  \   Returns a boolean whether or not there is a RichHandler attached to the\n    root
  logger.\n\n???+ source \"has_rich_handler <em class='small'>source</em>\"\n    ```python\n
  \   def has_rich_handler() -> bool:\n        \"\"\"\n        Returns a boolean whether
  or not there is a RichHandler attached to the\n        root logger.\n        \"\"\"\n\n
  \       logger = logging.getLogger()\n        return bool([h for h in logger.handlers
  if isinstance(h, RichHandler)])\n    ```\n!!! function\n    <h2 id=\"setup_text_log\"
  class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">setup_text_log
  <em class=\"small\">function</em></h2>\n\n    sets up a plain text log in markata's
  configured `log_dir`, or\n    `output_dir/_logs` if `log_dir` is not configured.
  \ The log file will be\n    named after the `<levelname>.log`\n\n???+ source \"setup_text_log
  <em class='small'>source</em>\"\n    ```python\n    def setup_text_log(markata:
  \"Markata\", level: int = logging.INFO) -> Path:\n        \"\"\"\n        sets up
  a plain text log in markata's configured `log_dir`, or\n        `output_dir/_logs`
  if `log_dir` is not configured.  The log file will be\n        named after the `<levelname>.log`\n
  \       \"\"\"\n        log_file = markata.config.logging.log_dir / (\n            logging.getLevelName(level).lower()
  + \".log\"\n        )\n\n        if has_file_handler(log_file):\n            return
  log_file\n\n        if not log_file.parent.exists():\n            log_file.parent.mkdir(parents=True)\n
  \       fh = logging.FileHandler(log_file)\n        fh.setLevel(level)\n        fh_formatter
  = logging.Formatter(\n            \"%(asctime)s %(name)-12s %(levelname)-8s %(message)s\",\n
  \       )\n        fh.setFormatter(fh_formatter)\n        logging.getLogger(\"\").addHandler(fh)\n\n
  \       return log_file\n    ```\n!!! function\n    <h2 id=\"setup_html_log\" class=\"admonition-title\"
  style=\"margin: 0; padding: .5rem 1rem;\">setup_html_log <em class=\"small\">function</em></h2>\n\n
  \   sets up an html log in markata's configured `log_dir`, or\n    `output_dir/_logs`
  if `log_dir` is not configured.  The log file will be\n    named after the `<levelname>/index.html`.
  \ The goal of this is to give\n\n???+ source \"setup_html_log <em class='small'>source</em>\"\n
  \   ```python\n    def setup_html_log(markata: \"Markata\", level: int = logging.INFO)
  -> Path:\n        \"\"\"\n        sets up an html log in markata's configured `log_dir`,
  or\n        `output_dir/_logs` if `log_dir` is not configured.  The log file will
  be\n        named after the `<levelname>/index.html`.  The goal of this is to give\n
  \       \"\"\"\n\n        log_file = (\n            markata.config.logging.log_dir\n
  \           / logging.getLevelName(level).lower()\n            / \"index.html\"\n
  \       )\n\n        if has_file_handler(log_file):\n            return log_file\n\n
  \       log_file.parent.mkdir(parents=True, exist_ok=True)\n\n        if not log_file.exists():\n
  \           template = Template(\n                markata.config.logging.template.read_text(),
  undefined=SilentUndefined\n            )\n            log_header = template.render(\n
  \               title=markata.config.title + \" logs\",\n                config=markata.config,\n
  \           )\n            log_file.write_text(log_header)\n        with open(log_file,
  \"a\") as f:\n            command = Path(sys.argv[0]).name + \" \" + \" \".join(sys.argv[1:])\n
  \           f.write(\n                f\"\"\"\n                <div style=\"\n                width:
  100%;\n                height: 20px;\n                margin-top: 5rem;\n                border-bottom:
  1px solid goldenrod;\n                text-align: center\">\n                    <span
  style=\"padding: 0 10px;\">\n                        {datetime.datetime.now()} running
  \"{command}\"\n                    </span>\n                </div>\n        \"\"\",\n
  \           )\n        fh = logging.FileHandler(log_file)\n        fh.setLevel(level)\n
  \       fh_formatter = logging.Formatter(\n            \"\"\"\n            <li>\n
  \               <p>\n                    <span class=\"time\">%(asctime)s</span>\n
  \                   <span class=\"name %(name)s\">%(name)-12s</span>\n                    <span
  class=\"levelname %(levelname)s\">%(levelname)-8s</span>\n                </p>\n
  \               <p class=\"message\">%(message)s</p>\n            </li>\n            \"\"\",\n
  \       )\n        fh.setFormatter(fh_formatter)\n        logging.getLogger(\"\").addHandler(fh)\n\n
  \       return log_file\n    ```"
date: 2025-05-05
description: "Setup Logging hook sets up the RichHandler for pretty console logs,
  and file logs to the configured markata's configured  , or   if  is not configured.
  \ The\u2026"
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>setup_logging.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Setup Logging hook sets up the RichHandler
    for pretty console logs, and file logs to the configured markata's configured
    \ , or   if  is not configured.  The\u2026\" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n
    \       <meta property=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    <head>\n        <script src=\"https://unpkg.com/@tailwindcss/browser@4\"></script>\n<title>setup_logging.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Setup Logging hook sets up the RichHandler
    for pretty console logs, and file logs to the configured markata's configured
    \ , or   if  is not configured.  The\u2026\" />\n <link href=\"/favicon.ico\"
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        setup_logging.py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <hr />\n<p>Setup Logging
    hook sets up the RichHandler for pretty console logs, and file\nlogs to the configured
    markata's configured <code>log_dir</code>, or <code>output_dir/_logs</code> if\n<code>log_dir</code>
    is not configured.  The log file will be named after the\n<code>&lt;levelname&gt;.log</code></p>\n<h1
    id=\"the-log-files\">The log files <a class=\"header-anchor\" href=\"#the-log-files\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>There will be 6 log
    files created based on log level and file type.</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span>markout/_logs\n\u251C\u2500\u2500
    debug\n\u2502   \u2514\u2500\u2500 index.html\n\u251C\u2500\u2500 debug.log\n\u251C\u2500\u2500
    info\n\u2502   \u2514\u2500\u2500 index.html\n\u251C\u2500\u2500 info.log\n\u251C\u2500\u2500
    warning\n\u2502   \u2514\u2500\u2500 index.html\n\u2514\u2500\u2500 warning.log\n</pre></div>\n\n</pre>\n\n<h2
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Ensure that setup_logging
    is in your hooks.  You can check if <code>setup_logging</code>\nis in your hooks
    by running <code>markata list --hooks</code> from your terminal and\nchecking
    the output, or creating an instance of <code>Markata()</code> and checking the\n<code>Markata().hooks</code>
    attribute.  If its missing or you wan to be more explicit,\nyou can add <code>setup_logging</code>
    to your <code>markata.toml</code> <code>[markata.hooks]</code>.</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n\n<span
    class=\"c1\"># make sure its in your list of hooks</span>\n<span class=\"n\">hooks</span><span
    class=\"o\">=</span><span class=\"p\">[</span>\n<span class=\"w\">   </span><span
    class=\"s2\">&quot;markata.plugins.setup_logging&quot;</span><span class=\"p\">,</span>\n<span
    class=\"w\">   </span><span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<h1
    id=\"log-template\">Log Template <a class=\"header-anchor\" href=\"#log-template\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n\n<span
    class=\"c1\"># make sure its in your list of hooks</span>\n<span class=\"n\">hooks</span><span
    class=\"o\">=</span><span class=\"p\">[</span>\n<span class=\"w\">   </span><span
    class=\"s2\">&quot;markata.plugins.setup_logging&quot;</span><span class=\"p\">,</span>\n<span
    class=\"w\">   </span><span class=\"p\">]</span>\n\n<span class=\"c1\"># point
    log template to the path of your logging template</span>\n<span class=\"n\">log_template</span><span
    class=\"o\">=</span><span class=\"s1\">&#39;templates/log_template.html&#39;</span>\n</pre></div>\n\n</pre>\n\n<p>You
    can see the latest default <code>log_template</code> on\n<a href=\"https://github.com/WaylonWalker/markata/blob/main/markata/plugins/default_log_template.html\">GitHub</a></p>\n<h1
    id=\"disable-logging\">Disable Logging <a class=\"header-anchor\" href=\"#disable-logging\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>If you do not want logging,
    you can explicityly disable it by adding it to your\n<code>[markata.disabled_hooks]</code>
    array in your <code>[markata.toml]</code></p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n\n<span
    class=\"c1\"># make sure its in your list of hooks</span>\n<span class=\"n\">disabled_hooks</span><span
    class=\"o\">=</span><span class=\"p\">[</span>\n<span class=\"w\">   </span><span
    class=\"s2\">&quot;markata.plugins.setup_logging&quot;</span><span class=\"p\">,</span>\n<span
    class=\"w\">   </span><span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<hr
    />\n<div class=\"admonition function\">\n<p class=\"admonition-title\">Function</p>\n<h2
    id=\"has_rich_handler\" class=\"admonition-title\" style=\"margin: 0; padding:
    .5rem 1rem;\">has_rich_handler <em class=\"small\">function</em></h2>\n<p>Returns
    a boolean whether or not there is a RichHandler attached to the\nroot logger.</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">has_rich_handler
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
    class=\"w\"> </span><span class=\"nf\">has_rich_handler</span><span class=\"p\">()</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">bool</span><span class=\"p\">:</span>\n<span
    class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span class=\"sd\">
    \   Returns a boolean whether or not there is a RichHandler attached to the</span>\n<span
    class=\"sd\">    root logger.</span>\n<span class=\"sd\">    &quot;&quot;&quot;</span>\n\n
    \   <span class=\"n\">logger</span> <span class=\"o\">=</span> <span class=\"n\">logging</span><span
    class=\"o\">.</span><span class=\"n\">getLogger</span><span class=\"p\">()</span>\n
    \   <span class=\"k\">return</span> <span class=\"nb\">bool</span><span class=\"p\">([</span><span
    class=\"n\">h</span> <span class=\"k\">for</span> <span class=\"n\">h</span> <span
    class=\"ow\">in</span> <span class=\"n\">logger</span><span class=\"o\">.</span><span
    class=\"n\">handlers</span> <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">h</span><span class=\"p\">,</span> <span
    class=\"n\">RichHandler</span><span class=\"p\">)])</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition function\">\n<p class=\"admonition-title\">Function</p>\n<h2
    id=\"setup_text_log\" class=\"admonition-title\" style=\"margin: 0; padding: .5rem
    1rem;\">setup_text_log <em class=\"small\">function</em></h2>\n<p>sets up a plain
    text log in markata's configured <code>log_dir</code>, or\n<code>output_dir/_logs</code>
    if <code>log_dir</code> is not configured.  The log file will be\nnamed after
    the <code>&lt;levelname&gt;.log</code></p>\n</div>\n<div class=\"admonition source
    is-collapsible collapsible-open\">\n<p class=\"admonition-title\">setup_text_log
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
    class=\"w\"> </span><span class=\"nf\">setup_text_log</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">level</span><span class=\"p\">:</span>
    <span class=\"nb\">int</span> <span class=\"o\">=</span> <span class=\"n\">logging</span><span
    class=\"o\">.</span><span class=\"n\">INFO</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Path</span><span class=\"p\">:</span>\n<span
    class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span class=\"sd\">
    \   sets up a plain text log in markata&#39;s configured `log_dir`, or</span>\n<span
    class=\"sd\">    `output_dir/_logs` if `log_dir` is not configured.  The log file
    will be</span>\n<span class=\"sd\">    named after the `&lt;levelname&gt;.log`</span>\n<span
    class=\"sd\">    &quot;&quot;&quot;</span>\n    <span class=\"n\">log_file</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">logging</span><span
    class=\"o\">.</span><span class=\"n\">log_dir</span> <span class=\"o\">/</span>
    <span class=\"p\">(</span>\n        <span class=\"n\">logging</span><span class=\"o\">.</span><span
    class=\"n\">getLevelName</span><span class=\"p\">(</span><span class=\"n\">level</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">lower</span><span
    class=\"p\">()</span> <span class=\"o\">+</span> <span class=\"s2\">&quot;.log&quot;</span>\n
    \   <span class=\"p\">)</span>\n\n    <span class=\"k\">if</span> <span class=\"n\">has_file_handler</span><span
    class=\"p\">(</span><span class=\"n\">log_file</span><span class=\"p\">):</span>\n
    \       <span class=\"k\">return</span> <span class=\"n\">log_file</span>\n\n
    \   <span class=\"k\">if</span> <span class=\"ow\">not</span> <span class=\"n\">log_file</span><span
    class=\"o\">.</span><span class=\"n\">parent</span><span class=\"o\">.</span><span
    class=\"n\">exists</span><span class=\"p\">():</span>\n        <span class=\"n\">log_file</span><span
    class=\"o\">.</span><span class=\"n\">parent</span><span class=\"o\">.</span><span
    class=\"n\">mkdir</span><span class=\"p\">(</span><span class=\"n\">parents</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \   <span class=\"n\">fh</span> <span class=\"o\">=</span> <span class=\"n\">logging</span><span
    class=\"o\">.</span><span class=\"n\">FileHandler</span><span class=\"p\">(</span><span
    class=\"n\">log_file</span><span class=\"p\">)</span>\n    <span class=\"n\">fh</span><span
    class=\"o\">.</span><span class=\"n\">setLevel</span><span class=\"p\">(</span><span
    class=\"n\">level</span><span class=\"p\">)</span>\n    <span class=\"n\">fh_formatter</span>
    <span class=\"o\">=</span> <span class=\"n\">logging</span><span class=\"o\">.</span><span
    class=\"n\">Formatter</span><span class=\"p\">(</span>\n        <span class=\"s2\">&quot;</span><span
    class=\"si\">%(asctime)s</span><span class=\"s2\"> </span><span class=\"si\">%(name)-12s</span><span
    class=\"s2\"> </span><span class=\"si\">%(levelname)-8s</span><span class=\"s2\">
    </span><span class=\"si\">%(message)s</span><span class=\"s2\">&quot;</span><span
    class=\"p\">,</span>\n    <span class=\"p\">)</span>\n    <span class=\"n\">fh</span><span
    class=\"o\">.</span><span class=\"n\">setFormatter</span><span class=\"p\">(</span><span
    class=\"n\">fh_formatter</span><span class=\"p\">)</span>\n    <span class=\"n\">logging</span><span
    class=\"o\">.</span><span class=\"n\">getLogger</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;&quot;</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">addHandler</span><span class=\"p\">(</span><span class=\"n\">fh</span><span
    class=\"p\">)</span>\n\n    <span class=\"k\">return</span> <span class=\"n\">log_file</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition function\">\n<p class=\"admonition-title\">Function</p>\n<h2
    id=\"setup_html_log\" class=\"admonition-title\" style=\"margin: 0; padding: .5rem
    1rem;\">setup_html_log <em class=\"small\">function</em></h2>\n<p>sets up an html
    log in markata's configured <code>log_dir</code>, or\n<code>output_dir/_logs</code>
    if <code>log_dir</code> is not configured.  The log file will be\nnamed after
    the <code>&lt;levelname&gt;/index.html</code>.  The goal of this is to give</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">setup_html_log
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
    class=\"w\"> </span><span class=\"nf\">setup_html_log</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">level</span><span class=\"p\">:</span>
    <span class=\"nb\">int</span> <span class=\"o\">=</span> <span class=\"n\">logging</span><span
    class=\"o\">.</span><span class=\"n\">INFO</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Path</span><span class=\"p\">:</span>\n<span
    class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span class=\"sd\">
    \   sets up an html log in markata&#39;s configured `log_dir`, or</span>\n<span
    class=\"sd\">    `output_dir/_logs` if `log_dir` is not configured.  The log file
    will be</span>\n<span class=\"sd\">    named after the `&lt;levelname&gt;/index.html`.
    \ The goal of this is to give</span>\n<span class=\"sd\">    &quot;&quot;&quot;</span>\n\n
    \   <span class=\"n\">log_file</span> <span class=\"o\">=</span> <span class=\"p\">(</span>\n
    \       <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">logging</span><span class=\"o\">.</span><span
    class=\"n\">log_dir</span>\n        <span class=\"o\">/</span> <span class=\"n\">logging</span><span
    class=\"o\">.</span><span class=\"n\">getLevelName</span><span class=\"p\">(</span><span
    class=\"n\">level</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">lower</span><span class=\"p\">()</span>\n        <span class=\"o\">/</span>
    <span class=\"s2\">&quot;index.html&quot;</span>\n    <span class=\"p\">)</span>\n\n
    \   <span class=\"k\">if</span> <span class=\"n\">has_file_handler</span><span
    class=\"p\">(</span><span class=\"n\">log_file</span><span class=\"p\">):</span>\n
    \       <span class=\"k\">return</span> <span class=\"n\">log_file</span>\n\n
    \   <span class=\"n\">log_file</span><span class=\"o\">.</span><span class=\"n\">parent</span><span
    class=\"o\">.</span><span class=\"n\">mkdir</span><span class=\"p\">(</span><span
    class=\"n\">parents</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span> <span class=\"n\">exist_ok</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n\n    <span class=\"k\">if</span>
    <span class=\"ow\">not</span> <span class=\"n\">log_file</span><span class=\"o\">.</span><span
    class=\"n\">exists</span><span class=\"p\">():</span>\n        <span class=\"n\">template</span>
    <span class=\"o\">=</span> <span class=\"n\">Template</span><span class=\"p\">(</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">logging</span><span class=\"o\">.</span><span
    class=\"n\">template</span><span class=\"o\">.</span><span class=\"n\">read_text</span><span
    class=\"p\">(),</span> <span class=\"n\">undefined</span><span class=\"o\">=</span><span
    class=\"n\">SilentUndefined</span>\n        <span class=\"p\">)</span>\n        <span
    class=\"n\">log_header</span> <span class=\"o\">=</span> <span class=\"n\">template</span><span
    class=\"o\">.</span><span class=\"n\">render</span><span class=\"p\">(</span>\n
    \           <span class=\"n\">title</span><span class=\"o\">=</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">title</span> <span class=\"o\">+</span> <span class=\"s2\">&quot;
    logs&quot;</span><span class=\"p\">,</span>\n            <span class=\"n\">config</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"p\">,</span>\n        <span class=\"p\">)</span>\n
    \       <span class=\"n\">log_file</span><span class=\"o\">.</span><span class=\"n\">write_text</span><span
    class=\"p\">(</span><span class=\"n\">log_header</span><span class=\"p\">)</span>\n
    \   <span class=\"k\">with</span> <span class=\"nb\">open</span><span class=\"p\">(</span><span
    class=\"n\">log_file</span><span class=\"p\">,</span> <span class=\"s2\">&quot;a&quot;</span><span
    class=\"p\">)</span> <span class=\"k\">as</span> <span class=\"n\">f</span><span
    class=\"p\">:</span>\n        <span class=\"n\">command</span> <span class=\"o\">=</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">sys</span><span
    class=\"o\">.</span><span class=\"n\">argv</span><span class=\"p\">[</span><span
    class=\"mi\">0</span><span class=\"p\">])</span><span class=\"o\">.</span><span
    class=\"n\">name</span> <span class=\"o\">+</span> <span class=\"s2\">&quot; &quot;</span>
    <span class=\"o\">+</span> <span class=\"s2\">&quot; &quot;</span><span class=\"o\">.</span><span
    class=\"n\">join</span><span class=\"p\">(</span><span class=\"n\">sys</span><span
    class=\"o\">.</span><span class=\"n\">argv</span><span class=\"p\">[</span><span
    class=\"mi\">1</span><span class=\"p\">:])</span>\n        <span class=\"n\">f</span><span
    class=\"o\">.</span><span class=\"n\">write</span><span class=\"p\">(</span>\n
    \           <span class=\"sa\">f</span><span class=\"s2\">&quot;&quot;&quot;</span>\n<span
    class=\"s2\">            &lt;div style=&quot;</span>\n<span class=\"s2\">            width:
    100%;</span>\n<span class=\"s2\">            height: 20px;</span>\n<span class=\"s2\">
    \           margin-top: 5rem;</span>\n<span class=\"s2\">            border-bottom:
    1px solid goldenrod;</span>\n<span class=\"s2\">            text-align: center&quot;&gt;</span>\n<span
    class=\"s2\">                &lt;span style=&quot;padding: 0 10px;&quot;&gt;</span>\n<span
    class=\"s2\">                    </span><span class=\"si\">{</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">now</span><span class=\"p\">()</span><span class=\"si\">}</span><span
    class=\"s2\"> running &quot;</span><span class=\"si\">{</span><span class=\"n\">command</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span>\n<span class=\"s2\">                &lt;/span&gt;</span>\n<span
    class=\"s2\">            &lt;/div&gt;</span>\n<span class=\"s2\">    &quot;&quot;&quot;</span><span
    class=\"p\">,</span>\n        <span class=\"p\">)</span>\n    <span class=\"n\">fh</span>
    <span class=\"o\">=</span> <span class=\"n\">logging</span><span class=\"o\">.</span><span
    class=\"n\">FileHandler</span><span class=\"p\">(</span><span class=\"n\">log_file</span><span
    class=\"p\">)</span>\n    <span class=\"n\">fh</span><span class=\"o\">.</span><span
    class=\"n\">setLevel</span><span class=\"p\">(</span><span class=\"n\">level</span><span
    class=\"p\">)</span>\n    <span class=\"n\">fh_formatter</span> <span class=\"o\">=</span>
    <span class=\"n\">logging</span><span class=\"o\">.</span><span class=\"n\">Formatter</span><span
    class=\"p\">(</span>\n<span class=\"w\">        </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">        &lt;li&gt;</span>\n<span class=\"sd\">            &lt;p&gt;</span>\n<span
    class=\"sd\">                &lt;span class=&quot;time&quot;&gt;%(asctime)s&lt;/span&gt;</span>\n<span
    class=\"sd\">                &lt;span class=&quot;name %(name)s&quot;&gt;%(name)-12s&lt;/span&gt;</span>\n<span
    class=\"sd\">                &lt;span class=&quot;levelname %(levelname)s&quot;&gt;%(levelname)-8s&lt;/span&gt;</span>\n<span
    class=\"sd\">            &lt;/p&gt;</span>\n<span class=\"sd\">            &lt;p
    class=&quot;message&quot;&gt;%(message)s&lt;/p&gt;</span>\n<span class=\"sd\">
    \       &lt;/li&gt;</span>\n<span class=\"sd\">        &quot;&quot;&quot;</span><span
    class=\"p\">,</span>\n    <span class=\"p\">)</span>\n    <span class=\"n\">fh</span><span
    class=\"o\">.</span><span class=\"n\">setFormatter</span><span class=\"p\">(</span><span
    class=\"n\">fh_formatter</span><span class=\"p\">)</span>\n    <span class=\"n\">logging</span><span
    class=\"o\">.</span><span class=\"n\">getLogger</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;&quot;</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">addHandler</span><span class=\"p\">(</span><span class=\"n\">fh</span><span
    class=\"p\">)</span>\n\n    <span class=\"k\">return</span> <span class=\"n\">log_file</span>\n</pre></div>\n\n</pre>\n\n</div>\n\n
    \   </section>\n</article>    </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>setup_logging.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Setup Logging hook sets up the RichHandler
    for pretty console logs, and file logs to the configured markata's configured
    \ , or   if  is not configured.  The\u2026\" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n
    \       <meta property=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    <head>\n        <script src=\"https://unpkg.com/@tailwindcss/browser@4\"></script>\n<title>setup_logging.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Setup Logging hook sets up the RichHandler
    for pretty console logs, and file logs to the configured markata's configured
    \ , or   if  is not configured.  The\u2026\" />\n <link href=\"/favicon.ico\"
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
    class=\"title\">\n    <h1 id=\"title\">\n        setup_logging.py\n    </h1>\n</section></article>\n
    \    </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       setup_logging.py\n    </h1>\n</section>    <section class=\"body\">\n
    \       <hr />\n<p>Setup Logging hook sets up the RichHandler for pretty console
    logs, and file\nlogs to the configured markata's configured <code>log_dir</code>,
    or <code>output_dir/_logs</code> if\n<code>log_dir</code> is not configured.  The
    log file will be named after the\n<code>&lt;levelname&gt;.log</code></p>\n<h1
    id=\"the-log-files\">The log files <a class=\"header-anchor\" href=\"#the-log-files\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>There will be 6 log
    files created based on log level and file type.</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span>markout/_logs\n\u251C\u2500\u2500
    debug\n\u2502   \u2514\u2500\u2500 index.html\n\u251C\u2500\u2500 debug.log\n\u251C\u2500\u2500
    info\n\u2502   \u2514\u2500\u2500 index.html\n\u251C\u2500\u2500 info.log\n\u251C\u2500\u2500
    warning\n\u2502   \u2514\u2500\u2500 index.html\n\u2514\u2500\u2500 warning.log\n</pre></div>\n\n</pre>\n\n<h2
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Ensure that setup_logging
    is in your hooks.  You can check if <code>setup_logging</code>\nis in your hooks
    by running <code>markata list --hooks</code> from your terminal and\nchecking
    the output, or creating an instance of <code>Markata()</code> and checking the\n<code>Markata().hooks</code>
    attribute.  If its missing or you wan to be more explicit,\nyou can add <code>setup_logging</code>
    to your <code>markata.toml</code> <code>[markata.hooks]</code>.</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n\n<span
    class=\"c1\"># make sure its in your list of hooks</span>\n<span class=\"n\">hooks</span><span
    class=\"o\">=</span><span class=\"p\">[</span>\n<span class=\"w\">   </span><span
    class=\"s2\">&quot;markata.plugins.setup_logging&quot;</span><span class=\"p\">,</span>\n<span
    class=\"w\">   </span><span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<h1
    id=\"log-template\">Log Template <a class=\"header-anchor\" href=\"#log-template\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n\n<span
    class=\"c1\"># make sure its in your list of hooks</span>\n<span class=\"n\">hooks</span><span
    class=\"o\">=</span><span class=\"p\">[</span>\n<span class=\"w\">   </span><span
    class=\"s2\">&quot;markata.plugins.setup_logging&quot;</span><span class=\"p\">,</span>\n<span
    class=\"w\">   </span><span class=\"p\">]</span>\n\n<span class=\"c1\"># point
    log template to the path of your logging template</span>\n<span class=\"n\">log_template</span><span
    class=\"o\">=</span><span class=\"s1\">&#39;templates/log_template.html&#39;</span>\n</pre></div>\n\n</pre>\n\n<p>You
    can see the latest default <code>log_template</code> on\n<a href=\"https://github.com/WaylonWalker/markata/blob/main/markata/plugins/default_log_template.html\">GitHub</a></p>\n<h1
    id=\"disable-logging\">Disable Logging <a class=\"header-anchor\" href=\"#disable-logging\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>If you do not want logging,
    you can explicityly disable it by adding it to your\n<code>[markata.disabled_hooks]</code>
    array in your <code>[markata.toml]</code></p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n\n<span
    class=\"c1\"># make sure its in your list of hooks</span>\n<span class=\"n\">disabled_hooks</span><span
    class=\"o\">=</span><span class=\"p\">[</span>\n<span class=\"w\">   </span><span
    class=\"s2\">&quot;markata.plugins.setup_logging&quot;</span><span class=\"p\">,</span>\n<span
    class=\"w\">   </span><span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<hr
    />\n<div class=\"admonition function\">\n<p class=\"admonition-title\">Function</p>\n<h2
    id=\"has_rich_handler\" class=\"admonition-title\" style=\"margin: 0; padding:
    .5rem 1rem;\">has_rich_handler <em class=\"small\">function</em></h2>\n<p>Returns
    a boolean whether or not there is a RichHandler attached to the\nroot logger.</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">has_rich_handler
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
    class=\"w\"> </span><span class=\"nf\">has_rich_handler</span><span class=\"p\">()</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">bool</span><span class=\"p\">:</span>\n<span
    class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span class=\"sd\">
    \   Returns a boolean whether or not there is a RichHandler attached to the</span>\n<span
    class=\"sd\">    root logger.</span>\n<span class=\"sd\">    &quot;&quot;&quot;</span>\n\n
    \   <span class=\"n\">logger</span> <span class=\"o\">=</span> <span class=\"n\">logging</span><span
    class=\"o\">.</span><span class=\"n\">getLogger</span><span class=\"p\">()</span>\n
    \   <span class=\"k\">return</span> <span class=\"nb\">bool</span><span class=\"p\">([</span><span
    class=\"n\">h</span> <span class=\"k\">for</span> <span class=\"n\">h</span> <span
    class=\"ow\">in</span> <span class=\"n\">logger</span><span class=\"o\">.</span><span
    class=\"n\">handlers</span> <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">h</span><span class=\"p\">,</span> <span
    class=\"n\">RichHandler</span><span class=\"p\">)])</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition function\">\n<p class=\"admonition-title\">Function</p>\n<h2
    id=\"setup_text_log\" class=\"admonition-title\" style=\"margin: 0; padding: .5rem
    1rem;\">setup_text_log <em class=\"small\">function</em></h2>\n<p>sets up a plain
    text log in markata's configured <code>log_dir</code>, or\n<code>output_dir/_logs</code>
    if <code>log_dir</code> is not configured.  The log file will be\nnamed after
    the <code>&lt;levelname&gt;.log</code></p>\n</div>\n<div class=\"admonition source
    is-collapsible collapsible-open\">\n<p class=\"admonition-title\">setup_text_log
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
    class=\"w\"> </span><span class=\"nf\">setup_text_log</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">level</span><span class=\"p\">:</span>
    <span class=\"nb\">int</span> <span class=\"o\">=</span> <span class=\"n\">logging</span><span
    class=\"o\">.</span><span class=\"n\">INFO</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Path</span><span class=\"p\">:</span>\n<span
    class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span class=\"sd\">
    \   sets up a plain text log in markata&#39;s configured `log_dir`, or</span>\n<span
    class=\"sd\">    `output_dir/_logs` if `log_dir` is not configured.  The log file
    will be</span>\n<span class=\"sd\">    named after the `&lt;levelname&gt;.log`</span>\n<span
    class=\"sd\">    &quot;&quot;&quot;</span>\n    <span class=\"n\">log_file</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">logging</span><span
    class=\"o\">.</span><span class=\"n\">log_dir</span> <span class=\"o\">/</span>
    <span class=\"p\">(</span>\n        <span class=\"n\">logging</span><span class=\"o\">.</span><span
    class=\"n\">getLevelName</span><span class=\"p\">(</span><span class=\"n\">level</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">lower</span><span
    class=\"p\">()</span> <span class=\"o\">+</span> <span class=\"s2\">&quot;.log&quot;</span>\n
    \   <span class=\"p\">)</span>\n\n    <span class=\"k\">if</span> <span class=\"n\">has_file_handler</span><span
    class=\"p\">(</span><span class=\"n\">log_file</span><span class=\"p\">):</span>\n
    \       <span class=\"k\">return</span> <span class=\"n\">log_file</span>\n\n
    \   <span class=\"k\">if</span> <span class=\"ow\">not</span> <span class=\"n\">log_file</span><span
    class=\"o\">.</span><span class=\"n\">parent</span><span class=\"o\">.</span><span
    class=\"n\">exists</span><span class=\"p\">():</span>\n        <span class=\"n\">log_file</span><span
    class=\"o\">.</span><span class=\"n\">parent</span><span class=\"o\">.</span><span
    class=\"n\">mkdir</span><span class=\"p\">(</span><span class=\"n\">parents</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \   <span class=\"n\">fh</span> <span class=\"o\">=</span> <span class=\"n\">logging</span><span
    class=\"o\">.</span><span class=\"n\">FileHandler</span><span class=\"p\">(</span><span
    class=\"n\">log_file</span><span class=\"p\">)</span>\n    <span class=\"n\">fh</span><span
    class=\"o\">.</span><span class=\"n\">setLevel</span><span class=\"p\">(</span><span
    class=\"n\">level</span><span class=\"p\">)</span>\n    <span class=\"n\">fh_formatter</span>
    <span class=\"o\">=</span> <span class=\"n\">logging</span><span class=\"o\">.</span><span
    class=\"n\">Formatter</span><span class=\"p\">(</span>\n        <span class=\"s2\">&quot;</span><span
    class=\"si\">%(asctime)s</span><span class=\"s2\"> </span><span class=\"si\">%(name)-12s</span><span
    class=\"s2\"> </span><span class=\"si\">%(levelname)-8s</span><span class=\"s2\">
    </span><span class=\"si\">%(message)s</span><span class=\"s2\">&quot;</span><span
    class=\"p\">,</span>\n    <span class=\"p\">)</span>\n    <span class=\"n\">fh</span><span
    class=\"o\">.</span><span class=\"n\">setFormatter</span><span class=\"p\">(</span><span
    class=\"n\">fh_formatter</span><span class=\"p\">)</span>\n    <span class=\"n\">logging</span><span
    class=\"o\">.</span><span class=\"n\">getLogger</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;&quot;</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">addHandler</span><span class=\"p\">(</span><span class=\"n\">fh</span><span
    class=\"p\">)</span>\n\n    <span class=\"k\">return</span> <span class=\"n\">log_file</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition function\">\n<p class=\"admonition-title\">Function</p>\n<h2
    id=\"setup_html_log\" class=\"admonition-title\" style=\"margin: 0; padding: .5rem
    1rem;\">setup_html_log <em class=\"small\">function</em></h2>\n<p>sets up an html
    log in markata's configured <code>log_dir</code>, or\n<code>output_dir/_logs</code>
    if <code>log_dir</code> is not configured.  The log file will be\nnamed after
    the <code>&lt;levelname&gt;/index.html</code>.  The goal of this is to give</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">setup_html_log
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
    class=\"w\"> </span><span class=\"nf\">setup_html_log</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">level</span><span class=\"p\">:</span>
    <span class=\"nb\">int</span> <span class=\"o\">=</span> <span class=\"n\">logging</span><span
    class=\"o\">.</span><span class=\"n\">INFO</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Path</span><span class=\"p\">:</span>\n<span
    class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span class=\"sd\">
    \   sets up an html log in markata&#39;s configured `log_dir`, or</span>\n<span
    class=\"sd\">    `output_dir/_logs` if `log_dir` is not configured.  The log file
    will be</span>\n<span class=\"sd\">    named after the `&lt;levelname&gt;/index.html`.
    \ The goal of this is to give</span>\n<span class=\"sd\">    &quot;&quot;&quot;</span>\n\n
    \   <span class=\"n\">log_file</span> <span class=\"o\">=</span> <span class=\"p\">(</span>\n
    \       <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">logging</span><span class=\"o\">.</span><span
    class=\"n\">log_dir</span>\n        <span class=\"o\">/</span> <span class=\"n\">logging</span><span
    class=\"o\">.</span><span class=\"n\">getLevelName</span><span class=\"p\">(</span><span
    class=\"n\">level</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">lower</span><span class=\"p\">()</span>\n        <span class=\"o\">/</span>
    <span class=\"s2\">&quot;index.html&quot;</span>\n    <span class=\"p\">)</span>\n\n
    \   <span class=\"k\">if</span> <span class=\"n\">has_file_handler</span><span
    class=\"p\">(</span><span class=\"n\">log_file</span><span class=\"p\">):</span>\n
    \       <span class=\"k\">return</span> <span class=\"n\">log_file</span>\n\n
    \   <span class=\"n\">log_file</span><span class=\"o\">.</span><span class=\"n\">parent</span><span
    class=\"o\">.</span><span class=\"n\">mkdir</span><span class=\"p\">(</span><span
    class=\"n\">parents</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span> <span class=\"n\">exist_ok</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n\n    <span class=\"k\">if</span>
    <span class=\"ow\">not</span> <span class=\"n\">log_file</span><span class=\"o\">.</span><span
    class=\"n\">exists</span><span class=\"p\">():</span>\n        <span class=\"n\">template</span>
    <span class=\"o\">=</span> <span class=\"n\">Template</span><span class=\"p\">(</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">logging</span><span class=\"o\">.</span><span
    class=\"n\">template</span><span class=\"o\">.</span><span class=\"n\">read_text</span><span
    class=\"p\">(),</span> <span class=\"n\">undefined</span><span class=\"o\">=</span><span
    class=\"n\">SilentUndefined</span>\n        <span class=\"p\">)</span>\n        <span
    class=\"n\">log_header</span> <span class=\"o\">=</span> <span class=\"n\">template</span><span
    class=\"o\">.</span><span class=\"n\">render</span><span class=\"p\">(</span>\n
    \           <span class=\"n\">title</span><span class=\"o\">=</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">title</span> <span class=\"o\">+</span> <span class=\"s2\">&quot;
    logs&quot;</span><span class=\"p\">,</span>\n            <span class=\"n\">config</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"p\">,</span>\n        <span class=\"p\">)</span>\n
    \       <span class=\"n\">log_file</span><span class=\"o\">.</span><span class=\"n\">write_text</span><span
    class=\"p\">(</span><span class=\"n\">log_header</span><span class=\"p\">)</span>\n
    \   <span class=\"k\">with</span> <span class=\"nb\">open</span><span class=\"p\">(</span><span
    class=\"n\">log_file</span><span class=\"p\">,</span> <span class=\"s2\">&quot;a&quot;</span><span
    class=\"p\">)</span> <span class=\"k\">as</span> <span class=\"n\">f</span><span
    class=\"p\">:</span>\n        <span class=\"n\">command</span> <span class=\"o\">=</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">sys</span><span
    class=\"o\">.</span><span class=\"n\">argv</span><span class=\"p\">[</span><span
    class=\"mi\">0</span><span class=\"p\">])</span><span class=\"o\">.</span><span
    class=\"n\">name</span> <span class=\"o\">+</span> <span class=\"s2\">&quot; &quot;</span>
    <span class=\"o\">+</span> <span class=\"s2\">&quot; &quot;</span><span class=\"o\">.</span><span
    class=\"n\">join</span><span class=\"p\">(</span><span class=\"n\">sys</span><span
    class=\"o\">.</span><span class=\"n\">argv</span><span class=\"p\">[</span><span
    class=\"mi\">1</span><span class=\"p\">:])</span>\n        <span class=\"n\">f</span><span
    class=\"o\">.</span><span class=\"n\">write</span><span class=\"p\">(</span>\n
    \           <span class=\"sa\">f</span><span class=\"s2\">&quot;&quot;&quot;</span>\n<span
    class=\"s2\">            &lt;div style=&quot;</span>\n<span class=\"s2\">            width:
    100%;</span>\n<span class=\"s2\">            height: 20px;</span>\n<span class=\"s2\">
    \           margin-top: 5rem;</span>\n<span class=\"s2\">            border-bottom:
    1px solid goldenrod;</span>\n<span class=\"s2\">            text-align: center&quot;&gt;</span>\n<span
    class=\"s2\">                &lt;span style=&quot;padding: 0 10px;&quot;&gt;</span>\n<span
    class=\"s2\">                    </span><span class=\"si\">{</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">now</span><span class=\"p\">()</span><span class=\"si\">}</span><span
    class=\"s2\"> running &quot;</span><span class=\"si\">{</span><span class=\"n\">command</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span>\n<span class=\"s2\">                &lt;/span&gt;</span>\n<span
    class=\"s2\">            &lt;/div&gt;</span>\n<span class=\"s2\">    &quot;&quot;&quot;</span><span
    class=\"p\">,</span>\n        <span class=\"p\">)</span>\n    <span class=\"n\">fh</span>
    <span class=\"o\">=</span> <span class=\"n\">logging</span><span class=\"o\">.</span><span
    class=\"n\">FileHandler</span><span class=\"p\">(</span><span class=\"n\">log_file</span><span
    class=\"p\">)</span>\n    <span class=\"n\">fh</span><span class=\"o\">.</span><span
    class=\"n\">setLevel</span><span class=\"p\">(</span><span class=\"n\">level</span><span
    class=\"p\">)</span>\n    <span class=\"n\">fh_formatter</span> <span class=\"o\">=</span>
    <span class=\"n\">logging</span><span class=\"o\">.</span><span class=\"n\">Formatter</span><span
    class=\"p\">(</span>\n<span class=\"w\">        </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">        &lt;li&gt;</span>\n<span class=\"sd\">            &lt;p&gt;</span>\n<span
    class=\"sd\">                &lt;span class=&quot;time&quot;&gt;%(asctime)s&lt;/span&gt;</span>\n<span
    class=\"sd\">                &lt;span class=&quot;name %(name)s&quot;&gt;%(name)-12s&lt;/span&gt;</span>\n<span
    class=\"sd\">                &lt;span class=&quot;levelname %(levelname)s&quot;&gt;%(levelname)-8s&lt;/span&gt;</span>\n<span
    class=\"sd\">            &lt;/p&gt;</span>\n<span class=\"sd\">            &lt;p
    class=&quot;message&quot;&gt;%(message)s&lt;/p&gt;</span>\n<span class=\"sd\">
    \       &lt;/li&gt;</span>\n<span class=\"sd\">        &quot;&quot;&quot;</span><span
    class=\"p\">,</span>\n    <span class=\"p\">)</span>\n    <span class=\"n\">fh</span><span
    class=\"o\">.</span><span class=\"n\">setFormatter</span><span class=\"p\">(</span><span
    class=\"n\">fh_formatter</span><span class=\"p\">)</span>\n    <span class=\"n\">logging</span><span
    class=\"o\">.</span><span class=\"n\">getLogger</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;&quot;</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">addHandler</span><span class=\"p\">(</span><span class=\"n\">fh</span><span
    class=\"p\">)</span>\n\n    <span class=\"k\">return</span> <span class=\"n\">log_file</span>\n</pre></div>\n\n</pre>\n\n</div>\n\n
    \   </section>\n</article>"
  raw.md: ''
published: false
slug: markata/plugins/setup-logging
title: setup_logging.py


---

---

Setup Logging hook sets up the RichHandler for pretty console logs, and file
logs to the configured markata's configured `log_dir`, or `output_dir/_logs` if
`log_dir` is not configured.  The log file will be named after the
`<levelname>.log`

# The log files

There will be 6 log files created based on log level and file type.

```
markout/_logs
├── debug
│   └── index.html
├── debug.log
├── info
│   └── index.html
├── info.log
├── warning
│   └── index.html
└── warning.log
```

## Configuration

Ensure that setup_logging is in your hooks.  You can check if `setup_logging`
is in your hooks by running `markata list --hooks` from your terminal and
checking the output, or creating an instance of `Markata()` and checking the
`Markata().hooks` attribute.  If its missing or you wan to be more explicit,
you can add `setup_logging` to your `markata.toml` `[markata.hooks]`.

``` toml
[markata]

# make sure its in your list of hooks
hooks=[
   "markata.plugins.setup_logging",
   ]
```

# Log Template
``` toml
[markata]

# make sure its in your list of hooks
hooks=[
   "markata.plugins.setup_logging",
   ]

# point log template to the path of your logging template
log_template='templates/log_template.html'
```

You can see the latest default `log_template` on
[GitHub](https://github.com/WaylonWalker/markata/blob/main/markata/plugins/default_log_template.html)

# Disable Logging

If you do not want logging, you can explicityly disable it by adding it to your
`[markata.disabled_hooks]` array in your `[markata.toml]`

``` toml
[markata]

# make sure its in your list of hooks
disabled_hooks=[
   "markata.plugins.setup_logging",
   ]
```

---

!!! function
    <h2 id="has_rich_handler" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">has_rich_handler <em class="small">function</em></h2>

    Returns a boolean whether or not there is a RichHandler attached to the
    root logger.

???+ source "has_rich_handler <em class='small'>source</em>"
    ```python
    def has_rich_handler() -> bool:
        """
        Returns a boolean whether or not there is a RichHandler attached to the
        root logger.
        """

        logger = logging.getLogger()
        return bool([h for h in logger.handlers if isinstance(h, RichHandler)])
    ```
!!! function
    <h2 id="setup_text_log" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">setup_text_log <em class="small">function</em></h2>

    sets up a plain text log in markata's configured `log_dir`, or
    `output_dir/_logs` if `log_dir` is not configured.  The log file will be
    named after the `<levelname>.log`

???+ source "setup_text_log <em class='small'>source</em>"
    ```python
    def setup_text_log(markata: "Markata", level: int = logging.INFO) -> Path:
        """
        sets up a plain text log in markata's configured `log_dir`, or
        `output_dir/_logs` if `log_dir` is not configured.  The log file will be
        named after the `<levelname>.log`
        """
        log_file = markata.config.logging.log_dir / (
            logging.getLevelName(level).lower() + ".log"
        )

        if has_file_handler(log_file):
            return log_file

        if not log_file.parent.exists():
            log_file.parent.mkdir(parents=True)
        fh = logging.FileHandler(log_file)
        fh.setLevel(level)
        fh_formatter = logging.Formatter(
            "%(asctime)s %(name)-12s %(levelname)-8s %(message)s",
        )
        fh.setFormatter(fh_formatter)
        logging.getLogger("").addHandler(fh)

        return log_file
    ```
!!! function
    <h2 id="setup_html_log" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">setup_html_log <em class="small">function</em></h2>

    sets up an html log in markata's configured `log_dir`, or
    `output_dir/_logs` if `log_dir` is not configured.  The log file will be
    named after the `<levelname>/index.html`.  The goal of this is to give

???+ source "setup_html_log <em class='small'>source</em>"
    ```python
    def setup_html_log(markata: "Markata", level: int = logging.INFO) -> Path:
        """
        sets up an html log in markata's configured `log_dir`, or
        `output_dir/_logs` if `log_dir` is not configured.  The log file will be
        named after the `<levelname>/index.html`.  The goal of this is to give
        """

        log_file = (
            markata.config.logging.log_dir
            / logging.getLevelName(level).lower()
            / "index.html"
        )

        if has_file_handler(log_file):
            return log_file

        log_file.parent.mkdir(parents=True, exist_ok=True)

        if not log_file.exists():
            template = Template(
                markata.config.logging.template.read_text(), undefined=SilentUndefined
            )
            log_header = template.render(
                title=markata.config.title + " logs",
                config=markata.config,
            )
            log_file.write_text(log_header)
        with open(log_file, "a") as f:
            command = Path(sys.argv[0]).name + " " + " ".join(sys.argv[1:])
            f.write(
                f"""
                <div style="
                width: 100%;
                height: 20px;
                margin-top: 5rem;
                border-bottom: 1px solid goldenrod;
                text-align: center">
                    <span style="padding: 0 10px;">
                        {datetime.datetime.now()} running "{command}"
                    </span>
                </div>
        """,
            )
        fh = logging.FileHandler(log_file)
        fh.setLevel(level)
        fh_formatter = logging.Formatter(
            """
            <li>
                <p>
                    <span class="time">%(asctime)s</span>
                    <span class="name %(name)s">%(name)-12s</span>
                    <span class="levelname %(levelname)s">%(levelname)-8s</span>
                </p>
                <p class="message">%(message)s</p>
            </li>
            """,
        )
        fh.setFormatter(fh_formatter)
        logging.getLogger("").addHandler(fh)

        return log_file
    ```