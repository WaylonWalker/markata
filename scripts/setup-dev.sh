#!/usr/bin/env bash
# Setup script for markata-go development environment
# Usage: ./scripts/setup-dev.sh

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print functions
info() { echo -e "${BLUE}ℹ${NC} $1"; }
success() { echo -e "${GREEN}✓${NC} $1"; }
warn() { echo -e "${YELLOW}⚠${NC} $1"; }
error() { echo -e "${RED}✗${NC} $1"; }

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  markata-go Development Environment Setup"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Check if we're in the right directory
if [[ ! -f "go.mod" ]] || ! grep -qE "^module github.com/WaylonWalker/markata$" go.mod 2>/dev/null; then
    error "This script must be run from the markata-go repository root"
    exit 1
fi

# ─────────────────────────────────────────────────────────────────────────────
# Check prerequisites
# ─────────────────────────────────────────────────────────────────────────────
info "Checking prerequisites..."

# Check Go
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    success "Go found: $GO_VERSION"
else
    error "Go is not installed. Please install Go 1.27+ from https://go.dev/dl/"
    exit 1
fi

# Check git
if command -v git &> /dev/null; then
    success "Git found: $(git --version | awk '{print $3}')"
else
    error "Git is not installed"
    exit 1
fi

echo ""

# ─────────────────────────────────────────────────────────────────────────────
# Install Go dependencies
# ─────────────────────────────────────────────────────────────────────────────
info "Downloading Go dependencies..."
go mod download
success "Go dependencies downloaded"

echo ""

# ─────────────────────────────────────────────────────────────────────────────
# Install development tools
# ─────────────────────────────────────────────────────────────────────────────
info "Installing development tools..."

# golangci-lint
if command -v golangci-lint &> /dev/null; then
    success "golangci-lint already installed"
else
    info "Installing golangci-lint..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    success "golangci-lint installed"
fi

# goimports
if command -v goimports &> /dev/null; then
    success "goimports already installed"
else
    info "Installing goimports..."
    go install golang.org/x/tools/cmd/goimports@latest
    success "goimports installed"
fi

# goreleaser (optional)
if command -v goreleaser &> /dev/null; then
    success "goreleaser already installed"
else
    info "Installing goreleaser..."
    go install github.com/goreleaser/goreleaser/v2@latest
    success "goreleaser installed"
fi

echo ""

# ─────────────────────────────────────────────────────────────────────────────
# Install pre-commit
# ─────────────────────────────────────────────────────────────────────────────
info "Setting up pre-commit hooks..."

if command -v pre-commit &> /dev/null; then
    success "pre-commit already installed"
else
    warn "pre-commit not found. Installing..."

    # Try pip first
    if command -v pip3 &> /dev/null; then
        pip3 install --user pre-commit
        success "pre-commit installed via pip3"
    elif command -v pip &> /dev/null; then
        pip install --user pre-commit
        success "pre-commit installed via pip"
    elif command -v brew &> /dev/null; then
        brew install pre-commit
        success "pre-commit installed via brew"
    else
        warn "Could not install pre-commit automatically."
        warn "Please install it manually: https://pre-commit.com/#installation"
        warn "  - pip install pre-commit"
        warn "  - brew install pre-commit"
        warn "  - apt install pre-commit"
    fi
fi

# Install the hooks
if command -v pre-commit &> /dev/null; then
    info "Installing git hooks..."
    pre-commit install
    pre-commit install --hook-type commit-msg
    success "Git hooks installed"

    # Create secrets baseline if it doesn't exist
    if [[ ! -f ".secrets.baseline" ]]; then
        info "Creating secrets baseline..."
        if command -v detect-secrets &> /dev/null; then
            detect-secrets scan > .secrets.baseline
            success "Secrets baseline created"
        else
            # Create an empty baseline
            echo '{"version": "1.5.0", "plugins_used": [], "results": {}}' > .secrets.baseline
            success "Empty secrets baseline created"
        fi
    fi
else
    warn "pre-commit not available, skipping hook installation"
fi

echo ""

# ─────────────────────────────────────────────────────────────────────────────
# Verify setup
# ─────────────────────────────────────────────────────────────────────────────
info "Verifying setup..."

# Try to build
if go build ./cmd/markata-go 2>/dev/null; then
    success "Build successful"
    rm -f markata-go
else
    warn "Build failed - there may be issues to fix"
fi

# Run go vet
if go vet ./... 2>/dev/null; then
    success "go vet passed"
else
    warn "go vet found issues"
fi

echo ""

# ─────────────────────────────────────────────────────────────────────────────
# Summary
# ─────────────────────────────────────────────────────────────────────────────
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  Setup Complete!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "You can now use the following commands:"
echo ""
echo "  ${GREEN}just build${NC}     - Build the binary"
echo "  ${GREEN}just test${NC}      - Run tests"
echo "  ${GREEN}just lint${NC}      - Run linter"
echo "  ${GREEN}just check${NC}     - Run all quality checks"
echo "  ${GREEN}just ci${NC}        - Run full CI checks"
echo ""
echo "Pre-commit hooks are installed and will run automatically"
echo "on git commit. To run them manually:"
echo ""
echo "  ${GREEN}pre-commit run${NC}              - Run on staged files"
echo "  ${GREEN}pre-commit run --all-files${NC} - Run on all files"
echo ""
echo "For VS Code users, open the workspace and install the"
echo "recommended extensions for the best development experience."
echo ""
