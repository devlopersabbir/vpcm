#!/bin/bash
set -e

# Repository configuration
REPO="devlopersabbir/vpcm"

# Text Styling & Colors (using printf-friendly variables)
BOLD=$(tput bold 2>/dev/null || echo "")
NORMAL=$(tput sgr0 2>/dev/null || echo "")
RED=$(tput setaf 1 2>/dev/null || echo "")
GREEN=$(tput setaf 2 2>/dev/null || echo "")
YELLOW=$(tput setaf 3 2>/dev/null || echo "")
BLUE=$(tput setaf 4 2>/dev/null || echo "")
CYAN=$(tput setaf 6 2>/dev/null || echo "")
WHITE=$(tput setaf 7 2>/dev/null || echo "")

# Helper functions to print messages portably
info() {
    printf "%s[info]%s %s\n" "${BLUE}" "${NORMAL}" "$1"
}
success() {
    printf "%s[✓]%s %s\n" "${GREEN}" "${NORMAL}" "$1"
}
warn() {
    printf "%s[!]%s %s\n" "${YELLOW}" "${NORMAL}" "$1"
}
error() {
    printf "%s[error]%s %s\n" "${RED}" "${NORMAL}" "$1"
}

# Print Banner
printf "%s%s\n" "${CYAN}" "${BOLD}"
echo "  __   _____  ____  __  __ "
echo "  \ \ / /  _ \/ ___||  \/  |"
echo "   \ V /| |_) \___ \| |\/| |"
echo "    \ / |  __/ ___) | |  | |"
echo "     V  |_|   |____/|_|  |_|"
printf "%s\n" "${NORMAL}"
printf "%s%sWelcome to the interactive VPSM & VPCM installer!%s\n" "${WHITE}" "${BOLD}" "${NORMAL}"
printf "This script will retrieve, unpack, and configure the latest release.\n\n"

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

# Default variables
INSTALL_DIR="/usr/local/bin"
ENABLE_WRAPPER="y"
NON_INTERACTIVE="n"

# Parse CLI arguments
while [[ $# -gt 0 ]]; do
    case "$1" in
        -d|--dir)
            INSTALL_DIR="$2"
            shift 2
            ;;
        -y|--yes|--non-interactive)
            NON_INTERACTIVE="y"
            shift
            ;;
        --no-wrapper)
            ENABLE_WRAPPER="n"
            shift
            ;;
        -h|--help)
            echo "Usage: install.sh [options]"
            echo "Options:"
            echo "  -d, --dir <path>     Installation directory (default: /usr/local/bin)"
            echo "  -y, --yes            Non-interactive installation (accept all defaults)"
            echo "  --no-wrapper         Do not configure the shell wrapper override"
            echo "  -h, --help           Show this help message"
            exit 0
            ;;
        *)
            shift
            ;;
    esac
done

# Check curl/wget availability
download_file() {
    local url="$1"
    local dest="$2"
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "$url" -o "$dest"
    elif command -v wget >/dev/null 2>&1; then
        wget -qO "$dest" "$url"
    else
        error "Neither curl nor wget is installed. Please install one of them to proceed."
        exit 1
    fi
}

