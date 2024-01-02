---
content: "Markata is a tool for handling directories of markdown.\n\n\n!! class <h2
  id='HooksConfig' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>HooksConfig
  <em class='small'>class</em></h2>\n\n???+ source \"HooksConfig <em class='small'>source</em>\"\n\n```python\n\n
  \       class HooksConfig(pydantic.BaseModel):\n            hooks: list = [\"default\"]\n
  \           disabled_hooks: list = []\n```\n\n\n!! class <h2 id='Markata' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>Markata <em class='small'>class</em></h2>\n\n???+
  source \"Markata <em class='small'>source</em>\"\n\n```python\n\n        class Markata:\n
  \           def __init__(self: \"Markata\", console: Console = None, config=None)
  -> None:\n                self.stages_ran = set()\n                self.threded
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
  for attr in self.registered_attrs[item]],\n                    ).name\n                    self.run(stage_to_run_to)\n
  \                   return getattr(self, item)\n                else:\n                    #
  Markata does not know what this is, raise\n                    raise AttributeError(f\"'Markata'
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
  default\n\n            def make_hash(self, *keys: str) -> str:\n                str_keys
  = [str(key) for key in keys]\n                return hashlib.md5(\"\".join(str_keys).encode(\"utf-8\")).hexdigest()\n\n
  \           @property\n            def content_dir_hash(self: \"Markata\") -> str:\n
  \               hashes = [\n                    dirhash(dir)\n                    for
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
  \                   transient=True,\n                    console=self.console,\n
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
  posts\n```\n\n\n!! function <h2 id='load_ipython_extension' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>load_ipython_extension <em class='small'>function</em></h2>\n\n???+
  source \"load_ipython_extension <em class='small'>source</em>\"\n\n```python\n\n
  \       def load_ipython_extension(ipython):\n            ipython.user_ns[\"m\"]
  = Markata()\n            ipython.user_ns[\"markata\"] = ipython.user_ns[\"m\"]\n
  \           ipython.user_ns[\"markata\"] = ipython.user_ns[\"m\"]\n```\n\n\n!! method
  <h2 id='__init__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__init__
  <em class='small'>method</em></h2>\n\n???+ source \"__init__ <em class='small'>source</em>\"\n\n```python\n\n
  \       def __init__(self: \"Markata\", console: Console = None, config=None) ->
  None:\n                self.stages_ran = set()\n                self.threded = False\n
  \               self._cache = None\n                self._precache = None\n                self.MARKATA_CACHE_DIR
  = Path(\".\") / \".markata.cache\"\n                self.MARKATA_CACHE_DIR.mkdir(exist_ok=True)\n
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
  for attr in self.registered_attrs[item]],\n                    ).name\n                    self.run(stage_to_run_to)\n
  \                   return getattr(self, item)\n                else:\n                    #
  Markata does not know what this is, raise\n                    raise AttributeError(f\"'Markata'
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
  -> str:\n                str_keys = [str(key) for key in keys]\n                return
  hashlib.md5(\"\".join(str_keys).encode(\"utf-8\")).hexdigest()\n```\n\n\n!! method
  <h2 id='content_dir_hash' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>content_dir_hash <em class='small'>method</em></h2>\n\n???+ source \"content_dir_hash
  <em class='small'>source</em>\"\n\n```python\n\n        def content_dir_hash(self:
  \"Markata\") -> str:\n                hashes = [\n                    dirhash(dir)\n
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
  \                   transient=True,\n                    console=self.console,\n
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
  posts\n```\n\n\n!! function <h2 id='evalr' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>evalr <em class='small'>function</em></h2>\n\n???+ source \"evalr <em class='small'>source</em>\"\n\n```python\n\n
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
  \                                   return -1\n```\n"
date: 0001-01-01
description: 'Markata is a tool for handling directories of markdown. ! ???+ source  !
  ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
  source '
html: "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n<title>__Init__.Py</title>\n<meta
  charset=\"utf-8\"/>\n<meta content=\"width=device-width, initial-scale=1\" name=\"viewport\"/>\n<meta
  content=\"Markata is a tool for handling directories of markdown. ! ???+ source
  \ ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
  source \" name=\"description\"/>\n<link href=\"/static/favicon.ico\" rel=\"icon\"
  type=\"image/png\"/>\n<script>\n        function setTheme(theme) {\n            document.documentElement.setAttribute(\"data-theme\",
  theme);\n        }\n\n        function detectColorSchemeOnLoad() {\n            //local
  storage is used to override OS theme settings\n            if (localStorage.getItem(\"theme\"))
  {\n                if (localStorage.getItem(\"theme\") == \"dark\") {\n                    setTheme(\"dark\");\n
  \               } else if (localStorage.getItem(\"theme\") == \"light\") {\n                    setTheme(\"light\");\n
  \               }\n            } else if (!window.matchMedia) {\n                //matchMedia
  method not supported\n                setTheme(\"light\");\n                return
  false;\n            } else if (window.matchMedia(\"(prefers-color-scheme: dark)\").matches)
  {\n                //OS theme setting detected as dark\n                setTheme(\"dark\");\n
  \           } else {\n                setTheme(\"light\");\n            }\n        }\n
  \       detectColorSchemeOnLoad();\n        document.addEventListener(\n            \"DOMContentLoaded\",\n
  \           function () {\n                //identify the toggle switch HTML element\n
  \               const toggleSwitch = document.querySelector(\n                    '#theme-switch
  input[type=\"checkbox\"]',\n                );\n\n                //function that
  changes the theme, and sets a localStorage variable to track the theme between page
  loads\n                function switchTheme(e) {\n                    if (e.target.checked)
  {\n                        localStorage.setItem(\"theme\", \"dark\");\n                        document.documentElement.setAttribute(\"data-theme\",
  \"dark\");\n                        toggleSwitch.checked = true;\n                    }
  else {\n                        localStorage.setItem(\"theme\", \"light\");\n                        document.documentElement.setAttribute(\"data-theme\",
  \"light\");\n                        toggleSwitch.checked = false;\n                    }\n
  \               }\n\n                //listener for changing themes\n                toggleSwitch.addEventListener(\"change\",
  switchTheme, false);\n\n                //pre-check the dark-theme checkbox if dark-theme
  is set\n                if (document.documentElement.getAttribute(\"data-theme\")
  == \"dark\") {\n                    toggleSwitch.checked = true;\n                }\n
  \           },\n            false,\n        );\n    </script>\n<style>\n      :root
  {\n        --color-bg: #1f2022;\n        --color-bg-2: ;\n        --color-bg-code:
  #1f2022;\n        --color-text: #eefbfe;\n        --color-link: #fb30c4; \n        --color-accent:
  #e1bd00c9;\n        --overlay-brightness: .85;\n        --body-width: 800px;\n      }\n
  \     [data-theme=\"dark\"] {\n        --color-bg: #1f2022;\n        --color-bg-2:
  ;\n        --color-bg-code: #1f2022;\n        --color-text: #eefbfe;\n        --color-link:
  #fb30c4; \n        --color-accent: #e1bd00c9;\n        --overlay-brightness: .85;\n
  \       --body-width: 800px;\n      }\n      [data-theme=\"light\"] {\n        --color-bg:
  #eefbfe;\n        --color-bg-2: ;\n        --color-bg-code: #eefbfe;\n        --color-text:
  #1f2022;\n        --color-link: #fb30c4; \n        --color-accent: #ffeb00;\n        --overlay-brightness:
  .95;\n      }\n\n        html {\n            font-family: \"Space Mono\", monospace;\n
  \           background: var(--color-bg);\n            color: var(--color-text);\n
  \       }\n\n        a {\n            color: var(--color-link);\n        }\n\n        main
  a {\n            max-width: 100%;\n        }\n\n        .heading-permalink {\n            font-size:
  .7em;\n        }\n\n        body {\n            max-width: var(--body-width);\n
  \           margin: 5rem auto;\n            padding: 0 .5rem;\n            font-size:
  1rem;\n            line-height: 1.56;\n        }\n\n        blockquote {\n            background:
  var(--color-bg);\n            filter: brightness(var(--overlay-brightness));\n            border-left:
  4px solid var(--color-accent);\n            border-radius: 4px;\n            box-shadow:\n
  \               -0.8rem 0rem 1rem -1rem #f1fa8c,\n                0.2rem 0rem 1rem
  rgb(0, 0, 0, .4);\n            padding-left: 1rem;\n            margin: 1rem;\n
  \       }\n\n        li.post {\n            list-style-type: None;\n            padding:
  .2rem 0;\n        }\n\n        pre.wrapper {\n            padding: 0;\n            box-shadow:
  0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n            display: flex;\n            flex-direction:
  column;\n            position: relative;\n            margin: 2rem;\n        }\n\n
  \       pre {\n            margin: 0;\n            padding: 1rem;\n            min-width:
  -webkit-fill-available;\n            max-width: fit-content;\n            overflow-x:
  auto;\n        }\n\n        pre .filepath {\n            margin: 0;\n            padding-left:
  1rem;\n            border-radius: 4px 4px 0 0;\n            background: black;\n
  \           display: flex;\n            justify-content: space-between;\n            align-items:
  center;\n        }\n\n        pre .filepath p {\n            margin: 0\n        }\n\n
  \       pre .filepath .right {\n            display: flex;\n            gap: .2rem;\n
  \           align-items: center;\n        }\n\n        pre::-webkit-scrollbar {\n
  \           height: 4px;\n            background-color: transparent;\n        }\n\n
  \       pre::-webkit-scrollbar-thumb {\n            background-color: #d3d3d32e;\n
  \           border-radius: 2px;\n        }\n\n        pre::-webkit-scrollbar-track
  {\n            background-color: transparent;\n        }\n\n        .copy-wrapper
  {\n            background: none;\n            position: absolute;\n            width:
  100%;\n            z-index: 100;\n            display: flex;\n            justify-content:
  flex-end;\n        }\n\n        button.copy {\n            z-index: 100;\n            background:
  none;\n            fill: #ffffff45;\n            border: none;\n            width:
  32px;\n            align-self: flex-end;\n            top: 0;\n            right:
  0;\n            margin: 0.5rem 0.2rem;\n\n        }\n\n        button.copy:hover
  {\n            fill: white\n        }\n\n        a.help {\n            fill: #ffffff45;\n
  \       }\n\n        a.help:hover {\n            fill: white;\n        }\n\n        a.help
  svg {\n            height: 24px;\n            width: 24px;\n        }\n\n        .highlight
  {\n            background: var(--color-bg-code);\n            color: var(--color-text);\n
  \           filter: brightness(var(--overlay-brightness));\n            border-radius:
  0 0 4px 4px;\n        }\n\n        .highlight .c {\n            color: #8b8b8b\n
  \       }\n\n        /* Comment */\n        .highlight .err {\n            color:
  #960050;\n            background-color: #1e0010\n        }\n\n        /* Error */\n
  \       .highlight .k {\n            color: #c678dd\n        }\n\n        /* Keyword
  */\n        .highlight .l {\n            color: #ae81ff\n        }\n\n        /*
  Literal */\n        .highlight .n {\n            color: #abb2bf\n        }\n\n        /*
  Name */\n        .highlight .o {\n            color: #c678dd\n        }\n\n        /*
  Operator */\n        .highlight .p {\n            color: #abb2bf\n        }\n\n
  \       /* Punctuation */\n        .highlight .ch {\n            color: #8b8b8b\n
  \       }\n\n        /* Comment.Hashbang */\n        .highlight .cm {\n            color:
  #8b8b8b\n        }\n\n        /* Comment.Multiline */\n        .highlight .cp {\n
  \           color: #8b8b8b\n        }\n\n        /* Comment.Preproc */\n        .highlight
  .cpf {\n            color: #8b8b8b\n        }\n\n        /* Comment.PreprocFile
  */\n        .highlight .c1 {\n            color: #8b8b8b\n        }\n\n        /*
  Comment.Single */\n        .highlight .cs {\n            color: #8b8b8b\n        }\n\n
  \       /* Comment.Special */\n        .highlight .gd {\n            color: #c678dd\n
  \       }\n\n        /* Generic.Deleted */\n        .highlight .ge {\n            font-style:
  italic\n        }\n\n        /* Generic.Emph */\n        .highlight .gi {\n            color:
  #a6e22e\n        }\n\n        /* Generic.Inserted */\n        .highlight .gs {\n
  \           font-weight: bold\n        }\n\n        /* Generic.Strong */\n        .highlight
  .gu {\n            color: #8b8b8b\n        }\n\n        /* Generic.Subheading */\n
  \       .highlight .kc {\n            color: #c678dd\n        }\n\n        /* Keyword.Constant
  */\n        .highlight .kd {\n            color: #c678dd\n        }\n\n        /*
  Keyword.Declaration */\n        .highlight .kn {\n            color: #c678dd\n        }\n\n
  \       /* Keyword.Namespace */\n        .highlight .kp {\n            color: #c678dd\n
  \       }\n\n        /* Keyword.Pseudo */\n        .highlight .kr {\n            color:
  #c678dd\n        }\n\n        /* Keyword.Reserved */\n        .highlight .kt {\n
  \           color: #c678dd\n        }\n\n        /* Keyword.Type */\n        .highlight
  .ld {\n            color: #e6db74\n        }\n\n        /* Literal.Date */\n        .highlight
  .m {\n            color: #ae81ff\n        }\n\n        /* Literal.Number */\n        .highlight
  .s {\n            color: #e6db74\n        }\n\n        /* Literal.String */\n        .highlight
  .na {\n            color: #a6e22e\n        }\n\n        /* Name.Attribute */\n        .highlight
  .nb {\n            color: #98c379\n        }\n\n        /* Name.Builtin */\n        .highlight
  .nc {\n            color: #abb2bf\n        }\n\n        /* Name.Class */\n        .highlight
  .no {\n            color: #c678dd\n        }\n\n        /* Name.Constant */\n        .highlight
  .nd {\n            color: #abb2bf\n        }\n\n        /* Name.Decorator */\n        .highlight
  .ni {\n            color: #abb2bf\n        }\n\n        /* Name.Entity */\n        .highlight
  .ne {\n            color: #a6e22e\n        }\n\n        /* Name.Exception */\n        .highlight
  .nf {\n            color: #61afef\n        }\n\n        /* Name.Function */\n        .highlight
  .nl {\n            color: #abb2bf\n        }\n\n        /* Name.Label */\n        .highlight
  .nn {\n            color: #abb2bf\n        }\n\n        /* Name.Namespace */\n        .highlight
  .nx {\n            color: #a6e22e\n        }\n\n        /* Name.Other */\n        .highlight
  .py {\n            color: #abb2bf\n        }\n\n        /* Name.Property */\n        .highlight
  .nt {\n            color: #c678dd\n        }\n\n        /* Name.Tag */\n        .highlight
  .nv {\n            color: #abb2bf\n        }\n\n        /* Name.Variable */\n        .highlight
  .ow {\n            color: #c678dd\n        }\n\n        /* Operator.Word */\n        .highlight
  .w {\n            color: #abb2bf\n        }\n\n        /* Text.Whitespace */\n        .highlight
  .mb {\n            color: #ae81ff\n        }\n\n        /* Literal.Number.Bin */\n
  \       .highlight .mf {\n            color: #ae81ff\n        }\n\n        /* Literal.Number.Float
  */\n        .highlight .mh {\n            color: #ae81ff\n        }\n\n        /*
  Literal.Number.Hex */\n        .highlight .mi {\n            color: #ae81ff\n        }\n\n
  \       /* Literal.Number.Integer */\n        .highlight .mo {\n            color:
  #ae81ff\n        }\n\n        /* Literal.Number.Oct */\n        .highlight .sa {\n
  \           color: #e6db74\n        }\n\n        /* Literal.String.Affix */\n        .highlight
  .sb {\n            color: #e6db74\n        }\n\n        /* Literal.String.Backtick
  */\n        .highlight .sc {\n            color: #e6db74\n        }\n\n        /*
  Literal.String.Char */\n        .highlight .dl {\n            color: #e6db74\n        }\n\n
  \       /* Literal.String.Delimiter */\n        .highlight .sd {\n            color:
  #98c379\n        }\n\n        /* Literal.String.Doc */\n        .highlight .s2 {\n
  \           color: #98c379\n        }\n\n        /* Literal.String.Double */\n        .highlight
  .se {\n            color: #ae81ff\n        }\n\n        /* Literal.String.Escape
  */\n        .highlight .sh {\n            color: #e6db74\n        }\n\n        /*
  Literal.String.Heredoc */\n        .highlight .si {\n            color: #e6db74\n
  \       }\n\n        /* Literal.String.Interpol */\n        .highlight .sx {\n            color:
  #e6db74\n        }\n\n        /* Literal.String.Other */\n        .highlight .sr
  {\n            color: #e6db74\n        }\n\n        /* Literal.String.Regex */\n
  \       .highlight .s1 {\n            color: #e6db74\n        }\n\n        /* Literal.String.Single
  */\n        .highlight .ss {\n            color: #e6db74\n        }\n\n        /*
  Literal.String.Symbol */\n        .highlight .bp {\n            color: #abb2bf\n
  \       }\n\n        /* Name.Builtin.Pseudo */\n        .highlight .fm {\n            color:
  #61afef\n        }\n\n        /* Name.Function.Magic */\n        .highlight .vc
  {\n            color: #abb2bf\n        }\n\n        /* Name.Variable.Class */\n
  \       .highlight .vg {\n            color: #abb2bf\n        }\n\n        /* Name.Variable.Global
  */\n        .highlight .vi {\n            color: #abb2bf\n        }\n\n        /*
  Name.Variable.Instance */\n        .highlight .vm {\n            color: #abb2bf\n
  \       }\n\n        /* Name.Variable.Magic */\n        .highlight .il {\n            color:
  #ae81ff\n        }\n\n        /* Literal.Number.Integer.Long */\n\n        /* Tab
  style starts here */\n        .tabbed-set {\n            position: relative;\n            display:
  flex;\n            flex-wrap: wrap;\n            margin: 1em 0;\n            border-radius:
  0.1rem;\n        }\n\n        .tabbed-set>input {\n            display: none;\n
  \       }\n\n        .tabbed-set label {\n            width: auto;\n            padding:
  0.9375em 1.25em 0.78125em;\n            font-weight: 700;\n            font-size:
  0.84em;\n            white-space: nowrap;\n            border-bottom: 0.15rem solid
  transparent;\n            border-top-left-radius: 0.1rem;\n            border-top-right-radius:
  0.1rem;\n            cursor: pointer;\n            transition: background-color
  250ms, color 250ms;\n        }\n\n        .tabbed-set .tabbed-content {\n            width:
  100%;\n            display: none;\n            box-shadow: 0 -.05rem #ddd;\n        }\n\n
  \       .tabbed-set input {\n            position: absolute;\n            opacity:
  0;\n        }\n\n        /* fonts */\n        h1 {\n            font-weight: 700;\n
  \       }\n\n        h1#title a {\n            font-size: 16px;\n        }\n\n        h1,\n
  \       h2,\n        h3,\n        h4,\n        h5,\n        h6 {\n            margin-top:
  3rem;\n        }\n\n        h1 {\n            font-size: 2.5em;\n            margin-top:
  5rem;\n        }\n\n        h2 {\n            font-size: 1.63rem;\n            margin-top:
  5rem;\n        }\n\n\n\n        p {\n            font-size: 21px;\n            font-style:
  normal;\n            font-variant: normal;\n            font-weight: 400;\n            line-height:
  1.5;\n        }\n\n        @media only screen and (max-width: 700px) {\n            p
  {\n                font-size: 18px;\n            }\n        }\n\n        @media
  only screen and (max-width: 600px) {\n            p {\n                font-size:
  16px;\n            }\n        }\n\n        @media only screen and (max-width: 500px)
  {\n            p {\n                font-size: 14px;\n            }\n        }\n\n
  \       @media only screen and (max-width: 400px) {\n            p {\n                font-size:
  12px;\n            }\n        }\n\n\n        pre {\n            font-style: normal;\n
  \           font-variant: normal;\n            font-weight: 400;\n            line-height:
  18.5714px;\n            */\n        }\n\n        a {\n            font-weight: 600;\n
  \           text-decoration-color: var(--color-accent);\n            color: var(--color-link);\n
  \           padding: .3rem .5rem;\n            display: inline-block;\n        }\n\n
  \       .admonition,\n        details {\n            box-shadow: 0.2rem 0rem 1rem
  rgb(0, 0, 0, .4);\n            margin: 5rem 0;\n            border: 1px solid transparent;\n
  \           border-radius: 4px;\n            text-align: left;\n            padding:
  0;\n            border: 0;\n\n        }\n\n        .admonition {\n            padding-bottom:
  1rem;\n        }\n\n        details[open] {\n            padding-bottom: .5rem;\n
  \       }\n\n        .admonition p {\n            padding: .2rem .6rem;\n        }\n\n
  \       .admonition-title,\n        .details-title,\n        summary {\n            background:
  var(--color-bg-2);\n            padding: 0;\n            margin: 0;\n            position:
  sticky;\n            top: 0;\n            z-index: 10;\n        }\n\n        summary:hover
  {\n            cursor: pointer;\n        }\n\n        summary.admonition-title,\n
  \       summary.details-title {\n            padding: .5rem;\n            padding-left:
  1rem;\n        }\n\n        .note {\n            border-left: 4px solid #f1fa8c;\n
  \           box-shadow:\n                -0.8rem 0rem 1rem -1rem #f1fa8c,\n                0.2rem
  0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n        .note>.admonition-title {\n            border-bottom:
  1px solid #3c3d2d;\n        }\n\n        .abstract {\n            border-left: 4px
  solid #8be9fd;\n            box-shadow:\n                -0.8rem 0rem 1rem -1rem
  #8be9fd,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n        .abstract>.admonition-title
  {\n            border-bottom: 1px solid #2c3a3f;\n        }\n\n        .info {\n
  \           border-left: 4px solid;\n            box-shadow:\n                -0.8rem
  0rem 1rem -1rem #8bb0fd,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n
  \       .info>.admonition-title {\n            border-bottom: 1px solid #2c313f;\n
  \       }\n\n        .tip {\n            border-left: 4px solid #008080;\n            box-shadow:\n
  \               -0.8rem 0rem 1rem -1rem #008080,\n                0.2rem 0rem 1rem
  rgb(0, 0, 0, .4);\n        }\n\n        .tip>.admonition-title {\n            border-bottom:
  1px solid #1b2a2b;\n        }\n\n        .success {\n            border-left: 4px
  solid #50fa7b;\n            box-shadow:\n                -0.8rem 0rem 1rem -1rem
  #50fa7b,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n        .success>.admonition-title
  {\n            border-bottom: 1px solid #263e2b;\n        }\n\n        .question
  {\n            border-left: 4px solid #a7fcbd;\n            box-shadow:\n                -0.8rem
  0rem 1rem -1rem #a7fcbd,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n
  \       .question>.admonition-title {\n            border-bottom: 1px solid #303e35;\n
  \       }\n\n        .warning {\n            border-left: 4px solid #ffb86c;\n            box-shadow:\n
  \               -0.8rem 0rem 1rem -1rem #ffb86c,\n                0.2rem 0rem 1rem
  rgb(0, 0, 0, .4);\n        }\n\n        .warning>.admonition-title {\n            border-bottom:
  1px solid #3f3328;\n        }\n\n        .failure {\n            border-left: 4px
  solid #b23b3b;\n            box-shadow:\n                -0.8rem 0rem 1rem -1rem
  #b23b3b,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n        .failure>.admonition-title
  {\n            border-bottom: 1px solid #34201f;\n        }\n\n        .danger {\n
  \           border-left: 4px solid #ff5555;\n            box-shadow:\n                -0.8rem
  0rem 1rem -1rem #ff5555,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n
  \       .danger>.admonition-title {\n            border-bottom: 1px solid #402523;\n
  \       }\n\n        .bug {\n            border-left: 4px solid #b2548a;\n            box-shadow:\n
  \               -0.8rem 0rem 1rem -1rem #b2548a,\n                0.2rem 0rem 1rem
  rgb(0, 0, 0, .4);\n        }\n\n        .bug>.admonition-title {\n            border-bottom:
  1px solid #32232c;\n        }\n\n        .example {\n            border-left: 4px
  solid #bd93f9;\n            box-shadow:\n                -0.8rem 0rem 1rem -1rem
  #bd93f9,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n        .example>.admonition-title
  {\n            border-bottom: 1px solid #332d3e;\n        }\n\n        .source {\n
  \           border-left: 4px solid #bd93f9;\n            box-shadow:\n                -0.8rem
  0rem 1rem -1rem #bd93f9,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n
  \       .source>.admonition-title {\n            border-bottom: 1px solid #332d3e;\n
  \       }\n\n        .quote {\n            border-left: 4px solid #999;\n            box-shadow:\n
  \               -0.8rem 0rem 1rem -1rem #999,\n                0.2rem 0rem 1rem
  rgb(0, 0, 0, .4);\n        }\n\n        .quote>.admonition-title {\n            border-bottom:
  1px solid #2d2e2f;\n        }\n\n        table {\n            margin: 1rem 0;\n
  \           border-collapse: collapse;\n            border-spacing: 0;\n            display:
  block;\n            max-width: -moz-fit-content;\n            max-width: fit-content;\n
  \           overflow-x: auto;\n            white-space: nowrap;\n        }\n\n        table
  thead th {\n            border: solid 1px var(--color-text);\n            padding:
  10px;\n            text-align: left;\n        }\n\n        table tbody td {\n            border:
  solid 1px var(--color-text);\n            padding: 10px;\n        }\n\n        .theme-switch
  {\n            z-index: 10;\n            display: inline-block;\n            height:
  34px;\n            position: relative;\n            width: 60px;\n\n            display:
  flex;\n            justify-content: flex-end;\n            margin-right: 1rem;\n
  \           margin-left: auto;\n            position: fixed;\n            right:
  1rem;\n            top: 1rem;\n        }\n\n        .theme-switch input {\n            display:
  none;\n\n        }\n\n        .slider {\n            background-color: #ccc;\n            bottom:
  0;\n            cursor: pointer;\n            left: 0;\n            position: absolute;\n
  \           right: 0;\n            top: 0;\n            transition: .4s;\n        }\n\n
  \       .slider:before {\n            background-color: #fff;\n            bottom:
  4px;\n            content: \"\";\n            height: 26px;\n            left: 4px;\n
  \           position: absolute;\n            transition: .4s;\n            width:
  26px;\n        }\n\n        input:checked+.slider {\n            background-color:
  #343434;\n        }\n\n        input:checked+.slider:before {\n            background-color:
  #848484;\n        }\n\n        input:checked+.slider:before {\n            transform:
  translateX(26px);\n        }\n\n        .slider.round {\n            border-radius:
  34px;\n        }\n\n        .slider.round:before {\n            border-radius: 50%;\n
  \       }\n\n        main p img {\n            width: 100%;\n            width:
  -moz-available;\n            width: -webkit-fill-available;\n            width:
  fill-available;\n        }\n\n        details>* {\n            margin: 1rem;\n        }\n\n
  \       .admonition>* {\n            margin: 1rem;\n        }\n\n        p.admonition-title,\n
  \       summary {\n            margin: 0;\n            padding-left: 1.2rem;\n        }\n\n
  \       .small {\n            font-size: .9rem;\n            color: #888;\n        }\n\n
  \       admonition+admonition {\n            margin-top: 20rem;\n        }\n\n        ::-webkit-scrollbar
  {\n            height: 12px;\n            background-color: transparent;\n        }\n\n
  \       ::-webkit-scrollbar-thumb {\n            background-color: #d3d3d32e;\n
  \           border-radius: 6px;\n        }\n\n        ::-webkit-scrollbar-track
  {\n            background-color: transparent;\n        }\n    </style>\n<script>\n
  \       if (\"serviceWorker\" in navigator) {\n            navigator.serviceWorker.register(\"/service-worker.js\");\n
  \           navigator.serviceWorker.addEventListener(\"controllerchange\", () =>
  {\n                console.log(\"new worker\");\n                window.location.reload();\n
  \           });\n        }\n    </script>\n<meta content=\"waylon@waylonwalker.com\"
  name=\"og:author_email\"/>\n<meta content=\"waylon@waylonwalker.com\" name=\"og:author_email\"/>\n<meta
  content=\"Waylon Walker\" name=\"og:author\" property=\"og:author\"/><meta content=\"waylon@waylonwalaker.com\"
  name=\"og:author_email\" property=\"og:author_email\"/><meta content=\"website\"
  name=\"og:type\" property=\"og:type\"/><meta content=\"Markata is a tool for handling
  directories of markdown. ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source
  \ ! ???+ source  ! ???+ source  ! ???+ source \" name=\"description\" property=\"description\"/><meta
  content=\"Markata is a tool for handling directories of markdown. ! ???+ source
  \ ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
  source \" name=\"og:description\" property=\"og:description\"/><meta content=\"Markata
  is a tool for handling directories of markdown. ! ???+ source  ! ???+ source  !
  ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source \" name=\"twitter:description\"
  property=\"twitter:description\"/><meta content=\"__Init__.Py | Markata\" name=\"og:title\"
  property=\"og:title\"/><meta content=\"__Init__.Py | Markata\" name=\"twitter:title\"
  property=\"twitter:title\"/><meta content=\"https://markata.dev//markata/init-og.png\"
  name=\"og:image\" property=\"og:image\"/><meta content=\"https://markata.dev//markata/init-og.png\"
  name=\"twitter:image\" property=\"twitter:image\"/><meta content=\"1600\" name=\"og:image:width\"
  property=\"og:image:width\"/><meta content=\"900\" name=\"og:image:width\" property=\"og:image:width\"/><meta
  content=\"summary_large_image\" name=\"twitter:card\" property=\"twitter:card\"/><meta
  content=\"Markata\" name=\"og:site_name\" property=\"og:site_name\"/><meta content=\"@_waylonwalker\"
  name=\"twitter:creator\" property=\"twitter:creator\"/><meta content=\"__Init__.Py\"
  name=\"title\" property=\"title\"/><meta content=\"markata 0.8.0.dev6\" name=\"generator\"
  property=\"generator\"/><link href=\"https://markata.dev//markata/init/\" rel=\"canonical\"/><meta
  content=\"https://markata.dev//markata/init/\" name=\"og:url\" property=\"og:url\"/></head>\n<body><nav>\n<a
  href=\"/\">markata</a>\n<a href=\"https://github.com/WaylonWalker/markata\">GitHub</a>\n</nav>\n<div>\n<label
  class=\"theme-switch\" for=\"checkbox-theme\" id=\"theme-switch\" title=\"light/dark
  mode toggle\">\n<input id=\"checkbox-theme\" type=\"checkbox\"/>\n<div class=\"slider
  round\"></div>\n</label>\n</div>\n<section class=\"title\">\n<h1 id=\"title\">\n
  \           __Init__.Py \n            \n        </h1>\n</section>\n<main><p>Markata
  is a tool for handling directories of markdown.</p>\n<p>!! class </p><h2 class=\"admonition-title\"
  id=\"HooksConfig\" style=\"margin:0;padding:.5rem 1rem;\">HooksConfig <em class=\"small\">class</em></h2>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"HooksConfig
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
  <span class=\"nc\">HooksConfig</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
  class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
  \           <span class=\"n\">hooks</span><span class=\"p\">:</span> <span class=\"nb\">list</span>
  <span class=\"o\">=</span> <span class=\"p\">[</span><span class=\"s2\">\"default\"</span><span
  class=\"p\">]</span>\n            <span class=\"n\">disabled_hooks</span><span class=\"p\">:</span>
  <span class=\"nb\">list</span> <span class=\"o\">=</span> <span class=\"p\">[]</span>\n</pre></div>\n\n</pre>\n<p>!!
  class </p><h2 class=\"admonition-title\" id=\"Markata\" style=\"margin:0;padding:.5rem
  1rem;\">Markata <em class=\"small\">class</em></h2>\n<div class=\"admonition source
  is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"Markata <em
  class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
  <span class=\"nc\">Markata</span><span class=\"p\">:</span>\n            <span class=\"k\">def</span>
  <span class=\"fm\">__init__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">,</span>
  <span class=\"n\">console</span><span class=\"p\">:</span> <span class=\"n\">Console</span>
  <span class=\"o\">=</span> <span class=\"kc\">None</span><span class=\"p\">,</span>
  <span class=\"n\">config</span><span class=\"o\">=</span><span class=\"kc\">None</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">stages_ran</span> <span class=\"o\">=</span> <span class=\"nb\">set</span><span
  class=\"p\">()</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">threded</span> <span class=\"o\">=</span> <span class=\"kc\">False</span>\n
  \               <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_cache</span>
  <span class=\"o\">=</span> <span class=\"kc\">None</span>\n                <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_precache</span>
  <span class=\"o\">=</span> <span class=\"kc\">None</span>\n                <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">MARKATA_CACHE_DIR</span>
  <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
  class=\"s2\">\".\"</span><span class=\"p\">)</span> <span class=\"o\">/</span> <span
  class=\"s2\">\".markata.cache\"</span>\n                <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">MARKATA_CACHE_DIR</span><span class=\"o\">.</span><span
  class=\"n\">mkdir</span><span class=\"p\">(</span><span class=\"n\">exist_ok</span><span
  class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n                <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_pm</span> <span
  class=\"o\">=</span> <span class=\"n\">pluggy</span><span class=\"o\">.</span><span
  class=\"n\">PluginManager</span><span class=\"p\">(</span><span class=\"s2\">\"markata\"</span><span
  class=\"p\">)</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">_pm</span><span class=\"o\">.</span><span class=\"n\">add_hookspecs</span><span
  class=\"p\">(</span><span class=\"n\">hookspec</span><span class=\"o\">.</span><span
  class=\"n\">MarkataSpecs</span><span class=\"p\">)</span>\n                <span
  class=\"k\">if</span> <span class=\"n\">config</span> <span class=\"ow\">is</span>
  <span class=\"ow\">not</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">config</span> <span class=\"o\">=</span> <span class=\"n\">config</span>\n
  \               <span class=\"k\">with</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">cache</span> <span class=\"k\">as</span> <span
  class=\"n\">cache</span><span class=\"p\">:</span>\n                    <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">init_cache_stats</span> <span class=\"o\">=</span>
  <span class=\"n\">cache</span><span class=\"o\">.</span><span class=\"n\">stats</span><span
  class=\"p\">()</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">registered_attrs</span> <span class=\"o\">=</span> <span class=\"n\">hookspec</span><span
  class=\"o\">.</span><span class=\"n\">registered_attrs</span>\n                <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">post_models</span>
  <span class=\"o\">=</span> <span class=\"p\">[]</span>\n                <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">config_models</span> <span class=\"o\">=</span>
  <span class=\"p\">[]</span>\n                <span class=\"k\">if</span> <span class=\"n\">config</span>
  <span class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n                    <span class=\"n\">raw_hooks</span> <span
  class=\"o\">=</span> <span class=\"n\">config</span>\n                <span class=\"k\">else</span><span
  class=\"p\">:</span>\n                    <span class=\"n\">raw_hooks</span> <span
  class=\"o\">=</span> <span class=\"n\">standard_config</span><span class=\"o\">.</span><span
  class=\"n\">load</span><span class=\"p\">(</span><span class=\"s2\">\"markata\"</span><span
  class=\"p\">)</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">hooks_conf</span> <span class=\"o\">=</span> <span class=\"n\">HooksConfig</span><span
  class=\"o\">.</span><span class=\"n\">parse_obj</span><span class=\"p\">(</span><span
  class=\"n\">raw_hooks</span><span class=\"p\">)</span>\n                <span class=\"k\">try</span><span
  class=\"p\">:</span>\n                    <span class=\"n\">default_index</span>
  <span class=\"o\">=</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">hooks_conf</span><span class=\"o\">.</span><span class=\"n\">hooks</span><span
  class=\"o\">.</span><span class=\"n\">index</span><span class=\"p\">(</span><span
  class=\"s2\">\"default\"</span><span class=\"p\">)</span>\n                    <span
  class=\"n\">hooks</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
  \                       <span class=\"o\">*</span><span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">hooks_conf</span><span class=\"o\">.</span><span
  class=\"n\">hooks</span><span class=\"p\">[:</span><span class=\"n\">default_index</span><span
  class=\"p\">],</span>\n                        <span class=\"o\">*</span><span class=\"n\">DEFAULT_HOOKS</span><span
  class=\"p\">,</span>\n                        <span class=\"o\">*</span><span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">hooks_conf</span><span class=\"o\">.</span><span
  class=\"n\">hooks</span><span class=\"p\">[</span><span class=\"n\">default_index</span>
  <span class=\"o\">+</span> <span class=\"mi\">1</span> <span class=\"p\">:],</span>\n
  \                   <span class=\"p\">]</span>\n                    <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">hooks_conf</span><span class=\"o\">.</span><span
  class=\"n\">hooks</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
  \                       <span class=\"n\">hook</span> <span class=\"k\">for</span>
  <span class=\"n\">hook</span> <span class=\"ow\">in</span> <span class=\"n\">hooks</span>
  <span class=\"k\">if</span> <span class=\"n\">hook</span> <span class=\"ow\">not</span>
  <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">hooks_conf</span><span class=\"o\">.</span><span class=\"n\">disabled_hooks</span>\n
  \                   <span class=\"p\">]</span>\n                <span class=\"k\">except</span>
  <span class=\"ne\">ValueError</span><span class=\"p\">:</span>\n                    <span
  class=\"c1\"># 'default' is not in hooks , do not replace with default_hooks</span>\n
  \                   <span class=\"k\">pass</span>\n\n                <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">_register_hooks</span><span class=\"p\">()</span>\n
  \               <span class=\"k\">if</span> <span class=\"n\">console</span> <span
  class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n                    <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">_console</span> <span class=\"o\">=</span> <span class=\"n\">console</span>\n
  \               <span class=\"n\">atexit</span><span class=\"o\">.</span><span class=\"n\">register</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">teardown</span><span class=\"p\">)</span>\n                <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">precache</span>\n\n            <span class=\"nd\">@property</span>\n
  \           <span class=\"k\">def</span> <span class=\"nf\">cache</span><span class=\"p\">(</span><span
  class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Cache</span><span
  class=\"p\">:</span>\n                <span class=\"c1\"># if self.threded:</span>\n
  \               <span class=\"c1\">#     FanoutCache(self.MARKATA_CACHE_DIR, statistics=True)</span>\n
  \               <span class=\"k\">if</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">_cache</span> <span class=\"ow\">is</span>
  <span class=\"ow\">not</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \                   <span class=\"k\">return</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">_cache</span>\n                <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">_cache</span> <span class=\"o\">=</span> <span
  class=\"n\">Cache</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">MARKATA_CACHE_DIR</span><span class=\"p\">,</span>
  <span class=\"n\">statistics</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
  class=\"p\">)</span>\n\n                <span class=\"k\">return</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">_cache</span>\n\n            <span class=\"nd\">@property</span>\n
  \           <span class=\"k\">def</span> <span class=\"nf\">precache</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">:</span> <span
  class=\"s2\">\"Markata\"</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
  class=\"k\">if</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">_precache</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n                    <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">cache</span><span class=\"o\">.</span><span class=\"n\">expire</span><span
  class=\"p\">()</span>\n                    <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">_precache</span> <span class=\"o\">=</span> <span class=\"p\">{</span><span
  class=\"n\">k</span><span class=\"p\">:</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">cache</span><span class=\"o\">.</span><span
  class=\"n\">get</span><span class=\"p\">(</span><span class=\"n\">k</span><span
  class=\"p\">)</span> <span class=\"k\">for</span> <span class=\"n\">k</span> <span
  class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">cache</span><span class=\"o\">.</span><span class=\"n\">iterkeys</span><span
  class=\"p\">()}</span>\n                <span class=\"k\">return</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">_precache</span>\n\n            <span class=\"k\">def</span>
  <span class=\"fm\">__getattr__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">,</span>
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
  class=\"p\">(</span>\n                        <span class=\"p\">[</span><span class=\"n\">attr</span><span
  class=\"p\">[</span><span class=\"s2\">\"lifecycle\"</span><span class=\"p\">]</span>
  <span class=\"k\">for</span> <span class=\"n\">attr</span> <span class=\"ow\">in</span>
  <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">registered_attrs</span><span
  class=\"p\">[</span><span class=\"n\">item</span><span class=\"p\">]],</span>\n
  \                   <span class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">name</span>\n
  \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">run</span><span class=\"p\">(</span><span class=\"n\">stage_to_run_to</span><span
  class=\"p\">)</span>\n                    <span class=\"k\">return</span> <span
  class=\"nb\">getattr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">,</span> <span class=\"n\">item</span><span class=\"p\">)</span>\n                <span
  class=\"k\">else</span><span class=\"p\">:</span>\n                    <span class=\"c1\">#
  Markata does not know what this is, raise</span>\n                    <span class=\"k\">raise</span>
  <span class=\"ne\">AttributeError</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
  class=\"s2\">\"'Markata' object has no attribute '</span><span class=\"si\">{</span><span
  class=\"n\">item</span><span class=\"si\">}</span><span class=\"s2\">'\"</span><span
  class=\"p\">)</span>\n\n            <span class=\"k\">def</span> <span class=\"nf\">__rich__</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">:</span> <span
  class=\"s2\">\"Markata\"</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"n\">Table</span><span class=\"p\">:</span>\n                <span
  class=\"n\">grid</span> <span class=\"o\">=</span> <span class=\"n\">Table</span><span
  class=\"o\">.</span><span class=\"n\">grid</span><span class=\"p\">()</span>\n                <span
  class=\"n\">grid</span><span class=\"o\">.</span><span class=\"n\">add_column</span><span
  class=\"p\">(</span><span class=\"s2\">\"label\"</span><span class=\"p\">)</span>\n
  \               <span class=\"n\">grid</span><span class=\"o\">.</span><span class=\"n\">add_column</span><span
  class=\"p\">(</span><span class=\"s2\">\"value\"</span><span class=\"p\">)</span>\n\n
  \               <span class=\"k\">for</span> <span class=\"n\">label</span><span
  class=\"p\">,</span> <span class=\"n\">value</span> <span class=\"ow\">in</span>
  <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">describe</span><span
  class=\"p\">()</span><span class=\"o\">.</span><span class=\"n\">items</span><span
  class=\"p\">():</span>\n                    <span class=\"n\">grid</span><span class=\"o\">.</span><span
  class=\"n\">add_row</span><span class=\"p\">(</span><span class=\"n\">label</span><span
  class=\"p\">,</span> <span class=\"n\">value</span><span class=\"p\">)</span>\n\n
  \               <span class=\"k\">return</span> <span class=\"n\">grid</span>\n\n
  \           <span class=\"k\">def</span> <span class=\"nf\">bust_cache</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">:</span> <span
  class=\"s2\">\"Markata\"</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"n\">Markata</span><span class=\"p\">:</span>\n                <span
  class=\"k\">with</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">cache</span> <span class=\"k\">as</span> <span class=\"n\">cache</span><span
  class=\"p\">:</span>\n                    <span class=\"n\">cache</span><span class=\"o\">.</span><span
  class=\"n\">clear</span><span class=\"p\">()</span>\n                <span class=\"k\">return</span>
  <span class=\"bp\">self</span>\n\n            <span class=\"k\">def</span> <span
  class=\"nf\">configure</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Markata</span><span
  class=\"p\">:</span>\n                <span class=\"n\">sys</span><span class=\"o\">.</span><span
  class=\"n\">path</span><span class=\"o\">.</span><span class=\"n\">append</span><span
  class=\"p\">(</span><span class=\"n\">os</span><span class=\"o\">.</span><span class=\"n\">getcwd</span><span
  class=\"p\">())</span>\n                <span class=\"c1\"># self.config = {**DEFUALT_CONFIG,
  **standard_config.load(\"markata\")}</span>\n                <span class=\"c1\">#
  if isinstance(self.config[\"glob_patterns\"], str):</span>\n                <span
  class=\"c1\">#     self.config[\"glob_patterns\"] = self.config[\"glob_patterns\"].split(\",\")</span>\n
  \               <span class=\"c1\"># elif isinstance(self.config[\"glob_patterns\"],
  list):</span>\n                <span class=\"c1\">#     self.config[\"glob_patterns\"]
  = list(self.config[\"glob_patterns\"])</span>\n                <span class=\"c1\">#
  else:</span>\n                <span class=\"c1\">#     raise TypeError(\"glob_patterns
  must be list or str\")</span>\n                <span class=\"c1\"># self.glob_patterns
  = self.config[\"glob_patterns\"]</span>\n\n                <span class=\"c1\">#
  self.hooks = self.config[\"hooks\"]</span>\n\n                <span class=\"c1\">#
  if \"disabled_hooks\" not in self.config:</span>\n                <span class=\"c1\">#
  \    self.disabled_hooks = [\"\"]</span>\n                <span class=\"c1\"># if
  isinstance(self.config[\"disabled_hooks\"], str):</span>\n                <span
  class=\"c1\">#     self.disabled_hooks = self.config[\"disabled_hooks\"].split(\",\")</span>\n
  \               <span class=\"c1\"># if isinstance(self.config[\"disabled_hooks\"],
  list):</span>\n                <span class=\"c1\">#     self.disabled_hooks = self.config[\"disabled_hooks\"]</span>\n\n
  \               <span class=\"c1\"># if not self.config.get(\"output_dir\", \"markout\").endswith(</span>\n
  \               <span class=\"c1\">#     self.config.get(\"path_prefix\", \"\")</span>\n
  \               <span class=\"c1\"># ):</span>\n                <span class=\"c1\">#
  \    self.config[\"output_dir\"] = (</span>\n                <span class=\"c1\">#
  \        self.config.get(\"output_dir\", \"markout\") +</span>\n                <span
  class=\"c1\">#         \"/\" +</span>\n                <span class=\"c1\">#         self.config.get(\"path_prefix\",
  \"\").rstrip(\"/\")</span>\n                <span class=\"c1\">#     )</span>\n
  \               <span class=\"c1\"># if (</span>\n                <span class=\"c1\">#
  \    len((output_split := self.config.get(\"output_dir\", \"markout\").split(\"/\")))
  &gt;</span>\n                <span class=\"c1\">#     1</span>\n                <span
  class=\"c1\"># ):</span>\n                <span class=\"c1\">#     if \"path_prefix\"
  not in self.config.keys():</span>\n                <span class=\"c1\">#         self.config[\"path_prefix\"]
  = \"/\".join(output_split[1:]) + \"/\"</span>\n                <span class=\"c1\">#
  if not self.config.get(\"path_prefix\", \"\").endswith(\"/\"):</span>\n                <span
  class=\"c1\">#     self.config[\"path_prefix\"] = self.config.get(\"path_prefix\",
  \"\") + \"/\"</span>\n\n                <span class=\"c1\"># self.config[\"output_dir\"]
  = self.config[\"output_dir\"].lstrip(\"/\")</span>\n                <span class=\"c1\">#
  self.config[\"path_prefix\"] = self.config[\"path_prefix\"].lstrip(\"/\")</span>\n\n
  \               <span class=\"k\">try</span><span class=\"p\">:</span>\n                    <span
  class=\"n\">default_index</span> <span class=\"o\">=</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">hooks_conf</span><span class=\"o\">.</span><span
  class=\"n\">hooks</span><span class=\"o\">.</span><span class=\"n\">index</span><span
  class=\"p\">(</span><span class=\"s2\">\"default\"</span><span class=\"p\">)</span>\n
  \                   <span class=\"n\">hooks</span> <span class=\"o\">=</span> <span
  class=\"p\">[</span>\n                        <span class=\"o\">*</span><span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">hooks_conf</span><span class=\"o\">.</span><span
  class=\"n\">hooks</span><span class=\"p\">[:</span><span class=\"n\">default_index</span><span
  class=\"p\">],</span>\n                        <span class=\"o\">*</span><span class=\"n\">DEFAULT_HOOKS</span><span
  class=\"p\">,</span>\n                        <span class=\"o\">*</span><span class=\"bp\">self</span><span
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
  class=\"c1\"># 'default' is not in hooks , do not replace with default_hooks</span>\n
  \                   <span class=\"k\">pass</span>\n\n                <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">_pm</span> <span class=\"o\">=</span> <span
  class=\"n\">pluggy</span><span class=\"o\">.</span><span class=\"n\">PluginManager</span><span
  class=\"p\">(</span><span class=\"s2\">\"markata\"</span><span class=\"p\">)</span>\n
  \               <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_pm</span><span
  class=\"o\">.</span><span class=\"n\">add_hookspecs</span><span class=\"p\">(</span><span
  class=\"n\">hookspec</span><span class=\"o\">.</span><span class=\"n\">MarkataSpecs</span><span
  class=\"p\">)</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">_register_hooks</span><span class=\"p\">()</span>\n\n                <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_pm</span><span
  class=\"o\">.</span><span class=\"n\">hook</span><span class=\"o\">.</span><span
  class=\"n\">configure</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
  class=\"o\">=</span><span class=\"bp\">self</span><span class=\"p\">)</span>\n                <span
  class=\"k\">return</span> <span class=\"bp\">self</span>\n\n            <span class=\"k\">def</span>
  <span class=\"nf\">get_plugin_config</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">,</span> <span class=\"n\">path_or_name</span><span class=\"p\">:</span>
  <span class=\"nb\">str</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"n\">Dict</span><span class=\"p\">:</span>\n                <span class=\"n\">key</span>
  <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
  class=\"n\">path_or_name</span><span class=\"p\">)</span><span class=\"o\">.</span><span
  class=\"n\">stem</span>\n\n                <span class=\"n\">config</span> <span
  class=\"o\">=</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">get</span><span
  class=\"p\">(</span><span class=\"n\">key</span><span class=\"p\">,</span> <span
  class=\"p\">{})</span>\n\n                <span class=\"k\">if</span> <span class=\"ow\">not</span>
  <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">config</span><span
  class=\"p\">,</span> <span class=\"nb\">dict</span><span class=\"p\">):</span>\n
  \                   <span class=\"k\">raise</span> <span class=\"ne\">TypeError</span><span
  class=\"p\">(</span><span class=\"s2\">\"must use dict\"</span><span class=\"p\">)</span>\n
  \               <span class=\"k\">if</span> <span class=\"s2\">\"cache_expire\"</span>
  <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">keys</span><span class=\"p\">():</span>\n
  \                   <span class=\"n\">config</span><span class=\"p\">[</span><span
  class=\"s2\">\"cache_expire\"</span><span class=\"p\">]</span> <span class=\"o\">=</span>
  <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"p\">[</span><span class=\"s2\">\"default_cache_expire\"</span><span class=\"p\">]</span>\n
  \               <span class=\"k\">if</span> <span class=\"s2\">\"config_key\"</span>
  <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">keys</span><span class=\"p\">():</span>\n
  \                   <span class=\"n\">config</span><span class=\"p\">[</span><span
  class=\"s2\">\"config_key\"</span><span class=\"p\">]</span> <span class=\"o\">=</span>
  <span class=\"n\">key</span>\n                <span class=\"k\">return</span> <span
  class=\"n\">config</span>\n\n            <span class=\"k\">def</span> <span class=\"nf\">get_config</span><span
  class=\"p\">(</span>\n                <span class=\"bp\">self</span><span class=\"p\">,</span>\n
  \               <span class=\"n\">key</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
  class=\"p\">,</span>\n                <span class=\"n\">default</span><span class=\"p\">:</span>
  <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">\"\"</span><span
  class=\"p\">,</span>\n                <span class=\"n\">warn</span><span class=\"p\">:</span>
  <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"kc\">True</span><span
  class=\"p\">,</span>\n                <span class=\"n\">suggested</span><span class=\"p\">:</span>
  <span class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
  class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span><span
  class=\"p\">,</span>\n            <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"n\">Any</span><span class=\"p\">:</span>\n                <span class=\"k\">if</span>
  <span class=\"n\">key</span> <span class=\"ow\">in</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
  class=\"n\">keys</span><span class=\"p\">():</span>\n                    <span class=\"k\">return</span>
  <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"p\">[</span><span class=\"n\">key</span><span class=\"p\">]</span>\n                <span
  class=\"k\">else</span><span class=\"p\">:</span>\n                    <span class=\"k\">if</span>
  <span class=\"n\">suggested</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n                        <span class=\"n\">suggested</span>
  <span class=\"o\">=</span> <span class=\"n\">textwrap</span><span class=\"o\">.</span><span
  class=\"n\">dedent</span><span class=\"p\">(</span>\n                            <span
  class=\"sa\">f</span><span class=\"s2\">\"\"\"</span>\n<span class=\"s2\">                            [markata]</span>\n<span
  class=\"s2\">                            </span><span class=\"si\">{</span><span
  class=\"n\">key</span><span class=\"si\">}</span><span class=\"s2\"> = '</span><span
  class=\"si\">{</span><span class=\"n\">default</span><span class=\"si\">}</span><span
  class=\"s2\">'</span>\n<span class=\"s2\">                            \"\"\"</span>\n
  \                       <span class=\"p\">)</span>\n                    <span class=\"k\">if</span>
  <span class=\"n\">warn</span><span class=\"p\">:</span>\n                        <span
  class=\"n\">logger</span><span class=\"o\">.</span><span class=\"n\">warning</span><span
  class=\"p\">(</span>\n                            <span class=\"n\">textwrap</span><span
  class=\"o\">.</span><span class=\"n\">dedent</span><span class=\"p\">(</span>\n
  \                               <span class=\"sa\">f</span><span class=\"s2\">\"\"\"</span>\n<span
  class=\"s2\">                                Warning </span><span class=\"si\">{</span><span
  class=\"n\">key</span><span class=\"si\">}</span><span class=\"s2\"> is not set
  in markata config, sitemap will</span>\n<span class=\"s2\">                                be
  missing root site_name</span>\n<span class=\"s2\">                                to
  resolve this open your markata.toml and add</span>\n\n<span class=\"s2\">                                </span><span
  class=\"si\">{</span><span class=\"n\">suggested</span><span class=\"si\">}</span>\n\n<span
  class=\"s2\">                                \"\"\"</span>\n                            <span
  class=\"p\">),</span>\n                        <span class=\"p\">)</span>\n                <span
  class=\"k\">return</span> <span class=\"n\">default</span>\n\n            <span
  class=\"k\">def</span> <span class=\"nf\">make_hash</span><span class=\"p\">(</span><span
  class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
  class=\"n\">keys</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span
  class=\"p\">:</span>\n                <span class=\"n\">str_keys</span> <span class=\"o\">=</span>
  <span class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">(</span><span
  class=\"n\">key</span><span class=\"p\">)</span> <span class=\"k\">for</span> <span
  class=\"n\">key</span> <span class=\"ow\">in</span> <span class=\"n\">keys</span><span
  class=\"p\">]</span>\n                <span class=\"k\">return</span> <span class=\"n\">hashlib</span><span
  class=\"o\">.</span><span class=\"n\">md5</span><span class=\"p\">(</span><span
  class=\"s2\">\"\"</span><span class=\"o\">.</span><span class=\"n\">join</span><span
  class=\"p\">(</span><span class=\"n\">str_keys</span><span class=\"p\">)</span><span
  class=\"o\">.</span><span class=\"n\">encode</span><span class=\"p\">(</span><span
  class=\"s2\">\"utf-8\"</span><span class=\"p\">))</span><span class=\"o\">.</span><span
  class=\"n\">hexdigest</span><span class=\"p\">()</span>\n\n            <span class=\"nd\">@property</span>\n
  \           <span class=\"k\">def</span> <span class=\"nf\">content_dir_hash</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">:</span> <span
  class=\"s2\">\"Markata\"</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"nb\">str</span><span class=\"p\">:</span>\n                <span class=\"n\">hashes</span>
  <span class=\"o\">=</span> <span class=\"p\">[</span>\n                    <span
  class=\"n\">dirhash</span><span class=\"p\">(</span><span class=\"nb\">dir</span><span
  class=\"p\">)</span>\n                    <span class=\"k\">for</span> <span class=\"nb\">dir</span>
  <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">content_directories</span>\n                    <span class=\"k\">if</span>
  <span class=\"nb\">dir</span><span class=\"o\">.</span><span class=\"n\">absolute</span><span
  class=\"p\">()</span> <span class=\"o\">!=</span> <span class=\"n\">Path</span><span
  class=\"p\">(</span><span class=\"s2\">\".\"</span><span class=\"p\">)</span><span
  class=\"o\">.</span><span class=\"n\">absolute</span><span class=\"p\">()</span>\n
  \               <span class=\"p\">]</span>\n                <span class=\"k\">return</span>
  <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">make_hash</span><span
  class=\"p\">(</span><span class=\"o\">*</span><span class=\"n\">hashes</span><span
  class=\"p\">)</span>\n\n            <span class=\"nd\">@property</span>\n            <span
  class=\"k\">def</span> <span class=\"nf\">console</span><span class=\"p\">(</span><span
  class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span
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
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"nb\">dict</span><span class=\"p\">[</span><span
  class=\"nb\">str</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
  class=\"p\">]:</span>\n                <span class=\"k\">return</span> <span class=\"p\">{</span><span
  class=\"s2\">\"version\"</span><span class=\"p\">:</span> <span class=\"n\">__version__</span><span
  class=\"p\">}</span>\n\n            <span class=\"k\">def</span> <span class=\"nf\">_to_dict</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">:</span> <span
  class=\"s2\">\"Markata\"</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"nb\">dict</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
  class=\"p\">,</span> <span class=\"n\">Iterable</span><span class=\"p\">]:</span>\n
  \               <span class=\"k\">return</span> <span class=\"p\">{</span><span
  class=\"s2\">\"config\"</span><span class=\"p\">:</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">,</span> <span
  class=\"s2\">\"articles\"</span><span class=\"p\">:</span> <span class=\"p\">[</span><span
  class=\"n\">a</span><span class=\"o\">.</span><span class=\"n\">to_dict</span><span
  class=\"p\">()</span> <span class=\"k\">for</span> <span class=\"n\">a</span> <span
  class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">articles</span><span class=\"p\">]}</span>\n\n            <span class=\"k\">def</span>
  <span class=\"nf\">to_dict</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"nb\">dict</span><span class=\"p\">:</span>\n
  \               <span class=\"k\">return</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">_to_dict</span><span class=\"p\">()</span>\n\n
  \           <span class=\"k\">def</span> <span class=\"nf\">to_json</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">:</span> <span
  class=\"s2\">\"Markata\"</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"nb\">str</span><span class=\"p\">:</span>\n                <span class=\"kn\">import</span>
  <span class=\"nn\">json</span>\n\n                <span class=\"k\">return</span>
  <span class=\"n\">json</span><span class=\"o\">.</span><span class=\"n\">dumps</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">to_dict</span><span class=\"p\">(),</span> <span class=\"n\">indent</span><span
  class=\"o\">=</span><span class=\"mi\">4</span><span class=\"p\">,</span> <span
  class=\"n\">sort_keys</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
  class=\"p\">,</span> <span class=\"n\">default</span><span class=\"o\">=</span><span
  class=\"nb\">str</span><span class=\"p\">)</span>\n\n            <span class=\"k\">def</span>
  <span class=\"nf\">_register_hooks</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">)</span>
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
  class=\"n\">e</span><span class=\"p\">:</span>\n                        <span class=\"c1\">#
  class style plugins</span>\n                        <span class=\"k\">if</span>
  <span class=\"s2\">\".\"</span> <span class=\"ow\">in</span> <span class=\"n\">hook</span><span
  class=\"p\">:</span>\n                            <span class=\"k\">try</span><span
  class=\"p\">:</span>\n                                <span class=\"n\">mod</span>
  <span class=\"o\">=</span> <span class=\"n\">importlib</span><span class=\"o\">.</span><span
  class=\"n\">import_module</span><span class=\"p\">(</span><span class=\"s2\">\".\"</span><span
  class=\"o\">.</span><span class=\"n\">join</span><span class=\"p\">(</span><span
  class=\"n\">hook</span><span class=\"o\">.</span><span class=\"n\">split</span><span
  class=\"p\">(</span><span class=\"s2\">\".\"</span><span class=\"p\">)[:</span><span
  class=\"o\">-</span><span class=\"mi\">1</span><span class=\"p\">]))</span>\n                                <span
  class=\"n\">plugin</span> <span class=\"o\">=</span> <span class=\"nb\">getattr</span><span
  class=\"p\">(</span><span class=\"n\">mod</span><span class=\"p\">,</span> <span
  class=\"n\">hook</span><span class=\"o\">.</span><span class=\"n\">split</span><span
  class=\"p\">(</span><span class=\"s2\">\".\"</span><span class=\"p\">)[</span><span
  class=\"o\">-</span><span class=\"mi\">1</span><span class=\"p\">])</span>\n                            <span
  class=\"k\">except</span> <span class=\"ne\">ModuleNotFoundError</span> <span class=\"k\">as</span>
  <span class=\"n\">e</span><span class=\"p\">:</span>\n                                <span
  class=\"k\">raise</span> <span class=\"ne\">ModuleNotFoundError</span><span class=\"p\">(</span>\n
  \                                   <span class=\"sa\">f</span><span class=\"s2\">\"module
  </span><span class=\"si\">{</span><span class=\"n\">hook</span><span class=\"si\">}</span><span
  class=\"s2\"> not found</span><span class=\"se\">\\n</span><span class=\"si\">{</span><span
  class=\"n\">sys</span><span class=\"o\">.</span><span class=\"n\">path</span><span
  class=\"si\">}</span><span class=\"s2\">\"</span>\n                                <span
  class=\"p\">)</span> <span class=\"kn\">from</span> <span class=\"nn\">e</span>\n
  \                       <span class=\"k\">else</span><span class=\"p\">:</span>\n
  \                           <span class=\"k\">raise</span> <span class=\"n\">e</span>\n\n
  \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">_pm</span><span class=\"o\">.</span><span class=\"n\">register</span><span
  class=\"p\">(</span><span class=\"n\">plugin</span><span class=\"p\">)</span>\n\n
  \           <span class=\"k\">def</span> <span class=\"fm\">__iter__</span><span
  class=\"p\">(</span>\n                <span class=\"bp\">self</span><span class=\"p\">:</span>
  <span class=\"s2\">\"Markata\"</span><span class=\"p\">,</span> <span class=\"n\">description</span><span
  class=\"p\">:</span> <span class=\"nb\">str</span> <span class=\"o\">=</span> <span
  class=\"s2\">\"working...\"</span>\n            <span class=\"p\">)</span> <span
  class=\"o\">-&gt;</span> <span class=\"n\">Iterable</span><span class=\"p\">[</span><span
  class=\"s2\">\"Markata.Post\"</span><span class=\"p\">]:</span>\n                <span
  class=\"n\">articles</span><span class=\"p\">:</span> <span class=\"n\">Iterable</span><span
  class=\"p\">[</span><span class=\"n\">Markata</span><span class=\"o\">.</span><span
  class=\"n\">Post</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
  class=\"n\">track</span><span class=\"p\">(</span>\n                    <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">articles</span><span class=\"p\">,</span>\n
  \                   <span class=\"n\">description</span><span class=\"o\">=</span><span
  class=\"n\">description</span><span class=\"p\">,</span>\n                    <span
  class=\"n\">transient</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
  class=\"p\">,</span>\n                    <span class=\"n\">console</span><span
  class=\"o\">=</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">console</span><span class=\"p\">,</span>\n                <span class=\"p\">)</span>\n
  \               <span class=\"k\">return</span> <span class=\"n\">articles</span>\n\n
  \           <span class=\"k\">def</span> <span class=\"nf\">iter_articles</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">:</span> <span
  class=\"s2\">\"Markata\"</span><span class=\"p\">,</span> <span class=\"n\">description</span><span
  class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">)</span> <span
  class=\"o\">-&gt;</span> <span class=\"n\">Iterable</span><span class=\"p\">[</span><span
  class=\"n\">Markata</span><span class=\"o\">.</span><span class=\"n\">Post</span><span
  class=\"p\">]:</span>\n                <span class=\"n\">articles</span><span class=\"p\">:</span>
  <span class=\"n\">Iterable</span><span class=\"p\">[</span><span class=\"n\">Markata</span><span
  class=\"o\">.</span><span class=\"n\">Post</span><span class=\"p\">]</span> <span
  class=\"o\">=</span> <span class=\"n\">track</span><span class=\"p\">(</span>\n
  \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">articles</span><span class=\"p\">,</span>\n                    <span
  class=\"n\">description</span><span class=\"o\">=</span><span class=\"n\">description</span><span
  class=\"p\">,</span>\n                    <span class=\"n\">transient</span><span
  class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n                    <span
  class=\"n\">console</span><span class=\"o\">=</span><span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">console</span><span class=\"p\">,</span>\n
  \               <span class=\"p\">)</span>\n                <span class=\"k\">return</span>
  <span class=\"n\">articles</span>\n\n            <span class=\"k\">def</span> <span
  class=\"nf\">teardown</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"n\">Markata</span><span class=\"p\">:</span>\n<span
  class=\"w\">                </span><span class=\"sd\">\"\"\"give special access
  to the teardown lifecycle method\"\"\"</span>\n                <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">_pm</span><span class=\"o\">.</span><span
  class=\"n\">hook</span><span class=\"o\">.</span><span class=\"n\">teardown</span><span
  class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
  class=\"bp\">self</span><span class=\"p\">)</span>\n                <span class=\"k\">return</span>
  <span class=\"bp\">self</span>\n\n            <span class=\"k\">def</span> <span
  class=\"nf\">run</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">,</span>
  <span class=\"n\">lifecycle</span><span class=\"p\">:</span> <span class=\"n\">LifeCycle</span>
  <span class=\"o\">=</span> <span class=\"kc\">None</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"n\">Markata</span><span class=\"p\">:</span>\n
  \               <span class=\"k\">if</span> <span class=\"n\">lifecycle</span> <span
  class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \                   <span class=\"n\">lifecycle</span> <span class=\"o\">=</span>
  <span class=\"nb\">max</span><span class=\"p\">(</span><span class=\"n\">LifeCycle</span><span
  class=\"o\">.</span><span class=\"n\">_member_map_</span><span class=\"o\">.</span><span
  class=\"n\">values</span><span class=\"p\">())</span>\n\n                <span class=\"k\">if</span>
  <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">lifecycle</span><span
  class=\"p\">,</span> <span class=\"nb\">str</span><span class=\"p\">):</span>\n
  \                   <span class=\"n\">lifecycle</span> <span class=\"o\">=</span>
  <span class=\"n\">LifeCycle</span><span class=\"p\">[</span><span class=\"n\">lifecycle</span><span
  class=\"p\">]</span>\n\n                <span class=\"n\">stages_to_run</span> <span
  class=\"o\">=</span> <span class=\"p\">[</span>\n                    <span class=\"n\">m</span>\n
  \                   <span class=\"k\">for</span> <span class=\"n\">m</span> <span
  class=\"ow\">in</span> <span class=\"n\">LifeCycle</span><span class=\"o\">.</span><span
  class=\"n\">_member_map_</span>\n                    <span class=\"k\">if</span>
  <span class=\"p\">(</span><span class=\"n\">LifeCycle</span><span class=\"p\">[</span><span
  class=\"n\">m</span><span class=\"p\">]</span> <span class=\"o\">&lt;=</span> <span
  class=\"n\">lifecycle</span><span class=\"p\">)</span> <span class=\"ow\">and</span>
  <span class=\"p\">(</span><span class=\"n\">m</span> <span class=\"ow\">not</span>
  <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">stages_ran</span><span class=\"p\">)</span>\n                <span class=\"p\">]</span>\n\n
  \               <span class=\"k\">if</span> <span class=\"ow\">not</span> <span
  class=\"n\">stages_to_run</span><span class=\"p\">:</span>\n                    <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">console</span><span
  class=\"o\">.</span><span class=\"n\">log</span><span class=\"p\">(</span><span
  class=\"sa\">f</span><span class=\"s2\">\"</span><span class=\"si\">{</span><span
  class=\"n\">lifecycle</span><span class=\"o\">.</span><span class=\"n\">name</span><span
  class=\"si\">}</span><span class=\"s2\"> already ran\"</span><span class=\"p\">)</span>\n
  \                   <span class=\"k\">return</span> <span class=\"bp\">self</span>\n\n
  \               <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">console</span><span
  class=\"o\">.</span><span class=\"n\">log</span><span class=\"p\">(</span><span
  class=\"sa\">f</span><span class=\"s2\">\"running </span><span class=\"si\">{</span><span
  class=\"n\">stages_to_run</span><span class=\"si\">}</span><span class=\"s2\">\"</span><span
  class=\"p\">)</span>\n                <span class=\"k\">for</span> <span class=\"n\">stage</span>
  <span class=\"ow\">in</span> <span class=\"n\">stages_to_run</span><span class=\"p\">:</span>\n
  \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
  class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">\"</span><span
  class=\"si\">{</span><span class=\"n\">stage</span><span class=\"si\">}</span><span
  class=\"s2\"> running\"</span><span class=\"p\">)</span>\n                    <span
  class=\"nb\">getattr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">_pm</span><span class=\"o\">.</span><span
  class=\"n\">hook</span><span class=\"p\">,</span> <span class=\"n\">stage</span><span
  class=\"p\">)(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
  class=\"bp\">self</span><span class=\"p\">)</span>\n                    <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">stages_ran</span><span class=\"o\">.</span><span
  class=\"n\">add</span><span class=\"p\">(</span><span class=\"n\">stage</span><span
  class=\"p\">)</span>\n                    <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
  class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">\"</span><span
  class=\"si\">{</span><span class=\"n\">stage</span><span class=\"si\">}</span><span
  class=\"s2\"> complete\"</span><span class=\"p\">)</span>\n\n                <span
  class=\"k\">with</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">cache</span> <span class=\"k\">as</span> <span class=\"n\">cache</span><span
  class=\"p\">:</span>\n                    <span class=\"n\">hits</span><span class=\"p\">,</span>
  <span class=\"n\">misses</span> <span class=\"o\">=</span> <span class=\"n\">cache</span><span
  class=\"o\">.</span><span class=\"n\">stats</span><span class=\"p\">()</span>\n\n
  \               <span class=\"k\">if</span> <span class=\"n\">hits</span> <span
  class=\"o\">+</span> <span class=\"n\">misses</span> <span class=\"o\">&gt;</span>
  <span class=\"mi\">0</span><span class=\"p\">:</span>\n                    <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">console</span><span
  class=\"o\">.</span><span class=\"n\">log</span><span class=\"p\">(</span>\n                        <span
  class=\"sa\">f</span><span class=\"s2\">\"lifetime cache hit rate </span><span class=\"si\">{</span><span
  class=\"nb\">round</span><span class=\"p\">(</span><span class=\"n\">hits</span><span
  class=\"o\">/</span><span class=\"w\"> </span><span class=\"p\">(</span><span class=\"n\">hits</span><span
  class=\"w\"> </span><span class=\"o\">+</span><span class=\"w\"> </span><span class=\"n\">misses</span><span
  class=\"p\">)</span><span class=\"o\">*</span><span class=\"mi\">100</span><span
  class=\"p\">,</span><span class=\"w\"> </span><span class=\"mi\">2</span><span class=\"p\">)</span><span
  class=\"si\">}</span><span class=\"s2\">%\"</span><span class=\"p\">,</span>\n                    <span
  class=\"p\">)</span>\n\n                <span class=\"k\">if</span> <span class=\"n\">misses</span>
  <span class=\"o\">&gt;</span> <span class=\"mi\">0</span><span class=\"p\">:</span>\n
  \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
  class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">\"lifetime cache
  hits/misses </span><span class=\"si\">{</span><span class=\"n\">hits</span><span
  class=\"si\">}</span><span class=\"s2\">/</span><span class=\"si\">{</span><span
  class=\"n\">misses</span><span class=\"si\">}</span><span class=\"s2\">\"</span><span
  class=\"p\">)</span>\n\n                <span class=\"n\">hits</span> <span class=\"o\">-=</span>
  <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">init_cache_stats</span><span
  class=\"p\">[</span><span class=\"mi\">0</span><span class=\"p\">]</span>\n                <span
  class=\"n\">misses</span> <span class=\"o\">-=</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">init_cache_stats</span><span class=\"p\">[</span><span
  class=\"mi\">1</span><span class=\"p\">]</span>\n\n                <span class=\"k\">if</span>
  <span class=\"n\">hits</span> <span class=\"o\">+</span> <span class=\"n\">misses</span>
  <span class=\"o\">&gt;</span> <span class=\"mi\">0</span><span class=\"p\">:</span>\n
  \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
  class=\"p\">(</span>\n                        <span class=\"sa\">f</span><span class=\"s2\">\"run
  cache hit rate </span><span class=\"si\">{</span><span class=\"nb\">round</span><span
  class=\"p\">(</span><span class=\"n\">hits</span><span class=\"o\">/</span><span
  class=\"w\"> </span><span class=\"p\">(</span><span class=\"n\">hits</span><span
  class=\"w\"> </span><span class=\"o\">+</span><span class=\"w\"> </span><span class=\"n\">misses</span><span
  class=\"p\">)</span><span class=\"o\">*</span><span class=\"mi\">100</span><span
  class=\"p\">,</span><span class=\"w\"> </span><span class=\"mi\">2</span><span class=\"p\">)</span><span
  class=\"si\">}</span><span class=\"s2\">%\"</span><span class=\"p\">,</span>\n                    <span
  class=\"p\">)</span>\n\n                <span class=\"k\">if</span> <span class=\"n\">misses</span>
  <span class=\"o\">&gt;</span> <span class=\"mi\">0</span><span class=\"p\">:</span>\n
  \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
  class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">\"run cache hits/misses
  </span><span class=\"si\">{</span><span class=\"n\">hits</span><span class=\"si\">}</span><span
  class=\"s2\">/</span><span class=\"si\">{</span><span class=\"n\">misses</span><span
  class=\"si\">}</span><span class=\"s2\">\"</span><span class=\"p\">)</span>\n\n
  \               <span class=\"k\">return</span> <span class=\"bp\">self</span>\n\n
  \           <span class=\"k\">def</span> <span class=\"nf\">filter</span><span class=\"p\">(</span><span
  class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span
  class=\"p\">,</span> <span class=\"nb\">filter</span><span class=\"p\">:</span>
  <span class=\"nb\">str</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"nb\">list</span><span class=\"p\">:</span>\n                <span
  class=\"k\">def</span> <span class=\"nf\">evalr</span><span class=\"p\">(</span><span
  class=\"n\">a</span><span class=\"p\">:</span> <span class=\"n\">Markata</span><span
  class=\"o\">.</span><span class=\"n\">Post</span><span class=\"p\">)</span> <span
  class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span class=\"p\">:</span>\n
  \                   <span class=\"k\">try</span><span class=\"p\">:</span>\n                        <span
  class=\"k\">return</span> <span class=\"nb\">eval</span><span class=\"p\">(</span>\n
  \                           <span class=\"nb\">filter</span><span class=\"p\">,</span>\n
  \                           <span class=\"p\">{</span><span class=\"o\">**</span><span
  class=\"n\">a</span><span class=\"o\">.</span><span class=\"n\">to_dict</span><span
  class=\"p\">(),</span> <span class=\"s2\">\"timedelta\"</span><span class=\"p\">:</span>
  <span class=\"n\">timedelta</span><span class=\"p\">,</span> <span class=\"s2\">\"post\"</span><span
  class=\"p\">:</span> <span class=\"n\">a</span><span class=\"p\">,</span> <span
  class=\"s2\">\"m\"</span><span class=\"p\">:</span> <span class=\"bp\">self</span><span
  class=\"p\">},</span>\n                            <span class=\"p\">{},</span>\n
  \                       <span class=\"p\">)</span>\n                    <span class=\"k\">except</span>
  <span class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n                        <span
  class=\"k\">return</span> <span class=\"nb\">eval</span><span class=\"p\">(</span>\n
  \                           <span class=\"nb\">filter</span><span class=\"p\">,</span>\n
  \                           <span class=\"p\">{</span><span class=\"o\">**</span><span
  class=\"n\">a</span><span class=\"o\">.</span><span class=\"n\">to_dict</span><span
  class=\"p\">(),</span> <span class=\"s2\">\"timedelta\"</span><span class=\"p\">:</span>
  <span class=\"n\">timedelta</span><span class=\"p\">,</span> <span class=\"s2\">\"post\"</span><span
  class=\"p\">:</span> <span class=\"n\">a</span><span class=\"p\">,</span> <span
  class=\"s2\">\"m\"</span><span class=\"p\">:</span> <span class=\"bp\">self</span><span
  class=\"p\">},</span>\n                            <span class=\"p\">{},</span>\n
  \                       <span class=\"p\">)</span>\n\n                <span class=\"k\">return</span>
  <span class=\"p\">[</span><span class=\"n\">a</span> <span class=\"k\">for</span>
  <span class=\"n\">a</span> <span class=\"ow\">in</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">articles</span> <span class=\"k\">if</span>
  <span class=\"n\">evalr</span><span class=\"p\">(</span><span class=\"n\">a</span><span
  class=\"p\">)]</span>\n\n            <span class=\"k\">def</span> <span class=\"nf\">map</span><span
  class=\"p\">(</span>\n                <span class=\"bp\">self</span><span class=\"p\">:</span>
  <span class=\"s2\">\"Markata\"</span><span class=\"p\">,</span>\n                <span
  class=\"n\">func</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
  <span class=\"o\">=</span> <span class=\"s2\">\"title\"</span><span class=\"p\">,</span>\n
  \               <span class=\"nb\">filter</span><span class=\"p\">:</span> <span
  class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">\"True\"</span><span
  class=\"p\">,</span>\n                <span class=\"n\">sort</span><span class=\"p\">:</span>
  <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">\"True\"</span><span
  class=\"p\">,</span>\n                <span class=\"n\">reverse</span><span class=\"p\">:</span>
  <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"kc\">True</span><span
  class=\"p\">,</span>\n                <span class=\"o\">*</span><span class=\"n\">args</span><span
  class=\"p\">:</span> <span class=\"nb\">tuple</span><span class=\"p\">,</span>\n
  \               <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
  class=\"p\">:</span> <span class=\"nb\">dict</span><span class=\"p\">,</span>\n
  \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">list</span><span
  class=\"p\">:</span>\n                <span class=\"kn\">import</span> <span class=\"nn\">copy</span>\n\n
  \               <span class=\"k\">def</span> <span class=\"nf\">try_sort</span><span
  class=\"p\">(</span><span class=\"n\">a</span><span class=\"p\">:</span> <span class=\"n\">Any</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">int</span><span
  class=\"p\">:</span>\n                    <span class=\"k\">if</span> <span class=\"s2\">\"datetime\"</span>
  <span class=\"ow\">in</span> <span class=\"n\">sort</span><span class=\"o\">.</span><span
  class=\"n\">lower</span><span class=\"p\">():</span>\n                        <span
  class=\"k\">return</span> <span class=\"n\">a</span><span class=\"o\">.</span><span
  class=\"n\">get</span><span class=\"p\">(</span><span class=\"n\">sort</span><span
  class=\"p\">,</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
  class=\"n\">datetime</span><span class=\"p\">(</span><span class=\"mi\">1970</span><span
  class=\"p\">,</span> <span class=\"mi\">1</span><span class=\"p\">,</span> <span
  class=\"mi\">1</span><span class=\"p\">))</span>\n\n                    <span class=\"k\">if</span>
  <span class=\"s2\">\"date\"</span> <span class=\"ow\">in</span> <span class=\"n\">sort</span><span
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
  class=\"p\">(),</span> <span class=\"p\">{})</span>\n                    <span class=\"k\">except</span>
  <span class=\"ne\">NameError</span><span class=\"p\">:</span>\n                        <span
  class=\"k\">return</span> <span class=\"o\">-</span><span class=\"mi\">1</span>\n
  \                   <span class=\"k\">return</span> <span class=\"n\">value</span>\n
  \                   <span class=\"k\">try</span><span class=\"p\">:</span>\n                        <span
  class=\"k\">return</span> <span class=\"nb\">int</span><span class=\"p\">(</span><span
  class=\"n\">value</span><span class=\"p\">)</span>\n                    <span class=\"k\">except</span>
  <span class=\"ne\">TypeError</span><span class=\"p\">:</span>\n                        <span
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
  <span class=\"o\">-</span><span class=\"mi\">1</span>\n\n                <span class=\"n\">articles</span>
  <span class=\"o\">=</span> <span class=\"n\">copy</span><span class=\"o\">.</span><span
  class=\"n\">copy</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">articles</span><span class=\"p\">)</span>\n
  \               <span class=\"n\">articles</span><span class=\"o\">.</span><span
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
  class=\"p\">(),</span> <span class=\"s2\">\"timedelta\"</span><span class=\"p\">:</span>
  <span class=\"n\">timedelta</span><span class=\"p\">,</span> <span class=\"s2\">\"post\"</span><span
  class=\"p\">:</span> <span class=\"n\">a</span><span class=\"p\">,</span> <span
  class=\"s2\">\"m\"</span><span class=\"p\">:</span> <span class=\"bp\">self</span><span
  class=\"p\">},</span>\n                            <span class=\"p\">{},</span>\n
  \                       <span class=\"p\">)</span>\n                        <span
  class=\"k\">for</span> <span class=\"n\">a</span> <span class=\"ow\">in</span> <span
  class=\"n\">articles</span>\n                        <span class=\"k\">if</span>
  <span class=\"nb\">eval</span><span class=\"p\">(</span>\n                            <span
  class=\"nb\">filter</span><span class=\"p\">,</span>\n                            <span
  class=\"p\">{</span><span class=\"o\">**</span><span class=\"n\">a</span><span class=\"o\">.</span><span
  class=\"n\">to_dict</span><span class=\"p\">(),</span> <span class=\"s2\">\"timedelta\"</span><span
  class=\"p\">:</span> <span class=\"n\">timedelta</span><span class=\"p\">,</span>
  <span class=\"s2\">\"post\"</span><span class=\"p\">:</span> <span class=\"n\">a</span><span
  class=\"p\">,</span> <span class=\"s2\">\"m\"</span><span class=\"p\">:</span> <span
  class=\"bp\">self</span><span class=\"p\">},</span>\n                            <span
  class=\"p\">{},</span>\n                        <span class=\"p\">)</span>\n                    <span
  class=\"p\">]</span>\n\n                <span class=\"k\">except</span> <span class=\"ne\">NameError</span>
  <span class=\"k\">as</span> <span class=\"n\">e</span><span class=\"p\">:</span>\n
  \                   <span class=\"n\">variable</span> <span class=\"o\">=</span>
  <span class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">e</span><span
  class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">split</span><span
  class=\"p\">(</span><span class=\"s2\">\"'\"</span><span class=\"p\">)[</span><span
  class=\"mi\">1</span><span class=\"p\">]</span>\n\n                    <span class=\"n\">missing_in_posts</span>
  <span class=\"o\">=</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">map</span><span class=\"p\">(</span>\n                        <span
  class=\"s2\">\"path\"</span><span class=\"p\">,</span>\n                        <span
  class=\"nb\">filter</span><span class=\"o\">=</span><span class=\"sa\">f</span><span
  class=\"s1\">'\"</span><span class=\"si\">{</span><span class=\"n\">variable</span><span
  class=\"si\">}</span><span class=\"s1\">\" not in post.keys()'</span><span class=\"p\">,</span>\n
  \                   <span class=\"p\">)</span>\n                    <span class=\"n\">message</span>
  <span class=\"o\">=</span> <span class=\"p\">(</span>\n                        <span
  class=\"sa\">f</span><span class=\"s2\">\"variable: '</span><span class=\"si\">{</span><span
  class=\"n\">variable</span><span class=\"si\">}</span><span class=\"s2\">' is missing
  in </span><span class=\"si\">{</span><span class=\"nb\">len</span><span class=\"p\">(</span><span
  class=\"n\">missing_in_posts</span><span class=\"p\">)</span><span class=\"si\">}</span><span
  class=\"s2\"> posts\"</span>\n                    <span class=\"p\">)</span>\n                    <span
  class=\"k\">if</span> <span class=\"nb\">len</span><span class=\"p\">(</span><span
  class=\"n\">missing_in_posts</span><span class=\"p\">)</span> <span class=\"o\">&gt;</span>
  <span class=\"mi\">10</span><span class=\"p\">:</span>\n                        <span
  class=\"n\">message</span> <span class=\"o\">+=</span> <span class=\"p\">(</span>\n
  \                           <span class=\"sa\">f</span><span class=\"s2\">\"</span><span
  class=\"se\">\\n</span><span class=\"s2\">first 10 paths to posts missing </span><span
  class=\"si\">{</span><span class=\"n\">variable</span><span class=\"si\">}</span><span
  class=\"s2\">\"</span>\n                            <span class=\"sa\">f</span><span
  class=\"s2\">\"[</span><span class=\"si\">{</span><span class=\"s1\">','</span><span
  class=\"o\">.</span><span class=\"n\">join</span><span class=\"p\">([</span><span
  class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">p</span><span
  class=\"p\">)</span><span class=\"w\"> </span><span class=\"k\">for</span><span
  class=\"w\"> </span><span class=\"n\">p</span><span class=\"w\"> </span><span class=\"ow\">in</span><span
  class=\"w\"> </span><span class=\"n\">missing_in_posts</span><span class=\"p\">[:</span><span
  class=\"mi\">10</span><span class=\"p\">]])</span><span class=\"si\">}</span><span
  class=\"s2\">...\"</span>\n                        <span class=\"p\">)</span>\n
  \                   <span class=\"k\">else</span><span class=\"p\">:</span>\n                        <span
  class=\"n\">message</span> <span class=\"o\">+=</span> <span class=\"sa\">f</span><span
  class=\"s2\">\"</span><span class=\"se\">\\n</span><span class=\"s2\">paths to posts
  missing </span><span class=\"si\">{</span><span class=\"n\">variable</span><span
  class=\"si\">}</span><span class=\"s2\"> </span><span class=\"si\">{</span><span
  class=\"n\">missing_in_posts</span><span class=\"si\">}</span><span class=\"s2\">\"</span>\n\n
  \                   <span class=\"k\">raise</span> <span class=\"n\">MissingFrontMatter</span><span
  class=\"p\">(</span><span class=\"n\">message</span><span class=\"p\">)</span>\n\n
  \               <span class=\"k\">return</span> <span class=\"n\">posts</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"load_ipython_extension\" style=\"margin:0;padding:.5rem
  1rem;\">load_ipython_extension <em class=\"small\">function</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"load_ipython_extension
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">load_ipython_extension</span><span class=\"p\">(</span><span
  class=\"n\">ipython</span><span class=\"p\">):</span>\n            <span class=\"n\">ipython</span><span
  class=\"o\">.</span><span class=\"n\">user_ns</span><span class=\"p\">[</span><span
  class=\"s2\">\"m\"</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
  class=\"n\">Markata</span><span class=\"p\">()</span>\n            <span class=\"n\">ipython</span><span
  class=\"o\">.</span><span class=\"n\">user_ns</span><span class=\"p\">[</span><span
  class=\"s2\">\"markata\"</span><span class=\"p\">]</span> <span class=\"o\">=</span>
  <span class=\"n\">ipython</span><span class=\"o\">.</span><span class=\"n\">user_ns</span><span
  class=\"p\">[</span><span class=\"s2\">\"m\"</span><span class=\"p\">]</span>\n
  \           <span class=\"n\">ipython</span><span class=\"o\">.</span><span class=\"n\">user_ns</span><span
  class=\"p\">[</span><span class=\"s2\">\"markata\"</span><span class=\"p\">]</span>
  <span class=\"o\">=</span> <span class=\"n\">ipython</span><span class=\"o\">.</span><span
  class=\"n\">user_ns</span><span class=\"p\">[</span><span class=\"s2\">\"m\"</span><span
  class=\"p\">]</span>\n</pre></div>\n\n</pre>\n<p>!! method </p><h2 class=\"admonition-title\"
  id=\"__init__\" style=\"margin:0;padding:.5rem 1rem;\"><strong>init</strong> <em
  class=\"small\">method</em></h2>\n<div class=\"admonition source is-collapsible
  collapsible-open\">\n<p class=\"admonition-title\">\"<strong>init</strong> <em class=\"small\">source</em>\"</p>\n</div>\n<pre
  class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"fm\">__init__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">,</span>
  <span class=\"n\">console</span><span class=\"p\">:</span> <span class=\"n\">Console</span>
  <span class=\"o\">=</span> <span class=\"kc\">None</span><span class=\"p\">,</span>
  <span class=\"n\">config</span><span class=\"o\">=</span><span class=\"kc\">None</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">stages_ran</span> <span class=\"o\">=</span> <span class=\"nb\">set</span><span
  class=\"p\">()</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">threded</span> <span class=\"o\">=</span> <span class=\"kc\">False</span>\n
  \               <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_cache</span>
  <span class=\"o\">=</span> <span class=\"kc\">None</span>\n                <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_precache</span>
  <span class=\"o\">=</span> <span class=\"kc\">None</span>\n                <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">MARKATA_CACHE_DIR</span>
  <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
  class=\"s2\">\".\"</span><span class=\"p\">)</span> <span class=\"o\">/</span> <span
  class=\"s2\">\".markata.cache\"</span>\n                <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">MARKATA_CACHE_DIR</span><span class=\"o\">.</span><span
  class=\"n\">mkdir</span><span class=\"p\">(</span><span class=\"n\">exist_ok</span><span
  class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n                <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_pm</span> <span
  class=\"o\">=</span> <span class=\"n\">pluggy</span><span class=\"o\">.</span><span
  class=\"n\">PluginManager</span><span class=\"p\">(</span><span class=\"s2\">\"markata\"</span><span
  class=\"p\">)</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">_pm</span><span class=\"o\">.</span><span class=\"n\">add_hookspecs</span><span
  class=\"p\">(</span><span class=\"n\">hookspec</span><span class=\"o\">.</span><span
  class=\"n\">MarkataSpecs</span><span class=\"p\">)</span>\n                <span
  class=\"k\">if</span> <span class=\"n\">config</span> <span class=\"ow\">is</span>
  <span class=\"ow\">not</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">config</span> <span class=\"o\">=</span> <span class=\"n\">config</span>\n
  \               <span class=\"k\">with</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">cache</span> <span class=\"k\">as</span> <span
  class=\"n\">cache</span><span class=\"p\">:</span>\n                    <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">init_cache_stats</span> <span class=\"o\">=</span>
  <span class=\"n\">cache</span><span class=\"o\">.</span><span class=\"n\">stats</span><span
  class=\"p\">()</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">registered_attrs</span> <span class=\"o\">=</span> <span class=\"n\">hookspec</span><span
  class=\"o\">.</span><span class=\"n\">registered_attrs</span>\n                <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">post_models</span>
  <span class=\"o\">=</span> <span class=\"p\">[]</span>\n                <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">config_models</span> <span class=\"o\">=</span>
  <span class=\"p\">[]</span>\n                <span class=\"k\">if</span> <span class=\"n\">config</span>
  <span class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n                    <span class=\"n\">raw_hooks</span> <span
  class=\"o\">=</span> <span class=\"n\">config</span>\n                <span class=\"k\">else</span><span
  class=\"p\">:</span>\n                    <span class=\"n\">raw_hooks</span> <span
  class=\"o\">=</span> <span class=\"n\">standard_config</span><span class=\"o\">.</span><span
  class=\"n\">load</span><span class=\"p\">(</span><span class=\"s2\">\"markata\"</span><span
  class=\"p\">)</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">hooks_conf</span> <span class=\"o\">=</span> <span class=\"n\">HooksConfig</span><span
  class=\"o\">.</span><span class=\"n\">parse_obj</span><span class=\"p\">(</span><span
  class=\"n\">raw_hooks</span><span class=\"p\">)</span>\n                <span class=\"k\">try</span><span
  class=\"p\">:</span>\n                    <span class=\"n\">default_index</span>
  <span class=\"o\">=</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">hooks_conf</span><span class=\"o\">.</span><span class=\"n\">hooks</span><span
  class=\"o\">.</span><span class=\"n\">index</span><span class=\"p\">(</span><span
  class=\"s2\">\"default\"</span><span class=\"p\">)</span>\n                    <span
  class=\"n\">hooks</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
  \                       <span class=\"o\">*</span><span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">hooks_conf</span><span class=\"o\">.</span><span
  class=\"n\">hooks</span><span class=\"p\">[:</span><span class=\"n\">default_index</span><span
  class=\"p\">],</span>\n                        <span class=\"o\">*</span><span class=\"n\">DEFAULT_HOOKS</span><span
  class=\"p\">,</span>\n                        <span class=\"o\">*</span><span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">hooks_conf</span><span class=\"o\">.</span><span
  class=\"n\">hooks</span><span class=\"p\">[</span><span class=\"n\">default_index</span>
  <span class=\"o\">+</span> <span class=\"mi\">1</span> <span class=\"p\">:],</span>\n
  \                   <span class=\"p\">]</span>\n                    <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">hooks_conf</span><span class=\"o\">.</span><span
  class=\"n\">hooks</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
  \                       <span class=\"n\">hook</span> <span class=\"k\">for</span>
  <span class=\"n\">hook</span> <span class=\"ow\">in</span> <span class=\"n\">hooks</span>
  <span class=\"k\">if</span> <span class=\"n\">hook</span> <span class=\"ow\">not</span>
  <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">hooks_conf</span><span class=\"o\">.</span><span class=\"n\">disabled_hooks</span>\n
  \                   <span class=\"p\">]</span>\n                <span class=\"k\">except</span>
  <span class=\"ne\">ValueError</span><span class=\"p\">:</span>\n                    <span
  class=\"c1\"># 'default' is not in hooks , do not replace with default_hooks</span>\n
  \                   <span class=\"k\">pass</span>\n\n                <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">_register_hooks</span><span class=\"p\">()</span>\n
  \               <span class=\"k\">if</span> <span class=\"n\">console</span> <span
  class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n                    <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">_console</span> <span class=\"o\">=</span> <span class=\"n\">console</span>\n
  \               <span class=\"n\">atexit</span><span class=\"o\">.</span><span class=\"n\">register</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">teardown</span><span class=\"p\">)</span>\n                <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">precache</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"cache\" style=\"margin:0;padding:.5rem
  1rem;\">cache <em class=\"small\">method</em></h2>\n<div class=\"admonition source
  is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"cache <em class=\"small\">source</em>\"</p>\n</div>\n<pre
  class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">cache</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"n\">Cache</span><span class=\"p\">:</span>\n
  \               <span class=\"c1\"># if self.threded:</span>\n                <span
  class=\"c1\">#     FanoutCache(self.MARKATA_CACHE_DIR, statistics=True)</span>\n
  \               <span class=\"k\">if</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">_cache</span> <span class=\"ow\">is</span>
  <span class=\"ow\">not</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \                   <span class=\"k\">return</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">_cache</span>\n                <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">_cache</span> <span class=\"o\">=</span> <span
  class=\"n\">Cache</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">MARKATA_CACHE_DIR</span><span class=\"p\">,</span>
  <span class=\"n\">statistics</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
  class=\"p\">)</span>\n\n                <span class=\"k\">return</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">_cache</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"precache\" style=\"margin:0;padding:.5rem
  1rem;\">precache <em class=\"small\">method</em></h2>\n<div class=\"admonition source
  is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"precache <em
  class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">precache</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">)</span>
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
  class=\"p\">()}</span>\n                <span class=\"k\">return</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">_precache</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"__getattr__\" style=\"margin:0;padding:.5rem
  1rem;\"><strong>getattr</strong> <em class=\"small\">method</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"<strong>getattr</strong>
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"fm\">__getattr__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">,</span>
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
  class=\"p\">(</span>\n                        <span class=\"p\">[</span><span class=\"n\">attr</span><span
  class=\"p\">[</span><span class=\"s2\">\"lifecycle\"</span><span class=\"p\">]</span>
  <span class=\"k\">for</span> <span class=\"n\">attr</span> <span class=\"ow\">in</span>
  <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">registered_attrs</span><span
  class=\"p\">[</span><span class=\"n\">item</span><span class=\"p\">]],</span>\n
  \                   <span class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">name</span>\n
  \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">run</span><span class=\"p\">(</span><span class=\"n\">stage_to_run_to</span><span
  class=\"p\">)</span>\n                    <span class=\"k\">return</span> <span
  class=\"nb\">getattr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">,</span> <span class=\"n\">item</span><span class=\"p\">)</span>\n                <span
  class=\"k\">else</span><span class=\"p\">:</span>\n                    <span class=\"c1\">#
  Markata does not know what this is, raise</span>\n                    <span class=\"k\">raise</span>
  <span class=\"ne\">AttributeError</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
  class=\"s2\">\"'Markata' object has no attribute '</span><span class=\"si\">{</span><span
  class=\"n\">item</span><span class=\"si\">}</span><span class=\"s2\">'\"</span><span
  class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!! method </p><h2 class=\"admonition-title\"
  id=\"__rich__\" style=\"margin:0;padding:.5rem 1rem;\"><strong>rich</strong> <em
  class=\"small\">method</em></h2>\n<div class=\"admonition source is-collapsible
  collapsible-open\">\n<p class=\"admonition-title\">\"<strong>rich</strong> <em class=\"small\">source</em>\"</p>\n</div>\n<pre
  class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">__rich__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"n\">Table</span><span class=\"p\">:</span>\n
  \               <span class=\"n\">grid</span> <span class=\"o\">=</span> <span class=\"n\">Table</span><span
  class=\"o\">.</span><span class=\"n\">grid</span><span class=\"p\">()</span>\n                <span
  class=\"n\">grid</span><span class=\"o\">.</span><span class=\"n\">add_column</span><span
  class=\"p\">(</span><span class=\"s2\">\"label\"</span><span class=\"p\">)</span>\n
  \               <span class=\"n\">grid</span><span class=\"o\">.</span><span class=\"n\">add_column</span><span
  class=\"p\">(</span><span class=\"s2\">\"value\"</span><span class=\"p\">)</span>\n\n
  \               <span class=\"k\">for</span> <span class=\"n\">label</span><span
  class=\"p\">,</span> <span class=\"n\">value</span> <span class=\"ow\">in</span>
  <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">describe</span><span
  class=\"p\">()</span><span class=\"o\">.</span><span class=\"n\">items</span><span
  class=\"p\">():</span>\n                    <span class=\"n\">grid</span><span class=\"o\">.</span><span
  class=\"n\">add_row</span><span class=\"p\">(</span><span class=\"n\">label</span><span
  class=\"p\">,</span> <span class=\"n\">value</span><span class=\"p\">)</span>\n\n
  \               <span class=\"k\">return</span> <span class=\"n\">grid</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"bust_cache\" style=\"margin:0;padding:.5rem
  1rem;\">bust_cache <em class=\"small\">method</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"bust_cache
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">bust_cache</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"n\">Markata</span><span class=\"p\">:</span>\n
  \               <span class=\"k\">with</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">cache</span> <span class=\"k\">as</span> <span
  class=\"n\">cache</span><span class=\"p\">:</span>\n                    <span class=\"n\">cache</span><span
  class=\"o\">.</span><span class=\"n\">clear</span><span class=\"p\">()</span>\n
  \               <span class=\"k\">return</span> <span class=\"bp\">self</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"configure\" style=\"margin:0;padding:.5rem
  1rem;\">configure <em class=\"small\">method</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"configure
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">configure</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Markata</span><span
  class=\"p\">:</span>\n                <span class=\"n\">sys</span><span class=\"o\">.</span><span
  class=\"n\">path</span><span class=\"o\">.</span><span class=\"n\">append</span><span
  class=\"p\">(</span><span class=\"n\">os</span><span class=\"o\">.</span><span class=\"n\">getcwd</span><span
  class=\"p\">())</span>\n                <span class=\"c1\"># self.config = {**DEFUALT_CONFIG,
  **standard_config.load(\"markata\")}</span>\n                <span class=\"c1\">#
  if isinstance(self.config[\"glob_patterns\"], str):</span>\n                <span
  class=\"c1\">#     self.config[\"glob_patterns\"] = self.config[\"glob_patterns\"].split(\",\")</span>\n
  \               <span class=\"c1\"># elif isinstance(self.config[\"glob_patterns\"],
  list):</span>\n                <span class=\"c1\">#     self.config[\"glob_patterns\"]
  = list(self.config[\"glob_patterns\"])</span>\n                <span class=\"c1\">#
  else:</span>\n                <span class=\"c1\">#     raise TypeError(\"glob_patterns
  must be list or str\")</span>\n                <span class=\"c1\"># self.glob_patterns
  = self.config[\"glob_patterns\"]</span>\n\n                <span class=\"c1\">#
  self.hooks = self.config[\"hooks\"]</span>\n\n                <span class=\"c1\">#
  if \"disabled_hooks\" not in self.config:</span>\n                <span class=\"c1\">#
  \    self.disabled_hooks = [\"\"]</span>\n                <span class=\"c1\"># if
  isinstance(self.config[\"disabled_hooks\"], str):</span>\n                <span
  class=\"c1\">#     self.disabled_hooks = self.config[\"disabled_hooks\"].split(\",\")</span>\n
  \               <span class=\"c1\"># if isinstance(self.config[\"disabled_hooks\"],
  list):</span>\n                <span class=\"c1\">#     self.disabled_hooks = self.config[\"disabled_hooks\"]</span>\n\n
  \               <span class=\"c1\"># if not self.config.get(\"output_dir\", \"markout\").endswith(</span>\n
  \               <span class=\"c1\">#     self.config.get(\"path_prefix\", \"\")</span>\n
  \               <span class=\"c1\"># ):</span>\n                <span class=\"c1\">#
  \    self.config[\"output_dir\"] = (</span>\n                <span class=\"c1\">#
  \        self.config.get(\"output_dir\", \"markout\") +</span>\n                <span
  class=\"c1\">#         \"/\" +</span>\n                <span class=\"c1\">#         self.config.get(\"path_prefix\",
  \"\").rstrip(\"/\")</span>\n                <span class=\"c1\">#     )</span>\n
  \               <span class=\"c1\"># if (</span>\n                <span class=\"c1\">#
  \    len((output_split := self.config.get(\"output_dir\", \"markout\").split(\"/\")))
  &gt;</span>\n                <span class=\"c1\">#     1</span>\n                <span
  class=\"c1\"># ):</span>\n                <span class=\"c1\">#     if \"path_prefix\"
  not in self.config.keys():</span>\n                <span class=\"c1\">#         self.config[\"path_prefix\"]
  = \"/\".join(output_split[1:]) + \"/\"</span>\n                <span class=\"c1\">#
  if not self.config.get(\"path_prefix\", \"\").endswith(\"/\"):</span>\n                <span
  class=\"c1\">#     self.config[\"path_prefix\"] = self.config.get(\"path_prefix\",
  \"\") + \"/\"</span>\n\n                <span class=\"c1\"># self.config[\"output_dir\"]
  = self.config[\"output_dir\"].lstrip(\"/\")</span>\n                <span class=\"c1\">#
  self.config[\"path_prefix\"] = self.config[\"path_prefix\"].lstrip(\"/\")</span>\n\n
  \               <span class=\"k\">try</span><span class=\"p\">:</span>\n                    <span
  class=\"n\">default_index</span> <span class=\"o\">=</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">hooks_conf</span><span class=\"o\">.</span><span
  class=\"n\">hooks</span><span class=\"o\">.</span><span class=\"n\">index</span><span
  class=\"p\">(</span><span class=\"s2\">\"default\"</span><span class=\"p\">)</span>\n
  \                   <span class=\"n\">hooks</span> <span class=\"o\">=</span> <span
  class=\"p\">[</span>\n                        <span class=\"o\">*</span><span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">hooks_conf</span><span class=\"o\">.</span><span
  class=\"n\">hooks</span><span class=\"p\">[:</span><span class=\"n\">default_index</span><span
  class=\"p\">],</span>\n                        <span class=\"o\">*</span><span class=\"n\">DEFAULT_HOOKS</span><span
  class=\"p\">,</span>\n                        <span class=\"o\">*</span><span class=\"bp\">self</span><span
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
  class=\"c1\"># 'default' is not in hooks , do not replace with default_hooks</span>\n
  \                   <span class=\"k\">pass</span>\n\n                <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">_pm</span> <span class=\"o\">=</span> <span
  class=\"n\">pluggy</span><span class=\"o\">.</span><span class=\"n\">PluginManager</span><span
  class=\"p\">(</span><span class=\"s2\">\"markata\"</span><span class=\"p\">)</span>\n
  \               <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_pm</span><span
  class=\"o\">.</span><span class=\"n\">add_hookspecs</span><span class=\"p\">(</span><span
  class=\"n\">hookspec</span><span class=\"o\">.</span><span class=\"n\">MarkataSpecs</span><span
  class=\"p\">)</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">_register_hooks</span><span class=\"p\">()</span>\n\n                <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_pm</span><span
  class=\"o\">.</span><span class=\"n\">hook</span><span class=\"o\">.</span><span
  class=\"n\">configure</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
  class=\"o\">=</span><span class=\"bp\">self</span><span class=\"p\">)</span>\n                <span
  class=\"k\">return</span> <span class=\"bp\">self</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"get_plugin_config\" style=\"margin:0;padding:.5rem
  1rem;\">get_plugin_config <em class=\"small\">method</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"get_plugin_config
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">get_plugin_config</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">,</span> <span class=\"n\">path_or_name</span><span class=\"p\">:</span>
  <span class=\"nb\">str</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"n\">Dict</span><span class=\"p\">:</span>\n                <span class=\"n\">key</span>
  <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
  class=\"n\">path_or_name</span><span class=\"p\">)</span><span class=\"o\">.</span><span
  class=\"n\">stem</span>\n\n                <span class=\"n\">config</span> <span
  class=\"o\">=</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">get</span><span
  class=\"p\">(</span><span class=\"n\">key</span><span class=\"p\">,</span> <span
  class=\"p\">{})</span>\n\n                <span class=\"k\">if</span> <span class=\"ow\">not</span>
  <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">config</span><span
  class=\"p\">,</span> <span class=\"nb\">dict</span><span class=\"p\">):</span>\n
  \                   <span class=\"k\">raise</span> <span class=\"ne\">TypeError</span><span
  class=\"p\">(</span><span class=\"s2\">\"must use dict\"</span><span class=\"p\">)</span>\n
  \               <span class=\"k\">if</span> <span class=\"s2\">\"cache_expire\"</span>
  <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">keys</span><span class=\"p\">():</span>\n
  \                   <span class=\"n\">config</span><span class=\"p\">[</span><span
  class=\"s2\">\"cache_expire\"</span><span class=\"p\">]</span> <span class=\"o\">=</span>
  <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"p\">[</span><span class=\"s2\">\"default_cache_expire\"</span><span class=\"p\">]</span>\n
  \               <span class=\"k\">if</span> <span class=\"s2\">\"config_key\"</span>
  <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">keys</span><span class=\"p\">():</span>\n
  \                   <span class=\"n\">config</span><span class=\"p\">[</span><span
  class=\"s2\">\"config_key\"</span><span class=\"p\">]</span> <span class=\"o\">=</span>
  <span class=\"n\">key</span>\n                <span class=\"k\">return</span> <span
  class=\"n\">config</span>\n</pre></div>\n\n</pre>\n<p>!! method </p><h2 class=\"admonition-title\"
  id=\"get_config\" style=\"margin:0;padding:.5rem 1rem;\">get_config <em class=\"small\">method</em></h2>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"get_config
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">get_config</span><span class=\"p\">(</span>\n                <span
  class=\"bp\">self</span><span class=\"p\">,</span>\n                <span class=\"n\">key</span><span
  class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">,</span>\n                <span
  class=\"n\">default</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
  <span class=\"o\">=</span> <span class=\"s2\">\"\"</span><span class=\"p\">,</span>\n
  \               <span class=\"n\">warn</span><span class=\"p\">:</span> <span class=\"nb\">bool</span>
  <span class=\"o\">=</span> <span class=\"kc\">True</span><span class=\"p\">,</span>\n
  \               <span class=\"n\">suggested</span><span class=\"p\">:</span> <span
  class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
  class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span><span
  class=\"p\">,</span>\n            <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"n\">Any</span><span class=\"p\">:</span>\n                <span class=\"k\">if</span>
  <span class=\"n\">key</span> <span class=\"ow\">in</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
  class=\"n\">keys</span><span class=\"p\">():</span>\n                    <span class=\"k\">return</span>
  <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"p\">[</span><span class=\"n\">key</span><span class=\"p\">]</span>\n                <span
  class=\"k\">else</span><span class=\"p\">:</span>\n                    <span class=\"k\">if</span>
  <span class=\"n\">suggested</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n                        <span class=\"n\">suggested</span>
  <span class=\"o\">=</span> <span class=\"n\">textwrap</span><span class=\"o\">.</span><span
  class=\"n\">dedent</span><span class=\"p\">(</span>\n                            <span
  class=\"sa\">f</span><span class=\"s2\">\"\"\"</span>\n<span class=\"s2\">                            [markata]</span>\n<span
  class=\"s2\">                            </span><span class=\"si\">{</span><span
  class=\"n\">key</span><span class=\"si\">}</span><span class=\"s2\"> = '</span><span
  class=\"si\">{</span><span class=\"n\">default</span><span class=\"si\">}</span><span
  class=\"s2\">'</span>\n<span class=\"s2\">                            \"\"\"</span>\n
  \                       <span class=\"p\">)</span>\n                    <span class=\"k\">if</span>
  <span class=\"n\">warn</span><span class=\"p\">:</span>\n                        <span
  class=\"n\">logger</span><span class=\"o\">.</span><span class=\"n\">warning</span><span
  class=\"p\">(</span>\n                            <span class=\"n\">textwrap</span><span
  class=\"o\">.</span><span class=\"n\">dedent</span><span class=\"p\">(</span>\n
  \                               <span class=\"sa\">f</span><span class=\"s2\">\"\"\"</span>\n<span
  class=\"s2\">                                Warning </span><span class=\"si\">{</span><span
  class=\"n\">key</span><span class=\"si\">}</span><span class=\"s2\"> is not set
  in markata config, sitemap will</span>\n<span class=\"s2\">                                be
  missing root site_name</span>\n<span class=\"s2\">                                to
  resolve this open your markata.toml and add</span>\n\n<span class=\"s2\">                                </span><span
  class=\"si\">{</span><span class=\"n\">suggested</span><span class=\"si\">}</span>\n\n<span
  class=\"s2\">                                \"\"\"</span>\n                            <span
  class=\"p\">),</span>\n                        <span class=\"p\">)</span>\n                <span
  class=\"k\">return</span> <span class=\"n\">default</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"make_hash\" style=\"margin:0;padding:.5rem
  1rem;\">make_hash <em class=\"small\">method</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"make_hash
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">make_hash</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">,</span> <span class=\"o\">*</span><span class=\"n\">keys</span><span
  class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">)</span> <span
  class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n
  \               <span class=\"n\">str_keys</span> <span class=\"o\">=</span> <span
  class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">(</span><span
  class=\"n\">key</span><span class=\"p\">)</span> <span class=\"k\">for</span> <span
  class=\"n\">key</span> <span class=\"ow\">in</span> <span class=\"n\">keys</span><span
  class=\"p\">]</span>\n                <span class=\"k\">return</span> <span class=\"n\">hashlib</span><span
  class=\"o\">.</span><span class=\"n\">md5</span><span class=\"p\">(</span><span
  class=\"s2\">\"\"</span><span class=\"o\">.</span><span class=\"n\">join</span><span
  class=\"p\">(</span><span class=\"n\">str_keys</span><span class=\"p\">)</span><span
  class=\"o\">.</span><span class=\"n\">encode</span><span class=\"p\">(</span><span
  class=\"s2\">\"utf-8\"</span><span class=\"p\">))</span><span class=\"o\">.</span><span
  class=\"n\">hexdigest</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"content_dir_hash\" style=\"margin:0;padding:.5rem
  1rem;\">content_dir_hash <em class=\"small\">method</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"content_dir_hash
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">content_dir_hash</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n
  \               <span class=\"n\">hashes</span> <span class=\"o\">=</span> <span
  class=\"p\">[</span>\n                    <span class=\"n\">dirhash</span><span
  class=\"p\">(</span><span class=\"nb\">dir</span><span class=\"p\">)</span>\n                    <span
  class=\"k\">for</span> <span class=\"nb\">dir</span> <span class=\"ow\">in</span>
  <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">content_directories</span>\n
  \                   <span class=\"k\">if</span> <span class=\"nb\">dir</span><span
  class=\"o\">.</span><span class=\"n\">absolute</span><span class=\"p\">()</span>
  <span class=\"o\">!=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
  class=\"s2\">\".\"</span><span class=\"p\">)</span><span class=\"o\">.</span><span
  class=\"n\">absolute</span><span class=\"p\">()</span>\n                <span class=\"p\">]</span>\n
  \               <span class=\"k\">return</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">make_hash</span><span class=\"p\">(</span><span
  class=\"o\">*</span><span class=\"n\">hashes</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"console\" style=\"margin:0;padding:.5rem
  1rem;\">console <em class=\"small\">method</em></h2>\n<div class=\"admonition source
  is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"console <em
  class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">console</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"n\">Console</span><span class=\"p\">:</span>\n
  \               <span class=\"k\">try</span><span class=\"p\">:</span>\n                    <span
  class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">_console</span>\n                <span class=\"k\">except</span> <span
  class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n                    <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_console</span>
  <span class=\"o\">=</span> <span class=\"n\">Console</span><span class=\"p\">()</span>\n
  \                   <span class=\"k\">return</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">_console</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"describe\" style=\"margin:0;padding:.5rem
  1rem;\">describe <em class=\"small\">method</em></h2>\n<div class=\"admonition source
  is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"describe <em
  class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">describe</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"nb\">dict</span><span class=\"p\">[</span><span
  class=\"nb\">str</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
  class=\"p\">]:</span>\n                <span class=\"k\">return</span> <span class=\"p\">{</span><span
  class=\"s2\">\"version\"</span><span class=\"p\">:</span> <span class=\"n\">__version__</span><span
  class=\"p\">}</span>\n</pre></div>\n\n</pre>\n<p>!! method </p><h2 class=\"admonition-title\"
  id=\"_to_dict\" style=\"margin:0;padding:.5rem 1rem;\">_to_dict <em class=\"small\">method</em></h2>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"_to_dict
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">_to_dict</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"nb\">dict</span><span class=\"p\">[</span><span
  class=\"nb\">str</span><span class=\"p\">,</span> <span class=\"n\">Iterable</span><span
  class=\"p\">]:</span>\n                <span class=\"k\">return</span> <span class=\"p\">{</span><span
  class=\"s2\">\"config\"</span><span class=\"p\">:</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">,</span> <span
  class=\"s2\">\"articles\"</span><span class=\"p\">:</span> <span class=\"p\">[</span><span
  class=\"n\">a</span><span class=\"o\">.</span><span class=\"n\">to_dict</span><span
  class=\"p\">()</span> <span class=\"k\">for</span> <span class=\"n\">a</span> <span
  class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">articles</span><span class=\"p\">]}</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"to_dict\" style=\"margin:0;padding:.5rem
  1rem;\">to_dict <em class=\"small\">method</em></h2>\n<div class=\"admonition source
  is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"to_dict <em
  class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">to_dict</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"nb\">dict</span><span class=\"p\">:</span>\n
  \               <span class=\"k\">return</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">_to_dict</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"to_json\" style=\"margin:0;padding:.5rem
  1rem;\">to_json <em class=\"small\">method</em></h2>\n<div class=\"admonition source
  is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"to_json <em
  class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">to_json</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n
  \               <span class=\"kn\">import</span> <span class=\"nn\">json</span>\n\n
  \               <span class=\"k\">return</span> <span class=\"n\">json</span><span
  class=\"o\">.</span><span class=\"n\">dumps</span><span class=\"p\">(</span><span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">to_dict</span><span
  class=\"p\">(),</span> <span class=\"n\">indent</span><span class=\"o\">=</span><span
  class=\"mi\">4</span><span class=\"p\">,</span> <span class=\"n\">sort_keys</span><span
  class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
  class=\"n\">default</span><span class=\"o\">=</span><span class=\"nb\">str</span><span
  class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!! method </p><h2 class=\"admonition-title\"
  id=\"_register_hooks\" style=\"margin:0;padding:.5rem 1rem;\">_register_hooks <em
  class=\"small\">method</em></h2>\n<div class=\"admonition source is-collapsible
  collapsible-open\">\n<p class=\"admonition-title\">\"_register_hooks <em class=\"small\">source</em>\"</p>\n</div>\n<pre
  class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">_register_hooks</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">)</span>
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
  class=\"n\">e</span><span class=\"p\">:</span>\n                        <span class=\"c1\">#
  class style plugins</span>\n                        <span class=\"k\">if</span>
  <span class=\"s2\">\".\"</span> <span class=\"ow\">in</span> <span class=\"n\">hook</span><span
  class=\"p\">:</span>\n                            <span class=\"k\">try</span><span
  class=\"p\">:</span>\n                                <span class=\"n\">mod</span>
  <span class=\"o\">=</span> <span class=\"n\">importlib</span><span class=\"o\">.</span><span
  class=\"n\">import_module</span><span class=\"p\">(</span><span class=\"s2\">\".\"</span><span
  class=\"o\">.</span><span class=\"n\">join</span><span class=\"p\">(</span><span
  class=\"n\">hook</span><span class=\"o\">.</span><span class=\"n\">split</span><span
  class=\"p\">(</span><span class=\"s2\">\".\"</span><span class=\"p\">)[:</span><span
  class=\"o\">-</span><span class=\"mi\">1</span><span class=\"p\">]))</span>\n                                <span
  class=\"n\">plugin</span> <span class=\"o\">=</span> <span class=\"nb\">getattr</span><span
  class=\"p\">(</span><span class=\"n\">mod</span><span class=\"p\">,</span> <span
  class=\"n\">hook</span><span class=\"o\">.</span><span class=\"n\">split</span><span
  class=\"p\">(</span><span class=\"s2\">\".\"</span><span class=\"p\">)[</span><span
  class=\"o\">-</span><span class=\"mi\">1</span><span class=\"p\">])</span>\n                            <span
  class=\"k\">except</span> <span class=\"ne\">ModuleNotFoundError</span> <span class=\"k\">as</span>
  <span class=\"n\">e</span><span class=\"p\">:</span>\n                                <span
  class=\"k\">raise</span> <span class=\"ne\">ModuleNotFoundError</span><span class=\"p\">(</span>\n
  \                                   <span class=\"sa\">f</span><span class=\"s2\">\"module
  </span><span class=\"si\">{</span><span class=\"n\">hook</span><span class=\"si\">}</span><span
  class=\"s2\"> not found</span><span class=\"se\">\\n</span><span class=\"si\">{</span><span
  class=\"n\">sys</span><span class=\"o\">.</span><span class=\"n\">path</span><span
  class=\"si\">}</span><span class=\"s2\">\"</span>\n                                <span
  class=\"p\">)</span> <span class=\"kn\">from</span> <span class=\"nn\">e</span>\n
  \                       <span class=\"k\">else</span><span class=\"p\">:</span>\n
  \                           <span class=\"k\">raise</span> <span class=\"n\">e</span>\n\n
  \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">_pm</span><span class=\"o\">.</span><span class=\"n\">register</span><span
  class=\"p\">(</span><span class=\"n\">plugin</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"__iter__\" style=\"margin:0;padding:.5rem
  1rem;\"><strong>iter</strong> <em class=\"small\">method</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"<strong>iter</strong>
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"fm\">__iter__</span><span class=\"p\">(</span>\n                <span
  class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span
  class=\"p\">,</span> <span class=\"n\">description</span><span class=\"p\">:</span>
  <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">\"working...\"</span>\n
  \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Iterable</span><span
  class=\"p\">[</span><span class=\"s2\">\"Markata.Post\"</span><span class=\"p\">]:</span>\n
  \               <span class=\"n\">articles</span><span class=\"p\">:</span> <span
  class=\"n\">Iterable</span><span class=\"p\">[</span><span class=\"n\">Markata</span><span
  class=\"o\">.</span><span class=\"n\">Post</span><span class=\"p\">]</span> <span
  class=\"o\">=</span> <span class=\"n\">track</span><span class=\"p\">(</span>\n
  \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">articles</span><span class=\"p\">,</span>\n                    <span
  class=\"n\">description</span><span class=\"o\">=</span><span class=\"n\">description</span><span
  class=\"p\">,</span>\n                    <span class=\"n\">transient</span><span
  class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n                    <span
  class=\"n\">console</span><span class=\"o\">=</span><span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">console</span><span class=\"p\">,</span>\n
  \               <span class=\"p\">)</span>\n                <span class=\"k\">return</span>
  <span class=\"n\">articles</span>\n</pre></div>\n\n</pre>\n<p>!! method </p><h2
  class=\"admonition-title\" id=\"iter_articles\" style=\"margin:0;padding:.5rem 1rem;\">iter_articles
  <em class=\"small\">method</em></h2>\n<div class=\"admonition source is-collapsible
  collapsible-open\">\n<p class=\"admonition-title\">\"iter_articles <em class=\"small\">source</em>\"</p>\n</div>\n<pre
  class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">iter_articles</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">,</span>
  <span class=\"n\">description</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Iterable</span><span
  class=\"p\">[</span><span class=\"n\">Markata</span><span class=\"o\">.</span><span
  class=\"n\">Post</span><span class=\"p\">]:</span>\n                <span class=\"n\">articles</span><span
  class=\"p\">:</span> <span class=\"n\">Iterable</span><span class=\"p\">[</span><span
  class=\"n\">Markata</span><span class=\"o\">.</span><span class=\"n\">Post</span><span
  class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"n\">track</span><span
  class=\"p\">(</span>\n                    <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">articles</span><span class=\"p\">,</span>\n                    <span
  class=\"n\">description</span><span class=\"o\">=</span><span class=\"n\">description</span><span
  class=\"p\">,</span>\n                    <span class=\"n\">transient</span><span
  class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n                    <span
  class=\"n\">console</span><span class=\"o\">=</span><span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">console</span><span class=\"p\">,</span>\n
  \               <span class=\"p\">)</span>\n                <span class=\"k\">return</span>
  <span class=\"n\">articles</span>\n</pre></div>\n\n</pre>\n<p>!! method </p><h2
  class=\"admonition-title\" id=\"teardown\" style=\"margin:0;padding:.5rem 1rem;\">teardown
  <em class=\"small\">method</em></h2>\ngive special access to the teardown lifecycle
  method\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"teardown
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">teardown</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"n\">Markata</span><span class=\"p\">:</span>\n<span
  class=\"w\">                </span><span class=\"sd\">\"\"\"give special access
  to the teardown lifecycle method\"\"\"</span>\n                <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">_pm</span><span class=\"o\">.</span><span
  class=\"n\">hook</span><span class=\"o\">.</span><span class=\"n\">teardown</span><span
  class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
  class=\"bp\">self</span><span class=\"p\">)</span>\n                <span class=\"k\">return</span>
  <span class=\"bp\">self</span>\n</pre></div>\n\n</pre>\n<p>!! method </p><h2 class=\"admonition-title\"
  id=\"run\" style=\"margin:0;padding:.5rem 1rem;\">run <em class=\"small\">method</em></h2>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"run
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">run</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">,</span>
  <span class=\"n\">lifecycle</span><span class=\"p\">:</span> <span class=\"n\">LifeCycle</span>
  <span class=\"o\">=</span> <span class=\"kc\">None</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"n\">Markata</span><span class=\"p\">:</span>\n
  \               <span class=\"k\">if</span> <span class=\"n\">lifecycle</span> <span
  class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \                   <span class=\"n\">lifecycle</span> <span class=\"o\">=</span>
  <span class=\"nb\">max</span><span class=\"p\">(</span><span class=\"n\">LifeCycle</span><span
  class=\"o\">.</span><span class=\"n\">_member_map_</span><span class=\"o\">.</span><span
  class=\"n\">values</span><span class=\"p\">())</span>\n\n                <span class=\"k\">if</span>
  <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">lifecycle</span><span
  class=\"p\">,</span> <span class=\"nb\">str</span><span class=\"p\">):</span>\n
  \                   <span class=\"n\">lifecycle</span> <span class=\"o\">=</span>
  <span class=\"n\">LifeCycle</span><span class=\"p\">[</span><span class=\"n\">lifecycle</span><span
  class=\"p\">]</span>\n\n                <span class=\"n\">stages_to_run</span> <span
  class=\"o\">=</span> <span class=\"p\">[</span>\n                    <span class=\"n\">m</span>\n
  \                   <span class=\"k\">for</span> <span class=\"n\">m</span> <span
  class=\"ow\">in</span> <span class=\"n\">LifeCycle</span><span class=\"o\">.</span><span
  class=\"n\">_member_map_</span>\n                    <span class=\"k\">if</span>
  <span class=\"p\">(</span><span class=\"n\">LifeCycle</span><span class=\"p\">[</span><span
  class=\"n\">m</span><span class=\"p\">]</span> <span class=\"o\">&lt;=</span> <span
  class=\"n\">lifecycle</span><span class=\"p\">)</span> <span class=\"ow\">and</span>
  <span class=\"p\">(</span><span class=\"n\">m</span> <span class=\"ow\">not</span>
  <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">stages_ran</span><span class=\"p\">)</span>\n                <span class=\"p\">]</span>\n\n
  \               <span class=\"k\">if</span> <span class=\"ow\">not</span> <span
  class=\"n\">stages_to_run</span><span class=\"p\">:</span>\n                    <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">console</span><span
  class=\"o\">.</span><span class=\"n\">log</span><span class=\"p\">(</span><span
  class=\"sa\">f</span><span class=\"s2\">\"</span><span class=\"si\">{</span><span
  class=\"n\">lifecycle</span><span class=\"o\">.</span><span class=\"n\">name</span><span
  class=\"si\">}</span><span class=\"s2\"> already ran\"</span><span class=\"p\">)</span>\n
  \                   <span class=\"k\">return</span> <span class=\"bp\">self</span>\n\n
  \               <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">console</span><span
  class=\"o\">.</span><span class=\"n\">log</span><span class=\"p\">(</span><span
  class=\"sa\">f</span><span class=\"s2\">\"running </span><span class=\"si\">{</span><span
  class=\"n\">stages_to_run</span><span class=\"si\">}</span><span class=\"s2\">\"</span><span
  class=\"p\">)</span>\n                <span class=\"k\">for</span> <span class=\"n\">stage</span>
  <span class=\"ow\">in</span> <span class=\"n\">stages_to_run</span><span class=\"p\">:</span>\n
  \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
  class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">\"</span><span
  class=\"si\">{</span><span class=\"n\">stage</span><span class=\"si\">}</span><span
  class=\"s2\"> running\"</span><span class=\"p\">)</span>\n                    <span
  class=\"nb\">getattr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">_pm</span><span class=\"o\">.</span><span
  class=\"n\">hook</span><span class=\"p\">,</span> <span class=\"n\">stage</span><span
  class=\"p\">)(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
  class=\"bp\">self</span><span class=\"p\">)</span>\n                    <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">stages_ran</span><span class=\"o\">.</span><span
  class=\"n\">add</span><span class=\"p\">(</span><span class=\"n\">stage</span><span
  class=\"p\">)</span>\n                    <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
  class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">\"</span><span
  class=\"si\">{</span><span class=\"n\">stage</span><span class=\"si\">}</span><span
  class=\"s2\"> complete\"</span><span class=\"p\">)</span>\n\n                <span
  class=\"k\">with</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">cache</span> <span class=\"k\">as</span> <span class=\"n\">cache</span><span
  class=\"p\">:</span>\n                    <span class=\"n\">hits</span><span class=\"p\">,</span>
  <span class=\"n\">misses</span> <span class=\"o\">=</span> <span class=\"n\">cache</span><span
  class=\"o\">.</span><span class=\"n\">stats</span><span class=\"p\">()</span>\n\n
  \               <span class=\"k\">if</span> <span class=\"n\">hits</span> <span
  class=\"o\">+</span> <span class=\"n\">misses</span> <span class=\"o\">&gt;</span>
  <span class=\"mi\">0</span><span class=\"p\">:</span>\n                    <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">console</span><span
  class=\"o\">.</span><span class=\"n\">log</span><span class=\"p\">(</span>\n                        <span
  class=\"sa\">f</span><span class=\"s2\">\"lifetime cache hit rate </span><span class=\"si\">{</span><span
  class=\"nb\">round</span><span class=\"p\">(</span><span class=\"n\">hits</span><span
  class=\"o\">/</span><span class=\"w\"> </span><span class=\"p\">(</span><span class=\"n\">hits</span><span
  class=\"w\"> </span><span class=\"o\">+</span><span class=\"w\"> </span><span class=\"n\">misses</span><span
  class=\"p\">)</span><span class=\"o\">*</span><span class=\"mi\">100</span><span
  class=\"p\">,</span><span class=\"w\"> </span><span class=\"mi\">2</span><span class=\"p\">)</span><span
  class=\"si\">}</span><span class=\"s2\">%\"</span><span class=\"p\">,</span>\n                    <span
  class=\"p\">)</span>\n\n                <span class=\"k\">if</span> <span class=\"n\">misses</span>
  <span class=\"o\">&gt;</span> <span class=\"mi\">0</span><span class=\"p\">:</span>\n
  \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
  class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">\"lifetime cache
  hits/misses </span><span class=\"si\">{</span><span class=\"n\">hits</span><span
  class=\"si\">}</span><span class=\"s2\">/</span><span class=\"si\">{</span><span
  class=\"n\">misses</span><span class=\"si\">}</span><span class=\"s2\">\"</span><span
  class=\"p\">)</span>\n\n                <span class=\"n\">hits</span> <span class=\"o\">-=</span>
  <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">init_cache_stats</span><span
  class=\"p\">[</span><span class=\"mi\">0</span><span class=\"p\">]</span>\n                <span
  class=\"n\">misses</span> <span class=\"o\">-=</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">init_cache_stats</span><span class=\"p\">[</span><span
  class=\"mi\">1</span><span class=\"p\">]</span>\n\n                <span class=\"k\">if</span>
  <span class=\"n\">hits</span> <span class=\"o\">+</span> <span class=\"n\">misses</span>
  <span class=\"o\">&gt;</span> <span class=\"mi\">0</span><span class=\"p\">:</span>\n
  \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
  class=\"p\">(</span>\n                        <span class=\"sa\">f</span><span class=\"s2\">\"run
  cache hit rate </span><span class=\"si\">{</span><span class=\"nb\">round</span><span
  class=\"p\">(</span><span class=\"n\">hits</span><span class=\"o\">/</span><span
  class=\"w\"> </span><span class=\"p\">(</span><span class=\"n\">hits</span><span
  class=\"w\"> </span><span class=\"o\">+</span><span class=\"w\"> </span><span class=\"n\">misses</span><span
  class=\"p\">)</span><span class=\"o\">*</span><span class=\"mi\">100</span><span
  class=\"p\">,</span><span class=\"w\"> </span><span class=\"mi\">2</span><span class=\"p\">)</span><span
  class=\"si\">}</span><span class=\"s2\">%\"</span><span class=\"p\">,</span>\n                    <span
  class=\"p\">)</span>\n\n                <span class=\"k\">if</span> <span class=\"n\">misses</span>
  <span class=\"o\">&gt;</span> <span class=\"mi\">0</span><span class=\"p\">:</span>\n
  \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
  class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">\"run cache hits/misses
  </span><span class=\"si\">{</span><span class=\"n\">hits</span><span class=\"si\">}</span><span
  class=\"s2\">/</span><span class=\"si\">{</span><span class=\"n\">misses</span><span
  class=\"si\">}</span><span class=\"s2\">\"</span><span class=\"p\">)</span>\n\n
  \               <span class=\"k\">return</span> <span class=\"bp\">self</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"filter\" style=\"margin:0;padding:.5rem
  1rem;\">filter <em class=\"small\">method</em></h2>\n<div class=\"admonition source
  is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"filter <em class=\"small\">source</em>\"</p>\n</div>\n<pre
  class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">filter</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">,</span>
  <span class=\"nb\">filter</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">list</span><span
  class=\"p\">:</span>\n                <span class=\"k\">def</span> <span class=\"nf\">evalr</span><span
  class=\"p\">(</span><span class=\"n\">a</span><span class=\"p\">:</span> <span class=\"n\">Markata</span><span
  class=\"o\">.</span><span class=\"n\">Post</span><span class=\"p\">)</span> <span
  class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span class=\"p\">:</span>\n
  \                   <span class=\"k\">try</span><span class=\"p\">:</span>\n                        <span
  class=\"k\">return</span> <span class=\"nb\">eval</span><span class=\"p\">(</span>\n
  \                           <span class=\"nb\">filter</span><span class=\"p\">,</span>\n
  \                           <span class=\"p\">{</span><span class=\"o\">**</span><span
  class=\"n\">a</span><span class=\"o\">.</span><span class=\"n\">to_dict</span><span
  class=\"p\">(),</span> <span class=\"s2\">\"timedelta\"</span><span class=\"p\">:</span>
  <span class=\"n\">timedelta</span><span class=\"p\">,</span> <span class=\"s2\">\"post\"</span><span
  class=\"p\">:</span> <span class=\"n\">a</span><span class=\"p\">,</span> <span
  class=\"s2\">\"m\"</span><span class=\"p\">:</span> <span class=\"bp\">self</span><span
  class=\"p\">},</span>\n                            <span class=\"p\">{},</span>\n
  \                       <span class=\"p\">)</span>\n                    <span class=\"k\">except</span>
  <span class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n                        <span
  class=\"k\">return</span> <span class=\"nb\">eval</span><span class=\"p\">(</span>\n
  \                           <span class=\"nb\">filter</span><span class=\"p\">,</span>\n
  \                           <span class=\"p\">{</span><span class=\"o\">**</span><span
  class=\"n\">a</span><span class=\"o\">.</span><span class=\"n\">to_dict</span><span
  class=\"p\">(),</span> <span class=\"s2\">\"timedelta\"</span><span class=\"p\">:</span>
  <span class=\"n\">timedelta</span><span class=\"p\">,</span> <span class=\"s2\">\"post\"</span><span
  class=\"p\">:</span> <span class=\"n\">a</span><span class=\"p\">,</span> <span
  class=\"s2\">\"m\"</span><span class=\"p\">:</span> <span class=\"bp\">self</span><span
  class=\"p\">},</span>\n                            <span class=\"p\">{},</span>\n
  \                       <span class=\"p\">)</span>\n\n                <span class=\"k\">return</span>
  <span class=\"p\">[</span><span class=\"n\">a</span> <span class=\"k\">for</span>
  <span class=\"n\">a</span> <span class=\"ow\">in</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">articles</span> <span class=\"k\">if</span>
  <span class=\"n\">evalr</span><span class=\"p\">(</span><span class=\"n\">a</span><span
  class=\"p\">)]</span>\n</pre></div>\n\n</pre>\n<p>!! method </p><h2 class=\"admonition-title\"
  id=\"map\" style=\"margin:0;padding:.5rem 1rem;\">map <em class=\"small\">method</em></h2>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"map
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">map</span><span class=\"p\">(</span>\n                <span class=\"bp\">self</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">,</span>\n
  \               <span class=\"n\">func</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
  <span class=\"o\">=</span> <span class=\"s2\">\"title\"</span><span class=\"p\">,</span>\n
  \               <span class=\"nb\">filter</span><span class=\"p\">:</span> <span
  class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">\"True\"</span><span
  class=\"p\">,</span>\n                <span class=\"n\">sort</span><span class=\"p\">:</span>
  <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">\"True\"</span><span
  class=\"p\">,</span>\n                <span class=\"n\">reverse</span><span class=\"p\">:</span>
  <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"kc\">True</span><span
  class=\"p\">,</span>\n                <span class=\"o\">*</span><span class=\"n\">args</span><span
  class=\"p\">:</span> <span class=\"nb\">tuple</span><span class=\"p\">,</span>\n
  \               <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
  class=\"p\">:</span> <span class=\"nb\">dict</span><span class=\"p\">,</span>\n
  \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">list</span><span
  class=\"p\">:</span>\n                <span class=\"kn\">import</span> <span class=\"nn\">copy</span>\n\n
  \               <span class=\"k\">def</span> <span class=\"nf\">try_sort</span><span
  class=\"p\">(</span><span class=\"n\">a</span><span class=\"p\">:</span> <span class=\"n\">Any</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">int</span><span
  class=\"p\">:</span>\n                    <span class=\"k\">if</span> <span class=\"s2\">\"datetime\"</span>
  <span class=\"ow\">in</span> <span class=\"n\">sort</span><span class=\"o\">.</span><span
  class=\"n\">lower</span><span class=\"p\">():</span>\n                        <span
  class=\"k\">return</span> <span class=\"n\">a</span><span class=\"o\">.</span><span
  class=\"n\">get</span><span class=\"p\">(</span><span class=\"n\">sort</span><span
  class=\"p\">,</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
  class=\"n\">datetime</span><span class=\"p\">(</span><span class=\"mi\">1970</span><span
  class=\"p\">,</span> <span class=\"mi\">1</span><span class=\"p\">,</span> <span
  class=\"mi\">1</span><span class=\"p\">))</span>\n\n                    <span class=\"k\">if</span>
  <span class=\"s2\">\"date\"</span> <span class=\"ow\">in</span> <span class=\"n\">sort</span><span
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
  class=\"p\">(),</span> <span class=\"p\">{})</span>\n                    <span class=\"k\">except</span>
  <span class=\"ne\">NameError</span><span class=\"p\">:</span>\n                        <span
  class=\"k\">return</span> <span class=\"o\">-</span><span class=\"mi\">1</span>\n
  \                   <span class=\"k\">return</span> <span class=\"n\">value</span>\n
  \                   <span class=\"k\">try</span><span class=\"p\">:</span>\n                        <span
  class=\"k\">return</span> <span class=\"nb\">int</span><span class=\"p\">(</span><span
  class=\"n\">value</span><span class=\"p\">)</span>\n                    <span class=\"k\">except</span>
  <span class=\"ne\">TypeError</span><span class=\"p\">:</span>\n                        <span
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
  <span class=\"o\">-</span><span class=\"mi\">1</span>\n\n                <span class=\"n\">articles</span>
  <span class=\"o\">=</span> <span class=\"n\">copy</span><span class=\"o\">.</span><span
  class=\"n\">copy</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">articles</span><span class=\"p\">)</span>\n
  \               <span class=\"n\">articles</span><span class=\"o\">.</span><span
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
  class=\"p\">(),</span> <span class=\"s2\">\"timedelta\"</span><span class=\"p\">:</span>
  <span class=\"n\">timedelta</span><span class=\"p\">,</span> <span class=\"s2\">\"post\"</span><span
  class=\"p\">:</span> <span class=\"n\">a</span><span class=\"p\">,</span> <span
  class=\"s2\">\"m\"</span><span class=\"p\">:</span> <span class=\"bp\">self</span><span
  class=\"p\">},</span>\n                            <span class=\"p\">{},</span>\n
  \                       <span class=\"p\">)</span>\n                        <span
  class=\"k\">for</span> <span class=\"n\">a</span> <span class=\"ow\">in</span> <span
  class=\"n\">articles</span>\n                        <span class=\"k\">if</span>
  <span class=\"nb\">eval</span><span class=\"p\">(</span>\n                            <span
  class=\"nb\">filter</span><span class=\"p\">,</span>\n                            <span
  class=\"p\">{</span><span class=\"o\">**</span><span class=\"n\">a</span><span class=\"o\">.</span><span
  class=\"n\">to_dict</span><span class=\"p\">(),</span> <span class=\"s2\">\"timedelta\"</span><span
  class=\"p\">:</span> <span class=\"n\">timedelta</span><span class=\"p\">,</span>
  <span class=\"s2\">\"post\"</span><span class=\"p\">:</span> <span class=\"n\">a</span><span
  class=\"p\">,</span> <span class=\"s2\">\"m\"</span><span class=\"p\">:</span> <span
  class=\"bp\">self</span><span class=\"p\">},</span>\n                            <span
  class=\"p\">{},</span>\n                        <span class=\"p\">)</span>\n                    <span
  class=\"p\">]</span>\n\n                <span class=\"k\">except</span> <span class=\"ne\">NameError</span>
  <span class=\"k\">as</span> <span class=\"n\">e</span><span class=\"p\">:</span>\n
  \                   <span class=\"n\">variable</span> <span class=\"o\">=</span>
  <span class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">e</span><span
  class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">split</span><span
  class=\"p\">(</span><span class=\"s2\">\"'\"</span><span class=\"p\">)[</span><span
  class=\"mi\">1</span><span class=\"p\">]</span>\n\n                    <span class=\"n\">missing_in_posts</span>
  <span class=\"o\">=</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">map</span><span class=\"p\">(</span>\n                        <span
  class=\"s2\">\"path\"</span><span class=\"p\">,</span>\n                        <span
  class=\"nb\">filter</span><span class=\"o\">=</span><span class=\"sa\">f</span><span
  class=\"s1\">'\"</span><span class=\"si\">{</span><span class=\"n\">variable</span><span
  class=\"si\">}</span><span class=\"s1\">\" not in post.keys()'</span><span class=\"p\">,</span>\n
  \                   <span class=\"p\">)</span>\n                    <span class=\"n\">message</span>
  <span class=\"o\">=</span> <span class=\"p\">(</span>\n                        <span
  class=\"sa\">f</span><span class=\"s2\">\"variable: '</span><span class=\"si\">{</span><span
  class=\"n\">variable</span><span class=\"si\">}</span><span class=\"s2\">' is missing
  in </span><span class=\"si\">{</span><span class=\"nb\">len</span><span class=\"p\">(</span><span
  class=\"n\">missing_in_posts</span><span class=\"p\">)</span><span class=\"si\">}</span><span
  class=\"s2\"> posts\"</span>\n                    <span class=\"p\">)</span>\n                    <span
  class=\"k\">if</span> <span class=\"nb\">len</span><span class=\"p\">(</span><span
  class=\"n\">missing_in_posts</span><span class=\"p\">)</span> <span class=\"o\">&gt;</span>
  <span class=\"mi\">10</span><span class=\"p\">:</span>\n                        <span
  class=\"n\">message</span> <span class=\"o\">+=</span> <span class=\"p\">(</span>\n
  \                           <span class=\"sa\">f</span><span class=\"s2\">\"</span><span
  class=\"se\">\\n</span><span class=\"s2\">first 10 paths to posts missing </span><span
  class=\"si\">{</span><span class=\"n\">variable</span><span class=\"si\">}</span><span
  class=\"s2\">\"</span>\n                            <span class=\"sa\">f</span><span
  class=\"s2\">\"[</span><span class=\"si\">{</span><span class=\"s1\">','</span><span
  class=\"o\">.</span><span class=\"n\">join</span><span class=\"p\">([</span><span
  class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">p</span><span
  class=\"p\">)</span><span class=\"w\"> </span><span class=\"k\">for</span><span
  class=\"w\"> </span><span class=\"n\">p</span><span class=\"w\"> </span><span class=\"ow\">in</span><span
  class=\"w\"> </span><span class=\"n\">missing_in_posts</span><span class=\"p\">[:</span><span
  class=\"mi\">10</span><span class=\"p\">]])</span><span class=\"si\">}</span><span
  class=\"s2\">...\"</span>\n                        <span class=\"p\">)</span>\n
  \                   <span class=\"k\">else</span><span class=\"p\">:</span>\n                        <span
  class=\"n\">message</span> <span class=\"o\">+=</span> <span class=\"sa\">f</span><span
  class=\"s2\">\"</span><span class=\"se\">\\n</span><span class=\"s2\">paths to posts
  missing </span><span class=\"si\">{</span><span class=\"n\">variable</span><span
  class=\"si\">}</span><span class=\"s2\"> </span><span class=\"si\">{</span><span
  class=\"n\">missing_in_posts</span><span class=\"si\">}</span><span class=\"s2\">\"</span>\n\n
  \                   <span class=\"k\">raise</span> <span class=\"n\">MissingFrontMatter</span><span
  class=\"p\">(</span><span class=\"n\">message</span><span class=\"p\">)</span>\n\n
  \               <span class=\"k\">return</span> <span class=\"n\">posts</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"evalr\" style=\"margin:0;padding:.5rem
  1rem;\">evalr <em class=\"small\">function</em></h2>\n<div class=\"admonition source
  is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"evalr <em class=\"small\">source</em>\"</p>\n</div>\n<pre
  class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
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
  class=\"p\">(),</span> <span class=\"s2\">\"timedelta\"</span><span class=\"p\">:</span>
  <span class=\"n\">timedelta</span><span class=\"p\">,</span> <span class=\"s2\">\"post\"</span><span
  class=\"p\">:</span> <span class=\"n\">a</span><span class=\"p\">,</span> <span
  class=\"s2\">\"m\"</span><span class=\"p\">:</span> <span class=\"bp\">self</span><span
  class=\"p\">},</span>\n                            <span class=\"p\">{},</span>\n
  \                       <span class=\"p\">)</span>\n                    <span class=\"k\">except</span>
  <span class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n                        <span
  class=\"k\">return</span> <span class=\"nb\">eval</span><span class=\"p\">(</span>\n
  \                           <span class=\"nb\">filter</span><span class=\"p\">,</span>\n
  \                           <span class=\"p\">{</span><span class=\"o\">**</span><span
  class=\"n\">a</span><span class=\"o\">.</span><span class=\"n\">to_dict</span><span
  class=\"p\">(),</span> <span class=\"s2\">\"timedelta\"</span><span class=\"p\">:</span>
  <span class=\"n\">timedelta</span><span class=\"p\">,</span> <span class=\"s2\">\"post\"</span><span
  class=\"p\">:</span> <span class=\"n\">a</span><span class=\"p\">,</span> <span
  class=\"s2\">\"m\"</span><span class=\"p\">:</span> <span class=\"bp\">self</span><span
  class=\"p\">},</span>\n                            <span class=\"p\">{},</span>\n
  \                       <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"try_sort\" style=\"margin:0;padding:.5rem
  1rem;\">try_sort <em class=\"small\">function</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"try_sort
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">try_sort</span><span class=\"p\">(</span><span class=\"n\">a</span><span
  class=\"p\">:</span> <span class=\"n\">Any</span><span class=\"p\">)</span> <span
  class=\"o\">-&gt;</span> <span class=\"nb\">int</span><span class=\"p\">:</span>\n
  \                   <span class=\"k\">if</span> <span class=\"s2\">\"datetime\"</span>
  <span class=\"ow\">in</span> <span class=\"n\">sort</span><span class=\"o\">.</span><span
  class=\"n\">lower</span><span class=\"p\">():</span>\n                        <span
  class=\"k\">return</span> <span class=\"n\">a</span><span class=\"o\">.</span><span
  class=\"n\">get</span><span class=\"p\">(</span><span class=\"n\">sort</span><span
  class=\"p\">,</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
  class=\"n\">datetime</span><span class=\"p\">(</span><span class=\"mi\">1970</span><span
  class=\"p\">,</span> <span class=\"mi\">1</span><span class=\"p\">,</span> <span
  class=\"mi\">1</span><span class=\"p\">))</span>\n\n                    <span class=\"k\">if</span>
  <span class=\"s2\">\"date\"</span> <span class=\"ow\">in</span> <span class=\"n\">sort</span><span
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
  class=\"p\">(),</span> <span class=\"p\">{})</span>\n                    <span class=\"k\">except</span>
  <span class=\"ne\">NameError</span><span class=\"p\">:</span>\n                        <span
  class=\"k\">return</span> <span class=\"o\">-</span><span class=\"mi\">1</span>\n
  \                   <span class=\"k\">return</span> <span class=\"n\">value</span>\n
  \                   <span class=\"k\">try</span><span class=\"p\">:</span>\n                        <span
  class=\"k\">return</span> <span class=\"nb\">int</span><span class=\"p\">(</span><span
  class=\"n\">value</span><span class=\"p\">)</span>\n                    <span class=\"k\">except</span>
  <span class=\"ne\">TypeError</span><span class=\"p\">:</span>\n                        <span
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
  <span class=\"o\">-</span><span class=\"mi\">1</span>\n</pre></div>\n\n</pre>\n</main>\n<footer>\xA9
  2024</footer>\n</body></html>"
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
                    self.run(stage_to_run_to)
                    return getattr(self, item)
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
                str_keys = [str(key) for key in keys]
                return hashlib.md5("".join(str_keys).encode("utf-8")).hexdigest()

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
                    transient=True,
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
                    self.run(stage_to_run_to)
                    return getattr(self, item)
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
                str_keys = [str(key) for key in keys]
                return hashlib.md5("".join(str_keys).encode("utf-8")).hexdigest()
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
                    transient=True,
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
