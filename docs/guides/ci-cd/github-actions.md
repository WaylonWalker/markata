---
title: "GitHub Actions CI/CD"
description: "Complete guide to building and deploying markata-go sites with GitHub Actions"
date: 2026-01-24
published: true
slug: /docs/guides/ci-cd/github-actions/
tags:
  - documentation
  - ci-cd
  - deployment
  - github-actions
---


GitHub Actions provides powerful CI/CD capabilities directly integrated with your GitHub repository. This guide covers everything from basic deployments to advanced multi-environment setups.

## Quick Start

Create `.github/workflows/deploy.yml` in your repository:

```yaml
name: Deploy Site

on:
  push:
    branches: [main]
  workflow_dispatch:

permissions:
  contents: read
  pages: write
  id-token: write

jobs:
  build-and-deploy:
    runs-on: ubuntu-latest
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    steps:
      - uses: actions/checkout@v4

      - name: Install markata-go
        run: |
          wget -qO- "https://github.com/WaylonWalker/markata-go/releases/latest/download/markata-go_linux_x86_64.tar.gz" | tar xz
          sudo mv markata-go /usr/local/bin/

      - name: Build site
        run: markata-go build --clean
        env:
          MARKATA_GO_URL: https://${{ github.repository_owner }}.github.io/${{ github.event.repository.name }}

      - name: Setup Pages
        uses: actions/configure-pages@v4

      - name: Upload artifact
        uses: actions/upload-pages-artifact@v3
        with:
          path: ./public

      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

Enable GitHub Pages in your repository settings:

1. Go to **Settings** > **Pages**
2. Under **Source**, select **GitHub Actions**
3. Push to `main` to trigger the workflow

## GitHub Pages

### Basic Deployment

The workflow above deploys to GitHub Pages using the modern Actions-based approach. Key components:

**Permissions** - Required for the `deploy-pages` action:

```yaml
permissions:
  contents: read    # Read repository
  pages: write      # Deploy to Pages
  id-token: write   # OIDC token for deployment
```

**Environment** - Links the job to the GitHub Pages environment:

```yaml
environment:
  name: github-pages
  url: ${{ steps.deployment.outputs.page_url }}
```

### Custom Domain

To use a custom domain:

1. Create `static/CNAME` with your domain:

```
example.com
```

2. Update your workflow:

```yaml
- name: Build site
  run: markata-go build --clean
  env:
    MARKATA_GO_URL: https://example.com
```

3. Configure DNS at your registrar:
   - **Apex domain**: Add `A` records pointing to GitHub's IPs
   - **Subdomain**: Add a `CNAME` record to `<username>.github.io`

### Project Sites vs User Sites

| Type | Repository Name | URL | Branch |
|------|----------------|-----|--------|
| **User/Org** | `username.github.io` | `https://username.github.io` | `main` |
| **Project** | Any name | `https://username.github.io/repo-name` | `main` |

For project sites, ensure your base URL includes the repository name:

```yaml
env:
  MARKATA_GO_URL: https://username.github.io/repo-name
```

## Netlify

Deploy to Netlify using GitHub Actions for more control over the build process.

### Basic Netlify Deployment

```yaml
name: Deploy to Netlify

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  build-and-deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install markata-go
        run: |
          wget -qO- "https://github.com/WaylonWalker/markata-go/releases/latest/download/markata-go_linux_x86_64.tar.gz" | tar xz
          sudo mv markata-go /usr/local/bin/

      - name: Build site
        run: markata-go build --clean
        env:
          MARKATA_GO_URL: ${{ github.event_name == 'push' && 'https://example.netlify.app' || '' }}

      - name: Deploy to Netlify
        uses: nwtgck/actions-netlify@v3
        with:
          publish-dir: ./public
          production-branch: main
          production-deploy: ${{ github.event_name == 'push' }}
          deploy-message: "Deploy from GitHub Actions"
          github-token: ${{ secrets.GITHUB_TOKEN }}
        env:
          NETLIFY_AUTH_TOKEN: ${{ secrets.NETLIFY_AUTH_TOKEN }}
          NETLIFY_SITE_ID: ${{ secrets.NETLIFY_SITE_ID }}
```

