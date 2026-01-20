"""
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

"""

import datetime
import re
import shutil
import textwrap
import warnings
from functools import lru_cache
from pathlib import Path
from typing import TYPE_CHECKING
from urllib.request import urlopen

import frontmatter
import pydantic

from markata.hookspec import hook_impl
from markata import background
from typing import Optional, TYPE_CHECKING, Dict, List, Union, Any
import jinja2
from jinja2 import Environment, Undefined

if TYPE_CHECKING:
    pass  # rich imports available at runtime
else:
    from rich.jupyter import JupyterMixin
    from rich.table import Table
    from rich.console import Console
    from rich.pretty import Pretty
    from rich.jupyter import JupyterMixin
from rich.table import Table
from rich.console import Console
from rich.pretty import Pretty
import typer

# Import JupyterMixin at runtime when needed
if not TYPE_CHECKING:
    JupyterMixin = type("JupyterMixin", (), {})

from pydantic import ConfigDict, Field, field_validator

# Import Markata at module level for type annotations
Markata = None
if TYPE_CHECKING:
    from markata import Markata as MarkataType
from markata.hookspec import register_attr

if TYPE_CHECKING:
    from frontmatter import Post
    from rich.console import Console


def to_pythonic_identifier(name: str) -> str:
    """
    Convert a string to a valid Python identifier.

    This function handles various problematic characters that might appear
    in feed names or slugs, making them suitable for use as Python attribute
    names and dictionary keys.

    Rules applied:
    - Replace spaces, slashes, dots, and other non-alphanumeric characters with underscores
    - Convert to lowercase
    - Remove leading/trailing underscores
    - Ensure the result starts with a letter or underscore
    - Collapse multiple consecutive underscores to a single one

    Examples:
    'project-gallery' -> 'project_gallery'
    'tag/htmx' -> 'tag_htmx'
    'My Feed Name' -> 'my_feed_name'
    '123start' -> '_123start'
    """
    if not name:
        return "_unnamed"

    # Replace non-alphanumeric characters (except underscores) with underscores
    pythonic = re.sub(r"[^a-zA-Z0-9_]", "_", str(name))

    # Convert to lowercase
    pythonic = pythonic.lower()

    # Collapse multiple consecutive underscores
    pythonic = re.sub(r"_+", "_", pythonic)

    # Remove leading and trailing underscores
    pythonic = pythonic.strip("_")

    # Ensure it starts with a letter or underscore (not a digit)
    if pythonic and pythonic[0].isdigit():
        pythonic = "_" + pythonic

    # Handle empty result or result that became empty after processing
    if not pythonic:
        pythonic = "_unnamed"

    return pythonic


if TYPE_CHECKING:

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
    # feed_groups: Dict[str, List[str]] = Field(default_factory=dict)
    # sidebar_feeds: List[str] = Field(default_factory=list)
    card_template: str = "card.html"
    template: str = "feed.html"
    partial_template: str = "feed_partial.html"
    rss_template: str = "rss.xml"
    sitemap_template: str = "sitemap.xml"
    xsl_template: str = "rss.xsl"

    # Pagination configuration
    enabled: bool = False
    items_per_page: int = 10
    pagination_type: str = "htmx"  # htmx, manual, js
    per_page: int = 10  # backwards compatibility

    model_config = ConfigDict(
        validate_assignment=True,  # Config model
        arbitrary_types_allowed=True,
        extra="allow",
        str_strip_whitespace=True,
        validate_default=True,
        coerce_numbers_to_str=True,
        populate_by_name=True,
    )

    @field_validator("name", mode="before")
    @classmethod
    def default_name(cls, v, info) -> str:
        if v:
            return to_pythonic_identifier(str(v))
        slug = info.data.get("slug")
        if not slug:
            raise ValueError("Either name or slug must be provided")
        return to_pythonic_identifier(str(slug))

    @field_validator("slug", mode="before")
    @classmethod
    def default_slug(cls, v, info) -> str:
        if v:
            return v
        name = info.data.get("name")
        if not name:
            raise ValueError("Either name or slug must be provided")
        return str(name).replace("_", "-")

    @property
    def __rich_console__(self) -> "Console":
        return self.markata.console

    @property
    def __rich__(self):
        return lambda: Pretty(self)


