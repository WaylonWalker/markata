"""
Creates redirects for times when your backend server can't.

## Configuration

Enable the redirect hook by adding it to your list of hooks.

``` toml
[markata]

# Were you keep static assets to copy into the project, default is static
# the assets_dir will set the default _redirects file directory
assets_dir = "static" 

# You can override the default redirects file location
redirects = static/_redirects 

hooks = [
   "markata.plugins.redirects", # creates redirects from static/_redirects file
   "markata.plugins.copy_assets", # copies your static assets into the output_dir (default: `markout`)
  ...
]
```

## Syntax

Your `_redirects` file is a simplified version of what services like cloudflare
pages or netlify use.  In fact you can use the same redirects file!

Here is an example that will redirect `/old` to `/new` and `/CHANGELOG` to
`/changelog`

```
/old        /new
/CHANGELOG  /changelog
```

## Limitations
_no splats_

Since it is static generated this plugin cannot cover *'s.  * or splat
redirects need to be taken care of server side.

!! tip
    If you have a public site, pair this up with
    [ahrefs](https://app.ahrefs.com/dashboard) to keep up with pages that have
    moved without you realizing.

"""
from dataclasses import asdict, dataclass
from pathlib import Path
from typing import TYPE_CHECKING

from jinja2 import Template

from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from markata import Markata

DEFAULT_REDIRECT_TEMPLATE = Path(__file__).parent / "default_redirect_template.html"


@dataclass
class Redirect:
    original: str
    new: str


@hook_impl
def save(markata: "Markata") -> None:
    """
    saves an index.html in the directory called out by the redirect.
    """
    assets_dir: str = str(markata.config.get("assets_dir", "static"))
    redirects_file = Path(
        str(markata.config.get("redirects", Path(assets_dir) / "_redirects"))
    )
    if redirects_file.exists():
        raw_redirects = redirects_file.read_text().split("\n")
    else:
        raw_redirects = []

    redirects = [
        Redirect(*s)
        for r in 
        if "*" not in r and len(s := r.split()) == 2 and not r.strip().startswith("#")
    ]

    if "redirect_template" in markata.config:
        template_file = Path(str(markata.config.get("redirect_template")))
    else:
        template_file = DEFAULT_REDIRECT_TEMPLATE
    template = Template(template_file.read_text())

    output_dir = Path(markata.config["output_dir"])  # type: ignore
    output_dir.mkdir(parents=True, exist_ok=True)

    for redirect in redirects:
        file = output_dir / redirect.original.strip("/") / "index.html"
        file.parent.mkdir(parents=True, exist_ok=True)
        file.write_text(template.render(**asdict(redirect), config=markata.config))
