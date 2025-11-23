# ZShellCheck Roadmap

ZShellCheck is an evolving static analysis tool for Zsh. Our primary goal is to reach 1000 Katas (checks) to ensure comprehensive coverage of Zsh scripting pitfalls.

## 🚀 Milestones

### ✅ Version 0.0.x - The Foundation
- [x] Establish core architecture (Lexer, Parser, AST, Walker).
- [x] Implement basic linting framework.
- [x] Set up CI/CD pipeline (GitHub Actions).
- [x] Create initial set of ~70 Katas.

### 🚧 Version 0.1.0 - The First 100
- [ ] **Goal:** Implement 100 high-value Katas.
- [ ] Refine Parser to support complex Zsh constructs (arrays, arithmetic, expansions).
- [ ] Improve error reporting and output formatting.

### 🔮 Version 0.5.0 - Mid-Term Goals
- [ ] **Goal:** Reach 500 Katas.
- [ ] Deep analysis (control flow, variable scoping).
- [ ] Auto-fix capabilities for simple violations.
- [ ] Plugin system for custom checks.

### 🌟 Version 1.0.0 - The 1000 Kata Milestone
- [ ] **Goal:** Complete 1000 Katas covering:
    - Syntax errors
    - Portability issues
    - Performance bottlenecks
    - Security vulnerabilities
    - Best practices
- [ ] Stable API for integrations.
- [ ] Comprehensive documentation and wiki.

## 📈 Progress Tracking

**Current Progress:** 73 Katas implemented (Version 0.0.73).

We track progress by version number. Version `0.x.y` roughly corresponds to `x*100 + y` Katas (e.g., `0.0.73` = 73 Katas).

Check the [README](./README.md) for the list of currently implemented Katas.
