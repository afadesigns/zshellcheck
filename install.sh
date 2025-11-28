#!/usr/bin/env bash
set -euo pipefail

echo "Installing zshellcheck..."

if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.25 or later."
    exit 1
fi

echo "Building and installing..."
if go install -v ./cmd/zshellcheck; then
    echo "Successfully installed zshellcheck."
    
    # Check if installed binary is in PATH
    GOPATH="$(go env GOPATH)"
    if [ -z "$GOPATH" ]; then
        GOPATH="$HOME/go"
    fi
    GOBIN="$GOPATH/bin"
    
    if [[ ":$PATH:" != ":$GOBIN:"* ]]; then
        echo ""
        echo "WARNING: $GOBIN is not in your PATH."
        echo "You can run zshellcheck via:"
        echo "  $GOBIN/zshellcheck"
        echo ""
        echo "Or add it to your PATH by adding this to your shell config (e.g. ~/.zshrc):"
        echo "  export PATH=\"$PATH:$GOBIN\""
    else
        echo "You can now run 'zshellcheck' from your terminal."
    fi
else
    echo "Installation failed."
    exit 1
fi
