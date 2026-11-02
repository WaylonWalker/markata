---
title: "Quality Assurance Guide"
description: "Overview of content quality assurance tools and practices for markata-go sites"
date: 2026-01-24
published: true
slug: /docs/guides/quality-assurance/
tags:
  - documentation
  - quality-assurance
  - linting
  - ci-cd
---

## Quality Assurance for Your markata-go Site

Maintaining high-quality content is essential for any professional site. This guide covers tools and practices for validating frontmatter, linting Markdown, checking links, and ensuring accessibility in your markata-go site.

## Table of Contents

- [Why Quality Assurance?](#why-quality-assurance)
- [Quick Setup](#quick-setup)
- [Guide Sections](#guide-sections)
- [Recommended Tool Stack](#recommended-tool-stack)

---

## Why Quality Assurance?

Content quality issues can harm your site in several ways:

| Issue | Impact |
|-------|--------|
| Invalid frontmatter | Build failures, missing metadata |
| Broken links | Poor user experience, SEO penalties |
| Inconsistent Markdown | Rendering issues, maintenance burden |
| Missing alt text | Accessibility violations, SEO impact |
| YAML syntax errors | Build failures, data loss |

Automated quality checks catch these issues before they reach production.

---

## Quick Setup

Get started with content quality checks in under 5 minutes:

### 1. Install pre-commit

```bash
# macOS
brew install pre-commit

# pip (all platforms)
pip install pre-commit
```

### 2. Create Configuration

Create `.pre-commit-config.yaml` in your site root:

```yaml
repos:
  # YAML/Frontmatter validation
  - repo: https://github.com/adrienverge/yamllint
    rev: v1.35.1
    hooks:
      - id: yamllint
        args: [--config-file, .yamllint.yml]
        types: [markdown]

  # Markdown linting
  - repo: https://github.com/igorshubovych/markdownlint-cli
    rev: v0.42.0
    hooks:
      - id: markdownlint
        args: [--config, .markdownlint.json]
```

### 3. Add Linter Configs

Create `.yamllint.yml`:

```yaml
extends: relaxed
rules:
  line-length: disable
  document-start: disable
```

Create `.markdownlint.json`:

```json
{
  "MD013": false,
  "MD033": false,
  "MD041": false
}
```

### 4. Install and Run

```bash
pre-commit install
pre-commit run --all-files
```

Now every commit will automatically validate your content!

---

## Guide Sections

This quality assurance documentation is organized into the following sections:

### [Pre-commit Hooks](pre-commit-hooks/)

Set up automated checks that run before every commit:

- YAML frontmatter validation
- Markdown linting with markdownlint
- Link checking
- Image alt text validation
- Custom validation scripts

### [GitHub Actions](github-actions/)

Continuous integration workflows for GitHub:

- Build validation on pull requests
- Link checking across your entire site
- Content quality gates
- Automated reporting

### [GitLab CI](gitlab-ci/)

CI/CD pipelines for GitLab:

- Pipeline configuration
- Caching for faster builds
- Merge request quality gates
- Badge generation

### [Custom Rules](custom-rules/)

Configure linting rules for your specific needs:

- markdownlint rule customization
- Custom frontmatter validation
- Project-specific style guides
- Rule inheritance and overrides

### [Troubleshooting](troubleshooting/)

Common issues and solutions:

- False positives and how to handle them
- Performance optimization
- Integration debugging
- Migration from other tools

### [Editor Integration](editor-integration/)

Use lint output with your editor's quickfix features:

- Vim/Neovim quickfix integration
- VS Code problem matchers
- Emacs compile-mode
- Sublime Text build systems

---

## Recommended Tool Stack

Here's the recommended set of tools for comprehensive content quality assurance:

### Essential Tools

| Tool | Purpose | Installation |
|------|---------|--------------|
| [pre-commit](https://pre-commit.com/) | Git hooks framework | `pip install pre-commit` |
| [markdownlint-cli](https://github.com/igorshubovych/markdownlint-cli) | Markdown linting | `npm install -g markdownlint-cli` |
| [yamllint](https://github.com/adrienverge/yamllint) | YAML validation | `pip install yamllint` |

### Recommended Additions

| Tool | Purpose | Installation |
|------|---------|--------------|
| [lychee](https://github.com/lycheeverse/lychee) | Fast link checking | `cargo install lychee` |
| [htmltest](https://github.com/wjdp/htmltest) | HTML validation | `go install github.com/wjdp/htmltest@latest` |
| [pa11y](https://pa11y.org/) | Accessibility testing | `npm install -g pa11y` |

### CI-Specific Tools

| Tool | Purpose | Platform |
|------|---------|----------|
| [actionlint](https://github.com/rhysd/actionlint) | GitHub Actions validation | GitHub |
| [super-linter](https://github.com/super-linter/super-linter) | Multi-language linting | GitHub |

---

## Next Steps

1. **Start with pre-commit hooks** - They provide immediate feedback during development
2. **Add CI checks** - Ensure quality gates in your pull request workflow
3. **Customize rules** - Adjust linting rules to match your style guide
4. **Monitor and iterate** - Review false positives and refine your configuration

Choose your next section based on your development workflow:

- Using Git locally? Start with [Pre-commit Hooks](pre-commit-hooks/)
- Using GitHub? Jump to [GitHub Actions](github-actions/)
- Using GitLab? See [GitLab CI](gitlab-ci/)
