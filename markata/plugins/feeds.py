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

from __future__ import annotations

import datetime
import re
import shutil
import textwrap
import warnings
from pathlib import Path
from typing import TYPE_CHECKING
from typing import Any
from typing import List
from typing import Optional
from urllib.request import urlopen

import jinja2
import pydantic
import typer
from jinja2 import Template
from jinja2 import Undefined
from pydantic import ConfigDict
from pydantic import Field
from pydantic import field_validator
from rich.console import Console
from rich.jupyter import JupyterMixin
from rich.pretty import Pretty
from rich.table import Table

from markata import __version__
from markata import background
from markata.hookspec import hook_impl
from markata.hookspec import register_attr
from markata.plugins.jinja_env import get_template
from markata.plugins.jinja_env import get_templates_mtime

if TYPE_CHECKING:
    from frontmatter import Post

    from markata import Markata


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


class SilentUndefined(Undefined):
    """A Jinja2 Undefined subclass that silently returns empty string on errors."""

    def _fail_with_undefined_error(self, *args, **kwargs):
        return ""


class MarkataFilterError(RuntimeError):
    """Raised when a feed filter expression fails."""

    ...


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
    atom: bool = True
    atom_template: str = "atom.xml"
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
        # Get posts from instance state or compute normally
        return self._get_posts()

    def _get_posts(self, override_posts=None):
        """
        Get posts with optional override for pagination.

        Args:
            override_posts: If provided, returns these posts instead of computing

        Returns:
            PrettyList of posts
        """
        if override_posts is not None:
            return PrettyList(override_posts)

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
    skip_htmx_integrity_check: bool = False

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
    _copy_pagination_static_files(markata, Path(markata.config.output_dir))


