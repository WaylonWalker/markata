---
date: 2025-12-09
description: "The plugin is used to create feed pages, which are lists of posts. The
  list is generated using a , then each post in the list is rendered with a before
  being\u2026"
published: false
slug: markata/plugins/feeds
title: feeds.py


---

---

The `markata.plugins.feeds` plugin is used to create feed pages, which are lists of
posts.  The list is generated using a `filter`, then each post in the list is
rendered with a `card_template` before being applied to the `body` of the
`template`.

## Installation

This plugin is built-in and enabled by default, but in you want to be very
explicit you can add it to your list of existing plugins.

``` toml
hooks = [
   "markata.plugins.feeds",
   ]
```

## Configuration

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
!!! class
    <h2 id="MarkataTemplateCache" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">MarkataTemplateCache <em class="small">class</em></h2>

    Template bytecode cache for improved performance.

???+ source "MarkataTemplateCache <em class='small'>source</em>"
    ```python
    class MarkataTemplateCache(jinja2.BytecodeCache):
        """Template bytecode cache for improved performance."""

        def __init__(self, directory):
            self.directory = Path(directory)
            self.directory.mkdir(parents=True, exist_ok=True)

        def load_bytecode(self, bucket):
            filename = self.directory / f"{bucket.key}.cache"
            if filename.exists():
                with open(filename, "rb") as f:
                    bucket.bytecode_from_string(f.read())

        def dump_bytecode(self, bucket):
            filename = self.directory / f"{bucket.key}.cache"
            with open(filename, "wb") as f:
                f.write(bucket.bytecode_to_string())
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

        xsl_template = get_template(markata.jinja_env, feed.config.xsl_template)
        xsl = xsl_template.render(
            markata=markata,
            __version__=__version__,
            today=datetime.datetime.today(),
            config=markata.config,
        )
        xsl_file = Path(markata.config.output_dir) / "rss.xsl"
        # Only read file if it exists and we need to compare
        should_write = True
        if xsl_file.exists():
            current_xsl = xsl_file.read_text()
            should_write = current_xsl != xsl

        if should_write:
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

        template = get_template(markata.jinja_env, feed.config.template)
        partial_template = get_template(markata.jinja_env, feed.config.partial_template)
        canonical_url = f"{markata.config.url}/{feed.config.slug}/"

        # Get templates mtime to bust cache when any template changes
        templates_mtime = get_templates_mtime(markata.jinja_env)

        # Use simpler hash for posts instead of expensive str(post.to_dict())
        # Hash just the essential post identifiers: slug + content_hash
        cache_key_posts = f"feed_hash_posts_{feed.config.slug}"
        if not hasattr(markata, "_feed_hash_cache"):
            markata._feed_hash_cache = {}

        if cache_key_posts not in markata._feed_hash_cache:
            # Use post slugs and published dates instead of full to_dict()
            # This provides a stable, lightweight cache key
            posts_data = feed.map("(post.slug, str(getattr(post, 'date', '')), getattr(post, 'title', ''))")
            markata._feed_hash_cache[cache_key_posts] = str(sorted(posts_data))

        posts_hash_data = markata._feed_hash_cache[cache_key_posts]

        key = markata.make_hash(
            "feeds",
            template,
            __version__,
            markata.config.url,
            markata.config.description,
            feed.config.title,
            posts_hash_data,  # Use cached post data
            canonical_url,
            str(templates_mtime),  # Track template file changes
            # datetime.datetime.today(),
            # markata.config,
        )

        html_key = markata.make_hash(key, "html")
        html_partial_key = markata.make_hash(key, "partial_html")
        feed_rss_key = markata.make_hash(key, "rss")
        feed_sitemap_key = markata.make_hash(key, "sitemap")
        feed_atom_key = markata.make_hash(key, "atom")

        feed_html_from_cache = markata.precache.get(html_key)
        feed_html_partial_from_cache = markata.precache.get(html_partial_key)
        feed_rss_from_cache = markata.precache.get(feed_rss_key)
        feed_sitemap_from_cache = markata.precache.get(feed_sitemap_key)
        feed_atom_from_cache = markata.precache.get(feed_atom_key)

        output_file = Path(markata.config.output_dir) / feed.config.slug / "index.html"
        partial_output_file = (
            Path(markata.config.output_dir) / feed.config.slug / "partial" / "index.html"
        )
        rss_output_file = Path(markata.config.output_dir) / feed.config.slug / "rss.xml"
        sitemap_output_file = (
            Path(markata.config.output_dir) / feed.config.slug / "sitemap.xml"
        )
        atom_output_file = (
            Path(markata.config.output_dir) / feed.config.slug / "atom.xml"
        )

        # Create all directories in one batch
        partial_output_file.parent.mkdir(exist_ok=True, parents=True)

        from_cache = True

        # ---------- HTML ----------
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

        # ---------- Partial HTML ----------
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

        # ---------- RSS ----------
        if feed.config.rss:
            if feed_rss_from_cache is None:
                from_cache = False
                rss_template = get_template(markata.jinja_env, feed.config.rss_template)
                feed_rss = rss_template.render(markata=markata, feed=feed)
                cache.set(feed_rss_key, feed_rss)
            else:
                feed_rss = feed_rss_from_cache
        else:
            feed_rss = None

        # ---------- Sitemap ----------
        if feed.config.sitemap:
            if feed_sitemap_from_cache is None:
                from_cache = False
                sitemap_template = get_template(markata.jinja_env, feed.config.sitemap_template)
                feed_sitemap = sitemap_template.render(markata=markata, feed=feed)
                cache.set(feed_sitemap_key, feed_sitemap)
            else:
                feed_sitemap = feed_sitemap_from_cache
        else:
            feed_sitemap = None

        # ---------- Atom ----------
        if feed.config.atom:
            if feed_atom_from_cache is None:
                from_cache = False
                atom_template = get_template(markata.jinja_env, feed.config.atom_template)
                feed_atom = atom_template.render(
                    markata=markata,
                    feed=feed,
                    datetime=datetime,  # ⭐ so the template can use datetime
                )
                cache.set(feed_atom_key, feed_atom)
            else:
                feed_atom = feed_atom_from_cache
            # If everything came from cache and files exist, bail early
            if (
                from_cache
                and output_file.exists()
                and partial_output_file.exists()
                and (not feed.config.rss or rss_output_file.exists())
                and (not feed.config.sitemap or sitemap_output_file.exists())
                and (not feed.config.atom or atom_output_file.exists())
            ):
                return

        # Write HTML
        current_html = output_file.read_text() if output_file.exists() else ""
        if current_html != feed_html:
            output_file.write_text(feed_html)

        # Write partial HTML
        current_partial_html = (
            partial_output_file.read_text() if partial_output_file.exists() else ""
        )
        if current_partial_html != feed_html_partial:
            partial_output_file.write_text(feed_html_partial)

        # Write RSS (if enabled)
        if feed_rss is not None:
            current_rss = rss_output_file.read_text() if rss_output_file.exists() else ""
            if current_rss != feed_rss:
                rss_output_file.write_text(feed_rss)

        # Write sitemap (if enabled)
        if feed_sitemap is not None:
            current_sitemap = (
                sitemap_output_file.read_text() if sitemap_output_file.exists() else ""
            )
            if current_sitemap != feed_sitemap:
                sitemap_output_file.write_text(feed_sitemap)

        # Write Atom (if enabled)
        if feed_atom is not None:
            current_atom = atom_output_file.read_text() if atom_output_file.exists() else ""
            if current_atom != feed_atom:
                atom_output_file.write_text(feed_atom)
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
        if template is None:
            template = markata.config.get("feeds_config", {}).get("card_template", None)

        # Get templates mtime to bust cache when any template changes
        templates_mtime = get_templates_mtime(markata.jinja_env)

        key = markata.make_hash(
            "feeds", template, str(post.to_dict()), str(templates_mtime)
        )

        card = markata.precache.get(key)
        if card is not None:
            return card

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
            except OSError:  # thrown by File name too long
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