get_latest_release() {
    # Method 1: Try using curl/wget to fetch the effective redirect URL from github.com/releases/latest
    # This completely bypasses GitHub API rate limit problems.
    local redirect_url=""
    if command -v curl >/dev/null 2>&1; then
        redirect_url=$(curl -Ls -o /dev/null -w "%{url_effective}" "https://github.com/$REPO/releases/latest" 2>/dev/null || echo "")
    elif command -v wget >/dev/null 2>&1; then
        redirect_url=$(wget --max-redirect=0 "https://github.com/$REPO/releases/latest" 2>&1 | grep -i "Location:" | awk '{print $2}' || echo "")
    fi

    if [ -n "$redirect_url" ] && [[ "$redirect_url" == *"/tag/"* ]]; then
        echo "${redirect_url##*/}"
        return 0
    fi

    # Method 2: Fall back to GitHub API
    local api_res=""
    if command -v curl >/dev/null 2>&1; then
        api_res=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" 2>/dev/null || echo "")
    elif command -v wget >/dev/null 2>&1; then
        api_res=$(wget -qO- "https://api.github.com/repos/$REPO/releases/latest" 2>/dev/null || echo "")
    fi

    if [ -n "$api_res" ]; then
        local tag=""
        tag=$(echo "$api_res" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' || echo "")
        if [ -n "$tag" ]; then
            echo "$tag"
            return 0
        fi
    fi

    return 1
}

# Resolve latest release
info "Checking latest release of VPSM..."
LATEST_RELEASE=$(get_latest_release || echo "")

if [ -z "$LATEST_RELEASE" ]; then
    error "Could not retrieve latest release version. Check network/rate limits."
    exit 1
fi

success "Latest release found: ${GREEN}${BOLD}$LATEST_RELEASE${NORMAL}"

# Asset URL mapping
FILENAME="vpsm-${OS}-${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/$FILENAME"

# Interactive installation path step
if [ "$NON_INTERACTIVE" != "y" ] && [ -t 0 ] && [ -c /dev/tty ]; then
    printf "\n%sStep 1: Choose Installation Directory%s\n" "${BOLD}" "${NORMAL}"
    read -p "Install path [default: $INSTALL_DIR]: " USER_INSTALL_DIR < /dev/tty
    INSTALL_DIR=${USER_INSTALL_DIR:-$INSTALL_DIR}
else
    info "Non-interactive installation: using directory $INSTALL_DIR"
fi

# Expand tilde (~) if present
if [[ "$INSTALL_DIR" == ~* ]]; then
    INSTALL_DIR="${INSTALL_DIR/#\~/$HOME}"
fi

# Clean up path (remove trailing slashes, resolve relative path)
if [[ "$INSTALL_DIR" != /* ]]; then
    INSTALL_DIR="$(pwd)/$INSTALL_DIR"
fi

# Validate path format to prevent literal code/weird strings
if [[ "$INSTALL_DIR" =~ [\$\{\}\*\?\'\"\|\<\>\&\;] ]]; then
    error "Invalid installation path: $INSTALL_DIR contains illegal characters."
    exit 1
fi

# Temporary directory for download
TMP_DIR=$(mktemp -d)
clean_up() {
    rm -rf "$TMP_DIR"
}
trap clean_up EXIT

# Download
info "Downloading $FILENAME from GitHub..."
download_file "$DOWNLOAD_URL" "$TMP_DIR/$FILENAME"

# Extract
info "Extracting binaries..."
tar -xzf "$TMP_DIR/$FILENAME" -C "$TMP_DIR"

# Ensure install dir exists
if [ ! -d "$INSTALL_DIR" ]; then
    info "Creating installation directory $INSTALL_DIR..."
    local parent_dir
    parent_dir=$(dirname "$INSTALL_DIR")
    if [ ! -w "$parent_dir" ]; then
        if command -v sudo >/dev/null 2>&1; then
            sudo mkdir -p "$INSTALL_DIR"
        else
            error "Cannot create $INSTALL_DIR (permission denied) and 'sudo' is not available."
            exit 1
        fi
    else
        mkdir -p "$INSTALL_DIR"
    fi
fi

# Function to move and chmod a binary safely
install_binary() {
    local BIN_NAME="$1"
    local SRC="$TMP_DIR/$BIN_NAME"
    local DEST="$INSTALL_DIR/$BIN_NAME"
    
    if [ ! -f "$SRC" ]; then
        SRC=$(find "$TMP_DIR" -type f -name "$BIN_NAME*" | head -n 1)
    fi

    if [ -z "$SRC" ] || [ ! -f "$SRC" ]; then
        error "Failed to locate binary $BIN_NAME in the release archive."
        exit 1
    fi

    # Check if target file or directory is writable
    if [ ! -w "$INSTALL_DIR" ] || { [ -f "$DEST" ] && [ ! -w "$DEST" ]; }; then
        if command -v sudo >/dev/null 2>&1; then
            sudo mv "$SRC" "$DEST"
            sudo chmod +x "$DEST"
        else
            error "Permission denied writing to $DEST and 'sudo' is not available."
            exit 1
        fi
    else
        mv "$SRC" "$DEST"
        chmod +x "$DEST"
    fi
    success "Installed $BIN_NAME to $INSTALL_DIR"
}

# Install all three core binaries
install_binary "vpsm"
install_binary "vpsmd"
install_binary "vpsm-api"

# Create symlink/wrapper for vpcm if it doesn't exist
if [ -L "$INSTALL_DIR/vpcm" ] || [ -f "$INSTALL_DIR/vpcm" ]; then
    if [ ! -w "$INSTALL_DIR" ] || { [ -f "$INSTALL_DIR/vpcm" ] && [ ! -w "$INSTALL_DIR/vpcm" ]; }; then
        if command -v sudo >/dev/null 2>&1; then
            sudo rm -f "$INSTALL_DIR/vpcm"
        else
            error "Permission denied removing old vpcm link and 'sudo' is not available."
            exit 1
        fi
    else
        rm -f "$INSTALL_DIR/vpcm"
    fi
fi

if [ ! -w "$INSTALL_DIR" ]; then
    if command -v sudo >/dev/null 2>&1; then
        sudo ln -s "$INSTALL_DIR/vpsm" "$INSTALL_DIR/vpcm"
    else
        error "Permission denied creating vpcm symlink and 'sudo' is not available."
        exit 1
    fi
else
    ln -s "$INSTALL_DIR/vpsm" "$INSTALL_DIR/vpcm"
fi
success "Linked vpcm to vpsm"

# Wrapper shell insertion
if [[ "$ENABLE_WRAPPER" =~ ^[Yy]$ ]]; then
    info "Configuring shell wrappers..."
    WRAPPER_CONTENT=""
    # First, try to load it from the extracted folder (bundled in the release)
    if [ -f "$TMP_DIR/shell_wrapper.sh" ]; then
        WRAPPER_CONTENT=$(cat "$TMP_DIR/shell_wrapper.sh")
    else
        # Fall back to curl/wget
        temp_wrapper=$(mktemp)
        if download_file "https://raw.githubusercontent.com/$REPO/main/scripts/shell_wrapper.sh" "$temp_wrapper" 2>/dev/null; then
            WRAPPER_CONTENT=$(cat "$temp_wrapper")
        fi
        rm -f "$temp_wrapper"
    fi

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
        warn "Could not fetch shell_wrapper.sh content. Skipping wrapper config."
    fi
fi

printf "\n%s%s✨ VPSM (vpsm/vpcm) has been successfully installed to %s!%s\n" "${GREEN}" "${BOLD}" "$INSTALL_DIR" "${NORMAL}"
printf "Run %s'vpsm version'%s or %s'vpcm version'%s to verify your installation.\n\n" "${CYAN}" "${NORMAL}" "${CYAN}" "${NORMAL}"
