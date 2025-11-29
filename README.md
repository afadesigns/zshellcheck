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
[![codecov](https://codecov.io/gh/afadesigns/zshellcheck/graph/badge.svg)](https://codecov.io/gh/afadesigns/zshellcheck)
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

The easiest way to install ZShellCheck is via the automated installer script. It supports **Linux** and **macOS**.

### Automatic Install (Recommended)

This will install the binary, man pages, and shell completions. It detects if you have Go installed; if not, it downloads the latest pre-built binary.

```bash
# Clone the repository or download the script
./install.sh
```

**Features:**
*   **Binary Fallback:** No Go environment required. Downloads binaries automatically.
*   **Interactive:** GUIDes you through adding `zshellcheck` to your `PATH` and `fpath`.
*   **Automated:** Use `./install.sh -y` for non-interactive/CI environments.
*   **Version Control:** Install a specific version with `./install.sh -v v0.1.0`.
*   **Uninstall:** Remove cleanly with `./install.sh --uninstall`.

### From Go Modules

If you prefer standard Go tools:

```bash
go install github.com/afadesigns/zshellcheck/cmd/zshellcheck@latest
```

### Building from Source

For contributors:

1.  Clone the repository.
2.  Run `./install.sh` (it detects the source repo and builds locally).

## Usage

After installation, run ZShellCheck against your Zsh files:

```bash
zshellcheck my_script.zsh
```

### Output Formats

*   **Text (default)**: Human-readable with ANSI colors.
*   **JSON**: `zshellcheck -format json file.zsh`
*   **SARIF**: `zshellcheck -format sarif file.zsh` (Github Security integration)

### Pre-commit Hook

Add this to your `.pre-commit-config.yaml`:

```yaml
-   repo: https://github.com/afadesigns/zshellcheck
    rev: v0.1.1
    hooks:
    -   id: zshellcheck
```

## Configuration

Customize checks via `.zshellcheckrc`. See the [Configuration Guide](docs/USER_GUIDE.md#configuration).

## Integrations

See our [Integrations Guide](docs/USER_GUIDE.md#integrations) for VS Code, Vim, and Neovim setup.

## Shell Completions

The `./install.sh` script installs completions automatically for Zsh and Bash.

**Manual Setup (Zsh):**
If you installed manually, add the `completions/zsh` directory to your `$fpath`:
```zsh
fpath+=/path/to/zshellcheck/completions/zsh
autoload -Uz compinit && compinit
```

**Manual Setup (Bash):**
Source the script in your `.bashrc`:
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
