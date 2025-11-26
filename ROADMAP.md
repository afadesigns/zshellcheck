# ZShellCheck Roadmap

ZShellCheck is an evolving static analysis tool for Zsh. Our mission is to provide the most comprehensive, fast, and reliable tooling for the Zsh ecosystem.

## 🚀 Milestones

### ✅ Version 0.0.x - The Foundation
- [x] Establish core architecture (Lexer, Parser, AST, Walker).
- [x] Implement basic linting framework.
- [x] Set up CI/CD pipeline (GitHub Actions).
- [x] Create initial set of Katas (ZC1000+).
- [x] Comprehensive Documentation Suite.

### 🌟 Version 1.0.0 - The 1000 Kata Milestone
- [ ] **Goal:** Implement 1000 Katas (ZC1000 - ZC2000) covering:
    - Syntax errors
    - Portability issues
    - Performance bottlenecks
    - Security vulnerabilities
    - Best practices
- [ ] **Language Server Protocol (LSP)**: Build an official LSP implementation to support VS Code, Neovim, and other editors natively with inline diagnostics and "Quick Fix" actions.
- [ ] **Auto-Fixer**: Implement `zshellcheck --fix` to automatically apply corrections for common violations (e.g., changing `[ ]` to `[[ ]]`).
- [ ] **Plugin System**: Allow users to write their own custom checks in Lua or Wasm.

## 🔮 Long-Term Vision
- **Type Checking**: Experimental static type inference for Zsh scripts.
- **Formatter**: A strictly opinionated formatter (like `gofmt` or `prettier`) for Zsh.

## 📈 Progress Tracking

**Current Progress:** Version 0.0.92 (92/1000 Katas).

For the list of currently implemented Katas, please refer to [KATAS.md](KATAS.md).