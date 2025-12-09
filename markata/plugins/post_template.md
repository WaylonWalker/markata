---
date: 2025-12-09
description: "The plugin handles the rendering of posts using Jinja2 templates. It
  provides extensive configuration options for HTML head elements, styling, and template\u2026"
published: false
slug: markata/plugins/post-template
title: post_template.py


---

---

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

---

!!! function
    <h2 id="render_article" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">render_article <em class="small">function</em></h2>

    Render an article using cached templates.

???+ source "render_article <em class='small'>source</em>"
    ```python
    def render_article(markata, cache, article):
        """Render an article using cached templates."""
        templates_mtime = get_templates_mtime(markata.jinja_env)

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
            template = get_template(markata.jinja_env, article.template)
            html = render_template(markata, article, template)

        if isinstance(article.template, dict):
            html = {
                slug: render_template(markata, article, get_template(markata.jinja_env, template))
                for slug, template in article.template.items()
            }
        cache.set(key, html, expire=markata.config.default_cache_expire)
        return html
    ```
!!! function
    <h2 id="render_template" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">render_template <em class="small">function</em></h2>

    Render a template with article context.

???+ source "render_template <em class='small'>source</em>"
    ```python
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
    ```
!!! function
    <h2 id="cli" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">cli <em class="small">function</em></h2>

    Markata hook to implement base cli commands.

???+ source "cli <em class='small'>source</em>"
    ```python
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
    ```
!!! class
    <h2 id="MarkataTemplateCache" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">MarkataTemplateCache <em class="small">class</em></h2>

    Template bytecode cache for improved performance.

???+ source "MarkataTemplateCache <em class='small'>source</em>"
    ```python
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
    ```
!!! function
    <h2 id="configure" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">configure <em class="small">function</em></h2>

    Massages the configuration limitations of toml to allow a little bit easier
    experience to the end user making configurations while allowing an simpler
    jinja template.  This enablees the use of the `markata.head.text` list in
    configuration.

???+ source "configure <em class='small'>source</em>"
    ```python
    def configure(markata: "Markata") -> None:
        """
        Massages the configuration limitations of toml to allow a little bit easier
        experience to the end user making configurations while allowing an simpler
        jinja template.  This enablees the use of the `markata.head.text` list in
        configuration.
        """
    ```
!!! function
    <h2 id="pre_render" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">pre_render <em class="small">function</em></h2>

    FOR EACH POST: Massages the configuration limitations of toml/yaml to allow
    a little bit easier experience to the end user making configurations while
    allowing an simpler jinja template.  This enables the use of the
    `markata.head.text` list in configuration.

???+ source "pre_render <em class='small'>source</em>"
    ```python
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
    ```
!!! method
    <h2 id="dynamic_templates_in_templates_dir" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">dynamic_templates_in_templates_dir <em class="small">method</em></h2>

    Ensure all required template directories are included in templates_dir.

    Args:
        value: The input templates_dir value, can be Path, str, or List[Path]

    Returns:
        List[Path]: List of template directories including markata templates and dynamic templates

???+ source "dynamic_templates_in_templates_dir <em class='small'>source</em>"
    ```python
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
    ```
!!! function
    <h2 id="templates" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">templates <em class="small">function</em></h2>

    template management

???+ source "templates <em class='small'>source</em>"
    ```python
    def templates():
            "template management"
    ```