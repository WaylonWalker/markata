---
title: "Deploy to Vercel"
description: "Step-by-step guide to deploy markata-go sites on Vercel"
date: 2026-01-24
published: true
tags:
  - documentation
  - deployment
  - vercel
---


Vercel provides zero-configuration deployments with excellent performance and edge capabilities. It's known for its speed and developer experience.

## Prerequisites

- A Vercel account (free tier available)
- Your markata-go site in a Git repository (GitHub, GitLab, or Bitbucket)

## Cost

| Tier | Bandwidth | Builds | Edge Functions | Price |
|------|-----------|--------|----------------|-------|
| Hobby | 100 GB/mo | 6,000 min/mo | 100,000 exec/mo | $0 |
| Pro | 1 TB/mo | Unlimited | 1M exec/mo | $20/mo |
| Enterprise | Custom | Unlimited | Custom | Contact |

The Hobby tier is free but limited to personal, non-commercial use.

## Method 1: Git Integration (Recommended)

### Step 1: Connect Repository

1. Log in to [Vercel](https://vercel.com)
2. Click **Add New** > **Project**
3. Import your Git repository
4. Authorize Vercel to access your repo

### Step 2: Configure Build Settings

Set these options:

| Setting | Value |
|---------|-------|
| Framework Preset | Other |
| Build Command | See below |
| Output Directory | `public` |
| Install Command | `echo "No npm dependencies"` |

**Build Command:**
```bash
curl -sSL https://go.dev/dl/go1.22.0.linux-amd64.tar.gz | tar -xzf - && ./go/bin/go install github.com/WaylonWalker/markata-go/cmd/markata-go@latest && ~/go/bin/markata-go build --clean
```

### Step 3: Set Environment Variables

Add these environment variables:

| Key | Value |
|-----|-------|
| `MARKATA_GO_URL` | `https://your-project.vercel.app` |

### Step 4: Deploy

Click **Deploy**. Vercel will build and deploy your site.

## Method 2: vercel.json Configuration

Create `vercel.json` in your repository root for more control:

```json
{
  "buildCommand": "curl -sSL https://go.dev/dl/go1.22.0.linux-amd64.tar.gz | tar -xzf - && ./go/bin/go install github.com/WaylonWalker/markata-go/cmd/markata-go@latest && ~/go/bin/markata-go build --clean",
  "outputDirectory": "public",
  "installCommand": "echo 'No npm dependencies'",
  "framework": null,
  "headers": [
    {
      "source": "/static/(.*)",
      "headers": [
        {
          "key": "Cache-Control",
          "value": "public, max-age=31536000, immutable"
        }
      ]
    },
    {
      "source": "/(.*).html",
      "headers": [
        {
          "key": "Cache-Control",
          "value": "public, max-age=0, must-revalidate"
        }
      ]
    },
    {
      "source": "/(.*)",
      "headers": [
        {
          "key": "X-Frame-Options",
          "value": "DENY"
        },
        {
          "key": "X-Content-Type-Options",
          "value": "nosniff"
        }
      ]
    }
  ],
  "redirects": [
    {
      "source": "/old-path/:path*",
      "destination": "/new-path/:path*",
      "permanent": true
    }
  ]
}
```

## Method 3: CLI Deployment

Deploy manually using the Vercel CLI:

```bash
# Install CLI
npm install -g vercel

# Login
vercel login

# Build your site locally
markata-go build --clean

# Deploy preview (from public directory)
cd public && vercel

# Deploy to production
cd public && vercel --prod
```

## Custom Domain Setup

### Step 1: Add Domain

1. Go to your project's **Settings** > **Domains**
2. Enter your domain (e.g., `example.com`)
3. Click **Add**

### Step 2: Configure DNS

Vercel provides DNS records to add at your registrar:

**For apex domain:**
```
Type: A
Name: @
Value: 76.76.21.21
```

**For www subdomain:**
```
Type: CNAME
Name: www
Value: cname.vercel-dns.com
```

### Step 3: Verify and Enable HTTPS

Vercel automatically provisions SSL certificates once DNS is configured. This usually takes a few minutes.

### Step 4: Update Configuration

Update your `markata-go.toml`:

```toml
[markata-go]
url = "https://example.com"
```

Set the production environment variable:

```bash
vercel env add MARKATA_GO_URL production
# Enter: https://example.com
```

## Environment Variables

Manage environment variables per environment:

```bash
# Add to production
vercel env add MARKATA_GO_URL production
# Value: https://example.com

# Add to preview (for PR deployments)
vercel env add MARKATA_GO_URL preview
# Value: (leave empty to use Vercel's preview URL)

# Add to development
vercel env add MARKATA_GO_URL development
# Value: http://localhost:8080

# List all variables
vercel env ls

# Pull variables to local .env
vercel env pull
```

## Preview Deployments

Vercel automatically creates preview deployments for every push:

- **Production**: `your-project.vercel.app`
- **Preview**: `your-project-git-branch-username.vercel.app`
- **PR Preview**: Automatically commented on GitHub PRs

## Edge Functions

Add edge functions for dynamic functionality:

```javascript
// api/hello.js
export const config = {
  runtime: 'edge',
};

export default function handler(request) {
  return new Response(JSON.stringify({ message: 'Hello from the Edge!' }), {
    headers: { 'Content-Type': 'application/json' },
  });
}
```

Access at `/api/hello`.

## Analytics

Enable Vercel Analytics for performance insights:

1. Go to **Project Settings** > **Analytics**
2. Click **Enable**

Add the analytics script to your template:

```html
<script defer src="/_vercel/insights/script.js"></script>
```

## Speed Insights

Enable Speed Insights for Core Web Vitals:

1. Go to **Project Settings** > **Speed Insights**
2. Click **Enable**

## Troubleshooting

### Build Fails: Go Installation Issues

The build command needs to install Go first. Ensure your build command is correct:

```json
{
  "buildCommand": "curl -sSL https://go.dev/dl/go1.22.0.linux-amd64.tar.gz | tar -xzf - && ./go/bin/go install github.com/WaylonWalker/markata-go/cmd/markata-go@latest && ~/go/bin/markata-go build --clean"
}
```

### Build Timeout

Hobby tier has a 45-second build limit. Consider:

1. Upgrading to Pro for longer builds
2. Pre-building locally and deploying the output:
   ```bash
   markata-go build --clean
   cd public && vercel --prod
   ```

### 404 on Page Refresh

Vercel handles clean URLs automatically for static sites. Ensure your output directory is correct.

### Environment Variables Not Applied

Check that variables are set for the correct environment:

```bash
vercel env ls
```

Production, Preview, and Development are separate environments.

### Assets Not Loading

Verify your `MARKATA_GO_URL` matches your deployment URL:

```bash
# Check current value
vercel env ls | grep MARKATA_GO_URL
```

## Performance Optimization

### Enable Compression

Vercel automatically compresses responses with gzip and Brotli.

### Configure Caching

In `vercel.json`:

```json
{
  "headers": [
    {
      "source": "/static/(.*)",
      "headers": [
        {
          "key": "Cache-Control",
          "value": "public, max-age=31536000, immutable"
        }
      ]
    },
    {
      "source": "/(.*).xml",
      "headers": [
        {
          "key": "Cache-Control",
          "value": "public, max-age=3600"
        }
      ]
    }
  ]
}
```

### Use Edge Network

Vercel's edge network automatically serves your site from the nearest location to each visitor.

## Monorepo Support

For monorepos, configure the root directory:

```json
{
  "installCommand": "echo 'Skip install'",
  "buildCommand": "cd packages/my-site && curl -sSL ... && markata-go build",
  "outputDirectory": "packages/my-site/public"
}
```

Or set the **Root Directory** in project settings.

## Next Steps

- [Vercel Docs](https://vercel.com/docs) - Official documentation
- [Configuration Guide](../configuration/) - Customize your markata-go site
- [Themes Guide](../themes/) - Change your site's appearance
