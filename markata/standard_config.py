"""Standard Config.
A module to load tooling config from a users project space.

Inspired from frustrations that some tools have a tool.ini, .tool.ini,
setup.cfg, or pyproject.toml.  Some allow for global configs, some don't.  Some
properly follow the users home directory, others end up in a weird temp
directory.  Windows home directory is only more confusing.  Some will even
respect the users `$XDG_HOME` directory.

This file is for any project that can be configured in plain text such as `ini`
or `toml` and not requiring a .py file.  Just name your tool and let users put
config where it makes sense to them, no need to figure out resolution order.

## Usage:

``` python
from standard_config import load

# Retrieve any overrides from the user
overrides = {'setting': True}
config = load('my_tool', overrides)
```

## Resolution Order

* First global file with a tool key
* First local file with a tool key
* Environment variables prefixed with `TOOL`
* Overrides

### Tool Specific Ini files

Ini file formats must include a `<tool>` key.

``` ini
[my_tool]
setting = True
```

### pyproject.toml

Toml files must include a `tool.<tool>` key

``` toml
[tool.my_tool]
setting = True
```

### setup.cfg

setup.cfg files must include a `tool:<tool>` key

``` ini
[tool:my_tool]
setting = True
```

### global files to consider

* <home>/tool.ini
* <home>/.tool
* <home>/.tool.ini
* <home>/.config/tool.ini
* <home>/.config/.tool
* <home>/.config/.tool.ini

### local files to consider

* <project_home>/tool.ini
* <project_home>/.tool
* <project_home>/.tool.ini
* <project_home>/pyproject.toml
* <project_home>/setup.cfg

Markata's standard configuration system.

## Configuration Overview

Markata uses a hierarchical configuration system based on Pydantic models. Configuration
can be set through:
1. TOML files
2. Environment variables
3. Command line arguments

# Basic Configuration

Minimal `markata.toml`:
```toml
[markata]
# Site info
title = "My Site"
url = "https://example.com"
description = "Site description"

# Content locations
content_dir = "content"
output_dir = "markout"
assets_dir = "static"

# Plugin management
hooks = ["default"]
```

# Environment Variables

All settings can be overridden with environment variables:
```bash
# Override site URL
export MARKATA_URL="https://staging.example.com"

# Override output directory
export MARKATA_OUTPUT_DIR="dist"

# Enable debug mode
export MARKATA_DEBUG=1
```

# Detailed Configuration

## Core Settings

```toml
[markata]
# Site information
title = "My Site"                  # Site title
url = "https://example.com"        # Base URL
description = "Site description"   # Meta description
author_name = "Author Name"        # Author name
author_email = "me@example.com"    # Author email
icon = "favicon.ico"               # Site icon
lang = "en"                        # Site language

# Content locations
content_dir = "content"           # Source content location
output_dir = "markout"            # Build output location
assets_dir = "static"             # Static assets location
template_dir = "templates"        # Template location

# Plugin management
hooks = ["default"]               # Active plugins
disabled_hooks = []               # Disabled plugins
```

## Cache Settings

```toml
[markata]
# Cache configuration
default_cache_expire = 3600       # Default TTL (1 hour)
template_cache_expire = 86400     # Template TTL (24 hours)
markdown_cache_expire = 21600     # Markdown TTL (6 hours)
dynamic_cache_expire = 3600       # Dynamic TTL (1 hour)
```

## Development Settings

```toml
[markata]
# Development server
dev_server_port = 8000            # Local server port
dev_server_host = "localhost"     # Local server host
debug = false                     # Debug mode

# Performance
parallel = true                   # Enable parallel processing
workers = 4                       # Number of worker threads
```

## Content Settings

```toml
[markata]
# Content processing
default_template = "post.html"    # Default template
markdown_extensions = [           # Markdown extensions
    "fenced_code",
    "tables",
    "footnotes"
]

# Content filtering
draft = false                     # Include drafts
future = false                    # Include future posts
```

# Plugin Configuration

Each plugin can define its own configuration section:

```toml
# RSS feed configuration
[markata.feeds]
rss = { output = "rss.xml" }
atom = { output = "atom.xml" }
json = { output = "feed.json" }

# Template configuration
[markata.template]
engine = "jinja2"
cache_size = 100
autoescape = true

# Markdown configuration
[markata.markdown]
highlight_theme = "monokai"
line_numbers = true
```

## Configuration Validation

The configuration is validated using Pydantic models:

```python
from pydantic import BaseModel, Field

class MarkataConfig(BaseModel):
    \"\"\"Core configuration model.\"\"\"
    # Site info
    title: str = Field(..., description="Site title")
    url: str = Field(..., description="Site base URL")

    # Directories
    content_dir: Path = Field("content", description="Content directory")
    output_dir: Path = Field("markout", description="Output directory")

    # Features
    debug: bool = Field(False, description="Enable debug mode")
    parallel: bool = Field(True, description="Enable parallel processing")

    model_config = ConfigDict(
        validate_assignment=True,
        arbitrary_types_allowed=True,
        extra="allow",
        str_strip_whitespace=True,
        validate_default=True,
        populate_by_name=True,
    )
```

# Usage Example

```python
from markata import Markata

# Load config from file
markata = Markata.from_file("markata.toml")

# Access configuration
print(markata.config.title)         # Site title
print(markata.config.url)           # Site URL
print(markata.config.content_dir)   # Content directory

# Access plugin config
print(markata.config.feeds.rss)     # RSS feed config
print(markata.config.template)      # Template config

# Override config
markata.config.debug = True
markata.config.parallel = False
```

See hookspec.py for plugin development and lifecycle.py for build process details.
"""

