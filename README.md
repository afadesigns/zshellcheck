```
 mmmmmm  mmmm  #             ""#    ""#      mmm  #                    #
     #" #"   " # mm    mmm     #      #    m"   " # mm    mmm    mmm   #   m
   m#   "#mmm  #"  #  #"  #    #      #    #      #"  #  #"  #  #"  "  # m"
  m"        "# #   #  #""""    #      #    #      #   #  #""""  #      #"#
 ##mmmm "mmm#" #   #  "#mm"    "mm    "mm   "mmm" #   #  "#mm"  "#mm"  #  "m
```

![CI](https://github.com/afadesigns/zshellcheck/actions/workflows/ci.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/afadesigns/zshellcheck)](https://goreportcard.com/report/github.com/afadesigns/zshellcheck)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/afadesigns/zshellcheck/badge)](https://securityscorecards.dev/viewer/?uri=github.com/afadesigns/zshellcheck)
[![codecov](https://codecov.io/gh/afadesigns/zshellcheck/graph/badge.svg?token=placeholder)](https://codecov.io/gh/afadesigns/zshellcheck)
[![SLSA](https://img.shields.io/badge/SLSA-Level%203-brightgreen)](https://slsa.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![Release](https://img.shields.io/github/v/release/afadesigns/zshellcheck)

**ZShellCheck** is the definitive static analysis and comprehensive development suite for the entire Zsh ecosystem, meticulously engineered as the full Zsh equivalent of ShellCheck for Bash. It offers intelligent automatic fixes (planned), advanced formatting capabilities, and deep code analysis to deliver unparalleled quality, performance, and reliability for Zsh scripts, functions, and configurations.

## Inspiration

ZShellCheck draws significant inspiration from the esteemed `ShellCheck` project, a powerful static analysis tool for `sh`/`bash` scripts. While `ZShellCheck` is an independent development with a native focus on Zsh's unique syntax and semantics, `ShellCheck`'s commitment to improving shell script quality served as a guiding principle in our mission to provide an equally robust and tailored solution for the Zsh community.

## Comparison

Why use ZShellCheck over ShellCheck? See our **[Detailed Comparison](docs/REFERENCE.md#comparison-vs-shellcheck)**.

**TL;DR**: Use **ShellCheck** for portable scripts (`sh`/`bash`). Use **ZShellCheck** for native **Zsh** scripts, plugins, and configuration.

## Table of Contents

- [Inspiration](#inspiration)
- [Comparison](docs/REFERENCE.md#comparison-vs-shellcheck)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](docs/USER_GUIDE.md#configuration)
- [Integrations](docs/USER_GUIDE.md#integrations)
- [Shell Completions](#shell-completions)
- [Architecture](docs/DEVELOPER.md#architecture-overview)
- [Troubleshooting](docs/USER_GUIDE.md#troubleshooting)
- [Developer Guide](docs/DEVELOPER.md)
- [Documentation](#documentation)
- [Changelog](#changelog)
- [Contributing](#contributing)
- [Governance](docs/REFERENCE.md#governance)
- [License](#license)

## Features

*   **Zsh-Native Parsing:** Full understanding and handling of Zsh's unique constructs, including `[[ ... ]]`, `(( ... ))`, advanced arrays, associative arrays, and parameter expansion modifiers, applicable across scripts, functions, and configuration files.
*   **Extensible Katas:** A modular system where rules are implemented as independent "Katas," allowing for easy expansion, customization, and precise control over checks.
*   **Highly Configurable:** Tailor ZShellCheck's behavior to your project's needs by enabling or disabling specific checks via a flexible `.zshellcheckrc` configuration file.
*   **Seamless Integration:** Designed for effortless integration into modern development workflows, supporting `pre-commit` hooks and continuous integration (CI) pipelines to enforce quality at every stage.

## Installation

ZShellCheck is written in Go and can be easily installed from source if you have a Go development environment configured.

### Helper Script

For a quick installation, you can use the provided helper script:

```bash
./install.sh
```

### From Go Modules

To install `zshellcheck`, ensure you have Go (version 1.25 or higher) installed, then run:

```bash
go install github.com/afadesigns/zshellcheck/cmd/zshellcheck@latest
```

This will install the `zshellcheck` executable into your `$GOPATH/bin` directory. Make sure `$GOPATH/bin` is in your system's `PATH`.

### Building from Source

For developers who want to build `zshellcheck` from its source code:

1.  Clone the repository:
    ```bash
    git clone https://github.com/afadesigns/zshellcheck.git
    cd zshellcheck
    ```
2.  Build the executable:
    ```bash
    go build -o zshellcheck cmd/zshellcheck/main.go
    ```
    This will create an executable named `zshellcheck` in your current directory.

## Usage

After installation, you can run ZShellCheck against your Zsh code files from your terminal. You can specify one or more files, or a directory to check recursively.

### Analyzing Files

```bash
# Analyze a single file
zshellcheck my_script.zsh

# Analyze multiple files
zshellcheck script1.zsh another_script.zsh

# Analyze a directory recursively
zshellcheck ./path/to/my/scripts
```

### Try it out

We provide a demo file with intentional violations so you can see ZShellCheck in action immediately.

```bash
zshellcheck examples/demo.zsh
```

ZShellCheck will output any identified violations directly to your terminal.

### Output Formats

You can control the output format using the `-format` flag:

*   **Text (default)**: Human-readable output.
    ```bash
    zshellcheck -format text my_script.zsh
    ```

*   **JSON**: Machine-readable JSON output, useful for integration with other tools or CI systems.
    ```bash
    zshellcheck -format json my_script.zsh
    ```

### Pre-commit Hook (Recommended)

To integrate ZShellCheck seamlessly into your development workflow and ensure code quality before commits, you can use it as a `pre-commit` hook.

1.  **Install `pre-commit`**:
    ```bash
    pip install pre-commit
    # Or brew install pre-commit on macOS
    ```

2.  **Configure `.pre-commit-config.yaml`**: Add the following configuration to a file named `.pre-commit-config.yaml` in the root of your Zsh project:

    ```yaml
    # .pre-commit-config.yaml
    -   repo: https://github.com/afadesigns/zshellcheck
        rev: v0.1.1 # Check releases for the latest version
        hooks:
        -   id: zshellcheck
    ```

3.  **Install the Hook**:
    ```bash
    pre-commit install
    ```

## Configuration

Tailor ZShellCheck to your project by creating a `.zshellcheckrc` file. For detailed instructions, see the **[Configuration Guide](docs/USER_GUIDE.md#configuration)**.

**Example `.zshellcheckrc`**:

```yaml
disabled_katas:
  - ZC1005 # Example: Disable "Use whence instead of which"
  - ZC1011 # Example: Disable "Use git porcelain commands instead of plumbing commands"
```

## Integrations

Want to use ZShellCheck in VS Code, Vim, or Neovim? Check out our **[Integrations Guide](docs/USER_GUIDE.md#integrations)**.

## Shell Completions

ZShellCheck provides completion scripts for Zsh and Bash.

### Zsh
Copy `completions/zsh/_zshellcheck` to a directory in your `$fpath` (e.g., `~/.zfunc`).
```zsh
fpath+=~/.zfunc
autoload -Uz compinit && compinit
```

### Bash
Source the completion script in your `.bashrc`.
```bash
source /path/to/zshellcheck/completions/bash/zshellcheck-completion.bash
```

## Architecture

Curious about how ZShellCheck works under the hood? Check out our [Architecture Guide](docs/DEVELOPER.md#architecture-overview) to learn about the Lexer, Parser, AST, and Kata Registry.

## Troubleshooting

Encountering issues? Check our **[Troubleshooting Guide](docs/USER_GUIDE.md#troubleshooting)** for solutions to common problems like "command not found" or parser errors.

## Developer Guide

Want to contribute code? Read our [Developer Guide](docs/DEVELOPER.md) and [AST Reference](docs/DEVELOPER.md#ast-reference) to get started with building, testing, and understanding the codebase.

## Documentation

For a comprehensive list of all implemented Katas (checks), including detailed descriptions, **good/bad code examples**, and configuration options, please refer to:

👉 **[KATAS.md](KATAS.md)**

Unsure about a term? Check the **[Glossary](docs/REFERENCE.md#glossary)**.

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for a history of changes and releases.

## Support

Need help? Have a question? Check out our [Support Guide](docs/USER_GUIDE.md#support).

## Contributing

We welcome contributions! Whether it's adding new Katas, improving the parser, or fixing bugs, your help is appreciated. For detailed instructions, please see [CONTRIBUTING.md](CONTRIBUTING.md).

See our [Governance Model](docs/REFERENCE.md#governance) for information on how this project is managed.

## License

Distributed under the MIT License. See `LICENSE` for more information.
