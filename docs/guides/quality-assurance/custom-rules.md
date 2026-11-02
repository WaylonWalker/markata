---
title: "Custom Linting Rules"
description: "Configure and create custom linting rules for markdownlint, frontmatter validation, and project-specific style guides"
date: 2026-01-24
published: true
slug: /docs/guides/quality-assurance/custom-rules/
tags:
  - documentation
  - quality-assurance
  - linting
  - configuration
---


This guide covers how to configure linting rules for your specific needs, including markdownlint customization, custom frontmatter validation, and creating project-specific style guides.

## Table of Contents

- [markdownlint Configuration](#markdownlint-configuration)
- [Custom markdownlint Rules](#custom-markdownlint-rules)
- [Frontmatter Validation Rules](#frontmatter-validation-rules)
- [Project Style Guides](#project-style-guides)
- [Rule Inheritance and Overrides](#rule-inheritance-and-overrides)
- [Creating Custom Validators](#creating-custom-validators)

---

## markdownlint Configuration

### Configuration File Formats

markdownlint supports multiple configuration formats:

**JSON (`.markdownlint.json`)**

```json
{
  "default": true,
  "MD013": false,
  "MD033": {
    "allowed_elements": ["details", "summary", "kbd"]
  }
}
```

**YAML (`.markdownlint.yaml`)**

```yaml
default: true
MD013: false
MD033:
  allowed_elements:
    - details
    - summary
    - kbd
```

**JavaScript (`.markdownlint.cjs`)**

```javascript
module.exports = {
  default: true,
  MD013: false,
  MD033: {
    allowed_elements: ['details', 'summary', 'kbd'],
  },
};
```

### Complete Rule Reference

Here's a comprehensive configuration covering all common rules:

```json
{
  "default": true,

  "MD001": true,
  "MD003": { "style": "atx" },
  "MD004": { "style": "dash" },
  "MD005": true,
  "MD007": { "indent": 2, "start_indented": false },
  "MD009": { "br_spaces": 2, "list_item_empty_lines": false },
  "MD010": { "code_blocks": true, "spaces_per_tab": 2 },
  "MD011": true,
  "MD012": { "maximum": 2 },

  "MD013": {
    "line_length": 120,
    "heading_line_length": 80,
    "code_block_line_length": 120,
    "code_blocks": false,
    "tables": false
  },

  "MD014": true,

  "MD018": true,
  "MD019": true,
  "MD020": true,
  "MD021": true,
  "MD022": { "lines_above": 1, "lines_below": 1 },
  "MD023": true,
  "MD024": { "siblings_only": true },
  "MD025": { "level": 1, "front_matter_title": "^\\s*title\\s*[:=]" },
  "MD026": { "punctuation": ".,;:!" },
  "MD027": true,
  "MD028": true,

  "MD029": { "style": "ordered" },
  "MD030": { "ul_single": 1, "ol_single": 1, "ul_multi": 1, "ol_multi": 1 },
  "MD031": { "list_items": true },
  "MD032": true,

  "MD033": {
    "allowed_elements": [
      "a", "abbr", "audio", "b", "br", "caption",
      "cite", "code", "col", "colgroup", "dd", "del",
      "details", "dfn", "div", "dl", "dt", "em",
      "figcaption", "figure", "h1", "h2", "h3", "h4",
      "h5", "h6", "hr", "i", "iframe", "img", "ins",
      "kbd", "li", "mark", "ol", "p", "picture", "pre",
      "q", "s", "samp", "small", "source", "span",
      "strong", "sub", "summary", "sup", "table",
      "tbody", "td", "tfoot", "th", "thead", "tr",
      "u", "ul", "var", "video"
    ]
  },

  "MD034": true,
  "MD035": { "style": "---" },
  "MD036": { "punctuation": ".,;:!?" },
  "MD037": true,
  "MD038": true,
  "MD039": true,
  "MD040": true,
  "MD041": false,
  "MD042": true,
  "MD043": false,
  "MD044": {
    "names": ["markata-go", "GitHub", "GitLab", "JavaScript", "TypeScript"],
    "code_blocks": false
  },
  "MD045": true,
  "MD046": { "style": "fenced" },
  "MD047": true,
  "MD048": { "style": "backtick" },
  "MD049": { "style": "underscore" },
  "MD050": { "style": "asterisk" },
  "MD051": true,
  "MD052": true,
  "MD053": true
}
```

### Rule Categories

| Category | Rules | Description |
|----------|-------|-------------|
| Headings | MD001-MD003, MD018-MD025 | Heading structure and style |
| Lists | MD004-MD007, MD029-MD032 | List formatting |
| Whitespace | MD009-MD012, MD027-MD028 | Spacing and blank lines |
| Code | MD014, MD031, MD038, MD040, MD046, MD048 | Code blocks and inline code |
| Links | MD034, MD039, MD042, MD051-MD053 | Link formatting |
| Emphasis | MD036-MD037, MD049-MD050 | Bold/italic style |

---

## Custom markdownlint Rules

### Creating a Custom Rule

Custom rules are JavaScript modules. Create `custom-rules/no-todo-comments.js`:

```javascript
// custom-rules/no-todo-comments.js
module.exports = {
  names: ['no-todo-comments'],
  description: 'Disallow TODO comments in content',
  tags: ['content', 'todo'],
  function: function rule(params, onError) {
    params.tokens.forEach((token) => {
      if (token.type === 'inline') {
        const todoMatch = token.content.match(/\bTODO\b/i);
        if (todoMatch) {
          onError({
            lineNumber: token.lineNumber,
            detail: 'Remove TODO comment before publishing',
            context: token.content,
          });
        }
      }
    });
  },
};
```

### Registering Custom Rules

In `.markdownlint.cjs`:

```javascript
const noTodoComments = require('./custom-rules/no-todo-comments');
const requireDescription = require('./custom-rules/require-description');

module.exports = {
  default: true,
  MD013: false,
  customRules: [noTodoComments, requireDescription],
  'no-todo-comments': true,
  'require-description': { minLength: 50 },
};
```

### Example Custom Rules

**Require minimum content length:**

```javascript
// custom-rules/min-content-length.js
module.exports = {
  names: ['min-content-length'],
  description: 'Enforce minimum content length',
  tags: ['content', 'length'],
  function: function rule(params, onError) {
    const minLength = params.config.minLength || 300;
    const content = params.lines.join('\n');
    
    // Remove frontmatter
    const bodyMatch = content.match(/^---[\s\S]*?---\n([\s\S]*)$/);
    const body = bodyMatch ? bodyMatch[1] : content;
    
    // Remove code blocks and count words
    const textOnly = body.replace(/```[\s\S]*?```/g, '').replace(/`[^`]+`/g, '');
    const wordCount = textOnly.split(/\s+/).filter(Boolean).length;
    
    if (wordCount < minLength) {
      onError({
        lineNumber: 1,
        detail: `Content has ${wordCount} words, minimum is ${minLength}`,
      });
    }
  },
};
```

**Check heading hierarchy:**

```javascript
// custom-rules/heading-hierarchy.js
module.exports = {
  names: ['heading-hierarchy'],
  description: 'Ensure headings follow proper hierarchy',
  tags: ['headings', 'structure'],
  function: function rule(params, onError) {
    let lastLevel = 0;
    
    params.tokens
      .filter((token) => token.type === 'heading_open')
      .forEach((token) => {
        const level = parseInt(token.tag.substring(1), 10);
        
        if (lastLevel > 0 && level > lastLevel + 1) {
          onError({
            lineNumber: token.lineNumber,
            detail: `Heading level jumped from h${lastLevel} to h${level}`,
            context: token.line,
          });
        }
        
        lastLevel = level;
      });
  },
};
```

---

## Frontmatter Validation Rules

### YAML Schema Validation

Create a JSON Schema for frontmatter validation. Create `frontmatter-schema.json`:

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "markata-go Frontmatter",
  "type": "object",
  "required": ["title", "date", "published"],
  "properties": {
    "title": {
      "type": "string",
      "minLength": 1,
      "maxLength": 100
    },
    "description": {
      "type": "string",
      "minLength": 50,
      "maxLength": 160
    },
    "date": {
      "type": "string",
      "pattern": "^\\d{4}-\\d{2}-\\d{2}$"
    },
    "published": {
      "type": "boolean"
    },
    "slug": {
      "type": "string",
      "pattern": "^/[a-z0-9-/]+/?$"
    },
    "tags": {
      "type": "array",
      "items": {
        "type": "string",
        "pattern": "^[a-z0-9-]+$"
      },
      "minItems": 1,
      "uniqueItems": true
    },
    "author": {
      "type": "string"
    },
    "image": {
      "type": "string",
      "format": "uri-reference"
    },
    "draft": {
      "type": "boolean"
    }
  },
  "additionalProperties": true
}
```

### Validation Script

Create `scripts/validate-frontmatter-schema.sh`:

```bash
#!/bin/bash
# scripts/validate-frontmatter-schema.sh
# Validates frontmatter against JSON Schema

set -e

# Requires: pip install check-jsonschema pyyaml

schema_file="frontmatter-schema.json"
errors=0

for file in "$@"; do
    [[ "$file" != *.md ]] && continue
    
    # Extract frontmatter
    frontmatter=$(sed -n '/^---$/,/^---$/p' "$file" | sed '1d;$d')
    
    if [ -z "$frontmatter" ]; then
        echo "WARNING: $file - No frontmatter found"
        continue
    fi
    
    # Convert to JSON and validate
    echo "$frontmatter" | python3 -c "
import sys, yaml, json
data = yaml.safe_load(sys.stdin.read())
print(json.dumps(data))
" > /tmp/frontmatter.json
    
    if ! check-jsonschema --schemafile "$schema_file" /tmp/frontmatter.json 2>/dev/null; then
        echo "ERROR: $file - Frontmatter validation failed"
        check-jsonschema --schemafile "$schema_file" /tmp/frontmatter.json 2>&1 | sed 's/^/  /'
        ((errors++))
    fi
done

if [ $errors -gt 0 ]; then
    echo ""
    echo "Found $errors frontmatter validation error(s)"
    exit 1
fi

echo "All frontmatter valid!"
```

### Pre-commit Hook

```yaml
repos:
  - repo: local
    hooks:
      - id: validate-frontmatter-schema
        name: Validate Frontmatter Schema
        entry: bash scripts/validate-frontmatter-schema.sh
        language: system
        files: \.(md|markdown)$
        pass_filenames: true
        additional_dependencies:
          - check-jsonschema
          - pyyaml
```

---

## Project Style Guides

### Documentation Style Guide

Create a comprehensive style configuration for documentation projects:

```json
{
  "default": true,
  
  "MD003": { "style": "atx" },
  "MD004": { "style": "dash" },
  "MD007": { "indent": 2 },
  
  "MD013": false,
  
  "MD022": { "lines_above": 1, "lines_below": 1 },
  "MD024": { "siblings_only": true },
  "MD025": { "front_matter_title": "^\\s*title\\s*[:=]" },
  
  "MD033": {
    "allowed_elements": [
      "details", "summary", "kbd", "br", "sup", "sub",
      "img", "video", "audio", "source", "iframe",
      "table", "thead", "tbody", "tr", "th", "td"
    ]
  },
  
  "MD036": false,
  "MD041": false,
  
  "MD044": {
    "names": [
      "markata-go",
      "GitHub",
      "GitLab",
      "JavaScript",
      "TypeScript",
      "Go",
      "Markdown",
      "YAML",
      "TOML",
      "JSON"
    ],
    "code_blocks": false
  },
  
  "MD046": { "style": "fenced" },
  "MD048": { "style": "backtick" }
}
```

### Blog Style Guide

For blog-style content with more flexibility:

```json
{
  "default": true,
  
  "MD003": { "style": "atx" },
  "MD004": { "style": "dash" },
  
  "MD013": false,
  
  "MD022": { "lines_above": 1, "lines_below": 1 },
  "MD024": { "siblings_only": true },
  
  "MD033": false,
  
  "MD036": false,
  "MD041": false,
  
  "MD046": { "style": "fenced" },
  "MD048": { "style": "backtick" }
}
```

### API Documentation Style

For technical API documentation:

```json
{
  "default": true,
  
  "MD003": { "style": "atx" },
  "MD004": { "style": "dash" },
  "MD007": { "indent": 2 },
  
  "MD013": {
    "line_length": 100,
    "code_blocks": false,
    "tables": false
  },
  
  "MD022": { "lines_above": 1, "lines_below": 1 },
  "MD024": false,
  "MD025": { "front_matter_title": "^\\s*title\\s*[:=]" },
  
  "MD033": {
    "allowed_elements": ["br", "code", "pre"]
  },
  
  "MD036": false,
  "MD041": false,
  
  "MD040": true,
  "MD046": { "style": "fenced" },
  "MD048": { "style": "backtick" }
}
```

---

## Rule Inheritance and Overrides

### Base Configuration

Create a base configuration in `config/.markdownlint-base.json`:

```json
{
  "default": true,
  "MD013": false,
  "MD033": false,
  "MD041": false
}
```

### Extending Base Config

**In project root `.markdownlint.json`:**

```json
{
  "extends": "config/.markdownlint-base.json",
  "MD024": { "siblings_only": true },
  "MD044": {
    "names": ["MyProject"],
    "code_blocks": false
  }
}
```

### Directory-Specific Overrides

Create override files for specific directories:

**`docs/.markdownlint.json`:**

```json
{
  "extends": "../.markdownlint.json",
  "MD013": {
    "line_length": 100
  },
  "MD025": {
    "front_matter_title": "^\\s*title\\s*[:=]"
  }
}
```

**`blog/.markdownlint.json`:**

```json
{
  "extends": "../.markdownlint.json",
  "MD013": false,
  "MD033": false
}
```

### Inline Overrides

Disable rules for specific sections:

```markdown
<!-- markdownlint-disable MD033 -->
<details>
<summary>Click to expand</summary>

This content uses raw HTML.

</details>
<!-- markdownlint-enable MD033 -->
```

Disable rules for entire file:

```markdown
<!-- markdownlint-disable-file MD013 MD033 -->

# This File Has Custom Rules

Long lines are allowed here...
```

---

## Creating Custom Validators

### Python Validator

Create a reusable Python validator:

```python
#!/usr/bin/env python3
# scripts/validate_content.py
"""Content validation for markata-go sites."""

import argparse
import re
import sys
from pathlib import Path
from typing import List, Tuple

import yaml


def extract_frontmatter(content: str) -> Tuple[dict, str]:
    """Extract frontmatter and body from content."""
    match = re.match(r'^---\n(.*?)\n---\n(.*)$', content, re.DOTALL)
    if not match:
        return {}, content
    
    try:
        frontmatter = yaml.safe_load(match.group(1))
        return frontmatter or {}, match.group(2)
    except yaml.YAMLError:
        return {}, content


def check_required_fields(
    frontmatter: dict,
    required: List[str],
    file_path: Path
) -> List[str]:
    """Check for required frontmatter fields."""
    errors = []
    for field in required:
        if field not in frontmatter:
            errors.append(f"{file_path}: Missing required field '{field}'")
    return errors


def check_description_length(
    frontmatter: dict,
    min_length: int,
    max_length: int,
    file_path: Path
) -> List[str]:
    """Check description length for SEO."""
    errors = []
    description = frontmatter.get('description', '')
    
    if description and len(description) < min_length:
        errors.append(
            f"{file_path}: Description too short "
            f"({len(description)} < {min_length})"
        )
    
    if description and len(description) > max_length:
        errors.append(
            f"{file_path}: Description too long "
            f"({len(description)} > {max_length})"
        )
    
    return errors


def check_tags(
    frontmatter: dict,
    allowed_tags: List[str],
    file_path: Path
) -> List[str]:
    """Check that tags are from allowed list."""
    warnings = []
    tags = frontmatter.get('tags', [])
    
    for tag in tags:
        if allowed_tags and tag not in allowed_tags:
            warnings.append(
                f"{file_path}: Unknown tag '{tag}'"
            )
    
    return warnings


def check_images_have_alt(body: str, file_path: Path) -> List[str]:
    """Check that all images have alt text."""
    errors = []
    
    # Find images with empty alt text
    empty_alt = re.findall(r'!\[\s*\]\([^)]+\)', body)
    for img in empty_alt:
        errors.append(f"{file_path}: Image missing alt text: {img[:50]}...")
    
    return errors


def validate_file(
    file_path: Path,
    config: dict
) -> Tuple[List[str], List[str]]:
    """Validate a single file."""
    errors = []
    warnings = []
    
    content = file_path.read_text()
    frontmatter, body = extract_frontmatter(content)
    
    # Required fields
    required = config.get('required_fields', ['title', 'date', 'published'])
    errors.extend(check_required_fields(frontmatter, required, file_path))
    
    # Description length
    desc_min = config.get('description_min_length', 50)
    desc_max = config.get('description_max_length', 160)
    errors.extend(check_description_length(
        frontmatter, desc_min, desc_max, file_path
    ))
    
    # Tags
    allowed_tags = config.get('allowed_tags', [])
    warnings.extend(check_tags(frontmatter, allowed_tags, file_path))
    
    # Images
    if config.get('require_alt_text', True):
        errors.extend(check_images_have_alt(body, file_path))
    
    return errors, warnings


def main():
    parser = argparse.ArgumentParser(description='Validate content files')
    parser.add_argument('files', nargs='+', help='Files to validate')
    parser.add_argument(
        '--config', '-c',
        default='.content-lint.yaml',
        help='Configuration file'
    )
    parser.add_argument(
        '--strict',
        action='store_true',
        help='Treat warnings as errors'
    )
    args = parser.parse_args()
    
    # Load config
    config = {}
    config_path = Path(args.config)
    if config_path.exists():
        config = yaml.safe_load(config_path.read_text())
    
    all_errors = []
    all_warnings = []
    
    for file_arg in args.files:
        file_path = Path(file_arg)
        if not file_path.suffix == '.md':
            continue
        
        errors, warnings = validate_file(file_path, config)
        all_errors.extend(errors)
        all_warnings.extend(warnings)
    
    # Output results
    for error in all_errors:
        print(f"ERROR: {error}")
    
    for warning in all_warnings:
        print(f"WARNING: {warning}")
    
    # Exit code
    if all_errors or (args.strict and all_warnings):
        sys.exit(1)
    
    print(f"Validated {len(args.files)} files successfully")
    sys.exit(0)


if __name__ == '__main__':
    main()
```

### Configuration File

Create `.content-lint.yaml`:

```yaml
# .content-lint.yaml
required_fields:
  - title
  - date
  - published

description_min_length: 50
description_max_length: 160

require_alt_text: true

allowed_tags:
  - documentation
  - tutorial
  - guide
  - reference
  - blog
  - announcement
  - quality-assurance
```

### Pre-commit Integration

```yaml
repos:
  - repo: local
    hooks:
      - id: content-validator
        name: Validate Content
        entry: python scripts/validate_content.py
        language: python
        files: \.(md|markdown)$
        pass_filenames: true
        additional_dependencies:
          - pyyaml
```

---

## See Also

- [Pre-commit Hooks](pre-commit-hooks/) - Local quality checks
- [GitHub Actions](github-actions/) - CI/CD integration
- [GitLab CI](gitlab-ci/) - GitLab pipelines
- [Troubleshooting](troubleshooting/) - Common issues and solutions
