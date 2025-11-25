# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.0.74] - 2025-11-24

### Added
- **Kata ZC1004**: Use `return` instead of `exit` in functions.
- **Kata ZC1016**: Use `read -s` when reading sensitive information.
- **Kata ZC1074**: Prefer modifiers `:h`/:`t` over `dirname`/`basename`.
- **Kata ZC1075**: Quote variable expansions to prevent globbing.
- **Kata ZC1076**: Use `autoload -Uz` for lazy loading.
- **Kata ZC1077**: Prefer `${var:u/l}` over `tr` for case conversion.
- **Kata ZC1078**: Quote `$@` and `$*` when passing arguments.
- **Kata ZC1079**: Quote RHS of `==` in `[[ ... ]]` to prevent pattern matching.
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