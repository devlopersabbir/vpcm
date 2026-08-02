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

EXE_EXT=""
ARCHIVE_EXT="tar.gz"

case "$OS" in
    darwin)
        OS="darwin"
        ;;
    linux)
        OS="linux"
        ;;
    mingw*|msys*|cygwin*|windows*)
        OS="windows"
        EXE_EXT=".exe"
        ARCHIVE_EXT="zip"
        ;;
    *)
        error "Unsupported operating system $OS"
        exit 1
        ;;
esac

# Default variables
if [ "$OS" = "windows" ]; then
    if [ -n "$LOCALAPPDATA" ]; then
        WIN_DIR=$(echo "$LOCALAPPDATA/Programs/vpsm" | tr '\\' '/')
        if command -v cygpath >/dev/null 2>&1; then
            INSTALL_DIR=$(cygpath -u "$WIN_DIR")
        else
            INSTALL_DIR="$WIN_DIR"
        fi
    else
        INSTALL_DIR="$HOME/bin"
    fi
else
    INSTALL_DIR="/usr/local/bin"
fi

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
            echo "  -d, --dir <path>     Installation directory (default: /usr/local/bin or %LOCALAPPDATA%/Programs/vpsm)"
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
FILENAME="vpsm-${OS}-${ARCH}.${ARCHIVE_EXT}"
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
if [[ "$INSTALL_DIR" != /* ]] && [[ "$INSTALL_DIR" != [A-Za-z]:* ]]; then
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
if [ "$ARCHIVE_EXT" = "zip" ]; then
    if command -v unzip >/dev/null 2>&1; then
        unzip -q "$TMP_DIR/$FILENAME" -d "$TMP_DIR"
    elif command -v powershell.exe >/dev/null 2>&1; then
        WIN_TMP="$TMP_DIR/$FILENAME"
        if command -v cygpath >/dev/null 2>&1; then
            WIN_TMP=$(cygpath -w "$TMP_DIR/$FILENAME")
            WIN_DEST=$(cygpath -w "$TMP_DIR")
        else
            WIN_DEST="$TMP_DIR"
        fi
        powershell.exe -Command "Expand-Archive -Path '$WIN_TMP' -DestinationPath '$WIN_DEST' -Force"
    elif command -v powershell >/dev/null 2>&1; then
        powershell -Command "Expand-Archive -Path '$TMP_DIR/$FILENAME' -DestinationPath '$TMP_DIR' -Force"
    else
        error "Neither unzip nor powershell is installed to unpack $FILENAME."
        exit 1
    fi
else
    tar -xzf "$TMP_DIR/$FILENAME" -C "$TMP_DIR"
fi

# Ensure install dir exists
if [ ! -d "$INSTALL_DIR" ]; then
    info "Creating installation directory $INSTALL_DIR..."
    parent_dir=$(dirname "$INSTALL_DIR")
    if [ ! -w "$parent_dir" ] && [ "$OS" != "windows" ]; then
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
    local BIN_BASE="$1"
    local BIN_NAME="${BIN_BASE}${EXE_EXT}"
    local SRC="$TMP_DIR/$BIN_NAME"
    local DEST="$INSTALL_DIR/$BIN_NAME"
    
    if [ ! -f "$SRC" ]; then
        SRC=$(find "$TMP_DIR" -type f -name "$BIN_NAME" | head -n 1)
    fi

    if [ -z "$SRC" ] || [ ! -f "$SRC" ]; then
        error "Failed to locate binary $BIN_NAME in the release archive."
        exit 1
    fi

    # Check if target file or directory is writable
    if [ ! -w "$INSTALL_DIR" ] || { [ -f "$DEST" ] && [ ! -w "$DEST" ]; }; then
        if command -v sudo >/dev/null 2>&1; then
            sudo mv "$SRC" "$DEST"
            [ -z "$EXE_EXT" ] && sudo chmod +x "$DEST"
        else
            error "Permission denied writing to $DEST and 'sudo' is not available."
            exit 1
        fi
    else
        mv "$SRC" "$DEST"
        [ -z "$EXE_EXT" ] && chmod +x "$DEST"
    fi
    success "Installed $BIN_NAME to $INSTALL_DIR"
}

# Install all three core binaries
install_binary "vpsm"
install_binary "vpsmd"
install_binary "vpsm-api"

# Create symlink/wrapper/copy for vpcm if it doesn't exist
VPCM_DEST="$INSTALL_DIR/vpcm${EXE_EXT}"
VPSM_DEST="$INSTALL_DIR/vpsm${EXE_EXT}"

if [ -L "$VPCM_DEST" ] || [ -f "$VPCM_DEST" ]; then
    if [ ! -w "$INSTALL_DIR" ] || { [ -f "$VPCM_DEST" ] && [ ! -w "$VPCM_DEST" ]; }; then
        if command -v sudo >/dev/null 2>&1; then
            sudo rm -f "$VPCM_DEST"
        else
            error "Permission denied removing old vpcm and 'sudo' is not available."
            exit 1
        fi
    else
        rm -f "$VPCM_DEST"
    fi
fi

if [ "$OS" = "windows" ]; then
    cp "$VPSM_DEST" "$VPCM_DEST"
    success "Linked vpcm${EXE_EXT} to vpsm${EXE_EXT}"
else
    if [ ! -w "$INSTALL_DIR" ]; then
        if command -v sudo >/dev/null 2>&1; then
            sudo ln -s "$VPSM_DEST" "$VPCM_DEST"
        else
            error "Permission denied creating vpcm symlink and 'sudo' is not available."
            exit 1
        fi
    else
        ln -s "$VPSM_DEST" "$VPCM_DEST"
    fi
    success "Linked vpcm to vpsm"
fi

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
# Auto-initialize default configuration (SQLite & API server enabled on 127.0.0.1)
CONFIG_DIR="$HOME/.config/vpsm"
CONFIG_FILE="$CONFIG_DIR/config.yaml"
if [ ! -f "$CONFIG_FILE" ]; then
    mkdir -p "$CONFIG_DIR"
    cat <<EOF > "$CONFIG_FILE"
database:
  driver: sqlite
  path: $HOME/.local/share/vpsm/vpsm.db
api:
  enabled: true
  host: 127.0.0.1
  port: 8080
  mode: local
  global_url: http://127.0.0.1:8080
ssh:
  timeout: 10s
logging:
  level: info
  format: pretty
EOF
    success "Auto-initialized default configuration with SQLite & REST API enabled on 127.0.0.1"
fi

printf "\n%s%s✨ VPSM (vpsm/vpcm) has been successfully installed to %s!%s\n" "${GREEN}" "${BOLD}" "$INSTALL_DIR" "${NORMAL}"
printf "Run %s'vpsm version'%s or %s'vpcm version'%s to verify your installation.\n\n" "${CYAN}" "${NORMAL}" "${CYAN}" "${NORMAL}"
