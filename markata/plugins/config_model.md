---
date: 2025-12-09
description: "The plugin defines Markata's core configuration model, providing validation
  and type safety for all configuration options. Installation This plugin is\u2026"
published: false
slug: markata/plugins/config-model
title: config_model.py


---

---

The `markata.plugins.config_model` plugin defines Markata's core configuration model,
providing validation and type safety for all configuration options.

## Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.config_model",
]
```

## Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.config_model",
]
```

Note: Disabling this plugin will break most of Markata's functionality as the Config
model is fundamental to the system.

## Configuration

Configure Markata in `markata.toml`:

```toml
[markata]
# Core settings
output_dir = "markout"
assets_dir = "static"

# Plugin management
hooks = ["default"]
disabled_hooks = []

# Cache settings
default_cache_expire = 3600
template_cache_expire = 86400  # 24 hours
markdown_cache_expire = 21600  # 6 hours
dynamic_cache_expire = 3600   # 1 hour

# Markdown settings
markdown_extensions = []

# Development settings
dev_server_port = 8000
dev_server_host = "localhost"
```

## Functionality

### Configuration Model

Core settings:
- `output_dir`: Build output location
- `assets_dir`: Static assets location
- `hooks`: Active plugins
- `disabled_hooks`: Disabled plugins
- `markdown_extensions`: Markdown processors
- Cache expiration times
- Development server settings

## Validation

The model provides:
- Type checking and coercion
- Path validation
- URL validation
- Color validation
- Integer constraints
- Default values

## Settings Management

Features:
- Environment variable support
- TOML file loading
- Settings inheritance
- Dynamic updates
- Validation on change

## Performance

Uses optimized Pydantic config:
- Assignment validation
- Arbitrary types
- Extra fields
- String stripping
- Default validation
- Number coercion
- Name population

## Dependencies

This plugin depends on:
- pydantic for model definition
- pydantic-settings for settings management
- pydantic-extra-types for color support
- rich for console output

---

!!! method
    <h2 id="__getitem__" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">__getitem__ <em class="small">method</em></h2>

    for backwards compatability

???+ source "__getitem__ <em class='small'>source</em>"
    ```python
    def __getitem__(self, item):
            "for backwards compatability"
            return getattr(self, item)
    ```
!!! method
    <h2 id="__setitem__" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">__setitem__ <em class="small">method</em></h2>

    for backwards compatability

???+ source "__setitem__ <em class='small'>source</em>"
    ```python
    def __setitem__(self, key, item):
            "for backwards compatability"
            return setattr(self, key, item)
    ```
!!! method
    <h2 id="get" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">get <em class="small">method</em></h2>

    for backwards compatability

???+ source "get <em class='small'>source</em>"
    ```python
    def get(self, item, default):
            "for backwards compatability"
            return getattr(self, item, default)
    ```
!!! method
    <h2 id="keys" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">keys <em class="small">method</em></h2>

    for backwards compatability

???+ source "keys <em class='small'>source</em>"
    ```python
    def keys(self):
            "for backwards compatability"
            return self.__dict__.keys()
    ```