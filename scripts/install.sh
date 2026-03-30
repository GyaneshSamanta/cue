#!/bin/bash
# cue installer for Linux and macOS
# Usage: curl -fsSL https://raw.githubusercontent.com/GyaneshSamanta/cue/main/scripts/install.sh | bash

set -euo pipefail

REPO="GyaneshSamanta/cue"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="$HOME/.cue"
BINARY_NAME="cue"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

print_step()  { echo -e "${CYAN}▸ $1${NC}"; }
print_ok()    { echo -e "${GREEN}✔ $1${NC}"; }
print_warn()  { echo -e "${YELLOW}⚠ $1${NC}"; }
print_err()   { echo -e "${RED}✖ $1${NC}"; }

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case "$OS" in
        linux)  PLATFORM="linux" ;;
        darwin) PLATFORM="darwin" ;;
        *)      print_err "Unsupported OS: $OS"; exit 1 ;;
    esac

    case "$ARCH" in
        x86_64|amd64) ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
        *)             print_err "Unsupported architecture: $ARCH"; exit 1 ;;
    esac

    print_ok "Detected: ${PLATFORM}/${ARCH}"
}

# Get latest release version
get_latest_version() {
    VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"v?([^"]+)".*/\1/')
    if [ -z "$VERSION" ]; then
        VERSION="1.0.0"
        print_warn "Could not fetch latest version, using v${VERSION}"
    fi
    print_ok "Latest version: v${VERSION}"
}

# Download binary
download() {
    FILENAME="${BINARY_NAME}-${PLATFORM}-${ARCH}"
    URL="https://github.com/${REPO}/releases/download/v${VERSION}/${FILENAME}"

    print_step "Downloading ${BINARY_NAME} v${VERSION}..."
    TMPDIR=$(mktemp -d)
    curl -fsSL -o "${TMPDIR}/${BINARY_NAME}" "$URL" || {
        print_err "Download failed. URL: ${URL}"
        print_warn "Try manual download from: https://github.com/${REPO}/releases"
        rm -rf "$TMPDIR"
        exit 1
    }
    chmod +x "${TMPDIR}/${BINARY_NAME}"
    print_ok "Downloaded successfully"
}

# Install binary
install_binary() {
    print_step "Installing to ${INSTALL_DIR}..."

    if [ -w "$INSTALL_DIR" ]; then
        mv "${TMPDIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
    else
        print_warn "Need sudo to install to ${INSTALL_DIR}"
        sudo mv "${TMPDIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
    fi

    rm -rf "$TMPDIR"
    print_ok "Installed to ${INSTALL_DIR}/${BINARY_NAME}"
}

# Create config directory
setup_config() {
    if [ ! -d "$CONFIG_DIR" ]; then
        mkdir -p "$CONFIG_DIR"
        print_ok "Created config directory: ${CONFIG_DIR}"
    fi

    # Write default config if none exists
    if [ ! -f "${CONFIG_DIR}/config.toml" ]; then
        cat > "${CONFIG_DIR}/config.toml" << 'EOF'
[core]
lock_poll_interval_secs = 5
lock_timeout_mins = 30
adaptive_backoff = true
notify_on_completion = true

[network]
probe_host = "1.1.1.1"
probe_fallback_host = "8.8.8.8"
probe_fallback_port = 53
fail_threshold = 3
recovery_threshold = 1
probe_interval_secs = 10

[history]
max_entries = 50000
default_display_count = 20

[workspace]
github_repo_name = "dev-workspace-backup"
backup_shell_configs = true
backup_vscode = false
backup_history = false

[ui]
color = true
progress_style = "bar"
explain_after_macro = true
EOF
        print_ok "Created default config"
    fi
}

# Verify installation
verify() {
    if command -v "$BINARY_NAME" &> /dev/null; then
        INSTALLED_VERSION=$("$BINARY_NAME" --version 2>/dev/null || echo "unknown")
        print_ok "Installation verified: ${INSTALLED_VERSION}"
    else
        print_warn "${BINARY_NAME} installed but not in PATH. Add ${INSTALL_DIR} to your PATH."
    fi
}

# Main
echo ""
echo "  ╔══════════════════════════════════════╗"
echo "  ║       cue installer         ║"
echo "  ║   Cross-Platform CLI Dev Utility     ║"
echo "  ╚══════════════════════════════════════╝"
echo ""

detect_platform
get_latest_version
download
install_binary
setup_config
verify

echo ""
print_ok "Installation complete! Run 'cue --help' to get started."
echo ""
