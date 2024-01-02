"""
Slugify migration for projects moving from markata<0.5.0 into markata>=0.5.0

to run this script install markata>=0.5.0 and run the following.

``` bash
python -m markata.scripts.migrate_to_slugify
```

Then make sure that you do not explicity turn off slugify and your site is
going to be on to better urls.

``` toml
[markata]
slugify=false
```
"""

from slugify import slugify

from markata import Markata


def routed_slugify(text):
    return "/".join([slugify(s) for s in text.split("/")])


if __name__ == "__main__":
    m = Markata()
    m.config["slugify"] = False

    original_urls = m.map("slug")
    redirects = [
        "/" + o.ljust(60) + "/" + routed_slugify(o)
        for o in original_urls
        if routed_slugify(o) != o
    ]
    assets_dir: str = m.config.assets_dir
    redirects_file = m.config.get.redirects_file
    redirects_file.touch()
    with open(redirects_file, "a") as f:
        f.write("\n")
        f.write("# marakta migrate to slugify redirects\n")
        f.write("\n".join(redirects))
        f.write("\n")
