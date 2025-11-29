#!/usr/bin/env bash
set -euo pipefail

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}Installing zshellcheck...${NC}"

# Check for Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed. Please install Go 1.25 or later.${NC}"
    exit 1
fi

# Build
echo "Building binary..."
if ! go build -o zshellcheck cmd/zshellcheck/main.go; then
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

# --- Final Checks ---
echo ""
echo -e "${GREEN}Installation complete!${NC}"

# Path check
if [[ ":$PATH:" != ":$BIN_DIR:"* ]]; then
    echo ""
    echo -e "${YELLOW}WARNING: $BIN_DIR is not in your PATH.${NC}"
    echo "Add this to your shell configuration (e.g., ~/.zshrc or ~/.bashrc):"
    echo -e "  ${BLUE}export PATH=\"$PATH:$BIN_DIR\"${NC}"
fi

# Fpath check for Zsh user install
if [ "$EUID" -ne 0 ] && [[ "$SHELL" == *"zsh"* ]]; then
    # Ideally we'd check fpath but that requires running zsh.
    # Just giving a friendly reminder is safer.
    echo ""
    echo -e "${BLUE}Note for Zsh users:${NC}"
    echo "Ensure $ZSH_COMP_DIR is in your \$fpath to enable completions."
    echo "Add this to ~/.zshrc before 'compinit':"
    echo -e "  ${BLUE}fpath+=($ZSH_COMP_DIR)${NC}"
fi

echo ""
echo "Run 'zshellcheck --help' to get started."