import configparser
import os
from pathlib import Path
from typing import Any
from typing import Dict
from typing import List
from typing import Optional
from typing import Union

import tomli
import yaml

path_spec_type = List[Dict[str, Union[Path, str, List[str]]]]


def _get_global_path_specs(tool: str) -> path_spec_type:
    """
    Generate a list of standard pathspecs for global config files.

    Args:
        tool (str): name of the tool to configure
    """
    try:
        home = Path(os.environ["XDG_HOME"])
    except KeyError:
        home = Path.home()

    return [
        {"path_specs": home / f"{tool}.ini", "parser": "ini", "keys": [tool]},
        {"path_specs": home / f".{tool}", "parser": "ini", "keys": [tool]},
        {"path_specs": home / f".{tool}.ini", "parser": "ini", "keys": [tool]},
        {
            "path_specs": home / ".config" / f"{tool}.ini",
            "parser": "ini",
            "keys": [tool],
        },
        {
            "path_specs": home / ".config" / f".{tool}",
            "parser": "ini",
            "keys": [tool],
        },
        {
            "path_specs": home / ".config" / f".{tool}.ini",
            "parser": "ini",
            "keys": [tool],
        },
    ]


def _get_local_path_specs(tool: str, project_home: Union[str, Path]) -> path_spec_type:
    """
    Generate a list of standard pathspecs for local, project directory config files.

    Args:
        tool (str): name of the tool to configure
    """
    return [
        {
            "path_specs": Path(project_home) / f"{tool}.ini",
            "parser": "ini",
            "keys": [tool],
        },
        {
            "path_specs": Path(project_home) / f".{tool}",
            "parser": "ini",
            "keys": [tool],
        },
        {
            "path_specs": Path(project_home) / f".{tool}.ini",
            "parser": "ini",
            "keys": [tool],
        },
        {
            "path_specs": Path(project_home) / f"{tool}.yml",
            "parser": "yaml",
            "keys": [tool],
        },
        {
            "path_specs": Path(project_home) / f".{tool}.yml",
            "parser": "yaml",
            "keys": [tool],
        },
        {
            "path_specs": Path(project_home) / f"{tool}.toml",
            "parser": "toml",
            "keys": [tool],
        },
        {
            "path_specs": Path(project_home) / f".{tool}.toml",
            "parser": "toml",
            "keys": [tool],
        },
        {
            "path_specs": Path(project_home) / "pyproject.toml",
            "parser": "toml",
            "keys": ["tool", tool],
        },
        {
            "path_specs": Path(project_home) / "setup.cfg",
            "parser": "ini",
            "keys": [f"tool.{tool}"],
        },
    ]


def _get_attrs(attrs: list, config: Dict) -> Dict:
    """Get nested config data from a list of keys.

    specifically written for pyproject.toml which needs to get `tool` then `<tool>`
    """
    for attr in attrs:
        config = config[attr]
    return config


