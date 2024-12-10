---
content: "Markata is a tool for handling directories of markdown.\n\n\n!! class <h2
  id='HooksConfig' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>HooksConfig
  <em class='small'>class</em></h2>\n\n???+ source \"HooksConfig <em class='small'>source</em>\"\n\n```python\n\n
  \       class HooksConfig(pydantic.BaseModel):\n            hooks: list = [\"default\"]\n
  \           disabled_hooks: list = []\n```\n\n\n!! class <h2 id='Markata' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>Markata <em class='small'>class</em></h2>\n\n???+
  source \"Markata <em class='small'>source</em>\"\n\n```python\n\n        class Markata:\n
  \           def __init__(self: \"Markata\", console: Console = None, config=None)
  -> None:\n                self.__version__ = __version__\n                self.stages_ran
  = set()\n                self.threded = False\n                self._cache = None\n
  \               self._precache = None\n                self.MARKATA_CACHE_DIR =
  Path(\".\") / \".markata.cache\"\n                self.MARKATA_CACHE_DIR.mkdir(exist_ok=True)\n
  \               self._pm = pluggy.PluginManager(\"markata\")\n                self._pm.add_hookspecs(hookspec.MarkataSpecs)\n
  \               if config is not None:\n                    self.config = config\n
  \               with self.cache as cache:\n                    self.init_cache_stats
  = cache.stats()\n                self.registered_attrs = hookspec.registered_attrs\n
  \               self.post_models = []\n                self.config_models = []\n
  \               if config is not None:\n                    raw_hooks = config\n
  \               else:\n                    raw_hooks = standard_config.load(\"markata\")\n
  \               self.hooks_conf = HooksConfig.parse_obj(raw_hooks)\n                try:\n
  \                   default_index = self.hooks_conf.hooks.index(\"default\")\n                    hooks
  = [\n                        *self.hooks_conf.hooks[:default_index],\n                        *DEFAULT_HOOKS,\n
  \                       *self.hooks_conf.hooks[default_index + 1 :],\n                    ]\n
  \                   self.hooks_conf.hooks = [\n                        hook for
  hook in hooks if hook not in self.hooks_conf.disabled_hooks\n                    ]\n
  \               except ValueError:\n                    # 'default' is not in hooks
  , do not replace with default_hooks\n                    pass\n\n                self._register_hooks()\n
  \               if console is not None:\n                    self._console = console\n
  \               atexit.register(self.teardown)\n                self.precache\n\n
  \           @property\n            def cache(self: \"Markata\") -> Cache:\n                #
  if self.threded:\n                #     FanoutCache(self.MARKATA_CACHE_DIR, statistics=True)\n
  \               if self._cache is not None:\n                    return self._cache\n
  \               self._cache = Cache(self.MARKATA_CACHE_DIR, statistics=True)\n\n
  \               return self._cache\n\n            @property\n            def precache(self:
  \"Markata\") -> None:\n                if self._precache is None:\n                    self.cache.expire()\n
  \                   self._precache = {k: self.cache.get(k) for k in self.cache.iterkeys()}\n
  \               return self._precache\n\n            def __getattr__(self: \"Markata\",
  item: str) -> Any:\n                if item in self._pm.hook.__dict__:\n                    #
  item is a hook, return a callable function\n                    return lambda: self.run(item)\n\n
  \               if item in self.__dict__:\n                    # item is an attribute,
  return it\n                    return self.__getitem__(item)\n\n                elif
  item in self.registered_attrs:\n                    # item is created by a plugin,
  run it\n                    stage_to_run_to = max(\n                        [attr[\"lifecycle\"]
  for attr in self.registered_attrs[item]],\n                    ).name\n                    self.console.log(\n
  \                       f\"Running to [purple]{stage_to_run_to}[/] to retrieve [purple]{item}[/]\"\n
  \                   )\n                    self.run(stage_to_run_to)\n                    return
  getattr(self, item)\n                elif item == \"precache\":\n                    return
  self._precache or {}\n                else:\n                    # Markata does
  not know what this is, raise\n                    raise AttributeError(f\"'Markata'
  object has no attribute '{item}'\")\n\n            def __rich__(self: \"Markata\")
  -> Table:\n                grid = Table.grid()\n                grid.add_column(\"label\")\n
  \               grid.add_column(\"value\")\n\n                for label, value in
  self.describe().items():\n                    grid.add_row(label, value)\n\n                return
  grid\n\n            def bust_cache(self: \"Markata\") -> Markata:\n                with
  self.cache as cache:\n                    cache.clear()\n                return
  self\n\n            def configure(self) -> Markata:\n                sys.path.append(os.getcwd())\n
  \               # self.config = {**DEFUALT_CONFIG, **standard_config.load(\"markata\")}\n
  \               # if isinstance(self.config[\"glob_patterns\"], str):\n                #
  \    self.config[\"glob_patterns\"] = self.config[\"glob_patterns\"].split(\",\")\n
  \               # elif isinstance(self.config[\"glob_patterns\"], list):\n                #
  \    self.config[\"glob_patterns\"] = list(self.config[\"glob_patterns\"])\n                #
  else:\n                #     raise TypeError(\"glob_patterns must be list or str\")\n
  \               # self.glob_patterns = self.config[\"glob_patterns\"]\n\n                #
  self.hooks = self.config[\"hooks\"]\n\n                # if \"disabled_hooks\" not
  in self.config:\n                #     self.disabled_hooks = [\"\"]\n                #
  if isinstance(self.config[\"disabled_hooks\"], str):\n                #     self.disabled_hooks
  = self.config[\"disabled_hooks\"].split(\",\")\n                # if isinstance(self.config[\"disabled_hooks\"],
  list):\n                #     self.disabled_hooks = self.config[\"disabled_hooks\"]\n\n
  \               # if not self.config.get(\"output_dir\", \"markout\").endswith(\n
  \               #     self.config.get(\"path_prefix\", \"\")\n                #
  ):\n                #     self.config[\"output_dir\"] = (\n                #         self.config.get(\"output_dir\",
  \"markout\") +\n                #         \"/\" +\n                #         self.config.get(\"path_prefix\",
  \"\").rstrip(\"/\")\n                #     )\n                # if (\n                #
  \    len((output_split := self.config.get(\"output_dir\", \"markout\").split(\"/\")))
  >\n                #     1\n                # ):\n                #     if \"path_prefix\"
  not in self.config.keys():\n                #         self.config[\"path_prefix\"]
  = \"/\".join(output_split[1:]) + \"/\"\n                # if not self.config.get(\"path_prefix\",
  \"\").endswith(\"/\"):\n                #     self.config[\"path_prefix\"] = self.config.get(\"path_prefix\",
  \"\") + \"/\"\n\n                # self.config[\"output_dir\"] = self.config[\"output_dir\"].lstrip(\"/\")\n
  \               # self.config[\"path_prefix\"] = self.config[\"path_prefix\"].lstrip(\"/\")\n\n
  \               try:\n                    default_index = self.hooks_conf.hooks.index(\"default\")\n
  \                   hooks = [\n                        *self.hooks_conf.hooks[:default_index],\n
  \                       *DEFAULT_HOOKS,\n                        *self.hooks_conf.hooks[default_index
  + 1 :],\n                    ]\n                    self.config.hooks = [\n                        hook
  for hook in hooks if hook not in self.config.disabled_hooks\n                    ]\n
  \               except ValueError:\n                    # 'default' is not in hooks
  , do not replace with default_hooks\n                    pass\n\n                self._pm
  = pluggy.PluginManager(\"markata\")\n                self._pm.add_hookspecs(hookspec.MarkataSpecs)\n
  \               self._register_hooks()\n\n                self._pm.hook.configure(markata=self)\n
  \               return self\n\n            def get_plugin_config(self, path_or_name:
  str) -> Dict:\n                key = Path(path_or_name).stem\n\n                config
  = self.config.get(key, {})\n\n                if not isinstance(config, dict):\n
  \                   raise TypeError(\"must use dict\")\n                if \"cache_expire\"
  not in config.keys():\n                    config[\"cache_expire\"] = self.config[\"default_cache_expire\"]\n
  \               if \"config_key\" not in config.keys():\n                    config[\"config_key\"]
  = key\n                return config\n\n            def get_config(\n                self,\n
  \               key: str,\n                default: str = \"\",\n                warn:
  bool = True,\n                suggested: Optional[str] = None,\n            ) ->
  Any:\n                if key in self.config.keys():\n                    return
  self.config[key]\n                else:\n                    if suggested is None:\n
  \                       suggested = textwrap.dedent(\n                            f\"\"\"\n
  \                           [markata]\n                            {key} = '{default}'\n
  \                           \"\"\"\n                        )\n                    if
  warn:\n                        logger.warning(\n                            textwrap.dedent(\n
  \                               f\"\"\"\n                                Warning
  {key} is not set in markata config, sitemap will\n                                be
  missing root site_name\n                                to resolve this open your
  markata.toml and add\n\n                                {suggested}\n\n                                \"\"\"\n
  \                           ),\n                        )\n                return
  default\n\n            def make_hash(self, *keys: str) -> str:\n                import
  xxhash\n\n                str_keys = [str(key) for key in keys]\n                hash
  = xxhash.xxh64(\"\".join(str_keys).encode(\"utf-8\")).hexdigest()\n                return
  hash\n\n            @property\n            def content_dir_hash(self: \"Markata\")
  -> str:\n                hashes = [\n                    dirhash(dir)\n                    for
  dir in self.content_directories\n                    if dir.absolute() != Path(\".\").absolute()\n
  \               ]\n                return self.make_hash(*hashes)\n\n            @property\n
  \           def console(self: \"Markata\") -> Console:\n                try:\n                    return
  self._console\n                except AttributeError:\n                    self._console
  = Console()\n                    return self._console\n\n            def describe(self:
  \"Markata\") -> dict[str, str]:\n                return {\"version\": __version__}\n\n
  \           def _to_dict(self: \"Markata\") -> dict[str, Iterable]:\n                return
  {\"config\": self.config, \"articles\": [a.to_dict() for a in self.articles]}\n\n
  \           def to_dict(self: \"Markata\") -> dict:\n                return self._to_dict()\n\n
  \           def to_json(self: \"Markata\") -> str:\n                import json\n\n
  \               return json.dumps(self.to_dict(), indent=4, sort_keys=True, default=str)\n\n
  \           def _register_hooks(self: \"Markata\") -> None:\n                sys.path.append(os.getcwd())\n
  \               for hook in self.hooks_conf.hooks:\n                    try:\n                        #
  module style plugins\n                        plugin = importlib.import_module(hook)\n
  \                   except ModuleNotFoundError as e:\n                        #
  class style plugins\n                        if \".\" in hook:\n                            try:\n
  \                               mod = importlib.import_module(\".\".join(hook.split(\".\")[:-1]))\n
  \                               plugin = getattr(mod, hook.split(\".\")[-1])\n                            except
  ModuleNotFoundError as e:\n                                raise ModuleNotFoundError(\n
  \                                   f\"module {hook} not found\\n{sys.path}\"\n
  \                               ) from e\n                        else:\n                            raise
  e\n\n                    self._pm.register(plugin)\n\n            def __iter__(\n
  \               self: \"Markata\", description: str = \"working...\"\n            )
  -> Iterable[\"Markata.Post\"]:\n                articles: Iterable[Markata.Post]
  = track(\n                    self.articles,\n                    description=description,\n
  \                   transient=False,\n                    console=self.console,\n
  \               )\n                return articles\n\n            def iter_articles(self:
  \"Markata\", description: str) -> Iterable[Markata.Post]:\n                articles:
  Iterable[Markata.Post] = track(\n                    self.articles,\n                    description=description,\n
  \                   transient=True,\n                    console=self.console,\n
  \               )\n                return articles\n\n            def teardown(self:
  \"Markata\") -> Markata:\n                \"\"\"give special access to the teardown
  lifecycle method\"\"\"\n                self._pm.hook.teardown(markata=self)\n                return
  self\n\n            def run(self: \"Markata\", lifecycle: LifeCycle = None) -> Markata:\n
  \               if lifecycle is None:\n                    lifecycle = max(LifeCycle._member_map_.values())\n\n
  \               if isinstance(lifecycle, str):\n                    lifecycle =
  LifeCycle[lifecycle]\n\n                stages_to_run = [\n                    m\n
  \                   for m in LifeCycle._member_map_\n                    if (LifeCycle[m]
  <= lifecycle) and (m not in self.stages_ran)\n                ]\n\n                if
  not stages_to_run:\n                    self.console.log(f\"{lifecycle.name} already
  ran\")\n                    return self\n\n                self.console.log(f\"running
  {stages_to_run}\")\n                for stage in stages_to_run:\n                    self.console.log(f\"{stage}
  running\")\n                    getattr(self._pm.hook, stage)(markata=self)\n                    self.stages_ran.add(stage)\n
  \                   self.console.log(f\"{stage} complete\")\n\n                with
  self.cache as cache:\n                    hits, misses = cache.stats()\n\n                if
  hits + misses > 0:\n                    self.console.log(\n                        f\"lifetime
  cache hit rate {round(hits/ (hits + misses)*100, 2)}%\",\n                    )\n\n
  \               if misses > 0:\n                    self.console.log(f\"lifetime
  cache hits/misses {hits}/{misses}\")\n\n                hits -= self.init_cache_stats[0]\n
  \               misses -= self.init_cache_stats[1]\n\n                if hits +
  misses > 0:\n                    self.console.log(\n                        f\"run
  cache hit rate {round(hits/ (hits + misses)*100, 2)}%\",\n                    )\n\n
  \               if misses > 0:\n                    self.console.log(f\"run cache
  hits/misses {hits}/{misses}\")\n\n                return self\n\n            def
  filter(self: \"Markata\", filter: str) -> list:\n                def evalr(a: Markata.Post)
  -> Any:\n                    try:\n                        return eval(\n                            filter,\n
  \                           {**a.to_dict(), \"timedelta\": timedelta, \"post\":
  a, \"m\": self},\n                            {},\n                        )\n                    except
  AttributeError:\n                        return eval(\n                            filter,\n
  \                           {**a.to_dict(), \"timedelta\": timedelta, \"post\":
  a, \"m\": self},\n                            {},\n                        )\n\n
  \               return [a for a in self.articles if evalr(a)]\n\n            def
  map(\n                self: \"Markata\",\n                func: str = \"title\",\n
  \               filter: str = \"True\",\n                sort: str = \"True\",\n
  \               reverse: bool = True,\n                *args: tuple,\n                **kwargs:
  dict,\n            ) -> list:\n                import copy\n\n                def
  try_sort(a: Any) -> int:\n                    if \"datetime\" in sort.lower():\n
  \                       return a.get(sort, datetime.datetime(1970, 1, 1))\n\n                    if
  \"date\" in sort.lower():\n                        return a.get(sort, datetime.date(1970,
  1, 1))\n\n                    try:\n                        value = eval(sort, a.to_dict(),
  {})\n                    except NameError:\n                        return -1\n
  \                   return value\n                    try:\n                        return
  int(value)\n                    except TypeError:\n                        try:\n
  \                           return int(value.timestamp())\n                        except
  Exception:\n                            try:\n                                return
  int(\n                                    datetime.datetime.combine(\n                                        value,\n
  \                                       datetime.datetime.min.time(),\n                                    ).timestamp(),\n
  \                               )\n                            except Exception:\n
  \                               try:\n                                    return
  sum([ord(c) for c in str(value)])\n                                except Exception:\n
  \                                   return -1\n\n                articles = copy.copy(self.articles)\n
  \               articles.sort(key=try_sort)\n                if reverse:\n                    articles.reverse()\n\n
  \               try:\n                    posts = [\n                        eval(\n
  \                           func,\n                            {**a.to_dict(), \"timedelta\":
  timedelta, \"post\": a, \"m\": self},\n                            {},\n                        )\n
  \                       for a in articles\n                        if eval(\n                            filter,\n
  \                           {**a.to_dict(), \"timedelta\": timedelta, \"post\":
  a, \"m\": self},\n                            {},\n                        )\n                    ]\n\n
  \               except NameError as e:\n                    variable = str(e).split(\"'\")[1]\n\n
  \                   missing_in_posts = self.map(\n                        \"path\",\n
  \                       filter=f'\"{variable}\" not in post.keys()',\n                    )\n
  \                   message = (\n                        f\"variable: '{variable}'
  is missing in {len(missing_in_posts)} posts\"\n                    )\n                    if
  len(missing_in_posts) > 10:\n                        message += (\n                            f\"\\nfirst
  10 paths to posts missing {variable}\"\n                            f\"[{','.join([str(p)
  for p in missing_in_posts[:10]])}...\"\n                        )\n                    else:\n
  \                       message += f\"\\npaths to posts missing {variable} {missing_in_posts}\"\n\n
  \                   raise MissingFrontMatter(message)\n\n                return
  posts\n\n            def first(\n                self: \"Markata\",\n                filter:
  str = \"True\",\n                sort: str = \"True\",\n                reverse:
  bool = True,\n                *args: tuple,\n                **kwargs: dict,\n            )
  -> list:\n                return self.map(\"post\", filter, sort, reverse, *args,
  **kwargs)[0]\n\n            def last(\n                self: \"Markata\",\n                filter:
  str = \"True\",\n                sort: str = \"True\",\n                reverse:
  bool = True,\n                *args: tuple,\n                **kwargs: dict,\n            )
  -> list:\n                return self.map(\"post\", filter, sort, reverse, *args,
  **kwargs)[-1]\n\n            def one(\n                self: \"Markata\",\n                filter:
  str = \"True\",\n                *args: tuple,\n                **kwargs: dict,\n
  \           ) -> list:\n                posts = self.map(\"post\", filter, *args,
  **kwargs)\n                if len(posts) > 1:\n                    raise TooManyPosts(f\"found
  {len(posts)} posts, expected 1. {posts}\")\n                if len(posts) == 0:\n
  \                   raise NoPosts\n                return posts[0]\n```\n\n\n!!
  function <h2 id='load_ipython_extension' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>load_ipython_extension <em class='small'>function</em></h2>\n\n???+ source
  \"load_ipython_extension <em class='small'>source</em>\"\n\n```python\n\n        def
  load_ipython_extension(ipython):\n            ipython.user_ns[\"m\"] = Markata()\n
  \           ipython.user_ns[\"markata\"] = ipython.user_ns[\"m\"]\n            ipython.user_ns[\"markata\"]
  = ipython.user_ns[\"m\"]\n```\n\n\n!! method <h2 id='__init__' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>__init__ <em class='small'>method</em></h2>\n\n???+
  source \"__init__ <em class='small'>source</em>\"\n\n```python\n\n        def __init__(self:
  \"Markata\", console: Console = None, config=None) -> None:\n                self.__version__
  = __version__\n                self.stages_ran = set()\n                self.threded
  = False\n                self._cache = None\n                self._precache = None\n
  \               self.MARKATA_CACHE_DIR = Path(\".\") / \".markata.cache\"\n                self.MARKATA_CACHE_DIR.mkdir(exist_ok=True)\n
  \               self._pm = pluggy.PluginManager(\"markata\")\n                self._pm.add_hookspecs(hookspec.MarkataSpecs)\n
  \               if config is not None:\n                    self.config = config\n
  \               with self.cache as cache:\n                    self.init_cache_stats
  = cache.stats()\n                self.registered_attrs = hookspec.registered_attrs\n
  \               self.post_models = []\n                self.config_models = []\n
  \               if config is not None:\n                    raw_hooks = config\n
  \               else:\n                    raw_hooks = standard_config.load(\"markata\")\n
  \               self.hooks_conf = HooksConfig.parse_obj(raw_hooks)\n                try:\n
  \                   default_index = self.hooks_conf.hooks.index(\"default\")\n                    hooks
  = [\n                        *self.hooks_conf.hooks[:default_index],\n                        *DEFAULT_HOOKS,\n
  \                       *self.hooks_conf.hooks[default_index + 1 :],\n                    ]\n
  \                   self.hooks_conf.hooks = [\n                        hook for
  hook in hooks if hook not in self.hooks_conf.disabled_hooks\n                    ]\n
  \               except ValueError:\n                    # 'default' is not in hooks
  , do not replace with default_hooks\n                    pass\n\n                self._register_hooks()\n
  \               if console is not None:\n                    self._console = console\n
  \               atexit.register(self.teardown)\n                self.precache\n```\n\n\n!!
  method <h2 id='cache' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>cache
  <em class='small'>method</em></h2>\n\n???+ source \"cache <em class='small'>source</em>\"\n\n```python\n\n
  \       def cache(self: \"Markata\") -> Cache:\n                # if self.threded:\n
  \               #     FanoutCache(self.MARKATA_CACHE_DIR, statistics=True)\n                if
  self._cache is not None:\n                    return self._cache\n                self._cache
  = Cache(self.MARKATA_CACHE_DIR, statistics=True)\n\n                return self._cache\n```\n\n\n!!
  method <h2 id='precache' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>precache <em class='small'>method</em></h2>\n\n???+ source \"precache <em
  class='small'>source</em>\"\n\n```python\n\n        def precache(self: \"Markata\")
  -> None:\n                if self._precache is None:\n                    self.cache.expire()\n
  \                   self._precache = {k: self.cache.get(k) for k in self.cache.iterkeys()}\n
  \               return self._precache\n```\n\n\n!! method <h2 id='__getattr__' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>__getattr__ <em class='small'>method</em></h2>\n\n???+
  source \"__getattr__ <em class='small'>source</em>\"\n\n```python\n\n        def
  __getattr__(self: \"Markata\", item: str) -> Any:\n                if item in self._pm.hook.__dict__:\n
  \                   # item is a hook, return a callable function\n                    return
  lambda: self.run(item)\n\n                if item in self.__dict__:\n                    #
  item is an attribute, return it\n                    return self.__getitem__(item)\n\n
  \               elif item in self.registered_attrs:\n                    # item
  is created by a plugin, run it\n                    stage_to_run_to = max(\n                        [attr[\"lifecycle\"]
  for attr in self.registered_attrs[item]],\n                    ).name\n                    self.console.log(\n
  \                       f\"Running to [purple]{stage_to_run_to}[/] to retrieve [purple]{item}[/]\"\n
  \                   )\n                    self.run(stage_to_run_to)\n                    return
  getattr(self, item)\n                elif item == \"precache\":\n                    return
  self._precache or {}\n                else:\n                    # Markata does
  not know what this is, raise\n                    raise AttributeError(f\"'Markata'
  object has no attribute '{item}'\")\n```\n\n\n!! method <h2 id='__rich__' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>__rich__ <em class='small'>method</em></h2>\n\n???+
  source \"__rich__ <em class='small'>source</em>\"\n\n```python\n\n        def __rich__(self:
  \"Markata\") -> Table:\n                grid = Table.grid()\n                grid.add_column(\"label\")\n
  \               grid.add_column(\"value\")\n\n                for label, value in
  self.describe().items():\n                    grid.add_row(label, value)\n\n                return
  grid\n```\n\n\n!! method <h2 id='bust_cache' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>bust_cache <em class='small'>method</em></h2>\n\n???+ source \"bust_cache
  <em class='small'>source</em>\"\n\n```python\n\n        def bust_cache(self: \"Markata\")
  -> Markata:\n                with self.cache as cache:\n                    cache.clear()\n
  \               return self\n```\n\n\n!! method <h2 id='configure' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>configure <em class='small'>method</em></h2>\n\n???+
  source \"configure <em class='small'>source</em>\"\n\n```python\n\n        def configure(self)
  -> Markata:\n                sys.path.append(os.getcwd())\n                # self.config
  = {**DEFUALT_CONFIG, **standard_config.load(\"markata\")}\n                # if
  isinstance(self.config[\"glob_patterns\"], str):\n                #     self.config[\"glob_patterns\"]
  = self.config[\"glob_patterns\"].split(\",\")\n                # elif isinstance(self.config[\"glob_patterns\"],
  list):\n                #     self.config[\"glob_patterns\"] = list(self.config[\"glob_patterns\"])\n
  \               # else:\n                #     raise TypeError(\"glob_patterns must
  be list or str\")\n                # self.glob_patterns = self.config[\"glob_patterns\"]\n\n
  \               # self.hooks = self.config[\"hooks\"]\n\n                # if \"disabled_hooks\"
  not in self.config:\n                #     self.disabled_hooks = [\"\"]\n                #
  if isinstance(self.config[\"disabled_hooks\"], str):\n                #     self.disabled_hooks
  = self.config[\"disabled_hooks\"].split(\",\")\n                # if isinstance(self.config[\"disabled_hooks\"],
  list):\n                #     self.disabled_hooks = self.config[\"disabled_hooks\"]\n\n
  \               # if not self.config.get(\"output_dir\", \"markout\").endswith(\n
  \               #     self.config.get(\"path_prefix\", \"\")\n                #
  ):\n                #     self.config[\"output_dir\"] = (\n                #         self.config.get(\"output_dir\",
  \"markout\") +\n                #         \"/\" +\n                #         self.config.get(\"path_prefix\",
  \"\").rstrip(\"/\")\n                #     )\n                # if (\n                #
  \    len((output_split := self.config.get(\"output_dir\", \"markout\").split(\"/\")))
  >\n                #     1\n                # ):\n                #     if \"path_prefix\"
  not in self.config.keys():\n                #         self.config[\"path_prefix\"]
  = \"/\".join(output_split[1:]) + \"/\"\n                # if not self.config.get(\"path_prefix\",
  \"\").endswith(\"/\"):\n                #     self.config[\"path_prefix\"] = self.config.get(\"path_prefix\",
  \"\") + \"/\"\n\n                # self.config[\"output_dir\"] = self.config[\"output_dir\"].lstrip(\"/\")\n
  \               # self.config[\"path_prefix\"] = self.config[\"path_prefix\"].lstrip(\"/\")\n\n
  \               try:\n                    default_index = self.hooks_conf.hooks.index(\"default\")\n
  \                   hooks = [\n                        *self.hooks_conf.hooks[:default_index],\n
  \                       *DEFAULT_HOOKS,\n                        *self.hooks_conf.hooks[default_index
  + 1 :],\n                    ]\n                    self.config.hooks = [\n                        hook
  for hook in hooks if hook not in self.config.disabled_hooks\n                    ]\n
  \               except ValueError:\n                    # 'default' is not in hooks
  , do not replace with default_hooks\n                    pass\n\n                self._pm
  = pluggy.PluginManager(\"markata\")\n                self._pm.add_hookspecs(hookspec.MarkataSpecs)\n
  \               self._register_hooks()\n\n                self._pm.hook.configure(markata=self)\n
  \               return self\n```\n\n\n!! method <h2 id='get_plugin_config' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>get_plugin_config <em class='small'>method</em></h2>\n\n???+
  source \"get_plugin_config <em class='small'>source</em>\"\n\n```python\n\n        def
  get_plugin_config(self, path_or_name: str) -> Dict:\n                key = Path(path_or_name).stem\n\n
  \               config = self.config.get(key, {})\n\n                if not isinstance(config,
  dict):\n                    raise TypeError(\"must use dict\")\n                if
  \"cache_expire\" not in config.keys():\n                    config[\"cache_expire\"]
  = self.config[\"default_cache_expire\"]\n                if \"config_key\" not in
  config.keys():\n                    config[\"config_key\"] = key\n                return
  config\n```\n\n\n!! method <h2 id='get_config' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>get_config <em class='small'>method</em></h2>\n\n???+ source \"get_config
  <em class='small'>source</em>\"\n\n```python\n\n        def get_config(\n                self,\n
  \               key: str,\n                default: str = \"\",\n                warn:
  bool = True,\n                suggested: Optional[str] = None,\n            ) ->
  Any:\n                if key in self.config.keys():\n                    return
  self.config[key]\n                else:\n                    if suggested is None:\n
  \                       suggested = textwrap.dedent(\n                            f\"\"\"\n
  \                           [markata]\n                            {key} = '{default}'\n
  \                           \"\"\"\n                        )\n                    if
  warn:\n                        logger.warning(\n                            textwrap.dedent(\n
  \                               f\"\"\"\n                                Warning
  {key} is not set in markata config, sitemap will\n                                be
  missing root site_name\n                                to resolve this open your
  markata.toml and add\n\n                                {suggested}\n\n                                \"\"\"\n
  \                           ),\n                        )\n                return
  default\n```\n\n\n!! method <h2 id='make_hash' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>make_hash <em class='small'>method</em></h2>\n\n???+ source \"make_hash <em
  class='small'>source</em>\"\n\n```python\n\n        def make_hash(self, *keys: str)
  -> str:\n                import xxhash\n\n                str_keys = [str(key) for
  key in keys]\n                hash = xxhash.xxh64(\"\".join(str_keys).encode(\"utf-8\")).hexdigest()\n
  \               return hash\n```\n\n\n!! method <h2 id='content_dir_hash' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>content_dir_hash <em class='small'>method</em></h2>\n\n???+
  source \"content_dir_hash <em class='small'>source</em>\"\n\n```python\n\n        def
  content_dir_hash(self: \"Markata\") -> str:\n                hashes = [\n                    dirhash(dir)\n
  \                   for dir in self.content_directories\n                    if
  dir.absolute() != Path(\".\").absolute()\n                ]\n                return
  self.make_hash(*hashes)\n```\n\n\n!! method <h2 id='console' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>console <em class='small'>method</em></h2>\n\n???+
  source \"console <em class='small'>source</em>\"\n\n```python\n\n        def console(self:
  \"Markata\") -> Console:\n                try:\n                    return self._console\n
  \               except AttributeError:\n                    self._console = Console()\n
  \                   return self._console\n```\n\n\n!! method <h2 id='describe' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>describe <em class='small'>method</em></h2>\n\n???+
  source \"describe <em class='small'>source</em>\"\n\n```python\n\n        def describe(self:
  \"Markata\") -> dict[str, str]:\n                return {\"version\": __version__}\n```\n\n\n!!
  method <h2 id='_to_dict' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>_to_dict <em class='small'>method</em></h2>\n\n???+ source \"_to_dict <em
  class='small'>source</em>\"\n\n```python\n\n        def _to_dict(self: \"Markata\")
  -> dict[str, Iterable]:\n                return {\"config\": self.config, \"articles\":
  [a.to_dict() for a in self.articles]}\n```\n\n\n!! method <h2 id='to_dict' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>to_dict <em class='small'>method</em></h2>\n\n???+
  source \"to_dict <em class='small'>source</em>\"\n\n```python\n\n        def to_dict(self:
  \"Markata\") -> dict:\n                return self._to_dict()\n```\n\n\n!! method
  <h2 id='to_json' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>to_json
  <em class='small'>method</em></h2>\n\n???+ source \"to_json <em class='small'>source</em>\"\n\n```python\n\n
  \       def to_json(self: \"Markata\") -> str:\n                import json\n\n
  \               return json.dumps(self.to_dict(), indent=4, sort_keys=True, default=str)\n```\n\n\n!!
  method <h2 id='_register_hooks' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>_register_hooks <em class='small'>method</em></h2>\n\n???+ source \"_register_hooks
  <em class='small'>source</em>\"\n\n```python\n\n        def _register_hooks(self:
  \"Markata\") -> None:\n                sys.path.append(os.getcwd())\n                for
  hook in self.hooks_conf.hooks:\n                    try:\n                        #
  module style plugins\n                        plugin = importlib.import_module(hook)\n
  \                   except ModuleNotFoundError as e:\n                        #
  class style plugins\n                        if \".\" in hook:\n                            try:\n
  \                               mod = importlib.import_module(\".\".join(hook.split(\".\")[:-1]))\n
  \                               plugin = getattr(mod, hook.split(\".\")[-1])\n                            except
  ModuleNotFoundError as e:\n                                raise ModuleNotFoundError(\n
  \                                   f\"module {hook} not found\\n{sys.path}\"\n
  \                               ) from e\n                        else:\n                            raise
  e\n\n                    self._pm.register(plugin)\n```\n\n\n!! method <h2 id='__iter__'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__iter__ <em class='small'>method</em></h2>\n\n???+
  source \"__iter__ <em class='small'>source</em>\"\n\n```python\n\n        def __iter__(\n
  \               self: \"Markata\", description: str = \"working...\"\n            )
  -> Iterable[\"Markata.Post\"]:\n                articles: Iterable[Markata.Post]
  = track(\n                    self.articles,\n                    description=description,\n
  \                   transient=False,\n                    console=self.console,\n
  \               )\n                return articles\n```\n\n\n!! method <h2 id='iter_articles'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>iter_articles <em
  class='small'>method</em></h2>\n\n???+ source \"iter_articles <em class='small'>source</em>\"\n\n```python\n\n
  \       def iter_articles(self: \"Markata\", description: str) -> Iterable[Markata.Post]:\n
  \               articles: Iterable[Markata.Post] = track(\n                    self.articles,\n
  \                   description=description,\n                    transient=True,\n
  \                   console=self.console,\n                )\n                return
  articles\n```\n\n\n!! method <h2 id='teardown' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>teardown <em class='small'>method</em></h2>\n    give special access to the
  teardown lifecycle method\n???+ source \"teardown <em class='small'>source</em>\"\n\n```python\n\n
  \       def teardown(self: \"Markata\") -> Markata:\n                \"\"\"give
  special access to the teardown lifecycle method\"\"\"\n                self._pm.hook.teardown(markata=self)\n
  \               return self\n```\n\n\n!! method <h2 id='run' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>run <em class='small'>method</em></h2>\n\n???+
  source \"run <em class='small'>source</em>\"\n\n```python\n\n        def run(self:
  \"Markata\", lifecycle: LifeCycle = None) -> Markata:\n                if lifecycle
  is None:\n                    lifecycle = max(LifeCycle._member_map_.values())\n\n
  \               if isinstance(lifecycle, str):\n                    lifecycle =
  LifeCycle[lifecycle]\n\n                stages_to_run = [\n                    m\n
  \                   for m in LifeCycle._member_map_\n                    if (LifeCycle[m]
  <= lifecycle) and (m not in self.stages_ran)\n                ]\n\n                if
  not stages_to_run:\n                    self.console.log(f\"{lifecycle.name} already
  ran\")\n                    return self\n\n                self.console.log(f\"running
  {stages_to_run}\")\n                for stage in stages_to_run:\n                    self.console.log(f\"{stage}
  running\")\n                    getattr(self._pm.hook, stage)(markata=self)\n                    self.stages_ran.add(stage)\n
  \                   self.console.log(f\"{stage} complete\")\n\n                with
  self.cache as cache:\n                    hits, misses = cache.stats()\n\n                if
  hits + misses > 0:\n                    self.console.log(\n                        f\"lifetime
  cache hit rate {round(hits/ (hits + misses)*100, 2)}%\",\n                    )\n\n
  \               if misses > 0:\n                    self.console.log(f\"lifetime
  cache hits/misses {hits}/{misses}\")\n\n                hits -= self.init_cache_stats[0]\n
  \               misses -= self.init_cache_stats[1]\n\n                if hits +
  misses > 0:\n                    self.console.log(\n                        f\"run
  cache hit rate {round(hits/ (hits + misses)*100, 2)}%\",\n                    )\n\n
  \               if misses > 0:\n                    self.console.log(f\"run cache
  hits/misses {hits}/{misses}\")\n\n                return self\n```\n\n\n!! method
  <h2 id='filter' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>filter
  <em class='small'>method</em></h2>\n\n???+ source \"filter <em class='small'>source</em>\"\n\n```python\n\n
  \       def filter(self: \"Markata\", filter: str) -> list:\n                def
  evalr(a: Markata.Post) -> Any:\n                    try:\n                        return
  eval(\n                            filter,\n                            {**a.to_dict(),
  \"timedelta\": timedelta, \"post\": a, \"m\": self},\n                            {},\n
  \                       )\n                    except AttributeError:\n                        return
  eval(\n                            filter,\n                            {**a.to_dict(),
  \"timedelta\": timedelta, \"post\": a, \"m\": self},\n                            {},\n
  \                       )\n\n                return [a for a in self.articles if
  evalr(a)]\n```\n\n\n!! method <h2 id='map' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>map <em class='small'>method</em></h2>\n\n???+ source \"map <em class='small'>source</em>\"\n\n```python\n\n
  \       def map(\n                self: \"Markata\",\n                func: str
  = \"title\",\n                filter: str = \"True\",\n                sort: str
  = \"True\",\n                reverse: bool = True,\n                *args: tuple,\n
  \               **kwargs: dict,\n            ) -> list:\n                import
  copy\n\n                def try_sort(a: Any) -> int:\n                    if \"datetime\"
  in sort.lower():\n                        return a.get(sort, datetime.datetime(1970,
  1, 1))\n\n                    if \"date\" in sort.lower():\n                        return
  a.get(sort, datetime.date(1970, 1, 1))\n\n                    try:\n                        value
  = eval(sort, a.to_dict(), {})\n                    except NameError:\n                        return
  -1\n                    return value\n                    try:\n                        return
  int(value)\n                    except TypeError:\n                        try:\n
  \                           return int(value.timestamp())\n                        except
  Exception:\n                            try:\n                                return
  int(\n                                    datetime.datetime.combine(\n                                        value,\n
  \                                       datetime.datetime.min.time(),\n                                    ).timestamp(),\n
  \                               )\n                            except Exception:\n
  \                               try:\n                                    return
  sum([ord(c) for c in str(value)])\n                                except Exception:\n
  \                                   return -1\n\n                articles = copy.copy(self.articles)\n
  \               articles.sort(key=try_sort)\n                if reverse:\n                    articles.reverse()\n\n
  \               try:\n                    posts = [\n                        eval(\n
  \                           func,\n                            {**a.to_dict(), \"timedelta\":
  timedelta, \"post\": a, \"m\": self},\n                            {},\n                        )\n
  \                       for a in articles\n                        if eval(\n                            filter,\n
  \                           {**a.to_dict(), \"timedelta\": timedelta, \"post\":
  a, \"m\": self},\n                            {},\n                        )\n                    ]\n\n
  \               except NameError as e:\n                    variable = str(e).split(\"'\")[1]\n\n
  \                   missing_in_posts = self.map(\n                        \"path\",\n
  \                       filter=f'\"{variable}\" not in post.keys()',\n                    )\n
  \                   message = (\n                        f\"variable: '{variable}'
  is missing in {len(missing_in_posts)} posts\"\n                    )\n                    if
  len(missing_in_posts) > 10:\n                        message += (\n                            f\"\\nfirst
  10 paths to posts missing {variable}\"\n                            f\"[{','.join([str(p)
  for p in missing_in_posts[:10]])}...\"\n                        )\n                    else:\n
  \                       message += f\"\\npaths to posts missing {variable} {missing_in_posts}\"\n\n
  \                   raise MissingFrontMatter(message)\n\n                return
  posts\n```\n\n\n!! method <h2 id='first' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>first <em class='small'>method</em></h2>\n\n???+ source \"first <em class='small'>source</em>\"\n\n```python\n\n
  \       def first(\n                self: \"Markata\",\n                filter:
  str = \"True\",\n                sort: str = \"True\",\n                reverse:
  bool = True,\n                *args: tuple,\n                **kwargs: dict,\n            )
  -> list:\n                return self.map(\"post\", filter, sort, reverse, *args,
  **kwargs)[0]\n```\n\n\n!! method <h2 id='last' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>last <em class='small'>method</em></h2>\n\n???+ source \"last <em class='small'>source</em>\"\n\n```python\n\n
  \       def last(\n                self: \"Markata\",\n                filter: str
  = \"True\",\n                sort: str = \"True\",\n                reverse: bool
  = True,\n                *args: tuple,\n                **kwargs: dict,\n            )
  -> list:\n                return self.map(\"post\", filter, sort, reverse, *args,
  **kwargs)[-1]\n```\n\n\n!! method <h2 id='one' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>one <em class='small'>method</em></h2>\n\n???+ source \"one <em class='small'>source</em>\"\n\n```python\n\n
  \       def one(\n                self: \"Markata\",\n                filter: str
  = \"True\",\n                *args: tuple,\n                **kwargs: dict,\n            )
  -> list:\n                posts = self.map(\"post\", filter, *args, **kwargs)\n
  \               if len(posts) > 1:\n                    raise TooManyPosts(f\"found
  {len(posts)} posts, expected 1. {posts}\")\n                if len(posts) == 0:\n
  \                   raise NoPosts\n                return posts[0]\n```\n\n\n!!
  function <h2 id='evalr' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>evalr
  <em class='small'>function</em></h2>\n\n???+ source \"evalr <em class='small'>source</em>\"\n\n```python\n\n
  \       def evalr(a: Markata.Post) -> Any:\n                    try:\n                        return
  eval(\n                            filter,\n                            {**a.to_dict(),
  \"timedelta\": timedelta, \"post\": a, \"m\": self},\n                            {},\n
  \                       )\n                    except AttributeError:\n                        return
  eval(\n                            filter,\n                            {**a.to_dict(),
  \"timedelta\": timedelta, \"post\": a, \"m\": self},\n                            {},\n
  \                       )\n```\n\n\n!! function <h2 id='try_sort' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>try_sort <em class='small'>function</em></h2>\n\n???+
  source \"try_sort <em class='small'>source</em>\"\n\n```python\n\n        def try_sort(a:
  Any) -> int:\n                    if \"datetime\" in sort.lower():\n                        return
  a.get(sort, datetime.datetime(1970, 1, 1))\n\n                    if \"date\" in
  sort.lower():\n                        return a.get(sort, datetime.date(1970, 1,
  1))\n\n                    try:\n                        value = eval(sort, a.to_dict(),
  {})\n                    except NameError:\n                        return -1\n
  \                   return value\n                    try:\n                        return
  int(value)\n                    except TypeError:\n                        try:\n
  \                           return int(value.timestamp())\n                        except
  Exception:\n                            try:\n                                return
  int(\n                                    datetime.datetime.combine(\n                                        value,\n
  \                                       datetime.datetime.min.time(),\n                                    ).timestamp(),\n
  \                               )\n                            except Exception:\n
  \                               try:\n                                    return
  sum([ord(c) for c in str(value)])\n                                except Exception:\n
  \                                   return -1\n```\n\n"
date: 0001-01-01
description: 'Markata is a tool for handling directories of markdown. ! ???+ source  !
  ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
  source '
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>__Init__.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"Markata is a tool for handling directories of markdown.
    ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
    source  ! ???+ source \" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>__Init__.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Markata is a tool for handling directories
    of markdown. ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
    source  ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\" rel=\"icon\"
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        __Init__.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <p>Markata is a tool
    for handling directories of markdown.</p>\n<p>!! class <h2 id='HooksConfig' class='admonition-title'
    style='margin:0;padding:.5rem 1rem;'>HooksConfig <em class='small'>class</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">HooksConfig
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
    <span class=\"nc\">HooksConfig</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">hooks</span><span class=\"p\">:</span> <span class=\"nb\">list</span>
    <span class=\"o\">=</span> <span class=\"p\">[</span><span class=\"s2\">&quot;default&quot;</span><span
    class=\"p\">]</span>\n            <span class=\"n\">disabled_hooks</span><span
    class=\"p\">:</span> <span class=\"nb\">list</span> <span class=\"o\">=</span>
    <span class=\"p\">[]</span>\n</pre></div>\n\n</pre>\n\n<p>!! class <h2 id='Markata'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Markata <em class='small'>class</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">Markata
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
    <span class=\"nc\">Markata</span><span class=\"p\">:</span>\n            <span
    class=\"k\">def</span> <span class=\"fm\">__init__</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">console</span><span class=\"p\">:</span>
    <span class=\"n\">Console</span> <span class=\"o\">=</span> <span class=\"kc\">None</span><span
    class=\"p\">,</span> <span class=\"n\">config</span><span class=\"o\">=</span><span
    class=\"kc\">None</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">__version__</span>
    <span class=\"o\">=</span> <span class=\"n\">__version__</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">stages_ran</span>
    <span class=\"o\">=</span> <span class=\"nb\">set</span><span class=\"p\">()</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">threded</span> <span class=\"o\">=</span> <span class=\"kc\">False</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_cache</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_precache</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">MARKATA_CACHE_DIR</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;.&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">/</span> <span class=\"s2\">&quot;.markata.cache&quot;</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">MARKATA_CACHE_DIR</span><span class=\"o\">.</span><span class=\"n\">mkdir</span><span
    class=\"p\">(</span><span class=\"n\">exist_ok</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_pm</span> <span class=\"o\">=</span> <span
    class=\"n\">pluggy</span><span class=\"o\">.</span><span class=\"n\">PluginManager</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;markata&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_pm</span><span class=\"o\">.</span><span class=\"n\">add_hookspecs</span><span
    class=\"p\">(</span><span class=\"n\">hookspec</span><span class=\"o\">.</span><span
    class=\"n\">MarkataSpecs</span><span class=\"p\">)</span>\n                <span
    class=\"k\">if</span> <span class=\"n\">config</span> <span class=\"ow\">is</span>
    <span class=\"ow\">not</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">config</span> <span class=\"o\">=</span> <span class=\"n\">config</span>\n
    \               <span class=\"k\">with</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">cache</span> <span class=\"k\">as</span>
    <span class=\"n\">cache</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">init_cache_stats</span>
    <span class=\"o\">=</span> <span class=\"n\">cache</span><span class=\"o\">.</span><span
    class=\"n\">stats</span><span class=\"p\">()</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">registered_attrs</span> <span class=\"o\">=</span>
    <span class=\"n\">hookspec</span><span class=\"o\">.</span><span class=\"n\">registered_attrs</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">post_models</span> <span class=\"o\">=</span> <span class=\"p\">[]</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">config_models</span> <span class=\"o\">=</span> <span class=\"p\">[]</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">config</span> <span
    class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">raw_hooks</span> <span
    class=\"o\">=</span> <span class=\"n\">config</span>\n                <span class=\"k\">else</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">raw_hooks</span> <span
    class=\"o\">=</span> <span class=\"n\">standard_config</span><span class=\"o\">.</span><span
    class=\"n\">load</span><span class=\"p\">(</span><span class=\"s2\">&quot;markata&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">hooks_conf</span> <span class=\"o\">=</span> <span class=\"n\">HooksConfig</span><span
    class=\"o\">.</span><span class=\"n\">parse_obj</span><span class=\"p\">(</span><span
    class=\"n\">raw_hooks</span><span class=\"p\">)</span>\n                <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                    <span class=\"n\">default_index</span>
    <span class=\"o\">=</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">hooks_conf</span><span class=\"o\">.</span><span class=\"n\">hooks</span><span
    class=\"o\">.</span><span class=\"n\">index</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;default&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"n\">hooks</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
    \                       <span class=\"o\">*</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">hooks_conf</span><span class=\"o\">.</span><span
    class=\"n\">hooks</span><span class=\"p\">[:</span><span class=\"n\">default_index</span><span
    class=\"p\">],</span>\n                        <span class=\"o\">*</span><span
    class=\"n\">DEFAULT_HOOKS</span><span class=\"p\">,</span>\n                        <span
    class=\"o\">*</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">hooks_conf</span><span class=\"o\">.</span><span class=\"n\">hooks</span><span
    class=\"p\">[</span><span class=\"n\">default_index</span> <span class=\"o\">+</span>
    <span class=\"mi\">1</span> <span class=\"p\">:],</span>\n                    <span
    class=\"p\">]</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">hooks_conf</span><span class=\"o\">.</span><span
    class=\"n\">hooks</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
    \                       <span class=\"n\">hook</span> <span class=\"k\">for</span>
    <span class=\"n\">hook</span> <span class=\"ow\">in</span> <span class=\"n\">hooks</span>
    <span class=\"k\">if</span> <span class=\"n\">hook</span> <span class=\"ow\">not</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">hooks_conf</span><span class=\"o\">.</span><span class=\"n\">disabled_hooks</span>\n
    \                   <span class=\"p\">]</span>\n                <span class=\"k\">except</span>
    <span class=\"ne\">ValueError</span><span class=\"p\">:</span>\n                    <span
    class=\"c1\"># &#39;default&#39; is not in hooks , do not replace with default_hooks</span>\n
    \                   <span class=\"k\">pass</span>\n\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_register_hooks</span><span class=\"p\">()</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">console</span> <span
    class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_console</span> <span class=\"o\">=</span>
    <span class=\"n\">console</span>\n                <span class=\"n\">atexit</span><span
    class=\"o\">.</span><span class=\"n\">register</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">teardown</span><span
    class=\"p\">)</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">precache</span>\n\n            <span class=\"nd\">@property</span>\n
    \           <span class=\"k\">def</span> <span class=\"nf\">cache</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">:</span> <span
    class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Cache</span><span class=\"p\">:</span>\n                <span
    class=\"c1\"># if self.threded:</span>\n                <span class=\"c1\">#     FanoutCache(self.MARKATA_CACHE_DIR,
    statistics=True)</span>\n                <span class=\"k\">if</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_cache</span> <span class=\"ow\">is</span>
    <span class=\"ow\">not</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_cache</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_cache</span> <span class=\"o\">=</span>
    <span class=\"n\">Cache</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">MARKATA_CACHE_DIR</span><span class=\"p\">,</span>
    <span class=\"n\">statistics</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">return</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_cache</span>\n\n
    \           <span class=\"nd\">@property</span>\n            <span class=\"k\">def</span>
    <span class=\"nf\">precache</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">if</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_precache</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">cache</span><span
    class=\"o\">.</span><span class=\"n\">expire</span><span class=\"p\">()</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_precache</span> <span class=\"o\">=</span> <span class=\"p\">{</span><span
    class=\"n\">k</span><span class=\"p\">:</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">cache</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"n\">k</span><span
    class=\"p\">)</span> <span class=\"k\">for</span> <span class=\"n\">k</span> <span
    class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">cache</span><span class=\"o\">.</span><span class=\"n\">iterkeys</span><span
    class=\"p\">()}</span>\n                <span class=\"k\">return</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_precache</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"fm\">__getattr__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">:</span> <span
    class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span> <span class=\"n\">item</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">item</span> <span
    class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_pm</span><span class=\"o\">.</span><span class=\"n\">hook</span><span
    class=\"o\">.</span><span class=\"vm\">__dict__</span><span class=\"p\">:</span>\n
    \                   <span class=\"c1\"># item is a hook, return a callable function</span>\n
    \                   <span class=\"k\">return</span> <span class=\"k\">lambda</span><span
    class=\"p\">:</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">run</span><span class=\"p\">(</span><span class=\"n\">item</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">if</span> <span class=\"n\">item</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"vm\">__dict__</span><span class=\"p\">:</span>\n                    <span
    class=\"c1\"># item is an attribute, return it</span>\n                    <span
    class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"fm\">__getitem__</span><span class=\"p\">(</span><span class=\"n\">item</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">elif</span> <span class=\"n\">item</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">registered_attrs</span><span class=\"p\">:</span>\n                    <span
    class=\"c1\"># item is created by a plugin, run it</span>\n                    <span
    class=\"n\">stage_to_run_to</span> <span class=\"o\">=</span> <span class=\"nb\">max</span><span
    class=\"p\">(</span>\n                        <span class=\"p\">[</span><span
    class=\"n\">attr</span><span class=\"p\">[</span><span class=\"s2\">&quot;lifecycle&quot;</span><span
    class=\"p\">]</span> <span class=\"k\">for</span> <span class=\"n\">attr</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">registered_attrs</span><span class=\"p\">[</span><span class=\"n\">item</span><span
    class=\"p\">]],</span>\n                    <span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">name</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">log</span><span class=\"p\">(</span>\n                        <span
    class=\"sa\">f</span><span class=\"s2\">&quot;Running to [purple]</span><span
    class=\"si\">{</span><span class=\"n\">stage_to_run_to</span><span class=\"si\">}</span><span
    class=\"s2\">[/] to retrieve [purple]</span><span class=\"si\">{</span><span class=\"n\">item</span><span
    class=\"si\">}</span><span class=\"s2\">[/]&quot;</span>\n                    <span
    class=\"p\">)</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">run</span><span class=\"p\">(</span><span
    class=\"n\">stage_to_run_to</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">return</span> <span class=\"nb\">getattr</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"n\">item</span><span
    class=\"p\">)</span>\n                <span class=\"k\">elif</span> <span class=\"n\">item</span>
    <span class=\"o\">==</span> <span class=\"s2\">&quot;precache&quot;</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_precache</span>
    <span class=\"ow\">or</span> <span class=\"p\">{}</span>\n                <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                    <span class=\"c1\">#
    Markata does not know what this is, raise</span>\n                    <span class=\"k\">raise</span>
    <span class=\"ne\">AttributeError</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;&#39;Markata&#39; object has no attribute &#39;</span><span
    class=\"si\">{</span><span class=\"n\">item</span><span class=\"si\">}</span><span
    class=\"s2\">&#39;&quot;</span><span class=\"p\">)</span>\n\n            <span
    class=\"k\">def</span> <span class=\"nf\">__rich__</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Table</span><span
    class=\"p\">:</span>\n                <span class=\"n\">grid</span> <span class=\"o\">=</span>
    <span class=\"n\">Table</span><span class=\"o\">.</span><span class=\"n\">grid</span><span
    class=\"p\">()</span>\n                <span class=\"n\">grid</span><span class=\"o\">.</span><span
    class=\"n\">add_column</span><span class=\"p\">(</span><span class=\"s2\">&quot;label&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"n\">grid</span><span class=\"o\">.</span><span
    class=\"n\">add_column</span><span class=\"p\">(</span><span class=\"s2\">&quot;value&quot;</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">for</span> <span class=\"n\">label</span><span
    class=\"p\">,</span> <span class=\"n\">value</span> <span class=\"ow\">in</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">describe</span><span
    class=\"p\">()</span><span class=\"o\">.</span><span class=\"n\">items</span><span
    class=\"p\">():</span>\n                    <span class=\"n\">grid</span><span
    class=\"o\">.</span><span class=\"n\">add_row</span><span class=\"p\">(</span><span
    class=\"n\">label</span><span class=\"p\">,</span> <span class=\"n\">value</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">return</span> <span
    class=\"n\">grid</span>\n\n            <span class=\"k\">def</span> <span class=\"nf\">bust_cache</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">:</span> <span
    class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Markata</span><span class=\"p\">:</span>\n                <span
    class=\"k\">with</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">cache</span> <span class=\"k\">as</span> <span class=\"n\">cache</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">cache</span><span
    class=\"o\">.</span><span class=\"n\">clear</span><span class=\"p\">()</span>\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"nf\">configure</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Markata</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">sys</span><span class=\"o\">.</span><span class=\"n\">path</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">os</span><span class=\"o\">.</span><span class=\"n\">getcwd</span><span
    class=\"p\">())</span>\n                <span class=\"c1\"># self.config = {**DEFUALT_CONFIG,
    **standard_config.load(&quot;markata&quot;)}</span>\n                <span class=\"c1\">#
    if isinstance(self.config[&quot;glob_patterns&quot;], str):</span>\n                <span
    class=\"c1\">#     self.config[&quot;glob_patterns&quot;] = self.config[&quot;glob_patterns&quot;].split(&quot;,&quot;)</span>\n
    \               <span class=\"c1\"># elif isinstance(self.config[&quot;glob_patterns&quot;],
    list):</span>\n                <span class=\"c1\">#     self.config[&quot;glob_patterns&quot;]
    = list(self.config[&quot;glob_patterns&quot;])</span>\n                <span class=\"c1\">#
    else:</span>\n                <span class=\"c1\">#     raise TypeError(&quot;glob_patterns
    must be list or str&quot;)</span>\n                <span class=\"c1\"># self.glob_patterns
    = self.config[&quot;glob_patterns&quot;]</span>\n\n                <span class=\"c1\">#
    self.hooks = self.config[&quot;hooks&quot;]</span>\n\n                <span class=\"c1\">#
    if &quot;disabled_hooks&quot; not in self.config:</span>\n                <span
    class=\"c1\">#     self.disabled_hooks = [&quot;&quot;]</span>\n                <span
    class=\"c1\"># if isinstance(self.config[&quot;disabled_hooks&quot;], str):</span>\n
    \               <span class=\"c1\">#     self.disabled_hooks = self.config[&quot;disabled_hooks&quot;].split(&quot;,&quot;)</span>\n
    \               <span class=\"c1\"># if isinstance(self.config[&quot;disabled_hooks&quot;],
    list):</span>\n                <span class=\"c1\">#     self.disabled_hooks =
    self.config[&quot;disabled_hooks&quot;]</span>\n\n                <span class=\"c1\">#
    if not self.config.get(&quot;output_dir&quot;, &quot;markout&quot;).endswith(</span>\n
    \               <span class=\"c1\">#     self.config.get(&quot;path_prefix&quot;,
    &quot;&quot;)</span>\n                <span class=\"c1\"># ):</span>\n                <span
    class=\"c1\">#     self.config[&quot;output_dir&quot;] = (</span>\n                <span
    class=\"c1\">#         self.config.get(&quot;output_dir&quot;, &quot;markout&quot;)
    +</span>\n                <span class=\"c1\">#         &quot;/&quot; +</span>\n
    \               <span class=\"c1\">#         self.config.get(&quot;path_prefix&quot;,
    &quot;&quot;).rstrip(&quot;/&quot;)</span>\n                <span class=\"c1\">#
    \    )</span>\n                <span class=\"c1\"># if (</span>\n                <span
    class=\"c1\">#     len((output_split := self.config.get(&quot;output_dir&quot;,
    &quot;markout&quot;).split(&quot;/&quot;))) &gt;</span>\n                <span
    class=\"c1\">#     1</span>\n                <span class=\"c1\"># ):</span>\n
    \               <span class=\"c1\">#     if &quot;path_prefix&quot; not in self.config.keys():</span>\n
    \               <span class=\"c1\">#         self.config[&quot;path_prefix&quot;]
    = &quot;/&quot;.join(output_split[1:]) + &quot;/&quot;</span>\n                <span
    class=\"c1\"># if not self.config.get(&quot;path_prefix&quot;, &quot;&quot;).endswith(&quot;/&quot;):</span>\n
    \               <span class=\"c1\">#     self.config[&quot;path_prefix&quot;]
    = self.config.get(&quot;path_prefix&quot;, &quot;&quot;) + &quot;/&quot;</span>\n\n
    \               <span class=\"c1\"># self.config[&quot;output_dir&quot;] = self.config[&quot;output_dir&quot;].lstrip(&quot;/&quot;)</span>\n
    \               <span class=\"c1\"># self.config[&quot;path_prefix&quot;] = self.config[&quot;path_prefix&quot;].lstrip(&quot;/&quot;)</span>\n\n
    \               <span class=\"k\">try</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">default_index</span> <span class=\"o\">=</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">hooks_conf</span><span class=\"o\">.</span><span
    class=\"n\">hooks</span><span class=\"o\">.</span><span class=\"n\">index</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;default&quot;</span><span class=\"p\">)</span>\n
    \                   <span class=\"n\">hooks</span> <span class=\"o\">=</span>
    <span class=\"p\">[</span>\n                        <span class=\"o\">*</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">hooks_conf</span><span
    class=\"o\">.</span><span class=\"n\">hooks</span><span class=\"p\">[:</span><span
    class=\"n\">default_index</span><span class=\"p\">],</span>\n                        <span
    class=\"o\">*</span><span class=\"n\">DEFAULT_HOOKS</span><span class=\"p\">,</span>\n
    \                       <span class=\"o\">*</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">hooks_conf</span><span class=\"o\">.</span><span
    class=\"n\">hooks</span><span class=\"p\">[</span><span class=\"n\">default_index</span>
    <span class=\"o\">+</span> <span class=\"mi\">1</span> <span class=\"p\">:],</span>\n
    \                   <span class=\"p\">]</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">hooks</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
    \                       <span class=\"n\">hook</span> <span class=\"k\">for</span>
    <span class=\"n\">hook</span> <span class=\"ow\">in</span> <span class=\"n\">hooks</span>
    <span class=\"k\">if</span> <span class=\"n\">hook</span> <span class=\"ow\">not</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">disabled_hooks</span>\n
    \                   <span class=\"p\">]</span>\n                <span class=\"k\">except</span>
    <span class=\"ne\">ValueError</span><span class=\"p\">:</span>\n                    <span
    class=\"c1\"># &#39;default&#39; is not in hooks , do not replace with default_hooks</span>\n
    \                   <span class=\"k\">pass</span>\n\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_pm</span> <span class=\"o\">=</span> <span
    class=\"n\">pluggy</span><span class=\"o\">.</span><span class=\"n\">PluginManager</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;markata&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_pm</span><span class=\"o\">.</span><span class=\"n\">add_hookspecs</span><span
    class=\"p\">(</span><span class=\"n\">hookspec</span><span class=\"o\">.</span><span
    class=\"n\">MarkataSpecs</span><span class=\"p\">)</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_register_hooks</span><span
    class=\"p\">()</span>\n\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_pm</span><span class=\"o\">.</span><span class=\"n\">hook</span><span
    class=\"o\">.</span><span class=\"n\">configure</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"o\">=</span><span class=\"bp\">self</span><span
    class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"bp\">self</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"nf\">get_plugin_config</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"n\">path_or_name</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Dict</span><span
    class=\"p\">:</span>\n                <span class=\"n\">key</span> <span class=\"o\">=</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">path_or_name</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">stem</span>\n\n
    \               <span class=\"n\">config</span> <span class=\"o\">=</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">,</span> <span class=\"p\">{})</span>\n\n
    \               <span class=\"k\">if</span> <span class=\"ow\">not</span> <span
    class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">config</span><span
    class=\"p\">,</span> <span class=\"nb\">dict</span><span class=\"p\">):</span>\n
    \                   <span class=\"k\">raise</span> <span class=\"ne\">TypeError</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;must use dict&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">if</span> <span class=\"s2\">&quot;cache_expire&quot;</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">keys</span><span class=\"p\">():</span>\n
    \                   <span class=\"n\">config</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;cache_expire&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;default_cache_expire&quot;</span><span
    class=\"p\">]</span>\n                <span class=\"k\">if</span> <span class=\"s2\">&quot;config_key&quot;</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">keys</span><span class=\"p\">():</span>\n
    \                   <span class=\"n\">config</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;config_key&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"n\">key</span>\n                <span class=\"k\">return</span>
    <span class=\"n\">config</span>\n\n            <span class=\"k\">def</span> <span
    class=\"nf\">get_config</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">self</span><span class=\"p\">,</span>\n                <span class=\"n\">key</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">default</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">warn</span><span class=\"p\">:</span>
    <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"kc\">True</span><span
    class=\"p\">,</span>\n                <span class=\"n\">suggested</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"kc\">None</span><span class=\"p\">,</span>\n            <span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">key</span> <span
    class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">keys</span><span
    class=\"p\">():</span>\n                    <span class=\"k\">return</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"p\">[</span><span class=\"n\">key</span><span class=\"p\">]</span>\n                <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                    <span class=\"k\">if</span>
    <span class=\"n\">suggested</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                        <span class=\"n\">suggested</span>
    <span class=\"o\">=</span> <span class=\"n\">textwrap</span><span class=\"o\">.</span><span
    class=\"n\">dedent</span><span class=\"p\">(</span>\n                            <span
    class=\"sa\">f</span><span class=\"s2\">&quot;&quot;&quot;</span>\n<span class=\"s2\">
    \                           [markata]</span>\n<span class=\"s2\">                            </span><span
    class=\"si\">{</span><span class=\"n\">key</span><span class=\"si\">}</span><span
    class=\"s2\"> = &#39;</span><span class=\"si\">{</span><span class=\"n\">default</span><span
    class=\"si\">}</span><span class=\"s2\">&#39;</span>\n<span class=\"s2\">                            &quot;&quot;&quot;</span>\n
    \                       <span class=\"p\">)</span>\n                    <span
    class=\"k\">if</span> <span class=\"n\">warn</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">logger</span><span class=\"o\">.</span><span
    class=\"n\">warning</span><span class=\"p\">(</span>\n                            <span
    class=\"n\">textwrap</span><span class=\"o\">.</span><span class=\"n\">dedent</span><span
    class=\"p\">(</span>\n                                <span class=\"sa\">f</span><span
    class=\"s2\">&quot;&quot;&quot;</span>\n<span class=\"s2\">                                Warning
    </span><span class=\"si\">{</span><span class=\"n\">key</span><span class=\"si\">}</span><span
    class=\"s2\"> is not set in markata config, sitemap will</span>\n<span class=\"s2\">
    \                               be missing root site_name</span>\n<span class=\"s2\">
    \                               to resolve this open your markata.toml and add</span>\n\n<span
    class=\"s2\">                                </span><span class=\"si\">{</span><span
    class=\"n\">suggested</span><span class=\"si\">}</span>\n\n<span class=\"s2\">
    \                               &quot;&quot;&quot;</span>\n                            <span
    class=\"p\">),</span>\n                        <span class=\"p\">)</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">default</span>\n\n            <span
    class=\"k\">def</span> <span class=\"nf\">make_hash</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"n\">keys</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span
    class=\"p\">:</span>\n                <span class=\"kn\">import</span> <span class=\"nn\">xxhash</span>\n\n
    \               <span class=\"n\">str_keys</span> <span class=\"o\">=</span> <span
    class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">)</span> <span class=\"k\">for</span>
    <span class=\"n\">key</span> <span class=\"ow\">in</span> <span class=\"n\">keys</span><span
    class=\"p\">]</span>\n                <span class=\"nb\">hash</span> <span class=\"o\">=</span>
    <span class=\"n\">xxhash</span><span class=\"o\">.</span><span class=\"n\">xxh64</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;&quot;</span><span class=\"o\">.</span><span
    class=\"n\">join</span><span class=\"p\">(</span><span class=\"n\">str_keys</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">encode</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;utf-8&quot;</span><span class=\"p\">))</span><span
    class=\"o\">.</span><span class=\"n\">hexdigest</span><span class=\"p\">()</span>\n
    \               <span class=\"k\">return</span> <span class=\"nb\">hash</span>\n\n
    \           <span class=\"nd\">@property</span>\n            <span class=\"k\">def</span>
    <span class=\"nf\">content_dir_hash</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">hashes</span> <span class=\"o\">=</span> <span
    class=\"p\">[</span>\n                    <span class=\"n\">dirhash</span><span
    class=\"p\">(</span><span class=\"nb\">dir</span><span class=\"p\">)</span>\n
    \                   <span class=\"k\">for</span> <span class=\"nb\">dir</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">content_directories</span>\n                    <span class=\"k\">if</span>
    <span class=\"nb\">dir</span><span class=\"o\">.</span><span class=\"n\">absolute</span><span
    class=\"p\">()</span> <span class=\"o\">!=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;.&quot;</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">absolute</span><span class=\"p\">()</span>\n
    \               <span class=\"p\">]</span>\n                <span class=\"k\">return</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">make_hash</span><span
    class=\"p\">(</span><span class=\"o\">*</span><span class=\"n\">hashes</span><span
    class=\"p\">)</span>\n\n            <span class=\"nd\">@property</span>\n            <span
    class=\"k\">def</span> <span class=\"nf\">console</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Console</span><span
    class=\"p\">:</span>\n                <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_console</span>\n                <span class=\"k\">except</span>
    <span class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_console</span>
    <span class=\"o\">=</span> <span class=\"n\">Console</span><span class=\"p\">()</span>\n
    \                   <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_console</span>\n\n            <span class=\"k\">def</span>
    <span class=\"nf\">describe</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">dict</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
    class=\"p\">]:</span>\n                <span class=\"k\">return</span> <span class=\"p\">{</span><span
    class=\"s2\">&quot;version&quot;</span><span class=\"p\">:</span> <span class=\"n\">__version__</span><span
    class=\"p\">}</span>\n\n            <span class=\"k\">def</span> <span class=\"nf\">_to_dict</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">:</span> <span
    class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"nb\">dict</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"n\">Iterable</span><span class=\"p\">]:</span>\n
    \               <span class=\"k\">return</span> <span class=\"p\">{</span><span
    class=\"s2\">&quot;config&quot;</span><span class=\"p\">:</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;articles&quot;</span><span class=\"p\">:</span> <span
    class=\"p\">[</span><span class=\"n\">a</span><span class=\"o\">.</span><span
    class=\"n\">to_dict</span><span class=\"p\">()</span> <span class=\"k\">for</span>
    <span class=\"n\">a</span> <span class=\"ow\">in</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">articles</span><span class=\"p\">]}</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"nf\">to_dict</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">:</span> <span
    class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"nb\">dict</span><span class=\"p\">:</span>\n                <span
    class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_to_dict</span><span class=\"p\">()</span>\n\n            <span class=\"k\">def</span>
    <span class=\"nf\">to_json</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n
    \               <span class=\"kn\">import</span> <span class=\"nn\">json</span>\n\n
    \               <span class=\"k\">return</span> <span class=\"n\">json</span><span
    class=\"o\">.</span><span class=\"n\">dumps</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">to_dict</span><span
    class=\"p\">(),</span> <span class=\"n\">indent</span><span class=\"o\">=</span><span
    class=\"mi\">4</span><span class=\"p\">,</span> <span class=\"n\">sort_keys</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
    class=\"n\">default</span><span class=\"o\">=</span><span class=\"nb\">str</span><span
    class=\"p\">)</span>\n\n            <span class=\"k\">def</span> <span class=\"nf\">_register_hooks</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">:</span> <span
    class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
    class=\"n\">sys</span><span class=\"o\">.</span><span class=\"n\">path</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">os</span><span class=\"o\">.</span><span class=\"n\">getcwd</span><span
    class=\"p\">())</span>\n                <span class=\"k\">for</span> <span class=\"n\">hook</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">hooks_conf</span><span class=\"o\">.</span><span class=\"n\">hooks</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                       <span class=\"c1\"># module style plugins</span>\n                        <span
    class=\"n\">plugin</span> <span class=\"o\">=</span> <span class=\"n\">importlib</span><span
    class=\"o\">.</span><span class=\"n\">import_module</span><span class=\"p\">(</span><span
    class=\"n\">hook</span><span class=\"p\">)</span>\n                    <span class=\"k\">except</span>
    <span class=\"ne\">ModuleNotFoundError</span> <span class=\"k\">as</span> <span
    class=\"n\">e</span><span class=\"p\">:</span>\n                        <span
    class=\"c1\"># class style plugins</span>\n                        <span class=\"k\">if</span>
    <span class=\"s2\">&quot;.&quot;</span> <span class=\"ow\">in</span> <span class=\"n\">hook</span><span
    class=\"p\">:</span>\n                            <span class=\"k\">try</span><span
    class=\"p\">:</span>\n                                <span class=\"n\">mod</span>
    <span class=\"o\">=</span> <span class=\"n\">importlib</span><span class=\"o\">.</span><span
    class=\"n\">import_module</span><span class=\"p\">(</span><span class=\"s2\">&quot;.&quot;</span><span
    class=\"o\">.</span><span class=\"n\">join</span><span class=\"p\">(</span><span
    class=\"n\">hook</span><span class=\"o\">.</span><span class=\"n\">split</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;.&quot;</span><span class=\"p\">)[:</span><span
    class=\"o\">-</span><span class=\"mi\">1</span><span class=\"p\">]))</span>\n
    \                               <span class=\"n\">plugin</span> <span class=\"o\">=</span>
    <span class=\"nb\">getattr</span><span class=\"p\">(</span><span class=\"n\">mod</span><span
    class=\"p\">,</span> <span class=\"n\">hook</span><span class=\"o\">.</span><span
    class=\"n\">split</span><span class=\"p\">(</span><span class=\"s2\">&quot;.&quot;</span><span
    class=\"p\">)[</span><span class=\"o\">-</span><span class=\"mi\">1</span><span
    class=\"p\">])</span>\n                            <span class=\"k\">except</span>
    <span class=\"ne\">ModuleNotFoundError</span> <span class=\"k\">as</span> <span
    class=\"n\">e</span><span class=\"p\">:</span>\n                                <span
    class=\"k\">raise</span> <span class=\"ne\">ModuleNotFoundError</span><span class=\"p\">(</span>\n
    \                                   <span class=\"sa\">f</span><span class=\"s2\">&quot;module
    </span><span class=\"si\">{</span><span class=\"n\">hook</span><span class=\"si\">}</span><span
    class=\"s2\"> not found</span><span class=\"se\">\\n</span><span class=\"si\">{</span><span
    class=\"n\">sys</span><span class=\"o\">.</span><span class=\"n\">path</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span>\n                                <span
    class=\"p\">)</span> <span class=\"kn\">from</span> <span class=\"nn\">e</span>\n
    \                       <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \                           <span class=\"k\">raise</span> <span class=\"n\">e</span>\n\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_pm</span><span class=\"o\">.</span><span class=\"n\">register</span><span
    class=\"p\">(</span><span class=\"n\">plugin</span><span class=\"p\">)</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"fm\">__iter__</span><span
    class=\"p\">(</span>\n                <span class=\"bp\">self</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span> <span
    class=\"n\">description</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;working...&quot;</span>\n
    \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Iterable</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;Markata.Post&quot;</span><span class=\"p\">]:</span>\n
    \               <span class=\"n\">articles</span><span class=\"p\">:</span> <span
    class=\"n\">Iterable</span><span class=\"p\">[</span><span class=\"n\">Markata</span><span
    class=\"o\">.</span><span class=\"n\">Post</span><span class=\"p\">]</span> <span
    class=\"o\">=</span> <span class=\"n\">track</span><span class=\"p\">(</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">articles</span><span class=\"p\">,</span>\n                    <span
    class=\"n\">description</span><span class=\"o\">=</span><span class=\"n\">description</span><span
    class=\"p\">,</span>\n                    <span class=\"n\">transient</span><span
    class=\"o\">=</span><span class=\"kc\">False</span><span class=\"p\">,</span>\n
    \                   <span class=\"n\">console</span><span class=\"o\">=</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"p\">,</span>\n                <span class=\"p\">)</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">articles</span>\n\n            <span
    class=\"k\">def</span> <span class=\"nf\">iter_articles</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">description</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Iterable</span><span class=\"p\">[</span><span class=\"n\">Markata</span><span
    class=\"o\">.</span><span class=\"n\">Post</span><span class=\"p\">]:</span>\n
    \               <span class=\"n\">articles</span><span class=\"p\">:</span> <span
    class=\"n\">Iterable</span><span class=\"p\">[</span><span class=\"n\">Markata</span><span
    class=\"o\">.</span><span class=\"n\">Post</span><span class=\"p\">]</span> <span
    class=\"o\">=</span> <span class=\"n\">track</span><span class=\"p\">(</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">articles</span><span class=\"p\">,</span>\n                    <span
    class=\"n\">description</span><span class=\"o\">=</span><span class=\"n\">description</span><span
    class=\"p\">,</span>\n                    <span class=\"n\">transient</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \                   <span class=\"n\">console</span><span class=\"o\">=</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"p\">,</span>\n                <span class=\"p\">)</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">articles</span>\n\n            <span
    class=\"k\">def</span> <span class=\"nf\">teardown</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Markata</span><span
    class=\"p\">:</span>\n<span class=\"w\">                </span><span class=\"sd\">&quot;&quot;&quot;give
    special access to the teardown lifecycle method&quot;&quot;&quot;</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_pm</span><span
    class=\"o\">.</span><span class=\"n\">hook</span><span class=\"o\">.</span><span
    class=\"n\">teardown</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"bp\">self</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"nf\">run</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">lifecycle</span><span class=\"p\">:</span>
    <span class=\"n\">LifeCycle</span> <span class=\"o\">=</span> <span class=\"kc\">None</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Markata</span><span
    class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"n\">lifecycle</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">lifecycle</span> <span class=\"o\">=</span>
    <span class=\"nb\">max</span><span class=\"p\">(</span><span class=\"n\">LifeCycle</span><span
    class=\"o\">.</span><span class=\"n\">_member_map_</span><span class=\"o\">.</span><span
    class=\"n\">values</span><span class=\"p\">())</span>\n\n                <span
    class=\"k\">if</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">lifecycle</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
    class=\"p\">):</span>\n                    <span class=\"n\">lifecycle</span>
    <span class=\"o\">=</span> <span class=\"n\">LifeCycle</span><span class=\"p\">[</span><span
    class=\"n\">lifecycle</span><span class=\"p\">]</span>\n\n                <span
    class=\"n\">stages_to_run</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
    \                   <span class=\"n\">m</span>\n                    <span class=\"k\">for</span>
    <span class=\"n\">m</span> <span class=\"ow\">in</span> <span class=\"n\">LifeCycle</span><span
    class=\"o\">.</span><span class=\"n\">_member_map_</span>\n                    <span
    class=\"k\">if</span> <span class=\"p\">(</span><span class=\"n\">LifeCycle</span><span
    class=\"p\">[</span><span class=\"n\">m</span><span class=\"p\">]</span> <span
    class=\"o\">&lt;=</span> <span class=\"n\">lifecycle</span><span class=\"p\">)</span>
    <span class=\"ow\">and</span> <span class=\"p\">(</span><span class=\"n\">m</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">stages_ran</span><span class=\"p\">)</span>\n
    \               <span class=\"p\">]</span>\n\n                <span class=\"k\">if</span>
    <span class=\"ow\">not</span> <span class=\"n\">stages_to_run</span><span class=\"p\">:</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
    class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">lifecycle</span><span class=\"o\">.</span><span
    class=\"n\">name</span><span class=\"si\">}</span><span class=\"s2\"> already
    ran&quot;</span><span class=\"p\">)</span>\n                    <span class=\"k\">return</span>
    <span class=\"bp\">self</span>\n\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">log</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;running </span><span class=\"si\">{</span><span class=\"n\">stages_to_run</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">for</span> <span class=\"n\">stage</span> <span
    class=\"ow\">in</span> <span class=\"n\">stages_to_run</span><span class=\"p\">:</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
    class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">stage</span><span class=\"si\">}</span><span
    class=\"s2\"> running&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"nb\">getattr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_pm</span><span class=\"o\">.</span><span
    class=\"n\">hook</span><span class=\"p\">,</span> <span class=\"n\">stage</span><span
    class=\"p\">)(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"bp\">self</span><span class=\"p\">)</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">stages_ran</span><span
    class=\"o\">.</span><span class=\"n\">add</span><span class=\"p\">(</span><span
    class=\"n\">stage</span><span class=\"p\">)</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">log</span><span class=\"p\">(</span><span
    class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">stage</span><span class=\"si\">}</span><span class=\"s2\"> complete&quot;</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">with</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">cache</span> <span class=\"k\">as</span>
    <span class=\"n\">cache</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">hits</span><span class=\"p\">,</span> <span class=\"n\">misses</span>
    <span class=\"o\">=</span> <span class=\"n\">cache</span><span class=\"o\">.</span><span
    class=\"n\">stats</span><span class=\"p\">()</span>\n\n                <span class=\"k\">if</span>
    <span class=\"n\">hits</span> <span class=\"o\">+</span> <span class=\"n\">misses</span>
    <span class=\"o\">&gt;</span> <span class=\"mi\">0</span><span class=\"p\">:</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
    class=\"p\">(</span>\n                        <span class=\"sa\">f</span><span
    class=\"s2\">&quot;lifetime cache hit rate </span><span class=\"si\">{</span><span
    class=\"nb\">round</span><span class=\"p\">(</span><span class=\"n\">hits</span><span
    class=\"o\">/</span><span class=\"w\"> </span><span class=\"p\">(</span><span
    class=\"n\">hits</span><span class=\"w\"> </span><span class=\"o\">+</span><span
    class=\"w\"> </span><span class=\"n\">misses</span><span class=\"p\">)</span><span
    class=\"o\">*</span><span class=\"mi\">100</span><span class=\"p\">,</span><span
    class=\"w\"> </span><span class=\"mi\">2</span><span class=\"p\">)</span><span
    class=\"si\">}</span><span class=\"s2\">%&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"p\">)</span>\n\n                <span class=\"k\">if</span>
    <span class=\"n\">misses</span> <span class=\"o\">&gt;</span> <span class=\"mi\">0</span><span
    class=\"p\">:</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">log</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;lifetime cache hits/misses </span><span class=\"si\">{</span><span
    class=\"n\">hits</span><span class=\"si\">}</span><span class=\"s2\">/</span><span
    class=\"si\">{</span><span class=\"n\">misses</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span><span class=\"p\">)</span>\n\n                <span
    class=\"n\">hits</span> <span class=\"o\">-=</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">init_cache_stats</span><span class=\"p\">[</span><span
    class=\"mi\">0</span><span class=\"p\">]</span>\n                <span class=\"n\">misses</span>
    <span class=\"o\">-=</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">init_cache_stats</span><span class=\"p\">[</span><span class=\"mi\">1</span><span
    class=\"p\">]</span>\n\n                <span class=\"k\">if</span> <span class=\"n\">hits</span>
    <span class=\"o\">+</span> <span class=\"n\">misses</span> <span class=\"o\">&gt;</span>
    <span class=\"mi\">0</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">log</span><span class=\"p\">(</span>\n                        <span
    class=\"sa\">f</span><span class=\"s2\">&quot;run cache hit rate </span><span
    class=\"si\">{</span><span class=\"nb\">round</span><span class=\"p\">(</span><span
    class=\"n\">hits</span><span class=\"o\">/</span><span class=\"w\"> </span><span
    class=\"p\">(</span><span class=\"n\">hits</span><span class=\"w\"> </span><span
    class=\"o\">+</span><span class=\"w\"> </span><span class=\"n\">misses</span><span
    class=\"p\">)</span><span class=\"o\">*</span><span class=\"mi\">100</span><span
    class=\"p\">,</span><span class=\"w\"> </span><span class=\"mi\">2</span><span
    class=\"p\">)</span><span class=\"si\">}</span><span class=\"s2\">%&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"p\">)</span>\n\n                <span
    class=\"k\">if</span> <span class=\"n\">misses</span> <span class=\"o\">&gt;</span>
    <span class=\"mi\">0</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">log</span><span class=\"p\">(</span><span
    class=\"sa\">f</span><span class=\"s2\">&quot;run cache hits/misses </span><span
    class=\"si\">{</span><span class=\"n\">hits</span><span class=\"si\">}</span><span
    class=\"s2\">/</span><span class=\"si\">{</span><span class=\"n\">misses</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">)</span>\n\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"nf\">filter</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">:</span> <span
    class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span> <span class=\"nb\">filter</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"nb\">list</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">def</span> <span class=\"nf\">evalr</span><span
    class=\"p\">(</span><span class=\"n\">a</span><span class=\"p\">:</span> <span
    class=\"n\">Markata</span><span class=\"o\">.</span><span class=\"n\">Post</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                       <span class=\"k\">return</span> <span class=\"nb\">eval</span><span
    class=\"p\">(</span>\n                            <span class=\"nb\">filter</span><span
    class=\"p\">,</span>\n                            <span class=\"p\">{</span><span
    class=\"o\">**</span><span class=\"n\">a</span><span class=\"o\">.</span><span
    class=\"n\">to_dict</span><span class=\"p\">(),</span> <span class=\"s2\">&quot;timedelta&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">timedelta</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;post&quot;</span><span class=\"p\">:</span> <span class=\"n\">a</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;m&quot;</span><span class=\"p\">:</span>
    <span class=\"bp\">self</span><span class=\"p\">},</span>\n                            <span
    class=\"p\">{},</span>\n                        <span class=\"p\">)</span>\n                    <span
    class=\"k\">except</span> <span class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n
    \                       <span class=\"k\">return</span> <span class=\"nb\">eval</span><span
    class=\"p\">(</span>\n                            <span class=\"nb\">filter</span><span
    class=\"p\">,</span>\n                            <span class=\"p\">{</span><span
    class=\"o\">**</span><span class=\"n\">a</span><span class=\"o\">.</span><span
    class=\"n\">to_dict</span><span class=\"p\">(),</span> <span class=\"s2\">&quot;timedelta&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">timedelta</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;post&quot;</span><span class=\"p\">:</span> <span class=\"n\">a</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;m&quot;</span><span class=\"p\">:</span>
    <span class=\"bp\">self</span><span class=\"p\">},</span>\n                            <span
    class=\"p\">{},</span>\n                        <span class=\"p\">)</span>\n\n
    \               <span class=\"k\">return</span> <span class=\"p\">[</span><span
    class=\"n\">a</span> <span class=\"k\">for</span> <span class=\"n\">a</span> <span
    class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">articles</span> <span class=\"k\">if</span> <span class=\"n\">evalr</span><span
    class=\"p\">(</span><span class=\"n\">a</span><span class=\"p\">)]</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"nf\">map</span><span class=\"p\">(</span>\n
    \               <span class=\"bp\">self</span><span class=\"p\">:</span> <span
    class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span>\n                <span
    class=\"n\">func</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;title&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"nb\">filter</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;True&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">sort</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;True&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">reverse</span><span class=\"p\">:</span>
    <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"kc\">True</span><span
    class=\"p\">,</span>\n                <span class=\"o\">*</span><span class=\"n\">args</span><span
    class=\"p\">:</span> <span class=\"nb\">tuple</span><span class=\"p\">,</span>\n
    \               <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">:</span> <span class=\"nb\">dict</span><span class=\"p\">,</span>\n
    \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">list</span><span
    class=\"p\">:</span>\n                <span class=\"kn\">import</span> <span class=\"nn\">copy</span>\n\n
    \               <span class=\"k\">def</span> <span class=\"nf\">try_sort</span><span
    class=\"p\">(</span><span class=\"n\">a</span><span class=\"p\">:</span> <span
    class=\"n\">Any</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"nb\">int</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">if</span> <span class=\"s2\">&quot;datetime&quot;</span> <span class=\"ow\">in</span>
    <span class=\"n\">sort</span><span class=\"o\">.</span><span class=\"n\">lower</span><span
    class=\"p\">():</span>\n                        <span class=\"k\">return</span>
    <span class=\"n\">a</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"n\">sort</span><span class=\"p\">,</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"p\">(</span><span class=\"mi\">1970</span><span class=\"p\">,</span> <span
    class=\"mi\">1</span><span class=\"p\">,</span> <span class=\"mi\">1</span><span
    class=\"p\">))</span>\n\n                    <span class=\"k\">if</span> <span
    class=\"s2\">&quot;date&quot;</span> <span class=\"ow\">in</span> <span class=\"n\">sort</span><span
    class=\"o\">.</span><span class=\"n\">lower</span><span class=\"p\">():</span>\n
    \                       <span class=\"k\">return</span> <span class=\"n\">a</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"n\">sort</span><span class=\"p\">,</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">date</span><span class=\"p\">(</span><span
    class=\"mi\">1970</span><span class=\"p\">,</span> <span class=\"mi\">1</span><span
    class=\"p\">,</span> <span class=\"mi\">1</span><span class=\"p\">))</span>\n\n
    \                   <span class=\"k\">try</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">value</span> <span class=\"o\">=</span> <span class=\"nb\">eval</span><span
    class=\"p\">(</span><span class=\"n\">sort</span><span class=\"p\">,</span> <span
    class=\"n\">a</span><span class=\"o\">.</span><span class=\"n\">to_dict</span><span
    class=\"p\">(),</span> <span class=\"p\">{})</span>\n                    <span
    class=\"k\">except</span> <span class=\"ne\">NameError</span><span class=\"p\">:</span>\n
    \                       <span class=\"k\">return</span> <span class=\"o\">-</span><span
    class=\"mi\">1</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">value</span>\n                    <span class=\"k\">try</span><span
    class=\"p\">:</span>\n                        <span class=\"k\">return</span>
    <span class=\"nb\">int</span><span class=\"p\">(</span><span class=\"n\">value</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">except</span> <span
    class=\"ne\">TypeError</span><span class=\"p\">:</span>\n                        <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                            <span
    class=\"k\">return</span> <span class=\"nb\">int</span><span class=\"p\">(</span><span
    class=\"n\">value</span><span class=\"o\">.</span><span class=\"n\">timestamp</span><span
    class=\"p\">())</span>\n                        <span class=\"k\">except</span>
    <span class=\"ne\">Exception</span><span class=\"p\">:</span>\n                            <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                                <span
    class=\"k\">return</span> <span class=\"nb\">int</span><span class=\"p\">(</span>\n
    \                                   <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">combine</span><span
    class=\"p\">(</span>\n                                        <span class=\"n\">value</span><span
    class=\"p\">,</span>\n                                        <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">min</span><span class=\"o\">.</span><span class=\"n\">time</span><span
    class=\"p\">(),</span>\n                                    <span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">timestamp</span><span class=\"p\">(),</span>\n
    \                               <span class=\"p\">)</span>\n                            <span
    class=\"k\">except</span> <span class=\"ne\">Exception</span><span class=\"p\">:</span>\n
    \                               <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                                   <span class=\"k\">return</span> <span class=\"nb\">sum</span><span
    class=\"p\">([</span><span class=\"nb\">ord</span><span class=\"p\">(</span><span
    class=\"n\">c</span><span class=\"p\">)</span> <span class=\"k\">for</span> <span
    class=\"n\">c</span> <span class=\"ow\">in</span> <span class=\"nb\">str</span><span
    class=\"p\">(</span><span class=\"n\">value</span><span class=\"p\">)])</span>\n
    \                               <span class=\"k\">except</span> <span class=\"ne\">Exception</span><span
    class=\"p\">:</span>\n                                    <span class=\"k\">return</span>
    <span class=\"o\">-</span><span class=\"mi\">1</span>\n\n                <span
    class=\"n\">articles</span> <span class=\"o\">=</span> <span class=\"n\">copy</span><span
    class=\"o\">.</span><span class=\"n\">copy</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">articles</span><span
    class=\"p\">)</span>\n                <span class=\"n\">articles</span><span class=\"o\">.</span><span
    class=\"n\">sort</span><span class=\"p\">(</span><span class=\"n\">key</span><span
    class=\"o\">=</span><span class=\"n\">try_sort</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">reverse</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">articles</span><span
    class=\"o\">.</span><span class=\"n\">reverse</span><span class=\"p\">()</span>\n\n
    \               <span class=\"k\">try</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">posts</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
    \                       <span class=\"nb\">eval</span><span class=\"p\">(</span>\n
    \                           <span class=\"n\">func</span><span class=\"p\">,</span>\n
    \                           <span class=\"p\">{</span><span class=\"o\">**</span><span
    class=\"n\">a</span><span class=\"o\">.</span><span class=\"n\">to_dict</span><span
    class=\"p\">(),</span> <span class=\"s2\">&quot;timedelta&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">timedelta</span><span class=\"p\">,</span> <span class=\"s2\">&quot;post&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">a</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;m&quot;</span><span class=\"p\">:</span> <span class=\"bp\">self</span><span
    class=\"p\">},</span>\n                            <span class=\"p\">{},</span>\n
    \                       <span class=\"p\">)</span>\n                        <span
    class=\"k\">for</span> <span class=\"n\">a</span> <span class=\"ow\">in</span>
    <span class=\"n\">articles</span>\n                        <span class=\"k\">if</span>
    <span class=\"nb\">eval</span><span class=\"p\">(</span>\n                            <span
    class=\"nb\">filter</span><span class=\"p\">,</span>\n                            <span
    class=\"p\">{</span><span class=\"o\">**</span><span class=\"n\">a</span><span
    class=\"o\">.</span><span class=\"n\">to_dict</span><span class=\"p\">(),</span>
    <span class=\"s2\">&quot;timedelta&quot;</span><span class=\"p\">:</span> <span
    class=\"n\">timedelta</span><span class=\"p\">,</span> <span class=\"s2\">&quot;post&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">a</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;m&quot;</span><span class=\"p\">:</span> <span class=\"bp\">self</span><span
    class=\"p\">},</span>\n                            <span class=\"p\">{},</span>\n
    \                       <span class=\"p\">)</span>\n                    <span
    class=\"p\">]</span>\n\n                <span class=\"k\">except</span> <span
    class=\"ne\">NameError</span> <span class=\"k\">as</span> <span class=\"n\">e</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">variable</span> <span
    class=\"o\">=</span> <span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">e</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">split</span><span class=\"p\">(</span><span class=\"s2\">&quot;&#39;&quot;</span><span
    class=\"p\">)[</span><span class=\"mi\">1</span><span class=\"p\">]</span>\n\n
    \                   <span class=\"n\">missing_in_posts</span> <span class=\"o\">=</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">map</span><span
    class=\"p\">(</span>\n                        <span class=\"s2\">&quot;path&quot;</span><span
    class=\"p\">,</span>\n                        <span class=\"nb\">filter</span><span
    class=\"o\">=</span><span class=\"sa\">f</span><span class=\"s1\">&#39;&quot;</span><span
    class=\"si\">{</span><span class=\"n\">variable</span><span class=\"si\">}</span><span
    class=\"s1\">&quot; not in post.keys()&#39;</span><span class=\"p\">,</span>\n
    \                   <span class=\"p\">)</span>\n                    <span class=\"n\">message</span>
    <span class=\"o\">=</span> <span class=\"p\">(</span>\n                        <span
    class=\"sa\">f</span><span class=\"s2\">&quot;variable: &#39;</span><span class=\"si\">{</span><span
    class=\"n\">variable</span><span class=\"si\">}</span><span class=\"s2\">&#39;
    is missing in </span><span class=\"si\">{</span><span class=\"nb\">len</span><span
    class=\"p\">(</span><span class=\"n\">missing_in_posts</span><span class=\"p\">)</span><span
    class=\"si\">}</span><span class=\"s2\"> posts&quot;</span>\n                    <span
    class=\"p\">)</span>\n                    <span class=\"k\">if</span> <span class=\"nb\">len</span><span
    class=\"p\">(</span><span class=\"n\">missing_in_posts</span><span class=\"p\">)</span>
    <span class=\"o\">&gt;</span> <span class=\"mi\">10</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">message</span> <span class=\"o\">+=</span>
    <span class=\"p\">(</span>\n                            <span class=\"sa\">f</span><span
    class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span class=\"s2\">first
    10 paths to posts missing </span><span class=\"si\">{</span><span class=\"n\">variable</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span>\n                            <span
    class=\"sa\">f</span><span class=\"s2\">&quot;[</span><span class=\"si\">{</span><span
    class=\"s1\">&#39;,&#39;</span><span class=\"o\">.</span><span class=\"n\">join</span><span
    class=\"p\">([</span><span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">p</span><span class=\"p\">)</span><span class=\"w\"> </span><span
    class=\"k\">for</span><span class=\"w\"> </span><span class=\"n\">p</span><span
    class=\"w\"> </span><span class=\"ow\">in</span><span class=\"w\"> </span><span
    class=\"n\">missing_in_posts</span><span class=\"p\">[:</span><span class=\"mi\">10</span><span
    class=\"p\">]])</span><span class=\"si\">}</span><span class=\"s2\">...&quot;</span>\n
    \                       <span class=\"p\">)</span>\n                    <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">message</span> <span class=\"o\">+=</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span class=\"s2\">paths
    to posts missing </span><span class=\"si\">{</span><span class=\"n\">variable</span><span
    class=\"si\">}</span><span class=\"s2\"> </span><span class=\"si\">{</span><span
    class=\"n\">missing_in_posts</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span>\n\n
    \                   <span class=\"k\">raise</span> <span class=\"n\">MissingFrontMatter</span><span
    class=\"p\">(</span><span class=\"n\">message</span><span class=\"p\">)</span>\n\n
    \               <span class=\"k\">return</span> <span class=\"n\">posts</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"nf\">first</span><span
    class=\"p\">(</span>\n                <span class=\"bp\">self</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span>\n                <span
    class=\"nb\">filter</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;True&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">sort</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;True&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">reverse</span><span class=\"p\">:</span>
    <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"kc\">True</span><span
    class=\"p\">,</span>\n                <span class=\"o\">*</span><span class=\"n\">args</span><span
    class=\"p\">:</span> <span class=\"nb\">tuple</span><span class=\"p\">,</span>\n
    \               <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">:</span> <span class=\"nb\">dict</span><span class=\"p\">,</span>\n
    \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">list</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">map</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;post&quot;</span><span class=\"p\">,</span> <span class=\"nb\">filter</span><span
    class=\"p\">,</span> <span class=\"n\">sort</span><span class=\"p\">,</span> <span
    class=\"n\">reverse</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"n\">args</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"n\">kwargs</span><span class=\"p\">)[</span><span class=\"mi\">0</span><span
    class=\"p\">]</span>\n\n            <span class=\"k\">def</span> <span class=\"nf\">last</span><span
    class=\"p\">(</span>\n                <span class=\"bp\">self</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span>\n                <span
    class=\"nb\">filter</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;True&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">sort</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;True&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">reverse</span><span class=\"p\">:</span>
    <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"kc\">True</span><span
    class=\"p\">,</span>\n                <span class=\"o\">*</span><span class=\"n\">args</span><span
    class=\"p\">:</span> <span class=\"nb\">tuple</span><span class=\"p\">,</span>\n
    \               <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">:</span> <span class=\"nb\">dict</span><span class=\"p\">,</span>\n
    \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">list</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">map</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;post&quot;</span><span class=\"p\">,</span> <span class=\"nb\">filter</span><span
    class=\"p\">,</span> <span class=\"n\">sort</span><span class=\"p\">,</span> <span
    class=\"n\">reverse</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"n\">args</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"n\">kwargs</span><span class=\"p\">)[</span><span class=\"o\">-</span><span
    class=\"mi\">1</span><span class=\"p\">]</span>\n\n            <span class=\"k\">def</span>
    <span class=\"nf\">one</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"nb\">filter</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;True&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"o\">*</span><span class=\"n\">args</span><span
    class=\"p\">:</span> <span class=\"nb\">tuple</span><span class=\"p\">,</span>\n
    \               <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">:</span> <span class=\"nb\">dict</span><span class=\"p\">,</span>\n
    \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">list</span><span
    class=\"p\">:</span>\n                <span class=\"n\">posts</span> <span class=\"o\">=</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">map</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;post&quot;</span><span class=\"p\">,</span>
    <span class=\"nb\">filter</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"n\">args</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"n\">kwargs</span><span class=\"p\">)</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">len</span><span class=\"p\">(</span><span class=\"n\">posts</span><span
    class=\"p\">)</span> <span class=\"o\">&gt;</span> <span class=\"mi\">1</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">raise</span> <span
    class=\"n\">TooManyPosts</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;found </span><span class=\"si\">{</span><span class=\"nb\">len</span><span
    class=\"p\">(</span><span class=\"n\">posts</span><span class=\"p\">)</span><span
    class=\"si\">}</span><span class=\"s2\"> posts, expected 1. </span><span class=\"si\">{</span><span
    class=\"n\">posts</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"k\">if</span> <span class=\"nb\">len</span><span
    class=\"p\">(</span><span class=\"n\">posts</span><span class=\"p\">)</span> <span
    class=\"o\">==</span> <span class=\"mi\">0</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">raise</span> <span class=\"n\">NoPosts</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">posts</span><span
    class=\"p\">[</span><span class=\"mi\">0</span><span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='load_ipython_extension' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>load_ipython_extension <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">load_ipython_extension
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
    <span class=\"nf\">load_ipython_extension</span><span class=\"p\">(</span><span
    class=\"n\">ipython</span><span class=\"p\">):</span>\n            <span class=\"n\">ipython</span><span
    class=\"o\">.</span><span class=\"n\">user_ns</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;m&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"n\">Markata</span><span class=\"p\">()</span>\n            <span
    class=\"n\">ipython</span><span class=\"o\">.</span><span class=\"n\">user_ns</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;markata&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"n\">ipython</span><span class=\"o\">.</span><span
    class=\"n\">user_ns</span><span class=\"p\">[</span><span class=\"s2\">&quot;m&quot;</span><span
    class=\"p\">]</span>\n            <span class=\"n\">ipython</span><span class=\"o\">.</span><span
    class=\"n\">user_ns</span><span class=\"p\">[</span><span class=\"s2\">&quot;markata&quot;</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"n\">ipython</span><span
    class=\"o\">.</span><span class=\"n\">user_ns</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;m&quot;</span><span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='__init__' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'><strong>init</strong> <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\"><strong>init</strong>
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
    <span class=\"fm\">__init__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">console</span><span class=\"p\">:</span> <span class=\"n\">Console</span>
    <span class=\"o\">=</span> <span class=\"kc\">None</span><span class=\"p\">,</span>
    <span class=\"n\">config</span><span class=\"o\">=</span><span class=\"kc\">None</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">__version__</span> <span class=\"o\">=</span> <span class=\"n\">__version__</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">stages_ran</span> <span class=\"o\">=</span> <span class=\"nb\">set</span><span
    class=\"p\">()</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">threded</span> <span class=\"o\">=</span> <span class=\"kc\">False</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_cache</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_precache</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">MARKATA_CACHE_DIR</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;.&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">/</span> <span class=\"s2\">&quot;.markata.cache&quot;</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">MARKATA_CACHE_DIR</span><span class=\"o\">.</span><span class=\"n\">mkdir</span><span
    class=\"p\">(</span><span class=\"n\">exist_ok</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_pm</span> <span class=\"o\">=</span> <span
    class=\"n\">pluggy</span><span class=\"o\">.</span><span class=\"n\">PluginManager</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;markata&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_pm</span><span class=\"o\">.</span><span class=\"n\">add_hookspecs</span><span
    class=\"p\">(</span><span class=\"n\">hookspec</span><span class=\"o\">.</span><span
    class=\"n\">MarkataSpecs</span><span class=\"p\">)</span>\n                <span
    class=\"k\">if</span> <span class=\"n\">config</span> <span class=\"ow\">is</span>
    <span class=\"ow\">not</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">config</span> <span class=\"o\">=</span> <span class=\"n\">config</span>\n
    \               <span class=\"k\">with</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">cache</span> <span class=\"k\">as</span>
    <span class=\"n\">cache</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">init_cache_stats</span>
    <span class=\"o\">=</span> <span class=\"n\">cache</span><span class=\"o\">.</span><span
    class=\"n\">stats</span><span class=\"p\">()</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">registered_attrs</span> <span class=\"o\">=</span>
    <span class=\"n\">hookspec</span><span class=\"o\">.</span><span class=\"n\">registered_attrs</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">post_models</span> <span class=\"o\">=</span> <span class=\"p\">[]</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">config_models</span> <span class=\"o\">=</span> <span class=\"p\">[]</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">config</span> <span
    class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">raw_hooks</span> <span
    class=\"o\">=</span> <span class=\"n\">config</span>\n                <span class=\"k\">else</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">raw_hooks</span> <span
    class=\"o\">=</span> <span class=\"n\">standard_config</span><span class=\"o\">.</span><span
    class=\"n\">load</span><span class=\"p\">(</span><span class=\"s2\">&quot;markata&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">hooks_conf</span> <span class=\"o\">=</span> <span class=\"n\">HooksConfig</span><span
    class=\"o\">.</span><span class=\"n\">parse_obj</span><span class=\"p\">(</span><span
    class=\"n\">raw_hooks</span><span class=\"p\">)</span>\n                <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                    <span class=\"n\">default_index</span>
    <span class=\"o\">=</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">hooks_conf</span><span class=\"o\">.</span><span class=\"n\">hooks</span><span
    class=\"o\">.</span><span class=\"n\">index</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;default&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"n\">hooks</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
    \                       <span class=\"o\">*</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">hooks_conf</span><span class=\"o\">.</span><span
    class=\"n\">hooks</span><span class=\"p\">[:</span><span class=\"n\">default_index</span><span
    class=\"p\">],</span>\n                        <span class=\"o\">*</span><span
    class=\"n\">DEFAULT_HOOKS</span><span class=\"p\">,</span>\n                        <span
    class=\"o\">*</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">hooks_conf</span><span class=\"o\">.</span><span class=\"n\">hooks</span><span
    class=\"p\">[</span><span class=\"n\">default_index</span> <span class=\"o\">+</span>
    <span class=\"mi\">1</span> <span class=\"p\">:],</span>\n                    <span
    class=\"p\">]</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">hooks_conf</span><span class=\"o\">.</span><span
    class=\"n\">hooks</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
    \                       <span class=\"n\">hook</span> <span class=\"k\">for</span>
    <span class=\"n\">hook</span> <span class=\"ow\">in</span> <span class=\"n\">hooks</span>
    <span class=\"k\">if</span> <span class=\"n\">hook</span> <span class=\"ow\">not</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">hooks_conf</span><span class=\"o\">.</span><span class=\"n\">disabled_hooks</span>\n
    \                   <span class=\"p\">]</span>\n                <span class=\"k\">except</span>
    <span class=\"ne\">ValueError</span><span class=\"p\">:</span>\n                    <span
    class=\"c1\"># &#39;default&#39; is not in hooks , do not replace with default_hooks</span>\n
    \                   <span class=\"k\">pass</span>\n\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_register_hooks</span><span class=\"p\">()</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">console</span> <span
    class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_console</span> <span class=\"o\">=</span>
    <span class=\"n\">console</span>\n                <span class=\"n\">atexit</span><span
    class=\"o\">.</span><span class=\"n\">register</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">teardown</span><span
    class=\"p\">)</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">precache</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='cache'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>cache <em class='small'>method</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">cache
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
    <span class=\"nf\">cache</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">Cache</span><span class=\"p\">:</span>\n
    \               <span class=\"c1\"># if self.threded:</span>\n                <span
    class=\"c1\">#     FanoutCache(self.MARKATA_CACHE_DIR, statistics=True)</span>\n
    \               <span class=\"k\">if</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_cache</span> <span class=\"ow\">is</span>
    <span class=\"ow\">not</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_cache</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_cache</span> <span class=\"o\">=</span>
    <span class=\"n\">Cache</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">MARKATA_CACHE_DIR</span><span class=\"p\">,</span>
    <span class=\"n\">statistics</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">return</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_cache</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='precache' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>precache <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">precache
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
    <span class=\"nf\">precache</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">if</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_precache</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">cache</span><span
    class=\"o\">.</span><span class=\"n\">expire</span><span class=\"p\">()</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_precache</span> <span class=\"o\">=</span> <span class=\"p\">{</span><span
    class=\"n\">k</span><span class=\"p\">:</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">cache</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"n\">k</span><span
    class=\"p\">)</span> <span class=\"k\">for</span> <span class=\"n\">k</span> <span
    class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">cache</span><span class=\"o\">.</span><span class=\"n\">iterkeys</span><span
    class=\"p\">()}</span>\n                <span class=\"k\">return</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_precache</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='__getattr__' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'><strong>getattr</strong> <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\"><strong>getattr</strong>
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
    <span class=\"fm\">__getattr__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">item</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span
    class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"n\">item</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_pm</span><span class=\"o\">.</span><span class=\"n\">hook</span><span
    class=\"o\">.</span><span class=\"vm\">__dict__</span><span class=\"p\">:</span>\n
    \                   <span class=\"c1\"># item is a hook, return a callable function</span>\n
    \                   <span class=\"k\">return</span> <span class=\"k\">lambda</span><span
    class=\"p\">:</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">run</span><span class=\"p\">(</span><span class=\"n\">item</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">if</span> <span class=\"n\">item</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"vm\">__dict__</span><span class=\"p\">:</span>\n                    <span
    class=\"c1\"># item is an attribute, return it</span>\n                    <span
    class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"fm\">__getitem__</span><span class=\"p\">(</span><span class=\"n\">item</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">elif</span> <span class=\"n\">item</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">registered_attrs</span><span class=\"p\">:</span>\n                    <span
    class=\"c1\"># item is created by a plugin, run it</span>\n                    <span
    class=\"n\">stage_to_run_to</span> <span class=\"o\">=</span> <span class=\"nb\">max</span><span
    class=\"p\">(</span>\n                        <span class=\"p\">[</span><span
    class=\"n\">attr</span><span class=\"p\">[</span><span class=\"s2\">&quot;lifecycle&quot;</span><span
    class=\"p\">]</span> <span class=\"k\">for</span> <span class=\"n\">attr</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">registered_attrs</span><span class=\"p\">[</span><span class=\"n\">item</span><span
    class=\"p\">]],</span>\n                    <span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">name</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">log</span><span class=\"p\">(</span>\n                        <span
    class=\"sa\">f</span><span class=\"s2\">&quot;Running to [purple]</span><span
    class=\"si\">{</span><span class=\"n\">stage_to_run_to</span><span class=\"si\">}</span><span
    class=\"s2\">[/] to retrieve [purple]</span><span class=\"si\">{</span><span class=\"n\">item</span><span
    class=\"si\">}</span><span class=\"s2\">[/]&quot;</span>\n                    <span
    class=\"p\">)</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">run</span><span class=\"p\">(</span><span
    class=\"n\">stage_to_run_to</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">return</span> <span class=\"nb\">getattr</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"n\">item</span><span
    class=\"p\">)</span>\n                <span class=\"k\">elif</span> <span class=\"n\">item</span>
    <span class=\"o\">==</span> <span class=\"s2\">&quot;precache&quot;</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_precache</span>
    <span class=\"ow\">or</span> <span class=\"p\">{}</span>\n                <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                    <span class=\"c1\">#
    Markata does not know what this is, raise</span>\n                    <span class=\"k\">raise</span>
    <span class=\"ne\">AttributeError</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;&#39;Markata&#39; object has no attribute &#39;</span><span
    class=\"si\">{</span><span class=\"n\">item</span><span class=\"si\">}</span><span
    class=\"s2\">&#39;&quot;</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">Table</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">grid</span> <span class=\"o\">=</span> <span
    class=\"n\">Table</span><span class=\"o\">.</span><span class=\"n\">grid</span><span
    class=\"p\">()</span>\n                <span class=\"n\">grid</span><span class=\"o\">.</span><span
    class=\"n\">add_column</span><span class=\"p\">(</span><span class=\"s2\">&quot;label&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"n\">grid</span><span class=\"o\">.</span><span
    class=\"n\">add_column</span><span class=\"p\">(</span><span class=\"s2\">&quot;value&quot;</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">for</span> <span class=\"n\">label</span><span
    class=\"p\">,</span> <span class=\"n\">value</span> <span class=\"ow\">in</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">describe</span><span
    class=\"p\">()</span><span class=\"o\">.</span><span class=\"n\">items</span><span
    class=\"p\">():</span>\n                    <span class=\"n\">grid</span><span
    class=\"o\">.</span><span class=\"n\">add_row</span><span class=\"p\">(</span><span
    class=\"n\">label</span><span class=\"p\">,</span> <span class=\"n\">value</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">return</span> <span
    class=\"n\">grid</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='bust_cache'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>bust_cache <em class='small'>method</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">bust_cache
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
    <span class=\"nf\">bust_cache</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">Markata</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">with</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">cache</span> <span class=\"k\">as</span>
    <span class=\"n\">cache</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">cache</span><span class=\"o\">.</span><span class=\"n\">clear</span><span
    class=\"p\">()</span>\n                <span class=\"k\">return</span> <span class=\"bp\">self</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='configure' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>configure <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">configure
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
    <span class=\"nf\">configure</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Markata</span><span
    class=\"p\">:</span>\n                <span class=\"n\">sys</span><span class=\"o\">.</span><span
    class=\"n\">path</span><span class=\"o\">.</span><span class=\"n\">append</span><span
    class=\"p\">(</span><span class=\"n\">os</span><span class=\"o\">.</span><span
    class=\"n\">getcwd</span><span class=\"p\">())</span>\n                <span class=\"c1\">#
    self.config = {**DEFUALT_CONFIG, **standard_config.load(&quot;markata&quot;)}</span>\n
    \               <span class=\"c1\"># if isinstance(self.config[&quot;glob_patterns&quot;],
    str):</span>\n                <span class=\"c1\">#     self.config[&quot;glob_patterns&quot;]
    = self.config[&quot;glob_patterns&quot;].split(&quot;,&quot;)</span>\n                <span
    class=\"c1\"># elif isinstance(self.config[&quot;glob_patterns&quot;], list):</span>\n
    \               <span class=\"c1\">#     self.config[&quot;glob_patterns&quot;]
    = list(self.config[&quot;glob_patterns&quot;])</span>\n                <span class=\"c1\">#
    else:</span>\n                <span class=\"c1\">#     raise TypeError(&quot;glob_patterns
    must be list or str&quot;)</span>\n                <span class=\"c1\"># self.glob_patterns
    = self.config[&quot;glob_patterns&quot;]</span>\n\n                <span class=\"c1\">#
    self.hooks = self.config[&quot;hooks&quot;]</span>\n\n                <span class=\"c1\">#
    if &quot;disabled_hooks&quot; not in self.config:</span>\n                <span
    class=\"c1\">#     self.disabled_hooks = [&quot;&quot;]</span>\n                <span
    class=\"c1\"># if isinstance(self.config[&quot;disabled_hooks&quot;], str):</span>\n
    \               <span class=\"c1\">#     self.disabled_hooks = self.config[&quot;disabled_hooks&quot;].split(&quot;,&quot;)</span>\n
    \               <span class=\"c1\"># if isinstance(self.config[&quot;disabled_hooks&quot;],
    list):</span>\n                <span class=\"c1\">#     self.disabled_hooks =
    self.config[&quot;disabled_hooks&quot;]</span>\n\n                <span class=\"c1\">#
    if not self.config.get(&quot;output_dir&quot;, &quot;markout&quot;).endswith(</span>\n
    \               <span class=\"c1\">#     self.config.get(&quot;path_prefix&quot;,
    &quot;&quot;)</span>\n                <span class=\"c1\"># ):</span>\n                <span
    class=\"c1\">#     self.config[&quot;output_dir&quot;] = (</span>\n                <span
    class=\"c1\">#         self.config.get(&quot;output_dir&quot;, &quot;markout&quot;)
    +</span>\n                <span class=\"c1\">#         &quot;/&quot; +</span>\n
    \               <span class=\"c1\">#         self.config.get(&quot;path_prefix&quot;,
    &quot;&quot;).rstrip(&quot;/&quot;)</span>\n                <span class=\"c1\">#
    \    )</span>\n                <span class=\"c1\"># if (</span>\n                <span
    class=\"c1\">#     len((output_split := self.config.get(&quot;output_dir&quot;,
    &quot;markout&quot;).split(&quot;/&quot;))) &gt;</span>\n                <span
    class=\"c1\">#     1</span>\n                <span class=\"c1\"># ):</span>\n
    \               <span class=\"c1\">#     if &quot;path_prefix&quot; not in self.config.keys():</span>\n
    \               <span class=\"c1\">#         self.config[&quot;path_prefix&quot;]
    = &quot;/&quot;.join(output_split[1:]) + &quot;/&quot;</span>\n                <span
    class=\"c1\"># if not self.config.get(&quot;path_prefix&quot;, &quot;&quot;).endswith(&quot;/&quot;):</span>\n
    \               <span class=\"c1\">#     self.config[&quot;path_prefix&quot;]
    = self.config.get(&quot;path_prefix&quot;, &quot;&quot;) + &quot;/&quot;</span>\n\n
    \               <span class=\"c1\"># self.config[&quot;output_dir&quot;] = self.config[&quot;output_dir&quot;].lstrip(&quot;/&quot;)</span>\n
    \               <span class=\"c1\"># self.config[&quot;path_prefix&quot;] = self.config[&quot;path_prefix&quot;].lstrip(&quot;/&quot;)</span>\n\n
    \               <span class=\"k\">try</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">default_index</span> <span class=\"o\">=</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">hooks_conf</span><span class=\"o\">.</span><span
    class=\"n\">hooks</span><span class=\"o\">.</span><span class=\"n\">index</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;default&quot;</span><span class=\"p\">)</span>\n
    \                   <span class=\"n\">hooks</span> <span class=\"o\">=</span>
    <span class=\"p\">[</span>\n                        <span class=\"o\">*</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">hooks_conf</span><span
    class=\"o\">.</span><span class=\"n\">hooks</span><span class=\"p\">[:</span><span
    class=\"n\">default_index</span><span class=\"p\">],</span>\n                        <span
    class=\"o\">*</span><span class=\"n\">DEFAULT_HOOKS</span><span class=\"p\">,</span>\n
    \                       <span class=\"o\">*</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">hooks_conf</span><span class=\"o\">.</span><span
    class=\"n\">hooks</span><span class=\"p\">[</span><span class=\"n\">default_index</span>
    <span class=\"o\">+</span> <span class=\"mi\">1</span> <span class=\"p\">:],</span>\n
    \                   <span class=\"p\">]</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">hooks</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
    \                       <span class=\"n\">hook</span> <span class=\"k\">for</span>
    <span class=\"n\">hook</span> <span class=\"ow\">in</span> <span class=\"n\">hooks</span>
    <span class=\"k\">if</span> <span class=\"n\">hook</span> <span class=\"ow\">not</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">disabled_hooks</span>\n
    \                   <span class=\"p\">]</span>\n                <span class=\"k\">except</span>
    <span class=\"ne\">ValueError</span><span class=\"p\">:</span>\n                    <span
    class=\"c1\"># &#39;default&#39; is not in hooks , do not replace with default_hooks</span>\n
    \                   <span class=\"k\">pass</span>\n\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_pm</span> <span class=\"o\">=</span> <span
    class=\"n\">pluggy</span><span class=\"o\">.</span><span class=\"n\">PluginManager</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;markata&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_pm</span><span class=\"o\">.</span><span class=\"n\">add_hookspecs</span><span
    class=\"p\">(</span><span class=\"n\">hookspec</span><span class=\"o\">.</span><span
    class=\"n\">MarkataSpecs</span><span class=\"p\">)</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_register_hooks</span><span
    class=\"p\">()</span>\n\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_pm</span><span class=\"o\">.</span><span class=\"n\">hook</span><span
    class=\"o\">.</span><span class=\"n\">configure</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"o\">=</span><span class=\"bp\">self</span><span
    class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"bp\">self</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='get_plugin_config' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>get_plugin_config <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">get_plugin_config
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
    <span class=\"nf\">get_plugin_config</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">path_or_name</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Dict</span><span class=\"p\">:</span>\n                <span
    class=\"n\">key</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">path_or_name</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">stem</span>\n\n                <span class=\"n\">config</span>
    <span class=\"o\">=</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"n\">key</span><span class=\"p\">,</span> <span
    class=\"p\">{})</span>\n\n                <span class=\"k\">if</span> <span class=\"ow\">not</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">config</span><span
    class=\"p\">,</span> <span class=\"nb\">dict</span><span class=\"p\">):</span>\n
    \                   <span class=\"k\">raise</span> <span class=\"ne\">TypeError</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;must use dict&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">if</span> <span class=\"s2\">&quot;cache_expire&quot;</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">keys</span><span class=\"p\">():</span>\n
    \                   <span class=\"n\">config</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;cache_expire&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;default_cache_expire&quot;</span><span
    class=\"p\">]</span>\n                <span class=\"k\">if</span> <span class=\"s2\">&quot;config_key&quot;</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">keys</span><span class=\"p\">():</span>\n
    \                   <span class=\"n\">config</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;config_key&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"n\">key</span>\n                <span class=\"k\">return</span>
    <span class=\"n\">config</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='get_config'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get_config <em class='small'>method</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">get_config
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
    <span class=\"nf\">get_config</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">self</span><span class=\"p\">,</span>\n                <span class=\"n\">key</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">default</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">warn</span><span class=\"p\">:</span>
    <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"kc\">True</span><span
    class=\"p\">,</span>\n                <span class=\"n\">suggested</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"kc\">None</span><span class=\"p\">,</span>\n            <span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">key</span> <span
    class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">keys</span><span
    class=\"p\">():</span>\n                    <span class=\"k\">return</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"p\">[</span><span class=\"n\">key</span><span class=\"p\">]</span>\n                <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                    <span class=\"k\">if</span>
    <span class=\"n\">suggested</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                        <span class=\"n\">suggested</span>
    <span class=\"o\">=</span> <span class=\"n\">textwrap</span><span class=\"o\">.</span><span
    class=\"n\">dedent</span><span class=\"p\">(</span>\n                            <span
    class=\"sa\">f</span><span class=\"s2\">&quot;&quot;&quot;</span>\n<span class=\"s2\">
    \                           [markata]</span>\n<span class=\"s2\">                            </span><span
    class=\"si\">{</span><span class=\"n\">key</span><span class=\"si\">}</span><span
    class=\"s2\"> = &#39;</span><span class=\"si\">{</span><span class=\"n\">default</span><span
    class=\"si\">}</span><span class=\"s2\">&#39;</span>\n<span class=\"s2\">                            &quot;&quot;&quot;</span>\n
    \                       <span class=\"p\">)</span>\n                    <span
    class=\"k\">if</span> <span class=\"n\">warn</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">logger</span><span class=\"o\">.</span><span
    class=\"n\">warning</span><span class=\"p\">(</span>\n                            <span
    class=\"n\">textwrap</span><span class=\"o\">.</span><span class=\"n\">dedent</span><span
    class=\"p\">(</span>\n                                <span class=\"sa\">f</span><span
    class=\"s2\">&quot;&quot;&quot;</span>\n<span class=\"s2\">                                Warning
    </span><span class=\"si\">{</span><span class=\"n\">key</span><span class=\"si\">}</span><span
    class=\"s2\"> is not set in markata config, sitemap will</span>\n<span class=\"s2\">
    \                               be missing root site_name</span>\n<span class=\"s2\">
    \                               to resolve this open your markata.toml and add</span>\n\n<span
    class=\"s2\">                                </span><span class=\"si\">{</span><span
    class=\"n\">suggested</span><span class=\"si\">}</span>\n\n<span class=\"s2\">
    \                               &quot;&quot;&quot;</span>\n                            <span
    class=\"p\">),</span>\n                        <span class=\"p\">)</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">default</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='make_hash' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>make_hash <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">make_hash
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
    <span class=\"nf\">make_hash</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"o\">*</span><span class=\"n\">keys</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n
    \               <span class=\"kn\">import</span> <span class=\"nn\">xxhash</span>\n\n
    \               <span class=\"n\">str_keys</span> <span class=\"o\">=</span> <span
    class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">)</span> <span class=\"k\">for</span>
    <span class=\"n\">key</span> <span class=\"ow\">in</span> <span class=\"n\">keys</span><span
    class=\"p\">]</span>\n                <span class=\"nb\">hash</span> <span class=\"o\">=</span>
    <span class=\"n\">xxhash</span><span class=\"o\">.</span><span class=\"n\">xxh64</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;&quot;</span><span class=\"o\">.</span><span
    class=\"n\">join</span><span class=\"p\">(</span><span class=\"n\">str_keys</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">encode</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;utf-8&quot;</span><span class=\"p\">))</span><span
    class=\"o\">.</span><span class=\"n\">hexdigest</span><span class=\"p\">()</span>\n
    \               <span class=\"k\">return</span> <span class=\"nb\">hash</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='content_dir_hash' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>content_dir_hash <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">content_dir_hash
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
    <span class=\"nf\">content_dir_hash</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">hashes</span> <span class=\"o\">=</span> <span
    class=\"p\">[</span>\n                    <span class=\"n\">dirhash</span><span
    class=\"p\">(</span><span class=\"nb\">dir</span><span class=\"p\">)</span>\n
    \                   <span class=\"k\">for</span> <span class=\"nb\">dir</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">content_directories</span>\n                    <span class=\"k\">if</span>
    <span class=\"nb\">dir</span><span class=\"o\">.</span><span class=\"n\">absolute</span><span
    class=\"p\">()</span> <span class=\"o\">!=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;.&quot;</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">absolute</span><span class=\"p\">()</span>\n
    \               <span class=\"p\">]</span>\n                <span class=\"k\">return</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">make_hash</span><span
    class=\"p\">(</span><span class=\"o\">*</span><span class=\"n\">hashes</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='console'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>console <em class='small'>method</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">console
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
    <span class=\"nf\">console</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">Console</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">try</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_console</span>\n                <span class=\"k\">except</span> <span
    class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_console</span>
    <span class=\"o\">=</span> <span class=\"n\">Console</span><span class=\"p\">()</span>\n
    \                   <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_console</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='describe' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>describe <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">describe
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
    <span class=\"nf\">describe</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">dict</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
    class=\"p\">]:</span>\n                <span class=\"k\">return</span> <span class=\"p\">{</span><span
    class=\"s2\">&quot;version&quot;</span><span class=\"p\">:</span> <span class=\"n\">__version__</span><span
    class=\"p\">}</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='_to_dict'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_to_dict <em class='small'>method</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_to_dict
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
    <span class=\"nf\">_to_dict</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">dict</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">,</span> <span class=\"n\">Iterable</span><span
    class=\"p\">]:</span>\n                <span class=\"k\">return</span> <span class=\"p\">{</span><span
    class=\"s2\">&quot;config&quot;</span><span class=\"p\">:</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;articles&quot;</span><span class=\"p\">:</span> <span
    class=\"p\">[</span><span class=\"n\">a</span><span class=\"o\">.</span><span
    class=\"n\">to_dict</span><span class=\"p\">()</span> <span class=\"k\">for</span>
    <span class=\"n\">a</span> <span class=\"ow\">in</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">articles</span><span class=\"p\">]}</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='to_dict' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>to_dict <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">to_dict
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
    <span class=\"nf\">to_dict</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">dict</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_to_dict</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='to_json' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>to_json <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">to_json
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
    <span class=\"nf\">to_json</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n
    \               <span class=\"kn\">import</span> <span class=\"nn\">json</span>\n\n
    \               <span class=\"k\">return</span> <span class=\"n\">json</span><span
    class=\"o\">.</span><span class=\"n\">dumps</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">to_dict</span><span
    class=\"p\">(),</span> <span class=\"n\">indent</span><span class=\"o\">=</span><span
    class=\"mi\">4</span><span class=\"p\">,</span> <span class=\"n\">sort_keys</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
    class=\"n\">default</span><span class=\"o\">=</span><span class=\"nb\">str</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='_register_hooks'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_register_hooks
    <em class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">_register_hooks <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">_register_hooks</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">sys</span><span class=\"o\">.</span><span class=\"n\">path</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">os</span><span class=\"o\">.</span><span class=\"n\">getcwd</span><span
    class=\"p\">())</span>\n                <span class=\"k\">for</span> <span class=\"n\">hook</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">hooks_conf</span><span class=\"o\">.</span><span class=\"n\">hooks</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                       <span class=\"c1\"># module style plugins</span>\n                        <span
    class=\"n\">plugin</span> <span class=\"o\">=</span> <span class=\"n\">importlib</span><span
    class=\"o\">.</span><span class=\"n\">import_module</span><span class=\"p\">(</span><span
    class=\"n\">hook</span><span class=\"p\">)</span>\n                    <span class=\"k\">except</span>
    <span class=\"ne\">ModuleNotFoundError</span> <span class=\"k\">as</span> <span
    class=\"n\">e</span><span class=\"p\">:</span>\n                        <span
    class=\"c1\"># class style plugins</span>\n                        <span class=\"k\">if</span>
    <span class=\"s2\">&quot;.&quot;</span> <span class=\"ow\">in</span> <span class=\"n\">hook</span><span
    class=\"p\">:</span>\n                            <span class=\"k\">try</span><span
    class=\"p\">:</span>\n                                <span class=\"n\">mod</span>
    <span class=\"o\">=</span> <span class=\"n\">importlib</span><span class=\"o\">.</span><span
    class=\"n\">import_module</span><span class=\"p\">(</span><span class=\"s2\">&quot;.&quot;</span><span
    class=\"o\">.</span><span class=\"n\">join</span><span class=\"p\">(</span><span
    class=\"n\">hook</span><span class=\"o\">.</span><span class=\"n\">split</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;.&quot;</span><span class=\"p\">)[:</span><span
    class=\"o\">-</span><span class=\"mi\">1</span><span class=\"p\">]))</span>\n
    \                               <span class=\"n\">plugin</span> <span class=\"o\">=</span>
    <span class=\"nb\">getattr</span><span class=\"p\">(</span><span class=\"n\">mod</span><span
    class=\"p\">,</span> <span class=\"n\">hook</span><span class=\"o\">.</span><span
    class=\"n\">split</span><span class=\"p\">(</span><span class=\"s2\">&quot;.&quot;</span><span
    class=\"p\">)[</span><span class=\"o\">-</span><span class=\"mi\">1</span><span
    class=\"p\">])</span>\n                            <span class=\"k\">except</span>
    <span class=\"ne\">ModuleNotFoundError</span> <span class=\"k\">as</span> <span
    class=\"n\">e</span><span class=\"p\">:</span>\n                                <span
    class=\"k\">raise</span> <span class=\"ne\">ModuleNotFoundError</span><span class=\"p\">(</span>\n
    \                                   <span class=\"sa\">f</span><span class=\"s2\">&quot;module
    </span><span class=\"si\">{</span><span class=\"n\">hook</span><span class=\"si\">}</span><span
    class=\"s2\"> not found</span><span class=\"se\">\\n</span><span class=\"si\">{</span><span
    class=\"n\">sys</span><span class=\"o\">.</span><span class=\"n\">path</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span>\n                                <span
    class=\"p\">)</span> <span class=\"kn\">from</span> <span class=\"nn\">e</span>\n
    \                       <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \                           <span class=\"k\">raise</span> <span class=\"n\">e</span>\n\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_pm</span><span class=\"o\">.</span><span class=\"n\">register</span><span
    class=\"p\">(</span><span class=\"n\">plugin</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='__iter__' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'><strong>iter</strong> <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\"><strong>iter</strong>
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
    <span class=\"fm\">__iter__</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">description</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;working...&quot;</span>\n
    \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Iterable</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;Markata.Post&quot;</span><span class=\"p\">]:</span>\n
    \               <span class=\"n\">articles</span><span class=\"p\">:</span> <span
    class=\"n\">Iterable</span><span class=\"p\">[</span><span class=\"n\">Markata</span><span
    class=\"o\">.</span><span class=\"n\">Post</span><span class=\"p\">]</span> <span
    class=\"o\">=</span> <span class=\"n\">track</span><span class=\"p\">(</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">articles</span><span class=\"p\">,</span>\n                    <span
    class=\"n\">description</span><span class=\"o\">=</span><span class=\"n\">description</span><span
    class=\"p\">,</span>\n                    <span class=\"n\">transient</span><span
    class=\"o\">=</span><span class=\"kc\">False</span><span class=\"p\">,</span>\n
    \                   <span class=\"n\">console</span><span class=\"o\">=</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"p\">,</span>\n                <span class=\"p\">)</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">articles</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='iter_articles' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>iter_articles <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">iter_articles
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
    <span class=\"nf\">iter_articles</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">description</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Iterable</span><span
    class=\"p\">[</span><span class=\"n\">Markata</span><span class=\"o\">.</span><span
    class=\"n\">Post</span><span class=\"p\">]:</span>\n                <span class=\"n\">articles</span><span
    class=\"p\">:</span> <span class=\"n\">Iterable</span><span class=\"p\">[</span><span
    class=\"n\">Markata</span><span class=\"o\">.</span><span class=\"n\">Post</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"n\">track</span><span
    class=\"p\">(</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">articles</span><span class=\"p\">,</span>\n
    \                   <span class=\"n\">description</span><span class=\"o\">=</span><span
    class=\"n\">description</span><span class=\"p\">,</span>\n                    <span
    class=\"n\">transient</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span>\n                    <span class=\"n\">console</span><span
    class=\"o\">=</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"p\">,</span>\n                <span class=\"p\">)</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">articles</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='teardown' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>teardown <em class='small'>method</em></h2>\ngive special access to the
    teardown lifecycle method</p>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">teardown <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">teardown</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">Markata</span><span class=\"p\">:</span>\n<span
    class=\"w\">                </span><span class=\"sd\">&quot;&quot;&quot;give special
    access to the teardown lifecycle method&quot;&quot;&quot;</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_pm</span><span
    class=\"o\">.</span><span class=\"n\">hook</span><span class=\"o\">.</span><span
    class=\"n\">teardown</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"bp\">self</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='run' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>run
    <em class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">run <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">run</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">lifecycle</span><span class=\"p\">:</span> <span class=\"n\">LifeCycle</span>
    <span class=\"o\">=</span> <span class=\"kc\">None</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">Markata</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">lifecycle</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">lifecycle</span> <span class=\"o\">=</span>
    <span class=\"nb\">max</span><span class=\"p\">(</span><span class=\"n\">LifeCycle</span><span
    class=\"o\">.</span><span class=\"n\">_member_map_</span><span class=\"o\">.</span><span
    class=\"n\">values</span><span class=\"p\">())</span>\n\n                <span
    class=\"k\">if</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">lifecycle</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
    class=\"p\">):</span>\n                    <span class=\"n\">lifecycle</span>
    <span class=\"o\">=</span> <span class=\"n\">LifeCycle</span><span class=\"p\">[</span><span
    class=\"n\">lifecycle</span><span class=\"p\">]</span>\n\n                <span
    class=\"n\">stages_to_run</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
    \                   <span class=\"n\">m</span>\n                    <span class=\"k\">for</span>
    <span class=\"n\">m</span> <span class=\"ow\">in</span> <span class=\"n\">LifeCycle</span><span
    class=\"o\">.</span><span class=\"n\">_member_map_</span>\n                    <span
    class=\"k\">if</span> <span class=\"p\">(</span><span class=\"n\">LifeCycle</span><span
    class=\"p\">[</span><span class=\"n\">m</span><span class=\"p\">]</span> <span
    class=\"o\">&lt;=</span> <span class=\"n\">lifecycle</span><span class=\"p\">)</span>
    <span class=\"ow\">and</span> <span class=\"p\">(</span><span class=\"n\">m</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">stages_ran</span><span class=\"p\">)</span>\n
    \               <span class=\"p\">]</span>\n\n                <span class=\"k\">if</span>
    <span class=\"ow\">not</span> <span class=\"n\">stages_to_run</span><span class=\"p\">:</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
    class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">lifecycle</span><span class=\"o\">.</span><span
    class=\"n\">name</span><span class=\"si\">}</span><span class=\"s2\"> already
    ran&quot;</span><span class=\"p\">)</span>\n                    <span class=\"k\">return</span>
    <span class=\"bp\">self</span>\n\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">log</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;running </span><span class=\"si\">{</span><span class=\"n\">stages_to_run</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">for</span> <span class=\"n\">stage</span> <span
    class=\"ow\">in</span> <span class=\"n\">stages_to_run</span><span class=\"p\">:</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
    class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">stage</span><span class=\"si\">}</span><span
    class=\"s2\"> running&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"nb\">getattr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_pm</span><span class=\"o\">.</span><span
    class=\"n\">hook</span><span class=\"p\">,</span> <span class=\"n\">stage</span><span
    class=\"p\">)(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"bp\">self</span><span class=\"p\">)</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">stages_ran</span><span
    class=\"o\">.</span><span class=\"n\">add</span><span class=\"p\">(</span><span
    class=\"n\">stage</span><span class=\"p\">)</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">log</span><span class=\"p\">(</span><span
    class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">stage</span><span class=\"si\">}</span><span class=\"s2\"> complete&quot;</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">with</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">cache</span> <span class=\"k\">as</span>
    <span class=\"n\">cache</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">hits</span><span class=\"p\">,</span> <span class=\"n\">misses</span>
    <span class=\"o\">=</span> <span class=\"n\">cache</span><span class=\"o\">.</span><span
    class=\"n\">stats</span><span class=\"p\">()</span>\n\n                <span class=\"k\">if</span>
    <span class=\"n\">hits</span> <span class=\"o\">+</span> <span class=\"n\">misses</span>
    <span class=\"o\">&gt;</span> <span class=\"mi\">0</span><span class=\"p\">:</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
    class=\"p\">(</span>\n                        <span class=\"sa\">f</span><span
    class=\"s2\">&quot;lifetime cache hit rate </span><span class=\"si\">{</span><span
    class=\"nb\">round</span><span class=\"p\">(</span><span class=\"n\">hits</span><span
    class=\"o\">/</span><span class=\"w\"> </span><span class=\"p\">(</span><span
    class=\"n\">hits</span><span class=\"w\"> </span><span class=\"o\">+</span><span
    class=\"w\"> </span><span class=\"n\">misses</span><span class=\"p\">)</span><span
    class=\"o\">*</span><span class=\"mi\">100</span><span class=\"p\">,</span><span
    class=\"w\"> </span><span class=\"mi\">2</span><span class=\"p\">)</span><span
    class=\"si\">}</span><span class=\"s2\">%&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"p\">)</span>\n\n                <span class=\"k\">if</span>
    <span class=\"n\">misses</span> <span class=\"o\">&gt;</span> <span class=\"mi\">0</span><span
    class=\"p\">:</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">log</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;lifetime cache hits/misses </span><span class=\"si\">{</span><span
    class=\"n\">hits</span><span class=\"si\">}</span><span class=\"s2\">/</span><span
    class=\"si\">{</span><span class=\"n\">misses</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span><span class=\"p\">)</span>\n\n                <span
    class=\"n\">hits</span> <span class=\"o\">-=</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">init_cache_stats</span><span class=\"p\">[</span><span
    class=\"mi\">0</span><span class=\"p\">]</span>\n                <span class=\"n\">misses</span>
    <span class=\"o\">-=</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">init_cache_stats</span><span class=\"p\">[</span><span class=\"mi\">1</span><span
    class=\"p\">]</span>\n\n                <span class=\"k\">if</span> <span class=\"n\">hits</span>
    <span class=\"o\">+</span> <span class=\"n\">misses</span> <span class=\"o\">&gt;</span>
    <span class=\"mi\">0</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">log</span><span class=\"p\">(</span>\n                        <span
    class=\"sa\">f</span><span class=\"s2\">&quot;run cache hit rate </span><span
    class=\"si\">{</span><span class=\"nb\">round</span><span class=\"p\">(</span><span
    class=\"n\">hits</span><span class=\"o\">/</span><span class=\"w\"> </span><span
    class=\"p\">(</span><span class=\"n\">hits</span><span class=\"w\"> </span><span
    class=\"o\">+</span><span class=\"w\"> </span><span class=\"n\">misses</span><span
    class=\"p\">)</span><span class=\"o\">*</span><span class=\"mi\">100</span><span
    class=\"p\">,</span><span class=\"w\"> </span><span class=\"mi\">2</span><span
    class=\"p\">)</span><span class=\"si\">}</span><span class=\"s2\">%&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"p\">)</span>\n\n                <span
    class=\"k\">if</span> <span class=\"n\">misses</span> <span class=\"o\">&gt;</span>
    <span class=\"mi\">0</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">log</span><span class=\"p\">(</span><span
    class=\"sa\">f</span><span class=\"s2\">&quot;run cache hits/misses </span><span
    class=\"si\">{</span><span class=\"n\">hits</span><span class=\"si\">}</span><span
    class=\"s2\">/</span><span class=\"si\">{</span><span class=\"n\">misses</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">)</span>\n\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='filter' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>filter <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">filter
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
    <span class=\"nf\">filter</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span>
    <span class=\"nb\">filter</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">list</span><span
    class=\"p\">:</span>\n                <span class=\"k\">def</span> <span class=\"nf\">evalr</span><span
    class=\"p\">(</span><span class=\"n\">a</span><span class=\"p\">:</span> <span
    class=\"n\">Markata</span><span class=\"o\">.</span><span class=\"n\">Post</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                       <span class=\"k\">return</span> <span class=\"nb\">eval</span><span
    class=\"p\">(</span>\n                            <span class=\"nb\">filter</span><span
    class=\"p\">,</span>\n                            <span class=\"p\">{</span><span
    class=\"o\">**</span><span class=\"n\">a</span><span class=\"o\">.</span><span
    class=\"n\">to_dict</span><span class=\"p\">(),</span> <span class=\"s2\">&quot;timedelta&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">timedelta</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;post&quot;</span><span class=\"p\">:</span> <span class=\"n\">a</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;m&quot;</span><span class=\"p\">:</span>
    <span class=\"bp\">self</span><span class=\"p\">},</span>\n                            <span
    class=\"p\">{},</span>\n                        <span class=\"p\">)</span>\n                    <span
    class=\"k\">except</span> <span class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n
    \                       <span class=\"k\">return</span> <span class=\"nb\">eval</span><span
    class=\"p\">(</span>\n                            <span class=\"nb\">filter</span><span
    class=\"p\">,</span>\n                            <span class=\"p\">{</span><span
    class=\"o\">**</span><span class=\"n\">a</span><span class=\"o\">.</span><span
    class=\"n\">to_dict</span><span class=\"p\">(),</span> <span class=\"s2\">&quot;timedelta&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">timedelta</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;post&quot;</span><span class=\"p\">:</span> <span class=\"n\">a</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;m&quot;</span><span class=\"p\">:</span>
    <span class=\"bp\">self</span><span class=\"p\">},</span>\n                            <span
    class=\"p\">{},</span>\n                        <span class=\"p\">)</span>\n\n
    \               <span class=\"k\">return</span> <span class=\"p\">[</span><span
    class=\"n\">a</span> <span class=\"k\">for</span> <span class=\"n\">a</span> <span
    class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">articles</span> <span class=\"k\">if</span> <span class=\"n\">evalr</span><span
    class=\"p\">(</span><span class=\"n\">a</span><span class=\"p\">)]</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='map' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>map
    <em class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">map <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">map</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">func</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;title&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"nb\">filter</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;True&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">sort</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;True&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">reverse</span><span class=\"p\">:</span>
    <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"kc\">True</span><span
    class=\"p\">,</span>\n                <span class=\"o\">*</span><span class=\"n\">args</span><span
    class=\"p\">:</span> <span class=\"nb\">tuple</span><span class=\"p\">,</span>\n
    \               <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">:</span> <span class=\"nb\">dict</span><span class=\"p\">,</span>\n
    \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">list</span><span
    class=\"p\">:</span>\n                <span class=\"kn\">import</span> <span class=\"nn\">copy</span>\n\n
    \               <span class=\"k\">def</span> <span class=\"nf\">try_sort</span><span
    class=\"p\">(</span><span class=\"n\">a</span><span class=\"p\">:</span> <span
    class=\"n\">Any</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"nb\">int</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">if</span> <span class=\"s2\">&quot;datetime&quot;</span> <span class=\"ow\">in</span>
    <span class=\"n\">sort</span><span class=\"o\">.</span><span class=\"n\">lower</span><span
    class=\"p\">():</span>\n                        <span class=\"k\">return</span>
    <span class=\"n\">a</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"n\">sort</span><span class=\"p\">,</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"p\">(</span><span class=\"mi\">1970</span><span class=\"p\">,</span> <span
    class=\"mi\">1</span><span class=\"p\">,</span> <span class=\"mi\">1</span><span
    class=\"p\">))</span>\n\n                    <span class=\"k\">if</span> <span
    class=\"s2\">&quot;date&quot;</span> <span class=\"ow\">in</span> <span class=\"n\">sort</span><span
    class=\"o\">.</span><span class=\"n\">lower</span><span class=\"p\">():</span>\n
    \                       <span class=\"k\">return</span> <span class=\"n\">a</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"n\">sort</span><span class=\"p\">,</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">date</span><span class=\"p\">(</span><span
    class=\"mi\">1970</span><span class=\"p\">,</span> <span class=\"mi\">1</span><span
    class=\"p\">,</span> <span class=\"mi\">1</span><span class=\"p\">))</span>\n\n
    \                   <span class=\"k\">try</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">value</span> <span class=\"o\">=</span> <span class=\"nb\">eval</span><span
    class=\"p\">(</span><span class=\"n\">sort</span><span class=\"p\">,</span> <span
    class=\"n\">a</span><span class=\"o\">.</span><span class=\"n\">to_dict</span><span
    class=\"p\">(),</span> <span class=\"p\">{})</span>\n                    <span
    class=\"k\">except</span> <span class=\"ne\">NameError</span><span class=\"p\">:</span>\n
    \                       <span class=\"k\">return</span> <span class=\"o\">-</span><span
    class=\"mi\">1</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">value</span>\n                    <span class=\"k\">try</span><span
    class=\"p\">:</span>\n                        <span class=\"k\">return</span>
    <span class=\"nb\">int</span><span class=\"p\">(</span><span class=\"n\">value</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">except</span> <span
    class=\"ne\">TypeError</span><span class=\"p\">:</span>\n                        <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                            <span
    class=\"k\">return</span> <span class=\"nb\">int</span><span class=\"p\">(</span><span
    class=\"n\">value</span><span class=\"o\">.</span><span class=\"n\">timestamp</span><span
    class=\"p\">())</span>\n                        <span class=\"k\">except</span>
    <span class=\"ne\">Exception</span><span class=\"p\">:</span>\n                            <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                                <span
    class=\"k\">return</span> <span class=\"nb\">int</span><span class=\"p\">(</span>\n
    \                                   <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">combine</span><span
    class=\"p\">(</span>\n                                        <span class=\"n\">value</span><span
    class=\"p\">,</span>\n                                        <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">min</span><span class=\"o\">.</span><span class=\"n\">time</span><span
    class=\"p\">(),</span>\n                                    <span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">timestamp</span><span class=\"p\">(),</span>\n
    \                               <span class=\"p\">)</span>\n                            <span
    class=\"k\">except</span> <span class=\"ne\">Exception</span><span class=\"p\">:</span>\n
    \                               <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                                   <span class=\"k\">return</span> <span class=\"nb\">sum</span><span
    class=\"p\">([</span><span class=\"nb\">ord</span><span class=\"p\">(</span><span
    class=\"n\">c</span><span class=\"p\">)</span> <span class=\"k\">for</span> <span
    class=\"n\">c</span> <span class=\"ow\">in</span> <span class=\"nb\">str</span><span
    class=\"p\">(</span><span class=\"n\">value</span><span class=\"p\">)])</span>\n
    \                               <span class=\"k\">except</span> <span class=\"ne\">Exception</span><span
    class=\"p\">:</span>\n                                    <span class=\"k\">return</span>
    <span class=\"o\">-</span><span class=\"mi\">1</span>\n\n                <span
    class=\"n\">articles</span> <span class=\"o\">=</span> <span class=\"n\">copy</span><span
    class=\"o\">.</span><span class=\"n\">copy</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">articles</span><span
    class=\"p\">)</span>\n                <span class=\"n\">articles</span><span class=\"o\">.</span><span
    class=\"n\">sort</span><span class=\"p\">(</span><span class=\"n\">key</span><span
    class=\"o\">=</span><span class=\"n\">try_sort</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">reverse</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">articles</span><span
    class=\"o\">.</span><span class=\"n\">reverse</span><span class=\"p\">()</span>\n\n
    \               <span class=\"k\">try</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">posts</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
    \                       <span class=\"nb\">eval</span><span class=\"p\">(</span>\n
    \                           <span class=\"n\">func</span><span class=\"p\">,</span>\n
    \                           <span class=\"p\">{</span><span class=\"o\">**</span><span
    class=\"n\">a</span><span class=\"o\">.</span><span class=\"n\">to_dict</span><span
    class=\"p\">(),</span> <span class=\"s2\">&quot;timedelta&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">timedelta</span><span class=\"p\">,</span> <span class=\"s2\">&quot;post&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">a</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;m&quot;</span><span class=\"p\">:</span> <span class=\"bp\">self</span><span
    class=\"p\">},</span>\n                            <span class=\"p\">{},</span>\n
    \                       <span class=\"p\">)</span>\n                        <span
    class=\"k\">for</span> <span class=\"n\">a</span> <span class=\"ow\">in</span>
    <span class=\"n\">articles</span>\n                        <span class=\"k\">if</span>
    <span class=\"nb\">eval</span><span class=\"p\">(</span>\n                            <span
    class=\"nb\">filter</span><span class=\"p\">,</span>\n                            <span
    class=\"p\">{</span><span class=\"o\">**</span><span class=\"n\">a</span><span
    class=\"o\">.</span><span class=\"n\">to_dict</span><span class=\"p\">(),</span>
    <span class=\"s2\">&quot;timedelta&quot;</span><span class=\"p\">:</span> <span
    class=\"n\">timedelta</span><span class=\"p\">,</span> <span class=\"s2\">&quot;post&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">a</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;m&quot;</span><span class=\"p\">:</span> <span class=\"bp\">self</span><span
    class=\"p\">},</span>\n                            <span class=\"p\">{},</span>\n
    \                       <span class=\"p\">)</span>\n                    <span
    class=\"p\">]</span>\n\n                <span class=\"k\">except</span> <span
    class=\"ne\">NameError</span> <span class=\"k\">as</span> <span class=\"n\">e</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">variable</span> <span
    class=\"o\">=</span> <span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">e</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">split</span><span class=\"p\">(</span><span class=\"s2\">&quot;&#39;&quot;</span><span
    class=\"p\">)[</span><span class=\"mi\">1</span><span class=\"p\">]</span>\n\n
    \                   <span class=\"n\">missing_in_posts</span> <span class=\"o\">=</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">map</span><span
    class=\"p\">(</span>\n                        <span class=\"s2\">&quot;path&quot;</span><span
    class=\"p\">,</span>\n                        <span class=\"nb\">filter</span><span
    class=\"o\">=</span><span class=\"sa\">f</span><span class=\"s1\">&#39;&quot;</span><span
    class=\"si\">{</span><span class=\"n\">variable</span><span class=\"si\">}</span><span
    class=\"s1\">&quot; not in post.keys()&#39;</span><span class=\"p\">,</span>\n
    \                   <span class=\"p\">)</span>\n                    <span class=\"n\">message</span>
    <span class=\"o\">=</span> <span class=\"p\">(</span>\n                        <span
    class=\"sa\">f</span><span class=\"s2\">&quot;variable: &#39;</span><span class=\"si\">{</span><span
    class=\"n\">variable</span><span class=\"si\">}</span><span class=\"s2\">&#39;
    is missing in </span><span class=\"si\">{</span><span class=\"nb\">len</span><span
    class=\"p\">(</span><span class=\"n\">missing_in_posts</span><span class=\"p\">)</span><span
    class=\"si\">}</span><span class=\"s2\"> posts&quot;</span>\n                    <span
    class=\"p\">)</span>\n                    <span class=\"k\">if</span> <span class=\"nb\">len</span><span
    class=\"p\">(</span><span class=\"n\">missing_in_posts</span><span class=\"p\">)</span>
    <span class=\"o\">&gt;</span> <span class=\"mi\">10</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">message</span> <span class=\"o\">+=</span>
    <span class=\"p\">(</span>\n                            <span class=\"sa\">f</span><span
    class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span class=\"s2\">first
    10 paths to posts missing </span><span class=\"si\">{</span><span class=\"n\">variable</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span>\n                            <span
    class=\"sa\">f</span><span class=\"s2\">&quot;[</span><span class=\"si\">{</span><span
    class=\"s1\">&#39;,&#39;</span><span class=\"o\">.</span><span class=\"n\">join</span><span
    class=\"p\">([</span><span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">p</span><span class=\"p\">)</span><span class=\"w\"> </span><span
    class=\"k\">for</span><span class=\"w\"> </span><span class=\"n\">p</span><span
    class=\"w\"> </span><span class=\"ow\">in</span><span class=\"w\"> </span><span
    class=\"n\">missing_in_posts</span><span class=\"p\">[:</span><span class=\"mi\">10</span><span
    class=\"p\">]])</span><span class=\"si\">}</span><span class=\"s2\">...&quot;</span>\n
    \                       <span class=\"p\">)</span>\n                    <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">message</span> <span class=\"o\">+=</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span class=\"s2\">paths
    to posts missing </span><span class=\"si\">{</span><span class=\"n\">variable</span><span
    class=\"si\">}</span><span class=\"s2\"> </span><span class=\"si\">{</span><span
    class=\"n\">missing_in_posts</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span>\n\n
    \                   <span class=\"k\">raise</span> <span class=\"n\">MissingFrontMatter</span><span
    class=\"p\">(</span><span class=\"n\">message</span><span class=\"p\">)</span>\n\n
    \               <span class=\"k\">return</span> <span class=\"n\">posts</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='first' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>first
    <em class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">first <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">first</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"nb\">filter</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;True&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">sort</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;True&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">reverse</span><span class=\"p\">:</span>
    <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"kc\">True</span><span
    class=\"p\">,</span>\n                <span class=\"o\">*</span><span class=\"n\">args</span><span
    class=\"p\">:</span> <span class=\"nb\">tuple</span><span class=\"p\">,</span>\n
    \               <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">:</span> <span class=\"nb\">dict</span><span class=\"p\">,</span>\n
    \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">list</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">map</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;post&quot;</span><span class=\"p\">,</span> <span class=\"nb\">filter</span><span
    class=\"p\">,</span> <span class=\"n\">sort</span><span class=\"p\">,</span> <span
    class=\"n\">reverse</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"n\">args</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"n\">kwargs</span><span class=\"p\">)[</span><span class=\"mi\">0</span><span
    class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='last' class='admonition-title'
    style='margin:0;padding:.5rem 1rem;'>last <em class='small'>method</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">last
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
    <span class=\"nf\">last</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"nb\">filter</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;True&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">sort</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;True&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">reverse</span><span class=\"p\">:</span>
    <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"kc\">True</span><span
    class=\"p\">,</span>\n                <span class=\"o\">*</span><span class=\"n\">args</span><span
    class=\"p\">:</span> <span class=\"nb\">tuple</span><span class=\"p\">,</span>\n
    \               <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">:</span> <span class=\"nb\">dict</span><span class=\"p\">,</span>\n
    \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">list</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">map</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;post&quot;</span><span class=\"p\">,</span> <span class=\"nb\">filter</span><span
    class=\"p\">,</span> <span class=\"n\">sort</span><span class=\"p\">,</span> <span
    class=\"n\">reverse</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"n\">args</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"n\">kwargs</span><span class=\"p\">)[</span><span class=\"o\">-</span><span
    class=\"mi\">1</span><span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='one' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>one
    <em class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">one <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">one</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"nb\">filter</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;True&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"o\">*</span><span class=\"n\">args</span><span
    class=\"p\">:</span> <span class=\"nb\">tuple</span><span class=\"p\">,</span>\n
    \               <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">:</span> <span class=\"nb\">dict</span><span class=\"p\">,</span>\n
    \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">list</span><span
    class=\"p\">:</span>\n                <span class=\"n\">posts</span> <span class=\"o\">=</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">map</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;post&quot;</span><span class=\"p\">,</span>
    <span class=\"nb\">filter</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"n\">args</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"n\">kwargs</span><span class=\"p\">)</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">len</span><span class=\"p\">(</span><span class=\"n\">posts</span><span
    class=\"p\">)</span> <span class=\"o\">&gt;</span> <span class=\"mi\">1</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">raise</span> <span
    class=\"n\">TooManyPosts</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;found </span><span class=\"si\">{</span><span class=\"nb\">len</span><span
    class=\"p\">(</span><span class=\"n\">posts</span><span class=\"p\">)</span><span
    class=\"si\">}</span><span class=\"s2\"> posts, expected 1. </span><span class=\"si\">{</span><span
    class=\"n\">posts</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"k\">if</span> <span class=\"nb\">len</span><span
    class=\"p\">(</span><span class=\"n\">posts</span><span class=\"p\">)</span> <span
    class=\"o\">==</span> <span class=\"mi\">0</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">raise</span> <span class=\"n\">NoPosts</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">posts</span><span
    class=\"p\">[</span><span class=\"mi\">0</span><span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='evalr' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>evalr <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">evalr
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
    <span class=\"nf\">evalr</span><span class=\"p\">(</span><span class=\"n\">a</span><span
    class=\"p\">:</span> <span class=\"n\">Markata</span><span class=\"o\">.</span><span
    class=\"n\">Post</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Any</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                        <span
    class=\"k\">return</span> <span class=\"nb\">eval</span><span class=\"p\">(</span>\n
    \                           <span class=\"nb\">filter</span><span class=\"p\">,</span>\n
    \                           <span class=\"p\">{</span><span class=\"o\">**</span><span
    class=\"n\">a</span><span class=\"o\">.</span><span class=\"n\">to_dict</span><span
    class=\"p\">(),</span> <span class=\"s2\">&quot;timedelta&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">timedelta</span><span class=\"p\">,</span> <span class=\"s2\">&quot;post&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">a</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;m&quot;</span><span class=\"p\">:</span> <span class=\"bp\">self</span><span
    class=\"p\">},</span>\n                            <span class=\"p\">{},</span>\n
    \                       <span class=\"p\">)</span>\n                    <span
    class=\"k\">except</span> <span class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n
    \                       <span class=\"k\">return</span> <span class=\"nb\">eval</span><span
    class=\"p\">(</span>\n                            <span class=\"nb\">filter</span><span
    class=\"p\">,</span>\n                            <span class=\"p\">{</span><span
    class=\"o\">**</span><span class=\"n\">a</span><span class=\"o\">.</span><span
    class=\"n\">to_dict</span><span class=\"p\">(),</span> <span class=\"s2\">&quot;timedelta&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">timedelta</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;post&quot;</span><span class=\"p\">:</span> <span class=\"n\">a</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;m&quot;</span><span class=\"p\">:</span>
    <span class=\"bp\">self</span><span class=\"p\">},</span>\n                            <span
    class=\"p\">{},</span>\n                        <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='try_sort' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>try_sort <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">try_sort
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
    <span class=\"nf\">try_sort</span><span class=\"p\">(</span><span class=\"n\">a</span><span
    class=\"p\">:</span> <span class=\"n\">Any</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"nb\">int</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">if</span> <span class=\"s2\">&quot;datetime&quot;</span>
    <span class=\"ow\">in</span> <span class=\"n\">sort</span><span class=\"o\">.</span><span
    class=\"n\">lower</span><span class=\"p\">():</span>\n                        <span
    class=\"k\">return</span> <span class=\"n\">a</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"n\">sort</span><span
    class=\"p\">,</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">datetime</span><span class=\"p\">(</span><span class=\"mi\">1970</span><span
    class=\"p\">,</span> <span class=\"mi\">1</span><span class=\"p\">,</span> <span
    class=\"mi\">1</span><span class=\"p\">))</span>\n\n                    <span
    class=\"k\">if</span> <span class=\"s2\">&quot;date&quot;</span> <span class=\"ow\">in</span>
    <span class=\"n\">sort</span><span class=\"o\">.</span><span class=\"n\">lower</span><span
    class=\"p\">():</span>\n                        <span class=\"k\">return</span>
    <span class=\"n\">a</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"n\">sort</span><span class=\"p\">,</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">date</span><span
    class=\"p\">(</span><span class=\"mi\">1970</span><span class=\"p\">,</span> <span
    class=\"mi\">1</span><span class=\"p\">,</span> <span class=\"mi\">1</span><span
    class=\"p\">))</span>\n\n                    <span class=\"k\">try</span><span
    class=\"p\">:</span>\n                        <span class=\"n\">value</span> <span
    class=\"o\">=</span> <span class=\"nb\">eval</span><span class=\"p\">(</span><span
    class=\"n\">sort</span><span class=\"p\">,</span> <span class=\"n\">a</span><span
    class=\"o\">.</span><span class=\"n\">to_dict</span><span class=\"p\">(),</span>
    <span class=\"p\">{})</span>\n                    <span class=\"k\">except</span>
    <span class=\"ne\">NameError</span><span class=\"p\">:</span>\n                        <span
    class=\"k\">return</span> <span class=\"o\">-</span><span class=\"mi\">1</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">value</span>\n
    \                   <span class=\"k\">try</span><span class=\"p\">:</span>\n                        <span
    class=\"k\">return</span> <span class=\"nb\">int</span><span class=\"p\">(</span><span
    class=\"n\">value</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">except</span> <span class=\"ne\">TypeError</span><span class=\"p\">:</span>\n
    \                       <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                           <span class=\"k\">return</span> <span class=\"nb\">int</span><span
    class=\"p\">(</span><span class=\"n\">value</span><span class=\"o\">.</span><span
    class=\"n\">timestamp</span><span class=\"p\">())</span>\n                        <span
    class=\"k\">except</span> <span class=\"ne\">Exception</span><span class=\"p\">:</span>\n
    \                           <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                               <span class=\"k\">return</span> <span class=\"nb\">int</span><span
    class=\"p\">(</span>\n                                    <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">combine</span><span class=\"p\">(</span>\n                                        <span
    class=\"n\">value</span><span class=\"p\">,</span>\n                                        <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">min</span><span class=\"o\">.</span><span
    class=\"n\">time</span><span class=\"p\">(),</span>\n                                    <span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">timestamp</span><span
    class=\"p\">(),</span>\n                                <span class=\"p\">)</span>\n
    \                           <span class=\"k\">except</span> <span class=\"ne\">Exception</span><span
    class=\"p\">:</span>\n                                <span class=\"k\">try</span><span
    class=\"p\">:</span>\n                                    <span class=\"k\">return</span>
    <span class=\"nb\">sum</span><span class=\"p\">([</span><span class=\"nb\">ord</span><span
    class=\"p\">(</span><span class=\"n\">c</span><span class=\"p\">)</span> <span
    class=\"k\">for</span> <span class=\"n\">c</span> <span class=\"ow\">in</span>
    <span class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">value</span><span
    class=\"p\">)])</span>\n                                <span class=\"k\">except</span>
    <span class=\"ne\">Exception</span><span class=\"p\">:</span>\n                                    <span
    class=\"k\">return</span> <span class=\"o\">-</span><span class=\"mi\">1</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>    </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>__Init__.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"Markata is a tool for handling directories of markdown.
    ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
    source  ! ???+ source \" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>__Init__.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Markata is a tool for handling directories
    of markdown. ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
    source  ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<article style=\"text-align: center;\">\n    <style>\n        section
    {\n            font-size: 200%;\n        }\n\n\n        .edit {\n            display:
    none;\n        }\n    </style>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       __Init__.Py\n    </h1>\n</section></article>\n     </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       __Init__.Py\n    </h1>\n</section>    <section class=\"body\">\n        <p>Markata
    is a tool for handling directories of markdown.</p>\n<p>!! class <h2 id='HooksConfig'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>HooksConfig <em
    class='small'>class</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">HooksConfig <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nc\">HooksConfig</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">hooks</span><span class=\"p\">:</span> <span class=\"nb\">list</span>
    <span class=\"o\">=</span> <span class=\"p\">[</span><span class=\"s2\">&quot;default&quot;</span><span
    class=\"p\">]</span>\n            <span class=\"n\">disabled_hooks</span><span
    class=\"p\">:</span> <span class=\"nb\">list</span> <span class=\"o\">=</span>
    <span class=\"p\">[]</span>\n</pre></div>\n\n</pre>\n\n<p>!! class <h2 id='Markata'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Markata <em class='small'>class</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">Markata
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
    <span class=\"nc\">Markata</span><span class=\"p\">:</span>\n            <span
    class=\"k\">def</span> <span class=\"fm\">__init__</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">console</span><span class=\"p\">:</span>
    <span class=\"n\">Console</span> <span class=\"o\">=</span> <span class=\"kc\">None</span><span
    class=\"p\">,</span> <span class=\"n\">config</span><span class=\"o\">=</span><span
    class=\"kc\">None</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">__version__</span>
    <span class=\"o\">=</span> <span class=\"n\">__version__</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">stages_ran</span>
    <span class=\"o\">=</span> <span class=\"nb\">set</span><span class=\"p\">()</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">threded</span> <span class=\"o\">=</span> <span class=\"kc\">False</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_cache</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_precache</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">MARKATA_CACHE_DIR</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;.&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">/</span> <span class=\"s2\">&quot;.markata.cache&quot;</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">MARKATA_CACHE_DIR</span><span class=\"o\">.</span><span class=\"n\">mkdir</span><span
    class=\"p\">(</span><span class=\"n\">exist_ok</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_pm</span> <span class=\"o\">=</span> <span
    class=\"n\">pluggy</span><span class=\"o\">.</span><span class=\"n\">PluginManager</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;markata&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_pm</span><span class=\"o\">.</span><span class=\"n\">add_hookspecs</span><span
    class=\"p\">(</span><span class=\"n\">hookspec</span><span class=\"o\">.</span><span
    class=\"n\">MarkataSpecs</span><span class=\"p\">)</span>\n                <span
    class=\"k\">if</span> <span class=\"n\">config</span> <span class=\"ow\">is</span>
    <span class=\"ow\">not</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">config</span> <span class=\"o\">=</span> <span class=\"n\">config</span>\n
    \               <span class=\"k\">with</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">cache</span> <span class=\"k\">as</span>
    <span class=\"n\">cache</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">init_cache_stats</span>
    <span class=\"o\">=</span> <span class=\"n\">cache</span><span class=\"o\">.</span><span
    class=\"n\">stats</span><span class=\"p\">()</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">registered_attrs</span> <span class=\"o\">=</span>
    <span class=\"n\">hookspec</span><span class=\"o\">.</span><span class=\"n\">registered_attrs</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">post_models</span> <span class=\"o\">=</span> <span class=\"p\">[]</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">config_models</span> <span class=\"o\">=</span> <span class=\"p\">[]</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">config</span> <span
    class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">raw_hooks</span> <span
    class=\"o\">=</span> <span class=\"n\">config</span>\n                <span class=\"k\">else</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">raw_hooks</span> <span
    class=\"o\">=</span> <span class=\"n\">standard_config</span><span class=\"o\">.</span><span
    class=\"n\">load</span><span class=\"p\">(</span><span class=\"s2\">&quot;markata&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">hooks_conf</span> <span class=\"o\">=</span> <span class=\"n\">HooksConfig</span><span
    class=\"o\">.</span><span class=\"n\">parse_obj</span><span class=\"p\">(</span><span
    class=\"n\">raw_hooks</span><span class=\"p\">)</span>\n                <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                    <span class=\"n\">default_index</span>
    <span class=\"o\">=</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">hooks_conf</span><span class=\"o\">.</span><span class=\"n\">hooks</span><span
    class=\"o\">.</span><span class=\"n\">index</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;default&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"n\">hooks</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
    \                       <span class=\"o\">*</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">hooks_conf</span><span class=\"o\">.</span><span
    class=\"n\">hooks</span><span class=\"p\">[:</span><span class=\"n\">default_index</span><span
    class=\"p\">],</span>\n                        <span class=\"o\">*</span><span
    class=\"n\">DEFAULT_HOOKS</span><span class=\"p\">,</span>\n                        <span
    class=\"o\">*</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">hooks_conf</span><span class=\"o\">.</span><span class=\"n\">hooks</span><span
    class=\"p\">[</span><span class=\"n\">default_index</span> <span class=\"o\">+</span>
    <span class=\"mi\">1</span> <span class=\"p\">:],</span>\n                    <span
    class=\"p\">]</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">hooks_conf</span><span class=\"o\">.</span><span
    class=\"n\">hooks</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
    \                       <span class=\"n\">hook</span> <span class=\"k\">for</span>
    <span class=\"n\">hook</span> <span class=\"ow\">in</span> <span class=\"n\">hooks</span>
    <span class=\"k\">if</span> <span class=\"n\">hook</span> <span class=\"ow\">not</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">hooks_conf</span><span class=\"o\">.</span><span class=\"n\">disabled_hooks</span>\n
    \                   <span class=\"p\">]</span>\n                <span class=\"k\">except</span>
    <span class=\"ne\">ValueError</span><span class=\"p\">:</span>\n                    <span
    class=\"c1\"># &#39;default&#39; is not in hooks , do not replace with default_hooks</span>\n
    \                   <span class=\"k\">pass</span>\n\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_register_hooks</span><span class=\"p\">()</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">console</span> <span
    class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_console</span> <span class=\"o\">=</span>
    <span class=\"n\">console</span>\n                <span class=\"n\">atexit</span><span
    class=\"o\">.</span><span class=\"n\">register</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">teardown</span><span
    class=\"p\">)</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">precache</span>\n\n            <span class=\"nd\">@property</span>\n
    \           <span class=\"k\">def</span> <span class=\"nf\">cache</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">:</span> <span
    class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Cache</span><span class=\"p\">:</span>\n                <span
    class=\"c1\"># if self.threded:</span>\n                <span class=\"c1\">#     FanoutCache(self.MARKATA_CACHE_DIR,
    statistics=True)</span>\n                <span class=\"k\">if</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_cache</span> <span class=\"ow\">is</span>
    <span class=\"ow\">not</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_cache</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_cache</span> <span class=\"o\">=</span>
    <span class=\"n\">Cache</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">MARKATA_CACHE_DIR</span><span class=\"p\">,</span>
    <span class=\"n\">statistics</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">return</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_cache</span>\n\n
    \           <span class=\"nd\">@property</span>\n            <span class=\"k\">def</span>
    <span class=\"nf\">precache</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">if</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_precache</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">cache</span><span
    class=\"o\">.</span><span class=\"n\">expire</span><span class=\"p\">()</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_precache</span> <span class=\"o\">=</span> <span class=\"p\">{</span><span
    class=\"n\">k</span><span class=\"p\">:</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">cache</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"n\">k</span><span
    class=\"p\">)</span> <span class=\"k\">for</span> <span class=\"n\">k</span> <span
    class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">cache</span><span class=\"o\">.</span><span class=\"n\">iterkeys</span><span
    class=\"p\">()}</span>\n                <span class=\"k\">return</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_precache</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"fm\">__getattr__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">:</span> <span
    class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span> <span class=\"n\">item</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">item</span> <span
    class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_pm</span><span class=\"o\">.</span><span class=\"n\">hook</span><span
    class=\"o\">.</span><span class=\"vm\">__dict__</span><span class=\"p\">:</span>\n
    \                   <span class=\"c1\"># item is a hook, return a callable function</span>\n
    \                   <span class=\"k\">return</span> <span class=\"k\">lambda</span><span
    class=\"p\">:</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">run</span><span class=\"p\">(</span><span class=\"n\">item</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">if</span> <span class=\"n\">item</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"vm\">__dict__</span><span class=\"p\">:</span>\n                    <span
    class=\"c1\"># item is an attribute, return it</span>\n                    <span
    class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"fm\">__getitem__</span><span class=\"p\">(</span><span class=\"n\">item</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">elif</span> <span class=\"n\">item</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">registered_attrs</span><span class=\"p\">:</span>\n                    <span
    class=\"c1\"># item is created by a plugin, run it</span>\n                    <span
    class=\"n\">stage_to_run_to</span> <span class=\"o\">=</span> <span class=\"nb\">max</span><span
    class=\"p\">(</span>\n                        <span class=\"p\">[</span><span
    class=\"n\">attr</span><span class=\"p\">[</span><span class=\"s2\">&quot;lifecycle&quot;</span><span
    class=\"p\">]</span> <span class=\"k\">for</span> <span class=\"n\">attr</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">registered_attrs</span><span class=\"p\">[</span><span class=\"n\">item</span><span
    class=\"p\">]],</span>\n                    <span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">name</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">log</span><span class=\"p\">(</span>\n                        <span
    class=\"sa\">f</span><span class=\"s2\">&quot;Running to [purple]</span><span
    class=\"si\">{</span><span class=\"n\">stage_to_run_to</span><span class=\"si\">}</span><span
    class=\"s2\">[/] to retrieve [purple]</span><span class=\"si\">{</span><span class=\"n\">item</span><span
    class=\"si\">}</span><span class=\"s2\">[/]&quot;</span>\n                    <span
    class=\"p\">)</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">run</span><span class=\"p\">(</span><span
    class=\"n\">stage_to_run_to</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">return</span> <span class=\"nb\">getattr</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"n\">item</span><span
    class=\"p\">)</span>\n                <span class=\"k\">elif</span> <span class=\"n\">item</span>
    <span class=\"o\">==</span> <span class=\"s2\">&quot;precache&quot;</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_precache</span>
    <span class=\"ow\">or</span> <span class=\"p\">{}</span>\n                <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                    <span class=\"c1\">#
    Markata does not know what this is, raise</span>\n                    <span class=\"k\">raise</span>
    <span class=\"ne\">AttributeError</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;&#39;Markata&#39; object has no attribute &#39;</span><span
    class=\"si\">{</span><span class=\"n\">item</span><span class=\"si\">}</span><span
    class=\"s2\">&#39;&quot;</span><span class=\"p\">)</span>\n\n            <span
    class=\"k\">def</span> <span class=\"nf\">__rich__</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Table</span><span
    class=\"p\">:</span>\n                <span class=\"n\">grid</span> <span class=\"o\">=</span>
    <span class=\"n\">Table</span><span class=\"o\">.</span><span class=\"n\">grid</span><span
    class=\"p\">()</span>\n                <span class=\"n\">grid</span><span class=\"o\">.</span><span
    class=\"n\">add_column</span><span class=\"p\">(</span><span class=\"s2\">&quot;label&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"n\">grid</span><span class=\"o\">.</span><span
    class=\"n\">add_column</span><span class=\"p\">(</span><span class=\"s2\">&quot;value&quot;</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">for</span> <span class=\"n\">label</span><span
    class=\"p\">,</span> <span class=\"n\">value</span> <span class=\"ow\">in</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">describe</span><span
    class=\"p\">()</span><span class=\"o\">.</span><span class=\"n\">items</span><span
    class=\"p\">():</span>\n                    <span class=\"n\">grid</span><span
    class=\"o\">.</span><span class=\"n\">add_row</span><span class=\"p\">(</span><span
    class=\"n\">label</span><span class=\"p\">,</span> <span class=\"n\">value</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">return</span> <span
    class=\"n\">grid</span>\n\n            <span class=\"k\">def</span> <span class=\"nf\">bust_cache</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">:</span> <span
    class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Markata</span><span class=\"p\">:</span>\n                <span
    class=\"k\">with</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">cache</span> <span class=\"k\">as</span> <span class=\"n\">cache</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">cache</span><span
    class=\"o\">.</span><span class=\"n\">clear</span><span class=\"p\">()</span>\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"nf\">configure</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Markata</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">sys</span><span class=\"o\">.</span><span class=\"n\">path</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">os</span><span class=\"o\">.</span><span class=\"n\">getcwd</span><span
    class=\"p\">())</span>\n                <span class=\"c1\"># self.config = {**DEFUALT_CONFIG,
    **standard_config.load(&quot;markata&quot;)}</span>\n                <span class=\"c1\">#
    if isinstance(self.config[&quot;glob_patterns&quot;], str):</span>\n                <span
    class=\"c1\">#     self.config[&quot;glob_patterns&quot;] = self.config[&quot;glob_patterns&quot;].split(&quot;,&quot;)</span>\n
    \               <span class=\"c1\"># elif isinstance(self.config[&quot;glob_patterns&quot;],
    list):</span>\n                <span class=\"c1\">#     self.config[&quot;glob_patterns&quot;]
    = list(self.config[&quot;glob_patterns&quot;])</span>\n                <span class=\"c1\">#
    else:</span>\n                <span class=\"c1\">#     raise TypeError(&quot;glob_patterns
    must be list or str&quot;)</span>\n                <span class=\"c1\"># self.glob_patterns
    = self.config[&quot;glob_patterns&quot;]</span>\n\n                <span class=\"c1\">#
    self.hooks = self.config[&quot;hooks&quot;]</span>\n\n                <span class=\"c1\">#
    if &quot;disabled_hooks&quot; not in self.config:</span>\n                <span
    class=\"c1\">#     self.disabled_hooks = [&quot;&quot;]</span>\n                <span
    class=\"c1\"># if isinstance(self.config[&quot;disabled_hooks&quot;], str):</span>\n
    \               <span class=\"c1\">#     self.disabled_hooks = self.config[&quot;disabled_hooks&quot;].split(&quot;,&quot;)</span>\n
    \               <span class=\"c1\"># if isinstance(self.config[&quot;disabled_hooks&quot;],
    list):</span>\n                <span class=\"c1\">#     self.disabled_hooks =
    self.config[&quot;disabled_hooks&quot;]</span>\n\n                <span class=\"c1\">#
    if not self.config.get(&quot;output_dir&quot;, &quot;markout&quot;).endswith(</span>\n
    \               <span class=\"c1\">#     self.config.get(&quot;path_prefix&quot;,
    &quot;&quot;)</span>\n                <span class=\"c1\"># ):</span>\n                <span
    class=\"c1\">#     self.config[&quot;output_dir&quot;] = (</span>\n                <span
    class=\"c1\">#         self.config.get(&quot;output_dir&quot;, &quot;markout&quot;)
    +</span>\n                <span class=\"c1\">#         &quot;/&quot; +</span>\n
    \               <span class=\"c1\">#         self.config.get(&quot;path_prefix&quot;,
    &quot;&quot;).rstrip(&quot;/&quot;)</span>\n                <span class=\"c1\">#
    \    )</span>\n                <span class=\"c1\"># if (</span>\n                <span
    class=\"c1\">#     len((output_split := self.config.get(&quot;output_dir&quot;,
    &quot;markout&quot;).split(&quot;/&quot;))) &gt;</span>\n                <span
    class=\"c1\">#     1</span>\n                <span class=\"c1\"># ):</span>\n
    \               <span class=\"c1\">#     if &quot;path_prefix&quot; not in self.config.keys():</span>\n
    \               <span class=\"c1\">#         self.config[&quot;path_prefix&quot;]
    = &quot;/&quot;.join(output_split[1:]) + &quot;/&quot;</span>\n                <span
    class=\"c1\"># if not self.config.get(&quot;path_prefix&quot;, &quot;&quot;).endswith(&quot;/&quot;):</span>\n
    \               <span class=\"c1\">#     self.config[&quot;path_prefix&quot;]
    = self.config.get(&quot;path_prefix&quot;, &quot;&quot;) + &quot;/&quot;</span>\n\n
    \               <span class=\"c1\"># self.config[&quot;output_dir&quot;] = self.config[&quot;output_dir&quot;].lstrip(&quot;/&quot;)</span>\n
    \               <span class=\"c1\"># self.config[&quot;path_prefix&quot;] = self.config[&quot;path_prefix&quot;].lstrip(&quot;/&quot;)</span>\n\n
    \               <span class=\"k\">try</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">default_index</span> <span class=\"o\">=</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">hooks_conf</span><span class=\"o\">.</span><span
    class=\"n\">hooks</span><span class=\"o\">.</span><span class=\"n\">index</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;default&quot;</span><span class=\"p\">)</span>\n
    \                   <span class=\"n\">hooks</span> <span class=\"o\">=</span>
    <span class=\"p\">[</span>\n                        <span class=\"o\">*</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">hooks_conf</span><span
    class=\"o\">.</span><span class=\"n\">hooks</span><span class=\"p\">[:</span><span
    class=\"n\">default_index</span><span class=\"p\">],</span>\n                        <span
    class=\"o\">*</span><span class=\"n\">DEFAULT_HOOKS</span><span class=\"p\">,</span>\n
    \                       <span class=\"o\">*</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">hooks_conf</span><span class=\"o\">.</span><span
    class=\"n\">hooks</span><span class=\"p\">[</span><span class=\"n\">default_index</span>
    <span class=\"o\">+</span> <span class=\"mi\">1</span> <span class=\"p\">:],</span>\n
    \                   <span class=\"p\">]</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">hooks</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
    \                       <span class=\"n\">hook</span> <span class=\"k\">for</span>
    <span class=\"n\">hook</span> <span class=\"ow\">in</span> <span class=\"n\">hooks</span>
    <span class=\"k\">if</span> <span class=\"n\">hook</span> <span class=\"ow\">not</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">disabled_hooks</span>\n
    \                   <span class=\"p\">]</span>\n                <span class=\"k\">except</span>
    <span class=\"ne\">ValueError</span><span class=\"p\">:</span>\n                    <span
    class=\"c1\"># &#39;default&#39; is not in hooks , do not replace with default_hooks</span>\n
    \                   <span class=\"k\">pass</span>\n\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_pm</span> <span class=\"o\">=</span> <span
    class=\"n\">pluggy</span><span class=\"o\">.</span><span class=\"n\">PluginManager</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;markata&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_pm</span><span class=\"o\">.</span><span class=\"n\">add_hookspecs</span><span
    class=\"p\">(</span><span class=\"n\">hookspec</span><span class=\"o\">.</span><span
    class=\"n\">MarkataSpecs</span><span class=\"p\">)</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_register_hooks</span><span
    class=\"p\">()</span>\n\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_pm</span><span class=\"o\">.</span><span class=\"n\">hook</span><span
    class=\"o\">.</span><span class=\"n\">configure</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"o\">=</span><span class=\"bp\">self</span><span
    class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"bp\">self</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"nf\">get_plugin_config</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"n\">path_or_name</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Dict</span><span
    class=\"p\">:</span>\n                <span class=\"n\">key</span> <span class=\"o\">=</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">path_or_name</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">stem</span>\n\n
    \               <span class=\"n\">config</span> <span class=\"o\">=</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">,</span> <span class=\"p\">{})</span>\n\n
    \               <span class=\"k\">if</span> <span class=\"ow\">not</span> <span
    class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">config</span><span
    class=\"p\">,</span> <span class=\"nb\">dict</span><span class=\"p\">):</span>\n
    \                   <span class=\"k\">raise</span> <span class=\"ne\">TypeError</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;must use dict&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">if</span> <span class=\"s2\">&quot;cache_expire&quot;</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">keys</span><span class=\"p\">():</span>\n
    \                   <span class=\"n\">config</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;cache_expire&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;default_cache_expire&quot;</span><span
    class=\"p\">]</span>\n                <span class=\"k\">if</span> <span class=\"s2\">&quot;config_key&quot;</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">keys</span><span class=\"p\">():</span>\n
    \                   <span class=\"n\">config</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;config_key&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"n\">key</span>\n                <span class=\"k\">return</span>
    <span class=\"n\">config</span>\n\n            <span class=\"k\">def</span> <span
    class=\"nf\">get_config</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">self</span><span class=\"p\">,</span>\n                <span class=\"n\">key</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">default</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">warn</span><span class=\"p\">:</span>
    <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"kc\">True</span><span
    class=\"p\">,</span>\n                <span class=\"n\">suggested</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"kc\">None</span><span class=\"p\">,</span>\n            <span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">key</span> <span
    class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">keys</span><span
    class=\"p\">():</span>\n                    <span class=\"k\">return</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"p\">[</span><span class=\"n\">key</span><span class=\"p\">]</span>\n                <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                    <span class=\"k\">if</span>
    <span class=\"n\">suggested</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                        <span class=\"n\">suggested</span>
    <span class=\"o\">=</span> <span class=\"n\">textwrap</span><span class=\"o\">.</span><span
    class=\"n\">dedent</span><span class=\"p\">(</span>\n                            <span
    class=\"sa\">f</span><span class=\"s2\">&quot;&quot;&quot;</span>\n<span class=\"s2\">
    \                           [markata]</span>\n<span class=\"s2\">                            </span><span
    class=\"si\">{</span><span class=\"n\">key</span><span class=\"si\">}</span><span
    class=\"s2\"> = &#39;</span><span class=\"si\">{</span><span class=\"n\">default</span><span
    class=\"si\">}</span><span class=\"s2\">&#39;</span>\n<span class=\"s2\">                            &quot;&quot;&quot;</span>\n
    \                       <span class=\"p\">)</span>\n                    <span
    class=\"k\">if</span> <span class=\"n\">warn</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">logger</span><span class=\"o\">.</span><span
    class=\"n\">warning</span><span class=\"p\">(</span>\n                            <span
    class=\"n\">textwrap</span><span class=\"o\">.</span><span class=\"n\">dedent</span><span
    class=\"p\">(</span>\n                                <span class=\"sa\">f</span><span
    class=\"s2\">&quot;&quot;&quot;</span>\n<span class=\"s2\">                                Warning
    </span><span class=\"si\">{</span><span class=\"n\">key</span><span class=\"si\">}</span><span
    class=\"s2\"> is not set in markata config, sitemap will</span>\n<span class=\"s2\">
    \                               be missing root site_name</span>\n<span class=\"s2\">
    \                               to resolve this open your markata.toml and add</span>\n\n<span
    class=\"s2\">                                </span><span class=\"si\">{</span><span
    class=\"n\">suggested</span><span class=\"si\">}</span>\n\n<span class=\"s2\">
    \                               &quot;&quot;&quot;</span>\n                            <span
    class=\"p\">),</span>\n                        <span class=\"p\">)</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">default</span>\n\n            <span
    class=\"k\">def</span> <span class=\"nf\">make_hash</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"n\">keys</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span
    class=\"p\">:</span>\n                <span class=\"kn\">import</span> <span class=\"nn\">xxhash</span>\n\n
    \               <span class=\"n\">str_keys</span> <span class=\"o\">=</span> <span
    class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">)</span> <span class=\"k\">for</span>
    <span class=\"n\">key</span> <span class=\"ow\">in</span> <span class=\"n\">keys</span><span
    class=\"p\">]</span>\n                <span class=\"nb\">hash</span> <span class=\"o\">=</span>
    <span class=\"n\">xxhash</span><span class=\"o\">.</span><span class=\"n\">xxh64</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;&quot;</span><span class=\"o\">.</span><span
    class=\"n\">join</span><span class=\"p\">(</span><span class=\"n\">str_keys</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">encode</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;utf-8&quot;</span><span class=\"p\">))</span><span
    class=\"o\">.</span><span class=\"n\">hexdigest</span><span class=\"p\">()</span>\n
    \               <span class=\"k\">return</span> <span class=\"nb\">hash</span>\n\n
    \           <span class=\"nd\">@property</span>\n            <span class=\"k\">def</span>
    <span class=\"nf\">content_dir_hash</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">hashes</span> <span class=\"o\">=</span> <span
    class=\"p\">[</span>\n                    <span class=\"n\">dirhash</span><span
    class=\"p\">(</span><span class=\"nb\">dir</span><span class=\"p\">)</span>\n
    \                   <span class=\"k\">for</span> <span class=\"nb\">dir</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">content_directories</span>\n                    <span class=\"k\">if</span>
    <span class=\"nb\">dir</span><span class=\"o\">.</span><span class=\"n\">absolute</span><span
    class=\"p\">()</span> <span class=\"o\">!=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;.&quot;</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">absolute</span><span class=\"p\">()</span>\n
    \               <span class=\"p\">]</span>\n                <span class=\"k\">return</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">make_hash</span><span
    class=\"p\">(</span><span class=\"o\">*</span><span class=\"n\">hashes</span><span
    class=\"p\">)</span>\n\n            <span class=\"nd\">@property</span>\n            <span
    class=\"k\">def</span> <span class=\"nf\">console</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Console</span><span
    class=\"p\">:</span>\n                <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_console</span>\n                <span class=\"k\">except</span>
    <span class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_console</span>
    <span class=\"o\">=</span> <span class=\"n\">Console</span><span class=\"p\">()</span>\n
    \                   <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_console</span>\n\n            <span class=\"k\">def</span>
    <span class=\"nf\">describe</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">dict</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
    class=\"p\">]:</span>\n                <span class=\"k\">return</span> <span class=\"p\">{</span><span
    class=\"s2\">&quot;version&quot;</span><span class=\"p\">:</span> <span class=\"n\">__version__</span><span
    class=\"p\">}</span>\n\n            <span class=\"k\">def</span> <span class=\"nf\">_to_dict</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">:</span> <span
    class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"nb\">dict</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"n\">Iterable</span><span class=\"p\">]:</span>\n
    \               <span class=\"k\">return</span> <span class=\"p\">{</span><span
    class=\"s2\">&quot;config&quot;</span><span class=\"p\">:</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;articles&quot;</span><span class=\"p\">:</span> <span
    class=\"p\">[</span><span class=\"n\">a</span><span class=\"o\">.</span><span
    class=\"n\">to_dict</span><span class=\"p\">()</span> <span class=\"k\">for</span>
    <span class=\"n\">a</span> <span class=\"ow\">in</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">articles</span><span class=\"p\">]}</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"nf\">to_dict</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">:</span> <span
    class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"nb\">dict</span><span class=\"p\">:</span>\n                <span
    class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_to_dict</span><span class=\"p\">()</span>\n\n            <span class=\"k\">def</span>
    <span class=\"nf\">to_json</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n
    \               <span class=\"kn\">import</span> <span class=\"nn\">json</span>\n\n
    \               <span class=\"k\">return</span> <span class=\"n\">json</span><span
    class=\"o\">.</span><span class=\"n\">dumps</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">to_dict</span><span
    class=\"p\">(),</span> <span class=\"n\">indent</span><span class=\"o\">=</span><span
    class=\"mi\">4</span><span class=\"p\">,</span> <span class=\"n\">sort_keys</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
    class=\"n\">default</span><span class=\"o\">=</span><span class=\"nb\">str</span><span
    class=\"p\">)</span>\n\n            <span class=\"k\">def</span> <span class=\"nf\">_register_hooks</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">:</span> <span
    class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
    class=\"n\">sys</span><span class=\"o\">.</span><span class=\"n\">path</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">os</span><span class=\"o\">.</span><span class=\"n\">getcwd</span><span
    class=\"p\">())</span>\n                <span class=\"k\">for</span> <span class=\"n\">hook</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">hooks_conf</span><span class=\"o\">.</span><span class=\"n\">hooks</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                       <span class=\"c1\"># module style plugins</span>\n                        <span
    class=\"n\">plugin</span> <span class=\"o\">=</span> <span class=\"n\">importlib</span><span
    class=\"o\">.</span><span class=\"n\">import_module</span><span class=\"p\">(</span><span
    class=\"n\">hook</span><span class=\"p\">)</span>\n                    <span class=\"k\">except</span>
    <span class=\"ne\">ModuleNotFoundError</span> <span class=\"k\">as</span> <span
    class=\"n\">e</span><span class=\"p\">:</span>\n                        <span
    class=\"c1\"># class style plugins</span>\n                        <span class=\"k\">if</span>
    <span class=\"s2\">&quot;.&quot;</span> <span class=\"ow\">in</span> <span class=\"n\">hook</span><span
    class=\"p\">:</span>\n                            <span class=\"k\">try</span><span
    class=\"p\">:</span>\n                                <span class=\"n\">mod</span>
    <span class=\"o\">=</span> <span class=\"n\">importlib</span><span class=\"o\">.</span><span
    class=\"n\">import_module</span><span class=\"p\">(</span><span class=\"s2\">&quot;.&quot;</span><span
    class=\"o\">.</span><span class=\"n\">join</span><span class=\"p\">(</span><span
    class=\"n\">hook</span><span class=\"o\">.</span><span class=\"n\">split</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;.&quot;</span><span class=\"p\">)[:</span><span
    class=\"o\">-</span><span class=\"mi\">1</span><span class=\"p\">]))</span>\n
    \                               <span class=\"n\">plugin</span> <span class=\"o\">=</span>
    <span class=\"nb\">getattr</span><span class=\"p\">(</span><span class=\"n\">mod</span><span
    class=\"p\">,</span> <span class=\"n\">hook</span><span class=\"o\">.</span><span
    class=\"n\">split</span><span class=\"p\">(</span><span class=\"s2\">&quot;.&quot;</span><span
    class=\"p\">)[</span><span class=\"o\">-</span><span class=\"mi\">1</span><span
    class=\"p\">])</span>\n                            <span class=\"k\">except</span>
    <span class=\"ne\">ModuleNotFoundError</span> <span class=\"k\">as</span> <span
    class=\"n\">e</span><span class=\"p\">:</span>\n                                <span
    class=\"k\">raise</span> <span class=\"ne\">ModuleNotFoundError</span><span class=\"p\">(</span>\n
    \                                   <span class=\"sa\">f</span><span class=\"s2\">&quot;module
    </span><span class=\"si\">{</span><span class=\"n\">hook</span><span class=\"si\">}</span><span
    class=\"s2\"> not found</span><span class=\"se\">\\n</span><span class=\"si\">{</span><span
    class=\"n\">sys</span><span class=\"o\">.</span><span class=\"n\">path</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span>\n                                <span
    class=\"p\">)</span> <span class=\"kn\">from</span> <span class=\"nn\">e</span>\n
    \                       <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \                           <span class=\"k\">raise</span> <span class=\"n\">e</span>\n\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_pm</span><span class=\"o\">.</span><span class=\"n\">register</span><span
    class=\"p\">(</span><span class=\"n\">plugin</span><span class=\"p\">)</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"fm\">__iter__</span><span
    class=\"p\">(</span>\n                <span class=\"bp\">self</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span> <span
    class=\"n\">description</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;working...&quot;</span>\n
    \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Iterable</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;Markata.Post&quot;</span><span class=\"p\">]:</span>\n
    \               <span class=\"n\">articles</span><span class=\"p\">:</span> <span
    class=\"n\">Iterable</span><span class=\"p\">[</span><span class=\"n\">Markata</span><span
    class=\"o\">.</span><span class=\"n\">Post</span><span class=\"p\">]</span> <span
    class=\"o\">=</span> <span class=\"n\">track</span><span class=\"p\">(</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">articles</span><span class=\"p\">,</span>\n                    <span
    class=\"n\">description</span><span class=\"o\">=</span><span class=\"n\">description</span><span
    class=\"p\">,</span>\n                    <span class=\"n\">transient</span><span
    class=\"o\">=</span><span class=\"kc\">False</span><span class=\"p\">,</span>\n
    \                   <span class=\"n\">console</span><span class=\"o\">=</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"p\">,</span>\n                <span class=\"p\">)</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">articles</span>\n\n            <span
    class=\"k\">def</span> <span class=\"nf\">iter_articles</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">description</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Iterable</span><span class=\"p\">[</span><span class=\"n\">Markata</span><span
    class=\"o\">.</span><span class=\"n\">Post</span><span class=\"p\">]:</span>\n
    \               <span class=\"n\">articles</span><span class=\"p\">:</span> <span
    class=\"n\">Iterable</span><span class=\"p\">[</span><span class=\"n\">Markata</span><span
    class=\"o\">.</span><span class=\"n\">Post</span><span class=\"p\">]</span> <span
    class=\"o\">=</span> <span class=\"n\">track</span><span class=\"p\">(</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">articles</span><span class=\"p\">,</span>\n                    <span
    class=\"n\">description</span><span class=\"o\">=</span><span class=\"n\">description</span><span
    class=\"p\">,</span>\n                    <span class=\"n\">transient</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \                   <span class=\"n\">console</span><span class=\"o\">=</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"p\">,</span>\n                <span class=\"p\">)</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">articles</span>\n\n            <span
    class=\"k\">def</span> <span class=\"nf\">teardown</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Markata</span><span
    class=\"p\">:</span>\n<span class=\"w\">                </span><span class=\"sd\">&quot;&quot;&quot;give
    special access to the teardown lifecycle method&quot;&quot;&quot;</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_pm</span><span
    class=\"o\">.</span><span class=\"n\">hook</span><span class=\"o\">.</span><span
    class=\"n\">teardown</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"bp\">self</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"nf\">run</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">lifecycle</span><span class=\"p\">:</span>
    <span class=\"n\">LifeCycle</span> <span class=\"o\">=</span> <span class=\"kc\">None</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Markata</span><span
    class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"n\">lifecycle</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">lifecycle</span> <span class=\"o\">=</span>
    <span class=\"nb\">max</span><span class=\"p\">(</span><span class=\"n\">LifeCycle</span><span
    class=\"o\">.</span><span class=\"n\">_member_map_</span><span class=\"o\">.</span><span
    class=\"n\">values</span><span class=\"p\">())</span>\n\n                <span
    class=\"k\">if</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">lifecycle</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
    class=\"p\">):</span>\n                    <span class=\"n\">lifecycle</span>
    <span class=\"o\">=</span> <span class=\"n\">LifeCycle</span><span class=\"p\">[</span><span
    class=\"n\">lifecycle</span><span class=\"p\">]</span>\n\n                <span
    class=\"n\">stages_to_run</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
    \                   <span class=\"n\">m</span>\n                    <span class=\"k\">for</span>
    <span class=\"n\">m</span> <span class=\"ow\">in</span> <span class=\"n\">LifeCycle</span><span
    class=\"o\">.</span><span class=\"n\">_member_map_</span>\n                    <span
    class=\"k\">if</span> <span class=\"p\">(</span><span class=\"n\">LifeCycle</span><span
    class=\"p\">[</span><span class=\"n\">m</span><span class=\"p\">]</span> <span
    class=\"o\">&lt;=</span> <span class=\"n\">lifecycle</span><span class=\"p\">)</span>
    <span class=\"ow\">and</span> <span class=\"p\">(</span><span class=\"n\">m</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">stages_ran</span><span class=\"p\">)</span>\n
    \               <span class=\"p\">]</span>\n\n                <span class=\"k\">if</span>
    <span class=\"ow\">not</span> <span class=\"n\">stages_to_run</span><span class=\"p\">:</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
    class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">lifecycle</span><span class=\"o\">.</span><span
    class=\"n\">name</span><span class=\"si\">}</span><span class=\"s2\"> already
    ran&quot;</span><span class=\"p\">)</span>\n                    <span class=\"k\">return</span>
    <span class=\"bp\">self</span>\n\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">log</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;running </span><span class=\"si\">{</span><span class=\"n\">stages_to_run</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">for</span> <span class=\"n\">stage</span> <span
    class=\"ow\">in</span> <span class=\"n\">stages_to_run</span><span class=\"p\">:</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
    class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">stage</span><span class=\"si\">}</span><span
    class=\"s2\"> running&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"nb\">getattr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_pm</span><span class=\"o\">.</span><span
    class=\"n\">hook</span><span class=\"p\">,</span> <span class=\"n\">stage</span><span
    class=\"p\">)(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"bp\">self</span><span class=\"p\">)</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">stages_ran</span><span
    class=\"o\">.</span><span class=\"n\">add</span><span class=\"p\">(</span><span
    class=\"n\">stage</span><span class=\"p\">)</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">log</span><span class=\"p\">(</span><span
    class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">stage</span><span class=\"si\">}</span><span class=\"s2\"> complete&quot;</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">with</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">cache</span> <span class=\"k\">as</span>
    <span class=\"n\">cache</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">hits</span><span class=\"p\">,</span> <span class=\"n\">misses</span>
    <span class=\"o\">=</span> <span class=\"n\">cache</span><span class=\"o\">.</span><span
    class=\"n\">stats</span><span class=\"p\">()</span>\n\n                <span class=\"k\">if</span>
    <span class=\"n\">hits</span> <span class=\"o\">+</span> <span class=\"n\">misses</span>
    <span class=\"o\">&gt;</span> <span class=\"mi\">0</span><span class=\"p\">:</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
    class=\"p\">(</span>\n                        <span class=\"sa\">f</span><span
    class=\"s2\">&quot;lifetime cache hit rate </span><span class=\"si\">{</span><span
    class=\"nb\">round</span><span class=\"p\">(</span><span class=\"n\">hits</span><span
    class=\"o\">/</span><span class=\"w\"> </span><span class=\"p\">(</span><span
    class=\"n\">hits</span><span class=\"w\"> </span><span class=\"o\">+</span><span
    class=\"w\"> </span><span class=\"n\">misses</span><span class=\"p\">)</span><span
    class=\"o\">*</span><span class=\"mi\">100</span><span class=\"p\">,</span><span
    class=\"w\"> </span><span class=\"mi\">2</span><span class=\"p\">)</span><span
    class=\"si\">}</span><span class=\"s2\">%&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"p\">)</span>\n\n                <span class=\"k\">if</span>
    <span class=\"n\">misses</span> <span class=\"o\">&gt;</span> <span class=\"mi\">0</span><span
    class=\"p\">:</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">log</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;lifetime cache hits/misses </span><span class=\"si\">{</span><span
    class=\"n\">hits</span><span class=\"si\">}</span><span class=\"s2\">/</span><span
    class=\"si\">{</span><span class=\"n\">misses</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span><span class=\"p\">)</span>\n\n                <span
    class=\"n\">hits</span> <span class=\"o\">-=</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">init_cache_stats</span><span class=\"p\">[</span><span
    class=\"mi\">0</span><span class=\"p\">]</span>\n                <span class=\"n\">misses</span>
    <span class=\"o\">-=</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">init_cache_stats</span><span class=\"p\">[</span><span class=\"mi\">1</span><span
    class=\"p\">]</span>\n\n                <span class=\"k\">if</span> <span class=\"n\">hits</span>
    <span class=\"o\">+</span> <span class=\"n\">misses</span> <span class=\"o\">&gt;</span>
    <span class=\"mi\">0</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">log</span><span class=\"p\">(</span>\n                        <span
    class=\"sa\">f</span><span class=\"s2\">&quot;run cache hit rate </span><span
    class=\"si\">{</span><span class=\"nb\">round</span><span class=\"p\">(</span><span
    class=\"n\">hits</span><span class=\"o\">/</span><span class=\"w\"> </span><span
    class=\"p\">(</span><span class=\"n\">hits</span><span class=\"w\"> </span><span
    class=\"o\">+</span><span class=\"w\"> </span><span class=\"n\">misses</span><span
    class=\"p\">)</span><span class=\"o\">*</span><span class=\"mi\">100</span><span
    class=\"p\">,</span><span class=\"w\"> </span><span class=\"mi\">2</span><span
    class=\"p\">)</span><span class=\"si\">}</span><span class=\"s2\">%&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"p\">)</span>\n\n                <span
    class=\"k\">if</span> <span class=\"n\">misses</span> <span class=\"o\">&gt;</span>
    <span class=\"mi\">0</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">log</span><span class=\"p\">(</span><span
    class=\"sa\">f</span><span class=\"s2\">&quot;run cache hits/misses </span><span
    class=\"si\">{</span><span class=\"n\">hits</span><span class=\"si\">}</span><span
    class=\"s2\">/</span><span class=\"si\">{</span><span class=\"n\">misses</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">)</span>\n\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"nf\">filter</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">:</span> <span
    class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span> <span class=\"nb\">filter</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"nb\">list</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">def</span> <span class=\"nf\">evalr</span><span
    class=\"p\">(</span><span class=\"n\">a</span><span class=\"p\">:</span> <span
    class=\"n\">Markata</span><span class=\"o\">.</span><span class=\"n\">Post</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                       <span class=\"k\">return</span> <span class=\"nb\">eval</span><span
    class=\"p\">(</span>\n                            <span class=\"nb\">filter</span><span
    class=\"p\">,</span>\n                            <span class=\"p\">{</span><span
    class=\"o\">**</span><span class=\"n\">a</span><span class=\"o\">.</span><span
    class=\"n\">to_dict</span><span class=\"p\">(),</span> <span class=\"s2\">&quot;timedelta&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">timedelta</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;post&quot;</span><span class=\"p\">:</span> <span class=\"n\">a</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;m&quot;</span><span class=\"p\">:</span>
    <span class=\"bp\">self</span><span class=\"p\">},</span>\n                            <span
    class=\"p\">{},</span>\n                        <span class=\"p\">)</span>\n                    <span
    class=\"k\">except</span> <span class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n
    \                       <span class=\"k\">return</span> <span class=\"nb\">eval</span><span
    class=\"p\">(</span>\n                            <span class=\"nb\">filter</span><span
    class=\"p\">,</span>\n                            <span class=\"p\">{</span><span
    class=\"o\">**</span><span class=\"n\">a</span><span class=\"o\">.</span><span
    class=\"n\">to_dict</span><span class=\"p\">(),</span> <span class=\"s2\">&quot;timedelta&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">timedelta</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;post&quot;</span><span class=\"p\">:</span> <span class=\"n\">a</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;m&quot;</span><span class=\"p\">:</span>
    <span class=\"bp\">self</span><span class=\"p\">},</span>\n                            <span
    class=\"p\">{},</span>\n                        <span class=\"p\">)</span>\n\n
    \               <span class=\"k\">return</span> <span class=\"p\">[</span><span
    class=\"n\">a</span> <span class=\"k\">for</span> <span class=\"n\">a</span> <span
    class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">articles</span> <span class=\"k\">if</span> <span class=\"n\">evalr</span><span
    class=\"p\">(</span><span class=\"n\">a</span><span class=\"p\">)]</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"nf\">map</span><span class=\"p\">(</span>\n
    \               <span class=\"bp\">self</span><span class=\"p\">:</span> <span
    class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span>\n                <span
    class=\"n\">func</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;title&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"nb\">filter</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;True&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">sort</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;True&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">reverse</span><span class=\"p\">:</span>
    <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"kc\">True</span><span
    class=\"p\">,</span>\n                <span class=\"o\">*</span><span class=\"n\">args</span><span
    class=\"p\">:</span> <span class=\"nb\">tuple</span><span class=\"p\">,</span>\n
    \               <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">:</span> <span class=\"nb\">dict</span><span class=\"p\">,</span>\n
    \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">list</span><span
    class=\"p\">:</span>\n                <span class=\"kn\">import</span> <span class=\"nn\">copy</span>\n\n
    \               <span class=\"k\">def</span> <span class=\"nf\">try_sort</span><span
    class=\"p\">(</span><span class=\"n\">a</span><span class=\"p\">:</span> <span
    class=\"n\">Any</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"nb\">int</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">if</span> <span class=\"s2\">&quot;datetime&quot;</span> <span class=\"ow\">in</span>
    <span class=\"n\">sort</span><span class=\"o\">.</span><span class=\"n\">lower</span><span
    class=\"p\">():</span>\n                        <span class=\"k\">return</span>
    <span class=\"n\">a</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"n\">sort</span><span class=\"p\">,</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"p\">(</span><span class=\"mi\">1970</span><span class=\"p\">,</span> <span
    class=\"mi\">1</span><span class=\"p\">,</span> <span class=\"mi\">1</span><span
    class=\"p\">))</span>\n\n                    <span class=\"k\">if</span> <span
    class=\"s2\">&quot;date&quot;</span> <span class=\"ow\">in</span> <span class=\"n\">sort</span><span
    class=\"o\">.</span><span class=\"n\">lower</span><span class=\"p\">():</span>\n
    \                       <span class=\"k\">return</span> <span class=\"n\">a</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"n\">sort</span><span class=\"p\">,</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">date</span><span class=\"p\">(</span><span
    class=\"mi\">1970</span><span class=\"p\">,</span> <span class=\"mi\">1</span><span
    class=\"p\">,</span> <span class=\"mi\">1</span><span class=\"p\">))</span>\n\n
    \                   <span class=\"k\">try</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">value</span> <span class=\"o\">=</span> <span class=\"nb\">eval</span><span
    class=\"p\">(</span><span class=\"n\">sort</span><span class=\"p\">,</span> <span
    class=\"n\">a</span><span class=\"o\">.</span><span class=\"n\">to_dict</span><span
    class=\"p\">(),</span> <span class=\"p\">{})</span>\n                    <span
    class=\"k\">except</span> <span class=\"ne\">NameError</span><span class=\"p\">:</span>\n
    \                       <span class=\"k\">return</span> <span class=\"o\">-</span><span
    class=\"mi\">1</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">value</span>\n                    <span class=\"k\">try</span><span
    class=\"p\">:</span>\n                        <span class=\"k\">return</span>
    <span class=\"nb\">int</span><span class=\"p\">(</span><span class=\"n\">value</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">except</span> <span
    class=\"ne\">TypeError</span><span class=\"p\">:</span>\n                        <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                            <span
    class=\"k\">return</span> <span class=\"nb\">int</span><span class=\"p\">(</span><span
    class=\"n\">value</span><span class=\"o\">.</span><span class=\"n\">timestamp</span><span
    class=\"p\">())</span>\n                        <span class=\"k\">except</span>
    <span class=\"ne\">Exception</span><span class=\"p\">:</span>\n                            <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                                <span
    class=\"k\">return</span> <span class=\"nb\">int</span><span class=\"p\">(</span>\n
    \                                   <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">combine</span><span
    class=\"p\">(</span>\n                                        <span class=\"n\">value</span><span
    class=\"p\">,</span>\n                                        <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">min</span><span class=\"o\">.</span><span class=\"n\">time</span><span
    class=\"p\">(),</span>\n                                    <span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">timestamp</span><span class=\"p\">(),</span>\n
    \                               <span class=\"p\">)</span>\n                            <span
    class=\"k\">except</span> <span class=\"ne\">Exception</span><span class=\"p\">:</span>\n
    \                               <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                                   <span class=\"k\">return</span> <span class=\"nb\">sum</span><span
    class=\"p\">([</span><span class=\"nb\">ord</span><span class=\"p\">(</span><span
    class=\"n\">c</span><span class=\"p\">)</span> <span class=\"k\">for</span> <span
    class=\"n\">c</span> <span class=\"ow\">in</span> <span class=\"nb\">str</span><span
    class=\"p\">(</span><span class=\"n\">value</span><span class=\"p\">)])</span>\n
    \                               <span class=\"k\">except</span> <span class=\"ne\">Exception</span><span
    class=\"p\">:</span>\n                                    <span class=\"k\">return</span>
    <span class=\"o\">-</span><span class=\"mi\">1</span>\n\n                <span
    class=\"n\">articles</span> <span class=\"o\">=</span> <span class=\"n\">copy</span><span
    class=\"o\">.</span><span class=\"n\">copy</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">articles</span><span
    class=\"p\">)</span>\n                <span class=\"n\">articles</span><span class=\"o\">.</span><span
    class=\"n\">sort</span><span class=\"p\">(</span><span class=\"n\">key</span><span
    class=\"o\">=</span><span class=\"n\">try_sort</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">reverse</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">articles</span><span
    class=\"o\">.</span><span class=\"n\">reverse</span><span class=\"p\">()</span>\n\n
    \               <span class=\"k\">try</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">posts</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
    \                       <span class=\"nb\">eval</span><span class=\"p\">(</span>\n
    \                           <span class=\"n\">func</span><span class=\"p\">,</span>\n
    \                           <span class=\"p\">{</span><span class=\"o\">**</span><span
    class=\"n\">a</span><span class=\"o\">.</span><span class=\"n\">to_dict</span><span
    class=\"p\">(),</span> <span class=\"s2\">&quot;timedelta&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">timedelta</span><span class=\"p\">,</span> <span class=\"s2\">&quot;post&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">a</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;m&quot;</span><span class=\"p\">:</span> <span class=\"bp\">self</span><span
    class=\"p\">},</span>\n                            <span class=\"p\">{},</span>\n
    \                       <span class=\"p\">)</span>\n                        <span
    class=\"k\">for</span> <span class=\"n\">a</span> <span class=\"ow\">in</span>
    <span class=\"n\">articles</span>\n                        <span class=\"k\">if</span>
    <span class=\"nb\">eval</span><span class=\"p\">(</span>\n                            <span
    class=\"nb\">filter</span><span class=\"p\">,</span>\n                            <span
    class=\"p\">{</span><span class=\"o\">**</span><span class=\"n\">a</span><span
    class=\"o\">.</span><span class=\"n\">to_dict</span><span class=\"p\">(),</span>
    <span class=\"s2\">&quot;timedelta&quot;</span><span class=\"p\">:</span> <span
    class=\"n\">timedelta</span><span class=\"p\">,</span> <span class=\"s2\">&quot;post&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">a</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;m&quot;</span><span class=\"p\">:</span> <span class=\"bp\">self</span><span
    class=\"p\">},</span>\n                            <span class=\"p\">{},</span>\n
    \                       <span class=\"p\">)</span>\n                    <span
    class=\"p\">]</span>\n\n                <span class=\"k\">except</span> <span
    class=\"ne\">NameError</span> <span class=\"k\">as</span> <span class=\"n\">e</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">variable</span> <span
    class=\"o\">=</span> <span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">e</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">split</span><span class=\"p\">(</span><span class=\"s2\">&quot;&#39;&quot;</span><span
    class=\"p\">)[</span><span class=\"mi\">1</span><span class=\"p\">]</span>\n\n
    \                   <span class=\"n\">missing_in_posts</span> <span class=\"o\">=</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">map</span><span
    class=\"p\">(</span>\n                        <span class=\"s2\">&quot;path&quot;</span><span
    class=\"p\">,</span>\n                        <span class=\"nb\">filter</span><span
    class=\"o\">=</span><span class=\"sa\">f</span><span class=\"s1\">&#39;&quot;</span><span
    class=\"si\">{</span><span class=\"n\">variable</span><span class=\"si\">}</span><span
    class=\"s1\">&quot; not in post.keys()&#39;</span><span class=\"p\">,</span>\n
    \                   <span class=\"p\">)</span>\n                    <span class=\"n\">message</span>
    <span class=\"o\">=</span> <span class=\"p\">(</span>\n                        <span
    class=\"sa\">f</span><span class=\"s2\">&quot;variable: &#39;</span><span class=\"si\">{</span><span
    class=\"n\">variable</span><span class=\"si\">}</span><span class=\"s2\">&#39;
    is missing in </span><span class=\"si\">{</span><span class=\"nb\">len</span><span
    class=\"p\">(</span><span class=\"n\">missing_in_posts</span><span class=\"p\">)</span><span
    class=\"si\">}</span><span class=\"s2\"> posts&quot;</span>\n                    <span
    class=\"p\">)</span>\n                    <span class=\"k\">if</span> <span class=\"nb\">len</span><span
    class=\"p\">(</span><span class=\"n\">missing_in_posts</span><span class=\"p\">)</span>
    <span class=\"o\">&gt;</span> <span class=\"mi\">10</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">message</span> <span class=\"o\">+=</span>
    <span class=\"p\">(</span>\n                            <span class=\"sa\">f</span><span
    class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span class=\"s2\">first
    10 paths to posts missing </span><span class=\"si\">{</span><span class=\"n\">variable</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span>\n                            <span
    class=\"sa\">f</span><span class=\"s2\">&quot;[</span><span class=\"si\">{</span><span
    class=\"s1\">&#39;,&#39;</span><span class=\"o\">.</span><span class=\"n\">join</span><span
    class=\"p\">([</span><span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">p</span><span class=\"p\">)</span><span class=\"w\"> </span><span
    class=\"k\">for</span><span class=\"w\"> </span><span class=\"n\">p</span><span
    class=\"w\"> </span><span class=\"ow\">in</span><span class=\"w\"> </span><span
    class=\"n\">missing_in_posts</span><span class=\"p\">[:</span><span class=\"mi\">10</span><span
    class=\"p\">]])</span><span class=\"si\">}</span><span class=\"s2\">...&quot;</span>\n
    \                       <span class=\"p\">)</span>\n                    <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">message</span> <span class=\"o\">+=</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span class=\"s2\">paths
    to posts missing </span><span class=\"si\">{</span><span class=\"n\">variable</span><span
    class=\"si\">}</span><span class=\"s2\"> </span><span class=\"si\">{</span><span
    class=\"n\">missing_in_posts</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span>\n\n
    \                   <span class=\"k\">raise</span> <span class=\"n\">MissingFrontMatter</span><span
    class=\"p\">(</span><span class=\"n\">message</span><span class=\"p\">)</span>\n\n
    \               <span class=\"k\">return</span> <span class=\"n\">posts</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"nf\">first</span><span
    class=\"p\">(</span>\n                <span class=\"bp\">self</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span>\n                <span
    class=\"nb\">filter</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;True&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">sort</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;True&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">reverse</span><span class=\"p\">:</span>
    <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"kc\">True</span><span
    class=\"p\">,</span>\n                <span class=\"o\">*</span><span class=\"n\">args</span><span
    class=\"p\">:</span> <span class=\"nb\">tuple</span><span class=\"p\">,</span>\n
    \               <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">:</span> <span class=\"nb\">dict</span><span class=\"p\">,</span>\n
    \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">list</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">map</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;post&quot;</span><span class=\"p\">,</span> <span class=\"nb\">filter</span><span
    class=\"p\">,</span> <span class=\"n\">sort</span><span class=\"p\">,</span> <span
    class=\"n\">reverse</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"n\">args</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"n\">kwargs</span><span class=\"p\">)[</span><span class=\"mi\">0</span><span
    class=\"p\">]</span>\n\n            <span class=\"k\">def</span> <span class=\"nf\">last</span><span
    class=\"p\">(</span>\n                <span class=\"bp\">self</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span>\n                <span
    class=\"nb\">filter</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;True&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">sort</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;True&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">reverse</span><span class=\"p\">:</span>
    <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"kc\">True</span><span
    class=\"p\">,</span>\n                <span class=\"o\">*</span><span class=\"n\">args</span><span
    class=\"p\">:</span> <span class=\"nb\">tuple</span><span class=\"p\">,</span>\n
    \               <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">:</span> <span class=\"nb\">dict</span><span class=\"p\">,</span>\n
    \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">list</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">map</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;post&quot;</span><span class=\"p\">,</span> <span class=\"nb\">filter</span><span
    class=\"p\">,</span> <span class=\"n\">sort</span><span class=\"p\">,</span> <span
    class=\"n\">reverse</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"n\">args</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"n\">kwargs</span><span class=\"p\">)[</span><span class=\"o\">-</span><span
    class=\"mi\">1</span><span class=\"p\">]</span>\n\n            <span class=\"k\">def</span>
    <span class=\"nf\">one</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"nb\">filter</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;True&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"o\">*</span><span class=\"n\">args</span><span
    class=\"p\">:</span> <span class=\"nb\">tuple</span><span class=\"p\">,</span>\n
    \               <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">:</span> <span class=\"nb\">dict</span><span class=\"p\">,</span>\n
    \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">list</span><span
    class=\"p\">:</span>\n                <span class=\"n\">posts</span> <span class=\"o\">=</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">map</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;post&quot;</span><span class=\"p\">,</span>
    <span class=\"nb\">filter</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"n\">args</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"n\">kwargs</span><span class=\"p\">)</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">len</span><span class=\"p\">(</span><span class=\"n\">posts</span><span
    class=\"p\">)</span> <span class=\"o\">&gt;</span> <span class=\"mi\">1</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">raise</span> <span
    class=\"n\">TooManyPosts</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;found </span><span class=\"si\">{</span><span class=\"nb\">len</span><span
    class=\"p\">(</span><span class=\"n\">posts</span><span class=\"p\">)</span><span
    class=\"si\">}</span><span class=\"s2\"> posts, expected 1. </span><span class=\"si\">{</span><span
    class=\"n\">posts</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"k\">if</span> <span class=\"nb\">len</span><span
    class=\"p\">(</span><span class=\"n\">posts</span><span class=\"p\">)</span> <span
    class=\"o\">==</span> <span class=\"mi\">0</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">raise</span> <span class=\"n\">NoPosts</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">posts</span><span
    class=\"p\">[</span><span class=\"mi\">0</span><span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='load_ipython_extension' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>load_ipython_extension <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">load_ipython_extension
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
    <span class=\"nf\">load_ipython_extension</span><span class=\"p\">(</span><span
    class=\"n\">ipython</span><span class=\"p\">):</span>\n            <span class=\"n\">ipython</span><span
    class=\"o\">.</span><span class=\"n\">user_ns</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;m&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"n\">Markata</span><span class=\"p\">()</span>\n            <span
    class=\"n\">ipython</span><span class=\"o\">.</span><span class=\"n\">user_ns</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;markata&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"n\">ipython</span><span class=\"o\">.</span><span
    class=\"n\">user_ns</span><span class=\"p\">[</span><span class=\"s2\">&quot;m&quot;</span><span
    class=\"p\">]</span>\n            <span class=\"n\">ipython</span><span class=\"o\">.</span><span
    class=\"n\">user_ns</span><span class=\"p\">[</span><span class=\"s2\">&quot;markata&quot;</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"n\">ipython</span><span
    class=\"o\">.</span><span class=\"n\">user_ns</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;m&quot;</span><span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='__init__' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'><strong>init</strong> <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\"><strong>init</strong>
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
    <span class=\"fm\">__init__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">console</span><span class=\"p\">:</span> <span class=\"n\">Console</span>
    <span class=\"o\">=</span> <span class=\"kc\">None</span><span class=\"p\">,</span>
    <span class=\"n\">config</span><span class=\"o\">=</span><span class=\"kc\">None</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">__version__</span> <span class=\"o\">=</span> <span class=\"n\">__version__</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">stages_ran</span> <span class=\"o\">=</span> <span class=\"nb\">set</span><span
    class=\"p\">()</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">threded</span> <span class=\"o\">=</span> <span class=\"kc\">False</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_cache</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_precache</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">MARKATA_CACHE_DIR</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;.&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">/</span> <span class=\"s2\">&quot;.markata.cache&quot;</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">MARKATA_CACHE_DIR</span><span class=\"o\">.</span><span class=\"n\">mkdir</span><span
    class=\"p\">(</span><span class=\"n\">exist_ok</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_pm</span> <span class=\"o\">=</span> <span
    class=\"n\">pluggy</span><span class=\"o\">.</span><span class=\"n\">PluginManager</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;markata&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_pm</span><span class=\"o\">.</span><span class=\"n\">add_hookspecs</span><span
    class=\"p\">(</span><span class=\"n\">hookspec</span><span class=\"o\">.</span><span
    class=\"n\">MarkataSpecs</span><span class=\"p\">)</span>\n                <span
    class=\"k\">if</span> <span class=\"n\">config</span> <span class=\"ow\">is</span>
    <span class=\"ow\">not</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">config</span> <span class=\"o\">=</span> <span class=\"n\">config</span>\n
    \               <span class=\"k\">with</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">cache</span> <span class=\"k\">as</span>
    <span class=\"n\">cache</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">init_cache_stats</span>
    <span class=\"o\">=</span> <span class=\"n\">cache</span><span class=\"o\">.</span><span
    class=\"n\">stats</span><span class=\"p\">()</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">registered_attrs</span> <span class=\"o\">=</span>
    <span class=\"n\">hookspec</span><span class=\"o\">.</span><span class=\"n\">registered_attrs</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">post_models</span> <span class=\"o\">=</span> <span class=\"p\">[]</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">config_models</span> <span class=\"o\">=</span> <span class=\"p\">[]</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">config</span> <span
    class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">raw_hooks</span> <span
    class=\"o\">=</span> <span class=\"n\">config</span>\n                <span class=\"k\">else</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">raw_hooks</span> <span
    class=\"o\">=</span> <span class=\"n\">standard_config</span><span class=\"o\">.</span><span
    class=\"n\">load</span><span class=\"p\">(</span><span class=\"s2\">&quot;markata&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">hooks_conf</span> <span class=\"o\">=</span> <span class=\"n\">HooksConfig</span><span
    class=\"o\">.</span><span class=\"n\">parse_obj</span><span class=\"p\">(</span><span
    class=\"n\">raw_hooks</span><span class=\"p\">)</span>\n                <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                    <span class=\"n\">default_index</span>
    <span class=\"o\">=</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">hooks_conf</span><span class=\"o\">.</span><span class=\"n\">hooks</span><span
    class=\"o\">.</span><span class=\"n\">index</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;default&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"n\">hooks</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
    \                       <span class=\"o\">*</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">hooks_conf</span><span class=\"o\">.</span><span
    class=\"n\">hooks</span><span class=\"p\">[:</span><span class=\"n\">default_index</span><span
    class=\"p\">],</span>\n                        <span class=\"o\">*</span><span
    class=\"n\">DEFAULT_HOOKS</span><span class=\"p\">,</span>\n                        <span
    class=\"o\">*</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">hooks_conf</span><span class=\"o\">.</span><span class=\"n\">hooks</span><span
    class=\"p\">[</span><span class=\"n\">default_index</span> <span class=\"o\">+</span>
    <span class=\"mi\">1</span> <span class=\"p\">:],</span>\n                    <span
    class=\"p\">]</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">hooks_conf</span><span class=\"o\">.</span><span
    class=\"n\">hooks</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
    \                       <span class=\"n\">hook</span> <span class=\"k\">for</span>
    <span class=\"n\">hook</span> <span class=\"ow\">in</span> <span class=\"n\">hooks</span>
    <span class=\"k\">if</span> <span class=\"n\">hook</span> <span class=\"ow\">not</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">hooks_conf</span><span class=\"o\">.</span><span class=\"n\">disabled_hooks</span>\n
    \                   <span class=\"p\">]</span>\n                <span class=\"k\">except</span>
    <span class=\"ne\">ValueError</span><span class=\"p\">:</span>\n                    <span
    class=\"c1\"># &#39;default&#39; is not in hooks , do not replace with default_hooks</span>\n
    \                   <span class=\"k\">pass</span>\n\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_register_hooks</span><span class=\"p\">()</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">console</span> <span
    class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_console</span> <span class=\"o\">=</span>
    <span class=\"n\">console</span>\n                <span class=\"n\">atexit</span><span
    class=\"o\">.</span><span class=\"n\">register</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">teardown</span><span
    class=\"p\">)</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">precache</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='cache'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>cache <em class='small'>method</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">cache
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
    <span class=\"nf\">cache</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">Cache</span><span class=\"p\">:</span>\n
    \               <span class=\"c1\"># if self.threded:</span>\n                <span
    class=\"c1\">#     FanoutCache(self.MARKATA_CACHE_DIR, statistics=True)</span>\n
    \               <span class=\"k\">if</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_cache</span> <span class=\"ow\">is</span>
    <span class=\"ow\">not</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_cache</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_cache</span> <span class=\"o\">=</span>
    <span class=\"n\">Cache</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">MARKATA_CACHE_DIR</span><span class=\"p\">,</span>
    <span class=\"n\">statistics</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">return</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_cache</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='precache' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>precache <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">precache
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
    <span class=\"nf\">precache</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">if</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_precache</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">cache</span><span
    class=\"o\">.</span><span class=\"n\">expire</span><span class=\"p\">()</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_precache</span> <span class=\"o\">=</span> <span class=\"p\">{</span><span
    class=\"n\">k</span><span class=\"p\">:</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">cache</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"n\">k</span><span
    class=\"p\">)</span> <span class=\"k\">for</span> <span class=\"n\">k</span> <span
    class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">cache</span><span class=\"o\">.</span><span class=\"n\">iterkeys</span><span
    class=\"p\">()}</span>\n                <span class=\"k\">return</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_precache</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='__getattr__' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'><strong>getattr</strong> <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\"><strong>getattr</strong>
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
    <span class=\"fm\">__getattr__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">item</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span
    class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"n\">item</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_pm</span><span class=\"o\">.</span><span class=\"n\">hook</span><span
    class=\"o\">.</span><span class=\"vm\">__dict__</span><span class=\"p\">:</span>\n
    \                   <span class=\"c1\"># item is a hook, return a callable function</span>\n
    \                   <span class=\"k\">return</span> <span class=\"k\">lambda</span><span
    class=\"p\">:</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">run</span><span class=\"p\">(</span><span class=\"n\">item</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">if</span> <span class=\"n\">item</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"vm\">__dict__</span><span class=\"p\">:</span>\n                    <span
    class=\"c1\"># item is an attribute, return it</span>\n                    <span
    class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"fm\">__getitem__</span><span class=\"p\">(</span><span class=\"n\">item</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">elif</span> <span class=\"n\">item</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">registered_attrs</span><span class=\"p\">:</span>\n                    <span
    class=\"c1\"># item is created by a plugin, run it</span>\n                    <span
    class=\"n\">stage_to_run_to</span> <span class=\"o\">=</span> <span class=\"nb\">max</span><span
    class=\"p\">(</span>\n                        <span class=\"p\">[</span><span
    class=\"n\">attr</span><span class=\"p\">[</span><span class=\"s2\">&quot;lifecycle&quot;</span><span
    class=\"p\">]</span> <span class=\"k\">for</span> <span class=\"n\">attr</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">registered_attrs</span><span class=\"p\">[</span><span class=\"n\">item</span><span
    class=\"p\">]],</span>\n                    <span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">name</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">log</span><span class=\"p\">(</span>\n                        <span
    class=\"sa\">f</span><span class=\"s2\">&quot;Running to [purple]</span><span
    class=\"si\">{</span><span class=\"n\">stage_to_run_to</span><span class=\"si\">}</span><span
    class=\"s2\">[/] to retrieve [purple]</span><span class=\"si\">{</span><span class=\"n\">item</span><span
    class=\"si\">}</span><span class=\"s2\">[/]&quot;</span>\n                    <span
    class=\"p\">)</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">run</span><span class=\"p\">(</span><span
    class=\"n\">stage_to_run_to</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">return</span> <span class=\"nb\">getattr</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"n\">item</span><span
    class=\"p\">)</span>\n                <span class=\"k\">elif</span> <span class=\"n\">item</span>
    <span class=\"o\">==</span> <span class=\"s2\">&quot;precache&quot;</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_precache</span>
    <span class=\"ow\">or</span> <span class=\"p\">{}</span>\n                <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                    <span class=\"c1\">#
    Markata does not know what this is, raise</span>\n                    <span class=\"k\">raise</span>
    <span class=\"ne\">AttributeError</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;&#39;Markata&#39; object has no attribute &#39;</span><span
    class=\"si\">{</span><span class=\"n\">item</span><span class=\"si\">}</span><span
    class=\"s2\">&#39;&quot;</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">Table</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">grid</span> <span class=\"o\">=</span> <span
    class=\"n\">Table</span><span class=\"o\">.</span><span class=\"n\">grid</span><span
    class=\"p\">()</span>\n                <span class=\"n\">grid</span><span class=\"o\">.</span><span
    class=\"n\">add_column</span><span class=\"p\">(</span><span class=\"s2\">&quot;label&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"n\">grid</span><span class=\"o\">.</span><span
    class=\"n\">add_column</span><span class=\"p\">(</span><span class=\"s2\">&quot;value&quot;</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">for</span> <span class=\"n\">label</span><span
    class=\"p\">,</span> <span class=\"n\">value</span> <span class=\"ow\">in</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">describe</span><span
    class=\"p\">()</span><span class=\"o\">.</span><span class=\"n\">items</span><span
    class=\"p\">():</span>\n                    <span class=\"n\">grid</span><span
    class=\"o\">.</span><span class=\"n\">add_row</span><span class=\"p\">(</span><span
    class=\"n\">label</span><span class=\"p\">,</span> <span class=\"n\">value</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">return</span> <span
    class=\"n\">grid</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='bust_cache'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>bust_cache <em class='small'>method</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">bust_cache
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
    <span class=\"nf\">bust_cache</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">Markata</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">with</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">cache</span> <span class=\"k\">as</span>
    <span class=\"n\">cache</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">cache</span><span class=\"o\">.</span><span class=\"n\">clear</span><span
    class=\"p\">()</span>\n                <span class=\"k\">return</span> <span class=\"bp\">self</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='configure' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>configure <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">configure
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
    <span class=\"nf\">configure</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Markata</span><span
    class=\"p\">:</span>\n                <span class=\"n\">sys</span><span class=\"o\">.</span><span
    class=\"n\">path</span><span class=\"o\">.</span><span class=\"n\">append</span><span
    class=\"p\">(</span><span class=\"n\">os</span><span class=\"o\">.</span><span
    class=\"n\">getcwd</span><span class=\"p\">())</span>\n                <span class=\"c1\">#
    self.config = {**DEFUALT_CONFIG, **standard_config.load(&quot;markata&quot;)}</span>\n
    \               <span class=\"c1\"># if isinstance(self.config[&quot;glob_patterns&quot;],
    str):</span>\n                <span class=\"c1\">#     self.config[&quot;glob_patterns&quot;]
    = self.config[&quot;glob_patterns&quot;].split(&quot;,&quot;)</span>\n                <span
    class=\"c1\"># elif isinstance(self.config[&quot;glob_patterns&quot;], list):</span>\n
    \               <span class=\"c1\">#     self.config[&quot;glob_patterns&quot;]
    = list(self.config[&quot;glob_patterns&quot;])</span>\n                <span class=\"c1\">#
    else:</span>\n                <span class=\"c1\">#     raise TypeError(&quot;glob_patterns
    must be list or str&quot;)</span>\n                <span class=\"c1\"># self.glob_patterns
    = self.config[&quot;glob_patterns&quot;]</span>\n\n                <span class=\"c1\">#
    self.hooks = self.config[&quot;hooks&quot;]</span>\n\n                <span class=\"c1\">#
    if &quot;disabled_hooks&quot; not in self.config:</span>\n                <span
    class=\"c1\">#     self.disabled_hooks = [&quot;&quot;]</span>\n                <span
    class=\"c1\"># if isinstance(self.config[&quot;disabled_hooks&quot;], str):</span>\n
    \               <span class=\"c1\">#     self.disabled_hooks = self.config[&quot;disabled_hooks&quot;].split(&quot;,&quot;)</span>\n
    \               <span class=\"c1\"># if isinstance(self.config[&quot;disabled_hooks&quot;],
    list):</span>\n                <span class=\"c1\">#     self.disabled_hooks =
    self.config[&quot;disabled_hooks&quot;]</span>\n\n                <span class=\"c1\">#
    if not self.config.get(&quot;output_dir&quot;, &quot;markout&quot;).endswith(</span>\n
    \               <span class=\"c1\">#     self.config.get(&quot;path_prefix&quot;,
    &quot;&quot;)</span>\n                <span class=\"c1\"># ):</span>\n                <span
    class=\"c1\">#     self.config[&quot;output_dir&quot;] = (</span>\n                <span
    class=\"c1\">#         self.config.get(&quot;output_dir&quot;, &quot;markout&quot;)
    +</span>\n                <span class=\"c1\">#         &quot;/&quot; +</span>\n
    \               <span class=\"c1\">#         self.config.get(&quot;path_prefix&quot;,
    &quot;&quot;).rstrip(&quot;/&quot;)</span>\n                <span class=\"c1\">#
    \    )</span>\n                <span class=\"c1\"># if (</span>\n                <span
    class=\"c1\">#     len((output_split := self.config.get(&quot;output_dir&quot;,
    &quot;markout&quot;).split(&quot;/&quot;))) &gt;</span>\n                <span
    class=\"c1\">#     1</span>\n                <span class=\"c1\"># ):</span>\n
    \               <span class=\"c1\">#     if &quot;path_prefix&quot; not in self.config.keys():</span>\n
    \               <span class=\"c1\">#         self.config[&quot;path_prefix&quot;]
    = &quot;/&quot;.join(output_split[1:]) + &quot;/&quot;</span>\n                <span
    class=\"c1\"># if not self.config.get(&quot;path_prefix&quot;, &quot;&quot;).endswith(&quot;/&quot;):</span>\n
    \               <span class=\"c1\">#     self.config[&quot;path_prefix&quot;]
    = self.config.get(&quot;path_prefix&quot;, &quot;&quot;) + &quot;/&quot;</span>\n\n
    \               <span class=\"c1\"># self.config[&quot;output_dir&quot;] = self.config[&quot;output_dir&quot;].lstrip(&quot;/&quot;)</span>\n
    \               <span class=\"c1\"># self.config[&quot;path_prefix&quot;] = self.config[&quot;path_prefix&quot;].lstrip(&quot;/&quot;)</span>\n\n
    \               <span class=\"k\">try</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">default_index</span> <span class=\"o\">=</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">hooks_conf</span><span class=\"o\">.</span><span
    class=\"n\">hooks</span><span class=\"o\">.</span><span class=\"n\">index</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;default&quot;</span><span class=\"p\">)</span>\n
    \                   <span class=\"n\">hooks</span> <span class=\"o\">=</span>
    <span class=\"p\">[</span>\n                        <span class=\"o\">*</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">hooks_conf</span><span
    class=\"o\">.</span><span class=\"n\">hooks</span><span class=\"p\">[:</span><span
    class=\"n\">default_index</span><span class=\"p\">],</span>\n                        <span
    class=\"o\">*</span><span class=\"n\">DEFAULT_HOOKS</span><span class=\"p\">,</span>\n
    \                       <span class=\"o\">*</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">hooks_conf</span><span class=\"o\">.</span><span
    class=\"n\">hooks</span><span class=\"p\">[</span><span class=\"n\">default_index</span>
    <span class=\"o\">+</span> <span class=\"mi\">1</span> <span class=\"p\">:],</span>\n
    \                   <span class=\"p\">]</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">hooks</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
    \                       <span class=\"n\">hook</span> <span class=\"k\">for</span>
    <span class=\"n\">hook</span> <span class=\"ow\">in</span> <span class=\"n\">hooks</span>
    <span class=\"k\">if</span> <span class=\"n\">hook</span> <span class=\"ow\">not</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">disabled_hooks</span>\n
    \                   <span class=\"p\">]</span>\n                <span class=\"k\">except</span>
    <span class=\"ne\">ValueError</span><span class=\"p\">:</span>\n                    <span
    class=\"c1\"># &#39;default&#39; is not in hooks , do not replace with default_hooks</span>\n
    \                   <span class=\"k\">pass</span>\n\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_pm</span> <span class=\"o\">=</span> <span
    class=\"n\">pluggy</span><span class=\"o\">.</span><span class=\"n\">PluginManager</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;markata&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_pm</span><span class=\"o\">.</span><span class=\"n\">add_hookspecs</span><span
    class=\"p\">(</span><span class=\"n\">hookspec</span><span class=\"o\">.</span><span
    class=\"n\">MarkataSpecs</span><span class=\"p\">)</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_register_hooks</span><span
    class=\"p\">()</span>\n\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_pm</span><span class=\"o\">.</span><span class=\"n\">hook</span><span
    class=\"o\">.</span><span class=\"n\">configure</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"o\">=</span><span class=\"bp\">self</span><span
    class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"bp\">self</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='get_plugin_config' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>get_plugin_config <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">get_plugin_config
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
    <span class=\"nf\">get_plugin_config</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">path_or_name</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Dict</span><span class=\"p\">:</span>\n                <span
    class=\"n\">key</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">path_or_name</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">stem</span>\n\n                <span class=\"n\">config</span>
    <span class=\"o\">=</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"n\">key</span><span class=\"p\">,</span> <span
    class=\"p\">{})</span>\n\n                <span class=\"k\">if</span> <span class=\"ow\">not</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">config</span><span
    class=\"p\">,</span> <span class=\"nb\">dict</span><span class=\"p\">):</span>\n
    \                   <span class=\"k\">raise</span> <span class=\"ne\">TypeError</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;must use dict&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">if</span> <span class=\"s2\">&quot;cache_expire&quot;</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">keys</span><span class=\"p\">():</span>\n
    \                   <span class=\"n\">config</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;cache_expire&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;default_cache_expire&quot;</span><span
    class=\"p\">]</span>\n                <span class=\"k\">if</span> <span class=\"s2\">&quot;config_key&quot;</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">keys</span><span class=\"p\">():</span>\n
    \                   <span class=\"n\">config</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;config_key&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"n\">key</span>\n                <span class=\"k\">return</span>
    <span class=\"n\">config</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='get_config'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get_config <em class='small'>method</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">get_config
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
    <span class=\"nf\">get_config</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">self</span><span class=\"p\">,</span>\n                <span class=\"n\">key</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">default</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">warn</span><span class=\"p\">:</span>
    <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"kc\">True</span><span
    class=\"p\">,</span>\n                <span class=\"n\">suggested</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"kc\">None</span><span class=\"p\">,</span>\n            <span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">key</span> <span
    class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">keys</span><span
    class=\"p\">():</span>\n                    <span class=\"k\">return</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"p\">[</span><span class=\"n\">key</span><span class=\"p\">]</span>\n                <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                    <span class=\"k\">if</span>
    <span class=\"n\">suggested</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                        <span class=\"n\">suggested</span>
    <span class=\"o\">=</span> <span class=\"n\">textwrap</span><span class=\"o\">.</span><span
    class=\"n\">dedent</span><span class=\"p\">(</span>\n                            <span
    class=\"sa\">f</span><span class=\"s2\">&quot;&quot;&quot;</span>\n<span class=\"s2\">
    \                           [markata]</span>\n<span class=\"s2\">                            </span><span
    class=\"si\">{</span><span class=\"n\">key</span><span class=\"si\">}</span><span
    class=\"s2\"> = &#39;</span><span class=\"si\">{</span><span class=\"n\">default</span><span
    class=\"si\">}</span><span class=\"s2\">&#39;</span>\n<span class=\"s2\">                            &quot;&quot;&quot;</span>\n
    \                       <span class=\"p\">)</span>\n                    <span
    class=\"k\">if</span> <span class=\"n\">warn</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">logger</span><span class=\"o\">.</span><span
    class=\"n\">warning</span><span class=\"p\">(</span>\n                            <span
    class=\"n\">textwrap</span><span class=\"o\">.</span><span class=\"n\">dedent</span><span
    class=\"p\">(</span>\n                                <span class=\"sa\">f</span><span
    class=\"s2\">&quot;&quot;&quot;</span>\n<span class=\"s2\">                                Warning
    </span><span class=\"si\">{</span><span class=\"n\">key</span><span class=\"si\">}</span><span
    class=\"s2\"> is not set in markata config, sitemap will</span>\n<span class=\"s2\">
    \                               be missing root site_name</span>\n<span class=\"s2\">
    \                               to resolve this open your markata.toml and add</span>\n\n<span
    class=\"s2\">                                </span><span class=\"si\">{</span><span
    class=\"n\">suggested</span><span class=\"si\">}</span>\n\n<span class=\"s2\">
    \                               &quot;&quot;&quot;</span>\n                            <span
    class=\"p\">),</span>\n                        <span class=\"p\">)</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">default</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='make_hash' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>make_hash <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">make_hash
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
    <span class=\"nf\">make_hash</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"o\">*</span><span class=\"n\">keys</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n
    \               <span class=\"kn\">import</span> <span class=\"nn\">xxhash</span>\n\n
    \               <span class=\"n\">str_keys</span> <span class=\"o\">=</span> <span
    class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">)</span> <span class=\"k\">for</span>
    <span class=\"n\">key</span> <span class=\"ow\">in</span> <span class=\"n\">keys</span><span
    class=\"p\">]</span>\n                <span class=\"nb\">hash</span> <span class=\"o\">=</span>
    <span class=\"n\">xxhash</span><span class=\"o\">.</span><span class=\"n\">xxh64</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;&quot;</span><span class=\"o\">.</span><span
    class=\"n\">join</span><span class=\"p\">(</span><span class=\"n\">str_keys</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">encode</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;utf-8&quot;</span><span class=\"p\">))</span><span
    class=\"o\">.</span><span class=\"n\">hexdigest</span><span class=\"p\">()</span>\n
    \               <span class=\"k\">return</span> <span class=\"nb\">hash</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='content_dir_hash' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>content_dir_hash <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">content_dir_hash
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
    <span class=\"nf\">content_dir_hash</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">hashes</span> <span class=\"o\">=</span> <span
    class=\"p\">[</span>\n                    <span class=\"n\">dirhash</span><span
    class=\"p\">(</span><span class=\"nb\">dir</span><span class=\"p\">)</span>\n
    \                   <span class=\"k\">for</span> <span class=\"nb\">dir</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">content_directories</span>\n                    <span class=\"k\">if</span>
    <span class=\"nb\">dir</span><span class=\"o\">.</span><span class=\"n\">absolute</span><span
    class=\"p\">()</span> <span class=\"o\">!=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;.&quot;</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">absolute</span><span class=\"p\">()</span>\n
    \               <span class=\"p\">]</span>\n                <span class=\"k\">return</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">make_hash</span><span
    class=\"p\">(</span><span class=\"o\">*</span><span class=\"n\">hashes</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='console'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>console <em class='small'>method</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">console
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
    <span class=\"nf\">console</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">Console</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">try</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_console</span>\n                <span class=\"k\">except</span> <span
    class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_console</span>
    <span class=\"o\">=</span> <span class=\"n\">Console</span><span class=\"p\">()</span>\n
    \                   <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_console</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='describe' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>describe <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">describe
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
    <span class=\"nf\">describe</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">dict</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
    class=\"p\">]:</span>\n                <span class=\"k\">return</span> <span class=\"p\">{</span><span
    class=\"s2\">&quot;version&quot;</span><span class=\"p\">:</span> <span class=\"n\">__version__</span><span
    class=\"p\">}</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='_to_dict'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_to_dict <em class='small'>method</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_to_dict
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
    <span class=\"nf\">_to_dict</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">dict</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">,</span> <span class=\"n\">Iterable</span><span
    class=\"p\">]:</span>\n                <span class=\"k\">return</span> <span class=\"p\">{</span><span
    class=\"s2\">&quot;config&quot;</span><span class=\"p\">:</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;articles&quot;</span><span class=\"p\">:</span> <span
    class=\"p\">[</span><span class=\"n\">a</span><span class=\"o\">.</span><span
    class=\"n\">to_dict</span><span class=\"p\">()</span> <span class=\"k\">for</span>
    <span class=\"n\">a</span> <span class=\"ow\">in</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">articles</span><span class=\"p\">]}</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='to_dict' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>to_dict <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">to_dict
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
    <span class=\"nf\">to_dict</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">dict</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_to_dict</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='to_json' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>to_json <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">to_json
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
    <span class=\"nf\">to_json</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n
    \               <span class=\"kn\">import</span> <span class=\"nn\">json</span>\n\n
    \               <span class=\"k\">return</span> <span class=\"n\">json</span><span
    class=\"o\">.</span><span class=\"n\">dumps</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">to_dict</span><span
    class=\"p\">(),</span> <span class=\"n\">indent</span><span class=\"o\">=</span><span
    class=\"mi\">4</span><span class=\"p\">,</span> <span class=\"n\">sort_keys</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
    class=\"n\">default</span><span class=\"o\">=</span><span class=\"nb\">str</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='_register_hooks'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_register_hooks
    <em class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">_register_hooks <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">_register_hooks</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">sys</span><span class=\"o\">.</span><span class=\"n\">path</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">os</span><span class=\"o\">.</span><span class=\"n\">getcwd</span><span
    class=\"p\">())</span>\n                <span class=\"k\">for</span> <span class=\"n\">hook</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">hooks_conf</span><span class=\"o\">.</span><span class=\"n\">hooks</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                       <span class=\"c1\"># module style plugins</span>\n                        <span
    class=\"n\">plugin</span> <span class=\"o\">=</span> <span class=\"n\">importlib</span><span
    class=\"o\">.</span><span class=\"n\">import_module</span><span class=\"p\">(</span><span
    class=\"n\">hook</span><span class=\"p\">)</span>\n                    <span class=\"k\">except</span>
    <span class=\"ne\">ModuleNotFoundError</span> <span class=\"k\">as</span> <span
    class=\"n\">e</span><span class=\"p\">:</span>\n                        <span
    class=\"c1\"># class style plugins</span>\n                        <span class=\"k\">if</span>
    <span class=\"s2\">&quot;.&quot;</span> <span class=\"ow\">in</span> <span class=\"n\">hook</span><span
    class=\"p\">:</span>\n                            <span class=\"k\">try</span><span
    class=\"p\">:</span>\n                                <span class=\"n\">mod</span>
    <span class=\"o\">=</span> <span class=\"n\">importlib</span><span class=\"o\">.</span><span
    class=\"n\">import_module</span><span class=\"p\">(</span><span class=\"s2\">&quot;.&quot;</span><span
    class=\"o\">.</span><span class=\"n\">join</span><span class=\"p\">(</span><span
    class=\"n\">hook</span><span class=\"o\">.</span><span class=\"n\">split</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;.&quot;</span><span class=\"p\">)[:</span><span
    class=\"o\">-</span><span class=\"mi\">1</span><span class=\"p\">]))</span>\n
    \                               <span class=\"n\">plugin</span> <span class=\"o\">=</span>
    <span class=\"nb\">getattr</span><span class=\"p\">(</span><span class=\"n\">mod</span><span
    class=\"p\">,</span> <span class=\"n\">hook</span><span class=\"o\">.</span><span
    class=\"n\">split</span><span class=\"p\">(</span><span class=\"s2\">&quot;.&quot;</span><span
    class=\"p\">)[</span><span class=\"o\">-</span><span class=\"mi\">1</span><span
    class=\"p\">])</span>\n                            <span class=\"k\">except</span>
    <span class=\"ne\">ModuleNotFoundError</span> <span class=\"k\">as</span> <span
    class=\"n\">e</span><span class=\"p\">:</span>\n                                <span
    class=\"k\">raise</span> <span class=\"ne\">ModuleNotFoundError</span><span class=\"p\">(</span>\n
    \                                   <span class=\"sa\">f</span><span class=\"s2\">&quot;module
    </span><span class=\"si\">{</span><span class=\"n\">hook</span><span class=\"si\">}</span><span
    class=\"s2\"> not found</span><span class=\"se\">\\n</span><span class=\"si\">{</span><span
    class=\"n\">sys</span><span class=\"o\">.</span><span class=\"n\">path</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span>\n                                <span
    class=\"p\">)</span> <span class=\"kn\">from</span> <span class=\"nn\">e</span>\n
    \                       <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \                           <span class=\"k\">raise</span> <span class=\"n\">e</span>\n\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_pm</span><span class=\"o\">.</span><span class=\"n\">register</span><span
    class=\"p\">(</span><span class=\"n\">plugin</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='__iter__' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'><strong>iter</strong> <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\"><strong>iter</strong>
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
    <span class=\"fm\">__iter__</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">description</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;working...&quot;</span>\n
    \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Iterable</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;Markata.Post&quot;</span><span class=\"p\">]:</span>\n
    \               <span class=\"n\">articles</span><span class=\"p\">:</span> <span
    class=\"n\">Iterable</span><span class=\"p\">[</span><span class=\"n\">Markata</span><span
    class=\"o\">.</span><span class=\"n\">Post</span><span class=\"p\">]</span> <span
    class=\"o\">=</span> <span class=\"n\">track</span><span class=\"p\">(</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">articles</span><span class=\"p\">,</span>\n                    <span
    class=\"n\">description</span><span class=\"o\">=</span><span class=\"n\">description</span><span
    class=\"p\">,</span>\n                    <span class=\"n\">transient</span><span
    class=\"o\">=</span><span class=\"kc\">False</span><span class=\"p\">,</span>\n
    \                   <span class=\"n\">console</span><span class=\"o\">=</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"p\">,</span>\n                <span class=\"p\">)</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">articles</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='iter_articles' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>iter_articles <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">iter_articles
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
    <span class=\"nf\">iter_articles</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">description</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Iterable</span><span
    class=\"p\">[</span><span class=\"n\">Markata</span><span class=\"o\">.</span><span
    class=\"n\">Post</span><span class=\"p\">]:</span>\n                <span class=\"n\">articles</span><span
    class=\"p\">:</span> <span class=\"n\">Iterable</span><span class=\"p\">[</span><span
    class=\"n\">Markata</span><span class=\"o\">.</span><span class=\"n\">Post</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"n\">track</span><span
    class=\"p\">(</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">articles</span><span class=\"p\">,</span>\n
    \                   <span class=\"n\">description</span><span class=\"o\">=</span><span
    class=\"n\">description</span><span class=\"p\">,</span>\n                    <span
    class=\"n\">transient</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span>\n                    <span class=\"n\">console</span><span
    class=\"o\">=</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"p\">,</span>\n                <span class=\"p\">)</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">articles</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='teardown' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>teardown <em class='small'>method</em></h2>\ngive special access to the
    teardown lifecycle method</p>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">teardown <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">teardown</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">Markata</span><span class=\"p\">:</span>\n<span
    class=\"w\">                </span><span class=\"sd\">&quot;&quot;&quot;give special
    access to the teardown lifecycle method&quot;&quot;&quot;</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_pm</span><span
    class=\"o\">.</span><span class=\"n\">hook</span><span class=\"o\">.</span><span
    class=\"n\">teardown</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"bp\">self</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='run' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>run
    <em class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">run <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">run</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">lifecycle</span><span class=\"p\">:</span> <span class=\"n\">LifeCycle</span>
    <span class=\"o\">=</span> <span class=\"kc\">None</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">Markata</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">lifecycle</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">lifecycle</span> <span class=\"o\">=</span>
    <span class=\"nb\">max</span><span class=\"p\">(</span><span class=\"n\">LifeCycle</span><span
    class=\"o\">.</span><span class=\"n\">_member_map_</span><span class=\"o\">.</span><span
    class=\"n\">values</span><span class=\"p\">())</span>\n\n                <span
    class=\"k\">if</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">lifecycle</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
    class=\"p\">):</span>\n                    <span class=\"n\">lifecycle</span>
    <span class=\"o\">=</span> <span class=\"n\">LifeCycle</span><span class=\"p\">[</span><span
    class=\"n\">lifecycle</span><span class=\"p\">]</span>\n\n                <span
    class=\"n\">stages_to_run</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
    \                   <span class=\"n\">m</span>\n                    <span class=\"k\">for</span>
    <span class=\"n\">m</span> <span class=\"ow\">in</span> <span class=\"n\">LifeCycle</span><span
    class=\"o\">.</span><span class=\"n\">_member_map_</span>\n                    <span
    class=\"k\">if</span> <span class=\"p\">(</span><span class=\"n\">LifeCycle</span><span
    class=\"p\">[</span><span class=\"n\">m</span><span class=\"p\">]</span> <span
    class=\"o\">&lt;=</span> <span class=\"n\">lifecycle</span><span class=\"p\">)</span>
    <span class=\"ow\">and</span> <span class=\"p\">(</span><span class=\"n\">m</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">stages_ran</span><span class=\"p\">)</span>\n
    \               <span class=\"p\">]</span>\n\n                <span class=\"k\">if</span>
    <span class=\"ow\">not</span> <span class=\"n\">stages_to_run</span><span class=\"p\">:</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
    class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">lifecycle</span><span class=\"o\">.</span><span
    class=\"n\">name</span><span class=\"si\">}</span><span class=\"s2\"> already
    ran&quot;</span><span class=\"p\">)</span>\n                    <span class=\"k\">return</span>
    <span class=\"bp\">self</span>\n\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">log</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;running </span><span class=\"si\">{</span><span class=\"n\">stages_to_run</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">for</span> <span class=\"n\">stage</span> <span
    class=\"ow\">in</span> <span class=\"n\">stages_to_run</span><span class=\"p\">:</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
    class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">stage</span><span class=\"si\">}</span><span
    class=\"s2\"> running&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"nb\">getattr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_pm</span><span class=\"o\">.</span><span
    class=\"n\">hook</span><span class=\"p\">,</span> <span class=\"n\">stage</span><span
    class=\"p\">)(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"bp\">self</span><span class=\"p\">)</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">stages_ran</span><span
    class=\"o\">.</span><span class=\"n\">add</span><span class=\"p\">(</span><span
    class=\"n\">stage</span><span class=\"p\">)</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">log</span><span class=\"p\">(</span><span
    class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">stage</span><span class=\"si\">}</span><span class=\"s2\"> complete&quot;</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">with</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">cache</span> <span class=\"k\">as</span>
    <span class=\"n\">cache</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">hits</span><span class=\"p\">,</span> <span class=\"n\">misses</span>
    <span class=\"o\">=</span> <span class=\"n\">cache</span><span class=\"o\">.</span><span
    class=\"n\">stats</span><span class=\"p\">()</span>\n\n                <span class=\"k\">if</span>
    <span class=\"n\">hits</span> <span class=\"o\">+</span> <span class=\"n\">misses</span>
    <span class=\"o\">&gt;</span> <span class=\"mi\">0</span><span class=\"p\">:</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
    class=\"p\">(</span>\n                        <span class=\"sa\">f</span><span
    class=\"s2\">&quot;lifetime cache hit rate </span><span class=\"si\">{</span><span
    class=\"nb\">round</span><span class=\"p\">(</span><span class=\"n\">hits</span><span
    class=\"o\">/</span><span class=\"w\"> </span><span class=\"p\">(</span><span
    class=\"n\">hits</span><span class=\"w\"> </span><span class=\"o\">+</span><span
    class=\"w\"> </span><span class=\"n\">misses</span><span class=\"p\">)</span><span
    class=\"o\">*</span><span class=\"mi\">100</span><span class=\"p\">,</span><span
    class=\"w\"> </span><span class=\"mi\">2</span><span class=\"p\">)</span><span
    class=\"si\">}</span><span class=\"s2\">%&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"p\">)</span>\n\n                <span class=\"k\">if</span>
    <span class=\"n\">misses</span> <span class=\"o\">&gt;</span> <span class=\"mi\">0</span><span
    class=\"p\">:</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">log</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;lifetime cache hits/misses </span><span class=\"si\">{</span><span
    class=\"n\">hits</span><span class=\"si\">}</span><span class=\"s2\">/</span><span
    class=\"si\">{</span><span class=\"n\">misses</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span><span class=\"p\">)</span>\n\n                <span
    class=\"n\">hits</span> <span class=\"o\">-=</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">init_cache_stats</span><span class=\"p\">[</span><span
    class=\"mi\">0</span><span class=\"p\">]</span>\n                <span class=\"n\">misses</span>
    <span class=\"o\">-=</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">init_cache_stats</span><span class=\"p\">[</span><span class=\"mi\">1</span><span
    class=\"p\">]</span>\n\n                <span class=\"k\">if</span> <span class=\"n\">hits</span>
    <span class=\"o\">+</span> <span class=\"n\">misses</span> <span class=\"o\">&gt;</span>
    <span class=\"mi\">0</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">log</span><span class=\"p\">(</span>\n                        <span
    class=\"sa\">f</span><span class=\"s2\">&quot;run cache hit rate </span><span
    class=\"si\">{</span><span class=\"nb\">round</span><span class=\"p\">(</span><span
    class=\"n\">hits</span><span class=\"o\">/</span><span class=\"w\"> </span><span
    class=\"p\">(</span><span class=\"n\">hits</span><span class=\"w\"> </span><span
    class=\"o\">+</span><span class=\"w\"> </span><span class=\"n\">misses</span><span
    class=\"p\">)</span><span class=\"o\">*</span><span class=\"mi\">100</span><span
    class=\"p\">,</span><span class=\"w\"> </span><span class=\"mi\">2</span><span
    class=\"p\">)</span><span class=\"si\">}</span><span class=\"s2\">%&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"p\">)</span>\n\n                <span
    class=\"k\">if</span> <span class=\"n\">misses</span> <span class=\"o\">&gt;</span>
    <span class=\"mi\">0</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">log</span><span class=\"p\">(</span><span
    class=\"sa\">f</span><span class=\"s2\">&quot;run cache hits/misses </span><span
    class=\"si\">{</span><span class=\"n\">hits</span><span class=\"si\">}</span><span
    class=\"s2\">/</span><span class=\"si\">{</span><span class=\"n\">misses</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">)</span>\n\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='filter' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>filter <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">filter
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
    <span class=\"nf\">filter</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span>
    <span class=\"nb\">filter</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">list</span><span
    class=\"p\">:</span>\n                <span class=\"k\">def</span> <span class=\"nf\">evalr</span><span
    class=\"p\">(</span><span class=\"n\">a</span><span class=\"p\">:</span> <span
    class=\"n\">Markata</span><span class=\"o\">.</span><span class=\"n\">Post</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                       <span class=\"k\">return</span> <span class=\"nb\">eval</span><span
    class=\"p\">(</span>\n                            <span class=\"nb\">filter</span><span
    class=\"p\">,</span>\n                            <span class=\"p\">{</span><span
    class=\"o\">**</span><span class=\"n\">a</span><span class=\"o\">.</span><span
    class=\"n\">to_dict</span><span class=\"p\">(),</span> <span class=\"s2\">&quot;timedelta&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">timedelta</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;post&quot;</span><span class=\"p\">:</span> <span class=\"n\">a</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;m&quot;</span><span class=\"p\">:</span>
    <span class=\"bp\">self</span><span class=\"p\">},</span>\n                            <span
    class=\"p\">{},</span>\n                        <span class=\"p\">)</span>\n                    <span
    class=\"k\">except</span> <span class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n
    \                       <span class=\"k\">return</span> <span class=\"nb\">eval</span><span
    class=\"p\">(</span>\n                            <span class=\"nb\">filter</span><span
    class=\"p\">,</span>\n                            <span class=\"p\">{</span><span
    class=\"o\">**</span><span class=\"n\">a</span><span class=\"o\">.</span><span
    class=\"n\">to_dict</span><span class=\"p\">(),</span> <span class=\"s2\">&quot;timedelta&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">timedelta</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;post&quot;</span><span class=\"p\">:</span> <span class=\"n\">a</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;m&quot;</span><span class=\"p\">:</span>
    <span class=\"bp\">self</span><span class=\"p\">},</span>\n                            <span
    class=\"p\">{},</span>\n                        <span class=\"p\">)</span>\n\n
    \               <span class=\"k\">return</span> <span class=\"p\">[</span><span
    class=\"n\">a</span> <span class=\"k\">for</span> <span class=\"n\">a</span> <span
    class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">articles</span> <span class=\"k\">if</span> <span class=\"n\">evalr</span><span
    class=\"p\">(</span><span class=\"n\">a</span><span class=\"p\">)]</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='map' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>map
    <em class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">map <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">map</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">func</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;title&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"nb\">filter</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;True&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">sort</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;True&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">reverse</span><span class=\"p\">:</span>
    <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"kc\">True</span><span
    class=\"p\">,</span>\n                <span class=\"o\">*</span><span class=\"n\">args</span><span
    class=\"p\">:</span> <span class=\"nb\">tuple</span><span class=\"p\">,</span>\n
    \               <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">:</span> <span class=\"nb\">dict</span><span class=\"p\">,</span>\n
    \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">list</span><span
    class=\"p\">:</span>\n                <span class=\"kn\">import</span> <span class=\"nn\">copy</span>\n\n
    \               <span class=\"k\">def</span> <span class=\"nf\">try_sort</span><span
    class=\"p\">(</span><span class=\"n\">a</span><span class=\"p\">:</span> <span
    class=\"n\">Any</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"nb\">int</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">if</span> <span class=\"s2\">&quot;datetime&quot;</span> <span class=\"ow\">in</span>
    <span class=\"n\">sort</span><span class=\"o\">.</span><span class=\"n\">lower</span><span
    class=\"p\">():</span>\n                        <span class=\"k\">return</span>
    <span class=\"n\">a</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"n\">sort</span><span class=\"p\">,</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"p\">(</span><span class=\"mi\">1970</span><span class=\"p\">,</span> <span
    class=\"mi\">1</span><span class=\"p\">,</span> <span class=\"mi\">1</span><span
    class=\"p\">))</span>\n\n                    <span class=\"k\">if</span> <span
    class=\"s2\">&quot;date&quot;</span> <span class=\"ow\">in</span> <span class=\"n\">sort</span><span
    class=\"o\">.</span><span class=\"n\">lower</span><span class=\"p\">():</span>\n
    \                       <span class=\"k\">return</span> <span class=\"n\">a</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"n\">sort</span><span class=\"p\">,</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">date</span><span class=\"p\">(</span><span
    class=\"mi\">1970</span><span class=\"p\">,</span> <span class=\"mi\">1</span><span
    class=\"p\">,</span> <span class=\"mi\">1</span><span class=\"p\">))</span>\n\n
    \                   <span class=\"k\">try</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">value</span> <span class=\"o\">=</span> <span class=\"nb\">eval</span><span
    class=\"p\">(</span><span class=\"n\">sort</span><span class=\"p\">,</span> <span
    class=\"n\">a</span><span class=\"o\">.</span><span class=\"n\">to_dict</span><span
    class=\"p\">(),</span> <span class=\"p\">{})</span>\n                    <span
    class=\"k\">except</span> <span class=\"ne\">NameError</span><span class=\"p\">:</span>\n
    \                       <span class=\"k\">return</span> <span class=\"o\">-</span><span
    class=\"mi\">1</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">value</span>\n                    <span class=\"k\">try</span><span
    class=\"p\">:</span>\n                        <span class=\"k\">return</span>
    <span class=\"nb\">int</span><span class=\"p\">(</span><span class=\"n\">value</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">except</span> <span
    class=\"ne\">TypeError</span><span class=\"p\">:</span>\n                        <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                            <span
    class=\"k\">return</span> <span class=\"nb\">int</span><span class=\"p\">(</span><span
    class=\"n\">value</span><span class=\"o\">.</span><span class=\"n\">timestamp</span><span
    class=\"p\">())</span>\n                        <span class=\"k\">except</span>
    <span class=\"ne\">Exception</span><span class=\"p\">:</span>\n                            <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                                <span
    class=\"k\">return</span> <span class=\"nb\">int</span><span class=\"p\">(</span>\n
    \                                   <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">combine</span><span
    class=\"p\">(</span>\n                                        <span class=\"n\">value</span><span
    class=\"p\">,</span>\n                                        <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">min</span><span class=\"o\">.</span><span class=\"n\">time</span><span
    class=\"p\">(),</span>\n                                    <span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">timestamp</span><span class=\"p\">(),</span>\n
    \                               <span class=\"p\">)</span>\n                            <span
    class=\"k\">except</span> <span class=\"ne\">Exception</span><span class=\"p\">:</span>\n
    \                               <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                                   <span class=\"k\">return</span> <span class=\"nb\">sum</span><span
    class=\"p\">([</span><span class=\"nb\">ord</span><span class=\"p\">(</span><span
    class=\"n\">c</span><span class=\"p\">)</span> <span class=\"k\">for</span> <span
    class=\"n\">c</span> <span class=\"ow\">in</span> <span class=\"nb\">str</span><span
    class=\"p\">(</span><span class=\"n\">value</span><span class=\"p\">)])</span>\n
    \                               <span class=\"k\">except</span> <span class=\"ne\">Exception</span><span
    class=\"p\">:</span>\n                                    <span class=\"k\">return</span>
    <span class=\"o\">-</span><span class=\"mi\">1</span>\n\n                <span
    class=\"n\">articles</span> <span class=\"o\">=</span> <span class=\"n\">copy</span><span
    class=\"o\">.</span><span class=\"n\">copy</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">articles</span><span
    class=\"p\">)</span>\n                <span class=\"n\">articles</span><span class=\"o\">.</span><span
    class=\"n\">sort</span><span class=\"p\">(</span><span class=\"n\">key</span><span
    class=\"o\">=</span><span class=\"n\">try_sort</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">reverse</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">articles</span><span
    class=\"o\">.</span><span class=\"n\">reverse</span><span class=\"p\">()</span>\n\n
    \               <span class=\"k\">try</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">posts</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
    \                       <span class=\"nb\">eval</span><span class=\"p\">(</span>\n
    \                           <span class=\"n\">func</span><span class=\"p\">,</span>\n
    \                           <span class=\"p\">{</span><span class=\"o\">**</span><span
    class=\"n\">a</span><span class=\"o\">.</span><span class=\"n\">to_dict</span><span
    class=\"p\">(),</span> <span class=\"s2\">&quot;timedelta&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">timedelta</span><span class=\"p\">,</span> <span class=\"s2\">&quot;post&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">a</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;m&quot;</span><span class=\"p\">:</span> <span class=\"bp\">self</span><span
    class=\"p\">},</span>\n                            <span class=\"p\">{},</span>\n
    \                       <span class=\"p\">)</span>\n                        <span
    class=\"k\">for</span> <span class=\"n\">a</span> <span class=\"ow\">in</span>
    <span class=\"n\">articles</span>\n                        <span class=\"k\">if</span>
    <span class=\"nb\">eval</span><span class=\"p\">(</span>\n                            <span
    class=\"nb\">filter</span><span class=\"p\">,</span>\n                            <span
    class=\"p\">{</span><span class=\"o\">**</span><span class=\"n\">a</span><span
    class=\"o\">.</span><span class=\"n\">to_dict</span><span class=\"p\">(),</span>
    <span class=\"s2\">&quot;timedelta&quot;</span><span class=\"p\">:</span> <span
    class=\"n\">timedelta</span><span class=\"p\">,</span> <span class=\"s2\">&quot;post&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">a</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;m&quot;</span><span class=\"p\">:</span> <span class=\"bp\">self</span><span
    class=\"p\">},</span>\n                            <span class=\"p\">{},</span>\n
    \                       <span class=\"p\">)</span>\n                    <span
    class=\"p\">]</span>\n\n                <span class=\"k\">except</span> <span
    class=\"ne\">NameError</span> <span class=\"k\">as</span> <span class=\"n\">e</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">variable</span> <span
    class=\"o\">=</span> <span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">e</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">split</span><span class=\"p\">(</span><span class=\"s2\">&quot;&#39;&quot;</span><span
    class=\"p\">)[</span><span class=\"mi\">1</span><span class=\"p\">]</span>\n\n
    \                   <span class=\"n\">missing_in_posts</span> <span class=\"o\">=</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">map</span><span
    class=\"p\">(</span>\n                        <span class=\"s2\">&quot;path&quot;</span><span
    class=\"p\">,</span>\n                        <span class=\"nb\">filter</span><span
    class=\"o\">=</span><span class=\"sa\">f</span><span class=\"s1\">&#39;&quot;</span><span
    class=\"si\">{</span><span class=\"n\">variable</span><span class=\"si\">}</span><span
    class=\"s1\">&quot; not in post.keys()&#39;</span><span class=\"p\">,</span>\n
    \                   <span class=\"p\">)</span>\n                    <span class=\"n\">message</span>
    <span class=\"o\">=</span> <span class=\"p\">(</span>\n                        <span
    class=\"sa\">f</span><span class=\"s2\">&quot;variable: &#39;</span><span class=\"si\">{</span><span
    class=\"n\">variable</span><span class=\"si\">}</span><span class=\"s2\">&#39;
    is missing in </span><span class=\"si\">{</span><span class=\"nb\">len</span><span
    class=\"p\">(</span><span class=\"n\">missing_in_posts</span><span class=\"p\">)</span><span
    class=\"si\">}</span><span class=\"s2\"> posts&quot;</span>\n                    <span
    class=\"p\">)</span>\n                    <span class=\"k\">if</span> <span class=\"nb\">len</span><span
    class=\"p\">(</span><span class=\"n\">missing_in_posts</span><span class=\"p\">)</span>
    <span class=\"o\">&gt;</span> <span class=\"mi\">10</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">message</span> <span class=\"o\">+=</span>
    <span class=\"p\">(</span>\n                            <span class=\"sa\">f</span><span
    class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span class=\"s2\">first
    10 paths to posts missing </span><span class=\"si\">{</span><span class=\"n\">variable</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span>\n                            <span
    class=\"sa\">f</span><span class=\"s2\">&quot;[</span><span class=\"si\">{</span><span
    class=\"s1\">&#39;,&#39;</span><span class=\"o\">.</span><span class=\"n\">join</span><span
    class=\"p\">([</span><span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">p</span><span class=\"p\">)</span><span class=\"w\"> </span><span
    class=\"k\">for</span><span class=\"w\"> </span><span class=\"n\">p</span><span
    class=\"w\"> </span><span class=\"ow\">in</span><span class=\"w\"> </span><span
    class=\"n\">missing_in_posts</span><span class=\"p\">[:</span><span class=\"mi\">10</span><span
    class=\"p\">]])</span><span class=\"si\">}</span><span class=\"s2\">...&quot;</span>\n
    \                       <span class=\"p\">)</span>\n                    <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">message</span> <span class=\"o\">+=</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span class=\"s2\">paths
    to posts missing </span><span class=\"si\">{</span><span class=\"n\">variable</span><span
    class=\"si\">}</span><span class=\"s2\"> </span><span class=\"si\">{</span><span
    class=\"n\">missing_in_posts</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span>\n\n
    \                   <span class=\"k\">raise</span> <span class=\"n\">MissingFrontMatter</span><span
    class=\"p\">(</span><span class=\"n\">message</span><span class=\"p\">)</span>\n\n
    \               <span class=\"k\">return</span> <span class=\"n\">posts</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='first' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>first
    <em class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">first <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">first</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"nb\">filter</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;True&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">sort</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;True&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">reverse</span><span class=\"p\">:</span>
    <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"kc\">True</span><span
    class=\"p\">,</span>\n                <span class=\"o\">*</span><span class=\"n\">args</span><span
    class=\"p\">:</span> <span class=\"nb\">tuple</span><span class=\"p\">,</span>\n
    \               <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">:</span> <span class=\"nb\">dict</span><span class=\"p\">,</span>\n
    \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">list</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">map</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;post&quot;</span><span class=\"p\">,</span> <span class=\"nb\">filter</span><span
    class=\"p\">,</span> <span class=\"n\">sort</span><span class=\"p\">,</span> <span
    class=\"n\">reverse</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"n\">args</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"n\">kwargs</span><span class=\"p\">)[</span><span class=\"mi\">0</span><span
    class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='last' class='admonition-title'
    style='margin:0;padding:.5rem 1rem;'>last <em class='small'>method</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">last
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
    <span class=\"nf\">last</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"nb\">filter</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;True&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">sort</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;True&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"n\">reverse</span><span class=\"p\">:</span>
    <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"kc\">True</span><span
    class=\"p\">,</span>\n                <span class=\"o\">*</span><span class=\"n\">args</span><span
    class=\"p\">:</span> <span class=\"nb\">tuple</span><span class=\"p\">,</span>\n
    \               <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">:</span> <span class=\"nb\">dict</span><span class=\"p\">,</span>\n
    \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">list</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">map</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;post&quot;</span><span class=\"p\">,</span> <span class=\"nb\">filter</span><span
    class=\"p\">,</span> <span class=\"n\">sort</span><span class=\"p\">,</span> <span
    class=\"n\">reverse</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"n\">args</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"n\">kwargs</span><span class=\"p\">)[</span><span class=\"o\">-</span><span
    class=\"mi\">1</span><span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='one' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>one
    <em class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">one <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">one</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"nb\">filter</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;True&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"o\">*</span><span class=\"n\">args</span><span
    class=\"p\">:</span> <span class=\"nb\">tuple</span><span class=\"p\">,</span>\n
    \               <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">:</span> <span class=\"nb\">dict</span><span class=\"p\">,</span>\n
    \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">list</span><span
    class=\"p\">:</span>\n                <span class=\"n\">posts</span> <span class=\"o\">=</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">map</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;post&quot;</span><span class=\"p\">,</span>
    <span class=\"nb\">filter</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"n\">args</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"n\">kwargs</span><span class=\"p\">)</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">len</span><span class=\"p\">(</span><span class=\"n\">posts</span><span
    class=\"p\">)</span> <span class=\"o\">&gt;</span> <span class=\"mi\">1</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">raise</span> <span
    class=\"n\">TooManyPosts</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;found </span><span class=\"si\">{</span><span class=\"nb\">len</span><span
    class=\"p\">(</span><span class=\"n\">posts</span><span class=\"p\">)</span><span
    class=\"si\">}</span><span class=\"s2\"> posts, expected 1. </span><span class=\"si\">{</span><span
    class=\"n\">posts</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"k\">if</span> <span class=\"nb\">len</span><span
    class=\"p\">(</span><span class=\"n\">posts</span><span class=\"p\">)</span> <span
    class=\"o\">==</span> <span class=\"mi\">0</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">raise</span> <span class=\"n\">NoPosts</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">posts</span><span
    class=\"p\">[</span><span class=\"mi\">0</span><span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='evalr' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>evalr <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">evalr
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
    <span class=\"nf\">evalr</span><span class=\"p\">(</span><span class=\"n\">a</span><span
    class=\"p\">:</span> <span class=\"n\">Markata</span><span class=\"o\">.</span><span
    class=\"n\">Post</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Any</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                        <span
    class=\"k\">return</span> <span class=\"nb\">eval</span><span class=\"p\">(</span>\n
    \                           <span class=\"nb\">filter</span><span class=\"p\">,</span>\n
    \                           <span class=\"p\">{</span><span class=\"o\">**</span><span
    class=\"n\">a</span><span class=\"o\">.</span><span class=\"n\">to_dict</span><span
    class=\"p\">(),</span> <span class=\"s2\">&quot;timedelta&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">timedelta</span><span class=\"p\">,</span> <span class=\"s2\">&quot;post&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">a</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;m&quot;</span><span class=\"p\">:</span> <span class=\"bp\">self</span><span
    class=\"p\">},</span>\n                            <span class=\"p\">{},</span>\n
    \                       <span class=\"p\">)</span>\n                    <span
    class=\"k\">except</span> <span class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n
    \                       <span class=\"k\">return</span> <span class=\"nb\">eval</span><span
    class=\"p\">(</span>\n                            <span class=\"nb\">filter</span><span
    class=\"p\">,</span>\n                            <span class=\"p\">{</span><span
    class=\"o\">**</span><span class=\"n\">a</span><span class=\"o\">.</span><span
    class=\"n\">to_dict</span><span class=\"p\">(),</span> <span class=\"s2\">&quot;timedelta&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">timedelta</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;post&quot;</span><span class=\"p\">:</span> <span class=\"n\">a</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;m&quot;</span><span class=\"p\">:</span>
    <span class=\"bp\">self</span><span class=\"p\">},</span>\n                            <span
    class=\"p\">{},</span>\n                        <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='try_sort' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>try_sort <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">try_sort
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
    <span class=\"nf\">try_sort</span><span class=\"p\">(</span><span class=\"n\">a</span><span
    class=\"p\">:</span> <span class=\"n\">Any</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"nb\">int</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">if</span> <span class=\"s2\">&quot;datetime&quot;</span>
    <span class=\"ow\">in</span> <span class=\"n\">sort</span><span class=\"o\">.</span><span
    class=\"n\">lower</span><span class=\"p\">():</span>\n                        <span
    class=\"k\">return</span> <span class=\"n\">a</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"n\">sort</span><span
    class=\"p\">,</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">datetime</span><span class=\"p\">(</span><span class=\"mi\">1970</span><span
    class=\"p\">,</span> <span class=\"mi\">1</span><span class=\"p\">,</span> <span
    class=\"mi\">1</span><span class=\"p\">))</span>\n\n                    <span
    class=\"k\">if</span> <span class=\"s2\">&quot;date&quot;</span> <span class=\"ow\">in</span>
    <span class=\"n\">sort</span><span class=\"o\">.</span><span class=\"n\">lower</span><span
    class=\"p\">():</span>\n                        <span class=\"k\">return</span>
    <span class=\"n\">a</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"n\">sort</span><span class=\"p\">,</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">date</span><span
    class=\"p\">(</span><span class=\"mi\">1970</span><span class=\"p\">,</span> <span
    class=\"mi\">1</span><span class=\"p\">,</span> <span class=\"mi\">1</span><span
    class=\"p\">))</span>\n\n                    <span class=\"k\">try</span><span
    class=\"p\">:</span>\n                        <span class=\"n\">value</span> <span
    class=\"o\">=</span> <span class=\"nb\">eval</span><span class=\"p\">(</span><span
    class=\"n\">sort</span><span class=\"p\">,</span> <span class=\"n\">a</span><span
    class=\"o\">.</span><span class=\"n\">to_dict</span><span class=\"p\">(),</span>
    <span class=\"p\">{})</span>\n                    <span class=\"k\">except</span>
    <span class=\"ne\">NameError</span><span class=\"p\">:</span>\n                        <span
    class=\"k\">return</span> <span class=\"o\">-</span><span class=\"mi\">1</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">value</span>\n
    \                   <span class=\"k\">try</span><span class=\"p\">:</span>\n                        <span
    class=\"k\">return</span> <span class=\"nb\">int</span><span class=\"p\">(</span><span
    class=\"n\">value</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">except</span> <span class=\"ne\">TypeError</span><span class=\"p\">:</span>\n
    \                       <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                           <span class=\"k\">return</span> <span class=\"nb\">int</span><span
    class=\"p\">(</span><span class=\"n\">value</span><span class=\"o\">.</span><span
    class=\"n\">timestamp</span><span class=\"p\">())</span>\n                        <span
    class=\"k\">except</span> <span class=\"ne\">Exception</span><span class=\"p\">:</span>\n
    \                           <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                               <span class=\"k\">return</span> <span class=\"nb\">int</span><span
    class=\"p\">(</span>\n                                    <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">combine</span><span class=\"p\">(</span>\n                                        <span
    class=\"n\">value</span><span class=\"p\">,</span>\n                                        <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">min</span><span class=\"o\">.</span><span
    class=\"n\">time</span><span class=\"p\">(),</span>\n                                    <span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">timestamp</span><span
    class=\"p\">(),</span>\n                                <span class=\"p\">)</span>\n
    \                           <span class=\"k\">except</span> <span class=\"ne\">Exception</span><span
    class=\"p\">:</span>\n                                <span class=\"k\">try</span><span
    class=\"p\">:</span>\n                                    <span class=\"k\">return</span>
    <span class=\"nb\">sum</span><span class=\"p\">([</span><span class=\"nb\">ord</span><span
    class=\"p\">(</span><span class=\"n\">c</span><span class=\"p\">)</span> <span
    class=\"k\">for</span> <span class=\"n\">c</span> <span class=\"ow\">in</span>
    <span class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">value</span><span
    class=\"p\">)])</span>\n                                <span class=\"k\">except</span>
    <span class=\"ne\">Exception</span><span class=\"p\">:</span>\n                                    <span
    class=\"k\">return</span> <span class=\"o\">-</span><span class=\"mi\">1</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>"
  raw.md: ''
published: true
slug: markata/init
title: __Init__.Py


---

Markata is a tool for handling directories of markdown.


!! class <h2 id='HooksConfig' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>HooksConfig <em class='small'>class</em></h2>

???+ source "HooksConfig <em class='small'>source</em>"

```python

        class HooksConfig(pydantic.BaseModel):
            hooks: list = ["default"]
            disabled_hooks: list = []
```


!! class <h2 id='Markata' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Markata <em class='small'>class</em></h2>

???+ source "Markata <em class='small'>source</em>"

```python

        class Markata:
            def __init__(self: "Markata", console: Console = None, config=None) -> None:
                self.__version__ = __version__
                self.stages_ran = set()
                self.threded = False
                self._cache = None
                self._precache = None
                self.MARKATA_CACHE_DIR = Path(".") / ".markata.cache"
                self.MARKATA_CACHE_DIR.mkdir(exist_ok=True)
                self._pm = pluggy.PluginManager("markata")
                self._pm.add_hookspecs(hookspec.MarkataSpecs)
                if config is not None:
                    self.config = config
                with self.cache as cache:
                    self.init_cache_stats = cache.stats()
                self.registered_attrs = hookspec.registered_attrs
                self.post_models = []
                self.config_models = []
                if config is not None:
                    raw_hooks = config
                else:
                    raw_hooks = standard_config.load("markata")
                self.hooks_conf = HooksConfig.parse_obj(raw_hooks)
                try:
                    default_index = self.hooks_conf.hooks.index("default")
                    hooks = [
                        *self.hooks_conf.hooks[:default_index],
                        *DEFAULT_HOOKS,
                        *self.hooks_conf.hooks[default_index + 1 :],
                    ]
                    self.hooks_conf.hooks = [
                        hook for hook in hooks if hook not in self.hooks_conf.disabled_hooks
                    ]
                except ValueError:
                    # 'default' is not in hooks , do not replace with default_hooks
                    pass

                self._register_hooks()
                if console is not None:
                    self._console = console
                atexit.register(self.teardown)
                self.precache

            @property
            def cache(self: "Markata") -> Cache:
                # if self.threded:
                #     FanoutCache(self.MARKATA_CACHE_DIR, statistics=True)
                if self._cache is not None:
                    return self._cache
                self._cache = Cache(self.MARKATA_CACHE_DIR, statistics=True)

                return self._cache

            @property
            def precache(self: "Markata") -> None:
                if self._precache is None:
                    self.cache.expire()
                    self._precache = {k: self.cache.get(k) for k in self.cache.iterkeys()}
                return self._precache

            def __getattr__(self: "Markata", item: str) -> Any:
                if item in self._pm.hook.__dict__:
                    # item is a hook, return a callable function
                    return lambda: self.run(item)

                if item in self.__dict__:
                    # item is an attribute, return it
                    return self.__getitem__(item)

                elif item in self.registered_attrs:
                    # item is created by a plugin, run it
                    stage_to_run_to = max(
                        [attr["lifecycle"] for attr in self.registered_attrs[item]],
                    ).name
                    self.console.log(
                        f"Running to [purple]{stage_to_run_to}[/] to retrieve [purple]{item}[/]"
                    )
                    self.run(stage_to_run_to)
                    return getattr(self, item)
                elif item == "precache":
                    return self._precache or {}
                else:
                    # Markata does not know what this is, raise
                    raise AttributeError(f"'Markata' object has no attribute '{item}'")

            def __rich__(self: "Markata") -> Table:
                grid = Table.grid()
                grid.add_column("label")
                grid.add_column("value")

                for label, value in self.describe().items():
                    grid.add_row(label, value)

                return grid

            def bust_cache(self: "Markata") -> Markata:
                with self.cache as cache:
                    cache.clear()
                return self

            def configure(self) -> Markata:
                sys.path.append(os.getcwd())
                # self.config = {**DEFUALT_CONFIG, **standard_config.load("markata")}
                # if isinstance(self.config["glob_patterns"], str):
                #     self.config["glob_patterns"] = self.config["glob_patterns"].split(",")
                # elif isinstance(self.config["glob_patterns"], list):
                #     self.config["glob_patterns"] = list(self.config["glob_patterns"])
                # else:
                #     raise TypeError("glob_patterns must be list or str")
                # self.glob_patterns = self.config["glob_patterns"]

                # self.hooks = self.config["hooks"]

                # if "disabled_hooks" not in self.config:
                #     self.disabled_hooks = [""]
                # if isinstance(self.config["disabled_hooks"], str):
                #     self.disabled_hooks = self.config["disabled_hooks"].split(",")
                # if isinstance(self.config["disabled_hooks"], list):
                #     self.disabled_hooks = self.config["disabled_hooks"]

                # if not self.config.get("output_dir", "markout").endswith(
                #     self.config.get("path_prefix", "")
                # ):
                #     self.config["output_dir"] = (
                #         self.config.get("output_dir", "markout") +
                #         "/" +
                #         self.config.get("path_prefix", "").rstrip("/")
                #     )
                # if (
                #     len((output_split := self.config.get("output_dir", "markout").split("/"))) >
                #     1
                # ):
                #     if "path_prefix" not in self.config.keys():
                #         self.config["path_prefix"] = "/".join(output_split[1:]) + "/"
                # if not self.config.get("path_prefix", "").endswith("/"):
                #     self.config["path_prefix"] = self.config.get("path_prefix", "") + "/"

                # self.config["output_dir"] = self.config["output_dir"].lstrip("/")
                # self.config["path_prefix"] = self.config["path_prefix"].lstrip("/")

                try:
                    default_index = self.hooks_conf.hooks.index("default")
                    hooks = [
                        *self.hooks_conf.hooks[:default_index],
                        *DEFAULT_HOOKS,
                        *self.hooks_conf.hooks[default_index + 1 :],
                    ]
                    self.config.hooks = [
                        hook for hook in hooks if hook not in self.config.disabled_hooks
                    ]
                except ValueError:
                    # 'default' is not in hooks , do not replace with default_hooks
                    pass

                self._pm = pluggy.PluginManager("markata")
                self._pm.add_hookspecs(hookspec.MarkataSpecs)
                self._register_hooks()

                self._pm.hook.configure(markata=self)
                return self

            def get_plugin_config(self, path_or_name: str) -> Dict:
                key = Path(path_or_name).stem

                config = self.config.get(key, {})

                if not isinstance(config, dict):
                    raise TypeError("must use dict")
                if "cache_expire" not in config.keys():
                    config["cache_expire"] = self.config["default_cache_expire"]
                if "config_key" not in config.keys():
                    config["config_key"] = key
                return config

            def get_config(
                self,
                key: str,
                default: str = "",
                warn: bool = True,
                suggested: Optional[str] = None,
            ) -> Any:
                if key in self.config.keys():
                    return self.config[key]
                else:
                    if suggested is None:
                        suggested = textwrap.dedent(
                            f"""
                            [markata]
                            {key} = '{default}'
                            """
                        )
                    if warn:
                        logger.warning(
                            textwrap.dedent(
                                f"""
                                Warning {key} is not set in markata config, sitemap will
                                be missing root site_name
                                to resolve this open your markata.toml and add

                                {suggested}

                                """
                            ),
                        )
                return default

            def make_hash(self, *keys: str) -> str:
                import xxhash

                str_keys = [str(key) for key in keys]
                hash = xxhash.xxh64("".join(str_keys).encode("utf-8")).hexdigest()
                return hash

            @property
            def content_dir_hash(self: "Markata") -> str:
                hashes = [
                    dirhash(dir)
                    for dir in self.content_directories
                    if dir.absolute() != Path(".").absolute()
                ]
                return self.make_hash(*hashes)

            @property
            def console(self: "Markata") -> Console:
                try:
                    return self._console
                except AttributeError:
                    self._console = Console()
                    return self._console

            def describe(self: "Markata") -> dict[str, str]:
                return {"version": __version__}

            def _to_dict(self: "Markata") -> dict[str, Iterable]:
                return {"config": self.config, "articles": [a.to_dict() for a in self.articles]}

            def to_dict(self: "Markata") -> dict:
                return self._to_dict()

            def to_json(self: "Markata") -> str:
                import json

                return json.dumps(self.to_dict(), indent=4, sort_keys=True, default=str)

            def _register_hooks(self: "Markata") -> None:
                sys.path.append(os.getcwd())
                for hook in self.hooks_conf.hooks:
                    try:
                        # module style plugins
                        plugin = importlib.import_module(hook)
                    except ModuleNotFoundError as e:
                        # class style plugins
                        if "." in hook:
                            try:
                                mod = importlib.import_module(".".join(hook.split(".")[:-1]))
                                plugin = getattr(mod, hook.split(".")[-1])
                            except ModuleNotFoundError as e:
                                raise ModuleNotFoundError(
                                    f"module {hook} not found\n{sys.path}"
                                ) from e
                        else:
                            raise e

                    self._pm.register(plugin)

            def __iter__(
                self: "Markata", description: str = "working..."
            ) -> Iterable["Markata.Post"]:
                articles: Iterable[Markata.Post] = track(
                    self.articles,
                    description=description,
                    transient=False,
                    console=self.console,
                )
                return articles

            def iter_articles(self: "Markata", description: str) -> Iterable[Markata.Post]:
                articles: Iterable[Markata.Post] = track(
                    self.articles,
                    description=description,
                    transient=True,
                    console=self.console,
                )
                return articles

            def teardown(self: "Markata") -> Markata:
                """give special access to the teardown lifecycle method"""
                self._pm.hook.teardown(markata=self)
                return self

            def run(self: "Markata", lifecycle: LifeCycle = None) -> Markata:
                if lifecycle is None:
                    lifecycle = max(LifeCycle._member_map_.values())

                if isinstance(lifecycle, str):
                    lifecycle = LifeCycle[lifecycle]

                stages_to_run = [
                    m
                    for m in LifeCycle._member_map_
                    if (LifeCycle[m] <= lifecycle) and (m not in self.stages_ran)
                ]

                if not stages_to_run:
                    self.console.log(f"{lifecycle.name} already ran")
                    return self

                self.console.log(f"running {stages_to_run}")
                for stage in stages_to_run:
                    self.console.log(f"{stage} running")
                    getattr(self._pm.hook, stage)(markata=self)
                    self.stages_ran.add(stage)
                    self.console.log(f"{stage} complete")

                with self.cache as cache:
                    hits, misses = cache.stats()

                if hits + misses > 0:
                    self.console.log(
                        f"lifetime cache hit rate {round(hits/ (hits + misses)*100, 2)}%",
                    )

                if misses > 0:
                    self.console.log(f"lifetime cache hits/misses {hits}/{misses}")

                hits -= self.init_cache_stats[0]
                misses -= self.init_cache_stats[1]

                if hits + misses > 0:
                    self.console.log(
                        f"run cache hit rate {round(hits/ (hits + misses)*100, 2)}%",
                    )

                if misses > 0:
                    self.console.log(f"run cache hits/misses {hits}/{misses}")

                return self

            def filter(self: "Markata", filter: str) -> list:
                def evalr(a: Markata.Post) -> Any:
                    try:
                        return eval(
                            filter,
                            {**a.to_dict(), "timedelta": timedelta, "post": a, "m": self},
                            {},
                        )
                    except AttributeError:
                        return eval(
                            filter,
                            {**a.to_dict(), "timedelta": timedelta, "post": a, "m": self},
                            {},
                        )

                return [a for a in self.articles if evalr(a)]

            def map(
                self: "Markata",
                func: str = "title",
                filter: str = "True",
                sort: str = "True",
                reverse: bool = True,
                *args: tuple,
                **kwargs: dict,
            ) -> list:
                import copy

                def try_sort(a: Any) -> int:
                    if "datetime" in sort.lower():
                        return a.get(sort, datetime.datetime(1970, 1, 1))

                    if "date" in sort.lower():
                        return a.get(sort, datetime.date(1970, 1, 1))

                    try:
                        value = eval(sort, a.to_dict(), {})
                    except NameError:
                        return -1
                    return value
                    try:
                        return int(value)
                    except TypeError:
                        try:
                            return int(value.timestamp())
                        except Exception:
                            try:
                                return int(
                                    datetime.datetime.combine(
                                        value,
                                        datetime.datetime.min.time(),
                                    ).timestamp(),
                                )
                            except Exception:
                                try:
                                    return sum([ord(c) for c in str(value)])
                                except Exception:
                                    return -1

                articles = copy.copy(self.articles)
                articles.sort(key=try_sort)
                if reverse:
                    articles.reverse()

                try:
                    posts = [
                        eval(
                            func,
                            {**a.to_dict(), "timedelta": timedelta, "post": a, "m": self},
                            {},
                        )
                        for a in articles
                        if eval(
                            filter,
                            {**a.to_dict(), "timedelta": timedelta, "post": a, "m": self},
                            {},
                        )
                    ]

                except NameError as e:
                    variable = str(e).split("'")[1]

                    missing_in_posts = self.map(
                        "path",
                        filter=f'"{variable}" not in post.keys()',
                    )
                    message = (
                        f"variable: '{variable}' is missing in {len(missing_in_posts)} posts"
                    )
                    if len(missing_in_posts) > 10:
                        message += (
                            f"\nfirst 10 paths to posts missing {variable}"
                            f"[{','.join([str(p) for p in missing_in_posts[:10]])}..."
                        )
                    else:
                        message += f"\npaths to posts missing {variable} {missing_in_posts}"

                    raise MissingFrontMatter(message)

                return posts

            def first(
                self: "Markata",
                filter: str = "True",
                sort: str = "True",
                reverse: bool = True,
                *args: tuple,
                **kwargs: dict,
            ) -> list:
                return self.map("post", filter, sort, reverse, *args, **kwargs)[0]

            def last(
                self: "Markata",
                filter: str = "True",
                sort: str = "True",
                reverse: bool = True,
                *args: tuple,
                **kwargs: dict,
            ) -> list:
                return self.map("post", filter, sort, reverse, *args, **kwargs)[-1]

            def one(
                self: "Markata",
                filter: str = "True",
                *args: tuple,
                **kwargs: dict,
            ) -> list:
                posts = self.map("post", filter, *args, **kwargs)
                if len(posts) > 1:
                    raise TooManyPosts(f"found {len(posts)} posts, expected 1. {posts}")
                if len(posts) == 0:
                    raise NoPosts
                return posts[0]
```


!! function <h2 id='load_ipython_extension' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>load_ipython_extension <em class='small'>function</em></h2>

???+ source "load_ipython_extension <em class='small'>source</em>"

```python

        def load_ipython_extension(ipython):
            ipython.user_ns["m"] = Markata()
            ipython.user_ns["markata"] = ipython.user_ns["m"]
            ipython.user_ns["markata"] = ipython.user_ns["m"]
```


!! method <h2 id='__init__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__init__ <em class='small'>method</em></h2>

???+ source "__init__ <em class='small'>source</em>"

```python

        def __init__(self: "Markata", console: Console = None, config=None) -> None:
                self.__version__ = __version__
                self.stages_ran = set()
                self.threded = False
                self._cache = None
                self._precache = None
                self.MARKATA_CACHE_DIR = Path(".") / ".markata.cache"
                self.MARKATA_CACHE_DIR.mkdir(exist_ok=True)
                self._pm = pluggy.PluginManager("markata")
                self._pm.add_hookspecs(hookspec.MarkataSpecs)
                if config is not None:
                    self.config = config
                with self.cache as cache:
                    self.init_cache_stats = cache.stats()
                self.registered_attrs = hookspec.registered_attrs
                self.post_models = []
                self.config_models = []
                if config is not None:
                    raw_hooks = config
                else:
                    raw_hooks = standard_config.load("markata")
                self.hooks_conf = HooksConfig.parse_obj(raw_hooks)
                try:
                    default_index = self.hooks_conf.hooks.index("default")
                    hooks = [
                        *self.hooks_conf.hooks[:default_index],
                        *DEFAULT_HOOKS,
                        *self.hooks_conf.hooks[default_index + 1 :],
                    ]
                    self.hooks_conf.hooks = [
                        hook for hook in hooks if hook not in self.hooks_conf.disabled_hooks
                    ]
                except ValueError:
                    # 'default' is not in hooks , do not replace with default_hooks
                    pass

                self._register_hooks()
                if console is not None:
                    self._console = console
                atexit.register(self.teardown)
                self.precache
```


!! method <h2 id='cache' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>cache <em class='small'>method</em></h2>

???+ source "cache <em class='small'>source</em>"

```python

        def cache(self: "Markata") -> Cache:
                # if self.threded:
                #     FanoutCache(self.MARKATA_CACHE_DIR, statistics=True)
                if self._cache is not None:
                    return self._cache
                self._cache = Cache(self.MARKATA_CACHE_DIR, statistics=True)

                return self._cache
```


!! method <h2 id='precache' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>precache <em class='small'>method</em></h2>

???+ source "precache <em class='small'>source</em>"

```python

        def precache(self: "Markata") -> None:
                if self._precache is None:
                    self.cache.expire()
                    self._precache = {k: self.cache.get(k) for k in self.cache.iterkeys()}
                return self._precache
```


!! method <h2 id='__getattr__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__getattr__ <em class='small'>method</em></h2>

???+ source "__getattr__ <em class='small'>source</em>"

```python

        def __getattr__(self: "Markata", item: str) -> Any:
                if item in self._pm.hook.__dict__:
                    # item is a hook, return a callable function
                    return lambda: self.run(item)

                if item in self.__dict__:
                    # item is an attribute, return it
                    return self.__getitem__(item)

                elif item in self.registered_attrs:
                    # item is created by a plugin, run it
                    stage_to_run_to = max(
                        [attr["lifecycle"] for attr in self.registered_attrs[item]],
                    ).name
                    self.console.log(
                        f"Running to [purple]{stage_to_run_to}[/] to retrieve [purple]{item}[/]"
                    )
                    self.run(stage_to_run_to)
                    return getattr(self, item)
                elif item == "precache":
                    return self._precache or {}
                else:
                    # Markata does not know what this is, raise
                    raise AttributeError(f"'Markata' object has no attribute '{item}'")
```


!! method <h2 id='__rich__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__rich__ <em class='small'>method</em></h2>

???+ source "__rich__ <em class='small'>source</em>"

```python

        def __rich__(self: "Markata") -> Table:
                grid = Table.grid()
                grid.add_column("label")
                grid.add_column("value")

                for label, value in self.describe().items():
                    grid.add_row(label, value)

                return grid
```


!! method <h2 id='bust_cache' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>bust_cache <em class='small'>method</em></h2>

???+ source "bust_cache <em class='small'>source</em>"

```python

        def bust_cache(self: "Markata") -> Markata:
                with self.cache as cache:
                    cache.clear()
                return self
```


!! method <h2 id='configure' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>configure <em class='small'>method</em></h2>

???+ source "configure <em class='small'>source</em>"

```python

        def configure(self) -> Markata:
                sys.path.append(os.getcwd())
                # self.config = {**DEFUALT_CONFIG, **standard_config.load("markata")}
                # if isinstance(self.config["glob_patterns"], str):
                #     self.config["glob_patterns"] = self.config["glob_patterns"].split(",")
                # elif isinstance(self.config["glob_patterns"], list):
                #     self.config["glob_patterns"] = list(self.config["glob_patterns"])
                # else:
                #     raise TypeError("glob_patterns must be list or str")
                # self.glob_patterns = self.config["glob_patterns"]

                # self.hooks = self.config["hooks"]

                # if "disabled_hooks" not in self.config:
                #     self.disabled_hooks = [""]
                # if isinstance(self.config["disabled_hooks"], str):
                #     self.disabled_hooks = self.config["disabled_hooks"].split(",")
                # if isinstance(self.config["disabled_hooks"], list):
                #     self.disabled_hooks = self.config["disabled_hooks"]

                # if not self.config.get("output_dir", "markout").endswith(
                #     self.config.get("path_prefix", "")
                # ):
                #     self.config["output_dir"] = (
                #         self.config.get("output_dir", "markout") +
                #         "/" +
                #         self.config.get("path_prefix", "").rstrip("/")
                #     )
                # if (
                #     len((output_split := self.config.get("output_dir", "markout").split("/"))) >
                #     1
                # ):
                #     if "path_prefix" not in self.config.keys():
                #         self.config["path_prefix"] = "/".join(output_split[1:]) + "/"
                # if not self.config.get("path_prefix", "").endswith("/"):
                #     self.config["path_prefix"] = self.config.get("path_prefix", "") + "/"

                # self.config["output_dir"] = self.config["output_dir"].lstrip("/")
                # self.config["path_prefix"] = self.config["path_prefix"].lstrip("/")

                try:
                    default_index = self.hooks_conf.hooks.index("default")
                    hooks = [
                        *self.hooks_conf.hooks[:default_index],
                        *DEFAULT_HOOKS,
                        *self.hooks_conf.hooks[default_index + 1 :],
                    ]
                    self.config.hooks = [
                        hook for hook in hooks if hook not in self.config.disabled_hooks
                    ]
                except ValueError:
                    # 'default' is not in hooks , do not replace with default_hooks
                    pass

                self._pm = pluggy.PluginManager("markata")
                self._pm.add_hookspecs(hookspec.MarkataSpecs)
                self._register_hooks()

                self._pm.hook.configure(markata=self)
                return self
```


!! method <h2 id='get_plugin_config' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get_plugin_config <em class='small'>method</em></h2>

???+ source "get_plugin_config <em class='small'>source</em>"

```python

        def get_plugin_config(self, path_or_name: str) -> Dict:
                key = Path(path_or_name).stem

                config = self.config.get(key, {})

                if not isinstance(config, dict):
                    raise TypeError("must use dict")
                if "cache_expire" not in config.keys():
                    config["cache_expire"] = self.config["default_cache_expire"]
                if "config_key" not in config.keys():
                    config["config_key"] = key
                return config
```


!! method <h2 id='get_config' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get_config <em class='small'>method</em></h2>

???+ source "get_config <em class='small'>source</em>"

```python

        def get_config(
                self,
                key: str,
                default: str = "",
                warn: bool = True,
                suggested: Optional[str] = None,
            ) -> Any:
                if key in self.config.keys():
                    return self.config[key]
                else:
                    if suggested is None:
                        suggested = textwrap.dedent(
                            f"""
                            [markata]
                            {key} = '{default}'
                            """
                        )
                    if warn:
                        logger.warning(
                            textwrap.dedent(
                                f"""
                                Warning {key} is not set in markata config, sitemap will
                                be missing root site_name
                                to resolve this open your markata.toml and add

                                {suggested}

                                """
                            ),
                        )
                return default
```


!! method <h2 id='make_hash' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>make_hash <em class='small'>method</em></h2>

???+ source "make_hash <em class='small'>source</em>"

```python

        def make_hash(self, *keys: str) -> str:
                import xxhash

                str_keys = [str(key) for key in keys]
                hash = xxhash.xxh64("".join(str_keys).encode("utf-8")).hexdigest()
                return hash
```


!! method <h2 id='content_dir_hash' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>content_dir_hash <em class='small'>method</em></h2>

???+ source "content_dir_hash <em class='small'>source</em>"

```python

        def content_dir_hash(self: "Markata") -> str:
                hashes = [
                    dirhash(dir)
                    for dir in self.content_directories
                    if dir.absolute() != Path(".").absolute()
                ]
                return self.make_hash(*hashes)
```


!! method <h2 id='console' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>console <em class='small'>method</em></h2>

???+ source "console <em class='small'>source</em>"

```python

        def console(self: "Markata") -> Console:
                try:
                    return self._console
                except AttributeError:
                    self._console = Console()
                    return self._console
```


!! method <h2 id='describe' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>describe <em class='small'>method</em></h2>

???+ source "describe <em class='small'>source</em>"

```python

        def describe(self: "Markata") -> dict[str, str]:
                return {"version": __version__}
```


!! method <h2 id='_to_dict' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_to_dict <em class='small'>method</em></h2>

???+ source "_to_dict <em class='small'>source</em>"

```python

        def _to_dict(self: "Markata") -> dict[str, Iterable]:
                return {"config": self.config, "articles": [a.to_dict() for a in self.articles]}
```


!! method <h2 id='to_dict' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>to_dict <em class='small'>method</em></h2>

???+ source "to_dict <em class='small'>source</em>"

```python

        def to_dict(self: "Markata") -> dict:
                return self._to_dict()
```


!! method <h2 id='to_json' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>to_json <em class='small'>method</em></h2>

???+ source "to_json <em class='small'>source</em>"

```python

        def to_json(self: "Markata") -> str:
                import json

                return json.dumps(self.to_dict(), indent=4, sort_keys=True, default=str)
```


!! method <h2 id='_register_hooks' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_register_hooks <em class='small'>method</em></h2>

???+ source "_register_hooks <em class='small'>source</em>"

```python

        def _register_hooks(self: "Markata") -> None:
                sys.path.append(os.getcwd())
                for hook in self.hooks_conf.hooks:
                    try:
                        # module style plugins
                        plugin = importlib.import_module(hook)
                    except ModuleNotFoundError as e:
                        # class style plugins
                        if "." in hook:
                            try:
                                mod = importlib.import_module(".".join(hook.split(".")[:-1]))
                                plugin = getattr(mod, hook.split(".")[-1])
                            except ModuleNotFoundError as e:
                                raise ModuleNotFoundError(
                                    f"module {hook} not found\n{sys.path}"
                                ) from e
                        else:
                            raise e

                    self._pm.register(plugin)
```


!! method <h2 id='__iter__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__iter__ <em class='small'>method</em></h2>

???+ source "__iter__ <em class='small'>source</em>"

```python

        def __iter__(
                self: "Markata", description: str = "working..."
            ) -> Iterable["Markata.Post"]:
                articles: Iterable[Markata.Post] = track(
                    self.articles,
                    description=description,
                    transient=False,
                    console=self.console,
                )
                return articles
```


!! method <h2 id='iter_articles' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>iter_articles <em class='small'>method</em></h2>

???+ source "iter_articles <em class='small'>source</em>"

```python

        def iter_articles(self: "Markata", description: str) -> Iterable[Markata.Post]:
                articles: Iterable[Markata.Post] = track(
                    self.articles,
                    description=description,
                    transient=True,
                    console=self.console,
                )
                return articles
```


!! method <h2 id='teardown' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>teardown <em class='small'>method</em></h2>
    give special access to the teardown lifecycle method
???+ source "teardown <em class='small'>source</em>"

```python

        def teardown(self: "Markata") -> Markata:
                """give special access to the teardown lifecycle method"""
                self._pm.hook.teardown(markata=self)
                return self
```


!! method <h2 id='run' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>run <em class='small'>method</em></h2>

???+ source "run <em class='small'>source</em>"

```python

        def run(self: "Markata", lifecycle: LifeCycle = None) -> Markata:
                if lifecycle is None:
                    lifecycle = max(LifeCycle._member_map_.values())

                if isinstance(lifecycle, str):
                    lifecycle = LifeCycle[lifecycle]

                stages_to_run = [
                    m
                    for m in LifeCycle._member_map_
                    if (LifeCycle[m] <= lifecycle) and (m not in self.stages_ran)
                ]

                if not stages_to_run:
                    self.console.log(f"{lifecycle.name} already ran")
                    return self

                self.console.log(f"running {stages_to_run}")
                for stage in stages_to_run:
                    self.console.log(f"{stage} running")
                    getattr(self._pm.hook, stage)(markata=self)
                    self.stages_ran.add(stage)
                    self.console.log(f"{stage} complete")

                with self.cache as cache:
                    hits, misses = cache.stats()

                if hits + misses > 0:
                    self.console.log(
                        f"lifetime cache hit rate {round(hits/ (hits + misses)*100, 2)}%",
                    )

                if misses > 0:
                    self.console.log(f"lifetime cache hits/misses {hits}/{misses}")

                hits -= self.init_cache_stats[0]
                misses -= self.init_cache_stats[1]

                if hits + misses > 0:
                    self.console.log(
                        f"run cache hit rate {round(hits/ (hits + misses)*100, 2)}%",
                    )

                if misses > 0:
                    self.console.log(f"run cache hits/misses {hits}/{misses}")

                return self
```


!! method <h2 id='filter' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>filter <em class='small'>method</em></h2>

???+ source "filter <em class='small'>source</em>"

```python

        def filter(self: "Markata", filter: str) -> list:
                def evalr(a: Markata.Post) -> Any:
                    try:
                        return eval(
                            filter,
                            {**a.to_dict(), "timedelta": timedelta, "post": a, "m": self},
                            {},
                        )
                    except AttributeError:
                        return eval(
                            filter,
                            {**a.to_dict(), "timedelta": timedelta, "post": a, "m": self},
                            {},
                        )

                return [a for a in self.articles if evalr(a)]
```


!! method <h2 id='map' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>map <em class='small'>method</em></h2>

???+ source "map <em class='small'>source</em>"

```python

        def map(
                self: "Markata",
                func: str = "title",
                filter: str = "True",
                sort: str = "True",
                reverse: bool = True,
                *args: tuple,
                **kwargs: dict,
            ) -> list:
                import copy

                def try_sort(a: Any) -> int:
                    if "datetime" in sort.lower():
                        return a.get(sort, datetime.datetime(1970, 1, 1))

                    if "date" in sort.lower():
                        return a.get(sort, datetime.date(1970, 1, 1))

                    try:
                        value = eval(sort, a.to_dict(), {})
                    except NameError:
                        return -1
                    return value
                    try:
                        return int(value)
                    except TypeError:
                        try:
                            return int(value.timestamp())
                        except Exception:
                            try:
                                return int(
                                    datetime.datetime.combine(
                                        value,
                                        datetime.datetime.min.time(),
                                    ).timestamp(),
                                )
                            except Exception:
                                try:
                                    return sum([ord(c) for c in str(value)])
                                except Exception:
                                    return -1

                articles = copy.copy(self.articles)
                articles.sort(key=try_sort)
                if reverse:
                    articles.reverse()

                try:
                    posts = [
                        eval(
                            func,
                            {**a.to_dict(), "timedelta": timedelta, "post": a, "m": self},
                            {},
                        )
                        for a in articles
                        if eval(
                            filter,
                            {**a.to_dict(), "timedelta": timedelta, "post": a, "m": self},
                            {},
                        )
                    ]

                except NameError as e:
                    variable = str(e).split("'")[1]

                    missing_in_posts = self.map(
                        "path",
                        filter=f'"{variable}" not in post.keys()',
                    )
                    message = (
                        f"variable: '{variable}' is missing in {len(missing_in_posts)} posts"
                    )
                    if len(missing_in_posts) > 10:
                        message += (
                            f"\nfirst 10 paths to posts missing {variable}"
                            f"[{','.join([str(p) for p in missing_in_posts[:10]])}..."
                        )
                    else:
                        message += f"\npaths to posts missing {variable} {missing_in_posts}"

                    raise MissingFrontMatter(message)

                return posts
```


!! method <h2 id='first' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>first <em class='small'>method</em></h2>

???+ source "first <em class='small'>source</em>"

```python

        def first(
                self: "Markata",
                filter: str = "True",
                sort: str = "True",
                reverse: bool = True,
                *args: tuple,
                **kwargs: dict,
            ) -> list:
                return self.map("post", filter, sort, reverse, *args, **kwargs)[0]
```


!! method <h2 id='last' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>last <em class='small'>method</em></h2>

???+ source "last <em class='small'>source</em>"

```python

        def last(
                self: "Markata",
                filter: str = "True",
                sort: str = "True",
                reverse: bool = True,
                *args: tuple,
                **kwargs: dict,
            ) -> list:
                return self.map("post", filter, sort, reverse, *args, **kwargs)[-1]
```


!! method <h2 id='one' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>one <em class='small'>method</em></h2>

???+ source "one <em class='small'>source</em>"

```python

        def one(
                self: "Markata",
                filter: str = "True",
                *args: tuple,
                **kwargs: dict,
            ) -> list:
                posts = self.map("post", filter, *args, **kwargs)
                if len(posts) > 1:
                    raise TooManyPosts(f"found {len(posts)} posts, expected 1. {posts}")
                if len(posts) == 0:
                    raise NoPosts
                return posts[0]
```


!! function <h2 id='evalr' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>evalr <em class='small'>function</em></h2>

???+ source "evalr <em class='small'>source</em>"

```python

        def evalr(a: Markata.Post) -> Any:
                    try:
                        return eval(
                            filter,
                            {**a.to_dict(), "timedelta": timedelta, "post": a, "m": self},
                            {},
                        )
                    except AttributeError:
                        return eval(
                            filter,
                            {**a.to_dict(), "timedelta": timedelta, "post": a, "m": self},
                            {},
                        )
```


!! function <h2 id='try_sort' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>try_sort <em class='small'>function</em></h2>

???+ source "try_sort <em class='small'>source</em>"

```python

        def try_sort(a: Any) -> int:
                    if "datetime" in sort.lower():
                        return a.get(sort, datetime.datetime(1970, 1, 1))

                    if "date" in sort.lower():
                        return a.get(sort, datetime.date(1970, 1, 1))

                    try:
                        value = eval(sort, a.to_dict(), {})
                    except NameError:
                        return -1
                    return value
                    try:
                        return int(value)
                    except TypeError:
                        try:
                            return int(value.timestamp())
                        except Exception:
                            try:
                                return int(
                                    datetime.datetime.combine(
                                        value,
                                        datetime.datetime.min.time(),
                                    ).timestamp(),
                                )
                            except Exception:
                                try:
                                    return sum([ord(c) for c in str(value)])
                                except Exception:
                                    return -1
```

