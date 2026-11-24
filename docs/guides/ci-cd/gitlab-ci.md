---
title: "GitLab CI/CD"
description: "Complete guide to building and deploying markata-go sites with GitLab CI"
date: 2026-01-24
published: true
slug: /docs/guides/ci-cd/gitlab-ci/
tags:
  - documentation
  - ci-cd
  - deployment
  - gitlab
---


GitLab CI/CD provides integrated continuous integration and deployment directly within GitLab. This guide covers deploying markata-go sites to GitLab Pages and other platforms.

## Quick Start

Create `.gitlab-ci.yml` in your repository root:

```yaml
image: alpine:latest

stages:
  - build
  - deploy

variables:
  MARKATA_VERSION: "0.1.0"

build:
  stage: build
  before_script:
    - apk add --no-cache wget tar
    - wget -qO- "https://github.com/WaylonWalker/markata-go/releases/download/v${MARKATA_VERSION}/markata-go_${MARKATA_VERSION}_linux_x86_64.tar.gz" | tar xz
    - mv markata-go /usr/local/bin/
  script:
    - markata-go build --clean
  artifacts:
    paths:
      - public/
    expire_in: 1 hour

pages:
  stage: deploy
  script:
    - echo "Deploying to GitLab Pages..."
  artifacts:
    paths:
      - public/
  only:
    - main
  environment:
    name: production
    url: https://$CI_PROJECT_NAMESPACE.gitlab.io/$CI_PROJECT_NAME
```

Push to your repository, and GitLab will automatically build and deploy your site.

## GitLab Pages

### Basic Deployment

GitLab Pages requires a job named `pages` that produces artifacts in the `public/` directory:

```yaml
pages:
  stage: deploy
  script:
    - echo "Deploying to GitLab Pages"
  artifacts:
    paths:
      - public/
  only:
    - main
```

The key requirements:
- Job must be named `pages`
- Artifacts must include `public/` directory
- Only runs on your default branch

### With Site URL Configuration

```yaml
variables:
  MARKATA_VERSION: "0.1.0"
  SITE_URL: "https://mygroup.gitlab.io/myproject"

build:
  stage: build
  image: alpine:latest
  before_script:
    - apk add --no-cache wget tar
    - wget -qO- "https://github.com/WaylonWalker/markata-go/releases/download/v${MARKATA_VERSION}/markata-go_${MARKATA_VERSION}_linux_x86_64.tar.gz" | tar xz
    - mv markata-go /usr/local/bin/
  script:
    - MARKATA_GO_URL=$SITE_URL markata-go build --clean
  artifacts:
    paths:
      - public/
    expire_in: 1 hour

pages:
  stage: deploy
  dependencies:
    - build
  script:
    - echo "Deploying to GitLab Pages"
  artifacts:
    paths:
      - public/
  only:
    - main
```

### Custom Domain

1. Go to **Settings** > **Pages** > **New Domain**
2. Add your domain and verify ownership
3. Update your configuration:

```yaml
variables:
  SITE_URL: "https://example.com"
```

4. Configure DNS:
   - Add a `CNAME` record pointing to `<namespace>.gitlab.io`
   - Or use `A` records for apex domains

### Access Control

Control who can view your GitLab Pages site:

1. **Settings** > **General** > **Visibility**
2. Choose visibility level for Pages

For private Pages with authentication:

```yaml
pages:
  stage: deploy
  script:
    - echo "Deploying private site"
  artifacts:
    paths:
      - public/
  only:
    - main
```

Then enable **Pages Access Control** in project settings.

## Complete Pipeline

A full-featured pipeline with testing and multiple environments:

