---
title: "Troubleshooting Quality Checks"
description: "Common issues and solutions for linting, pre-commit hooks, and CI/CD quality pipelines"
date: 2026-01-24
published: true
slug: qa-troubleshooting
tags:
  - quality-assurance
  - documentation
  - troubleshooting
---


This guide covers common issues you may encounter when setting up and running content quality checks for your markata-go site, along with their solutions.

## Table of Contents

- [Pre-commit Hook Issues](#pre-commit-hook-issues)
- [Markdown Linting Issues](#markdown-linting-issues)
- [YAML/Frontmatter Issues](#yamlfrontmatter-issues)
- [Link Checking Issues](#link-checking-issues)
- [CI/CD Pipeline Issues](#cicd-pipeline-issues)
- [Performance Issues](#performance-issues)
- [Migration from Other Tools](#migration-from-other-tools)

---

## Pre-commit Hook Issues

### Hooks Not Running

**Symptom:** Pre-commit hooks don't execute when you commit.

**Solutions:**

1. **Reinstall hooks:**
   ```bash
   pre-commit uninstall
   pre-commit install
   ```

2. **Verify git hooks directory:**
   ```bash
   ls -la .git/hooks/
   # Should see pre-commit -> /path/to/pre-commit-hook
   ```

3. **Check pre-commit installation:**
   ```bash
   pre-commit --version
   # If not found, reinstall: pip install pre-commit
   ```

4. **Verify configuration file exists:**
   ```bash
   ls -la .pre-commit-config.yaml
   ```

### Hook Fails with "Command Not Found"

**Symptom:** Hook fails because it can't find a command like `markdownlint` or `yamllint`.

**Solutions:**

1. **For system hooks**, install the tool globally:
   ```bash
   # markdownlint
   npm install -g markdownlint-cli

   # yamllint
   pip install yamllint
   ```

2. **Use repository-based hooks** instead of local hooks:
   ```yaml
   # Instead of local hook:
   - repo: https://github.com/igorshubovych/markdownlint-cli
     rev: v0.42.0
     hooks:
       - id: markdownlint
   ```

3. **Check PATH in your shell:**
   ```bash
   echo $PATH
   which markdownlint
   ```

### Hooks Running Too Slowly

**Symptom:** Commits take a long time because hooks are slow.

**Solutions:**

1. **Limit files processed:**
   ```yaml
   hooks:
     - id: markdownlint
       files: ^docs/  # Only check docs directory
   ```

2. **Skip slow hooks for WIP commits:**
   ```bash
   SKIP=check-links git commit -m "WIP: quick save"
   ```

3. **Move slow checks to CI only:**
   ```yaml
   # In .pre-commit-config.yaml
   hooks:
     - id: check-links
       stages: [push]  # Only run on push, not commit
   ```

4. **Clear pre-commit cache:**
   ```bash
   pre-commit clean
   pre-commit gc
   ```

### Skipping Hooks Temporarily

**When to skip:** Emergency fixes, work-in-progress commits, or when hooks have false positives.

```bash
# Skip all hooks (use sparingly!)
git commit --no-verify -m "Emergency fix"

# Skip specific hooks
SKIP=markdownlint,check-links git commit -m "Quick fix"

# Skip hooks for a single file
git add -f file-with-issues.md
```

**Warning:** Skipping hooks bypasses quality checks. Use CI/CD as a safety net.

---

## Markdown Linting Issues

### MD013: Line Length Errors

**Symptom:** `MD013/line-length: Line length [expected: 80, actual: 150]`

**Solutions:**

1. **Disable for prose** (recommended for content sites):
   ```json
   {
     "MD013": false
   }
   ```

2. **Increase line length limit:**
   ```json
   {
     "MD013": {
       "line_length": 120,
       "code_blocks": false,
       "tables": false
     }
   }
   ```

3. **Disable for specific file:**
   ```markdown
   <!-- markdownlint-disable-file MD013 -->

   # My Document

   This file allows long lines...
   ```

### MD033: Inline HTML Warnings

**Symptom:** `MD033/no-inline-html: Inline HTML [Element: details]`

**Solutions:**

1. **Allow specific HTML elements:**
   ```json
   {
     "MD033": {
       "allowed_elements": [
         "details", "summary", "kbd", "br",
         "sup", "sub", "img", "video"
       ]
     }
   }
   ```

2. **Disable for a section:**
   ```markdown
   <!-- markdownlint-disable MD033 -->
   <details>
   <summary>Click to expand</summary>

   Content here...

   </details>
   <!-- markdownlint-enable MD033 -->
   ```

3. **Disable entirely** (not recommended):
   ```json
   {
     "MD033": false
   }
   ```

### MD041: First Line Should Be Heading

**Symptom:** `MD041/first-line-heading: First line in file should be a top-level heading`

**Why it happens:** Files with YAML frontmatter trigger this because the first content line isn't a heading.

**Solution:**

```json
{
  "MD041": false
}
```

Or configure to recognize frontmatter:

```json
{
  "MD025": {
    "front_matter_title": "^\\s*title\\s*[:=]"
  },
  "MD041": false
}
```

### MD024: Duplicate Headings

**Symptom:** `MD024/no-duplicate-heading: Multiple headings with the same content`

**Solution:** Allow duplicate headings for siblings only:

```json
{
  "MD024": {
    "siblings_only": true
  }
}
```

This allows:
```markdown
## Installation  <!-- OK: different parent -->

### Step 1

## Configuration

### Step 1  <!-- OK: different parent section -->
```

### False Positives in Code Blocks

**Symptom:** Linting errors appear for content inside code blocks.

**Solutions:**

1. **Ensure code blocks are properly fenced:**
   ```markdown
   ```bash
   # This is code, not a heading
   echo "hello"
   ```
   ```

2. **Check for missing language identifier:**
   ```json
   {
     "MD040": true  // Enforces language specification
   }
   ```

3. **Verify backtick count matches:**
   ````markdown
   ```python
   code here
   ```  <!-- Must have same number of backticks -->
   ````

---

## YAML/Frontmatter Issues

### Invalid YAML Syntax

**Symptom:** `yaml: line 3: mapping values are not allowed here`

**Common causes and fixes:**

1. **Unquoted special characters:**
   ```yaml
   # Wrong
   title: My Post: A Journey

   # Correct
   title: "My Post: A Journey"
   ```

2. **Improper indentation:**
   ```yaml
   # Wrong (tabs)
   tags:
   	- documentation

   # Correct (2 spaces)
   tags:
     - documentation
   ```

3. **Missing quotes around dates:**
   ```yaml
   # Wrong (may be interpreted as number)
   version: 1.0

   # Correct
   version: "1.0"
   ```

### Frontmatter Not Detected

**Symptom:** Frontmatter is treated as content, not metadata.

**Causes:**

1. **No opening delimiter on first line:**
   ```markdown
   <!-- Wrong: whitespace before --- -->
    ---
   title: My Post
   ---

   <!-- Correct: --- must be on line 1 -->
   ---
   title: My Post
   ---
   ```

2. **Wrong delimiter:**
   ```markdown
   <!-- Wrong -->
   ----
   title: My Post
   ----

   <!-- Correct: exactly three dashes -->
   ---
   title: My Post
   ---
   ```

3. **BOM or hidden characters:**
   ```bash
   # Check for BOM
   file your-file.md
   # Should show: UTF-8 Unicode text
   # NOT: UTF-8 Unicode (with BOM) text

   # Remove BOM
   sed -i '1s/^\xEF\xBB\xBF//' your-file.md
   ```

### Required Field Missing

**Symptom:** Build fails with "missing required field: title"

**Solutions:**

1. **Add missing field:**
   ```yaml
   ---
   title: "My Post Title"
   date: 2026-01-24
   published: true
   ---
   ```

2. **Check field spelling:**
   ```yaml
   # Wrong
   titel: "My Post"

   # Correct
   title: "My Post"
   ```

3. **Check for invisible characters:**
   ```bash
   # View hex dump of frontmatter
   head -5 your-file.md | xxd
   ```

### Date Format Issues

**Symptom:** `Invalid date format` or dates not sorting correctly.

**Correct formats:**

```yaml
# ISO 8601 (recommended)
date: 2026-01-24

# With time
date: 2026-01-24T10:30:00

# With timezone
date: 2026-01-24T10:30:00-05:00
```

**Wrong formats:**

```yaml
# Wrong: needs quotes or different format
date: January 24, 2026

# Wrong: ambiguous
date: 01/24/2026
```

---

## Link Checking Issues

### Rate Limiting

**Symptom:** Link checks fail with `429 Too Many Requests`.

**Solutions:**

1. **Configure retry behavior:**
   ```json
   {
     "retryOn429": true,
     "retryCount": 5,
     "fallbackRetryDelay": "30s"
   }
   ```

2. **Add delays between requests:**
   ```toml
   # .lychee.toml
   max_concurrency = 5
   delay = 1
   ```

3. **Exclude problematic domains:**
   ```json
   {
     "ignorePatterns": [
       { "pattern": "^https://twitter\\.com" },
       { "pattern": "^https://x\\.com" }
     ]
   }
   ```

### False Positives on Valid Links

**Symptom:** Link checker reports broken links that work in browser.

**Common causes:**

1. **Bot detection:**
   ```toml
   # .lychee.toml
   user_agent = "Mozilla/5.0 (compatible; link-checker)"
   ```

2. **Authentication required:**
   ```json
   {
     "ignorePatterns": [
       { "pattern": "^https://private-site\\.com" }
     ]
   }
   ```

3. **JavaScript-rendered content:**
   - Link checkers can't execute JavaScript
   - Add to ignore list if the link is valid

4. **Redirect chains:**
   ```json
   {
     "aliveStatusCodes": [200, 206, 301, 302, 307, 308]
   }
   ```

### Internal Links Not Resolving

**Symptom:** Relative links like `/docs/guide/` reported as broken.

**Solutions:**

1. **Configure base URL:**
   ```json
   {
     "replacementPatterns": [
       {
         "pattern": "^/",
         "replacement": "https://example.com/"
       }
     ]
   }
   ```

2. **Run link checker on built output:**
   ```bash
   # Build first
   markata-go build

   # Check built HTML
   lychee ./public/**/*.html
   ```

3. **Use relative path resolution:**
   ```toml
   # .lychee.toml
   base = "https://example.com"
   ```

### Mailto and Tel Links

**Symptom:** `mailto:` and `tel:` links reported as broken.

**Solution:** Exclude these patterns:

```json
{
  "ignorePatterns": [
    { "pattern": "^mailto:" },
    { "pattern": "^tel:" }
  ]
}
```

Or in lychee:

```toml
# .lychee.toml
exclude = [
  "^mailto:",
  "^tel:",
]
```

---

## CI/CD Pipeline Issues

### Pipeline Not Triggering

**Symptom:** Push/MR doesn't start the pipeline.

**GitHub Actions solutions:**

1. **Check workflow file location:**
   ```
   .github/workflows/quality.yml  # Correct
   .github/workflow/quality.yml   # Wrong (missing 's')
   ```

2. **Verify trigger configuration:**
   ```yaml
   on:
     push:
       branches: [main]
     pull_request:
       branches: [main]
   ```

3. **Check for YAML syntax errors:**
   ```bash
   # Validate workflow file
   yamllint .github/workflows/quality.yml
   ```

**GitLab CI solutions:**

1. **Check file name:**
   ```
   .gitlab-ci.yml  # Correct
   gitlab-ci.yml   # Wrong (missing leading dot)
   ```

2. **Validate configuration:**
   ```bash
   # In GitLab UI: CI/CD > Pipelines > CI Lint
   # Or use API
   ```

### Job Fails with Exit Code 1

**Symptom:** Job fails but output doesn't show clear error.

**Solutions:**

1. **Add verbose output:**
   ```yaml
   script:
     - markdownlint '**/*.md' --verbose
   ```

2. **Check previous command exit codes:**
   ```yaml
   script:
     - set -e  # Exit on first error
     - markdownlint '**/*.md'
     - echo "Lint passed"
   ```

3. **Use `|| true` for non-blocking checks:**
   ```yaml
   script:
     - markdownlint '**/*.md' || true  # Continue even if fails
   ```

### Artifacts Not Available

**Symptom:** Built files not available in later jobs.

**GitHub Actions:**

```yaml
jobs:
  build:
    steps:
      - uses: actions/upload-artifact@v4
        with:
          name: site
          path: public/

  test:
    needs: build
    steps:
      - uses: actions/download-artifact@v4
        with:
          name: site
          path: public/
```

**GitLab CI:**

```yaml
build:
  artifacts:
    paths:
      - public/
    expire_in: 1 hour

test:
  needs:
    - job: build
      artifacts: true
```

### Cache Not Working

**Symptom:** Dependencies re-download every run.

**GitHub Actions:**

```yaml
- uses: actions/setup-node@v4
  with:
    node-version: '20'
    cache: 'npm'  # Enable built-in caching
```

**GitLab CI:**

```yaml
variables:
  npm_config_cache: $CI_PROJECT_DIR/.npm

cache:
  key: $CI_COMMIT_REF_SLUG
  paths:
    - .npm/
  policy: pull-push  # Important: must be pull-push
```

---

## Performance Issues

### Pre-commit Runs Slowly

**Diagnosis:**

```bash
# Time individual hooks
time pre-commit run markdownlint --all-files
time pre-commit run yamllint --all-files
```

**Optimizations:**

1. **Limit scope:**
   ```yaml
   hooks:
     - id: markdownlint
       files: ^docs/  # Only docs directory
   ```

2. **Use faster alternatives:**
   - `markdownlint-cli2` is faster than `markdownlint-cli`
   - `lychee` is faster than `markdown-link-check`

3. **Skip expensive checks on commit:**
   ```yaml
   hooks:
     - id: check-links
       stages: [push]  # Only on push
   ```

### CI/CD Takes Too Long

**Optimizations:**

1. **Run jobs in parallel:**
   ```yaml
   # GitHub Actions
   jobs:
     lint:
       runs-on: ubuntu-latest
     build:
       runs-on: ubuntu-latest
       needs: []  # No dependency = parallel
   ```

2. **Use caching effectively:**
   ```yaml
   - uses: actions/cache@v4
     with:
       path: ~/.npm
       key: npm-${{ hashFiles('package-lock.json') }}
   ```

3. **Only run on relevant changes:**
   ```yaml
   on:
     push:
       paths:
         - '**/*.md'
         - '.markdownlint.json'
   ```

4. **Use lightweight images:**
   ```yaml
   # Instead of: image: node:20
   image: node:20-alpine  # Much smaller
   ```

---

## Migration from Other Tools

### From Jekyll

**Frontmatter differences:**

```yaml
# Jekyll
---
layout: post
title: My Post
categories: blog
---

# markata-go
---
title: "My Post"
date: 2026-01-24
published: true
tags:
  - blog
---
```

**Conversion script:**

```bash
#!/bin/bash
# Convert Jekyll frontmatter to markata-go format
for file in _posts/*.md; do
  # Add published: true if missing
  if ! grep -q "^published:" "$file"; then
    sed -i '/^---$/,/^---$/{/^title:/a published: true
    }' "$file"
  fi
done
```

### From Hugo

**Frontmatter differences:**

```yaml
# Hugo
---
title: "My Post"
date: 2026-01-24T10:00:00-05:00
draft: true
---

# markata-go
---
title: "My Post"
date: 2026-01-24
published: false
---
```

**Key changes:**
- `draft: true` becomes `published: false`
- Date format simplified (time optional)
- Remove Hugo-specific fields like `weight`, `type`

### From Gatsby/MDX

**MDX components won't work directly.** Convert to standard Markdown:

```jsx
// Gatsby MDX
<Callout type="info">
  Important information here
</Callout>

// markata-go Markdown
> **Info:** Important information here
```

Or use HTML:

```html
<details>
<summary>Important information</summary>

Content here...

</details>
```

---

## Getting Help

If you can't resolve an issue:

1. **Check existing issues:** Search [GitHub Issues](https://github.com/waylonwalker/markata-go/issues)

2. **Enable verbose output:** Add `--verbose` or `-v` flags to commands

3. **Isolate the problem:** Create a minimal reproduction case

4. **Gather information:**
   ```bash
   # System info
   uname -a
   go version
   node --version
   python --version

   # Tool versions
   pre-commit --version
   markdownlint --version
   ```

5. **Open an issue** with:
   - What you expected
   - What actually happened
   - Steps to reproduce
   - Relevant configuration files
   - Error messages (full output)

---

## See Also

- [Pre-commit Hooks](pre-commit-hooks/) - Hook setup and configuration
- [GitHub Actions](github-actions/) - CI/CD for GitHub
- [GitLab CI](gitlab-ci/) - CI/CD for GitLab
- [Custom Rules](custom-rules/) - Creating custom linting rules
