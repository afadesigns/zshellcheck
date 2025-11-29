#!/usr/bin/env bash
set -euo pipefail

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
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

# Determine install location
INSTALL_DIR=""
if [ "$EUID" -eq 0 ]; then
    INSTALL_DIR="/usr/local/bin"
else
    # Prefer ~/.local/bin for non-root users
    INSTALL_DIR="$HOME/.local/bin"
fi

echo "Installing to $INSTALL_DIR..."

# Ensure directory exists
if [ ! -d "$INSTALL_DIR" ]; then
    echo "Creating directory $INSTALL_DIR..."
    mkdir -p "$INSTALL_DIR"
fi

# Move binary
if mv zshellcheck "$INSTALL_DIR/zshellcheck"; then
    echo -e "${GREEN}Successfully installed zshellcheck to $INSTALL_DIR/zshellcheck${NC}"
else
    echo -e "${RED}Failed to move binary to $INSTALL_DIR.${NC}"
    if [ "$EUID" -ne 0 ]; then
        echo "Permission denied. You may need to run with sudo to install to a system directory, or check permissions on ~/.local/bin."
    fi
    rm zshellcheck # Cleanup
    exit 1
fi

# Path check
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo ""
    echo -e "${YELLOW}WARNING: $INSTALL_DIR is not in your PATH.${NC}"
    echo "To run zshellcheck, add this to your shell configuration (e.g., ~/.zshrc or ~/.bashrc):"
    echo ""
    echo "  export PATH=\"$PATH:$INSTALL_DIR\""
    echo ""
    echo "Then restart your shell or run 'source <config_file>'."
else
    echo ""
    echo -e "${GREEN}Installation complete! You can now run 'zshellcheck' from your terminal.${NC}"
fi