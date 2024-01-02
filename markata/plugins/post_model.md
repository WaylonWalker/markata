---
content: "None\n\n\n!! class <h2 id='Post' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>Post <em class='small'>class</em></h2>\n\n???+ source \"Post <em class='small'>source</em>\"\n\n```python\n\n
  \       class Post(pydantic.BaseModel):\n            markata: Any = None\n            path:
  Path\n            slug: Optional[str] = None\n            href: Optional[str] =
  None\n            published: bool = False\n            description: Optional[str]
  = None\n            content: str = None\n            # date: Union[datetime.date,
  str]=None\n            date: Optional[Union[datetime.date, str]] = None\n            #
  pydantic.Field(\n            # default_factory=lambda: datetime.date.min\n            #
  )\n            date_time: Optional[datetime.datetime] = None\n            today:
  datetime.date = pydantic.Field(default_factory=datetime.date.today)\n            now:
  datetime.datetime = pydantic.Field(default_factory=datetime.datetime.utcnow)\n            load_time:
  float = 0\n            profile: Optional[str] = None\n            title: str = None\n
  \           model_config = ConfigDict(\n                validate_assignment=True,\n
  \               arbitrary_types_allowed=True,\n                extra=\"allow\",\n
  \           )\n\n            def __repr_args__(self: \"Post\") -> \"ReprArgs\":\n
  \               return [\n                    (key, value)\n                    for
  key, value in self.__dict__.items()\n                    if key in self.markata.config.post_model.repr_include\n
  \               ]\n\n            @property\n            def metadata(self: \"Post\")
  -> Dict:\n                \"for backwards compatability\"\n                return
  self.__dict__\n\n            def to_dict(self: \"Post\") -> Dict:\n                \"for
  backwards compatability\"\n                return self.__dict__\n\n            def
  __getitem__(self: \"Post\", item: str) -> Any:\n                \"for backwards
  compatability\"\n                return getattr(self, item)\n\n            def __setitem__(self:
  \"Post\", key: str, item: Any) -> None:\n                \"for backwards compatability\"\n
  \               setattr(self, key, item)\n\n            def get(self: \"Post\",
  item: str, default: Any) -> Any:\n                \"for backwards compatability\"\n
  \               return getattr(self, item, default)\n\n            def keys(self:
  \"Post\") -> List[str]:\n                \"for backwards compatability\"\n                return
  self.__dict__.keys()\n\n            # def json(\n            #     self: \"Post\",\n
  \           #     include: Iterable = None,\n            #     all: bool = False,\n
  \           #     **kwargs,\n            # ) -> str:\n            #     \"\"\"\n
  \           #     override function to give a default include value that will include\n
  \           #     user configured includes.\n            #     \"\"\"\n            #
  \    if all:\n            #         return pydantic.create_model(\"Post\", **self)(**self).json(\n
  \           #             **kwargs,\n            #         )\n            #     if
  include:\n            #         return pydantic.create_model(\"Post\", **self)(**self).json(\n
  \           #             include=include,\n            #             **kwargs,\n
  \           #         )\n            #     return pydantic.create_model(\"Post\",
  **self)(**self).json(\n            #         include={i: True for i in self.markata.config.post_model.include},\n
  \           #         **kwargs,\n            #     )\n\n            def yaml(self:
  \"Post\") -> str:\n                \"\"\"\n                dump model to yaml\n
  \               \"\"\"\n                import yaml\n\n                return yaml.dump(\n
  \                   self.dict(\n                        include={i: True for i in
  self.markata.config.post_model.include}\n                    ),\n                    Dumper=yaml.CDumper,\n
  \               )\n\n            def markdown(self: \"Post\") -> str:\n                \"\"\"\n
  \               dump model to markdown\n                \"\"\"\n\n                import
  yaml\n\n                frontmatter = yaml.dump(\n                    self.dict(\n
  \                       include={\n                            i: True\n                            for
  i in [\n                                _i\n                                for
  _i in self.markata.config.post_model.include\n                                if
  _i != \"content\"\n                            ]\n                        }\n                    ),\n
  \                   Dumper=yaml.CDumper,\n                )\n                post
  = \"---\\n\"\n                post += frontmatter\n                post += \"---\\n\\n\"\n\n
  \               if self.content:\n                    post += self.content\n                return
  post\n\n            @classmethod\n            def parse_file(cls, markata, path:
  Union[Path, str], **kwargs) -> \"Post\":\n                if isinstance(path, Path):\n
  \                   if path.suffix in [\".md\", \".markdown\"]:\n                        return
  cls.parse_markdown(markata=markata, path=path, **kwargs)\n                elif isinstance(path,
  str):\n                    if path.endswith(\".md\") or path.endswith(\".markdown\"):\n
  \                       return cls.parse_markdown(markata=markata, path=path, **kwargs)\n
  \               return super(Post, cls).parse_file(path, **kwargs)\n\n            @classmethod\n
  \           def parse_markdown(cls, markata, path: Union[Path, str], **kwargs) ->
  \"Post\":\n                if isinstance(path, str):\n                    path =
  Path(path)\n                text = path.read_text()\n                try:\n                    _,
  fm, *content = text.split(\"---\\n\")\n                    content = \"---\\n\".join(content)\n
  \                   try:\n                        fm = yaml.load(fm, Loader=yaml.CBaseLoader)\n
  \                   except yaml.YAMLError:\n                        fm = {}\n                except
  ValueError:\n                    fm = {}\n                    content = text\n                if
  fm is None or isinstance(fm, str):\n                    fm = {}\n\n                post_args
  = {\n                    \"markata\": markata,\n                    \"path\": path,\n
  \                   \"content\": content,\n                    **fm,\n                }\n\n
  \               return markata.Post(**post_args)\n\n            def dumps(self):\n
  \               \"\"\"\n                dumps raw article back out\n                \"\"\"\n
  \               return f\"---\\n{self.yaml()}\\n\\n---\\n\\n{self.content}\"\n\n
  \           @pydantic.validator(\"slug\", pre=True, always=True)\n            def
  default_slug(cls, v, *, values):\n                return v or slugify(str(values[\"path\"].stem))\n\n
  \           @pydantic.validator(\"slug\", pre=True, always=True)\n            def
  index_slug_is_empty(cls, v, *, values):\n                if v == \"index\":\n                    return
  \"\"\n                return v\n\n            @pydantic.validator(\"href\", pre=True,
  always=True)\n            def default_href(cls, v, *, values):\n                if
  v:\n                    return v\n                return f\"/{values['slug'].strip('/')}/\".replace(\"//\",
  \"/\")\n\n            @pydantic.validator(\"title\", pre=True, always=True)\n            def
  title_title(cls, v, *, values):\n                title = v or Path(values[\"path\"]).stem.replace(\"-\",
  \" \")\n                return title.title()\n\n            @pydantic.validator(\"date_time\",
  pre=True, always=True)\n            def dateparser_datetime(cls, v, *, values):\n
  \               if isinstance(v, str):\n                    d = dateparser.parse(v)\n
  \                   if d is None:\n                        raise ValueError(f'\"{v}\"
  is not a valid date')\n                return v\n\n            @pydantic.validator(\"date_time\",
  pre=True, always=True)\n            def date_is_datetime(cls, v, *, values):\n                if
  v is None and \"date\" not in values:\n                    values[\"markata\"].console.log(f\"{values['path']}
  has no date\")\n                    return datetime.datetime.now()\n                if
  v is None and values[\"date\"] is None:\n                    values[\"markata\"].console.log(f\"{values['path']}
  has no date\")\n                    return datetime.datetime.now()\n                if
  isinstance(v, datetime.datetime):\n                    return v\n                if
  isinstance(values[\"date\"], datetime.datetime):\n                    return values[\"date\"]\n
  \               if isinstance(v, datetime.date):\n                    return datetime.datetime.combine(v,
  datetime.time.min)\n                if isinstance(values[\"date\"], datetime.date):\n
  \                   return datetime.datetime.combine(values[\"date\"], datetime.time.min)\n
  \               return v\n\n            @pydantic.validator(\"date_time\", pre=True,
  always=True)\n            def mindate_time(cls, v, *, values):\n                if
  v is None and \"date\" not in values:\n                    values[\"markata\"].console.log(f\"{values['path']}
  has no date\")\n                    return datetime.datetime.min\n                if
  values[\"date\"] is None:\n                    values[\"markata\"].console.log(f\"{values['path']}
  has no date\")\n                    return datetime.datetime.min\n                if
  isinstance(v, datetime.datetime):\n                    return v\n                if
  isinstance(values[\"date\"], datetime.datetime):\n                    return values[\"date\"]\n
  \               if isinstance(v, datetime.date):\n                    return datetime.datetime.combine(v,
  datetime.time.min)\n                if isinstance(values[\"date\"], datetime.date):\n
  \                   return datetime.datetime.combine(values[\"date\"], datetime.time.min)\n
  \               return v\n\n            @pydantic.validator(\"date\", pre=True,
  always=True)\n            def dateparser_date(cls, v, *, values):\n                if
  v is None:\n                    return datetime.date.min\n                if isinstance(v,
  str):\n                    d = cls.markata.precache.get(v)\n                    if
  d is not None:\n                        return d\n                    d = dateparser.parse(v)\n
  \                   if d is None:\n                        raise ValueError(f'\"{v}\"
  is not a valid date')\n                    d = d.date()\n                    with
  cls.markata.cache as cache:\n                        cache.add(v, d)\n                    return
  d\n                return v\n```\n\n\n!! class <h2 id='PostModelConfig' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>PostModelConfig <em class='small'>class</em></h2>\n
  \   Configuration for the Post model\n???+ source \"PostModelConfig <em class='small'>source</em>\"\n\n```python\n\n
  \       class PostModelConfig(pydantic.BaseModel):\n            \"Configuration
  for the Post model\"\n\n            def __init__(self, **data) -> None:\n                \"\"\"\n\n
  \               include: post attributes to include by default in Post\n                model
  serialization.\n                repr_include: post attributes to include by default
  in Post\n                repr.  If `repr_include` is None, it will default to\n
  \               `include`, but it is likely that you want less in the repr\n                than
  serialized output.\n\n                example:\n\n                ``` toml title='markata.toml'\n
  \               [markata.post_model]\n                include = ['date', 'description',
  'published',\n                    'slug', 'title', 'content', 'html']\n                repr_include
  = ['date', 'description', 'published', 'slug', 'title']\n                ```\n                \"\"\"\n
  \               super().__init__(**data)\n\n            include: List[str] = [\n
  \               \"date\",\n                \"description\",\n                \"published\",\n
  \               \"slug\",\n                \"title\",\n                \"content\",\n
  \               \"html\",\n            ]\n            repr_include: Optional[List[str]]
  = [\n                \"date\",\n                \"description\",\n                \"published\",\n
  \               \"slug\",\n                \"title\",\n            ]\n\n            @pydantic.validator(\"repr_include\",
  pre=True, always=True)\n            def repr_include_validator(cls, v, *, values):\n
  \               if v:\n                    return v\n                return values.get(\"include\",
  None)\n```\n\n\n!! class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>Config <em class='small'>class</em></h2>\n\n???+ source \"Config <em class='small'>source</em>\"\n\n```python\n\n
  \       class Config(pydantic.BaseModel):\n            post_model: PostModelConfig
  = pydantic.Field(default_factory=PostModelConfig)\n```\n\n\n!! function <h2 id='post_model'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>post_model <em class='small'>function</em></h2>\n\n???+
  source \"post_model <em class='small'>source</em>\"\n\n```python\n\n        def
  post_model(markata: \"Markata\") -> None:\n            markata.post_models.append(Post)\n```\n\n\n!!
  function <h2 id='config_model' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>config_model <em class='small'>function</em></h2>\n\n???+ source \"config_model
  <em class='small'>source</em>\"\n\n```python\n\n        def config_model(markata:
  \"Markata\") -> None:\n            markata.config_models.append(Config)\n```\n\n\n!!
  class <h2 id='PostFactory' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>PostFactory <em class='small'>class</em></h2>\n\n???+ source \"PostFactory
  <em class='small'>source</em>\"\n\n```python\n\n        class PostFactory(ModelFactory):\n
  \           __model__ = Post\n```\n\n\n!! method <h2 id='__repr_args__' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>__repr_args__ <em class='small'>method</em></h2>\n\n???+
  source \"__repr_args__ <em class='small'>source</em>\"\n\n```python\n\n        def
  __repr_args__(self: \"Post\") -> \"ReprArgs\":\n                return [\n                    (key,
  value)\n                    for key, value in self.__dict__.items()\n                    if
  key in self.markata.config.post_model.repr_include\n                ]\n```\n\n\n!!
  method <h2 id='metadata' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>metadata <em class='small'>method</em></h2>\n    for backwards compatability\n???+
  source \"metadata <em class='small'>source</em>\"\n\n```python\n\n        def metadata(self:
  \"Post\") -> Dict:\n                \"for backwards compatability\"\n                return
  self.__dict__\n```\n\n\n!! method <h2 id='to_dict' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>to_dict <em class='small'>method</em></h2>\n    for backwards compatability\n???+
  source \"to_dict <em class='small'>source</em>\"\n\n```python\n\n        def to_dict(self:
  \"Post\") -> Dict:\n                \"for backwards compatability\"\n                return
  self.__dict__\n```\n\n\n!! method <h2 id='__getitem__' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>__getitem__ <em class='small'>method</em></h2>\n
  \   for backwards compatability\n???+ source \"__getitem__ <em class='small'>source</em>\"\n\n```python\n\n
  \       def __getitem__(self: \"Post\", item: str) -> Any:\n                \"for
  backwards compatability\"\n                return getattr(self, item)\n```\n\n\n!!
  method <h2 id='__setitem__' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>__setitem__ <em class='small'>method</em></h2>\n    for backwards compatability\n???+
  source \"__setitem__ <em class='small'>source</em>\"\n\n```python\n\n        def
  __setitem__(self: \"Post\", key: str, item: Any) -> None:\n                \"for
  backwards compatability\"\n                setattr(self, key, item)\n```\n\n\n!!
  method <h2 id='get' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get
  <em class='small'>method</em></h2>\n    for backwards compatability\n???+ source
  \"get <em class='small'>source</em>\"\n\n```python\n\n        def get(self: \"Post\",
  item: str, default: Any) -> Any:\n                \"for backwards compatability\"\n
  \               return getattr(self, item, default)\n```\n\n\n!! method <h2 id='keys'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>keys <em class='small'>method</em></h2>\n
  \   for backwards compatability\n???+ source \"keys <em class='small'>source</em>\"\n\n```python\n\n
  \       def keys(self: \"Post\") -> List[str]:\n                \"for backwards
  compatability\"\n                return self.__dict__.keys()\n```\n\n\n!! method
  <h2 id='yaml' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>yaml
  <em class='small'>method</em></h2>\n    dump model to yaml\n???+ source \"yaml <em
  class='small'>source</em>\"\n\n```python\n\n        def yaml(self: \"Post\") ->
  str:\n                \"\"\"\n                dump model to yaml\n                \"\"\"\n
  \               import yaml\n\n                return yaml.dump(\n                    self.dict(\n
  \                       include={i: True for i in self.markata.config.post_model.include}\n
  \                   ),\n                    Dumper=yaml.CDumper,\n                )\n```\n\n\n!!
  method <h2 id='markdown' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>markdown <em class='small'>method</em></h2>\n    dump model to markdown\n???+
  source \"markdown <em class='small'>source</em>\"\n\n```python\n\n        def markdown(self:
  \"Post\") -> str:\n                \"\"\"\n                dump model to markdown\n
  \               \"\"\"\n\n                import yaml\n\n                frontmatter
  = yaml.dump(\n                    self.dict(\n                        include={\n
  \                           i: True\n                            for i in [\n                                _i\n
  \                               for _i in self.markata.config.post_model.include\n
  \                               if _i != \"content\"\n                            ]\n
  \                       }\n                    ),\n                    Dumper=yaml.CDumper,\n
  \               )\n                post = \"---\\n\"\n                post += frontmatter\n
  \               post += \"---\\n\\n\"\n\n                if self.content:\n                    post
  += self.content\n                return post\n```\n\n\n!! method <h2 id='parse_file'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>parse_file <em class='small'>method</em></h2>\n\n???+
  source \"parse_file <em class='small'>source</em>\"\n\n```python\n\n        def
  parse_file(cls, markata, path: Union[Path, str], **kwargs) -> \"Post\":\n                if
  isinstance(path, Path):\n                    if path.suffix in [\".md\", \".markdown\"]:\n
  \                       return cls.parse_markdown(markata=markata, path=path, **kwargs)\n
  \               elif isinstance(path, str):\n                    if path.endswith(\".md\")
  or path.endswith(\".markdown\"):\n                        return cls.parse_markdown(markata=markata,
  path=path, **kwargs)\n                return super(Post, cls).parse_file(path, **kwargs)\n```\n\n\n!!
  method <h2 id='parse_markdown' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>parse_markdown <em class='small'>method</em></h2>\n\n???+ source \"parse_markdown
  <em class='small'>source</em>\"\n\n```python\n\n        def parse_markdown(cls,
  markata, path: Union[Path, str], **kwargs) -> \"Post\":\n                if isinstance(path,
  str):\n                    path = Path(path)\n                text = path.read_text()\n
  \               try:\n                    _, fm, *content = text.split(\"---\\n\")\n
  \                   content = \"---\\n\".join(content)\n                    try:\n
  \                       fm = yaml.load(fm, Loader=yaml.CBaseLoader)\n                    except
  yaml.YAMLError:\n                        fm = {}\n                except ValueError:\n
  \                   fm = {}\n                    content = text\n                if
  fm is None or isinstance(fm, str):\n                    fm = {}\n\n                post_args
  = {\n                    \"markata\": markata,\n                    \"path\": path,\n
  \                   \"content\": content,\n                    **fm,\n                }\n\n
  \               return markata.Post(**post_args)\n```\n\n\n!! method <h2 id='dumps'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>dumps <em class='small'>method</em></h2>\n
  \   dumps raw article back out\n???+ source \"dumps <em class='small'>source</em>\"\n\n```python\n\n
  \       def dumps(self):\n                \"\"\"\n                dumps raw article
  back out\n                \"\"\"\n                return f\"---\\n{self.yaml()}\\n\\n---\\n\\n{self.content}\"\n```\n\n\n!!
  method <h2 id='default_slug' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>default_slug <em class='small'>method</em></h2>\n\n???+ source \"default_slug
  <em class='small'>source</em>\"\n\n```python\n\n        def default_slug(cls, v,
  *, values):\n                return v or slugify(str(values[\"path\"].stem))\n```\n\n\n!!
  method <h2 id='index_slug_is_empty' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>index_slug_is_empty <em class='small'>method</em></h2>\n\n???+ source \"index_slug_is_empty
  <em class='small'>source</em>\"\n\n```python\n\n        def index_slug_is_empty(cls,
  v, *, values):\n                if v == \"index\":\n                    return \"\"\n
  \               return v\n```\n\n\n!! method <h2 id='default_href' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>default_href <em class='small'>method</em></h2>\n\n???+
  source \"default_href <em class='small'>source</em>\"\n\n```python\n\n        def
  default_href(cls, v, *, values):\n                if v:\n                    return
  v\n                return f\"/{values['slug'].strip('/')}/\".replace(\"//\", \"/\")\n```\n\n\n!!
  method <h2 id='title_title' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>title_title <em class='small'>method</em></h2>\n\n???+ source \"title_title
  <em class='small'>source</em>\"\n\n```python\n\n        def title_title(cls, v,
  *, values):\n                title = v or Path(values[\"path\"]).stem.replace(\"-\",
  \" \")\n                return title.title()\n```\n\n\n!! method <h2 id='dateparser_datetime'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>dateparser_datetime
  <em class='small'>method</em></h2>\n\n???+ source \"dateparser_datetime <em class='small'>source</em>\"\n\n```python\n\n
  \       def dateparser_datetime(cls, v, *, values):\n                if isinstance(v,
  str):\n                    d = dateparser.parse(v)\n                    if d is
  None:\n                        raise ValueError(f'\"{v}\" is not a valid date')\n
  \               return v\n```\n\n\n!! method <h2 id='date_is_datetime' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>date_is_datetime <em class='small'>method</em></h2>\n\n???+
  source \"date_is_datetime <em class='small'>source</em>\"\n\n```python\n\n        def
  date_is_datetime(cls, v, *, values):\n                if v is None and \"date\"
  not in values:\n                    values[\"markata\"].console.log(f\"{values['path']}
  has no date\")\n                    return datetime.datetime.now()\n                if
  v is None and values[\"date\"] is None:\n                    values[\"markata\"].console.log(f\"{values['path']}
  has no date\")\n                    return datetime.datetime.now()\n                if
  isinstance(v, datetime.datetime):\n                    return v\n                if
  isinstance(values[\"date\"], datetime.datetime):\n                    return values[\"date\"]\n
  \               if isinstance(v, datetime.date):\n                    return datetime.datetime.combine(v,
  datetime.time.min)\n                if isinstance(values[\"date\"], datetime.date):\n
  \                   return datetime.datetime.combine(values[\"date\"], datetime.time.min)\n
  \               return v\n```\n\n\n!! method <h2 id='mindate_time' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>mindate_time <em class='small'>method</em></h2>\n\n???+
  source \"mindate_time <em class='small'>source</em>\"\n\n```python\n\n        def
  mindate_time(cls, v, *, values):\n                if v is None and \"date\" not
  in values:\n                    values[\"markata\"].console.log(f\"{values['path']}
  has no date\")\n                    return datetime.datetime.min\n                if
  values[\"date\"] is None:\n                    values[\"markata\"].console.log(f\"{values['path']}
  has no date\")\n                    return datetime.datetime.min\n                if
  isinstance(v, datetime.datetime):\n                    return v\n                if
  isinstance(values[\"date\"], datetime.datetime):\n                    return values[\"date\"]\n
  \               if isinstance(v, datetime.date):\n                    return datetime.datetime.combine(v,
  datetime.time.min)\n                if isinstance(values[\"date\"], datetime.date):\n
  \                   return datetime.datetime.combine(values[\"date\"], datetime.time.min)\n
  \               return v\n```\n\n\n!! method <h2 id='dateparser_date' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>dateparser_date <em class='small'>method</em></h2>\n\n???+
  source \"dateparser_date <em class='small'>source</em>\"\n\n```python\n\n        def
  dateparser_date(cls, v, *, values):\n                if v is None:\n                    return
  datetime.date.min\n                if isinstance(v, str):\n                    d
  = cls.markata.precache.get(v)\n                    if d is not None:\n                        return
  d\n                    d = dateparser.parse(v)\n                    if d is None:\n
  \                       raise ValueError(f'\"{v}\" is not a valid date')\n                    d
  = d.date()\n                    with cls.markata.cache as cache:\n                        cache.add(v,
  d)\n                    return d\n                return v\n```\n\n\n!! method <h2
  id='__init__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__init__
  <em class='small'>method</em></h2>\n    include: post attributes to include by default
  in Post\n    model serialization.\n    repr_include: post attributes to include
  by default in Post\n    repr.  If `repr_include` is None, it will default to\n    `include`,
  but it is likely that you want less in the repr\n    than serialized output.\n\n
  \   example:\n\n    ``` toml title='markata.toml'\n    [markata.post_model]\n    include
  = ['date', 'description', 'published',\n        'slug', 'title', 'content', 'html']\n
  \   repr_include = ['date', 'description', 'published', 'slug', 'title']\n    ```\n???+
  source \"__init__ <em class='small'>source</em>\"\n\n```python\n\n        def __init__(self,
  **data) -> None:\n                \"\"\"\n\n                include: post attributes
  to include by default in Post\n                model serialization.\n                repr_include:
  post attributes to include by default in Post\n                repr.  If `repr_include`
  is None, it will default to\n                `include`, but it is likely that you
  want less in the repr\n                than serialized output.\n\n                example:\n\n
  \               ``` toml title='markata.toml'\n                [markata.post_model]\n
  \               include = ['date', 'description', 'published',\n                    'slug',
  'title', 'content', 'html']\n                repr_include = ['date', 'description',
  'published', 'slug', 'title']\n                ```\n                \"\"\"\n                super().__init__(**data)\n```\n\n\n!!
  method <h2 id='repr_include_validator' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>repr_include_validator <em class='small'>method</em></h2>\n\n???+ source
  \"repr_include_validator <em class='small'>source</em>\"\n\n```python\n\n        def
  repr_include_validator(cls, v, *, values):\n                if v:\n                    return
  v\n                return values.get(\"include\", None)\n```\n"
date: 0001-01-01
description: 'None ! ???+ source  ! ! ???+ source  ! ???+ source  ! ???+ source  !
  ???+ source  ! ???+ source  ! ! ! ! ! ! ! ! ! ???+ source  ! ???+ source  ! ! ???+
  source  '
html: "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n<title>Post_Model.Py</title>\n<meta
  charset=\"utf-8\"/>\n<meta content=\"width=device-width, initial-scale=1\" name=\"viewport\"/>\n<meta
  content=\"None ! ???+ source  ! ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
  source  ! ???+ source  ! ! ! ! ! ! ! ! ! ???+ source  ! ???+ source  ! ! ???+ source
  \ \" name=\"description\"/>\n<link href=\"/static/favicon.ico\" rel=\"icon\" type=\"image/png\"/>\n<script>\n
  \       function setTheme(theme) {\n            document.documentElement.setAttribute(\"data-theme\",
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
  name=\"og:type\" property=\"og:type\"/><meta content=\"None ! ???+ source  ! ! ???+
  source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ! ! ! ! !
  ! ! ! ???+ source  ! ???+ source  ! ! ???+ source  \" name=\"description\" property=\"description\"/><meta
  content=\"None ! ???+ source  ! ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
  source  ! ???+ source  ! ! ! ! ! ! ! ! ! ???+ source  ! ???+ source  ! ! ???+ source
  \ \" name=\"og:description\" property=\"og:description\"/><meta content=\"None !
  ???+ source  ! ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
  source  ! ! ! ! ! ! ! ! ! ???+ source  ! ???+ source  ! ! ???+ source  \" name=\"twitter:description\"
  property=\"twitter:description\"/><meta content=\"Post_Model.Py | Markata\" name=\"og:title\"
  property=\"og:title\"/><meta content=\"Post_Model.Py | Markata\" name=\"twitter:title\"
  property=\"twitter:title\"/><meta content=\"https://markata.dev//markata/plugins/post-model-og.png\"
  name=\"og:image\" property=\"og:image\"/><meta content=\"https://markata.dev//markata/plugins/post-model-og.png\"
  name=\"twitter:image\" property=\"twitter:image\"/><meta content=\"1600\" name=\"og:image:width\"
  property=\"og:image:width\"/><meta content=\"900\" name=\"og:image:width\" property=\"og:image:width\"/><meta
  content=\"summary_large_image\" name=\"twitter:card\" property=\"twitter:card\"/><meta
  content=\"Markata\" name=\"og:site_name\" property=\"og:site_name\"/><meta content=\"@_waylonwalker\"
  name=\"twitter:creator\" property=\"twitter:creator\"/><meta content=\"Post_Model.Py\"
  name=\"title\" property=\"title\"/><meta content=\"markata 0.8.0.dev6\" name=\"generator\"
  property=\"generator\"/><link href=\"https://markata.dev//markata/plugins/post-model/\"
  rel=\"canonical\"/><meta content=\"https://markata.dev//markata/plugins/post-model/\"
  name=\"og:url\" property=\"og:url\"/></head>\n<body><nav>\n<a href=\"/\">markata</a>\n<a
  href=\"https://github.com/WaylonWalker/markata\">GitHub</a>\n</nav>\n<div>\n<label
  class=\"theme-switch\" for=\"checkbox-theme\" id=\"theme-switch\" title=\"light/dark
  mode toggle\">\n<input id=\"checkbox-theme\" type=\"checkbox\"/>\n<div class=\"slider
  round\"></div>\n</label>\n</div>\n<section class=\"title\">\n<h1 id=\"title\">\n
  \           Post_Model.Py \n            \n        </h1>\n</section>\n<main><p>None</p>\n<p>!!
  class </p><h2 class=\"admonition-title\" id=\"Post\" style=\"margin:0;padding:.5rem
  1rem;\">Post <em class=\"small\">class</em></h2>\n<div class=\"admonition source
  is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"Post <em class=\"small\">source</em>\"</p>\n</div>\n<pre
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
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
  <span class=\"nc\">Post</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
  class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
  \           <span class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"n\">Any</span>
  <span class=\"o\">=</span> <span class=\"kc\">None</span>\n            <span class=\"n\">path</span><span
  class=\"p\">:</span> <span class=\"n\">Path</span>\n            <span class=\"n\">slug</span><span
  class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
  class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
  class=\"kc\">None</span>\n            <span class=\"n\">href</span><span class=\"p\">:</span>
  <span class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
  class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n
  \           <span class=\"n\">published</span><span class=\"p\">:</span> <span class=\"nb\">bool</span>
  <span class=\"o\">=</span> <span class=\"kc\">False</span>\n            <span class=\"n\">description</span><span
  class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
  class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
  class=\"kc\">None</span>\n            <span class=\"n\">content</span><span class=\"p\">:</span>
  <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n
  \           <span class=\"c1\"># date: Union[datetime.date, str]=None</span>\n            <span
  class=\"n\">date</span><span class=\"p\">:</span> <span class=\"n\">Optional</span><span
  class=\"p\">[</span><span class=\"n\">Union</span><span class=\"p\">[</span><span
  class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">date</span><span
  class=\"p\">,</span> <span class=\"nb\">str</span><span class=\"p\">]]</span> <span
  class=\"o\">=</span> <span class=\"kc\">None</span>\n            <span class=\"c1\">#
  pydantic.Field(</span>\n            <span class=\"c1\"># default_factory=lambda:
  datetime.date.min</span>\n            <span class=\"c1\"># )</span>\n            <span
  class=\"n\">date_time</span><span class=\"p\">:</span> <span class=\"n\">Optional</span><span
  class=\"p\">[</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
  class=\"n\">datetime</span><span class=\"p\">]</span> <span class=\"o\">=</span>
  <span class=\"kc\">None</span>\n            <span class=\"n\">today</span><span
  class=\"p\">:</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
  class=\"n\">date</span> <span class=\"o\">=</span> <span class=\"n\">pydantic</span><span
  class=\"o\">.</span><span class=\"n\">Field</span><span class=\"p\">(</span><span
  class=\"n\">default_factory</span><span class=\"o\">=</span><span class=\"n\">datetime</span><span
  class=\"o\">.</span><span class=\"n\">date</span><span class=\"o\">.</span><span
  class=\"n\">today</span><span class=\"p\">)</span>\n            <span class=\"n\">now</span><span
  class=\"p\">:</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
  class=\"n\">datetime</span> <span class=\"o\">=</span> <span class=\"n\">pydantic</span><span
  class=\"o\">.</span><span class=\"n\">Field</span><span class=\"p\">(</span><span
  class=\"n\">default_factory</span><span class=\"o\">=</span><span class=\"n\">datetime</span><span
  class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
  class=\"n\">utcnow</span><span class=\"p\">)</span>\n            <span class=\"n\">load_time</span><span
  class=\"p\">:</span> <span class=\"nb\">float</span> <span class=\"o\">=</span>
  <span class=\"mi\">0</span>\n            <span class=\"n\">profile</span><span class=\"p\">:</span>
  <span class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
  class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n
  \           <span class=\"n\">title</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
  <span class=\"o\">=</span> <span class=\"kc\">None</span>\n            <span class=\"n\">model_config</span>
  <span class=\"o\">=</span> <span class=\"n\">ConfigDict</span><span class=\"p\">(</span>\n
  \               <span class=\"n\">validate_assignment</span><span class=\"o\">=</span><span
  class=\"kc\">True</span><span class=\"p\">,</span>\n                <span class=\"n\">arbitrary_types_allowed</span><span
  class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n                <span
  class=\"n\">extra</span><span class=\"o\">=</span><span class=\"s2\">\"allow\"</span><span
  class=\"p\">,</span>\n            <span class=\"p\">)</span>\n\n            <span
  class=\"k\">def</span> <span class=\"nf\">__repr_args__</span><span class=\"p\">(</span><span
  class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">\"Post\"</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"s2\">\"ReprArgs\"</span><span
  class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"p\">[</span>\n
  \                   <span class=\"p\">(</span><span class=\"n\">key</span><span
  class=\"p\">,</span> <span class=\"n\">value</span><span class=\"p\">)</span>\n
  \                   <span class=\"k\">for</span> <span class=\"n\">key</span><span
  class=\"p\">,</span> <span class=\"n\">value</span> <span class=\"ow\">in</span>
  <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"vm\">__dict__</span><span
  class=\"o\">.</span><span class=\"n\">items</span><span class=\"p\">()</span>\n
  \                   <span class=\"k\">if</span> <span class=\"n\">key</span> <span
  class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">post_model</span><span class=\"o\">.</span><span
  class=\"n\">repr_include</span>\n                <span class=\"p\">]</span>\n\n
  \           <span class=\"nd\">@property</span>\n            <span class=\"k\">def</span>
  <span class=\"nf\">metadata</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Post\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"n\">Dict</span><span class=\"p\">:</span>\n
  \               <span class=\"s2\">\"for backwards compatability\"</span>\n                <span
  class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"vm\">__dict__</span>\n\n            <span class=\"k\">def</span> <span class=\"nf\">to_dict</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">:</span> <span
  class=\"s2\">\"Post\"</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"n\">Dict</span><span class=\"p\">:</span>\n                <span class=\"s2\">\"for
  backwards compatability\"</span>\n                <span class=\"k\">return</span>
  <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"vm\">__dict__</span>\n\n
  \           <span class=\"k\">def</span> <span class=\"fm\">__getitem__</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">:</span> <span
  class=\"s2\">\"Post\"</span><span class=\"p\">,</span> <span class=\"n\">item</span><span
  class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">)</span> <span
  class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span class=\"p\">:</span>\n
  \               <span class=\"s2\">\"for backwards compatability\"</span>\n                <span
  class=\"k\">return</span> <span class=\"nb\">getattr</span><span class=\"p\">(</span><span
  class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"n\">item</span><span
  class=\"p\">)</span>\n\n            <span class=\"k\">def</span> <span class=\"fm\">__setitem__</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">:</span> <span
  class=\"s2\">\"Post\"</span><span class=\"p\">,</span> <span class=\"n\">key</span><span
  class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">,</span> <span
  class=\"n\">item</span><span class=\"p\">:</span> <span class=\"n\">Any</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n                <span class=\"s2\">\"for backwards compatability\"</span>\n
  \               <span class=\"nb\">setattr</span><span class=\"p\">(</span><span
  class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"n\">key</span><span
  class=\"p\">,</span> <span class=\"n\">item</span><span class=\"p\">)</span>\n\n
  \           <span class=\"k\">def</span> <span class=\"nf\">get</span><span class=\"p\">(</span><span
  class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">\"Post\"</span><span
  class=\"p\">,</span> <span class=\"n\">item</span><span class=\"p\">:</span> <span
  class=\"nb\">str</span><span class=\"p\">,</span> <span class=\"n\">default</span><span
  class=\"p\">:</span> <span class=\"n\">Any</span><span class=\"p\">)</span> <span
  class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span class=\"p\">:</span>\n
  \               <span class=\"s2\">\"for backwards compatability\"</span>\n                <span
  class=\"k\">return</span> <span class=\"nb\">getattr</span><span class=\"p\">(</span><span
  class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"n\">item</span><span
  class=\"p\">,</span> <span class=\"n\">default</span><span class=\"p\">)</span>\n\n
  \           <span class=\"k\">def</span> <span class=\"nf\">keys</span><span class=\"p\">(</span><span
  class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">\"Post\"</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">List</span><span
  class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">]:</span>\n                <span
  class=\"s2\">\"for backwards compatability\"</span>\n                <span class=\"k\">return</span>
  <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"vm\">__dict__</span><span
  class=\"o\">.</span><span class=\"n\">keys</span><span class=\"p\">()</span>\n\n
  \           <span class=\"c1\"># def json(</span>\n            <span class=\"c1\">#
  \    self: \"Post\",</span>\n            <span class=\"c1\">#     include: Iterable
  = None,</span>\n            <span class=\"c1\">#     all: bool = False,</span>\n
  \           <span class=\"c1\">#     **kwargs,</span>\n            <span class=\"c1\">#
  ) -&gt; str:</span>\n            <span class=\"c1\">#     \"\"\"</span>\n            <span
  class=\"c1\">#     override function to give a default include value that will include</span>\n
  \           <span class=\"c1\">#     user configured includes.</span>\n            <span
  class=\"c1\">#     \"\"\"</span>\n            <span class=\"c1\">#     if all:</span>\n
  \           <span class=\"c1\">#         return pydantic.create_model(\"Post\",
  **self)(**self).json(</span>\n            <span class=\"c1\">#             **kwargs,</span>\n
  \           <span class=\"c1\">#         )</span>\n            <span class=\"c1\">#
  \    if include:</span>\n            <span class=\"c1\">#         return pydantic.create_model(\"Post\",
  **self)(**self).json(</span>\n            <span class=\"c1\">#             include=include,</span>\n
  \           <span class=\"c1\">#             **kwargs,</span>\n            <span
  class=\"c1\">#         )</span>\n            <span class=\"c1\">#     return pydantic.create_model(\"Post\",
  **self)(**self).json(</span>\n            <span class=\"c1\">#         include={i:
  True for i in self.markata.config.post_model.include},</span>\n            <span
  class=\"c1\">#         **kwargs,</span>\n            <span class=\"c1\">#     )</span>\n\n
  \           <span class=\"k\">def</span> <span class=\"nf\">yaml</span><span class=\"p\">(</span><span
  class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">\"Post\"</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span
  class=\"p\">:</span>\n<span class=\"w\">                </span><span class=\"sd\">\"\"\"</span>\n<span
  class=\"sd\">                dump model to yaml</span>\n<span class=\"sd\">                \"\"\"</span>\n
  \               <span class=\"kn\">import</span> <span class=\"nn\">yaml</span>\n\n
  \               <span class=\"k\">return</span> <span class=\"n\">yaml</span><span
  class=\"o\">.</span><span class=\"n\">dump</span><span class=\"p\">(</span>\n                    <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">dict</span><span
  class=\"p\">(</span>\n                        <span class=\"n\">include</span><span
  class=\"o\">=</span><span class=\"p\">{</span><span class=\"n\">i</span><span class=\"p\">:</span>
  <span class=\"kc\">True</span> <span class=\"k\">for</span> <span class=\"n\">i</span>
  <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">post_model</span><span class=\"o\">.</span><span
  class=\"n\">include</span><span class=\"p\">}</span>\n                    <span
  class=\"p\">),</span>\n                    <span class=\"n\">Dumper</span><span
  class=\"o\">=</span><span class=\"n\">yaml</span><span class=\"o\">.</span><span
  class=\"n\">CDumper</span><span class=\"p\">,</span>\n                <span class=\"p\">)</span>\n\n
  \           <span class=\"k\">def</span> <span class=\"nf\">markdown</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">:</span> <span
  class=\"s2\">\"Post\"</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"nb\">str</span><span class=\"p\">:</span>\n<span class=\"w\">                </span><span
  class=\"sd\">\"\"\"</span>\n<span class=\"sd\">                dump model to markdown</span>\n<span
  class=\"sd\">                \"\"\"</span>\n\n                <span class=\"kn\">import</span>
  <span class=\"nn\">yaml</span>\n\n                <span class=\"n\">frontmatter</span>
  <span class=\"o\">=</span> <span class=\"n\">yaml</span><span class=\"o\">.</span><span
  class=\"n\">dump</span><span class=\"p\">(</span>\n                    <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">dict</span><span class=\"p\">(</span>\n                        <span
  class=\"n\">include</span><span class=\"o\">=</span><span class=\"p\">{</span>\n
  \                           <span class=\"n\">i</span><span class=\"p\">:</span>
  <span class=\"kc\">True</span>\n                            <span class=\"k\">for</span>
  <span class=\"n\">i</span> <span class=\"ow\">in</span> <span class=\"p\">[</span>\n
  \                               <span class=\"n\">_i</span>\n                                <span
  class=\"k\">for</span> <span class=\"n\">_i</span> <span class=\"ow\">in</span>
  <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
  class=\"n\">post_model</span><span class=\"o\">.</span><span class=\"n\">include</span>\n
  \                               <span class=\"k\">if</span> <span class=\"n\">_i</span>
  <span class=\"o\">!=</span> <span class=\"s2\">\"content\"</span>\n                            <span
  class=\"p\">]</span>\n                        <span class=\"p\">}</span>\n                    <span
  class=\"p\">),</span>\n                    <span class=\"n\">Dumper</span><span
  class=\"o\">=</span><span class=\"n\">yaml</span><span class=\"o\">.</span><span
  class=\"n\">CDumper</span><span class=\"p\">,</span>\n                <span class=\"p\">)</span>\n
  \               <span class=\"n\">post</span> <span class=\"o\">=</span> <span class=\"s2\">\"---</span><span
  class=\"se\">\\n</span><span class=\"s2\">\"</span>\n                <span class=\"n\">post</span>
  <span class=\"o\">+=</span> <span class=\"n\">frontmatter</span>\n                <span
  class=\"n\">post</span> <span class=\"o\">+=</span> <span class=\"s2\">\"---</span><span
  class=\"se\">\\n\\n</span><span class=\"s2\">\"</span>\n\n                <span
  class=\"k\">if</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">content</span><span class=\"p\">:</span>\n                    <span
  class=\"n\">post</span> <span class=\"o\">+=</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">content</span>\n                <span class=\"k\">return</span>
  <span class=\"n\">post</span>\n\n            <span class=\"nd\">@classmethod</span>\n
  \           <span class=\"k\">def</span> <span class=\"nf\">parse_file</span><span
  class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
  class=\"n\">markata</span><span class=\"p\">,</span> <span class=\"n\">path</span><span
  class=\"p\">:</span> <span class=\"n\">Union</span><span class=\"p\">[</span><span
  class=\"n\">Path</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
  class=\"p\">],</span> <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"s2\">\"Post\"</span><span
  class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
  class=\"p\">(</span><span class=\"n\">path</span><span class=\"p\">,</span> <span
  class=\"n\">Path</span><span class=\"p\">):</span>\n                    <span class=\"k\">if</span>
  <span class=\"n\">path</span><span class=\"o\">.</span><span class=\"n\">suffix</span>
  <span class=\"ow\">in</span> <span class=\"p\">[</span><span class=\"s2\">\".md\"</span><span
  class=\"p\">,</span> <span class=\"s2\">\".markdown\"</span><span class=\"p\">]:</span>\n
  \                       <span class=\"k\">return</span> <span class=\"bp\">cls</span><span
  class=\"o\">.</span><span class=\"n\">parse_markdown</span><span class=\"p\">(</span><span
  class=\"n\">markata</span><span class=\"o\">=</span><span class=\"n\">markata</span><span
  class=\"p\">,</span> <span class=\"n\">path</span><span class=\"o\">=</span><span
  class=\"n\">path</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
  class=\"n\">kwargs</span><span class=\"p\">)</span>\n                <span class=\"k\">elif</span>
  <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">path</span><span
  class=\"p\">,</span> <span class=\"nb\">str</span><span class=\"p\">):</span>\n
  \                   <span class=\"k\">if</span> <span class=\"n\">path</span><span
  class=\"o\">.</span><span class=\"n\">endswith</span><span class=\"p\">(</span><span
  class=\"s2\">\".md\"</span><span class=\"p\">)</span> <span class=\"ow\">or</span>
  <span class=\"n\">path</span><span class=\"o\">.</span><span class=\"n\">endswith</span><span
  class=\"p\">(</span><span class=\"s2\">\".markdown\"</span><span class=\"p\">):</span>\n
  \                       <span class=\"k\">return</span> <span class=\"bp\">cls</span><span
  class=\"o\">.</span><span class=\"n\">parse_markdown</span><span class=\"p\">(</span><span
  class=\"n\">markata</span><span class=\"o\">=</span><span class=\"n\">markata</span><span
  class=\"p\">,</span> <span class=\"n\">path</span><span class=\"o\">=</span><span
  class=\"n\">path</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
  class=\"n\">kwargs</span><span class=\"p\">)</span>\n                <span class=\"k\">return</span>
  <span class=\"nb\">super</span><span class=\"p\">(</span><span class=\"n\">Post</span><span
  class=\"p\">,</span> <span class=\"bp\">cls</span><span class=\"p\">)</span><span
  class=\"o\">.</span><span class=\"n\">parse_file</span><span class=\"p\">(</span><span
  class=\"n\">path</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
  class=\"n\">kwargs</span><span class=\"p\">)</span>\n\n            <span class=\"nd\">@classmethod</span>\n
  \           <span class=\"k\">def</span> <span class=\"nf\">parse_markdown</span><span
  class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
  class=\"n\">markata</span><span class=\"p\">,</span> <span class=\"n\">path</span><span
  class=\"p\">:</span> <span class=\"n\">Union</span><span class=\"p\">[</span><span
  class=\"n\">Path</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
  class=\"p\">],</span> <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"s2\">\"Post\"</span><span
  class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
  class=\"p\">(</span><span class=\"n\">path</span><span class=\"p\">,</span> <span
  class=\"nb\">str</span><span class=\"p\">):</span>\n                    <span class=\"n\">path</span>
  <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
  class=\"n\">path</span><span class=\"p\">)</span>\n                <span class=\"n\">text</span>
  <span class=\"o\">=</span> <span class=\"n\">path</span><span class=\"o\">.</span><span
  class=\"n\">read_text</span><span class=\"p\">()</span>\n                <span class=\"k\">try</span><span
  class=\"p\">:</span>\n                    <span class=\"n\">_</span><span class=\"p\">,</span>
  <span class=\"n\">fm</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
  class=\"n\">content</span> <span class=\"o\">=</span> <span class=\"n\">text</span><span
  class=\"o\">.</span><span class=\"n\">split</span><span class=\"p\">(</span><span
  class=\"s2\">\"---</span><span class=\"se\">\\n</span><span class=\"s2\">\"</span><span
  class=\"p\">)</span>\n                    <span class=\"n\">content</span> <span
  class=\"o\">=</span> <span class=\"s2\">\"---</span><span class=\"se\">\\n</span><span
  class=\"s2\">\"</span><span class=\"o\">.</span><span class=\"n\">join</span><span
  class=\"p\">(</span><span class=\"n\">content</span><span class=\"p\">)</span>\n
  \                   <span class=\"k\">try</span><span class=\"p\">:</span>\n                        <span
  class=\"n\">fm</span> <span class=\"o\">=</span> <span class=\"n\">yaml</span><span
  class=\"o\">.</span><span class=\"n\">load</span><span class=\"p\">(</span><span
  class=\"n\">fm</span><span class=\"p\">,</span> <span class=\"n\">Loader</span><span
  class=\"o\">=</span><span class=\"n\">yaml</span><span class=\"o\">.</span><span
  class=\"n\">CBaseLoader</span><span class=\"p\">)</span>\n                    <span
  class=\"k\">except</span> <span class=\"n\">yaml</span><span class=\"o\">.</span><span
  class=\"n\">YAMLError</span><span class=\"p\">:</span>\n                        <span
  class=\"n\">fm</span> <span class=\"o\">=</span> <span class=\"p\">{}</span>\n                <span
  class=\"k\">except</span> <span class=\"ne\">ValueError</span><span class=\"p\">:</span>\n
  \                   <span class=\"n\">fm</span> <span class=\"o\">=</span> <span
  class=\"p\">{}</span>\n                    <span class=\"n\">content</span> <span
  class=\"o\">=</span> <span class=\"n\">text</span>\n                <span class=\"k\">if</span>
  <span class=\"n\">fm</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span>
  <span class=\"ow\">or</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
  class=\"n\">fm</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
  class=\"p\">):</span>\n                    <span class=\"n\">fm</span> <span class=\"o\">=</span>
  <span class=\"p\">{}</span>\n\n                <span class=\"n\">post_args</span>
  <span class=\"o\">=</span> <span class=\"p\">{</span>\n                    <span
  class=\"s2\">\"markata\"</span><span class=\"p\">:</span> <span class=\"n\">markata</span><span
  class=\"p\">,</span>\n                    <span class=\"s2\">\"path\"</span><span
  class=\"p\">:</span> <span class=\"n\">path</span><span class=\"p\">,</span>\n                    <span
  class=\"s2\">\"content\"</span><span class=\"p\">:</span> <span class=\"n\">content</span><span
  class=\"p\">,</span>\n                    <span class=\"o\">**</span><span class=\"n\">fm</span><span
  class=\"p\">,</span>\n                <span class=\"p\">}</span>\n\n                <span
  class=\"k\">return</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">Post</span><span class=\"p\">(</span><span class=\"o\">**</span><span
  class=\"n\">post_args</span><span class=\"p\">)</span>\n\n            <span class=\"k\">def</span>
  <span class=\"nf\">dumps</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">):</span>\n<span class=\"w\">                </span><span class=\"sd\">\"\"\"</span>\n<span
  class=\"sd\">                dumps raw article back out</span>\n<span class=\"sd\">
  \               \"\"\"</span>\n                <span class=\"k\">return</span> <span
  class=\"sa\">f</span><span class=\"s2\">\"---</span><span class=\"se\">\\n</span><span
  class=\"si\">{</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">yaml</span><span class=\"p\">()</span><span class=\"si\">}</span><span
  class=\"se\">\\n\\n</span><span class=\"s2\">---</span><span class=\"se\">\\n\\n</span><span
  class=\"si\">{</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">content</span><span class=\"si\">}</span><span class=\"s2\">\"</span>\n\n
  \           <span class=\"nd\">@pydantic</span><span class=\"o\">.</span><span class=\"n\">validator</span><span
  class=\"p\">(</span><span class=\"s2\">\"slug\"</span><span class=\"p\">,</span>
  <span class=\"n\">pre</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
  class=\"p\">,</span> <span class=\"n\">always</span><span class=\"o\">=</span><span
  class=\"kc\">True</span><span class=\"p\">)</span>\n            <span class=\"k\">def</span>
  <span class=\"nf\">default_slug</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
  class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
  class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
  class=\"p\">):</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span>
  <span class=\"ow\">or</span> <span class=\"n\">slugify</span><span class=\"p\">(</span><span
  class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">values</span><span
  class=\"p\">[</span><span class=\"s2\">\"path\"</span><span class=\"p\">]</span><span
  class=\"o\">.</span><span class=\"n\">stem</span><span class=\"p\">))</span>\n\n
  \           <span class=\"nd\">@pydantic</span><span class=\"o\">.</span><span class=\"n\">validator</span><span
  class=\"p\">(</span><span class=\"s2\">\"slug\"</span><span class=\"p\">,</span>
  <span class=\"n\">pre</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
  class=\"p\">,</span> <span class=\"n\">always</span><span class=\"o\">=</span><span
  class=\"kc\">True</span><span class=\"p\">)</span>\n            <span class=\"k\">def</span>
  <span class=\"nf\">index_slug_is_empty</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
  class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
  class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
  class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span>
  <span class=\"o\">==</span> <span class=\"s2\">\"index\"</span><span class=\"p\">:</span>\n
  \                   <span class=\"k\">return</span> <span class=\"s2\">\"\"</span>\n
  \               <span class=\"k\">return</span> <span class=\"n\">v</span>\n\n            <span
  class=\"nd\">@pydantic</span><span class=\"o\">.</span><span class=\"n\">validator</span><span
  class=\"p\">(</span><span class=\"s2\">\"href\"</span><span class=\"p\">,</span>
  <span class=\"n\">pre</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
  class=\"p\">,</span> <span class=\"n\">always</span><span class=\"o\">=</span><span
  class=\"kc\">True</span><span class=\"p\">)</span>\n            <span class=\"k\">def</span>
  <span class=\"nf\">default_href</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
  class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
  class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
  class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span><span
  class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
  class=\"n\">v</span>\n                <span class=\"k\">return</span> <span class=\"sa\">f</span><span
  class=\"s2\">\"/</span><span class=\"si\">{</span><span class=\"n\">values</span><span
  class=\"p\">[</span><span class=\"s1\">'slug'</span><span class=\"p\">]</span><span
  class=\"o\">.</span><span class=\"n\">strip</span><span class=\"p\">(</span><span
  class=\"s1\">'/'</span><span class=\"p\">)</span><span class=\"si\">}</span><span
  class=\"s2\">/\"</span><span class=\"o\">.</span><span class=\"n\">replace</span><span
  class=\"p\">(</span><span class=\"s2\">\"//\"</span><span class=\"p\">,</span> <span
  class=\"s2\">\"/\"</span><span class=\"p\">)</span>\n\n            <span class=\"nd\">@pydantic</span><span
  class=\"o\">.</span><span class=\"n\">validator</span><span class=\"p\">(</span><span
  class=\"s2\">\"title\"</span><span class=\"p\">,</span> <span class=\"n\">pre</span><span
  class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
  class=\"n\">always</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
  class=\"p\">)</span>\n            <span class=\"k\">def</span> <span class=\"nf\">title_title</span><span
  class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
  class=\"n\">v</span><span class=\"p\">,</span> <span class=\"o\">*</span><span class=\"p\">,</span>
  <span class=\"n\">values</span><span class=\"p\">):</span>\n                <span
  class=\"n\">title</span> <span class=\"o\">=</span> <span class=\"n\">v</span> <span
  class=\"ow\">or</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
  class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">\"path\"</span><span
  class=\"p\">])</span><span class=\"o\">.</span><span class=\"n\">stem</span><span
  class=\"o\">.</span><span class=\"n\">replace</span><span class=\"p\">(</span><span
  class=\"s2\">\"-\"</span><span class=\"p\">,</span> <span class=\"s2\">\" \"</span><span
  class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"n\">title</span><span
  class=\"o\">.</span><span class=\"n\">title</span><span class=\"p\">()</span>\n\n
  \           <span class=\"nd\">@pydantic</span><span class=\"o\">.</span><span class=\"n\">validator</span><span
  class=\"p\">(</span><span class=\"s2\">\"date_time\"</span><span class=\"p\">,</span>
  <span class=\"n\">pre</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
  class=\"p\">,</span> <span class=\"n\">always</span><span class=\"o\">=</span><span
  class=\"kc\">True</span><span class=\"p\">)</span>\n            <span class=\"k\">def</span>
  <span class=\"nf\">dateparser_datetime</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
  class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
  class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
  class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
  class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
  class=\"p\">):</span>\n                    <span class=\"n\">d</span> <span class=\"o\">=</span>
  <span class=\"n\">dateparser</span><span class=\"o\">.</span><span class=\"n\">parse</span><span
  class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">)</span>\n                    <span
  class=\"k\">if</span> <span class=\"n\">d</span> <span class=\"ow\">is</span> <span
  class=\"kc\">None</span><span class=\"p\">:</span>\n                        <span
  class=\"k\">raise</span> <span class=\"ne\">ValueError</span><span class=\"p\">(</span><span
  class=\"sa\">f</span><span class=\"s1\">'\"</span><span class=\"si\">{</span><span
  class=\"n\">v</span><span class=\"si\">}</span><span class=\"s1\">\" is not a valid
  date'</span><span class=\"p\">)</span>\n                <span class=\"k\">return</span>
  <span class=\"n\">v</span>\n\n            <span class=\"nd\">@pydantic</span><span
  class=\"o\">.</span><span class=\"n\">validator</span><span class=\"p\">(</span><span
  class=\"s2\">\"date_time\"</span><span class=\"p\">,</span> <span class=\"n\">pre</span><span
  class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
  class=\"n\">always</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
  class=\"p\">)</span>\n            <span class=\"k\">def</span> <span class=\"nf\">date_is_datetime</span><span
  class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
  class=\"n\">v</span><span class=\"p\">,</span> <span class=\"o\">*</span><span class=\"p\">,</span>
  <span class=\"n\">values</span><span class=\"p\">):</span>\n                <span
  class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"ow\">is</span> <span
  class=\"kc\">None</span> <span class=\"ow\">and</span> <span class=\"s2\">\"date\"</span>
  <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">values</span><span
  class=\"p\">:</span>\n                    <span class=\"n\">values</span><span class=\"p\">[</span><span
  class=\"s2\">\"markata\"</span><span class=\"p\">]</span><span class=\"o\">.</span><span
  class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
  class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">\"</span><span
  class=\"si\">{</span><span class=\"n\">values</span><span class=\"p\">[</span><span
  class=\"s1\">'path'</span><span class=\"p\">]</span><span class=\"si\">}</span><span
  class=\"s2\"> has no date\"</span><span class=\"p\">)</span>\n                    <span
  class=\"k\">return</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
  class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">now</span><span
  class=\"p\">()</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span>
  <span class=\"ow\">is</span> <span class=\"kc\">None</span> <span class=\"ow\">and</span>
  <span class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">\"date\"</span><span
  class=\"p\">]</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n                    <span class=\"n\">values</span><span class=\"p\">[</span><span
  class=\"s2\">\"markata\"</span><span class=\"p\">]</span><span class=\"o\">.</span><span
  class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
  class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">\"</span><span
  class=\"si\">{</span><span class=\"n\">values</span><span class=\"p\">[</span><span
  class=\"s1\">'path'</span><span class=\"p\">]</span><span class=\"si\">}</span><span
  class=\"s2\"> has no date\"</span><span class=\"p\">)</span>\n                    <span
  class=\"k\">return</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
  class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">now</span><span
  class=\"p\">()</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
  class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span class=\"n\">datetime</span><span
  class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"p\">):</span>\n
  \                   <span class=\"k\">return</span> <span class=\"n\">v</span>\n
  \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
  class=\"p\">(</span><span class=\"n\">values</span><span class=\"p\">[</span><span
  class=\"s2\">\"date\"</span><span class=\"p\">],</span> <span class=\"n\">datetime</span><span
  class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"p\">):</span>\n
  \                   <span class=\"k\">return</span> <span class=\"n\">values</span><span
  class=\"p\">[</span><span class=\"s2\">\"date\"</span><span class=\"p\">]</span>\n
  \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
  class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span class=\"n\">datetime</span><span
  class=\"o\">.</span><span class=\"n\">date</span><span class=\"p\">):</span>\n                    <span
  class=\"k\">return</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
  class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">combine</span><span
  class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span class=\"n\">datetime</span><span
  class=\"o\">.</span><span class=\"n\">time</span><span class=\"o\">.</span><span
  class=\"n\">min</span><span class=\"p\">)</span>\n                <span class=\"k\">if</span>
  <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">values</span><span
  class=\"p\">[</span><span class=\"s2\">\"date\"</span><span class=\"p\">],</span>
  <span class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">date</span><span
  class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
  class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
  class=\"o\">.</span><span class=\"n\">combine</span><span class=\"p\">(</span><span
  class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">\"date\"</span><span
  class=\"p\">],</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
  class=\"n\">time</span><span class=\"o\">.</span><span class=\"n\">min</span><span
  class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span>\n\n
  \           <span class=\"nd\">@pydantic</span><span class=\"o\">.</span><span class=\"n\">validator</span><span
  class=\"p\">(</span><span class=\"s2\">\"date_time\"</span><span class=\"p\">,</span>
  <span class=\"n\">pre</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
  class=\"p\">,</span> <span class=\"n\">always</span><span class=\"o\">=</span><span
  class=\"kc\">True</span><span class=\"p\">)</span>\n            <span class=\"k\">def</span>
  <span class=\"nf\">mindate_time</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
  class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
  class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
  class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span>
  <span class=\"ow\">is</span> <span class=\"kc\">None</span> <span class=\"ow\">and</span>
  <span class=\"s2\">\"date\"</span> <span class=\"ow\">not</span> <span class=\"ow\">in</span>
  <span class=\"n\">values</span><span class=\"p\">:</span>\n                    <span
  class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">\"markata\"</span><span
  class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">console</span><span
  class=\"o\">.</span><span class=\"n\">log</span><span class=\"p\">(</span><span
  class=\"sa\">f</span><span class=\"s2\">\"</span><span class=\"si\">{</span><span
  class=\"n\">values</span><span class=\"p\">[</span><span class=\"s1\">'path'</span><span
  class=\"p\">]</span><span class=\"si\">}</span><span class=\"s2\"> has no date\"</span><span
  class=\"p\">)</span>\n                    <span class=\"k\">return</span> <span
  class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
  class=\"o\">.</span><span class=\"n\">min</span>\n                <span class=\"k\">if</span>
  <span class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">\"date\"</span><span
  class=\"p\">]</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n                    <span class=\"n\">values</span><span class=\"p\">[</span><span
  class=\"s2\">\"markata\"</span><span class=\"p\">]</span><span class=\"o\">.</span><span
  class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
  class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">\"</span><span
  class=\"si\">{</span><span class=\"n\">values</span><span class=\"p\">[</span><span
  class=\"s1\">'path'</span><span class=\"p\">]</span><span class=\"si\">}</span><span
  class=\"s2\"> has no date\"</span><span class=\"p\">)</span>\n                    <span
  class=\"k\">return</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
  class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">min</span>\n
  \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
  class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span class=\"n\">datetime</span><span
  class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"p\">):</span>\n
  \                   <span class=\"k\">return</span> <span class=\"n\">v</span>\n
  \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
  class=\"p\">(</span><span class=\"n\">values</span><span class=\"p\">[</span><span
  class=\"s2\">\"date\"</span><span class=\"p\">],</span> <span class=\"n\">datetime</span><span
  class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"p\">):</span>\n
  \                   <span class=\"k\">return</span> <span class=\"n\">values</span><span
  class=\"p\">[</span><span class=\"s2\">\"date\"</span><span class=\"p\">]</span>\n
  \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
  class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span class=\"n\">datetime</span><span
  class=\"o\">.</span><span class=\"n\">date</span><span class=\"p\">):</span>\n                    <span
  class=\"k\">return</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
  class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">combine</span><span
  class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span class=\"n\">datetime</span><span
  class=\"o\">.</span><span class=\"n\">time</span><span class=\"o\">.</span><span
  class=\"n\">min</span><span class=\"p\">)</span>\n                <span class=\"k\">if</span>
  <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">values</span><span
  class=\"p\">[</span><span class=\"s2\">\"date\"</span><span class=\"p\">],</span>
  <span class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">date</span><span
  class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
  class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
  class=\"o\">.</span><span class=\"n\">combine</span><span class=\"p\">(</span><span
  class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">\"date\"</span><span
  class=\"p\">],</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
  class=\"n\">time</span><span class=\"o\">.</span><span class=\"n\">min</span><span
  class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span>\n\n
  \           <span class=\"nd\">@pydantic</span><span class=\"o\">.</span><span class=\"n\">validator</span><span
  class=\"p\">(</span><span class=\"s2\">\"date\"</span><span class=\"p\">,</span>
  <span class=\"n\">pre</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
  class=\"p\">,</span> <span class=\"n\">always</span><span class=\"o\">=</span><span
  class=\"kc\">True</span><span class=\"p\">)</span>\n            <span class=\"k\">def</span>
  <span class=\"nf\">dateparser_date</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
  class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
  class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
  class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span>
  <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \                   <span class=\"k\">return</span> <span class=\"n\">datetime</span><span
  class=\"o\">.</span><span class=\"n\">date</span><span class=\"o\">.</span><span
  class=\"n\">min</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
  class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
  class=\"p\">):</span>\n                    <span class=\"n\">d</span> <span class=\"o\">=</span>
  <span class=\"bp\">cls</span><span class=\"o\">.</span><span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">precache</span><span class=\"o\">.</span><span
  class=\"n\">get</span><span class=\"p\">(</span><span class=\"n\">v</span><span
  class=\"p\">)</span>\n                    <span class=\"k\">if</span> <span class=\"n\">d</span>
  <span class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n                        <span class=\"k\">return</span> <span
  class=\"n\">d</span>\n                    <span class=\"n\">d</span> <span class=\"o\">=</span>
  <span class=\"n\">dateparser</span><span class=\"o\">.</span><span class=\"n\">parse</span><span
  class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">)</span>\n                    <span
  class=\"k\">if</span> <span class=\"n\">d</span> <span class=\"ow\">is</span> <span
  class=\"kc\">None</span><span class=\"p\">:</span>\n                        <span
  class=\"k\">raise</span> <span class=\"ne\">ValueError</span><span class=\"p\">(</span><span
  class=\"sa\">f</span><span class=\"s1\">'\"</span><span class=\"si\">{</span><span
  class=\"n\">v</span><span class=\"si\">}</span><span class=\"s1\">\" is not a valid
  date'</span><span class=\"p\">)</span>\n                    <span class=\"n\">d</span>
  <span class=\"o\">=</span> <span class=\"n\">d</span><span class=\"o\">.</span><span
  class=\"n\">date</span><span class=\"p\">()</span>\n                    <span class=\"k\">with</span>
  <span class=\"bp\">cls</span><span class=\"o\">.</span><span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">cache</span> <span class=\"k\">as</span> <span
  class=\"n\">cache</span><span class=\"p\">:</span>\n                        <span
  class=\"n\">cache</span><span class=\"o\">.</span><span class=\"n\">add</span><span
  class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span class=\"n\">d</span><span
  class=\"p\">)</span>\n                    <span class=\"k\">return</span> <span
  class=\"n\">d</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n<p>!!
  class </p><h2 class=\"admonition-title\" id=\"PostModelConfig\" style=\"margin:0;padding:.5rem
  1rem;\">PostModelConfig <em class=\"small\">class</em></h2>\nConfiguration for the
  Post model\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
  class=\"admonition-title\">\"PostModelConfig <em class=\"small\">source</em>\"</p>\n</div>\n<pre
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
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
  <span class=\"nc\">PostModelConfig</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
  class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
  \           <span class=\"s2\">\"Configuration for the Post model\"</span>\n\n            <span
  class=\"k\">def</span> <span class=\"fm\">__init__</span><span class=\"p\">(</span><span
  class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
  class=\"n\">data</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span class=\"w\">                </span><span
  class=\"sd\">\"\"\"</span>\n\n<span class=\"sd\">                include: post attributes
  to include by default in Post</span>\n<span class=\"sd\">                model serialization.</span>\n<span
  class=\"sd\">                repr_include: post attributes to include by default
  in Post</span>\n<span class=\"sd\">                repr.  If `repr_include` is None,
  it will default to</span>\n<span class=\"sd\">                `include`, but it
  is likely that you want less in the repr</span>\n<span class=\"sd\">                than
  serialized output.</span>\n\n<span class=\"sd\">                example:</span>\n\n<span
  class=\"sd\">                ``` toml title='markata.toml'</span>\n<span class=\"sd\">
  \               [markata.post_model]</span>\n<span class=\"sd\">                include
  = ['date', 'description', 'published',</span>\n<span class=\"sd\">                    'slug',
  'title', 'content', 'html']</span>\n<span class=\"sd\">                repr_include
  = ['date', 'description', 'published', 'slug', 'title']</span>\n<span class=\"sd\">
  \               ```</span>\n<span class=\"sd\">                \"\"\"</span>\n                <span
  class=\"nb\">super</span><span class=\"p\">()</span><span class=\"o\">.</span><span
  class=\"fm\">__init__</span><span class=\"p\">(</span><span class=\"o\">**</span><span
  class=\"n\">data</span><span class=\"p\">)</span>\n\n            <span class=\"n\">include</span><span
  class=\"p\">:</span> <span class=\"n\">List</span><span class=\"p\">[</span><span
  class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
  class=\"p\">[</span>\n                <span class=\"s2\">\"date\"</span><span class=\"p\">,</span>\n
  \               <span class=\"s2\">\"description\"</span><span class=\"p\">,</span>\n
  \               <span class=\"s2\">\"published\"</span><span class=\"p\">,</span>\n
  \               <span class=\"s2\">\"slug\"</span><span class=\"p\">,</span>\n                <span
  class=\"s2\">\"title\"</span><span class=\"p\">,</span>\n                <span class=\"s2\">\"content\"</span><span
  class=\"p\">,</span>\n                <span class=\"s2\">\"html\"</span><span class=\"p\">,</span>\n
  \           <span class=\"p\">]</span>\n            <span class=\"n\">repr_include</span><span
  class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
  class=\"n\">List</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
  class=\"p\">]]</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n                <span
  class=\"s2\">\"date\"</span><span class=\"p\">,</span>\n                <span class=\"s2\">\"description\"</span><span
  class=\"p\">,</span>\n                <span class=\"s2\">\"published\"</span><span
  class=\"p\">,</span>\n                <span class=\"s2\">\"slug\"</span><span class=\"p\">,</span>\n
  \               <span class=\"s2\">\"title\"</span><span class=\"p\">,</span>\n
  \           <span class=\"p\">]</span>\n\n            <span class=\"nd\">@pydantic</span><span
  class=\"o\">.</span><span class=\"n\">validator</span><span class=\"p\">(</span><span
  class=\"s2\">\"repr_include\"</span><span class=\"p\">,</span> <span class=\"n\">pre</span><span
  class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
  class=\"n\">always</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
  class=\"p\">)</span>\n            <span class=\"k\">def</span> <span class=\"nf\">repr_include_validator</span><span
  class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
  class=\"n\">v</span><span class=\"p\">,</span> <span class=\"o\">*</span><span class=\"p\">,</span>
  <span class=\"n\">values</span><span class=\"p\">):</span>\n                <span
  class=\"k\">if</span> <span class=\"n\">v</span><span class=\"p\">:</span>\n                    <span
  class=\"k\">return</span> <span class=\"n\">v</span>\n                <span class=\"k\">return</span>
  <span class=\"n\">values</span><span class=\"o\">.</span><span class=\"n\">get</span><span
  class=\"p\">(</span><span class=\"s2\">\"include\"</span><span class=\"p\">,</span>
  <span class=\"kc\">None</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  class </p><h2 class=\"admonition-title\" id=\"Config\" style=\"margin:0;padding:.5rem
  1rem;\">Config <em class=\"small\">class</em></h2>\n<div class=\"admonition source
  is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"Config <em class=\"small\">source</em>\"</p>\n</div>\n<pre
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
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
  <span class=\"nc\">Config</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
  class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
  \           <span class=\"n\">post_model</span><span class=\"p\">:</span> <span
  class=\"n\">PostModelConfig</span> <span class=\"o\">=</span> <span class=\"n\">pydantic</span><span
  class=\"o\">.</span><span class=\"n\">Field</span><span class=\"p\">(</span><span
  class=\"n\">default_factory</span><span class=\"o\">=</span><span class=\"n\">PostModelConfig</span><span
  class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!! function </p><h2 class=\"admonition-title\"
  id=\"post_model\" style=\"margin:0;padding:.5rem 1rem;\">post_model <em class=\"small\">function</em></h2>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"post_model
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
  <span class=\"nf\">post_model</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">post_models</span><span
  class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
  class=\"n\">Post</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"config_model\" style=\"margin:0;padding:.5rem
  1rem;\">config_model <em class=\"small\">function</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"config_model
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
  <span class=\"nf\">config_model</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config_models</span><span
  class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
  class=\"n\">Config</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  class </p><h2 class=\"admonition-title\" id=\"PostFactory\" style=\"margin:0;padding:.5rem
  1rem;\">PostFactory <em class=\"small\">class</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"PostFactory
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
  <span class=\"nc\">PostFactory</span><span class=\"p\">(</span><span class=\"n\">ModelFactory</span><span
  class=\"p\">):</span>\n            <span class=\"n\">__model__</span> <span class=\"o\">=</span>
  <span class=\"n\">Post</span>\n</pre></div>\n\n</pre>\n<p>!! method </p><h2 class=\"admonition-title\"
  id=\"__repr_args__\" style=\"margin:0;padding:.5rem 1rem;\"><strong>repr_args</strong>
  <em class=\"small\">method</em></h2>\n<div class=\"admonition source is-collapsible
  collapsible-open\">\n<p class=\"admonition-title\">\"<strong>repr_args</strong>
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
  <span class=\"nf\">__repr_args__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Post\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"s2\">\"ReprArgs\"</span><span class=\"p\">:</span>\n
  \               <span class=\"k\">return</span> <span class=\"p\">[</span>\n                    <span
  class=\"p\">(</span><span class=\"n\">key</span><span class=\"p\">,</span> <span
  class=\"n\">value</span><span class=\"p\">)</span>\n                    <span class=\"k\">for</span>
  <span class=\"n\">key</span><span class=\"p\">,</span> <span class=\"n\">value</span>
  <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"vm\">__dict__</span><span class=\"o\">.</span><span class=\"n\">items</span><span
  class=\"p\">()</span>\n                    <span class=\"k\">if</span> <span class=\"n\">key</span>
  <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">post_model</span><span class=\"o\">.</span><span
  class=\"n\">repr_include</span>\n                <span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"metadata\" style=\"margin:0;padding:.5rem
  1rem;\">metadata <em class=\"small\">method</em></h2>\nfor backwards compatability\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"metadata
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
  <span class=\"nf\">metadata</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Post\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"n\">Dict</span><span class=\"p\">:</span>\n
  \               <span class=\"s2\">\"for backwards compatability\"</span>\n                <span
  class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"vm\">__dict__</span>\n</pre></div>\n\n</pre>\n<p>!! method </p><h2 class=\"admonition-title\"
  id=\"to_dict\" style=\"margin:0;padding:.5rem 1rem;\">to_dict <em class=\"small\">method</em></h2>\nfor
  backwards compatability\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
  class=\"admonition-title\">\"to_dict <em class=\"small\">source</em>\"</p>\n</div>\n<pre
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
  <span class=\"nf\">to_dict</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Post\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"n\">Dict</span><span class=\"p\">:</span>\n
  \               <span class=\"s2\">\"for backwards compatability\"</span>\n                <span
  class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"vm\">__dict__</span>\n</pre></div>\n\n</pre>\n<p>!! method </p><h2 class=\"admonition-title\"
  id=\"__getitem__\" style=\"margin:0;padding:.5rem 1rem;\"><strong>getitem</strong>
  <em class=\"small\">method</em></h2>\nfor backwards compatability\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"<strong>getitem</strong>
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
  <span class=\"fm\">__getitem__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Post\"</span><span class=\"p\">,</span>
  <span class=\"n\">item</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span
  class=\"p\">:</span>\n                <span class=\"s2\">\"for backwards compatability\"</span>\n
  \               <span class=\"k\">return</span> <span class=\"nb\">getattr</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
  class=\"n\">item</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"__setitem__\" style=\"margin:0;padding:.5rem
  1rem;\"><strong>setitem</strong> <em class=\"small\">method</em></h2>\nfor backwards
  compatability\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
  class=\"admonition-title\">\"<strong>setitem</strong> <em class=\"small\">source</em>\"</p>\n</div>\n<pre
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
  <span class=\"fm\">__setitem__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Post\"</span><span class=\"p\">,</span>
  <span class=\"n\">key</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
  class=\"p\">,</span> <span class=\"n\">item</span><span class=\"p\">:</span> <span
  class=\"n\">Any</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
  class=\"s2\">\"for backwards compatability\"</span>\n                <span class=\"nb\">setattr</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
  class=\"n\">key</span><span class=\"p\">,</span> <span class=\"n\">item</span><span
  class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!! method </p><h2 class=\"admonition-title\"
  id=\"get\" style=\"margin:0;padding:.5rem 1rem;\">get <em class=\"small\">method</em></h2>\nfor
  backwards compatability\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
  class=\"admonition-title\">\"get <em class=\"small\">source</em>\"</p>\n</div>\n<pre
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
  <span class=\"nf\">get</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Post\"</span><span class=\"p\">,</span>
  <span class=\"n\">item</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
  class=\"p\">,</span> <span class=\"n\">default</span><span class=\"p\">:</span>
  <span class=\"n\">Any</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"n\">Any</span><span class=\"p\">:</span>\n                <span class=\"s2\">\"for
  backwards compatability\"</span>\n                <span class=\"k\">return</span>
  <span class=\"nb\">getattr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">,</span> <span class=\"n\">item</span><span class=\"p\">,</span> <span
  class=\"n\">default</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"keys\" style=\"margin:0;padding:.5rem
  1rem;\">keys <em class=\"small\">method</em></h2>\nfor backwards compatability\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"keys
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
  <span class=\"nf\">keys</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Post\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"n\">List</span><span class=\"p\">[</span><span
  class=\"nb\">str</span><span class=\"p\">]:</span>\n                <span class=\"s2\">\"for
  backwards compatability\"</span>\n                <span class=\"k\">return</span>
  <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"vm\">__dict__</span><span
  class=\"o\">.</span><span class=\"n\">keys</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"yaml\" style=\"margin:0;padding:.5rem
  1rem;\">yaml <em class=\"small\">method</em></h2>\ndump model to yaml\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"yaml
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
  <span class=\"nf\">yaml</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Post\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n<span
  class=\"w\">                </span><span class=\"sd\">\"\"\"</span>\n<span class=\"sd\">
  \               dump model to yaml</span>\n<span class=\"sd\">                \"\"\"</span>\n
  \               <span class=\"kn\">import</span> <span class=\"nn\">yaml</span>\n\n
  \               <span class=\"k\">return</span> <span class=\"n\">yaml</span><span
  class=\"o\">.</span><span class=\"n\">dump</span><span class=\"p\">(</span>\n                    <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">dict</span><span
  class=\"p\">(</span>\n                        <span class=\"n\">include</span><span
  class=\"o\">=</span><span class=\"p\">{</span><span class=\"n\">i</span><span class=\"p\">:</span>
  <span class=\"kc\">True</span> <span class=\"k\">for</span> <span class=\"n\">i</span>
  <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">post_model</span><span class=\"o\">.</span><span
  class=\"n\">include</span><span class=\"p\">}</span>\n                    <span
  class=\"p\">),</span>\n                    <span class=\"n\">Dumper</span><span
  class=\"o\">=</span><span class=\"n\">yaml</span><span class=\"o\">.</span><span
  class=\"n\">CDumper</span><span class=\"p\">,</span>\n                <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"markdown\" style=\"margin:0;padding:.5rem
  1rem;\">markdown <em class=\"small\">method</em></h2>\ndump model to markdown\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"markdown
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
  <span class=\"nf\">markdown</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Post\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n<span
  class=\"w\">                </span><span class=\"sd\">\"\"\"</span>\n<span class=\"sd\">
  \               dump model to markdown</span>\n<span class=\"sd\">                \"\"\"</span>\n\n
  \               <span class=\"kn\">import</span> <span class=\"nn\">yaml</span>\n\n
  \               <span class=\"n\">frontmatter</span> <span class=\"o\">=</span>
  <span class=\"n\">yaml</span><span class=\"o\">.</span><span class=\"n\">dump</span><span
  class=\"p\">(</span>\n                    <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">dict</span><span class=\"p\">(</span>\n                        <span
  class=\"n\">include</span><span class=\"o\">=</span><span class=\"p\">{</span>\n
  \                           <span class=\"n\">i</span><span class=\"p\">:</span>
  <span class=\"kc\">True</span>\n                            <span class=\"k\">for</span>
  <span class=\"n\">i</span> <span class=\"ow\">in</span> <span class=\"p\">[</span>\n
  \                               <span class=\"n\">_i</span>\n                                <span
  class=\"k\">for</span> <span class=\"n\">_i</span> <span class=\"ow\">in</span>
  <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
  class=\"n\">post_model</span><span class=\"o\">.</span><span class=\"n\">include</span>\n
  \                               <span class=\"k\">if</span> <span class=\"n\">_i</span>
  <span class=\"o\">!=</span> <span class=\"s2\">\"content\"</span>\n                            <span
  class=\"p\">]</span>\n                        <span class=\"p\">}</span>\n                    <span
  class=\"p\">),</span>\n                    <span class=\"n\">Dumper</span><span
  class=\"o\">=</span><span class=\"n\">yaml</span><span class=\"o\">.</span><span
  class=\"n\">CDumper</span><span class=\"p\">,</span>\n                <span class=\"p\">)</span>\n
  \               <span class=\"n\">post</span> <span class=\"o\">=</span> <span class=\"s2\">\"---</span><span
  class=\"se\">\\n</span><span class=\"s2\">\"</span>\n                <span class=\"n\">post</span>
  <span class=\"o\">+=</span> <span class=\"n\">frontmatter</span>\n                <span
  class=\"n\">post</span> <span class=\"o\">+=</span> <span class=\"s2\">\"---</span><span
  class=\"se\">\\n\\n</span><span class=\"s2\">\"</span>\n\n                <span
  class=\"k\">if</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">content</span><span class=\"p\">:</span>\n                    <span
  class=\"n\">post</span> <span class=\"o\">+=</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">content</span>\n                <span class=\"k\">return</span>
  <span class=\"n\">post</span>\n</pre></div>\n\n</pre>\n<p>!! method </p><h2 class=\"admonition-title\"
  id=\"parse_file\" style=\"margin:0;padding:.5rem 1rem;\">parse_file <em class=\"small\">method</em></h2>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"parse_file
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
  <span class=\"nf\">parse_file</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
  class=\"p\">,</span> <span class=\"n\">markata</span><span class=\"p\">,</span>
  <span class=\"n\">path</span><span class=\"p\">:</span> <span class=\"n\">Union</span><span
  class=\"p\">[</span><span class=\"n\">Path</span><span class=\"p\">,</span> <span
  class=\"nb\">str</span><span class=\"p\">],</span> <span class=\"o\">**</span><span
  class=\"n\">kwargs</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"s2\">\"Post\"</span><span class=\"p\">:</span>\n                <span
  class=\"k\">if</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
  class=\"n\">path</span><span class=\"p\">,</span> <span class=\"n\">Path</span><span
  class=\"p\">):</span>\n                    <span class=\"k\">if</span> <span class=\"n\">path</span><span
  class=\"o\">.</span><span class=\"n\">suffix</span> <span class=\"ow\">in</span>
  <span class=\"p\">[</span><span class=\"s2\">\".md\"</span><span class=\"p\">,</span>
  <span class=\"s2\">\".markdown\"</span><span class=\"p\">]:</span>\n                        <span
  class=\"k\">return</span> <span class=\"bp\">cls</span><span class=\"o\">.</span><span
  class=\"n\">parse_markdown</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
  class=\"o\">=</span><span class=\"n\">markata</span><span class=\"p\">,</span> <span
  class=\"n\">path</span><span class=\"o\">=</span><span class=\"n\">path</span><span
  class=\"p\">,</span> <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
  class=\"p\">)</span>\n                <span class=\"k\">elif</span> <span class=\"nb\">isinstance</span><span
  class=\"p\">(</span><span class=\"n\">path</span><span class=\"p\">,</span> <span
  class=\"nb\">str</span><span class=\"p\">):</span>\n                    <span class=\"k\">if</span>
  <span class=\"n\">path</span><span class=\"o\">.</span><span class=\"n\">endswith</span><span
  class=\"p\">(</span><span class=\"s2\">\".md\"</span><span class=\"p\">)</span>
  <span class=\"ow\">or</span> <span class=\"n\">path</span><span class=\"o\">.</span><span
  class=\"n\">endswith</span><span class=\"p\">(</span><span class=\"s2\">\".markdown\"</span><span
  class=\"p\">):</span>\n                        <span class=\"k\">return</span> <span
  class=\"bp\">cls</span><span class=\"o\">.</span><span class=\"n\">parse_markdown</span><span
  class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
  class=\"n\">markata</span><span class=\"p\">,</span> <span class=\"n\">path</span><span
  class=\"o\">=</span><span class=\"n\">path</span><span class=\"p\">,</span> <span
  class=\"o\">**</span><span class=\"n\">kwargs</span><span class=\"p\">)</span>\n
  \               <span class=\"k\">return</span> <span class=\"nb\">super</span><span
  class=\"p\">(</span><span class=\"n\">Post</span><span class=\"p\">,</span> <span
  class=\"bp\">cls</span><span class=\"p\">)</span><span class=\"o\">.</span><span
  class=\"n\">parse_file</span><span class=\"p\">(</span><span class=\"n\">path</span><span
  class=\"p\">,</span> <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
  class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!! method </p><h2 class=\"admonition-title\"
  id=\"parse_markdown\" style=\"margin:0;padding:.5rem 1rem;\">parse_markdown <em
  class=\"small\">method</em></h2>\n<div class=\"admonition source is-collapsible
  collapsible-open\">\n<p class=\"admonition-title\">\"parse_markdown <em class=\"small\">source</em>\"</p>\n</div>\n<pre
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
  <span class=\"nf\">parse_markdown</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
  class=\"p\">,</span> <span class=\"n\">markata</span><span class=\"p\">,</span>
  <span class=\"n\">path</span><span class=\"p\">:</span> <span class=\"n\">Union</span><span
  class=\"p\">[</span><span class=\"n\">Path</span><span class=\"p\">,</span> <span
  class=\"nb\">str</span><span class=\"p\">],</span> <span class=\"o\">**</span><span
  class=\"n\">kwargs</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"s2\">\"Post\"</span><span class=\"p\">:</span>\n                <span
  class=\"k\">if</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
  class=\"n\">path</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
  class=\"p\">):</span>\n                    <span class=\"n\">path</span> <span class=\"o\">=</span>
  <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">path</span><span
  class=\"p\">)</span>\n                <span class=\"n\">text</span> <span class=\"o\">=</span>
  <span class=\"n\">path</span><span class=\"o\">.</span><span class=\"n\">read_text</span><span
  class=\"p\">()</span>\n                <span class=\"k\">try</span><span class=\"p\">:</span>\n
  \                   <span class=\"n\">_</span><span class=\"p\">,</span> <span class=\"n\">fm</span><span
  class=\"p\">,</span> <span class=\"o\">*</span><span class=\"n\">content</span>
  <span class=\"o\">=</span> <span class=\"n\">text</span><span class=\"o\">.</span><span
  class=\"n\">split</span><span class=\"p\">(</span><span class=\"s2\">\"---</span><span
  class=\"se\">\\n</span><span class=\"s2\">\"</span><span class=\"p\">)</span>\n
  \                   <span class=\"n\">content</span> <span class=\"o\">=</span>
  <span class=\"s2\">\"---</span><span class=\"se\">\\n</span><span class=\"s2\">\"</span><span
  class=\"o\">.</span><span class=\"n\">join</span><span class=\"p\">(</span><span
  class=\"n\">content</span><span class=\"p\">)</span>\n                    <span
  class=\"k\">try</span><span class=\"p\">:</span>\n                        <span
  class=\"n\">fm</span> <span class=\"o\">=</span> <span class=\"n\">yaml</span><span
  class=\"o\">.</span><span class=\"n\">load</span><span class=\"p\">(</span><span
  class=\"n\">fm</span><span class=\"p\">,</span> <span class=\"n\">Loader</span><span
  class=\"o\">=</span><span class=\"n\">yaml</span><span class=\"o\">.</span><span
  class=\"n\">CBaseLoader</span><span class=\"p\">)</span>\n                    <span
  class=\"k\">except</span> <span class=\"n\">yaml</span><span class=\"o\">.</span><span
  class=\"n\">YAMLError</span><span class=\"p\">:</span>\n                        <span
  class=\"n\">fm</span> <span class=\"o\">=</span> <span class=\"p\">{}</span>\n                <span
  class=\"k\">except</span> <span class=\"ne\">ValueError</span><span class=\"p\">:</span>\n
  \                   <span class=\"n\">fm</span> <span class=\"o\">=</span> <span
  class=\"p\">{}</span>\n                    <span class=\"n\">content</span> <span
  class=\"o\">=</span> <span class=\"n\">text</span>\n                <span class=\"k\">if</span>
  <span class=\"n\">fm</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span>
  <span class=\"ow\">or</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
  class=\"n\">fm</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
  class=\"p\">):</span>\n                    <span class=\"n\">fm</span> <span class=\"o\">=</span>
  <span class=\"p\">{}</span>\n\n                <span class=\"n\">post_args</span>
  <span class=\"o\">=</span> <span class=\"p\">{</span>\n                    <span
  class=\"s2\">\"markata\"</span><span class=\"p\">:</span> <span class=\"n\">markata</span><span
  class=\"p\">,</span>\n                    <span class=\"s2\">\"path\"</span><span
  class=\"p\">:</span> <span class=\"n\">path</span><span class=\"p\">,</span>\n                    <span
  class=\"s2\">\"content\"</span><span class=\"p\">:</span> <span class=\"n\">content</span><span
  class=\"p\">,</span>\n                    <span class=\"o\">**</span><span class=\"n\">fm</span><span
  class=\"p\">,</span>\n                <span class=\"p\">}</span>\n\n                <span
  class=\"k\">return</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">Post</span><span class=\"p\">(</span><span class=\"o\">**</span><span
  class=\"n\">post_args</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"dumps\" style=\"margin:0;padding:.5rem
  1rem;\">dumps <em class=\"small\">method</em></h2>\ndumps raw article back out\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"dumps
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
  <span class=\"nf\">dumps</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">):</span>\n<span class=\"w\">                </span><span class=\"sd\">\"\"\"</span>\n<span
  class=\"sd\">                dumps raw article back out</span>\n<span class=\"sd\">
  \               \"\"\"</span>\n                <span class=\"k\">return</span> <span
  class=\"sa\">f</span><span class=\"s2\">\"---</span><span class=\"se\">\\n</span><span
  class=\"si\">{</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">yaml</span><span class=\"p\">()</span><span class=\"si\">}</span><span
  class=\"se\">\\n\\n</span><span class=\"s2\">---</span><span class=\"se\">\\n\\n</span><span
  class=\"si\">{</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">content</span><span class=\"si\">}</span><span class=\"s2\">\"</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"default_slug\" style=\"margin:0;padding:.5rem
  1rem;\">default_slug <em class=\"small\">method</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"default_slug
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
  <span class=\"nf\">default_slug</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
  class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
  class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
  class=\"p\">):</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span>
  <span class=\"ow\">or</span> <span class=\"n\">slugify</span><span class=\"p\">(</span><span
  class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">values</span><span
  class=\"p\">[</span><span class=\"s2\">\"path\"</span><span class=\"p\">]</span><span
  class=\"o\">.</span><span class=\"n\">stem</span><span class=\"p\">))</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"index_slug_is_empty\" style=\"margin:0;padding:.5rem
  1rem;\">index_slug_is_empty <em class=\"small\">method</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"index_slug_is_empty
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
  <span class=\"nf\">index_slug_is_empty</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
  class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
  class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
  class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span>
  <span class=\"o\">==</span> <span class=\"s2\">\"index\"</span><span class=\"p\">:</span>\n
  \                   <span class=\"k\">return</span> <span class=\"s2\">\"\"</span>\n
  \               <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"default_href\" style=\"margin:0;padding:.5rem
  1rem;\">default_href <em class=\"small\">method</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"default_href
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
  <span class=\"nf\">default_href</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
  class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
  class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
  class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span><span
  class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
  class=\"n\">v</span>\n                <span class=\"k\">return</span> <span class=\"sa\">f</span><span
  class=\"s2\">\"/</span><span class=\"si\">{</span><span class=\"n\">values</span><span
  class=\"p\">[</span><span class=\"s1\">'slug'</span><span class=\"p\">]</span><span
  class=\"o\">.</span><span class=\"n\">strip</span><span class=\"p\">(</span><span
  class=\"s1\">'/'</span><span class=\"p\">)</span><span class=\"si\">}</span><span
  class=\"s2\">/\"</span><span class=\"o\">.</span><span class=\"n\">replace</span><span
  class=\"p\">(</span><span class=\"s2\">\"//\"</span><span class=\"p\">,</span> <span
  class=\"s2\">\"/\"</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"title_title\" style=\"margin:0;padding:.5rem
  1rem;\">title_title <em class=\"small\">method</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"title_title
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
  <span class=\"nf\">title_title</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
  class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
  class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
  class=\"p\">):</span>\n                <span class=\"n\">title</span> <span class=\"o\">=</span>
  <span class=\"n\">v</span> <span class=\"ow\">or</span> <span class=\"n\">Path</span><span
  class=\"p\">(</span><span class=\"n\">values</span><span class=\"p\">[</span><span
  class=\"s2\">\"path\"</span><span class=\"p\">])</span><span class=\"o\">.</span><span
  class=\"n\">stem</span><span class=\"o\">.</span><span class=\"n\">replace</span><span
  class=\"p\">(</span><span class=\"s2\">\"-\"</span><span class=\"p\">,</span> <span
  class=\"s2\">\" \"</span><span class=\"p\">)</span>\n                <span class=\"k\">return</span>
  <span class=\"n\">title</span><span class=\"o\">.</span><span class=\"n\">title</span><span
  class=\"p\">()</span>\n</pre></div>\n\n</pre>\n<p>!! method </p><h2 class=\"admonition-title\"
  id=\"dateparser_datetime\" style=\"margin:0;padding:.5rem 1rem;\">dateparser_datetime
  <em class=\"small\">method</em></h2>\n<div class=\"admonition source is-collapsible
  collapsible-open\">\n<p class=\"admonition-title\">\"dateparser_datetime <em class=\"small\">source</em>\"</p>\n</div>\n<pre
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
  <span class=\"nf\">dateparser_datetime</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
  class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
  class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
  class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
  class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
  class=\"p\">):</span>\n                    <span class=\"n\">d</span> <span class=\"o\">=</span>
  <span class=\"n\">dateparser</span><span class=\"o\">.</span><span class=\"n\">parse</span><span
  class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">)</span>\n                    <span
  class=\"k\">if</span> <span class=\"n\">d</span> <span class=\"ow\">is</span> <span
  class=\"kc\">None</span><span class=\"p\">:</span>\n                        <span
  class=\"k\">raise</span> <span class=\"ne\">ValueError</span><span class=\"p\">(</span><span
  class=\"sa\">f</span><span class=\"s1\">'\"</span><span class=\"si\">{</span><span
  class=\"n\">v</span><span class=\"si\">}</span><span class=\"s1\">\" is not a valid
  date'</span><span class=\"p\">)</span>\n                <span class=\"k\">return</span>
  <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n<p>!! method </p><h2 class=\"admonition-title\"
  id=\"date_is_datetime\" style=\"margin:0;padding:.5rem 1rem;\">date_is_datetime
  <em class=\"small\">method</em></h2>\n<div class=\"admonition source is-collapsible
  collapsible-open\">\n<p class=\"admonition-title\">\"date_is_datetime <em class=\"small\">source</em>\"</p>\n</div>\n<pre
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
  <span class=\"nf\">date_is_datetime</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
  class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
  class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
  class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span>
  <span class=\"ow\">is</span> <span class=\"kc\">None</span> <span class=\"ow\">and</span>
  <span class=\"s2\">\"date\"</span> <span class=\"ow\">not</span> <span class=\"ow\">in</span>
  <span class=\"n\">values</span><span class=\"p\">:</span>\n                    <span
  class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">\"markata\"</span><span
  class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">console</span><span
  class=\"o\">.</span><span class=\"n\">log</span><span class=\"p\">(</span><span
  class=\"sa\">f</span><span class=\"s2\">\"</span><span class=\"si\">{</span><span
  class=\"n\">values</span><span class=\"p\">[</span><span class=\"s1\">'path'</span><span
  class=\"p\">]</span><span class=\"si\">}</span><span class=\"s2\"> has no date\"</span><span
  class=\"p\">)</span>\n                    <span class=\"k\">return</span> <span
  class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
  class=\"o\">.</span><span class=\"n\">now</span><span class=\"p\">()</span>\n                <span
  class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"ow\">is</span> <span
  class=\"kc\">None</span> <span class=\"ow\">and</span> <span class=\"n\">values</span><span
  class=\"p\">[</span><span class=\"s2\">\"date\"</span><span class=\"p\">]</span>
  <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \                   <span class=\"n\">values</span><span class=\"p\">[</span><span
  class=\"s2\">\"markata\"</span><span class=\"p\">]</span><span class=\"o\">.</span><span
  class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
  class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">\"</span><span
  class=\"si\">{</span><span class=\"n\">values</span><span class=\"p\">[</span><span
  class=\"s1\">'path'</span><span class=\"p\">]</span><span class=\"si\">}</span><span
  class=\"s2\"> has no date\"</span><span class=\"p\">)</span>\n                    <span
  class=\"k\">return</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
  class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">now</span><span
  class=\"p\">()</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
  class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span class=\"n\">datetime</span><span
  class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"p\">):</span>\n
  \                   <span class=\"k\">return</span> <span class=\"n\">v</span>\n
  \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
  class=\"p\">(</span><span class=\"n\">values</span><span class=\"p\">[</span><span
  class=\"s2\">\"date\"</span><span class=\"p\">],</span> <span class=\"n\">datetime</span><span
  class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"p\">):</span>\n
  \                   <span class=\"k\">return</span> <span class=\"n\">values</span><span
  class=\"p\">[</span><span class=\"s2\">\"date\"</span><span class=\"p\">]</span>\n
  \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
  class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span class=\"n\">datetime</span><span
  class=\"o\">.</span><span class=\"n\">date</span><span class=\"p\">):</span>\n                    <span
  class=\"k\">return</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
  class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">combine</span><span
  class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span class=\"n\">datetime</span><span
  class=\"o\">.</span><span class=\"n\">time</span><span class=\"o\">.</span><span
  class=\"n\">min</span><span class=\"p\">)</span>\n                <span class=\"k\">if</span>
  <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">values</span><span
  class=\"p\">[</span><span class=\"s2\">\"date\"</span><span class=\"p\">],</span>
  <span class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">date</span><span
  class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
  class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
  class=\"o\">.</span><span class=\"n\">combine</span><span class=\"p\">(</span><span
  class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">\"date\"</span><span
  class=\"p\">],</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
  class=\"n\">time</span><span class=\"o\">.</span><span class=\"n\">min</span><span
  class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"mindate_time\" style=\"margin:0;padding:.5rem
  1rem;\">mindate_time <em class=\"small\">method</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"mindate_time
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
  <span class=\"nf\">mindate_time</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
  class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
  class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
  class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span>
  <span class=\"ow\">is</span> <span class=\"kc\">None</span> <span class=\"ow\">and</span>
  <span class=\"s2\">\"date\"</span> <span class=\"ow\">not</span> <span class=\"ow\">in</span>
  <span class=\"n\">values</span><span class=\"p\">:</span>\n                    <span
  class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">\"markata\"</span><span
  class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">console</span><span
  class=\"o\">.</span><span class=\"n\">log</span><span class=\"p\">(</span><span
  class=\"sa\">f</span><span class=\"s2\">\"</span><span class=\"si\">{</span><span
  class=\"n\">values</span><span class=\"p\">[</span><span class=\"s1\">'path'</span><span
  class=\"p\">]</span><span class=\"si\">}</span><span class=\"s2\"> has no date\"</span><span
  class=\"p\">)</span>\n                    <span class=\"k\">return</span> <span
  class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
  class=\"o\">.</span><span class=\"n\">min</span>\n                <span class=\"k\">if</span>
  <span class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">\"date\"</span><span
  class=\"p\">]</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n                    <span class=\"n\">values</span><span class=\"p\">[</span><span
  class=\"s2\">\"markata\"</span><span class=\"p\">]</span><span class=\"o\">.</span><span
  class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
  class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">\"</span><span
  class=\"si\">{</span><span class=\"n\">values</span><span class=\"p\">[</span><span
  class=\"s1\">'path'</span><span class=\"p\">]</span><span class=\"si\">}</span><span
  class=\"s2\"> has no date\"</span><span class=\"p\">)</span>\n                    <span
  class=\"k\">return</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
  class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">min</span>\n
  \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
  class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span class=\"n\">datetime</span><span
  class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"p\">):</span>\n
  \                   <span class=\"k\">return</span> <span class=\"n\">v</span>\n
  \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
  class=\"p\">(</span><span class=\"n\">values</span><span class=\"p\">[</span><span
  class=\"s2\">\"date\"</span><span class=\"p\">],</span> <span class=\"n\">datetime</span><span
  class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"p\">):</span>\n
  \                   <span class=\"k\">return</span> <span class=\"n\">values</span><span
  class=\"p\">[</span><span class=\"s2\">\"date\"</span><span class=\"p\">]</span>\n
  \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
  class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span class=\"n\">datetime</span><span
  class=\"o\">.</span><span class=\"n\">date</span><span class=\"p\">):</span>\n                    <span
  class=\"k\">return</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
  class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">combine</span><span
  class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span class=\"n\">datetime</span><span
  class=\"o\">.</span><span class=\"n\">time</span><span class=\"o\">.</span><span
  class=\"n\">min</span><span class=\"p\">)</span>\n                <span class=\"k\">if</span>
  <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">values</span><span
  class=\"p\">[</span><span class=\"s2\">\"date\"</span><span class=\"p\">],</span>
  <span class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">date</span><span
  class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
  class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
  class=\"o\">.</span><span class=\"n\">combine</span><span class=\"p\">(</span><span
  class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">\"date\"</span><span
  class=\"p\">],</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
  class=\"n\">time</span><span class=\"o\">.</span><span class=\"n\">min</span><span
  class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"dateparser_date\" style=\"margin:0;padding:.5rem
  1rem;\">dateparser_date <em class=\"small\">method</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"dateparser_date
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
  <span class=\"nf\">dateparser_date</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
  class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
  class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
  class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span>
  <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \                   <span class=\"k\">return</span> <span class=\"n\">datetime</span><span
  class=\"o\">.</span><span class=\"n\">date</span><span class=\"o\">.</span><span
  class=\"n\">min</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
  class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
  class=\"p\">):</span>\n                    <span class=\"n\">d</span> <span class=\"o\">=</span>
  <span class=\"bp\">cls</span><span class=\"o\">.</span><span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">precache</span><span class=\"o\">.</span><span
  class=\"n\">get</span><span class=\"p\">(</span><span class=\"n\">v</span><span
  class=\"p\">)</span>\n                    <span class=\"k\">if</span> <span class=\"n\">d</span>
  <span class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n                        <span class=\"k\">return</span> <span
  class=\"n\">d</span>\n                    <span class=\"n\">d</span> <span class=\"o\">=</span>
  <span class=\"n\">dateparser</span><span class=\"o\">.</span><span class=\"n\">parse</span><span
  class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">)</span>\n                    <span
  class=\"k\">if</span> <span class=\"n\">d</span> <span class=\"ow\">is</span> <span
  class=\"kc\">None</span><span class=\"p\">:</span>\n                        <span
  class=\"k\">raise</span> <span class=\"ne\">ValueError</span><span class=\"p\">(</span><span
  class=\"sa\">f</span><span class=\"s1\">'\"</span><span class=\"si\">{</span><span
  class=\"n\">v</span><span class=\"si\">}</span><span class=\"s1\">\" is not a valid
  date'</span><span class=\"p\">)</span>\n                    <span class=\"n\">d</span>
  <span class=\"o\">=</span> <span class=\"n\">d</span><span class=\"o\">.</span><span
  class=\"n\">date</span><span class=\"p\">()</span>\n                    <span class=\"k\">with</span>
  <span class=\"bp\">cls</span><span class=\"o\">.</span><span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">cache</span> <span class=\"k\">as</span> <span
  class=\"n\">cache</span><span class=\"p\">:</span>\n                        <span
  class=\"n\">cache</span><span class=\"o\">.</span><span class=\"n\">add</span><span
  class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span class=\"n\">d</span><span
  class=\"p\">)</span>\n                    <span class=\"k\">return</span> <span
  class=\"n\">d</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"__init__\" style=\"margin:0;padding:.5rem
  1rem;\"><strong>init</strong> <em class=\"small\">method</em></h2>\ninclude: post
  attributes to include by default in Post\nmodel serialization.\nrepr_include: post
  attributes to include by default in Post\nrepr.  If <code>repr_include</code> is
  None, it will default to\n<code>include</code>, but it is likely that you want less
  in the repr\nthan serialized output.\n<pre><code>example:\n\n``` toml title='markata.toml'\n[markata.post_model]\ninclude
  = ['date', 'description', 'published',\n    'slug', 'title', 'content', 'html']\nrepr_include
  = ['date', 'description', 'published', 'slug', 'title']\n```\n</code></pre>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"<strong>init</strong>
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
  <span class=\"fm\">__init__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">,</span> <span class=\"o\">**</span><span class=\"n\">data</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n<span class=\"w\">                </span><span class=\"sd\">\"\"\"</span>\n\n<span
  class=\"sd\">                include: post attributes to include by default in Post</span>\n<span
  class=\"sd\">                model serialization.</span>\n<span class=\"sd\">                repr_include:
  post attributes to include by default in Post</span>\n<span class=\"sd\">                repr.
  \ If `repr_include` is None, it will default to</span>\n<span class=\"sd\">                `include`,
  but it is likely that you want less in the repr</span>\n<span class=\"sd\">                than
  serialized output.</span>\n\n<span class=\"sd\">                example:</span>\n\n<span
  class=\"sd\">                ``` toml title='markata.toml'</span>\n<span class=\"sd\">
  \               [markata.post_model]</span>\n<span class=\"sd\">                include
  = ['date', 'description', 'published',</span>\n<span class=\"sd\">                    'slug',
  'title', 'content', 'html']</span>\n<span class=\"sd\">                repr_include
  = ['date', 'description', 'published', 'slug', 'title']</span>\n<span class=\"sd\">
  \               ```</span>\n<span class=\"sd\">                \"\"\"</span>\n                <span
  class=\"nb\">super</span><span class=\"p\">()</span><span class=\"o\">.</span><span
  class=\"fm\">__init__</span><span class=\"p\">(</span><span class=\"o\">**</span><span
  class=\"n\">data</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"repr_include_validator\" style=\"margin:0;padding:.5rem
  1rem;\">repr_include_validator <em class=\"small\">method</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"repr_include_validator
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
  <span class=\"nf\">repr_include_validator</span><span class=\"p\">(</span><span
  class=\"bp\">cls</span><span class=\"p\">,</span> <span class=\"n\">v</span><span
  class=\"p\">,</span> <span class=\"o\">*</span><span class=\"p\">,</span> <span
  class=\"n\">values</span><span class=\"p\">):</span>\n                <span class=\"k\">if</span>
  <span class=\"n\">v</span><span class=\"p\">:</span>\n                    <span
  class=\"k\">return</span> <span class=\"n\">v</span>\n                <span class=\"k\">return</span>
  <span class=\"n\">values</span><span class=\"o\">.</span><span class=\"n\">get</span><span
  class=\"p\">(</span><span class=\"s2\">\"include\"</span><span class=\"p\">,</span>
  <span class=\"kc\">None</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n</main>\n<footer>\xA9
  2024</footer>\n</body></html>"
published: true
slug: markata/plugins/post-model
title: Post_Model.Py


---

None


!! class <h2 id='Post' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Post <em class='small'>class</em></h2>

???+ source "Post <em class='small'>source</em>"

```python

        class Post(pydantic.BaseModel):
            markata: Any = None
            path: Path
            slug: Optional[str] = None
            href: Optional[str] = None
            published: bool = False
            description: Optional[str] = None
            content: str = None
            # date: Union[datetime.date, str]=None
            date: Optional[Union[datetime.date, str]] = None
            # pydantic.Field(
            # default_factory=lambda: datetime.date.min
            # )
            date_time: Optional[datetime.datetime] = None
            today: datetime.date = pydantic.Field(default_factory=datetime.date.today)
            now: datetime.datetime = pydantic.Field(default_factory=datetime.datetime.utcnow)
            load_time: float = 0
            profile: Optional[str] = None
            title: str = None
            model_config = ConfigDict(
                validate_assignment=True,
                arbitrary_types_allowed=True,
                extra="allow",
            )

            def __repr_args__(self: "Post") -> "ReprArgs":
                return [
                    (key, value)
                    for key, value in self.__dict__.items()
                    if key in self.markata.config.post_model.repr_include
                ]

            @property
            def metadata(self: "Post") -> Dict:
                "for backwards compatability"
                return self.__dict__

            def to_dict(self: "Post") -> Dict:
                "for backwards compatability"
                return self.__dict__

            def __getitem__(self: "Post", item: str) -> Any:
                "for backwards compatability"
                return getattr(self, item)

            def __setitem__(self: "Post", key: str, item: Any) -> None:
                "for backwards compatability"
                setattr(self, key, item)

            def get(self: "Post", item: str, default: Any) -> Any:
                "for backwards compatability"
                return getattr(self, item, default)

            def keys(self: "Post") -> List[str]:
                "for backwards compatability"
                return self.__dict__.keys()

            # def json(
            #     self: "Post",
            #     include: Iterable = None,
            #     all: bool = False,
            #     **kwargs,
            # ) -> str:
            #     """
            #     override function to give a default include value that will include
            #     user configured includes.
            #     """
            #     if all:
            #         return pydantic.create_model("Post", **self)(**self).json(
            #             **kwargs,
            #         )
            #     if include:
            #         return pydantic.create_model("Post", **self)(**self).json(
            #             include=include,
            #             **kwargs,
            #         )
            #     return pydantic.create_model("Post", **self)(**self).json(
            #         include={i: True for i in self.markata.config.post_model.include},
            #         **kwargs,
            #     )

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

            @classmethod
            def parse_file(cls, markata, path: Union[Path, str], **kwargs) -> "Post":
                if isinstance(path, Path):
                    if path.suffix in [".md", ".markdown"]:
                        return cls.parse_markdown(markata=markata, path=path, **kwargs)
                elif isinstance(path, str):
                    if path.endswith(".md") or path.endswith(".markdown"):
                        return cls.parse_markdown(markata=markata, path=path, **kwargs)
                return super(Post, cls).parse_file(path, **kwargs)

            @classmethod
            def parse_markdown(cls, markata, path: Union[Path, str], **kwargs) -> "Post":
                if isinstance(path, str):
                    path = Path(path)
                text = path.read_text()
                try:
                    _, fm, *content = text.split("---\n")
                    content = "---\n".join(content)
                    try:
                        fm = yaml.load(fm, Loader=yaml.CBaseLoader)
                    except yaml.YAMLError:
                        fm = {}
                except ValueError:
                    fm = {}
                    content = text
                if fm is None or isinstance(fm, str):
                    fm = {}

                post_args = {
                    "markata": markata,
                    "path": path,
                    "content": content,
                    **fm,
                }

                return markata.Post(**post_args)

            def dumps(self):
                """
                dumps raw article back out
                """
                return f"---\n{self.yaml()}\n\n---\n\n{self.content}"

            @pydantic.validator("slug", pre=True, always=True)
            def default_slug(cls, v, *, values):
                return v or slugify(str(values["path"].stem))

            @pydantic.validator("slug", pre=True, always=True)
            def index_slug_is_empty(cls, v, *, values):
                if v == "index":
                    return ""
                return v

            @pydantic.validator("href", pre=True, always=True)
            def default_href(cls, v, *, values):
                if v:
                    return v
                return f"/{values['slug'].strip('/')}/".replace("//", "/")

            @pydantic.validator("title", pre=True, always=True)
            def title_title(cls, v, *, values):
                title = v or Path(values["path"]).stem.replace("-", " ")
                return title.title()

            @pydantic.validator("date_time", pre=True, always=True)
            def dateparser_datetime(cls, v, *, values):
                if isinstance(v, str):
                    d = dateparser.parse(v)
                    if d is None:
                        raise ValueError(f'"{v}" is not a valid date')
                return v

            @pydantic.validator("date_time", pre=True, always=True)
            def date_is_datetime(cls, v, *, values):
                if v is None and "date" not in values:
                    values["markata"].console.log(f"{values['path']} has no date")
                    return datetime.datetime.now()
                if v is None and values["date"] is None:
                    values["markata"].console.log(f"{values['path']} has no date")
                    return datetime.datetime.now()
                if isinstance(v, datetime.datetime):
                    return v
                if isinstance(values["date"], datetime.datetime):
                    return values["date"]
                if isinstance(v, datetime.date):
                    return datetime.datetime.combine(v, datetime.time.min)
                if isinstance(values["date"], datetime.date):
                    return datetime.datetime.combine(values["date"], datetime.time.min)
                return v

            @pydantic.validator("date_time", pre=True, always=True)
            def mindate_time(cls, v, *, values):
                if v is None and "date" not in values:
                    values["markata"].console.log(f"{values['path']} has no date")
                    return datetime.datetime.min
                if values["date"] is None:
                    values["markata"].console.log(f"{values['path']} has no date")
                    return datetime.datetime.min
                if isinstance(v, datetime.datetime):
                    return v
                if isinstance(values["date"], datetime.datetime):
                    return values["date"]
                if isinstance(v, datetime.date):
                    return datetime.datetime.combine(v, datetime.time.min)
                if isinstance(values["date"], datetime.date):
                    return datetime.datetime.combine(values["date"], datetime.time.min)
                return v

            @pydantic.validator("date", pre=True, always=True)
            def dateparser_date(cls, v, *, values):
                if v is None:
                    return datetime.date.min
                if isinstance(v, str):
                    d = cls.markata.precache.get(v)
                    if d is not None:
                        return d
                    d = dateparser.parse(v)
                    if d is None:
                        raise ValueError(f'"{v}" is not a valid date')
                    d = d.date()
                    with cls.markata.cache as cache:
                        cache.add(v, d)
                    return d
                return v
```


!! class <h2 id='PostModelConfig' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>PostModelConfig <em class='small'>class</em></h2>
    Configuration for the Post model
???+ source "PostModelConfig <em class='small'>source</em>"

```python

        class PostModelConfig(pydantic.BaseModel):
            "Configuration for the Post model"

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

            @pydantic.validator("repr_include", pre=True, always=True)
            def repr_include_validator(cls, v, *, values):
                if v:
                    return v
                return values.get("include", None)
```


!! class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Config <em class='small'>class</em></h2>

???+ source "Config <em class='small'>source</em>"

```python

        class Config(pydantic.BaseModel):
            post_model: PostModelConfig = pydantic.Field(default_factory=PostModelConfig)
```


!! function <h2 id='post_model' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>post_model <em class='small'>function</em></h2>

???+ source "post_model <em class='small'>source</em>"

```python

        def post_model(markata: "Markata") -> None:
            markata.post_models.append(Post)
```


!! function <h2 id='config_model' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>config_model <em class='small'>function</em></h2>

???+ source "config_model <em class='small'>source</em>"

```python

        def config_model(markata: "Markata") -> None:
            markata.config_models.append(Config)
```


!! class <h2 id='PostFactory' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>PostFactory <em class='small'>class</em></h2>

???+ source "PostFactory <em class='small'>source</em>"

```python

        class PostFactory(ModelFactory):
            __model__ = Post
```


!! method <h2 id='__repr_args__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__repr_args__ <em class='small'>method</em></h2>

???+ source "__repr_args__ <em class='small'>source</em>"

```python

        def __repr_args__(self: "Post") -> "ReprArgs":
                return [
                    (key, value)
                    for key, value in self.__dict__.items()
                    if key in self.markata.config.post_model.repr_include
                ]
```


!! method <h2 id='metadata' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>metadata <em class='small'>method</em></h2>
    for backwards compatability
???+ source "metadata <em class='small'>source</em>"

```python

        def metadata(self: "Post") -> Dict:
                "for backwards compatability"
                return self.__dict__
```


!! method <h2 id='to_dict' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>to_dict <em class='small'>method</em></h2>
    for backwards compatability
???+ source "to_dict <em class='small'>source</em>"

```python

        def to_dict(self: "Post") -> Dict:
                "for backwards compatability"
                return self.__dict__
```


!! method <h2 id='__getitem__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__getitem__ <em class='small'>method</em></h2>
    for backwards compatability
???+ source "__getitem__ <em class='small'>source</em>"

```python

        def __getitem__(self: "Post", item: str) -> Any:
                "for backwards compatability"
                return getattr(self, item)
```


!! method <h2 id='__setitem__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__setitem__ <em class='small'>method</em></h2>
    for backwards compatability
???+ source "__setitem__ <em class='small'>source</em>"

```python

        def __setitem__(self: "Post", key: str, item: Any) -> None:
                "for backwards compatability"
                setattr(self, key, item)
```


!! method <h2 id='get' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get <em class='small'>method</em></h2>
    for backwards compatability
???+ source "get <em class='small'>source</em>"

```python

        def get(self: "Post", item: str, default: Any) -> Any:
                "for backwards compatability"
                return getattr(self, item, default)
```


!! method <h2 id='keys' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>keys <em class='small'>method</em></h2>
    for backwards compatability
???+ source "keys <em class='small'>source</em>"

```python

        def keys(self: "Post") -> List[str]:
                "for backwards compatability"
                return self.__dict__.keys()
```


!! method <h2 id='yaml' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>yaml <em class='small'>method</em></h2>
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


!! method <h2 id='markdown' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>markdown <em class='small'>method</em></h2>
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


!! method <h2 id='parse_file' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>parse_file <em class='small'>method</em></h2>

???+ source "parse_file <em class='small'>source</em>"

```python

        def parse_file(cls, markata, path: Union[Path, str], **kwargs) -> "Post":
                if isinstance(path, Path):
                    if path.suffix in [".md", ".markdown"]:
                        return cls.parse_markdown(markata=markata, path=path, **kwargs)
                elif isinstance(path, str):
                    if path.endswith(".md") or path.endswith(".markdown"):
                        return cls.parse_markdown(markata=markata, path=path, **kwargs)
                return super(Post, cls).parse_file(path, **kwargs)
```


!! method <h2 id='parse_markdown' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>parse_markdown <em class='small'>method</em></h2>

???+ source "parse_markdown <em class='small'>source</em>"

```python

        def parse_markdown(cls, markata, path: Union[Path, str], **kwargs) -> "Post":
                if isinstance(path, str):
                    path = Path(path)
                text = path.read_text()
                try:
                    _, fm, *content = text.split("---\n")
                    content = "---\n".join(content)
                    try:
                        fm = yaml.load(fm, Loader=yaml.CBaseLoader)
                    except yaml.YAMLError:
                        fm = {}
                except ValueError:
                    fm = {}
                    content = text
                if fm is None or isinstance(fm, str):
                    fm = {}

                post_args = {
                    "markata": markata,
                    "path": path,
                    "content": content,
                    **fm,
                }

                return markata.Post(**post_args)
```


!! method <h2 id='dumps' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>dumps <em class='small'>method</em></h2>
    dumps raw article back out
???+ source "dumps <em class='small'>source</em>"

```python

        def dumps(self):
                """
                dumps raw article back out
                """
                return f"---\n{self.yaml()}\n\n---\n\n{self.content}"
```


!! method <h2 id='default_slug' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>default_slug <em class='small'>method</em></h2>

???+ source "default_slug <em class='small'>source</em>"

```python

        def default_slug(cls, v, *, values):
                return v or slugify(str(values["path"].stem))
```


!! method <h2 id='index_slug_is_empty' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>index_slug_is_empty <em class='small'>method</em></h2>

???+ source "index_slug_is_empty <em class='small'>source</em>"

```python

        def index_slug_is_empty(cls, v, *, values):
                if v == "index":
                    return ""
                return v
```


!! method <h2 id='default_href' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>default_href <em class='small'>method</em></h2>

???+ source "default_href <em class='small'>source</em>"

```python

        def default_href(cls, v, *, values):
                if v:
                    return v
                return f"/{values['slug'].strip('/')}/".replace("//", "/")
```


!! method <h2 id='title_title' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>title_title <em class='small'>method</em></h2>

???+ source "title_title <em class='small'>source</em>"

```python

        def title_title(cls, v, *, values):
                title = v or Path(values["path"]).stem.replace("-", " ")
                return title.title()
```


!! method <h2 id='dateparser_datetime' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>dateparser_datetime <em class='small'>method</em></h2>

???+ source "dateparser_datetime <em class='small'>source</em>"

```python

        def dateparser_datetime(cls, v, *, values):
                if isinstance(v, str):
                    d = dateparser.parse(v)
                    if d is None:
                        raise ValueError(f'"{v}" is not a valid date')
                return v
```


!! method <h2 id='date_is_datetime' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>date_is_datetime <em class='small'>method</em></h2>

???+ source "date_is_datetime <em class='small'>source</em>"

```python

        def date_is_datetime(cls, v, *, values):
                if v is None and "date" not in values:
                    values["markata"].console.log(f"{values['path']} has no date")
                    return datetime.datetime.now()
                if v is None and values["date"] is None:
                    values["markata"].console.log(f"{values['path']} has no date")
                    return datetime.datetime.now()
                if isinstance(v, datetime.datetime):
                    return v
                if isinstance(values["date"], datetime.datetime):
                    return values["date"]
                if isinstance(v, datetime.date):
                    return datetime.datetime.combine(v, datetime.time.min)
                if isinstance(values["date"], datetime.date):
                    return datetime.datetime.combine(values["date"], datetime.time.min)
                return v
```


!! method <h2 id='mindate_time' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>mindate_time <em class='small'>method</em></h2>

???+ source "mindate_time <em class='small'>source</em>"

```python

        def mindate_time(cls, v, *, values):
                if v is None and "date" not in values:
                    values["markata"].console.log(f"{values['path']} has no date")
                    return datetime.datetime.min
                if values["date"] is None:
                    values["markata"].console.log(f"{values['path']} has no date")
                    return datetime.datetime.min
                if isinstance(v, datetime.datetime):
                    return v
                if isinstance(values["date"], datetime.datetime):
                    return values["date"]
                if isinstance(v, datetime.date):
                    return datetime.datetime.combine(v, datetime.time.min)
                if isinstance(values["date"], datetime.date):
                    return datetime.datetime.combine(values["date"], datetime.time.min)
                return v
```


!! method <h2 id='dateparser_date' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>dateparser_date <em class='small'>method</em></h2>

???+ source "dateparser_date <em class='small'>source</em>"

```python

        def dateparser_date(cls, v, *, values):
                if v is None:
                    return datetime.date.min
                if isinstance(v, str):
                    d = cls.markata.precache.get(v)
                    if d is not None:
                        return d
                    d = dateparser.parse(v)
                    if d is None:
                        raise ValueError(f'"{v}" is not a valid date')
                    d = d.date()
                    with cls.markata.cache as cache:
                        cache.add(v, d)
                    return d
                return v
```


!! method <h2 id='__init__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__init__ <em class='small'>method</em></h2>
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


!! method <h2 id='repr_include_validator' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>repr_include_validator <em class='small'>method</em></h2>

???+ source "repr_include_validator <em class='small'>source</em>"

```python

        def repr_include_validator(cls, v, *, values):
                if v:
                    return v
                return values.get("include", None)
```