def _download_htmx_if_needed(markata: Markata) -> None:
    """
    Download HTMX library to static directory if needed with integrity verification.
    """
    import hashlib
    from urllib.error import HTTPError
    from urllib.error import URLError
    from urllib.request import Request

    htmx_version = markata.config.htmx_version
    htmx_filename = "htmx.min.js"
    htmx_static_path = Path(markata.config.output_dir) / "static" / "js" / htmx_filename
    htmx_url = f"https://unpkg.com/htmx.org@{htmx_version}/dist/htmx.min.js"

    # Known SHA-256 hashes for HTMX versions
    HTMX_INTEGRITY_HASHES = {
        "1.9.10": "b3bdcf5c741897a53648b1207fff0469a0d61901429ba1f6e88f98ebd84e669e",
        "2.0.8": "22283ef68cb7545914f0a88a1bdedc7256a703d1d580c1d255217d0a50d31313",
        "2.0.7": "60231ae6ba9db3825eb15a261122d5f55921c4d53b66bf637dc18b4ee27c79f9",
        "2.0.6": "b6768eed4f3af85b73a75054701bd60e17cac718aef2b7f6b254e5e0e2045616",
        "2.0.5": "f601807715bde32e458b73821e16c5641a3d90dfb670f6ebd986f128b8222fcf",
        "2.0.4": "e209dda5c8235479f3166defc7750e1dbcd5a5c1808b7792fc2e6733768fb447",
        "2.0.3": "491955cd1810747d7d7b9ccb936400afb760e06d25d53e4572b64b6563b2784e",
        "2.0.2": "e1746d9759ec0d43c5c284452333a310bb5fd7285ebac4b2dc9bf44d72b5a887",
        "2.0.1": "6d4aaa4b0d3e8b4c91f8d97b92a361a19b1bd4544dea3f668fdc3e62a63995df",
        "2.0.0": "0fc57ba0e655504d282bb6ec1c3d89240cde9f2ce1c393d5b38a95c5bc6da875",
        "1.9.12": "449317ade7881e949510db614991e195c3a099c4c791c24dacec55f9f4a2a452",
        "1.9.11": "d15107cc7f040a9e83b1b66176fd927ad40b5e0255813a03f8ccfeed46ee42b0",
        "1.9.9": "96a334a9570a382cf9c61a1f86d55870ba1c65e166cc5bcae98ddd8cdabeb886",
        "1.9.8": "c4fce4dc5cc9c8c3c9bf1aa788d54bb2cb25cd27114eb06551494ff61c30d6fb",
        "1.9.7": "30c95cb75e7f7c9471c2bf43fa3db0a30a39077764295b15c405869fed7e5764",
        "1.9.6": "cbb723c305cf6d6315c890909815523588509e2e092a59f8cfc4a885829689d5",
        "1.9.5": "76a9887f1ce3bf8f88bea3b327f1e74b9d9b42e1dd9cb8237a87a74261d5d042",
        "1.9.4": "5c88af44013df62fde8a5e4fdf524d8a16834a28b1d15e34ae0994ac27cd4c7e",
        "1.9.3": "8f567d21cbe0553643db48866b2377a3bbb9247f8d924428002c2b847f28b23c",
        "1.9.2": "fd346e9c8639d4624893fc455f2407a09b418301736dd18ebbb07764637fb478",
        "1.9.1": "d7bff1d0f45e3418fa820d8a6f0de1ca5e87562f218a0f06add08652c7691a9c",
        "1.9.0": "97df3adfbf23b873d9a3a80f7143d801a32604ba29de9a33f21a92a171076aa8",
        "1.8.5": "705fb60063bf5270b7077409b848b57ea24d2277b806aa04efea513287bf63a6",
        "1.8.4": "df72edb141a16578945a0356c8a6a37239015251962071639b99b0184691ed1d",
        "1.8.3": "df811b5d27b3dddfec9a858b437b0c7302a56959450f0f9c133ef356c25fcf1c",
        "1.8.2": "91e7fb193c4a6a5d3bb56ed0a7007933664e7803da389a696de61147a6f66058",
        "1.8.1": "1a1c942f7bb50dcc2198b2f3c6cc64199332e32a5ba08e7bd2215aa0a1966a55",
        "1.8.0": "914e05e274362f2e166fc5a8cf6272e2042d9b9e50647678c64c579dcb5fa441",
    }

    expected_hash = HTMX_INTEGRITY_HASHES.get(htmx_version)
    if not expected_hash:
        if markata.config.skip_htmx_integrity_check:
            markata.console.warn(
                f"No integrity hash available for HTMX version {htmx_version}, skipping verification"
            )
            expected_hash = None
        else:
            raise ValueError(
                f"No integrity hash available for HTMX version {htmx_version}. "
                f"You can add 'skip_htmx_integrity_check: true' to your config to skip verification, "
                f"or add the hash to HTMX_INTEGRITY_HASHES in markata/plugins/feeds.py"
            )

    # Download if file doesn't exist
    if not htmx_static_path.exists():
        try:
            with warnings.catch_warnings():
                warnings.simplefilter("ignore", category=ResourceWarning)

                # Ensure static/js directory exists
                htmx_static_path.parent.mkdir(parents=True, exist_ok=True)

                # Download HTMX with timeout and integrity verification
                request = Request(htmx_url, headers={"User-Agent": "Markata/1.0"})
                with urlopen(request, timeout=10) as response:
                    content = response.read()

                    # Verify content integrity if hash is available
                    if expected_hash:
                        actual_hash = hashlib.sha256(content).hexdigest()
                        if actual_hash != expected_hash:
                            raise ValueError(
                                f"HTMX integrity check failed. Expected: {expected_hash}, Got: {actual_hash}"
                            )

                    htmx_static_path.write_bytes(content)

            verification_status = (
                "verified" if expected_hash else "without verification"
            )
            markata.console.print(
                f"Downloaded HTMX {htmx_version} to {htmx_static_path} ({verification_status})"
            )

        except (URLError, HTTPError, ValueError) as e:
            markata.console.error(f"Failed to download HTMX: {e}")
            # Critical security: no fallback to CDN
            raise RuntimeError(
                f"HTMX download failed: {e}. Cannot proceed without verified HTMX."
            )
        except Exception as e:
            markata.console.error(f"Unexpected error downloading HTMX: {e}")
            raise RuntimeError(f"HTMX download failed: {e}")

    return True


