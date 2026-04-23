## Description

<!--
Summarise the change. If this PR adds a new kata, describe the pattern
it detects and the Zsh semantics behind it.
-->

Closes # (issue)

## Type of change

- [ ] `feat`: new feature or enhancement
- [ ] `fix`: bug fix
- [ ] `docs`: documentation
- [ ] `ci`: CI/CD
- [ ] `deps`: dependency bump
- [ ] `refactor`: restructuring without behavior change
- [ ] `perf`: performance
- [ ] `test`: test additions or fixes
- [ ] `chore`: maintenance

## Checklist

- [ ] Read [CONTRIBUTING.md](../CONTRIBUTING.md).
- [ ] `go test -count=1 ./...` passes locally.
- [ ] `golangci-lint run ./...` clean.
- [ ] Relevant documentation updated (`README.md`, `docs/USER_GUIDE.md`, `docs/DEVELOPER.md`, `CHANGELOG.md`).
- [ ] Commits are GPG-signed.

## For new katas

- [ ] Detection file at `pkg/katas/zc<NNNN>.go` (self-registers via `init()` — no central file edit needed).
- [ ] Test file at `pkg/katas/katatests/zc<NNNN>_test.go` with both a violation case and a no-violation case.
- [ ] `Severity` set on the `Kata{…}` literal (`SeverityError` / `Warning` / `Info` / `Style`).
- [ ] Pattern is **Zsh-specific** — generic POSIX-sh anti-patterns belong in ShellCheck.
- [ ] Grepped existing katas for overlap: `grep -rn 'Title:' pkg/katas/ | grep -i '<keyword>'`.
- [ ] `Check` function uses `ok`-checked type assertions; never calls `panic()`.
