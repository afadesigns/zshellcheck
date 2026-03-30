# ZShellCheck Roadmap

ZShellCheck is an evolving static analysis tool for Zsh. Our mission is to provide the most comprehensive, fast, and reliable tooling for the Zsh ecosystem.

## 🚀 Milestones

### Version 0.0.x - The Foundation (Complete)
- [x] Establish core architecture (Lexer, Parser, AST, Walker).
- [x] Implement basic linting framework.
- [x] Set up CI/CD pipeline (GitHub Actions).
- [x] Create initial set of Katas (ZC1000+).
- [x] Comprehensive Documentation Suite.
- [x] Robust Installation Script (install.sh) with binary fallback.

### Version 0.1.x - Growing the Rule Set (In Progress)
- [x] Implement severity levels (error, warning, info, style) for all Katas.
- [x] Add `--severity` flag for filtering violations.
- [x] Add `--no-color` flag and verbose output mode.
- [x] Add SARIF output format for GitHub Security integration.
- [x] Add source code context to violation output.
- [x] Enhanced installer with auto-sudo, CI detection, man pages, and completions.
- [x] Reach 166 Katas (ZC1001 - ZC1169).
- [ ] Reach 250 Katas covering core Zsh idioms and common anti-patterns.
- [ ] Reach 500 Katas as the halfway milestone to 1.0.

### Version 1.0.0 - The 1000 Kata Milestone
- [ ] **Goal:** Implement 1000 Katas (ZC1000 - ZC2000) covering:
    - Syntax errors
    - Portability issues
    - Performance bottlenecks
    - Security vulnerabilities
    - Best practices
- [ ] **Language Server Protocol (LSP)**: Build an official LSP implementation to support VS Code, Neovim, and other editors natively with inline diagnostics and "Quick Fix" actions.
- [ ] **Auto-Fixer**: Implement `zshellcheck --fix` to automatically apply corrections for common violations (e.g., changing `[ ]` to `[[ ]]`).
- [ ] **Plugin System**: Allow users to write their own custom checks in Lua or Wasm.

## Long-Term Vision
- **Type Checking**: Experimental static type inference for Zsh scripts.
- **Formatter**: A strictly opinionated formatter (like `gofmt` or `prettier`) for Zsh.

## Progress Tracking

**Current Progress:** Version 0.1.66 (166/1000 Katas -- 16.6%).

```
[=================>                                                                  ] 166/1000
```

For the list of currently implemented Katas, please refer to [KATAS.md](KATAS.md).