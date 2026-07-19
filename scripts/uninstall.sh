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

# Helper functions
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
skip() {
    echo -e "${YELLOW}[-]${NORMAL} $1"
}

# Print Banner
echo -e "${RED}${BOLD}"
echo "  __   _____  ____  __  __ "
echo "  \ \ / /  _ \/ ___||  \/  |"
echo "   \ V /| |_) \___ \| |\/| |"
echo "    \ / |  __/ ___) | |  | |"
echo "     V  |_|   |____/|_|  |_|"
echo -e "${NORMAL}"
echo -e "${WHITE}${BOLD}VPSM / VPCM Uninstaller${NORMAL}"
echo -e "This script will remove all installed VPSM binaries and shell wrappers.\n"

# Interactive Step 1: Confirm before removing
echo -e "${YELLOW}${BOLD}Warning:${NORMAL} This will permanently remove VPSM (vpsm/vpcm/vpsmd/vpsm-api) from your system."
read -p "Are you sure you want to uninstall? [y/N]: " CONFIRM
CONFIRM=${CONFIRM:-"n"}
if [[ ! "$CONFIRM" =~ ^[Yy]$ ]]; then
    warn "Uninstall cancelled."
    exit 0
fi

# Interactive Step 2: Installation directory to clean
echo -e "\n${BOLD}Step 1: Provide the installation directory to clean${NORMAL}"
DEFAULT_INSTALL_DIR="/usr/local/bin"
read -p "Installation path [default: $DEFAULT_INSTALL_DIR]: " INSTALL_DIR
INSTALL_DIR=${INSTALL_DIR:-$DEFAULT_INSTALL_DIR}

# Interactive Step 3: Remove shell wrappers
echo -e "\n${BOLD}Step 2: Remove shell wrapper from shell profiles${NORMAL}"
read -p "Remove the VPSM ssh wrapper from your shell profiles? [Y/n]: " REMOVE_WRAPPER
REMOVE_WRAPPER=${REMOVE_WRAPPER:-"y"}

echo ""

# Remove binaries
BINARIES=("vpsm" "vpcm" "vpsmd" "vpsm-api")
for bin in "${BINARIES[@]}"; do
    BIN_PATH="$INSTALL_DIR/$bin"
    if [ -f "$BIN_PATH" ] || [ -L "$BIN_PATH" ]; then
        if [ ! -w "$INSTALL_DIR" ]; then
            sudo rm -f "$BIN_PATH"
        else
            rm -f "$BIN_PATH"
        fi
        success "Removed $BIN_PATH"
    else
        skip "$bin not found in $INSTALL_DIR — skipping"
    fi
done

# Remove shell wrapper blocks from shell profiles
remove_wrapper_from_file() {
    local SHELL_RC="$1"
    if [ ! -f "$SHELL_RC" ]; then
        return
    fi
    if ! grep -q "VPSM ssh wrapper override" "$SHELL_RC"; then
        skip "No VPSM wrapper found in $SHELL_RC — skipping"
        return
    fi

    # Use a portable sed approach to strip the wrapper block
    TMP_FILE=$(mktemp)
    # Remove from blank line before comment through the closing brace of the ssh() function
    awk '
        /# VPSM ssh wrapper override/ { skip = 1; next }
        skip && /^}/ { skip = 0; next }
        skip { next }
        { print }
    ' "$SHELL_RC" > "$TMP_FILE"
    mv "$TMP_FILE" "$SHELL_RC"
    success "Removed VPSM wrapper from $SHELL_RC"
}

if [[ "$REMOVE_WRAPPER" =~ ^[Yy]$ ]]; then
    info "Removing shell wrappers..."
    remove_wrapper_from_file "$HOME/.zshrc"
    remove_wrapper_from_file "$HOME/.bashrc"
fi

# Remove local data and config directories (optional)
echo -e "\n${BOLD}Step 3: Clean local config and data${NORMAL}"
CONFIG_DIR="$HOME/.config/vpsm"
if [ -d "$CONFIG_DIR" ]; then
    read -p "Remove config and data directory ($CONFIG_DIR)? [y/N]: " REMOVE_CONFIG
    REMOVE_CONFIG=${REMOVE_CONFIG:-"n"}
    if [[ "$REMOVE_CONFIG" =~ ^[Yy]$ ]]; then
        rm -rf "$CONFIG_DIR"
        success "Removed config directory $CONFIG_DIR"
    else
        skip "Keeping config directory $CONFIG_DIR"
    fi
else
    skip "No config directory found at $CONFIG_DIR"
fi

echo -e "\n${GREEN}${BOLD}✨ VPSM has been successfully uninstalled.${NORMAL}"
echo -e "Thank you for using VPSM! Visit ${CYAN}https://github.com/$REPO${NORMAL} for more information.\n"
