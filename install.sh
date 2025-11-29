#!/usr/bin/env bash
set -euo pipefail

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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
    if [ -t 0 ]; then # Only ask if interactive
        read -p "$prompt [y/N] " -r
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            return 1
        fi
        return 0
    else
        return 1
    fi
}

echo -e "${GREEN}Installing zshellcheck...${NC}"

# Check for Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed. Please install Go 1.25 or later.${NC}"
    exit 1
fi

# Determine Version
VERSION="dev"
if command -v git &> /dev/null && [ -d .git ]; then
    VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
fi

# Build
echo -e "Building binary (Version: ${BLUE}${VERSION}${NC})..."
# Inject version if possible. Assumes 'github.com/afadesigns/zshellcheck/pkg/version.Version' exists.
if ! go build -ldflags "-X github.com/afadesigns/zshellcheck/pkg/version.Version=${VERSION}" -o zshellcheck cmd/zshellcheck/main.go; then
    echo -e "${RED}Build failed.${NC}"
    exit 1
fi

# Determine install locations
if [ "$EUID" -eq 0 ]; then
    BIN_DIR="/usr/local/bin"
    MAN_DIR="/usr/local/share/man/man1"
    ZSH_COMP_DIR="/usr/local/share/zsh/site-functions"
    BASH_COMP_DIR="/usr/local/share/bash-completion/completions"
else
    # Prefer ~/.local for non-root users
    BIN_DIR="$HOME/.local/bin"
    MAN_DIR="$HOME/.local/share/man/man1"
    ZSH_COMP_DIR="$HOME/.local/share/zsh/site-functions"
    BASH_COMP_DIR="$HOME/.local/share/bash-completion/completions"
fi

# --- Install Binary ---
echo -e "Installing binary to ${BLUE}$BIN_DIR${NC}..."
mkdir -p "$BIN_DIR"
if mv zshellcheck "$BIN_DIR/zshellcheck"; then
    echo -e "${GREEN}✓ Binary installed.${NC}"
else
    echo -e "${RED}Failed to install binary.${NC}"
    rm -f zshellcheck
    exit 1
fi

# --- Install Man Page ---
if [ -f "man/man1/zshellcheck.1" ]; then
    echo -e "Installing man page to ${BLUE}$MAN_DIR${NC}..."
    mkdir -p "$MAN_DIR"
    cp "man/man1/zshellcheck.1" "$MAN_DIR/zshellcheck.1"
    echo -e "${GREEN}✓ Man page installed.${NC}"
else
    echo -e "${YELLOW}Man page not found, skipping.${NC}"
fi

# --- Install Zsh Completion ---
if [ -f "completions/zsh/_zshellcheck" ]; then
    echo -e "Installing Zsh completions to ${BLUE}$ZSH_COMP_DIR${NC}..."
    mkdir -p "$ZSH_COMP_DIR"
    cp "completions/zsh/_zshellcheck" "$ZSH_COMP_DIR/_zshellcheck"
    echo -e "${GREEN}✓ Zsh completions installed.${NC}"
else
    echo -e "${YELLOW}Zsh completions not found, skipping.${NC}"
fi

# --- Install Bash Completion ---
if [ -f "completions/bash/zshellcheck-completion.bash" ]; then
    echo -e "Installing Bash completions to ${BLUE}$BASH_COMP_DIR${NC}..."
    mkdir -p "$BASH_COMP_DIR"
    cp "completions/bash/zshellcheck-completion.bash" "$BASH_COMP_DIR/zshellcheck"
    echo -e "${GREEN}✓ Bash completions installed.${NC}"
else
    echo -e "${YELLOW}Bash completions not found, skipping.${NC}"
fi

# --- Final Checks & Auto-Config ---
echo ""
echo -e "${GREEN}Installation complete!${NC}"

SHELL_CONFIG=$(detect_shell_config)

# Path check
if [[ ":$PATH:" != ".*:*:*:$BIN_DIR:"* ]]; then
    echo ""
    echo -e "${YELLOW}WARNING: $BIN_DIR is not in your PATH.${NC}"
    
    EXPORT_CMD="export PATH=\"
$PATH:$BIN_DIR\""
    
    if [ -n "$SHELL_CONFIG" ]; then
        echo -e "Detected shell config: ${BLUE}$SHELL_CONFIG${NC}"
        if ask_yes_no "Would you like to append the PATH export to $SHELL_CONFIG?"; then
            echo "" >> "$SHELL_CONFIG"
            echo "# Added by zshellcheck installer" >> "$SHELL_CONFIG"
            echo "$EXPORT_CMD" >> "$SHELL_CONFIG"
            echo -e "${GREEN}✓ Added to $SHELL_CONFIG.${NC} Please restart your shell or run 'source $SHELL_CONFIG'."
        else
            echo "Please add the following line to your shell configuration manually:"
            echo -e "  ${BLUE}$EXPORT_CMD${NC}"
        fi
    else
        echo "Please add the following line to your shell configuration:"
        echo -e "  ${BLUE}$EXPORT_CMD${NC}"
    fi
fi

# Fpath check for Zsh user install
if [ "$EUID" -ne 0 ] && [[ "$SHELL" == *"zsh"* ]]; then
    # We can't easily check internal fpath of the running shell, so we check if standard user dir is commonly set up
    # or just advise.
    
    FPATH_CMD="fpath+=($ZSH_COMP_DIR)"
    
    echo ""
    echo -e "${BLUE}Zsh Completions:${NC}"
    if [ -n "$SHELL_CONFIG" ]; then
         if ask_yes_no "Would you like to add the completion directory to your fpath in $SHELL_CONFIG?"; then
            echo "" >> "$SHELL_CONFIG"
            echo "# Added by zshellcheck installer" >> "$SHELL_CONFIG"
            echo "$FPATH_CMD" >> "$SHELL_CONFIG"
            # Note: This needs to be before compinit usually.
            echo -e "${GREEN}✓ Added to $SHELL_CONFIG.${NC}"
            echo -e "${YELLOW}Note: Ensure this line appears BEFORE 'autoload -Uz compinit && compinit' in your config.${NC}"
        else
            echo "Ensure $ZSH_COMP_DIR is in your \$fpath."
            echo "Add this to ~/.zshrc before 'compinit':"
            echo -e "  ${BLUE}$FPATH_CMD${NC}"
        fi
    else
        echo "Ensure $ZSH_COMP_DIR is in your \$fpath."
        echo "Add this to ~/.zshrc before 'compinit':"
        echo -e "  ${BLUE}$FPATH_CMD${NC}"
    fi
fi

echo ""
echo "Run 'zshellcheck --help' to get started."