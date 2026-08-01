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
skip() {
    printf "%s[-]%s %s\n" "${YELLOW}" "${NORMAL}" "$1"
}

# Print Banner
printf "%s%s\n" "${RED}" "${BOLD}"
echo "  __   _____  ____  __  __ "
echo "  \ \ / /  _ \/ ___||  \/  |"
echo "   \ V /| |_) \___ \| |\/| |"
echo "    \ / |  __/ ___) | |  | |"
echo "     V  |_|   |____/|_|  |_|"
printf "%s\n" "${NORMAL}"
printf "%s%sVPSM / VPCM Uninstaller%s\n" "${WHITE}" "${BOLD}" "${NORMAL}"
printf "This script will remove all installed VPSM binaries and shell wrappers.\n\n"

# Default variables
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
if [[ "$OS" =~ ^(mingw|msys|cygwin|windows) ]]; then
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

REMOVE_WRAPPER="y"
REMOVE_CONFIG="n"
AUTO_CONFIRM="n"

# Parse CLI arguments
while [[ $# -gt 0 ]]; do
    case "$1" in
        -d|--dir)
            INSTALL_DIR="$2"
            shift 2
            ;;
        -y|--yes)
            AUTO_CONFIRM="y"
            shift
            ;;
        --keep-wrapper)
            REMOVE_WRAPPER="n"
            shift
            ;;
        --keep-config)
            REMOVE_CONFIG="n"
            shift
            ;;
        --remove-config)
            REMOVE_CONFIG="y"
            shift
            ;;
        -h|--help)
            echo "Usage: uninstall.sh [options]"
            echo "Options:"
            echo "  -d, --dir <path>     Installation directory (default: /usr/local/bin or %LOCALAPPDATA%/Programs/vpsm)"
            echo "  -y, --yes            Skip confirmation and uninstall VPSM"
            echo "  --keep-wrapper       Do not remove the shell wrapper override"
            echo "  --keep-config        Do not remove config and data directory (default)"
            echo "  --remove-config      Remove config and data directory"
            echo "  -h, --help           Show this help message"
            exit 0
            ;;
        *)
            shift
            ;;
    esac
done

# Confirm uninstallation if not auto-confirmed
if [ "$AUTO_CONFIRM" != "y" ]; then
    if [ -c /dev/tty ]; then
        printf "%sWarning:%s This will permanently remove VPSM (vpsm/vpcm/vpsmd/vpsm-api) from your system.\n" "${YELLOW}${BOLD}" "${NORMAL}"
        read -p "Are you sure you want to uninstall? [y/N]: " CONFIRM < /dev/tty
        CONFIRM=${CONFIRM:-"n"}
        if [[ ! "$CONFIRM" =~ ^[Yy]$ ]]; then
            warn "Uninstall cancelled."
            exit 0
        fi
    else
        warn "Running in non-interactive environment without -y/--yes flag. Aborting to prevent accidental removal."
        exit 1
    fi
fi

# Expand tilde (~) if present
if [[ "$INSTALL_DIR" == ~* ]]; then
    INSTALL_DIR="${INSTALL_DIR/#\~/$HOME}"
fi

# Clean up path
if [[ "$INSTALL_DIR" != /* ]] && [[ "$INSTALL_DIR" != [A-Za-z]:* ]]; then
    INSTALL_DIR="$(pwd)/$INSTALL_DIR"
fi

# Validate path format
if [[ "$INSTALL_DIR" =~ [\$\{\}\*\?\'\"\|\<\>\&\;] ]]; then
    error "Invalid installation path: $INSTALL_DIR contains illegal characters."
    exit 1
fi

# Prompt for wrapper removal if not auto-confirmed
if [ "$AUTO_CONFIRM" != "y" ] && [ "$REMOVE_WRAPPER" = "y" ]; then
    if [ -c /dev/tty ]; then
        read -p "Remove the VPSM ssh wrapper from your shell profiles? [Y/n]: " RM_WRAP < /dev/tty
        REMOVE_WRAPPER=${RM_WRAP:-"y"}
    fi
fi

printf "\n"

# Remove binaries
BINARIES=("vpsm" "vpsm.exe" "vpcm" "vpcm.exe" "vpsmd" "vpsmd.exe" "vpsm-api" "vpsm-api.exe")
for bin in "${BINARIES[@]}"; do
    BIN_PATH="$INSTALL_DIR/$bin"
    if [ -f "$BIN_PATH" ] || [ -L "$BIN_PATH" ]; then
        if [ ! -w "$INSTALL_DIR" ] || { [ -f "$BIN_PATH" ] && [ ! -w "$BIN_PATH" ]; }; then
            if command -v sudo >/dev/null 2>&1; then
                sudo rm -f "$BIN_PATH"
            else
                error "Permission denied removing $BIN_PATH and 'sudo' is not available."
                exit 1
            fi
        else
            rm -f "$BIN_PATH"
        fi
        success "Removed $BIN_PATH"
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

    # Use a portable awk/sed approach to strip the wrapper block
    local tmp_file
    tmp_file=$(mktemp)
    awk '
        /# VPSM ssh wrapper override/ { skip = 1; next }
        skip && /^}/ { skip = 0; next }
        skip { next }
        { print }
    ' "$SHELL_RC" > "$tmp_file"

    if [ ! -w "$SHELL_RC" ]; then
        if command -v sudo >/dev/null 2>&1; then
            sudo mv "$tmp_file" "$SHELL_RC"
        else
            error "Permission denied modifying $SHELL_RC and 'sudo' is not available."
            rm -f "$tmp_file"
            exit 1
        fi
    else
        mv "$tmp_file" "$SHELL_RC"
    fi
    success "Removed VPSM wrapper from $SHELL_RC"
}

if [[ "$REMOVE_WRAPPER" =~ ^[Yy]$ ]]; then
    info "Removing shell wrappers..."
    remove_wrapper_from_file "$HOME/.zshrc"
    remove_wrapper_from_file "$HOME/.bashrc"
    remove_wrapper_from_file "$HOME/.bash_profile"
    remove_wrapper_from_file "$HOME/.profile"
fi

# Remove local config and data directories
CONFIG_DIR="$HOME/.config/vpsm"
if [ -d "$CONFIG_DIR" ]; then
    if [ "$AUTO_CONFIRM" != "y" ] && [ "$REMOVE_CONFIG" = "n" ]; then
        if [ -t 0 ] && [ -c /dev/tty ]; then
            read -p "Remove config and data directory ($CONFIG_DIR)? [y/N]: " RM_CONF < /dev/tty
            REMOVE_CONFIG=${RM_CONF:-"n"}
        fi
    fi

    if [[ "$REMOVE_CONFIG" =~ ^[Yy]$ ]]; then
        rm -rf "$CONFIG_DIR"
        success "Removed config directory $CONFIG_DIR"
    else
        skip "Keeping config directory $CONFIG_DIR"
    fi
else
    skip "No config directory found at $CONFIG_DIR"
fi

printf "\n%s%s✨ VPSM has been successfully uninstalled.%s\n" "${GREEN}" "${BOLD}" "${NORMAL}"
printf "Thank you for using VPSM! Visit %shttps://github.com/%s%s for more information.\n\n" "${CYAN}" "$REPO" "${NORMAL}"
