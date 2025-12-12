#!/bin/bash
#
# Build Markata site with all available themes
# Creates separate builds in markout/<theme> and generates an index
#

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Available themes
THEMES=(
    "tokyo-night"
    "catppuccin"
    "everforest"
    "gruvbox"
    "kanagwa"
    "nord"
)

# Base output directory
OUTPUT_BASE="markout"

echo -e "${BLUE}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   Building Markata Site with Multiple Themes          ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════╝${NC}"
echo ""

# Clean old builds
if [ -d "$OUTPUT_BASE" ]; then
    echo -e "${YELLOW}→ Cleaning previous builds in $OUTPUT_BASE${NC}"
    rm -rf "$OUTPUT_BASE"
fi

mkdir -p "$OUTPUT_BASE"

# Track build results
declare -a SUCCESSFUL_BUILDS
declare -a FAILED_BUILDS

# Build each theme
for theme in "${THEMES[@]}"; do
    echo ""
    echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
    echo -e "${BLUE}Building theme: ${GREEN}$theme${NC}"
    echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
    
    OUTPUT_DIR="$OUTPUT_BASE/$theme"
    
    # Clear cache before each build to ensure theme changes take effect
    rm -rf .markata.cache
    
    # Run the build with theme override
    # Note: Suppress stderr to hide harmless Rich library recursion errors during teardown
    if uv run markata build \
        -s "output_dir=$OUTPUT_DIR" \
        -s "style.theme=$theme" \
        --quiet 2>/dev/null; then
        
        SUCCESSFUL_BUILDS+=("$theme")
        echo -e "${GREEN}✓ Successfully built $theme theme${NC}"
        echo -e "  Output: $OUTPUT_DIR"
        
        # Count generated files
        if [ -d "$OUTPUT_DIR" ]; then
            FILE_COUNT=$(find "$OUTPUT_DIR" -type f | wc -l)
            echo -e "  Files: $FILE_COUNT"
        fi
    else
        FAILED_BUILDS+=("$theme")
        echo -e "${RED}✗ Failed to build $theme theme${NC}"
    fi
done

