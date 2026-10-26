---
title: "Themes and Aesthetics"
description: "Configure color palettes and visual aesthetics for your site's appearance"
date: 2024-01-15
published: true
slug: /docs/guides/themes/
tags:
  - documentation
  - themes
  - styling
  - customization
  - aesthetics
---

# Themes, Palettes, and Aesthetics

markata-go makes it easy to customize your site's appearance. You can go from zero configuration to a beautiful site, then progressively customize as needed.

> **Prerequisites:** This guide assumes you've completed the [Getting Started](/docs/getting-started/) guide and have a working markata-go site.

## Quick Start

The fastest way to change your site's look is to pick a color palette:

```toml
[markata-go.theme]
palette = "catppuccin-mocha"
```

That's it! Your entire site now uses the Catppuccin Mocha color scheme.

---

## Header Preferences

The default header keeps search prominent and groups secondary settings under
the **Appearance and preferences** sliders icon. Open it to select a light or
dark palette, customize the style and font, or change reading size. Expand
**Loading: Auto** at the bottom to choose **Auto** (recommended),
**Save Data**, or **Full Quality** for media. The browser remembers your
preferences, and Auto can suggest a better choice based on the connection.

Press `/` to focus search, `t` to open the appearance panel, or `\\` to
toggle light/dark mode. The panel closes with Escape. On sites where the full
palette picker is disabled, the loading control stays available in the header.

## Available Palettes

markata-go includes 10 built-in color palettes. Use `markata-go palette list` to see them all:

| Palette | Variant | Description |
|---------|---------|-------------|
| `default-light` | light | Clean, minimal light theme |
| `default-dark` | dark | Clean, minimal dark theme |
| `catppuccin-mocha` | dark | Soothing pastel colors, high contrast |
| `catppuccin-latte` | light | Soothing pastel colors, light variant |
| `nord-dark` | dark | Arctic, north-bluish colors |
| `gruvbox-dark` | dark | Retro groove, warm colors |
| `dracula` | dark | Vibrant purple on dark |
| `rose-pine` | dark | Soho vibes, natural colors |
| `solarized-dark` | dark | Scientifically designed colors |
| `tokyo-night` | dark | Tokyo city lights inspired |

### Previewing Palettes

Preview a palette before using it:

```bash
# Show palette colors and contrast info
markata-go palette info catppuccin-mocha

# Check WCAG accessibility compliance
markata-go palette check catppuccin-mocha
```

---

## Configuration Levels

markata-go supports progressive customization - start simple and add complexity only when needed.

### Level 1: Just Pick a Palette

```toml
[markata-go.theme]
palette = "dracula"
```

### Level 2: Override Specific Colors

Keep your palette but tweak a few colors:

```toml
[markata-go.theme]
palette = "nord-dark"

[markata-go.theme.variables]
"--color-primary" = "#88c0d0"
"--color-link" = "#8fbcbb"
```

### Level 3: Add Custom CSS

Add a custom CSS file that loads after the theme:

```toml
[markata-go.theme]
palette = "catppuccin-mocha"
custom_css = "my-styles.css"
```

Then create `static/my-styles.css`:

```css
/* Override any theme styles */
.post-title {
  font-family: 'Georgia', serif;
}

.site-header {
  border-bottom: 2px solid var(--color-primary);
}
```

### Level 4: Override Templates

Override specific template files by creating them in your `templates/` directory:

```
my-site/
├── templates/
│   └── partials/
│       └── footer.html    # Your custom footer
└── markata-go.toml
```

Your custom templates take precedence over theme templates.

---

## Theme Configuration Reference

Full configuration options:

```toml
[markata-go.theme]
# Theme name (currently only "default" is available)
name = "default"

# Color palette to use
palette = "catppuccin-mocha"

# CSS variable overrides
[markata-go.theme.variables]
"--color-primary" = "#8b5cf6"
"--color-background" = "#1a1a2e"
"--color-text" = "#eaeaea"
"--color-link" = "#06b6d4"
"--color-link-hover" = "#22d3ee"
"--content-width" = "800px"
"--font-family" = "'Inter', sans-serif"

# Custom CSS file (relative to static/ directory)
custom_css = "custom.css"
```

### Available CSS Variables

These CSS custom properties can be overridden:

| Variable | Description | Default |
|----------|-------------|---------|
| `--color-background` | Page background | Depends on palette |
| `--color-text` | Body text | Depends on palette |
| `--color-primary` | Primary accent color | Depends on palette |
| `--color-link` | Link color | Depends on palette |
| `--color-link-hover` | Link hover color | Depends on palette |
| `--color-border` | Soft hairline border (text ink mixed 16% into the background) | Depends on palette |
| `--color-border-strong` | Full-contrast border ink for focus rings and emphasis | Depends on palette |
| `--color-border-soft` | `--color-border` at 60% alpha, for section dividers | Derived |
| `--heading-rule` | Color of the short accent bar under `h1`/`h2` | `color-mix(in srgb, var(--color-primary) 55%, transparent)` |
| `--radius-sm` / `--radius` / `--radius-lg` / `--radius-xl` | Corner radius scale | `0.375rem` / `0.5rem` / `0.75rem` / `1rem` |
| `--leading-prose` | Article body line-height | `1.7` |
| `--color-code-bg` | Code block background | Depends on palette |
| `--color-code-text` | Code block text | Depends on palette |
| `--color-code-comment` | Code comments | Depends on palette |
| `--color-code-keyword` | Code keywords | Depends on palette |
| `--color-code-string` | Code strings | Depends on palette |
| `--color-code-number` | Code numbers | Depends on palette |
| `--color-code-function` | Code functions | Depends on palette |
| `--color-code-type` | Code types and tags | Depends on palette |
| `--color-code-operator` | Code operators | Depends on palette |
| `--content-width` | Article measure in `ch` (relative to article text size) | `64ch` at the `large` preset |
| `--font-family` | Body font | System fonts |
| `--font-family-mono` | Code font | Monospace fonts |
| `--article-progress-height` | Sticky article progress bar height | `3px` |
| `--article-progress-track` | Background for the empty progress track | `color-mix(in srgb, var(--color-text) 8%, transparent)` |
| `--article-progress-start` | Gradient start color for the filled portion | `var(--color-primary)` |
| `--article-progress-end` | Gradient end color for the filled portion | `var(--color-primary-light, var(--color-primary))` |
| `--article-progress-glow` | Glow color around the indicator | `color-mix(in srgb, var(--color-primary) 35%, transparent)` |

---

## Article Reading Progress Indicator

Each article page renders a slim, sticky progress track (`.article-progress`) that follows the reader as they scroll. The indicator is powered by the `initArticleProgressIndicator` script, which throttles scroll events with `requestAnimationFrame` and updates the fill amount by transforming `.article-progress__indicator`. The fill tracks the scroll position directly, with no easing lag, and uses the palette's primary colors in light and dark mode. The track is hidden on non-post pages and honors `prefers-reduced-motion` via CSS.

Override the CSS variables above to tune the look. Example:

```toml
[markata-go.theme.variables]
"--article-progress-height" = "5px"
"--article-progress-track" = "rgba(255, 255, 255, 0.35)"
"--article-progress-start" = "#facc15"
"--article-progress-end" = "#fb923c"
"--article-progress-glow" = "rgba(250, 204, 21, 0.8)"
```

Because the indicator is just another element in the post template, you can also restyle it by targeting `.article-progress` and `.article-progress__indicator` from a custom CSS file loaded via `markata-go.theme.custom_css`.

---

## Dark Mode Support

markata-go uses an explicit theme mode model:
- use dark by default (or `theme.fallback_mode`), and
- switch to light only when the visitor explicitly chooses light mode.

### Using Different Palettes for Light/Dark

```toml
[markata-go.theme]
palette = "catppuccin-latte"        # Light mode
palette_dark = "catppuccin-mocha"   # Dark mode palette
```

