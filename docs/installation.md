---
title: "Installation"
description: "How to install Markata and preserve markata-go compatibility"
date: 2024-01-15
published: true
tags:
  - getting-started
  - installation
---

# Installation

**Markata** is the new command name for the Go static site generator. The installer installs both `markata` and the backward-compatible `markata-go` alias. Existing configuration files and environment variables continue to work; see [the name migration guide](guides/name-migration.md).

markata-go is distributed as a single binary with no dependencies. Choose the installation method that works best for you.

## Quick Install (Recommended)

### One-liner Install Script

The easiest way to install markata-go on Linux and macOS:

```bash
curl -sSL https://raw.githubusercontent.com/WaylonWalker/markata-go/main/install.sh | bash
```

This script:
- Automatically detects your OS and architecture
- Downloads the correct binary from GitHub releases
- Installs to `/usr/local/bin` or `~/.local/bin`
- Verifies the installation

**Custom installation directory:**

```bash
curl -sSL https://raw.githubusercontent.com/WaylonWalker/markata-go/main/install.sh | MARKATA_INSTALL_DIR=/opt/bin bash
```

**Install a specific version:**

```bash
curl -sSL https://raw.githubusercontent.com/WaylonWalker/markata-go/main/install.sh | MARKATA_VERSION=v0.1.0 bash
```

### Using jpillora/installer

The fastest way to install on Linux and macOS:

```bash
curl -sL https://i.jpillora.com/WaylonWalker/markata-go | bash
```

This automatically detects your OS and architecture and installs the latest release.

### Using eget

[eget](https://github.com/zyedidia/eget) is a convenient tool for downloading pre-built binaries:

```bash
# Install eget first (if you haven't)
go install github.com/zyedidia/eget@latest

# Install markata-go
eget WaylonWalker/markata-go
```

### Using mise

[mise](https://mise.jdx.dev/) is a polyglot runtime manager that can install tools directly from GitHub releases:

```bash
# Install globally
mise use -g github:WaylonWalker/markata-go

# Or add to your project's mise.toml
mise use github:WaylonWalker/markata-go
```

This uses mise's [GitHub backend](https://mise.jdx.dev/dev-tools/backends/github.html) which automatically detects the correct binary for your platform from GitHub releases.

## Package Managers

### Homebrew (macOS/Linux)

Coming soon! For now, use one of the methods above.

```bash
# Future support planned:
# brew install example/tap/markata-go
```

### Go Install

The `markata` entry point is available from the current development branch until it appears in a tagged release:

```bash
go install github.com/WaylonWalker/markata-go/cmd/markata@main
```

Previously tagged versions support `go install github.com/WaylonWalker/markata-go/cmd/markata-go@latest`.

Note: This installs to `$GOPATH/bin` (usually `~/go/bin`). Ensure this is in your `PATH`.

## Manual Download

Download the appropriate archive for your platform from the [GitHub Releases](https://github.com/WaylonWalker/markata-go/releases) page.

### Available Platforms

| OS | Architecture | Archive |
|---|---|---|
| Linux | x86_64 (amd64) | `markata-go_*_linux_x86_64.tar.gz` |
| Linux | ARM64 | `markata-go_*_linux_arm64.tar.gz` |
| Linux | ARMv7 | `markata-go_*_linux_armv7.tar.gz` |
| macOS | x86_64 (Intel) | `markata-go_*_darwin_x86_64.tar.gz` |
| macOS | ARM64 (Apple Silicon) | `markata-go_*_darwin_arm64.tar.gz` |
| Windows | x86_64 | `markata-go_*_windows_x86_64.zip` |
| FreeBSD | x86_64 | `markata-go_*_freebsd_x86_64.tar.gz` |
| Android | ARM64 (Termux) | `markata-go_*_android_arm64.tar.gz` |

### Manual Installation Steps

```bash
# Download (replace VERSION and PLATFORM)
curl -LO https://github.com/WaylonWalker/markata-go/releases/download/v0.1.0/markata-go_0.1.0_linux_x86_64.tar.gz

# Verify checksum (recommended)
curl -LO https://github.com/WaylonWalker/markata-go/releases/download/v0.1.0/checksums.txt
sha256sum -c checksums.txt --ignore-missing

# Extract
tar xzf markata-go_*.tar.gz

# Move to PATH
sudo mv markata-go /usr/local/bin/

# Verify installation
markata-go version
```

### Windows

1. Download the `.zip` file from releases
2. Extract the archive
3. Move `markata-go.exe` to a directory in your PATH
4. Or add the extraction directory to your PATH

## Building from Source

Requirements:
- Go 1.22 or later
- Git

```bash
# Clone the repository
git clone https://github.com/WaylonWalker/markata-go.git
cd markata-go

# Build with version info
go build -ldflags="-s -w" -o markata-go ./cmd/markata-go

# Or use just (if installed)
just build

# Install to GOPATH/bin
just install
```

### Development Build

For development with full version information:

```bash
# Using just
just build

# Or manually with ldflags
VERSION=$(git describe --tags --always --dirty)
COMMIT=$(git rev-parse --short HEAD)
DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

go build -ldflags "-s -w \
  -X github.com/WaylonWalker/markata-go/cmd/markata-go/cmd.Version=$VERSION \
  -X github.com/WaylonWalker/markata-go/cmd/markata-go/cmd.Commit=$COMMIT \
  -X github.com/WaylonWalker/markata-go/cmd/markata-go/cmd.Date=$DATE" \
  -o markata-go ./cmd/markata-go
```

## Verifying Installation

After installation, verify it works:

```bash
# Check version
markata-go version

# Show help
markata-go --help
```

Expected output:

```
markata-go 0.1.0
  commit:    abc1234
  built:     2024-01-15T12:00:00Z
  go:        go1.22.2
  os/arch:   linux/amd64
```

## Updating

### Quick Update

```bash
# Using jpillora/installer
curl -sL https://i.jpillora.com/WaylonWalker/markata-go | bash

# Using eget
eget WaylonWalker/markata-go

# Using mise
mise upgrade github:WaylonWalker/markata-go
```

### Go Install Update

```bash
go install github.com/WaylonWalker/markata-go/cmd/markata-go@latest
```

## Uninstalling

Remove the binary from your system:

```bash
# If installed to /usr/local/bin
sudo rm /usr/local/bin/markata-go

# If installed via go install
rm $(go env GOPATH)/bin/markata-go

# If using mise
mise uninstall github:WaylonWalker/markata-go
```

## Troubleshooting

### "command not found"

Ensure the installation directory is in your PATH:

```bash
# Check if markata-go is in PATH
which markata-go

# Add to PATH (bash)
echo 'export PATH="$PATH:/usr/local/bin"' >> ~/.bashrc
source ~/.bashrc

# Add Go bin to PATH
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bashrc
source ~/.bashrc
```

### Permission Denied

```bash
# Make binary executable
chmod +x markata-go

# Or use sudo for system-wide install
sudo mv markata-go /usr/local/bin/
```

### macOS Security Warning

If macOS blocks the binary:

1. Go to System Preferences > Security & Privacy
2. Click "Allow Anyway" for markata-go
3. Or remove the quarantine attribute:
   ```bash
   xattr -d com.apple.quarantine markata-go
   ```

## Next Steps

- [[quickstart|Quick Start Guide]] - Build your first site
- [[configuration-guide|Configuration]] - Configure your site
- [[cli-reference|CLI Reference]] - All available commands
