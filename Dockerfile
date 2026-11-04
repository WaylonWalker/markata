# Multi-stage Dockerfile for markata-go
# Produces a minimal, secure container image for production use

# =============================================================================
# Build Stage
# =============================================================================
FROM golang:1.27.3-alpine AS builder

# Install ca-certificates for HTTPS support in final image
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy dependency files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build arguments for version injection
ARG VERSION=dev
ARG COMMIT=none
ARG BUILD_DATE=unknown

# Build static binary with version information
# CGO_ENABLED=0 ensures fully static binary
# -s -w strips debug information for smaller binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w \
        -X github.com/WaylonWalker/markata/cmd/markata-go/cmd.Version=${VERSION} \
        -X github.com/WaylonWalker/markata/cmd/markata-go/cmd.Commit=${COMMIT} \
        -X github.com/WaylonWalker/markata/cmd/markata-go/cmd.Date=${BUILD_DATE}" \
    -o markata-go \
    ./cmd/markata-go

# Verify binary was built
RUN ./markata-go version

# =============================================================================
# Builder Runtime Stage
# =============================================================================
FROM alpine:3.20 AS builder-runtime

# Labels for container registry metadata
LABEL org.opencontainers.image.title="markata-go-builder"
LABEL org.opencontainers.image.description="Builder image for markata-go with shell and publish tooling"
LABEL org.opencontainers.image.url="https://github.com/WaylonWalker/markata-go"
LABEL org.opencontainers.image.source="https://github.com/WaylonWalker/markata-go"
LABEL org.opencontainers.image.vendor="Waylon Walker"
LABEL org.opencontainers.image.licenses="MIT"

# Install shell, TLS roots, and build/publish tooling
# Note: nodejs/npm/@mermaid-js/mermaid-cli are NOT needed here because
# markata-go's "chromium" mode uses chromedp (Go-native CDP) directly,
# and pagefind is installed as a standalone binary below.
# Browser-compatible Mermaid builds use the separate browser-builder-runtime target.
RUN apk add --no-cache \
    aws-cli \
    ca-certificates \
    coreutils \
    findutils \
    gawk \
    git \
    libavif-apps \
    libwebp-tools \
    openssl \
    openssh-client \
    rsync \
    tzdata

# Select architecture-matching standalone assets for multi-platform images.
# Release tags publish both x86_64 and aarch64 Pagefind MUSL archives.
ARG TARGETARCH
ARG PAGEFIND_VERSION=v1.4.0
RUN set -eu; \
    case "${TARGETARCH}" in amd64) pagefind_arch=x86_64 ;; arm64) pagefind_arch=aarch64 ;; *) echo "Unsupported TARGETARCH: ${TARGETARCH}" >&2; exit 1 ;; esac; \
    artifact="pagefind-${PAGEFIND_VERSION}-${pagefind_arch}-unknown-linux-musl.tar.gz"; \
    url="https://github.com/Pagefind/pagefind/releases/download/${PAGEFIND_VERSION}/${artifact}"; \
    wget -q "${url}" -O /tmp/pagefind.tar.gz; \
    wget -q "${url}.sha256" -O /tmp/pagefind.tar.gz.sha256; \
    expected="$(cut -d' ' -f1 /tmp/pagefind.tar.gz.sha256)"; \
    printf '%s  %s\n' "${expected}" /tmp/pagefind.tar.gz | sha256sum -c -; \
    tar xzf /tmp/pagefind.tar.gz -C /usr/local/bin/ pagefind; \
    chmod +x /usr/local/bin/pagefind; \
    rm /tmp/pagefind.tar.gz /tmp/pagefind.tar.gz.sha256

# Tailwind releases architecture-specific MUSL executables and checksums.
ARG TAILWIND_VERSION=v4.2.1
RUN set -eu; \
    case "${TARGETARCH}" in amd64) tailwind_bin=tailwindcss-linux-x64-musl ;; arm64) tailwind_bin=tailwindcss-linux-arm64-musl ;; *) echo "Unsupported TARGETARCH: ${TARGETARCH}" >&2; exit 1 ;; esac; \
    wget -q "https://github.com/tailwindlabs/tailwindcss/releases/download/${TAILWIND_VERSION}/${tailwind_bin}" -O "/tmp/${tailwind_bin}"; \
    wget -q "https://github.com/tailwindlabs/tailwindcss/releases/download/${TAILWIND_VERSION}/sha256sums.txt" -O /tmp/sha256sums.txt; \
    cd /tmp; \
    grep "  ${tailwind_bin}$" sha256sums.txt | sha256sum -c -; \
    install -m 0755 "${tailwind_bin}" /usr/local/bin/tailwindcss; \
    rm "${tailwind_bin}" sha256sums.txt

# Copy the binary
COPY --from=builder /app/markata-go /usr/local/bin/markata-go

# Pre-seed the builder image with the CDN asset cache and MermaidJS source so
# markata-go builds can run without internet access at runtime.
ARG BUNDLED_ASSETS_CACHE_DIR=/usr/local/share/markata-go/assets-cache
ARG BUNDLED_MERMAID_DIR=/usr/local/share/markata-go/mermaid
ARG MERMAID_JS_VERSION=10.9.5
RUN mkdir -p /tmp/markata-assets \
    && cd /tmp/markata-assets \
    && /usr/local/bin/markata-go assets download \
    && mkdir -p "${BUNDLED_ASSETS_CACHE_DIR}" "${BUNDLED_MERMAID_DIR}" \
    && cp -a .markata/assets-cache/. "${BUNDLED_ASSETS_CACHE_DIR}/" \
    && wget -q "https://cdn.jsdelivr.net/npm/mermaid@${MERMAID_JS_VERSION}/dist/mermaid.min.js" \
        -O "${BUNDLED_MERMAID_DIR}/mermaid-v${MERMAID_JS_VERSION}.min.js" \
    && rm -rf /tmp/markata-assets

ENV MARKATA_GO_OFFLINE=true \
    MARKATA_GO_BUNDLED_ASSETS_CACHE_DIR=${BUNDLED_ASSETS_CACHE_DIR} \
    MARKATA_GO_BUNDLED_MERMAID_DIR=${BUNDLED_MERMAID_DIR}

# Set working directory for user content
WORKDIR /site

# Default command opens a shell for scripts
CMD ["sh"]

# Optional compatibility image for authors who explicitly choose Mermaid
# Chromium mode. Ordinary native builds use browser-free builder-runtime.
FROM builder-runtime AS browser-builder-runtime
RUN apk add --no-cache chromium

# =============================================================================
# Runtime Stage
# =============================================================================
# Use scratch for the smallest possible runtime image
FROM scratch

# Labels for container registry metadata
LABEL org.opencontainers.image.title="markata-go"
LABEL org.opencontainers.image.description="A plugin-driven static site generator written in Go"
LABEL org.opencontainers.image.url="https://github.com/WaylonWalker/markata-go"
LABEL org.opencontainers.image.source="https://github.com/WaylonWalker/markata-go"
LABEL org.opencontainers.image.vendor="Waylon Walker"
LABEL org.opencontainers.image.licenses="MIT"

# Copy CA certificates for HTTPS support
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Copy the binary
COPY --from=builder /app/markata-go /usr/local/bin/markata-go

# Set working directory for user content
WORKDIR /site

# Run as non-root user (numeric UID/GID)
USER 65532:65532

# Default entrypoint
ENTRYPOINT ["markata-go"]

# Default command shows help
CMD ["--help"]
