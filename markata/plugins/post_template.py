"""


The `markata.plugins.post_template` plugin handles the rendering of posts using Jinja2
templates. It provides extensive configuration options for HTML head elements, styling,
and template customization.

## Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.post_template",
]
```

## Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.post_template",
]
```

## Configuration

## Head Elements

Configure HTML head elements in `markata.toml`:

```toml
# Meta tags
[[markata.head.meta]]
name = "og:type"
content = "article"

[[markata.head.meta]]
name = "og:author"
content = "Your Name"

# Links
[[markata.head.link]]
rel = "canonical"
href = "https://example.com"

# Scripts
[[markata.head.script]]
src = "/assets/main.js"

# Raw HTML
markata.head.text = '''
<style>
  /* Custom CSS */
</style>
'''
```

## Styling

Configure default styles:

```toml
[markata.style]
color_bg = "#1f2022"
color_text = "#eefbfe"
color_link = "#fb30c4"
color_accent = "#e1bd00c9"
body_width = "800px"
```

## Templates

Configure template settings:

```toml
[markata]
# Default template
post_template = "post.html"

# Template directories
templates_dir = "templates"
dynamic_templates_dir = ".markata.cache/templates"
template_cache_dir = ".markata.cache/template_bytecode"

# Jinja environment options
env_options = { trim_blocks = true }
```

## Functionality

## Template Rendering

The plugin:
1. Loads templates from configured directories
2. Compiles and caches templates for performance
3. Renders posts with Jinja2 templating
4. Supports template inheritance and includes
5. Provides template bytecode caching

## Post-Specific Overrides

Each post can override global settings:

```yaml
---
template: custom.html
config_overrides:
  head:
    meta:
      - name: og:type
        content: video
  style:
    color_bg: "#000000"
---
```

## Template Context

Templates have access to:
- Post attributes
- Global configuration
- Custom filters and functions
- Markata instance

## Performance Features

- Template bytecode caching
- Template compilation caching
- Configurable Jinja2 environment
- Efficient head element rendering

## Dependencies

This plugin depends on:
- jinja2 for templating
- pydantic for configuration
- typer for CLI commands

# Add head configuration

This snippet allows users to configure their head in `markata.toml`.

``` html
{{ config.get('head', {}).pop('text') if 'text' in config.get('head',{}).keys() }}
{% for tag, meta in config.get('head', {}).items() %}
    {% for _meta in meta %}
        <{{ tag }}
            {% for attr, value in _meta.items() %}{{ attr }}="{{ value }}"{% endfor %}
        />
    {% endfor %}
{% endfor %}
```

Users can specify any sort of tag in their `markata.toml`

``` toml
[[markata.head.meta]]
name = "og:type"
content = "article"

[[markata.head.meta]]
name = "og:author"
content = "Waylon Walker"
```

The above configuration becomes this once rendered.

``` html
<meta name='og:type' content='article' />
<meta name='og:Author' content='Waylon Walker' />
```

!! Note

    Article variables can be used for dynamic entries like canonical_url
    ``` toml
    [markata]
    url = "markata.dev"

    [[markata.head.meta]]
    href="{{ config.url }}/{{ slug }}/"
    rel="canonical"
    ```

Optionally users can also specify plain text to be appended to the head of
their documents.  This works well for things that involve full blocks.

``` toml
[[markata.head.text]]
value = '''
<script>
    console.log('hello world')
</script>
'''

[[markata.head.text]]
value='''
html  {
    font-family: "Space Mono", monospace;
    background: var(--color-bg);
    color: var(--color-text);
}
'''
```

## Add scripts to head

Markata config also supports adding scripts to the head via configuration.

``` toml
[[ markata.head.script ]]
    src = "https://cdn.tailwindcss.com"

```

"""

import inspect
from pathlib import Path
from typing import TYPE_CHECKING
from typing import Any
from typing import Dict
from typing import List
from typing import Optional
from typing import Union

import jinja2
import pydantic
import typer
from jinja2 import Undefined
from more_itertools import flatten
from pydantic import ConfigDict
from pydantic import Field
from pydantic import field_validator
from pydantic import model_validator

from markata import __version__
from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from markata import Markata


class SilentUndefined(Undefined):
    def _fail_with_undefined_error(self, *args, **kwargs):
        return ""


def optional(*fields):
    def dec(_cls):
        for field in fields:
            _cls.__fields__[field].default = None
        return _cls

    if (
        fields
        and inspect.isclass(fields[0])
        and issubclass(fields[0], pydantic.BaseModel)
    ):
        cls = fields[0]
        fields = cls.__fields__
        return dec(cls)
    return dec


class Style(pydantic.BaseModel):
    color_bg: str = "#1f2022"
    color_bg_code: str = "#1f2022"
    color_text: str = "#eefbfe"
    color_link: str = "#fb30c4"
    color_accent: str = "#e1bd00c9"
    overlay_brightness: str = ".85"
    body_width: str = "800px"
    color_bg_light: str = "#eefbfe"
    color_bg_code_light: str = "#eefbfe"
    color_text_light: str = "#1f2022"
    color_link_light: str = "#fb30c4"
    color_accent_light: str = "#ffeb00"
    overlay_brightness_light: str = ".95"


