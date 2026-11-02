---
title: "CI/CD Deployment Guide"
description: "Automate building and deploying markata-go sites with continuous integration pipelines"
date: 2026-01-24
published: true
slug: /docs/guides/ci-cd/
tags:
  - documentation
  - ci-cd
  - deployment
---


Automate your markata-go site builds and deployments using continuous integration (CI) and continuous deployment (CD) pipelines. This guide covers the most popular platforms and deployment patterns.

## Why CI/CD?

Setting up automated deployments provides several benefits:

- **Consistency** - Every deployment uses the same build process
- **Speed** - Sites deploy automatically on every push
- **Safety** - Preview deployments let you review changes before production
- **Collaboration** - Team members can contribute without manual deployment steps

## Platform Guides

Choose your CI/CD platform:

| Platform | Guide | Best For |
|----------|-------|----------|
| **GitHub Actions** | [[github-actions|GitHub Actions Guide]] | GitHub repositories, GitHub Pages |
| **GitLab CI** | [[gitlab-ci|GitLab CI Guide]] | GitLab repositories, GitLab Pages |

## Quick Start

The basic CI/CD workflow for markata-go is:

1. **Install markata-go** - Download the binary from GitHub releases
2. **Build the site** - Run `markata-go build --clean`
3. **Deploy** - Upload the `public/` directory to your hosting platform

### Installing markata-go in CI

For container-based CI pipelines, the easiest option is the builder image:

```bash
docker run --rm \
  -v "$PWD":/site \
  -w /site \
  ghcr.io/waylonwalker/markata-go-builder:latest \
  sh -c 'markata-go build --clean'
```

If you prefer installing a binary directly, download the pre-built release:

```bash
# Download latest release (Linux x86_64)
MARKATA_VERSION="v0.1.0"
wget -qO- "https://github.com/WaylonWalker/markata-go/releases/download/${MARKATA_VERSION}/markata-go_${MARKATA_VERSION#v}_linux_x86_64.tar.gz" | tar xz

# Run the build
./markata-go build --clean
```

Alternatively, if Go is available in your CI environment:

```bash
go install github.com/WaylonWalker/markata-go/cmd/markata-go@latest
markata-go build --clean
```

### Setting the Site URL

For production deployments, always set the site URL via environment variable:

```bash
MARKATA_GO_URL=https://example.com markata-go build --clean
```

This ensures all absolute URLs in feeds, sitemaps, and canonical tags are correct.

## Deployment Targets

markata-go generates static files that can be deployed anywhere. Common targets include:

| Target | Description | Guide Section |
|--------|-------------|---------------|
| **GitHub Pages** | Free hosting from GitHub | [[github-actions#github-pages|GitHub Actions]] |
| **GitLab Pages** | Free hosting from GitLab | [[gitlab-ci#gitlab-pages|GitLab CI]] |
| **Netlify** | CDN with deploy previews | [[github-actions#netlify|GitHub Actions]] |
| **Cloudflare Pages** | Fast global CDN | [[github-actions#cloudflare-pages|GitHub Actions]] |
| **AWS S3** | Scalable object storage | [[github-actions#aws-s3|GitHub Actions]] |
| **Self-hosted** | Your own server | [[deployment-guide|Deployment Guide]] |

## Key Concepts

### Build Artifacts

After running `markata-go build`, the `public/` directory contains your complete static site. This directory should be:

- **Uploaded as an artifact** between CI jobs (if using separate build/deploy jobs)
- **Deployed to your hosting platform** (GitHub Pages, Netlify, etc.)

### Environment Variables

markata-go supports environment variable configuration with the `MARKATA_GO_` prefix:

| Variable | Purpose | Example |
|----------|---------|---------|
| `MARKATA_GO_URL` | Base URL for the site | `https://example.com` |
| `MARKATA_GO_OUTPUT_DIR` | Output directory | `dist` |
| `MARKATA_GO_TITLE` | Site title | `My Blog` |

### Caching

To speed up CI builds, cache the markata-go binary and any Go modules:

```yaml
# GitHub Actions example
- uses: actions/cache@v4
  with:
    path: |
      ~/go/bin/markata-go
      ~/go/pkg/mod
    key: markata-go-${{ runner.os }}-${{ hashFiles('go.sum') }}
```

See the platform-specific guides for detailed caching configurations.

## Multi-Environment Deployments

For projects with staging and production environments:

```yaml
# Simplified example
jobs:
  build:
    # Build once, deploy to multiple environments

  deploy-staging:
    needs: build
    if: github.ref == 'refs/heads/develop'
    environment: staging

  deploy-production:
    needs: build
    if: github.ref == 'refs/heads/main'
    environment: production
```

See [[github-actions#multi-environment|Multi-Environment Deployments]] for complete examples.

## Preview Deployments

Preview deployments create temporary sites for pull requests, allowing reviewers to see changes before merging:

- **Netlify** - Automatic deploy previews for PRs
- **Vercel** - Automatic preview deployments
- **GitHub Pages** - Custom workflow with PR-specific paths

See [[github-actions#preview-deployments|Preview Deployments]] for implementation details.

## Troubleshooting

### Common Issues

**Build fails with "command not found"**

Ensure markata-go is installed and in your PATH:

```bash
# Check if installed
which markata-go || echo "Not found"

# If using wget install, run from current directory
./markata-go build
```

**Assets have wrong URLs**

Set `MARKATA_GO_URL` to match your deployment URL:

```bash
MARKATA_GO_URL=https://username.github.io/repo-name markata-go build
```

**Build succeeds but site is empty**

Check that your content directory contains Markdown files:

```bash
# List markdown files
find . -name "*.md" -type f
```

**RSS feeds show localhost URLs**

The base URL must be set for feeds to work correctly:

```bash
MARKATA_GO_URL=https://example.com markata-go build
```

### Debug Mode

Enable verbose output to diagnose build issues:

```bash
markata-go build --clean -v
```

### Validating Configuration

Before deploying, validate your configuration:

```bash
markata-go config validate
markata-go config show --sources
```

## Next Steps

- [[github-actions|GitHub Actions Guide]] - Complete GitHub Actions workflows
- [[gitlab-ci|GitLab CI Guide]] - GitLab CI/CD pipelines
- [[deployment-guide|Deployment Guide]] - Manual deployment options
- [[troubleshooting|Troubleshooting]] - Common issues and solutions
