---
title: "Deploy to GitHub Pages"
description: "Step-by-step guide to deploy markata-go sites on GitHub Pages"
date: 2026-01-24
published: true
tags:
  - documentation
  - deployment
  - github-pages
---


GitHub Pages provides free hosting for static sites directly from a GitHub repository. It's the easiest way to deploy if your code is already on GitHub.

## Prerequisites

- A GitHub account
- Your markata-go site in a GitHub repository
- Git installed locally

## Cost

| Tier | Storage | Bandwidth | Custom Domain |
|------|---------|-----------|---------------|
| Free | 1 GB | 100 GB/mo | Yes (HTTPS) |

GitHub Pages is completely free for public repositories. Private repositories require GitHub Pro ($4/mo) or a paid organization plan.

## Method 1: GitHub Actions (Recommended)

This method automatically builds and deploys your site on every push to `main`.

### Step 1: Create the Workflow

Create `.github/workflows/deploy.yml`:

```yaml
name: Deploy to GitHub Pages

on:
  push:
    branches:
      - main
  workflow_dispatch:

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: "pages"
  cancel-in-progress: false

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
        run: go install github.com/WaylonWalker/markata-go/cmd/markata-go@latest

      - name: Build site
        run: markata-go build --clean
        env:
          MARKATA_GO_URL: https://${{ github.repository_owner }}.github.io/${{ github.event.repository.name }}

      - name: Upload artifact
        uses: actions/upload-pages-artifact@v3
        with:
          path: ./public

  deploy:
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    runs-on: ubuntu-latest
    needs: build
    steps:
      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

### Step 2: Configure Repository Settings

1. Go to your repository on GitHub
2. Navigate to **Settings** > **Pages**
3. Under **Source**, select **GitHub Actions**
4. Push your changes to trigger the workflow

### Step 3: Access Your Site

After the workflow completes, your site will be available at:

```
https://<username>.github.io/<repository-name>/
```

## Method 2: Manual Deployment

For more control or if you can't use GitHub Actions.

### Step 1: Build Locally

```bash
# Set the URL for your GitHub Pages site
export MARKATA_GO_URL=https://username.github.io/repo-name

# Build the site
markata-go build --clean
```

### Step 2: Deploy to gh-pages Branch

```bash
# Navigate to output directory
cd public

# Initialize git and push to gh-pages branch
git init
git add -A
git commit -m "Deploy site"
git push -f git@github.com:username/repo-name.git main:gh-pages
```

Or use the `gh-pages` npm package:

```bash
npm install -g gh-pages
markata-go build --clean
gh-pages -d public
```

### Step 3: Configure Repository

1. Go to **Settings** > **Pages**
2. Under **Source**, select **Deploy from a branch**
3. Select the `gh-pages` branch and `/ (root)` folder
4. Click **Save**

## Custom Domain Setup

### Step 1: Add CNAME File

Create `static/CNAME` (no extension) with your domain:

```
example.com
```

This file will be copied to `public/CNAME` during build.

### Step 2: Update Configuration

In your `markata-go.toml`:

```toml
[markata-go]
url = "https://example.com"
```

### Step 3: Configure DNS

**For apex domain (example.com):**

Add these A records pointing to GitHub's servers:

```
185.199.108.153
185.199.109.153
185.199.110.153
185.199.111.153
```

**For subdomain (www.example.com or blog.example.com):**

Add a CNAME record:

```
Type: CNAME
Name: www (or blog)
Value: username.github.io
```

### Step 4: Enable HTTPS

1. Go to **Settings** > **Pages**
2. Under **Custom domain**, enter your domain
3. Check **Enforce HTTPS** (may take up to 24 hours)

## Troubleshooting

### Build Fails: "command not found: markata-go"

Ensure Go is set up correctly in your workflow:

```yaml
- name: Setup Go
  uses: actions/setup-go@v5
  with:
    go-version: '1.22'
    cache: true

- name: Install markata-go
  run: |
    go install github.com/WaylonWalker/markata-go/cmd/markata-go@latest
    echo "$(go env GOPATH)/bin" >> $GITHUB_PATH
```

### 404 Errors on Subpages

Make sure your site URL matches the GitHub Pages URL:

```bash
# For project sites (username.github.io/repo-name)
MARKATA_GO_URL=https://username.github.io/repo-name markata-go build

# For user sites (username.github.io)
MARKATA_GO_URL=https://username.github.io markata-go build
```

### CSS/JS Not Loading

Check that asset paths are relative or use the correct base URL. The workflow above sets `MARKATA_GO_URL` automatically.

### Custom Domain Not Working

1. Verify DNS propagation: `dig example.com`
2. Ensure CNAME file exists in output
3. Wait up to 24 hours for HTTPS certificate

### Workflow Not Triggering

Ensure the workflow file is in the correct location:
```
.github/workflows/deploy.yml
```

And that you're pushing to the `main` branch (or adjust the workflow trigger).

## Advanced Configuration

### Caching for Faster Builds

Add caching for the markata-go binary:

```yaml
- name: Cache markata-go
  uses: actions/cache@v4
  with:
    path: ~/go/bin/markata-go
    key: markata-go-${{ runner.os }}-v1

- name: Install markata-go
  run: |
    if [ ! -f ~/go/bin/markata-go ]; then
      go install github.com/WaylonWalker/markata-go/cmd/markata-go@latest
    fi
```

### Preview Deployments for Pull Requests

Deploy preview versions for PRs:

```yaml
name: Preview

on:
  pull_request:
    types: [opened, synchronize]

jobs:
  preview:
    runs-on: ubuntu-latest
    permissions:
      contents: write
      pull-requests: write
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'

      - name: Build
        run: |
          go install github.com/WaylonWalker/markata-go/cmd/markata-go@latest
          markata-go build --clean
        env:
          MARKATA_GO_URL: https://username.github.io/repo/pr-${{ github.event.number }}

      - name: Deploy Preview
        uses: peaceiris/actions-gh-pages@v4
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./public
          destination_dir: pr-${{ github.event.number }}

      - name: Comment PR
        uses: actions/github-script@v7
        with:
          script: |
            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: 'Preview: https://username.github.io/repo/pr-${{ github.event.number }}/'
            })
```

## Next Steps

- [Custom Domain Setup](https://docs.github.com/en/pages/configuring-a-custom-domain-for-your-github-pages-site) - GitHub's official docs
- [Configuration Guide](../configuration/) - Customize your markata-go site
- [Themes Guide](../themes/) - Change your site's appearance