class Feed(pydantic.BaseModel, JupyterMixin):
    """A storage class for markata feed objects.

    ## Usage

    ``` python
    if not TYPE_CHECKING:
        from markata import Markata
        m = Markata()

    # access posts for a feed
    m.feeds.docs.posts

    # access config for a feed
    m.feeds.docs.config
    ```
    """

    config: FeedConfig
    markata: Any = Field(exclude=True)

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
        # If this is a paginated page with specific posts, return those
        if hasattr(self, "_page_posts"):
            return PrettyList(self._page_posts)

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


class FeedsConfig(pydantic.BaseModel):
    feeds: List[FeedConfig] = [FeedConfig(slug="archive")]
    htmx_version: str = "2.0.8"

    @property
    def jinja_env(self):
        warnings.warn(
            "The FeedsConfig.jinja_env property is deprecated and will be removed in a future release. "
            "Please use the Markata.jinja_env property instead.",
            DeprecationWarning,
        )

        if hasattr(self, "_jinja_env"):
            return self._jinja_env

        self.env_options.setdefault("loader", self.jinja_loader)
        self.env_options.setdefault("undefined", SilentUndefined)
        self.env_options.setdefault("lstrip_blocks", True)
        self.env_options.setdefault("trim_blocks", True)
        self.env_options.setdefault(
            "bytecode_cache", MarkataTemplateCache(self.template_cache_dir)
        )
        self.env_options.setdefault(
            "auto_reload", False
        )  # Disable auto reload in production

        env = jinja2.Environment(**self.env_options)
        self._jinja_env = env
        return env


class PrettyList(list, JupyterMixin):
    def _repr_pretty_(self):
        return self.__rich__()

    def __rich__(self) -> Pretty:
        return Pretty(self)


@hook_impl(tryfirst=True)
@register_attr("config_models")
def config_model(markata: Markata) -> None:
    markata.config_models.append(FeedsConfig)


@hook_impl(tryfirst=True)
def htmx_config_model(markata: Markata) -> None:
    """Register HTMX configuration model with validation."""

    class HtmxConfig(pydantic.BaseModel):
        version: str = "2.0.8"

        model_config = ConfigDict(
            validate_assignment=True,
            extra="forbid",
        )

    markata.config_models.append(HtmxConfig)


@hook_impl
def configure(markata: Markata) -> None:
    """
    Configure feeds during configuration phase.
    """
    _download_htmx_if_needed(markata)


def _download_htmx_if_needed(markata: Markata) -> None:
    """
    Download HTMX library to static directory if needed.
    """
    htmx_version = markata.config.htmx_version
    htmx_filename = f"htmx.org@{htmx_version}.min.js"
    htmx_static_path = Path(markata.config.output_dir) / "static" / "js" / htmx_filename
    htmx_url = f"https://unpkg.com/htmx.org@{htmx_version}/dist/htmx.min.js"

    # Download if file doesn't exist
    if not htmx_static_path.exists():
        try:
            with warnings.catch_warnings():
                warnings.simplefilter("ignore", category=ResourceWarning)

                # Ensure static/js directory exists
                htmx_static_path.parent.mkdir(parents=True, exist_ok=True)

                # Download HTMX
                with urlopen(htmx_url) as response:
                    content = response.read()
                    htmx_static_path.write_bytes(content)

            markata.console.print(
                f"Downloaded HTMX {htmx_version} to {htmx_static_path}"
            )

        except Exception as e:
            markata.console.warn(f"Failed to download HTMX: {e}")
            # Fallback to CDN if download fails
            return False

    return True


