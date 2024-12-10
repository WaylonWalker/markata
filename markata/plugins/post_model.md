---
content: "None\n\n\n!! class <h2 id='Post' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>Post <em class='small'>class</em></h2>\n\n???+ source \"Post <em class='small'>source</em>\"\n\n```python\n\n
  \       class Post(pydantic.BaseModel, JupyterMixin):\n            markata: Any
  = Field(None, exclude=True)\n            path: Path\n            slug: Optional[str]
  = None\n            href: Optional[str] = None\n            published: bool = False\n
  \           description: Optional[str] = None\n            content: str = None\n
  \           # date: Union[datetime.date, str]=None\n            date: Optional[Union[datetime.date,
  str]] = None\n            # pydantic.Field(\n            # default_factory=lambda:
  datetime.date.min\n            # )\n            date_time: Optional[datetime.datetime]
  = None\n            today: datetime.date = pydantic.Field(default_factory=datetime.date.today)\n
  \           now: datetime.datetime = pydantic.Field(default_factory=datetime.datetime.utcnow)\n
  \           load_time: float = 0\n            profile: Optional[str] = None\n            title:
  str = None\n            model_config = ConfigDict(\n                validate_assignment=True,\n
  \               arbitrary_types_allowed=True,\n                extra=\"allow\",\n
  \           )\n            template: Optional[str | Dict[str, str]] = \"post.html\"\n
  \           sidebar: Optional[Any] = None\n\n            def __rich__(self) -> Pretty:\n
  \               return Pretty(self)\n\n            def __repr_args__(self: \"Post\")
  -> \"ReprArgs\":\n                return [\n                    (key, value)\n                    for
  key, value in self.__dict__.items()\n                    if key in self.markata.config.post_model.repr_include\n
  \               ]\n\n            @property\n            def key(self: \"Post\")
  -> List[str]:\n                return self.markata.make_hash(\n                    self.slug,\n
  \                   self.href,\n                    self.published,\n                    self.description,\n
  \                   self.content,\n                    self.date,\n                    self.title,\n
  \                   self.template,\n                )\n\n            @property\n
  \           def metadata(self: \"Post\") -> Dict:\n                \"for backwards
  compatability\"\n                return self.__dict__\n\n            def to_dict(self:
  \"Post\") -> Dict:\n                \"for backwards compatability\"\n                return
  self.__dict__\n\n            def __getitem__(self: \"Post\", item: str) -> Any:\n
  \               \"for backwards compatability\"\n                return getattr(self,
  item)\n\n            def __setitem__(self: \"Post\", key: str, item: Any) -> None:\n
  \               \"for backwards compatability\"\n                setattr(self, key,
  item)\n\n            def get(self: \"Post\", item: str, default: Any) -> Any:\n
  \               \"for backwards compatability\"\n                return getattr(self,
  item, default)\n\n            def keys(self: \"Post\") -> List[str]:\n                \"for
  backwards compatability\"\n                return self.__dict__.keys()\n\n            #
  def json(\n            #     self: \"Post\",\n            #     include: Iterable
  = None,\n            #     all: bool = False,\n            #     **kwargs,\n            #
  ) -> str:\n            #     \"\"\"\n            #     override function to give
  a default include value that will include\n            #     user configured includes.\n
  \           #     \"\"\"\n            #     if all:\n            #         return
  pydantic.create_model(\"Post\", **self)(**self).json(\n            #             **kwargs,\n
  \           #         )\n            #     if include:\n            #         return
  pydantic.create_model(\"Post\", **self)(**self).json(\n            #             include=include,\n
  \           #             **kwargs,\n            #         )\n            #     return
  pydantic.create_model(\"Post\", **self)(**self).json(\n            #         include={i:
  True for i in self.markata.config.post_model.include},\n            #         **kwargs,\n
  \           #     )\n\n            def yaml(self: \"Post\") -> str:\n                \"\"\"\n
  \               dump model to yaml\n                \"\"\"\n                import
  yaml\n\n                return yaml.dump(\n                    self.dict(\n                        include={i:
  True for i in self.markata.config.post_model.include}\n                    ),\n
  \                   Dumper=yaml.CDumper,\n                )\n\n            def markdown(self:
  \"Post\") -> str:\n                \"\"\"\n                dump model to markdown\n
  \               \"\"\"\n\n                import yaml\n\n                frontmatter
  = yaml.dump(\n                    self.dict(\n                        include={\n
  \                           i: True\n                            for i in [\n                                _i\n
  \                               for _i in self.markata.config.post_model.include\n
  \                               if _i != \"content\"\n                            ]\n
  \                       }\n                    ),\n                    Dumper=yaml.CDumper,\n
  \               )\n                post = \"---\\n\"\n                post += frontmatter\n
  \               post += \"---\\n\\n\"\n\n                if self.content:\n                    post
  += self.content\n                return post\n\n            @classmethod\n            def
  parse_file(cls, markata, path: Union[Path, str], **kwargs) -> \"Post\":\n                if
  isinstance(path, Path):\n                    if path.suffix in [\".md\", \".markdown\"]:\n
  \                       return cls.parse_markdown(markata=markata, path=path, **kwargs)\n
  \               elif isinstance(path, str):\n                    if path.endswith(\".md\")
  or path.endswith(\".markdown\"):\n                        return cls.parse_markdown(markata=markata,
  path=path, **kwargs)\n                return super(Post, cls).parse_file(path, **kwargs)\n\n
  \           @classmethod\n            def parse_markdown(cls, markata, path: Union[Path,
  str], **kwargs) -> \"Post\":\n                if isinstance(path, str):\n                    path
  = Path(path)\n                text = path.read_text()\n                try:\n                    _,
  fm, *content = text.split(\"---\\n\")\n                    content = \"---\\n\".join(content)\n
  \                   try:\n                        fm = yaml.load(fm, Loader=yaml.CBaseLoader)\n
  \                   except yaml.YAMLError:\n                        fm = {}\n                except
  ValueError:\n                    fm = {}\n                    content = text\n                if
  fm is None or isinstance(fm, str):\n                    fm = {}\n\n                post_args
  = {\n                    \"markata\": markata,\n                    \"path\": path,\n
  \                   \"content\": content,\n                    \"raw\": text,\n
  \                   **fm,\n                }\n\n                return markata.Post.parse_obj(post_args)\n\n
  \           def dumps(self):\n                \"\"\"\n                dumps raw
  article back out\n                \"\"\"\n                return f\"---\\n{self.yaml()}\\n\\n---\\n\\n{self.content}\"\n\n
  \           @pydantic.validator(\"slug\", pre=True, always=True)\n            def
  default_slug(cls, v, *, values):\n                return v or slugify(str(values[\"path\"].stem))\n\n
  \           @pydantic.validator(\"slug\", pre=True, always=True)\n            def
  index_slug_is_empty(cls, v, *, values):\n                if v == \"index\":\n                    return
  \"\"\n                return v\n\n            @pydantic.validator(\"slug\", pre=True,
  always=True)\n            def no_double_slash_in_slug(cls, v, *, values):\n                if
  v is None:\n                    return v\n                return v.replace(\"//\",
  \"/\")\n\n            @pydantic.validator(\"href\", pre=True, always=True)\n            def
  default_href(cls, v, *, values):\n                if v:\n                    return
  v\n                return f\"/{values['slug'].strip('/')}/\".replace(\"//\", \"/\")\n\n
  \           @pydantic.validator(\"title\", pre=True, always=True)\n            def
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
  \               \"slug\",\n                \"title\",\n            ]\n            export_include:
  Optional[List[str]] = [\n                \"date\",\n                \"description\",\n
  \               \"published\",\n                \"slug\",\n                \"title\",\n
  \           ]\n\n            @pydantic.validator(\"repr_include\", pre=True, always=True)\n
  \           def repr_include_validator(cls, v, *, values):\n                if v:\n
  \                   return v\n                return values.get(\"include\", None)\n```\n\n\n!!
  class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Config
  <em class='small'>class</em></h2>\n\n???+ source \"Config <em class='small'>source</em>\"\n\n```python\n\n
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
  \           __model__ = Post\n            __model__ = Post\n```\n\n\n!! method <h2
  id='__rich__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__rich__
  <em class='small'>method</em></h2>\n\n???+ source \"__rich__ <em class='small'>source</em>\"\n\n```python\n\n
  \       def __rich__(self) -> Pretty:\n                return Pretty(self)\n```\n\n\n!!
  method <h2 id='__repr_args__' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>__repr_args__ <em class='small'>method</em></h2>\n\n???+ source \"__repr_args__
  <em class='small'>source</em>\"\n\n```python\n\n        def __repr_args__(self:
  \"Post\") -> \"ReprArgs\":\n                return [\n                    (key,
  value)\n                    for key, value in self.__dict__.items()\n                    if
  key in self.markata.config.post_model.repr_include\n                ]\n```\n\n\n!!
  method <h2 id='key' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>key
  <em class='small'>method</em></h2>\n\n???+ source \"key <em class='small'>source</em>\"\n\n```python\n\n
  \       def key(self: \"Post\") -> List[str]:\n                return self.markata.make_hash(\n
  \                   self.slug,\n                    self.href,\n                    self.published,\n
  \                   self.description,\n                    self.content,\n                    self.date,\n
  \                   self.title,\n                    self.template,\n                )\n```\n\n\n!!
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
  \                   \"content\": content,\n                    \"raw\": text,\n
  \                   **fm,\n                }\n\n                return markata.Post.parse_obj(post_args)\n```\n\n\n!!
  method <h2 id='dumps' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>dumps
  <em class='small'>method</em></h2>\n    dumps raw article back out\n???+ source
  \"dumps <em class='small'>source</em>\"\n\n```python\n\n        def dumps(self):\n
  \               \"\"\"\n                dumps raw article back out\n                \"\"\"\n
  \               return f\"---\\n{self.yaml()}\\n\\n---\\n\\n{self.content}\"\n```\n\n\n!!
  method <h2 id='default_slug' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>default_slug <em class='small'>method</em></h2>\n\n???+ source \"default_slug
  <em class='small'>source</em>\"\n\n```python\n\n        def default_slug(cls, v,
  *, values):\n                return v or slugify(str(values[\"path\"].stem))\n```\n\n\n!!
  method <h2 id='index_slug_is_empty' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>index_slug_is_empty <em class='small'>method</em></h2>\n\n???+ source \"index_slug_is_empty
  <em class='small'>source</em>\"\n\n```python\n\n        def index_slug_is_empty(cls,
  v, *, values):\n                if v == \"index\":\n                    return \"\"\n
  \               return v\n```\n\n\n!! method <h2 id='no_double_slash_in_slug' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>no_double_slash_in_slug <em class='small'>method</em></h2>\n\n???+
  source \"no_double_slash_in_slug <em class='small'>source</em>\"\n\n```python\n\n
  \       def no_double_slash_in_slug(cls, v, *, values):\n                if v is
  None:\n                    return v\n                return v.replace(\"//\", \"/\")\n```\n\n\n!!
  method <h2 id='default_href' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>default_href <em class='small'>method</em></h2>\n\n???+ source \"default_href
  <em class='small'>source</em>\"\n\n```python\n\n        def default_href(cls, v,
  *, values):\n                if v:\n                    return v\n                return
  f\"/{values['slug'].strip('/')}/\".replace(\"//\", \"/\")\n```\n\n\n!! method <h2
  id='title_title' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>title_title
  <em class='small'>method</em></h2>\n\n???+ source \"title_title <em class='small'>source</em>\"\n\n```python\n\n
  \       def title_title(cls, v, *, values):\n                title = v or Path(values[\"path\"]).stem.replace(\"-\",
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
  v\n                return values.get(\"include\", None)\n```\n\n"
date: 0001-01-01
description: 'None ! ???+ source  ! ! ???+ source  ! ???+ source  ! ???+ source  !
  ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ! ! ! ! ! ! ! ! ???+
  source  ! '
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Post_Model.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None ! ???+ source  ! ! ???+ source  !
    ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
    source  ! ! ! ! ! ! ! ! ! ???+ source  ! \" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Post_Model.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None ! ???+ source  ! ! ???+ source  !
    ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
    source  ! ! ! ! ! ! ! ! ! ???+ source  ! \" />\n <link href=\"/favicon.ico\" rel=\"icon\"
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Post_Model.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <p>None</p>\n<p>!!
    class <h2 id='Post' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Post
    <em class='small'>class</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">Post <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nc\">Post</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">,</span>
    <span class=\"n\">JupyterMixin</span><span class=\"p\">):</span>\n            <span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"n\">Any</span>
    <span class=\"o\">=</span> <span class=\"n\">Field</span><span class=\"p\">(</span><span
    class=\"kc\">None</span><span class=\"p\">,</span> <span class=\"n\">exclude</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">path</span><span class=\"p\">:</span> <span class=\"n\">Path</span>\n
    \           <span class=\"n\">slug</span><span class=\"p\">:</span> <span class=\"n\">Optional</span><span
    class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">]</span> <span
    class=\"o\">=</span> <span class=\"kc\">None</span>\n            <span class=\"n\">href</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"kc\">None</span>\n            <span class=\"n\">published</span><span
    class=\"p\">:</span> <span class=\"nb\">bool</span> <span class=\"o\">=</span>
    <span class=\"kc\">False</span>\n            <span class=\"n\">description</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"kc\">None</span>\n            <span class=\"n\">content</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n
    \           <span class=\"c1\"># date: Union[datetime.date, str]=None</span>\n
    \           <span class=\"n\">date</span><span class=\"p\">:</span> <span class=\"n\">Optional</span><span
    class=\"p\">[</span><span class=\"n\">Union</span><span class=\"p\">[</span><span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">date</span><span
    class=\"p\">,</span> <span class=\"nb\">str</span><span class=\"p\">]]</span>
    <span class=\"o\">=</span> <span class=\"kc\">None</span>\n            <span class=\"c1\">#
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
    <span class=\"mi\">0</span>\n            <span class=\"n\">profile</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"kc\">None</span>\n            <span class=\"n\">title</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n
    \           <span class=\"n\">model_config</span> <span class=\"o\">=</span> <span
    class=\"n\">ConfigDict</span><span class=\"p\">(</span>\n                <span
    class=\"n\">validate_assignment</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span>\n                <span class=\"n\">arbitrary_types_allowed</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">extra</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;allow&quot;</span><span class=\"p\">,</span>\n            <span
    class=\"p\">)</span>\n            <span class=\"n\">template</span><span class=\"p\">:</span>
    <span class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"nb\">str</span>
    <span class=\"o\">|</span> <span class=\"n\">Dict</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
    class=\"p\">]]</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;post.html&quot;</span>\n
    \           <span class=\"n\">sidebar</span><span class=\"p\">:</span> <span class=\"n\">Optional</span><span
    class=\"p\">[</span><span class=\"n\">Any</span><span class=\"p\">]</span> <span
    class=\"o\">=</span> <span class=\"kc\">None</span>\n\n            <span class=\"k\">def</span>
    <span class=\"nf\">__rich__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Pretty</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"n\">Pretty</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">)</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"nf\">__repr_args__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">:</span> <span
    class=\"s2\">&quot;Post&quot;</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"s2\">&quot;ReprArgs&quot;</span><span class=\"p\">:</span>\n                <span
    class=\"k\">return</span> <span class=\"p\">[</span>\n                    <span
    class=\"p\">(</span><span class=\"n\">key</span><span class=\"p\">,</span> <span
    class=\"n\">value</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">for</span> <span class=\"n\">key</span><span class=\"p\">,</span>
    <span class=\"n\">value</span> <span class=\"ow\">in</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"vm\">__dict__</span><span class=\"o\">.</span><span
    class=\"n\">items</span><span class=\"p\">()</span>\n                    <span
    class=\"k\">if</span> <span class=\"n\">key</span> <span class=\"ow\">in</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">post_model</span><span class=\"o\">.</span><span class=\"n\">repr_include</span>\n
    \               <span class=\"p\">]</span>\n\n            <span class=\"nd\">@property</span>\n
    \           <span class=\"k\">def</span> <span class=\"nf\">key</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">List</span><span
    class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">]:</span>\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">make_hash</span><span class=\"p\">(</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">slug</span><span
    class=\"p\">,</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">href</span><span class=\"p\">,</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">published</span><span class=\"p\">,</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">description</span><span
    class=\"p\">,</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">content</span><span class=\"p\">,</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">date</span><span class=\"p\">,</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">title</span><span class=\"p\">,</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">template</span><span class=\"p\">,</span>\n                <span class=\"p\">)</span>\n\n
    \           <span class=\"nd\">@property</span>\n            <span class=\"k\">def</span>
    <span class=\"nf\">metadata</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">Dict</span><span class=\"p\">:</span>\n
    \               <span class=\"s2\">&quot;for backwards compatability&quot;</span>\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"vm\">__dict__</span>\n\n            <span class=\"k\">def</span>
    <span class=\"nf\">to_dict</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">Dict</span><span class=\"p\">:</span>\n
    \               <span class=\"s2\">&quot;for backwards compatability&quot;</span>\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"vm\">__dict__</span>\n\n            <span class=\"k\">def</span>
    <span class=\"fm\">__getitem__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">item</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span
    class=\"p\">:</span>\n                <span class=\"s2\">&quot;for backwards compatability&quot;</span>\n
    \               <span class=\"k\">return</span> <span class=\"nb\">getattr</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"n\">item</span><span class=\"p\">)</span>\n\n            <span class=\"k\">def</span>
    <span class=\"fm\">__setitem__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">key</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"n\">item</span><span class=\"p\">:</span> <span
    class=\"n\">Any</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
    class=\"s2\">&quot;for backwards compatability&quot;</span>\n                <span
    class=\"nb\">setattr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">key</span><span class=\"p\">,</span> <span
    class=\"n\">item</span><span class=\"p\">)</span>\n\n            <span class=\"k\">def</span>
    <span class=\"nf\">get</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">item</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"n\">default</span><span class=\"p\">:</span>
    <span class=\"n\">Any</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Any</span><span class=\"p\">:</span>\n                <span
    class=\"s2\">&quot;for backwards compatability&quot;</span>\n                <span
    class=\"k\">return</span> <span class=\"nb\">getattr</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"n\">item</span><span
    class=\"p\">,</span> <span class=\"n\">default</span><span class=\"p\">)</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"nf\">keys</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">List</span><span
    class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">]:</span>\n
    \               <span class=\"s2\">&quot;for backwards compatability&quot;</span>\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"vm\">__dict__</span><span class=\"o\">.</span><span
    class=\"n\">keys</span><span class=\"p\">()</span>\n\n            <span class=\"c1\">#
    def json(</span>\n            <span class=\"c1\">#     self: &quot;Post&quot;,</span>\n
    \           <span class=\"c1\">#     include: Iterable = None,</span>\n            <span
    class=\"c1\">#     all: bool = False,</span>\n            <span class=\"c1\">#
    \    **kwargs,</span>\n            <span class=\"c1\"># ) -&gt; str:</span>\n
    \           <span class=\"c1\">#     &quot;&quot;&quot;</span>\n            <span
    class=\"c1\">#     override function to give a default include value that will
    include</span>\n            <span class=\"c1\">#     user configured includes.</span>\n
    \           <span class=\"c1\">#     &quot;&quot;&quot;</span>\n            <span
    class=\"c1\">#     if all:</span>\n            <span class=\"c1\">#         return
    pydantic.create_model(&quot;Post&quot;, **self)(**self).json(</span>\n            <span
    class=\"c1\">#             **kwargs,</span>\n            <span class=\"c1\">#
    \        )</span>\n            <span class=\"c1\">#     if include:</span>\n            <span
    class=\"c1\">#         return pydantic.create_model(&quot;Post&quot;, **self)(**self).json(</span>\n
    \           <span class=\"c1\">#             include=include,</span>\n            <span
    class=\"c1\">#             **kwargs,</span>\n            <span class=\"c1\">#
    \        )</span>\n            <span class=\"c1\">#     return pydantic.create_model(&quot;Post&quot;,
    **self)(**self).json(</span>\n            <span class=\"c1\">#         include={i:
    True for i in self.markata.config.post_model.include},</span>\n            <span
    class=\"c1\">#         **kwargs,</span>\n            <span class=\"c1\">#     )</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"nf\">yaml</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span
    class=\"p\">:</span>\n<span class=\"w\">                </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">                dump model to yaml</span>\n<span class=\"sd\">                &quot;&quot;&quot;</span>\n
    \               <span class=\"kn\">import</span> <span class=\"nn\">yaml</span>\n\n
    \               <span class=\"k\">return</span> <span class=\"n\">yaml</span><span
    class=\"o\">.</span><span class=\"n\">dump</span><span class=\"p\">(</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">dict</span><span class=\"p\">(</span>\n                        <span
    class=\"n\">include</span><span class=\"o\">=</span><span class=\"p\">{</span><span
    class=\"n\">i</span><span class=\"p\">:</span> <span class=\"kc\">True</span>
    <span class=\"k\">for</span> <span class=\"n\">i</span> <span class=\"ow\">in</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">post_model</span><span class=\"o\">.</span><span class=\"n\">include</span><span
    class=\"p\">}</span>\n                    <span class=\"p\">),</span>\n                    <span
    class=\"n\">Dumper</span><span class=\"o\">=</span><span class=\"n\">yaml</span><span
    class=\"o\">.</span><span class=\"n\">CDumper</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">)</span>\n\n            <span class=\"k\">def</span>
    <span class=\"nf\">markdown</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n<span
    class=\"w\">                </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">                dump model to markdown</span>\n<span class=\"sd\">
    \               &quot;&quot;&quot;</span>\n\n                <span class=\"kn\">import</span>
    <span class=\"nn\">yaml</span>\n\n                <span class=\"n\">frontmatter</span>
    <span class=\"o\">=</span> <span class=\"n\">yaml</span><span class=\"o\">.</span><span
    class=\"n\">dump</span><span class=\"p\">(</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">dict</span><span class=\"p\">(</span>\n
    \                       <span class=\"n\">include</span><span class=\"o\">=</span><span
    class=\"p\">{</span>\n                            <span class=\"n\">i</span><span
    class=\"p\">:</span> <span class=\"kc\">True</span>\n                            <span
    class=\"k\">for</span> <span class=\"n\">i</span> <span class=\"ow\">in</span>
    <span class=\"p\">[</span>\n                                <span class=\"n\">_i</span>\n
    \                               <span class=\"k\">for</span> <span class=\"n\">_i</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">post_model</span><span class=\"o\">.</span><span
    class=\"n\">include</span>\n                                <span class=\"k\">if</span>
    <span class=\"n\">_i</span> <span class=\"o\">!=</span> <span class=\"s2\">&quot;content&quot;</span>\n
    \                           <span class=\"p\">]</span>\n                        <span
    class=\"p\">}</span>\n                    <span class=\"p\">),</span>\n                    <span
    class=\"n\">Dumper</span><span class=\"o\">=</span><span class=\"n\">yaml</span><span
    class=\"o\">.</span><span class=\"n\">CDumper</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">)</span>\n                <span class=\"n\">post</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;---</span><span class=\"se\">\\n</span><span
    class=\"s2\">&quot;</span>\n                <span class=\"n\">post</span> <span
    class=\"o\">+=</span> <span class=\"n\">frontmatter</span>\n                <span
    class=\"n\">post</span> <span class=\"o\">+=</span> <span class=\"s2\">&quot;---</span><span
    class=\"se\">\\n\\n</span><span class=\"s2\">&quot;</span>\n\n                <span
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
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"s2\">&quot;Post&quot;</span><span
    class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">path</span><span class=\"p\">,</span> <span
    class=\"n\">Path</span><span class=\"p\">):</span>\n                    <span
    class=\"k\">if</span> <span class=\"n\">path</span><span class=\"o\">.</span><span
    class=\"n\">suffix</span> <span class=\"ow\">in</span> <span class=\"p\">[</span><span
    class=\"s2\">&quot;.md&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;.markdown&quot;</span><span
    class=\"p\">]:</span>\n                        <span class=\"k\">return</span>
    <span class=\"bp\">cls</span><span class=\"o\">.</span><span class=\"n\">parse_markdown</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"p\">,</span> <span class=\"n\">path</span><span
    class=\"o\">=</span><span class=\"n\">path</span><span class=\"p\">,</span> <span
    class=\"o\">**</span><span class=\"n\">kwargs</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">elif</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">path</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">):</span>\n                    <span
    class=\"k\">if</span> <span class=\"n\">path</span><span class=\"o\">.</span><span
    class=\"n\">endswith</span><span class=\"p\">(</span><span class=\"s2\">&quot;.md&quot;</span><span
    class=\"p\">)</span> <span class=\"ow\">or</span> <span class=\"n\">path</span><span
    class=\"o\">.</span><span class=\"n\">endswith</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;.markdown&quot;</span><span class=\"p\">):</span>\n                        <span
    class=\"k\">return</span> <span class=\"bp\">cls</span><span class=\"o\">.</span><span
    class=\"n\">parse_markdown</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">path</span><span class=\"o\">=</span><span class=\"n\">path</span><span
    class=\"p\">,</span> <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"nb\">super</span><span
    class=\"p\">(</span><span class=\"n\">Post</span><span class=\"p\">,</span> <span
    class=\"bp\">cls</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">parse_file</span><span class=\"p\">(</span><span class=\"n\">path</span><span
    class=\"p\">,</span> <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">)</span>\n\n            <span class=\"nd\">@classmethod</span>\n            <span
    class=\"k\">def</span> <span class=\"nf\">parse_markdown</span><span class=\"p\">(</span><span
    class=\"bp\">cls</span><span class=\"p\">,</span> <span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"n\">path</span><span class=\"p\">:</span> <span
    class=\"n\">Union</span><span class=\"p\">[</span><span class=\"n\">Path</span><span
    class=\"p\">,</span> <span class=\"nb\">str</span><span class=\"p\">],</span>
    <span class=\"o\">**</span><span class=\"n\">kwargs</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"s2\">&quot;Post&quot;</span><span
    class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">path</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">path</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">path</span><span class=\"p\">)</span>\n
    \               <span class=\"n\">text</span> <span class=\"o\">=</span> <span
    class=\"n\">path</span><span class=\"o\">.</span><span class=\"n\">read_text</span><span
    class=\"p\">()</span>\n                <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">_</span><span class=\"p\">,</span> <span
    class=\"n\">fm</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"n\">content</span> <span class=\"o\">=</span> <span class=\"n\">text</span><span
    class=\"o\">.</span><span class=\"n\">split</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;---</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span><span
    class=\"p\">)</span>\n                    <span class=\"n\">content</span> <span
    class=\"o\">=</span> <span class=\"s2\">&quot;---</span><span class=\"se\">\\n</span><span
    class=\"s2\">&quot;</span><span class=\"o\">.</span><span class=\"n\">join</span><span
    class=\"p\">(</span><span class=\"n\">content</span><span class=\"p\">)</span>\n
    \                   <span class=\"k\">try</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">fm</span> <span class=\"o\">=</span> <span class=\"n\">yaml</span><span
    class=\"o\">.</span><span class=\"n\">load</span><span class=\"p\">(</span><span
    class=\"n\">fm</span><span class=\"p\">,</span> <span class=\"n\">Loader</span><span
    class=\"o\">=</span><span class=\"n\">yaml</span><span class=\"o\">.</span><span
    class=\"n\">CBaseLoader</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">except</span> <span class=\"n\">yaml</span><span class=\"o\">.</span><span
    class=\"n\">YAMLError</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">fm</span> <span class=\"o\">=</span> <span class=\"p\">{}</span>\n
    \               <span class=\"k\">except</span> <span class=\"ne\">ValueError</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">fm</span> <span class=\"o\">=</span>
    <span class=\"p\">{}</span>\n                    <span class=\"n\">content</span>
    <span class=\"o\">=</span> <span class=\"n\">text</span>\n                <span
    class=\"k\">if</span> <span class=\"n\">fm</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span> <span class=\"ow\">or</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">fm</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">fm</span> <span class=\"o\">=</span> <span class=\"p\">{}</span>\n\n
    \               <span class=\"n\">post_args</span> <span class=\"o\">=</span>
    <span class=\"p\">{</span>\n                    <span class=\"s2\">&quot;markata&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">markata</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;path&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">path</span><span class=\"p\">,</span>\n                    <span
    class=\"s2\">&quot;content&quot;</span><span class=\"p\">:</span> <span class=\"n\">content</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;raw&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">text</span><span class=\"p\">,</span>\n
    \                   <span class=\"o\">**</span><span class=\"n\">fm</span><span
    class=\"p\">,</span>\n                <span class=\"p\">}</span>\n\n                <span
    class=\"k\">return</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">Post</span><span class=\"o\">.</span><span class=\"n\">parse_obj</span><span
    class=\"p\">(</span><span class=\"n\">post_args</span><span class=\"p\">)</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"nf\">dumps</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">):</span>\n<span
    class=\"w\">                </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">                dumps raw article back out</span>\n<span class=\"sd\">
    \               &quot;&quot;&quot;</span>\n                <span class=\"k\">return</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;---</span><span class=\"se\">\\n</span><span
    class=\"si\">{</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">yaml</span><span class=\"p\">()</span><span class=\"si\">}</span><span
    class=\"se\">\\n\\n</span><span class=\"s2\">---</span><span class=\"se\">\\n\\n</span><span
    class=\"si\">{</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">content</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span>\n\n
    \           <span class=\"nd\">@pydantic</span><span class=\"o\">.</span><span
    class=\"n\">validator</span><span class=\"p\">(</span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">pre</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span> <span class=\"n\">always</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">def</span> <span class=\"nf\">default_slug</span><span
    class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"p\">,</span> <span class=\"n\">values</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">v</span> <span
    class=\"ow\">or</span> <span class=\"n\">slugify</span><span class=\"p\">(</span><span
    class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;path&quot;</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">stem</span><span class=\"p\">))</span>\n\n
    \           <span class=\"nd\">@pydantic</span><span class=\"o\">.</span><span
    class=\"n\">validator</span><span class=\"p\">(</span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">pre</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span> <span class=\"n\">always</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">def</span> <span class=\"nf\">index_slug_is_empty</span><span
    class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"p\">,</span> <span class=\"n\">values</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"o\">==</span>
    <span class=\"s2\">&quot;index&quot;</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"s2\">&quot;&quot;</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">v</span>\n\n            <span class=\"nd\">@pydantic</span><span
    class=\"o\">.</span><span class=\"n\">validator</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;slug&quot;</span><span class=\"p\">,</span> <span class=\"n\">pre</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
    class=\"n\">always</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n            <span class=\"k\">def</span> <span class=\"nf\">no_double_slash_in_slug</span><span
    class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"p\">,</span> <span class=\"n\">values</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">v</span>\n                <span class=\"k\">return</span>
    <span class=\"n\">v</span><span class=\"o\">.</span><span class=\"n\">replace</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;//&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;/&quot;</span><span class=\"p\">)</span>\n\n            <span
    class=\"nd\">@pydantic</span><span class=\"o\">.</span><span class=\"n\">validator</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;href&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">pre</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span> <span class=\"n\">always</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n            <span class=\"k\">def</span>
    <span class=\"nf\">default_href</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
    class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">v</span>\n                <span class=\"k\">return</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;/</span><span class=\"si\">{</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s1\">&#39;slug&#39;</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">strip</span><span class=\"p\">(</span><span
    class=\"s1\">&#39;/&#39;</span><span class=\"p\">)</span><span class=\"si\">}</span><span
    class=\"s2\">/&quot;</span><span class=\"o\">.</span><span class=\"n\">replace</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;//&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;/&quot;</span><span class=\"p\">)</span>\n\n            <span
    class=\"nd\">@pydantic</span><span class=\"o\">.</span><span class=\"n\">validator</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;title&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">pre</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span> <span class=\"n\">always</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n            <span class=\"k\">def</span>
    <span class=\"nf\">title_title</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
    class=\"p\">):</span>\n                <span class=\"n\">title</span> <span class=\"o\">=</span>
    <span class=\"n\">v</span> <span class=\"ow\">or</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;path&quot;</span><span class=\"p\">])</span><span class=\"o\">.</span><span
    class=\"n\">stem</span><span class=\"o\">.</span><span class=\"n\">replace</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;-&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot; &quot;</span><span class=\"p\">)</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">title</span><span class=\"o\">.</span><span
    class=\"n\">title</span><span class=\"p\">()</span>\n\n            <span class=\"nd\">@pydantic</span><span
    class=\"o\">.</span><span class=\"n\">validator</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;date_time&quot;</span><span class=\"p\">,</span> <span class=\"n\">pre</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
    class=\"n\">always</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n            <span class=\"k\">def</span> <span class=\"nf\">dateparser_datetime</span><span
    class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"p\">,</span> <span class=\"n\">values</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">d</span> <span class=\"o\">=</span> <span class=\"n\">dateparser</span><span
    class=\"o\">.</span><span class=\"n\">parse</span><span class=\"p\">(</span><span
    class=\"n\">v</span><span class=\"p\">)</span>\n                    <span class=\"k\">if</span>
    <span class=\"n\">d</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                        <span class=\"k\">raise</span> <span
    class=\"ne\">ValueError</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s1\">&#39;&quot;</span><span class=\"si\">{</span><span class=\"n\">v</span><span
    class=\"si\">}</span><span class=\"s1\">&quot; is not a valid date&#39;</span><span
    class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span>\n\n
    \           <span class=\"nd\">@pydantic</span><span class=\"o\">.</span><span
    class=\"n\">validator</span><span class=\"p\">(</span><span class=\"s2\">&quot;date_time&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">pre</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span> <span class=\"n\">always</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">def</span> <span class=\"nf\">date_is_datetime</span><span
    class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"p\">,</span> <span class=\"n\">values</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span> <span class=\"ow\">and</span> <span class=\"s2\">&quot;date&quot;</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">values</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;markata&quot;</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">log</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;</span><span class=\"si\">{</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s1\">&#39;path&#39;</span><span class=\"p\">]</span><span
    class=\"si\">}</span><span class=\"s2\"> has no date&quot;</span><span class=\"p\">)</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">now</span><span class=\"p\">()</span>\n                <span class=\"k\">if</span>
    <span class=\"n\">v</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span>
    <span class=\"ow\">and</span> <span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;date&quot;</span><span class=\"p\">]</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">&quot;markata&quot;</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">log</span><span class=\"p\">(</span><span
    class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">values</span><span class=\"p\">[</span><span class=\"s1\">&#39;path&#39;</span><span
    class=\"p\">]</span><span class=\"si\">}</span><span class=\"s2\"> has no date&quot;</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">now</span><span class=\"p\">()</span>\n
    \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">v</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;date&quot;</span><span class=\"p\">],</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"p\">):</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;date&quot;</span><span class=\"p\">]</span>\n
    \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">date</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">combine</span><span class=\"p\">(</span><span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">time</span><span class=\"o\">.</span><span
    class=\"n\">min</span><span class=\"p\">)</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;date&quot;</span><span class=\"p\">],</span>
    <span class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">date</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">combine</span><span class=\"p\">(</span><span
    class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">&quot;date&quot;</span><span
    class=\"p\">],</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">time</span><span class=\"o\">.</span><span class=\"n\">min</span><span
    class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span>\n\n
    \           <span class=\"nd\">@pydantic</span><span class=\"o\">.</span><span
    class=\"n\">validator</span><span class=\"p\">(</span><span class=\"s2\">&quot;date_time&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">pre</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span> <span class=\"n\">always</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">def</span> <span class=\"nf\">mindate_time</span><span
    class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"p\">,</span> <span class=\"n\">values</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span> <span class=\"ow\">and</span> <span class=\"s2\">&quot;date&quot;</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">values</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;markata&quot;</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">log</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;</span><span class=\"si\">{</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s1\">&#39;path&#39;</span><span class=\"p\">]</span><span
    class=\"si\">}</span><span class=\"s2\"> has no date&quot;</span><span class=\"p\">)</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">min</span>\n                <span class=\"k\">if</span> <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;date&quot;</span><span class=\"p\">]</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;markata&quot;</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
    class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s1\">&#39;path&#39;</span><span class=\"p\">]</span><span class=\"si\">}</span><span
    class=\"s2\"> has no date&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">min</span>\n
    \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">v</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;date&quot;</span><span class=\"p\">],</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"p\">):</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;date&quot;</span><span class=\"p\">]</span>\n
    \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">date</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">combine</span><span class=\"p\">(</span><span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">time</span><span class=\"o\">.</span><span
    class=\"n\">min</span><span class=\"p\">)</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;date&quot;</span><span class=\"p\">],</span>
    <span class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">date</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">combine</span><span class=\"p\">(</span><span
    class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">&quot;date&quot;</span><span
    class=\"p\">],</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">time</span><span class=\"o\">.</span><span class=\"n\">min</span><span
    class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span>\n\n
    \           <span class=\"nd\">@pydantic</span><span class=\"o\">.</span><span
    class=\"n\">validator</span><span class=\"p\">(</span><span class=\"s2\">&quot;date&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">pre</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span> <span class=\"n\">always</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">def</span> <span class=\"nf\">dateparser_date</span><span
    class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"p\">,</span> <span class=\"n\">values</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">date</span><span class=\"o\">.</span><span class=\"n\">min</span>\n
    \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">d</span> <span class=\"o\">=</span> <span class=\"bp\">cls</span><span
    class=\"o\">.</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">precache</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">if</span> <span class=\"n\">d</span> <span class=\"ow\">is</span>
    <span class=\"ow\">not</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                       <span class=\"k\">return</span> <span class=\"n\">d</span>\n
    \                   <span class=\"n\">d</span> <span class=\"o\">=</span> <span
    class=\"n\">dateparser</span><span class=\"o\">.</span><span class=\"n\">parse</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">if</span> <span class=\"n\">d</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                        <span
    class=\"k\">raise</span> <span class=\"ne\">ValueError</span><span class=\"p\">(</span><span
    class=\"sa\">f</span><span class=\"s1\">&#39;&quot;</span><span class=\"si\">{</span><span
    class=\"n\">v</span><span class=\"si\">}</span><span class=\"s1\">&quot; is not
    a valid date&#39;</span><span class=\"p\">)</span>\n                    <span
    class=\"n\">d</span> <span class=\"o\">=</span> <span class=\"n\">d</span><span
    class=\"o\">.</span><span class=\"n\">date</span><span class=\"p\">()</span>\n
    \                   <span class=\"k\">with</span> <span class=\"bp\">cls</span><span
    class=\"o\">.</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">cache</span> <span class=\"k\">as</span> <span class=\"n\">cache</span><span
    class=\"p\">:</span>\n                        <span class=\"n\">cache</span><span
    class=\"o\">.</span><span class=\"n\">add</span><span class=\"p\">(</span><span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"n\">d</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">d</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='PostModelConfig' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>PostModelConfig <em class='small'>class</em></h2>\nConfiguration for the
    Post model</p>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">PostModelConfig <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nc\">PostModelConfig</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"s2\">&quot;Configuration for the Post model&quot;</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"fm\">__init__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"o\">**</span><span class=\"n\">data</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
    class=\"w\">                </span><span class=\"sd\">&quot;&quot;&quot;</span>\n\n<span
    class=\"sd\">                include: post attributes to include by default in
    Post</span>\n<span class=\"sd\">                model serialization.</span>\n<span
    class=\"sd\">                repr_include: post attributes to include by default
    in Post</span>\n<span class=\"sd\">                repr.  If `repr_include` is
    None, it will default to</span>\n<span class=\"sd\">                `include`,
    but it is likely that you want less in the repr</span>\n<span class=\"sd\">                than
    serialized output.</span>\n\n<span class=\"sd\">                example:</span>\n\n<span
    class=\"sd\">                ``` toml title=&#39;markata.toml&#39;</span>\n<span
    class=\"sd\">                [markata.post_model]</span>\n<span class=\"sd\">
    \               include = [&#39;date&#39;, &#39;description&#39;, &#39;published&#39;,</span>\n<span
    class=\"sd\">                    &#39;slug&#39;, &#39;title&#39;, &#39;content&#39;,
    &#39;html&#39;]</span>\n<span class=\"sd\">                repr_include = [&#39;date&#39;,
    &#39;description&#39;, &#39;published&#39;, &#39;slug&#39;, &#39;title&#39;]</span>\n<span
    class=\"sd\">                ```</span>\n<span class=\"sd\">                &quot;&quot;&quot;</span>\n
    \               <span class=\"nb\">super</span><span class=\"p\">()</span><span
    class=\"o\">.</span><span class=\"fm\">__init__</span><span class=\"p\">(</span><span
    class=\"o\">**</span><span class=\"n\">data</span><span class=\"p\">)</span>\n\n
    \           <span class=\"n\">include</span><span class=\"p\">:</span> <span class=\"n\">List</span><span
    class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">]</span> <span
    class=\"o\">=</span> <span class=\"p\">[</span>\n                <span class=\"s2\">&quot;date&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;description&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;published&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;title&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;content&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;html&quot;</span><span
    class=\"p\">,</span>\n            <span class=\"p\">]</span>\n            <span
    class=\"n\">repr_include</span><span class=\"p\">:</span> <span class=\"n\">Optional</span><span
    class=\"p\">[</span><span class=\"n\">List</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]]</span> <span class=\"o\">=</span>
    <span class=\"p\">[</span>\n                <span class=\"s2\">&quot;date&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;description&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;published&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;title&quot;</span><span
    class=\"p\">,</span>\n            <span class=\"p\">]</span>\n            <span
    class=\"n\">export_include</span><span class=\"p\">:</span> <span class=\"n\">Optional</span><span
    class=\"p\">[</span><span class=\"n\">List</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]]</span> <span class=\"o\">=</span>
    <span class=\"p\">[</span>\n                <span class=\"s2\">&quot;date&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;description&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;published&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;title&quot;</span><span
    class=\"p\">,</span>\n            <span class=\"p\">]</span>\n\n            <span
    class=\"nd\">@pydantic</span><span class=\"o\">.</span><span class=\"n\">validator</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;repr_include&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">pre</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span> <span class=\"n\">always</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n            <span class=\"k\">def</span>
    <span class=\"nf\">repr_include_validator</span><span class=\"p\">(</span><span
    class=\"bp\">cls</span><span class=\"p\">,</span> <span class=\"n\">v</span><span
    class=\"p\">,</span> <span class=\"o\">*</span><span class=\"p\">,</span> <span
    class=\"n\">values</span><span class=\"p\">):</span>\n                <span class=\"k\">if</span>
    <span class=\"n\">v</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">v</span>\n                <span class=\"k\">return</span>
    <span class=\"n\">values</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;include&quot;</span><span class=\"p\">,</span>
    <span class=\"kc\">None</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    <span class=\"nc\">Config</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">post_model</span><span class=\"p\">:</span> <span
    class=\"n\">PostModelConfig</span> <span class=\"o\">=</span> <span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">Field</span><span class=\"p\">(</span><span
    class=\"n\">default_factory</span><span class=\"o\">=</span><span class=\"n\">PostModelConfig</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='post_model'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>post_model <em class='small'>function</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">post_model
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
    <span class=\"nf\">post_model</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">post_models</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">Post</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    class <h2 id='PostFactory' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>PostFactory <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">PostFactory
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
    <span class=\"nc\">PostFactory</span><span class=\"p\">(</span><span class=\"n\">ModelFactory</span><span
    class=\"p\">):</span>\n            <span class=\"n\">__model__</span> <span class=\"o\">=</span>
    <span class=\"n\">Post</span>\n            <span class=\"n\">__model__</span>
    <span class=\"o\">=</span> <span class=\"n\">Post</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"n\">Pretty</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='__repr_args__' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'><strong>repr_args</strong> <em class='small'>method</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\"><strong>repr_args</strong>
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
    <span class=\"nf\">__repr_args__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"s2\">&quot;ReprArgs&quot;</span><span
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
    class=\"n\">repr_include</span>\n                <span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='key' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>key
    <em class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">key <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">key</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">List</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]:</span>\n                <span class=\"k\">return</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">make_hash</span><span class=\"p\">(</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">slug</span><span class=\"p\">,</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">href</span><span class=\"p\">,</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">published</span><span class=\"p\">,</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">description</span><span
    class=\"p\">,</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">content</span><span class=\"p\">,</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">date</span><span class=\"p\">,</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">title</span><span class=\"p\">,</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">template</span><span class=\"p\">,</span>\n                <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='metadata' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>metadata <em class='small'>method</em></h2>\nfor backwards compatability</p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">metadata
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
    <span class=\"nf\">metadata</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">Dict</span><span class=\"p\">:</span>\n
    \               <span class=\"s2\">&quot;for backwards compatability&quot;</span>\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"vm\">__dict__</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='to_dict' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>to_dict <em class='small'>method</em></h2>\nfor backwards compatability</p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">to_dict
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
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">Dict</span><span class=\"p\">:</span>\n
    \               <span class=\"s2\">&quot;for backwards compatability&quot;</span>\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"vm\">__dict__</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='__getitem__' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'><strong>getitem</strong> <em class='small'>method</em></h2>\nfor backwards
    compatability</p>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\"><strong>getitem</strong> <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"fm\">__getitem__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">item</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span
    class=\"p\">:</span>\n                <span class=\"s2\">&quot;for backwards compatability&quot;</span>\n
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
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">key</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"n\">item</span><span class=\"p\">:</span> <span
    class=\"n\">Any</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
    class=\"s2\">&quot;for backwards compatability&quot;</span>\n                <span
    class=\"nb\">setattr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
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
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">item</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"n\">default</span><span class=\"p\">:</span>
    <span class=\"n\">Any</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Any</span><span class=\"p\">:</span>\n                <span
    class=\"s2\">&quot;for backwards compatability&quot;</span>\n                <span
    class=\"k\">return</span> <span class=\"nb\">getattr</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"n\">item</span><span
    class=\"p\">,</span> <span class=\"n\">default</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">List</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]:</span>\n                <span class=\"s2\">&quot;for
    backwards compatability&quot;</span>\n                <span class=\"k\">return</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"vm\">__dict__</span><span
    class=\"o\">.</span><span class=\"n\">keys</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='yaml' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>yaml
    <em class='small'>method</em></h2>\ndump model to yaml</p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">yaml
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
    <span class=\"nf\">yaml</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n<span
    class=\"w\">                </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">                dump model to yaml</span>\n<span class=\"sd\">                &quot;&quot;&quot;</span>\n
    \               <span class=\"kn\">import</span> <span class=\"nn\">yaml</span>\n\n
    \               <span class=\"k\">return</span> <span class=\"n\">yaml</span><span
    class=\"o\">.</span><span class=\"n\">dump</span><span class=\"p\">(</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">dict</span><span class=\"p\">(</span>\n                        <span
    class=\"n\">include</span><span class=\"o\">=</span><span class=\"p\">{</span><span
    class=\"n\">i</span><span class=\"p\">:</span> <span class=\"kc\">True</span>
    <span class=\"k\">for</span> <span class=\"n\">i</span> <span class=\"ow\">in</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">post_model</span><span class=\"o\">.</span><span class=\"n\">include</span><span
    class=\"p\">}</span>\n                    <span class=\"p\">),</span>\n                    <span
    class=\"n\">Dumper</span><span class=\"o\">=</span><span class=\"n\">yaml</span><span
    class=\"o\">.</span><span class=\"n\">CDumper</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! method
    <h2 id='markdown' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>markdown
    <em class='small'>method</em></h2>\ndump model to markdown</p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">markdown
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
    <span class=\"nf\">markdown</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n<span
    class=\"w\">                </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">                dump model to markdown</span>\n<span class=\"sd\">
    \               &quot;&quot;&quot;</span>\n\n                <span class=\"kn\">import</span>
    <span class=\"nn\">yaml</span>\n\n                <span class=\"n\">frontmatter</span>
    <span class=\"o\">=</span> <span class=\"n\">yaml</span><span class=\"o\">.</span><span
    class=\"n\">dump</span><span class=\"p\">(</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">dict</span><span class=\"p\">(</span>\n
    \                       <span class=\"n\">include</span><span class=\"o\">=</span><span
    class=\"p\">{</span>\n                            <span class=\"n\">i</span><span
    class=\"p\">:</span> <span class=\"kc\">True</span>\n                            <span
    class=\"k\">for</span> <span class=\"n\">i</span> <span class=\"ow\">in</span>
    <span class=\"p\">[</span>\n                                <span class=\"n\">_i</span>\n
    \                               <span class=\"k\">for</span> <span class=\"n\">_i</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">post_model</span><span class=\"o\">.</span><span
    class=\"n\">include</span>\n                                <span class=\"k\">if</span>
    <span class=\"n\">_i</span> <span class=\"o\">!=</span> <span class=\"s2\">&quot;content&quot;</span>\n
    \                           <span class=\"p\">]</span>\n                        <span
    class=\"p\">}</span>\n                    <span class=\"p\">),</span>\n                    <span
    class=\"n\">Dumper</span><span class=\"o\">=</span><span class=\"n\">yaml</span><span
    class=\"o\">.</span><span class=\"n\">CDumper</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">)</span>\n                <span class=\"n\">post</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;---</span><span class=\"se\">\\n</span><span
    class=\"s2\">&quot;</span>\n                <span class=\"n\">post</span> <span
    class=\"o\">+=</span> <span class=\"n\">frontmatter</span>\n                <span
    class=\"n\">post</span> <span class=\"o\">+=</span> <span class=\"s2\">&quot;---</span><span
    class=\"se\">\\n\\n</span><span class=\"s2\">&quot;</span>\n\n                <span
    class=\"k\">if</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">content</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">post</span> <span class=\"o\">+=</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">content</span>\n                <span class=\"k\">return</span>
    <span class=\"n\">post</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='parse_file'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>parse_file <em class='small'>method</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">parse_file
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
    <span class=\"nf\">parse_file</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">path</span><span class=\"p\">:</span> <span class=\"n\">Union</span><span
    class=\"p\">[</span><span class=\"n\">Path</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">],</span> <span class=\"o\">**</span><span
    class=\"n\">kwargs</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">:</span>\n                <span
    class=\"k\">if</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">path</span><span class=\"p\">,</span> <span class=\"n\">Path</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">if</span> <span class=\"n\">path</span><span
    class=\"o\">.</span><span class=\"n\">suffix</span> <span class=\"ow\">in</span>
    <span class=\"p\">[</span><span class=\"s2\">&quot;.md&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;.markdown&quot;</span><span class=\"p\">]:</span>\n                        <span
    class=\"k\">return</span> <span class=\"bp\">cls</span><span class=\"o\">.</span><span
    class=\"n\">parse_markdown</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">path</span><span class=\"o\">=</span><span class=\"n\">path</span><span
    class=\"p\">,</span> <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">)</span>\n                <span class=\"k\">elif</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">path</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">):</span>\n                    <span
    class=\"k\">if</span> <span class=\"n\">path</span><span class=\"o\">.</span><span
    class=\"n\">endswith</span><span class=\"p\">(</span><span class=\"s2\">&quot;.md&quot;</span><span
    class=\"p\">)</span> <span class=\"ow\">or</span> <span class=\"n\">path</span><span
    class=\"o\">.</span><span class=\"n\">endswith</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;.markdown&quot;</span><span class=\"p\">):</span>\n                        <span
    class=\"k\">return</span> <span class=\"bp\">cls</span><span class=\"o\">.</span><span
    class=\"n\">parse_markdown</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">path</span><span class=\"o\">=</span><span class=\"n\">path</span><span
    class=\"p\">,</span> <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"nb\">super</span><span
    class=\"p\">(</span><span class=\"n\">Post</span><span class=\"p\">,</span> <span
    class=\"bp\">cls</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">parse_file</span><span class=\"p\">(</span><span class=\"n\">path</span><span
    class=\"p\">,</span> <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='parse_markdown'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>parse_markdown <em
    class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">parse_markdown <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">parse_markdown</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">path</span><span class=\"p\">:</span> <span class=\"n\">Union</span><span
    class=\"p\">[</span><span class=\"n\">Path</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">],</span> <span class=\"o\">**</span><span
    class=\"n\">kwargs</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">:</span>\n                <span
    class=\"k\">if</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">path</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
    class=\"p\">):</span>\n                    <span class=\"n\">path</span> <span
    class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">path</span><span class=\"p\">)</span>\n                <span class=\"n\">text</span>
    <span class=\"o\">=</span> <span class=\"n\">path</span><span class=\"o\">.</span><span
    class=\"n\">read_text</span><span class=\"p\">()</span>\n                <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                    <span class=\"n\">_</span><span
    class=\"p\">,</span> <span class=\"n\">fm</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"n\">content</span> <span class=\"o\">=</span>
    <span class=\"n\">text</span><span class=\"o\">.</span><span class=\"n\">split</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;---</span><span class=\"se\">\\n</span><span
    class=\"s2\">&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"n\">content</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;---</span><span
    class=\"se\">\\n</span><span class=\"s2\">&quot;</span><span class=\"o\">.</span><span
    class=\"n\">join</span><span class=\"p\">(</span><span class=\"n\">content</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">fm</span> <span class=\"o\">=</span>
    <span class=\"n\">yaml</span><span class=\"o\">.</span><span class=\"n\">load</span><span
    class=\"p\">(</span><span class=\"n\">fm</span><span class=\"p\">,</span> <span
    class=\"n\">Loader</span><span class=\"o\">=</span><span class=\"n\">yaml</span><span
    class=\"o\">.</span><span class=\"n\">CBaseLoader</span><span class=\"p\">)</span>\n
    \                   <span class=\"k\">except</span> <span class=\"n\">yaml</span><span
    class=\"o\">.</span><span class=\"n\">YAMLError</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">fm</span> <span class=\"o\">=</span>
    <span class=\"p\">{}</span>\n                <span class=\"k\">except</span> <span
    class=\"ne\">ValueError</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">fm</span> <span class=\"o\">=</span> <span class=\"p\">{}</span>\n
    \                   <span class=\"n\">content</span> <span class=\"o\">=</span>
    <span class=\"n\">text</span>\n                <span class=\"k\">if</span> <span
    class=\"n\">fm</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span>
    <span class=\"ow\">or</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">fm</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
    class=\"p\">):</span>\n                    <span class=\"n\">fm</span> <span class=\"o\">=</span>
    <span class=\"p\">{}</span>\n\n                <span class=\"n\">post_args</span>
    <span class=\"o\">=</span> <span class=\"p\">{</span>\n                    <span
    class=\"s2\">&quot;markata&quot;</span><span class=\"p\">:</span> <span class=\"n\">markata</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;path&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">path</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;content&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">content</span><span class=\"p\">,</span>\n                    <span
    class=\"s2\">&quot;raw&quot;</span><span class=\"p\">:</span> <span class=\"n\">text</span><span
    class=\"p\">,</span>\n                    <span class=\"o\">**</span><span class=\"n\">fm</span><span
    class=\"p\">,</span>\n                <span class=\"p\">}</span>\n\n                <span
    class=\"k\">return</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">Post</span><span class=\"o\">.</span><span class=\"n\">parse_obj</span><span
    class=\"p\">(</span><span class=\"n\">post_args</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='dumps' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>dumps
    <em class='small'>method</em></h2>\ndumps raw article back out</p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">dumps
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
    <span class=\"nf\">dumps</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">):</span>\n<span class=\"w\">                </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">                dumps raw article back out</span>\n<span class=\"sd\">
    \               &quot;&quot;&quot;</span>\n                <span class=\"k\">return</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;---</span><span class=\"se\">\\n</span><span
    class=\"si\">{</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">yaml</span><span class=\"p\">()</span><span class=\"si\">}</span><span
    class=\"se\">\\n\\n</span><span class=\"s2\">---</span><span class=\"se\">\\n\\n</span><span
    class=\"si\">{</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">content</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='default_slug' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>default_slug <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">default_slug
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
    <span class=\"nf\">default_slug</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
    class=\"p\">):</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span>
    <span class=\"ow\">or</span> <span class=\"n\">slugify</span><span class=\"p\">(</span><span
    class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;path&quot;</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">stem</span><span class=\"p\">))</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='index_slug_is_empty' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>index_slug_is_empty <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">index_slug_is_empty
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
    <span class=\"nf\">index_slug_is_empty</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
    class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span>
    <span class=\"o\">==</span> <span class=\"s2\">&quot;index&quot;</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"s2\">&quot;&quot;</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='no_double_slash_in_slug' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>no_double_slash_in_slug <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">no_double_slash_in_slug
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
    <span class=\"nf\">no_double_slash_in_slug</span><span class=\"p\">(</span><span
    class=\"bp\">cls</span><span class=\"p\">,</span> <span class=\"n\">v</span><span
    class=\"p\">,</span> <span class=\"o\">*</span><span class=\"p\">,</span> <span
    class=\"n\">values</span><span class=\"p\">):</span>\n                <span class=\"k\">if</span>
    <span class=\"n\">v</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">v</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span><span
    class=\"o\">.</span><span class=\"n\">replace</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;//&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;/&quot;</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='default_href'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>default_href <em
    class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">default_href <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">default_href</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
    class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">v</span>\n                <span class=\"k\">return</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;/</span><span class=\"si\">{</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s1\">&#39;slug&#39;</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">strip</span><span class=\"p\">(</span><span
    class=\"s1\">&#39;/&#39;</span><span class=\"p\">)</span><span class=\"si\">}</span><span
    class=\"s2\">/&quot;</span><span class=\"o\">.</span><span class=\"n\">replace</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;//&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;/&quot;</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='title_title' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>title_title <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">title_title
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
    <span class=\"nf\">title_title</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
    class=\"p\">):</span>\n                <span class=\"n\">title</span> <span class=\"o\">=</span>
    <span class=\"n\">v</span> <span class=\"ow\">or</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;path&quot;</span><span class=\"p\">])</span><span class=\"o\">.</span><span
    class=\"n\">stem</span><span class=\"o\">.</span><span class=\"n\">replace</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;-&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot; &quot;</span><span class=\"p\">)</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">title</span><span class=\"o\">.</span><span
    class=\"n\">title</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='dateparser_datetime' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>dateparser_datetime <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">dateparser_datetime
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
    <span class=\"nf\">dateparser_datetime</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
    class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">d</span> <span class=\"o\">=</span> <span class=\"n\">dateparser</span><span
    class=\"o\">.</span><span class=\"n\">parse</span><span class=\"p\">(</span><span
    class=\"n\">v</span><span class=\"p\">)</span>\n                    <span class=\"k\">if</span>
    <span class=\"n\">d</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                        <span class=\"k\">raise</span> <span
    class=\"ne\">ValueError</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s1\">&#39;&quot;</span><span class=\"si\">{</span><span class=\"n\">v</span><span
    class=\"si\">}</span><span class=\"s1\">&quot; is not a valid date&#39;</span><span
    class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='date_is_datetime' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>date_is_datetime <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">date_is_datetime
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
    <span class=\"nf\">date_is_datetime</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
    class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span> <span class=\"ow\">and</span>
    <span class=\"s2\">&quot;date&quot;</span> <span class=\"ow\">not</span> <span
    class=\"ow\">in</span> <span class=\"n\">values</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;markata&quot;</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
    class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s1\">&#39;path&#39;</span><span class=\"p\">]</span><span class=\"si\">}</span><span
    class=\"s2\"> has no date&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">now</span><span
    class=\"p\">()</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span> <span class=\"ow\">and</span>
    <span class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">&quot;date&quot;</span><span
    class=\"p\">]</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;markata&quot;</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">log</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;</span><span class=\"si\">{</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s1\">&#39;path&#39;</span><span class=\"p\">]</span><span
    class=\"si\">}</span><span class=\"s2\"> has no date&quot;</span><span class=\"p\">)</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">now</span><span class=\"p\">()</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">v</span><span
    class=\"p\">,</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">datetime</span><span class=\"p\">):</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">v</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;date&quot;</span><span class=\"p\">],</span>
    <span class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">&quot;date&quot;</span><span
    class=\"p\">]</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">date</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">combine</span><span class=\"p\">(</span><span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">time</span><span class=\"o\">.</span><span
    class=\"n\">min</span><span class=\"p\">)</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;date&quot;</span><span class=\"p\">],</span>
    <span class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">date</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">combine</span><span class=\"p\">(</span><span
    class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">&quot;date&quot;</span><span
    class=\"p\">],</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">time</span><span class=\"o\">.</span><span class=\"n\">min</span><span
    class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='mindate_time' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>mindate_time <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">mindate_time
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
    <span class=\"nf\">mindate_time</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
    class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span> <span class=\"ow\">and</span>
    <span class=\"s2\">&quot;date&quot;</span> <span class=\"ow\">not</span> <span
    class=\"ow\">in</span> <span class=\"n\">values</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;markata&quot;</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
    class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s1\">&#39;path&#39;</span><span class=\"p\">]</span><span class=\"si\">}</span><span
    class=\"s2\"> has no date&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">min</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;date&quot;</span><span class=\"p\">]</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;markata&quot;</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
    class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s1\">&#39;path&#39;</span><span class=\"p\">]</span><span class=\"si\">}</span><span
    class=\"s2\"> has no date&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">min</span>\n
    \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">v</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;date&quot;</span><span class=\"p\">],</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"p\">):</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;date&quot;</span><span class=\"p\">]</span>\n
    \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">date</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">combine</span><span class=\"p\">(</span><span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">time</span><span class=\"o\">.</span><span
    class=\"n\">min</span><span class=\"p\">)</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;date&quot;</span><span class=\"p\">],</span>
    <span class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">date</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">combine</span><span class=\"p\">(</span><span
    class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">&quot;date&quot;</span><span
    class=\"p\">],</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">time</span><span class=\"o\">.</span><span class=\"n\">min</span><span
    class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='dateparser_date' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>dateparser_date <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">dateparser_date
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
    <span class=\"nf\">dateparser_date</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
    class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">date</span><span class=\"o\">.</span><span
    class=\"n\">min</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">d</span> <span class=\"o\">=</span> <span class=\"bp\">cls</span><span
    class=\"o\">.</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">precache</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">if</span> <span class=\"n\">d</span> <span class=\"ow\">is</span>
    <span class=\"ow\">not</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                       <span class=\"k\">return</span> <span class=\"n\">d</span>\n
    \                   <span class=\"n\">d</span> <span class=\"o\">=</span> <span
    class=\"n\">dateparser</span><span class=\"o\">.</span><span class=\"n\">parse</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">if</span> <span class=\"n\">d</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                        <span
    class=\"k\">raise</span> <span class=\"ne\">ValueError</span><span class=\"p\">(</span><span
    class=\"sa\">f</span><span class=\"s1\">&#39;&quot;</span><span class=\"si\">{</span><span
    class=\"n\">v</span><span class=\"si\">}</span><span class=\"s1\">&quot; is not
    a valid date&#39;</span><span class=\"p\">)</span>\n                    <span
    class=\"n\">d</span> <span class=\"o\">=</span> <span class=\"n\">d</span><span
    class=\"o\">.</span><span class=\"n\">date</span><span class=\"p\">()</span>\n
    \                   <span class=\"k\">with</span> <span class=\"bp\">cls</span><span
    class=\"o\">.</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">cache</span> <span class=\"k\">as</span> <span class=\"n\">cache</span><span
    class=\"p\">:</span>\n                        <span class=\"n\">cache</span><span
    class=\"o\">.</span><span class=\"n\">add</span><span class=\"p\">(</span><span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"n\">d</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">d</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='__init__' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'><strong>init</strong> <em class='small'>method</em></h2>\ninclude: post
    attributes to include by default in Post\nmodel serialization.\nrepr_include:
    post attributes to include by default in Post\nrepr.  If <code>repr_include</code>
    is None, it will default to\n<code>include</code>, but it is likely that you want
    less in the repr\nthan serialized output.</p>\n<pre><code>example:\n\n``` toml
    title='markata.toml'\n[markata.post_model]\ninclude = ['date', 'description',
    'published',\n    'slug', 'title', 'content', 'html']\nrepr_include = ['date',
    'description', 'published', 'slug', 'title']\n```\n</code></pre>\n<div class=\"admonition
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
    class=\"p\">,</span> <span class=\"o\">**</span><span class=\"n\">data</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n<span class=\"w\">                </span><span class=\"sd\">&quot;&quot;&quot;</span>\n\n<span
    class=\"sd\">                include: post attributes to include by default in
    Post</span>\n<span class=\"sd\">                model serialization.</span>\n<span
    class=\"sd\">                repr_include: post attributes to include by default
    in Post</span>\n<span class=\"sd\">                repr.  If `repr_include` is
    None, it will default to</span>\n<span class=\"sd\">                `include`,
    but it is likely that you want less in the repr</span>\n<span class=\"sd\">                than
    serialized output.</span>\n\n<span class=\"sd\">                example:</span>\n\n<span
    class=\"sd\">                ``` toml title=&#39;markata.toml&#39;</span>\n<span
    class=\"sd\">                [markata.post_model]</span>\n<span class=\"sd\">
    \               include = [&#39;date&#39;, &#39;description&#39;, &#39;published&#39;,</span>\n<span
    class=\"sd\">                    &#39;slug&#39;, &#39;title&#39;, &#39;content&#39;,
    &#39;html&#39;]</span>\n<span class=\"sd\">                repr_include = [&#39;date&#39;,
    &#39;description&#39;, &#39;published&#39;, &#39;slug&#39;, &#39;title&#39;]</span>\n<span
    class=\"sd\">                ```</span>\n<span class=\"sd\">                &quot;&quot;&quot;</span>\n
    \               <span class=\"nb\">super</span><span class=\"p\">()</span><span
    class=\"o\">.</span><span class=\"fm\">__init__</span><span class=\"p\">(</span><span
    class=\"o\">**</span><span class=\"n\">data</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='repr_include_validator' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>repr_include_validator <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">repr_include_validator
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
    <span class=\"nf\">repr_include_validator</span><span class=\"p\">(</span><span
    class=\"bp\">cls</span><span class=\"p\">,</span> <span class=\"n\">v</span><span
    class=\"p\">,</span> <span class=\"o\">*</span><span class=\"p\">,</span> <span
    class=\"n\">values</span><span class=\"p\">):</span>\n                <span class=\"k\">if</span>
    <span class=\"n\">v</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">v</span>\n                <span class=\"k\">return</span>
    <span class=\"n\">values</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;include&quot;</span><span class=\"p\">,</span>
    <span class=\"kc\">None</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>    </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Post_Model.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"None ! ???+ source  ! ! ???+ source  ! ???+ source
    \ ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  !
    ! ! ! ! ! ! ! ! ???+ source  ! \" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Post_Model.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None ! ???+ source  ! ! ???+ source  !
    ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
    source  ! ! ! ! ! ! ! ! ! ???+ source  ! \" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<article style=\"text-align: center;\">\n    <style>\n        section
    {\n            font-size: 200%;\n        }\n\n\n        .edit {\n            display:
    none;\n        }\n    </style>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Post_Model.Py\n    </h1>\n</section></article>\n     </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Post_Model.Py\n    </h1>\n</section>    <section class=\"body\">\n        <p>None</p>\n<p>!!
    class <h2 id='Post' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Post
    <em class='small'>class</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">Post <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nc\">Post</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">,</span>
    <span class=\"n\">JupyterMixin</span><span class=\"p\">):</span>\n            <span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"n\">Any</span>
    <span class=\"o\">=</span> <span class=\"n\">Field</span><span class=\"p\">(</span><span
    class=\"kc\">None</span><span class=\"p\">,</span> <span class=\"n\">exclude</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">path</span><span class=\"p\">:</span> <span class=\"n\">Path</span>\n
    \           <span class=\"n\">slug</span><span class=\"p\">:</span> <span class=\"n\">Optional</span><span
    class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">]</span> <span
    class=\"o\">=</span> <span class=\"kc\">None</span>\n            <span class=\"n\">href</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"kc\">None</span>\n            <span class=\"n\">published</span><span
    class=\"p\">:</span> <span class=\"nb\">bool</span> <span class=\"o\">=</span>
    <span class=\"kc\">False</span>\n            <span class=\"n\">description</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"kc\">None</span>\n            <span class=\"n\">content</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n
    \           <span class=\"c1\"># date: Union[datetime.date, str]=None</span>\n
    \           <span class=\"n\">date</span><span class=\"p\">:</span> <span class=\"n\">Optional</span><span
    class=\"p\">[</span><span class=\"n\">Union</span><span class=\"p\">[</span><span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">date</span><span
    class=\"p\">,</span> <span class=\"nb\">str</span><span class=\"p\">]]</span>
    <span class=\"o\">=</span> <span class=\"kc\">None</span>\n            <span class=\"c1\">#
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
    <span class=\"mi\">0</span>\n            <span class=\"n\">profile</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"kc\">None</span>\n            <span class=\"n\">title</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n
    \           <span class=\"n\">model_config</span> <span class=\"o\">=</span> <span
    class=\"n\">ConfigDict</span><span class=\"p\">(</span>\n                <span
    class=\"n\">validate_assignment</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span>\n                <span class=\"n\">arbitrary_types_allowed</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">extra</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;allow&quot;</span><span class=\"p\">,</span>\n            <span
    class=\"p\">)</span>\n            <span class=\"n\">template</span><span class=\"p\">:</span>
    <span class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"nb\">str</span>
    <span class=\"o\">|</span> <span class=\"n\">Dict</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
    class=\"p\">]]</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;post.html&quot;</span>\n
    \           <span class=\"n\">sidebar</span><span class=\"p\">:</span> <span class=\"n\">Optional</span><span
    class=\"p\">[</span><span class=\"n\">Any</span><span class=\"p\">]</span> <span
    class=\"o\">=</span> <span class=\"kc\">None</span>\n\n            <span class=\"k\">def</span>
    <span class=\"nf\">__rich__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Pretty</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"n\">Pretty</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">)</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"nf\">__repr_args__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">:</span> <span
    class=\"s2\">&quot;Post&quot;</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"s2\">&quot;ReprArgs&quot;</span><span class=\"p\">:</span>\n                <span
    class=\"k\">return</span> <span class=\"p\">[</span>\n                    <span
    class=\"p\">(</span><span class=\"n\">key</span><span class=\"p\">,</span> <span
    class=\"n\">value</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">for</span> <span class=\"n\">key</span><span class=\"p\">,</span>
    <span class=\"n\">value</span> <span class=\"ow\">in</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"vm\">__dict__</span><span class=\"o\">.</span><span
    class=\"n\">items</span><span class=\"p\">()</span>\n                    <span
    class=\"k\">if</span> <span class=\"n\">key</span> <span class=\"ow\">in</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">post_model</span><span class=\"o\">.</span><span class=\"n\">repr_include</span>\n
    \               <span class=\"p\">]</span>\n\n            <span class=\"nd\">@property</span>\n
    \           <span class=\"k\">def</span> <span class=\"nf\">key</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">List</span><span
    class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">]:</span>\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">make_hash</span><span class=\"p\">(</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">slug</span><span
    class=\"p\">,</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">href</span><span class=\"p\">,</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">published</span><span class=\"p\">,</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">description</span><span
    class=\"p\">,</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">content</span><span class=\"p\">,</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">date</span><span class=\"p\">,</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">title</span><span class=\"p\">,</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">template</span><span class=\"p\">,</span>\n                <span class=\"p\">)</span>\n\n
    \           <span class=\"nd\">@property</span>\n            <span class=\"k\">def</span>
    <span class=\"nf\">metadata</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">Dict</span><span class=\"p\">:</span>\n
    \               <span class=\"s2\">&quot;for backwards compatability&quot;</span>\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"vm\">__dict__</span>\n\n            <span class=\"k\">def</span>
    <span class=\"nf\">to_dict</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">Dict</span><span class=\"p\">:</span>\n
    \               <span class=\"s2\">&quot;for backwards compatability&quot;</span>\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"vm\">__dict__</span>\n\n            <span class=\"k\">def</span>
    <span class=\"fm\">__getitem__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">item</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span
    class=\"p\">:</span>\n                <span class=\"s2\">&quot;for backwards compatability&quot;</span>\n
    \               <span class=\"k\">return</span> <span class=\"nb\">getattr</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"n\">item</span><span class=\"p\">)</span>\n\n            <span class=\"k\">def</span>
    <span class=\"fm\">__setitem__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">key</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"n\">item</span><span class=\"p\">:</span> <span
    class=\"n\">Any</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
    class=\"s2\">&quot;for backwards compatability&quot;</span>\n                <span
    class=\"nb\">setattr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">key</span><span class=\"p\">,</span> <span
    class=\"n\">item</span><span class=\"p\">)</span>\n\n            <span class=\"k\">def</span>
    <span class=\"nf\">get</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">item</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"n\">default</span><span class=\"p\">:</span>
    <span class=\"n\">Any</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Any</span><span class=\"p\">:</span>\n                <span
    class=\"s2\">&quot;for backwards compatability&quot;</span>\n                <span
    class=\"k\">return</span> <span class=\"nb\">getattr</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"n\">item</span><span
    class=\"p\">,</span> <span class=\"n\">default</span><span class=\"p\">)</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"nf\">keys</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">List</span><span
    class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">]:</span>\n
    \               <span class=\"s2\">&quot;for backwards compatability&quot;</span>\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"vm\">__dict__</span><span class=\"o\">.</span><span
    class=\"n\">keys</span><span class=\"p\">()</span>\n\n            <span class=\"c1\">#
    def json(</span>\n            <span class=\"c1\">#     self: &quot;Post&quot;,</span>\n
    \           <span class=\"c1\">#     include: Iterable = None,</span>\n            <span
    class=\"c1\">#     all: bool = False,</span>\n            <span class=\"c1\">#
    \    **kwargs,</span>\n            <span class=\"c1\"># ) -&gt; str:</span>\n
    \           <span class=\"c1\">#     &quot;&quot;&quot;</span>\n            <span
    class=\"c1\">#     override function to give a default include value that will
    include</span>\n            <span class=\"c1\">#     user configured includes.</span>\n
    \           <span class=\"c1\">#     &quot;&quot;&quot;</span>\n            <span
    class=\"c1\">#     if all:</span>\n            <span class=\"c1\">#         return
    pydantic.create_model(&quot;Post&quot;, **self)(**self).json(</span>\n            <span
    class=\"c1\">#             **kwargs,</span>\n            <span class=\"c1\">#
    \        )</span>\n            <span class=\"c1\">#     if include:</span>\n            <span
    class=\"c1\">#         return pydantic.create_model(&quot;Post&quot;, **self)(**self).json(</span>\n
    \           <span class=\"c1\">#             include=include,</span>\n            <span
    class=\"c1\">#             **kwargs,</span>\n            <span class=\"c1\">#
    \        )</span>\n            <span class=\"c1\">#     return pydantic.create_model(&quot;Post&quot;,
    **self)(**self).json(</span>\n            <span class=\"c1\">#         include={i:
    True for i in self.markata.config.post_model.include},</span>\n            <span
    class=\"c1\">#         **kwargs,</span>\n            <span class=\"c1\">#     )</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"nf\">yaml</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span
    class=\"p\">:</span>\n<span class=\"w\">                </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">                dump model to yaml</span>\n<span class=\"sd\">                &quot;&quot;&quot;</span>\n
    \               <span class=\"kn\">import</span> <span class=\"nn\">yaml</span>\n\n
    \               <span class=\"k\">return</span> <span class=\"n\">yaml</span><span
    class=\"o\">.</span><span class=\"n\">dump</span><span class=\"p\">(</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">dict</span><span class=\"p\">(</span>\n                        <span
    class=\"n\">include</span><span class=\"o\">=</span><span class=\"p\">{</span><span
    class=\"n\">i</span><span class=\"p\">:</span> <span class=\"kc\">True</span>
    <span class=\"k\">for</span> <span class=\"n\">i</span> <span class=\"ow\">in</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">post_model</span><span class=\"o\">.</span><span class=\"n\">include</span><span
    class=\"p\">}</span>\n                    <span class=\"p\">),</span>\n                    <span
    class=\"n\">Dumper</span><span class=\"o\">=</span><span class=\"n\">yaml</span><span
    class=\"o\">.</span><span class=\"n\">CDumper</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">)</span>\n\n            <span class=\"k\">def</span>
    <span class=\"nf\">markdown</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n<span
    class=\"w\">                </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">                dump model to markdown</span>\n<span class=\"sd\">
    \               &quot;&quot;&quot;</span>\n\n                <span class=\"kn\">import</span>
    <span class=\"nn\">yaml</span>\n\n                <span class=\"n\">frontmatter</span>
    <span class=\"o\">=</span> <span class=\"n\">yaml</span><span class=\"o\">.</span><span
    class=\"n\">dump</span><span class=\"p\">(</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">dict</span><span class=\"p\">(</span>\n
    \                       <span class=\"n\">include</span><span class=\"o\">=</span><span
    class=\"p\">{</span>\n                            <span class=\"n\">i</span><span
    class=\"p\">:</span> <span class=\"kc\">True</span>\n                            <span
    class=\"k\">for</span> <span class=\"n\">i</span> <span class=\"ow\">in</span>
    <span class=\"p\">[</span>\n                                <span class=\"n\">_i</span>\n
    \                               <span class=\"k\">for</span> <span class=\"n\">_i</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">post_model</span><span class=\"o\">.</span><span
    class=\"n\">include</span>\n                                <span class=\"k\">if</span>
    <span class=\"n\">_i</span> <span class=\"o\">!=</span> <span class=\"s2\">&quot;content&quot;</span>\n
    \                           <span class=\"p\">]</span>\n                        <span
    class=\"p\">}</span>\n                    <span class=\"p\">),</span>\n                    <span
    class=\"n\">Dumper</span><span class=\"o\">=</span><span class=\"n\">yaml</span><span
    class=\"o\">.</span><span class=\"n\">CDumper</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">)</span>\n                <span class=\"n\">post</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;---</span><span class=\"se\">\\n</span><span
    class=\"s2\">&quot;</span>\n                <span class=\"n\">post</span> <span
    class=\"o\">+=</span> <span class=\"n\">frontmatter</span>\n                <span
    class=\"n\">post</span> <span class=\"o\">+=</span> <span class=\"s2\">&quot;---</span><span
    class=\"se\">\\n\\n</span><span class=\"s2\">&quot;</span>\n\n                <span
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
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"s2\">&quot;Post&quot;</span><span
    class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">path</span><span class=\"p\">,</span> <span
    class=\"n\">Path</span><span class=\"p\">):</span>\n                    <span
    class=\"k\">if</span> <span class=\"n\">path</span><span class=\"o\">.</span><span
    class=\"n\">suffix</span> <span class=\"ow\">in</span> <span class=\"p\">[</span><span
    class=\"s2\">&quot;.md&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;.markdown&quot;</span><span
    class=\"p\">]:</span>\n                        <span class=\"k\">return</span>
    <span class=\"bp\">cls</span><span class=\"o\">.</span><span class=\"n\">parse_markdown</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"p\">,</span> <span class=\"n\">path</span><span
    class=\"o\">=</span><span class=\"n\">path</span><span class=\"p\">,</span> <span
    class=\"o\">**</span><span class=\"n\">kwargs</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">elif</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">path</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">):</span>\n                    <span
    class=\"k\">if</span> <span class=\"n\">path</span><span class=\"o\">.</span><span
    class=\"n\">endswith</span><span class=\"p\">(</span><span class=\"s2\">&quot;.md&quot;</span><span
    class=\"p\">)</span> <span class=\"ow\">or</span> <span class=\"n\">path</span><span
    class=\"o\">.</span><span class=\"n\">endswith</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;.markdown&quot;</span><span class=\"p\">):</span>\n                        <span
    class=\"k\">return</span> <span class=\"bp\">cls</span><span class=\"o\">.</span><span
    class=\"n\">parse_markdown</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">path</span><span class=\"o\">=</span><span class=\"n\">path</span><span
    class=\"p\">,</span> <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"nb\">super</span><span
    class=\"p\">(</span><span class=\"n\">Post</span><span class=\"p\">,</span> <span
    class=\"bp\">cls</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">parse_file</span><span class=\"p\">(</span><span class=\"n\">path</span><span
    class=\"p\">,</span> <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">)</span>\n\n            <span class=\"nd\">@classmethod</span>\n            <span
    class=\"k\">def</span> <span class=\"nf\">parse_markdown</span><span class=\"p\">(</span><span
    class=\"bp\">cls</span><span class=\"p\">,</span> <span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"n\">path</span><span class=\"p\">:</span> <span
    class=\"n\">Union</span><span class=\"p\">[</span><span class=\"n\">Path</span><span
    class=\"p\">,</span> <span class=\"nb\">str</span><span class=\"p\">],</span>
    <span class=\"o\">**</span><span class=\"n\">kwargs</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"s2\">&quot;Post&quot;</span><span
    class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">path</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">path</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">path</span><span class=\"p\">)</span>\n
    \               <span class=\"n\">text</span> <span class=\"o\">=</span> <span
    class=\"n\">path</span><span class=\"o\">.</span><span class=\"n\">read_text</span><span
    class=\"p\">()</span>\n                <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">_</span><span class=\"p\">,</span> <span
    class=\"n\">fm</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"n\">content</span> <span class=\"o\">=</span> <span class=\"n\">text</span><span
    class=\"o\">.</span><span class=\"n\">split</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;---</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span><span
    class=\"p\">)</span>\n                    <span class=\"n\">content</span> <span
    class=\"o\">=</span> <span class=\"s2\">&quot;---</span><span class=\"se\">\\n</span><span
    class=\"s2\">&quot;</span><span class=\"o\">.</span><span class=\"n\">join</span><span
    class=\"p\">(</span><span class=\"n\">content</span><span class=\"p\">)</span>\n
    \                   <span class=\"k\">try</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">fm</span> <span class=\"o\">=</span> <span class=\"n\">yaml</span><span
    class=\"o\">.</span><span class=\"n\">load</span><span class=\"p\">(</span><span
    class=\"n\">fm</span><span class=\"p\">,</span> <span class=\"n\">Loader</span><span
    class=\"o\">=</span><span class=\"n\">yaml</span><span class=\"o\">.</span><span
    class=\"n\">CBaseLoader</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">except</span> <span class=\"n\">yaml</span><span class=\"o\">.</span><span
    class=\"n\">YAMLError</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">fm</span> <span class=\"o\">=</span> <span class=\"p\">{}</span>\n
    \               <span class=\"k\">except</span> <span class=\"ne\">ValueError</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">fm</span> <span class=\"o\">=</span>
    <span class=\"p\">{}</span>\n                    <span class=\"n\">content</span>
    <span class=\"o\">=</span> <span class=\"n\">text</span>\n                <span
    class=\"k\">if</span> <span class=\"n\">fm</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span> <span class=\"ow\">or</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">fm</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">fm</span> <span class=\"o\">=</span> <span class=\"p\">{}</span>\n\n
    \               <span class=\"n\">post_args</span> <span class=\"o\">=</span>
    <span class=\"p\">{</span>\n                    <span class=\"s2\">&quot;markata&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">markata</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;path&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">path</span><span class=\"p\">,</span>\n                    <span
    class=\"s2\">&quot;content&quot;</span><span class=\"p\">:</span> <span class=\"n\">content</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;raw&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">text</span><span class=\"p\">,</span>\n
    \                   <span class=\"o\">**</span><span class=\"n\">fm</span><span
    class=\"p\">,</span>\n                <span class=\"p\">}</span>\n\n                <span
    class=\"k\">return</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">Post</span><span class=\"o\">.</span><span class=\"n\">parse_obj</span><span
    class=\"p\">(</span><span class=\"n\">post_args</span><span class=\"p\">)</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"nf\">dumps</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">):</span>\n<span
    class=\"w\">                </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">                dumps raw article back out</span>\n<span class=\"sd\">
    \               &quot;&quot;&quot;</span>\n                <span class=\"k\">return</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;---</span><span class=\"se\">\\n</span><span
    class=\"si\">{</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">yaml</span><span class=\"p\">()</span><span class=\"si\">}</span><span
    class=\"se\">\\n\\n</span><span class=\"s2\">---</span><span class=\"se\">\\n\\n</span><span
    class=\"si\">{</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">content</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span>\n\n
    \           <span class=\"nd\">@pydantic</span><span class=\"o\">.</span><span
    class=\"n\">validator</span><span class=\"p\">(</span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">pre</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span> <span class=\"n\">always</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">def</span> <span class=\"nf\">default_slug</span><span
    class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"p\">,</span> <span class=\"n\">values</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">v</span> <span
    class=\"ow\">or</span> <span class=\"n\">slugify</span><span class=\"p\">(</span><span
    class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;path&quot;</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">stem</span><span class=\"p\">))</span>\n\n
    \           <span class=\"nd\">@pydantic</span><span class=\"o\">.</span><span
    class=\"n\">validator</span><span class=\"p\">(</span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">pre</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span> <span class=\"n\">always</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">def</span> <span class=\"nf\">index_slug_is_empty</span><span
    class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"p\">,</span> <span class=\"n\">values</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"o\">==</span>
    <span class=\"s2\">&quot;index&quot;</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"s2\">&quot;&quot;</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">v</span>\n\n            <span class=\"nd\">@pydantic</span><span
    class=\"o\">.</span><span class=\"n\">validator</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;slug&quot;</span><span class=\"p\">,</span> <span class=\"n\">pre</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
    class=\"n\">always</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n            <span class=\"k\">def</span> <span class=\"nf\">no_double_slash_in_slug</span><span
    class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"p\">,</span> <span class=\"n\">values</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">v</span>\n                <span class=\"k\">return</span>
    <span class=\"n\">v</span><span class=\"o\">.</span><span class=\"n\">replace</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;//&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;/&quot;</span><span class=\"p\">)</span>\n\n            <span
    class=\"nd\">@pydantic</span><span class=\"o\">.</span><span class=\"n\">validator</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;href&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">pre</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span> <span class=\"n\">always</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n            <span class=\"k\">def</span>
    <span class=\"nf\">default_href</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
    class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">v</span>\n                <span class=\"k\">return</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;/</span><span class=\"si\">{</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s1\">&#39;slug&#39;</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">strip</span><span class=\"p\">(</span><span
    class=\"s1\">&#39;/&#39;</span><span class=\"p\">)</span><span class=\"si\">}</span><span
    class=\"s2\">/&quot;</span><span class=\"o\">.</span><span class=\"n\">replace</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;//&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;/&quot;</span><span class=\"p\">)</span>\n\n            <span
    class=\"nd\">@pydantic</span><span class=\"o\">.</span><span class=\"n\">validator</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;title&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">pre</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span> <span class=\"n\">always</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n            <span class=\"k\">def</span>
    <span class=\"nf\">title_title</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
    class=\"p\">):</span>\n                <span class=\"n\">title</span> <span class=\"o\">=</span>
    <span class=\"n\">v</span> <span class=\"ow\">or</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;path&quot;</span><span class=\"p\">])</span><span class=\"o\">.</span><span
    class=\"n\">stem</span><span class=\"o\">.</span><span class=\"n\">replace</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;-&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot; &quot;</span><span class=\"p\">)</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">title</span><span class=\"o\">.</span><span
    class=\"n\">title</span><span class=\"p\">()</span>\n\n            <span class=\"nd\">@pydantic</span><span
    class=\"o\">.</span><span class=\"n\">validator</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;date_time&quot;</span><span class=\"p\">,</span> <span class=\"n\">pre</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
    class=\"n\">always</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n            <span class=\"k\">def</span> <span class=\"nf\">dateparser_datetime</span><span
    class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"p\">,</span> <span class=\"n\">values</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">d</span> <span class=\"o\">=</span> <span class=\"n\">dateparser</span><span
    class=\"o\">.</span><span class=\"n\">parse</span><span class=\"p\">(</span><span
    class=\"n\">v</span><span class=\"p\">)</span>\n                    <span class=\"k\">if</span>
    <span class=\"n\">d</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                        <span class=\"k\">raise</span> <span
    class=\"ne\">ValueError</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s1\">&#39;&quot;</span><span class=\"si\">{</span><span class=\"n\">v</span><span
    class=\"si\">}</span><span class=\"s1\">&quot; is not a valid date&#39;</span><span
    class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span>\n\n
    \           <span class=\"nd\">@pydantic</span><span class=\"o\">.</span><span
    class=\"n\">validator</span><span class=\"p\">(</span><span class=\"s2\">&quot;date_time&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">pre</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span> <span class=\"n\">always</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">def</span> <span class=\"nf\">date_is_datetime</span><span
    class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"p\">,</span> <span class=\"n\">values</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span> <span class=\"ow\">and</span> <span class=\"s2\">&quot;date&quot;</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">values</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;markata&quot;</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">log</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;</span><span class=\"si\">{</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s1\">&#39;path&#39;</span><span class=\"p\">]</span><span
    class=\"si\">}</span><span class=\"s2\"> has no date&quot;</span><span class=\"p\">)</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">now</span><span class=\"p\">()</span>\n                <span class=\"k\">if</span>
    <span class=\"n\">v</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span>
    <span class=\"ow\">and</span> <span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;date&quot;</span><span class=\"p\">]</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">&quot;markata&quot;</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">log</span><span class=\"p\">(</span><span
    class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">values</span><span class=\"p\">[</span><span class=\"s1\">&#39;path&#39;</span><span
    class=\"p\">]</span><span class=\"si\">}</span><span class=\"s2\"> has no date&quot;</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">now</span><span class=\"p\">()</span>\n
    \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">v</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;date&quot;</span><span class=\"p\">],</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"p\">):</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;date&quot;</span><span class=\"p\">]</span>\n
    \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">date</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">combine</span><span class=\"p\">(</span><span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">time</span><span class=\"o\">.</span><span
    class=\"n\">min</span><span class=\"p\">)</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;date&quot;</span><span class=\"p\">],</span>
    <span class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">date</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">combine</span><span class=\"p\">(</span><span
    class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">&quot;date&quot;</span><span
    class=\"p\">],</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">time</span><span class=\"o\">.</span><span class=\"n\">min</span><span
    class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span>\n\n
    \           <span class=\"nd\">@pydantic</span><span class=\"o\">.</span><span
    class=\"n\">validator</span><span class=\"p\">(</span><span class=\"s2\">&quot;date_time&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">pre</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span> <span class=\"n\">always</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">def</span> <span class=\"nf\">mindate_time</span><span
    class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"p\">,</span> <span class=\"n\">values</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span> <span class=\"ow\">and</span> <span class=\"s2\">&quot;date&quot;</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">values</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;markata&quot;</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">log</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;</span><span class=\"si\">{</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s1\">&#39;path&#39;</span><span class=\"p\">]</span><span
    class=\"si\">}</span><span class=\"s2\"> has no date&quot;</span><span class=\"p\">)</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">min</span>\n                <span class=\"k\">if</span> <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;date&quot;</span><span class=\"p\">]</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;markata&quot;</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
    class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s1\">&#39;path&#39;</span><span class=\"p\">]</span><span class=\"si\">}</span><span
    class=\"s2\"> has no date&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">min</span>\n
    \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">v</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;date&quot;</span><span class=\"p\">],</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"p\">):</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;date&quot;</span><span class=\"p\">]</span>\n
    \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">date</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">combine</span><span class=\"p\">(</span><span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">time</span><span class=\"o\">.</span><span
    class=\"n\">min</span><span class=\"p\">)</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;date&quot;</span><span class=\"p\">],</span>
    <span class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">date</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">combine</span><span class=\"p\">(</span><span
    class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">&quot;date&quot;</span><span
    class=\"p\">],</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">time</span><span class=\"o\">.</span><span class=\"n\">min</span><span
    class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span>\n\n
    \           <span class=\"nd\">@pydantic</span><span class=\"o\">.</span><span
    class=\"n\">validator</span><span class=\"p\">(</span><span class=\"s2\">&quot;date&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">pre</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span> <span class=\"n\">always</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">def</span> <span class=\"nf\">dateparser_date</span><span
    class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"p\">,</span> <span class=\"n\">values</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">date</span><span class=\"o\">.</span><span class=\"n\">min</span>\n
    \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">d</span> <span class=\"o\">=</span> <span class=\"bp\">cls</span><span
    class=\"o\">.</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">precache</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">if</span> <span class=\"n\">d</span> <span class=\"ow\">is</span>
    <span class=\"ow\">not</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                       <span class=\"k\">return</span> <span class=\"n\">d</span>\n
    \                   <span class=\"n\">d</span> <span class=\"o\">=</span> <span
    class=\"n\">dateparser</span><span class=\"o\">.</span><span class=\"n\">parse</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">if</span> <span class=\"n\">d</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                        <span
    class=\"k\">raise</span> <span class=\"ne\">ValueError</span><span class=\"p\">(</span><span
    class=\"sa\">f</span><span class=\"s1\">&#39;&quot;</span><span class=\"si\">{</span><span
    class=\"n\">v</span><span class=\"si\">}</span><span class=\"s1\">&quot; is not
    a valid date&#39;</span><span class=\"p\">)</span>\n                    <span
    class=\"n\">d</span> <span class=\"o\">=</span> <span class=\"n\">d</span><span
    class=\"o\">.</span><span class=\"n\">date</span><span class=\"p\">()</span>\n
    \                   <span class=\"k\">with</span> <span class=\"bp\">cls</span><span
    class=\"o\">.</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">cache</span> <span class=\"k\">as</span> <span class=\"n\">cache</span><span
    class=\"p\">:</span>\n                        <span class=\"n\">cache</span><span
    class=\"o\">.</span><span class=\"n\">add</span><span class=\"p\">(</span><span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"n\">d</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">d</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='PostModelConfig' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>PostModelConfig <em class='small'>class</em></h2>\nConfiguration for the
    Post model</p>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">PostModelConfig <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nc\">PostModelConfig</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"s2\">&quot;Configuration for the Post model&quot;</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"fm\">__init__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"o\">**</span><span class=\"n\">data</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
    class=\"w\">                </span><span class=\"sd\">&quot;&quot;&quot;</span>\n\n<span
    class=\"sd\">                include: post attributes to include by default in
    Post</span>\n<span class=\"sd\">                model serialization.</span>\n<span
    class=\"sd\">                repr_include: post attributes to include by default
    in Post</span>\n<span class=\"sd\">                repr.  If `repr_include` is
    None, it will default to</span>\n<span class=\"sd\">                `include`,
    but it is likely that you want less in the repr</span>\n<span class=\"sd\">                than
    serialized output.</span>\n\n<span class=\"sd\">                example:</span>\n\n<span
    class=\"sd\">                ``` toml title=&#39;markata.toml&#39;</span>\n<span
    class=\"sd\">                [markata.post_model]</span>\n<span class=\"sd\">
    \               include = [&#39;date&#39;, &#39;description&#39;, &#39;published&#39;,</span>\n<span
    class=\"sd\">                    &#39;slug&#39;, &#39;title&#39;, &#39;content&#39;,
    &#39;html&#39;]</span>\n<span class=\"sd\">                repr_include = [&#39;date&#39;,
    &#39;description&#39;, &#39;published&#39;, &#39;slug&#39;, &#39;title&#39;]</span>\n<span
    class=\"sd\">                ```</span>\n<span class=\"sd\">                &quot;&quot;&quot;</span>\n
    \               <span class=\"nb\">super</span><span class=\"p\">()</span><span
    class=\"o\">.</span><span class=\"fm\">__init__</span><span class=\"p\">(</span><span
    class=\"o\">**</span><span class=\"n\">data</span><span class=\"p\">)</span>\n\n
    \           <span class=\"n\">include</span><span class=\"p\">:</span> <span class=\"n\">List</span><span
    class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">]</span> <span
    class=\"o\">=</span> <span class=\"p\">[</span>\n                <span class=\"s2\">&quot;date&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;description&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;published&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;title&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;content&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;html&quot;</span><span
    class=\"p\">,</span>\n            <span class=\"p\">]</span>\n            <span
    class=\"n\">repr_include</span><span class=\"p\">:</span> <span class=\"n\">Optional</span><span
    class=\"p\">[</span><span class=\"n\">List</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]]</span> <span class=\"o\">=</span>
    <span class=\"p\">[</span>\n                <span class=\"s2\">&quot;date&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;description&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;published&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;title&quot;</span><span
    class=\"p\">,</span>\n            <span class=\"p\">]</span>\n            <span
    class=\"n\">export_include</span><span class=\"p\">:</span> <span class=\"n\">Optional</span><span
    class=\"p\">[</span><span class=\"n\">List</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]]</span> <span class=\"o\">=</span>
    <span class=\"p\">[</span>\n                <span class=\"s2\">&quot;date&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;description&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;published&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"s2\">&quot;title&quot;</span><span
    class=\"p\">,</span>\n            <span class=\"p\">]</span>\n\n            <span
    class=\"nd\">@pydantic</span><span class=\"o\">.</span><span class=\"n\">validator</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;repr_include&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">pre</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span> <span class=\"n\">always</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n            <span class=\"k\">def</span>
    <span class=\"nf\">repr_include_validator</span><span class=\"p\">(</span><span
    class=\"bp\">cls</span><span class=\"p\">,</span> <span class=\"n\">v</span><span
    class=\"p\">,</span> <span class=\"o\">*</span><span class=\"p\">,</span> <span
    class=\"n\">values</span><span class=\"p\">):</span>\n                <span class=\"k\">if</span>
    <span class=\"n\">v</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">v</span>\n                <span class=\"k\">return</span>
    <span class=\"n\">values</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;include&quot;</span><span class=\"p\">,</span>
    <span class=\"kc\">None</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    <span class=\"nc\">Config</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">post_model</span><span class=\"p\">:</span> <span
    class=\"n\">PostModelConfig</span> <span class=\"o\">=</span> <span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">Field</span><span class=\"p\">(</span><span
    class=\"n\">default_factory</span><span class=\"o\">=</span><span class=\"n\">PostModelConfig</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='post_model'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>post_model <em class='small'>function</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">post_model
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
    <span class=\"nf\">post_model</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">post_models</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">Post</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    class <h2 id='PostFactory' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>PostFactory <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">PostFactory
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
    <span class=\"nc\">PostFactory</span><span class=\"p\">(</span><span class=\"n\">ModelFactory</span><span
    class=\"p\">):</span>\n            <span class=\"n\">__model__</span> <span class=\"o\">=</span>
    <span class=\"n\">Post</span>\n            <span class=\"n\">__model__</span>
    <span class=\"o\">=</span> <span class=\"n\">Post</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"n\">Pretty</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='__repr_args__' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'><strong>repr_args</strong> <em class='small'>method</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\"><strong>repr_args</strong>
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
    <span class=\"nf\">__repr_args__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"s2\">&quot;ReprArgs&quot;</span><span
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
    class=\"n\">repr_include</span>\n                <span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='key' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>key
    <em class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">key <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">key</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">List</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]:</span>\n                <span class=\"k\">return</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">make_hash</span><span class=\"p\">(</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">slug</span><span class=\"p\">,</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">href</span><span class=\"p\">,</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">published</span><span class=\"p\">,</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">description</span><span
    class=\"p\">,</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">content</span><span class=\"p\">,</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">date</span><span class=\"p\">,</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">title</span><span class=\"p\">,</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">template</span><span class=\"p\">,</span>\n                <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='metadata' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>metadata <em class='small'>method</em></h2>\nfor backwards compatability</p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">metadata
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
    <span class=\"nf\">metadata</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">Dict</span><span class=\"p\">:</span>\n
    \               <span class=\"s2\">&quot;for backwards compatability&quot;</span>\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"vm\">__dict__</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='to_dict' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>to_dict <em class='small'>method</em></h2>\nfor backwards compatability</p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">to_dict
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
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">Dict</span><span class=\"p\">:</span>\n
    \               <span class=\"s2\">&quot;for backwards compatability&quot;</span>\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"vm\">__dict__</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='__getitem__' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'><strong>getitem</strong> <em class='small'>method</em></h2>\nfor backwards
    compatability</p>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\"><strong>getitem</strong> <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"fm\">__getitem__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">item</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span
    class=\"p\">:</span>\n                <span class=\"s2\">&quot;for backwards compatability&quot;</span>\n
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
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">key</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"n\">item</span><span class=\"p\">:</span> <span
    class=\"n\">Any</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
    class=\"s2\">&quot;for backwards compatability&quot;</span>\n                <span
    class=\"nb\">setattr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
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
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">item</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"n\">default</span><span class=\"p\">:</span>
    <span class=\"n\">Any</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Any</span><span class=\"p\">:</span>\n                <span
    class=\"s2\">&quot;for backwards compatability&quot;</span>\n                <span
    class=\"k\">return</span> <span class=\"nb\">getattr</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"n\">item</span><span
    class=\"p\">,</span> <span class=\"n\">default</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">List</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">]:</span>\n                <span class=\"s2\">&quot;for
    backwards compatability&quot;</span>\n                <span class=\"k\">return</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"vm\">__dict__</span><span
    class=\"o\">.</span><span class=\"n\">keys</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='yaml' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>yaml
    <em class='small'>method</em></h2>\ndump model to yaml</p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">yaml
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
    <span class=\"nf\">yaml</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n<span
    class=\"w\">                </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">                dump model to yaml</span>\n<span class=\"sd\">                &quot;&quot;&quot;</span>\n
    \               <span class=\"kn\">import</span> <span class=\"nn\">yaml</span>\n\n
    \               <span class=\"k\">return</span> <span class=\"n\">yaml</span><span
    class=\"o\">.</span><span class=\"n\">dump</span><span class=\"p\">(</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">dict</span><span class=\"p\">(</span>\n                        <span
    class=\"n\">include</span><span class=\"o\">=</span><span class=\"p\">{</span><span
    class=\"n\">i</span><span class=\"p\">:</span> <span class=\"kc\">True</span>
    <span class=\"k\">for</span> <span class=\"n\">i</span> <span class=\"ow\">in</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">post_model</span><span class=\"o\">.</span><span class=\"n\">include</span><span
    class=\"p\">}</span>\n                    <span class=\"p\">),</span>\n                    <span
    class=\"n\">Dumper</span><span class=\"o\">=</span><span class=\"n\">yaml</span><span
    class=\"o\">.</span><span class=\"n\">CDumper</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! method
    <h2 id='markdown' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>markdown
    <em class='small'>method</em></h2>\ndump model to markdown</p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">markdown
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
    <span class=\"nf\">markdown</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n<span
    class=\"w\">                </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">                dump model to markdown</span>\n<span class=\"sd\">
    \               &quot;&quot;&quot;</span>\n\n                <span class=\"kn\">import</span>
    <span class=\"nn\">yaml</span>\n\n                <span class=\"n\">frontmatter</span>
    <span class=\"o\">=</span> <span class=\"n\">yaml</span><span class=\"o\">.</span><span
    class=\"n\">dump</span><span class=\"p\">(</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">dict</span><span class=\"p\">(</span>\n
    \                       <span class=\"n\">include</span><span class=\"o\">=</span><span
    class=\"p\">{</span>\n                            <span class=\"n\">i</span><span
    class=\"p\">:</span> <span class=\"kc\">True</span>\n                            <span
    class=\"k\">for</span> <span class=\"n\">i</span> <span class=\"ow\">in</span>
    <span class=\"p\">[</span>\n                                <span class=\"n\">_i</span>\n
    \                               <span class=\"k\">for</span> <span class=\"n\">_i</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">post_model</span><span class=\"o\">.</span><span
    class=\"n\">include</span>\n                                <span class=\"k\">if</span>
    <span class=\"n\">_i</span> <span class=\"o\">!=</span> <span class=\"s2\">&quot;content&quot;</span>\n
    \                           <span class=\"p\">]</span>\n                        <span
    class=\"p\">}</span>\n                    <span class=\"p\">),</span>\n                    <span
    class=\"n\">Dumper</span><span class=\"o\">=</span><span class=\"n\">yaml</span><span
    class=\"o\">.</span><span class=\"n\">CDumper</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">)</span>\n                <span class=\"n\">post</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;---</span><span class=\"se\">\\n</span><span
    class=\"s2\">&quot;</span>\n                <span class=\"n\">post</span> <span
    class=\"o\">+=</span> <span class=\"n\">frontmatter</span>\n                <span
    class=\"n\">post</span> <span class=\"o\">+=</span> <span class=\"s2\">&quot;---</span><span
    class=\"se\">\\n\\n</span><span class=\"s2\">&quot;</span>\n\n                <span
    class=\"k\">if</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">content</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">post</span> <span class=\"o\">+=</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">content</span>\n                <span class=\"k\">return</span>
    <span class=\"n\">post</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='parse_file'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>parse_file <em class='small'>method</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">parse_file
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
    <span class=\"nf\">parse_file</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">path</span><span class=\"p\">:</span> <span class=\"n\">Union</span><span
    class=\"p\">[</span><span class=\"n\">Path</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">],</span> <span class=\"o\">**</span><span
    class=\"n\">kwargs</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">:</span>\n                <span
    class=\"k\">if</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">path</span><span class=\"p\">,</span> <span class=\"n\">Path</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">if</span> <span class=\"n\">path</span><span
    class=\"o\">.</span><span class=\"n\">suffix</span> <span class=\"ow\">in</span>
    <span class=\"p\">[</span><span class=\"s2\">&quot;.md&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;.markdown&quot;</span><span class=\"p\">]:</span>\n                        <span
    class=\"k\">return</span> <span class=\"bp\">cls</span><span class=\"o\">.</span><span
    class=\"n\">parse_markdown</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">path</span><span class=\"o\">=</span><span class=\"n\">path</span><span
    class=\"p\">,</span> <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">)</span>\n                <span class=\"k\">elif</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">path</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">):</span>\n                    <span
    class=\"k\">if</span> <span class=\"n\">path</span><span class=\"o\">.</span><span
    class=\"n\">endswith</span><span class=\"p\">(</span><span class=\"s2\">&quot;.md&quot;</span><span
    class=\"p\">)</span> <span class=\"ow\">or</span> <span class=\"n\">path</span><span
    class=\"o\">.</span><span class=\"n\">endswith</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;.markdown&quot;</span><span class=\"p\">):</span>\n                        <span
    class=\"k\">return</span> <span class=\"bp\">cls</span><span class=\"o\">.</span><span
    class=\"n\">parse_markdown</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">path</span><span class=\"o\">=</span><span class=\"n\">path</span><span
    class=\"p\">,</span> <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"nb\">super</span><span
    class=\"p\">(</span><span class=\"n\">Post</span><span class=\"p\">,</span> <span
    class=\"bp\">cls</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">parse_file</span><span class=\"p\">(</span><span class=\"n\">path</span><span
    class=\"p\">,</span> <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='parse_markdown'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>parse_markdown <em
    class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">parse_markdown <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">parse_markdown</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">path</span><span class=\"p\">:</span> <span class=\"n\">Union</span><span
    class=\"p\">[</span><span class=\"n\">Path</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">],</span> <span class=\"o\">**</span><span
    class=\"n\">kwargs</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">:</span>\n                <span
    class=\"k\">if</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">path</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
    class=\"p\">):</span>\n                    <span class=\"n\">path</span> <span
    class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">path</span><span class=\"p\">)</span>\n                <span class=\"n\">text</span>
    <span class=\"o\">=</span> <span class=\"n\">path</span><span class=\"o\">.</span><span
    class=\"n\">read_text</span><span class=\"p\">()</span>\n                <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                    <span class=\"n\">_</span><span
    class=\"p\">,</span> <span class=\"n\">fm</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"n\">content</span> <span class=\"o\">=</span>
    <span class=\"n\">text</span><span class=\"o\">.</span><span class=\"n\">split</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;---</span><span class=\"se\">\\n</span><span
    class=\"s2\">&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"n\">content</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;---</span><span
    class=\"se\">\\n</span><span class=\"s2\">&quot;</span><span class=\"o\">.</span><span
    class=\"n\">join</span><span class=\"p\">(</span><span class=\"n\">content</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">fm</span> <span class=\"o\">=</span>
    <span class=\"n\">yaml</span><span class=\"o\">.</span><span class=\"n\">load</span><span
    class=\"p\">(</span><span class=\"n\">fm</span><span class=\"p\">,</span> <span
    class=\"n\">Loader</span><span class=\"o\">=</span><span class=\"n\">yaml</span><span
    class=\"o\">.</span><span class=\"n\">CBaseLoader</span><span class=\"p\">)</span>\n
    \                   <span class=\"k\">except</span> <span class=\"n\">yaml</span><span
    class=\"o\">.</span><span class=\"n\">YAMLError</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">fm</span> <span class=\"o\">=</span>
    <span class=\"p\">{}</span>\n                <span class=\"k\">except</span> <span
    class=\"ne\">ValueError</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">fm</span> <span class=\"o\">=</span> <span class=\"p\">{}</span>\n
    \                   <span class=\"n\">content</span> <span class=\"o\">=</span>
    <span class=\"n\">text</span>\n                <span class=\"k\">if</span> <span
    class=\"n\">fm</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span>
    <span class=\"ow\">or</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">fm</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
    class=\"p\">):</span>\n                    <span class=\"n\">fm</span> <span class=\"o\">=</span>
    <span class=\"p\">{}</span>\n\n                <span class=\"n\">post_args</span>
    <span class=\"o\">=</span> <span class=\"p\">{</span>\n                    <span
    class=\"s2\">&quot;markata&quot;</span><span class=\"p\">:</span> <span class=\"n\">markata</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;path&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">path</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;content&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">content</span><span class=\"p\">,</span>\n                    <span
    class=\"s2\">&quot;raw&quot;</span><span class=\"p\">:</span> <span class=\"n\">text</span><span
    class=\"p\">,</span>\n                    <span class=\"o\">**</span><span class=\"n\">fm</span><span
    class=\"p\">,</span>\n                <span class=\"p\">}</span>\n\n                <span
    class=\"k\">return</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">Post</span><span class=\"o\">.</span><span class=\"n\">parse_obj</span><span
    class=\"p\">(</span><span class=\"n\">post_args</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='dumps' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>dumps
    <em class='small'>method</em></h2>\ndumps raw article back out</p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">dumps
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
    <span class=\"nf\">dumps</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">):</span>\n<span class=\"w\">                </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">                dumps raw article back out</span>\n<span class=\"sd\">
    \               &quot;&quot;&quot;</span>\n                <span class=\"k\">return</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;---</span><span class=\"se\">\\n</span><span
    class=\"si\">{</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">yaml</span><span class=\"p\">()</span><span class=\"si\">}</span><span
    class=\"se\">\\n\\n</span><span class=\"s2\">---</span><span class=\"se\">\\n\\n</span><span
    class=\"si\">{</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">content</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='default_slug' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>default_slug <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">default_slug
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
    <span class=\"nf\">default_slug</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
    class=\"p\">):</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span>
    <span class=\"ow\">or</span> <span class=\"n\">slugify</span><span class=\"p\">(</span><span
    class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;path&quot;</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">stem</span><span class=\"p\">))</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='index_slug_is_empty' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>index_slug_is_empty <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">index_slug_is_empty
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
    <span class=\"nf\">index_slug_is_empty</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
    class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span>
    <span class=\"o\">==</span> <span class=\"s2\">&quot;index&quot;</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"s2\">&quot;&quot;</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='no_double_slash_in_slug' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>no_double_slash_in_slug <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">no_double_slash_in_slug
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
    <span class=\"nf\">no_double_slash_in_slug</span><span class=\"p\">(</span><span
    class=\"bp\">cls</span><span class=\"p\">,</span> <span class=\"n\">v</span><span
    class=\"p\">,</span> <span class=\"o\">*</span><span class=\"p\">,</span> <span
    class=\"n\">values</span><span class=\"p\">):</span>\n                <span class=\"k\">if</span>
    <span class=\"n\">v</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">v</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span><span
    class=\"o\">.</span><span class=\"n\">replace</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;//&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;/&quot;</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='default_href'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>default_href <em
    class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">default_href <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">default_href</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
    class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">v</span>\n                <span class=\"k\">return</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;/</span><span class=\"si\">{</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s1\">&#39;slug&#39;</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">strip</span><span class=\"p\">(</span><span
    class=\"s1\">&#39;/&#39;</span><span class=\"p\">)</span><span class=\"si\">}</span><span
    class=\"s2\">/&quot;</span><span class=\"o\">.</span><span class=\"n\">replace</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;//&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;/&quot;</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='title_title' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>title_title <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">title_title
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
    <span class=\"nf\">title_title</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
    class=\"p\">):</span>\n                <span class=\"n\">title</span> <span class=\"o\">=</span>
    <span class=\"n\">v</span> <span class=\"ow\">or</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;path&quot;</span><span class=\"p\">])</span><span class=\"o\">.</span><span
    class=\"n\">stem</span><span class=\"o\">.</span><span class=\"n\">replace</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;-&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot; &quot;</span><span class=\"p\">)</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">title</span><span class=\"o\">.</span><span
    class=\"n\">title</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='dateparser_datetime' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>dateparser_datetime <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">dateparser_datetime
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
    <span class=\"nf\">dateparser_datetime</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
    class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">d</span> <span class=\"o\">=</span> <span class=\"n\">dateparser</span><span
    class=\"o\">.</span><span class=\"n\">parse</span><span class=\"p\">(</span><span
    class=\"n\">v</span><span class=\"p\">)</span>\n                    <span class=\"k\">if</span>
    <span class=\"n\">d</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                        <span class=\"k\">raise</span> <span
    class=\"ne\">ValueError</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s1\">&#39;&quot;</span><span class=\"si\">{</span><span class=\"n\">v</span><span
    class=\"si\">}</span><span class=\"s1\">&quot; is not a valid date&#39;</span><span
    class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='date_is_datetime' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>date_is_datetime <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">date_is_datetime
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
    <span class=\"nf\">date_is_datetime</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
    class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span> <span class=\"ow\">and</span>
    <span class=\"s2\">&quot;date&quot;</span> <span class=\"ow\">not</span> <span
    class=\"ow\">in</span> <span class=\"n\">values</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;markata&quot;</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
    class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s1\">&#39;path&#39;</span><span class=\"p\">]</span><span class=\"si\">}</span><span
    class=\"s2\"> has no date&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">now</span><span
    class=\"p\">()</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span> <span class=\"ow\">and</span>
    <span class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">&quot;date&quot;</span><span
    class=\"p\">]</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;markata&quot;</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">log</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;</span><span class=\"si\">{</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s1\">&#39;path&#39;</span><span class=\"p\">]</span><span
    class=\"si\">}</span><span class=\"s2\"> has no date&quot;</span><span class=\"p\">)</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">now</span><span class=\"p\">()</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">v</span><span
    class=\"p\">,</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">datetime</span><span class=\"p\">):</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">v</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;date&quot;</span><span class=\"p\">],</span>
    <span class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">&quot;date&quot;</span><span
    class=\"p\">]</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">date</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">combine</span><span class=\"p\">(</span><span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">time</span><span class=\"o\">.</span><span
    class=\"n\">min</span><span class=\"p\">)</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;date&quot;</span><span class=\"p\">],</span>
    <span class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">date</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">combine</span><span class=\"p\">(</span><span
    class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">&quot;date&quot;</span><span
    class=\"p\">],</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">time</span><span class=\"o\">.</span><span class=\"n\">min</span><span
    class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='mindate_time' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>mindate_time <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">mindate_time
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
    <span class=\"nf\">mindate_time</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
    class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span> <span class=\"ow\">and</span>
    <span class=\"s2\">&quot;date&quot;</span> <span class=\"ow\">not</span> <span
    class=\"ow\">in</span> <span class=\"n\">values</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;markata&quot;</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
    class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s1\">&#39;path&#39;</span><span class=\"p\">]</span><span class=\"si\">}</span><span
    class=\"s2\"> has no date&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">min</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;date&quot;</span><span class=\"p\">]</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;markata&quot;</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
    class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s1\">&#39;path&#39;</span><span class=\"p\">]</span><span class=\"si\">}</span><span
    class=\"s2\"> has no date&quot;</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">min</span>\n
    \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">v</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;date&quot;</span><span class=\"p\">],</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"p\">):</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;date&quot;</span><span class=\"p\">]</span>\n
    \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">date</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">combine</span><span class=\"p\">(</span><span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">time</span><span class=\"o\">.</span><span
    class=\"n\">min</span><span class=\"p\">)</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;date&quot;</span><span class=\"p\">],</span>
    <span class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">date</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">combine</span><span class=\"p\">(</span><span
    class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">&quot;date&quot;</span><span
    class=\"p\">],</span> <span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">time</span><span class=\"o\">.</span><span class=\"n\">min</span><span
    class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='dateparser_date' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>dateparser_date <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">dateparser_date
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
    <span class=\"nf\">dateparser_date</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
    class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">date</span><span class=\"o\">.</span><span
    class=\"n\">min</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">d</span> <span class=\"o\">=</span> <span class=\"bp\">cls</span><span
    class=\"o\">.</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">precache</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">if</span> <span class=\"n\">d</span> <span class=\"ow\">is</span>
    <span class=\"ow\">not</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                       <span class=\"k\">return</span> <span class=\"n\">d</span>\n
    \                   <span class=\"n\">d</span> <span class=\"o\">=</span> <span
    class=\"n\">dateparser</span><span class=\"o\">.</span><span class=\"n\">parse</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">)</span>\n                    <span
    class=\"k\">if</span> <span class=\"n\">d</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                        <span
    class=\"k\">raise</span> <span class=\"ne\">ValueError</span><span class=\"p\">(</span><span
    class=\"sa\">f</span><span class=\"s1\">&#39;&quot;</span><span class=\"si\">{</span><span
    class=\"n\">v</span><span class=\"si\">}</span><span class=\"s1\">&quot; is not
    a valid date&#39;</span><span class=\"p\">)</span>\n                    <span
    class=\"n\">d</span> <span class=\"o\">=</span> <span class=\"n\">d</span><span
    class=\"o\">.</span><span class=\"n\">date</span><span class=\"p\">()</span>\n
    \                   <span class=\"k\">with</span> <span class=\"bp\">cls</span><span
    class=\"o\">.</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">cache</span> <span class=\"k\">as</span> <span class=\"n\">cache</span><span
    class=\"p\">:</span>\n                        <span class=\"n\">cache</span><span
    class=\"o\">.</span><span class=\"n\">add</span><span class=\"p\">(</span><span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"n\">d</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">d</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='__init__' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'><strong>init</strong> <em class='small'>method</em></h2>\ninclude: post
    attributes to include by default in Post\nmodel serialization.\nrepr_include:
    post attributes to include by default in Post\nrepr.  If <code>repr_include</code>
    is None, it will default to\n<code>include</code>, but it is likely that you want
    less in the repr\nthan serialized output.</p>\n<pre><code>example:\n\n``` toml
    title='markata.toml'\n[markata.post_model]\ninclude = ['date', 'description',
    'published',\n    'slug', 'title', 'content', 'html']\nrepr_include = ['date',
    'description', 'published', 'slug', 'title']\n```\n</code></pre>\n<div class=\"admonition
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
    class=\"p\">,</span> <span class=\"o\">**</span><span class=\"n\">data</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n<span class=\"w\">                </span><span class=\"sd\">&quot;&quot;&quot;</span>\n\n<span
    class=\"sd\">                include: post attributes to include by default in
    Post</span>\n<span class=\"sd\">                model serialization.</span>\n<span
    class=\"sd\">                repr_include: post attributes to include by default
    in Post</span>\n<span class=\"sd\">                repr.  If `repr_include` is
    None, it will default to</span>\n<span class=\"sd\">                `include`,
    but it is likely that you want less in the repr</span>\n<span class=\"sd\">                than
    serialized output.</span>\n\n<span class=\"sd\">                example:</span>\n\n<span
    class=\"sd\">                ``` toml title=&#39;markata.toml&#39;</span>\n<span
    class=\"sd\">                [markata.post_model]</span>\n<span class=\"sd\">
    \               include = [&#39;date&#39;, &#39;description&#39;, &#39;published&#39;,</span>\n<span
    class=\"sd\">                    &#39;slug&#39;, &#39;title&#39;, &#39;content&#39;,
    &#39;html&#39;]</span>\n<span class=\"sd\">                repr_include = [&#39;date&#39;,
    &#39;description&#39;, &#39;published&#39;, &#39;slug&#39;, &#39;title&#39;]</span>\n<span
    class=\"sd\">                ```</span>\n<span class=\"sd\">                &quot;&quot;&quot;</span>\n
    \               <span class=\"nb\">super</span><span class=\"p\">()</span><span
    class=\"o\">.</span><span class=\"fm\">__init__</span><span class=\"p\">(</span><span
    class=\"o\">**</span><span class=\"n\">data</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='repr_include_validator' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>repr_include_validator <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">repr_include_validator
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
    <span class=\"nf\">repr_include_validator</span><span class=\"p\">(</span><span
    class=\"bp\">cls</span><span class=\"p\">,</span> <span class=\"n\">v</span><span
    class=\"p\">,</span> <span class=\"o\">*</span><span class=\"p\">,</span> <span
    class=\"n\">values</span><span class=\"p\">):</span>\n                <span class=\"k\">if</span>
    <span class=\"n\">v</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">v</span>\n                <span class=\"k\">return</span>
    <span class=\"n\">values</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;include&quot;</span><span class=\"p\">,</span>
    <span class=\"kc\">None</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>"
  raw.md: ''
published: true
slug: markata/plugins/post-model
title: Post_Model.Py


---

None


!! class <h2 id='Post' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Post <em class='small'>class</em></h2>

???+ source "Post <em class='small'>source</em>"

```python

        class Post(pydantic.BaseModel, JupyterMixin):
            markata: Any = Field(None, exclude=True)
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
            template: Optional[str | Dict[str, str]] = "post.html"
            sidebar: Optional[Any] = None

            def __rich__(self) -> Pretty:
                return Pretty(self)

            def __repr_args__(self: "Post") -> "ReprArgs":
                return [
                    (key, value)
                    for key, value in self.__dict__.items()
                    if key in self.markata.config.post_model.repr_include
                ]

            @property
            def key(self: "Post") -> List[str]:
                return self.markata.make_hash(
                    self.slug,
                    self.href,
                    self.published,
                    self.description,
                    self.content,
                    self.date,
                    self.title,
                    self.template,
                )

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
                    "raw": text,
                    **fm,
                }

                return markata.Post.parse_obj(post_args)

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

            @pydantic.validator("slug", pre=True, always=True)
            def no_double_slash_in_slug(cls, v, *, values):
                if v is None:
                    return v
                return v.replace("//", "/")

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
            export_include: Optional[List[str]] = [
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
            __model__ = Post
```


!! method <h2 id='__rich__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__rich__ <em class='small'>method</em></h2>

???+ source "__rich__ <em class='small'>source</em>"

```python

        def __rich__(self) -> Pretty:
                return Pretty(self)
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


!! method <h2 id='key' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>key <em class='small'>method</em></h2>

???+ source "key <em class='small'>source</em>"

```python

        def key(self: "Post") -> List[str]:
                return self.markata.make_hash(
                    self.slug,
                    self.href,
                    self.published,
                    self.description,
                    self.content,
                    self.date,
                    self.title,
                    self.template,
                )
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
                    "raw": text,
                    **fm,
                }

                return markata.Post.parse_obj(post_args)
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


!! method <h2 id='no_double_slash_in_slug' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>no_double_slash_in_slug <em class='small'>method</em></h2>

???+ source "no_double_slash_in_slug <em class='small'>source</em>"

```python

        def no_double_slash_in_slug(cls, v, *, values):
                if v is None:
                    return v
                return v.replace("//", "/")
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

