# Roadmap

## Versioning Strategy
We aim to implement at least **1000 Katas** (checks).
- **Version 1.0.0** will be released upon reaching Kata #1000 (ZC2000).
- Version format: `Major.Minor.Patch`.
- `Major`: Thousands of Katas (0 for <1000).
- `Minor`: Hundreds of Katas.
- `Patch`: Tens of Katas.
- The exact version is updated with every Kata implemented, appending the units digit if necessary (e.g., `0.0.51`).

**Current Progress:** 67 Katas implemented (Version 0.0.67).

## Milestones

### Phase 1: Core Stability & Basic Katas (0.0.0 - 0.1.0)
- [x] Stabilize Parser (Arithmetic, Conditions, Loops, Command Substitution).
- [x] Robust Integration Test Suite (`tests/integration_test.zsh`).
- [x] Basic Checks (ZC1001-ZC1050).
- [ ] Expand to 100 Katas (ZC1051-ZC1100).

### Phase 2: Advanced Analysis (0.1.0 - 0.5.0)
- [ ] **Deep Variable Expansion Parsing:** Support `${name: ...}`, `${(f)...}`, etc.
- [ ] **Globbing Analysis:** Extended glob patterns.
- [ ] **Data Flow Analysis:** Track variable types (array vs scalar) and potential values.
- [ ] **Autofix:** Automatic code correction.

### Phase 3: Maturity (0.5.0 - 1.0.0)
- [ ] Full Zsh Grammar Support.
- [ ] LSP Server Implementation.
- [ ] Plugin Architecture.
- [ ] 1000+ Katas.

## Planned Katas (Upcoming)
- **ZC1051:** Check for `rm` variable expansion safety.
- **ZC1052:** Warn about `sed -i` portability.
- **ZC1053:** Prefer `builtin cd` or check `cd` behavior.
- **ZC1054:** Check for valid shebangs (more strict).
- **ZC1055:** Warn about `echo` flags (portability).