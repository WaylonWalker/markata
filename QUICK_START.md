# Quick Start: Multi-Theme Builds

## TL;DR

```bash
# Build all themes with gallery
./build-all-themes.sh

# View locally
cd markout && python -m http.server 8000
# Open http://localhost:8000
```

## What You Get

```
markout/
├── index.html          # Beautiful theme gallery
├── tokyo-night/        # Each theme gets its own directory
├── catppuccin/
├── everforest/
├── gruvbox/
├── kanagwa/
└── nord/
```

## One-Liners

```bash
# Build single theme
markata build -s style.theme=catppuccin -s output_dir=markout/catppuccin

# Build with custom output
markata build -s style.theme=nord -o dist/nord-theme

# Build from different config
markata build -c themes/nord.toml -o markout/nord

# Use environment variable
MARKATA_STYLE__THEME=gruvbox markata build
```

## Gallery Preview

The generated gallery includes:
- 🎨 Color-coded cards for each theme
- 📱 Responsive design
- 🔗 Direct links to each build
- 📊 Build statistics
- 💡 Usage instructions

See `BUILD_SCRIPTS_README.md` for full documentation.
