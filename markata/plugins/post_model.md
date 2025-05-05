---
content: "---\n\nThe `markata.plugins.post_model` plugin defines the core Post model
  used throughout\nMarkata. It provides robust validation, serialization, and configuration
  options for\nall post attributes.\n\n## Installation\n\nThis plugin is built-in
  and enabled by default through the 'default' plugin.\nIf you want to be explicit,
  you can add it to your list of plugins:\n\n```toml\nhooks = [\n    \"markata.plugins.post_model\",\n]\n```\n\n##
  Uninstallation\n\nSince this plugin is included in the default plugin set, to disable
  it you must explicitly\nadd it to the disabled_hooks list if you are using the 'default'
  plugin:\n\n```toml\ndisabled_hooks = [\n    \"markata.plugins.post_model\",\n]\n```\n\nNote:
  Disabling this plugin will break most of Markata's functionality as the Post\nmodel
  is fundamental to the system.\n\n## Configuration\n\nConfigure post model behavior
  in your `markata.toml`:\n\n```toml\n[markata.post_model]\n# Attributes to include
  when serializing posts\ninclude = [\n    \"date\",\n    \"description\",\n    \"published\",\n
  \   \"slug\",\n    \"title\",\n    \"content\",\n    \"html\"\n]\n\n# Attributes
  to show in post representations\nrepr_include = [\n    \"date\",\n    \"description\",\n
  \   \"published\",\n    \"slug\",\n    \"title\"\n]\n\n# Attributes to include when
  exporting\nexport_include = [\n    \"date\",\n    \"description\",\n    \"published\",\n
  \   \"slug\",\n    \"title\"\n]\n```\n\n## Functionality\n\n## Post Model\n\nCore
  attributes:\n- `path`: Path to source file\n- `slug`: URL-friendly identifier\n-
  `href`: Full URL path\n- `published`: Publication status\n- `description`: Post
  summary\n- `content`: Raw markdown content\n- `html`: Rendered HTML content\n- `tags`:
  List of post tags\n- `date`: Publication date\n- `title`: Post title\n\n## Validation\n\nThe
  model provides:\n- Type checking and coercion\n- Required field validation\n- Custom
  field validators\n- Default values\n- Rich error messages\n\n## Serialization\n\nSupports
  multiple output formats:\n- Full serialization (all fields)\n- Representation (subset
  for display)\n- Export (subset for external use)\n- JSON/YAML compatible\n\n## Performance\n\nUses
  optimized Pydantic config:\n- Disabled assignment validation\n- Arbitrary types
  allowed\n- Extra fields allowed\n- String whitespace stripping\n- Default value
  validation\n- Number to string coercion\n- Alias population\n\n## Dependencies\n\nThis
  plugin depends on:\n- pydantic for model definition\n- rich for console output\n-
  yaml for YAML handling\n\n---\n\n!!! class\n    <h2 id=\"PostModelConfig\" class=\"admonition-title\"
  style=\"margin: 0; padding: .5rem 1rem;\">PostModelConfig <em class=\"small\">class</em></h2>\n\n
  \   Configuration for the Post model\n\n???+ source \"PostModelConfig <em class='small'>source</em>\"\n
  \   ```python\n    class PostModelConfig(pydantic.BaseModel):\n        \"\"\"Configuration
  for the Post model\"\"\"\n\n        def __init__(self, **data) -> None:\n            \"\"\"\n\n
  \           include: post attributes to include by default in Post\n            model
  serialization.\n            repr_include: post attributes to include by default
  in Post\n            repr.  If `repr_include` is None, it will default to\n            `include`,
  but it is likely that you want less in the repr\n            than serialized output.\n\n
  \           example:\n\n            ``` toml title='markata.toml'\n            [markata.post_model]\n
  \           include = ['date', 'description', 'published',\n                'slug',
  'title', 'content', 'html']\n            repr_include = ['date', 'description',
  'published', 'slug', 'title']\n            ```\n            \"\"\"\n            super().__init__(**data)\n\n
  \       default_date: datetime.date = datetime.date.today()\n        include: List[str]
  = [\n            \"date\",\n            \"description\",\n            \"published\",\n
  \           \"slug\",\n            \"title\",\n            \"content\",\n            \"html\",\n
  \       ]\n        repr_include: Optional[List[str]] = [\n            \"date\",\n
  \           \"description\",\n            \"published\",\n            \"slug\",\n
  \           \"title\",\n        ]\n        export_include: Optional[List[str]] =
  [\n            \"date\",\n            \"description\",\n            \"published\",\n
  \           \"slug\",\n            \"title\",\n        ]\n\n        model_config
  = ConfigDict(\n            validate_assignment=True,  # Config model\n            arbitrary_types_allowed=True,\n
  \           extra=\"allow\",\n            str_strip_whitespace=True,\n            validate_default=True,\n
  \           coerce_numbers_to_str=True,\n            populate_by_name=True,\n        )\n\n
  \       @field_validator(\"repr_include\", mode=\"before\")\n        @classmethod\n
  \       def repr_include_validator(cls, v, info) -> Optional[List[str]]:\n            if
  v:\n                return v\n            return info.data.get(\"include\")\n    ```\n!!!
  method\n    <h2 id=\"metadata\" class=\"admonition-title\" style=\"margin: 0; padding:
  .5rem 1rem;\">metadata <em class=\"small\">method</em></h2>\n\n    for backwards
  compatability\n\n???+ source \"metadata <em class='small'>source</em>\"\n    ```python\n
  \   def metadata(self: \"Post\") -> Dict:\n            \"for backwards compatability\"\n
  \           return self.__dict__\n    ```\n!!! method\n    <h2 id=\"to_dict\" class=\"admonition-title\"
  style=\"margin: 0; padding: .5rem 1rem;\">to_dict <em class=\"small\">method</em></h2>\n\n
  \   for backwards compatability\n\n???+ source \"to_dict <em class='small'>source</em>\"\n
  \   ```python\n    def to_dict(self: \"Post\") -> Dict:\n            \"for backwards
  compatability\"\n            return self.__dict__\n    ```\n!!! method\n    <h2
  id=\"__getitem__\" class=\"admonition-title\" style=\"margin: 0; padding: .5rem
  1rem;\">__getitem__ <em class=\"small\">method</em></h2>\n\n    for backwards compatability\n\n???+
  source \"__getitem__ <em class='small'>source</em>\"\n    ```python\n    def __getitem__(self:
  \"Post\", item: str) -> Any:\n            \"for backwards compatability\"\n            return
  getattr(self, item)\n    ```\n!!! method\n    <h2 id=\"__setitem__\" class=\"admonition-title\"
  style=\"margin: 0; padding: .5rem 1rem;\">__setitem__ <em class=\"small\">method</em></h2>\n\n
  \   for backwards compatability\n\n???+ source \"__setitem__ <em class='small'>source</em>\"\n
  \   ```python\n    def __setitem__(self: \"Post\", key: str, item: Any) -> None:\n
  \           \"for backwards compatability\"\n            setattr(self, key, item)\n
  \   ```\n!!! method\n    <h2 id=\"get\" class=\"admonition-title\" style=\"margin:
  0; padding: .5rem 1rem;\">get <em class=\"small\">method</em></h2>\n\n    for backwards
  compatability\n\n???+ source \"get <em class='small'>source</em>\"\n    ```python\n
  \   def get(self: \"Post\", item: str, default: Any) -> Any:\n            \"for
  backwards compatability\"\n            return getattr(self, item, default)\n    ```\n!!!
  method\n    <h2 id=\"keys\" class=\"admonition-title\" style=\"margin: 0; padding:
  .5rem 1rem;\">keys <em class=\"small\">method</em></h2>\n\n    for backwards compatability\n\n???+
  source \"keys <em class='small'>source</em>\"\n    ```python\n    def keys(self:
  \"Post\") -> List[str]:\n            \"for backwards compatability\"\n            return
  self.__dict__.keys()\n    ```\n!!! method\n    <h2 id=\"yaml\" class=\"admonition-title\"
  style=\"margin: 0; padding: .5rem 1rem;\">yaml <em class=\"small\">method</em></h2>\n\n
  \   dump model to yaml\n\n???+ source \"yaml <em class='small'>source</em>\"\n    ```python\n
  \   def yaml(self: \"Post\") -> str:\n            \"\"\"\n            dump model
  to yaml\n            \"\"\"\n            import yaml\n\n            return yaml.dump(\n
  \               self.dict(\n                    include={i: True for i in self.markata.config.post_model.include}\n
  \               ),\n                Dumper=yaml.CDumper,\n            )\n    ```\n!!!
  method\n    <h2 id=\"markdown\" class=\"admonition-title\" style=\"margin: 0; padding:
  .5rem 1rem;\">markdown <em class=\"small\">method</em></h2>\n\n    dump model to
  markdown\n\n???+ source \"markdown <em class='small'>source</em>\"\n    ```python\n
  \   def markdown(self: \"Post\") -> str:\n            \"\"\"\n            dump model
  to markdown\n            \"\"\"\n\n            import yaml\n\n            frontmatter
  = yaml.dump(\n                self.dict(\n                    include={\n                        i:
  True\n                        for i in [\n                            _i\n                            for
  _i in self.markata.config.post_model.include\n                            if _i
  != \"content\"\n                        ]\n                    }\n                ),\n
  \               Dumper=yaml.CDumper,\n            )\n            post = \"---\\n\"\n
  \           post += frontmatter\n            post += \"---\\n\\n\"\n\n            if
  self.content:\n                post += self.content\n            return post\n    ```\n!!!
  method\n    <h2 id=\"dumps\" class=\"admonition-title\" style=\"margin: 0; padding:
  .5rem 1rem;\">dumps <em class=\"small\">method</em></h2>\n\n    dumps raw article
  back out\n\n???+ source \"dumps <em class='small'>source</em>\"\n    ```python\n
  \   def dumps(self):\n            \"\"\"\n            dumps raw article back out\n
  \           \"\"\"\n            return f\"---\\n{self.yaml()}\\n\\n---\\n\\n{self.content}\"\n
  \   ```\n!!! method\n    <h2 id=\"parse_date_time\" class=\"admonition-title\" style=\"margin:
  0; padding: .5rem 1rem;\">parse_date_time <em class=\"small\">method</em></h2>\n\n
  \   Single validator to handle all date_time parsing cases\n\n???+ source \"parse_date_time
  <em class='small'>source</em>\"\n    ```python\n    def parse_date_time(cls, v,
  info):\n            \"\"\"Single validator to handle all date_time parsing cases\"\"\"\n
  \           # If we have an explicit date_time value\n            if v is not None:\n
  \               if isinstance(v, datetime.datetime):\n                    return
  v\n                if isinstance(v, datetime.date):\n                    return
  datetime.datetime.combine(v, datetime.time.min)\n                if isinstance(v,
  str):\n                    try:\n                        # Try ISO format first\n
  \                       return datetime.datetime.fromisoformat(v.replace(\"Z\",
  \"+00:00\"))\n                    except ValueError:\n                        try:\n
  \                           return datetime.datetime.strptime(v, \"%Y-%m-%d %H:%M\")\n
  \                       except ValueError:\n                            try:\n                                return
  datetime.datetime.strptime(v, \"%Y-%m-%d\")\n                            except
  ValueError:\n                                # Try dateparser as last resort for
  explicit date_time\n                                import dateparser\n\n                                parsed
  = dateparser.parse(v)\n                                if parsed:\n                                    return
  parsed\n                                return datetime.datetime.now()\n\n            #
  Get the raw date string directly from raw_date field\n            raw_date = info.data.get(\"raw_date\")\n
  \           if raw_date and isinstance(raw_date, str):\n                try:\n                    #
  Try ISO format first\n                    return datetime.datetime.fromisoformat(raw_date.replace(\"Z\",
  \"+00:00\"))\n                except ValueError:\n                    try:\n                        #
  Try parsing raw_date with time first\n                        return datetime.datetime.strptime(raw_date,
  \"%Y-%m-%d %H:%M\")\n                    except ValueError:\n                        try:\n
  \                           # Fallback to date only\n                            return
  datetime.datetime.strptime(raw_date, \"%Y-%m-%d\")\n                        except
  ValueError:\n                            # Try dateparser as last resort\n                            import
  dateparser\n\n                            parsed = dateparser.parse(raw_date)\n
  \                           if parsed:\n                                return parsed\n\n
  \           # If no raw_date, try to derive from date field\n            date =
  info.data.get(\"date\")\n            if date:\n                if isinstance(date,
  datetime.datetime):\n                    return date\n                if isinstance(date,
  str):\n                    try:\n                        # Try ISO format first\n
  \                       return datetime.datetime.fromisoformat(date.replace(\"Z\",
  \"+00:00\"))\n                    except ValueError:\n                        try:\n
  \                           # Try parsing date with time first\n                            return
  datetime.datetime.strptime(date, \"%Y-%m-%d %H:%M\")\n                        except
  ValueError:\n                            try:\n                                #
  Fallback to date only\n                                return datetime.datetime.strptime(date,
  \"%Y-%m-%d\")\n                            except ValueError:\n                                #
  Try dateparser as last resort\n                                import dateparser\n\n
  \                               parsed = dateparser.parse(date)\n                                if
  parsed:\n                                    return parsed\n                if isinstance(date,
  datetime.date):\n                    return datetime.datetime.combine(date, datetime.time.min)\n\n
  \           # If we still don't have a date, use now\n            return datetime.datetime.now()\n
  \   ```\n!!! method\n    <h2 id=\"__init__\" class=\"admonition-title\" style=\"margin:
  0; padding: .5rem 1rem;\">__init__ <em class=\"small\">method</em></h2>\n\n    include:
  post attributes to include by default in Post\n    model serialization.\n    repr_include:
  post attributes to include by default in Post\n    repr.  If `repr_include` is None,
  it will default to\n    `include`, but it is likely that you want less in the repr\n
  \   than serialized output.\n\n    example:\n\n    ``` toml title='markata.toml'\n
  \   [markata.post_model]\n    include = ['date', 'description', 'published',\n        'slug',
  'title', 'content', 'html']\n    repr_include = ['date', 'description', 'published',
  'slug', 'title']\n    ```\n\n???+ source \"__init__ <em class='small'>source</em>\"\n
  \   ```python\n    def __init__(self, **data) -> None:\n            \"\"\"\n\n            include:
  post attributes to include by default in Post\n            model serialization.\n
  \           repr_include: post attributes to include by default in Post\n            repr.
  \ If `repr_include` is None, it will default to\n            `include`, but it is
  likely that you want less in the repr\n            than serialized output.\n\n            example:\n\n
  \           ``` toml title='markata.toml'\n            [markata.post_model]\n            include
  = ['date', 'description', 'published',\n                'slug', 'title', 'content',
  'html']\n            repr_include = ['date', 'description', 'published', 'slug',
  'title']\n            ```\n            \"\"\"\n            super().__init__(**data)\n
  \   ```"
date: 2025-05-05
description: "The   plugin defines the core Post model used throughout Markata. It
  provides robust validation, serialization, and configuration options for all post\u2026"
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>post_model.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"The   plugin defines the core Post model
    used throughout Markata. It provides robust validation, serialization, and configuration
    options for all post\u2026\" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n        <meta property=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    <head>\n        <script src=\"https://unpkg.com/@tailwindcss/browser@4\"></script>\n<title>post_model.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"The   plugin defines the core Post model
    used throughout Markata. It provides robust validation, serialization, and configuration
    options for all post\u2026\" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        post_model.py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <hr />\n<p>The <code>markata.plugins.post_model</code>
    plugin defines the core Post model used throughout\nMarkata. It provides robust
    validation, serialization, and configuration options for\nall post attributes.</p>\n<h2
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
    class=\"p\">[</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;markata.plugins.post_model&quot;</span><span
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
    class=\"p\">[</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;markata.plugins.post_model&quot;</span><span
    class=\"p\">,</span>\n<span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>Note:
    Disabling this plugin will break most of Markata's functionality as the Post\nmodel
    is fundamental to the system.</p>\n<h2 id=\"configuration\">Configuration <a class=\"header-anchor\"
    href=\"#configuration\"><svg class=\"heading-permalink\" aria-hidden=\"true\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Configure post model
    behavior in your <code>markata.toml</code>:</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata.post_model]</span>\n<span
    class=\"c1\"># Attributes to include when serializing posts</span>\n<span class=\"n\">include</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"p\">[</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;date&quot;</span><span
    class=\"p\">,</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;description&quot;</span><span
    class=\"p\">,</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;published&quot;</span><span
    class=\"p\">,</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">,</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;title&quot;</span><span
    class=\"p\">,</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;content&quot;</span><span
    class=\"p\">,</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;html&quot;</span>\n<span
    class=\"p\">]</span>\n\n<span class=\"c1\"># Attributes to show in post representations</span>\n<span
    class=\"n\">repr_include</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"p\">[</span>\n<span class=\"w\">    </span><span
    class=\"s2\">&quot;date&quot;</span><span class=\"p\">,</span>\n<span class=\"w\">
    \   </span><span class=\"s2\">&quot;description&quot;</span><span class=\"p\">,</span>\n<span
    class=\"w\">    </span><span class=\"s2\">&quot;published&quot;</span><span class=\"p\">,</span>\n<span
    class=\"w\">    </span><span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">,</span>\n<span
    class=\"w\">    </span><span class=\"s2\">&quot;title&quot;</span>\n<span class=\"p\">]</span>\n\n<span
    class=\"c1\"># Attributes to include when exporting</span>\n<span class=\"n\">export_include</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"p\">[</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;date&quot;</span><span
    class=\"p\">,</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;description&quot;</span><span
    class=\"p\">,</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;published&quot;</span><span
    class=\"p\">,</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">,</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;title&quot;</span>\n<span
    class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<h2 id=\"functionality\">Functionality
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<h2 id=\"post-model\">Post
    Model <a class=\"header-anchor\" href=\"#post-model\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Core attributes:</p>\n<ul>\n<li><code>path</code>:
    Path to source file</li>\n<li><code>slug</code>: URL-friendly identifier</li>\n<li><code>href</code>:
    Full URL path</li>\n<li><code>published</code>: Publication status</li>\n<li><code>description</code>:
    Post summary</li>\n<li><code>content</code>: Raw markdown content</li>\n<li><code>html</code>:
    Rendered HTML content</li>\n<li><code>tags</code>: List of post tags</li>\n<li><code>date</code>:
    Publication date</li>\n<li><code>title</code>: Post title</li>\n</ul>\n<h2 id=\"validation\">Validation
    <a class=\"header-anchor\" href=\"#validation\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>The model provides:</p>\n<ul>\n<li>Type
    checking and coercion</li>\n<li>Required field validation</li>\n<li>Custom field
    validators</li>\n<li>Default values</li>\n<li>Rich error messages</li>\n</ul>\n<h2
    id=\"serialization\">Serialization <a class=\"header-anchor\" href=\"#serialization\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Supports multiple output
    formats:</p>\n<ul>\n<li>Full serialization (all fields)</li>\n<li>Representation
    (subset for display)</li>\n<li>Export (subset for external use)</li>\n<li>JSON/YAML
    compatible</li>\n</ul>\n<h2 id=\"performance\">Performance <a class=\"header-anchor\"
    href=\"#performance\"><svg class=\"heading-permalink\" aria-hidden=\"true\" fill=\"currentColor\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Uses optimized Pydantic
    config:</p>\n<ul>\n<li>Disabled assignment validation</li>\n<li>Arbitrary types
    allowed</li>\n<li>Extra fields allowed</li>\n<li>String whitespace stripping</li>\n<li>Default
    value validation</li>\n<li>Number to string coercion</li>\n<li>Alias population</li>\n</ul>\n<h2
    id=\"dependencies\">Dependencies <a class=\"header-anchor\" href=\"#dependencies\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>This plugin depends
    on:</p>\n<ul>\n<li>pydantic for model definition</li>\n<li>rich for console output</li>\n<li>yaml
    for YAML handling</li>\n</ul>\n<hr />\n<div class=\"admonition class\">\n<p class=\"admonition-title\">Class</p>\n<h2
    id=\"PostModelConfig\" class=\"admonition-title\" style=\"margin: 0; padding:
    .5rem 1rem;\">PostModelConfig <em class=\"small\">class</em></h2>\n<p>Configuration
    for the Post model</p>\n</div>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">PostModelConfig <em class='small'>source</em></p>\n<pre
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
    class=\"w\"> </span><span class=\"nc\">PostModelConfig</span><span class=\"p\">(</span><span
    class=\"n\">pydantic</span><span class=\"o\">.</span><span class=\"n\">BaseModel</span><span
    class=\"p\">):</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Configuration
    for the Post model&quot;&quot;&quot;</span>\n\n    <span class=\"k\">def</span><span
    class=\"w\"> </span><span class=\"fm\">__init__</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"n\">data</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span class=\"w\">        </span><span
    class=\"sd\">&quot;&quot;&quot;</span>\n\n<span class=\"sd\">        include:
    post attributes to include by default in Post</span>\n<span class=\"sd\">        model
    serialization.</span>\n<span class=\"sd\">        repr_include: post attributes
    to include by default in Post</span>\n<span class=\"sd\">        repr.  If `repr_include`
    is None, it will default to</span>\n<span class=\"sd\">        `include`, but
    it is likely that you want less in the repr</span>\n<span class=\"sd\">        than
    serialized output.</span>\n\n<span class=\"sd\">        example:</span>\n\n<span
    class=\"sd\">        ``` toml title=&#39;markata.toml&#39;</span>\n<span class=\"sd\">
    \       [markata.post_model]</span>\n<span class=\"sd\">        include = [&#39;date&#39;,
    &#39;description&#39;, &#39;published&#39;,</span>\n<span class=\"sd\">            &#39;slug&#39;,
    &#39;title&#39;, &#39;content&#39;, &#39;html&#39;]</span>\n<span class=\"sd\">
    \       repr_include = [&#39;date&#39;, &#39;description&#39;, &#39;published&#39;,
    &#39;slug&#39;, &#39;title&#39;]</span>\n<span class=\"sd\">        ```</span>\n<span
    class=\"sd\">        &quot;&quot;&quot;</span>\n        <span class=\"nb\">super</span><span
    class=\"p\">()</span><span class=\"o\">.</span><span class=\"fm\">__init__</span><span
    class=\"p\">(</span><span class=\"o\">**</span><span class=\"n\">data</span><span
    class=\"p\">)</span>\n\n    <span class=\"n\">default_date</span><span class=\"p\">:</span>
    <span class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">date</span>
    <span class=\"o\">=</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">date</span><span class=\"o\">.</span><span class=\"n\">today</span><span
    class=\"p\">()</span>\n    <span class=\"n\">include</span><span class=\"p\">:</span>
    <span class=\"n\">List</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n        <span
    class=\"s2\">&quot;date&quot;</span><span class=\"p\">,</span>\n        <span
    class=\"s2\">&quot;description&quot;</span><span class=\"p\">,</span>\n        <span
    class=\"s2\">&quot;published&quot;</span><span class=\"p\">,</span>\n        <span
    class=\"s2\">&quot;slug&quot;</span><span class=\"p\">,</span>\n        <span
    class=\"s2\">&quot;title&quot;</span><span class=\"p\">,</span>\n        <span
    class=\"s2\">&quot;content&quot;</span><span class=\"p\">,</span>\n        <span
    class=\"s2\">&quot;html&quot;</span><span class=\"p\">,</span>\n    <span class=\"p\">]</span>\n
    \   <span class=\"n\">repr_include</span><span class=\"p\">:</span> <span class=\"n\">Optional</span><span
    class=\"p\">[</span><span class=\"n\">List</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]]</span> <span class=\"o\">=</span>
    <span class=\"p\">[</span>\n        <span class=\"s2\">&quot;date&quot;</span><span
    class=\"p\">,</span>\n        <span class=\"s2\">&quot;description&quot;</span><span
    class=\"p\">,</span>\n        <span class=\"s2\">&quot;published&quot;</span><span
    class=\"p\">,</span>\n        <span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">,</span>\n        <span class=\"s2\">&quot;title&quot;</span><span
    class=\"p\">,</span>\n    <span class=\"p\">]</span>\n    <span class=\"n\">export_include</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"n\">List</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
    class=\"p\">]]</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
    \       <span class=\"s2\">&quot;date&quot;</span><span class=\"p\">,</span>\n
    \       <span class=\"s2\">&quot;description&quot;</span><span class=\"p\">,</span>\n
    \       <span class=\"s2\">&quot;published&quot;</span><span class=\"p\">,</span>\n
    \       <span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">,</span>\n
    \       <span class=\"s2\">&quot;title&quot;</span><span class=\"p\">,</span>\n
    \   <span class=\"p\">]</span>\n\n    <span class=\"n\">model_config</span> <span
    class=\"o\">=</span> <span class=\"n\">ConfigDict</span><span class=\"p\">(</span>\n
    \       <span class=\"n\">validate_assignment</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>  <span class=\"c1\"># Config
    model</span>\n        <span class=\"n\">arbitrary_types_allowed</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n        <span class=\"n\">extra</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;allow&quot;</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">str_strip_whitespace</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n        <span class=\"n\">validate_default</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">coerce_numbers_to_str</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n        <span class=\"n\">populate_by_name</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \   <span class=\"p\">)</span>\n\n    <span class=\"nd\">@field_validator</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;repr_include&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">mode</span><span class=\"o\">=</span><span class=\"s2\">&quot;before&quot;</span><span
    class=\"p\">)</span>\n    <span class=\"nd\">@classmethod</span>\n    <span class=\"k\">def</span><span
    class=\"w\"> </span><span class=\"nf\">repr_include_validator</span><span class=\"p\">(</span><span
    class=\"bp\">cls</span><span class=\"p\">,</span> <span class=\"n\">v</span><span
    class=\"p\">,</span> <span class=\"n\">info</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"n\">List</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
    class=\"p\">]]:</span>\n        <span class=\"k\">if</span> <span class=\"n\">v</span><span
    class=\"p\">:</span>\n            <span class=\"k\">return</span> <span class=\"n\">v</span>\n
    \       <span class=\"k\">return</span> <span class=\"n\">info</span><span class=\"o\">.</span><span
    class=\"n\">data</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;include&quot;</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"metadata\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">metadata
    <em class=\"small\">method</em></h2>\n<p>for backwards compatability</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">metadata
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
    class=\"w\"> </span><span class=\"nf\">metadata</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Dict</span><span
    class=\"p\">:</span>\n        <span class=\"s2\">&quot;for backwards compatability&quot;</span>\n
    \       <span class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"vm\">__dict__</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div class=\"admonition
    method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"to_dict\" class=\"admonition-title\"
    style=\"margin: 0; padding: .5rem 1rem;\">to_dict <em class=\"small\">method</em></h2>\n<p>for
    backwards compatability</p>\n</div>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">to_dict <em class='small'>source</em></p>\n<pre
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
    class=\"w\"> </span><span class=\"nf\">to_dict</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Dict</span><span
    class=\"p\">:</span>\n        <span class=\"s2\">&quot;for backwards compatability&quot;</span>\n
    \       <span class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"vm\">__dict__</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div class=\"admonition
    method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"__getitem__\" class=\"admonition-title\"
    style=\"margin: 0; padding: .5rem 1rem;\">__getitem__ <em class=\"small\">method</em></h2>\n<p>for
    backwards compatability</p>\n</div>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\"><strong>getitem</strong> <em
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
    class=\"w\"> </span><span class=\"fm\">__getitem__</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">item</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Any</span><span class=\"p\">:</span>\n        <span class=\"s2\">&quot;for
    backwards compatability&quot;</span>\n        <span class=\"k\">return</span>
    <span class=\"nb\">getattr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">item</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"__setitem__\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">__setitem__
    <em class=\"small\">method</em></h2>\n<p>for backwards compatability</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\"><strong>setitem</strong>
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
    class=\"w\"> </span><span class=\"fm\">__setitem__</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">key</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span><span class=\"p\">,</span> <span class=\"n\">item</span><span
    class=\"p\">:</span> <span class=\"n\">Any</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \       <span class=\"s2\">&quot;for backwards compatability&quot;</span>\n        <span
    class=\"nb\">setattr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">key</span><span class=\"p\">,</span> <span
    class=\"n\">item</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"get\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">get <em class=\"small\">method</em></h2>\n<p>for
    backwards compatability</p>\n</div>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">get <em class='small'>source</em></p>\n<pre
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
    class=\"w\"> </span><span class=\"nf\">get</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">item</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span><span class=\"p\">,</span> <span class=\"n\">default</span><span
    class=\"p\">:</span> <span class=\"n\">Any</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span class=\"p\">:</span>\n
    \       <span class=\"s2\">&quot;for backwards compatability&quot;</span>\n        <span
    class=\"k\">return</span> <span class=\"nb\">getattr</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"n\">item</span><span
    class=\"p\">,</span> <span class=\"n\">default</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"keys\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">keys <em
    class=\"small\">method</em></h2>\n<p>for backwards compatability</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">keys
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
    class=\"w\"> </span><span class=\"nf\">keys</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">List</span><span
    class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">]:</span>\n
    \       <span class=\"s2\">&quot;for backwards compatability&quot;</span>\n        <span
    class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"vm\">__dict__</span><span class=\"o\">.</span><span class=\"n\">keys</span><span
    class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div class=\"admonition
    method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"yaml\" class=\"admonition-title\"
    style=\"margin: 0; padding: .5rem 1rem;\">yaml <em class=\"small\">method</em></h2>\n<p>dump
    model to yaml</p>\n</div>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">yaml <em class='small'>source</em></p>\n<pre class='wrapper'>\n\n<div
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
    class=\"w\"> </span><span class=\"nf\">yaml</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span
    class=\"p\">:</span>\n<span class=\"w\">        </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">        dump model to yaml</span>\n<span class=\"sd\">        &quot;&quot;&quot;</span>\n
    \       <span class=\"kn\">import</span><span class=\"w\"> </span><span class=\"nn\">yaml</span>\n\n
    \       <span class=\"k\">return</span> <span class=\"n\">yaml</span><span class=\"o\">.</span><span
    class=\"n\">dump</span><span class=\"p\">(</span>\n            <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">dict</span><span class=\"p\">(</span>\n
    \               <span class=\"n\">include</span><span class=\"o\">=</span><span
    class=\"p\">{</span><span class=\"n\">i</span><span class=\"p\">:</span> <span
    class=\"kc\">True</span> <span class=\"k\">for</span> <span class=\"n\">i</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">post_model</span><span class=\"o\">.</span><span
    class=\"n\">include</span><span class=\"p\">}</span>\n            <span class=\"p\">),</span>\n
    \           <span class=\"n\">Dumper</span><span class=\"o\">=</span><span class=\"n\">yaml</span><span
    class=\"o\">.</span><span class=\"n\">CDumper</span><span class=\"p\">,</span>\n
    \       <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div class=\"admonition
    method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"markdown\" class=\"admonition-title\"
    style=\"margin: 0; padding: .5rem 1rem;\">markdown <em class=\"small\">method</em></h2>\n<p>dump
    model to markdown</p>\n</div>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">markdown <em class='small'>source</em></p>\n<pre class='wrapper'>\n\n<div
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
    class=\"w\"> </span><span class=\"nf\">markdown</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span
    class=\"p\">:</span>\n<span class=\"w\">        </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">        dump model to markdown</span>\n<span class=\"sd\">        &quot;&quot;&quot;</span>\n\n
    \       <span class=\"kn\">import</span><span class=\"w\"> </span><span class=\"nn\">yaml</span>\n\n
    \       <span class=\"n\">frontmatter</span> <span class=\"o\">=</span> <span
    class=\"n\">yaml</span><span class=\"o\">.</span><span class=\"n\">dump</span><span
    class=\"p\">(</span>\n            <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">dict</span><span class=\"p\">(</span>\n                <span class=\"n\">include</span><span
    class=\"o\">=</span><span class=\"p\">{</span>\n                    <span class=\"n\">i</span><span
    class=\"p\">:</span> <span class=\"kc\">True</span>\n                    <span
    class=\"k\">for</span> <span class=\"n\">i</span> <span class=\"ow\">in</span>
    <span class=\"p\">[</span>\n                        <span class=\"n\">_i</span>\n
    \                       <span class=\"k\">for</span> <span class=\"n\">_i</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">post_model</span><span class=\"o\">.</span><span
    class=\"n\">include</span>\n                        <span class=\"k\">if</span>
    <span class=\"n\">_i</span> <span class=\"o\">!=</span> <span class=\"s2\">&quot;content&quot;</span>\n
    \                   <span class=\"p\">]</span>\n                <span class=\"p\">}</span>\n
    \           <span class=\"p\">),</span>\n            <span class=\"n\">Dumper</span><span
    class=\"o\">=</span><span class=\"n\">yaml</span><span class=\"o\">.</span><span
    class=\"n\">CDumper</span><span class=\"p\">,</span>\n        <span class=\"p\">)</span>\n
    \       <span class=\"n\">post</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;---</span><span
    class=\"se\">\\n</span><span class=\"s2\">&quot;</span>\n        <span class=\"n\">post</span>
    <span class=\"o\">+=</span> <span class=\"n\">frontmatter</span>\n        <span
    class=\"n\">post</span> <span class=\"o\">+=</span> <span class=\"s2\">&quot;---</span><span
    class=\"se\">\\n\\n</span><span class=\"s2\">&quot;</span>\n\n        <span class=\"k\">if</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">content</span><span
    class=\"p\">:</span>\n            <span class=\"n\">post</span> <span class=\"o\">+=</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">content</span>\n
    \       <span class=\"k\">return</span> <span class=\"n\">post</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"dumps\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">dumps <em
    class=\"small\">method</em></h2>\n<p>dumps raw article back out</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">dumps
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
    class=\"w\"> </span><span class=\"nf\">dumps</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">):</span>\n<span class=\"w\">        </span><span
    class=\"sd\">&quot;&quot;&quot;</span>\n<span class=\"sd\">        dumps raw article
    back out</span>\n<span class=\"sd\">        &quot;&quot;&quot;</span>\n        <span
    class=\"k\">return</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;---</span><span
    class=\"se\">\\n</span><span class=\"si\">{</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">yaml</span><span class=\"p\">()</span><span
    class=\"si\">}</span><span class=\"se\">\\n\\n</span><span class=\"s2\">---</span><span
    class=\"se\">\\n\\n</span><span class=\"si\">{</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">content</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div class=\"admonition
    method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"parse_date_time\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">parse_date_time
    <em class=\"small\">method</em></h2>\n<p>Single validator to handle all date_time
    parsing cases</p>\n</div>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">parse_date_time <em class='small'>source</em></p>\n<pre
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
    class=\"w\"> </span><span class=\"nf\">parse_date_time</span><span class=\"p\">(</span><span
    class=\"bp\">cls</span><span class=\"p\">,</span> <span class=\"n\">v</span><span
    class=\"p\">,</span> <span class=\"n\">info</span><span class=\"p\">):</span>\n<span
    class=\"w\">        </span><span class=\"sd\">&quot;&quot;&quot;Single validator
    to handle all date_time parsing cases&quot;&quot;&quot;</span>\n        <span
    class=\"c1\"># If we have an explicit date_time value</span>\n        <span class=\"k\">if</span>
    <span class=\"n\">v</span> <span class=\"ow\">is</span> <span class=\"ow\">not</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">v</span><span
    class=\"p\">,</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">datetime</span><span class=\"p\">):</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">v</span>\n            <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">v</span><span
    class=\"p\">,</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">date</span><span class=\"p\">):</span>\n                <span class=\"k\">return</span>
    <span class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">combine</span><span class=\"p\">(</span><span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">time</span><span class=\"o\">.</span><span
    class=\"n\">min</span><span class=\"p\">)</span>\n            <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">v</span><span
    class=\"p\">,</span> <span class=\"nb\">str</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">try</span><span class=\"p\">:</span>\n                    <span
    class=\"c1\"># Try ISO format first</span>\n                    <span class=\"k\">return</span>
    <span class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">fromisoformat</span><span class=\"p\">(</span><span
    class=\"n\">v</span><span class=\"o\">.</span><span class=\"n\">replace</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;Z&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;+00:00&quot;</span><span class=\"p\">))</span>\n                <span
    class=\"k\">except</span> <span class=\"ne\">ValueError</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">try</span><span class=\"p\">:</span>\n                        <span
    class=\"k\">return</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">strptime</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;%Y-%m-</span><span class=\"si\">%d</span><span class=\"s2\">
    %H:%M&quot;</span><span class=\"p\">)</span>\n                    <span class=\"k\">except</span>
    <span class=\"ne\">ValueError</span><span class=\"p\">:</span>\n                        <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                            <span
    class=\"k\">return</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">strptime</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;%Y-%m-</span><span class=\"si\">%d</span><span class=\"s2\">&quot;</span><span
    class=\"p\">)</span>\n                        <span class=\"k\">except</span>
    <span class=\"ne\">ValueError</span><span class=\"p\">:</span>\n                            <span
    class=\"c1\"># Try dateparser as last resort for explicit date_time</span>\n                            <span
    class=\"kn\">import</span><span class=\"w\"> </span><span class=\"nn\">dateparser</span>\n\n
    \                           <span class=\"n\">parsed</span> <span class=\"o\">=</span>
    <span class=\"n\">dateparser</span><span class=\"o\">.</span><span class=\"n\">parse</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">)</span>\n                            <span
    class=\"k\">if</span> <span class=\"n\">parsed</span><span class=\"p\">:</span>\n
    \                               <span class=\"k\">return</span> <span class=\"n\">parsed</span>\n
    \                           <span class=\"k\">return</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">now</span><span class=\"p\">()</span>\n\n        <span class=\"c1\">#
    Get the raw date string directly from raw_date field</span>\n        <span class=\"n\">raw_date</span>
    <span class=\"o\">=</span> <span class=\"n\">info</span><span class=\"o\">.</span><span
    class=\"n\">data</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;raw_date&quot;</span><span class=\"p\">)</span>\n
    \       <span class=\"k\">if</span> <span class=\"n\">raw_date</span> <span class=\"ow\">and</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">raw_date</span><span
    class=\"p\">,</span> <span class=\"nb\">str</span><span class=\"p\">):</span>\n
    \           <span class=\"k\">try</span><span class=\"p\">:</span>\n                <span
    class=\"c1\"># Try ISO format first</span>\n                <span class=\"k\">return</span>
    <span class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">fromisoformat</span><span class=\"p\">(</span><span
    class=\"n\">raw_date</span><span class=\"o\">.</span><span class=\"n\">replace</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;Z&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;+00:00&quot;</span><span class=\"p\">))</span>\n            <span
    class=\"k\">except</span> <span class=\"ne\">ValueError</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">try</span><span class=\"p\">:</span>\n                    <span
    class=\"c1\"># Try parsing raw_date with time first</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">strptime</span><span
    class=\"p\">(</span><span class=\"n\">raw_date</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;%Y-%m-</span><span class=\"si\">%d</span><span class=\"s2\">
    %H:%M&quot;</span><span class=\"p\">)</span>\n                <span class=\"k\">except</span>
    <span class=\"ne\">ValueError</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                        <span
    class=\"c1\"># Fallback to date only</span>\n                        <span class=\"k\">return</span>
    <span class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">strptime</span><span class=\"p\">(</span><span
    class=\"n\">raw_date</span><span class=\"p\">,</span> <span class=\"s2\">&quot;%Y-%m-</span><span
    class=\"si\">%d</span><span class=\"s2\">&quot;</span><span class=\"p\">)</span>\n
    \                   <span class=\"k\">except</span> <span class=\"ne\">ValueError</span><span
    class=\"p\">:</span>\n                        <span class=\"c1\"># Try dateparser
    as last resort</span>\n                        <span class=\"kn\">import</span><span
    class=\"w\"> </span><span class=\"nn\">dateparser</span>\n\n                        <span
    class=\"n\">parsed</span> <span class=\"o\">=</span> <span class=\"n\">dateparser</span><span
    class=\"o\">.</span><span class=\"n\">parse</span><span class=\"p\">(</span><span
    class=\"n\">raw_date</span><span class=\"p\">)</span>\n                        <span
    class=\"k\">if</span> <span class=\"n\">parsed</span><span class=\"p\">:</span>\n
    \                           <span class=\"k\">return</span> <span class=\"n\">parsed</span>\n\n
    \       <span class=\"c1\"># If no raw_date, try to derive from date field</span>\n
    \       <span class=\"n\">date</span> <span class=\"o\">=</span> <span class=\"n\">info</span><span
    class=\"o\">.</span><span class=\"n\">data</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;date&quot;</span><span
    class=\"p\">)</span>\n        <span class=\"k\">if</span> <span class=\"n\">date</span><span
    class=\"p\">:</span>\n            <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">date</span><span class=\"p\">,</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"p\">):</span>\n                <span class=\"k\">return</span> <span class=\"n\">date</span>\n
    \           <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">date</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">):</span>\n                <span class=\"k\">try</span><span
    class=\"p\">:</span>\n                    <span class=\"c1\"># Try ISO format
    first</span>\n                    <span class=\"k\">return</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">fromisoformat</span><span class=\"p\">(</span><span class=\"n\">date</span><span
    class=\"o\">.</span><span class=\"n\">replace</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;Z&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;+00:00&quot;</span><span
    class=\"p\">))</span>\n                <span class=\"k\">except</span> <span class=\"ne\">ValueError</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                       <span class=\"c1\"># Try parsing date with time first</span>\n
    \                       <span class=\"k\">return</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">strptime</span><span class=\"p\">(</span><span class=\"n\">date</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;%Y-%m-</span><span class=\"si\">%d</span><span
    class=\"s2\"> %H:%M&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">except</span> <span class=\"ne\">ValueError</span><span class=\"p\">:</span>\n
    \                       <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                           <span class=\"c1\"># Fallback to date only</span>\n
    \                           <span class=\"k\">return</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">strptime</span><span class=\"p\">(</span><span class=\"n\">date</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;%Y-%m-</span><span class=\"si\">%d</span><span
    class=\"s2\">&quot;</span><span class=\"p\">)</span>\n                        <span
    class=\"k\">except</span> <span class=\"ne\">ValueError</span><span class=\"p\">:</span>\n
    \                           <span class=\"c1\"># Try dateparser as last resort</span>\n
    \                           <span class=\"kn\">import</span><span class=\"w\">
    </span><span class=\"nn\">dateparser</span>\n\n                            <span
    class=\"n\">parsed</span> <span class=\"o\">=</span> <span class=\"n\">dateparser</span><span
    class=\"o\">.</span><span class=\"n\">parse</span><span class=\"p\">(</span><span
    class=\"n\">date</span><span class=\"p\">)</span>\n                            <span
    class=\"k\">if</span> <span class=\"n\">parsed</span><span class=\"p\">:</span>\n
    \                               <span class=\"k\">return</span> <span class=\"n\">parsed</span>\n
    \           <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">date</span><span class=\"p\">,</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">date</span><span
    class=\"p\">):</span>\n                <span class=\"k\">return</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">combine</span><span class=\"p\">(</span><span class=\"n\">date</span><span
    class=\"p\">,</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">time</span><span class=\"o\">.</span><span class=\"n\">min</span><span
    class=\"p\">)</span>\n\n        <span class=\"c1\"># If we still don&#39;t have
    a date, use now</span>\n        <span class=\"k\">return</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">now</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"__init__\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">__init__
    <em class=\"small\">method</em></h2>\n<p>include: post attributes to include by
    default in Post\nmodel serialization.\nrepr_include: post attributes to include
    by default in Post\nrepr.  If <code>repr_include</code> is None, it will default
    to\n<code>include</code>, but it is likely that you want less in the repr\nthan
    serialized output.</p>\n<p>example:</p>\n<pre class='wrapper'>\n\n<div class='filepath'>\n<p>markata.toml</p>\n\n<div
    class='right'>\n\n<button class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
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
    class=\"highlight\"><pre><span></span><span class=\"k\">[markata.post_model]</span>\n<span
    class=\"n\">include</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"p\">[</span><span class=\"s1\">&#39;date&#39;</span><span
    class=\"p\">,</span><span class=\"w\"> </span><span class=\"s1\">&#39;description&#39;</span><span
    class=\"p\">,</span><span class=\"w\"> </span><span class=\"s1\">&#39;published&#39;</span><span
    class=\"p\">,</span>\n<span class=\"w\">    </span><span class=\"s1\">&#39;slug&#39;</span><span
    class=\"p\">,</span><span class=\"w\"> </span><span class=\"s1\">&#39;title&#39;</span><span
    class=\"p\">,</span><span class=\"w\"> </span><span class=\"s1\">&#39;content&#39;</span><span
    class=\"p\">,</span><span class=\"w\"> </span><span class=\"s1\">&#39;html&#39;</span><span
    class=\"p\">]</span>\n<span class=\"n\">repr_include</span><span class=\"w\">
    </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"p\">[</span><span
    class=\"s1\">&#39;date&#39;</span><span class=\"p\">,</span><span class=\"w\">
    </span><span class=\"s1\">&#39;description&#39;</span><span class=\"p\">,</span><span
    class=\"w\"> </span><span class=\"s1\">&#39;published&#39;</span><span class=\"p\">,</span><span
    class=\"w\"> </span><span class=\"s1\">&#39;slug&#39;</span><span class=\"p\">,</span><span
    class=\"w\"> </span><span class=\"s1\">&#39;title&#39;</span><span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\"><strong>init</strong>
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
    class=\"w\"> </span><span class=\"fm\">__init__</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"n\">data</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span class=\"w\">        </span><span
    class=\"sd\">&quot;&quot;&quot;</span>\n\n<span class=\"sd\">        include:
    post attributes to include by default in Post</span>\n<span class=\"sd\">        model
    serialization.</span>\n<span class=\"sd\">        repr_include: post attributes
    to include by default in Post</span>\n<span class=\"sd\">        repr.  If `repr_include`
    is None, it will default to</span>\n<span class=\"sd\">        `include`, but
    it is likely that you want less in the repr</span>\n<span class=\"sd\">        than
    serialized output.</span>\n\n<span class=\"sd\">        example:</span>\n\n<span
    class=\"sd\">        ``` toml title=&#39;markata.toml&#39;</span>\n<span class=\"sd\">
    \       [markata.post_model]</span>\n<span class=\"sd\">        include = [&#39;date&#39;,
    &#39;description&#39;, &#39;published&#39;,</span>\n<span class=\"sd\">            &#39;slug&#39;,
    &#39;title&#39;, &#39;content&#39;, &#39;html&#39;]</span>\n<span class=\"sd\">
    \       repr_include = [&#39;date&#39;, &#39;description&#39;, &#39;published&#39;,
    &#39;slug&#39;, &#39;title&#39;]</span>\n<span class=\"sd\">        ```</span>\n<span
    class=\"sd\">        &quot;&quot;&quot;</span>\n        <span class=\"nb\">super</span><span
    class=\"p\">()</span><span class=\"o\">.</span><span class=\"fm\">__init__</span><span
    class=\"p\">(</span><span class=\"o\">**</span><span class=\"n\">data</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n\n    </section>\n</article>
    \   </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>post_model.py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"The   plugin defines the core Post model used throughout
    Markata. It provides robust validation, serialization, and configuration options
    for all post\u2026\" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n        <meta property=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    <head>\n        <script src=\"https://unpkg.com/@tailwindcss/browser@4\"></script>\n<title>post_model.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"The   plugin defines the core Post model
    used throughout Markata. It provides robust validation, serialization, and configuration
    options for all post\u2026\" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
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
    class=\"title\">\n    <h1 id=\"title\">\n        post_model.py\n    </h1>\n</section></article>\n
    \    </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       post_model.py\n    </h1>\n</section>    <section class=\"body\">\n        <hr
    />\n<p>The <code>markata.plugins.post_model</code> plugin defines the core Post
    model used throughout\nMarkata. It provides robust validation, serialization,
    and configuration options for\nall post attributes.</p>\n<h2 id=\"installation\">Installation
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
    class=\"p\">[</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;markata.plugins.post_model&quot;</span><span
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
    class=\"p\">[</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;markata.plugins.post_model&quot;</span><span
    class=\"p\">,</span>\n<span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>Note:
    Disabling this plugin will break most of Markata's functionality as the Post\nmodel
    is fundamental to the system.</p>\n<h2 id=\"configuration\">Configuration <a class=\"header-anchor\"
    href=\"#configuration\"><svg class=\"heading-permalink\" aria-hidden=\"true\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Configure post model
    behavior in your <code>markata.toml</code>:</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata.post_model]</span>\n<span
    class=\"c1\"># Attributes to include when serializing posts</span>\n<span class=\"n\">include</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"p\">[</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;date&quot;</span><span
    class=\"p\">,</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;description&quot;</span><span
    class=\"p\">,</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;published&quot;</span><span
    class=\"p\">,</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">,</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;title&quot;</span><span
    class=\"p\">,</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;content&quot;</span><span
    class=\"p\">,</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;html&quot;</span>\n<span
    class=\"p\">]</span>\n\n<span class=\"c1\"># Attributes to show in post representations</span>\n<span
    class=\"n\">repr_include</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"p\">[</span>\n<span class=\"w\">    </span><span
    class=\"s2\">&quot;date&quot;</span><span class=\"p\">,</span>\n<span class=\"w\">
    \   </span><span class=\"s2\">&quot;description&quot;</span><span class=\"p\">,</span>\n<span
    class=\"w\">    </span><span class=\"s2\">&quot;published&quot;</span><span class=\"p\">,</span>\n<span
    class=\"w\">    </span><span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">,</span>\n<span
    class=\"w\">    </span><span class=\"s2\">&quot;title&quot;</span>\n<span class=\"p\">]</span>\n\n<span
    class=\"c1\"># Attributes to include when exporting</span>\n<span class=\"n\">export_include</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"p\">[</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;date&quot;</span><span
    class=\"p\">,</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;description&quot;</span><span
    class=\"p\">,</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;published&quot;</span><span
    class=\"p\">,</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">,</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;title&quot;</span>\n<span
    class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<h2 id=\"functionality\">Functionality
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<h2 id=\"post-model\">Post
    Model <a class=\"header-anchor\" href=\"#post-model\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Core attributes:</p>\n<ul>\n<li><code>path</code>:
    Path to source file</li>\n<li><code>slug</code>: URL-friendly identifier</li>\n<li><code>href</code>:
    Full URL path</li>\n<li><code>published</code>: Publication status</li>\n<li><code>description</code>:
    Post summary</li>\n<li><code>content</code>: Raw markdown content</li>\n<li><code>html</code>:
    Rendered HTML content</li>\n<li><code>tags</code>: List of post tags</li>\n<li><code>date</code>:
    Publication date</li>\n<li><code>title</code>: Post title</li>\n</ul>\n<h2 id=\"validation\">Validation
    <a class=\"header-anchor\" href=\"#validation\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>The model provides:</p>\n<ul>\n<li>Type
    checking and coercion</li>\n<li>Required field validation</li>\n<li>Custom field
    validators</li>\n<li>Default values</li>\n<li>Rich error messages</li>\n</ul>\n<h2
    id=\"serialization\">Serialization <a class=\"header-anchor\" href=\"#serialization\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Supports multiple output
    formats:</p>\n<ul>\n<li>Full serialization (all fields)</li>\n<li>Representation
    (subset for display)</li>\n<li>Export (subset for external use)</li>\n<li>JSON/YAML
    compatible</li>\n</ul>\n<h2 id=\"performance\">Performance <a class=\"header-anchor\"
    href=\"#performance\"><svg class=\"heading-permalink\" aria-hidden=\"true\" fill=\"currentColor\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Uses optimized Pydantic
    config:</p>\n<ul>\n<li>Disabled assignment validation</li>\n<li>Arbitrary types
    allowed</li>\n<li>Extra fields allowed</li>\n<li>String whitespace stripping</li>\n<li>Default
    value validation</li>\n<li>Number to string coercion</li>\n<li>Alias population</li>\n</ul>\n<h2
    id=\"dependencies\">Dependencies <a class=\"header-anchor\" href=\"#dependencies\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>This plugin depends
    on:</p>\n<ul>\n<li>pydantic for model definition</li>\n<li>rich for console output</li>\n<li>yaml
    for YAML handling</li>\n</ul>\n<hr />\n<div class=\"admonition class\">\n<p class=\"admonition-title\">Class</p>\n<h2
    id=\"PostModelConfig\" class=\"admonition-title\" style=\"margin: 0; padding:
    .5rem 1rem;\">PostModelConfig <em class=\"small\">class</em></h2>\n<p>Configuration
    for the Post model</p>\n</div>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">PostModelConfig <em class='small'>source</em></p>\n<pre
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
    class=\"w\"> </span><span class=\"nc\">PostModelConfig</span><span class=\"p\">(</span><span
    class=\"n\">pydantic</span><span class=\"o\">.</span><span class=\"n\">BaseModel</span><span
    class=\"p\">):</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Configuration
    for the Post model&quot;&quot;&quot;</span>\n\n    <span class=\"k\">def</span><span
    class=\"w\"> </span><span class=\"fm\">__init__</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"n\">data</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span class=\"w\">        </span><span
    class=\"sd\">&quot;&quot;&quot;</span>\n\n<span class=\"sd\">        include:
    post attributes to include by default in Post</span>\n<span class=\"sd\">        model
    serialization.</span>\n<span class=\"sd\">        repr_include: post attributes
    to include by default in Post</span>\n<span class=\"sd\">        repr.  If `repr_include`
    is None, it will default to</span>\n<span class=\"sd\">        `include`, but
    it is likely that you want less in the repr</span>\n<span class=\"sd\">        than
    serialized output.</span>\n\n<span class=\"sd\">        example:</span>\n\n<span
    class=\"sd\">        ``` toml title=&#39;markata.toml&#39;</span>\n<span class=\"sd\">
    \       [markata.post_model]</span>\n<span class=\"sd\">        include = [&#39;date&#39;,
    &#39;description&#39;, &#39;published&#39;,</span>\n<span class=\"sd\">            &#39;slug&#39;,
    &#39;title&#39;, &#39;content&#39;, &#39;html&#39;]</span>\n<span class=\"sd\">
    \       repr_include = [&#39;date&#39;, &#39;description&#39;, &#39;published&#39;,
    &#39;slug&#39;, &#39;title&#39;]</span>\n<span class=\"sd\">        ```</span>\n<span
    class=\"sd\">        &quot;&quot;&quot;</span>\n        <span class=\"nb\">super</span><span
    class=\"p\">()</span><span class=\"o\">.</span><span class=\"fm\">__init__</span><span
    class=\"p\">(</span><span class=\"o\">**</span><span class=\"n\">data</span><span
    class=\"p\">)</span>\n\n    <span class=\"n\">default_date</span><span class=\"p\">:</span>
    <span class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">date</span>
    <span class=\"o\">=</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">date</span><span class=\"o\">.</span><span class=\"n\">today</span><span
    class=\"p\">()</span>\n    <span class=\"n\">include</span><span class=\"p\">:</span>
    <span class=\"n\">List</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n        <span
    class=\"s2\">&quot;date&quot;</span><span class=\"p\">,</span>\n        <span
    class=\"s2\">&quot;description&quot;</span><span class=\"p\">,</span>\n        <span
    class=\"s2\">&quot;published&quot;</span><span class=\"p\">,</span>\n        <span
    class=\"s2\">&quot;slug&quot;</span><span class=\"p\">,</span>\n        <span
    class=\"s2\">&quot;title&quot;</span><span class=\"p\">,</span>\n        <span
    class=\"s2\">&quot;content&quot;</span><span class=\"p\">,</span>\n        <span
    class=\"s2\">&quot;html&quot;</span><span class=\"p\">,</span>\n    <span class=\"p\">]</span>\n
    \   <span class=\"n\">repr_include</span><span class=\"p\">:</span> <span class=\"n\">Optional</span><span
    class=\"p\">[</span><span class=\"n\">List</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]]</span> <span class=\"o\">=</span>
    <span class=\"p\">[</span>\n        <span class=\"s2\">&quot;date&quot;</span><span
    class=\"p\">,</span>\n        <span class=\"s2\">&quot;description&quot;</span><span
    class=\"p\">,</span>\n        <span class=\"s2\">&quot;published&quot;</span><span
    class=\"p\">,</span>\n        <span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">,</span>\n        <span class=\"s2\">&quot;title&quot;</span><span
    class=\"p\">,</span>\n    <span class=\"p\">]</span>\n    <span class=\"n\">export_include</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"n\">List</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
    class=\"p\">]]</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
    \       <span class=\"s2\">&quot;date&quot;</span><span class=\"p\">,</span>\n
    \       <span class=\"s2\">&quot;description&quot;</span><span class=\"p\">,</span>\n
    \       <span class=\"s2\">&quot;published&quot;</span><span class=\"p\">,</span>\n
    \       <span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">,</span>\n
    \       <span class=\"s2\">&quot;title&quot;</span><span class=\"p\">,</span>\n
    \   <span class=\"p\">]</span>\n\n    <span class=\"n\">model_config</span> <span
    class=\"o\">=</span> <span class=\"n\">ConfigDict</span><span class=\"p\">(</span>\n
    \       <span class=\"n\">validate_assignment</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>  <span class=\"c1\"># Config
    model</span>\n        <span class=\"n\">arbitrary_types_allowed</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n        <span class=\"n\">extra</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;allow&quot;</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">str_strip_whitespace</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n        <span class=\"n\">validate_default</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">coerce_numbers_to_str</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n        <span class=\"n\">populate_by_name</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \   <span class=\"p\">)</span>\n\n    <span class=\"nd\">@field_validator</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;repr_include&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">mode</span><span class=\"o\">=</span><span class=\"s2\">&quot;before&quot;</span><span
    class=\"p\">)</span>\n    <span class=\"nd\">@classmethod</span>\n    <span class=\"k\">def</span><span
    class=\"w\"> </span><span class=\"nf\">repr_include_validator</span><span class=\"p\">(</span><span
    class=\"bp\">cls</span><span class=\"p\">,</span> <span class=\"n\">v</span><span
    class=\"p\">,</span> <span class=\"n\">info</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"n\">List</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
    class=\"p\">]]:</span>\n        <span class=\"k\">if</span> <span class=\"n\">v</span><span
    class=\"p\">:</span>\n            <span class=\"k\">return</span> <span class=\"n\">v</span>\n
    \       <span class=\"k\">return</span> <span class=\"n\">info</span><span class=\"o\">.</span><span
    class=\"n\">data</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;include&quot;</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"metadata\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">metadata
    <em class=\"small\">method</em></h2>\n<p>for backwards compatability</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">metadata
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
    class=\"w\"> </span><span class=\"nf\">metadata</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Dict</span><span
    class=\"p\">:</span>\n        <span class=\"s2\">&quot;for backwards compatability&quot;</span>\n
    \       <span class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"vm\">__dict__</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div class=\"admonition
    method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"to_dict\" class=\"admonition-title\"
    style=\"margin: 0; padding: .5rem 1rem;\">to_dict <em class=\"small\">method</em></h2>\n<p>for
    backwards compatability</p>\n</div>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">to_dict <em class='small'>source</em></p>\n<pre
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
    class=\"w\"> </span><span class=\"nf\">to_dict</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Dict</span><span
    class=\"p\">:</span>\n        <span class=\"s2\">&quot;for backwards compatability&quot;</span>\n
    \       <span class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"vm\">__dict__</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div class=\"admonition
    method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"__getitem__\" class=\"admonition-title\"
    style=\"margin: 0; padding: .5rem 1rem;\">__getitem__ <em class=\"small\">method</em></h2>\n<p>for
    backwards compatability</p>\n</div>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\"><strong>getitem</strong> <em
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
    class=\"w\"> </span><span class=\"fm\">__getitem__</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">item</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Any</span><span class=\"p\">:</span>\n        <span class=\"s2\">&quot;for
    backwards compatability&quot;</span>\n        <span class=\"k\">return</span>
    <span class=\"nb\">getattr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">item</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"__setitem__\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">__setitem__
    <em class=\"small\">method</em></h2>\n<p>for backwards compatability</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\"><strong>setitem</strong>
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
    class=\"w\"> </span><span class=\"fm\">__setitem__</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">key</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span><span class=\"p\">,</span> <span class=\"n\">item</span><span
    class=\"p\">:</span> <span class=\"n\">Any</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \       <span class=\"s2\">&quot;for backwards compatability&quot;</span>\n        <span
    class=\"nb\">setattr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">key</span><span class=\"p\">,</span> <span
    class=\"n\">item</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"get\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">get <em class=\"small\">method</em></h2>\n<p>for
    backwards compatability</p>\n</div>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">get <em class='small'>source</em></p>\n<pre
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
    class=\"w\"> </span><span class=\"nf\">get</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">item</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span><span class=\"p\">,</span> <span class=\"n\">default</span><span
    class=\"p\">:</span> <span class=\"n\">Any</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span class=\"p\">:</span>\n
    \       <span class=\"s2\">&quot;for backwards compatability&quot;</span>\n        <span
    class=\"k\">return</span> <span class=\"nb\">getattr</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"n\">item</span><span
    class=\"p\">,</span> <span class=\"n\">default</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"keys\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">keys <em
    class=\"small\">method</em></h2>\n<p>for backwards compatability</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">keys
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
    class=\"w\"> </span><span class=\"nf\">keys</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">List</span><span
    class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">]:</span>\n
    \       <span class=\"s2\">&quot;for backwards compatability&quot;</span>\n        <span
    class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"vm\">__dict__</span><span class=\"o\">.</span><span class=\"n\">keys</span><span
    class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div class=\"admonition
    method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"yaml\" class=\"admonition-title\"
    style=\"margin: 0; padding: .5rem 1rem;\">yaml <em class=\"small\">method</em></h2>\n<p>dump
    model to yaml</p>\n</div>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">yaml <em class='small'>source</em></p>\n<pre class='wrapper'>\n\n<div
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
    class=\"w\"> </span><span class=\"nf\">yaml</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span
    class=\"p\">:</span>\n<span class=\"w\">        </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">        dump model to yaml</span>\n<span class=\"sd\">        &quot;&quot;&quot;</span>\n
    \       <span class=\"kn\">import</span><span class=\"w\"> </span><span class=\"nn\">yaml</span>\n\n
    \       <span class=\"k\">return</span> <span class=\"n\">yaml</span><span class=\"o\">.</span><span
    class=\"n\">dump</span><span class=\"p\">(</span>\n            <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">dict</span><span class=\"p\">(</span>\n
    \               <span class=\"n\">include</span><span class=\"o\">=</span><span
    class=\"p\">{</span><span class=\"n\">i</span><span class=\"p\">:</span> <span
    class=\"kc\">True</span> <span class=\"k\">for</span> <span class=\"n\">i</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">post_model</span><span class=\"o\">.</span><span
    class=\"n\">include</span><span class=\"p\">}</span>\n            <span class=\"p\">),</span>\n
    \           <span class=\"n\">Dumper</span><span class=\"o\">=</span><span class=\"n\">yaml</span><span
    class=\"o\">.</span><span class=\"n\">CDumper</span><span class=\"p\">,</span>\n
    \       <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div class=\"admonition
    method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"markdown\" class=\"admonition-title\"
    style=\"margin: 0; padding: .5rem 1rem;\">markdown <em class=\"small\">method</em></h2>\n<p>dump
    model to markdown</p>\n</div>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">markdown <em class='small'>source</em></p>\n<pre class='wrapper'>\n\n<div
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
    class=\"w\"> </span><span class=\"nf\">markdown</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span
    class=\"p\">:</span>\n<span class=\"w\">        </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">        dump model to markdown</span>\n<span class=\"sd\">        &quot;&quot;&quot;</span>\n\n
    \       <span class=\"kn\">import</span><span class=\"w\"> </span><span class=\"nn\">yaml</span>\n\n
    \       <span class=\"n\">frontmatter</span> <span class=\"o\">=</span> <span
    class=\"n\">yaml</span><span class=\"o\">.</span><span class=\"n\">dump</span><span
    class=\"p\">(</span>\n            <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">dict</span><span class=\"p\">(</span>\n                <span class=\"n\">include</span><span
    class=\"o\">=</span><span class=\"p\">{</span>\n                    <span class=\"n\">i</span><span
    class=\"p\">:</span> <span class=\"kc\">True</span>\n                    <span
    class=\"k\">for</span> <span class=\"n\">i</span> <span class=\"ow\">in</span>
    <span class=\"p\">[</span>\n                        <span class=\"n\">_i</span>\n
    \                       <span class=\"k\">for</span> <span class=\"n\">_i</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">post_model</span><span class=\"o\">.</span><span
    class=\"n\">include</span>\n                        <span class=\"k\">if</span>
    <span class=\"n\">_i</span> <span class=\"o\">!=</span> <span class=\"s2\">&quot;content&quot;</span>\n
    \                   <span class=\"p\">]</span>\n                <span class=\"p\">}</span>\n
    \           <span class=\"p\">),</span>\n            <span class=\"n\">Dumper</span><span
    class=\"o\">=</span><span class=\"n\">yaml</span><span class=\"o\">.</span><span
    class=\"n\">CDumper</span><span class=\"p\">,</span>\n        <span class=\"p\">)</span>\n
    \       <span class=\"n\">post</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;---</span><span
    class=\"se\">\\n</span><span class=\"s2\">&quot;</span>\n        <span class=\"n\">post</span>
    <span class=\"o\">+=</span> <span class=\"n\">frontmatter</span>\n        <span
    class=\"n\">post</span> <span class=\"o\">+=</span> <span class=\"s2\">&quot;---</span><span
    class=\"se\">\\n\\n</span><span class=\"s2\">&quot;</span>\n\n        <span class=\"k\">if</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">content</span><span
    class=\"p\">:</span>\n            <span class=\"n\">post</span> <span class=\"o\">+=</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">content</span>\n
    \       <span class=\"k\">return</span> <span class=\"n\">post</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"dumps\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">dumps <em
    class=\"small\">method</em></h2>\n<p>dumps raw article back out</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">dumps
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
    class=\"w\"> </span><span class=\"nf\">dumps</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">):</span>\n<span class=\"w\">        </span><span
    class=\"sd\">&quot;&quot;&quot;</span>\n<span class=\"sd\">        dumps raw article
    back out</span>\n<span class=\"sd\">        &quot;&quot;&quot;</span>\n        <span
    class=\"k\">return</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;---</span><span
    class=\"se\">\\n</span><span class=\"si\">{</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">yaml</span><span class=\"p\">()</span><span
    class=\"si\">}</span><span class=\"se\">\\n\\n</span><span class=\"s2\">---</span><span
    class=\"se\">\\n\\n</span><span class=\"si\">{</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">content</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div class=\"admonition
    method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"parse_date_time\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">parse_date_time
    <em class=\"small\">method</em></h2>\n<p>Single validator to handle all date_time
    parsing cases</p>\n</div>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">parse_date_time <em class='small'>source</em></p>\n<pre
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
    class=\"w\"> </span><span class=\"nf\">parse_date_time</span><span class=\"p\">(</span><span
    class=\"bp\">cls</span><span class=\"p\">,</span> <span class=\"n\">v</span><span
    class=\"p\">,</span> <span class=\"n\">info</span><span class=\"p\">):</span>\n<span
    class=\"w\">        </span><span class=\"sd\">&quot;&quot;&quot;Single validator
    to handle all date_time parsing cases&quot;&quot;&quot;</span>\n        <span
    class=\"c1\"># If we have an explicit date_time value</span>\n        <span class=\"k\">if</span>
    <span class=\"n\">v</span> <span class=\"ow\">is</span> <span class=\"ow\">not</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">v</span><span
    class=\"p\">,</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">datetime</span><span class=\"p\">):</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">v</span>\n            <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">v</span><span
    class=\"p\">,</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">date</span><span class=\"p\">):</span>\n                <span class=\"k\">return</span>
    <span class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">combine</span><span class=\"p\">(</span><span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">time</span><span class=\"o\">.</span><span
    class=\"n\">min</span><span class=\"p\">)</span>\n            <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">v</span><span
    class=\"p\">,</span> <span class=\"nb\">str</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">try</span><span class=\"p\">:</span>\n                    <span
    class=\"c1\"># Try ISO format first</span>\n                    <span class=\"k\">return</span>
    <span class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">fromisoformat</span><span class=\"p\">(</span><span
    class=\"n\">v</span><span class=\"o\">.</span><span class=\"n\">replace</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;Z&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;+00:00&quot;</span><span class=\"p\">))</span>\n                <span
    class=\"k\">except</span> <span class=\"ne\">ValueError</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">try</span><span class=\"p\">:</span>\n                        <span
    class=\"k\">return</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">strptime</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;%Y-%m-</span><span class=\"si\">%d</span><span class=\"s2\">
    %H:%M&quot;</span><span class=\"p\">)</span>\n                    <span class=\"k\">except</span>
    <span class=\"ne\">ValueError</span><span class=\"p\">:</span>\n                        <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                            <span
    class=\"k\">return</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">strptime</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;%Y-%m-</span><span class=\"si\">%d</span><span class=\"s2\">&quot;</span><span
    class=\"p\">)</span>\n                        <span class=\"k\">except</span>
    <span class=\"ne\">ValueError</span><span class=\"p\">:</span>\n                            <span
    class=\"c1\"># Try dateparser as last resort for explicit date_time</span>\n                            <span
    class=\"kn\">import</span><span class=\"w\"> </span><span class=\"nn\">dateparser</span>\n\n
    \                           <span class=\"n\">parsed</span> <span class=\"o\">=</span>
    <span class=\"n\">dateparser</span><span class=\"o\">.</span><span class=\"n\">parse</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">)</span>\n                            <span
    class=\"k\">if</span> <span class=\"n\">parsed</span><span class=\"p\">:</span>\n
    \                               <span class=\"k\">return</span> <span class=\"n\">parsed</span>\n
    \                           <span class=\"k\">return</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">now</span><span class=\"p\">()</span>\n\n        <span class=\"c1\">#
    Get the raw date string directly from raw_date field</span>\n        <span class=\"n\">raw_date</span>
    <span class=\"o\">=</span> <span class=\"n\">info</span><span class=\"o\">.</span><span
    class=\"n\">data</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;raw_date&quot;</span><span class=\"p\">)</span>\n
    \       <span class=\"k\">if</span> <span class=\"n\">raw_date</span> <span class=\"ow\">and</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">raw_date</span><span
    class=\"p\">,</span> <span class=\"nb\">str</span><span class=\"p\">):</span>\n
    \           <span class=\"k\">try</span><span class=\"p\">:</span>\n                <span
    class=\"c1\"># Try ISO format first</span>\n                <span class=\"k\">return</span>
    <span class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">fromisoformat</span><span class=\"p\">(</span><span
    class=\"n\">raw_date</span><span class=\"o\">.</span><span class=\"n\">replace</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;Z&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;+00:00&quot;</span><span class=\"p\">))</span>\n            <span
    class=\"k\">except</span> <span class=\"ne\">ValueError</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">try</span><span class=\"p\">:</span>\n                    <span
    class=\"c1\"># Try parsing raw_date with time first</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">strptime</span><span
    class=\"p\">(</span><span class=\"n\">raw_date</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;%Y-%m-</span><span class=\"si\">%d</span><span class=\"s2\">
    %H:%M&quot;</span><span class=\"p\">)</span>\n                <span class=\"k\">except</span>
    <span class=\"ne\">ValueError</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                        <span
    class=\"c1\"># Fallback to date only</span>\n                        <span class=\"k\">return</span>
    <span class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">strptime</span><span class=\"p\">(</span><span
    class=\"n\">raw_date</span><span class=\"p\">,</span> <span class=\"s2\">&quot;%Y-%m-</span><span
    class=\"si\">%d</span><span class=\"s2\">&quot;</span><span class=\"p\">)</span>\n
    \                   <span class=\"k\">except</span> <span class=\"ne\">ValueError</span><span
    class=\"p\">:</span>\n                        <span class=\"c1\"># Try dateparser
    as last resort</span>\n                        <span class=\"kn\">import</span><span
    class=\"w\"> </span><span class=\"nn\">dateparser</span>\n\n                        <span
    class=\"n\">parsed</span> <span class=\"o\">=</span> <span class=\"n\">dateparser</span><span
    class=\"o\">.</span><span class=\"n\">parse</span><span class=\"p\">(</span><span
    class=\"n\">raw_date</span><span class=\"p\">)</span>\n                        <span
    class=\"k\">if</span> <span class=\"n\">parsed</span><span class=\"p\">:</span>\n
    \                           <span class=\"k\">return</span> <span class=\"n\">parsed</span>\n\n
    \       <span class=\"c1\"># If no raw_date, try to derive from date field</span>\n
    \       <span class=\"n\">date</span> <span class=\"o\">=</span> <span class=\"n\">info</span><span
    class=\"o\">.</span><span class=\"n\">data</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;date&quot;</span><span
    class=\"p\">)</span>\n        <span class=\"k\">if</span> <span class=\"n\">date</span><span
    class=\"p\">:</span>\n            <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">date</span><span class=\"p\">,</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"p\">):</span>\n                <span class=\"k\">return</span> <span class=\"n\">date</span>\n
    \           <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">date</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">):</span>\n                <span class=\"k\">try</span><span
    class=\"p\">:</span>\n                    <span class=\"c1\"># Try ISO format
    first</span>\n                    <span class=\"k\">return</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">fromisoformat</span><span class=\"p\">(</span><span class=\"n\">date</span><span
    class=\"o\">.</span><span class=\"n\">replace</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;Z&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;+00:00&quot;</span><span
    class=\"p\">))</span>\n                <span class=\"k\">except</span> <span class=\"ne\">ValueError</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                       <span class=\"c1\"># Try parsing date with time first</span>\n
    \                       <span class=\"k\">return</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">strptime</span><span class=\"p\">(</span><span class=\"n\">date</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;%Y-%m-</span><span class=\"si\">%d</span><span
    class=\"s2\"> %H:%M&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">except</span> <span class=\"ne\">ValueError</span><span class=\"p\">:</span>\n
    \                       <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                           <span class=\"c1\"># Fallback to date only</span>\n
    \                           <span class=\"k\">return</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">strptime</span><span class=\"p\">(</span><span class=\"n\">date</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;%Y-%m-</span><span class=\"si\">%d</span><span
    class=\"s2\">&quot;</span><span class=\"p\">)</span>\n                        <span
    class=\"k\">except</span> <span class=\"ne\">ValueError</span><span class=\"p\">:</span>\n
    \                           <span class=\"c1\"># Try dateparser as last resort</span>\n
    \                           <span class=\"kn\">import</span><span class=\"w\">
    </span><span class=\"nn\">dateparser</span>\n\n                            <span
    class=\"n\">parsed</span> <span class=\"o\">=</span> <span class=\"n\">dateparser</span><span
    class=\"o\">.</span><span class=\"n\">parse</span><span class=\"p\">(</span><span
    class=\"n\">date</span><span class=\"p\">)</span>\n                            <span
    class=\"k\">if</span> <span class=\"n\">parsed</span><span class=\"p\">:</span>\n
    \                               <span class=\"k\">return</span> <span class=\"n\">parsed</span>\n
    \           <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">date</span><span class=\"p\">,</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">date</span><span
    class=\"p\">):</span>\n                <span class=\"k\">return</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">combine</span><span class=\"p\">(</span><span class=\"n\">date</span><span
    class=\"p\">,</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">time</span><span class=\"o\">.</span><span class=\"n\">min</span><span
    class=\"p\">)</span>\n\n        <span class=\"c1\"># If we still don&#39;t have
    a date, use now</span>\n        <span class=\"k\">return</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">now</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"__init__\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">__init__
    <em class=\"small\">method</em></h2>\n<p>include: post attributes to include by
    default in Post\nmodel serialization.\nrepr_include: post attributes to include
    by default in Post\nrepr.  If <code>repr_include</code> is None, it will default
    to\n<code>include</code>, but it is likely that you want less in the repr\nthan
    serialized output.</p>\n<p>example:</p>\n<pre class='wrapper'>\n\n<div class='filepath'>\n<p>markata.toml</p>\n\n<div
    class='right'>\n\n<button class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
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
    class=\"highlight\"><pre><span></span><span class=\"k\">[markata.post_model]</span>\n<span
    class=\"n\">include</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"p\">[</span><span class=\"s1\">&#39;date&#39;</span><span
    class=\"p\">,</span><span class=\"w\"> </span><span class=\"s1\">&#39;description&#39;</span><span
    class=\"p\">,</span><span class=\"w\"> </span><span class=\"s1\">&#39;published&#39;</span><span
    class=\"p\">,</span>\n<span class=\"w\">    </span><span class=\"s1\">&#39;slug&#39;</span><span
    class=\"p\">,</span><span class=\"w\"> </span><span class=\"s1\">&#39;title&#39;</span><span
    class=\"p\">,</span><span class=\"w\"> </span><span class=\"s1\">&#39;content&#39;</span><span
    class=\"p\">,</span><span class=\"w\"> </span><span class=\"s1\">&#39;html&#39;</span><span
    class=\"p\">]</span>\n<span class=\"n\">repr_include</span><span class=\"w\">
    </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"p\">[</span><span
    class=\"s1\">&#39;date&#39;</span><span class=\"p\">,</span><span class=\"w\">
    </span><span class=\"s1\">&#39;description&#39;</span><span class=\"p\">,</span><span
    class=\"w\"> </span><span class=\"s1\">&#39;published&#39;</span><span class=\"p\">,</span><span
    class=\"w\"> </span><span class=\"s1\">&#39;slug&#39;</span><span class=\"p\">,</span><span
    class=\"w\"> </span><span class=\"s1\">&#39;title&#39;</span><span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\"><strong>init</strong>
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
    class=\"w\"> </span><span class=\"fm\">__init__</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"n\">data</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span class=\"w\">        </span><span
    class=\"sd\">&quot;&quot;&quot;</span>\n\n<span class=\"sd\">        include:
    post attributes to include by default in Post</span>\n<span class=\"sd\">        model
    serialization.</span>\n<span class=\"sd\">        repr_include: post attributes
    to include by default in Post</span>\n<span class=\"sd\">        repr.  If `repr_include`
    is None, it will default to</span>\n<span class=\"sd\">        `include`, but
    it is likely that you want less in the repr</span>\n<span class=\"sd\">        than
    serialized output.</span>\n\n<span class=\"sd\">        example:</span>\n\n<span
    class=\"sd\">        ``` toml title=&#39;markata.toml&#39;</span>\n<span class=\"sd\">
    \       [markata.post_model]</span>\n<span class=\"sd\">        include = [&#39;date&#39;,
    &#39;description&#39;, &#39;published&#39;,</span>\n<span class=\"sd\">            &#39;slug&#39;,
    &#39;title&#39;, &#39;content&#39;, &#39;html&#39;]</span>\n<span class=\"sd\">
    \       repr_include = [&#39;date&#39;, &#39;description&#39;, &#39;published&#39;,
    &#39;slug&#39;, &#39;title&#39;]</span>\n<span class=\"sd\">        ```</span>\n<span
    class=\"sd\">        &quot;&quot;&quot;</span>\n        <span class=\"nb\">super</span><span
    class=\"p\">()</span><span class=\"o\">.</span><span class=\"fm\">__init__</span><span
    class=\"p\">(</span><span class=\"o\">**</span><span class=\"n\">data</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n\n    </section>\n</article>"
  raw.md: ''
published: false
slug: markata/plugins/post-model
title: post_model.py


---

---

The `markata.plugins.post_model` plugin defines the core Post model used throughout
Markata. It provides robust validation, serialization, and configuration options for
all post attributes.

## Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.post_model",
]
```

## Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.post_model",
]
```

Note: Disabling this plugin will break most of Markata's functionality as the Post
model is fundamental to the system.

## Configuration

Configure post model behavior in your `markata.toml`:

```toml
[markata.post_model]
# Attributes to include when serializing posts
include = [
    "date",
    "description",
    "published",
    "slug",
    "title",
    "content",
    "html"
]

# Attributes to show in post representations
repr_include = [
    "date",
    "description",
    "published",
    "slug",
    "title"
]

# Attributes to include when exporting
export_include = [
    "date",
    "description",
    "published",
    "slug",
    "title"
]
```

## Functionality

## Post Model

Core attributes:
- `path`: Path to source file
- `slug`: URL-friendly identifier
- `href`: Full URL path
- `published`: Publication status
- `description`: Post summary
- `content`: Raw markdown content
- `html`: Rendered HTML content
- `tags`: List of post tags
- `date`: Publication date
- `title`: Post title

## Validation

The model provides:
- Type checking and coercion
- Required field validation
- Custom field validators
- Default values
- Rich error messages

## Serialization

Supports multiple output formats:
- Full serialization (all fields)
- Representation (subset for display)
- Export (subset for external use)
- JSON/YAML compatible

## Performance

Uses optimized Pydantic config:
- Disabled assignment validation
- Arbitrary types allowed
- Extra fields allowed
- String whitespace stripping
- Default value validation
- Number to string coercion
- Alias population

## Dependencies

This plugin depends on:
- pydantic for model definition
- rich for console output
- yaml for YAML handling

---

!!! class
    <h2 id="PostModelConfig" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">PostModelConfig <em class="small">class</em></h2>

    Configuration for the Post model

???+ source "PostModelConfig <em class='small'>source</em>"
    ```python
    class PostModelConfig(pydantic.BaseModel):
        """Configuration for the Post model"""

        def __init__(self, **data) -> None:
            """

            include: post attributes to include by default in Post
            model serialization.
            repr_include: post attributes to include by default in Post
            repr.  If `repr_include` is None, it will default to
            `include`, but it is likely that you want less in the repr
            than serialized output.

            example:

            ``` toml title='markata.toml'
            [markata.post_model]
            include = ['date', 'description', 'published',
                'slug', 'title', 'content', 'html']
            repr_include = ['date', 'description', 'published', 'slug', 'title']
            ```
            """
            super().__init__(**data)

        default_date: datetime.date = datetime.date.today()
        include: List[str] = [
            "date",
            "description",
            "published",
            "slug",
            "title",
            "content",
            "html",
        ]
        repr_include: Optional[List[str]] = [
            "date",
            "description",
            "published",
            "slug",
            "title",
        ]
        export_include: Optional[List[str]] = [
            "date",
            "description",
            "published",
            "slug",
            "title",
        ]

        model_config = ConfigDict(
            validate_assignment=True,  # Config model
            arbitrary_types_allowed=True,
            extra="allow",
            str_strip_whitespace=True,
            validate_default=True,
            coerce_numbers_to_str=True,
            populate_by_name=True,
        )

        @field_validator("repr_include", mode="before")
        @classmethod
        def repr_include_validator(cls, v, info) -> Optional[List[str]]:
            if v:
                return v
            return info.data.get("include")
    ```
!!! method
    <h2 id="metadata" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">metadata <em class="small">method</em></h2>

    for backwards compatability

???+ source "metadata <em class='small'>source</em>"
    ```python
    def metadata(self: "Post") -> Dict:
            "for backwards compatability"
            return self.__dict__
    ```
!!! method
    <h2 id="to_dict" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">to_dict <em class="small">method</em></h2>

    for backwards compatability

???+ source "to_dict <em class='small'>source</em>"
    ```python
    def to_dict(self: "Post") -> Dict:
            "for backwards compatability"
            return self.__dict__
    ```
!!! method
    <h2 id="__getitem__" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">__getitem__ <em class="small">method</em></h2>

    for backwards compatability

???+ source "__getitem__ <em class='small'>source</em>"
    ```python
    def __getitem__(self: "Post", item: str) -> Any:
            "for backwards compatability"
            return getattr(self, item)
    ```
!!! method
    <h2 id="__setitem__" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">__setitem__ <em class="small">method</em></h2>

    for backwards compatability

???+ source "__setitem__ <em class='small'>source</em>"
    ```python
    def __setitem__(self: "Post", key: str, item: Any) -> None:
            "for backwards compatability"
            setattr(self, key, item)
    ```
!!! method
    <h2 id="get" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">get <em class="small">method</em></h2>

    for backwards compatability

???+ source "get <em class='small'>source</em>"
    ```python
    def get(self: "Post", item: str, default: Any) -> Any:
            "for backwards compatability"
            return getattr(self, item, default)
    ```
!!! method
    <h2 id="keys" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">keys <em class="small">method</em></h2>

    for backwards compatability

???+ source "keys <em class='small'>source</em>"
    ```python
    def keys(self: "Post") -> List[str]:
            "for backwards compatability"
            return self.__dict__.keys()
    ```
!!! method
    <h2 id="yaml" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">yaml <em class="small">method</em></h2>

    dump model to yaml

???+ source "yaml <em class='small'>source</em>"
    ```python
    def yaml(self: "Post") -> str:
            """
            dump model to yaml
            """
            import yaml

            return yaml.dump(
                self.dict(
                    include={i: True for i in self.markata.config.post_model.include}
                ),
                Dumper=yaml.CDumper,
            )
    ```
!!! method
    <h2 id="markdown" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">markdown <em class="small">method</em></h2>

    dump model to markdown

???+ source "markdown <em class='small'>source</em>"
    ```python
    def markdown(self: "Post") -> str:
            """
            dump model to markdown
            """

            import yaml

            frontmatter = yaml.dump(
                self.dict(
                    include={
                        i: True
                        for i in [
                            _i
                            for _i in self.markata.config.post_model.include
                            if _i != "content"
                        ]
                    }
                ),
                Dumper=yaml.CDumper,
            )
            post = "---\n"
            post += frontmatter
            post += "---\n\n"

            if self.content:
                post += self.content
            return post
    ```
!!! method
    <h2 id="dumps" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">dumps <em class="small">method</em></h2>

    dumps raw article back out

???+ source "dumps <em class='small'>source</em>"
    ```python
    def dumps(self):
            """
            dumps raw article back out
            """
            return f"---\n{self.yaml()}\n\n---\n\n{self.content}"
    ```
!!! method
    <h2 id="parse_date_time" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">parse_date_time <em class="small">method</em></h2>

    Single validator to handle all date_time parsing cases

???+ source "parse_date_time <em class='small'>source</em>"
    ```python
    def parse_date_time(cls, v, info):
            """Single validator to handle all date_time parsing cases"""
            # If we have an explicit date_time value
            if v is not None:
                if isinstance(v, datetime.datetime):
                    return v
                if isinstance(v, datetime.date):
                    return datetime.datetime.combine(v, datetime.time.min)
                if isinstance(v, str):
                    try:
                        # Try ISO format first
                        return datetime.datetime.fromisoformat(v.replace("Z", "+00:00"))
                    except ValueError:
                        try:
                            return datetime.datetime.strptime(v, "%Y-%m-%d %H:%M")
                        except ValueError:
                            try:
                                return datetime.datetime.strptime(v, "%Y-%m-%d")
                            except ValueError:
                                # Try dateparser as last resort for explicit date_time
                                import dateparser

                                parsed = dateparser.parse(v)
                                if parsed:
                                    return parsed
                                return datetime.datetime.now()

            # Get the raw date string directly from raw_date field
            raw_date = info.data.get("raw_date")
            if raw_date and isinstance(raw_date, str):
                try:
                    # Try ISO format first
                    return datetime.datetime.fromisoformat(raw_date.replace("Z", "+00:00"))
                except ValueError:
                    try:
                        # Try parsing raw_date with time first
                        return datetime.datetime.strptime(raw_date, "%Y-%m-%d %H:%M")
                    except ValueError:
                        try:
                            # Fallback to date only
                            return datetime.datetime.strptime(raw_date, "%Y-%m-%d")
                        except ValueError:
                            # Try dateparser as last resort
                            import dateparser

                            parsed = dateparser.parse(raw_date)
                            if parsed:
                                return parsed

            # If no raw_date, try to derive from date field
            date = info.data.get("date")
            if date:
                if isinstance(date, datetime.datetime):
                    return date
                if isinstance(date, str):
                    try:
                        # Try ISO format first
                        return datetime.datetime.fromisoformat(date.replace("Z", "+00:00"))
                    except ValueError:
                        try:
                            # Try parsing date with time first
                            return datetime.datetime.strptime(date, "%Y-%m-%d %H:%M")
                        except ValueError:
                            try:
                                # Fallback to date only
                                return datetime.datetime.strptime(date, "%Y-%m-%d")
                            except ValueError:
                                # Try dateparser as last resort
                                import dateparser

                                parsed = dateparser.parse(date)
                                if parsed:
                                    return parsed
                if isinstance(date, datetime.date):
                    return datetime.datetime.combine(date, datetime.time.min)

            # If we still don't have a date, use now
            return datetime.datetime.now()
    ```
!!! method
    <h2 id="__init__" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">__init__ <em class="small">method</em></h2>

    include: post attributes to include by default in Post
    model serialization.
    repr_include: post attributes to include by default in Post
    repr.  If `repr_include` is None, it will default to
    `include`, but it is likely that you want less in the repr
    than serialized output.

    example:

    ``` toml title='markata.toml'
    [markata.post_model]
    include = ['date', 'description', 'published',
        'slug', 'title', 'content', 'html']
    repr_include = ['date', 'description', 'published', 'slug', 'title']
    ```

???+ source "__init__ <em class='small'>source</em>"
    ```python
    def __init__(self, **data) -> None:
            """

            include: post attributes to include by default in Post
            model serialization.
            repr_include: post attributes to include by default in Post
            repr.  If `repr_include` is None, it will default to
            `include`, but it is likely that you want less in the repr
            than serialized output.

            example:

            ``` toml title='markata.toml'
            [markata.post_model]
            include = ['date', 'description', 'published',
                'slug', 'title', 'content', 'html']
            repr_include = ['date', 'description', 'published', 'slug', 'title']
            ```
            """
            super().__init__(**data)
    ```