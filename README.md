# ZShellCheck

```
 mmmmmm  mmmm  #             ""#    ""#      mmm  #                    #
     #" #"   " # mm    mmm     #      #    m"   " # mm    mmm    mmm   #   m
   m#   "#mmm  #"  #  #"  #    #      #    #      #"  #  #"  #  #"  "  # m"
  m"        "# #   #  #""""    #      #    #      #   #  #""""  #      #"#
 ##mmmm "mmm#" #   #  "#mm"    "mm    "mm   "mmm" #   #  "#mm"  "#mm"  #  "m
```

`zshellcheck` is a static analysis tool (linter) specifically designed for Zsh scripts. Unlike `shellcheck`, which focuses on POSIX sh/bash compatibility, `zshellcheck` understands Zsh syntax, best practices, and common pitfalls.

It parses Zsh scripts into an Abstract Syntax Tree (AST) and runs a series of checks ("Katas") to identify issues.

## Features

-   **Zsh-Specific Parsing:** Handles Zsh constructs like `[[ ... ]]`, `(( ... ))`, arrays, associative arrays, and modifiers.
-   **Extensible Katas:** Rules are implemented as independent "Katas" that can be easily added or disabled.
-   **Configurable:** Disable specific checks via `.zshellcheckrc` configuration file.
-   **Integration Ready:** Designed to work with `pre-commit` and CI pipelines.

## Installation

### From Source

```bash
go install github.com/afadesigns/zshellcheck/cmd/zshellcheck@latest
```

### Pre-commit Hook

Add this to your `.pre-commit-config.yaml`:

```yaml
-   repo: https://github.com/afadesigns/zshellcheck
    rev: main # or specific tag
    hooks:
    -   id: zshellcheck
```

## Usage

```bash
zshellcheck [flags] <file1.zsh> [file2.zsh]...
```

**Flags:**
- `-format [text|json]`: Output format (default: "text").

## Implemented Checks (Katas)

<details>
<summary>Click to expand full list of implemented checks</summary>

| ID | Title |
| :--- | :--- |
| **ZC1001** | Use ${} for array element access |
| **ZC1002** | Use $(...) instead of backticks |
| **ZC1003** | Use `((...))` for arithmetic comparisons instead of `[` or `test` |
| **ZC1005** | Use whence instead of which |
| **ZC1006** | Prefer [[ over test for tests |
| **ZC1007** | Avoid using `chmod 777` |
| **ZC1008** | Use `\$(())` for arithmetic operations |
| **ZC1009** | Use `((...))` for C-style arithmetic |
| **ZC1010** | Use [[ ... ]] instead of [ ... ] |
| **ZC1011** | Use `git` porcelain commands instead of plumbing commands |
| **ZC1012** | Use `read -r` to prevent backslash escaping |
| **ZC1013** | Use `((...))` for arithmetic operations instead of `let` |
| **ZC1014** | Use `git switch` or `git restore` instead of `git checkout` |
| **ZC1015** | Use `$(...)` for command substitution instead of backticks |
| **ZC1017** | Use `print -r` to print strings literally |
| **ZC1018** | Use `((...))` for C-style arithmetic instead of `expr` |
| **ZC1019** | Use `whence` instead of `which` |
| **ZC1020** | Use `[[ ... ]]` for tests instead of `test` |
| **ZC1021** | Use symbolic permissions with `chmod` instead of octal |
| **ZC1022** | Use `$((...))` for arithmetic expansion |
| **ZC1023** | Use `$((...))` for arithmetic expansion |
| **ZC1024** | Use `$((...))` for arithmetic expansion |
| **ZC1025** | Use `$((...))` for arithmetic expansion |
| **ZC1026** | Use `$((...))` for arithmetic expansion |
| **ZC1027** | Use `$((...))` for arithmetic expansion |
| **ZC1028** | Use `$((...))` for arithmetic expansion |
| **ZC1029** | Use `$((...))` for arithmetic expansion |
| **ZC1030** | Use `printf` instead of `echo` |
| **ZC1031** | Use `#!/usr/bin/env zsh` for portability |
| **ZC1032** | Use `((...))` for C-style incrementing |
| **ZC1033** | Use `$((...))` for arithmetic expansion |
| **ZC1034** | Use `command -v` instead of `which` |
| **ZC1035** | Use `$((...))` for arithmetic expansion |
| **ZC1036** | Prefer `[[ ... ]]` over `test` command |
| **ZC1037** | Use 'print -r --' for variable expansion |
| **ZC1038** | Avoid useless use of cat |
| **ZC1039** | Avoid `rm` with root path |
| **ZC1040** | Use (N) nullglob qualifier for globs in loops |
| **ZC1041** | Do not use variables in printf format string |
| **ZC1042** | Use "$@" to iterate over arguments |
| **ZC1043** | Use `local` for variables in functions |
| **ZC1044** | Check for unchecked `cd` commands |
| **ZC1045** | Declare and assign separately to avoid masking return values |
| **ZC1046** | Avoid `eval` |
| **ZC1047** | Avoid `sudo` in scripts |
| **ZC1048** | Avoid `source` with relative paths |
| **ZC1049** | Prefer functions over aliases |
| **ZC1050** | Avoid iterating over `ls` output |
| **ZC1051** | Quote variables in `rm` to avoid globbing |
| **ZC1052** | Avoid `sed -i` for portability |
| **ZC1053** | Silence `grep` output in conditions |
| **ZC1054** | Use POSIX classes in regex/glob |
| **ZC1055** | Use `[[ -n/-z ]]` for empty string checks |
| **ZC1056** | Avoid `$((...))` as a statement |
| **ZC1057** | Avoid `ls` in assignments |
| **ZC1058** | Avoid `sudo` with redirection |

</details>

## Configuration

Create a `.zshellcheckrc` file in your project root:

```yaml
disabled_katas:
  - ZC1005
  - ZC1011
```