---
content: "None\n\n\n!! class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>Config <em class='small'>class</em></h2>\n\n???+ source \"Config <em class='small'>source</em>\"\n\n```python\n\n
  \       class Config(BaseSettings, JupyterMixin):\n            hooks: list[str]
  = [\"default\"]\n            disabled_hooks: list[str] = []\n            markdown_extensions:
  list[str] = []\n            default_cache_expire: PositiveInt = 3600\n            output_dir:
  pydantic.DirectoryPath = Path(\"markout\")\n            assets_dir: Path = pydantic.Field(\n
  \               Path(\"static\"),\n                description=\"The directory to
  store static assets\",\n            )\n            nav: dict[str, str] = {\"home\":
  \"/\"}\n            site_version: int = 1\n            markdown_backend: str = \"markdown-it-py\"\n
  \           url: Optional[AnyUrl] = None\n            title: Optional[str] = \"Markata
  Site\"\n            description: Optional[str] = None\n            rss_description:
  Optional[str] = None\n            author_name: Optional[str] = None\n            author_email:
  Optional[str] = None\n            lang: str = \"en\"\n            repo_url: Optional[AnyUrl]
  = None\n            repo_branch: str = \"main\"\n            theme_color: Color
  = \"#322D39\"\n            background_color: Color = \"#B73CF6\"\n            start_url:
  str = \"/\"\n            site_name: Optional[str] = None\n            short_name:
  Optional[str] = None\n            display: str = \"minimal-ui\"\n            twitter_card:
  str = \"summary_large_image\"\n            twitter_creator: Optional[str] = None\n
  \           twitter_site: Optional[str] = None\n            path_prefix: Optional[str]
  = \"\"\n            model_config = ConfigDict(env_prefix=\"markata_\", extra=\"allow\")\n
  \           today: datetime.date = pydantic.Field(default_factory=datetime.date.today)\n\n
  \           def __getitem__(self, item):\n                \"for backwards compatability\"\n
  \               return getattr(self, item)\n\n            def __setitem__(self,
  key, item):\n                \"for backwards compatability\"\n                return
  setattr(self, key, item)\n\n            def get(self, item, default):\n                \"for
  backwards compatability\"\n                return getattr(self, item, default)\n\n
  \           def keys(self):\n                \"for backwards compatability\"\n                return
  self.__dict__.keys()\n\n            def toml(self: \"Config\") -> str:\n                import
  tomlkit\n\n                doc = tomlkit.document()\n\n                for key,
  value in self.dict().items():\n                    doc.add(key, value)\n                    doc.add(tomlkit.comment(key))\n
  \                   if value:\n                        doc[key] = value\n                return
  tomlkit.dumps(doc)\n\n            @pydantic.validator(\"output_dir\", pre=True,
  always=True)\n            def validate_output_dir_exists(cls, value: Path) -> Path:\n
  \               if not isinstance(value, Path):\n                    value = Path(value)\n
  \               value.mkdir(parents=True, exist_ok=True)\n                return
  value\n\n            @property\n            def __rich__(self) -> Pretty:\n                return
  lambda: Pretty(self)\n```\n\n\n!! function <h2 id='config_model' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>config_model <em class='small'>function</em></h2>\n\n???+
  source \"config_model <em class='small'>source</em>\"\n\n```python\n\n        def
  config_model(markata: \"Markata\") -> None:\n            markata.config_models.append(Config)\n```\n\n\n!!
  function <h2 id='load_config' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>load_config <em class='small'>function</em></h2>\n\n???+ source \"load_config
  <em class='small'>source</em>\"\n\n```python\n\n        def load_config(markata:
  \"Markata\") -> None:\n            if \"config\" not in markata.__dict__.keys():\n
  \               config = standard_config.load(\"markata\")\n                if config
  == {}:\n                    markata.config = markata.Config()\n                else:\n
  \                   markata.config = markata.Config.parse_obj(config)\n```\n\n\n!!
  class <h2 id='ConfigFactory' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>ConfigFactory <em class='small'>class</em></h2>\n\n???+ source \"ConfigFactory
  <em class='small'>source</em>\"\n\n```python\n\n        class ConfigFactory(ModelFactory):\n
  \           __model__ = Config\n```\n\n\n!! method <h2 id='__getitem__' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>__getitem__ <em class='small'>method</em></h2>\n
  \   for backwards compatability\n???+ source \"__getitem__ <em class='small'>source</em>\"\n\n```python\n\n
  \       def __getitem__(self, item):\n                \"for backwards compatability\"\n
  \               return getattr(self, item)\n```\n\n\n!! method <h2 id='__setitem__'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__setitem__ <em class='small'>method</em></h2>\n
  \   for backwards compatability\n???+ source \"__setitem__ <em class='small'>source</em>\"\n\n```python\n\n
  \       def __setitem__(self, key, item):\n                \"for backwards compatability\"\n
  \               return setattr(self, key, item)\n```\n\n\n!! method <h2 id='get'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get <em class='small'>method</em></h2>\n
  \   for backwards compatability\n???+ source \"get <em class='small'>source</em>\"\n\n```python\n\n
  \       def get(self, item, default):\n                \"for backwards compatability\"\n
  \               return getattr(self, item, default)\n```\n\n\n!! method <h2 id='keys'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>keys <em class='small'>method</em></h2>\n
  \   for backwards compatability\n???+ source \"keys <em class='small'>source</em>\"\n\n```python\n\n
  \       def keys(self):\n                \"for backwards compatability\"\n                return
  self.__dict__.keys()\n```\n\n\n!! method <h2 id='toml' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>toml <em class='small'>method</em></h2>\n\n???+
  source \"toml <em class='small'>source</em>\"\n\n```python\n\n        def toml(self:
  \"Config\") -> str:\n                import tomlkit\n\n                doc = tomlkit.document()\n\n
  \               for key, value in self.dict().items():\n                    doc.add(key,
  value)\n                    doc.add(tomlkit.comment(key))\n                    if
  value:\n                        doc[key] = value\n                return tomlkit.dumps(doc)\n```\n\n\n!!
  method <h2 id='validate_output_dir_exists' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>validate_output_dir_exists <em class='small'>method</em></h2>\n\n???+ source
  \"validate_output_dir_exists <em class='small'>source</em>\"\n\n```python\n\n        def
  validate_output_dir_exists(cls, value: Path) -> Path:\n                if not isinstance(value,
  Path):\n                    value = Path(value)\n                value.mkdir(parents=True,
  exist_ok=True)\n                return value\n```\n\n\n!! method <h2 id='__rich__'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__rich__ <em class='small'>method</em></h2>\n\n???+
  source \"__rich__ <em class='small'>source</em>\"\n\n```python\n\n        def __rich__(self)
  -> Pretty:\n                return lambda: Pretty(self)\n```\n\n"
date: 0001-01-01
description: 'None ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! !
  ! ! ! ???+ source  ! ???+ source  ! ???+ source '
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Config_Model.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None ! ???+ source  ! ???+ source  !
    ???+ source  ! ???+ source  ! ! ! ! ! ???+ source  ! ???+ source  ! ???+ source
    \" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\" />\n\n<link
    rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\" href=\"/app.css\"
    />\n<script src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    <head>\n<title>Config_Model.Py</title>\n<meta charset=\"UTF-8\" />\n<meta
    name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta name=\"description\"
    content=\"None ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! !
    ! ! ! ???+ source  ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Config_Model.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <p>None</p>\n<p>!!
    class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Config
    <em class='small'>class</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">Config <em class='small'>source</em></p>\n</div>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
    <span class=\"nc\">Config</span><span class=\"p\">(</span><span class=\"n\">BaseSettings</span><span
    class=\"p\">,</span> <span class=\"n\">JupyterMixin</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">hooks</span><span class=\"p\">:</span> <span class=\"nb\">list</span><span
    class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">]</span> <span
    class=\"o\">=</span> <span class=\"p\">[</span><span class=\"s2\">&quot;default&quot;</span><span
    class=\"p\">]</span>\n            <span class=\"n\">disabled_hooks</span><span
    class=\"p\">:</span> <span class=\"nb\">list</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"p\">[]</span>\n            <span class=\"n\">markdown_extensions</span><span
    class=\"p\">:</span> <span class=\"nb\">list</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"p\">[]</span>\n            <span class=\"n\">default_cache_expire</span><span
    class=\"p\">:</span> <span class=\"n\">PositiveInt</span> <span class=\"o\">=</span>
    <span class=\"mi\">3600</span>\n            <span class=\"n\">output_dir</span><span
    class=\"p\">:</span> <span class=\"n\">pydantic</span><span class=\"o\">.</span><span
    class=\"n\">DirectoryPath</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;markout&quot;</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">assets_dir</span><span class=\"p\">:</span> <span
    class=\"n\">Path</span> <span class=\"o\">=</span> <span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">Field</span><span class=\"p\">(</span>\n
    \               <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"s2\">&quot;static&quot;</span><span
    class=\"p\">),</span>\n                <span class=\"n\">description</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;The directory to store static assets&quot;</span><span
    class=\"p\">,</span>\n            <span class=\"p\">)</span>\n            <span
    class=\"n\">nav</span><span class=\"p\">:</span> <span class=\"nb\">dict</span><span
    class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"p\">{</span><span class=\"s2\">&quot;home&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;/&quot;</span><span class=\"p\">}</span>\n            <span
    class=\"n\">site_version</span><span class=\"p\">:</span> <span class=\"nb\">int</span>
    <span class=\"o\">=</span> <span class=\"mi\">1</span>\n            <span class=\"n\">markdown_backend</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;markdown-it-py&quot;</span>\n            <span class=\"n\">url</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"n\">AnyUrl</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"kc\">None</span>\n            <span class=\"n\">title</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"s2\">&quot;Markata Site&quot;</span>\n            <span class=\"n\">description</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"kc\">None</span>\n            <span class=\"n\">rss_description</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"kc\">None</span>\n            <span class=\"n\">author_name</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"kc\">None</span>\n            <span class=\"n\">author_email</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"kc\">None</span>\n            <span class=\"n\">lang</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;en&quot;</span>\n
    \           <span class=\"n\">repo_url</span><span class=\"p\">:</span> <span
    class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"n\">AnyUrl</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n
    \           <span class=\"n\">repo_branch</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;main&quot;</span>\n
    \           <span class=\"n\">theme_color</span><span class=\"p\">:</span> <span
    class=\"n\">Color</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;#322D39&quot;</span>\n
    \           <span class=\"n\">background_color</span><span class=\"p\">:</span>
    <span class=\"n\">Color</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;#B73CF6&quot;</span>\n
    \           <span class=\"n\">start_url</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;/&quot;</span>\n
    \           <span class=\"n\">site_name</span><span class=\"p\">:</span> <span
    class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n
    \           <span class=\"n\">short_name</span><span class=\"p\">:</span> <span
    class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n
    \           <span class=\"n\">display</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;minimal-ui&quot;</span>\n
    \           <span class=\"n\">twitter_card</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;summary_large_image&quot;</span>\n
    \           <span class=\"n\">twitter_creator</span><span class=\"p\">:</span>
    <span class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n
    \           <span class=\"n\">twitter_site</span><span class=\"p\">:</span> <span
    class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n
    \           <span class=\"n\">path_prefix</span><span class=\"p\">:</span> <span
    class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;&quot;</span>\n
    \           <span class=\"n\">model_config</span> <span class=\"o\">=</span> <span
    class=\"n\">ConfigDict</span><span class=\"p\">(</span><span class=\"n\">env_prefix</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;markata_&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">extra</span><span class=\"o\">=</span><span class=\"s2\">&quot;allow&quot;</span><span
    class=\"p\">)</span>\n            <span class=\"n\">today</span><span class=\"p\">:</span>
    <span class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">date</span>
    <span class=\"o\">=</span> <span class=\"n\">pydantic</span><span class=\"o\">.</span><span
    class=\"n\">Field</span><span class=\"p\">(</span><span class=\"n\">default_factory</span><span
    class=\"o\">=</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">date</span><span class=\"o\">.</span><span class=\"n\">today</span><span
    class=\"p\">)</span>\n\n            <span class=\"k\">def</span> <span class=\"fm\">__getitem__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"n\">item</span><span class=\"p\">):</span>\n                <span class=\"s2\">&quot;for
    backwards compatability&quot;</span>\n                <span class=\"k\">return</span>
    <span class=\"nb\">getattr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">item</span><span class=\"p\">)</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"fm\">__setitem__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"n\">key</span><span class=\"p\">,</span> <span class=\"n\">item</span><span
    class=\"p\">):</span>\n                <span class=\"s2\">&quot;for backwards
    compatability&quot;</span>\n                <span class=\"k\">return</span> <span
    class=\"nb\">setattr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">key</span><span class=\"p\">,</span> <span
    class=\"n\">item</span><span class=\"p\">)</span>\n\n            <span class=\"k\">def</span>
    <span class=\"nf\">get</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">item</span><span class=\"p\">,</span> <span
    class=\"n\">default</span><span class=\"p\">):</span>\n                <span class=\"s2\">&quot;for
    backwards compatability&quot;</span>\n                <span class=\"k\">return</span>
    <span class=\"nb\">getattr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">item</span><span class=\"p\">,</span> <span
    class=\"n\">default</span><span class=\"p\">)</span>\n\n            <span class=\"k\">def</span>
    <span class=\"nf\">keys</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">):</span>\n                <span class=\"s2\">&quot;for backwards
    compatability&quot;</span>\n                <span class=\"k\">return</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"vm\">__dict__</span><span
    class=\"o\">.</span><span class=\"n\">keys</span><span class=\"p\">()</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"nf\">toml</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Config&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span
    class=\"p\">:</span>\n                <span class=\"kn\">import</span> <span class=\"nn\">tomlkit</span>\n\n
    \               <span class=\"n\">doc</span> <span class=\"o\">=</span> <span
    class=\"n\">tomlkit</span><span class=\"o\">.</span><span class=\"n\">document</span><span
    class=\"p\">()</span>\n\n                <span class=\"k\">for</span> <span class=\"n\">key</span><span
    class=\"p\">,</span> <span class=\"n\">value</span> <span class=\"ow\">in</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">dict</span><span
    class=\"p\">()</span><span class=\"o\">.</span><span class=\"n\">items</span><span
    class=\"p\">():</span>\n                    <span class=\"n\">doc</span><span
    class=\"o\">.</span><span class=\"n\">add</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">,</span> <span class=\"n\">value</span><span
    class=\"p\">)</span>\n                    <span class=\"n\">doc</span><span class=\"o\">.</span><span
    class=\"n\">add</span><span class=\"p\">(</span><span class=\"n\">tomlkit</span><span
    class=\"o\">.</span><span class=\"n\">comment</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">))</span>\n                    <span class=\"k\">if</span>
    <span class=\"n\">value</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">doc</span><span class=\"p\">[</span><span class=\"n\">key</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"n\">value</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">tomlkit</span><span
    class=\"o\">.</span><span class=\"n\">dumps</span><span class=\"p\">(</span><span
    class=\"n\">doc</span><span class=\"p\">)</span>\n\n            <span class=\"nd\">@pydantic</span><span
    class=\"o\">.</span><span class=\"n\">validator</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;output_dir&quot;</span><span class=\"p\">,</span> <span class=\"n\">pre</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
    class=\"n\">always</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n            <span class=\"k\">def</span> <span class=\"nf\">validate_output_dir_exists</span><span
    class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
    class=\"n\">value</span><span class=\"p\">:</span> <span class=\"n\">Path</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Path</span><span
    class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"ow\">not</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">value</span><span
    class=\"p\">,</span> <span class=\"n\">Path</span><span class=\"p\">):</span>\n
    \                   <span class=\"n\">value</span> <span class=\"o\">=</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">value</span><span
    class=\"p\">)</span>\n                <span class=\"n\">value</span><span class=\"o\">.</span><span
    class=\"n\">mkdir</span><span class=\"p\">(</span><span class=\"n\">parents</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
    class=\"n\">exist_ok</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"n\">value</span>\n\n
    \           <span class=\"nd\">@property</span>\n            <span class=\"k\">def</span>
    <span class=\"nf\">__rich__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Pretty</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"k\">lambda</span><span
    class=\"p\">:</span> <span class=\"n\">Pretty</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='config_model' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>config_model <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">config_model
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
    <span class=\"nf\">config_model</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config_models</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">Config</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='load_config' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>load_config <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">load_config
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
    <span class=\"nf\">load_config</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">if</span> <span class=\"s2\">&quot;config&quot;</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"vm\">__dict__</span><span class=\"o\">.</span><span
    class=\"n\">keys</span><span class=\"p\">():</span>\n                <span class=\"n\">config</span>
    <span class=\"o\">=</span> <span class=\"n\">standard_config</span><span class=\"o\">.</span><span
    class=\"n\">load</span><span class=\"p\">(</span><span class=\"s2\">&quot;markata&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"k\">if</span> <span class=\"n\">config</span>
    <span class=\"o\">==</span> <span class=\"p\">{}:</span>\n                    <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">Config</span><span class=\"p\">()</span>\n                <span class=\"k\">else</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">Config</span><span
    class=\"o\">.</span><span class=\"n\">parse_obj</span><span class=\"p\">(</span><span
    class=\"n\">config</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='ConfigFactory' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>ConfigFactory <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">ConfigFactory
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
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
    <span class=\"nc\">ConfigFactory</span><span class=\"p\">(</span><span class=\"n\">ModelFactory</span><span
    class=\"p\">):</span>\n            <span class=\"n\">__model__</span> <span class=\"o\">=</span>
    <span class=\"n\">Config</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='__getitem__'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'><strong>getitem</strong>
    <em class='small'>method</em></h2>\nfor backwards compatability</p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\"><strong>getitem</strong>
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
    <span class=\"fm\">__getitem__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">item</span><span class=\"p\">):</span>\n
    \               <span class=\"s2\">&quot;for backwards compatability&quot;</span>\n
    \               <span class=\"k\">return</span> <span class=\"nb\">getattr</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"n\">item</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='__setitem__' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'><strong>setitem</strong> <em class='small'>method</em></h2>\nfor backwards
    compatability</p>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\"><strong>setitem</strong> <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"fm\">__setitem__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">key</span><span class=\"p\">,</span> <span
    class=\"n\">item</span><span class=\"p\">):</span>\n                <span class=\"s2\">&quot;for
    backwards compatability&quot;</span>\n                <span class=\"k\">return</span>
    <span class=\"nb\">setattr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">key</span><span class=\"p\">,</span> <span
    class=\"n\">item</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='get' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get
    <em class='small'>method</em></h2>\nfor backwards compatability</p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">get <em
    class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    <span class=\"nf\">get</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">item</span><span class=\"p\">,</span> <span
    class=\"n\">default</span><span class=\"p\">):</span>\n                <span class=\"s2\">&quot;for
    backwards compatability&quot;</span>\n                <span class=\"k\">return</span>
    <span class=\"nb\">getattr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">item</span><span class=\"p\">,</span> <span
    class=\"n\">default</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='keys' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>keys
    <em class='small'>method</em></h2>\nfor backwards compatability</p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">keys
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
    <span class=\"nf\">keys</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">):</span>\n                <span class=\"s2\">&quot;for backwards
    compatability&quot;</span>\n                <span class=\"k\">return</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"vm\">__dict__</span><span
    class=\"o\">.</span><span class=\"n\">keys</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='toml' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>toml
    <em class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">toml <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">toml</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Config&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n
    \               <span class=\"kn\">import</span> <span class=\"nn\">tomlkit</span>\n\n
    \               <span class=\"n\">doc</span> <span class=\"o\">=</span> <span
    class=\"n\">tomlkit</span><span class=\"o\">.</span><span class=\"n\">document</span><span
    class=\"p\">()</span>\n\n                <span class=\"k\">for</span> <span class=\"n\">key</span><span
    class=\"p\">,</span> <span class=\"n\">value</span> <span class=\"ow\">in</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">dict</span><span
    class=\"p\">()</span><span class=\"o\">.</span><span class=\"n\">items</span><span
    class=\"p\">():</span>\n                    <span class=\"n\">doc</span><span
    class=\"o\">.</span><span class=\"n\">add</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">,</span> <span class=\"n\">value</span><span
    class=\"p\">)</span>\n                    <span class=\"n\">doc</span><span class=\"o\">.</span><span
    class=\"n\">add</span><span class=\"p\">(</span><span class=\"n\">tomlkit</span><span
    class=\"o\">.</span><span class=\"n\">comment</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">))</span>\n                    <span class=\"k\">if</span>
    <span class=\"n\">value</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">doc</span><span class=\"p\">[</span><span class=\"n\">key</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"n\">value</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">tomlkit</span><span
    class=\"o\">.</span><span class=\"n\">dumps</span><span class=\"p\">(</span><span
    class=\"n\">doc</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='validate_output_dir_exists' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>validate_output_dir_exists <em class='small'>method</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">validate_output_dir_exists
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
    <span class=\"nf\">validate_output_dir_exists</span><span class=\"p\">(</span><span
    class=\"bp\">cls</span><span class=\"p\">,</span> <span class=\"n\">value</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Path</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">if</span> <span class=\"ow\">not</span> <span
    class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">value</span><span
    class=\"p\">,</span> <span class=\"n\">Path</span><span class=\"p\">):</span>\n
    \                   <span class=\"n\">value</span> <span class=\"o\">=</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">value</span><span
    class=\"p\">)</span>\n                <span class=\"n\">value</span><span class=\"o\">.</span><span
    class=\"n\">mkdir</span><span class=\"p\">(</span><span class=\"n\">parents</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
    class=\"n\">exist_ok</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"n\">value</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='__rich__' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'><strong>rich</strong> <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\"><strong>rich</strong>
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
    <span class=\"nf\">__rich__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Pretty</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"k\">lambda</span><span
    class=\"p\">:</span> <span class=\"n\">Pretty</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>    </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Config_Model.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None ! ???+ source  ! ???+ source  !
    ???+ source  ! ???+ source  ! ! ! ! ! ???+ source  ! ???+ source  ! ???+ source
    \" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\" />\n\n<link
    rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\" href=\"/app.css\"
    />\n<script src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    <head>\n<title>Config_Model.Py</title>\n<meta charset=\"UTF-8\" />\n<meta
    name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta name=\"description\"
    content=\"None ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! !
    ! ! ! ???+ source  ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<article style=\"text-align: center;\">\n    <style>\n        section
    {\n            font-size: 200%;\n        }\n\n\n        .edit {\n            display:
    none;\n        }\n    </style>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Config_Model.Py\n    </h1>\n</section></article>\n     </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Config_Model.Py\n    </h1>\n</section>    <section class=\"body\">\n        <p>None</p>\n<p>!!
    class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Config
    <em class='small'>class</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">Config <em class='small'>source</em></p>\n</div>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
    <span class=\"nc\">Config</span><span class=\"p\">(</span><span class=\"n\">BaseSettings</span><span
    class=\"p\">,</span> <span class=\"n\">JupyterMixin</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">hooks</span><span class=\"p\">:</span> <span class=\"nb\">list</span><span
    class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">]</span> <span
    class=\"o\">=</span> <span class=\"p\">[</span><span class=\"s2\">&quot;default&quot;</span><span
    class=\"p\">]</span>\n            <span class=\"n\">disabled_hooks</span><span
    class=\"p\">:</span> <span class=\"nb\">list</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"p\">[]</span>\n            <span class=\"n\">markdown_extensions</span><span
    class=\"p\">:</span> <span class=\"nb\">list</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"p\">[]</span>\n            <span class=\"n\">default_cache_expire</span><span
    class=\"p\">:</span> <span class=\"n\">PositiveInt</span> <span class=\"o\">=</span>
    <span class=\"mi\">3600</span>\n            <span class=\"n\">output_dir</span><span
    class=\"p\">:</span> <span class=\"n\">pydantic</span><span class=\"o\">.</span><span
    class=\"n\">DirectoryPath</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;markout&quot;</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">assets_dir</span><span class=\"p\">:</span> <span
    class=\"n\">Path</span> <span class=\"o\">=</span> <span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">Field</span><span class=\"p\">(</span>\n
    \               <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"s2\">&quot;static&quot;</span><span
    class=\"p\">),</span>\n                <span class=\"n\">description</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;The directory to store static assets&quot;</span><span
    class=\"p\">,</span>\n            <span class=\"p\">)</span>\n            <span
    class=\"n\">nav</span><span class=\"p\">:</span> <span class=\"nb\">dict</span><span
    class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"p\">{</span><span class=\"s2\">&quot;home&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;/&quot;</span><span class=\"p\">}</span>\n            <span
    class=\"n\">site_version</span><span class=\"p\">:</span> <span class=\"nb\">int</span>
    <span class=\"o\">=</span> <span class=\"mi\">1</span>\n            <span class=\"n\">markdown_backend</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;markdown-it-py&quot;</span>\n            <span class=\"n\">url</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"n\">AnyUrl</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"kc\">None</span>\n            <span class=\"n\">title</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"s2\">&quot;Markata Site&quot;</span>\n            <span class=\"n\">description</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"kc\">None</span>\n            <span class=\"n\">rss_description</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"kc\">None</span>\n            <span class=\"n\">author_name</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"kc\">None</span>\n            <span class=\"n\">author_email</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"kc\">None</span>\n            <span class=\"n\">lang</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;en&quot;</span>\n
    \           <span class=\"n\">repo_url</span><span class=\"p\">:</span> <span
    class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"n\">AnyUrl</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n
    \           <span class=\"n\">repo_branch</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;main&quot;</span>\n
    \           <span class=\"n\">theme_color</span><span class=\"p\">:</span> <span
    class=\"n\">Color</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;#322D39&quot;</span>\n
    \           <span class=\"n\">background_color</span><span class=\"p\">:</span>
    <span class=\"n\">Color</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;#B73CF6&quot;</span>\n
    \           <span class=\"n\">start_url</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;/&quot;</span>\n
    \           <span class=\"n\">site_name</span><span class=\"p\">:</span> <span
    class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n
    \           <span class=\"n\">short_name</span><span class=\"p\">:</span> <span
    class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n
    \           <span class=\"n\">display</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;minimal-ui&quot;</span>\n
    \           <span class=\"n\">twitter_card</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;summary_large_image&quot;</span>\n
    \           <span class=\"n\">twitter_creator</span><span class=\"p\">:</span>
    <span class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n
    \           <span class=\"n\">twitter_site</span><span class=\"p\">:</span> <span
    class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n
    \           <span class=\"n\">path_prefix</span><span class=\"p\">:</span> <span
    class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;&quot;</span>\n
    \           <span class=\"n\">model_config</span> <span class=\"o\">=</span> <span
    class=\"n\">ConfigDict</span><span class=\"p\">(</span><span class=\"n\">env_prefix</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;markata_&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">extra</span><span class=\"o\">=</span><span class=\"s2\">&quot;allow&quot;</span><span
    class=\"p\">)</span>\n            <span class=\"n\">today</span><span class=\"p\">:</span>
    <span class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">date</span>
    <span class=\"o\">=</span> <span class=\"n\">pydantic</span><span class=\"o\">.</span><span
    class=\"n\">Field</span><span class=\"p\">(</span><span class=\"n\">default_factory</span><span
    class=\"o\">=</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">date</span><span class=\"o\">.</span><span class=\"n\">today</span><span
    class=\"p\">)</span>\n\n            <span class=\"k\">def</span> <span class=\"fm\">__getitem__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"n\">item</span><span class=\"p\">):</span>\n                <span class=\"s2\">&quot;for
    backwards compatability&quot;</span>\n                <span class=\"k\">return</span>
    <span class=\"nb\">getattr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">item</span><span class=\"p\">)</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"fm\">__setitem__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"n\">key</span><span class=\"p\">,</span> <span class=\"n\">item</span><span
    class=\"p\">):</span>\n                <span class=\"s2\">&quot;for backwards
    compatability&quot;</span>\n                <span class=\"k\">return</span> <span
    class=\"nb\">setattr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">key</span><span class=\"p\">,</span> <span
    class=\"n\">item</span><span class=\"p\">)</span>\n\n            <span class=\"k\">def</span>
    <span class=\"nf\">get</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">item</span><span class=\"p\">,</span> <span
    class=\"n\">default</span><span class=\"p\">):</span>\n                <span class=\"s2\">&quot;for
    backwards compatability&quot;</span>\n                <span class=\"k\">return</span>
    <span class=\"nb\">getattr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">item</span><span class=\"p\">,</span> <span
    class=\"n\">default</span><span class=\"p\">)</span>\n\n            <span class=\"k\">def</span>
    <span class=\"nf\">keys</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">):</span>\n                <span class=\"s2\">&quot;for backwards
    compatability&quot;</span>\n                <span class=\"k\">return</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"vm\">__dict__</span><span
    class=\"o\">.</span><span class=\"n\">keys</span><span class=\"p\">()</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"nf\">toml</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Config&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span
    class=\"p\">:</span>\n                <span class=\"kn\">import</span> <span class=\"nn\">tomlkit</span>\n\n
    \               <span class=\"n\">doc</span> <span class=\"o\">=</span> <span
    class=\"n\">tomlkit</span><span class=\"o\">.</span><span class=\"n\">document</span><span
    class=\"p\">()</span>\n\n                <span class=\"k\">for</span> <span class=\"n\">key</span><span
    class=\"p\">,</span> <span class=\"n\">value</span> <span class=\"ow\">in</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">dict</span><span
    class=\"p\">()</span><span class=\"o\">.</span><span class=\"n\">items</span><span
    class=\"p\">():</span>\n                    <span class=\"n\">doc</span><span
    class=\"o\">.</span><span class=\"n\">add</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">,</span> <span class=\"n\">value</span><span
    class=\"p\">)</span>\n                    <span class=\"n\">doc</span><span class=\"o\">.</span><span
    class=\"n\">add</span><span class=\"p\">(</span><span class=\"n\">tomlkit</span><span
    class=\"o\">.</span><span class=\"n\">comment</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">))</span>\n                    <span class=\"k\">if</span>
    <span class=\"n\">value</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">doc</span><span class=\"p\">[</span><span class=\"n\">key</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"n\">value</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">tomlkit</span><span
    class=\"o\">.</span><span class=\"n\">dumps</span><span class=\"p\">(</span><span
    class=\"n\">doc</span><span class=\"p\">)</span>\n\n            <span class=\"nd\">@pydantic</span><span
    class=\"o\">.</span><span class=\"n\">validator</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;output_dir&quot;</span><span class=\"p\">,</span> <span class=\"n\">pre</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
    class=\"n\">always</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n            <span class=\"k\">def</span> <span class=\"nf\">validate_output_dir_exists</span><span
    class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
    class=\"n\">value</span><span class=\"p\">:</span> <span class=\"n\">Path</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Path</span><span
    class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"ow\">not</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">value</span><span
    class=\"p\">,</span> <span class=\"n\">Path</span><span class=\"p\">):</span>\n
    \                   <span class=\"n\">value</span> <span class=\"o\">=</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">value</span><span
    class=\"p\">)</span>\n                <span class=\"n\">value</span><span class=\"o\">.</span><span
    class=\"n\">mkdir</span><span class=\"p\">(</span><span class=\"n\">parents</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
    class=\"n\">exist_ok</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"n\">value</span>\n\n
    \           <span class=\"nd\">@property</span>\n            <span class=\"k\">def</span>
    <span class=\"nf\">__rich__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Pretty</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"k\">lambda</span><span
    class=\"p\">:</span> <span class=\"n\">Pretty</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='config_model' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>config_model <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">config_model
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
    <span class=\"nf\">config_model</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config_models</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">Config</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='load_config' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>load_config <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">load_config
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
    <span class=\"nf\">load_config</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">if</span> <span class=\"s2\">&quot;config&quot;</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"vm\">__dict__</span><span class=\"o\">.</span><span
    class=\"n\">keys</span><span class=\"p\">():</span>\n                <span class=\"n\">config</span>
    <span class=\"o\">=</span> <span class=\"n\">standard_config</span><span class=\"o\">.</span><span
    class=\"n\">load</span><span class=\"p\">(</span><span class=\"s2\">&quot;markata&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"k\">if</span> <span class=\"n\">config</span>
    <span class=\"o\">==</span> <span class=\"p\">{}:</span>\n                    <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">Config</span><span class=\"p\">()</span>\n                <span class=\"k\">else</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">Config</span><span
    class=\"o\">.</span><span class=\"n\">parse_obj</span><span class=\"p\">(</span><span
    class=\"n\">config</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='ConfigFactory' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>ConfigFactory <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">ConfigFactory
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
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
    <span class=\"nc\">ConfigFactory</span><span class=\"p\">(</span><span class=\"n\">ModelFactory</span><span
    class=\"p\">):</span>\n            <span class=\"n\">__model__</span> <span class=\"o\">=</span>
    <span class=\"n\">Config</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='__getitem__'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'><strong>getitem</strong>
    <em class='small'>method</em></h2>\nfor backwards compatability</p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\"><strong>getitem</strong>
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
    <span class=\"fm\">__getitem__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">item</span><span class=\"p\">):</span>\n
    \               <span class=\"s2\">&quot;for backwards compatability&quot;</span>\n
    \               <span class=\"k\">return</span> <span class=\"nb\">getattr</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"n\">item</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='__setitem__' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'><strong>setitem</strong> <em class='small'>method</em></h2>\nfor backwards
    compatability</p>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\"><strong>setitem</strong> <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"fm\">__setitem__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">key</span><span class=\"p\">,</span> <span
    class=\"n\">item</span><span class=\"p\">):</span>\n                <span class=\"s2\">&quot;for
    backwards compatability&quot;</span>\n                <span class=\"k\">return</span>
    <span class=\"nb\">setattr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">key</span><span class=\"p\">,</span> <span
    class=\"n\">item</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='get' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get
    <em class='small'>method</em></h2>\nfor backwards compatability</p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">get <em
    class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    <span class=\"nf\">get</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">item</span><span class=\"p\">,</span> <span
    class=\"n\">default</span><span class=\"p\">):</span>\n                <span class=\"s2\">&quot;for
    backwards compatability&quot;</span>\n                <span class=\"k\">return</span>
    <span class=\"nb\">getattr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">item</span><span class=\"p\">,</span> <span
    class=\"n\">default</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='keys' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>keys
    <em class='small'>method</em></h2>\nfor backwards compatability</p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">keys
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
    <span class=\"nf\">keys</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">):</span>\n                <span class=\"s2\">&quot;for backwards
    compatability&quot;</span>\n                <span class=\"k\">return</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"vm\">__dict__</span><span
    class=\"o\">.</span><span class=\"n\">keys</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='toml' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>toml
    <em class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">toml <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">toml</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Config&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n
    \               <span class=\"kn\">import</span> <span class=\"nn\">tomlkit</span>\n\n
    \               <span class=\"n\">doc</span> <span class=\"o\">=</span> <span
    class=\"n\">tomlkit</span><span class=\"o\">.</span><span class=\"n\">document</span><span
    class=\"p\">()</span>\n\n                <span class=\"k\">for</span> <span class=\"n\">key</span><span
    class=\"p\">,</span> <span class=\"n\">value</span> <span class=\"ow\">in</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">dict</span><span
    class=\"p\">()</span><span class=\"o\">.</span><span class=\"n\">items</span><span
    class=\"p\">():</span>\n                    <span class=\"n\">doc</span><span
    class=\"o\">.</span><span class=\"n\">add</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">,</span> <span class=\"n\">value</span><span
    class=\"p\">)</span>\n                    <span class=\"n\">doc</span><span class=\"o\">.</span><span
    class=\"n\">add</span><span class=\"p\">(</span><span class=\"n\">tomlkit</span><span
    class=\"o\">.</span><span class=\"n\">comment</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">))</span>\n                    <span class=\"k\">if</span>
    <span class=\"n\">value</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">doc</span><span class=\"p\">[</span><span class=\"n\">key</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"n\">value</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">tomlkit</span><span
    class=\"o\">.</span><span class=\"n\">dumps</span><span class=\"p\">(</span><span
    class=\"n\">doc</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='validate_output_dir_exists' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>validate_output_dir_exists <em class='small'>method</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">validate_output_dir_exists
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
    <span class=\"nf\">validate_output_dir_exists</span><span class=\"p\">(</span><span
    class=\"bp\">cls</span><span class=\"p\">,</span> <span class=\"n\">value</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Path</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">if</span> <span class=\"ow\">not</span> <span
    class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">value</span><span
    class=\"p\">,</span> <span class=\"n\">Path</span><span class=\"p\">):</span>\n
    \                   <span class=\"n\">value</span> <span class=\"o\">=</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">value</span><span
    class=\"p\">)</span>\n                <span class=\"n\">value</span><span class=\"o\">.</span><span
    class=\"n\">mkdir</span><span class=\"p\">(</span><span class=\"n\">parents</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
    class=\"n\">exist_ok</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"n\">value</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='__rich__' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'><strong>rich</strong> <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\"><strong>rich</strong>
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
    <span class=\"nf\">__rich__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Pretty</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"k\">lambda</span><span
    class=\"p\">:</span> <span class=\"n\">Pretty</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>"
  raw.md: ''
published: true
slug: markata/plugins/config-model
title: Config_Model.Py


---

None


!! class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Config <em class='small'>class</em></h2>

???+ source "Config <em class='small'>source</em>"

```python

        class Config(BaseSettings, JupyterMixin):
            hooks: list[str] = ["default"]
            disabled_hooks: list[str] = []
            markdown_extensions: list[str] = []
            default_cache_expire: PositiveInt = 3600
            output_dir: pydantic.DirectoryPath = Path("markout")
            assets_dir: Path = pydantic.Field(
                Path("static"),
                description="The directory to store static assets",
            )
            nav: dict[str, str] = {"home": "/"}
            site_version: int = 1
            markdown_backend: str = "markdown-it-py"
            url: Optional[AnyUrl] = None
            title: Optional[str] = "Markata Site"
            description: Optional[str] = None
            rss_description: Optional[str] = None
            author_name: Optional[str] = None
            author_email: Optional[str] = None
            lang: str = "en"
            repo_url: Optional[AnyUrl] = None
            repo_branch: str = "main"
            theme_color: Color = "#322D39"
            background_color: Color = "#B73CF6"
            start_url: str = "/"
            site_name: Optional[str] = None
            short_name: Optional[str] = None
            display: str = "minimal-ui"
            twitter_card: str = "summary_large_image"
            twitter_creator: Optional[str] = None
            twitter_site: Optional[str] = None
            path_prefix: Optional[str] = ""
            model_config = ConfigDict(env_prefix="markata_", extra="allow")
            today: datetime.date = pydantic.Field(default_factory=datetime.date.today)

            def __getitem__(self, item):
                "for backwards compatability"
                return getattr(self, item)

            def __setitem__(self, key, item):
                "for backwards compatability"
                return setattr(self, key, item)

            def get(self, item, default):
                "for backwards compatability"
                return getattr(self, item, default)

            def keys(self):
                "for backwards compatability"
                return self.__dict__.keys()

            def toml(self: "Config") -> str:
                import tomlkit

                doc = tomlkit.document()

                for key, value in self.dict().items():
                    doc.add(key, value)
                    doc.add(tomlkit.comment(key))
                    if value:
                        doc[key] = value
                return tomlkit.dumps(doc)

            @pydantic.validator("output_dir", pre=True, always=True)
            def validate_output_dir_exists(cls, value: Path) -> Path:
                if not isinstance(value, Path):
                    value = Path(value)
                value.mkdir(parents=True, exist_ok=True)
                return value

            @property
            def __rich__(self) -> Pretty:
                return lambda: Pretty(self)
```


!! function <h2 id='config_model' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>config_model <em class='small'>function</em></h2>

???+ source "config_model <em class='small'>source</em>"

```python

        def config_model(markata: "Markata") -> None:
            markata.config_models.append(Config)
```


!! function <h2 id='load_config' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>load_config <em class='small'>function</em></h2>

???+ source "load_config <em class='small'>source</em>"

```python

        def load_config(markata: "Markata") -> None:
            if "config" not in markata.__dict__.keys():
                config = standard_config.load("markata")
                if config == {}:
                    markata.config = markata.Config()
                else:
                    markata.config = markata.Config.parse_obj(config)
```


!! class <h2 id='ConfigFactory' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>ConfigFactory <em class='small'>class</em></h2>

???+ source "ConfigFactory <em class='small'>source</em>"

```python

        class ConfigFactory(ModelFactory):
            __model__ = Config
```


!! method <h2 id='__getitem__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__getitem__ <em class='small'>method</em></h2>
    for backwards compatability
???+ source "__getitem__ <em class='small'>source</em>"

```python

        def __getitem__(self, item):
                "for backwards compatability"
                return getattr(self, item)
```


!! method <h2 id='__setitem__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__setitem__ <em class='small'>method</em></h2>
    for backwards compatability
???+ source "__setitem__ <em class='small'>source</em>"

```python

        def __setitem__(self, key, item):
                "for backwards compatability"
                return setattr(self, key, item)
```


!! method <h2 id='get' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get <em class='small'>method</em></h2>
    for backwards compatability
???+ source "get <em class='small'>source</em>"

```python

        def get(self, item, default):
                "for backwards compatability"
                return getattr(self, item, default)
```


!! method <h2 id='keys' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>keys <em class='small'>method</em></h2>
    for backwards compatability
???+ source "keys <em class='small'>source</em>"

```python

        def keys(self):
                "for backwards compatability"
                return self.__dict__.keys()
```


!! method <h2 id='toml' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>toml <em class='small'>method</em></h2>

???+ source "toml <em class='small'>source</em>"

```python

        def toml(self: "Config") -> str:
                import tomlkit

                doc = tomlkit.document()

                for key, value in self.dict().items():
                    doc.add(key, value)
                    doc.add(tomlkit.comment(key))
                    if value:
                        doc[key] = value
                return tomlkit.dumps(doc)
```


!! method <h2 id='validate_output_dir_exists' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>validate_output_dir_exists <em class='small'>method</em></h2>

???+ source "validate_output_dir_exists <em class='small'>source</em>"

```python

        def validate_output_dir_exists(cls, value: Path) -> Path:
                if not isinstance(value, Path):
                    value = Path(value)
                value.mkdir(parents=True, exist_ok=True)
                return value
```


!! method <h2 id='__rich__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__rich__ <em class='small'>method</em></h2>

???+ source "__rich__ <em class='small'>source</em>"

```python

        def __rich__(self) -> Pretty:
                return lambda: Pretty(self)
```

