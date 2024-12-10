---
content: "# Add head configuration\n\nThis snippet allows users to configure their
  head in `markata.toml`.\n\n``` html\n{{ config.get('head', {}).pop('text') if 'text'
  in config.get('head',{}).keys() }}\n{% for tag, meta in config.get('head', {}).items()
  %}\n    {% for _meta in meta %}\n        <{{ tag }}\n            {% for attr, value
  in _meta.items() %}{{ attr }}=\"{{ value }}\"{% endfor %}\n        />\n    {% endfor
  %}\n{% endfor %}\n```\n\nUsers can specify any sort of tag in their `markata.toml`\n\n```
  toml\n[[markata.head.meta]]\nname = \"og:type\"\ncontent = \"article\"\n\n[[markata.head.meta]]\nname
  = \"og:author\"\ncontent = \"Waylon Walker\"\n```\n\nThe above configuration becomes
  this once rendered.\n\n``` html\n<meta name='og:type' content='article' />\n<meta
  name='og:Author' content='Waylon Walker' />\n```\n\n!! Note\n\n    Article variables
  can be used for dynamic entries like canonical_url\n    ``` toml\n    [markata]\n
  \   url = \"markata.dev\"\n\n    [[markata.head.meta]]\n    href=\"{{ config.url
  }}/{{ slug }}/\"\n    rel=\"canonical\"\n    ```\n\nOptionally users can also specify
  plain text to be appended to the head of\ntheir documents.  This works well for
  things that involve full blocks.\n\n``` toml\n[[markata.head.text]]\nvalue = '''\n<script>\n
  \   console.log('hello world')\n</script>\n'''\n\n[[markata.head.text]]\nvalue='''\nhtml
  \ {\n    font-family: \"Space Mono\", monospace;\n    background: var(--color-bg);\n
  \   color: var(--color-text);\n}\n'''\n```\n\n## Add scripts to head\n\nMarkata
  config also supports adding scripts to the head via configuration.\n\n``` toml\n[[
  markata.head.script ]]\n    src = \"https://cdn.tailwindcss.com\"\n\n```\n\n\n!!
  class <h2 id='SilentUndefined' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>SilentUndefined <em class='small'>class</em></h2>\n\n???+ source \"SilentUndefined
  <em class='small'>source</em>\"\n\n```python\n\n        class SilentUndefined(Undefined):\n
  \           def _fail_with_undefined_error(self, *args, **kwargs):\n                return
  \"\"\n```\n\n\n!! function <h2 id='optional' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>optional <em class='small'>function</em></h2>\n\n???+ source \"optional <em
  class='small'>source</em>\"\n\n```python\n\n        def optional(*fields):\n            def
  dec(_cls):\n                for field in fields:\n                    _cls.__fields__[field].default
  = None\n                return _cls\n\n            if (\n                fields\n
  \               and inspect.isclass(fields[0])\n                and issubclass(fields[0],
  pydantic.BaseModel)\n            ):\n                cls = fields[0]\n                fields
  = cls.__fields__\n                return dec(cls)\n            return dec\n```\n\n\n!!
  class <h2 id='Style' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Style
  <em class='small'>class</em></h2>\n\n???+ source \"Style <em class='small'>source</em>\"\n\n```python\n\n
  \       class Style(pydantic.BaseModel):\n            color_bg: str = \"#1f2022\"\n
  \           color_bg_code: str = \"#1f2022\"\n            color_text: str = \"#eefbfe\"\n
  \           color_link: str = \"#fb30c4\"\n            color_accent: str = \"#e1bd00c9\"\n
  \           overlay_brightness: str = \".85\"\n            body_width: str = \"800px\"\n
  \           color_bg_light: str = \"#eefbfe\"\n            color_bg_code_light:
  str = \"#eefbfe\"\n            color_text_light: str = \"#1f2022\"\n            color_link_light:
  str = \"#fb30c4\"\n            color_accent_light: str = \"#ffeb00\"\n            overlay_brightness_light:
  str = \".95\"\n```\n\n\n!! class <h2 id='StyleOverrides' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>StyleOverrides <em class='small'>class</em></h2>\n\n???+
  source \"StyleOverrides <em class='small'>source</em>\"\n\n```python\n\n        class
  StyleOverrides(Style): ...\n```\n\n\n!! class <h2 id='Meta' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>Meta <em class='small'>class</em></h2>\n\n???+
  source \"Meta <em class='small'>source</em>\"\n\n```python\n\n        class Meta(pydantic.BaseModel):\n
  \           name: str\n            content: str\n```\n\n\n!! class <h2 id='Text'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Text <em class='small'>class</em></h2>\n\n???+
  source \"Text <em class='small'>source</em>\"\n\n```python\n\n        class Text(pydantic.BaseModel):\n
  \           value: str\n```\n\n\n!! class <h2 id='Link' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>Link <em class='small'>class</em></h2>\n\n???+
  source \"Link <em class='small'>source</em>\"\n\n```python\n\n        class Link(pydantic.BaseModel):\n
  \           rel: str = \"canonical\"\n            href: str\n```\n\n\n!! class <h2
  id='Script' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Script
  <em class='small'>class</em></h2>\n\n???+ source \"Script <em class='small'>source</em>\"\n\n```python\n\n
  \       class Script(pydantic.BaseModel):\n            src: str\n```\n\n\n!! class
  <h2 id='HeadConfig' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>HeadConfig
  <em class='small'>class</em></h2>\n\n???+ source \"HeadConfig <em class='small'>source</em>\"\n\n```python\n\n
  \       class HeadConfig(pydantic.BaseModel):\n            meta: List[Meta] = []\n
  \           link: List[Link] = []\n            script: List[Script] = []\n            text:
  Union[List[Text], str] = \"\"\n\n            @pydantic.validator(\"text\", pre=True)\n
  \           def text_to_list(cls, v):\n                if isinstance(v, list):\n
  \                   return \"\\n\".join([text[\"value\"] for text in v])\n                return
  v\n\n            @property\n            def html(self):\n                html =
  self.text\n                html += \"\\n\"\n                for meta in self.meta:\n
  \                   html += f'<meta name=\"{meta.name}\" content=\"{meta.content}\"
  />\\n'\n                for link in self.link:\n                    html += f'<link
  rel=\"{link.rel}\" href=\"{link.href}\" />\\n'\n                return html\n```\n\n\n!!
  class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Config
  <em class='small'>class</em></h2>\n\n???+ source \"Config <em class='small'>source</em>\"\n\n```python\n\n
  \       class Config(pydantic.BaseModel):\n            head: HeadConfig = HeadConfig()\n
  \           style: Style = Style()\n            post_template: Optional[Union[str
  | Dict[str, str]]] = \"post.html\"\n            dynamic_templates_dir: Path = Path(\".markata.cache/templates\")\n
  \           templates_dir: Union[Path, List[Path]] = pydantic.Field(Path(\"templates\"))\n\n
  \           env_options: dict = {}\n\n            @pydantic.model_validator(mode=\"after\")\n
  \           def dynamic_templates_in_templates_dir(self):\n                markata_templates
  = Path(__file__).parents[1] / \"templates\"\n\n                if isinstance(self.templates_dir,
  Path):\n                    self.templates_dir = [\n                        self.templates_dir,\n
  \                       markata_templates,\n                        self.dynamic_templates_dir,\n
  \                   ]\n\n                if markata_templates not in self.templates_dir:\n
  \                   self.templates_dir.append(markata_templates)\n\n                if
  self.dynamic_templates_dir not in self.templates_dir:\n                    self.templates_dir.append(self.dynamic_templates_dir)\n\n
  \               return self\n\n            @property\n            def jinja_loader(self):\n
  \               return jinja2.FileSystemLoader(self.templates_dir)\n\n            @property\n
  \           def jinja_env(\n                self,\n            ):\n                if
  hasattr(self, \"_jinja_env\"):\n                    return self._jinja_env\n                self.env_options.setdefault(\"loader\",
  self.jinja_loader)\n                self.env_options.setdefault(\"undefined\", SilentUndefined)\n
  \               self.env_options.setdefault(\"lstrip_blocks\", True)\n                self.env_options.setdefault(\"trim_blocks\",
  True)\n\n                env = jinja2.Environment(**self.env_options)\n\n                self._jinja_env
  = env\n                return env\n```\n\n\n!! class <h2 id='PostOverrides' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>PostOverrides <em class='small'>class</em></h2>\n\n???+
  source \"PostOverrides <em class='small'>source</em>\"\n\n```python\n\n        class
  PostOverrides(pydantic.BaseModel):\n            head: HeadConfig = HeadConfig()\n
  \           style: Style = StyleOverrides()\n```\n\n\n!! class <h2 id='Post' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>Post <em class='small'>class</em></h2>\n\n???+
  source \"Post <em class='small'>source</em>\"\n\n```python\n\n        class Post(pydantic.BaseModel):\n
  \           config_overrides: PostOverrides = PostOverrides()\n            template:
  Optional[str | Dict[str, str]] = None\n\n            @pydantic.validator(\"template\",
  pre=True, always=True)\n            def default_template(cls, v, *, values):\n                if
  v is None:\n                    return values[\"markata\"].config.post_template\n
  \               if isinstance(v, str):\n                    v = {\"index\": v}\n
  \               if isinstance(values[\"markata\"].config.post_template, str):\n
  \                   config_template = {\n                        \"index\": values[\"markata\"].config.post_template,\n
  \                   }\n                else:\n                    config_template
  = values[\"markata\"].config.post_template\n                return {**config_template,
  **v}\n```\n\n\n!! function <h2 id='config_model' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>config_model <em class='small'>function</em></h2>\n\n???+ source \"config_model
  <em class='small'>source</em>\"\n\n```python\n\n        def config_model(markata:
  \"Markata\") -> None:\n            markata.config_models.append(Config)\n```\n\n\n!!
  function <h2 id='post_model' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>post_model <em class='small'>function</em></h2>\n\n???+ source \"post_model
  <em class='small'>source</em>\"\n\n```python\n\n        def post_model(markata:
  \"Markata\") -> None:\n            markata.post_models.append(Post)\n```\n\n\n!!
  function <h2 id='configure' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>configure <em class='small'>function</em></h2>\n    Massages the configuration
  limitations of toml to allow a little bit easier\n    experience to the end user
  making configurations while allowing an simpler\n    jinja template.  This enablees
  the use of the `markata.head.text` list in\n    configuration.\n???+ source \"configure
  <em class='small'>source</em>\"\n\n```python\n\n        def configure(markata: \"Markata\")
  -> None:\n            \"\"\"\n            Massages the configuration limitations
  of toml to allow a little bit easier\n            experience to the end user making
  configurations while allowing an simpler\n            jinja template.  This enablees
  the use of the `markata.head.text` list in\n            configuration.\n            \"\"\"\n```\n\n\n!!
  function <h2 id='pre_render' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>pre_render <em class='small'>function</em></h2>\n    FOR EACH POST: Massages
  the configuration limitations of toml/yaml to allow\n    a little bit easier experience
  to the end user making configurations while\n    allowing an simpler jinja template.
  \ This enables the use of the\n    `markata.head.text` list in configuration.\n???+
  source \"pre_render <em class='small'>source</em>\"\n\n```python\n\n        def
  pre_render(markata: \"Markata\") -> None:\n            \"\"\"\n            FOR EACH
  POST: Massages the configuration limitations of toml/yaml to allow\n            a
  little bit easier experience to the end user making configurations while\n            allowing
  an simpler jinja template.  This enables the use of the\n            `markata.head.text`
  list in configuration.\n            \"\"\"\n\n            markata.config.dynamic_templates_dir.mkdir(parents=True,
  exist_ok=True)\n            head_template = markata.config.dynamic_templates_dir
  / \"head.html\"\n            head_template.write_text(\n                markata.config.jinja_env.get_template(\"dynamic_head.html\").render(\n
  \                   {\"markata\": markata}\n                ),\n            )\n\n
  \           for article in [a for a in markata.articles if \"config_overrides\"
  in a]:\n                raw_text = article.get(\"config_overrides\", {}).get(\"head\",
  {}).get(\"text\", \"\")\n\n                if isinstance(raw_text, list):\n                    article[\"config_overrides\"][\"head\"][\"text\"]
  = \"\\n\".join(\n                        flatten([t.values() for t in raw_text]),\n
  \                   )\n```\n\n\n!! function <h2 id='render' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>render <em class='small'>function</em></h2>\n\n???+
  source \"render <em class='small'>source</em>\"\n\n```python\n\n        def render(markata:
  \"Markata\") -> None:\n            with markata.cache as cache:\n                for
  article in markata.articles:\n                    html = render_article(markata=markata,
  cache=cache, article=article)\n                    article.html = html\n```\n\n\n!!
  function <h2 id='get_template' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>get_template <em class='small'>function</em></h2>\n\n???+ source \"get_template
  <em class='small'>source</em>\"\n\n```python\n\n        def get_template(markata,
  template):\n            try:\n                return markata.config.jinja_env.get_template(template)\n
  \           except jinja2.TemplateNotFound:\n                # try to load it as
  a file\n                ...\n\n            try:\n                return Template(Path(template).read_text(),
  undefined=SilentUndefined)\n            except FileNotFoundError:\n                #
  default to load it as a string\n                ...\n            return Template(template,
  undefined=SilentUndefined)\n```\n\n\n!! function <h2 id='render_article' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>render_article <em class='small'>function</em></h2>\n\n???+
  source \"render_article <em class='small'>source</em>\"\n\n```python\n\n        def
  render_article(markata, cache, article):\n            key = markata.make_hash(\n
  \               \"post_template\",\n                __version__,\n                article.key,\n
  \           )\n            html = markata.precache.get(key)\n\n            if html
  is not None:\n                return html\n\n            if isinstance(article.template,
  str):\n                template = get_template(markata, article.template)\n                html
  = render_template(markata, article, template)\n\n            if isinstance(article.template,
  dict):\n                html = {\n                    slug: render_template(markata,
  article, get_template(markata, template))\n                    for slug, template
  in article.template.items()\n                }\n            cache.add(key, html,
  expire=markata.config.default_cache_expire)\n            return html\n```\n\n\n!!
  function <h2 id='render_template' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>render_template <em class='small'>function</em></h2>\n\n???+ source \"render_template
  <em class='small'>source</em>\"\n\n```python\n\n        def render_template(markata,
  article, template):\n            template = get_template(markata, template)\n            merged_config
  = markata.config\n            # TODO do we need to handle merge??\n            #
  if head_template:\n            #     head = eval(\n            #         head_template.render(\n
  \           #             __version__=__version__,\n            #             config=_full_config,\n
  \           #             **article,\n            #         )\n            #     )\n\n
  \           # merged_config = {\n            #     **_full_config,\n            #
  \    **{\"head\": head},\n            # }\n\n            # merged_config = always_merger.merge(\n
  \           #     merged_config,\n            #     copy.deepcopy(\n            #
  \        article.get(\n            #             \"config_overrides\",\n            #
  \            {},\n            #         )\n            #     ),\n            # )\n\n
  \           html = template.render(\n                __version__=__version__,\n
  \               markata=markata,\n                body=article.article_html,\n                config=merged_config,\n
  \               post=article,\n            )\n            return html\n```\n\n\n!!
  function <h2 id='save' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>save
  <em class='small'>function</em></h2>\n\n???+ source \"save <em class='small'>source</em>\"\n\n```python\n\n
  \       def save(markata: \"Markata\") -> None:\n            linked_templates =
  [\n                t\n                for t in markata.config.jinja_env.list_templates()\n
  \               if t.endswith(\"css\") or t.endswith(\"js\") or t.endswith(\"xsl\")\n
  \           ]\n            for template in linked_templates:\n                template
  = get_template(markata, template)\n                css = template.render(markata=markata,
  __version__=__version__)\n                Path(markata.config.output_dir / Path(template.filename).name).write_text(css)\n```\n\n\n!!
  function <h2 id='cli' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>cli
  <em class='small'>function</em></h2>\n    Markata hook to implement base cli commands.\n???+
  source \"cli <em class='small'>source</em>\"\n\n```python\n\n        def cli(app:
  typer.Typer, markata: \"Markata\") -> None:\n            \"\"\"\n            Markata
  hook to implement base cli commands.\n            \"\"\"\n\n            templates_app
  = typer.Typer()\n            app.add_typer(templates_app)\n\n            @templates_app.callback()\n
  \           def templates():\n                \"template management\"\n\n            @templates_app.command()\n
  \           def show(\n                template: str = typer.Argument(None, help=\"template
  to show\"),\n                theme: str = typer.Option(None, help=\"pygments syntax
  theme\"),\n            ) -> None:\n                markata.console.quiet = True\n
  \               if template:\n                    template = get_template(markata,
  template)\n\n                    markata.console.quiet = False\n                    markata.console.print(template.filename)\n
  \                   if theme is None or theme.lower() == \"none\":\n                        markata.console.print(Path(template.filename).read_text())\n
  \                   else:\n                        syntax = Syntax.from_path(template.filename,
  theme=theme)\n                        markata.console.print(syntax)\n\n                    return\n
  \               templates = markata.config.jinja_env.list_templates()\n                markata.console.quiet
  = False\n                markata.console.print(\"Templates directories:\", style=\"green
  underline\")\n\n                markata_templates = Path(__file__).parents[1] /
  \"templates\"\n                for dir in markata.config.templates_dir:\n                    if
  dir == markata.config.dynamic_templates_dir:\n                        markata.console.print(\n
  \                           f\"[gold3]{dir}[/][grey50] (dynamically created templates
  from configuration)[/] [gold3]\\[markata.config.dynamic_templates_dir][/]\",\n                            style=\"red\",\n
  \                       )\n                    elif dir == markata_templates:\n
  \                       markata.console.print(\n                            f\"[cyan]{dir}[/][grey50]
  (built-in)[/]\", style=\"red\"\n                        )\n                    else:\n
  \                       markata.console.print(\n                            f\"[orchid]{dir}[/]
  [orchid]\\[markata.config.templates_dir][/]\",\n                            style=\"red\",\n
  \                       )\n\n                markata.console.print()\n                markata.console.print(\n
  \                   \"Available Templates: [white]name -> path[/]\", style=\"green
  underline\"\n                )\n                for template in templates:\n                    source,
  file, uptodate = markata.config.jinja_env.loader.get_source(\n                        markata.config.jinja_env,
  template\n                    )\n\n                    if Path(file).is_relative_to(markata.config.dynamic_templates_dir):\n
  \                       markata.console.print(\n                            f\"[gold3]{template}
  -> [red]{file}[/] [grey50](dynamic)[/]\"\n                        )\n                    elif
  Path(file).is_relative_to(markata_templates):\n                        markata.console.print(\n
  \                           f\"[cyan]{template} -> [red]{file}[/] [grey50](built-in)[/]\"\n
  \                       )\n                    else:\n                        markata.console.print(f\"[orchid]{template}[/]
  -> [red]{file}[/]\")\n```\n\n\n!! method <h2 id='_fail_with_undefined_error' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>_fail_with_undefined_error <em class='small'>method</em></h2>\n\n???+
  source \"_fail_with_undefined_error <em class='small'>source</em>\"\n\n```python\n\n
  \       def _fail_with_undefined_error(self, *args, **kwargs):\n                return
  \"\"\n```\n\n\n!! function <h2 id='dec' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>dec <em class='small'>function</em></h2>\n\n???+ source \"dec <em class='small'>source</em>\"\n\n```python\n\n
  \       def dec(_cls):\n                for field in fields:\n                    _cls.__fields__[field].default
  = None\n                return _cls\n```\n\n\n!! method <h2 id='text_to_list' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>text_to_list <em class='small'>method</em></h2>\n\n???+
  source \"text_to_list <em class='small'>source</em>\"\n\n```python\n\n        def
  text_to_list(cls, v):\n                if isinstance(v, list):\n                    return
  \"\\n\".join([text[\"value\"] for text in v])\n                return v\n```\n\n\n!!
  method <h2 id='html' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>html
  <em class='small'>method</em></h2>\n\n???+ source \"html <em class='small'>source</em>\"\n\n```python\n\n
  \       def html(self):\n                html = self.text\n                html
  += \"\\n\"\n                for meta in self.meta:\n                    html +=
  f'<meta name=\"{meta.name}\" content=\"{meta.content}\" />\\n'\n                for
  link in self.link:\n                    html += f'<link rel=\"{link.rel}\" href=\"{link.href}\"
  />\\n'\n                return html\n```\n\n\n!! method <h2 id='dynamic_templates_in_templates_dir'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>dynamic_templates_in_templates_dir
  <em class='small'>method</em></h2>\n\n???+ source \"dynamic_templates_in_templates_dir
  <em class='small'>source</em>\"\n\n```python\n\n        def dynamic_templates_in_templates_dir(self):\n
  \               markata_templates = Path(__file__).parents[1] / \"templates\"\n\n
  \               if isinstance(self.templates_dir, Path):\n                    self.templates_dir
  = [\n                        self.templates_dir,\n                        markata_templates,\n
  \                       self.dynamic_templates_dir,\n                    ]\n\n                if
  markata_templates not in self.templates_dir:\n                    self.templates_dir.append(markata_templates)\n\n
  \               if self.dynamic_templates_dir not in self.templates_dir:\n                    self.templates_dir.append(self.dynamic_templates_dir)\n\n
  \               return self\n```\n\n\n!! method <h2 id='jinja_loader' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>jinja_loader <em class='small'>method</em></h2>\n\n???+
  source \"jinja_loader <em class='small'>source</em>\"\n\n```python\n\n        def
  jinja_loader(self):\n                return jinja2.FileSystemLoader(self.templates_dir)\n```\n\n\n!!
  method <h2 id='jinja_env' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>jinja_env <em class='small'>method</em></h2>\n\n???+ source \"jinja_env <em
  class='small'>source</em>\"\n\n```python\n\n        def jinja_env(\n                self,\n
  \           ):\n                if hasattr(self, \"_jinja_env\"):\n                    return
  self._jinja_env\n                self.env_options.setdefault(\"loader\", self.jinja_loader)\n
  \               self.env_options.setdefault(\"undefined\", SilentUndefined)\n                self.env_options.setdefault(\"lstrip_blocks\",
  True)\n                self.env_options.setdefault(\"trim_blocks\", True)\n\n                env
  = jinja2.Environment(**self.env_options)\n\n                self._jinja_env = env\n
  \               return env\n```\n\n\n!! method <h2 id='default_template' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>default_template <em class='small'>method</em></h2>\n\n???+
  source \"default_template <em class='small'>source</em>\"\n\n```python\n\n        def
  default_template(cls, v, *, values):\n                if v is None:\n                    return
  values[\"markata\"].config.post_template\n                if isinstance(v, str):\n
  \                   v = {\"index\": v}\n                if isinstance(values[\"markata\"].config.post_template,
  str):\n                    config_template = {\n                        \"index\":
  values[\"markata\"].config.post_template,\n                    }\n                else:\n
  \                   config_template = values[\"markata\"].config.post_template\n
  \               return {**config_template, **v}\n```\n\n\n!! function <h2 id='templates'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>templates <em class='small'>function</em></h2>\n
  \   template management\n???+ source \"templates <em class='small'>source</em>\"\n\n```python\n\n
  \       def templates():\n                \"template management\"\n```\n\n\n!! function
  <h2 id='show' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>show
  <em class='small'>function</em></h2>\n\n???+ source \"show <em class='small'>source</em>\"\n\n```python\n\n
  \       def show(\n                template: str = typer.Argument(None, help=\"template
  to show\"),\n                theme: str = typer.Option(None, help=\"pygments syntax
  theme\"),\n            ) -> None:\n                markata.console.quiet = True\n
  \               if template:\n                    template = get_template(markata,
  template)\n\n                    markata.console.quiet = False\n                    markata.console.print(template.filename)\n
  \                   if theme is None or theme.lower() == \"none\":\n                        markata.console.print(Path(template.filename).read_text())\n
  \                   else:\n                        syntax = Syntax.from_path(template.filename,
  theme=theme)\n                        markata.console.print(syntax)\n\n                    return\n
  \               templates = markata.config.jinja_env.list_templates()\n                markata.console.quiet
  = False\n                markata.console.print(\"Templates directories:\", style=\"green
  underline\")\n\n                markata_templates = Path(__file__).parents[1] /
  \"templates\"\n                for dir in markata.config.templates_dir:\n                    if
  dir == markata.config.dynamic_templates_dir:\n                        markata.console.print(\n
  \                           f\"[gold3]{dir}[/][grey50] (dynamically created templates
  from configuration)[/] [gold3]\\[markata.config.dynamic_templates_dir][/]\",\n                            style=\"red\",\n
  \                       )\n                    elif dir == markata_templates:\n
  \                       markata.console.print(\n                            f\"[cyan]{dir}[/][grey50]
  (built-in)[/]\", style=\"red\"\n                        )\n                    else:\n
  \                       markata.console.print(\n                            f\"[orchid]{dir}[/]
  [orchid]\\[markata.config.templates_dir][/]\",\n                            style=\"red\",\n
  \                       )\n\n                markata.console.print()\n                markata.console.print(\n
  \                   \"Available Templates: [white]name -> path[/]\", style=\"green
  underline\"\n                )\n                for template in templates:\n                    source,
  file, uptodate = markata.config.jinja_env.loader.get_source(\n                        markata.config.jinja_env,
  template\n                    )\n\n                    if Path(file).is_relative_to(markata.config.dynamic_templates_dir):\n
  \                       markata.console.print(\n                            f\"[gold3]{template}
  -> [red]{file}[/] [grey50](dynamic)[/]\"\n                        )\n                    elif
  Path(file).is_relative_to(markata_templates):\n                        markata.console.print(\n
  \                           f\"[cyan]{template} -> [red]{file}[/] [grey50](built-in)[/]\"\n
  \                       )\n                    else:\n                        markata.console.print(f\"[orchid]{template}[/]
  -> [red]{file}[/]\")\n```\n\n"
date: 0001-01-01
description: This snippet allows users to configure their head in  Users can specify
  any sort of tag in their  The above configuration becomes this once rendered. !
  Optional
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Post_Template.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"This snippet allows users to configure
    their head in  Users can specify any sort of tag in their  The above configuration
    becomes this once rendered. ! Optional\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Post_Template.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"This snippet allows users to configure
    their head in  Users can specify any sort of tag in their  The above configuration
    becomes this once rendered. ! Optional\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Post_Template.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <h1 id=\"add-head-configuration\">Add
    head configuration <a class=\"header-anchor\" href=\"#add-head-configuration\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>This snippet allows
    users to configure their head in <code>markata.toml</code>.</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span>{{ config.get(&#39;head&#39;,
    {}).pop(&#39;text&#39;) if &#39;text&#39; in config.get(&#39;head&#39;,{}).keys()
    }}\n{% for tag, meta in config.get(&#39;head&#39;, {}).items() %}\n    {% for
    _meta in meta %}\n        <span class=\"err\">&lt;</span>{{ tag }}\n            {%
    for attr, value in _meta.items() %}{{ attr }}=&quot;{{ value }}&quot;{% endfor
    %}\n        /&gt;\n    {% endfor %}\n{% endfor %}\n</pre></div>\n\n</pre>\n\n<p>Users
    can specify any sort of tag in their <code>markata.toml</code></p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.head.meta]]</span>\n<span
    class=\"n\">name</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;og:type&quot;</span>\n<span class=\"n\">content</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;article&quot;</span>\n\n<span class=\"k\">[[markata.head.meta]]</span>\n<span
    class=\"n\">name</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;og:author&quot;</span>\n<span class=\"n\">content</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;Waylon Walker&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>The
    above configuration becomes this once rendered.</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"p\">&lt;</span><span
    class=\"nt\">meta</span> <span class=\"na\">name</span><span class=\"o\">=</span><span
    class=\"s\">&#39;og:type&#39;</span> <span class=\"na\">content</span><span class=\"o\">=</span><span
    class=\"s\">&#39;article&#39;</span> <span class=\"p\">/&gt;</span>\n<span class=\"p\">&lt;</span><span
    class=\"nt\">meta</span> <span class=\"na\">name</span><span class=\"o\">=</span><span
    class=\"s\">&#39;og:Author&#39;</span> <span class=\"na\">content</span><span
    class=\"o\">=</span><span class=\"s\">&#39;Waylon Walker&#39;</span> <span class=\"p\">/&gt;</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    Note</p>\n<pre><code>Article variables can be used for dynamic entries like canonical_url\n```
    toml\n[markata]\nurl = &quot;markata.dev&quot;\n\n[[markata.head.meta]]\nhref=&quot;{{
    config.url }}/{{ slug }}/&quot;\nrel=&quot;canonical&quot;\n```\n</code></pre>\n<p>Optionally
    users can also specify plain text to be appended to the head of\ntheir documents.
    \ This works well for things that involve full blocks.</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.head.text]]</span>\n<span
    class=\"n\">value</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s1\">&#39;&#39;&#39;</span>\n<span class=\"s1\">&lt;script&gt;</span>\n<span
    class=\"s1\">    console.log(&#39;hello world&#39;)</span>\n<span class=\"s1\">&lt;/script&gt;</span>\n<span
    class=\"s1\">&#39;&#39;&#39;</span>\n\n<span class=\"k\">[[markata.head.text]]</span>\n<span
    class=\"n\">value</span><span class=\"o\">=</span><span class=\"s1\">&#39;&#39;&#39;</span>\n<span
    class=\"s1\">html  {</span>\n<span class=\"s1\">    font-family: &quot;Space Mono&quot;,
    monospace;</span>\n<span class=\"s1\">    background: var(--color-bg);</span>\n<span
    class=\"s1\">    color: var(--color-text);</span>\n<span class=\"s1\">}</span>\n<span
    class=\"s1\">&#39;&#39;&#39;</span>\n</pre></div>\n\n</pre>\n\n<h2 id=\"add-scripts-to-head\">Add
    scripts to head <a class=\"header-anchor\" href=\"#add-scripts-to-head\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Markata config also
    supports adding scripts to the head via configuration.</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[</span><span
    class=\"w\"> </span><span class=\"k\">markata.head.script</span><span class=\"w\">
    </span><span class=\"k\">]]</span>\n<span class=\"w\">    </span><span class=\"n\">src</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;https://cdn.tailwindcss.com&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='SilentUndefined' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>SilentUndefined <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">SilentUndefined
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
    <span class=\"nc\">SilentUndefined</span><span class=\"p\">(</span><span class=\"n\">Undefined</span><span
    class=\"p\">):</span>\n            <span class=\"k\">def</span> <span class=\"nf\">_fail_with_undefined_error</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"n\">args</span><span class=\"p\">,</span> <span
    class=\"o\">**</span><span class=\"n\">kwargs</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">return</span> <span class=\"s2\">&quot;&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='optional' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>optional <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">optional
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
    <span class=\"nf\">optional</span><span class=\"p\">(</span><span class=\"o\">*</span><span
    class=\"n\">fields</span><span class=\"p\">):</span>\n            <span class=\"k\">def</span>
    <span class=\"nf\">dec</span><span class=\"p\">(</span><span class=\"n\">_cls</span><span
    class=\"p\">):</span>\n                <span class=\"k\">for</span> <span class=\"n\">field</span>
    <span class=\"ow\">in</span> <span class=\"n\">fields</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">_cls</span><span class=\"o\">.</span><span
    class=\"n\">__fields__</span><span class=\"p\">[</span><span class=\"n\">field</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">default</span>
    <span class=\"o\">=</span> <span class=\"kc\">None</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">_cls</span>\n\n            <span class=\"k\">if</span>
    <span class=\"p\">(</span>\n                <span class=\"n\">fields</span>\n
    \               <span class=\"ow\">and</span> <span class=\"n\">inspect</span><span
    class=\"o\">.</span><span class=\"n\">isclass</span><span class=\"p\">(</span><span
    class=\"n\">fields</span><span class=\"p\">[</span><span class=\"mi\">0</span><span
    class=\"p\">])</span>\n                <span class=\"ow\">and</span> <span class=\"nb\">issubclass</span><span
    class=\"p\">(</span><span class=\"n\">fields</span><span class=\"p\">[</span><span
    class=\"mi\">0</span><span class=\"p\">],</span> <span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">)</span>\n
    \           <span class=\"p\">):</span>\n                <span class=\"bp\">cls</span>
    <span class=\"o\">=</span> <span class=\"n\">fields</span><span class=\"p\">[</span><span
    class=\"mi\">0</span><span class=\"p\">]</span>\n                <span class=\"n\">fields</span>
    <span class=\"o\">=</span> <span class=\"bp\">cls</span><span class=\"o\">.</span><span
    class=\"n\">__fields__</span>\n                <span class=\"k\">return</span>
    <span class=\"n\">dec</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">)</span>\n            <span class=\"k\">return</span> <span class=\"n\">dec</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='Style' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Style
    <em class='small'>class</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">Style <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nc\">Style</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">color_bg</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;#1f2022&quot;</span>\n
    \           <span class=\"n\">color_bg_code</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;#1f2022&quot;</span>\n
    \           <span class=\"n\">color_text</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;#eefbfe&quot;</span>\n
    \           <span class=\"n\">color_link</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;#fb30c4&quot;</span>\n
    \           <span class=\"n\">color_accent</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;#e1bd00c9&quot;</span>\n
    \           <span class=\"n\">overlay_brightness</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;.85&quot;</span>\n
    \           <span class=\"n\">body_width</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;800px&quot;</span>\n
    \           <span class=\"n\">color_bg_light</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;#eefbfe&quot;</span>\n
    \           <span class=\"n\">color_bg_code_light</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;#eefbfe&quot;</span>\n
    \           <span class=\"n\">color_text_light</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;#1f2022&quot;</span>\n
    \           <span class=\"n\">color_link_light</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;#fb30c4&quot;</span>\n
    \           <span class=\"n\">color_accent_light</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;#ffeb00&quot;</span>\n
    \           <span class=\"n\">overlay_brightness_light</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;.95&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='StyleOverrides' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>StyleOverrides <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">StyleOverrides
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
    <span class=\"nc\">StyleOverrides</span><span class=\"p\">(</span><span class=\"n\">Style</span><span
    class=\"p\">):</span> <span class=\"o\">...</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='Meta' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Meta
    <em class='small'>class</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">Meta <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nc\">Meta</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">name</span><span class=\"p\">:</span> <span class=\"nb\">str</span>\n
    \           <span class=\"n\">content</span><span class=\"p\">:</span> <span class=\"nb\">str</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='Text' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Text
    <em class='small'>class</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">Text <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nc\">Text</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">value</span><span class=\"p\">:</span> <span class=\"nb\">str</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='Link' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Link
    <em class='small'>class</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">Link <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nc\">Link</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">rel</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;canonical&quot;</span>\n            <span
    class=\"n\">href</span><span class=\"p\">:</span> <span class=\"nb\">str</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='Script' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Script
    <em class='small'>class</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">Script <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nc\">Script</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">src</span><span class=\"p\">:</span> <span class=\"nb\">str</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='HeadConfig' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>HeadConfig <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">HeadConfig
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
    <span class=\"nc\">HeadConfig</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">meta</span><span class=\"p\">:</span> <span class=\"n\">List</span><span
    class=\"p\">[</span><span class=\"n\">Meta</span><span class=\"p\">]</span> <span
    class=\"o\">=</span> <span class=\"p\">[]</span>\n            <span class=\"n\">link</span><span
    class=\"p\">:</span> <span class=\"n\">List</span><span class=\"p\">[</span><span
    class=\"n\">Link</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"p\">[]</span>\n            <span class=\"n\">script</span><span class=\"p\">:</span>
    <span class=\"n\">List</span><span class=\"p\">[</span><span class=\"n\">Script</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"p\">[]</span>\n
    \           <span class=\"n\">text</span><span class=\"p\">:</span> <span class=\"n\">Union</span><span
    class=\"p\">[</span><span class=\"n\">List</span><span class=\"p\">[</span><span
    class=\"n\">Text</span><span class=\"p\">],</span> <span class=\"nb\">str</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;&quot;</span>\n\n
    \           <span class=\"nd\">@pydantic</span><span class=\"o\">.</span><span
    class=\"n\">validator</span><span class=\"p\">(</span><span class=\"s2\">&quot;text&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">pre</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n            <span class=\"k\">def</span>
    <span class=\"nf\">text_to_list</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">):</span>\n                <span
    class=\"k\">if</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"nb\">list</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span><span
    class=\"o\">.</span><span class=\"n\">join</span><span class=\"p\">([</span><span
    class=\"n\">text</span><span class=\"p\">[</span><span class=\"s2\">&quot;value&quot;</span><span
    class=\"p\">]</span> <span class=\"k\">for</span> <span class=\"n\">text</span>
    <span class=\"ow\">in</span> <span class=\"n\">v</span><span class=\"p\">])</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">v</span>\n\n
    \           <span class=\"nd\">@property</span>\n            <span class=\"k\">def</span>
    <span class=\"nf\">html</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">):</span>\n                <span class=\"n\">html</span> <span class=\"o\">=</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">text</span>\n
    \               <span class=\"n\">html</span> <span class=\"o\">+=</span> <span
    class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span>\n
    \               <span class=\"k\">for</span> <span class=\"n\">meta</span> <span
    class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">meta</span><span class=\"p\">:</span>\n                    <span class=\"n\">html</span>
    <span class=\"o\">+=</span> <span class=\"sa\">f</span><span class=\"s1\">&#39;&lt;meta
    name=&quot;</span><span class=\"si\">{</span><span class=\"n\">meta</span><span
    class=\"o\">.</span><span class=\"n\">name</span><span class=\"si\">}</span><span
    class=\"s1\">&quot; content=&quot;</span><span class=\"si\">{</span><span class=\"n\">meta</span><span
    class=\"o\">.</span><span class=\"n\">content</span><span class=\"si\">}</span><span
    class=\"s1\">&quot; /&gt;</span><span class=\"se\">\\n</span><span class=\"s1\">&#39;</span>\n
    \               <span class=\"k\">for</span> <span class=\"n\">link</span> <span
    class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">link</span><span class=\"p\">:</span>\n                    <span class=\"n\">html</span>
    <span class=\"o\">+=</span> <span class=\"sa\">f</span><span class=\"s1\">&#39;&lt;link
    rel=&quot;</span><span class=\"si\">{</span><span class=\"n\">link</span><span
    class=\"o\">.</span><span class=\"n\">rel</span><span class=\"si\">}</span><span
    class=\"s1\">&quot; href=&quot;</span><span class=\"si\">{</span><span class=\"n\">link</span><span
    class=\"o\">.</span><span class=\"n\">href</span><span class=\"si\">}</span><span
    class=\"s1\">&quot; /&gt;</span><span class=\"se\">\\n</span><span class=\"s1\">&#39;</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">html</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    \           <span class=\"n\">head</span><span class=\"p\">:</span> <span class=\"n\">HeadConfig</span>
    <span class=\"o\">=</span> <span class=\"n\">HeadConfig</span><span class=\"p\">()</span>\n
    \           <span class=\"n\">style</span><span class=\"p\">:</span> <span class=\"n\">Style</span>
    <span class=\"o\">=</span> <span class=\"n\">Style</span><span class=\"p\">()</span>\n
    \           <span class=\"n\">post_template</span><span class=\"p\">:</span> <span
    class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"n\">Union</span><span
    class=\"p\">[</span><span class=\"nb\">str</span> <span class=\"o\">|</span> <span
    class=\"n\">Dict</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"nb\">str</span><span class=\"p\">]]]</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;post.html&quot;</span>\n            <span
    class=\"n\">dynamic_templates_dir</span><span class=\"p\">:</span> <span class=\"n\">Path</span>
    <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;.markata.cache/templates&quot;</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">templates_dir</span><span class=\"p\">:</span> <span
    class=\"n\">Union</span><span class=\"p\">[</span><span class=\"n\">Path</span><span
    class=\"p\">,</span> <span class=\"n\">List</span><span class=\"p\">[</span><span
    class=\"n\">Path</span><span class=\"p\">]]</span> <span class=\"o\">=</span>
    <span class=\"n\">pydantic</span><span class=\"o\">.</span><span class=\"n\">Field</span><span
    class=\"p\">(</span><span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;templates&quot;</span><span class=\"p\">))</span>\n\n            <span
    class=\"n\">env_options</span><span class=\"p\">:</span> <span class=\"nb\">dict</span>
    <span class=\"o\">=</span> <span class=\"p\">{}</span>\n\n            <span class=\"nd\">@pydantic</span><span
    class=\"o\">.</span><span class=\"n\">model_validator</span><span class=\"p\">(</span><span
    class=\"n\">mode</span><span class=\"o\">=</span><span class=\"s2\">&quot;after&quot;</span><span
    class=\"p\">)</span>\n            <span class=\"k\">def</span> <span class=\"nf\">dynamic_templates_in_templates_dir</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">):</span>\n
    \               <span class=\"n\">markata_templates</span> <span class=\"o\">=</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"vm\">__file__</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">parents</span><span
    class=\"p\">[</span><span class=\"mi\">1</span><span class=\"p\">]</span> <span
    class=\"o\">/</span> <span class=\"s2\">&quot;templates&quot;</span>\n\n                <span
    class=\"k\">if</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">templates_dir</span><span
    class=\"p\">,</span> <span class=\"n\">Path</span><span class=\"p\">):</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">templates_dir</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
    \                       <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">templates_dir</span><span class=\"p\">,</span>\n                        <span
    class=\"n\">markata_templates</span><span class=\"p\">,</span>\n                        <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">dynamic_templates_dir</span><span
    class=\"p\">,</span>\n                    <span class=\"p\">]</span>\n\n                <span
    class=\"k\">if</span> <span class=\"n\">markata_templates</span> <span class=\"ow\">not</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">templates_dir</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">templates_dir</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">markata_templates</span><span class=\"p\">)</span>\n\n                <span
    class=\"k\">if</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">dynamic_templates_dir</span> <span class=\"ow\">not</span> <span class=\"ow\">in</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">templates_dir</span><span
    class=\"p\">:</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">templates_dir</span><span class=\"o\">.</span><span
    class=\"n\">append</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">dynamic_templates_dir</span><span class=\"p\">)</span>\n\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span>\n\n
    \           <span class=\"nd\">@property</span>\n            <span class=\"k\">def</span>
    <span class=\"nf\">jinja_loader</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">):</span>\n                <span class=\"k\">return</span> <span class=\"n\">jinja2</span><span
    class=\"o\">.</span><span class=\"n\">FileSystemLoader</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">templates_dir</span><span
    class=\"p\">)</span>\n\n            <span class=\"nd\">@property</span>\n            <span
    class=\"k\">def</span> <span class=\"nf\">jinja_env</span><span class=\"p\">(</span>\n
    \               <span class=\"bp\">self</span><span class=\"p\">,</span>\n            <span
    class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"nb\">hasattr</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;_jinja_env&quot;</span><span class=\"p\">):</span>\n                    <span
    class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_jinja_env</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">env_options</span><span class=\"o\">.</span><span
    class=\"n\">setdefault</span><span class=\"p\">(</span><span class=\"s2\">&quot;loader&quot;</span><span
    class=\"p\">,</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">jinja_loader</span><span class=\"p\">)</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">env_options</span><span
    class=\"o\">.</span><span class=\"n\">setdefault</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;undefined&quot;</span><span class=\"p\">,</span> <span class=\"n\">SilentUndefined</span><span
    class=\"p\">)</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">env_options</span><span class=\"o\">.</span><span class=\"n\">setdefault</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;lstrip_blocks&quot;</span><span class=\"p\">,</span>
    <span class=\"kc\">True</span><span class=\"p\">)</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">env_options</span><span
    class=\"o\">.</span><span class=\"n\">setdefault</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;trim_blocks&quot;</span><span class=\"p\">,</span> <span class=\"kc\">True</span><span
    class=\"p\">)</span>\n\n                <span class=\"n\">env</span> <span class=\"o\">=</span>
    <span class=\"n\">jinja2</span><span class=\"o\">.</span><span class=\"n\">Environment</span><span
    class=\"p\">(</span><span class=\"o\">**</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">env_options</span><span class=\"p\">)</span>\n\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_jinja_env</span> <span class=\"o\">=</span> <span class=\"n\">env</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">env</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='PostOverrides' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>PostOverrides <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">PostOverrides
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
    <span class=\"nc\">PostOverrides</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">head</span><span class=\"p\">:</span> <span class=\"n\">HeadConfig</span>
    <span class=\"o\">=</span> <span class=\"n\">HeadConfig</span><span class=\"p\">()</span>\n
    \           <span class=\"n\">style</span><span class=\"p\">:</span> <span class=\"n\">Style</span>
    <span class=\"o\">=</span> <span class=\"n\">StyleOverrides</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">config_overrides</span><span class=\"p\">:</span>
    <span class=\"n\">PostOverrides</span> <span class=\"o\">=</span> <span class=\"n\">PostOverrides</span><span
    class=\"p\">()</span>\n            <span class=\"n\">template</span><span class=\"p\">:</span>
    <span class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"nb\">str</span>
    <span class=\"o\">|</span> <span class=\"n\">Dict</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
    class=\"p\">]]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n\n
    \           <span class=\"nd\">@pydantic</span><span class=\"o\">.</span><span
    class=\"n\">validator</span><span class=\"p\">(</span><span class=\"s2\">&quot;template&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">pre</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span> <span class=\"n\">always</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">def</span> <span class=\"nf\">default_template</span><span
    class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"p\">,</span> <span class=\"n\">values</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;markata&quot;</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">post_template</span>\n
    \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">v</span> <span class=\"o\">=</span> <span class=\"p\">{</span><span
    class=\"s2\">&quot;index&quot;</span><span class=\"p\">:</span> <span class=\"n\">v</span><span
    class=\"p\">}</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;markata&quot;</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">post_template</span><span
    class=\"p\">,</span> <span class=\"nb\">str</span><span class=\"p\">):</span>\n
    \                   <span class=\"n\">config_template</span> <span class=\"o\">=</span>
    <span class=\"p\">{</span>\n                        <span class=\"s2\">&quot;index&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;markata&quot;</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">post_template</span><span
    class=\"p\">,</span>\n                    <span class=\"p\">}</span>\n                <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                    <span class=\"n\">config_template</span>
    <span class=\"o\">=</span> <span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;markata&quot;</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">post_template</span>\n
    \               <span class=\"k\">return</span> <span class=\"p\">{</span><span
    class=\"o\">**</span><span class=\"n\">config_template</span><span class=\"p\">,</span>
    <span class=\"o\">**</span><span class=\"n\">v</span><span class=\"p\">}</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    function <h2 id='post_model' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>post_model <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">post_model
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
    function <h2 id='configure' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>configure <em class='small'>function</em></h2>\nMassages the configuration
    limitations of toml to allow a little bit easier\nexperience to the end user making
    configurations while allowing an simpler\njinja template.  This enablees the use
    of the <code>markata.head.text</code> list in\nconfiguration.</p>\n<div class=\"admonition
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
    <span class=\"nf\">configure</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            Massages the configuration limitations of toml to allow
    a little bit easier</span>\n<span class=\"sd\">            experience to the end
    user making configurations while allowing an simpler</span>\n<span class=\"sd\">
    \           jinja template.  This enablees the use of the `markata.head.text`
    list in</span>\n<span class=\"sd\">            configuration.</span>\n<span class=\"sd\">
    \           &quot;&quot;&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!! function
    <h2 id='pre_render' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>pre_render
    <em class='small'>function</em></h2>\nFOR EACH POST: Massages the configuration
    limitations of toml/yaml to allow\na little bit easier experience to the end user
    making configurations while\nallowing an simpler jinja template.  This enables
    the use of the\n<code>markata.head.text</code> list in configuration.</p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">pre_render
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
    <span class=\"nf\">pre_render</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            FOR EACH POST: Massages the configuration limitations
    of toml/yaml to allow</span>\n<span class=\"sd\">            a little bit easier
    experience to the end user making configurations while</span>\n<span class=\"sd\">
    \           allowing an simpler jinja template.  This enables the use of the</span>\n<span
    class=\"sd\">            `markata.head.text` list in configuration.</span>\n<span
    class=\"sd\">            &quot;&quot;&quot;</span>\n\n            <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">dynamic_templates_dir</span><span class=\"o\">.</span><span class=\"n\">mkdir</span><span
    class=\"p\">(</span><span class=\"n\">parents</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span> <span class=\"n\">exist_ok</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">head_template</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">dynamic_templates_dir</span> <span class=\"o\">/</span>
    <span class=\"s2\">&quot;head.html&quot;</span>\n            <span class=\"n\">head_template</span><span
    class=\"o\">.</span><span class=\"n\">write_text</span><span class=\"p\">(</span>\n
    \               <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">jinja_env</span><span
    class=\"o\">.</span><span class=\"n\">get_template</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;dynamic_head.html&quot;</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">render</span><span class=\"p\">(</span>\n
    \                   <span class=\"p\">{</span><span class=\"s2\">&quot;markata&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">markata</span><span class=\"p\">}</span>\n
    \               <span class=\"p\">),</span>\n            <span class=\"p\">)</span>\n\n
    \           <span class=\"k\">for</span> <span class=\"n\">article</span> <span
    class=\"ow\">in</span> <span class=\"p\">[</span><span class=\"n\">a</span> <span
    class=\"k\">for</span> <span class=\"n\">a</span> <span class=\"ow\">in</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">articles</span>
    <span class=\"k\">if</span> <span class=\"s2\">&quot;config_overrides&quot;</span>
    <span class=\"ow\">in</span> <span class=\"n\">a</span><span class=\"p\">]:</span>\n
    \               <span class=\"n\">raw_text</span> <span class=\"o\">=</span> <span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;config_overrides&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">{})</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;head&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">{})</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;text&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;&quot;</span><span class=\"p\">)</span>\n\n
    \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">raw_text</span><span class=\"p\">,</span>
    <span class=\"nb\">list</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">article</span><span class=\"p\">[</span><span class=\"s2\">&quot;config_overrides&quot;</span><span
    class=\"p\">][</span><span class=\"s2\">&quot;head&quot;</span><span class=\"p\">][</span><span
    class=\"s2\">&quot;text&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span><span
    class=\"o\">.</span><span class=\"n\">join</span><span class=\"p\">(</span>\n
    \                       <span class=\"n\">flatten</span><span class=\"p\">([</span><span
    class=\"n\">t</span><span class=\"o\">.</span><span class=\"n\">values</span><span
    class=\"p\">()</span> <span class=\"k\">for</span> <span class=\"n\">t</span>
    <span class=\"ow\">in</span> <span class=\"n\">raw_text</span><span class=\"p\">]),</span>\n
    \                   <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='render' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>render <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">render
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
    <span class=\"nf\">render</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">with</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">cache</span> <span class=\"k\">as</span>
    <span class=\"n\">cache</span><span class=\"p\">:</span>\n                <span
    class=\"k\">for</span> <span class=\"n\">article</span> <span class=\"ow\">in</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">articles</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">html</span> <span
    class=\"o\">=</span> <span class=\"n\">render_article</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"o\">=</span><span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"n\">cache</span><span class=\"o\">=</span><span
    class=\"n\">cache</span><span class=\"p\">,</span> <span class=\"n\">article</span><span
    class=\"o\">=</span><span class=\"n\">article</span><span class=\"p\">)</span>\n
    \                   <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">html</span> <span class=\"o\">=</span> <span class=\"n\">html</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='get_template' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>get_template <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">get_template
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
    <span class=\"nf\">get_template</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"n\">template</span><span class=\"p\">):</span>\n
    \           <span class=\"k\">try</span><span class=\"p\">:</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">jinja_env</span><span
    class=\"o\">.</span><span class=\"n\">get_template</span><span class=\"p\">(</span><span
    class=\"n\">template</span><span class=\"p\">)</span>\n            <span class=\"k\">except</span>
    <span class=\"n\">jinja2</span><span class=\"o\">.</span><span class=\"n\">TemplateNotFound</span><span
    class=\"p\">:</span>\n                <span class=\"c1\"># try to load it as a
    file</span>\n                <span class=\"o\">...</span>\n\n            <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                <span class=\"k\">return</span>
    <span class=\"n\">Template</span><span class=\"p\">(</span><span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">template</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">read_text</span><span class=\"p\">(),</span>
    <span class=\"n\">undefined</span><span class=\"o\">=</span><span class=\"n\">SilentUndefined</span><span
    class=\"p\">)</span>\n            <span class=\"k\">except</span> <span class=\"ne\">FileNotFoundError</span><span
    class=\"p\">:</span>\n                <span class=\"c1\"># default to load it
    as a string</span>\n                <span class=\"o\">...</span>\n            <span
    class=\"k\">return</span> <span class=\"n\">Template</span><span class=\"p\">(</span><span
    class=\"n\">template</span><span class=\"p\">,</span> <span class=\"n\">undefined</span><span
    class=\"o\">=</span><span class=\"n\">SilentUndefined</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='render_article' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>render_article <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">render_article
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
    <span class=\"nf\">render_article</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"n\">cache</span><span class=\"p\">,</span>
    <span class=\"n\">article</span><span class=\"p\">):</span>\n            <span
    class=\"n\">key</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">make_hash</span><span class=\"p\">(</span>\n
    \               <span class=\"s2\">&quot;post_template&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">__version__</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">key</span><span class=\"p\">,</span>\n            <span class=\"p\">)</span>\n
    \           <span class=\"n\">html</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">precache</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"n\">key</span><span
    class=\"p\">)</span>\n\n            <span class=\"k\">if</span> <span class=\"n\">html</span>
    <span class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"n\">html</span>\n\n
    \           <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">template</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
    class=\"p\">):</span>\n                <span class=\"n\">template</span> <span
    class=\"o\">=</span> <span class=\"n\">get_template</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">,</span> <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">template</span><span class=\"p\">)</span>\n
    \               <span class=\"n\">html</span> <span class=\"o\">=</span> <span
    class=\"n\">render_template</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"n\">article</span><span class=\"p\">,</span>
    <span class=\"n\">template</span><span class=\"p\">)</span>\n\n            <span
    class=\"k\">if</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">template</span><span
    class=\"p\">,</span> <span class=\"nb\">dict</span><span class=\"p\">):</span>\n
    \               <span class=\"n\">html</span> <span class=\"o\">=</span> <span
    class=\"p\">{</span>\n                    <span class=\"n\">slug</span><span class=\"p\">:</span>
    <span class=\"n\">render_template</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"n\">article</span><span class=\"p\">,</span>
    <span class=\"n\">get_template</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"n\">template</span><span class=\"p\">))</span>\n
    \                   <span class=\"k\">for</span> <span class=\"n\">slug</span><span
    class=\"p\">,</span> <span class=\"n\">template</span> <span class=\"ow\">in</span>
    <span class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">template</span><span
    class=\"o\">.</span><span class=\"n\">items</span><span class=\"p\">()</span>\n
    \               <span class=\"p\">}</span>\n            <span class=\"n\">cache</span><span
    class=\"o\">.</span><span class=\"n\">add</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">,</span> <span class=\"n\">html</span><span
    class=\"p\">,</span> <span class=\"n\">expire</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">default_cache_expire</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">return</span> <span class=\"n\">html</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='render_template' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>render_template <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">render_template
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
    <span class=\"nf\">render_template</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"n\">article</span><span class=\"p\">,</span>
    <span class=\"n\">template</span><span class=\"p\">):</span>\n            <span
    class=\"n\">template</span> <span class=\"o\">=</span> <span class=\"n\">get_template</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">template</span><span class=\"p\">)</span>\n            <span
    class=\"n\">merged_config</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span>\n            <span class=\"c1\">#
    TODO do we need to handle merge??</span>\n            <span class=\"c1\"># if
    head_template:</span>\n            <span class=\"c1\">#     head = eval(</span>\n
    \           <span class=\"c1\">#         head_template.render(</span>\n            <span
    class=\"c1\">#             __version__=__version__,</span>\n            <span
    class=\"c1\">#             config=_full_config,</span>\n            <span class=\"c1\">#
    \            **article,</span>\n            <span class=\"c1\">#         )</span>\n
    \           <span class=\"c1\">#     )</span>\n\n            <span class=\"c1\">#
    merged_config = {</span>\n            <span class=\"c1\">#     **_full_config,</span>\n
    \           <span class=\"c1\">#     **{&quot;head&quot;: head},</span>\n            <span
    class=\"c1\"># }</span>\n\n            <span class=\"c1\"># merged_config = always_merger.merge(</span>\n
    \           <span class=\"c1\">#     merged_config,</span>\n            <span
    class=\"c1\">#     copy.deepcopy(</span>\n            <span class=\"c1\">#         article.get(</span>\n
    \           <span class=\"c1\">#             &quot;config_overrides&quot;,</span>\n
    \           <span class=\"c1\">#             {},</span>\n            <span class=\"c1\">#
    \        )</span>\n            <span class=\"c1\">#     ),</span>\n            <span
    class=\"c1\"># )</span>\n\n            <span class=\"n\">html</span> <span class=\"o\">=</span>
    <span class=\"n\">template</span><span class=\"o\">.</span><span class=\"n\">render</span><span
    class=\"p\">(</span>\n                <span class=\"n\">__version__</span><span
    class=\"o\">=</span><span class=\"n\">__version__</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"p\">,</span>\n                <span class=\"n\">body</span><span
    class=\"o\">=</span><span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">article_html</span><span class=\"p\">,</span>\n                <span
    class=\"n\">config</span><span class=\"o\">=</span><span class=\"n\">merged_config</span><span
    class=\"p\">,</span>\n                <span class=\"n\">post</span><span class=\"o\">=</span><span
    class=\"n\">article</span><span class=\"p\">,</span>\n            <span class=\"p\">)</span>\n
    \           <span class=\"k\">return</span> <span class=\"n\">html</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='save' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>save <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">save
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
    <span class=\"nf\">save</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">linked_templates</span> <span class=\"o\">=</span>
    <span class=\"p\">[</span>\n                <span class=\"n\">t</span>\n                <span
    class=\"k\">for</span> <span class=\"n\">t</span> <span class=\"ow\">in</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">jinja_env</span><span class=\"o\">.</span><span
    class=\"n\">list_templates</span><span class=\"p\">()</span>\n                <span
    class=\"k\">if</span> <span class=\"n\">t</span><span class=\"o\">.</span><span
    class=\"n\">endswith</span><span class=\"p\">(</span><span class=\"s2\">&quot;css&quot;</span><span
    class=\"p\">)</span> <span class=\"ow\">or</span> <span class=\"n\">t</span><span
    class=\"o\">.</span><span class=\"n\">endswith</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;js&quot;</span><span class=\"p\">)</span> <span class=\"ow\">or</span>
    <span class=\"n\">t</span><span class=\"o\">.</span><span class=\"n\">endswith</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;xsl&quot;</span><span class=\"p\">)</span>\n
    \           <span class=\"p\">]</span>\n            <span class=\"k\">for</span>
    <span class=\"n\">template</span> <span class=\"ow\">in</span> <span class=\"n\">linked_templates</span><span
    class=\"p\">:</span>\n                <span class=\"n\">template</span> <span
    class=\"o\">=</span> <span class=\"n\">get_template</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">,</span> <span class=\"n\">template</span><span
    class=\"p\">)</span>\n                <span class=\"n\">css</span> <span class=\"o\">=</span>
    <span class=\"n\">template</span><span class=\"o\">.</span><span class=\"n\">render</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"p\">,</span> <span class=\"n\">__version__</span><span
    class=\"o\">=</span><span class=\"n\">__version__</span><span class=\"p\">)</span>\n
    \               <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">output_dir</span> <span class=\"o\">/</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">template</span><span class=\"o\">.</span><span
    class=\"n\">filename</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">name</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">write_text</span><span class=\"p\">(</span><span class=\"n\">css</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='cli' class='admonition-title'
    style='margin:0;padding:.5rem 1rem;'>cli <em class='small'>function</em></h2>\nMarkata
    hook to implement base cli commands.</p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">cli <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">cli</span><span class=\"p\">(</span><span class=\"n\">app</span><span
    class=\"p\">:</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
    class=\"n\">Typer</span><span class=\"p\">,</span> <span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            Markata hook to implement base cli commands.</span>\n<span
    class=\"sd\">            &quot;&quot;&quot;</span>\n\n            <span class=\"n\">templates_app</span>
    <span class=\"o\">=</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
    class=\"n\">Typer</span><span class=\"p\">()</span>\n            <span class=\"n\">app</span><span
    class=\"o\">.</span><span class=\"n\">add_typer</span><span class=\"p\">(</span><span
    class=\"n\">templates_app</span><span class=\"p\">)</span>\n\n            <span
    class=\"nd\">@templates_app</span><span class=\"o\">.</span><span class=\"n\">callback</span><span
    class=\"p\">()</span>\n            <span class=\"k\">def</span> <span class=\"nf\">templates</span><span
    class=\"p\">():</span>\n                <span class=\"s2\">&quot;template management&quot;</span>\n\n
    \           <span class=\"nd\">@templates_app</span><span class=\"o\">.</span><span
    class=\"n\">command</span><span class=\"p\">()</span>\n            <span class=\"k\">def</span>
    <span class=\"nf\">show</span><span class=\"p\">(</span>\n                <span
    class=\"n\">template</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
    class=\"n\">Argument</span><span class=\"p\">(</span><span class=\"kc\">None</span><span
    class=\"p\">,</span> <span class=\"n\">help</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;template to show&quot;</span><span class=\"p\">),</span>\n
    \               <span class=\"n\">theme</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"n\">typer</span><span
    class=\"o\">.</span><span class=\"n\">Option</span><span class=\"p\">(</span><span
    class=\"kc\">None</span><span class=\"p\">,</span> <span class=\"n\">help</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;pygments syntax theme&quot;</span><span
    class=\"p\">),</span>\n            <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">quiet</span> <span class=\"o\">=</span>
    <span class=\"kc\">True</span>\n                <span class=\"k\">if</span> <span
    class=\"n\">template</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">template</span> <span class=\"o\">=</span> <span class=\"n\">get_template</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">template</span><span class=\"p\">)</span>\n\n                    <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">quiet</span> <span class=\"o\">=</span>
    <span class=\"kc\">False</span>\n                    <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">print</span><span class=\"p\">(</span><span class=\"n\">template</span><span
    class=\"o\">.</span><span class=\"n\">filename</span><span class=\"p\">)</span>\n
    \                   <span class=\"k\">if</span> <span class=\"n\">theme</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span> <span class=\"ow\">or</span>
    <span class=\"n\">theme</span><span class=\"o\">.</span><span class=\"n\">lower</span><span
    class=\"p\">()</span> <span class=\"o\">==</span> <span class=\"s2\">&quot;none&quot;</span><span
    class=\"p\">:</span>\n                        <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">print</span><span class=\"p\">(</span><span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">template</span><span class=\"o\">.</span><span
    class=\"n\">filename</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">read_text</span><span class=\"p\">())</span>\n                    <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">syntax</span> <span class=\"o\">=</span> <span class=\"n\">Syntax</span><span
    class=\"o\">.</span><span class=\"n\">from_path</span><span class=\"p\">(</span><span
    class=\"n\">template</span><span class=\"o\">.</span><span class=\"n\">filename</span><span
    class=\"p\">,</span> <span class=\"n\">theme</span><span class=\"o\">=</span><span
    class=\"n\">theme</span><span class=\"p\">)</span>\n                        <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">print</span><span class=\"p\">(</span><span
    class=\"n\">syntax</span><span class=\"p\">)</span>\n\n                    <span
    class=\"k\">return</span>\n                <span class=\"n\">templates</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">jinja_env</span><span
    class=\"o\">.</span><span class=\"n\">list_templates</span><span class=\"p\">()</span>\n
    \               <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">quiet</span>
    <span class=\"o\">=</span> <span class=\"kc\">False</span>\n                <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">print</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;Templates directories:&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">style</span><span class=\"o\">=</span><span class=\"s2\">&quot;green
    underline&quot;</span><span class=\"p\">)</span>\n\n                <span class=\"n\">markata_templates</span>
    <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"vm\">__file__</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">parents</span><span class=\"p\">[</span><span class=\"mi\">1</span><span
    class=\"p\">]</span> <span class=\"o\">/</span> <span class=\"s2\">&quot;templates&quot;</span>\n
    \               <span class=\"k\">for</span> <span class=\"nb\">dir</span> <span
    class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">templates_dir</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">if</span> <span class=\"nb\">dir</span>
    <span class=\"o\">==</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">dynamic_templates_dir</span><span
    class=\"p\">:</span>\n                        <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">print</span><span class=\"p\">(</span>\n                            <span
    class=\"sa\">f</span><span class=\"s2\">&quot;[gold3]</span><span class=\"si\">{</span><span
    class=\"nb\">dir</span><span class=\"si\">}</span><span class=\"s2\">[/][grey50]
    (dynamically created templates from configuration)[/] [gold3]\\[markata.config.dynamic_templates_dir][/]&quot;</span><span
    class=\"p\">,</span>\n                            <span class=\"n\">style</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;red&quot;</span><span class=\"p\">,</span>\n
    \                       <span class=\"p\">)</span>\n                    <span
    class=\"k\">elif</span> <span class=\"nb\">dir</span> <span class=\"o\">==</span>
    <span class=\"n\">markata_templates</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">print</span><span class=\"p\">(</span>\n
    \                           <span class=\"sa\">f</span><span class=\"s2\">&quot;[cyan]</span><span
    class=\"si\">{</span><span class=\"nb\">dir</span><span class=\"si\">}</span><span
    class=\"s2\">[/][grey50] (built-in)[/]&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">style</span><span class=\"o\">=</span><span class=\"s2\">&quot;red&quot;</span>\n
    \                       <span class=\"p\">)</span>\n                    <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">print</span><span class=\"p\">(</span>\n
    \                           <span class=\"sa\">f</span><span class=\"s2\">&quot;[orchid]</span><span
    class=\"si\">{</span><span class=\"nb\">dir</span><span class=\"si\">}</span><span
    class=\"s2\">[/] [orchid]\\[markata.config.templates_dir][/]&quot;</span><span
    class=\"p\">,</span>\n                            <span class=\"n\">style</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;red&quot;</span><span class=\"p\">,</span>\n
    \                       <span class=\"p\">)</span>\n\n                <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">print</span><span class=\"p\">()</span>\n                <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">print</span><span class=\"p\">(</span>\n                    <span
    class=\"s2\">&quot;Available Templates: [white]name -&gt; path[/]&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">style</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;green underline&quot;</span>\n                <span class=\"p\">)</span>\n
    \               <span class=\"k\">for</span> <span class=\"n\">template</span>
    <span class=\"ow\">in</span> <span class=\"n\">templates</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">source</span><span class=\"p\">,</span>
    <span class=\"n\">file</span><span class=\"p\">,</span> <span class=\"n\">uptodate</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">jinja_env</span><span
    class=\"o\">.</span><span class=\"n\">loader</span><span class=\"o\">.</span><span
    class=\"n\">get_source</span><span class=\"p\">(</span>\n                        <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">jinja_env</span><span class=\"p\">,</span>
    <span class=\"n\">template</span>\n                    <span class=\"p\">)</span>\n\n
    \                   <span class=\"k\">if</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">file</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">is_relative_to</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">dynamic_templates_dir</span><span class=\"p\">):</span>\n
    \                       <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">print</span><span
    class=\"p\">(</span>\n                            <span class=\"sa\">f</span><span
    class=\"s2\">&quot;[gold3]</span><span class=\"si\">{</span><span class=\"n\">template</span><span
    class=\"si\">}</span><span class=\"s2\"> -&gt; [red]</span><span class=\"si\">{</span><span
    class=\"n\">file</span><span class=\"si\">}</span><span class=\"s2\">[/] [grey50](dynamic)[/]&quot;</span>\n
    \                       <span class=\"p\">)</span>\n                    <span
    class=\"k\">elif</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">file</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">is_relative_to</span><span class=\"p\">(</span><span class=\"n\">markata_templates</span><span
    class=\"p\">):</span>\n                        <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">print</span><span class=\"p\">(</span>\n                            <span
    class=\"sa\">f</span><span class=\"s2\">&quot;[cyan]</span><span class=\"si\">{</span><span
    class=\"n\">template</span><span class=\"si\">}</span><span class=\"s2\"> -&gt;
    [red]</span><span class=\"si\">{</span><span class=\"n\">file</span><span class=\"si\">}</span><span
    class=\"s2\">[/] [grey50](built-in)[/]&quot;</span>\n                        <span
    class=\"p\">)</span>\n                    <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">print</span><span
    class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">&quot;[orchid]</span><span
    class=\"si\">{</span><span class=\"n\">template</span><span class=\"si\">}</span><span
    class=\"s2\">[/] -&gt; [red]</span><span class=\"si\">{</span><span class=\"n\">file</span><span
    class=\"si\">}</span><span class=\"s2\">[/]&quot;</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='_fail_with_undefined_error' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>_fail_with_undefined_error <em class='small'>method</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_fail_with_undefined_error
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
    <span class=\"nf\">_fail_with_undefined_error</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"n\">args</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"n\">kwargs</span><span class=\"p\">):</span>\n                <span class=\"k\">return</span>
    <span class=\"s2\">&quot;&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!! function
    <h2 id='dec' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>dec
    <em class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">dec <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">dec</span><span class=\"p\">(</span><span class=\"n\">_cls</span><span
    class=\"p\">):</span>\n                <span class=\"k\">for</span> <span class=\"n\">field</span>
    <span class=\"ow\">in</span> <span class=\"n\">fields</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">_cls</span><span class=\"o\">.</span><span
    class=\"n\">__fields__</span><span class=\"p\">[</span><span class=\"n\">field</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">default</span>
    <span class=\"o\">=</span> <span class=\"kc\">None</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">_cls</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='text_to_list' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>text_to_list <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">text_to_list
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
    <span class=\"nf\">text_to_list</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">):</span>\n                <span
    class=\"k\">if</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"nb\">list</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span><span
    class=\"o\">.</span><span class=\"n\">join</span><span class=\"p\">([</span><span
    class=\"n\">text</span><span class=\"p\">[</span><span class=\"s2\">&quot;value&quot;</span><span
    class=\"p\">]</span> <span class=\"k\">for</span> <span class=\"n\">text</span>
    <span class=\"ow\">in</span> <span class=\"n\">v</span><span class=\"p\">])</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='html' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>html
    <em class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">html <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">html</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">):</span>\n                <span class=\"n\">html</span> <span class=\"o\">=</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">text</span>\n
    \               <span class=\"n\">html</span> <span class=\"o\">+=</span> <span
    class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span>\n
    \               <span class=\"k\">for</span> <span class=\"n\">meta</span> <span
    class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">meta</span><span class=\"p\">:</span>\n                    <span class=\"n\">html</span>
    <span class=\"o\">+=</span> <span class=\"sa\">f</span><span class=\"s1\">&#39;&lt;meta
    name=&quot;</span><span class=\"si\">{</span><span class=\"n\">meta</span><span
    class=\"o\">.</span><span class=\"n\">name</span><span class=\"si\">}</span><span
    class=\"s1\">&quot; content=&quot;</span><span class=\"si\">{</span><span class=\"n\">meta</span><span
    class=\"o\">.</span><span class=\"n\">content</span><span class=\"si\">}</span><span
    class=\"s1\">&quot; /&gt;</span><span class=\"se\">\\n</span><span class=\"s1\">&#39;</span>\n
    \               <span class=\"k\">for</span> <span class=\"n\">link</span> <span
    class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">link</span><span class=\"p\">:</span>\n                    <span class=\"n\">html</span>
    <span class=\"o\">+=</span> <span class=\"sa\">f</span><span class=\"s1\">&#39;&lt;link
    rel=&quot;</span><span class=\"si\">{</span><span class=\"n\">link</span><span
    class=\"o\">.</span><span class=\"n\">rel</span><span class=\"si\">}</span><span
    class=\"s1\">&quot; href=&quot;</span><span class=\"si\">{</span><span class=\"n\">link</span><span
    class=\"o\">.</span><span class=\"n\">href</span><span class=\"si\">}</span><span
    class=\"s1\">&quot; /&gt;</span><span class=\"se\">\\n</span><span class=\"s1\">&#39;</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">html</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='dynamic_templates_in_templates_dir' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>dynamic_templates_in_templates_dir <em class='small'>method</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">dynamic_templates_in_templates_dir
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
    <span class=\"nf\">dynamic_templates_in_templates_dir</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">):</span>\n                <span class=\"n\">markata_templates</span>
    <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"vm\">__file__</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">parents</span><span class=\"p\">[</span><span class=\"mi\">1</span><span
    class=\"p\">]</span> <span class=\"o\">/</span> <span class=\"s2\">&quot;templates&quot;</span>\n\n
    \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">templates_dir</span><span class=\"p\">,</span> <span class=\"n\">Path</span><span
    class=\"p\">):</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">templates_dir</span> <span class=\"o\">=</span>
    <span class=\"p\">[</span>\n                        <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">templates_dir</span><span class=\"p\">,</span>\n
    \                       <span class=\"n\">markata_templates</span><span class=\"p\">,</span>\n
    \                       <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">dynamic_templates_dir</span><span class=\"p\">,</span>\n                    <span
    class=\"p\">]</span>\n\n                <span class=\"k\">if</span> <span class=\"n\">markata_templates</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">templates_dir</span><span class=\"p\">:</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">templates_dir</span><span class=\"o\">.</span><span class=\"n\">append</span><span
    class=\"p\">(</span><span class=\"n\">markata_templates</span><span class=\"p\">)</span>\n\n
    \               <span class=\"k\">if</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">dynamic_templates_dir</span> <span class=\"ow\">not</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">templates_dir</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">templates_dir</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">dynamic_templates_dir</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">return</span> <span
    class=\"bp\">self</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='jinja_loader'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>jinja_loader <em
    class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">jinja_loader <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">jinja_loader</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">):</span>\n                <span class=\"k\">return</span> <span class=\"n\">jinja2</span><span
    class=\"o\">.</span><span class=\"n\">FileSystemLoader</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">templates_dir</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='jinja_env'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>jinja_env <em class='small'>method</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">jinja_env
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
    <span class=\"nf\">jinja_env</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">self</span><span class=\"p\">,</span>\n            <span class=\"p\">):</span>\n
    \               <span class=\"k\">if</span> <span class=\"nb\">hasattr</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;_jinja_env&quot;</span><span class=\"p\">):</span>\n                    <span
    class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_jinja_env</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">env_options</span><span class=\"o\">.</span><span
    class=\"n\">setdefault</span><span class=\"p\">(</span><span class=\"s2\">&quot;loader&quot;</span><span
    class=\"p\">,</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">jinja_loader</span><span class=\"p\">)</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">env_options</span><span
    class=\"o\">.</span><span class=\"n\">setdefault</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;undefined&quot;</span><span class=\"p\">,</span> <span class=\"n\">SilentUndefined</span><span
    class=\"p\">)</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">env_options</span><span class=\"o\">.</span><span class=\"n\">setdefault</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;lstrip_blocks&quot;</span><span class=\"p\">,</span>
    <span class=\"kc\">True</span><span class=\"p\">)</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">env_options</span><span
    class=\"o\">.</span><span class=\"n\">setdefault</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;trim_blocks&quot;</span><span class=\"p\">,</span> <span class=\"kc\">True</span><span
    class=\"p\">)</span>\n\n                <span class=\"n\">env</span> <span class=\"o\">=</span>
    <span class=\"n\">jinja2</span><span class=\"o\">.</span><span class=\"n\">Environment</span><span
    class=\"p\">(</span><span class=\"o\">**</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">env_options</span><span class=\"p\">)</span>\n\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_jinja_env</span> <span class=\"o\">=</span> <span class=\"n\">env</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">env</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='default_template' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>default_template <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">default_template
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
    <span class=\"nf\">default_template</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
    class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;markata&quot;</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">post_template</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">v</span><span
    class=\"p\">,</span> <span class=\"nb\">str</span><span class=\"p\">):</span>\n
    \                   <span class=\"n\">v</span> <span class=\"o\">=</span> <span
    class=\"p\">{</span><span class=\"s2\">&quot;index&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">v</span><span class=\"p\">}</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;markata&quot;</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">post_template</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
    class=\"p\">):</span>\n                    <span class=\"n\">config_template</span>
    <span class=\"o\">=</span> <span class=\"p\">{</span>\n                        <span
    class=\"s2\">&quot;index&quot;</span><span class=\"p\">:</span> <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;markata&quot;</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">post_template</span><span class=\"p\">,</span>\n                    <span
    class=\"p\">}</span>\n                <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">config_template</span> <span class=\"o\">=</span>
    <span class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">&quot;markata&quot;</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">post_template</span>\n                <span
    class=\"k\">return</span> <span class=\"p\">{</span><span class=\"o\">**</span><span
    class=\"n\">config_template</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"n\">v</span><span class=\"p\">}</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='templates' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>templates <em class='small'>function</em></h2>\ntemplate management</p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">templates
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
    <span class=\"nf\">templates</span><span class=\"p\">():</span>\n                <span
    class=\"s2\">&quot;template management&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='show' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>show <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">show
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
    <span class=\"nf\">show</span><span class=\"p\">(</span>\n                <span
    class=\"n\">template</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
    class=\"n\">Argument</span><span class=\"p\">(</span><span class=\"kc\">None</span><span
    class=\"p\">,</span> <span class=\"n\">help</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;template to show&quot;</span><span class=\"p\">),</span>\n
    \               <span class=\"n\">theme</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"n\">typer</span><span
    class=\"o\">.</span><span class=\"n\">Option</span><span class=\"p\">(</span><span
    class=\"kc\">None</span><span class=\"p\">,</span> <span class=\"n\">help</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;pygments syntax theme&quot;</span><span
    class=\"p\">),</span>\n            <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">quiet</span> <span class=\"o\">=</span>
    <span class=\"kc\">True</span>\n                <span class=\"k\">if</span> <span
    class=\"n\">template</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">template</span> <span class=\"o\">=</span> <span class=\"n\">get_template</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">template</span><span class=\"p\">)</span>\n\n                    <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">quiet</span> <span class=\"o\">=</span>
    <span class=\"kc\">False</span>\n                    <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">print</span><span class=\"p\">(</span><span class=\"n\">template</span><span
    class=\"o\">.</span><span class=\"n\">filename</span><span class=\"p\">)</span>\n
    \                   <span class=\"k\">if</span> <span class=\"n\">theme</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span> <span class=\"ow\">or</span>
    <span class=\"n\">theme</span><span class=\"o\">.</span><span class=\"n\">lower</span><span
    class=\"p\">()</span> <span class=\"o\">==</span> <span class=\"s2\">&quot;none&quot;</span><span
    class=\"p\">:</span>\n                        <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">print</span><span class=\"p\">(</span><span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">template</span><span class=\"o\">.</span><span
    class=\"n\">filename</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">read_text</span><span class=\"p\">())</span>\n                    <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">syntax</span> <span class=\"o\">=</span> <span class=\"n\">Syntax</span><span
    class=\"o\">.</span><span class=\"n\">from_path</span><span class=\"p\">(</span><span
    class=\"n\">template</span><span class=\"o\">.</span><span class=\"n\">filename</span><span
    class=\"p\">,</span> <span class=\"n\">theme</span><span class=\"o\">=</span><span
    class=\"n\">theme</span><span class=\"p\">)</span>\n                        <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">print</span><span class=\"p\">(</span><span
    class=\"n\">syntax</span><span class=\"p\">)</span>\n\n                    <span
    class=\"k\">return</span>\n                <span class=\"n\">templates</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">jinja_env</span><span
    class=\"o\">.</span><span class=\"n\">list_templates</span><span class=\"p\">()</span>\n
    \               <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">quiet</span>
    <span class=\"o\">=</span> <span class=\"kc\">False</span>\n                <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">print</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;Templates directories:&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">style</span><span class=\"o\">=</span><span class=\"s2\">&quot;green
    underline&quot;</span><span class=\"p\">)</span>\n\n                <span class=\"n\">markata_templates</span>
    <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"vm\">__file__</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">parents</span><span class=\"p\">[</span><span class=\"mi\">1</span><span
    class=\"p\">]</span> <span class=\"o\">/</span> <span class=\"s2\">&quot;templates&quot;</span>\n
    \               <span class=\"k\">for</span> <span class=\"nb\">dir</span> <span
    class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">templates_dir</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">if</span> <span class=\"nb\">dir</span>
    <span class=\"o\">==</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">dynamic_templates_dir</span><span
    class=\"p\">:</span>\n                        <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">print</span><span class=\"p\">(</span>\n                            <span
    class=\"sa\">f</span><span class=\"s2\">&quot;[gold3]</span><span class=\"si\">{</span><span
    class=\"nb\">dir</span><span class=\"si\">}</span><span class=\"s2\">[/][grey50]
    (dynamically created templates from configuration)[/] [gold3]\\[markata.config.dynamic_templates_dir][/]&quot;</span><span
    class=\"p\">,</span>\n                            <span class=\"n\">style</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;red&quot;</span><span class=\"p\">,</span>\n
    \                       <span class=\"p\">)</span>\n                    <span
    class=\"k\">elif</span> <span class=\"nb\">dir</span> <span class=\"o\">==</span>
    <span class=\"n\">markata_templates</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">print</span><span class=\"p\">(</span>\n
    \                           <span class=\"sa\">f</span><span class=\"s2\">&quot;[cyan]</span><span
    class=\"si\">{</span><span class=\"nb\">dir</span><span class=\"si\">}</span><span
    class=\"s2\">[/][grey50] (built-in)[/]&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">style</span><span class=\"o\">=</span><span class=\"s2\">&quot;red&quot;</span>\n
    \                       <span class=\"p\">)</span>\n                    <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">print</span><span class=\"p\">(</span>\n
    \                           <span class=\"sa\">f</span><span class=\"s2\">&quot;[orchid]</span><span
    class=\"si\">{</span><span class=\"nb\">dir</span><span class=\"si\">}</span><span
    class=\"s2\">[/] [orchid]\\[markata.config.templates_dir][/]&quot;</span><span
    class=\"p\">,</span>\n                            <span class=\"n\">style</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;red&quot;</span><span class=\"p\">,</span>\n
    \                       <span class=\"p\">)</span>\n\n                <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">print</span><span class=\"p\">()</span>\n                <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">print</span><span class=\"p\">(</span>\n                    <span
    class=\"s2\">&quot;Available Templates: [white]name -&gt; path[/]&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">style</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;green underline&quot;</span>\n                <span class=\"p\">)</span>\n
    \               <span class=\"k\">for</span> <span class=\"n\">template</span>
    <span class=\"ow\">in</span> <span class=\"n\">templates</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">source</span><span class=\"p\">,</span>
    <span class=\"n\">file</span><span class=\"p\">,</span> <span class=\"n\">uptodate</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">jinja_env</span><span
    class=\"o\">.</span><span class=\"n\">loader</span><span class=\"o\">.</span><span
    class=\"n\">get_source</span><span class=\"p\">(</span>\n                        <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">jinja_env</span><span class=\"p\">,</span>
    <span class=\"n\">template</span>\n                    <span class=\"p\">)</span>\n\n
    \                   <span class=\"k\">if</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">file</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">is_relative_to</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">dynamic_templates_dir</span><span class=\"p\">):</span>\n
    \                       <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">print</span><span
    class=\"p\">(</span>\n                            <span class=\"sa\">f</span><span
    class=\"s2\">&quot;[gold3]</span><span class=\"si\">{</span><span class=\"n\">template</span><span
    class=\"si\">}</span><span class=\"s2\"> -&gt; [red]</span><span class=\"si\">{</span><span
    class=\"n\">file</span><span class=\"si\">}</span><span class=\"s2\">[/] [grey50](dynamic)[/]&quot;</span>\n
    \                       <span class=\"p\">)</span>\n                    <span
    class=\"k\">elif</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">file</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">is_relative_to</span><span class=\"p\">(</span><span class=\"n\">markata_templates</span><span
    class=\"p\">):</span>\n                        <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">print</span><span class=\"p\">(</span>\n                            <span
    class=\"sa\">f</span><span class=\"s2\">&quot;[cyan]</span><span class=\"si\">{</span><span
    class=\"n\">template</span><span class=\"si\">}</span><span class=\"s2\"> -&gt;
    [red]</span><span class=\"si\">{</span><span class=\"n\">file</span><span class=\"si\">}</span><span
    class=\"s2\">[/] [grey50](built-in)[/]&quot;</span>\n                        <span
    class=\"p\">)</span>\n                    <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">print</span><span
    class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">&quot;[orchid]</span><span
    class=\"si\">{</span><span class=\"n\">template</span><span class=\"si\">}</span><span
    class=\"s2\">[/] -&gt; [red]</span><span class=\"si\">{</span><span class=\"n\">file</span><span
    class=\"si\">}</span><span class=\"s2\">[/]&quot;</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>    </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Post_Template.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"This snippet allows users to configure
    their head in  Users can specify any sort of tag in their  The above configuration
    becomes this once rendered. ! Optional\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Post_Template.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"This snippet allows users to configure
    their head in  Users can specify any sort of tag in their  The above configuration
    becomes this once rendered. ! Optional\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<article style=\"text-align: center;\">\n    <style>\n        section
    {\n            font-size: 200%;\n        }\n\n\n        .edit {\n            display:
    none;\n        }\n    </style>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Post_Template.Py\n    </h1>\n</section></article>\n     </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Post_Template.Py\n    </h1>\n</section>    <section class=\"body\">\n
    \       <h1 id=\"add-head-configuration\">Add head configuration <a class=\"header-anchor\"
    href=\"#add-head-configuration\"><svg class=\"heading-permalink\" aria-hidden=\"true\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>This snippet allows
    users to configure their head in <code>markata.toml</code>.</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span>{{ config.get(&#39;head&#39;,
    {}).pop(&#39;text&#39;) if &#39;text&#39; in config.get(&#39;head&#39;,{}).keys()
    }}\n{% for tag, meta in config.get(&#39;head&#39;, {}).items() %}\n    {% for
    _meta in meta %}\n        <span class=\"err\">&lt;</span>{{ tag }}\n            {%
    for attr, value in _meta.items() %}{{ attr }}=&quot;{{ value }}&quot;{% endfor
    %}\n        /&gt;\n    {% endfor %}\n{% endfor %}\n</pre></div>\n\n</pre>\n\n<p>Users
    can specify any sort of tag in their <code>markata.toml</code></p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.head.meta]]</span>\n<span
    class=\"n\">name</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;og:type&quot;</span>\n<span class=\"n\">content</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;article&quot;</span>\n\n<span class=\"k\">[[markata.head.meta]]</span>\n<span
    class=\"n\">name</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;og:author&quot;</span>\n<span class=\"n\">content</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;Waylon Walker&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>The
    above configuration becomes this once rendered.</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"p\">&lt;</span><span
    class=\"nt\">meta</span> <span class=\"na\">name</span><span class=\"o\">=</span><span
    class=\"s\">&#39;og:type&#39;</span> <span class=\"na\">content</span><span class=\"o\">=</span><span
    class=\"s\">&#39;article&#39;</span> <span class=\"p\">/&gt;</span>\n<span class=\"p\">&lt;</span><span
    class=\"nt\">meta</span> <span class=\"na\">name</span><span class=\"o\">=</span><span
    class=\"s\">&#39;og:Author&#39;</span> <span class=\"na\">content</span><span
    class=\"o\">=</span><span class=\"s\">&#39;Waylon Walker&#39;</span> <span class=\"p\">/&gt;</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    Note</p>\n<pre><code>Article variables can be used for dynamic entries like canonical_url\n```
    toml\n[markata]\nurl = &quot;markata.dev&quot;\n\n[[markata.head.meta]]\nhref=&quot;{{
    config.url }}/{{ slug }}/&quot;\nrel=&quot;canonical&quot;\n```\n</code></pre>\n<p>Optionally
    users can also specify plain text to be appended to the head of\ntheir documents.
    \ This works well for things that involve full blocks.</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.head.text]]</span>\n<span
    class=\"n\">value</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s1\">&#39;&#39;&#39;</span>\n<span class=\"s1\">&lt;script&gt;</span>\n<span
    class=\"s1\">    console.log(&#39;hello world&#39;)</span>\n<span class=\"s1\">&lt;/script&gt;</span>\n<span
    class=\"s1\">&#39;&#39;&#39;</span>\n\n<span class=\"k\">[[markata.head.text]]</span>\n<span
    class=\"n\">value</span><span class=\"o\">=</span><span class=\"s1\">&#39;&#39;&#39;</span>\n<span
    class=\"s1\">html  {</span>\n<span class=\"s1\">    font-family: &quot;Space Mono&quot;,
    monospace;</span>\n<span class=\"s1\">    background: var(--color-bg);</span>\n<span
    class=\"s1\">    color: var(--color-text);</span>\n<span class=\"s1\">}</span>\n<span
    class=\"s1\">&#39;&#39;&#39;</span>\n</pre></div>\n\n</pre>\n\n<h2 id=\"add-scripts-to-head\">Add
    scripts to head <a class=\"header-anchor\" href=\"#add-scripts-to-head\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Markata config also
    supports adding scripts to the head via configuration.</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[</span><span
    class=\"w\"> </span><span class=\"k\">markata.head.script</span><span class=\"w\">
    </span><span class=\"k\">]]</span>\n<span class=\"w\">    </span><span class=\"n\">src</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;https://cdn.tailwindcss.com&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='SilentUndefined' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>SilentUndefined <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">SilentUndefined
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
    <span class=\"nc\">SilentUndefined</span><span class=\"p\">(</span><span class=\"n\">Undefined</span><span
    class=\"p\">):</span>\n            <span class=\"k\">def</span> <span class=\"nf\">_fail_with_undefined_error</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"n\">args</span><span class=\"p\">,</span> <span
    class=\"o\">**</span><span class=\"n\">kwargs</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">return</span> <span class=\"s2\">&quot;&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='optional' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>optional <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">optional
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
    <span class=\"nf\">optional</span><span class=\"p\">(</span><span class=\"o\">*</span><span
    class=\"n\">fields</span><span class=\"p\">):</span>\n            <span class=\"k\">def</span>
    <span class=\"nf\">dec</span><span class=\"p\">(</span><span class=\"n\">_cls</span><span
    class=\"p\">):</span>\n                <span class=\"k\">for</span> <span class=\"n\">field</span>
    <span class=\"ow\">in</span> <span class=\"n\">fields</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">_cls</span><span class=\"o\">.</span><span
    class=\"n\">__fields__</span><span class=\"p\">[</span><span class=\"n\">field</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">default</span>
    <span class=\"o\">=</span> <span class=\"kc\">None</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">_cls</span>\n\n            <span class=\"k\">if</span>
    <span class=\"p\">(</span>\n                <span class=\"n\">fields</span>\n
    \               <span class=\"ow\">and</span> <span class=\"n\">inspect</span><span
    class=\"o\">.</span><span class=\"n\">isclass</span><span class=\"p\">(</span><span
    class=\"n\">fields</span><span class=\"p\">[</span><span class=\"mi\">0</span><span
    class=\"p\">])</span>\n                <span class=\"ow\">and</span> <span class=\"nb\">issubclass</span><span
    class=\"p\">(</span><span class=\"n\">fields</span><span class=\"p\">[</span><span
    class=\"mi\">0</span><span class=\"p\">],</span> <span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">)</span>\n
    \           <span class=\"p\">):</span>\n                <span class=\"bp\">cls</span>
    <span class=\"o\">=</span> <span class=\"n\">fields</span><span class=\"p\">[</span><span
    class=\"mi\">0</span><span class=\"p\">]</span>\n                <span class=\"n\">fields</span>
    <span class=\"o\">=</span> <span class=\"bp\">cls</span><span class=\"o\">.</span><span
    class=\"n\">__fields__</span>\n                <span class=\"k\">return</span>
    <span class=\"n\">dec</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">)</span>\n            <span class=\"k\">return</span> <span class=\"n\">dec</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='Style' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Style
    <em class='small'>class</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">Style <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nc\">Style</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">color_bg</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;#1f2022&quot;</span>\n
    \           <span class=\"n\">color_bg_code</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;#1f2022&quot;</span>\n
    \           <span class=\"n\">color_text</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;#eefbfe&quot;</span>\n
    \           <span class=\"n\">color_link</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;#fb30c4&quot;</span>\n
    \           <span class=\"n\">color_accent</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;#e1bd00c9&quot;</span>\n
    \           <span class=\"n\">overlay_brightness</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;.85&quot;</span>\n
    \           <span class=\"n\">body_width</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;800px&quot;</span>\n
    \           <span class=\"n\">color_bg_light</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;#eefbfe&quot;</span>\n
    \           <span class=\"n\">color_bg_code_light</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;#eefbfe&quot;</span>\n
    \           <span class=\"n\">color_text_light</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;#1f2022&quot;</span>\n
    \           <span class=\"n\">color_link_light</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;#fb30c4&quot;</span>\n
    \           <span class=\"n\">color_accent_light</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;#ffeb00&quot;</span>\n
    \           <span class=\"n\">overlay_brightness_light</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;.95&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='StyleOverrides' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>StyleOverrides <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">StyleOverrides
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
    <span class=\"nc\">StyleOverrides</span><span class=\"p\">(</span><span class=\"n\">Style</span><span
    class=\"p\">):</span> <span class=\"o\">...</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='Meta' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Meta
    <em class='small'>class</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">Meta <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nc\">Meta</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">name</span><span class=\"p\">:</span> <span class=\"nb\">str</span>\n
    \           <span class=\"n\">content</span><span class=\"p\">:</span> <span class=\"nb\">str</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='Text' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Text
    <em class='small'>class</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">Text <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nc\">Text</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">value</span><span class=\"p\">:</span> <span class=\"nb\">str</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='Link' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Link
    <em class='small'>class</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">Link <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nc\">Link</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">rel</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;canonical&quot;</span>\n            <span
    class=\"n\">href</span><span class=\"p\">:</span> <span class=\"nb\">str</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='Script' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Script
    <em class='small'>class</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">Script <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nc\">Script</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">src</span><span class=\"p\">:</span> <span class=\"nb\">str</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='HeadConfig' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>HeadConfig <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">HeadConfig
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
    <span class=\"nc\">HeadConfig</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">meta</span><span class=\"p\">:</span> <span class=\"n\">List</span><span
    class=\"p\">[</span><span class=\"n\">Meta</span><span class=\"p\">]</span> <span
    class=\"o\">=</span> <span class=\"p\">[]</span>\n            <span class=\"n\">link</span><span
    class=\"p\">:</span> <span class=\"n\">List</span><span class=\"p\">[</span><span
    class=\"n\">Link</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"p\">[]</span>\n            <span class=\"n\">script</span><span class=\"p\">:</span>
    <span class=\"n\">List</span><span class=\"p\">[</span><span class=\"n\">Script</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"p\">[]</span>\n
    \           <span class=\"n\">text</span><span class=\"p\">:</span> <span class=\"n\">Union</span><span
    class=\"p\">[</span><span class=\"n\">List</span><span class=\"p\">[</span><span
    class=\"n\">Text</span><span class=\"p\">],</span> <span class=\"nb\">str</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;&quot;</span>\n\n
    \           <span class=\"nd\">@pydantic</span><span class=\"o\">.</span><span
    class=\"n\">validator</span><span class=\"p\">(</span><span class=\"s2\">&quot;text&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">pre</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n            <span class=\"k\">def</span>
    <span class=\"nf\">text_to_list</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">):</span>\n                <span
    class=\"k\">if</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"nb\">list</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span><span
    class=\"o\">.</span><span class=\"n\">join</span><span class=\"p\">([</span><span
    class=\"n\">text</span><span class=\"p\">[</span><span class=\"s2\">&quot;value&quot;</span><span
    class=\"p\">]</span> <span class=\"k\">for</span> <span class=\"n\">text</span>
    <span class=\"ow\">in</span> <span class=\"n\">v</span><span class=\"p\">])</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">v</span>\n\n
    \           <span class=\"nd\">@property</span>\n            <span class=\"k\">def</span>
    <span class=\"nf\">html</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">):</span>\n                <span class=\"n\">html</span> <span class=\"o\">=</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">text</span>\n
    \               <span class=\"n\">html</span> <span class=\"o\">+=</span> <span
    class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span>\n
    \               <span class=\"k\">for</span> <span class=\"n\">meta</span> <span
    class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">meta</span><span class=\"p\">:</span>\n                    <span class=\"n\">html</span>
    <span class=\"o\">+=</span> <span class=\"sa\">f</span><span class=\"s1\">&#39;&lt;meta
    name=&quot;</span><span class=\"si\">{</span><span class=\"n\">meta</span><span
    class=\"o\">.</span><span class=\"n\">name</span><span class=\"si\">}</span><span
    class=\"s1\">&quot; content=&quot;</span><span class=\"si\">{</span><span class=\"n\">meta</span><span
    class=\"o\">.</span><span class=\"n\">content</span><span class=\"si\">}</span><span
    class=\"s1\">&quot; /&gt;</span><span class=\"se\">\\n</span><span class=\"s1\">&#39;</span>\n
    \               <span class=\"k\">for</span> <span class=\"n\">link</span> <span
    class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">link</span><span class=\"p\">:</span>\n                    <span class=\"n\">html</span>
    <span class=\"o\">+=</span> <span class=\"sa\">f</span><span class=\"s1\">&#39;&lt;link
    rel=&quot;</span><span class=\"si\">{</span><span class=\"n\">link</span><span
    class=\"o\">.</span><span class=\"n\">rel</span><span class=\"si\">}</span><span
    class=\"s1\">&quot; href=&quot;</span><span class=\"si\">{</span><span class=\"n\">link</span><span
    class=\"o\">.</span><span class=\"n\">href</span><span class=\"si\">}</span><span
    class=\"s1\">&quot; /&gt;</span><span class=\"se\">\\n</span><span class=\"s1\">&#39;</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">html</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    \           <span class=\"n\">head</span><span class=\"p\">:</span> <span class=\"n\">HeadConfig</span>
    <span class=\"o\">=</span> <span class=\"n\">HeadConfig</span><span class=\"p\">()</span>\n
    \           <span class=\"n\">style</span><span class=\"p\">:</span> <span class=\"n\">Style</span>
    <span class=\"o\">=</span> <span class=\"n\">Style</span><span class=\"p\">()</span>\n
    \           <span class=\"n\">post_template</span><span class=\"p\">:</span> <span
    class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"n\">Union</span><span
    class=\"p\">[</span><span class=\"nb\">str</span> <span class=\"o\">|</span> <span
    class=\"n\">Dict</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"nb\">str</span><span class=\"p\">]]]</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;post.html&quot;</span>\n            <span
    class=\"n\">dynamic_templates_dir</span><span class=\"p\">:</span> <span class=\"n\">Path</span>
    <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;.markata.cache/templates&quot;</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">templates_dir</span><span class=\"p\">:</span> <span
    class=\"n\">Union</span><span class=\"p\">[</span><span class=\"n\">Path</span><span
    class=\"p\">,</span> <span class=\"n\">List</span><span class=\"p\">[</span><span
    class=\"n\">Path</span><span class=\"p\">]]</span> <span class=\"o\">=</span>
    <span class=\"n\">pydantic</span><span class=\"o\">.</span><span class=\"n\">Field</span><span
    class=\"p\">(</span><span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;templates&quot;</span><span class=\"p\">))</span>\n\n            <span
    class=\"n\">env_options</span><span class=\"p\">:</span> <span class=\"nb\">dict</span>
    <span class=\"o\">=</span> <span class=\"p\">{}</span>\n\n            <span class=\"nd\">@pydantic</span><span
    class=\"o\">.</span><span class=\"n\">model_validator</span><span class=\"p\">(</span><span
    class=\"n\">mode</span><span class=\"o\">=</span><span class=\"s2\">&quot;after&quot;</span><span
    class=\"p\">)</span>\n            <span class=\"k\">def</span> <span class=\"nf\">dynamic_templates_in_templates_dir</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">):</span>\n
    \               <span class=\"n\">markata_templates</span> <span class=\"o\">=</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"vm\">__file__</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">parents</span><span
    class=\"p\">[</span><span class=\"mi\">1</span><span class=\"p\">]</span> <span
    class=\"o\">/</span> <span class=\"s2\">&quot;templates&quot;</span>\n\n                <span
    class=\"k\">if</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">templates_dir</span><span
    class=\"p\">,</span> <span class=\"n\">Path</span><span class=\"p\">):</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">templates_dir</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
    \                       <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">templates_dir</span><span class=\"p\">,</span>\n                        <span
    class=\"n\">markata_templates</span><span class=\"p\">,</span>\n                        <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">dynamic_templates_dir</span><span
    class=\"p\">,</span>\n                    <span class=\"p\">]</span>\n\n                <span
    class=\"k\">if</span> <span class=\"n\">markata_templates</span> <span class=\"ow\">not</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">templates_dir</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">templates_dir</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">markata_templates</span><span class=\"p\">)</span>\n\n                <span
    class=\"k\">if</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">dynamic_templates_dir</span> <span class=\"ow\">not</span> <span class=\"ow\">in</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">templates_dir</span><span
    class=\"p\">:</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">templates_dir</span><span class=\"o\">.</span><span
    class=\"n\">append</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">dynamic_templates_dir</span><span class=\"p\">)</span>\n\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span>\n\n
    \           <span class=\"nd\">@property</span>\n            <span class=\"k\">def</span>
    <span class=\"nf\">jinja_loader</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">):</span>\n                <span class=\"k\">return</span> <span class=\"n\">jinja2</span><span
    class=\"o\">.</span><span class=\"n\">FileSystemLoader</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">templates_dir</span><span
    class=\"p\">)</span>\n\n            <span class=\"nd\">@property</span>\n            <span
    class=\"k\">def</span> <span class=\"nf\">jinja_env</span><span class=\"p\">(</span>\n
    \               <span class=\"bp\">self</span><span class=\"p\">,</span>\n            <span
    class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"nb\">hasattr</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;_jinja_env&quot;</span><span class=\"p\">):</span>\n                    <span
    class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_jinja_env</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">env_options</span><span class=\"o\">.</span><span
    class=\"n\">setdefault</span><span class=\"p\">(</span><span class=\"s2\">&quot;loader&quot;</span><span
    class=\"p\">,</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">jinja_loader</span><span class=\"p\">)</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">env_options</span><span
    class=\"o\">.</span><span class=\"n\">setdefault</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;undefined&quot;</span><span class=\"p\">,</span> <span class=\"n\">SilentUndefined</span><span
    class=\"p\">)</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">env_options</span><span class=\"o\">.</span><span class=\"n\">setdefault</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;lstrip_blocks&quot;</span><span class=\"p\">,</span>
    <span class=\"kc\">True</span><span class=\"p\">)</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">env_options</span><span
    class=\"o\">.</span><span class=\"n\">setdefault</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;trim_blocks&quot;</span><span class=\"p\">,</span> <span class=\"kc\">True</span><span
    class=\"p\">)</span>\n\n                <span class=\"n\">env</span> <span class=\"o\">=</span>
    <span class=\"n\">jinja2</span><span class=\"o\">.</span><span class=\"n\">Environment</span><span
    class=\"p\">(</span><span class=\"o\">**</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">env_options</span><span class=\"p\">)</span>\n\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_jinja_env</span> <span class=\"o\">=</span> <span class=\"n\">env</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">env</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='PostOverrides' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>PostOverrides <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">PostOverrides
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
    <span class=\"nc\">PostOverrides</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">head</span><span class=\"p\">:</span> <span class=\"n\">HeadConfig</span>
    <span class=\"o\">=</span> <span class=\"n\">HeadConfig</span><span class=\"p\">()</span>\n
    \           <span class=\"n\">style</span><span class=\"p\">:</span> <span class=\"n\">Style</span>
    <span class=\"o\">=</span> <span class=\"n\">StyleOverrides</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">config_overrides</span><span class=\"p\">:</span>
    <span class=\"n\">PostOverrides</span> <span class=\"o\">=</span> <span class=\"n\">PostOverrides</span><span
    class=\"p\">()</span>\n            <span class=\"n\">template</span><span class=\"p\">:</span>
    <span class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"nb\">str</span>
    <span class=\"o\">|</span> <span class=\"n\">Dict</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
    class=\"p\">]]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n\n
    \           <span class=\"nd\">@pydantic</span><span class=\"o\">.</span><span
    class=\"n\">validator</span><span class=\"p\">(</span><span class=\"s2\">&quot;template&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">pre</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span> <span class=\"n\">always</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">def</span> <span class=\"nf\">default_template</span><span
    class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"p\">,</span> <span class=\"n\">values</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;markata&quot;</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">post_template</span>\n
    \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">v</span> <span class=\"o\">=</span> <span class=\"p\">{</span><span
    class=\"s2\">&quot;index&quot;</span><span class=\"p\">:</span> <span class=\"n\">v</span><span
    class=\"p\">}</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;markata&quot;</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">post_template</span><span
    class=\"p\">,</span> <span class=\"nb\">str</span><span class=\"p\">):</span>\n
    \                   <span class=\"n\">config_template</span> <span class=\"o\">=</span>
    <span class=\"p\">{</span>\n                        <span class=\"s2\">&quot;index&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;markata&quot;</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">post_template</span><span
    class=\"p\">,</span>\n                    <span class=\"p\">}</span>\n                <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                    <span class=\"n\">config_template</span>
    <span class=\"o\">=</span> <span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;markata&quot;</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">post_template</span>\n
    \               <span class=\"k\">return</span> <span class=\"p\">{</span><span
    class=\"o\">**</span><span class=\"n\">config_template</span><span class=\"p\">,</span>
    <span class=\"o\">**</span><span class=\"n\">v</span><span class=\"p\">}</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    function <h2 id='post_model' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>post_model <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">post_model
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
    function <h2 id='configure' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>configure <em class='small'>function</em></h2>\nMassages the configuration
    limitations of toml to allow a little bit easier\nexperience to the end user making
    configurations while allowing an simpler\njinja template.  This enablees the use
    of the <code>markata.head.text</code> list in\nconfiguration.</p>\n<div class=\"admonition
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
    <span class=\"nf\">configure</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            Massages the configuration limitations of toml to allow
    a little bit easier</span>\n<span class=\"sd\">            experience to the end
    user making configurations while allowing an simpler</span>\n<span class=\"sd\">
    \           jinja template.  This enablees the use of the `markata.head.text`
    list in</span>\n<span class=\"sd\">            configuration.</span>\n<span class=\"sd\">
    \           &quot;&quot;&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!! function
    <h2 id='pre_render' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>pre_render
    <em class='small'>function</em></h2>\nFOR EACH POST: Massages the configuration
    limitations of toml/yaml to allow\na little bit easier experience to the end user
    making configurations while\nallowing an simpler jinja template.  This enables
    the use of the\n<code>markata.head.text</code> list in configuration.</p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">pre_render
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
    <span class=\"nf\">pre_render</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            FOR EACH POST: Massages the configuration limitations
    of toml/yaml to allow</span>\n<span class=\"sd\">            a little bit easier
    experience to the end user making configurations while</span>\n<span class=\"sd\">
    \           allowing an simpler jinja template.  This enables the use of the</span>\n<span
    class=\"sd\">            `markata.head.text` list in configuration.</span>\n<span
    class=\"sd\">            &quot;&quot;&quot;</span>\n\n            <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">dynamic_templates_dir</span><span class=\"o\">.</span><span class=\"n\">mkdir</span><span
    class=\"p\">(</span><span class=\"n\">parents</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span> <span class=\"n\">exist_ok</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">head_template</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">dynamic_templates_dir</span> <span class=\"o\">/</span>
    <span class=\"s2\">&quot;head.html&quot;</span>\n            <span class=\"n\">head_template</span><span
    class=\"o\">.</span><span class=\"n\">write_text</span><span class=\"p\">(</span>\n
    \               <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">jinja_env</span><span
    class=\"o\">.</span><span class=\"n\">get_template</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;dynamic_head.html&quot;</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">render</span><span class=\"p\">(</span>\n
    \                   <span class=\"p\">{</span><span class=\"s2\">&quot;markata&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">markata</span><span class=\"p\">}</span>\n
    \               <span class=\"p\">),</span>\n            <span class=\"p\">)</span>\n\n
    \           <span class=\"k\">for</span> <span class=\"n\">article</span> <span
    class=\"ow\">in</span> <span class=\"p\">[</span><span class=\"n\">a</span> <span
    class=\"k\">for</span> <span class=\"n\">a</span> <span class=\"ow\">in</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">articles</span>
    <span class=\"k\">if</span> <span class=\"s2\">&quot;config_overrides&quot;</span>
    <span class=\"ow\">in</span> <span class=\"n\">a</span><span class=\"p\">]:</span>\n
    \               <span class=\"n\">raw_text</span> <span class=\"o\">=</span> <span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;config_overrides&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">{})</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;head&quot;</span><span
    class=\"p\">,</span> <span class=\"p\">{})</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;text&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;&quot;</span><span class=\"p\">)</span>\n\n
    \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">raw_text</span><span class=\"p\">,</span>
    <span class=\"nb\">list</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">article</span><span class=\"p\">[</span><span class=\"s2\">&quot;config_overrides&quot;</span><span
    class=\"p\">][</span><span class=\"s2\">&quot;head&quot;</span><span class=\"p\">][</span><span
    class=\"s2\">&quot;text&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span><span
    class=\"o\">.</span><span class=\"n\">join</span><span class=\"p\">(</span>\n
    \                       <span class=\"n\">flatten</span><span class=\"p\">([</span><span
    class=\"n\">t</span><span class=\"o\">.</span><span class=\"n\">values</span><span
    class=\"p\">()</span> <span class=\"k\">for</span> <span class=\"n\">t</span>
    <span class=\"ow\">in</span> <span class=\"n\">raw_text</span><span class=\"p\">]),</span>\n
    \                   <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='render' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>render <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">render
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
    <span class=\"nf\">render</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">with</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">cache</span> <span class=\"k\">as</span>
    <span class=\"n\">cache</span><span class=\"p\">:</span>\n                <span
    class=\"k\">for</span> <span class=\"n\">article</span> <span class=\"ow\">in</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">articles</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">html</span> <span
    class=\"o\">=</span> <span class=\"n\">render_article</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"o\">=</span><span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"n\">cache</span><span class=\"o\">=</span><span
    class=\"n\">cache</span><span class=\"p\">,</span> <span class=\"n\">article</span><span
    class=\"o\">=</span><span class=\"n\">article</span><span class=\"p\">)</span>\n
    \                   <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">html</span> <span class=\"o\">=</span> <span class=\"n\">html</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='get_template' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>get_template <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">get_template
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
    <span class=\"nf\">get_template</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"n\">template</span><span class=\"p\">):</span>\n
    \           <span class=\"k\">try</span><span class=\"p\">:</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">jinja_env</span><span
    class=\"o\">.</span><span class=\"n\">get_template</span><span class=\"p\">(</span><span
    class=\"n\">template</span><span class=\"p\">)</span>\n            <span class=\"k\">except</span>
    <span class=\"n\">jinja2</span><span class=\"o\">.</span><span class=\"n\">TemplateNotFound</span><span
    class=\"p\">:</span>\n                <span class=\"c1\"># try to load it as a
    file</span>\n                <span class=\"o\">...</span>\n\n            <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                <span class=\"k\">return</span>
    <span class=\"n\">Template</span><span class=\"p\">(</span><span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">template</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">read_text</span><span class=\"p\">(),</span>
    <span class=\"n\">undefined</span><span class=\"o\">=</span><span class=\"n\">SilentUndefined</span><span
    class=\"p\">)</span>\n            <span class=\"k\">except</span> <span class=\"ne\">FileNotFoundError</span><span
    class=\"p\">:</span>\n                <span class=\"c1\"># default to load it
    as a string</span>\n                <span class=\"o\">...</span>\n            <span
    class=\"k\">return</span> <span class=\"n\">Template</span><span class=\"p\">(</span><span
    class=\"n\">template</span><span class=\"p\">,</span> <span class=\"n\">undefined</span><span
    class=\"o\">=</span><span class=\"n\">SilentUndefined</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='render_article' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>render_article <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">render_article
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
    <span class=\"nf\">render_article</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"n\">cache</span><span class=\"p\">,</span>
    <span class=\"n\">article</span><span class=\"p\">):</span>\n            <span
    class=\"n\">key</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">make_hash</span><span class=\"p\">(</span>\n
    \               <span class=\"s2\">&quot;post_template&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">__version__</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">key</span><span class=\"p\">,</span>\n            <span class=\"p\">)</span>\n
    \           <span class=\"n\">html</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">precache</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"n\">key</span><span
    class=\"p\">)</span>\n\n            <span class=\"k\">if</span> <span class=\"n\">html</span>
    <span class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"n\">html</span>\n\n
    \           <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">template</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
    class=\"p\">):</span>\n                <span class=\"n\">template</span> <span
    class=\"o\">=</span> <span class=\"n\">get_template</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">,</span> <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">template</span><span class=\"p\">)</span>\n
    \               <span class=\"n\">html</span> <span class=\"o\">=</span> <span
    class=\"n\">render_template</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"n\">article</span><span class=\"p\">,</span>
    <span class=\"n\">template</span><span class=\"p\">)</span>\n\n            <span
    class=\"k\">if</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">template</span><span
    class=\"p\">,</span> <span class=\"nb\">dict</span><span class=\"p\">):</span>\n
    \               <span class=\"n\">html</span> <span class=\"o\">=</span> <span
    class=\"p\">{</span>\n                    <span class=\"n\">slug</span><span class=\"p\">:</span>
    <span class=\"n\">render_template</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"n\">article</span><span class=\"p\">,</span>
    <span class=\"n\">get_template</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"n\">template</span><span class=\"p\">))</span>\n
    \                   <span class=\"k\">for</span> <span class=\"n\">slug</span><span
    class=\"p\">,</span> <span class=\"n\">template</span> <span class=\"ow\">in</span>
    <span class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">template</span><span
    class=\"o\">.</span><span class=\"n\">items</span><span class=\"p\">()</span>\n
    \               <span class=\"p\">}</span>\n            <span class=\"n\">cache</span><span
    class=\"o\">.</span><span class=\"n\">add</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">,</span> <span class=\"n\">html</span><span
    class=\"p\">,</span> <span class=\"n\">expire</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">default_cache_expire</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">return</span> <span class=\"n\">html</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='render_template' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>render_template <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">render_template
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
    <span class=\"nf\">render_template</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"n\">article</span><span class=\"p\">,</span>
    <span class=\"n\">template</span><span class=\"p\">):</span>\n            <span
    class=\"n\">template</span> <span class=\"o\">=</span> <span class=\"n\">get_template</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">template</span><span class=\"p\">)</span>\n            <span
    class=\"n\">merged_config</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span>\n            <span class=\"c1\">#
    TODO do we need to handle merge??</span>\n            <span class=\"c1\"># if
    head_template:</span>\n            <span class=\"c1\">#     head = eval(</span>\n
    \           <span class=\"c1\">#         head_template.render(</span>\n            <span
    class=\"c1\">#             __version__=__version__,</span>\n            <span
    class=\"c1\">#             config=_full_config,</span>\n            <span class=\"c1\">#
    \            **article,</span>\n            <span class=\"c1\">#         )</span>\n
    \           <span class=\"c1\">#     )</span>\n\n            <span class=\"c1\">#
    merged_config = {</span>\n            <span class=\"c1\">#     **_full_config,</span>\n
    \           <span class=\"c1\">#     **{&quot;head&quot;: head},</span>\n            <span
    class=\"c1\"># }</span>\n\n            <span class=\"c1\"># merged_config = always_merger.merge(</span>\n
    \           <span class=\"c1\">#     merged_config,</span>\n            <span
    class=\"c1\">#     copy.deepcopy(</span>\n            <span class=\"c1\">#         article.get(</span>\n
    \           <span class=\"c1\">#             &quot;config_overrides&quot;,</span>\n
    \           <span class=\"c1\">#             {},</span>\n            <span class=\"c1\">#
    \        )</span>\n            <span class=\"c1\">#     ),</span>\n            <span
    class=\"c1\"># )</span>\n\n            <span class=\"n\">html</span> <span class=\"o\">=</span>
    <span class=\"n\">template</span><span class=\"o\">.</span><span class=\"n\">render</span><span
    class=\"p\">(</span>\n                <span class=\"n\">__version__</span><span
    class=\"o\">=</span><span class=\"n\">__version__</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"p\">,</span>\n                <span class=\"n\">body</span><span
    class=\"o\">=</span><span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">article_html</span><span class=\"p\">,</span>\n                <span
    class=\"n\">config</span><span class=\"o\">=</span><span class=\"n\">merged_config</span><span
    class=\"p\">,</span>\n                <span class=\"n\">post</span><span class=\"o\">=</span><span
    class=\"n\">article</span><span class=\"p\">,</span>\n            <span class=\"p\">)</span>\n
    \           <span class=\"k\">return</span> <span class=\"n\">html</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='save' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>save <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">save
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
    <span class=\"nf\">save</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">linked_templates</span> <span class=\"o\">=</span>
    <span class=\"p\">[</span>\n                <span class=\"n\">t</span>\n                <span
    class=\"k\">for</span> <span class=\"n\">t</span> <span class=\"ow\">in</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">jinja_env</span><span class=\"o\">.</span><span
    class=\"n\">list_templates</span><span class=\"p\">()</span>\n                <span
    class=\"k\">if</span> <span class=\"n\">t</span><span class=\"o\">.</span><span
    class=\"n\">endswith</span><span class=\"p\">(</span><span class=\"s2\">&quot;css&quot;</span><span
    class=\"p\">)</span> <span class=\"ow\">or</span> <span class=\"n\">t</span><span
    class=\"o\">.</span><span class=\"n\">endswith</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;js&quot;</span><span class=\"p\">)</span> <span class=\"ow\">or</span>
    <span class=\"n\">t</span><span class=\"o\">.</span><span class=\"n\">endswith</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;xsl&quot;</span><span class=\"p\">)</span>\n
    \           <span class=\"p\">]</span>\n            <span class=\"k\">for</span>
    <span class=\"n\">template</span> <span class=\"ow\">in</span> <span class=\"n\">linked_templates</span><span
    class=\"p\">:</span>\n                <span class=\"n\">template</span> <span
    class=\"o\">=</span> <span class=\"n\">get_template</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">,</span> <span class=\"n\">template</span><span
    class=\"p\">)</span>\n                <span class=\"n\">css</span> <span class=\"o\">=</span>
    <span class=\"n\">template</span><span class=\"o\">.</span><span class=\"n\">render</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"p\">,</span> <span class=\"n\">__version__</span><span
    class=\"o\">=</span><span class=\"n\">__version__</span><span class=\"p\">)</span>\n
    \               <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">output_dir</span> <span class=\"o\">/</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">template</span><span class=\"o\">.</span><span
    class=\"n\">filename</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">name</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">write_text</span><span class=\"p\">(</span><span class=\"n\">css</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='cli' class='admonition-title'
    style='margin:0;padding:.5rem 1rem;'>cli <em class='small'>function</em></h2>\nMarkata
    hook to implement base cli commands.</p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">cli <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">cli</span><span class=\"p\">(</span><span class=\"n\">app</span><span
    class=\"p\">:</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
    class=\"n\">Typer</span><span class=\"p\">,</span> <span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            Markata hook to implement base cli commands.</span>\n<span
    class=\"sd\">            &quot;&quot;&quot;</span>\n\n            <span class=\"n\">templates_app</span>
    <span class=\"o\">=</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
    class=\"n\">Typer</span><span class=\"p\">()</span>\n            <span class=\"n\">app</span><span
    class=\"o\">.</span><span class=\"n\">add_typer</span><span class=\"p\">(</span><span
    class=\"n\">templates_app</span><span class=\"p\">)</span>\n\n            <span
    class=\"nd\">@templates_app</span><span class=\"o\">.</span><span class=\"n\">callback</span><span
    class=\"p\">()</span>\n            <span class=\"k\">def</span> <span class=\"nf\">templates</span><span
    class=\"p\">():</span>\n                <span class=\"s2\">&quot;template management&quot;</span>\n\n
    \           <span class=\"nd\">@templates_app</span><span class=\"o\">.</span><span
    class=\"n\">command</span><span class=\"p\">()</span>\n            <span class=\"k\">def</span>
    <span class=\"nf\">show</span><span class=\"p\">(</span>\n                <span
    class=\"n\">template</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
    class=\"n\">Argument</span><span class=\"p\">(</span><span class=\"kc\">None</span><span
    class=\"p\">,</span> <span class=\"n\">help</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;template to show&quot;</span><span class=\"p\">),</span>\n
    \               <span class=\"n\">theme</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"n\">typer</span><span
    class=\"o\">.</span><span class=\"n\">Option</span><span class=\"p\">(</span><span
    class=\"kc\">None</span><span class=\"p\">,</span> <span class=\"n\">help</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;pygments syntax theme&quot;</span><span
    class=\"p\">),</span>\n            <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">quiet</span> <span class=\"o\">=</span>
    <span class=\"kc\">True</span>\n                <span class=\"k\">if</span> <span
    class=\"n\">template</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">template</span> <span class=\"o\">=</span> <span class=\"n\">get_template</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">template</span><span class=\"p\">)</span>\n\n                    <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">quiet</span> <span class=\"o\">=</span>
    <span class=\"kc\">False</span>\n                    <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">print</span><span class=\"p\">(</span><span class=\"n\">template</span><span
    class=\"o\">.</span><span class=\"n\">filename</span><span class=\"p\">)</span>\n
    \                   <span class=\"k\">if</span> <span class=\"n\">theme</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span> <span class=\"ow\">or</span>
    <span class=\"n\">theme</span><span class=\"o\">.</span><span class=\"n\">lower</span><span
    class=\"p\">()</span> <span class=\"o\">==</span> <span class=\"s2\">&quot;none&quot;</span><span
    class=\"p\">:</span>\n                        <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">print</span><span class=\"p\">(</span><span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">template</span><span class=\"o\">.</span><span
    class=\"n\">filename</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">read_text</span><span class=\"p\">())</span>\n                    <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">syntax</span> <span class=\"o\">=</span> <span class=\"n\">Syntax</span><span
    class=\"o\">.</span><span class=\"n\">from_path</span><span class=\"p\">(</span><span
    class=\"n\">template</span><span class=\"o\">.</span><span class=\"n\">filename</span><span
    class=\"p\">,</span> <span class=\"n\">theme</span><span class=\"o\">=</span><span
    class=\"n\">theme</span><span class=\"p\">)</span>\n                        <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">print</span><span class=\"p\">(</span><span
    class=\"n\">syntax</span><span class=\"p\">)</span>\n\n                    <span
    class=\"k\">return</span>\n                <span class=\"n\">templates</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">jinja_env</span><span
    class=\"o\">.</span><span class=\"n\">list_templates</span><span class=\"p\">()</span>\n
    \               <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">quiet</span>
    <span class=\"o\">=</span> <span class=\"kc\">False</span>\n                <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">print</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;Templates directories:&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">style</span><span class=\"o\">=</span><span class=\"s2\">&quot;green
    underline&quot;</span><span class=\"p\">)</span>\n\n                <span class=\"n\">markata_templates</span>
    <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"vm\">__file__</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">parents</span><span class=\"p\">[</span><span class=\"mi\">1</span><span
    class=\"p\">]</span> <span class=\"o\">/</span> <span class=\"s2\">&quot;templates&quot;</span>\n
    \               <span class=\"k\">for</span> <span class=\"nb\">dir</span> <span
    class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">templates_dir</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">if</span> <span class=\"nb\">dir</span>
    <span class=\"o\">==</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">dynamic_templates_dir</span><span
    class=\"p\">:</span>\n                        <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">print</span><span class=\"p\">(</span>\n                            <span
    class=\"sa\">f</span><span class=\"s2\">&quot;[gold3]</span><span class=\"si\">{</span><span
    class=\"nb\">dir</span><span class=\"si\">}</span><span class=\"s2\">[/][grey50]
    (dynamically created templates from configuration)[/] [gold3]\\[markata.config.dynamic_templates_dir][/]&quot;</span><span
    class=\"p\">,</span>\n                            <span class=\"n\">style</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;red&quot;</span><span class=\"p\">,</span>\n
    \                       <span class=\"p\">)</span>\n                    <span
    class=\"k\">elif</span> <span class=\"nb\">dir</span> <span class=\"o\">==</span>
    <span class=\"n\">markata_templates</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">print</span><span class=\"p\">(</span>\n
    \                           <span class=\"sa\">f</span><span class=\"s2\">&quot;[cyan]</span><span
    class=\"si\">{</span><span class=\"nb\">dir</span><span class=\"si\">}</span><span
    class=\"s2\">[/][grey50] (built-in)[/]&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">style</span><span class=\"o\">=</span><span class=\"s2\">&quot;red&quot;</span>\n
    \                       <span class=\"p\">)</span>\n                    <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">print</span><span class=\"p\">(</span>\n
    \                           <span class=\"sa\">f</span><span class=\"s2\">&quot;[orchid]</span><span
    class=\"si\">{</span><span class=\"nb\">dir</span><span class=\"si\">}</span><span
    class=\"s2\">[/] [orchid]\\[markata.config.templates_dir][/]&quot;</span><span
    class=\"p\">,</span>\n                            <span class=\"n\">style</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;red&quot;</span><span class=\"p\">,</span>\n
    \                       <span class=\"p\">)</span>\n\n                <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">print</span><span class=\"p\">()</span>\n                <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">print</span><span class=\"p\">(</span>\n                    <span
    class=\"s2\">&quot;Available Templates: [white]name -&gt; path[/]&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">style</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;green underline&quot;</span>\n                <span class=\"p\">)</span>\n
    \               <span class=\"k\">for</span> <span class=\"n\">template</span>
    <span class=\"ow\">in</span> <span class=\"n\">templates</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">source</span><span class=\"p\">,</span>
    <span class=\"n\">file</span><span class=\"p\">,</span> <span class=\"n\">uptodate</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">jinja_env</span><span
    class=\"o\">.</span><span class=\"n\">loader</span><span class=\"o\">.</span><span
    class=\"n\">get_source</span><span class=\"p\">(</span>\n                        <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">jinja_env</span><span class=\"p\">,</span>
    <span class=\"n\">template</span>\n                    <span class=\"p\">)</span>\n\n
    \                   <span class=\"k\">if</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">file</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">is_relative_to</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">dynamic_templates_dir</span><span class=\"p\">):</span>\n
    \                       <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">print</span><span
    class=\"p\">(</span>\n                            <span class=\"sa\">f</span><span
    class=\"s2\">&quot;[gold3]</span><span class=\"si\">{</span><span class=\"n\">template</span><span
    class=\"si\">}</span><span class=\"s2\"> -&gt; [red]</span><span class=\"si\">{</span><span
    class=\"n\">file</span><span class=\"si\">}</span><span class=\"s2\">[/] [grey50](dynamic)[/]&quot;</span>\n
    \                       <span class=\"p\">)</span>\n                    <span
    class=\"k\">elif</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">file</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">is_relative_to</span><span class=\"p\">(</span><span class=\"n\">markata_templates</span><span
    class=\"p\">):</span>\n                        <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">print</span><span class=\"p\">(</span>\n                            <span
    class=\"sa\">f</span><span class=\"s2\">&quot;[cyan]</span><span class=\"si\">{</span><span
    class=\"n\">template</span><span class=\"si\">}</span><span class=\"s2\"> -&gt;
    [red]</span><span class=\"si\">{</span><span class=\"n\">file</span><span class=\"si\">}</span><span
    class=\"s2\">[/] [grey50](built-in)[/]&quot;</span>\n                        <span
    class=\"p\">)</span>\n                    <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">print</span><span
    class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">&quot;[orchid]</span><span
    class=\"si\">{</span><span class=\"n\">template</span><span class=\"si\">}</span><span
    class=\"s2\">[/] -&gt; [red]</span><span class=\"si\">{</span><span class=\"n\">file</span><span
    class=\"si\">}</span><span class=\"s2\">[/]&quot;</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='_fail_with_undefined_error' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>_fail_with_undefined_error <em class='small'>method</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_fail_with_undefined_error
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
    <span class=\"nf\">_fail_with_undefined_error</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"n\">args</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"n\">kwargs</span><span class=\"p\">):</span>\n                <span class=\"k\">return</span>
    <span class=\"s2\">&quot;&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!! function
    <h2 id='dec' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>dec
    <em class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">dec <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">dec</span><span class=\"p\">(</span><span class=\"n\">_cls</span><span
    class=\"p\">):</span>\n                <span class=\"k\">for</span> <span class=\"n\">field</span>
    <span class=\"ow\">in</span> <span class=\"n\">fields</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">_cls</span><span class=\"o\">.</span><span
    class=\"n\">__fields__</span><span class=\"p\">[</span><span class=\"n\">field</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">default</span>
    <span class=\"o\">=</span> <span class=\"kc\">None</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">_cls</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='text_to_list' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>text_to_list <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">text_to_list
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
    <span class=\"nf\">text_to_list</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">):</span>\n                <span
    class=\"k\">if</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"nb\">list</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
    class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span><span
    class=\"o\">.</span><span class=\"n\">join</span><span class=\"p\">([</span><span
    class=\"n\">text</span><span class=\"p\">[</span><span class=\"s2\">&quot;value&quot;</span><span
    class=\"p\">]</span> <span class=\"k\">for</span> <span class=\"n\">text</span>
    <span class=\"ow\">in</span> <span class=\"n\">v</span><span class=\"p\">])</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='html' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>html
    <em class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">html <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">html</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">):</span>\n                <span class=\"n\">html</span> <span class=\"o\">=</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">text</span>\n
    \               <span class=\"n\">html</span> <span class=\"o\">+=</span> <span
    class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span>\n
    \               <span class=\"k\">for</span> <span class=\"n\">meta</span> <span
    class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">meta</span><span class=\"p\">:</span>\n                    <span class=\"n\">html</span>
    <span class=\"o\">+=</span> <span class=\"sa\">f</span><span class=\"s1\">&#39;&lt;meta
    name=&quot;</span><span class=\"si\">{</span><span class=\"n\">meta</span><span
    class=\"o\">.</span><span class=\"n\">name</span><span class=\"si\">}</span><span
    class=\"s1\">&quot; content=&quot;</span><span class=\"si\">{</span><span class=\"n\">meta</span><span
    class=\"o\">.</span><span class=\"n\">content</span><span class=\"si\">}</span><span
    class=\"s1\">&quot; /&gt;</span><span class=\"se\">\\n</span><span class=\"s1\">&#39;</span>\n
    \               <span class=\"k\">for</span> <span class=\"n\">link</span> <span
    class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">link</span><span class=\"p\">:</span>\n                    <span class=\"n\">html</span>
    <span class=\"o\">+=</span> <span class=\"sa\">f</span><span class=\"s1\">&#39;&lt;link
    rel=&quot;</span><span class=\"si\">{</span><span class=\"n\">link</span><span
    class=\"o\">.</span><span class=\"n\">rel</span><span class=\"si\">}</span><span
    class=\"s1\">&quot; href=&quot;</span><span class=\"si\">{</span><span class=\"n\">link</span><span
    class=\"o\">.</span><span class=\"n\">href</span><span class=\"si\">}</span><span
    class=\"s1\">&quot; /&gt;</span><span class=\"se\">\\n</span><span class=\"s1\">&#39;</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">html</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='dynamic_templates_in_templates_dir' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>dynamic_templates_in_templates_dir <em class='small'>method</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">dynamic_templates_in_templates_dir
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
    <span class=\"nf\">dynamic_templates_in_templates_dir</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">):</span>\n                <span class=\"n\">markata_templates</span>
    <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"vm\">__file__</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">parents</span><span class=\"p\">[</span><span class=\"mi\">1</span><span
    class=\"p\">]</span> <span class=\"o\">/</span> <span class=\"s2\">&quot;templates&quot;</span>\n\n
    \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">templates_dir</span><span class=\"p\">,</span> <span class=\"n\">Path</span><span
    class=\"p\">):</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">templates_dir</span> <span class=\"o\">=</span>
    <span class=\"p\">[</span>\n                        <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">templates_dir</span><span class=\"p\">,</span>\n
    \                       <span class=\"n\">markata_templates</span><span class=\"p\">,</span>\n
    \                       <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">dynamic_templates_dir</span><span class=\"p\">,</span>\n                    <span
    class=\"p\">]</span>\n\n                <span class=\"k\">if</span> <span class=\"n\">markata_templates</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">templates_dir</span><span class=\"p\">:</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">templates_dir</span><span class=\"o\">.</span><span class=\"n\">append</span><span
    class=\"p\">(</span><span class=\"n\">markata_templates</span><span class=\"p\">)</span>\n\n
    \               <span class=\"k\">if</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">dynamic_templates_dir</span> <span class=\"ow\">not</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">templates_dir</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">templates_dir</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">dynamic_templates_dir</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">return</span> <span
    class=\"bp\">self</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='jinja_loader'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>jinja_loader <em
    class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">jinja_loader <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">jinja_loader</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">):</span>\n                <span class=\"k\">return</span> <span class=\"n\">jinja2</span><span
    class=\"o\">.</span><span class=\"n\">FileSystemLoader</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">templates_dir</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='jinja_env'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>jinja_env <em class='small'>method</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">jinja_env
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
    <span class=\"nf\">jinja_env</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">self</span><span class=\"p\">,</span>\n            <span class=\"p\">):</span>\n
    \               <span class=\"k\">if</span> <span class=\"nb\">hasattr</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;_jinja_env&quot;</span><span class=\"p\">):</span>\n                    <span
    class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_jinja_env</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">env_options</span><span class=\"o\">.</span><span
    class=\"n\">setdefault</span><span class=\"p\">(</span><span class=\"s2\">&quot;loader&quot;</span><span
    class=\"p\">,</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">jinja_loader</span><span class=\"p\">)</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">env_options</span><span
    class=\"o\">.</span><span class=\"n\">setdefault</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;undefined&quot;</span><span class=\"p\">,</span> <span class=\"n\">SilentUndefined</span><span
    class=\"p\">)</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">env_options</span><span class=\"o\">.</span><span class=\"n\">setdefault</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;lstrip_blocks&quot;</span><span class=\"p\">,</span>
    <span class=\"kc\">True</span><span class=\"p\">)</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">env_options</span><span
    class=\"o\">.</span><span class=\"n\">setdefault</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;trim_blocks&quot;</span><span class=\"p\">,</span> <span class=\"kc\">True</span><span
    class=\"p\">)</span>\n\n                <span class=\"n\">env</span> <span class=\"o\">=</span>
    <span class=\"n\">jinja2</span><span class=\"o\">.</span><span class=\"n\">Environment</span><span
    class=\"p\">(</span><span class=\"o\">**</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">env_options</span><span class=\"p\">)</span>\n\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_jinja_env</span> <span class=\"o\">=</span> <span class=\"n\">env</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">env</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='default_template' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>default_template <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">default_template
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
    <span class=\"nf\">default_template</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
    class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;markata&quot;</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">post_template</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">v</span><span
    class=\"p\">,</span> <span class=\"nb\">str</span><span class=\"p\">):</span>\n
    \                   <span class=\"n\">v</span> <span class=\"o\">=</span> <span
    class=\"p\">{</span><span class=\"s2\">&quot;index&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">v</span><span class=\"p\">}</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;markata&quot;</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">post_template</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
    class=\"p\">):</span>\n                    <span class=\"n\">config_template</span>
    <span class=\"o\">=</span> <span class=\"p\">{</span>\n                        <span
    class=\"s2\">&quot;index&quot;</span><span class=\"p\">:</span> <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;markata&quot;</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">post_template</span><span class=\"p\">,</span>\n                    <span
    class=\"p\">}</span>\n                <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">config_template</span> <span class=\"o\">=</span>
    <span class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">&quot;markata&quot;</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">post_template</span>\n                <span
    class=\"k\">return</span> <span class=\"p\">{</span><span class=\"o\">**</span><span
    class=\"n\">config_template</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"n\">v</span><span class=\"p\">}</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='templates' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>templates <em class='small'>function</em></h2>\ntemplate management</p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">templates
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
    <span class=\"nf\">templates</span><span class=\"p\">():</span>\n                <span
    class=\"s2\">&quot;template management&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='show' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>show <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">show
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
    <span class=\"nf\">show</span><span class=\"p\">(</span>\n                <span
    class=\"n\">template</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
    class=\"n\">Argument</span><span class=\"p\">(</span><span class=\"kc\">None</span><span
    class=\"p\">,</span> <span class=\"n\">help</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;template to show&quot;</span><span class=\"p\">),</span>\n
    \               <span class=\"n\">theme</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"n\">typer</span><span
    class=\"o\">.</span><span class=\"n\">Option</span><span class=\"p\">(</span><span
    class=\"kc\">None</span><span class=\"p\">,</span> <span class=\"n\">help</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;pygments syntax theme&quot;</span><span
    class=\"p\">),</span>\n            <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">quiet</span> <span class=\"o\">=</span>
    <span class=\"kc\">True</span>\n                <span class=\"k\">if</span> <span
    class=\"n\">template</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">template</span> <span class=\"o\">=</span> <span class=\"n\">get_template</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">template</span><span class=\"p\">)</span>\n\n                    <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">quiet</span> <span class=\"o\">=</span>
    <span class=\"kc\">False</span>\n                    <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">print</span><span class=\"p\">(</span><span class=\"n\">template</span><span
    class=\"o\">.</span><span class=\"n\">filename</span><span class=\"p\">)</span>\n
    \                   <span class=\"k\">if</span> <span class=\"n\">theme</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span> <span class=\"ow\">or</span>
    <span class=\"n\">theme</span><span class=\"o\">.</span><span class=\"n\">lower</span><span
    class=\"p\">()</span> <span class=\"o\">==</span> <span class=\"s2\">&quot;none&quot;</span><span
    class=\"p\">:</span>\n                        <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">print</span><span class=\"p\">(</span><span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">template</span><span class=\"o\">.</span><span
    class=\"n\">filename</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">read_text</span><span class=\"p\">())</span>\n                    <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">syntax</span> <span class=\"o\">=</span> <span class=\"n\">Syntax</span><span
    class=\"o\">.</span><span class=\"n\">from_path</span><span class=\"p\">(</span><span
    class=\"n\">template</span><span class=\"o\">.</span><span class=\"n\">filename</span><span
    class=\"p\">,</span> <span class=\"n\">theme</span><span class=\"o\">=</span><span
    class=\"n\">theme</span><span class=\"p\">)</span>\n                        <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">print</span><span class=\"p\">(</span><span
    class=\"n\">syntax</span><span class=\"p\">)</span>\n\n                    <span
    class=\"k\">return</span>\n                <span class=\"n\">templates</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">jinja_env</span><span
    class=\"o\">.</span><span class=\"n\">list_templates</span><span class=\"p\">()</span>\n
    \               <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">quiet</span>
    <span class=\"o\">=</span> <span class=\"kc\">False</span>\n                <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">print</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;Templates directories:&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">style</span><span class=\"o\">=</span><span class=\"s2\">&quot;green
    underline&quot;</span><span class=\"p\">)</span>\n\n                <span class=\"n\">markata_templates</span>
    <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"vm\">__file__</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">parents</span><span class=\"p\">[</span><span class=\"mi\">1</span><span
    class=\"p\">]</span> <span class=\"o\">/</span> <span class=\"s2\">&quot;templates&quot;</span>\n
    \               <span class=\"k\">for</span> <span class=\"nb\">dir</span> <span
    class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">templates_dir</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">if</span> <span class=\"nb\">dir</span>
    <span class=\"o\">==</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">dynamic_templates_dir</span><span
    class=\"p\">:</span>\n                        <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">print</span><span class=\"p\">(</span>\n                            <span
    class=\"sa\">f</span><span class=\"s2\">&quot;[gold3]</span><span class=\"si\">{</span><span
    class=\"nb\">dir</span><span class=\"si\">}</span><span class=\"s2\">[/][grey50]
    (dynamically created templates from configuration)[/] [gold3]\\[markata.config.dynamic_templates_dir][/]&quot;</span><span
    class=\"p\">,</span>\n                            <span class=\"n\">style</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;red&quot;</span><span class=\"p\">,</span>\n
    \                       <span class=\"p\">)</span>\n                    <span
    class=\"k\">elif</span> <span class=\"nb\">dir</span> <span class=\"o\">==</span>
    <span class=\"n\">markata_templates</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">print</span><span class=\"p\">(</span>\n
    \                           <span class=\"sa\">f</span><span class=\"s2\">&quot;[cyan]</span><span
    class=\"si\">{</span><span class=\"nb\">dir</span><span class=\"si\">}</span><span
    class=\"s2\">[/][grey50] (built-in)[/]&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">style</span><span class=\"o\">=</span><span class=\"s2\">&quot;red&quot;</span>\n
    \                       <span class=\"p\">)</span>\n                    <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                        <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">print</span><span class=\"p\">(</span>\n
    \                           <span class=\"sa\">f</span><span class=\"s2\">&quot;[orchid]</span><span
    class=\"si\">{</span><span class=\"nb\">dir</span><span class=\"si\">}</span><span
    class=\"s2\">[/] [orchid]\\[markata.config.templates_dir][/]&quot;</span><span
    class=\"p\">,</span>\n                            <span class=\"n\">style</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;red&quot;</span><span class=\"p\">,</span>\n
    \                       <span class=\"p\">)</span>\n\n                <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">print</span><span class=\"p\">()</span>\n                <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">print</span><span class=\"p\">(</span>\n                    <span
    class=\"s2\">&quot;Available Templates: [white]name -&gt; path[/]&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">style</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;green underline&quot;</span>\n                <span class=\"p\">)</span>\n
    \               <span class=\"k\">for</span> <span class=\"n\">template</span>
    <span class=\"ow\">in</span> <span class=\"n\">templates</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">source</span><span class=\"p\">,</span>
    <span class=\"n\">file</span><span class=\"p\">,</span> <span class=\"n\">uptodate</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">jinja_env</span><span
    class=\"o\">.</span><span class=\"n\">loader</span><span class=\"o\">.</span><span
    class=\"n\">get_source</span><span class=\"p\">(</span>\n                        <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">jinja_env</span><span class=\"p\">,</span>
    <span class=\"n\">template</span>\n                    <span class=\"p\">)</span>\n\n
    \                   <span class=\"k\">if</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">file</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">is_relative_to</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">dynamic_templates_dir</span><span class=\"p\">):</span>\n
    \                       <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">print</span><span
    class=\"p\">(</span>\n                            <span class=\"sa\">f</span><span
    class=\"s2\">&quot;[gold3]</span><span class=\"si\">{</span><span class=\"n\">template</span><span
    class=\"si\">}</span><span class=\"s2\"> -&gt; [red]</span><span class=\"si\">{</span><span
    class=\"n\">file</span><span class=\"si\">}</span><span class=\"s2\">[/] [grey50](dynamic)[/]&quot;</span>\n
    \                       <span class=\"p\">)</span>\n                    <span
    class=\"k\">elif</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">file</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">is_relative_to</span><span class=\"p\">(</span><span class=\"n\">markata_templates</span><span
    class=\"p\">):</span>\n                        <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">print</span><span class=\"p\">(</span>\n                            <span
    class=\"sa\">f</span><span class=\"s2\">&quot;[cyan]</span><span class=\"si\">{</span><span
    class=\"n\">template</span><span class=\"si\">}</span><span class=\"s2\"> -&gt;
    [red]</span><span class=\"si\">{</span><span class=\"n\">file</span><span class=\"si\">}</span><span
    class=\"s2\">[/] [grey50](built-in)[/]&quot;</span>\n                        <span
    class=\"p\">)</span>\n                    <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">print</span><span
    class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">&quot;[orchid]</span><span
    class=\"si\">{</span><span class=\"n\">template</span><span class=\"si\">}</span><span
    class=\"s2\">[/] -&gt; [red]</span><span class=\"si\">{</span><span class=\"n\">file</span><span
    class=\"si\">}</span><span class=\"s2\">[/]&quot;</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>"
  raw.md: ''
published: true
slug: markata/plugins/post-template
title: Post_Template.Py


---

# Add head configuration

This snippet allows users to configure their head in `markata.toml`.

``` html
{{ config.get('head', {}).pop('text') if 'text' in config.get('head',{}).keys() }}
{% for tag, meta in config.get('head', {}).items() %}
    {% for _meta in meta %}
        <{{ tag }}
            {% for attr, value in _meta.items() %}{{ attr }}="{{ value }}"{% endfor %}
        />
    {% endfor %}
{% endfor %}
```

Users can specify any sort of tag in their `markata.toml`

``` toml
[[markata.head.meta]]
name = "og:type"
content = "article"

[[markata.head.meta]]
name = "og:author"
content = "Waylon Walker"
```

The above configuration becomes this once rendered.

``` html
<meta name='og:type' content='article' />
<meta name='og:Author' content='Waylon Walker' />
```

!! Note

    Article variables can be used for dynamic entries like canonical_url
    ``` toml
    [markata]
    url = "markata.dev"

    [[markata.head.meta]]
    href="{{ config.url }}/{{ slug }}/"
    rel="canonical"
    ```

Optionally users can also specify plain text to be appended to the head of
their documents.  This works well for things that involve full blocks.

``` toml
[[markata.head.text]]
value = '''
<script>
    console.log('hello world')
</script>
'''

[[markata.head.text]]
value='''
html  {
    font-family: "Space Mono", monospace;
    background: var(--color-bg);
    color: var(--color-text);
}
'''
```

## Add scripts to head

Markata config also supports adding scripts to the head via configuration.

``` toml
[[ markata.head.script ]]
    src = "https://cdn.tailwindcss.com"

```


!! class <h2 id='SilentUndefined' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>SilentUndefined <em class='small'>class</em></h2>

???+ source "SilentUndefined <em class='small'>source</em>"

```python

        class SilentUndefined(Undefined):
            def _fail_with_undefined_error(self, *args, **kwargs):
                return ""
```


!! function <h2 id='optional' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>optional <em class='small'>function</em></h2>

???+ source "optional <em class='small'>source</em>"

```python

        def optional(*fields):
            def dec(_cls):
                for field in fields:
                    _cls.__fields__[field].default = None
                return _cls

            if (
                fields
                and inspect.isclass(fields[0])
                and issubclass(fields[0], pydantic.BaseModel)
            ):
                cls = fields[0]
                fields = cls.__fields__
                return dec(cls)
            return dec
```


!! class <h2 id='Style' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Style <em class='small'>class</em></h2>

???+ source "Style <em class='small'>source</em>"

```python

        class Style(pydantic.BaseModel):
            color_bg: str = "#1f2022"
            color_bg_code: str = "#1f2022"
            color_text: str = "#eefbfe"
            color_link: str = "#fb30c4"
            color_accent: str = "#e1bd00c9"
            overlay_brightness: str = ".85"
            body_width: str = "800px"
            color_bg_light: str = "#eefbfe"
            color_bg_code_light: str = "#eefbfe"
            color_text_light: str = "#1f2022"
            color_link_light: str = "#fb30c4"
            color_accent_light: str = "#ffeb00"
            overlay_brightness_light: str = ".95"
```


!! class <h2 id='StyleOverrides' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>StyleOverrides <em class='small'>class</em></h2>

???+ source "StyleOverrides <em class='small'>source</em>"

```python

        class StyleOverrides(Style): ...
```


!! class <h2 id='Meta' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Meta <em class='small'>class</em></h2>

???+ source "Meta <em class='small'>source</em>"

```python

        class Meta(pydantic.BaseModel):
            name: str
            content: str
```


!! class <h2 id='Text' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Text <em class='small'>class</em></h2>

???+ source "Text <em class='small'>source</em>"

```python

        class Text(pydantic.BaseModel):
            value: str
```


!! class <h2 id='Link' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Link <em class='small'>class</em></h2>

???+ source "Link <em class='small'>source</em>"

```python

        class Link(pydantic.BaseModel):
            rel: str = "canonical"
            href: str
```


!! class <h2 id='Script' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Script <em class='small'>class</em></h2>

???+ source "Script <em class='small'>source</em>"

```python

        class Script(pydantic.BaseModel):
            src: str
```


!! class <h2 id='HeadConfig' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>HeadConfig <em class='small'>class</em></h2>

???+ source "HeadConfig <em class='small'>source</em>"

```python

        class HeadConfig(pydantic.BaseModel):
            meta: List[Meta] = []
            link: List[Link] = []
            script: List[Script] = []
            text: Union[List[Text], str] = ""

            @pydantic.validator("text", pre=True)
            def text_to_list(cls, v):
                if isinstance(v, list):
                    return "\n".join([text["value"] for text in v])
                return v

            @property
            def html(self):
                html = self.text
                html += "\n"
                for meta in self.meta:
                    html += f'<meta name="{meta.name}" content="{meta.content}" />\n'
                for link in self.link:
                    html += f'<link rel="{link.rel}" href="{link.href}" />\n'
                return html
```


!! class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Config <em class='small'>class</em></h2>

???+ source "Config <em class='small'>source</em>"

```python

        class Config(pydantic.BaseModel):
            head: HeadConfig = HeadConfig()
            style: Style = Style()
            post_template: Optional[Union[str | Dict[str, str]]] = "post.html"
            dynamic_templates_dir: Path = Path(".markata.cache/templates")
            templates_dir: Union[Path, List[Path]] = pydantic.Field(Path("templates"))

            env_options: dict = {}

            @pydantic.model_validator(mode="after")
            def dynamic_templates_in_templates_dir(self):
                markata_templates = Path(__file__).parents[1] / "templates"

                if isinstance(self.templates_dir, Path):
                    self.templates_dir = [
                        self.templates_dir,
                        markata_templates,
                        self.dynamic_templates_dir,
                    ]

                if markata_templates not in self.templates_dir:
                    self.templates_dir.append(markata_templates)

                if self.dynamic_templates_dir not in self.templates_dir:
                    self.templates_dir.append(self.dynamic_templates_dir)

                return self

            @property
            def jinja_loader(self):
                return jinja2.FileSystemLoader(self.templates_dir)

            @property
            def jinja_env(
                self,
            ):
                if hasattr(self, "_jinja_env"):
                    return self._jinja_env
                self.env_options.setdefault("loader", self.jinja_loader)
                self.env_options.setdefault("undefined", SilentUndefined)
                self.env_options.setdefault("lstrip_blocks", True)
                self.env_options.setdefault("trim_blocks", True)

                env = jinja2.Environment(**self.env_options)

                self._jinja_env = env
                return env
```


!! class <h2 id='PostOverrides' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>PostOverrides <em class='small'>class</em></h2>

???+ source "PostOverrides <em class='small'>source</em>"

```python

        class PostOverrides(pydantic.BaseModel):
            head: HeadConfig = HeadConfig()
            style: Style = StyleOverrides()
```


!! class <h2 id='Post' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Post <em class='small'>class</em></h2>

???+ source "Post <em class='small'>source</em>"

```python

        class Post(pydantic.BaseModel):
            config_overrides: PostOverrides = PostOverrides()
            template: Optional[str | Dict[str, str]] = None

            @pydantic.validator("template", pre=True, always=True)
            def default_template(cls, v, *, values):
                if v is None:
                    return values["markata"].config.post_template
                if isinstance(v, str):
                    v = {"index": v}
                if isinstance(values["markata"].config.post_template, str):
                    config_template = {
                        "index": values["markata"].config.post_template,
                    }
                else:
                    config_template = values["markata"].config.post_template
                return {**config_template, **v}
```


!! function <h2 id='config_model' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>config_model <em class='small'>function</em></h2>

???+ source "config_model <em class='small'>source</em>"

```python

        def config_model(markata: "Markata") -> None:
            markata.config_models.append(Config)
```


!! function <h2 id='post_model' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>post_model <em class='small'>function</em></h2>

???+ source "post_model <em class='small'>source</em>"

```python

        def post_model(markata: "Markata") -> None:
            markata.post_models.append(Post)
```


!! function <h2 id='configure' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>configure <em class='small'>function</em></h2>
    Massages the configuration limitations of toml to allow a little bit easier
    experience to the end user making configurations while allowing an simpler
    jinja template.  This enablees the use of the `markata.head.text` list in
    configuration.
???+ source "configure <em class='small'>source</em>"

```python

        def configure(markata: "Markata") -> None:
            """
            Massages the configuration limitations of toml to allow a little bit easier
            experience to the end user making configurations while allowing an simpler
            jinja template.  This enablees the use of the `markata.head.text` list in
            configuration.
            """
```


!! function <h2 id='pre_render' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>pre_render <em class='small'>function</em></h2>
    FOR EACH POST: Massages the configuration limitations of toml/yaml to allow
    a little bit easier experience to the end user making configurations while
    allowing an simpler jinja template.  This enables the use of the
    `markata.head.text` list in configuration.
???+ source "pre_render <em class='small'>source</em>"

```python

        def pre_render(markata: "Markata") -> None:
            """
            FOR EACH POST: Massages the configuration limitations of toml/yaml to allow
            a little bit easier experience to the end user making configurations while
            allowing an simpler jinja template.  This enables the use of the
            `markata.head.text` list in configuration.
            """

            markata.config.dynamic_templates_dir.mkdir(parents=True, exist_ok=True)
            head_template = markata.config.dynamic_templates_dir / "head.html"
            head_template.write_text(
                markata.config.jinja_env.get_template("dynamic_head.html").render(
                    {"markata": markata}
                ),
            )

            for article in [a for a in markata.articles if "config_overrides" in a]:
                raw_text = article.get("config_overrides", {}).get("head", {}).get("text", "")

                if isinstance(raw_text, list):
                    article["config_overrides"]["head"]["text"] = "\n".join(
                        flatten([t.values() for t in raw_text]),
                    )
```


!! function <h2 id='render' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>render <em class='small'>function</em></h2>

???+ source "render <em class='small'>source</em>"

```python

        def render(markata: "Markata") -> None:
            with markata.cache as cache:
                for article in markata.articles:
                    html = render_article(markata=markata, cache=cache, article=article)
                    article.html = html
```


!! function <h2 id='get_template' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get_template <em class='small'>function</em></h2>

???+ source "get_template <em class='small'>source</em>"

```python

        def get_template(markata, template):
            try:
                return markata.config.jinja_env.get_template(template)
            except jinja2.TemplateNotFound:
                # try to load it as a file
                ...

            try:
                return Template(Path(template).read_text(), undefined=SilentUndefined)
            except FileNotFoundError:
                # default to load it as a string
                ...
            return Template(template, undefined=SilentUndefined)
```


!! function <h2 id='render_article' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>render_article <em class='small'>function</em></h2>

???+ source "render_article <em class='small'>source</em>"

```python

        def render_article(markata, cache, article):
            key = markata.make_hash(
                "post_template",
                __version__,
                article.key,
            )
            html = markata.precache.get(key)

            if html is not None:
                return html

            if isinstance(article.template, str):
                template = get_template(markata, article.template)
                html = render_template(markata, article, template)

            if isinstance(article.template, dict):
                html = {
                    slug: render_template(markata, article, get_template(markata, template))
                    for slug, template in article.template.items()
                }
            cache.add(key, html, expire=markata.config.default_cache_expire)
            return html
```


!! function <h2 id='render_template' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>render_template <em class='small'>function</em></h2>

???+ source "render_template <em class='small'>source</em>"

```python

        def render_template(markata, article, template):
            template = get_template(markata, template)
            merged_config = markata.config
            # TODO do we need to handle merge??
            # if head_template:
            #     head = eval(
            #         head_template.render(
            #             __version__=__version__,
            #             config=_full_config,
            #             **article,
            #         )
            #     )

            # merged_config = {
            #     **_full_config,
            #     **{"head": head},
            # }

            # merged_config = always_merger.merge(
            #     merged_config,
            #     copy.deepcopy(
            #         article.get(
            #             "config_overrides",
            #             {},
            #         )
            #     ),
            # )

            html = template.render(
                __version__=__version__,
                markata=markata,
                body=article.article_html,
                config=merged_config,
                post=article,
            )
            return html
```


!! function <h2 id='save' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>save <em class='small'>function</em></h2>

???+ source "save <em class='small'>source</em>"

```python

        def save(markata: "Markata") -> None:
            linked_templates = [
                t
                for t in markata.config.jinja_env.list_templates()
                if t.endswith("css") or t.endswith("js") or t.endswith("xsl")
            ]
            for template in linked_templates:
                template = get_template(markata, template)
                css = template.render(markata=markata, __version__=__version__)
                Path(markata.config.output_dir / Path(template.filename).name).write_text(css)
```


!! function <h2 id='cli' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>cli <em class='small'>function</em></h2>
    Markata hook to implement base cli commands.
???+ source "cli <em class='small'>source</em>"

```python

        def cli(app: typer.Typer, markata: "Markata") -> None:
            """
            Markata hook to implement base cli commands.
            """

            templates_app = typer.Typer()
            app.add_typer(templates_app)

            @templates_app.callback()
            def templates():
                "template management"

            @templates_app.command()
            def show(
                template: str = typer.Argument(None, help="template to show"),
                theme: str = typer.Option(None, help="pygments syntax theme"),
            ) -> None:
                markata.console.quiet = True
                if template:
                    template = get_template(markata, template)

                    markata.console.quiet = False
                    markata.console.print(template.filename)
                    if theme is None or theme.lower() == "none":
                        markata.console.print(Path(template.filename).read_text())
                    else:
                        syntax = Syntax.from_path(template.filename, theme=theme)
                        markata.console.print(syntax)

                    return
                templates = markata.config.jinja_env.list_templates()
                markata.console.quiet = False
                markata.console.print("Templates directories:", style="green underline")

                markata_templates = Path(__file__).parents[1] / "templates"
                for dir in markata.config.templates_dir:
                    if dir == markata.config.dynamic_templates_dir:
                        markata.console.print(
                            f"[gold3]{dir}[/][grey50] (dynamically created templates from configuration)[/] [gold3]\[markata.config.dynamic_templates_dir][/]",
                            style="red",
                        )
                    elif dir == markata_templates:
                        markata.console.print(
                            f"[cyan]{dir}[/][grey50] (built-in)[/]", style="red"
                        )
                    else:
                        markata.console.print(
                            f"[orchid]{dir}[/] [orchid]\[markata.config.templates_dir][/]",
                            style="red",
                        )

                markata.console.print()
                markata.console.print(
                    "Available Templates: [white]name -> path[/]", style="green underline"
                )
                for template in templates:
                    source, file, uptodate = markata.config.jinja_env.loader.get_source(
                        markata.config.jinja_env, template
                    )

                    if Path(file).is_relative_to(markata.config.dynamic_templates_dir):
                        markata.console.print(
                            f"[gold3]{template} -> [red]{file}[/] [grey50](dynamic)[/]"
                        )
                    elif Path(file).is_relative_to(markata_templates):
                        markata.console.print(
                            f"[cyan]{template} -> [red]{file}[/] [grey50](built-in)[/]"
                        )
                    else:
                        markata.console.print(f"[orchid]{template}[/] -> [red]{file}[/]")
```


!! method <h2 id='_fail_with_undefined_error' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_fail_with_undefined_error <em class='small'>method</em></h2>

???+ source "_fail_with_undefined_error <em class='small'>source</em>"

```python

        def _fail_with_undefined_error(self, *args, **kwargs):
                return ""
```


!! function <h2 id='dec' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>dec <em class='small'>function</em></h2>

???+ source "dec <em class='small'>source</em>"

```python

        def dec(_cls):
                for field in fields:
                    _cls.__fields__[field].default = None
                return _cls
```


!! method <h2 id='text_to_list' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>text_to_list <em class='small'>method</em></h2>

???+ source "text_to_list <em class='small'>source</em>"

```python

        def text_to_list(cls, v):
                if isinstance(v, list):
                    return "\n".join([text["value"] for text in v])
                return v
```


!! method <h2 id='html' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>html <em class='small'>method</em></h2>

???+ source "html <em class='small'>source</em>"

```python

        def html(self):
                html = self.text
                html += "\n"
                for meta in self.meta:
                    html += f'<meta name="{meta.name}" content="{meta.content}" />\n'
                for link in self.link:
                    html += f'<link rel="{link.rel}" href="{link.href}" />\n'
                return html
```


!! method <h2 id='dynamic_templates_in_templates_dir' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>dynamic_templates_in_templates_dir <em class='small'>method</em></h2>

???+ source "dynamic_templates_in_templates_dir <em class='small'>source</em>"

```python

        def dynamic_templates_in_templates_dir(self):
                markata_templates = Path(__file__).parents[1] / "templates"

                if isinstance(self.templates_dir, Path):
                    self.templates_dir = [
                        self.templates_dir,
                        markata_templates,
                        self.dynamic_templates_dir,
                    ]

                if markata_templates not in self.templates_dir:
                    self.templates_dir.append(markata_templates)

                if self.dynamic_templates_dir not in self.templates_dir:
                    self.templates_dir.append(self.dynamic_templates_dir)

                return self
```


!! method <h2 id='jinja_loader' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>jinja_loader <em class='small'>method</em></h2>

???+ source "jinja_loader <em class='small'>source</em>"

```python

        def jinja_loader(self):
                return jinja2.FileSystemLoader(self.templates_dir)
```


!! method <h2 id='jinja_env' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>jinja_env <em class='small'>method</em></h2>

???+ source "jinja_env <em class='small'>source</em>"

```python

        def jinja_env(
                self,
            ):
                if hasattr(self, "_jinja_env"):
                    return self._jinja_env
                self.env_options.setdefault("loader", self.jinja_loader)
                self.env_options.setdefault("undefined", SilentUndefined)
                self.env_options.setdefault("lstrip_blocks", True)
                self.env_options.setdefault("trim_blocks", True)

                env = jinja2.Environment(**self.env_options)

                self._jinja_env = env
                return env
```


!! method <h2 id='default_template' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>default_template <em class='small'>method</em></h2>

???+ source "default_template <em class='small'>source</em>"

```python

        def default_template(cls, v, *, values):
                if v is None:
                    return values["markata"].config.post_template
                if isinstance(v, str):
                    v = {"index": v}
                if isinstance(values["markata"].config.post_template, str):
                    config_template = {
                        "index": values["markata"].config.post_template,
                    }
                else:
                    config_template = values["markata"].config.post_template
                return {**config_template, **v}
```


!! function <h2 id='templates' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>templates <em class='small'>function</em></h2>
    template management
???+ source "templates <em class='small'>source</em>"

```python

        def templates():
                "template management"
```


!! function <h2 id='show' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>show <em class='small'>function</em></h2>

???+ source "show <em class='small'>source</em>"

```python

        def show(
                template: str = typer.Argument(None, help="template to show"),
                theme: str = typer.Option(None, help="pygments syntax theme"),
            ) -> None:
                markata.console.quiet = True
                if template:
                    template = get_template(markata, template)

                    markata.console.quiet = False
                    markata.console.print(template.filename)
                    if theme is None or theme.lower() == "none":
                        markata.console.print(Path(template.filename).read_text())
                    else:
                        syntax = Syntax.from_path(template.filename, theme=theme)
                        markata.console.print(syntax)

                    return
                templates = markata.config.jinja_env.list_templates()
                markata.console.quiet = False
                markata.console.print("Templates directories:", style="green underline")

                markata_templates = Path(__file__).parents[1] / "templates"
                for dir in markata.config.templates_dir:
                    if dir == markata.config.dynamic_templates_dir:
                        markata.console.print(
                            f"[gold3]{dir}[/][grey50] (dynamically created templates from configuration)[/] [gold3]\[markata.config.dynamic_templates_dir][/]",
                            style="red",
                        )
                    elif dir == markata_templates:
                        markata.console.print(
                            f"[cyan]{dir}[/][grey50] (built-in)[/]", style="red"
                        )
                    else:
                        markata.console.print(
                            f"[orchid]{dir}[/] [orchid]\[markata.config.templates_dir][/]",
                            style="red",
                        )

                markata.console.print()
                markata.console.print(
                    "Available Templates: [white]name -> path[/]", style="green underline"
                )
                for template in templates:
                    source, file, uptodate = markata.config.jinja_env.loader.get_source(
                        markata.config.jinja_env, template
                    )

                    if Path(file).is_relative_to(markata.config.dynamic_templates_dir):
                        markata.console.print(
                            f"[gold3]{template} -> [red]{file}[/] [grey50](dynamic)[/]"
                        )
                    elif Path(file).is_relative_to(markata_templates):
                        markata.console.print(
                            f"[cyan]{template} -> [red]{file}[/] [grey50](built-in)[/]"
                        )
                    else:
                        markata.console.print(f"[orchid]{template}[/] -> [red]{file}[/]")
```

