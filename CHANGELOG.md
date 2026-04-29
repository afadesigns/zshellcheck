# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.17] - 2026-04-30

### Added
- OpenSSF Best Practices badge promoted from silver to **gold**. Project ID 12657. Self-attestation reflects 95.0% project-wide statement coverage (Codecov 3-OS union) and an unassociated significant-contributor count above the gold floor.
- `pkg/katas/` and `pkg/katas/katatests/` source layout consolidated into hundred-bucket files (`zc1000s.go`, `zc1000s_test.go`, …) so GitHub's 1,000-entry directory cap stops truncating the source view. Helpers live in `scripts/consolidate-katas-source.py` and `scripts/consolidate-katatests.py`.
- CI gold-tier companion gates: gocyclo regression baseline (`.github/gocyclo-baseline`), SPDX + copyright header check on every tracked `.go` file, project-wide coverage floor at 95.0% (raised from 90.0%), `pkg/katas/` directory-cap check, libraries.io stale-dependency surfacing.
- Coverage drive: 80+ targeted unit tests across `pkg/parser/`, `pkg/lexer/`, and `pkg/katas/` raising the Codecov 3-OS union from 94.40% to 95.05%. New sweep tests cover every `offsetLineColZCxxxx` and `byteOffsetToLineColZCxxxx` duplicate, every `Check` and `Fix` type-guard branch, and the `parseCommandPipeline` redirection loop reachable only via the `&&`/`||` chain. (#1320, #1322, #1323, #1324, #1325, #1326)

### Changed
- `golang.org/x/sys` bumped from 0.26.0 to 0.43.0; `go` directive raised to 1.25.0; release pipeline pin updated to match.
- `goreleaser/goreleaser-action` 7.1.0 → 7.2.1, `release-drafter/release-drafter` weekly tracking pin, `crate-ci/typos` 1.45.1 → 1.45.2. (#1318)
- Codecov gate raised 90.0% → 95.0%; the local Linux single-OS sanity floor stays at 90.0% because the union-merge across Linux + macOS + Windows runs ~3-4 points above the single-OS measurement. (#1320)
- `codecov.yml` badge gradient flipped: `range: 80..95` so 95% renders as the GREEN end of the gradient. The previous `range: 95..100` placed the badge at red even when the status check passed. (#1327)
- Codecov upload step in `ci.yml` skips on dependabot PRs; the action's `fail_ci_if_error: true` was failing every dep-bump build because dependabot runs lack `secrets.CODECOV_TOKEN`. Human PRs still upload and enforce the threshold. (#1319)

### Fixed
- Lexer: comment-skip in `tryShebangOrComment` recursed via `return l.NextToken()`, causing the named-return defer to run twice on the same token. Each `${` after a skipped comment double-incremented `dollarBraceDepth`, leaving it stuck at 1; subsequent `#` bytes were classified as `${…}` length operators rather than comment openers, so inline comments after `;;` (or any other comment mid-statement) leaked into the parser as `# IDENT` tokens. Replaced with a loop in `NextToken`; `tryShebangOrComment` now returns `(zero token, true)` when a comment was skipped. Drains 12 latent errors in `fzf-tab/lib/zsh-ls-colors/ls-colors.zsh`. (#1316)
- Parser: typed-nil concrete pointer wrapped in an `ast.Statement` interface (e.g. `(*ast.ForLoopStatement)(nil)` returned from a sub-parser recovery path) panicked inside `keywordStmtToExpression` when `TokenLiteralNode` dispatched on a nil receiver. Added `isTypedNilStatement` to detect each variant and degrade to a stub Identifier so pipeline chaining continues. Surfaced by a 50-byte fuzz-minimised input. (#1314)
- Parser: `${1:-default}` / `${1:=default}` / `${1:+alt}` / `${1:?err}` and the rest of the modifier-tail forms over a positional-parameter subject parse cleanly. The tokenizer emits the digit subject as `INT` followed by `:`, so the modifier-tail walker now drives instead of bailing on `expected next token to be }`. Fixes #129.
- Parser: 38 → 0 errors across the pinned-corpus baseline (170 → 0 over the full drainage cycle, 37 PRs). All 78 corpus files parse without error. Notable individual fixes: `${pos:mod}` (#129), `${#}` special parameter, zsh `if cond cmd` shortcut, glob `#`/`##` qualifier in command words, `;|`/`;&` case fall-through, glob bracket-class in `[[ … ]]`, `function` keyword as assignment rhs, array literal `)` not subshell terminator, `(((` subshell+arith vs arith+group disambiguation, `((` after newline fuses to `DoubleLparen`, `<(`/`>(` not process-sub inside `[[ … ]]`, `$(( … ))` flags `consumedParenTerminator`, embedded `$(…)` walker in `"…"` strings, `|&` stderr-pipe fusion, PIPE as bitwise-OR in arithmetic, `typeset` array literal terminator handling.

### Closed
- #129 (`[BUG] incomplete support for variable substitution`) — resolved by the modifier-tail walker for positional subjects.
- #1232 (`Panic in ZC1122: typed-nil *ast.Identifier dereferenced`) — `checkZC1122` carries the `ident == nil` guard; the matching pattern in pipeline-head wrap is covered by #1314.

## [1.0.16] - 2026-04-26

### Added
- **Auto-fix coverage now at 137/1000 katas.** New rewrites since 1.0.15, every one deterministic and idempotent on a re-run:
  - `ZC1015` backticks → `$(...)`.
  - `ZC1016` inserts `-s` after `read` when the variable looks sensitive (`password`, `secret`, `token`, …).
  - `ZC1008` and `ZC1022` share ZC1013's `let NAME=EXPR` → `(( NAME = EXPR ))` rewrite.
  - `ZC1032` `let i=i+1` → `(( i++ ))` (and `i-1` → `i--`).
  - `ZC1037` shares ZC1092's `echo` → `print -r --` rewrite.
  - `ZC1043` prepends `local ` to unscoped function-body assignments.
  - `ZC1053` inserts `-q` after `grep` / `egrep` / `fgrep` / `zgrep` when used in an `if` or `while` condition.
  - `ZC1069` rewrites `local` to `typeset` when used at file scope (outside any function).
  - `ZC1095` `seq N` → `{1..N}` (reuses the ZC1061 brace-expansion rewrite).
  - `ZC1034` / `ZC1271` `which` → `command -v`.
  - `ZC1107` `[[ a -lt b ]]` → `(( a < b ))`.
  - `ZC1146` `cat F | sed/awk/sort/head/tail` → `tool ... F`.
  - `ZC1153` `diff -q F1 F2` → `cmp -s F1 F2`.
  - `ZC1155` `which -a` → `whence -a`.
  - `ZC1163` `grep PAT | head -1` → `grep -m 1 PAT`.
  - `ZC1172` `read -a ARR` → `read -A ARR` (Zsh array form).
  - `ZC1190` `grep -v p1 | grep -v p2` → `grep -v -e p1 -e p2`.
  - `ZC1191` `clear` → `print -rn $'\e[2J\e[H'`.
  - `ZC1201` `rsh`/`rlogin` → `ssh`, `rcp` → `scp`.
  - `ZC1202` `ifconfig` → `ip addr`.
  - `ZC1203` `netstat` → `ss`.
  - `ZC1216` `nslookup` → `host`.
  - `ZC1219` `wget -O- URL` / `wget -qO- URL` → `curl -fsSL URL`.
  - `ZC1215` `cat /etc/{os,lsb}-release` → `. /etc/{os,lsb}-release` (single-arg only).
  - `ZC1217` shares the `service UNIT VERB` → `systemctl VERB UNIT` rewrite with ZC1512.
  - `ZC1230` `ping URL` → `ping -c 4 URL`.
  - `ZC1235` `git push -f` → `git push --force-with-lease`.
  - `ZC1238` strips `-it` from `docker exec`.
  - `ZC1239` strips `-it` from `kubectl exec`.
  - `ZC1252` `cat /etc/{passwd,group,shadow}` → `getent {passwd,group,shadow}`.
  - `ZC1255` `curl URL` → `curl -L URL`.
  - `ZC1257` `docker stop X` → `docker stop -t 10 X`.
  - `ZC1260` `git branch -D` → `git branch -d`.
  - `ZC1268` inserts `--` before the first non-flag arg of `du -sh`.
  - `ZC1273` `grep PAT FILE /dev/null` → `grep -q PAT FILE` (insert `-q`, drop `/dev/null`).
  - `ZC1276` `seq M N` → `{M..N}`.
  - `ZC1293` `test EXPR…` → `[[ EXPR… ]]` (also wired to ZC1006 / ZC1020 / ZC1036 which fire on the same shape).
  - `ZC1279` `readlink -f PATH` → `realpath PATH` when `-f` is the first argument.
  - `ZC1297` `$BASH_SOURCE` → `${(%):-%x}`.
  - `ZC1319` `$BASH_ARGC` → `$#`.
  - `ZC1320` `$BASH_ARGV` → `$argv`.
  - `ZC1334` `type -p` / `type -P` → `whence -p`.
  - `ZC1374` `$FUNCNEST` → `${#funcstack}` inside echo / print / printf args.
  - `ZC1377` `BASH_ALIASES` → `aliases` inside echo / print / printf string args.
  - `ZC1378` uppercase `DIRSTACK` → `dirstack` inside echo / print / printf string args.
  - `ZC1381` `$COMP_WORDS` / `$COMP_CWORD` / `$COMP_LINE` / `$COMP_POINT` → `$words` / `$CURRENT` / `$BUFFER` / `$CURSOR` inside echo / print / printf args.
  - `ZC1382` `$READLINE_LINE` / `$READLINE_POINT` / `$READLINE_MARK` → `$BUFFER` / `$CURSOR` / `$MARK` inside echo / print / printf args.
  - `ZC1380` `export HISTIGNORE=…` → `export HISTORY_IGNORE=…`.
  - `ZC1383` `TIMEFORMAT` → `TIMEFMT` inside echo / print / printf string args.
  - `ZC1394` `$BASH` → `$ZSH_NAME` inside echo / print / printf string args.
  - `ZC1403` `HISTFILESIZE` → `SAVEHIST` inside echo / print / printf / export args.
  - `ZC1404` `$BASH_CMDS` → `$commands` inside echo / print / printf string args.
  - `ZC1413` `hash -t cmd` → `whence -p cmd` (rename + flag swap).
  - `ZC1411` `enable -n NAME` → `disable NAME`.
  - `ZC1448` inserts `-y` after `apt install` / `apt upgrade` / `apt dist-upgrade` / `apt full-upgrade`.
  - `ZC1502` inserts `-- ` before the first `$var` argument of `grep` / `egrep` / `fgrep` / `rg` / `ag` to block flag injection.
  - `ZC1501` `docker-compose` → `docker compose`.
  - `ZC1512` `service UNIT VERB` → `systemctl VERB UNIT` (rename + arg swap).
  - `ZC1565` `whereis` / `locate` / `mlocate` / `plocate` → `command -v`.
  - `ZC1591` `printf '%s\n' "${array[@]}"` → `print -l -r -- "${array[@]}"`.
  - `ZC1637` `readonly NAME=value` → `typeset -r NAME=value`.
  - `ZC1643` `$(cat FILE)` → `$(<FILE)` inside SimpleCommand argument strings.
  - `ZC1675` `export -f FUNC` → `typeset -fx FUNC`, `export -n VAR` → `typeset +x VAR`.
  - `ZC1685` `sleep infinity` → `exec tail -f /dev/null`.
  - `ZC1717` strips `--disable-content-trust` from `docker pull` / `push` / `build` / `create` / `run`.
  - `ZC1773` `xargs CMD` → `xargs -r CMD`.
  - `ZC1334` collapses `type -p`'s flag with the rename so it wins over `ZC1064`'s narrower `type` → `command -v` form.
  - `ZC1013` defers to `ZC1032` on the increment/decrement shape so the rewrite uses the C-style operator instead of the literal `(( name = name+1 ))` form.

### Changed
- `ZC1005`'s `which` → `whence` rewrite now yields `command -v` for the bare-statement case because the new `ZC1034` fix arrives ahead in walk order. Inside backticks / `$(...)`, `whence` still wins because the parent `ExpressionStatement` is absent.
- `ZC1263`'s `apt` → `apt-get` rewrite for `apt install` now runs alongside `ZC1448`'s `-y` insertion, producing `apt-get -y install ...` in a single pass.

## [1.0.15] - 2026-04-25

### Breaking
- **Inline directive renamed `# zshellcheck disable=…` → `# noka`.**
  The legacy form is no longer recognised — every silenced violation needs its directive rewritten.
  Three forms remain available:
    - `cmd  # noka`              — silences every kata on this line.
    - `cmd  # noka: ZC1234`      — silences one kata on this line.
    - `cmd  # noka: ZC1234, ZC1075` — multiple, comma- or space-separated.
  Standalone `# noka` directives still apply to the next non-blank code line; placed at file tail with no code after them, they apply file-wide.
  Rationale: shorter (18 vs 30 chars), distinctive ("no kata"), aligns with the python ecosystem's `# noqa` convention.
  Refactored cleanly while the project is still early — no fork-side migration to coordinate.

### Added
- **`-no-banner` CLI flag.** Suppresses the startup banner.
  Useful for CI runs, scripted invocations, and embedding zshellcheck output inside other tools where the banner is noise.
  Banner remains the default for interactive use.
- **Auto-fix coverage expanded to 67 katas.** The first-wave shipped 3 (`ZC1002`, `ZC1005`, `ZC1092`); the registry now ships rewrites for parameter-name renames (`$BASH_ALIASES` → `$aliases`, `$BASH_REMATCH` → `$match`, `$BASH_VERSION` → `$ZSH_VERSION`, `$TIMEFORMAT` → `$TIMEFMT`, `$BASH_CMDS` → `$commands`, ZSH array helpers, etc.), command/flag rewrites (`echo -E` → `print -r`, `read -a` → `read -A`), and several modernisations across the ZC1300 series.
- **Demo recording.** `docs/assets/demo.gif` showcases the lint → diff → fix → re-lint loop on a sample script, embedded in the README.
  Source tape at `docs/assets/demo.tape` for reproducible re-renders via [VHS](https://github.com/charmbracelet/vhs).
- **CLI banner refreshed.** Terminal-art rendering of the project logo replaces the prior block-letter ASCII.
  Tagline matches the project slogan: `The quiet linter for a quiet shell.`
- **`KATAS.md` shows fix coverage.** Generator now emits an `Auto-fix: yes/no` line per entry and a `with auto-fix` row in the summary table.
- **`--help` redesign.** Flags grouped by intent (OUTPUT / FILTER / AUTO-FIX / DIAGNOSTICS), six-entry EXAMPLES block, ANSI colour gated on TTY + `NO_COLOR`.
- **Windows installer (`install.ps1`).** PowerShell 5.1+ compatible, mirrors `install.sh`: SHA-256 + cosign verification, user-scoped install into `%LOCALAPPDATA%\Programs\zshellcheck\bin`, clean `-Uninstall` reversal.
- **Linux native packages.** Goreleaser `nfpms:` block emits signed `.deb`, `.rpm`, and `.apk` artifacts on every tag, alongside the existing tarballs.
- **Multi-arch Docker image.** `FROM scratch`, ≈ 2 MB, published to `ghcr.io/afadesigns/zshellcheck` for `linux/amd64` and `linux/arm64`. Manifest signed with cosign.
- **`INSTALL.md` canonical install guide.** Single source of truth split into macOS / Windows / Linux / Cross-platform sections with explicit uninstall paths everywhere.
- **`INTEGRATIONS.md` at repo root.** Per-project tables (frameworks, plugin managers, tooling, plugins, prompts) plus the targeted-next list.
- **Homebrew-eligibility tracker.** New weekly workflow opens a tracking issue when stargazers / forks / watchers cross the third-party homebrew-core notability threshold.

### Changed
- **Column pointer character.** Lint output now uses `↑` (U+2191) under the offending column instead of `^`.
  Matches the convention modern compilers (rustc, swiftc) use when pointing to a column.
- **Multi-pass `-fix`.** `applyFixesUntilStable` now loops `fix.Apply` while `collectEdits` keeps producing edits, capped at five passes.
  Nested rewrites (e.g.
  `` `which git` `` → `$(whence git)`) resolve in a single `-fix` invocation.
- **Fix summary footer.** Multi-file `-fix` runs now print `fix summary: N edit(s) across M file(s) (scanned K)` to stderr.
  Single-file invocations stay silent for backward compatibility.
- **Repo description + homepage.** Synced to the locked README slogan/subheader; homepage points at the repo root.
- **README structure.** New quick-link nav row above the fold, demo GIF replaces the prior Katas-at-a-glance teaser, Install section rebuilt around the three primary channels (macOS/Linux, Windows, Go), Integrations table trimmed to a featured spotlight + link to `INTEGRATIONS.md`, Documentation table reordered with `INSTALL.md` and `INTEGRATIONS.md` added.
- **All top-level docs reflowed to semantic line breaks.** One sentence per source line. Rendered output unchanged (CommonMark soft breaks); diff hygiene improved.

### Fixed
- **Typed-nil `ast.Node` handling.** Guarded `Walk` against typed-nil interface values so downstream visitors no longer panic on partially-constructed trees produced by parser recovery paths.
- **32 kata nil-guards.** Hardened `*ast.Identifier` dereferences across ZC1122, ZC1191, and the ZC1297–ZC1333 series so external projects that produce nil identifier values no longer crash the linter.
- **Parser compat.** Routed `cmd --flag arg` through the simple-command path (refactored 23 mangled-name katas onto `FlagArgPosition`); added bracket-cond glob-alt literal support; fixed inner `$()` `RPAREN` containment via `consumedParenTerminator`; multi-line `$(…)` newline drain; case-pattern glob-alt label advance.

## [1.0.14] - 2026-04-24

### Added
- **Auto-fixer core.** New `pkg/fix` package applies per-kata `Fix` edits to source files.
  Handles 1-based line/column to byte-offset resolution, conflict resolution when edits overlap (outer span wins, inner picked up on rerun), and a built-in unified-diff renderer for preview mode.
- **CLI fix flags** — `-fix` (apply in place), `-diff` (preview as unified diff), `-dry-run` (with `-fix`, report without writing).
  File permissions are preserved across in-place rewrites.
- **Kata `Fix` hook.** `Kata` now carries an optional `Fix func(ast.Node, Violation, []byte) []FixEdit`.
  Checks that declare a Fix participate in auto-fixing; those that do not continue to lint-only.
- **First-wave Fix coverage:**
  - `ZC1002` — `` `cmd` `` to `$(cmd)`.
  - `ZC1005` — `which` to `whence`.
  - `ZC1092` — `echo` to `print -r --` for the no-flag form.
- **Zsh-ecosystem compatibility harness.** New `scripts/test-zsh-compat.sh` clones a corpus of well-known Zsh projects (oh-my-zsh, powerlevel10k, prezto, zsh-autosuggestions, zsh-syntax-highlighting, zsh-completions, spaceship-prompt) into `testdata/external-corpora/` (git-ignored) and reports parser errors plus violation summaries.

### Changed
- `CheckAndFix` registry method added alongside `Check` so the walker can collect violations and their fix edits in a single pass.

## [1.0.13] - 2026-04-22

### Fixed
- **Parser**: bare `$+name` / `$+name[key]` inside `(( … ))` no longer errors with `expected next token to be IDENT, got + instead`.
  Equivalent shape to the working `${+name[key]}` path.
  (#1047)
- **Parser**: `(( A )) && (( B ))` / `||` chains (and mixed) no longer error with `no prefix parse function for && found`.
  Logical operators after an arithmetic command now parse into a normal `InfixExpression`.
  (#1047)

### Changed
- `.pre-commit-hooks.yaml` — `language: go` → `language: golang`, the canonical pre-commit language identifier.
  Fixes installation under `prek`.
  (#1046)

## [1.0.12] - 2026-04-20

### Changed
- `action.yml` — extend `name` to `ZshellCheck v1` (the `ZshellCheck` form from v1.0.11 still collided with an existing Marketplace registry entry).
  The action identifier (`afadesigns/zshellcheck@vX.Y.Z`) is unchanged.

## [1.0.11] - 2026-04-20

### Changed
- `action.yml` — rename `name` from `ZShellCheck` to `ZshellCheck` (lowercase `h`).
  Marketplace requires a unique action name; the original capitalization collided with an existing registry entry.
  The action identifier (`afadesigns/zshellcheck@vX.Y.Z`) is unchanged.

## [1.0.10] - 2026-04-20

**Versioning switch.** The kata-count formula (MAJOR = count/1000,
MINOR = (count%1000)/100, PATCH = count%100) retires here.
Going
forward ZShellCheck follows standard [semantic versioning](https://semver.org).
`pkg/version/version.go` is now hand-maintained; `scripts/update-version.sh`,
`scripts/HOTFIX`, and the `update-version` pre-commit hook are removed.
The `tag-release` workflow no longer auto-bumps PATCH on every main
push — tags are cut manually.

## [1.0.9] - 2026-04-20

Hotfix rollup for every fix merged between the v1.0.0 tag and the
GitHub Marketplace debut.
No new katas — kata count stays at 1000.

### Added
- **Inline `# zshellcheck disable=…` directives** — suppress katas
  per-line (trailing comment), per-next-line (standalone comment above
  code), or file-wide (standalone comment with no following code).
  Multiple IDs may be comma- or whitespace-separated.
Merges with the
  config-level `disabled_katas` list.
(#127)
- **XDG Base Directory support** — `$XDG_CONFIG_HOME/zshellcheck/config.{yml,yaml}`
  is now merged with `~/.zshellcheckrc` and `./.zshellcheckrc`, with the
  project-local file winning.
(#309)
- **`scripts/HOTFIX` offset file** — tracks monotonic patch-release
  bumps so hotfix releases can ship between kata additions without
  colliding with the kata-count formula.

### Fixed
- **Parser panic on `dd if=src of=dst`** — the lexer now demotes
  keyword tokens to `IDENT` when immediately followed by `=`, so
  `if=`, `of=`, etc. parse as ordinary key=value argument pairs.
(#435)
- **`elif` chain mis-nesting** — `parseIfStatement` now terminates the
  consequence block on `ELIF` and builds a right-nested `IfStatement`
  chain, fixing false-positives on multi-branch conditionals.
(#126)
- **Parser crash on `/dev/sdX` literals** in arithmetic and redirection
  positions.
(#347)
- **14 duplicate katas retired** as no-op stubs (ZC1022–1029, 1033,
  1035, 1018, 1019, 1277, 1278).
Canonical detections remain in the
  surviving sibling IDs; retired IDs still parse in `.zshellcheckrc`
  so legacy configs keep working.
(#341–#345)
- **5 overlapping kata pairs narrowed** — ZC1441 skips when
  `--volumes` is present (ZC1545 owns that case); ZC1978 narrows to
  `tftp` (ZC1200 owns `ftp`); ZC1327 drops `-c`/`-d` (ZC1487 owns);
  ZC1826 drops numeric modes (ZC1892 owns); ZC1999 rewritten as a
  typo-detector pointing at ZC1934 `AUTO_NAME_DIRS`.
- **10 Style katas rebalanced to Warning** — ZC1075, 1078, 1079, 1084,
  1085, 1090, 1136, 1139, 1141, 1258.
These flag patterns with real
  correctness or safety impact, not cosmetic preference.
(#346)

### CI
- **OSV-Scanner** workflow drops the removed `--skip-git` flag that
  OSV-Scanner v2 rejects; PRs no longer pre-fail on the vulnerability
  scan step.
- **`golangci-lint-action` pinned at v6.5.2** until the v1→v2 config
  migration lands — avoids surfacing ~20 pre-existing `staticcheck`
  QF1001-QF1003 findings that block unrelated Dependabot bumps.

### Documentation
- USER_GUIDE gains an **Inline Disable Directives** section covering
  the three directive forms.
- Author identity corrected across CITATION.cff, SECURITY.md, and
  CODE_OF_CONDUCT.md — contact email is now `github@afadesign.co`.

## [1.0.0] - 2026-04-20

**1000 Katas milestone.** The kata-count formula (MAJOR = count/1000,
MINOR = (count%1000)/100, PATCH = count%100) now resolves to exactly
`1.0.0`.
This is the first stable release of ZShellCheck, targeted at
the GitHub Marketplace launch.

### Added
- **665 new Katas** bringing the total from 335 (v0.3.35) to **1000**
  (ZC1339 through ZC2003).
Highlights:
  - **Zsh semantics & `setopt` subtleties** — `PROMPT_SUBST`,
    `GLOBAL_RCS`, `POSIX_IDENTIFIERS`, `CHASE_DOTS`, `SH_FILE_EXPANSION`,
    `CSH_JUNKIE_QUOTES`, `REMATCH_PCRE`, `KSH_TYPESET`, `BRACE_CCL`,
    `CSH_NULLCMD`, `AUTO_NAMED_DIRS`, `EVAL_LINENO`, `KSH_ZERO_SUBSCRIPT`,
    `HIST_NO_FUNCTIONS`, `HIST_FCNTL_LOCK`, `BG_NICE`, and many more.
  - **Storage & filesystem safety** — `zpool import -f`/`export -f`,
    `dmsetup remove_all`, `losetup -P`/`kpartx -a`/`partprobe`,
    `sgdisk -Z`/`-o`, `lvreduce -f`/`-y`, `exportfs -au`.
  - **Kernel/devices** — `udevadm trigger --action=remove`,
    `tpm2_clear`, `ipcrm -a`, `unshare -U`/`-r`.
  - **Platform ops** — `crictl rmi -a`/`rm -af`,
    `kubectl taint nodes …:NoExecute`, `dnf/yum versionlock add`.
  - **Shell hygiene** — `zsh -f`/`-d` bypassing startup files,
    `exec -a NAME` masking `argv[0]`, `touch -d`/`-t`/`-r` timestamp
    rewrite, `nsupdate -y` TSIG-in-argv, `openssl passwd -crypt`/`-1`/`-apr1`,
    `ftp`/`tftp` plaintext, `pkexec` script elevation.
- **Test triplet per kata** — `pkg/katas/katatests/zc####_test.go` with
  valid + invalid cases across every new ID.
- **Misspell ignore-words entry for `exportfs`** — prevents false
  positives on legitimate NFS-tool references.

### Changed
- README, USER_GUIDE, REFERENCE, ROADMAP, CITATION.cff, and the
  `zshellcheck(1)` man page updated for v1.0.0 and the 1000-kata total.
- `-severity`, `--no-color`, `--version`, and `-format sarif` now
  documented in the man page.

### Documentation
- CHANGELOG gains this 1.0.0 section covering the ZC1339–ZC2003 range.
- ROADMAP marks the 250, 500, and 1000-kata milestones complete and
  advances the LSP / auto-fixer / plugin work into the 1.x bucket.

## [0.3.35] - 2026-04-17

**Public beta.** First release with successfully built, signed, and
attested binaries.
The kata-count formula now maps correctly to the
`MAJOR.MINOR.PATCH` scheme (335 katas → 0.3.35); prior tag series was
produced before the release pipeline was functional and contained no
published artifacts.

### Added
- **169 new Katas** (ZC1170 through ZC1338) spanning: Zsh built-in
  preferences over external commands, parameter-expansion
  alternatives to `cut`/`sed`/`tr`/`sort`, Bash-ism detection for
  portability, and git/docker/grep flag recommendations.
- **SBOM workflow**: weekly SPDX + CycloneDX generation, attested on
  main for SLSA-aligned build metadata.
- **OSV-Scanner**: daily scheduled scan of Go module dependencies,
  SARIF uploaded to Security tab.
- **Nightly fuzz**: 5-minute budget per lexer/parser target with
  corpus caching and failing-corpus artifact upload.
- **Release pipeline**: Syft for SBOM in archives, Cosign keyless
  signing, `actions/attest-build-provenance` for SLSA provenance,
  locked to goreleaser `~> v2`.

### Changed
- **Codecov gate**: raised project + patch targets from 80% to 95%.
- **Release archives**: now include man page, Zsh and Bash shell
  completions, and CHANGELOG alongside LICENSE and README.
- **Label taxonomy**: removed duplicates (`documentation`,
  `enhancement`), renamed `starter` to `good first issue` for
  GitHub discoverability, added severity/component/topic/status
  labels.

### Fixed
- **install.sh** (#311): graceful HTTP error handling when
  `/releases/latest` returns 404, with `/tags` fallback; explicit
  error message when no tag can be resolved instead of silent exit.
- **Release workflow**: installed missing Syft and Cosign on the
  runner; locked goreleaser-action version syntax to `~> v2` to
  avoid ambiguous `latest`.
- **GoReleaser config**: updated `archives.format_overrides.format`
  to the plural `formats` key per goreleaser v2.

### Known Issues
Filed during release prep for community resolution:
- **#341-#345**: duplicate kata detections (ZC1038/ZC1093,
  ZC1005/ZC1019, ZC1009/ZC1018/ZC1278, ZC1108/ZC1277, and a 10-kata
  cluster in ZC1022-ZC1029 + ZC1033/ZC1035).
- **#346**: 11 Style-severity katas likely better categorised as
  Warning (quoting defects, `rm -rf`, `curl | sh`, `source` with URL).
- **#347**: parser crashes on common Zsh constructs
  (`for...in...do...done`, `||`) — reduces detection fidelity on
  real-world `.zshrc` files.

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
- **Kata ZC1074**: Prefer modifiers `:h` / `:t` over `dirname` / `basename`.
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