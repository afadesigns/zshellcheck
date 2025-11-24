# Reference

This document contains the glossary, governance model, and comparison with other tools.

## Table of Contents

- [Governance](#governance)
- [Comparison vs ShellCheck](#comparison-vs-shellcheck)
- [Glossary](#glossary)

---

## Governance

### Overview
ZShellCheck is a community-maintained project led by the founder (**@afadesigns**).

### Roles
-   **Lead Maintainer**: Final say on roadmap.
-   **Maintainers**: Merge PRs, triage issues.
-   **Contributors**: Submit PRs and issues.

### Decision Making
Decisions are made by consensus. The Lead Maintainer has the casting vote in deadlocks.

---

## Comparison vs ShellCheck

| Feature | ShellCheck | ZShellCheck |
| :--- | :--- | :--- |
| **Focus** | `sh`/`bash` (POSIX) | **`zsh`** (Native) |
| **Language** | Haskell | Go |
| **Philosophy** | Portability | Zsh Power |

**Use ZShellCheck for:** `.zshrc`, Zsh plugins, and scripts using specific Zsh features (`[[`, modifiers).

---

## Glossary

-   **Kata**: A specific check/rule (e.g., `ZC1001`).
-   **AST**: Abstract Syntax Tree.
-   **Lexer**: Tokenizer.
-   **Walker**: AST traverser.
-   **Registry**: Central store of Katas.
