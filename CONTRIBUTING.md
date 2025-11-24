# Contributing to zshellcheck

We welcome contributions! Whether it's adding new Katas, improving the parser, or fixing bugs, your help is appreciated.

## Pull Request Workflow

We follow a strict Pull Request (PR) workflow to ensure code quality and maintain a clear history. This workflow is designed to facilitate smooth collaboration and maintain an organized project.

1.  **Sync `main`**: Before starting new work, ensure your local `main` branch is up-to-date with the remote `main`.
    ```bash
    git checkout main
    git pull origin main
    ```
2.  **Create a Branch**: Always create a new, descriptive branch for your changes. Use a prefix that indicates the type of change (e.g., `feat/`, `fix/`, `docs/`, `chore/`).
    ```bash
    git checkout -b feat/your-feature-name
    ```
3.  **Implement & Test**: Make your changes, adhering to coding style and conventions. Run local tests to verify functionality.
    ```bash
    go test ./...
    ./tests/integration_test.zsh
    ```
4.  **Commit**: Commit your changes using [Conventional Commits](https://www.conventionalcommits.org/) for clear history. Examples:
    *   `feat: Implement new Kata ZCXXXX (Short description)`
    *   `fix: Resolve parser bug in arithmetic expressions`
    *   `docs: Update wiki links`
    *   `chore: Upgrade npm dependencies`
5.  **Push**: Push your local branch to the remote repository.
    ```bash
    git push origin your-branch-name
    ```
6.  **Create Pull Request**: Use the GitHub CLI to create a Pull Request from your branch to `main`.
    ```bash
    gh pr create --title "feat: Your feature title" --body "A detailed description of your changes." --base main
    ```
    *   Provide a clear title and body explaining the *why* and *what* of your changes.
    *   Link any relevant issues (e.g., `Closes #123`, `Fixes #45`).
        *   **Labels**: Apply [appropriate labels](#project-labels) to your PR.
    7.  **Review & Merge**: Address any review comments. Once approved and all CI checks pass, an administrator will merge the PR. We use squash merges to maintain a clean Git history.
    
    ## Documentation
    
    For comprehensive documentation, including detailed usage, configuration, and a full list of implemented Katas, please refer to [KATAS.md](KATAS.md).
    
    For developers, please refer to:
    *   [Developer Guide](docs/DEVELOPER.md) - How to build, test, and debug.
    *   [AST Reference](docs/DEVELOPER.md#ast-reference) - Detailed documentation of the Abstract Syntax Tree nodes.
    *   [Architecture](docs/DEVELOPER.md#architecture-overview) - High-level overview of the system.
    
    ## Coding Style        *   We use `gofmt` for Go code formatting.
    *   We follow the standard Go coding conventions.
    *   Please ensure that your code is well-documented and easy to understand.
    
    ### Running Linters and Formatters
    
    Before submitting a Pull Request, please ensure your code passes all linting and formatting checks:
    
    ```bash
    go fmt ./...       # Format Go code
    go vet ./...       # Run Go vet (static analysis)
    golangci-lint run  # Run golangci-lint (if installed)
    ```
    
    ## Adding a New Kata
    
    Katas are the core rules of `zshellcheck`. To add one:
    
    1.  **Define the Kata:** Create a new file `pkg/katas/zcXXXX.go`.
    2.  **Register:** In the `init()` function, register the Kata with the `RegisterKata` function, specifying the AST node type it targets.
    3.  **Implement Logic:** Write the check function that inspects the node and returns a list of `Violation`s.
    4.  **Add Tests:** Create `pkg/katas/katatests/zcXXXX_test.go` with test cases covering valid and invalid Zsh code.
    
    ### Example Kata
    
    ```go
    func init() {
        RegisterKata(ast.SimpleCommandNode, Kata{
            ID: "ZC1099",
            Title: "Avoid foo command",
            Description: "The foo command is deprecated.",
            Check: checkZC1099,
        })
    }
    ```
    
    ## Project Labels
    
    We use a specific set of labels to categorize issues and pull requests, helping us organize and prioritize work effectively. Please use them appropriately.
    
    | Label | Description |
    | :--- | :--- |
    | **`feat`** | New features or significant enhancements. |
    | **`fix`** | Bug fixes. |
    | **`docs`** | Documentation changes or improvements. |
    | **`ci`** | Updates to CI/CD configurations or workflows. |
    | **`deps`** | Dependency updates. |
    | **`refactor`** | Code restructuring without behavior changes. |
    | **`test`** | Additions or corrections to tests. |
    | **`chore`** | Routine maintenance tasks (e.g., updating build scripts, `.gitignore`). |
    | **`starter`** | Good entry-level tasks for new contributors. |
    | **`help`** | Requires extra attention or assistance. |
    | **`question`** | Seeking further information or clarification. |
    | **`nofix`** | The issue or request will not be addressed. |
    | **`duplicate`** | This issue or PR is a duplicate. |
    | **`invalid`** | The issue or PR is invalid or not applicable. |
    