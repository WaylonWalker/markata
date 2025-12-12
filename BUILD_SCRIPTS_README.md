# Multi-Theme Build Scripts

This directory contains scripts for building your Markata site with multiple themes.

## Scripts

### `build-all-themes.sh`

A comprehensive script that builds your site with every available theme and generates a beautiful gallery index page.

**Features:**
- ✨ Builds all 6 themes (tokyo-night, catppuccin, everforest, gruvbox, kanagwa, nord)
- 📁 Outputs each theme to `markout/<theme>/`
- 🎨 Generates a beautiful HTML gallery index at `markout/index.html`
- 📊 Provides build statistics and error reporting
- 🎯 Color-coded terminal output

**Usage:**

```bash
./build-all-themes.sh
```

**Output structure:**
```
markout/
├── index.html          # Theme gallery homepage
├── tokyo-night/        # Tokyo Night theme build
│   └── index.html
├── catppuccin/         # Catppuccin theme build
│   └── index.html
├── everforest/         # Everforest theme build
│   └── index.html
├── gruvbox/            # Gruvbox theme build
│   └── index.html
├── kanagwa/            # Kanagwa theme build
│   └── index.html
└── nord/               # Nord theme build
    └── index.html
```

**To view locally:**

```bash
cd markout
python -m http.server 8000
# Open http://localhost:8000 in your browser
```

### `example_theme_build.sh`

A simpler example script showing the basic pattern for building with multiple themes.

**Usage:**

```bash
./example_theme_build.sh
```

## Gallery Features

The generated `markout/index.html` includes:

- 🎨 **Beautiful gradient background** with glassmorphism effects
- 🖼️ **Theme preview cards** with unique colors for each theme
- 📱 **Responsive design** that works on all devices
- 📝 **Theme descriptions** explaining each color scheme
- 🔗 **Direct links** to each theme's build
- 📊 **Build statistics** (theme count, build date)
- 💡 **Usage instructions** showing how to activate themes

## Customization

### Adding Custom Themes

If you've created custom themes, add them to the `THEMES` array in `build-all-themes.sh`:

```bash
THEMES=(
    "tokyo-night"
    "catppuccin"
    "everforest"
    "gruvbox"
    "kanagwa"
    "nord"
    "your-custom-theme"  # Add your theme here
)
```

### Changing Output Directory

Modify the `OUTPUT_BASE` variable:

```bash
OUTPUT_BASE="dist"  # Change from "markout" to "dist"
```

### Customizing Theme Descriptions

Edit the `descriptions` object in the generated HTML:

```javascript
const descriptions = {
    'your-theme': 'Your custom description here',
};
```

## CI/CD Integration

### GitHub Actions

```yaml
name: Build All Themes

on:
  push:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Setup Python
        uses: actions/setup-python@v4
        with:
          python-version: '3.11'
      - name: Install dependencies
        run: pip install -e .
      - name: Build all themes
        run: ./build-all-themes.sh
      - name: Deploy to GitHub Pages
        uses: peaceiris/actions-gh-pages@v3
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./markout
```

### GitLab CI

```yaml
build-themes:
  stage: build
  script:
    - pip install -e .
    - ./build-all-themes.sh
  artifacts:
    paths:
      - markout/
```

## Troubleshooting

### Build Failures

If a theme fails to build, the script will:
1. Continue with remaining themes
2. Report failed themes in the summary
3. Exit with status code 0 (for partial success)

Check the build output for specific error messages.

### Permission Errors

Make sure the script is executable:

```bash
chmod +x build-all-themes.sh
```

### Missing Themes

If themes aren't showing up:
1. Check that the theme exists in `markata/plugins/theme.py`
2. Verify the theme name matches the `THEMES` array
3. Check console output for build errors

## Advanced Usage

### Build Specific Themes Only

Edit the `THEMES` array to include only the themes you want:

```bash
THEMES=(
    "catppuccin"
    "nord"
)
```

### Parallel Builds

For faster builds, you can use `xargs` for parallel execution:

```bash
echo "tokyo-night catppuccin everforest" | xargs -n 1 -P 3 -I {} bash -c '
    markata build -s output_dir=markout/{} -s style.theme={} --quiet
'
```

### Custom Build Options

Add additional overrides to each build:

```bash
markata build \
    -s "output_dir=$OUTPUT_DIR" \
    -s "style.theme=$theme" \
    -s "url=https://example.com/$theme" \
    -s "title=My Site - $theme" \
    --quiet
```

## See Also

- [Runtime Configuration Override Documentation](docs/runtime-config-override.md)
- [Theme Documentation](docs/themes.md)
- [Markata Documentation](https://markata.dev)
