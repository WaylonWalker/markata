---
content: "Standard Config.\nA module to load tooling config from a users project space.\n\nInspired
  from frustrations that some tools have a tool.ini, .tool.ini,\nsetup.cfg, or pyproject.toml.
  \ Some allow for global configs, some don't.  Some\nproperly follow the users home
  directory, others end up in a weird temp\ndirectory.  Windows home directory is
  only more confusing.  Some will even\nrespect the users `$XDG_HOME` directory.\n\n\nThis
  file is for any project that can be configured in plain text such as `ini`\nor `toml`
  and not requiring a .py file.  Just name your tool and let users put\nconfig where
  it makes sense to them, no need to figure out resolution order.\n\n## Usage:\n\n```
  python\nfrom standard_config import load\n\n# Retrieve any overrides from the user\noverrides
  = {'setting': True}\nconfig = load('my_tool', overrides)\n```\n\n## Resolution Order\n\n*
  First global file with a tool key\n* First local file with a tool key\n* Environment
  variables prefixed with `TOOL`\n* Overrides\n\n### Tool Specific Ini files\n\nIni
  file formats must include a `<tool>` key.\n\n``` ini\n[my_tool]\nsetting = True\n```\n\n###
  pyproject.toml\n\nToml files must include a `tool.<tool>` key\n\n``` toml\n[tool.my_tool]\nsetting
  = True\n```\n\n### setup.cfg\n\nsetup.cfg files must include a `tool:<tool>` key\n\n```
  ini\n[tool:my_tool]\nsetting = True\n```\n\n\n### global files to consider\n\n*
  <home>/tool.ini\n* <home>/.tool\n* <home>/.tool.ini\n* <home>/.config/tool.ini\n*
  <home>/.config/.tool\n* <home>/.config/.tool.ini\n\n### local files to consider\n\n*
  <project_home>/tool.ini\n* <project_home>/.tool\n* <project_home>/.tool.ini\n* <project_home>/pyproject.toml\n*
  <project_home>/setup.cfg\n\n\n!! function <h2 id='_get_global_path_specs' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>_get_global_path_specs <em class='small'>function</em></h2>\n
  \   Generate a list of standard pathspecs for global config files.\n\n    Args:\n
  \       tool (str): name of the tool to configure\n???+ source \"_get_global_path_specs
  <em class='small'>source</em>\"\n\n```python\n\n        def _get_global_path_specs(tool:
  str) -> path_spec_type:\n            \"\"\"\n            Generate a list of standard
  pathspecs for global config files.\n\n            Args:\n                tool (str):
  name of the tool to configure\n            \"\"\"\n            try:\n                home
  = Path(os.environ[\"XDG_HOME\"])\n            except KeyError:\n                home
  = Path.home()\n\n            return [\n                {\"path_specs\": home / f\"{tool}.ini\",
  \"ac_parser\": \"ini\", \"keys\": [tool]},\n                {\"path_specs\": home
  / f\".{tool}\", \"ac_parser\": \"ini\", \"keys\": [tool]},\n                {\"path_specs\":
  home / f\".{tool}.ini\", \"ac_parser\": \"ini\", \"keys\": [tool]},\n                {\n
  \                   \"path_specs\": home / \".config\" / f\"{tool}.ini\",\n                    \"ac_parser\":
  \"ini\",\n                    \"keys\": [tool],\n                },\n                {\n
  \                   \"path_specs\": home / \".config\" / f\".{tool}\",\n                    \"ac_parser\":
  \"ini\",\n                    \"keys\": [tool],\n                },\n                {\n
  \                   \"path_specs\": home / \".config\" / f\".{tool}.ini\",\n                    \"ac_parser\":
  \"ini\",\n                    \"keys\": [tool],\n                },\n            ]\n```\n\n\n!!
  function <h2 id='_get_local_path_specs' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>_get_local_path_specs <em class='small'>function</em></h2>\n    Generate
  a list of standard pathspecs for local, project directory config files.\n\n    Args:\n
  \       tool (str): name of the tool to configure\n???+ source \"_get_local_path_specs
  <em class='small'>source</em>\"\n\n```python\n\n        def _get_local_path_specs(tool:
  str, project_home: Union[str, Path]) -> path_spec_type:\n            \"\"\"\n            Generate
  a list of standard pathspecs for local, project directory config files.\n\n            Args:\n
  \               tool (str): name of the tool to configure\n            \"\"\"\n
  \           return [\n                {\n                    \"path_specs\": Path(project_home)
  / f\"{tool}.ini\",\n                    \"ac_parser\": \"ini\",\n                    \"keys\":
  [tool],\n                },\n                {\n                    \"path_specs\":
  Path(project_home) / f\".{tool}\",\n                    \"ac_parser\": \"ini\",\n
  \                   \"keys\": [tool],\n                },\n                {\n                    \"path_specs\":
  Path(project_home) / f\".{tool}.ini\",\n                    \"ac_parser\": \"ini\",\n
  \                   \"keys\": [tool],\n                },\n                {\n                    \"path_specs\":
  Path(project_home) / f\"{tool}.yml\",\n                    \"ac_parser\": \"yaml\",\n
  \                   \"keys\": [tool],\n                },\n                {\n                    \"path_specs\":
  Path(project_home) / f\".{tool}.yml\",\n                    \"ac_parser\": \"yaml\",\n
  \                   \"keys\": [tool],\n                },\n                {\n                    \"path_specs\":
  Path(project_home) / f\"{tool}.toml\",\n                    \"ac_parser\": \"toml\",\n
  \                   \"keys\": [tool],\n                },\n                {\n                    \"path_specs\":
  Path(project_home) / f\".{tool}.toml\",\n                    \"ac_parser\": \"toml\",\n
  \                   \"keys\": [tool],\n                },\n                {\n                    \"path_specs\":
  Path(project_home) / \"pyproject.toml\",\n                    \"ac_parser\": \"toml\",\n
  \                   \"keys\": [\"tool\", tool],\n                },\n                {\n
  \                   \"path_specs\": Path(project_home) / \"setup.cfg\",\n                    \"ac_parser\":
  \"ini\",\n                    \"keys\": [f\"tool.{tool}\"],\n                },\n
  \           ]\n```\n\n\n!! function <h2 id='_get_attrs' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>_get_attrs <em class='small'>function</em></h2>\n
  \   Get nested config data from a list of keys.\n\n    specifically written for
  pyproject.toml which needs to get `tool` then `<tool>`\n???+ source \"_get_attrs
  <em class='small'>source</em>\"\n\n```python\n\n        def _get_attrs(attrs: list,
  config: Dict) -> Dict:\n            \"\"\"Get nested config data from a list of
  keys.\n\n            specifically written for pyproject.toml which needs to get
  `tool` then `<tool>`\n            \"\"\"\n            for attr in attrs:\n                config
  = config[attr]\n            return config\n```\n\n\n!! function <h2 id='_load_files'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_load_files <em class='small'>function</em></h2>\n
  \   Use anyconfig to load config files stopping at the first one that exists.\n\n
  \   config_path_specs (list): a list of pathspecs and keys to load\n???+ source
  \"_load_files <em class='small'>source</em>\"\n\n```python\n\n        def _load_files(config_path_specs:
  path_spec_type) -> Dict:\n            \"\"\"Use anyconfig to load config files stopping
  at the first one that exists.\n\n            config_path_specs (list): a list of
  pathspecs and keys to load\n            \"\"\"\n            for file in config_path_specs:\n
  \               if file[\"path_specs\"].exists():\n                    config =
  anyconfig.load(**file)\n                else:\n                    # ignore missing
  files\n                    continue\n\n                try:\n                    return
  _get_attrs(file[\"keys\"], config)\n                except KeyError:\n                    #
  ignore incorrect keys\n                    continue\n\n            return {}\n```\n\n\n!!
  function <h2 id='_load_env' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>_load_env <em class='small'>function</em></h2>\n    Load config from environment
  variables.\n\n    Args:\n        tool (str): name of the tool to configure\n???+
  source \"_load_env <em class='small'>source</em>\"\n\n```python\n\n        def _load_env(tool:
  str) -> Dict:\n            \"\"\"Load config from environment variables.\n\n            Args:\n
  \               tool (str): name of the tool to configure\n            \"\"\"\n
  \           vars = [var for var in os.environ if var.startswith(tool.upper())]\n
  \           return {\n                var.lower().strip(tool.lower()).strip(\"_\").strip(\"-\"):
  os.environ[var]\n                for var in vars\n            }\n```\n\n\n!! function
  <h2 id='load' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>load
  <em class='small'>function</em></h2>\n    Load tool config from standard config
  files.\n\n    Resolution Order\n\n    * First global file with a tool key\n    *
  First local file with a tool key\n    * Environment variables prefixed with `TOOL`\n
  \   * Overrides\n\n    Args:\n        tool (str): name of the tool to configure\n???+
  source \"load <em class='small'>source</em>\"\n\n```python\n\n        def load(tool:
  str, project_home: Union[Path, str] = \".\", overrides: Dict = {}) -> Dict:\n            \"\"\"Load
  tool config from standard config files.\n\n            Resolution Order\n\n            *
  First global file with a tool key\n            * First local file with a tool key\n
  \           * Environment variables prefixed with `TOOL`\n            * Overrides\n\n
  \           Args:\n                tool (str): name of the tool to configure\n            \"\"\"\n
  \           global_config = _load_files(_get_global_path_specs(tool))\n            local_config
  = _load_files(_get_local_path_specs(tool, project_home))\n            env_config
  = _load_env(tool)\n            return {**global_config, **local_config, **env_config,
  **overrides}\n```\n\n"
date: 0001-01-01
description: Standard Config. Inspired from frustrations that some tools have a tool.ini,
  .tool.ini, This file is for any project that can be configured in plain text such
  a
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Standard_Config.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Standard Config. Inspired from frustrations
    that some tools have a tool.ini, .tool.ini, This file is for any project that
    can be configured in plain text such a\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Standard_Config.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Standard Config. Inspired from frustrations
    that some tools have a tool.ini, .tool.ini, This file is for any project that
    can be configured in plain text such a\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<div class='container flex flex-row min-h-screen'>\n    <div>\n    </div>\n
    \   <div class='flex-grow px-8 mx-auto min-h-screen'>\n<header class='flex justify-center
    items-center p-8'>\n\n    <nav class='flex justify-center items-center my-8'>\n
    \       <a\n            href='/'>markata</a>\n        <a\n            href='https://github.com/WaylonWalker/markata'>GitHub</a>\n
    \       <a\n            href='https://markata.dev/docs/'>docs</a>\n        <a\n
    \           href='https://markata.dev/plugins/'>plugins</a>\n    </nav>\n\n    <div>\n
    \       <label id=\"theme-switch\" class=\"theme-switch\" for=\"checkbox-theme\"
    title=\"light/dark mode toggle\">\n            <input type=\"checkbox\" id=\"checkbox-theme\"
    />\n            <div class=\"slider round\"></div>\n        </label>\n    </div>\n</header><article
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Standard_Config.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <p>Standard Config.\nA
    module to load tooling config from a users project space.</p>\n<p>Inspired from
    frustrations that some tools have a tool.ini, .tool.ini,\nsetup.cfg, or pyproject.toml.
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"kn\">from</span>
    <span class=\"nn\">standard_config</span> <span class=\"kn\">import</span> <span
    class=\"n\">load</span>\n\n<span class=\"c1\"># Retrieve any overrides from the
    user</span>\n<span class=\"n\">overrides</span> <span class=\"o\">=</span> <span
    class=\"p\">{</span><span class=\"s1\">&#39;setting&#39;</span><span class=\"p\">:</span>
    <span class=\"kc\">True</span><span class=\"p\">}</span>\n<span class=\"n\">config</span>
    <span class=\"o\">=</span> <span class=\"n\">load</span><span class=\"p\">(</span><span
    class=\"s1\">&#39;my_tool&#39;</span><span class=\"p\">,</span> <span class=\"n\">overrides</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<h2 id=\"resolution-order\">Resolution
    Order <a class=\"header-anchor\" href=\"#resolution-order\"><svg class=\"heading-permalink\"
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
    files to consider</h3>\n<ul>\n<li>&lt;project_home&gt;/tool.ini</li>\n<li>&lt;project_home&gt;/.tool</li>\n<li>&lt;project_home&gt;/.tool.ini</li>\n<li>&lt;project_home&gt;/pyproject.toml</li>\n<li>&lt;project_home&gt;/setup.cfg</li>\n</ul>\n<p>!!
    function <h2 id='_get_global_path_specs' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>_get_global_path_specs <em class='small'>function</em></h2>\nGenerate a
    list of standard pathspecs for global config files.</p>\n<pre><code>Args:\n    tool
    (str): name of the tool to configure\n</code></pre>\n<div class=\"admonition source
    is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_get_global_path_specs
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">_get_global_path_specs</span><span class=\"p\">(</span><span
    class=\"n\">tool</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">path_spec_type</span><span
    class=\"p\">:</span>\n<span class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            Generate a list of standard pathspecs for global config
    files.</span>\n\n<span class=\"sd\">            Args:</span>\n<span class=\"sd\">
    \               tool (str): name of the tool to configure</span>\n<span class=\"sd\">
    \           &quot;&quot;&quot;</span>\n            <span class=\"k\">try</span><span
    class=\"p\">:</span>\n                <span class=\"n\">home</span> <span class=\"o\">=</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">os</span><span
    class=\"o\">.</span><span class=\"n\">environ</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;XDG_HOME&quot;</span><span class=\"p\">])</span>\n            <span
    class=\"k\">except</span> <span class=\"ne\">KeyError</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">home</span> <span class=\"o\">=</span> <span
    class=\"n\">Path</span><span class=\"o\">.</span><span class=\"n\">home</span><span
    class=\"p\">()</span>\n\n            <span class=\"k\">return</span> <span class=\"p\">[</span>\n
    \               <span class=\"p\">{</span><span class=\"s2\">&quot;path_specs&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">home</span> <span class=\"o\">/</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">tool</span><span class=\"si\">}</span><span class=\"s2\">.ini&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;ac_parser&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;ini&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;keys&quot;</span><span
    class=\"p\">:</span> <span class=\"p\">[</span><span class=\"n\">tool</span><span
    class=\"p\">]},</span>\n                <span class=\"p\">{</span><span class=\"s2\">&quot;path_specs&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">home</span> <span class=\"o\">/</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;.</span><span class=\"si\">{</span><span
    class=\"n\">tool</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;ac_parser&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;ini&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;keys&quot;</span><span
    class=\"p\">:</span> <span class=\"p\">[</span><span class=\"n\">tool</span><span
    class=\"p\">]},</span>\n                <span class=\"p\">{</span><span class=\"s2\">&quot;path_specs&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">home</span> <span class=\"o\">/</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;.</span><span class=\"si\">{</span><span
    class=\"n\">tool</span><span class=\"si\">}</span><span class=\"s2\">.ini&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;ac_parser&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;ini&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;keys&quot;</span><span
    class=\"p\">:</span> <span class=\"p\">[</span><span class=\"n\">tool</span><span
    class=\"p\">]},</span>\n                <span class=\"p\">{</span>\n                    <span
    class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">:</span> <span class=\"n\">home</span>
    <span class=\"o\">/</span> <span class=\"s2\">&quot;.config&quot;</span> <span
    class=\"o\">/</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">tool</span><span class=\"si\">}</span><span
    class=\"s2\">.ini&quot;</span><span class=\"p\">,</span>\n                    <span
    class=\"s2\">&quot;ac_parser&quot;</span><span class=\"p\">:</span> <span class=\"s2\">&quot;ini&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;keys&quot;</span><span
    class=\"p\">:</span> <span class=\"p\">[</span><span class=\"n\">tool</span><span
    class=\"p\">],</span>\n                <span class=\"p\">},</span>\n                <span
    class=\"p\">{</span>\n                    <span class=\"s2\">&quot;path_specs&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">home</span> <span class=\"o\">/</span>
    <span class=\"s2\">&quot;.config&quot;</span> <span class=\"o\">/</span> <span
    class=\"sa\">f</span><span class=\"s2\">&quot;.</span><span class=\"si\">{</span><span
    class=\"n\">tool</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;ac_parser&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;ini&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span>
    <span class=\"p\">[</span><span class=\"n\">tool</span><span class=\"p\">],</span>\n
    \               <span class=\"p\">},</span>\n                <span class=\"p\">{</span>\n
    \                   <span class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">home</span> <span class=\"o\">/</span> <span class=\"s2\">&quot;.config&quot;</span>
    <span class=\"o\">/</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;.</span><span
    class=\"si\">{</span><span class=\"n\">tool</span><span class=\"si\">}</span><span
    class=\"s2\">.ini&quot;</span><span class=\"p\">,</span>\n                    <span
    class=\"s2\">&quot;ac_parser&quot;</span><span class=\"p\">:</span> <span class=\"s2\">&quot;ini&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;keys&quot;</span><span
    class=\"p\">:</span> <span class=\"p\">[</span><span class=\"n\">tool</span><span
    class=\"p\">],</span>\n                <span class=\"p\">},</span>\n            <span
    class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='_get_local_path_specs'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_get_local_path_specs
    <em class='small'>function</em></h2>\nGenerate a list of standard pathspecs for
    local, project directory config files.</p>\n<pre><code>Args:\n    tool (str):
    name of the tool to configure\n</code></pre>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">_get_local_path_specs <em class='small'>source</em></p>\n</div>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">_get_local_path_specs</span><span class=\"p\">(</span><span
    class=\"n\">tool</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"n\">project_home</span><span class=\"p\">:</span>
    <span class=\"n\">Union</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"n\">Path</span><span class=\"p\">])</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">path_spec_type</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            Generate a list of standard pathspecs for local, project
    directory config files.</span>\n\n<span class=\"sd\">            Args:</span>\n<span
    class=\"sd\">                tool (str): name of the tool to configure</span>\n<span
    class=\"sd\">            &quot;&quot;&quot;</span>\n            <span class=\"k\">return</span>
    <span class=\"p\">[</span>\n                <span class=\"p\">{</span>\n                    <span
    class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">:</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">project_home</span><span class=\"p\">)</span>
    <span class=\"o\">/</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">tool</span><span class=\"si\">}</span><span
    class=\"s2\">.ini&quot;</span><span class=\"p\">,</span>\n                    <span
    class=\"s2\">&quot;ac_parser&quot;</span><span class=\"p\">:</span> <span class=\"s2\">&quot;ini&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;keys&quot;</span><span
    class=\"p\">:</span> <span class=\"p\">[</span><span class=\"n\">tool</span><span
    class=\"p\">],</span>\n                <span class=\"p\">},</span>\n                <span
    class=\"p\">{</span>\n                    <span class=\"s2\">&quot;path_specs&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">project_home</span><span class=\"p\">)</span> <span class=\"o\">/</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;.</span><span class=\"si\">{</span><span
    class=\"n\">tool</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;ac_parser&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;ini&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span>
    <span class=\"p\">[</span><span class=\"n\">tool</span><span class=\"p\">],</span>\n
    \               <span class=\"p\">},</span>\n                <span class=\"p\">{</span>\n
    \                   <span class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">project_home</span><span
    class=\"p\">)</span> <span class=\"o\">/</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;.</span><span class=\"si\">{</span><span class=\"n\">tool</span><span
    class=\"si\">}</span><span class=\"s2\">.ini&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;ac_parser&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;ini&quot;</span><span class=\"p\">,</span>\n                    <span
    class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span> <span class=\"p\">[</span><span
    class=\"n\">tool</span><span class=\"p\">],</span>\n                <span class=\"p\">},</span>\n
    \               <span class=\"p\">{</span>\n                    <span class=\"s2\">&quot;path_specs&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">project_home</span><span class=\"p\">)</span> <span class=\"o\">/</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">tool</span><span class=\"si\">}</span><span class=\"s2\">.yml&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;ac_parser&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;yaml&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span>
    <span class=\"p\">[</span><span class=\"n\">tool</span><span class=\"p\">],</span>\n
    \               <span class=\"p\">},</span>\n                <span class=\"p\">{</span>\n
    \                   <span class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">project_home</span><span
    class=\"p\">)</span> <span class=\"o\">/</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;.</span><span class=\"si\">{</span><span class=\"n\">tool</span><span
    class=\"si\">}</span><span class=\"s2\">.yml&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;ac_parser&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;yaml&quot;</span><span class=\"p\">,</span>\n                    <span
    class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span> <span class=\"p\">[</span><span
    class=\"n\">tool</span><span class=\"p\">],</span>\n                <span class=\"p\">},</span>\n
    \               <span class=\"p\">{</span>\n                    <span class=\"s2\">&quot;path_specs&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">project_home</span><span class=\"p\">)</span> <span class=\"o\">/</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">tool</span><span class=\"si\">}</span><span class=\"s2\">.toml&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;ac_parser&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;toml&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span>
    <span class=\"p\">[</span><span class=\"n\">tool</span><span class=\"p\">],</span>\n
    \               <span class=\"p\">},</span>\n                <span class=\"p\">{</span>\n
    \                   <span class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">project_home</span><span
    class=\"p\">)</span> <span class=\"o\">/</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;.</span><span class=\"si\">{</span><span class=\"n\">tool</span><span
    class=\"si\">}</span><span class=\"s2\">.toml&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;ac_parser&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;toml&quot;</span><span class=\"p\">,</span>\n                    <span
    class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span> <span class=\"p\">[</span><span
    class=\"n\">tool</span><span class=\"p\">],</span>\n                <span class=\"p\">},</span>\n
    \               <span class=\"p\">{</span>\n                    <span class=\"s2\">&quot;path_specs&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">project_home</span><span class=\"p\">)</span> <span class=\"o\">/</span>
    <span class=\"s2\">&quot;pyproject.toml&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;ac_parser&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;toml&quot;</span><span class=\"p\">,</span>\n                    <span
    class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span> <span class=\"p\">[</span><span
    class=\"s2\">&quot;tool&quot;</span><span class=\"p\">,</span> <span class=\"n\">tool</span><span
    class=\"p\">],</span>\n                <span class=\"p\">},</span>\n                <span
    class=\"p\">{</span>\n                    <span class=\"s2\">&quot;path_specs&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">project_home</span><span class=\"p\">)</span> <span class=\"o\">/</span>
    <span class=\"s2\">&quot;setup.cfg&quot;</span><span class=\"p\">,</span>\n                    <span
    class=\"s2\">&quot;ac_parser&quot;</span><span class=\"p\">:</span> <span class=\"s2\">&quot;ini&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;keys&quot;</span><span
    class=\"p\">:</span> <span class=\"p\">[</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;tool.</span><span class=\"si\">{</span><span class=\"n\">tool</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">],</span>\n
    \               <span class=\"p\">},</span>\n            <span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='_get_attrs' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>_get_attrs <em class='small'>function</em></h2>\nGet nested config data
    from a list of keys.</p>\n<pre><code>specifically written for pyproject.toml which
    needs to get `tool` then `&lt;tool&gt;`\n</code></pre>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_get_attrs
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">_get_attrs</span><span class=\"p\">(</span><span class=\"n\">attrs</span><span
    class=\"p\">:</span> <span class=\"nb\">list</span><span class=\"p\">,</span>
    <span class=\"n\">config</span><span class=\"p\">:</span> <span class=\"n\">Dict</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Dict</span><span
    class=\"p\">:</span>\n<span class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;Get
    nested config data from a list of keys.</span>\n\n<span class=\"sd\">            specifically
    written for pyproject.toml which needs to get `tool` then `&lt;tool&gt;`</span>\n<span
    class=\"sd\">            &quot;&quot;&quot;</span>\n            <span class=\"k\">for</span>
    <span class=\"n\">attr</span> <span class=\"ow\">in</span> <span class=\"n\">attrs</span><span
    class=\"p\">:</span>\n                <span class=\"n\">config</span> <span class=\"o\">=</span>
    <span class=\"n\">config</span><span class=\"p\">[</span><span class=\"n\">attr</span><span
    class=\"p\">]</span>\n            <span class=\"k\">return</span> <span class=\"n\">config</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='_load_files' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>_load_files <em class='small'>function</em></h2>\nUse anyconfig to load
    config files stopping at the first one that exists.</p>\n<pre><code>config_path_specs
    (list): a list of pathspecs and keys to load\n</code></pre>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_load_files
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">_load_files</span><span class=\"p\">(</span><span class=\"n\">config_path_specs</span><span
    class=\"p\">:</span> <span class=\"n\">path_spec_type</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">Dict</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;Use anyconfig
    to load config files stopping at the first one that exists.</span>\n\n<span class=\"sd\">
    \           config_path_specs (list): a list of pathspecs and keys to load</span>\n<span
    class=\"sd\">            &quot;&quot;&quot;</span>\n            <span class=\"k\">for</span>
    <span class=\"n\">file</span> <span class=\"ow\">in</span> <span class=\"n\">config_path_specs</span><span
    class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"n\">file</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">():</span>\n
    \                   <span class=\"n\">config</span> <span class=\"o\">=</span>
    <span class=\"n\">anyconfig</span><span class=\"o\">.</span><span class=\"n\">load</span><span
    class=\"p\">(</span><span class=\"o\">**</span><span class=\"n\">file</span><span
    class=\"p\">)</span>\n                <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \                   <span class=\"c1\"># ignore missing files</span>\n                    <span
    class=\"k\">continue</span>\n\n                <span class=\"k\">try</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">_get_attrs</span><span class=\"p\">(</span><span class=\"n\">file</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;keys&quot;</span><span class=\"p\">],</span>
    <span class=\"n\">config</span><span class=\"p\">)</span>\n                <span
    class=\"k\">except</span> <span class=\"ne\">KeyError</span><span class=\"p\">:</span>\n
    \                   <span class=\"c1\"># ignore incorrect keys</span>\n                    <span
    class=\"k\">continue</span>\n\n            <span class=\"k\">return</span> <span
    class=\"p\">{}</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='_load_env'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_load_env <em class='small'>function</em></h2>\nLoad
    config from environment variables.</p>\n<pre><code>Args:\n    tool (str): name
    of the tool to configure\n</code></pre>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">_load_env <em class='small'>source</em></p>\n</div>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">_load_env</span><span class=\"p\">(</span><span class=\"n\">tool</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Dict</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;Load config
    from environment variables.</span>\n\n<span class=\"sd\">            Args:</span>\n<span
    class=\"sd\">                tool (str): name of the tool to configure</span>\n<span
    class=\"sd\">            &quot;&quot;&quot;</span>\n            <span class=\"nb\">vars</span>
    <span class=\"o\">=</span> <span class=\"p\">[</span><span class=\"n\">var</span>
    <span class=\"k\">for</span> <span class=\"n\">var</span> <span class=\"ow\">in</span>
    <span class=\"n\">os</span><span class=\"o\">.</span><span class=\"n\">environ</span>
    <span class=\"k\">if</span> <span class=\"n\">var</span><span class=\"o\">.</span><span
    class=\"n\">startswith</span><span class=\"p\">(</span><span class=\"n\">tool</span><span
    class=\"o\">.</span><span class=\"n\">upper</span><span class=\"p\">())]</span>\n
    \           <span class=\"k\">return</span> <span class=\"p\">{</span>\n                <span
    class=\"n\">var</span><span class=\"o\">.</span><span class=\"n\">lower</span><span
    class=\"p\">()</span><span class=\"o\">.</span><span class=\"n\">strip</span><span
    class=\"p\">(</span><span class=\"n\">tool</span><span class=\"o\">.</span><span
    class=\"n\">lower</span><span class=\"p\">())</span><span class=\"o\">.</span><span
    class=\"n\">strip</span><span class=\"p\">(</span><span class=\"s2\">&quot;_&quot;</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">strip</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;-&quot;</span><span class=\"p\">):</span>
    <span class=\"n\">os</span><span class=\"o\">.</span><span class=\"n\">environ</span><span
    class=\"p\">[</span><span class=\"n\">var</span><span class=\"p\">]</span>\n                <span
    class=\"k\">for</span> <span class=\"n\">var</span> <span class=\"ow\">in</span>
    <span class=\"nb\">vars</span>\n            <span class=\"p\">}</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='load' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>load <em class='small'>function</em></h2>\nLoad tool config from standard
    config files.</p>\n<pre><code>Resolution Order\n\n* First global file with a tool
    key\n* First local file with a tool key\n* Environment variables prefixed with
    `TOOL`\n* Overrides\n\nArgs:\n    tool (str): name of the tool to configure\n</code></pre>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">load
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">load</span><span class=\"p\">(</span><span class=\"n\">tool</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">,</span> <span
    class=\"n\">project_home</span><span class=\"p\">:</span> <span class=\"n\">Union</span><span
    class=\"p\">[</span><span class=\"n\">Path</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"s2\">&quot;.&quot;</span><span class=\"p\">,</span> <span class=\"n\">overrides</span><span
    class=\"p\">:</span> <span class=\"n\">Dict</span> <span class=\"o\">=</span>
    <span class=\"p\">{})</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Dict</span><span
    class=\"p\">:</span>\n<span class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;Load
    tool config from standard config files.</span>\n\n<span class=\"sd\">            Resolution
    Order</span>\n\n<span class=\"sd\">            * First global file with a tool
    key</span>\n<span class=\"sd\">            * First local file with a tool key</span>\n<span
    class=\"sd\">            * Environment variables prefixed with `TOOL`</span>\n<span
    class=\"sd\">            * Overrides</span>\n\n<span class=\"sd\">            Args:</span>\n<span
    class=\"sd\">                tool (str): name of the tool to configure</span>\n<span
    class=\"sd\">            &quot;&quot;&quot;</span>\n            <span class=\"n\">global_config</span>
    <span class=\"o\">=</span> <span class=\"n\">_load_files</span><span class=\"p\">(</span><span
    class=\"n\">_get_global_path_specs</span><span class=\"p\">(</span><span class=\"n\">tool</span><span
    class=\"p\">))</span>\n            <span class=\"n\">local_config</span> <span
    class=\"o\">=</span> <span class=\"n\">_load_files</span><span class=\"p\">(</span><span
    class=\"n\">_get_local_path_specs</span><span class=\"p\">(</span><span class=\"n\">tool</span><span
    class=\"p\">,</span> <span class=\"n\">project_home</span><span class=\"p\">))</span>\n
    \           <span class=\"n\">env_config</span> <span class=\"o\">=</span> <span
    class=\"n\">_load_env</span><span class=\"p\">(</span><span class=\"n\">tool</span><span
    class=\"p\">)</span>\n            <span class=\"k\">return</span> <span class=\"p\">{</span><span
    class=\"o\">**</span><span class=\"n\">global_config</span><span class=\"p\">,</span>
    <span class=\"o\">**</span><span class=\"n\">local_config</span><span class=\"p\">,</span>
    <span class=\"o\">**</span><span class=\"n\">env_config</span><span class=\"p\">,</span>
    <span class=\"o\">**</span><span class=\"n\">overrides</span><span class=\"p\">}</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>    </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Standard_Config.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Standard Config. Inspired from frustrations
    that some tools have a tool.ini, .tool.ini, This file is for any project that
    can be configured in plain text such a\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Standard_Config.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Standard Config. Inspired from frustrations
    that some tools have a tool.ini, .tool.ini, This file is for any project that
    can be configured in plain text such a\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<article style=\"text-align: center;\">\n    <style>\n        section
    {\n            font-size: 200%;\n        }\n\n\n        .edit {\n            display:
    none;\n        }\n    </style>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Standard_Config.Py\n    </h1>\n</section></article>\n     </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Standard_Config.Py\n    </h1>\n</section>    <section class=\"body\">\n
    \       <p>Standard Config.\nA module to load tooling config from a users project
    space.</p>\n<p>Inspired from frustrations that some tools have a tool.ini, .tool.ini,\nsetup.cfg,
    or pyproject.toml.  Some allow for global configs, some don't.  Some\nproperly
    follow the users home directory, others end up in a weird temp\ndirectory.  Windows
    home directory is only more confusing.  Some will even\nrespect the users <code>$XDG_HOME</code>
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"kn\">from</span>
    <span class=\"nn\">standard_config</span> <span class=\"kn\">import</span> <span
    class=\"n\">load</span>\n\n<span class=\"c1\"># Retrieve any overrides from the
    user</span>\n<span class=\"n\">overrides</span> <span class=\"o\">=</span> <span
    class=\"p\">{</span><span class=\"s1\">&#39;setting&#39;</span><span class=\"p\">:</span>
    <span class=\"kc\">True</span><span class=\"p\">}</span>\n<span class=\"n\">config</span>
    <span class=\"o\">=</span> <span class=\"n\">load</span><span class=\"p\">(</span><span
    class=\"s1\">&#39;my_tool&#39;</span><span class=\"p\">,</span> <span class=\"n\">overrides</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<h2 id=\"resolution-order\">Resolution
    Order <a class=\"header-anchor\" href=\"#resolution-order\"><svg class=\"heading-permalink\"
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
    files to consider</h3>\n<ul>\n<li>&lt;project_home&gt;/tool.ini</li>\n<li>&lt;project_home&gt;/.tool</li>\n<li>&lt;project_home&gt;/.tool.ini</li>\n<li>&lt;project_home&gt;/pyproject.toml</li>\n<li>&lt;project_home&gt;/setup.cfg</li>\n</ul>\n<p>!!
    function <h2 id='_get_global_path_specs' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>_get_global_path_specs <em class='small'>function</em></h2>\nGenerate a
    list of standard pathspecs for global config files.</p>\n<pre><code>Args:\n    tool
    (str): name of the tool to configure\n</code></pre>\n<div class=\"admonition source
    is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_get_global_path_specs
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">_get_global_path_specs</span><span class=\"p\">(</span><span
    class=\"n\">tool</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">path_spec_type</span><span
    class=\"p\">:</span>\n<span class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            Generate a list of standard pathspecs for global config
    files.</span>\n\n<span class=\"sd\">            Args:</span>\n<span class=\"sd\">
    \               tool (str): name of the tool to configure</span>\n<span class=\"sd\">
    \           &quot;&quot;&quot;</span>\n            <span class=\"k\">try</span><span
    class=\"p\">:</span>\n                <span class=\"n\">home</span> <span class=\"o\">=</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">os</span><span
    class=\"o\">.</span><span class=\"n\">environ</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;XDG_HOME&quot;</span><span class=\"p\">])</span>\n            <span
    class=\"k\">except</span> <span class=\"ne\">KeyError</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">home</span> <span class=\"o\">=</span> <span
    class=\"n\">Path</span><span class=\"o\">.</span><span class=\"n\">home</span><span
    class=\"p\">()</span>\n\n            <span class=\"k\">return</span> <span class=\"p\">[</span>\n
    \               <span class=\"p\">{</span><span class=\"s2\">&quot;path_specs&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">home</span> <span class=\"o\">/</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">tool</span><span class=\"si\">}</span><span class=\"s2\">.ini&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;ac_parser&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;ini&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;keys&quot;</span><span
    class=\"p\">:</span> <span class=\"p\">[</span><span class=\"n\">tool</span><span
    class=\"p\">]},</span>\n                <span class=\"p\">{</span><span class=\"s2\">&quot;path_specs&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">home</span> <span class=\"o\">/</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;.</span><span class=\"si\">{</span><span
    class=\"n\">tool</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;ac_parser&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;ini&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;keys&quot;</span><span
    class=\"p\">:</span> <span class=\"p\">[</span><span class=\"n\">tool</span><span
    class=\"p\">]},</span>\n                <span class=\"p\">{</span><span class=\"s2\">&quot;path_specs&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">home</span> <span class=\"o\">/</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;.</span><span class=\"si\">{</span><span
    class=\"n\">tool</span><span class=\"si\">}</span><span class=\"s2\">.ini&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;ac_parser&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;ini&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;keys&quot;</span><span
    class=\"p\">:</span> <span class=\"p\">[</span><span class=\"n\">tool</span><span
    class=\"p\">]},</span>\n                <span class=\"p\">{</span>\n                    <span
    class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">:</span> <span class=\"n\">home</span>
    <span class=\"o\">/</span> <span class=\"s2\">&quot;.config&quot;</span> <span
    class=\"o\">/</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">tool</span><span class=\"si\">}</span><span
    class=\"s2\">.ini&quot;</span><span class=\"p\">,</span>\n                    <span
    class=\"s2\">&quot;ac_parser&quot;</span><span class=\"p\">:</span> <span class=\"s2\">&quot;ini&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;keys&quot;</span><span
    class=\"p\">:</span> <span class=\"p\">[</span><span class=\"n\">tool</span><span
    class=\"p\">],</span>\n                <span class=\"p\">},</span>\n                <span
    class=\"p\">{</span>\n                    <span class=\"s2\">&quot;path_specs&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">home</span> <span class=\"o\">/</span>
    <span class=\"s2\">&quot;.config&quot;</span> <span class=\"o\">/</span> <span
    class=\"sa\">f</span><span class=\"s2\">&quot;.</span><span class=\"si\">{</span><span
    class=\"n\">tool</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;ac_parser&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;ini&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span>
    <span class=\"p\">[</span><span class=\"n\">tool</span><span class=\"p\">],</span>\n
    \               <span class=\"p\">},</span>\n                <span class=\"p\">{</span>\n
    \                   <span class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">home</span> <span class=\"o\">/</span> <span class=\"s2\">&quot;.config&quot;</span>
    <span class=\"o\">/</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;.</span><span
    class=\"si\">{</span><span class=\"n\">tool</span><span class=\"si\">}</span><span
    class=\"s2\">.ini&quot;</span><span class=\"p\">,</span>\n                    <span
    class=\"s2\">&quot;ac_parser&quot;</span><span class=\"p\">:</span> <span class=\"s2\">&quot;ini&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;keys&quot;</span><span
    class=\"p\">:</span> <span class=\"p\">[</span><span class=\"n\">tool</span><span
    class=\"p\">],</span>\n                <span class=\"p\">},</span>\n            <span
    class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='_get_local_path_specs'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_get_local_path_specs
    <em class='small'>function</em></h2>\nGenerate a list of standard pathspecs for
    local, project directory config files.</p>\n<pre><code>Args:\n    tool (str):
    name of the tool to configure\n</code></pre>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">_get_local_path_specs <em class='small'>source</em></p>\n</div>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">_get_local_path_specs</span><span class=\"p\">(</span><span
    class=\"n\">tool</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"n\">project_home</span><span class=\"p\">:</span>
    <span class=\"n\">Union</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"n\">Path</span><span class=\"p\">])</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">path_spec_type</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            Generate a list of standard pathspecs for local, project
    directory config files.</span>\n\n<span class=\"sd\">            Args:</span>\n<span
    class=\"sd\">                tool (str): name of the tool to configure</span>\n<span
    class=\"sd\">            &quot;&quot;&quot;</span>\n            <span class=\"k\">return</span>
    <span class=\"p\">[</span>\n                <span class=\"p\">{</span>\n                    <span
    class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">:</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">project_home</span><span class=\"p\">)</span>
    <span class=\"o\">/</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">tool</span><span class=\"si\">}</span><span
    class=\"s2\">.ini&quot;</span><span class=\"p\">,</span>\n                    <span
    class=\"s2\">&quot;ac_parser&quot;</span><span class=\"p\">:</span> <span class=\"s2\">&quot;ini&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;keys&quot;</span><span
    class=\"p\">:</span> <span class=\"p\">[</span><span class=\"n\">tool</span><span
    class=\"p\">],</span>\n                <span class=\"p\">},</span>\n                <span
    class=\"p\">{</span>\n                    <span class=\"s2\">&quot;path_specs&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">project_home</span><span class=\"p\">)</span> <span class=\"o\">/</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;.</span><span class=\"si\">{</span><span
    class=\"n\">tool</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;ac_parser&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;ini&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span>
    <span class=\"p\">[</span><span class=\"n\">tool</span><span class=\"p\">],</span>\n
    \               <span class=\"p\">},</span>\n                <span class=\"p\">{</span>\n
    \                   <span class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">project_home</span><span
    class=\"p\">)</span> <span class=\"o\">/</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;.</span><span class=\"si\">{</span><span class=\"n\">tool</span><span
    class=\"si\">}</span><span class=\"s2\">.ini&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;ac_parser&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;ini&quot;</span><span class=\"p\">,</span>\n                    <span
    class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span> <span class=\"p\">[</span><span
    class=\"n\">tool</span><span class=\"p\">],</span>\n                <span class=\"p\">},</span>\n
    \               <span class=\"p\">{</span>\n                    <span class=\"s2\">&quot;path_specs&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">project_home</span><span class=\"p\">)</span> <span class=\"o\">/</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">tool</span><span class=\"si\">}</span><span class=\"s2\">.yml&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;ac_parser&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;yaml&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span>
    <span class=\"p\">[</span><span class=\"n\">tool</span><span class=\"p\">],</span>\n
    \               <span class=\"p\">},</span>\n                <span class=\"p\">{</span>\n
    \                   <span class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">project_home</span><span
    class=\"p\">)</span> <span class=\"o\">/</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;.</span><span class=\"si\">{</span><span class=\"n\">tool</span><span
    class=\"si\">}</span><span class=\"s2\">.yml&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;ac_parser&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;yaml&quot;</span><span class=\"p\">,</span>\n                    <span
    class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span> <span class=\"p\">[</span><span
    class=\"n\">tool</span><span class=\"p\">],</span>\n                <span class=\"p\">},</span>\n
    \               <span class=\"p\">{</span>\n                    <span class=\"s2\">&quot;path_specs&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">project_home</span><span class=\"p\">)</span> <span class=\"o\">/</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">tool</span><span class=\"si\">}</span><span class=\"s2\">.toml&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;ac_parser&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;toml&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span>
    <span class=\"p\">[</span><span class=\"n\">tool</span><span class=\"p\">],</span>\n
    \               <span class=\"p\">},</span>\n                <span class=\"p\">{</span>\n
    \                   <span class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">project_home</span><span
    class=\"p\">)</span> <span class=\"o\">/</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;.</span><span class=\"si\">{</span><span class=\"n\">tool</span><span
    class=\"si\">}</span><span class=\"s2\">.toml&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;ac_parser&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;toml&quot;</span><span class=\"p\">,</span>\n                    <span
    class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span> <span class=\"p\">[</span><span
    class=\"n\">tool</span><span class=\"p\">],</span>\n                <span class=\"p\">},</span>\n
    \               <span class=\"p\">{</span>\n                    <span class=\"s2\">&quot;path_specs&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">project_home</span><span class=\"p\">)</span> <span class=\"o\">/</span>
    <span class=\"s2\">&quot;pyproject.toml&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;ac_parser&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;toml&quot;</span><span class=\"p\">,</span>\n                    <span
    class=\"s2\">&quot;keys&quot;</span><span class=\"p\">:</span> <span class=\"p\">[</span><span
    class=\"s2\">&quot;tool&quot;</span><span class=\"p\">,</span> <span class=\"n\">tool</span><span
    class=\"p\">],</span>\n                <span class=\"p\">},</span>\n                <span
    class=\"p\">{</span>\n                    <span class=\"s2\">&quot;path_specs&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">project_home</span><span class=\"p\">)</span> <span class=\"o\">/</span>
    <span class=\"s2\">&quot;setup.cfg&quot;</span><span class=\"p\">,</span>\n                    <span
    class=\"s2\">&quot;ac_parser&quot;</span><span class=\"p\">:</span> <span class=\"s2\">&quot;ini&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;keys&quot;</span><span
    class=\"p\">:</span> <span class=\"p\">[</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;tool.</span><span class=\"si\">{</span><span class=\"n\">tool</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">],</span>\n
    \               <span class=\"p\">},</span>\n            <span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='_get_attrs' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>_get_attrs <em class='small'>function</em></h2>\nGet nested config data
    from a list of keys.</p>\n<pre><code>specifically written for pyproject.toml which
    needs to get `tool` then `&lt;tool&gt;`\n</code></pre>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_get_attrs
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">_get_attrs</span><span class=\"p\">(</span><span class=\"n\">attrs</span><span
    class=\"p\">:</span> <span class=\"nb\">list</span><span class=\"p\">,</span>
    <span class=\"n\">config</span><span class=\"p\">:</span> <span class=\"n\">Dict</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Dict</span><span
    class=\"p\">:</span>\n<span class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;Get
    nested config data from a list of keys.</span>\n\n<span class=\"sd\">            specifically
    written for pyproject.toml which needs to get `tool` then `&lt;tool&gt;`</span>\n<span
    class=\"sd\">            &quot;&quot;&quot;</span>\n            <span class=\"k\">for</span>
    <span class=\"n\">attr</span> <span class=\"ow\">in</span> <span class=\"n\">attrs</span><span
    class=\"p\">:</span>\n                <span class=\"n\">config</span> <span class=\"o\">=</span>
    <span class=\"n\">config</span><span class=\"p\">[</span><span class=\"n\">attr</span><span
    class=\"p\">]</span>\n            <span class=\"k\">return</span> <span class=\"n\">config</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='_load_files' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>_load_files <em class='small'>function</em></h2>\nUse anyconfig to load
    config files stopping at the first one that exists.</p>\n<pre><code>config_path_specs
    (list): a list of pathspecs and keys to load\n</code></pre>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_load_files
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">_load_files</span><span class=\"p\">(</span><span class=\"n\">config_path_specs</span><span
    class=\"p\">:</span> <span class=\"n\">path_spec_type</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">Dict</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;Use anyconfig
    to load config files stopping at the first one that exists.</span>\n\n<span class=\"sd\">
    \           config_path_specs (list): a list of pathspecs and keys to load</span>\n<span
    class=\"sd\">            &quot;&quot;&quot;</span>\n            <span class=\"k\">for</span>
    <span class=\"n\">file</span> <span class=\"ow\">in</span> <span class=\"n\">config_path_specs</span><span
    class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"n\">file</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;path_specs&quot;</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">():</span>\n
    \                   <span class=\"n\">config</span> <span class=\"o\">=</span>
    <span class=\"n\">anyconfig</span><span class=\"o\">.</span><span class=\"n\">load</span><span
    class=\"p\">(</span><span class=\"o\">**</span><span class=\"n\">file</span><span
    class=\"p\">)</span>\n                <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \                   <span class=\"c1\"># ignore missing files</span>\n                    <span
    class=\"k\">continue</span>\n\n                <span class=\"k\">try</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">_get_attrs</span><span class=\"p\">(</span><span class=\"n\">file</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;keys&quot;</span><span class=\"p\">],</span>
    <span class=\"n\">config</span><span class=\"p\">)</span>\n                <span
    class=\"k\">except</span> <span class=\"ne\">KeyError</span><span class=\"p\">:</span>\n
    \                   <span class=\"c1\"># ignore incorrect keys</span>\n                    <span
    class=\"k\">continue</span>\n\n            <span class=\"k\">return</span> <span
    class=\"p\">{}</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='_load_env'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_load_env <em class='small'>function</em></h2>\nLoad
    config from environment variables.</p>\n<pre><code>Args:\n    tool (str): name
    of the tool to configure\n</code></pre>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">_load_env <em class='small'>source</em></p>\n</div>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">_load_env</span><span class=\"p\">(</span><span class=\"n\">tool</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Dict</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;Load config
    from environment variables.</span>\n\n<span class=\"sd\">            Args:</span>\n<span
    class=\"sd\">                tool (str): name of the tool to configure</span>\n<span
    class=\"sd\">            &quot;&quot;&quot;</span>\n            <span class=\"nb\">vars</span>
    <span class=\"o\">=</span> <span class=\"p\">[</span><span class=\"n\">var</span>
    <span class=\"k\">for</span> <span class=\"n\">var</span> <span class=\"ow\">in</span>
    <span class=\"n\">os</span><span class=\"o\">.</span><span class=\"n\">environ</span>
    <span class=\"k\">if</span> <span class=\"n\">var</span><span class=\"o\">.</span><span
    class=\"n\">startswith</span><span class=\"p\">(</span><span class=\"n\">tool</span><span
    class=\"o\">.</span><span class=\"n\">upper</span><span class=\"p\">())]</span>\n
    \           <span class=\"k\">return</span> <span class=\"p\">{</span>\n                <span
    class=\"n\">var</span><span class=\"o\">.</span><span class=\"n\">lower</span><span
    class=\"p\">()</span><span class=\"o\">.</span><span class=\"n\">strip</span><span
    class=\"p\">(</span><span class=\"n\">tool</span><span class=\"o\">.</span><span
    class=\"n\">lower</span><span class=\"p\">())</span><span class=\"o\">.</span><span
    class=\"n\">strip</span><span class=\"p\">(</span><span class=\"s2\">&quot;_&quot;</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">strip</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;-&quot;</span><span class=\"p\">):</span>
    <span class=\"n\">os</span><span class=\"o\">.</span><span class=\"n\">environ</span><span
    class=\"p\">[</span><span class=\"n\">var</span><span class=\"p\">]</span>\n                <span
    class=\"k\">for</span> <span class=\"n\">var</span> <span class=\"ow\">in</span>
    <span class=\"nb\">vars</span>\n            <span class=\"p\">}</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='load' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>load <em class='small'>function</em></h2>\nLoad tool config from standard
    config files.</p>\n<pre><code>Resolution Order\n\n* First global file with a tool
    key\n* First local file with a tool key\n* Environment variables prefixed with
    `TOOL`\n* Overrides\n\nArgs:\n    tool (str): name of the tool to configure\n</code></pre>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">load
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">load</span><span class=\"p\">(</span><span class=\"n\">tool</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">,</span> <span
    class=\"n\">project_home</span><span class=\"p\">:</span> <span class=\"n\">Union</span><span
    class=\"p\">[</span><span class=\"n\">Path</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"s2\">&quot;.&quot;</span><span class=\"p\">,</span> <span class=\"n\">overrides</span><span
    class=\"p\">:</span> <span class=\"n\">Dict</span> <span class=\"o\">=</span>
    <span class=\"p\">{})</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Dict</span><span
    class=\"p\">:</span>\n<span class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;Load
    tool config from standard config files.</span>\n\n<span class=\"sd\">            Resolution
    Order</span>\n\n<span class=\"sd\">            * First global file with a tool
    key</span>\n<span class=\"sd\">            * First local file with a tool key</span>\n<span
    class=\"sd\">            * Environment variables prefixed with `TOOL`</span>\n<span
    class=\"sd\">            * Overrides</span>\n\n<span class=\"sd\">            Args:</span>\n<span
    class=\"sd\">                tool (str): name of the tool to configure</span>\n<span
    class=\"sd\">            &quot;&quot;&quot;</span>\n            <span class=\"n\">global_config</span>
    <span class=\"o\">=</span> <span class=\"n\">_load_files</span><span class=\"p\">(</span><span
    class=\"n\">_get_global_path_specs</span><span class=\"p\">(</span><span class=\"n\">tool</span><span
    class=\"p\">))</span>\n            <span class=\"n\">local_config</span> <span
    class=\"o\">=</span> <span class=\"n\">_load_files</span><span class=\"p\">(</span><span
    class=\"n\">_get_local_path_specs</span><span class=\"p\">(</span><span class=\"n\">tool</span><span
    class=\"p\">,</span> <span class=\"n\">project_home</span><span class=\"p\">))</span>\n
    \           <span class=\"n\">env_config</span> <span class=\"o\">=</span> <span
    class=\"n\">_load_env</span><span class=\"p\">(</span><span class=\"n\">tool</span><span
    class=\"p\">)</span>\n            <span class=\"k\">return</span> <span class=\"p\">{</span><span
    class=\"o\">**</span><span class=\"n\">global_config</span><span class=\"p\">,</span>
    <span class=\"o\">**</span><span class=\"n\">local_config</span><span class=\"p\">,</span>
    <span class=\"o\">**</span><span class=\"n\">env_config</span><span class=\"p\">,</span>
    <span class=\"o\">**</span><span class=\"n\">overrides</span><span class=\"p\">}</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>"
  raw.md: ''
published: true
slug: markata/standard-config
title: Standard_Config.Py


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


!! function <h2 id='_get_global_path_specs' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_get_global_path_specs <em class='small'>function</em></h2>
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
                {"path_specs": home / f"{tool}.ini", "ac_parser": "ini", "keys": [tool]},
                {"path_specs": home / f".{tool}", "ac_parser": "ini", "keys": [tool]},
                {"path_specs": home / f".{tool}.ini", "ac_parser": "ini", "keys": [tool]},
                {
                    "path_specs": home / ".config" / f"{tool}.ini",
                    "ac_parser": "ini",
                    "keys": [tool],
                },
                {
                    "path_specs": home / ".config" / f".{tool}",
                    "ac_parser": "ini",
                    "keys": [tool],
                },
                {
                    "path_specs": home / ".config" / f".{tool}.ini",
                    "ac_parser": "ini",
                    "keys": [tool],
                },
            ]
```


!! function <h2 id='_get_local_path_specs' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_get_local_path_specs <em class='small'>function</em></h2>
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
                    "ac_parser": "ini",
                    "keys": [tool],
                },
                {
                    "path_specs": Path(project_home) / f".{tool}",
                    "ac_parser": "ini",
                    "keys": [tool],
                },
                {
                    "path_specs": Path(project_home) / f".{tool}.ini",
                    "ac_parser": "ini",
                    "keys": [tool],
                },
                {
                    "path_specs": Path(project_home) / f"{tool}.yml",
                    "ac_parser": "yaml",
                    "keys": [tool],
                },
                {
                    "path_specs": Path(project_home) / f".{tool}.yml",
                    "ac_parser": "yaml",
                    "keys": [tool],
                },
                {
                    "path_specs": Path(project_home) / f"{tool}.toml",
                    "ac_parser": "toml",
                    "keys": [tool],
                },
                {
                    "path_specs": Path(project_home) / f".{tool}.toml",
                    "ac_parser": "toml",
                    "keys": [tool],
                },
                {
                    "path_specs": Path(project_home) / "pyproject.toml",
                    "ac_parser": "toml",
                    "keys": ["tool", tool],
                },
                {
                    "path_specs": Path(project_home) / "setup.cfg",
                    "ac_parser": "ini",
                    "keys": [f"tool.{tool}"],
                },
            ]
```


!! function <h2 id='_get_attrs' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_get_attrs <em class='small'>function</em></h2>
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


!! function <h2 id='_load_files' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_load_files <em class='small'>function</em></h2>
    Use anyconfig to load config files stopping at the first one that exists.

    config_path_specs (list): a list of pathspecs and keys to load
???+ source "_load_files <em class='small'>source</em>"

```python

        def _load_files(config_path_specs: path_spec_type) -> Dict:
            """Use anyconfig to load config files stopping at the first one that exists.

            config_path_specs (list): a list of pathspecs and keys to load
            """
            for file in config_path_specs:
                if file["path_specs"].exists():
                    config = anyconfig.load(**file)
                else:
                    # ignore missing files
                    continue

                try:
                    return _get_attrs(file["keys"], config)
                except KeyError:
                    # ignore incorrect keys
                    continue

            return {}
```


!! function <h2 id='_load_env' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_load_env <em class='small'>function</em></h2>
    Load config from environment variables.

    Args:
        tool (str): name of the tool to configure
???+ source "_load_env <em class='small'>source</em>"

```python

        def _load_env(tool: str) -> Dict:
            """Load config from environment variables.

            Args:
                tool (str): name of the tool to configure
            """
            vars = [var for var in os.environ if var.startswith(tool.upper())]
            return {
                var.lower().strip(tool.lower()).strip("_").strip("-"): os.environ[var]
                for var in vars
            }
```


!! function <h2 id='load' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>load <em class='small'>function</em></h2>
    Load tool config from standard config files.

    Resolution Order

    * First global file with a tool key
    * First local file with a tool key
    * Environment variables prefixed with `TOOL`
    * Overrides

    Args:
        tool (str): name of the tool to configure
???+ source "load <em class='small'>source</em>"

```python

        def load(tool: str, project_home: Union[Path, str] = ".", overrides: Dict = {}) -> Dict:
            """Load tool config from standard config files.

            Resolution Order

            * First global file with a tool key
            * First local file with a tool key
            * Environment variables prefixed with `TOOL`
            * Overrides

            Args:
                tool (str): name of the tool to configure
            """
            global_config = _load_files(_get_global_path_specs(tool))
            local_config = _load_files(_get_local_path_specs(tool, project_home))
            env_config = _load_env(tool)
            return {**global_config, **local_config, **env_config, **overrides}
```