```yaml
image: alpine:latest

stages:
  - validate
  - build
  - test
  - deploy

variables:
  MARKATA_VERSION: "0.1.0"
  # Cache configuration
  GOPATH: ${CI_PROJECT_DIR}/.go

# Global cache for all jobs
cache:
  key: ${CI_COMMIT_REF_SLUG}
  paths:
    - .go/

# Reusable template for markata-go setup
.markata-setup: &markata-setup
  before_script:
    - apk add --no-cache wget tar
    - |
      if [ ! -f /usr/local/bin/markata-go ]; then
        wget -qO- "https://github.com/WaylonWalker/markata-go/releases/download/v${MARKATA_VERSION}/markata-go_${MARKATA_VERSION}_linux_x86_64.tar.gz" | tar xz
        mv markata-go /usr/local/bin/
      fi

# Validate configuration
validate:
  stage: validate
  <<: *markata-setup
  script:
    - markata-go config validate
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
    - if: $CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH

# Build the site
build:
  stage: build
  <<: *markata-setup
  script:
    - |
      if [ "$CI_COMMIT_BRANCH" == "$CI_DEFAULT_BRANCH" ]; then
        export MARKATA_GO_URL="https://${CI_PROJECT_NAMESPACE}.gitlab.io/${CI_PROJECT_NAME}"
      else
        export MARKATA_GO_URL=""
      fi
    - markata-go build --clean -v
  artifacts:
    paths:
      - public/
    expire_in: 1 day
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
    - if: $CI_COMMIT_BRANCH

# Test HTML output
test:html:
  stage: test
  image: node:20-alpine
  dependencies:
    - build
  script:
    - npm install -g html-validate
    - html-validate "public/**/*.html" || true
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
    - if: $CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH
  allow_failure: true

# Test for broken links
test:links:
  stage: test
  image: node:20-alpine
  dependencies:
    - build
  script:
    - npm install -g linkinator
    - linkinator public --recurse --skip "^(?!https?://)" || true
  rules:
    - if: $CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH
  allow_failure: true

# Deploy to GitLab Pages
pages:
  stage: deploy
  dependencies:
    - build
  script:
    - echo "Deploying to GitLab Pages"
  artifacts:
    paths:
      - public/
  environment:
    name: production
    url: https://${CI_PROJECT_NAMESPACE}.gitlab.io/${CI_PROJECT_NAME}
  rules:
    - if: $CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH
```

## Multi-Environment Deployments

Deploy to staging and production environments:

```yaml
stages:
  - build
  - deploy

variables:
  MARKATA_VERSION: "0.1.0"

.build-template: &build-template
  stage: build
  image: alpine:latest
  before_script:
    - apk add --no-cache wget tar
    - wget -qO- "https://github.com/WaylonWalker/markata-go/releases/download/v${MARKATA_VERSION}/markata-go_${MARKATA_VERSION}_linux_x86_64.tar.gz" | tar xz
    - mv markata-go /usr/local/bin/
  script:
    - markata-go build --clean

# Build for staging
build:staging:
  <<: *build-template
  variables:
    MARKATA_GO_URL: "https://staging.example.com"
  artifacts:
    paths:
      - public/
    expire_in: 1 day
  rules:
    - if: $CI_COMMIT_BRANCH == "develop"

# Build for production
build:production:
  <<: *build-template
  variables:
    MARKATA_GO_URL: "https://example.com"
  artifacts:
    paths:
      - public/
    expire_in: 1 day
  rules:
    - if: $CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH

# Deploy to staging environment
deploy:staging:
  stage: deploy
  dependencies:
    - build:staging
  script:
    - echo "Deploying to staging..."
    # Add your staging deployment commands here
  environment:
    name: staging
    url: https://staging.example.com
  rules:
    - if: $CI_COMMIT_BRANCH == "develop"

# Deploy to production (manual approval)
deploy:production:
  stage: deploy
  dependencies:
    - build:production
  script:
    - echo "Deploying to production..."
    # For GitLab Pages, this happens automatically
  environment:
    name: production
    url: https://example.com
  rules:
    - if: $CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH
  when: manual  # Require manual approval
```

## Review Apps (Preview Deployments)

Create dynamic preview environments for merge requests:

