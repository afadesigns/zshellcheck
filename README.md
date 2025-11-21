<div align="center">
  <!-- ASCII Art Placeholder -->
  <pre>
    ____  __          __               __   ______          __
   / __ \/ /_  ____  / /_  ____ ______/ /__/ ____/___  ____/ /__
  / /_/ / __ \/ __ \/ __ \/ __ `/ ___/ //_/ /   / __ \/ __  / _ \
 / ____/ / / / /_/ / / / / /_/ / /__/ ,< / /___/ /_/ / /_/ /  __/
/_/   /_/ /_/\____/_/ /_/\__,_/\___/_/|_|\____/\____/\__,_/\___/
  </pre>
  <!-- ASCII Art Placeholder -->

  <h1>ZShellCheck</h1>
  <p><strong>Your wise sensei for clean, fast, and safe Zsh.</strong></p>

  <p>
    <a href="https://github.com/afadesigns/zshellcheck/actions/workflows/test.yml"><img src="https://github.com/afadesigns/zshellcheck/actions/workflows/test.yml/badge.svg" alt="CI Status"></a>
    <a href="https://golang.org/"><img src="https://img.shields.io/badge/Go-1.18+-blue.svg" alt="Go Version"></a>
    <a href="https://github.com/afadesigns/zshellcheck/blob/main/LICENSE"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License"></a>
    <a href="https://github.com/afadesigns/zshellcheck/stargazers"><img src="https://img.shields.io/github/stars/afadesigns/zshellcheck.svg?style=social&label=Star" alt="GitHub Stars"></a>
  </p>
</div>

## Why ZShellCheck?

Zsh is a powerful, feature-rich shell, but its complexity can lead to subtle bugs, performance issues, and hard-to-maintain scripts. While `bash` has `ShellCheck`, Zsh has lacked a dedicated, modern static analysis tool that understands its unique syntax and semantics. ZShellCheck fills this gap, providing a friendly "sensei" to guide you toward writing better Zsh code.

## Show, Don't Tell

*(Animated GIF of ZShellCheck in action will be here)*

## Gallery of Wisdom (Features)

Here are a few examples of how ZShellCheck helps you improve your Zsh code.

### ❌ Before (The "Confused Grasshopper")
```zsh
# Example 1: Incorrect array access
my_array=(one two three)
echo $my_array[1] # This doesn't do what you think!

# Example 2: Using legacy command substitution
files=`ls *.zsh`

# Example 3: Unsafe use of '['
if [ $ZSH_VERSION > 5.0 ]; then
  echo "Modern Zsh"
fi
```

### ✅ After (The "Enlightened Master")
```zsh
# Example 1: Correct array access
my_array=(one two three)
echo $my_array[1] # ZShellCheck: [ZC1001] Use ${my_array[1]} for array element access.
echo ${my_array[1]}

# Example 2: Using modern command substitution
files=$(ls *.zsh) # ZShellCheck: [ZC1002] Use $(...) instead of backticks for command substitution.

# Example 3: Safe and powerful '[['
if [[ $ZSH_VERSION > 5.0 ]]; then # ZShellCheck: [ZC1003] Prefer [[...]] over [...] for modern Zsh tests.
  echo "Modern Zsh"
fi
```

## Getting Started

### Installation

#### Homebrew (macOS)
```sh
brew install afadesigns/zshellcheck/zshellcheck
```

#### Go
```sh
go install github.com/afadesigns/zshellcheck/cmd/zshellcheck@latest
```

#### From Binary
Download the latest release from the [Releases](https://github.com/afadesigns/zshellcheck/releases) page.

### Usage
```sh
zshellcheck your_script.zsh
```

## The Dojo (Usage)
```sh
zshellcheck your_script.zsh
```

## Learn the Katas (Rules)

ZShellCheck uses a system of "Katas" (rules) to teach you better Zsh. Each Kata has a unique ID and a detailed explanation.

➡️ [**See the full list of Katas in our Wiki**](https://github.com/afadesigns/zshellcheck/wiki)

## Development

Interested in contributing to ZShellCheck? Here's how to get started:

1.  **Clone the repository:**
    ```sh
    git clone https://github.com/afadesigns/zshellcheck.git
    ```
2.  **Install dependencies:**
    ```sh
    go mod tidy
    ```
3.  **Run the tests:**
    ```sh
    go test ./...
    ```
4.  **Build the binary:**
    ```sh
    go build ./cmd/zshellcheck
    ```

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for more details.

## License

ZShellCheck is licensed under the [MIT License](LICENSE).

## Code of Conduct

Please note that this project is released with a [Contributor Code of Conduct](CODE_OF_CONDUCT.md). By participating in this project you agree to abide by its terms.
