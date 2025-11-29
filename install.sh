#!/usr/bin/env bash
set -euo pipefail

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

YES_TO_ALL=false

# Usage
usage() {
    echo "Usage: $0 [OPTIONS]"
    echo "Options:"
    echo "  -y, --yes       Automatic yes to prompts (non-interactive mode)"
    echo "  --uninstall     Uninstall zshellcheck"
    echo "  -h, --help      Show this help message"
}

# Parse Arguments
while [[ $# -gt 0 ]]; do
    key="$1"
    case $key in
        -y|--yes) 
            YES_TO_ALL=true
            shift
            ;; 
        --uninstall) 
            # Handled later, but consume it here to avoid errors if combined
            shift
            ;; 
        -h|--help) 
            usage
            exit 0
            ;; 
        *) 
            # Check if it's the uninstall flag passed as $1 in previous logic (legacy support)
            if [[ "$1" == "--uninstall" ]]; then
                shift
            else
                echo "Unknown option: $1"
                usage
                exit 1
            fi
            ;; 
    esac
done

# Detect Shell Config
detect_shell_config() {
    local shell_name
    shell_name=$(basename "$SHELL")
    case "$shell_name" in
        zsh) echo "$HOME/.zshrc" ;; 
        bash) echo "$HOME/.bashrc" ;; 
        *) echo "" ;; 
    esac
}

# Ask for confirmation
ask_yes_no() {
    local prompt="$1"
    
    if [ "$YES_TO_ALL" = true ]; then
        return 0
    fi

    if [ -t 0 ]; then # Only ask if interactive
        read -p "$prompt [y/N] " -r
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            return 1
        fi
        return 0
    else
        # Default to NO in non-interactive mode without -y
        return 1
    fi
}

# Check if line exists in file
line_exists() {
    grep -Fq "$1" "$2" 2>/dev/null
}

# Uninstall function
uninstall() {
    echo -e "${YELLOW}Uninstalling zshellcheck...${NC}"
    
    if [ "$EUID" -eq 0 ]; then
        BIN_DIR="/usr/local/bin"
        MAN_DIR="/usr/local/share/man/man1"
        ZSH_COMP_DIR="/usr/local/share/zsh/site-functions"
        BASH_COMP_DIR="/usr/local/share/bash-completion/completions"
    else
        BIN_DIR="$HOME/.local/bin"
        MAN_DIR="$HOME/.local/share/man/man1"
        ZSH_COMP_DIR="$HOME/.local/share/zsh/site-functions"
        BASH_COMP_DIR="$HOME/.local/share/bash-completion/completions"
    fi

    rm -v "$BIN_DIR/zshellcheck" 2>/dev/null || true
    rm -v "$MAN_DIR/zshellcheck.1" 2>/dev/null || true
    rm -v "$ZSH_COMP_DIR/_zshellcheck" 2>/dev/null || true
    rm -v "$BASH_COMP_DIR/zshellcheck" 2>/dev/null || true
    
    echo -e "${GREEN}Uninstallation complete.${NC}"
}

# Re-check uninstall arg logic just in case
if [[ "${1:-}" == "--uninstall" ]] || [[ "${*:-}" == *"--uninstall"* ]]; then
    uninstall
    exit 0
fi

echo -e "${GREEN}Installing zshellcheck...${NC}"

# --- BUILD OR DOWNLOAD ---

BUILD_SUCCESS=false

# cleanup trap
cleanup() {
    rm -f zshellcheck zshellcheck.tar.gz checksums.txt
}
trap cleanup EXIT

# Check for Go
if command -v go &> /dev/null; then
    # Determine Version
    VERSION="dev"
    if command -v git &> /dev/null && [ -d .git ]; then
        VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
    fi

    echo -e "Go found. Building binary from source (Version: ${BLUE}${VERSION}${NC})..."
    if go build -ldflags "-X github.com/afadesigns/zshellcheck/pkg/version.Version=${VERSION}" -o zshellcheck cmd/zshellcheck/main.go; then
        BUILD_SUCCESS=true
    else
        echo -e "${RED}Build failed.${NC}"
    fi
fi

