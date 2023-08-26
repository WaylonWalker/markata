"""
Adds a service_worker to your site.  This will make it installable on mobile,
viewable offline, and potentially more responsive as the user goes between good
and bad connections.

# Configuration

Enable this plugin by adding it to your `markata.toml` hooks list.

``` toml
[markata]
hooks=[
  # your hooks
  "markata.plugins.service_worker",
]
```

If you have any content that you want to  precache, add it to the list of
precache.  You can use devtools, change your network to offline, and see what
files send 404's to the console.  These files likely need precache.

``` toml
[markata]
precache_urls = ['archive-styles.css', 'scroll.css', 'manifest.json']
```

# cache busting

Markata uses the checksum.dirhash of your output directory as the cache key.
This is likely to change and bust the cache on every build.

# pre-caching feeds

You can add and entire feed to your precache, this will automatically load
these posts into the cache anytime someone visits your site and their browser
installs the service worker.

Be nice to your users and don't try to install everything possible in their
cache, but maybe a few that they are most likely to click on.

``` toml
[markata.feeds.recent]
filter="date<today and date>today-timedelta(days=30) and published"
sort="slug"
precache=true
```

> note this assumes that the blog implements a published boolean in each posts
frontmatter.
"""
import copy
from pathlib import Path
from typing import List, Optional, TYPE_CHECKING

from checksumdir import dirhash
from jinja2 import Template
import pydantic

from markata import __version__
from markata.hookspec import hook_impl
from pydantic import ConfigDict

if TYPE_CHECKING:
    from markata import Markata


class ServiceWorkerConfig(pydantic.BaseModel):
    output_dir: pydantic.DirectoryPath = None
    precache_urls: List[str] = ["index.html", "./"]
    precache_posts: bool = False
    precache_feeds: bool = False
    template_file: Optional[Path] = None
    template: Optional[Template] = None
    model_config = ConfigDict(arbitrary_types_allowed=True)

    @pydantic.validator("template_file", always=True, pre=True)
    def validate_template_file(cls, v):
        if v is None:
            return Path(__file__).parent / "default_service_worker_template.js"
        return v


class Config(pydantic.BaseModel):
    service_worker: ServiceWorkerConfig = ServiceWorkerConfig()


@hook_impl(tryfirst=True)
def config_model(markata: "Markata") -> None:
    markata.config_models.append(Config)


@hook_impl(trylast=True)
def render(markata: "Markata") -> None:
    """
    sets precache_urls in markata.config to be used in
    `markata.plugins.service_worker.save`.
    """

    config = markata.config.service_worker

    if config.precache_feeds:
        for feed, config in markata.config.feeds:
            config.precache_urls.append(f"/{feed}/")

    if config.precache_posts:
        for post in markata.map("post", **config):
            config.precache_urls.append(f'/{post.get("slug", "")}/')

    config.precache_urls = list(set(config.precache_urls))


@hook_impl(trylast=True)
def save(markata: "Markata") -> None:
    """
    Renders the service-worker.js file with your precache urls, and dirhash.
    """

    template = Template(markata.config.service_worker.template_file.read_text())
    service_worker_js = template.render(
        __version__=__version__,
        config=copy.deepcopy(markata.config),
        output_dirhash=dirhash(markata.config.output_dir),
    )

    output_file = markata.config.output_dir / "service-worker.js"
    output_file.write_text(service_worker_js)