def _generate_pagination_js(
    markata: Markata, pagination_config: dict, output_dir: Path
) -> str:
    """
    Generate JavaScript file for pagination and return its path.

    Args:
        markata: Markata instance
        pagination_config: Pagination configuration data
        output_dir: Output directory for JS file

    Returns:
        Path to generated JS file relative to output_dir
    """
    import json

    js_content = f"""// Generated JavaScript for pagination
window.paginationData = {json.dumps(pagination_config)};
"""

    js_dir = output_dir / "static" / "js"
    js_dir.mkdir(parents=True, exist_ok=True)

    js_file = js_dir / "pagination-config.js"
    js_file.write_text(js_content)

    return "/static/js/pagination-config.js"


def _copy_pagination_static_files(markata: Markata, output_dir: Path) -> None:
    """
    Copy pagination static files (JS and CSS) from markata package to output directory.

    Args:
        markata: Markata instance
        output_dir: Output directory for static files
    """
    import importlib.resources

    # Get the markata static directory
    static_package = importlib.resources.files("markata") / "static"

    # Copy pagination.js
    js_src = static_package / "js" / "pagination.js"
    js_dst_dir = output_dir / "static" / "js"
    js_dst_dir.mkdir(parents=True, exist_ok=True)
    js_dst = js_dst_dir / "pagination.js"

    if js_src.is_file():
        js_dst.write_text(js_src.read_text())
        markata.console.print(f"Copied pagination.js to {js_dst}")

    # Copy pagination.css
    css_src = static_package / "css" / "pagination.css"
    css_dst_dir = output_dir / "static" / "css"
    css_dst_dir.mkdir(parents=True, exist_ok=True)
    css_dst = css_dst_dir / "pagination.css"

    if css_src.is_file():
        css_dst.write_text(css_src.read_text())
        markata.console.print(f"Copied pagination.css to {css_dst}")


def _sanitize_feed_slug(slug: str) -> str:
    """
    Sanitize feed slug to prevent path traversal attacks.

    Args:
        slug: User-provided feed slug

    Returns:
        Sanitized slug safe for filesystem use

    Raises:
        ValueError: If slug contains dangerous characters
    """
    import re

    if not slug:
        raise ValueError("Feed slug cannot be empty")

    # Remove path traversal sequences (allow forward slashes for nested paths)
    if ".." in slug or "\\" in slug:
        raise ValueError(f"Invalid characters in feed slug: {slug}")

    # Allow alphanumeric characters, hyphens, underscores, and forward slashes for nested paths
    if not re.match(r"^[a-zA-Z0-9_/-]+$", slug):
        raise ValueError(f"Feed slug contains invalid characters: {slug}")

    # Prevent leading or trailing slashes and double slashes
    if slug.startswith("/") or slug.endswith("/") or "//" in slug:
        raise ValueError(f"Feed slug has invalid slash usage: {slug}")

    # Sanitize by removing any path traversal attempts
    safe_slug = slug.replace("..", "")

    # Additional safety check
    if safe_slug != slug:
        raise ValueError(f"Feed slug attempts path traversal: {slug}")

    return safe_slug


