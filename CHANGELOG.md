# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.20] - 2026-03-30

### Added
- **Kata ZC1098**: Use `(q)` flag for quoting variables in `eval`.
- **Kata ZC1099**: Use `(f)` flag to split lines instead of `while read`.
- **Kata ZC1100**: Use parameter expansion instead of `dirname`/`basename`.
- **Kata ZC1101**: Use `$(( ))` instead of `bc` for simple arithmetic.
- **Kata ZC1102**: Redirecting output of `sudo` does not work as expected.
- **Kata ZC1103**: Suggest `path` array instead of `$PATH` string manipulation.
- **Kata ZC1104**: Suggest `path` array instead of `export PATH` string manipulation.
- **Kata ZC1105**: Avoid nested arithmetic expansions for clarity.
- **Kata ZC1106**: Avoid `set -x` in production scripts for sensitive data exposure.
- **Kata ZC1107**: Use `(( ... ))` for arithmetic conditions.
- **Kata ZC1108**: Use Zsh `${(U)var}`/`${(L)var}` case conversion instead of `tr`.
- **Kata ZC1109**: Use parameter expansion instead of `cut` for field extraction.
- **Kata ZC1110**: Use Zsh subscripts instead of `head -1` or `tail -1`.
- **Kata ZC1111**: Avoid `xargs` for simple command invocation.
- **Kata ZC1112**: Avoid `grep -c` -- use Zsh pattern matching for counting.
- **Kata ZC1113**: Use `${var:A}` instead of `realpath` or `readlink -f`.
- **Kata ZC1114**: Consider Zsh `=(...)` for temporary files instead of `mktemp`.
- **Kata ZC1115**: Use Zsh string manipulation instead of `rev`.
- **Kata ZC1116**: Use Zsh multios instead of `tee`.
- **Kata ZC1117**: Use `&!` or `disown` instead of `nohup`.
- **Kata ZC1118**: Use `print -rn` instead of `echo -n`.
- **Kata ZC1119**: Use `$EPOCHSECONDS` instead of `date +%s`.
- **Kata ZC1120**: Use `$PWD` instead of `pwd`.

### Fixed
- **CI**: Deleted unsigned and draft releases for OpenSSF Scorecard Signed-Releases compliance.
- **CI**: Updated auto-approve workflow to use `redteamx` PAT for Code-Review compliance.
- **CI**: Updated release-drafter to use `$RESOLVED_VERSION` for version consistency.

## [0.1.1] - 2025-11-27

### Changed
- **Versioning**: Aligned version number with the total count of implemented Katas (101 Katas = v0.1.1).
- **Core**: Updated Go version to 1.25.
- **Core**: Fixed critical AST type definitions and parser integration issues.

### Added
- Implemented additional Katas to reach a total of 101.

## [0.0.74] - 2025-11-24

### Added
- **Kata ZC1004**: Use `return` instead of `exit` in functions.
- **Kata ZC1016**: Use `read -s` when reading sensitive information.
- **Kata ZC1074**: Prefer modifiers `:h`/:`t` over `dirname`/`basename`.
- **Kata ZC1075**: Quote variable expansions to prevent globbing.
- **Kata ZC1076**: Use `autoload -Uz` for lazy loading.
- **Kata ZC1077**: Prefer `${var:u/l}` over `tr` for case conversion.
- **Kata ZC1078**: Quote `$@` and `$*` when passing arguments.
- **Kata ZC1097**: Declare loop variables as `local` in functions.
- **Kata ZC1079**: Quote RHS of `==` in `[[ ... ]]` to prevent pattern matching.
- **Kata ZC1080**: Use `(N)` nullglob qualifier for globs in loops.
- **Kata ZC1081**: Use `${#var}` to get string length instead of `wc -c`.
- **Kata ZC1082**: Prefer `${var//old/new}` over `sed` for simple replacements.
- **Documentation**: Added `TROUBLESHOOTING.md`, `GOVERNANCE.md`, `COMPARISON.md`, `GLOSSARY.md`, `CITATION.cff`.
- **Documentation**: Expanded `KATAS.md` with new Katas.

### Fixed
- **Parser**: Fixed regression in arithmetic command parsing impacting tests.

## [0.0.72] - 2024-05-20

### Added
- Initial release with 72 implemented Katas.
- Basic Lexer, Parser, and AST implementation for Zsh.
- Text and JSON reporters.
- Integration tests framework.