**Setup:**

1. Create a Netlify site and note the Site ID
2. Generate a Personal Access Token in Netlify (User Settings > Applications)
3. Add secrets to your repository:
   - `NETLIFY_AUTH_TOKEN` - Your personal access token
   - `NETLIFY_SITE_ID` - Your site's API ID

### Netlify with Deploy Previews

The workflow above automatically creates deploy previews for pull requests. The preview URL is posted as a comment on the PR.

## Cloudflare Pages

Deploy to Cloudflare Pages for fast global CDN delivery.

```yaml
name: Deploy to Cloudflare Pages

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  build-and-deploy:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      deployments: write
      pull-requests: write
    steps:
      - uses: actions/checkout@v4

      - name: Install markata-go
        run: |
          wget -qO- "https://github.com/WaylonWalker/markata-go/releases/latest/download/markata-go_linux_x86_64.tar.gz" | tar xz
          sudo mv markata-go /usr/local/bin/

      - name: Build site
        run: markata-go build --clean
        env:
          MARKATA_GO_URL: ${{ github.event_name == 'push' && 'https://example.pages.dev' || '' }}

      - name: Deploy to Cloudflare Pages
        uses: cloudflare/pages-action@v1
        with:
          apiToken: ${{ secrets.CLOUDFLARE_API_TOKEN }}
          accountId: ${{ secrets.CLOUDFLARE_ACCOUNT_ID }}
          projectName: my-site
          directory: public
          gitHubToken: ${{ secrets.GITHUB_TOKEN }}
```

**Setup:**

1. Create a Cloudflare Pages project
2. Generate an API token with "Cloudflare Pages:Edit" permission
3. Add secrets:
   - `CLOUDFLARE_API_TOKEN` - Your API token
   - `CLOUDFLARE_ACCOUNT_ID` - Your account ID

## AWS S3

Deploy to Amazon S3 for scalable static hosting.

```yaml
name: Deploy to S3

on:
  push:
    branches: [main]

jobs:
  build-and-deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install markata-go
        run: |
          wget -qO- "https://github.com/WaylonWalker/markata-go/releases/latest/download/markata-go_linux_x86_64.tar.gz" | tar xz
          sudo mv markata-go /usr/local/bin/

      - name: Build site
        run: markata-go build --clean
        env:
          MARKATA_GO_URL: https://example.com

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v4
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: us-east-1

      - name: Sync to S3
        run: |
          aws s3 sync ./public s3://${{ vars.S3_BUCKET }} \
            --delete \
            --cache-control "max-age=31536000" \
            --exclude "*.html" \
            --exclude "*.xml"
          
          # HTML and XML with shorter cache
          aws s3 sync ./public s3://${{ vars.S3_BUCKET }} \
            --cache-control "max-age=0, must-revalidate" \
            --include "*.html" \
            --include "*.xml"

      - name: Invalidate CloudFront
        if: vars.CLOUDFRONT_DISTRIBUTION_ID != ''
        run: |
          aws cloudfront create-invalidation \
            --distribution-id ${{ vars.CLOUDFRONT_DISTRIBUTION_ID }} \
            --paths "/*"
```

**Setup:**

1. Create an S3 bucket configured for static website hosting
2. Create an IAM user with S3 and CloudFront permissions
3. Add secrets and variables:
   - `AWS_ACCESS_KEY_ID` - IAM access key
   - `AWS_SECRET_ACCESS_KEY` - IAM secret key
   - `S3_BUCKET` - Bucket name (variable)
   - `CLOUDFRONT_DISTRIBUTION_ID` - Optional CloudFront distribution

## Preview Deployments

Create preview deployments for pull requests to review changes before merging.

### GitHub Pages Preview (Custom Approach)

Deploy PR previews to subdirectories on GitHub Pages:

```yaml
name: PR Preview

on:
  pull_request:
    types: [opened, synchronize, reopened]

permissions:
  contents: write
  pull-requests: write

jobs:
  preview:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install markata-go
        run: |
          wget -qO- "https://github.com/WaylonWalker/markata-go/releases/latest/download/markata-go_linux_x86_64.tar.gz" | tar xz
          sudo mv markata-go /usr/local/bin/

      - name: Build site
        run: markata-go build --clean
        env:
          MARKATA_GO_URL: https://${{ github.repository_owner }}.github.io/${{ github.event.repository.name }}/pr-${{ github.event.number }}

      - name: Deploy Preview
        uses: peaceiris/actions-gh-pages@v4
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./public
          destination_dir: pr-${{ github.event.number }}

      - name: Comment on PR
        uses: actions/github-script@v7
        with:
          script: |
            const url = `https://${{ github.repository_owner }}.github.io/${{ github.event.repository.name }}/pr-${{ github.event.number }}/`;
            const body = `## Preview Deployment\n\nYour preview is ready!\n\n${url}`;
            
            // Find existing comment
            const { data: comments } = await github.rest.issues.listComments({
              owner: context.repo.owner,
              repo: context.repo.repo,
              issue_number: context.issue.number,
            });
            
            const botComment = comments.find(c => 
              c.user.type === 'Bot' && c.body.includes('Preview Deployment')
            );
            
            if (botComment) {
              await github.rest.issues.updateComment({
                owner: context.repo.owner,
                repo: context.repo.repo,
                comment_id: botComment.id,
                body: body,
              });
            } else {
              await github.rest.issues.createComment({
                owner: context.repo.owner,
                repo: context.repo.repo,
                issue_number: context.issue.number,
                body: body,
              });
            }
```

### Cleanup Preview on PR Close

```yaml
name: Cleanup PR Preview

on:
  pull_request:
    types: [closed]

permissions:
  contents: write

jobs:
  cleanup:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          ref: gh-pages

      - name: Remove preview directory
        run: |
          rm -rf pr-${{ github.event.number }}
          git config user.name "github-actions[bot]"
          git config user.email "github-actions[bot]@users.noreply.github.com"
          git add -A
          git commit -m "Remove preview for PR #${{ github.event.number }}" || exit 0
          git push
```

## Multi-Environment Deployments

Deploy to staging and production environments with different configurations.

### Staging and Production

```yaml
name: Deploy

on:
  push:
    branches:
      - main
      - develop

jobs:
  build:
    runs-on: ubuntu-latest
    outputs:
      artifact-id: ${{ steps.upload.outputs.artifact-id }}
    steps:
      - uses: actions/checkout@v4

      - name: Install markata-go
        run: |
          wget -qO- "https://github.com/WaylonWalker/markata-go/releases/latest/download/markata-go_linux_x86_64.tar.gz" | tar xz
          sudo mv markata-go /usr/local/bin/

      - name: Determine environment
        id: env
        run: |
          if [[ "${{ github.ref }}" == "refs/heads/main" ]]; then
            echo "name=production" >> $GITHUB_OUTPUT
            echo "url=https://example.com" >> $GITHUB_OUTPUT
          else
            echo "name=staging" >> $GITHUB_OUTPUT
            echo "url=https://staging.example.com" >> $GITHUB_OUTPUT
          fi

      - name: Build site
        run: markata-go build --clean
        env:
          MARKATA_GO_URL: ${{ steps.env.outputs.url }}

      - name: Upload artifact
        id: upload
        uses: actions/upload-artifact@v4
        with:
          name: site-${{ steps.env.outputs.name }}
          path: ./public
          retention-days: 1

  deploy-staging:
    needs: build
    if: github.ref == 'refs/heads/develop'
    runs-on: ubuntu-latest
    environment:
      name: staging
      url: https://staging.example.com
    steps:
      - name: Download artifact
        uses: actions/download-artifact@v4
        with:
          name: site-staging
          path: ./public

      - name: Deploy to staging
        run: |
          # Your staging deployment command
          echo "Deploying to staging..."

  deploy-production:
    needs: build
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    environment:
      name: production
      url: https://example.com
    steps:
      - name: Download artifact
        uses: actions/download-artifact@v4
        with:
          name: site-production
          path: ./public

      - name: Deploy to production
        run: |
          # Your production deployment command
          echo "Deploying to production..."
