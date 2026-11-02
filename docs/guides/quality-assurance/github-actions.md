---
title: "GitHub Actions for Content Quality"
description: "CI/CD workflows for content validation, link checking, and quality gates in GitHub"
date: 2026-01-24
published: true
slug: /docs/guides/quality-assurance/github-actions/
tags:
  - documentation
  - quality-assurance
  - github-actions
  - ci-cd
---


GitHub Actions provides powerful CI/CD capabilities for validating your markata-go site content. This guide covers workflows for build validation, link checking, and content quality gates.

## Table of Contents

- [Quick Start](#quick-start)
- [Build Validation Workflow](#build-validation-workflow)
- [Link Checking Workflow](#link-checking-workflow)
- [Content Quality Gates](#content-quality-gates)
- [Complete Quality Pipeline](#complete-quality-pipeline)
- [Reusable Workflows](#reusable-workflows)
- [Status Badges](#status-badges)

---

## Quick Start

Create `.github/workflows/content-quality.yml`:

```yaml
name: Content Quality

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Lint Markdown
        uses: DavidAnson/markdownlint-cli2-action@v16
        with:
          globs: |
            **/*.md
            !node_modules
            !public

      - name: Check YAML frontmatter
        run: |
          pip install yamllint
          find . -name "*.md" -exec grep -l "^---" {} \; | \
            xargs -I {} sh -c 'sed -n "/^---$/,/^---$/p" {} | yamllint -'
```

---

## Build Validation Workflow

Ensure your site builds successfully on every push and pull request.

### Basic Build Check

```yaml
# .github/workflows/build.yml
name: Build Validation

on:
  push:
    branches: [main]
    paths:
      - '**.md'
      - '**.toml'
      - 'templates/**'
      - 'static/**'
  pull_request:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'
          cache: true

      - name: Install markata-go
        run: go install github.com/waylonwalker/markata-go/cmd/markata-go@latest

      - name: Validate config
        run: markata-go config validate

      - name: Build site
        run: markata-go build --clean
        env:
          MARKATA_GO_URL: ${{ github.event_name == 'pull_request' && format('https://preview-{0}.example.com', github.event.number) || 'https://example.com' }}

      - name: Upload build artifact
        uses: actions/upload-artifact@v4
        with:
          name: site-build
          path: public/
          retention-days: 7
```

### Build with Cache

Speed up builds by caching Go modules and the markata-go binary:

```yaml
# .github/workflows/build-cached.yml
name: Build (Cached)

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'
          cache: true

      - name: Cache markata-go binary
        uses: actions/cache@v4
        id: markata-cache
        with:
          path: ~/go/bin/markata-go
          key: markata-go-${{ runner.os }}-${{ hashFiles('go.sum') }}
          restore-keys: |
            markata-go-${{ runner.os }}-

      - name: Install markata-go
        if: steps.markata-cache.outputs.cache-hit != 'true'
        run: go install github.com/waylonwalker/markata-go/cmd/markata-go@latest

      - name: Build
        run: ~/go/bin/markata-go build --clean
```

---

## Link Checking Workflow

Check for broken links across your entire site.

### Using lychee (Recommended)

[lychee](https://github.com/lycheeverse/lychee) is fast and handles many edge cases:

```yaml
# .github/workflows/links.yml
name: Link Checker

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
  schedule:
    # Run weekly to catch external link rot
    - cron: '0 0 * * 0'

jobs:
  check-links:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'
          cache: true

      - name: Install and build
        run: |
          go install github.com/waylonwalker/markata-go/cmd/markata-go@latest
          markata-go build --clean
        env:
          MARKATA_GO_URL: https://example.com

      - name: Check links
        uses: lycheeverse/lychee-action@v1
        with:
          args: >-
            --verbose
            --no-progress
            --accept 200,204,301,302,307,308
            --timeout 30
            --max-retries 3
            --exclude-path node_modules
            --exclude-path .git
            --exclude 'https://localhost.*'
            --exclude 'https://127\.0\.0\.1.*'
            --exclude 'mailto:.*'
            './public/**/*.html'
          fail: true
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

      - name: Create issue on failure
        if: failure() && github.event_name == 'schedule'
        uses: peter-evans/create-issue-from-file@v5
        with:
          title: Broken links detected
          content-filepath: ./lychee/out.md
          labels: bug, documentation
```

### Using markdown-link-check

For checking Markdown files directly (without building):

```yaml
# .github/workflows/markdown-links.yml
name: Markdown Link Check

on:
  push:
    branches: [main]
    paths: ['**.md']
  pull_request:
    branches: [main]
    paths: ['**.md']

jobs:
  check-links:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4

      - name: Check Markdown links
        uses: gaurav-nelson/github-action-markdown-link-check@v1
        with:
          config-file: '.markdown-link-check.json'
          use-quiet-mode: 'yes'
          use-verbose-mode: 'yes'
          folder-path: 'docs/'
          file-extension: '.md'
```

### Link Check Configuration

Create `.markdown-link-check.json`:

```json
{
  "ignorePatterns": [
    { "pattern": "^https://localhost" },
    { "pattern": "^https://127\\.0\\.0\\.1" },
    { "pattern": "^#" },
    { "pattern": "^mailto:" }
  ],
  "replacementPatterns": [
    {
      "pattern": "^/docs/",
      "replacement": "https://example.com/docs/"
    }
  ],
  "httpHeaders": [
    {
      "urls": ["https://github.com", "https://api.github.com"],
      "headers": {
        "Accept": "text/html, application/vnd.github.v3+json"
      }
    }
  ],
  "timeout": "20s",
  "retryOn429": true,
  "retryCount": 3,
  "fallbackRetryDelay": "10s",
  "aliveStatusCodes": [200, 206, 301, 302, 307, 308]
}
```

---

## Content Quality Gates

Enforce quality standards before merging pull requests.

### Comprehensive Linting

```yaml
# .github/workflows/lint.yml
name: Content Lint

on:
  pull_request:
    branches: [main]
    paths: ['**.md']

jobs:
  markdown-lint:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4

      - name: Lint Markdown files
        uses: DavidAnson/markdownlint-cli2-action@v16
        with:
          config: .markdownlint.json
          globs: |
            **/*.md
            !node_modules/**
            !public/**

  yaml-lint:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4

      - name: Setup Python
        uses: actions/setup-python@v5
        with:
          python-version: '3.12'

      - name: Install yamllint
        run: pip install yamllint

      - name: Lint YAML frontmatter
        run: |
          for file in $(find . -name "*.md" -not -path "./node_modules/*" -not -path "./public/*"); do
            echo "Checking: $file"
            # Extract frontmatter and lint it
            frontmatter=$(sed -n '1,/^---$/p' "$file" | tail -n +2 | head -n -1)
            if [ -n "$frontmatter" ]; then
              echo "$frontmatter" | yamllint -c .yamllint.yml -
            fi
          done

  frontmatter-check:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4

      - name: Check required frontmatter
        run: |
          errors=0
          for file in $(find docs -name "*.md"); do
            # Check for required fields
            for field in title date published; do
              if ! grep -q "^${field}:" "$file"; then
                echo "::error file=$file::Missing required field: $field"
                errors=$((errors + 1))
              fi
            done
          done
          exit $errors

  alt-text-check:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4

      - name: Check image alt text
        run: |
          errors=0
          for file in $(find . -name "*.md" -not -path "./node_modules/*"); do
            # Find images with empty alt text
            if grep -Pn '!\[\s*\]\(' "$file"; then
              echo "::error file=$file::Found image without alt text"
              errors=$((errors + 1))
            fi
          done
          exit $errors
```

### Spell Checking

```yaml
# .github/workflows/spelling.yml
name: Spell Check

on:
  pull_request:
    branches: [main]
    paths: ['**.md']

jobs:
  spellcheck:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4

      - name: Check spelling
        uses: crate-ci/typos@master
        with:
          files: ./docs
          config: .typos.toml
```

Create `.typos.toml`:

```toml
# .typos.toml
[default.extend-words]
# Add custom words that aren't typos
markata = "markata"
frontmatter = "frontmatter"
goldmark = "goldmark"

[files]
extend-exclude = [
    "*.json",
    "*.toml",
    "public/",
    "node_modules/",
]
```

---

## Complete Quality Pipeline

A comprehensive workflow that runs all quality checks:

```yaml
# .github/workflows/quality.yml
name: Content Quality Pipeline

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

concurrency:
  group: ${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true

jobs:
  # ============================================
  # Stage 1: Fast checks (run in parallel)
  # ============================================
  
  markdown-lint:
    name: Markdown Lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: DavidAnson/markdownlint-cli2-action@v16
        with:
          config: .markdownlint.json
          globs: '**/*.md'

  yaml-lint:
    name: YAML Lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: pip install yamllint
      - run: yamllint -c .yamllint.yml .

  frontmatter:
    name: Frontmatter Validation
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Validate frontmatter
        run: |
          errors=0
          for file in $(find docs -name "*.md" 2>/dev/null || true); do
            for field in title date published; do
              if ! grep -q "^${field}:" "$file" 2>/dev/null; then
                echo "::error file=$file::Missing: $field"
                errors=$((errors + 1))
              fi
            done
          done
          [ $errors -eq 0 ] || exit 1

  spelling:
    name: Spell Check
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: crate-ci/typos@master

  # ============================================
  # Stage 2: Build (depends on lint)
  # ============================================
  
  build:
    name: Build Site
    needs: [markdown-lint, yaml-lint, frontmatter]
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'
          cache: true

      - name: Install markata-go
        run: go install github.com/waylonwalker/markata-go/cmd/markata-go@latest

      - name: Build site
        run: markata-go build --clean
        env:
          MARKATA_GO_URL: https://example.com

      - name: Upload artifact
        uses: actions/upload-artifact@v4
        with:
          name: site
          path: public/

  # ============================================
  # Stage 3: Post-build checks
  # ============================================
  
  html-validate:
    name: HTML Validation
    needs: build
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Download build
        uses: actions/download-artifact@v4
        with:
          name: site
          path: public/

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'

      - name: Install html-validate
        run: npm install -g html-validate

      - name: Validate HTML
        run: |
          html-validate "public/**/*.html" --config .htmlvalidate.json || true
        continue-on-error: true

  links:
    name: Link Check
    needs: build
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Download build
        uses: actions/download-artifact@v4
        with:
          name: site
          path: public/

      - name: Check links
        uses: lycheeverse/lychee-action@v1
        with:
          args: >-
            --verbose
            --no-progress
            --accept 200,204,301,302,307,308
            --exclude-path node_modules
            --exclude 'localhost'
            './public/**/*.html'
          fail: false
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

  # ============================================
  # Stage 4: Deploy (only on main)
  # ============================================
  
  deploy:
    name: Deploy
    needs: [build, html-validate, links]
    if: github.ref == 'refs/heads/main' && github.event_name == 'push'
    runs-on: ubuntu-latest
    
    permissions:
      contents: read
      pages: write
      id-token: write
    
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    
    steps:
      - name: Download build
        uses: actions/download-artifact@v4
        with:
          name: site
          path: public/

      - name: Setup Pages
        uses: actions/configure-pages@v5

      - name: Upload to Pages
        uses: actions/upload-pages-artifact@v3
        with:
          path: public/

      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

### HTML Validation Config

Create `.htmlvalidate.json`:

```json
{
  "extends": ["html-validate:recommended"],
  "rules": {
    "no-trailing-whitespace": "off",
    "void-style": "off",
    "attribute-boolean-style": "off"
  }
}
```

---

## Reusable Workflows

Create reusable workflows for consistency across multiple repositories.

### Reusable Quality Check

```yaml
# .github/workflows/reusable-quality.yml
name: Reusable Quality Check

on:
  workflow_call:
    inputs:
      site-url:
        description: 'Site URL for build'
        required: true
        type: string
      node-version:
        description: 'Node.js version'
        required: false
        type: string
        default: '20'

jobs:
  quality:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'
          cache: true

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: ${{ inputs.node-version }}

      - name: Install tools
        run: |
          go install github.com/waylonwalker/markata-go/cmd/markata-go@latest
          npm install -g markdownlint-cli

      - name: Lint
        run: markdownlint '**/*.md' --ignore node_modules --ignore public

      - name: Build
        run: markata-go build --clean
        env:
          MARKATA_GO_URL: ${{ inputs.site-url }}

      - name: Upload artifact
        uses: actions/upload-artifact@v4
        with:
          name: site
          path: public/
```

### Using the Reusable Workflow

```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [main]
  pull_request:

jobs:
  quality:
    uses: ./.github/workflows/reusable-quality.yml
    with:
      site-url: https://example.com
```

---

## Status Badges

Add status badges to your README:

```markdown
# My Site

[![Content Quality](https://github.com/username/repo/actions/workflows/quality.yml/badge.svg)](https://github.com/username/repo/actions/workflows/quality.yml)
[![Links](https://github.com/username/repo/actions/workflows/links.yml/badge.svg)](https://github.com/username/repo/actions/workflows/links.yml)
[![Deploy](https://github.com/username/repo/actions/workflows/deploy.yml/badge.svg)](https://github.com/username/repo/actions/workflows/deploy.yml)
```

---

## See Also

- [Pre-commit Hooks](pre-commit-hooks/) - Local quality checks
- [GitLab CI](gitlab-ci/) - GitLab CI/CD configuration
- [Troubleshooting](troubleshooting/) - Common issues and solutions
