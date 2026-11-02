---
title: "GitLab CI for Content Quality"
description: "CI/CD pipelines for content validation, link checking, and quality gates in GitLab"
date: 2026-01-24
published: true
slug: /docs/guides/quality-assurance/gitlab-ci/
tags:
  - documentation
  - quality-assurance
  - gitlab-ci
  - ci-cd
---


GitLab CI/CD provides powerful pipeline capabilities for validating your markata-go site content. This guide covers pipeline configuration, caching strategies, merge request quality gates, and badge generation.

## Table of Contents

- [Quick Start](#quick-start)
- [Basic Pipeline Configuration](#basic-pipeline-configuration)
- [Build Validation Pipeline](#build-validation-pipeline)
- [Link Checking Pipeline](#link-checking-pipeline)
- [Content Quality Gates](#content-quality-gates)
- [Complete Quality Pipeline](#complete-quality-pipeline)
- [Caching Strategies](#caching-strategies)
- [Merge Request Integration](#merge-request-integration)
- [Badge Generation](#badge-generation)
- [GitLab Pages Deployment](#gitlab-pages-deployment)

---

## Quick Start

Create `.gitlab-ci.yml` in your repository root:

```yaml
stages:
  - lint
  - build
  - test

lint:
  stage: lint
  image: python:3.12-slim
  before_script:
    - pip install yamllint
    - apt-get update && apt-get install -y nodejs npm
    - npm install -g markdownlint-cli
  script:
    - markdownlint '**/*.md' --ignore node_modules --ignore public
    - |
      for file in $(find . -name "*.md" -not -path "./node_modules/*"); do
        frontmatter=$(sed -n '1,/^---$/p' "$file" | tail -n +2 | head -n -1)
        if [ -n "$frontmatter" ]; then
          echo "$frontmatter" | yamllint -d relaxed -
        fi
      done
  rules:
    - changes:
        - "**/*.md"
```

---

## Basic Pipeline Configuration

### Pipeline Structure

```yaml
# .gitlab-ci.yml
stages:
  - lint        # Fast checks
  - build       # Build site
  - test        # Post-build validation
  - deploy      # Deploy to production

variables:
  SITE_URL: https://example.com
  GO_VERSION: "1.22"

default:
  image: golang:${GO_VERSION}-alpine
  before_script:
    - go install github.com/waylonwalker/markata-go/cmd/markata-go@latest
```

### Environment Variables

Configure these variables in GitLab Settings > CI/CD > Variables:

| Variable | Description | Example |
|----------|-------------|---------|
| `SITE_URL` | Production site URL | `https://example.com` |
| `DEPLOY_TOKEN` | Deployment credentials | `glpat-xxx` |
| `SLACK_WEBHOOK` | Notifications URL | `https://hooks.slack.com/...` |

---

## Build Validation Pipeline

### Basic Build Job

```yaml
# .gitlab-ci.yml
stages:
  - lint
  - build

variables:
  GO_VERSION: "1.22"

build:
  stage: build
  image: golang:${GO_VERSION}-alpine
  
  before_script:
    - go install github.com/waylonwalker/markata-go/cmd/markata-go@latest
  
  script:
    - markata-go config validate
    - markata-go build --clean
  
  artifacts:
    paths:
      - public/
    expire_in: 1 week
  
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
    - if: $CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH
```

### Build with Environment-Specific URLs

```yaml
build:
  stage: build
  image: golang:1.22-alpine
  
  before_script:
    - go install github.com/waylonwalker/markata-go/cmd/markata-go@latest
  
  script:
    - markata-go build --clean
  
  variables:
    MARKATA_GO_URL: >-
      ${CI_MERGE_REQUEST_IID:+https://preview-${CI_MERGE_REQUEST_IID}.example.com}
      ${CI_MERGE_REQUEST_IID:-https://example.com}
  
  artifacts:
    paths:
      - public/
    expire_in: 1 week
```

---

## Link Checking Pipeline

### Using lychee

```yaml
check-links:
  stage: test
  image: lycheeverse/lychee:latest
  needs: [build]
  
  script:
    - |
      lychee \
        --verbose \
        --no-progress \
        --accept 200,204,301,302,307,308 \
        --timeout 30 \
        --max-retries 3 \
        --exclude-path node_modules \
        --exclude-path .git \
        --exclude 'localhost' \
        --exclude '127\.0\.0\.1' \
        --exclude 'mailto:' \
        './public/**/*.html'
  
  allow_failure: true
  
  artifacts:
    reports:
      junit: lychee-report.xml
    when: always
  
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
    - if: $CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH
    - if: $CI_PIPELINE_SOURCE == "schedule"
```

### Scheduled Link Checks

Add to your pipeline for weekly link rot detection:

```yaml
check-external-links:
  stage: test
  image: lycheeverse/lychee:latest
  
  script:
    - |
      lychee \
        --verbose \
        --include-mail \
        --timeout 60 \
        --max-retries 5 \
        './public/**/*.html' || true
    
    # Create issue if links are broken
    - |
      if [ -f lychee/out.md ] && [ -s lychee/out.md ]; then
        echo "Broken links detected. See artifacts for details."
        exit 1
      fi
  
  artifacts:
    paths:
      - lychee/
    expire_in: 1 month
  
  rules:
    - if: $CI_PIPELINE_SOURCE == "schedule"
```

Configure a scheduled pipeline in GitLab Settings > CI/CD > Schedules:
- Description: "Weekly link check"
- Interval: `0 0 * * 0` (weekly on Sunday)
- Target branch: `main`

---

## Content Quality Gates

### Comprehensive Linting Job

```yaml
lint:
  stage: lint
  image: python:3.12-slim
  
  before_script:
    - pip install yamllint
    - apt-get update && apt-get install -y nodejs npm
    - npm install -g markdownlint-cli
  
  script:
    # Markdown linting
    - markdownlint '**/*.md' --ignore node_modules --ignore public --config .markdownlint.json
    
    # YAML frontmatter validation
    - |
      for file in $(find . -name "*.md" -not -path "./node_modules/*" -not -path "./public/*"); do
        echo "Checking: $file"
        frontmatter=$(sed -n '1,/^---$/p' "$file" | tail -n +2 | head -n -1)
        if [ -n "$frontmatter" ]; then
          echo "$frontmatter" | yamllint -c .yamllint.yml - || exit 1
        fi
      done
  
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
      changes:
        - "**/*.md"

frontmatter-validation:
  stage: lint
  image: alpine:latest
  
  script:
    - |
      errors=0
      for file in $(find docs -name "*.md" 2>/dev/null); do
        for field in title date published; do
          if ! grep -q "^${field}:" "$file"; then
            echo "ERROR: $file - Missing required field: $field"
            errors=$((errors + 1))
          fi
        done
      done
      
      if [ $errors -gt 0 ]; then
        echo "Found $errors frontmatter errors"
        exit 1
      fi
      echo "All frontmatter valid!"
  
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
      changes:
        - "**/*.md"

alt-text-check:
  stage: lint
  image: alpine:latest
  
  script:
    - |
      errors=0
      for file in $(find . -name "*.md" -not -path "./node_modules/*"); do
        if grep -Pn '!\[\s*\]\(' "$file"; then
          echo "ERROR: $file - Found image(s) without alt text"
          errors=$((errors + 1))
        fi
      done
      
      if [ $errors -gt 0 ]; then
        echo "Found $errors images missing alt text"
        exit 1
      fi
      echo "All images have alt text!"
  
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
      changes:
        - "**/*.md"
```

### Spell Checking

```yaml
spellcheck:
  stage: lint
  image: rust:latest
  
  before_script:
    - cargo install typos-cli
  
  script:
    - typos ./docs --config .typos.toml
  
  allow_failure: true
  
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
      changes:
        - "**/*.md"
```

---

## Complete Quality Pipeline

A comprehensive pipeline that runs all quality checks:

```yaml
# .gitlab-ci.yml
# Complete content quality pipeline for markata-go sites

stages:
  - lint
  - build
  - test
  - deploy

variables:
  GO_VERSION: "1.22"
  NODE_VERSION: "20"
  SITE_URL: https://example.com

# ============================================
# Templates for reuse
# ============================================

.go-setup:
  image: golang:${GO_VERSION}-alpine
  before_script:
    - go install github.com/waylonwalker/markata-go/cmd/markata-go@latest

.node-setup:
  image: node:${NODE_VERSION}-alpine
  before_script:
    - npm install -g markdownlint-cli

.python-setup:
  image: python:3.12-slim
  before_script:
    - pip install yamllint

# ============================================
# Stage 1: Lint (parallel jobs)
# ============================================

markdown-lint:
  extends: .node-setup
  stage: lint
  
  script:
    - markdownlint '**/*.md' --ignore node_modules --ignore public --config .markdownlint.json
  
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
      changes:
        - "**/*.md"
    - if: $CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH

yaml-lint:
  extends: .python-setup
  stage: lint
  
  script:
    - yamllint -c .yamllint.yml .
  
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
      changes:
        - "**/*.md"
        - "**/*.yml"
        - "**/*.yaml"

frontmatter:
  stage: lint
  image: alpine:latest
  
  script:
    - |
      errors=0
      for file in $(find docs -name "*.md" 2>/dev/null || true); do
        for field in title date published; do
          if ! grep -q "^${field}:" "$file" 2>/dev/null; then
            echo "ERROR: $file - Missing: $field"
            errors=$((errors + 1))
          fi
        done
      done
      [ $errors -eq 0 ] || exit 1
  
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
      changes:
        - "**/*.md"

# ============================================
# Stage 2: Build
# ============================================

build:
  extends: .go-setup
  stage: build
  needs: [markdown-lint, yaml-lint, frontmatter]
  
  script:
    - markata-go config validate
    - markata-go build --clean
  
  variables:
    MARKATA_GO_URL: $SITE_URL
  
  artifacts:
    paths:
      - public/
    expire_in: 1 week
  
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
    - if: $CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH

# ============================================
# Stage 3: Test (post-build validation)
# ============================================

html-validate:
  stage: test
  image: node:${NODE_VERSION}-alpine
  needs: [build]
  
  before_script:
    - npm install -g html-validate
  
  script:
    - html-validate "public/**/*.html" --config .htmlvalidate.json || true
  
  allow_failure: true
  
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
    - if: $CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH

check-links:
  stage: test
  image: lycheeverse/lychee:latest
  needs: [build]
  
  script:
    - |
      lychee \
        --verbose \
        --no-progress \
        --accept 200,204,301,302,307,308 \
        --exclude-path node_modules \
        --exclude 'localhost' \
        './public/**/*.html'
  
  allow_failure: true
  
  artifacts:
    paths:
      - lychee/
    when: always
  
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
    - if: $CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH

# ============================================
# Stage 4: Deploy
# ============================================

pages:
  stage: deploy
  needs: [build, html-validate, check-links]
  
  script:
    - echo "Deploying to GitLab Pages"
  
  artifacts:
    paths:
      - public/
  
  environment:
    name: production
    url: $SITE_URL
  
  rules:
    - if: $CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH
```

---

## Caching Strategies

Speed up your pipelines with effective caching.

### Go Module Cache

```yaml
variables:
  GOPATH: $CI_PROJECT_DIR/.go

build:
  stage: build
  image: golang:1.22-alpine
  
  cache:
    key: go-modules-$CI_COMMIT_REF_SLUG
    paths:
      - .go/pkg/mod/
    policy: pull-push
  
  script:
    - go install github.com/waylonwalker/markata-go/cmd/markata-go@latest
    - markata-go build --clean
```

### Node.js Cache

```yaml
lint:
  stage: lint
  image: node:20-alpine
  
  cache:
    key: node-modules-$CI_COMMIT_REF_SLUG
    paths:
      - node_modules/
    policy: pull-push
  
  script:
    - npm install markdownlint-cli
    - npx markdownlint '**/*.md'
```

### Combined Cache Strategy

```yaml
variables:
  GOPATH: $CI_PROJECT_DIR/.go
  npm_config_cache: $CI_PROJECT_DIR/.npm

default:
  cache:
    - key: go-$CI_COMMIT_REF_SLUG
      paths:
        - .go/pkg/mod/
      policy: pull-push
    - key: npm-$CI_COMMIT_REF_SLUG
      paths:
        - .npm/
      policy: pull-push
```

---

## Merge Request Integration

### Merge Request Pipelines

Configure pipelines to run on merge requests:

```yaml
workflow:
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
    - if: $CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH
    - if: $CI_PIPELINE_SOURCE == "schedule"

lint:
  stage: lint
  script:
    - echo "Running lint checks"
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
      changes:
        - "**/*.md"
```

### Merge Request Comments

Report quality issues directly in merge request comments:

```yaml
report-quality:
  stage: test
  image: alpine:latest
  needs: [lint, build]
  
  script:
    - |
      # Collect lint results
      if [ -f lint-results.json ]; then
        echo "## Content Quality Report" > report.md
        echo "" >> report.md
        echo "### Linting Results" >> report.md
        cat lint-results.json | jq -r '.[] | "- \(.file): \(.message)"' >> report.md
      fi
    
    # Post comment to MR
    - |
      if [ -f report.md ]; then
        curl --request POST \
          --header "PRIVATE-TOKEN: ${GITLAB_TOKEN}" \
          --form "body=$(cat report.md)" \
          "${CI_API_V4_URL}/projects/${CI_PROJECT_ID}/merge_requests/${CI_MERGE_REQUEST_IID}/notes"
      fi
  
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
```

### Required Approvals

Configure merge request approval rules in GitLab Settings > Merge Requests:

1. **Require pipeline to succeed** before merging
2. Set **required approvals** for content changes
3. Configure **code owners** for documentation directories

Example `CODEOWNERS` file:

```
# .gitlab/CODEOWNERS
docs/       @docs-team
*.md        @docs-team @content-reviewers
```

---

## Badge Generation

Add status badges to your README.

### Pipeline Status Badge

```markdown
[![Pipeline Status](https://gitlab.com/username/repo/badges/main/pipeline.svg)](https://gitlab.com/username/repo/-/pipelines)
```

### Coverage Badge (if applicable)

```markdown
[![Coverage](https://gitlab.com/username/repo/badges/main/coverage.svg)](https://gitlab.com/username/repo/-/jobs)
```

### Custom Badges

Generate custom badges for quality metrics:

```yaml
generate-badges:
  stage: deploy
  image: alpine:latest
  
  script:
    # Count markdown files
    - file_count=$(find docs -name "*.md" | wc -l)
    
    # Generate badge JSON for shields.io
    - |
      cat > public/badges/docs-count.json << EOF
      {
        "schemaVersion": 1,
        "label": "docs",
        "message": "${file_count} pages",
        "color": "blue"
      }
      EOF
  
  artifacts:
    paths:
      - public/badges/
  
  rules:
    - if: $CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH
```

Use with shields.io:

```markdown
![Docs Count](https://img.shields.io/endpoint?url=https://example.com/badges/docs-count.json)
```

---

## GitLab Pages Deployment

### Basic Pages Deployment

```yaml
pages:
  stage: deploy
  
  script:
    - echo "Deploying to GitLab Pages"
  
  artifacts:
    paths:
      - public/
  
  rules:
    - if: $CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH
```

### Pages with Review Apps

Deploy preview environments for merge requests:

```yaml
pages:review:
  stage: deploy
  needs: [build]
  
  script:
    - echo "Deploying review app"
  
  artifacts:
    paths:
      - public/
  
  environment:
    name: review/$CI_COMMIT_REF_SLUG
    url: https://$CI_PROJECT_NAMESPACE.gitlab.io/-/$CI_PROJECT_NAME/-/jobs/$CI_JOB_ID/artifacts/public/index.html
    on_stop: stop:review
    auto_stop_in: 1 week
  
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"

stop:review:
  stage: deploy
  
  script:
    - echo "Stopping review app"
  
  environment:
    name: review/$CI_COMMIT_REF_SLUG
    action: stop
  
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
      when: manual
```

---

## Troubleshooting

### Common Issues

**Pipeline not triggering on MR**

Ensure workflow rules are configured:

```yaml
workflow:
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
    - if: $CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH
```

**Cache not working**

Check cache key and paths:

```yaml
cache:
  key: $CI_COMMIT_REF_SLUG  # Use branch-specific keys
  paths:
    - .go/pkg/mod/
  policy: pull-push  # Ensure both pull and push
```

**Job stuck in pending**

Check runner availability and tags:

```yaml
build:
  tags:
    - docker  # Match available runners
```

**Artifacts not available in downstream jobs**

Use `needs` to explicitly depend on artifacts:

```yaml
deploy:
  needs:
    - job: build
      artifacts: true
```

---

## See Also

- [Pre-commit Hooks](pre-commit-hooks/) - Local quality checks
- [GitHub Actions](github-actions/) - GitHub CI/CD configuration
- [Custom Rules](custom-rules/) - Configure linting rules
- [Troubleshooting](troubleshooting/) - Common issues and solutions
