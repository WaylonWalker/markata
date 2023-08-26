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
   # creates redirects from static/_redirects file
   "markata.plugins.redirects",
   # copies your static assets into the output_dir (default: `markout`)
   "markata.plugins.copy_assets",
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
redirects need to be taken care of server side.  It also cannot change the http
code, this is only

## Features

The features of markata.plugins.redirect is pretty limited since it is
implemented only as a static page.  Other features require server side
implementation.

| Feature                             | Support | Example                                                         | Notes                                                                                             |
| ----------------------------------- | ------- | --------------------------------------------------------------- | ------------------------------------------------------------------------------------------------- |
| Force                               | Yes     | `/pagethatexists /otherpage`                                    | Creates an index.html with http-equiv and canonical                                               |
| Redirects (301, 302, 303, 307, 308) | No      | `/home / 301`                                                   | Ignored, requires server side implementation                                                      |
| Rewrites (other status codes)       | No      | `/blog/* /blog/404.html 404`                                    | ...                                                                                               |
| Splats                              | No      | `/blog/* /blog/:splat`                                          | ...                                                                                               |
| Placeholders                        | No      | `/blog/:year/:month/:date/:slug /news/:year/:month/:date/:slug` | ...                                                                                               |
| Query Parameters                    | No      | `/shop id=:id /blog/:id 301`                                    | ...                                                                                               |
| Proxying                            | No      | `/blog/* https://blog.my.domain/:splat 200`                     | ...                                                                                               |
| Domain-level redirects              | No      | `workers.example.com/* workers.example.com/blog/:splat 301`     | ...                                                                                               |
| Redirect by country or language     | No      | `/ /us 302 Country=us`                                          | ...                                                                                               |
| Redirect by cookie                  | No      | `/* /preview/:splat 302 Cookie=preview`                        | ...                                                                                               |

> Compare with
> [cloudflare-pages](https://developers.cloudflare.com/pages/platform/redirects/)

!!! tip
    If you have a public site, pair this up with
    [ahrefs](https://app.ahrefs.com/dashboard) to keep up with pages that have
    moved without you realizing.

"""
from pathlib import Path
from typing import Dict, Optional

from jinja2 import Template
import pydantic

from markata import Markata
from markata.hookspec import hook_impl
from pydantic import ConfigDict

DEFAULT_REDIRECT_TEMPLATE = Path(__file__).parent / "default_redirect_template.html"


class Redirect(pydantic.BaseModel):
    "DataClass to store the original and new url"
    original: str
    new: str
    markata: Markata
    model_config = ConfigDict(validate_assignment=True, arbitrary_types_allowed=True)


class RedirectsConfig(pydantic.BaseModel):
    assets_dir: Path = Path("static")
    redirects_file: Optional[Path] = None

    @pydantic.validator("redirects_file", always=True)
    def default_redirects_file(
        cls: "RedirectsConfig", v: Path, *, values: Dict
    ) -> Path:
        if not v:
            return Path(values["assets_dir"]) / "_redirects"
        return v


class Config(pydantic.BaseModel):
    redirects: RedirectsConfig = RedirectsConfig()


@hook_impl(tryfirst=True)
def config_model(markata: "Markata") -> None:
    markata.config_models.append(Config)


@hook_impl
def save(markata: "Markata") -> None:
    """
    saves an index.html in the directory called out by the redirect.
    """
    redirects_file = Path(markata.config.redirects.redirects_file)
    if redirects_file.exists():
        raw_redirects = redirects_file.read_text().split("\n")
    else:
        raw_redirects = []

    redirects = [
        Redirect(original=s[0], new=s[1], markata=markata)
        for r in raw_redirects
        if "*" not in r and len(s := r.split()) == 2 and not r.strip().startswith("#")
    ]

    if "redirect_template" in markata.config:
        template_file = Path(str(markata.config.get("redirect_template")))
    else:
        template_file = DEFAULT_REDIRECT_TEMPLATE
    template = Template(template_file.read_text())

    for redirect in redirects:
        file = markata.config.output_dir / redirect.original.strip("/") / "index.html"
        file.parent.mkdir(parents=True, exist_ok=True)
        file.write_text(template.render(redirect.dict(), config=markata.config))
