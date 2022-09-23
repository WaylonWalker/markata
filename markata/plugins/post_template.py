"""


## Add head configuration

This snippet allows users to configure their head in `markata.toml`.

``` html
{{ config.get('head', {}).pop('text') if 'text' in config.get('head',{}).keys() }}{% for tag, meta in config.get('head', {}).items() %}{% for _meta in meta %}
<{{ tag }} {% for attr, value in _meta.items() %}{{ attr }}="{{ value }}"{% endfor %}/> {% endfor %}{% endfor %}
```

Users can specify any sort of tag in their `markata.toml`

```
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
    """
    Massages the configuration limitations of toml to allow a little bit easier
    experience to the end user making configurations while allowing an simpler
    jinja template.  This enablees the use of the `markata.head.text` list in
    configuration.
    """

    raw_text = markata.config.get("head", {}).get("text", "")

    if isinstance(raw_text, list):
        markata.config["head"]["text"] = "\n".join(
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
    for article in markata.iter_articles("apply template"):

        article.html = template.render(
            body=article.html,
            toc=markata.md.toc,  # type: ignore
            config=markata.config,
            **article,
        )
