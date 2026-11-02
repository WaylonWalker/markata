---
title: "Self-Hosting Guide"
description: "Deploy markata-go sites on your own servers using Docker, systemd, or manual setups"
date: 2024-01-23
published: true
tags:
  - documentation
  - deployment
  - self-hosting
  - docker
  - systemd
---


This guide covers deploying markata-go sites on your own infrastructure. Whether you prefer Docker containers, systemd services, or manual setups, you'll find everything you need to self-host your static site.

## Overview

Self-hosting gives you complete control over your site's infrastructure. markata-go supports several deployment methods:

| Method | Best For | Complexity |
|--------|----------|------------|
| Docker Compose | Quick setup, portability | Low |
| systemd | Linux servers, integration | Medium |
| Manual | Custom setups, learning | High |

**Prerequisites:**
- A Linux server (VPS, dedicated, or home server)
- Domain name pointed to your server
- Basic command line knowledge

## Quick Start with Docker

The fastest way to self-host is using Docker Compose:

```bash
# Clone your site or create a new one
git clone https://github.com/you/your-site.git
cd your-site

# Copy the deployment files
cp -r deploy/docker ./docker-deploy
cd docker-deploy

# Configure environment
cp .env.example .env
# Edit .env with your domain: SITE_DOMAIN=example.com

# Start production stack
docker compose -f docker-compose.prod.yml up -d
```

Your site will be available at `https://your-domain.com` with automatic HTTPS.

