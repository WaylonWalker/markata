"""
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

"""

import datetime
import shutil
import textwrap
from dataclasses import dataclass
from functools import lru_cache
from pathlib import Path
from typing import TYPE_CHECKING, Any, List, Optional

import jinja2
import pydantic
import typer
from jinja2 import Template, Undefined
from rich.jupyter import JupyterMixin
from rich.pretty import Pretty
from rich.table import Table

from markata import Markata, __version__, background
from markata.hookspec import hook_impl, register_attr

if TYPE_CHECKING:
    from frontmatter import Post
    from rich.console import Console


class SilentUndefined(Undefined):
    def _fail_with_undefined_error(self, *args, **kwargs):
        return ""


class MarkataFilterError(RuntimeError): ...


class FeedConfig(pydantic.BaseModel, JupyterMixin):
    DEFAULT_TITLE: str = "All Posts"

    title: str = DEFAULT_TITLE
    slug: str = None
    description: Optional[str] = None
    name: Optional[str] = None
    filter: str = "True"
    sort: str = "date"
    reverse: bool = False
    head: Optional[int] = None
    tail: Optional[int] = None
    rss: bool = True
    sitemap: bool = True
    card_template: str = "card.html"
    template: str = "feed.html"
    partial_template: str = "feed_partial.html"
    rss_template: str = "rss.xml"
    sitemap_template: str = "sitemap.xml"
    xsl_template: str = "rss.xsl"

    @pydantic.validator("name", pre=True, always=True)
    def default_name(cls, v, *, values):
        return v or str(values.get("slug")).replace("-", "_")

    @property
    def __rich_console__(self) -> "Console":
        return self.markata.console

    @property
    def __rich__(self) -> Pretty:
        return lambda: Pretty(self)


class FeedsConfig(pydantic.BaseModel):
    feeds: List[FeedConfig] = [FeedConfig(slug="archive")]


class PrettyList(list, JupyterMixin):
    def _repr_pretty_(self):
        return self.__rich__()

    def __rich__(self) -> Pretty:
        return Pretty(self)


@dataclass
class Feed(JupyterMixin):
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
    def __rich_console__(self) -> "Console":
        return self._m.console

    @property
    def name(self):
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
        return self._m.map(func, **{**self.config.dict(), **args})

    def __rich__(self) -> Table:
        table = Table(title=f"Feed: {self.name}")

        table.add_column("Post", justify="right", style="cyan", no_wrap=True)
        table.add_column("slug", justify="left", style="green")
        table.add_column("published", justify="left", style="green")

        for post in self.posts:
            table.add_row(post.title, post.slug, str(post.published))

        return table


class Feeds(JupyterMixin):
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

    def get(self, key: str, default: Any = None) -> Any:
        return getattr(self, key.replace("-", "_").lower(), default)

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


@hook_impl(tryfirst=True)
def config_model(markata: Markata) -> None:
    markata.config_models.append(FeedsConfig)


@hook_impl
@register_attr("feeds")
def pre_render(markata: Markata) -> None:
    """
    Create the Feeds object and attach it to markata.
    """
    markata.feeds = Feeds(markata)


@lru_cache()
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


@hook_impl
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

    if feed_html_from_cache is None:
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
        rss_template = get_template(markata, feed.config.rss_template)
        feed_rss = rss_template.render(markata=markata, feed=feed)
        cache.set(feed_rss_key, feed_rss)
    else:
        feed_rss = feed_rss_from_cache

    if feed_sitemap_from_cache is None:
        sitemap_template = get_template(markata, feed.config.sitemap_template)
        feed_sitemap = sitemap_template.render(markata=markata, feed=feed)
        cache.set(feed_sitemap_key, feed_sitemap)
    else:
        feed_sitemap = feed_sitemap_from_cache

    output_file.write_text(feed_html)
    partial_output_file.write_text(feed_html_partial)
    rss_output_file.write_text(feed_rss)
    sitemap_output_file.write_text(feed_sitemap)


@background.task
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


@hook_impl
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
        markata.console.print("Feeds")
        markata.console.print(feeds)