def _ensure_head_links(markata: Markata) -> None:
    """
    Ensure pagination CSS and JS links are in markata.config.head.link
    without duplicating existing links.
    """
    pagination_css_href = "/static/css/pagination.css"
    pagination_js_href = "/static/js/pagination.js"
    htmx_version = markata.config.htmx_version
    htmx_filename = f"htmx.org@{htmx_version}.min.js"
    htmx_static_href = f"/static/js/{htmx_filename}"

    # Try to download HTMX first
    if not _download_htmx_if_needed(markata):
        # Fallback to CDN if download fails
        htmx_cdn_href = f"https://unpkg.com/htmx.org@{htmx_version}"
    else:
        htmx_cdn_href = htmx_static_href

    # Helper function to get href from link (supports both dicts and objects)
    def get_href(link):
        if hasattr(link, "href"):
            return link.href
        return link.get("href", "")

    # Helper function to get src from script (supports both dicts and objects)
    def get_src(script):
        if hasattr(script, "src"):
            return script.src
        return script.get("src", "")

    # Check if pagination CSS is already in head.links
    css_exists = any(
        get_href(link) == pagination_css_href for link in markata.config.head.link
    )

    # Add CSS link if not already present
    if not css_exists:
        markata.config.head.link.append(
            {"rel": "stylesheet", "href": pagination_css_href}
        )

    # Check if pagination JS is already in head.script
    js_exists = any(
        get_src(script) == pagination_js_href for script in markata.config.head.script
    )

    # Add JS link if not already present
    if not js_exists:
        markata.config.head.script.append({"src": pagination_js_href})

    # Check if HTMX is already in head.script
    htmx_exists = any(
        get_src(script) in [htmx_cdn_href, htmx_static_href]
        for script in markata.config.head.script
    )

    # Add HTMX link if not already present
    if not htmx_exists:
        markata.config.head.script.append({"src": htmx_cdn_href})

    # Add HTMX CDN link if not already present
    if not htmx_exists:
        markata.config.head.script.append({"src": htmx_cdn_href})

    # Add CSS link if not already present
    if not css_exists:
        markata.config.head.link.append(
            {"rel": "stylesheet", "href": pagination_css_href}
        )


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
        return markata.jinja_env.get_template(template)
    except jinja2.TemplateNotFound:
        # try to load it as a file
        ...

    try:
        return Template(Path(template).read_text(), undefined=SilentUndefined)
    except FileNotFoundError:
        # default to load it as a string
        ...
    except OSError:  # thrown by File name too long
        # default to load it as a string
        ...
    return Template(template, undefined=SilentUndefined)


@hook_impl
def save(markata: Markata) -> None:
    """
    Creates a new feed page for each page in the config.
    """
    _ensure_head_links(markata)
    with markata.cache as cache:
        for feed in markata.feeds.values():
            if feed.config.enabled:
                create_paginated_feed(
                    markata,
                    feed,
                    cache,
                )
            else:
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
    current_xsl = xsl_file.read_text() if xsl_file.exists() else ""
    if current_xsl != xsl:
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

    current_html = output_file.read_text() if output_file.exists() else ""
    if current_html != feed_html:
        output_file.write_text(feed_html)
    current_partial_html = (
        partial_output_file.read_text() if partial_output_file.exists() else ""
    )
    if current_partial_html != feed_html_partial:
        partial_output_file.write_text(feed_html_partial)
    current_rss = rss_output_file.read_text() if rss_output_file.exists() else ""
    if current_rss != feed_rss:
        rss_output_file.write_text(feed_rss)
    current_sitemap = (
        sitemap_output_file.read_text() if sitemap_output_file.exists() else ""
    )
    if current_sitemap != feed_sitemap:
        sitemap_output_file.write_text(feed_sitemap)


