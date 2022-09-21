from pathlib import Path
from typing import TYPE_CHECKING

from jinja2 import Template, Undefined
from more_itertools import flatten

from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from markata import Markata


class SilentUndefined(Undefined):
    def _fail_with_undefined_error(self, *args, **kwargs):
        return ""


@hook_impl
def configure(markata: "Markata") -> None:
    '''
    massages the configuration limitations of toml to allow a little bit easier
    experience to the end user making configurations while allowing an simpler
    jinja template.

    Adding this snippet to your post template will allow users to use the same
    configuration as the default template.

    ``` html
    {{ config.get('head', {}).pop('text') if 'text' in config.get('head',{}).keys() }}{% for tag, meta in config.get('head', {}).items() %}{% for _meta in meta %}
    <{{ tag }} {% for attr, value in _meta.items() %}{{ attr }}="{{ value }}"{% endfor %}/> {% endfor %}{% endfor %}
    ```

    Here is an example config

    ``` toml
    [[markata.head.text]]
    value = """
    some script
    """

    [[markata.head.text]]
    value="""
    some style
    """

    [[markata.head.meta]]
    name = "og:type"
    property = "og:type"
    content = "article"

    [[markata.head.meta]]
    name = "og:author"
    property = "og:author"
    content = "Waylon Walker"
    ```
    '''

    markata.config["head"]["text"] = "\n".join(
        flatten([t.values() for t in markata.config["head"]["text"]])
    )


@hook_impl
def render(markata: "Markata") -> None:
    if "post_template" in markata.config:
        template_file = markata.config["post_template"]
    else:
        template_file = Path(__file__).parent / "default_post_template.html"
    with open(template_file) as f:
        template = Template(f.read(), undefined=SilentUndefined)
    for article in markata.iter_articles("apply template"):

        article.html = template.render(
            body=article.html,
            toc=markata.md.toc,  # type: ignore
            config=markata.config,
            **article,
        )
