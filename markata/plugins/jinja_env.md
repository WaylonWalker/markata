---
date: 2025-05-24
description: "Jinja2 Environment Plugin Provides a centralized Jinja2 environment
  configuration for consistent template rendering across all Markata plugins. This
  plugin\u2026"
published: false
slug: markata/plugins/jinja-env
title: jinja_env.py


---

---

Jinja2 Environment Plugin

Provides a centralized Jinja2 environment configuration for consistent template rendering
across all Markata plugins. This plugin ensures template rendering behavior is consistent
and available even when specific template-using plugins are not enabled.

## Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

``` toml
hooks = [
    "markata.plugins.jinja_env",
]
```

## Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

``` toml
disabled_hooks = [
    "markata.plugins.jinja_env",
]
```

## Configuration

Configure Jinja environment settings in your `markata.toml`:

``` toml
[markata.jinja_env]
template_paths = ["templates"]  # Additional template paths to search
undefined_silent = true        # Return empty string for undefined variables
trim_blocks = true            # Remove first newline after block
lstrip_blocks = true          # Strip tabs/spaces from start of line
template_cache_dir = ".markata.cache/template_bytecode"
```

# Usage

The environment is automatically available to other plugins via `markata.config.jinja_env`.
Template loading follows this order:
1. Package templates (built-in Markata templates)
2. User template paths (configured via template_paths)

Example usage in a plugin:

``` python
def render_template(markata, content):
    template = markata.jinja_env.from_string(content)
    return template.render(markata=markata)
```

# Notes

- Template paths are resolved relative to the current working directory
- Package templates are always available and take precedence
- Silent undefined behavior means undefined variables render as empty strings

---

!!! class
    <h2 id="_SilentUndefined" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">_SilentUndefined <em class="small">class</em></h2>

    Custom undefined type that returns empty string for undefined variables.

???+ source "_SilentUndefined <em class='small'>source</em>"
    ```python
    class _SilentUndefined(jinja2.Undefined):
        """Custom undefined type that returns empty string for undefined variables."""

        def _fail_with_undefined_error(self, *args, **kwargs):
            return ""

        __add__ = __radd__ = __mul__ = __rmul__ = __div__ = __rdiv__ = __truediv__ = (
            __rtruediv__
        ) = __floordiv__ = __rfloordiv__ = __mod__ = __rmod__ = __pos__ = __neg__ = (
            __call__
        ) = __getitem__ = __lt__ = __le__ = __gt__ = __ge__ = __int__ = __float__ = (
            __complex__
        ) = __pow__ = __rpow__ = _fail_with_undefined_error
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
!!! class
    <h2 id="JinjaEnvConfig" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">JinjaEnvConfig <em class="small">class</em></h2>

    Configuration for the Jinja environment.

???+ source "JinjaEnvConfig <em class='small'>source</em>"
    ```python
    class JinjaEnvConfig(pydantic.BaseModel):
        """Configuration for the Jinja environment."""

        templates_dir: List[str] = []
        undefined_silent: bool = True
        trim_blocks: bool = True
        lstrip_blocks: bool = True
        template_cache_dir: Path = Path(".markata.cache/template_bytecode")

        model_config = pydantic.ConfigDict(
            validate_assignment=True,  # Config model
            arbitrary_types_allowed=True,
            extra="allow",
            str_strip_whitespace=True,
            validate_default=True,
            coerce_numbers_to_str=True,
            populate_by_name=True,
        )
    ```
!!! function
    <h2 id="config_model" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">config_model <em class="small">function</em></h2>

    Register configuration models.

???+ source "config_model <em class='small'>source</em>"
    ```python
    def config_model(markata: "Markata") -> None:
        """Register configuration models."""
        markata.config_models.append(JinjaEnvConfig)
    ```
!!! function
    <h2 id="configure" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">configure <em class="small">function</em></h2>

    Initialize and configure the Jinja2 environment for Markata.

    This hook runs early in the configuration stage to ensure the jinja environment
    is available for other plugins that need it during configuration.

    Args:
        markata: The Markata instance

???+ source "configure <em class='small'>source</em>"
    ```python
    def configure(markata: Markata) -> None:
        """Initialize and configure the Jinja2 environment for Markata.

        This hook runs early in the configuration stage to ensure the jinja environment
        is available for other plugins that need it during configuration.

        Args:
            markata: The Markata instance
        """
        # Get configuration, falling back to defaults
        config = JinjaEnvConfig()
        if hasattr(markata.config, "jinja_env"):
            if isinstance(markata.config.jinja_env, dict):
                config = JinjaEnvConfig(**markata.config.jinja_env)

        # TODO: setting up env twice could not get dynamic templates to be recognized on first pass
        loaders = []
        if markata.config.templates_dir:
            for path in markata.config.templates_dir:
                path = Path(path).expanduser().resolve()
                if path.exists():
                    loaders.append(FileSystemLoader(str(path)))
        # Create environment
        env_for_dynamic_render = Environment(
            loader=ChoiceLoader(loaders),
            undefined=_SilentUndefined if config.undefined_silent else jinja2.Undefined,
            trim_blocks=config.trim_blocks,
            lstrip_blocks=config.lstrip_blocks,
            bytecode_cache=MarkataTemplateCache(config.template_cache_dir),
            auto_reload=True,
        )

        markata.config.dynamic_templates_dir.mkdir(parents=True, exist_ok=True)
        head_template = markata.config.dynamic_templates_dir / "head.html"
        head_template.write_text(
            env_for_dynamic_render.get_template("dynamic_head.html").render(
                {"markata": markata}
            ),
        )

        # Set up loaders
        loaders = []

        # Add package templates first (lowest priority)
        # loaders.append(PackageLoader("markata", "templates"))

        # Add user template paths (medium priority)
        if markata.config.templates_dir:
            for path in markata.config.templates_dir:
                path = Path(path).expanduser().resolve()
                if path.exists():
                    loaders.append(FileSystemLoader(str(path)))

        # Add dynamic templates directory (highest priority)
        # dynamic_templates_dir = Path(".markata.cache/templates")
        # dynamic_templates_dir.mkdir(parents=True, exist_ok=True)
        # loaders.append(FileSystemLoader(str(dynamic_templates_dir)))

        # Create environment
        env = Environment(
            loader=ChoiceLoader(loaders),
            undefined=_SilentUndefined if config.undefined_silent else jinja2.Undefined,
            trim_blocks=config.trim_blocks,
            lstrip_blocks=config.lstrip_blocks,
            bytecode_cache=MarkataTemplateCache(config.template_cache_dir),
            auto_reload=True,
        )

        # Register the environment on the config's private attribute
        markata.jinja_env = env
    ```