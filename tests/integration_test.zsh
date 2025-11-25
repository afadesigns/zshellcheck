#!/usr/bin/env zsh

# integration_test.zsh
# Runs zshellcheck against synthetic test cases to verify correctness.

set -e

# Ensure we are in root
cd "${0:a:h}/.."

# Build binary
echo "Building zshellcheck..."
go build -o bin/zshellcheck ./cmd/zshellcheck || { echo "Build failed"; exit 1; }

FAILURES=0
TOTAL=0

run_test() {
    local input="$1"
    local expected_kata="$2"
    local name="$3"
    
    TOTAL=$((TOTAL + 1))
    
    local tmpfile=$(mktemp)
    mv "$tmpfile" "${tmpfile}.zsh"
    tmpfile="${tmpfile}.zsh"
    
    echo "$input" > "$tmpfile"
    
    # Run zshellcheck
    local output=$(./bin/zshellcheck -format json "$tmpfile" 2>&1)
    rm "$tmpfile"
    
    if [[ -z "$expected_kata" ]]; then
        # Expect Pass
        if [[ "$output" == "[]" ]] || [[ -z "$output" ]]; then
            echo "PASS: $name"
        else
            echo "FAIL: $name"
            echo "  Expected NO violations, got:"
            echo "  $output"
            FAILURES=$((FAILURES + 1))
        fi
    else
        # Expect Violation
        if [[ "$output" == *"$expected_kata"* ]]; then
            echo "PASS: $name"
        else
            echo "FAIL: $name"
            echo "  Expected violation $expected_kata, got:"
            echo "  $output"
            FAILURES=$((FAILURES + 1))
        fi
    fi
}

echo "Running integration tests..."

# --- ZC1001: Array Access ---
run_test 'val=${arr[1]}' "" "ZC1001: Valid array access"
run_test 'val=$arr[1]' "ZC1001" "ZC1001: Invalid array access"

# --- ZC1002: Backticks ---
run_test 'val=$(cmd)' "" "ZC1002: Valid command subst"
run_test 'val=`cmd`' "ZC1002" "ZC1002: Backticks"

# --- ZC1003: Arithmetic Comparisons ---
run_test 'if (( val > 0 )); then; fi' "" "ZC1003: Valid ((...))"
run_test 'if [ $val -gt 0 ]; then; fi' "ZC1003" "ZC1003: [ -gt ]"

# --- ZC1006: [[ vs test ---
run_test 'if [[ -f file ]]; then; fi' "" "ZC1006: Valid [["
run_test 'if test -f file; then; fi' "ZC1006" "ZC1006: test command"

# --- Operator Precedence / Parser Regression Checks ---
run_test 'return ( (5 + 5) * 2 )' "" "Parser: Precedence 1"
run_test 'a + add( (b * c) ) + d' "" "Parser: Precedence 2 (nested call)"
run_test 'if ((1 < 2)); then return true; fi' "" "Parser: If Statement"

# --- ZC1039: Dangerous rm ---
run_test 'rm /' "ZC1039" "ZC1039: rm /"
run_test 'rm /tmp' "" "ZC1039: rm /tmp (Valid)"

# --- ZC1040: Nullglob in loops ---
run_test 'for i in *.txt; do printf "%s\n" "$i"; done' "ZC1040" "ZC1040: Missing (N)"
run_test 'for i in *.txt(N); do printf "%s\n" "$i"; done' "" "ZC1040: With (N)"
run_test 'for i in *; do printf "%s\n" "$i"; done' "ZC1040" "ZC1040: * missing (N)"

# --- ZC1041: printf format string ---
run_test 'printf "$var"' "ZC1041" "ZC1041: Variable format string"
run_test 'printf "Hello %s" "$var"' "" "ZC1041: Static format string"
run_test 'printf $fmt "arg"' "ZC1041" "ZC1041: Identifier format string"

# --- ZC1042: "\