```

### Manual Production Deployment

Require manual approval for production deployments:

1. Go to **Settings** > **Environments**
2. Create a `production` environment
3. Enable **Required reviewers**
4. Add team members who can approve

The workflow will pause at the `deploy-production` job until approved.

## Caching Strategies

Speed up builds by caching dependencies and build artifacts.

### Cache markata-go Binary

```yaml
- name: Cache markata-go
  id: cache-markata
  uses: actions/cache@v4
  with:
    path: /usr/local/bin/markata-go
    key: markata-go-${{ runner.os }}-v0.1.0

- name: Install markata-go
  if: steps.cache-markata.outputs.cache-hit != 'true'
  run: |
    wget -qO- "https://github.com/WaylonWalker/markata-go/releases/download/v0.1.0/markata-go_0.1.0_linux_x86_64.tar.gz" | tar xz
    sudo mv markata-go /usr/local/bin/
```

### Cache Go Modules (if building from source)

```yaml
- name: Setup Go
  uses: actions/setup-go@v5
  with:
    go-version: '1.22'
    cache: true

- name: Install markata-go
  run: go install github.com/WaylonWalker/markata-go/cmd/markata-go@latest
```

### Full Workflow with Caching

```yaml
name: Deploy with Caching

on:
  push:
    branches: [main]

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: pages
  cancel-in-progress: true

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0  # Full history for git dates

      - name: Cache markata-go
        id: cache-markata
        uses: actions/cache@v4
        with:
          path: /usr/local/bin/markata-go
          key: markata-go-${{ runner.os }}-v0.1.0

      - name: Install markata-go
        if: steps.cache-markata.outputs.cache-hit != 'true'
        run: |
          wget -qO- "https://github.com/WaylonWalker/markata-go/releases/download/v0.1.0/markata-go_0.1.0_linux_x86_64.tar.gz" | tar xz
          sudo mv markata-go /usr/local/bin/

      - name: Build site
        run: markata-go build --clean
        env:
          MARKATA_GO_URL: https://example.com

      - name: Upload artifact
        uses: actions/upload-pages-artifact@v3
        with:
          path: ./public

  deploy:
    needs: build
    runs-on: ubuntu-latest
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    steps:
      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

## Scheduled Builds

Rebuild your site on a schedule (useful for dynamic content or expired dates):

```yaml
name: Scheduled Build

on:
  schedule:
    # Run daily at midnight UTC
    - cron: '0 0 * * *'
  workflow_dispatch:  # Allow manual trigger

permissions:
  contents: read
  pages: write
  id-token: write

jobs:
  build-and-deploy:
    runs-on: ubuntu-latest
    steps:
      # ... same as basic deployment
```

## Troubleshooting

### Workflow Not Running

- Check that the workflow file is in `.github/workflows/`
- Verify the branch name matches your trigger
- Check **Actions** tab for any errors

### Permission Denied

Ensure your workflow has the required permissions:

```yaml
permissions:
  contents: read
  pages: write
  id-token: write
```

### Deployment Shows Old Content

1. Check that `--clean` flag is used
2. Verify the build completed successfully
3. Check for caching issues - try clearing the cache

### 404 on GitHub Pages

- Verify GitHub Pages is enabled in repository settings
- Check that `index.html` exists in the output
- For project sites, ensure base URL includes repository name

### Build Fails with Memory Error

For large sites, increase available memory:

```yaml
jobs:
  build:
    runs-on: ubuntu-latest
    env:
      GOMEMLIMIT: 4GiB
```

## Next Steps

- [[gitlab-ci|GitLab CI Guide]] - GitLab CI/CD pipelines
- [[../deployment|Deployment Guide]] - Manual deployment options
- [GitHub Actions Documentation](https://docs.github.com/en/actions) - Official docs