@optional
class StyleOverrides(Style): ...


class Meta(pydantic.BaseModel):
    name: Optional[str] = None
    property: Optional[str] = None
    content: str

    @field_validator("name", mode="after")
    def check_og(cls, v):
        if v.startswith("og:"):
            raise ValueError("Meta names cannot start with og:, use property instead")
        return v

    # ensure one of name or property is set
    @model_validator(mode="after")
    def check_one(cls, values):
        if not values.name and not values.property:
            raise ValueError("One of name or property must be set")
        return values


class Text(pydantic.BaseModel):
    value: str


class Link(pydantic.BaseModel):
    rel: str = "canonical"
    href: str


class Script(pydantic.BaseModel):
    src: str


class HeadConfig(pydantic.BaseModel):
    meta: List[Meta] = []
    link: List[Link] = []
    script: List[Script] = []
    text: Union[List[Text], str] = ""

    @field_validator("text", mode="before")
    def text_to_list(cls, v):
        if isinstance(v, list):
            return "\n".join([text["value"] for text in v])
        return v

    @property
    def html(self):
        html = self.text
        html += "\n"
        for meta in self.meta:
            html += f'<meta name="{meta.name}" content="{meta.content}" />\n'
        for link in self.link:
            html += f'<link rel="{link.rel}" href="{link.href}" />\n'
        return html


class Config(pydantic.BaseModel):
    head: HeadConfig = HeadConfig()
    style: Style = Style()
    post_template: Optional[Union[str | Dict[str, str]]] = "post.html"
    dynamic_templates_dir: Path = Path(".markata.cache/templates")
    templates_dir: Union[Path, List[Path]] = pydantic.Field(Path("templates"))
    template_cache_dir: Path = Path(".markata.cache/template_bytecode")
    env_options: dict = {}

    @pydantic.validator("templates_dir", pre=True)
    def dynamic_templates_in_templates_dir(cls, value):
        """Ensure all required template directories are included in templates_dir.

        Args:
            value: The input templates_dir value, can be Path, str, or List[Path]

        Returns:
            List[Path]: List of template directories including markata templates and dynamic templates
        """
        markata_templates = Path(__file__).parents[1] / "templates"
        dynamic_templates_dir = Path(".markata.cache/templates")

        # Convert string to Path if needed
        if isinstance(value, str):
            value = Path(value)

        # Convert single Path to list
        if isinstance(value, Path):
            value = [value]

        # Ensure it's a list of Paths
        templates_dir = [Path(p) if isinstance(p, str) else p for p in value]

        # Add required directories if not present
        if markata_templates not in templates_dir:
            templates_dir.append(markata_templates)

        if dynamic_templates_dir not in templates_dir:
            templates_dir.append(dynamic_templates_dir)

        return templates_dir


_template_cache = {}


def get_template(markata, template):
    """Get a template from the cache or compile it."""
    cache_key = str(template)
    if cache_key in _template_cache:
        return _template_cache[cache_key]

    if isinstance(template, str):
        template = markata.jinja_env.get_template(template)
    _template_cache[cache_key] = template
    return template


def get_templates_mtime(markata):
    """Get latest mtime from all template directories.

    This tracks changes to any template file including includes, extends, and imports.
    """
    max_mtime = 0
    for template_dir in markata.jinja_env.template_paths:
        template_path = Path(template_dir)
        if template_path.exists():
            for path in template_path.rglob('*'):
                if path.is_file():
                    try:
                        max_mtime = max(max_mtime, path.stat().st_mtime)
                    except (OSError, FileNotFoundError):
                        continue
    return max_mtime


def render_article(markata, cache, article):
    """Render an article using cached templates."""
    templates_mtime = get_templates_mtime(markata)

    key = markata.make_hash(
        "post_template",
        __version__,
        article.key,
        str(templates_mtime),  # Track template file changes
    )
    html = markata.precache.get(key)

    if html is not None:
        return html

    if isinstance(article.template, str):
        template = get_template(markata, article.template)
        html = render_template(markata, article, template)

    if isinstance(article.template, dict):
        html = {
            slug: render_template(markata, article, get_template(markata, template))
            for slug, template in article.template.items()
        }
    cache.set(key, html, expire=markata.config.default_cache_expire)
    return html


def render_template(markata, article, template):
    """Render a template with article context."""
    merged_config = markata.config

    # Get the body content - prefer article_html, fallback to html
    body = getattr(article, "article_html", None)
    if body is None:
        body = getattr(article, "html", "")

    context = {
        "post": article,
        "markata": markata,
        "config": merged_config,
        "body": body,
    }

    try:
        return template.render(**context)
    except Exception as e:
        markata.console.print(f"[red]Error rendering template for {article.path}[/]")
        markata.console.print(f"[red]{str(e)}[/]")
        raise


@hook_impl()
def save(markata: "Markata") -> None:
    linked_templates = [
        t
        for t in markata.jinja_env.list_templates()
        if t.endswith("css") or t.endswith("js") or t.endswith("xsl")
    ]
    for template in linked_templates:
        template = get_template(markata, template)
        css = template.render(markata=markata, __version__=__version__)
        Path(markata.config.output_dir / Path(template.filename).name).write_text(css)


