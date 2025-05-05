---
content: "---\n\nStandard Config.\nA module to load tooling config from a users project
  space.\n\nInspired from frustrations that some tools have a tool.ini, .tool.ini,\nsetup.cfg,
  or pyproject.toml.  Some allow for global configs, some don't.  Some\nproperly follow
  the users home directory, others end up in a weird temp\ndirectory.  Windows home
  directory is only more confusing.  Some will even\nrespect the users `$XDG_HOME`
  directory.\n\nThis file is for any project that can be configured in plain text
  such as `ini`\nor `toml` and not requiring a .py file.  Just name your tool and
  let users put\nconfig where it makes sense to them, no need to figure out resolution
  order.\n\n## Usage:\n\n``` python\nfrom standard_config import load\n\n# Retrieve
  any overrides from the user\noverrides = {'setting': True}\nconfig = load('my_tool',
  overrides)\n```\n\n## Resolution Order\n\n* First global file with a tool key\n*
  First local file with a tool key\n* Environment variables prefixed with `TOOL`\n*
  Overrides\n\n### Tool Specific Ini files\n\nIni file formats must include a `<tool>`
  key.\n\n``` ini\n[my_tool]\nsetting = True\n```\n\n### pyproject.toml\n\nToml files
  must include a `tool.<tool>` key\n\n``` toml\n[tool.my_tool]\nsetting = True\n```\n\n###
  setup.cfg\n\nsetup.cfg files must include a `tool:<tool>` key\n\n``` ini\n[tool:my_tool]\nsetting
  = True\n```\n\n### global files to consider\n\n* <home>/tool.ini\n* <home>/.tool\n*
  <home>/.tool.ini\n* <home>/.config/tool.ini\n* <home>/.config/.tool\n* <home>/.config/.tool.ini\n\n###
  local files to consider\n\n* <project_home>/tool.ini\n* <project_home>/.tool\n*
  <project_home>/.tool.ini\n* <project_home>/pyproject.toml\n* <project_home>/setup.cfg\n\nMarkata's
  standard configuration system.\n\n## Configuration Overview\n\nMarkata uses a hierarchical
  configuration system based on Pydantic models. Configuration\ncan be set through:\n1.
  TOML files\n2. Environment variables\n3. Command line arguments\n\n# Basic Configuration\n\nMinimal
  `markata.toml`:\n```toml\n[markata]\n# Site info\ntitle = \"My Site\"\nurl = \"https://example.com\"\ndescription
  = \"Site description\"\n\n# Content locations\ncontent_dir = \"content\"\noutput_dir
  = \"markout\"\nassets_dir = \"static\"\n\n# Plugin management\nhooks = [\"default\"]\n```\n\n#
  Environment Variables\n\nAll settings can be overridden with environment variables:\n```bash\n#
  Override site URL\nexport MARKATA_URL=\"https://staging.example.com\"\n\n# Override
  output directory\nexport MARKATA_OUTPUT_DIR=\"dist\"\n\n# Enable debug mode\nexport
  MARKATA_DEBUG=1\n```\n\n# Detailed Configuration\n\n## Core Settings\n\n```toml\n[markata]\n#
  Site information\ntitle = \"My Site\"                  # Site title\nurl = \"https://example.com\"
  \       # Base URL\ndescription = \"Site description\"   # Meta description\nauthor_name
  = \"Author Name\"        # Author name\nauthor_email = \"me@example.com\"    # Author
  email\nicon = \"favicon.ico\"               # Site icon\nlang = \"en\"                        #
  Site language\n\n# Content locations\ncontent_dir = \"content\"           # Source
  content location\noutput_dir = \"markout\"            # Build output location\nassets_dir
  = \"static\"             # Static assets location\ntemplate_dir = \"templates\"
  \       # Template location\n\n# Plugin management\nhooks = [\"default\"]               #
  Active plugins\ndisabled_hooks = []               # Disabled plugins\n```\n\n##
  Cache Settings\n\n```toml\n[markata]\n# Cache configuration\ndefault_cache_expire
  = 3600       # Default TTL (1 hour)\ntemplate_cache_expire = 86400     # Template
  TTL (24 hours)\nmarkdown_cache_expire = 21600     # Markdown TTL (6 hours)\ndynamic_cache_expire
  = 3600       # Dynamic TTL (1 hour)\n```\n\n## Development Settings\n\n```toml\n[markata]\n#
  Development server\ndev_server_port = 8000            # Local server port\ndev_server_host
  = \"localhost\"     # Local server host\ndebug = false                     # Debug
  mode\n\n# Performance\nparallel = true                   # Enable parallel processing\nworkers
  = 4                       # Number of worker threads\n```\n\n## Content Settings\n\n```toml\n[markata]\n#
  Content processing\ndefault_template = \"post.html\"    # Default template\nmarkdown_extensions
  = [           # Markdown extensions\n    \"fenced_code\",\n    \"tables\",\n    \"footnotes\"\n]\n\n#
  Content filtering\ndraft = false                     # Include drafts\nfuture =
  false                    # Include future posts\n```\n\n# Plugin Configuration\n\nEach
  plugin can define its own configuration section:\n\n```toml\n# RSS feed configuration\n[markata.feeds]\nrss
  = { output = \"rss.xml\" }\natom = { output = \"atom.xml\" }\njson = { output =
  \"feed.json\" }\n\n# Template configuration\n[markata.template]\nengine = \"jinja2\"\ncache_size
  = 100\nautoescape = true\n\n# Markdown configuration\n[markata.markdown]\nhighlight_theme
  = \"monokai\"\nline_numbers = true\n```\n\n## Configuration Validation\n\nThe configuration
  is validated using Pydantic models:\n\n```python\nfrom pydantic import BaseModel,
  Field\n\nclass MarkataConfig(BaseModel):\n    \"\"\"Core configuration model.\"\"\"\n
  \   # Site info\n    title: str = Field(..., description=\"Site title\")\n    url:
  str = Field(..., description=\"Site base URL\")\n\n    # Directories\n    content_dir:
  Path = Field(\"content\", description=\"Content directory\")\n    output_dir: Path
  = Field(\"markout\", description=\"Output directory\")\n\n    # Features\n    debug:
  bool = Field(False, description=\"Enable debug mode\")\n    parallel: bool = Field(True,
  description=\"Enable parallel processing\")\n\n    model_config = ConfigDict(\n
  \       validate_assignment=True,\n        arbitrary_types_allowed=True,\n        extra=\"allow\",\n
  \       str_strip_whitespace=True,\n        validate_default=True,\n        populate_by_name=True,\n
  \   )\n```\n\n# Usage Example\n\n```python\nfrom markata import Markata\n\n# Load
  config from file\nmarkata = Markata.from_file(\"markata.toml\")\n\n# Access configuration\nprint(markata.config.title)
  \        # Site title\nprint(markata.config.url)           # Site URL\nprint(markata.config.content_dir)
  \  # Content directory\n\n# Access plugin config\nprint(markata.config.feeds.rss)
  \    # RSS feed config\nprint(markata.config.template)      # Template config\n\n#
  Override config\nmarkata.config.debug = True\nmarkata.config.parallel = False\n```\n\nSee
  hookspec.py for plugin development and lifecycle.py for build process details.\n\n---\n\n!!!
  function\n    <h2 id=\"_get_global_path_specs\" class=\"admonition-title\" style=\"margin:
  0; padding: .5rem 1rem;\">_get_global_path_specs <em class=\"small\">function</em></h2>\n\n
  \   Generate a list of standard pathspecs for global config files.\n\n    Args:\n
  \       tool (str): name of the tool to configure\n\n???+ source \"_get_global_path_specs
  <em class='small'>source</em>\"\n    ```python\n    def _get_global_path_specs(tool:
  str) -> path_spec_type:\n        \"\"\"\n        Generate a list of standard pathspecs
  for global config files.\n\n        Args:\n            tool (str): name of the tool
  to configure\n        \"\"\"\n        try:\n            home = Path(os.environ[\"XDG_HOME\"])\n
  \       except KeyError:\n            home = Path.home()\n\n        return [\n            {\"path_specs\":
  home / f\"{tool}.ini\", \"parser\": \"ini\", \"keys\": [tool]},\n            {\"path_specs\":
  home / f\".{tool}\", \"parser\": \"ini\", \"keys\": [tool]},\n            {\"path_specs\":
  home / f\".{tool}.ini\", \"parser\": \"ini\", \"keys\": [tool]},\n            {\n
  \               \"path_specs\": home / \".config\" / f\"{tool}.ini\",\n                \"parser\":
  \"ini\",\n                \"keys\": [tool],\n            },\n            {\n                \"path_specs\":
  home / \".config\" / f\".{tool}\",\n                \"parser\": \"ini\",\n                \"keys\":
  [tool],\n            },\n            {\n                \"path_specs\": home / \".config\"
  / f\".{tool}.ini\",\n                \"parser\": \"ini\",\n                \"keys\":
  [tool],\n            },\n        ]\n    ```\n!!! function\n    <h2 id=\"_get_local_path_specs\"
  class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">_get_local_path_specs
  <em class=\"small\">function</em></h2>\n\n    Generate a list of standard pathspecs
  for local, project directory config files.\n\n    Args:\n        tool (str): name
  of the tool to configure\n\n???+ source \"_get_local_path_specs <em class='small'>source</em>\"\n
  \   ```python\n    def _get_local_path_specs(tool: str, project_home: Union[str,
  Path]) -> path_spec_type:\n        \"\"\"\n        Generate a list of standard pathspecs
  for local, project directory config files.\n\n        Args:\n            tool (str):
  name of the tool to configure\n        \"\"\"\n        return [\n            {\n
  \               \"path_specs\": Path(project_home) / f\"{tool}.ini\",\n                \"parser\":
  \"ini\",\n                \"keys\": [tool],\n            },\n            {\n                \"path_specs\":
  Path(project_home) / f\".{tool}\",\n                \"parser\": \"ini\",\n                \"keys\":
  [tool],\n            },\n            {\n                \"path_specs\": Path(project_home)
  / f\".{tool}.ini\",\n                \"parser\": \"ini\",\n                \"keys\":
  [tool],\n            },\n            {\n                \"path_specs\": Path(project_home)
  / f\"{tool}.yml\",\n                \"parser\": \"yaml\",\n                \"keys\":
  [tool],\n            },\n            {\n                \"path_specs\": Path(project_home)
  / f\".{tool}.yml\",\n                \"parser\": \"yaml\",\n                \"keys\":
  [tool],\n            },\n            {\n                \"path_specs\": Path(project_home)
  / f\"{tool}.toml\",\n                \"parser\": \"toml\",\n                \"keys\":
  [tool],\n            },\n            {\n                \"path_specs\": Path(project_home)
  / f\".{tool}.toml\",\n                \"parser\": \"toml\",\n                \"keys\":
  [tool],\n            },\n            {\n                \"path_specs\": Path(project_home)
  / \"pyproject.toml\",\n                \"parser\": \"toml\",\n                \"keys\":
  [\"tool\", tool],\n            },\n            {\n                \"path_specs\":
  Path(project_home) / \"setup.cfg\",\n                \"parser\": \"ini\",\n                \"keys\":
  [f\"tool.{tool}\"],\n            },\n        ]\n    ```\n!!! function\n    <h2 id=\"_get_attrs\"
  class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">_get_attrs
  <em class=\"small\">function</em></h2>\n\n    Get nested config data from a list
  of keys.\n\n    specifically written for pyproject.toml which needs to get `tool`
  then `<tool>`\n\n???+ source \"_get_attrs <em class='small'>source</em>\"\n    ```python\n
  \   def _get_attrs(attrs: list, config: Dict) -> Dict:\n        \"\"\"Get nested
  config data from a list of keys.\n\n        specifically written for pyproject.toml
  which needs to get `tool` then `<tool>`\n        \"\"\"\n        for attr in attrs:\n
  \           config = config[attr]\n        return config\n    ```\n!!! function\n
  \   <h2 id=\"_load_config_file\" class=\"admonition-title\" style=\"margin: 0; padding:
  .5rem 1rem;\">_load_config_file <em class=\"small\">function</em></h2>\n\n    Load
  a configuration file using the appropriate parser.\n\n    Args:\n        file_spec:
  Dictionary containing path_specs, parser, and keys information\n\n    Returns:\n
  \       Optional[Dict[str, Any]]: Parsed configuration or None if file doesn't exist\n\n???+
  source \"_load_config_file <em class='small'>source</em>\"\n    ```python\n    def
  _load_config_file(file_spec: Dict[str, Any]) -> Optional[Dict[str, Any]]:\n        \"\"\"Load
  a configuration file using the appropriate parser.\n\n        Args:\n            file_spec:
  Dictionary containing path_specs, parser, and keys information\n\n        Returns:\n
  \           Optional[Dict[str, Any]]: Parsed configuration or None if file doesn't
  exist\n        \"\"\"\n        path = file_spec[\"path_specs\"]\n        if not
  path.exists():\n            return None\n\n        try:\n            if file_spec[\"parser\"]
  == \"toml\":\n                with open(path, \"rb\") as f:\n                    config
  = tomli.load(f)\n            elif file_spec[\"parser\"] == \"yaml\":\n                with
  open(path, \"r\") as f:\n                    config = yaml.safe_load(f)\n            elif
  file_spec[\"parser\"] == \"ini\":\n                config = configparser.ConfigParser()\n
  \               config.read(path)\n                # Convert ConfigParser to dict\n
  \               config = {s: dict(config.items(s)) for s in config.sections()}\n
  \           else:\n                return None\n\n            return _get_attrs(file_spec[\"keys\"],
  config)\n        except (\n            KeyError,\n            TypeError,\n            yaml.YAMLError,\n
  \           tomli.TOMLDecodeError,\n            configparser.Error,\n        ):\n
  \           return None\n    ```\n!!! function\n    <h2 id=\"_load_files\" class=\"admonition-title\"
  style=\"margin: 0; padding: .5rem 1rem;\">_load_files <em class=\"small\">function</em></h2>\n\n
  \   Load config files stopping at the first one that exists and can be parsed.\n\n
  \   Args:\n        config_path_specs: List of path specifications to try\n\n    Returns:\n
  \       Dict[str, Any]: Configuration dictionary\n\n???+ source \"_load_files <em
  class='small'>source</em>\"\n    ```python\n    def _load_files(config_path_specs:
  path_spec_type) -> Dict[str, Any]:\n        \"\"\"Load config files stopping at
  the first one that exists and can be parsed.\n\n        Args:\n            config_path_specs:
  List of path specifications to try\n\n        Returns:\n            Dict[str, Any]:
  Configuration dictionary\n        \"\"\"\n        for file_spec in config_path_specs:\n
  \           config = _load_config_file(file_spec)\n            if config:\n                return
  config\n        return {}\n    ```\n!!! function\n    <h2 id=\"_load_env\" class=\"admonition-title\"
  style=\"margin: 0; padding: .5rem 1rem;\">_load_env <em class=\"small\">function</em></h2>\n\n
  \   Load config from environment variables.\n\n    Args:\n        tool (str): name
  of the tool to configure\n\n???+ source \"_load_env <em class='small'>source</em>\"\n
  \   ```python\n    def _load_env(tool: str) -> Dict[str, Any]:\n        \"\"\"Load
  config from environment variables.\n\n        Args:\n            tool (str): name
  of the tool to configure\n        \"\"\"\n        env_prefix = tool.upper()\n        env_config
  = {\n            key.replace(f\"{env_prefix}_\", \"\").lower(): value\n            for
  key, value in os.environ.items()\n            if key.startswith(f\"{env_prefix}_\")\n
  \       }\n        return env_config\n    ```\n!!! function\n    <h2 id=\"load\"
  class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">load <em class=\"small\">function</em></h2>\n\n
  \   Load tool config from standard config files.\n\n    Resolution Order\n\n    *
  First global file with a tool key\n    * First local file with a tool key\n    *
  Environment variables prefixed with `TOOL`\n    * Overrides\n\n    Args:\n        tool
  (str): name of the tool to configure\n        project_home (Union[Path, str], optional):
  Project directory to search for config files. Defaults to \".\".\n        overrides
  (Dict, optional): Override values to apply last. Defaults to None.\n\n    Returns:\n
  \       Dict[str, Any]: Configuration object\n\n???+ source \"load <em class='small'>source</em>\"\n
  \   ```python\n    def load(\n        tool: str,\n        project_home: Union[Path,
  str] = \".\",\n        overrides: Optional[Dict[str, Any]] = None,\n    ) -> Dict[str,
  Any]:\n        \"\"\"Load tool config from standard config files.\n\n        Resolution
  Order\n\n        * First global file with a tool key\n        * First local file
  with a tool key\n        * Environment variables prefixed with `TOOL`\n        *
  Overrides\n\n        Args:\n            tool (str): name of the tool to configure\n
  \           project_home (Union[Path, str], optional): Project directory to search
  for config files. Defaults to \".\".\n            overrides (Dict, optional): Override
  values to apply last. Defaults to None.\n\n        Returns:\n            Dict[str,
  Any]: Configuration object\n        \"\"\"\n        overrides = overrides or {}\n
  \       config = {}\n\n        # Load from files in order of precedence\n        config.update(_load_files(_get_global_path_specs(tool))
  or {})\n        config.update(_load_files(_get_local_path_specs(tool, project_home))
  or {})\n        config.update(_load_env(tool))\n        config.update(overrides)\n\n
  \       # If no settings class is provided, return the raw dict\n        return
  config\n    ```"
date: 2025-05-05
description: "Standard Config. A module to load tooling config from a users project
  space. Inspired from frustrations that some tools have a tool.ini, .tool.ini, setup.cfg,\u2026"
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>standard_config.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Standard Config. A module to load tooling
    config from a users project space. Inspired from frustrations that some tools
    have a tool.ini, .tool.ini, setup.cfg,\u2026\" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n
    \       <meta property=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    <head>\n        <script src=\"https://unpkg.com/@tailwindcss/browser@4\"></script>\n<title>standard_config.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Standard Config. A module to load tooling
    config from a users project space. Inspired from frustrations that some tools
    have a tool.ini, .tool.ini, setup.cfg,\u2026\" />\n <link href=\"/favicon.ico\"
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        standard_config.py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <hr />\n<p>Standard
    Config.\nA module to load tooling config from a users project space.</p>\n<p>Inspired
    from frustrations that some tools have a tool.ini, .tool.ini,\nsetup.cfg, or pyproject.toml.
    \ Some allow for global configs, some don't.  Some\nproperly follow the users
    home directory, others end up in a weird temp\ndirectory.  Windows home directory
    is only more confusing.  Some will even\nrespect the users <code>$XDG_HOME</code>
    directory.</p>\n<p>This file is for any project that can be configured in plain
    text such as <code>ini</code>\nor <code>toml</code> and not requiring a .py file.
    \ Just name your tool and let users put\nconfig where it makes sense to them,
    no need to figure out resolution order.</p>\n<h2 id=\"usage\">Usage: <a class=\"header-anchor\"
    href=\"#usage\"><svg class=\"heading-permalink\" aria-hidden=\"true\" fill=\"currentColor\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"kn\">from</span><span
    class=\"w\"> </span><span class=\"nn\">standard_config</span><span class=\"w\">
    </span><span class=\"kn\">import</span> <span class=\"n\">load</span>\n\n<span
    class=\"c1\"># Retrieve any overrides from the user</span>\n<span class=\"n\">overrides</span>
    <span class=\"o\">=</span> <span class=\"p\">{</span><span class=\"s1\">&#39;setting&#39;</span><span
    class=\"p\">:</span> <span class=\"kc\">True</span><span class=\"p\">}</span>\n<span
    class=\"n\">config</span> <span class=\"o\">=</span> <span class=\"n\">load</span><span
    class=\"p\">(</span><span class=\"s1\">&#39;my_tool&#39;</span><span class=\"p\">,</span>
    <span class=\"n\">overrides</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<h2
    id=\"resolution-order\">Resolution Order <a class=\"header-anchor\" href=\"#resolution-order\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<ul>\n<li>First global
    file with a tool key</li>\n<li>First local file with a tool key</li>\n<li>Environment
    variables prefixed with <code>TOOL</code></li>\n<li>Overrides</li>\n</ul>\n<h3>Tool
    Specific Ini files</h3>\n<p>Ini file formats must include a <code>&lt;tool&gt;</code>
    key.</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button class='copy'
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[my_tool]</span>\n<span
    class=\"na\">setting</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s\">True</span>\n</pre></div>\n\n</pre>\n\n<h3>pyproject.toml</h3>\n<p>Toml
    files must include a <code>tool.&lt;tool&gt;</code> key</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[tool.my_tool]</span>\n<span
    class=\"n\">setting</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"err\">True</span>\n</pre></div>\n\n</pre>\n\n<h3>setup.cfg</h3>\n<p>setup.cfg
    files must include a <code>tool:&lt;tool&gt;</code> key</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[tool:my_tool]</span>\n<span
    class=\"na\">setting</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s\">True</span>\n</pre></div>\n\n</pre>\n\n<h3>global
    files to consider</h3>\n<ul>\n<li><home>/tool.ini</li>\n<li><home>/.tool</li>\n<li><home>/.tool.ini</li>\n<li><home>/.config/tool.ini</li>\n<li><home>/.config/.tool</li>\n<li><home>/.config/.tool.ini</li>\n</ul>\n<h3>local
    files to consider</h3>\n<ul>\n<li>&lt;project_home&gt;/tool.ini</li>\n<li>&lt;project_home&gt;/.tool</li>\n<li>&lt;project_home&gt;/.tool.ini</li>\n<li>&lt;project_home&gt;/pyproject.toml</li>\n<li>&lt;project_home&gt;/setup.cfg</li>\n</ul>\n<p>Markata's
    standard configuration system.</p>\n<h2 id=\"configuration-overview\">Configuration
    Overview <a class=\"header-anchor\" href=\"#configuration-overview\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Markata uses a hierarchical
    configuration system based on Pydantic models. Configuration\ncan be set through:</p>\n<ol>\n<li>TOML
    files</li>\n<li>Environment variables</li>\n<li>Command line arguments</li>\n</ol>\n<h1
    id=\"basic-configuration\">Basic Configuration <a class=\"header-anchor\" href=\"#basic-configuration\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>Minimal <code>markata.toml</code>:</p>\n<pre
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
    class=\"c1\"># Site info</span>\n<span class=\"n\">title</span><span class=\"w\">
    </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s2\">&quot;My
    Site&quot;</span>\n<span class=\"n\">url</span><span class=\"w\"> </span><span
    class=\"o\">=</span><span class=\"w\"> </span><span class=\"s2\">&quot;https://example.com&quot;</span>\n<span
    class=\"n\">description</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;Site description&quot;</span>\n\n<span
    class=\"c1\"># Content locations</span>\n<span class=\"n\">content_dir</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;content&quot;</span>\n<span class=\"n\">output_dir</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;markout&quot;</span>\n<span class=\"n\">assets_dir</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;static&quot;</span>\n\n<span class=\"c1\"># Plugin management</span>\n<span
    class=\"n\">hooks</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"p\">[</span><span class=\"s2\">&quot;default&quot;</span><span
    class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<h1 id=\"environment-variables\">Environment
    Variables <a class=\"header-anchor\" href=\"#environment-variables\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>All settings can be
    overridden with environment variables:</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"c1\"># Override
    site URL</span>\n<span class=\"nb\">export</span><span class=\"w\"> </span><span
    class=\"nv\">MARKATA_URL</span><span class=\"o\">=</span><span class=\"s2\">&quot;https://staging.example.com&quot;</span>\n\n<span
    class=\"c1\"># Override output directory</span>\n<span class=\"nb\">export</span><span
    class=\"w\"> </span><span class=\"nv\">MARKATA_OUTPUT_DIR</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;dist&quot;</span>\n\n<span class=\"c1\"># Enable debug mode</span>\n<span
    class=\"nb\">export</span><span class=\"w\"> </span><span class=\"nv\">MARKATA_DEBUG</span><span
    class=\"o\">=</span><span class=\"m\">1</span>\n</pre></div>\n\n</pre>\n\n<h1
    id=\"detailed-configuration\">Detailed Configuration <a class=\"header-anchor\"
    href=\"#detailed-configuration\"><svg class=\"heading-permalink\" aria-hidden=\"true\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<h2 id=\"core-settings\">Core
    Settings <a class=\"header-anchor\" href=\"#core-settings\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n<span
    class=\"c1\"># Site information</span>\n<span class=\"n\">title</span><span class=\"w\">
    </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s2\">&quot;My
    Site&quot;</span><span class=\"w\">                  </span><span class=\"c1\">#
    Site title</span>\n<span class=\"n\">url</span><span class=\"w\"> </span><span
    class=\"o\">=</span><span class=\"w\"> </span><span class=\"s2\">&quot;https://example.com&quot;</span><span
    class=\"w\">        </span><span class=\"c1\"># Base URL</span>\n<span class=\"n\">description</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;Site description&quot;</span><span class=\"w\">   </span><span
    class=\"c1\"># Meta description</span>\n<span class=\"n\">author_name</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;Author Name&quot;</span><span class=\"w\">        </span><span
    class=\"c1\"># Author name</span>\n<span class=\"n\">author_email</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;me@example.com&quot;</span><span class=\"w\">    </span><span
    class=\"c1\"># Author email</span>\n<span class=\"n\">icon</span><span class=\"w\">
    </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s2\">&quot;favicon.ico&quot;</span><span
    class=\"w\">               </span><span class=\"c1\"># Site icon</span>\n<span
    class=\"n\">lang</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;en&quot;</span><span class=\"w\">
    \                       </span><span class=\"c1\"># Site language</span>\n\n<span
    class=\"c1\"># Content locations</span>\n<span class=\"n\">content_dir</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;content&quot;</span><span class=\"w\">           </span><span
    class=\"c1\"># Source content location</span>\n<span class=\"n\">output_dir</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;markout&quot;</span><span class=\"w\">            </span><span
    class=\"c1\"># Build output location</span>\n<span class=\"n\">assets_dir</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;static&quot;</span><span class=\"w\">             </span><span
    class=\"c1\"># Static assets location</span>\n<span class=\"n\">template_dir</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;templates&quot;</span><span class=\"w\">        </span><span
    class=\"c1\"># Template location</span>\n\n<span class=\"c1\"># Plugin management</span>\n<span
    class=\"n\">hooks</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"p\">[</span><span class=\"s2\">&quot;default&quot;</span><span
    class=\"p\">]</span><span class=\"w\">               </span><span class=\"c1\">#
    Active plugins</span>\n<span class=\"n\">disabled_hooks</span><span class=\"w\">
    </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"p\">[]</span><span
    class=\"w\">               </span><span class=\"c1\"># Disabled plugins</span>\n</pre></div>\n\n</pre>\n\n<h2
    id=\"cache-settings\">Cache Settings <a class=\"header-anchor\" href=\"#cache-settings\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n<span
    class=\"c1\"># Cache configuration</span>\n<span class=\"n\">default_cache_expire</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"mi\">3600</span><span class=\"w\">       </span><span class=\"c1\"># Default
    TTL (1 hour)</span>\n<span class=\"n\">template_cache_expire</span><span class=\"w\">
    </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"mi\">86400</span><span
    class=\"w\">     </span><span class=\"c1\"># Template TTL (24 hours)</span>\n<span
    class=\"n\">markdown_cache_expire</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"mi\">21600</span><span class=\"w\">     </span><span
    class=\"c1\"># Markdown TTL (6 hours)</span>\n<span class=\"n\">dynamic_cache_expire</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"mi\">3600</span><span class=\"w\">       </span><span class=\"c1\"># Dynamic
    TTL (1 hour)</span>\n</pre></div>\n\n</pre>\n\n<h2 id=\"development-settings\">Development
    Settings <a class=\"header-anchor\" href=\"#development-settings\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n<span
    class=\"c1\"># Development server</span>\n<span class=\"n\">dev_server_port</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"mi\">8000</span><span class=\"w\">            </span><span class=\"c1\">#
    Local server port</span>\n<span class=\"n\">dev_server_host</span><span class=\"w\">
    </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s2\">&quot;localhost&quot;</span><span
    class=\"w\">     </span><span class=\"c1\"># Local server host</span>\n<span class=\"n\">debug</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"kc\">false</span><span class=\"w\">                     </span><span class=\"c1\">#
    Debug mode</span>\n\n<span class=\"c1\"># Performance</span>\n<span class=\"n\">parallel</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"kc\">true</span><span class=\"w\">                   </span><span class=\"c1\">#
    Enable parallel processing</span>\n<span class=\"n\">workers</span><span class=\"w\">
    </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"mi\">4</span><span
    class=\"w\">                       </span><span class=\"c1\"># Number of worker
    threads</span>\n</pre></div>\n\n</pre>\n\n<h2 id=\"content-settings\">Content
    Settings <a class=\"header-anchor\" href=\"#content-settings\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n<span
    class=\"c1\"># Content processing</span>\n<span class=\"n\">default_template</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;post.html&quot;</span><span class=\"w\">    </span><span class=\"c1\">#
    Default template</span>\n<span class=\"n\">markdown_extensions</span><span class=\"w\">
    </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"p\">[</span><span
    class=\"w\">           </span><span class=\"c1\"># Markdown extensions</span>\n<span
    class=\"w\">    </span><span class=\"s2\">&quot;fenced_code&quot;</span><span
    class=\"p\">,</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;tables&quot;</span><span
    class=\"p\">,</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;footnotes&quot;</span>\n<span
    class=\"p\">]</span>\n\n<span class=\"c1\"># Content filtering</span>\n<span class=\"n\">draft</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"kc\">false</span><span class=\"w\">                     </span><span class=\"c1\">#
    Include drafts</span>\n<span class=\"n\">future</span><span class=\"w\"> </span><span
    class=\"o\">=</span><span class=\"w\"> </span><span class=\"kc\">false</span><span
    class=\"w\">                    </span><span class=\"c1\"># Include future posts</span>\n</pre></div>\n\n</pre>\n\n<h1
    id=\"plugin-configuration\">Plugin Configuration <a class=\"header-anchor\" href=\"#plugin-configuration\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>Each plugin can define
    its own configuration section:</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"c1\"># RSS
    feed configuration</span>\n<span class=\"k\">[markata.feeds]</span>\n<span class=\"n\">rss</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"p\">{</span><span class=\"w\"> </span><span class=\"n\">output</span><span
    class=\"w\"> </span><span class=\"p\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;rss.xml&quot;</span><span class=\"w\"> </span><span class=\"p\">}</span>\n<span
    class=\"n\">atom</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"p\">{</span><span class=\"w\"> </span><span
    class=\"n\">output</span><span class=\"w\"> </span><span class=\"p\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;atom.xml&quot;</span><span class=\"w\">
    </span><span class=\"p\">}</span>\n<span class=\"n\">json</span><span class=\"w\">
    </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"p\">{</span><span
    class=\"w\"> </span><span class=\"n\">output</span><span class=\"w\"> </span><span
    class=\"p\">=</span><span class=\"w\"> </span><span class=\"s2\">&quot;feed.json&quot;</span><span
    class=\"w\"> </span><span class=\"p\">}</span>\n\n<span class=\"c1\"># Template
    configuration</span>\n<span class=\"k\">[markata.template]</span>\n<span class=\"n\">engine</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;jinja2&quot;</span>\n<span class=\"n\">cache_size</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"mi\">100</span>\n<span class=\"n\">autoescape</span><span class=\"w\">
    </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"kc\">true</span>\n\n<span
    class=\"c1\"># Markdown configuration</span>\n<span class=\"k\">[markata.markdown]</span>\n<span
    class=\"n\">highlight_theme</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;monokai&quot;</span>\n<span class=\"n\">line_numbers</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"kc\">true</span>\n</pre></div>\n\n</pre>\n\n<h2 id=\"configuration-validation\">Configuration
    Validation <a class=\"header-anchor\" href=\"#configuration-validation\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>The configuration is
    validated using Pydantic models:</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    class=\"w\"> </span><span class=\"nn\">pydantic</span><span class=\"w\"> </span><span
    class=\"kn\">import</span> <span class=\"n\">BaseModel</span><span class=\"p\">,</span>
    <span class=\"n\">Field</span>\n\n<span class=\"k\">class</span><span class=\"w\">
    </span><span class=\"nc\">MarkataConfig</span><span class=\"p\">(</span><span
    class=\"n\">BaseModel</span><span class=\"p\">):</span>\n<span class=\"w\">    </span><span
    class=\"sd\">&quot;&quot;&quot;Core configuration model.&quot;&quot;&quot;</span>\n
    \   <span class=\"c1\"># Site info</span>\n    <span class=\"n\">title</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span> <span class=\"o\">=</span>
    <span class=\"n\">Field</span><span class=\"p\">(</span><span class=\"o\">...</span><span
    class=\"p\">,</span> <span class=\"n\">description</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;Site title&quot;</span><span class=\"p\">)</span>\n    <span
    class=\"n\">url</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"n\">Field</span><span class=\"p\">(</span><span
    class=\"o\">...</span><span class=\"p\">,</span> <span class=\"n\">description</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;Site base URL&quot;</span><span class=\"p\">)</span>\n\n
    \   <span class=\"c1\"># Directories</span>\n    <span class=\"n\">content_dir</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span> <span class=\"o\">=</span>
    <span class=\"n\">Field</span><span class=\"p\">(</span><span class=\"s2\">&quot;content&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">description</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;Content directory&quot;</span><span class=\"p\">)</span>\n
    \   <span class=\"n\">output_dir</span><span class=\"p\">:</span> <span class=\"n\">Path</span>
    <span class=\"o\">=</span> <span class=\"n\">Field</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;markout&quot;</span><span class=\"p\">,</span> <span class=\"n\">description</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;Output directory&quot;</span><span
    class=\"p\">)</span>\n\n    <span class=\"c1\"># Features</span>\n    <span class=\"n\">debug</span><span
    class=\"p\">:</span> <span class=\"nb\">bool</span> <span class=\"o\">=</span>
    <span class=\"n\">Field</span><span class=\"p\">(</span><span class=\"kc\">False</span><span
    class=\"p\">,</span> <span class=\"n\">description</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;Enable debug mode&quot;</span><span class=\"p\">)</span>\n
    \   <span class=\"n\">parallel</span><span class=\"p\">:</span> <span class=\"nb\">bool</span>
    <span class=\"o\">=</span> <span class=\"n\">Field</span><span class=\"p\">(</span><span
    class=\"kc\">True</span><span class=\"p\">,</span> <span class=\"n\">description</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;Enable parallel processing&quot;</span><span
    class=\"p\">)</span>\n\n    <span class=\"n\">model_config</span> <span class=\"o\">=</span>
    <span class=\"n\">ConfigDict</span><span class=\"p\">(</span>\n        <span class=\"n\">validate_assignment</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">arbitrary_types_allowed</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n        <span class=\"n\">extra</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;allow&quot;</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">str_strip_whitespace</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n        <span class=\"n\">validate_default</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">populate_by_name</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n    <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<h1
    id=\"usage-example\">Usage Example <a class=\"header-anchor\" href=\"#usage-example\"><svg
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"kn\">from</span><span
    class=\"w\"> </span><span class=\"nn\">markata</span><span class=\"w\"> </span><span
    class=\"kn\">import</span> <span class=\"n\">Markata</span>\n\n<span class=\"c1\">#
    Load config from file</span>\n<span class=\"n\">markata</span> <span class=\"o\">=</span>
    <span class=\"n\">Markata</span><span class=\"o\">.</span><span class=\"n\">from_file</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;markata.toml&quot;</span><span class=\"p\">)</span>\n\n<span
    class=\"c1\"># Access configuration</span>\n<span class=\"nb\">print</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">title</span><span
    class=\"p\">)</span>         <span class=\"c1\"># Site title</span>\n<span class=\"nb\">print</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">url</span><span
    class=\"p\">)</span>           <span class=\"c1\"># Site URL</span>\n<span class=\"nb\">print</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">content_dir</span><span
    class=\"p\">)</span>   <span class=\"c1\"># Content directory</span>\n\n<span
    class=\"c1\"># Access plugin config</span>\n<span class=\"nb\">print</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">feeds</span><span
    class=\"o\">.</span><span class=\"n\">rss</span><span class=\"p\">)</span>     <span
    class=\"c1\"># RSS feed config</span>\n<span class=\"nb\">print</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">template</span><span class=\"p\">)</span>
    \     <span class=\"c1\"># Template config</span>\n\n<span class=\"c1\"># Override
    config</span>\n<span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">debug</span>
    <span class=\"o\">=</span> <span class=\"kc\">True</span>\n<span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">parallel</span> <span class=\"o\">=</span> <span class=\"kc\">False</span>\n</pre></div>\n\n</pre>\n\n<p>See
    <a href=\"http://hookspec.py\">hookspec.py</a> for plugin development and <a href=\"http://lifecycle.py\">lifecycle.py</a>
    for build process details.</p>\n<hr />\n<div class=\"admonition function\">\n<p
    class=\"admonition-title\">Function</p>\n<h2 id=\"_get_global_path_specs\" class=\"admonition-title\"
    style=\"margin: 0; padding: .5rem 1rem;\">_get_global_path_specs <em class=\"small\">function</em></h2>\n<p>Generate
    a list of standard pathspecs for global config files.</p>\n<p>Args:\ntool (str):
    name of the tool to configure</p>\n</div>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">_get_global_path_specs <em
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
    class=\"w\"> </span><span class=\"nf\">_get_global_path_specs</span><span class=\"p\">(</span><span
    class=\"n\">tool</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">path_spec_type</span><span
    class=\"p\">:</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">    Generate a list of standard pathspecs for global config files.</span>\n\n<span
    class=\"sd\">    Args:</span>\n<span class=\"sd\">        tool (str): name of
    the tool to configure</span>\n<span class=\"sd\">    &quot;&quot;&quot;</span>\n
    \   <span class=\"k\">try</span><span class=\"p\">:</span>\n        <span class=\"n\">home</span>
    <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">os</span><span class=\"o\">.</span><span class=\"n\">environ</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;XDG_HOME&quot;</span><span class=\"p\">])</span>\n
    \   <span class=\"k\">except</span> <span class=\"ne\">KeyError</span><span class=\"p\">:</span>\n
    \       <span class=\"n\">home</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"o\">.</span><span class=\"n\">home</span><span class=\"p\">()</span>\n\n
    \   <span class=\"k\">return</span> <span class=\"p\">[</span>\n        <span
    class=\"p\">{</span><span class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">home</span> <span class=\"o\">/</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;</span><span class=\"si\">{</span><span class=\"n\">tool</span><span
    class=\"si\">}</span><span class=\"s2\">.ini&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;parser&quot;</span><span class=\"p\">:</span> <span class=\"s2\">&quot;ini&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span>
    <span class=\"p\">[</span><span class=\"n\">tool</span><span class=\"p\">]},</span>\n
    \       <span class=\"p\">{</span><span class=\"s2\">&quot;path_specs&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">home</span> <span class=\"o\">/</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;.</span><span class=\"si\">{</span><span
    class=\"n\">tool</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;parser&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;ini&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;keys&quot;</span><span
    class=\"p\">:</span> <span class=\"p\">[</span><span class=\"n\">tool</span><span
    class=\"p\">]},</span>\n        <span class=\"p\">{</span><span class=\"s2\">&quot;path_specs&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">home</span> <span class=\"o\">/</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;.</span><span class=\"si\">{</span><span
    class=\"n\">tool</span><span class=\"si\">}</span><span class=\"s2\">.ini&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;parser&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;ini&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;keys&quot;</span><span
    class=\"p\">:</span> <span class=\"p\">[</span><span class=\"n\">tool</span><span
    class=\"p\">]},</span>\n        <span class=\"p\">{</span>\n            <span
    class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">:</span> <span class=\"n\">home</span>
    <span class=\"o\">/</span> <span class=\"s2\">&quot;.config&quot;</span> <span
    class=\"o\">/</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">tool</span><span class=\"si\">}</span><span
    class=\"s2\">.ini&quot;</span><span class=\"p\">,</span>\n            <span class=\"s2\">&quot;parser&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;ini&quot;</span><span class=\"p\">,</span>\n
    \           <span class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span>
    <span class=\"p\">[</span><span class=\"n\">tool</span><span class=\"p\">],</span>\n
    \       <span class=\"p\">},</span>\n        <span class=\"p\">{</span>\n            <span
    class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">:</span> <span class=\"n\">home</span>
    <span class=\"o\">/</span> <span class=\"s2\">&quot;.config&quot;</span> <span
    class=\"o\">/</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;.</span><span
    class=\"si\">{</span><span class=\"n\">tool</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span><span class=\"p\">,</span>\n            <span class=\"s2\">&quot;parser&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;ini&quot;</span><span class=\"p\">,</span>\n
    \           <span class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span>
    <span class=\"p\">[</span><span class=\"n\">tool</span><span class=\"p\">],</span>\n
    \       <span class=\"p\">},</span>\n        <span class=\"p\">{</span>\n            <span
    class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">:</span> <span class=\"n\">home</span>
    <span class=\"o\">/</span> <span class=\"s2\">&quot;.config&quot;</span> <span
    class=\"o\">/</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;.</span><span
    class=\"si\">{</span><span class=\"n\">tool</span><span class=\"si\">}</span><span
    class=\"s2\">.ini&quot;</span><span class=\"p\">,</span>\n            <span class=\"s2\">&quot;parser&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;ini&quot;</span><span class=\"p\">,</span>\n
    \           <span class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span>
    <span class=\"p\">[</span><span class=\"n\">tool</span><span class=\"p\">],</span>\n
    \       <span class=\"p\">},</span>\n    <span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition function\">\n<p class=\"admonition-title\">Function</p>\n<h2
    id=\"_get_local_path_specs\" class=\"admonition-title\" style=\"margin: 0; padding:
    .5rem 1rem;\">_get_local_path_specs <em class=\"small\">function</em></h2>\n<p>Generate
    a list of standard pathspecs for local, project directory config files.</p>\n<p>Args:\ntool
    (str): name of the tool to configure</p>\n</div>\n<div class=\"admonition source
    is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_get_local_path_specs
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
    class=\"w\"> </span><span class=\"nf\">_get_local_path_specs</span><span class=\"p\">(</span><span
    class=\"n\">tool</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"n\">project_home</span><span class=\"p\">:</span>
    <span class=\"n\">Union</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"n\">Path</span><span class=\"p\">])</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">path_spec_type</span><span class=\"p\">:</span>\n<span
    class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span class=\"sd\">
    \   Generate a list of standard pathspecs for local, project directory config
    files.</span>\n\n<span class=\"sd\">    Args:</span>\n<span class=\"sd\">        tool
    (str): name of the tool to configure</span>\n<span class=\"sd\">    &quot;&quot;&quot;</span>\n
    \   <span class=\"k\">return</span> <span class=\"p\">[</span>\n        <span
    class=\"p\">{</span>\n            <span class=\"s2\">&quot;path_specs&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">project_home</span><span class=\"p\">)</span> <span class=\"o\">/</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">tool</span><span class=\"si\">}</span><span class=\"s2\">.ini&quot;</span><span
    class=\"p\">,</span>\n            <span class=\"s2\">&quot;parser&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;ini&quot;</span><span class=\"p\">,</span>\n
    \           <span class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span>
    <span class=\"p\">[</span><span class=\"n\">tool</span><span class=\"p\">],</span>\n
    \       <span class=\"p\">},</span>\n        <span class=\"p\">{</span>\n            <span
    class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">:</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">project_home</span><span class=\"p\">)</span>
    <span class=\"o\">/</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;.</span><span
    class=\"si\">{</span><span class=\"n\">tool</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span><span class=\"p\">,</span>\n            <span class=\"s2\">&quot;parser&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;ini&quot;</span><span class=\"p\">,</span>\n
    \           <span class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span>
    <span class=\"p\">[</span><span class=\"n\">tool</span><span class=\"p\">],</span>\n
    \       <span class=\"p\">},</span>\n        <span class=\"p\">{</span>\n            <span
    class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">:</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">project_home</span><span class=\"p\">)</span>
    <span class=\"o\">/</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;.</span><span
    class=\"si\">{</span><span class=\"n\">tool</span><span class=\"si\">}</span><span
    class=\"s2\">.ini&quot;</span><span class=\"p\">,</span>\n            <span class=\"s2\">&quot;parser&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;ini&quot;</span><span class=\"p\">,</span>\n
    \           <span class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span>
    <span class=\"p\">[</span><span class=\"n\">tool</span><span class=\"p\">],</span>\n
    \       <span class=\"p\">},</span>\n        <span class=\"p\">{</span>\n            <span
    class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">:</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">project_home</span><span class=\"p\">)</span>
    <span class=\"o\">/</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">tool</span><span class=\"si\">}</span><span
    class=\"s2\">.yml&quot;</span><span class=\"p\">,</span>\n            <span class=\"s2\">&quot;parser&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;yaml&quot;</span><span class=\"p\">,</span>\n
    \           <span class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span>
    <span class=\"p\">[</span><span class=\"n\">tool</span><span class=\"p\">],</span>\n
    \       <span class=\"p\">},</span>\n        <span class=\"p\">{</span>\n            <span
    class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">:</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">project_home</span><span class=\"p\">)</span>
    <span class=\"o\">/</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;.</span><span
    class=\"si\">{</span><span class=\"n\">tool</span><span class=\"si\">}</span><span
    class=\"s2\">.yml&quot;</span><span class=\"p\">,</span>\n            <span class=\"s2\">&quot;parser&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;yaml&quot;</span><span class=\"p\">,</span>\n
    \           <span class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span>
    <span class=\"p\">[</span><span class=\"n\">tool</span><span class=\"p\">],</span>\n
    \       <span class=\"p\">},</span>\n        <span class=\"p\">{</span>\n            <span
    class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">:</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">project_home</span><span class=\"p\">)</span>
    <span class=\"o\">/</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">tool</span><span class=\"si\">}</span><span
    class=\"s2\">.toml&quot;</span><span class=\"p\">,</span>\n            <span class=\"s2\">&quot;parser&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;toml&quot;</span><span class=\"p\">,</span>\n
    \           <span class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span>
    <span class=\"p\">[</span><span class=\"n\">tool</span><span class=\"p\">],</span>\n
    \       <span class=\"p\">},</span>\n        <span class=\"p\">{</span>\n            <span
    class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">:</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">project_home</span><span class=\"p\">)</span>
    <span class=\"o\">/</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;.</span><span
    class=\"si\">{</span><span class=\"n\">tool</span><span class=\"si\">}</span><span
    class=\"s2\">.toml&quot;</span><span class=\"p\">,</span>\n            <span class=\"s2\">&quot;parser&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;toml&quot;</span><span class=\"p\">,</span>\n
    \           <span class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span>
    <span class=\"p\">[</span><span class=\"n\">tool</span><span class=\"p\">],</span>\n
    \       <span class=\"p\">},</span>\n        <span class=\"p\">{</span>\n            <span
    class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">:</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">project_home</span><span class=\"p\">)</span>
    <span class=\"o\">/</span> <span class=\"s2\">&quot;pyproject.toml&quot;</span><span
    class=\"p\">,</span>\n            <span class=\"s2\">&quot;parser&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;toml&quot;</span><span class=\"p\">,</span>\n
    \           <span class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span>
    <span class=\"p\">[</span><span class=\"s2\">&quot;tool&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">tool</span><span class=\"p\">],</span>\n        <span class=\"p\">},</span>\n
    \       <span class=\"p\">{</span>\n            <span class=\"s2\">&quot;path_specs&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">project_home</span><span class=\"p\">)</span> <span class=\"o\">/</span>
    <span class=\"s2\">&quot;setup.cfg&quot;</span><span class=\"p\">,</span>\n            <span
    class=\"s2\">&quot;parser&quot;</span><span class=\"p\">:</span> <span class=\"s2\">&quot;ini&quot;</span><span
    class=\"p\">,</span>\n            <span class=\"s2\">&quot;keys&quot;</span><span
    class=\"p\">:</span> <span class=\"p\">[</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;tool.</span><span class=\"si\">{</span><span class=\"n\">tool</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">],</span>\n
    \       <span class=\"p\">},</span>\n    <span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition function\">\n<p class=\"admonition-title\">Function</p>\n<h2
    id=\"_get_attrs\" class=\"admonition-title\" style=\"margin: 0; padding: .5rem
    1rem;\">_get_attrs <em class=\"small\">function</em></h2>\n<p>Get nested config
    data from a list of keys.</p>\n<p>specifically written for pyproject.toml which
    needs to get <code>tool</code> then <code>&lt;tool&gt;</code></p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_get_attrs
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
    class=\"w\"> </span><span class=\"nf\">_get_attrs</span><span class=\"p\">(</span><span
    class=\"n\">attrs</span><span class=\"p\">:</span> <span class=\"nb\">list</span><span
    class=\"p\">,</span> <span class=\"n\">config</span><span class=\"p\">:</span>
    <span class=\"n\">Dict</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Dict</span><span class=\"p\">:</span>\n<span class=\"w\">    </span><span
    class=\"sd\">&quot;&quot;&quot;Get nested config data from a list of keys.</span>\n\n<span
    class=\"sd\">    specifically written for pyproject.toml which needs to get `tool`
    then `&lt;tool&gt;`</span>\n<span class=\"sd\">    &quot;&quot;&quot;</span>\n
    \   <span class=\"k\">for</span> <span class=\"n\">attr</span> <span class=\"ow\">in</span>
    <span class=\"n\">attrs</span><span class=\"p\">:</span>\n        <span class=\"n\">config</span>
    <span class=\"o\">=</span> <span class=\"n\">config</span><span class=\"p\">[</span><span
    class=\"n\">attr</span><span class=\"p\">]</span>\n    <span class=\"k\">return</span>
    <span class=\"n\">config</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div class=\"admonition
    function\">\n<p class=\"admonition-title\">Function</p>\n<h2 id=\"_load_config_file\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">_load_config_file
    <em class=\"small\">function</em></h2>\n<p>Load a configuration file using the
    appropriate parser.</p>\n<p>Args:\nfile_spec: Dictionary containing path_specs,
    parser, and keys information</p>\n<p>Returns:\nOptional[Dict[str, Any]]: Parsed
    configuration or None if file doesn't exist</p>\n</div>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_load_config_file
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
    class=\"w\"> </span><span class=\"nf\">_load_config_file</span><span class=\"p\">(</span><span
    class=\"n\">file_spec</span><span class=\"p\">:</span> <span class=\"n\">Dict</span><span
    class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">,</span> <span
    class=\"n\">Any</span><span class=\"p\">])</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"n\">Dict</span><span
    class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">,</span> <span
    class=\"n\">Any</span><span class=\"p\">]]:</span>\n<span class=\"w\">    </span><span
    class=\"sd\">&quot;&quot;&quot;Load a configuration file using the appropriate
    parser.</span>\n\n<span class=\"sd\">    Args:</span>\n<span class=\"sd\">        file_spec:
    Dictionary containing path_specs, parser, and keys information</span>\n\n<span
    class=\"sd\">    Returns:</span>\n<span class=\"sd\">        Optional[Dict[str,
    Any]]: Parsed configuration or None if file doesn&#39;t exist</span>\n<span class=\"sd\">
    \   &quot;&quot;&quot;</span>\n    <span class=\"n\">path</span> <span class=\"o\">=</span>
    <span class=\"n\">file_spec</span><span class=\"p\">[</span><span class=\"s2\">&quot;path_specs&quot;</span><span
    class=\"p\">]</span>\n    <span class=\"k\">if</span> <span class=\"ow\">not</span>
    <span class=\"n\">path</span><span class=\"o\">.</span><span class=\"n\">exists</span><span
    class=\"p\">():</span>\n        <span class=\"k\">return</span> <span class=\"kc\">None</span>\n\n
    \   <span class=\"k\">try</span><span class=\"p\">:</span>\n        <span class=\"k\">if</span>
    <span class=\"n\">file_spec</span><span class=\"p\">[</span><span class=\"s2\">&quot;parser&quot;</span><span
    class=\"p\">]</span> <span class=\"o\">==</span> <span class=\"s2\">&quot;toml&quot;</span><span
    class=\"p\">:</span>\n            <span class=\"k\">with</span> <span class=\"nb\">open</span><span
    class=\"p\">(</span><span class=\"n\">path</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;rb&quot;</span><span class=\"p\">)</span> <span class=\"k\">as</span>
    <span class=\"n\">f</span><span class=\"p\">:</span>\n                <span class=\"n\">config</span>
    <span class=\"o\">=</span> <span class=\"n\">tomli</span><span class=\"o\">.</span><span
    class=\"n\">load</span><span class=\"p\">(</span><span class=\"n\">f</span><span
    class=\"p\">)</span>\n        <span class=\"k\">elif</span> <span class=\"n\">file_spec</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;parser&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">==</span> <span class=\"s2\">&quot;yaml&quot;</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">with</span> <span class=\"nb\">open</span><span
    class=\"p\">(</span><span class=\"n\">path</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;r&quot;</span><span class=\"p\">)</span> <span class=\"k\">as</span>
    <span class=\"n\">f</span><span class=\"p\">:</span>\n                <span class=\"n\">config</span>
    <span class=\"o\">=</span> <span class=\"n\">yaml</span><span class=\"o\">.</span><span
    class=\"n\">safe_load</span><span class=\"p\">(</span><span class=\"n\">f</span><span
    class=\"p\">)</span>\n        <span class=\"k\">elif</span> <span class=\"n\">file_spec</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;parser&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">==</span> <span class=\"s2\">&quot;ini&quot;</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">config</span> <span class=\"o\">=</span> <span class=\"n\">configparser</span><span
    class=\"o\">.</span><span class=\"n\">ConfigParser</span><span class=\"p\">()</span>\n
    \           <span class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">read</span><span
    class=\"p\">(</span><span class=\"n\">path</span><span class=\"p\">)</span>\n
    \           <span class=\"c1\"># Convert ConfigParser to dict</span>\n            <span
    class=\"n\">config</span> <span class=\"o\">=</span> <span class=\"p\">{</span><span
    class=\"n\">s</span><span class=\"p\">:</span> <span class=\"nb\">dict</span><span
    class=\"p\">(</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">items</span><span class=\"p\">(</span><span class=\"n\">s</span><span
    class=\"p\">))</span> <span class=\"k\">for</span> <span class=\"n\">s</span>
    <span class=\"ow\">in</span> <span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">sections</span><span class=\"p\">()}</span>\n        <span class=\"k\">else</span><span
    class=\"p\">:</span>\n            <span class=\"k\">return</span> <span class=\"kc\">None</span>\n\n
    \       <span class=\"k\">return</span> <span class=\"n\">_get_attrs</span><span
    class=\"p\">(</span><span class=\"n\">file_spec</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;keys&quot;</span><span class=\"p\">],</span> <span class=\"n\">config</span><span
    class=\"p\">)</span>\n    <span class=\"k\">except</span> <span class=\"p\">(</span>\n
    \       <span class=\"ne\">KeyError</span><span class=\"p\">,</span>\n        <span
    class=\"ne\">TypeError</span><span class=\"p\">,</span>\n        <span class=\"n\">yaml</span><span
    class=\"o\">.</span><span class=\"n\">YAMLError</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">tomli</span><span class=\"o\">.</span><span class=\"n\">TOMLDecodeError</span><span
    class=\"p\">,</span>\n        <span class=\"n\">configparser</span><span class=\"o\">.</span><span
    class=\"n\">Error</span><span class=\"p\">,</span>\n    <span class=\"p\">):</span>\n
    \       <span class=\"k\">return</span> <span class=\"kc\">None</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition function\">\n<p class=\"admonition-title\">Function</p>\n<h2
    id=\"_load_files\" class=\"admonition-title\" style=\"margin: 0; padding: .5rem
    1rem;\">_load_files <em class=\"small\">function</em></h2>\n<p>Load config files
    stopping at the first one that exists and can be parsed.</p>\n<p>Args:\nconfig_path_specs:
    List of path specifications to try</p>\n<p>Returns:\nDict[str, Any]: Configuration
    dictionary</p>\n</div>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">_load_files <em class='small'>source</em></p>\n<pre
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
    class=\"w\"> </span><span class=\"nf\">_load_files</span><span class=\"p\">(</span><span
    class=\"n\">config_path_specs</span><span class=\"p\">:</span> <span class=\"n\">path_spec_type</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Dict</span><span
    class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">,</span> <span
    class=\"n\">Any</span><span class=\"p\">]:</span>\n<span class=\"w\">    </span><span
    class=\"sd\">&quot;&quot;&quot;Load config files stopping at the first one that
    exists and can be parsed.</span>\n\n<span class=\"sd\">    Args:</span>\n<span
    class=\"sd\">        config_path_specs: List of path specifications to try</span>\n\n<span
    class=\"sd\">    Returns:</span>\n<span class=\"sd\">        Dict[str, Any]: Configuration
    dictionary</span>\n<span class=\"sd\">    &quot;&quot;&quot;</span>\n    <span
    class=\"k\">for</span> <span class=\"n\">file_spec</span> <span class=\"ow\">in</span>
    <span class=\"n\">config_path_specs</span><span class=\"p\">:</span>\n        <span
    class=\"n\">config</span> <span class=\"o\">=</span> <span class=\"n\">_load_config_file</span><span
    class=\"p\">(</span><span class=\"n\">file_spec</span><span class=\"p\">)</span>\n
    \       <span class=\"k\">if</span> <span class=\"n\">config</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">return</span> <span class=\"n\">config</span>\n
    \   <span class=\"k\">return</span> <span class=\"p\">{}</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition function\">\n<p class=\"admonition-title\">Function</p>\n<h2
    id=\"_load_env\" class=\"admonition-title\" style=\"margin: 0; padding: .5rem
    1rem;\">_load_env <em class=\"small\">function</em></h2>\n<p>Load config from
    environment variables.</p>\n<p>Args:\ntool (str): name of the tool to configure</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_load_env
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
    class=\"w\"> </span><span class=\"nf\">_load_env</span><span class=\"p\">(</span><span
    class=\"n\">tool</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Dict</span><span
    class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">,</span> <span
    class=\"n\">Any</span><span class=\"p\">]:</span>\n<span class=\"w\">    </span><span
    class=\"sd\">&quot;&quot;&quot;Load config from environment variables.</span>\n\n<span
    class=\"sd\">    Args:</span>\n<span class=\"sd\">        tool (str): name of
    the tool to configure</span>\n<span class=\"sd\">    &quot;&quot;&quot;</span>\n
    \   <span class=\"n\">env_prefix</span> <span class=\"o\">=</span> <span class=\"n\">tool</span><span
    class=\"o\">.</span><span class=\"n\">upper</span><span class=\"p\">()</span>\n
    \   <span class=\"n\">env_config</span> <span class=\"o\">=</span> <span class=\"p\">{</span>\n
    \       <span class=\"n\">key</span><span class=\"o\">.</span><span class=\"n\">replace</span><span
    class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">env_prefix</span><span class=\"si\">}</span><span
    class=\"s2\">_&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;&quot;</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">lower</span><span
    class=\"p\">():</span> <span class=\"n\">value</span>\n        <span class=\"k\">for</span>
    <span class=\"n\">key</span><span class=\"p\">,</span> <span class=\"n\">value</span>
    <span class=\"ow\">in</span> <span class=\"n\">os</span><span class=\"o\">.</span><span
    class=\"n\">environ</span><span class=\"o\">.</span><span class=\"n\">items</span><span
    class=\"p\">()</span>\n        <span class=\"k\">if</span> <span class=\"n\">key</span><span
    class=\"o\">.</span><span class=\"n\">startswith</span><span class=\"p\">(</span><span
    class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">env_prefix</span><span class=\"si\">}</span><span class=\"s2\">_&quot;</span><span
    class=\"p\">)</span>\n    <span class=\"p\">}</span>\n    <span class=\"k\">return</span>
    <span class=\"n\">env_config</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div class=\"admonition
    function\">\n<p class=\"admonition-title\">Function</p>\n<h2 id=\"load\" class=\"admonition-title\"
    style=\"margin: 0; padding: .5rem 1rem;\">load <em class=\"small\">function</em></h2>\n<p>Load
    tool config from standard config files.</p>\n<p>Resolution Order</p>\n<ul>\n<li>First
    global file with a tool key</li>\n<li>First local file with a tool key</li>\n<li>Environment
    variables prefixed with <code>TOOL</code></li>\n<li>Overrides</li>\n</ul>\n<p>Args:\ntool
    (str): name of the tool to configure\nproject_home (Union[Path, str], optional):
    Project directory to search for config files. Defaults to &quot;.&quot;.\noverrides
    (Dict, optional): Override values to apply last. Defaults to None.</p>\n<p>Returns:\nDict[str,
    Any]: Configuration object</p>\n</div>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">load <em class='small'>source</em></p>\n<pre
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
    class=\"w\"> </span><span class=\"nf\">load</span><span class=\"p\">(</span>\n
    \   <span class=\"n\">tool</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">,</span>\n    <span class=\"n\">project_home</span><span class=\"p\">:</span>
    <span class=\"n\">Union</span><span class=\"p\">[</span><span class=\"n\">Path</span><span
    class=\"p\">,</span> <span class=\"nb\">str</span><span class=\"p\">]</span> <span
    class=\"o\">=</span> <span class=\"s2\">&quot;.&quot;</span><span class=\"p\">,</span>\n
    \   <span class=\"n\">overrides</span><span class=\"p\">:</span> <span class=\"n\">Optional</span><span
    class=\"p\">[</span><span class=\"n\">Dict</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">,</span> <span class=\"n\">Any</span><span
    class=\"p\">]]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span><span
    class=\"p\">,</span>\n<span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Dict</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"n\">Any</span><span class=\"p\">]:</span>\n<span
    class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Load tool config from
    standard config files.</span>\n\n<span class=\"sd\">    Resolution Order</span>\n\n<span
    class=\"sd\">    * First global file with a tool key</span>\n<span class=\"sd\">
    \   * First local file with a tool key</span>\n<span class=\"sd\">    * Environment
    variables prefixed with `TOOL`</span>\n<span class=\"sd\">    * Overrides</span>\n\n<span
    class=\"sd\">    Args:</span>\n<span class=\"sd\">        tool (str): name of
    the tool to configure</span>\n<span class=\"sd\">        project_home (Union[Path,
    str], optional): Project directory to search for config files. Defaults to &quot;.&quot;.</span>\n<span
    class=\"sd\">        overrides (Dict, optional): Override values to apply last.
    Defaults to None.</span>\n\n<span class=\"sd\">    Returns:</span>\n<span class=\"sd\">
    \       Dict[str, Any]: Configuration object</span>\n<span class=\"sd\">    &quot;&quot;&quot;</span>\n
    \   <span class=\"n\">overrides</span> <span class=\"o\">=</span> <span class=\"n\">overrides</span>
    <span class=\"ow\">or</span> <span class=\"p\">{}</span>\n    <span class=\"n\">config</span>
    <span class=\"o\">=</span> <span class=\"p\">{}</span>\n\n    <span class=\"c1\">#
    Load from files in order of precedence</span>\n    <span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">update</span><span class=\"p\">(</span><span
    class=\"n\">_load_files</span><span class=\"p\">(</span><span class=\"n\">_get_global_path_specs</span><span
    class=\"p\">(</span><span class=\"n\">tool</span><span class=\"p\">))</span> <span
    class=\"ow\">or</span> <span class=\"p\">{})</span>\n    <span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">update</span><span class=\"p\">(</span><span
    class=\"n\">_load_files</span><span class=\"p\">(</span><span class=\"n\">_get_local_path_specs</span><span
    class=\"p\">(</span><span class=\"n\">tool</span><span class=\"p\">,</span> <span
    class=\"n\">project_home</span><span class=\"p\">))</span> <span class=\"ow\">or</span>
    <span class=\"p\">{})</span>\n    <span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">update</span><span class=\"p\">(</span><span class=\"n\">_load_env</span><span
    class=\"p\">(</span><span class=\"n\">tool</span><span class=\"p\">))</span>\n
    \   <span class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">update</span><span
    class=\"p\">(</span><span class=\"n\">overrides</span><span class=\"p\">)</span>\n\n
    \   <span class=\"c1\"># If no settings class is provided, return the raw dict</span>\n
    \   <span class=\"k\">return</span> <span class=\"n\">config</span>\n</pre></div>\n\n</pre>\n\n</div>\n\n
    \   </section>\n</article>    </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>standard_config.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Standard Config. A module to load tooling
    config from a users project space. Inspired from frustrations that some tools
    have a tool.ini, .tool.ini, setup.cfg,\u2026\" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n
    \       <meta property=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    <head>\n        <script src=\"https://unpkg.com/@tailwindcss/browser@4\"></script>\n<title>standard_config.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Standard Config. A module to load tooling
    config from a users project space. Inspired from frustrations that some tools
    have a tool.ini, .tool.ini, setup.cfg,\u2026\" />\n <link href=\"/favicon.ico\"
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
    class=\"title\">\n    <h1 id=\"title\">\n        standard_config.py\n    </h1>\n</section></article>\n
    \    </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       standard_config.py\n    </h1>\n</section>    <section class=\"body\">\n
    \       <hr />\n<p>Standard Config.\nA module to load tooling config from a users
    project space.</p>\n<p>Inspired from frustrations that some tools have a tool.ini,
    .tool.ini,\nsetup.cfg, or pyproject.toml.  Some allow for global configs, some
    don't.  Some\nproperly follow the users home directory, others end up in a weird
    temp\ndirectory.  Windows home directory is only more confusing.  Some will even\nrespect
    the users <code>$XDG_HOME</code> directory.</p>\n<p>This file is for any project
    that can be configured in plain text such as <code>ini</code>\nor <code>toml</code>
    and not requiring a .py file.  Just name your tool and let users put\nconfig where
    it makes sense to them, no need to figure out resolution order.</p>\n<h2 id=\"usage\">Usage:
    <a class=\"header-anchor\" href=\"#usage\"><svg class=\"heading-permalink\" aria-hidden=\"true\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"kn\">from</span><span
    class=\"w\"> </span><span class=\"nn\">standard_config</span><span class=\"w\">
    </span><span class=\"kn\">import</span> <span class=\"n\">load</span>\n\n<span
    class=\"c1\"># Retrieve any overrides from the user</span>\n<span class=\"n\">overrides</span>
    <span class=\"o\">=</span> <span class=\"p\">{</span><span class=\"s1\">&#39;setting&#39;</span><span
    class=\"p\">:</span> <span class=\"kc\">True</span><span class=\"p\">}</span>\n<span
    class=\"n\">config</span> <span class=\"o\">=</span> <span class=\"n\">load</span><span
    class=\"p\">(</span><span class=\"s1\">&#39;my_tool&#39;</span><span class=\"p\">,</span>
    <span class=\"n\">overrides</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<h2
    id=\"resolution-order\">Resolution Order <a class=\"header-anchor\" href=\"#resolution-order\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<ul>\n<li>First global
    file with a tool key</li>\n<li>First local file with a tool key</li>\n<li>Environment
    variables prefixed with <code>TOOL</code></li>\n<li>Overrides</li>\n</ul>\n<h3>Tool
    Specific Ini files</h3>\n<p>Ini file formats must include a <code>&lt;tool&gt;</code>
    key.</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button class='copy'
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[my_tool]</span>\n<span
    class=\"na\">setting</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s\">True</span>\n</pre></div>\n\n</pre>\n\n<h3>pyproject.toml</h3>\n<p>Toml
    files must include a <code>tool.&lt;tool&gt;</code> key</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[tool.my_tool]</span>\n<span
    class=\"n\">setting</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"err\">True</span>\n</pre></div>\n\n</pre>\n\n<h3>setup.cfg</h3>\n<p>setup.cfg
    files must include a <code>tool:&lt;tool&gt;</code> key</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[tool:my_tool]</span>\n<span
    class=\"na\">setting</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s\">True</span>\n</pre></div>\n\n</pre>\n\n<h3>global
    files to consider</h3>\n<ul>\n<li><home>/tool.ini</li>\n<li><home>/.tool</li>\n<li><home>/.tool.ini</li>\n<li><home>/.config/tool.ini</li>\n<li><home>/.config/.tool</li>\n<li><home>/.config/.tool.ini</li>\n</ul>\n<h3>local
    files to consider</h3>\n<ul>\n<li>&lt;project_home&gt;/tool.ini</li>\n<li>&lt;project_home&gt;/.tool</li>\n<li>&lt;project_home&gt;/.tool.ini</li>\n<li>&lt;project_home&gt;/pyproject.toml</li>\n<li>&lt;project_home&gt;/setup.cfg</li>\n</ul>\n<p>Markata's
    standard configuration system.</p>\n<h2 id=\"configuration-overview\">Configuration
    Overview <a class=\"header-anchor\" href=\"#configuration-overview\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Markata uses a hierarchical
    configuration system based on Pydantic models. Configuration\ncan be set through:</p>\n<ol>\n<li>TOML
    files</li>\n<li>Environment variables</li>\n<li>Command line arguments</li>\n</ol>\n<h1
    id=\"basic-configuration\">Basic Configuration <a class=\"header-anchor\" href=\"#basic-configuration\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>Minimal <code>markata.toml</code>:</p>\n<pre
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
    class=\"c1\"># Site info</span>\n<span class=\"n\">title</span><span class=\"w\">
    </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s2\">&quot;My
    Site&quot;</span>\n<span class=\"n\">url</span><span class=\"w\"> </span><span
    class=\"o\">=</span><span class=\"w\"> </span><span class=\"s2\">&quot;https://example.com&quot;</span>\n<span
    class=\"n\">description</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;Site description&quot;</span>\n\n<span
    class=\"c1\"># Content locations</span>\n<span class=\"n\">content_dir</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;content&quot;</span>\n<span class=\"n\">output_dir</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;markout&quot;</span>\n<span class=\"n\">assets_dir</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;static&quot;</span>\n\n<span class=\"c1\"># Plugin management</span>\n<span
    class=\"n\">hooks</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"p\">[</span><span class=\"s2\">&quot;default&quot;</span><span
    class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<h1 id=\"environment-variables\">Environment
    Variables <a class=\"header-anchor\" href=\"#environment-variables\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>All settings can be
    overridden with environment variables:</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"c1\"># Override
    site URL</span>\n<span class=\"nb\">export</span><span class=\"w\"> </span><span
    class=\"nv\">MARKATA_URL</span><span class=\"o\">=</span><span class=\"s2\">&quot;https://staging.example.com&quot;</span>\n\n<span
    class=\"c1\"># Override output directory</span>\n<span class=\"nb\">export</span><span
    class=\"w\"> </span><span class=\"nv\">MARKATA_OUTPUT_DIR</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;dist&quot;</span>\n\n<span class=\"c1\"># Enable debug mode</span>\n<span
    class=\"nb\">export</span><span class=\"w\"> </span><span class=\"nv\">MARKATA_DEBUG</span><span
    class=\"o\">=</span><span class=\"m\">1</span>\n</pre></div>\n\n</pre>\n\n<h1
    id=\"detailed-configuration\">Detailed Configuration <a class=\"header-anchor\"
    href=\"#detailed-configuration\"><svg class=\"heading-permalink\" aria-hidden=\"true\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<h2 id=\"core-settings\">Core
    Settings <a class=\"header-anchor\" href=\"#core-settings\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n<span
    class=\"c1\"># Site information</span>\n<span class=\"n\">title</span><span class=\"w\">
    </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s2\">&quot;My
    Site&quot;</span><span class=\"w\">                  </span><span class=\"c1\">#
    Site title</span>\n<span class=\"n\">url</span><span class=\"w\"> </span><span
    class=\"o\">=</span><span class=\"w\"> </span><span class=\"s2\">&quot;https://example.com&quot;</span><span
    class=\"w\">        </span><span class=\"c1\"># Base URL</span>\n<span class=\"n\">description</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;Site description&quot;</span><span class=\"w\">   </span><span
    class=\"c1\"># Meta description</span>\n<span class=\"n\">author_name</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;Author Name&quot;</span><span class=\"w\">        </span><span
    class=\"c1\"># Author name</span>\n<span class=\"n\">author_email</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;me@example.com&quot;</span><span class=\"w\">    </span><span
    class=\"c1\"># Author email</span>\n<span class=\"n\">icon</span><span class=\"w\">
    </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s2\">&quot;favicon.ico&quot;</span><span
    class=\"w\">               </span><span class=\"c1\"># Site icon</span>\n<span
    class=\"n\">lang</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;en&quot;</span><span class=\"w\">
    \                       </span><span class=\"c1\"># Site language</span>\n\n<span
    class=\"c1\"># Content locations</span>\n<span class=\"n\">content_dir</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;content&quot;</span><span class=\"w\">           </span><span
    class=\"c1\"># Source content location</span>\n<span class=\"n\">output_dir</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;markout&quot;</span><span class=\"w\">            </span><span
    class=\"c1\"># Build output location</span>\n<span class=\"n\">assets_dir</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;static&quot;</span><span class=\"w\">             </span><span
    class=\"c1\"># Static assets location</span>\n<span class=\"n\">template_dir</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;templates&quot;</span><span class=\"w\">        </span><span
    class=\"c1\"># Template location</span>\n\n<span class=\"c1\"># Plugin management</span>\n<span
    class=\"n\">hooks</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"p\">[</span><span class=\"s2\">&quot;default&quot;</span><span
    class=\"p\">]</span><span class=\"w\">               </span><span class=\"c1\">#
    Active plugins</span>\n<span class=\"n\">disabled_hooks</span><span class=\"w\">
    </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"p\">[]</span><span
    class=\"w\">               </span><span class=\"c1\"># Disabled plugins</span>\n</pre></div>\n\n</pre>\n\n<h2
    id=\"cache-settings\">Cache Settings <a class=\"header-anchor\" href=\"#cache-settings\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n<span
    class=\"c1\"># Cache configuration</span>\n<span class=\"n\">default_cache_expire</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"mi\">3600</span><span class=\"w\">       </span><span class=\"c1\"># Default
    TTL (1 hour)</span>\n<span class=\"n\">template_cache_expire</span><span class=\"w\">
    </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"mi\">86400</span><span
    class=\"w\">     </span><span class=\"c1\"># Template TTL (24 hours)</span>\n<span
    class=\"n\">markdown_cache_expire</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"mi\">21600</span><span class=\"w\">     </span><span
    class=\"c1\"># Markdown TTL (6 hours)</span>\n<span class=\"n\">dynamic_cache_expire</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"mi\">3600</span><span class=\"w\">       </span><span class=\"c1\"># Dynamic
    TTL (1 hour)</span>\n</pre></div>\n\n</pre>\n\n<h2 id=\"development-settings\">Development
    Settings <a class=\"header-anchor\" href=\"#development-settings\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n<span
    class=\"c1\"># Development server</span>\n<span class=\"n\">dev_server_port</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"mi\">8000</span><span class=\"w\">            </span><span class=\"c1\">#
    Local server port</span>\n<span class=\"n\">dev_server_host</span><span class=\"w\">
    </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s2\">&quot;localhost&quot;</span><span
    class=\"w\">     </span><span class=\"c1\"># Local server host</span>\n<span class=\"n\">debug</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"kc\">false</span><span class=\"w\">                     </span><span class=\"c1\">#
    Debug mode</span>\n\n<span class=\"c1\"># Performance</span>\n<span class=\"n\">parallel</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"kc\">true</span><span class=\"w\">                   </span><span class=\"c1\">#
    Enable parallel processing</span>\n<span class=\"n\">workers</span><span class=\"w\">
    </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"mi\">4</span><span
    class=\"w\">                       </span><span class=\"c1\"># Number of worker
    threads</span>\n</pre></div>\n\n</pre>\n\n<h2 id=\"content-settings\">Content
    Settings <a class=\"header-anchor\" href=\"#content-settings\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n<span
    class=\"c1\"># Content processing</span>\n<span class=\"n\">default_template</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;post.html&quot;</span><span class=\"w\">    </span><span class=\"c1\">#
    Default template</span>\n<span class=\"n\">markdown_extensions</span><span class=\"w\">
    </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"p\">[</span><span
    class=\"w\">           </span><span class=\"c1\"># Markdown extensions</span>\n<span
    class=\"w\">    </span><span class=\"s2\">&quot;fenced_code&quot;</span><span
    class=\"p\">,</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;tables&quot;</span><span
    class=\"p\">,</span>\n<span class=\"w\">    </span><span class=\"s2\">&quot;footnotes&quot;</span>\n<span
    class=\"p\">]</span>\n\n<span class=\"c1\"># Content filtering</span>\n<span class=\"n\">draft</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"kc\">false</span><span class=\"w\">                     </span><span class=\"c1\">#
    Include drafts</span>\n<span class=\"n\">future</span><span class=\"w\"> </span><span
    class=\"o\">=</span><span class=\"w\"> </span><span class=\"kc\">false</span><span
    class=\"w\">                    </span><span class=\"c1\"># Include future posts</span>\n</pre></div>\n\n</pre>\n\n<h1
    id=\"plugin-configuration\">Plugin Configuration <a class=\"header-anchor\" href=\"#plugin-configuration\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>Each plugin can define
    its own configuration section:</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"c1\"># RSS
    feed configuration</span>\n<span class=\"k\">[markata.feeds]</span>\n<span class=\"n\">rss</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"p\">{</span><span class=\"w\"> </span><span class=\"n\">output</span><span
    class=\"w\"> </span><span class=\"p\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;rss.xml&quot;</span><span class=\"w\"> </span><span class=\"p\">}</span>\n<span
    class=\"n\">atom</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"p\">{</span><span class=\"w\"> </span><span
    class=\"n\">output</span><span class=\"w\"> </span><span class=\"p\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;atom.xml&quot;</span><span class=\"w\">
    </span><span class=\"p\">}</span>\n<span class=\"n\">json</span><span class=\"w\">
    </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"p\">{</span><span
    class=\"w\"> </span><span class=\"n\">output</span><span class=\"w\"> </span><span
    class=\"p\">=</span><span class=\"w\"> </span><span class=\"s2\">&quot;feed.json&quot;</span><span
    class=\"w\"> </span><span class=\"p\">}</span>\n\n<span class=\"c1\"># Template
    configuration</span>\n<span class=\"k\">[markata.template]</span>\n<span class=\"n\">engine</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;jinja2&quot;</span>\n<span class=\"n\">cache_size</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"mi\">100</span>\n<span class=\"n\">autoescape</span><span class=\"w\">
    </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"kc\">true</span>\n\n<span
    class=\"c1\"># Markdown configuration</span>\n<span class=\"k\">[markata.markdown]</span>\n<span
    class=\"n\">highlight_theme</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;monokai&quot;</span>\n<span class=\"n\">line_numbers</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"kc\">true</span>\n</pre></div>\n\n</pre>\n\n<h2 id=\"configuration-validation\">Configuration
    Validation <a class=\"header-anchor\" href=\"#configuration-validation\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>The configuration is
    validated using Pydantic models:</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    class=\"w\"> </span><span class=\"nn\">pydantic</span><span class=\"w\"> </span><span
    class=\"kn\">import</span> <span class=\"n\">BaseModel</span><span class=\"p\">,</span>
    <span class=\"n\">Field</span>\n\n<span class=\"k\">class</span><span class=\"w\">
    </span><span class=\"nc\">MarkataConfig</span><span class=\"p\">(</span><span
    class=\"n\">BaseModel</span><span class=\"p\">):</span>\n<span class=\"w\">    </span><span
    class=\"sd\">&quot;&quot;&quot;Core configuration model.&quot;&quot;&quot;</span>\n
    \   <span class=\"c1\"># Site info</span>\n    <span class=\"n\">title</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span> <span class=\"o\">=</span>
    <span class=\"n\">Field</span><span class=\"p\">(</span><span class=\"o\">...</span><span
    class=\"p\">,</span> <span class=\"n\">description</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;Site title&quot;</span><span class=\"p\">)</span>\n    <span
    class=\"n\">url</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"n\">Field</span><span class=\"p\">(</span><span
    class=\"o\">...</span><span class=\"p\">,</span> <span class=\"n\">description</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;Site base URL&quot;</span><span class=\"p\">)</span>\n\n
    \   <span class=\"c1\"># Directories</span>\n    <span class=\"n\">content_dir</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span> <span class=\"o\">=</span>
    <span class=\"n\">Field</span><span class=\"p\">(</span><span class=\"s2\">&quot;content&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">description</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;Content directory&quot;</span><span class=\"p\">)</span>\n
    \   <span class=\"n\">output_dir</span><span class=\"p\">:</span> <span class=\"n\">Path</span>
    <span class=\"o\">=</span> <span class=\"n\">Field</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;markout&quot;</span><span class=\"p\">,</span> <span class=\"n\">description</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;Output directory&quot;</span><span
    class=\"p\">)</span>\n\n    <span class=\"c1\"># Features</span>\n    <span class=\"n\">debug</span><span
    class=\"p\">:</span> <span class=\"nb\">bool</span> <span class=\"o\">=</span>
    <span class=\"n\">Field</span><span class=\"p\">(</span><span class=\"kc\">False</span><span
    class=\"p\">,</span> <span class=\"n\">description</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;Enable debug mode&quot;</span><span class=\"p\">)</span>\n
    \   <span class=\"n\">parallel</span><span class=\"p\">:</span> <span class=\"nb\">bool</span>
    <span class=\"o\">=</span> <span class=\"n\">Field</span><span class=\"p\">(</span><span
    class=\"kc\">True</span><span class=\"p\">,</span> <span class=\"n\">description</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;Enable parallel processing&quot;</span><span
    class=\"p\">)</span>\n\n    <span class=\"n\">model_config</span> <span class=\"o\">=</span>
    <span class=\"n\">ConfigDict</span><span class=\"p\">(</span>\n        <span class=\"n\">validate_assignment</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">arbitrary_types_allowed</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n        <span class=\"n\">extra</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;allow&quot;</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">str_strip_whitespace</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n        <span class=\"n\">validate_default</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">populate_by_name</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n    <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<h1
    id=\"usage-example\">Usage Example <a class=\"header-anchor\" href=\"#usage-example\"><svg
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"kn\">from</span><span
    class=\"w\"> </span><span class=\"nn\">markata</span><span class=\"w\"> </span><span
    class=\"kn\">import</span> <span class=\"n\">Markata</span>\n\n<span class=\"c1\">#
    Load config from file</span>\n<span class=\"n\">markata</span> <span class=\"o\">=</span>
    <span class=\"n\">Markata</span><span class=\"o\">.</span><span class=\"n\">from_file</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;markata.toml&quot;</span><span class=\"p\">)</span>\n\n<span
    class=\"c1\"># Access configuration</span>\n<span class=\"nb\">print</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">title</span><span
    class=\"p\">)</span>         <span class=\"c1\"># Site title</span>\n<span class=\"nb\">print</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">url</span><span
    class=\"p\">)</span>           <span class=\"c1\"># Site URL</span>\n<span class=\"nb\">print</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">content_dir</span><span
    class=\"p\">)</span>   <span class=\"c1\"># Content directory</span>\n\n<span
    class=\"c1\"># Access plugin config</span>\n<span class=\"nb\">print</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">feeds</span><span
    class=\"o\">.</span><span class=\"n\">rss</span><span class=\"p\">)</span>     <span
    class=\"c1\"># RSS feed config</span>\n<span class=\"nb\">print</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">template</span><span class=\"p\">)</span>
    \     <span class=\"c1\"># Template config</span>\n\n<span class=\"c1\"># Override
    config</span>\n<span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">debug</span>
    <span class=\"o\">=</span> <span class=\"kc\">True</span>\n<span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">parallel</span> <span class=\"o\">=</span> <span class=\"kc\">False</span>\n</pre></div>\n\n</pre>\n\n<p>See
    <a href=\"http://hookspec.py\">hookspec.py</a> for plugin development and <a href=\"http://lifecycle.py\">lifecycle.py</a>
    for build process details.</p>\n<hr />\n<div class=\"admonition function\">\n<p
    class=\"admonition-title\">Function</p>\n<h2 id=\"_get_global_path_specs\" class=\"admonition-title\"
    style=\"margin: 0; padding: .5rem 1rem;\">_get_global_path_specs <em class=\"small\">function</em></h2>\n<p>Generate
    a list of standard pathspecs for global config files.</p>\n<p>Args:\ntool (str):
    name of the tool to configure</p>\n</div>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">_get_global_path_specs <em
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
    class=\"w\"> </span><span class=\"nf\">_get_global_path_specs</span><span class=\"p\">(</span><span
    class=\"n\">tool</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">path_spec_type</span><span
    class=\"p\">:</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">    Generate a list of standard pathspecs for global config files.</span>\n\n<span
    class=\"sd\">    Args:</span>\n<span class=\"sd\">        tool (str): name of
    the tool to configure</span>\n<span class=\"sd\">    &quot;&quot;&quot;</span>\n
    \   <span class=\"k\">try</span><span class=\"p\">:</span>\n        <span class=\"n\">home</span>
    <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">os</span><span class=\"o\">.</span><span class=\"n\">environ</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;XDG_HOME&quot;</span><span class=\"p\">])</span>\n
    \   <span class=\"k\">except</span> <span class=\"ne\">KeyError</span><span class=\"p\">:</span>\n
    \       <span class=\"n\">home</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"o\">.</span><span class=\"n\">home</span><span class=\"p\">()</span>\n\n
    \   <span class=\"k\">return</span> <span class=\"p\">[</span>\n        <span
    class=\"p\">{</span><span class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">home</span> <span class=\"o\">/</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;</span><span class=\"si\">{</span><span class=\"n\">tool</span><span
    class=\"si\">}</span><span class=\"s2\">.ini&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;parser&quot;</span><span class=\"p\">:</span> <span class=\"s2\">&quot;ini&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span>
    <span class=\"p\">[</span><span class=\"n\">tool</span><span class=\"p\">]},</span>\n
    \       <span class=\"p\">{</span><span class=\"s2\">&quot;path_specs&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">home</span> <span class=\"o\">/</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;.</span><span class=\"si\">{</span><span
    class=\"n\">tool</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;parser&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;ini&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;keys&quot;</span><span
    class=\"p\">:</span> <span class=\"p\">[</span><span class=\"n\">tool</span><span
    class=\"p\">]},</span>\n        <span class=\"p\">{</span><span class=\"s2\">&quot;path_specs&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">home</span> <span class=\"o\">/</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;.</span><span class=\"si\">{</span><span
    class=\"n\">tool</span><span class=\"si\">}</span><span class=\"s2\">.ini&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;parser&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;ini&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;keys&quot;</span><span
    class=\"p\">:</span> <span class=\"p\">[</span><span class=\"n\">tool</span><span
    class=\"p\">]},</span>\n        <span class=\"p\">{</span>\n            <span
    class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">:</span> <span class=\"n\">home</span>
    <span class=\"o\">/</span> <span class=\"s2\">&quot;.config&quot;</span> <span
    class=\"o\">/</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">tool</span><span class=\"si\">}</span><span
    class=\"s2\">.ini&quot;</span><span class=\"p\">,</span>\n            <span class=\"s2\">&quot;parser&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;ini&quot;</span><span class=\"p\">,</span>\n
    \           <span class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span>
    <span class=\"p\">[</span><span class=\"n\">tool</span><span class=\"p\">],</span>\n
    \       <span class=\"p\">},</span>\n        <span class=\"p\">{</span>\n            <span
    class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">:</span> <span class=\"n\">home</span>
    <span class=\"o\">/</span> <span class=\"s2\">&quot;.config&quot;</span> <span
    class=\"o\">/</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;.</span><span
    class=\"si\">{</span><span class=\"n\">tool</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span><span class=\"p\">,</span>\n            <span class=\"s2\">&quot;parser&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;ini&quot;</span><span class=\"p\">,</span>\n
    \           <span class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span>
    <span class=\"p\">[</span><span class=\"n\">tool</span><span class=\"p\">],</span>\n
    \       <span class=\"p\">},</span>\n        <span class=\"p\">{</span>\n            <span
    class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">:</span> <span class=\"n\">home</span>
    <span class=\"o\">/</span> <span class=\"s2\">&quot;.config&quot;</span> <span
    class=\"o\">/</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;.</span><span
    class=\"si\">{</span><span class=\"n\">tool</span><span class=\"si\">}</span><span
    class=\"s2\">.ini&quot;</span><span class=\"p\">,</span>\n            <span class=\"s2\">&quot;parser&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;ini&quot;</span><span class=\"p\">,</span>\n
    \           <span class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span>
    <span class=\"p\">[</span><span class=\"n\">tool</span><span class=\"p\">],</span>\n
    \       <span class=\"p\">},</span>\n    <span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition function\">\n<p class=\"admonition-title\">Function</p>\n<h2
    id=\"_get_local_path_specs\" class=\"admonition-title\" style=\"margin: 0; padding:
    .5rem 1rem;\">_get_local_path_specs <em class=\"small\">function</em></h2>\n<p>Generate
    a list of standard pathspecs for local, project directory config files.</p>\n<p>Args:\ntool
    (str): name of the tool to configure</p>\n</div>\n<div class=\"admonition source
    is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_get_local_path_specs
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
    class=\"w\"> </span><span class=\"nf\">_get_local_path_specs</span><span class=\"p\">(</span><span
    class=\"n\">tool</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"n\">project_home</span><span class=\"p\">:</span>
    <span class=\"n\">Union</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"n\">Path</span><span class=\"p\">])</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">path_spec_type</span><span class=\"p\">:</span>\n<span
    class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span class=\"sd\">
    \   Generate a list of standard pathspecs for local, project directory config
    files.</span>\n\n<span class=\"sd\">    Args:</span>\n<span class=\"sd\">        tool
    (str): name of the tool to configure</span>\n<span class=\"sd\">    &quot;&quot;&quot;</span>\n
    \   <span class=\"k\">return</span> <span class=\"p\">[</span>\n        <span
    class=\"p\">{</span>\n            <span class=\"s2\">&quot;path_specs&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">project_home</span><span class=\"p\">)</span> <span class=\"o\">/</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">tool</span><span class=\"si\">}</span><span class=\"s2\">.ini&quot;</span><span
    class=\"p\">,</span>\n            <span class=\"s2\">&quot;parser&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;ini&quot;</span><span class=\"p\">,</span>\n
    \           <span class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span>
    <span class=\"p\">[</span><span class=\"n\">tool</span><span class=\"p\">],</span>\n
    \       <span class=\"p\">},</span>\n        <span class=\"p\">{</span>\n            <span
    class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">:</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">project_home</span><span class=\"p\">)</span>
    <span class=\"o\">/</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;.</span><span
    class=\"si\">{</span><span class=\"n\">tool</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span><span class=\"p\">,</span>\n            <span class=\"s2\">&quot;parser&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;ini&quot;</span><span class=\"p\">,</span>\n
    \           <span class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span>
    <span class=\"p\">[</span><span class=\"n\">tool</span><span class=\"p\">],</span>\n
    \       <span class=\"p\">},</span>\n        <span class=\"p\">{</span>\n            <span
    class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">:</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">project_home</span><span class=\"p\">)</span>
    <span class=\"o\">/</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;.</span><span
    class=\"si\">{</span><span class=\"n\">tool</span><span class=\"si\">}</span><span
    class=\"s2\">.ini&quot;</span><span class=\"p\">,</span>\n            <span class=\"s2\">&quot;parser&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;ini&quot;</span><span class=\"p\">,</span>\n
    \           <span class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span>
    <span class=\"p\">[</span><span class=\"n\">tool</span><span class=\"p\">],</span>\n
    \       <span class=\"p\">},</span>\n        <span class=\"p\">{</span>\n            <span
    class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">:</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">project_home</span><span class=\"p\">)</span>
    <span class=\"o\">/</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">tool</span><span class=\"si\">}</span><span
    class=\"s2\">.yml&quot;</span><span class=\"p\">,</span>\n            <span class=\"s2\">&quot;parser&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;yaml&quot;</span><span class=\"p\">,</span>\n
    \           <span class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span>
    <span class=\"p\">[</span><span class=\"n\">tool</span><span class=\"p\">],</span>\n
    \       <span class=\"p\">},</span>\n        <span class=\"p\">{</span>\n            <span
    class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">:</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">project_home</span><span class=\"p\">)</span>
    <span class=\"o\">/</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;.</span><span
    class=\"si\">{</span><span class=\"n\">tool</span><span class=\"si\">}</span><span
    class=\"s2\">.yml&quot;</span><span class=\"p\">,</span>\n            <span class=\"s2\">&quot;parser&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;yaml&quot;</span><span class=\"p\">,</span>\n
    \           <span class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span>
    <span class=\"p\">[</span><span class=\"n\">tool</span><span class=\"p\">],</span>\n
    \       <span class=\"p\">},</span>\n        <span class=\"p\">{</span>\n            <span
    class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">:</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">project_home</span><span class=\"p\">)</span>
    <span class=\"o\">/</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">tool</span><span class=\"si\">}</span><span
    class=\"s2\">.toml&quot;</span><span class=\"p\">,</span>\n            <span class=\"s2\">&quot;parser&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;toml&quot;</span><span class=\"p\">,</span>\n
    \           <span class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span>
    <span class=\"p\">[</span><span class=\"n\">tool</span><span class=\"p\">],</span>\n
    \       <span class=\"p\">},</span>\n        <span class=\"p\">{</span>\n            <span
    class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">:</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">project_home</span><span class=\"p\">)</span>
    <span class=\"o\">/</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;.</span><span
    class=\"si\">{</span><span class=\"n\">tool</span><span class=\"si\">}</span><span
    class=\"s2\">.toml&quot;</span><span class=\"p\">,</span>\n            <span class=\"s2\">&quot;parser&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;toml&quot;</span><span class=\"p\">,</span>\n
    \           <span class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span>
    <span class=\"p\">[</span><span class=\"n\">tool</span><span class=\"p\">],</span>\n
    \       <span class=\"p\">},</span>\n        <span class=\"p\">{</span>\n            <span
    class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">:</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">project_home</span><span class=\"p\">)</span>
    <span class=\"o\">/</span> <span class=\"s2\">&quot;pyproject.toml&quot;</span><span
    class=\"p\">,</span>\n            <span class=\"s2\">&quot;parser&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;toml&quot;</span><span class=\"p\">,</span>\n
    \           <span class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span>
    <span class=\"p\">[</span><span class=\"s2\">&quot;tool&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">tool</span><span class=\"p\">],</span>\n        <span class=\"p\">},</span>\n
    \       <span class=\"p\">{</span>\n            <span class=\"s2\">&quot;path_specs&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">project_home</span><span class=\"p\">)</span> <span class=\"o\">/</span>
    <span class=\"s2\">&quot;setup.cfg&quot;</span><span class=\"p\">,</span>\n            <span
    class=\"s2\">&quot;parser&quot;</span><span class=\"p\">:</span> <span class=\"s2\">&quot;ini&quot;</span><span
    class=\"p\">,</span>\n            <span class=\"s2\">&quot;keys&quot;</span><span
    class=\"p\">:</span> <span class=\"p\">[</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;tool.</span><span class=\"si\">{</span><span class=\"n\">tool</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">],</span>\n
    \       <span class=\"p\">},</span>\n    <span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition function\">\n<p class=\"admonition-title\">Function</p>\n<h2
    id=\"_get_attrs\" class=\"admonition-title\" style=\"margin: 0; padding: .5rem
    1rem;\">_get_attrs <em class=\"small\">function</em></h2>\n<p>Get nested config
    data from a list of keys.</p>\n<p>specifically written for pyproject.toml which
    needs to get <code>tool</code> then <code>&lt;tool&gt;</code></p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_get_attrs
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
    class=\"w\"> </span><span class=\"nf\">_get_attrs</span><span class=\"p\">(</span><span
    class=\"n\">attrs</span><span class=\"p\">:</span> <span class=\"nb\">list</span><span
    class=\"p\">,</span> <span class=\"n\">config</span><span class=\"p\">:</span>
    <span class=\"n\">Dict</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Dict</span><span class=\"p\">:</span>\n<span class=\"w\">    </span><span
    class=\"sd\">&quot;&quot;&quot;Get nested config data from a list of keys.</span>\n\n<span
    class=\"sd\">    specifically written for pyproject.toml which needs to get `tool`
    then `&lt;tool&gt;`</span>\n<span class=\"sd\">    &quot;&quot;&quot;</span>\n
    \   <span class=\"k\">for</span> <span class=\"n\">attr</span> <span class=\"ow\">in</span>
    <span class=\"n\">attrs</span><span class=\"p\">:</span>\n        <span class=\"n\">config</span>
    <span class=\"o\">=</span> <span class=\"n\">config</span><span class=\"p\">[</span><span
    class=\"n\">attr</span><span class=\"p\">]</span>\n    <span class=\"k\">return</span>
    <span class=\"n\">config</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div class=\"admonition
    function\">\n<p class=\"admonition-title\">Function</p>\n<h2 id=\"_load_config_file\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">_load_config_file
    <em class=\"small\">function</em></h2>\n<p>Load a configuration file using the
    appropriate parser.</p>\n<p>Args:\nfile_spec: Dictionary containing path_specs,
    parser, and keys information</p>\n<p>Returns:\nOptional[Dict[str, Any]]: Parsed
    configuration or None if file doesn't exist</p>\n</div>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_load_config_file
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
    class=\"w\"> </span><span class=\"nf\">_load_config_file</span><span class=\"p\">(</span><span
    class=\"n\">file_spec</span><span class=\"p\">:</span> <span class=\"n\">Dict</span><span
    class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">,</span> <span
    class=\"n\">Any</span><span class=\"p\">])</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"n\">Dict</span><span
    class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">,</span> <span
    class=\"n\">Any</span><span class=\"p\">]]:</span>\n<span class=\"w\">    </span><span
    class=\"sd\">&quot;&quot;&quot;Load a configuration file using the appropriate
    parser.</span>\n\n<span class=\"sd\">    Args:</span>\n<span class=\"sd\">        file_spec:
    Dictionary containing path_specs, parser, and keys information</span>\n\n<span
    class=\"sd\">    Returns:</span>\n<span class=\"sd\">        Optional[Dict[str,
    Any]]: Parsed configuration or None if file doesn&#39;t exist</span>\n<span class=\"sd\">
    \   &quot;&quot;&quot;</span>\n    <span class=\"n\">path</span> <span class=\"o\">=</span>
    <span class=\"n\">file_spec</span><span class=\"p\">[</span><span class=\"s2\">&quot;path_specs&quot;</span><span
    class=\"p\">]</span>\n    <span class=\"k\">if</span> <span class=\"ow\">not</span>
    <span class=\"n\">path</span><span class=\"o\">.</span><span class=\"n\">exists</span><span
    class=\"p\">():</span>\n        <span class=\"k\">return</span> <span class=\"kc\">None</span>\n\n
    \   <span class=\"k\">try</span><span class=\"p\">:</span>\n        <span class=\"k\">if</span>
    <span class=\"n\">file_spec</span><span class=\"p\">[</span><span class=\"s2\">&quot;parser&quot;</span><span
    class=\"p\">]</span> <span class=\"o\">==</span> <span class=\"s2\">&quot;toml&quot;</span><span
    class=\"p\">:</span>\n            <span class=\"k\">with</span> <span class=\"nb\">open</span><span
    class=\"p\">(</span><span class=\"n\">path</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;rb&quot;</span><span class=\"p\">)</span> <span class=\"k\">as</span>
    <span class=\"n\">f</span><span class=\"p\">:</span>\n                <span class=\"n\">config</span>
    <span class=\"o\">=</span> <span class=\"n\">tomli</span><span class=\"o\">.</span><span
    class=\"n\">load</span><span class=\"p\">(</span><span class=\"n\">f</span><span
    class=\"p\">)</span>\n        <span class=\"k\">elif</span> <span class=\"n\">file_spec</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;parser&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">==</span> <span class=\"s2\">&quot;yaml&quot;</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">with</span> <span class=\"nb\">open</span><span
    class=\"p\">(</span><span class=\"n\">path</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;r&quot;</span><span class=\"p\">)</span> <span class=\"k\">as</span>
    <span class=\"n\">f</span><span class=\"p\">:</span>\n                <span class=\"n\">config</span>
    <span class=\"o\">=</span> <span class=\"n\">yaml</span><span class=\"o\">.</span><span
    class=\"n\">safe_load</span><span class=\"p\">(</span><span class=\"n\">f</span><span
    class=\"p\">)</span>\n        <span class=\"k\">elif</span> <span class=\"n\">file_spec</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;parser&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">==</span> <span class=\"s2\">&quot;ini&quot;</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">config</span> <span class=\"o\">=</span> <span class=\"n\">configparser</span><span
    class=\"o\">.</span><span class=\"n\">ConfigParser</span><span class=\"p\">()</span>\n
    \           <span class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">read</span><span
    class=\"p\">(</span><span class=\"n\">path</span><span class=\"p\">)</span>\n
    \           <span class=\"c1\"># Convert ConfigParser to dict</span>\n            <span
    class=\"n\">config</span> <span class=\"o\">=</span> <span class=\"p\">{</span><span
    class=\"n\">s</span><span class=\"p\">:</span> <span class=\"nb\">dict</span><span
    class=\"p\">(</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">items</span><span class=\"p\">(</span><span class=\"n\">s</span><span
    class=\"p\">))</span> <span class=\"k\">for</span> <span class=\"n\">s</span>
    <span class=\"ow\">in</span> <span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">sections</span><span class=\"p\">()}</span>\n        <span class=\"k\">else</span><span
    class=\"p\">:</span>\n            <span class=\"k\">return</span> <span class=\"kc\">None</span>\n\n
    \       <span class=\"k\">return</span> <span class=\"n\">_get_attrs</span><span
    class=\"p\">(</span><span class=\"n\">file_spec</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;keys&quot;</span><span class=\"p\">],</span> <span class=\"n\">config</span><span
    class=\"p\">)</span>\n    <span class=\"k\">except</span> <span class=\"p\">(</span>\n
    \       <span class=\"ne\">KeyError</span><span class=\"p\">,</span>\n        <span
    class=\"ne\">TypeError</span><span class=\"p\">,</span>\n        <span class=\"n\">yaml</span><span
    class=\"o\">.</span><span class=\"n\">YAMLError</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">tomli</span><span class=\"o\">.</span><span class=\"n\">TOMLDecodeError</span><span
    class=\"p\">,</span>\n        <span class=\"n\">configparser</span><span class=\"o\">.</span><span
    class=\"n\">Error</span><span class=\"p\">,</span>\n    <span class=\"p\">):</span>\n
    \       <span class=\"k\">return</span> <span class=\"kc\">None</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition function\">\n<p class=\"admonition-title\">Function</p>\n<h2
    id=\"_load_files\" class=\"admonition-title\" style=\"margin: 0; padding: .5rem
    1rem;\">_load_files <em class=\"small\">function</em></h2>\n<p>Load config files
    stopping at the first one that exists and can be parsed.</p>\n<p>Args:\nconfig_path_specs:
    List of path specifications to try</p>\n<p>Returns:\nDict[str, Any]: Configuration
    dictionary</p>\n</div>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">_load_files <em class='small'>source</em></p>\n<pre
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
    class=\"w\"> </span><span class=\"nf\">_load_files</span><span class=\"p\">(</span><span
    class=\"n\">config_path_specs</span><span class=\"p\">:</span> <span class=\"n\">path_spec_type</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Dict</span><span
    class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">,</span> <span
    class=\"n\">Any</span><span class=\"p\">]:</span>\n<span class=\"w\">    </span><span
    class=\"sd\">&quot;&quot;&quot;Load config files stopping at the first one that
    exists and can be parsed.</span>\n\n<span class=\"sd\">    Args:</span>\n<span
    class=\"sd\">        config_path_specs: List of path specifications to try</span>\n\n<span
    class=\"sd\">    Returns:</span>\n<span class=\"sd\">        Dict[str, Any]: Configuration
    dictionary</span>\n<span class=\"sd\">    &quot;&quot;&quot;</span>\n    <span
    class=\"k\">for</span> <span class=\"n\">file_spec</span> <span class=\"ow\">in</span>
    <span class=\"n\">config_path_specs</span><span class=\"p\">:</span>\n        <span
    class=\"n\">config</span> <span class=\"o\">=</span> <span class=\"n\">_load_config_file</span><span
    class=\"p\">(</span><span class=\"n\">file_spec</span><span class=\"p\">)</span>\n
    \       <span class=\"k\">if</span> <span class=\"n\">config</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">return</span> <span class=\"n\">config</span>\n
    \   <span class=\"k\">return</span> <span class=\"p\">{}</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition function\">\n<p class=\"admonition-title\">Function</p>\n<h2
    id=\"_load_env\" class=\"admonition-title\" style=\"margin: 0; padding: .5rem
    1rem;\">_load_env <em class=\"small\">function</em></h2>\n<p>Load config from
    environment variables.</p>\n<p>Args:\ntool (str): name of the tool to configure</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_load_env
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
    class=\"w\"> </span><span class=\"nf\">_load_env</span><span class=\"p\">(</span><span
    class=\"n\">tool</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Dict</span><span
    class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">,</span> <span
    class=\"n\">Any</span><span class=\"p\">]:</span>\n<span class=\"w\">    </span><span
    class=\"sd\">&quot;&quot;&quot;Load config from environment variables.</span>\n\n<span
    class=\"sd\">    Args:</span>\n<span class=\"sd\">        tool (str): name of
    the tool to configure</span>\n<span class=\"sd\">    &quot;&quot;&quot;</span>\n
    \   <span class=\"n\">env_prefix</span> <span class=\"o\">=</span> <span class=\"n\">tool</span><span
    class=\"o\">.</span><span class=\"n\">upper</span><span class=\"p\">()</span>\n
    \   <span class=\"n\">env_config</span> <span class=\"o\">=</span> <span class=\"p\">{</span>\n
    \       <span class=\"n\">key</span><span class=\"o\">.</span><span class=\"n\">replace</span><span
    class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">env_prefix</span><span class=\"si\">}</span><span
    class=\"s2\">_&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;&quot;</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">lower</span><span
    class=\"p\">():</span> <span class=\"n\">value</span>\n        <span class=\"k\">for</span>
    <span class=\"n\">key</span><span class=\"p\">,</span> <span class=\"n\">value</span>
    <span class=\"ow\">in</span> <span class=\"n\">os</span><span class=\"o\">.</span><span
    class=\"n\">environ</span><span class=\"o\">.</span><span class=\"n\">items</span><span
    class=\"p\">()</span>\n        <span class=\"k\">if</span> <span class=\"n\">key</span><span
    class=\"o\">.</span><span class=\"n\">startswith</span><span class=\"p\">(</span><span
    class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">env_prefix</span><span class=\"si\">}</span><span class=\"s2\">_&quot;</span><span
    class=\"p\">)</span>\n    <span class=\"p\">}</span>\n    <span class=\"k\">return</span>
    <span class=\"n\">env_config</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div class=\"admonition
    function\">\n<p class=\"admonition-title\">Function</p>\n<h2 id=\"load\" class=\"admonition-title\"
    style=\"margin: 0; padding: .5rem 1rem;\">load <em class=\"small\">function</em></h2>\n<p>Load
    tool config from standard config files.</p>\n<p>Resolution Order</p>\n<ul>\n<li>First
    global file with a tool key</li>\n<li>First local file with a tool key</li>\n<li>Environment
    variables prefixed with <code>TOOL</code></li>\n<li>Overrides</li>\n</ul>\n<p>Args:\ntool
    (str): name of the tool to configure\nproject_home (Union[Path, str], optional):
    Project directory to search for config files. Defaults to &quot;.&quot;.\noverrides
    (Dict, optional): Override values to apply last. Defaults to None.</p>\n<p>Returns:\nDict[str,
    Any]: Configuration object</p>\n</div>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">load <em class='small'>source</em></p>\n<pre
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
    class=\"w\"> </span><span class=\"nf\">load</span><span class=\"p\">(</span>\n
    \   <span class=\"n\">tool</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">,</span>\n    <span class=\"n\">project_home</span><span class=\"p\">:</span>
    <span class=\"n\">Union</span><span class=\"p\">[</span><span class=\"n\">Path</span><span
    class=\"p\">,</span> <span class=\"nb\">str</span><span class=\"p\">]</span> <span
    class=\"o\">=</span> <span class=\"s2\">&quot;.&quot;</span><span class=\"p\">,</span>\n
    \   <span class=\"n\">overrides</span><span class=\"p\">:</span> <span class=\"n\">Optional</span><span
    class=\"p\">[</span><span class=\"n\">Dict</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">,</span> <span class=\"n\">Any</span><span
    class=\"p\">]]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span><span
    class=\"p\">,</span>\n<span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Dict</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"n\">Any</span><span class=\"p\">]:</span>\n<span
    class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;Load tool config from
    standard config files.</span>\n\n<span class=\"sd\">    Resolution Order</span>\n\n<span
    class=\"sd\">    * First global file with a tool key</span>\n<span class=\"sd\">
    \   * First local file with a tool key</span>\n<span class=\"sd\">    * Environment
    variables prefixed with `TOOL`</span>\n<span class=\"sd\">    * Overrides</span>\n\n<span
    class=\"sd\">    Args:</span>\n<span class=\"sd\">        tool (str): name of
    the tool to configure</span>\n<span class=\"sd\">        project_home (Union[Path,
    str], optional): Project directory to search for config files. Defaults to &quot;.&quot;.</span>\n<span
    class=\"sd\">        overrides (Dict, optional): Override values to apply last.
    Defaults to None.</span>\n\n<span class=\"sd\">    Returns:</span>\n<span class=\"sd\">
    \       Dict[str, Any]: Configuration object</span>\n<span class=\"sd\">    &quot;&quot;&quot;</span>\n
    \   <span class=\"n\">overrides</span> <span class=\"o\">=</span> <span class=\"n\">overrides</span>
    <span class=\"ow\">or</span> <span class=\"p\">{}</span>\n    <span class=\"n\">config</span>
    <span class=\"o\">=</span> <span class=\"p\">{}</span>\n\n    <span class=\"c1\">#
    Load from files in order of precedence</span>\n    <span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">update</span><span class=\"p\">(</span><span
    class=\"n\">_load_files</span><span class=\"p\">(</span><span class=\"n\">_get_global_path_specs</span><span
    class=\"p\">(</span><span class=\"n\">tool</span><span class=\"p\">))</span> <span
    class=\"ow\">or</span> <span class=\"p\">{})</span>\n    <span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">update</span><span class=\"p\">(</span><span
    class=\"n\">_load_files</span><span class=\"p\">(</span><span class=\"n\">_get_local_path_specs</span><span
    class=\"p\">(</span><span class=\"n\">tool</span><span class=\"p\">,</span> <span
    class=\"n\">project_home</span><span class=\"p\">))</span> <span class=\"ow\">or</span>
    <span class=\"p\">{})</span>\n    <span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">update</span><span class=\"p\">(</span><span class=\"n\">_load_env</span><span
    class=\"p\">(</span><span class=\"n\">tool</span><span class=\"p\">))</span>\n
    \   <span class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">update</span><span
    class=\"p\">(</span><span class=\"n\">overrides</span><span class=\"p\">)</span>\n\n
    \   <span class=\"c1\"># If no settings class is provided, return the raw dict</span>\n
    \   <span class=\"k\">return</span> <span class=\"n\">config</span>\n</pre></div>\n\n</pre>\n\n</div>\n\n
    \   </section>\n</article>"
  raw.md: ''
published: false
slug: markata/standard-config
title: standard_config.py


---

---

Standard Config.
A module to load tooling config from a users project space.

Inspired from frustrations that some tools have a tool.ini, .tool.ini,
setup.cfg, or pyproject.toml.  Some allow for global configs, some don't.  Some
properly follow the users home directory, others end up in a weird temp
directory.  Windows home directory is only more confusing.  Some will even
respect the users `$XDG_HOME` directory.

This file is for any project that can be configured in plain text such as `ini`
or `toml` and not requiring a .py file.  Just name your tool and let users put
config where it makes sense to them, no need to figure out resolution order.

## Usage:

``` python
from standard_config import load

# Retrieve any overrides from the user
overrides = {'setting': True}
config = load('my_tool', overrides)
```

## Resolution Order

* First global file with a tool key
* First local file with a tool key
* Environment variables prefixed with `TOOL`
* Overrides

### Tool Specific Ini files

Ini file formats must include a `<tool>` key.

``` ini
[my_tool]
setting = True
```

### pyproject.toml

Toml files must include a `tool.<tool>` key

``` toml
[tool.my_tool]
setting = True
```

### setup.cfg

setup.cfg files must include a `tool:<tool>` key

``` ini
[tool:my_tool]
setting = True
```

### global files to consider

* <home>/tool.ini
* <home>/.tool
* <home>/.tool.ini
* <home>/.config/tool.ini
* <home>/.config/.tool
* <home>/.config/.tool.ini

### local files to consider

* <project_home>/tool.ini
* <project_home>/.tool
* <project_home>/.tool.ini
* <project_home>/pyproject.toml
* <project_home>/setup.cfg

Markata's standard configuration system.

## Configuration Overview

Markata uses a hierarchical configuration system based on Pydantic models. Configuration
can be set through:
1. TOML files
2. Environment variables
3. Command line arguments

# Basic Configuration

Minimal `markata.toml`:
```toml
[markata]
# Site info
title = "My Site"
url = "https://example.com"
description = "Site description"

# Content locations
content_dir = "content"
output_dir = "markout"
assets_dir = "static"

# Plugin management
hooks = ["default"]
```

# Environment Variables

All settings can be overridden with environment variables:
```bash
# Override site URL
export MARKATA_URL="https://staging.example.com"

# Override output directory
export MARKATA_OUTPUT_DIR="dist"

# Enable debug mode
export MARKATA_DEBUG=1
```

# Detailed Configuration

## Core Settings

```toml
[markata]
# Site information
title = "My Site"                  # Site title
url = "https://example.com"        # Base URL
description = "Site description"   # Meta description
author_name = "Author Name"        # Author name
author_email = "me@example.com"    # Author email
icon = "favicon.ico"               # Site icon
lang = "en"                        # Site language

# Content locations
content_dir = "content"           # Source content location
output_dir = "markout"            # Build output location
assets_dir = "static"             # Static assets location
template_dir = "templates"        # Template location

# Plugin management
hooks = ["default"]               # Active plugins
disabled_hooks = []               # Disabled plugins
```

## Cache Settings

```toml
[markata]
# Cache configuration
default_cache_expire = 3600       # Default TTL (1 hour)
template_cache_expire = 86400     # Template TTL (24 hours)
markdown_cache_expire = 21600     # Markdown TTL (6 hours)
dynamic_cache_expire = 3600       # Dynamic TTL (1 hour)
```

## Development Settings

```toml
[markata]
# Development server
dev_server_port = 8000            # Local server port
dev_server_host = "localhost"     # Local server host
debug = false                     # Debug mode

# Performance
parallel = true                   # Enable parallel processing
workers = 4                       # Number of worker threads
```

## Content Settings

```toml
[markata]
# Content processing
default_template = "post.html"    # Default template
markdown_extensions = [           # Markdown extensions
    "fenced_code",
    "tables",
    "footnotes"
]

# Content filtering
draft = false                     # Include drafts
future = false                    # Include future posts
```

# Plugin Configuration

Each plugin can define its own configuration section:

```toml
# RSS feed configuration
[markata.feeds]
rss = { output = "rss.xml" }
atom = { output = "atom.xml" }
json = { output = "feed.json" }

# Template configuration
[markata.template]
engine = "jinja2"
cache_size = 100
autoescape = true

# Markdown configuration
[markata.markdown]
highlight_theme = "monokai"
line_numbers = true
```

## Configuration Validation

The configuration is validated using Pydantic models:

```python
from pydantic import BaseModel, Field

class MarkataConfig(BaseModel):
    """Core configuration model."""
    # Site info
    title: str = Field(..., description="Site title")
    url: str = Field(..., description="Site base URL")

    # Directories
    content_dir: Path = Field("content", description="Content directory")
    output_dir: Path = Field("markout", description="Output directory")

    # Features
    debug: bool = Field(False, description="Enable debug mode")
    parallel: bool = Field(True, description="Enable parallel processing")

    model_config = ConfigDict(
        validate_assignment=True,
        arbitrary_types_allowed=True,
        extra="allow",
        str_strip_whitespace=True,
        validate_default=True,
        populate_by_name=True,
    )
```

# Usage Example

```python
from markata import Markata

# Load config from file
markata = Markata.from_file("markata.toml")

# Access configuration
print(markata.config.title)         # Site title
print(markata.config.url)           # Site URL
print(markata.config.content_dir)   # Content directory

# Access plugin config
print(markata.config.feeds.rss)     # RSS feed config
print(markata.config.template)      # Template config

# Override config
markata.config.debug = True
markata.config.parallel = False
```

See hookspec.py for plugin development and lifecycle.py for build process details.

---

!!! function
    <h2 id="_get_global_path_specs" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">_get_global_path_specs <em class="small">function</em></h2>

    Generate a list of standard pathspecs for global config files.

    Args:
        tool (str): name of the tool to configure

???+ source "_get_global_path_specs <em class='small'>source</em>"
    ```python
    def _get_global_path_specs(tool: str) -> path_spec_type:
        """
        Generate a list of standard pathspecs for global config files.

        Args:
            tool (str): name of the tool to configure
        """
        try:
            home = Path(os.environ["XDG_HOME"])
        except KeyError:
            home = Path.home()

        return [
            {"path_specs": home / f"{tool}.ini", "parser": "ini", "keys": [tool]},
            {"path_specs": home / f".{tool}", "parser": "ini", "keys": [tool]},
            {"path_specs": home / f".{tool}.ini", "parser": "ini", "keys": [tool]},
            {
                "path_specs": home / ".config" / f"{tool}.ini",
                "parser": "ini",
                "keys": [tool],
            },
            {
                "path_specs": home / ".config" / f".{tool}",
                "parser": "ini",
                "keys": [tool],
            },
            {
                "path_specs": home / ".config" / f".{tool}.ini",
                "parser": "ini",
                "keys": [tool],
            },
        ]
    ```
!!! function
    <h2 id="_get_local_path_specs" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">_get_local_path_specs <em class="small">function</em></h2>

    Generate a list of standard pathspecs for local, project directory config files.

    Args:
        tool (str): name of the tool to configure

???+ source "_get_local_path_specs <em class='small'>source</em>"
    ```python
    def _get_local_path_specs(tool: str, project_home: Union[str, Path]) -> path_spec_type:
        """
        Generate a list of standard pathspecs for local, project directory config files.

        Args:
            tool (str): name of the tool to configure
        """
        return [
            {
                "path_specs": Path(project_home) / f"{tool}.ini",
                "parser": "ini",
                "keys": [tool],
            },
            {
                "path_specs": Path(project_home) / f".{tool}",
                "parser": "ini",
                "keys": [tool],
            },
            {
                "path_specs": Path(project_home) / f".{tool}.ini",
                "parser": "ini",
                "keys": [tool],
            },
            {
                "path_specs": Path(project_home) / f"{tool}.yml",
                "parser": "yaml",
                "keys": [tool],
            },
            {
                "path_specs": Path(project_home) / f".{tool}.yml",
                "parser": "yaml",
                "keys": [tool],
            },
            {
                "path_specs": Path(project_home) / f"{tool}.toml",
                "parser": "toml",
                "keys": [tool],
            },
            {
                "path_specs": Path(project_home) / f".{tool}.toml",
                "parser": "toml",
                "keys": [tool],
            },
            {
                "path_specs": Path(project_home) / "pyproject.toml",
                "parser": "toml",
                "keys": ["tool", tool],
            },
            {
                "path_specs": Path(project_home) / "setup.cfg",
                "parser": "ini",
                "keys": [f"tool.{tool}"],
            },
        ]
    ```
!!! function
    <h2 id="_get_attrs" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">_get_attrs <em class="small">function</em></h2>

    Get nested config data from a list of keys.

    specifically written for pyproject.toml which needs to get `tool` then `<tool>`

???+ source "_get_attrs <em class='small'>source</em>"
    ```python
    def _get_attrs(attrs: list, config: Dict) -> Dict:
        """Get nested config data from a list of keys.

        specifically written for pyproject.toml which needs to get `tool` then `<tool>`
        """
        for attr in attrs:
            config = config[attr]
        return config
    ```
!!! function
    <h2 id="_load_config_file" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">_load_config_file <em class="small">function</em></h2>

    Load a configuration file using the appropriate parser.

    Args:
        file_spec: Dictionary containing path_specs, parser, and keys information

    Returns:
        Optional[Dict[str, Any]]: Parsed configuration or None if file doesn't exist

???+ source "_load_config_file <em class='small'>source</em>"
    ```python
    def _load_config_file(file_spec: Dict[str, Any]) -> Optional[Dict[str, Any]]:
        """Load a configuration file using the appropriate parser.

        Args:
            file_spec: Dictionary containing path_specs, parser, and keys information

        Returns:
            Optional[Dict[str, Any]]: Parsed configuration or None if file doesn't exist
        """
        path = file_spec["path_specs"]
        if not path.exists():
            return None

        try:
            if file_spec["parser"] == "toml":
                with open(path, "rb") as f:
                    config = tomli.load(f)
            elif file_spec["parser"] == "yaml":
                with open(path, "r") as f:
                    config = yaml.safe_load(f)
            elif file_spec["parser"] == "ini":
                config = configparser.ConfigParser()
                config.read(path)
                # Convert ConfigParser to dict
                config = {s: dict(config.items(s)) for s in config.sections()}
            else:
                return None

            return _get_attrs(file_spec["keys"], config)
        except (
            KeyError,
            TypeError,
            yaml.YAMLError,
            tomli.TOMLDecodeError,
            configparser.Error,
        ):
            return None
    ```
!!! function
    <h2 id="_load_files" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">_load_files <em class="small">function</em></h2>

    Load config files stopping at the first one that exists and can be parsed.

    Args:
        config_path_specs: List of path specifications to try

    Returns:
        Dict[str, Any]: Configuration dictionary

???+ source "_load_files <em class='small'>source</em>"
    ```python
    def _load_files(config_path_specs: path_spec_type) -> Dict[str, Any]:
        """Load config files stopping at the first one that exists and can be parsed.

        Args:
            config_path_specs: List of path specifications to try

        Returns:
            Dict[str, Any]: Configuration dictionary
        """
        for file_spec in config_path_specs:
            config = _load_config_file(file_spec)
            if config:
                return config
        return {}
    ```
!!! function
    <h2 id="_load_env" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">_load_env <em class="small">function</em></h2>

    Load config from environment variables.

    Args:
        tool (str): name of the tool to configure

???+ source "_load_env <em class='small'>source</em>"
    ```python
    def _load_env(tool: str) -> Dict[str, Any]:
        """Load config from environment variables.

        Args:
            tool (str): name of the tool to configure
        """
        env_prefix = tool.upper()
        env_config = {
            key.replace(f"{env_prefix}_", "").lower(): value
            for key, value in os.environ.items()
            if key.startswith(f"{env_prefix}_")
        }
        return env_config
    ```
!!! function
    <h2 id="load" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">load <em class="small">function</em></h2>

    Load tool config from standard config files.

    Resolution Order

    * First global file with a tool key
    * First local file with a tool key
    * Environment variables prefixed with `TOOL`
    * Overrides

    Args:
        tool (str): name of the tool to configure
        project_home (Union[Path, str], optional): Project directory to search for config files. Defaults to ".".
        overrides (Dict, optional): Override values to apply last. Defaults to None.

    Returns:
        Dict[str, Any]: Configuration object

???+ source "load <em class='small'>source</em>"
    ```python
    def load(
        tool: str,
        project_home: Union[Path, str] = ".",
        overrides: Optional[Dict[str, Any]] = None,
    ) -> Dict[str, Any]:
        """Load tool config from standard config files.

        Resolution Order

        * First global file with a tool key
        * First local file with a tool key
        * Environment variables prefixed with `TOOL`
        * Overrides

        Args:
            tool (str): name of the tool to configure
            project_home (Union[Path, str], optional): Project directory to search for config files. Defaults to ".".
            overrides (Dict, optional): Override values to apply last. Defaults to None.

        Returns:
            Dict[str, Any]: Configuration object
        """
        overrides = overrides or {}
        config = {}

        # Load from files in order of precedence
        config.update(_load_files(_get_global_path_specs(tool)) or {})
        config.update(_load_files(_get_local_path_specs(tool, project_home)) or {})
        config.update(_load_env(tool))
        config.update(overrides)

        # If no settings class is provided, return the raw dict
        return config
    ```