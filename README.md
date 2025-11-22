# zshellcheck

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
<summary>Click to expand list of checks</summary>

| ID | Title | Description |
| :--- | :--- | :--- |
| **ZC1001** | Use `${}` for array access | Enforces `${array[1]}` syntax instead of `$array[1]`. |
| **ZC1002** | Use `$(...)` for substitution | Recommends `$(cmd)` over backticks `` `cmd` ``. |
| **ZC1003** | Use `((...))` for arithmetic | Recommends `(( val > 0 ))` over `[ $val -gt 0 ]` for numeric comparisons. |
| **ZC1005** | Use `whence` instead of `which` | `whence` is the Zsh builtin for locating commands. |
| **ZC1006** | Prefer `[[ ... ]]` over `test` | `[[` is safer and faster than `[` or `test`. |
| **ZC1010** | Use `[[ ... ]]` instead of `[ ... ]` | Same as ZC1006, targeting `[` specifically. |
| **ZC1011** | Avoid `git` plumbing commands | Warns against using unstable plumbing commands in scripts. |
| **ZC1012** | Use `read -r` | Prevents backslash escaping in `read`. |
| **ZC1032** | Use `(( i++ ))` for increment | Recommends C-style increment over `let i=i+1`. |
| **ZC1037** | Use `print -r --` for output | Safer alternative to `echo` for variable expansion. |

</details>

## Configuration

Create a `.zshellcheckrc` file in your project root:

```yaml
disabled_katas:
  - ZC1005
  - ZC1011
```