By default, the site renders in dark mode.
When a visitor explicitly switches to light mode, that preference is saved in localStorage and reused on future visits.

```toml
[markata-go.theme]
palette = "catppuccin-latte"
palette_dark = "catppuccin-mocha"
fallback_mode = "dark"  # or "light"
```

### Every Palette Has Both Modes

You never have to pick a pair by hand. Every palette resolves to a light and a
dark variant:

1. **Explicit families** — `everforest-light`/`everforest-dark`,
   `catppuccin-latte`/`catppuccin-mocha`, `rose-pine-dawn`/`rose-pine`, and
   similar named pairs are used as-is.
2. **Derived counterparts** — a palette that ships only one variant
   (`dracula`, `matte-black`, `monokai`, the Lospec palettes, …) gets an
   automatically derived counterpart named `<palette>-light` or
   `<palette>-dark`. The derivation keeps every hue, compresses backgrounds
   into a soft near-white (or near-black) band, and pushes text and links
   until they meet WCAG AA against every derived surface.

```toml
[markata-go.theme]
palette = "dracula"   # dark mode: dracula, light mode: dracula-light (derived)
```

Derived names work anywhere a palette name is accepted (`palette_light`,
`palette_dark`, calendar rules, the switcher include/exclude lists) and show
up in the multi-palette switcher as the family's other variant. If you prefer
a hand-tuned light theme for a dark-only palette, set `palette_light`
explicitly and the derived one is ignored.

---

## Live Theme Picker

Every site ships with a live theme picker. Visitors can restyle the site in any built-in palette without changing your configured default. The header only gains two small round buttons:

- **Mode toggle** (sun/moon): switch between light and dark.
- **Theme button** (painter's palette icon, dotted with the current theme's accent colors): opens the picker.

The picker opens as a compact popover (a bottom sheet on phones, so the page stays visible above it) with three tabs. Every tab is a grid of live preview cards:

- **Colors**: each card renders in that theme's own background, surface, text, accent, and status colors. Includes search, Light/Dark buttons to browse either mode, Shuffle for a random theme, a **Seasonal** card, and a calendar button that shows the seasonal schedule.
- **Style**: each card draws its corners, borders, and shadows in that aesthetic.
- **Font**: each card shows a heading, body line, and code sample in that font pack. Text size buttons (Small to X-Large) sit above the grid.

With a mouse, hovering a card shows that theme, style, or font on the whole page right away. Move off the cards and the page goes back to your current choice; click to keep it. Hovering a holiday or season in the seasonal schedule previews its colors the same way. On touch screens, tap to choose.

To move fast, use the **‹ ›** buttons in the picker's top bar. They step through the current tab and show the name of the current choice, and the page updates on every tap. Arrow keys, Home/End, and PageUp/PageDown do the same from the keyboard. The picker remembers the last tab you used.

The bottom bar has a labeled **Reset** button, which goes back to the site's default theme, style, and font. While you run `markata-go serve`, it also shows **Bake**, which writes the current choices into your site config (see below). Published sites never show Bake.

On phones, every control is at least 40px tall, the search field uses 16px text so iOS does not zoom, and the sheet leaves room for the home indicator.

Sites are dark by default (`fallback_mode = "dark"`): a first-time visitor sees the dark palette regardless of their OS preference until they flip the mode toggle.

### Baking a look into your config

Run `markata-go serve`, open the picker, and choose a theme, style, font, and text size. Click **Bake**; the button changes to **Bake into markata-go.toml?** (naming the file it will edit, and marked **(override)** when that is a `--merge-config` file). Click again within four seconds to write it. The dev server updates the file, rebuilds, and live-reloads, so every visitor now starts with that look:

```toml
[markata-go.theme]
palette = "gruvbox-dark"
palette_light = "gruvbox-light"
palette_dark = "gruvbox-dark"
fallback_mode = "dark"
aesthetic = "brutal"
fontpack = "editorial"
text_size = "large"
```

`fallback_mode` and `text_size` are only baked when you picked them yourself in the picker; otherwise bake leaves the site's settings alone, because they are usually each visitor's preference.

Bake picks the file carefully:

- It looks at the root config, every file it pulls in with `include = [...]` (for example a `config/` directory), and any `--merge-config` files.
- If one or more already have a `[markata-go.theme]` section, it edits the one that takes effect (the last to load). A theme in `config/theme.toml` stays there.
- Otherwise it adds `[markata-go.theme]` to the root config, or creates `markata-go.toml` when the site has none.
- If the site has no config and serve is using your global `~/.config/markata-go/config.toml`, the button reads **Bake unavailable**: bake never edits your global defaults. Create a `markata-go.toml` in the site first.

Only the baked keys change. Comments, key order, and other settings are kept, and existing keys are updated in place. TOML, YAML, and JSON configs all work. If the theme is written in a shape bake cannot edit safely (a TOML inline table `theme = { ... }`, a YAML flow mapping, or a multi-line value), bake leaves the file alone and reports the error; convert it to a normal table and try again. `MARKATA_GO_THEME_*` environment variables still override the file, and bake warns when one is set.

After a successful bake the picker names the keys it wrote and clears your browser's saved picks, so what you see is the site default. Bake uses the same checks as the settings sidebar: if the config does not load back as written, every file is restored. A baked key also replaces any unsaved sidebar preview of it.

