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

"""

import os
from pathlib import Path
from typing import Dict, List, Union

import anyconfig

path_spec_type = List


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
        {"path_specs": home / f"{tool}.ini", "ac_parser": "ini", "keys": [tool]},
        {"path_specs": home / f".{tool}", "ac_parser": "ini", "keys": [tool]},
        {"path_specs": home / f".{tool}.ini", "ac_parser": "ini", "keys": [tool]},
        {
            "path_specs": home / ".config" / f"{tool}.ini",
            "ac_parser": "ini",
            "keys": [tool],
        },
        {
            "path_specs": home / ".config" / f".{tool}",
            "ac_parser": "ini",
            "keys": [tool],
        },
        {
            "path_specs": home / ".config" / f".{tool}.ini",
            "ac_parser": "ini",
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
            "ac_parser": "ini",
            "keys": [tool],
        },
        {
            "path_specs": Path(project_home) / f".{tool}",
            "ac_parser": "ini",
            "keys": [tool],
        },
        {
            "path_specs": Path(project_home) / f".{tool}.ini",
            "ac_parser": "ini",
            "keys": [tool],
        },
        {
            "path_specs": Path(project_home) / f"{tool}.yml",
            "ac_parser": "yaml",
            "keys": [tool],
        },
        {
            "path_specs": Path(project_home) / f".{tool}.yml",
            "ac_parser": "yaml",
            "keys": [tool],
        },
        {
            "path_specs": Path(project_home) / f"{tool}.toml",
            "ac_parser": "toml",
            "keys": [tool],
        },
        {
            "path_specs": Path(project_home) / f".{tool}.toml",
            "ac_parser": "toml",
            "keys": [tool],
        },
        {
            "path_specs": Path(project_home) / "pyproject.toml",
            "ac_parser": "toml",
            "keys": ["tool", tool],
        },
        {
            "path_specs": Path(project_home) / "setup.cfg",
            "ac_parser": "ini",
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


def _load_files(config_path_specs: path_spec_type) -> Dict:
    """Use anyconfig to load config files stopping at the first one that exists.

    config_path_specs (list): a list of pathspecs and keys to load
    """
    for file in config_path_specs:
        if file["path_specs"].exists():
            config = anyconfig.load(**file)
        else:
            # ignore missing files
            continue

        try:
            return _get_attrs(file["keys"], config)
        except KeyError:
            # ignore incorrect keys
            continue

    return {}


def _load_env(tool: str) -> Dict:
    """Load config from environment variables.

    Args:
        tool (str): name of the tool to configure
    """
    vars = [var for var in os.environ if var.startswith(tool.upper())]
    return {
        var.lower().strip(tool.lower()).strip("_").strip("-"): os.environ[var]
        for var in vars
    }


def load(tool: str, project_home: Union[Path, str] = ".", overrides: Dict = {}) -> Dict:
    """Load tool config from standard config files.

    Resolution Order

    * First global file with a tool key
    * First local file with a tool key
    * Environment variables prefixed with `TOOL`
    * Overrides

    Args:
        tool (str): name of the tool to configure
    """
    global_config = _load_files(_get_global_path_specs(tool))
    local_config = _load_files(_get_local_path_specs(tool, project_home))
    env_config = _load_env(tool)
    return {**global_config, **local_config, **env_config, **overrides}
