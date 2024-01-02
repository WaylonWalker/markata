---
content: "The `markata.plugins.feeds` plugin is used to create feed pages, which are
  lists of\nposts.  The list is generated using a `filter`, then each post in the
  list is\nrendered with a `card_template` before being applied to the `body` of the\n`template`.\n\n#
  Installation\n\nThis plugin is built-in and enabled by default, but in you want
  to be very\nexplicit you can add it to your list of existing plugins.\n\n``` toml\nhooks
  = [\n   \"markata.plugins.feeds\",\n   ]\n```\n\n# Configuration\n\n# set default
  template and card_template\n\nAt the root of the markata.feeds config you may set
  `template`, and\n`card_template`.  These will become your defaults for every feed
  you create.\nIf you do not set these, markata will use it's defaults.  The defaults
  are\ndesigned to work for a variety of use cases, but are not likely the best for\nall.\n\n```
  toml\n[markata.feeds_config]\ntemplate=\"pages/templates/archive_template.html\"\ncard_template=\"plugins/feed_card_template.html\"\n```\n\n#
  pages\n\nUnderneath of the `markata.feeds` we will create a new map for each page
  where\nthe name of the map will be the name of the page.\n\n\nThe following config
  will create a page at `/all-posts` that inclues every\nsingle post.\n\n``` toml\n[[markata.feeds]]\ntitle=\"All
  Posts\"\nslug='all'\nfilter=\"True\"\n```\n\n# template\n\nThe `template` configuration
  key is a file path to the template that you want\nto use to create the feed.  You
  may set the default template you want to use\nfor all feeds under `[markata.feeds]`,
  as well as override it inside of each\nfeeds config.\n\nThe template is a jinja
  style template that expects to fill in a `title` and\n`body` variable.\n\n``` html\n<!DOCTYPE
  html>\n<html lang=\"en\">\n  <head>\n    <title>{{ title }}</title>\n  </head>\n
  \ <body>\n    <ul>\n        {{ body }}\n    </ul>\n  </body>\n</html>\n```\n\n!!!
  note\n    I highly reccomend putting your `body` in a `<ul>`, and wrapping your\n
  \   `card_template`s in an `<li>`.\n\n# card_template\n\nAll keys available from
  each post is available to put into your jinja\ntemplate.  These can either be placed
  there in your post frontmatter, or\nthrough a plugin that automatically adds to
  the post before the save phase.\n\nHere is a very simple example that would give
  a link to each post with the\ntitle and date.\n\n``` toml\n[[markata.feeds]]\nslug='all'\ntitle='All
  Posts'\nfilter=\"True\"\ncard_template='''\n<li>\n    <a href={{markata.config.get('path_prefix',
  '')}}{{slug}}>\n        {{title}}-{{date}}\n    </a>\n</li>\n'''\n```\n\n# filter\n\nThe
  filter is a python expression ran on every post that expects to return a\nboolean.
  \ The variables available to this expression are every key in your\nfrontmatter,
  plus the `timedelta` function, and `parse` function to more easily\nwork with dates.\n\n#
  Feed Examples\n\nTrue can be passed in to make a feed of all the posts you have.\n\n```
  toml\n[[markata.feeds]]\nslug='all'\ntitle='All Posts'\nfilter=\"True\"\n```\n\nYou
  can compare against the values of the keys from your frontmatter.  This\nexample
  creates a feed that includes every post where published is `True`.\n\n``` toml\n[[markata.feeds]]\nslug='draft'\ntitle='Draft'\nfilter=\"published=='False'\"\n```\n\nWe
  can also compare against dates.  The\n[markata.plugins.datetime](https://markata.dev/markata/plugins/datetime/)\nplugin,
  automatically adds `today` as today's date and `now` as the current\ndatetime.  These
  are quite handy to create feeds for scheduled, recent, or\ntoday's posts.  The following
  two examples will create a feed for scheduled\nposts and for today's posts respectively.\n\n```
  toml\n[[markata.feeds]]\nslug='scheduled'\ntitle='Scheduled'\nfilter=\"date>today\"\n\n[[markata.feeds]]\nslug='today'\ntitle='Today'\nfilter=\"date==today\"\n```\n\nIf
  you have list of items in your frontmatter for something like `tags`, you\ncan check
  for the existence of a tag in the list.\n\n``` toml\n[[markata.feeds]]\nslug='python'\ntitle='Python'\nfilter=\"date<=today
  and 'python' in tags\"\n```\n\nAnd of course you can combine all the things into
  larger expressions.  Here is\none example of the main feed on my blog.\n\n``` toml\n[[markata.feeds]]\nslug='blog'\ntitle='Blog'\nfilter=\"date<=today
  and templateKey in ['blog-post'] and published =='True'\"\n```\n\nHere is another
  example that shows my drafts for a particular tag.\n\n``` toml\n[[markata.feeds]]\nslug='python-draft'\ntitle='Python
  Draft'\nfilter=\"date<=today and 'python' in tags and published=='False'\"\n```\n\n#
  Defaults\n\nBy default feeds will create one feed page at `/archive/` that includes
  all\nposts.\n\n[[markata.feeds]]\nslug='archive'\ntitle='All Posts'\nfilter=\"True\"\n\n\n!!
  class <h2 id='SilentUndefined' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>SilentUndefined <em class='small'>class</em></h2>\n\n???+ source \"SilentUndefined
  <em class='small'>source</em>\"\n\n```python\n\n        class SilentUndefined(Undefined):\n
  \           def _fail_with_undefined_error(self, *args, **kwargs):\n                return
  \"\"\n```\n\n\n!! class <h2 id='MarkataFilterError' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>MarkataFilterError <em class='small'>class</em></h2>\n\n???+ source \"MarkataFilterError
  <em class='small'>source</em>\"\n\n```python\n\n        class MarkataFilterError(RuntimeError):\n
  \           ...\n```\n\n\n!! class <h2 id='FeedConfig' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>FeedConfig <em class='small'>class</em></h2>\n\n???+
  source \"FeedConfig <em class='small'>source</em>\"\n\n```python\n\n        class
  FeedConfig(pydantic.BaseModel):\n            DEFAULT_TITLE: str = \"All Posts\"\n\n
  \           title: str = DEFAULT_TITLE\n            slug: str = None\n            name:
  Optional[str] = None\n            filter: str = \"True\"\n            sort: str
  = \"date\"\n            reverse: bool = False\n            rss: bool = True\n            sitemap:
  bool = True\n            card_template: str = \"\"\"\n                <li class='post'>\n
  \                   <a href=\"/{{ markata.config.path_prefix }}{{ post.slug }}/\">\n
  \                       {{ post.title }}\n                    </a>\n                </li>\n
  \               \"\"\"\n            template: str = Path(__file__).parent / \"default_post_template.html.jinja\"\n
  \           rss_template: str = Path(__file__).parent / \"default_rss_template.xml\"\n
  \           sitemap_template: str = Path(__file__).parent / \"default_sitemap_template.xml\"\n
  \           xsl_template: str = Path(__file__).parent / \"default_xsl_template.xsl\"\n\n
  \           @pydantic.validator(\"name\", pre=True, always=True)\n            def
  default_name(cls, v, *, values):\n                return v or str(values.get(\"slug\")).replace(\"-\",
  \"_\")\n\n            @pydantic.validator(\"card_template\", \"template\", pre=True,
  always=True)\n            def read_template(cls, v, *, values) -> str:\n                if
  isinstance(v, Path):\n                    return str(v.read_text())\n                return
  v\n```\n\n\n!! class <h2 id='FeedsConfig' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>FeedsConfig <em class='small'>class</em></h2>\n\n???+ source \"FeedsConfig
  <em class='small'>source</em>\"\n\n```python\n\n        class FeedsConfig(pydantic.BaseModel):\n
  \           feeds: List[FeedConfig] = [FeedConfig(slug=\"archive\")]\n```\n\n\n!!
  class <h2 id='Feed' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Feed
  <em class='small'>class</em></h2>\n    A storage class for markata feed objects.\n\n
  \   # Usage\n\n    ``` python\n    from markata import Markata\n    m = Markata()\n\n
  \   # access posts for a feed\n    m.feeds.docs.posts\n\n    # access config for
  a feed\n    m.feeds.docs.config\n    ```\n???+ source \"Feed <em class='small'>source</em>\"\n\n```python\n\n
  \       class Feed:\n            \"\"\"\n            A storage class for markata
  feed objects.\n\n            # Usage\n\n            ``` python\n            from
  markata import Markata\n            m = Markata()\n\n            # access posts
  for a feed\n            m.feeds.docs.posts\n\n            # access config for a
  feed\n            m.feeds.docs.config\n            ```\n            \"\"\"\n\n            config:
  FeedConfig\n            _m: Markata\n\n            @property\n            def name(self):\n
  \               return self.config.name\n\n            @property\n            def
  posts(self):\n                return self.map(\"post\")\n\n            def map(self,
  func=\"post\", **args):\n                return self._m.map(func, **{**self.config.dict(),
  **args})\n```\n\n\n!! class <h2 id='Feeds' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>Feeds <em class='small'>class</em></h2>\n    A storage class for all markata
  Feed objects\n\n    ``` python\n    from markata import Markata\n    m = Markata()\n\n
  \   m.feeds\n\n    # access all config\n    m.feeds.config\n\n    # refresh list
  of posts in all feeds\n    m.feeds.refresh()\n\n\n    # iterating over feeds gives
  the name of the feed\n    for k in m.feeds:\n         print(k)\n\n    # project-gallery\n
  \   # docs\n    # autodoc\n    # core_modules\n    # plugins\n    # archive\n\n
  \   # iterate over items like keys and values in a dict, items returns name of\n
  \   # feed and a feed object\n    for k, v in m.feeds.items():\n        print(k,
  len(v.posts))\n\n    # project-gallery 2\n    # docs 6\n    # autodoc 65\n    #
  core_modules 26\n    # plugins 39\n    # archive 65\n\n    # values can be iterated
  over in just the same way\n    for v in m.feeds.values():\n         print(len(v.posts))\n
  \   # 2\n    # 6\n    # 65\n    # 26\n    # 39\n    # 65\n    ```\n\n    Accessing
  feeds can be done using square brackets or dot notation.\n\n    ``` python\n    from
  markata import Markata\n    m = Markata()\n\n    # both of these will return the
  `docs` Feed object.\n    m.feeds.docs\n    m['docs']\n    ```\n???+ source \"Feeds
  <em class='small'>source</em>\"\n\n```python\n\n        class Feeds:\n            \"\"\"\n
  \           A storage class for all markata Feed objects\n\n            ``` python\n
  \           from markata import Markata\n            m = Markata()\n\n            m.feeds\n\n
  \           # access all config\n            m.feeds.config\n\n            # refresh
  list of posts in all feeds\n            m.feeds.refresh()\n\n\n            # iterating
  over feeds gives the name of the feed\n            for k in m.feeds:\n                 print(k)\n\n
  \           # project-gallery\n            # docs\n            # autodoc\n            #
  core_modules\n            # plugins\n            # archive\n\n            # iterate
  over items like keys and values in a dict, items returns name of\n            #
  feed and a feed object\n            for k, v in m.feeds.items():\n                print(k,
  len(v.posts))\n\n            # project-gallery 2\n            # docs 6\n            #
  autodoc 65\n            # core_modules 26\n            # plugins 39\n            #
  archive 65\n\n            # values can be iterated over in just the same way\n            for
  v in m.feeds.values():\n                 print(len(v.posts))\n            # 2\n
  \           # 6\n            # 65\n            # 26\n            # 39\n            #
  65\n            ```\n\n            Accessing feeds can be done using square brackets
  or dot notation.\n\n            ``` python\n            from markata import Markata\n
  \           m = Markata()\n\n            # both of these will return the `docs`
  Feed object.\n            m.feeds.docs\n            m['docs']\n            ```\n
  \           \"\"\"\n\n            def __init__(self, markata: Markata) -> None:\n
  \               self._m = markata\n                self.config = {f.name: f for
  f in markata.config.feeds}\n                self.refresh()\n\n            def refresh(self)
  -> None:\n                \"\"\"\n                Refresh all of the feeds objects\n
  \               \"\"\"\n                for feed in self._m.config.feeds:\n                    feed
  = Feed(config=feed, _m=self._m)\n                    self.__setattr__(feed.name,
  feed)\n\n            def __iter__(self):\n                return iter(self.config.keys())\n\n
  \           def keys(self):\n                return iter(self.config.keys())\n\n
  \           def values(self):\n                return [self[feed] for feed in self.config.keys()]\n\n
  \           def items(self):\n                return [(key, self[key]) for key in
  self.config]\n\n            def __getitem__(self, key: str) -> Any:\n                return
  getattr(self, key.replace(\"-\", \"_\").lower())\n\n            def _dict_panel(self,
  config) -> str:\n                \"\"\"\n                pretty print configs with
  rich\n                \"\"\"\n                msg = \"\"\n                for key,
  value in config.items():\n                    if isinstance(value, str):\n                        if
  len(value) > 50:\n                            value = value[:50] + \"...\"\n                        value
  = value\n                    msg = msg + f\"[grey46]{key}[/][magenta3]:[/] [grey66]{value}[/]\\n\"\n
  \               return msg\n\n            def __rich__(self) -> Table:\n                from
  rich.table import Table\n\n                table = Table(title=f\"Feeds {len(self.config)}\")\n\n
  \               table.add_column(\"Feed\", justify=\"right\", style=\"cyan\", no_wrap=True)\n
  \               table.add_column(\"posts\", justify=\"left\", style=\"green\")\n
  \               table.add_column(\"config\", style=\"magenta\")\n\n                for
  name in self.config:\n                    table.add_row(\n                        name,\n
  \                       str(len(self[name].posts)),\n                        self._dict_panel(self.config[name].dict()),\n
  \                   )\n                return table\n```\n\n\n!! function <h2 id='config_model'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>config_model <em class='small'>function</em></h2>\n\n???+
  source \"config_model <em class='small'>source</em>\"\n\n```python\n\n        def
  config_model(markata: Markata) -> None:\n            markata.config_models.append(FeedsConfig)\n```\n\n\n!!
  function <h2 id='pre_render' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>pre_render <em class='small'>function</em></h2>\n    Create the Feeds object
  and attach it to markata.\n???+ source \"pre_render <em class='small'>source</em>\"\n\n```python\n\n
  \       def pre_render(markata: Markata) -> None:\n            \"\"\"\n            Create
  the Feeds object and attach it to markata.\n            \"\"\"\n            markata.feeds
  = Feeds(markata)\n```\n\n\n!! function <h2 id='save' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>save <em class='small'>function</em></h2>\n    Creates a new feed page for
  each page in the config.\n???+ source \"save <em class='small'>source</em>\"\n\n```python\n\n
  \       def save(markata: Markata) -> None:\n            \"\"\"\n            Creates
  a new feed page for each page in the config.\n            \"\"\"\n            with
  markata.cache as cache:\n                for feed in markata.feeds.values():\n                    create_page(\n
  \                       markata,\n                        feed,\n                        cache,\n
  \                   )\n\n            home = Path(str(markata.config.output_dir))
  / \"index.html\"\n            archive = Path(str(markata.config.output_dir)) / \"archive\"
  / \"index.html\"\n            if not home.exists() and archive.exists():\n                shutil.copy(str(archive),
  str(home))\n\n            xsl_template = get_template(feed.config.xsl_template)\n
  \           xsl = xsl_template.render(\n                markata=markata,\n                __version__=__version__,\n
  \               today=datetime.datetime.today(),\n                config=markata.config,\n
  \           )\n            xsl_file = Path(markata.config.output_dir) / \"rss.xsl\"\n
  \           xsl_file.write_text(xsl)\n```\n\n\n!! function <h2 id='get_template'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get_template <em class='small'>function</em></h2>\n\n???+
  source \"get_template <em class='small'>source</em>\"\n\n```python\n\n        def
  get_template(src) -> Template:\n            try:\n                return Template(Path(src).read_text(),
  undefined=SilentUndefined)\n            except FileNotFoundError:\n                return
  Template(src, undefined=SilentUndefined)\n            except OSError:  # File name
  too long\n                return Template(src, undefined=SilentUndefined)\n```\n\n\n!!
  function <h2 id='create_page' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>create_page <em class='small'>function</em></h2>\n    create an html unorderd
  list of posts.\n???+ source \"create_page <em class='small'>source</em>\"\n\n```python\n\n
  \       def create_page(\n            markata: Markata,\n            feed: Feed,\n
  \           cache,\n        ) -> None:\n            \"\"\"\n            create an
  html unorderd list of posts.\n            \"\"\"\n\n            posts = feed.posts\n\n
  \           cards = [\n                create_card(markata, post, feed.config.card_template,
  cache) for post in posts\n            ]\n            cards.insert(0, \"<ul>\")\n
  \           cards.append(\"</ul>\")\n            cards = \"\".join(cards)\n\n            template
  = get_template(feed.config.template)\n            rss_template = get_template(feed.config.rss_template)\n
  \           sitemap_template = get_template(feed.config.sitemap_template)\n            output_file
  = Path(markata.config.output_dir) / feed.config.slug / \"index.html\"\n            canonical_url
  = f\"{markata.config.url}/{feed.config.slug}/\"\n            output_file.parent.mkdir(exist_ok=True,
  parents=True)\n\n            rss_output_file = Path(markata.config.output_dir) /
  feed.config.slug / \"rss.xml\"\n            rss_output_file.parent.mkdir(exist_ok=True,
  parents=True)\n\n            sitemap_output_file = (\n                Path(markata.config.output_dir)
  / feed.config.slug / \"sitemap.xml\"\n            )\n            sitemap_output_file.parent.mkdir(exist_ok=True,
  parents=True)\n\n            key = markata.make_hash(\n                \"feeds\",\n
  \               template,\n                __version__,\n                cards,\n
  \               markata.config.url,\n                markata.config.description,\n
  \               feed.config.title,\n                canonical_url,\n                datetime.datetime.today(),\n
  \               markata.config,\n            )\n\n            feed_html_from_cache
  = markata.precache.get(key)\n            if feed_html_from_cache is None:\n                feed_html
  = template.render(\n                    markata=markata,\n                    __version__=__version__,\n
  \                   body=cards,\n                    url=markata.config.url,\n                    description=markata.config.description,\n
  \                   title=feed.config.title,\n                    canonical_url=canonical_url,\n
  \                   today=datetime.datetime.today(),\n                    config=markata.config,\n
  \               )\n                with markata.cache as cache:\n                    markata.cache.set(key,
  feed_html)\n\n            feed_rss = rss_template.render(markata=markata, feed=feed)\n
  \           feed_sitemap = sitemap_template.render(markata=markata, feed=feed)\n\n
  \           output_file.write_text(feed_html)\n            rss_output_file.write_text(feed_rss)\n
  \           sitemap_output_file.write_text(feed_sitemap)\n```\n\n\n!! function <h2
  id='create_card' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>create_card
  <em class='small'>function</em></h2>\n    Creates a card for one post based on the
  configured template.  If no\n    template is configured it will create one with
  the post title and dates\n    (if present).\n???+ source \"create_card <em class='small'>source</em>\"\n\n```python\n\n
  \       def create_card(\n            markata: \"Markata\",\n            post: \"Post\",\n
  \           template: Optional[str] = None,\n            cache=None,\n        )
  -> Any:\n            \"\"\"\n            Creates a card for one post based on the
  configured template.  If no\n            template is configured it will create one
  with the post title and dates\n            (if present).\n            \"\"\"\n            key
  = markata.make_hash(\"feeds\", template, str(post), post.content)\n\n            card
  = markata.precache.get(key)\n            if card is not None:\n                return
  card\n\n            if template is None:\n                template = markata.config.get(\"feeds_config\",
  {}).get(\"card_template\", None)\n\n            if template is None:\n                if
  \"date\" in post:\n                    card = textwrap.dedent(\n                        f\"\"\"\n
  \                       <li class='post'>\n                        <a href=\"/{markata.config.path_prefix}{post.slug}/\">\n
  \                           {post.title}\n                            {post.date.year}-\n
  \                           {post.date.month}-\n                            {post.date.day}\n
  \                       </a>\n                        </li>\n                        \"\"\",\n
  \                   )\n                else:\n                    card = textwrap.dedent(\n
  \                       f\"\"\"\n                        <li class='post'>\n                        <a
  href=\"/{markata.config.path_prefix}{post.slug}/\">\n                            {post.title}\n
  \                       </a>\n                        </li>\n                        \"\"\",\n
  \                   )\n            else:\n                try:\n                    _template
  = Template(Path(template).read_text())\n                except FileNotFoundError:\n
  \                   _template = Template(template)\n                except OSError:
  \ # File name too long\n                    _template = Template(template)\n                card
  = _template.render(post=post, **post.to_dict())\n            cache.add(key, card)\n
  \           return card\n```\n\n\n!! function <h2 id='cli' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>cli <em class='small'>function</em></h2>\n\n???+
  source \"cli <em class='small'>source</em>\"\n\n```python\n\n        def cli(app:
  typer.Typer, markata: \"Markata\") -> None:\n            feeds_app = typer.Typer()\n
  \           app.add_typer(feeds_app)\n\n            @feeds_app.callback()\n            def
  feeds():\n                \"feeds cli\"\n\n            @feeds_app.command()\n            def
  show() -> None:\n                markata.console.quiet = True\n                feeds
  = markata.feeds\n                markata.console.quiet = False\n                rich_print(feeds)\n```\n\n\n!!
  method <h2 id='_fail_with_undefined_error' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>_fail_with_undefined_error <em class='small'>method</em></h2>\n\n???+ source
  \"_fail_with_undefined_error <em class='small'>source</em>\"\n\n```python\n\n        def
  _fail_with_undefined_error(self, *args, **kwargs):\n                return \"\"\n```\n\n\n!!
  method <h2 id='default_name' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>default_name <em class='small'>method</em></h2>\n\n???+ source \"default_name
  <em class='small'>source</em>\"\n\n```python\n\n        def default_name(cls, v,
  *, values):\n                return v or str(values.get(\"slug\")).replace(\"-\",
  \"_\")\n```\n\n\n!! method <h2 id='read_template' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>read_template <em class='small'>method</em></h2>\n\n???+ source \"read_template
  <em class='small'>source</em>\"\n\n```python\n\n        def read_template(cls, v,
  *, values) -> str:\n                if isinstance(v, Path):\n                    return
  str(v.read_text())\n                return v\n```\n\n\n!! method <h2 id='name' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>name <em class='small'>method</em></h2>\n\n???+
  source \"name <em class='small'>source</em>\"\n\n```python\n\n        def name(self):\n
  \               return self.config.name\n```\n\n\n!! method <h2 id='posts' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>posts <em class='small'>method</em></h2>\n\n???+
  source \"posts <em class='small'>source</em>\"\n\n```python\n\n        def posts(self):\n
  \               return self.map(\"post\")\n```\n\n\n!! method <h2 id='map' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>map <em class='small'>method</em></h2>\n\n???+
  source \"map <em class='small'>source</em>\"\n\n```python\n\n        def map(self,
  func=\"post\", **args):\n                return self._m.map(func, **{**self.config.dict(),
  **args})\n```\n\n\n!! method <h2 id='__init__' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>__init__ <em class='small'>method</em></h2>\n\n???+ source \"__init__ <em
  class='small'>source</em>\"\n\n```python\n\n        def __init__(self, markata:
  Markata) -> None:\n                self._m = markata\n                self.config
  = {f.name: f for f in markata.config.feeds}\n                self.refresh()\n```\n\n\n!!
  method <h2 id='refresh' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>refresh
  <em class='small'>method</em></h2>\n    Refresh all of the feeds objects\n???+ source
  \"refresh <em class='small'>source</em>\"\n\n```python\n\n        def refresh(self)
  -> None:\n                \"\"\"\n                Refresh all of the feeds objects\n
  \               \"\"\"\n                for feed in self._m.config.feeds:\n                    feed
  = Feed(config=feed, _m=self._m)\n                    self.__setattr__(feed.name,
  feed)\n```\n\n\n!! method <h2 id='__iter__' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>__iter__ <em class='small'>method</em></h2>\n\n???+ source \"__iter__ <em
  class='small'>source</em>\"\n\n```python\n\n        def __iter__(self):\n                return
  iter(self.config.keys())\n```\n\n\n!! method <h2 id='keys' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>keys <em class='small'>method</em></h2>\n\n???+
  source \"keys <em class='small'>source</em>\"\n\n```python\n\n        def keys(self):\n
  \               return iter(self.config.keys())\n```\n\n\n!! method <h2 id='values'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>values <em class='small'>method</em></h2>\n\n???+
  source \"values <em class='small'>source</em>\"\n\n```python\n\n        def values(self):\n
  \               return [self[feed] for feed in self.config.keys()]\n```\n\n\n!!
  method <h2 id='items' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>items
  <em class='small'>method</em></h2>\n\n???+ source \"items <em class='small'>source</em>\"\n\n```python\n\n
  \       def items(self):\n                return [(key, self[key]) for key in self.config]\n```\n\n\n!!
  method <h2 id='__getitem__' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>__getitem__ <em class='small'>method</em></h2>\n\n???+ source \"__getitem__
  <em class='small'>source</em>\"\n\n```python\n\n        def __getitem__(self, key:
  str) -> Any:\n                return getattr(self, key.replace(\"-\", \"_\").lower())\n```\n\n\n!!
  method <h2 id='_dict_panel' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>_dict_panel <em class='small'>method</em></h2>\n    pretty print configs
  with rich\n???+ source \"_dict_panel <em class='small'>source</em>\"\n\n```python\n\n
  \       def _dict_panel(self, config) -> str:\n                \"\"\"\n                pretty
  print configs with rich\n                \"\"\"\n                msg = \"\"\n                for
  key, value in config.items():\n                    if isinstance(value, str):\n
  \                       if len(value) > 50:\n                            value =
  value[:50] + \"...\"\n                        value = value\n                    msg
  = msg + f\"[grey46]{key}[/][magenta3]:[/] [grey66]{value}[/]\\n\"\n                return
  msg\n```\n\n\n!! method <h2 id='__rich__' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>__rich__ <em class='small'>method</em></h2>\n\n???+ source \"__rich__ <em
  class='small'>source</em>\"\n\n```python\n\n        def __rich__(self) -> Table:\n
  \               from rich.table import Table\n\n                table = Table(title=f\"Feeds
  {len(self.config)}\")\n\n                table.add_column(\"Feed\", justify=\"right\",
  style=\"cyan\", no_wrap=True)\n                table.add_column(\"posts\", justify=\"left\",
  style=\"green\")\n                table.add_column(\"config\", style=\"magenta\")\n\n
  \               for name in self.config:\n                    table.add_row(\n                        name,\n
  \                       str(len(self[name].posts)),\n                        self._dict_panel(self.config[name].dict()),\n
  \                   )\n                return table\n```\n\n\n!! function <h2 id='feeds'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>feeds <em class='small'>function</em></h2>\n
  \   feeds cli\n???+ source \"feeds <em class='small'>source</em>\"\n\n```python\n\n
  \       def feeds():\n                \"feeds cli\"\n```\n\n\n!! function <h2 id='show'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>show <em class='small'>function</em></h2>\n\n???+
  source \"show <em class='small'>source</em>\"\n\n```python\n\n        def show()
  -> None:\n                markata.console.quiet = True\n                feeds =
  markata.feeds\n                markata.console.quiet = False\n                rich_print(feeds)\n```\n\n"
date: 0001-01-01
description: The  This plugin is built-in and enabled by default, but in you want
  to be very At the root of the markata.feeds config you may set  Underneath of the  The
  foll
html: "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n<title>Feeds.Py</title>\n<meta
  charset=\"utf-8\"/>\n<meta content=\"width=device-width, initial-scale=1\" name=\"viewport\"/>\n<meta
  content=\"The  This plugin is built-in and enabled by default, but in you want to
  be very At the root of the markata.feeds config you may set  Underneath of the  The
  foll\" name=\"description\"/>\n<link href=\"/static/favicon.ico\" rel=\"icon\" type=\"image/png\"/>\n<script>\n
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
  name=\"og:type\" property=\"og:type\"/><meta content=\"The  This plugin is built-in
  and enabled by default, but in you want to be very At the root of the markata.feeds
  config you may set  Underneath of the  The foll\" name=\"description\" property=\"description\"/><meta
  content=\"The  This plugin is built-in and enabled by default, but in you want to
  be very At the root of the markata.feeds config you may set  Underneath of the  The
  foll\" name=\"og:description\" property=\"og:description\"/><meta content=\"The
  \ This plugin is built-in and enabled by default, but in you want to be very At
  the root of the markata.feeds config you may set  Underneath of the  The foll\"
  name=\"twitter:description\" property=\"twitter:description\"/><meta content=\"Feeds.Py
  | Markata\" name=\"og:title\" property=\"og:title\"/><meta content=\"Feeds.Py |
  Markata\" name=\"twitter:title\" property=\"twitter:title\"/><meta content=\"https://markata.dev//markata/plugins/feeds-og.png\"
  name=\"og:image\" property=\"og:image\"/><meta content=\"https://markata.dev//markata/plugins/feeds-og.png\"
  name=\"twitter:image\" property=\"twitter:image\"/><meta content=\"1600\" name=\"og:image:width\"
  property=\"og:image:width\"/><meta content=\"900\" name=\"og:image:width\" property=\"og:image:width\"/><meta
  content=\"summary_large_image\" name=\"twitter:card\" property=\"twitter:card\"/><meta
  content=\"Markata\" name=\"og:site_name\" property=\"og:site_name\"/><meta content=\"@_waylonwalker\"
  name=\"twitter:creator\" property=\"twitter:creator\"/><meta content=\"Feeds.Py\"
  name=\"title\" property=\"title\"/><meta content=\"markata 0.8.0.dev6\" name=\"generator\"
  property=\"generator\"/><link href=\"https://markata.dev//markata/plugins/feeds/\"
  rel=\"canonical\"/><meta content=\"https://markata.dev//markata/plugins/feeds/\"
  name=\"og:url\" property=\"og:url\"/></head>\n<body><nav>\n<a href=\"/\">markata</a>\n<a
  href=\"https://github.com/WaylonWalker/markata\">GitHub</a>\n</nav>\n<div>\n<label
  class=\"theme-switch\" for=\"checkbox-theme\" id=\"theme-switch\" title=\"light/dark
  mode toggle\">\n<input id=\"checkbox-theme\" type=\"checkbox\"/>\n<div class=\"slider
  round\"></div>\n</label>\n</div>\n<section class=\"title\">\n<h1 id=\"title\">\n
  \           Feeds.Py \n            \n        </h1>\n</section>\n<main><p>The <code>markata.plugins.feeds</code>
  plugin is used to create feed pages, which are lists of\nposts.  The list is generated
  using a <code>filter</code>, then each post in the list is\nrendered with a <code>card_template</code>
  before being applied to the <code>body</code> of the\n<code>template</code>.</p>\n<h1
  id=\"installation\">Installation <a class=\"header-anchor\" href=\"#installation\"><svg
  aria-hidden=\"true\" class=\"heading-permalink\" fill=\"currentColor\" focusable=\"false\"
  height=\"1em\" viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
  d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>This plugin is built-in
  and enabled by default, but in you want to be very\nexplicit you can add it to your
  list of existing plugins.</p>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
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
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"n\">hooks</span><span
  class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"p\">[</span>\n<span
  class=\"w\">   </span><span class=\"s2\">\"markata.plugins.feeds\"</span><span class=\"p\">,</span>\n<span
  class=\"w\">   </span><span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n<h1 id=\"configuration\">Configuration
  <a class=\"header-anchor\" href=\"#configuration\"><svg aria-hidden=\"true\" class=\"heading-permalink\"
  fill=\"currentColor\" focusable=\"false\" height=\"1em\" viewbox=\"0 0 24 24\" width=\"1em\"
  xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949
  2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992 5.992 0 0 0 1.695-4.285 5.976
  5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242 6.003 6.003 0 0 0-1.905 1.24l-1.731
  1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985 3.985 0 0 1 2.761-1.123 3.975
  3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005 3.006a3.982 3.982 0
  0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201 13.6zm5.602-3.198a5.99
  5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995 2.994a5.992 5.992 0
  0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0 0 6.431 1.242 6.003
  6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836 19.81a3.985
  3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0 0 1-.111-5.644l3.005-3.006a3.982
  3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632 1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<h1
  id=\"set-default-template-and-card_template\">set default template and card_template
  <a class=\"header-anchor\" href=\"#set-default-template-and-card_template\"><svg
  aria-hidden=\"true\" class=\"heading-permalink\" fill=\"currentColor\" focusable=\"false\"
  height=\"1em\" viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
  d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>At the root of the markata.feeds
  config you may set <code>template</code>, and\n<code>card_template</code>.  These
  will become your defaults for every feed you create.\nIf you do not set these, markata
  will use it's defaults.  The defaults are\ndesigned to work for a variety of use
  cases, but are not likely the best for\nall.</p>\n<pre class=\"wrapper\">\n\n<div
  class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
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
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata.feeds_config]</span>\n<span
  class=\"n\">template</span><span class=\"o\">=</span><span class=\"s2\">\"pages/templates/archive_template.html\"</span>\n<span
  class=\"n\">card_template</span><span class=\"o\">=</span><span class=\"s2\">\"plugins/feed_card_template.html\"</span>\n</pre></div>\n\n</pre>\n<h1
  id=\"pages\">pages <a class=\"header-anchor\" href=\"#pages\"><svg aria-hidden=\"true\"
  class=\"heading-permalink\" fill=\"currentColor\" focusable=\"false\" height=\"1em\"
  viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M9.199
  13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>Underneath of the <code>markata.feeds</code>
  we will create a new map for each page where\nthe name of the map will be the name
  of the page.</p>\n<p>The following config will create a page at <code>/all-posts</code>
  that inclues every\nsingle post.</p>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
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
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.feeds]]</span>\n<span
  class=\"n\">title</span><span class=\"o\">=</span><span class=\"s2\">\"All Posts\"</span>\n<span
  class=\"n\">slug</span><span class=\"o\">=</span><span class=\"s1\">'all'</span>\n<span
  class=\"n\">filter</span><span class=\"o\">=</span><span class=\"s2\">\"True\"</span>\n</pre></div>\n\n</pre>\n<h1
  id=\"template\">template <a class=\"header-anchor\" href=\"#template\"><svg aria-hidden=\"true\"
  class=\"heading-permalink\" fill=\"currentColor\" focusable=\"false\" height=\"1em\"
  viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M9.199
  13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>The <code>template</code>
  configuration key is a file path to the template that you want\nto use to create
  the feed.  You may set the default template you want to use\nfor all feeds under
  <code>[markata.feeds]</code>, as well as override it inside of each\nfeeds config.</p>\n<p>The
  template is a jinja style template that expects to fill in a <code>title</code>
  and\n<code>body</code> variable.</p>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
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
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"cp\">&lt;!DOCTYPE
  html&gt;</span>\n<span class=\"p\">&lt;</span><span class=\"nt\">html</span> <span
  class=\"na\">lang</span><span class=\"o\">=</span><span class=\"s\">\"en\"</span><span
  class=\"p\">&gt;</span>\n  <span class=\"p\">&lt;</span><span class=\"nt\">head</span><span
  class=\"p\">&gt;</span>\n    <span class=\"p\">&lt;</span><span class=\"nt\">title</span><span
  class=\"p\">&gt;</span>{{ title }}<span class=\"p\">&lt;/</span><span class=\"nt\">title</span><span
  class=\"p\">&gt;</span>\n  <span class=\"p\">&lt;/</span><span class=\"nt\">head</span><span
  class=\"p\">&gt;</span>\n  <span class=\"p\">&lt;</span><span class=\"nt\">body</span><span
  class=\"p\">&gt;</span>\n    <span class=\"p\">&lt;</span><span class=\"nt\">ul</span><span
  class=\"p\">&gt;</span>\n        {{ body }}\n    <span class=\"p\">&lt;/</span><span
  class=\"nt\">ul</span><span class=\"p\">&gt;</span>\n  <span class=\"p\">&lt;/</span><span
  class=\"nt\">body</span><span class=\"p\">&gt;</span>\n<span class=\"p\">&lt;/</span><span
  class=\"nt\">html</span><span class=\"p\">&gt;</span>\n</pre></div>\n\n</pre>\n<div
  class=\"admonition note\">\n<p class=\"admonition-title\">Note</p>\n<p>I highly
  reccomend putting your <code>body</code> in a <code>&lt;ul&gt;</code>, and wrapping
  your\n<code>card_template</code>s in an <code>&lt;li&gt;</code>.</p>\n</div>\n<h1
  id=\"card_template\">card_template <a class=\"header-anchor\" href=\"#card_template\"><svg
  aria-hidden=\"true\" class=\"heading-permalink\" fill=\"currentColor\" focusable=\"false\"
  height=\"1em\" viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
  d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>All keys available from
  each post is available to put into your jinja\ntemplate.  These can either be placed
  there in your post frontmatter, or\nthrough a plugin that automatically adds to
  the post before the save phase.</p>\n<p>Here is a very simple example that would
  give a link to each post with the\ntitle and date.</p>\n<pre class=\"wrapper\">\n\n<div
  class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
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
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.feeds]]</span>\n<span
  class=\"n\">slug</span><span class=\"o\">=</span><span class=\"s1\">'all'</span>\n<span
  class=\"n\">title</span><span class=\"o\">=</span><span class=\"s1\">'All Posts'</span>\n<span
  class=\"n\">filter</span><span class=\"o\">=</span><span class=\"s2\">\"True\"</span>\n<span
  class=\"n\">card_template</span><span class=\"o\">=</span><span class=\"s1\">'''</span>\n<span
  class=\"s1\">&lt;li&gt;</span>\n<span class=\"s1\">    &lt;a href={{markata.config.get('path_prefix',
  '')}}{{slug}}&gt;</span>\n<span class=\"s1\">        {{title}}-{{date}}</span>\n<span
  class=\"s1\">    &lt;/a&gt;</span>\n<span class=\"s1\">&lt;/li&gt;</span>\n<span
  class=\"s1\">'''</span>\n</pre></div>\n\n</pre>\n<h1 id=\"filter\">filter <a class=\"header-anchor\"
  href=\"#filter\"><svg aria-hidden=\"true\" class=\"heading-permalink\" fill=\"currentColor\"
  focusable=\"false\" height=\"1em\" viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
  d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>The filter is a python
  expression ran on every post that expects to return a\nboolean.  The variables available
  to this expression are every key in your\nfrontmatter, plus the <code>timedelta</code>
  function, and <code>parse</code> function to more easily\nwork with dates.</p>\n<h1
  id=\"feed-examples\">Feed Examples <a class=\"header-anchor\" href=\"#feed-examples\"><svg
  aria-hidden=\"true\" class=\"heading-permalink\" fill=\"currentColor\" focusable=\"false\"
  height=\"1em\" viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
  d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>True can be passed in
  to make a feed of all the posts you have.</p>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
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
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.feeds]]</span>\n<span
  class=\"n\">slug</span><span class=\"o\">=</span><span class=\"s1\">'all'</span>\n<span
  class=\"n\">title</span><span class=\"o\">=</span><span class=\"s1\">'All Posts'</span>\n<span
  class=\"n\">filter</span><span class=\"o\">=</span><span class=\"s2\">\"True\"</span>\n</pre></div>\n\n</pre>\n<p>You
  can compare against the values of the keys from your frontmatter.  This\nexample
  creates a feed that includes every post where published is <code>True</code>.</p>\n<pre
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
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.feeds]]</span>\n<span
  class=\"n\">slug</span><span class=\"o\">=</span><span class=\"s1\">'draft'</span>\n<span
  class=\"n\">title</span><span class=\"o\">=</span><span class=\"s1\">'Draft'</span>\n<span
  class=\"n\">filter</span><span class=\"o\">=</span><span class=\"s2\">\"published=='False'\"</span>\n</pre></div>\n\n</pre>\n<p>We
  can also compare against dates.  The\n<a href=\"https://markata.dev/markata/plugins/datetime/\">markata.plugins.datetime</a>\nplugin,
  automatically adds <code>today</code> as today's date and <code>now</code> as the
  current\ndatetime.  These are quite handy to create feeds for scheduled, recent,
  or\ntoday's posts.  The following two examples will create a feed for scheduled\nposts
  and for today's posts respectively.</p>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
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
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.feeds]]</span>\n<span
  class=\"n\">slug</span><span class=\"o\">=</span><span class=\"s1\">'scheduled'</span>\n<span
  class=\"n\">title</span><span class=\"o\">=</span><span class=\"s1\">'Scheduled'</span>\n<span
  class=\"n\">filter</span><span class=\"o\">=</span><span class=\"s2\">\"date&gt;today\"</span>\n\n<span
  class=\"k\">[[markata.feeds]]</span>\n<span class=\"n\">slug</span><span class=\"o\">=</span><span
  class=\"s1\">'today'</span>\n<span class=\"n\">title</span><span class=\"o\">=</span><span
  class=\"s1\">'Today'</span>\n<span class=\"n\">filter</span><span class=\"o\">=</span><span
  class=\"s2\">\"date==today\"</span>\n</pre></div>\n\n</pre>\n<p>If you have list
  of items in your frontmatter for something like <code>tags</code>, you\ncan check
  for the existence of a tag in the list.</p>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
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
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.feeds]]</span>\n<span
  class=\"n\">slug</span><span class=\"o\">=</span><span class=\"s1\">'python'</span>\n<span
  class=\"n\">title</span><span class=\"o\">=</span><span class=\"s1\">'Python'</span>\n<span
  class=\"n\">filter</span><span class=\"o\">=</span><span class=\"s2\">\"date&lt;=today
  and 'python' in tags\"</span>\n</pre></div>\n\n</pre>\n<p>And of course you can
  combine all the things into larger expressions.  Here is\none example of the main
  feed on my blog.</p>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
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
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.feeds]]</span>\n<span
  class=\"n\">slug</span><span class=\"o\">=</span><span class=\"s1\">'blog'</span>\n<span
  class=\"n\">title</span><span class=\"o\">=</span><span class=\"s1\">'Blog'</span>\n<span
  class=\"n\">filter</span><span class=\"o\">=</span><span class=\"s2\">\"date&lt;=today
  and templateKey in ['blog-post'] and published =='True'\"</span>\n</pre></div>\n\n</pre>\n<p>Here
  is another example that shows my drafts for a particular tag.</p>\n<pre class=\"wrapper\">\n\n<div
  class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
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
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.feeds]]</span>\n<span
  class=\"n\">slug</span><span class=\"o\">=</span><span class=\"s1\">'python-draft'</span>\n<span
  class=\"n\">title</span><span class=\"o\">=</span><span class=\"s1\">'Python Draft'</span>\n<span
  class=\"n\">filter</span><span class=\"o\">=</span><span class=\"s2\">\"date&lt;=today
  and 'python' in tags and published=='False'\"</span>\n</pre></div>\n\n</pre>\n<h1
  id=\"defaults\">Defaults <a class=\"header-anchor\" href=\"#defaults\"><svg aria-hidden=\"true\"
  class=\"heading-permalink\" fill=\"currentColor\" focusable=\"false\" height=\"1em\"
  viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M9.199
  13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>By default feeds will
  create one feed page at <code>/archive/</code> that includes all\nposts.</p>\n<p><a
  class=\"wikilink\" href=\"/markata.feeds\">markata.feeds</a>\nslug='archive'\ntitle='All
  Posts'\nfilter=\"True\"</p>\n<p>!! class </p><h2 class=\"admonition-title\" id=\"SilentUndefined\"
  style=\"margin:0;padding:.5rem 1rem;\">SilentUndefined <em class=\"small\">class</em></h2>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"SilentUndefined
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
  <span class=\"nc\">SilentUndefined</span><span class=\"p\">(</span><span class=\"n\">Undefined</span><span
  class=\"p\">):</span>\n            <span class=\"k\">def</span> <span class=\"nf\">_fail_with_undefined_error</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
  class=\"o\">*</span><span class=\"n\">args</span><span class=\"p\">,</span> <span
  class=\"o\">**</span><span class=\"n\">kwargs</span><span class=\"p\">):</span>\n
  \               <span class=\"k\">return</span> <span class=\"s2\">\"\"</span>\n</pre></div>\n\n</pre>\n<p>!!
  class </p><h2 class=\"admonition-title\" id=\"MarkataFilterError\" style=\"margin:0;padding:.5rem
  1rem;\">MarkataFilterError <em class=\"small\">class</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"MarkataFilterError
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
  <span class=\"nc\">MarkataFilterError</span><span class=\"p\">(</span><span class=\"ne\">RuntimeError</span><span
  class=\"p\">):</span>\n            <span class=\"o\">...</span>\n</pre></div>\n\n</pre>\n<p>!!
  class </p><h2 class=\"admonition-title\" id=\"FeedConfig\" style=\"margin:0;padding:.5rem
  1rem;\">FeedConfig <em class=\"small\">class</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"FeedConfig
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
  <span class=\"nc\">FeedConfig</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
  class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
  \           <span class=\"n\">DEFAULT_TITLE</span><span class=\"p\">:</span> <span
  class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">\"All Posts\"</span>\n\n
  \           <span class=\"n\">title</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
  <span class=\"o\">=</span> <span class=\"n\">DEFAULT_TITLE</span>\n            <span
  class=\"n\">slug</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
  <span class=\"o\">=</span> <span class=\"kc\">None</span>\n            <span class=\"n\">name</span><span
  class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
  class=\"nb\">str</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
  class=\"kc\">None</span>\n            <span class=\"nb\">filter</span><span class=\"p\">:</span>
  <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">\"True\"</span>\n
  \           <span class=\"n\">sort</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
  <span class=\"o\">=</span> <span class=\"s2\">\"date\"</span>\n            <span
  class=\"n\">reverse</span><span class=\"p\">:</span> <span class=\"nb\">bool</span>
  <span class=\"o\">=</span> <span class=\"kc\">False</span>\n            <span class=\"n\">rss</span><span
  class=\"p\">:</span> <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span
  class=\"kc\">True</span>\n            <span class=\"n\">sitemap</span><span class=\"p\">:</span>
  <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"kc\">True</span>\n
  \           <span class=\"n\">card_template</span><span class=\"p\">:</span> <span
  class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">\"\"\"</span>\n<span
  class=\"s2\">                &lt;li class='post'&gt;</span>\n<span class=\"s2\">
  \                   &lt;a href=\"/{{ markata.config.path_prefix }}{{ post.slug }}/\"&gt;</span>\n<span
  class=\"s2\">                        {{ post.title }}</span>\n<span class=\"s2\">
  \                   &lt;/a&gt;</span>\n<span class=\"s2\">                &lt;/li&gt;</span>\n<span
  class=\"s2\">                \"\"\"</span>\n            <span class=\"n\">template</span><span
  class=\"p\">:</span> <span class=\"nb\">str</span> <span class=\"o\">=</span> <span
  class=\"n\">Path</span><span class=\"p\">(</span><span class=\"vm\">__file__</span><span
  class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">parent</span> <span
  class=\"o\">/</span> <span class=\"s2\">\"default_post_template.html.jinja\"</span>\n
  \           <span class=\"n\">rss_template</span><span class=\"p\">:</span> <span
  class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
  class=\"p\">(</span><span class=\"vm\">__file__</span><span class=\"p\">)</span><span
  class=\"o\">.</span><span class=\"n\">parent</span> <span class=\"o\">/</span> <span
  class=\"s2\">\"default_rss_template.xml\"</span>\n            <span class=\"n\">sitemap_template</span><span
  class=\"p\">:</span> <span class=\"nb\">str</span> <span class=\"o\">=</span> <span
  class=\"n\">Path</span><span class=\"p\">(</span><span class=\"vm\">__file__</span><span
  class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">parent</span> <span
  class=\"o\">/</span> <span class=\"s2\">\"default_sitemap_template.xml\"</span>\n
  \           <span class=\"n\">xsl_template</span><span class=\"p\">:</span> <span
  class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
  class=\"p\">(</span><span class=\"vm\">__file__</span><span class=\"p\">)</span><span
  class=\"o\">.</span><span class=\"n\">parent</span> <span class=\"o\">/</span> <span
  class=\"s2\">\"default_xsl_template.xsl\"</span>\n\n            <span class=\"nd\">@pydantic</span><span
  class=\"o\">.</span><span class=\"n\">validator</span><span class=\"p\">(</span><span
  class=\"s2\">\"name\"</span><span class=\"p\">,</span> <span class=\"n\">pre</span><span
  class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
  class=\"n\">always</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
  class=\"p\">)</span>\n            <span class=\"k\">def</span> <span class=\"nf\">default_name</span><span
  class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
  class=\"n\">v</span><span class=\"p\">,</span> <span class=\"o\">*</span><span class=\"p\">,</span>
  <span class=\"n\">values</span><span class=\"p\">):</span>\n                <span
  class=\"k\">return</span> <span class=\"n\">v</span> <span class=\"ow\">or</span>
  <span class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">values</span><span
  class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
  class=\"s2\">\"slug\"</span><span class=\"p\">))</span><span class=\"o\">.</span><span
  class=\"n\">replace</span><span class=\"p\">(</span><span class=\"s2\">\"-\"</span><span
  class=\"p\">,</span> <span class=\"s2\">\"_\"</span><span class=\"p\">)</span>\n\n
  \           <span class=\"nd\">@pydantic</span><span class=\"o\">.</span><span class=\"n\">validator</span><span
  class=\"p\">(</span><span class=\"s2\">\"card_template\"</span><span class=\"p\">,</span>
  <span class=\"s2\">\"template\"</span><span class=\"p\">,</span> <span class=\"n\">pre</span><span
  class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
  class=\"n\">always</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
  class=\"p\">)</span>\n            <span class=\"k\">def</span> <span class=\"nf\">read_template</span><span
  class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
  class=\"n\">v</span><span class=\"p\">,</span> <span class=\"o\">*</span><span class=\"p\">,</span>
  <span class=\"n\">values</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"nb\">str</span><span class=\"p\">:</span>\n                <span class=\"k\">if</span>
  <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">v</span><span
  class=\"p\">,</span> <span class=\"n\">Path</span><span class=\"p\">):</span>\n
  \                   <span class=\"k\">return</span> <span class=\"nb\">str</span><span
  class=\"p\">(</span><span class=\"n\">v</span><span class=\"o\">.</span><span class=\"n\">read_text</span><span
  class=\"p\">())</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n<p>!!
  class </p><h2 class=\"admonition-title\" id=\"FeedsConfig\" style=\"margin:0;padding:.5rem
  1rem;\">FeedsConfig <em class=\"small\">class</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"FeedsConfig
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
  <span class=\"nc\">FeedsConfig</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
  class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
  \           <span class=\"n\">feeds</span><span class=\"p\">:</span> <span class=\"n\">List</span><span
  class=\"p\">[</span><span class=\"n\">FeedConfig</span><span class=\"p\">]</span>
  <span class=\"o\">=</span> <span class=\"p\">[</span><span class=\"n\">FeedConfig</span><span
  class=\"p\">(</span><span class=\"n\">slug</span><span class=\"o\">=</span><span
  class=\"s2\">\"archive\"</span><span class=\"p\">)]</span>\n</pre></div>\n\n</pre>\n<p>!!
  class </p><h2 class=\"admonition-title\" id=\"Feed\" style=\"margin:0;padding:.5rem
  1rem;\">Feed <em class=\"small\">class</em></h2>\nA storage class for markata feed
  objects.\n<pre><code># Usage\n\n``` python\nfrom markata import Markata\nm = Markata()\n\n#
  access posts for a feed\nm.feeds.docs.posts\n\n# access config for a feed\nm.feeds.docs.config\n```\n</code></pre>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"Feed
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
  <span class=\"nc\">Feed</span><span class=\"p\">:</span>\n<span class=\"w\">            </span><span
  class=\"sd\">\"\"\"</span>\n<span class=\"sd\">            A storage class for markata
  feed objects.</span>\n\n<span class=\"sd\">            # Usage</span>\n\n<span class=\"sd\">
  \           ``` python</span>\n<span class=\"sd\">            from markata import
  Markata</span>\n<span class=\"sd\">            m = Markata()</span>\n\n<span class=\"sd\">
  \           # access posts for a feed</span>\n<span class=\"sd\">            m.feeds.docs.posts</span>\n\n<span
  class=\"sd\">            # access config for a feed</span>\n<span class=\"sd\">
  \           m.feeds.docs.config</span>\n<span class=\"sd\">            ```</span>\n<span
  class=\"sd\">            \"\"\"</span>\n\n            <span class=\"n\">config</span><span
  class=\"p\">:</span> <span class=\"n\">FeedConfig</span>\n            <span class=\"n\">_m</span><span
  class=\"p\">:</span> <span class=\"n\">Markata</span>\n\n            <span class=\"nd\">@property</span>\n
  \           <span class=\"k\">def</span> <span class=\"nf\">name</span><span class=\"p\">(</span><span
  class=\"bp\">self</span><span class=\"p\">):</span>\n                <span class=\"k\">return</span>
  <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">name</span>\n\n            <span class=\"nd\">@property</span>\n
  \           <span class=\"k\">def</span> <span class=\"nf\">posts</span><span class=\"p\">(</span><span
  class=\"bp\">self</span><span class=\"p\">):</span>\n                <span class=\"k\">return</span>
  <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">map</span><span
  class=\"p\">(</span><span class=\"s2\">\"post\"</span><span class=\"p\">)</span>\n\n
  \           <span class=\"k\">def</span> <span class=\"nf\">map</span><span class=\"p\">(</span><span
  class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"n\">func</span><span
  class=\"o\">=</span><span class=\"s2\">\"post\"</span><span class=\"p\">,</span>
  <span class=\"o\">**</span><span class=\"n\">args</span><span class=\"p\">):</span>\n
  \               <span class=\"k\">return</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">_m</span><span class=\"o\">.</span><span class=\"n\">map</span><span
  class=\"p\">(</span><span class=\"n\">func</span><span class=\"p\">,</span> <span
  class=\"o\">**</span><span class=\"p\">{</span><span class=\"o\">**</span><span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">dict</span><span class=\"p\">(),</span> <span
  class=\"o\">**</span><span class=\"n\">args</span><span class=\"p\">})</span>\n</pre></div>\n\n</pre>\n<p>!!
  class </p><h2 class=\"admonition-title\" id=\"Feeds\" style=\"margin:0;padding:.5rem
  1rem;\">Feeds <em class=\"small\">class</em></h2>\nA storage class for all markata
  Feed objects\n<pre><code>``` python\nfrom markata import Markata\nm = Markata()\n\nm.feeds\n\n#
  access all config\nm.feeds.config\n\n# refresh list of posts in all feeds\nm.feeds.refresh()\n\n\n#
  iterating over feeds gives the name of the feed\nfor k in m.feeds:\n     print(k)\n\n#
  project-gallery\n# docs\n# autodoc\n# core_modules\n# plugins\n# archive\n\n# iterate
  over items like keys and values in a dict, items returns name of\n# feed and a feed
  object\nfor k, v in m.feeds.items():\n    print(k, len(v.posts))\n\n# project-gallery
  2\n# docs 6\n# autodoc 65\n# core_modules 26\n# plugins 39\n# archive 65\n\n# values
  can be iterated over in just the same way\nfor v in m.feeds.values():\n     print(len(v.posts))\n#
  2\n# 6\n# 65\n# 26\n# 39\n# 65\n```\n\nAccessing feeds can be done using square
  brackets or dot notation.\n\n``` python\nfrom markata import Markata\nm = Markata()\n\n#
  both of these will return the `docs` Feed object.\nm.feeds.docs\nm['docs']\n```\n</code></pre>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"Feeds
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
  <span class=\"nc\">Feeds</span><span class=\"p\">:</span>\n<span class=\"w\">            </span><span
  class=\"sd\">\"\"\"</span>\n<span class=\"sd\">            A storage class for all
  markata Feed objects</span>\n\n<span class=\"sd\">            ``` python</span>\n<span
  class=\"sd\">            from markata import Markata</span>\n<span class=\"sd\">
  \           m = Markata()</span>\n\n<span class=\"sd\">            m.feeds</span>\n\n<span
  class=\"sd\">            # access all config</span>\n<span class=\"sd\">            m.feeds.config</span>\n\n<span
  class=\"sd\">            # refresh list of posts in all feeds</span>\n<span class=\"sd\">
  \           m.feeds.refresh()</span>\n\n\n<span class=\"sd\">            # iterating
  over feeds gives the name of the feed</span>\n<span class=\"sd\">            for
  k in m.feeds:</span>\n<span class=\"sd\">                 print(k)</span>\n\n<span
  class=\"sd\">            # project-gallery</span>\n<span class=\"sd\">            #
  docs</span>\n<span class=\"sd\">            # autodoc</span>\n<span class=\"sd\">
  \           # core_modules</span>\n<span class=\"sd\">            # plugins</span>\n<span
  class=\"sd\">            # archive</span>\n\n<span class=\"sd\">            # iterate
  over items like keys and values in a dict, items returns name of</span>\n<span class=\"sd\">
  \           # feed and a feed object</span>\n<span class=\"sd\">            for
  k, v in m.feeds.items():</span>\n<span class=\"sd\">                print(k, len(v.posts))</span>\n\n<span
  class=\"sd\">            # project-gallery 2</span>\n<span class=\"sd\">            #
  docs 6</span>\n<span class=\"sd\">            # autodoc 65</span>\n<span class=\"sd\">
  \           # core_modules 26</span>\n<span class=\"sd\">            # plugins 39</span>\n<span
  class=\"sd\">            # archive 65</span>\n\n<span class=\"sd\">            #
  values can be iterated over in just the same way</span>\n<span class=\"sd\">            for
  v in m.feeds.values():</span>\n<span class=\"sd\">                 print(len(v.posts))</span>\n<span
  class=\"sd\">            # 2</span>\n<span class=\"sd\">            # 6</span>\n<span
  class=\"sd\">            # 65</span>\n<span class=\"sd\">            # 26</span>\n<span
  class=\"sd\">            # 39</span>\n<span class=\"sd\">            # 65</span>\n<span
  class=\"sd\">            ```</span>\n\n<span class=\"sd\">            Accessing
  feeds can be done using square brackets or dot notation.</span>\n\n<span class=\"sd\">
  \           ``` python</span>\n<span class=\"sd\">            from markata import
  Markata</span>\n<span class=\"sd\">            m = Markata()</span>\n\n<span class=\"sd\">
  \           # both of these will return the `docs` Feed object.</span>\n<span class=\"sd\">
  \           m.feeds.docs</span>\n<span class=\"sd\">            m['docs']</span>\n<span
  class=\"sd\">            ```</span>\n<span class=\"sd\">            \"\"\"</span>\n\n
  \           <span class=\"k\">def</span> <span class=\"fm\">__init__</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
  class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"n\">Markata</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">_m</span> <span class=\"o\">=</span> <span class=\"n\">markata</span>\n
  \               <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span>
  <span class=\"o\">=</span> <span class=\"p\">{</span><span class=\"n\">f</span><span
  class=\"o\">.</span><span class=\"n\">name</span><span class=\"p\">:</span> <span
  class=\"n\">f</span> <span class=\"k\">for</span> <span class=\"n\">f</span> <span
  class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">feeds</span><span
  class=\"p\">}</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">refresh</span><span class=\"p\">()</span>\n\n            <span class=\"k\">def</span>
  <span class=\"nf\">refresh</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n<span class=\"w\">                </span><span class=\"sd\">\"\"\"</span>\n<span
  class=\"sd\">                Refresh all of the feeds objects</span>\n<span class=\"sd\">
  \               \"\"\"</span>\n                <span class=\"k\">for</span> <span
  class=\"n\">feed</span> <span class=\"ow\">in</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">_m</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">feeds</span><span class=\"p\">:</span>\n                    <span
  class=\"n\">feed</span> <span class=\"o\">=</span> <span class=\"n\">Feed</span><span
  class=\"p\">(</span><span class=\"n\">config</span><span class=\"o\">=</span><span
  class=\"n\">feed</span><span class=\"p\">,</span> <span class=\"n\">_m</span><span
  class=\"o\">=</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">_m</span><span class=\"p\">)</span>\n                    <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"fm\">__setattr__</span><span class=\"p\">(</span><span
  class=\"n\">feed</span><span class=\"o\">.</span><span class=\"n\">name</span><span
  class=\"p\">,</span> <span class=\"n\">feed</span><span class=\"p\">)</span>\n\n
  \           <span class=\"k\">def</span> <span class=\"fm\">__iter__</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">):</span>\n
  \               <span class=\"k\">return</span> <span class=\"nb\">iter</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">keys</span><span
  class=\"p\">())</span>\n\n            <span class=\"k\">def</span> <span class=\"nf\">keys</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">):</span>\n
  \               <span class=\"k\">return</span> <span class=\"nb\">iter</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">keys</span><span
  class=\"p\">())</span>\n\n            <span class=\"k\">def</span> <span class=\"nf\">values</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">):</span>\n
  \               <span class=\"k\">return</span> <span class=\"p\">[</span><span
  class=\"bp\">self</span><span class=\"p\">[</span><span class=\"n\">feed</span><span
  class=\"p\">]</span> <span class=\"k\">for</span> <span class=\"n\">feed</span>
  <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">keys</span><span
  class=\"p\">()]</span>\n\n            <span class=\"k\">def</span> <span class=\"nf\">items</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">):</span>\n
  \               <span class=\"k\">return</span> <span class=\"p\">[(</span><span
  class=\"n\">key</span><span class=\"p\">,</span> <span class=\"bp\">self</span><span
  class=\"p\">[</span><span class=\"n\">key</span><span class=\"p\">])</span> <span
  class=\"k\">for</span> <span class=\"n\">key</span> <span class=\"ow\">in</span>
  <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"p\">]</span>\n\n            <span class=\"k\">def</span> <span class=\"fm\">__getitem__</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
  class=\"n\">key</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span
  class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"nb\">getattr</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
  class=\"n\">key</span><span class=\"o\">.</span><span class=\"n\">replace</span><span
  class=\"p\">(</span><span class=\"s2\">\"-\"</span><span class=\"p\">,</span> <span
  class=\"s2\">\"_\"</span><span class=\"p\">)</span><span class=\"o\">.</span><span
  class=\"n\">lower</span><span class=\"p\">())</span>\n\n            <span class=\"k\">def</span>
  <span class=\"nf\">_dict_panel</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">,</span> <span class=\"n\">config</span><span class=\"p\">)</span> <span
  class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n<span
  class=\"w\">                </span><span class=\"sd\">\"\"\"</span>\n<span class=\"sd\">
  \               pretty print configs with rich</span>\n<span class=\"sd\">                \"\"\"</span>\n
  \               <span class=\"n\">msg</span> <span class=\"o\">=</span> <span class=\"s2\">\"\"</span>\n
  \               <span class=\"k\">for</span> <span class=\"n\">key</span><span class=\"p\">,</span>
  <span class=\"n\">value</span> <span class=\"ow\">in</span> <span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">items</span><span class=\"p\">():</span>\n
  \                   <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
  class=\"p\">(</span><span class=\"n\">value</span><span class=\"p\">,</span> <span
  class=\"nb\">str</span><span class=\"p\">):</span>\n                        <span
  class=\"k\">if</span> <span class=\"nb\">len</span><span class=\"p\">(</span><span
  class=\"n\">value</span><span class=\"p\">)</span> <span class=\"o\">&gt;</span>
  <span class=\"mi\">50</span><span class=\"p\">:</span>\n                            <span
  class=\"n\">value</span> <span class=\"o\">=</span> <span class=\"n\">value</span><span
  class=\"p\">[:</span><span class=\"mi\">50</span><span class=\"p\">]</span> <span
  class=\"o\">+</span> <span class=\"s2\">\"...\"</span>\n                        <span
  class=\"n\">value</span> <span class=\"o\">=</span> <span class=\"n\">value</span>\n
  \                   <span class=\"n\">msg</span> <span class=\"o\">=</span> <span
  class=\"n\">msg</span> <span class=\"o\">+</span> <span class=\"sa\">f</span><span
  class=\"s2\">\"[grey46]</span><span class=\"si\">{</span><span class=\"n\">key</span><span
  class=\"si\">}</span><span class=\"s2\">[/][magenta3]:[/] [grey66]</span><span class=\"si\">{</span><span
  class=\"n\">value</span><span class=\"si\">}</span><span class=\"s2\">[/]</span><span
  class=\"se\">\\n</span><span class=\"s2\">\"</span>\n                <span class=\"k\">return</span>
  <span class=\"n\">msg</span>\n\n            <span class=\"k\">def</span> <span class=\"nf\">__rich__</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">)</span> <span
  class=\"o\">-&gt;</span> <span class=\"n\">Table</span><span class=\"p\">:</span>\n
  \               <span class=\"kn\">from</span> <span class=\"nn\">rich.table</span>
  <span class=\"kn\">import</span> <span class=\"n\">Table</span>\n\n                <span
  class=\"n\">table</span> <span class=\"o\">=</span> <span class=\"n\">Table</span><span
  class=\"p\">(</span><span class=\"n\">title</span><span class=\"o\">=</span><span
  class=\"sa\">f</span><span class=\"s2\">\"Feeds </span><span class=\"si\">{</span><span
  class=\"nb\">len</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">)</span><span
  class=\"si\">}</span><span class=\"s2\">\"</span><span class=\"p\">)</span>\n\n
  \               <span class=\"n\">table</span><span class=\"o\">.</span><span class=\"n\">add_column</span><span
  class=\"p\">(</span><span class=\"s2\">\"Feed\"</span><span class=\"p\">,</span>
  <span class=\"n\">justify</span><span class=\"o\">=</span><span class=\"s2\">\"right\"</span><span
  class=\"p\">,</span> <span class=\"n\">style</span><span class=\"o\">=</span><span
  class=\"s2\">\"cyan\"</span><span class=\"p\">,</span> <span class=\"n\">no_wrap</span><span
  class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n                <span
  class=\"n\">table</span><span class=\"o\">.</span><span class=\"n\">add_column</span><span
  class=\"p\">(</span><span class=\"s2\">\"posts\"</span><span class=\"p\">,</span>
  <span class=\"n\">justify</span><span class=\"o\">=</span><span class=\"s2\">\"left\"</span><span
  class=\"p\">,</span> <span class=\"n\">style</span><span class=\"o\">=</span><span
  class=\"s2\">\"green\"</span><span class=\"p\">)</span>\n                <span class=\"n\">table</span><span
  class=\"o\">.</span><span class=\"n\">add_column</span><span class=\"p\">(</span><span
  class=\"s2\">\"config\"</span><span class=\"p\">,</span> <span class=\"n\">style</span><span
  class=\"o\">=</span><span class=\"s2\">\"magenta\"</span><span class=\"p\">)</span>\n\n
  \               <span class=\"k\">for</span> <span class=\"n\">name</span> <span
  class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">config</span><span class=\"p\">:</span>\n                    <span class=\"n\">table</span><span
  class=\"o\">.</span><span class=\"n\">add_row</span><span class=\"p\">(</span>\n
  \                       <span class=\"n\">name</span><span class=\"p\">,</span>\n
  \                       <span class=\"nb\">str</span><span class=\"p\">(</span><span
  class=\"nb\">len</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">[</span><span class=\"n\">name</span><span class=\"p\">]</span><span
  class=\"o\">.</span><span class=\"n\">posts</span><span class=\"p\">)),</span>\n
  \                       <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">_dict_panel</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">[</span><span
  class=\"n\">name</span><span class=\"p\">]</span><span class=\"o\">.</span><span
  class=\"n\">dict</span><span class=\"p\">()),</span>\n                    <span
  class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"n\">table</span>\n</pre></div>\n\n</pre>\n<p>!!
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
  class=\"p\">:</span> <span class=\"n\">Markata</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config_models</span><span
  class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
  class=\"n\">FeedsConfig</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"pre_render\" style=\"margin:0;padding:.5rem
  1rem;\">pre_render <em class=\"small\">function</em></h2>\nCreate the Feeds object
  and attach it to markata.\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
  class=\"admonition-title\">\"pre_render <em class=\"small\">source</em>\"</p>\n</div>\n<pre
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
  <span class=\"nf\">pre_render</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
  class=\"p\">:</span> <span class=\"n\">Markata</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
  class=\"w\">            </span><span class=\"sd\">\"\"\"</span>\n<span class=\"sd\">
  \           Create the Feeds object and attach it to markata.</span>\n<span class=\"sd\">
  \           \"\"\"</span>\n            <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">feeds</span> <span class=\"o\">=</span> <span class=\"n\">Feeds</span><span
  class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"save\" style=\"margin:0;padding:.5rem
  1rem;\">save <em class=\"small\">function</em></h2>\nCreates a new feed page for
  each page in the config.\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
  class=\"admonition-title\">\"save <em class=\"small\">source</em>\"</p>\n</div>\n<pre
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
  <span class=\"nf\">save</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
  class=\"p\">:</span> <span class=\"n\">Markata</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
  class=\"w\">            </span><span class=\"sd\">\"\"\"</span>\n<span class=\"sd\">
  \           Creates a new feed page for each page in the config.</span>\n<span class=\"sd\">
  \           \"\"\"</span>\n            <span class=\"k\">with</span> <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">cache</span> <span class=\"k\">as</span> <span
  class=\"n\">cache</span><span class=\"p\">:</span>\n                <span class=\"k\">for</span>
  <span class=\"n\">feed</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">feeds</span><span class=\"o\">.</span><span
  class=\"n\">values</span><span class=\"p\">():</span>\n                    <span
  class=\"n\">create_page</span><span class=\"p\">(</span>\n                        <span
  class=\"n\">markata</span><span class=\"p\">,</span>\n                        <span
  class=\"n\">feed</span><span class=\"p\">,</span>\n                        <span
  class=\"n\">cache</span><span class=\"p\">,</span>\n                    <span class=\"p\">)</span>\n\n
  \           <span class=\"n\">home</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
  class=\"p\">(</span><span class=\"nb\">str</span><span class=\"p\">(</span><span
  class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">output_dir</span><span class=\"p\">))</span>
  <span class=\"o\">/</span> <span class=\"s2\">\"index.html\"</span>\n            <span
  class=\"n\">archive</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
  class=\"p\">(</span><span class=\"nb\">str</span><span class=\"p\">(</span><span
  class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">output_dir</span><span class=\"p\">))</span>
  <span class=\"o\">/</span> <span class=\"s2\">\"archive\"</span> <span class=\"o\">/</span>
  <span class=\"s2\">\"index.html\"</span>\n            <span class=\"k\">if</span>
  <span class=\"ow\">not</span> <span class=\"n\">home</span><span class=\"o\">.</span><span
  class=\"n\">exists</span><span class=\"p\">()</span> <span class=\"ow\">and</span>
  <span class=\"n\">archive</span><span class=\"o\">.</span><span class=\"n\">exists</span><span
  class=\"p\">():</span>\n                <span class=\"n\">shutil</span><span class=\"o\">.</span><span
  class=\"n\">copy</span><span class=\"p\">(</span><span class=\"nb\">str</span><span
  class=\"p\">(</span><span class=\"n\">archive</span><span class=\"p\">),</span>
  <span class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">home</span><span
  class=\"p\">))</span>\n\n            <span class=\"n\">xsl_template</span> <span
  class=\"o\">=</span> <span class=\"n\">get_template</span><span class=\"p\">(</span><span
  class=\"n\">feed</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">xsl_template</span><span class=\"p\">)</span>\n
  \           <span class=\"n\">xsl</span> <span class=\"o\">=</span> <span class=\"n\">xsl_template</span><span
  class=\"o\">.</span><span class=\"n\">render</span><span class=\"p\">(</span>\n
  \               <span class=\"n\">markata</span><span class=\"o\">=</span><span
  class=\"n\">markata</span><span class=\"p\">,</span>\n                <span class=\"n\">__version__</span><span
  class=\"o\">=</span><span class=\"n\">__version__</span><span class=\"p\">,</span>\n
  \               <span class=\"n\">today</span><span class=\"o\">=</span><span class=\"n\">datetime</span><span
  class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
  class=\"n\">today</span><span class=\"p\">(),</span>\n                <span class=\"n\">config</span><span
  class=\"o\">=</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">config</span><span class=\"p\">,</span>\n            <span class=\"p\">)</span>\n
  \           <span class=\"n\">xsl_file</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
  class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">output_dir</span><span
  class=\"p\">)</span> <span class=\"o\">/</span> <span class=\"s2\">\"rss.xsl\"</span>\n
  \           <span class=\"n\">xsl_file</span><span class=\"o\">.</span><span class=\"n\">write_text</span><span
  class=\"p\">(</span><span class=\"n\">xsl</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"get_template\" style=\"margin:0;padding:.5rem
  1rem;\">get_template <em class=\"small\">function</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"get_template
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
  <span class=\"nf\">get_template</span><span class=\"p\">(</span><span class=\"n\">src</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Template</span><span
  class=\"p\">:</span>\n            <span class=\"k\">try</span><span class=\"p\">:</span>\n
  \               <span class=\"k\">return</span> <span class=\"n\">Template</span><span
  class=\"p\">(</span><span class=\"n\">Path</span><span class=\"p\">(</span><span
  class=\"n\">src</span><span class=\"p\">)</span><span class=\"o\">.</span><span
  class=\"n\">read_text</span><span class=\"p\">(),</span> <span class=\"n\">undefined</span><span
  class=\"o\">=</span><span class=\"n\">SilentUndefined</span><span class=\"p\">)</span>\n
  \           <span class=\"k\">except</span> <span class=\"ne\">FileNotFoundError</span><span
  class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"n\">Template</span><span
  class=\"p\">(</span><span class=\"n\">src</span><span class=\"p\">,</span> <span
  class=\"n\">undefined</span><span class=\"o\">=</span><span class=\"n\">SilentUndefined</span><span
  class=\"p\">)</span>\n            <span class=\"k\">except</span> <span class=\"ne\">OSError</span><span
  class=\"p\">:</span>  <span class=\"c1\"># File name too long</span>\n                <span
  class=\"k\">return</span> <span class=\"n\">Template</span><span class=\"p\">(</span><span
  class=\"n\">src</span><span class=\"p\">,</span> <span class=\"n\">undefined</span><span
  class=\"o\">=</span><span class=\"n\">SilentUndefined</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"create_page\" style=\"margin:0;padding:.5rem
  1rem;\">create_page <em class=\"small\">function</em></h2>\ncreate an html unorderd
  list of posts.\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
  class=\"admonition-title\">\"create_page <em class=\"small\">source</em>\"</p>\n</div>\n<pre
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
  <span class=\"nf\">create_page</span><span class=\"p\">(</span>\n            <span
  class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"n\">Markata</span><span
  class=\"p\">,</span>\n            <span class=\"n\">feed</span><span class=\"p\">:</span>
  <span class=\"n\">Feed</span><span class=\"p\">,</span>\n            <span class=\"n\">cache</span><span
  class=\"p\">,</span>\n        <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span class=\"w\">            </span><span
  class=\"sd\">\"\"\"</span>\n<span class=\"sd\">            create an html unorderd
  list of posts.</span>\n<span class=\"sd\">            \"\"\"</span>\n\n            <span
  class=\"n\">posts</span> <span class=\"o\">=</span> <span class=\"n\">feed</span><span
  class=\"o\">.</span><span class=\"n\">posts</span>\n\n            <span class=\"n\">cards</span>
  <span class=\"o\">=</span> <span class=\"p\">[</span>\n                <span class=\"n\">create_card</span><span
  class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">,</span> <span
  class=\"n\">post</span><span class=\"p\">,</span> <span class=\"n\">feed</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
  class=\"n\">card_template</span><span class=\"p\">,</span> <span class=\"n\">cache</span><span
  class=\"p\">)</span> <span class=\"k\">for</span> <span class=\"n\">post</span>
  <span class=\"ow\">in</span> <span class=\"n\">posts</span>\n            <span class=\"p\">]</span>\n
  \           <span class=\"n\">cards</span><span class=\"o\">.</span><span class=\"n\">insert</span><span
  class=\"p\">(</span><span class=\"mi\">0</span><span class=\"p\">,</span> <span
  class=\"s2\">\"&lt;ul&gt;\"</span><span class=\"p\">)</span>\n            <span
  class=\"n\">cards</span><span class=\"o\">.</span><span class=\"n\">append</span><span
  class=\"p\">(</span><span class=\"s2\">\"&lt;/ul&gt;\"</span><span class=\"p\">)</span>\n
  \           <span class=\"n\">cards</span> <span class=\"o\">=</span> <span class=\"s2\">\"\"</span><span
  class=\"o\">.</span><span class=\"n\">join</span><span class=\"p\">(</span><span
  class=\"n\">cards</span><span class=\"p\">)</span>\n\n            <span class=\"n\">template</span>
  <span class=\"o\">=</span> <span class=\"n\">get_template</span><span class=\"p\">(</span><span
  class=\"n\">feed</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">template</span><span class=\"p\">)</span>\n
  \           <span class=\"n\">rss_template</span> <span class=\"o\">=</span> <span
  class=\"n\">get_template</span><span class=\"p\">(</span><span class=\"n\">feed</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
  class=\"n\">rss_template</span><span class=\"p\">)</span>\n            <span class=\"n\">sitemap_template</span>
  <span class=\"o\">=</span> <span class=\"n\">get_template</span><span class=\"p\">(</span><span
  class=\"n\">feed</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">sitemap_template</span><span class=\"p\">)</span>\n
  \           <span class=\"n\">output_file</span> <span class=\"o\">=</span> <span
  class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
  class=\"n\">output_dir</span><span class=\"p\">)</span> <span class=\"o\">/</span>
  <span class=\"n\">feed</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">slug</span> <span class=\"o\">/</span> <span
  class=\"s2\">\"index.html\"</span>\n            <span class=\"n\">canonical_url</span>
  <span class=\"o\">=</span> <span class=\"sa\">f</span><span class=\"s2\">\"</span><span
  class=\"si\">{</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">url</span><span
  class=\"si\">}</span><span class=\"s2\">/</span><span class=\"si\">{</span><span
  class=\"n\">feed</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">slug</span><span class=\"si\">}</span><span
  class=\"s2\">/\"</span>\n            <span class=\"n\">output_file</span><span class=\"o\">.</span><span
  class=\"n\">parent</span><span class=\"o\">.</span><span class=\"n\">mkdir</span><span
  class=\"p\">(</span><span class=\"n\">exist_ok</span><span class=\"o\">=</span><span
  class=\"kc\">True</span><span class=\"p\">,</span> <span class=\"n\">parents</span><span
  class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n\n
  \           <span class=\"n\">rss_output_file</span> <span class=\"o\">=</span>
  <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
  class=\"n\">output_dir</span><span class=\"p\">)</span> <span class=\"o\">/</span>
  <span class=\"n\">feed</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">slug</span> <span class=\"o\">/</span> <span
  class=\"s2\">\"rss.xml\"</span>\n            <span class=\"n\">rss_output_file</span><span
  class=\"o\">.</span><span class=\"n\">parent</span><span class=\"o\">.</span><span
  class=\"n\">mkdir</span><span class=\"p\">(</span><span class=\"n\">exist_ok</span><span
  class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
  class=\"n\">parents</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
  class=\"p\">)</span>\n\n            <span class=\"n\">sitemap_output_file</span>
  <span class=\"o\">=</span> <span class=\"p\">(</span>\n                <span class=\"n\">Path</span><span
  class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">output_dir</span><span
  class=\"p\">)</span> <span class=\"o\">/</span> <span class=\"n\">feed</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
  class=\"n\">slug</span> <span class=\"o\">/</span> <span class=\"s2\">\"sitemap.xml\"</span>\n
  \           <span class=\"p\">)</span>\n            <span class=\"n\">sitemap_output_file</span><span
  class=\"o\">.</span><span class=\"n\">parent</span><span class=\"o\">.</span><span
  class=\"n\">mkdir</span><span class=\"p\">(</span><span class=\"n\">exist_ok</span><span
  class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
  class=\"n\">parents</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
  class=\"p\">)</span>\n\n            <span class=\"n\">key</span> <span class=\"o\">=</span>
  <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">make_hash</span><span
  class=\"p\">(</span>\n                <span class=\"s2\">\"feeds\"</span><span class=\"p\">,</span>\n
  \               <span class=\"n\">template</span><span class=\"p\">,</span>\n                <span
  class=\"n\">__version__</span><span class=\"p\">,</span>\n                <span
  class=\"n\">cards</span><span class=\"p\">,</span>\n                <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
  class=\"n\">url</span><span class=\"p\">,</span>\n                <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
  class=\"n\">description</span><span class=\"p\">,</span>\n                <span
  class=\"n\">feed</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">title</span><span class=\"p\">,</span>\n                <span
  class=\"n\">canonical_url</span><span class=\"p\">,</span>\n                <span
  class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
  class=\"o\">.</span><span class=\"n\">today</span><span class=\"p\">(),</span>\n
  \               <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">config</span><span class=\"p\">,</span>\n            <span class=\"p\">)</span>\n\n
  \           <span class=\"n\">feed_html_from_cache</span> <span class=\"o\">=</span>
  <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">precache</span><span
  class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
  class=\"n\">key</span><span class=\"p\">)</span>\n            <span class=\"k\">if</span>
  <span class=\"n\">feed_html_from_cache</span> <span class=\"ow\">is</span> <span
  class=\"kc\">None</span><span class=\"p\">:</span>\n                <span class=\"n\">feed_html</span>
  <span class=\"o\">=</span> <span class=\"n\">template</span><span class=\"o\">.</span><span
  class=\"n\">render</span><span class=\"p\">(</span>\n                    <span class=\"n\">markata</span><span
  class=\"o\">=</span><span class=\"n\">markata</span><span class=\"p\">,</span>\n
  \                   <span class=\"n\">__version__</span><span class=\"o\">=</span><span
  class=\"n\">__version__</span><span class=\"p\">,</span>\n                    <span
  class=\"n\">body</span><span class=\"o\">=</span><span class=\"n\">cards</span><span
  class=\"p\">,</span>\n                    <span class=\"n\">url</span><span class=\"o\">=</span><span
  class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">url</span><span class=\"p\">,</span>\n                    <span
  class=\"n\">description</span><span class=\"o\">=</span><span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
  class=\"n\">description</span><span class=\"p\">,</span>\n                    <span
  class=\"n\">title</span><span class=\"o\">=</span><span class=\"n\">feed</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
  class=\"n\">title</span><span class=\"p\">,</span>\n                    <span class=\"n\">canonical_url</span><span
  class=\"o\">=</span><span class=\"n\">canonical_url</span><span class=\"p\">,</span>\n
  \                   <span class=\"n\">today</span><span class=\"o\">=</span><span
  class=\"n\">datetime</span><span class=\"o\">.</span><span class=\"n\">datetime</span><span
  class=\"o\">.</span><span class=\"n\">today</span><span class=\"p\">(),</span>\n
  \                   <span class=\"n\">config</span><span class=\"o\">=</span><span
  class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"p\">,</span>\n                <span class=\"p\">)</span>\n                <span
  class=\"k\">with</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">cache</span> <span class=\"k\">as</span> <span class=\"n\">cache</span><span
  class=\"p\">:</span>\n                    <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">cache</span><span class=\"o\">.</span><span
  class=\"n\">set</span><span class=\"p\">(</span><span class=\"n\">key</span><span
  class=\"p\">,</span> <span class=\"n\">feed_html</span><span class=\"p\">)</span>\n\n
  \           <span class=\"n\">feed_rss</span> <span class=\"o\">=</span> <span class=\"n\">rss_template</span><span
  class=\"o\">.</span><span class=\"n\">render</span><span class=\"p\">(</span><span
  class=\"n\">markata</span><span class=\"o\">=</span><span class=\"n\">markata</span><span
  class=\"p\">,</span> <span class=\"n\">feed</span><span class=\"o\">=</span><span
  class=\"n\">feed</span><span class=\"p\">)</span>\n            <span class=\"n\">feed_sitemap</span>
  <span class=\"o\">=</span> <span class=\"n\">sitemap_template</span><span class=\"o\">.</span><span
  class=\"n\">render</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
  class=\"o\">=</span><span class=\"n\">markata</span><span class=\"p\">,</span> <span
  class=\"n\">feed</span><span class=\"o\">=</span><span class=\"n\">feed</span><span
  class=\"p\">)</span>\n\n            <span class=\"n\">output_file</span><span class=\"o\">.</span><span
  class=\"n\">write_text</span><span class=\"p\">(</span><span class=\"n\">feed_html</span><span
  class=\"p\">)</span>\n            <span class=\"n\">rss_output_file</span><span
  class=\"o\">.</span><span class=\"n\">write_text</span><span class=\"p\">(</span><span
  class=\"n\">feed_rss</span><span class=\"p\">)</span>\n            <span class=\"n\">sitemap_output_file</span><span
  class=\"o\">.</span><span class=\"n\">write_text</span><span class=\"p\">(</span><span
  class=\"n\">feed_sitemap</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"create_card\" style=\"margin:0;padding:.5rem
  1rem;\">create_card <em class=\"small\">function</em></h2>\nCreates a card for one
  post based on the configured template.  If no\ntemplate is configured it will create
  one with the post title and dates\n(if present).\n<div class=\"admonition source
  is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"create_card
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
  <span class=\"nf\">create_card</span><span class=\"p\">(</span>\n            <span
  class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span
  class=\"p\">,</span>\n            <span class=\"n\">post</span><span class=\"p\">:</span>
  <span class=\"s2\">\"Post\"</span><span class=\"p\">,</span>\n            <span
  class=\"n\">template</span><span class=\"p\">:</span> <span class=\"n\">Optional</span><span
  class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">]</span> <span
  class=\"o\">=</span> <span class=\"kc\">None</span><span class=\"p\">,</span>\n
  \           <span class=\"n\">cache</span><span class=\"o\">=</span><span class=\"kc\">None</span><span
  class=\"p\">,</span>\n        <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"n\">Any</span><span class=\"p\">:</span>\n<span class=\"w\">            </span><span
  class=\"sd\">\"\"\"</span>\n<span class=\"sd\">            Creates a card for one
  post based on the configured template.  If no</span>\n<span class=\"sd\">            template
  is configured it will create one with the post title and dates</span>\n<span class=\"sd\">
  \           (if present).</span>\n<span class=\"sd\">            \"\"\"</span>\n
  \           <span class=\"n\">key</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">make_hash</span><span class=\"p\">(</span><span
  class=\"s2\">\"feeds\"</span><span class=\"p\">,</span> <span class=\"n\">template</span><span
  class=\"p\">,</span> <span class=\"nb\">str</span><span class=\"p\">(</span><span
  class=\"n\">post</span><span class=\"p\">),</span> <span class=\"n\">post</span><span
  class=\"o\">.</span><span class=\"n\">content</span><span class=\"p\">)</span>\n\n
  \           <span class=\"n\">card</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">precache</span><span class=\"o\">.</span><span
  class=\"n\">get</span><span class=\"p\">(</span><span class=\"n\">key</span><span
  class=\"p\">)</span>\n            <span class=\"k\">if</span> <span class=\"n\">card</span>
  <span class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"n\">card</span>\n\n
  \           <span class=\"k\">if</span> <span class=\"n\">template</span> <span
  class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \               <span class=\"n\">template</span> <span class=\"o\">=</span> <span
  class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
  class=\"s2\">\"feeds_config\"</span><span class=\"p\">,</span> <span class=\"p\">{})</span><span
  class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
  class=\"s2\">\"card_template\"</span><span class=\"p\">,</span> <span class=\"kc\">None</span><span
  class=\"p\">)</span>\n\n            <span class=\"k\">if</span> <span class=\"n\">template</span>
  <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \               <span class=\"k\">if</span> <span class=\"s2\">\"date\"</span> <span
  class=\"ow\">in</span> <span class=\"n\">post</span><span class=\"p\">:</span>\n
  \                   <span class=\"n\">card</span> <span class=\"o\">=</span> <span
  class=\"n\">textwrap</span><span class=\"o\">.</span><span class=\"n\">dedent</span><span
  class=\"p\">(</span>\n                        <span class=\"sa\">f</span><span class=\"s2\">\"\"\"</span>\n<span
  class=\"s2\">                        &lt;li class='post'&gt;</span>\n<span class=\"s2\">
  \                       &lt;a href=\"/</span><span class=\"si\">{</span><span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
  class=\"n\">path_prefix</span><span class=\"si\">}{</span><span class=\"n\">post</span><span
  class=\"o\">.</span><span class=\"n\">slug</span><span class=\"si\">}</span><span
  class=\"s2\">/\"&gt;</span>\n<span class=\"s2\">                            </span><span
  class=\"si\">{</span><span class=\"n\">post</span><span class=\"o\">.</span><span
  class=\"n\">title</span><span class=\"si\">}</span>\n<span class=\"s2\">                            </span><span
  class=\"si\">{</span><span class=\"n\">post</span><span class=\"o\">.</span><span
  class=\"n\">date</span><span class=\"o\">.</span><span class=\"n\">year</span><span
  class=\"si\">}</span><span class=\"s2\">-</span>\n<span class=\"s2\">                            </span><span
  class=\"si\">{</span><span class=\"n\">post</span><span class=\"o\">.</span><span
  class=\"n\">date</span><span class=\"o\">.</span><span class=\"n\">month</span><span
  class=\"si\">}</span><span class=\"s2\">-</span>\n<span class=\"s2\">                            </span><span
  class=\"si\">{</span><span class=\"n\">post</span><span class=\"o\">.</span><span
  class=\"n\">date</span><span class=\"o\">.</span><span class=\"n\">day</span><span
  class=\"si\">}</span>\n<span class=\"s2\">                        &lt;/a&gt;</span>\n<span
  class=\"s2\">                        &lt;/li&gt;</span>\n<span class=\"s2\">                        \"\"\"</span><span
  class=\"p\">,</span>\n                    <span class=\"p\">)</span>\n                <span
  class=\"k\">else</span><span class=\"p\">:</span>\n                    <span class=\"n\">card</span>
  <span class=\"o\">=</span> <span class=\"n\">textwrap</span><span class=\"o\">.</span><span
  class=\"n\">dedent</span><span class=\"p\">(</span>\n                        <span
  class=\"sa\">f</span><span class=\"s2\">\"\"\"</span>\n<span class=\"s2\">                        &lt;li
  class='post'&gt;</span>\n<span class=\"s2\">                        &lt;a href=\"/</span><span
  class=\"si\">{</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">path_prefix</span><span
  class=\"si\">}{</span><span class=\"n\">post</span><span class=\"o\">.</span><span
  class=\"n\">slug</span><span class=\"si\">}</span><span class=\"s2\">/\"&gt;</span>\n<span
  class=\"s2\">                            </span><span class=\"si\">{</span><span
  class=\"n\">post</span><span class=\"o\">.</span><span class=\"n\">title</span><span
  class=\"si\">}</span>\n<span class=\"s2\">                        &lt;/a&gt;</span>\n<span
  class=\"s2\">                        &lt;/li&gt;</span>\n<span class=\"s2\">                        \"\"\"</span><span
  class=\"p\">,</span>\n                    <span class=\"p\">)</span>\n            <span
  class=\"k\">else</span><span class=\"p\">:</span>\n                <span class=\"k\">try</span><span
  class=\"p\">:</span>\n                    <span class=\"n\">_template</span> <span
  class=\"o\">=</span> <span class=\"n\">Template</span><span class=\"p\">(</span><span
  class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">template</span><span
  class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">read_text</span><span
  class=\"p\">())</span>\n                <span class=\"k\">except</span> <span class=\"ne\">FileNotFoundError</span><span
  class=\"p\">:</span>\n                    <span class=\"n\">_template</span> <span
  class=\"o\">=</span> <span class=\"n\">Template</span><span class=\"p\">(</span><span
  class=\"n\">template</span><span class=\"p\">)</span>\n                <span class=\"k\">except</span>
  <span class=\"ne\">OSError</span><span class=\"p\">:</span>  <span class=\"c1\">#
  File name too long</span>\n                    <span class=\"n\">_template</span>
  <span class=\"o\">=</span> <span class=\"n\">Template</span><span class=\"p\">(</span><span
  class=\"n\">template</span><span class=\"p\">)</span>\n                <span class=\"n\">card</span>
  <span class=\"o\">=</span> <span class=\"n\">_template</span><span class=\"o\">.</span><span
  class=\"n\">render</span><span class=\"p\">(</span><span class=\"n\">post</span><span
  class=\"o\">=</span><span class=\"n\">post</span><span class=\"p\">,</span> <span
  class=\"o\">**</span><span class=\"n\">post</span><span class=\"o\">.</span><span
  class=\"n\">to_dict</span><span class=\"p\">())</span>\n            <span class=\"n\">cache</span><span
  class=\"o\">.</span><span class=\"n\">add</span><span class=\"p\">(</span><span
  class=\"n\">key</span><span class=\"p\">,</span> <span class=\"n\">card</span><span
  class=\"p\">)</span>\n            <span class=\"k\">return</span> <span class=\"n\">card</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"cli\" style=\"margin:0;padding:.5rem
  1rem;\">cli <em class=\"small\">function</em></h2>\n<div class=\"admonition source
  is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"cli <em class=\"small\">source</em>\"</p>\n</div>\n<pre
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
  <span class=\"nf\">cli</span><span class=\"p\">(</span><span class=\"n\">app</span><span
  class=\"p\">:</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
  class=\"n\">Typer</span><span class=\"p\">,</span> <span class=\"n\">markata</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \           <span class=\"n\">feeds_app</span> <span class=\"o\">=</span> <span
  class=\"n\">typer</span><span class=\"o\">.</span><span class=\"n\">Typer</span><span
  class=\"p\">()</span>\n            <span class=\"n\">app</span><span class=\"o\">.</span><span
  class=\"n\">add_typer</span><span class=\"p\">(</span><span class=\"n\">feeds_app</span><span
  class=\"p\">)</span>\n\n            <span class=\"nd\">@feeds_app</span><span class=\"o\">.</span><span
  class=\"n\">callback</span><span class=\"p\">()</span>\n            <span class=\"k\">def</span>
  <span class=\"nf\">feeds</span><span class=\"p\">():</span>\n                <span
  class=\"s2\">\"feeds cli\"</span>\n\n            <span class=\"nd\">@feeds_app</span><span
  class=\"o\">.</span><span class=\"n\">command</span><span class=\"p\">()</span>\n
  \           <span class=\"k\">def</span> <span class=\"nf\">show</span><span class=\"p\">()</span>
  <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \               <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">quiet</span>
  <span class=\"o\">=</span> <span class=\"kc\">True</span>\n                <span
  class=\"n\">feeds</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">feeds</span>\n                <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
  class=\"n\">quiet</span> <span class=\"o\">=</span> <span class=\"kc\">False</span>\n
  \               <span class=\"n\">rich_print</span><span class=\"p\">(</span><span
  class=\"n\">feeds</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"_fail_with_undefined_error\" style=\"margin:0;padding:.5rem
  1rem;\">_fail_with_undefined_error <em class=\"small\">method</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"_fail_with_undefined_error
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
  <span class=\"nf\">_fail_with_undefined_error</span><span class=\"p\">(</span><span
  class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
  class=\"n\">args</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
  class=\"n\">kwargs</span><span class=\"p\">):</span>\n                <span class=\"k\">return</span>
  <span class=\"s2\">\"\"</span>\n</pre></div>\n\n</pre>\n<p>!! method </p><h2 class=\"admonition-title\"
  id=\"default_name\" style=\"margin:0;padding:.5rem 1rem;\">default_name <em class=\"small\">method</em></h2>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"default_name
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
  <span class=\"nf\">default_name</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
  class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
  class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
  class=\"p\">):</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span>
  <span class=\"ow\">or</span> <span class=\"nb\">str</span><span class=\"p\">(</span><span
  class=\"n\">values</span><span class=\"o\">.</span><span class=\"n\">get</span><span
  class=\"p\">(</span><span class=\"s2\">\"slug\"</span><span class=\"p\">))</span><span
  class=\"o\">.</span><span class=\"n\">replace</span><span class=\"p\">(</span><span
  class=\"s2\">\"-\"</span><span class=\"p\">,</span> <span class=\"s2\">\"_\"</span><span
  class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!! method </p><h2 class=\"admonition-title\"
  id=\"read_template\" style=\"margin:0;padding:.5rem 1rem;\">read_template <em class=\"small\">method</em></h2>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"read_template
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
  <span class=\"nf\">read_template</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
  class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
  class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span
  class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
  class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span class=\"n\">Path</span><span
  class=\"p\">):</span>\n                    <span class=\"k\">return</span> <span
  class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">v</span><span
  class=\"o\">.</span><span class=\"n\">read_text</span><span class=\"p\">())</span>\n
  \               <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"name\" style=\"margin:0;padding:.5rem
  1rem;\">name <em class=\"small\">method</em></h2>\n<div class=\"admonition source
  is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"name <em class=\"small\">source</em>\"</p>\n</div>\n<pre
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
  <span class=\"nf\">name</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">):</span>\n                <span class=\"k\">return</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
  class=\"n\">name</span>\n</pre></div>\n\n</pre>\n<p>!! method </p><h2 class=\"admonition-title\"
  id=\"posts\" style=\"margin:0;padding:.5rem 1rem;\">posts <em class=\"small\">method</em></h2>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"posts
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
  <span class=\"nf\">posts</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">):</span>\n                <span class=\"k\">return</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">map</span><span class=\"p\">(</span><span
  class=\"s2\">\"post\"</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"map\" style=\"margin:0;padding:.5rem
  1rem;\">map <em class=\"small\">method</em></h2>\n<div class=\"admonition source
  is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"map <em class=\"small\">source</em>\"</p>\n</div>\n<pre
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
  <span class=\"nf\">map</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">,</span> <span class=\"n\">func</span><span class=\"o\">=</span><span
  class=\"s2\">\"post\"</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
  class=\"n\">args</span><span class=\"p\">):</span>\n                <span class=\"k\">return</span>
  <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_m</span><span
  class=\"o\">.</span><span class=\"n\">map</span><span class=\"p\">(</span><span
  class=\"n\">func</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
  class=\"p\">{</span><span class=\"o\">**</span><span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
  class=\"n\">dict</span><span class=\"p\">(),</span> <span class=\"o\">**</span><span
  class=\"n\">args</span><span class=\"p\">})</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"__init__\" style=\"margin:0;padding:.5rem
  1rem;\"><strong>init</strong> <em class=\"small\">method</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"<strong>init</strong>
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
  class=\"p\">,</span> <span class=\"n\">markata</span><span class=\"p\">:</span>
  <span class=\"n\">Markata</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_m</span> <span
  class=\"o\">=</span> <span class=\"n\">markata</span>\n                <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">config</span> <span class=\"o\">=</span> <span
  class=\"p\">{</span><span class=\"n\">f</span><span class=\"o\">.</span><span class=\"n\">name</span><span
  class=\"p\">:</span> <span class=\"n\">f</span> <span class=\"k\">for</span> <span
  class=\"n\">f</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
  class=\"n\">feeds</span><span class=\"p\">}</span>\n                <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">refresh</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"refresh\" style=\"margin:0;padding:.5rem
  1rem;\">refresh <em class=\"small\">method</em></h2>\nRefresh all of the feeds objects\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"refresh
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
  <span class=\"nf\">refresh</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n<span class=\"w\">                </span><span class=\"sd\">\"\"\"</span>\n<span
  class=\"sd\">                Refresh all of the feeds objects</span>\n<span class=\"sd\">
  \               \"\"\"</span>\n                <span class=\"k\">for</span> <span
  class=\"n\">feed</span> <span class=\"ow\">in</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">_m</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">feeds</span><span class=\"p\">:</span>\n                    <span
  class=\"n\">feed</span> <span class=\"o\">=</span> <span class=\"n\">Feed</span><span
  class=\"p\">(</span><span class=\"n\">config</span><span class=\"o\">=</span><span
  class=\"n\">feed</span><span class=\"p\">,</span> <span class=\"n\">_m</span><span
  class=\"o\">=</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">_m</span><span class=\"p\">)</span>\n                    <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"fm\">__setattr__</span><span class=\"p\">(</span><span
  class=\"n\">feed</span><span class=\"o\">.</span><span class=\"n\">name</span><span
  class=\"p\">,</span> <span class=\"n\">feed</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
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
  <span class=\"fm\">__iter__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">):</span>\n                <span class=\"k\">return</span> <span class=\"nb\">iter</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">keys</span><span
  class=\"p\">())</span>\n</pre></div>\n\n</pre>\n<p>!! method </p><h2 class=\"admonition-title\"
  id=\"keys\" style=\"margin:0;padding:.5rem 1rem;\">keys <em class=\"small\">method</em></h2>\n<div
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
  class=\"p\">):</span>\n                <span class=\"k\">return</span> <span class=\"nb\">iter</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">keys</span><span
  class=\"p\">())</span>\n</pre></div>\n\n</pre>\n<p>!! method </p><h2 class=\"admonition-title\"
  id=\"values\" style=\"margin:0;padding:.5rem 1rem;\">values <em class=\"small\">method</em></h2>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"values
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
  <span class=\"nf\">values</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">):</span>\n                <span class=\"k\">return</span> <span class=\"p\">[</span><span
  class=\"bp\">self</span><span class=\"p\">[</span><span class=\"n\">feed</span><span
  class=\"p\">]</span> <span class=\"k\">for</span> <span class=\"n\">feed</span>
  <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">keys</span><span
  class=\"p\">()]</span>\n</pre></div>\n\n</pre>\n<p>!! method </p><h2 class=\"admonition-title\"
  id=\"items\" style=\"margin:0;padding:.5rem 1rem;\">items <em class=\"small\">method</em></h2>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"items
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
  <span class=\"nf\">items</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">):</span>\n                <span class=\"k\">return</span> <span class=\"p\">[(</span><span
  class=\"n\">key</span><span class=\"p\">,</span> <span class=\"bp\">self</span><span
  class=\"p\">[</span><span class=\"n\">key</span><span class=\"p\">])</span> <span
  class=\"k\">for</span> <span class=\"n\">key</span> <span class=\"ow\">in</span>
  <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"p\">]</span>\n</pre></div>\n\n</pre>\n<p>!! method </p><h2 class=\"admonition-title\"
  id=\"__getitem__\" style=\"margin:0;padding:.5rem 1rem;\"><strong>getitem</strong>
  <em class=\"small\">method</em></h2>\n<div class=\"admonition source is-collapsible
  collapsible-open\">\n<p class=\"admonition-title\">\"<strong>getitem</strong> <em
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
  <span class=\"fm\">__getitem__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">,</span> <span class=\"n\">key</span><span class=\"p\">:</span> <span
  class=\"nb\">str</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"n\">Any</span><span class=\"p\">:</span>\n                <span class=\"k\">return</span>
  <span class=\"nb\">getattr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">,</span> <span class=\"n\">key</span><span class=\"o\">.</span><span
  class=\"n\">replace</span><span class=\"p\">(</span><span class=\"s2\">\"-\"</span><span
  class=\"p\">,</span> <span class=\"s2\">\"_\"</span><span class=\"p\">)</span><span
  class=\"o\">.</span><span class=\"n\">lower</span><span class=\"p\">())</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"_dict_panel\" style=\"margin:0;padding:.5rem
  1rem;\">_dict_panel <em class=\"small\">method</em></h2>\npretty print configs with
  rich\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"_dict_panel
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
  <span class=\"nf\">_dict_panel</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">,</span> <span class=\"n\">config</span><span class=\"p\">)</span> <span
  class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n<span
  class=\"w\">                </span><span class=\"sd\">\"\"\"</span>\n<span class=\"sd\">
  \               pretty print configs with rich</span>\n<span class=\"sd\">                \"\"\"</span>\n
  \               <span class=\"n\">msg</span> <span class=\"o\">=</span> <span class=\"s2\">\"\"</span>\n
  \               <span class=\"k\">for</span> <span class=\"n\">key</span><span class=\"p\">,</span>
  <span class=\"n\">value</span> <span class=\"ow\">in</span> <span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">items</span><span class=\"p\">():</span>\n
  \                   <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
  class=\"p\">(</span><span class=\"n\">value</span><span class=\"p\">,</span> <span
  class=\"nb\">str</span><span class=\"p\">):</span>\n                        <span
  class=\"k\">if</span> <span class=\"nb\">len</span><span class=\"p\">(</span><span
  class=\"n\">value</span><span class=\"p\">)</span> <span class=\"o\">&gt;</span>
  <span class=\"mi\">50</span><span class=\"p\">:</span>\n                            <span
  class=\"n\">value</span> <span class=\"o\">=</span> <span class=\"n\">value</span><span
  class=\"p\">[:</span><span class=\"mi\">50</span><span class=\"p\">]</span> <span
  class=\"o\">+</span> <span class=\"s2\">\"...\"</span>\n                        <span
  class=\"n\">value</span> <span class=\"o\">=</span> <span class=\"n\">value</span>\n
  \                   <span class=\"n\">msg</span> <span class=\"o\">=</span> <span
  class=\"n\">msg</span> <span class=\"o\">+</span> <span class=\"sa\">f</span><span
  class=\"s2\">\"[grey46]</span><span class=\"si\">{</span><span class=\"n\">key</span><span
  class=\"si\">}</span><span class=\"s2\">[/][magenta3]:[/] [grey66]</span><span class=\"si\">{</span><span
  class=\"n\">value</span><span class=\"si\">}</span><span class=\"s2\">[/]</span><span
  class=\"se\">\\n</span><span class=\"s2\">\"</span>\n                <span class=\"k\">return</span>
  <span class=\"n\">msg</span>\n</pre></div>\n\n</pre>\n<p>!! method </p><h2 class=\"admonition-title\"
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
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Table</span><span
  class=\"p\">:</span>\n                <span class=\"kn\">from</span> <span class=\"nn\">rich.table</span>
  <span class=\"kn\">import</span> <span class=\"n\">Table</span>\n\n                <span
  class=\"n\">table</span> <span class=\"o\">=</span> <span class=\"n\">Table</span><span
  class=\"p\">(</span><span class=\"n\">title</span><span class=\"o\">=</span><span
  class=\"sa\">f</span><span class=\"s2\">\"Feeds </span><span class=\"si\">{</span><span
  class=\"nb\">len</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">)</span><span
  class=\"si\">}</span><span class=\"s2\">\"</span><span class=\"p\">)</span>\n\n
  \               <span class=\"n\">table</span><span class=\"o\">.</span><span class=\"n\">add_column</span><span
  class=\"p\">(</span><span class=\"s2\">\"Feed\"</span><span class=\"p\">,</span>
  <span class=\"n\">justify</span><span class=\"o\">=</span><span class=\"s2\">\"right\"</span><span
  class=\"p\">,</span> <span class=\"n\">style</span><span class=\"o\">=</span><span
  class=\"s2\">\"cyan\"</span><span class=\"p\">,</span> <span class=\"n\">no_wrap</span><span
  class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n                <span
  class=\"n\">table</span><span class=\"o\">.</span><span class=\"n\">add_column</span><span
  class=\"p\">(</span><span class=\"s2\">\"posts\"</span><span class=\"p\">,</span>
  <span class=\"n\">justify</span><span class=\"o\">=</span><span class=\"s2\">\"left\"</span><span
  class=\"p\">,</span> <span class=\"n\">style</span><span class=\"o\">=</span><span
  class=\"s2\">\"green\"</span><span class=\"p\">)</span>\n                <span class=\"n\">table</span><span
  class=\"o\">.</span><span class=\"n\">add_column</span><span class=\"p\">(</span><span
  class=\"s2\">\"config\"</span><span class=\"p\">,</span> <span class=\"n\">style</span><span
  class=\"o\">=</span><span class=\"s2\">\"magenta\"</span><span class=\"p\">)</span>\n\n
  \               <span class=\"k\">for</span> <span class=\"n\">name</span> <span
  class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">config</span><span class=\"p\">:</span>\n                    <span class=\"n\">table</span><span
  class=\"o\">.</span><span class=\"n\">add_row</span><span class=\"p\">(</span>\n
  \                       <span class=\"n\">name</span><span class=\"p\">,</span>\n
  \                       <span class=\"nb\">str</span><span class=\"p\">(</span><span
  class=\"nb\">len</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">[</span><span class=\"n\">name</span><span class=\"p\">]</span><span
  class=\"o\">.</span><span class=\"n\">posts</span><span class=\"p\">)),</span>\n
  \                       <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">_dict_panel</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">[</span><span
  class=\"n\">name</span><span class=\"p\">]</span><span class=\"o\">.</span><span
  class=\"n\">dict</span><span class=\"p\">()),</span>\n                    <span
  class=\"p\">)</span>\n                <span class=\"k\">return</span> <span class=\"n\">table</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"feeds\" style=\"margin:0;padding:.5rem
  1rem;\">feeds <em class=\"small\">function</em></h2>\nfeeds cli\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"feeds
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
  <span class=\"nf\">feeds</span><span class=\"p\">():</span>\n                <span
  class=\"s2\">\"feeds cli\"</span>\n</pre></div>\n\n</pre>\n<p>!! function </p><h2
  class=\"admonition-title\" id=\"show\" style=\"margin:0;padding:.5rem 1rem;\">show
  <em class=\"small\">function</em></h2>\n<div class=\"admonition source is-collapsible
  collapsible-open\">\n<p class=\"admonition-title\">\"show <em class=\"small\">source</em>\"</p>\n</div>\n<pre
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
  <span class=\"nf\">show</span><span class=\"p\">()</span> <span class=\"o\">-&gt;</span>
  <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
  class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
  class=\"o\">.</span><span class=\"n\">quiet</span> <span class=\"o\">=</span> <span
  class=\"kc\">True</span>\n                <span class=\"n\">feeds</span> <span class=\"o\">=</span>
  <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">feeds</span>\n
  \               <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">quiet</span>
  <span class=\"o\">=</span> <span class=\"kc\">False</span>\n                <span
  class=\"n\">rich_print</span><span class=\"p\">(</span><span class=\"n\">feeds</span><span
  class=\"p\">)</span>\n</pre></div>\n\n</pre>\n</main>\n<footer>\xA9 2024</footer>\n</body></html>"
published: true
slug: markata/plugins/feeds
title: Feeds.Py


---

The `markata.plugins.feeds` plugin is used to create feed pages, which are lists of
posts.  The list is generated using a `filter`, then each post in the list is
rendered with a `card_template` before being applied to the `body` of the
`template`.

# Installation

This plugin is built-in and enabled by default, but in you want to be very
explicit you can add it to your list of existing plugins.

``` toml
hooks = [
   "markata.plugins.feeds",
   ]
```

# Configuration

# set default template and card_template

At the root of the markata.feeds config you may set `template`, and
`card_template`.  These will become your defaults for every feed you create.
If you do not set these, markata will use it's defaults.  The defaults are
designed to work for a variety of use cases, but are not likely the best for
all.

``` toml
[markata.feeds_config]
template="pages/templates/archive_template.html"
card_template="plugins/feed_card_template.html"
```

# pages

Underneath of the `markata.feeds` we will create a new map for each page where
the name of the map will be the name of the page.


The following config will create a page at `/all-posts` that inclues every
single post.

``` toml
[[markata.feeds]]
title="All Posts"
slug='all'
filter="True"
```

# template

The `template` configuration key is a file path to the template that you want
to use to create the feed.  You may set the default template you want to use
for all feeds under `[markata.feeds]`, as well as override it inside of each
feeds config.

The template is a jinja style template that expects to fill in a `title` and
`body` variable.

``` html
<!DOCTYPE html>
<html lang="en">
  <head>
    <title>{{ title }}</title>
  </head>
  <body>
    <ul>
        {{ body }}
    </ul>
  </body>
</html>
```

!!! note
    I highly reccomend putting your `body` in a `<ul>`, and wrapping your
    `card_template`s in an `<li>`.

# card_template

All keys available from each post is available to put into your jinja
template.  These can either be placed there in your post frontmatter, or
through a plugin that automatically adds to the post before the save phase.

Here is a very simple example that would give a link to each post with the
title and date.

``` toml
[[markata.feeds]]
slug='all'
title='All Posts'
filter="True"
card_template='''
<li>
    <a href={{markata.config.get('path_prefix', '')}}{{slug}}>
        {{title}}-{{date}}
    </a>
</li>
'''
```

# filter

The filter is a python expression ran on every post that expects to return a
boolean.  The variables available to this expression are every key in your
frontmatter, plus the `timedelta` function, and `parse` function to more easily
work with dates.

# Feed Examples

True can be passed in to make a feed of all the posts you have.

``` toml
[[markata.feeds]]
slug='all'
title='All Posts'
filter="True"
```

You can compare against the values of the keys from your frontmatter.  This
example creates a feed that includes every post where published is `True`.

``` toml
[[markata.feeds]]
slug='draft'
title='Draft'
filter="published=='False'"
```

We can also compare against dates.  The
[markata.plugins.datetime](https://markata.dev/markata/plugins/datetime/)
plugin, automatically adds `today` as today's date and `now` as the current
datetime.  These are quite handy to create feeds for scheduled, recent, or
today's posts.  The following two examples will create a feed for scheduled
posts and for today's posts respectively.

``` toml
[[markata.feeds]]
slug='scheduled'
title='Scheduled'
filter="date>today"

[[markata.feeds]]
slug='today'
title='Today'
filter="date==today"
```

If you have list of items in your frontmatter for something like `tags`, you
can check for the existence of a tag in the list.

``` toml
[[markata.feeds]]
slug='python'
title='Python'
filter="date<=today and 'python' in tags"
```

And of course you can combine all the things into larger expressions.  Here is
one example of the main feed on my blog.

``` toml
[[markata.feeds]]
slug='blog'
title='Blog'
filter="date<=today and templateKey in ['blog-post'] and published =='True'"
```

Here is another example that shows my drafts for a particular tag.

``` toml
[[markata.feeds]]
slug='python-draft'
title='Python Draft'
filter="date<=today and 'python' in tags and published=='False'"
```

# Defaults

By default feeds will create one feed page at `/archive/` that includes all
posts.

[[markata.feeds]]
slug='archive'
title='All Posts'
filter="True"


!! class <h2 id='SilentUndefined' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>SilentUndefined <em class='small'>class</em></h2>

???+ source "SilentUndefined <em class='small'>source</em>"

```python

        class SilentUndefined(Undefined):
            def _fail_with_undefined_error(self, *args, **kwargs):
                return ""
```


!! class <h2 id='MarkataFilterError' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>MarkataFilterError <em class='small'>class</em></h2>

???+ source "MarkataFilterError <em class='small'>source</em>"

```python

        class MarkataFilterError(RuntimeError):
            ...
```


!! class <h2 id='FeedConfig' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>FeedConfig <em class='small'>class</em></h2>

???+ source "FeedConfig <em class='small'>source</em>"

```python

        class FeedConfig(pydantic.BaseModel):
            DEFAULT_TITLE: str = "All Posts"

            title: str = DEFAULT_TITLE
            slug: str = None
            name: Optional[str] = None
            filter: str = "True"
            sort: str = "date"
            reverse: bool = False
            rss: bool = True
            sitemap: bool = True
            card_template: str = """
                <li class='post'>
                    <a href="/{{ markata.config.path_prefix }}{{ post.slug }}/">
                        {{ post.title }}
                    </a>
                </li>
                """
            template: str = Path(__file__).parent / "default_post_template.html.jinja"
            rss_template: str = Path(__file__).parent / "default_rss_template.xml"
            sitemap_template: str = Path(__file__).parent / "default_sitemap_template.xml"
            xsl_template: str = Path(__file__).parent / "default_xsl_template.xsl"

            @pydantic.validator("name", pre=True, always=True)
            def default_name(cls, v, *, values):
                return v or str(values.get("slug")).replace("-", "_")

            @pydantic.validator("card_template", "template", pre=True, always=True)
            def read_template(cls, v, *, values) -> str:
                if isinstance(v, Path):
                    return str(v.read_text())
                return v
```


!! class <h2 id='FeedsConfig' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>FeedsConfig <em class='small'>class</em></h2>

???+ source "FeedsConfig <em class='small'>source</em>"

```python

        class FeedsConfig(pydantic.BaseModel):
            feeds: List[FeedConfig] = [FeedConfig(slug="archive")]
```


!! class <h2 id='Feed' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Feed <em class='small'>class</em></h2>
    A storage class for markata feed objects.

    # Usage

    ``` python
    from markata import Markata
    m = Markata()

    # access posts for a feed
    m.feeds.docs.posts

    # access config for a feed
    m.feeds.docs.config
    ```
???+ source "Feed <em class='small'>source</em>"

```python

        class Feed:
            """
            A storage class for markata feed objects.

            # Usage

            ``` python
            from markata import Markata
            m = Markata()

            # access posts for a feed
            m.feeds.docs.posts

            # access config for a feed
            m.feeds.docs.config
            ```
            """

            config: FeedConfig
            _m: Markata

            @property
            def name(self):
                return self.config.name

            @property
            def posts(self):
                return self.map("post")

            def map(self, func="post", **args):
                return self._m.map(func, **{**self.config.dict(), **args})
```


!! class <h2 id='Feeds' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Feeds <em class='small'>class</em></h2>
    A storage class for all markata Feed objects

    ``` python
    from markata import Markata
    m = Markata()

    m.feeds

    # access all config
    m.feeds.config

    # refresh list of posts in all feeds
    m.feeds.refresh()


    # iterating over feeds gives the name of the feed
    for k in m.feeds:
         print(k)

    # project-gallery
    # docs
    # autodoc
    # core_modules
    # plugins
    # archive

    # iterate over items like keys and values in a dict, items returns name of
    # feed and a feed object
    for k, v in m.feeds.items():
        print(k, len(v.posts))

    # project-gallery 2
    # docs 6
    # autodoc 65
    # core_modules 26
    # plugins 39
    # archive 65

    # values can be iterated over in just the same way
    for v in m.feeds.values():
         print(len(v.posts))
    # 2
    # 6
    # 65
    # 26
    # 39
    # 65
    ```

    Accessing feeds can be done using square brackets or dot notation.

    ``` python
    from markata import Markata
    m = Markata()

    # both of these will return the `docs` Feed object.
    m.feeds.docs
    m['docs']
    ```
???+ source "Feeds <em class='small'>source</em>"

```python

        class Feeds:
            """
            A storage class for all markata Feed objects

            ``` python
            from markata import Markata
            m = Markata()

            m.feeds

            # access all config
            m.feeds.config

            # refresh list of posts in all feeds
            m.feeds.refresh()


            # iterating over feeds gives the name of the feed
            for k in m.feeds:
                 print(k)

            # project-gallery
            # docs
            # autodoc
            # core_modules
            # plugins
            # archive

            # iterate over items like keys and values in a dict, items returns name of
            # feed and a feed object
            for k, v in m.feeds.items():
                print(k, len(v.posts))

            # project-gallery 2
            # docs 6
            # autodoc 65
            # core_modules 26
            # plugins 39
            # archive 65

            # values can be iterated over in just the same way
            for v in m.feeds.values():
                 print(len(v.posts))
            # 2
            # 6
            # 65
            # 26
            # 39
            # 65
            ```

            Accessing feeds can be done using square brackets or dot notation.

            ``` python
            from markata import Markata
            m = Markata()

            # both of these will return the `docs` Feed object.
            m.feeds.docs
            m['docs']
            ```
            """

            def __init__(self, markata: Markata) -> None:
                self._m = markata
                self.config = {f.name: f for f in markata.config.feeds}
                self.refresh()

            def refresh(self) -> None:
                """
                Refresh all of the feeds objects
                """
                for feed in self._m.config.feeds:
                    feed = Feed(config=feed, _m=self._m)
                    self.__setattr__(feed.name, feed)

            def __iter__(self):
                return iter(self.config.keys())

            def keys(self):
                return iter(self.config.keys())

            def values(self):
                return [self[feed] for feed in self.config.keys()]

            def items(self):
                return [(key, self[key]) for key in self.config]

            def __getitem__(self, key: str) -> Any:
                return getattr(self, key.replace("-", "_").lower())

            def _dict_panel(self, config) -> str:
                """
                pretty print configs with rich
                """
                msg = ""
                for key, value in config.items():
                    if isinstance(value, str):
                        if len(value) > 50:
                            value = value[:50] + "..."
                        value = value
                    msg = msg + f"[grey46]{key}[/][magenta3]:[/] [grey66]{value}[/]\n"
                return msg

            def __rich__(self) -> Table:
                from rich.table import Table

                table = Table(title=f"Feeds {len(self.config)}")

                table.add_column("Feed", justify="right", style="cyan", no_wrap=True)
                table.add_column("posts", justify="left", style="green")
                table.add_column("config", style="magenta")

                for name in self.config:
                    table.add_row(
                        name,
                        str(len(self[name].posts)),
                        self._dict_panel(self.config[name].dict()),
                    )
                return table
```


!! function <h2 id='config_model' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>config_model <em class='small'>function</em></h2>

???+ source "config_model <em class='small'>source</em>"

```python

        def config_model(markata: Markata) -> None:
            markata.config_models.append(FeedsConfig)
```


!! function <h2 id='pre_render' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>pre_render <em class='small'>function</em></h2>
    Create the Feeds object and attach it to markata.
???+ source "pre_render <em class='small'>source</em>"

```python

        def pre_render(markata: Markata) -> None:
            """
            Create the Feeds object and attach it to markata.
            """
            markata.feeds = Feeds(markata)
```


!! function <h2 id='save' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>save <em class='small'>function</em></h2>
    Creates a new feed page for each page in the config.
???+ source "save <em class='small'>source</em>"

```python

        def save(markata: Markata) -> None:
            """
            Creates a new feed page for each page in the config.
            """
            with markata.cache as cache:
                for feed in markata.feeds.values():
                    create_page(
                        markata,
                        feed,
                        cache,
                    )

            home = Path(str(markata.config.output_dir)) / "index.html"
            archive = Path(str(markata.config.output_dir)) / "archive" / "index.html"
            if not home.exists() and archive.exists():
                shutil.copy(str(archive), str(home))

            xsl_template = get_template(feed.config.xsl_template)
            xsl = xsl_template.render(
                markata=markata,
                __version__=__version__,
                today=datetime.datetime.today(),
                config=markata.config,
            )
            xsl_file = Path(markata.config.output_dir) / "rss.xsl"
            xsl_file.write_text(xsl)
```


!! function <h2 id='get_template' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get_template <em class='small'>function</em></h2>

???+ source "get_template <em class='small'>source</em>"

```python

        def get_template(src) -> Template:
            try:
                return Template(Path(src).read_text(), undefined=SilentUndefined)
            except FileNotFoundError:
                return Template(src, undefined=SilentUndefined)
            except OSError:  # File name too long
                return Template(src, undefined=SilentUndefined)
```


!! function <h2 id='create_page' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>create_page <em class='small'>function</em></h2>
    create an html unorderd list of posts.
???+ source "create_page <em class='small'>source</em>"

```python

        def create_page(
            markata: Markata,
            feed: Feed,
            cache,
        ) -> None:
            """
            create an html unorderd list of posts.
            """

            posts = feed.posts

            cards = [
                create_card(markata, post, feed.config.card_template, cache) for post in posts
            ]
            cards.insert(0, "<ul>")
            cards.append("</ul>")
            cards = "".join(cards)

            template = get_template(feed.config.template)
            rss_template = get_template(feed.config.rss_template)
            sitemap_template = get_template(feed.config.sitemap_template)
            output_file = Path(markata.config.output_dir) / feed.config.slug / "index.html"
            canonical_url = f"{markata.config.url}/{feed.config.slug}/"
            output_file.parent.mkdir(exist_ok=True, parents=True)

            rss_output_file = Path(markata.config.output_dir) / feed.config.slug / "rss.xml"
            rss_output_file.parent.mkdir(exist_ok=True, parents=True)

            sitemap_output_file = (
                Path(markata.config.output_dir) / feed.config.slug / "sitemap.xml"
            )
            sitemap_output_file.parent.mkdir(exist_ok=True, parents=True)

            key = markata.make_hash(
                "feeds",
                template,
                __version__,
                cards,
                markata.config.url,
                markata.config.description,
                feed.config.title,
                canonical_url,
                datetime.datetime.today(),
                markata.config,
            )

            feed_html_from_cache = markata.precache.get(key)
            if feed_html_from_cache is None:
                feed_html = template.render(
                    markata=markata,
                    __version__=__version__,
                    body=cards,
                    url=markata.config.url,
                    description=markata.config.description,
                    title=feed.config.title,
                    canonical_url=canonical_url,
                    today=datetime.datetime.today(),
                    config=markata.config,
                )
                with markata.cache as cache:
                    markata.cache.set(key, feed_html)

            feed_rss = rss_template.render(markata=markata, feed=feed)
            feed_sitemap = sitemap_template.render(markata=markata, feed=feed)

            output_file.write_text(feed_html)
            rss_output_file.write_text(feed_rss)
            sitemap_output_file.write_text(feed_sitemap)
```


!! function <h2 id='create_card' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>create_card <em class='small'>function</em></h2>
    Creates a card for one post based on the configured template.  If no
    template is configured it will create one with the post title and dates
    (if present).
???+ source "create_card <em class='small'>source</em>"

```python

        def create_card(
            markata: "Markata",
            post: "Post",
            template: Optional[str] = None,
            cache=None,
        ) -> Any:
            """
            Creates a card for one post based on the configured template.  If no
            template is configured it will create one with the post title and dates
            (if present).
            """
            key = markata.make_hash("feeds", template, str(post), post.content)

            card = markata.precache.get(key)
            if card is not None:
                return card

            if template is None:
                template = markata.config.get("feeds_config", {}).get("card_template", None)

            if template is None:
                if "date" in post:
                    card = textwrap.dedent(
                        f"""
                        <li class='post'>
                        <a href="/{markata.config.path_prefix}{post.slug}/">
                            {post.title}
                            {post.date.year}-
                            {post.date.month}-
                            {post.date.day}
                        </a>
                        </li>
                        """,
                    )
                else:
                    card = textwrap.dedent(
                        f"""
                        <li class='post'>
                        <a href="/{markata.config.path_prefix}{post.slug}/">
                            {post.title}
                        </a>
                        </li>
                        """,
                    )
            else:
                try:
                    _template = Template(Path(template).read_text())
                except FileNotFoundError:
                    _template = Template(template)
                except OSError:  # File name too long
                    _template = Template(template)
                card = _template.render(post=post, **post.to_dict())
            cache.add(key, card)
            return card
```


!! function <h2 id='cli' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>cli <em class='small'>function</em></h2>

???+ source "cli <em class='small'>source</em>"

```python

        def cli(app: typer.Typer, markata: "Markata") -> None:
            feeds_app = typer.Typer()
            app.add_typer(feeds_app)

            @feeds_app.callback()
            def feeds():
                "feeds cli"

            @feeds_app.command()
            def show() -> None:
                markata.console.quiet = True
                feeds = markata.feeds
                markata.console.quiet = False
                rich_print(feeds)
```


!! method <h2 id='_fail_with_undefined_error' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_fail_with_undefined_error <em class='small'>method</em></h2>

???+ source "_fail_with_undefined_error <em class='small'>source</em>"

```python

        def _fail_with_undefined_error(self, *args, **kwargs):
                return ""
```


!! method <h2 id='default_name' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>default_name <em class='small'>method</em></h2>

???+ source "default_name <em class='small'>source</em>"

```python

        def default_name(cls, v, *, values):
                return v or str(values.get("slug")).replace("-", "_")
```


!! method <h2 id='read_template' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>read_template <em class='small'>method</em></h2>

???+ source "read_template <em class='small'>source</em>"

```python

        def read_template(cls, v, *, values) -> str:
                if isinstance(v, Path):
                    return str(v.read_text())
                return v
```


!! method <h2 id='name' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>name <em class='small'>method</em></h2>

???+ source "name <em class='small'>source</em>"

```python

        def name(self):
                return self.config.name
```


!! method <h2 id='posts' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>posts <em class='small'>method</em></h2>

???+ source "posts <em class='small'>source</em>"

```python

        def posts(self):
                return self.map("post")
```


!! method <h2 id='map' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>map <em class='small'>method</em></h2>

???+ source "map <em class='small'>source</em>"

```python

        def map(self, func="post", **args):
                return self._m.map(func, **{**self.config.dict(), **args})
```


!! method <h2 id='__init__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__init__ <em class='small'>method</em></h2>

???+ source "__init__ <em class='small'>source</em>"

```python

        def __init__(self, markata: Markata) -> None:
                self._m = markata
                self.config = {f.name: f for f in markata.config.feeds}
                self.refresh()
```


!! method <h2 id='refresh' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>refresh <em class='small'>method</em></h2>
    Refresh all of the feeds objects
???+ source "refresh <em class='small'>source</em>"

```python

        def refresh(self) -> None:
                """
                Refresh all of the feeds objects
                """
                for feed in self._m.config.feeds:
                    feed = Feed(config=feed, _m=self._m)
                    self.__setattr__(feed.name, feed)
```


!! method <h2 id='__iter__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__iter__ <em class='small'>method</em></h2>

???+ source "__iter__ <em class='small'>source</em>"

```python

        def __iter__(self):
                return iter(self.config.keys())
```


!! method <h2 id='keys' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>keys <em class='small'>method</em></h2>

???+ source "keys <em class='small'>source</em>"

```python

        def keys(self):
                return iter(self.config.keys())
```


!! method <h2 id='values' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>values <em class='small'>method</em></h2>

???+ source "values <em class='small'>source</em>"

```python

        def values(self):
                return [self[feed] for feed in self.config.keys()]
```


!! method <h2 id='items' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>items <em class='small'>method</em></h2>

???+ source "items <em class='small'>source</em>"

```python

        def items(self):
                return [(key, self[key]) for key in self.config]
```


!! method <h2 id='__getitem__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__getitem__ <em class='small'>method</em></h2>

???+ source "__getitem__ <em class='small'>source</em>"

```python

        def __getitem__(self, key: str) -> Any:
                return getattr(self, key.replace("-", "_").lower())
```


!! method <h2 id='_dict_panel' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_dict_panel <em class='small'>method</em></h2>
    pretty print configs with rich
???+ source "_dict_panel <em class='small'>source</em>"

```python

        def _dict_panel(self, config) -> str:
                """
                pretty print configs with rich
                """
                msg = ""
                for key, value in config.items():
                    if isinstance(value, str):
                        if len(value) > 50:
                            value = value[:50] + "..."
                        value = value
                    msg = msg + f"[grey46]{key}[/][magenta3]:[/] [grey66]{value}[/]\n"
                return msg
```


!! method <h2 id='__rich__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__rich__ <em class='small'>method</em></h2>

???+ source "__rich__ <em class='small'>source</em>"

```python

        def __rich__(self) -> Table:
                from rich.table import Table

                table = Table(title=f"Feeds {len(self.config)}")

                table.add_column("Feed", justify="right", style="cyan", no_wrap=True)
                table.add_column("posts", justify="left", style="green")
                table.add_column("config", style="magenta")

                for name in self.config:
                    table.add_row(
                        name,
                        str(len(self[name].posts)),
                        self._dict_panel(self.config[name].dict()),
                    )
                return table
```


!! function <h2 id='feeds' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>feeds <em class='small'>function</em></h2>
    feeds cli
???+ source "feeds <em class='small'>source</em>"

```python

        def feeds():
                "feeds cli"
```


!! function <h2 id='show' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>show <em class='small'>function</em></h2>

???+ source "show <em class='small'>source</em>"

```python

        def show() -> None:
                markata.console.quiet = True
                feeds = markata.feeds
                markata.console.quiet = False
                rich_print(feeds)
```

