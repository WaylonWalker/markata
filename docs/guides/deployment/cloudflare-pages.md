---
title: "Deploy to Cloudflare Pages"
description: "Step-by-step guide to deploy markata-go sites on Cloudflare Pages"
date: 2026-01-24
published: true
tags:
  - documentation
  - deployment
  - cloudflare-pages
---


Cloudflare Pages provides fast global CDN deployment with unlimited bandwidth and Workers integration. It's an excellent choice for high-traffic sites.

## Prerequisites

- A Cloudflare account (free tier available)
- Your markata-go site in a Git repository (GitHub or GitLab)

## Cost

| Tier | Bandwidth | Builds | Sites | Price |
|------|-----------|--------|-------|-------|
| Free | Unlimited | 500/mo | Unlimited | $0 |
| Pro | Unlimited | 5,000/mo | Unlimited | $20/mo |
| Business | Unlimited | 20,000/mo | Unlimited | $200/mo |

The free tier offers unlimited bandwidth, making it ideal for high-traffic sites.

## Method 1: Git Integration (Recommended)

### Step 1: Connect Repository

1. Log in to the [Cloudflare Dashboard](https://dash.cloudflare.com)
2. Go to **Workers & Pages** > **Create application** > **Pages**
3. Click **Connect to Git**
4. Select your Git provider and repository

### Step 2: Configure Build Settings

| Setting | Value |
|---------|-------|
| Framework preset | None |
| Build command | See below |
| Build output directory | `public` |
| Root directory | (leave empty) |

**Build Command:**
```bash
curl -sSL https://go.dev/dl/go1.22.0.linux-amd64.tar.gz | tar -xzf - && ./go/bin/go install github.com/WaylonWalker/markata-go/cmd/markata-go@latest && ~/go/bin/markata-go build --clean
```

### Step 3: Set Environment Variables

Click **Environment variables** and add:

| Variable | Value |
|----------|-------|
| `MARKATA_GO_URL` | `https://your-project.pages.dev` |

Set this for both **Production** and **Preview** environments.

### Step 4: Deploy

Click **Save and Deploy**. Cloudflare will build and deploy your site.

## Method 2: Wrangler CLI

Deploy using the Wrangler CLI for more control:

```bash
# Install wrangler
npm install -g wrangler

# Login to Cloudflare
wrangler login

# Build your site locally
MARKATA_GO_URL=https://example.com markata-go build --clean

# Create a new Pages project (first time)
wrangler pages project create my-site

# Deploy
wrangler pages deploy public --project-name=my-site
```

For subsequent deployments:

```bash
markata-go build --clean
wrangler pages deploy public --project-name=my-site
```

## Custom Domain Setup

### Step 1: Add Domain to Cloudflare

If your domain isn't already on Cloudflare:

1. Go to **Websites** > **Add a Site**
2. Enter your domain
3. Update nameservers at your registrar to Cloudflare's

### Step 2: Connect Domain to Pages

1. Go to your Pages project
2. Click **Custom domains** > **Set up a custom domain**
3. Enter your domain (e.g., `example.com`)
4. Click **Activate domain**

Cloudflare automatically configures DNS and provisions SSL certificates.

### Step 3: Add www Subdomain (Optional)

1. Add `www.example.com` as another custom domain
2. Or create a redirect rule in **Rules** > **Redirect Rules**:
   - If hostname equals `www.example.com`
   - Then redirect to `https://example.com`

### Step 4: Update Configuration

Update your `markata-go.toml`:

```toml
[markata-go]
url = "https://example.com"
```

## Headers and Redirects

### Using _headers File

Create `static/_headers` (copied to output during build):

```
/*
  X-Frame-Options: DENY
  X-Content-Type-Options: nosniff
  Referrer-Policy: strict-origin-when-cross-origin
  Permissions-Policy: accelerometer=(), camera=(), geolocation=(), gyroscope=(), magnetometer=(), microphone=(), payment=(), usb=()

/static/*
  Cache-Control: public, max-age=31536000, immutable

/*.html
  Cache-Control: public, max-age=0, must-revalidate

/*.xml
  Cache-Control: public, max-age=3600
```

### Using _redirects File

Create `static/_redirects`:

```
# Redirect old URLs
/old-path/* /new-path/:splat 301
/blog/old-post /blog/new-post 301

# Proxy API requests
/api/* https://api.example.com/:splat 200
```

## Preview Deployments

Cloudflare Pages automatically creates preview deployments for:

- Every branch push
- Every pull request

Preview URLs follow this pattern:
```
<commit-hash>.<project-name>.pages.dev
```

You can also access branch-specific previews:
```
<branch-name>.<project-name>.pages.dev
```

## Workers Integration

### Adding a Worker

Create `functions/` directory for Pages Functions (Workers):

```javascript
// functions/api/hello.js
export async function onRequest(context) {
  return new Response(JSON.stringify({
    message: "Hello from Cloudflare Workers!"
  }), {
    headers: { "Content-Type": "application/json" }
  });
}
```

Access at `/api/hello`.

### Adding Middleware

Create `functions/_middleware.js` for all routes:

```javascript
export async function onRequest(context) {
  // Add security headers
  const response = await context.next();
  response.headers.set("X-Custom-Header", "value");
  return response;
}
```

## Environment Variables

### Setting Variables

In the Cloudflare dashboard:

1. Go to your Pages project
2. Click **Settings** > **Environment variables**
3. Add variables for **Production** and **Preview** separately

| Variable | Production | Preview |
|----------|------------|---------|
| `MARKATA_GO_URL` | `https://example.com` | (empty) |

### Accessing in Workers

```javascript
// functions/api/config.js
export async function onRequest(context) {
  const url = context.env.MARKATA_GO_URL;
  return new Response(JSON.stringify({ url }));
}
```

## Troubleshooting

### Build Fails: Go Not Found

The build environment doesn't have Go pre-installed. Use the build command that installs Go:

```bash
curl -sSL https://go.dev/dl/go1.22.0.linux-amd64.tar.gz | tar -xzf - && ./go/bin/go install github.com/WaylonWalker/markata-go/cmd/markata-go@latest && ~/go/bin/markata-go build --clean
```

### Build Timeout

Free tier builds timeout after 20 minutes. Optimize by:

1. Caching Go installation (not currently supported)
2. Pre-building locally and deploying:
   ```bash
   markata-go build --clean
   wrangler pages deploy public
   ```

### 404 on Subpages

Ensure your `_redirects` file doesn't interfere. Cloudflare Pages serves `index.html` files automatically for clean URLs.

### Headers Not Applied

Verify `_headers` file is in your output directory:

```bash
markata-go build --clean
cat public/_headers
```

### Custom Domain Not Working

1. Ensure domain is added to Cloudflare (not just Pages)
2. Check DNS is proxied (orange cloud enabled)
3. Wait for SSL certificate provisioning

## Performance Features

### Automatic Optimizations

Cloudflare Pages automatically provides:

- Global CDN distribution
- Brotli compression
- HTTP/2 and HTTP/3
- Early Hints
- Smart caching

### Web Analytics

Enable Cloudflare Web Analytics (free):

1. Go to **Analytics** in your Cloudflare dashboard
2. Click **Web Analytics** > **Add a site**
3. Add the script to your template:

```html
<script defer src='https://static.cloudflareinsights.com/beacon.min.js' data-cf-beacon='{"token": "your-token"}'></script>
```

### Page Rules

Create Page Rules for advanced caching:

1. Go to **Rules** > **Page Rules**
2. Create rules for your domain

Example: Cache everything for static assets:
```
example.com/static/*
Cache Level: Cache Everything
Edge Cache TTL: 1 month
```

## Comparison with Cloudflare Workers Sites

| Feature | Pages | Workers Sites |
|---------|-------|---------------|
| Git integration | Yes | No |
| Preview deploys | Automatic | Manual |
| Build system | Built-in | External |
| Functions | Pages Functions | Workers |
| Best for | Static sites | Dynamic apps |

Pages is recommended for markata-go sites.

## Next Steps

- [Cloudflare Pages Docs](https://developers.cloudflare.com/pages/) - Official documentation
- [Configuration Guide](../configuration/) - Customize your markata-go site
- [Feeds Guide](../feeds/) - Set up RSS and Atom feeds
