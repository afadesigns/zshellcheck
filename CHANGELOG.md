# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.66] - 2026-03-30

### Added
- **Kata ZC1169**: Avoid `install` for simple copy+chmod -- use `cp` then `chmod`.

### Fixed
- **CI**: Add top-level read permissions to workflows for Scorecard compliance.
- **CI**: Pin govulncheck to v1.1.4 for Scorecard compliance.

### Changed
- **CI**: Optimize CI with concurrency, dependabot limits, and updated hooks.
- **Docs**: Update contact email to afadesign.official@gmail.com.
- **Docs**: Add Contributors section to README.

## [0.1.58] - 2026-03-30

### Added
- **Kata ZC1121**: Use `$HOST` instead of `hostname`.
- **Kata ZC1122**: Use `$USER` instead of `whoami`.
- **Kata ZC1123**: Use `$OSTYPE` instead of `uname`.
- **Kata ZC1124**: Use `: > file` instead of `cat /dev/null > file` to truncate.
- **Kata ZC1125**: Avoid `echo | grep` for string matching.
- **Kata ZC1126**: Use `sort -u` instead of `sort | uniq`.
- **Kata ZC1127**: Avoid `ls` for counting files.
- **Kata ZC1128**: Use `> file` instead of `touch file` for creation.
- **Kata ZC1129**: Use Zsh `stat` module instead of `wc -c` for file size.
- **Kata ZC1131**: Avoid `cat file | while read` -- use redirection.
- **Kata ZC1132**: Use Zsh pattern extraction instead of `grep -o`.
- **Kata ZC1133**: Avoid `kill -9` -- use `kill` first, then escalate.
- **Kata ZC1134**: Avoid `sleep` in tight loops.
- **Kata ZC1135**: Avoid `env VAR=val cmd` -- use inline assignment.
- **Kata ZC1136**: Avoid `rm -rf` without safeguard.
- **Kata ZC1137**: Avoid hardcoded `/tmp` paths.
- **Kata ZC1139**: Avoid `source` with URL -- use local files.
- **Kata ZC1140**: Use `command -v` instead of `hash` for command existence.
- **Kata ZC1141**: Avoid `curl | sh` pattern.
- **Kata ZC1142**: Avoid chained `grep | grep` -- combine patterns.
- **Kata ZC1143**: Avoid `set -e` -- use explicit error handling.
- **Kata ZC1144**: Avoid `trap` with signal numbers -- use names.
- **Kata ZC1145**: Avoid `tr -d` for character deletion -- use parameter expansion.
- **Kata ZC1146**: Avoid `cat file | awk` -- pass file to awk directly.
- **Kata ZC1147**: Avoid `mkdir` without `-p` for nested paths.
- **Kata ZC1148**: Use `compdef` instead of `compctl` for completions.
- **Kata ZC1149**: Avoid `echo` for error messages -- use `>&2`.
- **Kata ZC1151**: Avoid `cat -A` -- use Zsh builtins for non-printable characters.
- **Kata ZC1152**: Use Zsh PCRE module instead of `grep -P`.
- **Kata ZC1153**: Use `cmp -s` instead of `diff` for equality check.
- **Kata ZC1154**: Use `find -exec {} +` instead of `find -exec {} \;`.
- **Kata ZC1155**: Use `whence -a` instead of `which -a`.
- **Kata ZC1156**: Avoid `ln` without `-s` for symlinks.
- **Kata ZC1157**: Avoid `strings` command -- use Zsh expansion.
- **Kata ZC1158**: Avoid `chown -R` without `--no-dereference`.
- **Kata ZC1159**: Avoid `tar` without explicit compression flag.
- **Kata ZC1160**: Prefer `curl` over `wget` for portability.
- **Kata ZC1161**: Avoid `openssl` for simple hashing -- use Zsh modules.
- **Kata ZC1162**: Use `cp -a` instead of `cp -r` to preserve attributes.
- **Kata ZC1163**: Use `grep -m 1` instead of `grep | head -1`.
- **Kata ZC1164**: Avoid `sed -n 'Np'` -- use Zsh array subscript.
- **Kata ZC1165**: Use Zsh parameter expansion for simple `awk` field extraction.
- **Kata ZC1166**: Avoid `grep -i` for case-insensitive match -- use `(#i)` glob flag.
- **Kata ZC1167**: Avoid `timeout` command -- use Zsh `TMOUT` or `zsh/sched`.
- **Kata ZC1168**: Use `${(f)...}` instead of `readarray`/`mapfile`.
- **Severity Levels**: All katas now have assigned severity levels (error, warning, info, style).
- **CLI**: Added `--severity` flag for filtering violations by minimum severity.
- **CLI**: Added `--no-color` flag for text reporter.
- **CLI**: Added verbose output mode for text reporter.
- **Reporter**: Added source code context to violation output.
- **Reporter**: Added ANSI colors and file location to output.
- **SARIF**: Added SARIF output format for GitHub Security integration.

### Changed
- **Installer**: Enhanced `install.sh` with pipe support, CI detection, auto-sudo, binary download fallback, checksum verification, mktemp for downloads, banner, and version flags.
- **Installer**: Added man page and shell completion installation.

### Fixed
- **CI**: Fixed OpenSSF Scorecard issues and hardened CI settings.
- **CI**: Pinned all actions to SHAs and restricted permissions.
- **CI**: Added CodeQL analysis and Fuzz testing.
- **CI**: Added dependency review workflow and tuned dependabot.
- **CI**: Multiple golangci-lint configuration fixes.
- **CI**: Fixed badge and release issues.

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
- **CI**: Updated code review workflow for Scorecard Code-Review compliance.
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