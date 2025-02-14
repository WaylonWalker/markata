---
content: "---\n\nThe `markata.plugins.feeds` plugin is used to create feed pages,
  which are lists of\nposts.  The list is generated using a `filter`, then each post
  in the list is\nrendered with a `card_template` before being applied to the `body`
  of the\n`template`.\n\n# Installation\n\nThis plugin is built-in and enabled by
  default, but in you want to be very\nexplicit you can add it to your list of existing
  plugins.\n\n``` toml\nhooks = [\n   \"markata.plugins.feeds\",\n   ]\n```\n\n# Configuration\n\n#
  set default template and card_template\n\nAt the root of the markata.feeds config
  you may set `template`, and\n`card_template`.  These will become your defaults for
  every feed you create.\nIf you do not set these, markata will use it's defaults.
  \ The defaults are\ndesigned to work for a variety of use cases, but are not likely
  the best for\nall.\n\n``` toml\n[markata.feeds_config]\ntemplate=\"pages/templates/archive_template.html\"\ncard_template=\"plugins/feed_card_template.html\"\n```\n\n#
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
  all\nposts.\n\n``` toml\n[[markata.feeds]]\nslug='archive'\ntitle='All Posts'\nfilter=\"True\"\n```\n\n---\n\n\n\n\n\n\n\n\n\n!!!
  class\n    <h2 id=\"Feed\" class=\"admonition-title\" style=\"margin: 0; padding:
  .5rem 1rem;\">Feed <em class=\"small\">class</em></h2>\n\n    A storage class for
  markata feed objects.\n\n    ## Usage\n\n    ``` python\n    from markata import
  Markata\n    m = Markata()\n\n    # access posts for a feed\n    m.feeds.docs.posts\n\n
  \   # access config for a feed\n    m.feeds.docs.config\n    ```\n\n???+ source
  \"Feed <em class='small'>source</em>\"\n    ```python\n    class Feed(pydantic.BaseModel,
  JupyterMixin):\n        \"\"\"A storage class for markata feed objects.\n\n        ##
  Usage\n\n        ``` python\n        from markata import Markata\n        m = Markata()\n\n
  \       # access posts for a feed\n        m.feeds.docs.posts\n\n        # access
  config for a feed\n        m.feeds.docs.config\n        ```\n        \"\"\"\n\n
  \       config: FeedConfig\n        markata: Markata = pydantic.Field(exclude=True)\n\n
  \       model_config = ConfigDict(\n            validate_assignment=False,\n            arbitrary_types_allowed=True,\n
  \           extra=\"allow\",\n            str_strip_whitespace=True,\n            validate_default=True,\n
  \           coerce_numbers_to_str=True,\n            populate_by_name=True,\n        )\n\n
  \       @property\n        def name(self) -> str:\n            \"\"\"The name of
  the feed, used for accessing it in the feeds object.\"\"\"\n            return self.config.name\n\n
  \       @property\n        def posts(self):\n            posts = self.map(\"post\")\n
  \           if self.config.head is not None and self.config.tail is not None:\n
  \               head_posts = posts[: self.config.head]\n                tail_posts
  = posts[-self.config.tail :]\n                return PrettyList(head_posts + tail_posts)\n
  \           if self.config.head is not None:\n                return PrettyList(posts[:
  self.config.head])\n            if self.config.tail is not None:\n                return
  PrettyList(posts[-self.config.tail :])\n            return PrettyList(posts)\n\n
  \       def first(\n            self: \"Markata\",\n        ) -> list:\n            return
  self.posts[0]\n\n        def last(\n            self: \"Markata\",\n        ) ->
  list:\n            return self.posts[-1]\n\n        def map(self, func=\"post\",
  **args):\n            return self.markata.map(func, **{**self.config.dict(), **args})\n\n
  \       @property\n        def __rich_console__(self) -> \"Console\":\n            return
  self.markata.console\n\n        def __rich__(self) -> Table:\n            table
  = Table(title=f\"Feed: {self.name}\")\n\n            table.add_column(\"Post\",
  justify=\"right\", style=\"cyan\", no_wrap=True)\n            table.add_column(\"slug\",
  justify=\"left\", style=\"green\")\n            table.add_column(\"published\",
  justify=\"left\", style=\"green\")\n\n            for post in self.posts:\n                table.add_row(post.title,
  post.slug, str(post.published))\n\n            return table\n    ```\n\n\n\n\n\n\n\n\n\n!!!
  function\n    <h2 id=\"pre_render\" class=\"admonition-title\" style=\"margin: 0;
  padding: .5rem 1rem;\">pre_render <em class=\"small\">function</em></h2>\n\n    Create
  the Feeds object and attach it to markata.\n\n???+ source \"pre_render <em class='small'>source</em>\"\n
  \   ```python\n    def pre_render(markata: Markata) -> None:\n        \"\"\"\n        Create
  the Feeds object and attach it to markata.\n        \"\"\"\n        markata.feeds
  = Feeds(markata)\n    ```\n\n\n\n\n\n!!! function\n    <h2 id=\"save\" class=\"admonition-title\"
  style=\"margin: 0; padding: .5rem 1rem;\">save <em class=\"small\">function</em></h2>\n\n
  \   Creates a new feed page for each page in the config.\n\n???+ source \"save <em
  class='small'>source</em>\"\n    ```python\n    def save(markata: Markata) -> None:\n
  \       \"\"\"\n        Creates a new feed page for each page in the config.\n        \"\"\"\n
  \       with markata.cache as cache:\n            for feed in markata.feeds.values():\n
  \               create_page(\n                    markata,\n                    feed,\n
  \                   cache,\n                )\n\n        home = Path(str(markata.config.output_dir))
  / \"index.html\"\n        archive = Path(str(markata.config.output_dir)) / \"archive\"
  / \"index.html\"\n        if not home.exists() and archive.exists():\n            shutil.copy(str(archive),
  str(home))\n\n        xsl_template = get_template(markata, feed.config.xsl_template)\n
  \       xsl = xsl_template.render(\n            markata=markata,\n            __version__=__version__,\n
  \           today=datetime.datetime.today(),\n            config=markata.config,\n
  \       )\n        xsl_file = Path(markata.config.output_dir) / \"rss.xsl\"\n        xsl_file.write_text(xsl)\n
  \   ```\n\n\n\n!!! function\n    <h2 id=\"create_page\" class=\"admonition-title\"
  style=\"margin: 0; padding: .5rem 1rem;\">create_page <em class=\"small\">function</em></h2>\n\n
  \   create an html unorderd list of posts.\n\n???+ source \"create_page <em class='small'>source</em>\"\n
  \   ```python\n    def create_page(\n        markata: Markata,\n        feed: Feed,\n
  \       cache,\n    ) -> None:\n        \"\"\"\n        create an html unorderd
  list of posts.\n        \"\"\"\n\n        template = get_template(markata, feed.config.template)\n
  \       partial_template = get_template(markata, feed.config.partial_template)\n
  \       canonical_url = f\"{markata.config.url}/{feed.config.slug}/\"\n\n        key
  = markata.make_hash(\n            \"feeds\",\n            template,\n            __version__,\n
  \           # cards,\n            markata.config.url,\n            markata.config.description,\n
  \           feed.config.title,\n            feed.map(\"content\"),\n            canonical_url,\n
  \           # datetime.datetime.today(),\n            # markata.config,\n        )\n\n
  \       html_key = markata.make_hash(key, \"html\")\n        html_partial_key =
  markata.make_hash(key, \"partial_html\")\n        feed_rss_key = markata.make_hash(key,
  \"rss\")\n        feed_sitemap_key = markata.make_hash(key, \"sitemap\")\n\n        feed_html_from_cache
  = markata.precache.get(html_key)\n        feed_html_partial_from_cache = markata.precache.get(html_partial_key)\n
  \       feed_rss_from_cache = markata.precache.get(feed_rss_key)\n        feed_sitemap_from_cache
  = markata.precache.get(feed_sitemap_key)\n\n        output_file = Path(markata.config.output_dir)
  / feed.config.slug / \"index.html\"\n        output_file.parent.mkdir(exist_ok=True,
  parents=True)\n\n        partial_output_file = (\n            Path(markata.config.output_dir)
  / feed.config.slug / \"partial\" / \"index.html\"\n        )\n        partial_output_file.parent.mkdir(exist_ok=True,
  parents=True)\n\n        rss_output_file = Path(markata.config.output_dir) / feed.config.slug
  / \"rss.xml\"\n        rss_output_file.parent.mkdir(exist_ok=True, parents=True)\n\n
  \       sitemap_output_file = (\n            Path(markata.config.output_dir) / feed.config.slug
  / \"sitemap.xml\"\n        )\n        sitemap_output_file.parent.mkdir(exist_ok=True,
  parents=True)\n\n        from_cache = True\n        if feed_html_from_cache is None:\n
  \           from_cache = False\n            feed_html = template.render(\n                markata=markata,\n
  \               __version__=__version__,\n                post=feed.config.model_dump(),\n
  \               url=markata.config.url,\n                config=markata.config,\n
  \               feed=feed,\n            )\n            cache.set(html_key, feed_html)\n
  \       else:\n            feed_html = feed_html_from_cache\n\n        if feed_html_partial_from_cache
  is None:\n            from_cache = False\n            feed_html_partial = partial_template.render(\n
  \               markata=markata,\n                __version__=__version__,\n                post=feed.config.model_dump(),\n
  \               url=markata.config.url,\n                config=markata.config,\n
  \               feed=feed,\n            )\n            cache.set(html_partial_key,
  feed_html_partial)\n        else:\n            feed_html_partial = feed_html_partial_from_cache\n\n
  \       if feed_rss_from_cache is None:\n            from_cache = False\n            rss_template
  = get_template(markata, feed.config.rss_template)\n            feed_rss = rss_template.render(markata=markata,
  feed=feed)\n            cache.set(feed_rss_key, feed_rss)\n        else:\n            feed_rss
  = feed_rss_from_cache\n\n        if feed_sitemap_from_cache is None:\n            from_cache
  = False\n            sitemap_template = get_template(markata, feed.config.sitemap_template)\n
  \           feed_sitemap = sitemap_template.render(markata=markata, feed=feed)\n
  \           cache.set(feed_sitemap_key, feed_sitemap)\n        else:\n            feed_sitemap
  = feed_sitemap_from_cache\n\n        if (\n            from_cache\n            and
  output_file.exists()\n            and partial_output_file.exists()\n            and
  rss_output_file.exists()\n            and sitemap_output_file.exists()\n        ):\n
  \           return\n\n        output_file.write_text(feed_html)\n        partial_output_file.write_text(feed_html_partial)\n
  \       rss_output_file.write_text(feed_rss)\n        sitemap_output_file.write_text(feed_sitemap)\n
  \   ```\n\n\n\n!!! function\n    <h2 id=\"create_card\" class=\"admonition-title\"
  style=\"margin: 0; padding: .5rem 1rem;\">create_card <em class=\"small\">function</em></h2>\n\n
  \   Creates a card for one post based on the configured template.  If no\n    template
  is configured it will create one with the post title and dates\n    (if present).\n\n???+
  source \"create_card <em class='small'>source</em>\"\n    ```python\n    def create_card(\n
  \       markata: \"Markata\",\n        post: \"Post\",\n        template: Optional[str]
  = None,\n        cache=None,\n    ) -> Any:\n        \"\"\"\n        Creates a card
  for one post based on the configured template.  If no\n        template is configured
  it will create one with the post title and dates\n        (if present).\n        \"\"\"\n
  \       key = markata.make_hash(\"feeds\", template, str(post), post.content)\n\n
  \       card = markata.precache.get(key)\n        if card is not None:\n            return
  card\n\n        if template is None:\n            template = markata.config.get(\"feeds_config\",
  {}).get(\"card_template\", None)\n\n        if template is None:\n            if
  \"date\" in post:\n                card = textwrap.dedent(\n                    f\"\"\"\n
  \                   <li class='post'>\n                    <a href=\"/{markata.config.path_prefix}{post.slug}/\">\n
  \                       {post.title}\n                        {post.date.year}-\n
  \                       {post.date.month}-\n                        {post.date.day}\n
  \                   </a>\n                    </li>\n                    \"\"\",\n
  \               )\n            else:\n                card = textwrap.dedent(\n
  \                   f\"\"\"\n                    <li class='post'>\n                    <a
  href=\"/{markata.config.path_prefix}{post.slug}/\">\n                        {post.title}\n
  \                   </a>\n                    </li>\n                    \"\"\",\n
  \               )\n        else:\n            try:\n                _template =
  Template(Path(template).read_text())\n            except FileNotFoundError:\n                _template
  = Template(template)\n            except OSError:  # File name too long\n                _template
  = Template(template)\n            card = _template.render(post=post, **post.to_dict())\n
  \       cache.add(key, card)\n        return card\n    ```\n\n\n\n\n\n!!! class\n
  \   <h2 id=\"Feeds\" class=\"admonition-title\" style=\"margin: 0; padding: .5rem
  1rem;\">Feeds <em class=\"small\">class</em></h2>\n\n    A storage class for all
  markata Feed objects\n\n    ``` python\n    from markata import Markata\n    markata
  = Markata()\n\n    markata.feeds\n\n    # access all config\n    markata.feeds.config\n\n
  \   # refresh list of posts in all feeds\n    markata.feeds.refresh()\n\n    # iterating
  over feeds gives the name of the feed\n    for k in markata.feeds:\n         print(k)\n\n
  \   # project-gallery\n    # docs\n    # autodoc\n    # core_modules\n    # plugins\n
  \   # archive\n\n    # iterate over items like keys and values in a dict, items
  returns name of\n    # feed and a feed object\n    for k, v in markata.feeds.items():\n
  \       print(k, len(v.posts))\n\n    # project-gallery 2\n    # docs 6\n    # autodoc
  65\n    # core_modules 26\n    # plugins 39\n    # archive 65\n\n    # values can
  be iterated over in just the same way\n    for v in markata.feeds.values():\n         print(len(v.posts))\n
  \   # 2\n    # 6\n    # 65\n    # 26\n    # 39\n    # 65\n\n???+ source \"Feeds
  <em class='small'>source</em>\"\n    ```python\n    class Feeds(JupyterMixin):\n
  \       \"\"\"A storage class for all markata Feed objects\n\n        ``` python\n
  \       from markata import Markata\n        markata = Markata()\n\n        markata.feeds\n\n
  \       # access all config\n        markata.feeds.config\n\n        # refresh list
  of posts in all feeds\n        markata.feeds.refresh()\n\n        # iterating over
  feeds gives the name of the feed\n        for k in markata.feeds:\n             print(k)\n\n
  \       # project-gallery\n        # docs\n        # autodoc\n        # core_modules\n
  \       # plugins\n        # archive\n\n        # iterate over items like keys and
  values in a dict, items returns name of\n        # feed and a feed object\n        for
  k, v in markata.feeds.items():\n            print(k, len(v.posts))\n\n        #
  project-gallery 2\n        # docs 6\n        # autodoc 65\n        # core_modules
  26\n        # plugins 39\n        # archive 65\n\n        # values can be iterated
  over in just the same way\n        for v in markata.feeds.values():\n             print(len(v.posts))\n
  \       # 2\n        # 6\n        # 65\n        # 26\n        # 39\n        # 65\n
  \       \"\"\"\n\n        def __init__(self, markata: Markata) -> None:\n            self.markata
  = markata\n            self.config = {f.name: f for f in markata.config.feeds}\n
  \           self.refresh()\n\n        def refresh(self):\n            \"\"\"Refresh
  all of the feeds objects\"\"\"\n            for feed_config in self.markata.config.feeds:\n
  \               # Ensure feed has a name, falling back to slug if needed\n                if
  feed_config.name is None and feed_config.slug is not None:\n                    feed_config.name
  = feed_config.slug.replace(\"-\", \"_\")\n                elif feed_config.name
  is None and feed_config.slug is None:\n                    feed_config.slug = \"archive\"\n
  \                   feed_config.name = \"archive\"\n\n                feed = Feed(config=feed_config,
  markata=self.markata)\n                self.__setattr__(feed.name, feed)\n\n        def
  __iter__(self):\n            return iter(self.config.keys())\n\n        def keys(self):\n
  \           return iter(self.config.keys())\n\n        def values(self):\n            return
  [self[feed] for feed in self.config.keys()]\n\n        def items(self):\n            return
  [(key, self[key]) for key in self.config]\n\n        def __getitem__(self, key:
  str) -> Any:\n            return getattr(self, key.replace(\"-\", \"_\").lower())\n\n
  \       def get(self, key: str, default: Any = None) -> Any:\n            return
  getattr(self, key.replace(\"-\", \"_\").lower(), default)\n\n        def _dict_panel(self,
  config) -> str:\n            \"\"\"pretty print configs with rich\"\"\"\n            msg
  = \"\"\n            for key, value in config.items():\n                if isinstance(value,
  str):\n                    if len(value) > 50:\n                        value =
  value[:50] + \"...\"\n                    value = value\n                msg = msg
  + f\"[grey46]{key}[/][magenta3]:[/] [grey66]{value}[/]\\n\"\n            return
  msg\n\n        def __rich__(self) -> Table:\n            table = Table(title=f\"Feeds
  {len(self.config)}\")\n\n            table.add_column(\"Feed\", justify=\"right\",
  style=\"cyan\", no_wrap=True)\n            table.add_column(\"posts\", justify=\"left\",
  style=\"green\")\n            table.add_column(\"config\", style=\"magenta\")\n\n
  \           for name in self.config:\n                table.add_row(\n                    name,\n
  \                   str(len(self[name].posts)),\n                    self._dict_panel(self.config[name].dict()),\n
  \               )\n            return table\n    ```\n\n\n\n\n\n\n\n\n\n\n\n\n\n!!!
  method\n    <h2 id=\"name\" class=\"admonition-title\" style=\"margin: 0; padding:
  .5rem 1rem;\">name <em class=\"small\">method</em></h2>\n\n    The name of the feed,
  used for accessing it in the feeds object.\n\n???+ source \"name <em class='small'>source</em>\"\n
  \   ```python\n    def name(self) -> str:\n            \"\"\"The name of the feed,
  used for accessing it in the feeds object.\"\"\"\n            return self.config.name\n
  \   ```\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n!!! function\n    <h2 id=\"feeds\"
  class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">feeds <em class=\"small\">function</em></h2>\n\n
  \   feeds cli\n\n???+ source \"feeds <em class='small'>source</em>\"\n    ```python\n
  \   def feeds():\n            \"feeds cli\"\n    ```\n\n\n\n\n\n\n\n!!! method\n
  \   <h2 id=\"refresh\" class=\"admonition-title\" style=\"margin: 0; padding: .5rem
  1rem;\">refresh <em class=\"small\">method</em></h2>\n\n    Refresh all of the feeds
  objects\n\n???+ source \"refresh <em class='small'>source</em>\"\n    ```python\n
  \   def refresh(self):\n            \"\"\"Refresh all of the feeds objects\"\"\"\n
  \           for feed_config in self.markata.config.feeds:\n                # Ensure
  feed has a name, falling back to slug if needed\n                if feed_config.name
  is None and feed_config.slug is not None:\n                    feed_config.name
  = feed_config.slug.replace(\"-\", \"_\")\n                elif feed_config.name
  is None and feed_config.slug is None:\n                    feed_config.slug = \"archive\"\n
  \                   feed_config.name = \"archive\"\n\n                feed = Feed(config=feed_config,
  markata=self.markata)\n                self.__setattr__(feed.name, feed)\n    ```\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n!!!
  method\n    <h2 id=\"_dict_panel\" class=\"admonition-title\" style=\"margin: 0;
  padding: .5rem 1rem;\">_dict_panel <em class=\"small\">method</em></h2>\n\n    pretty
  print configs with rich\n\n???+ source \"_dict_panel <em class='small'>source</em>\"\n
  \   ```python\n    def _dict_panel(self, config) -> str:\n            \"\"\"pretty
  print configs with rich\"\"\"\n            msg = \"\"\n            for key, value
  in config.items():\n                if isinstance(value, str):\n                    if
  len(value) > 50:\n                        value = value[:50] + \"...\"\n                    value
  = value\n                msg = msg + f\"[grey46]{key}[/][magenta3]:[/] [grey66]{value}[/]\\n\"\n
  \           return msg\n    ```"
date: 2025-02-14
description: 'The `markata.plugins.feeds` plugin is used to create feed pages, which
  are lists of

  posts.  The list is generated using a `filter`, then each post in the list i'
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>feeds.py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"The `markata.plugins.feeds` plugin is used to create
    feed pages, which are lists of\nposts.  The list is generated using a `filter`,
    then each post in the list i\" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n        <meta property=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    <head>\n        <script src=\"https://unpkg.com/@tailwindcss/browser@4\"></script>\n<title>feeds.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"The `markata.plugins.feeds` plugin is
    used to create feed pages, which are lists of\nposts.  The list is generated using
    a `filter`, then each post in the list i\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        feeds.py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <hr />\n<p>The <code>markata.plugins.feeds</code>
    plugin is used to create feed pages, which are lists of\nposts.  The list is generated
    using a <code>filter</code>, then each post in the list is\nrendered with a <code>card_template</code>
    before being applied to the <code>body</code> of the\n<code>template</code>.</p>\n<h1
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>This plugin is built-in
    and enabled by default, but in you want to be very\nexplicit you can add it to
    your list of existing plugins.</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    class=\"p\">[</span>\n<span class=\"w\">   </span><span class=\"s2\">&quot;markata.plugins.feeds&quot;</span><span
    class=\"p\">,</span>\n<span class=\"w\">   </span><span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<h1
    id=\"configuration\">Configuration <a class=\"header-anchor\" href=\"#configuration\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<h1 id=\"set-default-template-and-card_template\">set
    default template and card_template <a class=\"header-anchor\" href=\"#set-default-template-and-card_template\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>At the root of the markata.feeds
    config you may set <code>template</code>, and\n<code>card_template</code>.  These
    will become your defaults for every feed you create.\nIf you do not set these,
    markata will use it's defaults.  The defaults are\ndesigned to work for a variety
    of use cases, but are not likely the best for\nall.</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata.feeds_config]</span>\n<span
    class=\"n\">template</span><span class=\"o\">=</span><span class=\"s2\">&quot;pages/templates/archive_template.html&quot;</span>\n<span
    class=\"n\">card_template</span><span class=\"o\">=</span><span class=\"s2\">&quot;plugins/feed_card_template.html&quot;</span>\n</pre></div>\n\n</pre>\n\n<h1
    id=\"pages\">pages <a class=\"header-anchor\" href=\"#pages\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>Underneath of the <code>markata.feeds</code>
    we will create a new map for each page where\nthe name of the map will be the
    name of the page.</p>\n<p>The following config will create a page at <code>/all-posts</code>
    that inclues every\nsingle post.</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.feeds]]</span>\n<span
    class=\"n\">title</span><span class=\"o\">=</span><span class=\"s2\">&quot;All
    Posts&quot;</span>\n<span class=\"n\">slug</span><span class=\"o\">=</span><span
    class=\"s1\">&#39;all&#39;</span>\n<span class=\"n\">filter</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;True&quot;</span>\n</pre></div>\n\n</pre>\n\n<h1 id=\"template\">template
    <a class=\"header-anchor\" href=\"#template\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>The <code>template</code>
    configuration key is a file path to the template that you want\nto use to create
    the feed.  You may set the default template you want to use\nfor all feeds under
    <code>[markata.feeds]</code>, as well as override it inside of each\nfeeds config.</p>\n<p>The
    template is a jinja style template that expects to fill in a <code>title</code>
    and\n<code>body</code> variable.</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"cp\">&lt;!DOCTYPE
    html&gt;</span>\n<span class=\"p\">&lt;</span><span class=\"nt\">html</span> <span
    class=\"na\">lang</span><span class=\"o\">=</span><span class=\"s\">&quot;en&quot;</span><span
    class=\"p\">&gt;</span>\n  <span class=\"p\">&lt;</span><span class=\"nt\">head</span><span
    class=\"p\">&gt;</span>\n    <span class=\"p\">&lt;</span><span class=\"nt\">title</span><span
    class=\"p\">&gt;</span>{{ title }}<span class=\"p\">&lt;/</span><span class=\"nt\">title</span><span
    class=\"p\">&gt;</span>\n  <span class=\"p\">&lt;/</span><span class=\"nt\">head</span><span
    class=\"p\">&gt;</span>\n  <span class=\"p\">&lt;</span><span class=\"nt\">body</span><span
    class=\"p\">&gt;</span>\n    <span class=\"p\">&lt;</span><span class=\"nt\">ul</span><span
    class=\"p\">&gt;</span>\n        {{ body }}\n    <span class=\"p\">&lt;/</span><span
    class=\"nt\">ul</span><span class=\"p\">&gt;</span>\n  <span class=\"p\">&lt;/</span><span
    class=\"nt\">body</span><span class=\"p\">&gt;</span>\n<span class=\"p\">&lt;/</span><span
    class=\"nt\">html</span><span class=\"p\">&gt;</span>\n</pre></div>\n\n</pre>\n\n<div
    class=\"admonition note\">\n<p class=\"admonition-title\">Note</p>\n<p>I highly
    reccomend putting your <code>body</code> in a <code>&lt;ul&gt;</code>, and wrapping
    your\n<code>card_template</code>s in an <code>&lt;li&gt;</code>.</p>\n</div>\n<h1
    id=\"card_template\">card_template <a class=\"header-anchor\" href=\"#card_template\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>All keys available from
    each post is available to put into your jinja\ntemplate.  These can either be
    placed there in your post frontmatter, or\nthrough a plugin that automatically
    adds to the post before the save phase.</p>\n<p>Here is a very simple example
    that would give a link to each post with the\ntitle and date.</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.feeds]]</span>\n<span
    class=\"n\">slug</span><span class=\"o\">=</span><span class=\"s1\">&#39;all&#39;</span>\n<span
    class=\"n\">title</span><span class=\"o\">=</span><span class=\"s1\">&#39;All
    Posts&#39;</span>\n<span class=\"n\">filter</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;True&quot;</span>\n<span class=\"n\">card_template</span><span
    class=\"o\">=</span><span class=\"s1\">&#39;&#39;&#39;</span>\n<span class=\"s1\">&lt;li&gt;</span>\n<span
    class=\"s1\">    &lt;a href={{markata.config.get(&#39;path_prefix&#39;, &#39;&#39;)}}{{slug}}&gt;</span>\n<span
    class=\"s1\">        {{title}}-{{date}}</span>\n<span class=\"s1\">    &lt;/a&gt;</span>\n<span
    class=\"s1\">&lt;/li&gt;</span>\n<span class=\"s1\">&#39;&#39;&#39;</span>\n</pre></div>\n\n</pre>\n\n<h1
    id=\"filter\">filter <a class=\"header-anchor\" href=\"#filter\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>The filter is a python
    expression ran on every post that expects to return a\nboolean.  The variables
    available to this expression are every key in your\nfrontmatter, plus the <code>timedelta</code>
    function, and <code>parse</code> function to more easily\nwork with dates.</p>\n<h1
    id=\"feed-examples\">Feed Examples <a class=\"header-anchor\" href=\"#feed-examples\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>True can be passed in
    to make a feed of all the posts you have.</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.feeds]]</span>\n<span
    class=\"n\">slug</span><span class=\"o\">=</span><span class=\"s1\">&#39;all&#39;</span>\n<span
    class=\"n\">title</span><span class=\"o\">=</span><span class=\"s1\">&#39;All
    Posts&#39;</span>\n<span class=\"n\">filter</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;True&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>You can compare
    against the values of the keys from your frontmatter.  This\nexample creates a
    feed that includes every post where published is <code>True</code>.</p>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.feeds]]</span>\n<span
    class=\"n\">slug</span><span class=\"o\">=</span><span class=\"s1\">&#39;draft&#39;</span>\n<span
    class=\"n\">title</span><span class=\"o\">=</span><span class=\"s1\">&#39;Draft&#39;</span>\n<span
    class=\"n\">filter</span><span class=\"o\">=</span><span class=\"s2\">&quot;published==&#39;False&#39;&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>We
    can also compare against dates.  The\n<a href=\"https://markata.dev/markata/plugins/datetime/\">markata.plugins.datetime</a>\nplugin,
    automatically adds <code>today</code> as today's date and <code>now</code> as
    the current\ndatetime.  These are quite handy to create feeds for scheduled, recent,
    or\ntoday's posts.  The following two examples will create a feed for scheduled\nposts
    and for today's posts respectively.</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.feeds]]</span>\n<span
    class=\"n\">slug</span><span class=\"o\">=</span><span class=\"s1\">&#39;scheduled&#39;</span>\n<span
    class=\"n\">title</span><span class=\"o\">=</span><span class=\"s1\">&#39;Scheduled&#39;</span>\n<span
    class=\"n\">filter</span><span class=\"o\">=</span><span class=\"s2\">&quot;date&gt;today&quot;</span>\n\n<span
    class=\"k\">[[markata.feeds]]</span>\n<span class=\"n\">slug</span><span class=\"o\">=</span><span
    class=\"s1\">&#39;today&#39;</span>\n<span class=\"n\">title</span><span class=\"o\">=</span><span
    class=\"s1\">&#39;Today&#39;</span>\n<span class=\"n\">filter</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;date==today&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>If you
    have list of items in your frontmatter for something like <code>tags</code>, you\ncan
    check for the existence of a tag in the list.</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.feeds]]</span>\n<span
    class=\"n\">slug</span><span class=\"o\">=</span><span class=\"s1\">&#39;python&#39;</span>\n<span
    class=\"n\">title</span><span class=\"o\">=</span><span class=\"s1\">&#39;Python&#39;</span>\n<span
    class=\"n\">filter</span><span class=\"o\">=</span><span class=\"s2\">&quot;date&lt;=today
    and &#39;python&#39; in tags&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>And of
    course you can combine all the things into larger expressions.  Here is\none example
    of the main feed on my blog.</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.feeds]]</span>\n<span
    class=\"n\">slug</span><span class=\"o\">=</span><span class=\"s1\">&#39;blog&#39;</span>\n<span
    class=\"n\">title</span><span class=\"o\">=</span><span class=\"s1\">&#39;Blog&#39;</span>\n<span
    class=\"n\">filter</span><span class=\"o\">=</span><span class=\"s2\">&quot;date&lt;=today
    and templateKey in [&#39;blog-post&#39;] and published ==&#39;True&#39;&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>Here
    is another example that shows my drafts for a particular tag.</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.feeds]]</span>\n<span
    class=\"n\">slug</span><span class=\"o\">=</span><span class=\"s1\">&#39;python-draft&#39;</span>\n<span
    class=\"n\">title</span><span class=\"o\">=</span><span class=\"s1\">&#39;Python
    Draft&#39;</span>\n<span class=\"n\">filter</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;date&lt;=today and &#39;python&#39; in tags and published==&#39;False&#39;&quot;</span>\n</pre></div>\n\n</pre>\n\n<h1
    id=\"defaults\">Defaults <a class=\"header-anchor\" href=\"#defaults\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>By default feeds will
    create one feed page at <code>/archive/</code> that includes all\nposts.</p>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.feeds]]</span>\n<span
    class=\"n\">slug</span><span class=\"o\">=</span><span class=\"s1\">&#39;archive&#39;</span>\n<span
    class=\"n\">title</span><span class=\"o\">=</span><span class=\"s1\">&#39;All
    Posts&#39;</span>\n<span class=\"n\">filter</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;True&quot;</span>\n</pre></div>\n\n</pre>\n\n<hr />\n<div class=\"admonition
    class\">\n<p class=\"admonition-title\">Class</p>\n<h2 id=\"Feed\" class=\"admonition-title\"
    style=\"margin: 0; padding: .5rem 1rem;\">Feed <em class=\"small\">class</em></h2>\n<p>A
    storage class for markata feed objects.</p>\n<h2 id=\"usage\">Usage <a class=\"header-anchor\"
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
    class=\"w\"> </span><span class=\"nn\">markata</span><span class=\"w\"> </span><span
    class=\"kn\">import</span> <span class=\"n\">Markata</span>\n<span class=\"n\">m</span>
    <span class=\"o\">=</span> <span class=\"n\">Markata</span><span class=\"p\">()</span>\n\n<span
    class=\"c1\"># access posts for a feed</span>\n<span class=\"n\">m</span><span
    class=\"o\">.</span><span class=\"n\">feeds</span><span class=\"o\">.</span><span
    class=\"n\">docs</span><span class=\"o\">.</span><span class=\"n\">posts</span>\n\n<span
    class=\"c1\"># access config for a feed</span>\n<span class=\"n\">m</span><span
    class=\"o\">.</span><span class=\"n\">feeds</span><span class=\"o\">.</span><span
    class=\"n\">docs</span><span class=\"o\">.</span><span class=\"n\">config</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">Feed
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">class</span><span
    class=\"w\"> </span><span class=\"nc\">Feed</span><span class=\"p\">(</span><span
    class=\"n\">pydantic</span><span class=\"o\">.</span><span class=\"n\">BaseModel</span><span
    class=\"p\">,</span> <span class=\"n\">JupyterMixin</span><span class=\"p\">):</span>\n<span
    class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;A storage class for
    markata feed objects.</span>\n\n<span class=\"sd\">    ## Usage</span>\n\n<span
    class=\"sd\">    ``` python</span>\n<span class=\"sd\">    from markata import
    Markata</span>\n<span class=\"sd\">    m = Markata()</span>\n\n<span class=\"sd\">
    \   # access posts for a feed</span>\n<span class=\"sd\">    m.feeds.docs.posts</span>\n\n<span
    class=\"sd\">    # access config for a feed</span>\n<span class=\"sd\">    m.feeds.docs.config</span>\n<span
    class=\"sd\">    ```</span>\n<span class=\"sd\">    &quot;&quot;&quot;</span>\n\n
    \   <span class=\"n\">config</span><span class=\"p\">:</span> <span class=\"n\">FeedConfig</span>\n
    \   <span class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"n\">Markata</span>
    <span class=\"o\">=</span> <span class=\"n\">pydantic</span><span class=\"o\">.</span><span
    class=\"n\">Field</span><span class=\"p\">(</span><span class=\"n\">exclude</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n\n
    \   <span class=\"n\">model_config</span> <span class=\"o\">=</span> <span class=\"n\">ConfigDict</span><span
    class=\"p\">(</span>\n        <span class=\"n\">validate_assignment</span><span
    class=\"o\">=</span><span class=\"kc\">False</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">arbitrary_types_allowed</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n        <span class=\"n\">extra</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;allow&quot;</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">str_strip_whitespace</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n        <span class=\"n\">validate_default</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">coerce_numbers_to_str</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n        <span class=\"n\">populate_by_name</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \   <span class=\"p\">)</span>\n\n    <span class=\"nd\">@property</span>\n    <span
    class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">name</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n<span
    class=\"w\">        </span><span class=\"sd\">&quot;&quot;&quot;The name of the
    feed, used for accessing it in the feeds object.&quot;&quot;&quot;</span>\n        <span
    class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">name</span>\n\n
    \   <span class=\"nd\">@property</span>\n    <span class=\"k\">def</span><span
    class=\"w\"> </span><span class=\"nf\">posts</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">):</span>\n        <span class=\"n\">posts</span>
    <span class=\"o\">=</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">map</span><span class=\"p\">(</span><span class=\"s2\">&quot;post&quot;</span><span
    class=\"p\">)</span>\n        <span class=\"k\">if</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">head</span> <span class=\"ow\">is</span> <span class=\"ow\">not</span>
    <span class=\"kc\">None</span> <span class=\"ow\">and</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">tail</span> <span class=\"ow\">is</span> <span class=\"ow\">not</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"n\">head_posts</span>
    <span class=\"o\">=</span> <span class=\"n\">posts</span><span class=\"p\">[:</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">head</span><span class=\"p\">]</span>\n
    \           <span class=\"n\">tail_posts</span> <span class=\"o\">=</span> <span
    class=\"n\">posts</span><span class=\"p\">[</span><span class=\"o\">-</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">tail</span> <span class=\"p\">:]</span>\n
    \           <span class=\"k\">return</span> <span class=\"n\">PrettyList</span><span
    class=\"p\">(</span><span class=\"n\">head_posts</span> <span class=\"o\">+</span>
    <span class=\"n\">tail_posts</span><span class=\"p\">)</span>\n        <span class=\"k\">if</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">head</span> <span class=\"ow\">is</span>
    <span class=\"ow\">not</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">return</span> <span class=\"n\">PrettyList</span><span
    class=\"p\">(</span><span class=\"n\">posts</span><span class=\"p\">[:</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">head</span><span class=\"p\">])</span>\n
    \       <span class=\"k\">if</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">tail</span>
    <span class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n            <span class=\"k\">return</span> <span class=\"n\">PrettyList</span><span
    class=\"p\">(</span><span class=\"n\">posts</span><span class=\"p\">[</span><span
    class=\"o\">-</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">tail</span>
    <span class=\"p\">:])</span>\n        <span class=\"k\">return</span> <span class=\"n\">PrettyList</span><span
    class=\"p\">(</span><span class=\"n\">posts</span><span class=\"p\">)</span>\n\n
    \   <span class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">first</span><span
    class=\"p\">(</span>\n        <span class=\"bp\">self</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span>\n    <span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">list</span><span
    class=\"p\">:</span>\n        <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">posts</span><span class=\"p\">[</span><span
    class=\"mi\">0</span><span class=\"p\">]</span>\n\n    <span class=\"k\">def</span><span
    class=\"w\"> </span><span class=\"nf\">last</span><span class=\"p\">(</span>\n
    \       <span class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span>\n    <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"nb\">list</span><span class=\"p\">:</span>\n        <span class=\"k\">return</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">posts</span><span
    class=\"p\">[</span><span class=\"o\">-</span><span class=\"mi\">1</span><span
    class=\"p\">]</span>\n\n    <span class=\"k\">def</span><span class=\"w\"> </span><span
    class=\"nf\">map</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">func</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;post&quot;</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"n\">args</span><span class=\"p\">):</span>\n        <span class=\"k\">return</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">map</span><span class=\"p\">(</span><span
    class=\"n\">func</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"p\">{</span><span class=\"o\">**</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">dict</span><span class=\"p\">(),</span> <span class=\"o\">**</span><span
    class=\"n\">args</span><span class=\"p\">})</span>\n\n    <span class=\"nd\">@property</span>\n
    \   <span class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">__rich_console__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"s2\">&quot;Console&quot;</span><span class=\"p\">:</span>\n
    \       <span class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span>\n\n
    \   <span class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">__rich__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Table</span><span class=\"p\">:</span>\n
    \       <span class=\"n\">table</span> <span class=\"o\">=</span> <span class=\"n\">Table</span><span
    class=\"p\">(</span><span class=\"n\">title</span><span class=\"o\">=</span><span
    class=\"sa\">f</span><span class=\"s2\">&quot;Feed: </span><span class=\"si\">{</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">name</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">)</span>\n\n
    \       <span class=\"n\">table</span><span class=\"o\">.</span><span class=\"n\">add_column</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">justify</span><span class=\"o\">=</span><span class=\"s2\">&quot;right&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">style</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;cyan&quot;</span><span class=\"p\">,</span> <span class=\"n\">no_wrap</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \       <span class=\"n\">table</span><span class=\"o\">.</span><span class=\"n\">add_column</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">justify</span><span class=\"o\">=</span><span class=\"s2\">&quot;left&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">style</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;green&quot;</span><span class=\"p\">)</span>\n        <span
    class=\"n\">table</span><span class=\"o\">.</span><span class=\"n\">add_column</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;published&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">justify</span><span class=\"o\">=</span><span class=\"s2\">&quot;left&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">style</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;green&quot;</span><span class=\"p\">)</span>\n\n        <span
    class=\"k\">for</span> <span class=\"n\">post</span> <span class=\"ow\">in</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">posts</span><span
    class=\"p\">:</span>\n            <span class=\"n\">table</span><span class=\"o\">.</span><span
    class=\"n\">add_row</span><span class=\"p\">(</span><span class=\"n\">post</span><span
    class=\"o\">.</span><span class=\"n\">title</span><span class=\"p\">,</span> <span
    class=\"n\">post</span><span class=\"o\">.</span><span class=\"n\">slug</span><span
    class=\"p\">,</span> <span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">post</span><span class=\"o\">.</span><span class=\"n\">published</span><span
    class=\"p\">))</span>\n\n        <span class=\"k\">return</span> <span class=\"n\">table</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition function\">\n<p class=\"admonition-title\">Function</p>\n<h2
    id=\"pre_render\" class=\"admonition-title\" style=\"margin: 0; padding: .5rem
    1rem;\">pre_render <em class=\"small\">function</em></h2>\n<p>Create the Feeds
    object and attach it to markata.</p>\n</div>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">pre_render <em class='small'>source</em></p>\n<pre
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
    class=\"w\"> </span><span class=\"nf\">pre_render</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"n\">Markata</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">    Create the Feeds object and attach it to markata.</span>\n<span
    class=\"sd\">    &quot;&quot;&quot;</span>\n    <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">feeds</span> <span class=\"o\">=</span>
    <span class=\"n\">Feeds</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div class=\"admonition
    function\">\n<p class=\"admonition-title\">Function</p>\n<h2 id=\"save\" class=\"admonition-title\"
    style=\"margin: 0; padding: .5rem 1rem;\">save <em class=\"small\">function</em></h2>\n<p>Creates
    a new feed page for each page in the config.</p>\n</div>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">save
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
    class=\"w\"> </span><span class=\"nf\">save</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"n\">Markata</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">    Creates a new feed page for each page in the config.</span>\n<span
    class=\"sd\">    &quot;&quot;&quot;</span>\n    <span class=\"k\">with</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">cache</span>
    <span class=\"k\">as</span> <span class=\"n\">cache</span><span class=\"p\">:</span>\n
    \       <span class=\"k\">for</span> <span class=\"n\">feed</span> <span class=\"ow\">in</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">feeds</span><span
    class=\"o\">.</span><span class=\"n\">values</span><span class=\"p\">():</span>\n
    \           <span class=\"n\">create_page</span><span class=\"p\">(</span>\n                <span
    class=\"n\">markata</span><span class=\"p\">,</span>\n                <span class=\"n\">feed</span><span
    class=\"p\">,</span>\n                <span class=\"n\">cache</span><span class=\"p\">,</span>\n
    \           <span class=\"p\">)</span>\n\n    <span class=\"n\">home</span> <span
    class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">output_dir</span><span class=\"p\">))</span> <span class=\"o\">/</span>
    <span class=\"s2\">&quot;index.html&quot;</span>\n    <span class=\"n\">archive</span>
    <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">output_dir</span><span class=\"p\">))</span> <span class=\"o\">/</span>
    <span class=\"s2\">&quot;archive&quot;</span> <span class=\"o\">/</span> <span
    class=\"s2\">&quot;index.html&quot;</span>\n    <span class=\"k\">if</span> <span
    class=\"ow\">not</span> <span class=\"n\">home</span><span class=\"o\">.</span><span
    class=\"n\">exists</span><span class=\"p\">()</span> <span class=\"ow\">and</span>
    <span class=\"n\">archive</span><span class=\"o\">.</span><span class=\"n\">exists</span><span
    class=\"p\">():</span>\n        <span class=\"n\">shutil</span><span class=\"o\">.</span><span
    class=\"n\">copy</span><span class=\"p\">(</span><span class=\"nb\">str</span><span
    class=\"p\">(</span><span class=\"n\">archive</span><span class=\"p\">),</span>
    <span class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">home</span><span
    class=\"p\">))</span>\n\n    <span class=\"n\">xsl_template</span> <span class=\"o\">=</span>
    <span class=\"n\">get_template</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"n\">feed</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">xsl_template</span><span
    class=\"p\">)</span>\n    <span class=\"n\">xsl</span> <span class=\"o\">=</span>
    <span class=\"n\">xsl_template</span><span class=\"o\">.</span><span class=\"n\">render</span><span
    class=\"p\">(</span>\n        <span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"p\">,</span>\n        <span class=\"n\">__version__</span><span
    class=\"o\">=</span><span class=\"n\">__version__</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">today</span><span class=\"o\">=</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">today</span><span class=\"p\">(),</span>\n        <span class=\"n\">config</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"p\">,</span>\n    <span class=\"p\">)</span>\n
    \   <span class=\"n\">xsl_file</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">output_dir</span><span
    class=\"p\">)</span> <span class=\"o\">/</span> <span class=\"s2\">&quot;rss.xsl&quot;</span>\n
    \   <span class=\"n\">xsl_file</span><span class=\"o\">.</span><span class=\"n\">write_text</span><span
    class=\"p\">(</span><span class=\"n\">xsl</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition function\">\n<p class=\"admonition-title\">Function</p>\n<h2
    id=\"create_page\" class=\"admonition-title\" style=\"margin: 0; padding: .5rem
    1rem;\">create_page <em class=\"small\">function</em></h2>\n<p>create an html
    unorderd list of posts.</p>\n</div>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">create_page <em class='small'>source</em></p>\n<pre
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
    class=\"w\"> </span><span class=\"nf\">create_page</span><span class=\"p\">(</span>\n
    \   <span class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"n\">Markata</span><span
    class=\"p\">,</span>\n    <span class=\"n\">feed</span><span class=\"p\">:</span>
    <span class=\"n\">Feed</span><span class=\"p\">,</span>\n    <span class=\"n\">cache</span><span
    class=\"p\">,</span>\n<span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span class=\"w\">    </span><span
    class=\"sd\">&quot;&quot;&quot;</span>\n<span class=\"sd\">    create an html
    unorderd list of posts.</span>\n<span class=\"sd\">    &quot;&quot;&quot;</span>\n\n
    \   <span class=\"n\">template</span> <span class=\"o\">=</span> <span class=\"n\">get_template</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">feed</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">template</span><span class=\"p\">)</span>\n
    \   <span class=\"n\">partial_template</span> <span class=\"o\">=</span> <span
    class=\"n\">get_template</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"n\">feed</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">partial_template</span><span
    class=\"p\">)</span>\n    <span class=\"n\">canonical_url</span> <span class=\"o\">=</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">url</span><span class=\"si\">}</span><span
    class=\"s2\">/</span><span class=\"si\">{</span><span class=\"n\">feed</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">slug</span><span class=\"si\">}</span><span class=\"s2\">/&quot;</span>\n\n
    \   <span class=\"n\">key</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">make_hash</span><span class=\"p\">(</span>\n
    \       <span class=\"s2\">&quot;feeds&quot;</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">template</span><span class=\"p\">,</span>\n        <span
    class=\"n\">__version__</span><span class=\"p\">,</span>\n        <span class=\"c1\">#
    cards,</span>\n        <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">url</span><span
    class=\"p\">,</span>\n        <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">description</span><span
    class=\"p\">,</span>\n        <span class=\"n\">feed</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">title</span><span
    class=\"p\">,</span>\n        <span class=\"n\">feed</span><span class=\"o\">.</span><span
    class=\"n\">map</span><span class=\"p\">(</span><span class=\"s2\">&quot;content&quot;</span><span
    class=\"p\">),</span>\n        <span class=\"n\">canonical_url</span><span class=\"p\">,</span>\n
    \       <span class=\"c1\"># datetime.datetime.today(),</span>\n        <span
    class=\"c1\"># markata.config,</span>\n    <span class=\"p\">)</span>\n\n    <span
    class=\"n\">html_key</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">make_hash</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">,</span> <span class=\"s2\">&quot;html&quot;</span><span
    class=\"p\">)</span>\n    <span class=\"n\">html_partial_key</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">make_hash</span><span
    class=\"p\">(</span><span class=\"n\">key</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;partial_html&quot;</span><span class=\"p\">)</span>\n    <span
    class=\"n\">feed_rss_key</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">make_hash</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">,</span> <span class=\"s2\">&quot;rss&quot;</span><span
    class=\"p\">)</span>\n    <span class=\"n\">feed_sitemap_key</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">make_hash</span><span
    class=\"p\">(</span><span class=\"n\">key</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;sitemap&quot;</span><span class=\"p\">)</span>\n\n    <span
    class=\"n\">feed_html_from_cache</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">precache</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"n\">html_key</span><span
    class=\"p\">)</span>\n    <span class=\"n\">feed_html_partial_from_cache</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">precache</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"n\">html_partial_key</span><span class=\"p\">)</span>\n
    \   <span class=\"n\">feed_rss_from_cache</span> <span class=\"o\">=</span> <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">precache</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"n\">feed_rss_key</span><span class=\"p\">)</span>\n    <span class=\"n\">feed_sitemap_from_cache</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">precache</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"n\">feed_sitemap_key</span><span class=\"p\">)</span>\n\n
    \   <span class=\"n\">output_file</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">output_dir</span><span
    class=\"p\">)</span> <span class=\"o\">/</span> <span class=\"n\">feed</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">slug</span> <span class=\"o\">/</span> <span class=\"s2\">&quot;index.html&quot;</span>\n
    \   <span class=\"n\">output_file</span><span class=\"o\">.</span><span class=\"n\">parent</span><span
    class=\"o\">.</span><span class=\"n\">mkdir</span><span class=\"p\">(</span><span
    class=\"n\">exist_ok</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span> <span class=\"n\">parents</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n\n    <span class=\"n\">partial_output_file</span>
    <span class=\"o\">=</span> <span class=\"p\">(</span>\n        <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">output_dir</span><span
    class=\"p\">)</span> <span class=\"o\">/</span> <span class=\"n\">feed</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">slug</span> <span class=\"o\">/</span> <span class=\"s2\">&quot;partial&quot;</span>
    <span class=\"o\">/</span> <span class=\"s2\">&quot;index.html&quot;</span>\n
    \   <span class=\"p\">)</span>\n    <span class=\"n\">partial_output_file</span><span
    class=\"o\">.</span><span class=\"n\">parent</span><span class=\"o\">.</span><span
    class=\"n\">mkdir</span><span class=\"p\">(</span><span class=\"n\">exist_ok</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
    class=\"n\">parents</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n\n    <span class=\"n\">rss_output_file</span> <span class=\"o\">=</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">output_dir</span><span class=\"p\">)</span> <span class=\"o\">/</span>
    <span class=\"n\">feed</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">slug</span> <span class=\"o\">/</span> <span
    class=\"s2\">&quot;rss.xml&quot;</span>\n    <span class=\"n\">rss_output_file</span><span
    class=\"o\">.</span><span class=\"n\">parent</span><span class=\"o\">.</span><span
    class=\"n\">mkdir</span><span class=\"p\">(</span><span class=\"n\">exist_ok</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
    class=\"n\">parents</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n\n    <span class=\"n\">sitemap_output_file</span> <span
    class=\"o\">=</span> <span class=\"p\">(</span>\n        <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">output_dir</span><span
    class=\"p\">)</span> <span class=\"o\">/</span> <span class=\"n\">feed</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">slug</span> <span class=\"o\">/</span> <span class=\"s2\">&quot;sitemap.xml&quot;</span>\n
    \   <span class=\"p\">)</span>\n    <span class=\"n\">sitemap_output_file</span><span
    class=\"o\">.</span><span class=\"n\">parent</span><span class=\"o\">.</span><span
    class=\"n\">mkdir</span><span class=\"p\">(</span><span class=\"n\">exist_ok</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
    class=\"n\">parents</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n\n    <span class=\"n\">from_cache</span> <span class=\"o\">=</span>
    <span class=\"kc\">True</span>\n    <span class=\"k\">if</span> <span class=\"n\">feed_html_from_cache</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \       <span class=\"n\">from_cache</span> <span class=\"o\">=</span> <span class=\"kc\">False</span>\n
    \       <span class=\"n\">feed_html</span> <span class=\"o\">=</span> <span class=\"n\">template</span><span
    class=\"o\">.</span><span class=\"n\">render</span><span class=\"p\">(</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">=</span><span class=\"n\">markata</span><span
    class=\"p\">,</span>\n            <span class=\"n\">__version__</span><span class=\"o\">=</span><span
    class=\"n\">__version__</span><span class=\"p\">,</span>\n            <span class=\"n\">post</span><span
    class=\"o\">=</span><span class=\"n\">feed</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">model_dump</span><span
    class=\"p\">(),</span>\n            <span class=\"n\">url</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">url</span><span class=\"p\">,</span>\n            <span
    class=\"n\">config</span><span class=\"o\">=</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">,</span>\n
    \           <span class=\"n\">feed</span><span class=\"o\">=</span><span class=\"n\">feed</span><span
    class=\"p\">,</span>\n        <span class=\"p\">)</span>\n        <span class=\"n\">cache</span><span
    class=\"o\">.</span><span class=\"n\">set</span><span class=\"p\">(</span><span
    class=\"n\">html_key</span><span class=\"p\">,</span> <span class=\"n\">feed_html</span><span
    class=\"p\">)</span>\n    <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \       <span class=\"n\">feed_html</span> <span class=\"o\">=</span> <span class=\"n\">feed_html_from_cache</span>\n\n
    \   <span class=\"k\">if</span> <span class=\"n\">feed_html_partial_from_cache</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \       <span class=\"n\">from_cache</span> <span class=\"o\">=</span> <span class=\"kc\">False</span>\n
    \       <span class=\"n\">feed_html_partial</span> <span class=\"o\">=</span>
    <span class=\"n\">partial_template</span><span class=\"o\">.</span><span class=\"n\">render</span><span
    class=\"p\">(</span>\n            <span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"p\">,</span>\n            <span class=\"n\">__version__</span><span
    class=\"o\">=</span><span class=\"n\">__version__</span><span class=\"p\">,</span>\n
    \           <span class=\"n\">post</span><span class=\"o\">=</span><span class=\"n\">feed</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">model_dump</span><span class=\"p\">(),</span>\n            <span class=\"n\">url</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">url</span><span
    class=\"p\">,</span>\n            <span class=\"n\">config</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"p\">,</span>\n            <span class=\"n\">feed</span><span class=\"o\">=</span><span
    class=\"n\">feed</span><span class=\"p\">,</span>\n        <span class=\"p\">)</span>\n
    \       <span class=\"n\">cache</span><span class=\"o\">.</span><span class=\"n\">set</span><span
    class=\"p\">(</span><span class=\"n\">html_partial_key</span><span class=\"p\">,</span>
    <span class=\"n\">feed_html_partial</span><span class=\"p\">)</span>\n    <span
    class=\"k\">else</span><span class=\"p\">:</span>\n        <span class=\"n\">feed_html_partial</span>
    <span class=\"o\">=</span> <span class=\"n\">feed_html_partial_from_cache</span>\n\n
    \   <span class=\"k\">if</span> <span class=\"n\">feed_rss_from_cache</span> <span
    class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \       <span class=\"n\">from_cache</span> <span class=\"o\">=</span> <span class=\"kc\">False</span>\n
    \       <span class=\"n\">rss_template</span> <span class=\"o\">=</span> <span
    class=\"n\">get_template</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"n\">feed</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">rss_template</span><span
    class=\"p\">)</span>\n        <span class=\"n\">feed_rss</span> <span class=\"o\">=</span>
    <span class=\"n\">rss_template</span><span class=\"o\">.</span><span class=\"n\">render</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"p\">,</span> <span class=\"n\">feed</span><span
    class=\"o\">=</span><span class=\"n\">feed</span><span class=\"p\">)</span>\n
    \       <span class=\"n\">cache</span><span class=\"o\">.</span><span class=\"n\">set</span><span
    class=\"p\">(</span><span class=\"n\">feed_rss_key</span><span class=\"p\">,</span>
    <span class=\"n\">feed_rss</span><span class=\"p\">)</span>\n    <span class=\"k\">else</span><span
    class=\"p\">:</span>\n        <span class=\"n\">feed_rss</span> <span class=\"o\">=</span>
    <span class=\"n\">feed_rss_from_cache</span>\n\n    <span class=\"k\">if</span>
    <span class=\"n\">feed_sitemap_from_cache</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n        <span class=\"n\">from_cache</span>
    <span class=\"o\">=</span> <span class=\"kc\">False</span>\n        <span class=\"n\">sitemap_template</span>
    <span class=\"o\">=</span> <span class=\"n\">get_template</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">,</span> <span class=\"n\">feed</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">sitemap_template</span><span class=\"p\">)</span>\n        <span class=\"n\">feed_sitemap</span>
    <span class=\"o\">=</span> <span class=\"n\">sitemap_template</span><span class=\"o\">.</span><span
    class=\"n\">render</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">feed</span><span class=\"o\">=</span><span class=\"n\">feed</span><span
    class=\"p\">)</span>\n        <span class=\"n\">cache</span><span class=\"o\">.</span><span
    class=\"n\">set</span><span class=\"p\">(</span><span class=\"n\">feed_sitemap_key</span><span
    class=\"p\">,</span> <span class=\"n\">feed_sitemap</span><span class=\"p\">)</span>\n
    \   <span class=\"k\">else</span><span class=\"p\">:</span>\n        <span class=\"n\">feed_sitemap</span>
    <span class=\"o\">=</span> <span class=\"n\">feed_sitemap_from_cache</span>\n\n
    \   <span class=\"k\">if</span> <span class=\"p\">(</span>\n        <span class=\"n\">from_cache</span>\n
    \       <span class=\"ow\">and</span> <span class=\"n\">output_file</span><span
    class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">()</span>\n
    \       <span class=\"ow\">and</span> <span class=\"n\">partial_output_file</span><span
    class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">()</span>\n
    \       <span class=\"ow\">and</span> <span class=\"n\">rss_output_file</span><span
    class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">()</span>\n
    \       <span class=\"ow\">and</span> <span class=\"n\">sitemap_output_file</span><span
    class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">()</span>\n
    \   <span class=\"p\">):</span>\n        <span class=\"k\">return</span>\n\n    <span
    class=\"n\">output_file</span><span class=\"o\">.</span><span class=\"n\">write_text</span><span
    class=\"p\">(</span><span class=\"n\">feed_html</span><span class=\"p\">)</span>\n
    \   <span class=\"n\">partial_output_file</span><span class=\"o\">.</span><span
    class=\"n\">write_text</span><span class=\"p\">(</span><span class=\"n\">feed_html_partial</span><span
    class=\"p\">)</span>\n    <span class=\"n\">rss_output_file</span><span class=\"o\">.</span><span
    class=\"n\">write_text</span><span class=\"p\">(</span><span class=\"n\">feed_rss</span><span
    class=\"p\">)</span>\n    <span class=\"n\">sitemap_output_file</span><span class=\"o\">.</span><span
    class=\"n\">write_text</span><span class=\"p\">(</span><span class=\"n\">feed_sitemap</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div class=\"admonition
    function\">\n<p class=\"admonition-title\">Function</p>\n<h2 id=\"create_card\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">create_card
    <em class=\"small\">function</em></h2>\n<p>Creates a card for one post based on
    the configured template.  If no\ntemplate is configured it will create one with
    the post title and dates\n(if present).</p>\n</div>\n<div class=\"admonition source
    is-collapsible collapsible-open\">\n<p class=\"admonition-title\">create_card
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
    class=\"w\"> </span><span class=\"nf\">create_card</span><span class=\"p\">(</span>\n
    \   <span class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span>\n    <span class=\"n\">post</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">,</span>\n    <span
    class=\"n\">template</span><span class=\"p\">:</span> <span class=\"n\">Optional</span><span
    class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">]</span> <span
    class=\"o\">=</span> <span class=\"kc\">None</span><span class=\"p\">,</span>\n
    \   <span class=\"n\">cache</span><span class=\"o\">=</span><span class=\"kc\">None</span><span
    class=\"p\">,</span>\n<span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Any</span><span class=\"p\">:</span>\n<span class=\"w\">    </span><span
    class=\"sd\">&quot;&quot;&quot;</span>\n<span class=\"sd\">    Creates a card
    for one post based on the configured template.  If no</span>\n<span class=\"sd\">
    \   template is configured it will create one with the post title and dates</span>\n<span
    class=\"sd\">    (if present).</span>\n<span class=\"sd\">    &quot;&quot;&quot;</span>\n
    \   <span class=\"n\">key</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">make_hash</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;feeds&quot;</span><span class=\"p\">,</span> <span class=\"n\">template</span><span
    class=\"p\">,</span> <span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">post</span><span class=\"p\">),</span> <span class=\"n\">post</span><span
    class=\"o\">.</span><span class=\"n\">content</span><span class=\"p\">)</span>\n\n
    \   <span class=\"n\">card</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">precache</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"n\">key</span><span
    class=\"p\">)</span>\n    <span class=\"k\">if</span> <span class=\"n\">card</span>
    <span class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n        <span class=\"k\">return</span> <span class=\"n\">card</span>\n\n
    \   <span class=\"k\">if</span> <span class=\"n\">template</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n        <span class=\"n\">template</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;feeds_config&quot;</span><span class=\"p\">,</span>
    <span class=\"p\">{})</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;card_template&quot;</span><span class=\"p\">,</span>
    <span class=\"kc\">None</span><span class=\"p\">)</span>\n\n    <span class=\"k\">if</span>
    <span class=\"n\">template</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n        <span class=\"k\">if</span> <span class=\"s2\">&quot;date&quot;</span>
    <span class=\"ow\">in</span> <span class=\"n\">post</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">card</span> <span class=\"o\">=</span> <span class=\"n\">textwrap</span><span
    class=\"o\">.</span><span class=\"n\">dedent</span><span class=\"p\">(</span>\n
    \               <span class=\"sa\">f</span><span class=\"s2\">&quot;&quot;&quot;</span>\n<span
    class=\"s2\">                &lt;li class=&#39;post&#39;&gt;</span>\n<span class=\"s2\">
    \               &lt;a href=&quot;/</span><span class=\"si\">{</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">path_prefix</span><span class=\"si\">}{</span><span class=\"n\">post</span><span
    class=\"o\">.</span><span class=\"n\">slug</span><span class=\"si\">}</span><span
    class=\"s2\">/&quot;&gt;</span>\n<span class=\"s2\">                    </span><span
    class=\"si\">{</span><span class=\"n\">post</span><span class=\"o\">.</span><span
    class=\"n\">title</span><span class=\"si\">}</span>\n<span class=\"s2\">                    </span><span
    class=\"si\">{</span><span class=\"n\">post</span><span class=\"o\">.</span><span
    class=\"n\">date</span><span class=\"o\">.</span><span class=\"n\">year</span><span
    class=\"si\">}</span><span class=\"s2\">-</span>\n<span class=\"s2\">                    </span><span
    class=\"si\">{</span><span class=\"n\">post</span><span class=\"o\">.</span><span
    class=\"n\">date</span><span class=\"o\">.</span><span class=\"n\">month</span><span
    class=\"si\">}</span><span class=\"s2\">-</span>\n<span class=\"s2\">                    </span><span
    class=\"si\">{</span><span class=\"n\">post</span><span class=\"o\">.</span><span
    class=\"n\">date</span><span class=\"o\">.</span><span class=\"n\">day</span><span
    class=\"si\">}</span>\n<span class=\"s2\">                &lt;/a&gt;</span>\n<span
    class=\"s2\">                &lt;/li&gt;</span>\n<span class=\"s2\">                &quot;&quot;&quot;</span><span
    class=\"p\">,</span>\n            <span class=\"p\">)</span>\n        <span class=\"k\">else</span><span
    class=\"p\">:</span>\n            <span class=\"n\">card</span> <span class=\"o\">=</span>
    <span class=\"n\">textwrap</span><span class=\"o\">.</span><span class=\"n\">dedent</span><span
    class=\"p\">(</span>\n                <span class=\"sa\">f</span><span class=\"s2\">&quot;&quot;&quot;</span>\n<span
    class=\"s2\">                &lt;li class=&#39;post&#39;&gt;</span>\n<span class=\"s2\">
    \               &lt;a href=&quot;/</span><span class=\"si\">{</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">path_prefix</span><span class=\"si\">}{</span><span class=\"n\">post</span><span
    class=\"o\">.</span><span class=\"n\">slug</span><span class=\"si\">}</span><span
    class=\"s2\">/&quot;&gt;</span>\n<span class=\"s2\">                    </span><span
    class=\"si\">{</span><span class=\"n\">post</span><span class=\"o\">.</span><span
    class=\"n\">title</span><span class=\"si\">}</span>\n<span class=\"s2\">                &lt;/a&gt;</span>\n<span
    class=\"s2\">                &lt;/li&gt;</span>\n<span class=\"s2\">                &quot;&quot;&quot;</span><span
    class=\"p\">,</span>\n            <span class=\"p\">)</span>\n    <span class=\"k\">else</span><span
    class=\"p\">:</span>\n        <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">_template</span> <span class=\"o\">=</span> <span
    class=\"n\">Template</span><span class=\"p\">(</span><span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">template</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">read_text</span><span class=\"p\">())</span>\n
    \       <span class=\"k\">except</span> <span class=\"ne\">FileNotFoundError</span><span
    class=\"p\">:</span>\n            <span class=\"n\">_template</span> <span class=\"o\">=</span>
    <span class=\"n\">Template</span><span class=\"p\">(</span><span class=\"n\">template</span><span
    class=\"p\">)</span>\n        <span class=\"k\">except</span> <span class=\"ne\">OSError</span><span
    class=\"p\">:</span>  <span class=\"c1\"># File name too long</span>\n            <span
    class=\"n\">_template</span> <span class=\"o\">=</span> <span class=\"n\">Template</span><span
    class=\"p\">(</span><span class=\"n\">template</span><span class=\"p\">)</span>\n
    \       <span class=\"n\">card</span> <span class=\"o\">=</span> <span class=\"n\">_template</span><span
    class=\"o\">.</span><span class=\"n\">render</span><span class=\"p\">(</span><span
    class=\"n\">post</span><span class=\"o\">=</span><span class=\"n\">post</span><span
    class=\"p\">,</span> <span class=\"o\">**</span><span class=\"n\">post</span><span
    class=\"o\">.</span><span class=\"n\">to_dict</span><span class=\"p\">())</span>\n
    \   <span class=\"n\">cache</span><span class=\"o\">.</span><span class=\"n\">add</span><span
    class=\"p\">(</span><span class=\"n\">key</span><span class=\"p\">,</span> <span
    class=\"n\">card</span><span class=\"p\">)</span>\n    <span class=\"k\">return</span>
    <span class=\"n\">card</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div class=\"admonition
    class\">\n<p class=\"admonition-title\">Class</p>\n<h2 id=\"Feeds\" class=\"admonition-title\"
    style=\"margin: 0; padding: .5rem 1rem;\">Feeds <em class=\"small\">class</em></h2>\n<p>A
    storage class for all markata Feed objects</p>\n<pre class='wrapper'>\n\n<div
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
    class=\"kn\">import</span> <span class=\"n\">Markata</span>\n<span class=\"n\">markata</span>
    <span class=\"o\">=</span> <span class=\"n\">Markata</span><span class=\"p\">()</span>\n\n<span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">feeds</span>\n\n<span
    class=\"c1\"># access all config</span>\n<span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">feeds</span><span class=\"o\">.</span><span
    class=\"n\">config</span>\n\n<span class=\"c1\"># refresh list of posts in all
    feeds</span>\n<span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">feeds</span><span class=\"o\">.</span><span class=\"n\">refresh</span><span
    class=\"p\">()</span>\n\n<span class=\"c1\"># iterating over feeds gives the name
    of the feed</span>\n<span class=\"k\">for</span> <span class=\"n\">k</span> <span
    class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">feeds</span><span class=\"p\">:</span>\n     <span class=\"nb\">print</span><span
    class=\"p\">(</span><span class=\"n\">k</span><span class=\"p\">)</span>\n\n<span
    class=\"c1\"># project-gallery</span>\n<span class=\"c1\"># docs</span>\n<span
    class=\"c1\"># autodoc</span>\n<span class=\"c1\"># core_modules</span>\n<span
    class=\"c1\"># plugins</span>\n<span class=\"c1\"># archive</span>\n\n<span class=\"c1\">#
    iterate over items like keys and values in a dict, items returns name of</span>\n<span
    class=\"c1\"># feed and a feed object</span>\n<span class=\"k\">for</span> <span
    class=\"n\">k</span><span class=\"p\">,</span> <span class=\"n\">v</span> <span
    class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">feeds</span><span class=\"o\">.</span><span class=\"n\">items</span><span
    class=\"p\">():</span>\n    <span class=\"nb\">print</span><span class=\"p\">(</span><span
    class=\"n\">k</span><span class=\"p\">,</span> <span class=\"nb\">len</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"o\">.</span><span
    class=\"n\">posts</span><span class=\"p\">))</span>\n\n<span class=\"c1\"># project-gallery
    2</span>\n<span class=\"c1\"># docs 6</span>\n<span class=\"c1\"># autodoc 65</span>\n<span
    class=\"c1\"># core_modules 26</span>\n<span class=\"c1\"># plugins 39</span>\n<span
    class=\"c1\"># archive 65</span>\n\n<span class=\"c1\"># values can be iterated
    over in just the same way</span>\n<span class=\"k\">for</span> <span class=\"n\">v</span>
    <span class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">feeds</span><span class=\"o\">.</span><span class=\"n\">values</span><span
    class=\"p\">():</span>\n     <span class=\"nb\">print</span><span class=\"p\">(</span><span
    class=\"nb\">len</span><span class=\"p\">(</span><span class=\"n\">v</span><span
    class=\"o\">.</span><span class=\"n\">posts</span><span class=\"p\">))</span>\n<span
    class=\"c1\"># 2</span>\n<span class=\"c1\"># 6</span>\n<span class=\"c1\"># 65</span>\n<span
    class=\"c1\"># 26</span>\n<span class=\"c1\"># 39</span>\n<span class=\"c1\">#
    65</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">Feeds <em class='small'>source</em></p>\n<pre
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
    class=\"w\"> </span><span class=\"nc\">Feeds</span><span class=\"p\">(</span><span
    class=\"n\">JupyterMixin</span><span class=\"p\">):</span>\n<span class=\"w\">
    \   </span><span class=\"sd\">&quot;&quot;&quot;A storage class for all markata
    Feed objects</span>\n\n<span class=\"sd\">    ``` python</span>\n<span class=\"sd\">
    \   from markata import Markata</span>\n<span class=\"sd\">    markata = Markata()</span>\n\n<span
    class=\"sd\">    markata.feeds</span>\n\n<span class=\"sd\">    # access all config</span>\n<span
    class=\"sd\">    markata.feeds.config</span>\n\n<span class=\"sd\">    # refresh
    list of posts in all feeds</span>\n<span class=\"sd\">    markata.feeds.refresh()</span>\n\n<span
    class=\"sd\">    # iterating over feeds gives the name of the feed</span>\n<span
    class=\"sd\">    for k in markata.feeds:</span>\n<span class=\"sd\">         print(k)</span>\n\n<span
    class=\"sd\">    # project-gallery</span>\n<span class=\"sd\">    # docs</span>\n<span
    class=\"sd\">    # autodoc</span>\n<span class=\"sd\">    # core_modules</span>\n<span
    class=\"sd\">    # plugins</span>\n<span class=\"sd\">    # archive</span>\n\n<span
    class=\"sd\">    # iterate over items like keys and values in a dict, items returns
    name of</span>\n<span class=\"sd\">    # feed and a feed object</span>\n<span
    class=\"sd\">    for k, v in markata.feeds.items():</span>\n<span class=\"sd\">
    \       print(k, len(v.posts))</span>\n\n<span class=\"sd\">    # project-gallery
    2</span>\n<span class=\"sd\">    # docs 6</span>\n<span class=\"sd\">    # autodoc
    65</span>\n<span class=\"sd\">    # core_modules 26</span>\n<span class=\"sd\">
    \   # plugins 39</span>\n<span class=\"sd\">    # archive 65</span>\n\n<span class=\"sd\">
    \   # values can be iterated over in just the same way</span>\n<span class=\"sd\">
    \   for v in markata.feeds.values():</span>\n<span class=\"sd\">         print(len(v.posts))</span>\n<span
    class=\"sd\">    # 2</span>\n<span class=\"sd\">    # 6</span>\n<span class=\"sd\">
    \   # 65</span>\n<span class=\"sd\">    # 26</span>\n<span class=\"sd\">    #
    39</span>\n<span class=\"sd\">    # 65</span>\n<span class=\"sd\">    &quot;&quot;&quot;</span>\n\n
    \   <span class=\"k\">def</span><span class=\"w\"> </span><span class=\"fm\">__init__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"n\">Markata</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n        <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">markata</span> <span class=\"o\">=</span> <span class=\"n\">markata</span>\n
    \       <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span>
    <span class=\"o\">=</span> <span class=\"p\">{</span><span class=\"n\">f</span><span
    class=\"o\">.</span><span class=\"n\">name</span><span class=\"p\">:</span> <span
    class=\"n\">f</span> <span class=\"k\">for</span> <span class=\"n\">f</span> <span
    class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">feeds</span><span
    class=\"p\">}</span>\n        <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">refresh</span><span class=\"p\">()</span>\n\n    <span class=\"k\">def</span><span
    class=\"w\"> </span><span class=\"nf\">refresh</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">):</span>\n<span class=\"w\">        </span><span
    class=\"sd\">&quot;&quot;&quot;Refresh all of the feeds objects&quot;&quot;&quot;</span>\n
    \       <span class=\"k\">for</span> <span class=\"n\">feed_config</span> <span
    class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">feeds</span><span class=\"p\">:</span>\n
    \           <span class=\"c1\"># Ensure feed has a name, falling back to slug
    if needed</span>\n            <span class=\"k\">if</span> <span class=\"n\">feed_config</span><span
    class=\"o\">.</span><span class=\"n\">name</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span> <span class=\"ow\">and</span> <span class=\"n\">feed_config</span><span
    class=\"o\">.</span><span class=\"n\">slug</span> <span class=\"ow\">is</span>
    <span class=\"ow\">not</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">feed_config</span><span class=\"o\">.</span><span
    class=\"n\">name</span> <span class=\"o\">=</span> <span class=\"n\">feed_config</span><span
    class=\"o\">.</span><span class=\"n\">slug</span><span class=\"o\">.</span><span
    class=\"n\">replace</span><span class=\"p\">(</span><span class=\"s2\">&quot;-&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;_&quot;</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">elif</span> <span class=\"n\">feed_config</span><span
    class=\"o\">.</span><span class=\"n\">name</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span> <span class=\"ow\">and</span> <span class=\"n\">feed_config</span><span
    class=\"o\">.</span><span class=\"n\">slug</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
    class=\"n\">feed_config</span><span class=\"o\">.</span><span class=\"n\">slug</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;archive&quot;</span>\n                <span
    class=\"n\">feed_config</span><span class=\"o\">.</span><span class=\"n\">name</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;archive&quot;</span>\n\n            <span
    class=\"n\">feed</span> <span class=\"o\">=</span> <span class=\"n\">Feed</span><span
    class=\"p\">(</span><span class=\"n\">config</span><span class=\"o\">=</span><span
    class=\"n\">feed_config</span><span class=\"p\">,</span> <span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">markata</span><span class=\"p\">)</span>\n            <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"fm\">__setattr__</span><span class=\"p\">(</span><span
    class=\"n\">feed</span><span class=\"o\">.</span><span class=\"n\">name</span><span
    class=\"p\">,</span> <span class=\"n\">feed</span><span class=\"p\">)</span>\n\n
    \   <span class=\"k\">def</span><span class=\"w\"> </span><span class=\"fm\">__iter__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">):</span>\n
    \       <span class=\"k\">return</span> <span class=\"nb\">iter</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">keys</span><span class=\"p\">())</span>\n\n
    \   <span class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">keys</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">):</span>\n
    \       <span class=\"k\">return</span> <span class=\"nb\">iter</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">keys</span><span class=\"p\">())</span>\n\n
    \   <span class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">values</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">):</span>\n
    \       <span class=\"k\">return</span> <span class=\"p\">[</span><span class=\"bp\">self</span><span
    class=\"p\">[</span><span class=\"n\">feed</span><span class=\"p\">]</span> <span
    class=\"k\">for</span> <span class=\"n\">feed</span> <span class=\"ow\">in</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">keys</span><span class=\"p\">()]</span>\n\n
    \   <span class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">items</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">):</span>\n
    \       <span class=\"k\">return</span> <span class=\"p\">[(</span><span class=\"n\">key</span><span
    class=\"p\">,</span> <span class=\"bp\">self</span><span class=\"p\">[</span><span
    class=\"n\">key</span><span class=\"p\">])</span> <span class=\"k\">for</span>
    <span class=\"n\">key</span> <span class=\"ow\">in</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">]</span>\n\n
    \   <span class=\"k\">def</span><span class=\"w\"> </span><span class=\"fm\">__getitem__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"n\">key</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span
    class=\"p\">:</span>\n        <span class=\"k\">return</span> <span class=\"nb\">getattr</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"n\">key</span><span class=\"o\">.</span><span class=\"n\">replace</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;-&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;_&quot;</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">lower</span><span class=\"p\">())</span>\n\n    <span class=\"k\">def</span><span
    class=\"w\"> </span><span class=\"nf\">get</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"n\">key</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">,</span> <span
    class=\"n\">default</span><span class=\"p\">:</span> <span class=\"n\">Any</span>
    <span class=\"o\">=</span> <span class=\"kc\">None</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span class=\"p\">:</span>\n
    \       <span class=\"k\">return</span> <span class=\"nb\">getattr</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"n\">key</span><span class=\"o\">.</span><span class=\"n\">replace</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;-&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;_&quot;</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">lower</span><span class=\"p\">(),</span> <span class=\"n\">default</span><span
    class=\"p\">)</span>\n\n    <span class=\"k\">def</span><span class=\"w\"> </span><span
    class=\"nf\">_dict_panel</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">config</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n<span
    class=\"w\">        </span><span class=\"sd\">&quot;&quot;&quot;pretty print configs
    with rich&quot;&quot;&quot;</span>\n        <span class=\"n\">msg</span> <span
    class=\"o\">=</span> <span class=\"s2\">&quot;&quot;</span>\n        <span class=\"k\">for</span>
    <span class=\"n\">key</span><span class=\"p\">,</span> <span class=\"n\">value</span>
    <span class=\"ow\">in</span> <span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">items</span><span class=\"p\">():</span>\n            <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">value</span><span
    class=\"p\">,</span> <span class=\"nb\">str</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">if</span> <span class=\"nb\">len</span><span
    class=\"p\">(</span><span class=\"n\">value</span><span class=\"p\">)</span> <span
    class=\"o\">&gt;</span> <span class=\"mi\">50</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">value</span> <span class=\"o\">=</span>
    <span class=\"n\">value</span><span class=\"p\">[:</span><span class=\"mi\">50</span><span
    class=\"p\">]</span> <span class=\"o\">+</span> <span class=\"s2\">&quot;...&quot;</span>\n
    \               <span class=\"n\">value</span> <span class=\"o\">=</span> <span
    class=\"n\">value</span>\n            <span class=\"n\">msg</span> <span class=\"o\">=</span>
    <span class=\"n\">msg</span> <span class=\"o\">+</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;[grey46]</span><span class=\"si\">{</span><span class=\"n\">key</span><span
    class=\"si\">}</span><span class=\"s2\">[/][magenta3]:[/] [grey66]</span><span
    class=\"si\">{</span><span class=\"n\">value</span><span class=\"si\">}</span><span
    class=\"s2\">[/]</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span>\n
    \       <span class=\"k\">return</span> <span class=\"n\">msg</span>\n\n    <span
    class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">__rich__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Table</span><span class=\"p\">:</span>\n
    \       <span class=\"n\">table</span> <span class=\"o\">=</span> <span class=\"n\">Table</span><span
    class=\"p\">(</span><span class=\"n\">title</span><span class=\"o\">=</span><span
    class=\"sa\">f</span><span class=\"s2\">&quot;Feeds </span><span class=\"si\">{</span><span
    class=\"nb\">len</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">)</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">)</span>\n\n
    \       <span class=\"n\">table</span><span class=\"o\">.</span><span class=\"n\">add_column</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;Feed&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">justify</span><span class=\"o\">=</span><span class=\"s2\">&quot;right&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">style</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;cyan&quot;</span><span class=\"p\">,</span> <span class=\"n\">no_wrap</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \       <span class=\"n\">table</span><span class=\"o\">.</span><span class=\"n\">add_column</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;posts&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">justify</span><span class=\"o\">=</span><span class=\"s2\">&quot;left&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">style</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;green&quot;</span><span class=\"p\">)</span>\n        <span
    class=\"n\">table</span><span class=\"o\">.</span><span class=\"n\">add_column</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;config&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">style</span><span class=\"o\">=</span><span class=\"s2\">&quot;magenta&quot;</span><span
    class=\"p\">)</span>\n\n        <span class=\"k\">for</span> <span class=\"n\">name</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"p\">:</span>\n            <span class=\"n\">table</span><span
    class=\"o\">.</span><span class=\"n\">add_row</span><span class=\"p\">(</span>\n
    \               <span class=\"n\">name</span><span class=\"p\">,</span>\n                <span
    class=\"nb\">str</span><span class=\"p\">(</span><span class=\"nb\">len</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">[</span><span
    class=\"n\">name</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">posts</span><span class=\"p\">)),</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_dict_panel</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"p\">[</span><span class=\"n\">name</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">dict</span><span class=\"p\">()),</span>\n
    \           <span class=\"p\">)</span>\n        <span class=\"k\">return</span>
    <span class=\"n\">table</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div class=\"admonition
    method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"name\" class=\"admonition-title\"
    style=\"margin: 0; padding: .5rem 1rem;\">name <em class=\"small\">method</em></h2>\n<p>The
    name of the feed, used for accessing it in the feeds object.</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">name
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
    class=\"w\"> </span><span class=\"nf\">name</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"nb\">str</span><span class=\"p\">:</span>\n<span class=\"w\">        </span><span
    class=\"sd\">&quot;&quot;&quot;The name of the feed, used for accessing it in
    the feeds object.&quot;&quot;&quot;</span>\n        <span class=\"k\">return</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">name</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition function\">\n<p class=\"admonition-title\">Function</p>\n<h2
    id=\"feeds\" class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">feeds
    <em class=\"small\">function</em></h2>\n<p>feeds cli</p>\n</div>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">feeds
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
    class=\"w\"> </span><span class=\"nf\">feeds</span><span class=\"p\">():</span>\n
    \       <span class=\"s2\">&quot;feeds cli&quot;</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"refresh\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">refresh <em
    class=\"small\">method</em></h2>\n<p>Refresh all of the feeds objects</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">refresh
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
    class=\"w\"> </span><span class=\"nf\">refresh</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">):</span>\n<span class=\"w\">        </span><span
    class=\"sd\">&quot;&quot;&quot;Refresh all of the feeds objects&quot;&quot;&quot;</span>\n
    \       <span class=\"k\">for</span> <span class=\"n\">feed_config</span> <span
    class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">feeds</span><span class=\"p\">:</span>\n
    \           <span class=\"c1\"># Ensure feed has a name, falling back to slug
    if needed</span>\n            <span class=\"k\">if</span> <span class=\"n\">feed_config</span><span
    class=\"o\">.</span><span class=\"n\">name</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span> <span class=\"ow\">and</span> <span class=\"n\">feed_config</span><span
    class=\"o\">.</span><span class=\"n\">slug</span> <span class=\"ow\">is</span>
    <span class=\"ow\">not</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">feed_config</span><span class=\"o\">.</span><span
    class=\"n\">name</span> <span class=\"o\">=</span> <span class=\"n\">feed_config</span><span
    class=\"o\">.</span><span class=\"n\">slug</span><span class=\"o\">.</span><span
    class=\"n\">replace</span><span class=\"p\">(</span><span class=\"s2\">&quot;-&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;_&quot;</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">elif</span> <span class=\"n\">feed_config</span><span
    class=\"o\">.</span><span class=\"n\">name</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span> <span class=\"ow\">and</span> <span class=\"n\">feed_config</span><span
    class=\"o\">.</span><span class=\"n\">slug</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
    class=\"n\">feed_config</span><span class=\"o\">.</span><span class=\"n\">slug</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;archive&quot;</span>\n                <span
    class=\"n\">feed_config</span><span class=\"o\">.</span><span class=\"n\">name</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;archive&quot;</span>\n\n            <span
    class=\"n\">feed</span> <span class=\"o\">=</span> <span class=\"n\">Feed</span><span
    class=\"p\">(</span><span class=\"n\">config</span><span class=\"o\">=</span><span
    class=\"n\">feed_config</span><span class=\"p\">,</span> <span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">markata</span><span class=\"p\">)</span>\n            <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"fm\">__setattr__</span><span class=\"p\">(</span><span
    class=\"n\">feed</span><span class=\"o\">.</span><span class=\"n\">name</span><span
    class=\"p\">,</span> <span class=\"n\">feed</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"_dict_panel\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">_dict_panel
    <em class=\"small\">method</em></h2>\n<p>pretty print configs with rich</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_dict_panel
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
    class=\"w\"> </span><span class=\"nf\">_dict_panel</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"n\">config</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span
    class=\"p\">:</span>\n<span class=\"w\">        </span><span class=\"sd\">&quot;&quot;&quot;pretty
    print configs with rich&quot;&quot;&quot;</span>\n        <span class=\"n\">msg</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;&quot;</span>\n        <span
    class=\"k\">for</span> <span class=\"n\">key</span><span class=\"p\">,</span>
    <span class=\"n\">value</span> <span class=\"ow\">in</span> <span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">items</span><span class=\"p\">():</span>\n
    \           <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">value</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">):</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">len</span><span class=\"p\">(</span><span class=\"n\">value</span><span
    class=\"p\">)</span> <span class=\"o\">&gt;</span> <span class=\"mi\">50</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">value</span> <span
    class=\"o\">=</span> <span class=\"n\">value</span><span class=\"p\">[:</span><span
    class=\"mi\">50</span><span class=\"p\">]</span> <span class=\"o\">+</span> <span
    class=\"s2\">&quot;...&quot;</span>\n                <span class=\"n\">value</span>
    <span class=\"o\">=</span> <span class=\"n\">value</span>\n            <span class=\"n\">msg</span>
    <span class=\"o\">=</span> <span class=\"n\">msg</span> <span class=\"o\">+</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;[grey46]</span><span class=\"si\">{</span><span
    class=\"n\">key</span><span class=\"si\">}</span><span class=\"s2\">[/][magenta3]:[/]
    [grey66]</span><span class=\"si\">{</span><span class=\"n\">value</span><span
    class=\"si\">}</span><span class=\"s2\">[/]</span><span class=\"se\">\\n</span><span
    class=\"s2\">&quot;</span>\n        <span class=\"k\">return</span> <span class=\"n\">msg</span>\n</pre></div>\n\n</pre>\n\n</div>\n\n
    \   </section>\n</article>    </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>feeds.py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"The `markata.plugins.feeds` plugin is used to create
    feed pages, which are lists of\nposts.  The list is generated using a `filter`,
    then each post in the list i\" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n        <meta property=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    <head>\n        <script src=\"https://unpkg.com/@tailwindcss/browser@4\"></script>\n<title>feeds.py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"The `markata.plugins.feeds` plugin is
    used to create feed pages, which are lists of\nposts.  The list is generated using
    a `filter`, then each post in the list i\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n
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
    class=\"title\">\n    <h1 id=\"title\">\n        feeds.py\n    </h1>\n</section></article>\n
    \    </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       feeds.py\n    </h1>\n</section>    <section class=\"body\">\n        <hr
    />\n<p>The <code>markata.plugins.feeds</code> plugin is used to create feed pages,
    which are lists of\nposts.  The list is generated using a <code>filter</code>,
    then each post in the list is\nrendered with a <code>card_template</code> before
    being applied to the <code>body</code> of the\n<code>template</code>.</p>\n<h1
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>This plugin is built-in
    and enabled by default, but in you want to be very\nexplicit you can add it to
    your list of existing plugins.</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    class=\"p\">[</span>\n<span class=\"w\">   </span><span class=\"s2\">&quot;markata.plugins.feeds&quot;</span><span
    class=\"p\">,</span>\n<span class=\"w\">   </span><span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<h1
    id=\"configuration\">Configuration <a class=\"header-anchor\" href=\"#configuration\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<h1 id=\"set-default-template-and-card_template\">set
    default template and card_template <a class=\"header-anchor\" href=\"#set-default-template-and-card_template\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>At the root of the markata.feeds
    config you may set <code>template</code>, and\n<code>card_template</code>.  These
    will become your defaults for every feed you create.\nIf you do not set these,
    markata will use it's defaults.  The defaults are\ndesigned to work for a variety
    of use cases, but are not likely the best for\nall.</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata.feeds_config]</span>\n<span
    class=\"n\">template</span><span class=\"o\">=</span><span class=\"s2\">&quot;pages/templates/archive_template.html&quot;</span>\n<span
    class=\"n\">card_template</span><span class=\"o\">=</span><span class=\"s2\">&quot;plugins/feed_card_template.html&quot;</span>\n</pre></div>\n\n</pre>\n\n<h1
    id=\"pages\">pages <a class=\"header-anchor\" href=\"#pages\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>Underneath of the <code>markata.feeds</code>
    we will create a new map for each page where\nthe name of the map will be the
    name of the page.</p>\n<p>The following config will create a page at <code>/all-posts</code>
    that inclues every\nsingle post.</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.feeds]]</span>\n<span
    class=\"n\">title</span><span class=\"o\">=</span><span class=\"s2\">&quot;All
    Posts&quot;</span>\n<span class=\"n\">slug</span><span class=\"o\">=</span><span
    class=\"s1\">&#39;all&#39;</span>\n<span class=\"n\">filter</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;True&quot;</span>\n</pre></div>\n\n</pre>\n\n<h1 id=\"template\">template
    <a class=\"header-anchor\" href=\"#template\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>The <code>template</code>
    configuration key is a file path to the template that you want\nto use to create
    the feed.  You may set the default template you want to use\nfor all feeds under
    <code>[markata.feeds]</code>, as well as override it inside of each\nfeeds config.</p>\n<p>The
    template is a jinja style template that expects to fill in a <code>title</code>
    and\n<code>body</code> variable.</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"cp\">&lt;!DOCTYPE
    html&gt;</span>\n<span class=\"p\">&lt;</span><span class=\"nt\">html</span> <span
    class=\"na\">lang</span><span class=\"o\">=</span><span class=\"s\">&quot;en&quot;</span><span
    class=\"p\">&gt;</span>\n  <span class=\"p\">&lt;</span><span class=\"nt\">head</span><span
    class=\"p\">&gt;</span>\n    <span class=\"p\">&lt;</span><span class=\"nt\">title</span><span
    class=\"p\">&gt;</span>{{ title }}<span class=\"p\">&lt;/</span><span class=\"nt\">title</span><span
    class=\"p\">&gt;</span>\n  <span class=\"p\">&lt;/</span><span class=\"nt\">head</span><span
    class=\"p\">&gt;</span>\n  <span class=\"p\">&lt;</span><span class=\"nt\">body</span><span
    class=\"p\">&gt;</span>\n    <span class=\"p\">&lt;</span><span class=\"nt\">ul</span><span
    class=\"p\">&gt;</span>\n        {{ body }}\n    <span class=\"p\">&lt;/</span><span
    class=\"nt\">ul</span><span class=\"p\">&gt;</span>\n  <span class=\"p\">&lt;/</span><span
    class=\"nt\">body</span><span class=\"p\">&gt;</span>\n<span class=\"p\">&lt;/</span><span
    class=\"nt\">html</span><span class=\"p\">&gt;</span>\n</pre></div>\n\n</pre>\n\n<div
    class=\"admonition note\">\n<p class=\"admonition-title\">Note</p>\n<p>I highly
    reccomend putting your <code>body</code> in a <code>&lt;ul&gt;</code>, and wrapping
    your\n<code>card_template</code>s in an <code>&lt;li&gt;</code>.</p>\n</div>\n<h1
    id=\"card_template\">card_template <a class=\"header-anchor\" href=\"#card_template\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>All keys available from
    each post is available to put into your jinja\ntemplate.  These can either be
    placed there in your post frontmatter, or\nthrough a plugin that automatically
    adds to the post before the save phase.</p>\n<p>Here is a very simple example
    that would give a link to each post with the\ntitle and date.</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.feeds]]</span>\n<span
    class=\"n\">slug</span><span class=\"o\">=</span><span class=\"s1\">&#39;all&#39;</span>\n<span
    class=\"n\">title</span><span class=\"o\">=</span><span class=\"s1\">&#39;All
    Posts&#39;</span>\n<span class=\"n\">filter</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;True&quot;</span>\n<span class=\"n\">card_template</span><span
    class=\"o\">=</span><span class=\"s1\">&#39;&#39;&#39;</span>\n<span class=\"s1\">&lt;li&gt;</span>\n<span
    class=\"s1\">    &lt;a href={{markata.config.get(&#39;path_prefix&#39;, &#39;&#39;)}}{{slug}}&gt;</span>\n<span
    class=\"s1\">        {{title}}-{{date}}</span>\n<span class=\"s1\">    &lt;/a&gt;</span>\n<span
    class=\"s1\">&lt;/li&gt;</span>\n<span class=\"s1\">&#39;&#39;&#39;</span>\n</pre></div>\n\n</pre>\n\n<h1
    id=\"filter\">filter <a class=\"header-anchor\" href=\"#filter\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>The filter is a python
    expression ran on every post that expects to return a\nboolean.  The variables
    available to this expression are every key in your\nfrontmatter, plus the <code>timedelta</code>
    function, and <code>parse</code> function to more easily\nwork with dates.</p>\n<h1
    id=\"feed-examples\">Feed Examples <a class=\"header-anchor\" href=\"#feed-examples\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>True can be passed in
    to make a feed of all the posts you have.</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.feeds]]</span>\n<span
    class=\"n\">slug</span><span class=\"o\">=</span><span class=\"s1\">&#39;all&#39;</span>\n<span
    class=\"n\">title</span><span class=\"o\">=</span><span class=\"s1\">&#39;All
    Posts&#39;</span>\n<span class=\"n\">filter</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;True&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>You can compare
    against the values of the keys from your frontmatter.  This\nexample creates a
    feed that includes every post where published is <code>True</code>.</p>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.feeds]]</span>\n<span
    class=\"n\">slug</span><span class=\"o\">=</span><span class=\"s1\">&#39;draft&#39;</span>\n<span
    class=\"n\">title</span><span class=\"o\">=</span><span class=\"s1\">&#39;Draft&#39;</span>\n<span
    class=\"n\">filter</span><span class=\"o\">=</span><span class=\"s2\">&quot;published==&#39;False&#39;&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>We
    can also compare against dates.  The\n<a href=\"https://markata.dev/markata/plugins/datetime/\">markata.plugins.datetime</a>\nplugin,
    automatically adds <code>today</code> as today's date and <code>now</code> as
    the current\ndatetime.  These are quite handy to create feeds for scheduled, recent,
    or\ntoday's posts.  The following two examples will create a feed for scheduled\nposts
    and for today's posts respectively.</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.feeds]]</span>\n<span
    class=\"n\">slug</span><span class=\"o\">=</span><span class=\"s1\">&#39;scheduled&#39;</span>\n<span
    class=\"n\">title</span><span class=\"o\">=</span><span class=\"s1\">&#39;Scheduled&#39;</span>\n<span
    class=\"n\">filter</span><span class=\"o\">=</span><span class=\"s2\">&quot;date&gt;today&quot;</span>\n\n<span
    class=\"k\">[[markata.feeds]]</span>\n<span class=\"n\">slug</span><span class=\"o\">=</span><span
    class=\"s1\">&#39;today&#39;</span>\n<span class=\"n\">title</span><span class=\"o\">=</span><span
    class=\"s1\">&#39;Today&#39;</span>\n<span class=\"n\">filter</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;date==today&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>If you
    have list of items in your frontmatter for something like <code>tags</code>, you\ncan
    check for the existence of a tag in the list.</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.feeds]]</span>\n<span
    class=\"n\">slug</span><span class=\"o\">=</span><span class=\"s1\">&#39;python&#39;</span>\n<span
    class=\"n\">title</span><span class=\"o\">=</span><span class=\"s1\">&#39;Python&#39;</span>\n<span
    class=\"n\">filter</span><span class=\"o\">=</span><span class=\"s2\">&quot;date&lt;=today
    and &#39;python&#39; in tags&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>And of
    course you can combine all the things into larger expressions.  Here is\none example
    of the main feed on my blog.</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.feeds]]</span>\n<span
    class=\"n\">slug</span><span class=\"o\">=</span><span class=\"s1\">&#39;blog&#39;</span>\n<span
    class=\"n\">title</span><span class=\"o\">=</span><span class=\"s1\">&#39;Blog&#39;</span>\n<span
    class=\"n\">filter</span><span class=\"o\">=</span><span class=\"s2\">&quot;date&lt;=today
    and templateKey in [&#39;blog-post&#39;] and published ==&#39;True&#39;&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>Here
    is another example that shows my drafts for a particular tag.</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.feeds]]</span>\n<span
    class=\"n\">slug</span><span class=\"o\">=</span><span class=\"s1\">&#39;python-draft&#39;</span>\n<span
    class=\"n\">title</span><span class=\"o\">=</span><span class=\"s1\">&#39;Python
    Draft&#39;</span>\n<span class=\"n\">filter</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;date&lt;=today and &#39;python&#39; in tags and published==&#39;False&#39;&quot;</span>\n</pre></div>\n\n</pre>\n\n<h1
    id=\"defaults\">Defaults <a class=\"header-anchor\" href=\"#defaults\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>By default feeds will
    create one feed page at <code>/archive/</code> that includes all\nposts.</p>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.feeds]]</span>\n<span
    class=\"n\">slug</span><span class=\"o\">=</span><span class=\"s1\">&#39;archive&#39;</span>\n<span
    class=\"n\">title</span><span class=\"o\">=</span><span class=\"s1\">&#39;All
    Posts&#39;</span>\n<span class=\"n\">filter</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;True&quot;</span>\n</pre></div>\n\n</pre>\n\n<hr />\n<div class=\"admonition
    class\">\n<p class=\"admonition-title\">Class</p>\n<h2 id=\"Feed\" class=\"admonition-title\"
    style=\"margin: 0; padding: .5rem 1rem;\">Feed <em class=\"small\">class</em></h2>\n<p>A
    storage class for markata feed objects.</p>\n<h2 id=\"usage\">Usage <a class=\"header-anchor\"
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
    class=\"w\"> </span><span class=\"nn\">markata</span><span class=\"w\"> </span><span
    class=\"kn\">import</span> <span class=\"n\">Markata</span>\n<span class=\"n\">m</span>
    <span class=\"o\">=</span> <span class=\"n\">Markata</span><span class=\"p\">()</span>\n\n<span
    class=\"c1\"># access posts for a feed</span>\n<span class=\"n\">m</span><span
    class=\"o\">.</span><span class=\"n\">feeds</span><span class=\"o\">.</span><span
    class=\"n\">docs</span><span class=\"o\">.</span><span class=\"n\">posts</span>\n\n<span
    class=\"c1\"># access config for a feed</span>\n<span class=\"n\">m</span><span
    class=\"o\">.</span><span class=\"n\">feeds</span><span class=\"o\">.</span><span
    class=\"n\">docs</span><span class=\"o\">.</span><span class=\"n\">config</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">Feed
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">class</span><span
    class=\"w\"> </span><span class=\"nc\">Feed</span><span class=\"p\">(</span><span
    class=\"n\">pydantic</span><span class=\"o\">.</span><span class=\"n\">BaseModel</span><span
    class=\"p\">,</span> <span class=\"n\">JupyterMixin</span><span class=\"p\">):</span>\n<span
    class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;A storage class for
    markata feed objects.</span>\n\n<span class=\"sd\">    ## Usage</span>\n\n<span
    class=\"sd\">    ``` python</span>\n<span class=\"sd\">    from markata import
    Markata</span>\n<span class=\"sd\">    m = Markata()</span>\n\n<span class=\"sd\">
    \   # access posts for a feed</span>\n<span class=\"sd\">    m.feeds.docs.posts</span>\n\n<span
    class=\"sd\">    # access config for a feed</span>\n<span class=\"sd\">    m.feeds.docs.config</span>\n<span
    class=\"sd\">    ```</span>\n<span class=\"sd\">    &quot;&quot;&quot;</span>\n\n
    \   <span class=\"n\">config</span><span class=\"p\">:</span> <span class=\"n\">FeedConfig</span>\n
    \   <span class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"n\">Markata</span>
    <span class=\"o\">=</span> <span class=\"n\">pydantic</span><span class=\"o\">.</span><span
    class=\"n\">Field</span><span class=\"p\">(</span><span class=\"n\">exclude</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n\n
    \   <span class=\"n\">model_config</span> <span class=\"o\">=</span> <span class=\"n\">ConfigDict</span><span
    class=\"p\">(</span>\n        <span class=\"n\">validate_assignment</span><span
    class=\"o\">=</span><span class=\"kc\">False</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">arbitrary_types_allowed</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n        <span class=\"n\">extra</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;allow&quot;</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">str_strip_whitespace</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n        <span class=\"n\">validate_default</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">coerce_numbers_to_str</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n        <span class=\"n\">populate_by_name</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \   <span class=\"p\">)</span>\n\n    <span class=\"nd\">@property</span>\n    <span
    class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">name</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n<span
    class=\"w\">        </span><span class=\"sd\">&quot;&quot;&quot;The name of the
    feed, used for accessing it in the feeds object.&quot;&quot;&quot;</span>\n        <span
    class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">name</span>\n\n
    \   <span class=\"nd\">@property</span>\n    <span class=\"k\">def</span><span
    class=\"w\"> </span><span class=\"nf\">posts</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">):</span>\n        <span class=\"n\">posts</span>
    <span class=\"o\">=</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">map</span><span class=\"p\">(</span><span class=\"s2\">&quot;post&quot;</span><span
    class=\"p\">)</span>\n        <span class=\"k\">if</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">head</span> <span class=\"ow\">is</span> <span class=\"ow\">not</span>
    <span class=\"kc\">None</span> <span class=\"ow\">and</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">tail</span> <span class=\"ow\">is</span> <span class=\"ow\">not</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"n\">head_posts</span>
    <span class=\"o\">=</span> <span class=\"n\">posts</span><span class=\"p\">[:</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">head</span><span class=\"p\">]</span>\n
    \           <span class=\"n\">tail_posts</span> <span class=\"o\">=</span> <span
    class=\"n\">posts</span><span class=\"p\">[</span><span class=\"o\">-</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">tail</span> <span class=\"p\">:]</span>\n
    \           <span class=\"k\">return</span> <span class=\"n\">PrettyList</span><span
    class=\"p\">(</span><span class=\"n\">head_posts</span> <span class=\"o\">+</span>
    <span class=\"n\">tail_posts</span><span class=\"p\">)</span>\n        <span class=\"k\">if</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">head</span> <span class=\"ow\">is</span>
    <span class=\"ow\">not</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">return</span> <span class=\"n\">PrettyList</span><span
    class=\"p\">(</span><span class=\"n\">posts</span><span class=\"p\">[:</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">head</span><span class=\"p\">])</span>\n
    \       <span class=\"k\">if</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">tail</span>
    <span class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n            <span class=\"k\">return</span> <span class=\"n\">PrettyList</span><span
    class=\"p\">(</span><span class=\"n\">posts</span><span class=\"p\">[</span><span
    class=\"o\">-</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">tail</span>
    <span class=\"p\">:])</span>\n        <span class=\"k\">return</span> <span class=\"n\">PrettyList</span><span
    class=\"p\">(</span><span class=\"n\">posts</span><span class=\"p\">)</span>\n\n
    \   <span class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">first</span><span
    class=\"p\">(</span>\n        <span class=\"bp\">self</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span>\n    <span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">list</span><span
    class=\"p\">:</span>\n        <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">posts</span><span class=\"p\">[</span><span
    class=\"mi\">0</span><span class=\"p\">]</span>\n\n    <span class=\"k\">def</span><span
    class=\"w\"> </span><span class=\"nf\">last</span><span class=\"p\">(</span>\n
    \       <span class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span>\n    <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"nb\">list</span><span class=\"p\">:</span>\n        <span class=\"k\">return</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">posts</span><span
    class=\"p\">[</span><span class=\"o\">-</span><span class=\"mi\">1</span><span
    class=\"p\">]</span>\n\n    <span class=\"k\">def</span><span class=\"w\"> </span><span
    class=\"nf\">map</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">func</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;post&quot;</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"n\">args</span><span class=\"p\">):</span>\n        <span class=\"k\">return</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">map</span><span class=\"p\">(</span><span
    class=\"n\">func</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
    class=\"p\">{</span><span class=\"o\">**</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">dict</span><span class=\"p\">(),</span> <span class=\"o\">**</span><span
    class=\"n\">args</span><span class=\"p\">})</span>\n\n    <span class=\"nd\">@property</span>\n
    \   <span class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">__rich_console__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"s2\">&quot;Console&quot;</span><span class=\"p\">:</span>\n
    \       <span class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span>\n\n
    \   <span class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">__rich__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Table</span><span class=\"p\">:</span>\n
    \       <span class=\"n\">table</span> <span class=\"o\">=</span> <span class=\"n\">Table</span><span
    class=\"p\">(</span><span class=\"n\">title</span><span class=\"o\">=</span><span
    class=\"sa\">f</span><span class=\"s2\">&quot;Feed: </span><span class=\"si\">{</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">name</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">)</span>\n\n
    \       <span class=\"n\">table</span><span class=\"o\">.</span><span class=\"n\">add_column</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">justify</span><span class=\"o\">=</span><span class=\"s2\">&quot;right&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">style</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;cyan&quot;</span><span class=\"p\">,</span> <span class=\"n\">no_wrap</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \       <span class=\"n\">table</span><span class=\"o\">.</span><span class=\"n\">add_column</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">justify</span><span class=\"o\">=</span><span class=\"s2\">&quot;left&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">style</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;green&quot;</span><span class=\"p\">)</span>\n        <span
    class=\"n\">table</span><span class=\"o\">.</span><span class=\"n\">add_column</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;published&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">justify</span><span class=\"o\">=</span><span class=\"s2\">&quot;left&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">style</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;green&quot;</span><span class=\"p\">)</span>\n\n        <span
    class=\"k\">for</span> <span class=\"n\">post</span> <span class=\"ow\">in</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">posts</span><span
    class=\"p\">:</span>\n            <span class=\"n\">table</span><span class=\"o\">.</span><span
    class=\"n\">add_row</span><span class=\"p\">(</span><span class=\"n\">post</span><span
    class=\"o\">.</span><span class=\"n\">title</span><span class=\"p\">,</span> <span
    class=\"n\">post</span><span class=\"o\">.</span><span class=\"n\">slug</span><span
    class=\"p\">,</span> <span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">post</span><span class=\"o\">.</span><span class=\"n\">published</span><span
    class=\"p\">))</span>\n\n        <span class=\"k\">return</span> <span class=\"n\">table</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition function\">\n<p class=\"admonition-title\">Function</p>\n<h2
    id=\"pre_render\" class=\"admonition-title\" style=\"margin: 0; padding: .5rem
    1rem;\">pre_render <em class=\"small\">function</em></h2>\n<p>Create the Feeds
    object and attach it to markata.</p>\n</div>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">pre_render <em class='small'>source</em></p>\n<pre
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
    class=\"w\"> </span><span class=\"nf\">pre_render</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"n\">Markata</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">    Create the Feeds object and attach it to markata.</span>\n<span
    class=\"sd\">    &quot;&quot;&quot;</span>\n    <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">feeds</span> <span class=\"o\">=</span>
    <span class=\"n\">Feeds</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div class=\"admonition
    function\">\n<p class=\"admonition-title\">Function</p>\n<h2 id=\"save\" class=\"admonition-title\"
    style=\"margin: 0; padding: .5rem 1rem;\">save <em class=\"small\">function</em></h2>\n<p>Creates
    a new feed page for each page in the config.</p>\n</div>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">save
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
    class=\"w\"> </span><span class=\"nf\">save</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"n\">Markata</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n<span class=\"w\">    </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">    Creates a new feed page for each page in the config.</span>\n<span
    class=\"sd\">    &quot;&quot;&quot;</span>\n    <span class=\"k\">with</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">cache</span>
    <span class=\"k\">as</span> <span class=\"n\">cache</span><span class=\"p\">:</span>\n
    \       <span class=\"k\">for</span> <span class=\"n\">feed</span> <span class=\"ow\">in</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">feeds</span><span
    class=\"o\">.</span><span class=\"n\">values</span><span class=\"p\">():</span>\n
    \           <span class=\"n\">create_page</span><span class=\"p\">(</span>\n                <span
    class=\"n\">markata</span><span class=\"p\">,</span>\n                <span class=\"n\">feed</span><span
    class=\"p\">,</span>\n                <span class=\"n\">cache</span><span class=\"p\">,</span>\n
    \           <span class=\"p\">)</span>\n\n    <span class=\"n\">home</span> <span
    class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">output_dir</span><span class=\"p\">))</span> <span class=\"o\">/</span>
    <span class=\"s2\">&quot;index.html&quot;</span>\n    <span class=\"n\">archive</span>
    <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">output_dir</span><span class=\"p\">))</span> <span class=\"o\">/</span>
    <span class=\"s2\">&quot;archive&quot;</span> <span class=\"o\">/</span> <span
    class=\"s2\">&quot;index.html&quot;</span>\n    <span class=\"k\">if</span> <span
    class=\"ow\">not</span> <span class=\"n\">home</span><span class=\"o\">.</span><span
    class=\"n\">exists</span><span class=\"p\">()</span> <span class=\"ow\">and</span>
    <span class=\"n\">archive</span><span class=\"o\">.</span><span class=\"n\">exists</span><span
    class=\"p\">():</span>\n        <span class=\"n\">shutil</span><span class=\"o\">.</span><span
    class=\"n\">copy</span><span class=\"p\">(</span><span class=\"nb\">str</span><span
    class=\"p\">(</span><span class=\"n\">archive</span><span class=\"p\">),</span>
    <span class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">home</span><span
    class=\"p\">))</span>\n\n    <span class=\"n\">xsl_template</span> <span class=\"o\">=</span>
    <span class=\"n\">get_template</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"n\">feed</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">xsl_template</span><span
    class=\"p\">)</span>\n    <span class=\"n\">xsl</span> <span class=\"o\">=</span>
    <span class=\"n\">xsl_template</span><span class=\"o\">.</span><span class=\"n\">render</span><span
    class=\"p\">(</span>\n        <span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"p\">,</span>\n        <span class=\"n\">__version__</span><span
    class=\"o\">=</span><span class=\"n\">__version__</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">today</span><span class=\"o\">=</span><span class=\"n\">datetime</span><span
    class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
    class=\"n\">today</span><span class=\"p\">(),</span>\n        <span class=\"n\">config</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"p\">,</span>\n    <span class=\"p\">)</span>\n
    \   <span class=\"n\">xsl_file</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">output_dir</span><span
    class=\"p\">)</span> <span class=\"o\">/</span> <span class=\"s2\">&quot;rss.xsl&quot;</span>\n
    \   <span class=\"n\">xsl_file</span><span class=\"o\">.</span><span class=\"n\">write_text</span><span
    class=\"p\">(</span><span class=\"n\">xsl</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition function\">\n<p class=\"admonition-title\">Function</p>\n<h2
    id=\"create_page\" class=\"admonition-title\" style=\"margin: 0; padding: .5rem
    1rem;\">create_page <em class=\"small\">function</em></h2>\n<p>create an html
    unorderd list of posts.</p>\n</div>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">create_page <em class='small'>source</em></p>\n<pre
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
    class=\"w\"> </span><span class=\"nf\">create_page</span><span class=\"p\">(</span>\n
    \   <span class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"n\">Markata</span><span
    class=\"p\">,</span>\n    <span class=\"n\">feed</span><span class=\"p\">:</span>
    <span class=\"n\">Feed</span><span class=\"p\">,</span>\n    <span class=\"n\">cache</span><span
    class=\"p\">,</span>\n<span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span class=\"w\">    </span><span
    class=\"sd\">&quot;&quot;&quot;</span>\n<span class=\"sd\">    create an html
    unorderd list of posts.</span>\n<span class=\"sd\">    &quot;&quot;&quot;</span>\n\n
    \   <span class=\"n\">template</span> <span class=\"o\">=</span> <span class=\"n\">get_template</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">feed</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">template</span><span class=\"p\">)</span>\n
    \   <span class=\"n\">partial_template</span> <span class=\"o\">=</span> <span
    class=\"n\">get_template</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"n\">feed</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">partial_template</span><span
    class=\"p\">)</span>\n    <span class=\"n\">canonical_url</span> <span class=\"o\">=</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">url</span><span class=\"si\">}</span><span
    class=\"s2\">/</span><span class=\"si\">{</span><span class=\"n\">feed</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">slug</span><span class=\"si\">}</span><span class=\"s2\">/&quot;</span>\n\n
    \   <span class=\"n\">key</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">make_hash</span><span class=\"p\">(</span>\n
    \       <span class=\"s2\">&quot;feeds&quot;</span><span class=\"p\">,</span>\n
    \       <span class=\"n\">template</span><span class=\"p\">,</span>\n        <span
    class=\"n\">__version__</span><span class=\"p\">,</span>\n        <span class=\"c1\">#
    cards,</span>\n        <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">url</span><span
    class=\"p\">,</span>\n        <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">description</span><span
    class=\"p\">,</span>\n        <span class=\"n\">feed</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">title</span><span
    class=\"p\">,</span>\n        <span class=\"n\">feed</span><span class=\"o\">.</span><span
    class=\"n\">map</span><span class=\"p\">(</span><span class=\"s2\">&quot;content&quot;</span><span
    class=\"p\">),</span>\n        <span class=\"n\">canonical_url</span><span class=\"p\">,</span>\n
    \       <span class=\"c1\"># datetime.datetime.today(),</span>\n        <span
    class=\"c1\"># markata.config,</span>\n    <span class=\"p\">)</span>\n\n    <span
    class=\"n\">html_key</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">make_hash</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">,</span> <span class=\"s2\">&quot;html&quot;</span><span
    class=\"p\">)</span>\n    <span class=\"n\">html_partial_key</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">make_hash</span><span
    class=\"p\">(</span><span class=\"n\">key</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;partial_html&quot;</span><span class=\"p\">)</span>\n    <span
    class=\"n\">feed_rss_key</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">make_hash</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">,</span> <span class=\"s2\">&quot;rss&quot;</span><span
    class=\"p\">)</span>\n    <span class=\"n\">feed_sitemap_key</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">make_hash</span><span
    class=\"p\">(</span><span class=\"n\">key</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;sitemap&quot;</span><span class=\"p\">)</span>\n\n    <span
    class=\"n\">feed_html_from_cache</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">precache</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"n\">html_key</span><span
    class=\"p\">)</span>\n    <span class=\"n\">feed_html_partial_from_cache</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">precache</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"n\">html_partial_key</span><span class=\"p\">)</span>\n
    \   <span class=\"n\">feed_rss_from_cache</span> <span class=\"o\">=</span> <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">precache</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"n\">feed_rss_key</span><span class=\"p\">)</span>\n    <span class=\"n\">feed_sitemap_from_cache</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">precache</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"n\">feed_sitemap_key</span><span class=\"p\">)</span>\n\n
    \   <span class=\"n\">output_file</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">output_dir</span><span
    class=\"p\">)</span> <span class=\"o\">/</span> <span class=\"n\">feed</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">slug</span> <span class=\"o\">/</span> <span class=\"s2\">&quot;index.html&quot;</span>\n
    \   <span class=\"n\">output_file</span><span class=\"o\">.</span><span class=\"n\">parent</span><span
    class=\"o\">.</span><span class=\"n\">mkdir</span><span class=\"p\">(</span><span
    class=\"n\">exist_ok</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span> <span class=\"n\">parents</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n\n    <span class=\"n\">partial_output_file</span>
    <span class=\"o\">=</span> <span class=\"p\">(</span>\n        <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">output_dir</span><span
    class=\"p\">)</span> <span class=\"o\">/</span> <span class=\"n\">feed</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">slug</span> <span class=\"o\">/</span> <span class=\"s2\">&quot;partial&quot;</span>
    <span class=\"o\">/</span> <span class=\"s2\">&quot;index.html&quot;</span>\n
    \   <span class=\"p\">)</span>\n    <span class=\"n\">partial_output_file</span><span
    class=\"o\">.</span><span class=\"n\">parent</span><span class=\"o\">.</span><span
    class=\"n\">mkdir</span><span class=\"p\">(</span><span class=\"n\">exist_ok</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
    class=\"n\">parents</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n\n    <span class=\"n\">rss_output_file</span> <span class=\"o\">=</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">output_dir</span><span class=\"p\">)</span> <span class=\"o\">/</span>
    <span class=\"n\">feed</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">slug</span> <span class=\"o\">/</span> <span
    class=\"s2\">&quot;rss.xml&quot;</span>\n    <span class=\"n\">rss_output_file</span><span
    class=\"o\">.</span><span class=\"n\">parent</span><span class=\"o\">.</span><span
    class=\"n\">mkdir</span><span class=\"p\">(</span><span class=\"n\">exist_ok</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
    class=\"n\">parents</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n\n    <span class=\"n\">sitemap_output_file</span> <span
    class=\"o\">=</span> <span class=\"p\">(</span>\n        <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">output_dir</span><span
    class=\"p\">)</span> <span class=\"o\">/</span> <span class=\"n\">feed</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">slug</span> <span class=\"o\">/</span> <span class=\"s2\">&quot;sitemap.xml&quot;</span>\n
    \   <span class=\"p\">)</span>\n    <span class=\"n\">sitemap_output_file</span><span
    class=\"o\">.</span><span class=\"n\">parent</span><span class=\"o\">.</span><span
    class=\"n\">mkdir</span><span class=\"p\">(</span><span class=\"n\">exist_ok</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
    class=\"n\">parents</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n\n    <span class=\"n\">from_cache</span> <span class=\"o\">=</span>
    <span class=\"kc\">True</span>\n    <span class=\"k\">if</span> <span class=\"n\">feed_html_from_cache</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \       <span class=\"n\">from_cache</span> <span class=\"o\">=</span> <span class=\"kc\">False</span>\n
    \       <span class=\"n\">feed_html</span> <span class=\"o\">=</span> <span class=\"n\">template</span><span
    class=\"o\">.</span><span class=\"n\">render</span><span class=\"p\">(</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">=</span><span class=\"n\">markata</span><span
    class=\"p\">,</span>\n            <span class=\"n\">__version__</span><span class=\"o\">=</span><span
    class=\"n\">__version__</span><span class=\"p\">,</span>\n            <span class=\"n\">post</span><span
    class=\"o\">=</span><span class=\"n\">feed</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">model_dump</span><span
    class=\"p\">(),</span>\n            <span class=\"n\">url</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">url</span><span class=\"p\">,</span>\n            <span
    class=\"n\">config</span><span class=\"o\">=</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">,</span>\n
    \           <span class=\"n\">feed</span><span class=\"o\">=</span><span class=\"n\">feed</span><span
    class=\"p\">,</span>\n        <span class=\"p\">)</span>\n        <span class=\"n\">cache</span><span
    class=\"o\">.</span><span class=\"n\">set</span><span class=\"p\">(</span><span
    class=\"n\">html_key</span><span class=\"p\">,</span> <span class=\"n\">feed_html</span><span
    class=\"p\">)</span>\n    <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \       <span class=\"n\">feed_html</span> <span class=\"o\">=</span> <span class=\"n\">feed_html_from_cache</span>\n\n
    \   <span class=\"k\">if</span> <span class=\"n\">feed_html_partial_from_cache</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \       <span class=\"n\">from_cache</span> <span class=\"o\">=</span> <span class=\"kc\">False</span>\n
    \       <span class=\"n\">feed_html_partial</span> <span class=\"o\">=</span>
    <span class=\"n\">partial_template</span><span class=\"o\">.</span><span class=\"n\">render</span><span
    class=\"p\">(</span>\n            <span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"p\">,</span>\n            <span class=\"n\">__version__</span><span
    class=\"o\">=</span><span class=\"n\">__version__</span><span class=\"p\">,</span>\n
    \           <span class=\"n\">post</span><span class=\"o\">=</span><span class=\"n\">feed</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">model_dump</span><span class=\"p\">(),</span>\n            <span class=\"n\">url</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">url</span><span
    class=\"p\">,</span>\n            <span class=\"n\">config</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"p\">,</span>\n            <span class=\"n\">feed</span><span class=\"o\">=</span><span
    class=\"n\">feed</span><span class=\"p\">,</span>\n        <span class=\"p\">)</span>\n
    \       <span class=\"n\">cache</span><span class=\"o\">.</span><span class=\"n\">set</span><span
    class=\"p\">(</span><span class=\"n\">html_partial_key</span><span class=\"p\">,</span>
    <span class=\"n\">feed_html_partial</span><span class=\"p\">)</span>\n    <span
    class=\"k\">else</span><span class=\"p\">:</span>\n        <span class=\"n\">feed_html_partial</span>
    <span class=\"o\">=</span> <span class=\"n\">feed_html_partial_from_cache</span>\n\n
    \   <span class=\"k\">if</span> <span class=\"n\">feed_rss_from_cache</span> <span
    class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \       <span class=\"n\">from_cache</span> <span class=\"o\">=</span> <span class=\"kc\">False</span>\n
    \       <span class=\"n\">rss_template</span> <span class=\"o\">=</span> <span
    class=\"n\">get_template</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"n\">feed</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">rss_template</span><span
    class=\"p\">)</span>\n        <span class=\"n\">feed_rss</span> <span class=\"o\">=</span>
    <span class=\"n\">rss_template</span><span class=\"o\">.</span><span class=\"n\">render</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"p\">,</span> <span class=\"n\">feed</span><span
    class=\"o\">=</span><span class=\"n\">feed</span><span class=\"p\">)</span>\n
    \       <span class=\"n\">cache</span><span class=\"o\">.</span><span class=\"n\">set</span><span
    class=\"p\">(</span><span class=\"n\">feed_rss_key</span><span class=\"p\">,</span>
    <span class=\"n\">feed_rss</span><span class=\"p\">)</span>\n    <span class=\"k\">else</span><span
    class=\"p\">:</span>\n        <span class=\"n\">feed_rss</span> <span class=\"o\">=</span>
    <span class=\"n\">feed_rss_from_cache</span>\n\n    <span class=\"k\">if</span>
    <span class=\"n\">feed_sitemap_from_cache</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n        <span class=\"n\">from_cache</span>
    <span class=\"o\">=</span> <span class=\"kc\">False</span>\n        <span class=\"n\">sitemap_template</span>
    <span class=\"o\">=</span> <span class=\"n\">get_template</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">,</span> <span class=\"n\">feed</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">sitemap_template</span><span class=\"p\">)</span>\n        <span class=\"n\">feed_sitemap</span>
    <span class=\"o\">=</span> <span class=\"n\">sitemap_template</span><span class=\"o\">.</span><span
    class=\"n\">render</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">feed</span><span class=\"o\">=</span><span class=\"n\">feed</span><span
    class=\"p\">)</span>\n        <span class=\"n\">cache</span><span class=\"o\">.</span><span
    class=\"n\">set</span><span class=\"p\">(</span><span class=\"n\">feed_sitemap_key</span><span
    class=\"p\">,</span> <span class=\"n\">feed_sitemap</span><span class=\"p\">)</span>\n
    \   <span class=\"k\">else</span><span class=\"p\">:</span>\n        <span class=\"n\">feed_sitemap</span>
    <span class=\"o\">=</span> <span class=\"n\">feed_sitemap_from_cache</span>\n\n
    \   <span class=\"k\">if</span> <span class=\"p\">(</span>\n        <span class=\"n\">from_cache</span>\n
    \       <span class=\"ow\">and</span> <span class=\"n\">output_file</span><span
    class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">()</span>\n
    \       <span class=\"ow\">and</span> <span class=\"n\">partial_output_file</span><span
    class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">()</span>\n
    \       <span class=\"ow\">and</span> <span class=\"n\">rss_output_file</span><span
    class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">()</span>\n
    \       <span class=\"ow\">and</span> <span class=\"n\">sitemap_output_file</span><span
    class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">()</span>\n
    \   <span class=\"p\">):</span>\n        <span class=\"k\">return</span>\n\n    <span
    class=\"n\">output_file</span><span class=\"o\">.</span><span class=\"n\">write_text</span><span
    class=\"p\">(</span><span class=\"n\">feed_html</span><span class=\"p\">)</span>\n
    \   <span class=\"n\">partial_output_file</span><span class=\"o\">.</span><span
    class=\"n\">write_text</span><span class=\"p\">(</span><span class=\"n\">feed_html_partial</span><span
    class=\"p\">)</span>\n    <span class=\"n\">rss_output_file</span><span class=\"o\">.</span><span
    class=\"n\">write_text</span><span class=\"p\">(</span><span class=\"n\">feed_rss</span><span
    class=\"p\">)</span>\n    <span class=\"n\">sitemap_output_file</span><span class=\"o\">.</span><span
    class=\"n\">write_text</span><span class=\"p\">(</span><span class=\"n\">feed_sitemap</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div class=\"admonition
    function\">\n<p class=\"admonition-title\">Function</p>\n<h2 id=\"create_card\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">create_card
    <em class=\"small\">function</em></h2>\n<p>Creates a card for one post based on
    the configured template.  If no\ntemplate is configured it will create one with
    the post title and dates\n(if present).</p>\n</div>\n<div class=\"admonition source
    is-collapsible collapsible-open\">\n<p class=\"admonition-title\">create_card
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
    class=\"w\"> </span><span class=\"nf\">create_card</span><span class=\"p\">(</span>\n
    \   <span class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span
    class=\"p\">,</span>\n    <span class=\"n\">post</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;Post&quot;</span><span class=\"p\">,</span>\n    <span
    class=\"n\">template</span><span class=\"p\">:</span> <span class=\"n\">Optional</span><span
    class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">]</span> <span
    class=\"o\">=</span> <span class=\"kc\">None</span><span class=\"p\">,</span>\n
    \   <span class=\"n\">cache</span><span class=\"o\">=</span><span class=\"kc\">None</span><span
    class=\"p\">,</span>\n<span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Any</span><span class=\"p\">:</span>\n<span class=\"w\">    </span><span
    class=\"sd\">&quot;&quot;&quot;</span>\n<span class=\"sd\">    Creates a card
    for one post based on the configured template.  If no</span>\n<span class=\"sd\">
    \   template is configured it will create one with the post title and dates</span>\n<span
    class=\"sd\">    (if present).</span>\n<span class=\"sd\">    &quot;&quot;&quot;</span>\n
    \   <span class=\"n\">key</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">make_hash</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;feeds&quot;</span><span class=\"p\">,</span> <span class=\"n\">template</span><span
    class=\"p\">,</span> <span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">post</span><span class=\"p\">),</span> <span class=\"n\">post</span><span
    class=\"o\">.</span><span class=\"n\">content</span><span class=\"p\">)</span>\n\n
    \   <span class=\"n\">card</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">precache</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"n\">key</span><span
    class=\"p\">)</span>\n    <span class=\"k\">if</span> <span class=\"n\">card</span>
    <span class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n        <span class=\"k\">return</span> <span class=\"n\">card</span>\n\n
    \   <span class=\"k\">if</span> <span class=\"n\">template</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n        <span class=\"n\">template</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;feeds_config&quot;</span><span class=\"p\">,</span>
    <span class=\"p\">{})</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;card_template&quot;</span><span class=\"p\">,</span>
    <span class=\"kc\">None</span><span class=\"p\">)</span>\n\n    <span class=\"k\">if</span>
    <span class=\"n\">template</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n        <span class=\"k\">if</span> <span class=\"s2\">&quot;date&quot;</span>
    <span class=\"ow\">in</span> <span class=\"n\">post</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">card</span> <span class=\"o\">=</span> <span class=\"n\">textwrap</span><span
    class=\"o\">.</span><span class=\"n\">dedent</span><span class=\"p\">(</span>\n
    \               <span class=\"sa\">f</span><span class=\"s2\">&quot;&quot;&quot;</span>\n<span
    class=\"s2\">                &lt;li class=&#39;post&#39;&gt;</span>\n<span class=\"s2\">
    \               &lt;a href=&quot;/</span><span class=\"si\">{</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">path_prefix</span><span class=\"si\">}{</span><span class=\"n\">post</span><span
    class=\"o\">.</span><span class=\"n\">slug</span><span class=\"si\">}</span><span
    class=\"s2\">/&quot;&gt;</span>\n<span class=\"s2\">                    </span><span
    class=\"si\">{</span><span class=\"n\">post</span><span class=\"o\">.</span><span
    class=\"n\">title</span><span class=\"si\">}</span>\n<span class=\"s2\">                    </span><span
    class=\"si\">{</span><span class=\"n\">post</span><span class=\"o\">.</span><span
    class=\"n\">date</span><span class=\"o\">.</span><span class=\"n\">year</span><span
    class=\"si\">}</span><span class=\"s2\">-</span>\n<span class=\"s2\">                    </span><span
    class=\"si\">{</span><span class=\"n\">post</span><span class=\"o\">.</span><span
    class=\"n\">date</span><span class=\"o\">.</span><span class=\"n\">month</span><span
    class=\"si\">}</span><span class=\"s2\">-</span>\n<span class=\"s2\">                    </span><span
    class=\"si\">{</span><span class=\"n\">post</span><span class=\"o\">.</span><span
    class=\"n\">date</span><span class=\"o\">.</span><span class=\"n\">day</span><span
    class=\"si\">}</span>\n<span class=\"s2\">                &lt;/a&gt;</span>\n<span
    class=\"s2\">                &lt;/li&gt;</span>\n<span class=\"s2\">                &quot;&quot;&quot;</span><span
    class=\"p\">,</span>\n            <span class=\"p\">)</span>\n        <span class=\"k\">else</span><span
    class=\"p\">:</span>\n            <span class=\"n\">card</span> <span class=\"o\">=</span>
    <span class=\"n\">textwrap</span><span class=\"o\">.</span><span class=\"n\">dedent</span><span
    class=\"p\">(</span>\n                <span class=\"sa\">f</span><span class=\"s2\">&quot;&quot;&quot;</span>\n<span
    class=\"s2\">                &lt;li class=&#39;post&#39;&gt;</span>\n<span class=\"s2\">
    \               &lt;a href=&quot;/</span><span class=\"si\">{</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">path_prefix</span><span class=\"si\">}{</span><span class=\"n\">post</span><span
    class=\"o\">.</span><span class=\"n\">slug</span><span class=\"si\">}</span><span
    class=\"s2\">/&quot;&gt;</span>\n<span class=\"s2\">                    </span><span
    class=\"si\">{</span><span class=\"n\">post</span><span class=\"o\">.</span><span
    class=\"n\">title</span><span class=\"si\">}</span>\n<span class=\"s2\">                &lt;/a&gt;</span>\n<span
    class=\"s2\">                &lt;/li&gt;</span>\n<span class=\"s2\">                &quot;&quot;&quot;</span><span
    class=\"p\">,</span>\n            <span class=\"p\">)</span>\n    <span class=\"k\">else</span><span
    class=\"p\">:</span>\n        <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">_template</span> <span class=\"o\">=</span> <span
    class=\"n\">Template</span><span class=\"p\">(</span><span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">template</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">read_text</span><span class=\"p\">())</span>\n
    \       <span class=\"k\">except</span> <span class=\"ne\">FileNotFoundError</span><span
    class=\"p\">:</span>\n            <span class=\"n\">_template</span> <span class=\"o\">=</span>
    <span class=\"n\">Template</span><span class=\"p\">(</span><span class=\"n\">template</span><span
    class=\"p\">)</span>\n        <span class=\"k\">except</span> <span class=\"ne\">OSError</span><span
    class=\"p\">:</span>  <span class=\"c1\"># File name too long</span>\n            <span
    class=\"n\">_template</span> <span class=\"o\">=</span> <span class=\"n\">Template</span><span
    class=\"p\">(</span><span class=\"n\">template</span><span class=\"p\">)</span>\n
    \       <span class=\"n\">card</span> <span class=\"o\">=</span> <span class=\"n\">_template</span><span
    class=\"o\">.</span><span class=\"n\">render</span><span class=\"p\">(</span><span
    class=\"n\">post</span><span class=\"o\">=</span><span class=\"n\">post</span><span
    class=\"p\">,</span> <span class=\"o\">**</span><span class=\"n\">post</span><span
    class=\"o\">.</span><span class=\"n\">to_dict</span><span class=\"p\">())</span>\n
    \   <span class=\"n\">cache</span><span class=\"o\">.</span><span class=\"n\">add</span><span
    class=\"p\">(</span><span class=\"n\">key</span><span class=\"p\">,</span> <span
    class=\"n\">card</span><span class=\"p\">)</span>\n    <span class=\"k\">return</span>
    <span class=\"n\">card</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div class=\"admonition
    class\">\n<p class=\"admonition-title\">Class</p>\n<h2 id=\"Feeds\" class=\"admonition-title\"
    style=\"margin: 0; padding: .5rem 1rem;\">Feeds <em class=\"small\">class</em></h2>\n<p>A
    storage class for all markata Feed objects</p>\n<pre class='wrapper'>\n\n<div
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
    class=\"kn\">import</span> <span class=\"n\">Markata</span>\n<span class=\"n\">markata</span>
    <span class=\"o\">=</span> <span class=\"n\">Markata</span><span class=\"p\">()</span>\n\n<span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">feeds</span>\n\n<span
    class=\"c1\"># access all config</span>\n<span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">feeds</span><span class=\"o\">.</span><span
    class=\"n\">config</span>\n\n<span class=\"c1\"># refresh list of posts in all
    feeds</span>\n<span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">feeds</span><span class=\"o\">.</span><span class=\"n\">refresh</span><span
    class=\"p\">()</span>\n\n<span class=\"c1\"># iterating over feeds gives the name
    of the feed</span>\n<span class=\"k\">for</span> <span class=\"n\">k</span> <span
    class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">feeds</span><span class=\"p\">:</span>\n     <span class=\"nb\">print</span><span
    class=\"p\">(</span><span class=\"n\">k</span><span class=\"p\">)</span>\n\n<span
    class=\"c1\"># project-gallery</span>\n<span class=\"c1\"># docs</span>\n<span
    class=\"c1\"># autodoc</span>\n<span class=\"c1\"># core_modules</span>\n<span
    class=\"c1\"># plugins</span>\n<span class=\"c1\"># archive</span>\n\n<span class=\"c1\">#
    iterate over items like keys and values in a dict, items returns name of</span>\n<span
    class=\"c1\"># feed and a feed object</span>\n<span class=\"k\">for</span> <span
    class=\"n\">k</span><span class=\"p\">,</span> <span class=\"n\">v</span> <span
    class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">feeds</span><span class=\"o\">.</span><span class=\"n\">items</span><span
    class=\"p\">():</span>\n    <span class=\"nb\">print</span><span class=\"p\">(</span><span
    class=\"n\">k</span><span class=\"p\">,</span> <span class=\"nb\">len</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"o\">.</span><span
    class=\"n\">posts</span><span class=\"p\">))</span>\n\n<span class=\"c1\"># project-gallery
    2</span>\n<span class=\"c1\"># docs 6</span>\n<span class=\"c1\"># autodoc 65</span>\n<span
    class=\"c1\"># core_modules 26</span>\n<span class=\"c1\"># plugins 39</span>\n<span
    class=\"c1\"># archive 65</span>\n\n<span class=\"c1\"># values can be iterated
    over in just the same way</span>\n<span class=\"k\">for</span> <span class=\"n\">v</span>
    <span class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">feeds</span><span class=\"o\">.</span><span class=\"n\">values</span><span
    class=\"p\">():</span>\n     <span class=\"nb\">print</span><span class=\"p\">(</span><span
    class=\"nb\">len</span><span class=\"p\">(</span><span class=\"n\">v</span><span
    class=\"o\">.</span><span class=\"n\">posts</span><span class=\"p\">))</span>\n<span
    class=\"c1\"># 2</span>\n<span class=\"c1\"># 6</span>\n<span class=\"c1\"># 65</span>\n<span
    class=\"c1\"># 26</span>\n<span class=\"c1\"># 39</span>\n<span class=\"c1\">#
    65</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">Feeds <em class='small'>source</em></p>\n<pre
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
    class=\"w\"> </span><span class=\"nc\">Feeds</span><span class=\"p\">(</span><span
    class=\"n\">JupyterMixin</span><span class=\"p\">):</span>\n<span class=\"w\">
    \   </span><span class=\"sd\">&quot;&quot;&quot;A storage class for all markata
    Feed objects</span>\n\n<span class=\"sd\">    ``` python</span>\n<span class=\"sd\">
    \   from markata import Markata</span>\n<span class=\"sd\">    markata = Markata()</span>\n\n<span
    class=\"sd\">    markata.feeds</span>\n\n<span class=\"sd\">    # access all config</span>\n<span
    class=\"sd\">    markata.feeds.config</span>\n\n<span class=\"sd\">    # refresh
    list of posts in all feeds</span>\n<span class=\"sd\">    markata.feeds.refresh()</span>\n\n<span
    class=\"sd\">    # iterating over feeds gives the name of the feed</span>\n<span
    class=\"sd\">    for k in markata.feeds:</span>\n<span class=\"sd\">         print(k)</span>\n\n<span
    class=\"sd\">    # project-gallery</span>\n<span class=\"sd\">    # docs</span>\n<span
    class=\"sd\">    # autodoc</span>\n<span class=\"sd\">    # core_modules</span>\n<span
    class=\"sd\">    # plugins</span>\n<span class=\"sd\">    # archive</span>\n\n<span
    class=\"sd\">    # iterate over items like keys and values in a dict, items returns
    name of</span>\n<span class=\"sd\">    # feed and a feed object</span>\n<span
    class=\"sd\">    for k, v in markata.feeds.items():</span>\n<span class=\"sd\">
    \       print(k, len(v.posts))</span>\n\n<span class=\"sd\">    # project-gallery
    2</span>\n<span class=\"sd\">    # docs 6</span>\n<span class=\"sd\">    # autodoc
    65</span>\n<span class=\"sd\">    # core_modules 26</span>\n<span class=\"sd\">
    \   # plugins 39</span>\n<span class=\"sd\">    # archive 65</span>\n\n<span class=\"sd\">
    \   # values can be iterated over in just the same way</span>\n<span class=\"sd\">
    \   for v in markata.feeds.values():</span>\n<span class=\"sd\">         print(len(v.posts))</span>\n<span
    class=\"sd\">    # 2</span>\n<span class=\"sd\">    # 6</span>\n<span class=\"sd\">
    \   # 65</span>\n<span class=\"sd\">    # 26</span>\n<span class=\"sd\">    #
    39</span>\n<span class=\"sd\">    # 65</span>\n<span class=\"sd\">    &quot;&quot;&quot;</span>\n\n
    \   <span class=\"k\">def</span><span class=\"w\"> </span><span class=\"fm\">__init__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"n\">Markata</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n        <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">markata</span> <span class=\"o\">=</span> <span class=\"n\">markata</span>\n
    \       <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span>
    <span class=\"o\">=</span> <span class=\"p\">{</span><span class=\"n\">f</span><span
    class=\"o\">.</span><span class=\"n\">name</span><span class=\"p\">:</span> <span
    class=\"n\">f</span> <span class=\"k\">for</span> <span class=\"n\">f</span> <span
    class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">feeds</span><span
    class=\"p\">}</span>\n        <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">refresh</span><span class=\"p\">()</span>\n\n    <span class=\"k\">def</span><span
    class=\"w\"> </span><span class=\"nf\">refresh</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">):</span>\n<span class=\"w\">        </span><span
    class=\"sd\">&quot;&quot;&quot;Refresh all of the feeds objects&quot;&quot;&quot;</span>\n
    \       <span class=\"k\">for</span> <span class=\"n\">feed_config</span> <span
    class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">feeds</span><span class=\"p\">:</span>\n
    \           <span class=\"c1\"># Ensure feed has a name, falling back to slug
    if needed</span>\n            <span class=\"k\">if</span> <span class=\"n\">feed_config</span><span
    class=\"o\">.</span><span class=\"n\">name</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span> <span class=\"ow\">and</span> <span class=\"n\">feed_config</span><span
    class=\"o\">.</span><span class=\"n\">slug</span> <span class=\"ow\">is</span>
    <span class=\"ow\">not</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">feed_config</span><span class=\"o\">.</span><span
    class=\"n\">name</span> <span class=\"o\">=</span> <span class=\"n\">feed_config</span><span
    class=\"o\">.</span><span class=\"n\">slug</span><span class=\"o\">.</span><span
    class=\"n\">replace</span><span class=\"p\">(</span><span class=\"s2\">&quot;-&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;_&quot;</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">elif</span> <span class=\"n\">feed_config</span><span
    class=\"o\">.</span><span class=\"n\">name</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span> <span class=\"ow\">and</span> <span class=\"n\">feed_config</span><span
    class=\"o\">.</span><span class=\"n\">slug</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
    class=\"n\">feed_config</span><span class=\"o\">.</span><span class=\"n\">slug</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;archive&quot;</span>\n                <span
    class=\"n\">feed_config</span><span class=\"o\">.</span><span class=\"n\">name</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;archive&quot;</span>\n\n            <span
    class=\"n\">feed</span> <span class=\"o\">=</span> <span class=\"n\">Feed</span><span
    class=\"p\">(</span><span class=\"n\">config</span><span class=\"o\">=</span><span
    class=\"n\">feed_config</span><span class=\"p\">,</span> <span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">markata</span><span class=\"p\">)</span>\n            <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"fm\">__setattr__</span><span class=\"p\">(</span><span
    class=\"n\">feed</span><span class=\"o\">.</span><span class=\"n\">name</span><span
    class=\"p\">,</span> <span class=\"n\">feed</span><span class=\"p\">)</span>\n\n
    \   <span class=\"k\">def</span><span class=\"w\"> </span><span class=\"fm\">__iter__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">):</span>\n
    \       <span class=\"k\">return</span> <span class=\"nb\">iter</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">keys</span><span class=\"p\">())</span>\n\n
    \   <span class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">keys</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">):</span>\n
    \       <span class=\"k\">return</span> <span class=\"nb\">iter</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">keys</span><span class=\"p\">())</span>\n\n
    \   <span class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">values</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">):</span>\n
    \       <span class=\"k\">return</span> <span class=\"p\">[</span><span class=\"bp\">self</span><span
    class=\"p\">[</span><span class=\"n\">feed</span><span class=\"p\">]</span> <span
    class=\"k\">for</span> <span class=\"n\">feed</span> <span class=\"ow\">in</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">keys</span><span class=\"p\">()]</span>\n\n
    \   <span class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">items</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">):</span>\n
    \       <span class=\"k\">return</span> <span class=\"p\">[(</span><span class=\"n\">key</span><span
    class=\"p\">,</span> <span class=\"bp\">self</span><span class=\"p\">[</span><span
    class=\"n\">key</span><span class=\"p\">])</span> <span class=\"k\">for</span>
    <span class=\"n\">key</span> <span class=\"ow\">in</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">]</span>\n\n
    \   <span class=\"k\">def</span><span class=\"w\"> </span><span class=\"fm\">__getitem__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"n\">key</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span
    class=\"p\">:</span>\n        <span class=\"k\">return</span> <span class=\"nb\">getattr</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"n\">key</span><span class=\"o\">.</span><span class=\"n\">replace</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;-&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;_&quot;</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">lower</span><span class=\"p\">())</span>\n\n    <span class=\"k\">def</span><span
    class=\"w\"> </span><span class=\"nf\">get</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"n\">key</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">,</span> <span
    class=\"n\">default</span><span class=\"p\">:</span> <span class=\"n\">Any</span>
    <span class=\"o\">=</span> <span class=\"kc\">None</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">Any</span><span class=\"p\">:</span>\n
    \       <span class=\"k\">return</span> <span class=\"nb\">getattr</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"n\">key</span><span class=\"o\">.</span><span class=\"n\">replace</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;-&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;_&quot;</span><span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">lower</span><span class=\"p\">(),</span> <span class=\"n\">default</span><span
    class=\"p\">)</span>\n\n    <span class=\"k\">def</span><span class=\"w\"> </span><span
    class=\"nf\">_dict_panel</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">config</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span class=\"p\">:</span>\n<span
    class=\"w\">        </span><span class=\"sd\">&quot;&quot;&quot;pretty print configs
    with rich&quot;&quot;&quot;</span>\n        <span class=\"n\">msg</span> <span
    class=\"o\">=</span> <span class=\"s2\">&quot;&quot;</span>\n        <span class=\"k\">for</span>
    <span class=\"n\">key</span><span class=\"p\">,</span> <span class=\"n\">value</span>
    <span class=\"ow\">in</span> <span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">items</span><span class=\"p\">():</span>\n            <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">value</span><span
    class=\"p\">,</span> <span class=\"nb\">str</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">if</span> <span class=\"nb\">len</span><span
    class=\"p\">(</span><span class=\"n\">value</span><span class=\"p\">)</span> <span
    class=\"o\">&gt;</span> <span class=\"mi\">50</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">value</span> <span class=\"o\">=</span>
    <span class=\"n\">value</span><span class=\"p\">[:</span><span class=\"mi\">50</span><span
    class=\"p\">]</span> <span class=\"o\">+</span> <span class=\"s2\">&quot;...&quot;</span>\n
    \               <span class=\"n\">value</span> <span class=\"o\">=</span> <span
    class=\"n\">value</span>\n            <span class=\"n\">msg</span> <span class=\"o\">=</span>
    <span class=\"n\">msg</span> <span class=\"o\">+</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;[grey46]</span><span class=\"si\">{</span><span class=\"n\">key</span><span
    class=\"si\">}</span><span class=\"s2\">[/][magenta3]:[/] [grey66]</span><span
    class=\"si\">{</span><span class=\"n\">value</span><span class=\"si\">}</span><span
    class=\"s2\">[/]</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span>\n
    \       <span class=\"k\">return</span> <span class=\"n\">msg</span>\n\n    <span
    class=\"k\">def</span><span class=\"w\"> </span><span class=\"nf\">__rich__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Table</span><span class=\"p\">:</span>\n
    \       <span class=\"n\">table</span> <span class=\"o\">=</span> <span class=\"n\">Table</span><span
    class=\"p\">(</span><span class=\"n\">title</span><span class=\"o\">=</span><span
    class=\"sa\">f</span><span class=\"s2\">&quot;Feeds </span><span class=\"si\">{</span><span
    class=\"nb\">len</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">)</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">)</span>\n\n
    \       <span class=\"n\">table</span><span class=\"o\">.</span><span class=\"n\">add_column</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;Feed&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">justify</span><span class=\"o\">=</span><span class=\"s2\">&quot;right&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">style</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;cyan&quot;</span><span class=\"p\">,</span> <span class=\"n\">no_wrap</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \       <span class=\"n\">table</span><span class=\"o\">.</span><span class=\"n\">add_column</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;posts&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">justify</span><span class=\"o\">=</span><span class=\"s2\">&quot;left&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">style</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;green&quot;</span><span class=\"p\">)</span>\n        <span
    class=\"n\">table</span><span class=\"o\">.</span><span class=\"n\">add_column</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;config&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">style</span><span class=\"o\">=</span><span class=\"s2\">&quot;magenta&quot;</span><span
    class=\"p\">)</span>\n\n        <span class=\"k\">for</span> <span class=\"n\">name</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"p\">:</span>\n            <span class=\"n\">table</span><span
    class=\"o\">.</span><span class=\"n\">add_row</span><span class=\"p\">(</span>\n
    \               <span class=\"n\">name</span><span class=\"p\">,</span>\n                <span
    class=\"nb\">str</span><span class=\"p\">(</span><span class=\"nb\">len</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">[</span><span
    class=\"n\">name</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">posts</span><span class=\"p\">)),</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_dict_panel</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"p\">[</span><span class=\"n\">name</span><span class=\"p\">]</span><span
    class=\"o\">.</span><span class=\"n\">dict</span><span class=\"p\">()),</span>\n
    \           <span class=\"p\">)</span>\n        <span class=\"k\">return</span>
    <span class=\"n\">table</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div class=\"admonition
    method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"name\" class=\"admonition-title\"
    style=\"margin: 0; padding: .5rem 1rem;\">name <em class=\"small\">method</em></h2>\n<p>The
    name of the feed, used for accessing it in the feeds object.</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">name
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
    class=\"w\"> </span><span class=\"nf\">name</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"nb\">str</span><span class=\"p\">:</span>\n<span class=\"w\">        </span><span
    class=\"sd\">&quot;&quot;&quot;The name of the feed, used for accessing it in
    the feeds object.&quot;&quot;&quot;</span>\n        <span class=\"k\">return</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">name</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition function\">\n<p class=\"admonition-title\">Function</p>\n<h2
    id=\"feeds\" class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">feeds
    <em class=\"small\">function</em></h2>\n<p>feeds cli</p>\n</div>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">feeds
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
    class=\"w\"> </span><span class=\"nf\">feeds</span><span class=\"p\">():</span>\n
    \       <span class=\"s2\">&quot;feeds cli&quot;</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"refresh\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">refresh <em
    class=\"small\">method</em></h2>\n<p>Refresh all of the feeds objects</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">refresh
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
    class=\"w\"> </span><span class=\"nf\">refresh</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">):</span>\n<span class=\"w\">        </span><span
    class=\"sd\">&quot;&quot;&quot;Refresh all of the feeds objects&quot;&quot;&quot;</span>\n
    \       <span class=\"k\">for</span> <span class=\"n\">feed_config</span> <span
    class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">feeds</span><span class=\"p\">:</span>\n
    \           <span class=\"c1\"># Ensure feed has a name, falling back to slug
    if needed</span>\n            <span class=\"k\">if</span> <span class=\"n\">feed_config</span><span
    class=\"o\">.</span><span class=\"n\">name</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span> <span class=\"ow\">and</span> <span class=\"n\">feed_config</span><span
    class=\"o\">.</span><span class=\"n\">slug</span> <span class=\"ow\">is</span>
    <span class=\"ow\">not</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">feed_config</span><span class=\"o\">.</span><span
    class=\"n\">name</span> <span class=\"o\">=</span> <span class=\"n\">feed_config</span><span
    class=\"o\">.</span><span class=\"n\">slug</span><span class=\"o\">.</span><span
    class=\"n\">replace</span><span class=\"p\">(</span><span class=\"s2\">&quot;-&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;_&quot;</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">elif</span> <span class=\"n\">feed_config</span><span
    class=\"o\">.</span><span class=\"n\">name</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span> <span class=\"ow\">and</span> <span class=\"n\">feed_config</span><span
    class=\"o\">.</span><span class=\"n\">slug</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
    class=\"n\">feed_config</span><span class=\"o\">.</span><span class=\"n\">slug</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;archive&quot;</span>\n                <span
    class=\"n\">feed_config</span><span class=\"o\">.</span><span class=\"n\">name</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;archive&quot;</span>\n\n            <span
    class=\"n\">feed</span> <span class=\"o\">=</span> <span class=\"n\">Feed</span><span
    class=\"p\">(</span><span class=\"n\">config</span><span class=\"o\">=</span><span
    class=\"n\">feed_config</span><span class=\"p\">,</span> <span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">markata</span><span class=\"p\">)</span>\n            <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"fm\">__setattr__</span><span class=\"p\">(</span><span
    class=\"n\">feed</span><span class=\"o\">.</span><span class=\"n\">name</span><span
    class=\"p\">,</span> <span class=\"n\">feed</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n<div
    class=\"admonition method\">\n<p class=\"admonition-title\">Method</p>\n<h2 id=\"_dict_panel\"
    class=\"admonition-title\" style=\"margin: 0; padding: .5rem 1rem;\">_dict_panel
    <em class=\"small\">method</em></h2>\n<p>pretty print configs with rich</p>\n</div>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_dict_panel
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
    class=\"w\"> </span><span class=\"nf\">_dict_panel</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"n\">config</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span
    class=\"p\">:</span>\n<span class=\"w\">        </span><span class=\"sd\">&quot;&quot;&quot;pretty
    print configs with rich&quot;&quot;&quot;</span>\n        <span class=\"n\">msg</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;&quot;</span>\n        <span
    class=\"k\">for</span> <span class=\"n\">key</span><span class=\"p\">,</span>
    <span class=\"n\">value</span> <span class=\"ow\">in</span> <span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">items</span><span class=\"p\">():</span>\n
    \           <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">value</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">):</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">len</span><span class=\"p\">(</span><span class=\"n\">value</span><span
    class=\"p\">)</span> <span class=\"o\">&gt;</span> <span class=\"mi\">50</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">value</span> <span
    class=\"o\">=</span> <span class=\"n\">value</span><span class=\"p\">[:</span><span
    class=\"mi\">50</span><span class=\"p\">]</span> <span class=\"o\">+</span> <span
    class=\"s2\">&quot;...&quot;</span>\n                <span class=\"n\">value</span>
    <span class=\"o\">=</span> <span class=\"n\">value</span>\n            <span class=\"n\">msg</span>
    <span class=\"o\">=</span> <span class=\"n\">msg</span> <span class=\"o\">+</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;[grey46]</span><span class=\"si\">{</span><span
    class=\"n\">key</span><span class=\"si\">}</span><span class=\"s2\">[/][magenta3]:[/]
    [grey66]</span><span class=\"si\">{</span><span class=\"n\">value</span><span
    class=\"si\">}</span><span class=\"s2\">[/]</span><span class=\"se\">\\n</span><span
    class=\"s2\">&quot;</span>\n        <span class=\"k\">return</span> <span class=\"n\">msg</span>\n</pre></div>\n\n</pre>\n\n</div>\n\n
    \   </section>\n</article>"
  raw.md: ''
