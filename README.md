# ZShellCheck

`zshellcheck` is a static analysis tool (linter) specifically designed for Zsh scripts. Unlike `shellcheck`, which focuses on POSIX sh/bash compatibility, `zshellcheck` understands Zsh syntax, best practices, and common pitfalls.

It parses Zsh scripts into an Abstract Syntax Tree (AST) and runs a series of checks ("Katas") to identify issues.

## Features

*   **Zsh-Specific Parsing:** Handles Zsh constructs like `[[ ... ]]`, `(( ... ))`, arrays, associative arrays, and modifiers.
*   **Extensible Katas:** Rules are implemented as independent "Katas" that can be easily added or disabled.
*   **Configurable:** Disable specific checks via `.zshellcheckrc` configuration file.
*   **Integration Ready:** Designed to work with `pre-commit` and CI pipelines.

## Installation

To install `zshellcheck`, ensure you have Go (version 1.18 or higher) installed, then run:

```bash
go install github.com/afadesigns/zshellcheck/cmd/zshellcheck@latest
```

This will install the `zshellcheck` executable into your `$GOPATH/bin` directory. Make sure `$GOPATH/bin` is in your system's `PATH`.

## Usage

To run `zshellcheck` on a file:

```bash
zshellcheck myscript.zsh
```

To run on multiple files or a directory recursively:

```bash
zshellcheck ./path/to/my/scripts
```

## Configuration

`zshellcheck` can be configured using a `.zshellcheckrc` file in YAML format. This file allows you to enable or disable specific Katas (checks).

Example `.zshellcheckrc`:

```yaml
disabled-katas:
  - ZC1001 # Disable the 'Prefer local over typeset' check
```

## Contributing Katas

The core of `zshellcheck`'s linting logic resides in its "Katas." A Kata is an independent Go file in the `pkg/katas/` directory that defines a specific check for a Zsh anti-pattern or style issue.

To contribute a new Kata:

1.  Create a new file `pkg/katas/zcXXXX.go` (where `XXXX` is a unique number)
2.  Implement the `Kata` interface.
3.  Register your Kata in `pkg/katas/katas.go`.

For more detailed information on contributing, please see `CONTRIBUTING.md`.

