#!/bin/bash
set -e

# Repository configuration
REPO="devlopersabbir/vpcm"

# Text Styling & Colors
BOLD=$(tput bold 2>/dev/null || echo "")
NORMAL=$(tput sgr0 2>/dev/null || echo "")
RED=$(tput setaf 1 2>/dev/null || echo "")
GREEN=$(tput setaf 2 2>/dev/null || echo "")
YELLOW=$(tput setaf 3 2>/dev/null || echo "")
BLUE=$(tput setaf 4 2>/dev/null || echo "")
CYAN=$(tput setaf 6 2>/dev/null || echo "")
WHITE=$(tput setaf 7 2>/dev/null || echo "")

# Helper functions to print fancy messages
info() {
    echo -e "${BLUE}[info]${NORMAL} $1"
}
success() {
    echo -e "${GREEN}[✓]${NORMAL} $1"
}
warn() {
    echo -e "${YELLOW}[!]${NORMAL} $1"
}
error() {
    echo -e "${RED}[error]${NORMAL} $1"
}

# Print Banner
echo -e "${CYAN}${BOLD}"
echo "  __   _____  ____  __  __ "
echo "  \ \ / /  _ \/ ___||  \/  |"
echo "   \ V /| |_) \___ \| |\/| |"
echo "    \ / |  __/ ___) | |  | |"
echo "     V  |_|   |____/|_|  |_|"
echo -e "${NORMAL}"
echo -e "${WHITE}${BOLD}Welcome to the interactive VPSM & VPCM installer!${NORMAL}"
echo -e "This script will retrieve, unpack, and configure the latest release.\n"

# Detect OS and Arch
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64|amd64)
        ARCH="amd64"
        ;;
    arm64|aarch64)
        ARCH="arm64"
        ;;
    i386|i686)
        ARCH="386"
        ;;
    *)
        error "Unsupported architecture $ARCH"
        exit 1
        ;;
esac

case "$OS" in
    darwin)
        OS="darwin"
        ;;
    linux)
        OS="linux"
        ;;
    *)
        error "Unsupported operating system $OS"
        exit 1
        ;;
esac

# Get latest release tag from GitHub API
info "Checking latest release of VPSM..."
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_RELEASE" ]; then
    error "Could not retrieve latest release version."
    exit 1
fi

success "Latest release found: ${GREEN}${BOLD}$LATEST_RELEASE${NORMAL}"

# Asset URL mapping
FILENAME="vpsm-${OS}-${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/$FILENAME"

# Interactive Step 1: Install Directory
echo -e "\n${BOLD}Step 1: Choose Installation Directory${NORMAL}"
DEFAULT_INSTALL_DIR="/usr/local/bin"
read -p "Install path [default: $DEFAULT_INSTALL_DIR]: " INSTALL_DIR < /dev/tty
INSTALL_DIR=${INSTALL_DIR:-$DEFAULT_INSTALL_DIR}

# Expand tilde (~) if present
if [[ "$INSTALL_DIR" == ~* ]]; then
    INSTALL_DIR="${INSTALL_DIR/#\~/$HOME}"
fi