if [ "$BUILD_SUCCESS" = false ]; then
    echo -e "${YELLOW}Go not found or build failed.${NC}"
    echo -e "Attempting to download latest binary release..."

    if ! command -v curl &> /dev/null; then
        echo -e "${RED}Error: curl is required for downloading releases.${NC}"
        exit 1
    fi
    if ! command -v tar &> /dev/null; then
        echo -e "${RED}Error: tar is required for extracting releases.${NC}"
        exit 1
    fi

    # Detect OS/Arch
    OS=$(uname -s)
    ARCH=$(uname -m)
    
    # Map to Goreleaser names
    case "$OS" in
        Linux) GOOS="Linux" ;; 
        Darwin) GOOS="Darwin" ;; 
        *) echo -e "${RED}Unsupported OS: $OS${NC}"; exit 1 ;; 
    esac

    case "$ARCH" in
        x86_64) GOARCH="x86_64" ;; 
        aarch64|arm64) GOARCH="arm64" ;; 
        i386) GOARCH="i386" ;; 
        *) echo -e "${RED}Unsupported Arch: $ARCH${NC}"; exit 1 ;; 
    esac

    echo -e "Detected platform: ${BLUE}$GOOS $GOARCH${NC}"

    # Fetch latest tag
    LATEST_TAG=$(curl -s https://api.github.com/repos/afadesigns/zshellcheck/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    
    if [ -z "$LATEST_TAG" ]; then
        echo -e "${RED}Failed to fetch latest release info from GitHub.${NC}"
        exit 1
    fi

    # Construct URL
    FILENAME="zshellcheck_${GOOS}_${GOARCH}.tar.gz"
    URL="https://github.com/afadesigns/zshellcheck/releases/download/${LATEST_TAG}/${FILENAME}"
    CHECKSUM_URL="https://github.com/afadesigns/zshellcheck/releases/download/${LATEST_TAG}/checksums.txt"

    echo -e "Downloading version ${BLUE}${LATEST_TAG}${NC}..."
    
    # Download Checksums
    if curl -sL -o "checksums.txt" "$CHECKSUM_URL"; then
        echo -e "Verifying checksums..."
    else 
        echo -e "${YELLOW}Warning: Could not download checksums.txt. Skipping verification.${NC}"
    fi

    # Download Binary
    if curl -L -o "$FILENAME" "$URL"; then
        # Verify Checksum if file exists
        if [ -f "checksums.txt" ]; then
            # Use sha256sum if available, else shasum
            if command -v sha256sum &> /dev/null; then
                if grep "$FILENAME" checksums.txt | sha256sum -c - --status; then
                    echo -e "${GREEN}Checksum verified.${NC}"
                else
                    echo -e "${RED}Checksum verification failed! Aborting.${NC}"
                    exit 1
                fi
            elif command -v shasum &> /dev/null; then
                 if grep "$FILENAME" checksums.txt | shasum -a 256 -c - --status; then
                    echo -e "${GREEN}Checksum verified.${NC}"
                else
                    echo -e "${RED}Checksum verification failed! Aborting.${NC}"
                    exit 1
                fi
            else
                echo -e "${YELLOW}sha256sum not found, skipping verification.${NC}"
            fi
        fi

        echo "Extracting..."
        tar -xzf "$FILENAME" zshellcheck
        BUILD_SUCCESS=true
    else
        echo -e "${RED}Download failed.${NC}"
        exit 1
    fi
fi

# --- INSTALLATION ---

# Determine install locations
if [ "$EUID" -eq 0 ]; then
    BIN_DIR="/usr/local/bin"
    MAN_DIR="/usr/local/share/man/man1"
    ZSH_COMP_DIR="/usr/local/share/zsh/site-functions"
    BASH_COMP_DIR="/usr/local/share/bash-completion/completions"
else
    BIN_DIR="$HOME/.local/bin"
    MAN_DIR="$HOME/.local/share/man/man1"
    ZSH_COMP_DIR="$HOME/.local/share/zsh/site-functions"
    BASH_COMP_DIR="$HOME/.local/share/bash-completion/completions"
fi

echo -e "Installing binary to ${BLUE}$BIN_DIR${NC}..."
mkdir -p "$BIN_DIR"
if mv zshellcheck "$BIN_DIR/zshellcheck"; then
    echo -e "${GREEN}✓ Binary installed.${NC}"
else
    echo -e "${RED}Failed to install binary.${NC}"
    rm -f zshellcheck
    exit 1
fi

# Install resources (Man page & Completions)
# These exist in the repo, so we can install them even if we downloaded the binary
# (assuming we are running install.sh from the repo)

if [ -f "man/man1/zshellcheck.1" ]; then
    echo -e "Installing man page to ${BLUE}$MAN_DIR${NC}..."
    mkdir -p "$MAN_DIR"
    cp "man/man1/zshellcheck.1" "$MAN_DIR/zshellcheck.1"
    echo -e "${GREEN}✓ Man page installed.${NC}"
fi

if [ -f "completions/zsh/_zshellcheck" ]; then
    echo -e "Installing Zsh completions to ${BLUE}$ZSH_COMP_DIR${NC}..."
    mkdir -p "$ZSH_COMP_DIR"
    cp "completions/zsh/_zshellcheck" "$ZSH_COMP_DIR/_zshellcheck"
    echo -e "${GREEN}✓ Zsh completions installed.${NC}"
fi

if [ -f "completions/bash/zshellcheck-completion.bash" ]; then
    echo -e "Installing Bash completions to ${BLUE}$BASH_COMP_DIR${NC}..."
    mkdir -p "$BASH_COMP_DIR"
    cp "completions/bash/zshellcheck-completion.bash" "$BASH_COMP_DIR/zshellcheck"
    echo -e "${GREEN}✓ Bash completions installed.${NC}"
fi

# --- CONFIGURATION ---

echo ""
echo -e "${GREEN}Installation complete!${NC}"

SHELL_CONFIG=$(detect_shell_config)

# Path check
if [[ ":$PATH:" != ".*:${BIN_DIR}:"* ]]; then
    echo ""
    echo -e "${YELLOW}WARNING: $BIN_DIR is not in your PATH.${NC}"
    EXPORT_CMD="export PATH=\"
$PATH:$BIN_DIR\""
    
    if [ -n "$SHELL_CONFIG" ]; then
        if line_exists "$EXPORT_CMD" "$SHELL_CONFIG"; then
             echo -e "${GREEN}✓ PATH export already exists in $SHELL_CONFIG.${NC}"
        elif ask_yes_no "Would you like to append the PATH export to $SHELL_CONFIG?"; then
            echo "" >> "$SHELL_CONFIG"
            echo "# Added by zshellcheck installer" >> "$SHELL_CONFIG"
            echo "$EXPORT_CMD" >> "$SHELL_CONFIG"
            echo -e "${GREEN}✓ Added to $SHELL_CONFIG.${NC} Please restart your shell or run 'source $SHELL_CONFIG'."
        else
            echo "Please add this to your config:"
            echo -e "  ${BLUE}$EXPORT_CMD${NC}"
        fi
    else
        echo "Please add this to your config:"
        echo -e "  ${BLUE}$EXPORT_CMD${NC}"
    fi
fi

# Fpath check
if [ "$EUID" -ne 0 ] && [[ "$SHELL" == *"zsh"* ]]; then
    FPATH_CMD="fpath+=($ZSH_COMP_DIR)"
    
    echo ""
    echo -e "${BLUE}Zsh Completions:${NC}"
    if [ -n "$SHELL_CONFIG" ]; then
         if line_exists "$FPATH_CMD" "$SHELL_CONFIG"; then
            echo -e "${GREEN}✓ fpath update already exists in $SHELL_CONFIG.${NC}"
         elif ask_yes_no "Would you like to add the completion directory to your fpath in $SHELL_CONFIG?"; then
            echo "" >> "$SHELL_CONFIG"
            echo "# Added by zshellcheck installer" >> "$SHELL_CONFIG"
            echo "$FPATH_CMD" >> "$SHELL_CONFIG"
            echo -e "${GREEN}✓ Added to $SHELL_CONFIG.${NC}"
            echo -e "${YELLOW}Note: Ensure this line appears BEFORE 'autoload -Uz compinit && compinit'.${NC}"
        else
            echo "Add this to ~/.zshrc before 'compinit':"
            echo -e "  ${BLUE}$FPATH_CMD${NC}"
        fi
    else
        echo "Add this to ~/.zshrc before 'compinit':"
        echo -e "  ${BLUE}$FPATH_CMD${NC}"
    fi
fi

echo ""
echo "Run 'zshellcheck --help' to get started."
