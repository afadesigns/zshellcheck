# Getting Started with ZShellCheck

This guide will walk you through the essential steps to get ZShellCheck up and running, allowing you to start analyzing your Zsh code for quality and best practices.

## 1. Installation

ZShellCheck is written in Go and can be easily installed from source if you have a Go development environment configured.

### From Go Modules

To install `zshellcheck`, ensure you have Go (version 1.18 or higher) installed, then run:

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

### Pre-commit Hook (Recommended)

To integrate ZShellCheck seamlessly into your development workflow and ensure code quality before commits, you can use it as a `pre-commit` hook.

1.  **Install `pre-commit`**: If you haven't already, install the `pre-commit` framework:
    ```bash
    pip install pre-commit
    # Or brew install pre-commit on macOS
    ```

2.  **Configure `.pre-commit-config.yaml`**: Add the following configuration to a file named `.pre-commit-config.yaml` in the root of your Zsh project:

    ```yaml
    # .pre-commit-config.yaml
    -   repo: https://github.com/afadesigns/zshellcheck
        rev: v0.0.93 # Always use a specific release tag for stability, e.g., v0.0.93
        hooks:
        -   id: zshellcheck
    ```
    *   **`repo`**: Points to the ZShellCheck GitHub repository.
    *   **`rev`**: Specifies the exact [release tag](https://github.com/afadesigns/zshellcheck/releases) of ZShellCheck to use. **It is highly recommended to use a specific release tag (e.g., `v0.0.93`) for stability and reproducible builds.** Avoid `main` or `master` for production use.
    *   **`id`**: Refers to the `zshellcheck` hook defined within the repository.

3.  **Install the Hook**: From your project root, install the `pre-commit` hooks:
    ```bash
    pre-commit install
    ```
    Now, ZShellCheck will automatically run on your Zsh scripts every time you make a `git commit`.

## 2. Basic Usage

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

### Configuration (Disabling Checks)

If you need to disable specific Katas (checks) for your project, create a `.zshellcheckrc` file in your project's root directory. This file uses YAML syntax.

**Example `.zshellcheckrc`**:

```yaml
disabled_katas:
  - ZC1005 # Example: Disable "Use whence instead of which"
  - ZC1011 # Example: Disable "Use git porcelain commands instead of plumbing commands"
```

Any Katas listed under `disabled_katas` will be skipped by ZShellCheck when analyzing your project.

## What's Next?

*   Explore the [[Katas | All Katas]] page to understand the types of issues ZShellCheck identifies.
*   Learn how to [Contribute](/Contributing) by adding new Katas or improving the parser.
*   Dive deeper into the project's vision on the [Roadmap](/Roadmap).