# Clean up path (remove trailing slashes, resolve relative path)
if [[ "$INSTALL_DIR" != /* ]]; then
    # Convert relative to absolute
    INSTALL_DIR="$(pwd)/$INSTALL_DIR"
fi

# Validate path format to prevent literal code/weird strings
if [[ "$INSTALL_DIR" =~ [\$\{\}\*\?\'\"\|\<\>\&\;] ]]; then
    error "Invalid installation path: $INSTALL_DIR contains illegal characters."
    exit 1
fi

# Interactive Step 2: Configure SSH Wrapper
echo -e "\n${BOLD}Step 2: Shell Wrapper Override Configuration${NORMAL}"
read -p "Do you want to enable the shell wrapper to intercept ssh commands automatically? [Y/n]: " ENABLE_WRAPPER < /dev/tty
ENABLE_WRAPPER=${ENABLE_WRAPPER:-"y"}

# Temporary directory for download
TMP_DIR=$(mktemp -d)
clean_up() {
    rm -rf "$TMP_DIR"
}
trap clean_up EXIT

# Download
info "Downloading $FILENAME from GitHub..."
curl -fsSL "$DOWNLOAD_URL" -o "$TMP_DIR/$FILENAME"

# Extract
info "Extracting binary..."
tar -xzf "$TMP_DIR/$FILENAME" -C "$TMP_DIR"

# Locate the binary (resolves the naming mismatch e.g. vpsm-darwin-arm64)
EXTRACTED_BIN=$(find "$TMP_DIR" -type f -name "vpsm*" | head -n 1)

if [ -z "$EXTRACTED_BIN" ]; then
    error "Failed to locate extracted binary."
    exit 1
fi

# Ensure install dir exists
if [ ! -d "$INSTALL_DIR" ]; then
    info "Creating installation directory $INSTALL_DIR..."
    if [ ! -w "$(dirname "$INSTALL_DIR")" ]; then
        sudo mkdir -p "$INSTALL_DIR"
    else
        mkdir -p "$INSTALL_DIR"
    fi
fi

# Install
if [ ! -w "$INSTALL_DIR" ]; then
    warn "Requesting administrator privileges to write to $INSTALL_DIR..."
    sudo mv "$EXTRACTED_BIN" "$INSTALL_DIR/vpsm"
else
    mv "$EXTRACTED_BIN" "$INSTALL_DIR/vpsm"
fi

# Ensure executable permissions
if [ ! -w "$INSTALL_DIR" ]; then
    sudo chmod +x "$INSTALL_DIR/vpsm"
else
    chmod +x "$INSTALL_DIR/vpsm"
fi

# Create symlink/wrapper for vpcm if it doesn't exist
if [ -L "$INSTALL_DIR/vpcm" ] || [ -f "$INSTALL_DIR/vpcm" ]; then
    if [ ! -w "$INSTALL_DIR" ]; then
        sudo rm -f "$INSTALL_DIR/vpcm"
    else
        rm -f "$INSTALL_DIR/vpcm"
    fi
fi

if [ ! -w "$INSTALL_DIR" ]; then
    sudo ln -s "$INSTALL_DIR/vpsm" "$INSTALL_DIR/vpcm"
else
    ln -s "$INSTALL_DIR/vpsm" "$INSTALL_DIR/vpcm"
fi

# Wrapper shell insertion
if [[ "$ENABLE_WRAPPER" =~ ^[Yy]$ ]]; then
    info "Configuring shell wrappers..."
    # We will try to download or read shell_wrapper.sh from the repository or local repo if available
    WRAPPER_CONTENT=$(curl -fsSL "https://raw.githubusercontent.com/$REPO/main/scripts/shell_wrapper.sh" || echo "")
    if [ -n "$WRAPPER_CONTENT" ]; then
        if [ -f "$HOME/.zshrc" ] && ! grep -q "VPSM ssh wrapper override" "$HOME/.zshrc"; then
            echo "$WRAPPER_CONTENT" >> "$HOME/.zshrc"
            success "Configured SSH wrapper in ~/.zshrc"
        fi
        if [ -f "$HOME/.bashrc" ] && ! grep -q "VPSM ssh wrapper override" "$HOME/.bashrc"; then
            echo "$WRAPPER_CONTENT" >> "$HOME/.bashrc"
            success "Configured SSH wrapper in ~/.bashrc"
        fi
    else
        warn "Could not fetch shell_wrapper.sh content from GitHub. Skipping wrapper config."
    fi
fi

echo -e "\n${GREEN}${BOLD}✨ VPSM (vpsm/vpcm) has been successfully installed to $INSTALL_DIR!${NORMAL}"
echo -e "Run ${CYAN}'vpsm version'${NORMAL} or ${CYAN}'vpcm version'${NORMAL} to verify your installation.\n"
