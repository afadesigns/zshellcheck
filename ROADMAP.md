# ZShellCheck Roadmap

ZShellCheck is an evolving static analysis tool for Zsh.
Our mission is to provide the most comprehensive, fast, and reliable tooling for the Zsh ecosystem.

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
- [x] Reach 250 Katas covering core Zsh idioms and common anti-patterns.
- [x] Reach 500 Katas as the halfway milestone to 1.0.

### Version 1.0.0 - The 1000 Kata Milestone (Complete)
- [x] **Goal:** Implement 1000 Katas (ZC1001 - ZC2003) covering:
    - Syntax errors
    - Portability issues
    - Performance bottlenecks
    - Security vulnerabilities
    - Best practices

### Version 1.x - Beyond the Milestone
- [ ] **Language Server Protocol (LSP)**: Build an official LSP implementation to support VS Code, Neovim, and other editors natively with inline diagnostics and "Quick Fix" actions.
- [x] **Auto-Fixer core** (v1.0.14+): `-fix`, `-diff`, `-dry-run` flags applying deterministic per-kata rewrites.
  The set of fix-enabled katas grows with each release.
- [ ] **Auto-Fixer coverage**: 137 of 1000 katas (13.7%) ship a deterministic rewrite as of the latest tag.
  Expansion continues per release; the structural ceiling is the subset of detections that admit a context-free, idempotent, byte-exact rewrite — many advisory or context-dependent detections will remain detection-only.
- [ ] **Plugin System**: Allow users to write their own custom checks in Lua or Wasm.
- [ ] **Distribution channels** — broaden install paths beyond `./install.sh`, `go install`, and the signed Releases archive:
  - [ ] **Homebrew (homebrew-core)** — the project will be submitted once the [Homebrew notability threshold](https://docs.brew.sh/Acceptable-Formulae) is met (≥75 stars / ≥30 forks / ≥30 watchers for a third-party submission, or ≥225 stars / ≥90 forks / ≥90 watchers for a self-submission).
    The `homebrew-eligibility` workflow files a tracking issue automatically the week the threshold is crossed.
    A community member is welcome to file the new-formula PR on the project's behalf — third-party submissions face the lower bar.
  - [x] **Linux packages** — `.deb`, `.rpm`, `.apk` produced via goreleaser `nfpms:` and shipped on every Releases tag.
  - [ ] **AUR** — `zshellcheck-bin` published via goreleaser `aurs:` for Arch / Manjaro users.
  - [x] **PowerShell installer** — `install.ps1` mirroring `install.sh` for Windows users.
  - [x] **Docker image** — `FROM scratch` static binary published to `ghcr.io/afadesigns/zshellcheck` on tag, signed by cosign.

## Long-Term Vision
- **Type Checking**: Experimental static type inference for Zsh scripts.
- **Formatter**: A strictly opinionated formatter (like `gofmt` or `prettier`) for Zsh.

## Progress Tracking

**Current release:** v1.0.16 — 1000 katas.

```
[================================================================================] 1000/1000
```

v1.x hotfixes (parser, dedup, severity rebalance, XDG config, inline disable directives, Marketplace action rename) shipped between v1.0.0 and the current tag.
See `CHANGELOG.md` for the full per-release history.

For the list of currently implemented Katas, please refer to [KATAS.md](KATAS.md).