published: false
slug: markata/plugins/feeds
title: feeds.py


---

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

``` toml
[[markata.feeds]]
slug='archive'
title='All Posts'
filter="True"
```

---









!!! class
    <h2 id="Feed" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">Feed <em class="small">class</em></h2>

    A storage class for markata feed objects.

    ## Usage

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
    class Feed(pydantic.BaseModel, JupyterMixin):
        """A storage class for markata feed objects.

        ## Usage

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
        markata: Markata = pydantic.Field(exclude=True)

        model_config = ConfigDict(
            validate_assignment=False,
            arbitrary_types_allowed=True,
            extra="allow",
            str_strip_whitespace=True,
            validate_default=True,
            coerce_numbers_to_str=True,
            populate_by_name=True,
        )

        @property
        def name(self) -> str:
            """The name of the feed, used for accessing it in the feeds object."""
            return self.config.name

        @property
        def posts(self):
            posts = self.map("post")
            if self.config.head is not None and self.config.tail is not None:
                head_posts = posts[: self.config.head]
                tail_posts = posts[-self.config.tail :]
                return PrettyList(head_posts + tail_posts)
            if self.config.head is not None:
                return PrettyList(posts[: self.config.head])
            if self.config.tail is not None:
                return PrettyList(posts[-self.config.tail :])
            return PrettyList(posts)

        def first(
            self: "Markata",
        ) -> list:
            return self.posts[0]

        def last(
            self: "Markata",
        ) -> list:
            return self.posts[-1]

        def map(self, func="post", **args):
            return self.markata.map(func, **{**self.config.dict(), **args})

        @property
        def __rich_console__(self) -> "Console":
            return self.markata.console

        def __rich__(self) -> Table:
            table = Table(title=f"Feed: {self.name}")

            table.add_column("Post", justify="right", style="cyan", no_wrap=True)
            table.add_column("slug", justify="left", style="green")
            table.add_column("published", justify="left", style="green")

            for post in self.posts:
                table.add_row(post.title, post.slug, str(post.published))

            return table
    ```









!!! function
    <h2 id="pre_render" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">pre_render <em class="small">function</em></h2>

    Create the Feeds object and attach it to markata.

???+ source "pre_render <em class='small'>source</em>"
    ```python
    def pre_render(markata: Markata) -> None:
        """
        Create the Feeds object and attach it to markata.
        """
        markata.feeds = Feeds(markata)
    ```





!!! function
    <h2 id="save" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">save <em class="small">function</em></h2>

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

        xsl_template = get_template(markata, feed.config.xsl_template)
        xsl = xsl_template.render(
            markata=markata,
            __version__=__version__,
            today=datetime.datetime.today(),
            config=markata.config,
        )
        xsl_file = Path(markata.config.output_dir) / "rss.xsl"
        xsl_file.write_text(xsl)
    ```



!!! function
    <h2 id="create_page" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">create_page <em class="small">function</em></h2>

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

        template = get_template(markata, feed.config.template)
        partial_template = get_template(markata, feed.config.partial_template)
        canonical_url = f"{markata.config.url}/{feed.config.slug}/"

        key = markata.make_hash(
            "feeds",
            template,
            __version__,
            # cards,
            markata.config.url,
            markata.config.description,
            feed.config.title,
            feed.map("content"),
            canonical_url,
            # datetime.datetime.today(),
            # markata.config,
        )

        html_key = markata.make_hash(key, "html")
        html_partial_key = markata.make_hash(key, "partial_html")
        feed_rss_key = markata.make_hash(key, "rss")
        feed_sitemap_key = markata.make_hash(key, "sitemap")

        feed_html_from_cache = markata.precache.get(html_key)
        feed_html_partial_from_cache = markata.precache.get(html_partial_key)
        feed_rss_from_cache = markata.precache.get(feed_rss_key)
        feed_sitemap_from_cache = markata.precache.get(feed_sitemap_key)

        output_file = Path(markata.config.output_dir) / feed.config.slug / "index.html"
        output_file.parent.mkdir(exist_ok=True, parents=True)

        partial_output_file = (
            Path(markata.config.output_dir) / feed.config.slug / "partial" / "index.html"
        )
        partial_output_file.parent.mkdir(exist_ok=True, parents=True)

        rss_output_file = Path(markata.config.output_dir) / feed.config.slug / "rss.xml"
        rss_output_file.parent.mkdir(exist_ok=True, parents=True)

        sitemap_output_file = (
            Path(markata.config.output_dir) / feed.config.slug / "sitemap.xml"
        )
        sitemap_output_file.parent.mkdir(exist_ok=True, parents=True)

        from_cache = True
        if feed_html_from_cache is None:
            from_cache = False
            feed_html = template.render(
                markata=markata,
                __version__=__version__,
                post=feed.config.model_dump(),
                url=markata.config.url,
                config=markata.config,
                feed=feed,
            )
            cache.set(html_key, feed_html)
        else:
            feed_html = feed_html_from_cache

        if feed_html_partial_from_cache is None:
            from_cache = False
            feed_html_partial = partial_template.render(
                markata=markata,
                __version__=__version__,
                post=feed.config.model_dump(),
                url=markata.config.url,
                config=markata.config,
                feed=feed,
            )
            cache.set(html_partial_key, feed_html_partial)
        else:
            feed_html_partial = feed_html_partial_from_cache

        if feed_rss_from_cache is None:
            from_cache = False
            rss_template = get_template(markata, feed.config.rss_template)
            feed_rss = rss_template.render(markata=markata, feed=feed)
            cache.set(feed_rss_key, feed_rss)
        else:
            feed_rss = feed_rss_from_cache

        if feed_sitemap_from_cache is None:
            from_cache = False
            sitemap_template = get_template(markata, feed.config.sitemap_template)
            feed_sitemap = sitemap_template.render(markata=markata, feed=feed)
            cache.set(feed_sitemap_key, feed_sitemap)
        else:
            feed_sitemap = feed_sitemap_from_cache

        if (
            from_cache
            and output_file.exists()
            and partial_output_file.exists()
            and rss_output_file.exists()
            and sitemap_output_file.exists()
        ):
            return

        output_file.write_text(feed_html)
        partial_output_file.write_text(feed_html_partial)
        rss_output_file.write_text(feed_rss)
        sitemap_output_file.write_text(feed_sitemap)
    ```



!!! function
    <h2 id="create_card" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">create_card <em class="small">function</em></h2>

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





!!! class
    <h2 id="Feeds" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">Feeds <em class="small">class</em></h2>

    A storage class for all markata Feed objects

    ``` python
    from markata import Markata
    markata = Markata()

    markata.feeds

    # access all config
    markata.feeds.config

    # refresh list of posts in all feeds
    markata.feeds.refresh()

    # iterating over feeds gives the name of the feed
    for k in markata.feeds:
         print(k)

    # project-gallery
    # docs
    # autodoc
    # core_modules
    # plugins
    # archive

    # iterate over items like keys and values in a dict, items returns name of
    # feed and a feed object
    for k, v in markata.feeds.items():
        print(k, len(v.posts))

    # project-gallery 2
    # docs 6
    # autodoc 65
    # core_modules 26
    # plugins 39
    # archive 65

    # values can be iterated over in just the same way
    for v in markata.feeds.values():
         print(len(v.posts))
    # 2
    # 6
    # 65
    # 26
    # 39
    # 65

???+ source "Feeds <em class='small'>source</em>"
    ```python
    class Feeds(JupyterMixin):
        """A storage class for all markata Feed objects

        ``` python
        from markata import Markata
        markata = Markata()

        markata.feeds

        # access all config
        markata.feeds.config

        # refresh list of posts in all feeds
        markata.feeds.refresh()

        # iterating over feeds gives the name of the feed
        for k in markata.feeds:
             print(k)

        # project-gallery
        # docs
        # autodoc
        # core_modules
        # plugins
        # archive

        # iterate over items like keys and values in a dict, items returns name of
        # feed and a feed object
        for k, v in markata.feeds.items():
            print(k, len(v.posts))

        # project-gallery 2
        # docs 6
        # autodoc 65
        # core_modules 26
        # plugins 39
        # archive 65

        # values can be iterated over in just the same way
        for v in markata.feeds.values():
             print(len(v.posts))
        # 2
        # 6
        # 65
        # 26
        # 39
        # 65
        """

        def __init__(self, markata: Markata) -> None:
            self.markata = markata
            self.config = {f.name: f for f in markata.config.feeds}
            self.refresh()

        def refresh(self):
            """Refresh all of the feeds objects"""
            for feed_config in self.markata.config.feeds:
                # Ensure feed has a name, falling back to slug if needed
                if feed_config.name is None and feed_config.slug is not None:
                    feed_config.name = feed_config.slug.replace("-", "_")
                elif feed_config.name is None and feed_config.slug is None:
                    feed_config.slug = "archive"
                    feed_config.name = "archive"

                feed = Feed(config=feed_config, markata=self.markata)
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

        def get(self, key: str, default: Any = None) -> Any:
            return getattr(self, key.replace("-", "_").lower(), default)

        def _dict_panel(self, config) -> str:
            """pretty print configs with rich"""
            msg = ""
            for key, value in config.items():
                if isinstance(value, str):
                    if len(value) > 50:
                        value = value[:50] + "..."
                    value = value
                msg = msg + f"[grey46]{key}[/][magenta3]:[/] [grey66]{value}[/]\n"
            return msg

        def __rich__(self) -> Table:
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













!!! method
    <h2 id="name" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">name <em class="small">method</em></h2>

    The name of the feed, used for accessing it in the feeds object.

???+ source "name <em class='small'>source</em>"
    ```python
    def name(self) -> str:
            """The name of the feed, used for accessing it in the feeds object."""
            return self.config.name
    ```



















!!! function
    <h2 id="feeds" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">feeds <em class="small">function</em></h2>

    feeds cli

???+ source "feeds <em class='small'>source</em>"
    ```python
    def feeds():
            "feeds cli"
    ```







!!! method
    <h2 id="refresh" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">refresh <em class="small">method</em></h2>

    Refresh all of the feeds objects

???+ source "refresh <em class='small'>source</em>"
    ```python
    def refresh(self):
            """Refresh all of the feeds objects"""
            for feed_config in self.markata.config.feeds:
                # Ensure feed has a name, falling back to slug if needed
                if feed_config.name is None and feed_config.slug is not None:
                    feed_config.name = feed_config.slug.replace("-", "_")
                elif feed_config.name is None and feed_config.slug is None:
                    feed_config.slug = "archive"
                    feed_config.name = "archive"

                feed = Feed(config=feed_config, markata=self.markata)
                self.__setattr__(feed.name, feed)
    ```















!!! method
    <h2 id="_dict_panel" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">_dict_panel <em class="small">method</em></h2>

    pretty print configs with rich

???+ source "_dict_panel <em class='small'>source</em>"
    ```python
    def _dict_panel(self, config) -> str:
            """pretty print configs with rich"""
            msg = ""
            for key, value in config.items():
                if isinstance(value, str):
                    if len(value) > 50:
                        value = value[:50] + "..."
                    value = value
                msg = msg + f"[grey46]{key}[/][magenta3]:[/] [grey66]{value}[/]\n"
            return msg
    ```