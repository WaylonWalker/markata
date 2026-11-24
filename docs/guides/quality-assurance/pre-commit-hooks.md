---
title: "Pre-commit Hooks for Content Quality"
description: "Set up pre-commit hooks for frontmatter validation, Markdown linting, and link checking"
date: 2026-01-24
published: true
slug: /docs/guides/quality-assurance/pre-commit-hooks/
tags:
  - documentation
  - quality-assurance
  - pre-commit
  - linting
---


Pre-commit hooks automatically validate your content before each commit, catching issues early in the development cycle. This guide shows you how to set up comprehensive content quality checks for your markata-go site.

## Table of Contents

- [Installation](#installation)
- [Basic Configuration](#basic-configuration)
- [YAML Frontmatter Validation](#yaml-frontmatter-validation)
- [Markdown Linting](#markdown-linting)
- [Link Checking](#link-checking)
- [Image Alt Text Validation](#image-alt-text-validation)
- [Custom Validation Scripts](#custom-validation-scripts)
- [Complete Configuration Example](#complete-configuration-example)

---

## Installation

### Install pre-commit

```bash
# macOS (Homebrew)
brew install pre-commit

# Python (pip) - works on all platforms
pip install pre-commit

# Python (pipx) - isolated installation
pipx install pre-commit

# Conda
conda install -c conda-forge pre-commit
```

### Verify Installation

```bash
pre-commit --version
# pre-commit 3.7.0
```

### Initialize in Your Project

```bash
# Create initial configuration
pre-commit sample-config > .pre-commit-config.yaml

# Install git hooks
pre-commit install

# Optional: Install for commit messages too
pre-commit install --hook-type commit-msg
```

---

## Basic Configuration

Create `.pre-commit-config.yaml` in your site's root directory:

```yaml
# .pre-commit-config.yaml
default_language_version:
  python: python3

repos:
  # Basic file checks
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.6.0
    hooks:
      - id: trailing-whitespace
        args: [--markdown-linebreak-ext=md]
      - id: end-of-file-fixer
      - id: check-yaml
      - id: check-added-large-files
        args: [--maxkb=1024]
      - id: mixed-line-ending
        args: [--fix=lf]
```

### Running Pre-commit

```bash
# Run on staged files (automatic on commit)
pre-commit run

# Run on all files
pre-commit run --all-files

# Run specific hook
pre-commit run markdownlint --all-files

# Update hooks to latest versions
pre-commit autoupdate
```

---

## YAML Frontmatter Validation

markata-go uses YAML frontmatter for post metadata. Validating this frontmatter catches syntax errors before they break your build.

### Using yamllint

Add yamllint to your pre-commit config:

```yaml
repos:
  - repo: https://github.com/adrienverge/yamllint
    rev: v1.35.1
    hooks:
      - id: yamllint
        args: [--config-file, .yamllint.yml]
        files: \.(md|markdown)$
```

Create `.yamllint.yml` for Markdown frontmatter:

```yaml
# .yamllint.yml
extends: relaxed

rules:
  # Frontmatter doesn't need document start markers
  document-start: disable
  
  # Allow long lines in frontmatter (descriptions, etc.)
  line-length: disable
  
  # Allow trailing spaces (some editors add them)
  trailing-spaces: disable
  
  # Require consistent indentation
  indentation:
    spaces: 2
    indent-sequences: true
  
  # Ensure proper quoting
  quoted-strings:
    quote-type: any
    required: only-when-needed
  
  # Check for duplicate keys
  key-duplicates: enable
  
  # Ensure proper boolean values
  truthy:
    allowed-values: ['true', 'false', 'yes', 'no']
```

### Custom Frontmatter Validation Script

For more specific frontmatter validation, create a custom script:

```bash
#!/bin/bash
# scripts/validate-frontmatter.sh
# Validates that required frontmatter fields exist

set -e

required_fields=("title" "date" "published")
errors=0

for file in "$@"; do
    # Skip non-markdown files
    [[ "$file" != *.md ]] && continue
    
    # Extract frontmatter (between first two ---)
    frontmatter=$(sed -n '/^---$/,/^---$/p' "$file" | sed '1d;$d')
    
    if [ -z "$frontmatter" ]; then
        echo "ERROR: $file - No frontmatter found"
        ((errors++))
        continue
    fi
    
    for field in "${required_fields[@]}"; do
        if ! echo "$frontmatter" | grep -q "^${field}:"; then
            echo "ERROR: $file - Missing required field: $field"
            ((errors++))
        fi
    done
    
    # Check that published is boolean
    published=$(echo "$frontmatter" | grep "^published:" | cut -d: -f2 | tr -d ' ')
    if [ -n "$published" ] && [[ ! "$published" =~ ^(true|false)$ ]]; then
        echo "ERROR: $file - 'published' must be true or false, got: $published"
        ((errors++))
    fi
    
    # Check date format (YYYY-MM-DD)
    date_val=$(echo "$frontmatter" | grep "^date:" | cut -d: -f2 | tr -d ' ')
    if [ -n "$date_val" ] && [[ ! "$date_val" =~ ^[0-9]{4}-[0-9]{2}-[0-9]{2} ]]; then
        echo "WARNING: $file - Date format should be YYYY-MM-DD, got: $date_val"
    fi
done

if [ $errors -gt 0 ]; then
    echo ""
    echo "Found $errors frontmatter error(s)"
    exit 1
fi

echo "All frontmatter valid!"
```

Add the script to pre-commit:

```yaml
repos:
  - repo: local
    hooks:
      - id: validate-frontmatter
        name: Validate Frontmatter
        entry: bash scripts/validate-frontmatter.sh
        language: system
        files: \.(md|markdown)$
        pass_filenames: true
```

---

## Markdown Linting

Consistent Markdown formatting improves maintainability and ensures reliable rendering.

### Using markdownlint-cli

```yaml
repos:
  - repo: https://github.com/igorshubovych/markdownlint-cli
    rev: v0.42.0
    hooks:
      - id: markdownlint
        args: [--config, .markdownlint.json, --fix]
```

Create `.markdownlint.json`:

```json
{
  "default": true,
  
  "MD001": true,
  "MD003": { "style": "atx" },
  "MD004": { "style": "dash" },
  "MD007": { "indent": 2 },
  "MD009": { "br_spaces": 2 },
  "MD010": true,
  "MD012": { "maximum": 2 },
  
  "MD013": false,
  
  "MD022": { "lines_above": 1, "lines_below": 1 },
  "MD024": { "siblings_only": true },
  "MD025": { "front_matter_title": "^\\s*title\\s*[:=]" },
  
  "MD026": { "punctuation": ".,;:!" },
  
  "MD033": {
    "allowed_elements": [
      "details", "summary", "kbd", "br", "sup", "sub",
      "img", "video", "audio", "source", "iframe"
    ]
  },
  
  "MD034": true,
  
  "MD036": false,
  
  "MD040": true,
  "MD041": false,
  "MD046": { "style": "fenced" },
  "MD048": { "style": "backtick" }
}
```

### Rule Reference

Common rules you may want to configure:

| Rule | Description | Recommended |
|------|-------------|-------------|
| MD013 | Line length | Disable for prose |
| MD033 | No inline HTML | Allow specific tags |
| MD041 | First line H1 | Disable (frontmatter) |
| MD025 | Single H1 | Enable with frontmatter exception |
| MD024 | No duplicate headings | Enable siblings only |

### Per-File Overrides

Disable rules for specific files using `.markdownlintignore`:

```
# .markdownlintignore
# Ignore generated files
public/
node_modules/

# Ignore specific files
CHANGELOG.md
```

Or use inline comments in Markdown:

```markdown
<!-- markdownlint-disable MD033 -->
<details>
<summary>Click to expand</summary>

Content here...

</details>
<!-- markdownlint-enable MD033 -->
```

---

## Link Checking

Broken links hurt user experience and SEO. Check links before committing.

### Using markdown-link-check

```yaml
repos:
  - repo: https://github.com/tcort/markdown-link-check
    rev: v3.12.2
    hooks:
      - id: markdown-link-check
        args: [--config, .markdown-link-check.json]
```

Create `.markdown-link-check.json`:

```json
{
  "ignorePatterns": [
    { "pattern": "^https://localhost" },
    { "pattern": "^https://127\\.0\\.0\\.1" },
    { "pattern": "^#" }
  ],
  "replacementPatterns": [
    {
      "pattern": "^/",
      "replacement": "{{BASEURL}}/"
    }
  ],
  "httpHeaders": [
    {
      "urls": ["https://github.com"],
      "headers": {
        "Accept": "text/html"
      }
    }
  ],
  "timeout": "10s",
  "retryOn429": true,
  "retryCount": 3,
  "fallbackRetryDelay": "5s",
  "aliveStatusCodes": [200, 206, 301, 302, 307, 308]
}
```

### Using lychee (Faster Alternative)

[lychee](https://github.com/lycheeverse/lychee) is a fast, async link checker written in Rust:

```yaml
repos:
  - repo: local
    hooks:
      - id: lychee
        name: Check links with lychee
        entry: lychee
        language: system
        files: \.(md|html)$
        args:
          - --config=.lychee.toml
          - --no-progress
```

Create `.lychee.toml`:

```toml
# .lychee.toml

# Maximum number of concurrent requests
max_concurrency = 32

# Timeout for requests
timeout = 20

# User agent
user_agent = "lychee/0.15"

# Accept codes
accept = [200, 204, 301, 302, 307, 308]

# Exclude patterns
exclude = [
    "^https://localhost",
    "^https://127\\.0\\.0\\.1",
    "^mailto:",
    "^tel:",
]

# Exclude paths
exclude_path = [
    "node_modules",
    "public",
    ".git",
]

# Skip private IPs
skip_missing = false

# Include fragments
include_fragments = true
```

### Link Checking Best Practices

1. **Cache results** - Link checking is slow; cache in CI
2. **Allow retries** - Some sites rate-limit requests
3. **Exclude known-good patterns** - Skip localhost, anchors
4. **Run separately** - Consider running link checks only in CI, not on every commit

---

## Image Alt Text Validation

Alt text is essential for accessibility. Validate that all images have meaningful descriptions.

### Custom Alt Text Checker

```bash
#!/bin/bash
# scripts/check-alt-text.sh
# Ensures all images have alt text

set -e

errors=0

for file in "$@"; do
    [[ "$file" != *.md ]] && continue
    
    # Find images without alt text: ![](url) or ![ ](url)
    # Correct format: ![alt text](url)
    
    # Check for empty alt text
    if grep -Pn '!\[\s*\]\(' "$file"; then
        echo "ERROR: $file - Found image(s) with empty alt text"
        ((errors++))
    fi
    
    # Check for placeholder alt text
    if grep -Pin '!\[(image|img|photo|picture|screenshot)\]\(' "$file"; then
        echo "WARNING: $file - Found image(s) with placeholder alt text"
    fi
done

if [ $errors -gt 0 ]; then
    echo ""
    echo "Found $errors image(s) missing alt text"
    echo "Add descriptive alt text: ![description of image](url)"
    exit 1
fi

echo "All images have alt text!"
```

Add to pre-commit:

```yaml
repos:
  - repo: local
    hooks:
      - id: check-alt-text
        name: Check image alt text
        entry: bash scripts/check-alt-text.sh
        language: system
        files: \.(md|markdown)$
        pass_filenames: true
```

### Using remark-lint

For more comprehensive Markdown validation including alt text:

```yaml
repos:
  - repo: local
    hooks:
      - id: remark-lint
        name: Remark lint
        entry: npx remark
        language: node
        files: \.(md|markdown)$
        args: [--frail, --quiet]
        additional_dependencies:
          - remark-cli
          - remark-preset-lint-recommended
          - remark-lint-no-empty-image-alt
```

Create `.remarkrc.js`:

```javascript
// .remarkrc.js
module.exports = {
  plugins: [
    'preset-lint-recommended',
    'lint-no-empty-image-alt',
    ['lint-no-undefined-references', false], // Allow [[wikilinks]]
  ],
};
```

---

## Custom Validation Scripts

### Validate Slug Uniqueness

Prevent duplicate slugs that cause build conflicts:

```bash
#!/bin/bash
# scripts/check-unique-slugs.sh

set -e

# Extract all slugs from frontmatter
declare -A slugs

for file in $(find . -name "*.md" -not -path "./node_modules/*" -not -path "./public/*"); do
    slug=$(sed -n '/^---$/,/^---$/p' "$file" | grep "^slug:" | cut -d: -f2- | tr -d ' "'"'"'')
    
    if [ -n "$slug" ]; then
        if [ -n "${slugs[$slug]}" ]; then
            echo "ERROR: Duplicate slug '$slug'"
            echo "  - ${slugs[$slug]}"
            echo "  - $file"
            exit 1
        fi
        slugs[$slug]="$file"
    fi
done

echo "All slugs are unique!"
```

### Validate Tags

Ensure consistent tag usage:

```bash
#!/bin/bash
# scripts/validate-tags.sh

# Allowed tags (customize for your site)
allowed_tags=(
    "documentation"
    "tutorial"
    "guide"
    "reference"
    "blog"
    "announcement"
)

errors=0

for file in "$@"; do
    [[ "$file" != *.md ]] && continue
    
    # Extract tags from frontmatter
    tags=$(sed -n '/^tags:/,/^[a-z]/p' "$file" | grep "^\s*-" | sed 's/^\s*-\s*//')
    
    while IFS= read -r tag; do
        tag=$(echo "$tag" | tr -d ' ')
        [ -z "$tag" ] && continue
        
        if [[ ! " ${allowed_tags[*]} " =~ " ${tag} " ]]; then
            echo "WARNING: $file - Unknown tag: $tag"
            echo "  Allowed tags: ${allowed_tags[*]}"
        fi
    done <<< "$tags"
done

exit 0  # Warnings only, don't block commit
```

### Check Required Sections

Ensure documentation has required sections:

```bash
#!/bin/bash
# scripts/check-doc-structure.sh

for file in "$@"; do
    [[ "$file" != docs/*.md ]] && continue
    
    # Check for table of contents
    if ! grep -q "## Table of Contents" "$file"; then
        echo "WARNING: $file - Missing Table of Contents"
    fi
    
    # Check for See Also section
    if ! grep -q "## See Also" "$file"; then
        echo "INFO: $file - Consider adding a See Also section"
    fi
done
```

---

## Complete Configuration Example

Here's a comprehensive `.pre-commit-config.yaml` for a markata-go site:

```yaml
# .pre-commit-config.yaml
# Complete content quality configuration for markata-go sites

default_language_version:
  python: python3

repos:
  # ============================================
  # Basic file checks
  # ============================================
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.6.0
    hooks:
      - id: trailing-whitespace
        args: [--markdown-linebreak-ext=md]
      - id: end-of-file-fixer
      - id: check-yaml
      - id: check-json
      - id: check-toml
      - id: check-added-large-files
        args: [--maxkb=1024]
      - id: mixed-line-ending
        args: [--fix=lf]
      - id: check-merge-conflict
      - id: detect-private-key

  # ============================================
  # YAML/Frontmatter validation
  # ============================================
  - repo: https://github.com/adrienverge/yamllint
    rev: v1.35.1
    hooks:
      - id: yamllint
        args: [--config-file, .yamllint.yml]
        files: \.(md|yaml|yml)$

  # ============================================
  # Markdown linting
  # ============================================
  - repo: https://github.com/igorshubovych/markdownlint-cli
    rev: v0.42.0
    hooks:
      - id: markdownlint
        args: [--config, .markdownlint.json, --fix]
        exclude: ^(CHANGELOG|node_modules/)

  # ============================================
  # Link checking (optional - can be slow)
  # ============================================
  # - repo: https://github.com/tcort/markdown-link-check
  #   rev: v3.12.2
  #   hooks:
  #     - id: markdown-link-check
  #       args: [--config, .markdown-link-check.json]

  # ============================================
  # Custom validations
  # ============================================
  - repo: local
    hooks:
      # Validate required frontmatter fields
      - id: validate-frontmatter
        name: Validate Frontmatter
        entry: bash scripts/validate-frontmatter.sh
        language: system
        files: \.(md|markdown)$
        pass_filenames: true

      # Check image alt text
      - id: check-alt-text
        name: Check image alt text
        entry: bash scripts/check-alt-text.sh
        language: system
        files: \.(md|markdown)$
        pass_filenames: true

      # Validate build (optional - can be slow)
      # - id: markata-build
      #   name: Test build
      #   entry: markata-go build --dry-run
      #   language: system
      #   pass_filenames: false
      #   stages: [push]
```

### Directory Structure

After setup, your project should have:

```
your-site/
├── .pre-commit-config.yaml     # Pre-commit configuration
├── .yamllint.yml               # YAML linting rules
├── .markdownlint.json          # Markdown linting rules
├── .markdown-link-check.json   # Link checker config (optional)
├── scripts/
│   ├── validate-frontmatter.sh
│   └── check-alt-text.sh
├── docs/
│   └── ...your content...
└── markata-go.toml
```

---

## Troubleshooting

### Common Issues

**Hook not running**

```bash
# Reinstall hooks
pre-commit uninstall
pre-commit install
```

**Skipping hooks temporarily**

```bash
# Skip all hooks (use sparingly!)
git commit --no-verify -m "WIP: work in progress"

# Skip specific hook
SKIP=markdownlint git commit -m "Quick fix"
```

**Hooks too slow**

```bash
# Run only on changed files (default behavior)
pre-commit run

# Cache is stored in ~/.cache/pre-commit/
```

**False positives**

Use inline comments to disable specific rules:

```markdown
<!-- markdownlint-disable MD033 -->
<custom-element>content</custom-element>
<!-- markdownlint-enable MD033 -->
```

Or exclude files in `.markdownlintignore`.

---

## See Also

- [GitHub Actions](github-actions/) - CI/CD quality checks
- [Custom Rules](custom-rules/) - Configure linting rules
- [Troubleshooting](troubleshooting/) - Common issues and solutions
