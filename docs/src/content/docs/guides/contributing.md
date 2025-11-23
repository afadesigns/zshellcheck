---
title: Contributing
description: How to contribute to ZShellCheck.
---

## Pull Request Workflow

We follow a strict Pull Request (PR) workflow to ensure code quality and maintain a clear history.

1.  **Create a Branch**: Always start by creating a new branch for your changes. Use a descriptive name like `feat/new-kata`, `fix/parser-bug`, `docs/update-readme`, etc.
    ```bash
    git checkout -b feat/your-feature-name
    ```
2.  **Make Changes**: Implement your changes, ensuring they follow the project's coding style and conventions.
3.  **Test**: Run tests locally to ensure everything works as expected.
    ```bash
    go test ./...
    ./tests/integration_test.zsh
    ```
4.  **Commit**: Commit your changes with clear and concise messages. Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification.
    *   `feat: ...` for new features (ZC katas).
    *   `fix: ...` for bug fixes.
    *   `docs: ...` for documentation updates.
    *   `chore: ...` for maintenance tasks.
    *   `refactor: ...` for code restructuring.
    *   `test: ...` for test-related changes.
5.  **Push**: Push your branch to the repository.
    ```bash
    git push origin feat/your-feature-name
    ```
6.  **Create PR**: Create a Pull Request targeting the `main` branch.
    *   Provide a descriptive title and body explaining *why* the change is needed and *what* it does.
    *   Link any relevant issues (e.g., "Fixes #42").
    *   **Labels**: Apply appropriate labels to your PR (see below).
7.  **Review & Merge**: Wait for review. PRs require approval and passing CI checks. Only administrators can merge PRs into `main`.

## Labels

We use a specific set of labels to categorize issues and pull requests. Please use them appropriately:

| Label | Description |
| :--- | :--- |
| **`feat`** | New features or enhancements (e.g., adding a new Kata). |
| **`fix`** | Bug fixes. |
| **`docs`** | Documentation changes or improvements. |
| **`ci`** | Changes to CI/CD configuration or workflows. |
| **`deps`** | Dependency updates. |
| **`refactor`** | Code refactoring without changing behavior. |
| **`test`** | Adding or correcting tests. |
| **`chore`** | Routine maintenance tasks. |
| **`starter`** | Good tasks for newcomers to the project. |
| **`help`** | Extra attention or assistance is needed. |
| **`question`** | Further information is requested. |
| **`nofix`** | The issue or request will not be worked on. |
| **`duplicate`** | This issue or PR already exists. |
| **`invalid`** | The issue or PR is invalid or not applicable. |

## Coding Style

*   We use `gofmt` for Go code formatting.
*   We follow the standard Go coding conventions.
*   Please ensure that your code is well-documented and easy to understand.