def _load_config_file(file_spec: Dict[str, Any]) -> Optional[Dict[str, Any]]:
    """Load a configuration file using the appropriate parser.

    Args:
        file_spec: Dictionary containing path_specs, parser, and keys information

    Returns:
        Optional[Dict[str, Any]]: Parsed configuration or None if file doesn't exist
    """
    path = file_spec["path_specs"]
    if not path.exists():
        return None

    try:
        if file_spec["parser"] == "toml":
            with open(path, "rb") as f:
                config = tomli.load(f)
        elif file_spec["parser"] == "yaml":
            with open(path, "r") as f:
                config = yaml.safe_load(f)
        elif file_spec["parser"] == "ini":
            config = configparser.ConfigParser()
            config.read(path)
            # Convert ConfigParser to dict
            config = {s: dict(config.items(s)) for s in config.sections()}
        else:
            return None

        return _get_attrs(file_spec["keys"], config)
    except (
        KeyError,
        TypeError,
        yaml.YAMLError,
        tomli.TOMLDecodeError,
        configparser.Error,
    ) as e:
        # warn if var tool name in file name
        if file_spec["keys"][0] in str(path):
            message = f"Failed to load config file: {path}: {e}"
            raise UserWarning(message)

        return None


def _load_files(config_path_specs: path_spec_type) -> Dict[str, Any]:
    """Load config files stopping at the first one that exists and can be parsed.

    Args:
        config_path_specs: List of path specifications to try

    Returns:
        Dict[str, Any]: Configuration dictionary
    """
    for file_spec in config_path_specs:
        config = _load_config_file(file_spec)
        if config:
            return config
    return {}


def _load_env(tool: str) -> Dict[str, Any]:
    """Load config from environment variables.
    
    Supports nested configuration with double underscore:
    MARKATA_STYLE__THEME=nord -> {"style": {"theme": "nord"}}

    Args:
        tool (str): name of the tool to configure
    """
    env_prefix = tool.upper()
    env_config = {}

    for key, value in os.environ.items():
        if key.startswith(f"{env_prefix}_"):
            # Remove prefix
            config_key = key.replace(f"{env_prefix}_", "").lower()

            # Handle nested keys with double underscore
            if "__" in config_key:
                keys = config_key.split("__")
                current = env_config
                for k in keys[:-1]:
                    if k not in current:
                        current[k] = {}
                    current = current[k]
                current[keys[-1]] = value
            else:
                env_config[config_key] = value

    return env_config


def load(
    tool: str,
    project_home: Union[Path, str] = ".",
    overrides: Optional[Dict[str, Any]] = None,
    config_file: Optional[Union[Path, str]] = None,
) -> Dict[str, Any]:
    """Load tool config from standard config files.

    Resolution Order

    * First global file with a tool key
    * First local file with a tool key (or specific config_file if provided)
    * Environment variables prefixed with `TOOL`
    * Overrides

    Args:
        tool (str): name of the tool to configure
        project_home (Union[Path, str], optional): Project directory to search for config files. Defaults to ".".
        overrides (Dict, optional): Override values to apply last. Defaults to None.
        config_file (Union[Path, str], optional): Specific config file to load instead of searching. Defaults to None.

    Returns:
        Dict[str, Any]: Configuration object
    """
    overrides = overrides or {}
    config = {}

    # Load from files in order of precedence
    config.update(_load_files(_get_global_path_specs(tool)) or {})

    # If a specific config file is provided, use it instead of searching
    if config_file:
        config_path = Path(config_file)
        if not config_path.exists():
            raise FileNotFoundError(f"Config file not found: {config_file}")

        # Determine parser from file extension
        suffix = config_path.suffix.lower()
        if suffix == ".toml":
            parser = "toml"
        elif suffix in (".yml", ".yaml"):
            parser = "yaml"
        elif suffix in (".ini", ".cfg"):
            parser = "ini"
        else:
            # Try toml as default
            parser = "toml"

        file_spec = {
            "path_specs": config_path,
            "parser": parser,
            "keys": [tool] if parser == "ini" else (["tool", tool] if parser == "toml" else [tool]),
        }

        file_config = _load_config_file(file_spec)
        if file_config:
            config.update(file_config)
    else:
        config.update(_load_files(_get_local_path_specs(tool, project_home)) or {})

    config.update(_load_env(tool))
    config.update(overrides)

    # If no settings class is provided, return the raw dict
    return config