See [Docker Deployment README](https://github.com/WaylonWalker/markata-go/tree/main/deploy/docker) for detailed configuration options.

## Docker Compose Configurations

### Development Mode

For local development with hot reload:

```bash
cd deploy/docker
docker compose -f docker-compose.dev.yml up
```

Features:
- Hot reload on file changes
- Local preview at http://localhost:8000
- Go module caching

### Production Mode

For production with Caddy reverse proxy:

```bash
cd deploy/docker

# Configure
cp .env.example .env
vim .env  # Set SITE_DOMAIN and SITE_URL

# Deploy
docker compose -f docker-compose.prod.yml up -d
```

Features:
- Automatic HTTPS via Let's Encrypt
- HTTP/2 and HTTP/3 support
- Security headers
- Gzip/Zstd compression
- Optimized caching

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SITE_DOMAIN` | Your domain name | `example.com` |
| `SITE_URL` | Full URL for feeds/links | `https://example.com` |
| `SITE_DIR` | Path to site content | `./../..` |
| `PORT` | Development port | `8000` |

## systemd Services

For traditional Linux server deployments, use systemd services.

### Quick Setup

```bash
cd deploy/systemd
sudo ./install.sh
```

The installer prompts for:
- Site directory (default: `/var/www/mysite`)
- Site URL (default: `https://example.com`)
- Service user (default: `markata`)

### Manual Setup

```bash
# Create service user
sudo useradd -r -s /bin/false -d /var/www/mysite markata

# Create directories
sudo mkdir -p /var/www/mysite/{public,.go}
sudo chown -R markata:markata /var/www/mysite

# Install markata-go
sudo -u markata bash -c '
    export GOPATH=/var/www/mysite/.go
    go install github.com/WaylonWalker/markata-go/cmd/markata-go@latest
'

# Install services
sudo cp deploy/systemd/markata-go.service /etc/systemd/system/
sudo cp deploy/systemd/markata-go-watch.service /etc/systemd/system/
# Edit paths in service files

sudo systemctl daemon-reload
```

### Available Services

**markata-go.service** - One-time build:
```bash
sudo systemctl enable markata-go
sudo systemctl start markata-go

# Rebuild site
sudo systemctl restart markata-go
```

**markata-go-watch.service** - Auto-rebuild on changes:
```bash
sudo systemctl enable markata-go-watch
sudo systemctl start markata-go-watch

# View logs
sudo journalctl -u markata-go-watch -f
```

See [systemd Deployment README](https://github.com/WaylonWalker/markata-go/tree/main/deploy/systemd) for detailed setup and troubleshooting.

## Web Server Configuration

### Using nginx

If using systemd with nginx:

```nginx
server {
    listen 80;
    listen [::]:80;
    server_name example.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name example.com;

    ssl_certificate /etc/letsencrypt/live/example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/example.com/privkey.pem;

    root /var/www/mysite/public;
    index index.html;

    # Security headers
    add_header X-Frame-Options "DENY" always;
    add_header X-Content-Type-Options "nosniff" always;

    # MIME types for txt/md files
    location ~ \.(txt|md)$ {
        default_type text/plain;
        charset utf-8;
    }

    # Caching
    location /static/ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # Try exact file first, then index.html, then directory
    # This supports reversed redirects where /robots.txt is canonical
    location / {
        try_files $uri $uri/index.html $uri/ =404;
    }
}
```

**Testing nginx configuration:**

```bash
# Test configuration syntax
sudo nginx -t

# Reload configuration
sudo systemctl reload nginx

# Verify txt files are served correctly
curl -I https://example.com/robots.txt
# Should return: Content-Type: text/plain; charset=utf-8
```

### Using Caddy

Simpler alternative with automatic HTTPS:

```caddyfile
example.com {
    root * /var/www/mysite/public
    file_server
    encode gzip

    header {
        X-Frame-Options "DENY"
        X-Content-Type-Options "nosniff"
    }

    # MIME types for txt/md files
    @txtmd path *.txt *.md
    header @txtmd Content-Type "text/plain; charset=utf-8"

    @static path /static/*
    header @static Cache-Control "public, max-age=31536000, immutable"

    # Try exact file first (for /robots.txt), then index.html, then directory
    try_files {path} {path}/index.html {path}/
}
```

**Testing Caddy configuration:**

```bash
# Validate Caddyfile syntax
caddy validate --config /etc/caddy/Caddyfile

# Reload configuration
sudo systemctl reload caddy

# Verify txt files are served correctly
curl -I https://example.com/robots.txt
# Should return: Content-Type: text/plain; charset=utf-8
```

## Live Update Patterns

There are several approaches to keeping your self-hosted site up-to-date with content changes.

### Development Mode (Watch)

For development, use the serve command which watches for file changes by default:

```bash
# Watch is enabled by default
markata-go serve --port 8000 --host 0.0.0.0

# Explicitly disable watch for static serving
markata-go serve --no-watch
```

**Note:** The `serve` command has file watching enabled by default. Use `--no-watch` to disable it.

### Git Webhooks

For production, trigger rebuilds via webhooks when content is pushed to your repository.

**1. Create a webhook handler script:**

```bash
#!/bin/bash
# /usr/local/bin/webhook-rebuild.sh

set -e

SITE_DIR="/var/www/mysite"
LOG_FILE="/var/log/markata-rebuild.log"
LOCK_FILE="/tmp/markata-rebuild.lock"

# Prevent concurrent rebuilds
if [ -f "$LOCK_FILE" ]; then
    echo "$(date): Rebuild already in progress, skipping" >> "$LOG_FILE"
    exit 0
fi

trap "rm -f $LOCK_FILE" EXIT
touch "$LOCK_FILE"

echo "$(date): Starting rebuild" >> "$LOG_FILE"

cd "$SITE_DIR"
git fetch origin main
git reset --hard origin/main

# Rebuild the site
/var/www/mysite/.go/bin/markata-go build --clean >> "$LOG_FILE" 2>&1

echo "$(date): Rebuild complete" >> "$LOG_FILE"
```

**2. Set up a lightweight webhook server:**

Create `/etc/systemd/system/webhook.service`:

```ini
[Unit]
Description=Webhook server for site rebuilds
After=network.target

[Service]
Type=simple
User=markata
ExecStart=/usr/bin/webhook -hooks /etc/webhook/hooks.json -port 9000
Restart=always

[Install]
WantedBy=multi-user.target
```

**3. Configure the webhook:**

Create `/etc/webhook/hooks.json`:

```json
[
  {
    "id": "rebuild-site",
    "execute-command": "/usr/local/bin/webhook-rebuild.sh",
    "command-working-directory": "/var/www/mysite",
    "pass-arguments-to-command": [],
    "trigger-rule": {
      "match": {
        "type": "payload-hmac-sha256",
        "secret": "your-webhook-secret",  # pragma: allowlist secret
        "parameter": {
          "source": "header",
          "name": "X-Hub-Signature-256"
        }
      }
    }
  }
]
```

**4. Configure your Git host:**

- GitHub: Settings > Webhooks > Add webhook
- URL: `https://your-domain.com/hooks/rebuild-site`
- Secret: Same as in hooks.json
- Events: Push events

### Scheduled Rebuilds

For sites that pull data from external sources, use scheduled rebuilds.

**1. Create the rebuild service:**

```ini
# /etc/systemd/system/markata-go-rebuild.service
[Unit]
Description=Rebuild markata-go site
After=network.target

[Service]
Type=oneshot
User=markata
WorkingDirectory=/var/www/mysite
Environment=HOME=/var/www/mysite
Environment=GOPATH=/var/www/mysite/.go
Environment=PATH=/var/www/mysite/.go/bin:/usr/local/go/bin:/usr/bin:/bin
ExecStart=/var/www/mysite/.go/bin/markata-go build --clean
StandardOutput=journal
StandardError=journal
```

**2. Create the timer:**

```ini
# /etc/systemd/system/markata-go-rebuild.timer
[Unit]
Description=Rebuild site on schedule

[Timer]
# Rebuild hourly
OnCalendar=hourly
# Or use specific times: OnCalendar=*-*-* 06,12,18:00:00
# Catch up on missed runs
Persistent=true
# Add randomized delay to avoid thundering herd
RandomizedDelaySec=300

[Install]
WantedBy=timers.target
```

**3. Enable the timer:**

```bash
sudo systemctl daemon-reload
sudo systemctl enable markata-go-rebuild.timer
sudo systemctl start markata-go-rebuild.timer

# Check timer status
sudo systemctl list-timers markata-go-rebuild.timer
```

### CI/CD Integration

Build in CI and deploy via rsync:

```yaml
# .github/workflows/deploy.yml
name: Deploy Site

on:
  push:
    branches: [main]
  schedule:
    # Rebuild daily at 6am UTC
    - cron: '0 6 * * *'
  workflow_dispatch:

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'

      - name: Build site
        run: |
          go install github.com/WaylonWalker/markata-go/cmd/markata-go@latest
          markata-go build --clean
        env:
          MARKATA_GO_URL: https://example.com

      - name: Deploy via rsync
        run: |
          rsync -avz --delete public/ user@server:/var/www/mysite/public/
        env:
          SSH_PRIVATE_KEY: ${{ secrets.SSH_PRIVATE_KEY }}
```

## Health Checks

Health checks help ensure your site is running correctly and enable automated recovery.

### HTTP Health Checks

**Basic curl check:**

```bash
#!/bin/bash
# /usr/local/bin/health-check.sh
if curl -sf http://localhost:8000/ > /dev/null; then
    echo "Site is healthy"
    exit 0
else
    echo "Site is unhealthy"
    exit 1
fi
```

**systemd health monitoring:**

Add to your service file:

```ini
[Service]
# ... existing config ...

# Health check with automatic restart
ExecStartPost=/bin/sleep 5
ExecStartPost=/usr/bin/curl -sf http://localhost:8000/ || exit 1

# Watchdog (systemd will restart if no response)
WatchdogSec=60
NotifyAccess=main
```

### Docker Health Checks

**docker-compose.yml health check:**

```yaml
services:
  markata:
    # ... other config ...
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:8000/"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 30s
```

**Multi-endpoint health check:**

```yaml
healthcheck:
  test: |
    wget --spider -q http://localhost:8000/ &&
    wget --spider -q http://localhost:8000/blog/ &&
    wget --spider -q http://localhost:8000/blog/rss.xml
  interval: 60s
  timeout: 15s
  retries: 3
```

### External Monitoring

Use external services to monitor your site from outside your infrastructure:

| Service | Free Tier | Features |
|---------|-----------|----------|
| [Uptime Robot](https://uptimerobot.com/) | 50 monitors | HTTP, keyword, ping |
| [Healthchecks.io](https://healthchecks.io/) | 20 checks | Cron monitoring, alerts |
| [Better Stack](https://betterstack.com/) | 10 monitors | Status pages, incidents |

**Healthchecks.io integration for scheduled rebuilds:**

```bash
#!/bin/bash
# /usr/local/bin/webhook-rebuild.sh

HEALTHCHECK_URL="https://hc-ping.com/your-uuid"

# Ping start
curl -fsS -m 10 --retry 5 "${HEALTHCHECK_URL}/start" > /dev/null

# Do the rebuild
cd /var/www/mysite
git pull origin main
/var/www/mysite/.go/bin/markata-go build --clean

# Ping success or failure
if [ $? -eq 0 ]; then
    curl -fsS -m 10 --retry 5 "${HEALTHCHECK_URL}" > /dev/null
else
    curl -fsS -m 10 --retry 5 "${HEALTHCHECK_URL}/fail" > /dev/null
fi
```

### Container Orchestration Health

**Kubernetes readiness probe:**

```yaml
apiVersion: v1
kind: Pod
spec:
  containers:
    - name: markata
      readinessProbe:
        httpGet:
          path: /
          port: 8000
        initialDelaySeconds: 10
        periodSeconds: 5
      livenessProbe:
        httpGet:
          path: /
          port: 8000
        initialDelaySeconds: 30
        periodSeconds: 10
```

**Docker Swarm health check:**

```yaml
services:
  markata:
    deploy:
      replicas: 2
      update_config:
        parallelism: 1
        delay: 10s
        failure_action: rollback
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:8000/"]
      interval: 30s
      timeout: 10s
      retries: 3
```

## Security Checklist

Before going live:

- [ ] Firewall configured (only ports 80, 443, and SSH)
- [ ] SSH key authentication only (disable password auth)
- [ ] Automatic security updates enabled
- [ ] HTTPS working with valid certificate
- [ ] Security headers configured
- [ ] Service running as non-root user
- [ ] File permissions restricted

### Firewall Setup (ufw)

```bash
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow http
sudo ufw allow https
sudo ufw enable
```

## Monitoring

### Basic Health Checks

```bash
# Check if site is responding
curl -I https://example.com

# Check systemd service
sudo systemctl status markata-go

# View recent logs
sudo journalctl -u markata-go -n 100
```

### Uptime Monitoring

Consider using:
- [Uptime Robot](https://uptimerobot.com/) (free tier available)
- [Healthchecks.io](https://healthchecks.io/) for cron monitoring
- Self-hosted: [Uptime Kuma](https://github.com/louislam/uptime-kuma)

## Troubleshooting

### Site Not Loading

1. Check if the service is running:
   ```bash
   sudo systemctl status markata-go
   docker compose ps
   ```

2. Check if web server is running:
   ```bash
   sudo systemctl status nginx  # or caddy
   ```

3. Check firewall:
   ```bash
   sudo ufw status
   ```

### HTTPS Not Working

1. Verify DNS is pointing to your server
2. Check Let's Encrypt logs:
   ```bash
   # Caddy
   docker compose logs caddy
   # Certbot
   sudo certbot certificates
   ```

### Build Failures

1. Check build logs:
   ```bash
   sudo journalctl -u markata-go
   docker compose logs builder
   ```

2. Run build manually:
   ```bash
   cd /var/www/mysite
   markata-go build --clean -v
   ```

## Resource Requirements

| Setup | CPU | RAM | Disk |
|-------|-----|-----|------|
| Minimum | 1 core | 512 MB | 1 GB |
| Recommended | 2 cores | 1 GB | 5 GB |
| With build cache | 2 cores | 2 GB | 10 GB |

## Next Steps

- [Configuration Guide](/docs/guides/configuration/) - Customize your site
- [Deployment Guide](/docs/guides/deployment/) - Other hosting options
- [Themes Guide](/docs/guides/themes/) - Customize appearance
