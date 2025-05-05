---
content: "---\n\nJinja2 Environment Plugin\n\nProvides a centralized Jinja2 environment
  configuration for consistent template rendering\nacross all Markata plugins. This
  plugin ensures template rendering behavior is consistent\nand available even when
  specific template-using plugins are not enabled.\n\n## Installation\n\nThis plugin
  is built-in and enabled by default through the 'default' plugin.\nIf you want to
  be explicit, you can add it to your list of plugins:\n\n``` toml\nhooks = [\n    \"markata.plugins.jinja_env\",\n]\n```\n\n##
  Uninstallation\n\nSince this plugin is included in the default plugin set, to disable
  it you must explicitly\nadd it to the disabled_hooks list if you are using the 'default'
  plugin:\n\n``` toml\ndisabled_hooks = [\n    \"markata.plugins.jinja_env\",\n]\n```\n\n##
  Configuration\n\nConfigure Jinja environment settings in your `markata.toml`:\n\n```
  toml\n[markata.jinja_env]\ntemplate_paths = [\"templates\"]  # Additional template
  paths to search\nundefined_silent = true        # Return empty string for undefined
  variables\ntrim_blocks = true            # Remove first newline after block\nlstrip_blocks
  = true          # Strip tabs/spaces from start of line\ntemplate_cache_dir = \".markata.cache/template_bytecode\"\n```\n\n#
  Usage\n\nThe environment is automatically available to other plugins via `markata.config.jinja_env`.\nTemplate
  loading follows this order:\n1. Package templates (built-in Markata templates)\n2.
  User template paths (configured via template_paths)\n\nExample usage in a plugin:\n\n```
  python\ndef render_template(markata, content):\n    template = markata.jinja_env.from_string(content)\n
  \   return template.render(markata=markata)\n```\n\n# Notes\n\n- Template paths
  are resolved relative to the current working directory\n- Package templates are
  always available and take precedence\n- Silent undefined behavior means undefined
  variables render as empty strings\n\n---\n\n!!! class\n    <h2 id=\"_SilentUndefined\"
  class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">_SilentUndefined
  <em class=\"small\">class</em></h2>\n\n    Custom undefined type that returns empty
  string for undefined variables.\n\n???+ source \"_SilentUndefined <em class='small'>source</em>\"\n
  \   ```python\n    class _SilentUndefined(jinja2.Undefined):\n        \"\"\"Custom
  undefined type that returns empty string for undefined variables.\"\"\"\n\n        def
  _fail_with_undefined_error(self, *args, **kwargs):\n            return \"\"\n\n
  \       __add__ = __radd__ = __mul__ = __rmul__ = __div__ = __rdiv__ = __truediv__
  = (\n            __rtruediv__\n        ) = __floordiv__ = __rfloordiv__ = __mod__
  = __rmod__ = __pos__ = __neg__ = (\n            __call__\n        ) = __getitem__
  = __lt__ = __le__ = __gt__ = __ge__ = __int__ = __float__ = (\n            __complex__\n
  \       ) = __pow__ = __rpow__ = _fail_with_undefined_error\n    ```\n!!! class\n
  \   <h2 id=\"MarkataTemplateCache\" class=\"admonition-title\" style=\"margin: 0;
  padding: .5rem 1rem;\">MarkataTemplateCache <em class=\"small\">class</em></h2>\n\n
  \   Template bytecode cache for improved performance.\n\n???+ source \"MarkataTemplateCache
  <em class='small'>source</em>\"\n    ```python\n    class MarkataTemplateCache(jinja2.BytecodeCache):\n
  \       \"\"\"Template bytecode cache for improved performance.\"\"\"\n\n        def
  __init__(self, directory):\n            self.directory = Path(directory)\n            self.directory.mkdir(parents=True,
  exist_ok=True)\n\n        def load_bytecode(self, bucket):\n            filename
  = self.directory / f\"{bucket.key}.cache\"\n            if filename.exists():\n
  \               with open(filename, \"rb\") as f:\n                    bucket.bytecode_from_string(f.read())\n\n
  \       def dump_bytecode(self, bucket):\n            filename = self.directory
  / f\"{bucket.key}.cache\"\n            with open(filename, \"wb\") as f:\n                f.write(bucket.bytecode_to_string())\n
  \   ```\n!!! class\n    <h2 id=\"JinjaEnvConfig\" class=\"admonition-title\" style=\"margin:
  0; padding: .5rem 1rem;\">JinjaEnvConfig <em class=\"small\">class</em></h2>\n\n
  \   Configuration for the Jinja environment.\n\n???+ source \"JinjaEnvConfig <em
  class='small'>source</em>\"\n    ```python\n    class JinjaEnvConfig(pydantic.BaseModel):\n
  \       \"\"\"Configuration for the Jinja environment.\"\"\"\n\n        templates_dir:
  List[str] = []\n        undefined_silent: bool = True\n        trim_blocks: bool
  = True\n        lstrip_blocks: bool = True\n        template_cache_dir: Path = Path(\".markata.cache/template_bytecode\")\n\n
  \       model_config = pydantic.ConfigDict(\n            validate_assignment=True,
  \ # Config model\n            arbitrary_types_allowed=True,\n            extra=\"allow\",\n
  \           str_strip_whitespace=True,\n            validate_default=True,\n            coerce_numbers_to_str=True,\n
  \           populate_by_name=True,\n        )\n    ```\n!!! function\n    <h2 id=\"config_model\"
  class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">config_model
  <em class=\"small\">function</em></h2>\n\n    Register configuration models.\n\n???+
  source \"config_model <em class='small'>source</em>\"\n    ```python\n    def config_model(markata:
  \"Markata\") -> None:\n        \"\"\"Register configuration models.\"\"\"\n        markata.config_models.append(JinjaEnvConfig)\n
  \   ```\n!!! function\n    <h2 id=\"configure\" class=\"admonition-title\" style=\"margin:
  0; padding: .5rem 1rem;\">configure <em class=\"small\">function</em></h2>\n\n    Initialize
  and configure the Jinja2 environment for Markata.\n\n    This hook runs early in
  the configuration stage to ensure the jinja environment\n    is available for other
  plugins that need it during configuration.\n\n    Args:\n        markata: The Markata
  instance\n\n???+ source \"configure <em class='small'>source</em>\"\n    ```python\n
  \   def configure(markata: Markata) -> None:\n        \"\"\"Initialize and configure
  the Jinja2 environment for Markata.\n\n        This hook runs early in the configuration
  stage to ensure the jinja environment\n        is available for other plugins that
  need it during configuration.\n\n        Args:\n            markata: The Markata
  instance\n        \"\"\"\n        # Get configuration, falling back to defaults\n
  \       config = JinjaEnvConfig()\n        if hasattr(markata.config, \"jinja_env\"):\n
  \           if isinstance(markata.config.jinja_env, dict):\n                config
  = JinjaEnvConfig(**markata.config.jinja_env)\n\n        # TODO: setting up env twice
  could not get dynamic templates to be recognized on first pass\n        loaders
  = []\n        if markata.config.templates_dir:\n            for path in markata.config.templates_dir:\n
  \               path = Path(path).expanduser().resolve()\n                if path.exists():\n
  \                   loaders.append(FileSystemLoader(str(path)))\n        # Create
  environment\n        env_for_dynamic_render = Environment(\n            loader=ChoiceLoader(loaders),\n
  \           undefined=_SilentUndefined if config.undefined_silent else jinja2.Undefined,\n
  \           trim_blocks=config.trim_blocks,\n            lstrip_blocks=config.lstrip_blocks,\n
  \           bytecode_cache=MarkataTemplateCache(config.template_cache_dir),\n            auto_reload=True,\n
  \       )\n\n        markata.config.dynamic_templates_dir.mkdir(parents=True, exist_ok=True)\n
  \       head_template = markata.config.dynamic_templates_dir / \"head.html\"\n        head_template.write_text(\n
  \           env_for_dynamic_render.get_template(\"dynamic_head.html\").render(\n
  \               {\"markata\": markata}\n            ),\n        )\n\n        # Set
  up loaders\n        loaders = []\n\n        # Add package templates first (lowest
  priority)\n        # loaders.append(PackageLoader(\"markata\", \"templates\"))\n\n
  \       # Add user template paths (medium priority)\n        if markata.config.templates_dir:\n
  \           for path in markata.config.templates_dir:\n                path = Path(path).expanduser().resolve()\n
  \               if path.exists():\n                    loaders.append(FileSystemLoader(str(path)))\n\n
  \       # Add dynamic templates directory (highest priority)\n        # dynamic_templates_dir
  = Path(\".markata.cache/templates\")\n        # dynamic_templates_dir.mkdir(parents=True,
  exist_ok=True)\n        # loaders.append(FileSystemLoader(str(dynamic_templates_dir)))\n\n
  \       # Create environment\n        env = Environment(\n            loader=ChoiceLoader(loaders),\n
  \           undefined=_SilentUndefined if config.undefined_silent else jinja2.Undefined,\n
  \           trim_blocks=config.trim_blocks,\n            lstrip_blocks=config.lstrip_blocks,\n
  \           bytecode_cache=MarkataTemplateCache(config.template_cache_dir),\n            auto_reload=True,\n
  \       )\n\n        # Register the environment on the config's private attribute\n
  \       markata.jinja_env = env\n    ```"
date: 2025-05-05
description: "Jinja2 Environment Plugin Provides a centralized Jinja2 environment
  configuration for consistent template rendering across all Markata plugins. This
  plugin\u2026"
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>jinja_env.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Jinja2 Environment Plugin Provides a
    centralized Jinja2 environment configuration for consistent template rendering
    across all Markata plugins. This plugin\u2026\" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n
    \       <meta property=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    <head>\n        <script src=\"https://unpkg.com/@tailwindcss/browser@4\"></script>\n<title>jinja_env.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Jinja2 Environment Plugin Provides a
    centralized Jinja2 environment configuration for consistent template rendering
    across all Markata plugins. This plugin\u2026\" />\n <link href=\"/favicon.ico\"
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        jinja_env.py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <hr />\n<p>Jinja2 Environment
    Plugin</p>\n<p>Provides a centralized Jinja2 environment configuration for consistent
    template rendering\nacross all Markata plugins. This plugin ensures template rendering
    behavior is consistent\nand available even when specific template-using plugins
    are not enabled.</p>\n<h2 id=\"installation\">Installation <a class=\"header-anchor\"
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
    class=\"p\">[</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;markata.plugins.jinja_env&quot;</span><span
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
    class=\"p\">[</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;markata.plugins.jinja_env&quot;</span><span
    class=\"p\">,</span>\n<span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<h2
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Configure Jinja environment
    settings in your <code>markata.toml</code>:</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata.jinja_env]</span>\n<span
    class=\"n\">template_paths</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"p\">[</span><span class=\"s2\">&quot;templates&quot;</span><span
    class=\"p\">]</span><span class=\"w\">  </span><span class=\"c1\"># Additional
    template paths to search</span>\n<span class=\"n\">undefined_silent</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"kc\">true</span><span class=\"w\">        </span><span class=\"c1\">#
    Return empty string for undefined variables</span>\n<span class=\"n\">trim_blocks</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"kc\">true</span><span class=\"w\">            </span><span class=\"c1\">#
    Remove first newline after block</span>\n<span class=\"n\">lstrip_blocks</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"kc\">true</span><span class=\"w\">          </span><span class=\"c1\">#
    Strip tabs/spaces from start of line</span>\n<span class=\"n\">template_cache_dir</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;.markata.cache/template_bytecode&quot;</span>\n</pre></div>\n\n</pre>\n\n<h1
    id=\"usage\">Usage <a class=\"header-anchor\" href=\"#usage\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>The environment is automatically
    available to other plugins via <code>markata.config.jinja_env</code>.\nTemplate
    loading follows this order:</p>\n<ol>\n<li>Package templates (built-in Markata
    templates)</li>\n<li>User template paths (configured via template_paths)</li>\n</ol>\n<p>Example
    usage in a plugin:</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    class=\"w\"> </span><span class=\"nf\">render_template</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">,</span> <span class=\"n\">content</span><span
    class=\"p\">):</span>\n    <span class=\"n\">template</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">jinja_env</span><span
    class=\"o\">.</span><span class=\"n\">from_string</span><span class=\"p\">(</span><span
    class=\"n\">content</span><span class=\"p\">)</span>\n    <span class=\"k\">return</span>
    <span class=\"n\">template</span><span class=\"o\">.</span><span class=\"n\">render</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<h1
    id=\"notes\">Notes <a class=\"header-anchor\" href=\"#notes\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<ul>\n<li>Template paths
    are resolved relative to the current working directory</li>\n<li>Package templates
    are always available and take precedence</li>\n<li>Silent undefined behavior means
    undefined variables render as empty strings</li>\n</ul>\n<hr />\n<div class=\"admonition
    class\">\n<p class=\"admonition-title\">Class</p>\n<h2 id=\"_SilentUndefined\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">_SilentUndefined
    <em class=\"small\">class</em></h2>\n<p>Custom undefined type that returns empty
    string for undefined variables.</p>\n</div>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">_SilentUndefined <em class='small'>source</em></p>\n<pre
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
    class=\"w\"> </span><span class=\"nc\">_SilentUndefined</span><span class=\"p\">(</span><span
    class=\"n\">jinja2</span><span class=\"o\">.</span><span class=\"n\">Undefined</span><span
    class=\"p\">):</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Custom
    undefined type that returns empty string for undefined variables.&quot;&quot;&quot;</span>\n\n
    \   <span class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">_fail_with_undefined_error</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"n\">args</span><span class=\"p\">,</span> <span
    class=\"o\">**</span><span class=\"n\">kwargs</span><span class=\"p\">):</span>\n
    \       <span class=\"k\">return</span> <span class=\"s2\">&quot;&quot;</span>\n\n
    \   <span class=\"fm\">__add__</span> <span class=\"o\">=</span> <span class=\"fm\">__radd__</span>
    <span class=\"o\">=</span> <span class=\"fm\">__mul__</span> <span class=\"o\">=</span>
    <span class=\"fm\">__rmul__</span> <span class=\"o\">=</span> <span class=\"n\">__div__</span>
    <span class=\"o\">=</span> <span class=\"n\">__rdiv__</span> <span class=\"o\">=</span>
    <span class=\"fm\">__truediv__</span> <span class=\"o\">=</span> <span class=\"p\">(</span>\n
    \       <span class=\"fm\">__rtruediv__</span>\n    <span class=\"p\">)</span>
    <span class=\"o\">=</span> <span class=\"fm\">__floordiv__</span> <span class=\"o\">=</span>
    <span class=\"fm\">__rfloordiv__</span> <span class=\"o\">=</span> <span class=\"fm\">__mod__</span>
    <span class=\"o\">=</span> <span class=\"fm\">__rmod__</span> <span class=\"o\">=</span>
    <span class=\"fm\">__pos__</span> <span class=\"o\">=</span> <span class=\"fm\">__neg__</span>
    <span class=\"o\">=</span> <span class=\"p\">(</span>\n        <span class=\"fm\">__call__</span>\n
    \   <span class=\"p\">)</span> <span class=\"o\">=</span> <span class=\"fm\">__getitem__</span>
    <span class=\"o\">=</span> <span class=\"fm\">__lt__</span> <span class=\"o\">=</span>
    <span class=\"fm\">__le__</span> <span class=\"o\">=</span> <span class=\"fm\">__gt__</span>
    <span class=\"o\">=</span> <span class=\"fm\">__ge__</span> <span class=\"o\">=</span>
    <span class=\"fm\">__int__</span> <span class=\"o\">=</span> <span class=\"fm\">__float__</span>
    <span class=\"o\">=</span> <span class=\"p\">(</span>\n        <span class=\"fm\">__complex__</span>\n
    \   <span class=\"p\">)</span> <span class=\"o\">=</span> <span class=\"fm\">__pow__</span>
    <span class=\"o\">=</span> <span class=\"fm\">__rpow__</span> <span class=\"o\">=</span>
    <span class=\"n\">_fail_with_undefined_error</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition class\">\n<p class=\"admonition-title\">Class</p>\n<h2 id=\"MarkataTemplateCache\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">MarkataTemplateCache
    <em class=\"small\">class</em></h2>\n<p>Template bytecode cache for improved performance.</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">MarkataTemplateCache
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
    class=\"w\"> </span><span class=\"nc\">MarkataTemplateCache</span><span class=\"p\">(</span><span
    class=\"n\">jinja2</span><span class=\"o\">.</span><span class=\"n\">BytecodeCache</span><span
    class=\"p\">):</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Template
    bytecode cache for improved performance.&quot;&quot;&quot;</span>\n\n    <span
    class=\"k\">def</span><span class=\"w\"> </span><span class=\"fm\">__init__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"n\">directory</span><span class=\"p\">):</span>\n        <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">directory</span> <span class=\"o\">=</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">directory</span><span
    class=\"p\">)</span>\n        <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">directory</span><span class=\"o\">.</span><span class=\"n\">mkdir</span><span
    class=\"p\">(</span><span class=\"n\">parents</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span> <span class=\"n\">exist_ok</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n\n
    \   <span class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">load_bytecode</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"n\">bucket</span><span class=\"p\">):</span>\n        <span class=\"n\">filename</span>
    <span class=\"o\">=</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">directory</span> <span class=\"o\">/</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;</span><span class=\"si\">{</span><span class=\"n\">bucket</span><span
    class=\"o\">.</span><span class=\"n\">key</span><span class=\"si\">}</span><span
    class=\"s2\">.cache&quot;</span>\n        <span class=\"k\">if</span> <span class=\"n\">filename</span><span
    class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">():</span>\n
    \           <span class=\"k\">with</span> <span class=\"nb\">open</span><span
    class=\"p\">(</span><span class=\"n\">filename</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;rb&quot;</span><span class=\"p\">)</span> <span class=\"k\">as</span>
    <span class=\"n\">f</span><span class=\"p\">:</span>\n                <span class=\"n\">bucket</span><span
    class=\"o\">.</span><span class=\"n\">bytecode_from_string</span><span class=\"p\">(</span><span
    class=\"n\">f</span><span class=\"o\">.</span><span class=\"n\">read</span><span
    class=\"p\">())</span>\n\n    <span class=\"k\">def</span><span class=\"w\"> </span><span
    class=\"nf\">dump_bytecode</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">bucket</span><span class=\"p\">):</span>\n
    \       <span class=\"n\">filename</span> <span class=\"o\">=</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">directory</span> <span class=\"o\">/</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">bucket</span><span class=\"o\">.</span><span class=\"n\">key</span><span
    class=\"si\">}</span><span class=\"s2\">.cache&quot;</span>\n        <span class=\"k\">with</span>
    <span class=\"nb\">open</span><span class=\"p\">(</span><span class=\"n\">filename</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;wb&quot;</span><span class=\"p\">)</span>
    <span class=\"k\">as</span> <span class=\"n\">f</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">f</span><span class=\"o\">.</span><span class=\"n\">write</span><span
    class=\"p\">(</span><span class=\"n\">bucket</span><span class=\"o\">.</span><span
    class=\"n\">bytecode_to_string</span><span class=\"p\">())</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition class\">\n<p class=\"admonition-title\">Class</p>\n<h2 id=\"JinjaEnvConfig\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">JinjaEnvConfig
    <em class=\"small\">class</em></h2>\n<p>Configuration for the Jinja environment.</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">JinjaEnvConfig
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
    class=\"w\"> </span><span class=\"nc\">JinjaEnvConfig</span><span class=\"p\">(</span><span
    class=\"n\">pydantic</span><span class=\"o\">.</span><span class=\"n\">BaseModel</span><span
    class=\"p\">):</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Configuration
    for the Jinja environment.&quot;&quot;&quot;</span>\n\n    <span class=\"n\">templates_dir</span><span
    class=\"p\">:</span> <span class=\"n\">List</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"p\">[]</span>\n    <span class=\"n\">undefined_silent</span><span class=\"p\">:</span>
    <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"kc\">True</span>\n
    \   <span class=\"n\">trim_blocks</span><span class=\"p\">:</span> <span class=\"nb\">bool</span>
    <span class=\"o\">=</span> <span class=\"kc\">True</span>\n    <span class=\"n\">lstrip_blocks</span><span
    class=\"p\">:</span> <span class=\"nb\">bool</span> <span class=\"o\">=</span>
    <span class=\"kc\">True</span>\n    <span class=\"n\">template_cache_dir</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span> <span class=\"o\">=</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"s2\">&quot;.markata.cache/template_bytecode&quot;</span><span
    class=\"p\">)</span>\n\n    <span class=\"n\">model_config</span> <span class=\"o\">=</span>
    <span class=\"n\">pydantic</span><span class=\"o\">.</span><span class=\"n\">ConfigDict</span><span
    class=\"p\">(</span>\n        <span class=\"n\">validate_assignment</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>  <span
    class=\"c1\"># Config model</span>\n        <span class=\"n\">arbitrary_types_allowed</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">extra</span><span class=\"o\">=</span><span class=\"s2\">&quot;allow&quot;</span><span
    class=\"p\">,</span>\n        <span class=\"n\">str_strip_whitespace</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">validate_default</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n        <span class=\"n\">coerce_numbers_to_str</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">populate_by_name</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n    <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition function\">\n<p class=\"admonition-title\">Function</p>\n<h2
    id=\"config_model\" class=\"admonition-title\" style=\"margin: 0; padding: .5rem
    1rem;\">config_model <em class=\"small\">function</em></h2>\n<p>Register configuration
    models.</p>\n</div>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">config_model <em class='small'>source</em></p>\n<pre
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
    class=\"w\"> </span><span class=\"nf\">config_model</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Register
    configuration models.&quot;&quot;&quot;</span>\n    <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config_models</span><span class=\"o\">.</span><span
    class=\"n\">append</span><span class=\"p\">(</span><span class=\"n\">JinjaEnvConfig</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div class=\"admonition
    function\">\n<p class=\"admonition-title\">Function</p>\n<h2 id=\"configure\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">configure
    <em class=\"small\">function</em></h2>\n<p>Initialize and configure the Jinja2
    environment for Markata.</p>\n<p>This hook runs early in the configuration stage
    to ensure the jinja environment\nis available for other plugins that need it during
    configuration.</p>\n<p>Args:\nmarkata: The Markata instance</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">configure
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
    class=\"w\"> </span><span class=\"nf\">configure</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"n\">Markata</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Initialize
    and configure the Jinja2 environment for Markata.</span>\n\n<span class=\"sd\">
    \   This hook runs early in the configuration stage to ensure the jinja environment</span>\n<span
    class=\"sd\">    is available for other plugins that need it during configuration.</span>\n\n<span
    class=\"sd\">    Args:</span>\n<span class=\"sd\">        markata: The Markata
    instance</span>\n<span class=\"sd\">    &quot;&quot;&quot;</span>\n    <span class=\"c1\">#
    Get configuration, falling back to defaults</span>\n    <span class=\"n\">config</span>
    <span class=\"o\">=</span> <span class=\"n\">JinjaEnvConfig</span><span class=\"p\">()</span>\n
    \   <span class=\"k\">if</span> <span class=\"nb\">hasattr</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;jinja_env&quot;</span><span class=\"p\">):</span>\n
    \       <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">jinja_env</span><span
    class=\"p\">,</span> <span class=\"nb\">dict</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">config</span> <span class=\"o\">=</span> <span class=\"n\">JinjaEnvConfig</span><span
    class=\"p\">(</span><span class=\"o\">**</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">jinja_env</span><span class=\"p\">)</span>\n\n    <span class=\"c1\">#
    TODO: setting up env twice could not get dynamic templates to be recognized on
    first pass</span>\n    <span class=\"n\">loaders</span> <span class=\"o\">=</span>
    <span class=\"p\">[]</span>\n    <span class=\"k\">if</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">templates_dir</span><span class=\"p\">:</span>\n        <span class=\"k\">for</span>
    <span class=\"n\">path</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">templates_dir</span><span class=\"p\">:</span>\n            <span
    class=\"n\">path</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">path</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">expanduser</span><span class=\"p\">()</span><span
    class=\"o\">.</span><span class=\"n\">resolve</span><span class=\"p\">()</span>\n
    \           <span class=\"k\">if</span> <span class=\"n\">path</span><span class=\"o\">.</span><span
    class=\"n\">exists</span><span class=\"p\">():</span>\n                <span class=\"n\">loaders</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">FileSystemLoader</span><span class=\"p\">(</span><span class=\"nb\">str</span><span
    class=\"p\">(</span><span class=\"n\">path</span><span class=\"p\">)))</span>\n
    \   <span class=\"c1\"># Create environment</span>\n    <span class=\"n\">env_for_dynamic_render</span>
    <span class=\"o\">=</span> <span class=\"n\">Environment</span><span class=\"p\">(</span>\n
    \       <span class=\"n\">loader</span><span class=\"o\">=</span><span class=\"n\">ChoiceLoader</span><span
    class=\"p\">(</span><span class=\"n\">loaders</span><span class=\"p\">),</span>\n
    \       <span class=\"n\">undefined</span><span class=\"o\">=</span><span class=\"n\">_SilentUndefined</span>
    <span class=\"k\">if</span> <span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">undefined_silent</span> <span class=\"k\">else</span> <span class=\"n\">jinja2</span><span
    class=\"o\">.</span><span class=\"n\">Undefined</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">trim_blocks</span><span class=\"o\">=</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">trim_blocks</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">lstrip_blocks</span><span class=\"o\">=</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">lstrip_blocks</span><span
    class=\"p\">,</span>\n        <span class=\"n\">bytecode_cache</span><span class=\"o\">=</span><span
    class=\"n\">MarkataTemplateCache</span><span class=\"p\">(</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">template_cache_dir</span><span class=\"p\">),</span>\n
    \       <span class=\"n\">auto_reload</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span>\n    <span class=\"p\">)</span>\n\n    <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">dynamic_templates_dir</span><span class=\"o\">.</span><span class=\"n\">mkdir</span><span
    class=\"p\">(</span><span class=\"n\">parents</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span> <span class=\"n\">exist_ok</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \   <span class=\"n\">head_template</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">dynamic_templates_dir</span> <span class=\"o\">/</span> <span class=\"s2\">&quot;head.html&quot;</span>\n
    \   <span class=\"n\">head_template</span><span class=\"o\">.</span><span class=\"n\">write_text</span><span
    class=\"p\">(</span>\n        <span class=\"n\">env_for_dynamic_render</span><span
    class=\"o\">.</span><span class=\"n\">get_template</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;dynamic_head.html&quot;</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">render</span><span class=\"p\">(</span>\n
    \           <span class=\"p\">{</span><span class=\"s2\">&quot;markata&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">markata</span><span class=\"p\">}</span>\n
    \       <span class=\"p\">),</span>\n    <span class=\"p\">)</span>\n\n    <span
    class=\"c1\"># Set up loaders</span>\n    <span class=\"n\">loaders</span> <span
    class=\"o\">=</span> <span class=\"p\">[]</span>\n\n    <span class=\"c1\"># Add
    package templates first (lowest priority)</span>\n    <span class=\"c1\"># loaders.append(PackageLoader(&quot;markata&quot;,
    &quot;templates&quot;))</span>\n\n    <span class=\"c1\"># Add user template paths
    (medium priority)</span>\n    <span class=\"k\">if</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">templates_dir</span><span class=\"p\">:</span>\n        <span class=\"k\">for</span>
    <span class=\"n\">path</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">templates_dir</span><span class=\"p\">:</span>\n            <span
    class=\"n\">path</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">path</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">expanduser</span><span class=\"p\">()</span><span
    class=\"o\">.</span><span class=\"n\">resolve</span><span class=\"p\">()</span>\n
    \           <span class=\"k\">if</span> <span class=\"n\">path</span><span class=\"o\">.</span><span
    class=\"n\">exists</span><span class=\"p\">():</span>\n                <span class=\"n\">loaders</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">FileSystemLoader</span><span class=\"p\">(</span><span class=\"nb\">str</span><span
    class=\"p\">(</span><span class=\"n\">path</span><span class=\"p\">)))</span>\n\n
    \   <span class=\"c1\"># Add dynamic templates directory (highest priority)</span>\n
    \   <span class=\"c1\"># dynamic_templates_dir = Path(&quot;.markata.cache/templates&quot;)</span>\n
    \   <span class=\"c1\"># dynamic_templates_dir.mkdir(parents=True, exist_ok=True)</span>\n
    \   <span class=\"c1\"># loaders.append(FileSystemLoader(str(dynamic_templates_dir)))</span>\n\n
    \   <span class=\"c1\"># Create environment</span>\n    <span class=\"n\">env</span>
    <span class=\"o\">=</span> <span class=\"n\">Environment</span><span class=\"p\">(</span>\n
    \       <span class=\"n\">loader</span><span class=\"o\">=</span><span class=\"n\">ChoiceLoader</span><span
    class=\"p\">(</span><span class=\"n\">loaders</span><span class=\"p\">),</span>\n
    \       <span class=\"n\">undefined</span><span class=\"o\">=</span><span class=\"n\">_SilentUndefined</span>
    <span class=\"k\">if</span> <span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">undefined_silent</span> <span class=\"k\">else</span> <span class=\"n\">jinja2</span><span
    class=\"o\">.</span><span class=\"n\">Undefined</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">trim_blocks</span><span class=\"o\">=</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">trim_blocks</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">lstrip_blocks</span><span class=\"o\">=</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">lstrip_blocks</span><span
    class=\"p\">,</span>\n        <span class=\"n\">bytecode_cache</span><span class=\"o\">=</span><span
    class=\"n\">MarkataTemplateCache</span><span class=\"p\">(</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">template_cache_dir</span><span class=\"p\">),</span>\n
    \       <span class=\"n\">auto_reload</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span>\n    <span class=\"p\">)</span>\n\n    <span class=\"c1\">#
    Register the environment on the config&#39;s private attribute</span>\n    <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">jinja_env</span>
    <span class=\"o\">=</span> <span class=\"n\">env</span>\n</pre></div>\n\n</pre>\n\n</div>\n\n
    \   </section>\n</article>    </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>jinja_env.py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"Jinja2 Environment Plugin Provides a centralized
    Jinja2 environment configuration for consistent template rendering across all
    Markata plugins. This plugin\u2026\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n
    \       <meta property=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    <head>\n        <script src=\"https://unpkg.com/@tailwindcss/browser@4\"></script>\n<title>jinja_env.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Jinja2 Environment Plugin Provides a
    centralized Jinja2 environment configuration for consistent template rendering
    across all Markata plugins. This plugin\u2026\" />\n <link href=\"/favicon.ico\"
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
    class=\"title\">\n    <h1 id=\"title\">\n        jinja_env.py\n    </h1>\n</section></article>\n
    \    </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       jinja_env.py\n    </h1>\n</section>    <section class=\"body\">\n        <hr
    />\n<p>Jinja2 Environment Plugin</p>\n<p>Provides a centralized Jinja2 environment
    configuration for consistent template rendering\nacross all Markata plugins. This
    plugin ensures template rendering behavior is consistent\nand available even when
    specific template-using plugins are not enabled.</p>\n<h2 id=\"installation\">Installation
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
    class=\"p\">[</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;markata.plugins.jinja_env&quot;</span><span
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
    class=\"p\">[</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;markata.plugins.jinja_env&quot;</span><span
    class=\"p\">,</span>\n<span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<h2
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Configure Jinja environment
    settings in your <code>markata.toml</code>:</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata.jinja_env]</span>\n<span
    class=\"n\">template_paths</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"p\">[</span><span class=\"s2\">&quot;templates&quot;</span><span
    class=\"p\">]</span><span class=\"w\">  </span><span class=\"c1\"># Additional
    template paths to search</span>\n<span class=\"n\">undefined_silent</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"kc\">true</span><span class=\"w\">        </span><span class=\"c1\">#
    Return empty string for undefined variables</span>\n<span class=\"n\">trim_blocks</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"kc\">true</span><span class=\"w\">            </span><span class=\"c1\">#
    Remove first newline after block</span>\n<span class=\"n\">lstrip_blocks</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"kc\">true</span><span class=\"w\">          </span><span class=\"c1\">#
    Strip tabs/spaces from start of line</span>\n<span class=\"n\">template_cache_dir</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;.markata.cache/template_bytecode&quot;</span>\n</pre></div>\n\n</pre>\n\n<h1
    id=\"usage\">Usage <a class=\"header-anchor\" href=\"#usage\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>The environment is automatically
    available to other plugins via <code>markata.config.jinja_env</code>.\nTemplate
    loading follows this order:</p>\n<ol>\n<li>Package templates (built-in Markata
    templates)</li>\n<li>User template paths (configured via template_paths)</li>\n</ol>\n<p>Example
    usage in a plugin:</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    class=\"w\"> </span><span class=\"nf\">render_template</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">,</span> <span class=\"n\">content</span><span
    class=\"p\">):</span>\n    <span class=\"n\">template</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">jinja_env</span><span
    class=\"o\">.</span><span class=\"n\">from_string</span><span class=\"p\">(</span><span
    class=\"n\">content</span><span class=\"p\">)</span>\n    <span class=\"k\">return</span>
    <span class=\"n\">template</span><span class=\"o\">.</span><span class=\"n\">render</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<h1
    id=\"notes\">Notes <a class=\"header-anchor\" href=\"#notes\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<ul>\n<li>Template paths
    are resolved relative to the current working directory</li>\n<li>Package templates
    are always available and take precedence</li>\n<li>Silent undefined behavior means
    undefined variables render as empty strings</li>\n</ul>\n<hr />\n<div class=\"admonition
    class\">\n<p class=\"admonition-title\">Class</p>\n<h2 id=\"_SilentUndefined\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">_SilentUndefined
    <em class=\"small\">class</em></h2>\n<p>Custom undefined type that returns empty
    string for undefined variables.</p>\n</div>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">_SilentUndefined <em class='small'>source</em></p>\n<pre
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
    class=\"w\"> </span><span class=\"nc\">_SilentUndefined</span><span class=\"p\">(</span><span
    class=\"n\">jinja2</span><span class=\"o\">.</span><span class=\"n\">Undefined</span><span
    class=\"p\">):</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Custom
    undefined type that returns empty string for undefined variables.&quot;&quot;&quot;</span>\n\n
    \   <span class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">_fail_with_undefined_error</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"n\">args</span><span class=\"p\">,</span> <span
    class=\"o\">**</span><span class=\"n\">kwargs</span><span class=\"p\">):</span>\n
    \       <span class=\"k\">return</span> <span class=\"s2\">&quot;&quot;</span>\n\n
    \   <span class=\"fm\">__add__</span> <span class=\"o\">=</span> <span class=\"fm\">__radd__</span>
    <span class=\"o\">=</span> <span class=\"fm\">__mul__</span> <span class=\"o\">=</span>
    <span class=\"fm\">__rmul__</span> <span class=\"o\">=</span> <span class=\"n\">__div__</span>
    <span class=\"o\">=</span> <span class=\"n\">__rdiv__</span> <span class=\"o\">=</span>
    <span class=\"fm\">__truediv__</span> <span class=\"o\">=</span> <span class=\"p\">(</span>\n
    \       <span class=\"fm\">__rtruediv__</span>\n    <span class=\"p\">)</span>
    <span class=\"o\">=</span> <span class=\"fm\">__floordiv__</span> <span class=\"o\">=</span>
    <span class=\"fm\">__rfloordiv__</span> <span class=\"o\">=</span> <span class=\"fm\">__mod__</span>
    <span class=\"o\">=</span> <span class=\"fm\">__rmod__</span> <span class=\"o\">=</span>
    <span class=\"fm\">__pos__</span> <span class=\"o\">=</span> <span class=\"fm\">__neg__</span>
    <span class=\"o\">=</span> <span class=\"p\">(</span>\n        <span class=\"fm\">__call__</span>\n
    \   <span class=\"p\">)</span> <span class=\"o\">=</span> <span class=\"fm\">__getitem__</span>
    <span class=\"o\">=</span> <span class=\"fm\">__lt__</span> <span class=\"o\">=</span>
    <span class=\"fm\">__le__</span> <span class=\"o\">=</span> <span class=\"fm\">__gt__</span>
    <span class=\"o\">=</span> <span class=\"fm\">__ge__</span> <span class=\"o\">=</span>
    <span class=\"fm\">__int__</span> <span class=\"o\">=</span> <span class=\"fm\">__float__</span>
    <span class=\"o\">=</span> <span class=\"p\">(</span>\n        <span class=\"fm\">__complex__</span>\n
    \   <span class=\"p\">)</span> <span class=\"o\">=</span> <span class=\"fm\">__pow__</span>
    <span class=\"o\">=</span> <span class=\"fm\">__rpow__</span> <span class=\"o\">=</span>
    <span class=\"n\">_fail_with_undefined_error</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition class\">\n<p class=\"admonition-title\">Class</p>\n<h2 id=\"MarkataTemplateCache\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">MarkataTemplateCache
    <em class=\"small\">class</em></h2>\n<p>Template bytecode cache for improved performance.</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">MarkataTemplateCache
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
    class=\"w\"> </span><span class=\"nc\">MarkataTemplateCache</span><span class=\"p\">(</span><span
    class=\"n\">jinja2</span><span class=\"o\">.</span><span class=\"n\">BytecodeCache</span><span
    class=\"p\">):</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Template
    bytecode cache for improved performance.&quot;&quot;&quot;</span>\n\n    <span
    class=\"k\">def</span><span class=\"w\"> </span><span class=\"fm\">__init__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"n\">directory</span><span class=\"p\">):</span>\n        <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">directory</span> <span class=\"o\">=</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">directory</span><span
    class=\"p\">)</span>\n        <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">directory</span><span class=\"o\">.</span><span class=\"n\">mkdir</span><span
    class=\"p\">(</span><span class=\"n\">parents</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span> <span class=\"n\">exist_ok</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n\n
    \   <span class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">load_bytecode</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"n\">bucket</span><span class=\"p\">):</span>\n        <span class=\"n\">filename</span>
    <span class=\"o\">=</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">directory</span> <span class=\"o\">/</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;</span><span class=\"si\">{</span><span class=\"n\">bucket</span><span
    class=\"o\">.</span><span class=\"n\">key</span><span class=\"si\">}</span><span
    class=\"s2\">.cache&quot;</span>\n        <span class=\"k\">if</span> <span class=\"n\">filename</span><span
    class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">():</span>\n
    \           <span class=\"k\">with</span> <span class=\"nb\">open</span><span
    class=\"p\">(</span><span class=\"n\">filename</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;rb&quot;</span><span class=\"p\">)</span> <span class=\"k\">as</span>
    <span class=\"n\">f</span><span class=\"p\">:</span>\n                <span class=\"n\">bucket</span><span
    class=\"o\">.</span><span class=\"n\">bytecode_from_string</span><span class=\"p\">(</span><span
    class=\"n\">f</span><span class=\"o\">.</span><span class=\"n\">read</span><span
    class=\"p\">())</span>\n\n    <span class=\"k\">def</span><span class=\"w\"> </span><span
    class=\"nf\">dump_bytecode</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">bucket</span><span class=\"p\">):</span>\n
    \       <span class=\"n\">filename</span> <span class=\"o\">=</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">directory</span> <span class=\"o\">/</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">bucket</span><span class=\"o\">.</span><span class=\"n\">key</span><span
    class=\"si\">}</span><span class=\"s2\">.cache&quot;</span>\n        <span class=\"k\">with</span>
    <span class=\"nb\">open</span><span class=\"p\">(</span><span class=\"n\">filename</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;wb&quot;</span><span class=\"p\">)</span>
    <span class=\"k\">as</span> <span class=\"n\">f</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">f</span><span class=\"o\">.</span><span class=\"n\">write</span><span
    class=\"p\">(</span><span class=\"n\">bucket</span><span class=\"o\">.</span><span
    class=\"n\">bytecode_to_string</span><span class=\"p\">())</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition class\">\n<p class=\"admonition-title\">Class</p>\n<h2 id=\"JinjaEnvConfig\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">JinjaEnvConfig
    <em class=\"small\">class</em></h2>\n<p>Configuration for the Jinja environment.</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">JinjaEnvConfig
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
    class=\"w\"> </span><span class=\"nc\">JinjaEnvConfig</span><span class=\"p\">(</span><span
    class=\"n\">pydantic</span><span class=\"o\">.</span><span class=\"n\">BaseModel</span><span
    class=\"p\">):</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Configuration
    for the Jinja environment.&quot;&quot;&quot;</span>\n\n    <span class=\"n\">templates_dir</span><span
    class=\"p\">:</span> <span class=\"n\">List</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"p\">[]</span>\n    <span class=\"n\">undefined_silent</span><span class=\"p\">:</span>
    <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"kc\">True</span>\n
    \   <span class=\"n\">trim_blocks</span><span class=\"p\">:</span> <span class=\"nb\">bool</span>
    <span class=\"o\">=</span> <span class=\"kc\">True</span>\n    <span class=\"n\">lstrip_blocks</span><span
    class=\"p\">:</span> <span class=\"nb\">bool</span> <span class=\"o\">=</span>
    <span class=\"kc\">True</span>\n    <span class=\"n\">template_cache_dir</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span> <span class=\"o\">=</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"s2\">&quot;.markata.cache/template_bytecode&quot;</span><span
    class=\"p\">)</span>\n\n    <span class=\"n\">model_config</span> <span class=\"o\">=</span>
    <span class=\"n\">pydantic</span><span class=\"o\">.</span><span class=\"n\">ConfigDict</span><span
    class=\"p\">(</span>\n        <span class=\"n\">validate_assignment</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>  <span
    class=\"c1\"># Config model</span>\n        <span class=\"n\">arbitrary_types_allowed</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">extra</span><span class=\"o\">=</span><span class=\"s2\">&quot;allow&quot;</span><span
    class=\"p\">,</span>\n        <span class=\"n\">str_strip_whitespace</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">validate_default</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n        <span class=\"n\">coerce_numbers_to_str</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">populate_by_name</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n    <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition function\">\n<p class=\"admonition-title\">Function</p>\n<h2
    id=\"config_model\" class=\"admonition-title\" style=\"margin: 0; padding: .5rem
    1rem;\">config_model <em class=\"small\">function</em></h2>\n<p>Register configuration
    models.</p>\n</div>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">config_model <em class='small'>source</em></p>\n<pre
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
    class=\"w\"> </span><span class=\"nf\">config_model</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Register
    configuration models.&quot;&quot;&quot;</span>\n    <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config_models</span><span class=\"o\">.</span><span
    class=\"n\">append</span><span class=\"p\">(</span><span class=\"n\">JinjaEnvConfig</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div class=\"admonition
    function\">\n<p class=\"admonition-title\">Function</p>\n<h2 id=\"configure\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">configure
    <em class=\"small\">function</em></h2>\n<p>Initialize and configure the Jinja2
    environment for Markata.</p>\n<p>This hook runs early in the configuration stage
    to ensure the jinja environment\nis available for other plugins that need it during
    configuration.</p>\n<p>Args:\nmarkata: The Markata instance</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">configure
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
    class=\"w\"> </span><span class=\"nf\">configure</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"n\">Markata</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Initialize
    and configure the Jinja2 environment for Markata.</span>\n\n<span class=\"sd\">
    \   This hook runs early in the configuration stage to ensure the jinja environment</span>\n<span
    class=\"sd\">    is available for other plugins that need it during configuration.</span>\n\n<span
    class=\"sd\">    Args:</span>\n<span class=\"sd\">        markata: The Markata
    instance</span>\n<span class=\"sd\">    &quot;&quot;&quot;</span>\n    <span class=\"c1\">#
    Get configuration, falling back to defaults</span>\n    <span class=\"n\">config</span>
    <span class=\"o\">=</span> <span class=\"n\">JinjaEnvConfig</span><span class=\"p\">()</span>\n
    \   <span class=\"k\">if</span> <span class=\"nb\">hasattr</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;jinja_env&quot;</span><span class=\"p\">):</span>\n
    \       <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">jinja_env</span><span
    class=\"p\">,</span> <span class=\"nb\">dict</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">config</span> <span class=\"o\">=</span> <span class=\"n\">JinjaEnvConfig</span><span
    class=\"p\">(</span><span class=\"o\">**</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">jinja_env</span><span class=\"p\">)</span>\n\n    <span class=\"c1\">#
    TODO: setting up env twice could not get dynamic templates to be recognized on
    first pass</span>\n    <span class=\"n\">loaders</span> <span class=\"o\">=</span>
    <span class=\"p\">[]</span>\n    <span class=\"k\">if</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">templates_dir</span><span class=\"p\">:</span>\n        <span class=\"k\">for</span>
    <span class=\"n\">path</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">templates_dir</span><span class=\"p\">:</span>\n            <span
    class=\"n\">path</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">path</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">expanduser</span><span class=\"p\">()</span><span
    class=\"o\">.</span><span class=\"n\">resolve</span><span class=\"p\">()</span>\n
    \           <span class=\"k\">if</span> <span class=\"n\">path</span><span class=\"o\">.</span><span
    class=\"n\">exists</span><span class=\"p\">():</span>\n                <span class=\"n\">loaders</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">FileSystemLoader</span><span class=\"p\">(</span><span class=\"nb\">str</span><span
    class=\"p\">(</span><span class=\"n\">path</span><span class=\"p\">)))</span>\n
    \   <span class=\"c1\"># Create environment</span>\n    <span class=\"n\">env_for_dynamic_render</span>
    <span class=\"o\">=</span> <span class=\"n\">Environment</span><span class=\"p\">(</span>\n
    \       <span class=\"n\">loader</span><span class=\"o\">=</span><span class=\"n\">ChoiceLoader</span><span
    class=\"p\">(</span><span class=\"n\">loaders</span><span class=\"p\">),</span>\n
    \       <span class=\"n\">undefined</span><span class=\"o\">=</span><span class=\"n\">_SilentUndefined</span>
    <span class=\"k\">if</span> <span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">undefined_silent</span> <span class=\"k\">else</span> <span class=\"n\">jinja2</span><span
    class=\"o\">.</span><span class=\"n\">Undefined</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">trim_blocks</span><span class=\"o\">=</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">trim_blocks</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">lstrip_blocks</span><span class=\"o\">=</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">lstrip_blocks</span><span
    class=\"p\">,</span>\n        <span class=\"n\">bytecode_cache</span><span class=\"o\">=</span><span
    class=\"n\">MarkataTemplateCache</span><span class=\"p\">(</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">template_cache_dir</span><span class=\"p\">),</span>\n
    \       <span class=\"n\">auto_reload</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span>\n    <span class=\"p\">)</span>\n\n    <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">dynamic_templates_dir</span><span class=\"o\">.</span><span class=\"n\">mkdir</span><span
    class=\"p\">(</span><span class=\"n\">parents</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span> <span class=\"n\">exist_ok</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \   <span class=\"n\">head_template</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">dynamic_templates_dir</span> <span class=\"o\">/</span> <span class=\"s2\">&quot;head.html&quot;</span>\n
    \   <span class=\"n\">head_template</span><span class=\"o\">.</span><span class=\"n\">write_text</span><span
    class=\"p\">(</span>\n        <span class=\"n\">env_for_dynamic_render</span><span
    class=\"o\">.</span><span class=\"n\">get_template</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;dynamic_head.html&quot;</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">render</span><span class=\"p\">(</span>\n
    \           <span class=\"p\">{</span><span class=\"s2\">&quot;markata&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">markata</span><span class=\"p\">}</span>\n
    \       <span class=\"p\">),</span>\n    <span class=\"p\">)</span>\n\n    <span
    class=\"c1\"># Set up loaders</span>\n    <span class=\"n\">loaders</span> <span
    class=\"o\">=</span> <span class=\"p\">[]</span>\n\n    <span class=\"c1\"># Add
    package templates first (lowest priority)</span>\n    <span class=\"c1\"># loaders.append(PackageLoader(&quot;markata&quot;,
    &quot;templates&quot;))</span>\n\n    <span class=\"c1\"># Add user template paths
    (medium priority)</span>\n    <span class=\"k\">if</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">templates_dir</span><span class=\"p\">:</span>\n        <span class=\"k\">for</span>
    <span class=\"n\">path</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">templates_dir</span><span class=\"p\">:</span>\n            <span
    class=\"n\">path</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">path</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">expanduser</span><span class=\"p\">()</span><span
    class=\"o\">.</span><span class=\"n\">resolve</span><span class=\"p\">()</span>\n
    \           <span class=\"k\">if</span> <span class=\"n\">path</span><span class=\"o\">.</span><span
    class=\"n\">exists</span><span class=\"p\">():</span>\n                <span class=\"n\">loaders</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">FileSystemLoader</span><span class=\"p\">(</span><span class=\"nb\">str</span><span
    class=\"p\">(</span><span class=\"n\">path</span><span class=\"p\">)))</span>\n\n
    \   <span class=\"c1\"># Add dynamic templates directory (highest priority)</span>\n
    \   <span class=\"c1\"># dynamic_templates_dir = Path(&quot;.markata.cache/templates&quot;)</span>\n
    \   <span class=\"c1\"># dynamic_templates_dir.mkdir(parents=True, exist_ok=True)</span>\n
    \   <span class=\"c1\"># loaders.append(FileSystemLoader(str(dynamic_templates_dir)))</span>\n\n
    \   <span class=\"c1\"># Create environment</span>\n    <span class=\"n\">env</span>
    <span class=\"o\">=</span> <span class=\"n\">Environment</span><span class=\"p\">(</span>\n
    \       <span class=\"n\">loader</span><span class=\"o\">=</span><span class=\"n\">ChoiceLoader</span><span
    class=\"p\">(</span><span class=\"n\">loaders</span><span class=\"p\">),</span>\n
    \       <span class=\"n\">undefined</span><span class=\"o\">=</span><span class=\"n\">_SilentUndefined</span>
    <span class=\"k\">if</span> <span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">undefined_silent</span> <span class=\"k\">else</span> <span class=\"n\">jinja2</span><span
    class=\"o\">.</span><span class=\"n\">Undefined</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">trim_blocks</span><span class=\"o\">=</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">trim_blocks</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">lstrip_blocks</span><span class=\"o\">=</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">lstrip_blocks</span><span
    class=\"p\">,</span>\n        <span class=\"n\">bytecode_cache</span><span class=\"o\">=</span><span
    class=\"n\">MarkataTemplateCache</span><span class=\"p\">(</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">template_cache_dir</span><span class=\"p\">),</span>\n
    \       <span class=\"n\">auto_reload</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span>\n    <span class=\"p\">)</span>\n\n    <span class=\"c1\">#
    Register the environment on the config&#39;s private attribute</span>\n    <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">jinja_env</span>
    <span class=\"o\">=</span> <span class=\"n\">env</span>\n</pre></div>\n\n</pre>\n\n</div>\n\n
    \   </section>\n</article>"
  raw.md: ''
published: false
slug: markata/plugins/jinja-env
title: jinja_env.py


---

---

Jinja2 Environment Plugin

Provides a centralized Jinja2 environment configuration for consistent template rendering
across all Markata plugins. This plugin ensures template rendering behavior is consistent
and available even when specific template-using plugins are not enabled.

## Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

``` toml
hooks = [
    "markata.plugins.jinja_env",
]
```

## Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

``` toml
disabled_hooks = [
    "markata.plugins.jinja_env",
]
```

## Configuration

Configure Jinja environment settings in your `markata.toml`:

``` toml
[markata.jinja_env]
template_paths = ["templates"]  # Additional template paths to search
undefined_silent = true        # Return empty string for undefined variables
trim_blocks = true            # Remove first newline after block
lstrip_blocks = true          # Strip tabs/spaces from start of line
template_cache_dir = ".markata.cache/template_bytecode"
```

# Usage

The environment is automatically available to other plugins via `markata.config.jinja_env`.
Template loading follows this order:
1. Package templates (built-in Markata templates)
2. User template paths (configured via template_paths)

Example usage in a plugin:

``` python
def render_template(markata, content):
    template = markata.jinja_env.from_string(content)
    return template.render(markata=markata)
```

# Notes

- Template paths are resolved relative to the current working directory
- Package templates are always available and take precedence
- Silent undefined behavior means undefined variables render as empty strings

---

!!! class
    <h2 id="_SilentUndefined" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">_SilentUndefined <em class="small">class</em></h2>

    Custom undefined type that returns empty string for undefined variables.

???+ source "_SilentUndefined <em class='small'>source</em>"
    ```python
    class _SilentUndefined(jinja2.Undefined):
        """Custom undefined type that returns empty string for undefined variables."""

        def _fail_with_undefined_error(self, *args, **kwargs):
            return ""

        __add__ = __radd__ = __mul__ = __rmul__ = __div__ = __rdiv__ = __truediv__ = (
            __rtruediv__
        ) = __floordiv__ = __rfloordiv__ = __mod__ = __rmod__ = __pos__ = __neg__ = (
            __call__
        ) = __getitem__ = __lt__ = __le__ = __gt__ = __ge__ = __int__ = __float__ = (
            __complex__
        ) = __pow__ = __rpow__ = _fail_with_undefined_error
    ```
!!! class
    <h2 id="MarkataTemplateCache" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">MarkataTemplateCache <em class="small">class</em></h2>

    Template bytecode cache for improved performance.

???+ source "MarkataTemplateCache <em class='small'>source</em>"
    ```python
    class MarkataTemplateCache(jinja2.BytecodeCache):
        """Template bytecode cache for improved performance."""

        def __init__(self, directory):
            self.directory = Path(directory)
            self.directory.mkdir(parents=True, exist_ok=True)

        def load_bytecode(self, bucket):
            filename = self.directory / f"{bucket.key}.cache"
            if filename.exists():
                with open(filename, "rb") as f:
                    bucket.bytecode_from_string(f.read())

        def dump_bytecode(self, bucket):
            filename = self.directory / f"{bucket.key}.cache"
            with open(filename, "wb") as f:
                f.write(bucket.bytecode_to_string())
    ```
!!! class
    <h2 id="JinjaEnvConfig" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">JinjaEnvConfig <em class="small">class</em></h2>

    Configuration for the Jinja environment.

???+ source "JinjaEnvConfig <em class='small'>source</em>"
    ```python
    class JinjaEnvConfig(pydantic.BaseModel):
        """Configuration for the Jinja environment."""

        templates_dir: List[str] = []
        undefined_silent: bool = True
        trim_blocks: bool = True
        lstrip_blocks: bool = True
        template_cache_dir: Path = Path(".markata.cache/template_bytecode")

        model_config = pydantic.ConfigDict(
            validate_assignment=True,  # Config model
            arbitrary_types_allowed=True,
            extra="allow",
            str_strip_whitespace=True,
            validate_default=True,
            coerce_numbers_to_str=True,
            populate_by_name=True,
        )
    ```
!!! function
    <h2 id="config_model" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">config_model <em class="small">function</em></h2>

    Register configuration models.

???+ source "config_model <em class='small'>source</em>"
    ```python
    def config_model(markata: "Markata") -> None:
        """Register configuration models."""
        markata.config_models.append(JinjaEnvConfig)
    ```
!!! function
    <h2 id="configure" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">configure <em class="small">function</em></h2>

    Initialize and configure the Jinja2 environment for Markata.

    This hook runs early in the configuration stage to ensure the jinja environment
    is available for other plugins that need it during configuration.

    Args:
        markata: The Markata instance

???+ source "configure <em class='small'>source</em>"
    ```python
    def configure(markata: Markata) -> None:
        """Initialize and configure the Jinja2 environment for Markata.

        This hook runs early in the configuration stage to ensure the jinja environment
        is available for other plugins that need it during configuration.

        Args:
            markata: The Markata instance
        """
        # Get configuration, falling back to defaults
        config = JinjaEnvConfig()
        if hasattr(markata.config, "jinja_env"):
            if isinstance(markata.config.jinja_env, dict):
                config = JinjaEnvConfig(**markata.config.jinja_env)

        # TODO: setting up env twice could not get dynamic templates to be recognized on first pass
        loaders = []
        if markata.config.templates_dir:
            for path in markata.config.templates_dir:
                path = Path(path).expanduser().resolve()
                if path.exists():
                    loaders.append(FileSystemLoader(str(path)))
        # Create environment
        env_for_dynamic_render = Environment(
            loader=ChoiceLoader(loaders),
            undefined=_SilentUndefined if config.undefined_silent else jinja2.Undefined,
            trim_blocks=config.trim_blocks,
            lstrip_blocks=config.lstrip_blocks,
            bytecode_cache=MarkataTemplateCache(config.template_cache_dir),
            auto_reload=True,
        )

        markata.config.dynamic_templates_dir.mkdir(parents=True, exist_ok=True)
        head_template = markata.config.dynamic_templates_dir / "head.html"
        head_template.write_text(
            env_for_dynamic_render.get_template("dynamic_head.html").render(
                {"markata": markata}
            ),
        )

        # Set up loaders
        loaders = []

        # Add package templates first (lowest priority)
        # loaders.append(PackageLoader("markata", "templates"))

        # Add user template paths (medium priority)
        if markata.config.templates_dir:
            for path in markata.config.templates_dir:
                path = Path(path).expanduser().resolve()
                if path.exists():
                    loaders.append(FileSystemLoader(str(path)))

        # Add dynamic templates directory (highest priority)
        # dynamic_templates_dir = Path(".markata.cache/templates")
        # dynamic_templates_dir.mkdir(parents=True, exist_ok=True)
        # loaders.append(FileSystemLoader(str(dynamic_templates_dir)))

        # Create environment
        env = Environment(
            loader=ChoiceLoader(loaders),
            undefined=_SilentUndefined if config.undefined_silent else jinja2.Undefined,
            trim_blocks=config.trim_blocks,
            lstrip_blocks=config.lstrip_blocks,
            bytecode_cache=MarkataTemplateCache(config.template_cache_dir),
            auto_reload=True,
        )

        # Register the environment on the config's private attribute
        markata.jinja_env = env
    ```