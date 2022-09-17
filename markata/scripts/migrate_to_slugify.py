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
from pathlib import Path

from slugify import slugify

from markata import Markata

if __name__ == "__main__":
    m = Markata()
    m.config["slugify"] = False

    original_urls = m.map("slug")
    redirects = [o.ljust(60) + slugify(o) for o in original_urls if slugify(o) != o]
    assets_dir: str = str(m.config.get("assets_dir", "static"))
    redirects_file = Path(
        str(m.config.get("redirects", Path(assets_dir) / "_redirects"))
    )
    redirects_file.touch()
    with open(redirects_file, "a") as f:
        f.write("\n")
        f.write("# marakta migrate to slugify redirects")
        f.write("\n".join(redirects))
        f.write("\n")
