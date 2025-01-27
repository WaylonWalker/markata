"""


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
from typing import Any, Dict, List, Optional, TYPE_CHECKING, Union

import jinja2
from jinja2 import Template, Undefined
from more_itertools import flatten
import pydantic
from pydantic import ConfigDict, Field, field_validator, model_validator, root_validator
from rich.syntax import Syntax
import typer

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
    name: str
    content: str


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

    @model_validator(mode="after")
    @classmethod
    def dynamic_templates_in_templates_dir(cls, model) -> "Config":
        markata_templates = Path(__file__).parents[1] / "templates"
        templates_dir = model.templates_dir
        dynamic_templates_dir = model.dynamic_templates_dir

        if isinstance(templates_dir, Path):
            model.templates_dir = [
                templates_dir,
                markata_templates,
                dynamic_templates_dir,
            ]

        if markata_templates not in model.templates_dir:
            model.templates_dir.append(markata_templates)

        if dynamic_templates_dir not in model.templates_dir:
            model.templates_dir.append(dynamic_templates_dir)

        return model

    @property
    def jinja_loader(self):
        return jinja2.FileSystemLoader(self.templates_dir)

    @property
    def jinja_env(self):
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


_template_cache = {}


def get_template(markata, template):
    """Get a template from the cache or compile it."""
    try:
        return markata.config.jinja_env.get_template(template)
    except jinja2.TemplateNotFound:
        if template.startswith("/"):
            template_path = Path(template)
        else:
            template_path = (
                Path(__file__).parent.parent / "templates" / template
            ).resolve()

        if not template_path.exists():
            template_path = Path(template)

        if not template_path.exists():
            raise FileNotFoundError(f"Template not found: {template}")

        template_content = template_path.read_text()
        return Template(template_content, undefined=SilentUndefined)


def render_article(markata, cache, article):
    """Render an article using cached templates."""
    key = markata.make_hash(
        "post_template",
        __version__,
        article.key,
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
        for t in markata.config.jinja_env.list_templates()
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
        markata.console.quiet = True
        if template:
            template = get_template(markata, template)

            markata.console.quiet = False
            markata.console.print(template.filename)
            if theme is None or theme.lower() == "none":
                markata.console.print(Path(template.filename).read_text())
            else:
                syntax = Syntax.from_path(template.filename, theme=theme)
                markata.console.print(syntax)

            return
        templates = markata.config.jinja_env.list_templates()
        markata.console.quiet = False
        markata.console.print("Templates directories:", style="green underline")

        markata_templates = Path(__file__).parents[1] / "templates"
        for dir in markata.config.templates_dir:
            if dir == markata.config.dynamic_templates_dir:
                markata.console.print(
                    f"[gold3]{dir}[/][grey50] (dynamically created templates from configuration)[/] [gold3]\[markata.config.dynamic_templates_dir][/]",
                    style="red",
                )
            elif dir == markata_templates:
                markata.console.print(
                    f"[cyan]{dir}[/][grey50] (built-in)[/]", style="red"
                )
            else:
                markata.console.print(
                    f"[orchid]{dir}[/] [orchid]\[markata.config.templates_dir][/]",
                    style="red",
                )

        markata.console.print()
        markata.console.print(
            "Available Templates: [white]name -> path[/]", style="green underline"
        )
        for template in templates:
            source, file, uptodate = markata.config.jinja_env.loader.get_source(
                markata.config.jinja_env, template
            )

            if Path(file).is_relative_to(markata.config.dynamic_templates_dir):
                markata.console.print(
                    f"[gold3]{template} -> [red]{file}[/] [grey50](dynamic)[/]"
                )
            elif Path(file).is_relative_to(markata_templates):
                markata.console.print(
                    f"[cyan]{template} -> [red]{file}[/] [grey50](built-in)[/]"
                )
            else:
                markata.console.print(f"[orchid]{template}[/] -> [red]{file}[/]")


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


@hook_impl
def pre_render(markata: "Markata") -> None:
    """
    FOR EACH POST: Massages the configuration limitations of toml/yaml to allow
    a little bit easier experience to the end user making configurations while
    allowing an simpler jinja template.  This enables the use of the
    `markata.head.text` list in configuration.
    """

    markata.config.dynamic_templates_dir.mkdir(parents=True, exist_ok=True)
    head_template = markata.config.dynamic_templates_dir / "head.html"
    head_template.write_text(
        markata.config.jinja_env.get_template("dynamic_head.html").render(
            {"markata": markata}
        ),
    )

    for article in [a for a in markata.articles if "config_overrides" in a]:
        raw_text = article.get("config_overrides", {}).get("head", {}).get("text", "")

        if isinstance(raw_text, list):
            article["config_overrides"]["head"]["text"] = "\n".join(
                flatten([t.values() for t in raw_text]),
            )


@hook_impl
def render(markata: "Markata") -> None:
    with markata.cache as cache:
        for article in markata.filter("skip==False"):
            html = render_article(markata=markata, cache=cache, article=article)
            article.html = html
