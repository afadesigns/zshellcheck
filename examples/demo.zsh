#!/bin/zsh

# This script contains intentional violations to demonstrate ZShellCheck's capabilities.
# Run: zshellcheck examples/demo.zsh

# ZC1001: Invalid array access
my_array=(alpha beta gamma)
echo "Item: $my_array[1]"

# ZC1002: Use of backticks
current_date=`date`

# ZC1003: Use of [ ] for arithmetic
if [ $1 -gt 10 ]; then
    echo "Big number"
fi

# ZC1005: Use of 'which'
which git

# ZC1006: Use of test command
if test -f "/tmp/foo"; then
    echo "Found"
fi

# ZC1007: chmod 777
chmod 777 /tmp/testfile

# ZC1013: Use of let
let "x = x + 1"

# ZC1011: Git plumbing
git checkout main

# ZC1030: Use of echo (portability)
echo "Processing..."
