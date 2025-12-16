# Runtime Configuration Override

Markata now supports overriding configuration at runtime through multiple methods, allowing you to build your site with different configurations without modifying your config files.

## Override Methods

Configuration can be overridden using three methods (applied in order):

1. **Alternate Config File** (`-c/--config`)
2. **Environment Variables** (`MARKATA_*`)
3. **CLI Flags** (`-s/--set` and `-o/--output-dir`)

Later methods override earlier ones, so CLI flags take precedence over environment variables, which take precedence over config files.

## Usage Examples

### 1. Alternate Config File

Use a different configuration file entirely:

```bash
# Use a theme-specific config
markata build -c themes/catppuccin.toml

# Use a staging config
markata build -c configs/staging.toml
```

### 2. Environment Variables

Override any config value using environment variables with the `MARKATA_` prefix:

```bash
# Simple top-level config
MARKATA_OUTPUT_DIR=dist markata build

# Nested config using double underscore
MARKATA_STYLE__THEME=nord markata build

# Multiple overrides
MARKATA_OUTPUT_DIR=dist MARKATA_STYLE__THEME=gruvbox markata build
```

### 3. CLI Flags

#### Output Directory Override

Quick override for output directory:

```bash
markata build -o dist/theme-everforest
```

#### Generic Config Override (`--set`)

Set any configuration value using dot notation. Config overrides use **deep merge**, so nested values are merged rather than replaced:

```bash
# Single value
markata build -s output_dir=dist

# Nested config (merges with existing style config)
markata build -s style.theme=nord

# Multiple values
markata build -s output_dir=dist -s style.theme=catppuccin

# Override single template while keeping others
markata build -s post_template.index=custom.html

# Add new template variant
markata build -s post_template.summary=summary.html

# Complex values (use JSON)
markata build -s 'nav={"home":"/","docs":"/docs"}'
```

### 4. Combining Methods

All override methods can be combined:

```bash
# Environment + CLI overrides
MARKATA_STYLE__THEME=nord markata build -s output_dir=custom

# Config file + CLI overrides
markata build -c base.toml -s output_dir=dist -s style.theme=gruvbox

# All three methods
MARKATA_TITLE="My Site" markata build -c base.toml -s output_dir=dist
```

## Use Case: Multiple Theme Builds

Build your site with different themes without modifying config files:

```bash
#!/bin/bash
THEMES=("catppuccin" "everforest" "gruvbox" "nord")

for theme in "${THEMES[@]}"; do
    markata build \
        -s "output_dir=markout-$theme" \
        -s "style.theme=$theme"
done
```

See [example_theme_build.sh](../example_theme_build.sh) for a complete example.

## Use Case: Template Overrides

Override or add post templates without modifying config files:

```bash
# Override just the index template (keeps partial, og, etc.)
markata build -s post_template.index=custom.html

# Add a new summary template
markata build -s post_template.summary=summary_template.html

# Override multiple templates
markata build \
    -s post_template.index=custom.html \
    -s post_template.partial=custom_partial.html

# Combine with theme override
markata build \
    -s style.theme=nord \
    -s post_template.index=nord_post.html
```

**Note:** Template overrides use deep merge, so specifying `-s post_template.index=custom.html` will only override the `index` template while preserving other template variants like `partial`, `og`, and `raw.md`.

## Config Value Format

### Simple Values

```bash
# Strings
markata build -s title="My Site"

# Numbers (auto-detected)
markata build -s site_version=2

# Booleans (use JSON)
markata build -s debug=true
```

### Nested Values

Use dot notation for nested configuration:

```bash
# style.theme -> {"style": {"theme": "value"}}
markata build -s style.theme=nord

# nav.home -> {"nav": {"home": "value"}}
markata build -s nav.home="/"
```

### Complex Values

Use JSON for complex structures:

```bash
# Objects
markata build -s 'nav={"home":"/","docs":"/docs"}'

# Arrays
markata build -s 'hooks=["default","markata.plugins.custom"]'
```

## Environment Variable Format

Environment variables follow the pattern `MARKATA_<KEY>` where:

- Top-level keys use single underscore: `MARKATA_OUTPUT_DIR`
- Nested keys use double underscore: `MARKATA_STYLE__THEME`
- Keys are case-insensitive (converted to lowercase)

```bash
# Top-level
export MARKATA_OUTPUT_DIR=dist
export MARKATA_TITLE="My Site"

# Nested (double underscore)
export MARKATA_STYLE__THEME=nord
export MARKATA_NAV__HOME="/"

# Use them
markata build
```

## Resolution Order

Configuration is loaded and merged in this order:

1. Global config files (`~/.config/markata.toml`)
2. Local config files (`./markata.toml` or specified with `-c`)
3. Environment variables (`MARKATA_*`)
4. CLI flags (`-o`, `-s`)

Each layer overrides the previous, so CLI flags have the highest priority.

## Best Practices

1. **Use `-c` for complete config variations** (dev/staging/prod environments)
2. **Use environment variables for deployment-specific values** (URLs, API keys)
3. **Use CLI flags for one-off builds** (testing themes, custom output dirs)

## Migration Guide

If you previously used multiple config files for different environments:

**Before:**
```bash
cp markata-staging.toml markata.toml
markata build
cp markata-prod.toml markata.toml
markata build
```

**After:**
```bash
markata build -c markata-staging.toml
markata build -c markata-prod.toml
```

## Advanced: Scripting Builds

Create a build matrix for testing:

```bash
#!/bin/bash
# Build with all theme variations

THEMES=("catppuccin" "everforest" "gruvbox" "nord" "tokyo-night")
MODES=("light" "dark")

for theme in "${THEMES[@]}"; do
    for mode in "${MODES[@]}"; do
        OUTPUT="dist/${theme}-${mode}"
        echo "Building $theme ($mode) -> $OUTPUT"
        
        markata build \
            -s "output_dir=$OUTPUT" \
            -s "style.theme=$theme" \
            -s "style.mode=$mode"
    done
done
```

## See Also

- [Configuration Documentation](./configuration.md)
- [Plugin Development](./plugin-development.md)
- [Example Theme Build Script](../example_theme_build.sh)
