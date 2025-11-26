- **ZShellCheck Versioning:** Versioning follows `Major.Minor.Patch` based on the total number of implemented Katas.
  - **Major:** Represents each full thousand of Katas. (e.g., 1000 Katas = Major `1`, 2500 Katas = Major `2`).
  - **Minor:** Represents each full hundred of Katas *after* the Major component is accounted for. (e.g., for 120 Katas, after 0 Major, there's 1 hundred, so Minor `1`).
  - **Patch:** Represents the remaining Katas (tens and units) *after* Major and Minor components are accounted for. (e.g., for 120 Katas, after 0 Major and 1 Minor, there are 20 remaining, so Patch `20`).
  - Example 1: 120 Katas = `0.1.20` (0 thousands, 1 hundred, 20 remaining)
  - Example 2: 52 Katas = `0.0.52` (0 thousands, 0 hundreds, 52 remaining)
  - Example 3: 1005 Katas = `1.0.5` (1 thousand, 0 hundreds, 5 remaining)

  - **Release Strategy:** Draft releases can be created at any point. However, the **publication of Version `1.0.0` is a significant milestone, to be released to the GitHub Marketplace when 1000 Katas have been implemented (corresponding to reaching Kata ZC2000, assuming Kata numbering started at ZC1000).**

- **ZShellCheck Workflow:** Follow these steps for development:
  1. **Create Branch:** Create a new, descriptive branch. The prefix of the branch name should reflect the task type using one of the [Project Labels](#project-labels-and-descriptions).
     - **Command:** `git checkout -b <type>/<description-of-task>` (e.g., `git checkout -b feat/implement-zc1001` or `git checkout -b fix/resolve-parser-bug`).
  2. **Implement & Test:** Implement the Kata (or other changes) and ensure it's fully tested. For new Katas, this includes adding integration tests (`tests/integration_test.zsh`).
  3. **Commit:** Commit changes using [Conventional Commits](https://www.conventionalcommits.org/). The commit type should align with the branch prefix and task (e.g., `feat: Implement ZCXXXX to improve foo`, `fix: Resolve issue with bar parsing`).
  4. **Push Branch:** Push your local branch to the remote repository.
  5. **Create Pull Request:** Create a Pull Request (PR) against the `main` branch. Ensure the PR title and body clearly describe the changes.
     - **Labels:** Apply the most appropriate [Project Labels](#project-labels-and-descriptions) to the PR.
  6. **CI Verification:** Before merging, ensure all Continuous Integration (CI) checks pass. This includes:
     - Security policy
     - Security advisories
     - Private vulnerability reporting
     - Dependabot alerts
     - Code scanning alerts
       - **Details:** The CodeQL scan is configured via `.github/workflows/ci.yml:security`. It uses `CodeQL (2.23.5)` with extensions `codeql/go-queries (1.4.8)`, `codeql/go-all (5.0.1)`, and `codeql/threat-models (1.0.34)`.
     - Secret scanning alerts
     - **CI Setup Details:**
       - **Setup type:** Actions workflow
       - **Workflow path:** `.github/workflows/ci.yml`
       - **First scan:** 11 hours ago
       - **Last scan:** 9 hours ago
       - **Scan events:** Push to `main`, Pull request to `main`
     - **Action:** Apply a `sleep` command (e.g., `sleep 60s`) after pushing changes and before attempting to merge/tag, to allow GitHub Actions and security scans to complete and report their status.
  7. **Merge & Tag (Admin Bypass):** Once approved and all CI checks pass, an administrator will merge the PR (squash merge is preferred to maintain a clean Git history).
     - **Action:** For *each new implemented Kata*, after merging, **Gemini must update the version tag** according to the [Versioning](#zshellcheck-versioning) rules.
  8. **Delete Branch:** Delete the local and remote feature branch after a successful merge.

- **Project Labels and Descriptions:**
  - `feat`: New features or significant enhancements.
  - `fix`: Bug fixes.
  - `docs`: Documentation changes or improvements.
  - `ci`: Updates to CI/CD configurations or workflows.
  - `deps`: Dependency updates.
  - `refactor`: Code restructuring without behavior changes.
  - `test`: Additions or corrections to tests.
  - `chore`: Routine maintenance tasks (e.g., updating build scripts, `.gitignore`).
  - `starter`: Good entry-level tasks for new contributors.
  - `help`: Requires extra attention or assistance.
  - `question`: Seeking further information or clarification.
  - `nofix`: The issue or request will not be addressed.
  - `duplicate`: This issue or PR is a duplicate.
  - `invalid`: The issue or PR is invalid or not applicable.

- **Contributing to zshellcheck:**
  We welcome contributions! Whether it's adding new Katas, improving the parser, or fixing bugs, your help is appreciated.

  **Pull Request Workflow:**
  We follow a strict Pull Request (PR) workflow to ensure code quality and maintain a clear history. This workflow is designed to facilitate smooth collaboration and maintain an organized project.

  **Sync main:** Before starting new work, ensure your local main branch is up-to-date with the remote main.
  ```bash
  git checkout main
  git pull origin main
  ```
  **Create a Branch:** Always create a new, descriptive branch for your changes. Use a prefix that indicates the type of change (e.g., `feat/`, `fix/`, `docs/`, `chore/`).
  ```bash
  git checkout -b feat/your-feature-name
  ```
  **Implement & Test:** Make your changes, adhering to coding style and conventions. Run local tests to verify functionality.
  ```bash
  go test ./...
  ./tests/integration_test.zsh
  ```
  **Commit:** Commit your changes using [Conventional Commits](https://www.conventionalcommits.org/) for clear history. Examples:
  - `feat: Implement new Kata ZCXXXX (Short description)`
  - `fix: Resolve parser bug in arithmetic expressions`
  - `docs: Update wiki links`
  - `chore: Upgrade npm dependencies`
  **Push:** Push your local branch to the remote repository.
  ```bash
  git push origin your-branch-name
  ```
  **Create Pull Request:** Use the GitHub CLI to create a Pull Request from your branch to main.
  ```bash
  gh pr create --title "feat: Your feature title" --body "A detailed description of your changes." --base main
  ```
  - Provide a clear title and body explaining the why and what of your changes.
  - Link any relevant issues (e.g., `Closes #123`, `Fixes #45`).
  - **Labels:** Apply appropriate labels to your PR.
  **Review & Merge:** Address any review comments. Once approved and all CI checks pass, an administrator will merge the PR. We use squash merges to maintain a clean Git history.
  **Documentation:**
  For comprehensive documentation, including detailed usage, configuration, and a full list of implemented Katas, please refer to the [[Katas/ZC1000-ZC1099/Index | Foundational Zsh Checks]] page.

  **Coding Style:**
  - We use `gofmt` for Go code formatting.
  - We follow the standard Go coding conventions.
  - Please ensure that your code is well-documented and easy to understand.
  **Running Linters and Formatters:**
  Before submitting a Pull Request, please ensure your code passes all linting and formatting checks:

  ```bash
  go fmt ./...       # Format Go code
  go vet ./...       # Run Go vet (static analysis)
  golangci-lint run  # Run golangci-lint (if installed)
  ```
  **Adding a New Kata:**
  Katas are the core rules of `zshellcheck`. To add one:

  1.  **Define the Kata:** Create a new file `pkg/katas/zcXXXX.go`.
  2.  **Register:** In the `init()` function, register the Kata with the `RegisterKata` function, specifying the AST node type it targets.
  3.  **Implement Logic:** Write the check function that inspects the node and returns a list of `Violation`s.
  4.  **Add Tests:** Create `pkg/katas/katatests/zcXXXX_test.go` with test cases covering valid and invalid Zsh code.
  **Example Kata:**
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

- **Minimizing Friction (Agent Notes):**
  - This section is for notes and configurations aimed at minimizing friction during development and interaction within the `zshellcheck` project. It can be extended collaboratively as new insights or preferences emerge.
