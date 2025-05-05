---
content: "---\n\n[DEPRECATED] The `markata.plugins.sitemap` plugin is deprecated and
  will be removed in a\nfuture version. Please use `markata.plugins.feeds` instead,
  which provides more\ncomprehensive sitemap generation capabilities.\n\n## Installation\n\nThis
  plugin is deprecated. Use `markata.plugins.feeds` instead:\n\n```toml\nhooks = [\n
  \   \"markata.plugins.feeds\",  # Use this instead\n    # \"markata.plugins.sitemap\",
  \ # Deprecated\n]\n```\n\n# Migration Guide\n\nTo migrate to the new feeds plugin:\n\n1.
  Remove sitemap plugin from hooks:\n```toml\n# Remove or comment out\n# \"markata.plugins.sitemap\"\n```\n\n2.
  Add feeds plugin:\n```toml\nhooks = [\n    \"markata.plugins.feeds\"\n]\n```\n\n3.
  Update configuration:\n```toml\n[markata.feeds]\n# Sitemap configuration\nsitemap
  = { output = \"sitemap.xml\" }\n\n# Optional: Configure sitemap settings\n[markata.feeds.sitemap.options]\nchangefreq
  = \"daily\"\npriority = \"0.7\"\n```\n\nSee the feeds plugin documentation for more
  configuration options.\n\n# Legacy Configuration\n\nIf you must continue using this
  plugin temporarily, configure in `markata.toml`:\n\n```toml\n[markata]\nurl = \"https://example.com\"\n\n[markata.sitemap]\nchangefreq
  = \"daily\"\npriority = \"0.7\"\n```\n\n# Dependencies\n\nThis plugin depends on:\n-
  pydantic for configuration\n\nWARNING: This plugin is deprecated and will be removed
  in a future version.\nPlease migrate to `markata.plugins.feeds` as soon as possible.\n\n---\n\n!!!
  class\n    <h2 id=\"SiteMapUrl\" class=\"admonition-title\" style=\"margin: 0; padding:
  .5rem 1rem;\">SiteMapUrl <em class=\"small\">class</em></h2>\n\n    [DEPRECATED]
  A model representing a URL entry in the sitemap.xml file.\n\n    WARNING: This class
  is part of the deprecated sitemap plugin. Please migrate to\n    the feeds plugin
  which provides more comprehensive sitemap generation capabilities.\n\n    To configure
  the base URL for your site, set the 'url' field in your markata config:\n    ```yaml\n
  \   url: https://example.com\n    ```\n\n    If no base URL is set, relative URLs
  will be used.\n\n???+ source \"SiteMapUrl <em class='small'>source</em>\"\n    ```python\n
  \   class SiteMapUrl(pydantic.BaseModel):\n        \"\"\"[DEPRECATED] A model representing
  a URL entry in the sitemap.xml file.\n\n        WARNING: This class is part of the
  deprecated sitemap plugin. Please migrate to\n        the feeds plugin which provides
  more comprehensive sitemap generation capabilities.\n\n        To configure the
  base URL for your site, set the 'url' field in your markata config:\n        ```yaml\n
  \       url: https://example.com\n        ```\n\n        If no base URL is set,
  relative URLs will be used.\n        \"\"\"\n\n        slug: str = Field(..., exclude=True)\n
  \       loc: str = Field(\n            None,\n            include=True,\n            description=\"The
  full URL for this page in the sitemap. Generated automatically from config.url +
  slug.\",\n        )\n        changefreq: str = Field(\"daily\", include=True)\n
  \       priority: str = Field(\"0.7\", include=True)\n        markata: Any = Field(None,
  exclude=True)\n\n        model_config = pydantic.ConfigDict(\n            validate_assignment=False,\n
  \           arbitrary_types_allowed=True,\n            extra=\"allow\",\n            str_strip_whitespace=True,\n
  \           validate_default=True,\n            coerce_numbers_to_str=True,\n            populate_by_name=True,\n
  \       )\n\n        @pydantic.field_validator(\"loc\", mode=\"before\")\n        @classmethod\n
  \       def validate_loc(cls, v, info) -> str:\n            \"\"\"Generate the URL
  for the sitemap entry.\n\n            Uses markata.config.url as the base URL if
  set, otherwise uses relative URLs.\n            Example: https://example.com/my-page/
  or /my-page/\n            \"\"\"\n            if v is None:\n                markata
  = info.data.get(\"markata\")\n                slug = info.data.get(\"slug\")\n                if
  markata is None or slug is None:\n                    raise ValueError(\n                        \"Could
  not generate sitemap URL: markata and slug are required. \"\n                        \"This
  usually means the Post model is missing required fields. \"\n                        \"Check
  that your post has a valid slug and markata instance.\"\n                    )\n\n
  \               # Get base URL from config, default to empty string if not set\n
  \               base_url = getattr(markata.config, \"url\", \"\")\n                if
  not base_url:\n                    return f\"/{slug}/\"\n\n                # Ensure
  URL has a trailing slash for consistency\n                return f\"{base_url.rstrip('/')}/{slug}/\"\n
  \           return v\n\n        def dict(self, *args, **kwargs):\n            return
  {\"url\": {**super().dict(*args, **kwargs)}}\n    ```\n!!! class\n    <h2 id=\"SiteMapPost\"
  class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">SiteMapPost
  <em class=\"small\">class</em></h2>\n\n    [DEPRECATED] A model for posts that will
  be included in the sitemap.\n\n    WARNING: This class is part of the deprecated
  sitemap plugin. Please migrate to\n    the feeds plugin which provides more comprehensive
  sitemap generation capabilities.\n\n    To configure the base URL for your site,
  set the 'url' field in your markata config:\n    ```yaml\n    url: https://example.com\n
  \   ```\n\n???+ source \"SiteMapPost <em class='small'>source</em>\"\n    ```python\n
  \   class SiteMapPost(pydantic.BaseModel):\n        \"\"\"[DEPRECATED] A model for
  posts that will be included in the sitemap.\n\n        WARNING: This class is part
  of the deprecated sitemap plugin. Please migrate to\n        the feeds plugin which
  provides more comprehensive sitemap generation capabilities.\n\n        To configure
  the base URL for your site, set the 'url' field in your markata config:\n        ```yaml\n
  \       url: https://example.com\n        ```\n        \"\"\"\n\n        slug: str
  = None\n        published: bool = True\n        sitemap_url: Optional[SiteMapUrl]
  = None\n        markata: Any = Field(None, exclude=True)\n\n        model_config
  = pydantic.ConfigDict(\n            validate_assignment=False,\n            arbitrary_types_allowed=True,\n
  \           extra=\"allow\",\n            str_strip_whitespace=True,\n            validate_default=True,\n
  \           coerce_numbers_to_str=True,\n            populate_by_name=True,\n        )\n\n
  \       @pydantic.field_validator(\"sitemap_url\", mode=\"before\")\n        @classmethod\n
  \       def validate_sitemap_url(cls, v, info) -> Optional[SiteMapUrl]:\n            \"\"\"Initialize
  sitemap_url if not provided.\"\"\"\n            markata = info.data.get(\"markata\")\n
  \           slug = info.data.get(\"slug\")\n            if markata is None or slug
  is None:\n                raise ValueError(\n                    \"Could not create
  sitemap entry: markata and slug are required. \"\n                    \"This usually
  means the Post model is missing required fields. \"\n                    \"Check
  that your post has a valid slug and markata instance.\"\n                )\n\n            if
  v is None:\n                return SiteMapUrl(markata=markata, slug=slug)\n            if
  isinstance(v, dict):\n                return SiteMapUrl(**v, markata=markata, slug=slug)\n
  \           if v.markata is None:\n                v.markata = markata\n            if
  v.slug is None:\n                v.slug = slug\n            return v\n    ```\n!!!
  method\n    <h2 id=\"validate_loc\" class=\"admonition-title\" style=\"margin: 0;
  padding: .5rem 1rem;\">validate_loc <em class=\"small\">method</em></h2>\n\n    Generate
  the URL for the sitemap entry.\n\n    Uses markata.config.url as the base URL if
  set, otherwise uses relative URLs.\n    Example: https://example.com/my-page/ or
  /my-page/\n\n???+ source \"validate_loc <em class='small'>source</em>\"\n    ```python\n
  \   def validate_loc(cls, v, info) -> str:\n            \"\"\"Generate the URL for
  the sitemap entry.\n\n            Uses markata.config.url as the base URL if set,
  otherwise uses relative URLs.\n            Example: https://example.com/my-page/
  or /my-page/\n            \"\"\"\n            if v is None:\n                markata
  = info.data.get(\"markata\")\n                slug = info.data.get(\"slug\")\n                if
  markata is None or slug is None:\n                    raise ValueError(\n                        \"Could
  not generate sitemap URL: markata and slug are required. \"\n                        \"This
  usually means the Post model is missing required fields. \"\n                        \"Check
  that your post has a valid slug and markata instance.\"\n                    )\n\n
  \               # Get base URL from config, default to empty string if not set\n
  \               base_url = getattr(markata.config, \"url\", \"\")\n                if
  not base_url:\n                    return f\"/{slug}/\"\n\n                # Ensure
  URL has a trailing slash for consistency\n                return f\"{base_url.rstrip('/')}/{slug}/\"\n
  \           return v\n    ```\n!!! method\n    <h2 id=\"validate_sitemap_url\" class=\"admonition-title\"
  style=\"margin: 0; padding: .5rem 1rem;\">validate_sitemap_url <em class=\"small\">method</em></h2>\n\n
  \   Initialize sitemap_url if not provided.\n\n???+ source \"validate_sitemap_url
  <em class='small'>source</em>\"\n    ```python\n    def validate_sitemap_url(cls,
  v, info) -> Optional[SiteMapUrl]:\n            \"\"\"Initialize sitemap_url if not
  provided.\"\"\"\n            markata = info.data.get(\"markata\")\n            slug
  = info.data.get(\"slug\")\n            if markata is None or slug is None:\n                raise
  ValueError(\n                    \"Could not create sitemap entry: markata and slug
  are required. \"\n                    \"This usually means the Post model is missing
  required fields. \"\n                    \"Check that your post has a valid slug
  and markata instance.\"\n                )\n\n            if v is None:\n                return
  SiteMapUrl(markata=markata, slug=slug)\n            if isinstance(v, dict):\n                return
  SiteMapUrl(**v, markata=markata, slug=slug)\n            if v.markata is None:\n
  \               v.markata = markata\n            if v.slug is None:\n                v.slug
  = slug\n            return v\n    ```"
date: 2025-05-05
description: "[DEPRECATED] The   plugin is deprecated and will be removed in a future
  version. Please use   instead, which provides more comprehensive sitemap generation\u2026"
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>sitemap.py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"[DEPRECATED] The   plugin is deprecated and will
    be removed in a future version. Please use   instead, which provides more comprehensive
    sitemap generation\u2026\" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n        <meta property=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    <head>\n        <script src=\"https://unpkg.com/@tailwindcss/browser@4\"></script>\n<title>sitemap.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"[DEPRECATED] The   plugin is deprecated
    and will be removed in a future version. Please use   instead, which provides
    more comprehensive sitemap generation\u2026\" />\n <link href=\"/favicon.ico\"
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        sitemap.py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <hr />\n<p>[DEPRECATED]
    The <code>markata.plugins.sitemap</code> plugin is deprecated and will be removed
    in a\nfuture version. Please use <code>markata.plugins.feeds</code> instead, which
    provides more\ncomprehensive sitemap generation capabilities.</p>\n<h2 id=\"installation\">Installation
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>This plugin is deprecated.
    Use <code>markata.plugins.feeds</code> instead:</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"n\">hooks</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"p\">[</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;markata.plugins.feeds&quot;</span><span
    class=\"p\">,</span><span class=\"w\">  </span><span class=\"c1\"># Use this instead</span>\n<span
    class=\"w\">    </span><span class=\"c1\"># &quot;markata.plugins.sitemap&quot;,
    \ # Deprecated</span>\n<span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<h1
    id=\"migration-guide\">Migration Guide <a class=\"header-anchor\" href=\"#migration-guide\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>To migrate to the new
    feeds plugin:</p>\n<ol>\n<li>Remove sitemap plugin from hooks:</li>\n</ol>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"c1\"># Remove
    or comment out</span>\n<span class=\"c1\"># &quot;markata.plugins.sitemap&quot;</span>\n</pre></div>\n\n</pre>\n\n<ol
    start=\"2\">\n<li>Add feeds plugin:</li>\n</ol>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"n\">hooks</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"p\">[</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;markata.plugins.feeds&quot;</span>\n<span
    class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<ol start=\"3\">\n<li>Update configuration:</li>\n</ol>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata.feeds]</span>\n<span
    class=\"c1\"># Sitemap configuration</span>\n<span class=\"n\">sitemap</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"p\">{</span><span class=\"w\"> </span><span class=\"n\">output</span><span
    class=\"w\"> </span><span class=\"p\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;sitemap.xml&quot;</span><span class=\"w\"> </span><span class=\"p\">}</span>\n\n<span
    class=\"c1\"># Optional: Configure sitemap settings</span>\n<span class=\"k\">[markata.feeds.sitemap.options]</span>\n<span
    class=\"n\">changefreq</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;daily&quot;</span>\n<span class=\"n\">priority</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;0.7&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>See the feeds
    plugin documentation for more configuration options.</p>\n<h1 id=\"legacy-configuration\">Legacy
    Configuration <a class=\"header-anchor\" href=\"#legacy-configuration\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>If you must continue
    using this plugin temporarily, configure in <code>markata.toml</code>:</p>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n<span
    class=\"n\">url</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;https://example.com&quot;</span>\n\n<span
    class=\"k\">[markata.sitemap]</span>\n<span class=\"n\">changefreq</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;daily&quot;</span>\n<span class=\"n\">priority</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;0.7&quot;</span>\n</pre></div>\n\n</pre>\n\n<h1 id=\"dependencies\">Dependencies
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>This plugin depends
    on:</p>\n<ul>\n<li>pydantic for configuration</li>\n</ul>\n<p>WARNING: This plugin
    is deprecated and will be removed in a future version.\nPlease migrate to <code>markata.plugins.feeds</code>
    as soon as possible.</p>\n<hr />\n<div class=\"admonition class\">\n<p class=\"admonition-title\">Class</p>\n<h2
    id=\"SiteMapUrl\" class=\"admonition-title\" style=\"margin: 0; padding: .5rem
    1rem;\">SiteMapUrl <em class=\"small\">class</em></h2>\n<p>[DEPRECATED] A model
    representing a URL entry in the sitemap.xml file.</p>\n<p>WARNING: This class
    is part of the deprecated sitemap plugin. Please migrate to\nthe feeds plugin
    which provides more comprehensive sitemap generation capabilities.</p>\n<p>To
    configure the base URL for your site, set the 'url' field in your markata config:</p>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"nt\">url</span><span
    class=\"p\">:</span><span class=\"w\"> </span><span class=\"l l-Scalar l-Scalar-Plain\">https://example.com</span>\n</pre></div>\n\n</pre>\n\n<p>If
    no base URL is set, relative URLs will be used.</p>\n</div>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">SiteMapUrl
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
    class=\"w\"> </span><span class=\"nc\">SiteMapUrl</span><span class=\"p\">(</span><span
    class=\"n\">pydantic</span><span class=\"o\">.</span><span class=\"n\">BaseModel</span><span
    class=\"p\">):</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;[DEPRECATED]
    A model representing a URL entry in the sitemap.xml file.</span>\n\n<span class=\"sd\">
    \   WARNING: This class is part of the deprecated sitemap plugin. Please migrate
    to</span>\n<span class=\"sd\">    the feeds plugin which provides more comprehensive
    sitemap generation capabilities.</span>\n\n<span class=\"sd\">    To configure
    the base URL for your site, set the &#39;url&#39; field in your markata config:</span>\n<span
    class=\"sd\">    ```yaml</span>\n<span class=\"sd\">    url: https://example.com</span>\n<span
    class=\"sd\">    ```</span>\n\n<span class=\"sd\">    If no base URL is set, relative
    URLs will be used.</span>\n<span class=\"sd\">    &quot;&quot;&quot;</span>\n\n
    \   <span class=\"n\">slug</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"n\">Field</span><span class=\"p\">(</span><span
    class=\"o\">...</span><span class=\"p\">,</span> <span class=\"n\">exclude</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \   <span class=\"n\">loc</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"n\">Field</span><span class=\"p\">(</span>\n
    \       <span class=\"kc\">None</span><span class=\"p\">,</span>\n        <span
    class=\"n\">include</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span>\n        <span class=\"n\">description</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;The full URL for this page in the sitemap. Generated automatically
    from config.url + slug.&quot;</span><span class=\"p\">,</span>\n    <span class=\"p\">)</span>\n
    \   <span class=\"n\">changefreq</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"n\">Field</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;daily&quot;</span><span class=\"p\">,</span> <span class=\"n\">include</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \   <span class=\"n\">priority</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"n\">Field</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;0.7&quot;</span><span class=\"p\">,</span> <span class=\"n\">include</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \   <span class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"n\">Any</span>
    <span class=\"o\">=</span> <span class=\"n\">Field</span><span class=\"p\">(</span><span
    class=\"kc\">None</span><span class=\"p\">,</span> <span class=\"n\">exclude</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n\n
    \   <span class=\"n\">model_config</span> <span class=\"o\">=</span> <span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">ConfigDict</span><span class=\"p\">(</span>\n
    \       <span class=\"n\">validate_assignment</span><span class=\"o\">=</span><span
    class=\"kc\">False</span><span class=\"p\">,</span>\n        <span class=\"n\">arbitrary_types_allowed</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">extra</span><span class=\"o\">=</span><span class=\"s2\">&quot;allow&quot;</span><span
    class=\"p\">,</span>\n        <span class=\"n\">str_strip_whitespace</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">validate_default</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n        <span class=\"n\">coerce_numbers_to_str</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">populate_by_name</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n    <span class=\"p\">)</span>\n\n
    \   <span class=\"nd\">@pydantic</span><span class=\"o\">.</span><span class=\"n\">field_validator</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;loc&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">mode</span><span class=\"o\">=</span><span class=\"s2\">&quot;before&quot;</span><span
    class=\"p\">)</span>\n    <span class=\"nd\">@classmethod</span>\n    <span class=\"k\">def</span><span
    class=\"w\"> </span><span class=\"nf\">validate_loc</span><span class=\"p\">(</span><span
    class=\"bp\">cls</span><span class=\"p\">,</span> <span class=\"n\">v</span><span
    class=\"p\">,</span> <span class=\"n\">info</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n<span
    class=\"w\">        </span><span class=\"sd\">&quot;&quot;&quot;Generate the URL
    for the sitemap entry.</span>\n\n<span class=\"sd\">        Uses markata.config.url
    as the base URL if set, otherwise uses relative URLs.</span>\n<span class=\"sd\">
    \       Example: https://example.com/my-page/ or /my-page/</span>\n<span class=\"sd\">
    \       &quot;&quot;&quot;</span>\n        <span class=\"k\">if</span> <span class=\"n\">v</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">markata</span> <span class=\"o\">=</span> <span
    class=\"n\">info</span><span class=\"o\">.</span><span class=\"n\">data</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;markata&quot;</span><span class=\"p\">)</span>\n            <span
    class=\"n\">slug</span> <span class=\"o\">=</span> <span class=\"n\">info</span><span
    class=\"o\">.</span><span class=\"n\">data</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">)</span>\n            <span class=\"k\">if</span> <span class=\"n\">markata</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span> <span class=\"ow\">or</span>
    <span class=\"n\">slug</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"k\">raise</span> <span class=\"ne\">ValueError</span><span
    class=\"p\">(</span>\n                    <span class=\"s2\">&quot;Could not generate
    sitemap URL: markata and slug are required. &quot;</span>\n                    <span
    class=\"s2\">&quot;This usually means the Post model is missing required fields.
    &quot;</span>\n                    <span class=\"s2\">&quot;Check that your post
    has a valid slug and markata instance.&quot;</span>\n                <span class=\"p\">)</span>\n\n
    \           <span class=\"c1\"># Get base URL from config, default to empty string
    if not set</span>\n            <span class=\"n\">base_url</span> <span class=\"o\">=</span>
    <span class=\"nb\">getattr</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;url&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;&quot;</span><span
    class=\"p\">)</span>\n            <span class=\"k\">if</span> <span class=\"ow\">not</span>
    <span class=\"n\">base_url</span><span class=\"p\">:</span>\n                <span
    class=\"k\">return</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;/</span><span
    class=\"si\">{</span><span class=\"n\">slug</span><span class=\"si\">}</span><span
    class=\"s2\">/&quot;</span>\n\n            <span class=\"c1\"># Ensure URL has
    a trailing slash for consistency</span>\n            <span class=\"k\">return</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">base_url</span><span class=\"o\">.</span><span class=\"n\">rstrip</span><span
    class=\"p\">(</span><span class=\"s1\">&#39;/&#39;</span><span class=\"p\">)</span><span
    class=\"si\">}</span><span class=\"s2\">/</span><span class=\"si\">{</span><span
    class=\"n\">slug</span><span class=\"si\">}</span><span class=\"s2\">/&quot;</span>\n
    \       <span class=\"k\">return</span> <span class=\"n\">v</span>\n\n    <span
    class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">dict</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"n\">args</span><span class=\"p\">,</span> <span
    class=\"o\">**</span><span class=\"n\">kwargs</span><span class=\"p\">):</span>\n
    \       <span class=\"k\">return</span> <span class=\"p\">{</span><span class=\"s2\">&quot;url&quot;</span><span
    class=\"p\">:</span> <span class=\"p\">{</span><span class=\"o\">**</span><span
    class=\"nb\">super</span><span class=\"p\">()</span><span class=\"o\">.</span><span
    class=\"n\">dict</span><span class=\"p\">(</span><span class=\"o\">*</span><span
    class=\"n\">args</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"n\">kwargs</span><span class=\"p\">)}}</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition class\">\n<p class=\"admonition-title\">Class</p>\n<h2 id=\"SiteMapPost\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">SiteMapPost
    <em class=\"small\">class</em></h2>\n<p>[DEPRECATED] A model for posts that will
    be included in the sitemap.</p>\n<p>WARNING: This class is part of the deprecated
    sitemap plugin. Please migrate to\nthe feeds plugin which provides more comprehensive
    sitemap generation capabilities.</p>\n<p>To configure the base URL for your site,
    set the 'url' field in your markata config:</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"nt\">url</span><span
    class=\"p\">:</span><span class=\"w\"> </span><span class=\"l l-Scalar l-Scalar-Plain\">https://example.com</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">SiteMapPost
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
    class=\"w\"> </span><span class=\"nc\">SiteMapPost</span><span class=\"p\">(</span><span
    class=\"n\">pydantic</span><span class=\"o\">.</span><span class=\"n\">BaseModel</span><span
    class=\"p\">):</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;[DEPRECATED]
    A model for posts that will be included in the sitemap.</span>\n\n<span class=\"sd\">
    \   WARNING: This class is part of the deprecated sitemap plugin. Please migrate
    to</span>\n<span class=\"sd\">    the feeds plugin which provides more comprehensive
    sitemap generation capabilities.</span>\n\n<span class=\"sd\">    To configure
    the base URL for your site, set the &#39;url&#39; field in your markata config:</span>\n<span
    class=\"sd\">    ```yaml</span>\n<span class=\"sd\">    url: https://example.com</span>\n<span
    class=\"sd\">    ```</span>\n<span class=\"sd\">    &quot;&quot;&quot;</span>\n\n
    \   <span class=\"n\">slug</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"kc\">None</span>\n    <span class=\"n\">published</span><span
    class=\"p\">:</span> <span class=\"nb\">bool</span> <span class=\"o\">=</span>
    <span class=\"kc\">True</span>\n    <span class=\"n\">sitemap_url</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"n\">SiteMapUrl</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"kc\">None</span>\n    <span class=\"n\">markata</span><span class=\"p\">:</span>
    <span class=\"n\">Any</span> <span class=\"o\">=</span> <span class=\"n\">Field</span><span
    class=\"p\">(</span><span class=\"kc\">None</span><span class=\"p\">,</span> <span
    class=\"n\">exclude</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n\n    <span class=\"n\">model_config</span> <span class=\"o\">=</span>
    <span class=\"n\">pydantic</span><span class=\"o\">.</span><span class=\"n\">ConfigDict</span><span
    class=\"p\">(</span>\n        <span class=\"n\">validate_assignment</span><span
    class=\"o\">=</span><span class=\"kc\">False</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">arbitrary_types_allowed</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n        <span class=\"n\">extra</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;allow&quot;</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">str_strip_whitespace</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n        <span class=\"n\">validate_default</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">coerce_numbers_to_str</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n        <span class=\"n\">populate_by_name</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \   <span class=\"p\">)</span>\n\n    <span class=\"nd\">@pydantic</span><span
    class=\"o\">.</span><span class=\"n\">field_validator</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;sitemap_url&quot;</span><span class=\"p\">,</span> <span class=\"n\">mode</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;before&quot;</span><span class=\"p\">)</span>\n
    \   <span class=\"nd\">@classmethod</span>\n    <span class=\"k\">def</span><span
    class=\"w\"> </span><span class=\"nf\">validate_sitemap_url</span><span class=\"p\">(</span><span
    class=\"bp\">cls</span><span class=\"p\">,</span> <span class=\"n\">v</span><span
    class=\"p\">,</span> <span class=\"n\">info</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"n\">SiteMapUrl</span><span class=\"p\">]:</span>\n<span class=\"w\">        </span><span
    class=\"sd\">&quot;&quot;&quot;Initialize sitemap_url if not provided.&quot;&quot;&quot;</span>\n
    \       <span class=\"n\">markata</span> <span class=\"o\">=</span> <span class=\"n\">info</span><span
    class=\"o\">.</span><span class=\"n\">data</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;markata&quot;</span><span
    class=\"p\">)</span>\n        <span class=\"n\">slug</span> <span class=\"o\">=</span>
    <span class=\"n\">info</span><span class=\"o\">.</span><span class=\"n\">data</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;slug&quot;</span><span class=\"p\">)</span>\n        <span
    class=\"k\">if</span> <span class=\"n\">markata</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span> <span class=\"ow\">or</span> <span class=\"n\">slug</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">raise</span> <span class=\"ne\">ValueError</span><span
    class=\"p\">(</span>\n                <span class=\"s2\">&quot;Could not create
    sitemap entry: markata and slug are required. &quot;</span>\n                <span
    class=\"s2\">&quot;This usually means the Post model is missing required fields.
    &quot;</span>\n                <span class=\"s2\">&quot;Check that your post has
    a valid slug and markata instance.&quot;</span>\n            <span class=\"p\">)</span>\n\n
    \       <span class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"k\">return</span>
    <span class=\"n\">SiteMapUrl</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">slug</span><span class=\"o\">=</span><span class=\"n\">slug</span><span
    class=\"p\">)</span>\n        <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"nb\">dict</span><span class=\"p\">):</span>\n            <span class=\"k\">return</span>
    <span class=\"n\">SiteMapUrl</span><span class=\"p\">(</span><span class=\"o\">**</span><span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">slug</span><span class=\"o\">=</span><span class=\"n\">slug</span><span
    class=\"p\">)</span>\n        <span class=\"k\">if</span> <span class=\"n\">v</span><span
    class=\"o\">.</span><span class=\"n\">markata</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"n\">v</span><span
    class=\"o\">.</span><span class=\"n\">markata</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span>\n        <span class=\"k\">if</span> <span class=\"n\">v</span><span
    class=\"o\">.</span><span class=\"n\">slug</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"n\">v</span><span
    class=\"o\">.</span><span class=\"n\">slug</span> <span class=\"o\">=</span> <span
    class=\"n\">slug</span>\n        <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"validate_loc\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">validate_loc
    <em class=\"small\">method</em></h2>\n<p>Generate the URL for the sitemap entry.</p>\n<p>Uses
    markata.config.url as the base URL if set, otherwise uses relative URLs.\nExample:
    <a href=\"https://example.com/my-page/\">https://example.com/my-page/</a> or /my-page/</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">validate_loc
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
    class=\"w\"> </span><span class=\"nf\">validate_loc</span><span class=\"p\">(</span><span
    class=\"bp\">cls</span><span class=\"p\">,</span> <span class=\"n\">v</span><span
    class=\"p\">,</span> <span class=\"n\">info</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n<span
    class=\"w\">        </span><span class=\"sd\">&quot;&quot;&quot;Generate the URL
    for the sitemap entry.</span>\n\n<span class=\"sd\">        Uses markata.config.url
    as the base URL if set, otherwise uses relative URLs.</span>\n<span class=\"sd\">
    \       Example: https://example.com/my-page/ or /my-page/</span>\n<span class=\"sd\">
    \       &quot;&quot;&quot;</span>\n        <span class=\"k\">if</span> <span class=\"n\">v</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">markata</span> <span class=\"o\">=</span> <span
    class=\"n\">info</span><span class=\"o\">.</span><span class=\"n\">data</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;markata&quot;</span><span class=\"p\">)</span>\n            <span
    class=\"n\">slug</span> <span class=\"o\">=</span> <span class=\"n\">info</span><span
    class=\"o\">.</span><span class=\"n\">data</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">)</span>\n            <span class=\"k\">if</span> <span class=\"n\">markata</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span> <span class=\"ow\">or</span>
    <span class=\"n\">slug</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"k\">raise</span> <span class=\"ne\">ValueError</span><span
    class=\"p\">(</span>\n                    <span class=\"s2\">&quot;Could not generate
    sitemap URL: markata and slug are required. &quot;</span>\n                    <span
    class=\"s2\">&quot;This usually means the Post model is missing required fields.
    &quot;</span>\n                    <span class=\"s2\">&quot;Check that your post
    has a valid slug and markata instance.&quot;</span>\n                <span class=\"p\">)</span>\n\n
    \           <span class=\"c1\"># Get base URL from config, default to empty string
    if not set</span>\n            <span class=\"n\">base_url</span> <span class=\"o\">=</span>
    <span class=\"nb\">getattr</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;url&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;&quot;</span><span
    class=\"p\">)</span>\n            <span class=\"k\">if</span> <span class=\"ow\">not</span>
    <span class=\"n\">base_url</span><span class=\"p\">:</span>\n                <span
    class=\"k\">return</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;/</span><span
    class=\"si\">{</span><span class=\"n\">slug</span><span class=\"si\">}</span><span
    class=\"s2\">/&quot;</span>\n\n            <span class=\"c1\"># Ensure URL has
    a trailing slash for consistency</span>\n            <span class=\"k\">return</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">base_url</span><span class=\"o\">.</span><span class=\"n\">rstrip</span><span
    class=\"p\">(</span><span class=\"s1\">&#39;/&#39;</span><span class=\"p\">)</span><span
    class=\"si\">}</span><span class=\"s2\">/</span><span class=\"si\">{</span><span
    class=\"n\">slug</span><span class=\"si\">}</span><span class=\"s2\">/&quot;</span>\n
    \       <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"validate_sitemap_url\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">validate_sitemap_url
    <em class=\"small\">method</em></h2>\n<p>Initialize sitemap_url if not provided.</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">validate_sitemap_url
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
    class=\"w\"> </span><span class=\"nf\">validate_sitemap_url</span><span class=\"p\">(</span><span
    class=\"bp\">cls</span><span class=\"p\">,</span> <span class=\"n\">v</span><span
    class=\"p\">,</span> <span class=\"n\">info</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"n\">SiteMapUrl</span><span class=\"p\">]:</span>\n<span class=\"w\">        </span><span
    class=\"sd\">&quot;&quot;&quot;Initialize sitemap_url if not provided.&quot;&quot;&quot;</span>\n
    \       <span class=\"n\">markata</span> <span class=\"o\">=</span> <span class=\"n\">info</span><span
    class=\"o\">.</span><span class=\"n\">data</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;markata&quot;</span><span
    class=\"p\">)</span>\n        <span class=\"n\">slug</span> <span class=\"o\">=</span>
    <span class=\"n\">info</span><span class=\"o\">.</span><span class=\"n\">data</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;slug&quot;</span><span class=\"p\">)</span>\n        <span
    class=\"k\">if</span> <span class=\"n\">markata</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span> <span class=\"ow\">or</span> <span class=\"n\">slug</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">raise</span> <span class=\"ne\">ValueError</span><span
    class=\"p\">(</span>\n                <span class=\"s2\">&quot;Could not create
    sitemap entry: markata and slug are required. &quot;</span>\n                <span
    class=\"s2\">&quot;This usually means the Post model is missing required fields.
    &quot;</span>\n                <span class=\"s2\">&quot;Check that your post has
    a valid slug and markata instance.&quot;</span>\n            <span class=\"p\">)</span>\n\n
    \       <span class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"k\">return</span>
    <span class=\"n\">SiteMapUrl</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">slug</span><span class=\"o\">=</span><span class=\"n\">slug</span><span
    class=\"p\">)</span>\n        <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"nb\">dict</span><span class=\"p\">):</span>\n            <span class=\"k\">return</span>
    <span class=\"n\">SiteMapUrl</span><span class=\"p\">(</span><span class=\"o\">**</span><span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">slug</span><span class=\"o\">=</span><span class=\"n\">slug</span><span
    class=\"p\">)</span>\n        <span class=\"k\">if</span> <span class=\"n\">v</span><span
    class=\"o\">.</span><span class=\"n\">markata</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"n\">v</span><span
    class=\"o\">.</span><span class=\"n\">markata</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span>\n        <span class=\"k\">if</span> <span class=\"n\">v</span><span
    class=\"o\">.</span><span class=\"n\">slug</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"n\">v</span><span
    class=\"o\">.</span><span class=\"n\">slug</span> <span class=\"o\">=</span> <span
    class=\"n\">slug</span>\n        <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n</div>\n\n
    \   </section>\n</article>    </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>sitemap.py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"[DEPRECATED] The   plugin is deprecated and will
    be removed in a future version. Please use   instead, which provides more comprehensive
    sitemap generation\u2026\" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n        <meta property=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    <head>\n        <script src=\"https://unpkg.com/@tailwindcss/browser@4\"></script>\n<title>sitemap.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"[DEPRECATED] The   plugin is deprecated
    and will be removed in a future version. Please use   instead, which provides
    more comprehensive sitemap generation\u2026\" />\n <link href=\"/favicon.ico\"
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
    class=\"title\">\n    <h1 id=\"title\">\n        sitemap.py\n    </h1>\n</section></article>\n
    \    </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       sitemap.py\n    </h1>\n</section>    <section class=\"body\">\n        <hr
    />\n<p>[DEPRECATED] The <code>markata.plugins.sitemap</code> plugin is deprecated
    and will be removed in a\nfuture version. Please use <code>markata.plugins.feeds</code>
    instead, which provides more\ncomprehensive sitemap generation capabilities.</p>\n<h2
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>This plugin is deprecated.
    Use <code>markata.plugins.feeds</code> instead:</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"n\">hooks</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"p\">[</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;markata.plugins.feeds&quot;</span><span
    class=\"p\">,</span><span class=\"w\">  </span><span class=\"c1\"># Use this instead</span>\n<span
    class=\"w\">    </span><span class=\"c1\"># &quot;markata.plugins.sitemap&quot;,
    \ # Deprecated</span>\n<span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<h1
    id=\"migration-guide\">Migration Guide <a class=\"header-anchor\" href=\"#migration-guide\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>To migrate to the new
    feeds plugin:</p>\n<ol>\n<li>Remove sitemap plugin from hooks:</li>\n</ol>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"c1\"># Remove
    or comment out</span>\n<span class=\"c1\"># &quot;markata.plugins.sitemap&quot;</span>\n</pre></div>\n\n</pre>\n\n<ol
    start=\"2\">\n<li>Add feeds plugin:</li>\n</ol>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"n\">hooks</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"p\">[</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;markata.plugins.feeds&quot;</span>\n<span
    class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<ol start=\"3\">\n<li>Update configuration:</li>\n</ol>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata.feeds]</span>\n<span
    class=\"c1\"># Sitemap configuration</span>\n<span class=\"n\">sitemap</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"p\">{</span><span class=\"w\"> </span><span class=\"n\">output</span><span
    class=\"w\"> </span><span class=\"p\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;sitemap.xml&quot;</span><span class=\"w\"> </span><span class=\"p\">}</span>\n\n<span
    class=\"c1\"># Optional: Configure sitemap settings</span>\n<span class=\"k\">[markata.feeds.sitemap.options]</span>\n<span
    class=\"n\">changefreq</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;daily&quot;</span>\n<span class=\"n\">priority</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;0.7&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>See the feeds
    plugin documentation for more configuration options.</p>\n<h1 id=\"legacy-configuration\">Legacy
    Configuration <a class=\"header-anchor\" href=\"#legacy-configuration\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>If you must continue
    using this plugin temporarily, configure in <code>markata.toml</code>:</p>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n<span
    class=\"n\">url</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;https://example.com&quot;</span>\n\n<span
    class=\"k\">[markata.sitemap]</span>\n<span class=\"n\">changefreq</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;daily&quot;</span>\n<span class=\"n\">priority</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;0.7&quot;</span>\n</pre></div>\n\n</pre>\n\n<h1 id=\"dependencies\">Dependencies
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>This plugin depends
    on:</p>\n<ul>\n<li>pydantic for configuration</li>\n</ul>\n<p>WARNING: This plugin
    is deprecated and will be removed in a future version.\nPlease migrate to <code>markata.plugins.feeds</code>
    as soon as possible.</p>\n<hr />\n<div class=\"admonition class\">\n<p class=\"admonition-title\">Class</p>\n<h2
    id=\"SiteMapUrl\" class=\"admonition-title\" style=\"margin: 0; padding: .5rem
    1rem;\">SiteMapUrl <em class=\"small\">class</em></h2>\n<p>[DEPRECATED] A model
    representing a URL entry in the sitemap.xml file.</p>\n<p>WARNING: This class
    is part of the deprecated sitemap plugin. Please migrate to\nthe feeds plugin
    which provides more comprehensive sitemap generation capabilities.</p>\n<p>To
    configure the base URL for your site, set the 'url' field in your markata config:</p>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"nt\">url</span><span
    class=\"p\">:</span><span class=\"w\"> </span><span class=\"l l-Scalar l-Scalar-Plain\">https://example.com</span>\n</pre></div>\n\n</pre>\n\n<p>If
    no base URL is set, relative URLs will be used.</p>\n</div>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">SiteMapUrl
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
    class=\"w\"> </span><span class=\"nc\">SiteMapUrl</span><span class=\"p\">(</span><span
    class=\"n\">pydantic</span><span class=\"o\">.</span><span class=\"n\">BaseModel</span><span
    class=\"p\">):</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;[DEPRECATED]
    A model representing a URL entry in the sitemap.xml file.</span>\n\n<span class=\"sd\">
    \   WARNING: This class is part of the deprecated sitemap plugin. Please migrate
    to</span>\n<span class=\"sd\">    the feeds plugin which provides more comprehensive
    sitemap generation capabilities.</span>\n\n<span class=\"sd\">    To configure
    the base URL for your site, set the &#39;url&#39; field in your markata config:</span>\n<span
    class=\"sd\">    ```yaml</span>\n<span class=\"sd\">    url: https://example.com</span>\n<span
    class=\"sd\">    ```</span>\n\n<span class=\"sd\">    If no base URL is set, relative
    URLs will be used.</span>\n<span class=\"sd\">    &quot;&quot;&quot;</span>\n\n
    \   <span class=\"n\">slug</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"n\">Field</span><span class=\"p\">(</span><span
    class=\"o\">...</span><span class=\"p\">,</span> <span class=\"n\">exclude</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \   <span class=\"n\">loc</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"n\">Field</span><span class=\"p\">(</span>\n
    \       <span class=\"kc\">None</span><span class=\"p\">,</span>\n        <span
    class=\"n\">include</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span>\n        <span class=\"n\">description</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;The full URL for this page in the sitemap. Generated automatically
    from config.url + slug.&quot;</span><span class=\"p\">,</span>\n    <span class=\"p\">)</span>\n
    \   <span class=\"n\">changefreq</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"n\">Field</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;daily&quot;</span><span class=\"p\">,</span> <span class=\"n\">include</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \   <span class=\"n\">priority</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"n\">Field</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;0.7&quot;</span><span class=\"p\">,</span> <span class=\"n\">include</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \   <span class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"n\">Any</span>
    <span class=\"o\">=</span> <span class=\"n\">Field</span><span class=\"p\">(</span><span
    class=\"kc\">None</span><span class=\"p\">,</span> <span class=\"n\">exclude</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n\n
    \   <span class=\"n\">model_config</span> <span class=\"o\">=</span> <span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">ConfigDict</span><span class=\"p\">(</span>\n
    \       <span class=\"n\">validate_assignment</span><span class=\"o\">=</span><span
    class=\"kc\">False</span><span class=\"p\">,</span>\n        <span class=\"n\">arbitrary_types_allowed</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">extra</span><span class=\"o\">=</span><span class=\"s2\">&quot;allow&quot;</span><span
    class=\"p\">,</span>\n        <span class=\"n\">str_strip_whitespace</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">validate_default</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n        <span class=\"n\">coerce_numbers_to_str</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">populate_by_name</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n    <span class=\"p\">)</span>\n\n
    \   <span class=\"nd\">@pydantic</span><span class=\"o\">.</span><span class=\"n\">field_validator</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;loc&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">mode</span><span class=\"o\">=</span><span class=\"s2\">&quot;before&quot;</span><span
    class=\"p\">)</span>\n    <span class=\"nd\">@classmethod</span>\n    <span class=\"k\">def</span><span
    class=\"w\"> </span><span class=\"nf\">validate_loc</span><span class=\"p\">(</span><span
    class=\"bp\">cls</span><span class=\"p\">,</span> <span class=\"n\">v</span><span
    class=\"p\">,</span> <span class=\"n\">info</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n<span
    class=\"w\">        </span><span class=\"sd\">&quot;&quot;&quot;Generate the URL
    for the sitemap entry.</span>\n\n<span class=\"sd\">        Uses markata.config.url
    as the base URL if set, otherwise uses relative URLs.</span>\n<span class=\"sd\">
    \       Example: https://example.com/my-page/ or /my-page/</span>\n<span class=\"sd\">
    \       &quot;&quot;&quot;</span>\n        <span class=\"k\">if</span> <span class=\"n\">v</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">markata</span> <span class=\"o\">=</span> <span
    class=\"n\">info</span><span class=\"o\">.</span><span class=\"n\">data</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;markata&quot;</span><span class=\"p\">)</span>\n            <span
    class=\"n\">slug</span> <span class=\"o\">=</span> <span class=\"n\">info</span><span
    class=\"o\">.</span><span class=\"n\">data</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">)</span>\n            <span class=\"k\">if</span> <span class=\"n\">markata</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span> <span class=\"ow\">or</span>
    <span class=\"n\">slug</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"k\">raise</span> <span class=\"ne\">ValueError</span><span
    class=\"p\">(</span>\n                    <span class=\"s2\">&quot;Could not generate
    sitemap URL: markata and slug are required. &quot;</span>\n                    <span
    class=\"s2\">&quot;This usually means the Post model is missing required fields.
    &quot;</span>\n                    <span class=\"s2\">&quot;Check that your post
    has a valid slug and markata instance.&quot;</span>\n                <span class=\"p\">)</span>\n\n
    \           <span class=\"c1\"># Get base URL from config, default to empty string
    if not set</span>\n            <span class=\"n\">base_url</span> <span class=\"o\">=</span>
    <span class=\"nb\">getattr</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;url&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;&quot;</span><span
    class=\"p\">)</span>\n            <span class=\"k\">if</span> <span class=\"ow\">not</span>
    <span class=\"n\">base_url</span><span class=\"p\">:</span>\n                <span
    class=\"k\">return</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;/</span><span
    class=\"si\">{</span><span class=\"n\">slug</span><span class=\"si\">}</span><span
    class=\"s2\">/&quot;</span>\n\n            <span class=\"c1\"># Ensure URL has
    a trailing slash for consistency</span>\n            <span class=\"k\">return</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">base_url</span><span class=\"o\">.</span><span class=\"n\">rstrip</span><span
    class=\"p\">(</span><span class=\"s1\">&#39;/&#39;</span><span class=\"p\">)</span><span
    class=\"si\">}</span><span class=\"s2\">/</span><span class=\"si\">{</span><span
    class=\"n\">slug</span><span class=\"si\">}</span><span class=\"s2\">/&quot;</span>\n
    \       <span class=\"k\">return</span> <span class=\"n\">v</span>\n\n    <span
    class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">dict</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"n\">args</span><span class=\"p\">,</span> <span
    class=\"o\">**</span><span class=\"n\">kwargs</span><span class=\"p\">):</span>\n
    \       <span class=\"k\">return</span> <span class=\"p\">{</span><span class=\"s2\">&quot;url&quot;</span><span
    class=\"p\">:</span> <span class=\"p\">{</span><span class=\"o\">**</span><span
    class=\"nb\">super</span><span class=\"p\">()</span><span class=\"o\">.</span><span
    class=\"n\">dict</span><span class=\"p\">(</span><span class=\"o\">*</span><span
    class=\"n\">args</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"n\">kwargs</span><span class=\"p\">)}}</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition class\">\n<p class=\"admonition-title\">Class</p>\n<h2 id=\"SiteMapPost\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">SiteMapPost
    <em class=\"small\">class</em></h2>\n<p>[DEPRECATED] A model for posts that will
    be included in the sitemap.</p>\n<p>WARNING: This class is part of the deprecated
    sitemap plugin. Please migrate to\nthe feeds plugin which provides more comprehensive
    sitemap generation capabilities.</p>\n<p>To configure the base URL for your site,
    set the 'url' field in your markata config:</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"nt\">url</span><span
    class=\"p\">:</span><span class=\"w\"> </span><span class=\"l l-Scalar l-Scalar-Plain\">https://example.com</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">SiteMapPost
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
    class=\"w\"> </span><span class=\"nc\">SiteMapPost</span><span class=\"p\">(</span><span
    class=\"n\">pydantic</span><span class=\"o\">.</span><span class=\"n\">BaseModel</span><span
    class=\"p\">):</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;[DEPRECATED]
    A model for posts that will be included in the sitemap.</span>\n\n<span class=\"sd\">
    \   WARNING: This class is part of the deprecated sitemap plugin. Please migrate
    to</span>\n<span class=\"sd\">    the feeds plugin which provides more comprehensive
    sitemap generation capabilities.</span>\n\n<span class=\"sd\">    To configure
    the base URL for your site, set the &#39;url&#39; field in your markata config:</span>\n<span
    class=\"sd\">    ```yaml</span>\n<span class=\"sd\">    url: https://example.com</span>\n<span
    class=\"sd\">    ```</span>\n<span class=\"sd\">    &quot;&quot;&quot;</span>\n\n
    \   <span class=\"n\">slug</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"kc\">None</span>\n    <span class=\"n\">published</span><span
    class=\"p\">:</span> <span class=\"nb\">bool</span> <span class=\"o\">=</span>
    <span class=\"kc\">True</span>\n    <span class=\"n\">sitemap_url</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"n\">SiteMapUrl</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"kc\">None</span>\n    <span class=\"n\">markata</span><span class=\"p\">:</span>
    <span class=\"n\">Any</span> <span class=\"o\">=</span> <span class=\"n\">Field</span><span
    class=\"p\">(</span><span class=\"kc\">None</span><span class=\"p\">,</span> <span
    class=\"n\">exclude</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n\n    <span class=\"n\">model_config</span> <span class=\"o\">=</span>
    <span class=\"n\">pydantic</span><span class=\"o\">.</span><span class=\"n\">ConfigDict</span><span
    class=\"p\">(</span>\n        <span class=\"n\">validate_assignment</span><span
    class=\"o\">=</span><span class=\"kc\">False</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">arbitrary_types_allowed</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n        <span class=\"n\">extra</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;allow&quot;</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">str_strip_whitespace</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n        <span class=\"n\">validate_default</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">coerce_numbers_to_str</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n        <span class=\"n\">populate_by_name</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \   <span class=\"p\">)</span>\n\n    <span class=\"nd\">@pydantic</span><span
    class=\"o\">.</span><span class=\"n\">field_validator</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;sitemap_url&quot;</span><span class=\"p\">,</span> <span class=\"n\">mode</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;before&quot;</span><span class=\"p\">)</span>\n
    \   <span class=\"nd\">@classmethod</span>\n    <span class=\"k\">def</span><span
    class=\"w\"> </span><span class=\"nf\">validate_sitemap_url</span><span class=\"p\">(</span><span
    class=\"bp\">cls</span><span class=\"p\">,</span> <span class=\"n\">v</span><span
    class=\"p\">,</span> <span class=\"n\">info</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"n\">SiteMapUrl</span><span class=\"p\">]:</span>\n<span class=\"w\">        </span><span
    class=\"sd\">&quot;&quot;&quot;Initialize sitemap_url if not provided.&quot;&quot;&quot;</span>\n
    \       <span class=\"n\">markata</span> <span class=\"o\">=</span> <span class=\"n\">info</span><span
    class=\"o\">.</span><span class=\"n\">data</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;markata&quot;</span><span
    class=\"p\">)</span>\n        <span class=\"n\">slug</span> <span class=\"o\">=</span>
    <span class=\"n\">info</span><span class=\"o\">.</span><span class=\"n\">data</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;slug&quot;</span><span class=\"p\">)</span>\n        <span
    class=\"k\">if</span> <span class=\"n\">markata</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span> <span class=\"ow\">or</span> <span class=\"n\">slug</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">raise</span> <span class=\"ne\">ValueError</span><span
    class=\"p\">(</span>\n                <span class=\"s2\">&quot;Could not create
    sitemap entry: markata and slug are required. &quot;</span>\n                <span
    class=\"s2\">&quot;This usually means the Post model is missing required fields.
    &quot;</span>\n                <span class=\"s2\">&quot;Check that your post has
    a valid slug and markata instance.&quot;</span>\n            <span class=\"p\">)</span>\n\n
    \       <span class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"k\">return</span>
    <span class=\"n\">SiteMapUrl</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">slug</span><span class=\"o\">=</span><span class=\"n\">slug</span><span
    class=\"p\">)</span>\n        <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"nb\">dict</span><span class=\"p\">):</span>\n            <span class=\"k\">return</span>
    <span class=\"n\">SiteMapUrl</span><span class=\"p\">(</span><span class=\"o\">**</span><span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">slug</span><span class=\"o\">=</span><span class=\"n\">slug</span><span
    class=\"p\">)</span>\n        <span class=\"k\">if</span> <span class=\"n\">v</span><span
    class=\"o\">.</span><span class=\"n\">markata</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"n\">v</span><span
    class=\"o\">.</span><span class=\"n\">markata</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span>\n        <span class=\"k\">if</span> <span class=\"n\">v</span><span
    class=\"o\">.</span><span class=\"n\">slug</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"n\">v</span><span
    class=\"o\">.</span><span class=\"n\">slug</span> <span class=\"o\">=</span> <span
    class=\"n\">slug</span>\n        <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"validate_loc\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">validate_loc
    <em class=\"small\">method</em></h2>\n<p>Generate the URL for the sitemap entry.</p>\n<p>Uses
    markata.config.url as the base URL if set, otherwise uses relative URLs.\nExample:
    <a href=\"https://example.com/my-page/\">https://example.com/my-page/</a> or /my-page/</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">validate_loc
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
    class=\"w\"> </span><span class=\"nf\">validate_loc</span><span class=\"p\">(</span><span
    class=\"bp\">cls</span><span class=\"p\">,</span> <span class=\"n\">v</span><span
    class=\"p\">,</span> <span class=\"n\">info</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n<span
    class=\"w\">        </span><span class=\"sd\">&quot;&quot;&quot;Generate the URL
    for the sitemap entry.</span>\n\n<span class=\"sd\">        Uses markata.config.url
    as the base URL if set, otherwise uses relative URLs.</span>\n<span class=\"sd\">
    \       Example: https://example.com/my-page/ or /my-page/</span>\n<span class=\"sd\">
    \       &quot;&quot;&quot;</span>\n        <span class=\"k\">if</span> <span class=\"n\">v</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">markata</span> <span class=\"o\">=</span> <span
    class=\"n\">info</span><span class=\"o\">.</span><span class=\"n\">data</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;markata&quot;</span><span class=\"p\">)</span>\n            <span
    class=\"n\">slug</span> <span class=\"o\">=</span> <span class=\"n\">info</span><span
    class=\"o\">.</span><span class=\"n\">data</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">)</span>\n            <span class=\"k\">if</span> <span class=\"n\">markata</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span> <span class=\"ow\">or</span>
    <span class=\"n\">slug</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"k\">raise</span> <span class=\"ne\">ValueError</span><span
    class=\"p\">(</span>\n                    <span class=\"s2\">&quot;Could not generate
    sitemap URL: markata and slug are required. &quot;</span>\n                    <span
    class=\"s2\">&quot;This usually means the Post model is missing required fields.
    &quot;</span>\n                    <span class=\"s2\">&quot;Check that your post
    has a valid slug and markata instance.&quot;</span>\n                <span class=\"p\">)</span>\n\n
    \           <span class=\"c1\"># Get base URL from config, default to empty string
    if not set</span>\n            <span class=\"n\">base_url</span> <span class=\"o\">=</span>
    <span class=\"nb\">getattr</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;url&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;&quot;</span><span
    class=\"p\">)</span>\n            <span class=\"k\">if</span> <span class=\"ow\">not</span>
    <span class=\"n\">base_url</span><span class=\"p\">:</span>\n                <span
    class=\"k\">return</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;/</span><span
    class=\"si\">{</span><span class=\"n\">slug</span><span class=\"si\">}</span><span
    class=\"s2\">/&quot;</span>\n\n            <span class=\"c1\"># Ensure URL has
    a trailing slash for consistency</span>\n            <span class=\"k\">return</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">base_url</span><span class=\"o\">.</span><span class=\"n\">rstrip</span><span
    class=\"p\">(</span><span class=\"s1\">&#39;/&#39;</span><span class=\"p\">)</span><span
    class=\"si\">}</span><span class=\"s2\">/</span><span class=\"si\">{</span><span
    class=\"n\">slug</span><span class=\"si\">}</span><span class=\"s2\">/&quot;</span>\n
    \       <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"validate_sitemap_url\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">validate_sitemap_url
    <em class=\"small\">method</em></h2>\n<p>Initialize sitemap_url if not provided.</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">validate_sitemap_url
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
    class=\"w\"> </span><span class=\"nf\">validate_sitemap_url</span><span class=\"p\">(</span><span
    class=\"bp\">cls</span><span class=\"p\">,</span> <span class=\"n\">v</span><span
    class=\"p\">,</span> <span class=\"n\">info</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"n\">SiteMapUrl</span><span class=\"p\">]:</span>\n<span class=\"w\">        </span><span
    class=\"sd\">&quot;&quot;&quot;Initialize sitemap_url if not provided.&quot;&quot;&quot;</span>\n
    \       <span class=\"n\">markata</span> <span class=\"o\">=</span> <span class=\"n\">info</span><span
    class=\"o\">.</span><span class=\"n\">data</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;markata&quot;</span><span
    class=\"p\">)</span>\n        <span class=\"n\">slug</span> <span class=\"o\">=</span>
    <span class=\"n\">info</span><span class=\"o\">.</span><span class=\"n\">data</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;slug&quot;</span><span class=\"p\">)</span>\n        <span
    class=\"k\">if</span> <span class=\"n\">markata</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span> <span class=\"ow\">or</span> <span class=\"n\">slug</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">raise</span> <span class=\"ne\">ValueError</span><span
    class=\"p\">(</span>\n                <span class=\"s2\">&quot;Could not create
    sitemap entry: markata and slug are required. &quot;</span>\n                <span
    class=\"s2\">&quot;This usually means the Post model is missing required fields.
    &quot;</span>\n                <span class=\"s2\">&quot;Check that your post has
    a valid slug and markata instance.&quot;</span>\n            <span class=\"p\">)</span>\n\n
    \       <span class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"k\">return</span>
    <span class=\"n\">SiteMapUrl</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">slug</span><span class=\"o\">=</span><span class=\"n\">slug</span><span
    class=\"p\">)</span>\n        <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"nb\">dict</span><span class=\"p\">):</span>\n            <span class=\"k\">return</span>
    <span class=\"n\">SiteMapUrl</span><span class=\"p\">(</span><span class=\"o\">**</span><span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">slug</span><span class=\"o\">=</span><span class=\"n\">slug</span><span
    class=\"p\">)</span>\n        <span class=\"k\">if</span> <span class=\"n\">v</span><span
    class=\"o\">.</span><span class=\"n\">markata</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"n\">v</span><span
    class=\"o\">.</span><span class=\"n\">markata</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span>\n        <span class=\"k\">if</span> <span class=\"n\">v</span><span
    class=\"o\">.</span><span class=\"n\">slug</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"n\">v</span><span
    class=\"o\">.</span><span class=\"n\">slug</span> <span class=\"o\">=</span> <span
    class=\"n\">slug</span>\n        <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n</div>\n\n
    \   </section>\n</article>"
  raw.md: ''
published: false
slug: markata/plugins/sitemap
title: sitemap.py


---

---

[DEPRECATED] The `markata.plugins.sitemap` plugin is deprecated and will be removed in a
future version. Please use `markata.plugins.feeds` instead, which provides more
comprehensive sitemap generation capabilities.

## Installation

This plugin is deprecated. Use `markata.plugins.feeds` instead:

```toml
hooks = [
    "markata.plugins.feeds",  # Use this instead
    # "markata.plugins.sitemap",  # Deprecated
]
```

# Migration Guide

To migrate to the new feeds plugin:

1. Remove sitemap plugin from hooks:
```toml
# Remove or comment out
# "markata.plugins.sitemap"
```

2. Add feeds plugin:
```toml
hooks = [
    "markata.plugins.feeds"
]
```

3. Update configuration:
```toml
[markata.feeds]
# Sitemap configuration
sitemap = { output = "sitemap.xml" }

# Optional: Configure sitemap settings
[markata.feeds.sitemap.options]
changefreq = "daily"
priority = "0.7"
```

See the feeds plugin documentation for more configuration options.

# Legacy Configuration

If you must continue using this plugin temporarily, configure in `markata.toml`:

```toml
[markata]
url = "https://example.com"

[markata.sitemap]
changefreq = "daily"
priority = "0.7"
```

# Dependencies

This plugin depends on:
- pydantic for configuration

WARNING: This plugin is deprecated and will be removed in a future version.
Please migrate to `markata.plugins.feeds` as soon as possible.

---

!!! class
    <h2 id="SiteMapUrl" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">SiteMapUrl <em class="small">class</em></h2>

    [DEPRECATED] A model representing a URL entry in the sitemap.xml file.

    WARNING: This class is part of the deprecated sitemap plugin. Please migrate to
    the feeds plugin which provides more comprehensive sitemap generation capabilities.

    To configure the base URL for your site, set the 'url' field in your markata config:
    ```yaml
    url: https://example.com
    ```

    If no base URL is set, relative URLs will be used.

???+ source "SiteMapUrl <em class='small'>source</em>"
    ```python
    class SiteMapUrl(pydantic.BaseModel):
        """[DEPRECATED] A model representing a URL entry in the sitemap.xml file.

        WARNING: This class is part of the deprecated sitemap plugin. Please migrate to
        the feeds plugin which provides more comprehensive sitemap generation capabilities.

        To configure the base URL for your site, set the 'url' field in your markata config:
        ```yaml
        url: https://example.com
        ```

        If no base URL is set, relative URLs will be used.
        """

        slug: str = Field(..., exclude=True)
        loc: str = Field(
            None,
            include=True,
            description="The full URL for this page in the sitemap. Generated automatically from config.url + slug.",
        )
        changefreq: str = Field("daily", include=True)
        priority: str = Field("0.7", include=True)
        markata: Any = Field(None, exclude=True)

        model_config = pydantic.ConfigDict(
            validate_assignment=False,
            arbitrary_types_allowed=True,
            extra="allow",
            str_strip_whitespace=True,
            validate_default=True,
            coerce_numbers_to_str=True,
            populate_by_name=True,
        )

        @pydantic.field_validator("loc", mode="before")
        @classmethod
        def validate_loc(cls, v, info) -> str:
            """Generate the URL for the sitemap entry.

            Uses markata.config.url as the base URL if set, otherwise uses relative URLs.
            Example: https://example.com/my-page/ or /my-page/
            """
            if v is None:
                markata = info.data.get("markata")
                slug = info.data.get("slug")
                if markata is None or slug is None:
                    raise ValueError(
                        "Could not generate sitemap URL: markata and slug are required. "
                        "This usually means the Post model is missing required fields. "
                        "Check that your post has a valid slug and markata instance."
                    )

                # Get base URL from config, default to empty string if not set
                base_url = getattr(markata.config, "url", "")
                if not base_url:
                    return f"/{slug}/"

                # Ensure URL has a trailing slash for consistency
                return f"{base_url.rstrip('/')}/{slug}/"
            return v

        def dict(self, *args, **kwargs):
            return {"url": {**super().dict(*args, **kwargs)}}
    ```
!!! class
    <h2 id="SiteMapPost" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">SiteMapPost <em class="small">class</em></h2>

    [DEPRECATED] A model for posts that will be included in the sitemap.

    WARNING: This class is part of the deprecated sitemap plugin. Please migrate to
    the feeds plugin which provides more comprehensive sitemap generation capabilities.

    To configure the base URL for your site, set the 'url' field in your markata config:
    ```yaml
    url: https://example.com
    ```

???+ source "SiteMapPost <em class='small'>source</em>"
    ```python
    class SiteMapPost(pydantic.BaseModel):
        """[DEPRECATED] A model for posts that will be included in the sitemap.

        WARNING: This class is part of the deprecated sitemap plugin. Please migrate to
        the feeds plugin which provides more comprehensive sitemap generation capabilities.

        To configure the base URL for your site, set the 'url' field in your markata config:
        ```yaml
        url: https://example.com
        ```
        """

        slug: str = None
        published: bool = True
        sitemap_url: Optional[SiteMapUrl] = None
        markata: Any = Field(None, exclude=True)

        model_config = pydantic.ConfigDict(
            validate_assignment=False,
            arbitrary_types_allowed=True,
            extra="allow",
            str_strip_whitespace=True,
            validate_default=True,
            coerce_numbers_to_str=True,
            populate_by_name=True,
        )

        @pydantic.field_validator("sitemap_url", mode="before")
        @classmethod
        def validate_sitemap_url(cls, v, info) -> Optional[SiteMapUrl]:
            """Initialize sitemap_url if not provided."""
            markata = info.data.get("markata")
            slug = info.data.get("slug")
            if markata is None or slug is None:
                raise ValueError(
                    "Could not create sitemap entry: markata and slug are required. "
                    "This usually means the Post model is missing required fields. "
                    "Check that your post has a valid slug and markata instance."
                )

            if v is None:
                return SiteMapUrl(markata=markata, slug=slug)
            if isinstance(v, dict):
                return SiteMapUrl(**v, markata=markata, slug=slug)
            if v.markata is None:
                v.markata = markata
            if v.slug is None:
                v.slug = slug
            return v
    ```
!!! method
    <h2 id="validate_loc" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">validate_loc <em class="small">method</em></h2>

    Generate the URL for the sitemap entry.

    Uses markata.config.url as the base URL if set, otherwise uses relative URLs.
    Example: https://example.com/my-page/ or /my-page/

???+ source "validate_loc <em class='small'>source</em>"
    ```python
    def validate_loc(cls, v, info) -> str:
            """Generate the URL for the sitemap entry.

            Uses markata.config.url as the base URL if set, otherwise uses relative URLs.
            Example: https://example.com/my-page/ or /my-page/
            """
            if v is None:
                markata = info.data.get("markata")
                slug = info.data.get("slug")
                if markata is None or slug is None:
                    raise ValueError(
                        "Could not generate sitemap URL: markata and slug are required. "
                        "This usually means the Post model is missing required fields. "
                        "Check that your post has a valid slug and markata instance."
                    )

                # Get base URL from config, default to empty string if not set
                base_url = getattr(markata.config, "url", "")
                if not base_url:
                    return f"/{slug}/"

                # Ensure URL has a trailing slash for consistency
                return f"{base_url.rstrip('/')}/{slug}/"
            return v
    ```
!!! method
    <h2 id="validate_sitemap_url" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">validate_sitemap_url <em class="small">method</em></h2>

    Initialize sitemap_url if not provided.

???+ source "validate_sitemap_url <em class='small'>source</em>"
    ```python
    def validate_sitemap_url(cls, v, info) -> Optional[SiteMapUrl]:
            """Initialize sitemap_url if not provided."""
            markata = info.data.get("markata")
            slug = info.data.get("slug")
            if markata is None or slug is None:
                raise ValueError(
                    "Could not create sitemap entry: markata and slug are required. "
                    "This usually means the Post model is missing required fields. "
                    "Check that your post has a valid slug and markata instance."
                )

            if v is None:
                return SiteMapUrl(markata=markata, slug=slug)
            if isinstance(v, dict):
                return SiteMapUrl(**v, markata=markata, slug=slug)
            if v.markata is None:
                v.markata = markata
            if v.slug is None:
                v.slug = slug
            return v
    ```