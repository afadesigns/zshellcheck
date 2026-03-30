# ZShellCheck Katas

Comprehensive list of all 166 implemented checks, migrated from the Wiki.

## Table of Contents

- [ZC1001: Use ${} for array element access](#zc1001)
- [ZC1002: Use $(...) instead of backticks](#zc1002)
- [ZC1003: Use `((...))` for arithmetic comparisons instead of `[` or `test`](#zc1003)
- [ZC1004: Use `return` instead of `exit` in functions](#zc1004)
- [ZC1005: Use `whence` instead of `which`](#zc1005)
- [ZC1006: Prefer `[[` over `test` for tests](#zc1006)
- [ZC1007: Avoid using `chmod 777`](#zc1007)
- [ZC1008: Use `$(())` for arithmetic operations](#zc1008)
- [ZC1009: Use `((...))` for C-style arithmetic](#zc1009)
- [ZC1010: Use `[[ ... ]]` instead of `[ ... ]`](#zc1010)
- [ZC1011: Use `git` porcelain commands instead of plumbing commands](#zc1011)
- [ZC1012: Use `read -r` to prevent backslash escaping](#zc1012)
- [ZC1013: Use `((...))` for arithmetic operations instead of `let`](#zc1013)
- [ZC1014: Use `git switch` or `git restore` instead of `git checkout`](#zc1014)
- [ZC1015: Use `$(...)` for command substitution instead of backticks](#zc1015)
- [ZC1016: Use `read -s` when reading sensitive information](#zc1016)
- [ZC1017: Use `print -r` to print strings literally](#zc1017)
- [ZC1018: Use `((...))` for C-style arithmetic instead of `expr`](#zc1018)
- [ZC1019: Use `whence` instead of `which`](#zc1019)
- [ZC1020: Use `[[ ... ]]` for tests instead of `test`](#zc1020)
- [ZC1021: Use symbolic permissions with `chmod` instead of octal](#zc1021)
- [ZC1022: Use `$((...))` for arithmetic expansion](#zc1022)
- [ZC1023: Use `$((...))` for arithmetic expansion](#zc1023)
- [ZC1024: Use `$((...))` for arithmetic expansion](#zc1024)
- [ZC1025: Use `$((...))` for arithmetic expansion](#zc1025)
- [ZC1026: Use `$((...))` for arithmetic expansion](#zc1026)
- [ZC1027: Use `$((...))` for arithmetic expansion](#zc1027)
- [ZC1028: Use `$((...))` for arithmetic expansion](#zc1028)
- [ZC1029: Use `$((...))` for arithmetic expansion](#zc1029)
- [ZC1030: Use `printf` instead of `echo`](#zc1030)
- [ZC1031: Use `#!/usr/bin/env zsh` for portability](#zc1031)
- [ZC1032: Use `((...))` for C-style incrementing](#zc1032)
- [ZC1033: Use `$((...))` for arithmetic expansion](#zc1033)
- [ZC1034: Use `command -v` instead of `which`](#zc1034)
- [ZC1035: Use `$((...))` for arithmetic expansion](#zc1035)
- [ZC1036: Prefer `[[ ... ]]` over `test` command](#zc1036)
- [ZC1037: Use `print -r --` for variable expansion](#zc1037)
- [ZC1038: Avoid useless use of `cat`](#zc1038)
- [ZC1039: Avoid `rm` with root path](#zc1039)
- [ZC1040: Use `(N)` nullglob qualifier for globs in loops](#zc1040)
- [ZC1041: Do not use variables in `printf` format string](#zc1041)
- [ZC1042: Use "$@" to iterate over arguments](#zc1042)
- [ZC1043: Use `local` for variables in functions](#zc1043)
- [ZC1044: Check for unchecked `cd` commands](#zc1044)
- [ZC1045: Declare and assign separately to avoid masking return values](#zc1045)
- [ZC1046: Avoid `eval`](#zc1046)
- [ZC1047: Avoid `sudo` in scripts](#zc1047)
- [ZC1048: Avoid `source` with relative paths](#zc1048)
- [ZC1049: Prefer functions over aliases](#zc1049)
- [ZC1050: Avoid iterating over `ls` output](#zc1050)
- [ZC1051: Quote variables in `rm` to avoid globbing](#zc1051)
- [ZC1052: Avoid `sed -i` for portability](#zc1052)
- [ZC1053: Silence `grep` output in conditions](#zc1053)
- [ZC1054: Use POSIX classes in regex/glob](#zc1054)
- [ZC1055: Use `[[ -n/-z ]]` for empty string checks](#zc1055)
- [ZC1056: Avoid `$((...))` as a statement](#zc1056)
- [ZC1057: Avoid `ls` in assignments](#zc1057)
- [ZC1058: Avoid `sudo` with redirection](#zc1058)
- [ZC1059: Use `${var:?}` for `rm` arguments](#zc1059)
- [ZC1060: Avoid `ps | grep` without exclusion](#zc1060)
- [ZC1061: Prefer `{start..end}` over `seq`](#zc1061)
- [ZC1062: Prefer `grep -E` over `egrep`](#zc1062)
- [ZC1063: Prefer `grep -F` over `fgrep`](#zc1063)
- [ZC1064: Prefer `command -v` over `type`](#zc1064)
- [ZC1065: Ensure spaces around `[` and `[[`](#zc1065)
- [ZC1066: Avoid iterating over `cat` output](#zc1066)
- [ZC1067: Separate `export` and assignment to avoid masking return codes](#zc1067)
- [ZC1068: Use `add-zsh-hook` instead of defining hook functions directly](#zc1068)
- [ZC1069: Avoid `local` outside of functions](#zc1069)
- [ZC1070: Use `builtin` or `command` to avoid infinite recursion in wrapper functions](#zc1070)
- [ZC1071: Use `+=` for appending to arrays](#zc1071)
- [ZC1072: Use `awk` instead of `grep | awk`](#zc1072)
- [ZC1073: Unnecessary use of `$` in arithmetic expressions](#zc1073)
- [ZC1074: Prefer modifiers :h/:t over dirname/basename](#zc1074)
- [ZC1075: Quote variable expansions to prevent globbing](#zc1075)
- [ZC1076: Use `autoload -Uz` for lazy loading](#zc1076)
- [ZC1077: Prefer `${var:u/l}` over `tr` for case conversion](#zc1077)
- [ZC1078: Quote `$@` and `$*` when passing arguments](#zc1078)
- [ZC1079: Quote RHS of `==` in `[[ ... ]]` to prevent pattern matching](#zc1079)
- [ZC1080: Use `(N)` nullglob qualifier for globs in loops](#zc1080)
- [ZC1081: Use `${#var}` to get string length instead of `wc -c`](#zc1081)
- [ZC1082: Prefer `${var//old/new}` over `sed` for simple replacements](#zc1082)
- [ZC1083: Brace expansion limits cannot be variables](#zc1083)
- [ZC1084: Quote globs in `find` commands](#zc1084)
- [ZC1085: Quote variable expansions in `for` loops](#zc1085)
- [ZC1086: Prefer `func() { ... }` over `function func { ... }`](#zc1086)
- [ZC1087: Output redirection overwrites input file](#zc1087)
- [ZC1088: Subshell isolates state changes](#zc1088)
- [ZC1089: Redirection order matters (`2>&1 > file`)](#zc1089)
- [ZC1090: Quoted regex pattern in `=~`](#zc1090)
- [ZC1091: Use `((...))` for arithmetic comparisons in `[[...]]`](#zc1091)
- [ZC1092: Arrays cannot be exported](#zc1092)
- [ZC1093: Use `IFS=` with `read` in loops](#zc1093)
- [ZC1094: Consider `grep -P` for advanced regex](#zc1094)
- [ZC1095: Quote here-string content](#zc1095)
- [ZC1096: Avoid using `test -e` or `[ -e ... ]` for file existence checks](#zc1096)
- [ZC1097: Declare loop variables as `local` in functions](#zc1097)
- [ZC1098: Use `(q)` flag for quoting variables in `eval`](#zc1098)
- [ZC1099: Use `(f)` flag to split lines instead of `while read`](#zc1099)
- [ZC1100: Use parameter expansion instead of `dirname`/`basename`](#zc1100)
- [ZC1101: Use `$(( ))` instead of `bc` for simple arithmetic](#zc1101)
- [ZC1102: Redirecting output of `sudo` does not work as expected](#zc1102)
- [ZC1103: Suggest `path` array instead of `$PATH` string manipulation](#zc1103)
- [ZC1104: Suggest `path` array instead of `export PATH` string manipulation](#zc1104)
- [ZC1105: Avoid nested arithmetic expansions for clarity](#zc1105)
- [ZC1106: Avoid `set -x` in production scripts](#zc1106)
- [ZC1107: Use `(( ... ))` for arithmetic conditions](#zc1107)
- [ZC1108: Use Zsh case conversion instead of `tr`](#zc1108)
- [ZC1109: Use parameter expansion instead of `cut` for field extraction](#zc1109)
- [ZC1110: Use Zsh subscripts instead of `head -1` or `tail -1`](#zc1110)
- [ZC1111: Avoid `xargs` for simple command invocation](#zc1111)
- [ZC1112: Avoid `grep -c` -- use Zsh pattern matching for counting](#zc1112)
- [ZC1113: Use `${var:A}` instead of `realpath` or `readlink -f`](#zc1113)
- [ZC1114: Consider Zsh `=(...)` for temporary files](#zc1114)
- [ZC1115: Use Zsh string manipulation instead of `rev`](#zc1115)
- [ZC1116: Use Zsh multios instead of `tee`](#zc1116)
- [ZC1117: Use `&!` or `disown` instead of `nohup`](#zc1117)
- [ZC1118: Use `print -rn` instead of `echo -n`](#zc1118)
- [ZC1119: Use `$EPOCHSECONDS` instead of `date +%s`](#zc1119)
- [ZC1120: Use `$PWD` instead of `pwd`](#zc1120)
- [ZC1121: Use `$HOST` instead of `hostname`](#zc1121)
- [ZC1122: Use `$USER` instead of `whoami`](#zc1122)
- [ZC1123: Use `$OSTYPE` instead of `uname`](#zc1123)
- [ZC1124: Use `: > file` instead of `cat /dev/null > file` to truncate](#zc1124)
- [ZC1125: Avoid `echo | grep` for string matching](#zc1125)
- [ZC1126: Use `sort -u` instead of `sort | uniq`](#zc1126)
- [ZC1127: Avoid `ls` for counting files](#zc1127)
- [ZC1128: Use `> file` instead of `touch file` for creation](#zc1128)
- [ZC1129: Use Zsh `stat` module instead of `wc -c` for file size](#zc1129)
- [ZC1131: Avoid `cat file | while read` -- use redirection](#zc1131)
- [ZC1132: Use Zsh pattern extraction instead of `grep -o`](#zc1132)
- [ZC1133: Avoid `kill -9` -- use `kill` first, then escalate](#zc1133)
- [ZC1134: Avoid `sleep` in tight loops](#zc1134)
- [ZC1135: Avoid `env VAR=val cmd` -- use inline assignment](#zc1135)
- [ZC1136: Avoid `rm -rf` without safeguard](#zc1136)
- [ZC1137: Avoid hardcoded `/tmp` paths](#zc1137)
- [ZC1139: Avoid `source` with URL -- use local files](#zc1139)
- [ZC1140: Use `command -v` instead of `hash` for command existence](#zc1140)
- [ZC1141: Avoid `curl | sh` pattern](#zc1141)
- [ZC1142: Avoid chained `grep | grep` -- combine patterns](#zc1142)
- [ZC1143: Avoid `set -e` -- use explicit error handling](#zc1143)
- [ZC1144: Avoid `trap` with signal numbers -- use names](#zc1144)
- [ZC1145: Avoid `tr -d` for character deletion -- use parameter expansion](#zc1145)
- [ZC1146: Avoid `cat file | awk` -- pass file to awk directly](#zc1146)
- [ZC1147: Avoid `mkdir` without `-p` for nested paths](#zc1147)
- [ZC1148: Use `compdef` instead of `compctl` for completions](#zc1148)
- [ZC1149: Avoid `echo` for error messages -- use `>&2`](#zc1149)
- [ZC1151: Avoid `cat -A` -- use `print -v` or `od` for non-printable characters](#zc1151)
- [ZC1152: Use Zsh PCRE module instead of `grep -P`](#zc1152)
- [ZC1153: Use `cmp -s` instead of `diff` for equality check](#zc1153)
- [ZC1154: Use `find -exec {} +` instead of `find -exec {} \;`](#zc1154)
- [ZC1155: Use `whence -a` instead of `which -a`](#zc1155)
- [ZC1156: Avoid `ln` without `-s` for symlinks](#zc1156)
- [ZC1157: Avoid `strings` command -- use Zsh expansion](#zc1157)
- [ZC1158: Avoid `chown -R` without `--no-dereference`](#zc1158)
- [ZC1159: Avoid `tar` without explicit compression flag](#zc1159)
- [ZC1160: Prefer `curl` over `wget` for portability](#zc1160)
- [ZC1161: Avoid `openssl` for simple hashing -- use Zsh modules](#zc1161)
- [ZC1162: Use `cp -a` instead of `cp -r` to preserve attributes](#zc1162)
- [ZC1163: Use `grep -m 1` instead of `grep | head -1`](#zc1163)
- [ZC1164: Avoid `sed -n 'Np'` -- use Zsh array subscript](#zc1164)
- [ZC1165: Use Zsh parameter expansion for simple `awk` field extraction](#zc1165)
- [ZC1166: Avoid `grep -i` for case-insensitive match -- use `(#i)` glob flag](#zc1166)
- [ZC1167: Avoid `timeout` command -- use Zsh `TMOUT` or `zsh/sched`](#zc1167)
- [ZC1168: Use `${(f)...}` instead of `readarray`/`mapfile`](#zc1168)
- [ZC1169: Avoid `install` for simple copy+chmod -- use `cp` then `chmod`](#zc1169)

---

<div id="zc1001"></div>

<details>
<summary><strong>ZC1001</strong>: Use ${} for array element access <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

In Zsh, accessing array elements using `$array[index]` can sometimes behave unexpectedly or be misinterpreted in complex expansions. The more explicit and safer syntax is `${array[index]}`, which clearly delimits the array access from surrounding text or other expansions. This ensures the correct element is accessed, especially when dealing with nested expansions or when the array name might be followed by characters that could be part of a variable name.

### Bad Example

```zsh
my_array=(alpha beta gamma)
echo $my_array[2]suffix # Might output "beta" then "suffix", or error
```

### Good Example

```zsh
my_array=(alpha beta gamma)
echo "The second element is ${my_array[2]}."
## Expected output: The second element is beta.

## Or, to safely concatenate:
echo "${my_array[2]}suffix" # Clearly outputs "betasuffix"
```

### Configuration

To disable this Kata, add `ZC1001` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1002"></div>

<details>
<summary><strong>ZC1002</strong>: Use $(...) instead of backticks <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Backticks (`` `command` ``) for command substitution are an older, deprecated syntax with several significant disadvantages compared to `$(command)`. The `$(...)` form is the modern, preferred, and **POSIX-standard** way to perform command substitution in Zsh and other compatible shells. Preferring `$(...)` offers substantial benefits:

*   **Readability and Clarity:** `$(...)` clearly delimits the command being substituted, improving visual parsing and understanding of the script's logic.
*   **Arbitrary Nesting:** `$(...)` allows for arbitrary nesting of command substitutions without complex and error-prone backslash escaping. Backticks require cumbersome escaping (e.g., `` `echo \`date\`` ``) for nesting, which quickly becomes unmanageable.
*   **Reduced Ambiguity:** Backticks can sometimes be ambiguous with string literals or glob patterns, leading to unexpected behavior. `$(...)` avoids this ambiguity.
*   **Robustness:** Scripts using `$(...)` are generally more robust and less prone to subtle parsing errors across different shell environments or with complex inputs.

Adopting `$(...)` consistently leads to more readable, robust, and portable Zsh scripts.

### Bad Example

```zsh
## Old-style backticks
file_count=`ls | wc -l`
timestamp=`date +"%Y-%m-%d"`

## Difficult and error-prone nesting with backticks
nested_output=`echo \`date\``
```

### Good Example

```zsh
## Modern and clear command substitution
file_count=$(ls | wc -l)
timestamp=$(date +"%Y-%m-%d")

## Easy and readable nesting with $(...)
nested_output=$(echo $(date))
```

### Configuration

To disable this Kata, add `ZC1002` to the `disabled_katas` list in your `.zshellcheckrc` file.

```

```

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1003"></div>

<details>
<summary><strong>ZC1003</strong>: Use `((...))` for arithmetic comparisons instead of `[` or `test` <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

For pure **integer arithmetic comparisons** in Zsh, `((...))` provides a more natural, C-style syntax and often better performance and robustness than `[` or `test`. The `[[...]]` construct also supports arithmetic evaluation with `((...))`, but for direct integer comparisons, `((...))` is more idiomatic and efficient. Using `[` or `test` for arithmetic comparisons relies on external commands or string evaluation, which can be slower, less type-safe (as it treats numbers as strings), and prone to subtle quoting or parsing issues.

### Bad Example

```zsh
a=10; b=5
if [ "$a" -gt "$b" ]; then echo "a > b"; fi  # Relies on external `test` or `[` command
if test "$a" -le "$b"; then echo "a <= b"; fi # String comparison, potential issues with non-integers or edge cases
```

### Good Example

```zsh
a=10; b=5
if (( a > b )); then echo "a > b"; fi      # Direct integer arithmetic comparison
if (( a <= b )); then echo "a <= b"; fi   # Clearer, more performant, no quoting needed
```

### Configuration

To disable this Kata, add `ZC1003` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>


<div id="zc1004"></div>

<details>
<summary><strong>ZC1004</strong>: Use `return` instead of `exit` in functions <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Using `exit` in a function terminates the entire shell, which is often unintended in interactive sessions or sourced scripts. Use `return` to exit the function.

### Bad Example

```zsh
my_func() {
  if [[ -z $1 ]]: then
    exit 1 # Kills the shell!
  fi
}
```

### Good Example

```zsh
my_func() {
  if [[ -z $1 ]]: then
    return 1
  fi
}
```

### Configuration

To disable this Kata, add `ZC1004` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>


<div id="zc1005"></div>

<details>
<summary><strong>ZC1005</strong>: Use `whence` instead of `which` <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

In Zsh, `whence` is a built-in command that is significantly more powerful and accurate than the external `which` utility. `whence` reports exactly how Zsh will interpret a command, including whether it's an alias, a function, a built-in, or an executable found in your `PATH`. This provides a complete and reliable picture of command resolution within your current Zsh environment.

`which`, on the other hand, is an external utility that only searches your `PATH` for executable files. It will not show aliases, functions, or built-ins, which can lead to confusion or incorrect assumptions about which command will actually be executed by Zsh.

### Bad Example

```zsh
which ls    # Might show /bin/ls, but 'ls' could be an alias like 'ls --color=auto'
which my_zsh_function # Will not find Zsh functions
```

### Good Example

```zsh
whence ls   # Shows if 'ls' is an alias, function, or executable
whence -c ls # Provides more verbose, shell-like output (e.g., 'ls is alias ls --color=auto')
whence my_zsh_function # Correctly identifies Zsh functions
```

### Configuration

To disable this Kata, add `ZC1005` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1006"></div>

<details>
<summary><strong>ZC1006</strong>: Prefer `[[` over `test` for tests <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

In Zsh, `[[...]]` is a more powerful, safer, and idiomatic conditional construct than the traditional `test` command or single brackets `[...]`. The primary advantages of `[[...]]` include:

*   **Intelligent Word Splitting and Globbing:** Unlike `[` or `test`, `[[...]]` does not perform word splitting or pathname expansion (globbing) on unquoted variables, avoiding many common pitfalls and unexpected behaviors. This means you generally don't need to quote variables within `[[...]]` unless you specifically want literal string matching for patterns that might otherwise be interpreted as globs.
*   **Zsh-Specific Features:** It supports Zsh-specific features like regular expression matching (`=~`), globbing with extended patterns (`==`), and compound conditions (`&&`, `||`).
*   **Built-in Efficiency:** `[[...]]` is a keyword, not an external command, making it often faster and more efficient.

Using `test` or `[` requires meticulous quoting of variables to prevent unintended word splitting and pathname expansion, which can lead to bugs, security vulnerabilities, or simply incorrect logic. For arithmetic comparisons, `((...))` is generally preferred (see [[Katas/ZC1000-ZC1099/ZC1003 | ZC1003]]).

### Bad Example

```zsh
my_var="foo bar"
if [ -n $my_var ]; then echo "This might split into two arguments"; fi
## Expected to fail if $my_var contains spaces, or behave unexpectedly if it expands to a file glob.

file="my file.txt"
if [ -f $file ]; then echo "File exists"; fi # Will fail if 'my file.txt' is split

if test "$str1" = "$str2"; then echo "Strings are equal"; fi # Correctly quoted, but `[[...]]` is still preferred
```

### Good Example

```zsh
my_var="foo bar"
if [[ -n $my_var ]]; then echo "Correctly handles spaces and globs"; fi
## No quotes needed for $my_var, `[[...]]` handles it correctly.

file="my file.txt"
if [[ -f $file ]]; then echo "File exists"; fi # Safely handles spaces in file names

if [[ "$str1" = "$str2" ]]; then echo "Strings are equal"; fi

## Zsh-specific regex matching
if [[ "hello world" =~ "hello (.*)" ]]; then echo "Regex match!"; fi
```

### Configuration

To disable this Kata, add `ZC1006` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1007"></div>

<details>
<summary><strong>ZC1007</strong>: Avoid using `chmod 777` <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Setting file or directory permissions to `777` (read, write, execute for owner, group, and others) is a significant security risk. It grants unrestricted access to everyone on the system, which can lead to:

*   **Sensitive Data Exposure:** Confidential files can be read by anyone.
*   **Unauthorized Modification:** Files can be altered or deleted by any user.
*   **Malicious Code Execution:** Executable files or scripts can be run by unauthorized individuals, potentially compromising the system.

Best practice dictates adhering to the **principle of least privilege**, meaning you should assign only the minimum necessary permissions. Here are recommended secure alternatives:

*   **For Directories:**
    *   `775` (rwx for owner/group, rx for others): Common for shared directories where group members need to create/delete files, but others only need to read/traverse.
    *   `770` (rwx for owner/group, no access for others): Stricter for sensitive shared directories.
    *   `755` (rwx for owner, rx for group/others): Standard for public web directories or user home directories.
*   **For Files:**
    *   `664` (rw for owner/group, r for others): Common for shared data files.
    *   `660` (rw for owner/group, no access for others): Stricter for sensitive shared files.
    *   `644` (rw for owner, r for group/others): Standard for configuration files or publicly readable data.
    *   `755` (rwx for owner, rx for group/others): Only for executable scripts.

### Bad Example

```zsh
chmod 777 sensitive_script.sh # Grants execute to everyone
chmod -R 777 public_html/     # Makes all files and subdirectories writable by everyone
```

### Good Example

```zsh
chmod 755 executable_script.zsh # Executable for owner, readable/executable for group and others
chmod 644 config.yaml          # Readable for everyone, writable only by owner
chmod -R 775 shared_project_data/ # Shared directory: group members can modify, others can read/traverse
```

### Configuration

To disable this Kata, add `ZC1007` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1008"></div>

<details>
<summary><strong>ZC1008</strong>: Use `$(())` for arithmetic operations <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

In Zsh, `$(())` is the preferred and most idiomatic syntax for **arithmetic expansion**. It allows for direct C-style arithmetic evaluation, returning the result of the expression without requiring external commands or explicit variable assignments for the result. Compared to alternatives like `expr` or `let`:

*   **Conciseness & Integration:** `$(())` is more compact and integrates seamlessly into other command substitutions or variable assignments.
*   **Efficiency:** It's a built-in shell feature, making it generally more efficient than invoking an external command like `expr`.
*   **Clarity:** It clearly signals an arithmetic context, improving script readability.

Using `expr` involves executing an external command, which adds overhead and can be less efficient for frequent operations. `let` is a built-in that performs arithmetic and assigns the result to a shell variable, but it doesn't directly return a value for use in pipelines or other command contexts, making `$(())` more versatile for expansion.

### Bad Example

```zsh
## Using expr (external command, less efficient)
val_expr=$(expr 10 + 5)
echo "Result from expr: $val_expr"

## Using let (assigns, doesn't return for expansion)
x=5; y=3
let "sum = x + y"
echo "Result from let: $sum"
```

### Good Example

```zsh
## Using $((...)) for direct arithmetic expansion
val_arith=$(( 10 + 5 ))
echo "Result from arithmetic expansion: $val_arith"

x=5; y=3
sum_arith=$(( x + y ))
echo "Result from arithmetic expansion: $sum_arith"

## Example in a conditional context
if (( a > b )); then
  echo "a is greater than b"
fi
```

### Configuration

To disable this Kata, add `ZC1008` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1009"></div>

<details>
<summary><strong>ZC1009</strong>: Use `((...))` for C-style arithmetic <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

In Zsh, the `((...))` construct is specifically designed for performing **C-style integer arithmetic evaluation and comparisons**. This is a powerful and efficient built-in mechanism to handle arithmetic directly within the shell, offering features like variable assignment, increment/decrement operators, and various logical/bitwise operations.

Key aspects of `((...))`:

*   **Conditional Context:** When used in a conditional context (e.g., `if (( ... ))`), `((...))` returns an exit status of `0` (true) if the arithmetic result is non-zero, and `1` (false) if the result is zero. This makes it ideal for conditional logic based on arithmetic.
*   **Efficiency:** As a built-in feature, it is generally more efficient than invoking external utilities like `expr`.
*   **Conciseness:** It provides a clean, C-like syntax for arithmetic operations.

This Kata promotes using `((...))` for arithmetic contexts where its **exit status** is relevant (e.g., in `if` or `while` statements) or for direct C-style assignments/increments. For cases where you need the **result of an arithmetic expression** as a value (e.g., to assign to another variable or use in a command substitution), the arithmetic expansion `$(())` is more appropriate (see [[Katas/ZC1000-ZC1099/ZC1008 | ZC1008]]).

### Bad Example

```zsh
## Using expr for a conditional check (external command, inefficient)
if expr $a + $b > /dev/null; then echo "Sum is non-zero"; fi

## Using let for incrementing (less idiomatic for C-style operation)
let i=i+1

## Performing comparison with [ (string comparison, error-prone)
if [ "$count" -lt 10 ]; then echo "Count is less than 10"; fi
```

### Good Example

```zsh
## Using ((...)) for a conditional check based on arithmetic result
if (( a + b )); then echo "Sum is non-zero"; fi

## Using ((...)) for C-style increment/decrement
(( i++ ))
(( --j ))

## Using ((...)) for arithmetic comparison in a conditional context
count=5
if (( count < 10 )); then echo "Count is less than 10"; fi
if (( num_files >= 100 )); then echo "Many files!"; fi

## Direct assignment within ((...))
(( result = a * b ))
```

### Configuration

To disable this Kata, add `ZC1009` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1010"></div>

<details>
<summary><strong>ZC1010</strong>: Use `[[ ... ]]` instead of `[ ... ]` <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

In Zsh, `[[ ... ]]` is the preferred, more robust, and more feature-rich conditional construct compared to `[ ... ]` (which is typically an alias for the external `test` command or a shell built-in with similar behavior). The primary advantages of `[[...]]` include:

*   **Built-in Keyword:** `[[...]]` is a shell keyword, not an external command. This means it behaves more predictably and efficiently, without the overhead of invoking a separate process.
*   **Intelligent Word Splitting and Globbing:** `[[...]]` handles word splitting and pathname expansion (globbing) intelligently. You generally **do not need to quote variables** within `[[...]]` unless you specifically intend for literal string matching against patterns that could otherwise be interpreted as globs. This significantly reduces common scripting errors and improves robustness.
*   **Enhanced Features:** It offers Zsh-specific features and improved syntax:
    *   **Regular Expression Matching (`=~`):** Allows direct regex matching within the conditional.
    *   **Glob Pattern Matching:** Supports extended glob patterns without explicit glob qualifiers.
    *   **Logical Operators:** Uses C-style logical operators (`&&` for AND, `||` for OR) directly, avoiding the need for separate `-a` or `-o` flags which can have surprising precedence issues.

Using `[ ... ]` requires meticulous quoting of variables and expressions to prevent unintended word splitting and pathname expansion, which can lead to bugs, security vulnerabilities, or simply incorrect logic. For arithmetic comparisons, `((...))` is often more appropriate (see [[Katas/ZC1000-ZC1099/ZC1003 | ZC1003]]).

### Bad Example

```zsh
my_string="hello world"
if [ -n $my_string ]; then echo "This might incorrectly split arguments"; fi
## Fails if $my_string contains spaces (e.g., test -n hello world) or expands to a glob pattern

file_path="/path/to/my file.txt"
if [ -f $file_path ]; then echo "File exists"; fi # Will fail if file_path has spaces and is unquoted

## Complex logic with -a (prone to precedence issues)
if [ $a -gt 10 -a $b -lt 20 ]; then echo "Range"; fi
```

### Good Example

```zsh
my_string="hello world"
if [[ -n $my_string ]]; then echo "Correctly handles spaces and globs"; fi
## No quotes needed for $my_string here, `[[...]]` handles it robustly.

file_path="/path/to/my file.txt"
if [[ -f $file_path ]]; then echo "File exists"; fi # Safely handles spaces in file names

## Clearer logical operators
if [[ $a -gt 10 && $b -lt 20 ]]; then echo "Range"; fi

## Regular expression matching
if [[ "$version" =~ "^v[0-9]+\.[0-9]+$" ]]; then echo "Valid version format"; fi
```

### Configuration

To disable this Kata, add `ZC1010` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1011"></div>

<details>
<summary><strong>ZC1011</strong>: Use `git` porcelain commands instead of plumbing commands <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Git commands are broadly categorized into two main types:

*   **Porcelain commands:** These are high-level, user-friendly commands designed for common day-to-day operations (e.g., `git status`, `git commit`, `git pull`). Their output and behavior are generally stable across Git versions, making them reliable for scripting and user interaction.
*   **Plumbing commands:** These are low-level commands that interact directly with Git's internal data structures (e.g., `git hash-object`, `git cat-file`, `git update-index`). Their output and exact behavior can be more prone to change between Git versions, and they require a deeper understanding of Git's internals.

In typical shell scripts, it is **highly recommended to use porcelain commands** unless you have a very specific and well-understood need to interact with Git's internals. Relying on plumbing commands for common tasks can lead to fragile scripts that break with Git updates or are harder to understand and maintain.

### Bad Example

```zsh
## Using a plumbing command to get the current commit hash (less robust)
commit_hash=$(git rev-parse HEAD^{commit})

## Using an older or less direct method for current branch name
current_branch=$(git branch --show-current) # deprecated in some contexts / for older Git versions

## Using plumbing to check if a file is tracked (less readable)
is_tracked=$(git ls-files --error-unmatch -- "$file" &>/dev/null; echo $?)
```

### Good Example

```zsh
## Standard and robust way to get the current commit hash
commit_hash=$(git rev-parse HEAD)

## Recommended way to get the current branch name (porcelain)
current_branch=$(git rev-parse --abbrev-ref HEAD)

## Using porcelain to check if a file is tracked (more common and readable)
if git ls-files --error-unmatch -- "$file" &>/dev/null; then
  echo "$file is tracked"
else
  echo "$file is not tracked"
fi
```

### Configuration

To disable this Kata, add `ZC1011` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1012"></div>

<details>
<summary><strong>ZC1012</strong>: Use `read -r` to prevent backslash escaping <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

When using the `read` built-in command in Zsh (and other shells), it's crucial to use the `-r` option (`raw` mode). By default, `read` interprets backslashes (`\`) as escape characters, which can lead to unexpected behavior if the input contains backslashes (e.g., file paths, special characters, or multi-line input). The `-r` option prevents `read` from interpreting backslashes, ensuring that the input is stored literally as it was entered. This significantly improves script robustness and predictability, especially when reading user input or arbitrary data from files.

### Bad Example

```zsh
read -p "Enter a string (e.g., C:\Program Files or line1\\nline2): " user_input
echo "You entered: $user_input"
## If user enters "C:\Program Files", it might be interpreted as "C:Program Files"
## If user enters "line1\\nline2", it might be interpreted as "line1\nline2"
```

### Good Example

```zsh
read -r -p "Enter a string (e.g., C:\Program Files or line1\\nline2): " user_input
echo "You entered: $user_input"
## "C:\Program Files" is stored exactly as entered, "line1\\nline2" is stored as "line1\nline2"
```

### Configuration

To disable this Kata, add `ZC1012` to the `disabled_katas` list in your `.zshellcheckrc` file.

```

```

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1013"></div>

<details>
<summary><strong>ZC1013</strong>: Use `((...))` for arithmetic operations instead of `let` <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

While `let` is a built-in command for performing arithmetic operations and assignments, Zsh's `((...))` construct offers a more modern, idiomatic, and versatile C-style syntax for arithmetic. Preferring `((...))` over `let` brings several benefits:

*   **C-style Syntax:** `((...))` provides a familiar C-like syntax for arithmetic expressions, which is often more intuitive for developers coming from other programming languages.
*   **Versatility:** It can be used directly for both arithmetic evaluation (where its exit status is based on the result) and assignments, including compound assignments (`+=`, `-=`) and increment/decrement operators (`++`, `--`).
*   **Readability:** `((...))` makes arithmetic expressions clearer and more concise, especially for complex calculations or when integrating with conditional logic.
*   **Efficiency:** As a built-in shell feature, `((...))` is generally efficient.

`let` can sometimes be less readable, especially when expressions involve multiple variables or operators, and might require more careful quoting. This Kata encourages consistency and readability by promoting `((...))` for all C-style arithmetic needs.

### Bad Example

```zsh
count=0
let count=count+1      # Less idiomatic, requires re-typing variable name

total=10
price=2
quantity=5
let "total = price * quantity" # More verbose quoting for expressions

## Using let for comparison (often less clear)
if let "x < 10"; then echo "x is less"; fi
```

### Good Example

```zsh
count=0
(( count++ ))          # Concise C-style increment
(( count += 5 ))       # Compound assignment

total=10
price=2
quantity=5
(( total = price * quantity )) # Clean assignment

## Using ((...)) for comparison (clearer exit status)
x=5
if (( x < 10 )); then echo "x is less"; fi
```

### Configuration

To disable this Kata, add `ZC1013` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1014"></div>

<details>
<summary><strong>ZC1014</strong>: Use `git switch` or `git restore` instead of `git checkout` <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Git has introduced the specialized commands `git switch` and `git restore` to explicitly handle the distinct operations previously managed by the single, overloaded `git checkout` command. Adopting these newer commands offers significant benefits:

*   **Improved Clarity:** Each command now has a single, clear responsibility:
    *   `git switch`: Dedicated to changing branches.
    *   `git restore`: Dedicated to reverting modified files in the working tree or restoring files from the staging area.
*   **Reduced Ambiguity:** The original `git checkout` command had three main functions (switching branches, restoring working tree files, restoring staging index files), which could lead to confusion and accidental data loss. `git switch` and `git restore` eliminate this ambiguity.
*   **Modern Git Practices:** Using `git switch` and `git restore` aligns with current Git best practices, making scripts and workflows easier to understand and more resilient to future Git changes.

While scripts still using `git checkout` for these purposes remain functional, migrating to `git switch` and `git restore` improves readability, maintainability, and safety.

### Bad Example

```zsh
## Overloaded `git checkout` for switching branches
git checkout feature-branch

## Overloaded `git checkout` for discarding local changes (risky if not careful)
git checkout -- my_modified_file.txt

## Overloaded `git checkout` to restore a file from staging (less intuitive)
git checkout HEAD -- staged_file.txt
```

### Good Example

```zsh
## Clearer `git switch` for changing branches
git switch feature-branch

## Clearer `git restore` for discarding local changes in the working tree
git restore my_modified_file.txt

## Clearer `git restore` to unstage a file (restore from index to working tree, often implicit)
git restore --staged another_file.txt

## Restore a specific file from a commit (equivalent to checkout <commit> -- <file>)
git restore --source=HEAD~1 file_from_prev_commit.txt
```

### Configuration

To disable this Kata, add `ZC1014` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1015"></div>

<details>
<summary><strong>ZC1015</strong>: Use <code>$(...)</code> for command substitution instead of backticks <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

The `$(...)` syntax is the modern, recommended way to perform command substitution. It is easier to nest and less error-prone than backticks.

### Bad Example

```zsh
result=`ls -la`
nested=`echo \`date\``
```

### Good Example

```zsh
result=$(ls -la)
nested=$(echo $(date))
```

### Configuration

To disable this Kata, add `ZC1015` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1016"></div>

<details>
<summary><strong>ZC1016</strong>: Use `read -s` when reading sensitive information <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

When asking for passwords or secrets, use `read -s` to prevent the input from being echoed to the terminal.

### Bad Example

```zsh
read password
read "token?Enter API Token: "
```

### Good Example

```zsh
read -s password
read -s "token?Enter API Token: "
```

### Configuration

To disable this Kata, add `ZC1016` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>


<div id="zc1017"></div>

<details>
<summary><strong>ZC1017</strong>: Use `print -r` to print strings literally <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

When printing arbitrary strings or variable contents, especially those that might contain backslashes (`\`) or other characters that `echo` could interpret as escape sequences, it is best practice to use `print -r`. The `print` builtin with the `-r` (raw) option ensures that backslashes are treated literally and not as escape characters. This guarantees that the output matches the input, preventing unintended interpretations or formatting issues. `echo`'s behavior can vary between shells and even different versions of the same shell, making `print -r` a more reliable and portable choice for literal output.

### Bad Example

```zsh
echo "Path: C:\Program Files"
echo -e "Hello\nWorld" # -e is an echo extension
```

### Good Example

```zsh
print -r "Path: C:\Program Files"
print -r "Hello\\nWorld" # Backslashes are printed literally
print "Hello\nWorld"    # print without -r interprets escapes
```

### Configuration

To disable this Kata, add `ZC1017` to the `disabled_katas` list in your `.zshellcheckrc` file.

```

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1018"></div>

<details>
<summary><strong>ZC1018</strong>: Use `((...))` for C-style arithmetic instead of `expr` <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

This Kata reinforces the use of Zsh's native `((...))` arithmetic construct over the external `expr` command. `((...))` is a powerful built-in that provides C-style arithmetic evaluation, integer comparisons, and variable assignments directly within the shell. It is generally more efficient than invoking an external `expr` command, which incurs process overhead. Furthermore, `((...))` syntax is cleaner and less prone to quoting issues that often plague `expr` expressions. Consistent use of `((...))` improves script performance, readability, and robustness.

### Bad Example

```zsh
value=$(expr $x \* $y + 5)
expr $i = $j > /dev/null
```

### Good Example

```zsh
value=$(( x * y + 5 ))
(( i == j ))
```

### Configuration

To disable this Kata, add `ZC1018` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1019"></div>

<details>
<summary><strong>ZC1019</strong>: Use `whence` instead of `which` <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

This Kata is a reinforcement of ZC1005. It re-emphasizes the importance of using the Zsh built-in `whence` over the external `which` utility. `whence` is superior because it accurately reports how a command will be found and executed by the shell, including considerations for aliases, functions, and built-ins. `which` only searches `PATH` for executable files, making it unreliable for determining the actual command Zsh will run. Preferring `whence` ensures that scripts accurately reflect Zsh's command lookup behavior.

### Bad Example

```zsh
which my_alias
## If 'my_alias' is an alias, 'which' won't show its definition
```

### Good Example

```zsh
whence my_alias
## 'whence' will show if 'my_alias' is an alias, function, or builtin
```

### Configuration

To disable this Kata, add `ZC1019` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1020"></div>

<details>
<summary><strong>ZC1020</strong>: Use `[[ ... ]]` for tests instead of `test` <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

This Kata is a reinforcement of ZC1010. It reiterates the best practice of using Zsh's `[[ ... ]]` conditional command instead of the traditional `test` utility. `[[ ... ]]` is a shell keyword, providing superior safety and functionality. It automatically handles word splitting and globbing, preventing common errors that arise from unquoted variables in `test` or `[ ... ]`. Additionally, `[[ ... ]]` offers Zsh-specific features like regular expression matching and extended globbing patterns. Consistently using `[[ ... ]]` leads to more robust, readable, and less error-prone Zsh scripts.

### Bad Example

```zsh
if test -z "$var"; then echo "Var is empty"; fi
if test $num -eq 10; then echo "Num is 10"; fi
```

### Good Example

```zsh
if [[ -z "$var" ]]; then echo "Var is empty"; fi
if [[ $num -eq 10 ]]; then echo "Num is 10"; fi
```

### Configuration

To disable this Kata, add `ZC1020` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1021"></div>

<details>
<summary><strong>ZC1021</strong>: Use symbolic permissions with `chmod` instead of octal <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

While octal notation (`755`, `644`) for `chmod` is concise, symbolic mode (`u+rwx`, `go=rx`) is often more readable, less error-prone, and clearer about the *intent* of the permission change, especially for complex adjustments. Symbolic mode allows you to add, remove, or set specific permissions for user, group, or others without needing to calculate the new octal value. This reduces the risk of accidentally granting unwanted permissions and makes scripts easier to understand for collaborators.

### Bad Example

```zsh
chmod 755 script.sh
chmod 600 config.txt
```

### Good Example

```zsh
chmod u=rwx,go=rx script.sh # Equivalent to 755
chmod u=rw,go= config.txt  # Equivalent to 600
chmod +x install.sh         # Add execute permission for all
```

### Configuration

To disable this Kata, add `ZC1021` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1022"></div>

<details>
<summary><strong>ZC1022</strong>: Use `$((...))` for arithmetic expansion <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

This Kata is a reinforcement of ZC1008. It reiterates that `$(())` is the canonical Zsh syntax for performing arithmetic expansion. This construct evaluates an arithmetic expression and substitutes its result into the command line. It's concise, efficient, and avoids the complexities and overhead of older methods like `expr` or `let` when you need the result of an. It's concise, efficient, and avoids the complexities and overhead of older methods like `expr` or `let` when you need the result of an arithmetic calculation. Consistently using `$(())` enhances script readability, performance, and compatibility across modern Zsh environments.

### Bad Example

```zsh
VAL=$(( 1 + 2 ))
echo "The value is $(( VAL ))"
```

### Good Example

```zsh
VAL=$((1 + 2))
echo "The value is $((VAL))"
```

### Configuration

To disable this Kata, add `ZC1022` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1023"></div>

<details>
<summary><strong>ZC1023</strong>: Use `$((...))` for arithmetic expansion <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

This Kata is another reinforcement of ZC1008 and ZC1022. It consistently promotes the use of `$(())` for all arithmetic expansions within Zsh scripts. This syntax is the most idiomatic and robust way to perform integer calculations and retrieve their results. Using `$(())` simplifies complex arithmetic, improves script readability, and leverages Zsh's built-in capabilities efficiently, avoiding reliance on external tools or less clear syntaxes.

### Bad Example

```zsh
x=$(( 1 + 1 )) # Bad due to spacing, although functional
total=$(let "a = 1 + 2"; echo $a)
```

### Good Example

```zsh
x=$((1+1))
total=$((1+2))
```

### Configuration

To disable this Kata, add `ZC1023` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1024"></div>

<details>
<summary><strong>ZC1024</strong>: Use `$((...))` for arithmetic expansion <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

This Kata serves as another reinforcement of `$(())` as the standard for arithmetic expansion in Zsh. Its purpose is to ensure consistent adoption of this modern, efficient, and readable syntax throughout Zsh code. Avoiding older, less clear, or external methods for arithmetic calculations contributes to more maintainable and performant scripts. This emphasis helps prevent subtle bugs and promotes a unified coding style.

### Bad Example

```zsh
VAL=`expr $X + $Y`
let "TOTAL = 10 * 5" && echo $TOTAL
```

### Good Example

```zsh
VAL=$((X + Y))
echo $((10 * 5))
```

### Configuration

To disable this Kata, add `ZC1024` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1025"></div>

<details>
<summary><strong>ZC1025</strong>: Use `$((...))` for arithmetic expansion <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

This Kata consistently emphasizes the use of `$(())` for all arithmetic expansions within Zsh. It reinforces the principle that this construct is the most robust, efficient, and idiomatic way to perform integer calculations and obtain their results in Zsh. By promoting its consistent use, ZShellCheck aims to standardize arithmetic operations, improve script readability, and prevent issues arising from less standard or more cumbersome syntaxes.

### Bad Example

```zsh
result=`echo "2 * 3" | bc`
var=$(( 1 + 1 )) # Bad spacing
```

### Good Example

```zsh
result=$(( 2 * 3 ))
var=$((1+1))
```

### Configuration

To disable this Kata, add `ZC1025` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1026"></div>

<details>
<summary><strong>ZC1026</strong>: Use `$((...))` for arithmetic expansion <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

This Kata, a further reiteration, stresses the importance of exclusively using the `$(())` syntax for arithmetic expansion in Zsh. It highlights that `$(())` is the native and most efficient method for performing integer calculations within Zsh. Adhering to this practice ensures code clarity, consistency, and optimal performance, minimizing reliance on slower external processes or less robust shell-specific arithmetic features.

### Bad Example

```zsh
num=$(($VAL + 1)) # Still functional but inconsistent spacing
```

### Good Example

```zsh
num=$((VAL + 1))
```

### Configuration

To disable this Kata, add `ZC1026` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1027"></div>

<details>
<summary><strong>ZC1027</strong>: Use `$((...))` for arithmetic expansion <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

This Kata is another reinforcement of the standard for arithmetic expansion in Zsh. It emphasizes the use of `$(())` for all integer calculations where the result needs to be substituted into a command. This is the idiomatic Zsh way, providing superior readability, efficiency, and robustness compared to alternative methods. Consistent application of this syntax helps maintain high code quality and avoids potential parsing ambiguities or performance overheads.

### Bad Example

```zsh
RESULT=$(("$VAR" + 1)) # Unnecessary quoting
```

### Good Example

```zsh
RESULT=$((VAR + 1))
```

### Configuration

To disable this Kata, add `ZC1027` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1028"></div>

<details>
<summary><strong>ZC1028</strong>: Use `$((...))` for arithmetic expansion <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

This Kata is part of a series reinforcing the use of Zsh's `$(())` for arithmetic expansion. It underscores that this built-in construct is the most effective and idiomatic way to perform integer calculations and obtain their results within Zsh. Promoting its consistent usage throughout Zsh code leads to enhanced clarity, maintainability, and ensures optimal performance by leveraging native shell capabilities rather than external utilities or deprecated syntaxes.

### Bad Example

```zsh
VAL=$(expr $X + 1)
## Using 'expr' for simple increment
```

### Good Example

```zsh
VAL=$((X + 1))
```

### Configuration

To disable this Kata, add `ZC1028` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1029"></div>

<details>
<summary><strong>ZC1029</strong>: Use `$((...))` for arithmetic expansion <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

As part of the series emphasizing modern Zsh practices, this Kata advocates for the consistent and exclusive use of `$(())` for arithmetic expansion. This built-in feature offers the most efficient, readable, and robust method for performing calculations and substituting their results into commands. By adopting `$(())`, scripts become more predictable, easier to maintain, and avoid common pitfalls associated with older or less integrated arithmetic approaches in Zsh.

### Bad Example

```zsh
RESULT=$(( VAR + 1 )) # Unnecessary spacing, though functional
```

### Good Example

```zsh
RESULT=$((VAR + 1))
```

### Configuration

To disable this Kata, add `ZC1029` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1030"></div>

<details>
<summary><strong>ZC1030</strong>: Use `printf` instead of `echo` <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

For printing formatted output, `printf` is generally superior and more portable than `echo` in Zsh (and other shells). `printf` provides precise control over output formatting through format specifiers (e.g., `%s` for string, `%d` for integer), handles escape sequences predictably, and prevents issues with variable expansion that `echo` might encounter. `echo`'s behavior can vary significantly between shells and even versions, making it less reliable for consistent output. Using `printf` ensures robust and predictable formatting.

### Bad Example

```zsh
echo "Hello, $name"
echo -n "Progress: " # -n for no newline, not universally supported by echo
```

### Good Example

```zsh
printf "Hello, %s\n" "$name"
printf "Progress: "
```

### Configuration

To disable this Kata, add `ZC1030` to the `disabled_katas` list in your `.zshellcheckrc` file.

```

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1031"></div>

<details>
<summary><strong>ZC1031</strong>: Use `#!/usr/bin/env zsh` for portability <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

For Zsh scripts, it is best practice to use `#!/usr/bin/env zsh` as the shebang line instead of a hardcoded path like `#!/bin/zsh` or `#!/usr/bin/zsh`. The `env` utility locates the `zsh` executable in the user's `PATH`, ensuring that the script uses the `zsh` version intended by the user, rather than a specific system-installed version that might not exist or be outdated on all systems. This significantly improves script portability across different environments.

### Bad Example

```zsh
#!/bin/zsh
## Your script...
```

### Good Example

```zsh
#!/usr/bin/env zsh
## Your script...
```

### Configuration

To disable this Kata, add `ZC1031` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1032"></div>

<details>
<summary><strong>ZC1032</strong>: Use `((...))` for C-style incrementing <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

This Kata emphasizes using Zsh's native `((...))` construct for C-style variable incrementing and decrementing (e.g., `i++`, `j--`). This is the most idiomatic, efficient, and concise way to perform such operations in Zsh. It leverages Zsh's built-in arithmetic capabilities, avoiding the need for `let` or manual assignments, which can be more verbose and less readable. Consistent use of `((...))` for increments/decrements improves code clarity and performance.

### Bad Example

```zsh
let i=i+1
j=$((j-1))
```

### Good Example

```zsh
(( i++ ))
(( j-- ))
```

### Configuration

To disable this Kata, add `ZC1032` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1033"></div>

<details>
<summary><strong>ZC1033</strong>: Use `$((...))` for arithmetic expansion <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

This Kata is part of the recurring theme to promote the `$(())` syntax for all arithmetic expansions in Zsh. It highlights that this built-in construct is the most robust, efficient, and readable method for performing integer calculations and substituting their results into the command line. By consistently using `$(())`, scripts become more predictable, easier to maintain, and avoid common pitfalls associated with older or less integrated arithmetic approaches in Zsh.

### Bad Example

```zsh
VAL=`expr $X + 1` # Using backticks with expr
let result=$((5*5)) # Combining let with $((...)) unnecessarily
```

### Good Example

```zsh
VAL=$((X + 1))
result=$((5*5))
```

### Configuration

To disable this Kata, add `ZC1033` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1034"></div>

<details>
<summary><strong>ZC1034</strong>: Use `command -v` instead of `which` <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

In Zsh (and other shells), `command -v` is the most robust and portable way to determine the full pathname or definition of a command, including handling aliases, functions, and built-ins. While `type` can provide similar information, its output format can vary more widely across shells, making it less suitable for programmatic parsing in scripts. `command -v` is specifically designed for reliable output in scripts and adheres to POSIX standards, ensuring consistent behavior.

### Bad Example

```zsh
type my_command # Output format can vary
```

### Good Example

```zsh
command -v my_command # Reliable output for scripting
```

### Configuration

To disable this Kata, add `ZC1034` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1035"></div>

<details>
<summary><strong>ZC1035</strong>: Use `$((...))` for arithmetic expansion <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

This Kata is another in a series emphasizing the consistent use of Zsh's `$(())` for all arithmetic expansions. This built-in construct is the idiomatic, efficient, and reliable method for performing integer calculations and substituting their results into command arguments or variable assignments. Adhering to this standard improves script readability, maintainability, and ensures predictable behavior across different Zsh environments by avoiding external tools or less standard syntaxes.

### Bad Example

```zsh
VAL=$(echo $((1+1))) # Unnecessary echo and nested $((...))
```

### Good Example

```zsh
VAL=$((1+1))
```

### Configuration

To disable this Kata, add `ZC1035` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1036"></div>

<details>
<summary><strong>ZC1036</strong>: Prefer `[[ ... ]]` over `test` command <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

This Kata is another reinforcement of the best practice to use `[[ ... ]]` for conditional expressions in Zsh, instead of the older `test` command or single brackets `[ ... ]`. As a shell keyword, `[[ ... ]]` offers superior reliability by preventing issues with word splitting and globbing, and provides extended functionality like regular expression matching and robust pattern matching without complex quoting. Its consistent use makes Zsh scripts more resilient, readable, and aligned with modern shell scripting conventions.

### Bad Example

```zsh
if test -d /tmp; then echo "Dir exists"; fi
if [ -n "$VAR" ]; then echo "Var is set"; fi
```

### Good Example

```zsh
if [[ -d /tmp ]]; then echo "Dir exists"; fi
if [[ -n "$VAR" ]]; then echo "Var is set"; fi
```

### Configuration

To disable this Kata, add `ZC1036` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1037"></div>

<details>
<summary><strong>ZC1037</strong>: Use `print -r --` for variable expansion <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

When printing the raw, uninterpreted value of a variable, especially one that might begin with a hyphen (`-`) or contain backslashes (`\`), it's best practice to use `print -r --`. The `-r` (raw) option prevents backslash interpretation, and the `--` argument explicitly signals the end of options, preventing any subsequent variable content that starts with a hyphen from being mistakenly parsed as another option to `print`. This ensures literal output and enhances script robustness against unpredictable variable content.

### Bad Example

```zsh
file_name="-my-file.txt"
echo $file_name # Could be interpreted as an option
print $file_name # Might interpret escapes
```

### Good Example

```zsh
file_name="-my-file.txt"
print -r -- "$file_name" # Prints literally, handles hyphens safely
```

### Configuration

To disable this Kata, add `ZC1037` to the `disabled_katas` list in your `.zshellcheckrc` file.

```

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1038"></div>

<details>
<summary><strong>ZC1038</strong>: Avoid useless use of `cat` <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

It's considered a bad practice, often called a "Useless Use of Cat" (UUOC), to use `cat` to pipe the content of a single file to another command when that command can read the file directly. This adds an unnecessary process to the pipeline, incurring overhead without providing any benefit. Many commands (e.g., `grep`, `sed`, `awk`) can accept a filename as an argument, making `cat file | command` redundant.

### Bad Example

```zsh
cat file.txt | grep "pattern"
cat another.log | sed 's/old/new/'
```

### Good Example

```zsh
grep "pattern" file.txt
sed 's/old/new/' another.log
```

### Configuration

To disable this Kata, add `ZC1038` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1039"></div>

<details>
<summary><strong>ZC1039</strong>: Avoid `rm` with root path <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Executing `rm` on the root directory (`/`) is extremely dangerous and can lead to irreversible data loss, rendering the system unusable. This Kata warns against using `rm` with `/` as an argument, even if other options like `-rf` are not explicitly present. While modern `rm` implementations often have built-in safeguards, relying on them is not a substitute for careful scripting. It's a critical safety measure to prevent accidental deletion of the entire filesystem.

### Bad Example

```zsh
rm -rf /
rm / # Even without -rf, this is dangerous
```

### Good Example

```zsh
rm -rf /tmp/my_dir
rm -rf ${MY_VAR}/subdir # Ensure MY_VAR is not "/"
```

### Configuration

To disable this Kata, add `ZC1039` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1040"></div>

<details>
<summary><strong>ZC1040</strong>: Use `(N)` nullglob qualifier for globs in loops <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

In Zsh, if a glob pattern (like `*.txt`) does not match any files, it traditionally expands to itself (e.g., `*.txt`). This can lead to unexpected errors in loops or commands that assume a match. The `(N)` (nullglob) glob qualifier is a Zsh-specific feature that makes the glob expand to *nothing* if no matches are found, behaving more predictably in such scenarios. Using `(N)` ensures loops over globs don't mistakenly iterate over the literal pattern.

### Bad Example

```zsh
for file in *.txt; do # If no .txt files, loops once with literal "*.txt"
    echo "Processing $file"
done
```

### Good Example

```zsh
for file in *.txt(N); do # If no .txt files, loop is skipped
    echo "Processing $file"
done
```

### Configuration

To disable this Kata, add `ZC1040` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1041"></div>

<details>
<summary><strong>ZC1041</strong>: Do not use variables in `printf` format string <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Using a variable directly as the format string in `printf` can be a security vulnerability, especially if the variable's content comes from untrusted input. Malicious input could inject format specifiers (e.g., `%x`, `%n`) that might reveal stack contents, write to arbitrary memory locations, or cause a crash. Always use a static, literal string for the `printf` format and pass variables as separate arguments. If dynamic formatting is truly needed, sanitize the input or carefully construct the format string to prevent injection.

### Bad Example

```zsh
user_input="%s %s %s"
printf $user_input "hello" "world" # Vulnerable to format string attacks
```

### Good Example

```zsh
user_input="hello world"
printf "%s\n" "$user_input" # Safely prints the variable content
```

### Configuration

To disable this Kata, add `ZC1041` to the `disabled_katas` list in your `.zshellcheckrc` file.

```

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1042"></div>

<details>
<summary><strong>ZC1042</strong>: Use "$@" to iterate over arguments <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

When iterating over all positional parameters (arguments passed to a script or function), it is best practice to use `"$@"`. This construct correctly preserves word boundaries and individual arguments, even if they contain spaces or special characters. Using `$*` (unquoted) or `$@` (unquoted) will typically cause word splitting, treating spaces within an argument as delimiters and leading to incorrect iteration. Always quote `"$@"` to ensure each argument is passed as a distinct word.

### Bad Example

```zsh
for arg in $*; do echo "Arg: $arg"; done # Splits arguments with spaces
for arg in $@; do echo "Arg: $arg"; done # Also splits
```

### Good Example

```zsh
for arg in "$@"; do echo "Arg: $arg"; done # Correctly preserves arguments
my_func() {
    for arg in "$@"; do
        echo "Function Arg: $arg"
    }
}
my_func "arg one" "arg two"
```

### Configuration

To disable this Kata, add `ZC1042` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1043"></div>

<details>
<summary><strong>ZC1043</strong>: Use `local` for variables in functions <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

In Zsh (and other shells), variables defined within a function are global by default unless explicitly declared as local. This can lead to unintended side effects, where a function modifies a global variable with the same name, causing unexpected behavior in other parts of the script. Best practice dictates using the `local` keyword (or `typeset` for more options) to scope variables to the function, preventing name collisions and making functions self-contained and predictable.

### Bad Example

```zsh
global_var="hello"
my_func() {
    global_var="world" # Modifies global_var
    temp_var="inside"  # Also global by default
}
```

### Good Example

```zsh
global_var="hello"
my_func() {
    local global_var="world" # This is a new local variable
    local temp_var="inside"  # Scoped to function
}
```

### Configuration

To disable this Kata, add `ZC1043` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1044"></div>

<details>
<summary><strong>ZC1044</strong>: Check for unchecked `cd` commands <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

A `cd` command can fail for various reasons (e.g., directory does not exist, permission denied). If a `cd` command is not checked for success, subsequent commands in the script might operate in an unexpected directory, potentially leading to errors, data corruption, or security vulnerabilities. It is crucial to always check the exit status of `cd` and handle potential failures, typically by exiting the script, returning from a function, or providing an error message.

### Bad Example

```zsh
cd /path/to/nonexistent_dir
rm * # This might now run in the wrong directory!
```

### Good Example

```zsh
cd /path/to/my_dir || exit 1 # Exit if cd fails
## Or in a function:
cd_safe() {
    cd "$1" || { print "Error: Cannot change to $1"; return 1; }
}
cd_safe /path/to/my_dir && echo "Successfully changed directory"
```

### Configuration

To disable this Kata, add `ZC1044` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1045"></div>

<details>
<summary><strong>ZC1045</strong>: Declare and assign separately to avoid masking return values <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

When declaring and assigning a variable simultaneously on the same line (e.g., `local var=$(command)`), the exit status of the `command` might not be directly propagated to the `$?` variable of the outer shell. Instead, the exit status might reflect the success or failure of the `local` (or `typeset`) command itself, masking the actual success of the command substitution. To reliably capture the command's exit status, it's safer to perform the assignment in two separate steps: first execute the command, then assign its output to the variable.

### Bad Example

```zsh
local output=$(my_command)
if (( $? != 0 )); then echo "my_command failed"; fi # $? might be local's exit status
```

### Good Example

```zsh
output=$(my_command)
local output # Declare local after command execution
if (( $? != 0 )); then echo "my_command failed"; fi # $? is my_command's exit status
```

### Configuration

To disable this Kata, add `ZC1045` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1046"></div>

<details>
<summary><strong>ZC1046</strong>: Avoid `eval` <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

The `eval` command executes its arguments as a shell command. While powerful, `eval` is extremely dangerous because it can introduce severe security vulnerabilities if its input is not fully trusted and sanitized. Malicious input can inject arbitrary commands, leading to remote code execution or other exploits. Due to these risks, `eval` should be avoided whenever possible. If dynamic command construction is unavoidable, consider safer alternatives like arrays for arguments or `printf %q` for quoting.

### Bad Example

```zsh
user_input="rm -rf /" # Imagine this came from user
eval "$user_input"   # Danger!
```

### Good Example

```zsh
## Avoid eval. Use explicit command structure.
## Example for dynamic command:
cmd=("ls" "-l")
if [[ $show_all = true ]]; then
    cmd+="-a"
fi
"${cmd[@]}" # Safe execution
```

### Configuration

To disable this Kata, add `ZC1046` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1047"></div>

<details>
<summary><strong>ZC1047</strong>: Avoid `sudo` in scripts <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Hardcoding `sudo` directly into scripts can be a security risk. It grants elevated privileges, often without explicit user consent or a clear understanding of the full impact of the elevated commands. This can lead to unintended system modifications, security vulnerabilities if the script is compromised, or unexpected behavior if the script's assumptions about the environment change. Scripts should ideally be run with the necessary privileges from the outset or prompt the user for `sudo` only when absolutely necessary and with clear warnings. If `sudo` is unavoidable, minimize its scope (e.g., `sudo command -arg`) rather than running an entire subshell with elevated privileges.

### Bad Example

```zsh
sudo apt update && sudo apt upgrade -y
sudo sh -c 'echo "secret" > /etc/sensitive_file'
```

### Good Example

```zsh
## User runs script with sudo
## sudo ./install.sh

## Or prompt for sudo when necessary
if (( EUID != 0 )); then
    print "This script requires root privileges. Please run with sudo."
    exit 1
fi

## Minimal sudo scope
sudo apt update
sudo tee /etc/privileged_file <<< "secret"
```

### Configuration

To disable this Kata, add `ZC1047` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1048"></div>

<details>
<summary><strong>ZC1048</strong>: Avoid `source` with relative paths <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Using `source` (or `.`) with relative paths to include other script files can be problematic, especially in scripts that might change the current working directory (`cd`) or be invoked from different locations. A relative path for `source` is resolved *relative to the current working directory at the time `source` is executed*, not relative to the script containing the `source` command. This can lead to scripts failing to find their dependencies or sourcing unintended files. It's safer to use an absolute path or a path relative to the *script's own directory*, usually determined using `dirname $0` or similar techniques.

### Bad Example

```zsh
## script.sh
## inside script.sh:
source ./config.zsh # Might fail if script.sh is called from other dir
```

### Good Example

```zsh
## script.sh
## inside script.sh:
script_dir=$(dirname "${(%):-%x}") # Zsh-specific way to get script's dir
source "${script_dir}/config.zsh"
## Or:
source /absolute/path/to/config.zsh
```

### Configuration

To disable this Kata, add `ZC1048` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1049"></div>

<details>
<summary><strong>ZC1049</strong>: Prefer functions over aliases <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

While aliases are convenient for short command shortcuts, functions offer superior flexibility, robustness, and readability in Zsh scripts and configurations. Functions allow for positional parameters (`$1`, `$2`), conditional logic, multiple commands, and local variables, making them much more powerful than simple text substitutions provided by aliases. Aliases can also interact unpredictably with quoting and command-line parsing. For anything beyond a trivial, single-word substitution, a function is the preferred choice for maintainability and predictable behavior.

### Bad Example

```zsh
alias ll="ls -lh" # Simple alias
alias commit_all="git add . && git commit -m 'Auto commit'" # Complex alias
```

### Good Example

```zsh
## Simple alias is acceptable for interactive use, but prefer function in scripts:
ll() { ls -lh "$@"; } # Function version of ll

commit_all() {
    git add .
    git commit -m "Auto commit" "$@"
}
```

### Configuration

To disable this Kata, add `ZC1049` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1050"></div>

<details>
<summary><strong>ZC1050</strong>: Avoid iterating over `ls` output <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Parsing the output of `ls` in a loop (e.g., `for file in $(ls)`) is highly unreliable and considered a bad practice. `ls` output is primarily designed for human readability, not machine parsing. Filenames can contain spaces, newlines, or special characters that will cause `for` loops (which typically split on whitespace) to misinterpret filenames, leading to incorrect processing or even unexpected command execution. Instead of `ls`, use globbing (`for file in *`) for simple file lists or `find` with `-print0` and `xargs -0` for robust, null-delimited processing of arbitrary filenames.

### Bad Example

```zsh
for file in $(ls *.txt); do
    echo "Processing $file"
done
```

### Good Example

```zsh
for file in *.txt; do # Uses globbing, handles spaces in filenames
    echo "Processing $file"
done

## For more complex scenarios, especially with arbitrary filenames
find . -name "*.txt" -print0 | while IFS= read -r -d $'\0' file; do
    echo "Processing $file"
done
```

### Configuration

To disable this Kata, add `ZC1050` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1051"></div>

<details>
<summary><strong>ZC1051</strong>: Quote variables in `rm` to avoid globbing <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

When passing variable content to `rm`, especially if that content represents a filename, it is critical to quote the variable (e.g., `rm "$file"`) to prevent unintended globbing (pathname expansion). If a variable contains characters like `*`, `?`, or `[`, and it's unquoted, Zsh will attempt to interpret it as a glob pattern before passing it to `rm`. This can lead to deleting multiple files or the wrong files entirely if the variable's value happens to match existing files. Always quote variables when passing them as arguments to `rm` to ensure they are treated as literal filenames.

### Bad Example

```zsh
file_to_delete="backup-*.zip"
rm $file_to_delete # If "backup-*.zip" matches actual files, they'll be deleted
```

### Good Example

```zsh
file_to_delete="backup-*.zip"
rm "$file_to_delete" # Ensures only the literal "backup-*.zip" file is targeted
```

### Configuration

To disable this Kata, add `ZC1051` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1052"></div>

<details>
<summary><strong>ZC1052</strong>: Avoid `sed -i` for portability <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

The `-i` option of `sed` (for in-place editing) behaves differently across various `sed` implementations (e.g., GNU sed vs. BSD/macOS sed). Specifically, BSD `sed` requires a mandatory backup extension (e.g., `sed -i .bak`), while GNU `sed` allows `-i` without an extension (or `-i''` to explicitly avoid a backup). This discrepancy makes `sed -i` non-portable across systems. For cross-platform compatibility, it's safer to perform explicit redirection: write `sed`'s output to a temporary file, then move the temporary file over the original.

### Bad Example

```zsh
sed -i 's/foo/bar/' file.txt # Fails on BSD sed without backup extension
```

### Good Example

```zsh
## Portable in-place editing
sed 's/foo/bar/' file.txt > file.tmp && mv file.tmp file.txt

## If targeting specific sed (e.g., GNU sed on Linux only)
## sed -i'' 's/foo/bar/' file.txt # GNU compatible
## sed -i '.bak' 's/foo/bar/' file.txt # BSD compatible
```

### Configuration

To disable this Kata, add `ZC1052` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1053"></div>

<details>
<summary><strong>ZC1053</strong>: Silence `grep` output in conditions <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

When using `grep` within a conditional statement (e.g., `if grep ...`, `while grep ...`) to check for the existence of a pattern, its output to `stdout` is usually undesirable. The purpose of `grep` in such contexts is typically its exit status (0 for match, non-zero for no match or error), not the matching lines themselves. Redirecting `grep`'s output to `/dev/null` using `grep -q` (quiet mode) or `&>/dev/null` prevents cluttering the terminal or pipeline with extraneous information, improving script cleanliness and preventing unintended side effects if subsequent commands process `grep`'s `stdout`.

### Bad Example

```zsh
if grep "error" log.txt; then # Prints matching lines to stdout
    echo "Found errors!"
fi
```

### Good Example

```zsh
if grep -q "error" log.txt; then # -q suppresses stdout
    echo "Found errors!"
fi
## Or:
if grep "error" log.txt &>/dev/null; then # Redirects both stdout and stderr
    echo "Found errors!"
fi
```

### Configuration

To disable this Kata, add `ZC1053` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1054"></div>

<details>
<summary><strong>ZC1054</strong>: Use POSIX classes in regex/glob <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

When specifying character sets in regular expressions or glob patterns (e.g., with `grep`, `sed`, or Zsh's extended globbing), it's generally more portable and often clearer to use POSIX character classes (e.g., `[[:digit:]]`, `[[:alpha:]]`). These classes abstract character ranges based on locale and standards, ensuring consistent behavior across different systems. Using literal ranges like `[0-9]` or `[a-zA-Z]` might not always cover all expected characters in non-ASCII locales and can be less readable than their POSIX equivalents.

### Bad Example

```zsh
grep '[a-zA-Z0-9]' file.txt
ls *.[0-9] # Globbing for files ending in a digit
```

### Good Example

```zsh
grep '[[:alnum:]]' file.txt
ls *.[[:digit:]] # Globbing for files ending in a digit
```

### Configuration

To disable this Kata, add `ZC1054` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1055"></div>

<details>
<summary><strong>ZC1055</strong>: Use `[[ -n/-z ]]` for empty string checks <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

In Zsh (and other shells), the `[[ -n string ]]` and `[[ -z string ]]` constructs are the idiomatic and most reliable ways to check if a string is non-empty or empty, respectively. These operators are designed specifically for string length checks and handle various string contents (including empty strings, zero, or whitespace) predictably. Directly comparing a string to an empty string (e.g., `[[ "$var" = "" ]]`) works, but `-n` and `-z` are more concise and semantically clearer for this specific purpose.

### Bad Example

```zsh
if [[ "$VAR" = "" ]]; then echo "VAR is empty"; fi
if [[ ! "$VAR" = "" ]]; then echo "VAR is not empty"; fi
```

### Good Example

```zsh
if [[ -z "$VAR" ]]; then echo "VAR is empty"; fi
if [[ -n "$VAR" ]]; then echo "VAR is not empty"; fi
```

### Configuration

To disable this Kata, add `ZC1055` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1056"></div>

<details>
<summary><strong>ZC1056</strong>: Avoid `$((...))` as a statement <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

In Zsh, `$(())` is primarily an *arithmetic expansion* used to substitute the result of a calculation into a command. It is not intended to be used as a standalone statement for its side effects (like variable assignment within the expression) or for its exit status. For C-style arithmetic statements or to leverage the exit status of an arithmetic expression, the `((...))` construct should be used. Using `$(())` as a statement can lead to confusion or subtle bugs, as its primary purpose is value substitution.

### Bad Example

```zsh
$(( i++ )) # Result is substituted but not used, might be confusing
$(( var = 1 + 2 )) # Assignment inside is typically not useful as a statement
```

### Good Example

```zsh
(( i++ )) # Correctly performs increment as a statement
result=$(( i + 1 )) # Correctly uses arithmetic expansion for substitution
```

### Configuration

To disable this Kata, add `ZC1056` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1057"></div>

<details>
<summary><strong>ZC1057</strong>: Avoid `ls` in assignments <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Assigning the output of `ls` to a variable (e.g., `files=$(ls)`) is unreliable and problematic. Similar to iterating over `ls` output (ZC1050), the output of `ls` is meant for human consumption and is not safely parsable by scripts, especially when filenames contain spaces, newlines, or special characters. These characters will cause word splitting and globbing issues when the variable is later expanded, leading to incorrect file handling. For safely storing lists of filenames, consider using arrays populated by globbing (`files=(*.txt)`) or by using `find` with null-delimited output.

### Bad Example

```zsh
files=$(ls)
for f in $files; do echo "$f"; done # Unreliable if filenames have spaces
```

### Good Example

```zsh
files=(*) # Populates array with filenames from current directory
for f in "${files[@]}"; do echo "$f"; done

## For specific patterns
txt_files=(*.txt)
for f in "${txt_files[@]}"; do echo "$f"; done
```

### Configuration

To disable this Kata, add `ZC1057` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1058"></div>

<details>
<summary><strong>ZC1058</strong>: Avoid `sudo` with redirection <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Using `sudo command > file` for redirection is often problematic because the redirection (`> file`) is performed by the *unprivileged shell* running `sudo`, not by the `command` itself (which runs with elevated privileges). This means the unprivileged shell attempts to write to `file`, which will likely fail if `file` is in a restricted location. To write to a privileged location with `sudo`, the entire redirection operation needs to be executed by the privileged shell, typically by piping to `sudo tee` or using `sudo sh -c '...'`.

### Bad Example

```zsh
echo "content" > /etc/privileged_file # Fails if /etc is root-owned
sudo echo "content" > /etc/privileged_file # Still fails, redirection is unprivileged
```

### Good Example

```zsh
echo "content" | sudo tee /etc/privileged_file # Pipe to privileged tee
## Or for more complex operations:
sudo sh -c 'echo "content" > /etc/privileged_file' # Entire operation is privileged
```

### Configuration

To disable this Kata, add `ZC1058` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1059"></div>

<details>
<summary><strong>ZC1059</strong>: Use `${var:?}` for `rm` arguments <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

When `rm` is used with a variable that might expand to an empty string, it can unexpectedly try to delete the current directory or other critical locations, leading to data loss. The `${var:?message}` parameter expansion in Zsh (and POSIX shells) provides a crucial safeguard: if `var` is null or unset, the shell prints `message` to standard error and exits (or returns from the function, depending on context). This prevents `rm` from executing with an empty argument and ensures that important variables are defined before destructive operations.

### Bad Example

```zsh
dir_to_delete="$EMPTY_VAR"
rm -rf "$dir_to_delete"/temp_files # If EMPTY_VAR is unset, becomes "rm -rf /temp_files"
```

### Good Example

```zsh
dir_to_delete="${EMPTY_VAR:?Error: Directory variable is not set}"
rm -rf "$dir_to_delete"/temp_files # Exits if EMPTY_VAR is unset
```

### Configuration

To disable this Kata, add `ZC1059` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1060"></div>

<details>
<summary><strong>ZC1060</strong>: Avoid `ps | grep` without exclusion <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Using `ps | grep pattern` to find a specific process can lead to a "grep-for-grep" problem: `grep` itself might appear in the `ps` output, causing a false positive. This is particularly common if the pattern being searched for is part of the `grep` command line. To avoid this, it's best practice to exclude the `grep` process itself from the results. Common methods include using `grep -v grep` or a regex that makes the `grep` pattern unique (e.g., `[p]attern`).

### Bad Example

```zsh
ps aux | grep "my_script" # May match the grep command itself
```

### Good Example

```zsh
ps aux | grep "[m]y_script" # Excludes "grep my_script"
## Or:
ps aux | grep "my_script" | grep -v grep
```

### Configuration

To disable this Kata, add `ZC1060` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1061"></div>

<details>
<summary><strong>ZC1061</strong>: Prefer `{start..end}` over `seq` <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

In Zsh (and Bash 4+), brace expansion with ranges (e.g., `{1..10}`) is a highly efficient and idiomatic way to generate sequences of numbers or characters. This is a shell built-in feature, avoiding the need to invoke an external utility like `seq`. Using `seq` incurs the overhead of launching an external process, which can be slower, especially in loops. Brace expansion is generally faster, more convenient, and promotes native shell feature usage.

### Bad Example

```zsh
for i in $(seq 1 5); do echo $i; done
```

### Good Example

```zsh
for i in {1..5}; do echo $i; done
```

### Configuration

To disable this Kata, add `ZC1061` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1062"></div>

<details>
<summary><strong>ZC1062</strong>: Prefer `grep -E` over `egrep` <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

The `egrep` command is largely deprecated. Its functionality, which is to interpret patterns as extended regular expressions (EREs), is fully encompassed by `grep -E`. Using `grep -E` promotes consistency across commands (`grep`, `grep -F`, `grep -G`), simplifies tooling, and avoids reliance on a separate executable that might not be available or have inconsistent behavior on all systems. This is a matter of modern best practice and command consolidation.

### Bad Example

```zsh
egrep 'word1|word2' file.txt
```

### Good Example

```zsh
grep -E 'word1|word2' file.txt
```

### Configuration

To disable this Kata, add `ZC1062` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1063"></div>

<details>
<summary><strong>ZC1063</strong>: Prefer `grep -F` over `fgrep` <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Similar to `egrep`, the `fgrep` command is largely deprecated. Its functionality, which is to interpret patterns as fixed strings (not regular expressions), is fully encompassed by `grep -F`. Using `grep -F` promotes consistency across `grep` variants, simplifies command recall, and avoids relying on a separate executable that might not be consistently available or behave identically across all systems. This is a best practice for modern shell scripting.

### Bad Example

```zsh
fgrep '$.var' file.txt # Searches for literal "$.var"
```

### Good Example

```zsh
grep -F '$.var' file.txt # Searches for literal "$.var"
```

### Configuration

To disable this Kata, add `ZC1063` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1064"></div>

<details>
<summary><strong>ZC1064</strong>: Prefer `command -v` over `type` <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

In Zsh (and other shells), `command -v` is the most robust and portable way to determine the full pathname or definition of a command, including handling aliases, functions, and built-ins. While `type` can provide similar information, its output format can vary more widely across shells, making it less suitable for programmatic parsing in scripts. `command -v` is specifically designed for reliable output in scripts and adheres to POSIX standards, ensuring consistent behavior.

### Bad Example

```zsh
type my_command # Output format can vary
```

### Good Example

```zsh
command -v my_command # Reliable output for scripting
```

### Configuration

To disable this Kata, add `ZC1064` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1065"></div>

<details>
<summary><strong>ZC1065</strong>: Ensure spaces around `[` and `[[` <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

In both Zsh and other POSIX-like shells, it is a strict requirement to have whitespace (spaces or tabs) around the `[` and `[[` constructs, as well as their closing `]` and `]]`. These are treated as commands or keywords, and without proper spacing, the shell will interpret them as part of adjacent words, leading to syntax errors or unexpected behavior. This is a fundamental rule for shell scripting.

### Bad Example

```zsh
if [condition]; then echo "Bad"; fi
if [[$var = "value"]]; then echo "Also bad"; fi
```

### Good Example

```zsh
if [ condition ]; then echo "Good"; fi
if [[ $var = "value" ]]; then echo "Also good"; fi
```

### Configuration

To disable this Kata, add `ZC1065` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1066"></div>

<details>
<summary><strong>ZC1066</strong>: Avoid iterating over `cat` output <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Iterating over the output of `cat` (e.g., `for line in $(cat file.txt)`) is an unreliable and inefficient practice, often referred to as a "Useless Use of Cat" (UUOC) combined with parsing issues. The `for` loop, by default, splits input on whitespace, meaning that lines containing spaces will be broken into multiple items, and empty lines might be skipped. Additionally, `cat` adds an unnecessary process to the pipeline. For iterating over lines, `while IFS= read -r` loop is the robust method. For iterating over filenames, globbing or `find -print0` is preferred.

### Bad Example

```zsh
for item in $(cat my_list.txt); do
    echo "Item: $item"
done
```

### Good Example

```zsh
while IFS= read -r line; do
    echo "Line: $line"
done < my_list.txt
```

### Configuration

To disable this Kata, add `ZC1066` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1067"></div>

<details>
<summary><strong>ZC1067</strong>: Separate `export` and assignment to avoid masking return codes <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

When you combine variable assignment and `export` on the same line (e.g., `export VAR=$(command)`), the exit status of the `command` might be masked by the exit status of the `export` builtin itself. This can lead to scripts incorrectly assuming success for `command` when it actually failed. To reliably capture the exit status of the command substitution, it's safer to perform the assignment first, then export the variable in a separate step.

### Bad Example

```zsh
export BUILD_ID=$(run_build_process)
if (( $? != 0 )); then echo "Build failed"; fi # $? might be export's exit status
```

### Good Example

```zsh
BUILD_ID=$(run_build_process)
if (( $? != 0 )); then echo "Build failed"; exit 1; fi
export BUILD_ID # Export after checking success
```

### Configuration

To disable this Kata, add `ZC1067` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1068"></div>

<details>
<summary><strong>ZC1068</strong>: Use `add-zsh-hook` instead of defining hook functions directly <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Zsh provides the `add-zsh-hook` function to manage hook functions (like `precmd`, `chpwd`, `preexec`) in a robust and extensible manner. Instead of directly defining a function with the hook's name (e.g., `precmd() { ... }`), which overwrites any existing hook, `add-zsh-hook` allows multiple functions to be registered for the same hook. This ensures compatibility with other Zsh configurations or plugins that might also define hooks, preventing unexpected behavior or lost functionality. It's a best practice for cooperative and modular Zsh configurations.

### Bad Example

```zsh
precmd() {
    echo "Running precmd hook" # Overwrites any previous precmd()
}
```

### Good Example

```zsh
my_precmd_function() {
    echo "Running my custom precmd hook"
}
add-zsh-hook precmd my_precmd_function # Adds to existing hooks
```

### Configuration

To disable this Kata, add `ZC1068` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1069"></div>

<details>
<summary><strong>ZC1069</strong>: Avoid `local` outside of functions <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

The `local` (or `typeset`) keyword is specifically designed to create function-local variables, limiting their scope to the function they are defined within. Using `local` outside of a function (e.g., directly in a script or in a subshell that is not a function body) has no effect on scope or can lead to unexpected behavior, as variables are typically global in script scope. This Kata aims to prevent the misuse of `local` and clarify variable scoping expectations, ensuring that developers correctly apply variable declarations based on their intended scope.

### Bad Example

```zsh
## In a script:
local script_var="hello" # Has no effect on scope
if true; then
    local if_var="world" # No effect, 'if' is not a function scope
fi
```

### Good Example

```zsh
## In a function:
my_func() {
    local func_var="hello" # Correctly scoped to my_func
}
## In a script, simply assign:
script_var="hello"
```

### Configuration

To disable this Kata, add `ZC1069` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1070"></div>

<details>
<summary><strong>ZC1070</strong>: Use `builtin` or `command` to avoid infinite recursion in wrapper functions <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

When you create a Zsh function that has the same name as an existing command (e.g., `ls() { ... }`), calling `ls` inside your function will recursively call your function, leading to infinite recursion and a stack overflow. To call the *original* `ls` command (or any other built-in/external command) from within a function that wraps it, you must use `builtin` or `command`. `builtin ls` executes the built-in version, and `command ls` searches the `PATH` for the external command, bypassing function lookups. This prevents infinite recursion and allows you to augment existing commands safely.

### Bad Example

```zsh
ls() {
    echo "My custom ls"
    ls "$@" # Infinite recursion!
}
```

### Good Example

```zsh
ls() {
    echo "My custom ls"
    command ls "$@" # Calls the external ls command
    # Or if ls is a builtin you want to extend:
    # builtin ls "$@"
}
```

### Configuration

To disable this Kata, add `ZC1070` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1071"></div>

<details>
<summary><strong>ZC1071</strong>: Use `+=` for appending to arrays <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Appending elements to a Zsh array using `array=($array new_element ...)` is generally less efficient and more verbose than using the `+=` operator. The `array=(...)` syntax re-creates the entire array, which can be slower for large arrays, while `array+=(...)` is optimized for appending and is more concise. This Kata promotes the idiomatic use of `+=` for array concatenation.

### Bad Example

```zsh
my_array=(element1 element2)
my_array=($my_array new_element) # Recreates array
```

### Good Example

```zsh
my_array=(element1 element2)
my_array+=(new_element) # Appends efficiently
```

### Configuration

To disable this Kata, add `ZC1071` to the `disabled_katas` list in your `.zshellcheckrc` file.

*Note: This Kata is currently limited in its detection capabilities due to parser limitations with complex array literals. See [Issue #41](https://github.com/afadesigns/zshellcheck/issues/41) for details.*

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1072"></div>

<details>
<summary><strong>ZC1072</strong>: Use `awk` instead of `grep | awk` <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Piping the output of `grep` to `awk` (e.g., `grep pattern file | awk '{...}'`) is often inefficient because `awk` itself has powerful pattern matching capabilities. By combining the pattern matching directly within `awk` (e.g., `awk '/pattern/ { ... }' file`), you eliminate the need for an extra process in the pipeline, reducing overhead and improving performance. This is a common optimization for shell scripts.

### Bad Example

```zsh
grep "error" log.txt | awk '{print $NF}'
```

### Good Example

```zsh
awk '/error/ {print $NF}' log.txt
```

### Configuration

To disable this Kata, add `ZC1072` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1073"></div>

<details>
<summary><strong>ZC1073</strong>: Unnecessary use of `$` in arithmetic expressions <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

In Zsh's arithmetic contexts, specifically within `((...))` expressions, variables do not require a leading `$` to be evaluated for their numeric value. For example, `(( var > 0 ))` is sufficient; `(( $var > 0 ))` is redundant and can be confusing. The shell automatically dereferences variables within arithmetic contexts. Using `$` unnecessarily can clutter the code and deviate from idiomatic Zsh style. However, special parameters like `$#` (number of arguments) still require the `$` as they are not simple variables.

### Bad Example

```zsh
count=5
(( $count++ )) # $ is unnecessary here
if (( $my_var < 10 )); then echo "Low"; fi
```

### Good Example

```zsh
count=5
(( count++ ))
if (( my_var < 10 )); then echo "Low"; fi
if (( $# > 0 )); then echo "Args exist"; fi # $# still requires $
```

### Configuration

To disable this Kata, add `ZC1073` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>


<div id="zc1074"></div>

<details>
<summary><strong>ZC1074</strong>: Prefer modifiers :h/:t over dirname/basename <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Zsh provides modifiers like `:h` (head/dirname) and `:t` (tail/basename) that are faster and more idiomatic than spawning external commands.

### Bad Example

```zsh
dir=$(dirname $path)
file=$(basename $path)
```

### Good Example

```zsh
dir=${path:h}
file=${path:t}
```

### Configuration

To disable this Kata, add `ZC1074` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>


<div id="zc1075"></div>

<details>
<summary><strong>ZC1075</strong>: Quote variable expansions to prevent globbing <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Unquoted variable expansions in Zsh are subject to globbing (filename generation). If the variable contains characters like `*` or `?`, it might match files unexpectedly. Use quotes `"$var"` to prevent this.

### Bad Example

```zsh
rm $file
ls ${files[1]}
```

### Good Example

```zsh
rm "$file"
ls "${files[1]}"
```

### Configuration

To disable this Kata, add `ZC1075` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>


<div id="zc1076"></div>

<details>
<summary><strong>ZC1076</strong>: Use `autoload -Uz` for lazy loading <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

When using `autoload` to lazy-load functions, it is best practice to use the `-Uz` flags.
*   `-U`: Marks the function as "undefined" and prevents alias expansion during definition.
*   `-z`: Ensures Zsh style autoloading (as opposed to ksh style).

### Bad Example

```zsh
autoload my_func
autoload -U my_func
```

### Good Example

```zsh
autoload -Uz my_func
```

### Configuration

To disable this Kata, add `ZC1076` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>


<div id="zc1077"></div>

<details>
<summary><strong>ZC1077</strong>: Prefer `${var:u/l}` over `tr` for case conversion <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Using `tr` in a pipeline for simple case conversion is slower than using Zsh's built-in parameter expansion flags `:u` (upper) and `:l` (lower).

### Bad Example

```zsh
upper=$(echo $var | tr 'a-z' 'A-Z')
lower=$(echo $var | tr '[:upper:]' '[:lower:]')
```

### Good Example

```zsh
upper=${var:u}
lower=${var:l}
```

### Configuration

To disable this Kata, add `ZC1077` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>


<div id="zc1078"></div>

<details>
<summary><strong>ZC1078</strong>: Quote `$@` and `$*` when passing arguments <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Using unquoted `$@` or `$*` splits arguments by IFS (usually space). Use `"$@"` to preserve the original argument grouping, or `"$*"` to join them into a single string.

### Bad Example

```zsh
my_cmd $@
my_cmd $*
```

### Good Example

```zsh
my_cmd "$@"
my_cmd "$*"
```

### Configuration

To disable this Kata, add `ZC1078` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>


<div id="zc1079"></div>

<details>
<summary><strong>ZC1079</strong>: Quote RHS of `==` in `[[ ... ]]` to prevent pattern matching <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

In `[[ ... ]]`, unquoted variable expansions on the right-hand side of `==` or `!=` are treated as patterns (globbing). If you intend to compare strings literally, quote the variable.

### Bad Example

```zsh
[[ $var == $other ]]  # Matches if $other contains wildcards
[[ $var != $other ]]
```

### Good Example

```zsh
[[ $var == "$other" ]] # Literal string comparison
[[ $var == pattern* ]] # Unquoted literals are fine for patterns
```

### Configuration

To disable this Kata, add `ZC1079` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>


<div id="zc1080"></div>

<details>
<summary><strong>ZC1080</strong>: Use `(N)` nullglob qualifier for globs in loops <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

In Zsh, if a glob matches no files, it throws an error by default (`zsh: no matches found: ...`). When iterating over a glob in a `for` loop, use the `(N)` glob qualifier to allow it to match nothing (nullglob). This prevents the script from crashing or printing an error if the directory is empty.

### Bad Example

```zsh
for f in *.txt; do
  echo "Found $f"
done
```

### Good Example

```zsh
for f in *.txt(N); do
  echo "Found $f"
done
```

### Configuration

To disable this Kata, add `ZC1080` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>


<div id="zc1081"></div>

<details>
<summary><strong>ZC1081</strong>: Use `${#var}` to get string length instead of `wc -c` <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Using `echo $var | wc -c` involves a subshell and external command overhead. Zsh has a built-in operator `${#var}` to get the length of a string instantly.

### Bad Example

```zsh
len=$(echo $var | wc -c)
len=$(print -r $var | wc -m)
```

### Good Example

```zsh
len=${#var}
```

### Configuration

To disable this Kata, add `ZC1081` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>


<div id="zc1082"></div>

<details>
<summary><strong>ZC1082</strong>: Prefer `${var//old/new}` over `sed` for simple replacements <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Using `sed` for simple string replacement is slower than Zsh's built-in parameter expansion. Use `${var/old/new}` (replace first) or `${var//old/new}` (replace all).

### Bad Example

```zsh
new=$(echo $var | sed 's/foo/bar/g')
```

### Good Example

```zsh
new=${var//foo/bar}
```

### Configuration

To disable this Kata, add `ZC1082` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1083"></div>

<details>
<summary><strong>ZC1083</strong>: Brace expansion limits cannot be variables <img src="https://img.shields.io/badge/Severity-error-red?style=flat-square" height="15"/></summary>

### Description

Brace expansion `{x..y}` happens before variable expansion. Using variables as limits does not work as expected -- the braces are treated literally.

### Bad Example

```zsh
n=10
for i in {1..$n}; do echo $i; done
```

### Good Example

```zsh
n=10
for i in $(seq 1 $n); do echo $i; done
## Or use a C-style loop:
for (( i = 1; i <= n; i++ )); do echo $i; done
```

### Configuration

To disable this Kata, add `ZC1083` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1084"></div>

<details>
<summary><strong>ZC1084</strong>: Quote globs in <code>find</code> commands <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

Unquoted globs in `find` commands are expanded by the shell before `find` runs. Quote the pattern so `find` receives it literally.

### Bad Example

```zsh
find . -name *.txt
```

### Good Example

```zsh
find . -name '*.txt'
```

### Configuration

To disable this Kata, add `ZC1084` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1085"></div>

<details>
<summary><strong>ZC1085</strong>: Quote variable expansions in <code>for</code> loops <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

Unquoted variable expansions in `for` loops are split by IFS (usually spaces). Quote variables or use arrays to iterate safely.

### Bad Example

```zsh
files="file one.txt file two.txt"
for f in $files; do echo "$f"; done
```

### Good Example

```zsh
files=("file one.txt" "file two.txt")
for f in "${files[@]}"; do echo "$f"; done
```

### Configuration

To disable this Kata, add `ZC1085` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1086"></div>

<details>
<summary><strong>ZC1086</strong>: Prefer <code>func() { ... }</code> over <code>function func { ... }</code> <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

The `function` keyword is optional in Zsh and non-standard in POSIX sh. Use the `name() { ... }` form for portability and consistency.

### Bad Example

```zsh
function greet {
  echo "Hello"
}
```

### Good Example

```zsh
greet() {
  echo "Hello"
}
```

### Configuration

To disable this Kata, add `ZC1086` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1087"></div>

<details>
<summary><strong>ZC1087</strong>: Output redirection overwrites input file <img src="https://img.shields.io/badge/Severity-error-red?style=flat-square" height="15"/></summary>

### Description

Redirecting output to a file that is also being read as input causes the file to be truncated before it is read. Use a temporary file or `sponge`.

### Bad Example

```zsh
sort file.txt > file.txt
```

### Good Example

```zsh
sort file.txt > file.tmp && mv file.tmp file.txt
## Or with sponge:
sort file.txt | sponge file.txt
```

### Configuration

To disable this Kata, add `ZC1087` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1088"></div>

<details>
<summary><strong>ZC1088</strong>: Subshell isolates state changes <img src="https://img.shields.io/badge/Severity-warning-orange?style=flat-square" height="15"/></summary>

### Description

Commands inside `( ... )` run in a subshell. Variable assignments, `cd`, and other state changes inside a subshell do not affect the parent shell.

### Bad Example

```zsh
(cd /tmp && myvar="changed")
echo $myvar  # Still the old value
```

### Good Example

```zsh
cd /tmp && myvar="changed"
echo $myvar
## Or use braces for grouping without a subshell:
{ cd /tmp && myvar="changed"; }
```

### Configuration

To disable this Kata, add `ZC1088` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1089"></div>

<details>
<summary><strong>ZC1089</strong>: Redirection order matters (<code>2>&1 > file</code>) <img src="https://img.shields.io/badge/Severity-error-red?style=flat-square" height="15"/></summary>

### Description

Redirecting stderr to stdout (`2>&1`) before redirecting stdout to a file (`> file`) means stderr goes to the original stdout (terminal), not the file. Reverse the order.

### Bad Example

```zsh
command 2>&1 > output.log
```

### Good Example

```zsh
command > output.log 2>&1
```

### Configuration

To disable this Kata, add `ZC1089` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1090"></div>

<details>
<summary><strong>ZC1090</strong>: Quoted regex pattern in <code>=~</code> <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

Quoting the pattern on the right side of `=~` forces literal string matching in Zsh/Bash. Leave the regex pattern unquoted for proper regex evaluation.

### Bad Example

```zsh
if [[ $var =~ "^[0-9]+$" ]]; then echo "numeric"; fi
```

### Good Example

```zsh
if [[ $var =~ ^[0-9]+$ ]]; then echo "numeric"; fi
```

### Configuration

To disable this Kata, add `ZC1090` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1091"></div>

<details>
<summary><strong>ZC1091</strong>: Use <code>(( ... ))</code> for arithmetic comparisons in <code>[[...]]</code> <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

The `[[ ... ]]` construct is primarily for string comparisons and file tests. For arithmetic comparisons, use `(( ... ))` which provides natural C-style syntax.

### Bad Example

```zsh
if [[ $a -gt $b ]]; then echo "a is greater"; fi
```

### Good Example

```zsh
if (( a > b )); then echo "a is greater"; fi
```

### Configuration

To disable this Kata, add `ZC1091` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1092"></div>

<details>
<summary><strong>ZC1092</strong>: Arrays cannot be exported <img src="https://img.shields.io/badge/Severity-warning-orange?style=flat-square" height="15"/></summary>

### Description

In Zsh, arrays cannot be exported to child processes via `export`. Serialize arrays to a string when passing data to subprocesses.

### Bad Example

```zsh
export my_array=(one two three)
```

### Good Example

```zsh
my_array=(one two three)
export MY_ARRAY_STR="${(j/:/)my_array}"
```

### Configuration

To disable this Kata, add `ZC1092` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>







<div id="zc1093"></div>

<details>
<summary><strong>ZC1093</strong>: Use `IFS=` with `read` in loops <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

When reading lines in a loop using `while read ...`, the default behavior of `read` is to trim leading and trailing whitespace from the input line. To preserve the whitespace (read the line exactly as is), you should set `IFS` to an empty string for the `read` command.

### Bad Example

```zsh
while read -r line; do
  echo "$line"
done < file.txt
```

### Good Example

```zsh
while IFS= read -r line; do
  echo "$line"
done < file.txt
```

### Configuration

To disable this Kata, add `ZC1093` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>


<div id="zc1094"></div>

<details>
<summary><strong>ZC1094</strong>: Consider `grep -P` for advanced regex <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

While `grep -E` provides Extended Regular Expressions, there are situations where more advanced pattern matching capabilities are required. Perl Compatible Regular Expressions (PCRE), enabled by `grep -P`, offer features like lookarounds (positive/negative lookahead/lookbehind), backreferences within alternation, and non-greedy quantifiers. If your regex needs extend beyond what ERE offers, `grep -P` is a powerful alternative.

### Bad Example

```zsh
grep -E '(a|b)c' file.txt
```
*(This example is not "bad" in functionality, but serves to illustrate a basic `-E` usage that *could* be enhanced if advanced PCRE features were desired. The Kata suggests `-P` as a potential improvement path.)*

### Good Example

```zsh
grep -P '(?<=prefix)pattern' file.txt
```
*(This example shows a PCRE-specific feature (lookbehind) demonstrating the added capability of `grep -P`.)*

### Configuration

To disable this Kata, add `ZC1094` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>


<div id="zc1095"></div>

<details>
<summary><strong>ZC1095</strong>: Quote here-string content <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

When providing content to a here-string (`<<<`), it is generally best practice to quote the content. Unquoted strings can be subject to word splitting and globbing, leading to unexpected behavior if the string contains spaces, special characters, or patterns that could be expanded by the shell. Quoting ensures that the content is passed literally, preserving its exact form.

Use single quotes (`'...'`) for literal strings, dollar-single quotes (`$'...'`) if you need escape sequence interpretation (e.g., `
` for a newline), or double quotes (`"..."`) if you need variable expansion within the string.

### Bad Example

```zsh
grep pattern <<< foo bar
grep pattern <<< $my_var
```

### Good Example

```zsh
grep pattern <<< 'foo bar'
grep pattern <<< "$my_var"
grep pattern <<< $'	foo
bar'
```

### Configuration

To disable this Kata, add `ZC1095` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1096"></div>

<details>
<summary><strong>ZC1096</strong>: Avoid using `test -e` or `[ -e ... ]` for file existence checks <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Prefer more specific checks like `[[ -f ... ]]` for regular files or `[[ -d ... ]]` for directories. `test -e` is less specific and can mask intentions.

### Bad Example

```zsh
if [ -e "$file" ]; then ... fi
```

### Good Example

```zsh
if [[ -f "$file" ]]; then ... fi
if [[ -d "$dir" ]]; then ... fi
```

### Configuration

To disable this Kata, add `ZC1096` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1097"></div>

<details>
<summary><strong>ZC1097</strong>: Declare loop variables as <code>local</code> in functions <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

Loop variables in `for` loops are global by default in Zsh functions. Declare them as `local` to prevent leaking state to the caller.

### Bad Example

```zsh
my_func() {
  for item in a b c; do
    echo $item
  done
}
```

### Good Example

```zsh
my_func() {
  local item
  for item in a b c; do
    echo $item
  done
}
```

### Configuration

To disable this Kata, add `ZC1097` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1098"></div>

<details>
<summary><strong>ZC1098</strong>: Use `(q)` flag for quoting variables in `eval` <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

When constructing a command string for `eval`, use the `(q)` flag (or `(qq)`, `(q-)`) to safely quote variables and prevent command injection.

### Bad Example

```zsh
eval "ls $dir"
```

### Good Example

```zsh
eval "ls ${(q)dir}"
```

### Configuration

To disable this Kata, add `ZC1098` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1099"></div>

<details>
<summary><strong>ZC1099</strong>: Use `(f)` flag to split lines instead of `while read` <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Zsh provides the `(f)` parameter expansion flag to split a string into lines. Iterating over `${(f)variable}` is often cleaner and faster than piping to `while read`.

### Bad Example

```zsh
cat file | while read line; do ...; done
```

### Good Example

```zsh
for line in ${(f)"$(<file)"}; do ...; done
```

### Configuration

To disable this Kata, add `ZC1099` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1100"></div>

<details>
<summary><strong>ZC1100</strong>: Use parameter expansion instead of `dirname`/`basename` <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Zsh parameter expansion `${var%/*}` (dirname) and `${var##*/}` (basename) avoid spawning external processes for simple path manipulation.

### Bad Example

```zsh
dir=$(dirname "$filepath")
name=$(basename "$filepath")
```

### Good Example

```zsh
dir=${filepath%/*}
name=${filepath##*/}
```

### Configuration

To disable this Kata, add `ZC1100` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1101"></div>

<details>
<summary><strong>ZC1101</strong>: Use `$(( ))` instead of `bc` for simple arithmetic <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Zsh supports arithmetic expansion with `$(( ))` and floating point via `zmodload zsh/mathfunc`. Avoid piping to `bc` for simple calculations.

### Bad Example

```zsh
result=$(echo "1.5 + 2.5" | bc)
```

### Good Example

```zsh
result=$(( 1.5 + 2.5 ))
```

### Configuration

To disable this Kata, add `ZC1101` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1102"></div>

<details>
<summary><strong>ZC1102</strong>: Redirecting output of `sudo` does not work as expected <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Redirections are performed by the current shell before `sudo` is started. So `sudo echo > /root/file` will try to open `/root/file` as the current user, failing.

### Bad Example

```zsh
sudo echo "content" > /root/file
```

### Good Example

```zsh
echo "content" | sudo tee /root/file > /dev/null
```

### Configuration

To disable this Kata, add `ZC1102` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1103"></div>

<details>
<summary><strong>ZC1103</strong>: Suggest `path` array instead of `$PATH` string manipulation <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Zsh automatically maps the `$PATH` environment variable to the `$path` array. Modifying `$path` is cleaner and less error-prone than manipulating the colon-separated `$PATH` string.

### Bad Example

```zsh
PATH="$PATH:/usr/local/bin"
```

### Good Example

```zsh
path+=('/usr/local/bin')
```

### Configuration

To disable this Kata, add `ZC1103` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1104"></div>

<details>
<summary><strong>ZC1104</strong>: Suggest `path` array instead of `export PATH` string manipulation <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Zsh automatically maps the `$PATH` environment variable to the `$path` array. Modifying `$path` is cleaner and less error-prone than manipulating the colon-separated `$PATH` string via `export PATH=...`.

### Bad Example

```zsh
export PATH=$PATH:/usr/local/bin
```

### Good Example

```zsh
path+=('/usr/local/bin')
```

### Configuration

To disable this Kata, add `ZC1104` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1105"></div>

<details>
<summary><strong>ZC1105</strong>: Avoid nested arithmetic expansions for clarity <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

While Zsh supports nested arithmetic expansions like `(( $((...)) ))`, they can make code harder to read and reason about. Prefer flatter expressions or temporary variables for intermediate results to improve clarity.

### Bad Example

```zsh
(( result = $(( a + b )) * c ))
```

### Good Example

```zsh
(( sum = a + b ))
(( result = sum * c ))
```

### Configuration

To disable this Kata, add `ZC1105` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1106"></div>

<details>
<summary><strong>ZC1106</strong>: Avoid `set -x` in production scripts <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Using `set -x` (xtrace) in production environments can expose sensitive information, such as API keys or passwords, in logs. While useful for debugging, it should be avoided in production. Consider using targeted debugging or secure logging.

### Bad Example

```zsh
set -x
curl -H "Authorization: Bearer $TOKEN" https://api.example.com
```

### Good Example

```zsh
# Use targeted debugging or secure logging
print -r -- "Making API request to api.example.com"
curl -s -H "Authorization: Bearer $TOKEN" https://api.example.com
```

### Configuration

To disable this Kata, add `ZC1106` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1107"></div>

<details>
<summary><strong>ZC1107</strong>: Use `(( ... ))` for arithmetic conditions <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

In Zsh, arithmetic comparisons using `[[ ... -gt ... ]]` or `[ ... -eq ... ]` are valid but less idiomatic and slower than `(( ... ))`. The double parenthesis syntax supports standard math operators (`>`, `<`, `==`, `!=`) and is optimized.

### Bad Example

```zsh
if [ "$a" -eq "$b" ]; then echo "Equal"; fi
if [[ "$count" -gt 10 ]]; then echo "Greater"; fi
```

### Good Example

```zsh
if (( a == b )); then echo "Equal"; fi
if (( count > 10 )); then echo "Greater"; fi
```

### Configuration

To disable this Kata, add `ZC1107` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1108"></div>

<details>
<summary><strong>ZC1108</strong>: Use Zsh case conversion instead of `tr` <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Zsh provides `${(U)var}` for uppercase and `${(L)var}` for lowercase. Avoid piping through `tr '[:lower:]' '[:upper:]'` for simple case conversion.

### Bad Example

```zsh
upper=$(echo "$var" | tr '[:lower:]' '[:upper:]')
lower=$(echo "$var" | tr 'A-Z' 'a-z')
```

### Good Example

```zsh
upper=${(U)var}
lower=${(L)var}
```

### Configuration

To disable this Kata, add `ZC1108` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1109"></div>

<details>
<summary><strong>ZC1109</strong>: Use parameter expansion instead of `cut` for field extraction <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

For simple field extraction from variables, use Zsh parameter expansion like `${var%%:*}` or `${(s.:.)var}` instead of piping through `cut`.

### Bad Example

```zsh
field=$(echo "$line" | cut -d: -f1)
```

### Good Example

```zsh
field=${line%%:*}
# Or split into array:
fields=(${(s.:.)line})
```

### Configuration

To disable this Kata, add `ZC1109` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1110"></div>

<details>
<summary><strong>ZC1110</strong>: Use Zsh subscripts instead of `head -1` or `tail -1` <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Zsh array subscripts `${lines[1]}` and `${lines[-1]}` can extract the first or last element without spawning `head` or `tail` as external processes.

### Bad Example

```zsh
first=$(echo "$data" | head -1)
last=$(echo "$data" | tail -1)
```

### Good Example

```zsh
lines=(${(f)data})
first=${lines[1]}
last=${lines[-1]}
```

### Configuration

To disable this Kata, add `ZC1110` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1111"></div>

<details>
<summary><strong>ZC1111</strong>: Avoid `xargs` for simple command invocation <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Zsh can iterate arrays directly with `for` loops or use `${(f)...}` to split command output by newlines. Avoid `xargs` when processing lines one at a time.

### Bad Example

```zsh
find . -name "*.txt" | xargs rm
```

### Good Example

```zsh
for f in ${(f)$(find . -name "*.txt")}; do
    rm "$f"
done
# Or better: rm **/*.txt(N)
```

### Configuration

To disable this Kata, add `ZC1111` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1112"></div>

<details>
<summary><strong>ZC1112</strong>: Avoid `grep -c` -- use Zsh pattern matching for counting <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

For counting matches in a variable, use Zsh `${#${(f)...}}` or array filtering with `${(M)array:#pattern}` instead of piping through `grep -c`.

### Bad Example

```zsh
count=$(echo "$data" | grep -c "pattern")
```

### Good Example

```zsh
lines=(${(f)data})
matches=(${(M)lines:#*pattern*})
count=${#matches}
```

### Configuration

To disable this Kata, add `ZC1112` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1113"></div>

<details>
<summary><strong>ZC1113</strong>: Use `${var:A}` instead of `realpath` or `readlink -f` <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Zsh provides the `:A` modifier to resolve a path to its absolute form, following symlinks. Avoid spawning `realpath` or `readlink -f` as external processes.

### Bad Example

```zsh
abs=$(realpath "$path")
abs=$(readlink -f "$path")
```

### Good Example

```zsh
abs=${path:A}
```

### Configuration

To disable this Kata, add `ZC1113` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1114"></div>

<details>
<summary><strong>ZC1114</strong>: Consider Zsh `=(...)` for temporary files <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Zsh `=(cmd)` creates a temporary file with the command output that is automatically cleaned up. Consider this instead of manual `mktemp` and cleanup patterns.

### Bad Example

```zsh
tmpfile=$(mktemp)
cmd > "$tmpfile"
diff "$tmpfile" other_file
rm "$tmpfile"
```

### Good Example

```zsh
diff =(cmd) other_file
```

### Configuration

To disable this Kata, add `ZC1114` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1115"></div>

<details>
<summary><strong>ZC1115</strong>: Use Zsh string manipulation instead of `rev` <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Zsh can reverse strings using parameter expansion. Avoid spawning `rev` as an external process for simple string reversal.

### Bad Example

```zsh
reversed=$(echo "$string" | rev)
```

### Good Example

```zsh
reversed=${(j::)${(Oas::)string}}
```

### Configuration

To disable this Kata, add `ZC1115` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1116"></div>

<details>
<summary><strong>ZC1116</strong>: Use Zsh multios instead of `tee` <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Zsh `setopt multios` allows redirecting output to multiple files with `cmd > file1 > file2`. Avoid spawning `tee` for simple output duplication.

### Bad Example

```zsh
cmd | tee output.log
cmd | tee file1 file2
```

### Good Example

```zsh
setopt multios
cmd > file1 > file2
cmd > output.log > /dev/tty
```

### Configuration

To disable this Kata, add `ZC1116` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1117"></div>

<details>
<summary><strong>ZC1117</strong>: Use `&!` or `disown` instead of `nohup` <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Zsh provides `&!` (shorthand for `& disown`) to run a command in the background immune to hangups. Avoid spawning `nohup` as an external process.

### Bad Example

```zsh
nohup long_running_cmd &
```

### Good Example

```zsh
long_running_cmd &!
# Or equivalently:
long_running_cmd & disown
```

### Configuration

To disable this Kata, add `ZC1117` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1118"></div>

<details>
<summary><strong>ZC1118</strong>: Use `print -rn` instead of `echo -n` <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

The behavior of `echo -n` varies across shells and platforms. In Zsh, `print -rn` is the reliable way to output text without a trailing newline.

### Bad Example

```zsh
echo -n "Enter your name: "
```

### Good Example

```zsh
print -rn "Enter your name: "
```

### Configuration

To disable this Kata, add `ZC1118` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1119"></div>

<details>
<summary><strong>ZC1119</strong>: Use `$EPOCHSECONDS` instead of `date +%s` <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Zsh provides `$EPOCHSECONDS` and `$EPOCHREALTIME` via `zsh/datetime` module. Avoid spawning `date` for simple Unix timestamp retrieval.

### Bad Example

```zsh
timestamp=$(date +%s)
```

### Good Example

```zsh
zmodload zsh/datetime
timestamp=$EPOCHSECONDS
# For sub-second precision:
timestamp=$EPOCHREALTIME
```

### Configuration

To disable this Kata, add `ZC1119` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1120"></div>

<details>
<summary><strong>ZC1120</strong>: Use `$PWD` instead of `pwd` <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" height="15"/></summary>

### Description

Zsh maintains `$PWD` as a built-in variable tracking the current directory. Avoid spawning `pwd` as an external process.

### Bad Example

```zsh
current=$(pwd)
cd "$(pwd)/subdir"
```

### Good Example

```zsh
current=$PWD
cd "$PWD/subdir"
```

### Configuration

To disable this Kata, add `ZC1120` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[⬆ Back to Top](#table-of-contents)
</details>

<div id="zc1121"></div>

<details>
<summary><strong>ZC1121</strong>: Use <code>$HOST</code> instead of <code>hostname</code> <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

Zsh provides `$HOST` as a built-in variable containing the hostname. Avoid spawning the `hostname` command.

### Bad Example

```zsh
host=$(hostname)
```

### Good Example

```zsh
host=$HOST
```

### Configuration

To disable this Kata, add `ZC1121` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1122"></div>

<details>
<summary><strong>ZC1122</strong>: Use <code>$USER</code> instead of <code>whoami</code> <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

Zsh provides `$USER` as a built-in variable containing the current username. Avoid spawning the `whoami` command.

### Bad Example

```zsh
me=$(whoami)
```

### Good Example

```zsh
me=$USER
```

### Configuration

To disable this Kata, add `ZC1122` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1123"></div>

<details>
<summary><strong>ZC1123</strong>: Use <code>$OSTYPE</code> instead of <code>uname</code> <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

Zsh provides `$OSTYPE` (e.g., `linux-gnu`, `darwin`) as a built-in variable. Avoid spawning the `uname` command for OS detection.

### Bad Example

```zsh
os=$(uname -s)
```

### Good Example

```zsh
case $OSTYPE in
  linux*) echo "Linux" ;;
  darwin*) echo "macOS" ;;
esac
```

### Configuration

To disable this Kata, add `ZC1123` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1124"></div>

<details>
<summary><strong>ZC1124</strong>: Use <code>: > file</code> instead of <code>cat /dev/null > file</code> to truncate <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

Truncating a file with `cat /dev/null > file` spawns an unnecessary process. Use `: > file` or `> file` instead.

### Bad Example

```zsh
cat /dev/null > logfile.txt
```

### Good Example

```zsh
: > logfile.txt
## Or simply:
> logfile.txt
```

### Configuration

To disable this Kata, add `ZC1124` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1125"></div>

<details>
<summary><strong>ZC1125</strong>: Avoid <code>echo | grep</code> for string matching <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

Using `echo $var | grep pattern` spawns two unnecessary processes. Use Zsh's built-in pattern matching with `[[ $var == *pattern* ]]` instead.

### Bad Example

```zsh
if echo "$var" | grep -q "error"; then
  echo "found"
fi
```

### Good Example

```zsh
if [[ $var == *error* ]]; then
  echo "found"
fi
```

### Configuration

To disable this Kata, add `ZC1125` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1126"></div>

<details>
<summary><strong>ZC1126</strong>: Use <code>sort -u</code> instead of <code>sort | uniq</code> <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

`sort | uniq` spawns two processes when `sort -u` does the same in one.

### Bad Example

```zsh
cat data.txt | sort | uniq
```

### Good Example

```zsh
sort -u data.txt
```

### Configuration

To disable this Kata, add `ZC1126` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1127"></div>

<details>
<summary><strong>ZC1127</strong>: Avoid <code>ls</code> for counting files <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

Using `ls | wc -l` to count files spawns unnecessary processes and fails with filenames containing newlines. Use Zsh glob qualifiers instead.

### Bad Example

```zsh
count=$(ls *.txt | wc -l)
```

### Good Example

```zsh
files=(*.txt(N))
count=${#files}
```

### Configuration

To disable this Kata, add `ZC1127` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1128"></div>

<details>
<summary><strong>ZC1128</strong>: Use <code>> file</code> instead of <code>touch file</code> for creation <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

If the goal is to create an empty file, `> file` does it without spawning an external process. Use `touch` only when you need to update timestamps.

### Bad Example

```zsh
touch newfile.txt
```

### Good Example

```zsh
> newfile.txt
```

### Configuration

To disable this Kata, add `ZC1128` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1129"></div>

<details>
<summary><strong>ZC1129</strong>: Use Zsh <code>stat</code> module instead of <code>wc -c</code> for file size <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

Zsh's `zstat` (via `zmodload zsh/stat`) provides file size without spawning external processes like `wc -c`.

### Bad Example

```zsh
size=$(wc -c < file.txt)
```

### Good Example

```zsh
zmodload zsh/stat
zstat -A size +size file.txt
```

### Configuration

To disable this Kata, add `ZC1129` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1131"></div>

<details>
<summary><strong>ZC1131</strong>: Avoid <code>cat file | while read</code> -- use redirection <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

`cat file | while read line` spawns an unnecessary `cat` process and runs the loop in a subshell. Use input redirection instead.

### Bad Example

```zsh
cat file.txt | while read line; do
  echo "$line"
done
```

### Good Example

```zsh
while read line; do
  echo "$line"
done < file.txt
```

### Configuration

To disable this Kata, add `ZC1131` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1132"></div>

<details>
<summary><strong>ZC1132</strong>: Use Zsh pattern extraction instead of <code>grep -o</code> <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

For extracting matching parts from variables, use Zsh `${(M)var:#pattern}` or parameter expansion instead of spawning `grep -o`.

### Bad Example

```zsh
result=$(echo "$var" | grep -o '[0-9]*')
```

### Good Example

```zsh
if [[ $var =~ '([0-9]+)' ]]; then
  result=$match[1]
fi
```

### Configuration

To disable this Kata, add `ZC1132` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1133"></div>

<details>
<summary><strong>ZC1133</strong>: Avoid <code>kill -9</code> -- use <code>kill</code> first, then escalate <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

`kill -9` (SIGKILL) cannot be caught or ignored. Always try `kill` (SIGTERM) first to allow processes to clean up gracefully, then escalate if needed.

### Bad Example

```zsh
kill -9 $pid
```

### Good Example

```zsh
kill $pid
sleep 2
kill -0 $pid 2>/dev/null && kill -9 $pid
```

### Configuration

To disable this Kata, add `ZC1133` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1134"></div>

<details>
<summary><strong>ZC1134</strong>: Avoid <code>sleep</code> in tight loops <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

Using `sleep` inside a loop for polling creates busy-wait patterns. Consider using inotifywait, signals, or other event-driven approaches.

### Bad Example

```zsh
while true; do
  check_status && break
  sleep 1
done
```

### Good Example

```zsh
inotifywait -e modify /path/to/file
## Or use a reasonable interval with a maximum retry count
```

### Configuration

To disable this Kata, add `ZC1134` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1135"></div>

<details>
<summary><strong>ZC1135</strong>: Avoid <code>env VAR=val cmd</code> -- use inline assignment <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

Zsh supports inline environment variable assignment with `VAR=val cmd`. The `env` command is unnecessary.

### Bad Example

```zsh
env LANG=C sort file.txt
```

### Good Example

```zsh
LANG=C sort file.txt
```

### Configuration

To disable this Kata, add `ZC1135` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1136"></div>

<details>
<summary><strong>ZC1136</strong>: Avoid <code>rm -rf</code> without safeguard <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

`rm -rf` with a variable path is dangerous if the variable is empty. Always validate the variable or use `${var:?}` to abort on empty values.

### Bad Example

```zsh
rm -rf "$dir/"
```

### Good Example

```zsh
rm -rf "${dir:?'dir is not set'}/"
```

### Configuration

To disable this Kata, add `ZC1136` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1137"></div>

<details>
<summary><strong>ZC1137</strong>: Avoid hardcoded <code>/tmp</code> paths <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

Hardcoded `/tmp` paths are predictable and may cause race conditions or security issues. Use `mktemp` for safe temporary file creation.

### Bad Example

```zsh
echo "data" > /tmp/myapp.log
```

### Good Example

```zsh
tmpfile=$(mktemp)
echo "data" > "$tmpfile"
```

### Configuration

To disable this Kata, add `ZC1137` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1139"></div>

<details>
<summary><strong>ZC1139</strong>: Avoid <code>source</code> with URL -- use local files <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

Sourcing scripts from URLs (curl | source) is a security risk. Download first, inspect, then source from a local file.

### Bad Example

```zsh
source <(curl -s https://example.com/script.zsh)
```

### Good Example

```zsh
curl -o script.zsh https://example.com/script.zsh
## Inspect the script, then:
source script.zsh
```

### Configuration

To disable this Kata, add `ZC1139` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1140"></div>

<details>
<summary><strong>ZC1140</strong>: Use <code>command -v</code> instead of <code>hash</code> for command existence <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

`hash cmd` is a POSIX way to check command existence but provides less clear output. `command -v` is more portable and explicit.

### Bad Example

```zsh
if hash git 2>/dev/null; then echo "git found"; fi
```

### Good Example

```zsh
if command -v git &>/dev/null; then echo "git found"; fi
```

### Configuration

To disable this Kata, add `ZC1140` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1141"></div>

<details>
<summary><strong>ZC1141</strong>: Avoid <code>curl | sh</code> pattern <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

Piping curl output to sh/bash/zsh is a security risk. Download first, inspect the script, then execute it.

### Bad Example

```zsh
curl -s https://example.com/install.sh | sh
```

### Good Example

```zsh
curl -o install.sh https://example.com/install.sh
less install.sh  # Inspect first
sh install.sh
```

### Configuration

To disable this Kata, add `ZC1141` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1142"></div>

<details>
<summary><strong>ZC1142</strong>: Avoid chained <code>grep | grep</code> -- combine patterns <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

Chaining `grep pattern1 | grep pattern2` spawns multiple processes. Combine patterns with `grep -E 'p1.*p2'` or use `awk`.

### Bad Example

```zsh
grep "error" log.txt | grep "fatal"
```

### Good Example

```zsh
grep -E "error.*fatal" log.txt
## Or for independent patterns:
grep -E "error|fatal" log.txt
```

### Configuration

To disable this Kata, add `ZC1142` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1143"></div>

<details>
<summary><strong>ZC1143</strong>: Avoid <code>set -e</code> -- use explicit error handling <img src="https://img.shields.io/badge/Severity-info-yellow?style=flat-square" height="15"/></summary>

### Description

`set -e` (errexit) has surprising behavior in Zsh with conditionals, subshells, and pipelines. Prefer explicit error handling with `||`, `&&`, or `trap ERR`.

### Bad Example

```zsh
set -e
cd /nonexistent  # Script exits here
```

### Good Example

```zsh
cd /nonexistent || { echo "Failed to cd" >&2; return 1; }
```

### Configuration

To disable this Kata, add `ZC1143` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1144"></div>

<details>
<summary><strong>ZC1144</strong>: Avoid <code>trap</code> with signal numbers -- use names <img src="https://img.shields.io/badge/Severity-info-yellow?style=flat-square" height="15"/></summary>

### Description

Signal numbers vary across platforms. Use signal names like `SIGTERM`, `SIGINT`, `SIGHUP` for portability and readability.

### Bad Example

```zsh
trap 'cleanup' 2 15
```

### Good Example

```zsh
trap 'cleanup' INT TERM
```

### Configuration

To disable this Kata, add `ZC1144` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1145"></div>

<details>
<summary><strong>ZC1145</strong>: Avoid <code>tr -d</code> for character deletion -- use parameter expansion <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

For simple character deletion from variables, use Zsh `${var//char/}` instead of spawning `tr`.

### Bad Example

```zsh
clean=$(echo "$var" | tr -d ' ')
```

### Good Example

```zsh
clean=${var// /}
```

### Configuration

To disable this Kata, add `ZC1145` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1146"></div>

<details>
<summary><strong>ZC1146</strong>: Avoid <code>cat file | awk</code> -- pass file to awk directly <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

`cat file | awk` spawns an unnecessary `cat` process. Pass the file directly to awk.

### Bad Example

```zsh
cat data.txt | awk '{print $1}'
```

### Good Example

```zsh
awk '{print $1}' data.txt
```

### Configuration

To disable this Kata, add `ZC1146` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1147"></div>

<details>
<summary><strong>ZC1147</strong>: Avoid <code>mkdir</code> without <code>-p</code> for nested paths <img src="https://img.shields.io/badge/Severity-info-yellow?style=flat-square" height="15"/></summary>

### Description

Using `mkdir` without `-p` fails if parent directories do not exist. Use `mkdir -p` to create nested directory structures safely.

### Bad Example

```zsh
mkdir /path/to/nested/dir
```

### Good Example

```zsh
mkdir -p /path/to/nested/dir
```

### Configuration

To disable this Kata, add `ZC1147` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1148"></div>

<details>
<summary><strong>ZC1148</strong>: Use <code>compdef</code> instead of <code>compctl</code> for completions <img src="https://img.shields.io/badge/Severity-info-yellow?style=flat-square" height="15"/></summary>

### Description

`compctl` is the old Zsh completion system. Use `compdef` with the new completion system (`compsys`) for more powerful and maintainable completions.

### Bad Example

```zsh
compctl -f my_command
```

### Good Example

```zsh
compdef _files my_command
```

### Configuration

To disable this Kata, add `ZC1148` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1149"></div>

<details>
<summary><strong>ZC1149</strong>: Avoid <code>echo</code> for error messages -- use <code>>&2</code> <img src="https://img.shields.io/badge/Severity-info-yellow?style=flat-square" height="15"/></summary>

### Description

Error messages should go to stderr, not stdout. Use `>&2` redirection or `print -u2` so error output does not interfere with piped data.

### Bad Example

```zsh
echo "Error: file not found"
```

### Good Example

```zsh
print -u2 "Error: file not found"
## Or:
echo "Error: file not found" >&2
```

### Configuration

To disable this Kata, add `ZC1149` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1151"></div>

<details>
<summary><strong>ZC1151</strong>: Avoid <code>cat -A</code> -- use <code>print -v</code> or <code>od</code> for non-printable characters <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

`cat -A` shows non-printable characters but varies across platforms. Use `od`, `hexdump`, or Zsh `print -v` for consistent behavior.

### Bad Example

```zsh
cat -A file.txt
```

### Good Example

```zsh
od -c file.txt
## Or for hex:
hexdump -C file.txt
```

### Configuration

To disable this Kata, add `ZC1151` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1152"></div>

<details>
<summary><strong>ZC1152</strong>: Use Zsh PCRE module instead of <code>grep -P</code> <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

`grep -P` (Perl regex) is not available on all platforms (e.g., macOS). Use Zsh's `zmodload zsh/pcre` for portable PCRE support.

### Bad Example

```zsh
grep -P '\d{3}-\d{4}' file.txt
```

### Good Example

```zsh
zmodload zsh/pcre
pcre_compile '\d{3}-\d{4}'
pcre_match "$line"
```

### Configuration

To disable this Kata, add `ZC1152` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1153"></div>

<details>
<summary><strong>ZC1153</strong>: Use <code>cmp -s</code> instead of <code>diff</code> for equality check <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

When only checking if two files are identical (not viewing differences), `cmp -s` is faster and produces no output.

### Bad Example

```zsh
if diff file1 file2 > /dev/null 2>&1; then echo "same"; fi
```

### Good Example

```zsh
if cmp -s file1 file2; then echo "same"; fi
```

### Configuration

To disable this Kata, add `ZC1153` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1154"></div>

<details>
<summary><strong>ZC1154</strong>: Use <code>find -exec {} +</code> instead of <code>find -exec {} \;</code> <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

`find -exec cmd {} \;` runs cmd once per file. Using `+` instead of `\;` batches files into fewer invocations, improving performance.

### Bad Example

```zsh
find . -name '*.log' -exec rm {} \;
```

### Good Example

```zsh
find . -name '*.log' -exec rm {} +
```

### Configuration

To disable this Kata, add `ZC1154` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1155"></div>

<details>
<summary><strong>ZC1155</strong>: Use <code>whence -a</code> instead of <code>which -a</code> <img src="https://img.shields.io/badge/Severity-info-yellow?style=flat-square" height="15"/></summary>

### Description

`which -a` may be an external command on some systems. In Zsh, `whence -a` is a built-in that provides the same functionality.

### Bad Example

```zsh
which -a python
```

### Good Example

```zsh
whence -a python
```

### Configuration

To disable this Kata, add `ZC1155` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1156"></div>

<details>
<summary><strong>ZC1156</strong>: Avoid <code>ln</code> without <code>-s</code> for symlinks <img src="https://img.shields.io/badge/Severity-info-yellow?style=flat-square" height="15"/></summary>

### Description

Hard links (`ln` without `-s`) share inodes and can cause confusion. Prefer symbolic links with `ln -s` unless hard links are specifically required.

### Bad Example

```zsh
ln target linkname
```

### Good Example

```zsh
ln -s target linkname
```

### Configuration

To disable this Kata, add `ZC1156` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1157"></div>

<details>
<summary><strong>ZC1157</strong>: Avoid <code>strings</code> command -- use Zsh expansion <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

The `strings` command extracts printable strings from binaries. For text processing in variables, Zsh parameter expansion `${(ps:\0:)var}` is more efficient.

### Bad Example

```zsh
strings binary_file | grep "version"
```

### Good Example

```zsh
## For variable content:
print -r -- ${(ps:\0:)var}
```

### Configuration

To disable this Kata, add `ZC1157` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1158"></div>

<details>
<summary><strong>ZC1158</strong>: Avoid <code>chown -R</code> without <code>--no-dereference</code> <img src="https://img.shields.io/badge/Severity-warning-orange?style=flat-square" height="15"/></summary>

### Description

`chown -R` follows symlinks by default, potentially changing ownership of files outside the intended directory tree. Use `--no-dereference` or `-h` to operate on symlinks themselves.

### Bad Example

```zsh
chown -R user:group /var/www
```

### Good Example

```zsh
chown -Rh user:group /var/www
```

### Configuration

To disable this Kata, add `ZC1158` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1159"></div>

<details>
<summary><strong>ZC1159</strong>: Avoid <code>tar</code> without explicit compression flag <img src="https://img.shields.io/badge/Severity-info-yellow?style=flat-square" height="15"/></summary>

### Description

Use explicit compression flags (`-z` for gzip, `-j` for bzip2, `-J` for xz) with `tar` instead of relying on auto-detection from file extension.

### Bad Example

```zsh
tar cf archive.tar.gz /path
```

### Good Example

```zsh
tar czf archive.tar.gz /path
```

### Configuration

To disable this Kata, add `ZC1159` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1160"></div>

<details>
<summary><strong>ZC1160</strong>: Prefer <code>curl</code> over <code>wget</code> for portability <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

`wget` is not installed by default on macOS. `curl` is available on both Linux and macOS, making it a more portable choice.

### Bad Example

```zsh
wget https://example.com/file.tar.gz
```

### Good Example

```zsh
curl -LO https://example.com/file.tar.gz
```

### Configuration

To disable this Kata, add `ZC1160` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1161"></div>

<details>
<summary><strong>ZC1161</strong>: Avoid <code>openssl</code> for simple hashing -- use Zsh modules <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

For simple SHA/MD5 hashing, Zsh provides `zmodload zsh/sha256` and related modules. Avoid spawning `openssl` or `sha256sum` for simple hash operations.

### Bad Example

```zsh
hash=$(echo -n "$data" | sha256sum | cut -d' ' -f1)
```

### Good Example

```zsh
zmodload zsh/sha256
hash=$(sha256 -r "$data")
```

### Configuration

To disable this Kata, add `ZC1161` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1162"></div>

<details>
<summary><strong>ZC1162</strong>: Use <code>cp -a</code> instead of <code>cp -r</code> to preserve attributes <img src="https://img.shields.io/badge/Severity-info-yellow?style=flat-square" height="15"/></summary>

### Description

`cp -r` copies recursively but may not preserve permissions, timestamps, or symlinks. Use `cp -a` (archive mode) to preserve all file attributes.

### Bad Example

```zsh
cp -r /src /dest
```

### Good Example

```zsh
cp -a /src /dest
```

### Configuration

To disable this Kata, add `ZC1162` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1163"></div>

<details>
<summary><strong>ZC1163</strong>: Use <code>grep -m 1</code> instead of <code>grep | head -1</code> <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

`grep pattern | head -1` spawns two processes when `grep -m 1` does the same in one, stopping after the first match.

### Bad Example

```zsh
result=$(grep "error" log.txt | head -1)
```

### Good Example

```zsh
result=$(grep -m 1 "error" log.txt)
```

### Configuration

To disable this Kata, add `ZC1163` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1164"></div>

<details>
<summary><strong>ZC1164</strong>: Avoid <code>sed -n 'Np'</code> -- use Zsh array subscript <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

Extracting a specific line with `sed -n 'Np'` spawns a process. Read the file into a Zsh array and use subscript indexing instead.

### Bad Example

```zsh
line=$(sed -n '5p' file.txt)
```

### Good Example

```zsh
lines=("${(@f)$(< file.txt)}")
line=$lines[5]
```

### Configuration

To disable this Kata, add `ZC1164` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1165"></div>

<details>
<summary><strong>ZC1165</strong>: Use Zsh parameter expansion for simple <code>awk</code> field extraction <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

Simple `awk '{print $1}'` or `awk '{print $NF}'` can often be replaced with Zsh parameter expansion, avoiding a process spawn.

### Bad Example

```zsh
first=$(echo "$line" | awk '{print $1}')
```

### Good Example

```zsh
first=${line%% *}
## Or for the last field:
last=${line##* }
```

### Configuration

To disable this Kata, add `ZC1165` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1166"></div>

<details>
<summary><strong>ZC1166</strong>: Avoid <code>grep -i</code> for case-insensitive match -- use <code>(#i)</code> glob flag <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

Zsh provides the `(#i)` glob flag for case-insensitive matching. Use it instead of spawning `grep -i` for variable content.

### Bad Example

```zsh
if echo "$var" | grep -qi "error"; then echo "found"; fi
```

### Good Example

```zsh
if [[ $var == (#i)*error* ]]; then echo "found"; fi
```

### Configuration

To disable this Kata, add `ZC1166` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1167"></div>

<details>
<summary><strong>ZC1167</strong>: Avoid <code>timeout</code> command -- use Zsh <code>TMOUT</code> or <code>zsh/sched</code> <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

`timeout` is not available on all systems (macOS lacks it by default). Use Zsh's `TMOUT` variable or `zsh/sched` module for portable timeouts.

### Bad Example

```zsh
timeout 5 command
```

### Good Example

```zsh
TMOUT=5 read response
## Or use zsh/sched for scheduling
```

### Configuration

To disable this Kata, add `ZC1167` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1168"></div>

<details>
<summary><strong>ZC1168</strong>: Use <code>${(f)...}</code> instead of <code>readarray</code>/<code>mapfile</code> <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

`readarray` and `mapfile` are Bash builtins not available in Zsh. Use Zsh's `${(f)...}` flag to split output into an array by newlines.

### Bad Example

```zsh
readarray -t lines < file.txt
```

### Good Example

```zsh
lines=("${(@f)$(< file.txt)}")
```

### Configuration

To disable this Kata, add `ZC1168` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>

<div id="zc1169"></div>

<details>
<summary><strong>ZC1169</strong>: Avoid <code>install</code> for simple copy+chmod -- use <code>cp</code> then <code>chmod</code> <img src="https://img.shields.io/badge/Severity-style-blue?style=flat-square" height="15"/></summary>

### Description

`install` command is less common and may confuse readers. For clarity, use separate `cp` and `chmod` commands or `install` only in Makefiles.

### Bad Example

```zsh
install -m 755 my_script /usr/local/bin/my_script
```

### Good Example

```zsh
cp my_script /usr/local/bin/my_script
chmod 755 /usr/local/bin/my_script
```

### Configuration

To disable this Kata, add `ZC1169` to the `disabled_katas` list in your `.zshellcheckrc` file.

---

[Back to Top](#table-of-contents)
</details>