def _ensure_head_links(markata: Markata) -> None:
    """
    Ensure pagination CSS and JS links are in markata.config.head.link
    without duplicating existing links.
    """
    pagination_css_href = "/static/css/pagination.css"
    pagination_js_config_href = "/static/js/pagination-config.js"
    pagination_js_href = "/static/js/pagination.js"
    htmx_version = markata.config.htmx_version
    htmx_static_href = "/static/js/htmx.min.js"

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

    # Check if pagination JS config is already in head.script
    js_config_exists = any(
        get_src(script) == pagination_js_config_href
        for script in markata.config.head.script
    )

    # Check if pagination JS is already in head.script
    js_exists = any(
        get_src(script) == pagination_js_href for script in markata.config.head.script
    )

    # Add JS config link if not already present
    if not js_config_exists:
        markata.config.head.script.append({"src": pagination_js_config_href})

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


@hook_impl
@register_attr("feeds")
def pre_render(markata: Markata) -> None:
    """
    Create the Feeds object and attach it to markata.
    """
    markata.feeds = Feeds(markata)


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

    # Security: Sanitize feed slug to prevent path traversal attacks
    safe_slug = _sanitize_feed_slug(feed.config.slug)
    canonical_url = f"{markata.config.url}/{safe_slug}/"

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
        posts_data = feed.map(
            "(post.slug, str(getattr(post, 'date', '')), getattr(post, 'title', ''))"
        )
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

    output_file = Path(markata.config.output_dir) / safe_slug / "index.html"
    partial_output_file = (
        Path(markata.config.output_dir) / safe_slug / "partial" / "index.html"
    )
    rss_output_file = Path(markata.config.output_dir) / safe_slug / "rss.xml"
    sitemap_output_file = Path(markata.config.output_dir) / safe_slug / "sitemap.xml"
    atom_output_file = Path(markata.config.output_dir) / safe_slug / "atom.xml"

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
            sitemap_template = get_template(
                markata.jinja_env, feed.config.sitemap_template
            )
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

    # Validate per_page to prevent division by zero
    if per_page <= 0:
        raise ValueError(
            f"items_per_page must be a positive integer, got {per_page} for feed '{feed.config.slug}'"
        )

    total_posts = len(posts)

    # Handle empty feeds gracefully
    if total_posts == 0:
        total_pages = 1  # Still create one empty page
    else:
        total_pages = (total_posts + per_page - 1) // per_page

    # Security: Sanitize feed slug to prevent path traversal attacks
    safe_slug = _sanitize_feed_slug(feed.config.slug)

    template = get_template(markata, feed.config.template)
    canonical_url = f"{markata.config.url}/{safe_slug}/"

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

        # Generate JS config file if JS pagination is used
        pagination_js_url = None
        if feed.config.pagination_type == "js":
            pagination_config = {
                "enabled": True,
                "type": feed.config.pagination_type,
                "page": page_num,
                "totalPages": total_pages,
                "totalPosts": total_posts,
                "itemsShown": len(page_posts),
                "feedName": safe_slug,
                "hasNext": page_num < total_pages,
                "config": {
                    "pagination_type": feed.config.pagination_type,
                    "posts_per_page": getattr(feed.config, "posts_per_page", None),
                    "template": getattr(feed.config, "template", None),
                },
            }
            pagination_js_url = _generate_pagination_js(
                markata, pagination_config, Path(markata.config.output_dir)
            )

        # Create a feed object for this page (no state mutation)
        page_feed = Feed(config=feed.config, markata=feed.markata)

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
            output_file = Path(markata.config.output_dir) / safe_slug / "index.html"
        else:
            # Subsequent pages go to numbered subdirectories
            output_file = (
                Path(markata.config.output_dir)
                / safe_slug
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
                pagination_context=pagination_context,
                title=feed.config.title,
                page=page_num,
                total_pages=total_pages,
                total_posts=total_posts,
                has_next=pagination_context["has_next"],
                has_prev=pagination_context["has_prev"],
                next_page=pagination_context["next_page"],
                prev_page=pagination_context["prev_page"],
                feed_name=safe_slug,
                posts=page_posts,
                page_posts=page_posts,
                pagination_js_url=pagination_js_url,
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
                feed_name=safe_slug,
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
