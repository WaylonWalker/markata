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

### set default template and card_template

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

### pages

Underneath of the `markata.feeds` we will create a new map for each page where
the name of the map will be the name of the page.


The following config will create a page at `/all-posts` that inclues every
single post.

``` toml
[markata.feeds.all-posts]
filter="True"
```

### template

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

### card_template

All keys available from each post is available to put into your jinja
template.  These can either be placed there in your post frontmatter, or
through a plugin that automatically adds to the post before the save phase.

Here is a very simple example that would give a link to each post with the
title and date.

``` toml
[markata.feeds.all]
filter="True"
card_template='<li><a href={{markata.config.get('path_prefix', '')}}{{slug}}>{{title}}-{{date}}</a></li>'
```

### filter

The filter is a python expression ran on every post that expects to return a
boolean.  The variables available to this expression are every key in your
frontmatter, plus the `timedelta` function, and `parse` function to more easily
work with dates.

## Feed Examples

True can be passed in to make a feed of all the posts you have.

``` toml
[markata.feeds.archive]
filter="True"
```

You can compare against the values of the keys from your frontmatter.  This
example creates a feed that includes every post where published is `True`.

``` toml
[markata.feeds.draft]
filter="published=='False'"
```

We can also compare against dates.  The
[markata.plugins.datetime](https://markata.dev/markata/plugins/datetime/)
plugin, automatically adds `today` as today's date and `now` as the current
datetime.  These are quite handy to create feeds for scheduled, recent, or
today's posts.  The following two examples will create a feed for scheduled
posts and for today's posts respectively.

``` toml
[markata.feeds.scheduled]
filter="date>today"

[markata.feeds.today]
filter="date==today"
```

If you have list of items in your frontmatter for something like `tags`, you
can check for the existence of a tag in the list.

``` toml
[markata.feeds.python]
filter="date<=today and 'python' in tags"
```

And of course you can combine all the things into larger expressions.  Here is one example of the main feed on my blog.

``` toml
[markata.feeds.blog]
filter="date<=today and templateKey in ['blog-post'] and published =='True'"
```

Here is another example that shows my drafts for a particular tag.

``` toml
[markata.feeds.python-draft]
filter="date<=today and 'python' in tags and published=='False'"
```

## Defaults

By default feeds will create one feed page at `/archive/` that includes all
posts.

[markata.feeds.archive]
filter="True"

"""
import datetime
import shutil
import textwrap
from dataclasses import dataclass
from pathlib import Path
from typing import TYPE_CHECKING, Any, Dict, List, Optional, Union

import typer
from jinja2 import Template, Undefined
from rich import print as rich_print
from rich.table import Table

from markata import Markata, __version__
from markata.hookspec import hook_impl, register_attr

if TYPE_CHECKING:
    from frontmatter import Post


class SilentUndefined(Undefined):
    def _fail_with_undefined_error(self, *args, **kwargs):
        return ""


class MarkataFilterError(RuntimeError):
    ...


@dataclass
class Feed:
    """
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
    """

    name: str
    config: Dict
    posts: list
    _m: Markata

    def map(self, func="post", **args):
        return self._m.map(func, **{**self.config, **args})


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

    # iterate over items like keys and values in a dict, items returns name of feed and a feed object
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

    def __init__(self, markata: Markata):
        self._m = markata
        self.refresh()

    def refresh(self) -> None:
        """
        Refresh all of the feeds objects
        """
        self.config = self._m.config.get("feeds", dict())
        for page, page_conf in self.config.items():
            name = page.replace("-", "_").lower()
            feed = Feed(
                name=name,
                posts=self._m.map("post", **page_conf),
                config=page_conf,
                _m=self._m,
            )
            self.__setattr__(name, feed)

    def __iter__(self):
        return iter(self.config.keys())

    def keys(self):
        return iter(self.config.keys())

    def values(self):
        return [self[feed] for feed in self.config.keys()]

    def items(self):
        return [(key, self[key]) for key in self.config.keys()]

    def __getitem__(self, key: str) -> Any:
        return getattr(self, key.replace("-", "_").lower())

    def _dict_panel(self, config) -> str:
        """
        pretty print configs with rich
        """
        msg = ""
        for key, value in config.items():
            msg = msg + f"[grey46]{key}[/][magenta3]:[/] [grey66]{value}[/]\n"
        return msg

    def __rich__(self) -> Table:
        from rich.table import Table

        table = Table(title=f"Feeds {len(self.config)}")

        table.add_column("Feed", justify="right", style="cyan", no_wrap=True)
        table.add_column("posts", justify="left", style="green")
        table.add_column("config", style="magenta")

        for name in self.config.keys():
            table.add_row(
                name, str(len(self[name].posts)), self._dict_panel(self.config[name])
            )
        return table


@hook_impl
@register_attr("feeds")
def configure(markata: Markata) -> None:
    """
    configure the default values for the feeds plugin
    """
    if "feeds" not in markata.config.keys():
        markata.config["feeds"] = dict()
    config = markata.config.get("feeds", dict())
    if "archive" not in config.keys():
        config["archive"] = dict()
        config["archive"]["filter"] = "True"

    default_post_template = markata.config.get("feeds_config", {}).get(
        "template", Path(__file__).parent / "default_post_template.html"
    )

    for page, page_conf in config.items():
        if "template" not in page_conf.keys():
            page_conf["template"] = default_post_template


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
    feeds = markata.config.get("feeds", {})

    description = markata.get_config("description") or ""
    url = markata.get_config("url") or ""

    for page, page_conf in feeds.items():

        create_page(
            markata,
            page,
            description=description,
            url=url,
            **page_conf,
        )

    home = Path(str(markata.config["output_dir"])) / "index.html"
    archive = Path(str(markata.config["output_dir"])) / "archive" / "index.html"
    if not home.exists() and archive.exists():
        shutil.copy(str(archive), str(home))


def create_page(
    markata: Markata,
    page: str,
    tags: Optional[List] = None,
    published: str = "True",
    template: Optional[Union[Path, str]] = None,
    card_template: Optional[str] = None,
    filter: Optional[str] = None,
    description: Optional[str] = None,
    url: Optional[str] = None,
    title: Optional[str] = "feed",
    sort: str = "True",
    reverse: bool = False,
    **rest,
) -> None:
    """
    create an html unorderd list of posts.
    """

    posts = markata.feeds[page].posts
    cards = [create_card(markata, post, card_template) for post in posts]
    cards.insert(0, "<ul>")
    cards.append("</ul>")

    # if template is None:
    #     template = Path(__file__).parent / "default_post_template.html"

    with open(template) as f:
        template = Template(f.read(), undefined=SilentUndefined)
    output_file = Path(markata.config["output_dir"]) / page / "index.html"
    canonical_url = f"{url}/{page}/"
    output_file.parent.mkdir(exist_ok=True, parents=True)

    with open(output_file, "w+") as f:
        f.write(
            template.render(
                __version__=__version__,
                body="".join(cards),
                url=url,
                description=description,
                title=title,
                canonical_url=canonical_url,
                today=datetime.datetime.today(),
                config=markata.config,
            )
        )


def create_card(
    markata: "Markata", post: "Post", template: Optional[str] = None
) -> Any:
    """
    Creates a card for one post based on the configured template.  If no
    template is configured it will create one with the post title and dates (if present).
    """
    if template is None:
        template = markata.config.get("feeds_config", {}).get("card_template", None)

    if template is None:
        if "date" in post.keys():
            return textwrap.dedent(
                f"""
                <li class='post'>
                <a href="/{markata.config.get('path_prefix', '')}{post['slug']}/">
                    {post['title']} {post['date'].year}-{post['date'].month}-{post['date'].day}
                </a>
                </li>
                """
            )
        else:
            return textwrap.dedent(
                f"""
                <li class='post'>
                <a href="/{markata.config.get('path_prefix', '')}{post['slug']}/">
                    {post['title']}
                </a>
                </li>
                """
            )
    try:
        _template = Template(Path(template).read_text())
    except FileNotFoundError:
        _template = Template(template)
    return _template.render(**post.to_dict())


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
        rich_print(feeds)
