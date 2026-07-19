#!/bin/bash
set -e

# Repository configuration
REPO="devlopersabbir/vpcm"

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
        echo "Error: Unsupported architecture $ARCH"
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
        echo "Error: Unsupported operating system $OS"
        exit 1
        ;;
esac

# Get latest release tag from GitHub API
echo "Checking latest release of VPSM..."
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_RELEASE" ]; then
    echo "Error: Could not retrieve latest release version."
    exit 1
fi

echo "Latest release is $LATEST_RELEASE"

# Asset URL mapping
FILENAME="vpsm-${OS}-${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/$FILENAME"

# Temporary directory
TMP_DIR=$(mktemp -d)
clean_up() {
    rm -rf "$TMP_DIR"
}
trap clean_up EXIT

echo "Downloading $FILENAME from GitHub..."
curl -fsSL "$DOWNLOAD_URL" -o "$TMP_DIR/$FILENAME"

echo "Extracting binary..."
tar -xzf "$TMP_DIR/$FILENAME" -C "$TMP_DIR"

# Find install path
INSTALL_DIR="/usr/local/bin"
if [ ! -w "$INSTALL_DIR" ]; then
    echo "Requesting root privileges to install to $INSTALL_DIR..."
    sudo mv "$TMP_DIR/vpsm" "$INSTALL_DIR/vpsm"
else
    mv "$TMP_DIR/vpsm" "$INSTALL_DIR/vpsm"
fi

# Ensure executable permissions
chmod +x "$INSTALL_DIR/vpsm"

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

echo "✨ VPSM (vpsm/vpcm) has been successfully installed to $INSTALL_DIR!"
echo "Run 'vpsm version' or 'vpcm version' to verify installation."
