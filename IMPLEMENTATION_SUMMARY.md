# Configuration Override Implementation Summary

## Overview

Implemented a comprehensive runtime configuration override system for Markata that allows overriding any configuration value through:
1. Alternate config files (`-c/--config`)
2. Environment variables (`MARKATA_*`)
3. CLI flags (`-o/--output-dir` and `-s/--set`)

## Changes Made

### 1. Core Infrastructure (`markata/__init__.py`)

**Added parameters to `Markata.__init__`:**
- `config_overrides`: Dict of config values to override
- `config_file`: Path to alternate config file

**Stores these for use by plugins:**
- `self._config_overrides`: Available to load_config hook
- `self._config_file`: Available to load_config hook

### 2. Standard Config Loader (`markata/standard_config.py`)

**Enhanced `load()` function:**
- Added `config_file` parameter for explicit config file loading
- Enhanced `_load_env()` to support nested config with double underscore (`MARKATA_STYLE__THEME`)
- Properly handles resolution order: global → local/file → env → overrides

**Environment variable format:**
```bash
MARKATA_OUTPUT_DIR=dist              # Top-level: output_dir
MARKATA_STYLE__THEME=nord            # Nested: style.theme
```

### 3. Config Model Plugin (`markata/plugins/config_model.py`)

**Updated `load_config()` hook:**
- Retrieves `_config_overrides` and `_config_file` from markata instance
- Passes them to `standard_config.load()` for proper override handling

### 4. CLI Plugin (`markata/plugins/base_cli.py`)

**Added new options to `build` command:**
- `-c/--config PATH`: Alternate config file
- `-o/--output-dir TEXT`: Quick output directory override
- `-s/--set TEXT`: Generic config override with dot notation (repeatable)

**New helper functions:**
- `parse_set_options()`: Parses `-s key=value` into nested dicts
- `_deep_merge()`: Deep merges override dicts

**Build command now:**
1. Collects all overrides from CLI flags
2. Creates a new Markata instance with overrides if any provided
3. Replaces the existing instance's state

### 5. Documentation

**Created comprehensive docs:**
- `docs/runtime-config-override.md`: Full usage guide with examples
- Updated `markata/plugins/base_cli.py` docstring with override examples
- Created `example_theme_build.sh`: Practical example for building with multiple themes

## Usage Examples

### 1. Alternate Config File
```bash
markata build -c themes/catppuccin.toml
```

### 2. Environment Variables
```bash
# Simple
MARKATA_OUTPUT_DIR=dist markata build

# Nested
MARKATA_STYLE__THEME=nord markata build

# Multiple
MARKATA_OUTPUT_DIR=dist MARKATA_STYLE__THEME=gruvbox markata build
```

### 3. CLI Overrides
```bash
# Output directory shorthand
markata build -o dist/theme-everforest

# Generic overrides
markata build -s output_dir=dist -s style.theme=nord

# Nested config
markata build -s style.theme=catppuccin -s style.mode=dark

# JSON values
markata build -s 'nav={"home":"/","docs":"/docs"}'
```

### 4. Combined
```bash
# All three methods work together
MARKATA_TITLE="My Site" markata build -c base.toml -s output_dir=dist
```

## Resolution Order

Configuration is loaded and merged in this priority order (later overrides earlier):

1. **Global config files** (`~/.config/markata.toml`)
2. **Local config files** (`./markata.toml` or path from `-c`)
3. **Environment variables** (`MARKATA_*`)
4. **CLI flags** (`-o`, `-s`)

## Primary Use Case: Multi-Theme Builds

Build your site with different themes without modifying config:

```bash
#!/bin/bash
THEMES=("catppuccin" "everforest" "gruvbox" "nord")

for theme in "${THEMES[@]}"; do
    markata build \
        -s "output_dir=markout-$theme" \
        -s "style.theme=$theme"
done
```

This allows you to:
- Test themes without editing config files
- Build multiple theme variants in CI/CD
- Create theme previews automatically
- Maintain a single source config with runtime variations

## Testing

All features tested and verified:
- ✅ Config overrides via Python API
- ✅ Environment variable parsing (simple and nested)
- ✅ CLI flag parsing (`-s` with dot notation)
- ✅ Config file loading (`-c`)
- ✅ Proper resolution order
- ✅ Deep merging of nested configs

## Backward Compatibility

All changes are backward compatible:
- Existing code without overrides continues to work unchanged
- New parameters are optional with sensible defaults
- No breaking changes to existing API

## Files Modified

1. `markata/__init__.py` - Added override parameters
2. `markata/standard_config.py` - Enhanced config loading
3. `markata/plugins/config_model.py` - Updated load_config hook
4. `markata/plugins/base_cli.py` - Added CLI override options

## Files Created

1. `docs/runtime-config-override.md` - Complete documentation
2. `example_theme_build.sh` - Example build script

## Future Enhancements

Potential improvements:
1. Add `--dry-run` to show final config without building
2. Add config validation before building
3. Support loading multiple config files with layering
4. Add `markata config validate` command
5. Support YAML config files in addition to TOML
