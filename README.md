# ZShellCheck

`zshellcheck` is a static analysis tool (linter) specifically designed for Zsh scripts. Unlike `shellcheck`, which focuses on POSIX sh/bash compatibility, `zshellcheck` understands Zsh syntax, best practices, and common pitfalls.

It parses Zsh scripts into an Abstract Syntax Tree (AST) and runs a series of checks ("Katas") to identify issues.

## Features

*   **Zsh-Specific Parsing:** Handles Zsh constructs like `[[ ... ]]`, `(( ... ))`, arrays, associative arrays, and modifiers.
*   **Extensible Katas:** Rules are implemented as independent "Katas" that can be easily added or disabled.
*   **Configurable:** Disable specific checks via `.zshellcheckrc` configuration file.
*   **Integration Ready:** Designed to work with `pre-commit` and CI pipelines.

## Documentation

For comprehensive documentation, including detailed usage, configuration, and a full list of implemented Katas, please visit the [GitHub Wiki](https://github.com/afadesigns/zshellcheck/wiki).