@hook_impl()
def cli(app: typer.Typer, markata: "Markata") -> None:
    """
    Markata hook to implement base cli commands.
    """

    templates_app = typer.Typer()
    app.add_typer(templates_app, name="templates")

    @templates_app.callback()
    def templates():
        "template management"

    @templates_app.command()
    def show(
        template: str = typer.Argument(None, help="template to show"),
        theme: str = typer.Option(None, help="pygments syntax theme"),
    ) -> None:
        markata.console.quiet = False
        if template is not None:
            # Show specific template
            try:
                syntax = markata.jinja_env.get_template(template).source
                markata.console.print(syntax)
            except Exception as e:
                markata.console.print(
                    f"Error loading template {template}: {str(e)}", style="red"
                )

            return

        try:
            templates = markata.jinja_env.list_templates()
            markata.console.quiet = False
            markata.console.print("Templates directories:", style="green underline")

            # Show built-in templates directory
            markata_templates = Path(__file__).parents[1] / "templates"
            # markata.console.print(f"  {markata_templates} [grey50](built-in)[/]")

            # Show user template paths
            for path in markata.config.templates_dir:
                if path == markata_templates:
                    markata.console.print(f"  {path} [grey50](built-in)[/]")
                elif path == markata.config.dynamic_templates_dir:
                    markata.console.print(f"  {path} [grey50](dynamic)[/]")
                else:
                    markata.console.print(f"  {path}")

            markata.console.print("\nAvailable templates:", style="green underline")
            for template in sorted(templates):
                try:
                    source, file, uptodate = markata.jinja_env.loader.get_source(
                        markata.jinja_env, template
                    )
                    if Path(file).is_relative_to(markata.config.dynamic_templates_dir):
                        markata.console.print(
                            f"  {template} -> {file} [grey50](dynamic)[/]"
                        )
                    elif Path(file).is_relative_to(markata_templates):
                        markata.console.print(
                            f"  {template} -> {file} [grey50](built-in)[/]"
                        )
                    else:
                        markata.console.print(f"  {template} -> {file}")
                except Exception:
                    markata.console.print(f"  {template}")
        except Exception as e:
            markata.console.print(f"Error listing templates: {str(e)}", style="red")


class PostOverrides(pydantic.BaseModel):
    head: HeadConfig = HeadConfig()
    style: Style = StyleOverrides()


class Post(pydantic.BaseModel):
    config_overrides: PostOverrides = PostOverrides()
    template: Optional[str | Dict[str, str]] = None
    markata: Any = Field(None, exclude=True)

    model_config = ConfigDict(
        validate_assignment=True,  # Config model
        arbitrary_types_allowed=True,
        extra="allow",
        str_strip_whitespace=True,
        validate_default=True,
        coerce_numbers_to_str=True,
        populate_by_name=True,
    )

    @field_validator("template", mode="before")
    @classmethod
    def default_template(cls, v, info) -> Union[str, Dict[str, str]]:
        markata = info.data.get("markata")
        if v is None:
            return markata.config.post_template
        if isinstance(v, str):
            v = {"index": v}
        if isinstance(markata.config.post_template, str):
            config_template = {
                "index": markata.config.post_template,
            }
        else:
            config_template = markata.config.post_template
        return {**config_template, **v}


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


@hook_impl(tryfirst=True)
def config_model(markata: "Markata") -> None:
    markata.config_models.append(Config)


@hook_impl(tryfirst=True)
def post_model(markata: "Markata") -> None:
    markata.post_models.append(Post)


@hook_impl
def configure(markata: "Markata") -> None:
    """
    Massages the configuration limitations of toml to allow a little bit easier
    experience to the end user making configurations while allowing an simpler
    jinja template.  This enablees the use of the `markata.head.text` list in
    configuration.
    """


@hook_impl(tryfirst=True)
def pre_render(markata: "Markata") -> None:
    """
    FOR EACH POST: Massages the configuration limitations of toml/yaml to allow
    a little bit easier experience to the end user making configurations while
    allowing an simpler jinja template.  This enables the use of the
    `markata.head.text` list in configuration.
    """

    # markata.config.dynamic_templates_dir.mkdir(parents=True, exist_ok=True)
    # head_template = markata.config.dynamic_templates_dir / "head.html"
    # head_template.write_text(
    #     markata.jinja_env.get_template("dynamic_head.html").render(
    #         {"markata": markata}
    #     ),
    # )

    for article in [a for a in markata.articles if "config_overrides" in a]:
        raw_text = article.get("config_overrides", {}).get("head", {}).get("text", "")

        if isinstance(raw_text, list):
            article["config_overrides"]["head"]["text"] = "\n".join(
                flatten([t.values() for t in raw_text]),
            )


@hook_impl
def render(markata: "Markata") -> None:
    with markata.cache as cache:
        for article in markata.filter("not skip"):
            html = render_article(markata=markata, cache=cache, article=article)
            article.html = html
