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

"""
import inspect
from pathlib import Path
from typing import TYPE_CHECKING, List, Union

import jinja2
import pydantic
from jinja2 import Template, Undefined
from more_itertools import flatten

from markata import __version__
from markata.hookspec import hook_impl

env = jinja2.Environment()

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
class StyleOverrides(Style):
    ...


class Meta(pydantic.BaseModel):
    name: str
    content: str


class Text(pydantic.BaseModel):
    value: str


class Link(pydantic.BaseModel):
    rel: str = "canonical"
    href: str


class HeadConfig(pydantic.BaseModel):
    meta: List[Meta] = []
    link: List[Link] = []
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
    post_template: str = None

    @pydantic.validator("post_template", pre=True, always=True)
    def default_post_template(cls, v):
        if v is None:
            return (
                Path(__file__).parent / "default_post_template.html.jinja"
            ).read_text()
        if isinstance(v, Path):
            return v.read_text()
        if isinstance(v, str) and Path(v).exists():
            return Path(v).read_text()
        return v


class PostOverrides(pydantic.BaseModel):
    head: HeadConfig = HeadConfig()
    style: Style = StyleOverrides()


class Post(pydantic.BaseModel):
    config_overrides: PostOverrides = PostOverrides()


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

    # raw_text = "\n".join([t.value for t in markata.config.head.text])

    # if isinstance(raw_text, list):
    #     markata.config["head"]["text"] = "\n".join(
    #         flatten([t.values() for t in raw_text]),
    #     )


@hook_impl
def pre_render(markata: "Markata") -> None:
    """
    FOR EACH POST: Massages the configuration limitations of toml/yaml to allow
    a little bit easier experience to the end user making configurations while
    allowing an simpler jinja template.  This enablees the use of the
    `markata.head.text` list in configuration.
    """
    for article in [a for a in markata.articles if "config_overrides" in a]:
        raw_text = article.get("config_overrides", {}).get("head", {}).get("text", "")

        if isinstance(raw_text, list):
            article["config_overrides"]["head"]["text"] = "\n".join(
                flatten([t.values() for t in raw_text]),
            )


@hook_impl
def render(markata: "Markata") -> None:
    template = Template(markata.config.post_template, undefined=SilentUndefined)

    if "{{" in str(markata.config.get("head", {})):
        Template(
            str(markata.config.get("head", {})),
            undefined=SilentUndefined,
        )
    else:
        pass

    merged_config = markata.config
    for article in [a for a in markata.articles if hasattr(a, "html")]:
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

        article.html = template.render(
            __version__=__version__,
            body=article.html,
            toc=markata.md.toc,  # type: ignore
            config=merged_config,
            post=article,
            **article.metadata,
        )
