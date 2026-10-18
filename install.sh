#!/bin/bash
#
# Markata installer (supports legacy markata-go releases)
# https://github.com/WaylonWalker/markata
#
# Usage:
#   curl -sSL https://raw.githubusercontent.com/WaylonWalker/markata/main/install.sh | bash
#
# Environment variables:
#   MARKATA_INSTALL_DIR    - Custom install directory (legacy: MARKATA_GO_INSTALL_DIR)
#   MARKATA_VERSION        - Specific version (legacy: MARKATA_GO_VERSION)
#
# This script:
#   1. Detects OS (Linux, macOS, Windows/WSL) and architecture (amd64, arm64)
#   2. Downloads the latest release from GitHub
#   3. Installs to appropriate directory
#   4. Verifies the installation

set -e

# Configuration
REPO="WaylonWalker/markata"
BINARY_NAME="markata"
GITHUB_API="https://api.github.com/repos/${REPO}"
GITHUB_RELEASES="https://github.com/${REPO}/releases"

# Colors (disabled if not a terminal)
if [ -t 1 ]; then
    RED='\033[0;31m'
    GREEN='\033[0;32m'
    YELLOW='\033[0;33m'
    BLUE='\033[0;34m'
    NC='\033[0m' # No Color
else
    RED=''
    GREEN=''
    YELLOW=''
    BLUE=''
    NC=''
fi

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[OK]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
}

# Check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Detect OS
detect_os() {
    local os
    os="$(uname -s)"

    case "${os}" in
        Linux*)
            # Check for Android (Termux)
            if [ -n "${ANDROID_DATA}" ] || [ -n "${PREFIX}" ] && [ -d "${PREFIX}/bin" ]; then
                echo "android"
            else
                echo "linux"
            fi
            ;;
        Darwin*)
            echo "darwin"
            ;;
        CYGWIN*|MINGW*|MSYS*)
            echo "windows"
            ;;
        FreeBSD*)
            echo "freebsd"
            ;;
        *)
            log_error "Unsupported operating system: ${os}"
            exit 1
            ;;
    esac
}

# Detect architecture
detect_arch() {
    local arch
    arch="$(uname -m)"

    case "${arch}" in
        x86_64|amd64)
            echo "x86_64"
            ;;
        aarch64|arm64)
            echo "arm64"
            ;;
        armv7l|armv7)
            echo "armv7"
            ;;
        *)
            log_error "Unsupported architecture: ${arch}"
            exit 1
            ;;
    esac
}