```yaml
stages:
  - build
  - review
  - deploy
  - cleanup

variables:
  MARKATA_VERSION: "0.1.0"

# Build for merge request review
build:review:
  stage: build
  image: alpine:latest
  before_script:
    - apk add --no-cache wget tar
    - wget -qO- "https://github.com/WaylonWalker/markata-go/releases/download/v${MARKATA_VERSION}/markata-go_${MARKATA_VERSION}_linux_x86_64.tar.gz" | tar xz
    - mv markata-go /usr/local/bin/
  script:
    - MARKATA_GO_URL="https://${CI_PROJECT_NAMESPACE}.gitlab.io/-/${CI_PROJECT_NAME}/-/jobs/${CI_JOB_ID}/artifacts/public" markata-go build --clean
  artifacts:
    paths:
      - public/
    expire_in: 1 week
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"

# Deploy review app
review:
  stage: review
  dependencies:
    - build:review
  script:
    - echo "Review app deployed"
  artifacts:
    paths:
      - public/
  environment:
    name: review/$CI_COMMIT_REF_SLUG
    url: https://${CI_PROJECT_NAMESPACE}.gitlab.io/-/${CI_PROJECT_NAME}/-/jobs/${CI_JOB_ID}/artifacts/public/index.html
    on_stop: stop:review
    auto_stop_in: 1 week
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"

# Stop review app
stop:review:
  stage: cleanup
  script:
    - echo "Stopping review app"
  environment:
    name: review/$CI_COMMIT_REF_SLUG
    action: stop
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
      when: manual
  allow_failure: true
```

For more robust review apps, consider deploying to external services like Netlify or Vercel.

## Caching Strategies

### Cache markata-go Binary

```yaml
variables:
  MARKATA_VERSION: "0.1.0"

cache:
  key: markata-go-${MARKATA_VERSION}
  paths:
    - .cache/markata-go

build:
  stage: build
  image: alpine:latest
  before_script:
    - mkdir -p .cache
    - |
      if [ -f .cache/markata-go ]; then
        cp .cache/markata-go /usr/local/bin/
      else
        apk add --no-cache wget tar
        wget -qO- "https://github.com/WaylonWalker/markata-go/releases/download/v${MARKATA_VERSION}/markata-go_${MARKATA_VERSION}_linux_x86_64.tar.gz" | tar xz
        cp markata-go .cache/
        mv markata-go /usr/local/bin/
      fi
  script:
    - markata-go build --clean
```

### Cache with Go Install

If you prefer building from source:

```yaml
variables:
  GOPATH: ${CI_PROJECT_DIR}/.go

cache:
  key: ${CI_COMMIT_REF_SLUG}
  paths:
    - .go/pkg/mod/
    - .go/bin/

build:
  stage: build
  image: golang:1.22
  script:
    - |
      if [ ! -f ${GOPATH}/bin/markata-go ]; then
        go install github.com/WaylonWalker/markata-go/cmd/markata-go@latest
      fi
    - ${GOPATH}/bin/markata-go build --clean
```

### Per-Branch Cache

```yaml
cache:
  key: ${CI_COMMIT_REF_SLUG}
  paths:
    - .cache/

# Fallback to default branch cache if branch cache doesn't exist
cache:
  key:
    files:
      - go.sum
    prefix: ${CI_COMMIT_REF_SLUG}
  paths:
    - .go/pkg/mod/
  policy: pull-push
```

## Scheduled Pipelines

Rebuild your site on a schedule:

1. Go to **CI/CD** > **Schedules**
2. Create a new schedule
3. Set the interval (e.g., daily at midnight)
4. Target the main branch

Or configure in `.gitlab-ci.yml`:

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
    - if: $CI_PIPELINE_SOURCE == "schedule"
```

## Deploying to External Services

### Netlify

```yaml
deploy:netlify:
  stage: deploy
  image: node:20-alpine
  dependencies:
    - build
  before_script:
    - npm install -g netlify-cli
  script:
    - netlify deploy --dir=public --prod
  environment:
    name: production
    url: https://example.netlify.app
  rules:
    - if: $CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH
  variables:
    NETLIFY_AUTH_TOKEN: ${NETLIFY_AUTH_TOKEN}
    NETLIFY_SITE_ID: ${NETLIFY_SITE_ID}
