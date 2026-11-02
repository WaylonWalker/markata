---
title: "Deploy to Netlify"
description: "Step-by-step guide to deploy markata-go sites on Netlify"
date: 2026-01-24
published: true
tags:
  - documentation
  - deployment
  - netlify
---


Netlify offers continuous deployment, serverless functions, form handling, and a powerful CDN. It's excellent for teams and sites that need more than just static hosting.

## Prerequisites

- A Netlify account (free tier available)
- Your markata-go site in a Git repository (GitHub, GitLab, or Bitbucket)

## Cost

| Tier | Bandwidth | Build Minutes | Forms | Price |
|------|-----------|---------------|-------|-------|
| Free | 100 GB/mo | 300 min/mo | 100/mo | $0 |
| Pro | 1 TB/mo | 25,000 min/mo | Unlimited | $19/mo |
| Business | 1 TB/mo | 25,000 min/mo | Unlimited | $99/mo |

The free tier is generous enough for most personal and small business sites.

## Method 1: Git Integration (Recommended)

### Step 1: Connect Repository

1. Log in to [Netlify](https://app.netlify.com)
2. Click **Add new site** > **Import an existing project**
3. Choose your Git provider and authorize Netlify
4. Select your repository

### Step 2: Configure Build Settings

Enter these settings:

| Setting | Value |
|---------|-------|
| Base directory | (leave empty) |
| Build command | `go install github.com/WaylonWalker/markata-go/cmd/markata-go@latest && markata-go build --clean` |
| Publish directory | `public` |

### Step 3: Set Environment Variables

Click **Advanced** and add:

| Key | Value |
|-----|-------|
| `GO_VERSION` | `1.22` |
| `MARKATA_GO_URL` | `https://your-site.netlify.app` |

### Step 4: Deploy

Click **Deploy site**. Netlify will build and deploy your site automatically.

## Method 2: netlify.toml Configuration

For more control, create `netlify.toml` in your repository root:

```toml
[build]
  command = "go install github.com/WaylonWalker/markata-go/cmd/markata-go@latest && markata-go build --clean"
  publish = "public"

[build.environment]
  GO_VERSION = "1.22"
  MARKATA_GO_URL = "https://your-site.netlify.app"

# Production context
[context.production]
  environment = { MARKATA_GO_URL = "https://example.com" }

# Deploy previews (automatic for PRs)
[context.deploy-preview]
  command = "go install github.com/WaylonWalker/markata-go/cmd/markata-go@latest && markata-go build"

# Branch deploys
[context.branch-deploy]
  command = "go install github.com/WaylonWalker/markata-go/cmd/markata-go@latest && markata-go build"

# Security and caching headers
[[headers]]
  for = "/*"
  [headers.values]
    X-Frame-Options = "DENY"
    X-Content-Type-Options = "nosniff"
    Referrer-Policy = "strict-origin-when-cross-origin"

[[headers]]
  for = "/static/*"
  [headers.values]
    Cache-Control = "public, max-age=31536000, immutable"

[[headers]]
  for = "/*.html"
  [headers.values]
    Cache-Control = "public, max-age=0, must-revalidate"

# Redirects
[[redirects]]
  from = "/old-path/*"
  to = "/new-path/:splat"
  status = 301
```

## Method 3: CLI Deployment

Deploy manually using the Netlify CLI:

```bash
# Install CLI
npm install -g netlify-cli

# Login
netlify login

# Initialize project (first time)
netlify init

# Build your site
markata-go build --clean

# Deploy preview
netlify deploy --dir=public

# Deploy to production
netlify deploy --dir=public --prod
```

## Custom Domain Setup

### Step 1: Add Domain in Netlify

1. Go to **Site settings** > **Domain management**
2. Click **Add custom domain**
3. Enter your domain (e.g., `example.com`)

### Step 2: Configure DNS

**Option A: Use Netlify DNS (Recommended)**

1. Click **Set up Netlify DNS** for your domain
2. Update your domain registrar's nameservers to Netlify's:
   ```
   dns1.p01.nsone.net
   dns2.p01.nsone.net
   dns3.p01.nsone.net
   dns4.p01.nsone.net
   ```

**Option B: External DNS**

Add these records at your DNS provider:

For apex domain:
```
Type: A
Name: @
Value: 75.2.60.5
```

For www subdomain:
```
Type: CNAME
Name: www
Value: your-site.netlify.app
```

### Step 3: Enable HTTPS

Netlify automatically provisions Let's Encrypt certificates. Go to **Domain management** > **HTTPS** and click **Verify DNS configuration**.

### Step 4: Update Configuration

Update your `markata-go.toml`:

```toml
[markata-go]
url = "https://example.com"
```

And your `netlify.toml`:

```toml
[context.production.environment]
  MARKATA_GO_URL = "https://example.com"
```

## Deploy Previews

Netlify automatically creates preview deployments for pull requests. Each PR gets a unique URL like:

```
https://deploy-preview-123--your-site.netlify.app
```

To customize preview behavior, add to `netlify.toml`:

```toml
[context.deploy-preview]
  command = "go install github.com/WaylonWalker/markata-go/cmd/markata-go@latest && markata-go build"
  
[context.deploy-preview.environment]
  # Leave MARKATA_GO_URL empty to use Netlify's preview URL
```

## Forms (No Backend Required)

Netlify can handle form submissions without a backend:

```html
<form name="contact" method="POST" data-netlify="true">
  <input type="hidden" name="form-name" value="contact" />
  <input type="text" name="name" required />
  <input type="email" name="email" required />
  <textarea name="message" required></textarea>
  <button type="submit">Send</button>
</form>
```

Form submissions appear in **Site settings** > **Forms**.

## Serverless Functions

Add serverless functions in `netlify/functions/`:

```javascript
// netlify/functions/hello.js
exports.handler = async (event, context) => {
  return {
    statusCode: 200,
    body: JSON.stringify({ message: "Hello from Netlify Functions!" })
  };
};
```

Access at `/.netlify/functions/hello`.

## Troubleshooting

### Build Fails: Go Not Found

Ensure `GO_VERSION` is set in environment variables:

```toml
[build.environment]
  GO_VERSION = "1.22"
```

### Build Timeout

Free tier has a 15-minute build limit. Optimize your build:

```toml
[build]
  command = """
    if [ ! -f $HOME/go/bin/markata-go ]; then
      go install github.com/WaylonWalker/markata-go/cmd/markata-go@latest
    fi
    markata-go build --clean
  """
```

### Assets Not Loading

Check your site URL configuration:

```bash
# View current config
netlify env:list

# Set production URL
netlify env:set MARKATA_GO_URL https://example.com
```

### Deploy Preview URL Issues

For deploy previews, let Netlify handle the URL:

```toml
[context.deploy-preview.environment]
  # Don't set MARKATA_GO_URL - let Netlify use the preview URL
```

### Custom Headers Not Applied

Ensure `netlify.toml` is in your repository root (not in a subdirectory).

## Performance Optimization

### Enable Asset Optimization

In **Site settings** > **Build & deploy** > **Post processing**:

- Enable **Pretty URLs**
- Enable **Asset optimization** (minify CSS/JS)
- Enable **Prerendering**

### Configure Caching

```toml
[[headers]]
  for = "/static/*"
  [headers.values]
    Cache-Control = "public, max-age=31536000, immutable"

[[headers]]
  for = "/*.xml"
  [headers.values]
    Cache-Control = "public, max-age=3600"
```

## Next Steps

- [Netlify Docs](https://docs.netlify.com/) - Official documentation
- [Configuration Guide](../configuration/) - Customize your markata-go site
- [Feeds Guide](../feeds/) - Set up RSS and Atom feeds
