# Issue #002: Incomplete Spec File Implementation

## Summary

When implementing the specification, not all spec files are being processed and implemented. The implementation is missing coverage for several specification documents, resulting in incomplete feature implementations.

## Current Behavior

The implementation process reads and implements some spec files but misses others. For example:

- `spec/spec/THEMES.md` - Not implemented (themes, CSS variables, admonitions styling, theme CLI commands)
- Potentially other spec files may be partially or fully skipped

This results in:
1. Missing theme system functionality
2. Missing CSS customization features
3. Missing admonition styling
4. Missing `theme list`, `theme info`, `theme install`, `theme new` CLI commands
5. Incomplete feature parity with the specification

## Expected Behavior

All spec files listed in the README should be fully implemented:

| File | Status |
|------|--------|
| INSTALL.md | ? |
| SPEC.md | ? |
| CONFIG.md | ? |
| **THEMES.md** | **Missing** |
| LIFECYCLE.md | ? |
| FEEDS.md | ? |
| DEFAULT_PLUGINS.md | ? |
| PLUGINS.md | ? |
| DATA_MODEL.md | ? |
| CONTENT.md | ? |
| TEMPLATES.md | ? |
| OPTIONAL_PLUGINS.md | ? |
| HEAD_STYLE.md | ? |
| REDIRECTS.md | ? |
| AUTO_TITLE.md | ? |
| IMPLEMENTATION.md | ? |

## Missing Features from THEMES.md

### Theme System
- [ ] Theme directory structure (`themes/[name]/`)
- [ ] `theme.toml` metadata parsing
- [ ] Theme resolution order (project local -> project theme -> installed -> built-in -> default)
- [ ] Theme inheritance (`[theme.extends]`)

### Configuration
- [ ] `[markata-go.theme]` config section
- [ ] `[markata-go.theme.options]` for theme-specific options
- [ ] `[markata-go.theme.variables]` CSS variable overrides
- [ ] `custom_css` option

### Built-in Themes
- [ ] Default theme with templates and CSS
- [ ] Minimal theme (optional)
- [ ] CSS custom properties (colors, typography, spacing)

### Admonition Styles
- [ ] All admonition type styles (note, tip, warning, danger, etc.)
- [ ] Aside/sidebar styles
- [ ] Chat/conversation styles
- [ ] Collapsible admonitions
- [ ] Dark mode support for all admonitions

### Code Block Styles
- [ ] Base code styles (inline and blocks)
- [ ] Line numbers support
- [ ] Language labels
- [ ] Syntax highlighting theme options

### Template Requirements
- [ ] `base.html` with required blocks
- [ ] `post.html` template
- [ ] `feed.html` template
- [ ] `card.html` template
- [ ] Template filters (`theme_asset`, `asset_url`, etc.)

### CLI Commands
- [ ] `[name] theme list` - List available themes
- [ ] `[name] theme info [theme]` - Show theme details
- [ ] `[name] theme install [url]` - Install theme from URL
- [ ] `[name] theme new [name]` - Scaffold new theme

## Root Cause

The spec processing workflow needs to:
1. Enumerate ALL spec files in `spec/spec/` directory
2. Track which files have been processed
3. Verify each file's requirements have been implemented
4. Report any gaps in implementation

## Proposed Solution

### Option A: Spec Checklist Generation

Generate a checklist from all spec files before implementation:

```bash
# Discover all spec files
ls spec/spec/*.md

# Generate implementation checklist
[name] spec check --generate-checklist
```

### Option B: Spec Coverage Report

After implementation, generate a coverage report:

```bash
[name] spec coverage
# Output:
# SPEC.md: 45/50 requirements (90%)
# THEMES.md: 0/35 requirements (0%) <- MISSING
# ...
```

### Option C: Structured Spec Processing

Ensure spec processing follows a structured approach:
1. Parse spec README to get list of all spec files
2. Process each file in order
3. Track implementation status
4. Block completion until all specs addressed

## Implementation Priority

High - Missing entire feature areas results in an incomplete implementation that doesn't match the specification.

## Files to Audit

All spec files should be reviewed for implementation coverage:

```
spec/spec/
├── AUTO_TITLE.md
├── CONFIG.md
├── CONTENT.md
├── DATA_MODEL.md
├── DEFAULT_PLUGINS.md
├── FEEDS.md
├── HEAD_STYLE.md
├── IMPLEMENTATION.md
├── INSTALL.md
├── LIFECYCLE.md
├── OPTIONAL_PLUGINS.md
├── PLUGINS.md
├── REDIRECTS.md
├── SPEC.md
├── TEMPLATES.md
└── THEMES.md        <- Confirmed missing
```

## Related

- spec/README.md lists all specification files
- Each spec file defines MUST/SHOULD/MAY requirements
- tests.yaml should have test cases for each spec area

## Labels

`bug`, `spec`, `implementation-gap`, `high-priority`
