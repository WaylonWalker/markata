---
title: "Deploy with Docker"
description: "Step-by-step guide to deploy markata-go sites using Docker containers"
date: 2026-01-24
published: true
tags:
  - documentation
  - deployment
  - docker
---


Docker provides a consistent, reproducible environment for building and serving your markata-go site. It's ideal for self-hosting, local development, and integration with container orchestration platforms.

## Prerequisites

- Docker installed ([Get Docker](https://docs.docker.com/get-docker/))
- Docker Compose (optional, for multi-container setups)
- Your markata-go site ready to build

## Official markata-go Images

markata-go publishes two official container images for different workflows:

- `ghcr.io/waylonwalker/markata-go:<version>`: Minimal runtime image (scratch) that runs the `markata-go` binary directly.
- `ghcr.io/waylonwalker/markata-go-builder:<version>`: Builder image with `/bin/sh`, `aws`, `openssl`, core utilities, `rsync`, image encoders (`avifenc`, `cwebp`), Pagefind (standalone binary), Tailwind CLI (standalone binary), bundled CDN asset cache, bundled Mermaid JS, and Chromium for mermaid rendering (via Go-native chromedp, no Node.js required).

### Builder Image Quick Start

```bash
docker run --rm \
  -v "$PWD":/site \
  -w /site \
  ghcr.io/waylonwalker/markata-go-builder:latest \
  sh -c 'markata-go build --clean'
```

### Builder Image Publish Script

```bash
docker run --rm \
  -v "$PWD":/site \
  -v /webroot:/webroot \
  -w /site \
  ghcr.io/waylonwalker/markata-go-builder:latest \
  sh -c 'set -eu; ts=$(date -u +%Y%m%d%H%M%S); markata-go build --clean; rsync -a --delete public/ /webroot/releases/$ts/; ln -sfn releases/$ts /webroot/current'
```

### Chromium Mermaid Rendering in Containers

If your site uses `mode = "chromium"` for pre-rendered Mermaid diagrams, set
`no_sandbox = true` in your config. The Chromium sandbox requires kernel
capabilities that Docker restricts by default:

```toml
[markata-go.mermaid]
mode = "chromium"

[markata-go.mermaid.chromium]
no_sandbox = true
```

The builder image includes Chromium. For custom images, see
[Chromium in Containers](../configuration.md#chromium-in-containers-docker-distrobox-podman)
in the configuration guide for installation options including the lightweight
`chrome-headless-shell` binary that requires no root access.

## Cost

| Setup | Cost | Best For |
|-------|------|----------|
| Local Docker | Free | Development, testing |
| Self-hosted server | $5-20/mo (VPS) | Personal sites, full control |
| Container platforms | Varies | Production, scaling |

## Method 1: Simple Nginx Container

The simplest approach: build locally and serve with nginx.

### Step 1: Build Your Site

```bash
markata-go build --clean
```

### Step 2: Create Dockerfile

Create `Dockerfile` in your project root:

```dockerfile
FROM nginx:alpine

# Copy built site to nginx
COPY public/ /usr/share/nginx/html/

# Custom nginx config for clean URLs
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

### Step 3: Create nginx Configuration

Create `nginx.conf`:

```nginx
server {
    listen 80;
    server_name localhost;
    root /usr/share/nginx/html;
    index index.html;

    # Enable gzip compression
    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml text/javascript;

    # MIME types for txt/md files
    location ~ \.(txt|md)$ {
        default_type text/plain;
        charset utf-8;
    }

    # Try exact file first (for /robots.txt), then index.html, then directory
    # This supports reversed redirects where canonical files are at /slug.txt
    location / {
        try_files $uri $uri/index.html $uri/ =404;
    }

    # Cache static assets
    location /static/ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # Don't cache HTML
    location ~* \.html$ {
        expires -1;
        add_header Cache-Control "no-store, no-cache, must-revalidate";
    }

    # Security headers
    add_header X-Frame-Options "DENY" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;
}
```

### Step 4: Build and Run

```bash
# Build the image
docker build -t my-site .

# Run the container
docker run -d -p 8080:80 --name my-site my-site

# View your site at http://localhost:8080
```

## Method 2: Multi-Stage Build

Build the site inside Docker for fully reproducible builds.

### Dockerfile with Multi-Stage Build

```dockerfile
# Stage 1: Build the site
FROM golang:1.22-alpine AS builder

# Install git (needed for go install)
RUN apk add --no-cache git

# Install markata-go
RUN go install github.com/WaylonWalker/markata-go/cmd/markata-go@latest

# Set working directory
WORKDIR /site

# Copy site source
COPY . .

# Build the site
ENV MARKATA_GO_URL=https://example.com
RUN markata-go build --clean

# Stage 2: Serve with nginx
FROM nginx:alpine

# Copy built site from builder stage
COPY --from=builder /site/public/ /usr/share/nginx/html/

# Custom nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

### Build and Run

```bash
# Build with your production URL
docker build \
  --build-arg MARKATA_GO_URL=https://example.com \
  -t my-site .

# Run
docker run -d -p 8080:80 --name my-site my-site
```

## Method 3: Docker Compose

For development with live reload or multi-service setups.

### docker-compose.yml

```yaml
version: '3.8'

services:
  site:
    build: .
    ports:
      - "8080:80"
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost/"]
      interval: 30s
      timeout: 10s
      retries: 3

  # Optional: Add Watchtower for automatic updates
  watchtower:
    image: containrrr/watchtower
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    command: --interval 3600 site
```

### Commands

```bash
# Start services
docker compose up -d

# View logs
docker compose logs -f site

# Rebuild after changes
docker compose build && docker compose up -d

# Stop services
docker compose down
```

## Method 4: Development with Live Reload

For local development with automatic rebuilds.

### docker-compose.dev.yml

```yaml
version: '3.8'

services:
  dev:
    image: golang:1.22-alpine
    working_dir: /site
    volumes:
      - .:/site
      - go-cache:/go
    ports:
      - "8080:8080"
    command: |
      sh -c '
        apk add --no-cache git
        go install github.com/WaylonWalker/markata-go/cmd/markata-go@latest
        markata-go serve
      '
    environment:
      - MARKATA_GO_URL=http://localhost:8080

volumes:
  go-cache:
```

### Run Development Server

```bash
docker compose -f docker-compose.dev.yml up
```

## Production Deployment

### Deploying to a VPS

1. **Build and Push to Registry**

```bash
# Build image
docker build -t your-registry/my-site:latest .

# Push to registry (Docker Hub, GitHub Container Registry, etc.)
docker push your-registry/my-site:latest
```

2. **Pull and Run on Server**

```bash
# On your server
docker pull your-registry/my-site:latest
docker run -d \
  --name my-site \
  -p 80:80 \
  --restart unless-stopped \
  your-registry/my-site:latest
```

### With HTTPS (Traefik)

Use Traefik as a reverse proxy with automatic HTTPS.

#### docker-compose.prod.yml

```yaml
version: '3.8'

services:
  traefik:
    image: traefik:v2.10
    command:
      - "--api.insecure=true"
      - "--providers.docker=true"
      - "--providers.docker.exposedbydefault=false"
      - "--entrypoints.web.address=:80"
      - "--entrypoints.websecure.address=:443"
      - "--certificatesresolvers.letsencrypt.acme.httpchallenge=true"
      - "--certificatesresolvers.letsencrypt.acme.httpchallenge.entrypoint=web"
      - "--certificatesresolvers.letsencrypt.acme.email=you@example.com"
      - "--certificatesresolvers.letsencrypt.acme.storage=/letsencrypt/acme.json"
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - letsencrypt:/letsencrypt
    restart: unless-stopped

  site:
    build: .
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.site.rule=Host(`example.com`)"
      - "traefik.http.routers.site.entrypoints=websecure"
      - "traefik.http.routers.site.tls.certresolver=letsencrypt"
      - "traefik.http.routers.site-http.rule=Host(`example.com`)"
      - "traefik.http.routers.site-http.entrypoints=web"
      - "traefik.http.routers.site-http.middlewares=redirect-to-https"
      - "traefik.http.middlewares.redirect-to-https.redirectscheme.scheme=https"
    restart: unless-stopped

volumes:
  letsencrypt:
```

### With HTTPS (Caddy)

Caddy provides automatic HTTPS with simpler configuration.

#### Caddyfile

```
example.com {
    root * /srv
    file_server
    encode gzip

    # MIME types for txt/md files
    @txtmd path *.txt *.md
    header @txtmd Content-Type "text/plain; charset=utf-8"

    header /static/* Cache-Control "public, max-age=31536000, immutable"
    header *.html Cache-Control "no-cache, must-revalidate"

    # Try exact file first (for /robots.txt), then index.html, then directory
    try_files {path} {path}/index.html {path}/
}
```

#### docker-compose.caddy.yml

```yaml
version: '3.8'

services:
  caddy:
    image: caddy:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile
      - ./public:/srv
      - caddy_data:/data
      - caddy_config:/config
    restart: unless-stopped

volumes:
  caddy_data:
  caddy_config:
```

## CI/CD Integration

### GitHub Actions

```yaml
name: Build and Deploy Docker

on:
  push:
    branches: [main]

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  build-and-push:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write

    steps:
      - uses: actions/checkout@v4

      - name: Log in to Container Registry
        uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Build and push Docker image
        uses: docker/build-push-action@v5
        with:
          context: .
          push: true
          tags: |
            ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:latest
            ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:${{ github.sha }}
          build-args: |
            MARKATA_GO_URL=https://example.com

  deploy:
    needs: build-and-push
    runs-on: ubuntu-latest
    steps:
      - name: Deploy to server
        uses: appleboy/ssh-action@v1.0.0
        with:
          host: ${{ secrets.SERVER_HOST }}
          username: ${{ secrets.SERVER_USER }}
          key: ${{ secrets.SERVER_SSH_KEY }}
          script: |
            docker pull ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:latest
            docker stop my-site || true
            docker rm my-site || true
            docker run -d \
              --name my-site \
              -p 80:80 \
              --restart unless-stopped \
              ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:latest
```

## Container Orchestration

### Kubernetes Deployment

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: markata-site
spec:
  replicas: 2
  selector:
    matchLabels:
      app: markata-site
  template:
    metadata:
      labels:
        app: markata-site
    spec:
      containers:
        - name: site
          image: your-registry/my-site:latest
          ports:
            - containerPort: 80
          resources:
            limits:
              memory: "128Mi"
              cpu: "100m"
          livenessProbe:
            httpGet:
              path: /
              port: 80
            initialDelaySeconds: 5
            periodSeconds: 10
---
apiVersion: v1
kind: Service
metadata:
  name: markata-site
spec:
  selector:
    app: markata-site
  ports:
    - port: 80
      targetPort: 80
  type: ClusterIP
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: markata-site
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
spec:
  tls:
    - hosts:
        - example.com
      secretName: markata-site-tls
  rules:
    - host: example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: markata-site
                port:
                  number: 80
```

## Troubleshooting

### Container Won't Start

Check logs:
```bash
docker logs my-site
```

Common issues:
- Port already in use: Change the host port (`-p 8081:80`)
- Permission denied: Check file permissions in copied directories

### 404 Errors on Subpages

Ensure nginx is configured for the reversed redirect structure:
```nginx
# MIME types for txt/md files
location ~ \.(txt|md)$ {
    default_type text/plain;
    charset utf-8;
}

# Try exact file first, then index.html, then directory
location / {
    try_files $uri $uri/index.html $uri/ =404;
}
```

This order ensures:
- `/robots.txt` serves the canonical file directly
- `/my-post/` serves `/my-post/index.html`
- Directory requests fall back correctly

### Large Image Size

Optimize your Dockerfile:
```dockerfile
# Use alpine-based images
FROM nginx:alpine

# Use .dockerignore to exclude unnecessary files
```

Create `.dockerignore`:
```
.git
node_modules
*.md
Dockerfile
docker-compose*.yml
```

### Build Fails in Multi-Stage

Ensure Go modules are available:
```dockerfile
# If using go modules
COPY go.mod go.sum ./
RUN go mod download
```

### CSS/JS Not Loading

Verify MARKATA_GO_URL matches your deployment:
```bash
docker build --build-arg MARKATA_GO_URL=https://example.com -t my-site .
```

## Performance Optimization

### Enable Brotli Compression

```nginx
# In nginx.conf
brotli on;
brotli_types text/plain text/css application/json application/javascript text/xml application/xml text/javascript;
```

Requires nginx with brotli module:
```dockerfile
FROM fholzer/nginx-brotli:latest
```

### Resource Limits

```yaml
# In docker-compose.yml
services:
  site:
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 128M
        reservations:
          cpus: '0.1'
          memory: 64M
```

### Health Checks

```dockerfile
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget -q --spider http://localhost/ || exit 1
```

## Security Best Practices

### Run as Non-Root User

```dockerfile
FROM nginx:alpine

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy files
COPY --chown=appuser:appgroup public/ /usr/share/nginx/html/

# Switch to non-root user (for processes that don't need root)
# Note: nginx master process needs root for port 80
```

### Read-Only Filesystem

```yaml
services:
  site:
    read_only: true
    tmpfs:
      - /var/cache/nginx
      - /var/run
```

### Security Scanning

```bash
# Scan image for vulnerabilities
docker scout cve my-site:latest

# Or use Trivy
trivy image my-site:latest
```

## Next Steps

- [Self-Hosting Guide](../self-hosting/) - More self-hosting options
- [Configuration Guide](../configuration/) - Customize your markata-go site
- [Themes Guide](../themes/) - Change your site's appearance
