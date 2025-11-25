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

# --- ZC1042: "$@" over "$*" ---
run_test 'for arg in "$*"; do printf "%s\n" "$arg"; done' "ZC1042" "ZC1042: Quoted dollar star"
# run_test 'for arg in $*; do printf "%s\n" "$arg"; done' "ZC1042" "ZC1042: Unquoted dollar star"
run_test 'for arg in "$@"; do printf "%s\n" "$arg"; done' "" "ZC1042: Quoted dollar at (Valid)"

# --- ZC1043: Local variables in functions ---
# run_test 'fn() { var=1; }' "ZC1043" "ZC1043: Global var"
run_test 'fn() { local var=1; }' "" "ZC1043: Local var"
run_test 'var=1' "" "ZC1043: Global scope (Valid)"

# --- ZC1044: Unchecked cd ---
run_test 'cd /tmp' "ZC1044" "ZC1044: Unchecked cd"
run_test 'cd /tmp || exit' "" "ZC1044: cd || exit"
run_test 'cd /tmp || return' "" "ZC1044: cd || return"
run_test 'if cd /tmp; then printf "ok\n"; fi' "" "ZC1044: if cd"
run_test 'if cd /tmp; echo ok; then echo ok; fi' "ZC1044" "ZC1044: if cd; echo"
run_test '! cd /tmp' "" "ZC1044: ! cd"
run_test 'cd /tmp && echo ok' "ZC1044" "ZC1044: cd && echo (Unsafe)"
run_test '( cd /tmp )' "ZC1044" "ZC1044: Subshell cd unchecked"
run_test 'cd /tmp || printf "fail\n"' "" "ZC1044: cd || echo (Checked)"

# --- ZC1045: Masked return values ---
run_test 'local x=$(cmd)' "ZC1045" "ZC1045: local x=\$(cmd)"
run_test 'typeset y=$(cmd)' "ZC1045" "ZC1045: typeset y=\$(cmd)"
run_test 'local x="foo $(cmd)"' "ZC1045" "ZC1045: local x=\"... \$(cmd)\""
run_test 'fn() { local x; x=$(cmd); }' "" "ZC1045: Split declaration (Valid)"
run_test 'export x=$(cmd)' "ZC1067" "ZC1045: export (Caught by ZC1067)"

# --- ZC1046: Avoid eval ---
run_test 'eval "ls -l"' "ZC1046" "ZC1046: eval"
run_test 'builtin eval "ls -l"' "ZC1046" "ZC1046: builtin eval"
run_test 'command eval "ls -l"' "ZC1046" "ZC1046: command eval"
run_test 'printf "eval\n"' "" "ZC1046: echo word eval (Valid)"

# --- ZC1047: Avoid sudo ---
run_test 'sudo ls' "ZC1047" "ZC1047: sudo"
run_test 'printf "sudo\n"' "" "ZC1047: echo sudo (Valid)"

# --- ZC1048: Relative source ---
run_test 'source ./lib.zsh' "ZC1048" "ZC1048: source ./"
run_test '. ../lib.zsh' "ZC1048" "ZC1048: . ../"
run_test 'source "${0:a:h}/lib.zsh"' "" "ZC1048: Absolute source (Valid)"
run_test 'source /etc/profile' "" "ZC1048: Absolute path (Valid)"

# --- ZC1049: Aliases ---
run_test 'alias foo="echo bar"' "ZC1049" "ZC1049: alias definition"
run_test 'alias -g G="| grep"' "ZC1049" "ZC1049: global alias"
run_test 'unalias foo' "" "ZC1049: unalias (Valid)"
run_test 'function foo() { printf "bar\n"; }' "" "ZC1049: function (Valid)"

# --- ZC1050: Iterating over ls ---
run_test 'for f in $(ls *.txt(N) ); do printf "%s\n" "$f"; done' "ZC1050" "ZC1050: for in \$(ls)"
run_test 'for f in `ls *.txt(N)`; do printf "%s\n" "$f"; done' "ZC1050" "ZC1050: for in \`ls\`"
run_test 'for f in *.txt(N); do printf "%s\n" "$f"; done' "" "ZC1050: for in glob (Valid)"
run_test 'for f in $(find .); do printf "%s\n" "$f"; done' "" "ZC1050: for in find (Valid - specific to ls)"

# --- ZC1051: Unquoted rm ---
# run_test 'rm $var' "ZC1051" "ZC1051: rm variable"
# run_test 'rm "$var"' "" "ZC1051: rm \"$var\" (Valid)"
run_test 'rm ${var}' "ZC1051" "ZC1051: rm braces"
run_test 'rm *' "" "ZC1051: rm * (Valid glob)"

# --- ZC1052: sed -i ---
run_test 'sed -i "s/foo/bar/" file' "ZC1052" "ZC1052: sed -i"
run_test 'sed -e "s/foo/bar/" file' "" "ZC1052: sed -e (Valid)"
run_test 'sed "-i" "s/foo/bar/" file' "ZC1052" "ZC1052: sed \"-i\""

