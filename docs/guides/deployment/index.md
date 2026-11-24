---
title: "Deployment Guide"
description: "Choose the best platform to deploy your markata-go static site"
date: 2026-01-24
published: true
tags:
  - documentation
  - deployment
  - hosting
---


markata-go generates static sites that can be deployed anywhere that serves HTML files. This guide helps you choose the best platform for your needs.

## Platform Comparison

| Platform | Free Tier | Custom Domain | Build Time | Best For |
|----------|-----------|---------------|------------|----------|
| [GitHub Pages](./github-pages/) | Unlimited | Yes (HTTPS) | ~2 min | Open source projects |
| [Netlify](./netlify/) | 100 GB/mo | Yes (HTTPS) | ~1 min | Teams, forms, functions |
| [Vercel](./vercel/) | 100 GB/mo | Yes (HTTPS) | ~1 min | Performance, edge |
| [Cloudflare Pages](./cloudflare-pages/) | Unlimited | Yes (HTTPS) | ~2 min | Global CDN, Workers |
| [AWS S3](./aws-s3/) | 5 GB (12 mo) | Yes (via CF) | Manual | Enterprise, AWS stack |
| [Docker](./docker/) | N/A | Self-managed | ~30 sec | Self-hosted, control |
| Kubernetes Helm | N/A | Yes | Cluster-dependent | Homelabs, ArgoCD, in-cluster builds |

## Quick Decision Guide

**Choose GitHub Pages if:**
- Your code is on GitHub
- You want zero configuration
- You don't need server-side features

**Choose Netlify if:**
- You need forms, functions, or split testing
- You want automatic deploy previews
- You're working with a team

**Choose Vercel if:**
- Performance is your top priority
- You want edge functions
- You need analytics

**Choose Cloudflare Pages if:**
- You want unlimited bandwidth
- You need Workers integration
- Global performance matters

**Choose AWS S3 if:**
- You're already in the AWS ecosystem
- You need fine-grained access control
- Enterprise compliance is required

**Choose Docker if:**
- You want full control
- You're self-hosting other services
- You need offline capabilities

**Choose Kubernetes Helm if:**
- You already deploy with ArgoCD or Helm
- You want in-cluster builds from a source archive
- You want same-origin `/api/search` served with the site

## Building for Production

Before deploying to any platform, build your site:

```bash
# Standard production build
markata-go build --clean

# Build with custom URL
MARKATA_GO_URL=https://example.com markata-go build --clean

# Verbose build for debugging
markata-go build --clean -v
```

### Output Structure

After building, your `public/` directory contains:

```
public/
├── index.html              # Home page
├── blog/
│   ├── index.html          # Blog listing
│   ├── rss.xml             # RSS feed
│   └── atom.xml            # Atom feed
├── posts/
│   └── my-post/
│       └── index.html      # Individual post
├── static/                 # Assets (CSS, JS, images)
└── sitemap.xml
```

## Environment Variables

All platforms support these environment variables:

| Variable | Purpose | Example |
|----------|---------|---------|
| `MARKATA_GO_URL` | Base URL for links and feeds | `https://example.com` |
| `MARKATA_GO_OUTPUT` | Output directory | `public` |

## Next Steps

Choose a platform guide to get started:

- [GitHub Pages](./github-pages/) - Free hosting for GitHub repos
- [Netlify](./netlify/) - Feature-rich with generous free tier
- [Vercel](./vercel/) - Performance-focused deployment
- [Cloudflare Pages](./cloudflare-pages/) - Unlimited bandwidth
- [AWS S3](./aws-s3/) - Enterprise-grade static hosting
- [Docker](./docker/) - Self-hosted containers

For Kubernetes and source-archive deployments, see the chart README under `helm-chart/README.md` in the markata-go repository.

For self-hosting on your own servers, see the [Self-Hosting Guide](../self-hosting/).
