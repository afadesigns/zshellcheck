# Roadmap

## Short Term
- [x] Stabilize parser for basic Zsh constructs (loops, arithmetic, conditions).
- [x] Implement initial set of high-value Katas (ZC1001-ZC1037).
- [x] Integrate with `pre-commit`.
- [ ] Expand test coverage for parser edge cases.

## Medium Term
- [ ] **Variable Expansion Parsing:** Deep support for Zsh's complex parameter expansion `${name: ...}`.
- [ ] **Globbing Analysis:** better understanding of extended glob patterns.
- [ ] **Autofix:** Implement automatic fixing for simple violations (e.g., replacing `[` with `[[`).

## Long Term
- [ ] Full Zsh Grammar support (including obscure builtins and modifiers).
- [ ] Type inference for variables (to detect array vs scalar misuse).
- [ ] LSP (Language Server Protocol) implementation for editor integration.
