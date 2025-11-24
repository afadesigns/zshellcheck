# ZShellCheck

```
 mmmmmm  mmmm  #             ""#    ""#      mmm  #                    #
     #" #"   " # mm    mmm     #      #    m"   " # mm    mmm    mmm   #   m
   m#   "#mmm  #"  #  #"  #    #      #    #      #"  #  #"  #  #"  "  # m"
  m"        "# #   #  #""""    #      #    #      #   #  #""""  #      #"#
 ##mmmm "mmm#" #   #  "#mm"    "mm    "mm   "mmm" #   #  "#mm"  "#mm"  #  "m
```

ZShellCheck is the definitive static analysis and comprehensive development suite for the entire Zsh ecosystem, meticulously engineered as the full Zsh equivalent of ShellCheck for Bash, uniquely offering intelligent automatic fixes, advanced formatting capabilities, and deep code analysis to deliver unparalleled quality, performance, and reliability for Zsh scripts, functions, and configurations.

## Inspiration

ZShellCheck draws significant inspiration from the esteemed `ShellCheck` project, a powerful static analysis tool for `sh`/`bash` scripts. While `ZShellCheck` is an independent development with a native focus on Zsh's unique syntax and semantics, `ShellCheck`'s commitment to improving shell script quality served as a guiding principle in our mission to provide an equally robust and tailored solution for the Zsh community.

## Features

*   **Zsh-Native Parsing:** Full understanding and handling of Zsh's unique constructs, including `[[ ... ]]`, `(( ... ))`, advanced arrays, associative arrays, and parameter expansion modifiers, applicable across scripts, functions, and configuration files.
*   **Extensible Katas:** A modular system where rules are implemented as independent "Katas," allowing for easy expansion, customization, and precise control over checks.
*   **Highly Configurable:** Tailor ZShellCheck's behavior to your project's needs by enabling or disabling specific checks via a flexible `.zshellcheckrc` configuration file.
*   **Seamless Integration:** Designed for effortless integration into modern development workflows, supporting `pre-commit` hooks and continuous integration (CI) pipelines to enforce quality at every stage.
*   **Automated Fixes (Planned):** Future versions will include capabilities for automatically resolving common issues, streamlining code maintenance.
*   **Code Formatting (Planned):** Upcoming features will provide built-in Zsh script and configuration formatting to ensure consistent style across your codebase.

## Installation & Usage

For detailed instructions on how to install and use `zshellcheck`, please refer to the [[Getting Started]] page.

## Implemented Checks (Katas)

For a complete and navigable list of all implemented checks, including detailed descriptions, bad/good examples, and configuration options, please refer to the [[Katas | All Katas]] page.

## Configuration

Tailor ZShellCheck to your project by creating a `.zshellcheckrc` file in your project root:

```yaml
disabled_katas:
  - ZC1005
  - ZC1011
```