# --- ZC1053: Silence grep ---
run_test 'if grep foo file; then :; fi' "ZC1053" "ZC1053: if grep"
run_test 'while grep foo file; do :; done' "ZC1053" "ZC1053: while grep"
run_test 'if grep -q foo file; then :; fi' "" "ZC1053: grep -q (Valid)"
# run_test 'if grep foo file > /dev/null; then :; fi' "" "ZC1053: grep > /dev/null (Valid)"
run_test 'if grep foo file | wc -l; then :; fi' "" "ZC1053: grep piped (Valid)"

# --- ZC1054: POSIX classes ---
# run_test 'ls [a-z]*' "ZC1054" "ZC1054: glob [a-z]"
# run_test 'ls [[:lower:]]*' "" "ZC1054: glob [[:lower:]] (Valid)"
# run_test '[[ $v =~ [0-9] ]]' "ZC1054" "ZC1054: regex [0-9]"
run_test 'tr "[A-Z]" "[a-z]"' "ZC1054" "ZC1054: tr ranges"
run_test 'tr "[[:upper:]]" "[[:lower:]]"' "" "ZC1054: tr POSIX (Valid)"

# --- ZC1055: Null checks ---
run_test '[[ $var == "" ]]' "ZC1055" "ZC1055: == empty"
run_test '[[ "" != $var ]]' "ZC1055" "ZC1055: != empty"
run_test '[[ -z $var ]]' "" "ZC1055: -z (Valid)"
run_test '[[ -n $var ]]' "" "ZC1055: -n (Valid)"

# --- ZC1056: Arithmetic statement ---
# run_test '$(( i++ ))' "ZC1056" "ZC1056: \$\(\( i++ \)\)"
run_test '(( i++ ))' "" "ZC1056: (( i++ )) (Valid)"
# run_test '$(( 1+1 ))' "ZC1056" "ZC1056: \$\(\( 1+1 \)\)"
run_test '$(ls)' "" "ZC1056: \$(ls) (Valid)"
run_test 'val=$(( 1+1 ))' "" "ZC1056: Assignment (Valid)"

# --- ZC1057: ls assignment ---
run_test 'files=$(ls)' "ZC1057" "ZC1057: files=\$(ls)"
run_test 'files=`ls *.txt`' "ZC1057" "ZC1057: files=\`ls\`"
run_test 'local files=$(ls)' "ZC1057" "ZC1057: local files=\$(ls)"
# run_test 'files=(*)' "" "ZC1057: files=(*) (Valid)"

# --- ZC1058: sudo redirect ---
run_test 'sudo echo "foo" > /etc/file' "ZC1058" "ZC1058: sudo > file"
run_test 'sudo echo "foo" >> /etc/file' "ZC1058" "ZC1058: sudo >> file"
run_test 'printf "foo\n" | sudo tee /etc/file' "ZC1047" "ZC1058: sudo tee (Valid - ZC1047 expected)"
run_test 'sudo ls < /input' "ZC1047" "ZC1058: sudo < input (Valid - ZC1047 expected)"

# --- ZC1059: Unsafe rm variable ---
# run_test 'rm $VAR' "ZC1059" "ZC1059: rm \$VAR (Unsafe)"
run_test 'rm "$VAR"' "ZC1059" "ZC1059: rm \"\$VAR\" (Unsafe)"
run_test 'rm ${VAR}' "ZC1059" "ZC1059: rm \\${VAR}\\ (Unsafe)"
run_test 'rm "${VAR}"' "ZC1059" "ZC1059: rm \"\\${VAR}\" (Unsafe)"
# run_test 'rm /tmp/$VAR' "" "ZC1059: rm path (Valid)"

# --- ZC1060: ps | grep ---
run_test 'ps ax | grep foo' "ZC1060" "ZC1060: ps | grep"
run_test 'ps ax | grep "[f]oo"' "" "ZC1060: ps | grep [] (Valid)"
# run_test 'ps | grep -v grep' "ZC1060" "ZC1060: ps | grep -v (Warn: incomplete chain logic, still risky pattern)"

# --- ZC1061: seq vs range ---
run_test 'for i in $(seq 1 10); do :; done' "ZC1061" "ZC1061: for seq"
run_test 'for i in {1..10}; do :; done' "" "ZC1061: for range (Valid)"
run_test 'seq 5' "ZC1061" "ZC1061: seq command"

# --- ZC1062: egrep ---
run_test 'egrep "pattern" file' "ZC1062" "ZC1062: egrep"
run_test 'grep -E "pattern" file' "" "ZC1062: grep -E (Valid)"

# --- ZC1063: fgrep ---
run_test 'fgrep "pattern" file' "ZC1063" "ZC1063: fgrep"
run_test 'grep -F "pattern" file' "" "ZC1063: grep -F (Valid)"

# --- ZC1064: type vs command -v ---
run_test 'type ls' "ZC1064" "ZC1064: type"
run_test 'command -v ls' "" "ZC1064: command -v (Valid)"

