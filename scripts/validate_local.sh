#!/bin/bash
set -e

ZSHELLCHECK_BIN="./zshellcheck"
TARGET_DIR="/home/rtx/.config/zsh"

if [ ! -f "$ZSHELLCHECK_BIN" ]; then
echo "Building zshellcheck..."
go build -o zshellcheck ./cmd/zshellcheck
if [ $? -ne 0 ]; then
    echo "Build failed."
    exit 1
fi
fi

echo "Running zshellcheck against $TARGET_DIR..."

# Find zsh files and run zshellcheck
# Use fd if available, else find
if command -v fd >/dev/null; then
    FILES=$(fd -e zsh . "$TARGET_DIR")
else
    FILES=$(find "$TARGET_DIR" -name "*.zsh")
fi

PASS=0
FAIL=0

for f in $FILES; do
    # echo "Checking $f..."
    OUTPUT=$($ZSHELLCHECK_BIN "$f" 2>&1 || true)
    EXIT_CODE=$?
    
    # zshellcheck might not return non-zero on violations yet?
    # Check if output is empty.
    
    if [ -n "$OUTPUT" ]; then
        echo "FAIL: $f"
        echo "$OUTPUT"
        FAIL=$((FAIL+1))
    else
        # echo "PASS: $f"
        PASS=$((PASS+1))
    fi
done

echo "------------------------------------------------"
echo "Summary: PASS=$PASS, FAIL=$FAIL"

if [ $FAIL -gt 0 ]; then
    exit 1
fi