```

Add `NETLIFY_AUTH_TOKEN` and `NETLIFY_SITE_ID` as CI/CD variables in **Settings** > **CI/CD** > **Variables**.

### AWS S3

```yaml
deploy:s3:
  stage: deploy
  image: amazon/aws-cli:latest
  dependencies:
    - build
  script:
    - aws s3 sync ./public s3://${S3_BUCKET} --delete
    - |
      if [ -n "${CLOUDFRONT_DISTRIBUTION_ID}" ]; then
        aws cloudfront create-invalidation --distribution-id ${CLOUDFRONT_DISTRIBUTION_ID} --paths "/*"
      fi
  environment:
    name: production
    url: https://example.com
  rules:
    - if: $CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH
  variables:
    AWS_ACCESS_KEY_ID: ${AWS_ACCESS_KEY_ID}
    AWS_SECRET_ACCESS_KEY: ${AWS_SECRET_ACCESS_KEY}
    AWS_DEFAULT_REGION: us-east-1
```

## Using Docker Images

### Pre-built Image with markata-go

Create a custom Docker image for faster builds:

```dockerfile
# Dockerfile
FROM alpine:latest

ARG MARKATA_VERSION=0.1.0

RUN apk add --no-cache wget tar ca-certificates \
    && wget -qO- "https://github.com/WaylonWalker/markata-go/releases/download/v${MARKATA_VERSION}/markata-go_${MARKATA_VERSION}_linux_x86_64.tar.gz" | tar xz \
    && mv markata-go /usr/local/bin/ \
    && apk del wget tar

ENTRYPOINT ["markata-go"]
```

Use in your pipeline:

```yaml
build:
  stage: build
  image: registry.gitlab.com/mygroup/myproject/markata-go:latest
  script:
    - markata-go build --clean
```

### Go Image

```yaml
build:
  stage: build
  image: golang:1.22-alpine
  before_script:
    - go install github.com/WaylonWalker/markata-go/cmd/markata-go@latest
  script:
    - markata-go build --clean
```

## Troubleshooting

### Pipeline Not Running

- Check `.gitlab-ci.yml` syntax: **CI/CD** > **Editor**
- Verify the file is in the repository root
- Check branch rules match your branch

### Pages Not Updating

1. Verify the `pages` job completed successfully
2. Check artifacts include `public/` directory
3. Wait a few minutes for propagation
4. Check **Settings** > **Pages** for status

### 404 Error on Pages

- Ensure `public/index.html` exists
- Check the URL matches your project path
- Verify GitLab Pages is enabled

### Build Fails with Memory Error

Increase job resources using tags:

```yaml
build:
  tags:
    - high-memory
  script:
    - markata-go build --clean
```

Or contact your GitLab administrator to configure runner resources.

### Permission Denied

Check that protected variables are available:

1. **Settings** > **CI/CD** > **Variables**
2. Uncheck "Protected" or run on protected branches

## Environment Variables Reference

| Variable | Purpose | Example |
|----------|---------|---------|
| `MARKATA_GO_URL` | Site base URL | `https://example.com` |
| `MARKATA_GO_OUTPUT_DIR` | Output directory | `dist` |
| `MARKATA_GO_TITLE` | Site title override | `My Blog` |
| `CI_PROJECT_NAMESPACE` | GitLab group/user | `mygroup` |
| `CI_PROJECT_NAME` | Repository name | `myproject` |
| `CI_COMMIT_BRANCH` | Current branch | `main` |
| `CI_DEFAULT_BRANCH` | Default branch | `main` |
| `CI_PIPELINE_SOURCE` | Pipeline trigger | `push`, `merge_request_event` |

## Next Steps

- [[github-actions|GitHub Actions Guide]] - GitHub Actions workflows
- [[../deployment|Deployment Guide]] - Manual deployment options
- [GitLab CI/CD Documentation](https://docs.gitlab.com/ee/ci/) - Official docs
