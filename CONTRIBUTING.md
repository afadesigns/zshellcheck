# Contributing to ZShellCheck

Thanks for helping improve ZShellCheck. This guide covers the PR workflow, how to add a kata, and the local checks you should run before pushing.

For deeper internals (lexer/parser/AST design, release process, architecture diagrams), see the [Developer Guide](docs/DEVELOPER.md).

## Quick Start

```bash
git clone https://github.com/afadesigns/zshellcheck.git
cd zshellcheck
./install.sh
```

The installer builds from source when run inside the repo, or downloads the signed release binary otherwise. See [Developer Guide — Getting Started](docs/DEVELOPER.md#getting-started) for prerequisites.

## Pull Request Workflow

1. **Sync `main`**
   ```bash
   git switch main
   git pull origin main
   ```
2. **Branch** with a conventional prefix (`feat/`, `fix/`, `docs/`, `chore/`, `refactor/`, `perf/`, `test/`, `ci/`):
   ```bash
   git switch -c fix/short-description
   ```
3. **Implement + test locally.** See [Local Checks](#local-checks) below.
4. **Commit** using [Conventional Commits](https://www.conventionalcommits.org/):
   - `feat: ZC#### — detect <pattern>`
   - `fix: ZC#### false positive on <case>`
   - `docs: update USER_GUIDE for inline directives`
   - `ci: tighten golangci timeout`
   - `chore: bump go-release action pin`
   - Commits are **GPG-signed**. `commit.gpgsign=true` or append `-S`.
5. **Push + PR:**
   ```bash
   git push -u origin <branch>
   gh pr create --fill
   ```
6. **Review:** CODEOWNERS (@afadesigns) must approve. All required checks (`test`, `security`, `sbom`) must pass. CI will reject unsigned commits.
7. **Merge:** maintainer squash-merges on green.

### Link issues
Use `Closes #N` / `Fixes #N` in the PR body so the issue auto-closes on merge.

## Local Checks

Before pushing, run:

```bash
go test -count=1 ./...
/srv/tools/go/bin/golangci-lint run ./...
go vet ./...
```

Or run `pre-commit run --all-files` if you have `pre-commit` installed — the project ships `.pre-commit-config.yaml` and `.pre-commit-hooks.yaml` covering lint, format, tests, and a trace scan.

Fuzz tests are time-boxed; run them when touching lexer/parser:

```bash
go test -fuzz=FuzzLexer -fuzztime=10s ./pkg/lexer
go test -fuzz=FuzzParser -fuzztime=10s ./pkg/parser
```

## Adding a New Kata

A kata is a Zsh-specific detection rule. Full scaffold + conventions live in the [Developer Guide — Creating a New Kata](docs/DEVELOPER.md#creating-a-new-kata). Short form:

1. Pick the next ID: `ls pkg/katas/zc*.go | sort | tail -1`
2. Create `pkg/katas/zc<NNNN>.go` registering the kata.
3. Create `pkg/katas/katatests/zc<NNNN>_test.go` with valid + invalid fixtures.
4. Once committed, **fix — don't remove** a kata. Retire duplicates as no-op stubs (pattern: `ZC1018`, `ZC1022`).

### Kata Conventions

- **Zsh-specific only.** Reject generic POSIX-sh anti-patterns — ShellCheck covers those.
- **Severity required.** One of `SeverityError`, `SeverityWarning`, `SeverityInfo`, `SeverityStyle`. See [Severity Levels](docs/USER_GUIDE.md#severity-levels).
- **Never `panic()` in `Check`.** Use `ok`-checked type assertions. A kata panic kills the linter.
- **No duplicates.** Grep existing katas before writing a new one.
- **Backtick-quote shell syntax** in titles, descriptions, and messages. End sentences with a period.

## Adding an Auto-Fix

Any kata whose rewrite is context-free, deterministic, and preserves semantics can ship a `Fix`. The registry entry takes an optional `Fix` hook alongside `Check`:

```go
Fix: func(node ast.Node, v Violation, source []byte) []FixEdit {
    // return [] or nil to skip; never panic
}
```

Each returned `FixEdit` is a 1-based `Line` / `Column`, a byte `Length` to replace at that position, and a `Replace` string. Multiple edits in one kata are allowed but must not overlap each other.

**Checklist before shipping a `Fix`:**

- **Idempotent on re-run.** Running `zshellcheck -fix` twice in a row must produce the same output. The second pass is re-parsed, so arrange detection so it won't fire on the rewritten source.
- **Byte-exact preservation.** Source bytes outside the rewritten span stay unchanged — no whitespace normalization, no quote changes, no re-indent.
- **Detector gates the shape.** If the rewrite only fits specific arg shapes, guard that in the `Check` (e.g. ZC1128 only fires on `touch` with exactly one non-flag arg).
- **Conflict-aware.** If another kata can fire at the same position, yield or narrow the detector — see ZC1055 ↔ ZC1079 for how `[[ == "" ]]` vs unquoted-RHS coordinate.
- **Integration test.** Add a case to `pkg/fix/integration_test.go` covering the rewrite and an idempotent already-fixed input. Use `runFix(t, src)` from the test helper.
- **Never bypass safety flags.** Hook failures stay failures — don't use `--no-verify` / `--no-gpg-sign`.
- **Regenerate `KATAS.md`** when a kata gains or loses a `Fix`: `go run ./internal/tools/gen-katas-md`. CI will fail otherwise.

## Security

Do not file vulnerabilities as public issues. See [SECURITY.md](SECURITY.md) for the reporting process.

## Labels

| Label | Meaning |
|---|---|
| `feat` | New feature or significant enhancement |
| `fix` | Bug fix |
| `docs` | Documentation change |
| `ci` | CI/CD change |
| `deps` | Dependency bump |
| `refactor` | Restructuring without behavior change |
| `perf` | Performance improvement |
| `test` | Test additions or fixes |
| `chore` | Maintenance |
| `starter` | Good first issue |
| `help wanted` | Needs community input |
| `duplicate` | Supersedes another issue/PR |
