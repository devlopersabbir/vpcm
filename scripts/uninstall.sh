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
printf "%s%sComplete VPSM / VPCM Uninstaller%s\n" "${WHITE}" "${BOLD}" "${NORMAL}"
printf "This script cleanly removes all VPSM binaries, background daemons, configurations,\n"
printf "desktop entries, and shell wrappers while safely preserving your SQLite database.\n\n"

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

AUTO_CONFIRM="n"
PURGE_DB="n"

# Parse CLI arguments
while [[ $# -gt 0 ]]; do
    case "$1" in
        -d|--dir)
            INSTALL_DIR="$2"
            shift 2
            ;;
        -y|--yes|--non-interactive)
            AUTO_CONFIRM="y"
            shift
            ;;
        --purge-db)
            PURGE_DB="y"
            shift
            ;;
        -h|--help)
            echo "Usage: uninstall.sh [options]"
            echo "Options:"
            echo "  -d, --dir <path>     Primary installation directory (default: /usr/local/bin or %LOCALAPPDATA%/Programs/vpsm)"
            echo "  -y, --yes            Skip confirmation and perform clean uninstall"
            echo "  --purge-db           Purge the SQLite database as well (by default it is preserved)"
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
        printf "%sNotice:%s This will remove all VPSM/VPCM binaries, running daemons, configs, and shell wrappers.\n" "${YELLOW}${BOLD}" "${NORMAL}"
        printf "%sNote:%s Your SQLite database will NOT be removed unless --purge-db is specified.\n\n" "${GREEN}${BOLD}" "${NORMAL}"
        read -p "Are you sure you want to proceed with uninstallation? [y/N]: " CONFIRM < /dev/tty
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

printf "\n"

# 1. Stop background processes & daemons
info "Stopping running background processes and daemons..."
pkill -f "vpsm-api" 2>/dev/null || true
pkill -f "vpsmd" 2>/dev/null || true
pkill -f "vpsm-desktop" 2>/dev/null || true
pkill -x "vpsm" 2>/dev/null || true
pkill -x "vpcm" 2>/dev/null || true
success "Terminated active VPSM processes and background daemons."

# 2. Remove binaries and symlinks
info "Scanning and removing binaries from system paths..."
SEARCH_DIRS=(
    "$INSTALL_DIR"
    "/usr/local/bin"
    "/opt/homebrew/bin"
    "/usr/bin"
    "$HOME/.local/bin"
    "$HOME/bin"
    "$HOME/go/bin"
    "${GOPATH:+${GOPATH}/bin}"
)

BINARIES=("vpsm" "vpsm.exe" "vpcm" "vpcm.exe" "vpsmd" "vpsmd.exe" "vpsm-api" "vpsm-api.exe" "vpsm-desktop" "vpsm-desktop.exe")

for dir in "${SEARCH_DIRS[@]}"; do
    if [ -z "$dir" ] || [ ! -d "$dir" ]; then
        continue
    fi
    for bin in "${BINARIES[@]}"; do
        BIN_PATH="$dir/$bin"
        if [ -f "$BIN_PATH" ] || [ -L "$BIN_PATH" ]; then
            if [ ! -w "$dir" ] || { [ -f "$BIN_PATH" ] && [ ! -w "$BIN_PATH" ]; }; then
                if command -v sudo >/dev/null 2>&1; then
                    sudo rm -f "$BIN_PATH"
                else
                    warn "Permission denied removing $BIN_PATH and 'sudo' is not available."
                    continue
                fi
            else
                rm -f "$BIN_PATH"
            fi
            success "Removed $BIN_PATH"
        fi
    done
done

# 3. Clean Shell configuration files (wrappers, completion, aliases)
info "Cleaning shell profiles and RC files..."

