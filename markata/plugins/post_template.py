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
from functools import lru_cache
from pathlib import Path
from typing import TYPE_CHECKING, Dict, List, Optional, Union

import jinja2
import pydantic
import typer
from jinja2 import Template, Undefined
from more_itertools import flatten
from rich.syntax import Syntax

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

    @pydantic.validator("text", pre=True)
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

    env_options: dict = {}

    @pydantic.model_validator(mode="after")
    def dynamic_templates_in_templates_dir(self):
        markata_templates = Path(__file__).parents[1] / "templates"

        if isinstance(self.templates_dir, Path):
            self.templates_dir = [
                self.templates_dir,
                markata_templates,
                self.dynamic_templates_dir,
            ]

        if markata_templates not in self.templates_dir:
            self.templates_dir.append(markata_templates)

        if self.dynamic_templates_dir not in self.templates_dir:
            self.templates_dir.append(self.dynamic_templates_dir)

        return self

    @property
    def jinja_loader(self):
        return jinja2.FileSystemLoader(self.templates_dir)

    @property
    def jinja_env(
        self,
    ):
        if hasattr(self, "_jinja_env"):
            return self._jinja_env
        self.env_options.setdefault("loader", self.jinja_loader)
        self.env_options.setdefault("undefined", SilentUndefined)
        self.env_options.setdefault("lstrip_blocks", True)
        self.env_options.setdefault("trim_blocks", True)

        env = jinja2.Environment(**self.env_options)

        self._jinja_env = env
        return env


class PostOverrides(pydantic.BaseModel):
    head: HeadConfig = HeadConfig()
    style: Style = StyleOverrides()


class Post(pydantic.BaseModel):
    config_overrides: PostOverrides = PostOverrides()
    template: Optional[str | Dict[str, str]] = None

    @pydantic.validator("template", pre=True, always=True)
    def default_template(cls, v, *, values):
        if v is None:
            return values["markata"].config.post_template
        if isinstance(v, str):
            v = {"index": v}
        if isinstance(values["markata"].config.post_template, str):
            config_template = {
                "index": values["markata"].config.post_template,
            }
        else:
            config_template = values["markata"].config.post_template
        return {**config_template, **v}


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
        for article in markata.articles:
            html = render_article(markata=markata, cache=cache, article=article)
            article.html = html


@lru_cache()
def get_template(markata, template):
    try:
        return markata.config.jinja_env.get_template(template)
    except jinja2.TemplateNotFound:
        # try to load it as a file
        ...

    try:
        return Template(Path(template).read_text(), undefined=SilentUndefined)
    except FileNotFoundError:
        # default to load it as a string
        ...
    return Template(template, undefined=SilentUndefined)


# def render_article(markata, article):
def render_article(markata, cache, article):
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
    cache.add(key, html, expire=markata.config.default_cache_expire)
    return html


def render_template(markata, article, template):
    template = get_template(markata, template)
    merged_config = markata.config
    # TODO do we need to handle merge??
    # if head_template:
    #     head = eval(
    #         head_template.render(
    #             __version__=__version__,
    #             config=_full_config,
    #             **article,
    #         )
    #     )

    # merged_config = {
    #     **_full_config,
    #     **{"head": head},
    # }

    # merged_config = always_merger.merge(
    #     merged_config,
    #     copy.deepcopy(
    #         article.get(
    #             "config_overrides",
    #             {},
    #         )
    #     ),
    # )

    html = template.render(
        __version__=__version__,
        markata=markata,
        body=article.article_html,
        config=merged_config,
        post=article,
    )
    return html


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
    app.add_typer(templates_app)

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