def create_paginated_feed(
    markata: Markata,
    feed: Feed,
    cache,
) -> None:
    """
    Create paginated feed pages.
    """
    posts = feed.posts
    per_page = getattr(feed.config, "items_per_page", feed.config.per_page)
    total_posts = len(posts)
    total_pages = (total_posts + per_page - 1) // per_page

    template = get_template(markata, feed.config.template)
    canonical_url = f"{markata.config.url}/{feed.config.slug}/"

    for page_num in range(1, total_pages + 1):
        start_idx = (page_num - 1) * per_page
        end_idx = start_idx + per_page
        page_posts = posts[start_idx:end_idx]

        # Create pagination context
        pagination_context = {
            "current_page": page_num,
            "total_pages": total_pages,
            "total_posts": total_posts,
            "per_page": per_page,
            "has_prev": page_num > 1,
            "has_next": page_num < total_pages,
            "prev_page": page_num - 1 if page_num > 1 else None,
            "next_page": page_num + 1 if page_num < total_pages else None,
            "pagination_type": feed.config.pagination_type,
        }

        # Create a feed object with just the posts for this page
        page_feed = Feed(config=feed.config, markata=feed.markata)
        # Override the posts property for this page
        page_feed._page_posts = page_posts

        key = markata.make_hash(
            "feeds",
            "paginated",
            template,
            __version__,
            markata.config.url,
            markata.config.description,
            feed.config.title,
            [p.content for p in page_posts],
            canonical_url,
            page_num,
            pagination_context,
        )

        html_key = markata.make_hash(key, "html")
        html_partial_key = markata.make_hash(key, "partial_html")

        # Determine output file paths
        if page_num == 1:
            # First page goes to the main feed index
            output_file = (
                Path(markata.config.output_dir) / feed.config.slug / "index.html"
            )
        else:
            # Subsequent pages go to numbered subdirectories
            output_file = (
                Path(markata.config.output_dir)
                / feed.config.slug
                / str(page_num)
                / "index.html"
            )

        partial_output_file = output_file.parent / "partial" / "index.html"
        output_file.parent.mkdir(exist_ok=True, parents=True)
        partial_output_file.parent.mkdir(exist_ok=True, parents=True)

        # Check cache
        feed_html_from_cache = markata.precache.get(html_key)
        feed_html_partial_from_cache = markata.precache.get(html_partial_key)

        from_cache = True
        if feed_html_from_cache is None:
            from_cache = False
            feed_html = template.render(
                markata=markata,
                __version__=__version__,
                post=feed.config.model_dump(),
                url=markata.config.url,
                config=markata.config,
                feed=page_feed,
                pagination_enabled=True,
                pagination_config=pagination_context,
                title=feed.config.title,
                page=page_num,
                total_pages=total_pages,
                total_posts=total_posts,
                has_next=pagination_context["has_next"],
                has_prev=pagination_context["has_prev"],
                next_page=pagination_context["next_page"],
                prev_page=pagination_context["prev_page"],
                feed_name=feed.config.slug,
                posts=page_posts,
                page_posts=page_posts,
            )
            cache.set(html_key, feed_html)
        else:
            feed_html = feed_html_from_cache

        if feed_html_partial_from_cache is None:
            from_cache = False
            # For HTMX partials, use items-only template to avoid duplicating page structure
            items_partial_template = get_template(markata, "feed_items_partial.html")
            feed_html_partial = items_partial_template.render(
                markata=markata,
                __version__=__version__,
                post=feed.config.model_dump(),
                url=markata.config.url,
                config=markata.config,
                feed=page_feed,
                card_template=feed.config.card_template,
                posts=page_posts,
                page_posts=page_posts,
                has_next=pagination_context["has_next"],
                next_page=pagination_context["next_page"],
                feed_name=feed.config.slug,
                page=page_num,
                total_pages=total_pages,
                total_posts=total_posts,
                pagination_context=pagination_context,
            )
            cache.set(html_partial_key, feed_html_partial)
        else:
            feed_html_partial = feed_html_partial_from_cache

        if from_cache and output_file.exists() and partial_output_file.exists():
            continue

        current_html = output_file.read_text() if output_file.exists() else ""
        if current_html != feed_html:
            output_file.write_text(feed_html)

        current_partial_html = (
            partial_output_file.read_text() if partial_output_file.exists() else ""
        )
        if current_partial_html != feed_html_partial:
            partial_output_file.write_text(feed_html_partial)


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
    if template is None:
        template = markata.config.get("feeds_config", {}).get("card_template", None)

    # Get template modification time if template exists
    template_mtime = 0
    if template:
        template_path = None
        # Check user template paths first
        for path in markata.jinja_env.template_paths:
            potential_path = Path(path) / template
            if potential_path.exists():
                template_path = potential_path
                break

        # Check package templates if not found in user paths
        if not template_path:
            import importlib

            package_template = (
                importlib.resources.files("markata") / "templates" / template
            )
            if package_template.exists():
                template_path = package_template

        if template_path:
            template_mtime = template_path.stat().st_mtime

    key = markata.make_hash(
        "feeds", template, str(post), post.content, str(template_mtime)
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


@hook_impl
def cli(app: typer.Typer, markata: "Markata") -> None:
    feeds_app = typer.Typer()
    app.add_typer(feeds_app, name="feeds")

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
                feed_config.name = to_pythonic_identifier(str(feed_config.slug))
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
        return getattr(self, to_pythonic_identifier(str(key)))

    def get(self, key: str, default: Any = None) -> Any:
        return getattr(self, to_pythonic_identifier(str(key)), default)

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