clean_shell_file() {
    local SHELL_RC="$1"
    if [ ! -f "$SHELL_RC" ]; then
        return
    fi

    # Check if the file contains any vpsm/vpcm references or wrappers
    if ! grep -q -E "(VPSM ssh wrapper override|VPCM ssh wrapper override|vpsm|vpcm)" "$SHELL_RC" 2>/dev/null; then
        return
    fi

    local tmp_file
    tmp_file=$(mktemp)

    awk '
    BEGIN { in_vpsm_comment = 0; in_ssh_func = 0; ssh_buf = ""; }
    /# VPSM ssh wrapper override/ || /# VPCM ssh wrapper override/ {
        in_vpsm_comment = 1
        next
    }
    in_vpsm_comment {
        if ($0 ~ /^}/) {
            in_vpsm_comment = 0
        }
        next
    }
    /^(function[[:space:]]+)?ssh([[:space:]]*\(\))?[[:space:]]*\{/ {
        in_ssh_func = 1
        ssh_buf = $0 "\n"
        next
    }
    in_ssh_func {
        ssh_buf = ssh_buf $0 "\n"
        if ($0 ~ /^}/) {
            in_ssh_func = 0
            if (ssh_buf ~ /vpsm/ || ssh_buf ~ /vpcm/) {
                ssh_buf = ""
                next
            } else {
                printf "%s", ssh_buf
                ssh_buf = ""
                next
            }
        }
        next
    }
    /vpsm[[:space:]]+completion/ || /vpcm[[:space:]]+completion/ || /alias[[:space:]]+vpsm=/ || /alias[[:space:]]+vpcm=/ {
        next
    }
    {
        print
    }
    ' "$SHELL_RC" > "$tmp_file"

    if [ ! -w "$SHELL_RC" ]; then
        if command -v sudo >/dev/null 2>&1; then
            sudo mv "$tmp_file" "$SHELL_RC"
        else
            warn "Permission denied modifying $SHELL_RC."
            rm -f "$tmp_file"
            return
        fi
    else
        mv "$tmp_file" "$SHELL_RC"
    fi
    success "Cleaned shell profile: $SHELL_RC"
}

SHELL_PROFILES=(
    "$HOME/.zshrc"
    "$HOME/.bashrc"
    "$HOME/.bash_profile"
    "$HOME/.profile"
    "$HOME/.zprofile"
    "$HOME/.zshenv"
    "$HOME/.config/fish/config.fish"
)

for profile in "${SHELL_PROFILES[@]}"; do
    clean_shell_file "$profile"
done

# 4. Remove configuration directories and log files
info "Removing configuration directories and logs..."
CONFIG_DIRS=(
    "$HOME/.config/vpsm"
    "$HOME/.config/vpcm"
)

for cdir in "${CONFIG_DIRS[@]}"; do
    if [ -d "$cdir" ]; then
        rm -rf "$cdir"
        success "Removed config directory: $cdir"
    fi
done

# 5. Remove Desktop Launchers, App Bundles, and Icons
info "Removing desktop entries, applications, and icons..."
DESKTOP_FILES=(
    "$HOME/.local/share/applications/vpsm-desktop.desktop"
    "$HOME/.local/share/applications/vpcm-desktop.desktop"
    "$HOME/.local/share/applications/vpsm.desktop"
    "/usr/share/applications/vpsm-desktop.desktop"
    "/Applications/vpsm-desktop.app"
    "$HOME/Applications/vpsm-desktop.app"
)

for dfile in "${DESKTOP_FILES[@]}"; do
    if [ -e "$dfile" ]; then
        if [ ! -w "$dfile" ] && command -v sudo >/dev/null 2>&1; then
            sudo rm -rf "$dfile"
        else
            rm -rf "$dfile"
        fi
        success "Removed desktop entry / application: $dfile"
    fi
done

# Clean icon files if present
find "$HOME/.local/share/icons" -name "*vpsm*" -exec rm -f {} + 2>/dev/null || true

# 6. Remove temporary / cache files
rm -f /tmp/vpsm* /tmp/vpcm* 2>/dev/null || true

# 7. SQLite database preservation / handling
DEFAULT_DB_DIR="$HOME/.local/share/vpsm"
DEFAULT_DB_PATH="$DEFAULT_DB_DIR/vpsm.db"

if [ "$PURGE_DB" = "y" ]; then
    if [ -d "$DEFAULT_DB_DIR" ]; then
        rm -rf "$DEFAULT_DB_DIR"
        warn "Purged SQLite database directory at: $DEFAULT_DB_DIR"
    fi
else
    if [ -f "$DEFAULT_DB_PATH" ]; then
        success "Preserved SQLite database at: ${CYAN}${DEFAULT_DB_PATH}${NORMAL}"
        info "  (Your server inventory and data are safe and intact)"
    else
        info "No local SQLite database found at $DEFAULT_DB_PATH."
    fi
fi

printf "\n%s%s✨ VPSM / VPCM has been completely uninstalled from everywhere.%s\n" "${GREEN}" "${BOLD}" "${NORMAL}"

printf "\n%s%s[!] Important Notice for Active Shell Sessions:%s\n" "${YELLOW}" "${BOLD}" "${NORMAL}"
printf "Your current terminal session may still have the SSH wrapper loaded in memory.\n"
printf "To refresh immediately, run:\n"
printf "  %sunfunction ssh%s    (or restart your terminal with %sexec zsh%s / %sexec bash%s)\n\n" "${CYAN}" "${NORMAL}" "${CYAN}" "${NORMAL}" "${CYAN}" "${NORMAL}"