# Get the latest version from GitHub API
get_latest_version() {
    local version

    if command_exists curl; then
        version=$(curl -sS "${GITHUB_API}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    elif command_exists wget; then
        version=$(wget -qO- "${GITHUB_API}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    else
        log_error "Neither curl nor wget found. Please install one of them."
        exit 1
    fi

    if [ -z "${version}" ]; then
        log_error "Failed to fetch latest version from GitHub API"
        exit 1
    fi

    echo "${version}"
}

# Determine the install directory
get_install_dir() {
    # Use custom install dir if specified
    if [ -n "${MARKATA_INSTALL_DIR:-${MARKATA_GO_INSTALL_DIR:-}}" ]; then
        echo "${MARKATA_INSTALL_DIR:-${MARKATA_GO_INSTALL_DIR}}"
        return
    fi

    # Check if running as root
    if [ "$(id -u)" -eq 0 ]; then
        echo "/usr/local/bin"
    else
        # Try /usr/local/bin first if writable
        if [ -w "/usr/local/bin" ]; then
            echo "/usr/local/bin"
        else
            # Fall back to user's local bin
            echo "${HOME}/.local/bin"
        fi
    fi
}

# Ensure the install directory exists
ensure_install_dir() {
    local dir="$1"

    if [ ! -d "${dir}" ]; then
        log_info "Creating install directory: ${dir}"
        mkdir -p "${dir}"
    fi

    if [ ! -w "${dir}" ]; then
        log_error "Install directory is not writable: ${dir}"
        log_error "Try running with sudo or set MARKATA_INSTALL_DIR (legacy: MARKATA_GO_INSTALL_DIR) to a writable path"
        exit 1
    fi
}

# Download the preferred Markata artifact, falling back to historical releases.
# Both executable names are installed so existing scripts keep working.
download_and_install() {
    local os="$1"
    local arch="$2"
    local version="$3"
    local install_dir="$4"
    local version_num="${version#v}"
    local ext="tar.gz"
    if [ "${os}" = "windows" ]; then
        ext="zip"
    fi

    local tmp_dir
    tmp_dir=$(mktemp -d)
    trap "rm -rf -- '${tmp_dir}'" EXIT

    local variant archive_name archive_path download_url=""
    for variant in markata markata-go; do
        archive_name="${variant}_${version_num}_${os}_${arch}.${ext}"
        archive_path="${tmp_dir}/${archive_name}"
        download_url="${GITHUB_RELEASES}/download/${version}/${archive_name}"
        log_info "Trying ${archive_name}..."
        if command_exists curl; then
            if curl -fsSL "${download_url}" -o "${archive_path}"; then
                break
            fi
        elif command_exists wget; then
            if wget -q "${download_url}" -O "${archive_path}"; then
                break
            fi
        else
            log_error "Neither curl nor wget found."
            exit 1
        fi
        archive_path=""
    done

    if [ -z "${archive_path}" ]; then
        log_error "Could not download either Markata or legacy markata-go archive for ${version} (${os}/${arch})."
        log_error "Browse available releases at ${GITHUB_RELEASES}"
        exit 1
    fi

    log_success "Downloaded ${archive_name}"
    if [ "${ext}" = "tar.gz" ]; then
        tar -xzf "${archive_path}" -C "${tmp_dir}"
    elif command_exists unzip; then
        unzip -q "${archive_path}" -d "${tmp_dir}"
    else
        log_error "unzip command not found."
        exit 1
    fi

    local suffix=""
    if [ "${os}" = "windows" ]; then
        suffix=".exe"
    fi
    local binary_path="${tmp_dir}/${variant}${suffix}"
    if [ ! -f "${binary_path}" ]; then
        log_error "Expected ${variant}${suffix} inside ${archive_name}."
        exit 1
    fi

    local canonical="${install_dir}/markata${suffix}"
    local legacy="${install_dir}/markata-go${suffix}"

    # The old Python distribution also owns a `markata` command. Never
    # silently overwrite it (or another unrelated executable).
    local target
    for target in "${canonical}" "${legacy}"; do
        if [ -e "${target}" ] || [ -L "${target}" ]; then
            if ! "${target}" version 2>/dev/null | grep -q '^  go:'; then
                log_error "Refusing to overwrite non-Go executable: ${target}"
                log_error "Move or uninstall the conflicting command before installing Markata."
                exit 1
            fi
        fi
    done

    # Stage and rename each command to avoid truncating a running executable.
    # A copy (rather than a symlink) also supports Windows-compatible shells.
    local canonical_tmp="${install_dir}/.markata-install-$${suffix}"
    local legacy_tmp="${install_dir}/.markata-go-install-$${suffix}"
    cp "${binary_path}" "${canonical_tmp}"
    chmod +x "${canonical_tmp}"
    cp "${canonical_tmp}" "${legacy_tmp}"
    chmod +x "${legacy_tmp}"
    mv -f "${canonical_tmp}" "${canonical}"
    mv -f "${legacy_tmp}" "${legacy}"
    log_success "Installed ${canonical} and compatibility command ${legacy}"
}

# Check if install dir is in PATH
check_path() {
    local install_dir="$1"

    case ":${PATH}:" in
        *":${install_dir}:"*)
            return 0
            ;;
    esac

    return 1
}

# Suggest PATH setup
suggest_path_setup() {
    local install_dir="$1"

    log_warn "${install_dir} is not in your PATH"
    echo ""
    echo "Add it to your PATH by adding this line to your shell config:"
    echo ""

    local shell_name
    shell_name=$(basename "${SHELL}")

    case "${shell_name}" in
        bash)
            echo "  echo 'export PATH=\"${install_dir}:\$PATH\"' >> ~/.bashrc"
            echo ""
            echo "Then reload your shell:"
            echo "  source ~/.bashrc"
            ;;
        zsh)
            echo "  echo 'export PATH=\"${install_dir}:\$PATH\"' >> ~/.zshrc"
            echo ""
            echo "Then reload your shell:"
            echo "  source ~/.zshrc"
            ;;
        fish)
            echo "  fish_add_path ${install_dir}"
            ;;
        *)
            echo "  export PATH=\"${install_dir}:\$PATH\""
            ;;
    esac
    echo ""
}

# Verify the installation
verify_installation() {
    local install_dir="$1"
    local binary_path="${install_dir}/${BINARY_NAME}"

    if [ ! -x "${binary_path}" ]; then
        log_error "Installation verification failed: binary not found or not executable"
        exit 1
    fi

    # Try to run the binary
    local version_output
    if version_output=$("${binary_path}" version --short 2>&1); then
        log_success "Verified installation: ${BINARY_NAME} ${version_output}"
    else
        log_warn "Binary installed but verification failed. You may need to add ${install_dir} to your PATH."
    fi
}

# Main installation function
main() {
    echo ""
    echo "Markata installer"
    echo "===================="
    echo ""

    # Detect platform
    local os arch
    os=$(detect_os)
    arch=$(detect_arch)

    log_info "Detected platform: ${os}/${arch}"

    # Check for Android-specific limitations
    if [ "${os}" = "android" ] && [ "${arch}" != "arm64" ]; then
        log_error "Only arm64 architecture is supported on Android"
        exit 1
    fi

    # Get version to install
    local version
    if [ -n "${MARKATA_VERSION:-${MARKATA_GO_VERSION:-}}" ]; then
        version="${MARKATA_VERSION:-${MARKATA_GO_VERSION}}"
        # Add 'v' prefix if missing
        if [[ "${version}" != v* ]]; then
            version="v${version}"
        fi
        log_info "Installing requested version: ${version}"
    else
        log_info "Fetching latest version..."
        version=$(get_latest_version)
        log_info "Latest version: ${version}"
    fi

    # Determine install directory
    local install_dir
    install_dir=$(get_install_dir)
    log_info "Install directory: ${install_dir}"

    # Ensure install directory exists and is writable
    ensure_install_dir "${install_dir}"

    # Download and install
    download_and_install "${os}" "${arch}" "${version}" "${install_dir}"

    # Check PATH
    if ! check_path "${install_dir}"; then
        suggest_path_setup "${install_dir}"
    fi

    # Verify installation
    verify_installation "${install_dir}"

    echo ""
    log_success "Installation complete!"
    echo ""
    echo "Get started:"
    echo "  ${BINARY_NAME} --help"
    echo "  ${BINARY_NAME} new \"My First Post\""
    echo "  ${BINARY_NAME} build"
    echo "  ${BINARY_NAME} serve"
    echo "  markata-go build   # Backward-compatible alias"
    echo ""
}

# Run main
main
