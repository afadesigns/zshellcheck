# Contributing to ZShellCheck

Thanks for helping improve ZShellCheck.
This guide covers the PR workflow, how to add a kata, and the local checks to run before pushing.

For deeper internals — lexer, parser, AST design, release process, architecture diagrams — see the [developer guide](docs/DEVELOPER.md).

## Quick start

```bash
git clone https://github.com/afadesigns/zshellcheck.git
cd zshellcheck
./install.sh
```

The installer builds from source when run inside the repo, or downloads the signed release binary otherwise.
See [developer guide → getting started](docs/DEVELOPER.md#getting-started) for prerequisites.

## Pull request workflow

1. Sync `main`.
   ```bash
   git switch main
   git pull origin main
   ```
2. Branch with a conventional prefix (`feat/`, `fix/`, `docs/`, `chore/`, `refactor/`, `perf/`, `test/`, `ci/`).
   ```bash
   git switch -c fix/short-description
   ```
3. Implement and test locally.
   See [local checks](#local-checks).
4. Commit using [Conventional Commits](https://www.conventionalcommits.org/).
   - `feat: ZC#### — detect <pattern>`
   - `fix: ZC#### false positive on <case>`
   - `docs: update USER_GUIDE for inline directives`
   - `ci: tighten golangci timeout`
   - `chore: bump go-release action pin`

   Commits must be GPG-signed.
   Set `commit.gpgsign=true` or pass `-S`.

   Commits must also carry a [Developer Certificate of Origin](https://developercertificate.org/) sign-off line.
   Pass `-s` (or set `format.signOff = true` in your git config) so each commit ends with `Signed-off-by: Your Name <you@example.com>`.
   The sign-off certifies that you wrote the patch or otherwise have the right to submit it under the project's MIT license.
5. Push and open the PR.
   ```bash
   git push -u origin <branch>
   gh pr create --fill
   ```
6. Review.
   CODEOWNERS (@afadesigns) must approve.
   Required checks (`test`, `security`, `sbom`) must pass.
   CI rejects unsigned commits.
7. Merge.
   The maintainer squash-merges on green.

### Linking issues

Use `Closes #N` or `Fixes #N` in the PR body so the issue auto-closes on merge.

## Local checks

Before pushing, run:

```bash
go test -count=1 ./...
golangci-lint run ./...
go vet ./...
```

The project ships `.pre-commit-config.yaml` and `.pre-commit-hooks.yaml` covering lint, format, tests, and a trace scan; running `pre-commit run --all-files` exercises the same gates.

Fuzz tests are time-boxed; run them when touching the lexer or parser:

```bash
go test -fuzz=FuzzLexer -fuzztime=10s ./pkg/lexer
go test -fuzz=FuzzParser -fuzztime=10s ./pkg/parser
```

## Adding a new kata

A kata is a Zsh-specific detection rule.
The full scaffold and conventions live in the [developer guide → creating a new kata](docs/DEVELOPER.md#creating-a-new-kata).

Short form:

1. Pick the next ID: `ls pkg/katas/zc*.go | sort | tail -1`.
2. Create `pkg/katas/zc<NNNN>.go` registering the kata.
3. Create `pkg/katas/katatests/zc<NNNN>_test.go` with valid and invalid fixtures.
4. Once committed, **fix the kata, do not remove it**.
   Retire duplicates as no-op stubs; see `ZC1018` and `ZC1022` for the pattern.

### Kata conventions

- **Zsh-specific only.**
  Reject generic POSIX-sh anti-patterns; ShellCheck covers those.
- **Severity required.**
  One of `SeverityError`, `SeverityWarning`, `SeverityInfo`, or `SeverityStyle`.
  See [severity levels](docs/USER_GUIDE.md#severity-levels).
- **Never `panic()` inside `Check`.**
  Use `ok`-checked type assertions.
  A kata panic kills the linter.
- **No duplicates.**
  Grep existing katas before writing a new one.
- **Backtick-quote shell syntax** in titles, descriptions, and messages.
  End sentences with a period.

### Coding standards

Go code follows [Effective Go](https://go.dev/doc/effective_go) and the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).
Compliance is enforced automatically by [.golangci.yml](.golangci.yml), which runs gofumpt, govet, errcheck, staticcheck, gosec, ineffassign, unparam, gocyclo, dupl, revive, thelper, unconvert, and prealloc on every push and PR.

Every Go source file carries a copyright statement and an SPDX licence identifier on the two lines immediately preceding the `package` declaration:

```go
// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.

package katas
```

### Code review standards

Reviewers check every PR against the rubric below.
A PR is mergeable when each item passes or has an explicit waiver in the PR description.

- **Scope.** Single concern per PR. Title matches the [Conventional Commits](https://www.conventionalcommits.org/) prefix.
- **Tests.** Behaviour changes ship with regression tests. Coverage on the touched files does not regress.
- **Style.** `gofumpt`, `goimports`, and the full `golangci-lint` suite pass with zero violations.
- **Safety.** No `panic()` outside `main` or test files; `Check` functions use `ok`-checked type assertions; no new `unsafe` imports; no shell-out from linter logic.
- **Public-surface hygiene.** Exported identifiers carry doc comments; backward-incompatible changes are flagged in CHANGELOG and require a version-major or version-minor bump.
- **Trace hygiene.** No banned tokens (claude, anthropic, gemini, etc.) anywhere in tracked surface; the trace-scan gate passes.
- **Headers.** Every new Go file carries the SPDX and copyright lines from the [coding standards](#coding-standards) section.
- **Security.** No leaked secrets, no plaintext credentials, no broken cryptography. Vulnerability fixes follow the disclosure process in [SECURITY.md](SECURITY.md).
- **Documentation.** User-visible changes update README, INSTALL, USER_GUIDE, or KATAS.md as appropriate; KATAS.md is regenerated rather than hand-edited.
- **Signatures.** Every commit is GPG-signed and carries a DCO `Signed-off-by:` trailer.

### Test policy

Formal policy: every PR that adds or modifies non-trivial functionality must add or update tests in the same change.
For new katas, this is enforced structurally by the kata-test pairing rule above.
For non-kata code, the `test` required check fails when coverage drops; reviewers reject PRs that ship behaviour without a regression test.
Bug fixes ship a regression test that fails before the fix and passes after.

### Adding an auto-fix

A kata becomes auto-fixable when its rewrite is context-free, idempotent, and byte-exact.
When any condition fails, leave `Fix` nil and ship detection-only.

1. Set the `Fix` field on the kata struct.

   ```go
   RegisterKata(ast.SimpleCommandNode, Kata{
       ID:       "ZCXXXX",
       Title:    "...",
       Severity: SeverityWarning,
       Check:    checkZCXXXX,
       Fix:      fixZCXXXX,
   })
   ```

2. Implement `fixZCXXXX(node ast.Node, v Violation, source []byte) []FixEdit`.
   `FixEdit` carries a 1-based `Line` and `Column`, a byte `Length`, and the replacement string.
   `pkg/katas/fixutil.go` exposes helpers including `LineColToByteOffset`.
3. Re-confirm the rewrite is safe across whitespace, quoting, and trailing-comment variants.
   The fixer runs multi-pass (up to five iterations) so nested rewrites resolve in a single invocation.
4. Add a fix-side test in `pkg/katas/katatests/zcXXXX_test.go` covering at least one applied-edit case and one no-op case.
5. Re-run `go run ./internal/tools/gen-katas-md` to refresh `KATAS.md`.
   The new entry reports `Auto-fix: yes` and the summary count increments.

Reference rewrite shapes already in the catalog:

| Pattern | Example |
| --- | --- |
| Token substitution (single byte span) | `ZC1002` `` `cmd` `` → `$(cmd)` |
| Identifier rename | `ZC1005` `which` → `whence` |
| Command and flag collapse | `ZC1355` `echo -E …` → `print -r …` |
| Parameter-name rename | `ZC1313` `$BASH_ALIASES` → `$aliases` |
| Quote-insertion around an expansion | `ZC1075` `rm -rf $var` → `rm -rf "$var"` |

When a rewrite does not fit one of these shapes, document the new pattern in this list in the same PR.

## Helping with distribution

Two install channels benefit from community submission.
A third party filing the upstream PR avoids the higher self-submission notability bar.

- **Homebrew (homebrew-core).**
  Third-party submissions face the lower bar (≥75 stars / ≥30 forks / ≥30 watchers); self-submissions face ≥225 / ≥90 / ≥90.
  Once the repo crosses the lower bar, the `homebrew-eligibility` workflow opens a tracking issue.
  Community members can file the new-formula PR against [Homebrew/homebrew-core](https://github.com/Homebrew/homebrew-core) at that point — credit goes to the submitter.
  The author does not self-submit until the higher bar is met to keep the path of least resistance open.
- **AUR (`zshellcheck-bin`).**
  Anyone with an AUR account can host the package; coordination happens through an issue tagged `distribution`.
  Goreleaser writes the `PKGBUILD` automatically once the release workflow holds the credentials; until then, a community-maintained AUR package is welcome.

## Security

Do not file vulnerabilities as public issues.
See [SECURITY.md](SECURITY.md) for the reporting process.

## Labels

| Label | Meaning |
| --- | --- |
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
| `duplicate` | Supersedes another issue or PR |