# --- ZC1065: Spaces in test ---
run_test '[foo]' "ZC1065" "ZC1065: [foo]"
run_test '[[foo]]' "ZC1065" "ZC1065: [[foo]]"
# run_test '[ foo ]' "" "ZC1065: [ foo ] (Valid)"
run_test '[[ foo ]]' "" "ZC1065: [[ foo ]] (Valid)"

# --- ZC1066: cat iteration ---
run_test 'for l in $(cat file); do :; done' "ZC1066" "ZC1066: for in \$(cat)"
run_test 'for l in `cat file`; do :; done' "ZC1066" "ZC1066: for in \`cat\`"
# run_test 'while IFS= read -r l; do :; done < file' "" "ZC1066: while read (Valid)"

# --- ZC1067: export masking return value ---
run_test 'export var=$(cmd)' "ZC1067" "ZC1067: export var=\$(cmd)"
run_test 'export var=`cmd`' "ZC1067" "ZC1067: export var=\`cmd\`"
run_test 'var=$(cmd); export var' "" "ZC1067: var=\$(cmd); export var (Valid)"
run_test 'export var="value"' "" "ZC1067: export var=\"value\" (Valid)"

# --- ZC1068: Hook functions ---
run_test 'precmd() { :; }' "ZC1068" "ZC1068: precmd()"
run_test 'function chpwd() { :; }' "ZC1068" "ZC1068: function chpwd()"
run_test 'autoload -Uz add-zsh-hook; add-zsh-hook precmd my_func' "" "ZC1068: add-zsh-hook (Valid)"
run_test 'my_func() { :; }' "" "ZC1068: normal function (Valid)"

# --- ZC1069: local scope ---
run_test 'local x=1' "ZC1069" "ZC1069: local global"
run_test 'typeset x=1' "" "ZC1069: typeset global (Valid)"
run_test 'fn() { local x=1; }' "" "ZC1069: local in func (Valid)"
run_test 'if true; then local x=1; fi' "ZC1069" "ZC1069: local in if (global)"
run_test '( local x=1 )' "ZC1069" "ZC1069: local in subshell (global)"

# --- ZC1070: Infinite recursion ---
run_test 'cd() { cd "$@"; }' "ZC1070" "ZC1070: recursive cd"
run_test 'ls() { ls -G; }' "ZC1070" "ZC1070: recursive ls"
run_test 'cd() { builtin cd "$@"; }' "" "ZC1070: builtin cd (Valid)"
run_test 'myfunc() { myfunc; }' "" "ZC1070: recursive custom (Valid/Ignored)"

# --- ZC1072: grep | awk ---
run_test 'grep pattern file | awk "{print $1}"' "ZC1072" "ZC1072: grep | awk"
run_test 'grep -i pattern file | gawk "{print}"' "ZC1072" "ZC1072: grep -i | gawk"
run_test 'grep -r pattern . | awk "{print}"' "" "ZC1072: grep -r | awk (Valid)"
run_test 'awk "/pattern/ {print}" file' "" "ZC1072: awk only (Valid)"

# --- ZC1073: Unnecessary $ in arithmetic ---
run_test '(( $x > 5 ))' "ZC1073" "ZC1073: (( \$x ))"
run_test '(( x > 5 ))' "" "ZC1073: (( x )) (Valid)"
run_test '(( $# > 0 ))' "" "ZC1073: (( \$# )) (Valid)"

# --- ZC1083: Brace expansion variables ---
run_test 'echo {1..$n}' "ZC1083" "ZC1083: variable range end"
run_test 'echo {$n..10}' "ZC1083" "ZC1083: variable range start"
run_test 'printf "%s\n" {1..10}' "" "ZC1083: valid range"
run_test 'printf "%s\n" {a,b,$c}' "" "ZC1083: valid list expansion"
run_test 'echo "{1..$n}"' "ZC1083" "ZC1083: quoted variable range"

# --- ZC1084: find unquoted glob ---
run_test 'find . -name *.txt' "ZC1084" "ZC1084: *.txt"
run_test 'find . -name "*.txt"' "" "ZC1084: \"*.txt\" (Valid)"
run_test 'find . -name [a-z]' "ZC1084" "ZC1084: [a-z]"
run_test 'find . -name "[a-z]"' "" "ZC1084: \"[a-z]\" (Valid)"
run_test 'find . -name \*.txt' "" "ZC1084: \\*.txt (Valid)"

# --- ZC1085: for loop variable expansion ---
run_test 'for i in $items; do :; done' "ZC1085" "ZC1085: \$items"
run_test 'for i in "$items"; do :; done' "" "ZC1085: \"\$items\" (Valid)"
run_test 'for i in ${arr[@]}; do :; done' "ZC1085" "ZC1085: \${arr[@]}"
run_test 'for i in "${arr[@]}"; do :; done' "" "ZC1085: \"\${arr[@]}\" (Valid)"
run_test 'for i in *.txt; do :; done' "" "ZC1085: glob (Valid)"

# --- Summary ---
echo "------------------------------------------------"
if [[ $FAILURES -eq 0 ]]; then
    echo "All $TOTAL tests passed."
    exit 0
else
    echo "$FAILURES out of $TOTAL tests failed."
    exit 1
fi