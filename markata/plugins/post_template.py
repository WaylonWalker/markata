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
import copy
from pathlib import Path
from typing import TYPE_CHECKING, List, Optional

import jinja2
import pydantic
from deepmerge import always_merger
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


class Style(pydantic.BaseModel):
    color_bg: str = '#1f2022'
    color_bg_code: str = '#1f2022'
    color_text: str = '#eefbfe'
    color_link: str = '#fb30c4'
    color_accent: str = '#e1bd00c9'
    overlay_brightness: str = '.85'
    body_width: str = '800px'
    color_bg_light: str = '#eefbfe'
    color_bg_code_light: str = '#eefbfe'
    color_text_light: str = '#1f2022'
    color_link_light: str = '#fb30c4'
    color_accent_light: str = '#ffeb00'
    overlay_brightness_light: str = '.95'


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
    text: List[Text] = []
    link: List[Link] = []


class Config(pydantic.BaseModel):
    head: HeadConfig = HeadConfig()
    style: Style = Style()


@hook_impl(tryfirst=True)
def config_model(markata: "Markata") -> None:
    markata.config_models.append(Config)


@hook_impl
def configure(markata: "Markata") -> None:
    """
    Massages the configuration limitations of toml to allow a little bit easier
    experience to the end user making configurations while allowing an simpler
    jinja template.  This enablees the use of the `markata.head.text` list in
    configuration.
    """

    raw_text = "\n".join([t.value for t in markata.config.head.text])

    if isinstance(raw_text, list):
        markata.config["head"]["text"] = "\n".join(
            flatten([t.values() for t in raw_text])
        )


@hook_impl
def pre_render(markata: "Markata") -> None:
    """
    FOR EACH POST: Massages the configuration limitations of toml/yaml to allow
    a little bit easier experience to the end user making configurations while
    allowing an simpler jinja template.  This enablees the use of the
    `markata.head.text` list in configuration.
    """
    for article in [a for a in markata.articles if "config_overrides" in a.keys()]:
        raw_text = article.get("config_overrides", {}).get("head", {}).get("text", "")

        if isinstance(raw_text, list):
            article["config_overrides"]["head"]["text"] = "\n".join(
                flatten([t.values() for t in raw_text])
            )


@hook_impl
def render(markata: "Markata") -> None:
    if "post_template" in markata.config:
        template_file = markata.config["post_template"]
    else:
        template_file = Path(__file__).parent / "default_post_template.html"
    with open(template_file) as f:
        template = Template(f.read(), undefined=SilentUndefined)

    if "{{" in str(markata.config.get("head", {})):
        head_template = Template(
            str(markata.config.get("head", {})), undefined=SilentUndefined
        )
    else:
        head_template = None
        head = {}

    # _full_config = copy.deepcopy(markata.config)

    merged_config = markata.config
    for article in [a for a in markata.articles if hasattr(a, "html")]:
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

        markata.console.log(f'rendering {article.title}')
        article.html = template.render(
            __version__=__version__,
            body=article.html,
            toc=markata.md.toc,  # type: ignore
            config=merged_config,
            **article.metadata,
        )
        markata.console.log(f'rendered {article.title}')
