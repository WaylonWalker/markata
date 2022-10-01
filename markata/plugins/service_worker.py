"""
service_worker

Enable this plugin by adding it to your `markata.toml` hooks list.

``` toml
[[ markata ]]
hooks=[
  # your hooks
  "markata.plugins.service_worker",
]
```

"""
import copy
from pathlib import Path
from typing import TYPE_CHECKING

from checksumdir import dirhash
from jinja2 import Template

from markata import __version__
from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from markata import Markata

DEFAULT_PRECACHE_URLS = ["index.html", "./"]


@hook_impl(trylast=True)
def render(markata: "Markata") -> None:

    markata.config["precache_urls"] = markata.config.get("precache_urls", [])
    markata.config["precache_urls"].extend(DEFAULT_PRECACHE_URLS)

    for feed, config in markata.config.get("feeds").items():
        markata.config["precache_urls"].append(f'/{feed}/')

        if config.get('precache', False):
            for post in markata.map("post", **config):
                markata.config["precache_urls"].append(f'/{post.get("slug", "")}/')

    markata.config["precache_urls"] = list(set(markata.config["precache_urls"]))


@hook_impl(trylast=True)
def save(markata: "Markata") -> None:

    if "service_worker_template" in markata.config:
        template_file = markata.config["service_worker_template"]
    else:
        template_file = Path(__file__).parent / "default_service_worker_template.js"
    with open(template_file) as f:
        template = Template(f.read())

    output_dir = Path(markata.config.get("output_dir", "markout"))
    service_worker_file = output_dir / "service-worker.js"
    service_worker_js = template.render(
        __version__=__version__,
        config=copy.deepcopy(markata.config),
        output_dirhash=dirhash(output_dir),
    )

    service_worker_file.write_text(service_worker_js)