echo ""
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}Generating Theme Index${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"

# Generate index.html
cat > "$OUTPUT_BASE/index.html" << 'EOF'
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Markata Theme Gallery</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 2rem;
            color: #333;
        }
        
        .container {
            max-width: 1400px;
            margin: 0 auto;
        }
        
        header {
            text-align: center;
            margin-bottom: 4rem;
            color: white;
        }
        
        h1 {
            font-size: 3.5rem;
            font-weight: 800;
            margin-bottom: 1rem;
            text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
        }
        
        .subtitle {
            font-size: 1.3rem;
            opacity: 0.9;
            font-weight: 300;
        }
        
        .stats {
            display: flex;
            justify-content: center;
            gap: 2rem;
            margin-top: 1.5rem;
            flex-wrap: wrap;
        }
        
        .stat {
            background: rgba(255,255,255,0.2);
            padding: 0.75rem 1.5rem;
            border-radius: 2rem;
            backdrop-filter: blur(10px);
        }
        
        .stat-number {
            font-size: 2rem;
            font-weight: bold;
            display: block;
        }
        
        .stat-label {
            font-size: 0.9rem;
            opacity: 0.9;
        }
        
        .theme-grid {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
            gap: 2rem;
            margin-bottom: 3rem;
        }
        
        .theme-card {
            background: white;
            border-radius: 1rem;
            overflow: hidden;
            box-shadow: 0 10px 30px rgba(0,0,0,0.2);
            transition: all 0.3s ease;
            text-decoration: none;
            color: inherit;
            display: block;
        }
        
        .theme-card:hover {
            transform: translateY(-5px);
            box-shadow: 0 15px 40px rgba(0,0,0,0.3);
        }
        
        .theme-preview {
            width: 100%;
            height: 200px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            position: relative;
            overflow: hidden;
        }
        
        .theme-preview::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><text y="50" font-size="50" fill="white" fill-opacity="0.1">Aa</text></svg>');
            background-size: 50px;
            opacity: 0.3;
        }
        
        .theme-card[data-theme="tokyo-night"] .theme-preview {
            background: linear-gradient(135deg, #1a1b26 0%, #24283b 100%);
        }
        
        .theme-card[data-theme="catppuccin"] .theme-preview {
            background: linear-gradient(135deg, #f5c2e7 0%, #cba6f7 100%);
        }
        
        .theme-card[data-theme="everforest"] .theme-preview {
            background: linear-gradient(135deg, #a7c080 0%, #7fbbb3 100%);
        }
        
        .theme-card[data-theme="gruvbox"] .theme-preview {
            background: linear-gradient(135deg, #d79921 0%, #cc241d 100%);
        }
        
        .theme-card[data-theme="kanagwa"] .theme-preview {
            background: linear-gradient(135deg, #7e9cd8 0%, #938aa9 100%);
        }
        
        .theme-card[data-theme="nord"] .theme-preview {
            background: linear-gradient(135deg, #88c0d0 0%, #5e81ac 100%);
        }
        
        .theme-content {
            padding: 1.5rem;
        }
        
        .theme-name {
            font-size: 1.5rem;
            font-weight: 700;
            margin-bottom: 0.5rem;
            text-transform: capitalize;
        }
        
        .theme-description {
            color: #666;
            line-height: 1.6;
            margin-bottom: 1rem;
        }
        
        .theme-meta {
            display: flex;
            gap: 1rem;
            font-size: 0.9rem;
            color: #888;
        }
        
        .theme-badge {
            display: inline-flex;
            align-items: center;
            gap: 0.3rem;
        }
        
        .badge-icon {
            width: 16px;
            height: 16px;
        }
        
        footer {
            text-align: center;
            color: white;
            margin-top: 4rem;
            padding-top: 2rem;
            border-top: 1px solid rgba(255,255,255,0.2);
        }
        
        footer a {
            color: white;
            text-decoration: underline;
        }
        
        .info-box {
            background: rgba(255,255,255,0.95);
            border-radius: 1rem;
            padding: 2rem;
            margin-bottom: 3rem;
            box-shadow: 0 10px 30px rgba(0,0,0,0.2);
        }
        
        .info-box h2 {
            margin-bottom: 1rem;
            color: #667eea;
        }
        
        .info-box code {
            background: #f5f5f5;
            padding: 0.2rem 0.5rem;
            border-radius: 0.3rem;
            font-family: 'Monaco', 'Menlo', monospace;
            font-size: 0.9rem;
        }
        
        @media (max-width: 768px) {
            h1 {
                font-size: 2.5rem;
            }
            
            .theme-grid {
                grid-template-columns: 1fr;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>🎨 Markata Theme Gallery</h1>
            <p class="subtitle">Explore all available themes for your Markata site</p>
            <div class="stats">
                <div class="stat">
                    <span class="stat-number" id="theme-count">0</span>
                    <span class="stat-label">Themes</span>
                </div>
                <div class="stat">
                    <span class="stat-number" id="build-date">Today</span>
                    <span class="stat-label">Built</span>
                </div>
            </div>
        </header>
        
        <div class="info-box">
            <h2>📋 How to Use These Themes</h2>
            <p>Each theme can be activated in your <code>markata.toml</code> configuration:</p>
            <br>
            <pre><code>[markata.style]
theme = "catppuccin"  # or any theme name below</code></pre>
            <br>
            <p>Or override at build time:</p>
            <pre><code>markata build -s style.theme=nord</code></pre>
        </div>
        
        <div class="theme-grid" id="theme-grid">
            <!-- Theme cards will be inserted here -->
        </div>
        
        <footer>
            <p>Built with ❤️ using <a href="https://markata.dev" target="_blank">Markata</a></p>
            <p style="margin-top: 0.5rem; opacity: 0.8;">
                View the <a href="https://github.com/WaylonWalker/markata" target="_blank">source code on GitHub</a>
            </p>
        </footer>
    </div>
    
    <script>
        // Theme data
        const themes = [
THEME_DATA_PLACEHOLDER
        ];
        
        // Theme descriptions
        const descriptions = {
            'tokyo-night': 'A dark theme inspired by Tokyo\'s vibrant nightlife, featuring cool blues and purples.',
            'catppuccin': 'A soothing pastel theme with soft pinks and purples, perfect for long reading sessions.',
            'everforest': 'A comfortable green theme inspired by forests, easy on the eyes for extended use.',
            'gruvbox': 'A retro groove theme with warm, earthy colors and high contrast.',
            'kanagwa': 'A theme inspired by the famous Japanese wave painting, featuring blues and teals.',
            'nord': 'An arctic, north-bluish color palette with cool blues and crisp contrast.',
        };
        
        // Update stats
        document.getElementById('theme-count').textContent = themes.length;
        document.getElementById('build-date').textContent = new Date().toLocaleDateString('en-US', {
            month: 'short',
            day: 'numeric',
            year: 'numeric'
        });
        
        // Generate theme cards
        const grid = document.getElementById('theme-grid');
        themes.forEach(theme => {
            const card = document.createElement('a');
            card.className = 'theme-card';
            card.setAttribute('data-theme', theme);
            // Link to docs/index.html which is the main landing page
            card.href = `./${theme}/docs/index.html`;
            
            const themeName = theme.split('-').map(w => 
                w.charAt(0).toUpperCase() + w.slice(1)
            ).join(' ');
            
            card.innerHTML = `
                <div class="theme-preview">
                    <svg style="position: absolute; bottom: 1rem; right: 1rem; opacity: 0.5;" width="60" height="60" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2">
                        <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"/>
                    </svg>
                </div>
                <div class="theme-content">
                    <div class="theme-name">${themeName}</div>
                    <div class="theme-description">${descriptions[theme] || 'A beautiful Markata theme.'}</div>
                    <div class="theme-meta">
                        <span class="theme-badge">
                            <svg class="badge-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <circle cx="12" cy="12" r="10"/>
                                <path d="M12 6v6l4 2"/>
                            </svg>
                            View Demo
                        </span>
                        <span class="theme-badge">
                            <svg class="badge-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                                <polyline points="7 10 12 15 17 10"/>
                                <line x1="12" y1="15" x2="12" y2="3"/>
                            </svg>
                            Ready
                        </span>
                    </div>
                </div>
            `;
            
            grid.appendChild(card);
        });
    </script>
</body>
</html>
EOF

# Replace placeholder with actual theme data
THEME_JS_ARRAY=""
for theme in "${SUCCESSFUL_BUILDS[@]}"; do
    THEME_JS_ARRAY="${THEME_JS_ARRAY}            '${theme}',\n"
done

# Update the index.html with actual themes
sed -i "s|THEME_DATA_PLACEHOLDER|${THEME_JS_ARRAY}|g" "$OUTPUT_BASE/index.html"

echo -e "${GREEN}✓ Generated theme index at $OUTPUT_BASE/index.html${NC}"

# Print summary
echo ""
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}Build Summary${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo ""
echo -e "${GREEN}Successful builds: ${#SUCCESSFUL_BUILDS[@]}${NC}"
for theme in "${SUCCESSFUL_BUILDS[@]}"; do
    echo -e "  ${GREEN}✓${NC} $theme"
done

if [ ${#FAILED_BUILDS[@]} -gt 0 ]; then
    echo ""
    echo -e "${RED}Failed builds: ${#FAILED_BUILDS[@]}${NC}"
    for theme in "${FAILED_BUILDS[@]}"; do
        echo -e "  ${RED}✗${NC} $theme"
    done
fi

echo ""
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}All Done!${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo ""
echo -e "View the theme gallery: ${YELLOW}$OUTPUT_BASE/index.html${NC}"
echo ""
echo -e "To serve locally, run:"
echo -e "  ${YELLOW}cd $OUTPUT_BASE && python -m http.server 8000${NC}"
echo ""
