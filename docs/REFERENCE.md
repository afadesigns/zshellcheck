# Reference

This document contains the glossary, governance model, and comparison with other tools.

## Table of Contents

- [Governance](#governance)
- [Comparison vs ShellCheck](#comparison-vs-shellcheck)
- [Glossary](#glossary)

---

## Governance

ZShellCheck is led by **@afadesigns** (Andreas Fahl), the lead maintainer and sole current committer. Everything below describes how community contribution is expected to work as the project grows.

### Roles

- **Lead Maintainer** — final say on roadmap, release cadence, and architectural direction. Currently @afadesigns.
- **Maintainers** — review and merge PRs, triage issues. No additional maintainers at time of writing; CODEOWNERS routes all reviews to @afadesigns.
- **Contributors** — anyone who opens an issue or PR.

### Decision Making

Non-trivial direction is discussed on GitHub issues or PRs. Consensus preferred; the lead maintainer has the casting vote in deadlocks. Breaking changes and major features require a design-discussion issue before implementation.

### Review Requirements (enforced)

- Every PR to `main` needs one approving review from the CODEOWNERS set.
- Commits must be GPG-signed.
- Required status checks (`test`, `security`, `sbom`) must pass.
- No force-push, no branch deletion, no unsigned merges.

---

## Comparison vs ShellCheck

| Feature | ShellCheck | ZShellCheck |
| :--- | :--- | :--- |
| **Focus** | `sh`/`bash` (POSIX) | **`zsh`** (Native) |
| **Language** | Haskell | Go |
| **Philosophy** | Portability | Zsh Power |
| **Checks** | ~500 | 1000 |
| **Output** | Text, JSON, GCC, TTY | Text, JSON, SARIF |
| **Severity** | error, warning, info, style | error, warning, info, style |
| **Auto-fix** | Partial | Planned |

**Use ZShellCheck for:** `.zshrc`, Zsh plugins, and scripts using specific Zsh features (`[[`, modifiers).

---

## Glossary

-   **Kata**: A specific check/rule (e.g., `ZC1001`).
-   **AST**: Abstract Syntax Tree.
-   **Lexer**: Tokenizer.
-   **Walker**: AST traverser.
-   **Registry**: Central store of Katas.
-   **Severity**: The impact level of a Kata violation (`error`, `warning`, `info`, `style`).
-   **SARIF**: Static Analysis Results Interchange Format -- used for GitHub Security integration.