To change other theme settings, such as `palette_light`, `custom_css`, or `switcher.enabled`, click **All settings** next to Bake. It opens the serve-only settings sidebar at the Theme section (see [Edit Settings While Serving](/docs/guides/configuration/#edit-settings-while-serving)).

### Seasonal

The **Seasonal** card, first in the Colors grid, picks a palette from today's date. It follows the northern hemisphere seasons:

| Season | Starts | Palette |
|--------|--------|---------|
| Spring | Mar 20 | `pollen8` |
| Summer | Jun 21 | `summer-beach` |
| Autumn | Sep 22 | `autumn` |
| Winter | Dec 21 | `winter-frost` |

It switches to a holiday palette for the three days before a holiday and on the day itself:

| Holiday | Date | Palette |
|---------|------|---------|
| New Year | Jan 1 | `white-gold` / `black-gold` |
| Lunar New Year | varies | `lunar-new-year` |
| Valentine's Day | Feb 14 | `valentine` |
| St. Patrick's Day | Mar 17 | `st-patricks` |
| Easter | varies (Western) | `blessing` |
| Earth Day | Apr 22 | `everforest-light` / `everforest-dark` |
| Halloween | Oct 31 | `halloween` |
| Diwali | varies (Lakshmi Puja) | `diwali` |
| Hanukkah | varies, all 8 days | `hanukkah` |
| Christmas | Dec 25 | `christmas` |

To see the schedule, open the picker and click the calendar button next to Shuffle on the Colors tab. It shows:

- a strip for the current year, colored with each season's and holiday's theme, with a line at today
- a **Coming up** list: what is showing now, until when, and how long is left, then each upcoming holiday (its dates and how many days it shows, lead days included) and season change (its length in weeks) over the next year, each with a color swatch
- a **Use Seasonal** button, if you picked a different theme

Hover a band in the strip or a row in the list to preview that period's colors on the page. The hovered band and its row are outlined together, a line under the strip spells out the period's dates and length, and the picker's top label reads **Preview: …** while anything is being previewed. On touch screens, tap a band or row to see its dates without changing the page.

The schedule follows the Light/Dark buttons, so you see the palettes for the mode you are browsing. Click the calendar button again, search, or step with ‹ › to go back to the grid.

Seasonal follows the light/dark toggle and uses each palette's matching variant. The date comes from the visitor's clock, and the head script resolves it before the first paint, so nothing flashes. Easter is calculated. Other moving holidays come from a built-in table: Lunar New Year and Hanukkah through 2035, Diwali through 2033. Holidays and seasons whose palettes are excluded by `[markata-go.theme.switcher]` are skipped.

To make seasonal the default for every visitor who has not picked a theme, set `seasonal = true`. Your `palette` settings stay as the fallback (for example, when the picker is disabled):

```toml
[markata-go.theme]
palette = "ayu-dark"
seasonal = true
```

**Bake** writes `seasonal = true` when Seasonal is selected and leaves your `palette` settings as the fallback. Baking a regular theme sets `seasonal = false` if seasonal was on.

### Styles

The **Style** tab switches the aesthetic, which sets the shape and depth of code blocks, cards, admonitions, tables, images, and form controls:

| Style | Look |
|-------|------|
| `minimal` | Default. Soft 1px borders, small radii, no shadows |
| `balanced` | Medium radii with a light drop shadow |
| `elevated` | Large radii, borderless surfaces lifted by soft shadows |
| `precision` | Near-square corners with crisp hairline borders |
| `brutal` | Square corners, 2px ink borders, and hard offset shadows in the palette's accent color, applied to cards, code, tables, buttons, inputs, search, and share buttons |

### Fonts

The **Font** tab lists every bundled font pack as a preview card with its heading, body, and code faces. The default is the zero-download `system` pack: platform-native **system-ui** headings/body and **ui-monospace** code. The expressive **Brush** (Knewave, Space Grotesk, DM Mono) remains opt-in. Bundled alternative fonts are self-hosted, and a preview only downloads a face if selected or scrolled into view. With the font picker enabled, the build may still copy bundled font files to `assets/fonts/` so visitors can switch packs. Set `fontpack` in `[markata-go.theme]` to override the default for a site.

The active font pack preloads its distinct body, display, and heading faces.
Packs with separate display and heading fonts keep both ready for the first
paint. Shared faces are fetched once; code fonts and extended character
coverage load when needed.

### No flash on load

A visitor's choice is restored by a tiny inline script in `<head>` before any stylesheet paints, so return visits render in the chosen theme from the first frame. Transitions are suppressed during a switch, and view-transition navigation keeps the chosen theme.

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `t` | Open or close the theme picker |
| `.` | Next theme (shows a toast) |
| `,` | Previous theme |
| `>` | Next style (aesthetic) |
| `<` | Previous style |
| `f` | Next font |
| `F` | Previous font |
| `\` | Toggle dark/light mode |

Inside the picker:

- `1`, `2`, and `3` switch to the Colors, Style, and Font tabs.
- Arrow keys preview.
- `Enter` keeps the choice and closes.
- `Escape` closes.

On the tab bar, `←`/`→` move between tabs and `↓` moves into the grid.

### How choices are stored

Choices are saved per mode in `localStorage`:

| Key | Value |
|-----|-------|
| `color-mode` | `light` or `dark` |
| `theme-palette-light` | Palette used in light mode, or `seasonal` |
| `theme-palette-dark` | Palette used in dark mode, or `seasonal` |
| `theme-aesthetic` | Chosen aesthetic |
| `theme-fontpack` | Chosen font pack |
| `text-size` | Chosen reading size |
| `theme-picker-tab` | Last picker tab (`colors`, `style`, `font`) |

Picking a theme also stores its counterpart for the other mode. For example, choosing Gruvbox in dark mode makes the mode toggle switch to Gruvbox Light. Picking your site's default pair clears the stored keys, so the visitor follows future default changes again.

### Disabling the picker

```toml
[markata-go.theme.switcher]
enabled = false      # remove the picker; ship only your palette
mode_toggle = true   # keep the sun/moon toggle
```

With the picker disabled, only the configured palette's CSS and font pack are generated, and the text-size control goes back to the header.

### Filtering Palettes

By default, all discovered palettes are included. You can control which palettes appear:

**Exclude specific palettes:**

```toml
[markata-go.theme.switcher]
include_all = true  # Default
exclude = ["default-light", "default-dark"]  # Hide these palettes
```

**Include only specific palettes:**

```toml
[markata-go.theme.switcher]
include_all = false
include = ["catppuccin-mocha", "catppuccin-latte", "nord-dark", "nord-light"]
```

Your configured default palettes are always available in the picker.

### Switcher Configuration Reference

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `enabled` | boolean | `true` | Show the live theme picker |
| `mode_toggle` | boolean | `true` | Show the dark/light mode toggle |
| `include_all` | boolean | `true` | Include all discovered palettes |
| `include` | array | `[]` | Palettes to include (when `include_all` is false) |
| `exclude` | array | `[]` | Palettes to exclude (when `include_all` is true) |
| `position` | string | `"header"` | Where to place the switcher |

The mode toggle is also gated by `[markata-go.header].show_theme_toggle` for backward compatibility.

All palettes add roughly 30KB (gzip) to `palette.css`. Use `include` to trim that on very size-sensitive sites.

### Reading-size control

The default theme uses a large, comfortable reading size by default. Visitors
can choose a smaller or larger preset (Small, Medium, Large, X-Large) from the theme picker (or the header control when the picker is off), and their
choice is saved for later visits on the same site.

```toml
[markata-go.theme]
text_size = "large"              # small, medium, large, or x-large
show_text_size_control = true    # default: true
```

`text_size` is used when the visitor has not made a choice. The available
presets are:

| Preset | Site base | Article text | Article measure |
|--------|-----------|--------------|-----------------|
| `small` | 16px | 18px | 68ch |
| `medium` | 17px | 20px | 66ch |
| `large` | 18px | 22px | 64ch |
| `x-large` | 19px | 24px | 62ch |

Article text also scales with the viewport so wide desktop displays do not
render a narrow strip of small type: `--reading-scale` multiplies the article
font size by 1.08 from 1800px, 1.16 from 2200px (1440p), and 1.25 from
3000px (4K). The `ch`-based measure follows the scaled font, so line length
stays comfortable. Site chrome (`--text-base`) is not scaled.

Set `show_text_size_control = false` when a site should keep the configured
default without rendering the selector. Browser zoom remains available in
all modes.

Presets are defaults, not overrides. If `[theme.variables]` sets
`--text-base`, `--post-text-size`, or `--content-width` explicitly, the
author's value wins for that token in every preset; the selector then only
affects the tokens the site has not pinned.

### How It Works

1. **Palette manifest**: `palette.css` includes a JSON manifest (`--palette-manifest`) listing each theme's `name`, `displayName`, `variant` (`light`/`dark`), and mode `counterpart`.

2. **Scoped palettes**: each palette is emitted as a `[data-palette="name"]` block. The page sets `data-palette` on `<html>`, and each preview card sets it on itself, so cards render with the real palette CSS.

3. **Style tokens**: per-aesthetic surface tokens (`--radius-*`, `--surface-border`, `--surface-shadow`) are keyed on `[data-aesthetic="name"]`, so each Style card previews its own aesthetic.

4. **Font manifest**: `css/fonts.css` includes `--fontpack-manifest` and `--fontpack-default`, plus one `[data-fontpack="name"]` block per pack. Each manifest entry carries the CSS font stacks (`headingFont`, `bodyFont`, `codeFont`) that Font cards apply inline.

5. **Seasonal calendar**: when the picker is on, `base.html` inlines the seasonal schedule as `window.__markataSeasonal`, resolved against the palettes that ship.

6. **Motifs follow the theme**: background motifs are off by default. When you enable one (`[markata-go.theme.motif] kind = "block-w"`), it is repainted from the active palette's colors whenever the page runs a palette other than the configured one.

### Styling the Picker

```css
.theme-picker-toggle { /* round palette button in the header */ }
.theme-picker-panel  { /* popover / bottom sheet */ }
.theme-picker-tabs   { /* Colors / Style / Font tabs */ }
.theme-picker-step   { /* ‹ › step buttons */ }
.theme-card          { /* any preview card */ }
.palette-card, .seasonal-card, .style-card, .font-card { /* per-tab cards */ }
.theme-picker-schedule { /* seasonal schedule view (.ss-year strip, .ss-list) */ }
.theme-card[aria-selected="true"] { /* current theme */ }
.theme-picker-bake   { /* "Bake" button (markata-go serve only) */ }
.theme-picker-settings { /* "All settings" button (markata-go serve only) */ }
```

### JavaScript API

```javascript
const picker = window.markata.paletteSwitcher;

document.documentElement.dataset.palette; // current theme, e.g. "gruvbox-dark"
picker.selectTheme("nord-dark");  // apply and remember a theme
picker.nextFamily();              // next theme in the current mode
picker.prevFamily();
picker.randomTheme();
picker.toggleColorMode();         // flip light/dark
picker.setAesthetic("brutal");
picker.setFontpack("editorial");
picker.nextFont();                // cycle fonts (also prevFont, cycleFont)
picker.setTextSize("x-large");    // small, medium, large, x-large
picker.selectSeasonal();          // follow the seasonal calendar
picker.getSeasonal("dark");       // today's pick: { label, name }
picker.getBakeSettings();         // current choices as [markata-go.theme] keys
picker.bake();                    // serve only: write them to the site config
picker.open("font");              // show the picker (optionally on a tab)
picker.close();
picker.resetTheme();              // back to the site default
picker.getManifest();             // all themes with variant and counterpart
```

### Event Handling

```javascript
window.addEventListener('palette-change', (e) => {
  console.log('Palette:', e.detail.palette);
});

window.addEventListener('color-mode-change', (e) => {
  console.log('Mode:', e.detail.mode);
});

window.addEventListener('aesthetic-change', (e) => {
  console.log('Aesthetic:', e.detail.aesthetic);
});

window.addEventListener('fontpack-change', (e) => {
  console.log('Font pack:', e.detail.fontpack);
});

window.addEventListener('text-size-change', (e) => {
  console.log('Text size:', e.detail.size);
});
```

---

## Palette CLI Commands

### List All Palettes

```bash
markata-go palette list
```

Output:
```
NAME                      VARIANT  SOURCE     DESCRIPTION
----------------------------------------------------------------------
Catppuccin Mocha          dark     builtin    Soothing pastel theme
Catppuccin Latte          light    builtin    Soothing pastel theme (light)
Nord Dark                 dark     builtin    Arctic, north-bluish
...
```

### Get Palette Info

```bash
markata-go palette info catppuccin-mocha
```

Shows all colors in the palette with their hex values.

### Check Accessibility

```bash
markata-go palette check catppuccin-mocha
```

Checks WCAG 2.1 AA contrast requirements for text readability.

### Export Palette

Export a palette to different formats:

```bash
# Export as CSS custom properties
markata-go palette export catppuccin-mocha --format css

# Export as SCSS variables
markata-go palette export catppuccin-mocha --format scss

# Export as JSON
markata-go palette export catppuccin-mocha --format json

# Export as Tailwind config
markata-go palette export catppuccin-mocha --format tailwind
```

### Create New Palette

Generate a starter palette file:

```bash
markata-go palette new my-palette
```

Creates `palettes/my-palette.toml` that you can customize.

### Clone Existing Palette

Clone an existing palette as a starting point for customization:

```bash
# Interactive mode - opens fuzzy picker to select a palette
markata-go palette clone

# Clone a specific palette by name
markata-go palette clone catppuccin-mocha

# Clone with explicit new name
markata-go palette clone catppuccin-mocha --name "my-custom-theme"
```

The cloned palette is saved to `~/.config/markata-go/palettes/` and includes:
- All colors from the source palette
- Semantic and component color mappings
- A description noting the source palette

This is useful when you want to:
- Start from a well-designed palette and make minor adjustments
- Create variants of existing palettes (e.g., higher contrast)
- Experiment without modifying the original

### Fetch Palette from Lospec

Import color palettes directly from [Lospec.com](https://lospec.com/palette-list), a popular source for pixel art and retro color palettes:

```bash
# Fetch a palette by URL
markata-go palette fetch https://lospec.com/palette-list/sweetie-16.txt

# Use a custom name
markata-go palette fetch https://lospec.com/palette-list/cheese-palette.txt --name "My Cheese"

# Save to a custom directory
markata-go palette fetch https://lospec.com/palette-list/tokyo-night.txt -o palettes/
```

The fetched palette is saved to your user palettes directory (`~/.config/markata-go/palettes/`) by default. markata-go automatically:

- Downloads the color list from Lospec
- Analyzes colors to determine if it's a light or dark theme
- Generates semantic mappings (text-primary, bg-primary, accent, etc.)
- Caches the download to avoid repeated network requests

**Example output:**
```
Fetching palette from: https://lospec.com/palette-list/sweetie-16.txt
Palette saved to: /home/user/.config/markata-go/palettes/sweetie-16.toml

Palette details:
  Name:        Sweetie 16
  Variant:     dark
  Colors:      16
  Source:      https://lospec.com/palette-list/sweetie-16.txt

Semantic mappings:
  bg-primary      -> color0      (#1a1c2c)
  text-primary    -> color15     (#f4f4f4)
  accent          -> color8      (#b13e53)

Use this palette in your config:
  [markata-go.theme]
  palette = "sweetie-16"
```

### Pick Palette Interactively

Browse all available palettes in a full-screen interactive TUI with live color previews:

```bash
# Pick a palette and set it in your config (default)
markata-go palette pick

# Only print the name without updating config
markata-go palette pick --no-set
```

The picker shows a two-panel layout:

- **Left panel** -- Fuzzy-filterable list of all palettes with variant badges (`[dark]`/`[light]`). Type to filter, use arrow keys to navigate.
- **Right panel** -- Live preview of the highlighted palette showing color swatches, semantic roles, and a contrast preview block.

Press **Enter** to select a palette and set it in your config, or **Esc** to cancel.

**Compose with other commands:**

```bash
# View detailed info for the palette you pick (without setting it)
markata-go palette info "$(markata-go palette pick --no-set)"

# Pick, set, and rebuild
markata-go palette pick && markata-go build
```

---

## Creating Custom Palettes

Create a custom palette by adding a TOML file to `palettes/` in your project:

```toml
# palettes/my-brand.toml
[palette]
name = "My Brand"
variant = "light"  # or "dark"
author = "Your Name"

# Raw colors
[palette.colors]
brand-primary = "#3b82f6"
brand-secondary = "#8b5cf6"
brand-accent = "#06b6d4"
white = "#ffffff"
gray-50 = "#f9fafb"
gray-100 = "#f3f4f6"
gray-700 = "#374151"
gray-800 = "#1f2937"
gray-900 = "#111827"

# Semantic mapping
[palette.semantic]
text-primary = "gray-800"
text-secondary = "gray-700"
bg-primary = "white"
bg-secondary = "gray-50"
accent = "brand-primary"
link = "brand-primary"
link-hover = "brand-secondary"
```

Then use it:

```toml
[markata-go.theme]
palette = "my-brand"
```

---

## Aesthetics

While palettes control **colors**, aesthetics control **form and layout** - the non-color design tokens that define your site's visual personality.

> **Think of it this way:** Palettes are *what colors* you use. Aesthetics are *how things are shaped*.

### Quick Start

Set an aesthetic in one line:

```toml
[markata-go]
aesthetic = "elevated"
```

Your entire site now uses generous rounded corners, layered shadows, and comfortable spacing.

### Understanding the Separation

| Aspect | Palettes | Aesthetics |
|--------|----------|------------|
| **Controls** | Colors | Shape, spacing, depth |
| **Tokens** | `--color-*` variables | `--radius-*`, `--shadow-*`, `--spacing-*` |
| **Examples** | Background, text, accent colors | Border radius, shadow depth, spacing scale |
| **Switch** | Changes color scheme | Changes visual "feel" |

This separation means you can:
- Use any palette with any aesthetic
- Switch aesthetics without changing colors
- Fine-tune form independently of color

### Available Aesthetics

markata-go includes 5 built-in aesthetics:

| Aesthetic | Description | Best For |
|-----------|-------------|----------|
| `balanced` | **Default.** Comfortable rounding, subtle shadows, normal spacing | General purpose, blogs |
| `brutal` | Sharp corners, thick borders, tight spacing, hard accent-colored offset shadows | Bold statements, portfolios |
| `minimal` | Restrained rounding, airy whitespace, flat surfaces and selective borders | Documentation, reading-focused |
| `elevated` | Generous rounding, layered shadows, generous spacing | Premium/SaaS, card-heavy layouts |
| `precision` | Subtle corners, compact spacing, hairline borders, minimal shadows | Technical docs, data-heavy sites |

### Visual Comparison

**Brutal:**
```
┌────────────────────────┐
│ No rounding            │
│ Thick 3px borders      │
│ Compact spacing        │
│ Hard accent shadows    │
└────────────────────────┘
```

**Balanced (default):**
```
╭────────────────────────╮
│ Subtle 4-8px rounding  │
│ Normal 1px borders     │
│ Standard spacing       │
│ Light shadows          │
╰────────────────────────╯
```

**Elevated:**
```
╭────────────────────────╮
│                        │
│ Generous 16px rounding │
│ Minimal borders        │
│ Generous spacing       │
│ Layered shadows ▓▒░    │
│                        │
╰────────────────────────╯
```

### Configuration Examples

Styles affect the feed and post cards, code, admonitions, quotes, tables and
images, buttons and inputs, tags, header edge, picker panel, prose rhythm and
open doc/feed drawers. Balanced keeps the theme's native surface treatment.
Precision is compact and crisp; elevated is soft and lifted; minimal stays
flat and airy. Style cards preview their own shape and depth without changing
the active palette. Site `custom_css` and `[theme.variables]` can still override
the layered aesthetic rules.

**Basic usage:**

```toml
[markata-go]
aesthetic = "brutal"
```

**Combine with a palette:**

```toml
[markata-go]
aesthetic = "elevated"

[markata-go.theme]
palette = "catppuccin-mocha"
```

**Override specific tokens:**

```toml
[markata-go]
aesthetic = "balanced"

[markata-go.aesthetic_overrides]
border_radius = "12px"      # Override just the radius
shadow_intensity = 1.5      # Make shadows 50% stronger
spacing_scale = 1.1         # Slightly more spacing
```

### Aesthetic Overrides Reference

Fine-tune any aesthetic with these overrides:

| Override | Type | Description | Example |
|----------|------|-------------|---------|
| `border_radius` | string | Base border radius | `"8px"`, `"0.5rem"` |
| `border_width` | string | Border thickness | `"1px"`, `"2px"` |
| `border_style` | string | Border style | `"solid"`, `"dashed"`, `"none"` |
| `spacing_scale` | float | Multiplier for spacing | `0.9` (tighter), `1.2` (looser) |
| `shadow_intensity` | float | Multiplier for shadow opacity | `0` (none), `1.5` (stronger) |
| `shadow_size` | string | Shadow size preset | `"sm"`, `"md"`, `"lg"`, `"none"` |

**Example: Make "balanced" more spacious:**

```toml
[markata-go]
aesthetic = "balanced"

[markata-go.aesthetic_overrides]
spacing_scale = 1.25
shadow_intensity = 0.5  # Lighter shadows
```

**Example: Soften "brutal":**

```toml
[markata-go]
aesthetic = "brutal"

[markata-go.aesthetic_overrides]
border_radius = "4px"   # Add slight rounding
border_width = "2px"    # Thinner borders
```

### Aesthetic CLI Commands

List all available aesthetics:

```bash
markata-go aesthetic list
```

Output:
```
Available aesthetics:
  brutal     - Brutalist design: harsh, uncompromising, raw
  precision  - Technical/engineering: clean, exact, minimal
  balanced   - Default harmonious: comfortable, balanced
  elevated   - Layered/premium: depth, floating cards
  minimal    - Maximum whitespace: sparse, intentional
```

Show details of a specific aesthetic:

```bash
markata-go aesthetic show elevated
```

Output:
```
Aesthetic: elevated
Description: Layered/premium: depth, floating cards

Tokens:
  radius:  0.5rem (sm), 0.75rem (md), 1rem (lg)
  spacing: 1.25x scale
  border:  none
  shadow:  0 4px 12px rgba(0,0,0,0.15)

CSS Preview:
  --radius-sm: 0.5rem;
  --radius-md: 0.75rem;
  --radius-lg: 1rem;
  --radius-full: 9999px;
  --shadow: 0 4px 12px rgba(0,0,0,0.15);
```

### Creating Custom Aesthetics

Create a custom aesthetic by adding a TOML file to `aesthetics/` in your project:

```toml
# aesthetics/my-aesthetic.toml
name = "My Aesthetic"
description = "Custom design for my brand"

[tokens.radius]
none = "0"
sm = "6px"
md = "10px"
lg = "16px"
xl = "24px"
full = "9999px"

[tokens.spacing]
scale = 1.1

[tokens.border]
width_thin = "1px"
width_normal = "2px"
width_thick = "3px"
style = "solid"

[tokens.shadow]
sm = "0 1px 2px rgba(0,0,0,0.05)"
md = "0 4px 8px rgba(0,0,0,0.1)"
lg = "0 8px 16px rgba(0,0,0,0.12)"
xl = "0 16px 32px rgba(0,0,0,0.15)"

[tokens.typography]
font_primary = "var(--font-sans)"
leading_scale = 1.0
```

Then use it:

```toml
[markata-go]
aesthetic = "my-aesthetic"
```

### Keyboard Shortcuts

When the palette switcher is enabled, these shortcuts also work for aesthetics:

| Key | Action |
|-----|--------|
| `.` | Next palette family (full switcher only) |
| `,` | Previous palette family (full switcher only) |
| `>` | Next aesthetic (full switcher only) |
| `<` | Previous aesthetic (full switcher only) |
| `\` | Toggle dark/light mode (mode toggle only) |

---

## Template Overrides

Override any template by placing it in your `templates/` directory.

### Template Search Order

1. `templates/` - Your project templates (highest priority)
2. `themes/{theme}/templates/` - Theme templates
3. Embedded default templates (fallback)

### Available Templates

| Template | Purpose |
|----------|---------|
| `base.html` | HTML skeleton, head, header, footer |
| `post.html` | Single post/article layout |
| `feed.html` | List of posts (index, archive, tags) |
| `card.html` | Post preview card in feeds |
| `partials/header.html` | Site header/navigation |
| `partials/footer.html` | Site footer |
| `partials/head.html` | Additional head content |

### Example: Custom Footer

Create `templates/partials/footer.html`:

```html
<footer class="site-footer">
  <div class="container">
    <p>&copy; {{ now | date:"2006" }} {{ config.title }}. Built with markata-go.</p>
    <nav>
      <a href="/about/">About</a>
      <a href="/contact/">Contact</a>
      <a href="https://github.com/yourusername">GitHub</a>
    </nav>
  </div>
</footer>
```

---

## Static Assets

Add custom CSS, JavaScript, images, and fonts to the `static/` directory:

```
my-site/
├── static/
│   ├── css/
│   │   └── custom.css
│   ├── js/
│   │   └── analytics.js
│   ├── images/
│   │   └── logo.png
│   └── fonts/
│       └── MyFont.woff2
└── markata-go.toml
```

Files in `static/` are copied directly to the output directory.

Reference them in templates:

```html
<link rel="stylesheet" href="/css/custom.css">
<script src="/js/analytics.js"></script>
<img src="/images/logo.png" alt="Logo">
```

---

## Per-Post Styling

Override styles for specific posts using frontmatter:

```yaml
---
title: "Special Post"
template: landing.html  # Use a different template
---
```

Or add custom CSS classes:

```yaml
---
title: "Featured Article"
css_class: featured-post
---
```

Then style it:

```css
.featured-post {
  background: linear-gradient(to right, var(--color-primary), var(--color-accent));
}
```

---

## CSS Optimization

### Bundle a homepage

To reduce stylesheet requests on a specific page, enable page bundling:

```toml
[markata-go.css_bundle]
enabled = true
pages = ["index.html"]
```

Use output-relative HTML paths. Production builds combine local head
stylesheets and inline styles in their original order into hashed CSS files.
Scripts and special stylesheets separate bundles. Remote stylesheets, media
conditions, integrity attributes, CSS imports, and relative asset URLs stay
untouched. Fast builds skip this step. Original CSS files remain available for
navigation and custom templates. The runtime loader recognizes bundled sources
and avoids downloading them again.

The navbar reserves space for its search input before deferred search styles
load, so the header stays in place while search initializes.

Title fitting and highlighted headings wait for initial layout before measuring
when ResizeObserver is available. They still update after fonts load and the
page resizes, keeping the same typography with less synchronous startup work.
The article progress indicator also refreshes as images or other content change
the article's height.

Closed feed and table-of-contents drawers skip rendering their contents when
the browser supports content visibility and discrete transitions. Opening
restores the contents immediately, while closing preserves the slide.
Reduced-motion preferences remain supported. Other browsers keep the existing
rendering behavior.

Mobile drawers center the active feed item when opened without measuring
closed-drawer geometry during startup. Desktop drawers
continue to restore their saved state and follow the visible header edge.

With the built-in theme picker, the default template initially loads
`palette-core.css` with colors and picker metadata for the configured light/dark
pair. Generated palette changes also invalidate cached HTML, so warm builds
reference the current stylesheet. The remaining palette catalog and colors load when
opening the picker or restoring a saved or seasonal choice. Custom templates
can continue loading the complete `palette.css` stylesheet.

Image zoom CSS and JavaScript load on the first image click. The first click
opens that image once the library is ready; later clicks reuse the library.
Photo and article cards use responsive image candidates on trusted resize
hosts. For text-only feed previews, use `render_feed("til", 8, "summary")`.

markata-go automatically optimizes CSS loading by only including stylesheets that are actually needed for each page. This reduces page size and improves load times.

### How It Works

When rendering a page, markata-go scans the HTML content and detects which CSS features are used:

| CSS File | Loaded When |
|----------|-------------|
| `variables.css` | Always (core theme variables) |
| `main.css` | Always (core layout styles) |
| `components.css` | Always (navigation, footer, etc.) |
| `cards.css` | Feed/index pages with post cards |
| `admonitions.css` | Posts containing admonition blocks |
| `code.css` | Posts containing code blocks |
| `chroma.css` | Posts with syntax-highlighted code |
| `webmentions.css` | When webmentions are enabled |
| `palette-switcher.css` | When palette switcher is enabled |
| `search.css` | When search is enabled |

### Content Detection

The CSS detection works by analyzing the rendered HTML:

- **Admonitions**: Detected when `class="admonition` is present
- **Code blocks**: Detected when syntax highlighting classes (`class="chroma"`, `class="highlight"`) or code elements (`<pre><code`, `<code class="language-`) are present
- **Cards**: Included on feed pages (where `feed` variable exists in template context)

### Benefits

- **Smaller page sizes**: Simple pages without code blocks or admonitions skip those CSS files
- **Faster load times**: Less CSS to download and parse
- **Better caching**: Core CSS files are shared across all pages

### Custom CSS

Your custom CSS (via `theme.custom_css`) is always loaded when configured. If you need conditional loading for custom styles, consider using CSS custom properties or JavaScript-based loading.

---

## Media Borders and Gradient Effects

markata-go provides beautiful, configurable borders for images and videos. From subtle solid borders to animated gradients, you have full control over how your media looks.

### Default Media Styling

By default, images and videos in your content get:

- Rounded corners (`--media-border-radius`)
- A subtle border (`--media-border-width`, `--media-border-color`)
- Proper spacing and centering

### Configuring Media Borders

Customize the default borders via CSS variables:

```css
/* In your custom CSS or via theme.variables */
:root {
  --media-border-width: 3px;      /* Border thickness */
  --media-border-color: #e5e7eb;  /* Border color */
  --media-border-radius: 0.5rem;  /* Corner rounding */
}
```

Or via config:

```toml
[markata-go.theme.variables]
"--media-border-width" = "4px"
"--media-border-color" = "#8b5cf6"
"--media-border-radius" = "1rem"
```

### Gradient Borders

Enable colorful gradient borders for a modern, eye-catching look. Add a class to your post content to enable gradients for all media:

```yaml
---
title: "My Post with Gradient Borders"
css_class: gradient-borders
---
```

This applies the default accent gradient to all images and videos in that post.

### Available Gradient Presets

markata-go includes several beautiful gradient presets:

| Class | Colors | Best For |
|-------|--------|----------|
| `gradient-borders` | Primary to primary-dark | Brand-consistent |
| `gradient-vibrant` | Purple to pink | Creative, artistic |
| `gradient-warm` | Pink to orange | Energetic, warm |
| `gradient-cool` | Blue to cyan | Professional, tech |
| `gradient-sunset` | Pink to yellow | Warm, inviting |
| `gradient-ocean` | Teal to light blue | Calm, refreshing |

Use them in frontmatter:

```yaml
---
title: "Ocean-Themed Post"
css_class: gradient-ocean
---
```

### Animated Gradient Borders

For extra visual impact, use animated gradients that slowly shift colors:

```yaml
---
title: "Attention-Grabbing Post"
css_class: gradient-animated
---
```

The animation cycles through purple, pink, and blue over 6 seconds.

### Glow Effects

Add a subtle glow behind your media:

```yaml
---
title: "Glowing Media"
css_class: glow
---
```

Combine glow with gradients:

```yaml
---
title: "Maximum Impact"
css_class: gradient-vibrant glow
---
```

### Per-Image Styling

For fine-grained control, add classes directly to images in your Markdown using HTML:

```html
<img src="/images/hero.jpg" alt="Hero" class="gradient-vibrant glow">
```

Or use a wrapper div:

```html
<div class="media-frame gradient-sunset glow">
  <img src="/images/featured.jpg" alt="Featured">
</div>
```

### CSS Variable Reference for Media

| Variable | Description | Default |
|----------|-------------|---------|
| `--media-border-width` | Border thickness | `3px` |
| `--media-border-style` | Border style | `solid` |
| `--media-border-color` | Border color | `var(--color-border)` |
| `--media-border-radius` | Corner rounding | `0.5rem` |
| `--gradient-accent` | Default gradient | Primary colors |
| `--gradient-vibrant` | Purple-pink gradient | `#667eea` to `#f093fb` |
| `--gradient-warm` | Pink-orange gradient | `#f093fb` to `#f8b500` |
| `--gradient-cool` | Blue-cyan gradient | `#4facfe` to `#00f2fe` |
| `--gradient-sunset` | Pink-yellow gradient | `#fa709a` to `#fee140` |
| `--gradient-ocean` | Teal-blue gradient | `#2193b0` to `#6dd5ed` |

### Palette-Matching Gradients

If you're using a specific color palette, use the matching gradient for visual consistency:

| Class | Palette | Colors |
|-------|---------|--------|
| `gradient-catppuccin` | Catppuccin | Mauve → Pink → Red |
| `gradient-nord` | Nord | Frost colors (cyan → blue) |
| `gradient-dracula` | Dracula | Purple → Pink → Cyan |
| `gradient-gruvbox` | Gruvbox | Yellow → Orange → Red |
| `gradient-rose-pine` | Rosé Pine | Iris → Rose → Gold |
| `gradient-solarized` | Solarized | Blue → Cyan → Green |
| `gradient-tokyo-night` | Tokyo Night | Blue → Purple → Pink |

Example: If your site uses `catppuccin-mocha` palette, use `gradient-catppuccin` for borders:

```toml
# markata-go.toml
[markata-go.theme]
palette = "catppuccin-mocha"
```

```yaml
# In your post frontmatter
---
title: "Catppuccin-Styled Gallery"
css_class: gradient-catppuccin
---
```

### Custom Gradients

Create your own gradient by overriding the variables:

```css
/* In static/custom.css */
:root {
  --gradient-accent: linear-gradient(135deg, #ff6b6b, #feca57, #48dbfb);
}
```

Or define a completely new one:

```css
.post-content.gradient-custom img,
.post-content.gradient-custom video {
  border: none;
  padding: 3px;
  background: linear-gradient(45deg, #12c2e9, #c471ed, #f64f59);
  background-origin: border-box;
}
```

Then use in frontmatter:

```yaml
---
css_class: gradient-custom
---
```

### Dark Mode Considerations

Gradient borders adapt to dark mode:
- Glow effects become more prominent
- Border colors adjust automatically
- Gradients remain vibrant on dark backgrounds

Test your gradient choices in both light and dark mode.

---

## Background Decorations

Add multi-layered background decorations to your site for visual effects like snow, particles, stars, or animated elements.

### Basic Configuration

```toml
[markata-go.theme.background]
enabled = true

backgrounds = [
  { html = '<snow-fall count="200"></snow-fall>' },
]

scripts = ["/static/js/snow-fall.js"]
```

This adds a snow effect using a custom web component with its supporting JavaScript.

### Multiple Layers

Stack multiple background layers with different z-index values:

```toml
[markata-go.theme.background]
enabled = true

backgrounds = [
  { html = '<div class="stars"></div>', z_index = -20 },
  { html = '<div class="clouds"></div>', z_index = -10 },
  { html = '<snow-fall count="100"></snow-fall>', z_index = -5 },
]

scripts = ["/static/js/background-effects.js"]

css = '''
.stars {
  position: absolute;
  inset: 0;
  background: url("/images/stars.png") repeat;
  opacity: 0.3;
}

.clouds {
  position: absolute;
  inset: 0;
  background: url("/images/clouds.png") repeat-x;
  animation: drift 60s linear infinite;
}

@keyframes drift {
  from { background-position: 0 0; }
  to { background-position: 100% 0; }
}
'''
```

### Configuration Reference

| Option | Type | Description |
|--------|------|-------------|
| `enabled` | boolean | Enable/disable background decorations (default: false) |
| `backgrounds` | array | List of background elements |
| `backgrounds[].html` | string | HTML content for this layer |
| `backgrounds[].z_index` | integer | Stacking order (-1 is default, behind content) |
| `scripts` | array | Script URLs to load for background functionality |
| `css` | string | Custom CSS for styling background elements |
| `article_bg` | string | Background applied to post/feed content when decorations are enabled |
| `article_blur_enabled` | boolean | Opt in to backdrop blur on content containers (default: false) |
| `article_blur` | string | Blur amount used only when `article_blur_enabled = true` |
| `article_shadow` | string | Box shadow for post/feed content wrappers |
| `article_border` | string | Border for post/feed content wrappers |
| `article_radius` | string | Border radius for post/feed content wrappers |

### Tips for Background Decorations

1. **Performance**: Complex animations can impact performance. Test on lower-powered devices.

2. **Fullscreen video**: Keep `article_blur_enabled = false` unless you specifically want frosted-glass content. Even a no-op blur can hurt fullscreen video playback on some browsers.

   Also avoid adding global `::backdrop` blur, since that affects fullscreen media as well as dialogs.

3. **Accessibility**: Ensure backgrounds don't interfere with content readability. Use `pointer-events: none` (applied automatically).

4. **Z-Index**: Use negative values to place backgrounds behind content. Positive values overlay content.

5. **Web Components**: Custom elements like `<snow-fall>` provide encapsulated, reusable effects.

6. **Reduced Motion**: Consider respecting `prefers-reduced-motion` in your CSS:

```css
@media (prefers-reduced-motion: reduce) {
  .background-layer * {
    animation: none !important;
  }
}
```

### Readable Content Over Decorations

Use article wrapper styling to keep content readable without forcing blur:

```toml
[markata-go.theme.background]
enabled = true
article_bg = "color-mix(in srgb, var(--color-background) 88%, transparent)"
article_shadow = "0 12px 32px rgba(0, 0, 0, 0.12)"
article_border = "1px solid color-mix(in srgb, var(--color-border) 60%, transparent)"
article_radius = "1rem"

backgrounds = [
  { html = '<div class="stars"></div>', z_index = -10 },
]
```

If you do want a frosted-glass effect, enable it explicitly:

```toml
[markata-go.theme.background]
enabled = true
article_bg = "rgba(255, 255, 255, 0.72)"
article_blur_enabled = true
article_blur = "12px"
```

### Example: Particle Background

Using [particles.js](https://vincentgarreau.com/particles.js/):

```toml
[markata-go.theme.background]
enabled = true

backgrounds = [
  { html = '<div id="particles-js"></div>' },
]

scripts = [
  "https://cdn.jsdelivr.net/particles.js/2.0.0/particles.min.js",
  "/static/js/particles-config.js",
]

css = '''
#particles-js {
  position: absolute;
  inset: 0;
}
'''
```

---

## Seasonal Theme Calendar

Automatically apply different themes based on the date. Perfect for seasonal decorations, holidays, or promotional periods.

### Basic Configuration

Enable date-based theme switching:

```toml
[markata-go.theme_calendar]
enabled = true

[[markata-go.theme_calendar.rules]]
name = "Christmas Season"
start_date = "12-15"
end_date = "12-26"
palette = "christmas"

[[markata-go.theme_calendar.rules]]
name = "Winter Frost"
start_date = "12-01"
end_date = "02-28"
palette = "winter-frost"
```

During December 15-26, the site uses the "christmas" palette. From December 1 to February 28, it uses "winter-frost" (excluding the more specific Christmas period).

### How It Works

1. **First match wins**: Rules are evaluated in order. The first matching rule is applied.
2. **MM-DD format**: Dates use month-day format (no year), so rules repeat annually.
3. **Year boundary support**: Ranges like `12-01` to `02-28` correctly span December through February.
4. **Runs early**: The calendar plugin runs at the Configure stage before other theme plugins, ensuring palette overrides are applied correctly.

### Rule Configuration

Each rule can override various theme settings:

```toml
[[markata-go.theme_calendar.rules]]
name = "Spooky October"
start_date = "10-01"
end_date = "10-31"
palette = "halloween"            # Single palette for both modes
# Or use separate light/dark palettes:
# palette_light = "halloween-light"
# palette_dark = "halloween-dark"

# CSS variable overrides
[markata-go.theme_calendar.rules.variables]
"--color-accent" = "#ff6600"
"--color-link" = "#9b59b6"

# Custom CSS file for the season
custom_css = "halloween.css"
```

### Rule Options Reference

| Option | Type | Description |
|--------|------|-------------|
| `name` | string | Display name for the rule (used in logs and CLI) |
| `start_date` | string | Start date in MM-DD format (e.g., "12-01") |
| `end_date` | string | End date in MM-DD format (e.g., "02-28") |
| `palette` | string | Palette to use (overrides both modes) |
| `palette_light` | string | Light mode palette (if different from dark) |
| `palette_dark` | string | Dark mode palette |
| `custom_css` | string | Additional CSS file for this period |
| `variables` | table | CSS variable overrides |
| `background` | table | Background decoration override (see below) |
| `font` | table | Font configuration override |

### Background Overrides

Apply seasonal background decorations:

```toml
[[markata-go.theme_calendar.rules]]
name = "Winter Wonderland"
start_date = "12-01"
end_date = "02-28"
palette = "winter-frost"

[markata-go.theme_calendar.rules.background]
enabled = true
css = """
.snowflake {
  position: absolute;
  color: white;
  animation: fall linear infinite;
}
"""

[[markata-go.theme_calendar.rules.background.backgrounds]]
html = '<snow-fall count="100"></snow-fall>'
z_index = -5
```

### Font Overrides

Use a different font for special periods:

```toml
[[markata-go.theme_calendar.rules]]
name = "Christmas Season"
start_date = "12-15"
end_date = "12-26"
palette = "christmas"

[markata-go.theme_calendar.rules.font]
family = "Mountains of Christmas"
heading_family = "Snowburst One"
google_fonts = ["Mountains of Christmas", "Snowburst One"]
```

### Multiple Rules Example

Create a full year of seasonal themes:

```toml
[markata-go.theme_calendar]
enabled = true
default_palette = "catppuccin-mocha"  # Fallback when no rule matches

# Spring (March 20 - June 20)
[[markata-go.theme_calendar.rules]]
name = "Spring"
start_date = "03-20"
end_date = "06-20"
palette = "spring-garden"

# Summer (June 21 - September 22)
[[markata-go.theme_calendar.rules]]
name = "Summer"
start_date = "06-21"
end_date = "09-22"
palette = "summer-sunset"

# Halloween (October 15 - November 1)
[[markata-go.theme_calendar.rules]]
name = "Halloween"
start_date = "10-15"
end_date = "11-01"
palette = "spooky"

# Fall (September 23 - December 20, but after Halloween ends)
[[markata-go.theme_calendar.rules]]
name = "Fall"
start_date = "11-02"
end_date = "12-20"
palette = "autumn-leaves"

# Also need early fall before Halloween
[[markata-go.theme_calendar.rules]]
name = "Early Fall"
start_date = "09-23"
end_date = "10-14"
palette = "autumn-leaves"

# Winter (December 21 - March 19)
[[markata-go.theme_calendar.rules]]
name = "Winter"
start_date = "12-21"
end_date = "03-19"
palette = "winter-frost"
```

### CLI Commands

**List all calendar rules:**

```bash
markata-go theme calendar list
```

Output:
```
Theme Calendar Rules (3 configured)
============================================================

Christmas Season [ACTIVE]
  Date Range: 12-15 to 12-26
  Palette: christmas

Winter Frost
  Date Range: 12-01 to 02-28
  Palette: winter-frost

Halloween
  Date Range: 10-15 to 11-01
  Palette: spooky
```

**Preview theme for a specific date:**

```bash
# Check today's theme
markata-go theme calendar preview

# Check what theme applies on December 25
markata-go theme calendar preview 12-25

# Check New Year's Day
markata-go theme calendar preview 01-01
```

Output:
```
Checking theme for date: 12-25
----------------------------------------

Matching Rule: Christmas Season
Date Range: 12-15 to 12-26

Theme Overrides:
  Palette: christmas
  Font Family: Mountains of Christmas
```

### Tips

1. **Order matters**: Put more specific rules before general ones. Christmas should come before Winter.

2. **Testing**: Use `markata-go theme calendar preview MM-DD` to test any date without waiting for that day.

3. **Smooth transitions**: Consider overlapping date ranges with similar themes for gradual transitions.

4. **Performance**: Each rule is checked in order; keep the number of rules reasonable.

5. **Combine with switcher**: The calendar sets a default theme, but users can still override via the palette switcher if enabled.

---

## Font Configuration

markata-go provides flexible font configuration to customize your site's typography without writing CSS.

### Quick Start

Add custom fonts via Google Fonts:

```toml
[markata-go.theme.font]
google_fonts = ["Inter", "Fira Code"]
family = "'Inter', sans-serif"
code_family = "'Fira Code', monospace"
```

### Font Options

| Option | Description | Default |
|--------|-------------|---------|
| `family` | Body text font | System fonts |
| `heading_family` | Heading font (inherits from `family` if not set) | Same as `family` |
| `code_family` | Code/monospace font | System monospace |
| `size` | Base font size | `16px` |
| `line_height` | Base line height | `1.6` |
| `google_fonts` | Array of Google Fonts to load | `[]` |
| `custom_urls` | Array of custom font CSS URLs | `[]` |

### Using Google Fonts

Specify fonts to load from Google Fonts:

```toml
[markata-go.theme.font]
# Load these fonts from Google Fonts
google_fonts = ["Inter", "Playfair Display", "JetBrains Mono"]

# Use them in your font families
family = "'Inter', sans-serif"
heading_family = "'Playfair Display', serif"
code_family = "'JetBrains Mono', monospace"
```

The `google_fonts` array automatically generates the Google Fonts CSS URL with weights 400, 500, 600, and 700.

### Using Custom Fonts

Load fonts from any URL:

```toml
[markata-go.theme.font]
custom_urls = [
  "https://fonts.example.com/my-font.css",
  "/fonts/local-font.css"
]
family = "'My Custom Font', sans-serif"
```

### Typography Variables

Font configuration generates CSS custom properties that you can use in your custom CSS:

| Variable | Description |
|----------|-------------|
| `--font-family` | Body text font stack |
| `--font-heading` | Heading font stack |
| `--font-code` | Code/monospace font stack |
| `--font-size` | Base font size |
| `--line-height` | Base line height |

### Complete Example

```toml
[markata-go.theme]
palette = "catppuccin-mocha"

[markata-go.theme.font]
# Google Fonts to load
google_fonts = ["Source Sans Pro", "Source Serif Pro", "Source Code Pro"]

# Font assignments
family = "'Source Sans Pro', sans-serif"
heading_family = "'Source Serif Pro', serif"
code_family = "'Source Code Pro', monospace"

# Typography settings
size = "18px"
line_height = "1.7"
```

### Using Self-Hosted Fonts

For better performance and privacy, you can self-host fonts:

1. Download font files to `static/fonts/`
2. Create a CSS file defining `@font-face` rules
3. Reference it in `custom_urls`

```css
/* static/fonts/fonts.css */
@font-face {
  font-family: 'MyFont';
  src: url('/fonts/MyFont-Regular.woff2') format('woff2');
  font-weight: 400;
}

@font-face {
  font-family: 'MyFont';
  src: url('/fonts/MyFont-Bold.woff2') format('woff2');
  font-weight: 700;
}
```

```toml
[markata-go.theme.font]
custom_urls = ["/fonts/fonts.css"]
family = "'MyFont', sans-serif"
```

---

## Best Practices

### 1. Start with a Palette

Don't write CSS from scratch. Pick the closest palette and customize from there.

### 2. Use CSS Variables

Override `--color-*` variables instead of hardcoding colors. This ensures consistency and makes future changes easier.

### 3. Keep Customizations Minimal

The less you customize, the easier upgrades will be. Only override what you need.

### 4. Test Dark Mode

If you customize colors, test both light and dark mode to ensure readability.

### 5. Check Accessibility

Use `markata-go palette check` to verify your color choices meet WCAG guidelines.

---

## Troubleshooting

### Styles Not Loading

1. Check that CSS files exist in `public/css/` after building
2. Verify your browser's network tab shows CSS loading
3. Clear browser cache and rebuild: `markata-go build`

### Custom CSS Not Applying

1. Ensure `custom_css` path is relative to `static/`
2. Check for CSS specificity issues (theme styles may override yours)
3. Use browser dev tools to inspect applied styles

### Template Not Found

1. Verify the template file exists in `templates/`
2. Check the filename matches exactly (case-sensitive)
3. Ensure frontmatter `template:` value includes `.html` extension

---

---

## Next Steps

Now that you've styled your site, here are recommended next steps:

**Customize your templates:**
- [Templates Guide](/docs/guides/templates/) - Modify HTML structure, add custom partials, and use template inheritance

**Organize your content:**
- [Feeds Guide](/docs/guides/feeds/) - Create filtered collections, archives, and tag pages

**Deploy your site:**
- [Deployment Guide](/docs/guides/deployment/) - Deploy to GitHub Pages, Netlify, Vercel, or self-host

---

## See Also

- [Keyboard Navigation Guide](/docs/guides/keyboard-navigation/) - Comprehensive keyboard shortcuts for site navigation
- [Configuration Guide](/docs/guides/configuration/) - All configuration options
- [Templates Guide](/docs/guides/templates/) - Template syntax and customization
- [Frontmatter Guide](/docs/guides/frontmatter/) - Post-level configuration
- [Quick Reference](/docs/guides/quick-reference/) - Theme snippets and CLI commands
