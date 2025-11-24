# ZShellCheck vs ShellCheck

ZShellCheck is often compared to [ShellCheck](https://github.com/koalaman/shellcheck), the industry-standard linter for shell scripts. This document explains the key differences and when to use which.

## At a Glance

| Feature | ShellCheck | ZShellCheck |
| :--- | :--- | :--- |
| **Primary Focus** | POSIX `sh` and `bash` | **`zsh`** |
| **Zsh Support** | Limited / "Bash-mode" | **Native / Full** |
| **Language** | Haskell | **Go** |
| **Syntax Parsing** | Bash/Sh compatible | Zsh specific (`[[`, `((`, modifiers) |
| **Philosophy** | Compatibility & Safety | Zsh Idioms & Power |

## Detailed Comparison

### 1. Zsh Syntax Support

**ShellCheck** treats scripts as essentially Bash or Sh. It often flags valid Zsh syntax as errors because it doesn't understand them.

*   **Associative Arrays**:
    *   *Zsh:* `typeset -A my_map; my_map[key]=val`
    *   *ShellCheck:* May flag syntax errors or warn about bash compatibility.
    *   *ZShellCheck:* Parses and checks correctly.
*   **Modifiers**:
    *   *Zsh:* `${file:h}`, `${path:a}`
    *   *ShellCheck:* Often flags as "Bad substitution".
    *   *ZShellCheck:* Understands these are valid Zsh modifiers.
*   **Glob Qualifiers**:
    *   *Zsh:* `ls *(.)` (list only files)
    *   *ShellCheck:* Syntax error.
    *   *ZShellCheck:* Native support.

### 2. Focus and Philosophy

**ShellCheck** is excellent for ensuring your scripts are portable (run on dash, bash, sh) and safe from common pitfalls like quoting issues.

**ZShellCheck** acknowledges that you have chosen Zsh for its power. It encourages **idiomatic Zsh code**:
*   Promotes `[[ ... ]]` over `[ ... ]`.
*   Promotes `(( ... ))` for arithmetic.
*   Promotes usage of Zsh builtins (`whence`) over external commands (`which`).

### 3. When to use which?

| Use Case | Recommendation |
| :--- | :--- |
| **Writing portable scripts** (must run on Debian/Ubuntu/CentOS default shells) | Use **ShellCheck**. |
| **Writing `.zshrc` or `.zprofile`** | Use **ZShellCheck**. |
| **Writing Zsh plugins** (Oh-My-Zsh, antigen, etc.) | Use **ZShellCheck**. |
| **Writing complex automation explicitly for Zsh** | Use **ZShellCheck**. |

## Can I use both?

**Yes!** However, you will likely need to disable many checks in ShellCheck if you use Zsh features heavily, or it will be very noisy. ZShellCheck is designed to fill the gap where ShellCheck falls short for Zsh users.
